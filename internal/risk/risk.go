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
	"errors"
	"fmt"
	"sync"
	"time"
)

// BlockReason enumerates why the manager currently refuses new opens.
type BlockReason string

const (
	BlockNone        BlockReason = ""
	BlockDailyLoss   BlockReason = "daily_loss"
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
	// FeedSilenceSec: trip the feed-silence breaker if no book/trade event for
	// this many seconds. SPEC §6: 30s → close all positions.
	FeedSilenceSec int
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
		FeedSilenceSec:      30,
		Loc:                 loc,
	}
}

// State is the serializable snapshot emitted for logs / heartbeat / reports.
type State struct {
	Day             string // YYYY-MM-DD in cfg.Loc
	DayRealizedPnL  float64
	DayLossCapUSD   float64 // cfg.StartingBankrollUSD × DailyLossPct (positive number)
	Blocked         bool
	BlockReason     BlockReason
	BlockedAt       time.Time
	LastFeedAt      time.Time
	FeedSilentFor   time.Duration
	SingleLossFlags int // count of trades whose loss exceeded MaxSingleLossUSD today
}

// Manager is concurrent-safe. Fields mutated under mu.
type Manager struct {
	cfg Config
	mu  sync.Mutex

	day             string
	dayRealized     float64
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
	m := &Manager{cfg: cfg}
	m.rolloverLocked(now)
	// Seed feed heartbeat so CheckFeed doesn't immediately trip before we've
	// seen our first book event.
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
	if pnlUSD < -m.cfg.MaxSingleLossUSD {
		m.singleLossFlags++
	}

	cap := -m.dailyCapAbs() // negative threshold
	if !m.blocked && m.dayRealized <= cap {
		m.blocked = true
		m.blockReason = BlockDailyLoss
		m.blockedAt = at
		return true
	}
	return false
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

// CheckFeed evaluates feed health. If silent ≥ cfg.FeedSilenceSec and not
// already blocked, trips BlockFeedSilence. Returns (silentFor, trippedNow).
func (m *Manager) CheckFeed(now time.Time) (silentFor time.Duration, trippedNow bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	silentFor = now.Sub(m.lastFeedAt)
	if m.blocked {
		return silentFor, false
	}
	if silentFor >= time.Duration(m.cfg.FeedSilenceSec)*time.Second {
		m.blocked = true
		m.blockReason = BlockFeedSilence
		m.blockedAt = now
		trippedNow = true
	}
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
	return State{
		Day:             m.day,
		DayRealizedPnL:  m.dayRealized,
		DayLossCapUSD:   m.dailyCapAbs(),
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
