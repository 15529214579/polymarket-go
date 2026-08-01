package main

import (
	"strings"
	"testing"
)

func TestValidateOperationFlagsRejectsMixedMaintenance(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "wrap plus redeem", err: validateOperationFlags(true, true, false, false, false, false, false)},
		{name: "wrap plus dry run", err: validateOperationFlags(true, false, false, false, false, false, true)},
		{name: "redeem plus readiness", err: validateOperationFlags(false, true, true, false, false, false, false)},
		{name: "readiness plus cancel", err: validateOperationFlags(false, false, true, false, false, true, false)},
		{name: "public readiness plus open orders", err: validateOperationFlags(false, false, false, true, true, false, false)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil || !strings.Contains(tc.err.Error(), "cannot") && !strings.Contains(tc.err.Error(), "standalone") {
				t.Fatalf("expected operation conflict, got %v", tc.err)
			}
		})
	}
}

func TestValidateOperationFlagsAllowsSupportedModes(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "wrap", err: validateOperationFlags(true, false, false, false, false, false, false)},
		{name: "redeem", err: validateOperationFlags(false, true, false, false, false, false, false)},
		{name: "dry order", err: validateOperationFlags(false, false, false, false, false, false, true)},
		{name: "query then cancel", err: validateOperationFlags(false, false, false, false, true, true, false)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err != nil {
				t.Fatal(tc.err)
			}
		})
	}
}

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
