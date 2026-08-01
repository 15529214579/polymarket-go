package main

import "testing"

func TestRedeemValuePrefersCurrentValue(t *testing.T) {
	p := dataAPIPosition{Size: 10, CurPrice: 0.9, CurrentValue: 8.75}
	if got := redeemValue(p); got != 8.75 {
		t.Fatalf("redeemValue = %v", got)
	}
}

func TestRedeemValueFallsBackToMark(t *testing.T) {
	p := dataAPIPosition{Size: 10, CurPrice: 1}
	if got := redeemValue(p); got != 10 {
		t.Fatalf("redeemValue = %v", got)
	}
}

func TestRedeemValueIsZeroForLosingPosition(t *testing.T) {
	p := dataAPIPosition{Size: 10, CurPrice: 0, CurrentValue: 0}
	if got := redeemValue(p); got != 0 {
		t.Fatalf("redeemValue = %v", got)
	}
}
