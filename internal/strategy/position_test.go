package strategy

import (
	"errors"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func tick(mid float64, t time.Time) feed.Tick {
	return feed.Tick{Mid: mid, Time: t}
}

func TestPositionManager_OpenAndClose(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	p, err := pm.Open("asset-A", "cond-1", tick(0.40, now))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if p.SizeUSD != 5.0 || p.EntryMid != 0.40 {
		t.Fatalf("unexpected position: %+v", p)
	}
	if !pm.Has("asset-A") || !pm.HasMarket("cond-1") {
		t.Fatalf("expected Has/HasMarket true")
	}

	exit := ExitSignal{
		AssetID: "asset-A", Market: "cond-1", Time: now.Add(2 * time.Minute),
		EntryMid: 0.40, ExitMid: 0.46, Reason: ExitReversalTicks,
	}
	closed, err := pm.Close(p.ID, exit)
	if err != nil {
		t.Fatalf("close: %v", err)
	}
	wantPnL := (5.0 / 0.40) * (0.46 - 0.40) // = 12.5 * 0.06 = 0.75
	if diff := closed.PnLUSD - wantPnL; diff > 1e-9 || diff < -1e-9 {
		t.Fatalf("pnl: got %v want %v", closed.PnLUSD, wantPnL)
	}
	if pm.Has("asset-A") || pm.HasMarket("cond-1") {
		t.Fatal("still tracking after close")
	}
	if s := pm.Stats(); s.Open != 0 || s.Closed != 1 || s.RealizedPnLUSD != wantPnL {
		t.Fatalf("stats: %+v", s)
	}
}

// Stacking: same market on opposite sides may hold concurrent positions.
func TestPositionManager_StacksByMarket(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	yes, err := pm.Open("yes-token", "cond-1", tick(0.6, now))
	if err != nil {
		t.Fatalf("open yes: %v", err)
	}
	no, err := pm.Open("no-token", "cond-1", tick(0.4, now))
	if err != nil {
		t.Fatalf("open no (stacking should be allowed): %v", err)
	}
	if yes.ID == no.ID {
		t.Fatalf("expected distinct position IDs, got %q/%q", yes.ID, no.ID)
	}
	if got := pm.Stats().Open; got != 2 {
		t.Fatalf("expected 2 open, got %d", got)
	}
}

// Stacking: same asset can hold many concurrent positions (used by manual
// click-to-buy when the boss fires 1U then 5U then 10U on one prompt).
func TestPositionManager_StacksByAsset(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	p1, err := pm.Open("asset-A", "", tick(0.5, now))
	if err != nil {
		t.Fatalf("open1: %v", err)
	}
	p2, err := pm.Open("asset-A", "", tick(0.5, now))
	if err != nil {
		t.Fatalf("open2 (stacking should be allowed): %v", err)
	}
	if p1.ID == p2.ID {
		t.Fatalf("expected distinct IDs, got %q/%q", p1.ID, p2.ID)
	}
	if got := pm.Stats().Open; got != 2 {
		t.Fatalf("expected 2 open, got %d", got)
	}
}

// CloseFirstByAsset closes the oldest stacked position and leaves the rest.
func TestPositionManager_CloseFirstByAsset(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	oldest, err := pm.Open("asset-A", "m", tick(0.5, now))
	if err != nil {
		t.Fatal(err)
	}
	newer, err := pm.Open("asset-A", "m", tick(0.55, now.Add(5*time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	closed, err := pm.CloseFirstByAsset("asset-A", ExitSignal{
		AssetID: "asset-A", Market: "m", Time: now.Add(time.Minute),
		EntryMid: 0.5, ExitMid: 0.6, Reason: ExitReversalTicks,
	})
	if err != nil {
		t.Fatalf("close first: %v", err)
	}
	if closed.ID != oldest.ID {
		t.Fatalf("expected oldest %q closed, got %q", oldest.ID, closed.ID)
	}
	if got := pm.Stats().Open; got != 1 {
		t.Fatalf("expected 1 remaining, got %d", got)
	}
	if !pm.Has("asset-A") {
		t.Fatalf("newer %q should still be open", newer.ID)
	}
}

func TestPositionManager_MaxPositionsAndExposure(t *testing.T) {
	cfg := PositionConfig{PerPositionUSD: 5, MaxTotalOpenUSD: 12, MaxOpenPositions: 10}
	pm := NewPositionManager(cfg)
	now := time.Now()
	// 5 + 5 = 10 ≤ 12 OK; next would be 15 > 12 → exposure trip.
	if _, err := pm.Open("a1", "m1", tick(0.5, now)); err != nil {
		t.Fatalf("o1: %v", err)
	}
	if _, err := pm.Open("a2", "m2", tick(0.5, now)); err != nil {
		t.Fatalf("o2: %v", err)
	}
	if _, err := pm.Open("a3", "m3", tick(0.5, now)); !errors.Is(err, ErrMaxExposure) {
		t.Fatalf("want ErrMaxExposure, got %v", err)
	}

	// Now test MaxOpenPositions trip by allowing exposure but capping count.
	cfg2 := PositionConfig{PerPositionUSD: 1, MaxTotalOpenUSD: 1000, MaxOpenPositions: 2}
	pm2 := NewPositionManager(cfg2)
	if _, err := pm2.Open("x1", "", tick(0.5, now)); err != nil {
		t.Fatalf("x1: %v", err)
	}
	if _, err := pm2.Open("x2", "", tick(0.5, now)); err != nil {
		t.Fatalf("x2: %v", err)
	}
	if _, err := pm2.Open("x3", "", tick(0.5, now)); !errors.Is(err, ErrMaxPositions) {
		t.Fatalf("want ErrMaxPositions, got %v", err)
	}
}

func TestPositionManager_InvalidEntry(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	for _, m := range []float64{0, -0.1, 1.0, 1.5} {
		if _, err := pm.Open("a", "m", tick(m, now)); !errors.Is(err, ErrInvalidEntry) {
			t.Fatalf("mid=%v want ErrInvalidEntry got %v", m, err)
		}
	}
}

func TestPositionManager_PartialClose_Tranches(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	// 5 USDC / 0.50 mid = 10 units initial.
	p, err := pm.Open("asset-A", "m", tick(0.50, now))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	_ = pm.SetOpenFee(p.ID, 0.02)

	// Close 4 units at 0.575 — gross PnL = 4 × (0.575 - 0.50) = 0.30.
	ex1 := ExitSignal{Time: now.Add(time.Second), EntryMid: 0.50, ExitMid: 0.575, Reason: "ladder_tp1"}
	tr1, err := pm.PartialClose(p.ID, 4, ex1)
	if err != nil {
		t.Fatalf("partial close: %v", err)
	}
	if absDiff(tr1.PnLUSD, 0.30) > 1e-9 {
		t.Fatalf("tranche1 pnl: got %v want 0.30", tr1.PnLUSD)
	}
	if tr1.Status != PosClosed {
		t.Fatalf("tranche1 must be closed, got %v", tr1.Status)
	}

	// Verify remaining open: 6 units, SizeUSD = 6 × 0.50 = 3.0.
	snap := pm.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expect 1 open remaining, got %d", len(snap))
	}
	rem := snap[0]
	if absDiff(rem.Units, 6) > 1e-9 {
		t.Fatalf("remaining units: got %v want 6", rem.Units)
	}
	if absDiff(rem.SizeUSD, 3.0) > 1e-9 {
		t.Fatalf("remaining size: got %v want 3.0", rem.SizeUSD)
	}
	if rem.InitUnits != 10 {
		t.Fatalf("InitUnits should stay at 10, got %v", rem.InitUnits)
	}

	// Final close at 0.60 — gross = 6 × (0.60 - 0.50) = 0.60.
	ex2 := ExitSignal{Time: now.Add(2 * time.Second), EntryMid: 0.50, ExitMid: 0.60, Reason: "ladder_tp2"}
	tr2, err := pm.PartialClose(p.ID, 6, ex2)
	if err != nil {
		t.Fatalf("final close: %v", err)
	}
	if absDiff(tr2.PnLUSD, 0.60) > 1e-9 {
		t.Fatalf("tranche2 pnl: got %v want 0.60", tr2.PnLUSD)
	}
	if pm.Has("asset-A") {
		t.Fatalf("should be fully closed")
	}
	if s := pm.Stats(); s.Closed != 2 || s.Open != 0 {
		t.Fatalf("stats: %+v", s)
	}
}

// Oversized closeUnits collapses to a full close (no negative remainder).
func TestPositionManager_PartialClose_Oversized(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	p, err := pm.Open("asset-A", "m", tick(0.50, now))
	if err != nil {
		t.Fatal(err)
	}
	// 10 initial units. Ask for 999 — should fully close.
	tr, err := pm.PartialClose(p.ID, 999, ExitSignal{
		Time: now.Add(time.Minute), EntryMid: 0.50, ExitMid: 0.55, Reason: "ladder_tp2",
	})
	if err != nil {
		t.Fatalf("oversized close: %v", err)
	}
	if pm.Has("asset-A") {
		t.Fatal("oversized should close fully")
	}
	// PnL uses the position's full Units (10), not the request (999).
	wantPnL := 10 * (0.55 - 0.50)
	if absDiff(tr.PnLUSD, wantPnL) > 1e-9 {
		t.Fatalf("pnl: got %v want %v", tr.PnLUSD, wantPnL)
	}
}

func TestPositionManager_ReopenAfterClose(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	p, err := pm.Open("a", "m", tick(0.5, now))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pm.Close(p.ID, ExitSignal{AssetID: "a", Market: "m", Time: now.Add(time.Minute), EntryMid: 0.5, ExitMid: 0.55, Reason: ExitReversalTicks}); err != nil {
		t.Fatal(err)
	}
	// Re-opening the same asset is fine once the previous closed.
	if _, err := pm.Open("a", "m", tick(0.6, now.Add(2*time.Minute))); err != nil {
		t.Fatalf("reopen: %v", err)
	}
}
