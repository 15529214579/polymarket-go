package walletdiscover

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestScoreCandidatesPreservesPreviousScoreWhenAPIsRemainIncomplete(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"
	candidate := &Candidate{Address: address, Sources: map[string]int{"leaderboard_profit_7d": 1}}
	cfg := DefaultConfig()
	cfg.DataBase = "https://data.test"
	cfg.OutputDir = t.TempDir()
	cfg.Concurrency = 1
	cfg.ActivityPages = 1
	cfg.HTTPMaxAttempts = 1
	previous := map[string]WalletScore{
		address: {
			Address: address, Tier: "A", FollowAction: "auto-small", Reason: "strong signal",
			SmartMoneyScore: 88, Stats: WalletStats{ValidTrades: 100, ClosedPositions: 25},
		},
	}

	client := NewClient(cfg)
	client.http.Transport = unavailableTransport()
	scores := scoreCandidates(context.Background(), client, []*Candidate{candidate}, cfg, previous)
	if len(scores) != 1 {
		t.Fatalf("scores=%d", len(scores))
	}
	score := scores[0]
	if score.Tier != "A" || score.FollowAction != "auto-small" || score.SmartMoneyScore != 88 {
		t.Fatalf("previous score was not preserved: %+v", score)
	}
	if score.DataStatus != "preserved_previous" || len(score.DataIssues) != 2 {
		t.Fatalf("data status=%q issues=%v", score.DataStatus, score.DataIssues)
	}
	if score.Sources["leaderboard_profit_7d"] != 1 {
		t.Fatalf("candidate sources were not refreshed: %v", score.Sources)
	}
}

func TestScoreCandidatesRejectsIncompleteWalletWithoutHistory(t *testing.T) {
	address := "0x0000000000000000000000000000000000000002"
	cfg := DefaultConfig()
	cfg.DataBase = "https://data.test"
	cfg.OutputDir = t.TempDir()
	cfg.Concurrency = 1
	cfg.ActivityPages = 1
	cfg.HTTPMaxAttempts = 1
	client := NewClient(cfg)
	client.http.Transport = unavailableTransport()
	scores := scoreCandidates(context.Background(), client, []*Candidate{{Address: address, Sources: map[string]int{}}}, cfg, nil)
	if len(scores) != 1 || scores[0].Tier != "D" || scores[0].FollowAction != "reject" || scores[0].DataStatus != "incomplete" {
		t.Fatalf("score=%+v", scores)
	}
}

func unavailableTransport() http.RoundTripper {
	return roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return testHTTPResponse(req, http.StatusServiceUnavailable, "temporarily unavailable", nil), nil
	})
}

func TestCandidatePriorityPrefersRecentProfitLeaderboard(t *testing.T) {
	profit := &Candidate{
		Address: "0x0000000000000000000000000000000000000001",
		Sources: map[string]int{"leaderboard_profit_7d": 1},
	}
	volume := &Candidate{
		Address:          "0x0000000000000000000000000000000000000002",
		Sources:          map[string]int{"leaderboard_volume_all": 1},
		ObservedNotional: 20000,
	}
	if candidatePriority(profit) <= candidatePriority(volume) {
		t.Fatalf("recent profit leaderboard should outrank pure volume leaderboard: profit=%f volume=%f", candidatePriority(profit), candidatePriority(volume))
	}
}

func TestAddSourceWalletsCarriesSportsTapeNotionalIntoPriority(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.sports-tape.txt")
	wallet := "0xd72804b664e82152476670ba32f3e75f1cbb9e9b"
	if err := os.WriteFile(path, []byte(wallet+" # source=sports_tape buys=1 buyNotional=$55000 maxBuy=$55000\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	candidates := map[string]*Candidate{}

	addSourceWallets(candidates, path, "sports_tape")

	c := candidates[wallet]
	if c == nil {
		t.Fatalf("candidate missing for %s", wallet)
	}
	if c.ObservedTrades != 1 || c.ObservedNotional != 55000 {
		t.Fatalf("observed=%d/%.0f, want 1/55000", c.ObservedTrades, c.ObservedNotional)
	}
	profit7d := &Candidate{
		Address: "0x0000000000000000000000000000000000000001",
		Sources: map[string]int{"leaderboard_profit_7d": 1},
	}
	if candidatePriority(c) <= candidatePriority(profit7d) {
		t.Fatalf("huge sports tape should outrank plain 7d profit candidate: tape=%f profit=%f", candidatePriority(c), candidatePriority(profit7d))
	}
}
