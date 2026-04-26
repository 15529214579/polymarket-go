package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func TestNotifySL_ExtendsCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownPerAsset = 5 * time.Minute
	cfg.CooldownAfterSL = 15 * time.Minute
	sampler := feed.NewSampler(60)
	det := NewDetector(cfg, sampler)

	asset := "asset-abc"

	// Simulate a normal fire
	det.lastFire[asset] = time.Now()

	// After 6 minutes, normal cooldown would have expired
	det.lastFire[asset] = time.Now().Add(-6 * time.Minute)

	// Now simulate SL: NotifySL should push the cooldown to 15 min from now
	det.NotifySL(asset, "")

	// Check that lastFire was set such that 15 min cooldown applies
	elapsed := time.Since(det.lastFire[asset])
	if elapsed > 0 {
		t.Errorf("expected lastFire to be in the future after NotifySL, got elapsed=%v", elapsed)
	}

	expiresAt := det.lastFire[asset].Add(cfg.CooldownPerAsset)
	remaining := time.Until(expiresAt)
	if remaining < 14*time.Minute || remaining > 16*time.Minute {
		t.Errorf("expected ~15 min cooldown remaining, got %v", remaining)
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

	det.NotifySL(asset, "")

	if det.lastFire[asset] != original {
		t.Error("NotifySL with CooldownAfterSL=0 should not modify lastFire")
	}
}

func TestNotifySL_PerMarketCooldown(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownPerAsset = 5 * time.Minute
	cfg.CooldownAfterSL = 1 * time.Hour
	sampler := feed.NewSampler(60)
	det := NewDetector(cfg, sampler)

	assetA := "magic-win"
	assetB := "pistons-win"
	condID := "cond-0x123"
	det.RegisterMarket(condID, []string{assetA, assetB})

	// SL on asset A should cool down both A and B
	det.NotifySL(assetA, condID)

	for _, id := range []string{assetA, assetB} {
		expiresAt := det.lastFire[id].Add(cfg.CooldownPerAsset)
		remaining := time.Until(expiresAt)
		if remaining < 59*time.Minute || remaining > 61*time.Minute {
			t.Errorf("asset %s: expected ~1h cooldown, got %v", id, remaining)
		}
	}
}

func TestNotifySL_UnknownMarketFallback(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CooldownPerAsset = 5 * time.Minute
	cfg.CooldownAfterSL = 1 * time.Hour
	sampler := feed.NewSampler(60)
	det := NewDetector(cfg, sampler)

	asset := "orphan-asset"
	det.NotifySL(asset, "unknown-cond")

	// Should still cool down the single asset
	expiresAt := det.lastFire[asset].Add(cfg.CooldownPerAsset)
	remaining := time.Until(expiresAt)
	if remaining < 59*time.Minute || remaining > 61*time.Minute {
		t.Errorf("expected ~1h cooldown for orphan asset, got %v", remaining)
	}

	// Other assets should not be affected
	if _, exists := det.lastFire["other-asset"]; exists {
		t.Error("unrelated asset should not have cooldown")
	}
}
