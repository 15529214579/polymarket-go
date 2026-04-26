package btc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	gammaEventsURL = "https://gamma-api.polymarket.com/events"
	// slug pattern used by Polymarket's BTC-price-level markets
	btcSlug = "what-price-will-bitcoin-hit-before-2027"
)

// PMMarket represents one BTC price-level market fetched from the Gamma API.
type PMMarket struct {
	MarketID    string
	Question    string
	Strike      float64 // parsed strike price from question
	YesPrice    float64 // 0..1
	NoPrice     float64 // 0..1
	FetchedAt   time.Time
}

// ImpliedCurvePoint is one point on the implied-probability ladder.
type ImpliedCurvePoint struct {
	Strike   float64
	YesPrice float64 // raw PM price for Yes
	Implied  float64 // adjusted implied probability (accounting for PM fee)
}

// ---------------------------------------------------------------------------
// Gamma API client
// ---------------------------------------------------------------------------

type gammaEvent struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Slug     string        `json:"slug"`
	Markets  []gammaMarket `json:"markets"`
}

type gammaMarket struct {
	ID           string `json:"id"`
	Question     string `json:"question"`
	OutcomeYes   string `json:"outcomePrices"`
	Outcomes     string `json:"outcomes"`
	ClobTokenIDs string `json:"clobTokenIds"`
}

// FetchBTCMarkets retrieves all active BTC price-level markets from
// the Polymarket Gamma API.
func FetchBTCMarkets(ctx context.Context) ([]PMMarket, error) {
	url := fmt.Sprintf("%s?slug=%s&closed=false", gammaEventsURL, btcSlug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch BTC markets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("gamma HTTP %d: %s", resp.StatusCode, body)
	}

	var events []gammaEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		// Gamma sometimes returns a single object instead of array
		return nil, fmt.Errorf("decode events: %w", err)
	}

	now := time.Now().UTC()
	var markets []PMMarket
	for _, ev := range events {
		for _, m := range ev.Markets {
			strike := parseStrikeFromQuestion(m.Question)
			yesPrice, noPrice := parsePrices(m.OutcomeYes)
			markets = append(markets, PMMarket{
				MarketID:  m.ID,
				Question:  m.Question,
				Strike:    strike,
				YesPrice:  yesPrice,
				NoPrice:   noPrice,
				FetchedAt: now,
			})
		}
	}

	// Sort by strike ascending for curve display.
	sort.Slice(markets, func(i, j int) bool {
		return markets[i].Strike < markets[j].Strike
	})
	return markets, nil
}

// parseStrikeFromQuestion tries to extract a dollar amount from the question
// string, e.g. "Will Bitcoin reach $120,000 before 2027?" → 120000.
func parseStrikeFromQuestion(q string) float64 {
	q = strings.ReplaceAll(q, ",", "")
	for _, field := range strings.Fields(q) {
		field = strings.Trim(field, "$?.!,")
		if strings.HasPrefix(field, "$") {
			field = field[1:]
		}
		// try removing trailing K/k
		multiplier := 1.0
		upper := strings.ToUpper(field)
		if strings.HasSuffix(upper, "K") {
			multiplier = 1000
			field = field[:len(field)-1]
		}
		f, err := strconv.ParseFloat(field, 64)
		if err == nil && f > 1000 {
			return f * multiplier
		}
	}
	return 0
}

// parsePrices parses the OutcomeYes/OutcomePrices JSON field.
// Gamma encodes it as a JSON array string: e.g. "[\"0.72\",\"0.28\"]"
func parsePrices(raw string) (yes, no float64) {
	// Strip outer quotes if present (double-encoded JSON)
	raw = strings.TrimSpace(raw)
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		raw = strings.ReplaceAll(raw, `\"`, `"`)
	}

	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		// try float array
		var farr []float64
		if err2 := json.Unmarshal([]byte(raw), &farr); err2 == nil && len(farr) >= 2 {
			return farr[0], farr[1]
		}
		return 0, 0
	}
	if len(arr) >= 1 {
		yes, _ = strconv.ParseFloat(arr[0], 64)
	}
	if len(arr) >= 2 {
		no, _ = strconv.ParseFloat(arr[1], 64)
	}
	return yes, no
}

// ---------------------------------------------------------------------------
// Implied probability curve
// ---------------------------------------------------------------------------

// BuildImpliedCurve converts raw PM prices to an implied-probability curve.
// PM Yes price ≈ implied probability (after the ~2% fee haircut).
func BuildImpliedCurve(markets []PMMarket) []ImpliedCurvePoint {
	pts := make([]ImpliedCurvePoint, 0, len(markets))
	for _, m := range markets {
		if m.Strike <= 0 || m.YesPrice <= 0 {
			continue
		}
		// Polymarket takes ~2% fee; adjust raw price upward to get implied prob.
		implied := m.YesPrice / 0.98
		if implied > 1 {
			implied = 1
		}
		pts = append(pts, ImpliedCurvePoint{
			Strike:   m.Strike,
			YesPrice: m.YesPrice,
			Implied:  implied,
		})
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].Strike < pts[j].Strike })
	return pts
}

// ---------------------------------------------------------------------------
// Gap analysis: PM implied vs Markov-predicted probability
// ---------------------------------------------------------------------------

// MarkovVsPM compares the Markov model's predicted BTC direction (bull/bear)
// against PM market implied probabilities, returning opportunities where the
// gap exceeds minGapPP percentage points.
type GapOpportunity struct {
	Strike        float64
	YesPrice      float64 // raw PM Yes price
	ImpliedProb   float64 // implied PM probability
	MarkovBullPct float64 // model's predicted bull probability (%)
	GapPP         float64 // markov_bull_pct - implied_prob*100
	Direction     string  // "BUY_YES" or "BUY_NO"
}

// FindGaps compares a Markov prediction against the implied probability curve.
// btcPredBull is the model's predicted probability that BTC is bullish (0..1).
func FindGaps(curve []ImpliedCurvePoint, btcPredBull float64, minGapPP float64) []GapOpportunity {
	var opps []GapOpportunity
	for _, pt := range curve {
		markovBullPct := btcPredBull * 100
		impliedPct := pt.Implied * 100
		gap := markovBullPct - impliedPct
		absGap := math.Abs(gap)
		if absGap < minGapPP {
			continue
		}
		dir := "BUY_YES"
		if gap < 0 {
			dir = "BUY_NO"
		}
		opps = append(opps, GapOpportunity{
			Strike:        pt.Strike,
			YesPrice:      pt.YesPrice,
			ImpliedProb:   pt.Implied,
			MarkovBullPct: markovBullPct,
			GapPP:         gap,
			Direction:     dir,
		})
	}
	sort.Slice(opps, func(i, j int) bool {
		return math.Abs(opps[i].GapPP) > math.Abs(opps[j].GapPP)
	})
	return opps
}

// ---------------------------------------------------------------------------
// SQLite persistence
// ---------------------------------------------------------------------------

// InitDB creates the required tables in db if they do not exist.
func InitDB(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS btc_candles (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,          -- unix seconds UTC
    open      REAL    NOT NULL,
    high      REAL    NOT NULL,
    low       REAL    NOT NULL,
    close     REAL    NOT NULL,
    volume    REAL    NOT NULL,
    interval  TEXT    NOT NULL DEFAULT '1h',
    UNIQUE(timestamp, interval)
);

CREATE TABLE IF NOT EXISTS pm_btc_prices (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp  INTEGER NOT NULL,         -- unix seconds UTC
    market_id  TEXT    NOT NULL,
    question   TEXT,
    strike     REAL,
    yes_price  REAL,
    no_price   REAL
);
CREATE INDEX IF NOT EXISTS pm_btc_prices_ts ON pm_btc_prices(timestamp);

CREATE TABLE IF NOT EXISTS btc_predictions (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp        INTEGER NOT NULL,   -- unix seconds UTC
    predicted_state  INTEGER NOT NULL,
    actual_state     INTEGER,            -- filled in retrospectively
    predicted_return REAL,
    actual_return    REAL,               -- filled in retrospectively
    bull_prob        REAL,
    bear_prob        REAL
);
CREATE INDEX IF NOT EXISTS btc_predictions_ts ON btc_predictions(timestamp);
`
	_, err := db.Exec(ddl)
	return err
}

// SaveCandles upserts candles into btc_candles.
func SaveCandles(ctx context.Context, db *sql.DB, candles []Candle, interval Interval) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx, `
INSERT OR REPLACE INTO btc_candles(timestamp, open, high, low, close, volume, interval)
VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range candles {
		if _, err := stmt.ExecContext(ctx,
			c.Timestamp.Unix(), c.Open, c.High, c.Low, c.Close, c.Volume, string(interval),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// SavePMPrices stores PM market snapshots.
func SavePMPrices(ctx context.Context, db *sql.DB, markets []PMMarket) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO pm_btc_prices(timestamp, market_id, question, strike, yes_price, no_price)
VALUES(?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range markets {
		if _, err := stmt.ExecContext(ctx,
			m.FetchedAt.Unix(), m.MarketID, m.Question, m.Strike, m.YesPrice, m.NoPrice,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// SavePrediction stores one Markov prediction row.
func SavePrediction(ctx context.Context, db *sql.DB, ts time.Time, pred Prediction) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO btc_predictions(timestamp, predicted_state, predicted_return, bull_prob, bear_prob)
VALUES(?,?,?,?,?)`,
		ts.Unix(), pred.CurrentState, pred.ExpectedReturn, pred.BullProb, pred.BearProb,
	)
	return err
}

// LoadCandles retrieves candles from btc_candles table for the given interval.
func LoadCandles(ctx context.Context, db *sql.DB, interval Interval) ([]Candle, error) {
	rows, err := db.QueryContext(ctx, `
SELECT timestamp, open, high, low, close, volume
FROM btc_candles
WHERE interval = ?
ORDER BY timestamp ASC`, string(interval))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []Candle
	for rows.Next() {
		var ts int64
		var c Candle
		if err := rows.Scan(&ts, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume); err != nil {
			return nil, err
		}
		c.Timestamp = time.Unix(ts, 0).UTC()
		candles = append(candles, c)
	}
	return candles, rows.Err()
}
