package main

import (
	"strings"
	"testing"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func TestQualifiesLeaderboardWhaleRequiresRecentProfitSource(t *testing.T) {
	base := walletdiscover.WalletScore{
		SmartMoneyScore: 90,
		BotScore:        20,
		Sources:         map[string]int{"leaderboard_volume_7d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:      30,
			TargetTrades:     3,
			AvgTradeNotional: 1000,
		},
	}
	if qualifiesLeaderboardWhale(base, 70, 45, 20, 500, 1, 0) {
		t.Fatal("volume-only leaderboard wallet should not qualify as recommended whale")
	}

	base.Sources = map[string]int{"leaderboard_profit_7d": 1}
	if !qualifiesLeaderboardWhale(base, 70, 45, 20, 500, 1, 0) {
		t.Fatal("recent profit leaderboard wallet with clean stats should qualify")
	}

	base.Stats.TargetTrades = 0
	if qualifiesLeaderboardWhale(base, 70, 45, 20, 500, 1, 0) {
		t.Fatal("leaderboard wallet with no target-category trades should not qualify")
	}
}

func TestRenderReportExcludesReviewNoiseFromRecommendations(t *testing.T) {
	addr := "0x620d7e06ce27d16532c061eba9b46c7e1833c67f"
	scores := []walletdiscover.WalletScore{{
		Address:         addr,
		Tier:            "B",
		SmartMoneyScore: 90,
		BotScore:        20,
		Sources:         map[string]int{"leaderboard_profit_7d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:      40,
			TargetTrades:     4,
			AvgTradeNotional: 2500,
		},
	}}

	scan := renderReport(scores, nil, map[string]string{addr: "review_noise"}, "scores.json", "push.txt", "review-noise.txt", 25, 70, 45, 20, 500, 1, 0, 80, 45, 100, 300, 20, "B", 80, 35, 50, 1000, 5, 1, "C", 55, 45, 5, 300, 3, 1)
	report := scan.Report
	excludedSection := section(report, "## Excluded Leaderboard Whales")
	recommendedSection := section(report, "## Recommended Leaderboard Whales")
	if !strings.Contains(excludedSection, addr) {
		t.Fatalf("excluded section missing wallet:\n%s", excludedSection)
	}
	if strings.Contains(recommendedSection, addr) {
		t.Fatalf("recommended section contains excluded wallet:\n%s", recommendedSection)
	}
	if !strings.Contains(report, "- Excluded by quarantine/review-noise: 1") {
		t.Fatalf("summary missing excluded count:\n%s", report)
	}
}

func TestQualifiesStrictLeaderboardPushRequiresBTierAndCleanerFlow(t *testing.T) {
	base := walletdiscover.WalletScore{
		Tier:            "B",
		SmartMoneyScore: 90,
		BotScore:        25,
		Sources:         map[string]int{"leaderboard_profit_7d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:       80,
			TargetTrades:      8,
			TargetLargeTrades: 2,
			AvgTradeNotional:  2500,
		},
	}
	if !qualifiesStrictLeaderboardPush(base, "B", 80, 35, 50, 1000, 5, 1) {
		t.Fatal("clean B-tier recent-profit whale should qualify for strict push")
	}

	base.Tier = "C"
	if qualifiesStrictLeaderboardPush(base, "B", 80, 35, 50, 1000, 5, 1) {
		t.Fatal("C-tier wallet should not qualify for strict push")
	}

	base.Tier = "B"
	base.RiskFlags = []string{"open_copy_exposure"}
	if qualifiesStrictLeaderboardPush(base, "B", 80, 35, 50, 1000, 5, 1) {
		t.Fatal("open-copy exposure should keep leaderboard wallet out of strict push")
	}

	base.RiskFlags = nil
	base.Stats.TargetLargeTrades = 0
	if qualifiesStrictLeaderboardPush(base, "B", 80, 35, 50, 1000, 5, 1) {
		t.Fatal("strict push should require target-category large trades")
	}
}

func section(report, header string) string {
	start := strings.Index(report, header)
	if start < 0 {
		return ""
	}
	rest := report[start+len(header):]
	next := strings.Index(rest, "\n## ")
	if next < 0 {
		return rest
	}
	return rest[:next]
}

func TestHighValueWatchRequiresRecentProfitSource(t *testing.T) {
	volumeOnly := walletdiscover.WalletScore{
		SmartMoneyScore: 80,
		Sources:         map[string]int{"leaderboard_volume_30d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:      25,
			TargetTrades:     2,
			AvgTradeNotional: 750,
		},
	}
	if highValueWatch(volumeOnly, 1, 0) {
		t.Fatal("volume-only leaderboard wallet should not enter high-value watch")
	}

	volumeOnly.Sources = map[string]int{"leaderboard_profit_30d": 1}
	if !highValueWatch(volumeOnly, 1, 0) {
		t.Fatal("recent profit leaderboard wallet should enter high-value watch")
	}

	volumeOnly.Stats.TargetTrades = 0
	if highValueWatch(volumeOnly, 1, 0) {
		t.Fatal("high-value watch should require target-category activity")
	}
}

func TestQualifiesLeaderboardWhaleWatchAllowsCleanLargeLeaderboardFlow(t *testing.T) {
	score := walletdiscover.WalletScore{
		SmartMoneyScore: 92,
		BotScore:        40,
		Sources:         map[string]int{"leaderboard_volume_7d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:       160,
			TargetTrades:      70,
			TargetLargeTrades: 35,
			AvgTradeNotional:  450,
		},
		RiskFlags: []string{"burst_trading", "open_copy_exposure"},
	}
	if !qualifiesLeaderboardWhaleWatch(score, 80, 45, 100, 300, 1, 20) {
		t.Fatal("large low-bot leaderboard wallet should qualify for whale-watch")
	}
	score.RiskFlags = append(score.RiskFlags, "bot_like_flow")
	if qualifiesLeaderboardWhaleWatch(score, 80, 45, 100, 300, 1, 20) {
		t.Fatal("bot-like flow should not qualify for whale-watch")
	}
}

func TestQualifiesLeaderboardSportsPushAllowsBroadTargetCategoryLeaderboardFlow(t *testing.T) {
	score := walletdiscover.WalletScore{
		Tier:            "C",
		SmartMoneyScore: 62,
		BotScore:        38,
		Sources:         map[string]int{"leaderboard_volume_30d": 1},
		Stats: walletdiscover.WalletStats{
			LargeTrades:       12,
			TargetTrades:      8,
			TargetLargeTrades: 2,
			AvgTradeNotional:  420,
		},
	}
	if !qualifiesLeaderboardSportsPush(score, "C", 55, 45, 5, 300, 3, 1) {
		t.Fatal("broad sports leaderboard wallet should qualify for sports push")
	}

	score.Stats.TargetLargeTrades = 0
	if qualifiesLeaderboardSportsPush(score, "C", 55, 45, 5, 300, 3, 1) {
		t.Fatal("sports push should require large target-category positions")
	}

	score.Stats.TargetLargeTrades = 2
	score.RiskFlags = []string{"fixed_price"}
	if qualifiesLeaderboardSportsPush(score, "C", 55, 45, 5, 300, 3, 1) {
		t.Fatal("fixed-price flow should not qualify for sports push")
	}
}
