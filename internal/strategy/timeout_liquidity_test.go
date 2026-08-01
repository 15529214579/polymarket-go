package strategy

import (
	"testing"
	"time"
)

func TestTimeoutLiquidityTrackerThrottlesAndSummarizes(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	tracker := NewTimeoutLiquidityTracker(5 * time.Minute)
	pos := Position{
		ID:                 "p1",
		Market:             "market-1",
		SizeUSD:            20,
		OpenFeeUSD:         0.60,
		EntryFeeChargedUSD: 0.10,
	}

	first, log := tracker.Observe(pos, now)
	if !log || first.Attempts != 1 {
		t.Fatalf("first=%+v log=%v", first, log)
	}
	if got := first.ConservativeNetPnLUSD; got != -20.5 {
		t.Fatalf("conservative net=%v, want -20.5", got)
	}
	second, log := tracker.Observe(pos, now.Add(time.Minute))
	if log || second.Attempts != 2 {
		t.Fatalf("second=%+v log=%v", second, log)
	}
	third, log := tracker.Observe(pos, now.Add(5*time.Minute))
	if !log || third.Attempts != 3 {
		t.Fatalf("third=%+v log=%v", third, log)
	}

	summary := tracker.Summary()
	if summary.Positions != 1 || summary.ExposureUSD != 20 || summary.ConservativeNetPnLUSD != -20.5 {
		t.Fatalf("summary=%+v", summary)
	}
	resolved, ok := tracker.Resolve(pos.ID)
	if !ok || resolved.Attempts != 3 || tracker.Summary().Positions != 0 {
		t.Fatalf("resolved=%+v ok=%v summary=%+v", resolved, ok, tracker.Summary())
	}
}
