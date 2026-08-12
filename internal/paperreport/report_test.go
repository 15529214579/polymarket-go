package paperreport

import (
	"math"
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
	if report.Tradable.Closed.Positions != 3 || math.Abs(report.Tradable.Closed.NetPnL-1.7) > 1e-9 {
		t.Fatalf("tradable=%+v", report.Tradable)
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

func TestAnalyzeWalletPolicyCountsIndependentWalletMarketSamples(t *testing.T) {
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	repeatedWallet := "0x1111111111111111111111111111111111111111"
	independentWallet := "0x2222222222222222222222222222222222222222"
	var trades []journal.TradeRecord
	for i := 0; i < 10; i++ {
		trades = append(trades,
			journal.TradeRecord{
				ID: "repeat-" + string(rune('a'+i)), Market: "same-market", EntryTime: entry.Add(time.Duration(i) * time.Minute),
				SizeUSD: 20, NetPnLUSD: 1, SignalSource: "copytrade_wallet:" + repeatedWallet,
			},
			journal.TradeRecord{
				ID: "independent-" + string(rune('a'+i)), Market: "market-" + string(rune('a'+i)), EntryTime: entry.Add(time.Duration(i) * time.Minute),
				SizeUSD: 20, NetPnLUSD: 1, SignalSource: "copytrade_wallet:" + independentWallet,
			},
		)
	}

	report := AnalyzeWalletPolicy(trades, nil, WalletPolicyConfig{})
	repeated := findWallet(report.Wallets, repeatedWallet)
	if repeated.Positions != 10 || repeated.IndependentSamples != 1 || repeated.Decision != "collect" || repeated.Wins != 1 {
		t.Fatalf("repeated=%+v", repeated)
	}
	independent := findWallet(report.Wallets, independentWallet)
	if independent.Positions != 10 || independent.IndependentSamples != 10 || independent.Decision != "promote" || independent.Wins != 10 {
		t.Fatalf("independent=%+v", independent)
	}
	if report.Promoted != 1 {
		t.Fatalf("promoted=%d", report.Promoted)
	}
}

func TestAnalyzeWalletPolicyEmergencyDemotesSevereSmallSampleLoss(t *testing.T) {
	wallet := "0x7777777777777777777777777777777777777777"
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	trades := []journal.TradeRecord{
		{ID: "loss-a", Market: "market-a", EntryTime: entry, SizeUSD: 20, NetPnLUSD: -11, SignalSource: "copytrade_collect_wallet:" + wallet},
		{ID: "loss-b", Market: "market-b", EntryTime: entry.Add(time.Hour), SizeUSD: 20, NetPnLUSD: -10, SignalSource: "copytrade_collect_wallet:" + wallet},
	}
	report := AnalyzeWalletPolicy(trades, nil, WalletPolicyConfig{})
	row := findWallet(report.Wallets, wallet)
	if row.Decision != "demote" || report.Demoted != 1 {
		t.Fatalf("severe small-sample wallet=%+v report=%+v", row, report)
	}
}

func TestAnalyzeWalletPolicyKeepsOutlierDependentWalletOutOfCore(t *testing.T) {
	wallet := "0x4444444444444444444444444444444444444444"
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	var trades []journal.TradeRecord
	for i := 0; i < 10; i++ {
		net := -1.0
		if i == 0 {
			net = 20
		}
		trades = append(trades, journal.TradeRecord{
			ID: "outlier-" + string(rune('a'+i)), Market: "market-" + string(rune('a'+i)),
			EntryTime: entry.Add(time.Duration(i) * time.Minute), SizeUSD: 20, NetPnLUSD: net,
			SignalSource: "copytrade_wallet:" + wallet,
		})
	}

	report := AnalyzeWalletPolicy(trades, nil, WalletPolicyConfig{})
	row := findWallet(report.Wallets, wallet)
	if row.Decision != "keep" || row.TrimmedNetPnL != -9 || row.BestSampleShare <= 100 {
		t.Fatalf("outlier wallet=%+v", row)
	}
}

func TestAnalyzeWalletPolicyKeepsTwoSidedWalletOutOfCore(t *testing.T) {
	wallet := "0x5555555555555555555555555555555555555555"
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	var trades []journal.TradeRecord
	for i := 0; i < 10; i++ {
		market := "market-" + string(rune('a'+i))
		trades = append(trades, journal.TradeRecord{
			ID: "side-a-" + string(rune('a'+i)), Market: market, AssetID: "yes-" + market,
			EntryTime: entry.Add(time.Duration(i) * time.Minute), SizeUSD: 20, NetPnLUSD: 1,
			SignalSource: "copytrade_wallet:" + wallet,
		})
		if i == 0 {
			trades = append(trades, journal.TradeRecord{
				ID: "side-b", Market: market, AssetID: "no-" + market,
				EntryTime: entry.Add(30 * time.Second), SizeUSD: 20, NetPnLUSD: 1,
				SignalSource: "copytrade_wallet:" + wallet,
			})
		}
	}

	report := AnalyzeWalletPolicy(trades, nil, WalletPolicyConfig{})
	row := findWallet(report.Wallets, wallet)
	if row.Decision != "keep" || row.TwoSidedMarkets != 1 || row.IndependentSamples != 10 {
		t.Fatalf("two-sided wallet=%+v", row)
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

func TestWalletFromSourceResolvesLegacyBroadCollectionAlias(t *testing.T) {
	wallet := "0x6666666666666666666666666666666666666666"
	aliases := map[string]string{"legacy-label": wallet}
	for _, source := range []string{
		"copytrade_collect_legacy-label",
		"copytrade_collect_football_score_legacy-label",
	} {
		if got := walletFromSource(source, aliases); got != wallet {
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
	if report.Tradable.Closed.Positions != 1 || report.Tradable.Closed.NetPnL != 1 {
		t.Fatalf("tradable=%+v", report.Tradable)
	}
	if report.BroadCollection.Closed.Positions != 2 || report.BroadCollection.Closed.NetPnL != 1 {
		t.Fatalf("collection=%+v", report.BroadCollection)
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

func findWallet(rows []WalletPerformance, wallet string) WalletPerformance {
	for _, row := range rows {
		if row.Wallet == wallet {
			return row
		}
	}
	return WalletPerformance{}
}
