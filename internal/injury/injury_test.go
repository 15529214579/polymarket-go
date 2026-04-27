package injury

import (
	"context"
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	cases := []struct {
		in   string
		want PlayerStatus
	}{
		{"Out", StatusOut},
		{"out", StatusOut},
		{"Day-To-Day (Out)", StatusOut},
		{"Doubtful", StatusDoubtful},
		{"Questionable", StatusQuest},
		{"Probable", StatusProb},
		{"Active", StatusAvail},
		{"", StatusAvail},
	}
	for _, c := range cases {
		got := parseStatus(c.in)
		if got != c.want {
			t.Errorf("parseStatus(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsStar(t *testing.T) {
	if !isStar("Denver Nuggets", "Nikola Jokic") {
		t.Error("Jokic should be a star")
	}
	if !isStar("Denver Nuggets", "nikola jokic") {
		t.Error("case-insensitive match should work")
	}
	if isStar("Denver Nuggets", "Random Bench Player") {
		t.Error("non-star should not match")
	}
	if isStar("Fake Team", "Anyone") {
		t.Error("unknown team should not match")
	}
}

func TestAssessImpact(t *testing.T) {
	e := InjuryEntry{Player: "Nikola Jokic", Team: "Denver Nuggets", Status: StatusOut}
	if got := assessImpact("Denver Nuggets", e); got != "franchise_player_out" {
		t.Errorf("Jokic impact = %q, want franchise_player_out", got)
	}

	e2 := InjuryEntry{Player: "Jamal Murray", Team: "Denver Nuggets", Status: StatusOut}
	if got := assessImpact("Denver Nuggets", e2); got != "co_star_out" {
		t.Errorf("Murray impact = %q, want co_star_out", got)
	}

	e3 := InjuryEntry{Player: "Aaron Gordon", Team: "Denver Nuggets", Status: StatusOut}
	if got := assessImpact("Denver Nuggets", e3); got != "rotation_star_out" {
		t.Errorf("Gordon impact = %q, want rotation_star_out", got)
	}
}

func TestScanFiltersCorrectly(t *testing.T) {
	entries := []InjuryEntry{
		{Player: "Nikola Jokic", Team: "Denver Nuggets", Status: StatusOut, Reason: "knee soreness"},
		{Player: "Random Bench Guy", Team: "Denver Nuggets", Status: StatusOut, Reason: "ankle"},
		{Player: "Jamal Murray", Team: "Denver Nuggets", Status: StatusQuest, Reason: "hamstring"},
	}

	// Simulate the scan filter logic
	var alerts []InjuryAlert
	byTeam := make(map[string][]InjuryEntry)
	for _, e := range entries {
		byTeam[e.Team] = append(byTeam[e.Team], e)
	}

	for team, teamEntries := range byTeam {
		for _, e := range teamEntries {
			if e.Status != StatusOut && e.Status != StatusDoubtful {
				continue
			}
			if !isStar(team, e.Player) {
				continue
			}
			alerts = append(alerts, InjuryAlert{
				Team:       team,
				StarPlayer: e.Player,
				Status:     e.Status,
				Impact:     assessImpact(team, e),
			})
		}
	}

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert (Jokic), got %d", len(alerts))
	}
	if alerts[0].StarPlayer != "Nikola Jokic" {
		t.Errorf("expected Jokic alert, got %s", alerts[0].StarPlayer)
	}
	if alerts[0].Impact != "franchise_player_out" {
		t.Errorf("expected franchise_player_out, got %s", alerts[0].Impact)
	}
}

func TestScanDedup(t *testing.T) {
	cfg := Config{Enabled: true, StarOnly: true}
	scanner := NewScanner(cfg, t.TempDir())

	// Manually add to seen with today's date-based key
	key := scanner.seenKey("Denver Nuggets", "Nikola Jokic", StatusOut)
	scanner.seen[key] = time.Now()

	entries := []InjuryEntry{
		{Player: "Nikola Jokic", Team: "Denver Nuggets", Status: StatusOut},
	}

	byTeam := map[string][]InjuryEntry{"Denver Nuggets": entries}
	var alerts []InjuryAlert
	for team, teamEntries := range byTeam {
		for _, e := range teamEntries {
			if e.Status != StatusOut && e.Status != StatusDoubtful {
				continue
			}
			if !isStar(team, e.Player) {
				continue
			}
			k := scanner.seenKey(team, e.Player, e.Status)
			if _, ok := scanner.seen[k]; ok {
				continue
			}
			alerts = append(alerts, InjuryAlert{StarPlayer: e.Player})
		}
	}

	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts (dedup), got %d", len(alerts))
	}
}

func TestSeenPersistence(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{Enabled: true, StarOnly: true}

	s1 := NewScanner(cfg, dir)
	key := s1.seenKey("Denver Nuggets", "Nikola Jokic", StatusOut)
	s1.seen[key] = time.Now()
	s1.saveSeen()

	s2 := NewScanner(cfg, dir)
	if _, ok := s2.seen[key]; !ok {
		t.Error("seen state not persisted across restarts")
	}
}

func TestDisabledScanReturnsNil(t *testing.T) {
	scanner := NewScanner(DefaultConfig(), t.TempDir())
	alerts, err := scanner.Scan(context.Background())
	if err != nil || alerts != nil {
		t.Errorf("disabled scan should return nil, nil; got %v, %v", alerts, err)
	}
}
