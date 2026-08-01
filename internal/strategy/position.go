package strategy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

	HoldProfileShort = "short"
	HoldProfileEvent = "event"
)

// Position represents a single paper (and eventually real) position.
// Paper mode: no orders hit the network; we just book entry/exit at the tick mid.
//
// Units shrinks as ladder tranches close; InitUnits is the original size at
// open time — keep it around so fee apportionment (entry fee × tranche share)
// stays consistent after a partial close.
type Position struct {
	ID                 string
	AssetID            string
	Market             string // Polymarket conditionID
	SizeUSD            float64
	Units              float64 // current remaining units
	InitUnits          float64 // units at open — invariant after partial closes
	OpenFeeUSD         float64 // fee paid on the buy; apportioned across tranches
	EntryFeeChargedUSD float64 // buy fee already assigned to earlier tranches
	EntryFeeUSD        float64 // apportioned buy fee on a closed position/tranche
	ExitFeeUSD         float64 // fee paid on the closing fill
	NetPnLUSD          float64 // gross PnL minus entry and exit fees
	Accounted          bool    `json:",omitempty"` // distinguishes exact zero net from legacy records
	EntryMid           float64
	EntryTime          time.Time
	ExitMid            float64
	ExitTime           time.Time
	ExitReason         ExitReason
	PnLUSD             float64
	Status             PositionStatus
	Question           string    `json:",omitempty"`
	Outcome            string    `json:",omitempty"`
	Source             string    `json:",omitempty"`
	WalletLabel        string    `json:",omitempty"`
	SignalSource       string    `json:",omitempty"`
	OpenOrderID        string    `json:",omitempty"`
	HoldProfile        string    `json:",omitempty"`
	EventStart         time.Time `json:",omitempty"`
	ExitDeadline       time.Time `json:",omitempty"`
}

// PositionConfig drives sizing + exposure caps. SPEC §2 / §6.
type PositionConfig struct {
	PerPositionUSD   float64 // Paper: 5 USDC
	MaxTotalOpenUSD  float64 // Cap on sum(open SizeUSD)
	MaxOpenPositions int     // Hard cap on concurrent positions
	MaxPerMarketUSD  float64 // Cap per conditionID (0 = unlimited)
}

func DefaultPositionConfig() PositionConfig {
	return PositionConfig{
		PerPositionUSD:   5.0,
		MaxTotalOpenUSD:  300.0,
		MaxOpenPositions: 60,
		MaxPerMarketUSD:  30.0,
	}
}

// PositionStats is a point-in-time summary of the book.
type PositionStats struct {
	Open           int
	Closed         int
	TotalExposure  float64 // sum of open SizeUSD
	RealizedPnLUSD float64
}

// ClosedAccounting is the fee-aware accounting persisted in the trade
// journal. It is used to migrate position snapshots written by older builds.
type ClosedAccounting struct {
	ID          string
	ExitTime    time.Time
	EntryFeeUSD float64
	ExitFeeUSD  float64
	NetPnLUSD   float64
}

var (
	ErrMaxPositions     = errors.New("max concurrent positions reached")
	ErrMaxExposure      = errors.New("max total exposure reached")
	ErrMaxPerMarket     = errors.New("max per-market exposure reached")
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
	saveMu   sync.Mutex
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
	if pm.cfg.MaxPerMarketUSD > 0 && market != "" {
		if pm.marketExposureLocked(market)+sizeUSD > pm.cfg.MaxPerMarketUSD+1e-9 {
			return nil, fmt.Errorf("%w: market=%s existing=$%.2f new=$%.2f cap=$%.2f",
				ErrMaxPerMarket, market[:min(len(market), 12)],
				pm.marketExposureLocked(market), sizeUSD, pm.cfg.MaxPerMarketUSD)
		}
	}

	pm.nextID++
	units := sizeUSD / entry.Mid
	p := &Position{
		ID:        fmt.Sprintf("p%d", pm.nextID),
		AssetID:   assetID,
		Market:    market,
		SizeUSD:   sizeUSD,
		Units:     units,
		InitUnits: units,
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

// SetOpenFee records the fee paid on the buy leg for later apportionment
// across ladder tranches. Safe to call once per position right after Open.
func (pm *PositionManager) SetOpenFee(posID string, feeUSD float64) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return ErrPositionNotFound
	}
	p.OpenFeeUSD = feeUSD
	return nil
}

// SetOpenAttribution persists the signal source and opening order ID with the
// position so attribution survives process restarts.
func (pm *PositionManager) SetOpenAttribution(posID, source, openOrderID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return ErrPositionNotFound
	}
	p.SignalSource = source
	p.OpenOrderID = openOrderID
	return nil
}

// ConfigureOpenHold persists the timeout policy for a filled position. A
// pre-event sports position is held until EventStart + eventPostStartHold;
// other positions retain the ordinary EntryTime + maxHold deadline.
func (pm *PositionManager) ConfigureOpenHold(posID string, eventStart time.Time, maxHold, eventPostStartHold time.Duration) (Position, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return Position{}, ErrPositionNotFound
	}
	p.HoldProfile, p.ExitDeadline = PlannedHold(p.EntryTime, eventStart, maxHold, eventPostStartHold)
	if p.HoldProfile == HoldProfileEvent {
		p.EventStart = eventStart
	} else {
		p.EventStart = time.Time{}
	}
	return *p, nil
}

// PlannedHold returns a durable timeout plan. Event positioning only applies
// when the event starts after entry and extends the ordinary short deadline.
func PlannedHold(entryTime, eventStart time.Time, maxHold, eventPostStartHold time.Duration) (string, time.Time) {
	if maxHold <= 0 || entryTime.IsZero() {
		return HoldProfileShort, time.Time{}
	}
	deadline := entryTime.Add(maxHold)
	if eventStart.After(entryTime) {
		eventDeadline := eventStart.Add(eventPostStartHold)
		if eventDeadline.After(deadline) {
			return HoldProfileEvent, eventDeadline
		}
	}
	return HoldProfileShort, deadline
}

// PositionTimeoutDue supports legacy snapshots without ExitDeadline while
// using the persisted deadline for newly opened event-aware positions.
func PositionTimeoutDue(p Position, now time.Time, maxHold time.Duration) bool {
	deadline := p.ExitDeadline
	if deadline.IsZero() && maxHold > 0 && !p.EntryTime.IsZero() {
		deadline = p.EntryTime.Add(maxHold)
	}
	return !deadline.IsZero() && !now.Before(deadline)
}

// ApplyOpenFill replaces the provisional signal price and size with the
// actual fill after the order succeeds.
func (pm *PositionManager) ApplyOpenFill(posID string, fillPrice, filledUnits float64, filledAt time.Time) error {
	if fillPrice <= 0 || fillPrice >= 1 || filledUnits <= 0 {
		return fmt.Errorf("%w: fill price=%v units=%v", ErrInvalidEntry, fillPrice, filledUnits)
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return ErrPositionNotFound
	}
	p.EntryMid = fillPrice
	p.Units = filledUnits
	p.InitUnits = filledUnits
	p.SizeUSD = fillPrice * filledUnits
	if !filledAt.IsZero() {
		p.EntryTime = filledAt
	}
	return nil
}

// CancelOpen releases a provisional position after order submission fails.
// It does not create a closed trade or realized PnL.
func (pm *PositionManager) CancelOpen(posID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return ErrPositionNotFound
	}
	pm.removeOpenLocked(p)
	return nil
}

// OldestByAsset returns the oldest open position for an asset.
func (pm *PositionManager) OldestByAsset(assetID string) (Position, error) {
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
	return *oldest, nil
}

// PartialClose closes closeUnits from an open position at the given exit
// signal and records the closed portion as its own tranche in closed[].
// If closeUnits covers (effectively) all remaining Units, the position is
// fully closed — equivalent to Close. Returns a copy of the tranche record.
func (pm *PositionManager) PartialClose(posID string, closeUnits float64, exit ExitSignal) (Position, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p, ok := pm.open[posID]
	if !ok {
		return Position{}, ErrPositionNotFound
	}
	if closeUnits <= 0 {
		return Position{}, fmt.Errorf("%w: close units=%v", ErrInvalidEntry, closeUnits)
	}
	if closeUnits >= p.Units-1e-9 {
		pm.closeLocked(p, exit)
		return *p, nil
	}
	tranche := &Position{
		ID:           p.ID,
		AssetID:      p.AssetID,
		Market:       p.Market,
		SizeUSD:      closeUnits * p.EntryMid,
		Units:        closeUnits,
		InitUnits:    p.InitUnits,
		OpenFeeUSD:   p.OpenFeeUSD,
		EntryFeeUSD:  apportionedOpenFee(p, closeUnits),
		ExitFeeUSD:   exit.ExitFeeUSD,
		EntryMid:     p.EntryMid,
		EntryTime:    p.EntryTime,
		ExitMid:      exit.ExitMid,
		ExitTime:     exit.Time,
		ExitReason:   exit.Reason,
		PnLUSD:       closeUnits * (exit.ExitMid - p.EntryMid),
		Accounted:    true,
		Status:       PosClosed,
		Question:     p.Question,
		Outcome:      p.Outcome,
		Source:       p.Source,
		WalletLabel:  p.WalletLabel,
		SignalSource: p.SignalSource,
		OpenOrderID:  p.OpenOrderID,
		HoldProfile:  p.HoldProfile,
		EventStart:   p.EventStart,
		ExitDeadline: p.ExitDeadline,
	}
	tranche.NetPnLUSD = tranche.PnLUSD - tranche.EntryFeeUSD - tranche.ExitFeeUSD
	p.EntryFeeChargedUSD += tranche.EntryFeeUSD
	p.Units -= closeUnits
	p.SizeUSD = p.Units * p.EntryMid
	pm.closed = append(pm.closed, tranche)
	return *tranche, nil
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
	p.EntryFeeUSD = remainingOpenFee(p)
	p.ExitFeeUSD = exit.ExitFeeUSD
	p.NetPnLUSD = p.PnLUSD - p.EntryFeeUSD - p.ExitFeeUSD
	p.Accounted = true
	p.Status = PosClosed

	pm.removeOpenLocked(p)
	pm.closed = append(pm.closed, p)
}

func apportionedOpenFee(p *Position, units float64) float64 {
	if p.InitUnits <= 0 || units <= 0 {
		return 0
	}
	return p.OpenFeeUSD * units / p.InitUnits
}

func remainingOpenFee(p *Position) float64 {
	remaining := p.OpenFeeUSD - p.EntryFeeChargedUSD
	if remaining < 0 {
		return 0
	}
	return remaining
}

// removeOpenLocked removes all open-position indexes. Caller must hold pm.mu.
func (pm *PositionManager) removeOpenLocked(p *Position) {
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
	closedCount := 0
	for _, p := range pm.closed {
		if p.ExitReason == "submit_failed" {
			continue
		}
		closedCount++
		if p.Accounted {
			realized += p.NetPnLUSD
		} else {
			realized += p.PnLUSD
		}
	}
	return PositionStats{
		Open:           len(pm.open),
		Closed:         closedCount,
		TotalExposure:  pm.totalExposureLocked(),
		RealizedPnLUSD: realized,
	}
}

// ReconcileClosedAccounting applies journal costs to legacy closed positions.
// Matching uses position ID plus exact exit time, which also distinguishes
// multiple partial-close tranches sharing one position ID.
func (pm *PositionManager) ReconcileClosedAccounting(records []ClosedAccounting) int {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	type key struct {
		id string
		t  int64
	}
	byKey := make(map[key]ClosedAccounting, len(records))
	for _, rec := range records {
		byKey[key{id: rec.ID, t: rec.ExitTime.UnixNano()}] = rec
	}
	updated := 0
	for _, p := range pm.closed {
		if p.Accounted {
			continue
		}
		rec, ok := byKey[key{id: p.ID, t: p.ExitTime.UnixNano()}]
		if !ok {
			continue
		}
		p.EntryFeeUSD = rec.EntryFeeUSD
		p.ExitFeeUSD = rec.ExitFeeUSD
		p.NetPnLUSD = rec.NetPnLUSD
		p.Accounted = true
		updated++
	}
	return updated
}

// ReconcileOpenEntryFees restores the entry-fee amount already assigned to
// partial closes before EntryFeeChargedUSD was added to position snapshots.
func (pm *PositionManager) ReconcileOpenEntryFees(records []ClosedAccounting) int {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	chargedByID := make(map[string]float64)
	for _, rec := range records {
		chargedByID[rec.ID] += rec.EntryFeeUSD
	}
	updated := 0
	for _, p := range pm.open {
		if p.EntryFeeChargedUSD > 0 {
			continue
		}
		charged := chargedByID[p.ID]
		if charged <= 0 {
			continue
		}
		if charged > p.OpenFeeUSD {
			charged = p.OpenFeeUSD
		}
		p.EntryFeeChargedUSD = charged
		updated++
	}
	return updated
}

func (pm *PositionManager) totalExposureLocked() float64 {
	var s float64
	for _, p := range pm.open {
		s += p.SizeUSD
	}
	return s
}

func (pm *PositionManager) marketExposureLocked(market string) float64 {
	var s float64
	for _, p := range pm.byMarket[market] {
		s += p.SizeUSD
	}
	return s
}

func (pm *PositionManager) ExposureForMarket(market string) float64 {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.marketExposureLocked(market)
}

// positionState is the JSON-serializable snapshot for persistence.
type positionState struct {
	NextID int         `json:"next_id"`
	Open   []*Position `json:"open"`
	Closed []*Position `json:"closed"`
}

// SaveState writes all open+closed positions to a JSON file so they survive
// daemon restarts. Called after every Open/Close.
func (pm *PositionManager) SaveState(path string) error {
	pm.saveMu.Lock()
	defer pm.saveMu.Unlock()

	pm.mu.Lock()
	open := make([]*Position, 0, len(pm.open))
	for _, p := range pm.open {
		cp := *p
		open = append(open, &cp)
	}
	closed := make([]*Position, 0, len(pm.closed))
	for _, p := range pm.closed {
		cp := *p
		closed = append(closed, &cp)
	}
	nextID := pm.nextID
	pm.mu.Unlock()

	st := positionState{NextID: nextID, Open: open, Closed: closed}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return writeStateAtomic(path, data)
}

func writeStateAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".positions-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// LoadState restores positions from a JSON file written by SaveState.
// Rebuilds the in-memory index maps. No-op if file doesn't exist.
func (pm *PositionManager) LoadState(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var st positionState
	if err := json.Unmarshal(data, &st); err != nil {
		return err
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.nextID = st.NextID
	pm.open = make(map[string]*Position)
	pm.byAsset = make(map[string]map[string]*Position)
	pm.byMarket = make(map[string]map[string]*Position)
	pm.closed = st.Closed
	if pm.closed == nil {
		pm.closed = []*Position{}
	}

	for _, p := range st.Open {
		pm.open[p.ID] = p
		if pm.byAsset[p.AssetID] == nil {
			pm.byAsset[p.AssetID] = map[string]*Position{}
		}
		pm.byAsset[p.AssetID][p.ID] = p
		if p.Market != "" {
			if pm.byMarket[p.Market] == nil {
				pm.byMarket[p.Market] = map[string]*Position{}
			}
			pm.byMarket[p.Market][p.ID] = p
		}
	}
	return nil
}
