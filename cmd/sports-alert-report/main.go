package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

type alertEvent struct {
	Key           string    `json:"key"`
	SentAt        time.Time `json:"sent_at"`
	TradeTime     time.Time `json:"trade_time"`
	DryRun        bool      `json:"dry_run,omitempty"`
	Reconciled    bool      `json:"reconciled,omitempty"`
	Mode          string    `json:"mode"`
	Wallet        string    `json:"wallet"`
	KnownList     string    `json:"known_list,omitempty"`
	Tier          string    `json:"tier,omitempty"`
	Bot           float64   `json:"bot,omitempty"`
	Category      string    `json:"category"`
	Notional      float64   `json:"notional"`
	Price         float64   `json:"price"`
	Outcome       string    `json:"outcome"`
	Market        string    `json:"market"`
	Slug          string    `json:"slug,omitempty"`
	Asset         string    `json:"asset"`
	Transaction   string    `json:"transaction,omitempty"`
	TargetCopyROI float64   `json:"target_copy_roi,omitempty"`
	TargetCopyT   int       `json:"target_copy_t,omitempty"`
}

type alertResult struct {
	Event    alertEvent
	StakeUSD float64
	Units    float64
	Mid      float64
	Marked   bool
	PnLUSD   float64
	ReturnPC float64
}

type resultSummary struct {
	Signals   int
	Marked    int
	Unmarked  int
	Wins      int
	StakeUSD  float64
	PnLUSD    float64
	ReturnPC  float64
	AvgDeltaP float64
}

type policyDecisionReport struct {
	GeneratedAt                   time.Time        `json:"generated_at"`
	Log                           string           `json:"log"`
	ExtraLog                      string           `json:"extra_log,omitempty"`
	StakeUSD                      float64          `json:"stake_usd"`
	CurrentPolicy                 decisionBucket   `json:"current_policy"`
	CurrentPolicyPositionCapped   decisionBucket   `json:"current_policy_position_capped"`
	EffectivePolicy               decisionBucket   `json:"effective_policy"`
	EffectivePolicyPositionCapped decisionBucket   `json:"effective_policy_position_capped"`
	EffectiveModes                []string         `json:"effective_modes"`
	Modes                         []decisionBucket `json:"modes"`
	RecommendedModes              []string         `json:"recommended_modes"`
	PositiveModes                 []string         `json:"positive_modes"`
	CutModes                      []string         `json:"cut_modes"`
	ProbationModes                []string         `json:"probation_modes"`
}

type decisionBucket struct {
	Key        string  `json:"key"`
	Action     string  `json:"action"`
	Reason     string  `json:"reason"`
	Alerts     int     `json:"alerts"`
	Marked     int     `json:"marked"`
	Unmarked   int     `json:"unmarked"`
	Wins       int     `json:"wins"`
	WinRate    float64 `json:"win_rate"`
	ROI        float64 `json:"roi"`
	PnLUSD     float64 `json:"pnl_usd"`
	AvgDeltaPP float64 `json:"avg_delta_pp"`
}

type midpointCache struct {
	Assets map[string]cachedMidpoint `json:"assets"`
}

type cachedMidpoint struct {
	Mid       float64   `json:"mid"`
	FetchedAt time.Time `json:"fetched_at"`
}

func main() {
	logPath := flag.String("log", "db/strategy_iteration/sports_tape_alerts.jsonl", "sports tape alert audit JSONL")
	extraLogPath := flag.String("extra_log", "", "comma-separated extra alert JSONL logs merged into evaluation")
	reportPath := flag.String("report", "reports/sports_alert_performance.md", "markdown report path")
	markCachePath := flag.String("mark_cache", "db/strategy_iteration/sports_alert_midpoints.json", "JSON cache for last known token midpoints; empty disables cache")
	decisionPath := flag.String("decision_json", "", "optional JSON policy decision output for automation")
	gammaBase := flag.String("gamma_base", "https://gamma-api.polymarket.com", "Gamma API base URL used for settled-market mark fallback")
	stakeUSD := flag.Float64("stake", 10, "fixed paper stake per alert")
	currentPolicyModesRaw := flag.String("current_policy_modes", "FLOW-SCOUT,EDGE-HOT,FOLLOW-READY,CANDIDATE,PROBATION,CONSENSUS", "comma-separated modes counted as the currently enabled alert policy")
	currentExcludeWalletsPath := flag.String("current_exclude_wallets", "", "comma-separated wallet files excluded from current policy metrics")
	timeout := flag.Duration("timeout", 20*time.Second, "overall midpoint fetch timeout")
	flag.Parse()

	events, err := loadAlertEvents(*logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-alert-report: load alerts: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(*extraLogPath) != "" {
		extraEvents, err := loadExtraAlertEvents(*extraLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sports-alert-report: load extra alerts: %v\n", err)
			os.Exit(1)
		}
		events = mergeAlertEvents(events, extraEvents)
	}
	currentExcludeWallets, err := loadWalletSet(*currentExcludeWalletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-alert-report: load current exclude wallets: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	client := &http.Client{Timeout: minDuration(*timeout, 10*time.Second)}
	cache, err := loadMidpointCache(*markCachePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-alert-report: load mark cache: %v\n", err)
		os.Exit(1)
	}
	assetSlugs := alertAssetSlugs(events)
	cacheDirty := false
	results := evaluateAlerts(ctx, events, *stakeUSD, func(ctx context.Context, asset string) (float64, error) {
		mid, err := fetchMidpoint(ctx, client, asset)
		if err == nil && mid > 0 {
			if cache != nil {
				cache.Set(asset, mid, time.Now())
				cacheDirty = true
			}
			return mid, nil
		}
		mid, settlementErr := fetchSettledTokenPrice(ctx, client, *gammaBase, asset)
		if settlementErr != nil {
			mid, settlementErr = fetchSettledTokenPriceByEventSlug(ctx, client, *gammaBase, assetSlugs[asset], asset)
		}
		if settlementErr == nil && mid >= 0 {
			if cache != nil {
				cache.Set(asset, mid, time.Now())
				cacheDirty = true
			}
			return mid, nil
		}
		if cache != nil {
			if cached, ok := cache.Get(asset); ok {
				return cached, nil
			}
		}
		return 0, err
	})
	if cacheDirty {
		if err := saveMidpointCache(*markCachePath, cache); err != nil {
			fmt.Fprintf(os.Stderr, "sports-alert-report: save mark cache: %v\n", err)
			os.Exit(1)
		}
	}
	currentPolicyModes := parseModeSet(*currentPolicyModesRaw)
	if err := writeReport(*reportPath, *logPath, *extraLogPath, results, *stakeUSD, currentPolicyModes, currentExcludeWallets, *currentExcludeWalletsPath); err != nil {
		fmt.Fprintf(os.Stderr, "sports-alert-report: write report: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(*decisionPath) != "" {
		if err := writeDecisionJSON(*decisionPath, *logPath, *extraLogPath, results, *stakeUSD, currentPolicyModes, currentExcludeWallets); err != nil {
			fmt.Fprintf(os.Stderr, "sports-alert-report: write decision json: %v\n", err)
			os.Exit(1)
		}
	}
	sum := summarize(results)
	fmt.Printf("sports-alert-report done: alerts=%d marked=%d roi=%.1f%% report=%s\n", sum.Signals, sum.Marked, sum.ReturnPC, *reportPath)
}

func loadAlertEvents(path string) ([]alertEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []alertEvent
	seen := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var ev alertEvent
		if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
			continue
		}
		ev.Key = strings.ToLower(strings.TrimSpace(ev.Key))
		ev.Wallet = strings.ToLower(strings.TrimSpace(ev.Wallet))
		ev.Asset = strings.TrimSpace(ev.Asset)
		if ev.Key == "" || ev.Wallet == "" || ev.Asset == "" || ev.Price <= 0 {
			continue
		}
		if _, ok := seen[ev.Key]; ok {
			continue
		}
		seen[ev.Key] = struct{}{}
		out = append(out, ev)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SentAt.Before(out[j].SentAt) })
	return out, nil
}

func loadExtraAlertEvents(raw string) ([]alertEvent, error) {
	var out []alertEvent
	for _, path := range strings.Split(raw, ",") {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		events, err := loadAlertEvents(path)
		if err != nil {
			return nil, err
		}
		out = append(out, events...)
	}
	return out, nil
}

func mergeAlertEvents(primary, extra []alertEvent) []alertEvent {
	seen := map[string]struct{}{}
	out := make([]alertEvent, 0, len(primary)+len(extra))
	for _, ev := range append(primary, extra...) {
		key := strings.ToLower(strings.TrimSpace(ev.Key))
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, ev)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SentAt.Before(out[j].SentAt) })
	return out
}

func alertAssetSlugs(events []alertEvent) map[string]string {
	out := map[string]string{}
	for _, ev := range events {
		asset := strings.TrimSpace(ev.Asset)
		slug := strings.TrimSpace(ev.Slug)
		if asset == "" || slug == "" {
			continue
		}
		if _, ok := out[asset]; !ok {
			out[asset] = slug
		}
	}
	return out
}

func loadMidpointCache(path string) (*midpointCache, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &midpointCache{Assets: map[string]cachedMidpoint{}}, nil
		}
		return nil, err
	}
	defer f.Close()

	var c midpointCache
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	if c.Assets == nil {
		c.Assets = map[string]cachedMidpoint{}
	}
	return &c, nil
}

func saveMidpointCache(path string, c *midpointCache) error {
	if strings.TrimSpace(path) == "" || c == nil {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (c *midpointCache) Get(asset string) (float64, bool) {
	if c == nil {
		return 0, false
	}
	row, ok := c.Assets[strings.TrimSpace(asset)]
	if !ok || row.Mid < 0 {
		return 0, false
	}
	return row.Mid, true
}

func (c *midpointCache) Set(asset string, mid float64, fetchedAt time.Time) {
	if c == nil || strings.TrimSpace(asset) == "" || mid < 0 {
		return
	}
	if c.Assets == nil {
		c.Assets = map[string]cachedMidpoint{}
	}
	c.Assets[strings.TrimSpace(asset)] = cachedMidpoint{Mid: mid, FetchedAt: fetchedAt}
}

func loadWalletSet(path string) (map[string]struct{}, error) {
	out := map[string]struct{}{}
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
			for addr := range set {
				out[addr] = struct{}{}
			}
		}
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		body, comment, _ := strings.Cut(sc.Text(), "#")
		line := strings.TrimSpace(body)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(strings.TrimSpace(fields[0]))
		if strings.HasPrefix(addr, "0x") && len(addr) == 42 {
			if excludeWalletRow(comment) {
				out[addr] = struct{}{}
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func excludeWalletRow(comment string) bool {
	fields := parseCommentFields(comment)
	status := strings.ToLower(strings.TrimSpace(fields["status"]))
	if status == "" {
		return true
	}
	return strings.HasPrefix(status, "reject-") || strings.HasPrefix(status, "blocked-")
}

func parseCommentFields(comment string) map[string]string {
	out := map[string]string{}
	for _, field := range strings.Fields(comment) {
		k, v, ok := strings.Cut(field, "=")
		if ok {
			out[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
		}
	}
	return out
}

func evaluateAlerts(ctx context.Context, events []alertEvent, stakeUSD float64, fetch func(context.Context, string) (float64, error)) []alertResult {
	mids := map[string]float64{}
	out := make([]alertResult, 0, len(events))
	for _, ev := range events {
		res := alertResult{Event: ev, StakeUSD: stakeUSD}
		if ev.Price > 0 && stakeUSD > 0 {
			res.Units = stakeUSD / ev.Price
		}
		mid, ok := mids[ev.Asset]
		marked := ok
		if !ok {
			var err error
			mid, err = fetch(ctx, ev.Asset)
			if err == nil && mid >= 0 {
				mids[ev.Asset] = mid
				marked = true
			}
		}
		if marked && res.Units > 0 {
			res.Mid = mid
			res.Marked = true
			res.PnLUSD = res.Units * (mid - ev.Price)
			res.ReturnPC = res.PnLUSD / stakeUSD * 100
		}
		out = append(out, res)
	}
	return out
}

func fetchMidpoint(ctx context.Context, client *http.Client, asset string) (float64, error) {
	reqURL := "https://clob.polymarket.com/midpoint?token_id=" + url.QueryEscape(asset)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("midpoint status %d", resp.StatusCode)
	}
	var body struct {
		Mid string `json:"mid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}
	mid, err := strconv.ParseFloat(body.Mid, 64)
	if err != nil {
		return 0, err
	}
	return mid, nil
}

type gammaMarket struct {
	Closed           bool   `json:"closed"`
	ClobTokenIDsRaw  string `json:"clobTokenIds"`
	OutcomePricesRaw string `json:"outcomePrices"`
}

type gammaEvent struct {
	Markets []gammaMarket `json:"markets"`
}

func fetchSettledTokenPrice(ctx context.Context, client *http.Client, gammaBase, asset string) (float64, error) {
	asset = strings.TrimSpace(asset)
	if asset == "" {
		return 0, fmt.Errorf("empty asset")
	}
	base := strings.TrimRight(strings.TrimSpace(gammaBase), "/")
	if base == "" {
		return 0, fmt.Errorf("empty gamma base")
	}
	reqURL := base + "/markets?clob_token_ids=" + url.QueryEscape(asset)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("gamma status %d", resp.StatusCode)
	}
	var markets []gammaMarket
	if err := json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		return 0, err
	}
	return settledTokenPriceFromMarkets(markets, asset)
}

func fetchSettledTokenPriceByEventSlug(ctx context.Context, client *http.Client, gammaBase, slug, asset string) (float64, error) {
	slug = strings.TrimSpace(slug)
	asset = strings.TrimSpace(asset)
	if slug == "" || asset == "" {
		return 0, fmt.Errorf("empty slug or asset")
	}
	base := strings.TrimRight(strings.TrimSpace(gammaBase), "/")
	if base == "" {
		return 0, fmt.Errorf("empty gamma base")
	}
	reqURL := base + "/events?slug=" + url.QueryEscape(slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("gamma event status %d", resp.StatusCode)
	}
	var events []gammaEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return 0, err
	}
	for _, ev := range events {
		if price, err := settledTokenPriceFromMarkets(ev.Markets, asset); err == nil {
			return price, nil
		}
	}
	return 0, fmt.Errorf("settled event price not found")
}

func settledTokenPriceFromMarkets(markets []gammaMarket, asset string) (float64, error) {
	for _, m := range markets {
		if !m.Closed {
			continue
		}
		tokens := parseJSONStringArray(m.ClobTokenIDsRaw)
		prices := parseJSONStringArray(m.OutcomePricesRaw)
		for i, token := range tokens {
			if strings.TrimSpace(token) != asset || i >= len(prices) {
				continue
			}
			price, err := strconv.ParseFloat(strings.TrimSpace(prices[i]), 64)
			if err != nil || price < 0 || price > 1 {
				continue
			}
			return price, nil
		}
	}
	return 0, fmt.Errorf("settled price not found")
}

func parseJSONStringArray(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func writeReport(path, logPath, extraLogPath string, results []alertResult, stakeUSD float64, currentPolicyModes map[string]struct{}, currentExcludeWallets map[string]struct{}, currentExcludePath string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	sum := summarize(results)
	activeResults, excludedWallets, excludedMarkets := filterActiveResults(results, currentExcludeWallets)
	currentResults := filterResultsByModes(activeResults, currentPolicyModes)
	currentSum := summarize(currentResults)
	currentAction, currentReason := gateAction(currentSum)
	currentPositionResults := capFirstByWalletAsset(currentResults)
	currentPositionSum := summarize(currentPositionResults)
	currentPositionAction, currentPositionReason := gateAction(currentPositionSum)
	modeBuckets := modeDecisionBuckets(activeResults)
	effectivePolicyModes := effectiveModeSet(currentPolicyModes, blockedModes(modeBuckets))
	effectiveResults := filterResultsByModes(activeResults, effectivePolicyModes)
	effectiveSum := summarize(effectiveResults)
	effectiveAction, effectiveReason := gateAction(effectiveSum)
	effectivePositionResults := capFirstByWalletAsset(effectiveResults)
	effectivePositionSum := summarize(effectivePositionResults)
	effectivePositionAction, effectivePositionReason := gateAction(effectivePositionSum)
	observeResults := filterResultsByModes(activeResults, parseModeSet("OBSERVE,OBSERVE-BURST"))
	observeSum := summarize(observeResults)
	observeAction, observeReason := gateAction(observeSum)
	observePositionResults := capFirstByWalletAsset(observeResults)
	observePositionSum := summarize(observePositionResults)
	observePositionAction, observePositionReason := gateAction(observePositionSum)
	observeBurstResults := filterResultsByModes(activeResults, parseModeSet("OBSERVE-BURST"))
	observeBurstSum := summarize(observeBurstResults)
	observeBurstAction, observeBurstReason := gateAction(observeBurstSum)
	observeBurstPositionResults := capFirstByWalletAsset(observeBurstResults)
	observeBurstPositionSum := summarize(observeBurstPositionResults)
	observeBurstPositionAction, observeBurstPositionReason := gateAction(observeBurstPositionSum)
	unknownFlowResults := filterResultsByModes(activeResults, parseModeSet("UNKNOWN-FLOW"))
	unknownFlowSum := summarize(unknownFlowResults)
	unknownFlowAction, unknownFlowReason := gateAction(unknownFlowSum)
	unknownFlowPositionResults := capFirstByWalletAsset(unknownFlowResults)
	unknownFlowPositionSum := summarize(unknownFlowPositionResults)
	unknownFlowPositionAction, unknownFlowPositionReason := gateAction(unknownFlowPositionSum)
	seedFlowResults := filterResultsByModes(activeResults, parseModeSet("SEED-FLOW"))
	seedFlowSum := summarize(seedFlowResults)
	seedFlowAction, seedFlowReason := gateAction(seedFlowSum)
	seedFlowPositionResults := capFirstByWalletAsset(seedFlowResults)
	seedFlowPositionSum := summarize(seedFlowPositionResults)
	seedFlowPositionAction, seedFlowPositionReason := gateAction(seedFlowPositionSum)
	scoredFlowResults := filterResultsByModes(activeResults, parseModeSet("SCORED-FLOW"))
	scoredFlowSum := summarize(scoredFlowResults)
	scoredFlowAction, scoredFlowReason := gateAction(scoredFlowSum)
	scoredFlowPositionResults := capFirstByWalletAsset(scoredFlowResults)
	scoredFlowPositionSum := summarize(scoredFlowPositionResults)
	scoredFlowPositionAction, scoredFlowPositionReason := gateAction(scoredFlowPositionSum)
	insiderResults := filterResultsByModes(activeResults, parseModeSet("INSIDER-SCOUT"))
	insiderSum := summarize(insiderResults)
	insiderAction, insiderReason := gateAction(insiderSum)
	insiderPositionResults := capFirstByWalletAsset(insiderResults)
	insiderPositionSum := summarize(insiderPositionResults)
	insiderPositionAction, insiderPositionReason := gateAction(insiderPositionSum)
	consensusResults := filterResultsByModes(activeResults, parseModeSet("CONSENSUS"))
	consensusSum := summarize(consensusResults)
	consensusAction, consensusReason := gateAction(consensusSum)
	consensusPositionResults := capFirstByWalletAsset(consensusResults)
	consensusPositionSum := summarize(consensusPositionResults)
	consensusPositionAction, consensusPositionReason := gateAction(consensusPositionSum)
	var b strings.Builder
	fmt.Fprintf(&b, "# Sports Alert Performance\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Alert log: `%s`\n", logPath)
	if strings.TrimSpace(extraLogPath) != "" {
		fmt.Fprintf(&b, "- Extra alert log: `%s`\n", extraLogPath)
	}
	fmt.Fprintf(&b, "- Fixed paper stake: $%.2f per Telegram alert\n", stakeUSD)
	fmt.Fprintf(&b, "- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices\n\n")
	if currentExcludePath != "" {
		fmt.Fprintf(&b, "- Current policy exclude wallets: `%s`\n", currentExcludePath)
		fmt.Fprintf(&b, "- Historical alerts excluded by wallet lists: %d\n", excludedWallets)
	}
	fmt.Fprintf(&b, "- Historical alerts excluded by current market filter: %d\n\n", excludedMarkets)
	fmt.Fprintf(&b, "## Summary\n\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", sum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", sum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", sum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(sum.Wins), float64(sum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", sum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", sum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n\n", sum.AvgDeltaP)

	fmt.Fprintf(&b, "## Current Policy Performance\n\n")
	fmt.Fprintf(&b, "- Modes: %s\n", formatModeSet(currentPolicyModes))
	fmt.Fprintf(&b, "- Alerts: %d\n", currentSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", currentSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", currentSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(currentSum.Wins), float64(currentSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", currentSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", currentSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", currentSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", currentAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", currentReason)

	fmt.Fprintf(&b, "## Current Policy Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Modes: %s\n", formatModeSet(currentPolicyModes))
	fmt.Fprintf(&b, "- Positions: %d\n", currentPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", currentPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", currentPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(currentPositionSum.Wins), float64(currentPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", currentPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", currentPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", currentPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", currentPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", currentPositionReason)

	fmt.Fprintf(&b, "## Effective Tradable Policy Performance\n\n")
	fmt.Fprintf(&b, "- Rule: current policy modes after removing CUT/PROBATION modes\n")
	fmt.Fprintf(&b, "- Modes: %s\n", formatModeSet(effectivePolicyModes))
	fmt.Fprintf(&b, "- Alerts: %d\n", effectiveSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", effectiveSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", effectiveSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(effectiveSum.Wins), float64(effectiveSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", effectiveSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", effectiveSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", effectiveSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", effectiveAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", effectiveReason)

	fmt.Fprintf(&b, "## Effective Tradable Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first alert per wallet + asset after removing CUT/PROBATION modes\n")
	fmt.Fprintf(&b, "- Modes: %s\n", formatModeSet(effectivePolicyModes))
	fmt.Fprintf(&b, "- Positions: %d\n", effectivePositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", effectivePositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", effectivePositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(effectivePositionSum.Wins), float64(effectivePositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", effectivePositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", effectivePositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", effectivePositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", effectivePositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", effectivePositionReason)

	fmt.Fprintf(&b, "## Experimental OBSERVE Performance\n\n")
	fmt.Fprintf(&b, "- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted\n")
	fmt.Fprintf(&b, "- Modes: OBSERVE,OBSERVE-BURST\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", observeSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", observeSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", observeSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(observeSum.Wins), float64(observeSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", observeSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", observeSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", observeSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", observeAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", observeReason)

	fmt.Fprintf(&b, "## Experimental OBSERVE Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first OBSERVE/OBSERVE-BURST alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", observePositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", observePositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", observePositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(observePositionSum.Wins), float64(observePositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", observePositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", observePositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", observePositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", observePositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", observePositionReason)

	fmt.Fprintf(&b, "## Experimental OBSERVE-BURST Performance\n\n")
	fmt.Fprintf(&b, "- Rule: same-wallet split-order sports/esports BUY bursts; not counted in current policy until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: OBSERVE-BURST\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", observeBurstSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", observeBurstSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", observeBurstSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(observeBurstSum.Wins), float64(observeBurstSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", observeBurstSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", observeBurstSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", observeBurstSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", observeBurstAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", observeBurstReason)

	fmt.Fprintf(&b, "## Experimental OBSERVE-BURST Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first OBSERVE-BURST alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", observeBurstPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", observeBurstPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", observeBurstPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(observeBurstPositionSum.Wins), float64(observeBurstPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", observeBurstPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", observeBurstPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", observeBurstPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", observeBurstPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", observeBurstPositionReason)

	fmt.Fprintf(&b, "## Experimental UNKNOWN-FLOW Performance\n\n")
	fmt.Fprintf(&b, "- Rule: shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: UNKNOWN-FLOW\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", unknownFlowSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", unknownFlowSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", unknownFlowSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(unknownFlowSum.Wins), float64(unknownFlowSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", unknownFlowSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", unknownFlowSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", unknownFlowSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", unknownFlowAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", unknownFlowReason)

	fmt.Fprintf(&b, "## Experimental UNKNOWN-FLOW Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first UNKNOWN-FLOW alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", unknownFlowPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", unknownFlowPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", unknownFlowPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(unknownFlowPositionSum.Wins), float64(unknownFlowPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", unknownFlowPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", unknownFlowPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", unknownFlowPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", unknownFlowPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", unknownFlowPositionReason)

	fmt.Fprintf(&b, "## Experimental SEED-FLOW Performance\n\n")
	fmt.Fprintf(&b, "- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: SEED-FLOW\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", seedFlowSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", seedFlowSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", seedFlowSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(seedFlowSum.Wins), float64(seedFlowSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", seedFlowSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", seedFlowSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", seedFlowSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", seedFlowAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", seedFlowReason)

	fmt.Fprintf(&b, "## Experimental SEED-FLOW Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first SEED-FLOW alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", seedFlowPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", seedFlowPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", seedFlowPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(seedFlowPositionSum.Wins), float64(seedFlowPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", seedFlowPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", seedFlowPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", seedFlowPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", seedFlowPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", seedFlowPositionReason)

	fmt.Fprintf(&b, "## Experimental SCORED-FLOW Performance\n\n")
	fmt.Fprintf(&b, "- Rule: shadow-only scored low-bot leaderboard wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: SCORED-FLOW\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", scoredFlowSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", scoredFlowSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", scoredFlowSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(scoredFlowSum.Wins), float64(scoredFlowSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", scoredFlowSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", scoredFlowSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", scoredFlowSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", scoredFlowAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", scoredFlowReason)

	fmt.Fprintf(&b, "## Experimental SCORED-FLOW Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first SCORED-FLOW alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", scoredFlowPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", scoredFlowPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", scoredFlowPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(scoredFlowPositionSum.Wins), float64(scoredFlowPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", scoredFlowPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", scoredFlowPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", scoredFlowPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", scoredFlowPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", scoredFlowPositionReason)

	fmt.Fprintf(&b, "## Experimental INSIDER-SCOUT Performance\n\n")
	fmt.Fprintf(&b, "- Rule: very large low-bot sports/esports whale BUY alerts; observation only until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: INSIDER-SCOUT\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", insiderSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", insiderSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", insiderSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(insiderSum.Wins), float64(insiderSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", insiderSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", insiderSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", insiderSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", insiderAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", insiderReason)

	fmt.Fprintf(&b, "## Experimental INSIDER-SCOUT Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first INSIDER-SCOUT alert per wallet + asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", insiderPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", insiderPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", insiderPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(insiderPositionSum.Wins), float64(insiderPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", insiderPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", insiderPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", insiderPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", insiderPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", insiderPositionReason)

	fmt.Fprintf(&b, "## Experimental CONSENSUS Performance\n\n")
	fmt.Fprintf(&b, "- Rule: cross-wallet same-asset sports/esports BUY bursts; research/observation only until repeated positive ROI is proven\n")
	fmt.Fprintf(&b, "- Modes: CONSENSUS\n")
	fmt.Fprintf(&b, "- Alerts: %d\n", consensusSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", consensusSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", consensusSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(consensusSum.Wins), float64(consensusSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", consensusSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", consensusSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", consensusSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", consensusAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", consensusReason)

	fmt.Fprintf(&b, "## Experimental CONSENSUS Position-Capped Performance\n\n")
	fmt.Fprintf(&b, "- Rule: first CONSENSUS alert per asset\n")
	fmt.Fprintf(&b, "- Positions: %d\n", consensusPositionSum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", consensusPositionSum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", consensusPositionSum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(consensusPositionSum.Wins), float64(consensusPositionSum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", consensusPositionSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", consensusPositionSum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n", consensusPositionSum.AvgDeltaP)
	fmt.Fprintf(&b, "- Gate action: %s\n", consensusPositionAction)
	fmt.Fprintf(&b, "- Gate reason: %s\n\n", consensusPositionReason)

	writeGateSection(&b, "Mode Gates", activeResults, func(r alertResult) string { return firstNonEmpty(r.Event.Mode, "OBSERVE") })
	writeGateSection(&b, "Wallet Gates", activeResults, func(r alertResult) string { return shortAddr(r.Event.Wallet) })
	writeGroupSection(&b, "By Mode", activeResults, func(r alertResult) string { return firstNonEmpty(r.Event.Mode, "OBSERVE") })
	writeGroupSection(&b, "By Wallet", activeResults, func(r alertResult) string { return shortAddr(r.Event.Wallet) })
	writeRecentSection(&b, results)
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeDecisionJSON(path, logPath, extraLogPath string, results []alertResult, stakeUSD float64, currentPolicyModes map[string]struct{}, currentExcludeWallets map[string]struct{}) error {
	activeResults, _, _ := filterActiveResults(results, currentExcludeWallets)
	currentResults := filterResultsByModes(activeResults, currentPolicyModes)
	currentPositionResults := capFirstByWalletAsset(currentResults)
	modes := modeDecisionBuckets(activeResults)
	effectiveModes := effectiveModeSet(currentPolicyModes, blockedModes(modes))
	effectiveResults := filterResultsByModes(activeResults, effectiveModes)
	effectivePositionResults := capFirstByWalletAsset(effectiveResults)

	recommended := []string{}
	positive := []string{}
	cut := []string{}
	probation := []string{}
	for _, b := range modes {
		switch b.Action {
		case "PROMOTE_CANDIDATE":
			recommended = append(recommended, b.Key)
		case "COLLECT_POSITIVE":
			positive = append(positive, b.Key)
		case "CUT":
			cut = append(cut, b.Key)
		case "PROBATION":
			probation = append(probation, b.Key)
		}
	}

	report := policyDecisionReport{
		GeneratedAt:                   time.Now(),
		Log:                           logPath,
		ExtraLog:                      strings.TrimSpace(extraLogPath),
		StakeUSD:                      stakeUSD,
		CurrentPolicy:                 decisionBucketFor("CURRENT_POLICY", currentResults),
		CurrentPolicyPositionCapped:   decisionBucketFor("CURRENT_POLICY_POSITION_CAPPED", currentPositionResults),
		EffectivePolicy:               decisionBucketFor("EFFECTIVE_POLICY", effectiveResults),
		EffectivePolicyPositionCapped: decisionBucketFor("EFFECTIVE_POLICY_POSITION_CAPPED", effectivePositionResults),
		EffectiveModes:                sortedModeSet(effectiveModes),
		Modes:                         modes,
		RecommendedModes:              recommended,
		PositiveModes:                 positive,
		CutModes:                      cut,
		ProbationModes:                probation,
	}
	return writeJSONFile(path, report)
}

func modeDecisionBuckets(results []alertResult) []decisionBucket {
	groups := map[string][]alertResult{}
	for _, r := range results {
		mode := firstNonEmpty(r.Event.Mode, "OBSERVE")
		groups[mode] = append(groups[mode], r)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		ai := decisionBucketFor(keys[i], groups[keys[i]])
		aj := decisionBucketFor(keys[j], groups[keys[j]])
		if actionRank(ai.Action) != actionRank(aj.Action) {
			return actionRank(ai.Action) > actionRank(aj.Action)
		}
		if ai.Marked != aj.Marked {
			return ai.Marked > aj.Marked
		}
		if ai.ROI != aj.ROI {
			return ai.ROI > aj.ROI
		}
		return keys[i] < keys[j]
	})

	out := make([]decisionBucket, 0, len(keys))
	for _, key := range keys {
		out = append(out, decisionBucketFor(key, groups[key]))
	}
	return out
}

func blockedModes(buckets []decisionBucket) map[string]struct{} {
	out := map[string]struct{}{}
	for _, b := range buckets {
		switch strings.ToUpper(strings.TrimSpace(b.Action)) {
		case "CUT", "PROBATION":
			key := strings.ToUpper(strings.TrimSpace(b.Key))
			if key != "" {
				out[key] = struct{}{}
			}
		}
	}
	return out
}

func effectiveModeSet(currentModes map[string]struct{}, blocked map[string]struct{}) map[string]struct{} {
	out := map[string]struct{}{}
	for mode := range currentModes {
		mode = strings.ToUpper(strings.TrimSpace(mode))
		if mode == "" {
			continue
		}
		if _, ok := blocked[mode]; ok {
			continue
		}
		out[mode] = struct{}{}
	}
	return out
}

func sortedModeSet(modes map[string]struct{}) []string {
	out := make([]string, 0, len(modes))
	for mode := range modes {
		out = append(out, mode)
	}
	sort.Strings(out)
	return out
}

func decisionBucketFor(key string, results []alertResult) decisionBucket {
	sum := summarize(results)
	action, reason := gateAction(sum)
	return decisionBucket{
		Key:        key,
		Action:     action,
		Reason:     reason,
		Alerts:     sum.Signals,
		Marked:     sum.Marked,
		Unmarked:   sum.Unmarked,
		Wins:       sum.Wins,
		WinRate:    pct(float64(sum.Wins), float64(sum.Marked)),
		ROI:        sum.ReturnPC,
		PnLUSD:     sum.PnLUSD,
		AvgDeltaPP: sum.AvgDeltaP,
	}
}

func actionRank(action string) int {
	switch action {
	case "PROMOTE_CANDIDATE":
		return 5
	case "COLLECT_POSITIVE":
		return 4
	case "COLLECT":
		return 3
	case "PROBATION":
		return 2
	case "CUT":
		return 1
	default:
		return 0
	}
}

func writeJSONFile(path string, v any) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

func filterActiveResults(results []alertResult, blocked map[string]struct{}) ([]alertResult, int, int) {
	out := make([]alertResult, 0, len(results))
	var excludedWallets int
	var excludedMarkets int
	for _, r := range results {
		addr := strings.ToLower(strings.TrimSpace(r.Event.Wallet))
		if _, ok := blocked[addr]; ok {
			excludedWallets++
			continue
		}
		if !isCurrentPolicyMarket(r.Event) {
			excludedMarkets++
			continue
		}
		out = append(out, r)
	}
	return out, excludedWallets, excludedMarkets
}

func isCurrentPolicyMarket(ev alertEvent) bool {
	q := strings.TrimSpace(ev.Market)
	slug := strings.TrimSpace(ev.Slug)
	text := strings.ToLower(q + " " + slug + " " + ev.Category)
	if feed.IsDerivativeFollowMarketText(text) || feed.IsOutrightFollowMarketText(text) {
		return false
	}
	if strings.Contains(text, "tennis") ||
		strings.Contains(text, "wimbledon") ||
		strings.Contains(text, " atp") ||
		strings.Contains(text, " wta") ||
		strings.HasPrefix(strings.ToLower(slug), "atp-") ||
		strings.HasPrefix(strings.ToLower(slug), "wta-") {
		return false
	}
	if feed.IsFollowTargetMarket(feed.Market{Question: q, Slug: slug}) {
		return true
	}
	category := strings.ToLower(strings.TrimSpace(ev.Category))
	switch category {
	case "basketball", "soccer", "football", "esports", "lol", "dota", "dota2":
		return true
	}
	for _, group := range [][]string{
		{"nba", "wnba", "basketball"},
		{"epl", "premier league", "la liga", "bundesliga", "serie a", "ligue 1", "champions league", "ucl", "uefa", "fifa", "fifwc", "fifa world cup", "copa ", "concacaf", "conmebol", "eredivisie", "liga mx", "mls", "soccer", "futbol"},
		{"lol", "league of legends", "lck", "lpl", "msi", "dota", "dota2", "esport"},
	} {
		for _, k := range group {
			if strings.Contains(text, k) {
				return true
			}
		}
	}
	return false
}

func filterResultsByModes(results []alertResult, modes map[string]struct{}) []alertResult {
	if len(modes) == 0 {
		return nil
	}
	out := make([]alertResult, 0, len(results))
	for _, r := range results {
		mode := strings.ToUpper(strings.TrimSpace(firstNonEmpty(r.Event.Mode, "OBSERVE")))
		if _, ok := modes[mode]; ok {
			out = append(out, r)
		}
	}
	return out
}

func capFirstByWalletAsset(results []alertResult) []alertResult {
	rows := append([]alertResult(nil), results...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].Event.SentAt.Before(rows[j].Event.SentAt) })
	seen := map[string]struct{}{}
	out := make([]alertResult, 0, len(rows))
	for _, r := range rows {
		mode := strings.ToUpper(strings.TrimSpace(firstNonEmpty(r.Event.Mode, "OBSERVE")))
		key := strings.ToLower(strings.TrimSpace(r.Event.Wallet)) + "|" + strings.TrimSpace(r.Event.Asset)
		if mode == "CONSENSUS" {
			key = "consensus|" + strings.TrimSpace(r.Event.Asset)
		}
		if key == "|" {
			key = strings.ToLower(strings.TrimSpace(r.Event.Key))
		}
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, r)
	}
	return out
}

func writeGateSection(b *strings.Builder, title string, results []alertResult, keyFn func(alertResult) string) {
	groups := map[string][]alertResult{}
	for _, r := range results {
		groups[keyFn(r)] = append(groups[keyFn(r)], r)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ai := summarize(groups[keys[i]])
		aj := summarize(groups[keys[j]])
		if ai.Marked != aj.Marked {
			return ai.Marked > aj.Marked
		}
		if ai.ReturnPC != aj.ReturnPC {
			return ai.ReturnPC > aj.ReturnPC
		}
		return keys[i] < keys[j]
	})
	fmt.Fprintf(b, "## %s\n\n", title)
	fmt.Fprintf(b, "- Promote gate: marked>=5, ROI>=5.0%%, win>=60.0%%\n")
	fmt.Fprintf(b, "- Cut gate: marked>=3 and ROI<=-20.0%% or win<=35.0%%; severe single loss is probation\n\n")
	fmt.Fprintf(b, "| Key | Alerts | Marked | Win | ROI | PnL | Action | Reason |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---|---|\n")
	if len(keys) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | COLLECT | no alerts yet |\n\n")
		return
	}
	for _, k := range keys {
		s := summarize(groups[k])
		action, reason := gateAction(s)
		fmt.Fprintf(b, "| `%s` | %d | %d | %.1f%% | %.1f%% | $%+.2f | %s | %s |\n",
			k, s.Signals, s.Marked, pct(float64(s.Wins), float64(s.Marked)), s.ReturnPC, s.PnLUSD, action, reason)
	}
	fmt.Fprintf(b, "\n")
}

func gateAction(s resultSummary) (string, string) {
	winRate := pct(float64(s.Wins), float64(s.Marked))
	switch {
	case s.Marked == 0:
		return "COLLECT", "no marked alerts yet"
	case s.Marked >= 5 && s.ReturnPC >= 5 && winRate >= 60:
		return "PROMOTE_CANDIDATE", "enough positive marked alerts"
	case s.Marked >= 3 && (s.ReturnPC <= -20 || winRate <= 35):
		return "CUT", "negative marked sample"
	case s.Marked >= 1 && s.ReturnPC <= -50:
		return "PROBATION", "severe drawdown on limited sample"
	case s.ReturnPC > 0:
		return "COLLECT_POSITIVE", "positive but sample below promote gate"
	default:
		return "COLLECT", "sample below promote gate"
	}
}

func writeGroupSection(b *strings.Builder, title string, results []alertResult, keyFn func(alertResult) string) {
	groups := map[string][]alertResult{}
	for _, r := range results {
		groups[keyFn(r)] = append(groups[keyFn(r)], r)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ai := summarize(groups[keys[i]])
		aj := summarize(groups[keys[j]])
		if ai.PnLUSD != aj.PnLUSD {
			return ai.PnLUSD > aj.PnLUSD
		}
		return keys[i] < keys[j]
	})
	fmt.Fprintf(b, "## %s\n\n", title)
	fmt.Fprintf(b, "| Key | Alerts | Marked | Win | ROI | PnL | AvgDeltaPP |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|\n")
	if len(keys) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | +0.00 |\n\n")
		return
	}
	for _, k := range keys {
		s := summarize(groups[k])
		fmt.Fprintf(b, "| `%s` | %d | %d | %.1f%% | %.1f%% | $%+.2f | %+.2f |\n",
			k, s.Signals, s.Marked, pct(float64(s.Wins), float64(s.Marked)), s.ReturnPC, s.PnLUSD, s.AvgDeltaP)
	}
	fmt.Fprintf(b, "\n")
}

func writeRecentSection(b *strings.Builder, results []alertResult) {
	rows := append([]alertResult(nil), results...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].Event.SentAt.After(rows[j].Event.SentAt) })
	limit := minInt(len(rows), 20)
	fmt.Fprintf(b, "## Recent Alerts\n\n")
	fmt.Fprintf(b, "| Sent | Mode | Wallet | Notional | Entry | Mid | ROI | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---|\n")
	if limit == 0 {
		fmt.Fprintf(b, "| n/a |  |  | $0 | 0.000 | 0.000 | 0.0%% |  |\n")
		return
	}
	for i := 0; i < limit; i++ {
		r := rows[i]
		mid := "-"
		if r.Marked {
			mid = fmt.Sprintf("%.3f", r.Mid)
		}
		fmt.Fprintf(b, "| %s | %s | `%s` | $%.0f | %.3f | %s | %.1f%% | %s |\n",
			formatTime(r.Event.SentAt), firstNonEmpty(r.Event.Mode, "OBSERVE"), shortAddr(r.Event.Wallet),
			r.Event.Notional, r.Event.Price, mid, r.ReturnPC, oneLine(r.Event.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func summarize(results []alertResult) resultSummary {
	var s resultSummary
	for _, r := range results {
		s.Signals++
		if !r.Marked {
			s.Unmarked++
			continue
		}
		s.Marked++
		if r.PnLUSD > 0 {
			s.Wins++
		}
		s.StakeUSD += r.StakeUSD
		s.PnLUSD += r.PnLUSD
		s.AvgDeltaP += (r.Mid - r.Event.Price) * 100
	}
	if s.StakeUSD > 0 {
		s.ReturnPC = s.PnLUSD / s.StakeUSD * 100
	}
	if s.Marked > 0 {
		s.AvgDeltaP /= float64(s.Marked)
	}
	return s
}

func pct(n, d float64) float64 {
	if d <= 0 {
		return 0
	}
	return n / d * 100
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("01-02 15:04")
}

func oneLine(s string, limit int) string {
	s = strings.Join(strings.Fields(s), " ")
	if limit > 0 && len(s) > limit {
		return s[:limit-1] + "..."
	}
	return s
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}

func parseModeSet(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.ToUpper(strings.TrimSpace(part))
		if part != "" {
			out[part] = struct{}{}
		}
	}
	return out
}

func formatModeSet(modes map[string]struct{}) string {
	if len(modes) == 0 {
		return "-"
	}
	out := make([]string, 0, len(modes))
	for mode := range modes {
		out = append(out, mode)
	}
	sort.Strings(out)
	return strings.Join(out, ",")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
