package order

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"
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

func TestPaperUsesFreshExecutableTopOfBook(t *testing.T) {
	p := NewPaperClient(50)
	now := time.Now().UTC()
	buy, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 20, LimitPx: 0.60,
		PaperReferencePx: 0.50, PaperBestBid: 0.51, PaperBestAsk: 0.52, PaperQuoteAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := 0.52 * 1.005; math.Abs(buy.AvgPrice-want) > 1e-9 {
		t.Fatalf("buy price=%v want=%v", buy.AvgPrice, want)
	}
	if buy.ReferencePrice != 0.52 || buy.ExecutionModel != "top_of_book" {
		t.Fatalf("buy execution=%+v", buy)
	}

	sell, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Sell, SizeUSD: 20, LimitPx: 0.40,
		PaperReferencePx: 0.50, PaperBestBid: 0.48, PaperBestAsk: 0.49, PaperQuoteAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := 0.48 * 0.995; math.Abs(sell.AvgPrice-want) > 1e-9 {
		t.Fatalf("sell price=%v want=%v", sell.AvgPrice, want)
	}
}

func TestPaperFreshQuoteWithoutExecutableSideExpires(t *testing.T) {
	p := NewPaperClient(0)
	result, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 20, LimitPx: 0.60,
		PaperReferencePx: 0.50, PaperBestBid: 0.49, PaperQuoteAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusExpired || result.Error != "no executable top-of-book quote" || result.FilledSize != 0 {
		t.Fatalf("result=%+v", result)
	}
}

func TestPaperStaleQuoteExpires(t *testing.T) {
	p := NewPaperClient(0)
	result, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 20, LimitPx: 0.60,
		PaperReferencePx: 0.50, PaperBestAsk: 0.58, PaperBestAskSize: 100,
		PaperQuoteAt: time.Now().Add(-time.Minute), PaperRequireQuote: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusExpired || result.ExecutionModel != "stale_quote" {
		t.Fatalf("result=%+v", result)
	}
}

func TestPaperNonStrictOrderPreservesStaleQuoteFallback(t *testing.T) {
	p := NewPaperClient(0)
	result, err := p.Submit(context.Background(), Intent{
		AssetID: "asset", Side: Buy, SizeUSD: 5, LimitPx: 0.5,
		PaperReferencePx: 0.49, PaperBestAsk: 0.5,
		PaperQuoteAt: time.Now().Add(-time.Minute),
	})
	if err != nil || result.Status != StatusFilled || result.ExecutionModel != "stale_quote_fallback" || result.AvgPrice != 0.49 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestPaperRequiredQuoteRejectsMissingAndShallowBook(t *testing.T) {
	p := NewPaperClient(0)
	missing, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 20, LimitPx: 0.60,
		PaperReferencePx: 0.50, PaperRequireQuote: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if missing.Status != StatusExpired || missing.ExecutionModel != "quote_required" {
		t.Fatalf("missing quote result=%+v", missing)
	}
	shallow, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 20, LimitPx: 0.60,
		PaperReferencePx: 0.50, PaperBestAsk: 0.50, PaperBestAskSize: 10,
		PaperQuoteAt: time.Now(), PaperRequireQuote: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if shallow.Status != StatusExpired || !strings.Contains(shallow.Error, "insufficient") {
		t.Fatalf("shallow result=%+v", shallow)
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

func TestPaperDynamicTakerFee(t *testing.T) {
	p := NewPaperClientWithFeeModel(0, 0, 0.05)
	r, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Buy, SizeUSD: 5, LimitPx: 0.50, Type: GTC,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 10 shares x 0.05 x 0.50 x 0.50 = 0.125 USDC.
	if math.Abs(r.FeeUSD-0.125) > 1e-9 {
		t.Fatalf("fee: want 0.125 got %v", r.FeeUSD)
	}
}

func TestPaperDynamicTakerFeeMarketOverride(t *testing.T) {
	p := NewPaperClientWithFeeModel(0, 0, 0.05)
	rate := 0.03
	r, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Market: "m", Side: Buy, SizeUSD: 10, LimitPx: 0.50,
		TakerFeeRateOverride: &rate,
	})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r.FeeUSD-0.15) > 1e-9 {
		t.Fatalf("fee: want 0.15 got %v", r.FeeUSD)
	}
	zero := 0.0
	r, err = p.Submit(context.Background(), Intent{
		AssetID: "a", Market: "m", Side: Sell, SizeUSD: 10, LimitPx: 0.50,
		TakerFeeRateOverride: &zero,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.FeeUSD != 0 {
		t.Fatalf("zero-fee override got %v", r.FeeUSD)
	}
}

func TestPaperSellFeeUsesExactShares(t *testing.T) {
	p := NewPaperClientWithFeeModel(100, 0, 0.05)
	r, err := p.Submit(context.Background(), Intent{
		AssetID: "a", Side: Sell, SizeUSD: 4, SizeShares: 10, LimitPx: 0.40, Type: GTC,
	})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r.FilledSize-10) > 1e-9 {
		t.Fatalf("filled shares: want 10 got %v", r.FilledSize)
	}
	// Sell fills at 0.396; dynamic fee rounds to five USDC decimals.
	if math.Abs(r.FeeUSD-0.11959) > 1e-9 {
		t.Fatalf("fee: want 0.11959 got %v", r.FeeUSD)
	}
}

func TestPaperRejectsOutOfRange(t *testing.T) {
	p := NewPaperClient(0)
	for _, px := range []float64{0, -0.1, 1, 1.5, math.NaN(), math.Inf(1)} {
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
		{AssetID: "a", Side: Buy, SizeUSD: math.NaN(), LimitPx: 0.5},
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
