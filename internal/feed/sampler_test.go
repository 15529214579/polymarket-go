package feed

import (
	"context"
	"testing"
	"time"
)

func TestSamplerEmitsPerSecondTickAndWindow(t *testing.T) {
	s := NewSampler(5)
	books := make(chan BookEvent, 16)
	trades := make(chan TradeEvent, 16)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = s.Run(ctx, books, trades) }()

	asset := "A"
	feedBook := func(bid, ask float64) {
		books <- BookEvent{
			AssetID: asset, Market: "m", Timestamp: time.Now(),
			Bids: []Level{{Price: bid, Size: 100}},
			Asks: []Level{{Price: ask, Size: 100}},
		}
	}

	feedBook(0.80, 0.82)
	trades <- TradeEvent{AssetID: asset, Market: "m", Timestamp: time.Now(), Price: 0.81, Size: 5, Side: "BUY"}

	// wait for 1st tick emission
	waitTick(t, s, time.Millisecond*1500)

	feedBook(0.83, 0.85)
	trades <- TradeEvent{AssetID: asset, Market: "m", Timestamp: time.Now(), Price: 0.84, Size: 3, Side: "BUY"}
	waitTick(t, s, time.Millisecond*1500)

	feedBook(0.86, 0.88)
	waitTick(t, s, time.Millisecond*1500)

	w, ok := s.Window(asset)
	if !ok {
		t.Fatalf("Window expected ok=true")
	}
	if w.Samples < 2 {
		t.Fatalf("expected >=2 samples, got %d", w.Samples)
	}
	if w.EndMid <= w.StartMid {
		t.Fatalf("expected rising mid, start=%.3f end=%.3f", w.StartMid, w.EndMid)
	}
	if w.DeltaPP <= 0 {
		t.Fatalf("expected positive delta_pp, got %.3f", w.DeltaPP)
	}
	if w.Upticks == 0 {
		t.Fatalf("expected upticks>0")
	}
	if w.BuyRatio != 1.0 {
		t.Fatalf("expected buy_ratio=1, got %.3f (buy=%.1f sell=%.1f)",
			w.BuyRatio, w.BuyVol, w.SellVol)
	}
}

func waitTick(t *testing.T, s *Sampler, d time.Duration) {
	t.Helper()
	select {
	case <-s.Ticks():
	case <-time.After(d):
		t.Fatalf("timed out waiting for tick")
	}
}

func TestWindowUnknownAsset(t *testing.T) {
	s := NewSampler(10)
	if _, ok := s.Window("nope"); ok {
		t.Fatalf("expected ok=false for unknown asset")
	}
}

func TestRingWrapsAndStartsFromOldest(t *testing.T) {
	s := NewSampler(3)
	st := s.ensure("A", "m")
	for i := 0; i < 5; i++ {
		st.ring[st.head] = Tick{Mid: float64(i)}
		st.head = (st.head + 1) % 3
		if st.head == 0 {
			st.full = true
		}
	}
	// ring now conceptually holds [2,3,4]
	got := st.orderedTicks(3)
	if len(got) != 3 {
		t.Fatalf("len %d", len(got))
	}
	want := []float64{2, 3, 4}
	for i, v := range want {
		if got[i].Mid != v {
			t.Fatalf("idx %d: want %v got %v", i, v, got[i].Mid)
		}
	}
}
