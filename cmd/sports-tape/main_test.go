package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func TestMergeRetainedTradesKeepsCleanRecentOnly(t *testing.T) {
	now := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)
	cfg := walletdiscover.DefaultConfig()
	cfg.TargetCategories = "basketball,soccer,esports"
	cfg.MinNotionalUSD = 500
	cfg.MinPrice = 0.05
	cfg.MaxPrice = 0.95

	current := []tapeTrade{
		{
			Time:        now.Add(-5 * time.Minute),
			Wallet:      "0x1111111111111111111111111111111111111111",
			Side:        "BUY",
			Notional:    1000,
			Price:       0.50,
			Size:        2000,
			Market:      "Dota 2: Team Spirit vs Team Nemesis - Game 2 Winner",
			Slug:        "dota2-spirit-nemesis-2026-07-09-game2",
			Category:    "esports",
			Asset:       "asset-current",
			ConditionID: "0xcurrent",
			Transaction: "0xcurrent",
		},
	}
	prior := []tapeTrade{
		current[0],
		{
			Time:        now.Add(-2 * time.Hour),
			Wallet:      "0x2222222222222222222222222222222222222222",
			Side:        "BUY",
			Notional:    6000,
			Price:       0.75,
			Size:        8000,
			Market:      "Golden State Valkyries vs. Toronto Tempo",
			Slug:        "wnba-gsv-tor-2026-07-08",
			Category:    "basketball",
			Asset:       "asset-keep",
			ConditionID: "0xkeep",
			Transaction: "0xkeep",
		},
		{
			Time:        now.Add(-30 * time.Minute),
			Wallet:      "0x3333333333333333333333333333333333333333",
			Side:        "BUY",
			Notional:    900,
			Price:       0.60,
			Size:        1500,
			Market:      "FC Petrocub Hincesti vs. KF Egnatia Rrogozhine: O/U 2.5",
			Slug:        "soccer-ou-2026-07-09",
			Category:    "soccer",
			Asset:       "asset-derivative",
			ConditionID: "0xderivative",
			Transaction: "0xderivative",
		},
		{
			Time:        now.Add(-8 * time.Hour),
			Wallet:      "0x4444444444444444444444444444444444444444",
			Side:        "BUY",
			Notional:    5000,
			Price:       0.50,
			Size:        10000,
			Market:      "Will England win on 2026-07-11?",
			Slug:        "fifwc-nor-eng-2026-07-11",
			Category:    "soccer",
			Asset:       "asset-stale",
			ConditionID: "0xstale",
			Transaction: "0xstale",
		},
	}

	got, retained := mergeRetainedTrades(current, prior, cfg, nil, now, 6*time.Hour)
	if retained != 1 {
		t.Fatalf("retained=%d, want 1", retained)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2: %#v", len(got), got)
	}
	if got[0].Transaction != "0xcurrent" || got[1].Transaction != "0xkeep" {
		t.Fatalf("unexpected retained ordering/content: %#v", got)
	}
}

func TestLoadPushMetas_CommaSeparatedFilesOverrideInOrder(t *testing.T) {
	dir := t.TempDir()
	observe := filepath.Join(dir, "observe.txt")
	reversal := filepath.Join(dir, "reversal.txt")
	wallet := "0x1111111111111111111111111111111111111111"
	if err := os.WriteFile(observe, []byte(wallet+" # list=tape_observe tier=D bot=46.5\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reversal, []byte(wallet+" # list=tape_reversal tier=D bot=46.5\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := loadPushMetas(observe + "," + reversal)
	if got[wallet].List != "tape_reversal" {
		t.Fatalf("List=%q, want tape_reversal", got[wallet].List)
	}
	if got[wallet].Tier != "D" || got[wallet].Bot != 46.5 {
		t.Fatalf("meta=%+v, want tier D bot 46.5", got[wallet])
	}
}

func TestLoadWalletSet_CommaSeparatedFiles(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.txt")
	second := filepath.Join(dir, "second.txt")
	walletA := "0x1111111111111111111111111111111111111111"
	walletB := "0x2222222222222222222222222222222222222222"
	if err := os.WriteFile(first, []byte(walletA+" # list=quarantine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte(walletB+" # list=review_noise\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := loadWalletSet(first + "," + filepath.Join(dir, "missing.txt") + "," + second)
	if _, ok := got[walletA]; !ok {
		t.Fatalf("missing first wallet in set: %#v", got)
	}
	if _, ok := got[walletB]; !ok {
		t.Fatalf("missing second wallet in set: %#v", got)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2: %#v", len(got), got)
	}
}

func TestApplyWalletMetas_RefreshesRetainedTradeStatus(t *testing.T) {
	wallet := "0x1111111111111111111111111111111111111111"
	trades := []tapeTrade{{
		Wallet:    wallet,
		KnownList: "tape_candidate",
		Tier:      "C",
		Bot:       30,
	}}
	metas := map[string]walletMeta{
		wallet: {List: "tape_reversal", Tier: "D", Bot: 46.5},
	}
	scores := map[string]scoreMeta{
		wallet: {Smart: 100, TargetCopyROI: 12.5, TargetCopyT: 3},
	}

	applyWalletMetas(trades, metas, scores)
	if trades[0].KnownList != "tape_reversal" {
		t.Fatalf("KnownList=%q, want tape_reversal", trades[0].KnownList)
	}
	if trades[0].Tier != "D" || trades[0].Bot != 46.5 || trades[0].Smart != 100 {
		t.Fatalf("trade meta=%+v, want refreshed status", trades[0])
	}
	if trades[0].TargetCopyROI != 12.5 || trades[0].TargetCopyT != 3 {
		t.Fatalf("target copy=%f/%d, want 12.5/3", trades[0].TargetCopyROI, trades[0].TargetCopyT)
	}
}

func TestMergeRetainedTradesSkipsExcludedWallets(t *testing.T) {
	now := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)
	cfg := walletdiscover.DefaultConfig()
	cfg.TargetCategories = "basketball,soccer,esports"
	cfg.MinNotionalUSD = 500
	cfg.MinPrice = 0.05
	cfg.MaxPrice = 0.95

	excludedWallet := "0x2222222222222222222222222222222222222222"
	prior := []tapeTrade{
		{
			Time:        now.Add(-2 * time.Hour),
			Wallet:      excludedWallet,
			Side:        "BUY",
			Notional:    6000,
			Price:       0.75,
			Size:        8000,
			Market:      "Golden State Valkyries vs. Toronto Tempo",
			Slug:        "wnba-gsv-tor-2026-07-08",
			Category:    "basketball",
			Asset:       "asset-excluded",
			ConditionID: "0xexcluded",
			Transaction: "0xexcluded",
		},
	}

	got, retained := mergeRetainedTrades(nil, prior, cfg, map[string]struct{}{excludedWallet: {}}, now, 6*time.Hour)
	if retained != 0 {
		t.Fatalf("retained=%d, want 0", retained)
	}
	if len(got) != 0 {
		t.Fatalf("len=%d, want 0: %#v", len(got), got)
	}
}
