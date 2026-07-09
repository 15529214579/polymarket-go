package main

import (
	"bufio"
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

type walletMeta struct {
	Address string
	List    string
	Tier    string
	Smart   float64
	Bot     float64
}

type whaleTrade struct {
	TS      string  `json:"ts"`
	Wallet  string  `json:"wallet"`
	Label   string  `json:"label"`
	Side    string  `json:"side"`
	Price   float64 `json:"price"`
	Size    float64 `json:"size"`
	AssetID string  `json:"asset_id"`
	TradeID string  `json:"trade_id"`
	Action  string  `json:"action"`
	Reason  string  `json:"reason"`
	List    string  `json:"list"`
	Tier    string  `json:"tier"`
	Time    time.Time
}

type walletPerf struct {
	Meta               walletMeta
	Signals            int
	Closed             int
	Open               int
	Wins               int
	PnL                float64
	Capital            float64
	ROI                float64
	EventSignals       int
	EventClosed        int
	EventSettled       int
	EventMarked        int
	EventWins          int
	EventPnL           float64
	EventCapital       float64
	EventROI           float64
	ProvenEventSignals int
	ProvenEventWins    int
	ProvenEventPnL     float64
	ProvenEventCapital float64
	ProvenEventROI     float64
	Pending            int
	AssetCD            int
	EventCD            int
	DerivativeFiltered int
	CategoryFiltered   int
	OtherNoise         int
	EdgeSamples        int
	EdgeWins           int
	EdgeDeltaPPSum     float64
	EdgeAvgPP          float64
	Edge5mSamples      int
	Edge5mDeltaPPSum   float64
	Edge5mAvgPP        float64
	Edge15mSamples     int
	Edge15mDeltaPPSum  float64
	Edge15mAvgPP       float64
	Decision           string
	Reason             string
	Score              *walletdiscover.WalletScore
	SignalRank         float64
}

type performanceSnapshot struct {
	EventCappedByWallet       []snapshotStats `json:"event_capped_by_wallet"`
	EventCappedProvenByWallet []snapshotStats `json:"event_capped_proven_by_wallet"`
}

type snapshotStats struct {
	Wallet    string  `json:"wallet"`
	Signals   int     `json:"signals"`
	Closed    int     `json:"closed"`
	Settled   int     `json:"settled"`
	Marked    int     `json:"marked"`
	Wins      int     `json:"wins"`
	PnLUSD    float64 `json:"pnl_usd"`
	StakeUSD  float64 `json:"stake_usd"`
	ReturnPct float64 `json:"return_pct"`
}

type edgeSnapshot struct {
	Wallet      string  `json:"wallet"`
	Label       string  `json:"label"`
	HorizonSec  int64   `json:"horizon_sec"`
	DeltaPP     float64 `json:"delta_pp"`
	NotionalUSD float64 `json:"notional_usd"`
	SampleTime  string  `json:"sample_time"`
}

type edgePromoteParams struct {
	MinSamples            int
	MinAvgPP              float64
	MinWin                float64
	MaxBot                float64
	ReversalMin15mSamples int
	ReversalMax15mAvgPP   float64
	SevereMinSamples      int
	SevereMaxAvgPP        float64
	NegativeMinSamples    int
	NegativeMaxAvgPP      float64
	NegativeMaxWin        float64
}

type tapeFollowParams struct {
	MinSamples  int
	MinAvgPP    float64
	MinWin      float64
	Min5mAvgPP  float64
	Min15mAvgPP float64
	MaxBot      float64
}

func main() {
	logPath := flag.String("log", "db/journal/whale_trades.jsonl", "whale trades JSONL")
	pushPath := flag.String("push", "wallets.strategy-push.txt", "current push wallet file")
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet score JSON")
	performanceJSONPath := flag.String("performance_json", "reports/whale_performance.json", "optional whale-performance JSON with event-capped live stats")
	edgeSnapshotsPath := flag.String("edge_snapshots", "db/strategy_iteration/whale_edge_snapshots.jsonl", "optional whale-edge snapshot JSONL")
	reportPath := flag.String("report", "reports/wallet_maintenance.md", "maintenance report path")
	quarantinePath := flag.String("quarantine", "wallets.strategy-quarantine.txt", "quarantine wallet output")
	reviewNoisePath := flag.String("review_noise", "wallets.strategy-review-noise.txt", "wallets excluded from strategy generation because recent BUYs are suppressed/noisy")
	tapeCandidatesPath := flag.String("tape_candidates", "wallets.strategy-tape-candidates.txt", "edge-promoted tape candidate wallet output")
	tapeFollowPath := flag.String("tape_follow", "wallets.strategy-tape-follow.txt", "strict follow-ready tape wallet output")
	tapeReversalPath := flag.String("tape_reversal", "wallets.strategy-tape-reversal.txt", "tape wallets demoted for late-window edge reversal")
	stakeUSD := flag.Float64("stake", 10, "fixed stake used to evaluate each BUY signal")
	minSignals := flag.Int("min_signals", 5, "minimum closed live signals before demotion/promotion")
	minROI := flag.Float64("min_roi", 0, "minimum live ROI before quarantine")
	promoteROI := flag.Float64("promote_roi", 5, "minimum live ROI for upgrade suggestion")
	edgePromoteMinSamples := flag.Int("edge_promote_min_samples", 3, "minimum edge samples before promoting tape observation candidates")
	edgePromoteMinAvgPP := flag.Float64("edge_promote_min_avg_pp", 1, "minimum average edge delta in percentage points for tape observation candidates")
	edgePromoteMinWinRate := flag.Float64("edge_promote_min_win_rate", 60, "minimum edge win rate for tape observation candidates")
	edgePromoteMaxBot := flag.Float64("edge_promote_max_bot", 45, "maximum bot score before promoting tape observation candidates")
	edgePromoteReversalMin15mSamples := flag.Int("edge_promote_reversal_min_15m_samples", 2, "minimum 15m edge samples before demoting tape candidates for reversal")
	edgePromoteReversalMax15mAvgPP := flag.Float64("edge_promote_reversal_max_15m_avg_pp", -1, "demote tape candidates when 15m average edge is at or below this value")
	edgePromoteSevereMinSamples := flag.Int("edge_promote_severe_min_samples", 1, "minimum edge samples before demoting tape wallets for severe drawdown")
	edgePromoteSevereMaxAvgPP := flag.Float64("edge_promote_severe_max_avg_pp", -20, "demote tape wallets when average edge is at or below this severe drawdown value")
	edgePromoteNegativeMinSamples := flag.Int("edge_promote_negative_min_samples", 5, "minimum edge samples before demoting tape wallets for persistent negative edge")
	edgePromoteNegativeMaxAvgPP := flag.Float64("edge_promote_negative_max_avg_pp", -0.25, "demote tape wallets when average edge is at or below this persistent negative threshold")
	edgePromoteNegativeMaxWinRate := flag.Float64("edge_promote_negative_max_win_rate", 20, "demote tape wallets when edge win rate is at or below this persistent negative threshold")
	tapeFollowMinSamples := flag.Int("tape_follow_min_samples", 6, "minimum edge samples before tape candidate can become follow-ready")
	tapeFollowMinAvgPP := flag.Float64("tape_follow_min_avg_pp", 1.5, "minimum average edge delta in percentage points for follow-ready tape candidates")
	tapeFollowMinWinRate := flag.Float64("tape_follow_min_win_rate", 65, "minimum edge win rate for follow-ready tape candidates")
	tapeFollowMin5mAvgPP := flag.Float64("tape_follow_min_5m_avg_pp", 0.5, "minimum 5m average edge delta for follow-ready tape candidates")
	tapeFollowMin15mAvgPP := flag.Float64("tape_follow_min_15m_avg_pp", 0, "minimum 15m average edge delta for follow-ready tape candidates")
	tapeFollowMaxBot := flag.Float64("tape_follow_max_bot", 45, "maximum bot score for follow-ready tape candidates")
	eventMinSignals := flag.Int("event_min_signals", 1, "minimum event-capped live entries before event-based quarantine")
	eventQuarantineROI := flag.Float64("event_quarantine_roi", -30, "quarantine non-core wallets when event-capped ROI is below this threshold")
	noiseMinSuppressed := flag.Int("noise_min_suppressed", 3, "review non-core wallets with at least this many suppressed BUYs and no evaluated entries")
	flag.Parse()

	metas, err := loadWalletMetas(*pushPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: load push wallets: %v\n", err)
		os.Exit(1)
	}
	scores, err := loadScores(*scoresPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: load scores: %v\n", err)
		os.Exit(1)
	}
	trades, err := loadTrades(*logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: load whale log: %v\n", err)
		os.Exit(1)
	}

	perfs := evaluate(metas, scores, trades, *stakeUSD)
	if err := applyPerformanceSnapshot(perfs, *performanceJSONPath); err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: load performance json: %v\n", err)
		os.Exit(1)
	}
	if err := applyEdgeSnapshots(perfs, *edgeSnapshotsPath); err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: load edge snapshots: %v\n", err)
		os.Exit(1)
	}
	edgePromote := edgePromoteParams{
		MinSamples:            *edgePromoteMinSamples,
		MinAvgPP:              *edgePromoteMinAvgPP,
		MinWin:                *edgePromoteMinWinRate,
		MaxBot:                *edgePromoteMaxBot,
		ReversalMin15mSamples: *edgePromoteReversalMin15mSamples,
		ReversalMax15mAvgPP:   *edgePromoteReversalMax15mAvgPP,
		SevereMinSamples:      *edgePromoteSevereMinSamples,
		SevereMaxAvgPP:        *edgePromoteSevereMaxAvgPP,
		NegativeMinSamples:    *edgePromoteNegativeMinSamples,
		NegativeMaxAvgPP:      *edgePromoteNegativeMaxAvgPP,
		NegativeMaxWin:        *edgePromoteNegativeMaxWinRate,
	}
	tapeFollow := tapeFollowParams{
		MinSamples:  *tapeFollowMinSamples,
		MinAvgPP:    *tapeFollowMinAvgPP,
		MinWin:      *tapeFollowMinWinRate,
		Min5mAvgPP:  *tapeFollowMin5mAvgPP,
		Min15mAvgPP: *tapeFollowMin15mAvgPP,
		MaxBot:      *tapeFollowMaxBot,
	}
	decide(perfs, *minSignals, *minROI, *promoteROI, edgePromote, tapeFollow, *eventMinSignals, *eventQuarantineROI, *noiseMinSuppressed)
	quarantineTotal, reviewNoiseTotal, candidateTotal, followTotal, reversalTotal, err := writeOutputs(*reportPath, *quarantinePath, *reviewNoisePath, *tapeCandidatesPath, *tapeFollowPath, *tapeReversalPath, perfs, *minSignals, *minROI, *promoteROI, edgePromote, tapeFollow, *eventMinSignals, *eventQuarantineROI, *noiseMinSuppressed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet-maintain: write outputs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("wallet-maintain done: wallets=%d quarantine=%d review_noise=%d tape_candidates=%d tape_follow=%d tape_reversal=%d report=%s\n", len(perfs), quarantineTotal, reviewNoiseTotal, candidateTotal, followTotal, reversalTotal, *reportPath)
}

func applyPerformanceSnapshot(perfs []*walletPerf, path string) error {
	if path == "" {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var snap performanceSnapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return err
	}
	byWallet := map[string]snapshotStats{}
	for _, st := range snap.EventCappedByWallet {
		w := strings.ToLower(strings.TrimSpace(st.Wallet))
		if w != "" {
			byWallet[w] = st
		}
	}
	provenByWallet := map[string]snapshotStats{}
	for _, st := range snap.EventCappedProvenByWallet {
		w := strings.ToLower(strings.TrimSpace(st.Wallet))
		if w != "" {
			provenByWallet[w] = st
		}
	}
	for _, p := range perfs {
		st, ok := byWallet[strings.ToLower(p.Meta.Address)]
		if ok {
			p.EventSignals = st.Signals
			p.EventClosed = st.Closed
			p.EventSettled = st.Settled
			p.EventMarked = st.Marked
			p.EventWins = st.Wins
			p.EventPnL = st.PnLUSD
			p.EventCapital = st.StakeUSD
			p.EventROI = st.ReturnPct
		}
		proven, ok := provenByWallet[strings.ToLower(p.Meta.Address)]
		if ok {
			p.ProvenEventSignals = proven.Signals
			p.ProvenEventWins = proven.Wins
			p.ProvenEventPnL = proven.PnLUSD
			p.ProvenEventCapital = proven.StakeUSD
			p.ProvenEventROI = proven.ReturnPct
		} else if p.EventClosed+p.EventSettled > 0 && len(snap.EventCappedProvenByWallet) == 0 {
			p.ProvenEventSignals = p.EventClosed + p.EventSettled
		}
	}
	return nil
}

func applyEdgeSnapshots(perfs []*walletPerf, path string) error {
	if path == "" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	byWallet := map[string]*walletPerf{}
	for _, p := range perfs {
		byWallet[strings.ToLower(p.Meta.Address)] = p
	}

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var snap edgeSnapshot
		if err := json.Unmarshal(sc.Bytes(), &snap); err != nil {
			continue
		}
		w := strings.ToLower(strings.TrimSpace(snap.Wallet))
		p := byWallet[w]
		if p == nil {
			continue
		}
		p.EdgeSamples++
		if snap.DeltaPP > 0 {
			p.EdgeWins++
		}
		p.EdgeDeltaPPSum += snap.DeltaPP
		switch snap.HorizonSec {
		case int64((5 * time.Minute).Seconds()):
			p.Edge5mSamples++
			p.Edge5mDeltaPPSum += snap.DeltaPP
		case int64((15 * time.Minute).Seconds()):
			p.Edge15mSamples++
			p.Edge15mDeltaPPSum += snap.DeltaPP
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	for _, p := range perfs {
		if p.EdgeSamples > 0 {
			p.EdgeAvgPP = p.EdgeDeltaPPSum / float64(p.EdgeSamples)
		}
		if p.Edge5mSamples > 0 {
			p.Edge5mAvgPP = p.Edge5mDeltaPPSum / float64(p.Edge5mSamples)
		}
		if p.Edge15mSamples > 0 {
			p.Edge15mAvgPP = p.Edge15mDeltaPPSum / float64(p.Edge15mSamples)
		}
	}
	return nil
}

func loadWalletMetas(path string) (map[string]walletMeta, error) {
	if strings.Contains(path, ",") {
		out := map[string]walletMeta{}
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			metas, err := loadWalletMetas(part)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
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
		fields := strings.Fields(strings.TrimSpace(parts[0]))
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		meta := walletMeta{Address: addr}
		if len(parts) > 1 {
			for _, field := range strings.Fields(parts[1]) {
				k, v, ok := strings.Cut(field, "=")
				if !ok {
					continue
				}
				switch k {
				case "list":
					meta.List = v
				case "tier":
					meta.Tier = v
				case "smart":
					fmt.Sscanf(v, "%f", &meta.Smart)
				case "bot":
					fmt.Sscanf(v, "%f", &meta.Bot)
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

func loadScores(path string) (map[string]walletdiscover.WalletScore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []walletdiscover.WalletScore
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	out := map[string]walletdiscover.WalletScore{}
	for _, s := range list {
		out[strings.ToLower(s.Address)] = s
	}
	return out, nil
}

func loadTrades(path string) ([]whaleTrade, error) {
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
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
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

func evaluate(metas map[string]walletMeta, scores map[string]walletdiscover.WalletScore, trades []whaleTrade, stakeUSD float64) []*walletPerf {
	perfs := map[string]*walletPerf{}
	for addr, meta := range metas {
		score := scores[addr]
		scorePtr := &score
		if score.Address == "" {
			scorePtr = nil
		}
		perfs[addr] = &walletPerf{Meta: meta, Score: scorePtr}
	}

	sells := map[string][]whaleTrade{}
	for _, tr := range trades {
		if _, ok := perfs[tr.Wallet]; !ok {
			continue
		}
		if tr.Side == "SELL" {
			key := tr.Wallet + "|" + tr.AssetID
			sells[key] = append(sells[key], tr)
		}
	}

	seen := map[string]struct{}{}
	for _, buy := range trades {
		p := perfs[buy.Wallet]
		if p == nil || buy.Side != "BUY" {
			continue
		}
		switch buy.Action {
		case "pending_consensus":
			p.Pending++
		case "cooldown":
			p.AssetCD++
		case "event_cooldown":
			p.EventCD++
		case "alert", "followed":
		case "skip":
			switch strings.ToLower(strings.TrimSpace(buy.Reason)) {
			case "derivative_filtered":
				p.DerivativeFiltered++
			case "category_filtered":
				p.CategoryFiltered++
			default:
				p.OtherNoise++
			}
		default:
			p.OtherNoise++
		}
		if buy.Action != "alert" && buy.Action != "followed" {
			continue
		}
		key := signalKey(buy)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		p.Signals++
		exit := firstSellAfter(sells[buy.Wallet+"|"+buy.AssetID], buy.Time)
		if exit == nil {
			p.Open++
			continue
		}
		p.Closed++
		units := stakeUSD / buy.Price
		pnl := units * (exit.Price - buy.Price)
		p.PnL += pnl
		p.Capital += stakeUSD
		if pnl > 0 {
			p.Wins++
		}
	}

	out := make([]*walletPerf, 0, len(perfs))
	for _, p := range perfs {
		if p.Capital > 0 {
			p.ROI = p.PnL / p.Capital * 100
		}
		p.SignalRank = float64(p.Closed)*100 + p.ROI
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Meta.List != out[j].Meta.List {
			return listRank(out[i].Meta.List) < listRank(out[j].Meta.List)
		}
		if out[i].SignalRank != out[j].SignalRank {
			return out[i].SignalRank > out[j].SignalRank
		}
		return out[i].Meta.Address < out[j].Meta.Address
	})
	return out
}

func firstSellAfter(sells []whaleTrade, t time.Time) *whaleTrade {
	for _, sell := range sells {
		if sell.Time.After(t) || sell.Time.Equal(t) {
			return &sell
		}
	}
	return nil
}

func signalKey(tr whaleTrade) string {
	if tr.TradeID != "" {
		return tr.Wallet + "|" + tr.AssetID + "|" + tr.Side + "|" + tr.TradeID
	}
	return fmt.Sprintf("%s|%s|%s|%d|%.8f|%.8f", tr.Wallet, tr.AssetID, tr.Side, tr.Time.Unix(), tr.Price, tr.Size)
}

func decide(perfs []*walletPerf, minSignals int, minROI, promoteROI float64, edgePromote edgePromoteParams, tapeFollow tapeFollowParams, eventMinSignals int, eventQuarantineROI float64, noiseMinSuppressed int) {
	for _, p := range perfs {
		if p.Meta.List == "tape_reversal" {
			p.Decision = "review_tape_reversal"
			if severeTapeEdgeDrawdown(p, edgePromote) {
				p.Reason = fmt.Sprintf("retained reversal: severe edge drawdown %+0.2fpp over %d samples", p.EdgeAvgPP, p.EdgeSamples)
			} else if tapeCandidateReversed(p, edgePromote) {
				if persistentNegativeTapeEdge(p, edgePromote) {
					edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
					p.Reason = fmt.Sprintf("retained reversal: persistent negative edge %.1f%% win avg %+0.2fpp over %d samples", edgeWin, p.EdgeAvgPP, p.EdgeSamples)
				} else {
					p.Reason = fmt.Sprintf("retained reversal: 15m edge reversed %+0.2fpp over %d samples", p.Edge15mAvgPP, p.Edge15mSamples)
				}
			} else {
				p.Reason = "retained reversal risk"
			}
			continue
		}
		if p.Meta.List != "core" && p.ProvenEventSignals >= eventMinSignals && p.ProvenEventROI < eventQuarantineROI {
			p.Decision = "quarantine"
			p.Reason = fmt.Sprintf("proven event-capped ROI %.1f%% < %.1f%% over %d entries", p.ProvenEventROI, eventQuarantineROI, p.ProvenEventSignals)
			continue
		}
		suppressed := p.Pending + p.AssetCD + p.EventCD + p.DerivativeFiltered + p.CategoryFiltered + p.OtherNoise
		if p.Meta.List != "core" && p.EventSignals == 0 && noiseMinSuppressed > 0 && suppressed >= noiseMinSuppressed && !positiveEdgeNoiseExempt(p, edgePromote) {
			p.Decision = "review_noise"
			p.Reason = fmt.Sprintf("suppressed BUYs %d >= %d with no evaluated entries", suppressed, noiseMinSuppressed)
			continue
		}
		if isTapeEdgeList(p.Meta.List) && tapeEdgeReviewable(p, edgePromote) {
			edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
			if tapeCandidateReversed(p, edgePromote) {
				p.Decision = "review_tape_reversal"
				if severeTapeEdgeDrawdown(p, edgePromote) {
					p.Reason = fmt.Sprintf("severe edge drawdown %+0.2fpp over %d samples", p.EdgeAvgPP, p.EdgeSamples)
				} else if persistentNegativeTapeEdge(p, edgePromote) {
					p.Reason = fmt.Sprintf("persistent negative edge %.1f%% win avg %+0.2fpp over %d samples", edgeWin, p.EdgeAvgPP, p.EdgeSamples)
				} else {
					p.Reason = fmt.Sprintf("15m edge reversed %+0.2fpp over %d samples", p.Edge15mAvgPP, p.Edge15mSamples)
				}
				continue
			}
			if tapeFollowReady(p, tapeFollow, edgeWin) {
				p.Decision = "promote_tape_follow"
				p.Reason = fmt.Sprintf("follow-ready edge %.1f%% win avg %+0.2fpp 5m %+0.2fpp 15m %+0.2fpp bot %.1f", edgeWin, p.EdgeAvgPP, p.Edge5mAvgPP, p.Edge15mAvgPP, walletBot(p))
				continue
			}
			if p.EdgeAvgPP >= edgePromote.MinAvgPP && edgeWin >= edgePromote.MinWin {
				if !tapeCandidateBotAllowed(p, edgePromote) {
					p.Decision = "learning"
					p.Reason = fmt.Sprintf("edge positive but bot %.1f > %.1f", walletBot(p), edgePromote.MaxBot)
					continue
				}
				p.Decision = "promote_tape_candidate"
				p.Reason = fmt.Sprintf("edge %.1f%% win and avg %+0.2fpp over %d samples bot %.1f", edgeWin, p.EdgeAvgPP, p.EdgeSamples, walletBot(p))
				continue
			}
		}
		if p.Closed < minSignals {
			p.Decision = "learning"
			p.Reason = fmt.Sprintf("closed signals %d < %d", p.Closed, minSignals)
			continue
		}
		if p.ROI < minROI {
			p.Decision = "quarantine"
			p.Reason = fmt.Sprintf("live ROI %.1f%% < %.1f%%", p.ROI, minROI)
			continue
		}
		switch p.Meta.List {
		case "scout", "target", "flow":
			if p.ROI >= promoteROI {
				p.Decision = "promote_watch"
				p.Reason = fmt.Sprintf("%s live ROI %.1f%% >= %.1f%%", p.Meta.List, p.ROI, promoteROI)
			} else {
				p.Decision = "keep_" + p.Meta.List
				p.Reason = "non-negative live ROI, needs stronger edge"
			}
		case "watch":
			if p.ROI >= promoteROI && p.Score != nil && p.Score.Stats.CopyClosedTrades >= 8 && p.Score.Stats.CopyROI >= 10 {
				p.Decision = "promote_core_candidate"
				p.Reason = fmt.Sprintf("watch live ROI %.1f%% and copyROI %.1f%%", p.ROI, p.Score.Stats.CopyROI)
			} else {
				p.Decision = "keep_watch"
				p.Reason = "non-negative live ROI"
			}
		case "sports":
			if p.ROI >= promoteROI && p.Score != nil && p.Score.Stats.TargetCopyClosed >= 5 && p.Score.Stats.TargetCopyROI >= 5 {
				p.Decision = "promote_core_candidate"
				p.Reason = fmt.Sprintf("sports live ROI %.1f%% and targetCopyROI %.1f%%", p.ROI, p.Score.Stats.TargetCopyROI)
			} else {
				p.Decision = "keep_sports"
				p.Reason = "non-negative live ROI"
			}
		default:
			p.Decision = "keep"
			p.Reason = "non-negative live ROI"
		}
	}
}

func positiveEdgeNoiseExempt(p *walletPerf, params edgePromoteParams) bool {
	if p == nil || params.MinSamples <= 0 {
		return false
	}
	edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
	return p.EdgeSamples >= params.MinSamples &&
		p.EdgeAvgPP >= params.MinAvgPP &&
		edgeWin >= params.MinWin &&
		walletBot(p) <= params.MaxBot
}

func writeOutputs(reportPath, quarantinePath, reviewNoisePath, tapeCandidatesPath, tapeFollowPath, tapeReversalPath string, perfs []*walletPerf, minSignals int, minROI, promoteROI float64, edgePromote edgePromoteParams, tapeFollow tapeFollowParams, eventMinSignals int, eventQuarantineROI float64, noiseMinSuppressed int) (int, int, int, int, int, error) {
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(quarantinePath), 0o755); err != nil && filepath.Dir(quarantinePath) != "." {
		return 0, 0, 0, 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(reviewNoisePath), 0o755); err != nil && filepath.Dir(reviewNoisePath) != "." {
		return 0, 0, 0, 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(tapeCandidatesPath), 0o755); err != nil && filepath.Dir(tapeCandidatesPath) != "." {
		return 0, 0, 0, 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(tapeFollowPath), 0o755); err != nil && filepath.Dir(tapeFollowPath) != "." {
		return 0, 0, 0, 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(tapeReversalPath), 0o755); err != nil && filepath.Dir(tapeReversalPath) != "." {
		return 0, 0, 0, 0, 0, err
	}

	perfWallets := map[string]struct{}{}
	for _, p := range perfs {
		perfWallets[strings.ToLower(p.Meta.Address)] = struct{}{}
	}
	retainedQuarantine := existingWalletLinesOutsideCurrentPush(quarantinePath, perfWallets)
	retainedReviewNoise := existingReviewNoiseOutsideCurrentPush(reviewNoisePath, perfWallets, edgePromote)
	var qLines []string
	var reviewNoiseLines []string
	var tapeCandidateLines []string
	var tapeFollowLines []string
	var tapeReversalLines []string
	var b strings.Builder
	fmt.Fprintf(&b, "# Wallet Maintenance Report\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Min closed signals: %d\n", minSignals)
	fmt.Fprintf(&b, "- Quarantine ROI: < %.1f%%\n", minROI)
	fmt.Fprintf(&b, "- Promotion ROI: >= %.1f%%\n\n", promoteROI)
	fmt.Fprintf(&b, "- Tape edge promotion: samples >= %d, avg edge >= %.2fpp, win >= %.1f%%, bot <= %.1f\n\n", edgePromote.MinSamples, edgePromote.MinAvgPP, edgePromote.MinWin, edgePromote.MaxBot)
	fmt.Fprintf(&b, "- Tape reversal review: 15m samples >= %d and 15m avg <= %.2fpp\n\n", edgePromote.ReversalMin15mSamples, edgePromote.ReversalMax15mAvgPP)
	fmt.Fprintf(&b, "- Tape severe drawdown review: samples >= %d and avg edge <= %.2fpp\n\n", edgePromote.SevereMinSamples, edgePromote.SevereMaxAvgPP)
	fmt.Fprintf(&b, "- Tape persistent negative review: samples >= %d, win <= %.1f%%, avg edge <= %.2fpp\n\n", edgePromote.NegativeMinSamples, edgePromote.NegativeMaxWin, edgePromote.NegativeMaxAvgPP)
	fmt.Fprintf(&b, "- Tape follow-ready: samples >= %d, avg edge >= %.2fpp, win >= %.1f%%, 5m avg >= %.2fpp, 15m avg >= %.2fpp, bot <= %.1f\n\n",
		tapeFollow.MinSamples, tapeFollow.MinAvgPP, tapeFollow.MinWin, tapeFollow.Min5mAvgPP, tapeFollow.Min15mAvgPP, tapeFollow.MaxBot)
	fmt.Fprintf(&b, "- Proven event-capped quarantine: non-core wallets with proven entries >= %d and ROI < %.1f%%\n\n", eventMinSignals, eventQuarantineROI)
	fmt.Fprintf(&b, "- Noise review: non-core wallets with suppressed BUYs >= %d and no evaluated entries\n\n", noiseMinSuppressed)
	fmt.Fprintf(&b, "## Summary\n\n")
	for _, decision := range []string{"quarantine", "review_noise", "review_tape_reversal", "promote_core_candidate", "promote_watch", "promote_tape_follow", "promote_tape_candidate", "keep", "keep_watch", "keep_sports", "keep_scout", "keep_target", "keep_flow", "learning"} {
		fmt.Fprintf(&b, "- %s: %d\n", decision, countDecision(perfs, decision))
	}
	fmt.Fprintf(&b, "- retained_existing_quarantine: %d\n", len(retainedQuarantine))
	fmt.Fprintf(&b, "- retained_existing_review_noise: %d\n", len(retainedReviewNoise))
	fmt.Fprintf(&b, "\n## Wallets\n\n")
	fmt.Fprintf(&b, "| Wallet | List | Tier | Decision | Closed | Open | Win | LiveROI | LivePnL | EventEntries | EventROI | EventPnL | ProvenEventEntries | ProvenEventROI | ProvenEventPnL | EdgeN | EdgeWin | EdgeAvgPP | Edge5mPP | Edge15mPP | Pending | AssetCD | EventCD | DerivSkip | CatSkip | OtherNoise | CopyROI | Bot | Reason |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, p := range perfs {
		copyROI, bot := 0.0, p.Meta.Bot
		if p.Score != nil {
			copyROI = p.Score.Stats.CopyROI
			bot = p.Score.BotScore
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %d | %d | %.1f%% | %.1f%% | $%+.2f | %d | %.1f%% | $%+.2f | %d | %.1f%% | $%+.2f | %d | %.1f%% | %+.2f | %+.2f | %+.2f | %d | %d | %d | %d | %d | %d | %.1f%% | %.1f | %s |\n",
			shortAddr(p.Meta.Address), p.Meta.List, p.Meta.Tier, p.Decision, p.Closed, p.Open,
			pct(float64(p.Wins), float64(p.Closed)), p.ROI, p.PnL, p.EventSignals, p.EventROI, p.EventPnL,
			p.ProvenEventSignals, p.ProvenEventROI, p.ProvenEventPnL, p.EdgeSamples,
			pct(float64(p.EdgeWins), float64(p.EdgeSamples)), p.EdgeAvgPP, p.Edge5mAvgPP, p.Edge15mAvgPP,
			p.Pending, p.AssetCD, p.EventCD, p.DerivativeFiltered, p.CategoryFiltered, p.OtherNoise, copyROI, bot, p.Reason)
		if p.Decision == "quarantine" {
			qLines = append(qLines, fmt.Sprintf("%s # list=%s tier=%s liveROI=%.1f%% provenEventROI=%.1f%% provenEventT=%d closed=%d pnl=$%+.2f reason=%s",
				p.Meta.Address, p.Meta.List, p.Meta.Tier, p.ROI, p.ProvenEventROI, p.ProvenEventSignals, p.Closed, p.PnL, p.Reason))
		}
		if p.Decision == "review_noise" {
			reviewNoiseLines = append(reviewNoiseLines, fmt.Sprintf("%s # list=%s tier=%s suppressed=%d pending=%d assetCD=%d eventCD=%d derivSkip=%d catSkip=%d otherNoise=%d edgeN=%d edgeAvgPP=%+.2f bot=%.1f copyROI=%.1f%% reason=%s",
				p.Meta.Address, p.Meta.List, p.Meta.Tier,
				p.Pending+p.AssetCD+p.EventCD+p.DerivativeFiltered+p.CategoryFiltered+p.OtherNoise,
				p.Pending, p.AssetCD, p.EventCD, p.DerivativeFiltered, p.CategoryFiltered, p.OtherNoise,
				p.EdgeSamples, p.EdgeAvgPP, bot, copyROI, p.Reason))
		}
		if p.Decision == "promote_tape_candidate" || p.Decision == "promote_tape_follow" {
			edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
			tapeCandidateLines = append(tapeCandidateLines, fmt.Sprintf("%s # list=tape_candidate source=%s tier=%s edgeN=%d edgeWin=%.1f%% edgeAvgPP=%+.2f edge5mPP=%+.2f edge15mPP=%+.2f bot=%.1f copyROI=%.1f%% reason=%s",
				p.Meta.Address, p.Meta.List, p.Meta.Tier, p.EdgeSamples, edgeWin, p.EdgeAvgPP, p.Edge5mAvgPP, p.Edge15mAvgPP, bot, copyROI, p.Reason))
		}
		if p.Decision == "promote_tape_follow" {
			edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
			tapeFollowLines = append(tapeFollowLines, fmt.Sprintf("%s # list=tape_follow source=%s tier=%s edgeN=%d edgeWin=%.1f%% edgeAvgPP=%+.2f edge5mPP=%+.2f edge15mPP=%+.2f bot=%.1f copyROI=%.1f%% reason=%s",
				p.Meta.Address, p.Meta.List, p.Meta.Tier, p.EdgeSamples, edgeWin, p.EdgeAvgPP, p.Edge5mAvgPP, p.Edge15mAvgPP, bot, copyROI, p.Reason))
		}
		if p.Decision == "review_tape_reversal" {
			edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
			tapeReversalLines = append(tapeReversalLines, fmt.Sprintf("%s # list=tape_reversal source=%s tier=%s edgeN=%d edgeWin=%.1f%% edgeAvgPP=%+.2f edge5mPP=%+.2f edge15mPP=%+.2f bot=%.1f copyROI=%.1f%% reason=%s",
				p.Meta.Address, p.Meta.List, p.Meta.Tier, p.EdgeSamples, edgeWin, p.EdgeAvgPP, p.Edge5mAvgPP, p.Edge15mAvgPP, bot, copyROI, p.Reason))
		}
	}
	qLines = append(qLines, retainedQuarantine...)
	reviewNoiseLines = append(reviewNoiseLines, retainedReviewNoise...)
	if err := os.WriteFile(reportPath, []byte(b.String()), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	sort.Strings(qLines)
	content := strings.Join(qLines, "\n")
	if content != "" {
		content += "\n"
	}
	if err := os.WriteFile(quarantinePath, []byte(content), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	sort.Strings(reviewNoiseLines)
	reviewNoiseContent := strings.Join(reviewNoiseLines, "\n")
	if reviewNoiseContent != "" {
		reviewNoiseContent += "\n"
	}
	if err := os.WriteFile(reviewNoisePath, []byte(reviewNoiseContent), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	sort.Strings(tapeCandidateLines)
	candidateContent := strings.Join(tapeCandidateLines, "\n")
	if candidateContent != "" {
		candidateContent += "\n"
	}
	if err := os.WriteFile(tapeCandidatesPath, []byte(candidateContent), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	sort.Strings(tapeFollowLines)
	followContent := strings.Join(tapeFollowLines, "\n")
	if followContent != "" {
		followContent += "\n"
	}
	if err := os.WriteFile(tapeFollowPath, []byte(followContent), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	sort.Strings(tapeReversalLines)
	reversalContent := strings.Join(tapeReversalLines, "\n")
	if reversalContent != "" {
		reversalContent += "\n"
	}
	if err := os.WriteFile(tapeReversalPath, []byte(reversalContent), 0o644); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	return len(qLines), len(reviewNoiseLines), len(tapeCandidateLines), len(tapeFollowLines), len(tapeReversalLines), nil
}

func existingWalletLinesOutsideCurrentPush(path string, current map[string]struct{}) []string {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []string
	seen := map[string]struct{}{}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(strings.Split(line, "#")[0])
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		if _, ok := current[addr]; ok {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		out = append(out, line)
	}
	return out
}

func existingReviewNoiseOutsideCurrentPush(path string, current map[string]struct{}, edgePromote edgePromoteParams) []string {
	lines := existingWalletLinesOutsideCurrentPush(path, current)
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if reviewNoiseLineHasStrongPositiveEdge(line, edgePromote) {
			continue
		}
		out = append(out, line)
	}
	return out
}

func reviewNoiseLineHasStrongPositiveEdge(line string, edgePromote edgePromoteParams) bool {
	if edgePromote.MinSamples <= 0 {
		return false
	}
	_, comment, ok := strings.Cut(line, "#")
	if !ok {
		return false
	}
	fields := map[string]string{}
	for _, field := range strings.Fields(comment) {
		k, v, ok := strings.Cut(field, "=")
		if ok {
			fields[k] = v
		}
	}
	edgeN := parseLineInt(fields["edgeN"])
	edgeAvg := parseLineFloat(fields["edgeAvgPP"])
	bot := parseLineFloat(fields["bot"])
	copyROI := parseLineFloat(fields["copyROI"])
	return edgeN >= edgePromote.MinSamples &&
		edgeAvg >= edgePromote.MinAvgPP &&
		bot > 0 && bot <= edgePromote.MaxBot &&
		copyROI > 0
}

func parseLineInt(raw string) int {
	var out int
	_, _ = fmt.Sscanf(strings.TrimSpace(raw), "%d", &out)
	return out
}

func parseLineFloat(raw string) float64 {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "$")
	raw = strings.TrimSuffix(raw, "%")
	raw = strings.ReplaceAll(raw, ",", "")
	var out float64
	_, _ = fmt.Sscanf(raw, "%f", &out)
	return out
}

func countDecision(perfs []*walletPerf, decision string) int {
	n := 0
	for _, p := range perfs {
		if p.Decision == decision {
			n++
		}
	}
	return n
}

func listRank(list string) int {
	switch list {
	case "core":
		return 0
	case "watch":
		return 1
	case "sports":
		return 2
	case "scout":
		return 3
	case "target":
		return 4
	case "flow":
		return 5
	case "tape_observe":
		return 6
	case "tape_probation":
		return 7
	case "tape_candidate":
		return 8
	case "tape_follow":
		return 9
	case "tape_reversal":
		return 10
	default:
		return 11
	}
}

func isTapeEdgeList(list string) bool {
	switch list {
	case "tape", "tape_observe", "tape_probation", "tape_candidate", "tape_follow", "tape_reversal":
		return true
	default:
		return false
	}
}

func tapeFollowReady(p *walletPerf, params tapeFollowParams, edgeWin float64) bool {
	if params.MinSamples <= 0 || p.EdgeSamples < params.MinSamples {
		return false
	}
	if p.EdgeAvgPP < params.MinAvgPP || edgeWin < params.MinWin {
		return false
	}
	if p.Edge5mSamples == 0 || p.Edge5mAvgPP < params.Min5mAvgPP {
		return false
	}
	if p.Edge15mSamples == 0 || p.Edge15mAvgPP < params.Min15mAvgPP {
		return false
	}
	if walletBot(p) > params.MaxBot {
		return false
	}
	return true
}

func tapeCandidateBotAllowed(p *walletPerf, params edgePromoteParams) bool {
	if params.MaxBot <= 0 {
		return true
	}
	return walletBot(p) <= params.MaxBot
}

func tapeEdgeReviewable(p *walletPerf, params edgePromoteParams) bool {
	if params.MinSamples > 0 && p.EdgeSamples >= params.MinSamples {
		return true
	}
	return severeTapeEdgeDrawdown(p, params)
}

func tapeCandidateReversed(p *walletPerf, params edgePromoteParams) bool {
	if severeTapeEdgeDrawdown(p, params) {
		return true
	}
	if persistentNegativeTapeEdge(p, params) {
		return true
	}
	if params.ReversalMin15mSamples <= 0 {
		return false
	}
	return p.Edge15mSamples >= params.ReversalMin15mSamples && p.Edge15mAvgPP <= params.ReversalMax15mAvgPP
}

func severeTapeEdgeDrawdown(p *walletPerf, params edgePromoteParams) bool {
	if params.SevereMinSamples <= 0 {
		return false
	}
	return p.EdgeSamples >= params.SevereMinSamples && p.EdgeAvgPP <= params.SevereMaxAvgPP
}

func persistentNegativeTapeEdge(p *walletPerf, params edgePromoteParams) bool {
	if params.NegativeMinSamples <= 0 || p.EdgeSamples < params.NegativeMinSamples {
		return false
	}
	edgeWin := pct(float64(p.EdgeWins), float64(p.EdgeSamples))
	return edgeWin <= params.NegativeMaxWin && p.EdgeAvgPP <= params.NegativeMaxAvgPP
}

func walletBot(p *walletPerf) float64 {
	if p == nil {
		return 0
	}
	if p.Score != nil {
		return p.Score.BotScore
	}
	return p.Meta.Bot
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
