package feed

import "testing"

func TestIsLoLMarket(t *testing.T) {
	cases := []struct {
		name string
		m    Market
		want bool
	}{
		{"lol prefix", Market{Question: "LoL: Gen.G vs T1", Slug: "lol-gen-g-vs-t1-2026-04-20"}, true},
		{"lol slug only", Market{Question: "Who wins?", Slug: "lol-worlds-final-2026"}, true},
		{"league of legends phrase", Market{Question: "League of Legends Worlds winner"}, true},
		{"election false positive", Market{Question: "2026 election winner", Slug: "election-2026"}, false},
	}
	for _, c := range cases {
		if got := IsLoLMarket(c.m); got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestIsBasketballMarket(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"nba-atl-nyk-2026-04-20", true},
		{"nba-phi-bos-2026-04-21", true},
		{"nba-playoffs-who-will-win-series-pistons-vs-magic", true},
		{"will-the-los-angeles-lakers-win-the-2026-nba-finals", false}, // seasonal
		{"nba-finals-2026-mvp", false},
		{"mlb-det-bos-2026-04-20", false},
		{"nba-min-den-2026-04-20-spread-home-6pt5", false},       // derivative
		{"nba-min-den-2026-04-20-total-231pt5", false},           // derivative
		{"nba-atl-nyk-2026-04-20-spread-home-5pt5", false},       // derivative
	}
	for _, c := range cases {
		if got := IsBasketballMarket(Market{Slug: c.slug}); got != c.want {
			t.Errorf("%s: got %v want %v", c.slug, got, c.want)
		}
	}
}

func TestIsFootballMarket(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"epl-cry-wes-2026-04-20-wes", true},
		{"epl-bur-mac-2026-04-22-mac", true},
		{"will-manchester-city-win-2025-26", false},                 // seasonal
		{"nfl-dal-nyg-2026-09-08", false},                           // not in scope
		{"epl-cry-wes-2026-04-20-spread-away-2pt5", false},          // derivative
	}
	for _, c := range cases {
		if got := IsFootballMarket(Market{Slug: c.slug}); got != c.want {
			t.Errorf("%s: got %v want %v", c.slug, got, c.want)
		}
	}
}

func TestFilterSports_UnionAndOrder(t *testing.T) {
	in := []Market{
		{Slug: "lol-lec-vit-gx-2026-04-19"},
		{Slug: "nba-atl-nyk-2026-04-20"},
		{Slug: "epl-cry-wes-2026-04-20-cry"},
		{Slug: "will-the-lakers-win-the-2026-nba-finals"}, // excluded
		{Slug: "election-2026"},                           // excluded
		{Slug: "nba-playoffs-who-will-win-series-foo-bar"},
	}
	out := FilterSports(in)
	if len(out) != 4 {
		t.Fatalf("want 4, got %d", len(out))
	}
	// Order preserved.
	want := []string{
		"lol-lec-vit-gx-2026-04-19",
		"nba-atl-nyk-2026-04-20",
		"epl-cry-wes-2026-04-20-cry",
		"nba-playoffs-who-will-win-series-foo-bar",
	}
	for i, m := range out {
		if m.Slug != want[i] {
			t.Errorf("pos %d: got %s want %s", i, m.Slug, want[i])
		}
	}
}
