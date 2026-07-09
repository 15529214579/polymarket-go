package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadAlertEvents_DedupesAndSorts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.jsonl")
	data := "" +
		`{"key":"k2","sent_at":"2026-07-09T03:01:00Z","wallet":"0x2222222222222222222222222222222222222222","asset":"asset-2","price":0.5,"mode":"OBSERVE"}` + "\n" +
		`{"key":"k1","sent_at":"2026-07-09T03:00:00Z","wallet":"0x1111111111111111111111111111111111111111","asset":"asset-1","price":0.4,"mode":"PROBATION"}` + "\n" +
		`{"key":"k1","sent_at":"2026-07-09T03:00:00Z","wallet":"0x1111111111111111111111111111111111111111","asset":"asset-1","price":0.4,"mode":"PROBATION"}` + "\n" +
		`{"key":"","sent_at":"2026-07-09T03:00:00Z","wallet":"0x3333333333333333333333333333333333333333","asset":"asset-3","price":0.4}` + "\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadAlertEvents(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("events len=%d, want 2", len(got))
	}
	if got[0].Key != "k1" || got[1].Key != "k2" {
		t.Fatalf("order=%s,%s want k1,k2", got[0].Key, got[1].Key)
	}
}

func TestMergeAlertEvents_DedupesAcrossLogs(t *testing.T) {
	a := []alertEvent{
		{Key: "k2", SentAt: time.Date(2026, 7, 9, 3, 2, 0, 0, time.UTC), Wallet: "0x2", Asset: "asset-2", Price: 0.5},
		{Key: "k1", SentAt: time.Date(2026, 7, 9, 3, 1, 0, 0, time.UTC), Wallet: "0x1", Asset: "asset-1", Price: 0.4},
	}
	b := []alertEvent{
		{Key: "K2", SentAt: time.Date(2026, 7, 9, 3, 3, 0, 0, time.UTC), Wallet: "0x2b", Asset: "asset-2b", Price: 0.6},
		{Key: "k3", SentAt: time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC), Wallet: "0x3", Asset: "asset-3", Price: 0.7},
	}

	got := mergeAlertEvents(a, b)
	if len(got) != 3 {
		t.Fatalf("merged len=%d, want 3: %#v", len(got), got)
	}
	if got[0].Key != "k3" || got[1].Key != "k1" || got[2].Key != "k2" {
		t.Fatalf("merged order=%s,%s,%s want k3,k1,k2", got[0].Key, got[1].Key, got[2].Key)
	}
	if got[2].Wallet != "0x2" {
		t.Fatalf("duplicate winner wallet=%s, want primary event", got[2].Wallet)
	}
}

func TestEvaluateAlerts_ComputesFixedStakeROI(t *testing.T) {
	events := []alertEvent{
		{
			Key:    "k1",
			Wallet: "0x1111111111111111111111111111111111111111",
			Asset:  "asset-1",
			Price:  0.40,
		},
		{
			Key:    "k2",
			Wallet: "0x2222222222222222222222222222222222222222",
			Asset:  "asset-2",
			Price:  0.50,
		},
	}
	results := evaluateAlerts(context.Background(), events, 10, func(_ context.Context, asset string) (float64, error) {
		if asset == "asset-1" {
			return 0.44, nil
		}
		return 0, errors.New("no mark")
	})

	if len(results) != 2 {
		t.Fatalf("results len=%d, want 2", len(results))
	}
	if !results[0].Marked || results[0].PnLUSD < 0.99 || results[0].PnLUSD > 1.01 {
		t.Fatalf("first result=%+v, want marked +$1 pnl", results[0])
	}
	if results[1].Marked {
		t.Fatalf("second result marked, want unmarked")
	}
	sum := summarize(results)
	if sum.Signals != 2 || sum.Marked != 1 || sum.Unmarked != 1 || sum.ReturnPC < 9.9 || sum.ReturnPC > 10.1 {
		t.Fatalf("summary=%+v, want 2/1/1/10%%", sum)
	}
}

func TestEvaluateAlerts_AllowsZeroSettlementMark(t *testing.T) {
	events := []alertEvent{{
		Key:    "k1",
		Wallet: "0x1111111111111111111111111111111111111111",
		Asset:  "asset-1",
		Price:  0.40,
	}}
	results := evaluateAlerts(context.Background(), events, 10, func(_ context.Context, asset string) (float64, error) {
		return 0, nil
	})

	if len(results) != 1 {
		t.Fatalf("results len=%d, want 1", len(results))
	}
	if !results[0].Marked || results[0].Mid != 0 || results[0].PnLUSD != -10 || results[0].ReturnPC != -100 {
		t.Fatalf("result=%+v, want marked total loss at settlement 0", results[0])
	}
}

func TestMidpointCache_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "midpoints.json")
	now := time.Date(2026, 7, 9, 2, 0, 0, 0, time.UTC)
	cache := &midpointCache{}
	cache.Set("asset-1", 0.44, now)
	cache.Set("", 0.50, now)
	cache.Set("asset-2", 0, now)

	if err := saveMidpointCache(path, cache); err != nil {
		t.Fatal(err)
	}
	got, err := loadMidpointCache(path)
	if err != nil {
		t.Fatal(err)
	}
	mid, ok := got.Get("asset-1")
	if !ok || mid != 0.44 {
		t.Fatalf("cached midpoint=%.2f/%v, want 0.44/true", mid, ok)
	}
	mid, ok = got.Get("asset-2")
	if !ok || mid != 0 {
		t.Fatalf("asset-2 midpoint=%.2f/%v, want 0/true", mid, ok)
	}
	if len(got.Assets) != 2 {
		t.Fatalf("cache assets=%d, want 2", len(got.Assets))
	}
}

func TestSettledTokenPriceFromMarkets_UsesClosedGammaOutcomePrice(t *testing.T) {
	markets := []gammaMarket{
		{Closed: false, ClobTokenIDsRaw: `["token-2"]`, OutcomePricesRaw: `["0.44"]`},
		{Closed: true, ClobTokenIDsRaw: `["token-1","token-2"]`, OutcomePricesRaw: `["1","0"]`},
	}
	got, err := settledTokenPriceFromMarkets(markets, "token-2")
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("settled price=%.2f, want 0", got)
	}
}

func TestWriteReport_IncludesSummaryGroupsAndRecentRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{{
		Event: alertEvent{
			SentAt:   time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC),
			Wallet:   "0x1111111111111111111111111111111111111111",
			Mode:     "PROBATION",
			Market:   "Dota 2: Team A vs Team B - Game 1 Winner",
			Notional: 6000,
			Price:    0.40,
		},
		StakeUSD: 10,
		Units:    25,
		Mid:      0.44,
		Marked:   true,
		PnLUSD:   1,
		ReturnPC: 10,
	}}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("PROBATION"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{"# Sports Alert Performance", "Alerts: 1", "ROI incl. midpoint marks: 10.0%", "## Current Policy Performance", "Modes: PROBATION", "Gate action: COLLECT_POSITIVE", "## Current Policy Position-Capped Performance", "Positions: 1", "## Experimental OBSERVE Performance", "## Experimental OBSERVE Position-Capped Performance", "## Mode Gates", "COLLECT_POSITIVE", "## By Mode", "PROBATION", "## Recent Alerts"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestCurrentPolicySummary_ExcludesDisabledObserve(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{
		{
			Event:    alertEvent{Mode: "FLOW-SCOUT", Wallet: "0x1111111111111111111111111111111111111111", Price: 0.20, Market: "Will Argentina win on 2026-07-11?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.22,
			Marked:   true,
			PnLUSD:   1,
			ReturnPC: 10,
		},
		{
			Event:    alertEvent{Mode: "OBSERVE", Wallet: "0x2222222222222222222222222222222222222222", Price: 0.50, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.001,
			Marked:   true,
			PnLUSD:   -9.98,
			ReturnPC: -99.8,
		},
		{
			Event:    alertEvent{Mode: "OBSERVE-BURST", Wallet: "0x4444444444444444444444444444444444444444", Price: 0.31, Market: "Indiana Fever vs. Los Angeles Sparks"},
			StakeUSD: 10,
			Mid:      0.34,
			Marked:   true,
			PnLUSD:   0.97,
			ReturnPC: 9.7,
		},
		{
			Event:    alertEvent{Mode: "INSIDER-SCOUT", Wallet: "0x3333333333333333333333333333333333333333", Price: 0.40, Market: "Will France win on 2026-07-09?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.44,
			Marked:   true,
			PnLUSD:   1,
			ReturnPC: 10,
		},
	}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("FLOW-SCOUT,EDGE-HOT"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{"## Current Policy Performance", "- Alerts: 1", "- PnL incl. midpoint marks: $+1.00", "- ROI incl. midpoint marks: 10.0%", "- Gate action: COLLECT_POSITIVE"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
	for _, want := range []string{"## Experimental OBSERVE Performance", "- Modes: OBSERVE,OBSERVE-BURST", "- Alerts: 2"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing OBSERVE metric %q:\n%s", want, report)
		}
	}
	for _, want := range []string{"## Experimental OBSERVE-BURST Performance", "- Modes: OBSERVE-BURST", "- Alerts: 1", "- PnL incl. midpoint marks: $+0.97", "- ROI incl. midpoint marks: 9.7%", "## Experimental OBSERVE-BURST Position-Capped Performance"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing OBSERVE-BURST metric %q:\n%s", want, report)
		}
	}
	for _, want := range []string{"## Experimental INSIDER-SCOUT Performance", "- Alerts: 1", "- PnL incl. midpoint marks: $+1.00", "- ROI incl. midpoint marks: 10.0%", "- Gate action: COLLECT_POSITIVE"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing INSIDER-SCOUT metric %q:\n%s", want, report)
		}
	}
}

func TestCurrentPolicyPositionCappedSummary_DedupesWalletAsset(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC), Mode: "FLOW-SCOUT", Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-1", Price: 0.20, Market: "Will Argentina win on 2026-07-11?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.22,
			Marked:   true,
			PnLUSD:   1,
			ReturnPC: 10,
		},
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 5, 0, 0, time.UTC), Mode: "FLOW-SCOUT", Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-1", Price: 0.21, Market: "Will Argentina win on 2026-07-11?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.22,
			Marked:   true,
			PnLUSD:   0.48,
			ReturnPC: 4.8,
		},
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 10, 0, 0, time.UTC), Mode: "OBSERVE", Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-2", Price: 0.50, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.001,
			Marked:   true,
			PnLUSD:   -9.98,
			ReturnPC: -99.8,
		},
	}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("FLOW-SCOUT,EDGE-HOT"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{"## Current Policy Performance", "- Alerts: 2", "## Current Policy Position-Capped Performance", "- Positions: 1", "- PnL incl. midpoint marks: $+1.00", "- ROI incl. midpoint marks: 10.0%", "- Gate action: COLLECT_POSITIVE"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
	if strings.Contains(report, "Positions: 2") {
		t.Fatalf("position-capped summary did not dedupe duplicate wallet+asset:\n%s", report)
	}
}

func TestConsensusExperimentSummary_IsSeparateFromCurrentPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC), Mode: "CONSENSUS", Wallet: "multi:2", Asset: "asset-1", Price: 0.30, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.36,
			Marked:   true,
			PnLUSD:   2,
			ReturnPC: 20,
		},
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 5, 0, 0, time.UTC), Mode: "CONSENSUS", Wallet: "multi:5", Asset: "asset-1", Price: 0.32, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.36,
			Marked:   true,
			PnLUSD:   1.25,
			ReturnPC: 12.5,
		},
	}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("FLOW-SCOUT,EDGE-HOT"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental CONSENSUS Performance",
		"- Alerts: 2",
		"- ROI incl. midpoint marks: 16.2%",
		"## Experimental CONSENSUS Position-Capped Performance",
		"- Positions: 1",
		"- ROI incl. midpoint marks: 20.0%",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestCurrentPolicySummary_IncludesConsensusWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC), Mode: "CONSENSUS", Wallet: "multi:2", Asset: "asset-1", Price: 0.30, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.36,
			Marked:   true,
			PnLUSD:   2,
			ReturnPC: 20,
		},
		{
			Event:    alertEvent{SentAt: time.Date(2026, 7, 9, 3, 5, 0, 0, time.UTC), Mode: "OBSERVE", Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-2", Price: 0.50, Market: "Will France win on 2026-07-09?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.49,
			Marked:   true,
			PnLUSD:   -0.2,
			ReturnPC: -2,
		},
	}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("FLOW-SCOUT,EDGE-HOT,CONSENSUS"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Current Policy Performance",
		"- Modes: CONSENSUS,EDGE-HOT,FLOW-SCOUT",
		"- Alerts: 1",
		"- PnL incl. midpoint marks: $+2.00",
		"- ROI incl. midpoint marks: 20.0%",
		"## Experimental OBSERVE Performance",
		"- PnL incl. midpoint marks: $-0.20",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestCurrentPolicySummary_ExcludesBlockedWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	blocked := "0x1111111111111111111111111111111111111111"
	active := "0x2222222222222222222222222222222222222222"
	results := []alertResult{
		{
			Event:    alertEvent{Mode: "FLOW-SCOUT", Wallet: blocked, Price: 0.20, Market: "Will Argentina win on 2026-07-11?", Category: "soccer"},
			StakeUSD: 10,
			Mid:      0.24,
			Marked:   true,
			PnLUSD:   2,
			ReturnPC: 20,
		},
		{
			Event:    alertEvent{Mode: "OBSERVE", Wallet: active, Price: 0.50, Market: "Dota 2: Team A vs Team B"},
			StakeUSD: 10,
			Mid:      0.55,
			Marked:   true,
			PnLUSD:   1,
			ReturnPC: 10,
		},
	}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("FLOW-SCOUT"), map[string]struct{}{blocked: {}}, "review-noise.txt"); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"Historical alerts excluded by wallet lists: 1",
		"## Summary",
		"- Alerts: 2",
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental OBSERVE Performance",
		"- Alerts: 1",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
	if strings.Contains(report, "`0x1111...1111` | 1 | 1") {
		t.Fatalf("blocked wallet still appears in active gate/group sections:\n%s", report)
	}
}

func TestCurrentPolicyMarketFilter_ExcludesOutrightWorldCupFromGates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{{
		Event: alertEvent{
			Mode:     "OBSERVE",
			Wallet:   "0x1111111111111111111111111111111111111111",
			Price:    0.326,
			Market:   "Will France win the 2026 FIFA World Cup?",
			Slug:     "world-cup-winner",
			Category: "soccer",
		},
		StakeUSD: 10,
		Mid:      0.323,
		Marked:   true,
		PnLUSD:   -0.11,
		ReturnPC: -1.1,
	}}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("OBSERVE"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Summary",
		"- Alerts: 1",
		"Historical alerts excluded by current market filter: 1",
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental OBSERVE Performance",
		"- Alerts: 0",
		"## Recent Alerts",
		"Will France win the 2026 FIFA World Cup?",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestUnknownFlowExperimentSummary_IsReportedSeparately(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{{
		Event: alertEvent{
			Mode:     "UNKNOWN-FLOW",
			Wallet:   "0x1111111111111111111111111111111111111111",
			Asset:    "asset-1",
			Price:    0.50,
			Market:   "Argentina vs. Switzerland: Team to Advance",
			Category: "soccer",
		},
		StakeUSD: 10,
		Mid:      0.55,
		Marked:   true,
		PnLUSD:   1,
		ReturnPC: 10,
	}}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("CONSENSUS"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental UNKNOWN-FLOW Performance",
		"- Alerts: 1",
		"- PnL incl. midpoint marks: $+1.00",
		"- ROI incl. midpoint marks: 10.0%",
		"## Experimental UNKNOWN-FLOW Position-Capped Performance",
		"- Positions: 1",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestScoredFlowExperimentSummary_IsReportedSeparately(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{{
		Event: alertEvent{
			Mode:     "SCORED-FLOW",
			Wallet:   "0xb3c1a1801a6c22813a2f701969820fa2e9d54837",
			Asset:    "asset-1",
			Price:    0.50,
			Market:   "Norway vs. England: Team to Advance",
			Category: "soccer",
		},
		StakeUSD: 10,
		Mid:      0.45,
		Marked:   true,
		PnLUSD:   -1,
		ReturnPC: -10,
	}}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("CONSENSUS"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental SCORED-FLOW Performance",
		"- Modes: SCORED-FLOW",
		"- Alerts: 1",
		"- PnL incl. midpoint marks: $-1.00",
		"- ROI incl. midpoint marks: -10.0%",
		"## Experimental SCORED-FLOW Position-Capped Performance",
		"- Positions: 1",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestSeedFlowExperimentSummary_IsReportedSeparately(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	results := []alertResult{{
		Event: alertEvent{
			Mode:     "SEED-FLOW",
			Wallet:   "0x2393e78ad6d21ea3dda86f98980461c971341a8d",
			Asset:    "asset-1",
			Price:    0.62,
			Market:   "Will France win on 2026-07-09?",
			Category: "soccer",
		},
		StakeUSD: 10,
		Mid:      0.67,
		Marked:   true,
		PnLUSD:   1,
		ReturnPC: 10,
	}}

	if err := writeReport(path, "alerts.jsonl", "", results, 10, parseModeSet("CONSENSUS"), nil, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{
		"## Current Policy Performance",
		"- Alerts: 0",
		"## Experimental SEED-FLOW Performance",
		"- Modes: SEED-FLOW",
		"- Alerts: 1",
		"- PnL incl. midpoint marks: $+1.00",
		"- ROI incl. midpoint marks: 10.0%",
		"## Experimental SEED-FLOW Position-Capped Performance",
		"- Positions: 1",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestLoadWalletSet_ExcludesOnlyRejectedStatusRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "observe.txt")
	rejected := "0x1111111111111111111111111111111111111111"
	watch := "0x2222222222222222222222222222222222222222"
	body := rejected + " # list=tape_observe status=reject-flow tier=D\n" +
		watch + " # list=tape_observe status=watch tier=B\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletSet(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got[rejected]; !ok {
		t.Fatalf("missing rejected wallet in exclude set")
	}
	if _, ok := got[watch]; ok {
		t.Fatalf("watch status wallet should not be excluded")
	}
}

func TestGateAction(t *testing.T) {
	tests := []struct {
		name   string
		sum    resultSummary
		action string
	}{
		{
			name:   "positive sample below promote gate",
			sum:    resultSummary{Marked: 2, Wins: 2, StakeUSD: 20, PnLUSD: 1, ReturnPC: 5},
			action: "COLLECT_POSITIVE",
		},
		{
			name:   "promote candidate",
			sum:    resultSummary{Marked: 5, Wins: 3, StakeUSD: 50, PnLUSD: 3, ReturnPC: 6},
			action: "PROMOTE_CANDIDATE",
		},
		{
			name:   "severe single drawdown",
			sum:    resultSummary{Marked: 1, Wins: 0, StakeUSD: 10, PnLUSD: -8, ReturnPC: -80},
			action: "PROBATION",
		},
		{
			name:   "negative marked sample",
			sum:    resultSummary{Marked: 3, Wins: 1, StakeUSD: 30, PnLUSD: -7, ReturnPC: -23.3},
			action: "CUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, _ := gateAction(tt.sum)
			if action != tt.action {
				t.Fatalf("action=%s, want %s", action, tt.action)
			}
		})
	}
}

func TestWriteDecisionJSON_SeparatesPromotedAndCollectModes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decision.json")

	var results []alertResult
	for i := 0; i < 5; i++ {
		results = append(results, alertResult{
			Event: alertEvent{
				Key:      "consensus-" + string(rune('a'+i)),
				SentAt:   time.Date(2026, 7, 9, 3, i, 0, 0, time.UTC),
				Mode:     "CONSENSUS",
				Wallet:   "multi:3",
				Category: "basketball",
				Market:   "Indiana Fever vs. Los Angeles Sparks",
				Asset:    "asset-consensus",
				Price:    0.40,
			},
			StakeUSD: 10,
			Mid:      0.50,
			Marked:   true,
			PnLUSD:   2,
			ReturnPC: 20,
		})
	}
	results = append(results, alertResult{
		Event: alertEvent{
			Key:      "unknown-flow",
			SentAt:   time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC),
			Mode:     "UNKNOWN-FLOW",
			Wallet:   "0xb748f5d2ad58fbb3facd9ea1b41f7be7ade4dbdd",
			Category: "soccer",
			Market:   "Argentina vs. Switzerland: Team to Advance",
			Asset:    "asset-unknown",
			Price:    0.74,
		},
		StakeUSD: 10,
		Mid:      0.735,
		Marked:   true,
		PnLUSD:   -0.07,
		ReturnPC: -0.7,
	})

	if err := writeDecisionJSON(path, "alerts.jsonl", "", results, 10, parseModeSet("CONSENSUS"), nil); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got policyDecisionReport
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got.CurrentPolicy.Action != "PROMOTE_CANDIDATE" {
		t.Fatalf("current action=%s, want PROMOTE_CANDIDATE", got.CurrentPolicy.Action)
	}
	if !containsString(got.RecommendedModes, "CONSENSUS") {
		t.Fatalf("recommended modes=%v, want CONSENSUS", got.RecommendedModes)
	}
	for _, mode := range got.RecommendedModes {
		if mode == "UNKNOWN-FLOW" {
			t.Fatalf("UNKNOWN-FLOW should not be recommended: %+v", got)
		}
	}
}

func TestWriteDecisionJSON_EffectivePolicyExcludesCutModes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decision.json")

	var results []alertResult
	for i := 0; i < 3; i++ {
		results = append(results, alertResult{
			Event: alertEvent{
				Key:      "consensus-loss-" + string(rune('a'+i)),
				SentAt:   time.Date(2026, 7, 9, 3, i, 0, 0, time.UTC),
				Mode:     "CONSENSUS",
				Wallet:   "multi:3",
				Category: "soccer",
				Market:   "Will France win on 2026-07-09?",
				Asset:    "asset-consensus",
				Price:    0.70,
			},
			StakeUSD: 10,
			Mid:      0.50,
			Marked:   true,
			PnLUSD:   -2,
			ReturnPC: -20,
		})
	}
	results = append(results, alertResult{
		Event: alertEvent{
			Key:      "observe-burst-win",
			SentAt:   time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC),
			Mode:     "OBSERVE-BURST",
			Wallet:   "0x1111111111111111111111111111111111111111",
			Category: "basketball",
			Market:   "Indiana Fever vs. Los Angeles Sparks",
			Asset:    "asset-observe",
			Price:    0.30,
		},
		StakeUSD: 10,
		Mid:      0.60,
		Marked:   true,
		PnLUSD:   10,
		ReturnPC: 100,
	})

	if err := writeDecisionJSON(path, "shadow.jsonl", "alerts.jsonl", results, 10, parseModeSet("CONSENSUS,OBSERVE-BURST"), nil); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got policyDecisionReport
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if !containsString(got.CutModes, "CONSENSUS") {
		t.Fatalf("cut modes=%v, want CONSENSUS", got.CutModes)
	}
	if len(got.EffectiveModes) != 1 || got.EffectiveModes[0] != "OBSERVE-BURST" {
		t.Fatalf("effective modes=%v, want OBSERVE-BURST only", got.EffectiveModes)
	}
	if got.EffectivePolicy.Alerts != 1 || got.EffectivePolicy.ROI <= 0 {
		t.Fatalf("effective policy=%+v, want one positive OBSERVE-BURST alert", got.EffectivePolicy)
	}
}

func containsString(vals []string, want string) bool {
	for _, v := range vals {
		if v == want {
			return true
		}
	}
	return false
}
