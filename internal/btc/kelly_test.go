package btc

import (
	"math"
	"testing"
)

func TestKellyFraction(t *testing.T) {
	// Buy at 0.40 when fair is 0.50 → positive edge
	f := KellyFraction(0.40, 0.50)
	if f <= 0 {
		t.Errorf("expected positive Kelly for underpriced bet, got %f", f)
	}
	if f > 0.5 {
		t.Errorf("half-Kelly should be <= 0.5, got %f", f)
	}

	// Buy at 0.50 when fair is 0.50 → zero edge
	f2 := KellyFraction(0.50, 0.50)
	if f2 != 0 {
		t.Errorf("expected zero Kelly at fair price, got %f", f2)
	}

	// Buy at 0.60 when fair is 0.50 → negative edge
	f3 := KellyFraction(0.60, 0.50)
	if f3 != 0 {
		t.Errorf("expected zero Kelly for overpriced bet, got %f", f3)
	}
}

func TestKellySizeUSD(t *testing.T) {
	size := KellySizeUSD(90.0, 0.40, 0.50, 15.0)
	if size <= 0 {
		t.Errorf("expected positive size, got %f", size)
	}
	if size > 15.0 {
		t.Errorf("should be capped at maxBet, got %f", size)
	}

	// No edge → no bet
	size2 := KellySizeUSD(90.0, 0.50, 0.50, 15.0)
	if size2 != 0 {
		t.Errorf("expected zero size at fair price, got %f", size2)
	}
}

func TestExpectedValue(t *testing.T) {
	// Buy at 0.40 when fair is 0.50
	ev := ExpectedValue(0.40, 0.50)
	if ev <= 0 {
		t.Errorf("expected positive EV, got %f", ev)
	}
	expected := 0.50*(1.0/0.40-1.0) - 0.50
	if math.Abs(ev-expected) > 0.001 {
		t.Errorf("EV mismatch: got %f, want %f", ev, expected)
	}

	// Fair price → zero EV
	ev2 := ExpectedValue(0.50, 0.50)
	if math.Abs(ev2) > 0.001 {
		t.Errorf("expected ~zero EV at fair price, got %f", ev2)
	}
}

func TestValueEdge(t *testing.T) {
	edge := ValueEdge(0.45, 0.50)
	if math.Abs(edge-5.0) > 0.001 {
		t.Errorf("expected 5pp edge, got %f", edge)
	}
}
