package strategy

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// PositionStatus is the lifecycle state of a paper/real position.
type PositionStatus string

const (
	PosOpen   PositionStatus = "open"
	PosClosed PositionStatus = "closed"
)

// Position represents a single paper (and eventually real) position.
// Paper mode: no orders hit the network; we just book entry/exit at the tick mid.
type Position struct {
	ID         string
	AssetID    string
	Market     string // Polymarket conditionID
	SizeUSD    float64
	Units      float64 // = SizeUSD / EntryMid
	EntryMid   float64
	EntryTime  time.Time
	ExitMid    float64
	ExitTime   time.Time
	ExitReason ExitReason
	PnLUSD     float64
	Status     PositionStatus
}

// PositionConfig drives sizing + exposure caps. SPEC §2 / §6.
type PositionConfig struct {
	PerPositionUSD   float64 // Paper: 5 USDC
	MaxTotalOpenUSD  float64 // Cap on sum(open SizeUSD)
	MaxOpenPositions int     // Hard cap on concurrent positions
}

func DefaultPositionConfig() PositionConfig {
	return PositionConfig{
		PerPositionUSD:   5.0,
		MaxTotalOpenUSD:  45.0, // 50% of 90 USDC starting bankroll (paper assumption)
		MaxOpenPositions: 6,
	}
}

// PositionStats is a point-in-time summary of the book.
type PositionStats struct {
	Open           int
	Closed         int
	TotalExposure  float64 // sum of open SizeUSD
	RealizedPnLUSD float64
}

var (
	ErrMaxPositions     = errors.New("max concurrent positions reached")
	ErrMaxExposure      = errors.New("max total exposure reached")
	ErrInvalidEntry     = errors.New("invalid entry mid")
	ErrPositionNotFound = errors.New("no open position for id/asset")
)

// PositionManager is the single source of truth for open/closed positions.
// Stacking is allowed: the same asset (and the same market) can hold multiple
// concurrent positions — dedupe is intentionally absent so the paper run can
// accumulate samples per market. Exposure and position-count caps still apply.
// Concurrent-safe.
type PositionManager struct {
	cfg      PositionConfig
	mu       sync.Mutex
	open     map[string]*Position            // by posID
	byAsset  map[string]map[string]*Position // assetID → posID set
	byMarket map[string]map[string]*Position // marketID → posID set
	closed   []*Position
	nextID   int
}

func NewPositionManager(cfg PositionConfig) *PositionManager {
	return &PositionManager{
		cfg:      cfg,
		open:     make(map[string]*Position),
		byAsset:  make(map[string]map[string]*Position),
		byMarket: make(map[string]map[string]*Position),
	}
}

// Has returns true if assetID currently holds at least one open position.
func (pm *PositionManager) Has(assetID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return len(pm.byAsset[assetID]) > 0
}

// HasMarket returns true if the market currently holds at least one open
// position on any side.
func (pm *PositionManager) HasMarket(market string) bool {
	if market == "" {
		return false
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return len(pm.byMarket[market]) > 0
}

// Open books a paper position at entry.Mid using the default PerPositionUSD
// from config. See OpenSized for the variable-size variant used by the manual
// prompt flow (Phase 3.5).
func (pm *PositionManager) Open(assetID, market string, entry feed.Tick) (*Position, error) {
	return pm.OpenSized(assetID, market, entry, pm.cfg.PerPositionUSD)
}

// OpenSized books a paper position at an explicit size. Stacking allowed —
// caller can open many positions on the same asset/market; only exposure and
// count caps fail. Used by both the auto signal loop and the manual
// button-select path (Phase 3.5, 1/5/10 USDC per click).
func (pm *PositionManager) OpenSized(assetID, market string, entry feed.Tick, sizeUSD float64) (*Position, error) {
	if entry.Mid <= 0 || entry.Mid >= 1 {
		return nil, fmt.Errorf("%w: mid=%v", ErrInvalidEntry, entry.Mid)
	}
	if sizeUSD <= 0 {
		return nil, fmt.Errorf("%w: size=%v", ErrInvalidEntry, sizeUSD)
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.open) >= pm.cfg.MaxOpenPositions {
		return nil, ErrMaxPositions
	}
	if pm.totalExposureLocked()+sizeUSD > pm.cfg.MaxTotalOpenUSD+1e-9 {
		return nil, ErrMaxExposure
	}

	pm.nextID++
	p := &Position{
		ID:        fmt.Sprintf("p%d", pm.nextID),
		AssetID:   assetID,
		Market:    market,
		SizeUSD:   sizeUSD,
		Units:     sizeUSD / entry.Mid,
		EntryMid:  entry.Mid,
		EntryTime: entry.Time,
		Status:    PosOpen,
	}
	pm.open[p.ID] = p
	if pm.byAsset[assetID] == nil {
		pm.byAsset[assetID] = map[string]*Position{}
	}
	pm.byAsset[assetID][p.ID] = p
	if market != "" {
		if pm.byMarket[market] == nil {
			pm.byMarket[market] = map[string]*Position{}
		}
		pm.byMarket[market][p.ID] = p
	}
	return p, nil
}

// Close realizes PnL against the exit signal for the given posID and moves
// the position to closed. PnL (paper): units × (exit - entry). Returns the
// closed copy.
func (pm *PositionManager) Close(posID string, exit ExitSignal) (Position, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return Position{}, ErrPositionNotFound
	}
	pm.closeLocked(p, exit)
	return *p, nil
}

// CloseFirstByAsset closes the oldest open position for assetID. Kept as a
// convenience for the exit-watch path which signals by asset, not by posID.
// When multiple positions are stacked the remainder stay open — subsequent
// exit events (or settlement) will close them.
func (pm *PositionManager) CloseFirstByAsset(assetID string, exit ExitSignal) (Position, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	set := pm.byAsset[assetID]
	if len(set) == 0 {
		return Position{}, ErrPositionNotFound
	}
	var oldest *Position
	for _, p := range set {
		if oldest == nil || p.EntryTime.Before(oldest.EntryTime) {
			oldest = p
		}
	}
	pm.closeLocked(oldest, exit)
	return *oldest, nil
}

// closeLocked mutates state; caller must hold pm.mu.
func (pm *PositionManager) closeLocked(p *Position, exit ExitSignal) {
	p.ExitMid = exit.ExitMid
	p.ExitTime = exit.Time
	p.ExitReason = exit.Reason
	p.PnLUSD = p.Units * (exit.ExitMid - p.EntryMid)
	p.Status = PosClosed

	delete(pm.open, p.ID)
	if set := pm.byAsset[p.AssetID]; set != nil {
		delete(set, p.ID)
		if len(set) == 0 {
			delete(pm.byAsset, p.AssetID)
		}
	}
	if p.Market != "" {
		if set := pm.byMarket[p.Market]; set != nil {
			delete(set, p.ID)
			if len(set) == 0 {
				delete(pm.byMarket, p.Market)
			}
		}
	}
	pm.closed = append(pm.closed, p)
}

// Snapshot returns a copy of all open positions sorted by entry time.
func (pm *PositionManager) Snapshot() []Position {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	out := make([]Position, 0, len(pm.open))
	for _, p := range pm.open {
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EntryTime.Before(out[j].EntryTime) })
	return out
}

// Closed returns a copy of all closed positions in close-time order.
func (pm *PositionManager) Closed() []Position {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	out := make([]Position, len(pm.closed))
	for i, p := range pm.closed {
		out[i] = *p
	}
	return out
}

// Stats returns a point-in-time summary.
func (pm *PositionManager) Stats() PositionStats {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	var realized float64
	for _, p := range pm.closed {
		realized += p.PnLUSD
	}
	return PositionStats{
		Open:           len(pm.open),
		Closed:         len(pm.closed),
		TotalExposure:  pm.totalExposureLocked(),
		RealizedPnLUSD: realized,
	}
}

func (pm *PositionManager) totalExposureLocked() float64 {
	var s float64
	for _, p := range pm.open {
		s += p.SizeUSD
	}
	return s
}
