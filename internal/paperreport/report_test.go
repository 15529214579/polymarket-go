package paperreport

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func TestAnalyzeGroupsTranchesAndOpenExposure(t *testing.T) {
	trades := []journal.TradeRecord{
		{ID: "p1", SizeUSD: 5, PnLUSD: 2, EntryFeeUSD: 0.1, ExitFeeUSD: 0.1, NetPnLUSD: 1.8, SignalSource: "copytrade_w_a", PolicyVersion: "v1"},
		{ID: "p1.t2", SizeUSD: 15, PnLUSD: -1, EntryFeeUSD: 0.3, ExitFeeUSD: 0.2, NetPnLUSD: -1.5, SignalSource: "copytrade_w_a", PolicyVersion: "v1"},
		{ID: "p2", SizeUSD: 5, PnLUSD: -2, SignalSource: "auto"},
		{ID: "p3", SizeUSD: 20, PnLUSD: 4, EntryFeeUSD: 0.4, ExitFeeUSD: 0.2, NetPnLUSD: 3.4, SignalSource: "copytrade_football_score_w_b", PolicyVersion: "v1"},
	}
	open := []strategy.Position{{ID: "p4", SizeUSD: 20, OpenFeeUSD: 0.5, EntryFeeChargedUSD: 0.1}}
	report := Analyze(trades, open)
	if report.ClosedPositions != 3 || report.Records != 4 {
		t.Fatalf("positions=%d records=%d", report.ClosedPositions, report.Records)
	}
	if report.RealizedNetPnLUSD != 1.7 || report.FeesUSD != 1.3 {
		t.Fatalf("net=%v fees=%v", report.RealizedNetPnLUSD, report.FeesUSD)
	}
	if report.OpenExposureUSD != 20 || report.ConservativeOpenNetPnLUSD != -20.4 || report.ConservativeTotalNetPnLUSD != -18.7 {
		t.Fatalf("open=%+v", report)
	}
	if got := findCohort(report.ByStake, "20U+"); got.Positions != 2 || got.NetPnL != 3.7 {
		t.Fatalf("20U=%+v", got)
	}
	if got := findCohort(report.ByStrategy, "football_score"); got.Positions != 1 || got.NetPnL != 3.4 {
		t.Fatalf("score=%+v", got)
	}
	if report.GeneratedAt.Before(time.Now().Add(-time.Minute)) {
		t.Fatal("generated timestamp is stale")
	}
}

func TestAnalyzeAttributesLegacyAutoRemainderToCopytradeSource(t *testing.T) {
	entry := time.Date(2026, 7, 26, 23, 38, 31, 0, time.FixedZone("SGT", 8*60*60))
	trades := []journal.TradeRecord{
		{ID: "p40", EntryTime: entry, SizeUSD: 4, PnLUSD: -3, SignalSource: "copytrade_w_141a...d05a"},
		{ID: "p40", EntryTime: entry, SizeUSD: 1, PnLUSD: 0, SignalSource: "auto"},
	}

	report := Analyze(trades, nil)
	if report.ClosedPositions != 1 {
		t.Fatalf("closed positions=%d", report.ClosedPositions)
	}
	if got := findCohort(report.BySignalSource, "copytrade_w_141a...d05a"); got.Positions != 1 || got.Records != 2 || got.NetPnL != -3 {
		t.Fatalf("source=%+v", got)
	}
	if got := findCohort(report.ByStrategy, "copytrade"); got.Positions != 1 || got.Records != 2 {
		t.Fatalf("strategy=%+v", got)
	}
	if got := findCohort(report.BySignalSource, "mixed"); got.Positions != 0 {
		t.Fatalf("unexpected mixed source=%+v", got)
	}
}

func TestAnalyzeWalletPolicyResolvesLegacyAndFullSources(t *testing.T) {
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	var trades []journal.TradeRecord
	for i := 0; i < 10; i++ {
		trades = append(trades, journal.TradeRecord{
			ID:           "winner-" + string(rune('a'+i)),
			EntryTime:    entry.Add(time.Duration(i) * time.Minute),
			SizeUSD:      20,
			PnLUSD:       1,
			NetPnLUSD:    0.8,
			EntryFeeUSD:  0.1,
			ExitFeeUSD:   0.1,
			SignalSource: "copytrade_wallet:0x1111111111111111111111111111111111111111",
		})
		trades = append(trades, journal.TradeRecord{
			ID:           "loser-" + string(rune('a'+i)),
			EntryTime:    entry.Add(time.Duration(i) * time.Minute),
			SizeUSD:      5,
			PnLUSD:       -1,
			SignalSource: "copytrade_w_2222...2222",
		})
	}

	report := AnalyzeWalletPolicy(trades, map[string]string{
		"w_2222...2222": "0x2222222222222222222222222222222222222222",
	}, WalletPolicyConfig{})
	if report.Promoted != 1 || report.Demoted != 1 || report.Unresolved != 0 {
		t.Fatalf("policy=%+v", report)
	}
	if report.Wallets[0].Decision != "promote" || report.Wallets[1].Decision != "demote" {
		t.Fatalf("wallets=%+v", report.Wallets)
	}
}

func TestWalletFromSourceSupportsBroadCollection(t *testing.T) {
	wallet := "0x3333333333333333333333333333333333333333"
	for _, source := range []string{
		"copytrade_collect_wallet:" + wallet,
		"copytrade_collect_football_score_wallet:" + wallet,
	} {
		if got := walletFromSource(source, nil); got != wallet {
			t.Fatalf("walletFromSource(%q)=%q, want %q", source, got, wallet)
		}
	}
}

func TestAnalyzeSeparatesBroadCollectionCohorts(t *testing.T) {
	trades := []journal.TradeRecord{
		{ID: "core", SizeUSD: 20, NetPnLUSD: 1, SignalSource: "copytrade_wallet:0x1111111111111111111111111111111111111111"},
		{ID: "collect", SizeUSD: 20, NetPnLUSD: -1, SignalSource: "copytrade_collect_wallet:0x2222222222222222222222222222222222222222"},
		{ID: "score", SizeUSD: 20, NetPnLUSD: 2, SignalSource: "copytrade_collect_football_score_wallet:0x3333333333333333333333333333333333333333"},
	}

	report := Analyze(trades, nil)
	if got := findCohort(report.ByStrategy, "copytrade"); got.Positions != 1 {
		t.Fatalf("core=%+v", got)
	}
	if got := findCohort(report.ByStrategy, "copytrade_collect"); got.Positions != 1 || got.NetPnL != -1 {
		t.Fatalf("collection=%+v", got)
	}
	if got := findCohort(report.ByStrategy, "football_score_collect"); got.Positions != 1 || got.NetPnL != 2 {
		t.Fatalf("score collection=%+v", got)
	}
}

func findCohort(rows []Cohort, name string) Cohort {
	for _, row := range rows {
		if row.Name == name {
			return row
		}
	}
	return Cohort{}
}
