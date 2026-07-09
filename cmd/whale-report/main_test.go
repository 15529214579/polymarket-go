package main

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func TestEventKey_GroupsSameEsportsMatch(t *testing.T) {
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
			name: "counter strike bo3",
			in:   "Counter-Strike: Isurus vs UNO MILLE (BO3) - Thunderpick World Championship",
			want: "counter-strike: isurus vs uno mille",
		},
		{
			name: "soccer prop",
			in:   "FC Flora vs. SK Iberia 1999: O/U 2.5",
			want: "fc flora vs. sk iberia 1999",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := eventKey(tc.in); got != tc.want {
				t.Fatalf("eventKey(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestEventCappedResults_KeepsFirstPerWalletEvent(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	results := []*signalResult{
		{
			Buy: whaleTrade{
				Wallet: "0xabc", Market: "Dota 2: REKONIX vs MOUZ - Game 2 Winner",
				Time: base.Add(2 * time.Minute),
			},
			ExitSource: "mark",
		},
		{
			Buy: whaleTrade{
				Wallet: "0xabc", Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
				Time: base,
			},
			ExitSource: "mark",
		},
		{
			Buy: whaleTrade{
				Wallet: "0xabc", Market: "Dota 2: REKONIX vs MOUZ (BO2) - Esports World Cup",
				Time: base.Add(time.Minute),
			},
			ExitSource: "mark",
		},
		{
			Buy: whaleTrade{
				Wallet: "0xdef", Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
				Time: base.Add(30 * time.Second),
			},
			ExitSource: "mark",
		},
	}

	got := eventCappedResults(results)
	if len(got) != 2 {
		t.Fatalf("eventCappedResults len = %d, want 2", len(got))
	}
	if got[0].Buy.Wallet != "0xabc" || got[0].Buy.Market != "Dota 2: REKONIX vs MOUZ - Game 1 Winner" {
		t.Fatalf("first capped result = %#v", got[0].Buy)
	}
	if got[1].Buy.Wallet != "0xdef" {
		t.Fatalf("second capped wallet = %s, want 0xdef", got[1].Buy.Wallet)
	}
}

func TestBuildSignalResults_SuppressingActionWinsOverAlertDuplicate(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xtrade",
			AssetID: "asset-1", Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
			Price: 0.6, Size: 1000, Time: base,
		},
		{
			Wallet: "0xabc", Side: "BUY", Action: "event_cooldown", TradeID: "0xtrade",
			AssetID: "asset-1", Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
			Price: 0.6, Size: 1000, Time: base,
		},
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xkeep",
			AssetID: "asset-2", Market: "Dota 2: L1ga Team vs PlayTime - Game 2 Winner",
			Price: 0.5, Size: 1000, Time: base.Add(time.Minute),
		},
	}
	actions := parseActionSet("alert,followed")

	eval := buildSignalResults(trades, nil, actions, evalOptions{StakeUSD: 10, MinNotional: 100})
	if eval.RawBuys != 2 {
		t.Fatalf("RawBuys=%d, want 2", eval.RawBuys)
	}
	if eval.SuppressedRepeats != 1 {
		t.Fatalf("SuppressedRepeats=%d, want 1", eval.SuppressedRepeats)
	}
	if eval.LoggedEventCooldowns != 1 {
		t.Fatalf("LoggedEventCooldowns=%d, want 1", eval.LoggedEventCooldowns)
	}
	if len(eval.SuppressedByWallet) != 1 || eval.SuppressedByWallet[0].Total != 1 || eval.SuppressedByWallet[0].EventCooldown != 1 {
		t.Fatalf("SuppressedByWallet=%#v, want one event cooldown", eval.SuppressedByWallet)
	}
	if len(eval.SuppressedByEvent) != 1 || eval.SuppressedByEvent[0].Event != "dota 2: rekonix vs mouz" {
		t.Fatalf("SuppressedByEvent=%#v, want rekonix event", eval.SuppressedByEvent)
	}
	if eval.DuplicateBuys != 0 {
		t.Fatalf("DuplicateBuys=%d, want 0", eval.DuplicateBuys)
	}
	if len(eval.Results) != 1 {
		t.Fatalf("Results len=%d, want 1", len(eval.Results))
	}
	if eval.Results[0].Buy.TradeID != "0xkeep" {
		t.Fatalf("kept trade=%s, want 0xkeep", eval.Results[0].Buy.TradeID)
	}
}

func TestBuildSignalResults_CountsIgnoredDuplicateAlerts(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xtrade",
			AssetID: "asset-1", Market: "Dota 2: Nigma Galaxy vs Aurora - Game 2 Winner",
			Price: 0.63, Size: 1000, Time: base,
		},
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xtrade",
			AssetID: "asset-1", Market: "Dota 2: Nigma Galaxy vs Aurora - Game 2 Winner",
			Price: 0.63, Size: 1000, Time: base,
		},
	}
	actions := parseActionSet("alert,followed")

	eval := buildSignalResults(trades, nil, actions, evalOptions{StakeUSD: 10, MinNotional: 100})
	if eval.RawBuys != 1 {
		t.Fatalf("RawBuys=%d, want 1", eval.RawBuys)
	}
	if eval.DuplicateBuys != 1 {
		t.Fatalf("DuplicateBuys=%d, want 1", eval.DuplicateBuys)
	}
	if len(eval.Results) != 1 {
		t.Fatalf("Results len=%d, want 1", len(eval.Results))
	}
}

func TestBuildSignalResults_UsesListMinimumNotional(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet: "0xsports", List: "sports", Side: "BUY", Action: "alert", TradeID: "0xsportssmall",
			AssetID: "asset-sports-small", Market: "Will Morocco win on 2026-07-09?",
			Price: 0.5, Size: 1000, Time: base,
		},
		{
			Wallet: "0xsports", List: "sports", Side: "BUY", Action: "alert", TradeID: "0xsportslarge",
			AssetID: "asset-sports-large", Market: "Will France win on 2026-07-09?",
			Price: 0.5, Size: 1500, Time: base.Add(time.Minute),
		},
		{
			Wallet: "0xwatch", List: "watch", Side: "BUY", Action: "alert", TradeID: "0xwatch",
			AssetID: "asset-watch", Market: "Dota 2: A vs B - Game 1 Winner",
			Price: 0.5, Size: 750, Time: base.Add(2 * time.Minute),
		},
		{
			Wallet: "0xunknown", Side: "BUY", Action: "alert", TradeID: "0xunknown",
			AssetID: "asset-unknown", Market: "LoL: A vs B (BO1)",
			Price: 0.5, Size: 500, Time: base.Add(3 * time.Minute),
		},
	}
	actions := parseActionSet("alert")

	eval := buildSignalResults(trades, nil, actions, evalOptions{
		StakeUSD:    10,
		MinNotional: 500,
		ListMinNotional: map[string]float64{
			"sports": 1500,
			"watch":  750,
		},
	})
	if eval.RawBuys != 3 {
		t.Fatalf("RawBuys=%d, want 3", eval.RawBuys)
	}
	if len(eval.Results) != 3 {
		t.Fatalf("Results len=%d, want 3", len(eval.Results))
	}
	got := map[string]bool{}
	for _, res := range eval.Results {
		got[res.Buy.TradeID] = true
	}
	if got["0xsportssmall"] {
		t.Fatalf("sports $1000 buy should be below list min: %#v", got)
	}
	for _, want := range []string{"0xsportslarge", "0xwatch", "0xunknown"} {
		if !got[want] {
			t.Fatalf("missing %s in evaluated results: %#v", want, got)
		}
	}
}

func TestBuildSignalResults_FlagsPolicyViolations(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet: "0xabc", List: "leaderboard_watch", Side: "BUY", Action: "alert", TradeID: "0xmlb",
			AssetID: "asset-mlb", Market: "Boston Red Sox vs. Chicago White Sox",
			Outcome: "Boston Red Sox", Price: 0.51, Size: 6000, Time: base,
		},
		{
			Wallet: "0xabc", List: "watch", Side: "BUY", Action: "alert", TradeID: "0xspread",
			AssetID: "asset-spread", Market: "Spread: Milwaukee Brewers (-1.5)",
			Outcome: "Milwaukee Brewers", Price: 0.44, Size: 1000, Time: base.Add(time.Minute),
		},
		{
			Wallet: "0xabc", List: "watch", Side: "BUY", Action: "alert", TradeID: "0xprice",
			AssetID: "asset-price", Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
			Outcome: "MOUZ", Price: 0.99, Size: 1000, Time: base.Add(2 * time.Minute),
		},
		{
			Wallet: "0xabc", List: "watch", Side: "BUY", Action: "alert", TradeID: "0xgood",
			AssetID: "asset-good", Market: "LoL: T1 vs Gen.G (BO3) - LCK Summer",
			Outcome: "T1", Price: 0.61, Size: 1000, Time: base.Add(3 * time.Minute),
		},
		{
			Wallet: "0xabc", List: "watch", Side: "BUY", Action: "alert", TradeID: "0xwnba",
			AssetID: "asset-wnba", Market: "Golden State Valkyries vs. Toronto Tempo",
			Outcome: "Golden State Valkyries", Price: 0.52, Size: 1000, Time: base.Add(4 * time.Minute),
		},
		{
			Wallet: "0xabc", List: "watch", Side: "SELL", Action: "alert", TradeID: "0xsell",
			AssetID: "asset-sell", Market: "Boston Red Sox vs. Chicago White Sox",
			Outcome: "Boston Red Sox", Price: 0.99, Size: 1000, Time: base.Add(5 * time.Minute),
		},
	}

	eval := buildSignalResults(trades, nil, parseActionSet("alert"), evalOptions{
		StakeUSD:    10,
		MinNotional: 750,
		ListMinNotional: map[string]float64{
			"leaderboard_watch": 5000,
			"watch":             750,
		},
	})
	if len(eval.PolicyViolations) != 3 {
		t.Fatalf("PolicyViolations len=%d, want 3: %#v", len(eval.PolicyViolations), eval.PolicyViolations)
	}
	reasons := map[string]int{}
	for _, v := range eval.PolicyViolations {
		reasons[v.Reason]++
	}
	if reasons["category_filtered"] != 1 || reasons["derivative_filtered"] != 1 || reasons["price_filtered"] != 1 {
		t.Fatalf("violation reasons=%#v, want one category/derivative/price", reasons)
	}
}

func TestBuildSignalResults_FiltersBuysBeforeSince(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	trades := []whaleTrade{
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xbefore",
			AssetID: "asset-before", Market: "Will England win on 2026-07-11?",
			Price: 0.5, Size: 1000, Time: base.Add(-time.Minute),
		},
		{
			Wallet: "0xabc", Side: "BUY", Action: "alert", TradeID: "0xafter",
			AssetID: "asset-after", Market: "Will Spain win on 2026-07-10?",
			Price: 0.5, Size: 1000, Time: base.Add(time.Minute),
		},
		{
			Wallet: "0xabc", Side: "SELL", Action: "alert", TradeID: "0xsell",
			AssetID: "asset-after", Market: "Will Spain win on 2026-07-10?",
			Price: 0.75, Size: 1000, Time: base.Add(2 * time.Minute),
		},
	}

	eval := buildSignalResults(trades, nil, parseActionSet("alert"), evalOptions{
		StakeUSD:    10,
		MinNotional: 500,
		Since:       base,
	})
	if eval.RawBuys != 1 {
		t.Fatalf("RawBuys=%d, want 1", eval.RawBuys)
	}
	if len(eval.Results) != 1 {
		t.Fatalf("Results len=%d, want 1", len(eval.Results))
	}
	res := eval.Results[0]
	if res.Buy.TradeID != "0xafter" {
		t.Fatalf("kept trade=%s, want 0xafter", res.Buy.TradeID)
	}
	if res.ExitSource != "sell" || res.ReturnPct != 50 {
		t.Fatalf("exit=%s return=%.1f, want sell/50", res.ExitSource, res.ReturnPct)
	}
}

func TestParseListMinNotional(t *testing.T) {
	got := parseListMinNotional("core=1000,sports=1500, watch = 750, bad, scout=-1")
	if got["core"] != 1000 || got["sports"] != 1500 || got["watch"] != 750 {
		t.Fatalf("unexpected list mins: %#v", got)
	}
	if _, ok := got["scout"]; ok {
		t.Fatalf("negative list min should be ignored: %#v", got)
	}
}

func TestApplyTradeMetas_UsesCurrentWalletList(t *testing.T) {
	trades := []whaleTrade{{
		Wallet: "0x1111111111111111111111111111111111111111",
		List:   "core",
		Tier:   "A",
	}}
	metas := map[string]walletMeta{
		"0x1111111111111111111111111111111111111111": {List: "sports", Tier: "B"},
	}

	applyTradeMetas(trades, metas)
	if trades[0].List != "sports" || trades[0].Tier != "B" {
		t.Fatalf("trade meta=%s/%s, want sports/B", trades[0].List, trades[0].Tier)
	}
}

func TestLoadWalletMetas_MergesCommaSeparatedFiles(t *testing.T) {
	dir := t.TempDir()
	push := filepath.Join(dir, "push.txt")
	leaderboard := filepath.Join(dir, "leaderboard.txt")
	if err := os.WriteFile(push, []byte("0x1111111111111111111111111111111111111111 # list=watch tier=B\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(leaderboard, []byte(strings.Join([]string{
		"0x2222222222222222222222222222222222222222 # list=leaderboard_watch tier=C",
		"0x1111111111111111111111111111111111111111 # list=leaderboard_push tier=A",
		"",
	}, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}

	metas, err := loadWalletMetas(push + "," + leaderboard)
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 2 {
		t.Fatalf("metas len=%d, want 2: %#v", len(metas), metas)
	}
	if got := metas["0x2222222222222222222222222222222222222222"]; got.List != "leaderboard_watch" || got.Tier != "C" {
		t.Fatalf("leaderboard meta=%#v, want leaderboard_watch/C", got)
	}
	if got := metas["0x1111111111111111111111111111111111111111"]; got.List != "leaderboard_push" || got.Tier != "A" {
		t.Fatalf("later file should override duplicate meta, got %#v", got)
	}
}

func TestBuildSnapshot_IncludesEventCappedBreakdowns(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	results := []*signalResult{
		{
			Buy: whaleTrade{
				Wallet: "0xabc", Label: "alpha", List: "core", Tier: "A",
				Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
				Time:   base,
			},
			ExitSource: "mark",
			StakeUSD:   10,
			PnLUSD:     2,
			ReturnPct:  20,
		},
		{
			Buy: whaleTrade{
				Wallet: "0xabc", Label: "alpha", List: "core", Tier: "A",
				Market: "Dota 2: REKONIX vs MOUZ - Game 2 Winner",
				Time:   base.Add(time.Minute),
			},
			ExitSource: "mark",
			StakeUSD:   10,
			PnLUSD:     -1,
			ReturnPct:  -10,
		},
		{
			Buy: whaleTrade{
				Wallet: "0xdef", Label: "beta", List: "watch", Tier: "B",
				Market: "Dota 2: L1ga Team vs PlayTime - Game 1 Winner",
				Time:   base.Add(2 * time.Minute),
			},
			ExitSource: "mark",
			StakeUSD:   10,
			PnLUSD:     3,
			ReturnPct:  30,
		},
	}

	snap := buildSnapshot("log.jsonl", "wallets.txt", "report.md", evaluation{Results: results}, evalOptions{StakeUSD: 10})
	if snap.EventCappedSignals != 2 {
		t.Fatalf("EventCappedSignals=%d, want 2", snap.EventCappedSignals)
	}
	if len(snap.EventCappedByList) != 2 {
		t.Fatalf("EventCappedByList len=%d, want 2", len(snap.EventCappedByList))
	}
	if len(snap.EventCappedByWallet) != 2 {
		t.Fatalf("EventCappedByWallet len=%d, want 2", len(snap.EventCappedByWallet))
	}
	if snap.EventCappedByList[0].List != "core" || snap.EventCappedByList[0].Signals != 1 || snap.EventCappedByList[0].PnLUSD != 2 {
		t.Fatalf("core capped stats=%#v, want one +$2 entry", snap.EventCappedByList[0])
	}
	if snap.EventCappedByWallet[0].Wallet != "0xdef" || snap.EventCappedByWallet[0].ReturnPct != 30 {
		t.Fatalf("top capped wallet=%#v, want 0xdef with 30%% ROI", snap.EventCappedByWallet[0])
	}
}

func TestBuildSnapshot_IncludesProvenMetrics(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	results := []*signalResult{
		{
			Buy: whaleTrade{
				Wallet: "0xabc", List: "core",
				Market: "Dota 2: REKONIX vs MOUZ - Game 1 Winner",
				Time:   base,
			},
			ExitSource: "mark",
			StakeUSD:   10,
			PnLUSD:     4,
			ReturnPct:  40,
		},
		{
			Buy: whaleTrade{
				Wallet: "0xabc", List: "core",
				Market: "Dota 2: REKONIX vs MOUZ - Game 2 Winner",
				Time:   base.Add(time.Minute),
			},
			ExitSource: "settlement",
			StakeUSD:   10,
			PnLUSD:     9,
			ReturnPct:  90,
		},
		{
			Buy: whaleTrade{
				Wallet: "0xdef", List: "watch",
				Market: "Dota 2: L1ga Team vs PlayTime - Game 1 Winner",
				Time:   base.Add(2 * time.Minute),
			},
			ExitSource: "sell",
			StakeUSD:   10,
			PnLUSD:     -3,
			ReturnPct:  -30,
		},
	}

	snap := buildSnapshot("log.jsonl", "wallets.txt", "report.md", evaluation{Results: results}, evalOptions{StakeUSD: 10})
	if snap.EvaluatedSignals != 3 || snap.PnLUSD != 10 || math.Abs(snap.ReturnPct-100.0/3.0) > 0.000001 {
		t.Fatalf("floating summary=%#v, want all three evaluated including mark", snap)
	}
	if snap.ProvenSignals != 2 {
		t.Fatalf("ProvenSignals=%d, want 2", snap.ProvenSignals)
	}
	if snap.ProvenPnLUSD != 6 {
		t.Fatalf("ProvenPnLUSD=%.2f, want 6", snap.ProvenPnLUSD)
	}
	if snap.ProvenReturnPct != 30 {
		t.Fatalf("ProvenReturnPct=%.1f, want 30", snap.ProvenReturnPct)
	}
	if snap.EventCappedProvenSignals != 2 || snap.EventCappedProvenPnLUSD != 6 {
		t.Fatalf("event capped proven signals/pnl=%d/%.2f, want 2/6", snap.EventCappedProvenSignals, snap.EventCappedProvenPnLUSD)
	}
	if len(snap.EventCappedProvenByWallet) != 2 {
		t.Fatalf("EventCappedProvenByWallet len=%d, want 2", len(snap.EventCappedProvenByWallet))
	}
	if snap.EventCappedProvenByWallet[0].Wallet != "0xabc" || snap.EventCappedProvenByWallet[0].PnLUSD != 9 {
		t.Fatalf("top proven wallet=%#v, want 0xabc +$9", snap.EventCappedProvenByWallet[0])
	}
}

func TestSettlementPriceForAsset(t *testing.T) {
	m := feed.Market{
		Closed:           true,
		ClobTokenIDsRaw:  `["asset-win","asset-lose"]`,
		OutcomePricesRaw: `["1","0"]`,
	}

	px, ok := settlementPriceForAsset(m, "asset-win")
	if !ok || px != 1 {
		t.Fatalf("winning asset settlement = %.2f/%v, want 1/true", px, ok)
	}
	px, ok = settlementPriceForAsset(m, "asset-lose")
	if !ok || px != 0 {
		t.Fatalf("losing asset settlement = %.2f/%v, want 0/true", px, ok)
	}
	if _, ok := settlementPriceForAsset(m, "missing"); ok {
		t.Fatalf("missing asset should not settle")
	}
}
