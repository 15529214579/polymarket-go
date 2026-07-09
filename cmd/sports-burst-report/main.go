package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tapeTrade struct {
	Time        time.Time `json:"time"`
	Wallet      string    `json:"wallet"`
	Side        string    `json:"side"`
	Notional    float64   `json:"notional"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
	Outcome     string    `json:"outcome"`
	Market      string    `json:"market"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	Asset       string    `json:"asset"`
	Transaction string    `json:"transaction"`
	KnownList   string    `json:"known_list,omitempty"`
	Tier        string    `json:"tier,omitempty"`
	Bot         float64   `json:"bot,omitempty"`
}

type walletStatus struct {
	List string
	Tier string
	Bot  float64
}

type edgeSnapshot struct {
	Wallet     string  `json:"wallet"`
	HorizonSec int64   `json:"horizon_sec"`
	DeltaPP    float64 `json:"delta_pp"`
}

type edgeStats struct {
	Samples int
	Wins    int
	SumPP   float64
}

type participantMeta struct {
	Tier string  `json:"tier,omitempty"`
	Bot  float64 `json:"bot,omitempty"`
}

type burstSignal struct {
	Wallet          string
	Scope           string
	Participants    []string
	ParticipantMeta map[string]participantMeta
	Asset           string
	Slug            string
	Mode            string
	KnownList       string
	Tier            string
	Bot             float64
	Category        string
	Market          string
	Outcome         string
	Trades          int
	Wallets         int
	TotalNotional   float64
	TotalSize       float64
	VWAP            float64
	LastNotional    float64
	FirstTime       time.Time
	LastTime        time.Time
}

type burstResult struct {
	Signal   burstSignal
	StakeUSD float64
	Units    float64
	Mid      float64
	Marked   bool
	PnLUSD   float64
	ReturnPC float64
}

type consensusAlertEvent struct {
	Key          string    `json:"key"`
	TradeTime    time.Time `json:"trade_time"`
	Mode         string    `json:"mode"`
	Wallet       string    `json:"wallet"`
	KnownList    string    `json:"known_list,omitempty"`
	Tier         string    `json:"tier,omitempty"`
	Bot          float64   `json:"bot,omitempty"`
	Category     string    `json:"category"`
	Notional     float64   `json:"notional"`
	Price        float64   `json:"price"`
	Outcome      string    `json:"outcome"`
	Market       string    `json:"market"`
	Slug         string    `json:"slug,omitempty"`
	Asset        string    `json:"asset"`
	Transaction  string    `json:"transaction,omitempty"`
	Participants []string  `json:"participants,omitempty"`
}

type summary struct {
	Signals   int
	Marked    int
	Unmarked  int
	Wins      int
	StakeUSD  float64
	PnLUSD    float64
	ReturnPC  float64
	AvgDeltaP float64
}

type participantSummary struct {
	Wallet        string
	Tier          string
	Bot           float64
	Signals       int
	Marked        int
	Wins          int
	TotalNotional float64
	PnLUSD        float64
	StakeUSD      float64
	ReturnPC      float64
	AvgDeltaP     float64
}

type consensusEvent struct {
	Key             string                     `json:"key"`
	Mode            string                     `json:"mode,omitempty"`
	FirstTime       time.Time                  `json:"first_time"`
	LastTime        time.Time                  `json:"last_time"`
	Asset           string                     `json:"asset"`
	Slug            string                     `json:"slug,omitempty"`
	Category        string                     `json:"category,omitempty"`
	Market          string                     `json:"market"`
	Outcome         string                     `json:"outcome"`
	Wallets         int                        `json:"wallets"`
	Trades          int                        `json:"trades"`
	TotalNotional   float64                    `json:"total_notional"`
	VWAP            float64                    `json:"vwap"`
	Marked          bool                       `json:"marked"`
	Mid             float64                    `json:"mid,omitempty"`
	PnLUSD          float64                    `json:"pnl_usd,omitempty"`
	ReturnPC        float64                    `json:"return_pc,omitempty"`
	Participants    []string                   `json:"participants"`
	ParticipantMeta map[string]participantMeta `json:"participant_meta,omitempty"`
	UpdatedAt       time.Time                  `json:"updated_at"`
}

type midpointCache struct {
	Assets map[string]cachedMidpoint `json:"assets"`
}

type cachedMidpoint struct {
	Mid       float64   `json:"mid"`
	FetchedAt time.Time `json:"fetched_at"`
}

func main() {
	tapePath := flag.String("tape", "db/strategy_iteration/sports_tape.jsonl", "sports tape JSONL")
	reportPath := flag.String("report", "reports/sports_burst_performance.md", "markdown report path")
	walletStatusPath := flag.String("wallet_statuses", "", "comma-separated wallet status files used to align burst metrics with live alert filtering")
	edgeSnapshotsPath := flag.String("edge_snapshots", "db/strategy_iteration/whale_edge_snapshots.jsonl", "optional whale-edge snapshot JSONL used to suppress consensus participants with negative measured edge")
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet score JSON used to backfill consensus participant tier/bot metadata")
	participantsOutPath := flag.String("participants_out", "", "optional research wallet file for positive CONSENSUS participants")
	participantsExcludeWalletsPath := flag.String("participants_exclude_wallets", "", "comma-separated wallet files excluded from participants_out")
	consensusEventsOutPath := flag.String("consensus_events_out", "db/strategy_iteration/sports_consensus_events.jsonl", "durable JSONL history for CONSENSUS bursts and marks")
	consensusWatchEventsOutPath := flag.String("consensus_watch_events_out", "db/strategy_iteration/sports_consensus_watch_events.jsonl", "durable JSONL history for research-only CONSENSUS-WATCH bursts")
	consensusAlertLogsPath := flag.String("consensus_alert_logs", "", "comma-separated sports tape alert JSONL logs imported into durable CONSENSUS event history")
	markCachePath := flag.String("mark_cache", "db/strategy_iteration/sports_alert_midpoints.json", "JSON cache for last known token marks; empty disables cache")
	gammaBase := flag.String("gamma_base", "https://gamma-api.polymarket.com", "Gamma API base URL used for settled-market mark fallback")
	stakeUSD := flag.Float64("stake", 10, "fixed paper stake per burst")
	window := flag.Duration("window", 15*time.Minute, "same-wallet same-asset burst window")
	maxAge := flag.Duration("max_age", 6*time.Hour, "maximum latest trade age included in report")
	minNotional := flag.Float64("min_notional", 5000, "minimum cumulative BUY notional")
	minTrades := flag.Int("min_trades", 2, "minimum BUY trades in burst")
	minLegNotional := flag.Float64("min_leg_notional", 1000, "minimum individual BUY notional included")
	consensusMinNotional := flag.Float64("consensus_min_notional", 10000, "minimum cumulative BUY notional for cross-wallet same-asset consensus bursts")
	consensusMinWallets := flag.Int("consensus_min_wallets", 2, "minimum unique wallets for cross-wallet same-asset consensus bursts")
	consensusMaxBot := flag.Float64("consensus_max_bot", 60, "maximum bot score for wallets included in consensus bursts; 0 disables bot-score filtering")
	consensusWatchMinNotional := flag.Float64("consensus_watch_min_notional", 5000, "research-only lower cumulative BUY notional for CONSENSUS-WATCH rows in the burst report; 0 disables")
	consensusWatchMinWallets := flag.Int("consensus_watch_min_wallets", 2, "research-only minimum unique wallets for CONSENSUS-WATCH rows")
	consensusHistoryMaxAge := flag.Duration("consensus_history_max_age", 0, "maximum latest trade age for durable consensus event history; 0 uses max_age")
	edgeBlock15mSamples := flag.Int("edge_block_15m_samples", 2, "minimum 15m edge samples needed to suppress a consensus participant")
	edgeBlock15mMaxAvg := flag.Float64("edge_block_15m_max_avg_pp", -1, "suppress consensus participant when 15m edge average is at or below this pp value")
	edgeBlock1hSamples := flag.Int("edge_block_1h_samples", 1, "minimum 1h edge samples needed to suppress a consensus participant")
	edgeBlock1hMaxAvg := flag.Float64("edge_block_1h_max_avg_pp", -5, "suppress consensus participant when 1h edge average is at or below this pp value")
	participantMinSignals := flag.Int("participant_min_signals", 1, "minimum positive CONSENSUS participant signals written to participants_out")
	participantMinROI := flag.Float64("participant_min_roi", 0, "minimum CONSENSUS participant ROI written to participants_out")
	timeout := flag.Duration("timeout", 20*time.Second, "overall midpoint fetch timeout")
	flag.Parse()

	trades, err := loadTapeTrades(*tapePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load tape: %v\n", err)
		os.Exit(1)
	}
	statuses, err := loadWalletStatuses(*walletStatusPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load wallet statuses: %v\n", err)
		os.Exit(1)
	}
	applyWalletStatuses(trades, statuses)
	edgeMetrics, err := loadEdgeMetrics(*edgeSnapshotsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load edge snapshots: %v\n", err)
		os.Exit(1)
	}
	edgeBlocks := negativeEdgeBlocks(edgeMetrics, *edgeBlock15mSamples, *edgeBlock15mMaxAvg, *edgeBlock1hSamples, *edgeBlock1hMaxAvg)
	participantsExclude, err := loadWalletSet(*participantsExcludeWalletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load participant excludes: %v\n", err)
		os.Exit(1)
	}
	scoreMeta, err := loadScoreMeta(*scoresPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load scores: %v\n", err)
		os.Exit(1)
	}
	now := time.Now()
	signals := buildBurstSignals(trades, now, *maxAge, *window, *minNotional, *minTrades, *minLegNotional)
	signals = append(signals, buildConsensusBurstSignals(trades, now, *maxAge, *window, *consensusMinNotional, *minTrades, *minLegNotional, *consensusMinWallets, *consensusMaxBot, edgeBlocks)...)
	sortBurstSignals(signals)
	var watchSignals []burstSignal
	if *consensusWatchMinNotional > 0 {
		watchSignals = consensusWatchSignals(
			buildConsensusBurstSignals(trades, now, *maxAge, *window, *consensusWatchMinNotional, *minTrades, *minLegNotional, *consensusWatchMinWallets, *consensusMaxBot, edgeBlocks),
			*consensusWatchMinNotional,
			*consensusMinNotional,
		)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	client := &http.Client{Timeout: minDuration(*timeout, 10*time.Second)}
	cache, err := loadMidpointCache(*markCachePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load mark cache: %v\n", err)
		os.Exit(1)
	}
	cacheDirty := false
	fetchMark := func(ctx context.Context, sig burstSignal) (float64, bool, error) {
		mid, err := fetchMidpoint(ctx, client, sig.Asset)
		if err == nil && mid >= 0 {
			if cache != nil {
				cache.Set(sig.Asset, mid, time.Now())
				cacheDirty = true
			}
			return mid, true, nil
		}
		mid, settlementErr := fetchSettledTokenPrice(ctx, client, *gammaBase, sig.Asset)
		if settlementErr != nil {
			mid, settlementErr = fetchSettledTokenPriceByEventSlug(ctx, client, *gammaBase, sig.Slug, sig.Asset)
		}
		if settlementErr == nil && mid >= 0 {
			if cache != nil {
				cache.Set(sig.Asset, mid, time.Now())
				cacheDirty = true
			}
			return mid, true, nil
		}
		if cache != nil {
			if cached, ok := cache.Get(sig.Asset); ok {
				return cached, true, nil
			}
		}
		return 0, false, err
	}
	results := evaluateBursts(ctx, signals, *stakeUSD, fetchMark)
	consensusWatchResults := evaluateBursts(ctx, watchSignals, *stakeUSD, fetchMark)
	historyMaxAge := *consensusHistoryMaxAge
	if historyMaxAge <= 0 {
		historyMaxAge = *maxAge
	}
	historyResults := results
	if historyMaxAge > *maxAge {
		historySignals := buildConsensusBurstSignals(trades, now, historyMaxAge, *window, *consensusMinNotional, *minTrades, *minLegNotional, *consensusMinWallets, *consensusMaxBot, edgeBlocks)
		historyResults = mergeConsensusHistoryResults(results, evaluateBursts(ctx, historySignals, *stakeUSD, fetchMark))
	}
	consensusAlerts, err := loadConsensusAlertEvents(*consensusAlertLogsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load consensus alert logs: %v\n", err)
		os.Exit(1)
	}
	if len(consensusAlerts) > 0 {
		alertResults := evaluateBursts(ctx, consensusAlertSignals(consensusAlerts), *stakeUSD, fetchMark)
		historyResults = mergeConsensusHistoryResults(historyResults, alertResults)
	}
	if cacheDirty {
		if err := saveMidpointCache(*markCachePath, cache); err != nil {
			fmt.Fprintf(os.Stderr, "sports-burst-report: save mark cache: %v\n", err)
			os.Exit(1)
		}
	}
	if err := writeConsensusEventsFile(*consensusEventsOutPath, historyResults, now); err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: write consensus events: %v\n", err)
		os.Exit(1)
	}
	if err := writeConsensusWatchEventsFile(*consensusWatchEventsOutPath, consensusWatchResults, now); err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: write consensus watch events: %v\n", err)
		os.Exit(1)
	}
	if err := writeConsensusParticipantsFile(*participantsOutPath, consensusParticipantSummaries(historyResults), participantsExclude, scoreMeta, *participantMinSignals, *participantMinROI); err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: write participants: %v\n", err)
		os.Exit(1)
	}
	consensusEvents, err := loadConsensusEventList(*consensusEventsOutPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load consensus events: %v\n", err)
		os.Exit(1)
	}
	consensusWatchEvents, err := loadConsensusEventList(*consensusWatchEventsOutPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load consensus watch events: %v\n", err)
		os.Exit(1)
	}
	consensusResearchRows, err := loadConsensusParticipantRows(*participantsOutPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: load participants: %v\n", err)
		os.Exit(1)
	}
	if err := writeReport(*reportPath, *tapePath, results, consensusWatchResults, consensusWatchEvents, consensusEvents, consensusResearchRows, *stakeUSD, *window, *maxAge, *minNotional, *minTrades, *minLegNotional, *consensusMinNotional, *consensusMinWallets, *consensusMaxBot, *consensusWatchMinNotional, *consensusWatchMinWallets, len(edgeBlocks), now); err != nil {
		fmt.Fprintf(os.Stderr, "sports-burst-report: write report: %v\n", err)
		os.Exit(1)
	}
	sum := summarize(results)
	fmt.Printf("sports-burst-report done: bursts=%d marked=%d roi=%.1f%% report=%s\n", sum.Signals, sum.Marked, sum.ReturnPC, *reportPath)
}

func loadTapeTrades(path string) ([]tapeTrade, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []tapeTrade
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var tr tapeTrade
		if err := json.Unmarshal(sc.Bytes(), &tr); err != nil {
			continue
		}
		tr.Wallet = strings.ToLower(strings.TrimSpace(tr.Wallet))
		tr.Asset = strings.TrimSpace(tr.Asset)
		tr.Side = strings.ToUpper(strings.TrimSpace(tr.Side))
		if tr.Wallet == "" || tr.Asset == "" || tr.Time.IsZero() {
			continue
		}
		out = append(out, tr)
	}
	return out, sc.Err()
}

func loadWalletStatuses(path string) (map[string]walletStatus, error) {
	out := map[string]walletStatus{}
	if path == "" {
		return out, nil
	}
	if strings.Contains(path, ",") {
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			set, err := loadWalletStatuses(part)
			if err != nil {
				return nil, err
			}
			for addr, status := range set {
				if shouldOverrideWalletStatus(out[addr], status) {
					out[addr] = status
				}
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
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		body, comment, _ := strings.Cut(line, "#")
		fields := strings.Fields(body)
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(strings.TrimSpace(fields[0]))
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		meta := parseCommentFields(comment)
		status := walletStatus{
			List: meta["list"],
			Tier: meta["tier"],
		}
		if isBlockedTapeGateStatus(meta["status"]) {
			status.List = "review_noise"
		}
		if raw := strings.TrimSpace(meta["bot"]); raw != "" {
			fmt.Sscanf(raw, "%f", &status.Bot)
		}
		out[addr] = status
	}
	return out, sc.Err()
}

func applyWalletStatuses(trades []tapeTrade, statuses map[string]walletStatus) {
	for i := range trades {
		status, ok := statuses[strings.ToLower(trades[i].Wallet)]
		if !ok {
			continue
		}
		if status.List != "" {
			trades[i].KnownList = status.List
		}
		if status.Tier != "" {
			trades[i].Tier = status.Tier
		}
		if status.Bot > 0 {
			trades[i].Bot = status.Bot
		}
	}
}

func isBlockedTapeGateStatus(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return strings.HasPrefix(status, "reject-") || strings.HasPrefix(status, "blocked-")
}

func shouldOverrideWalletStatus(existing, next walletStatus) bool {
	existingPriority := walletStatusPriority(existing.List)
	nextPriority := walletStatusPriority(next.List)
	if existingPriority != nextPriority {
		return nextPriority > existingPriority
	}
	return true
}

func walletStatusPriority(list string) int {
	switch strings.ToLower(strings.TrimSpace(list)) {
	case "tape_reversal":
		return 100
	case "review_noise":
		return 95
	case "tape_follow":
		return 90
	case "tape_candidate":
		return 80
	case "tape_edgehot":
		return 75
	case "tape_probation":
		return 70
	case "flow", "watch", "target", "sports", "scout":
		return 60
	case "consensus_research":
		return 20
	case "tape_observe":
		return 10
	default:
		return 0
	}
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
		body, _, _ := strings.Cut(sc.Text(), "#")
		fields := strings.Fields(body)
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(strings.TrimSpace(fields[0]))
		if strings.HasPrefix(addr, "0x") && len(addr) == 42 {
			_, comment, _ := strings.Cut(sc.Text(), "#")
			if participantExcludeRow(comment) {
				out[addr] = struct{}{}
			}
		}
	}
	return out, sc.Err()
}

func loadScoreMeta(path string) (map[string]participantMeta, error) {
	out := map[string]participantMeta{}
	if strings.TrimSpace(path) == "" {
		return out, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	var rows []struct {
		Address  string  `json:"address"`
		Tier     string  `json:"tier"`
		BotScore float64 `json:"bot_score"`
	}
	if err := json.Unmarshal(b, &rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		addr := strings.ToLower(strings.TrimSpace(row.Address))
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		out[addr] = participantMeta{Tier: row.Tier, Bot: row.BotScore}
	}
	return out, nil
}

func participantExcludeRow(comment string) bool {
	fields := parseCommentFields(comment)
	status := strings.ToLower(fields["status"])
	if status != "" {
		return strings.HasPrefix(status, "reject-") || strings.HasPrefix(status, "blocked-")
	}
	return true
}

func loadEdgeMetrics(path string) (map[string]map[int64]*edgeStats, error) {
	out := map[string]map[int64]*edgeStats{}
	if strings.TrimSpace(path) == "" {
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
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var snap edgeSnapshot
		if err := json.Unmarshal(sc.Bytes(), &snap); err != nil {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(snap.Wallet))
		if wallet == "" || snap.HorizonSec <= 0 {
			continue
		}
		byHorizon := out[wallet]
		if byHorizon == nil {
			byHorizon = map[int64]*edgeStats{}
			out[wallet] = byHorizon
		}
		st := byHorizon[snap.HorizonSec]
		if st == nil {
			st = &edgeStats{}
			byHorizon[snap.HorizonSec] = st
		}
		st.Samples++
		if snap.DeltaPP > 0 {
			st.Wins++
		}
		st.SumPP += snap.DeltaPP
	}
	return out, sc.Err()
}

func negativeEdgeBlocks(metrics map[string]map[int64]*edgeStats, min15mSamples int, max15mAvgPP float64, min1hSamples int, max1hAvgPP float64) map[string]string {
	out := map[string]string{}
	for wallet, byHorizon := range metrics {
		if st := byHorizon[int64((15 * time.Minute).Seconds())]; st != nil && min15mSamples > 0 && st.Samples >= min15mSamples {
			avg := avgEdge(st)
			if avg <= max15mAvgPP {
				out[wallet] = fmt.Sprintf("15m edge %.2fpp over %d samples", avg, st.Samples)
				continue
			}
		}
		if st := byHorizon[int64(time.Hour.Seconds())]; st != nil && min1hSamples > 0 && st.Samples >= min1hSamples {
			avg := avgEdge(st)
			if avg <= max1hAvgPP {
				out[wallet] = fmt.Sprintf("1h edge %.2fpp over %d samples", avg, st.Samples)
			}
		}
	}
	return out
}

func avgEdge(st *edgeStats) float64 {
	if st == nil || st.Samples == 0 {
		return 0
	}
	return st.SumPP / float64(st.Samples)
}

func parseCommentFields(comment string) map[string]string {
	out := map[string]string{}
	for _, field := range strings.Fields(comment) {
		k, v, ok := strings.Cut(field, "=")
		if ok {
			out[k] = v
		}
	}
	return out
}

func parseDecoratedInt(raw string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(raw))
	return n
}

func parseDecoratedFloat(raw string) float64 {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimSuffix(cleaned, "%")
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	if cleaned == "" || cleaned == "-" {
		return 0
	}
	n, _ := strconv.ParseFloat(cleaned, 64)
	return n
}

func buildBurstSignals(trades []tapeTrade, now time.Time, maxAge, window time.Duration, minNotional float64, minTrades int, minLegNotional float64) []burstSignal {
	type key struct {
		wallet string
		asset  string
	}
	groups := map[key][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Price <= 0 || tr.Notional < minLegNotional || tr.Size <= 0 {
			continue
		}
		if !isStrategyMode(alertMode(tr)) {
			continue
		}
		age := now.Sub(tr.Time)
		if maxAge > 0 && (age < 0 || age > maxAge) {
			continue
		}
		k := key{wallet: tr.Wallet, asset: tr.Asset}
		groups[k] = append(groups[k], tr)
	}

	var out []burstSignal
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= window {
				start--
			}
			windowRows := rows[start : end+1]
			if len(windowRows) < minTrades {
				continue
			}
			var totalNotional, totalSize float64
			for _, tr := range windowRows {
				totalNotional += tr.Notional
				totalSize += tr.Size
			}
			if totalNotional < minNotional || totalSize <= 0 {
				continue
			}
			first := windowRows[0]
			last := windowRows[len(windowRows)-1]
			out = append(out, burstSignal{
				Wallet:       last.Wallet,
				Scope:        "wallet",
				Participants: []string{last.Wallet},
				ParticipantMeta: map[string]participantMeta{
					last.Wallet: {Tier: last.Tier, Bot: last.Bot},
				},
				Asset:         last.Asset,
				Slug:          last.Slug,
				Mode:          alertMode(last),
				KnownList:     last.KnownList,
				Tier:          last.Tier,
				Bot:           last.Bot,
				Category:      last.Category,
				Market:        last.Market,
				Outcome:       last.Outcome,
				Trades:        len(windowRows),
				Wallets:       1,
				TotalNotional: totalNotional,
				TotalSize:     totalSize,
				VWAP:          totalNotional / totalSize,
				LastNotional:  last.Notional,
				FirstTime:     first.Time,
				LastTime:      last.Time,
			})
			break
		}
	}
	sortBurstSignals(out)
	return out
}

func buildConsensusBurstSignals(trades []tapeTrade, now time.Time, maxAge, window time.Duration, minNotional float64, minTrades int, minLegNotional float64, minWallets int, maxBot float64, edgeBlocks map[string]string) []burstSignal {
	groups := map[string][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Price <= 0 || tr.Notional < minLegNotional || tr.Size <= 0 || tr.Asset == "" {
			continue
		}
		if !isConsensusEligible(tr, maxBot, edgeBlocks) {
			continue
		}
		age := now.Sub(tr.Time)
		if maxAge > 0 && (age < 0 || age > maxAge) {
			continue
		}
		groups[tr.Asset] = append(groups[tr.Asset], tr)
	}

	var out []burstSignal
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= window {
				start--
			}
			windowRows := rows[start : end+1]
			if len(windowRows) < minTrades {
				continue
			}
			wallets := map[string]struct{}{}
			metas := map[string]participantMeta{}
			var totalNotional, totalSize, maxBotInWindow float64
			for _, tr := range windowRows {
				wallets[tr.Wallet] = struct{}{}
				meta := metas[tr.Wallet]
				if tierRank(tr.Tier) > tierRank(meta.Tier) {
					meta.Tier = tr.Tier
				}
				if tr.Bot > meta.Bot {
					meta.Bot = tr.Bot
				}
				metas[tr.Wallet] = meta
				totalNotional += tr.Notional
				totalSize += tr.Size
				if tr.Bot > maxBotInWindow {
					maxBotInWindow = tr.Bot
				}
			}
			if minWallets > 0 && len(wallets) < minWallets {
				continue
			}
			if totalNotional < minNotional || totalSize <= 0 {
				continue
			}
			participants := sortedWallets(wallets)
			first := windowRows[0]
			last := windowRows[len(windowRows)-1]
			out = append(out, burstSignal{
				Wallet:          fmt.Sprintf("multi:%d", len(wallets)),
				Scope:           "consensus",
				Participants:    participants,
				ParticipantMeta: metas,
				Asset:           last.Asset,
				Slug:            last.Slug,
				Mode:            "CONSENSUS",
				KnownList:       "consensus",
				Tier:            strongestTier(windowRows),
				Bot:             maxBotInWindow,
				Category:        last.Category,
				Market:          last.Market,
				Outcome:         last.Outcome,
				Trades:          len(windowRows),
				Wallets:         len(wallets),
				TotalNotional:   totalNotional,
				TotalSize:       totalSize,
				VWAP:            totalNotional / totalSize,
				LastNotional:    last.Notional,
				FirstTime:       first.Time,
				LastTime:        last.Time,
			})
			break
		}
	}
	sortBurstSignals(out)
	return out
}

func isConsensusEligible(tr tapeTrade, maxBot float64, edgeBlocks map[string]string) bool {
	switch strings.ToLower(strings.TrimSpace(tr.KnownList)) {
	case "review_noise", "tape_reversal":
		return false
	}
	if strings.EqualFold(strings.TrimSpace(tr.Tier), "BOT") {
		return false
	}
	if _, blocked := edgeBlocks[strings.ToLower(strings.TrimSpace(tr.Wallet))]; blocked {
		return false
	}
	if maxBot > 0 && tr.Bot > maxBot {
		return false
	}
	return true
}

func sortedWallets(wallets map[string]struct{}) []string {
	out := make([]string, 0, len(wallets))
	for wallet := range wallets {
		out = append(out, wallet)
	}
	sort.Strings(out)
	return out
}

func strongestTier(rows []tapeTrade) string {
	bestRank := 0
	best := ""
	for _, tr := range rows {
		tier := strings.ToUpper(strings.TrimSpace(tr.Tier))
		rank := tierRank(tier)
		if rank > bestRank {
			bestRank = rank
			best = tier
		}
	}
	return best
}

func tierRank(tier string) int {
	switch tier {
	case "A":
		return 4
	case "B":
		return 3
	case "C":
		return 2
	case "D":
		return 1
	default:
		return 0
	}
}

func sortBurstSignals(out []burstSignal) {
	sort.Slice(out, func(i, j int) bool {
		if out[i].TotalNotional != out[j].TotalNotional {
			return out[i].TotalNotional > out[j].TotalNotional
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
}

func consensusWatchSignals(signals []burstSignal, watchMinNotional, officialMinNotional float64) []burstSignal {
	if watchMinNotional <= 0 {
		return nil
	}
	out := make([]burstSignal, 0, len(signals))
	for _, sig := range signals {
		if !strings.EqualFold(sig.Mode, "CONSENSUS") {
			continue
		}
		if officialMinNotional > 0 && sig.TotalNotional >= officialMinNotional {
			continue
		}
		sig.Mode = "CONSENSUS-WATCH"
		sig.Scope = "consensus-watch"
		out = append(out, sig)
	}
	sortBurstSignals(out)
	return out
}

func evaluateBursts(ctx context.Context, signals []burstSignal, stakeUSD float64, fetch func(context.Context, burstSignal) (float64, bool, error)) []burstResult {
	type mark struct {
		mid float64
		ok  bool
	}
	marks := map[string]mark{}
	out := make([]burstResult, 0, len(signals))
	for _, sig := range signals {
		res := burstResult{Signal: sig, StakeUSD: stakeUSD}
		if sig.VWAP > 0 && stakeUSD > 0 {
			res.Units = stakeUSD / sig.VWAP
		}
		cached, ok := marks[sig.Asset]
		if !ok {
			mid, marked, err := fetch(ctx, sig)
			cached = mark{mid: mid, ok: err == nil && marked && mid >= 0}
			if cached.ok {
				marks[sig.Asset] = cached
			}
		}
		if cached.ok && res.Units > 0 {
			res.Mid = cached.mid
			res.Marked = true
			res.PnLUSD = res.Units * (cached.mid - sig.VWAP)
			res.ReturnPC = res.PnLUSD / stakeUSD * 100
		}
		out = append(out, res)
	}
	return out
}

func loadConsensusAlertEvents(paths string) ([]consensusAlertEvent, error) {
	var out []consensusAlertEvent
	for _, raw := range strings.Split(paths, ",") {
		path := strings.TrimSpace(raw)
		if path == "" {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
		for sc.Scan() {
			var ev consensusAlertEvent
			if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
				continue
			}
			if !strings.EqualFold(strings.TrimSpace(ev.Mode), "CONSENSUS") {
				continue
			}
			if strings.TrimSpace(ev.Asset) == "" || ev.TradeTime.IsZero() || ev.Price <= 0 {
				continue
			}
			out = append(out, ev)
		}
		if err := sc.Err(); err != nil {
			_ = f.Close()
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func consensusAlertSignals(events []consensusAlertEvent) []burstSignal {
	out := make([]burstSignal, 0, len(events))
	for _, ev := range events {
		wallets := consensusAlertWallets(ev)
		wallet := strings.TrimSpace(ev.Wallet)
		if wallet == "" && wallets > 0 {
			wallet = fmt.Sprintf("multi:%d", wallets)
		}
		trades := len(ev.Participants)
		if trades == 0 {
			trades = wallets
		}
		if trades <= 0 {
			trades = 1
		}
		out = append(out, burstSignal{
			Wallet:        wallet,
			Scope:         "alert-log",
			Participants:  append([]string(nil), ev.Participants...),
			Asset:         ev.Asset,
			Slug:          ev.Slug,
			Mode:          "CONSENSUS",
			KnownList:     ev.KnownList,
			Tier:          ev.Tier,
			Bot:           ev.Bot,
			Category:      ev.Category,
			Market:        ev.Market,
			Outcome:       ev.Outcome,
			Trades:        trades,
			Wallets:       wallets,
			TotalNotional: ev.Notional,
			VWAP:          ev.Price,
			FirstTime:     ev.TradeTime,
			LastTime:      ev.TradeTime,
		})
	}
	sortBurstSignals(out)
	return out
}

func consensusAlertWallets(ev consensusAlertEvent) int {
	if len(ev.Participants) > 0 {
		return len(ev.Participants)
	}
	wallet := strings.TrimSpace(ev.Wallet)
	if strings.HasPrefix(strings.ToLower(wallet), "multi:") {
		n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(strings.ToLower(wallet), "multi:")))
		if err == nil && n > 0 {
			return n
		}
	}
	if wallet != "" {
		return 1
	}
	return 0
}

func mergeConsensusHistoryResults(base, history []burstResult) []burstResult {
	out := append([]burstResult(nil), base...)
	seen := map[string]struct{}{}
	for _, res := range out {
		if !strings.EqualFold(res.Signal.Mode, "CONSENSUS") {
			continue
		}
		key := consensusEventKey(res.Signal.Asset, res.Signal.LastTime, res.Signal.Wallets, res.Signal.TotalNotional)
		if key != "" {
			seen[key] = struct{}{}
		}
	}
	for _, res := range history {
		if !strings.EqualFold(res.Signal.Mode, "CONSENSUS") {
			continue
		}
		key := consensusEventKey(res.Signal.Asset, res.Signal.LastTime, res.Signal.Wallets, res.Signal.TotalNotional)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
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
	return strconv.ParseFloat(body.Mid, 64)
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

func writeReport(path, tapePath string, results []burstResult, consensusWatchResults []burstResult, consensusWatchEvents []consensusEvent, consensusEvents []consensusEvent, consensusResearchRows []participantSummary, stakeUSD float64, window, maxAge time.Duration, minNotional float64, minTrades int, minLegNotional float64, consensusMinNotional float64, consensusMinWallets int, consensusMaxBot float64, consensusWatchMinNotional float64, consensusWatchMinWallets int, negativeEdgeBlocked int, now time.Time) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	sum := summarize(results)
	var b strings.Builder
	fmt.Fprintf(&b, "# Sports Burst Performance\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", now.Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Tape: `%s`\n", tapePath)
	fmt.Fprintf(&b, "- Fixed paper stake: $%.2f per burst\n", stakeUSD)
	fmt.Fprintf(&b, "- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices and local cache\n")
	fmt.Fprintf(&b, "- Rule: same wallet + same asset, window=%s, min total=$%.0f, min trades=%d, min leg=$%.0f\n\n", window, minNotional, minTrades, minLegNotional)
	fmt.Fprintf(&b, "- Consensus rule: same asset across wallets, min total=$%.0f, min wallets=%d, max bot=%.1f; excludes review-noise, reversal-risk, negative-edge, and BOT tiers\n", consensusMinNotional, consensusMinWallets, consensusMaxBot)
	fmt.Fprintf(&b, "- Negative-edge blocked wallets: %d\n\n", negativeEdgeBlocked)
	fmt.Fprintf(&b, "## Summary\n\n")
	fmt.Fprintf(&b, "- Bursts: %d\n", sum.Signals)
	fmt.Fprintf(&b, "- Marked to current midpoint: %d\n", sum.Marked)
	fmt.Fprintf(&b, "- Unmarked: %d\n", sum.Unmarked)
	fmt.Fprintf(&b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(sum.Wins), float64(sum.Marked)))
	fmt.Fprintf(&b, "- PnL incl. midpoint marks: $%+.2f\n", sum.PnLUSD)
	fmt.Fprintf(&b, "- ROI incl. midpoint marks: %.1f%%\n", sum.ReturnPC)
	fmt.Fprintf(&b, "- Avg price delta: %+.2fpp\n\n", sum.AvgDeltaP)
	writeGateSection(&b, "Mode Gates", results, func(r burstResult) string { return r.Signal.Mode })
	writeGateSection(&b, "Scope Gates", results, func(r burstResult) string { return r.Signal.Scope })
	writeGateSection(&b, "Wallet Gates", results, func(r burstResult) string { return shortAddr(r.Signal.Wallet) })
	writeConsensusParticipantSection(&b, results)
	writeConsensusWatchSection(&b, consensusWatchResults, consensusWatchEvents, consensusWatchMinNotional, consensusWatchMinWallets, stakeUSD)
	writeConsensusHistorySection(&b, consensusEvents, stakeUSD)
	writeConsensusResearchWalletSection(&b, consensusResearchRows)
	writeRecentSection(&b, results, now)
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeConsensusWatchSection(b *strings.Builder, results []burstResult, events []consensusEvent, minNotional float64, minWallets int, stakeUSD float64) {
	sum := summarize(results)
	eventSum := summarizeConsensusEvents(events, stakeUSD)
	fmt.Fprintf(b, "## Research-Only Consensus Watch\n\n")
	fmt.Fprintf(b, "- Rule: lower-threshold cross-wallet same-asset BUY bursts for sample discovery only; not Telegram-pushed and not counted as official CONSENSUS\n")
	fmt.Fprintf(b, "- Threshold: total>=$%.0f wallets>=%d\n", minNotional, minWallets)
	fmt.Fprintf(b, "- Watch bursts: %d\n", sum.Signals)
	fmt.Fprintf(b, "- Marked to current midpoint/settlement: %d\n", sum.Marked)
	fmt.Fprintf(b, "- ROI incl. midpoint marks: %.1f%%\n", sum.ReturnPC)
	fmt.Fprintf(b, "- Durable watch events: %d\n", eventSum.Signals)
	fmt.Fprintf(b, "- Durable watch marked: %d\n", eventSum.Marked)
	fmt.Fprintf(b, "- Durable watch ROI: %.1f%%\n\n", eventSum.ReturnPC)
	fmt.Fprintf(b, "| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Participants | Market |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|---|---|\n")
	if len(results) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | $0 | 0.000 | 0.000 | 0.0%% |  |  |\n\n")
		return
	}
	rows := append([]burstResult(nil), results...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Signal.LastTime.Equal(rows[j].Signal.LastTime) {
			return rows[i].Signal.TotalNotional > rows[j].Signal.TotalNotional
		}
		return rows[i].Signal.LastTime.After(rows[j].Signal.LastTime)
	})
	limit := minInt(len(rows), 15)
	for i := 0; i < limit; i++ {
		r := rows[i]
		mid := "-"
		if r.Marked {
			mid = fmt.Sprintf("%.3f", r.Mid)
		}
		fmt.Fprintf(b, "| %s | %d | %d | $%.0f | %.3f | %s | %.1f%% | %s | %s |\n",
			formatEventTime(r.Signal.LastTime), r.Signal.Wallets, r.Signal.Trades, r.Signal.TotalNotional, r.Signal.VWAP, mid, r.ReturnPC, formatParticipants(r.Signal.Participants, 6), oneLine(r.Signal.Market, 80))
	}
	fmt.Fprintf(b, "\n")

	fmt.Fprintf(b, "### Durable Watch History\n\n")
	fmt.Fprintf(b, "| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Market |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|---|\n")
	if len(events) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | $0 | 0.000 | 0.000 | 0.0%% |  |\n\n")
		return
	}
	limitEvents := minInt(len(events), 15)
	for i := 0; i < limitEvents; i++ {
		ev := events[i]
		mid := "-"
		if ev.Marked {
			mid = fmt.Sprintf("%.3f", ev.Mid)
		}
		fmt.Fprintf(b, "| %s | %d | %d | $%.0f | %.3f | %s | %.1f%% | %s |\n",
			formatEventTime(ev.LastTime), ev.Wallets, ev.Trades, ev.TotalNotional, ev.VWAP, mid, ev.ReturnPC, oneLine(ev.Market, 80))
	}
	fmt.Fprintf(b, "\n")

	writeConsensusEventParticipantSection(b, "Durable Watch Participants", events)
}

func writeConsensusHistorySection(b *strings.Builder, events []consensusEvent, stakeUSD float64) {
	sum := summarizeConsensusEvents(events, stakeUSD)
	needed := 0
	if sum.Marked < 5 {
		needed = 5 - sum.Marked
	}
	fmt.Fprintf(b, "## Durable Consensus Event History\n\n")
	fmt.Fprintf(b, "- Rule: persisted CONSENSUS events imported from live and shadow alert logs; this survives tape-window rollover\n")
	fmt.Fprintf(b, "- Events: %d\n", sum.Signals)
	fmt.Fprintf(b, "- Marked to current midpoint/settlement: %d\n", sum.Marked)
	fmt.Fprintf(b, "- Marked samples still needed for promotion review: %d\n", needed)
	fmt.Fprintf(b, "- Win rate incl. midpoint marks: %.1f%%\n", pct(float64(sum.Wins), float64(sum.Marked)))
	fmt.Fprintf(b, "- PnL incl. midpoint marks: $%+.2f\n", sum.PnLUSD)
	fmt.Fprintf(b, "- ROI incl. midpoint marks: %.1f%%\n", sum.ReturnPC)
	fmt.Fprintf(b, "- Gate action: %s\n\n", gateAction(sum))
	fmt.Fprintf(b, "| Last | Category | Wallets | Trades | Total | VWAP | Mid | ROI | Market |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---|\n")
	if len(events) == 0 {
		fmt.Fprintf(b, "| n/a |  | 0 | 0 | $0 | 0.000 | 0.000 | 0.0%% |  |\n\n")
		return
	}
	limit := minInt(len(events), 15)
	for i := 0; i < limit; i++ {
		ev := events[i]
		mid := "-"
		if ev.Marked {
			mid = fmt.Sprintf("%.3f", ev.Mid)
		}
		fmt.Fprintf(b, "| %s | %s | %d | %d | $%.0f | %.3f | %s | %.1f%% | %s |\n",
			formatEventTime(ev.LastTime), dash(ev.Category), ev.Wallets, ev.Trades, ev.TotalNotional, ev.VWAP, mid, ev.ReturnPC, oneLine(ev.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeConsensusEventParticipantSection(b *strings.Builder, title string, events []consensusEvent) {
	rows := consensusEventParticipantSummaries(events)
	fmt.Fprintf(b, "### %s\n\n", title)
	fmt.Fprintf(b, "| Wallet | Tier | Bot | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| n/a | - | 0.0 | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | $0 | +0.00 |\n\n")
		return
	}
	limit := minInt(len(rows), 25)
	for i := 0; i < limit; i++ {
		r := rows[i]
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %d | %d | %.1f%% | %.1f%% | $%+.2f | $%.0f | %+.2f |\n",
			shortAddr(r.Wallet), dash(r.Tier), r.Bot, r.Signals, r.Marked, pct(float64(r.Wins), float64(r.Marked)), r.ReturnPC, r.PnLUSD, r.TotalNotional, r.AvgDeltaP)
	}
	fmt.Fprintf(b, "\n")
}

func summarizeConsensusEvents(events []consensusEvent, stakeUSD float64) summary {
	var s summary
	for _, ev := range events {
		s.Signals++
		if !ev.Marked {
			s.Unmarked++
			continue
		}
		s.Marked++
		if ev.PnLUSD > 0 {
			s.Wins++
		}
		s.StakeUSD += stakeUSD
		s.PnLUSD += ev.PnLUSD
		s.AvgDeltaP += (ev.Mid - ev.VWAP) * 100
	}
	if s.StakeUSD > 0 {
		s.ReturnPC = s.PnLUSD / s.StakeUSD * 100
	}
	if s.Marked > 0 {
		s.AvgDeltaP /= float64(s.Marked)
	}
	return s
}

func consensusEventParticipantSummaries(events []consensusEvent) []participantSummary {
	byWallet := map[string]*participantSummary{}
	for _, ev := range events {
		for _, wallet := range ev.Participants {
			wallet = strings.ToLower(strings.TrimSpace(wallet))
			if wallet == "" {
				continue
			}
			row := byWallet[wallet]
			if row == nil {
				row = &participantSummary{Wallet: wallet}
				byWallet[wallet] = row
			}
			meta := ev.ParticipantMeta[wallet]
			if tierRank(meta.Tier) > tierRank(row.Tier) {
				row.Tier = meta.Tier
			}
			if meta.Bot > row.Bot {
				row.Bot = meta.Bot
			}
			row.Signals++
			row.TotalNotional += ev.TotalNotional
			if !ev.Marked {
				continue
			}
			row.Marked++
			if ev.PnLUSD > 0 {
				row.Wins++
			}
			row.PnLUSD += ev.PnLUSD
			row.StakeUSD += 10
			row.AvgDeltaP += (ev.Mid - ev.VWAP) * 100
		}
	}
	rows := make([]participantSummary, 0, len(byWallet))
	for _, row := range byWallet {
		if row.StakeUSD > 0 {
			row.ReturnPC = row.PnLUSD / row.StakeUSD * 100
		}
		if row.Marked > 0 {
			row.AvgDeltaP /= float64(row.Marked)
		}
		rows = append(rows, *row)
	}
	sortParticipantSummaries(rows)
	return rows
}

func formatEventTime(t time.Time) string {
	if t.IsZero() {
		return "n/a"
	}
	return t.Format("01-02 15:04")
}

func writeConsensusResearchWalletSection(b *strings.Builder, rows []participantSummary) {
	sortParticipantSummaries(rows)
	fmt.Fprintf(b, "## Durable Consensus Research Wallets\n\n")
	fmt.Fprintf(b, "- Rule: persisted address-level attribution from positive CONSENSUS events; wallets stay research-only until repeated marked samples pass promotion gates\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Bot | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| n/a | - | 0.0 | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | $0 | +0.00 |\n\n")
		return
	}
	limit := minInt(len(rows), 25)
	for i := 0; i < limit; i++ {
		r := rows[i]
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %d | %d | %.1f%% | %.1f%% | $%+.2f | $%.0f | %+.2f |\n",
			shortAddr(r.Wallet), dash(r.Tier), r.Bot, r.Signals, r.Marked, pct(float64(r.Wins), float64(r.Marked)), r.ReturnPC, r.PnLUSD, r.TotalNotional, r.AvgDeltaP)
	}
	fmt.Fprintf(b, "\n")
}

func writeGateSection(b *strings.Builder, title string, results []burstResult, keyFn func(burstResult) string) {
	groups := map[string][]burstResult{}
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
	fmt.Fprintf(b, "- Collect-positive gate: ROI>0 with sample below promote gate\n\n")
	fmt.Fprintf(b, "| Key | Bursts | Marked | Win | ROI | PnL | Action |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---|\n")
	if len(keys) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | COLLECT |\n\n")
		return
	}
	for _, k := range keys {
		s := summarize(groups[k])
		fmt.Fprintf(b, "| `%s` | %d | %d | %.1f%% | %.1f%% | $%+.2f | %s |\n",
			k, s.Signals, s.Marked, pct(float64(s.Wins), float64(s.Marked)), s.ReturnPC, s.PnLUSD, gateAction(s))
	}
	fmt.Fprintf(b, "\n")
}

func writeConsensusParticipantSection(b *strings.Builder, results []burstResult) {
	rows := consensusParticipantSummaries(results)
	fmt.Fprintf(b, "## Consensus Participants\n\n")
	fmt.Fprintf(b, "- Rule: attribution table for wallets appearing in CONSENSUS bursts; PnL is attributed to each participant for ranking only\n\n")
	fmt.Fprintf(b, "| Wallet | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|---:|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| n/a | 0 | 0 | 0.0%% | 0.0%% | $+0.00 | $0 | +0.00 |\n\n")
		return
	}
	limit := minInt(len(rows), 25)
	for i := 0; i < limit; i++ {
		r := rows[i]
		fmt.Fprintf(b, "| `%s` | %d | %d | %.1f%% | %.1f%% | $%+.2f | $%.0f | %+.2f |\n",
			shortAddr(r.Wallet), r.Signals, r.Marked, pct(float64(r.Wins), float64(r.Marked)), r.ReturnPC, r.PnLUSD, r.TotalNotional, r.AvgDeltaP)
	}
	fmt.Fprintf(b, "\n")
}

func consensusParticipantSummaries(results []burstResult) []participantSummary {
	byWallet := map[string]*participantSummary{}
	for _, res := range results {
		if !strings.EqualFold(res.Signal.Mode, "CONSENSUS") {
			continue
		}
		for _, wallet := range res.Signal.Participants {
			wallet = strings.ToLower(strings.TrimSpace(wallet))
			if wallet == "" {
				continue
			}
			row := byWallet[wallet]
			if row == nil {
				row = &participantSummary{Wallet: wallet}
				byWallet[wallet] = row
			}
			meta := res.Signal.ParticipantMeta[wallet]
			if tierRank(meta.Tier) > tierRank(row.Tier) {
				row.Tier = meta.Tier
			}
			if meta.Bot > row.Bot {
				row.Bot = meta.Bot
			}
			row.Signals++
			row.TotalNotional += res.Signal.TotalNotional
			if !res.Marked {
				continue
			}
			row.Marked++
			if res.PnLUSD > 0 {
				row.Wins++
			}
			row.PnLUSD += res.PnLUSD
			row.StakeUSD += res.StakeUSD
			row.AvgDeltaP += (res.Mid - res.Signal.VWAP) * 100
		}
	}
	rows := make([]participantSummary, 0, len(byWallet))
	for _, row := range byWallet {
		if row.StakeUSD > 0 {
			row.ReturnPC = row.PnLUSD / row.StakeUSD * 100
		}
		if row.Marked > 0 {
			row.AvgDeltaP /= float64(row.Marked)
		}
		rows = append(rows, *row)
	}
	sortParticipantSummaries(rows)
	return rows
}

func sortParticipantSummaries(rows []participantSummary) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Marked != rows[j].Marked {
			return rows[i].Marked > rows[j].Marked
		}
		if rows[i].ReturnPC != rows[j].ReturnPC {
			return rows[i].ReturnPC > rows[j].ReturnPC
		}
		if rows[i].TotalNotional != rows[j].TotalNotional {
			return rows[i].TotalNotional > rows[j].TotalNotional
		}
		return rows[i].Wallet < rows[j].Wallet
	})
}

func writeConsensusParticipantsFile(path string, rows []participantSummary, exclude map[string]struct{}, scoreMeta map[string]participantMeta, minSignals int, minROI float64) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	lines, order, err := loadExistingConsensusParticipantLines(path, exclude)
	if err != nil {
		return err
	}
	for wallet, line := range lines {
		if meta, ok := scoreMeta[wallet]; ok {
			lines[wallet] = backfillConsensusParticipantLine(line, meta)
		}
	}
	var b strings.Builder
	for _, row := range rows {
		wallet := strings.ToLower(strings.TrimSpace(row.Wallet))
		if _, blocked := exclude[wallet]; blocked {
			continue
		}
		if row.Signals < minSignals || row.Marked == 0 || row.ReturnPC <= minROI {
			continue
		}
		if meta, ok := scoreMeta[wallet]; ok {
			if strings.TrimSpace(row.Tier) == "" || row.Tier == "-" {
				row.Tier = meta.Tier
			}
			if row.Bot <= 0 {
				row.Bot = meta.Bot
			}
		}
		if _, exists := lines[wallet]; !exists {
			order = append(order, wallet)
		}
		lines[wallet] = consensusParticipantLine(row)
	}
	for _, wallet := range order {
		line, ok := lines[wallet]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "%s\n", line)
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func consensusParticipantLine(row participantSummary) string {
	return fmt.Sprintf("%s # list=consensus_research tier=%s bot=%.1f signals=%d marked=%d win=%.1f%% roi=%.1f%% pnl=$%+.2f notional=$%.0f avgDeltaPP=%+.2f source=sports_consensus",
		strings.ToLower(strings.TrimSpace(row.Wallet)), dash(row.Tier), row.Bot, row.Signals, row.Marked, pct(float64(row.Wins), float64(row.Marked)), row.ReturnPC, row.PnLUSD, row.TotalNotional, row.AvgDeltaP)
}

func backfillConsensusParticipantLine(line string, meta participantMeta) string {
	meta.Tier = strings.TrimSpace(meta.Tier)
	if meta.Tier == "" && meta.Bot <= 0 {
		return line
	}
	body, comment, ok := strings.Cut(line, "#")
	if !ok {
		return line
	}
	fields := parseCommentFields(comment)
	if !strings.EqualFold(fields["list"], "consensus_research") {
		return line
	}
	if meta.Tier != "" {
		fields["tier"] = meta.Tier
	}
	if meta.Bot > 0 {
		fields["bot"] = fmt.Sprintf("%.1f", meta.Bot)
	}
	parts := []string{"list=consensus_research"}
	if fields["tier"] != "" {
		parts = append(parts, "tier="+fields["tier"])
	}
	if fields["bot"] != "" {
		parts = append(parts, "bot="+fields["bot"])
	}
	for _, raw := range strings.Fields(comment) {
		k, _, ok := strings.Cut(raw, "=")
		if !ok {
			continue
		}
		switch k {
		case "list", "tier", "bot":
			continue
		default:
			parts = append(parts, raw)
		}
	}
	return strings.TrimSpace(body) + " # " + strings.Join(parts, " ")
}

func loadConsensusParticipantRows(path string) ([]participantSummary, error) {
	path = strings.TrimSpace(path)
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

	var rows []participantSummary
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		body, comment, _ := strings.Cut(line, "#")
		fields := strings.Fields(body)
		if len(fields) == 0 {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(fields[0]))
		if !strings.HasPrefix(wallet, "0x") || len(wallet) != 42 {
			continue
		}
		meta := parseCommentFields(comment)
		if !strings.EqualFold(meta["list"], "consensus_research") {
			continue
		}
		row := participantSummary{
			Wallet:        wallet,
			Tier:          meta["tier"],
			Bot:           parseDecoratedFloat(meta["bot"]),
			Signals:       parseDecoratedInt(meta["signals"]),
			Marked:        parseDecoratedInt(meta["marked"]),
			ReturnPC:      parseDecoratedFloat(meta["roi"]),
			PnLUSD:        parseDecoratedFloat(meta["pnl"]),
			TotalNotional: parseDecoratedFloat(meta["notional"]),
			AvgDeltaP:     parseDecoratedFloat(meta["avgDeltaPP"]),
		}
		winRate := parseDecoratedFloat(meta["win"])
		if row.Marked > 0 && winRate > 0 {
			row.Wins = int(math.Round(float64(row.Marked) * winRate / 100))
		}
		if row.Marked > 0 {
			row.StakeUSD = float64(row.Marked) * 10
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sortParticipantSummaries(rows)
	return rows, nil
}

func loadExistingConsensusParticipantLines(path string, exclude map[string]struct{}) (map[string]string, []string, error) {
	lines := map[string]string{}
	var order []string
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lines, order, nil
		}
		return nil, nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		body, comment, _ := strings.Cut(line, "#")
		fields := strings.Fields(body)
		if len(fields) == 0 {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(fields[0]))
		if !strings.HasPrefix(wallet, "0x") || len(wallet) != 42 {
			continue
		}
		if _, blocked := exclude[wallet]; blocked {
			continue
		}
		meta := parseCommentFields(comment)
		if !strings.EqualFold(meta["list"], "consensus_research") {
			continue
		}
		if _, exists := lines[wallet]; exists {
			continue
		}
		lines[wallet] = wallet + " #" + comment
		order = append(order, wallet)
	}
	return lines, order, sc.Err()
}

func writeConsensusEventsFile(path string, results []burstResult, now time.Time) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	existing, err := loadConsensusEvents(path)
	if err != nil {
		return err
	}
	for _, res := range results {
		if !strings.EqualFold(res.Signal.Mode, "CONSENSUS") {
			continue
		}
		ev := consensusEventFromResult(res, now)
		if ev.Key == "" {
			continue
		}
		existing[ev.Key] = ev
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	keys := make([]string, 0, len(existing))
	for key := range existing {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := existing[keys[i]], existing[keys[j]]
		if !a.LastTime.Equal(b.LastTime) {
			return a.LastTime.After(b.LastTime)
		}
		return keys[i] < keys[j]
	})
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for _, key := range keys {
		if err := enc.Encode(existing[key]); err != nil {
			_ = f.Close()
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func writeConsensusWatchEventsFile(path string, results []burstResult, now time.Time) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	existing, err := loadConsensusEvents(path)
	if err != nil {
		return err
	}
	for _, res := range results {
		if !strings.EqualFold(res.Signal.Mode, "CONSENSUS-WATCH") {
			continue
		}
		ev := consensusEventFromResult(res, now)
		if ev.Key == "" {
			continue
		}
		existing[ev.Key] = ev
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	keys := make([]string, 0, len(existing))
	for key := range existing {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := existing[keys[i]], existing[keys[j]]
		if !a.LastTime.Equal(b.LastTime) {
			return a.LastTime.After(b.LastTime)
		}
		return keys[i] < keys[j]
	})
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for _, key := range keys {
		if err := enc.Encode(existing[key]); err != nil {
			_ = f.Close()
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func loadConsensusEvents(path string) (map[string]consensusEvent, error) {
	out := map[string]consensusEvent{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var ev consensusEvent
		if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
			continue
		}
		if ev.Key == "" {
			ev.Key = consensusEventKey(ev.Asset, ev.LastTime, ev.Wallets, ev.TotalNotional)
		}
		if ev.Key == "" {
			continue
		}
		ev.Participants = nonNilStrings(ev.Participants)
		out[ev.Key] = ev
	}
	return out, sc.Err()
}

func loadConsensusEventList(path string) ([]consensusEvent, error) {
	events, err := loadConsensusEvents(path)
	if err != nil {
		return nil, err
	}
	out := make([]consensusEvent, 0, len(events))
	for _, ev := range events {
		out = append(out, ev)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].LastTime.Equal(out[j].LastTime) {
			return out[i].LastTime.After(out[j].LastTime)
		}
		return out[i].Key < out[j].Key
	})
	return out, nil
}

func consensusEventFromResult(res burstResult, now time.Time) consensusEvent {
	sig := res.Signal
	ev := consensusEvent{
		Key:             consensusEventKeyForMode(sig.Mode, sig.Asset, sig.LastTime, sig.Wallets, sig.TotalNotional),
		Mode:            strings.ToUpper(strings.TrimSpace(sig.Mode)),
		FirstTime:       sig.FirstTime,
		LastTime:        sig.LastTime,
		Asset:           sig.Asset,
		Slug:            sig.Slug,
		Category:        sig.Category,
		Market:          sig.Market,
		Outcome:         sig.Outcome,
		Wallets:         sig.Wallets,
		Trades:          sig.Trades,
		TotalNotional:   sig.TotalNotional,
		VWAP:            sig.VWAP,
		Marked:          res.Marked,
		Participants:    nonNilStrings(sig.Participants),
		ParticipantMeta: sig.ParticipantMeta,
		UpdatedAt:       now,
	}
	if res.Marked {
		ev.Mid = res.Mid
		ev.PnLUSD = res.PnLUSD
		ev.ReturnPC = res.ReturnPC
	}
	return ev
}

func nonNilStrings(in []string) []string {
	if len(in) == 0 {
		return []string{}
	}
	return append([]string(nil), in...)
}

func consensusEventKey(asset string, lastTime time.Time, wallets int, totalNotional float64) string {
	return consensusEventKeyForMode("CONSENSUS", asset, lastTime, wallets, totalNotional)
}

func consensusEventKeyForMode(mode, asset string, lastTime time.Time, wallets int, totalNotional float64) string {
	asset = strings.TrimSpace(asset)
	if asset == "" || lastTime.IsZero() {
		return ""
	}
	prefix := "consensus"
	if strings.EqualFold(strings.TrimSpace(mode), "CONSENSUS-WATCH") {
		prefix = "consensus-watch"
	}
	return strings.ToLower(fmt.Sprintf("%s|%s|%d|%d|%.4f", prefix, asset, lastTime.Unix(), wallets, totalNotional))
}

func writeRecentSection(b *strings.Builder, results []burstResult, now time.Time) {
	rows := append([]burstResult(nil), results...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].Signal.LastTime.After(rows[j].Signal.LastTime) })
	if len(rows) > 20 {
		rows = rows[:20]
	}
	fmt.Fprintf(b, "## Recent Bursts\n\n")
	fmt.Fprintf(b, "| Last | Scope | Mode | Actor | Participants | Wallets | Trades | Total | VWAP | Mid | ROI | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---|---|---:|---:|---:|---:|---:|---:|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| n/a |  |  |  |  | 0 | 0 | $0 | 0.000 | 0.000 | 0.0%% |  |\n")
		return
	}
	for _, r := range rows {
		mid := "-"
		if r.Marked {
			mid = fmt.Sprintf("%.3f", r.Mid)
		}
		fmt.Fprintf(b, "| %s | %s | %s | `%s` | %s | %d | %d | $%.0f | %.3f | %s | %.1f%% | %s |\n",
			formatDuration(now.Sub(r.Signal.LastTime)), r.Signal.Scope, r.Signal.Mode, shortAddr(r.Signal.Wallet), formatParticipants(r.Signal.Participants, 8), r.Signal.Wallets, r.Signal.Trades,
			r.Signal.TotalNotional, r.Signal.VWAP, mid, r.ReturnPC, oneLine(r.Signal.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func summarize(results []burstResult) summary {
	var s summary
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
		s.AvgDeltaP += (r.Mid - r.Signal.VWAP) * 100
	}
	if s.StakeUSD > 0 {
		s.ReturnPC = s.PnLUSD / s.StakeUSD * 100
	}
	if s.Marked > 0 {
		s.AvgDeltaP /= float64(s.Marked)
	}
	return s
}

func gateAction(s summary) string {
	winRate := pct(float64(s.Wins), float64(s.Marked))
	switch {
	case s.Marked >= 5 && s.ReturnPC >= 5 && winRate >= 60:
		return "PROMOTE_CANDIDATE"
	case s.ReturnPC > 0:
		return "COLLECT_POSITIVE"
	default:
		return "COLLECT"
	}
}

func alertMode(tr tapeTrade) string {
	switch strings.ToLower(strings.TrimSpace(tr.KnownList)) {
	case "tape_follow":
		return "FOLLOW-READY"
	case "tape_candidate":
		return "CANDIDATE"
	case "tape_probation":
		return "PROBATION"
	case "flow":
		return "FLOW-SCOUT"
	case "watch", "target", "sports", "scout":
		return "EDGE-HOT"
	default:
		return "OBSERVE"
	}
}

func isStrategyMode(mode string) bool {
	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case "FOLLOW-READY", "CANDIDATE", "PROBATION", "FLOW-SCOUT", "EDGE-HOT":
		return true
	default:
		return false
	}
}

func pct(n, d float64) float64 {
	if d <= 0 {
		return 0
	}
	return n / d * 100
}

func dash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func formatParticipants(wallets []string, limit int) string {
	if len(wallets) == 0 {
		return ""
	}
	if limit <= 0 || limit > len(wallets) {
		limit = len(wallets)
	}
	parts := make([]string, 0, limit+1)
	for i := 0; i < limit; i++ {
		parts = append(parts, "`"+shortAddr(wallets[i])+"`")
	}
	if len(wallets) > limit {
		parts = append(parts, fmt.Sprintf("+%d more", len(wallets)-limit))
	}
	return strings.Join(parts, ", ")
}

func oneLine(s string, limit int) string {
	s = strings.Join(strings.Fields(s), " ")
	if limit > 0 && len(s) > limit {
		return s[:limit-1] + "..."
	}
	return s
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d >= time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.0fm", d.Minutes())
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
