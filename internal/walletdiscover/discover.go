package walletdiscover

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Run(ctx context.Context, cfg Config) (*Result, error) {
	cfg = normalizeConfig(cfg)
	client := NewClient(cfg)

	markets, err := discoverMarkets(ctx, client, cfg)
	if err != nil {
		return nil, err
	}
	allowed := map[string]struct{}{}
	for _, m := range markets {
		allowed[m.ConditionID] = struct{}{}
	}

	candidates := map[string]*Candidate{}
	addExistingWallets(candidates, cfg)
	addSourceWallets(candidates, cfg.SportsTapeWallets, "sports_tape")
	addSourceWallets(candidates, cfg.RetainWallets, "retain")
	if err := addLeaderboardCandidates(ctx, client, candidates, cfg); err != nil {
		slog.Warn("wallet_discover.leaderboard_partial", "err", err)
	}
	if err := addHolderCandidates(ctx, client, markets, candidates, cfg); err != nil {
		slog.Warn("wallet_discover.holders_partial", "err", err)
	}
	if err := addTradeCandidates(ctx, client, candidates, allowed, cfg); err != nil {
		slog.Warn("wallet_discover.trades_partial", "err", err)
	}

	rankedCandidates := rankCandidates(candidates, cfg.MaxCandidates)
	previous, err := LoadPreviousScores(cfg.OutputDir)
	if err != nil {
		slog.Warn("wallet_discover.previous_scores_fail", "err", err)
		previous = map[string]WalletScore{}
	}
	scores := scoreCandidates(ctx, client, rankedCandidates, cfg, previous)
	SortScores(scores)
	result := &Result{
		Markets:     markets,
		Candidates:  rankedCandidates,
		Scores:      scores,
		HTTP:        client.Stats(),
		DataQuality: summarizeDataQuality(scores),
	}
	if err := SaveResult(result, cfg); err != nil {
		return nil, err
	}
	return result, nil
}

func normalizeConfig(cfg Config) Config {
	def := DefaultConfig()
	if cfg.GammaBase == "" {
		cfg.GammaBase = def.GammaBase
	}
	if cfg.DataBase == "" {
		cfg.DataBase = def.DataBase
	}
	if cfg.LeaderboardBase == "" {
		cfg.LeaderboardBase = def.LeaderboardBase
	}
	if cfg.HTTPMaxAttempts <= 0 {
		cfg.HTTPMaxAttempts = def.HTTPMaxAttempts
	}
	if cfg.HTTPRetryBase <= 0 {
		cfg.HTTPRetryBase = def.HTTPRetryBase
	}
	if cfg.HTTPRetryMax <= 0 {
		cfg.HTTPRetryMax = def.HTTPRetryMax
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = def.HTTPTimeout
	}
	if cfg.MarketsLimit <= 0 {
		cfg.MarketsLimit = def.MarketsLimit
	}
	if cfg.TradesPages <= 0 {
		cfg.TradesPages = def.TradesPages
	}
	if cfg.TradesLimit <= 0 {
		cfg.TradesLimit = def.TradesLimit
	}
	if cfg.HoldersLimit <= 0 {
		cfg.HoldersLimit = def.HoldersLimit
	}
	if cfg.ActivityPages <= 0 {
		cfg.ActivityPages = def.ActivityPages
	}
	if cfg.ActivityLimit <= 0 {
		cfg.ActivityLimit = def.ActivityLimit
	}
	if cfg.ClosedLimit <= 0 {
		cfg.ClosedLimit = def.ClosedLimit
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = def.Concurrency
	}
	if cfg.MaxCandidates <= 0 {
		cfg.MaxCandidates = def.MaxCandidates
	}
	if cfg.Days <= 0 {
		cfg.Days = def.Days
	}
	if cfg.MinNotionalUSD <= 0 {
		cfg.MinNotionalUSD = def.MinNotionalUSD
	}
	if cfg.MinPrice <= 0 {
		cfg.MinPrice = def.MinPrice
	}
	if cfg.MinHolderShares <= 0 {
		cfg.MinHolderShares = def.MinHolderShares
	}
	if cfg.MaxPrice <= 0 {
		cfg.MaxPrice = def.MaxPrice
	}
	if cfg.CopyStakeUSD <= 0 {
		cfg.CopyStakeUSD = def.CopyStakeUSD
	}
	if cfg.CopySlippageBP < 0 {
		cfg.CopySlippageBP = def.CopySlippageBP
	}
	if cfg.CopyFeeBP < 0 {
		cfg.CopyFeeBP = def.CopyFeeBP
	}
	if cfg.LeaderboardWindows == "" {
		cfg.LeaderboardWindows = def.LeaderboardWindows
	}
	if cfg.LeaderboardKinds == "" {
		cfg.LeaderboardKinds = def.LeaderboardKinds
	}
	if cfg.TargetCategories == "" {
		cfg.TargetCategories = def.TargetCategories
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = def.OutputDir
	}
	if cfg.ReportPath == "" {
		cfg.ReportPath = def.ReportPath
	}
	if cfg.GeneratedTier == "" {
		cfg.GeneratedTier = def.GeneratedTier
	}
	if cfg.GeneratedWalletsPath == "" {
		cfg.GeneratedWalletsPath = def.GeneratedWalletsPath
	}
	if cfg.AutoWalletsPath == "" {
		cfg.AutoWalletsPath = def.AutoWalletsPath
	}
	if cfg.PromptWalletsPath == "" {
		cfg.PromptWalletsPath = def.PromptWalletsPath
	}
	if cfg.PositiveWalletsPath == "" {
		cfg.PositiveWalletsPath = def.PositiveWalletsPath
	}
	return cfg
}

func discoverMarkets(ctx context.Context, client *Client, cfg Config) ([]Market, error) {
	var out []Market
	offset := 0
	now := time.Now()
	for len(out) < cfg.MarketsLimit {
		limit := 100
		if cfg.MarketsLimit-len(out) < limit {
			limit = cfg.MarketsLimit - len(out)
		}
		batch, err := client.ListMarkets(ctx, limit, offset)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		for _, m := range batch {
			if GoodDiscoveryMarket(m, now) && TargetCategoryAllowed(MarketTargetCategory(m), cfg.TargetCategories) {
				out = append(out, m)
			}
		}
		offset += len(batch)
		if len(batch) < limit {
			break
		}
	}
	slog.Info("wallet_discover.markets", "selected", len(out))
	return out, nil
}

func addExistingWallets(candidates map[string]*Candidate, cfg Config) {
	addSourceWallets(candidates, cfg.ExistingWallets, "existing")
}

func addSourceWallets(candidates map[string]*Candidate, path, source string) {
	if path == "" || source == "" {
		return
	}
	rows, err := loadSourceWalletRows(path)
	if err != nil {
		slog.Warn("wallet_discover.source_load_fail", "source", source, "file", path, "err", err)
		return
	}
	for _, row := range rows {
		c := ensureCandidate(candidates, row.Address)
		c.Sources[source]++
		if row.Buys > 0 {
			c.ObservedTrades += row.Buys
		}
		if row.BuyNotional > 0 {
			c.ObservedNotional += row.BuyNotional
		} else if row.MaxBuy > 0 {
			c.ObservedNotional += row.MaxBuy
		}
	}
}

type sourceWalletRow struct {
	Address     string
	Buys        int
	BuyNotional float64
	MaxBuy      float64
}

func loadSourceWalletRows(path string) ([]sourceWalletRow, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []sourceWalletRow
	seen := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		body, comment, _ := strings.Cut(sc.Text(), "#")
		fields := strings.Fields(strings.TrimSpace(body))
		if len(fields) == 0 {
			continue
		}
		addr := normalizeAddress(fields[0])
		if addr == "" {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		meta := parseSourceCommentFields(comment)
		out = append(out, sourceWalletRow{
			Address:     addr,
			Buys:        parseSourceInt(meta["buys"]),
			BuyNotional: parseSourceNumber(meta["buyNotional"]),
			MaxBuy:      parseSourceNumber(meta["maxBuy"]),
		})
	}
	return out, sc.Err()
}

func parseSourceCommentFields(comment string) map[string]string {
	out := map[string]string{}
	for _, field := range strings.Fields(comment) {
		k, v, ok := strings.Cut(field, "=")
		if ok {
			out[k] = v
		}
	}
	return out
}

func parseSourceInt(raw string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(raw))
	return v
}

func parseSourceNumber(raw string) float64 {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "$")
	raw = strings.TrimSuffix(raw, "%")
	raw = strings.ReplaceAll(raw, ",", "")
	if raw == "" || raw == "-" {
		return 0
	}
	v, _ := strconv.ParseFloat(raw, 64)
	return v
}

func addLeaderboardCandidates(ctx context.Context, client *Client, candidates map[string]*Candidate, cfg Config) error {
	if cfg.LeaderboardLimit <= 0 {
		return nil
	}
	windows := csvFields(cfg.LeaderboardWindows)
	kinds := csvFields(cfg.LeaderboardKinds)
	if len(windows) == 0 {
		windows = []string{"all"}
	}
	if len(kinds) == 0 {
		kinds = []string{"profit"}
	}
	total := 0
	for _, kind := range kinds {
		for _, window := range windows {
			fetched := 0
			for fetched < cfg.LeaderboardLimit {
				pageLimit := 50
				if remaining := cfg.LeaderboardLimit - fetched; remaining < pageLimit {
					pageLimit = remaining
				}
				entries, err := client.LeaderboardPage(ctx, kind, window, pageLimit, fetched)
				if err != nil {
					slog.Warn("wallet_discover.leaderboard_fail", "kind", kind, "window", window, "offset", fetched, "err", err)
					break
				}
				if len(entries) == 0 {
					break
				}
				for _, e := range entries {
					addr := normalizeAddress(e.ProxyWallet)
					if addr == "" {
						continue
					}
					c := ensureCandidate(candidates, addr)
					c.Sources["leaderboard_"+kind+"_"+window]++
					if e.Amount > 0 {
						c.ObservedNotional += e.Amount
					}
					if e.Name != "" {
						c.Names[e.Name]++
					}
					if e.Pseudonym != "" {
						c.Names[e.Pseudonym]++
					}
				}
				fetched += len(entries)
				if len(entries) < pageLimit {
					break
				}
			}
			total += fetched
			slog.Info("wallet_discover.leaderboard",
				"kind", kind,
				"window", window,
				"entries", fetched,
				"candidates", len(candidates),
			)
		}
	}
	if total == 0 {
		return fmt.Errorf("no leaderboard entries fetched")
	}
	return nil
}

func addHolderCandidates(ctx context.Context, client *Client, markets []Market, candidates map[string]*Candidate, cfg Config) error {
	for i, m := range markets {
		if err := ctx.Err(); err != nil {
			return err
		}
		resp, err := client.Holders(ctx, m.ConditionID, cfg.HoldersLimit)
		if err != nil {
			slog.Warn("wallet_discover.holders_fail", "market", m.ConditionID, "err", err)
			continue
		}
		for _, token := range resp {
			for _, h := range token.Holders {
				addr := normalizeAddress(h.ProxyWallet)
				if addr == "" || h.Amount < cfg.MinHolderShares {
					continue
				}
				c := ensureCandidate(candidates, addr)
				c.Sources["holder"]++
				c.ObservedHolders++
				if h.Amount > c.MaxHolderShares {
					c.MaxHolderShares = h.Amount
				}
				if h.Name != "" {
					c.Names[h.Name]++
				}
				if h.Pseudonym != "" {
					c.Names[h.Pseudonym]++
				}
				c.Markets[m.ConditionID] = h.Amount
			}
		}
		if (i+1)%25 == 0 {
			slog.Info("wallet_discover.holders_progress", "markets", i+1, "candidates", len(candidates))
		}
	}
	return nil
}

func addTradeCandidates(ctx context.Context, client *Client, candidates map[string]*Candidate, allowed map[string]struct{}, cfg Config) error {
	for page := 0; page < cfg.TradesPages; page++ {
		trades, err := client.RecentTrades(ctx, cfg.TradesLimit, page*cfg.TradesLimit)
		if err != nil {
			if IsMaxHistoricalOffsetError(err) {
				slog.Info("wallet_discover.trades_offset_limit", "page", page+1, "offset", page*cfg.TradesLimit)
				return nil
			}
			return err
		}
		for _, tr := range trades {
			if !QualifyingTrade(tr, cfg, allowed) {
				continue
			}
			addr := normalizeAddress(tr.ProxyWallet)
			c := ensureCandidate(candidates, addr)
			c.Sources["recent_trade"]++
			c.ObservedTrades++
			c.ObservedNotional += tr.NotionalUSD()
			if tr.ConditionID != "" {
				c.Markets[tr.ConditionID] += tr.NotionalUSD()
			}
		}
		slog.Info("wallet_discover.trades_page", "page", page+1, "candidates", len(candidates))
	}
	return nil
}

func scoreCandidates(ctx context.Context, client *Client, candidates []*Candidate, cfg Config, previous map[string]WalletScore) []WalletScore {
	type item struct {
		score WalletScore
	}
	jobs := make(chan *Candidate)
	results := make(chan item)
	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cand := range jobs {
				var issues []string
				usedCachedActivity := false
				trades, err := pullActivity(ctx, client, cand.Address, cfg)
				if err != nil {
					slog.Warn("wallet_discover.activity_fail", "wallet", short(cand.Address), "err", err)
					issues = append(issues, "activity_api_unavailable")
					cached, cacheErr := LoadCachedActivity(cfg.OutputDir, cand.Address)
					if cacheErr != nil {
						slog.Warn("wallet_discover.activity_cache_read_fail", "wallet", short(cand.Address), "err", cacheErr)
						issues = append(issues, "activity_cache_unavailable")
					} else if len(cached) > 0 {
						trades = cached
						usedCachedActivity = true
					}
				}
				activityIncomplete := err != nil && !usedCachedActivity
				closed, closedErr := client.ClosedPositions(ctx, cand.Address, cfg.ClosedLimit)
				if closedErr != nil {
					slog.Warn("wallet_discover.closed_fail", "wallet", short(cand.Address), "err", closedErr)
					issues = append(issues, "closed_positions_api_unavailable")
				}
				score := ScoreWallet(cand.Address, cand, trades, closed, cfg)
				switch {
				case activityIncomplete || closedErr != nil:
					if old, ok := previous[cand.Address]; ok {
						score = preservePreviousScore(old, cand, issues)
					} else {
						score.DataStatus = "incomplete"
						score.DataIssues = issues
						score.Tier = "D"
						score.FollowAction = "reject"
						score.Reason = "incomplete API data"
						score.RiskFlags = appendUnique(score.RiskFlags, "incomplete_data")
					}
				case usedCachedActivity:
					score.DataStatus = "cached_activity"
					score.DataIssues = issues
					score.RiskFlags = appendUnique(score.RiskFlags, "cached_activity")
				default:
					score.DataStatus = "complete"
				}
				results <- item{score: score}
			}
		}()
	}
	go func() {
		for _, cand := range candidates {
			jobs <- cand
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()
	var scores []WalletScore
	for r := range results {
		scores = append(scores, r.score)
	}
	return scores
}

func preservePreviousScore(previous WalletScore, cand *Candidate, issues []string) WalletScore {
	previous.Address = cand.Address
	previous.Sources = copySources(cand.Sources)
	previous.DataStatus = "preserved_previous"
	previous.DataIssues = append([]string(nil), issues...)
	previous.RiskFlags = appendUnique(append([]string(nil), previous.RiskFlags...), "stale_data_preserved")
	previous.Reason = strings.TrimSpace(previous.Reason)
	if previous.Reason == "" {
		previous.Reason = "previous score preserved"
	} else if !strings.Contains(previous.Reason, "previous score preserved") {
		previous.Reason += "; previous score preserved"
	}
	return previous
}

func copySources(sources map[string]int) map[string]int {
	out := make(map[string]int, len(sources))
	for key, value := range sources {
		out[key] = value
	}
	return out
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func summarizeDataQuality(scores []WalletScore) DataQuality {
	var quality DataQuality
	for _, score := range scores {
		switch score.DataStatus {
		case "cached_activity":
			quality.CachedActivity++
		case "preserved_previous":
			quality.PreservedPrevious++
		case "incomplete":
			quality.Incomplete++
		default:
			quality.Complete++
		}
	}
	return quality
}

func pullActivity(ctx context.Context, client *Client, addr string, cfg Config) ([]Trade, error) {
	cutoff := time.Now().Add(-time.Duration(cfg.Days) * 24 * time.Hour).Unix()
	var all []Trade
	for page := 0; page < cfg.ActivityPages; page++ {
		batch, err := client.Activity(ctx, addr, cfg.ActivityLimit, page*cfg.ActivityLimit)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		var oldest int64
		for _, tr := range batch {
			if oldest == 0 || tr.Timestamp < oldest {
				oldest = tr.Timestamp
			}
		}
		if oldest > 0 && oldest < cutoff {
			break
		}
	}
	var filtered []Trade
	for _, tr := range all {
		if tr.Timestamp >= cutoff {
			filtered = append(filtered, tr)
		}
	}
	if err := SaveActivity(cfg.OutputDir, addr, filtered); err != nil {
		slog.Warn("wallet_discover.activity_cache_fail", "wallet", short(addr), "err", err)
	}
	return filtered, nil
}

func ensureCandidate(candidates map[string]*Candidate, addr string) *Candidate {
	addr = normalizeAddress(addr)
	c := candidates[addr]
	if c == nil {
		c = &Candidate{
			Address: addr,
			Sources: map[string]int{},
			Names:   map[string]int{},
			Markets: map[string]float64{},
		}
		candidates[addr] = c
	}
	return c
}

func rankCandidates(m map[string]*Candidate, max int) []*Candidate {
	out := make([]*Candidate, 0, len(m))
	for _, c := range m {
		if c.Address != "" {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		si, sj := candidatePriority(out[i]), candidatePriority(out[j])
		if si != sj {
			return si > sj
		}
		return out[i].Address < out[j].Address
	})
	if max > 0 && len(out) > max {
		out = out[:max]
	}
	return out
}

func candidatePriority(c *Candidate) float64 {
	return float64(c.Sources["existing"])*100000 +
		float64(c.Sources["sports_tape"])*80000 +
		leaderboardSourcePriority(c.Sources) +
		float64(c.Sources["holder"])*1000 +
		c.ObservedNotional +
		float64(c.ObservedTrades)*100 +
		c.MaxHolderShares
}

func leaderboardSourcePriority(sources map[string]int) float64 {
	return float64(sources["leaderboard_profit_7d"])*70000 +
		float64(sources["leaderboard_profit_30d"])*60000 +
		float64(sources["leaderboard_profit_all"])*35000 +
		float64(sources["leaderboard_volume_7d"])*15000 +
		float64(sources["leaderboard_volume_30d"])*10000 +
		float64(sources["leaderboard_volume_all"])*5000
}

func csvFields(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func short(addr string) string {
	if len(addr) <= 10 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func RenderReport(result *Result, cfg Config) string {
	var b strings.Builder
	counts := map[string]int{}
	for _, s := range result.Scores {
		counts[s.Tier]++
	}
	fmt.Fprintf(&b, "# Wallet Discovery Report\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Markets scanned: %d\n", len(result.Markets))
	fmt.Fprintf(&b, "- Candidates scored: %d\n", len(result.Scores))
	fmt.Fprintf(&b, "- Tiers: A=%d B=%d C=%d D=%d BOT=%d\n", counts["A"], counts["B"], counts["C"], counts["D"], counts["BOT"])
	fmt.Fprintf(&b, "- Data quality: complete=%d cached_activity=%d preserved_previous=%d incomplete=%d\n",
		result.DataQuality.Complete, result.DataQuality.CachedActivity, result.DataQuality.PreservedPrevious, result.DataQuality.Incomplete)
	fmt.Fprintf(&b, "- HTTP resilience: retries=%d rate_limits=%d terminal_failures=%d\n", result.HTTP.Retries, result.HTTP.RateLimits, result.HTTP.Failures)
	if cfg.TargetCategories != "" {
		fmt.Fprintf(&b, "- Target categories: %s\n", cfg.TargetCategories)
	}
	fmt.Fprintf(&b, "- Filters: min_notional=$%.0f price=%.2f-%.2f days=%d copy_stake=$%.0f slippage=%.0fbp fee=%.0fbp\n\n",
		cfg.MinNotionalUSD, cfg.MinPrice, cfg.MaxPrice, cfg.Days, cfg.CopyStakeUSD, cfg.CopySlippageBP, cfg.CopyFeeBP)
	if cfg.LeaderboardLimit > 0 {
		fmt.Fprintf(&b, "- Leaderboard: kinds=%s windows=%s limit=%d\n\n", cfg.LeaderboardKinds, cfg.LeaderboardWindows, cfg.LeaderboardLimit)
	}
	fmt.Fprintf(&b, "## Top Wallets\n\n")
	fmt.Fprintf(&b, "| Rank | Wallet | Tier | Action | Data | Smart | Bot | Edge | Valid | Large | TargetT | TargetLarge | Target%% | TargetCopyROI | TargetCopyT | ROI | CopyROI | CopyPnL | CopyT | Focus | Reason |\n")
	fmt.Fprintf(&b, "|---:|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for i, s := range result.Scores {
		if i >= 50 {
			break
		}
		focus := s.Stats.TopCategory
		if focus == "" {
			focus = "-"
		} else {
			focus = fmt.Sprintf("%s %.0f%%", focus, s.Stats.TopCategoryRatio*100)
		}
		fmt.Fprintf(&b, "| %d | `%s` | %s | %s | %s | %.2f | %.2f | %.2f | %d | %d | %d | %d | %.0f%% | %.1f%% | %d | %.1f%% | %.1f%% | $%+.2f | %d | %s | %s |\n",
			i+1, short(s.Address), s.Tier, s.FollowAction, s.DataStatus, s.SmartMoneyScore, s.BotScore, s.Edge,
			s.Stats.ValidTrades, s.Stats.LargeTrades, s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed, s.Stats.ClosedROI, s.Stats.CopyROI, s.Stats.CopyPnL, s.Stats.CopyClosedTrades, focus, s.Reason)
	}
	return b.String()
}
