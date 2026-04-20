package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func mkTick(asset string, t time.Time, mid float64) feed.Tick {
	return feed.Tick{AssetID: asset, Market: "m", Time: t, Mid: mid, BestBid: mid - 0.005, BestAsk: mid + 0.005}
}

func TestExit_NoPosition(t *testing.T) {
	e := NewExitTracker(DefaultExitConfig())
	if _, fired := e.OnTick(mkTick("A", time.Now(), 0.5)); fired {
		t.Fatalf("no position → must not fire")
	}
}

func TestExit_StopLoss(t *testing.T) {
	cfg := DefaultExitConfig()
	e := NewExitTracker(cfg)
	t0 := time.Now()
	e.Open("A", "m", mkTick("A", t0, 0.60))
	// entry 0.60 → price falls to 0.56 → 4pp down, > 3pp stop
	sig, fired := e.OnTick(mkTick("A", t0.Add(10*time.Second), 0.56))
	if !fired {
		t.Fatalf("expected stop-loss to fire")
	}
	if sig.Reason != ExitStopLoss {
		t.Fatalf("reason=%s want stop_loss", sig.Reason)
	}
	if e.Has("A") {
		t.Fatalf("position should be closed after exit")
	}
}

func TestExit_ReversalDrawdown(t *testing.T) {
	cfg := DefaultExitConfig()
	e := NewExitTracker(cfg)
	t0 := time.Now()
	e.Open("A", "m", mkTick("A", t0, 0.60))
	// price peaks at 0.65, then drops to 0.62 → peak-current = 3pp > 2pp drawdown
	_, _ = e.OnTick(mkTick("A", t0.Add(1*time.Second), 0.63))
	_, _ = e.OnTick(mkTick("A", t0.Add(2*time.Second), 0.65))
	sig, fired := e.OnTick(mkTick("A", t0.Add(3*time.Second), 0.62))
	if !fired {
		t.Fatalf("expected reversal drawdown to fire")
	}
	if sig.Reason != ExitReversalDrawdown {
		t.Fatalf("reason=%s want reversal_drawdown", sig.Reason)
	}
}

func TestExit_ReversalTicks(t *testing.T) {
	cfg := DefaultExitConfig()
	cfg.DrawdownPP = 10.0 // disable drawdown path for this test
	cfg.StopLossPP = 10.0
	e := NewExitTracker(cfg)
	t0 := time.Now()
	e.Open("A", "m", mkTick("A", t0, 0.600))
	// 3 consecutive small downticks, each > MinTickPPForMove (0.1pp)
	_, fired := e.OnTick(mkTick("A", t0.Add(1*time.Second), 0.598))
	if fired {
		t.Fatalf("1 downtick should not fire")
	}
	_, fired = e.OnTick(mkTick("A", t0.Add(2*time.Second), 0.596))
	if fired {
		t.Fatalf("2 downticks should not fire")
	}
	sig, fired := e.OnTick(mkTick("A", t0.Add(3*time.Second), 0.594))
	if !fired {
		t.Fatalf("3 downticks should fire")
	}
	if sig.Reason != ExitReversalTicks {
		t.Fatalf("reason=%s want reversal_ticks", sig.Reason)
	}
}

func TestExit_ReversalTicksResetByUptick(t *testing.T) {
	cfg := DefaultExitConfig()
	cfg.DrawdownPP = 10.0
	cfg.StopLossPP = 10.0
	e := NewExitTracker(cfg)
	t0 := time.Now()
	e.Open("A", "m", mkTick("A", t0, 0.600))
	_, _ = e.OnTick(mkTick("A", t0.Add(1*time.Second), 0.598)) // down 1
	_, _ = e.OnTick(mkTick("A", t0.Add(2*time.Second), 0.600)) // up → reset
	_, _ = e.OnTick(mkTick("A", t0.Add(3*time.Second), 0.598)) // down 1
	_, fired := e.OnTick(mkTick("A", t0.Add(4*time.Second), 0.596)) // down 2 only
	if fired {
		t.Fatalf("uptick should reset consec-down counter")
	}
}

func TestExit_Timeout(t *testing.T) {
	cfg := DefaultExitConfig()
	cfg.MaxHold = 5 * time.Second
	e := NewExitTracker(cfg)
	t0 := time.Now()
	e.Open("A", "m", mkTick("A", t0, 0.60))
	// price unchanged but time elapses past MaxHold
	sig, fired := e.OnTick(mkTick("A", t0.Add(6*time.Second), 0.60))
	if !fired {
		t.Fatalf("expected timeout to fire")
	}
	if sig.Reason != ExitTimeout {
		t.Fatalf("reason=%s want timeout", sig.Reason)
	}
}

func TestExit_IgnoresTicksForUnknownAsset(t *testing.T) {
	e := NewExitTracker(DefaultExitConfig())
	e.Open("A", "m", mkTick("A", time.Now(), 0.60))
	if _, fired := e.OnTick(mkTick("B", time.Now(), 0.50)); fired {
		t.Fatalf("tick for non-open asset should be ignored")
	}
}
