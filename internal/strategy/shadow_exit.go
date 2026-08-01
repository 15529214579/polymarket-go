package strategy

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// ShadowExitConfig defines counterfactual exits that are observed but never
// sent to the order client. Price policies use the executable bid, not mid.
type ShadowExitConfig struct {
	Timeouts      []time.Duration
	StopLosses    []float64
	TakeProfits   []float64
	SLConfirmTime time.Duration
}

func DefaultShadowExitConfig() ShadowExitConfig {
	return ShadowExitConfig{
		Timeouts:      []time.Duration{10 * time.Minute, 20 * time.Minute, 30 * time.Minute},
		StopLosses:    []float64{0.20, 0.25},
		TakeProfits:   []float64{0.30, 0.50},
		SLConfirmTime: 15 * time.Second,
	}
}

// ShadowExitObservation is emitted once per position and policy. It is an
// analysis event only; callers must not treat it as an executable exit.
type ShadowExitObservation struct {
	PosID          string
	AssetID        string
	Market         string
	Policy         string
	ObservedAt     time.Time
	EntryTime      time.Time
	EntryMid       float64
	ExitPrice      float64
	GrossReturnPct float64
	HeldFor        time.Duration
	ThresholdPct   float64
	HoldProfile    string
	EventStart     time.Time
	ExitDeadline   time.Time
}

type shadowExitState struct {
	position   Position
	fired      map[string]bool
	belowSince map[string]time.Time
}

// ShadowExitTracker tracks counterfactual exits independently of the live
// ladder. It is safe to drive from the 1Hz recorder loop.
type ShadowExitTracker struct {
	cfg    ShadowExitConfig
	mu     sync.Mutex
	states map[string]*shadowExitState
}

func NewShadowExitTracker(cfg ShadowExitConfig) *ShadowExitTracker {
	cfg.Timeouts = append([]time.Duration(nil), cfg.Timeouts...)
	cfg.StopLosses = append([]float64(nil), cfg.StopLosses...)
	cfg.TakeProfits = append([]float64(nil), cfg.TakeProfits...)
	sort.Slice(cfg.Timeouts, func(i, j int) bool { return cfg.Timeouts[i] < cfg.Timeouts[j] })
	sort.Float64s(cfg.StopLosses)
	sort.Float64s(cfg.TakeProfits)
	return &ShadowExitTracker{cfg: cfg, states: map[string]*shadowExitState{}}
}

func (s *ShadowExitTracker) Open(p Position) {
	if p.ID == "" || p.EntryMid <= 0 || p.EntryTime.IsZero() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, ok := s.states[p.ID]; ok {
		state.position = p
		return
	}
	s.states[p.ID] = &shadowExitState{
		position:   p,
		fired:      map[string]bool{},
		belowSince: map[string]time.Time{},
	}
}

func (s *ShadowExitTracker) Close(posID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, posID)
}

func (s *ShadowExitTracker) Snapshot() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]string, len(s.states))
	for posID, state := range s.states {
		out[posID] = state.position.AssetID
	}
	return out
}

func (s *ShadowExitTracker) OnTick(posID string, tick feed.Tick) []ShadowExitObservation {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.states[posID]
	if !ok || tick.Time.Before(state.position.EntryTime) {
		return nil
	}

	price := tick.BestBid
	if price <= 0 {
		price = tick.Mid
	}
	var out []ShadowExitObservation
	for _, timeout := range s.cfg.Timeouts {
		policy := timeoutPolicy(timeout)
		if state.fired[policy] || tick.Time.Sub(state.position.EntryTime) < timeout {
			continue
		}
		state.fired[policy] = true
		out = append(out, shadowObservation(state.position, tick.Time, price, policy, 0))
	}

	// A price policy without a real bid would produce an optimistic and often
	// untradeable trigger, so mid-only samples are kept for time policies only.
	if tick.BestBid <= 0 {
		return out
	}
	for _, stop := range s.cfg.StopLosses {
		policy := percentPolicy("sl", stop)
		if state.fired[policy] {
			continue
		}
		if tick.BestBid > state.position.EntryMid*(1-stop) {
			delete(state.belowSince, policy)
			continue
		}
		since := state.belowSince[policy]
		if since.IsZero() {
			state.belowSince[policy] = tick.Time
			since = tick.Time
		}
		if tick.Time.Sub(since) < s.cfg.SLConfirmTime {
			continue
		}
		state.fired[policy] = true
		out = append(out, shadowObservation(state.position, tick.Time, tick.BestBid, policy, stop*100))
	}
	for _, take := range s.cfg.TakeProfits {
		policy := percentPolicy("tp", take)
		if state.fired[policy] || tick.BestBid < state.position.EntryMid*(1+take) {
			continue
		}
		state.fired[policy] = true
		out = append(out, shadowObservation(state.position, tick.Time, tick.BestBid, policy, take*100))
	}
	return out
}

func shadowObservation(p Position, observedAt time.Time, exitPrice float64, policy string, thresholdPct float64) ShadowExitObservation {
	return ShadowExitObservation{
		PosID:          p.ID,
		AssetID:        p.AssetID,
		Market:         p.Market,
		Policy:         policy,
		ObservedAt:     observedAt,
		EntryTime:      p.EntryTime,
		EntryMid:       p.EntryMid,
		ExitPrice:      exitPrice,
		GrossReturnPct: (exitPrice/p.EntryMid - 1) * 100,
		HeldFor:        observedAt.Sub(p.EntryTime),
		ThresholdPct:   thresholdPct,
		HoldProfile:    p.HoldProfile,
		EventStart:     p.EventStart,
		ExitDeadline:   p.ExitDeadline,
	}
}

func timeoutPolicy(timeout time.Duration) string {
	if timeout > 0 && timeout%time.Minute == 0 {
		return fmt.Sprintf("timeout_%dm", int(timeout/time.Minute))
	}
	return "timeout_" + timeout.String()
}

func percentPolicy(prefix string, value float64) string {
	return fmt.Sprintf("%s_%.0f", prefix, value*100)
}
