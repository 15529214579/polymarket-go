package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrLiveDisabled = errors.New("order: live trading disabled")
	ErrLiveNotArmed = errors.New("order: live trading not armed")
	ErrLiveLimit    = errors.New("order: live trading limit exceeded")
)

// LiveGuardConfig defines fail-closed controls around a live order client.
// The arm file is short-lived and wallet-specific; the disable file always wins.
type LiveGuardConfig struct {
	ArmFile          string
	DisableFile      string
	ExpectedWallet   string
	MaxOrderUSD      float64
	MaxSessionBuyUSD float64
	MaxArmDuration   time.Duration
	Now              func() time.Time
}

type liveArmFile struct {
	Wallet    string    `json:"wallet"`
	ArmedAt   time.Time `json:"armed_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GuardedClient serializes live submissions and re-checks the kill switch,
// arm expiry, wallet binding, and order limits immediately before signing.
type GuardedClient struct {
	inner         Client
	cfg           LiveGuardConfig
	mu            sync.Mutex
	sessionBuyUSD float64
}

func NewGuardedClient(inner Client, cfg LiveGuardConfig) (*GuardedClient, error) {
	if inner == nil {
		return nil, errors.New("order: guarded client requires an inner client")
	}
	var err error
	cfg, err = normalizeLiveGuardConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &GuardedClient{inner: inner, cfg: cfg}, nil
}

// CheckLiveGuard validates the guard configuration and current arm state
// without constructing an order client or signing anything.
func CheckLiveGuard(cfg LiveGuardConfig) error {
	cfg, err := normalizeLiveGuardConfig(cfg)
	if err != nil {
		return err
	}
	return checkLiveArm(cfg, cfg.Now())
}

func normalizeLiveGuardConfig(cfg LiveGuardConfig) (LiveGuardConfig, error) {
	if strings.TrimSpace(cfg.ArmFile) == "" || strings.TrimSpace(cfg.DisableFile) == "" {
		return LiveGuardConfig{}, errors.New("order: live arm and disable files are required")
	}
	if strings.TrimSpace(cfg.ExpectedWallet) == "" {
		return LiveGuardConfig{}, errors.New("order: expected live wallet is required")
	}
	if !finitePositive(cfg.MaxOrderUSD) || !finitePositive(cfg.MaxSessionBuyUSD) {
		return LiveGuardConfig{}, errors.New("order: positive live order and session limits are required")
	}
	if cfg.MaxSessionBuyUSD < cfg.MaxOrderUSD {
		return LiveGuardConfig{}, errors.New("order: live session limit must be at least the order limit")
	}
	if cfg.MaxArmDuration <= 0 {
		return LiveGuardConfig{}, errors.New("order: positive live arm duration is required")
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return cfg, nil
}

func (c *GuardedClient) Name() string { return c.inner.Name() + "-guarded" }

// CheckReady performs the same fail-closed checks used before every order.
func (c *GuardedClient) CheckReady() error {
	return checkLiveArm(c.cfg, c.cfg.Now())
}

func (c *GuardedClient) Submit(ctx context.Context, in Intent) (Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := checkLiveArm(c.cfg, c.cfg.Now()); err != nil {
		return rejectedResult(err), err
	}
	if !finitePositive(in.SizeUSD) {
		err := fmt.Errorf("%w: order size must be finite and positive", ErrLiveLimit)
		return rejectedResult(err), err
	}
	if in.SizeUSD > c.cfg.MaxOrderUSD+1e-9 {
		err := fmt.Errorf("%w: order %.2fU > %.2fU", ErrLiveLimit, in.SizeUSD, c.cfg.MaxOrderUSD)
		return rejectedResult(err), err
	}
	if in.Side == Buy && c.sessionBuyUSD+in.SizeUSD > c.cfg.MaxSessionBuyUSD+1e-9 {
		err := fmt.Errorf("%w: session buys %.2fU + %.2fU > %.2fU", ErrLiveLimit, c.sessionBuyUSD, in.SizeUSD, c.cfg.MaxSessionBuyUSD)
		return rejectedResult(err), err
	}

	result, err := c.inner.Submit(ctx, in)
	if err == nil && in.Side == Buy && result.Status == StatusFilled {
		c.sessionBuyUSD += in.SizeUSD
	}
	return result, err
}

func checkLiveArm(cfg LiveGuardConfig, now time.Time) error {
	if _, err := os.Stat(cfg.DisableFile); err == nil {
		return fmt.Errorf("%w: %s", ErrLiveDisabled, cfg.DisableFile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("%w: cannot check disable file: %v", ErrLiveDisabled, err)
	}

	info, err := os.Lstat(cfg.ArmFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s is missing", ErrLiveNotArmed, cfg.ArmFile)
		}
		return fmt.Errorf("%w: cannot inspect arm file: %v", ErrLiveNotArmed, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: arm file must be a regular file", ErrLiveNotArmed)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("%w: arm file permissions must be 0600", ErrLiveNotArmed)
	}
	if info.Size() <= 0 || info.Size() > 4096 {
		return fmt.Errorf("%w: arm file size is invalid", ErrLiveNotArmed)
	}

	raw, err := os.ReadFile(cfg.ArmFile)
	if err != nil {
		return fmt.Errorf("%w: cannot read arm file: %v", ErrLiveNotArmed, err)
	}
	var arm liveArmFile
	if err := json.Unmarshal(raw, &arm); err != nil {
		return fmt.Errorf("%w: invalid arm file JSON", ErrLiveNotArmed)
	}
	if !strings.EqualFold(strings.TrimSpace(arm.Wallet), strings.TrimSpace(cfg.ExpectedWallet)) {
		return fmt.Errorf("%w: arm file wallet does not match", ErrLiveNotArmed)
	}
	if arm.ArmedAt.IsZero() || arm.ExpiresAt.IsZero() || !arm.ExpiresAt.After(arm.ArmedAt) {
		return fmt.Errorf("%w: invalid arm time window", ErrLiveNotArmed)
	}
	if arm.ArmedAt.After(now.Add(5 * time.Minute)) {
		return fmt.Errorf("%w: armed_at is in the future", ErrLiveNotArmed)
	}
	if !arm.ExpiresAt.After(now) {
		return fmt.Errorf("%w: arm file expired", ErrLiveNotArmed)
	}
	if arm.ExpiresAt.Sub(arm.ArmedAt) > cfg.MaxArmDuration {
		return fmt.Errorf("%w: arm window exceeds %s", ErrLiveNotArmed, cfg.MaxArmDuration)
	}
	return nil
}

func rejectedResult(err error) Result {
	return Result{Status: StatusRejected, Error: err.Error()}
}

func finitePositive(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}
