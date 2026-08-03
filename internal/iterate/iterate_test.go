package iterate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/journal"
)

func TestAnalyzeUsesCompletedDaysAndNetPnL(t *testing.T) {
	dir := t.TempDir()
	j, err := journal.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = j.Close() })

	now := time.Date(2026, 8, 3, 22, 0, 0, 0, journal.SGT)
	yesterday := dayAtNoon(now.AddDate(0, 0, -1))
	prior := dayAtNoon(now.AddDate(0, 0, -2))
	today := dayAtNoon(now)

	records := []journal.TradeRecord{
		{
			ID: "fee-aware-win", Question: "Dota 2: A vs B", SignalSource: "copytrade_wallet:0x1",
			EntryMid: 0.50, EntryTime: yesterday.Add(-30 * time.Minute), ExitTime: yesterday,
			HeldSec: 600, ExitReason: "timeout", PnLUSD: 5, EntryFeeUSD: 1, ExitFeeUSD: 0.5, NetPnLUSD: 3.5,
		},
		{
			ID: "fee-aware-loss", Question: "Will A win?", SignalSource: "copytrade_wallet:0x2",
			EntryMid: 0.70, EntryTime: prior.Add(-30 * time.Minute), ExitTime: prior,
			HeldSec: 600, ExitReason: "timeout", PnLUSD: -2, EntryFeeUSD: 0.2, ExitFeeUSD: 0.3, NetPnLUSD: -2.5,
		},
		{
			ID: "legacy-win", Question: "Will B win?", SignalSource: "auto",
			EntryMid: 0.40, EntryTime: prior.Add(-10 * time.Minute), ExitTime: prior,
			HeldSec: 60, ExitReason: "timeout", PnLUSD: 1,
		},
		{
			ID: "manual-excluded", SignalSource: "manual", EntryTime: yesterday.Add(-time.Minute),
			ExitTime: yesterday, HeldSec: 60, PnLUSD: 99,
		},
		{
			ID: "today-excluded", SignalSource: "auto", EntryTime: today.Add(-time.Minute),
			ExitTime: today, HeldSec: 60, PnLUSD: 100,
		},
	}
	for _, record := range records {
		if err := j.Append(record); err != nil {
			t.Fatal(err)
		}
	}
	if err := j.Close(); err != nil {
		t.Fatal(err)
	}

	report, err := analyzeAt(dir, 2, now)
	if err != nil {
		t.Fatal(err)
	}
	if report.Day != yesterday.Format("2006-01-02") {
		t.Fatalf("day=%s want=%s", report.Day, yesterday.Format("2006-01-02"))
	}
	if report.TotalTrades != 3 || report.TotalWins != 2 || report.TotalLosses != 1 {
		t.Fatalf("trades/wins/losses=%d/%d/%d", report.TotalTrades, report.TotalWins, report.TotalLosses)
	}
	assertClose(t, report.GrossPnLUSD, 4)
	assertClose(t, report.EntryFeesUSD, 1.2)
	assertClose(t, report.ExitFeesUSD, 0.8)
	assertClose(t, report.FeesUSD, 2)
	assertClose(t, report.CumulativePnL, 2)
	assertClose(t, report.AvgPnLPerDay, 1)
	if len(report.DailyBreakdown) != 2 || report.DailyBreakdown[1].Day != report.Day {
		t.Fatalf("daily breakdown=%+v", report.DailyBreakdown)
	}

	markdown := FormatMarkdown(report)
	for _, want := range []string{"已实现净 PnL", "已记录手续费", "不包含未平仓浮动 PnL", "+$2.00"} {
		if !strings.Contains(markdown, want) {
			t.Fatalf("markdown missing %q:\n%s", want, markdown)
		}
	}
	telegram := FormatTelegram(report)
	for _, want := range []string{"已实现净 PnL +$2.00", "手续费 $2.00", "不含未平仓浮动 PnL"} {
		if !strings.Contains(telegram, want) {
			t.Fatalf("telegram missing %q:\n%s", want, telegram)
		}
	}
}

func TestAnalyzeRejectsInvalidWindowAndCorruptJournal(t *testing.T) {
	now := time.Date(2026, 8, 3, 22, 0, 0, 0, journal.SGT)
	if _, err := analyzeAt(t.TempDir(), 0, now); err == nil {
		t.Fatal("expected invalid-window error")
	}

	dir := t.TempDir()
	reportDay := now.AddDate(0, 0, -1).Format("2006-01-02")
	path := filepath.Join(dir, "trades-"+reportDay+".jsonl")
	if err := os.WriteFile(path, []byte("{not-json}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := analyzeAt(dir, 1, now); err == nil || !strings.Contains(err.Error(), reportDay) {
		t.Fatalf("expected dated journal error, got %v", err)
	}
}

func dayAtNoon(t time.Time) time.Time {
	t = t.In(journal.SGT)
	return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, journal.SGT)
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	const epsilon = 1e-9
	if got < want-epsilon || got > want+epsilon {
		t.Fatalf("got %.12f want %.12f", got, want)
	}
}
