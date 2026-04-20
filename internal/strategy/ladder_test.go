package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func lcfg() LadderConfig {
	// Shorter MaxHold than default so timeout tests run in ns.
	c := DefaultLadderConfig()
	c.MaxHold = time.Minute
	return c
}

func lt(mid float64, at time.Time) feed.Tick {
	return feed.Tick{AssetID: "A", Market: "M", Time: at, Mid: mid}
}

func TestLadder_UnknownPos_NoEmit(t *testing.T) {
	l := NewLadderTracker(lcfg())
	if _, fired := l.OnTick("ghost", lt(0.55, time.Now())); fired {
		t.Fatalf("expected no emit on unknown posID")
	}
}

func TestLadder_TP1_ThenTP2_Chain(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.40}, 100)

	// Below TP1 — nothing.
	if _, fired := l.OnTick("p1", lt(0.44, t0.Add(1*time.Second))); fired {
		t.Fatalf("premature emit below TP1")
	}
	// Exactly TP1 (0.40 × 1.15 = 0.46) — TP1 fires, 50 units close.
	ex, fired := l.OnTick("p1", lt(0.46, t0.Add(2*time.Second)))
	if !fired || ex.Tranche != "t1" || ex.Reason != ExitLadderTP1 {
		t.Fatalf("tp1 miss: fired=%v ex=%+v", fired, ex)
	}
	if absDiff(ex.CloseUnits, 50) > 1e-9 {
		t.Fatalf("tp1 wrong units: got %v want 50", ex.CloseUnits)
	}
	if ex.Final {
		t.Fatalf("tp1 should not be final with frac=0.50")
	}
	// Still below TP2 (0.40×1.30=0.52) — no emit.
	if _, fired := l.OnTick("p1", lt(0.50, t0.Add(3*time.Second))); fired {
		t.Fatalf("premature emit between TP1 and TP2")
	}
	// TP2 fires; closes remaining 50 units.
	ex2, fired := l.OnTick("p1", lt(0.55, t0.Add(4*time.Second)))
	if !fired || ex2.Tranche != "t2" || ex2.Reason != ExitLadderTP2 {
		t.Fatalf("tp2 miss: fired=%v ex=%+v", fired, ex2)
	}
	if absDiff(ex2.CloseUnits, 50) > 1e-9 {
		t.Fatalf("tp2 wrong units: got %v want 50", ex2.CloseUnits)
	}
	if !ex2.Final {
		t.Fatalf("tp2 should be final — nothing left")
	}
	if l.Has("p1") {
		t.Fatalf("tracker should have dropped posID after final tranche")
	}
}

// Price gaps past TP2 on the same tick without an intervening TP1 tick.
// The tracker should still emit TP1 first (closing TP1Frac of initial units)
// and defer TP2 to the next tick.
func TestLadder_GapsPastTP2_StillSplitsTranches(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.40}, 100)

	ex, fired := l.OnTick("p1", lt(0.60, t0.Add(1*time.Second)))
	if !fired || ex.Tranche != "t1" {
		t.Fatalf("expected tp1 first on gap, got %+v", ex)
	}
	// Next tick at same price should now emit TP2.
	ex2, fired := l.OnTick("p1", lt(0.60, t0.Add(2*time.Second)))
	if !fired || ex2.Tranche != "t2" {
		t.Fatalf("expected tp2 on follow-up, got %+v", ex2)
	}
	if !ex2.Final {
		t.Fatalf("second tranche should be final")
	}
}

func TestLadder_StopLoss_ClosesEverything(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 80)

	// SL threshold = 0.50 × (1 - 0.10) = 0.45. Go below.
	ex, fired := l.OnTick("p1", lt(0.44, t0.Add(500*time.Millisecond)))
	if !fired || ex.Tranche != "sl" || ex.Reason != ExitLadderSL {
		t.Fatalf("sl miss: fired=%v ex=%+v", fired, ex)
	}
	if absDiff(ex.CloseUnits, 80) > 1e-9 {
		t.Fatalf("sl should close all remaining, got %v", ex.CloseUnits)
	}
	if !ex.Final || l.Has("p1") {
		t.Fatalf("sl should be final and drop state")
	}
}

func TestLadder_Timeout_Fires(t *testing.T) {
	cfg := lcfg()
	cfg.MaxHold = 100 * time.Millisecond
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 40)

	ex, fired := l.OnTick("p1", lt(0.51, t0.Add(150*time.Millisecond)))
	if !fired || ex.Reason != ExitLadderTimeout {
		t.Fatalf("timeout miss: fired=%v ex=%+v", fired, ex)
	}
	if !ex.Final || l.Has("p1") {
		t.Fatalf("timeout should be final")
	}
}

// SL takes priority over timeout on the same tick.
func TestLadder_SL_BeforeTimeout(t *testing.T) {
	cfg := lcfg()
	cfg.MaxHold = 100 * time.Millisecond
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 10)

	ex, fired := l.OnTick("p1", lt(0.40, t0.Add(200*time.Millisecond)))
	if !fired || ex.Reason != ExitLadderSL {
		t.Fatalf("want SL, got %+v", ex)
	}
}

func TestLadder_Forget_DropsState(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 10)
	l.Forget("p1")
	if l.Has("p1") {
		t.Fatalf("forget should drop state")
	}
	if _, fired := l.OnTick("p1", lt(0.60, t0)); fired {
		t.Fatalf("forgotten pos should not emit")
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
