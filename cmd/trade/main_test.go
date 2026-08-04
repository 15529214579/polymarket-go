package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/order"
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

func TestRunRedeemAllRequiresStatePath(t *testing.T) {
	err := runRedeemAll(nil, projectWalletAddress, "  ")
	if err == nil || !strings.Contains(err.Error(), "state path is required") {
		t.Fatalf("expected state path error, got %v", err)
	}
}

func TestRunRedeemAllRejectsPaperStatePath(t *testing.T) {
	err := runRedeemAll(nil, projectWalletAddress, "db/redeemed.json")
	if err == nil || !strings.Contains(err.Error(), "under") {
		t.Fatalf("expected live-state isolation error, got %v", err)
	}
}

func TestValidateLiveRuntimeStatePathRejectsCrossModeAndNestedPaths(t *testing.T) {
	for _, path := range []string{"db/paper/orders.sqlite", "db/live/nested/orders.sqlite"} {
		if _, err := validateLiveRuntimeStatePath(path, "execution ledger"); err == nil {
			t.Fatalf("expected path %q to be rejected", path)
		}
	}
}

func TestValidateLiveRuntimeStatePathRejectsSymlinkedParent(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	target := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(target, "live"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(dir, "db")); err != nil {
		t.Fatal(err)
	}
	if _, err := validateLiveRuntimeStatePath("db/live/orders.sqlite", "execution ledger"); err == nil {
		t.Fatal("expected symlinked parent rejection")
	}
}

func TestApplyRecoveredTradeFillsPersistsOriginalFillTime(t *testing.T) {
	ledger, err := order.OpenExecutionLedger(filepath.Join(t.TempDir(), "orders.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()
	client, err := order.NewLedgerClient(order.NewPaperClient(0), ledger, "live")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), order.Intent{
		AssetID: "asset-1", Side: order.Buy, SizeUSD: 20, LimitPx: 0.5, Type: order.FAK,
	})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "buy_times.json")
	if err := applyRecoveredTradeFills(ledger, path); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]string
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	got, err := time.Parse(time.RFC3339Nano, data["asset-1"])
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(result.FilledAt) {
		t.Fatalf("buy time=%s fill time=%s", got, result.FilledAt)
	}
	records, err := ledger.UnappliedFills("live")
	if err != nil || len(records) != 0 {
		t.Fatalf("unapplied=%+v err=%v", records, err)
	}
}
