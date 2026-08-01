package strategy

import (
	"fmt"
	"math"
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
	SlippageBp    float64
	FlatFeeBp     float64
	TakerFeeRate  float64
}

func DefaultShadowExitConfig() ShadowExitConfig {
	return ShadowExitConfig{
		Timeouts:      []time.Duration{10 * time.Minute, 20 * time.Minute, 30 * time.Minute, 45 * time.Minute, 60 * time.Minute},
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
	ExitQuotePrice float64
	ExitPrice      float64
	GrossReturnPct float64
	GrossPnLUSD    float64
	EntryFeeUSD    float64
	ExitFeeUSD     float64
	NetPnLUSD      float64
	NetReturnPct   float64
	SlippageBp     float64
	TakerFeeRate   float64
	HeldFor        time.Duration
	ThresholdPct   float64
	HoldProfile    string
	EventStart     time.Time
	ExitDeadline   time.Time
	Question       string
	Outcome        string
	Source         string
	SignalSource   string
	ActualCloseAt  time.Time
	ActualReason   ExitReason
}

type shadowExitState struct {
	position   Position
	feeRate    float64
	fired      map[string]bool
	belowSince map[string]time.Time
	actualAt   time.Time
	actualWhy  ExitReason
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
	s.OpenWithFeeRate(p, s.cfg.TakerFeeRate)
}

// OpenWithFeeRate tracks a position using the same per-market fee rate as the
// paper order. Rehydrated positions use the configured fallback rate.
func (s *ShadowExitTracker) OpenWithFeeRate(p Position, feeRate float64) {
	if p.ID == "" || p.EntryMid <= 0 || p.EntryTime.IsZero() {
		return
	}
	if feeRate < 0 || feeRate > 1 {
		feeRate = s.cfg.TakerFeeRate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, ok := s.states[p.ID]; ok {
		state.position = p
		state.feeRate = feeRate
		return
	}
	s.states[p.ID] = &shadowExitState{
		position:   p,
		feeRate:    feeRate,
		fired:      map[string]bool{},
		belowSince: map[string]time.Time{},
	}
}

func (s *ShadowExitTracker) Close(posID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, posID)
}

// ActualClose keeps a counterfactual position alive until the longest timeout
// observation has fired, even if the running policy exits earlier.
func (s *ShadowExitTracker) ActualClose(p Position) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.states[p.ID]
	if !ok {
		return
	}
	state.actualAt = p.ExitTime
	state.actualWhy = p.ExitReason
	if s.allTimeoutsFired(state) {
		delete(s.states, p.ID)
	}
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
	if price <= 0 || price >= 1 {
		if !state.actualAt.IsZero() && tick.Time.Sub(state.position.EntryTime) >= s.shadowRetention() {
			delete(s.states, posID)
		}
		return nil
	}
	var out []ShadowExitObservation
	for _, timeout := range s.cfg.Timeouts {
		policy := timeoutPolicy(timeout)
		if state.fired[policy] || tick.Time.Sub(state.position.EntryTime) < timeout {
			continue
		}
		state.fired[policy] = true
		out = append(out, state.observation(tick.Time, price, policy, 0, s.cfg))
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
		out = append(out, state.observation(tick.Time, tick.BestBid, policy, stop*100, s.cfg))
	}
	for _, take := range s.cfg.TakeProfits {
		policy := percentPolicy("tp", take)
		if state.fired[policy] || tick.BestBid < state.position.EntryMid*(1+take) {
			continue
		}
		state.fired[policy] = true
		out = append(out, state.observation(tick.Time, tick.BestBid, policy, take*100, s.cfg))
	}
	if !state.actualAt.IsZero() && s.allTimeoutsFired(state) {
		delete(s.states, posID)
	}
	return out
}

func (s *ShadowExitTracker) allTimeoutsFired(state *shadowExitState) bool {
	for _, timeout := range s.cfg.Timeouts {
		if !state.fired[timeoutPolicy(timeout)] {
			return false
		}
	}
	return true
}

func (s *ShadowExitTracker) shadowRetention() time.Duration {
	if len(s.cfg.Timeouts) == 0 {
		return 0
	}
	return s.cfg.Timeouts[len(s.cfg.Timeouts)-1] + 30*time.Minute
}

func (s *shadowExitState) observation(observedAt time.Time, exitQuote float64, policy string, thresholdPct float64, cfg ShadowExitConfig) ShadowExitObservation {
	obs := shadowObservation(s.position, observedAt, exitQuote, policy, thresholdPct, cfg, s.feeRate)
	obs.ActualCloseAt = s.actualAt
	obs.ActualReason = s.actualWhy
	return obs
}

func shadowObservation(p Position, observedAt time.Time, exitQuote float64, policy string, thresholdPct float64, cfg ShadowExitConfig, feeRate float64) ShadowExitObservation {
	exitPrice := exitQuote * (1 - cfg.SlippageBp/10_000)
	units := p.Units
	if units <= 0 && p.SizeUSD > 0 {
		units = p.SizeUSD / p.EntryMid
	}
	entryFee := p.OpenFeeUSD
	if entryFee == 0 {
		entryFee = p.EntryFeeUSD + p.EntryFeeChargedUSD
	}
	grossPnL := (exitPrice - p.EntryMid) * units
	exitNotional := exitPrice * units
	flatFee := exitNotional * cfg.FlatFeeBp / 10_000
	platformFee := units * feeRate * exitPrice * (1 - exitPrice)
	if platformFee > 0 {
		platformFee = math.Round(platformFee*100_000) / 100_000
	}
	exitFee := flatFee + platformFee
	netPnL := grossPnL - entryFee - exitFee
	capital := p.SizeUSD
	if capital <= 0 {
		capital = p.EntryMid * units
	}
	netReturnPct := 0.0
	if capital > 0 {
		netReturnPct = netPnL / capital * 100
	}
	return ShadowExitObservation{
		PosID:          p.ID,
		AssetID:        p.AssetID,
		Market:         p.Market,
		Policy:         policy,
		ObservedAt:     observedAt,
		EntryTime:      p.EntryTime,
		EntryMid:       p.EntryMid,
		ExitQuotePrice: exitQuote,
		ExitPrice:      exitPrice,
		GrossReturnPct: (exitPrice/p.EntryMid - 1) * 100,
		GrossPnLUSD:    grossPnL,
		EntryFeeUSD:    entryFee,
		ExitFeeUSD:     exitFee,
		NetPnLUSD:      netPnL,
		NetReturnPct:   netReturnPct,
		SlippageBp:     cfg.SlippageBp,
		TakerFeeRate:   feeRate,
		HeldFor:        observedAt.Sub(p.EntryTime),
		ThresholdPct:   thresholdPct,
		HoldProfile:    p.HoldProfile,
		EventStart:     p.EventStart,
		ExitDeadline:   p.ExitDeadline,
		Question:       p.Question,
		Outcome:        p.Outcome,
		Source:         p.Source,
		SignalSource:   p.SignalSource,
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
