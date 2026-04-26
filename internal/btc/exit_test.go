package btc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func testExitDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_journal_mode=WAL")
	if err != nil {
		t.Fatal(err)
	}
	if err := InitDB(db); err != nil {
		t.Fatal(err)
	}
	if err := InitExitDB(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestRecordAndOpenPositions(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	sig := Signal{
		MarketID:  "test-market-1",
		Strike:    50000,
		Direction: "BUY_NO",
		GapPP:     -23.5,
		PMPrice:   0.425,
		Spot:      78000,
		Sigma:     0.27,
	}
	if err := RecordEntry(ctx, db, sig); err != nil {
		t.Fatal(err)
	}

	positions, err := OpenPositions(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}
	if positions[0].Strike != 50000 {
		t.Fatalf("expected strike 50000, got %f", positions[0].Strike)
	}
	if positions[0].Direction != "BUY_NO" {
		t.Fatalf("expected BUY_NO, got %s", positions[0].Direction)
	}
}

func TestExitGapNarrowed(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	sig := Signal{
		MarketID:  "market-55k",
		Strike:    55000,
		Direction: "BUY_NO",
		GapPP:     -30.6,
		PMPrice:   0.525,
		Spot:      78000,
		Sigma:     0.27,
	}
	if err := RecordEntry(ctx, db, sig); err != nil {
		t.Fatal(err)
	}

	// BS prob for 55K dip at spot=78K/sigma=0.27/yte=0.67 ≈ 0.21
	// Set PM YesPrice near BS prob so gap < 3pp → triggers gap_narrowed
	markets := []PMMarket{{
		MarketID: "market-55k",
		Strike:   55000,
		YesPrice: 0.20,
	}}

	cfg := DefaultExitConfig()
	exits := CheckExits(ctx, db, markets, 78000, 0.27, 0.67, cfg)
	if len(exits) != 1 {
		t.Fatalf("expected 1 exit, got %d", len(exits))
	}
	if exits[0].Reason != ExitGapNarrowed {
		t.Fatalf("expected gap_narrowed, got %s", exits[0].Reason)
	}
}

func TestExitStopLoss(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	// BUY_NO on dip market (strike < spot) — BTC dropping further hurts us
	sig := Signal{
		MarketID:  "market-50k",
		Strike:    50000,
		Direction: "BUY_NO",
		GapPP:     -25.0,
		PMPrice:   0.425,
		Spot:      78000,
		Sigma:     0.27,
	}
	if err := RecordEntry(ctx, db, sig); err != nil {
		t.Fatal(err)
	}

	// BTC dropped 6% to ~73,300 — stop loss should trigger for BUY_NO on dip
	markets := []PMMarket{{
		MarketID: "market-50k",
		Strike:   50000,
		YesPrice: 0.50,
	}}

	cfg := DefaultExitConfig()
	exits := CheckExits(ctx, db, markets, 73300, 0.27, 0.67, cfg)
	if len(exits) != 1 {
		t.Fatalf("expected 1 exit, got %d", len(exits))
	}
	if exits[0].Reason != ExitStopLoss {
		t.Fatalf("expected stop_loss, got %s", exits[0].Reason)
	}
}

func TestExitTimeout(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	// Use reach market (100K) where BS prob ≈ 0.28 and PM price=0.12 → gap ~16pp (still large, no gap_narrowed)
	eightDaysAgo := time.Now().Add(-8 * 24 * time.Hour).Unix()
	_, err := db.ExecContext(ctx, `
INSERT INTO btc_positions(market_id, strike, direction, entry_gap_pp, entry_pm_price, entry_spot, entry_sigma, entered_at)
VALUES(?,?,?,?,?,?,?,?)`,
		"market-old", 100000, "BUY_YES", 16.0, 0.12, 78000, 0.27, eightDaysAgo)
	if err != nil {
		t.Fatal(err)
	}

	markets := []PMMarket{{
		MarketID: "market-old",
		Strike:   100000,
		YesPrice: 0.12,
	}}

	cfg := DefaultExitConfig()
	exits := CheckExits(ctx, db, markets, 78000, 0.27, 0.67, cfg)
	if len(exits) != 1 {
		t.Fatalf("expected 1 exit, got %d", len(exits))
	}
	if exits[0].Reason != ExitTimeout {
		t.Fatalf("expected timeout, got %s", exits[0].Reason)
	}
}

func TestClosePositionRemovesFromOpen(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	sig := Signal{
		MarketID:  "market-close-test",
		Strike:    45000,
		Direction: "BUY_NO",
		GapPP:     -20.0,
		PMPrice:   0.335,
		Spot:      78000,
		Sigma:     0.27,
	}
	if err := RecordEntry(ctx, db, sig); err != nil {
		t.Fatal(err)
	}

	positions, _ := OpenPositions(ctx, db)
	if len(positions) != 1 {
		t.Fatal("expected 1 open position")
	}

	if err := ClosePosition(ctx, db, positions[0].ID, ExitGapNarrowed, -2.0, 78000, 5.5); err != nil {
		t.Fatal(err)
	}

	positions, _ = OpenPositions(ctx, db)
	if len(positions) != 0 {
		t.Fatalf("expected 0 open positions after close, got %d", len(positions))
	}
}

func TestNoExitWhenGapStillLarge(t *testing.T) {
	db := testExitDB(t)
	defer db.Close()
	ctx := context.Background()

	sig := Signal{
		MarketID:  "market-hold",
		Strike:    50000,
		Direction: "BUY_NO",
		GapPP:     -25.0,
		PMPrice:   0.425,
		Spot:      78000,
		Sigma:     0.27,
	}
	if err := RecordEntry(ctx, db, sig); err != nil {
		t.Fatal(err)
	}

	// Gap still large, no stop-loss, fresh position — should NOT exit
	markets := []PMMarket{{
		MarketID: "market-hold",
		Strike:   50000,
		YesPrice: 0.40,
	}}

	cfg := DefaultExitConfig()
	exits := CheckExits(ctx, db, markets, 78000, 0.27, 0.67, cfg)
	if len(exits) != 0 {
		t.Fatalf("expected 0 exits (should hold), got %d", len(exits))
	}
}
