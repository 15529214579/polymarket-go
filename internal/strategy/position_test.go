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
	closed, err := pm.Close("asset-A", exit)
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

func TestPositionManager_DedupeByMarket(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	if _, err := pm.Open("yes-token", "cond-1", tick(0.6, now)); err != nil {
		t.Fatalf("open yes: %v", err)
	}
	// NO side of the SAME market must be rejected.
	if _, err := pm.Open("no-token", "cond-1", tick(0.4, now)); !errors.Is(err, ErrMarketAlreadyOpen) {
		t.Fatalf("want ErrMarketAlreadyOpen, got %v", err)
	}
	// Different market OK.
	if _, err := pm.Open("yes-token-2", "cond-2", tick(0.5, now)); err != nil {
		t.Fatalf("open cond-2: %v", err)
	}
}

func TestPositionManager_DedupeByAsset(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	if _, err := pm.Open("asset-A", "", tick(0.5, now)); err != nil {
		t.Fatalf("open1: %v", err)
	}
	if _, err := pm.Open("asset-A", "", tick(0.5, now)); !errors.Is(err, ErrAssetAlreadyOpen) {
		t.Fatalf("want ErrAssetAlreadyOpen, got %v", err)
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

func TestPositionManager_ReopenAfterClose(t *testing.T) {
	pm := NewPositionManager(DefaultPositionConfig())
	now := time.Now()
	if _, err := pm.Open("a", "m", tick(0.5, now)); err != nil {
		t.Fatal(err)
	}
	if _, err := pm.Close("a", ExitSignal{AssetID: "a", Market: "m", Time: now.Add(time.Minute), EntryMid: 0.5, ExitMid: 0.55, Reason: ExitReversalTicks}); err != nil {
		t.Fatal(err)
	}
	// Re-opening the same asset is fine once the previous closed.
	if _, err := pm.Open("a", "m", tick(0.6, now.Add(2*time.Minute))); err != nil {
		t.Fatalf("reopen: %v", err)
	}
}
