package strategy

import (
	"math"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func shadowPosition(now time.Time) Position {
	return Position{
		ID: "p1", AssetID: "asset", Market: "market",
		SizeUSD: 10, Units: 20, OpenFeeUSD: 0.10,
		EntryMid: 0.50, EntryTime: now, Status: PosOpen,
		HoldProfile: HoldProfileEvent,
		EventStart:  now.Add(time.Hour), ExitDeadline: now.Add(70 * time.Minute),
	}
}

func TestShadowExitTrackerTimeoutsFireOnce(t *testing.T) {
	now := time.Now().UTC()
	tracker := NewShadowExitTracker(DefaultShadowExitConfig())
	tracker.Open(shadowPosition(now))

	if got := tracker.OnTick("p1", feed.Tick{Time: now.Add(10 * time.Minute), Mid: 0.55}); len(got) != 0 {
		t.Fatalf("mid-only timeout observations=%+v", got)
	}
	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(10 * time.Minute), BestBid: 0.54, Mid: 0.55})
	if len(got) != 1 || got[0].Policy != "timeout_10m" || got[0].HoldProfile != HoldProfileEvent {
		t.Fatalf("10m observations=%+v", got)
	}
	if duplicate := tracker.OnTick("p1", feed.Tick{Time: now.Add(15 * time.Minute), BestBid: 0.55, Mid: 0.56}); len(duplicate) != 0 {
		t.Fatalf("duplicate observations=%+v", duplicate)
	}
	got = tracker.OnTick("p1", feed.Tick{Time: now.Add(30 * time.Minute), BestBid: 0.59, Mid: 0.60})
	if len(got) != 2 || got[0].Policy != "timeout_20m" || got[1].Policy != "timeout_30m" {
		t.Fatalf("30m observations=%+v", got)
	}
}

func TestShadowExitTrackerContinuesAfterActualClose(t *testing.T) {
	now := time.Now().UTC()
	tracker := NewShadowExitTracker(DefaultShadowExitConfig())
	p := shadowPosition(now)
	tracker.Open(p)
	p.ExitTime = now.Add(12 * time.Minute)
	p.ExitReason = ExitLadderTimeout
	tracker.ActualClose(p)

	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(45 * time.Minute), BestBid: 0.60})
	if len(got) != 4 || got[3].Policy != "timeout_45m" {
		t.Fatalf("post-close observations=%+v", got)
	}
	if !got[3].ActualCloseAt.Equal(p.ExitTime) || got[3].ActualReason != ExitLadderTimeout {
		t.Fatalf("actual close context=%+v", got[3])
	}
	got = tracker.OnTick("p1", feed.Tick{Time: now.Add(60 * time.Minute), BestBid: 0.62})
	if len(got) != 1 || got[0].Policy != "timeout_60m" {
		t.Fatalf("60m observations=%+v", got)
	}
	if _, ok := tracker.Snapshot()["p1"]; ok {
		t.Fatal("tracker retained state after longest timeout")
	}
}

func TestShadowExitTrackerNetPnLIncludesSlippageAndBothFees(t *testing.T) {
	now := time.Now().UTC()
	cfg := DefaultShadowExitConfig()
	cfg.Timeouts = nil
	cfg.StopLosses = nil
	cfg.TakeProfits = []float64{0.30}
	cfg.SlippageBp = 50
	cfg.FlatFeeBp = 10
	cfg.TakerFeeRate = 0.05
	tracker := NewShadowExitTracker(cfg)
	tracker.Open(shadowPosition(now))

	got := tracker.OnTick("p1", feed.Tick{Time: now.Add(time.Second), BestBid: 0.65})
	if len(got) != 1 {
		t.Fatalf("observations=%+v", got)
	}
	obs := got[0]
	wantFill := 0.65 * 0.995
	wantGross := (wantFill - 0.50) * 20
	wantPlatformFee := math.Round((20*0.05*wantFill*(1-wantFill))*100_000) / 100_000
	wantExitFee := 20*wantFill*10/10_000 + wantPlatformFee
	wantNet := wantGross - 0.10 - wantExitFee
	if math.Abs(obs.ExitPrice-wantFill) > 1e-9 || math.Abs(obs.ExitFeeUSD-wantExitFee) > 1e-9 || math.Abs(obs.NetPnLUSD-wantNet) > 1e-9 {
		t.Fatalf("costed observation=%+v want fill=%.8f exit_fee=%.8f net=%.8f", obs, wantFill, wantExitFee, wantNet)
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
