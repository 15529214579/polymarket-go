package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func TestNotifySL_ExtendsCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownPerAsset = 5 * time.Minute
	cfg.CooldownAfterSL = 30 * time.Minute
	sampler := feed.NewSampler(60)
	det := NewDetector(cfg, sampler)

	asset := "asset-abc"

	// Simulate a normal fire
	det.lastFire[asset] = time.Now()

	// After 6 minutes, normal cooldown would have expired
	det.lastFire[asset] = time.Now().Add(-6 * time.Minute)
	// evaluate would pass the cooldown check (6m > 5m)

	// Now simulate SL: NotifySL should push the cooldown to 30 min from now
	det.NotifySL(asset)

	// Check that lastFire was set such that 30 min cooldown applies
	elapsed := time.Since(det.lastFire[asset])
	// lastFire = now + (30m - 5m) = now + 25m into the future
	// So elapsed should be negative (about -25 min)
	if elapsed > 0 {
		t.Errorf("expected lastFire to be in the future after NotifySL, got elapsed=%v", elapsed)
	}

	// The effective cooldown: lastFire + CooldownPerAsset should be ~30 min from now
	expiresAt := det.lastFire[asset].Add(cfg.CooldownPerAsset)
	remaining := time.Until(expiresAt)
	if remaining < 29*time.Minute || remaining > 31*time.Minute {
		t.Errorf("expected ~30 min cooldown remaining, got %v", remaining)
	}
}

func TestNotifySL_Disabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownAfterSL = 0 // disabled
	sampler := feed.NewSampler(60)
	det := NewDetector(cfg, sampler)

	asset := "asset-xyz"
	det.lastFire[asset] = time.Now().Add(-10 * time.Minute)
	original := det.lastFire[asset]

	det.NotifySL(asset)

	if det.lastFire[asset] != original {
		t.Error("NotifySL with CooldownAfterSL=0 should not modify lastFire")
	}
}
