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
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/config"
	"github.com/15529214579/polymarket-go/internal/notify"
)

type tapeTrade struct {
	Time          time.Time `json:"time"`
	Wallet        string    `json:"wallet"`
	Side          string    `json:"side"`
	Notional      float64   `json:"notional"`
	Price         float64   `json:"price"`
	Size          float64   `json:"size"`
	Outcome       string    `json:"outcome"`
	Market        string    `json:"market"`
	Slug          string    `json:"slug"`
	Category      string    `json:"category"`
	Asset         string    `json:"asset"`
	Transaction   string    `json:"transaction"`
	KnownList     string    `json:"known_list,omitempty"`
	Tier          string    `json:"tier,omitempty"`
	Smart         float64   `json:"smart,omitempty"`
	Bot           float64   `json:"bot,omitempty"`
	TargetCopyROI float64   `json:"target_copy_roi,omitempty"`
	TargetCopyT   int       `json:"target_copy_t,omitempty"`
	Participants  []string  `json:"participants,omitempty"`
}

type walletStatus struct {
	List    string
	Tier    string
	Smart   float64
	Bot     float64
	Reason  string
	Source  string
	EdgeN   int
	EdgeAvg float64
	Edge5m  float64
	Edge15m float64
}

type sentState struct {
	Sent    map[string]string     `json:"sent"`
	Details map[string]sentDetail `json:"details,omitempty"`
}

type sentDetail struct {
	Mode      string `json:"mode,omitempty"`
	KnownList string `json:"known_list,omitempty"`
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

type modePolicyDecision struct {
	GeneratedAt time.Time                  `json:"generated_at"`
	Modes       []modePolicyDecisionBucket `json:"modes"`
}

type modePolicyDecisionBucket struct {
	Key    string `json:"key"`
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type burstCandidate struct {
	Wallet        string
	Asset         string
	Mode          string
	KnownList     string
	Tier          string
	Bot           float64
	Category      string
	Market        string
	Outcome       string
	Slug          string
	Trades        int
	TotalNotional float64
	LastNotional  float64
	LastPrice     float64
	FirstTime     time.Time
	LastTime      time.Time
	AlreadySent   bool
	Reason        string
}

type consensusCandidate struct {
	Asset         string
	Category      string
	Market        string
	Outcome       string
	Slug          string
	Trades        int
	Wallets       int
	Participants  []string
	TotalNotional float64
	TotalSize     float64
	VWAP          float64
	LastNotional  float64
	FirstTime     time.Time
	LastTime      time.Time
	Tier          string
	Bot           float64
	AlreadySent   bool
	Reason        string
}

type unknownFlowCandidate struct {
	Wallet        string
	Category      string
	Market        string
	Outcome       string
	Slug          string
	Asset         string
	Trades        int
	Markets       int
	TotalNotional float64
	LastNotional  float64
	LastPrice     float64
	FirstTime     time.Time
	LastTime      time.Time
	Tier          string
	Bot           float64
	AlreadyLogged bool
	Reason        string
}

type diagnosticRow struct {
	Trade     tapeTrade
	Mode      string
	Eligible  bool
	Alertable bool
	Reason    string
	Age       time.Duration
	Sent      bool
}

type alertPolicy struct {
	AllowedModes            map[string]struct{}
	MinNotional             float64
	ObserveMinNotional      float64
	ObserveBurstMinNotional float64
	ObserveMaxBot           float64
	ObserveRequireKnown     bool
	ObserveMinTier          string
	InsiderMinNotional      float64
	InsiderMaxBot           float64
	PositionCooldown        time.Duration
	RepeatMinNotional       float64
	EdgeBlocks              map[string]string
	EdgeHot                 map[string]string
	RequirePositiveEdge     map[string]struct{}
	EdgeMetrics             map[string]map[int64]*edgeStats
	EdgeHotMinNotional      float64
	EdgeHotMaxBot           float64
	EdgeHotMinSamples       int
	EdgeHotMinAvgPP         float64
	EdgeHotMinWinRate       float64
	EdgeHotMin5mAvgPP       float64
	EdgeHotMin15mAvgPP      float64
	EdgeHotMax1hNegPP       float64
	ModeBlocks              map[string]string
	ModeActions             map[string]string
	ModeMinAction           string
	AllowReversalAlerts     bool
	ConsensusAlerts         bool
	ConsensusMinNotional    float64
	ConsensusMinWallets     int
	ConsensusMaxBot         float64
	UnknownFlowMinNotional  float64
	UnknownFlowMinMarkets   int
	UnknownFlowMaxBot       float64
	SeedFlowMinNotional     float64
	SeedFlowMinMarkets      int
	ScoredFlowMinNotional   float64
	ScoredFlowMinMarkets    int
	ScoredFlowMaxBot        float64
	ScoredFlowMinTier       string
}

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
	Participants  []string  `json:"participants,omitempty"`
}

func main() {
	tapePath := flag.String("tape", "db/strategy_iteration/sports_tape.jsonl", "sports tape JSONL")
	statePath := flag.String("state", "db/strategy_iteration/sports_tape_alert_sent.json", "sent-alert state JSON")
	sentLogPath := flag.String("sent_log", "db/strategy_iteration/sports_tape_alerts.jsonl", "append-only alert audit JSONL")
	shadowLogPath := flag.String("shadow_log", "", "optional append-only dry-run audit JSONL for stale OBSERVE-BURST candidates")
	diagnosticPath := flag.String("diagnostic_report", "reports/sports_alert_candidates.md", "markdown report explaining current sports alert candidates and blockers")
	walletStatusPath := flag.String("wallet_statuses", "wallets.strategy-push.txt,wallets.strategy-tape-follow.txt,wallets.strategy-tape-candidates.txt,wallets.strategy-tape-probation.txt,wallets.strategy-tape-reversal.txt,wallets.strategy-review-noise.txt,wallets.strategy-tape-observe.txt", "comma-separated wallet status files used to label alerts")
	edgeSnapshotsPath := flag.String("edge_snapshots", "db/strategy_iteration/whale_edge_snapshots.jsonl", "optional whale-edge snapshot JSONL used to suppress wallets with negative measured edge")
	modePolicyPath := flag.String("mode_policy", "", "optional sports policy decision JSON; modes with CUT/PROBATION are suppressed before Telegram")
	modePolicyMaxAge := flag.Duration("mode_policy_max_age", 2*time.Hour, "ignore mode_policy when generated_at is older than this; 0 disables age check")
	modePolicyMinAction := flag.String("mode_policy_min_action", "COLLECT_POSITIVE", "minimum policy action required before Telegram; lower or missing modes are shadow-only when mode_policy is active")
	envPath := flag.String("env", ".env.local", "dotenv file with Telegram credentials")
	minNotional := flag.Float64("min_notional", 5000, "minimum BUY notional to alert")
	observeMinNotional := flag.Float64("observe_min_notional", 0, "minimum BUY notional for OBSERVE-mode raw whale alerts; 0 disables OBSERVE alerts")
	observeBurstMinNotional := flag.Float64("observe_burst_min_notional", 0, "minimum same-wallet same-asset cumulative BUY notional for OBSERVE-BURST alerts; 0 disables")
	observeMaxBot := flag.Float64("observe_max_bot", 45, "maximum bot score for OBSERVE-mode raw whale alerts; 0 disables bot-score filtering")
	observeRequireKnown := flag.Bool("observe_require_known", true, "require OBSERVE-mode raw whale alerts to have wallet score/list metadata")
	observeMinTier := flag.String("observe_min_tier", "", "minimum wallet tier for OBSERVE-mode raw whale alerts; empty disables tier filtering")
	insiderMinNotional := flag.Float64("insider_min_notional", 0, "minimum BUY notional for low-sample insider OBSERVE alerts; 0 disables insider-scout alerts")
	insiderMaxBot := flag.Float64("insider_max_bot", 35, "maximum bot score for insider OBSERVE alerts when bot metadata is available; 0 disables bot-score filtering")
	edgeHotMinNotional := flag.Float64("edge_hot_min_notional", 1000, "minimum BUY notional for EDGE-HOT wallets")
	edgeHotMaxBot := flag.Float64("edge_hot_max_bot", 45, "maximum bot score for EDGE-HOT wallets; 0 disables bot-score filtering")
	edgeHotMinSamples := flag.Int("edge_hot_min_samples", 2, "minimum non-zero horizon edge samples needed for EDGE-HOT alerts")
	edgeHotMinAvg := flag.Float64("edge_hot_min_avg_pp", 2, "minimum non-zero horizon average edge in pp for EDGE-HOT alerts")
	edgeHotMinWinRate := flag.Float64("edge_hot_min_win_rate", 60, "minimum non-zero horizon win rate for EDGE-HOT alerts")
	edgeHotMin5mAvg := flag.Float64("edge_hot_min_5m_avg_pp", 0.5, "minimum 5m average edge in pp for EDGE-HOT alerts")
	edgeHotMin15mAvg := flag.Float64("edge_hot_min_15m_avg_pp", 0, "minimum 15m average edge in pp for EDGE-HOT alerts")
	edgeHotMax1hNeg := flag.Float64("edge_hot_max_1h_neg_pp", -5, "do not mark EDGE-HOT when 1h edge is at or below this pp value")
	edgeBlock15mSamples := flag.Int("edge_block_15m_samples", 2, "minimum 15m edge samples needed to suppress a wallet")
	edgeBlock15mMaxAvg := flag.Float64("edge_block_15m_max_avg_pp", -1, "suppress wallet when 15m edge average is at or below this pp value")
	edgeBlock1hSamples := flag.Int("edge_block_1h_samples", 1, "minimum 1h edge samples needed to suppress a wallet")
	edgeBlock1hMaxAvg := flag.Float64("edge_block_1h_max_avg_pp", -5, "suppress wallet when 1h edge average is at or below this pp value")
	requirePositiveEdgeRaw := flag.String("require_positive_edge_modes", "CANDIDATE,PROBATION", "comma-separated alert modes that require edge-hot positive edge evidence before sending")
	modesRaw := flag.String("modes", "FOLLOW-READY,CANDIDATE,PROBATION,EDGE-HOT,FLOW-SCOUT", "comma-separated alert modes to send at min_notional")
	maxAge := flag.Duration("max_age", 10*time.Minute, "maximum trade age to alert")
	diagnosticAge := flag.Duration("diagnostic_age", 6*time.Hour, "maximum trade age to include in the candidate diagnostic report")
	burstWindow := flag.Duration("burst_window", 15*time.Minute, "time window for diagnostic same-wallet same-asset accumulation bursts")
	burstMinNotional := flag.Float64("burst_min_notional", 5000, "minimum cumulative BUY notional for diagnostic accumulation bursts")
	burstMinTrades := flag.Int("burst_min_trades", 2, "minimum BUY trades for diagnostic accumulation bursts")
	burstMinLegNotional := flag.Float64("burst_min_leg_notional", 1000, "minimum individual BUY notional included in diagnostic accumulation bursts")
	consensusAlerts := flag.Bool("consensus_alerts", true, "send OBSERVE-style alerts for fresh cross-wallet same-asset consensus bursts")
	consensusMinNotional := flag.Float64("consensus_min_notional", 10000, "minimum cumulative BUY notional for consensus burst alerts")
	consensusMinWallets := flag.Int("consensus_min_wallets", 2, "minimum unique wallets for consensus burst alerts")
	consensusMaxBot := flag.Float64("consensus_max_bot", 60, "maximum bot score for wallets included in consensus bursts; 0 disables bot-score filtering")
	unknownFlowMinNotional := flag.Float64("unknown_flow_min_notional", 6000, "minimum cumulative BUY notional for shadow-only unknown-wallet multi-market flow; 0 disables")
	unknownFlowMinMarkets := flag.Int("unknown_flow_min_markets", 2, "minimum distinct markets/assets for shadow-only unknown-wallet flow")
	unknownFlowMaxBot := flag.Float64("unknown_flow_max_bot", 45, "maximum bot score for shadow-only unknown-wallet flow when metadata exists; 0 disables")
	seedFlowMinNotional := flag.Float64("seed_flow_min_notional", 3000, "minimum cumulative BUY notional for lower-threshold shadow-only unknown-wallet multi-market seed flow; 0 disables")
	seedFlowMinMarkets := flag.Int("seed_flow_min_markets", 2, "minimum distinct markets/assets for lower-threshold seed flow")
	scoredFlowMinNotional := flag.Float64("scored_flow_min_notional", 6000, "minimum cumulative BUY notional for shadow-only scored-wallet multi-market flow; 0 disables")
	scoredFlowMinMarkets := flag.Int("scored_flow_min_markets", 2, "minimum distinct markets/assets for shadow-only scored-wallet flow")
	scoredFlowMaxBot := flag.Float64("scored_flow_max_bot", 35, "maximum bot score for shadow-only scored-wallet flow; 0 disables")
	scoredFlowMinTier := flag.String("scored_flow_min_tier", "B", "minimum wallet tier for shadow-only scored-wallet flow")
	consensusMaxAge := flag.Duration("consensus_max_age", 15*time.Minute, "maximum last-leg age for consensus burst alerts; 0 uses max_age")
	positionCooldown := flag.Duration("position_cooldown", 30*time.Minute, "suppress repeat alerts for the same wallet+asset within this window; 0 disables")
	repeatMinNotional := flag.Float64("repeat_min_notional", 25000, "same wallet+asset alerts at or above this notional bypass position_cooldown; 0 disables bypass")
	maxAlerts := flag.Int("max_alerts", 5, "maximum alerts per run")
	backfill := flag.Bool("backfill", false, "send existing qualifying trades on first run")
	dryRun := flag.Bool("dry_run", false, "print alerts instead of sending Telegram")
	flag.Parse()

	_ = config.LoadDotEnv(*envPath)
	trades, err := loadTapeTrades(*tapePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: load tape: %v\n", err)
		os.Exit(1)
	}
	statuses, err := loadWalletStatuses(*walletStatusPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: load wallet statuses: %v\n", err)
		os.Exit(1)
	}
	applyWalletStatuses(trades, statuses)
	state, existed, err := loadState(*statePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: load state: %v\n", err)
		os.Exit(1)
	}
	if err := reconcileAlertLog(*sentLogPath, trades, state); err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: reconcile alert log: %v\n", err)
		os.Exit(1)
	}
	edgeMetrics, err := loadEdgeMetrics(*edgeSnapshotsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: load edge snapshots: %v\n", err)
		os.Exit(1)
	}
	edgeBlocks := negativeEdgeBlocks(edgeMetrics, *edgeBlock15mSamples, *edgeBlock15mMaxAvg, *edgeBlock1hSamples, *edgeBlock1hMaxAvg)
	edgeHot := positiveEdgeHot(edgeMetrics, *edgeHotMinSamples, *edgeHotMinAvg, *edgeHotMinWinRate, *edgeHotMin5mAvg, *edgeHotMin15mAvg, *edgeHotMax1hNeg)
	now := time.Now()
	modeActions, err := loadModePolicyActions(*modePolicyPath, *modePolicyMaxAge, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: load mode policy: %v\n", err)
		os.Exit(1)
	}
	modeBlocks := modePolicyBlocksFromActions(modeActions, *modePolicyMinAction)
	allowedModes := parseModeSet(*modesRaw)
	ensureRequiredAlertModes(allowedModes, "EDGE-HOT", "FLOW-SCOUT")
	policy := alertPolicy{
		AllowedModes:            allowedModes,
		MinNotional:             *minNotional,
		ObserveMinNotional:      *observeMinNotional,
		ObserveBurstMinNotional: *observeBurstMinNotional,
		ObserveMaxBot:           *observeMaxBot,
		ObserveRequireKnown:     *observeRequireKnown,
		ObserveMinTier:          *observeMinTier,
		InsiderMinNotional:      *insiderMinNotional,
		InsiderMaxBot:           *insiderMaxBot,
		PositionCooldown:        *positionCooldown,
		RepeatMinNotional:       *repeatMinNotional,
		EdgeBlocks:              edgeBlocks,
		EdgeHot:                 edgeHot,
		RequirePositiveEdge:     parseModeSet(*requirePositiveEdgeRaw),
		EdgeMetrics:             edgeMetrics,
		EdgeHotMinNotional:      *edgeHotMinNotional,
		EdgeHotMaxBot:           *edgeHotMaxBot,
		EdgeHotMinSamples:       *edgeHotMinSamples,
		EdgeHotMinAvgPP:         *edgeHotMinAvg,
		EdgeHotMinWinRate:       *edgeHotMinWinRate,
		EdgeHotMin5mAvgPP:       *edgeHotMin5mAvg,
		EdgeHotMin15mAvgPP:      *edgeHotMin15mAvg,
		EdgeHotMax1hNegPP:       *edgeHotMax1hNeg,
		ModeBlocks:              modeBlocks,
		ModeActions:             modeActions,
		ModeMinAction:           *modePolicyMinAction,
		ConsensusAlerts:         *consensusAlerts,
		ConsensusMinNotional:    *consensusMinNotional,
		ConsensusMinWallets:     *consensusMinWallets,
		ConsensusMaxBot:         *consensusMaxBot,
		UnknownFlowMinNotional:  *unknownFlowMinNotional,
		UnknownFlowMinMarkets:   *unknownFlowMinMarkets,
		UnknownFlowMaxBot:       *unknownFlowMaxBot,
		SeedFlowMinNotional:     *seedFlowMinNotional,
		SeedFlowMinMarkets:      *seedFlowMinMarkets,
		ScoredFlowMinNotional:   *scoredFlowMinNotional,
		ScoredFlowMinMarkets:    *scoredFlowMinMarkets,
		ScoredFlowMaxBot:        *scoredFlowMaxBot,
		ScoredFlowMinTier:       *scoredFlowMinTier,
	}
	consensusAlertAge := *consensusMaxAge
	if consensusAlertAge <= 0 {
		consensusAlertAge = *maxAge
	}
	if err := writeCandidateDiagnostics(*diagnosticPath, trades, state, policy, now, *diagnosticAge, *maxAge, consensusAlertAge, *burstWindow, *burstMinNotional, *burstMinTrades, *burstMinLegNotional, 50); err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: write diagnostic report: %v\n", err)
		os.Exit(1)
	}
	shadowLogged := 0
	shadowLoggedKeys := map[string]struct{}{}
	if *shadowLogPath != "" {
		logged, err := loadLoggedAlertKeys(*shadowLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sports-tape-alert: load shadow log: %v\n", err)
			os.Exit(1)
		}
		shadowLoggedKeys = logged
		shadowTrades := staleObserveBurstTrades(trades, state, policy, logged, now, *diagnosticAge, *maxAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)
		if policy.ConsensusAlerts {
			shadowTrades = append(shadowTrades, staleConsensusTrades(trades, state, policy, logged, now, *diagnosticAge, consensusAlertAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
		}
		shadowTrades = append(shadowTrades, shadowUnknownFlowTrades(trades, policy, logged, now, *diagnosticAge, *maxAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
		shadowTrades = append(shadowTrades, shadowSeedFlowTrades(trades, policy, logged, now, *diagnosticAge, *maxAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
		shadowTrades = append(shadowTrades, shadowScoredFlowTrades(trades, policy, logged, now, *diagnosticAge, *maxAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
		if err := appendAlertLog(*shadowLogPath, shadowTrades, now, true, false); err != nil {
			fmt.Fprintf(os.Stderr, "sports-tape-alert: append shadow log: %v\n", err)
			os.Exit(1)
		}
		shadowLogged = len(shadowTrades)
	}

	candidates := qualifyingTrades(trades, state, now, policy, *maxAge)
	pendingState := cloneSentState(state)
	markTrades(pendingState, candidates, now)
	candidates = append(candidates, observeBurstTrades(trades, pendingState, policy, now, *maxAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
	if policy.ConsensusAlerts {
		candidates = append(candidates, consensusTrades(trades, state, policy, now, consensusAlertAge, *burstWindow, *burstMinTrades, *burstMinLegNotional)...)
		sort.Slice(candidates, func(i, j int) bool {
			if candidates[i].Time.Equal(candidates[j].Time) {
				return candidates[i].Notional > candidates[j].Notional
			}
			return candidates[i].Time.Before(candidates[j].Time)
		})
	}
	if *shadowLogPath != "" {
		blockedTrades := policyBlockedModeTrades(candidates, policy, shadowLoggedKeys)
		if err := appendAlertLog(*shadowLogPath, blockedTrades, now, true, false); err != nil {
			fmt.Fprintf(os.Stderr, "sports-tape-alert: append policy-blocked shadow log: %v\n", err)
			os.Exit(1)
		}
		shadowLogged += len(blockedTrades)
	}
	candidates = suppressPolicyBlockedModes(candidates, policy)
	if !existed && !*backfill {
		if *dryRun {
			fmt.Printf("sports-tape-alert initialized dry-run: marked=%d sent=0 state=%s\n", len(candidates), *statePath)
			return
		}
		markTrades(state, candidates, now)
		if err := saveState(*statePath, state); err != nil {
			fmt.Fprintf(os.Stderr, "sports-tape-alert: init state: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("sports-tape-alert initialized: marked=%d sent=0 state=%s\n", len(candidates), *statePath)
		return
	}
	if len(candidates) > *maxAlerts {
		candidates = candidates[:*maxAlerts]
	}

	sent := 0
	if len(candidates) > 0 {
		if *dryRun {
			for _, tr := range candidates {
				fmt.Println(formatAlert(tr))
				fmt.Println("---")
				sent++
			}
		} else {
			notifier := buildNotifier()
			for _, tr := range candidates {
				notifier.SidecarAlert(formatAlert(tr))
				sent++
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = notifier.Close(ctx)
			if err := appendAlertLog(*sentLogPath, candidates, now, false, false); err != nil {
				fmt.Fprintf(os.Stderr, "sports-tape-alert: append alert log: %v\n", err)
				os.Exit(1)
			}
			markTrades(state, candidates, now)
		}
	}
	if err := saveState(*statePath, state); err != nil {
		fmt.Fprintf(os.Stderr, "sports-tape-alert: save state: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sports-tape-alert done: candidates=%d sent=%d shadow=%d state=%s\n", len(candidates), sent, shadowLogged, *statePath)
}

func buildNotifier() notify.Notifier {
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		return notify.Nop{}
	}
	return notify.NewTelegram(notify.TelegramConfig{
		BotToken:       tok,
		ChatID:         chat,
		PromptBotToken: os.Getenv("SIDECAR_BOT_TOKEN"),
		PushBotToken:   os.Getenv("PUSH_BOT_TOKEN"),
		QueueSize:      16,
	})
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
			statuses, err := loadWalletStatuses(part)
			if err != nil {
				return nil, err
			}
			for addr, status := range statuses {
				if existing, ok := out[addr]; !ok || shouldOverrideWalletStatus(existing, status) {
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
		parts := strings.SplitN(line, "#", 2)
		fields := strings.Fields(parts[0])
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		status := walletStatus{}
		if len(parts) > 1 {
			comment := strings.TrimSpace(parts[1])
			if idx := strings.Index(comment, "reason="); idx >= 0 {
				status.Reason = strings.TrimSpace(comment[idx+len("reason="):])
				comment = strings.TrimSpace(comment[:idx])
			}
			gateStatus := ""
			for _, field := range strings.Fields(comment) {
				k, v, ok := strings.Cut(field, "=")
				if !ok {
					continue
				}
				switch k {
				case "list":
					status.List = v
				case "status":
					gateStatus = v
				case "tier":
					status.Tier = v
				case "source":
					status.Source = v
				case "smart":
					fmt.Sscanf(v, "%f", &status.Smart)
				case "bot":
					fmt.Sscanf(v, "%f", &status.Bot)
				case "edgeN":
					fmt.Sscanf(v, "%d", &status.EdgeN)
				case "edgeAvgPP":
					fmt.Sscanf(v, "%f", &status.EdgeAvg)
				case "edge5mPP":
					fmt.Sscanf(v, "%f", &status.Edge5m)
				case "edge15mPP":
					fmt.Sscanf(v, "%f", &status.Edge15m)
				}
			}
			if isBlockedTapeGateStatus(gateStatus) {
				status.List = "review_noise"
				if status.Reason == "" {
					status.Reason = gateStatus
				}
			}
		}
		out[addr] = status
	}
	return out, sc.Err()
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
	case "leaderboard_push":
		return 65
	case "flow", "watch", "target", "sports", "scout":
		return 60
	case "leaderboard_watch":
		return 15
	case "consensus_research":
		return 20
	case "tape_observe":
		return 10
	default:
		return 0
	}
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
		if status.Smart > 0 {
			trades[i].Smart = status.Smart
		}
		if status.Bot > 0 {
			trades[i].Bot = status.Bot
		}
	}
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
		if tr.Wallet == "" || tr.Time.IsZero() {
			continue
		}
		tr.Wallet = strings.ToLower(strings.TrimSpace(tr.Wallet))
		tr.Side = strings.ToUpper(strings.TrimSpace(tr.Side))
		out = append(out, tr)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.After(out[j].Time) })
	return out, nil
}

func loadEdgeMetrics(path string) (map[string]map[int64]*edgeStats, error) {
	out := map[string]map[int64]*edgeStats{}
	if path == "" {
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
		if wallet == "" {
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
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func loadModePolicyBlocks(path string, maxAge time.Duration, now time.Time) (map[string]string, error) {
	actions, err := loadModePolicyActions(path, maxAge, now)
	if err != nil {
		return nil, err
	}
	return modePolicyBlocksFromActions(actions, "COLLECT_POSITIVE"), nil
}

func loadModePolicyActions(path string, maxAge time.Duration, now time.Time) (map[string]string, error) {
	out := map[string]string{}
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

	var decision modePolicyDecision
	if err := json.NewDecoder(f).Decode(&decision); err != nil {
		return nil, err
	}
	if maxAge > 0 && !decision.GeneratedAt.IsZero() {
		age := now.Sub(decision.GeneratedAt)
		if age < 0 || age > maxAge {
			return out, nil
		}
	}
	for _, mode := range decision.Modes {
		key := strings.ToUpper(strings.TrimSpace(mode.Key))
		action := strings.ToUpper(strings.TrimSpace(mode.Action))
		if key != "" && action != "" {
			out[key] = action
		}
	}
	return out, nil
}

func modePolicyBlocksFromActions(actions map[string]string, minAction string) map[string]string {
	out := map[string]string{}
	for mode, action := range actions {
		if !modeActionMeetsMin(action, minAction) {
			out[mode] = fmt.Sprintf("policy action %s below %s", strings.ToUpper(action), strings.ToUpper(strings.TrimSpace(minAction)))
		}
	}
	return out
}

func modePolicyBlockReason(mode string, policy alertPolicy) (string, bool) {
	mode = strings.ToUpper(strings.TrimSpace(mode))
	if mode == "" {
		return "", false
	}
	if reason, blocked := policy.ModeBlocks[mode]; blocked {
		return reason, true
	}
	if len(policy.ModeActions) == 0 || strings.TrimSpace(policy.ModeMinAction) == "" {
		return "", false
	}
	action, ok := policy.ModeActions[mode]
	if !ok {
		return "missing mode policy action", true
	}
	if !modeActionMeetsMin(action, policy.ModeMinAction) {
		return fmt.Sprintf("policy action %s below %s", strings.ToUpper(action), strings.ToUpper(strings.TrimSpace(policy.ModeMinAction))), true
	}
	return "", false
}

func modeActionMeetsMin(action, minAction string) bool {
	minRank := modeActionRank(minAction)
	if minRank <= 0 {
		return true
	}
	return modeActionRank(action) >= minRank
}

func modeActionRank(action string) int {
	switch strings.ToUpper(strings.TrimSpace(action)) {
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

func negativeEdgeBlocks(metrics map[string]map[int64]*edgeStats, min15mSamples int, max15mAvgPP float64, min1hSamples int, max1hAvgPP float64) map[string]string {
	out := map[string]string{}
	for wallet, byHorizon := range metrics {
		if st := byHorizon[int64((15 * time.Minute).Seconds())]; st != nil && min15mSamples > 0 && st.Samples >= min15mSamples {
			avg := st.SumPP / float64(st.Samples)
			if avg <= max15mAvgPP {
				out[wallet] = fmt.Sprintf("15m edge %.2fpp over %d samples", avg, st.Samples)
				continue
			}
		}
		if st := byHorizon[int64((time.Hour).Seconds())]; st != nil && min1hSamples > 0 && st.Samples >= min1hSamples {
			avg := st.SumPP / float64(st.Samples)
			if avg <= max1hAvgPP {
				out[wallet] = fmt.Sprintf("1h edge %.2fpp over %d samples", avg, st.Samples)
			}
		}
	}
	return out
}

func positiveEdgeHot(metrics map[string]map[int64]*edgeStats, minSamples int, minAvgPP, minWinRate, min5mAvgPP, min15mAvgPP, max1hNegPP float64) map[string]string {
	out := map[string]string{}
	for wallet, byHorizon := range metrics {
		total := edgeStats{}
		for horizon, st := range byHorizon {
			if horizon <= 0 {
				continue
			}
			total.Samples += st.Samples
			total.Wins += st.Wins
			total.SumPP += st.SumPP
		}
		if minSamples > 0 && total.Samples < minSamples {
			continue
		}
		avg := avgEdge(&total)
		winRate := 0.0
		if total.Samples > 0 {
			winRate = float64(total.Wins) / float64(total.Samples) * 100
		}
		if avg < minAvgPP || winRate < minWinRate {
			continue
		}
		if st := byHorizon[int64((5 * time.Minute).Seconds())]; min5mAvgPP > 0 && (st == nil || avgEdge(st) < min5mAvgPP) {
			continue
		}
		if st := byHorizon[int64((15 * time.Minute).Seconds())]; st == nil || avgEdge(st) < min15mAvgPP {
			continue
		}
		if st := byHorizon[int64((time.Hour).Seconds())]; st != nil && avgEdge(st) <= max1hNegPP {
			continue
		}
		out[wallet] = fmt.Sprintf("edge-hot %.1f%% win avg %+0.2fpp over %d samples", winRate, avg, total.Samples)
	}
	return out
}

func avgEdge(st *edgeStats) float64 {
	if st == nil || st.Samples == 0 {
		return 0
	}
	return st.SumPP / float64(st.Samples)
}

func loadState(path string) (sentState, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return sentState{Sent: map[string]string{}}, false, nil
		}
		return sentState{}, false, err
	}
	var st sentState
	if err := json.Unmarshal(b, &st); err != nil {
		return sentState{}, true, err
	}
	if st.Sent == nil {
		st.Sent = map[string]string{}
	}
	if st.Details == nil {
		st.Details = map[string]sentDetail{}
	}
	return st, true, nil
}

func saveState(path string, st sentState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func reconcileAlertLog(path string, trades []tapeTrade, st sentState) error {
	if path == "" || len(st.Sent) == 0 {
		return nil
	}
	logged, err := loadLoggedAlerts(path)
	if err != nil {
		return err
	}
	if st.Details == nil {
		st.Details = map[string]sentDetail{}
	}
	for key, ev := range logged {
		if _, ok := st.Sent[key]; !ok {
			continue
		}
		st.Details[key] = sentDetail{Mode: strings.ToUpper(ev.Mode), KnownList: ev.KnownList}
	}
	var missing []tapeTrade
	sentAt := map[string]time.Time{}
	for _, tr := range trades {
		key := tradeKey(tr)
		ts, ok := st.Sent[key]
		if !ok {
			continue
		}
		if _, ok := logged[key]; ok {
			continue
		}
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t = time.Now()
		}
		missing = append(missing, tr)
		sentAt[key] = t
	}
	if len(missing) == 0 {
		return nil
	}
	return appendAlertLogWithTimes(path, missing, sentAt, false, true)
}

func loadLoggedAlertKeys(path string) (map[string]struct{}, error) {
	logged, err := loadLoggedAlerts(path)
	if err != nil {
		return nil, err
	}
	out := map[string]struct{}{}
	for key := range logged {
		out[key] = struct{}{}
	}
	return out, nil
}

func loadLoggedAlerts(path string) (map[string]alertEvent, error) {
	out := map[string]alertEvent{}
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
		var ev alertEvent
		if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
			continue
		}
		if ev.Key != "" {
			out[strings.ToLower(ev.Key)] = ev
		}
	}
	return out, sc.Err()
}

func appendAlertLog(path string, trades []tapeTrade, sentAt time.Time, dryRun, reconciled bool) error {
	times := map[string]time.Time{}
	for _, tr := range trades {
		times[tradeKey(tr)] = sentAt
	}
	return appendAlertLogWithTimes(path, trades, times, dryRun, reconciled)
}

func appendAlertLogWithTimes(path string, trades []tapeTrade, sentAt map[string]time.Time, dryRun, reconciled bool) error {
	if path == "" || len(trades) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, tr := range trades {
		key := tradeKey(tr)
		t := sentAt[key]
		if t.IsZero() {
			t = time.Now()
		}
		if err := enc.Encode(alertEventFromTrade(tr, key, t, dryRun, reconciled)); err != nil {
			return err
		}
	}
	return nil
}

func alertEventFromTrade(tr tapeTrade, key string, sentAt time.Time, dryRun, reconciled bool) alertEvent {
	return alertEvent{
		Key:           strings.ToLower(key),
		SentAt:        sentAt,
		TradeTime:     tr.Time,
		DryRun:        dryRun,
		Reconciled:    reconciled,
		Mode:          alertMode(tr),
		Wallet:        strings.ToLower(tr.Wallet),
		KnownList:     tr.KnownList,
		Tier:          tr.Tier,
		Bot:           tr.Bot,
		Category:      tr.Category,
		Notional:      tr.Notional,
		Price:         tr.Price,
		Outcome:       tr.Outcome,
		Market:        tr.Market,
		Slug:          tr.Slug,
		Asset:         tr.Asset,
		Transaction:   tr.Transaction,
		TargetCopyROI: tr.TargetCopyROI,
		TargetCopyT:   tr.TargetCopyT,
		Participants:  append([]string(nil), tr.Participants...),
	}
}

func qualifyingTrades(trades []tapeTrade, st sentState, now time.Time, policy alertPolicy, maxAge time.Duration) []tapeTrade {
	var raw []tapeTrade
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 {
			continue
		}
		ok, reason := alertDecision(tr, policy)
		if !ok {
			continue
		}
		tr = annotateAlertTrade(tr, reason)
		if maxAge > 0 {
			age := now.Sub(tr.Time)
			if age < 0 || age > maxAge {
				continue
			}
		}
		if _, ok := st.Sent[tradeKey(tr)]; ok {
			continue
		}
		raw = append(raw, tr)
	}
	sort.Slice(raw, func(i, j int) bool {
		if raw[i].Time.Equal(raw[j].Time) {
			return raw[i].Notional > raw[j].Notional
		}
		return raw[i].Time.Before(raw[j].Time)
	})
	lastByPosition := lastSentByAssetWallet(st)
	out := make([]tapeTrade, 0, len(raw))
	for _, tr := range raw {
		if active, _ := positionCooldownActive(lastByPosition, tr, now, policy); active {
			continue
		}
		out = append(out, tr)
		if policy.PositionCooldown > 0 {
			lastByPosition[assetWalletKey(tr.Asset, tr.Wallet)] = now
		}
	}
	return out
}

func suppressPolicyBlockedModes(trades []tapeTrade, policy alertPolicy) []tapeTrade {
	if (len(policy.ModeBlocks) == 0 && (len(policy.ModeActions) == 0 || strings.TrimSpace(policy.ModeMinAction) == "")) || len(trades) == 0 {
		return trades
	}
	out := make([]tapeTrade, 0, len(trades))
	for _, tr := range trades {
		mode := strings.ToUpper(strings.TrimSpace(alertMode(tr)))
		if _, blocked := modePolicyBlockReason(mode, policy); blocked {
			continue
		}
		out = append(out, tr)
	}
	return out
}

func policyBlockedModeTrades(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}) []tapeTrade {
	if (len(policy.ModeBlocks) == 0 && (len(policy.ModeActions) == 0 || strings.TrimSpace(policy.ModeMinAction) == "")) || len(trades) == 0 {
		return nil
	}
	out := make([]tapeTrade, 0, len(trades))
	for _, tr := range trades {
		mode := strings.ToUpper(strings.TrimSpace(alertMode(tr)))
		if _, blocked := modePolicyBlockReason(mode, policy); !blocked {
			continue
		}
		if _, ok := logged[tradeKey(tr)]; ok {
			continue
		}
		out = append(out, tr)
		logged[tradeKey(tr)] = struct{}{}
	}
	return out
}

func annotateAlertTrade(tr tapeTrade, reason string) tapeTrade {
	if reason == "insider-scout huge whale" && strings.TrimSpace(tr.KnownList) == "" {
		tr.KnownList = "insider_scout"
	}
	return tr
}

func shouldAlertTrade(tr tapeTrade, policy alertPolicy) bool {
	ok, _ := alertDecision(tr, policy)
	return ok
}

func alertDecision(tr tapeTrade, policy alertPolicy) (bool, string) {
	mode := alertMode(tr)
	if mode == "REVIEW-NOISE" {
		return false, "review-noise excluded"
	}
	if mode == "REVERSAL-RISK" && !policy.AllowReversalAlerts {
		return false, "reversal risk disabled"
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		return false, reason
	}
	if mode == "EDGE-HOT" {
		if _, ok := policy.AllowedModes[mode]; !ok {
			return false, "EDGE-HOT mode disabled"
		}
		edgeReason, ok := policy.EdgeHot[strings.ToLower(tr.Wallet)]
		if !ok {
			return false, "edge-hot thresholds not met"
		}
		if strings.EqualFold(tr.Tier, "BOT") {
			return false, "BOT tier"
		}
		if policy.EdgeHotMaxBot > 0 && tr.Bot > policy.EdgeHotMaxBot {
			return false, fmt.Sprintf("bot %.1f > %.1f", tr.Bot, policy.EdgeHotMaxBot)
		}
		if tr.Notional < policy.EdgeHotMinNotional {
			return false, fmt.Sprintf("notional $%.0f < edge-hot min $%.0f", tr.Notional, policy.EdgeHotMinNotional)
		}
		return true, edgeReason
	}
	if _, ok := policy.AllowedModes[mode]; ok {
		if _, required := policy.RequirePositiveEdge[mode]; required {
			edgeReason, ok := policy.EdgeHot[strings.ToLower(tr.Wallet)]
			if !ok {
				return false, mode + " requires positive edge evidence"
			}
			if strings.EqualFold(tr.Tier, "BOT") {
				return false, "BOT tier"
			}
			if policy.EdgeHotMaxBot > 0 && tr.Bot > policy.EdgeHotMaxBot {
				return false, fmt.Sprintf("bot %.1f > %.1f", tr.Bot, policy.EdgeHotMaxBot)
			}
			if tr.Notional < policy.MinNotional {
				return false, fmt.Sprintf("notional $%.0f < min $%.0f", tr.Notional, policy.MinNotional)
			}
			return true, edgeReason
		}
		if tr.Notional < policy.MinNotional {
			return false, fmt.Sprintf("notional $%.0f < min $%.0f", tr.Notional, policy.MinNotional)
		}
		return true, "mode allowed"
	}
	if mode == "OBSERVE" && (policy.ObserveMinNotional > 0 || policy.InsiderMinNotional > 0) {
		if strings.EqualFold(tr.KnownList, "consensus_research") {
			return false, "consensus research wallet requires consensus burst"
		}
		if policy.InsiderMinNotional > 0 && tr.Notional >= policy.InsiderMinNotional {
			if strings.EqualFold(tr.Tier, "BOT") {
				return false, "BOT tier"
			}
			if policy.InsiderMaxBot > 0 && tr.Bot > policy.InsiderMaxBot {
				return false, fmt.Sprintf("bot %.1f > insider max %.1f", tr.Bot, policy.InsiderMaxBot)
			}
			return true, "insider-scout huge whale"
		}
		if policy.ObserveMinNotional <= 0 {
			return false, "observe mode disabled"
		}
		if policy.ObserveRequireKnown && !hasKnownWalletMeta(tr) {
			return false, "observe wallet unscored"
		}
		if strings.EqualFold(tr.Tier, "BOT") {
			return false, "BOT tier"
		}
		if policy.ObserveMaxBot > 0 && tr.Bot > policy.ObserveMaxBot {
			return false, fmt.Sprintf("bot %.1f > %.1f", tr.Bot, policy.ObserveMaxBot)
		}
		if !tierAtLeast(tr.Tier, policy.ObserveMinTier) {
			return false, fmt.Sprintf("tier %s below observe min %s", dash(tr.Tier), strings.ToUpper(strings.TrimSpace(policy.ObserveMinTier)))
		}
		if tr.Notional < policy.ObserveMinNotional {
			return false, fmt.Sprintf("notional $%.0f < observe min $%.0f", tr.Notional, policy.ObserveMinNotional)
		}
		return true, "huge observe non-bot"
	}
	return false, "mode disabled"
}

func hasKnownWalletMeta(tr tapeTrade) bool {
	return strings.TrimSpace(tr.KnownList) != "" || strings.TrimSpace(tr.Tier) != "" || tr.Bot > 0 || tr.TargetCopyT > 0 || tr.TargetCopyROI != 0
}

func hasScoredWalletMeta(tr tapeTrade) bool {
	return strings.TrimSpace(tr.Tier) != "" || tr.Bot > 0 || tr.TargetCopyT > 0 || tr.TargetCopyROI != 0
}

func writeCandidateDiagnostics(path string, trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, diagnosticAge, alertAge, consensusAlertAge, burstWindow time.Duration, burstMinNotional float64, burstMinTrades int, burstMinLegNotional float64, limit int) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var rows []diagnosticRow
	eligible := 0
	currentEligible := 0
	sent := 0
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 {
			continue
		}
		age := now.Sub(tr.Time)
		if diagnosticAge > 0 && (age < 0 || age > diagnosticAge) {
			continue
		}
		ok, reason := alertDecision(tr, policy)
		if ok && alertAge > 0 && age > alertAge {
			reason = fmt.Sprintf("stale: age %s > alert window %s", formatDuration(age), alertAge)
		}
		alreadySent := false
		mode := alertMode(tr)
		if _, exists := st.Sent[tradeKey(tr)]; exists {
			alreadySent = true
			sent++
			mode = sentMode(st, tradeKey(tr), mode)
			if ok {
				reason = "already sent as " + mode
				ok = false
			}
		}
		if ok {
			if active, cooldownReason := positionCooldownActive(lastSentByAssetWallet(st), tr, now, policy); active {
				reason = cooldownReason
				ok = false
			}
		}
		if ok {
			eligible++
			if alertAge <= 0 || age <= alertAge {
				currentEligible++
			}
		}
		alertable := ok && (alertAge <= 0 || age <= alertAge)
		rows = append(rows, diagnosticRow{Trade: tr, Mode: mode, Eligible: ok, Alertable: alertable, Reason: reason, Age: age, Sent: alreadySent})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Alertable != rows[j].Alertable {
			return rows[i].Alertable
		}
		if rows[i].Eligible != rows[j].Eligible {
			return rows[i].Eligible
		}
		if rows[i].Trade.Notional != rows[j].Trade.Notional {
			return rows[i].Trade.Notional > rows[j].Trade.Notional
		}
		return rows[i].Trade.Time.After(rows[j].Trade.Time)
	})
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Sports Alert Candidates\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", now.Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Recent BUY rows inspected: %d\n", len(rows))
	fmt.Fprintf(&b, "- Diagnostic window: %s\n", diagnosticAge)
	fmt.Fprintf(&b, "- Alert window: %s\n", alertAge)
	fmt.Fprintf(&b, "- Consensus alert window: %s\n", consensusAlertAge)
	fmt.Fprintf(&b, "- Currently alertable unsent rows: %d\n", currentEligible)
	fmt.Fprintf(&b, "- Eligible unsent rows in diagnostic window: %d\n", eligible)
	fmt.Fprintf(&b, "- Already sent rows in window: %d\n", sent)
	fmt.Fprintf(&b, "- Positive-edge required modes: %s\n", formatModeSet(policy.RequirePositiveEdge))
	fmt.Fprintf(&b, "- Observe min notional: $%.0f\n", policy.ObserveMinNotional)
	fmt.Fprintf(&b, "- Observe-burst min notional: $%.0f\n", policy.ObserveBurstMinNotional)
	fmt.Fprintf(&b, "- Unknown-flow min notional: $%.0f\n", policy.UnknownFlowMinNotional)
	fmt.Fprintf(&b, "- Unknown-flow min markets: %d\n", policy.UnknownFlowMinMarkets)
	fmt.Fprintf(&b, "- Seed-flow min notional: $%.0f\n", policy.SeedFlowMinNotional)
	fmt.Fprintf(&b, "- Seed-flow min markets: %d\n", policy.SeedFlowMinMarkets)
	fmt.Fprintf(&b, "- Scored-flow min notional: $%.0f\n", policy.ScoredFlowMinNotional)
	fmt.Fprintf(&b, "- Scored-flow min markets: %d\n", policy.ScoredFlowMinMarkets)
	fmt.Fprintf(&b, "- Scored-flow min tier: %s\n", dash(policy.ScoredFlowMinTier))
	fmt.Fprintf(&b, "- Scored-flow max bot: %.1f\n", policy.ScoredFlowMaxBot)
	fmt.Fprintf(&b, "- Observe min tier: %s\n", dash(policy.ObserveMinTier))
	fmt.Fprintf(&b, "- Insider-scout min notional: $%.0f\n", policy.InsiderMinNotional)
	fmt.Fprintf(&b, "- Insider-scout max bot: %.1f\n", policy.InsiderMaxBot)
	fmt.Fprintf(&b, "- Edge-hot wallets: %d\n", len(policy.EdgeHot))
	fmt.Fprintf(&b, "- Negative-edge blocked wallets: %d\n\n", len(policy.EdgeBlocks))
	writeWalletReasonSection(&b, "Edge-Hot Wallets", policy.EdgeHot)
	writeNearEdgeHotSection(&b, rows, policy, 15)
	writeWalletReasonSection(&b, "Negative-Edge Blocked Wallets", policy.EdgeBlocks)
	writeBurstSection(&b, buildBurstCandidates(trades, st, policy, now, diagnosticAge, alertAge, burstWindow, burstMinNotional, burstMinTrades, burstMinLegNotional), now, 20)
	writeObserveBurstSection(&b, buildObserveBurstCandidates(trades, st, policy, now, diagnosticAge, alertAge, burstWindow, burstMinTrades, burstMinLegNotional), now, 20)
	writeConsensusSection(&b, buildConsensusCandidates(trades, st, policy, now, diagnosticAge, consensusAlertAge, burstWindow, burstMinTrades, burstMinLegNotional), now, 20)
	writeUnknownFlowSection(&b, buildUnknownFlowCandidates(trades, policy, nil, now, diagnosticAge, alertAge, burstWindow, burstMinTrades, burstMinLegNotional), now, 20)
	writeSeedFlowSection(&b, buildSeedFlowCandidates(trades, policy, now, diagnosticAge, alertAge, burstWindow, burstMinTrades, burstMinLegNotional), now, 20)
	writeScoredFlowSection(&b, buildScoredFlowCandidates(trades, policy, nil, now, diagnosticAge, alertAge, burstWindow, burstMinTrades, burstMinLegNotional), now, 20)
	fmt.Fprintf(&b, "## Trade Rows\n\n")
	fmt.Fprintf(&b, "| Status | Mode | Wallet | List | Tier | Bot | Notional | Age | Reason | Market |\n")
	fmt.Fprintf(&b, "|---|---|---|---|---|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(&b, "| none |  |  |  |  | 0.0 | $0 |  |  |  |\n")
	}
	for _, r := range rows {
		status := "blocked"
		if r.Alertable {
			status = "eligible"
		} else if r.Eligible {
			status = "stale"
		} else if r.Sent {
			status = "sent"
		}
		fmt.Fprintf(&b, "| %s | %s | `%s` | %s | %s | %.1f | $%.0f | %s | %s | %s |\n",
			status, r.Mode, shortAddr(r.Trade.Wallet), dash(r.Trade.KnownList), dash(r.Trade.Tier), r.Trade.Bot,
			r.Trade.Notional, formatDuration(r.Age), oneLine(r.Reason, 60), oneLine(r.Trade.Market, 80))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeNearEdgeHotSection(b *strings.Builder, rows []diagnosticRow, policy alertPolicy, limit int) {
	type near struct {
		Trade   tapeTrade
		Samples int
		AvgPP   float64
		WinRate float64
		Avg5m   float64
		Avg15m  float64
		Avg1h   float64
		Reason  string
	}
	byWallet := map[string]near{}
	for _, r := range rows {
		if r.Mode != "EDGE-HOT" {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(r.Trade.Wallet))
		if wallet == "" {
			continue
		}
		if _, hot := policy.EdgeHot[wallet]; hot {
			continue
		}
		reason, samples, avg, winRate, avg5m, avg15m, avg1h := edgeHotGap(wallet, r.Trade, policy)
		n := near{
			Trade:   r.Trade,
			Samples: samples,
			AvgPP:   avg,
			WinRate: winRate,
			Avg5m:   avg5m,
			Avg15m:  avg15m,
			Avg1h:   avg1h,
			Reason:  reason,
		}
		prev, ok := byWallet[wallet]
		if !ok || n.Trade.Notional > prev.Trade.Notional {
			byWallet[wallet] = n
		}
	}
	out := make([]near, 0, len(byWallet))
	for _, n := range byWallet {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Samples != out[j].Samples {
			return out[i].Samples > out[j].Samples
		}
		if out[i].AvgPP != out[j].AvgPP {
			return out[i].AvgPP > out[j].AvgPP
		}
		return out[i].Trade.Notional > out[j].Trade.Notional
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}

	fmt.Fprintf(b, "## Near Edge-Hot Wallets\n\n")
	fmt.Fprintf(b, "- Rule: recent large BUYs from edge-hot candidate lists that are not yet eligible; this is diagnostics only and does not loosen Telegram pushes.\n\n")
	fmt.Fprintf(b, "| Wallet | List | Tier | Bot | Notional | Samples | Win | AvgPP | 5m | 15m | 1h | Gap | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	if len(out) == 0 {
		fmt.Fprintf(b, "| none |  |  | 0.0 | $0 | 0 | 0.0%% | +0.00 | +0.00 | +0.00 | +0.00 |  |  |\n\n")
		return
	}
	for _, n := range out {
		fmt.Fprintf(b, "| `%s` | %s | %s | %.1f | $%.0f | %d | %.1f%% | %+0.2f | %+0.2f | %+0.2f | %+0.2f | %s | %s |\n",
			shortAddr(n.Trade.Wallet), dash(n.Trade.KnownList), dash(n.Trade.Tier), n.Trade.Bot, n.Trade.Notional,
			n.Samples, n.WinRate, n.AvgPP, n.Avg5m, n.Avg15m, n.Avg1h, oneLine(n.Reason, 70), oneLine(n.Trade.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func edgeHotGap(wallet string, tr tapeTrade, policy alertPolicy) (string, int, float64, float64, float64, float64, float64) {
	if reason, blocked := policy.EdgeBlocks[wallet]; blocked {
		return reason, 0, 0, 0, 0, 0, 0
	}
	byHorizon := policy.EdgeMetrics[wallet]
	if byHorizon == nil {
		return "no edge snapshots", 0, 0, 0, 0, 0, 0
	}
	total := edgeStats{}
	for horizon, st := range byHorizon {
		if horizon <= 0 || st == nil {
			continue
		}
		total.Samples += st.Samples
		total.Wins += st.Wins
		total.SumPP += st.SumPP
	}
	avg := avgEdge(&total)
	winRate := 0.0
	if total.Samples > 0 {
		winRate = float64(total.Wins) / float64(total.Samples) * 100
	}
	avg5m := avgEdge(byHorizon[int64((5 * time.Minute).Seconds())])
	avg15m := avgEdge(byHorizon[int64((15 * time.Minute).Seconds())])
	avg1h := avgEdge(byHorizon[int64((time.Hour).Seconds())])

	switch {
	case policy.EdgeHotMinSamples > 0 && total.Samples < policy.EdgeHotMinSamples:
		return fmt.Sprintf("samples %d < %d", total.Samples, policy.EdgeHotMinSamples), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case avg < policy.EdgeHotMinAvgPP:
		return fmt.Sprintf("avg %+0.2fpp < %+0.2fpp", avg, policy.EdgeHotMinAvgPP), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case winRate < policy.EdgeHotMinWinRate:
		return fmt.Sprintf("win %.1f%% < %.1f%%", winRate, policy.EdgeHotMinWinRate), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case policy.EdgeHotMin5mAvgPP > 0 && (byHorizon[int64((5*time.Minute).Seconds())] == nil || avg5m < policy.EdgeHotMin5mAvgPP):
		return fmt.Sprintf("5m %+0.2fpp < %+0.2fpp", avg5m, policy.EdgeHotMin5mAvgPP), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case byHorizon[int64((15*time.Minute).Seconds())] == nil || avg15m < policy.EdgeHotMin15mAvgPP:
		return fmt.Sprintf("15m %+0.2fpp < %+0.2fpp", avg15m, policy.EdgeHotMin15mAvgPP), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case byHorizon[int64((time.Hour).Seconds())] != nil && avg1h <= policy.EdgeHotMax1hNegPP:
		return fmt.Sprintf("1h %+0.2fpp <= %+0.2fpp", avg1h, policy.EdgeHotMax1hNegPP), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case strings.EqualFold(tr.Tier, "BOT"):
		return "BOT tier", total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case policy.EdgeHotMaxBot > 0 && tr.Bot > policy.EdgeHotMaxBot:
		return fmt.Sprintf("bot %.1f > %.1f", tr.Bot, policy.EdgeHotMaxBot), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	case tr.Notional < policy.EdgeHotMinNotional:
		return fmt.Sprintf("notional $%.0f < $%.0f", tr.Notional, policy.EdgeHotMinNotional), total.Samples, avg, winRate, avg5m, avg15m, avg1h
	default:
		return "meets edge metrics but not selected", total.Samples, avg, winRate, avg5m, avg15m, avg1h
	}
}

func writeWalletReasonSection(b *strings.Builder, title string, rows map[string]string) {
	fmt.Fprintf(b, "## %s\n\n", title)
	if len(rows) == 0 {
		fmt.Fprintf(b, "No wallets in this section.\n\n")
		return
	}
	keys := make([]string, 0, len(rows))
	for wallet := range rows {
		keys = append(keys, wallet)
	}
	sort.Strings(keys)
	fmt.Fprintf(b, "| Wallet | Reason |\n")
	fmt.Fprintf(b, "|---|---|\n")
	for _, wallet := range keys {
		fmt.Fprintf(b, "| `%s` | %s |\n", shortAddr(wallet), oneLine(rows[wallet], 100))
	}
	fmt.Fprintf(b, "\n")
}

func buildBurstCandidates(trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minNotional float64, minTrades int, minLegNotional float64) []burstCandidate {
	type key struct {
		wallet string
		asset  string
	}
	groups := map[key][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 || tr.Notional < minLegNotional {
			continue
		}
		age := now.Sub(tr.Time)
		if diagnosticAge > 0 && (age < 0 || age > diagnosticAge) {
			continue
		}
		mode := alertMode(tr)
		if !isBurstMode(mode) {
			continue
		}
		if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
			_ = reason
			continue
		}
		if mode == "EDGE-HOT" {
			if _, ok := policy.EdgeHot[strings.ToLower(tr.Wallet)]; !ok {
				continue
			}
			if strings.EqualFold(tr.Tier, "BOT") || (policy.EdgeHotMaxBot > 0 && tr.Bot > policy.EdgeHotMaxBot) {
				continue
			}
		}
		if _, required := policy.RequirePositiveEdge[mode]; required {
			if _, ok := policy.EdgeHot[strings.ToLower(tr.Wallet)]; !ok {
				continue
			}
			if strings.EqualFold(tr.Tier, "BOT") || (policy.EdgeHotMaxBot > 0 && tr.Bot > policy.EdgeHotMaxBot) {
				continue
			}
		}
		if _, ok := policy.AllowedModes[mode]; !ok {
			continue
		}
		groups[key{wallet: strings.ToLower(tr.Wallet), asset: tr.Asset}] = append(groups[key{wallet: strings.ToLower(tr.Wallet), asset: tr.Asset}], tr)
	}

	var out []burstCandidate
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= burstWindow {
				start--
			}
			window := rows[start : end+1]
			if len(window) < minTrades {
				continue
			}
			var total float64
			for _, tr := range window {
				total += tr.Notional
			}
			if total < minNotional {
				continue
			}
			last := window[len(window)-1]
			first := window[0]
			age := now.Sub(last.Time)
			reason := "burst in diagnostic window"
			if alertAge > 0 && age > alertAge {
				reason = fmt.Sprintf("stale burst: last age %s > alert window %s", formatDuration(age), alertAge)
			}
			out = append(out, burstCandidate{
				Wallet:        last.Wallet,
				Asset:         last.Asset,
				Mode:          alertMode(last),
				KnownList:     last.KnownList,
				Tier:          last.Tier,
				Bot:           last.Bot,
				Category:      last.Category,
				Market:        last.Market,
				Outcome:       last.Outcome,
				Slug:          last.Slug,
				Trades:        len(window),
				TotalNotional: total,
				LastNotional:  last.Notional,
				LastPrice:     last.Price,
				FirstTime:     first.Time,
				LastTime:      last.Time,
				AlreadySent:   hasSentAssetWallet(st, last.Asset, last.Wallet),
				Reason:        reason,
			})
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AlreadySent != out[j].AlreadySent {
			return !out[i].AlreadySent
		}
		if out[i].TotalNotional != out[j].TotalNotional {
			return out[i].TotalNotional > out[j].TotalNotional
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
	return out
}

func observeBurstTrades(trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, maxAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildObserveBurstCandidates(trades, st, policy, now, burstScanAge(maxAge, burstWindow), maxAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadySent || strings.HasPrefix(c.Reason, "stale") || c.TotalNotional <= 0 {
			continue
		}
		out = append(out, observeBurstTrade(c))
	}
	return out
}

func staleObserveBurstTrades(trades []tapeTrade, st sentState, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildObserveBurstCandidates(trades, st, policy, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadySent || !strings.HasPrefix(c.Reason, "stale") || c.TotalNotional <= 0 {
			continue
		}
		tr := observeBurstTrade(c)
		if _, ok := logged[tradeKey(tr)]; ok {
			continue
		}
		out = append(out, tr)
	}
	return out
}

func buildObserveBurstCandidates(trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []burstCandidate {
	if policy.ObserveBurstMinNotional <= 0 {
		return nil
	}
	type key struct {
		wallet string
		asset  string
	}
	groups := map[key][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 || tr.Notional < minLegNotional {
			continue
		}
		age := now.Sub(tr.Time)
		if diagnosticAge > 0 && (age < 0 || age > diagnosticAge) {
			continue
		}
		if !observeBurstEligible(tr, policy) {
			continue
		}
		k := key{wallet: strings.ToLower(tr.Wallet), asset: tr.Asset}
		groups[k] = append(groups[k], tr)
	}

	var out []burstCandidate
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= burstWindow {
				start--
			}
			window := rows[start : end+1]
			if len(window) < minTrades {
				continue
			}
			var total, totalSize float64
			for _, tr := range window {
				total += tr.Notional
				totalSize += tr.Size
			}
			if total < policy.ObserveBurstMinNotional {
				continue
			}
			first := window[0]
			last := window[len(window)-1]
			age := now.Sub(last.Time)
			reason := "observe burst in diagnostic window"
			if alertAge > 0 && age > alertAge {
				reason = fmt.Sprintf("stale observe burst: last age %s > alert window %s", formatDuration(age), alertAge)
			}
			price := last.Price
			if totalSize > 0 {
				price = total / totalSize
			}
			out = append(out, burstCandidate{
				Wallet:        last.Wallet,
				Asset:         last.Asset,
				Mode:          "OBSERVE-BURST",
				KnownList:     "observe_burst",
				Tier:          last.Tier,
				Bot:           last.Bot,
				Category:      last.Category,
				Market:        last.Market,
				Outcome:       last.Outcome,
				Slug:          last.Slug,
				Trades:        len(window),
				TotalNotional: total,
				LastNotional:  last.Notional,
				LastPrice:     price,
				FirstTime:     first.Time,
				LastTime:      last.Time,
				AlreadySent:   hasSentAssetWallet(st, last.Asset, last.Wallet),
				Reason:        reason,
			})
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AlreadySent != out[j].AlreadySent {
			return !out[i].AlreadySent
		}
		if out[i].TotalNotional != out[j].TotalNotional {
			return out[i].TotalNotional > out[j].TotalNotional
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
	return out
}

func observeBurstEligible(tr tapeTrade, policy alertPolicy) bool {
	if alertMode(tr) != "OBSERVE" {
		return false
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		_ = reason
		return false
	}
	if policy.ObserveRequireKnown && !hasKnownWalletMeta(tr) {
		return false
	}
	if strings.EqualFold(tr.Tier, "BOT") {
		return false
	}
	if policy.ObserveMaxBot > 0 && tr.Bot > policy.ObserveMaxBot {
		return false
	}
	return tierAtLeast(tr.Tier, policy.ObserveMinTier)
}

func observeBurstTrade(c burstCandidate) tapeTrade {
	return tapeTrade{
		Time:        c.LastTime,
		Wallet:      c.Wallet,
		Side:        "BUY",
		Notional:    c.TotalNotional,
		Price:       c.LastPrice,
		Outcome:     c.Outcome,
		Market:      c.Market,
		Slug:        c.Slug,
		Category:    c.Category,
		Asset:       c.Asset,
		Transaction: observeBurstKey(c),
		KnownList:   "observe_burst",
		Tier:        c.Tier,
		Bot:         c.Bot,
	}
}

func observeBurstKey(c burstCandidate) string {
	return strings.ToLower(fmt.Sprintf("observe-burst|%s|%s|%d|%d|%.4f", c.Wallet, c.Asset, c.LastTime.Unix(), c.Trades, c.TotalNotional))
}

func shadowUnknownFlowTrades(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildUnknownFlowCandidates(trades, policy, logged, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadyLogged || c.TotalNotional <= 0 {
			continue
		}
		out = append(out, unknownFlowTrade(c))
	}
	return out
}

func buildUnknownFlowCandidates(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []unknownFlowCandidate {
	return buildMultiMarketFlowCandidates(trades, policy, logged, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional, policy.UnknownFlowMinNotional, policy.UnknownFlowMinMarkets, unknownFlowEligible, "unknown multi-market flow in diagnostic window", "stale unknown flow")
}

func shadowSeedFlowTrades(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildSeedFlowCandidates(trades, policy, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.TotalNotional <= 0 {
			continue
		}
		tr := seedFlowTrade(c)
		if logged != nil {
			if _, ok := logged[tradeKey(tr)]; ok {
				continue
			}
		}
		out = append(out, tr)
	}
	return out
}

func buildSeedFlowCandidates(trades []tapeTrade, policy alertPolicy, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []unknownFlowCandidate {
	rows := buildMultiMarketFlowCandidates(trades, policy, nil, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional, policy.SeedFlowMinNotional, policy.SeedFlowMinMarkets, seedFlowEligible, "seed multi-market flow in diagnostic window", "stale seed flow")
	if policy.UnknownFlowMinNotional <= 0 {
		return rows
	}
	out := rows[:0]
	for _, r := range rows {
		if r.TotalNotional < policy.UnknownFlowMinNotional {
			out = append(out, r)
		}
	}
	return out
}

func shadowScoredFlowTrades(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildScoredFlowCandidates(trades, policy, logged, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadyLogged || c.TotalNotional <= 0 {
			continue
		}
		out = append(out, scoredFlowTrade(c))
	}
	return out
}

func buildScoredFlowCandidates(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []unknownFlowCandidate {
	return buildMultiMarketFlowCandidates(trades, policy, logged, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional, policy.ScoredFlowMinNotional, policy.ScoredFlowMinMarkets, scoredFlowEligible, "scored multi-market flow in diagnostic window", "stale scored flow")
}

func buildMultiMarketFlowCandidates(trades []tapeTrade, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64, minNotional float64, minMarkets int, eligible func(tapeTrade, alertPolicy) bool, activeReason, stalePrefix string) []unknownFlowCandidate {
	if minNotional <= 0 {
		return nil
	}
	if minMarkets <= 0 {
		minMarkets = 2
	}
	groups := map[string][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 || tr.Notional < minLegNotional {
			continue
		}
		age := now.Sub(tr.Time)
		if diagnosticAge > 0 && (age < 0 || age > diagnosticAge) {
			continue
		}
		if !eligible(tr, policy) {
			continue
		}
		wallet := strings.ToLower(tr.Wallet)
		groups[wallet] = append(groups[wallet], tr)
	}

	var out []unknownFlowCandidate
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= burstWindow {
				start--
			}
			window := rows[start : end+1]
			if len(window) < minTrades {
				continue
			}
			var total float64
			markets := map[string]struct{}{}
			for _, tr := range window {
				total += tr.Notional
				markets[unknownFlowMarketKey(tr)] = struct{}{}
			}
			if len(markets) < minMarkets || total < minNotional {
				continue
			}
			first := window[0]
			last := window[len(window)-1]
			age := now.Sub(last.Time)
			reason := activeReason
			if alertAge > 0 && age > alertAge {
				reason = fmt.Sprintf("%s: last age %s > alert window %s", stalePrefix, formatDuration(age), alertAge)
			}
			c := unknownFlowCandidate{
				Wallet:        last.Wallet,
				Category:      last.Category,
				Market:        last.Market,
				Outcome:       last.Outcome,
				Slug:          last.Slug,
				Asset:         last.Asset,
				Trades:        len(window),
				Markets:       len(markets),
				TotalNotional: total,
				LastNotional:  last.Notional,
				LastPrice:     last.Price,
				FirstTime:     first.Time,
				LastTime:      last.Time,
				Tier:          last.Tier,
				Bot:           last.Bot,
				Reason:        reason,
			}
			tr := unknownFlowTrade(c)
			if strings.Contains(stalePrefix, "scored") {
				tr = scoredFlowTrade(c)
			}
			if logged != nil {
				if _, ok := logged[tradeKey(tr)]; ok {
					c.AlreadyLogged = true
				}
			}
			out = append(out, c)
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AlreadyLogged != out[j].AlreadyLogged {
			return !out[i].AlreadyLogged
		}
		if out[i].TotalNotional != out[j].TotalNotional {
			return out[i].TotalNotional > out[j].TotalNotional
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
	return out
}

func unknownFlowEligible(tr tapeTrade, policy alertPolicy) bool {
	if alertMode(tr) != "OBSERVE" {
		return false
	}
	if hasKnownWalletMeta(tr) {
		return false
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		_ = reason
		return false
	}
	if strings.EqualFold(tr.Tier, "BOT") {
		return false
	}
	if policy.UnknownFlowMaxBot > 0 && tr.Bot > policy.UnknownFlowMaxBot {
		return false
	}
	return true
}

func scoredFlowEligible(tr tapeTrade, policy alertPolicy) bool {
	if alertMode(tr) != "OBSERVE" {
		return false
	}
	if !hasScoredWalletMeta(tr) {
		return false
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		_ = reason
		return false
	}
	if strings.EqualFold(tr.Tier, "BOT") {
		return false
	}
	if policy.ScoredFlowMaxBot > 0 && tr.Bot > policy.ScoredFlowMaxBot {
		return false
	}
	return tierAtLeast(tr.Tier, policy.ScoredFlowMinTier)
}

func seedFlowEligible(tr tapeTrade, policy alertPolicy) bool {
	mode := alertMode(tr)
	if mode != "OBSERVE" && mode != "EDGE-HOT" {
		return false
	}
	if strings.EqualFold(tr.KnownList, "consensus_research") {
		return false
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		_ = reason
		return false
	}
	if strings.EqualFold(tr.Tier, "BOT") {
		return false
	}
	if policy.UnknownFlowMaxBot > 0 && tr.Bot > policy.UnknownFlowMaxBot {
		return false
	}
	if strings.TrimSpace(tr.Tier) != "" && !tierAtLeast(tr.Tier, "B") {
		return false
	}
	return true
}

func unknownFlowTrade(c unknownFlowCandidate) tapeTrade {
	return tapeTrade{
		Time:        c.LastTime,
		Wallet:      c.Wallet,
		Side:        "BUY",
		Notional:    c.TotalNotional,
		Price:       c.LastPrice,
		Outcome:     c.Outcome,
		Market:      c.Market,
		Slug:        c.Slug,
		Category:    c.Category,
		Asset:       c.Asset,
		Transaction: unknownFlowKey(c),
		KnownList:   "unknown_flow",
	}
}

func scoredFlowTrade(c unknownFlowCandidate) tapeTrade {
	return tapeTrade{
		Time:        c.LastTime,
		Wallet:      c.Wallet,
		Side:        "BUY",
		Notional:    c.TotalNotional,
		Price:       c.LastPrice,
		Outcome:     c.Outcome,
		Market:      c.Market,
		Slug:        c.Slug,
		Category:    c.Category,
		Asset:       c.Asset,
		Transaction: scoredFlowKey(c),
		KnownList:   "scored_flow",
		Tier:        c.Tier,
		Bot:         c.Bot,
	}
}

func unknownFlowKey(c unknownFlowCandidate) string {
	return strings.ToLower(fmt.Sprintf("unknown-flow|%s|%d|%d|%d|%.4f", c.Wallet, c.LastTime.Unix(), c.Trades, c.Markets, c.TotalNotional))
}

func seedFlowTrade(c unknownFlowCandidate) tapeTrade {
	return tapeTrade{
		Time:        c.LastTime,
		Wallet:      c.Wallet,
		Side:        "BUY",
		Notional:    c.TotalNotional,
		Price:       c.LastPrice,
		Outcome:     c.Outcome,
		Market:      c.Market,
		Slug:        c.Slug,
		Category:    c.Category,
		Asset:       c.Asset,
		Transaction: seedFlowKey(c),
		KnownList:   "seed_flow",
		Tier:        c.Tier,
		Bot:         c.Bot,
	}
}

func seedFlowKey(c unknownFlowCandidate) string {
	return strings.ToLower(fmt.Sprintf("seed-flow|%s|%d|%d|%d|%.4f", c.Wallet, c.LastTime.Unix(), c.Trades, c.Markets, c.TotalNotional))
}

func scoredFlowKey(c unknownFlowCandidate) string {
	return strings.ToLower(fmt.Sprintf("scored-flow|%s|%d|%d|%d|%.4f", c.Wallet, c.LastTime.Unix(), c.Trades, c.Markets, c.TotalNotional))
}

func unknownFlowMarketKey(tr tapeTrade) string {
	for _, v := range []string{tr.Asset, tr.Slug, tr.Market} {
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "" {
			return v
		}
	}
	return strings.ToLower(fmt.Sprintf("%s|%d|%.4f", tr.Wallet, tr.Time.Unix(), tr.Notional))
}

func consensusTrades(trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, maxAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildConsensusCandidates(trades, st, policy, now, burstScanAge(maxAge, burstWindow), maxAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadySent || strings.HasPrefix(c.Reason, "stale") || c.VWAP <= 0 {
			continue
		}
		out = append(out, consensusTrade(c))
	}
	return out
}

func burstScanAge(alertAge, burstWindow time.Duration) time.Duration {
	if alertAge <= 0 || burstWindow <= 0 {
		return alertAge
	}
	return alertAge + burstWindow
}

func staleConsensusTrades(trades []tapeTrade, st sentState, policy alertPolicy, logged map[string]struct{}, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []tapeTrade {
	candidates := buildConsensusCandidates(trades, st, policy, now, diagnosticAge, alertAge, burstWindow, minTrades, minLegNotional)
	out := make([]tapeTrade, 0, len(candidates))
	for _, c := range candidates {
		if c.AlreadySent || !strings.HasPrefix(c.Reason, "stale") || c.VWAP <= 0 {
			continue
		}
		tr := consensusTrade(c)
		if _, ok := logged[tradeKey(tr)]; ok {
			continue
		}
		out = append(out, tr)
	}
	return out
}

func buildConsensusCandidates(trades []tapeTrade, st sentState, policy alertPolicy, now time.Time, diagnosticAge, alertAge, burstWindow time.Duration, minTrades int, minLegNotional float64) []consensusCandidate {
	groups := map[string][]tapeTrade{}
	for _, tr := range trades {
		if tr.Side != "BUY" || tr.Asset == "" || tr.Price <= 0 || tr.Notional < minLegNotional || tr.Size <= 0 {
			continue
		}
		age := now.Sub(tr.Time)
		if diagnosticAge > 0 && (age < 0 || age > diagnosticAge) {
			continue
		}
		if !consensusEligible(tr, policy) {
			continue
		}
		groups[tr.Asset] = append(groups[tr.Asset], tr)
	}

	var out []consensusCandidate
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
		for end := len(rows) - 1; end >= 0; end-- {
			start := end
			for start > 0 && rows[end].Time.Sub(rows[start-1].Time) <= burstWindow {
				start--
			}
			window := rows[start : end+1]
			if len(window) < minTrades {
				continue
			}
			wallets := map[string]struct{}{}
			var totalNotional, totalSize, maxBot float64
			for _, tr := range window {
				wallets[strings.ToLower(tr.Wallet)] = struct{}{}
				totalNotional += tr.Notional
				totalSize += tr.Size
				if tr.Bot > maxBot {
					maxBot = tr.Bot
				}
			}
			if policy.ConsensusMinWallets > 0 && len(wallets) < policy.ConsensusMinWallets {
				continue
			}
			if totalNotional < policy.ConsensusMinNotional || totalSize <= 0 {
				continue
			}
			first := window[0]
			last := window[len(window)-1]
			age := now.Sub(last.Time)
			reason := "consensus burst in diagnostic window"
			if alertAge > 0 && age > alertAge {
				reason = fmt.Sprintf("stale consensus: last age %s > alert window %s", formatDuration(age), alertAge)
			}
			c := consensusCandidate{
				Asset:         last.Asset,
				Category:      last.Category,
				Market:        last.Market,
				Outcome:       last.Outcome,
				Slug:          last.Slug,
				Trades:        len(window),
				Wallets:       len(wallets),
				Participants:  sortedWallets(wallets),
				TotalNotional: totalNotional,
				TotalSize:     totalSize,
				VWAP:          totalNotional / totalSize,
				LastNotional:  last.Notional,
				FirstTime:     first.Time,
				LastTime:      last.Time,
				Tier:          strongestTier(window),
				Bot:           maxBot,
				Reason:        reason,
			}
			c.AlreadySent = consensusAlreadySent(st, c)
			out = append(out, c)
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AlreadySent != out[j].AlreadySent {
			return !out[i].AlreadySent
		}
		if out[i].TotalNotional != out[j].TotalNotional {
			return out[i].TotalNotional > out[j].TotalNotional
		}
		return out[i].LastTime.After(out[j].LastTime)
	})
	return out
}

func consensusEligible(tr tapeTrade, policy alertPolicy) bool {
	switch alertMode(tr) {
	case "REVIEW-NOISE", "REVERSAL-RISK":
		return false
	}
	if strings.EqualFold(tr.Tier, "BOT") {
		return false
	}
	if reason, blocked := policy.EdgeBlocks[strings.ToLower(tr.Wallet)]; blocked {
		_ = reason
		return false
	}
	if policy.ConsensusMaxBot > 0 && tr.Bot > policy.ConsensusMaxBot {
		return false
	}
	return true
}

func consensusTrade(c consensusCandidate) tapeTrade {
	return tapeTrade{
		Time:         c.LastTime,
		Wallet:       fmt.Sprintf("multi:%d", c.Wallets),
		Side:         "BUY",
		Notional:     c.TotalNotional,
		Price:        c.VWAP,
		Outcome:      c.Outcome,
		Market:       c.Market,
		Slug:         c.Slug,
		Category:     c.Category,
		Asset:        c.Asset,
		Transaction:  consensusKey(c),
		KnownList:    "consensus",
		Tier:         c.Tier,
		Bot:          c.Bot,
		Participants: append([]string(nil), c.Participants...),
	}
}

func consensusKey(c consensusCandidate) string {
	return strings.ToLower(fmt.Sprintf("consensus|%s|%d|%d|%.4f", c.Asset, c.LastTime.Unix(), c.Wallets, c.TotalNotional))
}

func consensusAlreadySent(st sentState, c consensusCandidate) bool {
	if st.Sent == nil {
		return false
	}
	_, ok := st.Sent[tradeKey(consensusTrade(c))]
	return ok
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

func tierAtLeast(tier, minTier string) bool {
	minTier = strings.ToUpper(strings.TrimSpace(minTier))
	if minTier == "" {
		return true
	}
	return tierRank(strings.ToUpper(strings.TrimSpace(tier))) >= tierRank(minTier) && tierRank(minTier) > 0
}

func isBurstMode(mode string) bool {
	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case "FOLLOW-READY", "CANDIDATE", "PROBATION", "EDGE-HOT", "FLOW-SCOUT":
		return true
	default:
		return false
	}
}

func hasSentAssetWallet(st sentState, asset, wallet string) bool {
	_, ok := lastSentByAssetWallet(st)[assetWalletKey(asset, wallet)]
	return ok
}

func lastSentByAssetWallet(st sentState) map[string]time.Time {
	out := map[string]time.Time{}
	for key, rawTs := range st.Sent {
		parts := strings.Split(strings.ToLower(key), "|")
		if len(parts) < 3 {
			continue
		}
		asset := parts[len(parts)-2]
		wallet := parts[len(parts)-1]
		posKey := assetWalletKey(asset, wallet)
		if posKey == "|" {
			continue
		}
		ts, err := time.Parse(time.RFC3339, rawTs)
		if err != nil {
			continue
		}
		if prev, ok := out[posKey]; !ok || ts.After(prev) {
			out[posKey] = ts
		}
	}
	return out
}

func positionCooldownActive(lastByPosition map[string]time.Time, tr tapeTrade, now time.Time, policy alertPolicy) (bool, string) {
	if policy.PositionCooldown <= 0 {
		return false, ""
	}
	if policy.RepeatMinNotional > 0 && tr.Notional >= policy.RepeatMinNotional {
		return false, ""
	}
	last, ok := lastByPosition[assetWalletKey(tr.Asset, tr.Wallet)]
	if !ok {
		return false, ""
	}
	age := now.Sub(last)
	if age < 0 {
		age = 0
	}
	if age >= policy.PositionCooldown {
		return false, ""
	}
	return true, fmt.Sprintf("position cooldown: same wallet+asset alerted %s ago < %s", formatDuration(age), policy.PositionCooldown)
}

func assetWalletKey(asset, wallet string) string {
	wallet = strings.ToLower(strings.TrimSpace(wallet))
	asset = strings.TrimSpace(asset)
	return strings.ToLower(asset) + "|" + wallet
}

func sentMode(st sentState, key, fallback string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" {
		return fallback
	}
	if detail, ok := st.Details[key]; ok {
		mode := strings.ToUpper(strings.TrimSpace(detail.Mode))
		if mode != "" {
			return mode
		}
	}
	return fallback
}

func writeBurstSection(b *strings.Builder, rows []burstCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Accumulation Bursts\n\n")
	fmt.Fprintf(b, "- Rule: same wallet + same asset, strategy mode only, cumulative BUY notional within burst window\n\n")
	fmt.Fprintf(b, "| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none |  |  | 0 | $0 | $0 |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "burst"
		if r.AlreadySent {
			status = "covered"
		} else if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | %s | `%s` | %d | $%.0f | $%.0f | %s | %s | %s |\n",
			status, r.Mode, shortAddr(r.Wallet), r.Trades, r.TotalNotional, r.LastNotional,
			formatDuration(now.Sub(r.LastTime)), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeObserveBurstSection(b *strings.Builder, rows []burstCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Observe Accumulation Bursts\n\n")
	fmt.Fprintf(b, "- Rule: same wallet + same asset, OBSERVE mode only, cumulative BUY notional within burst window; requires scored low-bot wallet unless insider threshold applies separately; consensus_research wallets remain single-wallet blocked and require a CONSENSUS burst\n\n")
	fmt.Fprintf(b, "| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none |  |  | 0 | $0 | $0 |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "observe-burst"
		if r.AlreadySent {
			status = "covered"
		} else if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | %s | `%s` | %d | $%.0f | $%.0f | %s | %s | %s |\n",
			status, r.Mode, shortAddr(r.Wallet), r.Trades, r.TotalNotional, r.LastNotional,
			formatDuration(now.Sub(r.LastTime)), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeConsensusSection(b *strings.Builder, rows []consensusCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Consensus Bursts\n\n")
	fmt.Fprintf(b, "- Rule: same asset across wallets, cumulative BUY notional within burst window; excludes review-noise, reversal-risk, negative-edge, and BOT tiers\n\n")
	fmt.Fprintf(b, "| Status | Wallets | Trades | Total | VWAP | LastAge | Participants | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none | 0 | 0 | $0 | 0.000 |  |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "consensus"
		if r.AlreadySent {
			status = "covered"
		} else if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | %d | %d | $%.0f | %.3f | %s | %s | %s | %s |\n",
			status, r.Wallets, r.Trades, r.TotalNotional, r.VWAP, formatDuration(now.Sub(r.LastTime)),
			formatParticipants(r.Participants, 8), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeUnknownFlowSection(b *strings.Builder, rows []unknownFlowCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Unknown Multi-Market Flow\n\n")
	fmt.Fprintf(b, "- Rule: shadow-only unknown wallets buying multiple target markets within the burst window; used to collect edge before any Telegram promotion.\n\n")
	fmt.Fprintf(b, "| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none |  | 0 | 0 | $0 | $0 |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "unknown-flow"
		if r.AlreadyLogged {
			status = "covered"
		} else if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | `%s` | %d | %d | $%.0f | $%.0f | %s | %s | %s |\n",
			status, shortAddr(r.Wallet), r.Trades, r.Markets, r.TotalNotional, r.LastNotional,
			formatDuration(now.Sub(r.LastTime)), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeSeedFlowSection(b *strings.Builder, rows []unknownFlowCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Seed Multi-Market Flow\n\n")
	fmt.Fprintf(b, "- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; used to catch early sports flow before UNKNOWN-FLOW size is reached.\n\n")
	fmt.Fprintf(b, "| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none |  | 0 | 0 | $0 | $0 |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "seed-flow"
		if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | `%s` | %d | %d | $%.0f | $%.0f | %s | %s | %s |\n",
			status, shortAddr(r.Wallet), r.Trades, r.Markets, r.TotalNotional, r.LastNotional,
			formatDuration(now.Sub(r.LastTime)), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
}

func writeScoredFlowSection(b *strings.Builder, rows []unknownFlowCandidate, now time.Time, limit int) {
	fmt.Fprintf(b, "## Scored Multi-Market Flow\n\n")
	fmt.Fprintf(b, "- Rule: shadow-only scored low-bot wallets buying multiple target markets within the burst window; used to find leaderboard whale flow before any Telegram promotion.\n\n")
	fmt.Fprintf(b, "| Status | Wallet | Tier | Bot | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---:|---:|---|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| none |  |  | 0.0 | 0 | 0 | $0 | $0 |  |  |  |\n\n")
		return
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	for _, r := range rows {
		status := "scored-flow"
		if r.AlreadyLogged {
			status = "covered"
		} else if strings.HasPrefix(r.Reason, "stale") {
			status = "stale"
		}
		fmt.Fprintf(b, "| %s | `%s` | %s | %.1f | %d | %d | $%.0f | $%.0f | %s | %s | %s |\n",
			status, shortAddr(r.Wallet), dash(r.Tier), r.Bot, r.Trades, r.Markets, r.TotalNotional, r.LastNotional,
			formatDuration(now.Sub(r.LastTime)), oneLine(r.Reason, 60), oneLine(r.Market, 80))
	}
	fmt.Fprintf(b, "\n")
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

func ensureRequiredAlertModes(modes map[string]struct{}, required ...string) {
	for _, mode := range required {
		mode = strings.ToUpper(strings.TrimSpace(mode))
		if mode != "" {
			modes[mode] = struct{}{}
		}
	}
}

func markTrades(st sentState, trades []tapeTrade, now time.Time) {
	if st.Sent == nil {
		st.Sent = map[string]string{}
	}
	if st.Details == nil {
		st.Details = map[string]sentDetail{}
	}
	ts := now.Format(time.RFC3339)
	for _, tr := range trades {
		key := tradeKey(tr)
		st.Sent[key] = ts
		st.Details[key] = sentDetail{Mode: alertMode(tr), KnownList: tr.KnownList}
	}
}

func cloneSentState(st sentState) sentState {
	out := sentState{
		Sent:    map[string]string{},
		Details: map[string]sentDetail{},
	}
	for k, v := range st.Sent {
		out.Sent[k] = v
	}
	for k, v := range st.Details {
		out.Details[k] = v
	}
	return out
}

func tradeKey(tr tapeTrade) string {
	if tr.Transaction != "" {
		return strings.ToLower(strings.Join([]string{tr.Transaction, tr.Asset, tr.Wallet}, "|"))
	}
	return strings.ToLower(fmt.Sprintf("%s|%s|%d|%.8f|%.4f", tr.Wallet, tr.Asset, tr.Time.Unix(), tr.Price, tr.Notional))
}

func formatAlert(tr tapeTrade) string {
	list := dash(firstNonEmpty(tr.KnownList, "observe"))
	tier := dash(tr.Tier)
	mode := alertMode(tr)
	note := alertNote(tr)
	lines := []string{
		"Polymarket sports whale tape",
		fmt.Sprintf("%s | %s | %s", mode, strings.ToUpper(tr.Category), tr.Time.Format("15:04:05 MST")),
		fmt.Sprintf("Wallet: %s  list=%s tier=%s bot=%.1f", shortAddr(tr.Wallet), list, tier, tr.Bot),
		fmt.Sprintf("Order: BUY $%.0f @ %.3f  outcome=%s", tr.Notional, tr.Price, oneLine(tr.Outcome, 48)),
		fmt.Sprintf("Market: %s", oneLine(tr.Market, 96)),
		fmt.Sprintf("TargetCopy: %.1f%% / %d", tr.TargetCopyROI, tr.TargetCopyT),
		note,
	}
	if tr.Slug != "" {
		lines = append(lines, "Link: https://polymarket.com/event/"+tr.Slug)
	}
	return strings.Join(lines, "\n")
}

func alertMode(tr tapeTrade) string {
	if strings.EqualFold(tr.KnownList, "consensus") {
		return "CONSENSUS"
	}
	if strings.EqualFold(tr.KnownList, "observe_burst") {
		return "OBSERVE-BURST"
	}
	if strings.EqualFold(tr.KnownList, "unknown_flow") {
		return "UNKNOWN-FLOW"
	}
	if strings.EqualFold(tr.KnownList, "seed_flow") {
		return "SEED-FLOW"
	}
	if strings.EqualFold(tr.KnownList, "scored_flow") {
		return "SCORED-FLOW"
	}
	if strings.EqualFold(tr.KnownList, "review_noise") {
		return "REVIEW-NOISE"
	}
	if strings.EqualFold(tr.KnownList, "tape_reversal") {
		return "REVERSAL-RISK"
	}
	if strings.EqualFold(tr.KnownList, "tape_follow") {
		return "FOLLOW-READY"
	}
	if strings.EqualFold(tr.KnownList, "tape_candidate") {
		return "CANDIDATE"
	}
	if strings.EqualFold(tr.KnownList, "tape_probation") {
		return "PROBATION"
	}
	if strings.EqualFold(tr.KnownList, "flow") {
		return "FLOW-SCOUT"
	}
	if strings.EqualFold(tr.KnownList, "insider_scout") {
		return "INSIDER-SCOUT"
	}
	if isEdgeHotList(tr.KnownList) {
		return "EDGE-HOT"
	}
	return "OBSERVE"
}

func alertNote(tr tapeTrade) string {
	if strings.EqualFold(tr.KnownList, "consensus") {
		return "Note: cross-wallet consensus burst; research/observation alert only until repeated positive ROI is proven."
	}
	if strings.EqualFold(tr.KnownList, "observe_burst") {
		return "Note: same-wallet split whale accumulation; observation only until repeated edge and ROI are proven."
	}
	if strings.EqualFold(tr.KnownList, "unknown_flow") {
		return "Note: unknown wallet multi-market sports flow; shadow-only until repeated edge and ROI are proven."
	}
	if strings.EqualFold(tr.KnownList, "seed_flow") {
		return "Note: early unknown wallet multi-market sports flow; research-only and below Telegram promotion thresholds."
	}
	if strings.EqualFold(tr.KnownList, "scored_flow") {
		return "Note: scored low-bot leaderboard whale multi-market sports flow; shadow-only until repeated edge and ROI are proven."
	}
	if strings.EqualFold(tr.KnownList, "consensus_research") {
		return "Note: wallet appeared in a positive consensus burst; research-only for single-wallet orders until repeated consensus edge is proven."
	}
	if strings.EqualFold(tr.KnownList, "tape_reversal") {
		return "Note: late-window edge reversed; observe only and do not treat as follow-ready."
	}
	if strings.EqualFold(tr.KnownList, "tape_follow") {
		return "Note: strict edge gates passed; still size manually until proven ROI is established."
	}
	if strings.EqualFold(tr.KnownList, "tape_probation") {
		return "Note: positive target-copy but soft flow risks; edge observation only until 5m/15m windows prove out."
	}
	if strings.EqualFold(tr.KnownList, "flow") {
		return "Note: recent leaderboard flow scout; observation alert only until copy and edge samples prove out."
	}
	if strings.EqualFold(tr.KnownList, "insider_scout") {
		return "Note: very large low-bot sports whale; observation only until repeated edge and ROI are proven."
	}
	if isEdgeHotList(tr.KnownList) {
		return "Note: measured early edge is positive; still observation-sized until more alerts prove ROI."
	}
	return "Note: observation alert only; wait for edge/proven promotion before treating as follow-ready."
}

func isEdgeHotList(list string) bool {
	switch strings.ToLower(strings.TrimSpace(list)) {
	case "watch", "target", "sports", "scout", "tape_edgehot":
		return true
	default:
		return false
	}
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
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

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
