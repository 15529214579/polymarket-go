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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

type whaleTrade struct {
	TS          string  `json:"ts"`
	Wallet      string  `json:"wallet"`
	Label       string  `json:"label"`
	Side        string  `json:"side"`
	Market      string  `json:"market"`
	Outcome     string  `json:"outcome"`
	Price       float64 `json:"price"`
	Size        float64 `json:"size"`
	Units       float64 `json:"units"`
	AssetID     string  `json:"asset_id"`
	ConditionID string  `json:"condition_id"`
	TradeID     string  `json:"trade_id"`
	Action      string  `json:"action"`
	Reason      string  `json:"reason"`
	List        string  `json:"list"`
	Tier        string  `json:"tier"`
	Smart       float64 `json:"smart"`
	Bot         float64 `json:"bot"`

	Time time.Time `json:"-"`
}

type signalResult struct {
	Buy        whaleTrade
	Exit       *whaleTrade
	ExitSource string
	ExitPrice  float64
	StakeUSD   float64
	Units      float64
	PnLUSD     float64
	ReturnPct  float64
	HoldHours  float64
}

type evalOptions struct {
	StakeUSD          float64
	MinNotional       float64
	ListMinNotional   map[string]float64
	Since             time.Time
	RepeatCooldown    time.Duration
	RepeatMinNotional float64
}

type evaluation struct {
	Results              []*signalResult
	RawBuys              int
	SuppressedRepeats    int
	LoggedCooldownBuys   int
	LoggedEventCooldowns int
	LoggedDuplicateBuys  int
	LoggedPendingBuys    int
	DuplicateBuys        int
	SuppressedByWallet   []suppressionStats
	SuppressedByEvent    []suppressionStats
	PolicyViolations     []policyViolation
}

type walletStats struct {
	Wallet       string
	Label        string
	List         string
	Tier         string
	Signals      int
	Closed       int
	Settled      int
	Marked       int
	Open         int
	Wins         int
	PnLUSD       float64
	StakeUSD     float64
	ReturnPct    float64
	AvgReturnPct float64
}

type eventStats struct {
	Event        string
	Wallet       string
	Label        string
	List         string
	Tier         string
	Signals      int
	Closed       int
	Settled      int
	Marked       int
	Wins         int
	Markets      int
	PnLUSD       float64
	StakeUSD     float64
	ReturnPct    float64
	AvgReturnPct float64
	FirstTime    time.Time
	LastTime     time.Time
}

type suppressionStats struct {
	Wallet        string
	Label         string
	List          string
	Tier          string
	Event         string
	AssetCooldown int
	EventCooldown int
	Duplicate     int
	Pending       int
	Total         int
	NotionalUSD   float64
	LastTime      time.Time
}

type policyViolation struct {
	Time        time.Time
	Wallet      string
	Label       string
	List        string
	Tier        string
	Reason      string
	Market      string
	Outcome     string
	Price       float64
	NotionalUSD float64
	TradeID     string
}

type walletMeta struct {
	List string
	Tier string
}

type snapshotStats struct {
	Wallet       string  `json:"wallet,omitempty"`
	Label        string  `json:"label,omitempty"`
	List         string  `json:"list,omitempty"`
	Tier         string  `json:"tier,omitempty"`
	Signals      int     `json:"signals"`
	Closed       int     `json:"closed"`
	Settled      int     `json:"settled"`
	Marked       int     `json:"marked"`
	Open         int     `json:"open"`
	Wins         int     `json:"wins"`
	PnLUSD       float64 `json:"pnl_usd"`
	StakeUSD     float64 `json:"stake_usd"`
	ReturnPct    float64 `json:"return_pct"`
	AvgReturnPct float64 `json:"avg_return_pct"`
}

type snapshotEventStats struct {
	Event        string  `json:"event"`
	Wallet       string  `json:"wallet,omitempty"`
	Label        string  `json:"label,omitempty"`
	List         string  `json:"list,omitempty"`
	Tier         string  `json:"tier,omitempty"`
	Signals      int     `json:"signals"`
	Closed       int     `json:"closed"`
	Settled      int     `json:"settled"`
	Marked       int     `json:"marked"`
	Wins         int     `json:"wins"`
	Markets      int     `json:"markets"`
	PnLUSD       float64 `json:"pnl_usd"`
	StakeUSD     float64 `json:"stake_usd"`
	ReturnPct    float64 `json:"return_pct"`
	AvgReturnPct float64 `json:"avg_return_pct"`
	FirstTime    string  `json:"first_time,omitempty"`
	LastTime     string  `json:"last_time,omitempty"`
}

type snapshotSuppressionStats struct {
	Wallet        string  `json:"wallet,omitempty"`
	Label         string  `json:"label,omitempty"`
	List          string  `json:"list,omitempty"`
	Tier          string  `json:"tier,omitempty"`
	Event         string  `json:"event,omitempty"`
	AssetCooldown int     `json:"asset_cooldown"`
	EventCooldown int     `json:"event_cooldown"`
	Duplicate     int     `json:"duplicate"`
	Pending       int     `json:"pending_consensus"`
	Total         int     `json:"total"`
	NotionalUSD   float64 `json:"notional_usd"`
	LastTime      string  `json:"last_time,omitempty"`
}

type snapshotPolicyViolation struct {
	Time        string  `json:"time"`
	Wallet      string  `json:"wallet,omitempty"`
	Label       string  `json:"label,omitempty"`
	List        string  `json:"list,omitempty"`
	Tier        string  `json:"tier,omitempty"`
	Reason      string  `json:"reason"`
	Market      string  `json:"market"`
	Outcome     string  `json:"outcome,omitempty"`
	Price       float64 `json:"price"`
	NotionalUSD float64 `json:"notional_usd"`
	TradeID     string  `json:"trade_id,omitempty"`
}

type performanceSnapshot struct {
	GeneratedAt                 string                     `json:"generated_at"`
	LogPath                     string                     `json:"log_path"`
	WalletsPath                 string                     `json:"wallets_path,omitempty"`
	ReportPath                  string                     `json:"report_path,omitempty"`
	StakeUSD                    float64                    `json:"stake_usd"`
	MinNotional                 float64                    `json:"min_notional"`
	ListMinNotional             map[string]float64         `json:"list_min_notional,omitempty"`
	Since                       string                     `json:"since,omitempty"`
	RepeatCooldown              string                     `json:"repeat_cooldown,omitempty"`
	RepeatMinNotional           float64                    `json:"repeat_min_notional,omitempty"`
	RawBuys                     int                        `json:"raw_buys"`
	SuppressedRepeats           int                        `json:"suppressed_repeats"`
	LoggedCooldownBuys          int                        `json:"logged_cooldown_buys"`
	LoggedEventCooldowns        int                        `json:"logged_event_cooldown_buys"`
	LoggedDuplicateBuys         int                        `json:"logged_duplicate_buys"`
	LoggedPendingBuys           int                        `json:"logged_pending_consensus_buys"`
	DuplicateBuys               int                        `json:"duplicate_buys"`
	EvaluatedSignals            int                        `json:"evaluated_signals"`
	Closed                      int                        `json:"closed"`
	Settled                     int                        `json:"settled"`
	Marked                      int                        `json:"marked"`
	Open                        int                        `json:"open"`
	Wins                        int                        `json:"wins"`
	PnLUSD                      float64                    `json:"pnl_usd"`
	StakeEvaluatedUSD           float64                    `json:"stake_evaluated_usd"`
	ReturnPct                   float64                    `json:"return_pct"`
	WinRatePct                  float64                    `json:"win_rate_pct"`
	ProvenSignals               int                        `json:"proven_signals"`
	ProvenWins                  int                        `json:"proven_wins"`
	ProvenPnLUSD                float64                    `json:"proven_pnl_usd"`
	ProvenStakeUSD              float64                    `json:"proven_stake_usd"`
	ProvenReturnPct             float64                    `json:"proven_return_pct"`
	ProvenWinRatePct            float64                    `json:"proven_win_rate_pct"`
	EventCappedSignals          int                        `json:"event_capped_signals"`
	EventCappedWins             int                        `json:"event_capped_wins"`
	EventCappedPnLUSD           float64                    `json:"event_capped_pnl_usd"`
	EventCappedStakeUSD         float64                    `json:"event_capped_stake_usd"`
	EventCappedReturnPct        float64                    `json:"event_capped_return_pct"`
	EventCappedWinRatePct       float64                    `json:"event_capped_win_rate_pct"`
	EventCappedProvenSignals    int                        `json:"event_capped_proven_signals"`
	EventCappedProvenWins       int                        `json:"event_capped_proven_wins"`
	EventCappedProvenPnLUSD     float64                    `json:"event_capped_proven_pnl_usd"`
	EventCappedProvenStakeUSD   float64                    `json:"event_capped_proven_stake_usd"`
	EventCappedProvenReturnPct  float64                    `json:"event_capped_proven_return_pct"`
	EventCappedProvenWinRatePct float64                    `json:"event_capped_proven_win_rate_pct"`
	EventCappedByList           []snapshotStats            `json:"event_capped_by_list,omitempty"`
	EventCappedByWallet         []snapshotStats            `json:"event_capped_by_wallet,omitempty"`
	EventCappedProvenByWallet   []snapshotStats            `json:"event_capped_proven_by_wallet,omitempty"`
	ByList                      []snapshotStats            `json:"by_list,omitempty"`
	ByWallet                    []snapshotStats            `json:"by_wallet,omitempty"`
	ByEvent                     []snapshotEventStats       `json:"by_event,omitempty"`
	SuppressedByWallet          []snapshotSuppressionStats `json:"suppressed_by_wallet,omitempty"`
	SuppressedByEvent           []snapshotSuppressionStats `json:"suppressed_by_event,omitempty"`
	PolicyViolationCount        int                        `json:"policy_violation_count"`
	PolicyViolations            []snapshotPolicyViolation  `json:"policy_violations,omitempty"`
}

func main() {
	logPath := flag.String("log", "db/journal/whale_trades.jsonl", "whale trades JSONL path")
	walletsPath := flag.String("wallets", "wallets.strategy-core.txt", "optional wallet allowlist")
	reportPath := flag.String("report", "reports/whale_performance.md", "markdown report path")
	summaryJSONPath := flag.String("summary_json", "", "optional path for latest structured summary JSON")
	snapshotJSONLPath := flag.String("snapshot_jsonl", "", "optional path to append structured summary snapshots as JSONL")
	stakeUSD := flag.Float64("stake", 10, "fixed stake used to evaluate each BUY signal")
	minNotional := flag.Float64("min_notional", 100, "minimum whale notional to evaluate")
	listMinNotional := flag.String("list_min_notional", "", "optional per-list minimum notional, e.g. core=1000,sports=1500")
	sinceRaw := flag.String("since", "", "only evaluate BUY/suppression signals at or after this RFC3339 timestamp")
	includeActions := flag.String("actions", "alert,followed", "comma-separated BUY actions to evaluate, or all")
	repeatCooldown := flag.Duration("repeat_cooldown", 0, "ignore repeat BUYs from the same wallet+asset inside this window; 0 disables")
	repeatMinNotional := flag.Float64("repeat_min_notional", 5000, "repeat BUYs at or above this notional bypass repeat cooldown")
	markCurrent := flag.Bool("mark_current", true, "mark open BUY signals to current CLOB midpoint")
	maxOpenAge := flag.Duration("max_open_age", 0, "if >0, only mark open signals at least this old")
	timeout := flag.Duration("timeout", 20*time.Second, "overall midpoint fetch timeout")
	flag.Parse()

	trades, err := loadWhaleTrades(*logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-report: %v\n", err)
		os.Exit(1)
	}
	allow, err := loadWalletMetas(*walletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-report: %v\n", err)
		os.Exit(1)
	}
	applyTradeMetas(trades, allow)
	actions := parseActionSet(*includeActions)
	listMins := parseListMinNotional(*listMinNotional)
	since, err := parseSince(*sinceRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-report: invalid -since: %v\n", err)
		os.Exit(1)
	}
	eval := buildSignalResults(trades, allow, actions, evalOptions{
		StakeUSD:          *stakeUSD,
		MinNotional:       *minNotional,
		ListMinNotional:   listMins,
		Since:             since,
		RepeatCooldown:    *repeatCooldown,
		RepeatMinNotional: *repeatMinNotional,
	})
	if *markCurrent {
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()
		markOpenSignals(ctx, eval.Results, *maxOpenAge)
	}
	if err := writeReport(*reportPath, *logPath, *walletsPath, eval, evalOptions{
		StakeUSD:          *stakeUSD,
		MinNotional:       *minNotional,
		ListMinNotional:   listMins,
		Since:             since,
		RepeatCooldown:    *repeatCooldown,
		RepeatMinNotional: *repeatMinNotional,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "whale-report: write report: %v\n", err)
		os.Exit(1)
	}
	snap := buildSnapshot(*logPath, *walletsPath, *reportPath, eval, evalOptions{
		StakeUSD:          *stakeUSD,
		MinNotional:       *minNotional,
		ListMinNotional:   listMins,
		Since:             since,
		RepeatCooldown:    *repeatCooldown,
		RepeatMinNotional: *repeatMinNotional,
	})
	if *summaryJSONPath != "" {
		if err := writeJSON(*summaryJSONPath, snap); err != nil {
			fmt.Fprintf(os.Stderr, "whale-report: write summary json: %v\n", err)
			os.Exit(1)
		}
	}
	if *snapshotJSONLPath != "" {
		if err := appendJSONL(*snapshotJSONLPath, snap); err != nil {
			fmt.Fprintf(os.Stderr, "whale-report: write snapshot jsonl: %v\n", err)
			os.Exit(1)
		}
	}
	sum := summarize(eval.Results)
	fmt.Printf("whale-report done: raw=%d suppressed=%d duplicate=%d evaluated=%d closed=%d marked=%d open=%d pnl=$%+.2f roi=%.1f%% report=%s\n",
		eval.RawBuys, eval.SuppressedRepeats, eval.DuplicateBuys, sum.Signals, countSource(eval.Results, "sell"), countSource(eval.Results, "mark"), countSource(eval.Results, ""),
		sum.PnLUSD, sum.ReturnPct, *reportPath)
}

func loadWhaleTrades(path string) ([]whaleTrade, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []whaleTrade
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 8*1024*1024)
	for sc.Scan() {
		var tr whaleTrade
		if err := json.Unmarshal(sc.Bytes(), &tr); err != nil {
			continue
		}
		ts, err := time.Parse(time.RFC3339, tr.TS)
		if err != nil {
			continue
		}
		tr.Time = ts
		tr.Wallet = strings.ToLower(strings.TrimSpace(tr.Wallet))
		tr.Side = strings.ToUpper(strings.TrimSpace(tr.Side))
		tr.Action = strings.ToLower(strings.TrimSpace(tr.Action))
		tr.List = strings.ToLower(strings.TrimSpace(tr.List))
		tr.Tier = strings.ToUpper(strings.TrimSpace(tr.Tier))
		if tr.Wallet == "" || tr.AssetID == "" || tr.Price <= 0 {
			continue
		}
		out = append(out, tr)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Before(out[j].Time) })
	return out, nil
}

func loadWalletMetas(path string) (map[string]walletMeta, error) {
	if path == "" {
		return nil, nil
	}
	if strings.Contains(path, ",") {
		out := map[string]walletMeta{}
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			metas, err := loadWalletMetas(part)
			if err != nil {
				return nil, err
			}
			for addr, meta := range metas {
				out[addr] = meta
			}
		}
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	out := map[string]walletMeta{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "#", 2)
		addr := strings.ToLower(strings.TrimSpace(parts[0]))
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		meta := walletMeta{}
		if len(parts) > 1 {
			for _, field := range strings.Fields(parts[1]) {
				k, v, ok := strings.Cut(field, "=")
				if !ok {
					continue
				}
				switch k {
				case "list":
					meta.List = strings.ToLower(strings.TrimSpace(v))
				case "tier":
					meta.Tier = strings.ToUpper(strings.TrimSpace(v))
				}
			}
		}
		out[addr] = meta
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func applyTradeMetas(trades []whaleTrade, metas map[string]walletMeta) {
	if len(metas) == 0 {
		return
	}
	for i := range trades {
		meta, ok := metas[trades[i].Wallet]
		if !ok {
			continue
		}
		if meta.List != "" {
			trades[i].List = meta.List
		}
		if meta.Tier != "" {
			trades[i].Tier = meta.Tier
		}
	}
}

func parseActionSet(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" {
			out[part] = struct{}{}
		}
	}
	return out
}

func parseListMinNotional(raw string) map[string]float64 {
	out := map[string]float64{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" {
			continue
		}
		minNotional, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil || minNotional < 0 {
			continue
		}
		out[key] = minNotional
	}
	return out
}

func parseSince(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, nil
	}
	ts, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, err
	}
	return ts, nil
}

func effectiveMinNotional(tr whaleTrade, opts evalOptions) float64 {
	list := strings.ToLower(strings.TrimSpace(tr.List))
	if list != "" {
		if v, ok := opts.ListMinNotional[list]; ok {
			return v
		}
	}
	return opts.MinNotional
}

func beforeSince(tr whaleTrade, opts evalOptions) bool {
	return !opts.Since.IsZero() && tr.Time.Before(opts.Since)
}

func buildSignalResults(trades []whaleTrade, allow map[string]walletMeta, actions map[string]struct{}, opts evalOptions) evaluation {
	var buys []whaleTrade
	var sellsByKey = map[string][]whaleTrade{}
	seenBuys := map[string]struct{}{}
	seenSells := map[string]struct{}{}
	suppressedBuyKeys := map[string]struct{}{}
	suppressedWallets := map[string]*suppressionStats{}
	suppressedEvents := map[string]*suppressionStats{}
	lastBuyByWalletAsset := map[string]time.Time{}
	eval := evaluation{}
	for _, tr := range trades {
		if allow != nil {
			if _, ok := allow[tr.Wallet]; !ok {
				continue
			}
		}
		if tr.Side == "BUY" && !beforeSince(tr, opts) && tr.Size >= effectiveMinNotional(tr, opts) && actionAllowed(actions, tr.Action) {
			if reason := reportPolicyViolationReason(tr); reason != "" {
				eval.PolicyViolations = append(eval.PolicyViolations, policyViolation{
					Time:        tr.Time,
					Wallet:      tr.Wallet,
					Label:       tr.Label,
					List:        tr.List,
					Tier:        tr.Tier,
					Reason:      reason,
					Market:      tr.Market,
					Outcome:     tr.Outcome,
					Price:       tr.Price,
					NotionalUSD: tr.Size,
					TradeID:     tr.TradeID,
				})
			}
		}
		if tr.Side == "BUY" && !beforeSince(tr, opts) && tr.Size >= effectiveMinNotional(tr, opts) {
			switch tr.Action {
			case "cooldown":
				eval.LoggedCooldownBuys++
			case "event_cooldown":
				eval.LoggedEventCooldowns++
			case "duplicate":
				eval.LoggedDuplicateBuys++
			case "pending_consensus":
				eval.LoggedPendingBuys++
			}
			if suppressingAction(tr.Action) {
				addSuppressionStats(suppressedWallets, tr.Wallet, tr, "")
				addSuppressionStats(suppressedEvents, tr.Wallet+"|"+eventKey(tr.Market), tr, eventKey(tr.Market))
			}
		}
		if tr.Side == "BUY" && !beforeSince(tr, opts) && suppressingAction(tr.Action) {
			suppressedBuyKeys[signalKey(tr)] = struct{}{}
		}
	}
	for _, tr := range trades {
		if allow != nil {
			if _, ok := allow[tr.Wallet]; !ok {
				continue
			}
		}
		key := tr.Wallet + "|" + tr.AssetID
		switch tr.Side {
		case "BUY":
			if beforeSince(tr, opts) {
				continue
			}
			if tr.Size < effectiveMinNotional(tr, opts) {
				continue
			}
			if !actionAllowed(actions, tr.Action) {
				continue
			}
			dedupeKey := signalKey(tr)
			if _, ok := seenBuys[dedupeKey]; ok {
				eval.DuplicateBuys++
				continue
			}
			seenBuys[dedupeKey] = struct{}{}
			eval.RawBuys++
			if _, ok := suppressedBuyKeys[dedupeKey]; ok {
				eval.SuppressedRepeats++
				continue
			}
			if opts.RepeatCooldown > 0 {
				if tr.Size >= opts.RepeatMinNotional {
					lastBuyByWalletAsset[key] = tr.Time
				} else if last, ok := lastBuyByWalletAsset[key]; ok && tr.Time.Sub(last) < opts.RepeatCooldown {
					eval.SuppressedRepeats++
					continue
				} else {
					lastBuyByWalletAsset[key] = tr.Time
				}
			}
			buys = append(buys, tr)
		case "SELL":
			dedupeKey := signalKey(tr)
			if _, ok := seenSells[dedupeKey]; ok {
				continue
			}
			seenSells[dedupeKey] = struct{}{}
			sellsByKey[key] = append(sellsByKey[key], tr)
		}
	}

	for _, buy := range buys {
		res := &signalResult{Buy: buy, StakeUSD: opts.StakeUSD}
		res.Units = opts.StakeUSD / buy.Price
		key := buy.Wallet + "|" + buy.AssetID
		for _, sell := range sellsByKey[key] {
			if sell.Time.After(buy.Time) || sell.Time.Equal(buy.Time) {
				sellCopy := sell
				res.Exit = &sellCopy
				res.ExitSource = "sell"
				res.ExitPrice = sell.Price
				res.HoldHours = sell.Time.Sub(buy.Time).Hours()
				break
			}
		}
		if res.ExitSource != "" {
			fillPnL(res)
		}
		eval.Results = append(eval.Results, res)
	}
	eval.SuppressedByWallet = sortedSuppressionStats(suppressedWallets)
	eval.SuppressedByEvent = sortedSuppressionStats(suppressedEvents)
	sort.Slice(eval.PolicyViolations, func(i, j int) bool {
		return eval.PolicyViolations[i].Time.After(eval.PolicyViolations[j].Time)
	})
	return eval
}

func reportPolicyViolationReason(tr whaleTrade) string {
	if ok, reason := reportMarketDecision(tr.Market); !ok {
		return reason
	}
	if tr.Price < 0.05 || tr.Price > 0.95 {
		return "price_filtered"
	}
	return ""
}

func reportMarketDecision(q string) (bool, string) {
	text := strings.ToLower(q)
	if reportDerivativeMarketText(text) {
		return false, "derivative_filtered"
	}
	if strings.Contains(text, "tennis") ||
		strings.Contains(text, "wimbledon") ||
		strings.Contains(text, " atp") ||
		strings.Contains(text, " wta") {
		return false, "category_filtered"
	}
	if feed.IsFollowTargetMarket(feed.Market{Question: q}) {
		return true, ""
	}
	basketball := []string{"nba", "wnba", "basketball"}
	soccer := []string{
		"epl", "premier league", "la liga", "bundesliga", "serie a", "ligue 1",
		"champions league", "ucl", "uefa", "fifa", "fifwc", "fifa world cup",
		"copa ", "concacaf", "conmebol", "eredivisie", "liga mx", "mls",
		"soccer", "fútbol", "futbol",
	}
	esports := []string{
		"lol", "league of legends", "lck", "lpl", "msi", "worlds",
		"dota", "dota2", "cs2", "csgo", "valorant", "esport",
	}
	for _, group := range [][]string{basketball, soccer, esports} {
		for _, k := range group {
			if strings.Contains(text, k) {
				return true, ""
			}
		}
	}
	return false, "category_filtered"
}

func reportDerivativeMarketText(text string) bool {
	return feed.IsDerivativeFollowMarketText(text)
}

func addSuppressionStats(dst map[string]*suppressionStats, key string, tr whaleTrade, event string) {
	if key == "" {
		return
	}
	st := dst[key]
	if st == nil {
		st = &suppressionStats{
			Wallet: tr.Wallet,
			Label:  tr.Label,
			List:   tr.List,
			Tier:   tr.Tier,
			Event:  event,
		}
		dst[key] = st
	}
	switch tr.Action {
	case "cooldown":
		st.AssetCooldown++
	case "event_cooldown":
		st.EventCooldown++
	case "duplicate":
		st.Duplicate++
	case "pending_consensus":
		st.Pending++
	}
	st.Total++
	st.NotionalUSD += tr.Size
	if tr.Time.After(st.LastTime) {
		st.LastTime = tr.Time
	}
}

func sortedSuppressionStats(src map[string]*suppressionStats) []suppressionStats {
	out := make([]suppressionStats, 0, len(src))
	for _, st := range src {
		out = append(out, *st)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total != out[j].Total {
			return out[i].Total > out[j].Total
		}
		if out[i].NotionalUSD != out[j].NotionalUSD {
			return out[i].NotionalUSD > out[j].NotionalUSD
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
	return out
}

func suppressingAction(action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "cooldown", "event_cooldown", "duplicate", "pending_consensus":
		return true
	default:
		return false
	}
}

func signalKey(tr whaleTrade) string {
	if tr.TradeID != "" {
		return tr.Wallet + "|" + tr.AssetID + "|" + tr.Side + "|" + tr.TradeID
	}
	return fmt.Sprintf("%s|%s|%s|%d|%.8f|%.8f", tr.Wallet, tr.AssetID, tr.Side, tr.Time.Unix(), tr.Price, tr.Units)
}

func actionAllowed(actions map[string]struct{}, action string) bool {
	if _, ok := actions["all"]; ok {
		return true
	}
	_, ok := actions[action]
	return ok
}

func markOpenSignals(ctx context.Context, results []*signalResult, minAge time.Duration) {
	settleClosedSignals(ctx, results, minAge)
	assets := map[string]struct{}{}
	now := time.Now()
	for _, res := range results {
		if res.ExitSource != "" {
			continue
		}
		if minAge > 0 && now.Sub(res.Buy.Time) < minAge {
			continue
		}
		assets[res.Buy.AssetID] = struct{}{}
	}
	if len(assets) == 0 {
		return
	}
	priceByAsset := fetchMidpoints(ctx, assets)
	for _, res := range results {
		if res.ExitSource != "" {
			continue
		}
		px := priceByAsset[res.Buy.AssetID]
		if px <= 0 {
			continue
		}
		res.ExitSource = "mark"
		res.ExitPrice = px
		res.HoldHours = now.Sub(res.Buy.Time).Hours()
		fillPnL(res)
	}
}

func settleClosedSignals(ctx context.Context, results []*signalResult, minAge time.Duration) {
	conditionIDs := map[string]struct{}{}
	assetIDs := map[string]struct{}{}
	now := time.Now()
	for _, res := range results {
		if res.ExitSource != "" {
			continue
		}
		if minAge > 0 && now.Sub(res.Buy.Time) < minAge {
			continue
		}
		if res.Buy.ConditionID != "" {
			conditionIDs[res.Buy.ConditionID] = struct{}{}
		} else if res.Buy.AssetID != "" {
			assetIDs[res.Buy.AssetID] = struct{}{}
		}
	}
	if len(conditionIDs) == 0 && len(assetIDs) == 0 {
		return
	}
	conditionList := make([]string, 0, len(conditionIDs))
	for id := range conditionIDs {
		conditionList = append(conditionList, id)
	}
	client := feed.NewGammaClient()
	markets := make([]feed.Market, 0, len(conditionIDs)+len(assetIDs))
	if len(conditionList) > 0 {
		byCondition, err := client.GetByConditionIDs(ctx, conditionList)
		if err == nil {
			markets = append(markets, byCondition...)
		}
	}
	if len(assetIDs) > 0 {
		assetList := make([]string, 0, len(assetIDs))
		for id := range assetIDs {
			assetList = append(assetList, id)
		}
		byAsset, err := client.GetByClobTokenIDs(ctx, assetList)
		if err == nil {
			markets = append(markets, byAsset...)
		}
	}
	byCond := map[string]feed.Market{}
	byAsset := map[string]feed.Market{}
	for _, m := range markets {
		if m.ConditionID != "" {
			byCond[m.ConditionID] = m
		}
		for _, token := range m.ClobTokenIDs() {
			byAsset[token] = m
		}
	}
	for _, res := range results {
		if res.ExitSource != "" {
			continue
		}
		var (
			m  feed.Market
			ok bool
		)
		if res.Buy.ConditionID != "" {
			m, ok = byCond[res.Buy.ConditionID]
		}
		if !ok && res.Buy.AssetID != "" {
			m, ok = byAsset[res.Buy.AssetID]
		}
		if !ok || !m.Closed {
			continue
		}
		px, ok := settlementPriceForAsset(m, res.Buy.AssetID)
		if !ok {
			continue
		}
		res.ExitSource = "settlement"
		res.ExitPrice = px
		res.HoldHours = now.Sub(res.Buy.Time).Hours()
		fillPnL(res)
	}
}

func settlementPriceForAsset(m feed.Market, asset string) (float64, bool) {
	tokens := m.ClobTokenIDs()
	prices := m.OutcomePrices()
	for i, token := range tokens {
		if token != asset || i >= len(prices) {
			continue
		}
		px, err := strconv.ParseFloat(prices[i], 64)
		if err != nil {
			return 0, false
		}
		return px, true
	}
	return 0, false
}

func fetchMidpoints(ctx context.Context, assets map[string]struct{}) map[string]float64 {
	out := map[string]float64{}
	var mu sync.Mutex
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 10 * time.Second}
	for asset := range assets {
		asset := asset
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}
			px := fetchMidpoint(ctx, client, asset)
			if px <= 0 {
				return
			}
			mu.Lock()
			out[asset] = px
			mu.Unlock()
		}()
	}
	wg.Wait()
	return out
}

func fetchMidpoint(ctx context.Context, client *http.Client, asset string) float64 {
	reqURL := "https://clob.polymarket.com/midpoint?token_id=" + url.QueryEscape(asset)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		return 0
	}
	var body struct {
		Mid string `json:"mid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0
	}
	px, _ := strconv.ParseFloat(body.Mid, 64)
	return px
}

func fillPnL(res *signalResult) {
	res.PnLUSD = res.Units * (res.ExitPrice - res.Buy.Price)
	if res.StakeUSD > 0 {
		res.ReturnPct = res.PnLUSD / res.StakeUSD * 100
	}
}

func summarize(results []*signalResult) walletStats {
	var s walletStats
	for _, r := range results {
		if r.ExitSource == "" {
			s.Open++
			continue
		}
		s.Signals++
		countExitSource(&s, r.ExitSource)
		if r.PnLUSD > 0 {
			s.Wins++
		}
		s.PnLUSD += r.PnLUSD
		s.StakeUSD += r.StakeUSD
		s.AvgReturnPct += r.ReturnPct
	}
	if s.StakeUSD > 0 {
		s.ReturnPct = s.PnLUSD / s.StakeUSD * 100
	}
	if s.Signals > 0 {
		s.AvgReturnPct /= float64(s.Signals)
	}
	return s
}

func isProvenExitSource(source string) bool {
	return source == "sell" || source == "settlement"
}

func provenResults(results []*signalResult) []*signalResult {
	out := make([]*signalResult, 0, len(results))
	for _, r := range results {
		if isProvenExitSource(r.ExitSource) {
			out = append(out, r)
		}
	}
	return out
}

func statsByWallet(results []*signalResult) []walletStats {
	byWallet := map[string]*walletStats{}
	for _, r := range results {
		if r.ExitSource == "" {
			continue
		}
		w := r.Buy.Wallet
		st := byWallet[w]
		if st == nil {
			st = &walletStats{Wallet: w, Label: r.Buy.Label, List: r.Buy.List, Tier: r.Buy.Tier}
			byWallet[w] = st
		}
		st.Signals++
		countExitSource(st, r.ExitSource)
		if r.PnLUSD > 0 {
			st.Wins++
		}
		st.PnLUSD += r.PnLUSD
		st.StakeUSD += r.StakeUSD
		st.AvgReturnPct += r.ReturnPct
	}
	out := make([]walletStats, 0, len(byWallet))
	for _, st := range byWallet {
		if st.StakeUSD > 0 {
			st.ReturnPct = st.PnLUSD / st.StakeUSD * 100
		}
		if st.Signals > 0 {
			st.AvgReturnPct /= float64(st.Signals)
		}
		out = append(out, *st)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ReturnPct != out[j].ReturnPct {
			return out[i].ReturnPct > out[j].ReturnPct
		}
		return out[i].Signals > out[j].Signals
	})
	return out
}

func statsByList(results []*signalResult) []walletStats {
	byList := map[string]*walletStats{}
	for _, r := range results {
		if r.ExitSource == "" {
			continue
		}
		list := r.Buy.List
		if list == "" {
			list = "unknown"
		}
		st := byList[list]
		if st == nil {
			st = &walletStats{List: list}
			byList[list] = st
		}
		st.Signals++
		countExitSource(st, r.ExitSource)
		if r.PnLUSD > 0 {
			st.Wins++
		}
		st.PnLUSD += r.PnLUSD
		st.StakeUSD += r.StakeUSD
		st.AvgReturnPct += r.ReturnPct
	}
	out := make([]walletStats, 0, len(byList))
	for _, st := range byList {
		if st.StakeUSD > 0 {
			st.ReturnPct = st.PnLUSD / st.StakeUSD * 100
		}
		if st.Signals > 0 {
			st.AvgReturnPct /= float64(st.Signals)
		}
		out = append(out, *st)
	}
	sort.Slice(out, func(i, j int) bool {
		order := map[string]int{
			"core": 0, "watch": 1, "sports": 2, "target": 3,
			"scout": 4, "flow": 5, "unknown": 6,
		}
		if order[out[i].List] != order[out[j].List] {
			return order[out[i].List] < order[out[j].List]
		}
		return out[i].Signals > out[j].Signals
	})
	return out
}

func statsByEvent(results []*signalResult) []eventStats {
	byEvent := map[string]*eventStats{}
	marketsByEvent := map[string]map[string]struct{}{}
	for _, r := range results {
		if r.ExitSource == "" {
			continue
		}
		event := eventKey(r.Buy.Market)
		key := r.Buy.Wallet + "|" + event
		st := byEvent[key]
		if st == nil {
			st = &eventStats{
				Event:     event,
				Wallet:    r.Buy.Wallet,
				Label:     r.Buy.Label,
				List:      r.Buy.List,
				Tier:      r.Buy.Tier,
				FirstTime: r.Buy.Time,
				LastTime:  r.Buy.Time,
			}
			byEvent[key] = st
			marketsByEvent[key] = map[string]struct{}{}
		}
		st.Signals++
		countEventExitSource(st, r.ExitSource)
		if r.PnLUSD > 0 {
			st.Wins++
		}
		st.PnLUSD += r.PnLUSD
		st.StakeUSD += r.StakeUSD
		st.AvgReturnPct += r.ReturnPct
		if r.Buy.Time.Before(st.FirstTime) {
			st.FirstTime = r.Buy.Time
		}
		if r.Buy.Time.After(st.LastTime) {
			st.LastTime = r.Buy.Time
		}
		marketsByEvent[key][r.Buy.Market] = struct{}{}
	}
	out := make([]eventStats, 0, len(byEvent))
	for key, st := range byEvent {
		st.Markets = len(marketsByEvent[key])
		if st.StakeUSD > 0 {
			st.ReturnPct = st.PnLUSD / st.StakeUSD * 100
		}
		if st.Signals > 0 {
			st.AvgReturnPct /= float64(st.Signals)
		}
		out = append(out, *st)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].LastTime.Equal(out[j].LastTime) {
			return out[i].LastTime.After(out[j].LastTime)
		}
		return out[i].PnLUSD > out[j].PnLUSD
	})
	return out
}

func countExitSource(st *walletStats, source string) {
	switch source {
	case "sell":
		st.Closed++
	case "settlement":
		st.Settled++
	case "mark":
		st.Marked++
	}
}

func countEventExitSource(st *eventStats, source string) {
	switch source {
	case "sell":
		st.Closed++
	case "settlement":
		st.Settled++
	case "mark":
		st.Marked++
	}
}

func summarizeEvents(events []eventStats) walletStats {
	var s walletStats
	for _, ev := range events {
		s.Signals++
		s.Closed += ev.Closed
		s.Settled += ev.Settled
		s.Marked += ev.Marked
		if ev.PnLUSD > 0 {
			s.Wins++
		}
		s.PnLUSD += ev.PnLUSD
		s.StakeUSD += ev.StakeUSD
		s.AvgReturnPct += ev.ReturnPct
	}
	if s.StakeUSD > 0 {
		s.ReturnPct = s.PnLUSD / s.StakeUSD * 100
	}
	if s.Signals > 0 {
		s.AvgReturnPct /= float64(s.Signals)
	}
	return s
}

func eventCappedResults(results []*signalResult) []*signalResult {
	eligible := make([]*signalResult, 0, len(results))
	for _, r := range results {
		if r.ExitSource != "" {
			eligible = append(eligible, r)
		}
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].Buy.Time.Before(eligible[j].Buy.Time) })
	seen := map[string]struct{}{}
	out := make([]*signalResult, 0, len(eligible))
	for _, r := range eligible {
		key := r.Buy.Wallet + "|" + eventKey(r.Buy.Market)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, r)
	}
	return out
}

func countSource(results []*signalResult, source string) int {
	n := 0
	for _, r := range results {
		if r.ExitSource == source {
			n++
		}
	}
	return n
}

func writeReport(path, logPath, walletsPath string, eval evaluation, opts evalOptions) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	results := eval.Results
	sum := summarize(results)
	proven := provenResults(results)
	provenSum := summarize(proven)
	events := statsByEvent(results)
	eventSum := summarizeEvents(events)
	capped := eventCappedResults(results)
	cappedSum := summarize(capped)
	cappedProven := eventCappedResults(proven)
	cappedProvenSum := summarize(cappedProven)
	var b strings.Builder
	fmt.Fprintf(&b, "# Whale Performance Report\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Log: `%s`\n", logPath)
	if walletsPath != "" {
		fmt.Fprintf(&b, "- Wallet filter: `%s`\n", walletsPath)
	}
	fmt.Fprintf(&b, "- Fixed stake: $%.2f per BUY signal\n", opts.StakeUSD)
	fmt.Fprintf(&b, "- Minimum whale notional: $%.0f\n", opts.MinNotional)
	if len(opts.ListMinNotional) > 0 {
		fmt.Fprintf(&b, "- List minimum notionals: %s\n", formatListMinNotional(opts.ListMinNotional))
	}
	if !opts.Since.IsZero() {
		fmt.Fprintf(&b, "- Policy since: %s\n", opts.Since.Format(time.RFC3339))
	}
	if opts.RepeatCooldown > 0 {
		fmt.Fprintf(&b, "- Repeat cooldown: %s per wallet+asset; BUYs >= $%.0f bypass cooldown\n", opts.RepeatCooldown, opts.RepeatMinNotional)
	}
	fmt.Fprintf(&b, "\n")
	fmt.Fprintf(&b, "## Summary\n\n")
	fmt.Fprintf(&b, "- Raw matched BUY signals: %d\n", eval.RawBuys)
	if opts.RepeatCooldown > 0 {
		fmt.Fprintf(&b, "- Suppressed repeat BUYs: %d\n", eval.SuppressedRepeats)
	}
	fmt.Fprintf(&b, "- Logged asset-cooldown BUYs: %d\n", eval.LoggedCooldownBuys)
	fmt.Fprintf(&b, "- Logged event-cooldown BUYs: %d\n", eval.LoggedEventCooldowns)
	fmt.Fprintf(&b, "- Logged pending-consensus BUYs: %d\n", eval.LoggedPendingBuys)
	fmt.Fprintf(&b, "- Logged duplicate BUYs: %d\n", eval.LoggedDuplicateBuys)
	fmt.Fprintf(&b, "- Duplicate BUY alerts ignored: %d\n", eval.DuplicateBuys)
	fmt.Fprintf(&b, "- Evaluated signals: %d\n", sum.Signals)
	fmt.Fprintf(&b, "- Realized via whale SELL: %d\n", sum.Closed)
	fmt.Fprintf(&b, "- Settled by market resolution: %d\n", sum.Settled)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", sum.Marked)
	fmt.Fprintf(&b, "- Still open/unmarked: %d\n", sum.Open)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(sum.Wins), float64(sum.Signals)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", sum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", sum.ReturnPct)
	fmt.Fprintf(&b, "- Proven signals: %d\n", provenSum.Signals)
	fmt.Fprintf(&b, "- Proven win rate: %.1f%%\n", pct(float64(provenSum.Wins), float64(provenSum.Signals)))
	fmt.Fprintf(&b, "- Proven PnL: $%+.2f\n", provenSum.PnLUSD)
	fmt.Fprintf(&b, "- Proven ROI: %.1f%%\n\n", provenSum.ReturnPct)
	writePolicyViolationSection(&b, eval.PolicyViolations)
	writeSuppressionSection(&b, eval)
	fmt.Fprintf(&b, "## Event Cluster Summary\n\n")
	fmt.Fprintf(&b, "- Independent wallet-event clusters: %d\n", eventSum.Signals)
	fmt.Fprintf(&b, "- Event-cluster win rate: %.1f%%\n", pct(float64(eventSum.Wins), float64(eventSum.Signals)))
	fmt.Fprintf(&b, "- Event-cluster ROI: %.1f%%\n", eventSum.ReturnPct)
	fmt.Fprintf(&b, "- Event-cluster PnL: $%+.2f\n\n", eventSum.PnLUSD)
	fmt.Fprintf(&b, "## Event-Capped Strategy\n\n")
	fmt.Fprintf(&b, "- Rule: one fixed-stake entry per wallet-event cluster, using the first evaluated signal\n")
	fmt.Fprintf(&b, "- Entries: %d\n", cappedSum.Signals)
	fmt.Fprintf(&b, "- Realized via whale SELL: %d\n", cappedSum.Closed)
	fmt.Fprintf(&b, "- Settled by market resolution: %d\n", cappedSum.Settled)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", cappedSum.Marked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(cappedSum.Wins), float64(cappedSum.Signals)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", cappedSum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", cappedSum.ReturnPct)
	fmt.Fprintf(&b, "- Proven entries: %d\n", cappedProvenSum.Signals)
	fmt.Fprintf(&b, "- Proven win rate: %.1f%%\n", pct(float64(cappedProvenSum.Wins), float64(cappedProvenSum.Signals)))
	fmt.Fprintf(&b, "- Proven PnL: $%+.2f\n", cappedProvenSum.PnLUSD)
	fmt.Fprintf(&b, "- Proven ROI: %.1f%%\n\n", cappedProvenSum.ReturnPct)

	fmt.Fprintf(&b, "## Event-Capped By List\n\n")
	fmt.Fprintf(&b, "| List | Entries | Closed | Settled | Marked | Win | ROI | PnL |\n")
	fmt.Fprintf(&b, "|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, st := range statsByList(capped) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %.1f%% | %.1f%% | $%+.2f |\n",
			st.List, st.Signals, st.Closed, st.Settled, st.Marked,
			pct(float64(st.Wins), float64(st.Signals)), st.ReturnPct, st.PnLUSD)
	}

	fmt.Fprintf(&b, "\n## Event-Capped By Wallet\n\n")
	fmt.Fprintf(&b, "| Wallet | List | Tier | Label | Entries | Closed | Settled | Marked | Win | ROI | PnL |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, st := range statsByWallet(capped) {
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %d | %d | %d | %d | %.1f%% | %.1f%% | $%+.2f |\n",
			shortAddr(st.Wallet), st.List, st.Tier, st.Label, st.Signals, st.Closed, st.Settled, st.Marked,
			pct(float64(st.Wins), float64(st.Signals)), st.ReturnPct, st.PnLUSD)
	}

	fmt.Fprintf(&b, "## By List\n\n")
	fmt.Fprintf(&b, "| List | Signals | Closed | Settled | Marked | Win | ROI | PnL |\n")
	fmt.Fprintf(&b, "|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, st := range statsByList(results) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %.1f%% | %.1f%% | $%+.2f |\n",
			st.List, st.Signals, st.Closed, st.Settled, st.Marked,
			pct(float64(st.Wins), float64(st.Signals)), st.ReturnPct, st.PnLUSD)
	}

	fmt.Fprintf(&b, "## By Wallet\n\n")
	fmt.Fprintf(&b, "| Wallet | List | Tier | Label | Signals | Closed | Settled | Marked | Win | ROI | PnL |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, st := range statsByWallet(results) {
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %d | %d | %d | %d | %.1f%% | %.1f%% | $%+.2f |\n",
			shortAddr(st.Wallet), st.List, st.Tier, st.Label, st.Signals, st.Closed, st.Settled, st.Marked,
			pct(float64(st.Wins), float64(st.Signals)), st.ReturnPct, st.PnLUSD)
	}

	fmt.Fprintf(&b, "\n## By Event Cluster\n\n")
	fmt.Fprintf(&b, "| Event | Wallet | List | Tier | Signals | Markets | Closed | Settled | Marked | Win | ROI | PnL |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	eventLimit := 30
	if len(events) < eventLimit {
		eventLimit = len(events)
	}
	for i := 0; i < eventLimit; i++ {
		st := events[i]
		fmt.Fprintf(&b, "| %s | `%s` | %s | %s | %d | %d | %d | %d | %d | %.1f%% | %.1f%% | $%+.2f |\n",
			truncate(st.Event, 46), shortAddr(st.Wallet), st.List, st.Tier, st.Signals, st.Markets,
			st.Closed, st.Settled, st.Marked, pct(float64(st.Wins), float64(st.Signals)), st.ReturnPct, st.PnLUSD)
	}

	fmt.Fprintf(&b, "\n## Event-Capped Entries\n\n")
	fmt.Fprintf(&b, "| Time | Event | Wallet | List | Outcome | Entry | Exit | Src | Ret | PnL | Market |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---|---:|---:|---|---:|---:|---|\n")
	cappedRecent := append([]*signalResult{}, capped...)
	sort.Slice(cappedRecent, func(i, j int) bool { return cappedRecent[i].Buy.Time.After(cappedRecent[j].Buy.Time) })
	cappedLimit := 30
	if len(cappedRecent) < cappedLimit {
		cappedLimit = len(cappedRecent)
	}
	for i := 0; i < cappedLimit; i++ {
		r := cappedRecent[i]
		fmt.Fprintf(&b, "| %s | %s | `%s` | %s | %s | %.4f | %.4f | %s | %+.1f%% | $%+.2f | %s |\n",
			r.Buy.Time.Format("01-02 15:04"), truncate(eventKey(r.Buy.Market), 34), shortAddr(r.Buy.Wallet),
			r.Buy.List, truncate(r.Buy.Outcome, 18), r.Buy.Price, r.ExitPrice, r.ExitSource, r.ReturnPct,
			r.PnLUSD, truncate(r.Buy.Market, 42))
	}

	fmt.Fprintf(&b, "\n## Recent Signals\n\n")
	fmt.Fprintf(&b, "| Time | Wallet | List | Tier | Side | Outcome | Entry | Exit | Src | Ret | PnL | Market |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---|---|---:|---:|---|---:|---:|---|\n")
	recent := make([]*signalResult, 0, len(results))
	for _, r := range results {
		if r.ExitSource != "" {
			recent = append(recent, r)
		}
	}
	sort.Slice(recent, func(i, j int) bool { return recent[i].Buy.Time.After(recent[j].Buy.Time) })
	limit := 40
	if len(recent) < limit {
		limit = len(recent)
	}
	for i := 0; i < limit; i++ {
		r := recent[i]
		fmt.Fprintf(&b, "| %s | `%s` | %s | %s | BUY | %s | %.4f | %.4f | %s | %+.1f%% | $%+.2f | %s |\n",
			r.Buy.Time.Format("01-02 15:04"), shortAddr(r.Buy.Wallet), r.Buy.List, r.Buy.Tier, truncate(r.Buy.Outcome, 18),
			r.Buy.Price, r.ExitPrice, r.ExitSource, r.ReturnPct, r.PnLUSD, truncate(r.Buy.Market, 48))
	}

	if sum.Signals == 0 {
		fmt.Fprintf(&b, "\nNo evaluated BUY signals matched the current filters yet. Keep the bot running; this report will become useful after the core list fires.\n")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func buildSnapshot(logPath, walletsPath, reportPath string, eval evaluation, opts evalOptions) performanceSnapshot {
	sum := summarize(eval.Results)
	proven := provenResults(eval.Results)
	provenSum := summarize(proven)
	capped := eventCappedResults(eval.Results)
	cappedSum := summarize(capped)
	cappedProven := eventCappedResults(proven)
	cappedProvenSum := summarize(cappedProven)
	snap := performanceSnapshot{
		GeneratedAt:                 time.Now().Format(time.RFC3339),
		LogPath:                     logPath,
		WalletsPath:                 walletsPath,
		ReportPath:                  reportPath,
		StakeUSD:                    opts.StakeUSD,
		MinNotional:                 opts.MinNotional,
		ListMinNotional:             opts.ListMinNotional,
		RawBuys:                     eval.RawBuys,
		SuppressedRepeats:           eval.SuppressedRepeats,
		LoggedCooldownBuys:          eval.LoggedCooldownBuys,
		LoggedEventCooldowns:        eval.LoggedEventCooldowns,
		LoggedDuplicateBuys:         eval.LoggedDuplicateBuys,
		LoggedPendingBuys:           eval.LoggedPendingBuys,
		DuplicateBuys:               eval.DuplicateBuys,
		EvaluatedSignals:            sum.Signals,
		Closed:                      sum.Closed,
		Settled:                     sum.Settled,
		Marked:                      sum.Marked,
		Open:                        sum.Open,
		Wins:                        sum.Wins,
		PnLUSD:                      sum.PnLUSD,
		StakeEvaluatedUSD:           sum.StakeUSD,
		ReturnPct:                   sum.ReturnPct,
		WinRatePct:                  pct(float64(sum.Wins), float64(sum.Signals)),
		ProvenSignals:               provenSum.Signals,
		ProvenWins:                  provenSum.Wins,
		ProvenPnLUSD:                provenSum.PnLUSD,
		ProvenStakeUSD:              provenSum.StakeUSD,
		ProvenReturnPct:             provenSum.ReturnPct,
		ProvenWinRatePct:            pct(float64(provenSum.Wins), float64(provenSum.Signals)),
		EventCappedSignals:          cappedSum.Signals,
		EventCappedWins:             cappedSum.Wins,
		EventCappedPnLUSD:           cappedSum.PnLUSD,
		EventCappedStakeUSD:         cappedSum.StakeUSD,
		EventCappedReturnPct:        cappedSum.ReturnPct,
		EventCappedWinRatePct:       pct(float64(cappedSum.Wins), float64(cappedSum.Signals)),
		EventCappedProvenSignals:    cappedProvenSum.Signals,
		EventCappedProvenWins:       cappedProvenSum.Wins,
		EventCappedProvenPnLUSD:     cappedProvenSum.PnLUSD,
		EventCappedProvenStakeUSD:   cappedProvenSum.StakeUSD,
		EventCappedProvenReturnPct:  cappedProvenSum.ReturnPct,
		EventCappedProvenWinRatePct: pct(float64(cappedProvenSum.Wins), float64(cappedProvenSum.Signals)),
		EventCappedByList:           snapshotStatsList(statsByList(capped)),
		EventCappedByWallet:         snapshotStatsList(statsByWallet(capped)),
		EventCappedProvenByWallet:   snapshotStatsList(statsByWallet(cappedProven)),
		ByList:                      snapshotStatsList(statsByList(eval.Results)),
		ByWallet:                    snapshotStatsList(statsByWallet(eval.Results)),
		ByEvent:                     snapshotEventStatsList(statsByEvent(eval.Results)),
		SuppressedByWallet:          snapshotSuppressionStatsList(eval.SuppressedByWallet),
		SuppressedByEvent:           snapshotSuppressionStatsList(eval.SuppressedByEvent),
		PolicyViolationCount:        len(eval.PolicyViolations),
		PolicyViolations:            snapshotPolicyViolationList(eval.PolicyViolations),
	}
	if !opts.Since.IsZero() {
		snap.Since = opts.Since.Format(time.RFC3339)
	}
	if opts.RepeatCooldown > 0 {
		snap.RepeatCooldown = opts.RepeatCooldown.String()
		snap.RepeatMinNotional = opts.RepeatMinNotional
	}
	return snap
}

func writePolicyViolationSection(b *strings.Builder, violations []policyViolation) {
	fmt.Fprintf(b, "## Policy Violations\n\n")
	fmt.Fprintf(b, "- Alerted BUYs outside current sports/esports/price policy: %d\n\n", len(violations))
	if len(violations) == 0 {
		return
	}
	fmt.Fprintf(b, "| Time | Reason | Wallet | List | Tier | Notional | Price | Outcome | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---|---|---:|---:|---|---|\n")
	limit := minInt(len(violations), 25)
	for i := 0; i < limit; i++ {
		v := violations[i]
		fmt.Fprintf(b, "| %s | %s | `%s` | %s | %s | $%.0f | %.4f | %s | %s |\n",
			formatShortTime(v.Time), v.Reason, shortAddr(v.Wallet), v.List, v.Tier,
			v.NotionalUSD, v.Price, truncate(v.Outcome, 18), truncate(v.Market, 56))
	}
	fmt.Fprintf(b, "\n")
}

func writeSuppressionSection(b *strings.Builder, eval evaluation) {
	if len(eval.SuppressedByWallet) == 0 && len(eval.SuppressedByEvent) == 0 {
		return
	}
	fmt.Fprintf(b, "## Suppressed BUY Noise\n\n")
	fmt.Fprintf(b, "These rows are detected large BUYs that were logged but not pushed/evaluated because duplicate, cooldown, or consensus gates suppressed them.\n\n")
	if len(eval.SuppressedByWallet) > 0 {
		fmt.Fprintf(b, "### By Wallet\n\n")
		fmt.Fprintf(b, "| Wallet | List | Tier | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |\n")
		fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---:|---:|---|\n")
		limit := minInt(len(eval.SuppressedByWallet), 12)
		for i := 0; i < limit; i++ {
			st := eval.SuppressedByWallet[i]
			fmt.Fprintf(b, "| `%s` | %s | %s | %d | %d | %d | %d | %d | $%.0f | %s |\n",
				shortAddr(st.Wallet), st.List, st.Tier, st.AssetCooldown, st.EventCooldown, st.Pending, st.Duplicate, st.Total, st.NotionalUSD, formatShortTime(st.LastTime))
		}
		fmt.Fprintf(b, "\n")
	}
	if len(eval.SuppressedByEvent) > 0 {
		fmt.Fprintf(b, "### By Wallet-Event\n\n")
		fmt.Fprintf(b, "| Event | Wallet | List | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |\n")
		fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---:|---:|---|\n")
		limit := minInt(len(eval.SuppressedByEvent), 12)
		for i := 0; i < limit; i++ {
			st := eval.SuppressedByEvent[i]
			fmt.Fprintf(b, "| %s | `%s` | %s | %d | %d | %d | %d | %d | $%.0f | %s |\n",
				truncate(st.Event, 42), shortAddr(st.Wallet), st.List, st.AssetCooldown, st.EventCooldown, st.Pending, st.Duplicate, st.Total, st.NotionalUSD, formatShortTime(st.LastTime))
		}
		fmt.Fprintf(b, "\n")
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func formatShortTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("01-02 15:04")
}

func formatListMinNotional(listMins map[string]float64) string {
	keys := make([]string, 0, len(listMins))
	for k := range listMins {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=$%.0f", k, listMins[k]))
	}
	return strings.Join(parts, ", ")
}

func snapshotStatsList(in []walletStats) []snapshotStats {
	out := make([]snapshotStats, 0, len(in))
	for _, st := range in {
		out = append(out, snapshotStats(st))
	}
	return out
}

func snapshotEventStatsList(in []eventStats) []snapshotEventStats {
	out := make([]snapshotEventStats, 0, len(in))
	for _, st := range in {
		ev := snapshotEventStats{
			Event:        st.Event,
			Wallet:       st.Wallet,
			Label:        st.Label,
			List:         st.List,
			Tier:         st.Tier,
			Signals:      st.Signals,
			Closed:       st.Closed,
			Settled:      st.Settled,
			Marked:       st.Marked,
			Wins:         st.Wins,
			Markets:      st.Markets,
			PnLUSD:       st.PnLUSD,
			StakeUSD:     st.StakeUSD,
			ReturnPct:    st.ReturnPct,
			AvgReturnPct: st.AvgReturnPct,
		}
		if !st.FirstTime.IsZero() {
			ev.FirstTime = st.FirstTime.Format(time.RFC3339)
		}
		if !st.LastTime.IsZero() {
			ev.LastTime = st.LastTime.Format(time.RFC3339)
		}
		out = append(out, ev)
	}
	return out
}

func snapshotSuppressionStatsList(in []suppressionStats) []snapshotSuppressionStats {
	out := make([]snapshotSuppressionStats, 0, len(in))
	for _, st := range in {
		ss := snapshotSuppressionStats{
			Wallet:        st.Wallet,
			Label:         st.Label,
			List:          st.List,
			Tier:          st.Tier,
			Event:         st.Event,
			AssetCooldown: st.AssetCooldown,
			EventCooldown: st.EventCooldown,
			Duplicate:     st.Duplicate,
			Pending:       st.Pending,
			Total:         st.Total,
			NotionalUSD:   st.NotionalUSD,
		}
		if !st.LastTime.IsZero() {
			ss.LastTime = st.LastTime.Format(time.RFC3339)
		}
		out = append(out, ss)
	}
	return out
}

func snapshotPolicyViolationList(in []policyViolation) []snapshotPolicyViolation {
	limit := minInt(len(in), 50)
	out := make([]snapshotPolicyViolation, 0, limit)
	for i := 0; i < limit; i++ {
		v := in[i]
		out = append(out, snapshotPolicyViolation{
			Time:        v.Time.Format(time.RFC3339),
			Wallet:      v.Wallet,
			Label:       v.Label,
			List:        v.List,
			Tier:        v.Tier,
			Reason:      v.Reason,
			Market:      v.Market,
			Outcome:     v.Outcome,
			Price:       v.Price,
			NotionalUSD: v.NotionalUSD,
			TradeID:     v.TradeID,
		})
	}
	return out
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func appendJSONL(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

func pct(num, den float64) float64 {
	if den == 0 {
		return 0
	}
	return num / den * 100
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}

var (
	parentheticalRE = regexp.MustCompile(`\s*\([^)]*\b(?:bo[0-9]+|game|map)\b[^)]*\)`)
	gameSuffixRE    = regexp.MustCompile(`\s*-\s*(?:game|map)\s*[0-9]+\s+winner\b.*$`)
	spaceRE         = regexp.MustCompile(`\s+`)
)

func eventKey(market string) string {
	s := strings.ToLower(strings.TrimSpace(market))
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "—", "-")
	s = gameSuffixRE.ReplaceAllString(s, "")
	if idx := strings.Index(s, " - "); idx >= 0 {
		s = s[:idx]
	}
	if vs := strings.Index(s, " vs"); vs >= 0 {
		if colon := strings.Index(s[vs:], ":"); colon >= 0 {
			s = s[:vs+colon]
		}
	}
	s = parentheticalRE.ReplaceAllString(s, "")
	s = strings.TrimSuffix(s, " winner")
	s = strings.TrimSuffix(s, " match winner")
	s = strings.TrimSpace(strings.Trim(s, "-:"))
	s = spaceRE.ReplaceAllString(s, " ")
	if s == "" {
		return "unknown"
	}
	return s
}
