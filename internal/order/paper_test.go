package order

import (
	"context"
	"math"
	"testing"
)

func TestPaperSubmitFillsAtMid(t *testing.T) {
	p := NewPaperClient(0)
	r, err := p.Submit(context.Background(), Intent{
		AssetID: "asset-1", Market: "mkt-1",
		Side: Buy, SizeUSD: 5, LimitPx: 0.42, Type: GTC,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != StatusFilled {
		t.Fatalf("want filled, got %s", r.Status)
	}
	wantUnits := 5.0 / 0.42
	if math.Abs(r.FilledSize-wantUnits) > 1e-9 {
		t.Fatalf("units: want %v got %v", wantUnits, r.FilledSize)
	}
	if r.AvgPrice != 0.42 {
		t.Fatalf("avg px: want 0.42 got %v", r.AvgPrice)
	}
	if r.OrderID == "" {
		t.Fatal("empty order id")
	}
}

func TestPaperSlippageBuyWorsensPrice(t *testing.T) {
	p := NewPaperClient(50) // 50bp = 0.5%
	r, _ := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 5, LimitPx: 0.50, Type: GTC,
	})
	// Buy fills 0.50 + 0.50*0.005 = 0.5025
	want := 0.5025
	if math.Abs(r.AvgPrice-want) > 1e-9 {
		t.Fatalf("buy slippage: want %v got %v", want, r.AvgPrice)
	}
}

func TestPaperSlippageSellImprovesThenSuffers(t *testing.T) {
	p := NewPaperClient(100) // 100bp
	r, _ := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Sell, SizeUSD: 5, LimitPx: 0.50, Type: GTC,
	})
	// Sell fills 0.50 - 0.50*0.01 = 0.495
	if math.Abs(r.AvgPrice-0.495) > 1e-9 {
		t.Fatalf("sell slippage: want 0.495 got %v", r.AvgPrice)
	}
}

func TestPaperRejectsOutOfRange(t *testing.T) {
	p := NewPaperClient(0)
	for _, px := range []float64{0, -0.1, 1, 1.5} {
		if _, err := p.Submit(context.Background(), Intent{
			AssetID: "a", Side: Buy, SizeUSD: 5, LimitPx: px,
		}); err == nil {
			t.Fatalf("px %v: want error", px)
		}
	}
}

func TestPaperRejectsBadIntent(t *testing.T) {
	p := NewPaperClient(0)
	cases := []Intent{
		{Side: Buy, SizeUSD: 5, LimitPx: 0.5},                // missing AssetID
		{AssetID: "a", Side: "??", SizeUSD: 5, LimitPx: 0.5}, // bad side
		{AssetID: "a", Side: Buy, SizeUSD: 0, LimitPx: 0.5},  // zero size
	}
	for i, c := range cases {
		if _, err := p.Submit(context.Background(), c); err == nil {
			t.Fatalf("case %d: want error", i)
		}
	}
}

func TestPaperHistoryCopy(t *testing.T) {
	p := NewPaperClient(0)
	for i := 0; i < 3; i++ {
		_, _ = p.Submit(context.Background(), Intent{
			AssetID: "a", Side: Buy, SizeUSD: 5, LimitPx: 0.5,
		})
	}
	h := p.History()
	if len(h) != 3 {
		t.Fatalf("history len: want 3 got %d", len(h))
	}
	// mutating returned slice mustn't affect internal state
	h[0] = Result{}
	h2 := p.History()
	if h2[0].Status != StatusFilled {
		t.Fatal("History returned non-copy")
	}
}

func TestPaperSlippageToEdgeRejects(t *testing.T) {
	// 0.999 BUY with 1000bp slippage lands >= 1.0 → reject
	p := NewPaperClient(1000)
	if _, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 5, LimitPx: 0.999,
	}); err == nil {
		t.Fatal("expected reject when slippage pushes price out of (0,1)")
	}
}
