package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func lcfg() LadderConfig {
	return LadderConfig{
		TP1Pct:  0.15,
		TP1Frac: 0.50,
		TP2Pct:  0.30,
		TP2Frac: 1.00,
		SLPct:   0.05,
		MaxHold: time.Minute,
	}
}

func lt(mid float64, at time.Time) feed.Tick {
	return feed.Tick{AssetID: "A", Market: "M", Time: at, Mid: mid}
}

func executableTick(mid, bid float64, at time.Time) feed.Tick {
	return feed.Tick{AssetID: "A", Market: "M", Time: at, QuoteTime: at, Mid: mid, BestBid: bid, BestBidSize: 1_000}
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
	l.Confirm("p1")
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
	l.Confirm("p1")
	if l.Has("p1") {
		t.Fatalf("tracker should have dropped posID after final tranche")
	}
}

func TestLadder_TP1PartialFillRetriesOnlyRemainder(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.40}, 100)
	ex, fired := l.OnTick("p1", lt(0.46, t0.Add(time.Second)))
	if !fired || ex.CloseUnits != 50 {
		t.Fatalf("first=%+v fired=%v", ex, fired)
	}
	l.Confirm("p1", 20)
	retry, fired := l.OnTick("p1", lt(0.46, t0.Add(2*time.Second)))
	if !fired || retry.Tranche != "t1" || retry.CloseUnits != 30 {
		t.Fatalf("retry=%+v fired=%v", retry, fired)
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
	l.Confirm("p1")
	// Next tick at same price should now emit TP2.
	ex2, fired := l.OnTick("p1", lt(0.60, t0.Add(2*time.Second)))
	if !fired || ex2.Tranche != "t2" {
		t.Fatalf("expected tp2 on follow-up, got %+v", ex2)
	}
	if !ex2.Final {
		t.Fatalf("second tranche should be final")
	}
	l.Confirm("p1")
}

func TestLadder_StopLoss_ClosesEverything(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 80)

	// SL threshold = 0.50 × (1 - 0.05) = 0.475. Go below.
	ex, fired := l.OnTick("p1", lt(0.44, t0.Add(500*time.Millisecond)))
	if !fired || ex.Tranche != "sl" || ex.Reason != ExitLadderSL {
		t.Fatalf("sl miss: fired=%v ex=%+v", fired, ex)
	}
	if absDiff(ex.CloseUnits, 80) > 1e-9 {
		t.Fatalf("sl should close all remaining, got %v", ex.CloseUnits)
	}
	if !ex.Final || !l.Has("p1") {
		t.Fatalf("sl should remain pending until confirmation")
	}
	l.Confirm("p1")
	if l.Has("p1") {
		t.Fatalf("confirmed sl should drop state")
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
	if !ex.Final || !l.Has("p1") {
		t.Fatalf("timeout should remain pending")
	}
	l.Retry("p1")
	if _, fired := l.OnTick("p1", lt(0.51, t0.Add(160*time.Millisecond))); !fired {
		t.Fatalf("failed timeout close should retry")
	}
	l.Confirm("p1")
	if l.Has("p1") {
		t.Fatalf("confirmed timeout should drop state")
	}
}

func TestLadder_EventDeadlineDefersTimeout(t *testing.T) {
	cfg := lcfg()
	cfg.MaxHold = 10 * time.Minute
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	deadline := t0.Add(time.Hour)
	l.OpenWithDeadline("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 40, deadline)

	if ex, fired := l.OnTick("p1", lt(0.51, t0.Add(11*time.Minute))); fired {
		t.Fatalf("event hold timed out early: %+v", ex)
	}
	ex, fired := l.OnTick("p1", lt(0.51, deadline))
	if !fired || ex.Reason != ExitLadderTimeout || !ex.Final {
		t.Fatalf("deadline timeout miss: fired=%v ex=%+v", fired, ex)
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

func TestLadder_FeeAwareTPUsesExecutableBidAndCosts(t *testing.T) {
	cfg := lcfg()
	cfg.TP1Pct = 0.30
	cfg.TP2Pct = 9.99
	cfg.SLPct = 0.35
	cfg.MaxHold = time.Hour
	cfg.FeeAware = true
	cfg.RequireExecutableBid = true
	cfg.SlippageBp = 50
	cfg.TakerFeeRate = 0.05
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	l.OpenPosition(Position{
		ID: "p1", AssetID: "A", Market: "M", EntryTime: t0,
		EntryMid: 0.50, Units: 40, InitUnits: 40, SizeUSD: 20, OpenFeeUSD: 0.50,
	}, 0.05)

	if _, fired := l.OnTick("p1", feed.Tick{Time: t0.Add(time.Second), Mid: 0.90}); fired {
		t.Fatal("fee-aware ladder must not trigger without an executable bid")
	}
	if _, fired := l.OnTick("p1", executableTick(0.68, 0.67, t0.Add(2*time.Second))); fired {
		t.Fatal("gross move reaches 30%, but net return after costs should not")
	}
	ex, fired := l.OnTick("p1", executableTick(0.69, 0.68, t0.Add(3*time.Second)))
	if !fired || ex.Reason != ExitLadderTP1 || ex.CloseUnits != 20 {
		t.Fatalf("fee-aware TP1=%+v fired=%v", ex, fired)
	}
	if absDiff(ex.ExitMid, 0.68*0.995) > 1e-9 || ex.TakerFeeRate != 0.05 {
		t.Fatalf("executable fill estimate=%+v", ex)
	}
}

func TestLadder_StopLossRequiresMinimumHoldAndConfirmation(t *testing.T) {
	cfg := lcfg()
	cfg.TP1Pct = 9.99
	cfg.TP2Pct = 9.99
	cfg.SLPct = 0.35
	cfg.MaxHold = time.Hour
	cfg.RequireExecutableBid = true
	cfg.MinHoldBeforeSL = 2 * time.Minute
	cfg.SLConfirmTime = 30 * time.Second
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.50}, 40)

	if _, fired := l.OnTick("p1", executableTick(0.21, 0.20, t0.Add(time.Minute))); fired {
		t.Fatal("SL fired during minimum hold")
	}
	if _, fired := l.OnTick("p1", executableTick(0.31, 0.30, t0.Add(2*time.Minute))); fired {
		t.Fatal("SL fired before confirmation")
	}
	if _, fired := l.OnTick("p1", executableTick(0.41, 0.40, t0.Add(2*time.Minute+20*time.Second))); fired {
		t.Fatal("recovery should reset SL confirmation")
	}
	if _, fired := l.OnTick("p1", executableTick(0.31, 0.30, t0.Add(3*time.Minute))); fired {
		t.Fatal("reset SL fired immediately")
	}
	ex, fired := l.OnTick("p1", executableTick(0.31, 0.30, t0.Add(3*time.Minute+31*time.Second)))
	if !fired || ex.Reason != ExitLadderSL {
		t.Fatalf("confirmed SL=%+v fired=%v", ex, fired)
	}
}

func TestLadder_FeeAwareTrailingClosesRemainderAfterTP1(t *testing.T) {
	cfg := lcfg()
	cfg.TP1Pct = 0.30
	cfg.TP2Pct = 9.99
	cfg.SLPct = 0.35
	cfg.TrailingPct = 0.20
	cfg.MaxHold = time.Hour
	cfg.FeeAware = true
	cfg.RequireExecutableBid = true
	cfg.SlippageBp = 50
	cfg.TakerFeeRate = 0.05
	l := NewLadderTracker(cfg)
	t0 := time.Now()
	l.OpenPosition(Position{
		ID: "p1", AssetID: "A", Market: "M", EntryTime: t0,
		EntryMid: 0.50, Units: 40, InitUnits: 40, SizeUSD: 20, OpenFeeUSD: 0.50,
	}, 0.05)

	tp1, fired := l.OnTick("p1", executableTick(0.69, 0.68, t0.Add(time.Second)))
	if !fired || tp1.Reason != ExitLadderTP1 {
		t.Fatalf("TP1=%+v fired=%v", tp1, fired)
	}
	l.Confirm("p1")
	if _, fired := l.OnTick("p1", executableTick(0.81, 0.80, t0.Add(2*time.Second))); fired {
		t.Fatal("new high should move the trailing peak, not exit")
	}
	trail, fired := l.OnTick("p1", executableTick(0.71, 0.70, t0.Add(3*time.Second)))
	if !fired || trail.Reason != ExitLadderTrailing || trail.Tranche != "trail" || !trail.Final || trail.CloseUnits != 20 {
		t.Fatalf("trailing exit=%+v fired=%v", trail, fired)
	}
}

func TestLadder_SyncPositionClearsPendingAndRestoresPartialProgress(t *testing.T) {
	l := NewLadderTracker(lcfg())
	t0 := time.Now()
	l.Open("p1", "M", "A", feed.Tick{Time: t0, Mid: 0.40}, 100)
	if _, fired := l.OnTick("p1", lt(0.46, t0.Add(time.Second))); !fired {
		t.Fatal("expected pending TP1")
	}
	l.SyncPosition(Position{
		ID: "p1", AssetID: "A", Market: "M", EntryTime: t0,
		EntryMid: 0.40, InitUnits: 100, Units: 80,
	})
	ex, fired := l.OnTick("p1", lt(0.46, t0.Add(2*time.Second)))
	if !fired || ex.Reason != ExitLadderTP1 || ex.CloseUnits != 30 {
		t.Fatalf("synced TP1=%+v fired=%v", ex, fired)
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
