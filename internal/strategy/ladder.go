package strategy

import (
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// LadderConfig parametrizes the aggressive tranche-TP / SL / timeout machine
// introduced in Phase 7.b. Defaults come from SPEC §2.4 — pulled in to pay
// early profits and cap downside at a hard -10%.
type LadderConfig struct {
	TP1Pct  float64       // TP1 trigger: ExitMid ≥ Entry × (1 + TP1Pct)
	TP1Frac float64       // fraction of InitUnits to close on TP1
	TP2Pct  float64       // TP2 trigger: ExitMid ≥ Entry × (1 + TP2Pct)
	TP2Frac float64       // fraction of InitUnits to close on TP2 (typically 1.0 = clear remaining)
	SLPct   float64       // stop-loss: ExitMid ≤ Entry × (1 - SLPct) — closes 100% of remaining
	MaxHold time.Duration // force-close after held duration reaches this
}

// DefaultLadderConfig matches SPEC §2.4: +15% TP1 (close 50%), +30% TP2
// (clear remaining), -5% SL, 4h MaxHold. SL tightened from -10% to -5% on
// 2026-04-20 22:42 SGT after Phase 7.d sweep (SL=5% topped 10/10 configs).
func DefaultLadderConfig() LadderConfig {
	return LadderConfig{
		TP1Pct:  0.15,
		TP1Frac: 0.50,
		TP2Pct:  0.30,
		TP2Frac: 1.00,
		SLPct:   0.05,
		MaxHold: 4 * time.Hour,
	}
}

// Ladder exit reasons; kept separate from the auto-mode ExitReason constants
// so journal / logs can distinguish ladder tranches from legacy exits.
const (
	ExitLadderTP1     ExitReason = "ladder_tp1"
	ExitLadderTP2     ExitReason = "ladder_tp2"
	ExitLadderSL      ExitReason = "ladder_sl"
	ExitLadderTimeout ExitReason = "ladder_timeout"
)

// LadderExit is emitted when a tranche (or the final remainder) of a
// position should close. One tick can produce at most one LadderExit per
// position; stacked conditions (e.g. TP1 and TP2 on the same tick) resolve
// on subsequent ticks in order.
type LadderExit struct {
	PosID      string
	AssetID    string
	Market     string
	Time       time.Time
	EntryMid   float64
	ExitMid    float64
	CloseUnits float64
	Tranche    string // "t1" | "t2" | "sl" | "timeout"
	Final      bool   // true when the tranche closes the last remaining units
	Reason     ExitReason
	HeldFor    time.Duration
}

type ladderState struct {
	PosID     string
	AssetID   string
	Market    string
	EntryTime time.Time
	EntryMid  float64
	InitUnits float64
	RemUnits  float64
	TP1Done   bool
}

// LadderTracker owns per-position ladder state, keyed by posID. The caller
// (main.go) drives it from the same 1s tick-tail used by the legacy
// ExitTracker but indexes by position rather than asset — stacking allowed.
type LadderTracker struct {
	cfg    LadderConfig
	mu     sync.Mutex
	states map[string]*ladderState
}

func NewLadderTracker(cfg LadderConfig) *LadderTracker {
	return &LadderTracker{cfg: cfg, states: map[string]*ladderState{}}
}

// Open registers a new position. No-op if posID is already tracked.
func (l *LadderTracker) Open(posID, market, assetID string, entry feed.Tick, initUnits float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.states[posID]; exists {
		return
	}
	l.states[posID] = &ladderState{
		PosID:     posID,
		AssetID:   assetID,
		Market:    market,
		EntryTime: entry.Time,
		EntryMid:  entry.Mid,
		InitUnits: initUnits,
		RemUnits:  initUnits,
	}
}

// Has reports whether the tracker still owns posID.
func (l *LadderTracker) Has(posID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.states[posID]
	return ok
}

// Forget drops posID from the tracker without emitting. Use when the
// settlement watcher closes the remainder out-of-band.
func (l *LadderTracker) Forget(posID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.states, posID)
}

// OnTick feeds one tick for posID and returns a LadderExit if any rule fires.
// Priority order on a single tick: SL → timeout → TP2 (if TP1 done) → TP1.
// Gaps past TP2 without an intervening TP1 tick still emit TP1 first and
// defer TP2 to the next tick — keeps tranches disjoint and journal clean.
func (l *LadderTracker) OnTick(posID string, t feed.Tick) (LadderExit, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	st, ok := l.states[posID]
	if !ok {
		return LadderExit{}, false
	}

	heldFor := t.Time.Sub(st.EntryTime)
	mid := t.Mid

	slPx := st.EntryMid * (1 - l.cfg.SLPct)
	tp1Px := st.EntryMid * (1 + l.cfg.TP1Pct)
	tp2Px := st.EntryMid * (1 + l.cfg.TP2Pct)

	if mid <= slPx {
		return l.emitLocked(st, t, st.RemUnits, "sl", ExitLadderSL, heldFor), true
	}
	if heldFor >= l.cfg.MaxHold {
		return l.emitLocked(st, t, st.RemUnits, "timeout", ExitLadderTimeout, heldFor), true
	}
	if st.TP1Done && mid >= tp2Px {
		units := st.InitUnits * l.cfg.TP2Frac
		if units > st.RemUnits {
			units = st.RemUnits
		}
		return l.emitLocked(st, t, units, "t2", ExitLadderTP2, heldFor), true
	}
	if !st.TP1Done && mid >= tp1Px {
		units := st.InitUnits * l.cfg.TP1Frac
		if units > st.RemUnits {
			units = st.RemUnits
		}
		st.TP1Done = true
		return l.emitLocked(st, t, units, "t1", ExitLadderTP1, heldFor), true
	}
	return LadderExit{}, false
}

func (l *LadderTracker) emitLocked(st *ladderState, t feed.Tick, units float64, tranche string, reason ExitReason, heldFor time.Duration) LadderExit {
	if units > st.RemUnits {
		units = st.RemUnits
	}
	st.RemUnits -= units
	final := st.RemUnits <= 1e-9
	if final {
		delete(l.states, st.PosID)
	}
	return LadderExit{
		PosID:      st.PosID,
		AssetID:    st.AssetID,
		Market:     st.Market,
		Time:       t.Time,
		EntryMid:   st.EntryMid,
		ExitMid:    t.Mid,
		CloseUnits: units,
		Tranche:    tranche,
		Final:      final,
		Reason:     reason,
		HeldFor:    heldFor,
	}
}
