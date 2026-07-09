package walletdiscover

import (
	"os"
	"path/filepath"
	"testing"
)

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
