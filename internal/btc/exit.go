package btc

import (
	"context"
	"database/sql"
	"log/slog"
	"math"
	"time"
)

// ExitReason describes why a position should be closed.
type ExitReason string

const (
	ExitGapNarrowed ExitReason = "gap_narrowed"  // BS-PM gap shrunk below threshold
	ExitStopLoss    ExitReason = "stop_loss"      // BTC moved >5% against position
	ExitTimeout     ExitReason = "timeout"        // held >7 days without gap recovery
)

// ExitConfig configures the exit strategy thresholds.
type ExitConfig struct {
	GapCloseThreshPP float64       // close when |gap| drops below this (default 3pp)
	StopLossPct      float64       // close when BTC moves this % against us (default 0.05)
	TimeoutDuration  time.Duration // close after this duration (default 7 days)
}

// DefaultExitConfig returns conservative exit defaults.
func DefaultExitConfig() ExitConfig {
	return ExitConfig{
		GapCloseThreshPP: 3.0,
		StopLossPct:      0.05,
		TimeoutDuration:  7 * 24 * time.Hour,
	}
}

// Position represents an open BTC market position.
type Position struct {
	ID           int64
	MarketID     string
	Strike       float64
	Direction    string  // BUY_YES or BUY_NO
	EntryGapPP   float64 // BS-PM gap at entry
	EntryPMPrice float64 // PM price at entry
	EntrySpot    float64 // BTC spot at entry
	EntrySigma   float64 // vol at entry
	EnteredAt    time.Time
}

// ExitSignal represents a recommendation to close a position.
type ExitSignal struct {
	Position    Position
	Reason      ExitReason
	CurrentGap  float64 // current BS-PM gap in pp
	CurrentSpot float64
	PnLEstimate float64 // rough estimate: (currentPM - entryPM) or (entryPM - currentPM)
}

// ExitCallback is called for each position that should be closed.
type ExitCallback func(sig ExitSignal)

// InitExitDB creates the btc_positions table.
func InitExitDB(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS btc_positions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    market_id       TEXT    NOT NULL,
    strike          REAL    NOT NULL,
    direction       TEXT    NOT NULL,
    entry_gap_pp    REAL    NOT NULL,
    entry_pm_price  REAL    NOT NULL,
    entry_spot      REAL    NOT NULL,
    entry_sigma     REAL    NOT NULL,
    entered_at      INTEGER NOT NULL,
    closed_at       INTEGER,
    close_reason    TEXT,
    close_gap_pp    REAL,
    close_spot      REAL,
    pnl_estimate    REAL,
    UNIQUE(market_id, direction, entered_at)
);
CREATE INDEX IF NOT EXISTS btc_positions_open ON btc_positions(closed_at) WHERE closed_at IS NULL;
`
	_, err := db.Exec(ddl)
	return err
}

// RecordEntry inserts a new open position.
func RecordEntry(ctx context.Context, db *sql.DB, sig Signal) error {
	_, err := db.ExecContext(ctx, `
INSERT OR IGNORE INTO btc_positions(market_id, strike, direction, entry_gap_pp, entry_pm_price, entry_spot, entry_sigma, entered_at)
VALUES(?,?,?,?,?,?,?,?)`,
		sig.MarketID, sig.Strike, sig.Direction, sig.GapPP, sig.PMPrice, sig.Spot, sig.Sigma,
		time.Now().Unix(),
	)
	return err
}

// OpenPositions returns all positions where closed_at IS NULL.
func OpenPositions(ctx context.Context, db *sql.DB) ([]Position, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id, market_id, strike, direction, entry_gap_pp, entry_pm_price, entry_spot, entry_sigma, entered_at
FROM btc_positions WHERE closed_at IS NULL ORDER BY entered_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []Position
	for rows.Next() {
		var p Position
		var ts int64
		if err := rows.Scan(&p.ID, &p.MarketID, &p.Strike, &p.Direction,
			&p.EntryGapPP, &p.EntryPMPrice, &p.EntrySpot, &p.EntrySigma, &ts); err != nil {
			return nil, err
		}
		p.EnteredAt = time.Unix(ts, 0)
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

// ClosePosition marks a position as closed with reason and metrics.
func ClosePosition(ctx context.Context, db *sql.DB, posID int64, reason ExitReason, currentGap, currentSpot, pnl float64) error {
	_, err := db.ExecContext(ctx, `
UPDATE btc_positions SET closed_at=?, close_reason=?, close_gap_pp=?, close_spot=?, pnl_estimate=?
WHERE id=? AND closed_at IS NULL`,
		time.Now().Unix(), string(reason), currentGap, currentSpot, pnl, posID,
	)
	return err
}

// CheckExits evaluates all open positions against current market state and
// returns exit signals for positions that should be closed.
func CheckExits(ctx context.Context, db *sql.DB, markets []PMMarket, spot, sigma, yearsToExpiry float64, cfg ExitConfig) []ExitSignal {
	positions, err := OpenPositions(ctx, db)
	if err != nil {
		slog.Warn("exit.open_positions_fail", "err", err.Error())
		return nil
	}
	if len(positions) == 0 {
		return nil
	}

	marketMap := make(map[string]PMMarket, len(markets))
	for _, m := range markets {
		if m.MarketID != "" {
			marketMap[m.MarketID] = m
		}
	}

	var exits []ExitSignal
	for _, pos := range positions {
		if sig, ok := checkOnePosition(pos, marketMap, spot, sigma, yearsToExpiry, cfg); ok {
			exits = append(exits, sig)
		}
	}
	return exits
}

func checkOnePosition(pos Position, marketMap map[string]PMMarket, spot, sigma, yearsToExpiry float64, cfg ExitConfig) (ExitSignal, bool) {
	m, found := marketMap[pos.MarketID]

	var currentGapPP float64
	var currentPMPrice float64
	if found && m.YesPrice > 0.01 && m.YesPrice < 0.99 {
		adjustedSigma := VolSmileAdjust(sigma, spot, pos.Strike)
		bsProb := FirstPassageProb(spot, pos.Strike, adjustedSigma, yearsToExpiry)
		currentGapPP = (bsProb - m.YesPrice) * 100
		currentPMPrice = m.YesPrice
	}

	pnl := estimatePnL(pos, currentPMPrice)

	// 1. Gap narrowed: alpha exhausted
	if found && math.Abs(currentGapPP) < cfg.GapCloseThreshPP && math.Abs(pos.EntryGapPP) >= cfg.GapCloseThreshPP {
		return ExitSignal{
			Position:    pos,
			Reason:      ExitGapNarrowed,
			CurrentGap:  currentGapPP,
			CurrentSpot: spot,
			PnLEstimate: pnl,
		}, true
	}

	// 2. Stop-loss: BTC moved against our position
	spotMove := (spot - pos.EntrySpot) / pos.EntrySpot
	shouldStop := false
	if pos.Direction == "BUY_YES" && pos.Strike > pos.EntrySpot {
		// Bought Yes on reach market — BTC dropping hurts
		shouldStop = spotMove < -cfg.StopLossPct
	} else if pos.Direction == "BUY_YES" && pos.Strike <= pos.EntrySpot {
		// Bought Yes on dip market — BTC rising hurts
		shouldStop = spotMove > cfg.StopLossPct
	} else if pos.Direction == "BUY_NO" && pos.Strike > pos.EntrySpot {
		// Bought No on reach market — BTC rising hurts
		shouldStop = spotMove > cfg.StopLossPct
	} else if pos.Direction == "BUY_NO" && pos.Strike <= pos.EntrySpot {
		// Bought No on dip market — BTC dropping hurts
		shouldStop = spotMove < -cfg.StopLossPct
	}
	if shouldStop {
		return ExitSignal{
			Position:    pos,
			Reason:      ExitStopLoss,
			CurrentGap:  currentGapPP,
			CurrentSpot: spot,
			PnLEstimate: pnl,
		}, true
	}

	// 3. Timeout: held too long
	if time.Since(pos.EnteredAt) > cfg.TimeoutDuration {
		return ExitSignal{
			Position:    pos,
			Reason:      ExitTimeout,
			CurrentGap:  currentGapPP,
			CurrentSpot: spot,
			PnLEstimate: pnl,
		}, true
	}

	return ExitSignal{}, false
}

func estimatePnL(pos Position, currentPMPrice float64) float64 {
	if currentPMPrice <= 0 {
		return 0
	}
	if pos.Direction == "BUY_YES" {
		return (currentPMPrice - pos.EntryPMPrice) * 100
	}
	// BUY_NO: profit when Yes price drops
	return (pos.EntryPMPrice - currentPMPrice) * 100
}
