package walletdiscover

import "testing"

func TestRateFootballScoreWalletUsesScoreSpecificEvidence(t *testing.T) {
	score := WalletScore{
		Tier: "BOT",
		Stats: WalletStats{
			FootballScoreTrades:      30,
			FootballScoreLargeTrades: 12,
			FootballScoreCopyClosed:  10,
			FootballScoreCopyROI:     18,
			FootballScoreCopyWinRate: 60,
		},
	}
	got := RateFootballScoreWallet(score)
	if got.Tier != "A" {
		t.Fatalf("rating=%+v, want score-specific A", got)
	}
}

func TestRateFootballScoreWalletRejectsNegativeCopyEvidence(t *testing.T) {
	score := WalletScore{Stats: WalletStats{
		FootballScoreTrades:      30,
		FootballScoreLargeTrades: 12,
		FootballScoreCopyClosed:  5,
		FootballScoreCopyROI:     -1,
	}}
	got := RateFootballScoreWallet(score)
	if got.Tier != "D" {
		t.Fatalf("rating=%+v, want D", got)
	}
}

func TestRateFootballScoreWalletKeepsOpenOnlyHolderInObservation(t *testing.T) {
	score := WalletScore{Stats: WalletStats{
		FootballScoreTrades:      20,
		FootballScoreLargeTrades: 8,
		FootballScoreCopyOpen:    6,
	}}
	got := RateFootballScoreWallet(score)
	if got.Tier != "C" {
		t.Fatalf("rating=%+v, want C", got)
	}
}
