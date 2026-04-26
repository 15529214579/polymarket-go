package strategy

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

type fakeSampler struct {
	last map[string]feed.Tick
}

func (f *fakeSampler) TickTail(id string, n int) ([]feed.Tick, bool) {
	t, ok := f.last[id]
	if !ok {
		return nil, false
	}
	return []feed.Tick{t}, true
}

func TestIsEligible_GlobalBand(t *testing.T) {
	cfg := DefaultLotteryConfig()
	cases := []struct {
		name  string
		sport SportFamily
		mid   float64
		want  bool
	}{
		{"basketball lower edge", SportBasketball, 0.05, true},
		{"basketball upper edge", SportBasketball, 0.30, true},
		{"basketball above ceiling", SportBasketball, 0.305, false},
		{"basketball below floor", SportBasketball, 0.049, false},
		{"football 0.20 ok", SportFootball, 0.20, true},
		{"zero mid rejected", SportBasketball, 0, false},
		{"negative mid rejected", SportBasketball, -0.1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsEligible(cfg, tc.sport, tc.mid)
			if got != tc.want {
				t.Fatalf("IsEligible(sport=%s, mid=%v) = %v, want %v", tc.sport, tc.mid, got, tc.want)
			}
		})
	}
}

func TestIsEligible_LoLHasTighterFloor(t *testing.T) {
	cfg := DefaultLotteryConfig()
	// LoL floor is 0.15 — a 0.08 LoL underdog is "压倒性低概率" and must be skipped.
	if IsEligible(cfg, SportLoL, 0.08) {
		t.Fatal("LoL @ 0.08 should be rejected (predictable stomp)")
	}
	// Same 0.08 for NBA is fine — upsets happen.
	if !IsEligible(cfg, SportBasketball, 0.08) {
		t.Fatal("NBA @ 0.08 should be eligible")
	}
	// LoL at floor 0.15 is eligible.
	if !IsEligible(cfg, SportLoL, 0.15) {
		t.Fatal("LoL @ 0.15 should be eligible (at floor)")
	}
}

func TestEffectiveFloor(t *testing.T) {
	cfg := DefaultLotteryConfig()
	if got := EffectiveFloor(cfg, SportLoL); got != 0.15 {
		t.Fatalf("LoL floor: got %v, want 0.15", got)
	}
	if got := EffectiveFloor(cfg, SportBasketball); got != 0.05 {
		t.Fatalf("NBA floor: got %v, want 0.05", got)
	}
	// When LoLMinPrice is lower than MinPrice, global floor wins.
	cfg.LoLMinPrice = 0.01
	if got := EffectiveFloor(cfg, SportLoL); got != 0.05 {
		t.Fatalf("LoL floor w/ LoLMin<Global: got %v, want 0.05", got)
	}
}

func TestScanEligible(t *testing.T) {
	now := time.Now()
	fs := &fakeSampler{last: map[string]feed.Tick{
		"lol_tough":  {AssetID: "lol_tough", Market: "M1", Mid: 0.08, Time: now},  // LoL stomp, reject
		"lol_close":  {AssetID: "lol_close", Market: "M2", Mid: 0.18, Time: now},  // LoL in band, accept
		"nba_upset":  {AssetID: "nba_upset", Market: "M3", Mid: 0.08, Time: now},  // NBA upset, accept
		"nba_heavy":  {AssetID: "nba_heavy", Market: "M4", Mid: 0.72, Time: now},  // NBA favorite, reject (above ceiling)
		"epl_mid":    {AssetID: "epl_mid", Market: "M5", Mid: 0.22, Time: now},    // football, accept
		"no_sampler": {AssetID: "no_sampler", Market: "M6", Mid: 0.10, Time: now}, // classified but sampler sees nothing below
	}}
	delete(fs.last, "no_sampler")
	assetSport := map[string]SportFamily{
		"lol_tough":  SportLoL,
		"lol_close":  SportLoL,
		"nba_upset":  SportBasketball,
		"nba_heavy":  SportBasketball,
		"epl_mid":    SportFootball,
		"no_sampler": SportFootball,
	}
	cfg := DefaultLotteryConfig()
	got := ScanEligible(fs, assetSport, cfg)
	if len(got) != 3 {
		t.Fatalf("expected 3 eligible, got %d: %+v", len(got), got)
	}
	seen := map[string]bool{}
	for _, c := range got {
		seen[c.AssetID] = true
		if c.Mid <= 0 {
			t.Fatalf("candidate %s has non-positive mid %v", c.AssetID, c.Mid)
		}
		if c.Sport == SportUnknown {
			t.Fatalf("candidate %s has unknown sport", c.AssetID)
		}
	}
	for _, want := range []string{"lol_close", "nba_upset", "epl_mid"} {
		if !seen[want] {
			t.Fatalf("expected %s in eligible set, missing; got %+v", want, got)
		}
	}
	for _, reject := range []string{"lol_tough", "nba_heavy", "no_sampler"} {
		if seen[reject] {
			t.Fatalf("unexpected eligible: %s", reject)
		}
	}
}

func TestScanEligible_EmptyInputs(t *testing.T) {
	if got := ScanEligible(nil, nil, DefaultLotteryConfig()); got != nil {
		t.Fatalf("nil sampler should yield nil, got %+v", got)
	}
	fs := &fakeSampler{last: map[string]feed.Tick{}}
	if got := ScanEligible(fs, nil, DefaultLotteryConfig()); got != nil {
		t.Fatalf("nil assetSport should yield nil, got %+v", got)
	}
}

func TestIsVolatile(t *testing.T) {
	cases := []struct {
		name string
		ws   feed.WindowStats
		want bool
	}{
		{"calm pre-match", feed.WindowStats{Samples: 60, DeltaPP: 0.5, Upticks: 3, Downticks: 2}, false},
		{"high delta in-play", feed.WindowStats{Samples: 60, DeltaPP: 6.0, Upticks: 8, Downticks: 5}, true},
		{"negative delta in-play", feed.WindowStats{Samples: 60, DeltaPP: -5.5, Upticks: 4, Downticks: 6}, true},
		{"choppy back-and-forth", feed.WindowStats{Samples: 60, DeltaPP: 1.0, Upticks: 12, Downticks: 10}, true},
		{"borderline uptick+downtick=19", feed.WindowStats{Samples: 60, DeltaPP: 2.0, Upticks: 10, Downticks: 9}, false},
		{"too few samples", feed.WindowStats{Samples: 5, DeltaPP: 8.0, Upticks: 3, Downticks: 2}, false},
		{"exact threshold delta=5", feed.WindowStats{Samples: 60, DeltaPP: 5.0, Upticks: 5, Downticks: 3}, true},
		{"exact threshold churn=20", feed.WindowStats{Samples: 60, DeltaPP: 1.0, Upticks: 10, Downticks: 10}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vr := IsVolatile(tc.ws)
			if vr.Volatile != tc.want {
				t.Fatalf("IsVolatile(%+v).Volatile = %v, want %v", tc.ws, vr.Volatile, tc.want)
			}
		})
	}
}

func TestClassifySport_FromMarket(t *testing.T) {
	lol := feed.Market{Question: "LoL: T1 vs Hanwha - LCK Spring"}
	nba := feed.Market{Slug: "nba-atl-nyk-2026-04-21"}
	epl := feed.Market{Slug: "epl-ars-che-2026-04-21"}
	other := feed.Market{Question: "Bitcoin > 100k by 2026?"}
	if got := ClassifySport(lol); got != SportLoL {
		t.Fatalf("LoL: got %v", got)
	}
	if got := ClassifySport(nba); got != SportBasketball {
		t.Fatalf("NBA: got %v", got)
	}
	if got := ClassifySport(epl); got != SportFootball {
		t.Fatalf("EPL: got %v", got)
	}
	if got := ClassifySport(other); got != SportUnknown {
		t.Fatalf("non-sport: got %v", got)
	}
}
