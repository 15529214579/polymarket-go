package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	nethttp "net/http"
	neturl "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/15529214579/polymarket-go/internal/btc"
	"github.com/15529214579/polymarket-go/internal/config"
	"github.com/15529214579/polymarket-go/internal/elon"
	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/injury"
	"github.com/15529214579/polymarket-go/internal/iterate"
	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/notify"
	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/15529214579/polymarket-go/internal/risk"
	"github.com/15529214579/polymarket-go/internal/sanitize"
	"github.com/15529214579/polymarket-go/internal/strategy"
	"github.com/15529214579/polymarket-go/internal/tickrec"
	"github.com/15529214579/polymarket-go/internal/whale"
)

func main() {
	mode := flag.String("mode", "run", "run | discover | feed | sample | detect | prompt-test | daily-report | daily-iterate | arb-scan(disabled)")
	iterateWindow := flag.Int("iterate_window", 7, "daily-iterate: rolling window days for analysis")
	maxMarkets := flag.Int("markets", 20, "top-N sports markets (LoL + NBA daily/playoffs + EPL daily) by vol24h to subscribe")
	windowSec := flag.Int("window", 60, "sampler window in seconds")
	slippageBp := flag.Float64("slippage_bp", 0, "paper fill slippage in bp applied against you")
	largeFillUSD := flag.Float64("large_fill_usd", 3.0, "DM notifier threshold on |realized pnl|")
	envFile := flag.String("env_file", ".env.local", "dotenv file to load before reading env")
	signalMode := flag.String("signal_mode", "auto", "auto (paper-submit + DM) | prompt (DM only) | whale (follow whale buys, auto-close on sells)")
	exitMode := flag.String("exit_mode", "hold", "hold (settlement only) | auto (SPEC §2 reversal/drawdown/stop/timeout) | ladder (Phase 7.b TP1/TP2/SL/timeout)")
	journalDir := flag.String("journal_dir", "db/journal", "trade-journal directory (one JSONL per SGT day)")
	tickPathDir := flag.String("tickpath_dir", "db/tickpath", "Phase 7.e tick-path persistence dir (one JSONL per posID; empty disables)")
	reportDay := flag.String("report_day", "", "daily-report mode: SGT day YYYY-MM-DD (default: yesterday SGT)")
	reportPush := flag.Bool("report_push", false, "daily-report mode: also push summary via Telegram alert bot")
	// Phase 7.a entry-price band filter: only emit SignalPrompt when sig.Mid is
	// inside [min, max]. Default 0.15–0.70 matches python-db winner distribution
	// (see reports/python_autopsy.md §4–5).
	minEntry := flag.Float64("min_entry_price", 0.15, "signals with mid < this are filtered out (reports/python_autopsy.md §2.1)")
	maxEntry := flag.Float64("max_entry_price", 0.60, "signals with mid > this are filtered out")
	// Phase 7.b ladder TP / SL / timeout + fee modeling. Defaults are SPEC §2.4.
	feeBp := flag.Float64("fee_bp", 0, "per-side fee in basis points of notional; default 0 matches CLOB V1 reality (update after V2 cutover)")
	ladderTP1Pct := flag.Float64("ladder_tp1_pct", 9.99, "ladder TP1 trigger: 9.99 = effectively disabled (ride to settlement/timeout)")
	ladderTP1Frac := flag.Float64("ladder_tp1_frac", 0.50, "fraction of initial units to close on TP1")
	ladderTP2Pct := flag.Float64("ladder_tp2_pct", 9.99, "ladder TP2 trigger: 9.99 = effectively disabled")
	ladderTP2Frac := flag.Float64("ladder_tp2_frac", 1.00, "fraction of initial units to close on TP2 (1.0 = all remaining)")
	ladderSLPct := flag.Float64("ladder_sl_pct", 0.15, "ladder stop-loss: exit ≤ entry × (1 - this) closes 100%")
	ladderMaxHold := flag.Duration("ladder_max_hold", 6*time.Hour, "ladder hard timeout — closes remainder")
	// Phase 7.g lottery comparison strategy (SPEC §2.5). Parallel to momentum:
	// scan subscribed assets, open a small paper position when mid is in the
	// low-price band, hold to settlement, journal with source=lottery so PnL
	// can be diffed vs momentum. LoL gets a tighter floor because LoL upsets
	// are rare once the game starts (predictable metagame).
	lotteryEnabled := flag.Bool("lottery_enabled", true, "Phase 7.g parallel lottery strategy (low-price + hold to settlement)")
	lotteryMin := flag.Float64("lottery_min_price", 0.05, "lottery global floor; skips ≤ this")
	lotteryMax := flag.Float64("lottery_max_price", 0.30, "lottery ceiling; skips > this")
	lotteryLoLMin := flag.Float64("lottery_lol_min", 0.15, "lottery LoL-only floor; skip below (overrides global when higher)")
	lotterySize := flag.Float64("lottery_size_usd", 1.0, "lottery entry size in USDC")
	lotteryScan := flag.Duration("lottery_scan_interval", 5*time.Minute, "lottery scanner cadence")
	injuryEnabled := flag.Bool("injury_enabled", false, "enable NBA injury report scanner (ESPN API)")
	injuryInterval := flag.Duration("injury_interval", 30*time.Minute, "injury scan interval")
	injuryStarOnly := flag.Bool("injury_star_only", true, "only alert on star players (top ~3-4 per team)")
	whaleEnabled := flag.Bool("whale_enabled", false, "enable smart-money whale trade tracker")
	fadeMode := flag.Bool("fade_mode", false, "fade (mean-reversion): buy the opposite outcome when momentum fires")
	confirmDelay := flag.Duration("confirm_delay", 10*time.Second, "wait N seconds after signal trigger, re-check price before entry")
	whaleWallets := flag.String("whale_wallets", "", "tracked wallets: addr|label|minUSD|profileURL,... (comma-separated)")
	walletsFile := flag.String("wallets_file", "", "path to wallets file (one address per line) for copytrade mode")
	copytradeSize := flag.Float64("copytrade_size", 5.0, "default per-trade paper size in USDC for copytrade mode")
	walletTiersFile := flag.String("wallet_tiers", "", "path to copytrade_backtest_results.json for tiered sizing (A=$20, B=$10, C/D=default)")
	minTier := flag.String("min_tier", "", "minimum wallet tier for copytrade (A=only A, B=A+B, C=A+B+C, empty=all)")
	whaleWallet := flag.String("whale_wallet", "", "(legacy) single target wallet address (hex 0x…)")
	whaleProfile := flag.String("whale_profile", "", "(legacy) whale's Polymarket profile URL")
	whaleMinUSD := flag.Float64("whale_min_usd", 1000, "(legacy) minimum notional USD to trigger alert")
	whaleInterval := flag.Duration("whale_interval", 30*time.Second, "whale trade poll interval")
	whaleReplayWindow := flag.Duration("whale_replay_window", 0, "startup replay window for recent whale trades; 0 disables replay")
	btcEnabled := flag.Bool("btc_enabled", false, "enable BTC prediction strategy (BS first-passage vs PM gap)")
	btcInterval := flag.Duration("btc_interval", 1*time.Hour, "BTC strategy scan interval")
	btcMinGapPP := flag.Float64("btc_min_gap_pp", 7.0, "BTC minimum gap in pp to signal")
	btcTopN := flag.Int("btc_top_n", 3, "BTC max signals per scan cycle")
	btcSizeUSD := flag.Float64("btc_size_usd", 5.0, "BTC default signal size hint")
	btcDBPath := flag.String("btc_db", "db/btc.db", "SQLite path for BTC strategy data")
	updownEnabled := flag.Bool("updown_enabled", false, "enable BTC Up/Down short-term auto-trading (volume strategy)")
	updownInterval := flag.Duration("updown_interval", 15*time.Minute, "Up/Down market scan interval")
	updownConfidence := flag.Float64("updown_confidence", 0.52, "minimum Markov confidence to enter Up/Down trade")
	updownSize := flag.Float64("updown_size", 5.0, "Up/Down per-trade size in USDC")
	updownMaxDaily := flag.Int("updown_max_daily", 20, "Up/Down max bets per day")
	updownDB := flag.String("updown_db", "db/btc.db", "SQLite path for Up/Down bet tracking")
	// Phase 10: new strategy scanners
	btcDailyEnabled := flag.Bool("btc_daily_enabled", false, "enable BTC daily threshold scanner")
	btcDailyInterval := flag.Duration("btc_daily_interval", 15*time.Minute, "BTC daily threshold scan interval")
	btcDailyMinEdge := flag.Float64("btc_daily_min_edge_pp", 5.0, "BTC daily minimum edge in pp to signal")
	elonEnabled := flag.Bool("elon_enabled", false, "enable Elon tweet count scanner")
	elonInterval := flag.Duration("elon_interval", 30*time.Minute, "Elon tweet count scan interval")
	elonMinEdge := flag.Float64("elon_min_edge_pp", 5.0, "Elon minimum edge in pp to signal")
	eurovisionEnabled := flag.Bool("eurovision_enabled", false, "enable Eurovision odds scanner")
	eurovisionInterval := flag.Duration("eurovision_interval", 6*time.Hour, "Eurovision scan interval")
	eurovisionMinEdge := flag.Float64("eurovision_min_edge_pp", 5.0, "Eurovision minimum edge in pp to signal")
	liveTrading := flag.Bool("live", false, "enable real V2 CLOB order submission (requires wallet mnemonic in Bitwarden)")
	initialCapital := flag.Float64("initial_capital", 200.0, "initial capital in USD for total P&L calculation")
	positionsStatePath := flag.String("positions_state", "db/positions.json", "paper/live position state JSON path")
	riskStatePath := flag.String("risk_state", "db/risk_state.json", "risk state JSON path")
	buyTimesStatePath := flag.String("buy_times_state", "db/buy_times.json", "buy-times state JSON path")
	posMaxTotalOpenUSD := flag.Float64("pos_max_total_open_usd", 300.0, "max total open paper exposure in USD")
	posMaxOpenPositions := flag.Int("pos_max_open_positions", 60, "max concurrent open paper positions")
	posMaxPerMarketUSD := flag.Float64("pos_max_per_market_usd", 30.0, "max open paper exposure per conditionID in USD; 0 disables")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := config.LoadDotEnv(*envFile); err != nil {
		slog.Warn("dotenv_load_warn", "path", *envFile, "err", err.Error())
	}

	order.InitProxy()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch *mode {
	case "discover":
		if err := runDiscover(ctx); err != nil {
			slog.Error("discover failed", "err", err)
			os.Exit(1)
		}
	case "feed":
		if err := runFeed(ctx, *maxMarkets); err != nil {
			slog.Error("feed failed", "err", err)
			os.Exit(1)
		}
	case "sample":
		if err := runSample(ctx, *maxMarkets, *windowSec); err != nil && ctx.Err() == nil {
			slog.Error("sample failed", "err", err)
			os.Exit(1)
		}
	case "detect":
		ladderCfg := strategy.LadderConfig{
			TP1Pct:  *ladderTP1Pct,
			TP1Frac: *ladderTP1Frac,
			TP2Pct:  *ladderTP2Pct,
			TP2Frac: *ladderTP2Frac,
			SLPct:   *ladderSLPct,
			MaxHold: *ladderMaxHold,
		}
		lottCfg := strategy.LotteryConfig{
			MinPrice:     *lotteryMin,
			MaxPrice:     *lotteryMax,
			LoLMinPrice:  *lotteryLoLMin,
			SizeUSD:      *lotterySize,
			ScanInterval: *lotteryScan,
		}
		injCfg := injury.Config{
			Enabled:      *injuryEnabled,
			ScanInterval: *injuryInterval,
			StarOnly:     *injuryStarOnly,
		}
		var whaleWalletEntries []whale.WalletEntry
		if *whaleWallets != "" {
			var err error
			whaleWalletEntries, err = whale.ParseWallets(*whaleWallets)
			if err != nil {
				slog.Error("invalid -whale_wallets", "err", err)
				os.Exit(1)
			}
		}
		whaleCfg := whale.Config{
			Enabled:      *whaleEnabled,
			Wallets:      whaleWalletEntries,
			Wallet:       *whaleWallet,
			ProfileURL:   *whaleProfile,
			MinSizeUSD:   *whaleMinUSD,
			PollInterval: *whaleInterval,
			ReplayWindow: *whaleReplayWindow,
		}
		if *walletsFile != "" {
			listMinUSD := parseWhaleListMinUSD(os.Getenv("WHALE_LIST_MIN_USD"))
			fileWallets, err := whale.LoadWalletsFileWithListMins(*walletsFile, whaleCfg.MinSizeUSD, listMinUSD)
			if err != nil {
				slog.Error("copytrade.wallets_load_fail", "file", *walletsFile, "err", err)
				os.Exit(1)
			}
			whaleCfg.Wallets = fileWallets
			whaleCfg.Enabled = true
			slog.Info("copytrade.wallets_loaded", "file", *walletsFile, "count", len(fileWallets), "list_min_usd", listMinUSD)
		}
		btcCfg := btc.StrategyConfig{
			Enabled:      *btcEnabled,
			ScanInterval: *btcInterval,
			MinGapPP:     *btcMinGapPP,
			TopN:         *btcTopN,
			SizeUSD:      *btcSizeUSD,
			DBPath:       *btcDBPath,
		}
		updownCfg := btc.UpDownConfig{
			Enabled:       *updownEnabled,
			ScanInterval:  *updownInterval,
			MinConfidence: *updownConfidence,
			SizeUSD:       *updownSize,
			MaxDailyBets:  *updownMaxDaily,
			DBPath:        *updownDB,
		}
		p10 := phase10Config{
			BTCDailyEnabled:    *btcDailyEnabled,
			BTCDailyInterval:   *btcDailyInterval,
			BTCDailyMinEdge:    *btcDailyMinEdge,
			ElonEnabled:        *elonEnabled,
			ElonInterval:       *elonInterval,
			ElonMinEdge:        *elonMinEdge,
			EurovisionEnabled:  *eurovisionEnabled,
			EurovisionInterval: *eurovisionInterval,
			EurovisionMinEdge:  *eurovisionMinEdge,
		}
		if err := runDetect(ctx, *maxMarkets, *windowSec, *slippageBp, *feeBp, *largeFillUSD, *signalMode, *exitMode, *journalDir, *tickPathDir, *minEntry, *maxEntry, ladderCfg, *lotteryEnabled, lottCfg, injCfg, whaleCfg, *confirmDelay, btcCfg, updownCfg, p10, *liveTrading, *fadeMode, *walletsFile, *copytradeSize, *walletTiersFile, *initialCapital, *minTier, *positionsStatePath, *riskStatePath, *buyTimesStatePath, *posMaxTotalOpenUSD, *posMaxOpenPositions, *posMaxPerMarketUSD); err != nil && ctx.Err() == nil {
			slog.Error("detect failed", "err", err)
			os.Exit(1)
		}
	case "daily-report":
		if err := runDailyReport(ctx, *journalDir, *reportDay, *reportPush); err != nil {
			slog.Error("daily-report failed", "err", err)
			os.Exit(1)
		}
	case "daily-iterate":
		if err := runDailyIterate(ctx, *journalDir, *iterateWindow, *reportPush); err != nil {
			slog.Error("daily-iterate failed", "err", err)
			os.Exit(1)
		}
	case "arb-scan":
		if err := runArbScan(ctx); err != nil {
			slog.Error("arb-scan failed", "err", err)
			os.Exit(1)
		}
	case "prompt-test":
		if err := runPromptTest(ctx, *slippageBp); err != nil && ctx.Err() == nil {
			slog.Error("prompt-test failed", "err", err)
			os.Exit(1)
		}
	case "run":
		slog.Info("polymarket-go starting", "mode", "paper")
		// Phase 2+: strategy loop. For now, bot -mode=feed exercises the data layer.
		<-ctx.Done()
		slog.Info("shutdown")
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n", *mode)
		os.Exit(2)
	}
}

func runDiscover(ctx context.Context) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	followMkts := feed.FilterFollowTargets(all)
	mkts := feed.FilterTradablePriceBand(followMkts, followMinEntryPrice, followMaxEntryPrice)
	slog.Info("gamma.discover",
		"total_active", len(all),
		"follow_targets", len(mkts),
		"lol", len(feed.FilterLoL(all)),
		"basketball", countBy(mkts, feed.IsBasketballMarket),
		"football", countBy(mkts, feed.IsFootballMarket),
		"dota2", countBy(mkts, feed.IsDota2Market),
	)
	for _, m := range mkts {
		tokens := m.ClobTokenIDs()
		slog.Info("sports_market",
			"q", m.Question,
			"slug", m.Slug,
			"vol24h", m.Volume24hr,
			"liq_clob", m.LiquidityClob,
			"accepting", m.AcceptingOrders,
			"end", m.EndDate,
			"tokens", len(tokens),
		)
	}
	return nil
}

func runFeed(ctx context.Context, topN int) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	followMkts := feed.FilterFollowTargets(all)
	mkts := feed.FilterTradablePriceBand(followMkts, followMinEntryPrice, followMaxEntryPrice)
	if len(mkts) == 0 {
		return fmt.Errorf("no active follow-target sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	slog.Info("feed.start", "markets", len(mkts), "assets", len(assetIDs))

	ws := feed.NewWSSClient(assetIDs)

	// consumer: log tick summaries
	go func() {
		throttle := map[string]time.Time{}
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-ws.Books():
				if !ok {
					return
				}
				// throttle per-asset book logs to 1/s to keep log volume sane
				if t, seen := throttle[ev.AssetID]; seen && time.Since(t) < time.Second {
					continue
				}
				throttle[ev.AssetID] = time.Now()
				bestBid, bestAsk := 0.0, 0.0
				if len(ev.Bids) > 0 {
					bestBid = ev.Bids[0].Price
				}
				if len(ev.Asks) > 0 {
					bestAsk = ev.Asks[0].Price
				}
				slog.Info("book",
					"asset", short(ev.AssetID),
					"q", metaQ(meta, ev.AssetID),
					"bid", bestBid,
					"ask", bestAsk,
					"n_bids", len(ev.Bids),
					"n_asks", len(ev.Asks),
				)
			case tr, ok := <-ws.Trades():
				if !ok {
					return
				}
				slog.Info("trade",
					"asset", short(tr.AssetID),
					"q", metaQ(meta, tr.AssetID),
					"price", tr.Price,
					"size", tr.Size,
					"side", tr.Side,
				)
			}
		}
	}()

	return ws.Run(ctx)
}

func runSample(ctx context.Context, topN, windowSec int) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	followMkts := feed.FilterFollowTargets(all)
	mkts := feed.FilterTradablePriceBand(followMkts, followMinEntryPrice, followMaxEntryPrice)
	if len(mkts) == 0 {
		return fmt.Errorf("no active follow-target sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	slog.Info("sample.start", "markets", len(mkts), "assets", len(assetIDs), "window_sec", windowSec)

	ws := feed.NewWSSClient(assetIDs)
	sampler := feed.NewSampler(windowSec)

	go func() {
		if err := sampler.Run(ctx, ws.Books(), ws.Trades()); err != nil && ctx.Err() == nil {
			slog.Error("sampler exited", "err", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-sampler.Ticks():
				if !ok {
					return
				}
				slog.Info("tick",
					"asset", short(t.AssetID),
					"q", metaQ(meta, t.AssetID),
					"bid", t.BestBid,
					"ask", t.BestAsk,
					"mid", t.Mid,
					"trades", t.Trades,
					"buy_vol", t.BuyVol,
					"sell_vol", t.SellVol,
				)
			}
		}
	}()

	// periodic window summary, every 10s
	go func() {
		tk := time.NewTicker(10 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				for _, w := range sampler.Snapshot() {
					slog.Info("window",
						"asset", short(w.AssetID),
						"q", metaQ(meta, w.AssetID),
						"samples", w.Samples,
						"start_mid", w.StartMid,
						"end_mid", w.EndMid,
						"delta_pp", w.DeltaPP,
						"up", w.Upticks,
						"down", w.Downticks,
						"flat", w.Flats,
						"buy_ratio", w.BuyRatio,
					)
				}
			}
		}
	}()

	return ws.Run(ctx)
}

type phase10Config struct {
	BTCDailyEnabled    bool
	BTCDailyInterval   time.Duration
	BTCDailyMinEdge    float64
	ElonEnabled        bool
	ElonInterval       time.Duration
	ElonMinEdge        float64
	EurovisionEnabled  bool
	EurovisionInterval time.Duration
	EurovisionMinEdge  float64
}

type walletFileMeta struct {
	List            string
	Tier            string
	SmartMoneyScore float64
	BotScore        float64
}

func momentumSignalsEnabledForMode(signalMode string) bool {
	return signalMode != "copytrade" && signalMode != "whale"
}

func lotteryScannerEnabledForMode(signalMode string, lotteryEnabled bool) bool {
	return lotteryEnabled && signalMode != "copytrade" && signalMode != "whale"
}

func copytradeAutoAllowedForAction(action string, liveTrading bool, paperFollowPrompt bool) (bool, string) {
	if action == "" {
		return true, "legacy_tier_file"
	}
	if action == "auto-small" {
		return true, action
	}
	if action == "prompt" && paperFollowPrompt && !liveTrading {
		return true, action
	}
	return false, action
}

func drainSamplerTicks(ctx context.Context, sampler *feed.Sampler) {
	ticks := sampler.Ticks()
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ticks:
			if !ok {
				return
			}
		}
	}
}

func runDetect(ctx context.Context, topN, windowSec int, slippageBp, feeBp, largeFillUSD float64, signalMode, exitMode, journalDir, tickPathDir string, minEntry, maxEntry float64, ladderCfg strategy.LadderConfig, lotteryEnabled bool, lotteryCfg strategy.LotteryConfig, injCfg injury.Config, whaleCfg whale.Config, confirmDelay time.Duration, btcCfg btc.StrategyConfig, updownCfg btc.UpDownConfig, p10 phase10Config, liveTrading bool, fadeMode bool, walletsFile string, copytradeSize float64, walletTiersFile string, initialCapital float64, minTierFilter string, positionsStatePath, riskStatePath, buyTimesStatePath string, posMaxTotalOpenUSD float64, posMaxOpenPositions int, posMaxPerMarketUSD float64) error {
	if signalMode != "auto" && signalMode != "prompt" && signalMode != "whale" && signalMode != "copytrade" {
		return fmt.Errorf("invalid signal_mode %q (want auto|prompt|whale|copytrade)", signalMode)
	}
	if exitMode != "hold" && exitMode != "auto" && exitMode != "ladder" {
		return fmt.Errorf("invalid exit_mode %q (want hold|auto|ladder)", exitMode)
	}
	paperFollowFootballScore := signalMode == "copytrade" && !liveTrading && os.Getenv("COPYTRADE_PAPER_FOLLOW_FOOTBALL_SCORE") == "1"
	paperFootballScoreSize := parseWhaleEnvFloat("COPYTRADE_PAPER_FOOTBALL_SCORE_SIZE", 5)
	if paperFootballScoreSize <= 0 {
		paperFootballScoreSize = 5
	}
	if signalMode == "copytrade" {
		slog.Info("copytrade_football_config",
			"score_enabled", paperFollowFootballScore,
			"score_size_usd", paperFootballScoreSize,
			"live", liveTrading,
		)
	}
	momentumSignalsEnabled := momentumSignalsEnabledForMode(signalMode)
	lotteryScannerEnabled := lotteryScannerEnabledForMode(signalMode, lotteryEnabled)
	// Load wallet metadata for copytrade gradient sizing and bot/smart-money gates.
	type walletMeta struct {
		Tier            string   `json:"tier"`
		FollowAction    string   `json:"follow_action"`
		SmartMoneyScore float64  `json:"smart_money_score"`
		BotScore        float64  `json:"bot_score"`
		RiskFlags       []string `json:"risk_flags"`
	}
	walletTiers := map[string]string{}
	walletMetas := map[string]walletMeta{}
	if walletTiersFile != "" {
		raw, err := os.ReadFile(walletTiersFile)
		if err != nil {
			slog.Warn("wallet_tiers_load_fail", "file", walletTiersFile, "err", err)
		} else {
			var parsed map[string]struct {
				Tier            string   `json:"tier"`
				FollowAction    string   `json:"follow_action"`
				SmartMoneyScore float64  `json:"smart_money_score"`
				BotScore        float64  `json:"bot_score"`
				RiskFlags       []string `json:"risk_flags"`
			}
			if err := json.Unmarshal(raw, &parsed); err != nil {
				slog.Warn("wallet_tiers_parse_fail", "err", err)
			} else {
				for addr, info := range parsed {
					key := strings.ToLower(addr)
					walletTiers[key] = info.Tier
					walletMetas[key] = walletMeta{
						Tier:            info.Tier,
						FollowAction:    info.FollowAction,
						SmartMoneyScore: info.SmartMoneyScore,
						BotScore:        info.BotScore,
						RiskFlags:       info.RiskFlags,
					}
				}
				tierCounts := map[string]int{}
				for _, t := range walletTiers {
					tierCounts[t]++
				}
				slog.Info("wallet_tiers_loaded", "file", walletTiersFile, "A", tierCounts["A"], "B", tierCounts["B"], "C", tierCounts["C"], "D", tierCounts["D"])
			}
		}
	}
	walletFileMetas := map[string]walletFileMeta{}
	if walletsFile != "" {
		metas, err := loadWalletFileMetas(walletsFile)
		if err != nil {
			slog.Warn("wallet_file_meta_load_fail", "file", walletsFile, "err", err)
		} else {
			walletFileMetas = metas
			listCounts := map[string]int{}
			for _, meta := range metas {
				listCounts[meta.List]++
			}
			slog.Info("wallet_file_meta_loaded",
				"file", walletsFile,
				"core", listCounts["core"],
				"watch", listCounts["watch"],
				"sports", listCounts["sports"],
				"scout", listCounts["scout"],
				"target", listCounts["target"],
				"flow", listCounts["flow"],
				"tape", listCounts["tape"],
				"football_score", listCounts["football_score_push"],
				"other", listCounts[""],
			)
		}
	}
	copytradeAutoAllowed := func(wallet string) (bool, string) {
		meta := walletMetas[strings.ToLower(wallet)]
		return copytradeAutoAllowedForAction(meta.FollowAction, liveTrading, os.Getenv("COPYTRADE_PAPER_FOLLOW_PROMPT") == "1")
	}
	copytradeForWallet := func(wallet string) float64 {
		tier := walletTiers[strings.ToLower(wallet)]
		meta := walletMetas[strings.ToLower(wallet)]
		if meta.FollowAction == "auto-small" && meta.SmartMoneyScore >= 80 {
			return 20.0
		}
		switch tier {
		case "A":
			return 10.0
		case "B":
			return 5.0
		default:
			return copytradeSize
		}
	}

	// hold & ladder both want the settlement watcher on — hold as primary,
	// ladder as safety net (a market resolving mid-tranche clears remainder).
	wantSettlement := exitMode == "hold" || exitMode == "ladder"
	jrn, err := journal.New(journalDir)
	if err != nil {
		return fmt.Errorf("journal init: %w", err)
	}
	defer jrn.Close()
	// Phase 7.e: per-position 1Hz tick path recorder. Empty dir → noop nil
	// recorder so unit tests / ad-hoc runs can opt out cleanly.
	var recorder *tickrec.Recorder
	if tickPathDir != "" {
		recorder, err = tickrec.New(tickPathDir)
		if err != nil {
			return fmt.Errorf("tickrec init: %w", err)
		}
	}
	src := newSourceTracker()
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	followMkts := feed.FilterFollowTargets(all)
	mkts := feed.FilterTradablePriceBand(followMkts, followMinEntryPrice, followMaxEntryPrice)
	if len(mkts) == 0 {
		return fmt.Errorf("no active follow-target sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	// Per-asset sport family for lottery-mode filtering (SPEC §2.5).
	assetSport := make(map[string]strategy.SportFamily, len(meta))
	for _, m := range mkts {
		family := strategy.ClassifySport(m)
		for _, tok := range m.ClobTokenIDs() {
			if tok == "" {
				continue
			}
			assetSport[tok] = family
		}
	}
	slog.Info("detect.start",
		"markets", len(mkts),
		"follow_targets_before_price", len(followMkts),
		"lol", countBy(mkts, feed.IsLoLMarket),
		"dota2", countBy(mkts, feed.IsDota2Market),
		"basketball", countBy(mkts, feed.IsBasketballMarket),
		"football", countBy(mkts, feed.IsFootballMarket),
		"assets", len(assetIDs),
		"window_sec", windowSec,
		"fade_mode", fadeMode,
	)

	ws := feed.NewWSSClient(assetIDs)
	sampler := feed.NewSampler(windowSec)

	cfg := strategy.DefaultConfig()
	cfg.WindowSec = windowSec
	cfg.ConfirmDelay = confirmDelay
	if windowSec < cfg.MinSamplesWarm {
		cfg.MinSamplesWarm = windowSec / 2
	}
	det := strategy.NewDetector(cfg, sampler)
	for _, m := range mkts {
		tokens := m.ClobTokenIDs()
		var ids []string
		for _, t := range tokens {
			if t != "" {
				ids = append(ids, t)
			}
		}
		if len(ids) > 0 {
			det.RegisterMarket(m.ConditionID, ids)
		}
	}
	exitCfg := strategy.DefaultExitConfig()
	exit := strategy.NewExitTracker(exitCfg)
	ladder := strategy.NewLadderTracker(ladderCfg)
	posCfg := strategy.DefaultPositionConfig()
	if posMaxTotalOpenUSD > 0 {
		posCfg.MaxTotalOpenUSD = posMaxTotalOpenUSD
	}
	if posMaxOpenPositions > 0 {
		posCfg.MaxOpenPositions = posMaxOpenPositions
	}
	if posMaxPerMarketUSD >= 0 {
		posCfg.MaxPerMarketUSD = posMaxPerMarketUSD
	}
	pm := strategy.NewPositionManager(posCfg)
	if positionsStatePath == "" {
		positionsStatePath = "db/positions.json"
	}
	if err := os.MkdirAll(filepath.Dir(positionsStatePath), 0755); err != nil {
		return fmt.Errorf("positions state dir: %w", err)
	}
	if err := pm.LoadState(positionsStatePath); err != nil {
		slog.Warn("positions_load_err", "path", positionsStatePath, "err", err.Error())
	} else {
		stats := pm.Stats()
		slog.Info("positions_loaded", "path", positionsStatePath, "open", stats.Open, "closed", stats.Closed, "exposure_usd", stats.TotalExposure)
	}
	savePositions := func() {
		if err := pm.SaveState(positionsStatePath); err != nil {
			slog.Warn("positions_save_err", "path", positionsStatePath, "err", err)
		}
	}

	buyTimesMap := make(map[string]time.Time)
	if buyTimesStatePath == "" {
		buyTimesStatePath = "db/buy_times.json"
	}
	if err := os.MkdirAll(filepath.Dir(buyTimesStatePath), 0755); err != nil {
		return fmt.Errorf("buy-times state dir: %w", err)
	}
	if raw, err := os.ReadFile(buyTimesStatePath); err == nil {
		_ = json.Unmarshal(raw, &buyTimesMap)
		slog.Info("buy_times_loaded", "count", len(buyTimesMap))
	}
	saveBuyTimes := func() {
		raw, _ := json.MarshalIndent(buyTimesMap, "", "  ")
		_ = os.WriteFile(buyTimesStatePath, raw, 0644)
	}

	var orderClient order.Client
	var walletAddress string
	var onChain *order.OnChain
	pnlTrigger := make(chan struct{}, 4)
	paper := order.NewPaperClientWithFee(slippageBp, feeBp)
	orderClient = paper
	if liveTrading {
		slog.Info("v2_live_init", "msg", "loading wallet from Bitwarden")
		mnemonic, err := order.LoadMnemonicFromBitwarden("Polymarket-Go Wallet", "mnemonic")
		if err != nil {
			slog.Error("v2_wallet_load_failed", "err", err)
			os.Exit(1)
		}
		wallet, err := order.NewWalletFromMnemonic(mnemonic, "")
		if err != nil {
			slog.Error("v2_wallet_derive_failed", "err", err)
			os.Exit(1)
		}
		walletAddress = wallet.Address().Hex()
		slog.Info("v2_wallet_loaded", "address", walletAddress)
		oc, ocErr := order.NewOnChain("", wallet)
		if ocErr != nil {
			slog.Warn("onchain_init_failed", "err", ocErr)
		} else {
			onChain = oc
			slog.Info("onchain_ready", "address", walletAddress)
		}
		creds, err := order.DeriveAPIKey(order.ClobBaseURL, wallet)
		if err != nil {
			slog.Error("v2_api_key_derive_failed", "err", err)
			os.Exit(1)
		}
		slog.Info("v2_api_key_derived")
		v2Client := order.NewV2Client(wallet, creds, false)
		orderClient = v2Client
		slog.Info("v2_live_ready", "client", v2Client.Name(), "exchange", order.V2ExchangeAddress)
		if err := v2Client.CancelAllOpen(context.Background()); err != nil {
			slog.Warn("v2_cancel_all_open_failed", "err", err)
		}
	}
	slog.Info("order_client_ready", "name", orderClient.Name())
	riskCfg := risk.DefaultConfig()
	if initialCapital > 0 {
		riskCfg.StartingBankrollUSD = initialCapital
	}
	riskCfg.FeedConnected = ws.Connected
	rm := risk.New(riskCfg, time.Now())
	if riskStatePath == "" {
		riskStatePath = "db/risk_state.json"
	}
	if err := os.MkdirAll(filepath.Dir(riskStatePath), 0755); err != nil {
		return fmt.Errorf("risk state dir: %w", err)
	}
	if err := rm.LoadState(riskStatePath, time.Now()); err != nil {
		slog.Warn("risk.load_state_failed", "err", err)
	} else {
		st := rm.State()
		slog.Info("risk.state_loaded",
			"day", st.Day,
			"day_pnl", st.DayRealizedPnL,
			"cumulative_pnl", st.CumulativePnL,
			"blocked", st.Blocked,
		)
	}
	notifier := buildNotifier()
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = notifier.Close(sctx)
	}()
	pending := notify.NewPendingStore(2 * time.Hour)
	closePending := notify.NewCloseStore(2 * time.Hour)
	// Admin trigger dir: external callers (e.g. `-mode=prompt-test`) drop a JSON
	// blob into db/admin/send-prompt.trigger; the daemon watcher below picks it
	// up, emits a synthetic signal prompt, and stores the nonce in its OWN
	// pending store — so the sidecar longpoll can Claim it on callback.
	// Without this, a short-lived prompt-test subprocess registers the nonce in
	// its own memory, exits, and the callback lands on a daemon that doesn't
	// know the nonce → "已过期或已点过" even on a fresh click.
	const adminTrigger = "db/admin/send-prompt.trigger"
	const adminResume = "db/admin/resume-risk.trigger"
	_ = os.MkdirAll(filepath.Dir(adminTrigger), 0o755)
	slog.Info("daemon.startup",
		"pid", os.Getpid(),
		"args", fmt.Sprintf("%v", os.Args[1:]),
		"reason", os.Getenv("RESTART_REASON"),
	)
	slog.Info("paper_client.ready", "slippage_bp", slippageBp, "per_pos_usd", posCfg.PerPositionUSD)
	slog.Info("risk.ready",
		"bankroll_usd", riskCfg.StartingBankrollUSD,
		"daily_loss_cap_usd", rm.State().DayLossCapUSD,
		"max_single_loss_usd", riskCfg.MaxSingleLossUSD,
		"feed_silence_sec", riskCfg.FeedSilenceSec,
		"large_fill_usd", largeFillUSD,
	)
	slog.Info("signal_mode.ready", "mode", signalMode)
	if lotteryScannerEnabled {
		slog.Info("lottery.ready",
			"min_price", lotteryCfg.MinPrice,
			"max_price", lotteryCfg.MaxPrice,
			"lol_min_price", lotteryCfg.LoLMinPrice,
			"size_usd", lotteryCfg.SizeUSD,
			"scan_interval", lotteryCfg.ScanInterval.String(),
		)
	} else if lotteryEnabled {
		slog.Info("lottery.disabled_for_signal_mode", "mode", signalMode)
	}
	slog.Info("exit_mode.ready",
		"mode", exitMode,
		"want_settlement", wantSettlement,
		"fee_bp", feeBp,
		"ladder_tp1_pct", ladderCfg.TP1Pct,
		"ladder_tp1_frac", ladderCfg.TP1Frac,
		"ladder_tp2_pct", ladderCfg.TP2Pct,
		"ladder_tp2_frac", ladderCfg.TP2Frac,
		"ladder_sl_pct", ladderCfg.SLPct,
		"ladder_max_hold", ladderCfg.MaxHold.String(),
	)

	// Inbound callback consumer (Phase 3.5.b). Only runs if a DEDICATED sidecar
	// bot token is configured — we never long-poll the alert bot's token because
	// OpenClaw may also be polling it, and Telegram delivers updates competitively.
	sidecarToken := os.Getenv("SIDECAR_BOT_TOKEN")
	sidecarChat := os.Getenv("SIDECAR_CHAT_ID")
	if sidecarChat == "" {
		sidecarChat = os.Getenv("TELEGRAM_CHAT_ID")
	}
	enableSidecarLongPoll := signalMode != "copytrade" && signalMode != "whale"
	if signalMode == "copytrade" && os.Getenv("COPYTRADE_LONGPOLL") == "1" {
		enableSidecarLongPoll = true
	}
	if signalMode == "whale" && os.Getenv("WHALE_LONGPOLL") == "1" {
		enableSidecarLongPoll = true
	}
	if sidecarToken != "" && sidecarChat != "" && enableSidecarLongPoll {
		chatID, err := strconv.ParseInt(sidecarChat, 10, 64)
		if err != nil {
			slog.Warn("sidecar_chat_id_parse_fail", "err", err.Error())
		} else {
			h := &buyHandler{
				pm:            pm,
				exit:          exit,
				ladder:        ladder,
				paper:         orderClient,
				rm:            rm,
				pending:       pending,
				closePending:  closePending,
				notifier:      notifier,
				meta:          meta,
				src:           src,
				recorder:      recorder,
				jrn:           jrn,
				largeFillUSD:  largeFillUSD,
				exitMode:      exitMode,
				riskStatePath: riskStatePath,
				savePositions: savePositions,
			}
			lp := notify.NewLongPoll(notify.LongPollConfig{
				BotToken:       sidecarToken,
				ExpectedChatID: chatID,
			}, h)
			go func() {
				slog.Info("sidecar_longpoll.ready", "chat_id", chatID)
				if err := lp.Run(ctx); err != nil && ctx.Err() == nil {
					slog.Warn("sidecar_longpoll_exit", "err", err.Error())
				}
			}()
		}
	} else if (signalMode == "copytrade" || signalMode == "whale") && !enableSidecarLongPoll {
		slog.Info("sidecar_longpoll.skip", "mode", signalMode, "reason", "push_only")
	} else if signalMode == "prompt" {
		slog.Warn("signal_mode_prompt_without_sidecar",
			"hint", "prompt mode needs SIDECAR_BOT_TOKEN + chat_id — buttons will arrive but clicks won't be consumed")
	}

	// Pending-store reaper so expired button prompts don't accumulate.
	// For each evicted entry we rewrite the original DM to "已过期" and strip
	// its keyboard — so the boss's chat history shows the outcome of every
	// prompt (Phase 3.5 B).
	go func() {
		tk := time.NewTicker(15 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				evicted := pending.Reap(now)
				if len(evicted) == 0 {
					continue
				}
				edited := 0
				for _, p := range evicted {
					if p.MessageID != 0 {
						notifier.EditSignalExpired(p.MessageID)
						edited++
					}
				}
				slog.Info("pending_reap",
					"expired", len(evicted),
					"edited_expired_dm", edited,
					"remaining", pending.Size(),
				)
			}
		}
	}()

	// Close-prompt reaper — same pattern as the buy-prompt reaper above.
	go func() {
		tk := time.NewTicker(15 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				evicted := closePending.Reap(now)
				for _, ci := range evicted {
					if ci.MessageID != 0 {
						notifier.EditCloseDone("⌛ 已过期 · 未操作", ci.MessageID)
					}
				}
			}
		}
	}()

	// Admin trigger watcher: 1 Hz poll of db/admin/send-prompt.trigger.
	// Any process (e.g. `-mode=prompt-test`) that drops a JSON file here gets
	// a synthetic prompt emitted by *this* daemon, with the nonce registered
	// in the shared pending store the longpoll consumer reads from.
	go func() {
		tk := time.NewTicker(1 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				if data, err := os.ReadFile(adminTrigger); err == nil {
					_ = os.Remove(adminTrigger)
					if err := sendAdminPrompt(data, mkts, meta, sampler, pending, notifier); err != nil {
						slog.Warn("admin_prompt_fail", "err", err.Error())
					}
				}
				if _, err := os.Stat(adminResume); err == nil {
					_ = os.Remove(adminResume)
					rm.Resume()
					if err := rm.SaveState(riskStatePath); err != nil {
						slog.Warn("risk_save_err", "err", err)
					}
					slog.Info("risk_admin_resume", "by", "trigger_file")
					notifier.RiskResume(notify.RiskResumeEvent{})
				}
			}
		}
	}()

	go func() {
		if err := sampler.Run(ctx, ws.Books(), ws.Trades()); err != nil && ctx.Err() == nil {
			slog.Error("sampler exited", "err", err)
		}
	}()
	if momentumSignalsEnabled {
		go func() {
			if err := det.Run(ctx); err != nil && ctx.Err() == nil {
				slog.Error("detector exited", "err", err)
			}
		}()
	} else {
		slog.Info("momentum_detector.disabled_for_signal_mode", "mode", signalMode)
		go drainSamplerTicks(ctx, sampler)
	}

	// Fan-out ticks to the exit tracker (only tracks opened positions).
	// Uses a fresh Sampler subscription via a side goroutine: we tap the detector's
	// upstream by reading ticks *through* the sampler's Ticks() channel which the
	// detector already consumes. To avoid a fight for one channel, run a dedicated
	// tap via TickTail polling on each detected open asset instead.
	// Simpler: have detect subscribe to ticks directly alongside the detector.
	//
	// Here we piggyback on the fact that only the detector consumes sampler.Ticks().
	// Instead of stealing them, expose a separate TickTail-based poller below.
	// Update: cleanest is to have the sampler fan-out — but we only have one
	// consumer right now. Workaround: poll sampler.Window for open positions.
	go func() {
		tk := time.NewTicker(1 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				for _, w := range sampler.Snapshot() {
					if !exit.Has(w.AssetID) {
						continue
					}
					tail, ok := sampler.TickTail(w.AssetID, 1)
					if !ok || len(tail) == 0 {
						continue
					}
					if sig, fired := exit.OnTick(tail[0]); fired {
						sellIntent := order.Intent{
							AssetID: sig.AssetID,
							Market:  sig.Market,
							Side:    order.Sell,
							SizeUSD: posCfg.PerPositionUSD,
							LimitPx: sig.ExitMid,
							Type:    order.GTC,
						}
						res, err := orderClient.Submit(ctx, sellIntent)
						if err != nil {
							slog.Warn("paper_sell_reject",
								"asset", short(sig.AssetID),
								"limit", sig.ExitMid,
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							slog.Warn("sell_not_filled",
								"asset", short(sig.AssetID),
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						sig.ExitMid = res.AvgPrice
						sig.ChangePP = (res.AvgPrice - sig.EntryMid) * 100
						closed, err := pm.CloseFirstByAsset(sig.AssetID, sig)
						if err != nil {
							slog.Warn("paper_close_miss", "asset", short(sig.AssetID), "err", err.Error())
							continue
						}
						savePositions()
						if recorder != nil {
							if rerr := recorder.Stop(closed.ID); rerr != nil {
								slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
							}
						}
						entryFee := closed.OpenFeeUSD
						exitFee := res.FeeUSD
						netPnL := closed.PnLUSD - entryFee - exitFee
						stats := pm.Stats()
						posSource, _ := src.Peek(closed.ID)
						if posSource != "manual" {
							if tripped := rm.OnClose(netPnL, sig.Time); tripped {
								rst := rm.State()
								slog.Error("risk_trip",
									"reason", string(rst.BlockReason),
									"day_pnl_usd", rst.DayRealizedPnL,
									"cap_usd", rst.DayLossCapUSD,
								)
								notifier.RiskTrip(notify.RiskTripEvent{
									Reason:        string(rst.BlockReason),
									DayPnLUSD:     rst.DayRealizedPnL,
									DayLossCapUSD: rst.DayLossCapUSD,
									DrawdownUSD:   rst.DrawdownUSD,
									DrawdownCap:   rst.DrawdownCapUSD,
									OpenPositions: stats.Open,
								})
							}
							if err := rm.SaveState(riskStatePath); err != nil {
								slog.Warn("risk_save_err", "err", err)
							}
						}
						if netPnL <= -largeFillUSD || netPnL >= largeFillUSD {
							notifier.LargeFill(notify.LargeFillEvent{
								Question: metaQ(meta, sig.AssetID),
								AssetID:  sig.AssetID,
								Side:     "sell",
								SizeUSD:  posCfg.PerPositionUSD,
								PnLUSD:   netPnL,
								EntryPx:  sig.EntryMid,
								ExitPx:   res.AvgPrice,
								Reason:   string(sig.Reason),
								HeldSec:  int(sig.HeldFor.Seconds()),
							})
						}
						source, openOID := src.Take(closed.ID)
						if err := jrn.Append(journal.TradeRecord{
							ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
							Question:     metaQ(meta, closed.AssetID),
							Outcome:      metaOutcome(meta, closed.AssetID),
							Side:         "buy",
							SizeUSD:      closed.SizeUSD,
							Units:        closed.Units,
							EntryMid:     closed.EntryMid,
							EntryTime:    closed.EntryTime,
							ExitMid:      closed.ExitMid,
							ExitTime:     closed.ExitTime,
							ExitReason:   string(closed.ExitReason),
							HeldSec:      int(sig.HeldFor.Seconds()),
							PnLUSD:       closed.PnLUSD,
							EntryFeeUSD:  entryFee,
							ExitFeeUSD:   exitFee,
							NetPnLUSD:    netPnL,
							OpenOrderID:  openOID,
							CloseOrderID: res.OrderID,
							Mode:         "paper",
							SignalSource: source,
						}); err != nil {
							slog.Warn("journal_append_fail", "asset", short(sig.AssetID), "err", err.Error())
						}
						slog.Info("exit",
							"asset", short(sig.AssetID),
							"q", metaQ(meta, sig.AssetID),
							"reason", string(sig.Reason),
							"order_id", res.OrderID,
							"entry", sig.EntryMid,
							"peak", sig.PeakMid,
							"exit_fill", res.AvgPrice,
							"delta_pp", sig.ChangePP,
							"drawdown_pp", sig.DrawdownPP,
							"held_sec", int(sig.HeldFor.Seconds()),
							"gross_pnl_usd", closed.PnLUSD,
							"entry_fee_usd", entryFee,
							"exit_fee_usd", exitFee,
							"net_pnl_usd", netPnL,
							"open_positions", stats.Open,
							"realized_pnl", stats.RealizedPnLUSD,
						)
					}
				}
			}
		}
	}()

	// Ladder exit-watch: runs in parallel to the auto exit-watch. Only
	// ladder-tracked positions fire here — it polls posID directly so
	// stacked positions on the same asset track independently.
	if exitMode == "ladder" {
		go func() {
			tk := time.NewTicker(1 * time.Second)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					for _, p := range pm.Snapshot() {
						if !ladder.Has(p.ID) {
							continue
						}
						tail, ok := sampler.TickTail(p.AssetID, 1)
						if !ok || len(tail) == 0 {
							continue
						}
						ex, fired := ladder.OnTick(p.ID, tail[0])
						if !fired {
							continue
						}
						notional := ex.CloseUnits * ex.ExitMid
						sellIntent := order.Intent{
							AssetID: ex.AssetID,
							Market:  ex.Market,
							Side:    order.Sell,
							SizeUSD: notional,
							LimitPx: ex.ExitMid,
							Type:    order.GTC,
						}
						res, err := orderClient.Submit(ctx, sellIntent)
						if err != nil {
							slog.Warn("paper_ladder_sell_reject",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"limit", ex.ExitMid,
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							slog.Warn("ladder_sell_not_filled",
								"pos", p.ID,
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						ex.ExitMid = res.AvgPrice
						esig := strategy.ExitSignal{
							AssetID:  ex.AssetID,
							Market:   ex.Market,
							Time:     ex.Time,
							EntryMid: ex.EntryMid,
							PeakMid:  ex.ExitMid,
							ExitMid:  ex.ExitMid,
							HeldFor:  ex.HeldFor,
							ChangePP: (ex.ExitMid - ex.EntryMid) * 100,
							Reason:   ex.Reason,
						}
						closedTranche, cerr := pm.PartialClose(p.ID, ex.CloseUnits, esig)
						if cerr != nil {
							slog.Warn("ladder_partial_close_fail",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"err", cerr.Error())
							continue
						}
						savePositions()
						if ex.Final && recorder != nil {
							if rerr := recorder.Stop(p.ID); rerr != nil {
								slog.Warn("tickrec_stop_fail", "pos", p.ID, "err", rerr.Error())
							}
						}
						// Apportion the open fee across tranches by unit share.
						entryFeeShare := 0.0
						if p.InitUnits > 0 {
							entryFeeShare = p.OpenFeeUSD * (ex.CloseUnits / p.InitUnits)
						}
						exitFee := res.FeeUSD
						netPnL := closedTranche.PnLUSD - entryFeeShare - exitFee
						stats := pm.Stats()
						ladderSource, _ := src.Peek(p.ID)
						if ladderSource != "manual" {
							if tripped := rm.OnClose(netPnL, ex.Time); tripped {
								rst := rm.State()
								slog.Error("risk_trip",
									"reason", string(rst.BlockReason),
									"day_pnl_usd", rst.DayRealizedPnL,
									"cap_usd", rst.DayLossCapUSD,
								)
								notifier.RiskTrip(notify.RiskTripEvent{
									Reason:        string(rst.BlockReason),
									DayPnLUSD:     rst.DayRealizedPnL,
									DayLossCapUSD: rst.DayLossCapUSD,
									DrawdownUSD:   rst.DrawdownUSD,
									DrawdownCap:   rst.DrawdownCapUSD,
									OpenPositions: stats.Open,
								})
							}
							if err := rm.SaveState(riskStatePath); err != nil {
								slog.Warn("risk_save_err", "err", err)
							}
						}
						if netPnL <= -largeFillUSD || netPnL >= largeFillUSD {
							notifier.LargeFill(notify.LargeFillEvent{
								Question: metaQ(meta, ex.AssetID),
								AssetID:  ex.AssetID,
								Side:     "sell",
								SizeUSD:  notional,
								PnLUSD:   netPnL,
								EntryPx:  ex.EntryMid,
								ExitPx:   res.AvgPrice,
								Reason:   string(ex.Reason),
								HeldSec:  int(ex.HeldFor.Seconds()),
							})
						}
						// Source stays keyed by posID; Take only on the final
						// tranche so earlier tranches can still attribute.
						var source, openOID string
						if ex.Final {
							source, openOID = src.Take(p.ID)
						} else {
							source, openOID = src.Peek(p.ID)
						}
						trancheID := closedTranche.ID + "." + ex.Tranche
						if err := jrn.Append(journal.TradeRecord{
							ID:           trancheID,
							AssetID:      closedTranche.AssetID,
							Market:       closedTranche.Market,
							Question:     metaQ(meta, closedTranche.AssetID),
							Outcome:      metaOutcome(meta, closedTranche.AssetID),
							Side:         "buy",
							SizeUSD:      closedTranche.SizeUSD,
							Units:        closedTranche.Units,
							EntryMid:     closedTranche.EntryMid,
							EntryTime:    closedTranche.EntryTime,
							ExitMid:      closedTranche.ExitMid,
							ExitTime:     closedTranche.ExitTime,
							ExitReason:   string(closedTranche.ExitReason),
							HeldSec:      int(ex.HeldFor.Seconds()),
							PnLUSD:       closedTranche.PnLUSD,
							EntryFeeUSD:  entryFeeShare,
							ExitFeeUSD:   exitFee,
							NetPnLUSD:    netPnL,
							Tranche:      ex.Tranche,
							OpenOrderID:  openOID,
							CloseOrderID: res.OrderID,
							Mode:         "paper",
							SignalSource: source,
						}); err != nil {
							slog.Warn("journal_append_fail",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"err", err.Error())
						}
						slog.Info("ladder_exit",
							"pos", p.ID,
							"asset", short(ex.AssetID),
							"q", metaQ(meta, ex.AssetID),
							"tranche", ex.Tranche,
							"reason", string(ex.Reason),
							"final", ex.Final,
							"order_id", res.OrderID,
							"entry", ex.EntryMid,
							"exit_fill", res.AvgPrice,
							"close_units", ex.CloseUnits,
							"held_sec", int(ex.HeldFor.Seconds()),
							"gross_pnl_usd", closedTranche.PnLUSD,
							"entry_fee_usd", entryFeeShare,
							"exit_fee_usd", exitFee,
							"net_pnl_usd", netPnL,
							"open_positions", stats.Open,
							"realized_pnl", stats.RealizedPnLUSD,
						)
						if ex.Reason == strategy.ExitLadderSL {
							var cid string
							if me := meta[ex.AssetID]; me != nil {
								cid = me.ConditionID
							}
							det.NotifySL(ex.AssetID, cid)
							slog.Info("sl_cooldown_extended",
								"asset", short(ex.AssetID),
								"market", short(cid),
								"cooldown", det.CooldownAfterSL().String())
						}
					}
				}
			}
		}()
	}

	// Phase 7.e: 1Hz tick-path recorder. Runs independent of exit_mode so we
	// capture the full post-open trajectory in hold-mode paper runs too. Each
	// second, snapshot open recordings and persist the latest sampler tick;
	// Recorder dedupes within-second on its side so this can safely over-fire.
	if recorder != nil {
		go func() {
			tk := time.NewTicker(1 * time.Second)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					for posID, assetID := range recorder.Snapshot() {
						tail, ok := sampler.TickTail(assetID, 1)
						if !ok || len(tail) == 0 {
							continue
						}
						if err := recorder.Record(posID, tail[0]); err != nil {
							slog.Warn("tickrec_record_fail", "pos", posID, "err", err.Error())
						}
					}
				}
			}
		}()
	}

	// Injury scanner: created before signal/lottery goroutines so both can
	// call injScanner.HasInjuredStar(). Returns nil/false when disabled.
	injScanner := injury.NewScanner(injCfg, "db")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case sig, ok := <-det.Signals():
				if !ok {
					return
				}
				if !momentumSignalsEnabled {
					slog.Debug("momentum_skip_signal_mode", "mode", signalMode, "asset", short(sig.AssetID))
					continue
				}
				slog.Info("signal",
					"asset", short(sig.AssetID),
					"q", metaQ(meta, sig.AssetID),
					"mid", sig.Mid,
					"delta_pp", sig.DeltaPP,
					"tail_ups", sig.TailUps,
					"tail_len", sig.TailLen,
					"buy_ratio", sig.BuyRatio,
					"reason", sig.Reason,
				)
				// Risk gate first — daily-loss breaker / feed-silence / manual pause.
				if err := rm.AllowOpen(time.Now()); err != nil {
					st := rm.State()
					slog.Warn("risk_block_open",
						"asset", short(sig.AssetID),
						"q", metaQ(meta, sig.AssetID),
						"reason", string(st.BlockReason),
						"day_pnl_usd", st.DayRealizedPnL,
						"cap_usd", st.DayLossCapUSD,
					)
					continue
				}
				// Phase 7.a entry-price band filter — winners in python DB clustered
				// in [0.15, 0.70]; edges (<0.15 bleed to zero, >0.70 favorites wipe
				// out) were losers. Signal still logs; only the prompt is suppressed.
				if sig.Mid < minEntry || sig.Mid > maxEntry {
					slog.Info("signal_filtered_price_band",
						"asset", short(sig.AssetID),
						"q", metaQ(meta, sig.AssetID),
						"mid", sig.Mid,
						"min", minEntry,
						"max", maxEntry,
					)
					continue
				}
				// Injury filter: if this is a basketball market and the team
				// we'd be betting on has a star OUT/Doubtful, block the trade.
				if injCfg.Enabled {
					if blocked, team, players := injuryBlocksMomentum(injScanner, assetSport, meta, sig.AssetID); blocked {
						slog.Info("signal_blocked_injury",
							"asset", short(sig.AssetID),
							"q", metaQ(meta, sig.AssetID),
							"mid", sig.Mid,
							"team", team,
							"injured_stars", players,
							"delta_pp", sig.DeltaPP,
						)
						continue
					}
				}
				// Fade mode: buy the opposite outcome (mean-reversion).
				buyAssetID := sig.AssetID
				buyMid := sig.Mid
				if fadeMode {
					me := meta[sig.AssetID]
					if me == nil || me.Sibling == "" {
						slog.Info("fade_skip_no_sibling",
							"asset", short(sig.AssetID),
							"q", metaQ(meta, sig.AssetID),
						)
						continue
					}
					buyAssetID = me.Sibling
					buyMid = 1.0 - sig.Mid
					if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
						buyMid = w.EndMid
					}
					if buyMid < minEntry || buyMid > maxEntry {
						slog.Info("fade_filtered_price_band",
							"asset", short(buyAssetID),
							"q", metaQ(meta, buyAssetID),
							"mid", buyMid,
							"min", minEntry,
							"max", maxEntry,
						)
						continue
					}
					slog.Info("fade_swap",
						"orig_asset", short(sig.AssetID),
						"fade_asset", short(buyAssetID),
						"orig_mid", sig.Mid,
						"fade_mid", buyMid,
						"q", metaQ(meta, buyAssetID),
					)
				}

				// Prompt mode: publish the signal as a DM with one button row per
				// outcome (YES/NO or team-A/team-B) and stash the full Choices slice
				// in the pending store. The callback longpoll (above) claims the
				// nonce, picks Choices[slot], and executes via buyHandler.
				if signalMode == "prompt" {
					me := meta[sig.AssetID]
					choices := []notify.Choice{{
						AssetID: sig.AssetID, Outcome: outcomeOrDefault(me, "Yes"),
						Mid: sig.Mid, IsSignal: true,
					}}
					sigChoices := []notify.SignalChoice{{
						Slot: 0, Outcome: choices[0].Outcome, Mid: sig.Mid, IsSignal: true,
					}}
					if me != nil && me.Sibling != "" {
						sibMid := 1.0 - sig.Mid // fallback: binary complement
						if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
							sibMid = w.EndMid
						}
						sibOutcome := me.SiblingOutcome
						if sibOutcome == "" {
							sibOutcome = "No"
						}
						choices = append(choices, notify.Choice{
							AssetID: me.Sibling, Outcome: sibOutcome, Mid: sibMid,
						})
						sigChoices = append(sigChoices, notify.SignalChoice{
							Slot: 1, Outcome: sibOutcome, Mid: sibMid,
						})
					}
					p := pending.Put(notify.PendingIntent{
						Market:   sig.Market,
						Question: metaQ(meta, sig.AssetID),
						Choices:  choices,
					}, time.Now())
					var match, ctxLine, endIn string
					if me != nil {
						match = me.Match
						ctxLine = me.Context
						endIn = notify.HumanizeEndIn(time.Now(), me.EndTime)
					}
					nonceSnap := p.Nonce
					slugVal := ""
					if me != nil {
						slugVal = me.Slug
					}
					notifier.SignalPrompt(notify.SignalPromptEvent{
						Nonce:     p.Nonce,
						Match:     match,
						Context:   ctxLine,
						EndIn:     endIn,
						Choices:   sigChoices,
						DeltaPP:   sig.DeltaPP,
						TailUps:   sig.TailUps,
						TailLen:   sig.TailLen,
						BuyRatio:  sig.BuyRatio,
						Slug:      slugVal,
						ExpiresIn: 10 * time.Minute,
						OnSent: func(msgID int64, err error) {
							if err != nil {
								slog.Warn("prompt_send_err", "err", err)
								return
							}
							if msgID == 0 {
								return
							}
							pending.SetMessageID(nonceSnap, msgID)
						},
					})
					slog.Info("signal_prompt_sent",
						"asset", short(sig.AssetID),
						"nonce", p.Nonce,
						"mid", sig.Mid,
						"choices", len(choices),
					)
					continue
				}

				buyIntent := order.Intent{
					AssetID: buyAssetID,
					Market:  sig.Market,
					Side:    order.Buy,
					SizeUSD: posCfg.PerPositionUSD,
					LimitPx: buyMid,
					Type:    order.GTC,
				}
				res, err := orderClient.Submit(ctx, buyIntent)
				if err != nil {
					slog.Warn("paper_buy_reject",
						"asset", short(buyAssetID),
						"limit", buyMid,
						"err", err.Error())
					continue
				}
				if res.Status != order.StatusFilled {
					slog.Warn("buy_not_filled",
						"asset", short(buyAssetID),
						"order_id", res.OrderID,
						"status", res.Status)
					continue
				}
				entryTick := feed.Tick{
					AssetID: buyAssetID, Market: sig.Market,
					Time: sig.Time, Mid: res.AvgPrice,
				}
				pos, err := pm.Open(buyAssetID, sig.Market, entryTick)
				if err != nil {
					slog.Info("paper_open_skip",
						"asset", short(buyAssetID),
						"q", metaQ(meta, buyAssetID),
						"order_id", res.OrderID,
						"reason", err.Error(),
					)
					continue
				}
				if err := pm.SetOpenFee(pos.ID, res.FeeUSD); err != nil {
					slog.Warn("set_open_fee_err", "pos", pos.ID, "err", err)
				}
				savePositions()
				switch exitMode {
				case "auto":
					exit.Open(buyAssetID, sig.Market, entryTick)
				case "ladder":
					ladder.Open(pos.ID, sig.Market, buyAssetID, entryTick, pos.Units)
				}
				if recorder != nil {
					if rerr := recorder.Start(pos.ID, buyAssetID); rerr != nil {
						slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
					}
				}
				src.Mark(pos.ID, "auto", res.OrderID)
				stats := pm.Stats()
				slog.Info("paper_open",
					"id", pos.ID,
					"order_id", res.OrderID,
					"asset", short(buyAssetID),
					"q", metaQ(meta, buyAssetID),
					"signal_mid", sig.Mid,
					"entry_fill", res.AvgPrice,
					"size_usd", pos.SizeUSD,
					"units", pos.Units,
					"open_positions", stats.Open,
					"total_exposure_usd", stats.TotalExposure,
				)

				// Auto mode also sends a signal DM with buttons so the boss
				// can see what happened and optionally add a manual position.
				me := meta[sig.AssetID]
				sigChoices := []notify.SignalChoice{{
					Slot: 0, Outcome: outcomeOrDefault(me, "Yes"),
					Mid: sig.Mid, IsSignal: true,
				}}
				dmChoices := []notify.Choice{{
					AssetID: sig.AssetID, Outcome: sigChoices[0].Outcome,
					Mid: sig.Mid, IsSignal: true,
				}}
				if me != nil && me.Sibling != "" {
					sibMid := 1.0 - sig.Mid
					if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
						sibMid = w.EndMid
					}
					sibOut := me.SiblingOutcome
					if sibOut == "" {
						sibOut = "No"
					}
					dmChoices = append(dmChoices, notify.Choice{
						AssetID: me.Sibling, Outcome: sibOut, Mid: sibMid,
					})
					sigChoices = append(sigChoices, notify.SignalChoice{
						Slot: 1, Outcome: sibOut, Mid: sibMid,
					})
				}
				p := pending.Put(notify.PendingIntent{
					Market:   sig.Market,
					Question: metaQ(meta, sig.AssetID),
					Choices:  dmChoices,
				}, time.Now())
				var match, ctxLine, endIn string
				if me != nil {
					match = me.Match
					ctxLine = me.Context
					endIn = notify.HumanizeEndIn(time.Now(), me.EndTime)
				}
				nonceSnap := p.Nonce
				autoSlug := ""
				if me != nil {
					autoSlug = me.Slug
				}
				notifier.SignalPrompt(notify.SignalPromptEvent{
					Nonce:     p.Nonce,
					Match:     match,
					Context:   ctxLine,
					EndIn:     endIn,
					Choices:   sigChoices,
					DeltaPP:   sig.DeltaPP,
					TailUps:   sig.TailUps,
					TailLen:   sig.TailLen,
					BuyRatio:  sig.BuyRatio,
					Slug:      autoSlug,
					ExpiresIn: 2 * time.Hour,
					OnSent: func(msgID int64, err error) {
						if err != nil {
							slog.Warn("notify_send_err", "err", err)
							return
						}
						if msgID == 0 {
							return
						}
						pending.SetMessageID(nonceSnap, msgID)
					},
				})
				slog.Info("auto_signal_dm_sent",
					"asset", short(sig.AssetID),
					"nonce", p.Nonce,
					"auto_order", res.OrderID,
					"mid", sig.Mid,
				)
			}
		}
	}()

	// Phase 7.g lottery scanner: periodically scan for low-price underdog
	// assets, open small paper positions, hold to settlement. Journal with
	// source=lottery so PnL can be compared vs momentum strategy.
	lotteryOpen := make(map[string]bool) // assetID → already has lottery position (guarded by single-writer goroutine)
	for _, p := range pm.Snapshot() {
		lotteryOpen[p.AssetID] = true
	}
	if lotteryScannerEnabled {
		go func() {
			tk := time.NewTicker(lotteryCfg.ScanInterval)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					if err := rm.AllowOpen(time.Now()); err != nil {
						continue
					}
					candidates := strategy.ScanEligible(sampler, assetSport, lotteryCfg)
					for _, c := range candidates {
						if lotteryOpen[c.AssetID] {
							continue
						}
						// Injury filter: skip basketball underdogs whose stars are OUT.
						if injCfg.Enabled && c.Sport == strategy.SportBasketball {
							if blocked, team, players := injuryBlocksLottery(injScanner, meta, c.AssetID); blocked {
								slog.Info("lottery_blocked_injury",
									"asset", short(c.AssetID),
									"q", metaQ(meta, c.AssetID),
									"mid", c.Mid,
									"team", team,
									"injured_stars", players,
								)
								continue
							}
						}
						// Volatility filter: skip assets with choppy in-play price action.
						if ws, ok := sampler.Window(c.AssetID); ok {
							vr := strategy.IsVolatile(ws)
							if vr.Volatile {
								slog.Info("lottery_blocked_volatile",
									"asset", short(c.AssetID),
									"q", metaQ(meta, c.AssetID),
									"mid", c.Mid,
									"sport", string(c.Sport),
									"delta_pp", vr.DeltaPP,
									"upticks", vr.Upticks,
									"downticks", vr.Downticks,
									"samples", vr.Samples,
								)
								continue
							}
						}
						buyIntent := order.Intent{
							AssetID: c.AssetID,
							Market:  c.Market,
							Side:    order.Buy,
							SizeUSD: lotteryCfg.SizeUSD,
							LimitPx: c.Mid,
							Type:    order.GTC,
						}
						res, err := orderClient.Submit(ctx, buyIntent)
						if err != nil {
							slog.Warn("lottery_buy_reject",
								"asset", short(c.AssetID),
								"mid", c.Mid,
								"sport", string(c.Sport),
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							slog.Warn("lottery_buy_not_filled",
								"asset", short(c.AssetID),
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						entryTick := feed.Tick{
							AssetID: c.AssetID, Market: c.Market,
							Time: c.Time, Mid: res.AvgPrice,
						}
						pos, err := pm.OpenSized(c.AssetID, c.Market, entryTick, lotteryCfg.SizeUSD)
						if err != nil {
							slog.Info("lottery_open_skip",
								"asset", short(c.AssetID),
								"q", metaQ(meta, c.AssetID),
								"reason", err.Error())
							continue
						}
						if err := pm.SetOpenFee(pos.ID, res.FeeUSD); err != nil {
							slog.Warn("set_open_fee_err", "pos", pos.ID, "err", err)
						}
						savePositions()
						lotteryOpen[c.AssetID] = true
						src.Mark(pos.ID, "lottery", res.OrderID)
						stats := pm.Stats()
						slog.Info("lottery_open",
							"id", pos.ID,
							"order_id", res.OrderID,
							"asset", short(c.AssetID),
							"q", metaQ(meta, c.AssetID),
							"mid", c.Mid,
							"sport", string(c.Sport),
							"entry_fill", res.AvgPrice,
							"size_usd", pos.SizeUSD,
							"units", pos.Units,
							"open_positions", stats.Open,
							"total_exposure_usd", stats.TotalExposure,
						)
					}
				}
			}
		}()
	}

	// NBA injury scanner: periodic ESPN API poll + DM notification.
	// injScanner is created unconditionally so momentum/lottery filters
	// can call HasInjuredStar() — it returns nil when disabled.
	if injCfg.Enabled {
		slog.Info("injury_scanner.ready",
			"interval", injCfg.ScanInterval.String(),
			"star_only", injCfg.StarOnly,
		)
		go func() {
			tk := time.NewTicker(injCfg.ScanInterval)
			defer tk.Stop()
			// Immediate first scan on startup (don't wait for first tick).
			processInjuryAlerts := func(alerts []injury.InjuryAlert) {
				// Group alerts by game matchup so we push ONE notification per game.
				type gameKey struct{ teamA, teamB string }
				type gameGroup struct {
					key    gameKey
					alerts []injury.InjuryAlert
				}
				groups := make(map[gameKey]*gameGroup)
				var order []gameKey
				for _, a := range alerts {
					if !injuryTeamInMarkets(a.Team, meta, assetSport) {
						// Fallback: also allow if ESPN scoreboard has an upcoming game for this team
						if _, hasGame := injScanner.GameFor(a.Team); !hasGame {
							slog.Info("injury_skip_no_market", "team", a.Team, "player", a.StarPlayer)
							continue
						}
					}
					if injuryGameFinished(a.Team, injScanner) {
						slog.Info("injury_skip_final", "team", a.Team, "player", a.StarPlayer)
						continue
					}
					slog.Info("injury_alert",
						"team", a.Team,
						"player", a.StarPlayer,
						"status", string(a.Status),
						"impact", a.Impact,
					)
					// Normalize key so both teams in a matchup collapse to the same group.
					opp := injuryFindOpponentName(a.Team, injScanner, meta, assetSport)
					k := gameKey{a.Team, opp}
					// Ensure consistent ordering: if opponent already started a group, merge.
					if opp != "" {
						if _, ok := groups[gameKey{opp, a.Team}]; ok {
							k = gameKey{opp, a.Team}
						}
					}
					if _, ok := groups[k]; !ok {
						groups[k] = &gameGroup{key: k}
						order = append(order, k)
					}
					groups[k].alerts = append(groups[k].alerts, a)
				}
				for _, k := range order {
					g := groups[k]
					ev := injuryBuildGameEvent(g.alerts, injScanner, meta, assetSport, sampler)
					notifier.InjuryAlert(ev)
					// Push opponent buy prompt once per game (use first alert with OUT/DTD).
					for _, a := range g.alerts {
						if a.Status == injury.StatusOut || a.Status == injury.StatusDTD {
							injuryPushOpponentPrompt(a, meta, assetSport, sampler, pending, notifier)
							break
						}
					}
				}
			}
			scanOnce := func() {
				alerts, err := injScanner.Scan(ctx)
				if err != nil {
					slog.Warn("injury_scan_fail", "err", err.Error())
					return
				}
				processInjuryAlerts(alerts)
			}
			scanOnce()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					alerts, err := injScanner.Scan(ctx)
					if err != nil {
						slog.Warn("injury_scan_fail", "err", err.Error())
						continue
					}
					processInjuryAlerts(alerts)
				}
			}
		}()
	}

	// Whale trades JSONL log — records every detected trade for analysis.
	var whaleLogMu sync.Mutex
	whaleLogPath := filepath.Join(journalDir, "whale_trades.jsonl")
	appendWhaleTrade := func(ev whale.AlertEvent, action, reason string) {
		walletKey := strings.ToLower(ev.Wallet)
		fileMeta := walletFileMetas[walletKey]
		tier := fileMeta.Tier
		if tier == "" {
			tier = walletTiers[walletKey]
		}
		rec := map[string]interface{}{
			"ts":           ev.Timestamp.In(journal.SGT).Format(time.RFC3339),
			"wallet":       ev.Wallet,
			"label":        ev.Label,
			"side":         ev.Side,
			"market":       ev.Question,
			"outcome":      ev.Outcome,
			"price":        ev.Price,
			"size":         ev.Notional,
			"units":        ev.SizeUnits,
			"asset_id":     ev.AssetID,
			"condition_id": ev.ConditionID,
			"trade_id":     ev.TradeID,
			"action":       action,
			"reason":       reason,
			"mode":         signalMode,
			"wallets":      walletsFile,
			"list":         fileMeta.List,
			"tier":         tier,
			"smart":        fileMeta.SmartMoneyScore,
			"bot":          fileMeta.BotScore,
		}
		line, _ := json.Marshal(rec)
		whaleLogMu.Lock()
		defer whaleLogMu.Unlock()
		if f, err := os.OpenFile(whaleLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			_, _ = f.Write(append(line, '\n'))
			_ = f.Close()
		}
	}

	type marketMeta struct {
		EndDate  time.Time
		Category string
		NegRisk  bool
	}
	var endDateMu sync.Mutex
	metaCache := make(map[string]marketMeta)
	lookupMarketMeta := func(conditionID string) (marketMeta, bool) {
		endDateMu.Lock()
		if m, ok := metaCache[conditionID]; ok {
			endDateMu.Unlock()
			return m, true
		}
		endDateMu.Unlock()
		qctx, qcancel := context.WithTimeout(ctx, 10*time.Second)
		defer qcancel()
		mkts, err := gc.GetByConditionIDs(qctx, []string{conditionID})
		if err != nil || len(mkts) == 0 {
			return marketMeta{}, false
		}
		ed, err := time.Parse(time.RFC3339, mkts[0].EndDate)
		if err != nil {
			ed, err = time.Parse("2006-01-02T15:04:05Z", mkts[0].EndDate)
			if err != nil {
				return marketMeta{}, false
			}
		}
		m := marketMeta{EndDate: ed, Category: mkts[0].Category, NegRisk: mkts[0].NegRisk}
		endDateMu.Lock()
		metaCache[conditionID] = m
		endDateMu.Unlock()
		return m, true
	}

	// Smart-money whale tracker: polls target wallet's CLOB trades and
	// pushes DM for large orders. Feature-flagged via -whale_enabled.
	// In whale signal_mode: BUY → SignalPrompt with buttons (boss clicks to follow);
	// SELL → auto-close matching positions.
	if whaleCfg.Enabled {
		whaleRepeatCooldown := 3 * time.Minute
		if raw := strings.TrimSpace(os.Getenv("WHALE_REPEAT_COOLDOWN")); raw != "" {
			d, err := time.ParseDuration(raw)
			if err != nil {
				slog.Warn("invalid WHALE_REPEAT_COOLDOWN", "value", raw, "err", err)
			} else {
				whaleRepeatCooldown = d
			}
		}
		whaleRepeatMinUSD := 0.0
		if raw := strings.TrimSpace(os.Getenv("WHALE_REPEAT_MIN_USD")); raw != "" {
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				slog.Warn("invalid WHALE_REPEAT_MIN_USD", "value", raw, "err", err)
			} else {
				whaleRepeatMinUSD = v
			}
		}
		whaleEventCooldown := 6 * time.Hour
		if raw := strings.TrimSpace(os.Getenv("WHALE_EVENT_COOLDOWN")); raw != "" {
			d, err := time.ParseDuration(raw)
			if err != nil {
				slog.Warn("invalid WHALE_EVENT_COOLDOWN", "value", raw, "err", err)
			} else {
				whaleEventCooldown = d
			}
		}
		whaleEventRepeatMinUSD := 50000.0
		if raw := strings.TrimSpace(os.Getenv("WHALE_EVENT_REPEAT_MIN_USD")); raw != "" {
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				slog.Warn("invalid WHALE_EVENT_REPEAT_MIN_USD", "value", raw, "err", err)
			} else {
				whaleEventRepeatMinUSD = v
			}
		}
		whaleConfirmLists := parseCSVSet(os.Getenv("WHALE_CONFIRM_LISTS"))
		if len(whaleConfirmLists) == 0 {
			whaleConfirmLists = parseCSVSet("watch,scout,target,tape,sports")
		}
		whaleConfirmWindow := 30 * time.Minute
		if raw := strings.TrimSpace(os.Getenv("WHALE_CONFIRM_WINDOW")); raw != "" {
			d, err := time.ParseDuration(raw)
			if err != nil {
				slog.Warn("invalid WHALE_CONFIRM_WINDOW", "value", raw, "err", err)
			} else {
				whaleConfirmWindow = d
			}
		}
		whaleConfirmMinWallets := 2
		if raw := strings.TrimSpace(os.Getenv("WHALE_CONFIRM_MIN_WALLETS")); raw != "" {
			v, err := strconv.Atoi(raw)
			if err != nil {
				slog.Warn("invalid WHALE_CONFIRM_MIN_WALLETS", "value", raw, "err", err)
			} else {
				whaleConfirmMinWallets = v
			}
		}
		whaleConfirmBypassUSD := 5000.0
		if raw := strings.TrimSpace(os.Getenv("WHALE_CONFIRM_BYPASS_USD")); raw != "" {
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				slog.Warn("invalid WHALE_CONFIRM_BYPASS_USD", "value", raw, "err", err)
			} else {
				whaleConfirmBypassUSD = v
			}
		}
		whaleConfirmMaxWorsePrice := 0.02
		if raw := strings.TrimSpace(os.Getenv("WHALE_CONFIRM_MAX_WORSE_PRICE")); raw != "" {
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				slog.Warn("invalid WHALE_CONFIRM_MAX_WORSE_PRICE", "value", raw, "err", err)
			} else {
				whaleConfirmMaxWorsePrice = v
			}
		}
		slog.Info("whale_consensus_config",
			"lists", sortedSetKeys(whaleConfirmLists),
			"window", whaleConfirmWindow.String(),
			"min_wallets", whaleConfirmMinWallets,
			"bypass_usd", whaleConfirmBypassUSD,
			"max_worse_price", whaleConfirmMaxWorsePrice,
		)
		whaleEdgeSnapshots := strings.TrimSpace(os.Getenv("WHALE_EDGE_SNAPSHOTS"))
		if whaleEdgeSnapshots == "" {
			whaleEdgeSnapshots = "db/strategy_iteration/whale_edge_snapshots.jsonl"
		}
		whaleEdgeTTL := 60 * time.Second
		if raw := strings.TrimSpace(os.Getenv("WHALE_EDGE_BLOCK_REFRESH")); raw != "" {
			d, err := time.ParseDuration(raw)
			if err != nil {
				slog.Warn("invalid WHALE_EDGE_BLOCK_REFRESH", "value", raw, "err", err)
			} else {
				whaleEdgeTTL = d
			}
		}
		whaleEdgeBlocker := newWhaleEdgeBlockCache(whaleEdgeSnapshots, whaleEdgeBlockConfig{
			Min15mSamples:  parseWhaleEnvInt("WHALE_EDGE_BLOCK_15M_SAMPLES", 2),
			Max15mAvgPP:    parseWhaleEnvFloat("WHALE_EDGE_BLOCK_15M_MAX_AVG_PP", -1),
			Min1hSamples:   parseWhaleEnvInt("WHALE_EDGE_BLOCK_1H_SAMPLES", 1),
			Max1hAvgPP:     parseWhaleEnvFloat("WHALE_EDGE_BLOCK_1H_MAX_AVG_PP", -5),
			HotMinUSD:      parseWhaleEnvFloat("WHALE_EDGE_HOT_MIN_USD", 750),
			HotMinSamples:  parseWhaleEnvInt("WHALE_EDGE_HOT_MIN_SAMPLES", 2),
			HotMinAvgPP:    parseWhaleEnvFloat("WHALE_EDGE_HOT_MIN_AVG_PP", 2),
			HotMinWinRate:  parseWhaleEnvFloat("WHALE_EDGE_HOT_MIN_WIN_RATE", 60),
			HotMin5mAvgPP:  parseWhaleEnvFloat("WHALE_EDGE_HOT_MIN_5M_AVG_PP", 0.5),
			HotMin15mAvgPP: parseWhaleEnvFloat("WHALE_EDGE_HOT_MIN_15M_AVG_PP", 0),
			HotMax1hNegPP:  parseWhaleEnvFloat("WHALE_EDGE_HOT_MAX_1H_NEG_PP", -5),
			Refresh:        whaleEdgeTTL,
		})
		if reason, ok := whaleEdgeBlocker.Reason(""); ok {
			slog.Warn("whale_edge_block_init", "reason", reason)
		} else {
			slog.Info("whale_edge_block_config",
				"snapshots", whaleEdgeSnapshots,
				"refresh", whaleEdgeTTL.String(),
				"block_15m_samples", whaleEdgeBlocker.cfg.Min15mSamples,
				"block_15m_max_avg_pp", whaleEdgeBlocker.cfg.Max15mAvgPP,
				"block_1h_samples", whaleEdgeBlocker.cfg.Min1hSamples,
				"block_1h_max_avg_pp", whaleEdgeBlocker.cfg.Max1hAvgPP,
				"hot_min_usd", whaleEdgeBlocker.cfg.HotMinUSD,
				"hot_min_samples", whaleEdgeBlocker.cfg.HotMinSamples,
				"hot_min_avg_pp", whaleEdgeBlocker.cfg.HotMinAvgPP,
				"hot_min_win_rate", whaleEdgeBlocker.cfg.HotMinWinRate,
				"hot_min_5m_avg_pp", whaleEdgeBlocker.cfg.HotMin5mAvgPP,
				"hot_min_15m_avg_pp", whaleEdgeBlocker.cfg.HotMin15mAvgPP,
				"hot_max_1h_neg_pp", whaleEdgeBlocker.cfg.HotMax1hNegPP,
			)
		}
		var whaleRepeatMu sync.Mutex
		whaleRepeatLast := map[string]time.Time{}
		shouldSuppressWhaleRepeat := func(ev whale.AlertEvent, side string) bool {
			if signalMode != "whale" || side != "BUY" || whaleRepeatCooldown <= 0 {
				return false
			}
			key := strings.ToLower(ev.Wallet) + "|" + ev.AssetID
			whaleRepeatMu.Lock()
			defer whaleRepeatMu.Unlock()
			if whaleRepeatBypassesCooldown(ev.Notional, whaleRepeatMinUSD) {
				whaleRepeatLast[key] = ev.Timestamp
				return false
			}
			if last, ok := whaleRepeatLast[key]; ok && ev.Timestamp.Sub(last) < whaleRepeatCooldown {
				return true
			}
			whaleRepeatLast[key] = ev.Timestamp
			return false
		}
		var whaleEventMu sync.Mutex
		whaleEventLast := map[string]time.Time{}
		var whaleSeenMu sync.Mutex
		whaleSeenTrades := loadWhaleSeenTradeIDs(whaleLogPath)
		markWhaleTradeSeen := func(ev whale.AlertEvent) bool {
			tradeID := strings.ToLower(strings.TrimSpace(ev.TradeID))
			if tradeID == "" {
				return true
			}
			key := strings.ToLower(strings.TrimSpace(ev.Wallet)) + "|" + tradeID
			whaleSeenMu.Lock()
			defer whaleSeenMu.Unlock()
			if _, ok := whaleSeenTrades[key]; ok {
				return false
			}
			whaleSeenTrades[key] = struct{}{}
			return true
		}
		{
			repeatSeed, eventSeed, repeatN, eventN := loadWhaleCooldownSeeds(whaleLogPath, whaleRepeatCooldown, whaleEventCooldown)
			for k, v := range repeatSeed {
				whaleRepeatLast[k] = v
			}
			for k, v := range eventSeed {
				whaleEventLast[k] = v
			}
			if repeatN > 0 || eventN > 0 {
				slog.Info("whale_cooldown_seeded", "repeat_keys", repeatN, "event_keys", eventN, "log", whaleLogPath)
			}
		}
		shouldSuppressWhaleEventRepeat := func(ev whale.AlertEvent, side string) (bool, string) {
			if signalMode != "whale" || side != "BUY" || whaleEventCooldown <= 0 {
				return false, ""
			}
			event := whaleEventKeyForAlert(ev)
			key := strings.ToLower(ev.Wallet) + "|" + event
			whaleEventMu.Lock()
			defer whaleEventMu.Unlock()
			if whaleNotionalAtLeast(ev.Notional, whaleEventRepeatMinUSD) {
				whaleEventLast[key] = ev.Timestamp
				return false, event
			}
			if last, ok := whaleEventLast[key]; ok && ev.Timestamp.Sub(last) < whaleEventCooldown {
				return true, event
			}
			whaleEventLast[key] = ev.Timestamp
			return false, event
		}
		whaleConfirmGate := newWhaleConfirmGate(whaleConfirmWindow, whaleConfirmMinWallets, whaleConfirmBypassUSD, whaleConfirmMaxWorsePrice)
		wt := whale.NewTracker(whaleCfg, func(ev whale.AlertEvent) {
			if !markWhaleTradeSeen(ev) {
				slog.Debug("whale_duplicate_trade_suppressed", "wallet", ev.Label, "trade_id", ev.TradeID, "market", ev.Question)
				return
			}
			side := strings.ToUpper(ev.Side)

			if signalMode == "copytrade" {
				baseAlert := notify.WhaleAlertEvent{
					Wallet:      ev.Wallet,
					Label:       ev.Label,
					Side:        ev.Side,
					SizeUnits:   ev.SizeUnits,
					Price:       ev.Price,
					Notional:    ev.Notional,
					Market:      ev.Question,
					Outcome:     ev.Outcome,
					TradeID:     ev.TradeID,
					LinkURL:     ev.LinkURL,
					ProfileURL:  ev.ProfileURL,
					Timestamp:   ev.Timestamp,
					TotalShares: ev.TotalShares,
					AvgPrice:    ev.AvgPrice,
					PctSold:     ev.PctSold,
				}

				footballScore := feed.IsFootballScoreMarketText(ev.Question + " " + ev.Slug)
				if ev.Price < copytradeEntryPriceFloor(footballScore, paperFollowFootballScore) || ev.Price > followMaxEntryPrice {
					appendWhaleTrade(ev, "skip", "price_filtered")
					slog.Info("copytrade_price_filtered", "wallet", ev.Label, "price", ev.Price, "market", ev.Question)
					return
				}

				switch side {
				case "BUY":
					if minTierFilter != "" {
						wTier := walletTiers[strings.ToLower(ev.Wallet)]
						allowed := false
						switch strings.ToUpper(minTierFilter) {
						case "A":
							allowed = wTier == "A"
						case "B":
							allowed = wTier == "A" || wTier == "B"
						case "C":
							allowed = wTier == "A" || wTier == "B" || wTier == "C"
						}
						if !allowed {
							appendWhaleTrade(ev, "skip", "tier_filtered:"+wTier)
							slog.Info("copytrade_tier_filtered", "wallet", ev.Label, "tier", wTier, "min_tier", minTierFilter, "market", ev.Question)
							tierAlert := baseAlert
							tierLabel := ev.Label
							if tierLabel == "" {
								tierLabel = ev.Wallet
							}
							tierAlert.Label = fmt.Sprintf("⚠️👀 [观察] %s (Tier %s)", tierLabel, wTier)
							notifier.WhaleAlert(tierAlert)
							return
						}
					}
					autoAllowed, followAction := copytradeAutoAllowed(ev.Wallet)
					if !autoAllowed {
						appendWhaleTrade(ev, "skip", "follow_action:"+followAction)
						slog.Info("copytrade_follow_action_filtered",
							"wallet", ev.Label,
							"action", followAction,
							"tier", walletTiers[strings.ToLower(ev.Wallet)],
							"market", ev.Question,
						)
						alert := baseAlert
						label := ev.Label
						if label == "" {
							label = ev.Wallet
						}
						alert.Label = fmt.Sprintf("👀 [%s] %s", followAction, label)
						notifier.WhaleAlert(alert)
						return
					}
					if ok, filterReason := copytradeMarketDecision(ev.Question, ev.Slug, ev.Outcome, paperFollowFootballScore); !ok {
						appendWhaleTrade(ev, "skip", filterReason)
						slog.Info("copytrade_market_filtered", "wallet", ev.Label, "reason", filterReason, "market", ev.Question, "slug", ev.Slug)
						return
					}
					// A-tier wallet passed smart-money gates — mark alert with strong recommendation.
					{
						aLabel := ev.Label
						if aLabel == "" {
							aLabel = ev.Wallet
						}
						meta := walletMetas[strings.ToLower(ev.Wallet)]
						baseAlert.Label = fmt.Sprintf("🔥💰 [自动小额] %s (Tier A · smart %.1f)", aLabel, meta.SmartMoneyScore)
					}
					// Extra DM via sidecar bot so A-tier alerts don't get buried
					aTierMsg := fmt.Sprintf("🔥💰 A级鲸鱼下单\n%s · %s\n💰 %.0f shares @ %.4f = $%.0f\n🐋 %s\n%s",
						ev.Question, ev.Outcome,
						ev.SizeUnits, ev.Price, ev.Notional,
						baseAlert.Label, ev.LinkURL)
					notifier.SidecarAlert(aTierMsg)
					isNegRisk := true
					if meta, ok := lookupMarketMeta(ev.ConditionID); ok {
						isNegRisk = meta.NegRisk
						if time.Until(meta.EndDate) > 30*24*time.Hour {
							appendWhaleTrade(ev, "skip", fmt.Sprintf("settlement_too_far:%s", meta.EndDate.Format("2006-01-02")))
							slog.Info("copytrade_settlement_filtered", "wallet", ev.Label, "market", ev.Question, "end_date", meta.EndDate.Format("2006-01-02"))
							notifier.WhaleAlert(baseAlert)
							return
						}
					}
					if err := rm.AllowOpen(time.Now()); err != nil {
						appendWhaleTrade(ev, "skip", "risk_blocked:"+err.Error())
						slog.Info("copytrade_blocked", "reason", err.Error(), "wallet", ev.Label, "market", ev.Question)
						return
					}
					tick := feed.Tick{Mid: ev.Price, Time: ev.Timestamp}
					sizeUSD := copytradeMarketSize(copytradeForWallet(ev.Wallet), footballScore, paperFollowFootballScore, paperFootballScoreSize)
					tier := walletTiers[strings.ToLower(ev.Wallet)]
					if tier == "" {
						tier = "?"
					}
					pos, err := pm.OpenSized(ev.AssetID, ev.ConditionID, tick, sizeUSD)
					if err != nil {
						appendWhaleTrade(ev, "skip", "open_rejected:"+err.Error())
						slog.Warn("copytrade_open_rejected", "wallet", ev.Label, "asset", short(ev.AssetID), "err", err.Error())
						return
					}
					pos.Question = ev.Question
					pos.Outcome = ev.Outcome
					pos.Source = "copytrade"
					if footballScore {
						pos.Source = "copytrade_football_score"
					}
					pos.WalletLabel = ev.Label
					crossPx := ev.Price * 1.05
					if crossPx > 0.99 {
						crossPx = 0.99
					}
					// Clear any stale orders before placing new one
					if v2, ok := orderClient.(*order.V2Client); ok {
						if err := v2.CancelAllOpen(ctx); err != nil {
							slog.Warn("copytrade_pre_cancel_err", "err", err)
						}
					}
					intent := order.Intent{
						AssetID: ev.AssetID,
						Market:  ev.ConditionID,
						Side:    order.Buy,
						SizeUSD: sizeUSD,
						LimitPx: crossPx,
						Type:    order.GTC,
						NegRisk: isNegRisk,
					}
					result, err := orderClient.Submit(ctx, intent)
					// Retry with reduced size on balance error
					if err != nil && strings.Contains(err.Error(), "not enough balance") {
						avail := parseAvailableBalance(err.Error())
						if avail > 1.0 {
							slog.Info("copytrade_balance_retry", "original_size", sizeUSD, "available", avail)
							intent.SizeUSD = avail * 0.95
							result, err = orderClient.Submit(ctx, intent)
							if err == nil && result.Status == order.StatusFilled {
								sizeUSD = intent.SizeUSD
							}
						}
					}
					if err != nil || result.Status != order.StatusFilled {
						reason := ""
						if err != nil {
							reason = err.Error()
						} else {
							reason = fmt.Sprintf("order %s: %s", result.Status, result.Error)
						}
						slog.Warn("copytrade_submit_err", "wallet", ev.Label, "err", reason, "status", result.Status)
						if !strings.Contains(reason, "not enough balance") {
							errMsg := fmt.Sprintf("❌ 跟单失败\n%s · %s\n💰 $%.0f @ %.4f · Tier %s\n🐋 %s\n⚠️ %s",
								ev.Question, ev.Outcome,
								sizeUSD, ev.Price, tier,
								ev.Label, reason)
							notifier.SidecarAlert(errMsg)
						}
						pm.Close(pos.ID, strategy.ExitSignal{Reason: "submit_failed"})
						savePositions()
						appendWhaleTrade(ev, "submit_failed", reason)
					} else {
						if err := pm.SetOpenFee(pos.ID, result.FeeUSD); err != nil {
							slog.Warn("copytrade_set_open_fee_fail", "pos", pos.ID, "err", err.Error())
						}
						slog.Info("copytrade_filled",
							"pos", pos.ID,
							"wallet", ev.Label,
							"tier", tier,
							"asset", short(ev.AssetID),
							"outcome", ev.Outcome,
							"price", ev.Price,
							"size_usd", sizeUSD,
							"order_id", result.OrderID,
							"fee_usd", result.FeeUSD,
						)
						fillMsg := fmt.Sprintf("✅ 跟单成功\n%s · %s\n💰 $%.0f @ %.4f · Tier %s\n🐋 %s 买入 $%.0f\n🆔 %s",
							ev.Question, ev.Outcome,
							sizeUSD, ev.Price, tier,
							ev.Label, ev.Notional,
							result.OrderID)
						notifier.SidecarAlert(fillMsg)
						signalSource := "copytrade_" + ev.Label
						if footballScore {
							signalSource = "copytrade_football_score_" + ev.Label
						}
						src.Mark(pos.ID, signalSource, result.OrderID)
						savePositions()
						buyTimesMap[ev.AssetID] = time.Now()
						saveBuyTimes()
						appendWhaleTrade(ev, "followed", "")
						select {
						case pnlTrigger <- struct{}{}:
						default:
						}
					}
					notifier.WhaleAlert(baseAlert)
					return

				case "SELL":
					var matches []strategy.Position
					for _, pos := range pm.Snapshot() {
						if pos.AssetID == ev.AssetID {
							matches = append(matches, pos)
						}
					}
					if len(matches) == 0 {
						appendWhaleTrade(ev, "sell_no_pos", "")
						notifier.WhaleAlert(baseAlert)
						return
					}
					sellPct := ev.PctSold / 100.0
					if sellPct <= 0 || sellPct > 1 {
						sellPct = 1.0
					}
					now := time.Now()
					closed := 0
					for _, pos := range matches {
						exitSig := strategy.ExitSignal{
							AssetID:  pos.AssetID,
							Market:   pos.Market,
							Time:     now,
							EntryMid: pos.EntryMid,
							ExitMid:  ev.Price,
							HeldFor:  now.Sub(pos.EntryTime),
							ChangePP: (ev.Price - pos.EntryMid) * 100,
							Reason:   strategy.ExitReason("whale_sell"),
						}
						closeUnits := pos.Units * sellPct
						var closedPos strategy.Position
						var err error
						if sellPct >= 0.95 {
							closedPos, err = pm.Close(pos.ID, exitSig)
						} else {
							closedPos, err = pm.PartialClose(pos.ID, closeUnits, exitSig)
						}
						if err != nil {
							slog.Warn("copytrade_close_miss", "pos", pos.ID, "err", err.Error())
							continue
						}
						entryFeeShare := 0.0
						if closedPos.InitUnits > 0 {
							entryFeeShare = closedPos.OpenFeeUSD * (closedPos.Units / closedPos.InitUnits)
						}
						netPnL := closedPos.PnLUSD - entryFeeShare
						source, openOID := src.Take(closedPos.ID)
						if tripped := rm.OnClose(netPnL, now); tripped {
							rst := rm.State()
							slog.Error("risk_trip",
								"reason", string(rst.BlockReason),
								"day_pnl_usd", rst.DayRealizedPnL,
								"cap_usd", rst.DayLossCapUSD,
							)
						}
						if err := rm.SaveState(riskStatePath); err != nil {
							slog.Warn("risk_save_err", "err", err)
						}
						closeReason := "copytrade_sell"
						if sellPct < 0.95 {
							closeReason = fmt.Sprintf("copytrade_partial_%.0f%%", ev.PctSold)
						}
						if jerr := jrn.Append(journal.TradeRecord{
							ID:           closedPos.ID,
							AssetID:      closedPos.AssetID,
							Market:       closedPos.Market,
							Question:     ev.Question,
							Outcome:      ev.Outcome,
							Side:         "buy",
							SizeUSD:      closedPos.SizeUSD,
							Units:        closedPos.Units,
							EntryMid:     closedPos.EntryMid,
							EntryTime:    closedPos.EntryTime,
							ExitMid:      closedPos.ExitMid,
							ExitTime:     closedPos.ExitTime,
							ExitReason:   string(closedPos.ExitReason),
							HeldSec:      int(closedPos.ExitTime.Sub(closedPos.EntryTime).Seconds()),
							PnLUSD:       closedPos.PnLUSD,
							EntryFeeUSD:  entryFeeShare,
							NetPnLUSD:    netPnL,
							OpenOrderID:  openOID,
							CloseOrderID: closeReason,
							Mode:         "paper",
							SignalSource: source,
						}); jerr != nil {
							slog.Warn("journal_append_fail", "asset", short(ev.AssetID), "err", jerr.Error())
						}
						closed++
						slog.Info("copytrade_sell_executed",
							"pos", pos.ID,
							"pct_sold", fmt.Sprintf("%.0f%%", ev.PctSold),
							"close_units", closeUnits,
							"remaining", pos.Units-closeUnits,
							"pnl", closedPos.PnLUSD,
						)
						pnlIcon := "🟢"
						if closedPos.PnLUSD < 0 {
							pnlIcon = "🔴"
						}
						sellType := "全部平仓"
						if sellPct < 0.95 {
							sellType = fmt.Sprintf("部分平仓 %.0f%%", ev.PctSold)
						}
						sellMsg := fmt.Sprintf("%s 跟卖 · %s\n%s · %s\n💰 PnL $%+.2f · %.0f shares\n🐋 %s 卖出 %.0f%%",
							pnlIcon, sellType,
							ev.Question, ev.Outcome,
							closedPos.PnLUSD, closeUnits,
							ev.Label, ev.PctSold)
						notifier.SidecarAlert(sellMsg)
					}
					if closed > 0 {
						savePositions()
						action := "closed"
						if sellPct < 0.95 {
							action = fmt.Sprintf("partial_%.0f%%", ev.PctSold)
						}
						appendWhaleTrade(ev, action, fmt.Sprintf("matched=%d", closed))
						slog.Info("copytrade_closed", "wallet", ev.Label, "asset", short(ev.AssetID), "count", closed, "pct", ev.PctSold)
					} else {
						appendWhaleTrade(ev, "sell_close_fail", "")
					}
					notifier.WhaleAlert(baseAlert)
					return

				default:
					appendWhaleTrade(ev, "other", side)
					notifier.WhaleAlert(baseAlert)
					return
				}
			}

			meta := walletFileMetas[strings.ToLower(ev.Wallet)]
			if ok, filterReason := whaleMarketDecision(ev.Question, ev.Slug, ev.Outcome, meta.List); !ok {
				appendWhaleTrade(ev, "skip", filterReason)
				slog.Info("whale_market_filtered", "wallet", ev.Label, "reason", filterReason, "market", ev.Question, "slug", ev.Slug)
				return
			}
			isFootballScore := feed.IsFootballScoreMarketText(ev.Question + " " + ev.Slug)
			if ok, filterReason := whalePriceDecision(side, ev.Price, isFootballScore); !ok {
				appendWhaleTrade(ev, "skip", filterReason)
				slog.Info("whale_price_filtered", "wallet", ev.Label, "side", side, "price", ev.Price, "market", ev.Question)
				return
			}
			if side == "BUY" && signalMode == "whale" {
				if reason, blocked := whaleEdgeBlocker.Reason(ev.Wallet); blocked {
					appendWhaleTrade(ev, "skip", "edge_blocked:"+reason)
					slog.Info("whale_edge_blocked",
						"wallet", ev.Label,
						"reason", reason,
						"market", ev.Question,
						"notional", ev.Notional,
					)
					return
				}
			}

			if shouldSuppressWhaleRepeat(ev, side) {
				appendWhaleTrade(ev, "cooldown", fmt.Sprintf("repeat_within:%s", whaleRepeatCooldown))
				slog.Info("whale_repeat_suppressed",
					"wallet", ev.Label,
					"asset", short(ev.AssetID),
					"notional", ev.Notional,
					"cooldown", whaleRepeatCooldown.String(),
					"repeat_min_usd", whaleRepeatMinUSD,
				)
				return
			}

			if suppressed, event := shouldSuppressWhaleEventRepeat(ev, side); suppressed {
				appendWhaleTrade(ev, "event_cooldown", fmt.Sprintf("event_repeat_within:%s event=%s", whaleEventCooldown, event))
				slog.Info("whale_event_repeat_suppressed",
					"wallet", ev.Label,
					"event", event,
					"asset", short(ev.AssetID),
					"notional", ev.Notional,
					"cooldown", whaleEventCooldown.String(),
					"repeat_min_usd", whaleEventRepeatMinUSD,
					"market", ev.Question,
				)
				return
			}

			confirmNote := ""
			if side == "BUY" && signalMode == "whale" {
				confirmNote = whaleDirectGateNote(meta.List)
				if _, needsConfirmation := whaleConfirmLists[strings.ToLower(meta.List)]; needsConfirmation {
					if hotReason, hot := whaleEdgeBlocker.HotReason(ev.Wallet); hot && whaleNotionalAtLeast(ev.Notional, whaleEdgeBlocker.cfg.HotMinUSD) {
						confirmNote = "gate edge-hot " + hotReason
						slog.Info("whale_consensus_edge_hot_bypass",
							"wallet", ev.Label,
							"list", meta.List,
							"reason", hotReason,
							"notional", ev.Notional,
							"min_usd", whaleEdgeBlocker.cfg.HotMinUSD,
							"market", ev.Question,
						)
					} else {
						decision := whaleConfirmGate.Observe(ev, meta.List)
						if !decision.Ready {
							appendWhaleTrade(ev, "pending_consensus", decision.Reason)
							slog.Info("whale_consensus_pending",
								"wallet", ev.Label,
								"list", meta.List,
								"event", decision.Event,
								"outcome", ev.Outcome,
								"wallets", decision.Wallets,
								"need", decision.Need,
								"notional", ev.Notional,
								"price", ev.Price,
							)
							return
						}
						confirmNote = whaleConfirmDecisionNote(decision)
						if decision.Reason != "" {
							slog.Info("whale_consensus_ready",
								"wallet", ev.Label,
								"list", meta.List,
								"event", decision.Event,
								"outcome", ev.Outcome,
								"wallets", decision.Wallets,
								"need", decision.Need,
								"reason", decision.Reason,
								"notional", ev.Notional,
								"price", ev.Price,
							)
						}
					}
				}
			}

			appendWhaleTrade(ev, "alert", signalMode)

			if signalMode != "whale" {
				notifier.WhaleAlert(notify.WhaleAlertEvent{
					Wallet:      ev.Wallet,
					Label:       ev.Label,
					Side:        ev.Side,
					SizeUnits:   ev.SizeUnits,
					Price:       ev.Price,
					Notional:    ev.Notional,
					Market:      ev.Question,
					Outcome:     ev.Outcome,
					TradeID:     ev.TradeID,
					LinkURL:     ev.LinkURL,
					ProfileURL:  ev.ProfileURL,
					Timestamp:   ev.Timestamp,
					TotalShares: ev.TotalShares,
					AvgPrice:    ev.AvgPrice,
					PctSold:     ev.PctSold,
				})
				return
			}

			// ---- whale-follow mode ----

			if side == "BUY" {
				// Send SignalPrompt with buy buttons so the boss can follow.
				choices := []notify.Choice{{
					AssetID:  ev.AssetID,
					Outcome:  ev.Outcome,
					Mid:      ev.Price,
					IsSignal: true,
				}}
				sigChoices := []notify.SignalChoice{{
					Slot: 0, Outcome: ev.Outcome, Mid: ev.Price, IsSignal: true,
				}}

				p := pending.Put(notify.PendingIntent{
					Market:   ev.ConditionID,
					Question: ev.Question,
					Choices:  choices,
				}, time.Now())

				whaleTag := ev.Label
				if whaleTag == "" {
					whaleTag = "鲸鱼"
				}
				if meta := walletFileMetas[strings.ToLower(ev.Wallet)]; meta.List != "" {
					tier := meta.Tier
					if tier == "" {
						tier = "?"
					}
					whaleTag = fmt.Sprintf("[%s · Tier %s] %s", meta.List, tier, whaleTag)
				}
				ctxLine := formatWhalePromptContext(ev, walletFileMetas[strings.ToLower(ev.Wallet)], confirmNote)

				nonceSnap := p.Nonce
				notifier.SignalPrompt(notify.SignalPromptEvent{
					Nonce:      p.Nonce,
					Match:      ev.Question,
					Context:    ctxLine,
					Slug:       ev.Slug,
					WhaleLabel: whaleTag,
					Choices:    sigChoices,
					ExpiresIn:  10 * time.Minute,
					OnSent: func(msgID int64, err error) {
						if err != nil {
							slog.Warn("notify_send_err", "err", err)
							return
						}
						if msgID == 0 {
							return
						}
						pending.SetMessageID(nonceSnap, msgID)
					},
				})
				slog.Info("whale_follow_prompt_sent",
					"asset", ev.AssetID,
					"outcome", ev.Outcome,
					"price", ev.Price,
					"notional", ev.Notional,
					"nonce", p.Nonce,
				)
				return
			}

			if side == "SELL" {
				// Collect matching open positions for this asset.
				var matchingPos []notify.ClosePosition
				for _, pos := range pm.Snapshot() {
					if pos.AssetID != ev.AssetID {
						continue
					}
					matchingPos = append(matchingPos, notify.ClosePosition{
						PosID:    pos.ID,
						SizeUSD:  pos.SizeUSD,
						Units:    pos.Units,
						EntryMid: pos.EntryMid,
					})
				}

				if len(matchingPos) == 0 {
					notifier.WhaleAlert(notify.WhaleAlertEvent{
						Wallet:      ev.Wallet,
						Label:       ev.Label,
						Side:        ev.Side,
						SizeUnits:   ev.SizeUnits,
						Price:       ev.Price,
						Notional:    ev.Notional,
						Market:      ev.Question,
						Outcome:     ev.Outcome,
						TradeID:     ev.TradeID,
						LinkURL:     ev.LinkURL,
						ProfileURL:  ev.ProfileURL,
						Timestamp:   ev.Timestamp,
						TotalShares: ev.TotalShares,
						AvgPrice:    ev.AvgPrice,
						PctSold:     ev.PctSold,
					})
					return
				}

				// We hold matching positions — push a close prompt with buttons.
				ci := closePending.Put(notify.CloseIntent{
					AssetID:    ev.AssetID,
					Market:     ev.ConditionID,
					Question:   ev.Question,
					Outcome:    ev.Outcome,
					WhalePrice: ev.Price,
				}, time.Now())

				nonceSnap := ci.Nonce
				notifier.ClosePrompt(notify.ClosePromptEvent{
					Nonce:            ci.Nonce,
					Market:           ev.Question,
					Outcome:          ev.Outcome,
					AssetID:          ev.AssetID,
					WhaleLabel:       ev.Label,
					WhaleSize:        ev.SizeUnits,
					WhaleNotl:        ev.Notional,
					WhalePrice:       ev.Price,
					LinkURL:          ev.LinkURL,
					ProfileURL:       ev.ProfileURL,
					Positions:        matchingPos,
					WhaleTotalShares: ev.TotalShares,
					WhaleAvgPrice:    ev.AvgPrice,
					WhalePctSold:     ev.PctSold,
					OnSent: func(msgID int64, err error) {
						if err != nil {
							slog.Warn("notify_send_err", "err", err)
							return
						}
						if msgID == 0 {
							return
						}
						closePending.SetMessageID(nonceSnap, msgID)
					},
				})
				slog.Info("whale_close_prompt_sent",
					"asset", ev.AssetID,
					"outcome", ev.Outcome,
					"positions", len(matchingPos),
					"nonce", ci.Nonce,
				)
				return
			}

			// Unrecognized side (e.g. MINT/REDEEM) — just notify.
			notifier.WhaleAlert(notify.WhaleAlertEvent{
				Wallet:      ev.Wallet,
				Label:       ev.Label,
				Side:        ev.Side,
				SizeUnits:   ev.SizeUnits,
				Price:       ev.Price,
				Notional:    ev.Notional,
				Market:      ev.Question,
				Outcome:     ev.Outcome,
				TradeID:     ev.TradeID,
				LinkURL:     ev.LinkURL,
				ProfileURL:  ev.ProfileURL,
				Timestamp:   ev.Timestamp,
				TotalShares: ev.TotalShares,
				AvgPrice:    ev.AvgPrice,
				PctSold:     ev.PctSold,
			})
		})
		go func() {
			if err := wt.Run(ctx); err != nil && ctx.Err() == nil {
				slog.Warn("whale_tracker_exit", "err", err.Error())
			}
		}()
	}

	// BTC prediction strategy: independent goroutine scanning PM BTC markets
	// vs Black-Scholes first-passage fair value. Fires SignalPrompt DMs for
	// gaps > btc_min_gap_pp. Completely independent of sports strategies.
	if btcCfg.Enabled {
		go func() {
			if err := btc.RunStrategy(ctx, btcCfg, func(sig btc.Signal) {
				dirEmoji := "📈"
				dirLabel := "BUY Yes"
				if sig.Direction == "BUY_NO" {
					dirEmoji = "📉"
					dirLabel = "BUY No"
				}

				// Build choice for the signal direction.
				choices := []notify.Choice{{
					AssetID:  sig.MarketID,
					Outcome:  sig.Direction,
					Mid:      sig.PMPrice,
					IsSignal: true,
				}}
				sigChoices := []notify.SignalChoice{{
					Slot: 0, Outcome: sig.Direction, Mid: sig.PMPrice, IsSignal: true,
				}}

				p := pending.Put(notify.PendingIntent{
					Market:   sig.MarketID,
					Question: sig.Question,
					Choices:  choices,
				}, time.Now())

				ctxLine := fmt.Sprintf(
					"%s BTC $%.0f · %s\nSpot: $%.0f · Vol: %.0f%% · Gap: %+.1fpp\nBS fair: %.1f%% vs PM: %.1f%%",
					dirEmoji, sig.Strike, dirLabel,
					sig.Spot, sig.Sigma*100, sig.GapPP,
					sig.BSProb*100, sig.PMPrice*100,
				)

				slug := "what-price-will-bitcoin-hit-before-2027"
				nonceSnap := p.Nonce
				notifier.SignalPrompt(notify.SignalPromptEvent{
					Nonce:      p.Nonce,
					Match:      sig.Question,
					Context:    ctxLine,
					Slug:       slug,
					WhaleLabel: "₿ BTC策略",
					Choices:    sigChoices,
					ExpiresIn:  2 * time.Hour,
					OnSent: func(msgID int64, err error) {
						if err != nil {
							slog.Warn("notify_send_err", "err", err)
							return
						}
						if msgID == 0 {
							return
						}
						pending.SetMessageID(nonceSnap, msgID)
					},
				})
				slog.Info("btc_strategy.signal_pushed",
					"strike", sig.Strike,
					"direction", sig.Direction,
					"gap_pp", sig.GapPP,
					"nonce", p.Nonce,
				)
			}); err != nil && ctx.Err() == nil {
				slog.Warn("btc_strategy_exit", "err", err.Error())
			}
		}()
	}

	if updownCfg.Enabled {
		go func() {
			if err := btc.RunUpDownStrategy(ctx, updownCfg, func(sig btc.UpDownSignal) {
				slog.Info("updown_strategy.auto_bet",
					"slug", sig.MarketSlug,
					"direction", sig.PredictedDirection,
					"confidence", fmt.Sprintf("%.3f", sig.Confidence),
					"pm_price", fmt.Sprintf("%.3f", sig.PMPrice),
					"size", fmt.Sprintf("%.2f", sig.SizeUSD),
					"spot", fmt.Sprintf("%.0f", sig.Spot),
				)
			}); err != nil && ctx.Err() == nil {
				slog.Warn("updown_strategy_exit", "err", err.Error())
			}
		}()
	}

	// Phase 10: BTC Daily Threshold scanner
	if p10.BTCDailyEnabled {
		go func() {
			slog.Info("btc_daily_scanner.start", "interval", p10.BTCDailyInterval.String(), "min_edge_pp", p10.BTCDailyMinEdge)
			if sigs, err := btc.ScanDailyBTC(ctx, p10.BTCDailyMinEdge); err != nil {
				slog.Warn("btc_daily_scan_err", "err", err.Error())
			} else {
				for _, s := range sigs {
					msg := btc.FormatDailySignal(s)
					slog.Info("btc_daily_signal", "slug", s.Market.Slug, "side", s.Side, "edge_pp", fmt.Sprintf("%.1f", s.Edge*100))
					notifier.TextAlert(msg)
				}
			}
			tk := time.NewTicker(p10.BTCDailyInterval)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					sigs, err := btc.ScanDailyBTC(ctx, p10.BTCDailyMinEdge)
					if err != nil {
						slog.Warn("btc_daily_scan_err", "err", err.Error())
						continue
					}
					slog.Info("btc_daily_scan", "signals", len(sigs))
					for _, s := range sigs {
						msg := btc.FormatDailySignal(s)
						slog.Info("btc_daily_signal", "slug", s.Market.Slug, "side", s.Side, "edge_pp", fmt.Sprintf("%.1f", s.Edge*100))
						notifier.TextAlert(msg)
					}
				}
			}
		}()
	}

	// Phase 10: Elon tweet count scanner
	if p10.ElonEnabled {
		go func() {
			xToken := os.Getenv("X_BEARER_TOKEN")
			if xToken == "" {
				slog.Warn("elon_scanner.disabled", "reason", "X_BEARER_TOKEN not set")
				return
			}
			slog.Info("elon_scanner.start", "interval", p10.ElonInterval.String(), "min_edge_pp", p10.ElonMinEdge)
			tk := time.NewTicker(p10.ElonInterval)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					markets, err := elon.FetchElonMarkets(ctx)
					if err != nil {
						slog.Warn("elon_fetch_err", "err", err.Error())
						continue
					}
					if len(markets) == 0 {
						slog.Info("elon_scan", "markets", 0)
						continue
					}
					m0 := markets[0]
					count, err := elon.CountTweetsViaSearch(ctx, xToken, m0.Start, time.Now().UTC())
					if err != nil {
						slog.Warn("elon_count_err", "err", err.Error())
						continue
					}
					hoursLeft := m0.End.Sub(time.Now().UTC()).Hours()
					sigs := elon.EvalSignals(markets, count, hoursLeft, p10.ElonMinEdge)
					slog.Info("elon_scan", "markets", len(markets), "count", count, "hours_left", fmt.Sprintf("%.1f", hoursLeft), "signals", len(sigs))
					for _, s := range sigs {
						msg := elon.FormatTweetSignal(s)
						slog.Info("elon_signal", "range", fmt.Sprintf("%d-%d", s.Market.RangeLo, s.Market.RangeHi), "side", s.Side, "edge_pp", fmt.Sprintf("%.1f", s.Edge*100))
						notifier.TextAlert(msg)
					}
				}
			}
		}()
	}

	// Phase 10: Eurovision odds scanner
	if p10.EurovisionEnabled {
		slog.Warn("eurovision_scanner.disabled", "reason", "third-party bookmaker odds API removed")
	}

	// Feed-silence watchdog + periodic risk snapshot. SPEC §6: >30s WSS
	// silence trips breaker. We also push a risk summary every 60s so the
	// heartbeat log has a recent snapshot.
	go func() {
		tk := time.NewTicker(5 * time.Second)
		defer tk.Stop()
		lastSummary := time.Now()
		prevBlocked := false
		prevReason := ""
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				if at := ws.LastEventAt(); !at.IsZero() {
					rm.OnFeedHeartbeat(at)
				}
				silent, tripped := rm.CheckFeed(now)
				st := rm.State()
				if tripped {
					slog.Error("risk_trip",
						"reason", string(risk.BlockFeedSilence),
						"silent_sec", int(silent.Seconds()),
						"connected", ws.Connected(),
					)
					notifier.RiskTrip(notify.RiskTripEvent{
						Reason:        string(risk.BlockFeedSilence),
						DayPnLUSD:     st.DayRealizedPnL,
						DayLossCapUSD: st.DayLossCapUSD,
						SilentSec:     int(silent.Seconds()),
						OpenPositions: pm.Stats().Open,
					})
				}
				// Auto-resume when the breaker tripped ONLY because of feed
				// silence and the WSS reconnected. Socket-back is the real
				// "feed is healthy again" signal — waiting for trade chatter
				// starves us during quiet markets. Daily-loss + manual-pause
				// still require an explicit human resume.
				if st.Blocked && st.BlockReason == risk.BlockFeedSilence && ws.Connected() {
					rm.Resume()
					slog.Info("risk_auto_resume",
						"prev_reason", string(risk.BlockFeedSilence),
						"silent_sec", int(silent.Seconds()),
					)
					st = rm.State()
				}
				// Detect resume transition (auto or manual).
				if prevBlocked && !st.Blocked {
					notifier.RiskResume(notify.RiskResumeEvent{
						PrevReason:    prevReason,
						DayPnLUSD:     st.DayRealizedPnL,
						DayLossCapUSD: st.DayLossCapUSD,
					})
				}
				prevBlocked = st.Blocked
				if st.Blocked {
					prevReason = string(st.BlockReason)
				}
				if now.Sub(lastSummary) >= 60*time.Second {
					lastSummary = now
					st := rm.State()
					slog.Info("risk_status",
						"day", st.Day,
						"day_pnl_usd", st.DayRealizedPnL,
						"cap_usd", st.DayLossCapUSD,
						"cum_pnl", st.CumulativePnL,
						"peak_equity", st.PeakEquity,
						"drawdown", st.DrawdownUSD,
						"dd_cap", st.DrawdownCapUSD,
						"blocked", st.Blocked,
						"reason", string(st.BlockReason),
						"single_loss_flags", st.SingleLossFlags,
						"feed_silent_sec", int(silent.Seconds()),
					)
				}
			}
		}
	}()

	// Settlement watcher (exit_mode=hold). Polls gamma for each open position's
	// market every 60s; when a market is `closed` we close the position using
	// OutcomePrices[SlotIdx] as the final fill — 1.0 for the winning side,
	// 0.0 for the loser. Does the same risk/journal/notify bookkeeping the
	// auto-exit tracker does. SPEC §2 exit_mode=hold.
	if wantSettlement {
		go func() {
			tk := time.NewTicker(60 * time.Second)
			defer tk.Stop()
			lastHeldLog := time.Time{}
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
				}
				open := pm.Snapshot()
				if len(open) == 0 {
					continue
				}
				// Collect unique conditionIDs from the open positions. Copytrade
				// can follow markets that were not part of the startup market
				// scan, so position.Market is the durable source of truth here.
				seen := make(map[string]struct{}, len(open))
				ids := make([]string, 0, len(open))
				for _, p := range open {
					conditionID := p.Market
					if conditionID == "" {
						if me := meta[p.AssetID]; me != nil {
							conditionID = me.ConditionID
						}
					}
					if conditionID == "" {
						continue
					}
					key := strings.ToLower(conditionID)
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}
					ids = append(ids, conditionID)
				}
				if len(ids) == 0 {
					continue
				}
				qctx, qcancel := context.WithTimeout(ctx, 15*time.Second)
				mkts2, err := gc.GetByConditionIDs(qctx, ids)
				qcancel()
				if err != nil {
					slog.Warn("settlement_poll_fail", "err", err.Error(), "ids", len(ids))
					continue
				}
				byCond := make(map[string]feed.Market, len(mkts2))
				for _, m := range mkts2 {
					byCond[strings.ToLower(m.ConditionID)] = m
				}
				// Periodic "still holding" log (once per 5 min) — easy to grep for.
				now := time.Now()
				if now.Sub(lastHeldLog) >= 5*time.Minute {
					lastHeldLog = now
					slog.Info("hold_status",
						"open", len(open),
						"markets_polled", len(ids),
						"markets_returned", len(mkts2),
						"resolved_seen", countResolved(mkts2),
					)
				}
				for _, p := range open {
					conditionID := p.Market
					if conditionID == "" {
						if me := meta[p.AssetID]; me != nil {
							conditionID = me.ConditionID
						}
					}
					if conditionID == "" {
						continue
					}
					m, ok := byCond[strings.ToLower(conditionID)]
					if !ok {
						if !liveTrading && exitMode == "ladder" && ladderCfg.MaxHold > 0 && now.Sub(p.EntryTime) >= ladderCfg.MaxHold {
							sig := strategy.ExitSignal{
								AssetID:  p.AssetID,
								Market:   p.Market,
								Time:     now,
								EntryMid: p.EntryMid,
								PeakMid:  p.EntryMid,
								ExitMid:  p.EntryMid,
								HeldFor:  now.Sub(p.EntryTime),
								Reason:   strategy.ExitTimeout,
							}
							closed, cerr := pm.Close(p.ID, sig)
							if cerr != nil {
								slog.Warn("paper_timeout_flat_close_miss", "pos", p.ID, "asset", short(p.AssetID), "err", cerr.Error())
								continue
							}
							savePositions()
							ladder.Forget(p.ID)
							if recorder != nil {
								if rerr := recorder.Stop(closed.ID); rerr != nil {
									slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
								}
							}
							source, openOID := src.Take(closed.ID)
							if jerr := jrn.Append(journal.TradeRecord{
								ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
								Question:     closed.Question,
								Outcome:      closed.Outcome,
								Side:         "buy",
								SizeUSD:      closed.SizeUSD,
								Units:        closed.Units,
								EntryMid:     closed.EntryMid,
								EntryTime:    closed.EntryTime,
								ExitMid:      closed.ExitMid,
								ExitTime:     closed.ExitTime,
								ExitReason:   string(closed.ExitReason),
								HeldSec:      int(sig.HeldFor.Seconds()),
								PnLUSD:       closed.PnLUSD,
								NetPnLUSD:    closed.PnLUSD,
								Tranche:      "timeout_flat",
								OpenOrderID:  openOID,
								CloseOrderID: fmt.Sprintf("timeout-flat-%s", short(p.AssetID)),
								Mode:         "paper",
								SignalSource: source,
							}); jerr != nil {
								slog.Warn("journal_append_fail", "asset", short(p.AssetID), "err", jerr.Error())
							}
							stats := pm.Stats()
							slog.Info("paper_timeout_flat_exit",
								"pos", closed.ID,
								"asset", short(p.AssetID),
								"q", closed.Question,
								"outcome", closed.Outcome,
								"entry", p.EntryMid,
								"held_sec", int(sig.HeldFor.Seconds()),
								"open_positions", stats.Open,
								"realized_pnl", stats.RealizedPnLUSD,
							)
						}
						continue
					}
					slotIdx, slotOK := settlementSlotForPosition(p, m, meta[p.AssetID])
					prices := m.OutcomePrices()
					if !slotOK || slotIdx < 0 || slotIdx >= len(prices) {
						slog.Warn("settlement_slot_miss", "pos", p.ID, "asset", short(p.AssetID), "market", short(conditionID), "prices", len(prices))
						continue
					}
					settleMid, perr := strconv.ParseFloat(prices[slotIdx], 64)
					if perr != nil {
						slog.Warn("settlement_price_parse_fail",
							"asset", short(p.AssetID),
							"raw", prices[slotIdx],
							"err", perr.Error())
						continue
					}
					reason := strategy.ExitSettlement
					if !m.Closed {
						if liveTrading || exitMode != "ladder" || ladderCfg.MaxHold <= 0 || now.Sub(p.EntryTime) < ladderCfg.MaxHold {
							continue
						}
						reason = strategy.ExitTimeout
					}
					sig := strategy.ExitSignal{
						AssetID:  p.AssetID,
						Market:   p.Market,
						Time:     now,
						EntryMid: p.EntryMid,
						PeakMid:  p.EntryMid,
						ExitMid:  settleMid,
						HeldFor:  now.Sub(p.EntryTime),
						ChangePP: (settleMid - p.EntryMid) * 100,
						Reason:   reason,
					}
					closed, cerr := pm.Close(p.ID, sig)
					if cerr != nil {
						slog.Warn("settlement_close_miss", "pos", p.ID, "asset", short(p.AssetID), "err", cerr.Error())
						continue
					}
					savePositions()
					// Drop any ladder state that was still tracking this
					// position — settlement supersedes TP/SL/timeout.
					ladder.Forget(p.ID)
					if recorder != nil {
						if rerr := recorder.Stop(closed.ID); rerr != nil {
							slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
						}
					}
					orderID := fmt.Sprintf("settle-%s", short(p.AssetID))
					// Apportion remaining open fee to the portion of units
					// still open at settlement (p.InitUnits may be > p.Units
					// if ladder TP1 already fired); settlement has no exit fee.
					entryFeeShare := 0.0
					if p.InitUnits > 0 {
						entryFeeShare = p.OpenFeeUSD * (p.Units / p.InitUnits)
					}
					netPnL := closed.PnLUSD - entryFeeShare
					stats := pm.Stats()
					settleSource, _ := src.Peek(closed.ID)
					if settleSource != "manual" {
						if tripped := rm.OnClose(netPnL, now); tripped {
							rst := rm.State()
							slog.Error("risk_trip",
								"reason", string(rst.BlockReason),
								"day_pnl_usd", rst.DayRealizedPnL,
								"cap_usd", rst.DayLossCapUSD,
							)
							notifier.RiskTrip(notify.RiskTripEvent{
								Reason:        string(rst.BlockReason),
								DayPnLUSD:     rst.DayRealizedPnL,
								DayLossCapUSD: rst.DayLossCapUSD,
								OpenPositions: stats.Open,
							})
						}
						if err := rm.SaveState(riskStatePath); err != nil {
							slog.Warn("risk_save_err", "err", err)
						}
					}
					if netPnL <= -largeFillUSD || netPnL >= largeFillUSD {
						notifier.LargeFill(notify.LargeFillEvent{
							Question: paperPositionQuestion(p, m, meta[p.AssetID]),
							AssetID:  p.AssetID,
							Side:     "sell",
							SizeUSD:  p.SizeUSD,
							PnLUSD:   netPnL,
							EntryPx:  p.EntryMid,
							ExitPx:   settleMid,
							Reason:   string(reason),
							HeldSec:  int(sig.HeldFor.Seconds()),
						})
					}
					source, openOID := src.Take(closed.ID)
					if jerr := jrn.Append(journal.TradeRecord{
						ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
						Question:     paperPositionQuestion(closed, m, meta[closed.AssetID]),
						Outcome:      paperPositionOutcome(closed, m, slotIdx, meta[closed.AssetID]),
						Side:         "buy",
						SizeUSD:      closed.SizeUSD,
						Units:        closed.Units,
						EntryMid:     closed.EntryMid,
						EntryTime:    closed.EntryTime,
						ExitMid:      closed.ExitMid,
						ExitTime:     closed.ExitTime,
						ExitReason:   string(closed.ExitReason),
						HeldSec:      int(sig.HeldFor.Seconds()),
						PnLUSD:       closed.PnLUSD,
						EntryFeeUSD:  entryFeeShare,
						ExitFeeUSD:   0,
						NetPnLUSD:    netPnL,
						Tranche:      "settle",
						OpenOrderID:  openOID,
						CloseOrderID: orderID,
						Mode:         "paper",
						SignalSource: source,
					}); jerr != nil {
						slog.Warn("journal_append_fail", "asset", short(p.AssetID), "err", jerr.Error())
					}
					slog.Info("settlement_exit",
						"asset", short(p.AssetID),
						"q", paperPositionQuestion(p, m, meta[p.AssetID]),
						"outcome", paperPositionOutcome(p, m, slotIdx, meta[p.AssetID]),
						"entry", p.EntryMid,
						"settle", settleMid,
						"gross_pnl_usd", closed.PnLUSD,
						"entry_fee_usd", entryFeeShare,
						"net_pnl_usd", netPnL,
						"held_sec", int(sig.HeldFor.Seconds()),
						"reason", string(reason),
						"open_positions", stats.Open,
						"realized_pnl", stats.RealizedPnLUSD,
					)
				}
			}
		}()
	}

	go func() {
		tk := time.NewTicker(30 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				top := topWindow(sampler.Snapshot(), 5)
				for _, w := range top {
					slog.Info("top_window",
						"asset", short(w.AssetID),
						"q", metaQ(meta, w.AssetID),
						"samples", w.Samples,
						"mid", w.EndMid,
						"delta_pp", w.DeltaPP,
						"up", w.Upticks,
						"buy_ratio", w.BuyRatio,
					)
				}
			}
		}
	}()

	// Hourly P&L summary push — every 1h via sidecar bot (@Murphyoderbot).
	// Uses Polymarket data-api /positions for real on-chain positions.
	go func() {
		tk := time.NewTicker(1 * time.Hour)
		defer tk.Stop()
		pushPnL := func() {
			sgt := time.FixedZone("SGT", 8*3600)
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("📊 P&L · %s SGT\n\n", time.Now().In(sgt).Format("15:04")))

			if walletAddress == "" {
				return
			}

			positions, err := fetchDataAPIPositions(walletAddress)
			if err != nil {
				sb.WriteString(fmt.Sprintf("⚠️ 拉取仓位失败: %s\n", err))
				notifier.SidecarAlert(sb.String())
				slog.Warn("pnl_push_fetch_err", "err", err)
				return
			}

			// Auto-close settled/zeroed positions in local position manager to free MaxOpen slots
			settledAssets := map[string]float64{}
			for _, p := range positions {
				if _, tracked := buyTimesMap[p.Asset]; !tracked {
					continue
				}
				if p.CurPrice < 0.001 || p.CurPrice >= 0.99 || p.Size < 0.01 {
					settledAssets[p.Asset] = p.CurPrice
				}
			}
			if len(settledAssets) > 0 {
				purged := pm.PurgeSettled(settledAssets)
				if purged > 0 {
					slog.Info("positions_purged_settled", "count", purged)
					savePositions()
				}
			}

			// --- Query wallet balance ---
			var walletPUSD float64
			if onChain != nil {
				if bal, berr := onChain.PUSDBalance(ctx); berr == nil {
					f, _ := new(big.Float).SetInt(bal).Float64()
					walletPUSD = f / 1e6
				} else {
					slog.Warn("pnl_wallet_balance_err", "err", berr)
				}
			}

			// --- Classify positions ---
			type posLine struct {
				emoji, title, outcome string
				avgPrice, curPrice    float64
				size                  float64
				cost, value           float64
				pnl, pct              float64
				endDate               string
				redeemable            bool
				buyTime               time.Time
			}
			var activeLines []posLine
			var activeCost, activeValue float64
			for _, p := range positions {
				if p.Size < 0.01 {
					continue
				}
				if _, tracked := buyTimesMap[p.Asset]; !tracked {
					continue
				}
				if p.CurPrice < 0.001 || p.CurPrice >= 0.99 || p.Redeemable {
					continue
				}
				emoji := "🟢"
				if p.CashPnL < 0 {
					emoji = "🔴"
				}
				title := p.Title
				if title == "" {
					title = p.Asset
					if len(title) > 20 {
						title = title[:8] + ".." + title[len(title)-4:]
					}
				}
				activeCost += p.InitialValue
				activeValue += p.CurrentValue
				bt := buyTimesMap[p.Asset]
				activeLines = append(activeLines, posLine{
					emoji: emoji, title: title, outcome: p.Outcome,
					avgPrice: p.AvgPrice, curPrice: p.CurPrice,
					size: p.Size, cost: p.InitialValue, value: p.CurrentValue,
					pnl: p.CashPnL, pct: p.PercentPnL,
					endDate: p.EndDate, redeemable: p.Redeemable,
					buyTime: bt,
				})
			}
			sort.Slice(activeLines, func(i, j int) bool { return activeLines[i].buyTime.After(activeLines[j].buyTime) })

			// --- Total assets & profit ---
			totalAssets := walletPUSD + activeValue
			capital := initialCapital
			totalProfit := totalAssets - capital
			totalProfitPct := 0.0
			if capital > 0 {
				totalProfitPct = totalProfit / capital * 100
			}

			// --- Daily profit (snapshot-based) ---
			snapshotFile := filepath.Join(filepath.Dir(journalDir), "daily_snapshot.json")
			todayStr := time.Now().In(sgt).Format("2006-01-02")
			type dailySnap struct {
				Date        string  `json:"date"`
				TotalAssets float64 `json:"total_assets"`
			}
			var snap dailySnap
			if raw, rerr := os.ReadFile(snapshotFile); rerr == nil {
				json.Unmarshal(raw, &snap)
			}
			if snap.Date != todayStr {
				snap = dailySnap{Date: todayStr, TotalAssets: totalAssets}
				if jb, jerr := json.Marshal(snap); jerr == nil {
					os.WriteFile(snapshotFile, jb, 0644)
				}
			}
			dailyProfit := totalAssets - snap.TotalAssets

			// --- Format header ---
			sb.WriteString(fmt.Sprintf("总资产: $%.2f (本金 $%.0f)\n", totalAssets, capital))
			sb.WriteString(fmt.Sprintf("总盈利: $%+.2f (%+.1f%%)\n", totalProfit, totalProfitPct))
			sb.WriteString(fmt.Sprintf("今日: $%+.2f\n", dailyProfit))
			sb.WriteString(fmt.Sprintf("闲钱: $%.2f pUSD\n", walletPUSD))
			sb.WriteString(fmt.Sprintf("持仓: %d 笔 · $%.2f 成本 · $%.2f 市值\n", len(activeLines), activeCost, activeValue))

			if len(activeLines) > 0 {
				sb.WriteString(fmt.Sprintf("\n--- 持仓明细 (%d) ---\n", len(activeLines)))
				shown := 0
				for _, l := range activeLines {
					title := l.title
					if len(title) > 40 {
						title = title[:40] + "…"
					}
					direction := l.outcome
					if direction == "" {
						direction = "YES"
					}
					timeTag := ""
					if !l.buyTime.IsZero() {
						timeTag = fmt.Sprintf(" · %s", l.buyTime.In(sgt).Format("01/02 15:04"))
					}
					sb.WriteString(fmt.Sprintf("%s %s · %s%s\n   %.1f份 · $%.2f成本 · 入%.3f→现%.3f · $%+.2f (%+.1f%%)\n",
						l.emoji, title, direction, timeTag,
						l.size, l.cost, l.avgPrice, l.curPrice, l.pnl, l.pct))
					shown++
					if shown >= 20 {
						sb.WriteString(fmt.Sprintf("... +%d more\n", len(activeLines)-20))
						break
					}
				}
			} else {
				sb.WriteString("\n无活跃持仓\n")
			}

			notifier.SidecarAlert(sb.String())
			slog.Info("hourly_pnl_pushed", "active", len(activeLines), "active_cost", activeCost, "active_value", activeValue, "wallet", walletPUSD, "total_assets", totalAssets, "total_profit", totalProfit, "daily_profit", dailyProfit)
		}
		time.Sleep(5 * time.Second)
		pushPnL()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				pushPnL()
			case <-pnlTrigger:
				time.Sleep(10 * time.Second)
				pushPnL()
			}
		}
	}()

	// Auto-redeem: every 1h check for redeemable winning positions and redeem on-chain.
	if onChain != nil && walletAddress != "" {
		go func() {
			tk := time.NewTicker(1 * time.Hour)
			defer tk.Stop()
			redeemedFile := filepath.Join("db", "redeemed.json")
			redeemed := make(map[string]bool)
			if data, err := os.ReadFile(redeemedFile); err == nil {
				json.Unmarshal(data, &redeemed)
			}
			saveRedeemed := func() {
				if data, err := json.Marshal(redeemed); err == nil {
					os.WriteFile(redeemedFile, data, 0644)
				}
			}
			checkRedeem := func() {
				positions, err := fetchDataAPIPositions(walletAddress)
				if err != nil {
					slog.Warn("redeem_fetch_err", "err", err)
					return
				}
				sgt := time.FixedZone("SGT", 8*3600)
				for _, p := range positions {
					if !p.Redeemable || p.CurPrice < 0.99 || p.Size < 0.01 {
						continue
					}
					if redeemed[p.Asset] {
						continue
					}
					slog.Info("redeem_candidate",
						"title", p.Title,
						"outcome", p.Outcome,
						"size", p.Size,
						"conditionId", p.ConditionID,
						"negRisk", p.NegativeRisk,
						"outcomeIndex", p.OutcomeIndex,
					)
					rCtx, rCancel := context.WithTimeout(ctx, 120*time.Second)
					err := onChain.RedeemPosition(rCtx, p.ConditionID, p.OutcomeIndex, p.Size, p.NegativeRisk)
					rCancel()
					if err != nil {
						slog.Error("redeem_failed", "title", p.Title, "err", err)
						notifier.SidecarAlert(fmt.Sprintf("❌ 赎回失败: %s\n%s · %.1f份\n错误: %s",
							p.Title, p.Outcome, p.Size, err))
					} else {
						redeemed[p.Asset] = true
						saveRedeemed()
						pnl := p.CurrentValue - p.InitialValue
						slog.Info("redeem_success", "title", p.Title, "size", p.Size, "pnl", pnl)
						notifier.SidecarAlert(fmt.Sprintf("💰 赎回成功: %s\n%s · %.1f份 · $%.2f→$%.2f · P&L $%+.2f\n%s SGT",
							p.Title, p.Outcome, p.Size, p.InitialValue, p.CurrentValue, pnl,
							time.Now().In(sgt).Format("01/02 15:04")))
					}
				}
			}
			time.Sleep(10 * time.Second)
			checkRedeem()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					checkRedeem()
				}
			}
		}()
	}

	if signalMode == "copytrade" {
		go func() {
			if err := ws.Run(ctx); err != nil && ctx.Err() == nil {
				slog.Warn("wss_run_exit", "err", err)
			}
		}()
		<-ctx.Done()
		return ctx.Err()
	}
	return ws.Run(ctx)
}

// countResolved returns the number of markets in the slice that have already
// settled on-chain (closed=true). Used for settlement-watcher status logging.
func countResolved(ms []feed.Market) int {
	n := 0
	for _, m := range ms {
		if m.Closed {
			n++
		}
	}
	return n
}

func settlementSlotForPosition(p strategy.Position, m feed.Market, me *assetMeta) (int, bool) {
	tokens := m.ClobTokenIDs()
	for i, token := range tokens {
		if token == p.AssetID {
			return i, true
		}
	}
	if me != nil && me.SlotIdx >= 0 {
		return me.SlotIdx, true
	}
	return -1, false
}

func paperPositionQuestion(p strategy.Position, m feed.Market, me *assetMeta) string {
	if p.Question != "" {
		return p.Question
	}
	if me != nil && me.Question != "" {
		return me.Question
	}
	return m.Question
}

func paperPositionOutcome(p strategy.Position, m feed.Market, slotIdx int, me *assetMeta) string {
	if p.Outcome != "" {
		return p.Outcome
	}
	if me != nil && me.Outcome != "" {
		return me.Outcome
	}
	outcomes := m.Outcomes()
	if slotIdx >= 0 && slotIdx < len(outcomes) {
		return outcomes[slotIdx]
	}
	return ""
}

func topWindow(ws []feed.WindowStats, n int) []feed.WindowStats {
	// sort by DeltaPP desc; naive O(n^2) good enough for small sets
	out := append([]feed.WindowStats(nil), ws...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].DeltaPP > out[i].DeltaPP {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if n > len(out) {
		n = len(out)
	}
	return out[:n]
}

func countBy(ms []feed.Market, pred func(feed.Market) bool) int {
	n := 0
	for _, m := range ms {
		if pred(m) {
			n++
		}
	}
	return n
}

func short(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:6] + ".." + id[len(id)-4:]
}

func fetchCLOBMidpoint(tokenID string) float64 {
	resp, err := nethttp.Get("https://clob.polymarket.com/midpoint?token_id=" + tokenID)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0
	}
	var result struct {
		Mid string `json:"mid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0
	}
	mid, _ := strconv.ParseFloat(result.Mid, 64)
	return mid
}

type dataAPIPosition struct {
	Size         float64 `json:"size"`
	AvgPrice     float64 `json:"avgPrice"`
	InitialValue float64 `json:"initialValue"`
	CurrentValue float64 `json:"currentValue"`
	CashPnL      float64 `json:"cashPnl"`
	PercentPnL   float64 `json:"percentPnl"`
	TotalBought  float64 `json:"totalBought"`
	RealizedPnL  float64 `json:"realizedPnl"`
	CurPrice     float64 `json:"curPrice"`
	Title        string  `json:"title"`
	Outcome      string  `json:"outcome"`
	Asset        string  `json:"asset"`
	ConditionID  string  `json:"conditionId"`
	EndDate      string  `json:"endDate"`
	Redeemable   bool    `json:"redeemable"`
	NegativeRisk bool    `json:"negativeRisk"`
	OutcomeIndex int     `json:"outcomeIndex"`
}

func fetchDataAPIPositions(walletAddr string) ([]dataAPIPosition, error) {
	reqURL := "https://data-api.polymarket.com/positions?user=" + strings.ToLower(walletAddr) + "&sizeThreshold=0.01&limit=200"
	resp, err := nethttp.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("data-api positions: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}
	var positions []dataAPIPosition
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, fmt.Errorf("data-api decode: %w", err)
	}
	return positions, nil
}

// assetMeta carries per-asset context used by signal prompts and log lines.
// Built once at startup from the gamma market list so the hot path never
// touches gamma again.
type assetMeta struct {
	Question       string
	ConditionID    string // market conditionId (0x…); needed for gamma settlement lookup
	SlotIdx        int    // index of this asset in market.Outcomes / OutcomePrices (0 or 1)
	Match          string // parsed title, e.g. "LoL: Shifters vs G2 Esports"
	Context        string // parsed context, e.g. "Game 1 Winner" or "BO3 · LCK ..."
	Outcome        string // this asset's outcome label ("Shifters", "Yes", ...)
	Sibling        string // sibling asset_id (the other outcome) — empty if market is non-binary
	SiblingOutcome string
	EndTime        time.Time // parsed from market.EndDate; zero if unparseable
	Slug           string    // market slug for newshare link
}

// buildAssetMeta walks a market list and produces an asset_id-keyed view that
// pairs each asset with its sibling outcome. Only binary markets get sibling
// info; multi-outcome markets (rare in LoL) degrade to "no sibling" which
// renders as a single-row prompt.
func buildAssetMeta(ms []feed.Market) map[string]*assetMeta {
	out := make(map[string]*assetMeta, len(ms)*2)
	for _, m := range ms {
		tokens := m.ClobTokenIDs()
		outcomes := m.Outcomes()
		match, ctx := notify.ParseMarketTitle(m.Question)
		var endTime time.Time
		if m.EndDate != "" {
			if t, err := time.Parse(time.RFC3339, m.EndDate); err == nil {
				endTime = t
			}
		}
		for i, id := range tokens {
			if id == "" {
				continue
			}
			me := &assetMeta{
				Question:    m.Question,
				ConditionID: m.ConditionID,
				SlotIdx:     i,
				Match:       match,
				Context:     ctx,
				EndTime:     endTime,
				Slug:        m.Slug,
			}
			if i < len(outcomes) {
				me.Outcome = outcomes[i]
			}
			// Sibling: for a 2-outcome market, the "other" token.
			if len(tokens) == 2 {
				sibIdx := 1 - i
				me.Sibling = tokens[sibIdx]
				if sibIdx < len(outcomes) {
					me.SiblingOutcome = outcomes[sibIdx]
				}
			}
			out[id] = me
		}
	}
	return out
}

// metaOutcome returns the Outcome label for an asset, or "" if unknown.
func metaOutcome(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Outcome
	}
	return ""
}

// metaQ returns the Question string for an asset, or "" if unknown. Used by log
// lines that previously indexed a plain map[string]string.
func metaQ(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Question
	}
	return ""
}

// metaMatch returns the parsed match title for an asset, or "" if unknown.
// Falls back to empty so the FillReceipt formatter can use Question instead.
func metaMatch(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Match
	}
	return ""
}

// outcomeOrDefault pulls the outcome label for a meta entry, falling back to
// def when the market has no outcome list (rare but defensive).
func outcomeOrDefault(me *assetMeta, def string) string {
	if me == nil || me.Outcome == "" {
		return def
	}
	return me.Outcome
}

// adminPromptReq is the payload for db/admin/send-prompt.trigger.
type adminPromptReq struct {
	AssetID string `json:"asset_id,omitempty"` // optional; defaults to top market's first token
	Note    string `json:"note,omitempty"`     // freeform tag appended to Context line
}

// sendAdminPrompt is invoked in-process by the daemon when the admin trigger
// file appears. It emits a SignalPrompt DM that routes through the SAME
// pending store the sidecar longpoll reads from, so callbacks Claim cleanly
// no matter which process wrote the trigger.
func sendAdminPrompt(raw []byte, mkts []feed.Market, meta map[string]*assetMeta, sampler *feed.Sampler, pending *notify.PendingStore, notifier notify.Notifier) error {
	var req adminPromptReq
	if len(strings.TrimSpace(string(raw))) > 0 {
		if err := json.Unmarshal(raw, &req); err != nil {
			return fmt.Errorf("parse trigger: %w", err)
		}
	}
	assetID := req.AssetID
	if assetID == "" {
		if len(mkts) == 0 {
			return fmt.Errorf("no markets subscribed")
		}
		tokens := mkts[0].ClobTokenIDs()
		if len(tokens) == 0 {
			return fmt.Errorf("top market has no clob tokens")
		}
		assetID = tokens[0]
	}
	me := meta[assetID]
	if me == nil {
		return fmt.Errorf("asset %s not in subscribed set", short(assetID))
	}

	mid := 0.50
	if w, ok := sampler.Window(assetID); ok && w.Samples > 0 {
		mid = w.EndMid
	}

	choices := []notify.Choice{{
		AssetID: assetID, Outcome: outcomeOrDefault(me, "Yes"),
		Mid: mid, IsSignal: true,
	}}
	sigChoices := []notify.SignalChoice{{
		Slot: 0, Outcome: choices[0].Outcome, Mid: mid, IsSignal: true,
	}}
	if me.Sibling != "" {
		sibMid := 1.0 - mid
		if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
			sibMid = w.EndMid
		}
		sibOutcome := me.SiblingOutcome
		if sibOutcome == "" {
			sibOutcome = "No"
		}
		choices = append(choices, notify.Choice{AssetID: me.Sibling, Outcome: sibOutcome, Mid: sibMid})
		sigChoices = append(sigChoices, notify.SignalChoice{Slot: 1, Outcome: sibOutcome, Mid: sibMid})
	}

	p := pending.Put(notify.PendingIntent{
		Market:   "admin-test",
		Question: me.Question,
		Choices:  choices,
	}, time.Now())

	ctxLine := me.Context
	note := strings.TrimSpace(req.Note)
	if note == "" {
		note = "PROMPT-TEST"
	}
	if ctxLine != "" {
		ctxLine += " · " + note
	} else {
		ctxLine = note
	}

	nonceSnap := p.Nonce
	notifier.SignalPrompt(notify.SignalPromptEvent{
		Nonce:     p.Nonce,
		Match:     me.Match,
		Context:   ctxLine,
		EndIn:     notify.HumanizeEndIn(time.Now(), me.EndTime),
		Slug:      me.Slug,
		Choices:   sigChoices,
		ExpiresIn: 10 * time.Minute,
		OnSent: func(msgID int64, err error) {
			if err != nil {
				slog.Warn("notify_send_err", "err", err)
				return
			}
			if msgID == 0 {
				return
			}
			pending.SetMessageID(nonceSnap, msgID)
		},
	})
	slog.Info("admin_prompt_sent",
		"asset", short(assetID),
		"nonce", p.Nonce,
		"mid", mid,
		"choices", len(choices),
		"note", note,
	)
	return nil
}

// runPromptTest now simply writes an admin trigger file and exits. The running
// daemon's watcher (see runDetect) will emit the prompt on its own pending
// store, so callbacks are Claim-able no matter how many daemon restarts or
// subprocess lifecycles happen between send and click.
func runPromptTest(ctx context.Context, _ float64) error {
	if _, err := os.Stat("db"); err != nil {
		return fmt.Errorf("db/ not found — run from the polymarket-go repo root (%w)", err)
	}
	if err := os.MkdirAll("db/admin", 0o755); err != nil {
		return fmt.Errorf("mkdir db/admin: %w", err)
	}
	payload, _ := json.Marshal(adminPromptReq{Note: "PROMPT-TEST"})
	triggerPath := "db/admin/send-prompt.trigger"
	if err := os.WriteFile(triggerPath, payload, 0o644); err != nil {
		return fmt.Errorf("write trigger: %w", err)
	}
	slog.Info("prompt_test.trigger_dropped",
		"path", triggerPath,
		"hint", "running daemon will pick it up within 1s and emit a prompt via its own pending store",
	)
	// Brief watch loop so the operator sees confirmation (trigger consumed).
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(triggerPath); err != nil {
			slog.Info("prompt_test.consumed", "ok", true)
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(500 * time.Millisecond):
		}
	}
	return fmt.Errorf("trigger still present after 8s — is the daemon running? (-mode=detect)")
}

// buyHandler wires a click on one outcome's Buy 1/5/10 → same paper-submit →
// pm.Open path the auto-mode signal loop uses, but honors the size the boss
// picked and the Choice (YES/NO) resolved from PendingIntent.Choices[slot].
// Executes synchronously on the longpoll goroutine; Telegram dispatch of the
// resulting DM is async via notifier.
type buyHandler struct {
	pm            *strategy.PositionManager
	exit          *strategy.ExitTracker
	ladder        *strategy.LadderTracker
	paper         order.Client
	rm            *risk.Manager
	pending       *notify.PendingStore
	closePending  *notify.CloseStore
	notifier      notify.Notifier
	meta          map[string]*assetMeta
	src           *sourceTracker
	recorder      *tickrec.Recorder
	jrn           *journal.Journal
	largeFillUSD  float64
	exitMode      string
	riskStatePath string
	savePositions func()
}

func (h *buyHandler) OnBuy(ctx context.Context, nonce string, slot int, sizeUSD float64, mode string, messageID int64) (string, error) {
	now := time.Now()
	p, ok := h.pending.Claim(nonce, now)
	if !ok {
		if h.notifier != nil && messageID != 0 {
			h.notifier.EditSignalExpired(messageID)
		}
		return "", fmt.Errorf("已过期或已点过")
	}
	if slot < 0 || slot >= len(p.Choices) {
		return "", fmt.Errorf("选项越界 slot=%d", slot)
	}
	choice := p.Choices[slot]
	if err := h.rm.AllowOpen(now); err != nil {
		st := h.rm.State()
		return "", fmt.Errorf("风控阻止: %s (day_pnl=%.2f dd=%.2f/%.2f)", st.BlockReason, st.DayRealizedPnL, st.DrawdownUSD, st.DrawdownCapUSD)
	}
	intent := order.Intent{
		AssetID: choice.AssetID,
		Market:  p.Market,
		Side:    order.Buy,
		SizeUSD: sizeUSD,
		LimitPx: choice.Mid,
		Type:    order.GTC,
	}
	res, err := h.paper.Submit(ctx, intent)
	if err != nil {
		return "", fmt.Errorf("下单失败: %s", err.Error())
	}
	entryTick := feed.Tick{
		AssetID: choice.AssetID, Market: p.Market,
		Time: now, Mid: res.AvgPrice,
	}
	pos, err := h.pm.OpenSized(choice.AssetID, p.Market, entryTick, sizeUSD)
	if err != nil {
		return "", fmt.Errorf("开仓失败: %s", err.Error())
	}
	_ = h.pm.SetOpenFee(pos.ID, res.FeeUSD)
	if h.savePositions != nil {
		h.savePositions()
	}

	// Mode branching: "hold" = hold-to-settlement (no SL, no timeout);
	// "ladder" = normal ladder with SL + 4h timeout.
	if mode != "hold" {
		switch h.exitMode {
		case "auto":
			if h.exit != nil {
				h.exit.Open(choice.AssetID, p.Market, entryTick)
			}
		case "ladder":
			if h.ladder != nil {
				h.ladder.Open(pos.ID, p.Market, choice.AssetID, entryTick, pos.Units)
			}
		}
	}

	if h.recorder != nil {
		if rerr := h.recorder.Start(pos.ID, choice.AssetID); rerr != nil {
			slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
		}
	}
	if h.src != nil {
		h.src.Mark(pos.ID, "manual", res.OrderID)
	}
	stats := h.pm.Stats()
	modeTag := "ladder"
	if mode == "hold" {
		modeTag = "hold"
	}
	slog.Info("manual_open",
		"id", pos.ID,
		"order_id", res.OrderID,
		"asset", short(choice.AssetID),
		"q", metaQ(h.meta, choice.AssetID),
		"outcome", choice.Outcome,
		"slot", slot,
		"size_usd", sizeUSD,
		"mode", modeTag,
		"signal_mid", choice.Mid,
		"entry_fill", res.AvgPrice,
		"units", pos.Units,
		"open_positions", stats.Open,
		"total_exposure_usd", stats.TotalExposure,
	)

	receipt := notify.FillReceiptEvent{
		Question: metaQ(h.meta, choice.AssetID),
		Match:    metaMatch(h.meta, choice.AssetID),
		Outcome:  choice.Outcome,
		SizeUSD:  sizeUSD,
		Units:    pos.Units,
		FillPx:   res.AvgPrice,
		OrderID:  res.OrderID,
		Source:   "manual",
	}
	if h.notifier != nil {
		if messageID != 0 {
			h.notifier.EditSignalFilled(receipt, messageID)
		}
		h.notifier.FillReceipt(receipt)
	}

	icon := "✅"
	if mode == "hold" {
		icon = "🔒"
	}
	return fmt.Sprintf("%s %s %gU @ %.4f · order %s",
		icon, choice.Outcome, sizeUSD, res.AvgPrice, short(res.OrderID)), nil
}

func (h *buyHandler) OnClose(ctx context.Context, nonce string, messageID int64) (string, error) {
	now := time.Now()
	ci, ok := h.closePending.Claim(nonce, now)
	if !ok {
		if h.notifier != nil && messageID != 0 {
			h.notifier.EditCloseDone("⌛ 已过期或已点过", messageID)
		}
		return "", fmt.Errorf("已过期或已点过")
	}

	var totalNetPnL float64
	var closedCount int
	for _, pos := range h.pm.Snapshot() {
		if pos.AssetID != ci.AssetID {
			continue
		}
		sellIntent := order.Intent{
			AssetID: pos.AssetID,
			Market:  pos.Market,
			Side:    order.Sell,
			SizeUSD: pos.SizeUSD,
			LimitPx: ci.WhalePrice,
			Type:    order.GTC,
		}
		res, serr := h.paper.Submit(ctx, sellIntent)
		if serr != nil {
			slog.Warn("whale_close_sell_reject", "pos", pos.ID, "err", serr.Error())
			continue
		}
		sig := strategy.ExitSignal{
			AssetID:  pos.AssetID,
			Market:   pos.Market,
			Time:     now,
			EntryMid: pos.EntryMid,
			PeakMid:  pos.EntryMid,
			ExitMid:  res.AvgPrice,
			HeldFor:  now.Sub(pos.EntryTime),
			ChangePP: (res.AvgPrice - pos.EntryMid) * 100,
			Reason:   strategy.ExitReason("whale_sell"),
		}
		closed, cerr := h.pm.Close(pos.ID, sig)
		if cerr != nil {
			slog.Warn("whale_close_miss", "pos", pos.ID, "err", cerr.Error())
			continue
		}
		if h.savePositions != nil {
			h.savePositions()
		}
		h.ladder.Forget(pos.ID)
		if h.recorder != nil {
			_ = h.recorder.Stop(closed.ID)
		}
		entryFeeShare := pos.OpenFeeUSD
		if pos.InitUnits > 0 {
			entryFeeShare = pos.OpenFeeUSD * (pos.Units / pos.InitUnits)
		}
		exitFee := res.FeeUSD
		netPnL := closed.PnLUSD - entryFeeShare - exitFee
		totalNetPnL += netPnL
		closedCount++
		stats := h.pm.Stats()
		closeSource, _ := h.src.Peek(closed.ID)
		if closeSource != "manual" {
			if tripped := h.rm.OnClose(netPnL, now); tripped {
				rst := h.rm.State()
				h.notifier.RiskTrip(notify.RiskTripEvent{
					Reason:        string(rst.BlockReason),
					DayPnLUSD:     rst.DayRealizedPnL,
					DayLossCapUSD: rst.DayLossCapUSD,
					DrawdownUSD:   rst.DrawdownUSD,
					DrawdownCap:   rst.DrawdownCapUSD,
					OpenPositions: stats.Open,
				})
			}
			if err := h.rm.SaveState(h.riskStatePath); err != nil {
				slog.Warn("risk_save_err", "err", err)
			}
		}
		if netPnL <= -h.largeFillUSD || netPnL >= h.largeFillUSD {
			h.notifier.LargeFill(notify.LargeFillEvent{
				Question: ci.Question,
				AssetID:  pos.AssetID,
				Side:     "sell",
				SizeUSD:  pos.SizeUSD,
				PnLUSD:   netPnL,
				EntryPx:  pos.EntryMid,
				ExitPx:   res.AvgPrice,
				Reason:   "whale_sell",
				HeldSec:  int(sig.HeldFor.Seconds()),
			})
		}
		source, openOID := h.src.Take(closed.ID)
		if h.jrn != nil {
			_ = h.jrn.Append(journal.TradeRecord{
				ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
				Question:     ci.Question,
				Outcome:      ci.Outcome,
				Side:         "buy",
				SizeUSD:      closed.SizeUSD,
				Units:        closed.Units,
				EntryMid:     closed.EntryMid,
				EntryTime:    closed.EntryTime,
				ExitMid:      closed.ExitMid,
				ExitTime:     closed.ExitTime,
				ExitReason:   "whale_sell",
				HeldSec:      int(sig.HeldFor.Seconds()),
				PnLUSD:       closed.PnLUSD,
				EntryFeeUSD:  entryFeeShare,
				ExitFeeUSD:   exitFee,
				NetPnLUSD:    netPnL,
				OpenOrderID:  openOID,
				CloseOrderID: res.OrderID,
				Mode:         "paper",
				SignalSource: source,
			})
		}
		slog.Info("whale_follow_close",
			"pos", pos.ID,
			"asset", short(pos.AssetID),
			"q", ci.Question,
			"outcome", ci.Outcome,
			"entry", pos.EntryMid,
			"exit_fill", res.AvgPrice,
			"net_pnl_usd", netPnL,
			"held_sec", int(sig.HeldFor.Seconds()),
			"open_positions", h.pm.Stats().Open,
		)
	}

	if closedCount == 0 {
		if messageID != 0 {
			h.notifier.EditCloseDone("⚠️ 无匹配持仓可平", messageID)
		}
		return "⚠️ 无匹配持仓", nil
	}
	result := fmt.Sprintf("✅ 已平仓 %d 笔 · pnl %+.2f USDC", closedCount, totalNetPnL)
	if messageID != 0 {
		h.notifier.EditCloseDone(result, messageID)
	}
	return result, nil
}

// buildNotifier returns a Telegram notifier when TELEGRAM_BOT_TOKEN + _CHAT_ID
// are present, otherwise a Nop so the trading loop is unconditional.
func parseAvailableBalance(errMsg string) float64 {
	// Parse: "balance: 22011730, sum of matched orders: 19986400, order amount..."
	balIdx := strings.Index(errMsg, "balance: ")
	matchIdx := strings.Index(errMsg, "sum of matched orders: ")
	if balIdx < 0 || matchIdx < 0 {
		return 0
	}
	balStr := errMsg[balIdx+len("balance: "):]
	if i := strings.IndexByte(balStr, ','); i > 0 {
		balStr = balStr[:i]
	}
	matchStr := errMsg[matchIdx+len("sum of matched orders: "):]
	if i := strings.IndexByte(matchStr, ','); i > 0 {
		matchStr = matchStr[:i]
	}
	bal, err1 := strconv.ParseFloat(strings.TrimSpace(balStr), 64)
	matched, err2 := strconv.ParseFloat(strings.TrimSpace(matchStr), 64)
	if err1 != nil || err2 != nil {
		return 0
	}
	return (bal - matched) / 1e6
}

func buildNotifier() notify.Notifier {
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		slog.Info("notify.ready", "mode", "nop", "reason", "telegram_env_missing")
		return notify.Nop{}
	}
	cfg := notify.TelegramConfig{
		BotToken:       tok,
		ChatID:         chat,
		PromptBotToken: os.Getenv("SIDECAR_BOT_TOKEN"),
		PushBotToken:   os.Getenv("PUSH_BOT_TOKEN"),
	}
	slog.Info("notify.ready",
		"mode", "telegram",
		"chat_id", chat,
		"prompt_via_sidecar", cfg.PromptBotToken != "",
		"push_via_separate", cfg.PushBotToken != "",
	)
	return notify.NewTelegram(cfg)
}

// sourceTracker remembers which path opened a position (auto detector vs manual
// click) so the journal can attribute closed trades correctly. Position state
// itself is in PositionManager; this is a small sidecar keyed by AssetID.
type sourceTracker struct {
	mu sync.Mutex
	m  map[string]sourceEntry
}

type sourceEntry struct {
	source      string // "auto" or "manual"
	openOrderID string
}

func newSourceTracker() *sourceTracker {
	return &sourceTracker{m: map[string]sourceEntry{}}
}

func (s *sourceTracker) Mark(assetID, source, openOID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[assetID] = sourceEntry{source: source, openOrderID: openOID}
}

// Take returns the recorded source + open order id and removes the entry.
// Missing entries default to "auto" with empty order id (safe for legacy).
// Keyed by posID so stacked positions on one asset keep distinct sources
// and ladder mode can Peek across partial closes without evicting prematurely.
func (s *sourceTracker) Take(posID string) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[posID]
	if !ok {
		return "auto", ""
	}
	delete(s.m, posID)
	return e.source, e.openOrderID
}

// runArbScan used to call paid third-party bookmaker odds APIs. Keep the CLI
// entrypoint as a clear no-op so old scripts fail safely without spending quota.
func runArbScan(_ context.Context) error {
	return fmt.Errorf("arb-scan disabled: third-party bookmaker odds API removed")
}

// Peek is like Take but leaves the entry in place — used for non-final
// ladder tranches so subsequent tranches of the same posID can still
// attribute their source + open order id.
func (s *sourceTracker) Peek(posID string) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[posID]
	if !ok {
		return "auto", ""
	}
	return e.source, e.openOrderID
}

// runDailyReport reads one SGT day's trade journal, prints the summary to
// stdout, and (when -report_push is set) DMs it via the Telegram alert bot.
// Default day = yesterday SGT — this matches the cron firing at 00:00:30 SGT
// to summarize the day that just ended.
func runDailyReport(ctx context.Context, dir, day string, push bool) error {
	if day == "" {
		yesterday := time.Now().In(journal.SGT).AddDate(0, 0, -1)
		day = yesterday.Format("2006-01-02")
	}
	trades, err := journal.Read(dir, day)
	if err != nil {
		return fmt.Errorf("read journal: %w", err)
	}
	summary := journal.Summarize(day, trades)
	out := journal.FormatTelegram(summary)
	fmt.Print(out)
	if !push {
		return nil
	}
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		return fmt.Errorf("report_push: TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID missing")
	}
	if err := sendTelegram(ctx, tok, chat, out); err != nil {
		return fmt.Errorf("telegram push: %w", err)
	}
	slog.Info("daily_report.pushed", "day", day, "trades", summary.Trades, "pnl_usd", summary.RealizedPnLUSD)
	return nil
}

func runDailyIterate(ctx context.Context, journalDir string, windowDays int, push bool) error {
	report, err := iterate.Analyze(journalDir, windowDays)
	if err != nil {
		return fmt.Errorf("iterate analyze: %w", err)
	}

	md := iterate.FormatMarkdown(report)

	reportsDir := "reports/daily"
	if err := os.MkdirAll(reportsDir, 0o755); err != nil {
		return fmt.Errorf("mkdir reports: %w", err)
	}
	mdPath := filepath.Join(reportsDir, report.Day+".md")
	if err := os.WriteFile(mdPath, []byte(md), 0o644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	slog.Info("daily_iterate.report_written", "path", mdPath, "trades", report.TotalTrades, "pnl", report.CumulativePnL)
	fmt.Print(md)

	if push {
		tok := os.Getenv("PUSH_BOT_TOKEN")
		if tok == "" {
			tok = os.Getenv("TELEGRAM_BOT_TOKEN")
		}
		chat := os.Getenv("TELEGRAM_CHAT_ID")
		if tok != "" && chat != "" {
			tgMsg := iterate.FormatTelegram(report)
			if err := sendTelegram(ctx, tok, chat, tgMsg); err != nil {
				slog.Warn("daily_iterate.push_fail", "err", sanitize.Error(err))
			} else {
				slog.Info("daily_iterate.pushed", "day", report.Day, "suggestions", len(report.Suggestions))
			}
		}
	}
	return nil
}

// sendTelegram fires a single sendMessage to the alert bot. We don't use the
// notify package because that's wired for fire-and-forget queues; the cron
// wants a synchronous send-and-exit.
func sendTelegram(ctx context.Context, token, chat, body string) error {
	api := "https://api.telegram.org/bot" + token + "/sendMessage"
	form := neturl.Values{}
	form.Set("chat_id", chat)
	form.Set("text", body)
	form.Set("disable_web_page_preview", "true")
	req, err := nethttp.NewRequestWithContext(ctx, "POST", api, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cl := &nethttp.Client{Timeout: 10 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("%s", sanitize.Error(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("telegram http %d", resp.StatusCode)
	}
	return nil
}

// injuryBlocksMomentum checks whether a momentum signal should be blocked
// because the team we'd bet on has injured stars. Only applies to basketball.
// Returns (blocked, teamName, playerList) for logging.
func injuryBlocksMomentum(sc *injury.Scanner, assetSport map[string]strategy.SportFamily, meta map[string]*assetMeta, assetID string) (bool, string, string) {
	if assetSport[assetID] != strategy.SportBasketball {
		return false, "", ""
	}
	me := meta[assetID]
	if me == nil || me.Outcome == "" {
		return false, "", ""
	}
	team := me.Outcome
	entries := sc.InjuredStars(team)
	if len(entries) == 0 {
		return false, "", ""
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Player + "(" + string(e.Status) + ")"
	}
	return true, team, strings.Join(names, ", ")
}

// injuryBlocksLottery checks whether a lottery candidate should be blocked
// because the underdog team has injured stars. Returns (blocked, teamName, playerList).
func injuryBlocksLottery(sc *injury.Scanner, meta map[string]*assetMeta, assetID string) (bool, string, string) {
	me := meta[assetID]
	if me == nil || me.Outcome == "" {
		return false, "", ""
	}
	team := me.Outcome
	entries := sc.InjuredStars(team)
	if len(entries) == 0 {
		return false, "", ""
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Player + "(" + string(e.Status) + ")"
	}
	return true, team, strings.Join(names, ", ")
}

func injuryGameFinished(team string, sc *injury.Scanner) bool {
	gi, ok := sc.GameFor(team)
	if !ok {
		return false
	}
	s := strings.ToLower(gi.Status)
	return strings.Contains(s, "final") || strings.Contains(s, "complete")
}

func injuryTeamInMarkets(team string, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily) bool {
	lt := strings.ToLower(team)
	for assetID, me := range meta {
		if assetSport[assetID] != strategy.SportBasketball {
			continue
		}
		if strings.Contains(strings.ToLower(me.Question), lt) || strings.Contains(lt, strings.ToLower(me.Outcome)) {
			return true
		}
	}
	return false
}

// injuryFindOpponent returns the opponent team name for a given team by scanning PM markets.
func injuryFindOpponent(team string, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily) string {
	opp, _, _ := injuryFindOpponentAndGame(team, meta, assetSport)
	return opp
}

func injuryFindOpponentAndGame(team string, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily) (opponent, gameCtx string, gameTime time.Time) {
	lt := strings.ToLower(team)
	for assetID, me := range meta {
		if assetSport[assetID] != strategy.SportBasketball {
			continue
		}
		if !strings.Contains(lt, strings.ToLower(me.Outcome)) {
			continue
		}
		if me.Sibling == "" {
			continue
		}
		if sib := meta[me.Sibling]; sib != nil {
			return sib.Outcome, me.Context, me.EndTime
		}
	}
	return "", "", time.Time{}
}

// injuryFindOpponentName resolves the opponent team name for a given team (ESPN first, PM fallback).
func injuryFindOpponentName(team string, injScanner *injury.Scanner, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily) string {
	if gi, ok := injScanner.GameFor(team); ok {
		if team == gi.HomeTeam || strings.Contains(strings.ToLower(team), strings.ToLower(gi.HomeTeam)) {
			return gi.AwayTeam
		}
		return gi.HomeTeam
	}
	opp, _, _ := injuryFindOpponentAndGame(team, meta, assetSport)
	return opp
}

// injuryBuildGameEvent builds a single InjuryAlertEvent for all alerts from the same game.
func injuryBuildGameEvent(alerts []injury.InjuryAlert, injScanner *injury.Scanner, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily, sampler *feed.Sampler) notify.InjuryAlertEvent {
	primary := alerts[0]
	ev := injuryBuildAlertEvent(primary, injScanner, meta, assetSport, sampler)
	// Populate TriggerPlayers from ALL alerts in this game group.
	buildTrigger := func(a injury.InjuryAlert) notify.InjuryInfo {
		return notify.InjuryInfo{
			Player:    a.StarPlayer,
			Status:    string(a.Status),
			Reason:    a.Reason,
			Role:      injury.PlayerRole(a.Team, a.StarPlayer),
			ImpactPct: injury.PlayerImpactPct(a.Team, a.StarPlayer),
		}
	}
	ev.TriggerPlayers = nil
	for _, a := range alerts {
		ev.TriggerPlayers = append(ev.TriggerPlayers, buildTrigger(a))
	}
	return ev
}

// injuryBuildAlertEvent constructs a rich InjuryAlertEvent with both teams' injury context.
func injuryBuildAlertEvent(a injury.InjuryAlert, injScanner *injury.Scanner, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily, sampler *feed.Sampler) notify.InjuryAlertEvent {
	ev := notify.InjuryAlertEvent{
		Team:       a.Team,
		StarPlayer: a.StarPlayer,
		Status:     string(a.Status),
		Reason:     a.Reason,
		Impact:     a.Impact,
	}

	// Build InjuryInfo with role and impact data
	buildInfo := func(team string, e injury.InjuryEntry) notify.InjuryInfo {
		return notify.InjuryInfo{
			Player:    e.Player,
			Status:    string(e.Status),
			Reason:    e.Reason,
			Role:      injury.PlayerRole(team, e.Player),
			ImpactPct: injury.PlayerImpactPct(team, e.Player),
		}
	}

	// Populate team injuries from scanner cache (all injuries, not just stars)
	teamEntries := injScanner.AllInjuries(a.Team)
	if len(teamEntries) > 0 {
		for _, e := range teamEntries {
			ev.TeamInjuries = append(ev.TeamInjuries, buildInfo(a.Team, e))
		}
	} else {
		for _, e := range a.Entries {
			ev.TeamInjuries = append(ev.TeamInjuries, buildInfo(a.Team, e))
		}
	}

	// Use ESPN scoreboard for accurate game time and series info
	if gi, ok := injScanner.GameFor(a.Team); ok {
		if a.Team == gi.HomeTeam || strings.Contains(strings.ToLower(a.Team), strings.ToLower(gi.HomeTeam)) {
			ev.OpponentName = gi.AwayTeam
			ev.MatchTitle = gi.AwayTeam + " @ " + gi.HomeTeam
		} else {
			ev.OpponentName = gi.HomeTeam
			ev.MatchTitle = a.Team + " @ " + gi.HomeTeam
		}
		ev.GameTime = gi.Tipoff
		ev.GameContext = gi.SeriesNote
		if gi.Status != "" && gi.Status != "Scheduled" {
			ev.GameContext += " (" + gi.Status + ")"
		}

		oppEntries := injScanner.AllInjuries(ev.OpponentName)
		for _, e := range oppEntries {
			ev.OpponentInj = append(ev.OpponentInj, buildInfo(ev.OpponentName, e))
		}
	} else {
		// Fallback to PM market data
		opponent, gameCtx, gameTime := injuryFindOpponentAndGame(a.Team, meta, assetSport)
		if opponent != "" {
			ev.OpponentName = opponent
			ev.MatchTitle = a.Team + " vs " + opponent
			ev.GameContext = gameCtx
			ev.GameTime = gameTime

			oppEntries := injScanner.AllInjuries(opponent)
			for _, e := range oppEntries {
				ev.OpponentInj = append(ev.OpponentInj, buildInfo(opponent, e))
			}
		}
	}

	// Populate PM prices from sampler
	if sampler != nil {
		lt := strings.ToLower(a.Team)
		for assetID, me := range meta {
			if assetSport[assetID] != strategy.SportBasketball {
				continue
			}
			lo := strings.ToLower(me.Outcome)
			if strings.Contains(lt, lo) || strings.Contains(lo, lt) {
				if w, ok := sampler.Window(assetID); ok && w.Samples > 0 {
					ev.TeamPrice = w.EndMid
				}
				if me.Sibling != "" {
					if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
						ev.OpponentPrice = w.EndMid
					}
				}
				ev.Slug = me.Slug
				break
			}
		}
	}

	// Fallback: generate slug from ESPN game data when not found in WSS-tracked markets
	if ev.Slug == "" && !ev.GameTime.IsZero() {
		if gi, ok := injScanner.GameFor(a.Team); ok {
			ev.Slug = injuryGuessSlug(gi.AwayTeam, gi.HomeTeam, ev.GameTime)
		}
	}

	return ev
}

// injuryGuessSlug generates a PM-style slug from team names and game date.
// PM NBA slugs follow the pattern: nba-{away_abbr}-{home_abbr}-YYYY-MM-DD
func injuryGuessSlug(teamA, teamB string, gameTime time.Time) string {
	abbr := map[string]string{
		"Atlanta Hawks": "atl", "Boston Celtics": "bos", "Brooklyn Nets": "bkn",
		"Charlotte Hornets": "cha", "Chicago Bulls": "chi", "Cleveland Cavaliers": "cle",
		"Dallas Mavericks": "dal", "Denver Nuggets": "den", "Detroit Pistons": "det",
		"Golden State Warriors": "gsw", "Houston Rockets": "hou", "Indiana Pacers": "ind",
		"LA Clippers": "lac", "Los Angeles Lakers": "lal", "Memphis Grizzlies": "mem",
		"Miami Heat": "mia", "Milwaukee Bucks": "mil", "Minnesota Timberwolves": "min",
		"New Orleans Pelicans": "nop", "New York Knicks": "nyk", "Oklahoma City Thunder": "okc",
		"Orlando Magic": "orl", "Philadelphia 76ers": "phi", "Phoenix Suns": "phx",
		"Portland Trail Blazers": "por", "Sacramento Kings": "sac", "San Antonio Spurs": "sas",
		"Toronto Raptors": "tor", "Utah Jazz": "uta", "Washington Wizards": "was",
	}
	a, b := abbr[teamA], abbr[teamB]
	if a == "" || b == "" {
		return ""
	}
	utcDate := gameTime.UTC().Format("2006-01-02")
	return fmt.Sprintf("nba-%s-%s-%s", a, b, utcDate)
}

// injuryPushOpponentPrompt finds the PM market for an injured team's game and
// pushes a SignalPrompt with buy buttons for the opposing team.
func injuryPushOpponentPrompt(a injury.InjuryAlert, meta map[string]*assetMeta, assetSport map[string]strategy.SportFamily, sampler *feed.Sampler, pending *notify.PendingStore, notifier notify.Notifier) {
	lowerTeam := strings.ToLower(a.Team)
	for assetID, me := range meta {
		if me.Sibling == "" {
			continue
		}
		if assetSport[assetID] != strategy.SportBasketball {
			continue
		}
		if !strings.Contains(lowerTeam, strings.ToLower(me.Outcome)) {
			continue
		}
		sibMe := meta[me.Sibling]
		if sibMe == nil {
			continue
		}

		sibMid := 0.50
		if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
			sibMid = w.EndMid
		}
		injMid := 1.0 - sibMid
		if w, ok := sampler.Window(assetID); ok && w.Samples > 0 {
			injMid = w.EndMid
		}

		choices := []notify.Choice{
			{AssetID: me.Sibling, Outcome: sibMe.Outcome, Mid: sibMid, IsSignal: true},
			{AssetID: assetID, Outcome: me.Outcome, Mid: injMid},
		}
		sigChoices := []notify.SignalChoice{
			{Slot: 0, Outcome: sibMe.Outcome, Mid: sibMid, IsSignal: true},
			{Slot: 1, Outcome: me.Outcome, Mid: injMid},
		}

		p := pending.Put(notify.PendingIntent{
			Market:   "injury-alert",
			Question: me.Question,
			Choices:  choices,
		}, time.Now())

		ctxLine := fmt.Sprintf("🚨 %s %s OUT", a.Team, a.StarPlayer)
		if a.Reason != "" {
			ctxLine += " · " + a.Reason
		}
		teamOut := 0
		for _, e := range a.Entries {
			if e.Status == injury.StatusOut {
				teamOut++
			}
		}
		if teamOut > 1 {
			ctxLine += fmt.Sprintf("\n%s 共 %d 名球员缺阵", a.Team, teamOut)
		}

		nonce := p.Nonce
		notifier.SignalPrompt(notify.SignalPromptEvent{
			Nonce:   p.Nonce,
			Match:   me.Match,
			Context: ctxLine,
			EndIn:   notify.HumanizeEndIn(time.Now(), me.EndTime),
			Slug:    me.Slug,
			Choices: sigChoices,
			OnSent: func(msgID int64, err error) {
				if err != nil {
					slog.Warn("notify_send_err", "err", err)
					return
				}
				if msgID == 0 {
					return
				}
				pending.SetMessageID(nonce, msgID)
			},
		})
		slog.Info("injury_opponent_prompt",
			"injured_team", a.Team,
			"player", a.StarPlayer,
			"opponent", sibMe.Outcome,
			"mid", sibMid,
			"nonce", nonce,
		)
		break
	}
}

type whaleConfirmDecision struct {
	Ready   bool
	Reason  string
	Event   string
	Wallets int
	Need    int
}

type whaleConfirmGate struct {
	mu            sync.Mutex
	window        time.Duration
	minWallets    int
	bypassUSD     float64
	maxWorsePrice float64
	events        map[string]*whaleConfirmState
}

type whaleConfirmState struct {
	firstTime time.Time
	minPrice  float64
	wallets   map[string]struct{}
}

func newWhaleConfirmGate(window time.Duration, minWallets int, bypassUSD, maxWorsePrice float64) *whaleConfirmGate {
	if minWallets <= 0 {
		minWallets = 1
	}
	return &whaleConfirmGate{
		window:        window,
		minWallets:    minWallets,
		bypassUSD:     bypassUSD,
		maxWorsePrice: maxWorsePrice,
		events:        map[string]*whaleConfirmState{},
	}
}

func whaleDirectGateNote(list string) string {
	list = strings.TrimSpace(strings.ToLower(list))
	if list == "" {
		return "gate direct"
	}
	return "gate " + list + " direct"
}

func whaleConfirmDecisionNote(decision whaleConfirmDecision) string {
	if strings.Contains(decision.Reason, "bypass_notional") {
		return "gate 5k+ bypass"
	}
	if decision.Ready && decision.Wallets > 0 && decision.Need > 0 && decision.Wallets >= decision.Need {
		return fmt.Sprintf("gate %d-wallet confirmed", decision.Wallets)
	}
	if decision.Reason == "confirmation_disabled" {
		return "gate disabled"
	}
	if decision.Reason != "" {
		return "gate " + decision.Reason
	}
	return "gate confirmed"
}

func formatWhalePromptContext(ev whale.AlertEvent, meta walletFileMeta, gateNote string) string {
	parts := make([]string, 0, 6)
	if feed.IsFootballScoreMarketText(ev.Question + " " + ev.Slug) {
		parts = append(parts, "⚽ 比分盘")
	}
	parts = append(parts,
		fmt.Sprintf("$%.0f", ev.Notional),
		fmt.Sprintf("%.0f shares", ev.SizeUnits),
	)
	if gateNote != "" {
		parts = append(parts, gateNote)
	}
	if meta.List == "scout" {
		parts = append(parts, "scout observe")
	}
	if meta.SmartMoneyScore > 0 || meta.BotScore > 0 {
		parts = append(parts, fmt.Sprintf("S%.0f/B%.0f", meta.SmartMoneyScore, meta.BotScore))
	}
	if ev.TotalShares > 0 {
		parts = append(parts, fmt.Sprintf("pos %.0f @ %.4f", ev.TotalShares, ev.AvgPrice))
	}
	return "🐋 " + strings.Join(parts, " · ")
}

func (g *whaleConfirmGate) Observe(ev whale.AlertEvent, list string) whaleConfirmDecision {
	event := whaleEventKeyForAlert(ev)
	need := g.minWallets
	if need <= 1 || g.window <= 0 {
		return whaleConfirmDecision{Ready: true, Reason: "confirmation_disabled", Event: event, Wallets: 1, Need: need}
	}
	if g.bypassUSD > 0 && whaleNotionalAtLeast(ev.Notional, g.bypassUSD) {
		return whaleConfirmDecision{Ready: true, Reason: fmt.Sprintf("bypass_notional:$%.0f", ev.Notional), Event: event, Wallets: 1, Need: need}
	}
	outcome := strings.ToLower(strings.TrimSpace(ev.Outcome))
	key := event + "|" + outcome
	now := ev.Timestamp
	if now.IsZero() {
		now = time.Now()
	}
	wallet := strings.ToLower(strings.TrimSpace(ev.Wallet))

	g.mu.Lock()
	defer g.mu.Unlock()
	st := g.events[key]
	if st == nil || (!st.firstTime.IsZero() && now.Sub(st.firstTime) > g.window) {
		st = &whaleConfirmState{
			firstTime: now,
			minPrice:  ev.Price,
			wallets:   map[string]struct{}{},
		}
		g.events[key] = st
	}
	if st.minPrice <= 0 || ev.Price < st.minPrice {
		st.minPrice = ev.Price
	}
	if wallet != "" {
		st.wallets[wallet] = struct{}{}
	}
	wallets := len(st.wallets)
	if wallets < need {
		return whaleConfirmDecision{
			Event:   event,
			Wallets: wallets,
			Need:    need,
			Reason:  fmt.Sprintf("confirm_wallets:%d/%d event=%s list=%s", wallets, need, event, list),
		}
	}
	if g.maxWorsePrice >= 0 && ev.Price > st.minPrice+g.maxWorsePrice {
		return whaleConfirmDecision{
			Event:   event,
			Wallets: wallets,
			Need:    need,
			Reason:  fmt.Sprintf("confirm_price_worse:%.4f>%.4f event=%s", ev.Price, st.minPrice+g.maxWorsePrice, event),
		}
	}
	return whaleConfirmDecision{
		Ready:   true,
		Event:   event,
		Wallets: wallets,
		Need:    need,
		Reason:  fmt.Sprintf("confirmed_wallets:%d/%d event=%s", wallets, need, event),
	}
}

func parseCSVSet(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" {
			out[part] = struct{}{}
		}
	}
	return out
}

func parseWhaleListMinUSD(raw string) map[string]float64 {
	out := map[string]float64{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			slog.Warn("invalid WHALE_LIST_MIN_USD entry", "entry", part)
			continue
		}
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" {
			slog.Warn("invalid WHALE_LIST_MIN_USD entry", "entry", part)
			continue
		}
		minUSD, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil || minUSD < 0 {
			slog.Warn("invalid WHALE_LIST_MIN_USD value", "entry", part, "err", err)
			continue
		}
		out[key] = minUSD
	}
	return out
}

func sortedSetKeys(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func whaleNotionalAtLeast(notional, threshold float64) bool {
	if threshold <= 0 {
		return true
	}
	const centsTolerance = 0.01
	return notional+centsTolerance >= threshold
}

func whaleRepeatBypassesCooldown(notional, threshold float64) bool {
	return threshold > 0 && whaleNotionalAtLeast(notional, threshold)
}

type whaleEdgeBlockConfig struct {
	Min15mSamples  int
	Max15mAvgPP    float64
	Min1hSamples   int
	Max1hAvgPP     float64
	HotMinUSD      float64
	HotMinSamples  int
	HotMinAvgPP    float64
	HotMinWinRate  float64
	HotMin5mAvgPP  float64
	HotMin15mAvgPP float64
	HotMax1hNegPP  float64
	Refresh        time.Duration
}

type whaleEdgeBlockCache struct {
	path     string
	cfg      whaleEdgeBlockConfig
	mu       sync.Mutex
	loadedAt time.Time
	blocks   map[string]string
	hot      map[string]string
	lastErr  error
}

func newWhaleEdgeBlockCache(path string, cfg whaleEdgeBlockConfig) *whaleEdgeBlockCache {
	return &whaleEdgeBlockCache{path: path, cfg: cfg, blocks: map[string]string{}, hot: map[string]string{}}
}

func (c *whaleEdgeBlockCache) Reason(wallet string) (string, bool) {
	if c == nil {
		return "", false
	}
	key := strings.ToLower(strings.TrimSpace(wallet))
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.shouldReloadLocked() {
		now := time.Now()
		blocks, hot, err := loadWhaleEdgeSignals(c.path, c.cfg)
		if err != nil {
			c.lastErr = err
			c.loadedAt = now
			slog.Warn("whale_edge_blocks_reload_fail", "path", c.path, "err", err)
		} else {
			c.blocks = blocks
			c.hot = hot
			c.lastErr = nil
			c.loadedAt = now
		}
	}
	if key == "" {
		return "", false
	}
	reason, ok := c.blocks[key]
	return reason, ok
}

func (c *whaleEdgeBlockCache) HotReason(wallet string) (string, bool) {
	if c == nil {
		return "", false
	}
	key := strings.ToLower(strings.TrimSpace(wallet))
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.shouldReloadLocked() {
		now := time.Now()
		blocks, hot, err := loadWhaleEdgeSignals(c.path, c.cfg)
		if err != nil {
			c.lastErr = err
			c.loadedAt = now
			slog.Warn("whale_edge_hot_reload_fail", "path", c.path, "err", err)
		} else {
			c.blocks = blocks
			c.hot = hot
			c.lastErr = nil
			c.loadedAt = now
		}
	}
	if key == "" {
		return "", false
	}
	reason, ok := c.hot[key]
	return reason, ok
}

func (c *whaleEdgeBlockCache) shouldReloadLocked() bool {
	if c.path == "" {
		return false
	}
	if c.loadedAt.IsZero() {
		return true
	}
	if c.cfg.Refresh <= 0 {
		return false
	}
	return time.Since(c.loadedAt) >= c.cfg.Refresh
}

type whaleEdgeSnapshot struct {
	Wallet     string  `json:"wallet"`
	HorizonSec int64   `json:"horizon_sec"`
	DeltaPP    float64 `json:"delta_pp"`
}

type whaleEdgeStats struct {
	Samples int
	SumPP   float64
	Wins    int
}

func loadWhaleNegativeEdgeBlocks(path string, cfg whaleEdgeBlockConfig) (map[string]string, error) {
	blocks, _, err := loadWhaleEdgeSignals(path, cfg)
	return blocks, err
}

func loadWhaleEdgeSignals(path string, cfg whaleEdgeBlockConfig) (map[string]string, map[string]string, error) {
	out := map[string]string{}
	hot := map[string]string{}
	if strings.TrimSpace(path) == "" {
		return out, hot, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, hot, nil
		}
		return nil, nil, err
	}
	defer f.Close()

	metrics := map[string]map[int64]*whaleEdgeStats{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var snap whaleEdgeSnapshot
		if err := json.Unmarshal(sc.Bytes(), &snap); err != nil {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(snap.Wallet))
		if wallet == "" {
			continue
		}
		byHorizon := metrics[wallet]
		if byHorizon == nil {
			byHorizon = map[int64]*whaleEdgeStats{}
			metrics[wallet] = byHorizon
		}
		st := byHorizon[snap.HorizonSec]
		if st == nil {
			st = &whaleEdgeStats{}
			byHorizon[snap.HorizonSec] = st
		}
		st.Samples++
		st.SumPP += snap.DeltaPP
		if snap.DeltaPP > 0 {
			st.Wins++
		}
	}
	if err := sc.Err(); err != nil {
		return nil, nil, err
	}

	for wallet, byHorizon := range metrics {
		if st := byHorizon[int64((15 * time.Minute).Seconds())]; st != nil && cfg.Min15mSamples > 0 && st.Samples >= cfg.Min15mSamples {
			avg := st.SumPP / float64(st.Samples)
			if avg <= cfg.Max15mAvgPP {
				out[wallet] = fmt.Sprintf("15m edge %.2fpp over %d samples", avg, st.Samples)
				continue
			}
		}
		if st := byHorizon[int64((time.Hour).Seconds())]; st != nil && cfg.Min1hSamples > 0 && st.Samples >= cfg.Min1hSamples {
			avg := st.SumPP / float64(st.Samples)
			if avg <= cfg.Max1hAvgPP {
				out[wallet] = fmt.Sprintf("1h edge %.2fpp over %d samples", avg, st.Samples)
			}
		}
		if _, blocked := out[wallet]; blocked {
			continue
		}
		if reason, ok := whaleEdgeHotReason(byHorizon, cfg); ok {
			hot[wallet] = reason
		}
	}
	return out, hot, nil
}

func whaleEdgeHotReason(byHorizon map[int64]*whaleEdgeStats, cfg whaleEdgeBlockConfig) (string, bool) {
	var total whaleEdgeStats
	for horizon, st := range byHorizon {
		if horizon <= 0 || st == nil {
			continue
		}
		total.Samples += st.Samples
		total.SumPP += st.SumPP
		total.Wins += st.Wins
	}
	if total.Samples == 0 {
		return "", false
	}
	avg := total.SumPP / float64(total.Samples)
	winRate := float64(total.Wins) / float64(total.Samples) * 100
	if cfg.HotMinSamples > 0 && total.Samples < cfg.HotMinSamples {
		return "", false
	}
	if avg < cfg.HotMinAvgPP {
		return "", false
	}
	if winRate < cfg.HotMinWinRate {
		return "", false
	}
	if st := byHorizon[int64((5 * time.Minute).Seconds())]; st == nil || st.Samples == 0 || st.SumPP/float64(st.Samples) < cfg.HotMin5mAvgPP {
		return "", false
	}
	if st := byHorizon[int64((15 * time.Minute).Seconds())]; st == nil || st.Samples == 0 || st.SumPP/float64(st.Samples) < cfg.HotMin15mAvgPP {
		return "", false
	}
	if st := byHorizon[int64((time.Hour).Seconds())]; st != nil && st.Samples > 0 {
		if st.SumPP/float64(st.Samples) <= cfg.HotMax1hNegPP {
			return "", false
		}
	}
	return fmt.Sprintf("avg %.2fpp win %.0f%% over %d samples", avg, winRate, total.Samples), true
}

func parseWhaleEnvInt(name string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		slog.Warn("invalid "+name, "value", raw, "err", err)
		return fallback
	}
	return v
}

func parseWhaleEnvFloat(name string, fallback float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		slog.Warn("invalid "+name, "value", raw, "err", err)
		return fallback
	}
	return v
}

func loadWalletFileMetas(path string) (map[string]walletFileMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := map[string]walletFileMeta{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "#", 2)
		addr := strings.Fields(strings.TrimSpace(parts[0]))
		if len(addr) == 0 {
			continue
		}
		key := strings.ToLower(addr[0])
		if !strings.HasPrefix(key, "0x") || len(key) != 42 {
			continue
		}
		meta := walletFileMeta{}
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
					meta.SmartMoneyScore, _ = strconv.ParseFloat(v, 64)
				case "bot":
					meta.BotScore, _ = strconv.ParseFloat(v, 64)
				}
			}
		}
		out[key] = meta
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

var (
	whaleEventParenRE = regexp.MustCompile(`\s*\([^)]*\b(?:bo[0-9]+|game|map)\b[^)]*\)`)
	whaleEventGameRE  = regexp.MustCompile(`\s*-\s*(?:game|map)\s*[0-9]+\s+winner\b.*$`)
	whaleEventSpaceRE = regexp.MustCompile(`\s+`)
)

const (
	followMinEntryPrice         = 0.05
	followMaxEntryPrice         = 0.95
	footballScoreMinEntryPrice  = 0.01
	whaleSettlementSellMinPrice = 0.98
)

func isTargetFollowMarket(q, slug string) bool {
	ok, _ := targetFollowMarketDecision(q, slug)
	return ok
}

func whaleMarketDecision(q, slug, outcome, _ string) (bool, string) {
	if feed.IsFootballScoreMarketText(q + " " + slug) {
		if strings.EqualFold(strings.TrimSpace(outcome), "No") {
			return false, "football_score_no_filtered"
		}
		return true, ""
	}
	return targetFollowMarketDecision(q, slug)
}

func copytradeMarketDecision(q, slug, outcome string, allowFootballScore bool) (bool, string) {
	if feed.IsFootballScoreMarketText(q + " " + slug) {
		if !allowFootballScore {
			return false, "derivative_filtered"
		}
		if strings.EqualFold(strings.TrimSpace(outcome), "No") {
			return false, "football_score_no_filtered"
		}
		return true, ""
	}
	return targetFollowMarketDecision(q, slug)
}

func copytradeEntryPriceFloor(footballScore, allowFootballScore bool) float64 {
	if footballScore && allowFootballScore {
		return footballScoreMinEntryPrice
	}
	return followMinEntryPrice
}

func copytradeMarketSize(base float64, footballScore, allowFootballScore bool, scoreCap float64) float64 {
	if !footballScore || !allowFootballScore || scoreCap <= 0 || base <= scoreCap {
		return base
	}
	return scoreCap
}

func targetFollowMarketDecision(q, slug string) (bool, string) {
	text := strings.ToLower(q + " " + slug)
	slug = strings.ToLower(slug)
	if isDerivativeFollowMarketText(text) {
		return false, "derivative_filtered"
	}
	if feed.IsOutrightFollowMarketText(text) {
		return false, "outright_filtered"
	}
	if strings.Contains(text, "tennis") ||
		strings.Contains(text, "wimbledon") ||
		strings.Contains(text, " atp") ||
		strings.Contains(text, " wta") ||
		strings.HasPrefix(slug, "atp-") ||
		strings.HasPrefix(slug, "wta-") {
		return false, "category_filtered"
	}
	if feed.IsFollowTargetMarket(feed.Market{Question: q, Slug: slug}) {
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

func whalePriceDecision(side string, price float64, footballScore bool) (bool, string) {
	switch strings.ToUpper(strings.TrimSpace(side)) {
	case "BUY":
		minPrice := followMinEntryPrice
		if footballScore {
			minPrice = footballScoreMinEntryPrice
		}
		if price < minPrice || price > followMaxEntryPrice {
			return false, "price_filtered"
		}
	case "SELL":
		if price >= whaleSettlementSellMinPrice {
			return false, "settlement_sell_filtered"
		}
	}
	return true, ""
}

func isDerivativeFollowMarketText(text string) bool {
	return feed.IsDerivativeFollowMarketText(text)
}

func whaleEventKey(market string) string {
	s := strings.ToLower(strings.TrimSpace(market))
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "—", "-")
	s = whaleEventGameRE.ReplaceAllString(s, "")
	if idx := strings.Index(s, " - "); idx >= 0 {
		s = s[:idx]
	}
	if vs := strings.Index(s, " vs"); vs >= 0 {
		if colon := strings.Index(s[vs:], ":"); colon >= 0 {
			s = s[:vs+colon]
		}
	}
	s = whaleEventParenRE.ReplaceAllString(s, "")
	s = strings.TrimSuffix(s, " winner")
	s = strings.TrimSuffix(s, " match winner")
	s = strings.TrimSpace(strings.Trim(s, "-:"))
	s = whaleEventSpaceRE.ReplaceAllString(s, " ")
	if s == "" {
		return "unknown"
	}
	return s
}

func whaleEventKeyForAlert(ev whale.AlertEvent) string {
	if !feed.IsFootballScoreMarketText(ev.Question + " " + ev.Slug) {
		return whaleEventKey(ev.Question)
	}
	if conditionID := strings.ToLower(strings.TrimSpace(ev.ConditionID)); conditionID != "" {
		return "football_score:" + conditionID
	}
	market := whaleEventSpaceRE.ReplaceAllString(strings.ToLower(strings.TrimSpace(ev.Question)), " ")
	outcome := strings.ToLower(strings.TrimSpace(ev.Outcome))
	return "football_score:" + market + "|" + outcome
}

func loadWhaleCooldownSeeds(path string, repeatCooldown, eventCooldown time.Duration) (map[string]time.Time, map[string]time.Time, int, int) {
	repeat := map[string]time.Time{}
	events := map[string]time.Time{}
	f, err := os.Open(path)
	if err != nil {
		return repeat, events, 0, 0
	}
	defer f.Close()

	now := time.Now()
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 8*1024*1024)
	for sc.Scan() {
		var rec struct {
			TS          string  `json:"ts"`
			Wallet      string  `json:"wallet"`
			Side        string  `json:"side"`
			Market      string  `json:"market"`
			Outcome     string  `json:"outcome"`
			Slug        string  `json:"slug"`
			AssetID     string  `json:"asset_id"`
			ConditionID string  `json:"condition_id"`
			Action      string  `json:"action"`
			Size        float64 `json:"size"`
		}
		if err := json.Unmarshal(sc.Bytes(), &rec); err != nil {
			continue
		}
		action := strings.ToLower(strings.TrimSpace(rec.Action))
		if action != "alert" && action != "followed" {
			continue
		}
		if strings.ToUpper(strings.TrimSpace(rec.Side)) != "BUY" {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(rec.Wallet))
		if wallet == "" {
			continue
		}
		ts, err := time.Parse(time.RFC3339, rec.TS)
		if err != nil {
			continue
		}
		if repeatCooldown > 0 && rec.AssetID != "" && now.Sub(ts) < repeatCooldown {
			key := wallet + "|" + rec.AssetID
			if ts.After(repeat[key]) {
				repeat[key] = ts
			}
		}
		if eventCooldown > 0 && rec.Market != "" && now.Sub(ts) < eventCooldown {
			key := wallet + "|" + whaleEventKeyForAlert(whale.AlertEvent{
				Question: rec.Market, Outcome: rec.Outcome, Slug: rec.Slug, ConditionID: rec.ConditionID,
			})
			if ts.After(events[key]) {
				events[key] = ts
			}
		}
	}
	return repeat, events, len(repeat), len(events)
}

func loadWhaleSeenTradeIDs(path string) map[string]struct{} {
	seen := map[string]struct{}{}
	f, err := os.Open(path)
	if err != nil {
		return seen
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 8*1024*1024)
	for sc.Scan() {
		var rec struct {
			Action  string `json:"action"`
			Wallet  string `json:"wallet"`
			TradeID string `json:"trade_id"`
		}
		if err := json.Unmarshal(sc.Bytes(), &rec); err != nil {
			continue
		}
		if strings.ToLower(strings.TrimSpace(rec.Action)) == "skip" {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(rec.Wallet))
		tradeID := strings.ToLower(strings.TrimSpace(rec.TradeID))
		if wallet == "" || tradeID == "" {
			continue
		}
		seen[wallet+"|"+tradeID] = struct{}{}
	}
	return seen
}
