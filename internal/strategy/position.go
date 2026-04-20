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
	ErrAssetAlreadyOpen  = errors.New("asset already has open position")
	ErrMarketAlreadyOpen = errors.New("market already has open position (different side)")
	ErrMaxPositions      = errors.New("max concurrent positions reached")
	ErrMaxExposure       = errors.New("max total exposure reached")
	ErrInvalidEntry      = errors.New("invalid entry mid")
	ErrPositionNotFound  = errors.New("no open position for asset")
)

// PositionManager is the single source of truth for open/closed positions.
// Concurrent-safe; the detect loop opens on signals, the exit watcher closes on exit signals.
type PositionManager struct {
	cfg      PositionConfig
	mu       sync.Mutex
	open     map[string]*Position // by AssetID
	byMarket map[string]string    // market → assetID (dedupe across YES/NO)
	closed   []*Position
	nextID   int
}

func NewPositionManager(cfg PositionConfig) *PositionManager {
	return &PositionManager{
		cfg:      cfg,
		open:     make(map[string]*Position),
		byMarket: make(map[string]string),
	}
}

// Has returns true if assetID currently holds an open position.
func (pm *PositionManager) Has(assetID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	_, ok := pm.open[assetID]
	return ok
}

// HasMarket returns true if conditionID has ANY side (YES/NO) open.
func (pm *PositionManager) HasMarket(market string) bool {
	if market == "" {
		return false
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	_, ok := pm.byMarket[market]
	return ok
}

// Open books a paper position at entry.Mid using the default PerPositionUSD
// from config. See OpenSized for the variable-size variant used by the manual
// prompt flow (Phase 3.5).
func (pm *PositionManager) Open(assetID, market string, entry feed.Tick) (*Position, error) {
	return pm.OpenSized(assetID, market, entry, pm.cfg.PerPositionUSD)
}

// OpenSized books a paper position at an explicit size. Used by the manual
// button-select path where the boss picks 1 / 5 / 10 USDC per click. All the
// dedupe + exposure caps that Open enforces still apply.
func (pm *PositionManager) OpenSized(assetID, market string, entry feed.Tick, sizeUSD float64) (*Position, error) {
	if entry.Mid <= 0 || entry.Mid >= 1 {
		return nil, fmt.Errorf("%w: mid=%v", ErrInvalidEntry, entry.Mid)
	}
	if sizeUSD <= 0 {
		return nil, fmt.Errorf("%w: size=%v", ErrInvalidEntry, sizeUSD)
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.open[assetID]; ok {
		return nil, ErrAssetAlreadyOpen
	}
	if market != "" {
		if _, ok := pm.byMarket[market]; ok {
			return nil, ErrMarketAlreadyOpen
		}
	}
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
	pm.open[assetID] = p
	if market != "" {
		pm.byMarket[market] = assetID
	}
	return p, nil
}

// Close realizes PnL against the exit signal and moves the position to closed.
// PnL (paper): units × (exit - entry). Returns the closed copy.
func (pm *PositionManager) Close(assetID string, exit ExitSignal) (Position, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[assetID]
	if !ok {
		return Position{}, ErrPositionNotFound
	}
	p.ExitMid = exit.ExitMid
	p.ExitTime = exit.Time
	p.ExitReason = exit.Reason
	p.PnLUSD = p.Units * (exit.ExitMid - p.EntryMid)
	p.Status = PosClosed

	delete(pm.open, assetID)
	if p.Market != "" {
		delete(pm.byMarket, p.Market)
	}
	pm.closed = append(pm.closed, p)
	return *p, nil
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
