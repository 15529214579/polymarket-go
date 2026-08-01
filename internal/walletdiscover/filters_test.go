package walletdiscover

import (
	"testing"
	"time"
)

func TestIsNoisyMarket_FiltersShortCycleBTC(t *testing.T) {
	m := Market{
		Question:    "Bitcoin Up or Down - May 18, 11:55AM-12:00PM ET",
		Slug:        "btc-updown-5m-1779119700",
		ConditionID: "0xabc",
		Active:      true,
		Liquidity:   1000,
	}
	if !IsNoisyMarket(m) {
		t.Fatal("expected BTC up/down market to be noisy")
	}
	if GoodDiscoveryMarket(m, time.Now()) {
		t.Fatal("noisy market should not be a discovery market")
	}
}

func TestQualifyingTrade_RequiresPriceSizeAndMarket(t *testing.T) {
	cfg := DefaultConfig()
	allowed := map[string]struct{}{"0xok": {}}
	tr := Trade{
		ProxyWallet: "0x1111111111111111111111111111111111111111",
		Side:        "BUY",
		Type:        "TRADE",
		Price:       0.42,
		Size:        300,
		ConditionID: "0xok",
		Title:       "Will Arsenal FC win on 2026-05-18?",
	}
	if !QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("expected trade to qualify")
	}
	tr.Price = 0.99
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("extreme price should not qualify")
	}
	tr.Price = 0.42
	tr.ConditionID = "0xother"
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("unscanned market should not qualify")
	}
}

func TestTargetCategoryAllowed_FocusesRequestedSports(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{"nba", "NBA: Knicks vs Celtics", "basketball"},
		{"nba teams only", "Thunder vs Pacers - Game 7 Winner", "basketball"},
		{"soccer", "EPL Arsenal vs Burnley", "soccer"},
		{"soccer ucl token", "UCL: Arsenal vs Barcelona", "soccer"},
		{"soccer clubs only", "Chelsea vs PSG - Club World Cup", "soccer"},
		{"chelsea clinton is not soccer", "Will Chelsea Clinton win the 2028 Democratic presidential nomination?", "other"},
		{"brazilian election is not soccer", "Will Geraldo Alckmin win the 2026 Brazilian presidential election?", "other"},
		{"soccer fifwc slug", "Will Spain win on 2026-07-10? fifwc-esp-bel-2026-07-10", "soccer"},
		{"soccer national question only", "Will Spain win on 2026-07-10?", "soccer"},
		{"soccer national england", "Will England win on 2026-07-11?", "soccer"},
		{"soccer world cup outright rejected", "Will France win the 2026 FIFA World Cup? world-cup-winner", "other"},
		{"soccer national combo rejected", "Will France win on 2026-07-09? AND Will Norway win on 2026-07-11?", "other"},
		{"esports combo rejected", "LoL: Eintracht Frankfurt vs BIG (BO1) - Prime League Regular Season AND LoL: VfB eSports vs Eintracht Spandau (BO1)", "other"},
		{"soccer spread rejected", "Spread: KF Egnatia Rrogozhinë (-1.5)", "other"},
		{"soccer total rejected", "FC Petrocub Hînceşti vs. KF Egnatia Rrogozhinë: O/U 2.5", "other"},
		{"soccer player prop rejected", "Will Kylian Mbappe score 13+ goals during the 2026 FIFA World Cup?", "other"},
		{"soccer exact score", "Exact Score: Liverpool 2 - 1 Arsenal?", "soccer"},
		{"soccer correct score", "France vs. Morocco: Correct Score 1-0", "soccer"},
		{"lol", "LoL: T1 vs Gen.G - LCK Spring", "esports"},
		{"dota", "dota2-gl-heroic-2026-04-22", "esports"},
		{"nuclear does not contain ucl token", "US-Iran Final Nuclear Deal by August 31, 2026? us-iran-final-nuclear-deal-by-20260621201254412", "other"},
		{"tennis excluded", "Wimbledon WTA: A vs B", "other"},
	}
	for _, c := range cases {
		if got := targetCategory(c.text); got != c.want {
			t.Fatalf("%s: got %s want %s", c.name, got, c.want)
		}
		if c.want == "other" {
			if TargetCategoryAllowed(c.want, "basketball,soccer,esports") {
				t.Fatalf("%s should not be allowed", c.name)
			}
			continue
		}
		if !TargetCategoryAllowed(c.want, "basketball,soccer,esports") {
			t.Fatalf("%s should be allowed", c.name)
		}
	}
}

func TestQualifyingTrade_TargetCategoryFilter(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TargetCategories = "basketball,soccer,esports"
	allowed := map[string]struct{}{"0xok": {}}
	tr := Trade{
		ProxyWallet: "0x1111111111111111111111111111111111111111",
		Side:        "BUY",
		Type:        "TRADE",
		Price:       0.42,
		Size:        300,
		ConditionID: "0xok",
		Title:       "Wimbledon WTA: A vs B",
		Slug:        "wta-a-b-2026-07-08",
	}
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("tennis should not qualify for target follow categories")
	}
	tr.Title = "LoL: T1 vs Gen.G - LCK Spring"
	tr.Slug = "lol-lck-t1-geng-2026-07-08"
	if !QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("LoL should qualify for target follow categories")
	}
	tr.Title = "Dota 2: Virtus.pro vs 1win - Game 2 Winner"
	tr.Slug = "dota2-vp-1win-2026-07-08"
	if !QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("Dota 2 game winner should qualify for target follow categories")
	}
	tr.Title = "Game 2: Both Teams Beat Roshan?"
	tr.Slug = "dota2-vp-1win-2026-07-08"
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("Dota 2 prop should not qualify for target follow categories")
	}
	tr.Title = "Will France win the 2026 FIFA World Cup?"
	tr.Slug = "world-cup-winner"
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("long-horizon World Cup outright should not qualify for target follow categories")
	}
	tr.Title = "Exact Score: Liverpool 2 - 1 Arsenal?"
	tr.Slug = "epl-liv-ars-correct-score"
	tr.Outcome = "Yes"
	if !QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("football exact-score Yes trade should qualify for push discovery")
	}
	tr.Outcome = "No"
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("football exact-score No trade should not qualify for push discovery")
	}
	tr.Title = "Counter-Strike: maybe vs Tricksters (BO3) - CCT Europe Contenders #6 Playoffs"
	tr.Slug = "counter-strike-maybe-tricksters-2026-07-08"
	if QualifyingTrade(tr, cfg, allowed) {
		t.Fatal("unsupported esports should not qualify for strict target follow categories")
	}
}
