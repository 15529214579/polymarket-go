package feed

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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
		question string
		slug     string
		want     bool
	}{
		{"", "nba-atl-nyk-2026-04-20", true},
		{"", "nba-phi-bos-2026-04-21", true},
		{"", "nba-playoffs-who-will-win-series-pistons-vs-magic", true},
		{"Golden State Valkyries vs. Toronto Tempo", "", true},
		{"Lakers vs Celtics", "", true},
		{"Golden State Valkyries vs. Unknown Team", "", false},
		{"Will the Lakers win the 2026 NBA Finals?", "", false},
		{"", "will-the-los-angeles-lakers-win-the-2026-nba-finals", false}, // seasonal
		{"Will the Thunder win the NBA Finals?", "will-the-thunder-win-the-nba-finals", false},
		{"", "nba-finals-2026-mvp", false},
		{"", "mlb-det-bos-2026-04-20", false},
		{"", "nba-min-den-2026-04-20-spread-home-6pt5", false}, // derivative
		{"", "nba-min-den-2026-04-20-total-231pt5", false},     // derivative
		{"", "nba-atl-nyk-2026-04-20-spread-home-5pt5", false}, // derivative
	}
	for _, c := range cases {
		if got := IsBasketballMarket(Market{Question: c.question, Slug: c.slug}); got != c.want {
			t.Errorf("%s / %s: got %v want %v", c.question, c.slug, got, c.want)
		}
	}
}

func TestIsFootballMarket(t *testing.T) {
	cases := []struct {
		question string
		slug     string
		want     bool
	}{
		{"", "epl-cry-wes-2026-04-20-wes", true},
		{"", "epl-bur-mac-2026-04-22-mac", true},
		{"", "fifwc-esp-bel-2026-07-10", true},
		{"Will France win the 2026 FIFA World Cup?", "world-cup-winner", false},
		{"", "will-manchester-city-win-2025-26", false},        // seasonal
		{"", "nfl-dal-nyg-2026-09-08", false},                  // not in scope
		{"", "epl-cry-wes-2026-04-20-spread-away-2pt5", false}, // derivative
		{"", "fifwc-esp-bel-2026-07-10-total-2pt5", false},     // derivative
	}
	for _, c := range cases {
		if got := IsFootballMarket(Market{Question: c.question, Slug: c.slug}); got != c.want {
			t.Errorf("%s / %s: got %v want %v", c.question, c.slug, got, c.want)
		}
	}
}

func TestIsFootballMarket_NationalWinQuestionWithoutSlug(t *testing.T) {
	if !IsFootballMarket(Market{Question: "Will Spain win on 2026-07-10?"}) {
		t.Fatal("national-team daily soccer question should pass even without slug")
	}
	if !IsFootballMarket(Market{Question: "Will England win on 2026-07-11?"}) {
		t.Fatal("England national-team daily soccer question should pass even without slug")
	}
	if IsFootballMarket(Market{Question: "Will Lorenzo Musetti win on 2026-07-10?"}) {
		t.Fatal("player daily question should not be treated as national-team soccer")
	}
	if IsFootballMarket(Market{Question: "Will France win on 2026-07-09? AND Will Norway win on 2026-07-11?"}) {
		t.Fatal("combo/parlay national-team question should not be treated as a clean soccer moneyline")
	}
}

func TestIsFootballMarket_NationalTeamAdvanceQuestion(t *testing.T) {
	if !IsFootballMarket(Market{Question: "Spain vs. Belgium: Team to Advance", Slug: "spain-belgium-team-to-advance"}) {
		t.Fatal("national-team soccer advance question should pass")
	}
	if IsFootballMarket(Market{Question: "Boston Red Sox vs. Chicago White Sox: Team to Advance", Slug: "mlb-bos-cws-team-to-advance"}) {
		t.Fatal("non-football team advance question should not pass")
	}
}

func TestIsFootballMarket_RejectsSoccerPropsAndDerivatives(t *testing.T) {
	cases := []Market{
		{Question: "Kylian Mbappé: 1+ goals", Slug: "france-morocco-mbappe-1-goals"},
		{Question: "Ousmane Dembélé: 3+ shots", Slug: "france-morocco-dembele-3-shots"},
		{Question: "France vs. Morocco: Will the Match Go to a Penalty Shootout?", Slug: "france-morocco-penalty-shootout"},
		{Question: "Norway vs. England: Draw at halftime?", Slug: "norway-england-draw-at-halftime"},
		{Question: "Morocco to score first vs. France?", Slug: "morocco-france-score-first"},
		{Question: "France vs. Morocco: Will the Match Go to Extra Time?", Slug: "france-morocco-extra-time"},
		{Question: "France to win the second half?", Slug: "france-morocco-second-half"},
	}
	for _, m := range cases {
		if IsFootballMarket(m) {
			t.Fatalf("soccer derivative should be rejected: %s / %s", m.Question, m.Slug)
		}
	}
}

func TestIsFootballScoreMarketText(t *testing.T) {
	for _, text := range []string{
		"Exact Score: Liverpool 2 - 1 Arsenal?",
		"France vs. Morocco: Correct Score 1-0",
		"Liverpool to win 2-0",
		"Draw (1-1)",
	} {
		if !IsFootballScoreMarketText(text) {
			t.Fatalf("football score market not recognized: %q", text)
		}
	}
	for _, text := range []string{
		"LoL series score 2-0",
		"Dota 2 exact score 2-0",
		"NBA exact score 110-105",
		"Both teams to score",
	} {
		if IsFootballScoreMarketText(text) {
			t.Fatalf("non-football score market recognized: %q", text)
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
		{"dota2-gl-xtreme-2026-07-08-match-result-team2", false},
		{"dota2-gl-xtreme-2026-07-08-match-result-draw", false},
	}
	for _, c := range cases {
		if got := IsDota2Market(Market{Question: c.slug, Slug: c.slug}); got != c.want {
			t.Errorf("%s: got %v want %v", c.slug, got, c.want)
		}
	}
}

func TestIsDota2Market_RejectsMatchResultQuestionDerivatives(t *testing.T) {
	cases := []Market{
		{Question: "Xtreme Gaming to win 2-0?", Slug: "dota2-gl-xtreme-2026-07-08-match-result-team2"},
		{Question: "GamerLegion vs Xtreme Gaming: Draw (1-1)?", Slug: "dota2-gl-xtreme-2026-07-08-match-result-draw"},
		{Question: "Game 2: Both Teams Beat Roshan?", Slug: "dota2-vp-1win-2026-07-08"},
		{Question: "First Blood in Game 2?", Slug: "dota2-vp-1win-2026-07-08"},
		{Question: "Total Kills Over/Under 55.5 in Game 2?", Slug: "dota2-vp-1win-2026-07-08"},
	}
	for _, m := range cases {
		if IsDota2Market(m) {
			t.Fatalf("derivative Dota market should be rejected: %s / %s", m.Question, m.Slug)
		}
	}
}

func TestIsTennisMarket(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"wta-sabalen-osaka-2026-04-27", true},
		{"wta-bencic-baptist-2026-04-27", true},
		{"wta-sierra-pliskov-2026-04-27", true},
		{"wta-you-lu-2026-04-27", true},
		{"wta-andreev-bondar-2026-04-27", true},
		{"atp-tsitsip-aguilar-2026-04-27", true},
		{"atp-atmane-zverev-2026-04-26", true},
		{"will-lorenzo-musetti-win-the-2026-mens-french-open", false}, // seasonal
		{"nba-atl-nyk-2026-04-20", false},
	}
	for _, c := range cases {
		if got := IsTennisMarket(Market{Slug: c.slug}); got != c.want {
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
		{Slug: "will-the-lakers-win-the-2026-nba-finals"}, // excluded
		{Slug: "election-2026"},                           // excluded
		{Slug: "nba-playoffs-who-will-win-series-foo-bar"},
		{Question: "LoL: Fnatic vs G2 - LEC Regular", Slug: "lol-lec-fnc-g2-2026-04-19"}, // excluded (LEC)
		{Slug: "wta-sabalen-osaka-2026-04-27"},
		{Slug: "atp-tsitsip-aguilar-2026-04-27"},
	}
	out := FilterSports(in)
	if len(out) != 7 {
		t.Fatalf("want 7, got %d", len(out))
	}
	want := []string{
		"lol-lck-t1-geng-2026-04-19",
		"nba-atl-nyk-2026-04-20",
		"epl-cry-wes-2026-04-20-cry",
		"dota2-gl-heroic-2026-04-22",
		"nba-playoffs-who-will-win-series-foo-bar",
		"wta-sabalen-osaka-2026-04-27",
		"atp-tsitsip-aguilar-2026-04-27",
	}
	for i, m := range out {
		if m.Slug != want[i] {
			t.Errorf("pos %d: got %s want %s", i, m.Slug, want[i])
		}
	}
}

func TestFilterFollowTargets_ExcludesTennis(t *testing.T) {
	in := []Market{
		{Question: "LoL: T1 vs Gen.G - LCK Spring", Slug: "lol-lck-t1-geng-2026-04-19"},
		{Slug: "nba-atl-nyk-2026-04-20"},
		{Slug: "epl-cry-wes-2026-04-20-cry"},
		{Slug: "dota2-gl-heroic-2026-04-22"},
		{Slug: "wta-sabalen-osaka-2026-04-27"},
		{Slug: "atp-tsitsip-aguilar-2026-04-27"},
	}
	out := FilterFollowTargets(in)
	if len(out) != 4 {
		t.Fatalf("want 4, got %d", len(out))
	}
	for _, m := range out {
		if IsTennisMarket(m) {
			t.Fatalf("follow targets should exclude tennis: %s", m.Slug)
		}
	}
}

func TestFilterTradablePriceBand_RemovesSettledLikeMarkets(t *testing.T) {
	in := []Market{
		{Slug: "extreme-high-low", OutcomePricesRaw: `["0.999","0.001"]`},
		{Slug: "inside-band", OutcomePricesRaw: `["0.74","0.26"]`},
		{Slug: "boundary-band", OutcomePricesRaw: `["0.95","0.05"]`},
		{Slug: "missing-prices"},
		{Slug: "bad-prices", OutcomePricesRaw: `["bad",""]`},
	}
	out := FilterTradablePriceBand(in, 0.05, 0.95)
	got := make([]string, 0, len(out))
	for _, m := range out {
		got = append(got, m.Slug)
	}
	want := []string{"inside-band", "boundary-band", "missing-prices", "bad-prices"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestListActiveMarkets_PaginatesWithGammaLimitCap(t *testing.T) {
	requestedOffsets := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/markets" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "100" {
			t.Errorf("limit=%q, want 100", got)
		}
		offset := r.URL.Query().Get("offset")
		requestedOffsets = append(requestedOffsets, offset)
		var rows []Market
		switch offset {
		case "0":
			for i := 0; i < 100; i++ {
				rows = append(rows, Market{ID: "page1"})
			}
		case "100":
			for i := 0; i < 25; i++ {
				rows = append(rows, Market{ID: "page2"})
			}
		default:
			t.Errorf("unexpected offset %q", offset)
		}
		_ = json.NewEncoder(w).Encode(rows)
	}))
	defer srv.Close()

	c := &GammaClient{http: &http.Client{Timeout: 3 * time.Second}, base: srv.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	got, err := c.ListActiveMarkets(ctx, 500)
	if err != nil {
		t.Fatalf("ListActiveMarkets: %v", err)
	}
	if len(got) != 125 {
		t.Fatalf("got %d markets, want 125", len(got))
	}
	wantOffsets := []string{"0", "100"}
	if len(requestedOffsets) != len(wantOffsets) {
		t.Fatalf("offsets=%v, want %v", requestedOffsets, wantOffsets)
	}
	for i := range wantOffsets {
		if requestedOffsets[i] != wantOffsets[i] {
			t.Fatalf("offsets=%v, want %v", requestedOffsets, wantOffsets)
		}
	}
}

func TestListActiveMarkets_StopsAtGammaOffsetLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("offset") == "100" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error":"offset too large, use /markets/keyset for deeper pagination"}`))
			return
		}
		rows := make([]Market, 100)
		_ = json.NewEncoder(w).Encode(rows)
	}))
	defer srv.Close()

	c := &GammaClient{http: &http.Client{Timeout: 3 * time.Second}, base: srv.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	got, err := c.ListActiveMarkets(ctx, 500)
	if err != nil {
		t.Fatalf("ListActiveMarkets: %v", err)
	}
	if len(got) != 100 {
		t.Fatalf("got %d markets, want partial first page 100", len(got))
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

func TestGetByClobTokenIDs(t *testing.T) {
	wantTokens := map[string]bool{"101": true, "202": true}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/markets" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		got := r.URL.Query()["clob_token_ids"]
		sort.Strings(got)
		want := []string{"101", "202"}
		if len(got) != len(want) {
			t.Errorf("clob_token_ids count: got %v want %v", got, want)
		}
		for _, id := range got {
			if !wantTokens[id] {
				t.Errorf("unexpected clob_token_id %q", id)
			}
		}
		rows := []map[string]any{
			{
				"conditionId":   "0xabc",
				"question":      "Dota 2: A vs B - Game 1 Winner",
				"closed":        true,
				"outcomes":      `["A","B"]`,
				"outcomePrices": `["0","1"]`,
				"clobTokenIds":  `["101","102"]`,
			},
			{
				"conditionId":   "0xdef",
				"question":      "LoL: C vs D - Game 2 Winner",
				"closed":        false,
				"outcomes":      `["C","D"]`,
				"outcomePrices": `["0.45","0.55"]`,
				"clobTokenIds":  `["201","202"]`,
			},
		}
		_ = json.NewEncoder(w).Encode(rows)
	}))
	defer srv.Close()

	c := &GammaClient{http: &http.Client{Timeout: 3 * time.Second}, base: srv.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	got, err := c.GetByClobTokenIDs(ctx, []string{"101", "", "202"})
	if err != nil {
		t.Fatalf("GetByClobTokenIDs: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d markets, want 2", len(got))
	}
	if got[0].ConditionID != "0xabc" {
		t.Fatalf("first condition=%q, want 0xabc", got[0].ConditionID)
	}
}

func TestMarketEventStartTime(t *testing.T) {
	direct := Market{GameStartTime: "2026-08-02T12:30:00Z"}
	if got := direct.EventStartTime(); got.Format(time.RFC3339) != "2026-08-02T12:30:00Z" {
		t.Fatalf("direct game start=%v", got)
	}
	event := Market{Events: []MarketEvent{{StartTime: "2026-08-03T09:15:00.123Z"}}}
	if got := event.EventStartTime(); got.Format(time.RFC3339Nano) != "2026-08-03T09:15:00.123Z" {
		t.Fatalf("event game start=%v", got)
	}
	spaceOffset := Market{GameStartTime: "2026-08-01 13:15:00+00"}
	if got := spaceOffset.EventStartTime(); got.Format(time.RFC3339) != "2026-08-01T13:15:00Z" {
		t.Fatalf("space-offset game start=%v", got)
	}
}

func TestGetCLOBEventStart(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/clob-markets/0xabc" {
			t.Fatalf("path=%q", r.URL.Path)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"gst":"2026-08-04T11:00:00Z","fd":{"r":0.03}}`)),
			Header:     make(http.Header),
		}, nil
	})}
	c := &GammaClient{http: client, clobBase: "https://clob.test"}
	got, err := c.GetCLOBEventStart(context.Background(), "0xabc")
	if err != nil {
		t.Fatal(err)
	}
	if got.Format(time.RFC3339) != "2026-08-04T11:00:00Z" {
		t.Fatalf("game start=%v", got)
	}
	info, err := c.GetCLOBMarketInfo(context.Background(), "0xabc")
	if err != nil {
		t.Fatal(err)
	}
	if !info.FeeRateKnown || info.TakerFeeRate != 0.03 {
		t.Fatalf("fee info=%+v", info)
	}
}
