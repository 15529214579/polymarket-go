package btc

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// CoinConfig defines parameters for a crypto price-level strategy.
type CoinConfig struct {
	Name         string   // display name: "ETH", "SOL"
	BinancePair  string   // e.g. "ETHUSDT", "SOLUSDT"
	GammaSlug    string   // PM event slug
	MinStrike    float64  // ignore strikes below this (filters noise)
	VolFloorPct  float64  // minimum annualized vol (e.g. 0.30 = 30%)
}

var (
	CoinETH = CoinConfig{
		Name:        "ETH",
		BinancePair: "ETHUSDT",
		GammaSlug:   "what-price-will-ethereum-hit-before-2027",
		MinStrike:   100,
		VolFloorPct: 0.30,
	}
	CoinSOL = CoinConfig{
		Name:        "SOL",
		BinancePair: "SOLUSDT",
		GammaSlug:   "what-price-will-solana-hit-before-2027",
		MinStrike:   5,
		VolFloorPct: 0.35,
	}
)

// FetchCoinMarkets retrieves price-level markets for any coin from the Gamma API.
func FetchCoinMarkets(ctx context.Context, slug string) ([]PMMarket, error) {
	url := fmt.Sprintf("%s?slug=%s&closed=false", gammaEventsURL, slug)
	return fetchMarketsFromURL(ctx, url)
}

// ScanCoinOnce runs the full BS gap pipeline for one coin.
// Returns signals and a summary log line.
func ScanCoinOnce(ctx context.Context, db *sql.DB, coin CoinConfig, cfg StrategyConfig) ([]Signal, error) {
	candles1h, err := FetchCandles(ctx, coin.BinancePair, Interval1h, 720)
	if err != nil {
		return nil, fmt.Errorf("%s fetch 1h candles: %w", coin.Name, err)
	}
	if len(candles1h) < 24 {
		return nil, fmt.Errorf("%s insufficient 1h candles: %d", coin.Name, len(candles1h))
	}

	if err := saveCoinCandles(ctx, db, coin.Name, candles1h, Interval1h); err != nil {
		slog.Warn("multicoin.save_candles_fail", "coin", coin.Name, "err", err.Error())
	}

	spot := candles1h[len(candles1h)-1].Close
	sigmaHist := HistoricalVolatility(candles1h)
	sigmaEWMA := EWMAVolatility(candles1h, 0.94)
	sigma := BlendedVolatility(candles1h, 0.94, 0.6)
	if sigma < coin.VolFloorPct {
		sigma = coin.VolFloorPct
	}

	markets, err := FetchCoinMarkets(ctx, coin.GammaSlug)
	if err != nil {
		return nil, fmt.Errorf("%s fetch PM markets: %w", coin.Name, err)
	}

	if err := saveCoinPMPrices(ctx, db, coin.Name, markets); err != nil {
		slog.Warn("multicoin.save_pm_fail", "coin", coin.Name, "err", err.Error())
	}

	regime, regimeLabel, regimeConf := DetectCurrentRegime(candles1h)

	yearsToExpiry := yearsUntilEnd2026()
	gaps := FindBSGaps(markets, spot, sigma, yearsToExpiry, cfg.MinGapPP)

	slog.Info("multicoin.scan_done",
		"coin", coin.Name,
		"spot", spot,
		"sigma_hist", fmt.Sprintf("%.1f%%", sigmaHist*100),
		"sigma_ewma", fmt.Sprintf("%.1f%%", sigmaEWMA*100),
		"sigma_blended", fmt.Sprintf("%.1f%%", sigma*100),
		"pm_markets", len(markets),
		"gaps_found", len(gaps),
		"regime", fmt.Sprintf("%s(conf=%.2f)", regimeLabel, regimeConf),
	)

	var signals []Signal
	for i, g := range gaps {
		if i >= cfg.TopN {
			break
		}

		isReach := g.Strike > spot
		regBias := RegimeDirectionBias(regime, regimeConf, "n/a", g.Direction, isReach)
		score := ScoreSignal(g.GapPP, 1.0, regBias, "n/a", 0, g.EdgeRatio)

		sig := Signal{
			Coin:         coin.Name,
			Strike:       g.Strike,
			Question:     g.Question,
			MarketID:     marketIDForStrike(markets, g.Strike),
			PMPrice:      g.PMPrice,
			BSProb:       g.BSProb,
			GapPP:        g.GapPP,
			Direction:    g.Direction,
			EdgeRatio:    g.EdgeRatio,
			Spot:         spot,
			Sigma:        sigma,
			SentimentMod: 1.0,
			RegimeBias:   regBias,
			Score:        score,
		}
		signals = append(signals, sig)

		slog.Info("multicoin.signal",
			"coin", coin.Name,
			"strike", g.Strike,
			"pm_price", fmt.Sprintf("%.3f", g.PMPrice),
			"bs_prob", fmt.Sprintf("%.3f", g.BSProb),
			"gap_pp", fmt.Sprintf("%.1f", g.GapPP),
			"direction", g.Direction,
			"regime", fmt.Sprintf("%s(%.2f)", regimeLabel, regimeConf),
			"regime_bias", fmt.Sprintf("%.2f", regBias),
			"score", fmt.Sprintf("%d/%s", score.Total, score.Tier),
		)
	}

	return signals, nil
}

// saveCoinCandles stores candles for non-BTC coins in a separate table.
func saveCoinCandles(ctx context.Context, db *sql.DB, coin string, candles []Candle, interval Interval) error {
	ddl := `CREATE TABLE IF NOT EXISTS coin_candles (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		coin      TEXT    NOT NULL,
		timestamp INTEGER NOT NULL,
		open      REAL    NOT NULL,
		high      REAL    NOT NULL,
		low       REAL    NOT NULL,
		close     REAL    NOT NULL,
		volume    REAL    NOT NULL,
		interval  TEXT    NOT NULL DEFAULT '1h',
		UNIQUE(coin, timestamp, interval)
	);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx,
		`INSERT OR REPLACE INTO coin_candles(coin, timestamp, open, high, low, close, volume, interval) VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range candles {
		if _, err := stmt.ExecContext(ctx, coin, c.Timestamp.Unix(), c.Open, c.High, c.Low, c.Close, c.Volume, string(interval)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// saveCoinPMPrices stores PM price snapshots for non-BTC coins.
func saveCoinPMPrices(ctx context.Context, db *sql.DB, coin string, markets []PMMarket) error {
	ddl := `CREATE TABLE IF NOT EXISTS coin_pm_prices (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		coin       TEXT    NOT NULL,
		timestamp  INTEGER NOT NULL,
		market_id  TEXT    NOT NULL,
		question   TEXT,
		strike     REAL,
		yes_price  REAL,
		no_price   REAL
	);
	CREATE INDEX IF NOT EXISTS coin_pm_prices_ts ON coin_pm_prices(coin, timestamp);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	now := time.Now().Unix()
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO coin_pm_prices(coin, timestamp, market_id, question, strike, yes_price, no_price) VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range markets {
		if _, err := stmt.ExecContext(ctx, coin, now, m.MarketID, m.Question, m.Strike, m.YesPrice, m.NoPrice); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// fetchMarketsFromURL is the generic Gamma API fetcher (shared with FetchBTCMarkets).
func fetchMarketsFromURL(ctx context.Context, url string) ([]PMMarket, error) {
	req, err := newGammaRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	return doGammaFetch(req)
}
