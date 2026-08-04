package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

type walletMeta struct {
	List  string
	Tier  string
	Smart float64
	Bot   float64
}

type tapeTrade struct {
	Time           time.Time `json:"time"`
	Wallet         string    `json:"wallet"`
	Side           string    `json:"side"`
	Notional       float64   `json:"notional"`
	Price          float64   `json:"price"`
	Size           float64   `json:"size"`
	Outcome        string    `json:"outcome"`
	Market         string    `json:"market"`
	Slug           string    `json:"slug"`
	Category       string    `json:"category"`
	Asset          string    `json:"asset"`
	ConditionID    string    `json:"condition_id"`
	Transaction    string    `json:"transaction"`
	KnownList      string    `json:"known_list,omitempty"`
	Tier           string    `json:"tier,omitempty"`
	Smart          float64   `json:"smart,omitempty"`
	Bot            float64   `json:"bot,omitempty"`
	TargetCopyROI  float64   `json:"target_copy_roi,omitempty"`
	TargetCopyT    int       `json:"target_copy_t,omitempty"`
	TargetCopyPnL  float64   `json:"target_copy_pnl,omitempty"`
	TargetTradePct float64   `json:"target_trade_pct,omitempty"`
}

type walletAgg struct {
	Wallet        string
	KnownList     string
	Tier          string
	Smart         float64
	Bot           float64
	Buys          int
	Sells         int
	BuyNotional   float64
	SellNotional  float64
	MaxBuy        float64
	Markets       map[string]struct{}
	Categories    map[string]struct{}
	TargetCopyROI float64
	TargetCopyT   int
	TargetCopyPnL float64
}

func main() {
	output := flag.String("output", "db/strategy_iteration/sports_tape.jsonl", "JSONL output path")
	report := flag.String("report", "reports/sports_tape.md", "markdown report path")
	walletsOut := flag.String("wallets_out", "wallets.sports-tape.txt", "wallet seed output path")
	walletsMax := flag.Int("wallets_max", 50, "maximum wallets to write to wallet seed output")
	walletsMinBuyNotional := flag.Float64("wallets_min_buy_notional", 500, "minimum aggregate buy notional for wallet seed output")
	pushWallets := flag.String("push_wallets", "wallets.strategy-push.txt", "known push wallet file")
	walletStatuses := flag.String("wallet_statuses", "", "comma-separated wallet status files that override tape labels in order")
	excludeWallets := flag.String("exclude_wallets", "", "comma-separated wallet files to exclude from tape scan and retained output")
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet scores JSON")
	targetCategories := flag.String("target_categories", "basketball,soccer,esports", "comma-separated target categories")
	pages := flag.Int("pages", 7, "recent trade pages to scan")
	limit := flag.Int("limit", 500, "recent trade page size")
	minNotional := flag.Float64("min_notional", 500, "minimum trade notional USD")
	minPrice := flag.Float64("min_price", 0.05, "minimum price")
	maxPrice := flag.Float64("max_price", 0.95, "maximum price")
	topN := flag.Int("top", 25, "top wallets/trades to render")
	timeout := flag.Duration("timeout", 3*time.Minute, "overall timeout")
	allowEmpty := flag.Bool("allow_empty", false, "allow empty scans to overwrite existing tape outputs")
	retainInput := flag.String("retain_input", "", "existing sports tape JSONL to merge before writing; default output path")
	retainWindow := flag.Duration("retain_window", 6*time.Hour, "retain prior qualifying sports tape trades for this long; 0 disables retention")
	flag.Parse()
	if *retainInput == "" {
		*retainInput = *output
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	cfg := walletdiscover.DefaultConfig()
	cfg.TargetCategories = *targetCategories
	cfg.MinNotionalUSD = *minNotional
	cfg.MinPrice = *minPrice
	cfg.MaxPrice = *maxPrice
	client := walletdiscover.NewClient(cfg)
	metas := loadPushMetas(*pushWallets)
	mergeWalletMetas(metas, loadPushMetas(*walletStatuses))
	excluded := loadWalletSet(*excludeWallets)
	scores := loadScoreMetas(*scoresPath)

	var trades []tapeTrade
	seen := map[string]struct{}{}
	for page := 0; page < *pages; page++ {
		batch, err := client.RecentTrades(ctx, *limit, page*(*limit))
		if err != nil {
			if walletdiscover.IsMaxHistoricalOffsetError(err) {
				fmt.Fprintf(os.Stderr, "sports-tape: reached recent trades offset limit at page %d; preserving %d qualifying trades\n", page+1, len(trades))
				break
			}
			fmt.Fprintf(os.Stderr, "sports-tape: recent trades page %d: %v\n", page+1, err)
			break
		}
		for _, tr := range batch {
			if !walletdiscover.QualifyingTrade(tr, cfg, nil) {
				continue
			}
			cat := walletdiscover.TradeTargetCategory(tr)
			addr := strings.ToLower(strings.TrimSpace(tr.ProxyWallet))
			if addr == "" {
				continue
			}
			if _, ok := excluded[addr]; ok {
				continue
			}
			key := tradeKey(tr)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			slug := tr.EventSlug
			if slug == "" {
				slug = tr.Slug
			}
			meta := metas[addr]
			score := scores[addr]
			tape := tapeTrade{
				Time:           time.Unix(tr.Timestamp, 0),
				Wallet:         addr,
				Side:           strings.ToUpper(tr.Side),
				Notional:       round2(tr.NotionalUSD()),
				Price:          tr.Price,
				Size:           tr.Size,
				Outcome:        tr.Outcome,
				Market:         tr.Title,
				Slug:           slug,
				Category:       cat,
				Asset:          tr.Asset,
				ConditionID:    tr.ConditionID,
				Transaction:    tr.TransactionHash,
				KnownList:      meta.List,
				Tier:           firstNonEmpty(meta.Tier, score.Tier),
				Smart:          firstNonZero(meta.Smart, score.Smart),
				Bot:            firstNonZero(meta.Bot, score.Bot),
				TargetCopyROI:  score.TargetCopyROI,
				TargetCopyT:    score.TargetCopyT,
				TargetCopyPnL:  score.TargetCopyPnL,
				TargetTradePct: score.TargetTradePct,
			}
			trades = append(trades, tape)
		}
	}

	sort.Slice(trades, func(i, j int) bool {
		if trades[i].Time.Equal(trades[j].Time) {
			return trades[i].Notional > trades[j].Notional
		}
		return trades[i].Time.After(trades[j].Time)
	})

	retained := 0
	if *retainWindow > 0 {
		prior, err := loadTapeTrades(*retainInput)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "sports-tape: load retained tape: %v\n", err)
		}
		trades, retained = mergeRetainedTrades(trades, prior, cfg, excluded, time.Now(), *retainWindow)
	}
	applyWalletMetas(trades, metas, scores)

	if len(trades) == 0 && !*allowEmpty {
		walletsWritten := countWalletSeeds(*walletsOut)
		fmt.Printf("sports-tape done: trades=0 wallets=0 seed_wallets=%d output=%s report=%s wallets=%s preserved=true\n", walletsWritten, *output, *report, *walletsOut)
		return
	}

	if err := writeJSONL(*output, trades); err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape: write output: %v\n", err)
		os.Exit(1)
	}
	if err := writeReport(*report, trades, *minNotional, *targetCategories, *topN, *retainWindow, retained); err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape: write report: %v\n", err)
		os.Exit(1)
	}
	walletsWritten, err := writeWalletSeeds(*walletsOut, trades, *walletsMax, *walletsMinBuyNotional)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape: write wallet seeds: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sports-tape done: trades=%d retained=%d wallets=%d seed_wallets=%d output=%s report=%s wallets=%s\n", len(trades), retained, len(aggregate(trades)), walletsWritten, *output, *report, *walletsOut)
}

func tradeKey(tr walletdiscover.Trade) string {
	if tr.TransactionHash != "" {
		return strings.ToLower(strings.Join([]string{tr.TransactionHash, tr.Asset, fmt.Sprint(tr.Timestamp)}, "|"))
	}
	return strings.ToLower(fmt.Sprintf("%s|%s|%s|%d|%.8f|%.8f",
		tr.ProxyWallet, tr.Side, tr.Asset, tr.Timestamp, tr.Price, tr.Size))
}

func tapeKey(tr tapeTrade) string {
	if tr.Transaction != "" {
		return strings.ToLower(strings.Join([]string{tr.Transaction, tr.Asset, fmt.Sprint(tr.Time.Unix())}, "|"))
	}
	return strings.ToLower(fmt.Sprintf("%s|%s|%s|%d|%.8f|%.8f",
		tr.Wallet, tr.Side, tr.Asset, tr.Time.Unix(), tr.Price, tr.Size))
}

func loadTapeTrades(path string) ([]tapeTrade, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []tapeTrade
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 8*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var tr tapeTrade
		if err := json.Unmarshal([]byte(line), &tr); err != nil {
			continue
		}
		if tr.Wallet == "" || tr.Time.IsZero() {
			continue
		}
		out = append(out, tr)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func mergeRetainedTrades(current, prior []tapeTrade, cfg walletdiscover.Config, excluded map[string]struct{}, now time.Time, retainWindow time.Duration) ([]tapeTrade, int) {
	if retainWindow <= 0 || len(prior) == 0 {
		return current, 0
	}
	seen := map[string]struct{}{}
	for _, tr := range current {
		seen[tapeKey(tr)] = struct{}{}
	}
	cutoff := now.Add(-retainWindow)
	retained := 0
	for _, tr := range prior {
		if _, ok := excluded[strings.ToLower(tr.Wallet)]; ok {
			continue
		}
		if tr.Time.Before(cutoff) || !qualifyingTapeTrade(tr, cfg) {
			continue
		}
		key := tapeKey(tr)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		current = append(current, tr)
		retained++
	}
	sort.Slice(current, func(i, j int) bool {
		if current[i].Time.Equal(current[j].Time) {
			return current[i].Notional > current[j].Notional
		}
		return current[i].Time.After(current[j].Time)
	})
	return current, retained
}

func loadWalletSet(path string) map[string]struct{} {
	out := map[string]struct{}{}
	if strings.TrimSpace(path) == "" {
		return out
	}
	if strings.Contains(path, ",") {
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			for wallet := range loadWalletSet(part) {
				out[wallet] = struct{}{}
			}
		}
		return out
	}
	wallets, err := walletdiscover.LoadWallets(path)
	if err != nil {
		return out
	}
	for _, wallet := range wallets {
		out[strings.ToLower(wallet)] = struct{}{}
	}
	return out
}

func qualifyingTapeTrade(tr tapeTrade, cfg walletdiscover.Config) bool {
	return walletdiscover.QualifyingTrade(walletdiscover.Trade{
		ProxyWallet:     tr.Wallet,
		Side:            tr.Side,
		Asset:           tr.Asset,
		ConditionID:     tr.ConditionID,
		Size:            tr.Size,
		Price:           tr.Price,
		Timestamp:       tr.Time.Unix(),
		Title:           tr.Market,
		Slug:            tr.Slug,
		EventSlug:       tr.Slug,
		Outcome:         tr.Outcome,
		TransactionHash: tr.Transaction,
		Type:            "TRADE",
	}, cfg, nil)
}

type scoreMeta struct {
	Tier           string
	Smart          float64
	Bot            float64
	TargetCopyROI  float64
	TargetCopyT    int
	TargetCopyPnL  float64
	TargetTradePct float64
}

func loadScoreMetas(path string) map[string]scoreMeta {
	out := map[string]scoreMeta{}
	b, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	var scores []walletdiscover.WalletScore
	if err := json.Unmarshal(b, &scores); err != nil {
		return out
	}
	for _, s := range scores {
		out[strings.ToLower(s.Address)] = scoreMeta{
			Tier:           s.Tier,
			Smart:          s.SmartMoneyScore,
			Bot:            s.BotScore,
			TargetCopyROI:  s.Stats.TargetCopyROI,
			TargetCopyT:    s.Stats.TargetCopyClosed,
			TargetCopyPnL:  s.Stats.TargetCopyPnL,
			TargetTradePct: s.Stats.TargetTradeRatio * 100,
		}
	}
	return out
}

func loadPushMetas(path string) map[string]walletMeta {
	out := map[string]walletMeta{}
	if strings.TrimSpace(path) == "" {
		return out
	}
	if strings.Contains(path, ",") {
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			mergeWalletMetas(out, loadPushMetas(part))
		}
		return out
	}
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.SplitN(strings.TrimSpace(sc.Text()), "#", 2)
		fields := strings.Fields(parts[0])
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		meta := walletMeta{}
		if len(parts) > 1 {
			for _, kv := range strings.Fields(parts[1]) {
				k, v, ok := strings.Cut(kv, "=")
				if !ok {
					continue
				}
				switch k {
				case "list":
					meta.List = v
				case "tier":
					meta.Tier = v
				case "smart":
					if value, err := strconv.ParseFloat(v, 64); err == nil {
						meta.Smart = value
					}
				case "bot":
					if value, err := strconv.ParseFloat(v, 64); err == nil {
						meta.Bot = value
					}
				}
			}
		}
		out[addr] = meta
	}
	return out
}

func mergeWalletMetas(dst, src map[string]walletMeta) {
	for addr, meta := range src {
		dst[strings.ToLower(addr)] = meta
	}
}

func applyWalletMetas(trades []tapeTrade, metas map[string]walletMeta, scores map[string]scoreMeta) {
	for i := range trades {
		addr := strings.ToLower(strings.TrimSpace(trades[i].Wallet))
		meta := metas[addr]
		score := scores[addr]
		if meta.List != "" {
			trades[i].KnownList = meta.List
		}
		if tier := firstNonEmpty(meta.Tier, score.Tier); tier != "" {
			trades[i].Tier = tier
		}
		if v := firstNonZero(meta.Smart, score.Smart); v != 0 {
			trades[i].Smart = v
		}
		if v := firstNonZero(meta.Bot, score.Bot); v != 0 {
			trades[i].Bot = v
		}
		if score.TargetCopyROI != 0 {
			trades[i].TargetCopyROI = score.TargetCopyROI
		}
		if score.TargetCopyT != 0 {
			trades[i].TargetCopyT = score.TargetCopyT
		}
		if score.TargetCopyPnL != 0 {
			trades[i].TargetCopyPnL = score.TargetCopyPnL
		}
		if score.TargetTradePct != 0 {
			trades[i].TargetTradePct = score.TargetTradePct
		}
	}
}

func writeJSONL(path string, trades []tapeTrade) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, tr := range trades {
		if err := enc.Encode(tr); err != nil {
			return err
		}
	}
	return nil
}

func writeReport(path string, trades []tapeTrade, minNotional float64, categories string, topN int, retainWindow time.Duration, retained int) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	aggs := aggregate(trades)
	sort.Slice(aggs, func(i, j int) bool {
		if aggs[i].BuyNotional != aggs[j].BuyNotional {
			return aggs[i].BuyNotional > aggs[j].BuyNotional
		}
		return aggs[i].Buys > aggs[j].Buys
	})
	var b strings.Builder
	fmt.Fprintf(&b, "# Sports Whale Tape\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Target categories: %s\n", categories)
	fmt.Fprintf(&b, "- Min notional: $%.0f\n", minNotional)
	fmt.Fprintf(&b, "- Retain window: %s\n", retainWindow)
	fmt.Fprintf(&b, "- Retained prior trades: %d\n", retained)
	fmt.Fprintf(&b, "- Trades: %d\n", len(trades))
	fmt.Fprintf(&b, "- Wallets: %d\n\n", len(aggs))

	fmt.Fprintf(&b, "## Top Wallets\n\n")
	fmt.Fprintf(&b, "| Wallet | List | Tier | Buys | BuyNotional | MaxBuy | Mkts | Cat | Smart | Bot | TargetCopyROI | TargetCopyT |\n")
	fmt.Fprintf(&b, "|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for i, a := range aggs {
		if i >= topN {
			break
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %d | $%.0f | $%.0f | %d | %d | %.1f | %.1f | %.1f%% | %d |\n",
			shortAddr(a.Wallet), dash(a.KnownList), dash(a.Tier), a.Buys, a.BuyNotional, a.MaxBuy,
			len(a.Markets), len(a.Categories), a.Smart, a.Bot, a.TargetCopyROI, a.TargetCopyT)
	}

	fmt.Fprintf(&b, "\n## Largest Trades\n\n")
	fmt.Fprintf(&b, "| Time | Wallet | List | Side | Notional | Price | Outcome | Category | Market |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---:|---:|---|---|---|\n")
	bySize := append([]tapeTrade{}, trades...)
	sort.Slice(bySize, func(i, j int) bool { return bySize[i].Notional > bySize[j].Notional })
	for i, tr := range bySize {
		if i >= topN {
			break
		}
		fmt.Fprintf(&b, "| %s | `%s` | %s | %s | $%.0f | %.3f | %s | %s | %s |\n",
			tr.Time.Format("01-02 15:04"), shortAddr(tr.Wallet), dash(tr.KnownList), tr.Side, tr.Notional,
			tr.Price, trunc(tr.Outcome, 24), tr.Category, trunc(tr.Market, 70))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeWalletSeeds(path string, trades []tapeTrade, maxWallets int, minBuyNotional float64) (int, error) {
	if maxWallets <= 0 {
		return 0, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return 0, err
	}
	aggs := aggregate(trades)
	sort.Slice(aggs, func(i, j int) bool {
		if aggs[i].BuyNotional != aggs[j].BuyNotional {
			return aggs[i].BuyNotional > aggs[j].BuyNotional
		}
		return aggs[i].MaxBuy > aggs[j].MaxBuy
	})
	var lines []string
	for _, a := range aggs {
		if len(lines) >= maxWallets {
			break
		}
		if a.BuyNotional < minBuyNotional || a.Buys == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s # source=sports_tape list=%s tier=%s buys=%d buyNotional=$%.0f maxBuy=$%.0f targetCopyROI=%.1f%% targetCopyT=%d bot=%.1f",
			a.Wallet, dash(a.KnownList), dash(a.Tier), a.Buys, a.BuyNotional, a.MaxBuy, a.TargetCopyROI, a.TargetCopyT, a.Bot))
	}
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return len(lines), os.WriteFile(path, []byte(content), 0o644)
}

func countWalletSeeds(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(strings.Split(sc.Text(), "#")[0])
		if strings.HasPrefix(strings.ToLower(line), "0x") && len(line) == 42 {
			n++
		}
	}
	return n
}

func aggregate(trades []tapeTrade) []walletAgg {
	m := map[string]*walletAgg{}
	for _, tr := range trades {
		a := m[tr.Wallet]
		if a == nil {
			a = &walletAgg{
				Wallet:        tr.Wallet,
				KnownList:     tr.KnownList,
				Tier:          tr.Tier,
				Smart:         tr.Smart,
				Bot:           tr.Bot,
				Markets:       map[string]struct{}{},
				Categories:    map[string]struct{}{},
				TargetCopyROI: tr.TargetCopyROI,
				TargetCopyT:   tr.TargetCopyT,
				TargetCopyPnL: tr.TargetCopyPnL,
			}
			m[tr.Wallet] = a
		}
		a.Markets[tr.ConditionID] = struct{}{}
		a.Categories[tr.Category] = struct{}{}
		if tr.Side == "BUY" {
			a.Buys++
			a.BuyNotional += tr.Notional
			if tr.Notional > a.MaxBuy {
				a.MaxBuy = tr.Notional
			}
		} else if tr.Side == "SELL" {
			a.Sells++
			a.SellNotional += tr.Notional
		}
	}
	out := make([]walletAgg, 0, len(m))
	for _, a := range m {
		out = append(out, *a)
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonZero(vals ...float64) float64 {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func trunc(s string, n int) string {
	s = strings.ReplaceAll(s, "|", "/")
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func round2(v float64) float64 {
	if v < 0 {
		return -round2(-v)
	}
	return float64(int(v*100+0.5)) / 100
}
