package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func TestBuildScoutlistRequiresRecentProfitLeaderboard(t *testing.T) {
	base := walletdiscover.WalletScore{
		Address:         "0x0000000000000000000000000000000000000001",
		Tier:            "C",
		SmartMoneyScore: 100,
		BotScore:        40,
		Edge:            80,
		Sources:         map[string]int{"leaderboard_profit_all": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:      100,
			AvgTradeNotional: 1000,
		},
	}
	params := scoutParams{
		MaxWallets:      5,
		MaxBot:          45,
		MinSmart:        80,
		MinEdge:         60,
		MinLargeTrades:  50,
		MinAvgNotional:  500,
		ExcludeRiskTags: map[string]struct{}{},
	}

	if got := buildScoutlist([]walletdiscover.WalletScore{base}, nil, params); len(got) != 0 {
		t.Fatalf("profit_all-only wallet should not enter scout list: got %d", len(got))
	}

	base.Sources = map[string]int{"leaderboard_profit_7d": 1}
	if got := buildScoutlist([]walletdiscover.WalletScore{base}, nil, params); len(got) != 1 {
		t.Fatalf("recent profit wallet should enter scout list: got %d", len(got))
	}
}

func TestDirectTapeObserveRows_IncludesRejectedHugeWalletsOnly(t *testing.T) {
	huge := walletdiscover.WalletScore{
		Address:         "0x1111111111111111111111111111111111111111",
		Tier:            "BOT",
		BotScore:        61.5,
		RiskFlags:       []string{"bot_like_flow", "burst_trading"},
		Stats:           walletdiscover.WalletStats{CopyROI: -0.5, CopyClosedTrades: 2},
		SmartMoneyScore: 100,
	}
	pushed := walletdiscover.WalletScore{
		Address: "0x2222222222222222222222222222222222222222",
		Tier:    "B",
	}
	inputs := []tapeInput{
		{Address: huge.Address, Tier: "BOT", Buys: 3, BuyNotional: 122538, MaxBuy: 45177, Bot: 61.5},
		{Address: pushed.Address, Tier: "B", Buys: 1, BuyNotional: 6000, MaxBuy: 6000, Bot: 25},
		{Address: "0x3333333333333333333333333333333333333333", Tier: "C", Buys: 1, BuyNotional: 1000, MaxBuy: 1000, Bot: 10},
	}

	rows := directTapeObserveRows([]walletdiscover.WalletScore{huge, pushed}, []walletdiscover.WalletScore{pushed}, inputs, 5000, 8000, 45)
	if len(rows) != 1 {
		t.Fatalf("rows len=%d, want 1: %#v", len(rows), rows)
	}
	if !strings.Contains(rows[0], strings.ToLower(huge.Address)) || !strings.Contains(rows[0], "status=reject-bot") {
		t.Fatalf("unexpected observe row: %s", rows[0])
	}
}

func TestFilterExcludedTapeInputs_RemovesReviewNoiseBeforeObserve(t *testing.T) {
	reviewNoise := "0x1111111111111111111111111111111111111111"
	kept := "0x2222222222222222222222222222222222222222"
	inputs := []tapeInput{
		{Address: reviewNoise, Tier: "B", Buys: 10, BuyNotional: 45119, MaxBuy: 14000, Bot: 23.7},
		{Address: kept, Tier: "B", Buys: 1, BuyNotional: 6000, MaxBuy: 6000, Bot: 25},
	}

	filtered := filterExcludedTapeInputs(inputs, map[string]struct{}{reviewNoise: {}})
	rows := directTapeObserveRows(nil, nil, filtered, 5000, 8000, 45)
	if len(rows) != 1 {
		t.Fatalf("rows len=%d, want 1: %#v", len(rows), rows)
	}
	if strings.Contains(rows[0], reviewNoise) {
		t.Fatalf("review-noise wallet leaked into observe rows: %s", rows[0])
	}
	if !strings.Contains(rows[0], kept) {
		t.Fatalf("kept wallet missing from observe rows: %s", rows[0])
	}
}

func TestDirectTapeObserveRows_IncludesScoredLowBotBurstWallet(t *testing.T) {
	burst := walletdiscover.WalletScore{
		Address:         "0x1111111111111111111111111111111111111111",
		Tier:            "B",
		BotScore:        34.3,
		SmartMoneyScore: 80,
	}
	inputs := []tapeInput{
		{Address: burst.Address, Tier: "B", Buys: 2, BuyNotional: 9084, MaxBuy: 4618, Bot: 34.3},
		{Address: "0x2222222222222222222222222222222222222222", Tier: "D", Buys: 3, BuyNotional: 12000, MaxBuy: 4500, Bot: 20},
		{Address: "0x3333333333333333333333333333333333333333", Tier: "B", Buys: 2, BuyNotional: 9000, MaxBuy: 4500, Bot: 55},
	}

	rows := directTapeObserveRows([]walletdiscover.WalletScore{burst}, nil, inputs, 5000, 8000, 45)
	if len(rows) != 1 {
		t.Fatalf("rows len=%d, want 1: %#v", len(rows), rows)
	}
	if !strings.Contains(rows[0], strings.ToLower(burst.Address)) || !strings.Contains(rows[0], "status=watch-burst") {
		t.Fatalf("unexpected burst observe row: %s", rows[0])
	}
}

func TestBuildTapeHotlistRequiresScoredPositiveTargetCopy(t *testing.T) {
	params := tapeParams{
		MaxWallets:       10,
		MaxBot:           45,
		MinDirectMaxBuy:  5000,
		MinScoredMaxBuy:  2500,
		MinSmart:         70,
		MinTargetCopyT:   2,
		MinTargetCopyROI: 25,
		ExcludeRiskTags:  map[string]struct{}{"bot_like_flow": {}, "fixed_price": {}, "negative_copy_sim": {}},
	}
	negative := walletdiscover.WalletScore{
		Address:         "0x1111111111111111111111111111111111111111",
		Tier:            "B",
		SmartMoneyScore: 95,
		BotScore:        5,
		Stats:           walletdiscover.WalletStats{TargetCopyClosed: 1, TargetCopyROI: -60},
	}
	positive := walletdiscover.WalletScore{
		Address:         "0x2222222222222222222222222222222222222222",
		Tier:            "C",
		SmartMoneyScore: 100,
		BotScore:        30,
		Stats:           walletdiscover.WalletStats{TargetCopyClosed: 4, TargetCopyROI: 60},
	}
	inputs := []tapeInput{
		{Address: "0x3333333333333333333333333333333333333333", Tier: "TAPE", MaxBuy: 7000, BuyNotional: 7000, Buys: 1},
		{Address: negative.Address, Tier: "B", MaxBuy: 10000, BuyNotional: 10000, Buys: 1},
		{Address: positive.Address, Tier: "C", MaxBuy: 3000, BuyNotional: 3000, Buys: 1},
	}

	got := buildTapeHotlist([]walletdiscover.WalletScore{negative, positive}, inputs, nil, params)
	if len(got) != 1 {
		t.Fatalf("hotlist len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Address != positive.Address {
		t.Fatalf("hotlist wallet=%s, want positive wallet", got[0].Address)
	}
}

func TestBuildTapeHotlistBlocksNegativeMeasuredEdge(t *testing.T) {
	addr := "0x2222222222222222222222222222222222222222"
	params := tapeParams{
		MaxWallets:       10,
		MaxBot:           45,
		MinDirectMaxBuy:  5000,
		MinScoredMaxBuy:  2500,
		MinSmart:         70,
		MinTargetCopyT:   2,
		MinTargetCopyROI: 25,
		ExcludeRiskTags:  map[string]struct{}{},
		EdgeBlocks:       map[string]edgeBlock{addr: {Reason: "1h edge -18.95pp over 1 samples"}},
	}
	score := walletdiscover.WalletScore{
		Address:         addr,
		Tier:            "C",
		SmartMoneyScore: 100,
		BotScore:        30,
		Stats:           walletdiscover.WalletStats{TargetCopyClosed: 4, TargetCopyROI: 60},
	}
	inputs := []tapeInput{{Address: addr, Tier: "C", MaxBuy: 5000, BuyNotional: 5000, Buys: 1}}

	if got := buildTapeHotlist([]walletdiscover.WalletScore{score}, inputs, nil, params); len(got) != 0 {
		t.Fatalf("negative measured edge wallet entered hotlist: %#v", got)
	}
	if got := buildTapeProbation([]walletdiscover.WalletScore{score}, inputs, nil, params); len(got) != 0 {
		t.Fatalf("negative measured edge wallet entered probation: %#v", got)
	}
}

func TestLoadEdgeBlocks_BlocksOneHourNegativeEdge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	data := `{"wallet":"0x2222222222222222222222222222222222222222","horizon_sec":3600,"delta_pp":-18.95}` + "\n" +
		`{"wallet":"0x3333333333333333333333333333333333333333","horizon_sec":900,"delta_pp":-0.5}` + "\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	blocks, err := loadEdgeBlocks(path, 2, -1, 1, -5)
	if err != nil {
		t.Fatal(err)
	}
	if reason := blocks["0x2222222222222222222222222222222222222222"].Reason; !strings.Contains(reason, "1h edge -18.95pp") {
		t.Fatalf("block reason=%q, want 1h negative edge", reason)
	}
	if _, ok := blocks["0x3333333333333333333333333333333333333333"]; ok {
		t.Fatalf("15m single mild negative sample should not block: %#v", blocks)
	}
}

func TestBuildTapeEdgeHotRequiresPositiveMeasuredEdge(t *testing.T) {
	hot := "0x1111111111111111111111111111111111111111"
	cold := "0x2222222222222222222222222222222222222222"
	edgeOnly := "0x3333333333333333333333333333333333333333"
	params := tapeEdgeHotParams{
		MaxWallets:   10,
		MaxBot:       45,
		MinMaxBuy:    500,
		MinSamples:   2,
		MinAvgPP:     2,
		MinWinRate:   60,
		Min5mAvgPP:   0.5,
		Min15mAvgPP:  0,
		Max1hNegPP:   -5,
		EdgeProfiles: map[string]edgeProfile{},
		EdgeBlocks:   map[string]edgeBlock{},
		ExcludedRisks: map[string]struct{}{
			"bot_like_flow": {},
			"fixed_price":   {},
		},
	}
	params.EdgeProfiles[hot] = edgeProfile{TotalSamples: 3, TotalWins: 3, AvgPP: 8, WinRate: 100, Samples5m: 1, Avg5mPP: 3, Samples15m: 1, Avg15mPP: 5, Samples1h: 1, Avg1hPP: 16}
	params.EdgeProfiles[cold] = edgeProfile{TotalSamples: 3, TotalWins: 3, AvgPP: 8, WinRate: 100, Samples5m: 1, Avg5mPP: 3, Samples15m: 1, Avg15mPP: 5, Samples1h: 1, Avg1hPP: -20}
	params.EdgeProfiles[edgeOnly] = edgeProfile{TotalSamples: 3, TotalWins: 3, AvgPP: 10, WinRate: 100, MaxNotional: 2200, TapeAction: true, Samples5m: 1, Avg5mPP: 4, Samples15m: 1, Avg15mPP: 7, Samples1h: 1, Avg1hPP: 19}
	scores := []walletdiscover.WalletScore{
		{Address: hot, Tier: "C", SmartMoneyScore: 100, BotScore: 34, Sources: map[string]int{"sports_tape": 1}, Stats: walletdiscover.WalletStats{TargetTradeRatio: 0.25, TargetCopyROI: 50, TargetCopyClosed: 3, AvgTradeNotional: 1000}},
		{Address: cold, Tier: "C", SmartMoneyScore: 100, BotScore: 20, Sources: map[string]int{"sports_tape": 1}, Stats: walletdiscover.WalletStats{TargetCopyROI: 50, TargetCopyClosed: 3, AvgTradeNotional: 1000}},
	}
	inputs := []tapeInput{
		{Address: hot, Tier: "C", MaxBuy: 1200, BuyNotional: 1200, Buys: 1},
		{Address: cold, Tier: "C", MaxBuy: 1500, BuyNotional: 1500, Buys: 1},
	}

	got := buildTapeEdgeHot(scores, inputs, nil, params)
	if len(got) != 2 {
		t.Fatalf("edge-hot wallets=%#v, want hot and edge-only wallets", got)
	}
	seen := map[string]walletdiscover.WalletScore{}
	for _, s := range got {
		seen[s.Address] = s
	}
	if _, ok := seen[hot]; !ok {
		t.Fatalf("edge-hot wallets=%#v, missing scored hot wallet", got)
	}
	if s, ok := seen[edgeOnly]; !ok || s.Stats.AvgTradeNotional != 2200 {
		t.Fatalf("edge-hot wallets=%#v, missing edge-only wallet with snapshot notional", got)
	}
}

func TestWriteTapeEdgeHotFile_IncludesEdgeMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edgehot.txt")
	addr := "0x1111111111111111111111111111111111111111"
	wallets := []walletdiscover.WalletScore{{Address: addr, Tier: "C", SmartMoneyScore: 100, BotScore: 34}}
	profiles := map[string]edgeProfile{addr: {TotalSamples: 3, WinRate: 100, AvgPP: 8.25, Avg5mPP: 2.5, Avg15mPP: 4.5, MaxNotional: 1200, Reason: "edge-hot 100.0% win avg +8.25pp over 3 samples"}}
	inputs := []tapeInput{}

	if err := writeTapeEdgeHotFile(path, wallets, profiles, inputs); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{addr, "list=tape_edgehot", "maxBuy=$1200", "edgeN=3", "edgeAvgPP=+8.25"} {
		if !strings.Contains(got, want) {
			t.Fatalf("edge-hot file missing %q:\n%s", want, got)
		}
	}
}

func TestFilterEdgeBlockedWallets_RemovesSelectedPushWallet(t *testing.T) {
	blocked := walletdiscover.WalletScore{Address: "0x1111111111111111111111111111111111111111", Tier: "B"}
	kept := walletdiscover.WalletScore{Address: "0x2222222222222222222222222222222222222222", Tier: "B"}
	blocks := map[string]edgeBlock{
		strings.ToLower(blocked.Address): {Reason: "1h edge -29.95pp over 2 samples"},
	}

	got := filterEdgeBlockedWallets([]walletdiscover.WalletScore{blocked, kept}, blocks)
	if len(got) != 1 || got[0].Address != kept.Address {
		t.Fatalf("filtered wallets=%#v, want only kept wallet", got)
	}
	blockedRows := edgeBlockedWallets([]walletdiscover.WalletScore{blocked, kept}, blocks)
	if len(blockedRows) != 1 || blockedRows[0].Address != blocked.Address {
		t.Fatalf("blocked wallets=%#v, want only blocked wallet", blockedRows)
	}
}

func TestWritePushWalletsFile_AllowsScoutOmission(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "push.txt")
	core := walletdiscover.WalletScore{Address: "0x1111111111111111111111111111111111111111", Tier: "A"}
	scout := walletdiscover.WalletScore{Address: "0x2222222222222222222222222222222222222222", Tier: "C"}

	if err := writePushWalletsFile(path, []walletdiscover.WalletScore{core}, nil, nil, nil, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), strings.ToLower(core.Address)) {
		t.Fatalf("core wallet missing from push file:\n%s", string(b))
	}
	if strings.Contains(string(b), strings.ToLower(scout.Address)) {
		t.Fatalf("omitted scout wallet leaked into push file:\n%s", string(b))
	}

	if err := writePushWalletsFile(path, nil, nil, nil, []walletdiscover.WalletScore{scout}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), strings.ToLower(scout.Address)) || !strings.Contains(string(b), "list=scout") {
		t.Fatalf("explicit scout push missing from push file:\n%s", string(b))
	}
}

func TestBuildTapeProbationAllowsSoftRiskObservationOnly(t *testing.T) {
	params := tapeParams{
		MaxWallets:       10,
		MaxBot:           45,
		MinScoredMaxBuy:  2500,
		MinSmart:         70,
		MinTargetCopyT:   2,
		MinTargetCopyROI: 25,
	}
	softRisk := walletdiscover.WalletScore{
		Address:         "0x1111111111111111111111111111111111111111",
		Tier:            "C",
		SmartMoneyScore: 100,
		BotScore:        43.6,
		RiskFlags:       []string{"fixed_price", "opposite_side_same_market"},
		Stats:           walletdiscover.WalletStats{TargetCopyClosed: 10, TargetCopyROI: 75.5, TargetLargeTrades: 28, AvgTradeNotional: 500},
	}
	hardRisk := walletdiscover.WalletScore{
		Address:         "0x2222222222222222222222222222222222222222",
		Tier:            "C",
		SmartMoneyScore: 100,
		BotScore:        40,
		RiskFlags:       []string{"bot_like_flow"},
		Stats:           walletdiscover.WalletStats{TargetCopyClosed: 10, TargetCopyROI: 80, TargetLargeTrades: 28, AvgTradeNotional: 500},
	}
	inputs := []tapeInput{
		{Address: softRisk.Address, Tier: "C", MaxBuy: 3000, BuyNotional: 5000, Buys: 2},
		{Address: hardRisk.Address, Tier: "C", MaxBuy: 3000, BuyNotional: 5000, Buys: 2},
	}

	got := buildTapeProbation([]walletdiscover.WalletScore{softRisk, hardRisk}, inputs, nil, params)
	if len(got) != 1 {
		t.Fatalf("probation len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Address != softRisk.Address {
		t.Fatalf("probation wallet=%s, want soft risk wallet", got[0].Address)
	}
}
