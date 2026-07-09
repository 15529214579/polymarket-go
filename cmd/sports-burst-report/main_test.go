package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildBurstSignals_UsesLatestWindowAndVWAP(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	asset := "asset-1"
	trades := []tapeTrade{
		{Time: now.Add(-30 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 10000, Price: 0.19, Size: 52631.578947, Asset: asset, KnownList: "flow"},
		{Time: now.Add(-10 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 1200, Price: 0.20, Size: 6000, Asset: asset, KnownList: "flow"},
		{Time: now.Add(-5 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3000, Price: 0.21, Size: 14285.714286, Asset: asset, KnownList: "flow"},
		{Time: now.Add(-2 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3400, Price: 0.2125, Size: 16000, Asset: asset, KnownList: "flow", Market: "Will Argentina win?"},
	}

	got := buildBurstSignals(trades, now, time.Hour, 15*time.Minute, 5000, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("bursts len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Trades != 3 || got[0].TotalNotional != 7600 {
		t.Fatalf("trades/total=%d/%.0f, want 3/7600", got[0].Trades, got[0].TotalNotional)
	}
	if got[0].VWAP < 0.2094 || got[0].VWAP > 0.2096 {
		t.Fatalf("vwap=%.6f, want about 0.2095", got[0].VWAP)
	}
	if got[0].Mode != "FLOW-SCOUT" {
		t.Fatalf("mode=%s, want FLOW-SCOUT", got[0].Mode)
	}
	if got[0].Scope != "wallet" || got[0].Wallets != 1 {
		t.Fatalf("scope/wallets=%s/%d, want wallet/1", got[0].Scope, got[0].Wallets)
	}
}

func TestBuildConsensusBurstSignals_FindsCrossWalletSameAsset(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	trades := []tapeTrade{
		{Time: now.Add(-8 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 9000, Price: 0.30, Size: 30000, Asset: "asset-1", Tier: "B", Bot: 20, Market: "Will France win?"},
		{Time: now.Add(-4 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "C", Bot: 30, Market: "Will France win?"},
		{Time: now.Add(-3 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-2", Tier: "BOT", Bot: 80, Market: "Will Spain win?"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x4444444444444444444444444444444444444444", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-2", KnownList: "review_noise", Tier: "D", Bot: 20, Market: "Will Spain win?"},
	}

	got := buildConsensusBurstSignals(trades, now, time.Hour, 15*time.Minute, 15000, 2, 1000, 2, 45, nil)
	if len(got) != 1 {
		t.Fatalf("consensus bursts len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Mode != "CONSENSUS" || got[0].Scope != "consensus" {
		t.Fatalf("mode/scope=%s/%s, want CONSENSUS/consensus", got[0].Mode, got[0].Scope)
	}
	if got[0].Wallets != 2 || got[0].Trades != 2 || got[0].TotalNotional != 17000 {
		t.Fatalf("wallets/trades/total=%d/%d/%.0f, want 2/2/17000", got[0].Wallets, got[0].Trades, got[0].TotalNotional)
	}
	if got[0].Tier != "B" {
		t.Fatalf("tier=%s, want strongest B", got[0].Tier)
	}
	if len(got[0].Participants) != 2 || got[0].Participants[0] != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("participants=%v, want sorted wallet list", got[0].Participants)
	}
	if got[0].ParticipantMeta["0x1111111111111111111111111111111111111111"].Tier != "B" || got[0].ParticipantMeta["0x2222222222222222222222222222222222222222"].Bot != 30 {
		t.Fatalf("participant meta=%+v, want per-wallet tier/bot", got[0].ParticipantMeta)
	}
}

func TestBuildConsensusBurstSignals_ExcludesRejectedWalletStatus(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	rejected := "0x1111111111111111111111111111111111111111"
	good := "0x2222222222222222222222222222222222222222"
	other := "0x3333333333333333333333333333333333333333"
	trades := []tapeTrade{
		{Time: now.Add(-8 * time.Minute), Wallet: rejected, Side: "BUY", Notional: 9000, Price: 0.30, Size: 30000, Asset: "asset-1", KnownList: "tape_observe", Tier: "D", Bot: 58, Market: "basketball"},
		{Time: now.Add(-7 * time.Minute), Wallet: good, Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "B", Bot: 20, Market: "basketball"},
		{Time: now.Add(-6 * time.Minute), Wallet: other, Side: "BUY", Notional: 7000, Price: 0.31, Size: 22580.645161, Asset: "asset-1", Tier: "B", Bot: 21, Market: "basketball"},
	}
	applyWalletStatuses(trades, map[string]walletStatus{
		rejected: {List: "review_noise", Tier: "D", Bot: 58},
	})

	got := buildConsensusBurstSignals(trades, now, time.Hour, 15*time.Minute, 10000, 2, 1000, 2, 60, nil)
	if len(got) != 1 {
		t.Fatalf("consensus bursts len=%d, want one good-only burst: %#v", len(got), got)
	}
	if got[0].Wallets != 2 || got[0].TotalNotional != 15000 {
		t.Fatalf("wallets/total=%d/%.0f, want rejected excluded and total 15000", got[0].Wallets, got[0].TotalNotional)
	}
	for _, participant := range got[0].Participants {
		if participant == rejected {
			t.Fatalf("rejected wallet included in participants: %#v", got[0].Participants)
		}
	}
}

func TestBuildConsensusBurstSignals_ExcludesNegativeEdgeWallet(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	blocked := "0x1111111111111111111111111111111111111111"
	good := "0x2222222222222222222222222222222222222222"
	other := "0x3333333333333333333333333333333333333333"
	trades := []tapeTrade{
		{Time: now.Add(-8 * time.Minute), Wallet: blocked, Side: "BUY", Notional: 9000, Price: 0.30, Size: 30000, Asset: "asset-1", Tier: "B", Bot: 20, Market: "basketball"},
		{Time: now.Add(-7 * time.Minute), Wallet: good, Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "B", Bot: 20, Market: "basketball"},
		{Time: now.Add(-6 * time.Minute), Wallet: other, Side: "BUY", Notional: 7000, Price: 0.31, Size: 22580.645161, Asset: "asset-1", Tier: "B", Bot: 21, Market: "basketball"},
	}

	got := buildConsensusBurstSignals(trades, now, time.Hour, 15*time.Minute, 10000, 2, 1000, 2, 60, map[string]string{
		blocked: "1h edge -18.95pp over 1 samples",
	})
	if len(got) != 1 {
		t.Fatalf("consensus bursts len=%d, want one good-only burst: %#v", len(got), got)
	}
	if got[0].Wallets != 2 || got[0].TotalNotional != 15000 {
		t.Fatalf("wallets/total=%d/%.0f, want blocked excluded and total 15000", got[0].Wallets, got[0].TotalNotional)
	}
	for _, participant := range got[0].Participants {
		if participant == blocked {
			t.Fatalf("negative-edge wallet included in participants: %#v", got[0].Participants)
		}
	}
}

func TestLoadWalletStatuses_RejectStatusBecomesReviewNoise(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "observe.txt")
	addr := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(path, []byte(addr+" # list=tape_observe status=reject-bot tier=D bot=58.5\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletStatuses(path)
	if err != nil {
		t.Fatal(err)
	}
	if got[addr].List != "review_noise" || got[addr].Tier != "D" || got[addr].Bot < 58 {
		t.Fatalf("status=%+v, want review_noise D bot", got[addr])
	}
}

func TestEvaluateBursts_ComputesFixedStakeROI(t *testing.T) {
	signals := []burstSignal{{
		Wallet: "0x1111111111111111111111111111111111111111",
		Asset:  "asset-1",
		VWAP:   0.20,
	}}
	results := evaluateBursts(context.Background(), signals, 10, func(_ context.Context, sig burstSignal) (float64, bool, error) {
		if sig.Asset != "asset-1" {
			t.Fatalf("asset=%s, want asset-1", sig.Asset)
		}
		return 0.22, true, nil
	})
	if len(results) != 1 || !results[0].Marked {
		t.Fatalf("results=%+v, want one marked result", results)
	}
	if results[0].ReturnPC < 9.9 || results[0].ReturnPC > 10.1 {
		t.Fatalf("roi=%.2f, want 10%%", results[0].ReturnPC)
	}
}

func TestEvaluateBursts_AllowsZeroSettlementMark(t *testing.T) {
	signals := []burstSignal{{
		Wallet: "multi:2",
		Asset:  "asset-1",
		VWAP:   0.25,
	}}
	results := evaluateBursts(context.Background(), signals, 10, func(_ context.Context, sig burstSignal) (float64, bool, error) {
		return 0, true, nil
	})
	if len(results) != 1 || !results[0].Marked {
		t.Fatalf("results=%+v, want zero mark treated as marked loss", results)
	}
	if results[0].ReturnPC > -99.9 || results[0].ReturnPC < -100.1 {
		t.Fatalf("roi=%.2f, want -100%%", results[0].ReturnPC)
	}
}

func TestWriteReport_IncludesGatesAndRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	results := []burstResult{
		{
			Signal: burstSignal{
				Wallet:        "0x1111111111111111111111111111111111111111",
				Scope:         "wallet",
				Participants:  []string{"0x1111111111111111111111111111111111111111"},
				Mode:          "FLOW-SCOUT",
				Wallets:       1,
				Trades:        3,
				TotalNotional: 7600,
				VWAP:          0.20,
				Market:        "Will Argentina win the 2026 FIFA World Cup?",
				LastTime:      now.Add(-2 * time.Minute),
			},
			StakeUSD: 10,
			Units:    50,
			Mid:      0.22,
			Marked:   true,
			PnLUSD:   1,
			ReturnPC: 10,
		},
		{
			Signal: burstSignal{
				Wallet:        "multi:2",
				Scope:         "consensus",
				Participants:  []string{"0x2222222222222222222222222222222222222222", "0x3333333333333333333333333333333333333333"},
				Mode:          "CONSENSUS",
				Wallets:       2,
				Trades:        2,
				TotalNotional: 15000,
				VWAP:          0.30,
				Market:        "Dota 2: Team A vs Team B",
				LastTime:      now.Add(-3 * time.Minute),
			},
			StakeUSD: 10,
			Units:    33.333333,
			Mid:      0.36,
			Marked:   true,
			PnLUSD:   2,
			ReturnPC: 20,
		},
	}
	events := []consensusEvent{{
		Key:           "consensus|asset-1|1|2|15000",
		LastTime:      now.Add(-2 * time.Minute),
		Category:      "esports",
		Market:        "Dota 2: Team A vs Team B",
		Wallets:       2,
		Trades:        2,
		TotalNotional: 15000,
		VWAP:          0.30,
		Marked:        true,
		Mid:           0.36,
		PnLUSD:        2,
		ReturnPC:      20,
	}}
	watchResults := []burstResult{{
		Signal: burstSignal{
			Mode:          "CONSENSUS-WATCH",
			Scope:         "consensus-watch",
			Participants:  []string{"0x5555555555555555555555555555555555555555", "0x6666666666666666666666666666666666666666"},
			Wallets:       2,
			Trades:        2,
			TotalNotional: 6500,
			VWAP:          0.40,
			Market:        "LoL: Team C vs Team D",
			LastTime:      now.Add(-4 * time.Minute),
		},
		StakeUSD: 10,
		Units:    25,
		Mid:      0.44,
		Marked:   true,
		PnLUSD:   1,
		ReturnPC: 10,
	}}
	watchEvents := []consensusEvent{{
		Key:           "consensus-watch|asset-watch|1|2|6500",
		Mode:          "CONSENSUS-WATCH",
		LastTime:      now.Add(-4 * time.Minute),
		Market:        "LoL: Team C vs Team D",
		Wallets:       2,
		Trades:        2,
		TotalNotional: 6500,
		VWAP:          0.40,
		Marked:        true,
		Mid:           0.44,
		PnLUSD:        1,
		ReturnPC:      10,
		Participants:  []string{"0x5555555555555555555555555555555555555555", "0x6666666666666666666666666666666666666666"},
		ParticipantMeta: map[string]participantMeta{
			"0x5555555555555555555555555555555555555555": {Tier: "C", Bot: 33},
		},
	}}
	researchRows := []participantSummary{{
		Wallet:        "0x4444444444444444444444444444444444444444",
		Tier:          "D",
		Bot:           55,
		Signals:       1,
		Marked:        1,
		Wins:          1,
		ReturnPC:      20,
		PnLUSD:        2,
		TotalNotional: 15000,
		AvgDeltaP:     6,
	}}
	if err := writeReport(path, "tape.jsonl", results, watchResults, watchEvents, events, researchRows, 10, 15*time.Minute, time.Hour, 5000, 2, 1000, 15000, 2, 45, 5000, 2, 1, now); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	report := string(b)
	for _, want := range []string{"# Sports Burst Performance", "Mark source: live CLOB midpoint", "Consensus rule", "Bursts: 2", "COLLECT_POSITIVE", "## Scope Gates", "## Consensus Participants", "0x2222...2222", "20.0%", "## Research-Only Consensus Watch", "Watch bursts: 1", "Durable watch events: 1", "### Durable Watch History", "LoL: Team C vs Team D", "### Durable Watch Participants", "0x5555...5555", "## Durable Consensus Event History", "Events: 1", "Marked samples still needed for promotion review: 4", "Dota 2: Team A vs Team B", "## Durable Consensus Research Wallets", "0x4444...4444", "## Recent Bursts", "FLOW-SCOUT", "0x1111...1111"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
}

func TestConsensusParticipantSummaries_RanksRepeatedPositiveWallets(t *testing.T) {
	results := []burstResult{
		{
			Signal:   burstSignal{Mode: "CONSENSUS", Participants: []string{"0x2222222222222222222222222222222222222222", "0x1111111111111111111111111111111111111111"}, TotalNotional: 10000, VWAP: 0.25},
			StakeUSD: 10,
			Mid:      0.30,
			Marked:   true,
			PnLUSD:   2,
		},
		{
			Signal:   burstSignal{Mode: "CONSENSUS", Participants: []string{"0x2222222222222222222222222222222222222222"}, TotalNotional: 12000, VWAP: 0.40},
			StakeUSD: 10,
			Mid:      0.44,
			Marked:   true,
			PnLUSD:   1,
		},
	}

	got := consensusParticipantSummaries(results)
	if len(got) != 2 {
		t.Fatalf("participants len=%d, want 2", len(got))
	}
	if got[0].Wallet != "0x2222222222222222222222222222222222222222" || got[0].Marked != 2 {
		t.Fatalf("top participant=%+v, want repeated 0x2222 marked=2", got[0])
	}
	if got[0].ReturnPC < 14.9 || got[0].ReturnPC > 15.1 {
		t.Fatalf("top ROI=%.2f, want 15%%", got[0].ReturnPC)
	}
}

func TestConsensusWatchSignals_ExcludesOfficialThreshold(t *testing.T) {
	signals := []burstSignal{
		{Mode: "CONSENSUS", Scope: "consensus", Asset: "asset-watch", TotalNotional: 7000},
		{Mode: "CONSENSUS", Scope: "consensus", Asset: "asset-official", TotalNotional: 12000},
		{Mode: "FLOW-SCOUT", Scope: "single", Asset: "asset-flow", TotalNotional: 7000},
	}
	got := consensusWatchSignals(signals, 5000, 10000)
	if len(got) != 1 {
		t.Fatalf("watch signals=%d, want one lower-threshold consensus", len(got))
	}
	if got[0].Asset != "asset-watch" || got[0].Mode != "CONSENSUS-WATCH" || got[0].Scope != "consensus-watch" {
		t.Fatalf("watch signal=%+v, want rewritten research-only consensus watch", got[0])
	}
	if disabled := consensusWatchSignals(signals, 0, 10000); len(disabled) != 0 {
		t.Fatalf("disabled watch signals=%d, want 0", len(disabled))
	}
}

func TestWriteConsensusWatchEventsFile_WritesResearchOnlyHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watch.jsonl")
	now := time.Date(2026, 7, 9, 5, 0, 0, 0, time.UTC)
	result := burstResult{
		Signal: burstSignal{
			Mode:          "CONSENSUS-WATCH",
			Asset:         "asset-watch",
			Market:        "Will Spain win?",
			Wallets:       2,
			Trades:        3,
			TotalNotional: 8397,
			VWAP:          0.40,
			FirstTime:     now.Add(-10 * time.Minute),
			LastTime:      now.Add(-5 * time.Minute),
		},
		StakeUSD: 10,
		Mid:      0.404,
		Marked:   true,
		PnLUSD:   0.1,
		ReturnPC: 1,
	}

	if err := writeConsensusWatchEventsFile(path, []burstResult{result}, now); err != nil {
		t.Fatal(err)
	}
	events, err := loadConsensusEvents(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want one watch event", len(events))
	}
	for key, ev := range events {
		if !strings.HasPrefix(key, "consensus-watch|") || ev.Mode != "CONSENSUS-WATCH" {
			t.Fatalf("key/event=%s/%+v, want consensus-watch research-only history", key, ev)
		}
	}
}

func TestConsensusEventParticipantSummaries_RanksWatchParticipants(t *testing.T) {
	wallet := "0x1111111111111111111111111111111111111111"
	rows := consensusEventParticipantSummaries([]consensusEvent{
		{
			Participants: []string{wallet, "0x2222222222222222222222222222222222222222"},
			ParticipantMeta: map[string]participantMeta{
				wallet: {Tier: "C", Bot: 31},
			},
			Marked:        true,
			Mid:           0.44,
			VWAP:          0.40,
			PnLUSD:        1,
			TotalNotional: 6500,
		},
		{
			Participants: []string{wallet},
			ParticipantMeta: map[string]participantMeta{
				wallet: {Tier: "B", Bot: 29},
			},
			Marked:        true,
			Mid:           0.50,
			VWAP:          0.45,
			PnLUSD:        1.1,
			TotalNotional: 7200,
		},
	})
	if len(rows) != 2 {
		t.Fatalf("rows=%d, want two participants", len(rows))
	}
	if rows[0].Wallet != wallet || rows[0].Signals != 2 || rows[0].Marked != 2 {
		t.Fatalf("top row=%+v, want repeated participant first", rows[0])
	}
	if rows[0].Tier != "B" || rows[0].ReturnPC < 10.4 || rows[0].ReturnPC > 10.6 {
		t.Fatalf("top row=%+v, want best tier and aggregated ROI", rows[0])
	}
}

func TestWriteConsensusParticipantsFile_WritesPositiveResearchWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	rows := []participantSummary{
		{Wallet: "0x1111111111111111111111111111111111111111", Tier: "B", Bot: 22.5, Signals: 1, Marked: 1, Wins: 1, StakeUSD: 10, PnLUSD: 2, ReturnPC: 20, TotalNotional: 15000, AvgDeltaP: 6},
		{Wallet: "0x2222222222222222222222222222222222222222", Signals: 1, Marked: 1, Wins: 0, StakeUSD: 10, PnLUSD: -1, ReturnPC: -10, TotalNotional: 12000, AvgDeltaP: -3},
		{Wallet: "0x3333333333333333333333333333333333333333", Signals: 1, Marked: 0, TotalNotional: 20000},
	}

	if err := writeConsensusParticipantsFile(path, rows, nil, nil, 1, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "0x1111111111111111111111111111111111111111 # list=consensus_research") {
		t.Fatalf("missing positive wallet:\n%s", got)
	}
	if !strings.Contains(got, "tier=B bot=22.5") {
		t.Fatalf("missing participant tier/bot:\n%s", got)
	}
	if strings.Contains(got, "0x222222") || strings.Contains(got, "0x333333") {
		t.Fatalf("wrote non-positive or unmarked wallets:\n%s", got)
	}
}

func TestWriteConsensusParticipantsFile_SkipsExcludedWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	rows := []participantSummary{
		{Wallet: "0x1111111111111111111111111111111111111111", Signals: 1, Marked: 1, Wins: 1, StakeUSD: 10, PnLUSD: 2, ReturnPC: 20, TotalNotional: 15000, AvgDeltaP: 6},
		{Wallet: "0x2222222222222222222222222222222222222222", Signals: 1, Marked: 1, Wins: 1, StakeUSD: 10, PnLUSD: 1, ReturnPC: 10, TotalNotional: 12000, AvgDeltaP: 3},
	}
	exclude := map[string]struct{}{
		"0x1111111111111111111111111111111111111111": {},
	}

	if err := writeConsensusParticipantsFile(path, rows, exclude, nil, 1, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if strings.Contains(got, "0x111111") {
		t.Fatalf("wrote excluded wallet:\n%s", got)
	}
	if !strings.Contains(got, "0x2222222222222222222222222222222222222222 # list=consensus_research") {
		t.Fatalf("missing non-excluded positive wallet:\n%s", got)
	}
}

func TestWriteConsensusParticipantsFile_PreservesExistingResearchWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	existing := "0x9999999999999999999999999999999999999999 # list=consensus_research signals=1 marked=1 win=100.0% roi=12.0% pnl=$+1.20 notional=$9000 avgDeltaPP=+3.00 source=sports_consensus\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writeConsensusParticipantsFile(path, nil, nil, nil, 1, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "0x9999999999999999999999999999999999999999 # list=consensus_research") {
		t.Fatalf("existing research wallet was not preserved:\n%s", got)
	}
}

func TestWriteConsensusParticipantsFile_BackfillsExistingResearchMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	wallet := "0x9999999999999999999999999999999999999999"
	existing := wallet + " # list=consensus_research signals=1 marked=1 win=100.0% roi=12.0% pnl=$+1.20 notional=$9000 avgDeltaPP=+3.00 source=sports_consensus\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	scoreMeta := map[string]participantMeta{
		wallet: {Tier: "B", Bot: 22.4},
	}

	if err := writeConsensusParticipantsFile(path, nil, nil, scoreMeta, 1, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "list=consensus_research tier=B bot=22.4 signals=1") {
		t.Fatalf("existing research wallet metadata was not backfilled:\n%s", got)
	}
}

func TestWriteConsensusParticipantsFile_UpdatesExistingResearchWallet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(path, []byte(wallet+" # list=consensus_research signals=1 marked=1 win=100.0% roi=12.0% pnl=$+1.20 notional=$9000 avgDeltaPP=+3.00 source=sports_consensus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rows := []participantSummary{
		{Wallet: wallet, Signals: 2, Marked: 2, Wins: 2, StakeUSD: 20, PnLUSD: 5, ReturnPC: 25, TotalNotional: 24000, AvgDeltaP: 8},
	}

	if err := writeConsensusParticipantsFile(path, rows, nil, nil, 1, 0); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "signals=2 marked=2") || !strings.Contains(got, "roi=25.0%") {
		t.Fatalf("existing research wallet was not updated:\n%s", got)
	}
	if strings.Count(got, wallet) != 1 {
		t.Fatalf("wallet duplicated after update:\n%s", got)
	}
}

func TestLoadConsensusParticipantRows_ParsesPersistedResearchWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	raw := "0x1111111111111111111111111111111111111111 # list=consensus_research tier=D bot=55.0 signals=1 marked=1 win=100.0% roi=32.6% pnl=$+3.26 notional=$10336 avgDeltaPP=+24.58 source=sports_consensus\n"
	raw += "0x2222222222222222222222222222222222222222 # list=other signals=9 marked=9 roi=99.0%\n"
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	rows, err := loadConsensusParticipantRows(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want one consensus research wallet", len(rows))
	}
	got := rows[0]
	if got.Wallet != "0x1111111111111111111111111111111111111111" || got.Tier != "D" || got.Bot != 55 || got.Marked != 1 || got.Wins != 1 {
		t.Fatalf("row=%+v, want parsed participant metadata", got)
	}
	if got.ReturnPC != 32.6 || got.PnLUSD != 3.26 || got.TotalNotional != 10336 || got.AvgDeltaP != 24.58 {
		t.Fatalf("row=%+v, want parsed participant performance", got)
	}
}

func TestWriteConsensusEventsFile_WritesAndUpdatesConsensusHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	now := time.Date(2026, 7, 9, 5, 0, 0, 0, time.UTC)
	result := burstResult{
		Signal: burstSignal{
			Mode:         "CONSENSUS",
			Participants: []string{"0x1111111111111111111111111111111111111111", "0x2222222222222222222222222222222222222222"},
			ParticipantMeta: map[string]participantMeta{
				"0x1111111111111111111111111111111111111111": {Tier: "B", Bot: 20},
				"0x2222222222222222222222222222222222222222": {Tier: "C", Bot: 30},
			},
			Asset:         "asset-1",
			Slug:          "dota-team-a-team-b",
			Category:      "esports",
			Market:        "Dota 2: Team A vs Team B",
			Outcome:       "Team A",
			Wallets:       2,
			Trades:        2,
			TotalNotional: 15000,
			VWAP:          0.30,
			FirstTime:     now.Add(-10 * time.Minute),
			LastTime:      now.Add(-5 * time.Minute),
		},
		StakeUSD: 10,
		Units:    33.333333,
		Mid:      0.36,
		Marked:   true,
		PnLUSD:   2,
		ReturnPC: 20,
	}

	if err := writeConsensusEventsFile(path, []burstResult{result}, now); err != nil {
		t.Fatal(err)
	}
	result.Mid = 0.39
	result.PnLUSD = 3
	result.ReturnPC = 30
	if err := writeConsensusEventsFile(path, []burstResult{result}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	events, err := loadConsensusEvents(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want one deduped event: %#v", len(events), events)
	}
	for _, ev := range events {
		if !ev.Marked || ev.ReturnPC != 30 || ev.Mid != 0.39 {
			t.Fatalf("event mark=%+v, want updated marked ROI", ev)
		}
		if len(ev.Participants) != 2 || ev.ParticipantMeta["0x1111111111111111111111111111111111111111"].Tier != "B" {
			t.Fatalf("event participants/meta=%+v/%+v", ev.Participants, ev.ParticipantMeta)
		}
	}
}

func TestConsensusAlertSignals_ImportsLegacyShadowConsensus(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "shadow.jsonl")
	raw := `{"key":"consensus|asset-1|1783542565|5|10336.1700|asset-1|multi:5","trade_time":"2026-07-09T04:29:25+08:00","mode":"CONSENSUS","wallet":"multi:5","known_list":"consensus","tier":"D","bot":56.61,"category":"esports","notional":10336.17,"price":0.7542159772660015,"outcome":"1win","market":"Dota 2: Virtus.pro vs 1win - Game 2 Winner","slug":"dota2-vp-1win-2026-07-08","asset":"asset-1","transaction":"consensus|asset-1|1783542565|5|10336.1700"}` + "\n"
	raw += `{"trade_time":"2026-07-09T04:29:25+08:00","mode":"OBSERVE-BURST","wallet":"0x1111111111111111111111111111111111111111","asset":"asset-2","notional":9000,"price":0.5}` + "\n"
	if err := os.WriteFile(logPath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	events, err := loadConsensusAlertEvents(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want one consensus event", len(events))
	}
	signals := consensusAlertSignals(events)
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want one", len(signals))
	}
	if signals[0].Wallets != 5 || signals[0].Trades != 5 || signals[0].Scope != "alert-log" {
		t.Fatalf("signal=%+v, want legacy multi:5 consensus import", signals[0])
	}

	outPath := filepath.Join(dir, "events.jsonl")
	results := []burstResult{{
		Signal:   signals[0],
		StakeUSD: 10,
		Units:    10 / signals[0].VWAP,
		Mid:      1,
		Marked:   true,
		PnLUSD:   3.26,
		ReturnPC: 32.6,
	}}
	if err := writeConsensusEventsFile(outPath, results, time.Date(2026, 7, 9, 6, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	stored, err := loadConsensusEvents(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != 1 {
		t.Fatalf("stored=%d, want one", len(stored))
	}
	for _, ev := range stored {
		if ev.Wallets != 5 || ev.Category != "esports" || ev.ReturnPC != 32.6 {
			t.Fatalf("event=%+v, want imported marked esports consensus", ev)
		}
		if ev.Participants == nil {
			t.Fatalf("event participants is nil, want normalized empty slice")
		}
	}
}

func TestMergeConsensusHistoryResults_DedupesCurrentConsensus(t *testing.T) {
	now := time.Date(2026, 7, 9, 5, 0, 0, 0, time.UTC)
	current := burstResult{Signal: burstSignal{
		Mode:          "CONSENSUS",
		Asset:         "asset-1",
		Wallets:       2,
		TotalNotional: 15000,
		LastTime:      now,
	}}
	old := burstResult{Signal: burstSignal{
		Mode:          "CONSENSUS",
		Asset:         "asset-2",
		Wallets:       3,
		TotalNotional: 22000,
		LastTime:      now.Add(-time.Hour),
	}}
	got := mergeConsensusHistoryResults([]burstResult{current}, []burstResult{current, old})
	if len(got) != 2 {
		t.Fatalf("merged len=%d, want current plus one old event: %#v", len(got), got)
	}
	if got[1].Signal.Asset != "asset-2" {
		t.Fatalf("second asset=%s, want old event asset-2", got[1].Signal.Asset)
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

func TestLoadWalletSet_OnlyExcludesBlockedStatusRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "observe.txt")
	blocked := "0x1111111111111111111111111111111111111111"
	watch := "0x2222222222222222222222222222222222222222"
	if err := os.WriteFile(path, []byte(
		blocked+" # list=tape_observe status=reject-bot tier=D\n"+
			watch+" # list=tape_observe status=watch tier=B\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletSet(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got[blocked]; !ok {
		t.Fatalf("missing blocked wallet in exclude set")
	}
	if _, ok := got[watch]; ok {
		t.Fatalf("watch wallet should not be excluded")
	}
}
