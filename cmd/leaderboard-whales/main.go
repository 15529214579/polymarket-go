package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func main() {
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet score JSON from wallet-discover")
	pushWalletsPath := flag.String("push_wallets", "wallets.strategy-push.txt", "current whale push wallet file")
	excludeWalletsPath := flag.String("exclude_wallets", "", "comma-separated wallet files excluded from recommendations")
	reportPath := flag.String("report", "reports/leaderboard_whales.md", "markdown report path")
	recommendWalletsPath := flag.String("recommend_wallets", "wallets.leaderboard-watch.txt", "wallet file for clean leaderboard whales worth large-order monitoring")
	strictPushWalletsPath := flag.String("push_wallets_out", "wallets.leaderboard-push.txt", "wallet file for strict leaderboard whales eligible for whale push monitoring")
	topN := flag.Int("top", 25, "rows per section")
	minSmart := flag.Float64("min_smart", 70, "minimum smart-money score for recommended leaderboard whales")
	maxBot := flag.Float64("max_bot", 45, "maximum bot score for recommended leaderboard whales")
	minLarge := flag.Int("min_large", 20, "minimum large trades for recommended leaderboard whales")
	minAvgNotional := flag.Float64("min_avg_notional", 500, "minimum average notional for recommended leaderboard whales")
	minTargetTrades := flag.Int("min_target_trades", 1, "minimum target-category trades for recommended leaderboard whales")
	minTargetLarge := flag.Int("min_target_large", 0, "minimum target-category large trades for recommended leaderboard whales")
	whaleWatchMinSmart := flag.Float64("whale_watch_min_smart", 80, "minimum smart-money score for leaderboard whale-watch candidates")
	whaleWatchMaxBot := flag.Float64("whale_watch_max_bot", 45, "maximum bot score for leaderboard whale-watch candidates")
	whaleWatchMinLarge := flag.Int("whale_watch_min_large", 100, "minimum large trades for leaderboard whale-watch candidates")
	whaleWatchMinAvgNotional := flag.Float64("whale_watch_min_avg_notional", 300, "minimum average notional for leaderboard whale-watch candidates")
	whaleWatchMinTargetLarge := flag.Int("whale_watch_min_target_large", 20, "minimum target-category large trades for leaderboard whale-watch candidates")
	recommendLimit := flag.Int("recommend_limit", 50, "maximum wallets written to recommend_wallets")
	strictPushLimit := flag.Int("push_limit", 25, "maximum wallets written to push_wallets_out")
	strictPushMinTier := flag.String("push_min_tier", "B", "minimum wallet tier for push_wallets_out")
	strictPushMinSmart := flag.Float64("push_min_smart", 80, "minimum smart-money score for push_wallets_out")
	strictPushMaxBot := flag.Float64("push_max_bot", 35, "maximum bot score for push_wallets_out")
	strictPushMinLarge := flag.Int("push_min_large", 50, "minimum large trades for push_wallets_out")
	strictPushMinAvgNotional := flag.Float64("push_min_avg_notional", 1000, "minimum average notional for push_wallets_out")
	strictPushMinTargetTrades := flag.Int("push_min_target_trades", 5, "minimum target-category trades for push_wallets_out")
	strictPushMinTargetLarge := flag.Int("push_min_target_large", 1, "minimum target-category large trades for push_wallets_out")
	flag.Parse()

	scores, err := loadScores(*scoresPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: %v\n", err)
		os.Exit(1)
	}
	push, err := loadWalletSet(*pushWalletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: load push wallets: %v\n", err)
		os.Exit(1)
	}
	exclude, err := loadWalletSet(*excludeWalletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: load exclude wallets: %v\n", err)
		os.Exit(1)
	}
	scan := renderReport(scores, push, exclude, *scoresPath, *pushWalletsPath, *excludeWalletsPath, *topN, *minSmart, *maxBot, *minLarge, *minAvgNotional, *minTargetTrades, *minTargetLarge, *whaleWatchMinSmart, *whaleWatchMaxBot, *whaleWatchMinLarge, *whaleWatchMinAvgNotional, *whaleWatchMinTargetLarge, *strictPushMinTier, *strictPushMinSmart, *strictPushMaxBot, *strictPushMinLarge, *strictPushMinAvgNotional, *strictPushMinTargetTrades, *strictPushMinTargetLarge)
	if err := os.MkdirAll(filepath.Dir(*reportPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: mkdir report: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*reportPath, []byte(scan.Report), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: write report: %v\n", err)
		os.Exit(1)
	}
	if err := writeWalletFile(*recommendWalletsPath, limitScores(mergeScores(scan.Recommended, scan.WhaleWatch), *recommendLimit), "leaderboard_watch"); err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: write recommend wallets: %v\n", err)
		os.Exit(1)
	}
	if err := writeWalletFile(*strictPushWalletsPath, limitScores(scan.StrictPush, *strictPushLimit), "leaderboard_push"); err != nil {
		fmt.Fprintf(os.Stderr, "leaderboard-whales: write push wallets: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("leaderboard-whales done: leaderboard=%d pushed=%d excluded=%d strict_push=%d recommended=%d report=%s recommend_wallets=%s push_wallets_out=%s\n",
		countLeaderboard(scores), countLeaderboardPushed(scores, push), countLeaderboardPushed(scores, exclude), len(scan.StrictPush), len(scan.Recommended)+len(scan.WhaleWatch), *reportPath, *recommendWalletsPath, *strictPushWalletsPath)
}

func loadScores(path string) ([]walletdiscover.WalletScore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var scores []walletdiscover.WalletScore
	if err := json.Unmarshal(b, &scores); err != nil {
		return nil, err
	}
	return scores, nil
}

func loadWalletSet(path string) (map[string]string, error) {
	out := map[string]string{}
	if path == "" {
		return out, nil
	}
	if strings.Contains(path, ",") {
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			set, err := loadWalletSet(part)
			if err != nil {
				return nil, err
			}
			for addr, list := range set {
				out[addr] = list
			}
		}
		return out, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	for _, raw := range strings.Split(string(b), "\n") {
		line, comment, _ := strings.Cut(raw, "#")
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(strings.TrimSpace(fields[0]))
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		out[addr] = parseList(comment)
	}
	return out, nil
}

func parseList(comment string) string {
	for _, f := range strings.Fields(comment) {
		k, v, ok := strings.Cut(f, "=")
		if ok && k == "list" {
			return v
		}
	}
	return ""
}

type scanResult struct {
	Report      string
	Recommended []walletdiscover.WalletScore
	WhaleWatch  []walletdiscover.WalletScore
	StrictPush  []walletdiscover.WalletScore
}

func renderReport(scores []walletdiscover.WalletScore, push, exclude map[string]string, scoresPath, pushPath, excludePath string, topN int, minSmart, maxBot float64, minLarge int, minAvgNotional float64, minTargetTrades, minTargetLarge int, whaleWatchMinSmart, whaleWatchMaxBot float64, whaleWatchMinLarge int, whaleWatchMinAvgNotional float64, whaleWatchMinTargetLarge int, strictPushMinTier string, strictPushMinSmart, strictPushMaxBot float64, strictPushMinLarge int, strictPushMinAvgNotional float64, strictPushMinTargetTrades, strictPushMinTargetLarge int) scanResult {
	var leaderboard []walletdiscover.WalletScore
	var pushed []walletdiscover.WalletScore
	var excluded []walletdiscover.WalletScore
	var recommended []walletdiscover.WalletScore
	var whaleWatch []walletdiscover.WalletScore
	var strictPush []walletdiscover.WalletScore
	var watch []walletdiscover.WalletScore
	var bots []walletdiscover.WalletScore

	for _, s := range scores {
		if !hasLeaderboardSource(s.Sources) {
			continue
		}
		leaderboard = append(leaderboard, s)
		addr := strings.ToLower(s.Address)
		_, inPush := push[addr]
		_, isExcluded := exclude[addr]
		if inPush {
			pushed = append(pushed, s)
		}
		if isExcluded {
			excluded = append(excluded, s)
			continue
		}
		if isBotLike(s) {
			bots = append(bots, s)
			continue
		}
		if qualifiesStrictLeaderboardPush(s, strictPushMinTier, strictPushMinSmart, strictPushMaxBot, strictPushMinLarge, strictPushMinAvgNotional, strictPushMinTargetTrades, strictPushMinTargetLarge) {
			if !inPush {
				strictPush = append(strictPush, s)
			}
			continue
		}
		if qualifiesLeaderboardWhale(s, minSmart, maxBot, minLarge, minAvgNotional, minTargetTrades, minTargetLarge) {
			if !inPush {
				recommended = append(recommended, s)
			}
			continue
		}
		if qualifiesLeaderboardWhaleWatch(s, whaleWatchMinSmart, whaleWatchMaxBot, whaleWatchMinLarge, whaleWatchMinAvgNotional, minTargetTrades, whaleWatchMinTargetLarge) {
			if !inPush {
				whaleWatch = append(whaleWatch, s)
			}
			continue
		}
		if highValueWatch(s, minTargetTrades, minTargetLarge) && !inPush {
			watch = append(watch, s)
		}
	}

	sortScores(pushed, leaderboardWhaleScore)
	sortScores(excluded, leaderboardWhaleScore)
	sortScores(strictPush, leaderboardWhaleScore)
	sortScores(recommended, leaderboardWhaleScore)
	sortScores(whaleWatch, leaderboardWhaleScore)
	sortScores(watch, watchScore)
	sortScores(bots, botSortScore)

	tierCounts := map[string]int{}
	for _, s := range leaderboard {
		tierCounts[s.Tier]++
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Leaderboard Whale Scan\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Scores: `%s`\n", scoresPath)
	fmt.Fprintf(&b, "- Push wallets: `%s`\n", pushPath)
	fmt.Fprintf(&b, "- Exclude wallets: `%s`\n", dash(excludePath))
	fmt.Fprintf(&b, "- Leaderboard-origin wallets scored: %d\n", len(leaderboard))
	fmt.Fprintf(&b, "- Tiers: A=%d B=%d C=%d D=%d BOT=%d\n", tierCounts["A"], tierCounts["B"], tierCounts["C"], tierCounts["D"], tierCounts["BOT"])
	fmt.Fprintf(&b, "- Already in whale push: %d\n", len(pushed))
	fmt.Fprintf(&b, "- Excluded by quarantine/review-noise: %d\n", len(excluded))
	fmt.Fprintf(&b, "- New strict push candidates: %d\n", len(strictPush))
	fmt.Fprintf(&b, "- New recommended watch candidates: %d\n", len(recommended))
	fmt.Fprintf(&b, "- New leaderboard whale-watch candidates: %d\n", len(whaleWatch))
	fmt.Fprintf(&b, "- Recommended target-category requirement: targetTrades>=%d targetLarge>=%d\n", minTargetTrades, minTargetLarge)
	fmt.Fprintf(&b, "- Whale-watch requirement: smart>=%.0f bot<%.0f large>=%d targetLarge>=%d avgNotional>=$%.0f\n", whaleWatchMinSmart, whaleWatchMaxBot, whaleWatchMinLarge, whaleWatchMinTargetLarge, whaleWatchMinAvgNotional)
	fmt.Fprintf(&b, "- Strict push target-category requirement: targetTrades>=%d targetLarge>=%d\n", strictPushMinTargetTrades, strictPushMinTargetLarge)
	fmt.Fprintf(&b, "- Bot/flow filtered: %d\n\n", len(bots))

	writeTable(&b, "Leaderboard Whales Already in Push", pushed, push, topN, true)
	writeTable(&b, "Excluded Leaderboard Whales", excluded, exclude, topN, true)
	writeTable(&b, "Strict Leaderboard Push Candidates", strictPush, push, topN, false)
	writeTable(&b, "Recommended Leaderboard Whales", recommended, push, topN, false)
	writeTable(&b, "Leaderboard Whale Watch", whaleWatch, push, topN, false)
	writeTable(&b, "High-Value Watch Only", watch, push, topN, false)
	writeTable(&b, "Filtered Bot-Like Leaderboard Wallets", bots, push, topN, false)
	return scanResult{Report: b.String(), Recommended: recommended, WhaleWatch: whaleWatch, StrictPush: strictPush}
}

func writeTable(b *strings.Builder, title string, rows []walletdiscover.WalletScore, push map[string]string, topN int, includeList bool) {
	fmt.Fprintf(b, "## %s\n\n", title)
	fmt.Fprintf(b, "- Wallets shown: %d\n\n", minInt(len(rows), topN))
	if len(rows) == 0 {
		fmt.Fprintf(b, "No wallets in this section.\n\n")
		return
	}
	if includeList {
		fmt.Fprintf(b, "| Wallet | List | Tier | Smart | Bot | WhaleScore | Large | TargetT | TargetLarge | AvgNotional | CopyROI | CopyT | Sources | Risks |\n")
		fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	} else {
		fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | WhaleScore | Large | TargetT | TargetLarge | AvgNotional | CopyROI | CopyT | Sources | Risks |\n")
		fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	}
	for i, s := range rows {
		if i >= topN {
			break
		}
		if includeList {
			fmt.Fprintf(b, "| `%s` | %s | %s | %.1f | %.1f | %.1f | %d | %d | %d | $%.0f | %.1f%% | %d | %s | %s |\n",
				s.Address, dash(push[strings.ToLower(s.Address)]), s.Tier, s.SmartMoneyScore, s.BotScore, leaderboardWhaleScore(s),
				s.Stats.LargeTrades, s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.AvgTradeNotional,
				s.Stats.CopyROI, s.Stats.CopyClosedTrades, sourceSummary(s.Sources), riskSummary(s.RiskFlags))
			continue
		}
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %d | $%.0f | %.1f%% | %d | %s | %s |\n",
			s.Address, s.Tier, s.SmartMoneyScore, s.BotScore, leaderboardWhaleScore(s),
			s.Stats.LargeTrades, s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.AvgTradeNotional,
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	fmt.Fprintf(b, "\n")
}

func hasLeaderboardSource(sources map[string]int) bool {
	for src, n := range sources {
		if n > 0 && strings.HasPrefix(src, "leaderboard_") {
			return true
		}
	}
	return false
}

func hasRecentProfitLeaderboardSource(sources map[string]int) bool {
	for _, src := range []string{"leaderboard_profit_7d", "leaderboard_profit_30d"} {
		if sources[src] > 0 {
			return true
		}
	}
	return false
}

func qualifiesLeaderboardWhale(s walletdiscover.WalletScore, minSmart, maxBot float64, minLarge int, minAvgNotional float64, minTargetTrades, minTargetLarge int) bool {
	if !hasRecentProfitLeaderboardSource(s.Sources) {
		return false
	}
	if !hasTargetCategoryActivity(s, minTargetTrades, minTargetLarge) {
		return false
	}
	if s.SmartMoneyScore < minSmart || s.BotScore >= maxBot {
		return false
	}
	if s.Stats.LargeTrades < minLarge || s.Stats.AvgTradeNotional < minAvgNotional {
		return false
	}
	if hasRisk(s.RiskFlags, "bot_like_flow", "fixed_amount", "fixed_price", "negative_copy_sim") {
		return false
	}
	if s.Stats.CopyClosedTrades >= 3 && (s.Stats.CopyROI < 0 || s.Stats.CopyPnL < 0) {
		return false
	}
	return true
}

func qualifiesStrictLeaderboardPush(s walletdiscover.WalletScore, minTier string, minSmart, maxBot float64, minLarge int, minAvgNotional float64, minTargetTrades, minTargetLarge int) bool {
	if !qualifiesLeaderboardWhale(s, minSmart, maxBot, minLarge, minAvgNotional, minTargetTrades, minTargetLarge) {
		return false
	}
	if !tierAtLeast(s.Tier, minTier) {
		return false
	}
	if hasRisk(s.RiskFlags, "burst_trading", "opposite_side_same_market", "extreme_price_heavy", "open_copy_exposure") {
		return false
	}
	return true
}

func isBotLike(s walletdiscover.WalletScore) bool {
	return s.Tier == "BOT" || s.BotScore >= 45 || hasRisk(s.RiskFlags, "bot_like_flow", "fixed_amount", "fixed_price")
}

func highValueWatch(s walletdiscover.WalletScore, minTargetTrades, minTargetLarge int) bool {
	return hasRecentProfitLeaderboardSource(s.Sources) && hasTargetCategoryActivity(s, minTargetTrades, minTargetLarge) && s.SmartMoneyScore >= 70 && s.Stats.LargeTrades >= 20 && s.Stats.AvgTradeNotional >= 500
}

func qualifiesLeaderboardWhaleWatch(s walletdiscover.WalletScore, minSmart, maxBot float64, minLarge int, minAvgNotional float64, minTargetTrades, minTargetLarge int) bool {
	if !hasLeaderboardSource(s.Sources) {
		return false
	}
	if !hasTargetCategoryActivity(s, minTargetTrades, minTargetLarge) {
		return false
	}
	if s.SmartMoneyScore < minSmart || s.BotScore >= maxBot {
		return false
	}
	if s.Stats.LargeTrades < minLarge || s.Stats.AvgTradeNotional < minAvgNotional {
		return false
	}
	if hasRisk(s.RiskFlags, "bot_like_flow", "fixed_amount", "fixed_price", "negative_copy_sim", "opposite_side_same_market") {
		return false
	}
	return true
}

func hasTargetCategoryActivity(s walletdiscover.WalletScore, minTargetTrades, minTargetLarge int) bool {
	if minTargetTrades > 0 && s.Stats.TargetTrades < minTargetTrades {
		return false
	}
	if minTargetLarge > 0 && s.Stats.TargetLargeTrades < minTargetLarge {
		return false
	}
	return true
}

func leaderboardWhaleScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return s.SmartMoneyScore*0.7 -
		s.BotScore*1.8 +
		float64(st.LargeTrades)*0.06 +
		float64(st.TargetLargeTrades)*0.12 +
		st.TargetTradeRatio*40 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*10 +
		copyBonus(st) +
		sourceBonus(s.Sources)
}

func watchScore(s walletdiscover.WalletScore) float64 {
	return leaderboardWhaleScore(s) - math.Abs(math.Min(s.Stats.CopyROI, 0))*2
}

func botSortScore(s walletdiscover.WalletScore) float64 {
	return s.BotScore*2 + float64(s.Stats.LargeTrades)*0.03 + math.Log1p(math.Max(s.Stats.AvgTradeNotional, 0))*5
}

func copyBonus(st walletdiscover.WalletStats) float64 {
	if st.CopyClosedTrades < 3 {
		return 0
	}
	return st.CopyROI*math.Log1p(float64(st.CopyClosedTrades))*0.5 + math.Log1p(math.Max(st.CopyPnL, 0))*5
}

func sourceBonus(sources map[string]int) float64 {
	var out float64
	for src, n := range sources {
		if n <= 0 {
			continue
		}
		switch {
		case strings.Contains(src, "leaderboard_profit_7d"):
			out += 18
		case strings.Contains(src, "leaderboard_profit_30d"):
			out += 12
		case strings.Contains(src, "leaderboard_profit_all"):
			out += 8
		case strings.Contains(src, "leaderboard_volume_7d"):
			out += 8
		case strings.Contains(src, "leaderboard_volume_30d"):
			out += 5
		}
	}
	return out
}

func hasRisk(flags []string, names ...string) bool {
	set := map[string]struct{}{}
	for _, n := range names {
		set[n] = struct{}{}
	}
	for _, f := range flags {
		if _, ok := set[f]; ok {
			return true
		}
	}
	return false
}

func sortScores(rows []walletdiscover.WalletScore, score func(walletdiscover.WalletScore) float64) {
	sort.Slice(rows, func(i, j int) bool {
		if score(rows[i]) != score(rows[j]) {
			return score(rows[i]) > score(rows[j])
		}
		return rows[i].Address < rows[j].Address
	})
}

func countLeaderboard(scores []walletdiscover.WalletScore) int {
	n := 0
	for _, s := range scores {
		if hasLeaderboardSource(s.Sources) {
			n++
		}
	}
	return n
}

func countLeaderboardPushed(scores []walletdiscover.WalletScore, push map[string]string) int {
	n := 0
	for _, s := range scores {
		if hasLeaderboardSource(s.Sources) {
			if _, ok := push[strings.ToLower(s.Address)]; ok {
				n++
			}
		}
	}
	return n
}

func sourceSummary(sources map[string]int) string {
	var keys []string
	for k, n := range sources {
		if n > 0 && strings.HasPrefix(k, "leaderboard_") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if len(keys) > 5 {
		keys = keys[:5]
	}
	if len(keys) == 0 {
		return "-"
	}
	return strings.Join(keys, ",")
}

func riskSummary(flags []string) string {
	if len(flags) == 0 {
		return "-"
	}
	out := append([]string{}, flags...)
	sort.Strings(out)
	if len(out) > 5 {
		out = out[:5]
	}
	return strings.Join(out, ",")
}

func tierAtLeast(got, minTier string) bool {
	rank := map[string]int{"D": 1, "C": 2, "B": 3, "A": 4}
	minTier = strings.ToUpper(strings.TrimSpace(minTier))
	if minTier == "" {
		return true
	}
	return rank[strings.ToUpper(strings.TrimSpace(got))] >= rank[minTier]
}

func writeWalletFile(path string, rows []walletdiscover.WalletScore, list string) error {
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# generated by leaderboard-whales at %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&b, "# list=%s rows=%d\n", list, len(rows))
	for _, s := range rows {
		fmt.Fprintf(&b, "%s # list=%s tier=%s smart=%.1f bot=%.1f whaleScore=%.1f large=%d targetT=%d targetLarge=%d avgNotional=$%.0f targetCopyROI=%.1f targetCopyT=%d copyROI=%.1f copyT=%d sources=%s\n",
			strings.ToLower(s.Address), list, dash(s.Tier), s.SmartMoneyScore, s.BotScore, leaderboardWhaleScore(s),
			s.Stats.LargeTrades, s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.AvgTradeNotional, s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed, s.Stats.CopyROI, s.Stats.CopyClosedTrades, sourceSummary(s.Sources))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func limitScores(rows []walletdiscover.WalletScore, n int) []walletdiscover.WalletScore {
	if n <= 0 || len(rows) <= n {
		return rows
	}
	return rows[:n]
}

func mergeScores(groups ...[]walletdiscover.WalletScore) []walletdiscover.WalletScore {
	seen := map[string]struct{}{}
	var out []walletdiscover.WalletScore
	for _, group := range groups {
		for _, row := range group {
			addr := strings.ToLower(row.Address)
			if addr == "" {
				continue
			}
			if _, ok := seen[addr]; ok {
				continue
			}
			seen[addr] = struct{}{}
			out = append(out, row)
		}
	}
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
