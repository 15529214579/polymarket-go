// Package risk implements the runtime risk controls described in SPEC §6:
//
//   - Daily realized PnL circuit breaker (pct-of-bankroll cap).
//   - Per-trade realized-loss ceiling (monitor only, trade is already exited).
//   - WSS feed-silence watchdog (>N sec without a book/trade event).
//
// The RiskManager is the single gate the detect loop consults before opening a
// new paper/real position. Once tripped, it stays tripped until Resume() is
// called (manual — per SPEC §6 "等老板手动恢复"). Day boundaries are SGT
// (Asia/Singapore) per the rest of the project's observability surfaces.
package risk

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// BlockReason enumerates why the manager currently refuses new opens.
type BlockReason string

const (
	BlockNone        BlockReason = ""
	BlockDailyLoss   BlockReason = "daily_loss"
	BlockDrawdown    BlockReason = "drawdown"
	BlockFeedSilence BlockReason = "feed_silence"
	BlockManualPause BlockReason = "manual_pause"
)

// Config is the static tuning knob set — see SPEC §6.
type Config struct {
	// StartingBankrollUSD anchors pct-based caps. Day 0 snapshot = 90.41 USDC.
	StartingBankrollUSD float64
	// DailyLossPct, e.g. 0.15 = 15%. Applied to StartingBankrollUSD.
	DailyLossPct float64
	// MaxSingleLossUSD: per-trade realized loss ceiling. Observational only
	// (the trade already closed by the time we see it).
	MaxSingleLossUSD float64
	// MaxDrawdownPct: portfolio-level max drawdown as fraction of StartingBankrollUSD.
	// Unlike DailyLossPct which resets each day, drawdown tracks cumulative PnL
	// from inception. Trips when (peakEquity - currentEquity) > pct × bankroll.
	MaxDrawdownPct float64
	// FeedSilenceSec: trip the feed-silence breaker if the WSS is disconnected
	// AND no book/trade event arrived in this many seconds. Silence alone on a
	// healthy connection (quiet off-hours market) does NOT trip — SPEC §6
	// revised 2026-04-20: nighttime LoL/NBA gaps routinely exceed 60s without
	// a live trade, so we require an actual socket drop as the primary trigger.
	FeedSilenceSec int
	// FeedConnected, when non-nil, reports whether the WSS session is live.
	// CheckFeed only trips when this returns false. If nil (e.g. unit tests
	// that don't wire a client), legacy "trip on pure silence" behavior is
	// preserved.
	FeedConnected func() bool
	// Location for "today" rollover. Defaults to Asia/Singapore.
	Loc *time.Location
}

// DefaultConfig tracks SPEC defaults. Callers may override pieces.
func DefaultConfig() Config {
	loc, _ := time.LoadLocation("Asia/Singapore")
	return Config{
		StartingBankrollUSD: 90.41,
		DailyLossPct:        0.15,
		MaxSingleLossUSD:    3.0,
		MaxDrawdownPct:      0.15,
		FeedSilenceSec:      120,
		Loc:                 loc,
	}
}

// State is the serializable snapshot emitted for logs / heartbeat / reports.
type State struct {
	Day             string // YYYY-MM-DD in cfg.Loc
	DayRealizedPnL  float64
	DayLossCapUSD   float64 // cfg.StartingBankrollUSD × DailyLossPct (positive number)
	CumulativePnL   float64
	PeakEquity      float64
	CurrentEquity   float64
	DrawdownUSD     float64
	DrawdownCapUSD  float64
	Blocked         bool
	BlockReason     BlockReason
	BlockedAt       time.Time
	LastFeedAt      time.Time
	FeedSilentFor   time.Duration
	SingleLossFlags int // count of trades whose loss exceeded MaxSingleLossUSD today
}

// RealizedResult is one durable close used to rebuild financial risk state.
type RealizedResult struct {
	PnLUSD float64
	At     time.Time
}

// Manager is concurrent-safe. Fields mutated under mu.
type Manager struct {
	cfg Config
	mu  sync.Mutex

	day             string
	dayRealized     float64
	cumulativePnL   float64
	peakEquity      float64
	blocked         bool
	blockReason     BlockReason
	blockedAt       time.Time
	lastFeedAt      time.Time
	singleLossFlags int
}

var (
	ErrBlocked = errors.New("risk blocked: opens paused")
)

// New constructs a manager primed to "now".
func New(cfg Config, now time.Time) *Manager {
	if cfg.Loc == nil {
		cfg.Loc = time.UTC
	}
	m := &Manager{
		cfg:        cfg,
		peakEquity: cfg.StartingBankrollUSD,
	}
	m.rolloverLocked(now)
	m.lastFeedAt = now
	return m
}

// AllowOpen returns nil if new positions may be opened, or a wrapped
// ErrBlocked with the current BlockReason. Callers should skip the open on
// error and surface the reason to logs.
func (m *Manager) AllowOpen(now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rolloverLocked(now)
	if m.blocked {
		return fmt.Errorf("%w: %s", ErrBlocked, m.blockReason)
	}
	return nil
}

// OnClose records a realized PnL result from the position manager and trips
// the daily-loss breaker if the accumulated day loss crosses the cap.
// Returns true if this call flipped the breaker to blocked.
func (m *Manager) OnClose(pnlUSD float64, at time.Time) (trippedNow bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rolloverLocked(at)

	m.dayRealized += pnlUSD
	m.cumulativePnL += pnlUSD
	if pnlUSD < -m.cfg.MaxSingleLossUSD {
		m.singleLossFlags++
	}

	equity := m.cfg.StartingBankrollUSD + m.cumulativePnL
	if equity > m.peakEquity {
		m.peakEquity = equity
	}

	if !m.blocked {
		dayCap := -m.dailyCapAbs()
		if m.dayRealized <= dayCap {
			m.blocked = true
			m.blockReason = BlockDailyLoss
			m.blockedAt = at
			return true
		}
		if m.cfg.MaxDrawdownPct > 0 {
			ddCap := m.cfg.StartingBankrollUSD * m.cfg.MaxDrawdownPct
			if m.peakEquity-equity > ddCap {
				m.blocked = true
				m.blockReason = BlockDrawdown
				m.blockedAt = at
				return true
			}
		}
	}
	return false
}

// RebuildRealized replaces financial counters from the durable trade journal.
// Operational pauses and feed-silence blocks remain in force; daily-loss and
// drawdown blocks are recomputed from the supplied PnL history.
func (m *Manager) RebuildRealized(results []RealizedResult, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ordered := append([]RealizedResult(nil), results...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].At.Before(ordered[j].At) })
	day := now.In(m.cfg.Loc).Format("2006-01-02")
	preserveOperationalBlock := m.blocked && (m.blockReason == BlockManualPause || m.blockReason == BlockFeedSilence)

	m.day = day
	m.dayRealized = 0
	m.cumulativePnL = 0
	m.peakEquity = m.cfg.StartingBankrollUSD
	m.singleLossFlags = 0
	for _, result := range ordered {
		at := result.At
		if at.IsZero() {
			at = now
		}
		m.cumulativePnL += result.PnLUSD
		if at.In(m.cfg.Loc).Format("2006-01-02") == day {
			m.dayRealized += result.PnLUSD
			if result.PnLUSD < -m.cfg.MaxSingleLossUSD {
				m.singleLossFlags++
			}
		}
		equity := m.cfg.StartingBankrollUSD + m.cumulativePnL
		if equity > m.peakEquity {
			m.peakEquity = equity
		}
	}

	if preserveOperationalBlock {
		return
	}
	m.blocked = false
	m.blockReason = BlockNone
	m.blockedAt = time.Time{}
	if m.dayRealized <= -m.dailyCapAbs() {
		m.blocked = true
		m.blockReason = BlockDailyLoss
		m.blockedAt = now
		return
	}
	if m.cfg.MaxDrawdownPct > 0 {
		equity := m.cfg.StartingBankrollUSD + m.cumulativePnL
		if m.peakEquity-equity > m.cfg.StartingBankrollUSD*m.cfg.MaxDrawdownPct {
			m.blocked = true
			m.blockReason = BlockDrawdown
			m.blockedAt = now
		}
	}
}

// OnFeedHeartbeat is called whenever a book/trade event arrives. Keeps the
// watchdog clock healthy.
func (m *Manager) OnFeedHeartbeat(at time.Time) {
	m.mu.Lock()
	m.lastFeedAt = at
	// A returned-from-silence heartbeat is NOT enough to auto-resume:
	// once the breaker has tripped, SPEC says boss resumes manually.
	m.mu.Unlock()
}

// CheckFeed evaluates feed health. Trips BlockFeedSilence only when the WSS
// probe reports disconnected AND silence has exceeded cfg.FeedSilenceSec. A
// healthy socket with a quiet market (common during off-hours) never trips.
// If no probe was configured, falls back to the legacy "trip on silence"
// behavior so older tests / wiring keep working. Returns (silentFor,
// trippedNow).
func (m *Manager) CheckFeed(now time.Time) (silentFor time.Duration, trippedNow bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rolloverLocked(now)
	silentFor = now.Sub(m.lastFeedAt)
	if m.blocked {
		return silentFor, false
	}
	if silentFor < time.Duration(m.cfg.FeedSilenceSec)*time.Second {
		return silentFor, false
	}
	// Connected + quiet is fine — don't trip.
	if m.cfg.FeedConnected != nil && m.cfg.FeedConnected() {
		return silentFor, false
	}
	m.blocked = true
	m.blockReason = BlockFeedSilence
	m.blockedAt = now
	trippedNow = true
	return silentFor, trippedNow
}

// Pause manually blocks new opens (e.g. boss sends a pause command).
// Idempotent — keeps the earliest BlockedAt if already blocked.
func (m *Manager) Pause(at time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.blocked {
		m.blocked = true
		m.blockReason = BlockManualPause
		m.blockedAt = at
	}
}

// Resume clears the breaker. Does NOT reset daily realized PnL — if the
// breaker tripped on daily loss and the day hasn't rolled over, resuming
// means the boss has explicitly acknowledged the drawdown.
func (m *Manager) Resume() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blocked = false
	m.blockReason = BlockNone
	m.blockedAt = time.Time{}
}

// State returns a copy for observability.
func (m *Manager) State() State {
	m.mu.Lock()
	defer m.mu.Unlock()
	silent := time.Duration(0)
	if !m.lastFeedAt.IsZero() {
		silent = time.Since(m.lastFeedAt)
	}
	equity := m.cfg.StartingBankrollUSD + m.cumulativePnL
	ddCapUSD := m.cfg.StartingBankrollUSD * m.cfg.MaxDrawdownPct
	dd := m.peakEquity - equity
	if dd < 0 {
		dd = 0
	}
	return State{
		Day:             m.day,
		DayRealizedPnL:  m.dayRealized,
		DayLossCapUSD:   m.dailyCapAbs(),
		CumulativePnL:   m.cumulativePnL,
		PeakEquity:      m.peakEquity,
		CurrentEquity:   equity,
		DrawdownUSD:     dd,
		DrawdownCapUSD:  ddCapUSD,
		Blocked:         m.blocked,
		BlockReason:     m.blockReason,
		BlockedAt:       m.blockedAt,
		LastFeedAt:      m.lastFeedAt,
		FeedSilentFor:   silent,
		SingleLossFlags: m.singleLossFlags,
	}
}

// rolloverLocked resets daily counters when the local date flips. Breaker
// state does NOT automatically reset on rollover — boss still has to Resume
// explicitly, matching the SPEC "wait for manual recovery" contract. Rollover
// only gives us a clean daily PnL ledger.
func (m *Manager) rolloverLocked(now time.Time) {
	day := now.In(m.cfg.Loc).Format("2006-01-02")
	if day == m.day {
		return
	}
	m.day = day
	m.dayRealized = 0
	m.singleLossFlags = 0
}

func (m *Manager) dailyCapAbs() float64 {
	return m.cfg.StartingBankrollUSD * m.cfg.DailyLossPct
}

// persistedState is the on-disk format for surviving restarts.
type persistedState struct {
	Day           string      `json:"day"`
	DayRealized   float64     `json:"day_realized"`
	CumulativePnL float64     `json:"cumulative_pnl"`
	PeakEquity    float64     `json:"peak_equity"`
	Blocked       bool        `json:"blocked"`
	BlockReason   BlockReason `json:"block_reason,omitempty"`
	BlockedAt     time.Time   `json:"blocked_at,omitempty"`
}

// SaveState writes the current risk state to a JSON file so it survives daemon restarts.
func (m *Manager) SaveState(path string) error {
	m.mu.Lock()
	ps := persistedState{
		Day:           m.day,
		DayRealized:   m.dayRealized,
		CumulativePnL: m.cumulativePnL,
		PeakEquity:    m.peakEquity,
		Blocked:       m.blocked,
		BlockReason:   m.blockReason,
		BlockedAt:     m.blockedAt,
	}
	m.mu.Unlock()
	data, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".risk-state-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	dirHandle, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	syncErr := dirHandle.Sync()
	closeErr := dirHandle.Close()
	return errors.Join(syncErr, closeErr)
}

// LoadState restores cumulative and breaker state across every restart. Daily
// PnL is restored only when the snapshot belongs to the current local day.
func (m *Manager) LoadState(path string, now time.Time) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var ps persistedState
	if err := json.Unmarshal(data, &ps); err != nil {
		return err
	}
	today := now.In(m.cfg.Loc).Format("2006-01-02")
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cumulativePnL = ps.CumulativePnL
	if ps.PeakEquity > 0 {
		m.peakEquity = ps.PeakEquity
	}
	m.blocked = ps.Blocked
	m.blockReason = ps.BlockReason
	m.blockedAt = ps.BlockedAt
	if m.blocked && m.blockReason == BlockNone {
		m.blockReason = BlockManualPause
	}
	if ps.Day == today {
		m.dayRealized = ps.DayRealized
	}
	return nil
}
