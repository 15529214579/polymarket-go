package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func shadowPosition(now time.Time) Position {
	return Position{
		ID: "p1", AssetID: "asset", Market: "market",
		EntryMid: 0.50, EntryTime: now, Status: PosOpen,
		HoldProfile: HoldProfileEvent,
		EventStart:  now.Add(time.Hour), ExitDeadline: now.Add(70 * time.Minute),
	}
}

func TestShadowExitTrackerTimeoutsFireOnce(t *testing.T) {
	now := time.Now().UTC()
	tracker := NewShadowExitTracker(DefaultShadowExitConfig())
	tracker.Open(shadowPosition(now))

	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(10 * time.Minute), Mid: 0.55})
	if len(got) != 1 || got[0].Policy != "timeout_10m" || got[0].HoldProfile != HoldProfileEvent {
		t.Fatalf("10m observations=%+v", got)
	}
	if duplicate := tracker.OnTick("p1", feed.Tick{Time: now.Add(15 * time.Minute), Mid: 0.56}); len(duplicate) != 0 {
		t.Fatalf("duplicate observations=%+v", duplicate)
	}
	got = tracker.OnTick("p1", feed.Tick{Time: now.Add(30 * time.Minute), Mid: 0.60})
	if len(got) != 2 || got[0].Policy != "timeout_20m" || got[1].Policy != "timeout_30m" {
		t.Fatalf("30m observations=%+v", got)
	}
}

func TestShadowExitTrackerSLRequiresExecutableBidAndConfirmation(t *testing.T) {
	now := time.Now().UTC()
	cfg := DefaultShadowExitConfig()
	cfg.Timeouts = nil
	cfg.TakeProfits = nil
	tracker := NewShadowExitTracker(cfg)
	tracker.Open(shadowPosition(now))

	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(time.Second), Mid: 0.35}); len(got) != 0 {
		t.Fatalf("mid-only tick triggered stop: %+v", got)
	}
	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(2 * time.Second), BestBid: 0.39, Mid: 0.40}); len(got) != 0 {
		t.Fatalf("first low bid triggered stop: %+v", got)
	}
	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(10 * time.Second), BestBid: 0.41, Mid: 0.42}); len(got) != 0 {
		t.Fatalf("recovery emitted stop: %+v", got)
	}
	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(20 * time.Second), BestBid: 0.39, Mid: 0.40}); len(got) != 0 {
		t.Fatalf("reset low bid triggered stop: %+v", got)
	}
	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(35 * time.Second), BestBid: 0.39, Mid: 0.40})
	if len(got) != 1 || got[0].Policy != "sl_20" {
		t.Fatalf("confirmed observations=%+v", got)
	}
}

func TestShadowExitTrackerTPUsesBid(t *testing.T) {
	now := time.Now().UTC()
	cfg := DefaultShadowExitConfig()
	cfg.Timeouts = nil
	cfg.StopLosses = nil
	tracker := NewShadowExitTracker(cfg)
	tracker.Open(shadowPosition(now))

	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(time.Second), BestBid: 0.64, Mid: 0.66}); len(got) != 0 {
		t.Fatalf("mid crossed but bid did not: %+v", got)
	}
	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(2 * time.Second), BestBid: 0.65, Mid: 0.66})
	if len(got) != 1 || got[0].Policy != "tp_30" {
		t.Fatalf("tp observations=%+v", got)
	}
	tracker.Close("p1")
	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(3 * time.Second), BestBid: 0.80}); len(got) != 0 {
		t.Fatalf("closed position emitted: %+v", got)
	}
}
