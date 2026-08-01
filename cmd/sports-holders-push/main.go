package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

type holderAgg struct {
	Wallet       string
	Shares       float64
	MaxShares    float64
	Markets      int
	Categories   map[string]int
	TopMarket    string
	TopMarketAmt float64
	Tier         string
	GlobalTier   string
	Smart        float64
	Bot          float64
	Risks        []string
	ScoreTrades  int
	ScoreLarge   int
	ScoreClosed  int
	ScoreROI     float64
	ScoreCopyT   int
	ScoreCopyROI float64
}

func main() {
	outPath := flag.String("out", "wallets.sports-holders-push.txt", "wallet output path")
	excludePath := flag.String("exclude_wallets", "", "comma-separated wallet files to exclude")
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet score JSON used to annotate tier/smart/bot")
	categories := flag.String("target_categories", "basketball,soccer,esports", "comma-separated target categories")
	listName := flag.String("list", "sports_holders_push", "wallet list metadata name")
	scoreOnly := flag.Bool("football_score_only", false, "inspect football correct-score markets only")
	marketLimit := flag.Int("markets", 300, "maximum active markets to inspect")
	holderLimit := flag.Int("holders", 100, "holders per market")
	maxWallets := flag.Int("max_wallets", 200, "maximum wallets to write")
	minShares := flag.Float64("min_shares", 1000, "minimum aggregate holder shares")
	timeout := flag.Duration("timeout", 5*time.Minute, "overall timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	cfg := walletdiscover.DefaultConfig()
	cfg.TargetCategories = *categories
	client := walletdiscover.NewClient(cfg)
	excluded := loadWalletSet(*excludePath)
	scores := loadScoreSet(*scoresPath)

	markets, err := listTargetMarkets(ctx, client, *marketLimit, *categories, *scoreOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-holders-push: markets: %v\n", err)
		os.Exit(1)
	}

	aggs := map[string]*holderAgg{}
	for i, m := range markets {
		resp, err := client.Holders(ctx, m.ConditionID, *holderLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sports-holders-push: holders market=%s: %v\n", m.ConditionID, err)
			continue
		}
		allowedTokens := positiveOutcomeTokens(m)
		if *scoreOnly && len(allowedTokens) == 0 {
			fmt.Fprintf(os.Stderr, "sports-holders-push: score market has no positive outcome token market=%s\n", m.ConditionID)
			continue
		}
		category := walletdiscover.MarketTargetCategory(m)
		for _, token := range resp {
			if *scoreOnly && !allowedTokens[strings.TrimSpace(token.Token)] {
				continue
			}
			for _, h := range token.Holders {
				addr := normalizeAddress(h.ProxyWallet)
				if addr == "" || excluded[addr] {
					continue
				}
				agg := aggs[addr]
				if agg == nil {
					agg = &holderAgg{Wallet: addr, Categories: map[string]int{}}
					aggs[addr] = agg
				}
				agg.Shares += h.Amount
				agg.Markets++
				agg.Categories[category]++
				if h.Amount > agg.MaxShares {
					agg.MaxShares = h.Amount
				}
				if h.Amount > agg.TopMarketAmt {
					agg.TopMarketAmt = h.Amount
					agg.TopMarket = m.Question
				}
			}
		}
		if (i+1)%50 == 0 {
			fmt.Printf("sports-holders-push.progress markets=%d wallets=%d\n", i+1, len(aggs))
		}
	}

	rows := make([]*holderAgg, 0, len(aggs))
	for _, agg := range aggs {
		if agg.Shares >= *minShares {
			if score, ok := scores[agg.Wallet]; ok {
				agg.Tier = score.Tier
				agg.GlobalTier = score.Tier
				agg.Smart = score.SmartMoneyScore
				agg.Bot = score.BotScore
				agg.Risks = score.RiskFlags
				agg.ScoreTrades = score.Stats.FootballScoreTrades
				agg.ScoreLarge = score.Stats.FootballScoreLargeTrades
				agg.ScoreClosed = score.Stats.FootballScoreClosed
				agg.ScoreROI = score.Stats.FootballScoreClosedROI
				agg.ScoreCopyT = score.Stats.FootballScoreCopyClosed
				agg.ScoreCopyROI = score.Stats.FootballScoreCopyROI
				if *scoreOnly {
					agg.Tier = walletdiscover.RateFootballScoreWallet(score).Tier
				}
			}
			rows = append(rows, agg)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Shares != rows[j].Shares {
			return rows[i].Shares > rows[j].Shares
		}
		if rows[i].Markets != rows[j].Markets {
			return rows[i].Markets > rows[j].Markets
		}
		return rows[i].Wallet < rows[j].Wallet
	})
	if len(rows) > *maxWallets {
		rows = rows[:*maxWallets]
	}
	if err := writeWallets(*outPath, rows, len(markets), *minShares, *listName); err != nil {
		fmt.Fprintf(os.Stderr, "sports-holders-push: write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sports-holders-push done: markets=%d wallets=%d min_shares=%.0f out=%s\n", len(markets), len(rows), *minShares, *outPath)
}

func positiveOutcomeTokens(m walletdiscover.Market) map[string]bool {
	ids := m.ClobTokenIDs()
	outcomes := m.Outcomes()
	out := map[string]bool{}
	if len(ids) != len(outcomes) {
		return out
	}
	for i, outcome := range outcomes {
		switch strings.ToLower(strings.TrimSpace(outcome)) {
		case "", "no", "false", "none", "other":
			continue
		default:
			out[strings.TrimSpace(ids[i])] = true
		}
	}
	return out
}

func listTargetMarkets(ctx context.Context, client *walletdiscover.Client, max int, categories string, scoreOnly bool) ([]walletdiscover.Market, error) {
	var out []walletdiscover.Market
	now := time.Now()
	for offset := 0; len(out) < max; {
		limit := 100
		batch, err := client.ListMarkets(ctx, limit, offset)
		if err != nil {
			if strings.Contains(err.Error(), "offset too large") {
				break
			}
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		for _, m := range batch {
			if scoreOnly && !walletdiscover.IsFootballScoreMarket(m) {
				continue
			}
			if walletdiscover.GoodDiscoveryMarket(m, now) && walletdiscover.TargetCategoryAllowed(walletdiscover.MarketTargetCategory(m), categories) {
				out = append(out, m)
				if len(out) >= max {
					break
				}
			}
		}
		offset += len(batch)
		if len(batch) < limit {
			break
		}
	}
	return out, nil
}

func loadWalletSet(path string) map[string]bool {
	out := map[string]bool{}
	for _, part := range strings.Split(path, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		b, err := os.ReadFile(part)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			fields := strings.Fields(strings.Split(line, "#")[0])
			if len(fields) == 0 {
				continue
			}
			addr := normalizeAddress(fields[0])
			if addr != "" {
				out[addr] = true
			}
		}
	}
	return out
}

func loadScoreSet(path string) map[string]walletdiscover.WalletScore {
	out := map[string]walletdiscover.WalletScore{}
	if path == "" {
		return out
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	var rows []walletdiscover.WalletScore
	if err := json.Unmarshal(b, &rows); err != nil {
		return out
	}
	for _, row := range rows {
		addr := normalizeAddress(row.Address)
		if addr != "" {
			out[addr] = row
		}
	}
	return out
}

func writeWallets(path string, rows []*holderAgg, markets int, minShares float64, list string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# generated by sports-holders-push at %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&b, "# list=%s rows=%d markets=%d minShares=%.0f\n", list, len(rows), markets, minShares)
	for _, row := range rows {
		fmt.Fprintf(&b, "%s # list=%s tier=%s globalTier=%s smart=%.1f bot=%.1f scoreTrades=%d scoreLarge=%d scoreClosed=%d scoreROI=%.1f scoreCopyT=%d scoreCopyROI=%.1f shares=%.0f maxShares=%.0f markets=%d categories=%s risks=%s top=%s\n",
			row.Wallet, list, dash(row.Tier), dash(row.GlobalTier), row.Smart, row.Bot,
			row.ScoreTrades, row.ScoreLarge, row.ScoreClosed, row.ScoreROI, row.ScoreCopyT, row.ScoreCopyROI,
			row.Shares, row.MaxShares, row.Markets, categorySummary(row.Categories), riskSummary(row.Risks), cleanComment(row.TopMarket))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func normalizeAddress(addr string) string {
	addr = strings.ToLower(strings.TrimSpace(addr))
	if len(addr) != 42 || !strings.HasPrefix(addr, "0x") {
		return ""
	}
	for _, ch := range addr[2:] {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') {
			return ""
		}
	}
	return addr
}

func categorySummary(in map[string]int) string {
	var keys []string
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s:%d", k, in[k]))
	}
	return strings.Join(parts, ",")
}

func cleanComment(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "#", "")
	if len(s) > 80 {
		return s[:80]
	}
	return s
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "?"
	}
	return strings.TrimSpace(s)
}

func riskSummary(risks []string) string {
	if len(risks) == 0 {
		return "-"
	}
	return strings.Join(risks, ",")
}
