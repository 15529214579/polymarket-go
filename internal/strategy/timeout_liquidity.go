package strategy

import (
	"sync"
	"time"
)

// TimeoutLiquidityState tracks an expired paper position that cannot be sold
// because the order book has no executable bid.
type TimeoutLiquidityState struct {
	PositionID            string
	Market                string
	FirstSeen             time.Time
	LastSeen              time.Time
	LastLogged            time.Time
	Attempts              int
	ExposureUSD           float64
	ConservativeNetPnLUSD float64
}

type TimeoutLiquiditySummary struct {
	Positions             int
	ExposureUSD           float64
	ConservativeNetPnLUSD float64
}

type TimeoutLiquidityTracker struct {
	mu          sync.Mutex
	logInterval time.Duration
	states      map[string]TimeoutLiquidityState
}

func NewTimeoutLiquidityTracker(logInterval time.Duration) *TimeoutLiquidityTracker {
	if logInterval <= 0 {
		logInterval = 5 * time.Minute
	}
	return &TimeoutLiquidityTracker{
		logInterval: logInterval,
		states:      make(map[string]TimeoutLiquidityState),
	}
}

// Observe records one failed timeout-exit attempt. The boolean is true for
// the first observation and then once per log interval.
func (t *TimeoutLiquidityTracker) Observe(p Position, now time.Time) (TimeoutLiquidityState, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	state, exists := t.states[p.ID]
	if !exists {
		state = TimeoutLiquidityState{
			PositionID: p.ID,
			Market:     p.Market,
			FirstSeen:  now,
		}
	}
	state.LastSeen = now
	state.Attempts++
	state.ExposureUSD = p.SizeUSD
	remainingEntryFee := p.OpenFeeUSD - p.EntryFeeChargedUSD
	if remainingEntryFee < 0 {
		remainingEntryFee = 0
	}
	state.ConservativeNetPnLUSD = -(p.SizeUSD + remainingEntryFee)
	shouldLog := !exists || state.LastLogged.IsZero() || now.Sub(state.LastLogged) >= t.logInterval
	if shouldLog {
		state.LastLogged = now
	}
	t.states[p.ID] = state
	return state, shouldLog
}

func (t *TimeoutLiquidityTracker) Resolve(positionID string) (TimeoutLiquidityState, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	state, ok := t.states[positionID]
	if ok {
		delete(t.states, positionID)
	}
	return state, ok
}

func (t *TimeoutLiquidityTracker) Summary() TimeoutLiquiditySummary {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out TimeoutLiquiditySummary
	for _, state := range t.states {
		out.Positions++
		out.ExposureUSD += state.ExposureUSD
		out.ConservativeNetPnLUSD += state.ConservativeNetPnLUSD
	}
	return out
}
