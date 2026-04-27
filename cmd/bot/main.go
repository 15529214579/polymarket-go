package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	nethttp "net/http"
	neturl "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/15529214579/polymarket-go/internal/arb"
	"github.com/15529214579/polymarket-go/internal/btc"
	"github.com/15529214579/polymarket-go/internal/config"
	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/injury"
	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/notify"
	"github.com/15529214579/polymarket-go/internal/odds"
	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/15529214579/polymarket-go/internal/risk"
	"github.com/15529214579/polymarket-go/internal/strategy"
	"github.com/15529214579/polymarket-go/internal/tickrec"
	"github.com/15529214579/polymarket-go/internal/whale"
)

func main() {
	mode := flag.String("mode", "run", "run | discover | feed | sample | detect | prompt-test | daily-report | arb-scan")
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
	arbEnabled := flag.Bool("arb_enabled", true, "enable arb scanner (odds vs Polymarket gap detection)")
	arbInterval := flag.Duration("arb_interval", 12*time.Hour, "arb scan interval (budget: 500 req/month free tier)")
	arbMinGapPP := flag.Float64("arb_min_gap_pp", 5.0, "arb minimum net EV in pp to flag")
	arbDBPath := flag.String("arb_db", "db/odds.db", "SQLite path for odds snapshots")
	// OddsPapi high-frequency sharp-line scanner (Pinnacle + 350 books).
	oddsPapiEnabled := flag.Bool("oddspapi_enabled", false, "enable OddsPapi high-freq football scanner (Pinnacle sharp line)")
	oddsPapiInterval := flag.Duration("oddspapi_interval", 3*time.Hour, "OddsPapi scan interval (budget: 250 req/month free tier)")
	oddsPapiBookmaker := flag.String("oddspapi_bookmaker", "pinnacle", "OddsPapi bookmaker to fetch (pinnacle, bet365, etc)")
	oddsPapiSports := flag.String("oddspapi_sports", "soccer_epl,soccer_spain_la_liga,soccer_uefa_champs_league", "comma-separated sport keys for OddsPapi")
	injuryEnabled := flag.Bool("injury_enabled", false, "enable NBA injury report scanner (ESPN API)")
	injuryInterval := flag.Duration("injury_interval", 30*time.Minute, "injury scan interval")
	injuryStarOnly := flag.Bool("injury_star_only", true, "only alert on star players (top ~3-4 per team)")
	whaleEnabled := flag.Bool("whale_enabled", false, "enable smart-money whale trade tracker")
	confirmDelay := flag.Duration("confirm_delay", 10*time.Second, "wait N seconds after signal trigger, re-check price before entry")
	whaleWallets := flag.String("whale_wallets", "", "tracked wallets: addr|label|minUSD|profileURL,... (comma-separated)")
	whaleWallet := flag.String("whale_wallet", "", "(legacy) single target wallet address (hex 0x…)")
	whaleProfile := flag.String("whale_profile", "", "(legacy) whale's Polymarket profile URL")
	whaleMinUSD := flag.Float64("whale_min_usd", 1000, "(legacy) minimum notional USD to trigger alert")
	whaleInterval := flag.Duration("whale_interval", 30*time.Second, "whale trade poll interval")
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
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := config.LoadDotEnv(*envFile); err != nil {
		slog.Warn("dotenv_load_warn", "path", *envFile, "err", err.Error())
	}

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
		}
		oddsPapiCfg := odds.OddsPapiConfig{
			Enabled:   *oddsPapiEnabled,
			Interval:  *oddsPapiInterval,
			Bookmaker: *oddsPapiBookmaker,
		}
		if *oddsPapiSports != "" {
			oddsPapiCfg.SportKeys = strings.Split(*oddsPapiSports, ",")
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
			if err := runDetect(ctx, *maxMarkets, *windowSec, *slippageBp, *feeBp, *largeFillUSD, *signalMode, *exitMode, *journalDir, *tickPathDir, *minEntry, *maxEntry, ladderCfg, *lotteryEnabled, lottCfg, *arbEnabled, *arbInterval, *arbMinGapPP, *arbDBPath, injCfg, whaleCfg, oddsPapiCfg, *confirmDelay, btcCfg, updownCfg); err != nil && ctx.Err() == nil {
			slog.Error("detect failed", "err", err)
			os.Exit(1)
		}
	case "daily-report":
		if err := runDailyReport(ctx, *journalDir, *reportDay, *reportPush); err != nil {
			slog.Error("daily-report failed", "err", err)
			os.Exit(1)
		}
	case "arb-scan":
		if err := runArbScan(ctx, *arbDBPath, *arbMinGapPP); err != nil {
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
	mkts := feed.FilterSports(all)
	slog.Info("gamma.discover",
		"total_active", len(all),
		"sports", len(mkts),
		"lol", len(feed.FilterLoL(all)),
		"nba_epl_playoffs", len(mkts)-len(feed.FilterLoL(all)),
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
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
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
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
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

func runDetect(ctx context.Context, topN, windowSec int, slippageBp, feeBp, largeFillUSD float64, signalMode, exitMode, journalDir, tickPathDir string, minEntry, maxEntry float64, ladderCfg strategy.LadderConfig, lotteryEnabled bool, lotteryCfg strategy.LotteryConfig, arbEnabled bool, arbInterval time.Duration, arbMinGapPP float64, arbDBPath string, injCfg injury.Config, whaleCfg whale.Config, oddsPapiCfg odds.OddsPapiConfig, confirmDelay time.Duration, btcCfg btc.StrategyConfig, updownCfg btc.UpDownConfig) error {
	if signalMode != "auto" && signalMode != "prompt" && signalMode != "whale" {
		return fmt.Errorf("invalid signal_mode %q (want auto|prompt|whale)", signalMode)
	}
	if exitMode != "hold" && exitMode != "auto" && exitMode != "ladder" {
		return fmt.Errorf("invalid exit_mode %q (want hold|auto|ladder)", exitMode)
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
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
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
		"lol", countBy(mkts, feed.IsLoLMarket),
		"dota2", countBy(mkts, feed.IsDota2Market),
		"basketball", countBy(mkts, feed.IsBasketballMarket),
		"football", countBy(mkts, feed.IsFootballMarket),
		"assets", len(assetIDs),
		"window_sec", windowSec,
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
	pm := strategy.NewPositionManager(posCfg)
	paper := order.NewPaperClientWithFee(slippageBp, feeBp)
	riskCfg := risk.DefaultConfig()
	riskCfg.FeedConnected = ws.Connected
	rm := risk.New(riskCfg, time.Now())
	const riskStatePath = "db/risk_state.json"
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
	if lotteryEnabled {
		slog.Info("lottery.ready",
			"min_price", lotteryCfg.MinPrice,
			"max_price", lotteryCfg.MaxPrice,
			"lol_min_price", lotteryCfg.LoLMinPrice,
			"size_usd", lotteryCfg.SizeUSD,
			"scan_interval", lotteryCfg.ScanInterval.String(),
		)
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
	if sidecarToken != "" && sidecarChat != "" {
		chatID, err := strconv.ParseInt(sidecarChat, 10, 64)
		if err != nil {
			slog.Warn("sidecar_chat_id_parse_fail", "err", err.Error())
		} else {
			h := &buyHandler{
				pm:            pm,
				exit:          exit,
				ladder:        ladder,
				paper:         paper,
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
					_ = rm.SaveState(riskStatePath)
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
	go func() {
		if err := det.Run(ctx); err != nil && ctx.Err() == nil {
			slog.Error("detector exited", "err", err)
		}
	}()

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
						res, err := paper.Submit(ctx, sellIntent)
						if err != nil {
							slog.Warn("paper_sell_reject",
								"asset", short(sig.AssetID),
								"limit", sig.ExitMid,
								"err", err.Error())
							continue
						}
						// Override exit mid with the actual fill price so realized PnL
						// reflects paper slippage.
						sig.ExitMid = res.AvgPrice
						sig.ChangePP = (res.AvgPrice - sig.EntryMid) * 100
						closed, err := pm.CloseFirstByAsset(sig.AssetID, sig)
						if err != nil {
							slog.Warn("paper_close_miss", "asset", short(sig.AssetID), "err", err.Error())
							continue
						}
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
							_ = rm.SaveState(riskStatePath)
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
						res, err := paper.Submit(ctx, sellIntent)
						if err != nil {
							slog.Warn("paper_ladder_sell_reject",
								"pos", p.ID,
								"asset", short(ex.AssetID),
								"tranche", ex.Tranche,
								"limit", ex.ExitMid,
								"err", err.Error())
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
							_ = rm.SaveState(riskStatePath)
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
	injScanner := injury.NewScanner(injCfg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case sig, ok := <-det.Signals():
				if !ok {
					return
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
				// Whale + momentum run as union — both signal sources active.
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
				// Paper stacking: no per-asset/per-market dedupe here. Auto mode
				// still has a 5-min cooldown in detector.go so one asset can't
				// spam opens, and pm enforces MaxOpenPositions + exposure caps.

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
							if err != nil || msgID == 0 {
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
					AssetID: sig.AssetID,
					Market:  sig.Market,
					Side:    order.Buy,
					SizeUSD: posCfg.PerPositionUSD,
					LimitPx: sig.Mid,
					Type:    order.GTC,
				}
				res, err := paper.Submit(ctx, buyIntent)
				if err != nil {
					slog.Warn("paper_buy_reject",
						"asset", short(sig.AssetID),
						"limit", sig.Mid,
						"err", err.Error())
					continue
				}
				// Book the position at the actual fill price so slippage is priced in.
				entryTick := feed.Tick{
					AssetID: sig.AssetID, Market: sig.Market,
					Time: sig.Time, Mid: res.AvgPrice,
				}
				pos, err := pm.Open(sig.AssetID, sig.Market, entryTick)
				if err != nil {
					slog.Info("paper_open_skip",
						"asset", short(sig.AssetID),
						"q", metaQ(meta, sig.AssetID),
						"order_id", res.OrderID,
						"reason", err.Error(),
					)
					continue
				}
				_ = pm.SetOpenFee(pos.ID, res.FeeUSD)
				switch exitMode {
				case "auto":
					exit.Open(sig.AssetID, sig.Market, entryTick)
				case "ladder":
					ladder.Open(pos.ID, sig.Market, sig.AssetID, entryTick, pos.Units)
				}
				if recorder != nil {
					if rerr := recorder.Start(pos.ID, sig.AssetID); rerr != nil {
						slog.Warn("tickrec_start_fail", "pos", pos.ID, "err", rerr.Error())
					}
				}
				src.Mark(pos.ID, "auto", res.OrderID)
				stats := pm.Stats()
				slog.Info("paper_open",
					"id", pos.ID,
					"order_id", res.OrderID,
					"asset", short(sig.AssetID),
					"q", metaQ(meta, sig.AssetID),
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
						if err != nil || msgID == 0 {
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
	if lotteryEnabled {
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
						res, err := paper.Submit(ctx, buyIntent)
						if err != nil {
							slog.Warn("lottery_buy_reject",
								"asset", short(c.AssetID),
								"mid", c.Mid,
								"sport", string(c.Sport),
								"err", err.Error())
							continue
						}
						entryTick := feed.Tick{
							AssetID: c.AssetID, Market: c.Market,
							Time: c.Time, Mid: res.AvgPrice,
						}
						pos, err := pm.Open(c.AssetID, c.Market, entryTick)
						if err != nil {
							slog.Info("lottery_open_skip",
								"asset", short(c.AssetID),
								"q", metaQ(meta, c.AssetID),
								"reason", err.Error())
							continue
						}
						_ = pm.SetOpenFee(pos.ID, res.FeeUSD)
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

	// Arb scanner: periodic cross-venue price comparison (Polymarket vs bookmaker odds).
	// Runs on arbInterval cadence, stores snapshots to SQLite, logs opportunities.
	if arbEnabled {
		arbStore, arbErr := arb.NewStore(arbDBPath)
		if arbErr != nil {
			slog.Warn("arb_store_init_fail", "err", arbErr.Error(), "path", arbDBPath)
		} else {
			oddsClient := odds.NewClient("", "")
			scanCfg := arb.DefaultScanConfig()
			scanCfg.MinGapPP = arbMinGapPP
			scanner := arb.NewScanner(oddsClient, arbStore, scanCfg)
			go func() {
				defer arbStore.Close()
				// Run one scan immediately on startup.
				slog.Info("arb_scanner.ready", "interval", arbInterval.String(), "min_gap_pp", arbMinGapPP, "db", arbDBPath)
				opps, err := scanner.Scan(ctx)
				if err != nil {
					slog.Warn("arb_scan_fail", "err", err.Error())
				} else {
					usage := oddsClient.Usage()
					slog.Info("arb_scan_done", "opportunities", len(opps),
						"api_remaining", usage.RequestsRemaining, "api_used", usage.RequestsUsed,
						"db_rows", arbStore.Count())
					for _, o := range opps {
						slog.Info("arb_opportunity",
							"sport", o.Sport,
							"event", o.EventName,
							"dir", o.Direction,
							"gap_pp", o.GapPP,
							"net_ev_pp", o.NetEvPP,
							"poly", o.PolymarketPrice,
							"bk_prob", o.BookmakerProb,
							"market", o.MarketTitle,
						)
					}
				}
				tk := time.NewTicker(arbInterval)
				defer tk.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-tk.C:
						opps, err := scanner.Scan(ctx)
						if err != nil {
							slog.Warn("arb_scan_fail", "err", err.Error())
							continue
						}
						usage := oddsClient.Usage()
						slog.Info("arb_scan_done", "opportunities", len(opps),
							"api_remaining", usage.RequestsRemaining, "api_used", usage.RequestsUsed,
							"db_rows", arbStore.Count())
						for _, o := range opps {
							slog.Info("arb_opportunity",
								"sport", o.Sport,
								"event", o.EventName,
								"dir", o.Direction,
								"gap_pp", o.GapPP,
								"net_ev_pp", o.NetEvPP,
								"poly", o.PolymarketPrice,
								"bk_prob", o.BookmakerProb,
								"market", o.MarketTitle,
							)
						}
					}
				}
			}()
		}
	}

	// OddsPapi high-frequency sharp-line scanner (Pinnacle / bet365).
	// Runs independently from the Odds API scanner at higher frequency,
	// targeting football (EPL, UCL, La Liga) with sharp bookmaker lines.
	if oddsPapiCfg.Enabled {
		papiClient := odds.NewOddsPapiClient("", oddsPapiCfg.Bookmaker, "")
		if !papiClient.HasKey() {
			slog.Warn("oddspapi_disabled_no_key", "env", "ODDSPAPI_API_KEY")
		} else {
			arbStoreP, arbErrP := arb.NewStore(arbDBPath)
			if arbErrP != nil {
				slog.Warn("oddspapi_store_init_fail", "err", arbErrP.Error())
			} else {
				scanCfgP := arb.DefaultScanConfig()
				scanCfgP.MinGapPP = arbMinGapPP
				scanCfgP.MinBookCount = 1 // single bookmaker, no consensus needed
				scannerP := arb.NewScanner(nil, arbStoreP, scanCfgP)
				go func() {
					defer arbStoreP.Close()
					slog.Info("oddspapi_scanner.ready",
						"interval", oddsPapiCfg.Interval.String(),
						"bookmaker", oddsPapiCfg.Bookmaker,
						"sports", oddsPapiCfg.SportKeys,
					)
					// Immediate scan on startup.
					papiOdds, err := papiClient.FetchFootballOdds(ctx, oddsPapiCfg.SportKeys)
					if err != nil {
						slog.Warn("oddspapi_scan_fail", "err", err.Error())
					} else if len(papiOdds) > 0 {
						opps, err := scannerP.ScanWithOdds(ctx, papiOdds, "oddspapi/"+oddsPapiCfg.Bookmaker)
						if err != nil {
							slog.Warn("oddspapi_match_fail", "err", err.Error())
						} else {
							usage := papiClient.Usage()
							slog.Info("oddspapi_scan_done",
								"opportunities", len(opps),
								"odds_items", len(papiOdds),
								"api_remaining", usage.RequestsRemaining,
								"api_used", usage.RequestsUsed,
							)
							for _, o := range opps {
								slog.Info("oddspapi_opportunity",
									"sport", o.Sport,
									"event", o.EventName,
									"dir", o.Direction,
									"gap_pp", o.GapPP,
									"net_ev_pp", o.NetEvPP,
									"poly", o.PolymarketPrice,
									"bk_prob", o.BookmakerProb,
									"bk", o.Bookmaker,
									"market", o.MarketTitle,
								)
							}
						}
					}
					tk := time.NewTicker(oddsPapiCfg.Interval)
					defer tk.Stop()
					for {
						select {
						case <-ctx.Done():
							return
						case <-tk.C:
							papiOdds, err := papiClient.FetchFootballOdds(ctx, oddsPapiCfg.SportKeys)
							if err != nil {
								slog.Warn("oddspapi_scan_fail", "err", err.Error())
								continue
							}
							if len(papiOdds) == 0 {
								continue
							}
							opps, err := scannerP.ScanWithOdds(ctx, papiOdds, "oddspapi/"+oddsPapiCfg.Bookmaker)
							if err != nil {
								slog.Warn("oddspapi_match_fail", "err", err.Error())
								continue
							}
							usage := papiClient.Usage()
							slog.Info("oddspapi_scan_done",
								"opportunities", len(opps),
								"odds_items", len(papiOdds),
								"api_remaining", usage.RequestsRemaining,
								"api_used", usage.RequestsUsed,
							)
							for _, o := range opps {
								slog.Info("oddspapi_opportunity",
									"sport", o.Sport,
									"event", o.EventName,
									"dir", o.Direction,
									"gap_pp", o.GapPP,
									"net_ev_pp", o.NetEvPP,
									"poly", o.PolymarketPrice,
									"bk_prob", o.BookmakerProb,
									"bk", o.Bookmaker,
									"market", o.MarketTitle,
								)
							}
						}
					}
				}()
			}
		}
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
			scanOnce := func() {
				alerts, err := injScanner.Scan(ctx)
				if err != nil {
					slog.Warn("injury_scan_fail", "err", err.Error())
					return
				}
				for _, a := range alerts {
					if !injuryTeamInMarkets(a.Team, meta, assetSport) {
						continue
					}
					slog.Info("injury_alert",
						"team", a.Team,
						"player", a.StarPlayer,
						"status", string(a.Status),
						"impact", a.Impact,
					)
					notifier.InjuryAlert(notify.InjuryAlertEvent{
						Team:       a.Team,
						StarPlayer: a.StarPlayer,
						Status:     string(a.Status),
						Reason:     a.Reason,
						Impact:     a.Impact,
					})
					if a.Status == injury.StatusOut {
						injuryPushOpponentPrompt(a, meta, assetSport, sampler, pending, notifier)
					}
				}
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
					for _, a := range alerts {
						if !injuryTeamInMarkets(a.Team, meta, assetSport) {
							continue
						}
						slog.Info("injury_alert",
							"team", a.Team,
							"player", a.StarPlayer,
							"status", string(a.Status),
							"impact", a.Impact,
						)
						notifier.InjuryAlert(notify.InjuryAlertEvent{
							Team:       a.Team,
							StarPlayer: a.StarPlayer,
							Status:     string(a.Status),
							Reason:     a.Reason,
							Impact:     a.Impact,
						})
						if a.Status == injury.StatusOut {
							injuryPushOpponentPrompt(a, meta, assetSport, sampler, pending, notifier)
						}
					}
				}
			}
		}()
	}

	// Smart-money whale tracker: polls target wallet's CLOB trades and
	// pushes DM for large orders. Feature-flagged via -whale_enabled.
	// In whale signal_mode: BUY → SignalPrompt with buttons (boss clicks to follow);
	// SELL → auto-close matching positions.
	if whaleCfg.Enabled {
		wt := whale.NewTracker(whaleCfg, func(ev whale.AlertEvent) {
			side := strings.ToUpper(ev.Side)

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
				ctxLine := fmt.Sprintf("🐋 %s 跟单 · $%.0f · %.0f shares", whaleTag, ev.Notional, ev.SizeUnits)
				if ev.TotalShares > 0 {
					ctxLine += fmt.Sprintf("\n持仓: %.0f shares (均价 $%.4f)", ev.TotalShares, ev.AvgPrice)
				}
				if ev.LinkURL != "" {
					ctxLine += "\n" + ev.LinkURL
				}

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
						if err != nil || msgID == 0 {
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
						if err != nil || msgID == 0 {
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
						if err != nil || msgID == 0 {
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
				// Collect unique conditionIDs from meta.
				seen := make(map[string]struct{}, len(open))
				ids := make([]string, 0, len(open))
				for _, p := range open {
					me := meta[p.AssetID]
					if me == nil || me.ConditionID == "" {
						continue
					}
					if _, ok := seen[me.ConditionID]; ok {
						continue
					}
					seen[me.ConditionID] = struct{}{}
					ids = append(ids, me.ConditionID)
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
					byCond[m.ConditionID] = m
				}
				// Periodic "still holding" log (once per 5 min) — easy to grep for.
				now := time.Now()
				if now.Sub(lastHeldLog) >= 5*time.Minute {
					lastHeldLog = now
					slog.Info("hold_status",
						"open", len(open),
						"markets_polled", len(ids),
						"resolved_seen", countResolved(mkts2),
					)
				}
				for _, p := range open {
					me := meta[p.AssetID]
					if me == nil || me.ConditionID == "" {
						continue
					}
					m, ok := byCond[me.ConditionID]
					if !ok {
						continue
					}
					if !m.Closed {
						continue
					}
					prices := m.OutcomePrices()
					if me.SlotIdx < 0 || me.SlotIdx >= len(prices) {
						continue
					}
					settleMid, perr := strconv.ParseFloat(prices[me.SlotIdx], 64)
					if perr != nil {
						slog.Warn("settlement_price_parse_fail",
							"asset", short(p.AssetID),
							"raw", prices[me.SlotIdx],
							"err", perr.Error())
						continue
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
						Reason:   strategy.ExitSettlement,
					}
					closed, cerr := pm.Close(p.ID, sig)
					if cerr != nil {
						slog.Warn("settlement_close_miss", "pos", p.ID, "asset", short(p.AssetID), "err", cerr.Error())
						continue
					}
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
						_ = rm.SaveState(riskStatePath)
					}
					if netPnL <= -largeFillUSD || netPnL >= largeFillUSD {
						notifier.LargeFill(notify.LargeFillEvent{
							Question: metaQ(meta, p.AssetID),
							AssetID:  p.AssetID,
							Side:     "sell",
							SizeUSD:  p.SizeUSD,
							PnLUSD:   netPnL,
							EntryPx:  p.EntryMid,
							ExitPx:   settleMid,
							Reason:   string(strategy.ExitSettlement),
							HeldSec:  int(sig.HeldFor.Seconds()),
						})
					}
					source, openOID := src.Take(closed.ID)
					if jerr := jrn.Append(journal.TradeRecord{
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
						"q", metaQ(meta, p.AssetID),
						"outcome", metaOutcome(meta, p.AssetID),
						"entry", p.EntryMid,
						"settle", settleMid,
						"gross_pnl_usd", closed.PnLUSD,
						"entry_fee_usd", entryFeeShare,
						"net_pnl_usd", netPnL,
						"held_sec", int(sig.HeldFor.Seconds()),
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
			if err != nil || msgID == 0 {
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
			_ = h.rm.SaveState(h.riskStatePath)
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
//
// SIDECAR_BOT_TOKEN, when set, routes SignalPrompt messages through that bot so
// the inline-keyboard message originates from the same bot the LongPoll watches
// for callback_query (Phase 3.5.b). Other events stay on the alert bot.
func buildNotifier() notify.Notifier {
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		slog.Info("notify.ready", "mode", "nop", "reason", "telegram_env_missing")
		return notify.Nop{}
	}
	cfg := notify.TelegramConfig{BotToken: tok, ChatID: chat, PromptBotToken: os.Getenv("SIDECAR_BOT_TOKEN")}
	slog.Info("notify.ready",
		"mode", "telegram",
		"chat_id", chat,
		"prompt_via_sidecar", cfg.PromptBotToken != "",
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

// runArbScan runs a single arb scan cycle (one-shot CLI mode).
func runArbScan(ctx context.Context, dbPath string, minGapPP float64) error {
	store, err := arb.NewStore(dbPath)
	if err != nil {
		return fmt.Errorf("arb store: %w", err)
	}
	defer store.Close()

	oddsClient := odds.NewClient("", "")
	cfg := arb.DefaultScanConfig()
	cfg.MinGapPP = minGapPP
	scanner := arb.NewScanner(oddsClient, store, cfg)

	opps, err := scanner.Scan(ctx)
	if err != nil {
		return err
	}

	usage := oddsClient.Usage()
	fmt.Printf("\narb-scan: %d opportunities (min_gap=%.0fpp)\n", len(opps), minGapPP)
	fmt.Printf("API usage: %d used, %d remaining\n", usage.RequestsUsed, usage.RequestsRemaining)
	fmt.Printf("DB rows: %d\n\n", store.Count())

	for _, o := range opps {
		fmt.Printf("  %s | %s | %s\n", o.Sport, o.EventName, o.Direction)
		fmt.Printf("    poly=%.3f bk=%.3f gap=%+.1fpp net_ev=%+.1fpp | %s\n",
			o.PolymarketPrice, o.BookmakerProb, o.GapPP, o.NetEvPP, o.MarketTitle)
	}
	return nil
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
		return err
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

		nonce := p.Nonce
		notifier.SignalPrompt(notify.SignalPromptEvent{
			Nonce:   p.Nonce,
			Match:   me.Match,
			Context: fmt.Sprintf("🚨 %s %s OUT · %s", a.Team, a.StarPlayer, a.Reason),
			EndIn:   notify.HumanizeEndIn(time.Now(), me.EndTime),
			Slug:    me.Slug,
			Choices: sigChoices,
			OnSent: func(msgID int64, err error) {
				if err != nil || msgID == 0 {
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
