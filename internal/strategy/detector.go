// Package strategy implements entry/exit signal detection.
//
// Detector consumes the 1-Hz tick stream + rolling window stats from feed.Sampler
// and emits Signal events when all configured thresholds are met.
//
// SPEC §2 momentum entry:
//   - 60s net delta ≥ 3pp
//   - ≥4 of the last 5 ticks are upticks
//   - 60s buy-volume ratio ≥ 60%
package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

type Signal struct {
	AssetID   string
	Market    string
	Time      time.Time
	Mid       float64
	WindowSec int
	DeltaPP   float64
	TailLen   int
	TailUps   int
	BuyRatio  float64
	Reason    string
}

func (s Signal) String() string {
	return fmt.Sprintf("SIGNAL %s mid=%.3f Δ=%+.2fpp tail=%d/%d buyR=%.2f [%s]",
		s.AssetID[:min(8, len(s.AssetID))], s.Mid, s.DeltaPP,
		s.TailUps, s.TailLen, s.BuyRatio, s.Market)
}

type Config struct {
	WindowSec        int
	MinDeltaPP       float64
	TailLen          int
	MinTailUpticks   int
	MinBuyRatio      float64
	MinSamplesWarm   int           // don't fire until window has this many samples
	CooldownPerAsset time.Duration // dedup per asset after a fire
}

func DefaultConfig() Config {
	return Config{
		WindowSec:        60,
		MinDeltaPP:       3.0,
		TailLen:          5,
		MinTailUpticks:   4,
		MinBuyRatio:      0.60,
		MinSamplesWarm:   30,
		CooldownPerAsset: 5 * time.Minute,
	}
}

type Detector struct {
	cfg     Config
	sampler *feed.Sampler
	out     chan Signal

	lastFire map[string]time.Time
}

func NewDetector(cfg Config, sampler *feed.Sampler) *Detector {
	return &Detector{
		cfg:      cfg,
		sampler:  sampler,
		out:      make(chan Signal, 64),
		lastFire: map[string]time.Time{},
	}
}

func (d *Detector) Signals() <-chan Signal { return d.out }

func (d *Detector) Run(ctx context.Context) error {
	ticks := d.sampler.Ticks()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t, ok := <-ticks:
			if !ok {
				return nil
			}
			if sig, fired := d.evaluate(t); fired {
				select {
				case d.out <- sig:
				default:
				}
			}
		}
	}
}

func (d *Detector) evaluate(t feed.Tick) (Signal, bool) {
	if last, seen := d.lastFire[t.AssetID]; seen && t.Time.Sub(last) < d.cfg.CooldownPerAsset {
		return Signal{}, false
	}
	w, ok := d.sampler.Window(t.AssetID)
	if !ok || w.Samples < d.cfg.MinSamplesWarm {
		return Signal{}, false
	}
	if w.DeltaPP < d.cfg.MinDeltaPP {
		return Signal{}, false
	}
	if w.BuyRatio < d.cfg.MinBuyRatio {
		return Signal{}, false
	}
	tail, ok := d.sampler.TickTail(t.AssetID, d.cfg.TailLen)
	if !ok || len(tail) < d.cfg.TailLen {
		return Signal{}, false
	}
	ups := 0
	prev := tail[0].Mid
	for i := 1; i < len(tail); i++ {
		if tail[i].Mid > prev {
			ups++
		}
		prev = tail[i].Mid
	}
	// The first tick has no prior; count it as up if mid > window start.
	if tail[0].Mid > w.StartMid {
		ups++
	}
	if ups < d.cfg.MinTailUpticks {
		return Signal{}, false
	}
	d.lastFire[t.AssetID] = t.Time
	return Signal{
		AssetID:   t.AssetID,
		Market:    t.Market,
		Time:      t.Time,
		Mid:       t.Mid,
		WindowSec: w.WindowSec,
		DeltaPP:   w.DeltaPP,
		TailLen:   len(tail),
		TailUps:   ups,
		BuyRatio:  w.BuyRatio,
		Reason:    "momentum",
	}, true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
