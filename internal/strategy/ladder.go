package strategy

import (
	"math"
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// LadderConfig parametrizes tranche TP, trailing, confirmed SL, and timeout.
type LadderConfig struct {
	TP1Pct               float64       // first take-profit return threshold
	TP1Frac              float64       // fraction of InitUnits to close on TP1
	TP2Pct               float64       // second take-profit return threshold
	TP2Frac              float64       // fraction of InitUnits to close on TP2
	SLPct                float64       // stop-loss return threshold
	TrailingPct          float64       // post-TP1 drawdown from peak return; 0 disables
	SLConfirmTime        time.Duration // loss must persist this long before closing
	MinHoldBeforeSL      time.Duration // ignore transient loss immediately after entry
	MaxHold              time.Duration // force-close after held duration reaches this
	FeeAware             bool          // evaluate thresholds after entry/exit costs
	RequireExecutableBid bool          // do not evaluate exits without a valid best bid
	SlippageBp           float64
	FlatFeeBp            float64
	TakerFeeRate         float64
}

// DefaultLadderConfig returns defaults with TP disabled (999% = ride to settlement/timeout),
// SL 15%, 4h MaxHold. Updated 2026-04-24 21:38 SGT after tickpath sweep
// (42 paths, 25K ticks): SL 5%→15% cuts false stops from 26→14; TP unlimited
// outperforms 15/30% by +$30 over sample; entry cap 0.70→0.50 cuts losing band.
func DefaultLadderConfig() LadderConfig {
	return LadderConfig{
		TP1Pct:  9.99,
		TP1Frac: 0.50,
		TP2Pct:  9.99,
		TP2Frac: 1.00,
		SLPct:   0.15,
		MaxHold: 4 * time.Hour,
	}
}

// Ladder exit reasons; kept separate from the auto-mode ExitReason constants
// so journal / logs can distinguish ladder tranches from legacy exits.
const (
	ExitLadderTP1      ExitReason = "ladder_tp1"
	ExitLadderTP2      ExitReason = "ladder_tp2"
	ExitLadderSL       ExitReason = "ladder_sl"
	ExitLadderTrailing ExitReason = "ladder_trailing"
	ExitLadderTimeout  ExitReason = "ladder_timeout"
)

// LadderExit is emitted when a tranche (or the final remainder) of a
// position should close. One tick can produce at most one LadderExit per
// position; stacked conditions (e.g. TP1 and TP2 on the same tick) resolve
// on subsequent ticks in order.
type LadderExit struct {
	PosID        string
	AssetID      string
	Market       string
	Time         time.Time
	EntryMid     float64
	ExitMid      float64
	CloseUnits   float64
	Tranche      string // "t1" | "t2" | "trail" | "sl" | "timeout"
	Final        bool   // true when the tranche closes the last remaining units
	Reason       ExitReason
	HeldFor      time.Duration
	TakerFeeRate float64
}

type ladderState struct {
	PosID              string
	AssetID            string
	Market             string
	EntryTime          time.Time
	Deadline           time.Time
	EntryMid           float64
	InitUnits          float64
	RemUnits           float64
	TP1Done            bool
	TP1Closed          float64
	OpenFeeUSD         float64
	EntryFeeChargedUSD float64
	FeeRate            float64
	PeakReturn         float64
	PeakSet            bool
	BelowSince         time.Time
	Pending            *LadderExit
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
	l.OpenWithDeadline(posID, market, assetID, entry, initUnits, time.Time{})
}

// OpenWithDeadline registers a position with a persisted timeout deadline.
// A zero deadline preserves the legacy EntryTime + MaxHold behavior.
func (l *LadderTracker) OpenWithDeadline(posID, market, assetID string, entry feed.Tick, initUnits float64, deadline time.Time) {
	l.openPosition(Position{
		ID: posID, Market: market, AssetID: assetID, EntryTime: entry.Time,
		EntryMid: entry.Mid, Units: initUnits, InitUnits: initUnits, ExitDeadline: deadline,
	}, l.cfg.TakerFeeRate)
}

// OpenPosition registers the durable position accounting used by fee-aware
// exits. It also restores partial TP progress after a process restart.
func (l *LadderTracker) OpenPosition(p Position, feeRate float64) {
	l.openPosition(p, feeRate)
}

// SyncPosition refreshes remaining units and fee accounting after an
// out-of-band partial close, using the configured fallback market fee rate.
func (l *LadderTracker) SyncPosition(p Position) {
	l.openPosition(p, l.cfg.TakerFeeRate)
}

func (l *LadderTracker) openPosition(p Position, feeRate float64) {
	if p.ID == "" || p.EntryMid <= 0 || p.EntryMid >= 1 || p.Units <= 0 || p.EntryTime.IsZero() {
		return
	}
	if feeRate < 0 || feeRate > 1 || math.IsNaN(feeRate) || math.IsInf(feeRate, 0) {
		feeRate = l.cfg.TakerFeeRate
	}
	initUnits := p.InitUnits
	if initUnits <= 0 {
		initUnits = p.Units
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if existing, exists := l.states[p.ID]; exists {
		existing.Deadline = p.ExitDeadline
		existing.RemUnits = p.Units
		existing.OpenFeeUSD = p.OpenFeeUSD
		existing.EntryFeeChargedUSD = p.EntryFeeChargedUSD
		existing.FeeRate = feeRate
		closedUnits := existing.InitUnits - p.Units
		if closedUnits > existing.TP1Closed {
			existing.TP1Closed = closedUnits
		}
		existing.TP1Done = existing.TP1Closed+1e-9 >= existing.InitUnits*l.cfg.TP1Frac
		existing.Pending = nil
		existing.BelowSince = time.Time{}
		return
	}
	closedUnits := initUnits - p.Units
	if closedUnits < 0 {
		closedUnits = 0
	}
	l.states[p.ID] = &ladderState{
		PosID:              p.ID,
		AssetID:            p.AssetID,
		Market:             p.Market,
		EntryTime:          p.EntryTime,
		Deadline:           p.ExitDeadline,
		EntryMid:           p.EntryMid,
		InitUnits:          initUnits,
		RemUnits:           p.Units,
		TP1Closed:          closedUnits,
		TP1Done:            closedUnits+1e-9 >= initUnits*l.cfg.TP1Frac,
		OpenFeeUSD:         p.OpenFeeUSD,
		EntryFeeChargedUSD: p.EntryFeeChargedUSD,
		FeeRate:            feeRate,
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
// Priority order on a single tick: confirmed SL, timeout, TP2, trailing, TP1.
// Gaps past TP2 without an intervening TP1 tick still emit TP1 first and
// defer TP2 to the next tick — keeps tranches disjoint and journal clean.
func (l *LadderTracker) OnTick(posID string, t feed.Tick) (LadderExit, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	st, ok := l.states[posID]
	if !ok || st.Pending != nil {
		return LadderExit{}, false
	}

	heldFor := t.Time.Sub(st.EntryTime)
	exitPrice, ok := l.executableExitPrice(t)
	if !ok {
		return LadderExit{}, false
	}
	currentReturn := l.returnAfterCostsLocked(st, exitPrice, st.RemUnits)

	slHit := exitPrice <= st.EntryMid*(1-l.cfg.SLPct)
	if l.cfg.FeeAware {
		slHit = currentReturn <= -l.cfg.SLPct
	}
	if heldFor >= l.cfg.MinHoldBeforeSL && l.cfg.SLPct > 0 && slHit {
		if st.BelowSince.IsZero() {
			st.BelowSince = t.Time
		}
		if t.Time.Sub(st.BelowSince) >= l.cfg.SLConfirmTime {
			return l.proposeLocked(st, t, exitPrice, st.RemUnits, "sl", ExitLadderSL, heldFor), true
		}
	} else {
		st.BelowSince = time.Time{}
	}
	timeoutDue := !st.Deadline.IsZero() && !t.Time.Before(st.Deadline)
	if st.Deadline.IsZero() {
		timeoutDue = heldFor >= l.cfg.MaxHold
	}
	if timeoutDue {
		return l.proposeLocked(st, t, exitPrice, st.RemUnits, "timeout", ExitLadderTimeout, heldFor), true
	}
	tp2Hit := exitPrice >= st.EntryMid*(1+l.cfg.TP2Pct)
	if l.cfg.FeeAware {
		tp2Hit = currentReturn >= l.cfg.TP2Pct
	}
	if st.TP1Done && tp2Hit {
		units := st.InitUnits * l.cfg.TP2Frac
		if units > st.RemUnits {
			units = st.RemUnits
		}
		return l.proposeLocked(st, t, exitPrice, units, "t2", ExitLadderTP2, heldFor), true
	}
	if st.TP1Done && l.cfg.TrailingPct > 0 {
		if !st.PeakSet || currentReturn > st.PeakReturn {
			st.PeakReturn = currentReturn
			st.PeakSet = true
		} else if st.PeakReturn-currentReturn >= l.cfg.TrailingPct {
			return l.proposeLocked(st, t, exitPrice, st.RemUnits, "trail", ExitLadderTrailing, heldFor), true
		}
	}
	tp1Hit := exitPrice >= st.EntryMid*(1+l.cfg.TP1Pct)
	if l.cfg.FeeAware {
		tp1Hit = currentReturn >= l.cfg.TP1Pct
	}
	if !st.TP1Done && tp1Hit {
		target := st.InitUnits * l.cfg.TP1Frac
		units := target - st.TP1Closed
		if units <= 1e-9 {
			st.TP1Done = true
			return LadderExit{}, false
		}
		if units > st.RemUnits {
			units = st.RemUnits
		}
		return l.proposeLocked(st, t, exitPrice, units, "t1", ExitLadderTP1, heldFor), true
	}
	return LadderExit{}, false
}

func (l *LadderTracker) executableExitPrice(t feed.Tick) (float64, bool) {
	quote := t.Mid
	if l.cfg.RequireExecutableBid || l.cfg.FeeAware {
		quote = t.BestBid
	}
	if l.cfg.RequireExecutableBid {
		quoteAge := t.Time.Sub(t.QuoteTime)
		if t.BestBidSize <= 0 || t.QuoteTime.IsZero() || quoteAge > 5*time.Second {
			return 0, false
		}
	}
	if quote <= 0 || quote >= 1 || math.IsNaN(quote) || math.IsInf(quote, 0) {
		return 0, false
	}
	price := quote * (1 - l.cfg.SlippageBp/10_000)
	if price <= 0 || price >= 1 {
		return 0, false
	}
	return price, true
}

func (l *LadderTracker) returnAfterCostsLocked(st *ladderState, exitPrice, units float64) float64 {
	if units <= 0 || st.EntryMid <= 0 {
		return 0
	}
	gross := units * (exitPrice - st.EntryMid)
	entryFee := 0.0
	remainingEntryFee := st.OpenFeeUSD - st.EntryFeeChargedUSD
	if remainingEntryFee < 0 {
		remainingEntryFee = 0
	}
	if st.RemUnits > 0 {
		entryFee = remainingEntryFee * units / st.RemUnits
	}
	exitNotional := units * exitPrice
	exitFee := exitNotional*l.cfg.FlatFeeBp/10_000 + units*st.FeeRate*exitPrice*(1-exitPrice)
	if exitFee > 0 {
		exitFee = math.Round(exitFee*100_000) / 100_000
	}
	return (gross - entryFee - exitFee) / (units * st.EntryMid)
}

func (l *LadderTracker) proposeLocked(st *ladderState, t feed.Tick, exitPrice, units float64, tranche string, reason ExitReason, heldFor time.Duration) LadderExit {
	if units > st.RemUnits {
		units = st.RemUnits
	}
	exit := LadderExit{
		PosID:        st.PosID,
		AssetID:      st.AssetID,
		Market:       st.Market,
		Time:         t.Time,
		EntryMid:     st.EntryMid,
		ExitMid:      exitPrice,
		CloseUnits:   units,
		Tranche:      tranche,
		Final:        st.RemUnits-units <= 1e-9,
		Reason:       reason,
		HeldFor:      heldFor,
		TakerFeeRate: st.FeeRate,
	}
	st.Pending = &exit
	return exit
}

// Confirm applies the pending tranche only after the closing fill is durable.
func (l *LadderTracker) Confirm(posID string, filledUnits ...float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	st, ok := l.states[posID]
	if !ok || st.Pending == nil {
		return
	}
	pending := st.Pending
	units := pending.CloseUnits
	if len(filledUnits) > 0 {
		units = filledUnits[0]
	}
	if units <= 0 {
		return
	}
	if units > pending.CloseUnits {
		units = pending.CloseUnits
	}
	if units > st.RemUnits {
		units = st.RemUnits
	}
	beforeUnits := st.RemUnits
	remainingEntryFee := st.OpenFeeUSD - st.EntryFeeChargedUSD
	if remainingEntryFee > 0 && beforeUnits > 0 {
		st.EntryFeeChargedUSD += remainingEntryFee * units / beforeUnits
	}
	st.RemUnits -= units
	if pending.Tranche == "t1" {
		st.TP1Closed += units
		st.TP1Done = st.TP1Closed+1e-9 >= st.InitUnits*l.cfg.TP1Frac
		if st.TP1Done && st.RemUnits > 1e-9 {
			st.PeakReturn = l.returnAfterCostsLocked(st, pending.ExitMid, st.RemUnits)
			st.PeakSet = true
		}
	}
	if st.RemUnits <= 1e-9 {
		delete(l.states, posID)
		return
	}
	st.Pending = nil
}

// Retry releases a pending tranche after a failed order without consuming state.
func (l *LadderTracker) Retry(posID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if st, ok := l.states[posID]; ok {
		st.Pending = nil
	}
}
