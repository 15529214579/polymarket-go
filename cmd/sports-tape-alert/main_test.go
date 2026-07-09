package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestQualifyingTrades_FiltersRecentUnsentBuyWhales(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		AllowedModes:       parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION,EDGE-HOT"),
		MinNotional:        5000,
		ObserveMinNotional: 25000,
		ObserveMaxBot:      45,
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 6000, Price: 0.4, Asset: "asset-1", Transaction: "tx-1", KnownList: "tape_probation"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "SELL", Notional: 9000, Price: 0.5, Asset: "asset-2", Transaction: "tx-2"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 4999, Price: 0.5, Asset: "asset-3", Transaction: "tx-3", KnownList: "tape_probation"},
		{Time: now.Add(-20 * time.Minute), Wallet: "0x4444444444444444444444444444444444444444", Side: "BUY", Notional: 8000, Price: 0.5, Asset: "asset-4", Transaction: "tx-4", KnownList: "tape_probation"},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 1 {
		t.Fatalf("qualifying len=%d, want 1", len(got))
	}
	if got[0].Wallet != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("wallet=%s, want first", got[0].Wallet)
	}

	st := sentState{Sent: map[string]string{}}
	markTrades(st, got, now)
	if again := qualifyingTrades(trades, st, now, policy, 10*time.Minute); len(again) != 0 {
		t.Fatalf("qualifying after mark len=%d, want 0", len(again))
	}
}

func TestLoadModePolicyBlocks_CutsFreshModeDecision(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	body := `{
  "generated_at": "` + now.Add(-time.Minute).Format(time.RFC3339) + `",
  "modes": [
    {"key": "CONSENSUS", "action": "CUT", "reason": "negative marked sample"},
    {"key": "OBSERVE-BURST", "action": "COLLECT_POSITIVE", "reason": "still collecting"},
    {"key": "UNKNOWN-FLOW", "action": "PROBATION", "reason": "severe drawdown"}
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadModePolicyBlocks(path, time.Hour, now)
	if err != nil {
		t.Fatal(err)
	}
	if got["CONSENSUS"] != "policy action CUT below COLLECT_POSITIVE" {
		t.Fatalf("CONSENSUS block=%q, want CUT below COLLECT_POSITIVE", got["CONSENSUS"])
	}
	if got["UNKNOWN-FLOW"] != "policy action PROBATION below COLLECT_POSITIVE" {
		t.Fatalf("UNKNOWN-FLOW block=%q, want PROBATION below COLLECT_POSITIVE", got["UNKNOWN-FLOW"])
	}
	if _, ok := got["OBSERVE-BURST"]; ok {
		t.Fatalf("OBSERVE-BURST should not be blocked: %v", got)
	}

	stale, err := loadModePolicyBlocks(path, time.Second, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(stale) != 0 {
		t.Fatalf("stale policy blocks=%v, want none", stale)
	}
}

func TestSuppressPolicyBlockedModes_RemovesCutConsensus(t *testing.T) {
	trades := []tapeTrade{
		{Wallet: "multi:3", Asset: "asset-1", Price: 0.5, KnownList: "consensus"},
		{Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-2", Price: 0.4, KnownList: "observe_burst"},
	}
	got := suppressPolicyBlockedModes(trades, alertPolicy{ModeBlocks: map[string]string{"CONSENSUS": "negative marked sample"}})
	if len(got) != 1 {
		t.Fatalf("filtered len=%d, want 1: %#v", len(got), got)
	}
	if alertMode(got[0]) != "OBSERVE-BURST" {
		t.Fatalf("remaining mode=%s, want OBSERVE-BURST", alertMode(got[0]))
	}
}

func TestPolicyBlockedModeTrades_ReturnsUnloggedBlockedTrades(t *testing.T) {
	consensus := tapeTrade{Wallet: "multi:3", Asset: "asset-1", Price: 0.5, KnownList: "consensus", Transaction: "consensus-key"}
	observeBurst := tapeTrade{Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-2", Price: 0.4, KnownList: "observe_burst", Transaction: "observe-key"}
	logged := map[string]struct{}{}

	got := policyBlockedModeTrades([]tapeTrade{consensus, observeBurst}, alertPolicy{ModeBlocks: map[string]string{"CONSENSUS": "negative marked sample"}}, logged)
	if len(got) != 1 || alertMode(got[0]) != "CONSENSUS" {
		t.Fatalf("blocked trades=%#v, want one CONSENSUS", got)
	}
	if _, ok := logged[tradeKey(consensus)]; !ok {
		t.Fatalf("blocked consensus key was not marked logged")
	}
	again := policyBlockedModeTrades([]tapeTrade{consensus, observeBurst}, alertPolicy{ModeBlocks: map[string]string{"CONSENSUS": "negative marked sample"}}, logged)
	if len(again) != 0 {
		t.Fatalf("blocked duplicate len=%d, want 0", len(again))
	}
}

func TestSuppressPolicyBlockedModes_RequiresMinimumPolicyAction(t *testing.T) {
	trades := []tapeTrade{
		{Wallet: "0x1111111111111111111111111111111111111111", Asset: "asset-1", Price: 0.5, KnownList: "observe_burst"},
		{Wallet: "0x2222222222222222222222222222222222222222", Asset: "asset-2", Price: 0.4, KnownList: "unknown_flow"},
		{Wallet: "0x3333333333333333333333333333333333333333", Asset: "asset-3", Price: 0.6, KnownList: "tape_edgehot"},
	}
	policy := alertPolicy{
		ModeActions: map[string]string{
			"OBSERVE-BURST": "COLLECT_POSITIVE",
			"UNKNOWN-FLOW":  "COLLECT",
		},
		ModeMinAction: "COLLECT_POSITIVE",
	}

	got := suppressPolicyBlockedModes(trades, policy)
	if len(got) != 1 || alertMode(got[0]) != "OBSERVE-BURST" {
		t.Fatalf("filtered=%#v, want only OBSERVE-BURST", got)
	}
	blocked := policyBlockedModeTrades(trades, policy, map[string]struct{}{})
	if len(blocked) != 2 {
		t.Fatalf("blocked len=%d, want COLLECT and missing-action modes: %#v", len(blocked), blocked)
	}
}

func TestQualifyingTrades_DefaultPolicyFiltersNoisyObserveAndReversal(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		AllowedModes:        parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION"),
		MinNotional:         5000,
		ObserveMinNotional:  25000,
		ObserveMaxBot:       45,
		ObserveRequireKnown: true,
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 7675, Price: 0.4, Asset: "asset-1", Transaction: "tx-1"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 26000, Price: 0.5, Asset: "asset-2", Transaction: "tx-2", Tier: "B", Bot: 20},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 7000, Price: 0.5, Asset: "asset-3", Transaction: "tx-3", KnownList: "tape_reversal"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x4444444444444444444444444444444444444444", Side: "BUY", Notional: 5000, Price: 0.5, Asset: "asset-4", Transaction: "tx-4", KnownList: "tape_candidate"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x5555555555555555555555555555555555555555", Side: "BUY", Notional: 50000, Price: 0.5, Asset: "asset-5", Transaction: "tx-5", Tier: "BOT", Bot: 62},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x6666666666666666666666666666666666666666", Side: "BUY", Notional: 55000, Price: 0.5, Asset: "asset-6", Transaction: "tx-6"},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 2 {
		t.Fatalf("qualifying len=%d, want raw huge observe + candidate", len(got))
	}
	if got[0].Wallet != "0x2222222222222222222222222222222222222222" || got[1].Wallet != "0x4444444444444444444444444444444444444444" {
		t.Fatalf("wallets=%s,%s want huge observe and candidate", got[0].Wallet, got[1].Wallet)
	}
}

func TestAlertDecision_BlocksUnscoredObserveByDefault(t *testing.T) {
	policy := alertPolicy{
		ObserveMinNotional:  10000,
		ObserveMaxBot:       35,
		ObserveRequireKnown: true,
	}
	tr := tapeTrade{
		Wallet:   "0x1111111111111111111111111111111111111111",
		Side:     "BUY",
		Notional: 55000,
		Price:    0.326,
		Asset:    "asset-1",
	}

	ok, reason := alertDecision(tr, policy)
	if ok || reason != "observe wallet unscored" {
		t.Fatalf("alertDecision=%v/%q, want unscored observe block", ok, reason)
	}
	tr.Tier = "B"
	tr.Bot = 20
	ok, reason = alertDecision(tr, policy)
	if !ok || reason != "huge observe non-bot" {
		t.Fatalf("alertDecision=%v/%q, want known observe allowed", ok, reason)
	}
}

func TestAlertDecision_AllowsHugeUnscoredInsiderObserve(t *testing.T) {
	policy := alertPolicy{
		ObserveMinNotional:  10000,
		ObserveMaxBot:       35,
		ObserveRequireKnown: true,
		ObserveMinTier:      "B",
		InsiderMinNotional:  25000,
		InsiderMaxBot:       35,
	}
	tr := tapeTrade{
		Wallet:   "0x1111111111111111111111111111111111111111",
		Side:     "BUY",
		Notional: 55000,
		Price:    0.326,
		Asset:    "asset-1",
	}

	ok, reason := alertDecision(tr, policy)
	if !ok || reason != "insider-scout huge whale" {
		t.Fatalf("alertDecision=%v/%q, want insider observe allowed", ok, reason)
	}
	tr.Notional = 12000
	ok, reason = alertDecision(tr, policy)
	if ok || reason != "observe wallet unscored" {
		t.Fatalf("small unscored alertDecision=%v/%q, want unscored block", ok, reason)
	}
	tr.Notional = 55000
	tr.Tier = "BOT"
	tr.Bot = 62
	ok, reason = alertDecision(tr, policy)
	if ok || reason != "BOT tier" {
		t.Fatalf("bot insider alertDecision=%v/%q, want BOT block", ok, reason)
	}
}

func TestAlertDecision_BlocksConsensusResearchSingleWalletInsider(t *testing.T) {
	policy := alertPolicy{
		ObserveMinNotional:  10000,
		ObserveMaxBot:       35,
		ObserveRequireKnown: true,
		ObserveMinTier:      "B",
		InsiderMinNotional:  25000,
		InsiderMaxBot:       35,
	}
	tr := tapeTrade{
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  55000,
		Price:     0.326,
		Asset:     "asset-1",
		KnownList: "consensus_research",
	}

	ok, reason := alertDecision(tr, policy)
	if ok || reason != "consensus research wallet requires consensus burst" {
		t.Fatalf("alertDecision=%v/%q, want consensus research single-wallet block", ok, reason)
	}
	tr.Tier = "D"
	tr.Bot = 20
	ok, reason = alertDecision(tr, policy)
	if ok || reason != "consensus research wallet requires consensus burst" {
		t.Fatalf("alertDecision=%v/%q, want consensus research single-wallet block", ok, reason)
	}
	tr.Tier = "B"
	ok, reason = alertDecision(tr, policy)
	if ok || reason != "consensus research wallet requires consensus burst" {
		t.Fatalf("alertDecision=%v/%q, want consensus research single-wallet block", ok, reason)
	}
}

func TestQualifyingTrades_TagsHugeUnscoredInsiderMode(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		ObserveRequireKnown: true,
		ObserveMinTier:      "B",
		InsiderMinNotional:  25000,
		InsiderMaxBot:       35,
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 55000, Price: 0.326, Asset: "asset-1", Transaction: "tx-1"},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 1 {
		t.Fatalf("qualifying len=%d, want 1", len(got))
	}
	if got[0].KnownList != "insider_scout" || alertMode(got[0]) != "INSIDER-SCOUT" {
		t.Fatalf("known/mode=%s/%s, want insider_scout/INSIDER-SCOUT", got[0].KnownList, alertMode(got[0]))
	}
}

func TestQualifyingTrades_PositionCooldownDedupesSameWalletAsset(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	policy := alertPolicy{
		AllowedModes:      parseModeSet("FLOW-SCOUT"),
		MinNotional:       5000,
		PositionCooldown:  30 * time.Minute,
		RepeatMinNotional: 25000,
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 6000, Price: 0.20, Asset: "asset-1", Transaction: "tx-1", KnownList: "flow"},
		{Time: now.Add(-2 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 14000, Price: 0.21, Asset: "asset-1", Transaction: "tx-2", KnownList: "flow"},
		{Time: now.Add(-1 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 26000, Price: 0.22, Asset: "asset-1", Transaction: "tx-3", KnownList: "flow"},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 2 {
		t.Fatalf("qualifying len=%d, want 2: %#v", len(got), got)
	}
	if got[0].Transaction != "tx-2" {
		t.Fatalf("first tx=%s, want largest same-time tx-2", got[0].Transaction)
	}
	if got[1].Transaction != "tx-3" {
		t.Fatalf("second tx=%s, want repeat-min bypass tx-3", got[1].Transaction)
	}
}

func TestQualifyingTrades_PositionCooldownUsesPersistedState(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	oldTrade := tapeTrade{Time: now.Add(-10 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 8000, Price: 0.20, Asset: "asset-1", Transaction: "tx-old", KnownList: "flow"}
	st := sentState{Sent: map[string]string{}, Details: map[string]sentDetail{}}
	markTrades(st, []tapeTrade{oldTrade}, now.Add(-9*time.Minute))
	policy := alertPolicy{
		AllowedModes:      parseModeSet("FLOW-SCOUT"),
		MinNotional:       5000,
		PositionCooldown:  30 * time.Minute,
		RepeatMinNotional: 25000,
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 12000, Price: 0.21, Asset: "asset-1", Transaction: "tx-small", KnownList: "flow"},
		{Time: now.Add(-1 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 26000, Price: 0.22, Asset: "asset-1", Transaction: "tx-huge", KnownList: "flow"},
	}

	got := qualifyingTrades(trades, st, now, policy, 10*time.Minute)
	if len(got) != 1 || got[0].Transaction != "tx-huge" {
		t.Fatalf("qualifying=%#v, want only repeat-min bypass tx-huge", got)
	}
}

func TestQualifyingTrades_BlocksNegativeEdgeWallet(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		AllowedModes: parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION"),
		MinNotional:  5000,
		EdgeBlocks: map[string]string{
			"0x1111111111111111111111111111111111111111": "1h edge -18.95pp over 1 samples",
		},
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 9000, Price: 0.4, Asset: "asset-1", Transaction: "tx-1", KnownList: "tape_probation"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 9000, Price: 0.4, Asset: "asset-2", Transaction: "tx-2", KnownList: "tape_probation"},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 1 || got[0].Wallet != "0x2222222222222222222222222222222222222222" {
		t.Fatalf("qualifying=%+v, want only unblocked wallet", got)
	}
}

func TestQualifyingTrades_AllowsEdgeHotKnownWallet(t *testing.T) {
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		AllowedModes:       parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION,EDGE-HOT"),
		MinNotional:        5000,
		ObserveMinNotional: 25000,
		ObserveMaxBot:      45,
		EdgeHotMinNotional: 1000,
		EdgeHotMaxBot:      45,
		EdgeHot: map[string]string{
			"0x1111111111111111111111111111111111111111": "edge-hot 100.0% win avg +4.50pp over 2 samples",
		},
	}
	trades := []tapeTrade{
		{Time: now.Add(-2 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 1200, Price: 0.4, Asset: "asset-1", Transaction: "tx-1", KnownList: "watch", Tier: "B", Bot: 34},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 1200, Price: 0.4, Asset: "asset-2", Transaction: "tx-2", KnownList: "watch", Tier: "B", Bot: 34},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 1200, Price: 0.4, Asset: "asset-3", Transaction: "tx-3", KnownList: "watch", Tier: "BOT", Bot: 62},
	}

	got := qualifyingTrades(trades, sentState{Sent: map[string]string{}}, now, policy, 10*time.Minute)
	if len(got) != 1 || got[0].Wallet != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("qualifying=%+v, want only edge-hot non-bot wallet", got)
	}
}

func TestTapeEdgeHotListMapsToEdgeHotMode(t *testing.T) {
	tr := tapeTrade{
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  1500,
		Price:     0.45,
		Asset:     "asset-1",
		KnownList: "tape_edgehot",
		Tier:      "C",
		Bot:       34,
	}
	policy := alertPolicy{
		AllowedModes:       parseModeSet("EDGE-HOT"),
		MinNotional:        5000,
		EdgeHotMinNotional: 1000,
		EdgeHotMaxBot:      45,
		EdgeHot: map[string]string{
			strings.ToLower(tr.Wallet): "edge-hot 100.0% win avg +8.00pp over 3 samples",
		},
	}

	if mode := alertMode(tr); mode != "EDGE-HOT" {
		t.Fatalf("mode=%s, want EDGE-HOT", mode)
	}
	ok, reason := alertDecision(tr, policy)
	if !ok || !strings.Contains(reason, "edge-hot") {
		t.Fatalf("alertDecision=%v/%q, want edge-hot alert", ok, reason)
	}
}

func TestAlertDecision_AllowsFlowScoutMode(t *testing.T) {
	modes := parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION")
	ensureRequiredAlertModes(modes, "EDGE-HOT", "FLOW-SCOUT")
	policy := alertPolicy{AllowedModes: modes, MinNotional: 5000}
	tr := tapeTrade{
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  10000,
		Price:     0.2,
		Asset:     "asset-1",
		KnownList: "flow",
		Tier:      "B",
		Bot:       23,
	}

	ok, reason := alertDecision(tr, policy)
	if !ok || reason != "mode allowed" {
		t.Fatalf("alertDecision=%v/%q, want allowed flow scout", ok, reason)
	}
}

func TestAlertDecision_BlocksReviewNoise(t *testing.T) {
	policy := alertPolicy{
		AllowedModes:       parseModeSet("FOLLOW-READY,CANDIDATE,PROBATION,OBSERVE"),
		MinNotional:        5000,
		ObserveMinNotional: 10000,
	}
	tr := tapeTrade{
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  50000,
		Price:     0.2,
		Asset:     "asset-1",
		KnownList: "review_noise",
	}

	if mode := alertMode(tr); mode != "REVIEW-NOISE" {
		t.Fatalf("mode=%q, want REVIEW-NOISE", mode)
	}
	ok, reason := alertDecision(tr, policy)
	if ok || reason != "review-noise excluded" {
		t.Fatalf("alertDecision=%v/%q, want review-noise block", ok, reason)
	}
}

func TestAlertDecision_ObserveRequiresConfiguredMinTier(t *testing.T) {
	policy := alertPolicy{
		ObserveMinNotional:  5000,
		ObserveMaxBot:       35,
		ObserveRequireKnown: true,
		ObserveMinTier:      "B",
	}
	base := tapeTrade{
		Wallet:   "0x1111111111111111111111111111111111111111",
		Side:     "BUY",
		Notional: 7000,
		Price:    0.2,
		Asset:    "asset-1",
		Tier:     "B",
		Bot:      34,
	}

	ok, reason := alertDecision(base, policy)
	if !ok || reason != "huge observe non-bot" {
		t.Fatalf("B-tier observe alertDecision=%v/%q, want allowed", ok, reason)
	}

	base.Tier = "C"
	ok, reason = alertDecision(base, policy)
	if ok || !strings.Contains(reason, "below observe min B") {
		t.Fatalf("C-tier observe alertDecision=%v/%q, want min-tier block", ok, reason)
	}
}

func TestLoadWalletStatuses_ReviewNoiseOverridesObserve(t *testing.T) {
	dir := t.TempDir()
	observe := filepath.Join(dir, "observe.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(observe, []byte(wallet+" # list=tape_observe tier=B bot=23.7\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reviewNoise, []byte(wallet+" # list=review_noise tier=B bot=23.7\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletStatuses(observe + "," + reviewNoise)
	if err != nil {
		t.Fatal(err)
	}
	if got[wallet].List != "review_noise" {
		t.Fatalf("List=%q, want review_noise", got[wallet].List)
	}
}

func TestLoadWalletStatuses_ConsensusResearchPriority(t *testing.T) {
	dir := t.TempDir()
	observe := filepath.Join(dir, "observe.txt")
	consensus := filepath.Join(dir, "consensus.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(observe, []byte(wallet+" # list=tape_observe tier=D bot=30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(consensus, []byte(wallet+" # list=consensus_research signals=1 roi=37.2% bot=30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reviewNoise, []byte(wallet+" # list=review_noise tier=D bot=30\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletStatuses(observe + "," + consensus)
	if err != nil {
		t.Fatal(err)
	}
	if got[wallet].List != "consensus_research" {
		t.Fatalf("List=%q, want consensus_research", got[wallet].List)
	}
	got, err = loadWalletStatuses(observe + "," + consensus + "," + reviewNoise)
	if err != nil {
		t.Fatal(err)
	}
	if got[wallet].List != "review_noise" {
		t.Fatalf("List=%q, want review_noise override", got[wallet].List)
	}
}

func TestAlertDecision_RequiresPositiveEdgeForProbation(t *testing.T) {
	wallet := "0x1111111111111111111111111111111111111111"
	tr := tapeTrade{
		Wallet:    wallet,
		Side:      "BUY",
		Notional:  9000,
		Price:     0.45,
		Asset:     "asset-1",
		KnownList: "tape_probation",
		Tier:      "C",
		Bot:       24,
	}
	policy := alertPolicy{
		AllowedModes:        parseModeSet("PROBATION"),
		RequirePositiveEdge: parseModeSet("PROBATION"),
		MinNotional:         5000,
		EdgeHotMaxBot:       45,
		EdgeHot:             map[string]string{},
	}

	ok, reason := alertDecision(tr, policy)
	if ok || reason != "PROBATION requires positive edge evidence" {
		t.Fatalf("alertDecision=%v/%q, want positive-edge block", ok, reason)
	}

	policy.EdgeHot[strings.ToLower(wallet)] = "edge-hot 100.0% win avg +4.50pp over 2 samples"
	ok, reason = alertDecision(tr, policy)
	if !ok || !strings.Contains(reason, "edge-hot") {
		t.Fatalf("alertDecision=%v/%q, want edge-hot allowed probation", ok, reason)
	}
}

func TestBuildBurstCandidates_FindsLatestCoveredFlowAccumulation(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	asset := "asset-1"
	trades := []tapeTrade{
		{Time: now.Add(-25 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 6000, Price: 0.19, Asset: asset, Transaction: "tx-old", KnownList: "flow", Market: "Will Argentina win?", Tier: "B", Bot: 20},
		{Time: now.Add(-10 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 1200, Price: 0.20, Asset: asset, Transaction: "tx-1", KnownList: "flow", Market: "Will Argentina win?", Tier: "B", Bot: 20},
		{Time: now.Add(-5 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3000, Price: 0.21, Asset: asset, Transaction: "tx-2", KnownList: "flow", Market: "Will Argentina win?", Tier: "B", Bot: 20},
		{Time: now.Add(-2 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3400, Price: 0.21, Asset: asset, Transaction: "tx-3", KnownList: "flow", Market: "Will Argentina win?", Tier: "B", Bot: 20},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 9000, Price: 0.40, Asset: "asset-2", KnownList: "tape_observe"},
	}
	policy := alertPolicy{
		AllowedModes: parseModeSet("FLOW-SCOUT"),
		MinNotional:  5000,
	}
	st := sentState{Sent: map[string]string{
		"tx-old|" + asset + "|" + wallet: now.Format(time.RFC3339),
	}}

	got := buildBurstCandidates(trades, st, policy, now, time.Hour, 10*time.Minute, 15*time.Minute, 5000, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("bursts len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Trades != 3 || got[0].TotalNotional != 7600 {
		t.Fatalf("burst trades/total=%d/%.0f, want 3/7600", got[0].Trades, got[0].TotalNotional)
	}
	if !got[0].AlreadySent {
		t.Fatalf("AlreadySent=false, want covered burst because same asset wallet was alerted")
	}
}

func TestObserveBurstTrades_BuildsSyntheticSplitWhaleAlert(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{
		{Time: now.Add(-9 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4618, Price: 0.40, Size: 11545, Asset: "asset-1", Transaction: "tx-1", Market: "Indiana Fever vs. Los Angeles Sparks", Outcome: "Indiana Fever", Category: "basketball", Tier: "B", Bot: 34.3},
		{Time: now.Add(-8 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4466, Price: 0.41, Size: 10892.682927, Asset: "asset-1", Transaction: "tx-2", Market: "Indiana Fever vs. Los Angeles Sparks", Outcome: "Indiana Fever", Category: "basketball", Tier: "B", Bot: 34.3},
		{Time: now.Add(-7 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 9000, Price: 0.50, Size: 18000, Asset: "asset-2", Transaction: "tx-3", Market: "Will France win?", Category: "soccer", Tier: "D", Bot: 12},
	}
	policy := alertPolicy{
		ObserveBurstMinNotional: 8000,
		ObserveRequireKnown:     true,
		ObserveMinTier:          "B",
		ObserveMaxBot:           35,
	}

	got := observeBurstTrades(trades, sentState{Sent: map[string]string{}}, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("observeBurstTrades len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Wallet != wallet || got[0].KnownList != "observe_burst" || alertMode(got[0]) != "OBSERVE-BURST" {
		t.Fatalf("synthetic=%+v, want observe burst for wallet", got[0])
	}
	if got[0].Notional != 9084 {
		t.Fatalf("notional=%.0f, want 9084", got[0].Notional)
	}
	if got[0].Price <= 0.404 || got[0].Price >= 0.406 {
		t.Fatalf("price=%.6f, want VWAP around 0.405", got[0].Price)
	}
}

func TestObserveBurstTrades_IncludesOlderLegWhenLastBuyFresh(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{
		{Time: now.Add(-14 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4618, Price: 0.40, Size: 11545, Asset: "asset-1", Transaction: "tx-1", Market: "Indiana Fever vs. Los Angeles Sparks", Outcome: "Indiana Fever", Category: "basketball", Tier: "B", Bot: 34.3},
		{Time: now.Add(-4 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4466, Price: 0.41, Size: 10892.682927, Asset: "asset-1", Transaction: "tx-2", Market: "Indiana Fever vs. Los Angeles Sparks", Outcome: "Indiana Fever", Category: "basketball", Tier: "B", Bot: 34.3},
	}
	policy := alertPolicy{
		ObserveBurstMinNotional: 8000,
		ObserveRequireKnown:     true,
		ObserveMinTier:          "B",
		ObserveMaxBot:           35,
	}

	got := observeBurstTrades(trades, sentState{Sent: map[string]string{}}, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("observeBurstTrades len=%d, want 1 with older first leg inside burst window: %#v", len(got), got)
	}
	if got[0].Notional != 9084 || got[0].Time != now.Add(-4*time.Minute) {
		t.Fatalf("burst notional/time=%.0f/%s, want 9084 and fresh last leg", got[0].Notional, got[0].Time)
	}
}

func TestObserveBurstTrades_SkipsCoveredWalletAsset(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	asset := "asset-1"
	trades := []tapeTrade{
		{Time: now.Add(-9 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4618, Price: 0.40, Size: 11545, Asset: asset, Transaction: "tx-1", Tier: "B", Bot: 34.3},
		{Time: now.Add(-8 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4466, Price: 0.41, Size: 10892.682927, Asset: asset, Transaction: "tx-2", Tier: "B", Bot: 34.3},
	}
	st := sentState{Sent: map[string]string{}, Details: map[string]sentDetail{}}
	markTrades(st, []tapeTrade{{Time: now.Add(-9 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 6000, Price: 0.40, Asset: asset, Transaction: "tx-old"}}, now.Add(-8*time.Minute))
	policy := alertPolicy{
		ObserveBurstMinNotional: 8000,
		ObserveRequireKnown:     true,
		ObserveMinTier:          "B",
		ObserveMaxBot:           35,
	}

	got := observeBurstTrades(trades, st, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 0 {
		t.Fatalf("observeBurstTrades len=%d, want covered burst skipped: %#v", len(got), got)
	}
}

func TestStaleObserveBurstTrades_LogsOnlyStaleUnseenBursts(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{
		{Time: now.Add(-20 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4618, Price: 0.40, Size: 11545, Asset: "asset-1", Transaction: "tx-1", Tier: "B", Bot: 34.3},
		{Time: now.Add(-19 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 4466, Price: 0.41, Size: 10892.682927, Asset: "asset-1", Transaction: "tx-2", Tier: "B", Bot: 34.3},
		{Time: now.Add(-9 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 4618, Price: 0.40, Size: 11545, Asset: "asset-2", Transaction: "tx-3", Tier: "B", Bot: 34.3},
		{Time: now.Add(-8 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 4466, Price: 0.41, Size: 10892.682927, Asset: "asset-2", Transaction: "tx-4", Tier: "B", Bot: 34.3},
	}
	policy := alertPolicy{
		ObserveBurstMinNotional: 8000,
		ObserveRequireKnown:     true,
		ObserveMinTier:          "B",
		ObserveMaxBot:           35,
	}

	got := staleObserveBurstTrades(trades, sentState{Sent: map[string]string{}}, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("staleObserveBurstTrades len=%d, want 1 stale burst: %#v", len(got), got)
	}
	if got[0].Wallet != wallet || got[0].Asset != "asset-1" || alertMode(got[0]) != "OBSERVE-BURST" {
		t.Fatalf("shadow burst=%+v, want stale observe burst for first wallet", got[0])
	}

	logged := map[string]struct{}{tradeKey(got[0]): {}}
	got = staleObserveBurstTrades(trades, sentState{Sent: map[string]string{}}, policy, logged, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 0 {
		t.Fatalf("staleObserveBurstTrades len=%d, want logged burst deduped: %#v", len(got), got)
	}
}

func TestShadowUnknownFlowTrades_BuildsShadowOnlyMultiMarketCandidate(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{
		{Time: now.Add(-7 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3000, Price: 0.63, Size: 4761.904762, Asset: "asset-1", Transaction: "tx-1", Market: "Will France win on 2026-07-09?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-4 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3000, Price: 0.74, Size: 4054.054054, Asset: "asset-2", Transaction: "tx-2", Market: "Argentina vs. Switzerland: Team to Advance", Outcome: "Argentina", Category: "soccer"},
		{Time: now.Add(-3 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 7000, Price: 0.50, Size: 14000, Asset: "asset-3", Transaction: "tx-3", Market: "Indiana Fever vs. Phoenix Mercury", Category: "basketball", Tier: "B", Bot: 20},
	}
	policy := alertPolicy{
		UnknownFlowMinNotional: 6000,
		UnknownFlowMinMarkets:  2,
		UnknownFlowMaxBot:      45,
	}

	got := shadowUnknownFlowTrades(trades, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("shadowUnknownFlowTrades len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Wallet != wallet || got[0].KnownList != "unknown_flow" || alertMode(got[0]) != "UNKNOWN-FLOW" {
		t.Fatalf("unknown flow trade=%+v, want UNKNOWN-FLOW synthetic for wallet", got[0])
	}
	if got[0].Notional != 6000 || got[0].Asset != "asset-2" {
		t.Fatalf("notional/asset=%.0f/%s, want 6000/latest asset-2", got[0].Notional, got[0].Asset)
	}
	if len(shadowUnknownFlowTrades(trades, policy, map[string]struct{}{tradeKey(got[0]): {}}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)) != 0 {
		t.Fatalf("logged unknown-flow candidate should be deduped")
	}
}

func TestShadowSeedFlowTrades_BuildsLowerThresholdUnknownMultiMarketCandidate(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0x2393e78ad6d21ea3dda86f98980461c971341a8d"
	trades := []tapeTrade{
		{Time: now.Add(-6 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 2000, Price: 0.62, Size: 3225.806452, Asset: "asset-france-win", Transaction: "tx-1", Market: "Will France win on 2026-07-09?", Outcome: "Yes", Category: "soccer", KnownList: "watch", Tier: "A", Bot: 22.4},
		{Time: now.Add(-5 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 1000, Price: 0.78, Size: 1282.051282, Asset: "asset-france-advance", Transaction: "tx-2", Market: "France vs. Morocco: Team to Advance", Outcome: "France", Category: "soccer", KnownList: "watch", Tier: "A", Bot: 22.4},
	}
	policy := alertPolicy{
		UnknownFlowMinNotional: 6000,
		UnknownFlowMinMarkets:  2,
		UnknownFlowMaxBot:      45,
		SeedFlowMinNotional:    3000,
		SeedFlowMinMarkets:     2,
	}

	got := shadowSeedFlowTrades(trades, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("shadowSeedFlowTrades len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Wallet != wallet || got[0].KnownList != "seed_flow" || alertMode(got[0]) != "SEED-FLOW" {
		t.Fatalf("seed flow trade=%+v, want SEED-FLOW synthetic for wallet", got[0])
	}
	if got[0].Notional != 3000 || got[0].Asset != "asset-france-advance" {
		t.Fatalf("notional/asset=%.0f/%s, want 3000/latest asset-france-advance", got[0].Notional, got[0].Asset)
	}
	if len(shadowUnknownFlowTrades(trades, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)) != 0 {
		t.Fatalf("seed-sized flow should stay below UNKNOWN-FLOW threshold")
	}
	if len(shadowSeedFlowTrades(trades, policy, map[string]struct{}{tradeKey(got[0]): {}}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)) != 0 {
		t.Fatalf("logged seed-flow candidate should be deduped")
	}
}

func TestShadowScoredFlowTrades_BuildsShadowOnlyLeaderboardWhaleFlow(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	wallet := "0xb3c1111111111111111111111111111111114837"
	trades := []tapeTrade{
		{Time: now.Add(-8 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3100, Price: 0.31, Size: 10000, Asset: "asset-spain-belgium", Transaction: "tx-1", Market: "Spain vs Belgium: Team to Advance", Outcome: "Spain", Category: "soccer", Tier: "B", Bot: 21.5},
		{Time: now.Add(-7 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 3000, Price: 0.60, Size: 5000, Asset: "asset-france-morocco", Transaction: "tx-2", Market: "France vs Morocco: Team to Advance", Outcome: "France", Category: "soccer", Tier: "B", Bot: 21.5},
		{Time: now.Add(-6 * time.Minute), Wallet: wallet, Side: "BUY", Notional: 1681, Price: 0.41, Size: 4100, Asset: "asset-norway-england", Transaction: "tx-3", Market: "Norway vs England: Team to Advance", Outcome: "Norway", Category: "soccer", Tier: "B", Bot: 21.5},
	}
	policy := alertPolicy{
		ScoredFlowMinNotional: 6000,
		ScoredFlowMinMarkets:  2,
		ScoredFlowMaxBot:      35,
		ScoredFlowMinTier:     "B",
	}

	got := shadowScoredFlowTrades(trades, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("shadowScoredFlowTrades len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Wallet != wallet || got[0].KnownList != "scored_flow" || alertMode(got[0]) != "SCORED-FLOW" {
		t.Fatalf("scored flow trade=%+v, want SCORED-FLOW synthetic for wallet", got[0])
	}
	if got[0].Notional != 7781 || got[0].Asset != "asset-norway-england" || got[0].Tier != "B" || got[0].Bot != 21.5 {
		t.Fatalf("notional/asset/tier/bot=%.0f/%s/%s/%.1f, want 7781/latest/B/21.5", got[0].Notional, got[0].Asset, got[0].Tier, got[0].Bot)
	}
	if len(shadowScoredFlowTrades(trades, policy, map[string]struct{}{tradeKey(got[0]): {}}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)) != 0 {
		t.Fatalf("logged scored-flow candidate should be deduped")
	}
}

func TestShadowScoredFlowTrades_FiltersBotsLowTierAndHighBot(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	trades := []tapeTrade{
		{Time: now.Add(-8 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 4000, Price: 0.31, Asset: "asset-1", Market: "A", Tier: "C", Bot: 20},
		{Time: now.Add(-7 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 4000, Price: 0.32, Asset: "asset-2", Market: "B", Tier: "C", Bot: 20},
		{Time: now.Add(-8 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 4000, Price: 0.31, Asset: "asset-3", Market: "C", Tier: "B", Bot: 40},
		{Time: now.Add(-7 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 4000, Price: 0.32, Asset: "asset-4", Market: "D", Tier: "B", Bot: 40},
		{Time: now.Add(-8 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 4000, Price: 0.31, Asset: "asset-5", Market: "E", Tier: "BOT", Bot: 20},
		{Time: now.Add(-7 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 4000, Price: 0.32, Asset: "asset-6", Market: "F", Tier: "BOT", Bot: 20},
	}
	policy := alertPolicy{
		ScoredFlowMinNotional: 6000,
		ScoredFlowMinMarkets:  2,
		ScoredFlowMaxBot:      35,
		ScoredFlowMinTier:     "B",
	}

	got := shadowScoredFlowTrades(trades, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 0 {
		t.Fatalf("shadowScoredFlowTrades len=%d, want filters to reject all: %#v", len(got), got)
	}
}

func TestConsensusTrades_BuildsSyntheticAlertAndDedupes(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		ConsensusAlerts:      true,
		ConsensusMinNotional: 10000,
		ConsensusMinWallets:  2,
		ConsensusMaxBot:      60,
		EdgeBlocks: map[string]string{
			"0x5555555555555555555555555555555555555555": "1h edge -10pp over 1 samples",
		},
	}
	trades := []tapeTrade{
		{Time: now.Add(-7 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 7000, Price: 0.30, Size: 23333.333333, Asset: "asset-1", Tier: "B", Bot: 20, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-4 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "C", Bot: 30, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-3 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-2", Tier: "BOT", Bot: 80, Market: "Will Spain win?"},
		{Time: now.Add(-2 * time.Minute), Wallet: "0x4444444444444444444444444444444444444444", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-2", KnownList: "review_noise", Tier: "D", Bot: 20, Market: "Will Spain win?"},
		{Time: now.Add(-1 * time.Minute), Wallet: "0x5555555555555555555555555555555555555555", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-3", Tier: "B", Bot: 20, Market: "Will Germany win?"},
		{Time: now.Add(-1 * time.Minute), Wallet: "0x6666666666666666666666666666666666666666", Side: "BUY", Notional: 20000, Price: 0.20, Size: 100000, Asset: "asset-3", Tier: "B", Bot: 20, Market: "Will Germany win?"},
	}

	got := consensusTrades(trades, sentState{Sent: map[string]string{}}, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("consensus trades len=%d, want 1: %#v", len(got), got)
	}
	if got[0].Wallet != "multi:2" || got[0].KnownList != "consensus" || alertMode(got[0]) != "CONSENSUS" {
		t.Fatalf("synthetic trade=%+v, want CONSENSUS multi:2", got[0])
	}
	if got[0].Notional != 15000 || got[0].Price < 0.310 || got[0].Price > 0.311 {
		t.Fatalf("notional/price=%.0f/%.6f, want 15000/~0.3103", got[0].Notional, got[0].Price)
	}
	if len(got[0].Participants) != 2 || got[0].Participants[0] != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("participants=%v, want sorted consensus participants", got[0].Participants)
	}
	ev := alertEventFromTrade(got[0], tradeKey(got[0]), now, true, false)
	if len(ev.Participants) != 2 || ev.Participants[1] != "0x2222222222222222222222222222222222222222" {
		t.Fatalf("event participants=%v, want consensus participants in audit log", ev.Participants)
	}

	st := sentState{Sent: map[string]string{}, Details: map[string]sentDetail{}}
	markTrades(st, got, now)
	again := consensusTrades(trades, st, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(again) != 0 {
		t.Fatalf("consensus after mark len=%d, want 0: %#v", len(again), again)
	}
}

func TestConsensusTrades_IncludesOlderLegWhenLastBuyFresh(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		ConsensusAlerts:      true,
		ConsensusMinNotional: 10000,
		ConsensusMinWallets:  2,
		ConsensusMaxBot:      60,
	}
	trades := []tapeTrade{
		{Time: now.Add(-14 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 7000, Price: 0.30, Size: 23333.333333, Asset: "asset-1", Tier: "B", Bot: 20, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-4 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "C", Bot: 30, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
	}

	got := consensusTrades(trades, sentState{Sent: map[string]string{}}, policy, now, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("consensusTrades len=%d, want 1 with older first leg inside burst window: %#v", len(got), got)
	}
	if got[0].Notional != 15000 || got[0].Time != now.Add(-4*time.Minute) {
		t.Fatalf("consensus notional/time=%.0f/%s, want 15000 and fresh last leg", got[0].Notional, got[0].Time)
	}
}

func TestStaleConsensusTrades_LogsOnlyStaleUnseenConsensus(t *testing.T) {
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		ConsensusAlerts:      true,
		ConsensusMinNotional: 10000,
		ConsensusMinWallets:  2,
		ConsensusMaxBot:      60,
	}
	trades := []tapeTrade{
		{Time: now.Add(-21 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 7000, Price: 0.30, Size: 23333.333333, Asset: "asset-1", Tier: "B", Bot: 20, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-20 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "C", Bot: 30, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-7 * time.Minute), Wallet: "0x3333333333333333333333333333333333333333", Side: "BUY", Notional: 7000, Price: 0.30, Size: 23333.333333, Asset: "asset-2", Tier: "B", Bot: 20, Market: "Will Spain win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-4 * time.Minute), Wallet: "0x4444444444444444444444444444444444444444", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-2", Tier: "C", Bot: 30, Market: "Will Spain win?", Outcome: "Yes", Category: "soccer"},
	}

	got := staleConsensusTrades(trades, sentState{Sent: map[string]string{}}, policy, map[string]struct{}{}, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 1 {
		t.Fatalf("staleConsensusTrades len=%d, want 1 stale consensus: %#v", len(got), got)
	}
	if got[0].Asset != "asset-1" || got[0].KnownList != "consensus" || alertMode(got[0]) != "CONSENSUS" {
		t.Fatalf("shadow consensus=%+v, want stale consensus for asset-1", got[0])
	}

	logged := map[string]struct{}{tradeKey(got[0]): {}}
	got = staleConsensusTrades(trades, sentState{Sent: map[string]string{}}, policy, logged, now, time.Hour, 10*time.Minute, 15*time.Minute, 2, 1000)
	if len(got) != 0 {
		t.Fatalf("staleConsensusTrades len=%d, want logged consensus deduped: %#v", len(got), got)
	}
}

func TestWriteCandidateDiagnostics_IncludesConsensusBursts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candidates.md")
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	policy := alertPolicy{
		ConsensusAlerts:      true,
		ConsensusMinNotional: 10000,
		ConsensusMinWallets:  2,
		ConsensusMaxBot:      60,
	}
	trades := []tapeTrade{
		{Time: now.Add(-7 * time.Minute), Wallet: "0x1111111111111111111111111111111111111111", Side: "BUY", Notional: 7000, Price: 0.30, Size: 23333.333333, Asset: "asset-1", Tier: "B", Bot: 20, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
		{Time: now.Add(-4 * time.Minute), Wallet: "0x2222222222222222222222222222222222222222", Side: "BUY", Notional: 8000, Price: 0.32, Size: 25000, Asset: "asset-1", Tier: "C", Bot: 30, Market: "Will France win?", Outcome: "Yes", Category: "soccer"},
	}

	if err := writeCandidateDiagnostics(path, trades, sentState{Sent: map[string]string{}}, policy, now, time.Hour, 10*time.Minute, 15*time.Minute, 15*time.Minute, 5000, 2, 1000, 50); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{"## Consensus Bursts", "| consensus | 2 | 2 | $15000", "0x1111...1111", "Will France win?"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diagnostic missing %q:\n%s", want, got)
		}
	}
}

func TestWriteCandidateDiagnostics_UsesSentModeFromAuditState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candidates.md")
	now := time.Date(2026, 7, 9, 4, 0, 0, 0, time.UTC)
	tr := tapeTrade{
		Time:        now.Add(-5 * time.Minute),
		Wallet:      "0x1111111111111111111111111111111111111111",
		Side:        "BUY",
		Notional:    6000,
		Price:       0.2,
		Asset:       "asset-1",
		Transaction: "0xabc",
		KnownList:   "scout",
		Tier:        "B",
		Bot:         20,
		Market:      "Will Argentina win?",
	}
	key := tradeKey(tr)
	st := sentState{
		Sent: map[string]string{key: now.Add(-4 * time.Minute).Format(time.RFC3339)},
		Details: map[string]sentDetail{
			key: {Mode: "FLOW-SCOUT", KnownList: "flow"},
		},
	}
	policy := alertPolicy{
		AllowedModes:       parseModeSet("EDGE-HOT,FLOW-SCOUT"),
		MinNotional:        5000,
		EdgeHotMinNotional: 1000,
		EdgeHotMaxBot:      45,
		EdgeHot: map[string]string{
			strings.ToLower(tr.Wallet): "edge-hot 100.0% win avg +5.00pp over 3 samples",
		},
	}

	if err := writeCandidateDiagnostics(path, []tapeTrade{tr}, st, policy, now, time.Hour, 10*time.Minute, 15*time.Minute, 15*time.Minute, 5000, 2, 1000, 50); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{"| sent | FLOW-SCOUT |", "already sent as FLOW-SCOUT"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diagnostic missing %q:\n%s", want, got)
		}
	}
}

func TestEdgeMetrics_BuildBlocksAndHotWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	data := "" +
		`{"wallet":"0x1111111111111111111111111111111111111111","horizon_sec":900,"delta_pp":-2.0}` + "\n" +
		`{"wallet":"0x1111111111111111111111111111111111111111","horizon_sec":900,"delta_pp":-1.0}` + "\n" +
		`{"wallet":"0x2222222222222222222222222222222222222222","horizon_sec":3600,"delta_pp":-18.95}` + "\n" +
		`{"wallet":"0x3333333333333333333333333333333333333333","horizon_sec":900,"delta_pp":-10.0}` + "\n" +
		`{"wallet":"0x4444444444444444444444444444444444444444","horizon_sec":3600,"delta_pp":-4.9}` + "\n" +
		`{"wallet":"0x5555555555555555555555555555555555555555","horizon_sec":300,"delta_pp":2.0}` + "\n" +
		`{"wallet":"0x5555555555555555555555555555555555555555","horizon_sec":900,"delta_pp":5.0}` + "\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	metrics, err := loadEdgeMetrics(path)
	if err != nil {
		t.Fatal(err)
	}
	blocks := negativeEdgeBlocks(metrics, 2, -1, 1, -5)
	if len(blocks) != 2 {
		t.Fatalf("blocks len=%d, want 2: %#v", len(blocks), blocks)
	}
	if !strings.Contains(blocks["0x1111111111111111111111111111111111111111"], "15m edge") {
		t.Fatalf("first block=%q, want 15m reason", blocks["0x1111111111111111111111111111111111111111"])
	}
	if !strings.Contains(blocks["0x2222222222222222222222222222222222222222"], "1h edge") {
		t.Fatalf("second block=%q, want 1h reason", blocks["0x2222222222222222222222222222222222222222"])
	}
	hot := positiveEdgeHot(metrics, 2, 2, 60, 0.5, 0, -5)
	if len(hot) != 1 {
		t.Fatalf("hot len=%d, want 1: %#v", len(hot), hot)
	}
	if !strings.Contains(hot["0x5555555555555555555555555555555555555555"], "edge-hot") {
		t.Fatalf("hot reason=%q, want edge-hot reason", hot["0x5555555555555555555555555555555555555555"])
	}
}

func TestWriteCandidateDiagnostics_IncludesNearEdgeHotGaps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candidates.md")
	now := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{{
		Time:      now.Add(-2 * time.Minute),
		Wallet:    wallet,
		Side:      "BUY",
		Notional:  2500,
		Price:     0.42,
		Asset:     "asset-1",
		KnownList: "watch",
		Tier:      "B",
		Bot:       20,
		Market:    "Will Spain win on 2026-07-10?",
	}}
	policy := alertPolicy{
		AllowedModes:       parseModeSet("EDGE-HOT"),
		EdgeMetrics:        map[string]map[int64]*edgeStats{wallet: {int64((15 * time.Minute).Seconds()): {Samples: 1, Wins: 1, SumPP: 4}}},
		EdgeHot:            map[string]string{},
		EdgeBlocks:         map[string]string{},
		EdgeHotMinSamples:  2,
		EdgeHotMinAvgPP:    2,
		EdgeHotMinWinRate:  60,
		EdgeHotMin15mAvgPP: 0,
		EdgeHotMax1hNegPP:  -5,
		EdgeHotMinNotional: 1000,
		EdgeHotMaxBot:      45,
	}
	if err := writeCandidateDiagnostics(path, trades, sentState{Sent: map[string]string{}}, policy, now, time.Hour, 10*time.Minute, 15*time.Minute, 15*time.Minute, 5000, 2, 1000, 50); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	for _, want := range []string{"## Near Edge-Hot Wallets", "samples 1 < 2", "Will Spain win on 2026-07-10?"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diagnostic report missing %q:\n%s", want, got)
		}
	}
}

func TestFormatAlert_RendersObservationWarning(t *testing.T) {
	tr := tapeTrade{
		Time:          time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC),
		Wallet:        "0x1111111111111111111111111111111111111111",
		Side:          "BUY",
		Notional:      6000,
		Price:         0.42,
		Outcome:       "Yes",
		Market:        "Will France win on 2026-07-09?",
		Slug:          "fifwc-fra-mar-2026-07-09",
		Category:      "soccer",
		Asset:         "asset-1",
		Tier:          "TAPE",
		Bot:           12.3,
		TargetCopyROI: 45.6,
		TargetCopyT:   7,
	}

	msg := formatAlert(tr)
	for _, want := range []string{"Polymarket sports whale tape", "OBSERVE", "BUY $6000", "observation alert only", "https://polymarket.com/event/fifwc-fra-mar-2026-07-09"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("alert missing %q:\n%s", want, msg)
		}
	}
}

func TestFormatAlert_RendersReversalRisk(t *testing.T) {
	tr := tapeTrade{
		Time:      time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC),
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  8000,
		Price:     0.37,
		Outcome:   "No",
		Market:    "Will a team win?",
		Category:  "basketball",
		KnownList: "tape_reversal",
		Tier:      "B",
		Bot:       28,
	}

	msg := formatAlert(tr)
	for _, want := range []string{"REVERSAL-RISK", "list=tape_reversal", "late-window edge reversed", "do not treat as follow-ready"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("alert missing %q:\n%s", want, msg)
		}
	}
}

func TestFormatAlert_RendersProbationWarning(t *testing.T) {
	tr := tapeTrade{
		Time:      time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC),
		Wallet:    "0x1111111111111111111111111111111111111111",
		Side:      "BUY",
		Notional:  5000,
		Price:     0.45,
		Outcome:   "Team Nemesis",
		Market:    "Dota 2: Team Spirit vs Team Nemesis - Game 2 Winner",
		Category:  "esports",
		KnownList: "tape_probation",
		Tier:      "C",
		Bot:       43.6,
	}

	msg := formatAlert(tr)
	for _, want := range []string{"PROBATION", "list=tape_probation", "soft flow risks", "edge observation only"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("alert missing %q:\n%s", want, msg)
		}
	}
}

func TestLoadWalletStatuses_CommaSeparatedFilesOverrideInOrder(t *testing.T) {
	dir := t.TempDir()
	candidate := filepath.Join(dir, "candidate.txt")
	reversal := filepath.Join(dir, "reversal.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(candidate, []byte(wallet+" # list=tape_candidate tier=C bot=18 edgeN=3 edgeAvgPP=+2.10 reason=good early edge\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reversal, []byte(wallet+" # list=tape_reversal tier=C bot=18 edgeN=5 edgeAvgPP=+0.20 edge15mPP=-2.49 reason=15m edge reversed -2.49pp over 2 samples\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	statuses, err := loadWalletStatuses(candidate + "," + reversal)
	if err != nil {
		t.Fatal(err)
	}
	got := statuses[wallet]
	if got.List != "tape_reversal" {
		t.Fatalf("List=%q, want tape_reversal", got.List)
	}
	if got.EdgeN != 5 || got.Edge15m != -2.49 {
		t.Fatalf("edge stats=%d/%.2f, want 5/-2.49", got.EdgeN, got.Edge15m)
	}
	if !strings.Contains(got.Reason, "15m edge reversed") {
		t.Fatalf("Reason=%q, want reversal reason", got.Reason)
	}
}

func TestLoadWalletStatuses_RiskStatusBeatsLaterObserve(t *testing.T) {
	dir := t.TempDir()
	reversal := filepath.Join(dir, "reversal.txt")
	observe := filepath.Join(dir, "observe.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(reversal, []byte(wallet+" # list=tape_reversal tier=B bot=24 reason=severe sports alert drawdown\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(observe, []byte(wallet+" # list=tape_observe tier=B bot=24\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	statuses, err := loadWalletStatuses(reversal + "," + observe)
	if err != nil {
		t.Fatal(err)
	}
	got := statuses[wallet]
	if got.List != "tape_reversal" {
		t.Fatalf("List=%q, want tape_reversal to beat later observe", got.List)
	}
	if !strings.Contains(got.Reason, "drawdown") {
		t.Fatalf("Reason=%q, want reversal reason", got.Reason)
	}
}

func TestLoadWalletStatuses_BlocksRejectedTapeObserveStatus(t *testing.T) {
	dir := t.TempDir()
	observe := filepath.Join(dir, "observe.txt")
	wallet := "0xd72804b664e82152476670ba32f3e75f1cbb9e9b"
	if err := os.WriteFile(observe, []byte(wallet+" # list=tape_observe status=reject-flow tier=D bot=27.5 maxBuy=$55000\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	statuses, err := loadWalletStatuses(observe)
	if err != nil {
		t.Fatal(err)
	}
	got := statuses[wallet]
	if got.List != "review_noise" {
		t.Fatalf("List=%q, want review_noise", got.List)
	}
	tr := tapeTrade{Wallet: wallet, Side: "BUY", Notional: 55000, Price: 0.326, Asset: "asset-1"}
	applyWalletStatuses([]tapeTrade{tr}, statuses)
	tr.KnownList = got.List
	tr.Tier = got.Tier
	tr.Bot = got.Bot
	ok, reason := alertDecision(tr, alertPolicy{ObserveMinNotional: 10000, ObserveRequireKnown: true})
	if ok || reason != "review-noise excluded" {
		t.Fatalf("alertDecision=%v/%q, want rejected tape observe blocked", ok, reason)
	}
}

func TestApplyWalletStatuses_OverridesTapeLabel(t *testing.T) {
	trades := []tapeTrade{{
		Wallet:    "0x1111111111111111111111111111111111111111",
		KnownList: "tape_candidate",
		Tier:      "C",
		Bot:       18,
	}}
	statuses := map[string]walletStatus{
		"0x1111111111111111111111111111111111111111": {List: "tape_reversal", Tier: "C", Bot: 28},
	}

	applyWalletStatuses(trades, statuses)
	if trades[0].KnownList != "tape_reversal" || trades[0].Bot != 28 {
		t.Fatalf("status=%s bot=%.1f, want tape_reversal/28", trades[0].KnownList, trades[0].Bot)
	}
}

func TestAppendAlertLog_WritesAuditEvent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.jsonl")
	sentAt := time.Date(2026, 7, 9, 3, 10, 0, 0, time.UTC)
	trades := []tapeTrade{{
		Time:          time.Date(2026, 7, 9, 3, 9, 0, 0, time.UTC),
		Wallet:        "0x1111111111111111111111111111111111111111",
		Side:          "BUY",
		Notional:      6000,
		Price:         0.42,
		Outcome:       "Yes",
		Market:        "Will a team win?",
		Slug:          "team-win",
		Category:      "soccer",
		Asset:         "asset-1",
		Transaction:   "0xabc",
		KnownList:     "tape_probation",
		Tier:          "C",
		Bot:           43.6,
		TargetCopyROI: 75.5,
		TargetCopyT:   10,
	}}

	if err := appendAlertLog(path, trades, sentAt, true, false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 1 {
		t.Fatalf("log lines=%d, want 1", len(lines))
	}
	var ev alertEvent
	if err := json.Unmarshal([]byte(lines[0]), &ev); err != nil {
		t.Fatal(err)
	}
	if ev.Mode != "PROBATION" || ev.KnownList != "tape_probation" || !ev.DryRun {
		t.Fatalf("mode/list/dry=%s/%s/%v, want PROBATION/tape_probation/true", ev.Mode, ev.KnownList, ev.DryRun)
	}
	if ev.Wallet != "0x1111111111111111111111111111111111111111" || ev.Market != "Will a team win?" {
		t.Fatalf("wallet/market=%s/%s, want audit trade details", ev.Wallet, ev.Market)
	}
	if ev.Key != "0xabc|asset-1|0x1111111111111111111111111111111111111111" {
		t.Fatalf("key=%s, want transaction asset wallet key", ev.Key)
	}
}

func TestReconcileAlertLog_BackfillsStateWithoutDuplicates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.jsonl")
	sentAt := time.Date(2026, 7, 9, 3, 10, 0, 0, time.UTC)
	tr := tapeTrade{
		Time:        time.Date(2026, 7, 9, 3, 9, 0, 0, time.UTC),
		Wallet:      "0x2222222222222222222222222222222222222222",
		Side:        "BUY",
		Notional:    9000,
		Price:       0.55,
		Outcome:     "No",
		Market:      "Will another team win?",
		Category:    "basketball",
		Asset:       "asset-2",
		Transaction: "0xdef",
		KnownList:   "tape_follow",
		Tier:        "B",
	}
	st := sentState{Sent: map[string]string{tradeKey(tr): sentAt.Format(time.RFC3339)}}

	if err := reconcileAlertLog(path, []tapeTrade{tr}, st); err != nil {
		t.Fatal(err)
	}
	if err := reconcileAlertLog(path, []tapeTrade{tr}, st); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 1 {
		t.Fatalf("log lines=%d, want 1 after duplicate reconcile", len(lines))
	}
	var ev alertEvent
	if err := json.Unmarshal([]byte(lines[0]), &ev); err != nil {
		t.Fatal(err)
	}
	if !ev.Reconciled || ev.DryRun {
		t.Fatalf("reconciled/dry=%v/%v, want true/false", ev.Reconciled, ev.DryRun)
	}
	if !ev.SentAt.Equal(sentAt) || ev.Mode != "FOLLOW-READY" {
		t.Fatalf("sentAt/mode=%s/%s, want %s/FOLLOW-READY", ev.SentAt, ev.Mode, sentAt)
	}
}
