package feed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestIsLoLMarket(t *testing.T) {
	cases := []struct {
		name string
		m    Market
		want bool
	}{
		{"lol lck match", Market{Question: "LoL: Gen.G vs T1 - LCK Spring", Slug: "lol-gen-g-vs-t1-2026-04-20"}, true},
		{"lol lpl match", Market{Question: "LoL: EDG vs JDG (BO3) - LPL Summer", Slug: "lol-edg-jdg-2026-04-20"}, true},
		{"lol lec blocked", Market{Question: "LoL: Fnatic vs G2 - LEC Regular Season", Slug: "lol-fnc-g2-2026-04-20"}, false},
		{"lol ljl blocked", Market{Question: "LoL: DFM vs SGB - LJL 2026 Spring", Slug: "lol-dfm-sgb-2026-04-20"}, false},
		{"lol no league blocked", Market{Question: "LoL: Gen.G vs T1", Slug: "lol-gen-g-vs-t1-2026-04-20"}, false},
		{"lol slug only no league", Market{Question: "Who wins?", Slug: "lol-worlds-final-2026"}, false},
		{"league of legends lck", Market{Question: "League of Legends LCK finals winner"}, true},
		{"lck challengers blocked", Market{Question: "LoL: HANJIN BRION Challengers vs Hanwha Life Esports Challengers (BO3) - LCK Challengers League Rounds 1-2", Slug: "lol-lck-challengers-hanjin-2026"}, false},
		{"lck academy blocked", Market{Question: "LoL: Gen.G Global Academy vs Nongshim Esports Academy (BO3) - LCK Challengers League Rounds 1-2", Slug: "lol-lck-academy-2026"}, false},
		{"lpl developing blocked", Market{Question: "LoL: TES.A vs WBG.A (BO3) - LPL Developing League", Slug: "lol-lpl-developing-2026"}, false},
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
		{"nba-min-den-2026-04-20-spread-home-6pt5", false}, // derivative
		{"nba-min-den-2026-04-20-total-231pt5", false},     // derivative
		{"nba-atl-nyk-2026-04-20-spread-home-5pt5", false}, // derivative
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
		{"will-manchester-city-win-2025-26", false},        // seasonal
		{"nfl-dal-nyg-2026-09-08", false},                  // not in scope
		{"epl-cry-wes-2026-04-20-spread-away-2pt5", false}, // derivative
	}
	for _, c := range cases {
		if got := IsFootballMarket(Market{Slug: c.slug}); got != c.want {
			t.Errorf("%s: got %v want %v", c.slug, got, c.want)
		}
	}
}

func TestIsDota2Market(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"dota2-gl-heroic-2026-04-22", true},
		{"dota2-gl-heroic-2026-04-22-game1", true},
		{"dota2-xtreme-ts8-2026-04-22", true},
		{"dota2-sar1-mouz-2026-04-22", true},
		{"dota2-satan-ivo-2026-04-21", true},
		{"will-a-team-from-china-win-dota-2-the-international-10", false}, // seasonal
		{"lol-gen-g-vs-t1-2026-04-20", false},
		{"nba-atl-nyk-2026-04-20", false},
		{"dota2-gl-heroic-2026-04-22-spread-home-5pt5", false}, // derivative
	}
	for _, c := range cases {
		if got := IsDota2Market(Market{Slug: c.slug}); got != c.want {
			t.Errorf("%s: got %v want %v", c.slug, got, c.want)
		}
	}
}

func TestFilterSports_UnionAndOrder(t *testing.T) {
	in := []Market{
		{Question: "LoL: T1 vs Gen.G - LCK Spring", Slug: "lol-lck-t1-geng-2026-04-19"},
		{Slug: "nba-atl-nyk-2026-04-20"},
		{Slug: "epl-cry-wes-2026-04-20-cry"},
		{Slug: "dota2-gl-heroic-2026-04-22"},
		{Slug: "will-the-lakers-win-the-2026-nba-finals"},  // excluded
		{Slug: "election-2026"},                            // excluded
		{Slug: "nba-playoffs-who-will-win-series-foo-bar"},
		{Question: "LoL: Fnatic vs G2 - LEC Regular", Slug: "lol-lec-fnc-g2-2026-04-19"}, // excluded (LEC)
	}
	out := FilterSports(in)
	if len(out) != 5 {
		t.Fatalf("want 5, got %d", len(out))
	}
	want := []string{
		"lol-lck-t1-geng-2026-04-19",
		"nba-atl-nyk-2026-04-20",
		"epl-cry-wes-2026-04-20-cry",
		"dota2-gl-heroic-2026-04-22",
		"nba-playoffs-who-will-win-series-foo-bar",
	}
	for i, m := range out {
		if m.Slug != want[i] {
			t.Errorf("pos %d: got %s want %s", i, m.Slug, want[i])
		}
	}
}

func TestGetByConditionIDs_FiltersAndDecodes(t *testing.T) {
	wantIDs := map[string]bool{
		"0xabc": true,
		"0xdef": true,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/markets" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		got := r.URL.Query()["condition_ids"]
		sort.Strings(got)
		want := []string{"0xabc", "0xdef"}
		if len(got) != len(want) {
			t.Errorf("condition_ids count: got %v want %v", got, want)
		}
		for _, id := range got {
			if !wantIDs[id] {
				t.Errorf("unexpected condition_id %q", id)
			}
		}
		// Return a mixed slice: one closed with outcomePrices, one still open.
		rows := []map[string]any{
			{
				"conditionId":   "0xabc",
				"question":      "LoL: A vs B - Game 1 Winner",
				"closed":        true,
				"outcomes":      `["A","B"]`,
				"outcomePrices": `["1","0"]`,
				"clobTokenIds":  `["101","102"]`,
			},
			{
				"conditionId":   "0xdef",
				"question":      "LoL: C vs D - Game 2 Winner",
				"closed":        false,
				"outcomes":      `["C","D"]`,
				"outcomePrices": `["0.65","0.35"]`,
				"clobTokenIds":  `["201","202"]`,
			},
		}
		_ = json.NewEncoder(w).Encode(rows)
	}))
	defer srv.Close()

	c := &GammaClient{http: &http.Client{Timeout: 3 * time.Second}, base: srv.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	got, err := c.GetByConditionIDs(ctx, []string{"0xabc", "", "0xdef"})
	if err != nil {
		t.Fatalf("GetByConditionIDs: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d markets, want 2", len(got))
	}
	byCond := map[string]Market{}
	for _, m := range got {
		byCond[m.ConditionID] = m
	}
	closed, ok := byCond["0xabc"]
	if !ok {
		t.Fatalf("0xabc missing")
	}
	if !closed.Closed {
		t.Errorf("0xabc should be Closed=true")
	}
	prices := closed.OutcomePrices()
	if len(prices) != 2 || prices[0] != "1" || prices[1] != "0" {
		t.Errorf("prices decode failed: %v", prices)
	}
	if open := byCond["0xdef"]; open.Closed {
		t.Errorf("0xdef should be Closed=false")
	}
}

func TestGetByConditionIDs_EmptyInput(t *testing.T) {
	c := NewGammaClient()
	got, err := c.GetByConditionIDs(context.Background(), nil)
	if err != nil {
		t.Errorf("empty should be no-op, got err=%v", err)
	}
	if len(got) != 0 {
		t.Errorf("empty should return nil slice, got len=%d", len(got))
	}
}
