package btc

import (
	"math"
	"testing"
	"time"
)

func TestMacroStateNormal(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC) // no event nearby
	state := GetMacroState(now)
	if state.Phase != "normal" {
		t.Fatalf("expected normal, got %s", state.Phase)
	}
	if state.VolMultiplier != 1.0 {
		t.Fatalf("expected vol_mult=1.0, got %.3f", state.VolMultiplier)
	}
}

func TestMacroStatePreFOMC(t *testing.T) {
	fomc := time.Date(2026, 5, 7, 18, 0, 0, 0, time.UTC)
	now := fomc.Add(-12 * time.Hour) // 12h before FOMC
	state := GetMacroState(now)
	if state.Phase != "pre_event" {
		t.Fatalf("expected pre_event, got %s", state.Phase)
	}
	if state.NextEvent.Type != EventFOMC {
		t.Fatalf("expected FOMC, got %s", state.NextEvent.Type)
	}
	// 12h before, ramp=0.5, bump=0.30, mult=1.15
	if math.Abs(state.VolMultiplier-1.15) > 0.01 {
		t.Fatalf("expected vol_mult≈1.15, got %.3f", state.VolMultiplier)
	}
}

func TestMacroStatePreCPI(t *testing.T) {
	cpi := time.Date(2026, 5, 13, 12, 30, 0, 0, time.UTC)
	now := cpi.Add(-6 * time.Hour) // 6h before CPI
	state := GetMacroState(now)
	if state.Phase != "pre_event" {
		t.Fatalf("expected pre_event, got %s", state.Phase)
	}
	// 6h before, ramp=0.75, bump=0.25, mult=1.1875
	if state.VolMultiplier < 1.15 || state.VolMultiplier > 1.22 {
		t.Fatalf("expected vol_mult≈1.19, got %.3f", state.VolMultiplier)
	}
}

func TestMacroStatePostEvent(t *testing.T) {
	cpi := time.Date(2026, 5, 13, 12, 30, 0, 0, time.UTC)
	now := cpi.Add(30 * time.Minute) // 30m after CPI
	state := GetMacroState(now)
	if state.Phase != "post_event" {
		t.Fatalf("expected post_event, got %s", state.Phase)
	}
	if state.VolMultiplier >= 1.0 {
		t.Fatalf("expected vol_mult<1.0 for cooldown, got %.3f", state.VolMultiplier)
	}
}

func TestMacroVolAdjust(t *testing.T) {
	baseVol := 0.27
	macro := MacroState{VolMultiplier: 1.25}
	adj := MacroVolAdjust(baseVol, macro)
	if math.Abs(adj-0.3375) > 0.001 {
		t.Fatalf("expected 0.3375, got %.4f", adj)
	}
}

func TestMacroFindNextEvent(t *testing.T) {
	now := time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC)
	state := GetMacroState(now)
	if state.NextEvent == nil {
		t.Fatal("expected to find Dec 23 PCE")
	}
	if state.NextEvent.Type != EventPCE {
		t.Fatalf("expected PCE, got %s", state.NextEvent.Type)
	}
}
