package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/whale"
)

func TestWhaleEventKey_GroupsRelatedMarkets(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "dota game one",
			in:   "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
			want: "dota 2: rekonix vs mouz",
		},
		{
			name: "dota bo2",
			in:   "Dota 2: REKONIX vs MOUZ (BO2) - Esports World Cup Group C",
			want: "dota 2: rekonix vs mouz",
		},
		{
			name: "lol bo3",
			in:   "LoL: T1 vs Gen.G (BO3) - LCK Summer",
			want: "lol: t1 vs gen.g",
		},
		{
			name: "soccer total",
			in:   "FC Flora vs. SK Iberia 1999: O/U 2.5",
			want: "fc flora vs. sk iberia 1999",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := whaleEventKey(tc.in); got != tc.want {
				t.Fatalf("whaleEventKey(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestIsTargetFollowMarket_IncludesFIFWCSoccer(t *testing.T) {
	if !isTargetFollowMarket("Will Spain win on 2026-07-10?", "fifwc-esp-bel-2026-07-10") {
		t.Fatal("FIFWC soccer moneyline should be a target follow market")
	}
}

func TestIsTargetFollowMarket_IncludesNationalSoccerWithoutSlug(t *testing.T) {
	if !isTargetFollowMarket("Will Spain win on 2026-07-10?", "") {
		t.Fatal("national soccer moneyline should be a target follow market even when activity has no slug")
	}
	if !isTargetFollowMarket("Will England win on 2026-07-11?", "") {
		t.Fatal("England national soccer moneyline should be a target follow market even when activity has no slug")
	}
	if isTargetFollowMarket("Will France win on 2026-07-09? AND Will Norway win on 2026-07-11?", "") {
		t.Fatal("combo/parlay national soccer market should not be a target follow market")
	}
}

func TestIsTargetFollowMarket_ExcludesComboEsports(t *testing.T) {
	q := "LoL: Eintracht Frankfurt vs BIG (BO1) - Prime League Regular Season AND LoL: VfB eSports vs Eintracht Spandau (BO1)"
	if isTargetFollowMarket(q, "") {
		t.Fatal("combo/parlay esports market should not be a target follow market")
	}
}

func TestIsTargetFollowMarket_ExcludesTennis(t *testing.T) {
	if isTargetFollowMarket("Wimbledon WTA: A vs B", "wta-a-b-2026-07-08") {
		t.Fatal("tennis should not be a target follow market")
	}
}

func TestTargetFollowMarketDecision_ExplainsFilteredMarkets(t *testing.T) {
	if ok, reason := targetFollowMarketDecision("Norway vs. England: O/U 2.5", ""); ok || reason != "derivative_filtered" {
		t.Fatalf("O/U decision=%v/%q, want derivative_filtered", ok, reason)
	}
	if ok, reason := targetFollowMarketDecision("Spread: Toronto Blue Jays (-1.5)", ""); ok || reason != "derivative_filtered" {
		t.Fatalf("spread decision=%v/%q, want derivative_filtered", ok, reason)
	}
	if ok, reason := targetFollowMarketDecision("Wimbledon ATP: Taylor Fritz vs Alexander Zverev", "atp-fritz-zverev"); ok || reason != "category_filtered" {
		t.Fatalf("tennis decision=%v/%q, want category_filtered", ok, reason)
	}
	if ok, reason := targetFollowMarketDecision("Will England win on 2026-07-11?", ""); !ok || reason != "" {
		t.Fatalf("soccer moneyline decision=%v/%q, want allowed", ok, reason)
	}
	for _, q := range []string{
		"Kylian Mbappé: 1+ goals",
		"Ousmane Dembélé: 3+ shots",
		"Will France vs. Morocco end in a draw?",
		"France vs. Morocco: Will the Match Go to a Penalty Shootout?",
		"Norway vs. England: Draw at halftime?",
		"Morocco to score first vs. France?",
		"France vs. Morocco: Will the Match Go to Extra Time?",
		"France to win the second half?",
	} {
		if ok, reason := targetFollowMarketDecision(q, "fifa-world-cup"); ok || reason != "derivative_filtered" {
			t.Fatalf("soccer derivative %q decision=%v/%q, want derivative_filtered", q, ok, reason)
		}
	}
}

func TestWhaleMarketDecision_LeaderboardKeepsTargetCategories(t *testing.T) {
	if ok, reason := whaleMarketDecision("Boston Red Sox vs. Chicago White Sox", "mlb-bos-cws-2026-07-09", "", "leaderboard_watch"); ok || reason != "category_filtered" {
		t.Fatalf("leaderboard MLB moneyline decision=%v/%q, want category_filtered", ok, reason)
	}
	if ok, reason := whaleMarketDecision("Spread: Milwaukee Brewers (-1.5)", "mlb-mil-stl-2026-07-08", "", "leaderboard_watch"); ok || reason != "derivative_filtered" {
		t.Fatalf("leaderboard spread decision=%v/%q, want derivative_filtered", ok, reason)
	}
	if ok, reason := whaleMarketDecision("Will France win the 2026 FIFA World Cup?", "world-cup-winner", "", "leaderboard_watch"); ok || reason != "outright_filtered" {
		t.Fatalf("leaderboard world cup outright decision=%v/%q, want outright_filtered", ok, reason)
	}
	if ok, reason := whaleMarketDecision("Spain vs. Belgium: Team to Advance", "spain-belgium-team-to-advance", "", "leaderboard_watch"); !ok || reason != "" {
		t.Fatalf("leaderboard soccer decision=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := whaleMarketDecision("Golden State Valkyries vs. Toronto Tempo", "", "", "leaderboard_watch"); !ok || reason != "" {
		t.Fatalf("leaderboard WNBA decision=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := whaleMarketDecision("Lakers vs Celtics", "", "", "watch"); !ok || reason != "" {
		t.Fatalf("NBA team-name decision=%v/%q, want allowed", ok, reason)
	}
}

func TestWhaleMarketDecision_AllowsFootballScoreYesOnly(t *testing.T) {
	q := "Exact Score: Liverpool 2 - 1 Arsenal?"
	if ok, reason := whaleMarketDecision(q, "epl-liv-ars-correct-score", "Yes", "football_score_push"); !ok || reason != "" {
		t.Fatalf("football score Yes decision=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := whaleMarketDecision(q, "epl-liv-ars-correct-score", "No", "football_score_push"); ok || reason != "football_score_no_filtered" {
		t.Fatalf("football score No decision=%v/%q, want filtered", ok, reason)
	}
	if ok, reason := targetFollowMarketDecision(q, "epl-liv-ars-correct-score"); ok || reason != "derivative_filtered" {
		t.Fatalf("automatic follow decision=%v/%q, want derivative_filtered", ok, reason)
	}
}

func TestCopytradeMarketDecision_FootballScoreIsPaperOptIn(t *testing.T) {
	q := "Exact Score: Liverpool 2 - 1 Arsenal?"
	if ok, reason := copytradeMarketDecision(q, "epl-liv-ars-correct-score", "Yes", false); ok || reason != "derivative_filtered" {
		t.Fatalf("disabled score decision=%v/%q, want derivative_filtered", ok, reason)
	}
	if ok, reason := copytradeMarketDecision(q, "epl-liv-ars-correct-score", "Yes", true); !ok || reason != "" {
		t.Fatalf("enabled score decision=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := copytradeMarketDecision(q, "epl-liv-ars-correct-score", "No", true); ok || reason != "football_score_no_filtered" {
		t.Fatalf("score No decision=%v/%q, want filtered", ok, reason)
	}
	if ok, reason := copytradeMarketDecision("Will Spain win on 2026-08-01?", "fifwc-esp-fra-2026-08-01", "Yes", true); !ok || reason != "" {
		t.Fatalf("soccer moneyline decision=%v/%q, want allowed", ok, reason)
	}
}

func TestCopytradeFootballScorePriceAndSize(t *testing.T) {
	if got := copytradeEntryPriceFloor(true, true); got != 0.01 {
		t.Fatalf("score price floor=%v, want 0.01", got)
	}
	if got := copytradeEntryPriceFloor(true, false); got != 0.05 {
		t.Fatalf("disabled score price floor=%v, want 0.05", got)
	}
	if got := copytradeMarketSize(20, true, true, 20); got != 20 {
		t.Fatalf("score size=%v, want 20", got)
	}
	if got := copytradeMarketSize(10, false, true, 5); got != 10 {
		t.Fatalf("moneyline size=%v, want 10", got)
	}
}

func TestCopytradeFootballScoreUsesDedicatedFileTier(t *testing.T) {
	meta := walletFileMeta{List: "football_score_push", Tier: "B"}
	if got := copytradeTierForMarket("BOT", meta, true); got != "B" {
		t.Fatalf("football score tier=%q, want B", got)
	}
	if got := copytradeTierForMarket("BOT", meta, false); got != "BOT" {
		t.Fatalf("ordinary market tier=%q, want BOT", got)
	}
	allowed, reason := copytradeAutoAllowedForMarket("reject-bot", false, true, true, true, "B")
	if !allowed || reason != "football_score_b" {
		t.Fatalf("score auto allowed=%v reason=%q", allowed, reason)
	}
	if allowed, _ := copytradeAutoAllowedForMarket("reject-bot", true, true, true, true, "B"); allowed {
		t.Fatal("dedicated score tier must not bypass live-trading gate")
	}
}

func TestPaperCollectionCannotAffectLive(t *testing.T) {
	if !paperCollectionEnabled(true, "copytrade", false) {
		t.Fatal("broad collection should be enabled for paper copytrade")
	}
	if paperCollectionEnabled(true, "copytrade", true) {
		t.Fatal("broad collection must be disabled for live trading")
	}
	if paperCollectionEnabled(true, "whale", false) {
		t.Fatal("broad collection must be limited to copytrade mode")
	}
}

func TestCopytradeTierAllowed(t *testing.T) {
	if !copytradeTierAllowed("B", "B") || copytradeTierAllowed("C", "B") {
		t.Fatal("B core gate should allow A/B and reject C")
	}
	if !copytradeTierAllowed("d", "") {
		t.Fatal("empty tier gate should allow every tier")
	}
}

func TestCopytradeCollectionMarketDecisionIsBroadButKeepsScoreNoFilter(t *testing.T) {
	if ok, reason := copytradeCollectionMarketDecision("Norway vs England: O/U 2.5", "", "Over", true); !ok || reason != "" {
		t.Fatalf("paper collection derivative=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := copytradeCollectionMarketDecision("Wimbledon ATP: A vs B", "atp-a-b", "A", true); !ok || reason != "" {
		t.Fatalf("paper collection tennis=%v/%q, want allowed", ok, reason)
	}
	q := "Exact Score: Liverpool 2 - 1 Arsenal?"
	if ok, reason := copytradeCollectionMarketDecision(q, "epl-liv-ars-correct-score", "No", true); ok || reason != "football_score_no_filtered" {
		t.Fatalf("paper collection score No=%v/%q, want filtered", ok, reason)
	}
}

func TestCopytradeWalletSize_PaperUsesConfiguredFixedSize(t *testing.T) {
	for _, tc := range []struct {
		tier, action string
		score        float64
	}{
		{tier: "A"},
		{tier: "B"},
		{tier: "A", action: "auto-small", score: 95},
	} {
		if got := copytradeWalletSize(20, false, tc.tier, tc.action, tc.score); got != 20 {
			t.Fatalf("tier=%s action=%s size=%v, want 20", tc.tier, tc.action, got)
		}
	}
}

func TestCopytradeWalletSize_LiveKeepsTierSizing(t *testing.T) {
	if got := copytradeWalletSize(20, true, "A", "", 0); got != 10 {
		t.Fatalf("live tier A size=%v, want 10", got)
	}
	if got := copytradeWalletSize(20, true, "B", "", 0); got != 5 {
		t.Fatalf("live tier B size=%v, want 5", got)
	}
	if got := copytradeWalletSize(20, true, "A", "auto-small", 90); got != 20 {
		t.Fatalf("live auto-small size=%v, want 20", got)
	}
}

func TestWhalePriceDecision_FiltersExtremeBuyPrices(t *testing.T) {
	for _, price := range []float64{0.001, 0.049, 0.951, 0.999} {
		if ok, reason := whalePriceDecision("BUY", price, false); ok || reason != "price_filtered" {
			t.Fatalf("BUY price %.3f decision=%v/%q, want price_filtered", price, ok, reason)
		}
	}
	for _, price := range []float64{0.05, 0.5, 0.95} {
		if ok, reason := whalePriceDecision("BUY", price, false); !ok || reason != "" {
			t.Fatalf("BUY price %.3f decision=%v/%q, want allowed", price, ok, reason)
		}
	}
	if ok, reason := whalePriceDecision("SELL", 0.999, false); ok || reason != "settlement_sell_filtered" {
		t.Fatalf("SELL settlement price decision=%v/%q, want settlement_sell_filtered", ok, reason)
	}
	if ok, reason := whalePriceDecision("SELL", 0.75, false); !ok || reason != "" {
		t.Fatalf("SELL normal exit price decision=%v/%q, want allowed", ok, reason)
	}
	if ok, reason := whalePriceDecision("BUY", 0.02, true); !ok || reason != "" {
		t.Fatalf("football score BUY decision=%v/%q, want low-price score allowed", ok, reason)
	}
}

func TestWhaleEventKeyForAlert_SeparatesFootballScores(t *testing.T) {
	a := whale.AlertEvent{Question: "Exact Score: Liverpool 1 - 0 Arsenal?", Outcome: "Yes", ConditionID: "0xscore10"}
	b := whale.AlertEvent{Question: "Exact Score: Liverpool 1 - 1 Arsenal?", Outcome: "Yes", ConditionID: "0xscore11"}
	if whaleEventKeyForAlert(a) == whaleEventKeyForAlert(b) {
		t.Fatal("different football score markets must not share an event cooldown key")
	}
}

func TestCopytradeExposureEventKey_GroupsFootballScoreBasket(t *testing.T) {
	a := whale.AlertEvent{Question: "Exact Score: Liverpool 1 - 0 Arsenal?", Slug: "liverpool-arsenal-2026-08-01", Outcome: "Yes", ConditionID: "0xscore10"}
	b := whale.AlertEvent{Question: "Exact Score: Liverpool 1 - 1 Arsenal?", Slug: "liverpool-arsenal-2026-08-01", Outcome: "Yes", ConditionID: "0xscore11"}
	if copytradeExposureEventKey(a) != copytradeExposureEventKey(b) {
		t.Fatalf("same match score baskets differ: %q / %q", copytradeExposureEventKey(a), copytradeExposureEventKey(b))
	}
}

func TestCopytradeAutoAllowedForAction_PaperPromptOptIn(t *testing.T) {
	if ok, action := copytradeAutoAllowedForAction("prompt", false, false); ok || action != "prompt" {
		t.Fatalf("default prompt allowed=%v action=%q, want blocked prompt", ok, action)
	}
	if ok, action := copytradeAutoAllowedForAction("prompt", false, true); !ok || action != "prompt" {
		t.Fatalf("paper prompt allowed=%v action=%q, want allowed prompt", ok, action)
	}
	if ok, action := copytradeAutoAllowedForAction("prompt", true, true); ok || action != "prompt" {
		t.Fatalf("live prompt allowed=%v action=%q, want blocked prompt", ok, action)
	}
	if ok, action := copytradeAutoAllowedForAction("auto-small", true, false); !ok || action != "auto-small" {
		t.Fatalf("auto-small allowed=%v action=%q, want allowed auto-small", ok, action)
	}
}

func TestLoadWhaleNegativeEdgeBlocks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	data := strings.Join([]string{
		`{"wallet":"0x1111111111111111111111111111111111111111","horizon_sec":900,"delta_pp":-2.0}`,
		`{"wallet":"0x1111111111111111111111111111111111111111","horizon_sec":900,"delta_pp":-1.0}`,
		`{"wallet":"0x2222222222222222222222222222222222222222","horizon_sec":3600,"delta_pp":-8.0}`,
		`{"wallet":"0x3333333333333333333333333333333333333333","horizon_sec":900,"delta_pp":-0.5}`,
		`{"wallet":"0x3333333333333333333333333333333333333333","horizon_sec":3600,"delta_pp":-4.0}`,
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	blocks, err := loadWhaleNegativeEdgeBlocks(path, whaleEdgeBlockConfig{
		Min15mSamples: 2,
		Max15mAvgPP:   -1,
		Min1hSamples:  1,
		Max1hAvgPP:    -5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := blocks["0x1111111111111111111111111111111111111111"]; !strings.Contains(got, "15m edge -1.50pp") {
		t.Fatalf("first block=%q, want 15m avg", got)
	}
	if got := blocks["0x2222222222222222222222222222222222222222"]; !strings.Contains(got, "1h edge -8.00pp") {
		t.Fatalf("second block=%q, want 1h avg", got)
	}
	if _, ok := blocks["0x3333333333333333333333333333333333333333"]; ok {
		t.Fatalf("third wallet should not be blocked: %#v", blocks)
	}
}

func TestWhaleEdgeBlockCacheRefreshes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	cache := newWhaleEdgeBlockCache(path, whaleEdgeBlockConfig{
		Min15mSamples: 1,
		Max15mAvgPP:   -1,
		Refresh:       time.Nanosecond,
	})
	wallet := "0x1111111111111111111111111111111111111111"
	if reason, ok := cache.Reason(wallet); ok || reason != "" {
		t.Fatalf("empty cache reason=%q ok=%v, want no block", reason, ok)
	}
	if err := os.WriteFile(path, []byte(`{"wallet":"`+wallet+`","horizon_sec":900,"delta_pp":-2.0}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reason, ok := cache.Reason(wallet)
	if !ok || !strings.Contains(reason, "15m edge -2.00pp") {
		t.Fatalf("refresh reason=%q ok=%v, want 15m block", reason, ok)
	}
}

func TestLoadWhaleEdgeSignalsIdentifiesHotWallets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	hot := "0x1111111111111111111111111111111111111111"
	cold := "0x2222222222222222222222222222222222222222"
	data := strings.Join([]string{
		`{"wallet":"` + hot + `","horizon_sec":0,"delta_pp":-4.0}`,
		`{"wallet":"` + hot + `","horizon_sec":300,"delta_pp":3.0}`,
		`{"wallet":"` + hot + `","horizon_sec":900,"delta_pp":4.0}`,
		`{"wallet":"` + hot + `","horizon_sec":1800,"delta_pp":5.0}`,
		`{"wallet":"` + hot + `","horizon_sec":3600,"delta_pp":6.0}`,
		`{"wallet":"` + cold + `","horizon_sec":300,"delta_pp":5.0}`,
		`{"wallet":"` + cold + `","horizon_sec":900,"delta_pp":-1.0}`,
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	blocks, hotWallets, err := loadWhaleEdgeSignals(path, whaleEdgeBlockConfig{
		Min15mSamples:  2,
		Max15mAvgPP:    -1,
		Min1hSamples:   1,
		Max1hAvgPP:     -5,
		HotMinSamples:  2,
		HotMinAvgPP:    2,
		HotMinWinRate:  60,
		HotMin5mAvgPP:  0.5,
		HotMin15mAvgPP: 0,
		HotMax1hNegPP:  -5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, blocked := blocks[hot]; blocked {
		t.Fatalf("hot wallet should not be blocked: %#v", blocks)
	}
	if got := hotWallets[hot]; !strings.Contains(got, "avg 4.50pp") || !strings.Contains(got, "over 4 samples") {
		t.Fatalf("hot reason=%q, want avg over four positive horizons", got)
	}
	if _, ok := hotWallets[cold]; ok {
		t.Fatalf("cold wallet should not be hot: %#v", hotWallets)
	}
}

func TestWhaleEdgeBlockCacheNegativeEdgeWinsOverHot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	wallet := "0x1111111111111111111111111111111111111111"
	data := strings.Join([]string{
		`{"wallet":"` + wallet + `","horizon_sec":300,"delta_pp":5.0}`,
		`{"wallet":"` + wallet + `","horizon_sec":900,"delta_pp":5.0}`,
		`{"wallet":"` + wallet + `","horizon_sec":3600,"delta_pp":-8.0}`,
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	cache := newWhaleEdgeBlockCache(path, whaleEdgeBlockConfig{
		Min1hSamples:   1,
		Max1hAvgPP:     -5,
		HotMinSamples:  2,
		HotMinAvgPP:    1,
		HotMinWinRate:  50,
		HotMin5mAvgPP:  0.5,
		HotMin15mAvgPP: 0,
		HotMax1hNegPP:  -5,
	})
	if reason, ok := cache.Reason(wallet); !ok || !strings.Contains(reason, "1h edge -8.00pp") {
		t.Fatalf("block reason=%q ok=%v, want negative 1h block", reason, ok)
	}
	if reason, ok := cache.HotReason(wallet); ok || reason != "" {
		t.Fatalf("negative-edge wallet should not be hot reason=%q ok=%v", reason, ok)
	}
}

func TestWhaleModeDisablesNonWhaleStrategies(t *testing.T) {
	if momentumSignalsEnabledForMode("whale") {
		t.Fatal("whale mode should not run generic momentum signals")
	}
	if lotteryScannerEnabledForMode("whale", true) {
		t.Fatal("whale mode should not run lottery scanner")
	}
	if !momentumSignalsEnabledForMode("auto") {
		t.Fatal("auto mode should keep generic momentum signals")
	}
	if !lotteryScannerEnabledForMode("auto", true) {
		t.Fatal("auto mode should keep lottery scanner when enabled")
	}
}

func TestLoadWhaleCooldownSeeds_UsesRecentBuyAlertsOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "whale_trades.jsonl")
	now := time.Now()
	lines := []string{
		`{"ts":"` + now.Add(-time.Hour).Format(time.RFC3339) + `","wallet":"0xABC","side":"BUY","market":"Dota 2: REKONIX vs MOUZ - Game 1 Winner","asset_id":"asset-1","action":"alert","size":1000}`,
		`{"ts":"` + now.Add(-30*time.Minute).Format(time.RFC3339) + `","wallet":"0xABC","side":"BUY","market":"Dota 2: REKONIX vs MOUZ (BO2) - Esports World Cup","asset_id":"asset-2","action":"event_cooldown","size":1000}`,
		`{"ts":"` + now.Add(-20*time.Minute).Format(time.RFC3339) + `","wallet":"0xABC","side":"SELL","market":"Dota 2: REKONIX vs MOUZ - Game 1 Winner","asset_id":"asset-1","action":"alert","size":1000}`,
		`{"ts":"` + now.Add(-8*time.Hour).Format(time.RFC3339) + `","wallet":"0xDEF","side":"BUY","market":"LoL: T1 vs Gen.G (BO3) - LCK Summer","asset_id":"asset-3","action":"alert","size":1000}`,
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	repeat, events, repeatN, eventN := loadWhaleCooldownSeeds(path, 3*time.Hour, 6*time.Hour)
	if repeatN != 1 || eventN != 1 {
		t.Fatalf("seed counts repeat=%d event=%d, want 1/1", repeatN, eventN)
	}
	if _, ok := repeat["0xabc|asset-1"]; !ok {
		t.Fatalf("missing repeat seed for asset-1: %#v", repeat)
	}
	if _, ok := events["0xabc|dota 2: rekonix vs mouz"]; !ok {
		t.Fatalf("missing event seed: %#v", events)
	}
}

func TestLoadWhaleSeenTradeIDs_NormalizesAndSkipsBadRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "whale_trades.jsonl")
	lines := []string{
		`{"wallet":"0xABC","trade_id":"0xTRADE"}`,
		`{"wallet":"0xABC","trade_id":"0xSKIP","action":"skip","reason":"category_filtered"}`,
		`{"wallet":"0xabc","trade_id":""}`,
		`{"wallet":"","trade_id":"0xmissingwallet"}`,
		`not-json`,
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	seen := loadWhaleSeenTradeIDs(path)
	if len(seen) != 1 {
		t.Fatalf("seen len=%d, want 1: %#v", len(seen), seen)
	}
	if _, ok := seen["0xabc|0xtrade"]; !ok {
		t.Fatalf("missing normalized trade key: %#v", seen)
	}
	if _, ok := seen["0xabc|0xskip"]; ok {
		t.Fatalf("skip rows should not seed seen trades: %#v", seen)
	}
}

func TestWhaleConfirmGate_RequiresDistinctWalletConfirmation(t *testing.T) {
	gate := newWhaleConfirmGate(30*time.Minute, 2, 5000, 0.02)
	base := whale.AlertEvent{
		Question:  "Dota 2: REKONIX vs MOUZ - Game 2 Winner",
		Outcome:   "MOUZ",
		Price:     0.67,
		Notional:  900,
		Timestamp: time.Unix(1000, 0),
		Wallet:    "0x1111111111111111111111111111111111111111",
	}
	first := gate.Observe(base, "watch")
	if first.Ready {
		t.Fatalf("first wallet should wait for confirmation: %#v", first)
	}
	if first.Wallets != 1 || first.Need != 2 {
		t.Fatalf("first decision wallets=%d need=%d, want 1/2", first.Wallets, first.Need)
	}

	base.Wallet = "0x2222222222222222222222222222222222222222"
	base.Price = 0.574
	base.Notional = 1500
	base.Timestamp = base.Timestamp.Add(4 * time.Minute)
	second := gate.Observe(base, "watch")
	if !second.Ready {
		t.Fatalf("second distinct wallet should confirm: %#v", second)
	}
	if !strings.Contains(second.Reason, "confirmed_wallets:2/2") {
		t.Fatalf("unexpected confirmation reason: %s", second.Reason)
	}
}

func TestWhaleConfirmGate_BypassesVeryLargeOrder(t *testing.T) {
	gate := newWhaleConfirmGate(30*time.Minute, 2, 5000, 0.02)
	decision := gate.Observe(whale.AlertEvent{
		Question:  "LoL: T1 vs Gen.G (BO3) - LCK Summer",
		Outcome:   "T1",
		Price:     0.61,
		Notional:  7500,
		Timestamp: time.Unix(1000, 0),
		Wallet:    "0x1111111111111111111111111111111111111111",
	}, "target")
	if !decision.Ready {
		t.Fatalf("very large order should bypass confirmation: %#v", decision)
	}
	if !strings.Contains(decision.Reason, "bypass_notional") {
		t.Fatalf("unexpected bypass reason: %s", decision.Reason)
	}
}

func TestWhaleConfirmGate_BypassesNearThresholdRounding(t *testing.T) {
	gate := newWhaleConfirmGate(30*time.Minute, 2, 5000, 0.02)
	decision := gate.Observe(whale.AlertEvent{
		Question:  "Dota 2: Nigma Galaxy vs Aurora - Game 2 Winner",
		Outcome:   "Aurora",
		Price:     0.63,
		Notional:  4999.999998,
		Timestamp: time.Unix(1000, 0),
		Wallet:    "0x1111111111111111111111111111111111111111",
	}, "target")
	if !decision.Ready {
		t.Fatalf("near-$5000 order should bypass confirmation despite floating point dust: %#v", decision)
	}
	if !strings.Contains(decision.Reason, "bypass_notional") {
		t.Fatalf("unexpected bypass reason: %s", decision.Reason)
	}
}

func TestWhaleConfirmDecisionNote(t *testing.T) {
	cases := []struct {
		name string
		in   whaleConfirmDecision
		want string
	}{
		{
			name: "large bypass",
			in:   whaleConfirmDecision{Ready: true, Reason: "bypass_notional:$7500", Wallets: 1, Need: 2},
			want: "gate 5k+ bypass",
		},
		{
			name: "multi wallet",
			in:   whaleConfirmDecision{Ready: true, Reason: "confirmed_wallets:2/2 event=dota", Wallets: 2, Need: 2},
			want: "gate 2-wallet confirmed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := whaleConfirmDecisionNote(tc.in); got != tc.want {
				t.Fatalf("whaleConfirmDecisionNote=%q, want %q", got, tc.want)
			}
		})
	}
}

func TestFormatWhalePromptContext_IncludesGateAndWalletQuality(t *testing.T) {
	got := formatWhalePromptContext(whale.AlertEvent{
		Notional:    6200,
		SizeUnits:   8000,
		TotalShares: 12000,
		AvgPrice:    0.4123,
	}, walletFileMeta{
		List:            "scout",
		SmartMoneyScore: 100,
		BotScore:        44.7,
	}, "gate 5k+ bypass")
	for _, want := range []string{"$6200", "8000 shares", "gate 5k+ bypass", "scout observe", "S100/B45", "pos 12000 @ 0.4123"} {
		if !strings.Contains(got, want) {
			t.Fatalf("context %q missing %q", got, want)
		}
	}
}

func TestWhaleNotionalAtLeast_CentToleranceOnly(t *testing.T) {
	if !whaleNotionalAtLeast(4999.999998, 5000) {
		t.Fatal("expected sub-cent floating point dust to count toward threshold")
	}
	if whaleNotionalAtLeast(4999.98, 5000) {
		t.Fatal("expected more than one cent below threshold to stay below threshold")
	}
}

func TestWhaleRepeatBypassesCooldown_DisabledByDefault(t *testing.T) {
	if whaleRepeatBypassesCooldown(50000, 0) {
		t.Fatal("zero repeat bypass threshold should keep cooldown active")
	}
	if !whaleRepeatBypassesCooldown(5000, 5000) {
		t.Fatal("explicit repeat bypass threshold should allow large orders through")
	}
}

func TestParseWhaleListMinUSD(t *testing.T) {
	got := parseWhaleListMinUSD("core=1000,sports=1500, watch = 500, bad, scout=-1")
	if got["core"] != 1000 || got["sports"] != 1500 || got["watch"] != 500 {
		t.Fatalf("unexpected parsed thresholds: %#v", got)
	}
	if _, ok := got["scout"]; ok {
		t.Fatalf("negative threshold should be ignored: %#v", got)
	}
}
