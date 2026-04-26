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
	Strike    float64
	Question  string
	MarketID  string
	PMPrice   float64 // current PM Yes price
	BSProb    float64 // model fair value
	GapPP     float64 // (BSProb - PMPrice) * 100; positive = BUY_YES
	Direction string  // BUY_YES or BUY_NO
	EdgeRatio float64 // |gap| / PMPrice
	Spot      float64 // BTC spot at signal time
	Sigma     float64 // annualized vol used
}

// SignalCallback is called for each actionable signal. The caller (main.go)
// wires this to push SignalPrompt DMs with buy buttons.
type SignalCallback func(sig Signal)

// RunStrategy is the live BTC strategy loop. It periodically:
// 1. Fetches BTC spot + 1h candles from Binance
// 2. Fetches PM BTC markets from Gamma API
// 3. Computes BS first-passage probabilities vs PM prices
// 4. Fires callback for gaps > MinGapPP
// 5. Persists data to SQLite for daily backtest iteration
func RunStrategy(ctx context.Context, cfg StrategyConfig, cb SignalCallback) error {
	db, err := sql.Open("sqlite", cfg.DBPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("btc strategy db open: %w", err)
	}
	defer db.Close()

	if err := InitDB(db); err != nil {
		return fmt.Errorf("btc strategy db init: %w", err)
	}

	slog.Info("btc_strategy.ready",
		"interval", cfg.ScanInterval.String(),
		"min_gap_pp", cfg.MinGapPP,
		"top_n", cfg.TopN,
	)

	scan := func() {
		signals, scanErr := scanOnce(ctx, db, cfg)
		if scanErr != nil {
			slog.Warn("btc_strategy.scan_fail", "err", scanErr.Error())
			return
		}
		for _, sig := range signals {
			cb(sig)
		}
	}

	scan()

	tk := time.NewTicker(cfg.ScanInterval)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tk.C:
			scan()
		}
	}
}

func scanOnce(ctx context.Context, db *sql.DB, cfg StrategyConfig) ([]Signal, error) {
	candles, err := FetchCandles(ctx, "BTCUSDT", Interval1h, 720)
	if err != nil {
		return nil, fmt.Errorf("fetch candles: %w", err)
	}
	if len(candles) < 24 {
		return nil, fmt.Errorf("insufficient candles: %d", len(candles))
	}

	if err := SaveCandles(ctx, db, candles, Interval1h); err != nil {
		slog.Warn("btc_strategy.save_candles_fail", "err", err.Error())
	}

	spot := candles[len(candles)-1].Close
	sigma := HistoricalVolatility(candles)

	markets, err := FetchBTCMarkets(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch PM markets: %w", err)
	}

	if err := SavePMPrices(ctx, db, markets); err != nil {
		slog.Warn("btc_strategy.save_pm_fail", "err", err.Error())
	}

	yearsToExpiry := yearsUntilEnd2026()

	gaps := FindBSGaps(markets, spot, sigma, yearsToExpiry, cfg.MinGapPP)

	slog.Info("btc_strategy.scan_done",
		"spot", spot,
		"sigma_pct", fmt.Sprintf("%.1f%%", sigma*100),
		"pm_markets", len(markets),
		"gaps_found", len(gaps),
		"min_gap_pp", cfg.MinGapPP,
	)

	var signals []Signal
	for i, g := range gaps {
		if i >= cfg.TopN {
			break
		}
		sig := Signal{
			Strike:    g.Strike,
			Question:  g.Question,
			MarketID:  marketIDForStrike(markets, g.Strike),
			PMPrice:   g.PMPrice,
			BSProb:    g.BSProb,
			GapPP:     g.GapPP,
			Direction: g.Direction,
			EdgeRatio: g.EdgeRatio,
			Spot:      spot,
			Sigma:     sigma,
		}
		signals = append(signals, sig)

		slog.Info("btc_strategy.signal",
			"strike", g.Strike,
			"pm_price", fmt.Sprintf("%.3f", g.PMPrice),
			"bs_prob", fmt.Sprintf("%.3f", g.BSProb),
			"gap_pp", fmt.Sprintf("%.1f", g.GapPP),
			"direction", g.Direction,
			"edge_ratio", fmt.Sprintf("%.2f", g.EdgeRatio),
		)
	}

	return signals, nil
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
