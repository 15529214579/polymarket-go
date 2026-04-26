package btc

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"time"

	_ "modernc.org/sqlite" //nolint:revive // register sqlite driver
)

// StrategyConfig configures the live BTC prediction strategy.
type StrategyConfig struct {
	Enabled      bool
	ScanInterval time.Duration // how often to re-scan PM vs BS gaps
	MinGapPP     float64       // minimum |gap| in percentage points to signal
	TopN         int           // max signals per scan cycle (prevent spam)
	SizeUSD      float64       // default per-signal size hint
	DBPath       string        // SQLite path for BTC data
}

// DefaultStrategyConfig returns sensible defaults.
func DefaultStrategyConfig() StrategyConfig {
	return StrategyConfig{
		ScanInterval: 1 * time.Hour,
		MinGapPP:     7.0,
		TopN:         3,
		SizeUSD:      5.0,
		DBPath:       "db/btc.db",
	}
}

// Signal is one actionable BTC market gap detected by the strategy.
type Signal struct {
	Strike       float64
	Question     string
	MarketID     string
	PMPrice      float64 // current PM Yes price
	BSProb       float64 // model fair value
	GapPP        float64 // (BSProb - PMPrice) * 100; positive = BUY_YES
	Direction    string  // BUY_YES or BUY_NO
	EdgeRatio    float64 // |gap| / PMPrice
	Spot         float64 // BTC spot at signal time
	Sigma        float64 // annualized vol used
	SentimentMod float64 // sentiment multiplier (>1 = amplified, <1 = dampened)
	FearGreed    int     // 0-100 F&G index at signal time
	FundingRate  float64 // perpetual funding rate at signal time
	RegimeBias       float64 // regime direction bias multiplier
	InstitutionalMod float64 // institutional flow modifier
	OnChainMod       float64 // on-chain metrics modifier
	DepthMod         float64 // orderbook depth modifier
	Score            SignalScore
}

// SignalCallback is called for each actionable signal. The caller (main.go)
// wires this to push SignalPrompt DMs with buy buttons.
type SignalCallback func(sig Signal)

// RunStrategy is the live BTC strategy loop. It periodically:
// 1. Fetches BTC spot + 1h candles from Binance
// 2. Fetches PM BTC markets from Gamma API
// 3. Computes BS first-passage probabilities vs PM prices
// 4. Fires callback for gaps > MinGapPP
// 5. Checks open positions for exit conditions
// 6. Persists data to SQLite for daily backtest iteration
func RunStrategy(ctx context.Context, cfg StrategyConfig, cb SignalCallback) error {
	return RunStrategyWithExit(ctx, cfg, cb, nil, DefaultExitConfig())
}

// RunStrategyWithExit is like RunStrategy but also checks open positions for
// exit conditions (gap narrowing, stop-loss, timeout) on each scan cycle.
func RunStrategyWithExit(ctx context.Context, cfg StrategyConfig, cb SignalCallback, exitCb ExitCallback, exitCfg ExitConfig) error {
	db, err := sql.Open("sqlite", cfg.DBPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("btc strategy db open: %w", err)
	}
	defer db.Close()

	if err := InitDB(db); err != nil {
		return fmt.Errorf("btc strategy db init: %w", err)
	}
	if err := InitExitDB(db); err != nil {
		return fmt.Errorf("btc exit db init: %w", err)
	}

	slog.Info("btc_strategy.ready",
		"interval", cfg.ScanInterval.String(),
		"min_gap_pp", cfg.MinGapPP,
		"top_n", cfg.TopN,
		"exit_gap_thresh", exitCfg.GapCloseThreshPP,
		"exit_stop_loss", fmt.Sprintf("%.0f%%", exitCfg.StopLossPct*100),
		"exit_timeout", exitCfg.TimeoutDuration.String(),
	)

	var scanCount int
	scan := func(trigger string) {
		scanCount++
		if trigger != "" {
			slog.Info("btc_strategy.momentum_triggered", "trigger", trigger)
		}
		signals, markets, spot, sigma, yte, scanErr := scanOnceWithState(ctx, db, cfg)
		if scanErr != nil {
			slog.Warn("btc_strategy.scan_fail", "err", scanErr.Error())
			return
		}
		for _, sig := range signals {
			if err := RecordEntry(ctx, db, sig); err != nil {
				slog.Warn("btc_exit.record_entry_fail", "err", err.Error())
			}
			cb(sig)
		}

		if exitCb != nil {
			exits := CheckExits(ctx, db, markets, spot, sigma, yte, exitCfg)
			for _, ex := range exits {
				slog.Info("btc_strategy.exit_signal",
					"position_id", ex.Position.ID,
					"strike", ex.Position.Strike,
					"direction", ex.Position.Direction,
					"reason", string(ex.Reason),
					"entry_gap", fmt.Sprintf("%.1f", ex.Position.EntryGapPP),
					"current_gap", fmt.Sprintf("%.1f", ex.CurrentGap),
					"entry_spot", ex.Position.EntrySpot,
					"current_spot", ex.CurrentSpot,
					"pnl_est", fmt.Sprintf("%.1f", ex.PnLEstimate),
				)
				if err := ClosePosition(ctx, db, ex.Position.ID, ex.Reason, ex.CurrentGap, ex.CurrentSpot, ex.PnLEstimate); err != nil {
					slog.Warn("btc_exit.close_fail", "err", err.Error())
				}
				exitCb(ex)
			}
		}

		if scanCount%6 == 0 {
			go func() {
				c1h, ferr := FetchCandles(ctx, "BTCUSDT", Interval1h, 720)
				if ferr != nil || len(c1h) < 200 {
					return
				}
				CheckDrift(c1h, 168) // 7-day window vs full
			}()
		}
	}

	scan("")

	momentumCh := make(chan struct{}, 1)
	watcher := NewMomentumWatcher(func() {
		select {
		case momentumCh <- struct{}{}:
		default:
		}
	})
	go watcher.Run(ctx)

	tk := time.NewTicker(cfg.ScanInterval)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tk.C:
			scan("")
		case <-momentumCh:
			scan("sharp_move")
		}
	}
}

func scanOnce(ctx context.Context, db *sql.DB, cfg StrategyConfig) ([]Signal, error) {
	signals, _, _, _, _, err := scanOnceWithState(ctx, db, cfg)
	return signals, err
}

func scanOnceWithState(ctx context.Context, db *sql.DB, cfg StrategyConfig) ([]Signal, []PMMarket, float64, float64, float64, error) {
	candles1h, err := FetchCandles(ctx, "BTCUSDT", Interval1h, 720)
	if err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("fetch 1h candles: %w", err)
	}
	if len(candles1h) < 24 {
		return nil, nil, 0, 0, 0, fmt.Errorf("insufficient 1h candles: %d", len(candles1h))
	}

	if err := SaveCandles(ctx, db, candles1h, Interval1h); err != nil {
		slog.Warn("btc_strategy.save_candles_fail", "interval", "1h", "err", err.Error())
	}

	candles5m, err := FetchCandles(ctx, "BTCUSDT", Interval5m, 1000)
	if err != nil {
		slog.Warn("btc_strategy.fetch_5m_fail", "err", err.Error())
	} else if len(candles5m) > 0 {
		if err := SaveCandles(ctx, db, candles5m, Interval5m); err != nil {
			slog.Warn("btc_strategy.save_candles_fail", "interval", "5m", "err", err.Error())
		}
	}

	candles15m, err := FetchCandles(ctx, "BTCUSDT", Interval15m, 1000)
	if err != nil {
		slog.Warn("btc_strategy.fetch_15m_fail", "err", err.Error())
	} else if len(candles15m) > 0 {
		if err := SaveCandles(ctx, db, candles15m, Interval15m); err != nil {
			slog.Warn("btc_strategy.save_candles_fail", "interval", "15m", "err", err.Error())
		}
	}

	spot := candles1h[len(candles1h)-1].Close
	sigmaHist := HistoricalVolatility(candles1h)
	sigmaEWMA := EWMAVolatility(candles1h, 0.94)
	sigma := BlendedVolatility(candles1h, 0.94, 0.6) // 60% EWMA + 40% hist, with 25% floor

	multiTF, err := PredictMultiTF(ctx)
	if err != nil {
		slog.Warn("btc_strategy.multi_tf_fail", "err", err.Error())
	}

	regime, regimeLabel, regimeConf := DetectCurrentRegime(candles1h)

	markets, err := FetchBTCMarkets(ctx)
	if err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("fetch PM markets: %w", err)
	}

	if err := SavePMPrices(ctx, db, markets); err != nil {
		slog.Warn("btc_strategy.save_pm_fail", "err", err.Error())
	}

	trackPMDeltas(ctx, db, markets, spot)

	sentiment := FetchSentiment(ctx)
	if err := SaveSentiment(ctx, db, sentiment); err != nil {
		slog.Warn("btc_strategy.save_sentiment_fail", "err", err.Error())
	}

	instFlow := FetchInstitutionalFlow(ctx)
	if err := SaveInstitutionalFlow(ctx, db, instFlow); err != nil {
		slog.Warn("btc_strategy.save_institutional_fail", "err", err.Error())
	}

	onchain := FetchOnChainMetrics(ctx)
	if err := SaveOnChainMetrics(ctx, db, onchain); err != nil {
		slog.Warn("btc_strategy.save_onchain_fail", "err", err.Error())
	}

	obDepth := FetchOrderbookDepth(ctx, markets)
	if err := SaveOrderbookDepth(ctx, db, obDepth); err != nil {
		slog.Warn("btc_strategy.save_orderbook_fail", "err", err.Error())
	}

	yearsToExpiry := yearsUntilEnd2026()

	gaps := FindBSGaps(markets, spot, sigma, yearsToExpiry, cfg.MinGapPP)

	tfLabel := "n/a"
	if multiTF != nil {
		tfLabel = fmt.Sprintf("%s(bull=%.0f%%,bear=%.0f%%,conf=%.2f)",
			multiTF.Alignment,
			multiTF.CombinedBull*100, multiTF.CombinedBear*100,
			multiTF.Confidence)
	}

	sentLabel := "n/a"
	if sentiment.FearGreed != nil {
		sentLabel = fmt.Sprintf("F&G=%d(%s)", sentiment.FearGreed.Value, sentiment.FearGreed.Classification)
	}
	if sentiment.FundingRate != nil {
		sentLabel += fmt.Sprintf(" FR=%.4f%%", sentiment.FundingRate.Rate*100)
	}

	slog.Info("btc_strategy.scan_done",
		"spot", spot,
		"sigma_hist", fmt.Sprintf("%.1f%%", sigmaHist*100),
		"sigma_ewma", fmt.Sprintf("%.1f%%", sigmaEWMA*100),
		"sigma_blended", fmt.Sprintf("%.1f%%", sigma*100),
		"pm_markets", len(markets),
		"gaps_found", len(gaps),
		"min_gap_pp", cfg.MinGapPP,
		"multi_tf", tfLabel,
		"sentiment", sentLabel,
		"institutional", fmt.Sprintf("OI=%.0f L/S=%.3f prem=%.4f%% %s(%.2f)", instFlow.OpenInterest, instFlow.LongShortRatio, instFlow.FuturesPremium, instFlow.FlowSignal, instFlow.FlowScore),
		"onchain", fmt.Sprintf("mempool=%d fee=$%.2f %s(%.2f)", onchain.MempoolTxs, onchain.AvgFee24h, onchain.OnChainSignal, onchain.OnChainScore),
		"orderbook", fmt.Sprintf("avg_depth=%.2f min_depth=%.2f markets=%d", obDepth.AvgScore, obDepth.MinScore, len(obDepth.Depths)),
		"regime", fmt.Sprintf("%s(conf=%.2f)", regimeLabel, regimeConf),
	)

	var signals []Signal
	for i, g := range gaps {
		if i >= cfg.TopN {
			break
		}

		if multiTF != nil && !multiTF.MultiTFEntryFilter(g.Direction) {
			slog.Info("btc_strategy.tf_filtered",
				"strike", g.Strike,
				"direction", g.Direction,
				"alignment", multiTF.Alignment,
				"combined_bull", fmt.Sprintf("%.1f%%", multiTF.CombinedBull*100),
			)
			continue
		}

		isReach := g.Strike > spot
		sentMod := sentiment.SentimentModifier(g.Direction, isReach)
		instMod := instFlow.InstitutionalModifier(g.Direction, isReach)
		chainMod := onchain.OnChainModifier(g.Direction, isReach)
		depthMod := DepthModifier(obDepth, g.Strike)

		tfAlignment := "n/a"
		tfConf := 0.0
		if multiTF != nil {
			tfAlignment = multiTF.Alignment
			tfConf = multiTF.Confidence
		}
		regBias := RegimeDirectionBias(regime, regimeConf, tfAlignment, g.Direction, isReach)

		var fng int
		var fr float64
		if sentiment.FearGreed != nil {
			fng = sentiment.FearGreed.Value
		}
		if sentiment.FundingRate != nil {
			fr = sentiment.FundingRate.Rate
		}

		score := ScoreSignal(g.GapPP, sentMod, regBias, tfAlignment, tfConf, g.EdgeRatio)

		sig := Signal{
			Strike:           g.Strike,
			Question:         g.Question,
			MarketID:         marketIDForStrike(markets, g.Strike),
			PMPrice:          g.PMPrice,
			BSProb:           g.BSProb,
			GapPP:            g.GapPP,
			Direction:        g.Direction,
			EdgeRatio:        g.EdgeRatio,
			Spot:             spot,
			Sigma:            sigma,
			SentimentMod:     sentMod,
			FearGreed:        fng,
			FundingRate:      fr,
			RegimeBias:       regBias,
			InstitutionalMod: instMod,
			OnChainMod:       chainMod,
			DepthMod:         depthMod,
			Score:            score,
		}
		signals = append(signals, sig)

		slog.Info("btc_strategy.signal",
			"strike", g.Strike,
			"pm_price", fmt.Sprintf("%.3f", g.PMPrice),
			"bs_prob", fmt.Sprintf("%.3f", g.BSProb),
			"gap_pp", fmt.Sprintf("%.1f", g.GapPP),
			"direction", g.Direction,
			"edge_ratio", fmt.Sprintf("%.2f", g.EdgeRatio),
			"multi_tf", tfLabel,
			"regime", fmt.Sprintf("%s(%.2f)", regimeLabel, regimeConf),
			"regime_bias", fmt.Sprintf("%.2f", regBias),
			"sent_mod", fmt.Sprintf("%.2f", sentMod),
			"inst_mod", fmt.Sprintf("%.2f", instMod),
			"chain_mod", fmt.Sprintf("%.2f", chainMod),
			"depth_mod", fmt.Sprintf("%.2f", depthMod),
			"score", fmt.Sprintf("%d/%s", score.Total, score.Tier),
			"fng", fng,
			"funding", fmt.Sprintf("%.6f", fr),
		)
	}

	return signals, markets, spot, sigma, yearsToExpiry, nil
}

func trackPMDeltas(ctx context.Context, db *sql.DB, markets []PMMarket, spot float64) {
	const ddl = `CREATE TABLE IF NOT EXISTS pm_btc_deltas (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp  INTEGER NOT NULL,
		market_id  TEXT    NOT NULL,
		strike     REAL,
		yes_price  REAL,
		btc_spot   REAL,
		UNIQUE(timestamp, market_id)
	);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		slog.Warn("btc_strategy.delta_ddl_fail", "err", err.Error())
		return
	}

	now := time.Now().Unix()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx,
		`INSERT OR IGNORE INTO pm_btc_deltas(timestamp, market_id, strike, yes_price, btc_spot) VALUES(?,?,?,?,?)`)
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, m := range markets {
		stmt.ExecContext(ctx, now, m.MarketID, m.Strike, m.YesPrice, spot) //nolint:errcheck
	}
	tx.Commit() //nolint:errcheck
}

func yearsUntilEnd2026() float64 {
	end := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
	remaining := time.Until(end)
	if remaining <= 0 {
		return 0.01
	}
	return remaining.Hours() / 8760.0
}

func marketIDForStrike(markets []PMMarket, strike float64) string {
	for _, m := range markets {
		if math.Abs(m.Strike-strike) < 1.0 {
			return m.MarketID
		}
	}
	return ""
}
