package journal

import (
	"encoding/json"
	"os"
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
	j.SetPolicyVersion("paper-v3")

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
	if got[0].ID != "p1" || got[0].PnLUSD != 0.3125 || got[0].SignalSource != "auto" || got[0].PolicyVersion != "paper-v3" {
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

func TestJournal_PartitionsAndReadsByExitDay(t *testing.T) {
	dir := t.TempDir()
	j, _ := New(dir)
	t.Cleanup(func() { _ = j.Close() })

	rec := TradeRecord{
		ID:        "cross-midnight",
		EntryTime: sgt(2026, 4, 20, 23, 58),
		ExitTime:  sgt(2026, 4, 21, 0, 3),
	}
	if err := j.Append(rec); err != nil {
		t.Fatal(err)
	}
	entryDay, _ := Read(dir, "2026-04-20")
	exitDay, _ := Read(dir, "2026-04-21")
	if len(entryDay) != 0 || len(exitDay) != 1 || exitDay[0].ID != rec.ID {
		t.Fatalf("entry=%+v exit=%+v", entryDay, exitDay)
	}

	// Simulate a legacy entry-day partition and verify reports still bucket it
	// by the realized day.
	legacy := TradeRecord{ID: "legacy", EntryTime: sgt(2026, 4, 19, 23, 0), ExitTime: sgt(2026, 4, 21, 1, 0)}
	legacyPath := Path(dir, "2026-04-19")
	legacyJSON, _ := json.Marshal(legacy)
	if err := os.WriteFile(legacyPath, append(legacyJSON, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	closed, err := ReadClosedDay(dir, "2026-04-21")
	if err != nil {
		t.Fatal(err)
	}
	if len(closed) != 2 || closed[0].ID != rec.ID || closed[1].ID != legacy.ID {
		t.Fatalf("closed day: %+v", closed)
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
	// Headline stats = auto only (3 trades, manual excluded).
	if s.Trades != 3 {
		t.Fatalf("trades=%d", s.Trades)
	}
	if s.Wins != 1 || s.Losses != 1 || s.Breakevens != 1 {
		t.Fatalf("w/l/b = %d/%d/%d", s.Wins, s.Losses, s.Breakevens)
	}
	if absDiff(s.RealizedPnLUSD, 0.30) > 1e-9 {
		t.Fatalf("realized=%v", s.RealizedPnLUSD)
	}
	if absDiff(s.WinRate, 1.0/2.0) > 1e-9 {
		t.Fatalf("winrate=%v", s.WinRate)
	}
	if s.BiggestWinUSD != 0.50 || s.BiggestLossUSD != -0.20 {
		t.Fatalf("biggest win=%v loss=%v", s.BiggestWinUSD, s.BiggestLossUSD)
	}
	if s.AvgHeldSec != (120+60+30)/3 {
		t.Fatalf("avg held=%d", s.AvgHeldSec)
	}
	if s.ExitReasonCount["timeout"] != 1 || s.ExitReasonCount["stop_loss"] != 1 {
		t.Fatalf("reasons=%+v", s.ExitReasonCount)
	}
	// Per-source stats.
	if s.Auto.Count != 3 || s.Manual.Count != 1 {
		t.Fatalf("auto/manual=%d/%d", s.Auto.Count, s.Manual.Count)
	}
	if absDiff(s.Manual.PnLUSD, 1.00) > 1e-9 {
		t.Fatalf("manual pnl=%v", s.Manual.PnLUSD)
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

func TestSummarizeHeadlineHoldExcludesCollectionTrades(t *testing.T) {
	s := Summarize("2026-08-03", []TradeRecord{
		{PnLUSD: 1, HeldSec: 60, SignalSource: "auto"},
		{PnLUSD: 2, HeldSec: 3600, SignalSource: "copytrade_collect_wallet:0x1"},
	})
	if s.Trades != 1 || s.AvgHeldSec != 60 {
		t.Fatalf("trades=%d avgHeld=%d, want 1/60", s.Trades, s.AvgHeldSec)
	}
}

func TestSummarizeSeparatesBroadCollectionFromHeadline(t *testing.T) {
	s := Summarize("2026-08-03", []TradeRecord{
		{ID: "tradable", SignalSource: "copytrade_wallet:0x1", PnLUSD: 2, NetPnLUSD: 1.5, EntryFeeUSD: 0.25, ExitFeeUSD: 0.25},
		{ID: "collection", SignalSource: "copytrade_collect_wallet:0x2", PnLUSD: 100, NetPnLUSD: 99, EntryFeeUSD: 0.5, ExitFeeUSD: 0.5},
	})
	if s.Trades != 1 || s.RealizedPnLUSD != 1.5 || s.Auto.Count != 2 || s.Auto.PnLUSD != 100.5 {
		t.Fatalf("summary=%+v", s)
	}
	if s.Tradable.Count != 1 || s.Collection.Count != 1 || s.Collection.PnLUSD != 99 {
		t.Fatalf("scopes tradable=%+v collection=%+v", s.Tradable, s.Collection)
	}
	out := FormatTelegram(s)
	if !strings.Contains(out, "可交易已实现净 PnL: +1.5000") || !strings.Contains(out, "宽采集研究: 1 笔，净 PnL +99.0000") {
		t.Fatalf("output=%q", out)
	}
}

func TestFormatTelegram_RendersWinSignAndReasons(t *testing.T) {
	s := Summarize("2026-04-20", []TradeRecord{
		{PnLUSD: 1.5, EntryFeeUSD: 0.10, ExitFeeUSD: 0.20, NetPnLUSD: 1.20, HeldSec: 90, ExitReason: "stop_loss"},
	})
	out := FormatTelegram(s)
	if !strings.Contains(out, "+1.2000 USDC") {
		t.Fatalf("missing positive PnL formatting: %q", out)
	}
	if s.EntryFeesUSD != 0.10 || s.ExitFeesUSD != 0.20 || !strings.Contains(out, "入场 0.1000 / 出场 0.2000") {
		t.Fatalf("missing fee breakdown: summary=%+v output=%q", s, out)
	}
	if !strings.Contains(out, "不含未平仓浮动 PnL") {
		t.Fatalf("missing PnL basis: %q", out)
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
