package strategy

import (
	"fmt"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// ExitReason identifies which exit rule fired.
type ExitReason string

const (
	ExitReversalTicks    ExitReason = "reversal_ticks"    // N consecutive downticks from peak
	ExitReversalDrawdown ExitReason = "reversal_drawdown" // peak - current ≥ drawdown_pp
	ExitStopLoss         ExitReason = "stop_loss"         // entry - current ≥ stop_pp
	ExitTimeout          ExitReason = "timeout"           // held longer than max_hold
)

// ExitSignal is emitted when an open position should be closed.
type ExitSignal struct {
	AssetID    string
	Market     string
	Time       time.Time
	EntryMid   float64
	PeakMid    float64
	ExitMid    float64
	HeldFor    time.Duration
	ChangePP   float64 // (exit - entry) in pp
	DrawdownPP float64 // (peak - exit) in pp
	Reason     ExitReason
}

func (e ExitSignal) String() string {
	id := e.AssetID
	if len(id) > 8 {
		id = id[:8]
	}
	return fmt.Sprintf("EXIT %s reason=%s entry=%.3f peak=%.3f exit=%.3f Δ=%+.2fpp dd=%.2fpp held=%s",
		id, e.Reason, e.EntryMid, e.PeakMid, e.ExitMid, e.ChangePP, e.DrawdownPP, e.HeldFor.Round(time.Second))
}

// ExitConfig holds SPEC §2 exit thresholds.
type ExitConfig struct {
	ReversalTicks    int           // consecutive downticks from peak → reversal exit
	DrawdownPP       float64       // peak - current ≥ this pp → reversal exit
	StopLossPP       float64       // entry - current ≥ this pp → stop-loss exit
	MaxHold          time.Duration // force-close after this duration
	MinTickPPForMove float64       // ignore jitter smaller than this pp when counting down/up
}

func DefaultExitConfig() ExitConfig {
	return ExitConfig{
		ReversalTicks:    3,
		DrawdownPP:       2.0,
		StopLossPP:       3.0,
		MaxHold:          30 * time.Minute,
		MinTickPPForMove: 0.1, // 0.001 mid = 0.1pp; anything smaller is noise
	}
}

// trackedPosition is ExitTracker's per-asset runtime state.
// Public position bookkeeping lives in PositionManager (position.go).
type trackedPosition struct {
	AssetID   string
	Market    string
	EntryTime time.Time
	EntryMid  float64

	// runtime
	peakMid    float64
	lastMid    float64
	consecDown int
}

// ExitTracker holds open positions and emits exits when rules trip.
// Thread model: feed ticks in from one goroutine; callers pull ExitSignal from Signals().
type ExitTracker struct {
	cfg ExitConfig
	pos map[string]*trackedPosition
	out chan ExitSignal
}

func NewExitTracker(cfg ExitConfig) *ExitTracker {
	return &ExitTracker{
		cfg: cfg,
		pos: map[string]*trackedPosition{},
		out: make(chan ExitSignal, 64),
	}
}

// Signals exposes emitted exits.
func (e *ExitTracker) Signals() <-chan ExitSignal { return e.out }

// Open registers a new position at the given entry tick. No-op if already open.
func (e *ExitTracker) Open(assetID, market string, entry feed.Tick) {
	if _, exists := e.pos[assetID]; exists {
		return
	}
	e.pos[assetID] = &trackedPosition{
		AssetID:   assetID,
		Market:    market,
		EntryTime: entry.Time,
		EntryMid:  entry.Mid,
		peakMid:   entry.Mid,
		lastMid:   entry.Mid,
	}
}

// Has reports whether a position is currently open for assetID.
func (e *ExitTracker) Has(assetID string) bool {
	_, ok := e.pos[assetID]
	return ok
}

// OnTick feeds a 1s tick through the tracker. Emits an ExitSignal and closes the
// position if any exit rule fires. Returns the emitted signal (zero if none).
func (e *ExitTracker) OnTick(t feed.Tick) (ExitSignal, bool) {
	p, ok := e.pos[t.AssetID]
	if !ok {
		return ExitSignal{}, false
	}

	if t.Mid > p.peakMid {
		p.peakMid = t.Mid
	}

	// 1) timeout — checked first so we always exit stale positions even if price is frozen
	if t.Time.Sub(p.EntryTime) >= e.cfg.MaxHold {
		return e.close(p, t, ExitTimeout), true
	}

	// 2) stop-loss: entry - current ≥ StopLossPP
	stopMoveDown := (p.EntryMid - t.Mid) * 100
	if stopMoveDown >= e.cfg.StopLossPP {
		return e.close(p, t, ExitStopLoss), true
	}

	// 3) reversal by drawdown from peak
	drawdown := (p.peakMid - t.Mid) * 100
	if drawdown >= e.cfg.DrawdownPP {
		return e.close(p, t, ExitReversalDrawdown), true
	}

	// 4) reversal by N consecutive downticks (ignoring jitter)
	moveDownPP := (p.lastMid - t.Mid) * 100
	switch {
	case moveDownPP >= e.cfg.MinTickPPForMove:
		p.consecDown++
	case (t.Mid-p.lastMid)*100 >= e.cfg.MinTickPPForMove:
		p.consecDown = 0 // cleared by any meaningful uptick
	}
	p.lastMid = t.Mid
	if p.consecDown >= e.cfg.ReversalTicks {
		return e.close(p, t, ExitReversalTicks), true
	}

	return ExitSignal{}, false
}

func (e *ExitTracker) close(p *trackedPosition, t feed.Tick, reason ExitReason) ExitSignal {
	sig := ExitSignal{
		AssetID:    p.AssetID,
		Market:     p.Market,
		Time:       t.Time,
		EntryMid:   p.EntryMid,
		PeakMid:    p.peakMid,
		ExitMid:    t.Mid,
		HeldFor:    t.Time.Sub(p.EntryTime),
		ChangePP:   (t.Mid - p.EntryMid) * 100,
		DrawdownPP: (p.peakMid - t.Mid) * 100,
		Reason:     reason,
	}
	delete(e.pos, p.AssetID)
	select {
	case e.out <- sig:
	default:
	}
	return sig
}
