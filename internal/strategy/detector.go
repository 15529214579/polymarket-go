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
	"log/slog"
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
	WindowSec            int
	MinDeltaPP           float64
	TailLen              int
	MinTailUpticks       int
	MinBuyRatio          float64
	MinSamplesWarm       int           // don't fire until window has this many samples
	CooldownPerAsset     time.Duration // dedup per asset after a fire
	CooldownAfterSL      time.Duration // extended cooldown after a stop-loss exit
	ConfirmDelay         time.Duration // wait this long after signal, re-check before firing
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
		CooldownAfterSL:  1 * time.Hour,
		ConfirmDelay:     10 * time.Second,
	}
}

type pendingConfirm struct {
	signal Signal
	readyAt time.Time
}

type Detector struct {
	cfg     Config
	sampler *feed.Sampler
	out     chan Signal

	lastFire      map[string]time.Time
	marketAssets  map[string][]string // conditionID → []assetID
	pending       map[string]*pendingConfirm // assetID → awaiting confirmation
}

func NewDetector(cfg Config, sampler *feed.Sampler) *Detector {
	return &Detector{
		cfg:          cfg,
		sampler:      sampler,
		out:          make(chan Signal, 64),
		lastFire:     map[string]time.Time{},
		marketAssets: map[string][]string{},
		pending:      map[string]*pendingConfirm{},
	}
}

// RegisterMarket tells the detector which assets belong to the same market
// (conditionID). NotifySL uses this to cool down all outcomes in a market.
func (d *Detector) RegisterMarket(conditionID string, assetIDs []string) {
	d.marketAssets[conditionID] = assetIDs
}

func (d *Detector) Signals() <-chan Signal { return d.out }

// NotifySL extends the cooldown for all assets in the same market after a
// stop-loss exit. conditionID identifies the market; if empty or unknown, falls
// back to cooling only the single assetID.
func (d *Detector) NotifySL(assetID, conditionID string) {
	if d.cfg.CooldownAfterSL <= 0 {
		return
	}
	coolUntil := time.Now().Add(d.cfg.CooldownAfterSL - d.cfg.CooldownPerAsset)
	siblings := d.marketAssets[conditionID]
	if len(siblings) == 0 {
		d.lastFire[assetID] = coolUntil
		return
	}
	for _, id := range siblings {
		d.lastFire[id] = coolUntil
	}
}

func (d *Detector) CooldownAfterSL() time.Duration { return d.cfg.CooldownAfterSL }

func (d *Detector) Run(ctx context.Context) error {
	ticks := d.sampler.Ticks()
	confirmTicker := time.NewTicker(1 * time.Second)
	defer confirmTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t, ok := <-ticks:
			if !ok {
				return nil
			}
			if sig, fired := d.evaluate(t); fired {
				if d.cfg.ConfirmDelay > 0 {
					d.pending[sig.AssetID] = &pendingConfirm{
						signal:  sig,
						readyAt: time.Now().Add(d.cfg.ConfirmDelay),
					}
				} else {
					select {
					case d.out <- sig:
					default:
					}
				}
			}
		case now := <-confirmTicker.C:
			for assetID, pc := range d.pending {
				if now.Before(pc.readyAt) {
					continue
				}
				delete(d.pending, assetID)
				w, ok := d.sampler.Window(assetID)
				if !ok || w.Samples < d.cfg.MinSamplesWarm {
					slog.Info("confirm_reject",
						"asset", assetID[:min(8, len(assetID))],
						"reason", "window_lost",
						"orig_mid", pc.signal.Mid,
					)
					continue
				}
				threshold := pc.signal.Mid - pc.signal.Mid*float64(d.cfg.MinDeltaPP)/100
				if w.EndMid <= threshold {
					slog.Info("confirm_reject",
						"asset", assetID[:min(8, len(assetID))],
						"reason", "price_retrace",
						"orig_mid", fmt.Sprintf("%.4f", pc.signal.Mid),
						"cur_mid", fmt.Sprintf("%.4f", w.EndMid),
						"threshold", fmt.Sprintf("%.4f", threshold),
					)
					continue
				}
				slog.Info("confirm_fire",
					"asset", assetID[:min(8, len(assetID))],
					"orig_mid", fmt.Sprintf("%.4f", pc.signal.Mid),
					"cur_mid", fmt.Sprintf("%.4f", w.EndMid),
					"delay", d.cfg.ConfirmDelay.String(),
				)
				pc.signal.Mid = w.EndMid
				select {
				case d.out <- pc.signal:
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
