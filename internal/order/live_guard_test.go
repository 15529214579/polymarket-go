package order

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type guardTestClient struct {
	calls int
}

func (c *guardTestClient) Name() string { return "test-live" }

func (c *guardTestClient) Submit(_ context.Context, in Intent) (Result, error) {
	c.calls++
	return Result{
		OrderID:    "filled",
		Status:     StatusFilled,
		FilledSize: in.SizeUSD / in.LimitPx,
		AvgPrice:   in.LimitPx,
	}, nil
}

type guardResultClient struct {
	result Result
	err    error
	calls  int
}

func (c *guardResultClient) Name() string { return "test-live-result" }

func (c *guardResultClient) Submit(context.Context, Intent) (Result, error) {
	c.calls++
	return c.result, c.err
}

func TestGuardedClientRejectsMissingArmFile(t *testing.T) {
	cfg, _ := testLiveGuardConfig(t)
	inner := &guardTestClient{}
	client, err := NewGuardedClient(inner, cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Submit(context.Background(), testBuyIntent(5))
	if !errors.Is(err, ErrLiveNotArmed) || result.Status != StatusRejected {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if inner.calls != 0 {
		t.Fatalf("inner calls=%d", inner.calls)
	}
}

func TestGuardedClientEnforcesOrderAndSessionLimits(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	cfg.MaxOrderUSD = 6
	cfg.MaxSessionBuyUSD = 10
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)
	inner := &guardTestClient{}
	client, err := NewGuardedClient(inner, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Submit(context.Background(), testBuyIntent(6)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Submit(context.Background(), testBuyIntent(5)); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("expected session limit, got %v", err)
	}
	if _, err := client.Submit(context.Background(), testBuyIntent(7)); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("expected order limit, got %v", err)
	}
	if _, err := client.Submit(context.Background(), testBuyIntent(math.NaN())); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("expected non-finite size rejection, got %v", err)
	}
	if _, err := client.Submit(context.Background(), Intent{
		AssetID: "123", Market: "market", Side: Sell, SizeUSD: 25, SizeShares: 50, LimitPx: 0.5, Type: GTC,
	}); err != nil {
		t.Fatalf("risk-reducing sell must not be blocked by the buy limit: %v", err)
	}
	if inner.calls != 2 {
		t.Fatalf("inner calls=%d", inner.calls)
	}
}

func TestGuardedClientRechecksDisableFileBeforeEverySubmit(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)
	inner := &guardTestClient{}
	client, err := NewGuardedClient(inner, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Submit(context.Background(), testBuyIntent(5)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfg.DisableFile, []byte("disabled\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Submit(context.Background(), testBuyIntent(5)); !errors.Is(err, ErrLiveDisabled) {
		t.Fatalf("expected disabled error, got %v", err)
	}
	if inner.calls != 1 {
		t.Fatalf("inner calls=%d", inner.calls)
	}
}

func TestCheckLiveGuardFailsClosed(t *testing.T) {
	tests := []struct {
		name       string
		wallet     string
		duration   time.Duration
		permission os.FileMode
	}{
		{name: "expired", wallet: "0xabc", duration: -time.Minute, permission: 0o600},
		{name: "wrong wallet", wallet: "0xdef", duration: time.Hour, permission: 0o600},
		{name: "wide permissions", wallet: "0xabc", duration: time.Hour, permission: 0o644},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, now := testLiveGuardConfig(t)
			writeLiveArmWithWallet(t, cfg, tc.wallet, now, tc.duration, tc.permission)
			if err := CheckLiveGuard(cfg); !errors.Is(err, ErrLiveNotArmed) {
				t.Fatalf("expected not armed, got %v", err)
			}
		})
	}
}

func testLiveGuardConfig(t *testing.T) (LiveGuardConfig, time.Time) {
	t.Helper()
	now := time.Date(2026, 8, 2, 1, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	return LiveGuardConfig{
		ArmFile:          filepath.Join(dir, "live.enabled"),
		DisableFile:      filepath.Join(dir, "live.disabled"),
		ExpectedWallet:   "0xabc",
		MaxOrderUSD:      20,
		MaxSessionBuyUSD: 100,
		MaxArmDuration:   24 * time.Hour,
		SessionStatePath: filepath.Join(dir, "session.json"),
		Now:              func() time.Time { return now },
	}, now
}

func TestGuardedClientRestoresSessionLimitAcrossRestart(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	cfg.MaxOrderUSD = 6
	cfg.MaxSessionBuyUSD = 10
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)

	first, err := NewGuardedClient(&guardTestClient{}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := first.Submit(context.Background(), testBuyIntent(6)); err != nil {
		t.Fatal(err)
	}
	second, err := NewGuardedClient(&guardTestClient{}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := second.Submit(context.Background(), testBuyIntent(5)); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("restarted guard should preserve session total: %v", err)
	}
}

func TestGuardedClientRollsBackTerminalNonFillReservation(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	cfg.MaxOrderUSD = 10
	cfg.MaxSessionBuyUSD = 10
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)

	terminal := &guardResultClient{result: Result{Status: StatusExpired}}
	first, err := NewGuardedClient(terminal, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if result, err := first.Submit(context.Background(), testBuyIntent(6)); err != nil || result.Status != StatusExpired {
		t.Fatalf("terminal result=%+v err=%v", result, err)
	}

	second, err := NewGuardedClient(&guardTestClient{}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := second.Submit(context.Background(), testBuyIntent(10)); err != nil {
		t.Fatalf("terminal non-fill should release reservation: %v", err)
	}
}

func TestGuardedClientKeepsUnknownBuyReservedAcrossRestart(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	cfg.MaxOrderUSD = 10
	cfg.MaxSessionBuyUSD = 10
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)

	unknown := &guardResultClient{
		result: Result{Status: StatusPending},
		err:    errors.New("network outcome unknown"),
	}
	first, err := NewGuardedClient(unknown, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if result, err := first.Submit(context.Background(), testBuyIntent(6)); err == nil || result.Status != StatusPending {
		t.Fatalf("pending result=%+v err=%v", result, err)
	}

	second, err := NewGuardedClient(&guardTestClient{}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := second.Submit(context.Background(), testBuyIntent(5)); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("unknown BUY must remain reserved after restart: %v", err)
	}
}

func TestGuardedClientRejectsSymlinkedSessionState(t *testing.T) {
	cfg, now := testLiveGuardConfig(t)
	writeLiveArm(t, cfg, now, 2*time.Hour, 0o600)
	target := filepath.Join(t.TempDir(), "session.json")
	if err := os.WriteFile(target, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, cfg.SessionStatePath); err != nil {
		t.Fatal(err)
	}
	inner := &guardTestClient{}
	client, err := NewGuardedClient(inner, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Submit(context.Background(), testBuyIntent(5)); !errors.Is(err, ErrLiveLimit) {
		t.Fatalf("expected session state rejection, got %v", err)
	}
	if inner.calls != 0 {
		t.Fatalf("inner calls=%d", inner.calls)
	}
}

func writeLiveArm(t *testing.T, cfg LiveGuardConfig, now time.Time, duration time.Duration, permission os.FileMode) {
	t.Helper()
	writeLiveArmWithWallet(t, cfg, cfg.ExpectedWallet, now, duration, permission)
}

func writeLiveArmWithWallet(t *testing.T, cfg LiveGuardConfig, wallet string, now time.Time, duration time.Duration, permission os.FileMode) {
	t.Helper()
	arm := liveArmFile{
		Wallet:    wallet,
		ArmedAt:   now.Add(-time.Minute),
		ExpiresAt: now.Add(duration),
	}
	raw, err := json.Marshal(arm)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfg.ArmFile, raw, permission); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(cfg.ArmFile, permission); err != nil {
		t.Fatal(err)
	}
}

func testBuyIntent(size float64) Intent {
	return Intent{
		AssetID: "123",
		Market:  "market",
		Side:    Buy,
		SizeUSD: size,
		LimitPx: 0.5,
		Type:    GTC,
	}
}
