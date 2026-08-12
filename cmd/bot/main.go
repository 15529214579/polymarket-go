package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
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
	takerFeeRate := flag.Float64("taker_fee_rate", 0, "paper dynamic taker fee rate; fee = shares * rate * price * (1-price)")
	ladderTP1Pct := flag.Float64("ladder_tp1_pct", 9.99, "ladder TP1 trigger: 9.99 = effectively disabled (ride to settlement/timeout)")
	ladderTP1Frac := flag.Float64("ladder_tp1_frac", 0.50, "fraction of initial units to close on TP1")
	ladderTP2Pct := flag.Float64("ladder_tp2_pct", 9.99, "ladder TP2 trigger: 9.99 = effectively disabled")
	ladderTP2Frac := flag.Float64("ladder_tp2_frac", 1.00, "fraction of initial units to close on TP2 (1.0 = all remaining)")
	ladderSLPct := flag.Float64("ladder_sl_pct", 0.15, "ladder stop-loss: exit ≤ entry × (1 - this) closes 100%")
	ladderTrailingPct := flag.Float64("ladder_trailing_pct", 0, "post-TP1 trailing drawdown in return terms; 0 disables")
	ladderSLConfirm := flag.Duration("ladder_sl_confirm", 0, "loss must persist for this duration before ladder SL")
	ladderMinHoldBeforeSL := flag.Duration("ladder_min_hold_before_sl", 0, "minimum hold before ladder SL can trigger")
	ladderFeeAware := flag.Bool("ladder_fee_aware", false, "evaluate TP/SL/trailing thresholds after entry fee, exit fee, and slippage")
	ladderRequireBid := flag.Bool("ladder_require_bid", false, "require an executable best bid before evaluating ladder exits")
	ladderMaxHold := flag.Duration("ladder_max_hold", 6*time.Hour, "ladder hard timeout — closes remainder")
	exitPollInterval := flag.Duration("exit_poll_interval", 5*time.Second, "deadline check interval for timed exits")
	eventPostStartHold := flag.Duration("event_post_start_hold", 30*time.Minute, "event positions: minimum observation hold after scheduled start or in-play entry")
	timeoutReentryCooldown := flag.Duration("timeout_reentry_cooldown", 30*time.Minute, "block the same market after a timeout exit")
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
	paperCollectBroad := flag.Bool("paper_collect_broad", false, "paper only: collect otherwise-filtered wallet and market samples without extra push alerts")
	paperPromotedOnly := flag.Bool("paper_promoted_only", false, "paper only: count only locally promoted wallets as tradable; collect all others as research")
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
	liveDisableFile := flag.String("live_disable_file", "db/live-trading.disabled", "live trading kill-switch file; ignored in paper/research modes")
	liveArmFile := flag.String("live_arm_file", "db/live-trading.enabled", "short-lived, wallet-bound live trading arm file")
	liveMaxOrderUSD := flag.Float64("live_max_order_usd", 20.0, "hard maximum notional for one live BUY; exits are not capped")
	liveMaxSessionBuyUSD := flag.Float64("live_max_session_buy_usd", 100.0, "hard maximum reserved/filled BUY notional for one arm window")
	liveMaxArmDuration := flag.Duration("live_max_arm_duration", 24*time.Hour, "maximum accepted live arm-file validity window")
	initialCapital := flag.Float64("initial_capital", 200.0, "initial capital in USD for total P&L calculation")
	positionsStatePath := flag.String("positions_state", "db/positions.json", "paper/live position state JSON path")
	riskStatePath := flag.String("risk_state", "db/risk_state.json", "risk state JSON path")
	buyTimesStatePath := flag.String("buy_times_state", "db/buy_times.json", "buy-times state JSON path")
	posMaxTotalOpenUSD := flag.Float64("pos_max_total_open_usd", 300.0, "max total open paper exposure in USD")
	posMaxOpenPositions := flag.Int("pos_max_open_positions", 60, "max concurrent open paper positions")
	posMaxPerMarketUSD := flag.Float64("pos_max_per_market_usd", 30.0, "max open paper exposure per conditionID in USD; 0 disables")
	posMaxPerEventUSD := flag.Float64("pos_max_per_event_usd", 100.0, "max open paper exposure across related event conditions; 0 disables")
	footballScoreMaxEventUSD := flag.Float64("football_score_max_event_usd", 60.0, "max open paper exposure across exact-score outcomes for one match")
	flag.Parse()
	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	if *mode == "detect" {
		stateRoot := filepath.Join("db", "paper")
		if *liveTrading {
			stateRoot = filepath.Join("db", "live")
		}
		if !setFlags["positions_state"] {
			*positionsStatePath = filepath.Join(stateRoot, "positions.json")
		}
		if !setFlags["risk_state"] {
			*riskStatePath = filepath.Join(stateRoot, "risk_state.json")
		}
		if !setFlags["buy_times_state"] {
			*buyTimesStatePath = filepath.Join(stateRoot, "buy_times.json")
		}
		if !setFlags["journal_dir"] {
			*journalDir = filepath.Join(stateRoot, "journal")
		}
		if !setFlags["tickpath_dir"] {
			*tickPathDir = filepath.Join(stateRoot, "tickpath")
		}
		if err := validateTradingStatePaths(*liveTrading, *positionsStatePath, *riskStatePath, *buyTimesStatePath, *journalDir, *tickPathDir); err != nil {
			slog.Error("trading_state_isolation_rejected", "err", err)
			os.Exit(1)
		}
	} else if (*mode == "daily-report" || *mode == "daily-iterate") && !setFlags["journal_dir"] {
		*journalDir = filepath.Join("db", "paper", "journal")
	}

	if err := config.LoadDotEnv(*envFile); err != nil {
		slog.Warn("dotenv_load_warn", "path", *envFile, "err", err.Error())
	}
	if *liveTrading {
		if _, err := os.Stat(*liveDisableFile); err == nil {
			slog.Error("live_trading_disabled", "file", *liveDisableFile)
			os.Exit(1)
		} else if !os.IsNotExist(err) {
			slog.Error("live_disable_file_check_failed", "file", *liveDisableFile, "err", err.Error())
			os.Exit(1)
		}
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
			TP1Pct:               *ladderTP1Pct,
			TP1Frac:              *ladderTP1Frac,
			TP2Pct:               *ladderTP2Pct,
			TP2Frac:              *ladderTP2Frac,
			SLPct:                *ladderSLPct,
			TrailingPct:          *ladderTrailingPct,
			SLConfirmTime:        *ladderSLConfirm,
			MinHoldBeforeSL:      *ladderMinHoldBeforeSL,
			MaxHold:              *ladderMaxHold,
			FeeAware:             *ladderFeeAware,
			RequireExecutableBid: *ladderRequireBid,
			SlippageBp:           *slippageBp,
			FlatFeeBp:            *feeBp,
			TakerFeeRate:         *takerFeeRate,
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
		liveGuardCfg := order.LiveGuardConfig{
			ArmFile:          *liveArmFile,
			DisableFile:      *liveDisableFile,
			MaxOrderUSD:      *liveMaxOrderUSD,
			MaxSessionBuyUSD: *liveMaxSessionBuyUSD,
			MaxArmDuration:   *liveMaxArmDuration,
		}
		if err := runDetect(ctx, *maxMarkets, *windowSec, *slippageBp, *feeBp, *takerFeeRate, *largeFillUSD, *signalMode, *exitMode, *journalDir, *tickPathDir, *minEntry, *maxEntry, ladderCfg, *exitPollInterval, *eventPostStartHold, *timeoutReentryCooldown, *lotteryEnabled, lottCfg, injCfg, whaleCfg, *confirmDelay, btcCfg, updownCfg, p10, *liveTrading, liveGuardCfg, *fadeMode, *walletsFile, *copytradeSize, *walletTiersFile, *initialCapital, *minTier, *paperCollectBroad, *paperPromotedOnly, *positionsStatePath, *riskStatePath, *buyTimesStatePath, *posMaxTotalOpenUSD, *posMaxOpenPositions, *posMaxPerMarketUSD, *posMaxPerEventUSD, *footballScoreMaxEventUSD); err != nil && ctx.Err() == nil {
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
		slog.Error("arb-scan disabled", "reason", "third-party bookmaker odds API removed")
		os.Exit(1)
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

func copytradeTierForMarket(globalTier string, fileMeta walletFileMeta, footballScore bool) string {
	if fileMeta.List == "paper_promoted" && fileMeta.Tier != "" && fileMeta.Tier != "?" {
		return strings.ToUpper(fileMeta.Tier)
	}
	if footballScore && fileMeta.List == "football_score_push" && fileMeta.Tier != "" && fileMeta.Tier != "?" {
		return strings.ToUpper(fileMeta.Tier)
	}
	return strings.ToUpper(globalTier)
}

func copytradeAutoAllowedForMarket(action string, liveTrading, paperFollowPrompt, footballScoreEnabled, footballScore bool, tier string) (bool, string) {
	allowed, reason := copytradeAutoAllowedForAction(action, liveTrading, paperFollowPrompt)
	if allowed {
		return true, reason
	}
	if !liveTrading && footballScoreEnabled && footballScore && (tier == "A" || tier == "B") {
		return true, "football_score_" + strings.ToLower(tier)
	}
	return false, reason
}

func paperCollectionEnabled(requested bool, signalMode string, liveTrading bool) bool {
	return requested && signalMode == "copytrade" && !liveTrading
}

func paperPromotedOnlyEnabled(requested bool, signalMode string, liveTrading bool) bool {
	return requested && signalMode == "copytrade" && !liveTrading
}

func paperWalletPromoted(meta walletFileMeta) bool {
	return meta.List == "paper_promoted"
}

func paperPromotionAutoAllowed(meta walletFileMeta, promotedOnly, liveTrading bool) bool {
	return promotedOnly && !liveTrading && paperWalletPromoted(meta)
}

func riskEligibleSignalSource(source string) bool {
	source = strings.ToLower(strings.TrimSpace(source))
	return source != "manual" && !strings.HasPrefix(source, "copytrade_collect")
}

func tradeRecordNetPnL(record journal.TradeRecord) float64 {
	if record.NetPnLUSD == 0 && record.EntryFeeUSD == 0 && record.ExitFeeUSD == 0 {
		return record.PnLUSD
	}
	return record.NetPnLUSD
}

func durableRiskResults(records []journal.TradeRecord, mode string) []risk.RealizedResult {
	mode = strings.ToLower(strings.TrimSpace(mode))
	results := make([]risk.RealizedResult, 0, len(records))
	for _, record := range records {
		recordMode := strings.ToLower(strings.TrimSpace(record.Mode))
		if recordMode != "" && recordMode != mode {
			continue
		}
		if !riskEligibleSignalSource(record.SignalSource) {
			continue
		}
		at := record.ExitTime
		if at.IsZero() {
			at = record.EntryTime
		}
		results = append(results, risk.RealizedResult{PnLUSD: tradeRecordNetPnL(record), At: at})
	}
	return results
}

func copytradeTierAllowed(tier, minTier string) bool {
	tier = strings.ToUpper(strings.TrimSpace(tier))
	switch strings.ToUpper(strings.TrimSpace(minTier)) {
	case "":
		return true
	case "A":
		return tier == "A"
	case "B":
		return tier == "A" || tier == "B"
	case "C":
		return tier == "A" || tier == "B" || tier == "C"
	default:
		return false
	}
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

func monitorLiveGuard(ctx context.Context, cancel context.CancelFunc, client *order.GuardedClient) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := client.CheckReady(); err != nil {
				slog.Error("v2_live_guard_tripped", "err", err)
				cancel()
				return
			}
		}
	}
}

func monitorExecutionLedger(ctx context.Context, cancel context.CancelFunc, ledger *order.ExecutionLedger) {
	const staleAfter = 90 * time.Second
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			cutoff := now.Add(-staleAfter)
			pending, err := ledger.StaleUnresolvedCount("live", cutoff)
			if err != nil {
				slog.Error("order_ledger_monitor_failed", "err", err)
				cancel()
				return
			}
			unapplied, err := ledger.StaleUnappliedCount("live", cutoff)
			if err != nil {
				slog.Error("order_ledger_monitor_failed", "err", err)
				cancel()
				return
			}
			if pending > 0 || unapplied > 0 {
				slog.Error("order_ledger_reconciliation_required",
					"stale_pending", pending,
					"stale_unapplied", unapplied,
					"stale_after", staleAfter.String(),
					"action", "stopping live process")
				cancel()
				return
			}
		}
	}
}

type onChainReader interface {
	PUSDBalance(context.Context) (*big.Int, error)
	ConditionalTokenBalance(context.Context, string) (*big.Int, error)
}

type copytradeHoldPolicy struct {
	Name      string
	MaxHold   time.Duration
	EventHold time.Duration
}

func selectCopytradeHoldPolicy(text string, footballScore, paper bool, defaultMaxHold, defaultEventHold, esportsHold, footballScoreHold time.Duration) copytradeHoldPolicy {
	switch {
	case footballScore:
		return copytradeHoldPolicy{Name: "football_score", MaxHold: footballScoreHold, EventHold: footballScoreHold}
	case paper && feed.IsEsportsMarketText(text):
		return copytradeHoldPolicy{Name: "esports", MaxHold: esportsHold, EventHold: esportsHold}
	default:
		return copytradeHoldPolicy{Name: "default", MaxHold: defaultMaxHold, EventHold: defaultEventHold}
	}
}

func runDetect(ctx context.Context, topN, windowSec int, slippageBp, feeBp, takerFeeRate, largeFillUSD float64, signalMode, exitMode, journalDir, tickPathDir string, minEntry, maxEntry float64, ladderCfg strategy.LadderConfig, exitPollInterval, eventPostStartHold, timeoutReentryCooldown time.Duration, lotteryEnabled bool, lotteryCfg strategy.LotteryConfig, injCfg injury.Config, whaleCfg whale.Config, confirmDelay time.Duration, btcCfg btc.StrategyConfig, updownCfg btc.UpDownConfig, p10 phase10Config, liveTrading bool, liveGuardCfg order.LiveGuardConfig, fadeMode bool, walletsFile string, copytradeSize float64, walletTiersFile string, initialCapital float64, minTierFilter string, paperCollectBroad, paperPromotedOnly bool, positionsStatePath, riskStatePath, buyTimesStatePath string, posMaxTotalOpenUSD float64, posMaxOpenPositions int, posMaxPerMarketUSD, posMaxPerEventUSD, footballScoreMaxEventUSD float64) error {
	ctx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()
	if signalMode != "auto" && signalMode != "prompt" && signalMode != "whale" && signalMode != "copytrade" {
		return fmt.Errorf("invalid signal_mode %q (want auto|prompt|whale|copytrade)", signalMode)
	}
	if exitMode != "hold" && exitMode != "auto" && exitMode != "ladder" {
		return fmt.Errorf("invalid exit_mode %q (want hold|auto|ladder)", exitMode)
	}
	if exitPollInterval <= 0 {
		return fmt.Errorf("exit_poll_interval must be positive")
	}
	if eventPostStartHold < 0 || timeoutReentryCooldown < 0 || posMaxPerEventUSD < 0 || footballScoreMaxEventUSD < 0 {
		return fmt.Errorf("event hold and timeout cooldown must not be negative")
	}
	if ladderCfg.TP1Pct < 0 || ladderCfg.TP2Pct < 0 || ladderCfg.SLPct < 0 || ladderCfg.TrailingPct < 0 ||
		ladderCfg.TP1Frac <= 0 || ladderCfg.TP1Frac > 1 || ladderCfg.TP2Frac <= 0 || ladderCfg.TP2Frac > 1 ||
		ladderCfg.SLConfirmTime < 0 || ladderCfg.MinHoldBeforeSL < 0 || ladderCfg.MaxHold <= 0 ||
		ladderCfg.SlippageBp < 0 || ladderCfg.SlippageBp >= 10_000 || ladderCfg.FlatFeeBp < 0 ||
		ladderCfg.TakerFeeRate < 0 || ladderCfg.TakerFeeRate > 1 {
		return fmt.Errorf("invalid ladder thresholds, durations, or cost model")
	}
	paperCollectBroad = paperCollectionEnabled(paperCollectBroad, signalMode, liveTrading)
	paperFollowFootballScore := signalMode == "copytrade" && !liveTrading && os.Getenv("COPYTRADE_PAPER_FOLLOW_FOOTBALL_SCORE") == "1"
	paperFootballScoreSize := parseWhaleEnvFloat("COPYTRADE_PAPER_FOOTBALL_SCORE_SIZE", 5)
	if paperFootballScoreSize <= 0 {
		paperFootballScoreSize = 5
	}
	footballScoreMaxSignalAge := 2 * time.Minute
	if raw := strings.TrimSpace(os.Getenv("COPYTRADE_FOOTBALL_SCORE_MAX_SIGNAL_AGE")); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed >= 0 {
			footballScoreMaxSignalAge = parsed
		} else {
			slog.Warn("copytrade_football_invalid_signal_age", "value", raw, "err", err)
		}
	}
	footballScoreHold := 150 * time.Minute
	if raw := strings.TrimSpace(os.Getenv("COPYTRADE_FOOTBALL_SCORE_HOLD")); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			footballScoreHold = parsed
		} else {
			slog.Warn("copytrade_football_invalid_hold", "value", raw, "err", err)
		}
	}
	esportsHold := 45 * time.Minute
	if raw := strings.TrimSpace(os.Getenv("COPYTRADE_ESPORTS_HOLD")); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			esportsHold = parsed
		} else {
			slog.Warn("copytrade_esports_invalid_hold", "value", raw, "err", err)
		}
	}
	paperCollectBroad = paperCollectionEnabled(paperCollectBroad, signalMode, liveTrading)
	paperPromotedOnly = paperPromotedOnlyEnabled(paperPromotedOnly, signalMode, liveTrading)
	if signalMode == "copytrade" {
		slog.Info("copytrade_collection_config",
			"broad", paperCollectBroad,
			"promoted_only", paperPromotedOnly,
			"core_min_tier", minTierFilter,
			"live", liveTrading,
		)
		slog.Info("copytrade_football_config",
			"score_enabled", paperFollowFootballScore,
			"score_size_usd", paperFootballScoreSize,
			"score_max_signal_age", footballScoreMaxSignalAge.String(),
			"score_hold", footballScoreHold.String(),
			"score_max_event_usd", footballScoreMaxEventUSD,
			"live", liveTrading,
		)
		slog.Info("copytrade_hold_config",
			"default_max_hold", ladderCfg.MaxHold.String(),
			"default_event_hold", eventPostStartHold.String(),
			"paper_esports_hold", esportsHold.String(),
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
	copytradeTier := func(wallet string, footballScore bool) string {
		key := strings.ToLower(wallet)
		return copytradeTierForMarket(walletTiers[key], walletFileMetas[key], footballScore)
	}
	copytradeAutoAllowed := func(wallet string, footballScore bool) (bool, string) {
		key := strings.ToLower(wallet)
		if paperPromotionAutoAllowed(walletFileMetas[key], paperPromotedOnly, liveTrading) {
			return true, "paper_promoted"
		}
		meta := walletMetas[key]
		tier := copytradeTier(wallet, footballScore)
		return copytradeAutoAllowedForMarket(meta.FollowAction, liveTrading, os.Getenv("COPYTRADE_PAPER_FOLLOW_PROMPT") == "1", paperFollowFootballScore, footballScore, tier)
	}
	copytradeForWallet := func(wallet string) float64 {
		tier := walletTiers[strings.ToLower(wallet)]
		meta := walletMetas[strings.ToLower(wallet)]
		return copytradeWalletSize(copytradeSize, liveTrading, tier, meta.FollowAction, meta.SmartMoneyScore)
	}

	// hold & ladder both want the settlement watcher on — hold as primary,
	// ladder as safety net (a market resolving mid-tranche clears remainder).
	wantSettlement := exitMode == "hold" || exitMode == "ladder"
	jrn, err := journal.New(journalDir)
	if err != nil {
		return fmt.Errorf("journal init: %w", err)
	}
	jrn.SetPolicyVersion(os.Getenv("PAPER_POLICY_VERSION"))
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
	paperExecutionIntent := func(intent order.Intent, referencePrice float64) order.Intent {
		if liveTrading {
			return intent
		}
		intent.PaperReferencePx = referencePrice
		intent.PaperRequireQuote = true
		tail, ok := sampler.TickTail(intent.AssetID, 1)
		if !ok || len(tail) == 0 {
			return intent
		}
		quote := tail[0]
		intent.PaperBestBid = quote.BestBid
		intent.PaperBestBidSize = quote.BestBidSize
		intent.PaperBestAsk = quote.BestAsk
		intent.PaperBestAskSize = quote.BestAskSize
		intent.PaperQuoteAt = quote.QuoteTime
		return intent
	}

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
	shadowExitCfg := strategy.DefaultShadowExitConfig()
	shadowExitCfg.SlippageBp = slippageBp
	shadowExitCfg.FlatFeeBp = feeBp
	shadowExitCfg.TakerFeeRate = takerFeeRate
	shadowExits := strategy.NewShadowExitTracker(shadowExitCfg)
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
	posCfg.MaxPerEventUSD = posMaxPerEventUSD
	posCfg.PolicyVersion = strings.TrimSpace(os.Getenv("PAPER_POLICY_VERSION"))
	pm := strategy.NewPositionManager(posCfg)
	if positionsStatePath == "" {
		positionsStatePath = "db/positions.json"
	}
	if whaleCfg.StatePath == "" {
		whaleCfg.StatePath = filepath.Join(filepath.Dir(positionsStatePath), "whale-watermarks.json")
	}
	if err := os.MkdirAll(filepath.Dir(positionsStatePath), 0755); err != nil {
		return fmt.Errorf("positions state dir: %w", err)
	}
	type persistedAttribution struct {
		source      string
		openOrderID string
	}
	journalAttribution := make(map[string]persistedAttribution)
	journalExecutions := make(map[string]struct{})
	journalCloseOrders := make(map[string]struct{})
	var accounting []strategy.ClosedAccounting
	var persistedTrades []journal.TradeRecord
	if trades, readErr := journal.ReadAll(journalDir); readErr != nil {
		if liveTrading {
			return fmt.Errorf("read live trade journal: %w", readErr)
		}
		slog.Warn("positions_accounting_reconcile_read_err", "err", readErr)
	} else {
		persistedTrades = trades
		accounting = make([]strategy.ClosedAccounting, 0, len(trades))
		for _, tr := range trades {
			if tr.ExecutionID != "" {
				journalExecutions[tr.ExecutionID] = struct{}{}
			}
			if tr.CloseOrderID != "" {
				journalCloseOrders[tr.CloseOrderID] = struct{}{}
			}
			id := tr.ID
			if dot := strings.IndexByte(id, '.'); dot > 0 {
				id = id[:dot]
			}
			net := tr.NetPnLUSD
			if net == 0 && tr.EntryFeeUSD == 0 && tr.ExitFeeUSD == 0 {
				net = tr.PnLUSD
			}
			accounting = append(accounting, strategy.ClosedAccounting{
				ID: id, EntryTime: tr.EntryTime, ExitTime: tr.ExitTime,
				EntryFeeUSD: tr.EntryFeeUSD, ExitFeeUSD: tr.ExitFeeUSD, NetPnLUSD: net,
			})
			if tr.SignalSource != "" || tr.OpenOrderID != "" {
				journalAttribution[id] = persistedAttribution{source: tr.SignalSource, openOrderID: tr.OpenOrderID}
			}
		}
	}
	if _, statErr := os.Lstat(positionsStatePath); os.IsNotExist(statErr) && liveTrading {
		ledgerExists := false
		if _, ledgerErr := os.Lstat(filepath.Join(filepath.Dir(positionsStatePath), "orders.sqlite")); ledgerErr == nil {
			ledgerExists = true
		} else if !os.IsNotExist(ledgerErr) {
			return fmt.Errorf("inspect live execution ledger: %w", ledgerErr)
		}
		if ledgerExists || len(persistedTrades) > 0 {
			return fmt.Errorf("live position state is missing while durable trading history exists")
		}
		if err := pm.SaveState(positionsStatePath); err != nil {
			return fmt.Errorf("initialize live position state: %w", err)
		}
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("inspect positions state: %w", statErr)
	}
	if err := pm.LoadState(positionsStatePath); err != nil {
		if liveTrading {
			return fmt.Errorf("load live positions: %w", err)
		}
		slog.Warn("positions_load_err", "path", positionsStatePath, "err", err.Error())
	} else {
		closedUpdated := pm.ReconcileClosedAccounting(accounting)
		openUpdated := pm.ReconcileOpenEntryFees(accounting)
		if closedUpdated > 0 || openUpdated > 0 {
			slog.Info("positions_accounting_reconciled", "closed_updated", closedUpdated, "open_updated", openUpdated)
			if saveErr := pm.SaveState(positionsStatePath); saveErr != nil {
				if liveTrading {
					return fmt.Errorf("persist live accounting reconciliation: %w", saveErr)
				}
				slog.Warn("positions_accounting_reconcile_save_err", "err", saveErr)
			}
		}
		stats := pm.Stats()
		if liveTrading && stats.Closed > 0 && len(persistedTrades) == 0 {
			return fmt.Errorf("live position state contains closed trades but the trade journal is empty")
		}
		slog.Info("positions_loaded",
			"path", positionsStatePath,
			"open", stats.Open,
			"closed", stats.Closed,
			"exposure_usd", stats.TotalExposure,
			"tradable_open", pm.OpenCountForScope(strategy.ExposureScopeTradable),
			"tradable_exposure_usd", pm.ExposureForScope(strategy.ExposureScopeTradable),
			"collection_open", pm.OpenCountForScope(strategy.ExposureScopeCollection),
			"collection_exposure_usd", pm.ExposureForScope(strategy.ExposureScopeCollection),
		)
	}
	rehydratedAttribution := 0
	for _, p := range pm.Snapshot() {
		source := persistedPositionSource(p)
		openOrderID := p.OpenOrderID
		if attr, ok := journalAttribution[p.ID]; ok {
			if p.SignalSource == "" && attr.source != "" {
				source = attr.source
			}
			if openOrderID == "" {
				openOrderID = attr.openOrderID
			}
		}
		src.Mark(p.ID, source, openOrderID)
		if strings.HasPrefix(source, "copytrade") {
			shadowExits.Open(p)
			if exitMode == "ladder" {
				ladder.OpenPosition(p, takerFeeRate)
			}
			if added, err := ws.SubscribeAssets(p.AssetID); err != nil {
				slog.Warn("copytrade_wss_subscribe_fail", "pos", p.ID, "asset", short(p.AssetID), "err", err.Error())
			} else if added > 0 {
				slog.Info("copytrade_wss_subscribed", "pos", p.ID, "asset", short(p.AssetID), "phase", "rehydrate")
			}
		} else if source == "auto" {
			entryTick := feed.Tick{AssetID: p.AssetID, Market: p.Market, Time: p.EntryTime, Mid: p.EntryMid}
			switch exitMode {
			case "auto":
				exit.Open(p.AssetID, p.Market, entryTick)
			case "ladder":
				ladder.OpenWithDeadline(p.ID, p.Market, p.AssetID, entryTick, p.Units, p.ExitDeadline)
			}
		}
		if p.SignalSource == "" || (p.OpenOrderID == "" && openOrderID != "") {
			if err := pm.SetOpenAttribution(p.ID, source, openOrderID); err == nil {
				rehydratedAttribution++
			}
		}
	}
	if rehydratedAttribution > 0 {
		if err := pm.SaveState(positionsStatePath); err != nil {
			slog.Warn("positions_attribution_rehydrate_save_err", "err", err)
		} else {
			slog.Info("positions_attribution_rehydrated", "updated", rehydratedAttribution)
		}
	}
	if recorder != nil {
		rehydratedRecordings := 0
		for _, p := range pm.Snapshot() {
			if err := recorder.Start(p.ID, p.AssetID); err != nil {
				slog.Warn("tickrec_rehydrate_fail", "pos", p.ID, "asset", short(p.AssetID), "err", err.Error())
				continue
			}
			rehydratedRecordings++
		}
		if rehydratedRecordings > 0 {
			slog.Info("tickrec_rehydrated", "open_positions", rehydratedRecordings)
		}
	}
	savePositionsDurable := func() error {
		return pm.SaveState(positionsStatePath)
	}
	savePositions := func() {
		if err := savePositionsDurable(); err != nil {
			slog.Warn("positions_save_err", "path", positionsStatePath, "err", err)
		}
	}
	configureHoldWithPolicy := func(posID string, eventStart time.Time, maxHold, eventHold time.Duration) strategy.Position {
		planned, err := pm.ConfigureOpenHold(posID, eventStart, maxHold, eventHold)
		if err != nil {
			slog.Warn("position_hold_config_fail", "pos", posID, "err", err)
			return strategy.Position{ID: posID}
		}
		slog.Info("position_hold_configured",
			"pos", posID,
			"profile", planned.HoldProfile,
			"event_start", planned.EventStart,
			"exit_deadline", planned.ExitDeadline,
			"max_hold", maxHold.String(),
			"event_hold", eventHold.String(),
		)
		return planned
	}
	configureHold := func(posID string, eventStart time.Time) strategy.Position {
		return configureHoldWithPolicy(posID, eventStart, ladderCfg.MaxHold, eventPostStartHold)
	}
	var timeoutCooldownMu sync.Mutex
	timeoutCooldownUntil := make(map[string]time.Time)
	markTimeoutCooldown := func(market string, exitedAt time.Time) {
		market = strings.ToLower(strings.TrimSpace(market))
		if market == "" || timeoutReentryCooldown <= 0 {
			return
		}
		until := exitedAt.Add(timeoutReentryCooldown)
		timeoutCooldownMu.Lock()
		if until.After(timeoutCooldownUntil[market]) {
			timeoutCooldownUntil[market] = until
		}
		timeoutCooldownMu.Unlock()
	}
	timeoutBlockedUntil := func(market string, now time.Time) (time.Time, bool) {
		market = strings.ToLower(strings.TrimSpace(market))
		timeoutCooldownMu.Lock()
		defer timeoutCooldownMu.Unlock()
		until := timeoutCooldownUntil[market]
		return until, !until.IsZero() && now.Before(until)
	}
	timeoutLiquidity := strategy.NewTimeoutLiquidityTracker(5 * time.Minute)
	recordTimeoutLiquidity := func(p strategy.Position, mid float64, reason string, now time.Time) {
		// Keep extending the market-level block while the expired position has
		// no executable exit. This prevents stacking more exposure into a market
		// whose order book cannot currently absorb a sell.
		markTimeoutCooldown(p.Market, now)
		state, shouldLog := timeoutLiquidity.Observe(p, now)
		if !shouldLog {
			return
		}
		phase := "ongoing"
		if state.Attempts == 1 {
			phase = "detected"
		}
		slog.Warn("paper_timeout_best_bid_unavailable",
			"pos", p.ID,
			"asset", short(p.AssetID),
			"market", short(p.Market),
			"phase", phase,
			"reason", reason,
			"attempts", state.Attempts,
			"unavailable_sec", int(now.Sub(state.FirstSeen).Seconds()),
			"mid", mid,
			"exposure_usd", state.ExposureUSD,
			"conservative_value_usd", 0,
			"conservative_net_pnl_usd", state.ConservativeNetPnLUSD,
		)
	}
	logTimeoutLiquidityRecovered := func(p strategy.Position, exitFill float64, now time.Time) {
		state, ok := timeoutLiquidity.Resolve(p.ID)
		if !ok {
			return
		}
		slog.Info("paper_timeout_best_bid_recovered",
			"pos", p.ID,
			"asset", short(p.AssetID),
			"market", short(p.Market),
			"attempts", state.Attempts,
			"unavailable_sec", int(now.Sub(state.FirstSeen).Seconds()),
			"exit_fill", exitFill,
		)
	}
	for _, tr := range persistedTrades {
		if isTimeoutExitReason(tr.ExitReason) {
			markTimeoutCooldown(tr.Market, tr.ExitTime)
		}
	}

	if buyTimesStatePath == "" {
		buyTimesStatePath = "db/buy_times.json"
	}
	buyTimes, err := loadBuyTimeStore(buyTimesStatePath)
	if err != nil {
		return fmt.Errorf("buy-times state: %w", err)
	}
	slog.Info("buy_times_loaded", "count", buyTimes.Len())

	tradeMode := "paper"
	if liveTrading {
		tradeMode = "live"
	}
	ledgerPath := filepath.Join(filepath.Dir(positionsStatePath), "orders.sqlite")
	executionLedger, err := order.OpenExecutionLedger(ledgerPath)
	if err != nil {
		return fmt.Errorf("order ledger init: %w", err)
	}
	defer func() {
		if closeErr := executionLedger.Close(); closeErr != nil {
			slog.Warn("order_ledger_close_err", "err", closeErr)
		}
	}()
	if !liveTrading {
		if err := executionLedger.ResolveInterruptedPaper(); err != nil {
			return fmt.Errorf("resolve interrupted paper orders: %w", err)
		}
		removedPending := 0
		releasedCloses := 0
		for _, p := range pm.Snapshot() {
			if p.ClosingUnits > 0 {
				pm.AbortClose(p.ID)
				releasedCloses++
			}
			if strings.HasSuffix(p.SignalSource, "_pending") {
				if err := pm.CancelOpen(p.ID); err == nil {
					removedPending++
				}
			}
		}
		if removedPending > 0 || releasedCloses > 0 {
			if err := pm.SaveState(positionsStatePath); err != nil {
				return fmt.Errorf("remove interrupted paper reservations: %w", err)
			}
			if removedPending > 0 {
				slog.Warn("paper_pending_reservations_removed", "count", removedPending)
			}
			if releasedCloses > 0 {
				slog.Warn("paper_close_reservations_released", "count", releasedCloses)
			}
		}
	}

	var orderClient order.Client
	var v2Client *order.V2Client
	var walletAddress string
	var chainReader onChainReader
	pnlTrigger := make(chan struct{}, 4)
	paper := order.NewPaperClientWithFeeModel(slippageBp, feeBp, takerFeeRate)
	orderClient = paper
	if liveTrading {
		liveGuardCfg.SessionStatePath = filepath.Join(filepath.Dir(positionsStatePath), "live-session.json")
		slog.Info("v2_live_init", "msg", "loading wallet from Bitwarden")
		wallet, err := order.LoadWalletFromBitwarden("Polymarket-Go Wallet", "mnemonic", "")
		if err != nil {
			return fmt.Errorf("v2 wallet load: %w", err)
		}
		walletAddress = wallet.Address().Hex()
		slog.Info("v2_wallet_loaded", "address", walletAddress)
		liveGuardCfg.ExpectedWallet = walletAddress
		if err := order.CheckLiveGuard(liveGuardCfg); err != nil {
			return fmt.Errorf("v2 live guard rejected: %w", err)
		}
		oc, ocErr := order.NewReadOnlyOnChain("", wallet.Address())
		if ocErr != nil {
			slog.Warn("onchain_read_only_init_failed", "err", ocErr)
		} else {
			chainReader = oc
			slog.Info("onchain_read_only_ready", "address", walletAddress)
		}
		creds, err := order.DeriveExistingAPIKey(order.ClobBaseURL, wallet)
		if err != nil {
			return fmt.Errorf("v2 api key derive: %w", err)
		}
		slog.Info("v2_api_key_derived")
		v2Client = order.NewV2Client(wallet, creds, false)
		guardedClient, err := order.NewGuardedClient(v2Client, liveGuardCfg)
		if err != nil {
			return fmt.Errorf("v2 live guard init: %w", err)
		}
		orderClient = guardedClient
		slog.Info("v2_live_ready",
			"client", guardedClient.Name(),
			"exchange", order.V2ExchangeAddress,
			"max_order_usd", liveGuardCfg.MaxOrderUSD,
			"max_session_buy_usd", liveGuardCfg.MaxSessionBuyUSD,
			"arm_expires_within", liveGuardCfg.MaxArmDuration.String())
		go monitorLiveGuard(ctx, cancelRun, guardedClient)
		if err := executionLedger.ReconcileLive(ctx, v2Client); err != nil {
			slog.Warn("order_ledger_live_reconcile_partial", "err", err)
		}
		unresolved, err := executionLedger.UnresolvedCount(tradeMode)
		if err != nil {
			return fmt.Errorf("count unresolved live orders: %w", err)
		}
		if unresolved > 0 {
			return fmt.Errorf("refusing live startup: %d order executions remain unresolved", unresolved)
		}
		if os.Getenv("POLYMARKET_CANCEL_OPEN_ON_START") == "1" {
			if err := v2Client.CancelAllOpen(context.Background()); err != nil {
				slog.Warn("v2_cancel_all_open_failed", "err", err)
			}
		} else {
			slog.Info("v2_cancel_all_open_skipped", "reason", "POLYMARKET_CANCEL_OPEN_ON_START is not 1")
		}
	}
	ledgerClient, err := order.NewLedgerClient(orderClient, executionLedger, tradeMode)
	if err != nil {
		return err
	}
	orderClient = ledgerClient
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
	if _, statErr := os.Lstat(riskStatePath); os.IsNotExist(statErr) && liveTrading {
		stats := pm.Stats()
		if len(persistedTrades) > 0 || stats.Open > 0 || stats.Closed > 0 {
			return fmt.Errorf("live risk state is missing while trading history or positions exist")
		}
		if err := rm.SaveState(riskStatePath); err != nil {
			return fmt.Errorf("initialize live risk state: %w", err)
		}
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("inspect risk state: %w", statErr)
	}
	riskNow := time.Now()
	if err := rm.LoadState(riskStatePath, riskNow); err != nil {
		if liveTrading {
			return fmt.Errorf("load live risk state: %w", err)
		}
		slog.Warn("risk.load_state_failed", "err", err)
	}
	if liveTrading {
		st := rm.State()
		slog.Info("risk.state_loaded",
			"day", st.Day,
			"day_pnl", st.DayRealizedPnL,
			"cumulative_pnl", st.CumulativePnL,
			"blocked", st.Blocked,
			"block_reason", st.BlockReason,
		)
	} else {
		riskResults := durableRiskResults(persistedTrades, tradeMode)
		rm.RebuildRealized(riskResults, riskNow)
		if err := rm.SaveState(riskStatePath); err != nil {
			return fmt.Errorf("persist reconciled risk state: %w", err)
		}
		st := rm.State()
		slog.Info("risk.state_reconciled",
			"records", len(riskResults),
			"day", st.Day,
			"day_pnl", st.DayRealizedPnL,
			"cumulative_pnl", st.CumulativePnL,
			"blocked", st.Blocked,
			"block_reason", st.BlockReason,
		)
	}
	markExecutionApplied := func(result order.Result) {
		if result.ExecutionID == "" {
			return
		}
		if err := savePositionsDurable(); err != nil {
			slog.Error("order_ledger_position_not_durable", "execution_id", result.ExecutionID, "order_id", result.OrderID, "err", err)
			return
		}
		if err := rm.SaveState(riskStatePath); err != nil {
			slog.Error("order_ledger_risk_not_durable", "execution_id", result.ExecutionID, "order_id", result.OrderID, "err", err)
			return
		}
		if err := executionLedger.MarkApplied(result.ExecutionID); err != nil {
			slog.Error("order_ledger_mark_applied_fail", "execution_id", result.ExecutionID, "order_id", result.OrderID, "err", err)
		}
	}
	unappliedNonFills, err := executionLedger.UnappliedNonFills(tradeMode)
	if err != nil {
		return fmt.Errorf("load unapplied non-fill executions: %w", err)
	}
	reservationsChanged := false
	for _, record := range unappliedNonFills {
		switch record.Intent.Side {
		case order.Buy:
			if p, ok := pm.OpenByID(record.Intent.ClientID); ok && strings.HasSuffix(p.SignalSource, "_pending") {
				if err := pm.CancelOpen(p.ID); err != nil {
					return fmt.Errorf("release BUY %s reservation: %w", record.ID, err)
				}
				reservationsChanged = true
			}
		case order.Sell:
			if p, ok := pm.OpenByID(record.Intent.ClientID); ok && p.ClosingUnits > 0 {
				pm.AbortClose(p.ID)
				reservationsChanged = true
			}
		default:
			return fmt.Errorf("ledger execution %s has unsupported side %q", record.ID, record.Intent.Side)
		}
		if reservationsChanged {
			if err := savePositionsDurable(); err != nil {
				return fmt.Errorf("persist released execution %s reservation: %w", record.ID, err)
			}
			reservationsChanged = false
		}
		if err := executionLedger.MarkApplied(record.ID); err != nil {
			return fmt.Errorf("mark non-fill execution %s applied: %w", record.ID, err)
		}
		slog.Warn("order_ledger_nonfill_recovered", "execution_id", record.ID, "order_id", record.OrderID, "status", record.Status, "side", record.Intent.Side)
	}
	unapplied, err := executionLedger.UnappliedFills(tradeMode)
	if err != nil {
		return fmt.Errorf("load unapplied order fills: %w", err)
	}
	for _, record := range unapplied {
		result := record.Result
		if result.Status != order.StatusFilled {
			return fmt.Errorf("ledger execution %s is not a fill", record.ID)
		}
		switch record.Intent.Side {
		case order.Buy:
			filledAt := result.FilledAt
			if filledAt.IsZero() {
				filledAt = record.UpdatedAt
			}
			if result.AvgPrice <= 0 || result.AvgPrice >= 1 || result.FilledSize <= 0 {
				return fmt.Errorf("ledger BUY %s has invalid fill price/size", record.ID)
			}
			posID := strings.TrimSpace(record.Intent.ClientID)
			if posID == "" {
				posID = "recovered-" + record.ID[:12]
			}
			_, exists := pm.OpenByID(posID)
			if !exists && pm.HasOrderID(result.OrderID) {
				markExecutionApplied(result)
				continue
			}
			if exists {
				if err := pm.ApplyOpenFill(posID, result.AvgPrice, result.FilledSize, filledAt); err != nil {
					return fmt.Errorf("recover BUY %s fill: %w", record.ID, err)
				}
			} else {
				me := meta[record.Intent.AssetID]
				recovered := strategy.Position{
					ID: posID, AssetID: record.Intent.AssetID, Market: record.Intent.Market,
					Units: result.FilledSize, EntryMid: result.AvgPrice, EntryTime: filledAt,
					OpenFeeUSD: result.FeeUSD, OpenOrderID: result.OrderID,
					SignalSource: recoveredEntrySource(record.Intent.Reason, strategy.Position{}),
				}
				if me != nil {
					recovered.Question = me.Question
					recovered.Outcome = me.Outcome
				}
				var inserted bool
				_, inserted, recoverErr := pm.RecoverOpen(recovered)
				if recoverErr != nil {
					return fmt.Errorf("recover BUY %s position: %w", record.ID, recoverErr)
				}
				if !inserted {
					markExecutionApplied(result)
					continue
				}
			}
			p, _ := pm.OpenByID(posID)
			source := recoveredEntrySource(record.Intent.Reason, p)
			if err := pm.SetOpenFee(posID, result.FeeUSD); err != nil {
				return fmt.Errorf("recover BUY %s fee: %w", record.ID, err)
			}
			if err := pm.SetOpenAttribution(posID, source, result.OrderID); err != nil {
				return fmt.Errorf("recover BUY %s attribution: %w", record.ID, err)
			}
			eventStart := time.Time{}
			if me := meta[record.Intent.AssetID]; me != nil {
				eventStart = me.EventStart
			}
			planned, err := pm.ConfigureOpenHold(posID, eventStart, ladderCfg.MaxHold, eventPostStartHold)
			if err != nil {
				return fmt.Errorf("recover BUY %s hold plan: %w", record.ID, err)
			}
			if err := savePositionsDurable(); err != nil {
				return fmt.Errorf("persist recovered BUY %s: %w", record.ID, err)
			}
			src.Mark(posID, source, result.OrderID)
			entryTick := feed.Tick{AssetID: planned.AssetID, Market: planned.Market, Time: planned.EntryTime, Mid: planned.EntryMid}
			if strings.HasPrefix(source, "copytrade") {
				shadowExits.Open(planned)
				if exitMode == "ladder" {
					ladder.OpenPosition(planned, takerFeeRate)
				}
				if _, err := ws.SubscribeAssets(planned.AssetID); err != nil {
					slog.Warn("ledger_recover_wss_subscribe_fail", "pos", posID, "err", err)
				}
			} else if exitMode == "auto" {
				exit.Open(planned.AssetID, planned.Market, entryTick)
			} else if exitMode == "ladder" {
				ladder.OpenWithDeadline(planned.ID, planned.Market, planned.AssetID, entryTick, planned.Units, planned.ExitDeadline)
			}
			if recorder != nil {
				if err := recorder.Start(planned.ID, planned.AssetID); err != nil {
					slog.Warn("ledger_recover_tickrec_fail", "pos", posID, "err", err)
				}
			}
			if err := buyTimes.Set(planned.AssetID, planned.EntryTime); err != nil {
				return fmt.Errorf("persist recovered BUY %s time: %w", record.ID, err)
			}
			markExecutionApplied(result)
			slog.Warn("order_ledger_buy_recovered", "execution_id", record.ID, "order_id", result.OrderID, "pos", posID)

		case order.Sell:
			if _, ok := journalExecutions[record.ID]; ok {
				markExecutionApplied(result)
				continue
			}
			if _, ok := journalCloseOrders[result.OrderID]; ok {
				markExecutionApplied(result)
				continue
			}
			closed, alreadyClosed := pm.ClosedByExecution(record.ID)
			if !alreadyClosed {
				p, ok := pm.OpenByID(record.Intent.ClientID)
				if !ok {
					return fmt.Errorf("recover SELL %s: position %q is missing", record.ID, record.Intent.ClientID)
				}
				closeUnits := result.FilledSize
				if closeUnits <= 0 {
					closeUnits = record.Intent.SizeShares
				}
				if closeUnits > p.Units+1e-9 {
					return fmt.Errorf("recover SELL %s: filled %.8f shares exceed local position %.8f", record.ID, closeUnits, p.Units)
				}
				if p.ClosingUnits > 0 {
					if err := pm.ApplyCloseFill(p.ID, closeUnits); err != nil {
						return fmt.Errorf("recover SELL %s fill size: %w", record.ID, err)
					}
				} else if _, err := pm.BeginClose(p.ID, closeUnits); err != nil {
					return fmt.Errorf("recover SELL %s reservation: %w", record.ID, err)
				}
				exitAt := result.FilledAt
				if exitAt.IsZero() {
					exitAt = record.UpdatedAt
				}
				closed, err = pm.CommitClose(p.ID, strategy.ExitSignal{
					AssetID: p.AssetID, Market: p.Market, Time: exitAt,
					EntryMid: p.EntryMid, PeakMid: p.EntryMid, ExitMid: result.AvgPrice,
					HeldFor: exitAt.Sub(p.EntryTime), ChangePP: (result.AvgPrice - p.EntryMid) * 100,
					ExitFeeUSD: result.FeeUSD, Reason: strategy.ExitReason(record.Intent.Reason),
					CloseOrderID: result.OrderID, CloseExecutionID: record.ID,
				})
				if err != nil {
					return fmt.Errorf("recover SELL %s close: %w", record.ID, err)
				}
				if err := savePositionsDurable(); err != nil {
					return fmt.Errorf("persist recovered SELL %s: %w", record.ID, err)
				}
				if riskEligibleSignalSource(recoveredEntrySource("", closed)) {
					rm.OnClose(closed.NetPnLUSD, closed.ExitTime)
					if err := rm.SaveState(riskStatePath); err != nil {
						return fmt.Errorf("persist recovered SELL %s risk: %w", record.ID, err)
					}
				}
			}
			remaining, stillOpen := pm.OpenByID(closed.ID)
			me := meta[closed.AssetID]
			question, outcome := closed.Question, closed.Outcome
			if me != nil {
				if question == "" {
					question = me.Question
				}
				if outcome == "" {
					outcome = me.Outcome
				}
			}
			held := closed.ExitTime.Sub(closed.EntryTime)
			if held < 0 {
				held = 0
			}
			if err := jrn.Append(journal.TradeRecord{
				ExecutionID: record.ID, ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
				Question: question, Outcome: outcome, Side: "buy", SizeUSD: closed.SizeUSD, Units: closed.Units,
				EntryMid: closed.EntryMid, EntryTime: closed.EntryTime, ExitMid: closed.ExitMid, ExitTime: closed.ExitTime,
				ExitReason: string(closed.ExitReason), HeldSec: int(held.Seconds()), PnLUSD: closed.PnLUSD,
				EntryFeeUSD: closed.EntryFeeUSD, ExitFeeUSD: closed.ExitFeeUSD, NetPnLUSD: closed.NetPnLUSD,
				OpenOrderID: closed.OpenOrderID, CloseOrderID: result.OrderID, Mode: tradeMode,
				SignalSource:  recoveredEntrySource("", closed),
				PolicyVersion: closed.PolicyVersion,
			}); err != nil {
				return fmt.Errorf("journal recovered SELL %s: %w", record.ID, err)
			}
			journalExecutions[record.ID] = struct{}{}
			markExecutionApplied(result)
			if stillOpen {
				if ladder.Has(closed.ID) {
					ladder.Forget(closed.ID)
					ladder.SyncPosition(remaining)
				}
			} else {
				src.Take(closed.ID)
				ladder.Forget(closed.ID)
				shadowExits.ActualClose(closed)
				if recorder != nil {
					_ = recorder.Stop(closed.ID)
				}
			}
			slog.Warn("order_ledger_sell_recovered", "execution_id", record.ID, "order_id", result.OrderID, "pos", closed.ID, "remaining_open", stillOpen)
		default:
			return fmt.Errorf("ledger execution %s has unsupported side %q", record.ID, record.Intent.Side)
		}
	}
	if liveTrading {
		removedPending := 0
		releasedCloses := 0
		for _, p := range pm.Snapshot() {
			if p.ClosingUnits > 0 {
				pm.AbortClose(p.ID)
				releasedCloses++
			}
			if strings.HasSuffix(p.SignalSource, "_pending") {
				if err := pm.CancelOpen(p.ID); err == nil {
					removedPending++
				}
			}
		}
		if removedPending > 0 || releasedCloses > 0 {
			if err := savePositionsDurable(); err != nil {
				return fmt.Errorf("persist released live reservations: %w", err)
			}
			slog.Warn("live_orphan_reservations_released", "pending_buys", removedPending, "closing_positions", releasedCloses)
		}
	}
	if liveTrading {
		go monitorExecutionLedger(ctx, cancelRun, executionLedger)
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
	slog.Info("paper_client.ready", "slippage_bp", slippageBp, "fee_bp", feeBp, "taker_fee_rate", takerFeeRate, "per_pos_usd", posCfg.PerPositionUSD)
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
		"taker_fee_rate", takerFeeRate,
		"ladder_tp1_pct", ladderCfg.TP1Pct,
		"ladder_tp1_frac", ladderCfg.TP1Frac,
		"ladder_tp2_pct", ladderCfg.TP2Pct,
		"ladder_tp2_frac", ladderCfg.TP2Frac,
		"ladder_sl_pct", ladderCfg.SLPct,
		"ladder_trailing_pct", ladderCfg.TrailingPct,
		"ladder_sl_confirm", ladderCfg.SLConfirmTime.String(),
		"ladder_min_hold_before_sl", ladderCfg.MinHoldBeforeSL.String(),
		"ladder_fee_aware", ladderCfg.FeeAware,
		"ladder_require_bid", ladderCfg.RequireExecutableBid,
		"ladder_max_hold", ladderCfg.MaxHold.String(),
		"exit_poll_interval", exitPollInterval.String(),
		"event_post_start_hold", eventPostStartHold.String(),
		"timeout_reentry_cooldown", timeoutReentryCooldown.String(),
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
				pm:                   pm,
				exit:                 exit,
				ladder:               ladder,
				paper:                orderClient,
				rm:                   rm,
				pending:              pending,
				closePending:         closePending,
				notifier:             notifier,
				meta:                 meta,
				src:                  src,
				recorder:             recorder,
				jrn:                  jrn,
				largeFillUSD:         largeFillUSD,
				exitMode:             exitMode,
				holdMax:              ladderCfg.MaxHold,
				eventPostHold:        eventPostStartHold,
				riskStatePath:        riskStatePath,
				savePositions:        savePositions,
				savePositionsDurable: savePositionsDurable,
				markExecutionApplied: markExecutionApplied,
				prepareIntent:        paperExecutionIntent,
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
						pos, perr := pm.OldestByAsset(sig.AssetID)
						if perr != nil {
							exit.ConfirmClose(sig.AssetID)
							slog.Warn("paper_close_miss", "asset", short(sig.AssetID), "err", perr.Error())
							continue
						}
						if _, err := pm.BeginClose(pos.ID, pos.Units); err != nil {
							if !errors.Is(err, strategy.ErrPositionClosing) {
								exit.RetryClose(sig.AssetID)
							}
							slog.Warn("paper_close_reserve_fail", "pos", pos.ID, "err", err)
							continue
						}
						sellIntent := paperExecutionIntent(order.Intent{
							ClientID:   pos.ID,
							Reason:     string(sig.Reason),
							AssetID:    sig.AssetID,
							Market:     pos.Market,
							Side:       order.Sell,
							SizeUSD:    pos.Units * sig.ExitMid,
							SizeShares: pos.Units,
							LimitPx:    sig.ExitMid,
							Type:       order.FAK,
						}, sig.ExitMid)
						res, err := orderClient.Submit(ctx, sellIntent)
						if err != nil {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(pos.ID)
								exit.RetryClose(sig.AssetID)
							}
							slog.Warn("paper_sell_reject",
								"asset", short(sig.AssetID),
								"limit", sig.ExitMid,
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(pos.ID)
								exit.RetryClose(sig.AssetID)
							}
							slog.Warn("sell_not_filled",
								"asset", short(sig.AssetID),
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						if err := pm.ApplyCloseFill(pos.ID, res.FilledSize); err != nil {
							slog.Error("paper_close_fill_size_invalid", "pos", pos.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
							continue
						}
						sig.ExitMid = res.AvgPrice
						if !res.FilledAt.IsZero() {
							sig.Time = res.FilledAt
							sig.HeldFor = sig.Time.Sub(pos.EntryTime)
						}
						sig.ChangePP = (res.AvgPrice - sig.EntryMid) * 100
						sig.ExitFeeUSD = res.FeeUSD
						sig.CloseOrderID = res.OrderID
						sig.CloseExecutionID = res.ExecutionID
						closed, err := pm.CommitClose(pos.ID, sig)
						if err != nil {
							slog.Warn("paper_close_miss", "asset", short(sig.AssetID), "err", err.Error())
							continue
						}
						_, stillOpen := pm.OpenByID(closed.ID)
						exit.ConfirmClose(sig.AssetID)
						var nextAuto *strategy.Position
						for _, candidate := range pm.Snapshot() {
							if candidate.AssetID != sig.AssetID || persistedPositionSource(candidate) != "auto" || candidate.ClosingUnits > 0 {
								continue
							}
							if nextAuto == nil || candidate.EntryTime.Before(nextAuto.EntryTime) {
								copy := candidate
								nextAuto = &copy
							}
						}
						if nextAuto != nil {
							exit.Open(nextAuto.AssetID, nextAuto.Market, feed.Tick{
								AssetID: nextAuto.AssetID, Market: nextAuto.Market,
								Time: nextAuto.EntryTime, Mid: nextAuto.EntryMid,
							})
						}
						savePositions()
						if isTimeoutExitReason(string(closed.ExitReason)) {
							markTimeoutCooldown(closed.Market, closed.ExitTime)
						}
						if !stillOpen && recorder != nil {
							if rerr := recorder.Stop(closed.ID); rerr != nil {
								slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
							}
						}
						entryFee := closed.EntryFeeUSD
						exitFee := closed.ExitFeeUSD
						netPnL := closed.NetPnLUSD
						stats := pm.Stats()
						posSource, _ := src.Peek(closed.ID)
						if riskEligibleSignalSource(posSource) {
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
								SizeUSD:  closed.SizeUSD,
								PnLUSD:   netPnL,
								EntryPx:  sig.EntryMid,
								ExitPx:   res.AvgPrice,
								Reason:   string(sig.Reason),
								HeldSec:  int(sig.HeldFor.Seconds()),
							})
						}
						var source, openOID string
						if stillOpen {
							source, openOID = src.Peek(closed.ID)
						} else {
							source, openOID = src.Take(closed.ID)
						}
						if err := jrn.Append(journal.TradeRecord{
							ExecutionID: res.ExecutionID,
							ID:          closed.ID, AssetID: closed.AssetID, Market: closed.Market,
							Question:      metaQ(meta, closed.AssetID),
							Outcome:       metaOutcome(meta, closed.AssetID),
							Side:          "buy",
							SizeUSD:       closed.SizeUSD,
							Units:         closed.Units,
							EntryMid:      closed.EntryMid,
							EntryTime:     closed.EntryTime,
							ExitMid:       closed.ExitMid,
							ExitTime:      closed.ExitTime,
							ExitReason:    string(closed.ExitReason),
							HeldSec:       int(sig.HeldFor.Seconds()),
							PnLUSD:        closed.PnLUSD,
							EntryFeeUSD:   entryFee,
							ExitFeeUSD:    exitFee,
							NetPnLUSD:     netPnL,
							OpenOrderID:   openOID,
							CloseOrderID:  res.OrderID,
							Mode:          tradeMode,
							SignalSource:  source,
							PolicyVersion: closed.PolicyVersion,
						}); err != nil {
							slog.Warn("journal_append_fail", "asset", short(sig.AssetID), "err", err.Error())
						} else {
							markExecutionApplied(res)
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
						if _, err := pm.BeginClose(p.ID, ex.CloseUnits); err != nil {
							if !errors.Is(err, strategy.ErrPositionClosing) {
								ladder.Retry(p.ID)
							}
							slog.Warn("ladder_close_reserve_fail", "pos", p.ID, "err", err)
							continue
						}
						notional := ex.CloseUnits * ex.ExitMid
						feeRate := ex.TakerFeeRate
						sellIntent := paperExecutionIntent(order.Intent{
							ClientID:             p.ID,
							Reason:               string(ex.Reason),
							AssetID:              ex.AssetID,
							Market:               ex.Market,
							Side:                 order.Sell,
							SizeUSD:              notional,
							SizeShares:           ex.CloseUnits,
							LimitPx:              ex.ExitMid,
							Type:                 order.FAK,
							TakerFeeRateOverride: &feeRate,
						}, ex.ExitMid)
						res, err := orderClient.Submit(ctx, sellIntent)
						if err != nil {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(p.ID)
								ladder.Retry(p.ID)
							}
							slog.Warn("paper_ladder_sell_reject",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"limit", ex.ExitMid,
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(p.ID)
								ladder.Retry(p.ID)
							}
							slog.Warn("ladder_sell_not_filled",
								"pos", p.ID,
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						if err := pm.ApplyCloseFill(p.ID, res.FilledSize); err != nil {
							slog.Error("ladder_close_fill_size_invalid", "pos", p.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
							continue
						}
						ex.ExitMid = res.AvgPrice
						esig := strategy.ExitSignal{
							CloseExecutionID: res.ExecutionID,
							CloseOrderID:     res.OrderID,
							AssetID:          ex.AssetID,
							Market:           ex.Market,
							Time:             filledAtOr(res.FilledAt, ex.Time),
							EntryMid:         ex.EntryMid,
							PeakMid:          ex.ExitMid,
							ExitMid:          ex.ExitMid,
							HeldFor:          ex.HeldFor,
							ChangePP:         (ex.ExitMid - ex.EntryMid) * 100,
							ExitFeeUSD:       res.FeeUSD,
							Reason:           ex.Reason,
						}
						closedTranche, cerr := pm.CommitClose(p.ID, esig)
						if cerr != nil {
							slog.Warn("ladder_partial_close_fail",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"err", cerr.Error())
							continue
						}
						_, stillOpen := pm.OpenByID(p.ID)
						ladder.Confirm(p.ID, res.FilledSize)
						savePositions()
						if !stillOpen && isTimeoutExitReason(string(closedTranche.ExitReason)) {
							markTimeoutCooldown(closedTranche.Market, closedTranche.ExitTime)
						}
						if !stillOpen && recorder != nil {
							if rerr := recorder.Stop(p.ID); rerr != nil {
								slog.Warn("tickrec_stop_fail", "pos", p.ID, "err", rerr.Error())
							}
						}
						if !stillOpen {
							shadowExits.ActualClose(closedTranche)
						}
						entryFeeShare := closedTranche.EntryFeeUSD
						exitFee := closedTranche.ExitFeeUSD
						netPnL := closedTranche.NetPnLUSD
						stats := pm.Stats()
						ladderSource, _ := src.Peek(p.ID)
						if riskEligibleSignalSource(ladderSource) {
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
								SizeUSD:  closedTranche.SizeUSD,
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
						if !stillOpen {
							source, openOID = src.Take(p.ID)
						} else {
							source, openOID = src.Peek(p.ID)
						}
						trancheID := closedTranche.ID + "." + ex.Tranche
						if err := jrn.Append(journal.TradeRecord{
							ExecutionID:   res.ExecutionID,
							ID:            trancheID,
							AssetID:       closedTranche.AssetID,
							Market:        closedTranche.Market,
							Question:      metaQ(meta, closedTranche.AssetID),
							Outcome:       metaOutcome(meta, closedTranche.AssetID),
							Side:          "buy",
							SizeUSD:       closedTranche.SizeUSD,
							Units:         closedTranche.Units,
							EntryMid:      closedTranche.EntryMid,
							EntryTime:     closedTranche.EntryTime,
							ExitMid:       closedTranche.ExitMid,
							ExitTime:      closedTranche.ExitTime,
							ExitReason:    string(closedTranche.ExitReason),
							HeldSec:       int(ex.HeldFor.Seconds()),
							PnLUSD:        closedTranche.PnLUSD,
							EntryFeeUSD:   entryFeeShare,
							ExitFeeUSD:    exitFee,
							NetPnLUSD:     netPnL,
							Tranche:       ex.Tranche,
							OpenOrderID:   openOID,
							CloseOrderID:  res.OrderID,
							Mode:          tradeMode,
							SignalSource:  source,
							PolicyVersion: closedTranche.PolicyVersion,
						}); err != nil {
							slog.Warn("journal_append_fail",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"err", err.Error())
						} else {
							markExecutionApplied(res)
						}
						slog.Info("ladder_exit",
							"pos", p.ID,
							"asset", short(ex.AssetID),
							"q", metaQ(meta, ex.AssetID),
							"tranche", ex.Tranche,
							"reason", string(ex.Reason),
							"final", !stillOpen,
							"order_id", res.OrderID,
							"entry", ex.EntryMid,
							"exit_fill", res.AvgPrice,
							"close_units", res.FilledSize,
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

	// Phase 7.e: 1Hz tick-path recorder plus non-executing exit-policy
	// observations. Copytrade shadow policies run even if disk recording is
	// disabled; Recorder dedupes within-second on its side.
	if recorder != nil || signalMode == "copytrade" {
		go func() {
			tk := time.NewTicker(1 * time.Second)
			defer tk.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					targets := shadowExits.Snapshot()
					if recorder != nil {
						for posID, assetID := range recorder.Snapshot() {
							targets[posID] = assetID
						}
					}
					for posID, assetID := range targets {
						tail, ok := sampler.TickTail(assetID, 1)
						if !ok || len(tail) == 0 {
							continue
						}
						if recorder != nil {
							if err := recorder.Record(posID, tail[0]); err != nil {
								slog.Warn("tickrec_record_fail", "pos", posID, "err", err.Error())
							}
						}
						for _, obs := range shadowExits.OnTick(posID, tail[0]) {
							slog.Info("copytrade_exit_shadow",
								"pos", obs.PosID,
								"asset", short(obs.AssetID),
								"market", short(obs.Market),
								"policy", obs.Policy,
								"entry_time", obs.EntryTime,
								"entry", obs.EntryMid,
								"best_bid", obs.ExitQuotePrice,
								"executable_exit", obs.ExitPrice,
								"gross_return_pct", obs.GrossReturnPct,
								"gross_pnl_usd", obs.GrossPnLUSD,
								"entry_fee_usd", obs.EntryFeeUSD,
								"exit_fee_usd", obs.ExitFeeUSD,
								"net_pnl_usd", obs.NetPnLUSD,
								"net_return_pct", obs.NetReturnPct,
								"slippage_bp", obs.SlippageBp,
								"taker_fee_rate", obs.TakerFeeRate,
								"held_sec", int(obs.HeldFor.Seconds()),
								"hold_profile", obs.HoldProfile,
								"event_start", obs.EventStart,
								"actual_exit_deadline", obs.ExitDeadline,
								"question", obs.Question,
								"outcome", obs.Outcome,
								"source", obs.Source,
								"signal_source", obs.SignalSource,
								"actual_close_at", obs.ActualCloseAt,
								"actual_exit_reason", obs.ActualReason,
							)
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

				reserveTick := feed.Tick{
					AssetID: buyAssetID, Market: sig.Market,
					Time: time.Now(), Mid: buyMid,
				}
				pos, err := pm.OpenSized(buyAssetID, sig.Market, reserveTick, posCfg.PerPositionUSD)
				if err != nil {
					slog.Info("paper_open_skip",
						"asset", short(buyAssetID),
						"q", metaQ(meta, buyAssetID),
						"reason", err.Error(),
					)
					continue
				}
				buyIntent := order.Intent{
					ClientID: pos.ID,
					Reason:   "auto_entry",
					AssetID:  buyAssetID,
					Market:   sig.Market,
					Side:     order.Buy,
					SizeUSD:  posCfg.PerPositionUSD,
					LimitPx:  buyMid,
					Type:     order.FAK,
				}
				buyIntent = paperExecutionIntent(buyIntent, buyMid)
				res, err := orderClient.Submit(ctx, buyIntent)
				if err != nil {
					if res.Status == order.StatusPending {
						_ = pm.SetOpenAttribution(pos.ID, "auto_pending", res.OrderID)
						savePositions()
					} else {
						_ = pm.CancelOpen(pos.ID)
					}
					slog.Warn("paper_buy_reject",
						"asset", short(buyAssetID),
						"limit", buyMid,
						"err", err.Error())
					continue
				}
				if res.Status != order.StatusFilled {
					_ = pm.CancelOpen(pos.ID)
					slog.Warn("buy_not_filled",
						"asset", short(buyAssetID),
						"order_id", res.OrderID,
						"status", res.Status)
					continue
				}
				filledAt := res.FilledAt
				if filledAt.IsZero() {
					filledAt = time.Now()
				}
				if err := pm.ApplyOpenFill(pos.ID, res.AvgPrice, res.FilledSize, filledAt); err != nil {
					slog.Error("paper_apply_fill_fail", "pos", pos.ID, "order_id", res.OrderID, "err", err)
					continue
				}
				entryTick := feed.Tick{
					AssetID: buyAssetID, Market: sig.Market,
					Time: filledAt, Mid: res.AvgPrice,
				}
				if err := pm.SetOpenFee(pos.ID, res.FeeUSD); err != nil {
					slog.Warn("set_open_fee_err", "pos", pos.ID, "err", err)
				}
				markPositionSource(pm, src, pos.ID, "auto", res.OrderID)
				eventStart := time.Time{}
				if me := meta[buyAssetID]; me != nil {
					eventStart = me.EventStart
				}
				planned := configureHold(pos.ID, eventStart)
				if err := savePositionsDurable(); err != nil {
					slog.Error("positions_save_after_fill_fail", "execution_id", res.ExecutionID, "pos", pos.ID, "err", err)
				} else {
					markExecutionApplied(res)
				}
				switch exitMode {
				case "auto":
					exit.Open(buyAssetID, sig.Market, entryTick)
				case "ladder":
					ladder.OpenWithDeadline(pos.ID, sig.Market, buyAssetID, entryTick, pos.Units, planned.ExitDeadline)
				}
				if recorder != nil {
					if rerr := recorder.Start(pos.ID, buyAssetID); rerr != nil {
						slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
					}
				}
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
					"hold_profile", planned.HoldProfile,
					"exit_deadline", planned.ExitDeadline,
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
						reserveTick := feed.Tick{
							AssetID: c.AssetID, Market: c.Market,
							Time: time.Now(), Mid: c.Mid,
						}
						pos, err := pm.OpenSized(c.AssetID, c.Market, reserveTick, lotteryCfg.SizeUSD)
						if err != nil {
							slog.Info("lottery_open_skip",
								"asset", short(c.AssetID),
								"q", metaQ(meta, c.AssetID),
								"reason", err.Error())
							continue
						}
						buyIntent := order.Intent{
							ClientID: pos.ID,
							Reason:   "lottery_entry",
							AssetID:  c.AssetID,
							Market:   c.Market,
							Side:     order.Buy,
							SizeUSD:  lotteryCfg.SizeUSD,
							LimitPx:  c.Mid,
							Type:     order.FAK,
						}
						buyIntent = paperExecutionIntent(buyIntent, c.Mid)
						res, err := orderClient.Submit(ctx, buyIntent)
						if err != nil {
							if res.Status == order.StatusPending {
								_ = pm.SetOpenAttribution(pos.ID, "lottery_pending", res.OrderID)
								savePositions()
							} else {
								_ = pm.CancelOpen(pos.ID)
							}
							slog.Warn("lottery_buy_reject",
								"asset", short(c.AssetID),
								"mid", c.Mid,
								"sport", string(c.Sport),
								"err", err.Error())
							continue
						}
						if res.Status != order.StatusFilled {
							_ = pm.CancelOpen(pos.ID)
							slog.Warn("lottery_buy_not_filled",
								"asset", short(c.AssetID),
								"order_id", res.OrderID,
								"status", res.Status)
							continue
						}
						filledAt := res.FilledAt
						if filledAt.IsZero() {
							filledAt = time.Now()
						}
						if err := pm.ApplyOpenFill(pos.ID, res.AvgPrice, res.FilledSize, filledAt); err != nil {
							slog.Error("lottery_apply_fill_fail", "pos", pos.ID, "order_id", res.OrderID, "err", err)
							continue
						}
						if err := pm.SetOpenFee(pos.ID, res.FeeUSD); err != nil {
							slog.Warn("set_open_fee_err", "pos", pos.ID, "err", err)
						}
						markPositionSource(pm, src, pos.ID, "lottery", res.OrderID)
						eventStart := time.Time{}
						if me := meta[c.AssetID]; me != nil {
							eventStart = me.EventStart
						}
						planned := configureHold(pos.ID, eventStart)
						if err := savePositionsDurable(); err != nil {
							slog.Error("positions_save_after_fill_fail", "execution_id", res.ExecutionID, "pos", pos.ID, "err", err)
						} else {
							markExecutionApplied(res)
						}
						lotteryOpen[c.AssetID] = true
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
							"hold_profile", planned.HoldProfile,
							"exit_deadline", planned.ExitDeadline,
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
		EndDate      time.Time
		GameStart    time.Time
		Category     string
		NegRisk      bool
		TakerFeeRate float64
		FeeRateKnown bool
		LoadedAt     time.Time
	}
	var endDateMu sync.Mutex
	metaCache := make(map[string]marketMeta)
	lookupMarketMeta := func(conditionID string) (marketMeta, bool) {
		endDateMu.Lock()
		if m, ok := metaCache[conditionID]; ok && (!m.GameStart.IsZero() || time.Since(m.LoadedAt) < 30*time.Second) {
			endDateMu.Unlock()
			return m, true
		}
		endDateMu.Unlock()
		qctx, qcancel := context.WithTimeout(ctx, 10*time.Second)
		mkts, err := gc.GetByConditionIDs(qctx, []string{conditionID})
		qcancel()
		if err != nil || len(mkts) == 0 {
			return marketMeta{}, false
		}
		ed, _ := time.Parse(time.RFC3339, mkts[0].EndDate)
		gameStart := mkts[0].EventStartTime()
		sctx, scancel := context.WithTimeout(ctx, 3*time.Second)
		clobInfo, clobErr := gc.GetCLOBMarketInfo(sctx, conditionID)
		scancel()
		if clobErr != nil {
			slog.Debug("clob_market_info_miss", "market", short(conditionID), "err", clobErr)
		} else if gameStart.IsZero() {
			gameStart = clobInfo.GameStart
		}
		m := marketMeta{
			EndDate: ed, GameStart: gameStart, Category: mkts[0].Category, NegRisk: mkts[0].NegRisk,
			TakerFeeRate: clobInfo.TakerFeeRate, FeeRateKnown: clobInfo.FeeRateKnown, LoadedAt: time.Now(),
		}
		if gameStart.IsZero() {
			slog.Warn("market_start_unavailable", "market", short(conditionID), "gamma_game_start", mkts[0].GameStartTime, "clob_error", clobErr)
		}
		endDateMu.Lock()
		metaCache[conditionID] = m
		endDateMu.Unlock()
		return m, true
	}
	marketFeeRateOverride := func(conditionID string) *float64 {
		m, ok := lookupMarketMeta(conditionID)
		if !ok || !m.FeeRateKnown {
			return nil
		}
		rate := m.TakerFeeRate
		return &rate
	}
	reconfiguredHolds := 0
	for _, openPos := range pm.Snapshot() {
		if !strings.HasPrefix(persistedPositionSource(openPos), "copytrade") {
			continue
		}
		eventStart := time.Time{}
		if marketInfo, ok := lookupMarketMeta(openPos.Market); ok {
			eventStart = marketInfo.GameStart
		}
		policy := selectCopytradeHoldPolicy(openPos.Question, feed.IsFootballScoreMarketText(openPos.Question), !liveTrading,
			ladderCfg.MaxHold, eventPostStartHold, esportsHold, footballScoreHold)
		_, candidateDeadline := strategy.PlannedHold(openPos.EntryTime, eventStart, policy.MaxHold, policy.EventHold)
		if candidateDeadline.IsZero() || !candidateDeadline.After(openPos.ExitDeadline) {
			continue
		}
		planned := configureHoldWithPolicy(openPos.ID, eventStart, policy.MaxHold, policy.EventHold)
		shadowExits.Open(planned)
		if exitMode == "ladder" {
			feeRate := takerFeeRate
			if override := marketFeeRateOverride(openPos.Market); override != nil {
				feeRate = *override
			}
			ladder.OpenPosition(planned, feeRate)
		}
		slog.Info("position_hold_policy_reconfigured", "pos", openPos.ID, "policy", policy.Name, "exit_deadline", planned.ExitDeadline)
		reconfiguredHolds++
	}
	if reconfiguredHolds > 0 {
		savePositions()
		slog.Info("positions_hold_reconfigured", "count", reconfiguredHolds)
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
				if side == "BUY" && footballScore && footballScoreMaxSignalAge > 0 && !ev.Timestamp.IsZero() {
					signalAge := time.Since(ev.Timestamp)
					if signalAge < 0 {
						signalAge = 0
					}
					if signalAge > footballScoreMaxSignalAge {
						appendWhaleTrade(ev, "skip", "football_score_stale:"+signalAge.Round(time.Second).String())
						slog.Info("copytrade_football_score_stale",
							"wallet", ev.Label,
							"signal_age", signalAge.Round(time.Second).String(),
							"max_age", footballScoreMaxSignalAge.String(),
							"market", ev.Question,
						)
						return
					}
				}
				if ev.Price < copytradeEntryPriceFloor(footballScore, paperFollowFootballScore) || ev.Price > followMaxEntryPrice {
					appendWhaleTrade(ev, "skip", "price_filtered")
					slog.Info("copytrade_price_filtered", "wallet", ev.Label, "price", ev.Price, "market", ev.Question)
					return
				}

				switch side {
				case "BUY":
					collectionOnly := false
					walletKey := strings.ToLower(ev.Wallet)
					if paperPromotedOnly && !paperWalletPromoted(walletFileMetas[walletKey]) {
						if paperCollectBroad {
							collectionOnly = true
							slog.Info("copytrade_collection_gate_bypassed", "gate", "paper_promotion", "wallet", ev.Label, "market", ev.Question)
						} else {
							appendWhaleTrade(ev, "skip", "paper_not_promoted")
							slog.Info("copytrade_paper_promotion_filtered", "wallet", ev.Label, "market", ev.Question)
							return
						}
					}
					wTier := copytradeTier(ev.Wallet, footballScore)
					if !copytradeTierAllowed(wTier, minTierFilter) {
						if paperCollectBroad {
							collectionOnly = true
							slog.Info("copytrade_collection_gate_bypassed", "gate", "tier", "wallet", ev.Label, "tier", wTier, "market", ev.Question)
						} else {
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
					autoAllowed, followAction := copytradeAutoAllowed(ev.Wallet, footballScore)
					if !autoAllowed {
						if paperCollectBroad {
							collectionOnly = true
							slog.Info("copytrade_collection_gate_bypassed", "gate", "follow_action", "wallet", ev.Label, "action", followAction, "tier", wTier, "market", ev.Question)
						} else {
							appendWhaleTrade(ev, "skip", "follow_action:"+followAction)
							slog.Info("copytrade_follow_action_filtered",
								"wallet", ev.Label,
								"action", followAction,
								"tier", wTier,
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
					}
					if ok, filterReason := copytradeMarketDecision(ev.Question, ev.Slug, ev.Outcome, paperFollowFootballScore); !ok {
						if broadOK, _ := copytradeCollectionMarketDecision(ev.Question, ev.Slug, ev.Outcome, paperFollowFootballScore); paperCollectBroad && broadOK {
							collectionOnly = true
							slog.Info("copytrade_collection_gate_bypassed", "gate", "market", "wallet", ev.Label, "reason", filterReason, "market", ev.Question, "slug", ev.Slug)
						} else {
							appendWhaleTrade(ev, "skip", filterReason)
							slog.Info("copytrade_market_filtered", "wallet", ev.Label, "reason", filterReason, "market", ev.Question, "slug", ev.Slug)
							return
						}
					}
					if until, blocked := timeoutBlockedUntil(ev.ConditionID, time.Now()); blocked {
						appendWhaleTrade(ev, "skip", "timeout_reentry_cooldown:"+until.Format(time.RFC3339))
						slog.Info("copytrade_timeout_reentry_blocked",
							"wallet", ev.Label,
							"market", ev.Question,
							"condition_id", ev.ConditionID,
							"blocked_until", until,
						)
						return
					}
					// Core A/B entries keep their push path. Broad collection entries are
					// persisted and evaluated without adding notification noise.
					if !collectionOnly {
						aLabel := ev.Label
						if aLabel == "" {
							aLabel = ev.Wallet
						}
						meta := walletMetas[strings.ToLower(ev.Wallet)]
						baseAlert.Label = fmt.Sprintf("🔥💰 [自动模拟] %s (Tier %s · smart %.1f)", aLabel, wTier, meta.SmartMoneyScore)
						// Extra DM via sidecar bot so qualified alerts don't get buried.
						aTierMsg := fmt.Sprintf("🔥💰 Tier %s 智能钱下单\n%s · %s\n💰 %.0f shares @ %.4f = $%.0f\n🐋 %s\n%s",
							wTier,
							ev.Question, ev.Outcome,
							ev.SizeUnits, ev.Price, ev.Notional,
							baseAlert.Label, ev.LinkURL)
						notifier.SidecarAlert(aTierMsg)
					}
					isNegRisk := true
					followedMeta := marketMeta{}
					if marketInfo, ok := lookupMarketMeta(ev.ConditionID); ok {
						followedMeta = marketInfo
						isNegRisk = marketInfo.NegRisk
						if !marketInfo.EndDate.IsZero() && time.Until(marketInfo.EndDate) > 30*24*time.Hour {
							appendWhaleTrade(ev, "skip", fmt.Sprintf("settlement_too_far:%s", marketInfo.EndDate.Format("2006-01-02")))
							slog.Info("copytrade_settlement_filtered", "wallet", ev.Label, "market", ev.Question, "end_date", marketInfo.EndDate.Format("2006-01-02"))
							if !collectionOnly {
								notifier.WhaleAlert(baseAlert)
							}
							return
						}
					}
					if !collectionOnly {
						if err := rm.AllowOpen(time.Now()); err != nil {
							appendWhaleTrade(ev, "skip", "risk_blocked:"+err.Error())
							slog.Info("copytrade_blocked", "reason", err.Error(), "wallet", ev.Label, "market", ev.Question)
							return
						}
					}
					tick := feed.Tick{Mid: ev.Price, Time: time.Now()}
					sizeUSD := copytradeMarketSize(copytradeForWallet(ev.Wallet), footballScore, paperFollowFootballScore, paperFootballScoreSize)
					tier := wTier
					if tier == "" {
						tier = "?"
					}
					eventKey := copytradeExposureEventKey(ev)
					eventCap := 0.0
					if footballScore {
						eventCap = footballScoreMaxEventUSD
					}
					exposureScope := strategy.ExposureScopeTradable
					dedupeKey := copytradeCoreDedupeKey(ev)
					if collectionOnly {
						exposureScope = strategy.ExposureScopeCollection
					}
					pos, err := pm.OpenSizedForEventScopedGuarded(ev.AssetID, ev.ConditionID, eventKey, tick, sizeUSD, eventCap, exposureScope, dedupeKey)
					if err != nil {
						appendWhaleTrade(ev, "skip", "open_rejected:"+err.Error())
						slog.Warn("copytrade_open_rejected", "wallet", ev.Label, "asset", short(ev.AssetID), "err", err.Error())
						return
					}
					positionSource := "copytrade"
					if footballScore {
						positionSource = "copytrade_football_score"
					}
					if collectionOnly {
						positionSource = "copytrade_collect"
						if footballScore {
							positionSource = "copytrade_collect_football_score"
						}
					}
					if err := pm.SetOpenMetadata(pos.ID, ev.Question, ev.Outcome, positionSource, ev.Label); err != nil {
						_ = pm.CancelOpen(pos.ID)
						slog.Warn("copytrade_metadata_fail", "pos", pos.ID, "err", err)
						return
					}
					crossPx := ev.Price * 1.05
					if crossPx > 0.99 {
						crossPx = 0.99
					}
					orderPx := crossPx
					// Clear any stale orders before placing a new live order.
					if admin, ok := orderClient.(interface{ CancelAllOpen(context.Context) error }); ok {
						if err := admin.CancelAllOpen(ctx); err != nil {
							slog.Warn("copytrade_pre_cancel_err", "err", err)
						}
					}
					intent := order.Intent{
						ClientID:             pos.ID,
						Reason:               "copytrade_entry",
						AssetID:              ev.AssetID,
						Market:               ev.ConditionID,
						Side:                 order.Buy,
						SizeUSD:              sizeUSD,
						LimitPx:              orderPx,
						Type:                 order.FAK,
						NegRisk:              isNegRisk,
						TakerFeeRateOverride: marketFeeRateOverride(ev.ConditionID),
					}
					intent = paperExecutionIntent(intent, ev.Price)
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
						if !collectionOnly && !strings.Contains(reason, "not enough balance") {
							errMsg := fmt.Sprintf("❌ 跟单失败\n%s · %s\n💰 $%.0f @ %.4f · Tier %s\n🐋 %s\n⚠️ %s",
								ev.Question, ev.Outcome,
								sizeUSD, ev.Price, tier,
								ev.Label, reason)
							notifier.SidecarAlert(errMsg)
						}
						if result.Status == order.StatusPending {
							_ = pm.SetOpenAttribution(pos.ID, "copytrade_pending", result.OrderID)
						} else if cancelErr := pm.CancelOpen(pos.ID); cancelErr != nil {
							slog.Warn("copytrade_cancel_reservation_fail", "pos", pos.ID, "err", cancelErr)
						}
						savePositions()
						appendWhaleTrade(ev, "submit_failed", reason)
					} else {
						filledAt := result.FilledAt
						if filledAt.IsZero() {
							filledAt = time.Now()
						}
						if err := pm.ApplyOpenFill(pos.ID, result.AvgPrice, result.FilledSize, filledAt); err != nil {
							slog.Warn("copytrade_apply_fill_fail", "pos", pos.ID, "err", err.Error())
						}
						if err := pm.SetOpenFee(pos.ID, result.FeeUSD); err != nil {
							slog.Warn("copytrade_set_open_fee_fail", "pos", pos.ID, "err", err.Error())
						}
						effectiveFeeRate := takerFeeRate
						if intent.TakerFeeRateOverride != nil {
							effectiveFeeRate = *intent.TakerFeeRateOverride
						}
						slog.Info("copytrade_filled",
							"pos", pos.ID,
							"wallet", ev.Label,
							"tier", tier,
							"asset", short(ev.AssetID),
							"outcome", ev.Outcome,
							"signal_price", ev.Price,
							"fill_price", result.AvgPrice,
							"size_usd", sizeUSD,
							"order_id", result.OrderID,
							"fee_usd", result.FeeUSD,
							"taker_fee_rate", effectiveFeeRate,
							"collection_only", collectionOnly,
							"execution_model", result.ExecutionModel,
							"execution_reference", result.ReferencePrice,
							"quote_age_ms", result.QuoteAge.Milliseconds(),
						)
						fillMsg := fmt.Sprintf("✅ 跟单成功\n%s · %s\n💰 $%.0f @ %.4f · fee %.3fU · Tier %s\n🐋 %s 买入 $%.0f\n🆔 %s",
							ev.Question, ev.Outcome,
							sizeUSD, result.AvgPrice, result.FeeUSD, tier,
							ev.Label, ev.Notional,
							result.OrderID)
						if !collectionOnly {
							notifier.SidecarAlert(fillMsg)
						}
						signalSource := "copytrade_wallet:" + strings.ToLower(strings.TrimSpace(ev.Wallet))
						if footballScore {
							signalSource = "copytrade_football_score_wallet:" + strings.ToLower(strings.TrimSpace(ev.Wallet))
						}
						if collectionOnly {
							signalSource = "copytrade_collect_wallet:" + strings.ToLower(strings.TrimSpace(ev.Wallet))
							if footballScore {
								signalSource = "copytrade_collect_football_score_wallet:" + strings.ToLower(strings.TrimSpace(ev.Wallet))
							}
						}
						markPositionSource(pm, src, pos.ID, signalSource, result.OrderID)
						holdPolicy := selectCopytradeHoldPolicy(ev.Question+" "+ev.Slug+" "+followedMeta.Category, footballScore, !liveTrading,
							ladderCfg.MaxHold, eventPostStartHold, esportsHold, footballScoreHold)
						planned := configureHoldWithPolicy(pos.ID, followedMeta.GameStart, holdPolicy.MaxHold, holdPolicy.EventHold)
						shadowExits.OpenWithFeeRate(planned, effectiveFeeRate)
						if exitMode == "ladder" {
							ladder.OpenPosition(planned, effectiveFeeRate)
						}
						if added, subErr := ws.SubscribeAssets(ev.AssetID); subErr != nil {
							slog.Warn("copytrade_wss_subscribe_fail", "pos", pos.ID, "asset", short(ev.AssetID), "err", subErr.Error())
						} else if added > 0 {
							slog.Info("copytrade_wss_subscribed", "pos", pos.ID, "asset", short(ev.AssetID), "phase", "open")
						}
						if recorder != nil {
							if rerr := recorder.Start(pos.ID, ev.AssetID); rerr != nil {
								slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
							} else {
								entryTick := feed.Tick{
									AssetID: ev.AssetID,
									Market:  ev.ConditionID,
									Time:    filledAt,
									Mid:     result.AvgPrice,
								}
								if rerr := recorder.Record(pos.ID, entryTick); rerr != nil {
									slog.Warn("tickrec_entry_fail", "pos", pos.ID, "err", rerr.Error())
								}
							}
						}
						if err := savePositionsDurable(); err != nil {
							slog.Error("positions_save_after_fill_fail", "execution_id", result.ExecutionID, "pos", pos.ID, "err", err)
						} else {
							markExecutionApplied(result)
						}
						if err := buyTimes.Set(ev.AssetID, filledAt); err != nil {
							slog.Warn("buy_times_save_err", "asset", short(ev.AssetID), "err", err)
						}
						if collectionOnly {
							appendWhaleTrade(ev, "followed_collection", "paper_collect_broad")
						} else {
							appendWhaleTrade(ev, "followed", "")
						}
						slog.Info("copytrade_hold_plan",
							"pos", pos.ID,
							"policy", holdPolicy.Name,
							"profile", planned.HoldProfile,
							"event_start", planned.EventStart,
							"exit_deadline", planned.ExitDeadline,
						)
						select {
						case pnlTrigger <- struct{}{}:
						default:
						}
					}
					if !collectionOnly {
						notifier.WhaleAlert(baseAlert)
					}
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
					notifyWhaleSell := false
					var durableCloseResults []order.Result
					for _, pos := range matches {
						closeUnits := pos.Units * sellPct
						if _, err := pm.BeginClose(pos.ID, closeUnits); err != nil {
							slog.Warn("copytrade_close_reserve_fail", "pos", pos.ID, "err", err)
							continue
						}
						sellIntent := paperExecutionIntent(order.Intent{
							ClientID:             pos.ID,
							Reason:               "whale_sell",
							AssetID:              pos.AssetID,
							Market:               pos.Market,
							Side:                 order.Sell,
							SizeUSD:              closeUnits * ev.Price,
							SizeShares:           closeUnits,
							LimitPx:              ev.Price,
							Type:                 order.FAK,
							TakerFeeRateOverride: marketFeeRateOverride(pos.Market),
						}, ev.Price)
						res, serr := orderClient.Submit(ctx, sellIntent)
						if serr != nil || res.Status != order.StatusFilled {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(pos.ID)
							}
							slog.Warn("copytrade_sell_submit_fail", "pos", pos.ID, "err", serr, "status", res.Status, "detail", res.Error, "execution_model", res.ExecutionModel)
							continue
						}
						if err := pm.ApplyCloseFill(pos.ID, res.FilledSize); err != nil {
							slog.Error("copytrade_close_fill_size_invalid", "pos", pos.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
							continue
						}
						exitPrice := res.AvgPrice
						exitFee := res.FeeUSD
						closeOrderID := res.OrderID
						closedAt := filledAtOr(res.FilledAt, now)
						exitSig := strategy.ExitSignal{
							CloseExecutionID: res.ExecutionID,
							CloseOrderID:     res.OrderID,
							AssetID:          pos.AssetID,
							Market:           pos.Market,
							Time:             closedAt,
							EntryMid:         pos.EntryMid,
							ExitMid:          exitPrice,
							HeldFor:          closedAt.Sub(pos.EntryTime),
							ChangePP:         (exitPrice - pos.EntryMid) * 100,
							ExitFeeUSD:       exitFee,
							Reason:           strategy.ExitReason("whale_sell"),
						}
						closedPos, err := pm.CommitClose(pos.ID, exitSig)
						if err != nil {
							slog.Warn("copytrade_close_miss", "pos", pos.ID, "err", err.Error())
							continue
						}
						_, stillOpen := pm.OpenByID(pos.ID)
						if stillOpen {
							remaining, _ := pm.OpenByID(pos.ID)
							ladder.SyncPosition(remaining)
						} else {
							ladder.Forget(pos.ID)
						}
						entryFeeShare := closedPos.EntryFeeUSD
						exitFee = closedPos.ExitFeeUSD
						netPnL := closedPos.NetPnLUSD
						var source, openOID string
						if !stillOpen {
							source, openOID = src.Take(closedPos.ID)
							shadowExits.ActualClose(closedPos)
							if recorder != nil {
								if rerr := recorder.Stop(closedPos.ID); rerr != nil {
									slog.Warn("tickrec_stop_fail", "pos", closedPos.ID, "err", rerr.Error())
								}
							}
						} else {
							source, openOID = src.Peek(closedPos.ID)
						}
						if riskEligibleSignalSource(source) {
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
						}
						if jerr := jrn.Append(journal.TradeRecord{
							ExecutionID:   res.ExecutionID,
							ID:            closedPos.ID,
							AssetID:       closedPos.AssetID,
							Market:        closedPos.Market,
							Question:      ev.Question,
							Outcome:       ev.Outcome,
							Side:          "buy",
							SizeUSD:       closedPos.SizeUSD,
							Units:         closedPos.Units,
							EntryMid:      closedPos.EntryMid,
							EntryTime:     closedPos.EntryTime,
							ExitMid:       closedPos.ExitMid,
							ExitTime:      closedPos.ExitTime,
							ExitReason:    string(closedPos.ExitReason),
							HeldSec:       int(closedPos.ExitTime.Sub(closedPos.EntryTime).Seconds()),
							PnLUSD:        closedPos.PnLUSD,
							EntryFeeUSD:   entryFeeShare,
							ExitFeeUSD:    exitFee,
							NetPnLUSD:     netPnL,
							OpenOrderID:   openOID,
							CloseOrderID:  closeOrderID,
							Mode:          tradeMode,
							SignalSource:  source,
							PolicyVersion: closedPos.PolicyVersion,
						}); jerr != nil {
							slog.Warn("journal_append_fail", "asset", short(ev.AssetID), "err", jerr.Error())
						} else {
							durableCloseResults = append(durableCloseResults, res)
						}
						closed++
						slog.Info("copytrade_sell_executed",
							"pos", pos.ID,
							"pct_sold", fmt.Sprintf("%.0f%%", ev.PctSold),
							"close_units", closeUnits,
							"remaining", pos.Units-closeUnits,
							"gross_pnl", closedPos.PnLUSD,
							"exit_fee", exitFee,
							"net_pnl", netPnL,
						)
						pnlIcon := "🟢"
						if netPnL < 0 {
							pnlIcon = "🔴"
						}
						sellType := "全部平仓"
						if sellPct < 0.95 {
							sellType = fmt.Sprintf("部分平仓 %.0f%%", ev.PctSold)
						}
						sellMsg := fmt.Sprintf("%s 跟卖 · %s\n%s · %s\n💰 PnL $%+.2f · %.0f shares\n🐋 %s 卖出 %.0f%%",
							pnlIcon, sellType,
							ev.Question, ev.Outcome,
							netPnL, closeUnits,
							ev.Label, ev.PctSold)
						if !strings.HasPrefix(strings.ToLower(source), "copytrade_collect") {
							notifyWhaleSell = true
							notifier.SidecarAlert(sellMsg)
						}
					}
					if closed > 0 {
						if err := savePositionsDurable(); err != nil {
							slog.Error("positions_save_after_close_fail", "err", err)
						} else {
							for _, result := range durableCloseResults {
								markExecutionApplied(result)
							}
						}
						action := "closed"
						if sellPct < 0.95 {
							action = fmt.Sprintf("partial_%.0f%%", ev.PctSold)
						}
						appendWhaleTrade(ev, action, fmt.Sprintf("matched=%d", closed))
						slog.Info("copytrade_closed", "wallet", ev.Label, "asset", short(ev.AssetID), "count", closed, "pct", ev.PctSold)
					} else {
						appendWhaleTrade(ev, "sell_close_fail", "")
					}
					if notifyWhaleSell {
						notifier.WhaleAlert(baseAlert)
					}
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
					illiquid := timeoutLiquidity.Summary()
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
						"timeout_illiquid_positions", illiquid.Positions,
						"timeout_illiquid_exposure_usd", illiquid.ExposureUSD,
						"timeout_illiquid_conservative_net_pnl_usd", illiquid.ConservativeNetPnLUSD,
					)
				}
			}
		}
	}()

	// Settlement watcher (exit_mode=hold/ladder). Deadline checks run at the
	// configured short interval while ordinary Gamma settlement polling remains
	// every 60s; when a market is `closed` we close the position using
	// OutcomePrices[SlotIdx] as the final fill — 1.0 for the winning side,
	// 0.0 for the loser. Does the same risk/journal/notify bookkeeping the
	// auto-exit tracker does. SPEC §2 exit_mode=hold.
	if wantSettlement {
		go func() {
			tk := time.NewTicker(exitPollInterval)
			defer tk.Stop()
			lastHeldLog := time.Time{}
			lastMarketPoll := time.Time{}
			eventDeferredLogged := make(map[string]bool)
			marketRetryAfter := make(map[string]time.Time)
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
				now := time.Now()
				timeoutDue := false
				for _, p := range open {
					legacyDeadline := p.EntryTime.Add(ladderCfg.MaxHold)
					if p.HoldProfile == strategy.HoldProfileEvent && !eventDeferredLogged[p.ID] &&
						!now.Before(legacyDeadline) && now.Before(p.ExitDeadline) {
						eventDeferredLogged[p.ID] = true
						slog.Info("event_hold_deferred",
							"pos", p.ID,
							"market", short(p.Market),
							"event_start", p.EventStart,
							"exit_deadline", p.ExitDeadline,
						)
					}
					if !liveTrading && exitMode == "ladder" && strategy.PositionTimeoutDue(p, now, ladderCfg.MaxHold) {
						timeoutDue = true
					}
				}
				marketPollDue := lastMarketPoll.IsZero() || now.Sub(lastMarketPoll) >= time.Minute
				if !timeoutDue && !marketPollDue {
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
					if retryAt := marketRetryAfter[key]; retryAt.After(now) {
						continue
					}
					seen[key] = struct{}{}
					ids = append(ids, conditionID)
				}
				byCond := make(map[string]feed.Market)
				var mkts2 []feed.Market
				clobFallback := 0
				tokenFallback := 0
				if marketPollDue {
					lastMarketPoll = now
					if len(ids) > 0 {
						qctx, qcancel := context.WithTimeout(ctx, 5*time.Second)
						var err error
						mkts2, err = gc.GetByConditionIDs(qctx, ids)
						qcancel()
						if err != nil {
							slog.Warn("settlement_poll_fail", "err", err.Error(), "ids", len(ids))
							mkts2 = nil
						}
						for _, m := range mkts2 {
							key := strings.ToLower(m.ConditionID)
							byCond[key] = m
							delete(marketRetryAfter, key)
						}
					}
					missingIDs := make([]string, 0, len(ids)-len(byCond))
					for _, conditionID := range ids {
						if _, ok := byCond[strings.ToLower(conditionID)]; !ok {
							missingIDs = append(missingIDs, conditionID)
						}
					}
					if len(missingIDs) > 0 {
						tokenIDs := settlementFallbackTokenIDs(open, missingIDs)
						if len(tokenIDs) > 0 {
							tokenCtx, tokenCancel := context.WithTimeout(ctx, 5*time.Second)
							tokenMarkets, tokenErr := gc.GetByClobTokenIDs(tokenCtx, tokenIDs)
							tokenCancel()
							if tokenErr != nil {
								slog.Warn("settlement_token_fallback_fail", "tokens", len(tokenIDs), "err", tokenErr)
							} else {
								tokenFallback = indexSettlementMarketsByToken(open, missingIDs, tokenMarkets, byCond)
							}
							missingIDs = missingSettlementConditionIDs(ids, byCond)
						}
					}
					if len(missingIDs) > 0 {
						fallbackCtx, fallbackCancel := context.WithTimeout(ctx, 20*time.Second)
						fallbackMarkets, fallbackErr := gc.GetCLOBMarkets(fallbackCtx, missingIDs)
						fallbackCancel()
						for _, market := range fallbackMarkets {
							key := strings.ToLower(market.ConditionID)
							byCond[key] = market
							delete(marketRetryAfter, key)
						}
						clobFallback = len(fallbackMarkets)
						if fallbackErr != nil {
							slog.Warn("settlement_clob_fallback_partial", "requested", len(missingIDs), "returned", clobFallback, "err", fallbackErr)
						}
						backoffUntil := now.Add(15 * time.Minute)
						backedOff := 0
						for _, conditionID := range missingIDs {
							key := strings.ToLower(conditionID)
							if _, found := byCond[key]; found {
								continue
							}
							marketRetryAfter[key] = backoffUntil
							backedOff++
						}
						if backedOff > 0 {
							slog.Warn("settlement_market_backoff", "markets", backedOff, "retry_after", backoffUntil)
						}
					}
				}
				// Periodic "still holding" log (once per 5 min) — easy to grep for.
				if now.Sub(lastHeldLog) >= 5*time.Minute {
					lastHeldLog = now
					slog.Info("hold_status",
						"open", len(open),
						"markets_polled", len(ids),
						"markets_returned", len(byCond),
						"gamma_returned", len(mkts2),
						"token_fallback_returned", tokenFallback,
						"clob_fallback_returned", clobFallback,
						"resolved_seen", countResolvedMap(byCond),
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
						if !liveTrading && exitMode == "ladder" && strategy.PositionTimeoutDue(p, now, ladderCfg.MaxHold) {
							if p.ExitDeadline.IsZero() {
								slog.Warn("legacy_hold_metadata_unavailable", "pos", p.ID, "asset", short(p.AssetID), "market", short(p.Market))
								continue
							}
							tail, tickOK := sampler.TickTail(p.AssetID, 1)
							if !tickOK || len(tail) == 0 {
								recordTimeoutLiquidity(p, 0, "tick_unavailable", now)
								continue
							}
							timeoutPrice, priceOK := paperTimeoutExitPrice(tail[0], now)
							if !priceOK {
								recordTimeoutLiquidity(p, tail[0].Mid, "no_executable_bid", now)
								continue
							}
							if _, reserveErr := pm.BeginClose(p.ID, p.Units); reserveErr != nil {
								slog.Warn("paper_timeout_flat_reserve_fail", "pos", p.ID, "err", reserveErr)
								continue
							}
							timeoutIntent := paperExecutionIntent(order.Intent{
								ClientID:             p.ID,
								Reason:               string(strategy.ExitTimeout),
								AssetID:              p.AssetID,
								Market:               p.Market,
								Side:                 order.Sell,
								SizeUSD:              p.Units * timeoutPrice,
								SizeShares:           p.Units,
								LimitPx:              timeoutPrice,
								Type:                 order.FAK,
								TakerFeeRateOverride: marketFeeRateOverride(p.Market),
							}, timeoutPrice)
							res, serr := orderClient.Submit(ctx, timeoutIntent)
							if serr != nil || res.Status != order.StatusFilled {
								if res.Status == order.StatusPending {
									savePositions()
								} else {
									pm.AbortClose(p.ID)
								}
								slog.Warn("paper_timeout_flat_sell_fail", "pos", p.ID, "asset", short(p.AssetID), "err", serr, "status", res.Status)
								continue
							}
							if err := pm.ApplyCloseFill(p.ID, res.FilledSize); err != nil {
								slog.Error("paper_timeout_flat_fill_size_invalid", "pos", p.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
								continue
							}
							sig := strategy.ExitSignal{
								CloseExecutionID: res.ExecutionID,
								CloseOrderID:     res.OrderID,
								AssetID:          p.AssetID,
								Market:           p.Market,
								Time:             filledAtOr(res.FilledAt, now),
								EntryMid:         p.EntryMid,
								PeakMid:          p.EntryMid,
								ExitMid:          res.AvgPrice,
								HeldFor:          now.Sub(p.EntryTime),
								ExitFeeUSD:       res.FeeUSD,
								Reason:           strategy.ExitTimeout,
							}
							closed, cerr := pm.CommitClose(p.ID, sig)
							if cerr != nil {
								slog.Warn("paper_timeout_flat_close_miss", "pos", p.ID, "asset", short(p.AssetID), "err", cerr.Error())
								continue
							}
							remaining, stillOpen := pm.OpenByID(p.ID)
							logTimeoutLiquidityRecovered(p, res.AvgPrice, now)
							savePositions()
							markTimeoutCooldown(closed.Market, closed.ExitTime)
							hadLadder := ladder.Has(p.ID)
							ladder.Forget(p.ID)
							if stillOpen && hadLadder {
								ladder.SyncPosition(remaining)
							}
							if !stillOpen {
								shadowExits.ActualClose(closed)
							}
							if !stillOpen && recorder != nil {
								if rerr := recorder.Stop(closed.ID); rerr != nil {
									slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
								}
							}
							entryFeeShare := closed.EntryFeeUSD
							netPnL := closed.NetPnLUSD
							flatSource, _ := src.Peek(closed.ID)
							if riskEligibleSignalSource(flatSource) {
								rm.OnClose(netPnL, now)
								if err := rm.SaveState(riskStatePath); err != nil {
									slog.Warn("risk_save_err", "err", err)
								}
							}
							var source, openOID string
							if stillOpen {
								source, openOID = src.Peek(closed.ID)
							} else {
								source, openOID = src.Take(closed.ID)
							}
							if jerr := jrn.Append(journal.TradeRecord{
								ExecutionID: res.ExecutionID,
								ID:          closed.ID, AssetID: closed.AssetID, Market: closed.Market,
								Question:      closed.Question,
								Outcome:       closed.Outcome,
								Side:          "buy",
								SizeUSD:       closed.SizeUSD,
								Units:         closed.Units,
								EntryMid:      closed.EntryMid,
								EntryTime:     closed.EntryTime,
								ExitMid:       closed.ExitMid,
								ExitTime:      closed.ExitTime,
								ExitReason:    string(closed.ExitReason),
								HeldSec:       int(sig.HeldFor.Seconds()),
								PnLUSD:        closed.PnLUSD,
								EntryFeeUSD:   entryFeeShare,
								ExitFeeUSD:    closed.ExitFeeUSD,
								NetPnLUSD:     netPnL,
								Tranche:       "timeout_flat",
								OpenOrderID:   openOID,
								CloseOrderID:  res.OrderID,
								Mode:          tradeMode,
								SignalSource:  source,
								PolicyVersion: closed.PolicyVersion,
							}); jerr != nil {
								slog.Warn("journal_append_fail", "asset", short(p.AssetID), "err", jerr.Error())
							} else {
								markExecutionApplied(res)
							}
							stats := pm.Stats()
							slog.Info("paper_timeout_flat_exit",
								"pos", closed.ID,
								"asset", short(p.AssetID),
								"q", closed.Question,
								"outcome", closed.Outcome,
								"entry", p.EntryMid,
								"exit_fill", res.AvgPrice,
								"entry_fee_usd", entryFeeShare,
								"exit_fee_usd", res.FeeUSD,
								"net_pnl_usd", netPnL,
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
					if p.ExitDeadline.IsZero() {
						eventStart := m.EventStartTime()
						if eventStart.IsZero() {
							if marketInfo, found := lookupMarketMeta(conditionID); found {
								eventStart = marketInfo.GameStart
							}
						}
						planned := configureHold(p.ID, eventStart)
						p = planned
						shadowExits.Open(planned)
						savePositions()
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
						if liveTrading || exitMode != "ladder" || !strategy.PositionTimeoutDue(p, now, ladderCfg.MaxHold) {
							continue
						}
						reason = strategy.ExitTimeout
					}
					exitMid := settleMid
					exitFee := 0.0
					executableBid := 0.0
					orderID := fmt.Sprintf("settle-%s", short(p.AssetID))
					var executionResult order.Result
					tranche := "settle"
					if reason == strategy.ExitTimeout {
						tail, tickOK := sampler.TickTail(p.AssetID, 1)
						if !tickOK || len(tail) == 0 {
							recordTimeoutLiquidity(p, 0, "tick_unavailable", now)
							continue
						}
						var priceOK bool
						executableBid, priceOK = paperTimeoutExitPrice(tail[0], now)
						if !priceOK {
							recordTimeoutLiquidity(p, tail[0].Mid, "no_executable_bid", now)
							continue
						}
						if _, reserveErr := pm.BeginClose(p.ID, p.Units); reserveErr != nil {
							slog.Warn("paper_timeout_reserve_fail", "pos", p.ID, "err", reserveErr)
							continue
						}
						timeoutIntent := paperExecutionIntent(order.Intent{
							ClientID:             p.ID,
							Reason:               string(strategy.ExitTimeout),
							AssetID:              p.AssetID,
							Market:               p.Market,
							Side:                 order.Sell,
							SizeUSD:              p.Units * executableBid,
							SizeShares:           p.Units,
							LimitPx:              executableBid,
							Type:                 order.FAK,
							TakerFeeRateOverride: marketFeeRateOverride(p.Market),
						}, executableBid)
						res, serr := orderClient.Submit(ctx, timeoutIntent)
						if serr != nil || res.Status != order.StatusFilled {
							if res.Status == order.StatusPending {
								savePositions()
							} else {
								pm.AbortClose(p.ID)
							}
							slog.Warn("paper_timeout_sell_fail", "pos", p.ID, "asset", short(p.AssetID), "err", serr, "status", res.Status)
							continue
						}
						if err := pm.ApplyCloseFill(p.ID, res.FilledSize); err != nil {
							slog.Error("paper_timeout_fill_size_invalid", "pos", p.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
							continue
						}
						exitMid = res.AvgPrice
						exitFee = res.FeeUSD
						orderID = res.OrderID
						executionResult = res
						tranche = "timeout"
					}
					if reason == strategy.ExitSettlement {
						if _, reserveErr := pm.BeginClose(p.ID, p.Units); reserveErr != nil {
							slog.Warn("settlement_reserve_fail", "pos", p.ID, "err", reserveErr)
							continue
						}
					}
					sig := strategy.ExitSignal{
						CloseExecutionID: executionResult.ExecutionID,
						CloseOrderID:     orderID,
						AssetID:          p.AssetID,
						Market:           p.Market,
						Time:             filledAtOr(executionResult.FilledAt, now),
						EntryMid:         p.EntryMid,
						PeakMid:          p.EntryMid,
						ExitMid:          exitMid,
						HeldFor:          now.Sub(p.EntryTime),
						ChangePP:         (exitMid - p.EntryMid) * 100,
						ExitFeeUSD:       exitFee,
						Reason:           reason,
					}
					closed, cerr := pm.CommitClose(p.ID, sig)
					if cerr != nil {
						slog.Warn("settlement_close_miss", "pos", p.ID, "asset", short(p.AssetID), "err", cerr.Error())
						continue
					}
					remaining, stillOpen := pm.OpenByID(p.ID)
					if reason == strategy.ExitTimeout {
						logTimeoutLiquidityRecovered(p, exitMid, now)
					}
					savePositions()
					if !stillOpen {
						shadowExits.ActualClose(closed)
					}
					if reason == strategy.ExitTimeout {
						markTimeoutCooldown(closed.Market, closed.ExitTime)
					}
					// Drop any ladder state that was still tracking this
					// position — settlement supersedes TP/SL/timeout.
					hadLadder := ladder.Has(p.ID)
					ladder.Forget(p.ID)
					if stillOpen && hadLadder {
						ladder.SyncPosition(remaining)
					}
					if !stillOpen && recorder != nil {
						if rerr := recorder.Stop(closed.ID); rerr != nil {
							slog.Warn("tickrec_stop_fail", "pos", closed.ID, "err", rerr.Error())
						}
					}
					entryFeeShare := closed.EntryFeeUSD
					exitFee = closed.ExitFeeUSD
					netPnL := closed.NetPnLUSD
					stats := pm.Stats()
					settleSource, _ := src.Peek(closed.ID)
					if riskEligibleSignalSource(settleSource) {
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
							SizeUSD:  closed.SizeUSD,
							PnLUSD:   netPnL,
							EntryPx:  p.EntryMid,
							ExitPx:   exitMid,
							Reason:   string(reason),
							HeldSec:  int(sig.HeldFor.Seconds()),
						})
					}
					var source, openOID string
					if stillOpen {
						source, openOID = src.Peek(closed.ID)
					} else {
						source, openOID = src.Take(closed.ID)
					}
					if jerr := jrn.Append(journal.TradeRecord{
						ExecutionID: executionResult.ExecutionID,
						ID:          closed.ID, AssetID: closed.AssetID, Market: closed.Market,
						Question:      paperPositionQuestion(closed, m, meta[closed.AssetID]),
						Outcome:       paperPositionOutcome(closed, m, slotIdx, meta[closed.AssetID]),
						Side:          "buy",
						SizeUSD:       closed.SizeUSD,
						Units:         closed.Units,
						EntryMid:      closed.EntryMid,
						EntryTime:     closed.EntryTime,
						ExitMid:       closed.ExitMid,
						ExitTime:      closed.ExitTime,
						ExitReason:    string(closed.ExitReason),
						HeldSec:       int(sig.HeldFor.Seconds()),
						PnLUSD:        closed.PnLUSD,
						EntryFeeUSD:   entryFeeShare,
						ExitFeeUSD:    exitFee,
						NetPnLUSD:     netPnL,
						Tranche:       tranche,
						OpenOrderID:   openOID,
						CloseOrderID:  orderID,
						Mode:          tradeMode,
						SignalSource:  source,
						PolicyVersion: closed.PolicyVersion,
					}); jerr != nil {
						slog.Warn("journal_append_fail", "asset", short(p.AssetID), "err", jerr.Error())
					} else {
						markExecutionApplied(executionResult)
					}
					slog.Info("settlement_exit",
						"asset", short(p.AssetID),
						"q", paperPositionQuestion(p, m, meta[p.AssetID]),
						"outcome", paperPositionOutcome(p, m, slotIdx, meta[p.AssetID]),
						"entry", p.EntryMid,
						"market_mid", settleMid,
						"executable_bid", executableBid,
						"exit_fill", exitMid,
						"gross_pnl_usd", closed.PnLUSD,
						"entry_fee_usd", entryFeeShare,
						"exit_fee_usd", exitFee,
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
			fmt.Fprintf(&sb, "📊 P&L · %s SGT\n\n", time.Now().In(sgt).Format("15:04"))

			if walletAddress == "" {
				return
			}

			positions, err := fetchDataAPIPositions(ctx, walletAddress)
			if err != nil {
				fmt.Fprintf(&sb, "⚠️ 拉取仓位失败: %s\n", err)
				notifier.SidecarAlert(sb.String())
				slog.Warn("pnl_push_fetch_err", "err", err)
				return
			}

			// --- Query wallet balance ---
			var walletPUSD float64
			if chainReader != nil {
				if bal, berr := chainReader.PUSDBalance(ctx); berr == nil {
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
			var activeCost, activeValue, portfolioValue float64
			for _, p := range positions {
				if p.Size < 0.01 {
					continue
				}
				positionValue := p.CurrentValue
				if positionValue <= 0 && p.CurPrice > 0 {
					positionValue = p.Size * p.CurPrice
				}
				portfolioValue += positionValue
				bt, tracked := buyTimes.Get(p.Asset)
				if !tracked {
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
			totalAssets := walletPUSD + portfolioValue
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
				if err := json.Unmarshal(raw, &snap); err != nil {
					slog.Warn("pnl_snapshot_decode_err", "err", err)
				}
			}
			if snap.Date != todayStr {
				snap = dailySnap{Date: todayStr, TotalAssets: totalAssets}
				if jb, jerr := json.Marshal(snap); jerr != nil {
					slog.Warn("pnl_snapshot_encode_err", "err", jerr)
				} else if werr := writeAtomicFile(snapshotFile, jb, 0o600); werr != nil {
					slog.Warn("pnl_snapshot_write_err", "err", werr)
				}
			}
			dailyProfit := totalAssets - snap.TotalAssets

			// --- Format header ---
			fmt.Fprintf(&sb, "总资产: $%.2f (本金 $%.0f)\n", totalAssets, capital)
			fmt.Fprintf(&sb, "总盈利: $%+.2f (%+.1f%%)\n", totalProfit, totalProfitPct)
			fmt.Fprintf(&sb, "今日: $%+.2f\n", dailyProfit)
			fmt.Fprintf(&sb, "闲钱: $%.2f pUSD\n", walletPUSD)
			fmt.Fprintf(&sb, "持仓总市值: $%.2f · 跟踪中 %d 笔 ($%.2f 成本 / $%.2f 市值)\n", portfolioValue, len(activeLines), activeCost, activeValue)

			if len(activeLines) > 0 {
				fmt.Fprintf(&sb, "\n--- 持仓明细 (%d) ---\n", len(activeLines))
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
					fmt.Fprintf(&sb, "%s %s · %s%s\n   %.1f份 · $%.2f成本 · 入%.3f→现%.3f · $%+.2f (%+.1f%%)\n",
						l.emoji, title, direction, timeTag,
						l.size, l.cost, l.avgPrice, l.curPrice, l.pnl, l.pct)
					shown++
					if shown >= 20 {
						fmt.Fprintf(&sb, "... +%d more\n", len(activeLines)-20)
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

	// Redemption observer: report redeemable positions, but never sign or send
	// an on-chain transaction from the automated bot process.
	if chainReader != nil && walletAddress != "" {
		go func() {
			tk := time.NewTicker(1 * time.Hour)
			defer tk.Stop()
			alertedFile := filepath.Join("db", "live", "redeem-alerted.json")
			if err := os.MkdirAll(filepath.Dir(alertedFile), 0700); err != nil {
				slog.Warn("redeem_observer_state_dir_failed", "err", err)
				return
			}
			alerted := make(map[string]bool)
			if data, err := os.ReadFile(alertedFile); err == nil {
				_ = json.Unmarshal(data, &alerted)
			}
			saveAlerted := func() {
				data, err := json.Marshal(alerted)
				if err != nil {
					slog.Warn("redeem_observer_state_encode_failed", "err", err)
					return
				}
				if err := writeAtomicFile(alertedFile, data, 0o600); err != nil {
					slog.Warn("redeem_observer_state_write_failed", "err", err)
				}
			}
			checkRedeemable := func() {
				positions, err := fetchDataAPIPositions(ctx, walletAddress)
				if err != nil {
					slog.Warn("redeem_observer_fetch_err", "err", err)
					return
				}
				sgt := time.FixedZone("SGT", 8*3600)
				for _, p := range positions {
					if !p.Redeemable || p.CurPrice < 0.99 || p.Size < 0.01 {
						continue
					}
					balanceCtx, balanceCancel := context.WithTimeout(ctx, 15*time.Second)
					chainBalance, balanceErr := chainReader.ConditionalTokenBalance(balanceCtx, p.Asset)
					balanceCancel()
					if balanceErr != nil {
						slog.Warn("redeem_observer_balance_failed", "asset", p.Asset, "err", balanceErr)
						continue
					}
					if chainBalance.Sign() <= 0 {
						alerted[p.Asset] = true
						continue
					}
					if alerted[p.Asset] {
						continue
					}
					value := p.CurrentValue
					if value <= 0 {
						value = p.Size * p.CurPrice
					}
					slog.Info("redeem_maintenance_pending",
						"title", p.Title,
						"outcome", p.Outcome,
						"size", p.Size,
						"value", value,
						"conditionId", p.ConditionID,
						"negRisk", p.NegativeRisk,
						"outcomeIndex", p.OutcomeIndex,
					)
					notifier.SidecarAlert(fmt.Sprintf("💰 可赎回 · 等待独立赎回任务\n%s\n%s · %.1f份 · 约 $%.2f\n%s SGT",
						p.Title, p.Outcome, p.Size, value,
						time.Now().In(sgt).Format("01/02 15:04")))
					alerted[p.Asset] = true
					saveAlerted()
				}
			}
			time.Sleep(10 * time.Second)
			checkRedeemable()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tk.C:
					checkRedeemable()
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

func missingSettlementConditionIDs(ids []string, markets map[string]feed.Market) []string {
	missing := make([]string, 0, len(ids))
	for _, conditionID := range ids {
		if _, ok := markets[strings.ToLower(conditionID)]; !ok {
			missing = append(missing, conditionID)
		}
	}
	return missing
}

func settlementFallbackTokenIDs(open []strategy.Position, missingConditionIDs []string) []string {
	missing := make(map[string]struct{}, len(missingConditionIDs))
	for _, conditionID := range missingConditionIDs {
		missing[strings.ToLower(strings.TrimSpace(conditionID))] = struct{}{}
	}
	seen := make(map[string]struct{})
	var tokens []string
	for _, p := range open {
		if _, ok := missing[strings.ToLower(strings.TrimSpace(p.Market))]; !ok || p.AssetID == "" {
			continue
		}
		if _, duplicate := seen[p.AssetID]; duplicate {
			continue
		}
		seen[p.AssetID] = struct{}{}
		tokens = append(tokens, p.AssetID)
	}
	return tokens
}

func indexSettlementMarketsByToken(open []strategy.Position, missingConditionIDs []string, markets []feed.Market, byCondition map[string]feed.Market) int {
	missing := make(map[string]struct{}, len(missingConditionIDs))
	for _, conditionID := range missingConditionIDs {
		missing[strings.ToLower(strings.TrimSpace(conditionID))] = struct{}{}
	}
	conditionByToken := make(map[string]map[string]struct{})
	for _, p := range open {
		conditionID := strings.ToLower(strings.TrimSpace(p.Market))
		if _, ok := missing[conditionID]; !ok || p.AssetID == "" {
			continue
		}
		if conditionByToken[p.AssetID] == nil {
			conditionByToken[p.AssetID] = map[string]struct{}{}
		}
		conditionByToken[p.AssetID][conditionID] = struct{}{}
	}
	recovered := 0
	for _, market := range markets {
		for _, tokenID := range market.ClobTokenIDs() {
			for conditionID := range conditionByToken[tokenID] {
				if _, exists := byCondition[conditionID]; exists {
					continue
				}
				byCondition[conditionID] = market
				recovered++
			}
		}
	}
	return recovered
}

// countResolved returns the number of markets in the slice that have already
// settled on-chain (closed=true). Used for settlement-watcher status logging.
func countResolvedMap(markets map[string]feed.Market) int {
	n := 0
	for _, m := range markets {
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

func fetchDataAPIPositions(ctx context.Context, walletAddr string) ([]dataAPIPosition, error) {
	client := &nethttp.Client{Timeout: 20 * time.Second}
	return fetchDataAPIPositionsFrom(ctx, client, "https://data-api.polymarket.com", walletAddr)
}

func fetchDataAPIPositionsFrom(ctx context.Context, client *nethttp.Client, baseURL, walletAddr string) ([]dataAPIPosition, error) {
	const pageLimit = 500
	positions := make([]dataAPIPosition, 0, pageLimit)
	seen := make(map[string]struct{})
	for _, redeemable := range []bool{false, true} {
		for offset := 0; offset <= 10000; offset += pageLimit {
			query := neturl.Values{
				"user":          {strings.ToLower(strings.TrimSpace(walletAddr))},
				"sizeThreshold": {"0.01"},
				"limit":         {strconv.Itoa(pageLimit)},
				"offset":        {strconv.Itoa(offset)},
				"redeemable":    {strconv.FormatBool(redeemable)},
			}
			req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, strings.TrimRight(baseURL, "/")+"/positions?"+query.Encode(), nil)
			if err != nil {
				return nil, err
			}
			resp, err := client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("data-api positions: %w", err)
			}
			body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
			closeErr := resp.Body.Close()
			if readErr != nil || closeErr != nil {
				return nil, errors.Join(readErr, closeErr)
			}
			if resp.StatusCode >= 400 {
				return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
			}
			var page []dataAPIPosition
			if err := json.Unmarshal(body, &page); err != nil {
				return nil, fmt.Errorf("data-api decode: %w", err)
			}
			for _, position := range page {
				key := position.ConditionID + "\x00" + position.Asset
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				positions = append(positions, position)
			}
			if len(page) < pageLimit {
				break
			}
			if offset == 10000 {
				return nil, errors.New("data-api positions exceeded pagination limit")
			}
		}
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
	EventStart     time.Time // scheduled event start; zero when unavailable
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
				EventStart:  m.EventStartTime(),
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
	pm                   *strategy.PositionManager
	exit                 *strategy.ExitTracker
	ladder               *strategy.LadderTracker
	paper                order.Client
	rm                   *risk.Manager
	pending              *notify.PendingStore
	closePending         *notify.CloseStore
	notifier             notify.Notifier
	meta                 map[string]*assetMeta
	src                  *sourceTracker
	recorder             *tickrec.Recorder
	jrn                  *journal.Journal
	largeFillUSD         float64
	exitMode             string
	holdMax              time.Duration
	eventPostHold        time.Duration
	riskStatePath        string
	savePositions        func()
	savePositionsDurable func() error
	markExecutionApplied func(order.Result)
	prepareIntent        func(order.Intent, float64) order.Intent
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
		Type:    order.FAK,
	}
	reserveTick := feed.Tick{
		AssetID: choice.AssetID, Market: p.Market,
		Time: now, Mid: choice.Mid,
	}
	pos, err := h.pm.OpenSized(choice.AssetID, p.Market, reserveTick, sizeUSD)
	if err != nil {
		return "", fmt.Errorf("开仓失败: %s", err.Error())
	}
	intent.ClientID = pos.ID
	intent.Reason = "manual_entry"
	if h.prepareIntent != nil {
		intent = h.prepareIntent(intent, choice.Mid)
	}
	res, err := h.paper.Submit(ctx, intent)
	if err != nil {
		if res.Status == order.StatusPending {
			_ = h.pm.SetOpenAttribution(pos.ID, "manual_pending", res.OrderID)
			if h.savePositions != nil {
				h.savePositions()
			}
		} else {
			_ = h.pm.CancelOpen(pos.ID)
		}
		return "", fmt.Errorf("下单失败: %s", err.Error())
	}
	if res.Status != order.StatusFilled {
		_ = h.pm.CancelOpen(pos.ID)
		return "", fmt.Errorf("下单未成交: %s", res.Status)
	}
	filledAt := res.FilledAt
	if filledAt.IsZero() {
		filledAt = time.Now()
	}
	if err := h.pm.ApplyOpenFill(pos.ID, res.AvgPrice, res.FilledSize, filledAt); err != nil {
		return "", fmt.Errorf("成交记账失败: %s", err.Error())
	}
	entryTick := feed.Tick{
		AssetID: choice.AssetID, Market: p.Market,
		Time: filledAt, Mid: res.AvgPrice,
	}
	_ = h.pm.SetOpenFee(pos.ID, res.FeeUSD)
	markPositionSource(h.pm, h.src, pos.ID, "manual", res.OrderID)
	eventStart := time.Time{}
	if me := h.meta[choice.AssetID]; me != nil {
		eventStart = me.EventStart
	}
	planned, err := h.pm.ConfigureOpenHold(pos.ID, eventStart, h.holdMax, h.eventPostHold)
	if err != nil {
		return "", fmt.Errorf("持仓计划失败: %s", err.Error())
	}
	if h.savePositionsDurable != nil {
		if err := h.savePositionsDurable(); err != nil {
			return "", fmt.Errorf("持仓落盘失败: %s", err.Error())
		}
		if h.markExecutionApplied != nil {
			h.markExecutionApplied(res)
		}
	} else if h.savePositions != nil {
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
				h.ladder.OpenWithDeadline(pos.ID, p.Market, choice.AssetID, entryTick, pos.Units, planned.ExitDeadline)
			}
		}
	}

	if h.recorder != nil {
		if rerr := h.recorder.Start(pos.ID, choice.AssetID); rerr != nil {
			slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
		}
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
		"hold_profile", planned.HoldProfile,
		"event_start", planned.EventStart,
		"exit_deadline", planned.ExitDeadline,
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
		if _, err := h.pm.BeginClose(pos.ID, pos.Units); err != nil {
			slog.Warn("whale_close_reserve_fail", "pos", pos.ID, "err", err)
			continue
		}
		sellIntent := order.Intent{
			ClientID:   pos.ID,
			Reason:     "manual_whale_sell",
			AssetID:    pos.AssetID,
			Market:     pos.Market,
			Side:       order.Sell,
			SizeUSD:    pos.SizeUSD,
			SizeShares: pos.Units,
			LimitPx:    ci.WhalePrice,
			Type:       order.FAK,
		}
		if h.prepareIntent != nil {
			sellIntent = h.prepareIntent(sellIntent, ci.WhalePrice)
		}
		res, serr := h.paper.Submit(ctx, sellIntent)
		if serr != nil || res.Status != order.StatusFilled {
			if res.Status == order.StatusPending {
				if h.savePositions != nil {
					h.savePositions()
				}
			} else {
				h.pm.AbortClose(pos.ID)
			}
			slog.Warn("whale_close_sell_reject", "pos", pos.ID, "err", serr, "status", res.Status)
			continue
		}
		if err := h.pm.ApplyCloseFill(pos.ID, res.FilledSize); err != nil {
			slog.Error("whale_close_fill_size_invalid", "pos", pos.ID, "execution_id", res.ExecutionID, "filled_size", res.FilledSize, "err", err)
			continue
		}
		closedAt := filledAtOr(res.FilledAt, now)
		sig := strategy.ExitSignal{
			CloseExecutionID: res.ExecutionID,
			CloseOrderID:     res.OrderID,
			AssetID:          pos.AssetID,
			Market:           pos.Market,
			Time:             closedAt,
			EntryMid:         pos.EntryMid,
			PeakMid:          pos.EntryMid,
			ExitMid:          res.AvgPrice,
			HeldFor:          closedAt.Sub(pos.EntryTime),
			ChangePP:         (res.AvgPrice - pos.EntryMid) * 100,
			ExitFeeUSD:       res.FeeUSD,
			Reason:           strategy.ExitReason("whale_sell"),
		}
		closed, cerr := h.pm.CommitClose(pos.ID, sig)
		if cerr != nil {
			slog.Warn("whale_close_miss", "pos", pos.ID, "err", cerr.Error())
			continue
		}
		remaining, stillOpen := h.pm.OpenByID(pos.ID)
		positionDurable := true
		if h.savePositionsDurable != nil {
			if err := h.savePositionsDurable(); err != nil {
				positionDurable = false
				slog.Error("positions_save_after_close_fail", "execution_id", res.ExecutionID, "pos", pos.ID, "err", err)
			}
		} else if h.savePositions != nil {
			h.savePositions()
		}
		hadLadder := h.ladder.Has(pos.ID)
		h.ladder.Forget(pos.ID)
		if stillOpen && hadLadder {
			h.ladder.SyncPosition(remaining)
		}
		if !stillOpen && h.recorder != nil {
			_ = h.recorder.Stop(closed.ID)
		}
		entryFeeShare := closed.EntryFeeUSD
		exitFee := closed.ExitFeeUSD
		netPnL := closed.NetPnLUSD
		totalNetPnL += netPnL
		closedCount++
		stats := h.pm.Stats()
		closeSource, _ := h.src.Peek(closed.ID)
		if riskEligibleSignalSource(closeSource) {
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
				SizeUSD:  closed.SizeUSD,
				PnLUSD:   netPnL,
				EntryPx:  pos.EntryMid,
				ExitPx:   res.AvgPrice,
				Reason:   "whale_sell",
				HeldSec:  int(sig.HeldFor.Seconds()),
			})
		}
		var source, openOID string
		if stillOpen {
			source, openOID = h.src.Peek(closed.ID)
		} else {
			source, openOID = h.src.Take(closed.ID)
		}
		journalDurable := h.jrn == nil
		if h.jrn != nil {
			if err := h.jrn.Append(journal.TradeRecord{
				ExecutionID: res.ExecutionID,
				ID:          closed.ID, AssetID: closed.AssetID, Market: closed.Market,
				Question:      ci.Question,
				Outcome:       ci.Outcome,
				Side:          "buy",
				SizeUSD:       closed.SizeUSD,
				Units:         closed.Units,
				EntryMid:      closed.EntryMid,
				EntryTime:     closed.EntryTime,
				ExitMid:       closed.ExitMid,
				ExitTime:      closed.ExitTime,
				ExitReason:    "whale_sell",
				HeldSec:       int(sig.HeldFor.Seconds()),
				PnLUSD:        closed.PnLUSD,
				EntryFeeUSD:   entryFeeShare,
				ExitFeeUSD:    exitFee,
				NetPnLUSD:     netPnL,
				OpenOrderID:   openOID,
				CloseOrderID:  res.OrderID,
				Mode:          orderMode(h.paper),
				SignalSource:  source,
				PolicyVersion: closed.PolicyVersion,
			}); err != nil {
				slog.Warn("journal_append_fail", "asset", short(pos.AssetID), "err", err)
			} else {
				journalDurable = true
			}
		}
		if positionDurable && journalDurable && h.markExecutionApplied != nil {
			h.markExecutionApplied(res)
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

// sourceTracker is the runtime lookup for journal attribution. The same values
// are persisted on Position and rehydrated at startup.
type sourceTracker struct {
	mu sync.Mutex
	m  map[string]sourceEntry
}

type sourceEntry struct {
	source      string // "auto" or "manual"
	openOrderID string
}

type buyTimeStore struct {
	mu    sync.RWMutex
	path  string
	times map[string]time.Time
}

func loadBuyTimeStore(path string) (*buyTimeStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	store := &buyTimeStore{path: path, times: make(map[string]time.Time)}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &store.times); err != nil {
			return nil, fmt.Errorf("decode %s: %w", path, err)
		}
	}
	return store, nil
}

func (s *buyTimeStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.times)
}

func (s *buyTimeStore) Get(assetID string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.times[assetID]
	return t, ok
}

func (s *buyTimeStore) Set(assetID string, boughtAt time.Time) error {
	if strings.TrimSpace(assetID) == "" || boughtAt.IsZero() {
		return errors.New("buy time requires asset id and timestamp")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.times[assetID] = boughtAt
	raw, err := json.MarshalIndent(s.times, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomicFile(s.path, raw, 0o600)
}

func writeAtomicFile(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".state-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Chmod(mode); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	dirHandle, err := os.Open(dir)
	if err != nil {
		return err
	}
	syncErr := dirHandle.Sync()
	closeErr := dirHandle.Close()
	return errors.Join(syncErr, closeErr)
}

func newSourceTracker() *sourceTracker {
	return &sourceTracker{m: map[string]sourceEntry{}}
}

func markPositionSource(pm *strategy.PositionManager, src *sourceTracker, posID, source, openOrderID string) {
	if err := pm.SetOpenAttribution(posID, source, openOrderID); err != nil {
		slog.Warn("position_attribution_fail", "pos", posID, "source", source, "err", err)
	}
	if src != nil {
		src.Mark(posID, source, openOrderID)
	}
}

func persistedPositionSource(p strategy.Position) string {
	if p.SignalSource != "" {
		return p.SignalSource
	}
	if p.Source == "copytrade" || p.Source == "copytrade_football_score" {
		source := p.Source
		if p.WalletLabel != "" {
			source += "_" + p.WalletLabel
		}
		return source
	}
	if p.Source != "" {
		return p.Source
	}
	return "auto"
}

func recoveredEntrySource(reason string, p strategy.Position) string {
	if source := strings.TrimSuffix(strings.TrimSpace(p.SignalSource), "_pending"); source != "" {
		return source
	}
	if p.Source != "" {
		return persistedPositionSource(p)
	}
	switch strings.TrimSpace(reason) {
	case "manual_entry":
		return "manual"
	case "lottery_entry":
		return "lottery"
	case "copytrade_entry":
		return "copytrade"
	default:
		return "auto"
	}
}

func orderMode(client order.Client) string {
	if client != nil && strings.Contains(client.Name(), "v2-live") {
		return "live"
	}
	return "paper"
}

func isTimeoutExitReason(reason string) bool {
	switch strategy.ExitReason(reason) {
	case strategy.ExitTimeout, strategy.ExitLadderTimeout:
		return true
	default:
		return false
	}
}

func filledAtOr(filledAt, fallback time.Time) time.Time {
	if !filledAt.IsZero() {
		return filledAt
	}
	return fallback
}

func validateTradingStatePaths(live bool, paths ...string) error {
	for _, path := range paths {
		cleaned := filepath.Clean(strings.TrimSpace(path))
		if cleaned == "." || cleaned == "" {
			continue
		}
		parts := strings.FieldsFunc(strings.ToLower(cleaned), func(r rune) bool {
			return r == '/' || r == '\\'
		})
		liveScoped := false
		for _, part := range parts {
			if part == "live" || strings.HasSuffix(part, "-live") || strings.HasPrefix(part, "live-") {
				liveScoped = true
			}
			if live && strings.Contains(part, "paper") {
				return fmt.Errorf("live state path must not use a paper directory: %s", path)
			}
			if !live && (part == "live" || strings.HasSuffix(part, "-live") || strings.HasPrefix(part, "live-")) {
				return fmt.Errorf("paper state path must not use a live directory: %s", path)
			}
		}
		if live && !liveScoped {
			return fmt.Errorf("live state path must include a dedicated live directory: %s", path)
		}
		if live {
			if err := rejectSymlinkPathComponents(cleaned); err != nil {
				return fmt.Errorf("live state path is unsafe: %s: %w", path, err)
			}
		}
	}
	return nil
}

func rejectSymlinkPathComponents(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for current := absPath; ; current = filepath.Dir(current) {
		info, statErr := os.Lstat(current)
		switch {
		case statErr == nil && info.Mode()&os.ModeSymlink != 0:
			return fmt.Errorf("%s is a symlink", current)
		case statErr != nil && !os.IsNotExist(statErr):
			return fmt.Errorf("inspect %s: %w", current, statErr)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return nil
		}
	}
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
	trades, err := journal.ReadClosedDay(dir, day)
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

func copytradeCollectionMarketDecision(q, slug, outcome string, allowFootballScore bool) (bool, string) {
	if feed.IsFootballScoreMarketText(q + " " + slug) {
		if !allowFootballScore {
			return false, "derivative_filtered"
		}
		if strings.EqualFold(strings.TrimSpace(outcome), "No") {
			return false, "football_score_no_filtered"
		}
	}
	return true, ""
}

func copytradeEntryPriceFloor(footballScore, allowFootballScore bool) float64 {
	if footballScore && allowFootballScore {
		return footballScoreMinEntryPrice
	}
	return followMinEntryPrice
}

func paperTimeoutExitPrice(tick feed.Tick, now time.Time) (float64, bool) {
	if tick.BestBid <= 0 || tick.BestBid >= 1 || tick.BestBidSize <= 0 || tick.QuoteTime.IsZero() {
		return 0, false
	}
	age := now.Sub(tick.QuoteTime)
	if age < 0 || age > 5*time.Second {
		return 0, false
	}
	return tick.BestBid, true
}

func copytradeMarketSize(base float64, footballScore, allowFootballScore bool, scoreCap float64) float64 {
	if !footballScore || !allowFootballScore || scoreCap <= 0 || base <= scoreCap {
		return base
	}
	return scoreCap
}

func copytradeWalletSize(configured float64, live bool, tier, followAction string, smartMoneyScore float64) float64 {
	if !live && configured > 0 {
		return configured
	}
	if followAction == "auto-small" && smartMoneyScore >= 80 {
		return 20
	}
	switch tier {
	case "A":
		return 10
	case "B":
		return 5
	default:
		return configured
	}
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

func copytradeExposureEventKey(ev whale.AlertEvent) string {
	if !feed.IsFootballScoreMarketText(ev.Question + " " + ev.Slug) {
		return "event:" + whaleEventKey(ev.Question)
	}
	if slug := strings.ToLower(strings.TrimSpace(ev.Slug)); slug != "" {
		return "football_score_event:" + slug
	}
	return "football_score_event:" + whaleEventKey(ev.Question)
}

func copytradeCoreDedupeKey(ev whale.AlertEvent) string {
	wallet := strings.ToLower(strings.TrimSpace(ev.Wallet))
	conditionID := strings.ToLower(strings.TrimSpace(ev.ConditionID))
	assetID := strings.ToLower(strings.TrimSpace(ev.AssetID))
	if wallet == "" || conditionID == "" || assetID == "" {
		return ""
	}
	return wallet + "|" + conditionID + "|" + assetID
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
