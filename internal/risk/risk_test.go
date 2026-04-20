package risk

import (
	"errors"
	"testing"
	"time"
)

func testCfg() Config {
	loc, _ := time.LoadLocation("Asia/Singapore")
	return Config{
		StartingBankrollUSD: 90.0,
		DailyLossPct:        0.15, // cap = -13.5 USDC
		MaxSingleLossUSD:    3.0,
		FeedSilenceSec:      30,
		Loc:                 loc,
	}
}

func TestAllowOpen_FreshManager(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), now)
	if err := m.AllowOpen(now); err != nil {
		t.Fatalf("fresh manager should allow opens, got %v", err)
	}
}

func TestOnClose_TripsDailyLossBreaker(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), now)

	// Two -5 USDC closes: still under -13.5 cap → not tripped.
	if tripped := m.OnClose(-5, now); tripped {
		t.Fatal("should not trip on first -5")
	}
	if tripped := m.OnClose(-5, now); tripped {
		t.Fatal("should not trip on second -5")
	}
	if err := m.AllowOpen(now); err != nil {
		t.Fatalf("should still allow opens at -10, got %v", err)
	}

	// Third -5 → -15 < -13.5 → trip.
	if tripped := m.OnClose(-5, now); !tripped {
		t.Fatal("should trip after crossing daily loss cap")
	}
	err := m.AllowOpen(now)
	if !errors.Is(err, ErrBlocked) {
		t.Fatalf("expected ErrBlocked, got %v", err)
	}
	if got := m.State().BlockReason; got != BlockDailyLoss {
		t.Fatalf("expected BlockDailyLoss, got %q", got)
	}
}

func TestOnClose_CountsSingleLossFlags(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), now)
	m.OnClose(-1, now)    // not flagged
	m.OnClose(-3, now)    // boundary: MaxSingleLossUSD=3 → -3 is NOT < -3
	m.OnClose(-3.01, now) // flagged
	m.OnClose(-5, now)    // flagged
	if got := m.State().SingleLossFlags; got != 2 {
		t.Fatalf("expected 2 single-loss flags, got %d", got)
	}
}

func TestCheckFeed_TripsOnSilence(t *testing.T) {
	start := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), start)
	m.OnFeedHeartbeat(start)

	// 29s silent: not tripped.
	silent, tripped := m.CheckFeed(start.Add(29 * time.Second))
	if tripped {
		t.Fatalf("29s should not trip, got silentFor=%v", silent)
	}

	// 30s silent: trips.
	silent, tripped = m.CheckFeed(start.Add(30 * time.Second))
	if !tripped {
		t.Fatalf("30s should trip, got silentFor=%v", silent)
	}
	if got := m.State().BlockReason; got != BlockFeedSilence {
		t.Fatalf("expected BlockFeedSilence, got %q", got)
	}
}

func TestCheckFeed_AlreadyBlockedIsNoop(t *testing.T) {
	start := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), start)
	m.Pause(start)
	// Even if silent, CheckFeed must not override a manual pause reason.
	_, tripped := m.CheckFeed(start.Add(60 * time.Second))
	if tripped {
		t.Fatal("CheckFeed should not re-trip when already blocked")
	}
	if got := m.State().BlockReason; got != BlockManualPause {
		t.Fatalf("expected BlockManualPause preserved, got %q", got)
	}
}

func TestResume_ClearsBreakerButKeepsDailyLedger(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), now)
	m.OnClose(-20, now) // trips
	if err := m.AllowOpen(now); !errors.Is(err, ErrBlocked) {
		t.Fatalf("expected blocked, got %v", err)
	}

	m.Resume()
	if err := m.AllowOpen(now); err != nil {
		t.Fatalf("should allow opens after Resume, got %v", err)
	}
	// Daily PnL still reflects the loss — not a new day.
	if got := m.State().DayRealizedPnL; got != -20 {
		t.Fatalf("Resume must not wipe ledger, got %v", got)
	}
}

func TestRollover_ResetsLedgerOnNewDay(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Singapore")
	cfg := testCfg()
	cfg.Loc = loc

	// 23:30 SGT day 1
	d1 := time.Date(2026, 4, 20, 15, 30, 0, 0, time.UTC) // 23:30 SGT
	m := New(cfg, d1)
	m.OnClose(-8, d1)
	if got := m.State().DayRealizedPnL; got != -8 {
		t.Fatalf("d1 PnL should be -8, got %v", got)
	}

	// 00:30 SGT next day
	d2 := time.Date(2026, 4, 20, 16, 30, 0, 0, time.UTC) // 00:30 SGT next day
	m.OnClose(-1, d2)
	st := m.State()
	if st.DayRealizedPnL != -1 {
		t.Fatalf("d2 PnL should be -1 after rollover, got %v", st.DayRealizedPnL)
	}
	if st.Day != "2026-04-21" {
		t.Fatalf("expected day=2026-04-21, got %q", st.Day)
	}
}

func TestRollover_DoesNotAutoResumeBreaker(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Singapore")
	cfg := testCfg()
	cfg.Loc = loc

	d1 := time.Date(2026, 4, 20, 15, 30, 0, 0, time.UTC) // 23:30 SGT
	m := New(cfg, d1)
	m.OnClose(-20, d1) // trip

	d2 := time.Date(2026, 4, 20, 16, 30, 0, 0, time.UTC) // 00:30 SGT next day
	if err := m.AllowOpen(d2); !errors.Is(err, ErrBlocked) {
		t.Fatalf("breaker must persist across day rollover until manual Resume, got %v", err)
	}
}

func TestDailyCap_Calculation(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	m := New(testCfg(), now)
	if got := m.State().DayLossCapUSD; got != 13.5 {
		t.Fatalf("expected cap 13.5, got %v", got)
	}
}
