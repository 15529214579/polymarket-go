package journal

import (
	"strings"
	"testing"
	"time"
)

func sgt(y int, m time.Month, d, h, mn int) time.Time {
	return time.Date(y, m, d, h, mn, 0, 0, SGT)
}

func TestJournal_AppendAndRead_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	j, err := New(dir)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	t.Cleanup(func() { _ = j.Close() })

	rec := TradeRecord{
		ID: "p1", AssetID: "A", Market: "M", Question: "Foo?",
		Outcome: "Yes", Side: "buy", SizeUSD: 5, Units: 6.25,
		EntryMid: 0.80, EntryTime: sgt(2026, 4, 20, 14, 0),
		ExitMid: 0.85, ExitTime: sgt(2026, 4, 20, 14, 15),
		ExitReason: "reversal_drawdown", HeldSec: 900,
		PnLUSD: 0.3125, OpenOrderID: "paper-aaa", CloseOrderID: "paper-bbb",
		Mode: "paper", SignalSource: "auto",
	}
	if err := j.Append(rec); err != nil {
		t.Fatalf("append: %v", err)
	}

	got, err := Read(dir, "2026-04-20")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	if got[0].ID != "p1" || got[0].PnLUSD != 0.3125 || got[0].SignalSource != "auto" {
		t.Fatalf("round-trip mismatch: %+v", got[0])
	}
}

func TestJournal_RotatesOnSGTDay(t *testing.T) {
	dir := t.TempDir()
	j, _ := New(dir)
	t.Cleanup(func() { _ = j.Close() })

	a := TradeRecord{ID: "a", EntryTime: sgt(2026, 4, 20, 23, 30)}
	b := TradeRecord{ID: "b", EntryTime: sgt(2026, 4, 21, 0, 5)}
	if err := j.Append(a); err != nil {
		t.Fatal(err)
	}
	if err := j.Append(b); err != nil {
		t.Fatal(err)
	}
	day1, _ := Read(dir, "2026-04-20")
	day2, _ := Read(dir, "2026-04-21")
	if len(day1) != 1 || day1[0].ID != "a" {
		t.Fatalf("day1: %+v", day1)
	}
	if len(day2) != 1 || day2[0].ID != "b" {
		t.Fatalf("day2: %+v", day2)
	}
}

func TestRead_MissingFile_NilNoErr(t *testing.T) {
	got, err := Read(t.TempDir(), "2026-04-20")
	if err != nil || got != nil {
		t.Fatalf("want (nil,nil), got (%v,%v)", got, err)
	}
}

func TestSummarize_BasicAggregate(t *testing.T) {
	day := "2026-04-20"
	trades := []TradeRecord{
		{ID: "1", PnLUSD: 0.50, HeldSec: 120, ExitReason: "reversal_ticks", SignalSource: "auto"},
		{ID: "2", PnLUSD: -0.20, HeldSec: 60, ExitReason: "stop_loss", SignalSource: "auto"},
		{ID: "3", PnLUSD: 1.00, HeldSec: 300, ExitReason: "reversal_drawdown", SignalSource: "manual"},
		{ID: "4", PnLUSD: 0.00, HeldSec: 30, ExitReason: "timeout", SignalSource: "auto"},
	}
	s := Summarize(day, trades)
	if s.Trades != 4 {
		t.Fatalf("trades=%d", s.Trades)
	}
	if s.Wins != 2 || s.Losses != 1 || s.Breakevens != 1 {
		t.Fatalf("w/l/b = %d/%d/%d", s.Wins, s.Losses, s.Breakevens)
	}
	if absDiff(s.RealizedPnLUSD, 1.30) > 1e-9 {
		t.Fatalf("realized=%v", s.RealizedPnLUSD)
	}
	if absDiff(s.WinRate, 2.0/3.0) > 1e-9 {
		t.Fatalf("winrate=%v", s.WinRate)
	}
	if s.BiggestWinUSD != 1.00 || s.BiggestLossUSD != -0.20 {
		t.Fatalf("biggest win=%v loss=%v", s.BiggestWinUSD, s.BiggestLossUSD)
	}
	if s.AvgHeldSec != (120+60+300+30)/4 {
		t.Fatalf("avg held=%d", s.AvgHeldSec)
	}
	if s.ExitReasonCount["timeout"] != 1 || s.ExitReasonCount["stop_loss"] != 1 {
		t.Fatalf("reasons=%+v", s.ExitReasonCount)
	}
	if s.AutoCount != 3 || s.ManualCount != 1 {
		t.Fatalf("auto/manual=%d/%d", s.AutoCount, s.ManualCount)
	}
}

func TestSummarize_EmptyDay(t *testing.T) {
	s := Summarize("2026-04-20", nil)
	if s.Trades != 0 {
		t.Fatalf("trades=%d", s.Trades)
	}
	out := FormatTelegram(s)
	if !strings.Contains(out, "无成交") {
		t.Fatalf("missing empty-day marker: %q", out)
	}
}

func TestFormatTelegram_RendersWinSignAndReasons(t *testing.T) {
	s := Summarize("2026-04-20", []TradeRecord{
		{PnLUSD: 1.5, HeldSec: 90, ExitReason: "stop_loss"},
	})
	out := FormatTelegram(s)
	if !strings.Contains(out, "+1.5000 USDC") {
		t.Fatalf("missing positive PnL formatting: %q", out)
	}
	if !strings.Contains(out, "stop_loss×1") {
		t.Fatalf("missing reason: %q", out)
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
