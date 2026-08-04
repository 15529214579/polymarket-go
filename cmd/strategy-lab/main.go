package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

type strategyParams struct {
	MinTier            string
	MaxBot             float64
	MinCopyTrades      int
	MinCopyROI         float64
	MinCopyPnL         float64
	MinCopyWinRate     float64
	MinClosedROI       float64
	MinSmartMoneyScore float64
}

type strategyResult struct {
	Params           strategyParams
	Wallets          []walletdiscover.WalletScore
	WalletCount      int
	TotalCopyTrades  int
	TotalCopyWins    int
	TotalCopyPnL     float64
	TotalCopyCapital float64
	TotalCopyROI     float64
	TotalCopyWinRate float64
	TotalOpenCost    float64
	OpenCostRatio    float64
	MedianCopyROI    float64
	WorstCopyROI     float64
	Score            float64
}

type watchParams struct {
	MaxWallets      int
	MaxBot          float64
	MinSmart        float64
	MinCopyTrades   int
	MinCopyROI      float64
	MinCopyPnL      float64
	MinCopyWinRate  float64
	MinLargeTrades  int
	MinAvgNotional  float64
	ExcludeRiskTags map[string]struct{}
}

type sportsParams struct {
	MaxWallets       int
	MaxBot           float64
	MinSmart         float64
	MinEdge          float64
	MinTargetTrades  int
	MinTargetRatio   float64
	MinTargetCopyT   int
	MinTargetCopyROI float64
	MinTargetCopyPnL float64
	MinTargetLarge   int
	ExcludeRiskTags  map[string]struct{}
}

type scoutParams struct {
	MaxWallets       int
	MaxBot           float64
	MinSmart         float64
	MinEdge          float64
	MinRecentTrades  int
	MinTargetTrades  int
	MinTargetRatio   float64
	MinTargetCopyT   int
	MinTargetCopyROI float64
	MinLargeTrades   int
	MinAvgNotional   float64
	ExcludeRiskTags  map[string]struct{}
}

type tapeParams struct {
	MaxWallets       int
	MaxBot           float64
	MinDirectMaxBuy  float64
	MinScoredMaxBuy  float64
	MinSmart         float64
	MinTargetCopyT   int
	MinTargetCopyROI float64
	ExcludeRiskTags  map[string]struct{}
	EdgeBlocks       map[string]edgeBlock
}

type tapeInput struct {
	Address     string
	Tier        string
	Buys        int
	BuyNotional float64
	MaxBuy      float64
	Bot         float64
}

type tapeEdgeHotParams struct {
	MaxWallets    int
	MaxBot        float64
	MinMaxBuy     float64
	MinSamples    int
	MinAvgPP      float64
	MinWinRate    float64
	Min5mAvgPP    float64
	Min15mAvgPP   float64
	Max1hNegPP    float64
	EdgeProfiles  map[string]edgeProfile
	EdgeBlocks    map[string]edgeBlock
	ExcludedRisks map[string]struct{}
}

type edgeBlock struct {
	Samples15m int
	Avg15mPP   float64
	Samples1h  int
	Avg1hPP    float64
	Reason     string
}

type edgeProfile struct {
	TotalSamples int
	TotalWins    int
	AvgPP        float64
	WinRate      float64
	MaxNotional  float64
	TapeAction   bool
	Samples5m    int
	Avg5mPP      float64
	Samples15m   int
	Avg15mPP     float64
	Samples1h    int
	Avg1hPP      float64
	Reason       string
}

func main() {
	scoresPath := flag.String("scores", "db/strategy_iteration/wallet_scores.json", "wallet score JSON from wallet-discover")
	reportPath := flag.String("report", "reports/strategy_lab.md", "strategy lab report path")
	excludeWalletsPath := flag.String("exclude_wallets", "", "wallet file to exclude from strategy generation")
	coreWalletsPath := flag.String("core_wallets", "wallets.strategy-core.txt", "core wallet list path")
	watchWalletsPath := flag.String("watch_wallets", "wallets.strategy-watch.txt", "active watchlist wallet path")
	sportsWalletsPath := flag.String("sports_wallets", "wallets.strategy-sports.txt", "target-category positive copy wallet path")
	scoutWalletsPath := flag.String("scout_wallets", "wallets.strategy-scout.txt", "leaderboard-only scout wallet path")
	targetWalletsPath := flag.String("target_wallets", "wallets.strategy-target.txt", "target-category whale scout wallet path")
	flowWalletsPath := flag.String("flow_wallets", "wallets.strategy-flow.txt", "recent-trade flow scout wallet path")
	tapeWalletsPath := flag.String("tape_wallets", "wallets.strategy-tape.txt", "recent sports-tape hot whale wallet path")
	tapeProbationWalletsPath := flag.String("tape_probation_wallets", "wallets.strategy-tape-probation.txt", "scored sports-tape wallets with positive target-copy but soft flow risks; edge observation only")
	tapeObserveWalletsPath := flag.String("tape_observe_wallets", "wallets.strategy-tape-observe.txt", "direct sports-tape whale observation-only wallet path")
	tapeEdgeHotWalletsPath := flag.String("tape_edge_hot_wallets", "wallets.strategy-tape-edgehot.txt", "sports-tape wallets with positive measured short-window edge")
	sportsTapeInputPath := flag.String("sports_tape_wallets", "wallets.sports-tape.txt", "recent target-category large-order wallet file from sports-tape")
	pushWalletsPath := flag.String("push_wallets", "wallets.strategy-push.txt", "combined push wallet path")
	minWallets := flag.Int("min_wallets", 5, "minimum wallets in a valid strategy")
	minTotalCopyTrades := flag.Int("min_total_copy_trades", 80, "minimum aggregate closed copy-simulation trades")
	floorMinCopyTrades := flag.Int("floor_min_copy_trades", 3, "hard lower bound for per-wallet closed copy-simulation trades")
	floorMinCopyROI := flag.Float64("floor_min_copy_roi", 0, "hard lower bound for per-wallet copy-simulation ROI")
	floorMinCopyPnL := flag.Float64("floor_min_copy_pnl", 0, "hard lower bound for per-wallet copy-simulation PnL")
	floorMinCopyWinRate := flag.Float64("floor_min_copy_win_rate", 0, "hard lower bound for per-wallet copy-simulation win rate")
	floorMinClosedROI := flag.Float64("floor_min_closed_roi", 0, "hard lower bound for wallet closed-position ROI")
	floorMinSmartMoneyScore := flag.Float64("floor_min_smart", 70, "hard lower bound for smart-money score")
	ceilingMaxBot := flag.Float64("ceiling_max_bot", 35, "hard upper bound for bot score")
	topN := flag.Int("top_n", 12, "number of strategies to show in the report")
	watchMaxWallets := flag.Int("watch_max_wallets", 20, "maximum non-core active watchlist wallets")
	watchMaxBot := flag.Float64("watch_max_bot", 35, "maximum bot score for active watchlist")
	watchMinSmart := flag.Float64("watch_min_smart", 60, "minimum smart-money score for active watchlist")
	watchMinCopyTrades := flag.Int("watch_min_copy_trades", 3, "minimum closed copy-simulation trades for active watchlist")
	watchMinCopyROI := flag.Float64("watch_min_copy_roi", 5, "minimum copy-simulation ROI for active watchlist")
	watchMinCopyPnL := flag.Float64("watch_min_copy_pnl", 10, "minimum copy-simulation PnL for active watchlist")
	watchMinCopyWinRate := flag.Float64("watch_min_copy_win_rate", 50, "minimum copy-simulation win rate for active watchlist")
	watchMinLargeTrades := flag.Int("watch_min_large_trades", 5, "minimum large-trade sample for active watchlist")
	watchMinAvgNotional := flag.Float64("watch_min_avg_notional", 250, "minimum average trade notional for active watchlist")
	sportsMaxWallets := flag.Int("sports_max_wallets", 10, "maximum target-category positive copy wallets")
	sportsMaxBot := flag.Float64("sports_max_bot", 35, "maximum bot score for target-category positive copy wallets")
	sportsMinSmart := flag.Float64("sports_min_smart", 70, "minimum smart-money score for target-category positive copy wallets")
	sportsMinEdge := flag.Float64("sports_min_edge", 50, "minimum edge score for target-category positive copy wallets")
	sportsMinTargetTrades := flag.Int("sports_min_target_trades", 10, "minimum target-category trade sample for sports layer")
	sportsMinTargetRatio := flag.Float64("sports_min_target_ratio", 0.05, "minimum target-category trade ratio for sports layer")
	sportsMinTargetCopyTrades := flag.Int("sports_min_target_copy_trades", 5, "minimum target-category closed copy-simulation trades for sports layer")
	sportsMinTargetCopyROI := flag.Float64("sports_min_target_copy_roi", 5, "minimum target-category copy-simulation ROI for sports layer")
	sportsMinTargetCopyPnL := flag.Float64("sports_min_target_copy_pnl", 5, "minimum target-category copy-simulation PnL for sports layer")
	sportsMinTargetLarge := flag.Int("sports_min_target_large", 5, "minimum target-category large trades for sports layer")
	scoutMaxWallets := flag.Int("scout_max_wallets", 10, "maximum leaderboard-only scout wallets")
	scoutMaxBot := flag.Float64("scout_max_bot", 45, "maximum bot score for leaderboard scout")
	scoutMinSmart := flag.Float64("scout_min_smart", 80, "minimum smart-money score for leaderboard scout")
	scoutMinEdge := flag.Float64("scout_min_edge", 60, "minimum edge score for leaderboard scout")
	scoutMinLargeTrades := flag.Int("scout_min_large_trades", 50, "minimum large-trade sample for leaderboard scout")
	scoutMinAvgNotional := flag.Float64("scout_min_avg_notional", 500, "minimum average trade notional for leaderboard scout")
	scoutPushEnabled := flag.Bool("scout_push_enabled", false, "include leaderboard-only scout wallets in whale push; default keeps them research-only")
	targetMaxWallets := flag.Int("target_max_wallets", 10, "maximum target-category scout wallets")
	targetMaxBot := flag.Float64("target_max_bot", 45, "maximum bot score for target-category scout")
	targetMinSmart := flag.Float64("target_min_smart", 70, "minimum smart-money score for target-category scout")
	targetMinEdge := flag.Float64("target_min_edge", 40, "minimum edge score for target-category scout")
	targetMinTrades := flag.Int("target_min_trades", 50, "minimum target-category trade sample")
	targetMinRatio := flag.Float64("target_min_ratio", 0.20, "minimum target-category trade ratio")
	targetMinCopyTrades := flag.Int("target_min_copy_trades", 5, "minimum target-category closed copy-simulation trades")
	targetMinCopyROI := flag.Float64("target_min_copy_roi", 10, "minimum target-category copy-simulation ROI")
	targetMinLargeTrades := flag.Int("target_min_large_trades", 50, "minimum target-category large-trade sample")
	targetMinAvgNotional := flag.Float64("target_min_avg_notional", 100, "minimum average trade notional for target-category scout")
	flowMaxWallets := flag.Int("flow_max_wallets", 15, "maximum recent-trade flow scout wallets")
	flowMaxBot := flag.Float64("flow_max_bot", 45, "maximum bot score for recent-trade flow scout")
	flowMinSmart := flag.Float64("flow_min_smart", 70, "minimum smart-money score for recent-trade flow scout")
	flowMinEdge := flag.Float64("flow_min_edge", 35, "minimum edge score for recent-trade flow scout")
	flowMinRecentTrades := flag.Int("flow_min_recent_trades", 2, "minimum qualifying recent-trade source hits for recent-trade flow scout")
	flowMinLargeTrades := flag.Int("flow_min_large_trades", 5, "minimum large-trade sample for recent-trade flow scout")
	flowMinAvgNotional := flag.Float64("flow_min_avg_notional", 500, "minimum average trade notional for recent-trade flow scout")
	tapeMaxWallets := flag.Int("tape_max_wallets", 8, "maximum recent sports-tape hot wallets")
	tapeMaxBot := flag.Float64("tape_max_bot", 45, "maximum bot score for scored sports-tape hot wallets")
	tapeMinDirectMaxBuy := flag.Float64("tape_min_direct_max_buy", 5000, "minimum single recent sports-tape BUY to include even before full scoring")
	tapeObserveMinBuyNotional := flag.Float64("tape_observe_min_buy_notional", 8000, "minimum cumulative recent sports-tape BUY notional to include scored low-bot wallets in observe-only list")
	tapeMinScoredMaxBuy := flag.Float64("tape_min_scored_max_buy", 2500, "minimum single recent sports-tape BUY to include scored low-bot wallets")
	tapeMinSmart := flag.Float64("tape_min_smart", 70, "minimum smart score for scored sports-tape hot wallets")
	tapeMinTargetCopyTrades := flag.Int("tape_min_target_copy_trades", 2, "minimum target copy trades for positive sports-tape hot wallets")
	tapeMinTargetCopyROI := flag.Float64("tape_min_target_copy_roi", 25, "minimum target copy ROI for positive sports-tape hot wallets")
	tapeEdgeHotMaxWallets := flag.Int("tape_edge_hot_max_wallets", 10, "maximum sports-tape edge-hot wallets")
	tapeEdgeHotMaxBot := flag.Float64("tape_edge_hot_max_bot", 45, "maximum bot score for sports-tape edge-hot wallets")
	tapeEdgeHotMinMaxBuy := flag.Float64("tape_edge_hot_min_max_buy", 500, "minimum recent sports-tape max BUY for edge-hot wallet tracking")
	tapeEdgeHotMinSamples := flag.Int("tape_edge_hot_min_samples", 2, "minimum non-zero horizon edge samples for sports-tape edge-hot wallets")
	tapeEdgeHotMinAvgPP := flag.Float64("tape_edge_hot_min_avg_pp", 2, "minimum non-zero horizon average edge in pp for sports-tape edge-hot wallets")
	tapeEdgeHotMinWinRate := flag.Float64("tape_edge_hot_min_win_rate", 60, "minimum non-zero horizon win rate for sports-tape edge-hot wallets")
	tapeEdgeHotMin5mAvgPP := flag.Float64("tape_edge_hot_min_5m_avg_pp", 0.5, "minimum 5m average edge in pp for sports-tape edge-hot wallets")
	tapeEdgeHotMin15mAvgPP := flag.Float64("tape_edge_hot_min_15m_avg_pp", 0, "minimum 15m average edge in pp for sports-tape edge-hot wallets")
	tapeEdgeHotMax1hNegPP := flag.Float64("tape_edge_hot_max_1h_neg_pp", -5, "do not mark sports-tape edge-hot when 1h edge is at or below this pp value")
	edgeSnapshotsPath := flag.String("edge_snapshots", "db/strategy_iteration/whale_edge_snapshots.jsonl", "whale-edge snapshot JSONL used to block negative measured sports-tape wallets")
	tapeEdgeBlock15mSamples := flag.Int("tape_edge_block_15m_samples", 2, "minimum 15m edge samples needed to block tape hotlist wallets")
	tapeEdgeBlock15mMaxAvgPP := flag.Float64("tape_edge_block_15m_max_avg_pp", -1, "maximum 15m average edge pp allowed before blocking tape hotlist wallets")
	tapeEdgeBlock1hSamples := flag.Int("tape_edge_block_1h_samples", 1, "minimum 1h edge samples needed to block tape hotlist wallets")
	tapeEdgeBlock1hMaxAvgPP := flag.Float64("tape_edge_block_1h_max_avg_pp", -5, "maximum 1h average edge pp allowed before blocking tape hotlist wallets")
	flag.Parse()

	scores, err := loadScores(*scoresPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: %v\n", err)
		os.Exit(1)
	}
	excluded, err := loadWalletSet(*excludeWalletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: load exclude wallets: %v\n", err)
		os.Exit(1)
	}
	if len(excluded) > 0 {
		scores = filterExcluded(scores, excluded)
	}
	floors := strategyParams{
		MaxBot:             *ceilingMaxBot,
		MinCopyTrades:      *floorMinCopyTrades,
		MinCopyROI:         *floorMinCopyROI,
		MinCopyPnL:         *floorMinCopyPnL,
		MinCopyWinRate:     *floorMinCopyWinRate,
		MinClosedROI:       *floorMinClosedROI,
		MinSmartMoneyScore: *floorMinSmartMoneyScore,
	}
	strategies := evaluateStrategies(scores, *minWallets, *minTotalCopyTrades, floors)
	if len(strategies) == 0 {
		fmt.Fprintf(os.Stderr, "strategy-lab: no valid strategy from %s\n", *scoresPath)
		os.Exit(1)
	}
	best := strategies[0]
	if err := writeCoreWallets(*coreWalletsPath, best.Wallets); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write core wallets: %v\n", err)
		os.Exit(1)
	}
	watch := buildWatchlist(scores, best.Wallets, watchParams{
		MaxWallets:     *watchMaxWallets,
		MaxBot:         *watchMaxBot,
		MinSmart:       *watchMinSmart,
		MinCopyTrades:  *watchMinCopyTrades,
		MinCopyROI:     *watchMinCopyROI,
		MinCopyPnL:     *watchMinCopyPnL,
		MinCopyWinRate: *watchMinCopyWinRate,
		MinLargeTrades: *watchMinLargeTrades,
		MinAvgNotional: *watchMinAvgNotional,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":       {},
			"burst_trading":       {},
			"extreme_price_heavy": {},
			"fixed_amount":        {},
			"fixed_price":         {},
			"negative_copy_sim":   {},
		},
	})
	if err := writeWalletsFile(*watchWalletsPath, watch, "watch"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write watch wallets: %v\n", err)
		os.Exit(1)
	}
	sports := buildSportslist(scores, append([]walletdiscover.WalletScore{}, append(best.Wallets, watch...)...), sportsParams{
		MaxWallets:       *sportsMaxWallets,
		MaxBot:           *sportsMaxBot,
		MinSmart:         *sportsMinSmart,
		MinEdge:          *sportsMinEdge,
		MinTargetTrades:  *sportsMinTargetTrades,
		MinTargetRatio:   *sportsMinTargetRatio,
		MinTargetCopyT:   *sportsMinTargetCopyTrades,
		MinTargetCopyROI: *sportsMinTargetCopyROI,
		MinTargetCopyPnL: *sportsMinTargetCopyPnL,
		MinTargetLarge:   *sportsMinTargetLarge,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeWalletsFile(*sportsWalletsPath, sports, "sports"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write sports wallets: %v\n", err)
		os.Exit(1)
	}
	scout := buildScoutlist(scores, append([]walletdiscover.WalletScore{}, append(append(best.Wallets, watch...), sports...)...), scoutParams{
		MaxWallets:      *scoutMaxWallets,
		MaxBot:          *scoutMaxBot,
		MinSmart:        *scoutMinSmart,
		MinEdge:         *scoutMinEdge,
		MinRecentTrades: 0,
		MinTargetTrades: 0,
		MinTargetRatio:  0,
		MinLargeTrades:  *scoutMinLargeTrades,
		MinAvgNotional:  *scoutMinAvgNotional,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeWalletsFile(*scoutWalletsPath, scout, "scout"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write scout wallets: %v\n", err)
		os.Exit(1)
	}
	target := buildTargetScoutlist(scores, append([]walletdiscover.WalletScore{}, append(append(append(best.Wallets, watch...), sports...), scout...)...), scoutParams{
		MaxWallets:       *targetMaxWallets,
		MaxBot:           *targetMaxBot,
		MinSmart:         *targetMinSmart,
		MinEdge:          *targetMinEdge,
		MinTargetTrades:  *targetMinTrades,
		MinTargetRatio:   *targetMinRatio,
		MinTargetCopyT:   *targetMinCopyTrades,
		MinTargetCopyROI: *targetMinCopyROI,
		MinLargeTrades:   *targetMinLargeTrades,
		MinAvgNotional:   *targetMinAvgNotional,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeWalletsFile(*targetWalletsPath, target, "target"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write target wallets: %v\n", err)
		os.Exit(1)
	}
	flow := buildFlowScoutlist(scores, append([]walletdiscover.WalletScore{}, append(append(append(append(best.Wallets, watch...), sports...), scout...), target...)...), scoutParams{
		MaxWallets:      *flowMaxWallets,
		MaxBot:          *flowMaxBot,
		MinSmart:        *flowMinSmart,
		MinEdge:         *flowMinEdge,
		MinRecentTrades: *flowMinRecentTrades,
		MinTargetTrades: 0,
		MinTargetRatio:  0,
		MinLargeTrades:  *flowMinLargeTrades,
		MinAvgNotional:  *flowMinAvgNotional,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeWalletsFile(*flowWalletsPath, flow, "flow"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write flow wallets: %v\n", err)
		os.Exit(1)
	}
	tapeInputs, err := loadTapeInputs(*sportsTapeInputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: load sports tape wallets: %v\n", err)
		os.Exit(1)
	}
	tapeInputs = filterExcludedTapeInputs(tapeInputs, excluded)
	edgeBlocks, err := loadEdgeBlocks(*edgeSnapshotsPath, *tapeEdgeBlock15mSamples, *tapeEdgeBlock15mMaxAvgPP, *tapeEdgeBlock1hSamples, *tapeEdgeBlock1hMaxAvgPP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: load edge snapshots: %v\n", err)
		os.Exit(1)
	}
	edgeProfiles, err := loadEdgeProfiles(*edgeSnapshotsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: load edge profiles: %v\n", err)
		os.Exit(1)
	}
	selectedBeforeTape := mergeWalletGroups(best.Wallets, watch, sports, scout, target, flow)
	tape := buildTapeHotlist(scores, tapeInputs, selectedBeforeTape, tapeParams{
		MaxWallets:       *tapeMaxWallets,
		MaxBot:           *tapeMaxBot,
		MinDirectMaxBuy:  *tapeMinDirectMaxBuy,
		MinScoredMaxBuy:  *tapeMinScoredMaxBuy,
		MinSmart:         *tapeMinSmart,
		MinTargetCopyT:   *tapeMinTargetCopyTrades,
		MinTargetCopyROI: *tapeMinTargetCopyROI,
		EdgeBlocks:       edgeBlocks,
		ExcludeRiskTags: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeWalletsFile(*tapeWalletsPath, tape, "tape"); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write tape wallets: %v\n", err)
		os.Exit(1)
	}
	tapeProbation := buildTapeProbation(scores, tapeInputs, mergeWalletGroups(selectedBeforeTape, tape), tapeParams{
		MaxWallets:       *tapeMaxWallets,
		MaxBot:           *tapeMaxBot,
		MinDirectMaxBuy:  *tapeMinDirectMaxBuy,
		MinScoredMaxBuy:  *tapeMinScoredMaxBuy,
		MinSmart:         *tapeMinSmart,
		MinTargetCopyT:   *tapeMinTargetCopyTrades,
		MinTargetCopyROI: *tapeMinTargetCopyROI,
		EdgeBlocks:       edgeBlocks,
	})
	if err := writeTapeProbationFile(*tapeProbationWalletsPath, tapeProbation); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write tape probation wallets: %v\n", err)
		os.Exit(1)
	}
	pushCore := filterEdgeBlockedWallets(best.Wallets, edgeBlocks)
	pushWatch := filterEdgeBlockedWallets(watch, edgeBlocks)
	pushSports := filterEdgeBlockedWallets(sports, edgeBlocks)
	var pushScout []walletdiscover.WalletScore
	if *scoutPushEnabled {
		pushScout = filterEdgeBlockedWallets(scout, edgeBlocks)
	}
	pushTarget := filterEdgeBlockedWallets(target, edgeBlocks)
	pushFlow := filterEdgeBlockedWallets(flow, edgeBlocks)
	pushTape := filterEdgeBlockedWallets(tape, edgeBlocks)
	blockedPush := edgeBlockedWallets(mergeWalletGroups(best.Wallets, watch, sports, scout, target, flow, tape), edgeBlocks)
	push := mergeWalletGroups(pushCore, pushWatch, pushSports, pushScout, pushTarget, pushFlow, pushTape)
	if err := writeDirectTapeObserveFile(*tapeObserveWalletsPath, scores, push, tapeInputs, *tapeMinDirectMaxBuy, *tapeObserveMinBuyNotional, *tapeMaxBot); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write tape observe wallets: %v\n", err)
		os.Exit(1)
	}
	tapeEdgeHot := buildTapeEdgeHot(scores, tapeInputs, nil, tapeEdgeHotParams{
		MaxWallets:   *tapeEdgeHotMaxWallets,
		MaxBot:       *tapeEdgeHotMaxBot,
		MinMaxBuy:    *tapeEdgeHotMinMaxBuy,
		MinSamples:   *tapeEdgeHotMinSamples,
		MinAvgPP:     *tapeEdgeHotMinAvgPP,
		MinWinRate:   *tapeEdgeHotMinWinRate,
		Min5mAvgPP:   *tapeEdgeHotMin5mAvgPP,
		Min15mAvgPP:  *tapeEdgeHotMin15mAvgPP,
		Max1hNegPP:   *tapeEdgeHotMax1hNegPP,
		EdgeProfiles: edgeProfiles,
		EdgeBlocks:   edgeBlocks,
		ExcludedRisks: map[string]struct{}{
			"bot_like_flow":     {},
			"fixed_amount":      {},
			"fixed_price":       {},
			"negative_copy_sim": {},
		},
	})
	if err := writeTapeEdgeHotFile(*tapeEdgeHotWalletsPath, tapeEdgeHot, edgeProfiles, tapeInputs); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write tape edge-hot wallets: %v\n", err)
		os.Exit(1)
	}
	if err := writePushWalletsFile(*pushWalletsPath, pushCore, pushWatch, pushSports, pushScout, pushTarget, pushFlow, pushTape); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write push wallets: %v\n", err)
		os.Exit(1)
	}
	if err := writeReport(*reportPath, *scoresPath, *excludeWalletsPath, len(excluded), scores, best, strategies, watch, sports, scout, target, flow, tape, tapeProbation, tapeEdgeHot, push, blockedPush, tapeInputs, edgeProfiles, edgeBlocks, *tapeMinDirectMaxBuy, *scoutPushEnabled, *topN); err != nil {
		fmt.Fprintf(os.Stderr, "strategy-lab: write report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("strategy-lab done: core=%d watch=%d sports=%d scout=%d target=%d flow=%d tape=%d tape_probation=%d push=%d copyT=%d copyROI=%.1f%% copyPnL=$%+.2f core_file=%s push_file=%s report=%s\n",
		best.WalletCount, len(watch), len(sports), len(scout), len(target), len(flow), len(tape), len(tapeProbation), len(push), best.TotalCopyTrades, best.TotalCopyROI, best.TotalCopyPnL,
		*coreWalletsPath, *pushWalletsPath, *reportPath)
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

func loadWalletSet(path string) (map[string]struct{}, error) {
	out := map[string]struct{}{}
	if path == "" {
		return out, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(strings.Split(line, "#")[0])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if strings.HasPrefix(addr, "0x") && len(addr) == 42 {
			out[addr] = struct{}{}
		}
	}
	return out, nil
}

func loadTapeInputs(path string) ([]tapeInput, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []tapeInput
	seen := map[string]struct{}{}
	for _, raw := range strings.Split(string(b), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		body, comment, _ := strings.Cut(line, "#")
		fields := strings.Fields(strings.TrimSpace(body))
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		meta := parseCommentFields(comment)
		out = append(out, tapeInput{
			Address:     addr,
			Tier:        meta["tier"],
			Buys:        parseIntField(meta["buys"]),
			BuyNotional: parseMoneyField(meta["buyNotional"]),
			MaxBuy:      parseMoneyField(meta["maxBuy"]),
			Bot:         parseMoneyField(meta["bot"]),
		})
	}
	return out, nil
}

func filterExcludedTapeInputs(inputs []tapeInput, excluded map[string]struct{}) []tapeInput {
	if len(inputs) == 0 || len(excluded) == 0 {
		return inputs
	}
	out := make([]tapeInput, 0, len(inputs))
	for _, in := range inputs {
		if _, ok := excluded[strings.ToLower(in.Address)]; ok {
			continue
		}
		out = append(out, in)
	}
	return out
}

func loadEdgeBlocks(path string, min15mSamples int, max15mAvgPP float64, min1hSamples int, max1hAvgPP float64) (map[string]edgeBlock, error) {
	out := map[string]edgeBlock{}
	profiles, err := loadEdgeProfiles(path)
	if err != nil {
		return nil, err
	}
	for addr, profile := range profiles {
		block := edgeBlock{
			Samples15m: profile.Samples15m,
			Avg15mPP:   profile.Avg15mPP,
			Samples1h:  profile.Samples1h,
			Avg1hPP:    profile.Avg1hPP,
		}
		switch {
		case min15mSamples > 0 && profile.Samples15m >= min15mSamples && profile.Avg15mPP <= max15mAvgPP:
			block.Reason = fmt.Sprintf("15m edge %.2fpp over %d samples", block.Avg15mPP, block.Samples15m)
		case min1hSamples > 0 && profile.Samples1h >= min1hSamples && profile.Avg1hPP <= max1hAvgPP:
			block.Reason = fmt.Sprintf("1h edge %.2fpp over %d samples", block.Avg1hPP, block.Samples1h)
		}
		if block.Reason != "" {
			out[addr] = block
		}
	}
	return out, nil
}

func loadEdgeProfiles(path string) (map[string]edgeProfile, error) {
	out := map[string]edgeProfile{}
	if path == "" {
		return out, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	type edgeAgg struct {
		samples int
		wins    int
		sum     float64
	}
	type walletAgg struct {
		byHorizon   map[int]*edgeAgg
		maxNotional float64
		tapeAction  bool
	}
	aggs := map[string]*walletAgg{}
	for _, raw := range strings.Split(string(b), "\n") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		var ev struct {
			Wallet     string  `json:"wallet"`
			HorizonSec int     `json:"horizon_sec"`
			DeltaPP    float64 `json:"delta_pp"`
			Notional   float64 `json:"notional_usd"`
			Action     string  `json:"action"`
		}
		if err := json.Unmarshal([]byte(raw), &ev); err != nil {
			continue
		}
		addr := strings.ToLower(strings.TrimSpace(ev.Wallet))
		if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
			continue
		}
		a := aggs[addr]
		if a == nil {
			a = &walletAgg{byHorizon: map[int]*edgeAgg{}}
			aggs[addr] = a
		}
		h := ev.HorizonSec
		st := a.byHorizon[h]
		if st == nil {
			st = &edgeAgg{}
			a.byHorizon[h] = st
		}
		st.samples++
		st.sum += ev.DeltaPP
		if ev.DeltaPP > 0 {
			st.wins++
		}
		if ev.Notional > 0 && ev.Notional > a.maxNotional {
			a.maxNotional = ev.Notional
		}
		if strings.EqualFold(ev.Action, "tape") {
			a.tapeAction = true
		}
	}
	for addr, a := range aggs {
		var total edgeAgg
		profile := edgeProfile{}
		for horizon, st := range a.byHorizon {
			if horizon > 0 {
				total.samples += st.samples
				total.wins += st.wins
				total.sum += st.sum
			}
			avg := 0.0
			if st.samples > 0 {
				avg = st.sum / float64(st.samples)
			}
			switch horizon {
			case 300:
				profile.Samples5m = st.samples
				profile.Avg5mPP = avg
			case 900:
				profile.Samples15m = st.samples
				profile.Avg15mPP = avg
			case 3600:
				profile.Samples1h = st.samples
				profile.Avg1hPP = avg
			}
		}
		profile.TotalSamples = total.samples
		profile.TotalWins = total.wins
		profile.MaxNotional = a.maxNotional
		profile.TapeAction = a.tapeAction
		if total.samples > 0 {
			profile.AvgPP = total.sum / float64(total.samples)
			profile.WinRate = float64(total.wins) / float64(total.samples) * 100
		}
		profile.Reason = fmt.Sprintf("edge-hot %.1f%% win avg %+0.2fpp over %d samples", profile.WinRate, profile.AvgPP, profile.TotalSamples)
		out[addr] = profile
	}
	return out, nil
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

func parseIntField(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return 0
	}
	v, _ := strconv.Atoi(raw)
	return v
}

func parseMoneyField(raw string) float64 {
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

func cleanTapeTier(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return "TAPE"
	}
	return raw
}

func tierAtLeast(tier, minTier string) bool {
	minTier = strings.ToUpper(strings.TrimSpace(minTier))
	if minTier == "" {
		return true
	}
	return tierRank(strings.ToUpper(strings.TrimSpace(tier))) >= tierRank(minTier) && tierRank(minTier) > 0
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

func filterExcluded(scores []walletdiscover.WalletScore, excluded map[string]struct{}) []walletdiscover.WalletScore {
	out := scores[:0]
	for _, s := range scores {
		if _, ok := excluded[strings.ToLower(s.Address)]; ok {
			continue
		}
		out = append(out, s)
	}
	return out
}

func evaluateStrategies(scores []walletdiscover.WalletScore, minWallets, minTotalCopyTrades int, floors strategyParams) []strategyResult {
	var out []strategyResult
	minTiers := []string{"A", "B"}
	maxBots := []float64{20, 25, 30, 35}
	minCopyTrades := []int{3, 5, 8, 10, 15, 20}
	minCopyROIs := []float64{0, 5, 10, 20, 40, 60, 100}
	minCopyPnLs := []float64{0, 5, 10, 25, 50, 100}
	minCopyWinRates := []float64{0, 50, 60, 70, 80}
	minClosedROIs := []float64{0, 5, 10, 20}
	minSmartScores := []float64{70, 85, 95}

	seen := map[string]struct{}{}
	for _, minTier := range minTiers {
		for _, maxBot := range maxBots {
			if maxBot > floors.MaxBot {
				continue
			}
			for _, minT := range minCopyTrades {
				if minT < floors.MinCopyTrades {
					continue
				}
				for _, minROI := range minCopyROIs {
					if minROI < floors.MinCopyROI {
						continue
					}
					for _, minPnL := range minCopyPnLs {
						if minPnL < floors.MinCopyPnL {
							continue
						}
						for _, minWin := range minCopyWinRates {
							if minWin < floors.MinCopyWinRate {
								continue
							}
							for _, minClosedROI := range minClosedROIs {
								if minClosedROI < floors.MinClosedROI {
									continue
								}
								for _, minSmart := range minSmartScores {
									if minSmart < floors.MinSmartMoneyScore {
										continue
									}
									params := strategyParams{
										MinTier:            minTier,
										MaxBot:             maxBot,
										MinCopyTrades:      minT,
										MinCopyROI:         minROI,
										MinCopyPnL:         minPnL,
										MinCopyWinRate:     minWin,
										MinClosedROI:       minClosedROI,
										MinSmartMoneyScore: minSmart,
									}
									res := buildStrategy(scores, params)
									if !validStrategy(res, minWallets, minTotalCopyTrades) {
										continue
									}
									key := walletKey(res.Wallets)
									if _, ok := seen[key]; ok {
										continue
									}
									seen[key] = struct{}{}
									out = append(out, res)
								}
							}
						}
					}
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].TotalCopyTrades != out[j].TotalCopyTrades {
			return out[i].TotalCopyTrades > out[j].TotalCopyTrades
		}
		return out[i].TotalCopyPnL > out[j].TotalCopyPnL
	})
	return out
}

func buildStrategy(scores []walletdiscover.WalletScore, params strategyParams) strategyResult {
	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if !walletdiscover.TierAllowed(s.Tier, params.MinTier) {
			continue
		}
		if s.FollowAction != "prompt" && s.FollowAction != "auto-small" {
			continue
		}
		if s.BotScore >= params.MaxBot {
			continue
		}
		if s.SmartMoneyScore < params.MinSmartMoneyScore {
			continue
		}
		if s.Stats.ClosedROI < params.MinClosedROI {
			continue
		}
		if s.Stats.CopyClosedTrades < params.MinCopyTrades {
			continue
		}
		if s.Stats.CopyROI < params.MinCopyROI || s.Stats.CopyPnL < params.MinCopyPnL {
			continue
		}
		if s.Stats.CopyWinRate < params.MinCopyWinRate {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if selected[i].Tier != selected[j].Tier {
			return walletdiscover.TierAllowed(selected[i].Tier, selected[j].Tier)
		}
		if selected[i].Stats.CopyPnL != selected[j].Stats.CopyPnL {
			return selected[i].Stats.CopyPnL > selected[j].Stats.CopyPnL
		}
		return selected[i].Address < selected[j].Address
	})
	return summarizeStrategy(params, selected)
}

func summarizeStrategy(params strategyParams, wallets []walletdiscover.WalletScore) strategyResult {
	res := strategyResult{Params: params, Wallets: wallets, WalletCount: len(wallets)}
	var rois []float64
	for _, s := range wallets {
		res.TotalCopyTrades += s.Stats.CopyClosedTrades
		res.TotalCopyWins += s.Stats.CopyWins
		res.TotalCopyPnL += s.Stats.CopyPnL
		res.TotalCopyCapital += s.Stats.CopyCapital
		res.TotalOpenCost += s.Stats.CopyOpenCost
		rois = append(rois, s.Stats.CopyROI)
	}
	if res.TotalCopyCapital > 0 {
		res.TotalCopyROI = res.TotalCopyPnL / res.TotalCopyCapital * 100
		res.OpenCostRatio = res.TotalOpenCost / res.TotalCopyCapital
	}
	if res.TotalCopyTrades > 0 {
		res.TotalCopyWinRate = float64(res.TotalCopyWins) / float64(res.TotalCopyTrades) * 100
	}
	sort.Float64s(rois)
	if len(rois) > 0 {
		res.WorstCopyROI = rois[0]
		if len(rois)%2 == 0 {
			res.MedianCopyROI = (rois[len(rois)/2-1] + rois[len(rois)/2]) / 2
		} else {
			res.MedianCopyROI = rois[len(rois)/2]
		}
	}
	res.Score = conservativeScore(res)
	return res
}

func validStrategy(res strategyResult, minWallets, minTotalCopyTrades int) bool {
	if res.WalletCount < minWallets || res.TotalCopyTrades < minTotalCopyTrades {
		return false
	}
	if res.TotalCopyCapital <= 0 || res.TotalCopyROI <= 0 || res.TotalCopyPnL <= 0 {
		return false
	}
	if res.WorstCopyROI <= 0 {
		return false
	}
	return true
}

func conservativeScore(res strategyResult) float64 {
	tradeConfidence := math.Log1p(float64(res.TotalCopyTrades))
	walletConfidence := math.Sqrt(float64(res.WalletCount))
	openPenalty := math.Min(res.OpenCostRatio*8, 35)
	return res.TotalCopyROI*tradeConfidence + res.MedianCopyROI*walletConfidence + res.TotalCopyWinRate - openPenalty
}

func buildWatchlist(scores, core []walletdiscover.WalletScore, params watchParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 {
		return nil
	}
	coreSet := map[string]struct{}{}
	for _, s := range core {
		coreSet[strings.ToLower(s.Address)] = struct{}{}
	}

	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if _, ok := coreSet[strings.ToLower(s.Address)]; ok {
			continue
		}
		if !walletdiscover.TierAllowed(s.Tier, "B") {
			continue
		}
		if s.FollowAction != "prompt" && s.FollowAction != "auto-small" {
			continue
		}
		if s.BotScore >= params.MaxBot || s.SmartMoneyScore < params.MinSmart {
			continue
		}
		if s.Stats.CopyClosedTrades < params.MinCopyTrades ||
			s.Stats.CopyROI < params.MinCopyROI ||
			s.Stats.CopyPnL < params.MinCopyPnL ||
			s.Stats.CopyWinRate < params.MinCopyWinRate ||
			s.Stats.LargeTrades < params.MinLargeTrades ||
			s.Stats.AvgTradeNotional < params.MinAvgNotional {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if activeWatchScore(selected[i]) != activeWatchScore(selected[j]) {
			return activeWatchScore(selected[i]) > activeWatchScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildSportslist(scores, excluded []walletdiscover.WalletScore, params sportsParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}

	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if _, ok := excludedSet[strings.ToLower(s.Address)]; ok {
			continue
		}
		if s.BotScore >= params.MaxBot ||
			s.SmartMoneyScore < params.MinSmart ||
			s.Edge < params.MinEdge ||
			s.Stats.TargetTrades < params.MinTargetTrades ||
			s.Stats.TargetTradeRatio < params.MinTargetRatio ||
			s.Stats.TargetCopyClosed < params.MinTargetCopyT ||
			s.Stats.TargetCopyROI < params.MinTargetCopyROI ||
			s.Stats.TargetCopyPnL < params.MinTargetCopyPnL ||
			s.Stats.TargetLargeTrades < params.MinTargetLarge {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if sportsScore(selected[i]) != sportsScore(selected[j]) {
			return sportsScore(selected[i]) > sportsScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildScoutlist(scores, excluded []walletdiscover.WalletScore, params scoutParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}

	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if _, ok := excludedSet[strings.ToLower(s.Address)]; ok {
			continue
		}
		if !hasRecentProfitLeaderboardSource(s.Sources) {
			continue
		}
		if s.BotScore >= params.MaxBot ||
			s.SmartMoneyScore < params.MinSmart ||
			s.Edge < params.MinEdge ||
			s.Stats.LargeTrades < params.MinLargeTrades ||
			s.Stats.AvgTradeNotional < params.MinAvgNotional {
			continue
		}
		if s.Stats.CopyClosedTrades >= 3 && (s.Stats.CopyROI < 0 || s.Stats.CopyPnL < 0) {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if scoutScore(selected[i]) != scoutScore(selected[j]) {
			return scoutScore(selected[i]) > scoutScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildTargetScoutlist(scores, excluded []walletdiscover.WalletScore, params scoutParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}

	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if _, ok := excludedSet[strings.ToLower(s.Address)]; ok {
			continue
		}
		if s.BotScore >= params.MaxBot ||
			s.SmartMoneyScore < params.MinSmart ||
			s.Edge < params.MinEdge ||
			s.Stats.TargetTrades < params.MinTargetTrades ||
			s.Stats.TargetLargeTrades < params.MinLargeTrades ||
			s.Stats.TargetTradeRatio < params.MinTargetRatio ||
			s.Stats.TargetCopyClosed < params.MinTargetCopyT ||
			s.Stats.TargetCopyROI < params.MinTargetCopyROI ||
			s.Stats.AvgTradeNotional < params.MinAvgNotional {
			continue
		}
		if s.Stats.CopyClosedTrades >= 3 && (s.Stats.CopyROI < 0 || s.Stats.CopyPnL < 0) {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if targetScoutScore(selected[i]) != targetScoutScore(selected[j]) {
			return targetScoutScore(selected[i]) > targetScoutScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildFlowScoutlist(scores, excluded []walletdiscover.WalletScore, params scoutParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}

	var selected []walletdiscover.WalletScore
	for _, s := range scores {
		if _, ok := excludedSet[strings.ToLower(s.Address)]; ok {
			continue
		}
		if s.Sources["recent_trade"] < params.MinRecentTrades {
			continue
		}
		if s.BotScore >= params.MaxBot ||
			s.SmartMoneyScore < params.MinSmart ||
			s.Edge < params.MinEdge ||
			s.Stats.LargeTrades < params.MinLargeTrades ||
			s.Stats.AvgTradeNotional < params.MinAvgNotional {
			continue
		}
		if s.Stats.CopyClosedTrades >= 3 && (s.Stats.CopyROI < 0 || s.Stats.CopyPnL < 0) {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if flowScore(selected[i]) != flowScore(selected[j]) {
			return flowScore(selected[i]) > flowScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildTapeHotlist(scores []walletdiscover.WalletScore, tapeInputs []tapeInput, excluded []walletdiscover.WalletScore, params tapeParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 || len(tapeInputs) == 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}
	scoreByAddress := map[string]walletdiscover.WalletScore{}
	for _, s := range scores {
		scoreByAddress[strings.ToLower(s.Address)] = s
	}

	var selected []walletdiscover.WalletScore
	for _, in := range tapeInputs {
		addr := strings.ToLower(in.Address)
		if addr == "" {
			continue
		}
		if _, ok := excludedSet[addr]; ok {
			continue
		}
		if edgeBlockReason(addr, params.EdgeBlocks) != "" {
			continue
		}
		s, scored := scoreByAddress[addr]
		if !scored {
			s = walletdiscover.WalletScore{
				Address:      addr,
				Tier:         cleanTapeTier(in.Tier),
				FollowAction: "watch",
				Sources:      map[string]int{"sports_tape": 1},
				Stats: walletdiscover.WalletStats{
					LargeTrades:       in.Buys,
					TargetTrades:      in.Buys,
					TargetLargeTrades: in.Buys,
					AvgTradeNotional:  in.MaxBuy,
				},
			}
			if in.Bot > 0 {
				s.BotScore = in.Bot
			}
		} else {
			if s.Sources == nil {
				s.Sources = map[string]int{}
			}
			s.Sources["sports_tape"]++
		}
		if s.Tier == "" {
			s.Tier = cleanTapeTier(in.Tier)
		}
		if in.Bot > 0 && s.BotScore == 0 {
			s.BotScore = in.Bot
		}
		if s.Stats.AvgTradeNotional <= 0 {
			s.Stats.AvgTradeNotional = in.MaxBuy
		}
		if !scored {
			continue
		}
		if strings.EqualFold(s.Tier, "BOT") {
			continue
		}
		if s.BotScore > 0 && s.BotScore >= params.MaxBot {
			continue
		}
		directWhale := in.MaxBuy >= params.MinDirectMaxBuy
		if directWhale && hasRiskFlag(s.RiskFlags, "bot_like_flow") {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludeRiskTags) {
			continue
		}
		scoredLarge := in.MaxBuy >= params.MinScoredMaxBuy &&
			s.SmartMoneyScore >= params.MinSmart &&
			walletdiscover.TierAllowed(s.Tier, "C")
		positiveTarget := (directWhale || scoredLarge) &&
			in.MaxBuy >= params.MinScoredMaxBuy &&
			s.SmartMoneyScore >= params.MinSmart &&
			s.Stats.TargetCopyClosed >= params.MinTargetCopyT &&
			s.Stats.TargetCopyROI >= params.MinTargetCopyROI
		if !positiveTarget {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if tapeHotScore(selected[i]) != tapeHotScore(selected[j]) {
			return tapeHotScore(selected[i]) > tapeHotScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildTapeProbation(scores []walletdiscover.WalletScore, tapeInputs []tapeInput, excluded []walletdiscover.WalletScore, params tapeParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 || len(tapeInputs) == 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}
	scoreByAddress := map[string]walletdiscover.WalletScore{}
	for _, s := range scores {
		scoreByAddress[strings.ToLower(s.Address)] = s
	}

	var selected []walletdiscover.WalletScore
	for _, in := range tapeInputs {
		addr := strings.ToLower(in.Address)
		if addr == "" || in.MaxBuy < params.MinScoredMaxBuy {
			continue
		}
		if _, ok := excludedSet[addr]; ok {
			continue
		}
		if edgeBlockReason(addr, params.EdgeBlocks) != "" {
			continue
		}
		s, scored := scoreByAddress[addr]
		if !scored {
			continue
		}
		if strings.EqualFold(s.Tier, "BOT") || (s.BotScore > 0 && s.BotScore >= params.MaxBot) {
			continue
		}
		if s.SmartMoneyScore < params.MinSmart || !walletdiscover.TierAllowed(s.Tier, "C") {
			continue
		}
		if s.Stats.TargetCopyClosed < params.MinTargetCopyT || s.Stats.TargetCopyROI < params.MinTargetCopyROI {
			continue
		}
		if hasRiskFlag(s.RiskFlags, "bot_like_flow") || hasRiskFlag(s.RiskFlags, "fixed_amount") || hasRiskFlag(s.RiskFlags, "negative_copy_sim") {
			continue
		}
		if s.Sources == nil {
			s.Sources = map[string]int{}
		}
		s.Sources["sports_tape"]++
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if tapeProbationScore(selected[i]) != tapeProbationScore(selected[j]) {
			return tapeProbationScore(selected[i]) > tapeProbationScore(selected[j])
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func buildTapeEdgeHot(scores []walletdiscover.WalletScore, tapeInputs []tapeInput, excluded []walletdiscover.WalletScore, params tapeEdgeHotParams) []walletdiscover.WalletScore {
	if params.MaxWallets <= 0 || len(params.EdgeProfiles) == 0 {
		return nil
	}
	excludedSet := map[string]struct{}{}
	for _, s := range excluded {
		excludedSet[strings.ToLower(s.Address)] = struct{}{}
	}
	scoreByAddress := map[string]walletdiscover.WalletScore{}
	for _, s := range scores {
		scoreByAddress[strings.ToLower(s.Address)] = s
	}

	inputByAddress := map[string]tapeInput{}
	for _, in := range tapeInputs {
		addr := strings.ToLower(in.Address)
		if addr == "" {
			continue
		}
		cur := inputByAddress[addr]
		if in.MaxBuy > cur.MaxBuy {
			inputByAddress[addr] = in
		}
	}
	for addr, profile := range params.EdgeProfiles {
		if _, ok := inputByAddress[addr]; ok || profile.MaxNotional <= 0 {
			continue
		}
		inputByAddress[addr] = tapeInput{
			Address:     addr,
			Buys:        1,
			BuyNotional: profile.MaxNotional,
			MaxBuy:      profile.MaxNotional,
		}
	}

	var selected []walletdiscover.WalletScore
	for _, in := range inputByAddress {
		addr := strings.ToLower(in.Address)
		if addr == "" || in.MaxBuy < params.MinMaxBuy {
			continue
		}
		if _, ok := excludedSet[addr]; ok {
			continue
		}
		if edgeBlockReason(addr, params.EdgeBlocks) != "" {
			continue
		}
		profile, ok := params.EdgeProfiles[addr]
		if !ok || !edgeHotProfilePasses(profile, params) {
			continue
		}
		s, scored := scoreByAddress[addr]
		if !scored && !profile.TapeAction {
			continue
		}
		if !scored {
			s = walletdiscover.WalletScore{
				Address:      addr,
				Tier:         cleanTapeTier(in.Tier),
				FollowAction: "watch",
				Sources:      map[string]int{"sports_tape": 1},
				Stats: walletdiscover.WalletStats{
					LargeTrades:       in.Buys,
					TargetTrades:      in.Buys,
					TargetLargeTrades: in.Buys,
					AvgTradeNotional:  in.MaxBuy,
				},
			}
			if in.Bot > 0 {
				s.BotScore = in.Bot
			}
		} else if !profile.TapeAction && s.Stats.TargetTradeRatio < 0.20 {
			continue
		} else {
			if s.Sources == nil {
				s.Sources = map[string]int{}
			}
			s.Sources["sports_tape"]++
		}
		if s.Tier == "" {
			s.Tier = cleanTapeTier(in.Tier)
		}
		if in.Bot > 0 && s.BotScore == 0 {
			s.BotScore = in.Bot
		}
		if s.Stats.AvgTradeNotional <= 0 {
			s.Stats.AvgTradeNotional = in.MaxBuy
		}
		if strings.EqualFold(s.Tier, "BOT") || (params.MaxBot > 0 && s.BotScore >= params.MaxBot) {
			continue
		}
		if hasExcludedRisk(s.RiskFlags, params.ExcludedRisks) {
			continue
		}
		selected = append(selected, s)
	}
	sort.Slice(selected, func(i, j int) bool {
		if tapeEdgeHotScore(selected[i], params.EdgeProfiles) != tapeEdgeHotScore(selected[j], params.EdgeProfiles) {
			return tapeEdgeHotScore(selected[i], params.EdgeProfiles) > tapeEdgeHotScore(selected[j], params.EdgeProfiles)
		}
		return selected[i].Address < selected[j].Address
	})
	if len(selected) > params.MaxWallets {
		selected = selected[:params.MaxWallets]
	}
	return selected
}

func edgeHotProfilePasses(profile edgeProfile, params tapeEdgeHotParams) bool {
	if params.MinSamples > 0 && profile.TotalSamples < params.MinSamples {
		return false
	}
	if profile.AvgPP < params.MinAvgPP || profile.WinRate < params.MinWinRate {
		return false
	}
	if params.Min5mAvgPP > 0 && (profile.Samples5m == 0 || profile.Avg5mPP < params.Min5mAvgPP) {
		return false
	}
	if profile.Samples15m == 0 || profile.Avg15mPP < params.Min15mAvgPP {
		return false
	}
	if profile.Samples1h > 0 && profile.Avg1hPP <= params.Max1hNegPP {
		return false
	}
	return true
}

func edgeBlockReason(addr string, blocks map[string]edgeBlock) string {
	if len(blocks) == 0 {
		return ""
	}
	block := blocks[strings.ToLower(strings.TrimSpace(addr))]
	return block.Reason
}

func activeWatchScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	openPenalty := 0.0
	if st.CopyCapital > 0 {
		openPenalty = math.Min(st.CopyOpenCost/st.CopyCapital*10, 25)
	}
	return s.SmartMoneyScore*0.8 +
		s.Edge*0.7 -
		s.BotScore*1.3 +
		st.CopyROI*math.Log1p(float64(st.CopyClosedTrades)) +
		st.CopyWinRate*0.4 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*5 +
		targetBonus(st) +
		sourceBonus(s.Sources) -
		openPenalty
}

func sportsScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	openPenalty := 0.0
	if st.TargetCopyCapital > 0 {
		openPenalty = math.Min(st.TargetCopyOpenCost/st.TargetCopyCapital*8, 25)
	}
	return s.SmartMoneyScore*0.55 +
		s.Edge*0.75 -
		s.BotScore*1.4 +
		st.TargetCopyROI*math.Log1p(float64(st.TargetCopyClosed)) +
		st.TargetCopyWinRate*0.35 +
		math.Log1p(math.Max(st.TargetCopyPnL, 0))*8 +
		targetBonus(st) -
		openPenalty
}

func scoutScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return s.SmartMoneyScore*0.6 +
		s.Edge*0.7 -
		s.BotScore*1.8 +
		float64(st.LargeTrades)*0.08 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*9 +
		targetBonus(st) +
		sourceBonus(s.Sources)
}

func targetScoutScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return s.SmartMoneyScore*0.7 +
		s.Edge*0.8 -
		s.BotScore*1.8 +
		float64(st.TargetLargeTrades)*0.12 +
		st.TargetTradeRatio*55 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*7 +
		sourceBonus(s.Sources)
}

func flowScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return s.SmartMoneyScore*0.7 +
		s.Edge*0.8 -
		s.BotScore*1.7 +
		float64(s.Sources["recent_trade"])*18 +
		float64(st.LargeTrades)*0.12 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*10 +
		targetBonus(st) +
		sourceBonus(s.Sources)
}

func tapeHotScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return math.Log1p(math.Max(st.AvgTradeNotional, 0))*22 +
		s.SmartMoneyScore*0.4 +
		s.Edge*0.4 -
		s.BotScore*1.5 +
		st.TargetCopyROI*math.Log1p(float64(st.TargetCopyClosed))*0.8 +
		sourceBonus(s.Sources)
}

func tapeProbationScore(s walletdiscover.WalletScore) float64 {
	score := tapeHotScore(s)
	if hasRiskFlag(s.RiskFlags, "fixed_price") {
		score -= 20
	}
	if hasRiskFlag(s.RiskFlags, "opposite_side_same_market") {
		score -= 15
	}
	if hasRiskFlag(s.RiskFlags, "open_copy_exposure") {
		score -= 10
	}
	return score
}

func tapeEdgeHotScore(s walletdiscover.WalletScore, profiles map[string]edgeProfile) float64 {
	st := s.Stats
	profile := profiles[strings.ToLower(s.Address)]
	return profile.AvgPP*8 +
		profile.WinRate*0.5 +
		math.Log1p(math.Max(st.AvgTradeNotional, 0))*16 +
		s.SmartMoneyScore*0.25 -
		s.BotScore*1.2 +
		st.TargetCopyROI*math.Log1p(float64(st.TargetCopyClosed))*0.35 +
		sourceBonus(s.Sources)
}

func targetBonus(st walletdiscover.WalletStats) float64 {
	if st.TargetTrades == 0 {
		return 0
	}
	return math.Min(float64(st.TargetLargeTrades)*2+st.TargetTradeRatio*35, 45)
}

func sourceBonus(sources map[string]int) float64 {
	var bonus float64
	for src, n := range sources {
		if n <= 0 {
			continue
		}
		switch {
		case strings.Contains(src, "leaderboard_profit_7d"):
			bonus += 18
		case strings.Contains(src, "leaderboard_profit_30d"):
			bonus += 12
		case strings.Contains(src, "leaderboard_volume_7d"):
			bonus += 10
		case strings.Contains(src, "leaderboard_volume_30d"):
			bonus += 6
		case strings.Contains(src, "recent_trade"):
			bonus += 8
		case strings.Contains(src, "sports_tape"):
			bonus += 8
		case strings.Contains(src, "holder"):
			bonus += 3
		case strings.Contains(src, "existing"):
			bonus += 2
		}
	}
	return bonus
}

func hasRecentProfitLeaderboardSource(sources map[string]int) bool {
	for _, src := range []string{"leaderboard_profit_7d", "leaderboard_profit_30d"} {
		if sources[src] > 0 {
			return true
		}
	}
	return false
}

func hasExcludedRisk(flags []string, excluded map[string]struct{}) bool {
	for _, flag := range flags {
		if _, ok := excluded[flag]; ok {
			return true
		}
	}
	return false
}

func hasRiskFlag(flags []string, target string) bool {
	for _, flag := range flags {
		if flag == target {
			return true
		}
	}
	return false
}

func mergeWalletGroups(groups ...[]walletdiscover.WalletScore) []walletdiscover.WalletScore {
	var capHint int
	for _, group := range groups {
		capHint += len(group)
	}
	out := make([]walletdiscover.WalletScore, 0, capHint)
	seen := map[string]struct{}{}
	for _, group := range groups {
		for _, s := range group {
			addr := strings.ToLower(s.Address)
			if _, ok := seen[addr]; ok {
				continue
			}
			seen[addr] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

func filterEdgeBlockedWallets(wallets []walletdiscover.WalletScore, blocks map[string]edgeBlock) []walletdiscover.WalletScore {
	if len(wallets) == 0 || len(blocks) == 0 {
		return wallets
	}
	out := make([]walletdiscover.WalletScore, 0, len(wallets))
	for _, s := range wallets {
		if edgeBlockReason(s.Address, blocks) != "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func edgeBlockedWallets(wallets []walletdiscover.WalletScore, blocks map[string]edgeBlock) []walletdiscover.WalletScore {
	if len(wallets) == 0 || len(blocks) == 0 {
		return nil
	}
	var out []walletdiscover.WalletScore
	for _, s := range wallets {
		if edgeBlockReason(s.Address, blocks) == "" {
			continue
		}
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		ri := edgeBlockReason(out[i].Address, blocks)
		rj := edgeBlockReason(out[j].Address, blocks)
		if ri != rj {
			return ri < rj
		}
		return out[i].Address < out[j].Address
	})
	return out
}

func walletKey(wallets []walletdiscover.WalletScore) string {
	addrs := make([]string, 0, len(wallets))
	for _, s := range wallets {
		addrs = append(addrs, s.Address)
	}
	sort.Strings(addrs)
	return strings.Join(addrs, ",")
}

func writeCoreWallets(path string, wallets []walletdiscover.WalletScore) error {
	return writeWalletsFile(path, wallets, "core")
}

func writeWalletsFile(path string, wallets []walletdiscover.WalletScore, label string) error {
	var lines []string
	for _, s := range wallets {
		lines = append(lines, walletLine(s, label))
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func writePushWalletsFile(path string, core, watch, sports, scout, target, flow, tape []walletdiscover.WalletScore) error {
	var lines []string
	seen := map[string]struct{}{}
	for _, group := range []struct {
		label   string
		wallets []walletdiscover.WalletScore
	}{
		{label: "core", wallets: core},
		{label: "watch", wallets: watch},
		{label: "sports", wallets: sports},
		{label: "scout", wallets: scout},
		{label: "target", wallets: target},
		{label: "flow", wallets: flow},
		{label: "tape", wallets: tape},
	} {
		for _, s := range group.wallets {
			addr := strings.ToLower(s.Address)
			if _, ok := seen[addr]; ok {
				continue
			}
			seen[addr] = struct{}{}
			lines = append(lines, walletLine(s, group.label))
		}
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func writeDirectTapeObserveFile(path string, scores, push []walletdiscover.WalletScore, inputs []tapeInput, minMaxBuy, minBuyNotional, maxBot float64) error {
	rows := directTapeObserveRows(scores, push, inputs, minMaxBuy, minBuyNotional, maxBot)
	lines := make([]string, 0, len(rows))
	lines = append(lines, rows...)
	sort.Strings(lines)
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func writeTapeProbationFile(path string, wallets []walletdiscover.WalletScore) error {
	var lines []string
	for _, s := range wallets {
		lines = append(lines, fmt.Sprintf("%s # list=tape_probation tier=%s smart=%.2f bot=%.2f targetCopyROI=%.1f%% targetCopyT=%d targetLarge=%d avgNotional=$%.0f risks=%s",
			strings.ToLower(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed,
			s.Stats.TargetLargeTrades, s.Stats.AvgTradeNotional, riskSummary(s.RiskFlags)))
	}
	sort.Strings(lines)
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func writeTapeEdgeHotFile(path string, wallets []walletdiscover.WalletScore, profiles map[string]edgeProfile, inputs []tapeInput) error {
	inputByAddress := map[string]tapeInput{}
	for _, in := range inputs {
		addr := strings.ToLower(in.Address)
		if addr == "" {
			continue
		}
		cur := inputByAddress[addr]
		if in.MaxBuy > cur.MaxBuy {
			inputByAddress[addr] = in
		}
	}
	var lines []string
	for _, s := range wallets {
		addr := strings.ToLower(s.Address)
		profile := profiles[addr]
		in := inputByAddress[addr]
		if in.MaxBuy <= 0 && profile.MaxNotional > 0 {
			in = tapeInput{
				Address:     addr,
				Buys:        1,
				BuyNotional: profile.MaxNotional,
				MaxBuy:      profile.MaxNotional,
			}
		}
		lines = append(lines, fmt.Sprintf("%s # list=tape_edgehot tier=%s smart=%.2f bot=%.2f maxBuy=$%.0f buyNotional=$%.0f buys=%d edgeN=%d edgeWin=%.1f%% edgeAvgPP=%+.2f edge5mPP=%+.2f edge15mPP=%+.2f reason=%s",
			addr, s.Tier, s.SmartMoneyScore, s.BotScore, in.MaxBuy, in.BuyNotional, in.Buys,
			profile.TotalSamples, profile.WinRate, profile.AvgPP, profile.Avg5mPP, profile.Avg15mPP, profile.Reason))
	}
	sort.Strings(lines)
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func directTapeObserveRows(scores, push []walletdiscover.WalletScore, inputs []tapeInput, minMaxBuy, minBuyNotional, maxBot float64) []string {
	if len(inputs) == 0 || (minMaxBuy <= 0 && minBuyNotional <= 0) {
		return nil
	}
	scoreByAddress := map[string]walletdiscover.WalletScore{}
	for _, s := range scores {
		scoreByAddress[strings.ToLower(s.Address)] = s
	}
	pushSet := map[string]struct{}{}
	for _, s := range push {
		pushSet[strings.ToLower(s.Address)] = struct{}{}
	}
	seen := map[string]struct{}{}
	var rows []string
	for _, in := range inputs {
		addr := strings.ToLower(in.Address)
		if addr == "" {
			continue
		}
		if _, ok := pushSet[addr]; ok {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		s, scored := scoreByAddress[addr]
		direct := minMaxBuy > 0 && in.MaxBuy >= minMaxBuy
		burst := !direct && minBuyNotional > 0 && in.BuyNotional >= minBuyNotional && in.Buys >= 2
		if !direct && !burst {
			continue
		}
		status := "observe-direct"
		tier := cleanTapeTier(in.Tier)
		bot := in.Bot
		copyROI := 0.0
		copyT := 0
		risks := "-"
		if scored {
			status = tapeGateStatus(s)
			tier = s.Tier
			bot = s.BotScore
			copyROI = s.Stats.CopyROI
			copyT = s.Stats.CopyClosedTrades
			risks = riskSummary(s.RiskFlags)
		}
		if burst {
			if !scored || !tierAtLeast(tier, "B") {
				continue
			}
			if maxBot > 0 && bot > maxBot {
				continue
			}
			if status == "watch" {
				status = "watch-burst"
			}
		}
		rows = append(rows, fmt.Sprintf("%s # list=tape_observe status=%s tier=%s bot=%.1f maxBuy=$%.0f buyNotional=$%.0f buys=%d copyROI=%.1f%% copyT=%d risks=%s",
			addr, status, tier, bot, in.MaxBuy, in.BuyNotional, in.Buys, copyROI, copyT, risks))
	}
	return rows
}

func walletLine(s walletdiscover.WalletScore, label string) string {
	return fmt.Sprintf("%s # list=%s tier=%s smart=%.2f bot=%.2f copyROI=%.1f%% copyPnL=$%+.2f copyT=%d copyWin=%.1f%% avgNotional=$%.0f sources=%s",
		s.Address, label, s.Tier, s.SmartMoneyScore, s.BotScore,
		s.Stats.CopyROI, s.Stats.CopyPnL, s.Stats.CopyClosedTrades, s.Stats.CopyWinRate,
		s.Stats.AvgTradeNotional, sourceSummary(s.Sources))
}

func sourceSummary(sources map[string]int) string {
	if len(sources) == 0 {
		return "-"
	}
	var keys []string
	for k, n := range sources {
		if n > 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if len(keys) > 4 {
		keys = keys[:4]
	}
	return strings.Join(keys, ",")
}

func writeReport(path, scoresPath, excludePath string, excludedCount int, scores []walletdiscover.WalletScore, best strategyResult, strategies []strategyResult, watch, sports, scout, target, flow, tape, tapeProbation, tapeEdgeHot, push, blockedPush []walletdiscover.WalletScore, tapeInputs []tapeInput, edgeProfiles map[string]edgeProfile, edgeBlocks map[string]edgeBlock, directTapeMin float64, scoutPushEnabled bool, topN int) error {
	if topN <= 0 || topN > len(strategies) {
		topN = len(strategies)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Strategy Lab Report\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Scores: `%s`\n", scoresPath)
	if excludePath != "" {
		fmt.Fprintf(&b, "- Exclude wallets: `%s` (%d)\n", excludePath, excludedCount)
	}
	fmt.Fprintf(&b, "- Valid strategies found: %d\n", len(strategies))
	fmt.Fprintf(&b, "- Candidate layers: %d core + %d watch + %d sports + %d scout + %d target + %d flow + %d tape\n", best.WalletCount, len(watch), len(sports), len(scout), len(target), len(flow), len(tape))
	fmt.Fprintf(&b, "- Push wallets after live-edge blocks: %d total\n", len(push))
	fmt.Fprintf(&b, "- Live-edge blocked push wallets: %d\n", len(blockedPush))
	fmt.Fprintf(&b, "- Leaderboard scout push enabled: %t\n", scoutPushEnabled)
	fmt.Fprintf(&b, "- Tape probation wallets: %d observation-only\n", len(tapeProbation))
	fmt.Fprintf(&b, "- Tape edge-hot wallets: %d observation-only\n\n", len(tapeEdgeHot))

	fmt.Fprintf(&b, "## Selected Core Strategy\n\n")
	writeStrategySummary(&b, best)
	fmt.Fprintf(&b, "\n## Core Wallets\n\n")
	fmt.Fprintf(&b, "| Wallet | Tier | Bot | TargetT | Target%% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |\n")
	fmt.Fprintf(&b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, s := range best.Wallets {
		fmt.Fprintf(&b, "| `%s` | %s | %.2f | %d | %.0f%% | %.1f%% | $%+.2f | %d | %.1f%% | %.1f%% |\n",
			shortAddr(s.Address), s.Tier, s.BotScore, s.Stats.TargetTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.CopyROI, s.Stats.CopyPnL,
			s.Stats.CopyClosedTrades, s.Stats.CopyWinRate, s.Stats.ClosedROI)
	}

	writeWatchlistReport(&b, watch)
	writeSportsReport(&b, sports)
	writeScoutReport(&b, scout)
	writeTargetScoutReport(&b, target)
	writeFlowReport(&b, flow)
	writeTapeHotReport(&b, tape)
	writeTapeProbationReport(&b, tapeProbation)
	writeTapeEdgeHotReport(&b, tapeEdgeHot, edgeProfiles, tapeInputs)
	writeEdgeBlockedPushReport(&b, blockedPush, edgeBlocks)
	writeDirectTapeGateReview(&b, scores, push, tapeInputs, edgeBlocks, directTapeMin)
	writeSportsTapeReview(&b, scores, push, edgeBlocks)

	fmt.Fprintf(&b, "\n## Top Strategies\n\n")
	fmt.Fprintf(&b, "| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |\n")
	fmt.Fprintf(&b, "|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for i := 0; i < topN; i++ {
		s := strategies[i]
		fmt.Fprintf(&b, "| %d | %d | %d | %.1f%% | $%+.2f | %.1f%% | %.1f%% | %.1f%% | %.2fx | %s |\n",
			i+1, s.WalletCount, s.TotalCopyTrades, s.TotalCopyROI, s.TotalCopyPnL,
			s.TotalCopyWinRate, s.MedianCopyROI, s.WorstCopyROI, s.OpenCostRatio, formatParams(s.Params))
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeSportsTapeReview(b *strings.Builder, scores, push []walletdiscover.WalletScore, edgeBlocks map[string]edgeBlock) {
	var tape []walletdiscover.WalletScore
	for _, s := range scores {
		if s.Sources["sports_tape"] > 0 {
			tape = append(tape, s)
		}
	}
	if len(tape) == 0 {
		return
	}
	pushSet := map[string]struct{}{}
	for _, s := range push {
		pushSet[strings.ToLower(s.Address)] = struct{}{}
	}
	sort.Slice(tape, func(i, j int) bool {
		if tapeCandidateScore(tape[i]) != tapeCandidateScore(tape[j]) {
			return tapeCandidateScore(tape[i]) > tapeCandidateScore(tape[j])
		}
		return tape[i].Address < tape[j].Address
	})
	if len(tape) > 15 {
		tape = tape[:15]
	}
	fmt.Fprintf(b, "\n## Sports Tape Candidate Review\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(tape))
	fmt.Fprintf(b, "- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters\n\n")
	fmt.Fprintf(b, "| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target%% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range tape {
		status := s.FollowAction
		blockReason := edgeBlockReason(s.Address, edgeBlocks)
		if blockReason != "" {
			status = "blocked-edge"
		}
		if _, ok := pushSet[strings.ToLower(s.Address)]; ok {
			status = "pushed"
		}
		fmt.Fprintf(b, "| `%s` | %s | %s | %s | %.1f | %.1f | %.1f | %d | %d | %.0f%% | %.1f%% | %d | %.1f%% | %d | %s | %s |\n",
			shortAddr(s.Address), status, dash(blockReason), s.Tier, s.SmartMoneyScore, s.BotScore, tapeCandidateScore(s),
			s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed, s.Stats.CopyROI, s.Stats.CopyClosedTrades,
			sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
}

func writeDirectTapeGateReview(b *strings.Builder, scores, push []walletdiscover.WalletScore, inputs []tapeInput, edgeBlocks map[string]edgeBlock, minMaxBuy float64) {
	if len(inputs) == 0 || minMaxBuy <= 0 {
		return
	}
	scoreByAddress := map[string]walletdiscover.WalletScore{}
	for _, s := range scores {
		scoreByAddress[strings.ToLower(s.Address)] = s
	}
	pushSet := map[string]struct{}{}
	for _, s := range push {
		pushSet[strings.ToLower(s.Address)] = struct{}{}
	}
	var rows []tapeInput
	for _, in := range inputs {
		if in.MaxBuy >= minMaxBuy {
			rows = append(rows, in)
		}
	}
	if len(rows) == 0 {
		return
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].MaxBuy != rows[j].MaxBuy {
			return rows[i].MaxBuy > rows[j].MaxBuy
		}
		return rows[i].Address < rows[j].Address
	})
	if len(rows) > 12 {
		rows = rows[:12]
	}

	fmt.Fprintf(b, "\n## Direct Sports Tape Gate Review\n\n")
	fmt.Fprintf(b, "- Rule: %.0f+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.\n\n", minMaxBuy)
	fmt.Fprintf(b, "| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, in := range rows {
		addr := strings.ToLower(in.Address)
		s, scored := scoreByAddress[addr]
		status := "direct-tape"
		tier := cleanTapeTier(in.Tier)
		bot := in.Bot
		copyROI := 0.0
		copyT := 0
		risks := "-"
		if scored {
			tier = s.Tier
			bot = s.BotScore
			copyROI = s.Stats.CopyROI
			copyT = s.Stats.CopyClosedTrades
			risks = riskSummary(s.RiskFlags)
			status = tapeGateStatus(s)
		}
		blockReason := edgeBlockReason(addr, edgeBlocks)
		if blockReason != "" {
			status = "blocked-edge"
		}
		if _, ok := pushSet[addr]; ok {
			status = "pushed"
		}
		fmt.Fprintf(b, "| `%s` | %s | %s | %s | %.1f | $%.0f | $%.0f | %d | %.1f%% | %d | %s |\n",
			shortAddr(addr), status, dash(blockReason), tier, bot, in.MaxBuy, in.BuyNotional, in.Buys, copyROI, copyT, risks)
	}
}

func tapeGateStatus(s walletdiscover.WalletScore) string {
	if strings.EqualFold(s.Tier, "BOT") || s.BotScore >= 45 {
		return "reject-bot"
	}
	if hasRiskFlag(s.RiskFlags, "bot_like_flow") || hasRiskFlag(s.RiskFlags, "fixed_price") || hasRiskFlag(s.RiskFlags, "fixed_amount") {
		return "reject-flow"
	}
	if hasRiskFlag(s.RiskFlags, "negative_copy_sim") {
		return "reject-copy"
	}
	return "watch"
}

func tapeCandidateScore(s walletdiscover.WalletScore) float64 {
	st := s.Stats
	return s.SmartMoneyScore*0.45 +
		s.Edge*0.4 -
		s.BotScore*1.8 +
		float64(st.TargetLargeTrades)*0.08 +
		st.TargetTradeRatio*25 +
		st.TargetCopyROI*math.Log1p(float64(st.TargetCopyClosed)) +
		st.CopyROI*math.Log1p(float64(st.CopyClosedTrades))*0.4
}

func writeFlowReport(b *strings.Builder, flow []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Recent Trade Flow Scout\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(flow))
	fmt.Fprintf(b, "- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target%% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range flow {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %.0f%% | %.1f%% | %d | %d | $%.0f | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, flowScore(s), s.Sources["recent_trade"],
			s.Stats.TargetTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, s.Stats.LargeTrades, s.Stats.AvgTradeNotional,
			sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(flow) == 0 {
		fmt.Fprintf(b, "\nNo recent-trade wallets passed the flow filters.\n")
	}
}

func writeTapeHotReport(b *strings.Builder, tape []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Sports Tape Hotlist\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(tape))
	fmt.Fprintf(b, "- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range tape {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %.1f%% | %d | %.1f%% | %d | $%.0f | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, tapeHotScore(s),
			s.Stats.TargetTrades, s.Stats.TargetLargeTrades,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed,
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, s.Stats.AvgTradeNotional,
			sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(tape) == 0 {
		fmt.Fprintf(b, "\nNo recent sports-tape wallets passed the hotlist filters.\n")
	}
}

func writeTapeProbationReport(b *strings.Builder, tape []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Sports Tape Probation\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(tape))
	fmt.Fprintf(b, "- Rule: scored sports-tape wallets with positive target-copy and sub-45 bot score, but soft flow risks; observation-only until edge windows prove out\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | ProbationScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, s := range tape {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %.1f%% | %d | %.1f%% | %d | $%.0f | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, tapeProbationScore(s),
			s.Stats.TargetTrades, s.Stats.TargetLargeTrades,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed,
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, s.Stats.AvgTradeNotional,
			riskSummary(s.RiskFlags))
	}
	if len(tape) == 0 {
		fmt.Fprintf(b, "\nNo sports-tape wallets qualified for probation edge observation.\n")
	}
}

func writeTapeEdgeHotReport(b *strings.Builder, wallets []walletdiscover.WalletScore, profiles map[string]edgeProfile, inputs []tapeInput) {
	inputByAddress := map[string]tapeInput{}
	for _, in := range inputs {
		addr := strings.ToLower(in.Address)
		if addr == "" {
			continue
		}
		cur := inputByAddress[addr]
		if in.MaxBuy > cur.MaxBuy {
			inputByAddress[addr] = in
		}
	}
	fmt.Fprintf(b, "\n## Sports Tape Edge-Hot\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(wallets))
	fmt.Fprintf(b, "- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, s := range wallets {
		addr := strings.ToLower(s.Address)
		profile := profiles[addr]
		in := inputByAddress[addr]
		if in.MaxBuy <= 0 && profile.MaxNotional > 0 {
			in = tapeInput{
				Address:     addr,
				Buys:        1,
				BuyNotional: profile.MaxNotional,
				MaxBuy:      profile.MaxNotional,
			}
		}
		fmt.Fprintf(b, "| `%s` | %s | %.1f | $%.0f | $%.0f | %d | %.1f%% | %+.2f | %+.2f | %+.2f | %+.2f | %.1f | %s |\n",
			shortAddr(addr), s.Tier, s.BotScore, in.MaxBuy, in.BuyNotional, profile.TotalSamples, profile.WinRate,
			profile.AvgPP, profile.Avg5mPP, profile.Avg15mPP, profile.Avg1hPP, tapeEdgeHotScore(s, profiles), riskSummary(s.RiskFlags))
	}
	if len(wallets) == 0 {
		fmt.Fprintf(b, "\nNo sports-tape wallets currently pass the measured edge-hot filters.\n")
	}
}

func writeEdgeBlockedPushReport(b *strings.Builder, wallets []walletdiscover.WalletScore, blocks map[string]edgeBlock) {
	fmt.Fprintf(b, "\n## Live Edge Blocked Push Wallets\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(wallets))
	fmt.Fprintf(b, "- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---|---:|---:|---|---|\n")
	for _, s := range wallets {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %s | %.1f%% | %d | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, dash(edgeBlockReason(s.Address, blocks)),
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(wallets) == 0 {
		fmt.Fprintf(b, "\nNo selected push wallets are blocked by live edge.\n")
	}
}

func writeSportsReport(b *strings.Builder, sports []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Sports Positive Copy\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(sports))
	fmt.Fprintf(b, "- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target%% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range sports {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %.0f%% | %.1f%% | $%+.2f | %d | %.1f%% | %.1f%% | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, sportsScore(s),
			s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyPnL, s.Stats.TargetCopyClosed, s.Stats.TargetCopyWinRate,
			s.Stats.CopyROI, sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(sports) == 0 {
		fmt.Fprintf(b, "\nNo target-category positive-copy wallets passed the sports filters.\n")
	}
}

func writeTargetScoutReport(b *strings.Builder, target []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Target Category Scout\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(target))
	fmt.Fprintf(b, "- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target%% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range target {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %d | %.0f%% | %.1f%% | %d | %.1f%% | %d | $%.0f | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, targetScoutScore(s),
			s.Stats.TargetTrades, s.Stats.TargetLargeTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.TargetCopyROI, s.Stats.TargetCopyClosed, s.Stats.CopyROI, s.Stats.CopyClosedTrades, s.Stats.AvgTradeNotional,
			sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(target) == 0 {
		fmt.Fprintf(b, "\nNo target-category wallets passed the scout filters.\n")
	}
}

func writeScoutReport(b *strings.Builder, scout []walletdiscover.WalletScore) {
	fmt.Fprintf(b, "\n## Leaderboard Scout\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(scout))
	fmt.Fprintf(b, "- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default\n\n")
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target%% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, s := range scout {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %.0f%% | %.1f%% | %d | %d | $%.0f | %s | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, scoutScore(s),
			s.Stats.TargetTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.CopyROI, s.Stats.CopyClosedTrades, s.Stats.LargeTrades, s.Stats.AvgTradeNotional,
			sourceSummary(s.Sources), riskSummary(s.RiskFlags))
	}
	if len(scout) == 0 {
		fmt.Fprintf(b, "\nNo leaderboard-only wallets passed the scout filters.\n")
	}
}

func writeWatchlistReport(b *strings.Builder, watch []walletdiscover.WalletScore) {
	sum := summarizeStrategy(strategyParams{}, watch)
	fmt.Fprintf(b, "\n## Active Watchlist\n\n")
	fmt.Fprintf(b, "- Wallets: %d\n", len(watch))
	fmt.Fprintf(b, "- Aggregate closed copy trades: %d\n", sum.TotalCopyTrades)
	fmt.Fprintf(b, "- Aggregate copy ROI: %.1f%%\n", sum.TotalCopyROI)
	fmt.Fprintf(b, "- Aggregate copy PnL: $%+.2f\n", sum.TotalCopyPnL)
	fmt.Fprintf(b, "- Aggregate copy win rate: %.1f%%\n", sum.TotalCopyWinRate)
	fmt.Fprintf(b, "- Worst included CopyROI: %.1f%%\n\n", sum.WorstCopyROI)
	fmt.Fprintf(b, "| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target%% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |\n")
	fmt.Fprintf(b, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, s := range watch {
		fmt.Fprintf(b, "| `%s` | %s | %.1f | %.1f | %.1f | %d | %.0f%% | %.1f%% | $%+.2f | %d | %.1f%% | $%.0f | %s |\n",
			shortAddr(s.Address), s.Tier, s.SmartMoneyScore, s.BotScore, activeWatchScore(s),
			s.Stats.TargetTrades, s.Stats.TargetTradeRatio*100,
			s.Stats.CopyROI, s.Stats.CopyPnL, s.Stats.CopyClosedTrades, s.Stats.CopyWinRate,
			s.Stats.AvgTradeNotional, sourceSummary(s.Sources))
	}
	if len(watch) == 0 {
		fmt.Fprintf(b, "\nNo non-core wallets passed the active watchlist filters.\n")
	}
}

func riskSummary(flags []string) string {
	if len(flags) == 0 {
		return "-"
	}
	out := append([]string{}, flags...)
	sort.Strings(out)
	if len(out) > 4 {
		out = out[:4]
	}
	return strings.Join(out, ",")
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func writeStrategySummary(b *strings.Builder, s strategyResult) {
	fmt.Fprintf(b, "- Wallets: %d\n", s.WalletCount)
	fmt.Fprintf(b, "- Aggregate closed copy trades: %d\n", s.TotalCopyTrades)
	fmt.Fprintf(b, "- Aggregate copy ROI: %.1f%%\n", s.TotalCopyROI)
	fmt.Fprintf(b, "- Aggregate copy PnL: $%+.2f\n", s.TotalCopyPnL)
	fmt.Fprintf(b, "- Aggregate copy win rate: %.1f%%\n", s.TotalCopyWinRate)
	fmt.Fprintf(b, "- Median wallet CopyROI: %.1f%%\n", s.MedianCopyROI)
	fmt.Fprintf(b, "- Worst included wallet CopyROI: %.1f%%\n", s.WorstCopyROI)
	fmt.Fprintf(b, "- Open copy cost / closed copy capital: %.2fx\n", s.OpenCostRatio)
	fmt.Fprintf(b, "- Params: %s\n", formatParams(s.Params))
}

func formatParams(p strategyParams) string {
	return fmt.Sprintf("tier>=%s bot<%.0f copyT>=%d copyROI>=%.0f copyPnL>=%.0f copyWin>=%.0f closedROI>=%.0f smart>=%.0f",
		p.MinTier, p.MaxBot, p.MinCopyTrades, p.MinCopyROI, p.MinCopyPnL, p.MinCopyWinRate, p.MinClosedROI, p.MinSmartMoneyScore)
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}
