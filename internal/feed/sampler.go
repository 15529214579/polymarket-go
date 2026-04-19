package feed

import (
	"context"
	"sync"
	"time"
)

// Tick is a 1-second summary for one asset.
// Emitted at the close of each second; counters cover that second only.
type Tick struct {
	AssetID  string
	Market   string
	Time     time.Time
	BestBid  float64
	BestAsk  float64
	Mid      float64
	Trades   int
	BuyVol   float64
	SellVol  float64
}

// WindowStats summarizes the last N 1-second ticks for one asset.
// Delta/ratios are only meaningful when Samples ≈ WindowSec (i.e. warm).
type WindowStats struct {
	AssetID   string
	Market    string
	WindowSec int
	Samples   int
	StartMid  float64
	EndMid    float64
	DeltaPP   float64 // (end-start) in percentage points (0..100 scale)
	Upticks   int
	Downticks int
	Flats     int
	BuyVol    float64
	SellVol   float64
	BuyRatio  float64 // BuyVol / (BuyVol+SellVol); 0 when no trades
}

// Sampler turns the async book/trade streams into a uniform 1-Hz tick stream
// and keeps a per-asset rolling window for strategy queries.
type Sampler struct {
	windowSec int
	out       chan Tick

	mu    sync.Mutex
	state map[string]*assetState
}

type assetState struct {
	market  string
	bestBid float64
	bestAsk float64
	lastMid float64
	seen    bool

	// accumulators for the current (in-progress) second
	buyVol  float64
	sellVol float64
	trades  int

	// ring of recent Ticks (len == windowSec once warm)
	ring []Tick
	head int // next write index
	full bool
}

func NewSampler(windowSec int) *Sampler {
	if windowSec <= 0 {
		windowSec = 60
	}
	return &Sampler{
		windowSec: windowSec,
		out:       make(chan Tick, 1024),
		state:     map[string]*assetState{},
	}
}

// Ticks streams per-asset 1-second summaries. Buffered; drop-on-full is
// avoided by sizing the chan generously — consumers should not block >1s.
func (s *Sampler) Ticks() <-chan Tick { return s.out }

// Run consumes books+trades and emits a Tick per asset per second until ctx is done.
func (s *Sampler) Run(ctx context.Context, books <-chan BookEvent, trades <-chan TradeEvent) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-books:
			if !ok {
				return nil
			}
			s.onBook(ev)
		case tr, ok := <-trades:
			if !ok {
				return nil
			}
			s.onTrade(tr)
		case now := <-ticker.C:
			s.flushSecond(now)
		}
	}
}

func (s *Sampler) onBook(ev BookEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.ensure(ev.AssetID, ev.Market)
	if len(ev.Bids) > 0 {
		st.bestBid = ev.Bids[0].Price
	} else {
		st.bestBid = 0
	}
	if len(ev.Asks) > 0 {
		st.bestAsk = ev.Asks[0].Price
	} else {
		st.bestAsk = 0
	}
	st.seen = true
}

func (s *Sampler) onTrade(tr TradeEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.ensure(tr.AssetID, tr.Market)
	st.trades++
	switch tr.Side {
	case "BUY":
		st.buyVol += tr.Size
	case "SELL":
		st.sellVol += tr.Size
	}
}

func (s *Sampler) ensure(assetID, market string) *assetState {
	st, ok := s.state[assetID]
	if !ok {
		st = &assetState{market: market, ring: make([]Tick, s.windowSec)}
		s.state[assetID] = st
	}
	if market != "" && st.market == "" {
		st.market = market
	}
	return st
}

func (s *Sampler) flushSecond(now time.Time) {
	s.mu.Lock()
	// Emit outside the lock.
	batch := make([]Tick, 0, len(s.state))
	for assetID, st := range s.state {
		if !st.seen {
			continue
		}
		mid := midOf(st.bestBid, st.bestAsk, st.lastMid)
		t := Tick{
			AssetID: assetID,
			Market:  st.market,
			Time:    now,
			BestBid: st.bestBid,
			BestAsk: st.bestAsk,
			Mid:     mid,
			Trades:  st.trades,
			BuyVol:  st.buyVol,
			SellVol: st.sellVol,
		}
		st.ring[st.head] = t
		st.head = (st.head + 1) % s.windowSec
		if st.head == 0 {
			st.full = true
		}
		st.lastMid = mid
		st.trades = 0
		st.buyVol = 0
		st.sellVol = 0
		batch = append(batch, t)
	}
	s.mu.Unlock()

	for _, t := range batch {
		select {
		case s.out <- t:
		default:
			// consumer stalled; drop oldest by draining one slot
			select {
			case <-s.out:
			default:
			}
			select {
			case s.out <- t:
			default:
			}
		}
	}
}

// Window returns stats over the last windowSec ticks for assetID.
// ok=false if the asset is unknown or has zero samples yet.
func (s *Sampler) Window(assetID string) (WindowStats, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.state[assetID]
	if !ok {
		return WindowStats{}, false
	}
	ws := WindowStats{
		AssetID:   assetID,
		Market:    st.market,
		WindowSec: s.windowSec,
	}
	ticks := st.orderedTicks(s.windowSec)
	ws.Samples = len(ticks)
	if ws.Samples == 0 {
		return ws, false
	}
	ws.StartMid = ticks[0].Mid
	ws.EndMid = ticks[len(ticks)-1].Mid
	ws.DeltaPP = (ws.EndMid - ws.StartMid) * 100

	prev := ws.StartMid
	for _, t := range ticks {
		switch {
		case t.Mid > prev:
			ws.Upticks++
		case t.Mid < prev:
			ws.Downticks++
		default:
			ws.Flats++
		}
		prev = t.Mid
		ws.BuyVol += t.BuyVol
		ws.SellVol += t.SellVol
	}
	if tot := ws.BuyVol + ws.SellVol; tot > 0 {
		ws.BuyRatio = ws.BuyVol / tot
	}
	return ws, true
}

// TickTail returns up to n most recent ticks for assetID in chronological order.
// ok=false if the asset is unknown.
func (s *Sampler) TickTail(assetID string, n int) ([]Tick, bool) {
	if n <= 0 {
		return nil, true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.state[assetID]
	if !ok {
		return nil, false
	}
	all := st.orderedTicks(s.windowSec)
	if len(all) <= n {
		return all, true
	}
	return all[len(all)-n:], true
}

// Snapshot returns a window summary for every asset that has at least one tick.
func (s *Sampler) Snapshot() []WindowStats {
	s.mu.Lock()
	ids := make([]string, 0, len(s.state))
	for id := range s.state {
		ids = append(ids, id)
	}
	s.mu.Unlock()
	out := make([]WindowStats, 0, len(ids))
	for _, id := range ids {
		if w, ok := s.Window(id); ok {
			out = append(out, w)
		}
	}
	return out
}

func (st *assetState) orderedTicks(windowSec int) []Tick {
	if !st.full {
		// ring filled 0..head-1 in order
		if st.head == 0 {
			return nil
		}
		out := make([]Tick, st.head)
		copy(out, st.ring[:st.head])
		return out
	}
	out := make([]Tick, windowSec)
	copy(out, st.ring[st.head:])
	copy(out[windowSec-st.head:], st.ring[:st.head])
	return out
}

func midOf(bid, ask, fallback float64) float64 {
	switch {
	case bid > 0 && ask > 0:
		return (bid + ask) / 2
	case bid > 0:
		return bid
	case ask > 0:
		return ask
	default:
		return fallback
	}
}
