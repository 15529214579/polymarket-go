package notify

import (
	"strings"
	"testing"
	"time"
)

func TestParseMarketTitle(t *testing.T) {
	cases := []struct {
		in        string
		wantMatch string
		wantCtx   string
	}{
		{
			"LoL: Shifters vs G2 Esports - Game 1 Winner",
			"LoL: Shifters vs G2 Esports",
			"Game 1 Winner",
		},
		{
			"LoL: Weibo Gaming vs Oh My God (BO3) - LPL Group Ascend",
			"LoL: Weibo Gaming vs Oh My God",
			"BO3 · LPL Group Ascend",
		},
		{
			"LoL: Gen.G Global Academy vs Nongshim Esports Academy (BO3) - LCK Challengers League Rounds 1-2",
			"LoL: Gen.G Global Academy vs Nongshim Esports Academy",
			"BO3 · LCK Challengers League Rounds 1-2",
		},
		{
			"Games Total: O/U 2.5",
			"Games Total: O/U 2.5",
			"",
		},
		{
			"Game Handicap: BLG (-1.5) vs Invictus Gaming (+1.5)",
			// No " - " separator at top level so returns as-is.
			"Game Handicap: BLG (-1.5) vs Invictus Gaming (+1.5)",
			"",
		},
	}
	for _, c := range cases {
		m, ctx := ParseMarketTitle(c.in)
		if m != c.wantMatch || ctx != c.wantCtx {
			t.Errorf("ParseMarketTitle(%q) = (%q, %q), want (%q, %q)", c.in, m, ctx, c.wantMatch, c.wantCtx)
		}
	}
}

func TestHumanizeEndIn(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		end  time.Time
		want string
	}{
		{now.Add(2*time.Hour + 5*time.Minute), "2h 05m"},
		{now.Add(45 * time.Minute), "45m"},
		{now.Add(30 * time.Second), "<1m"},
		{now.Add(-1 * time.Minute), ""},
		{time.Time{}, ""},
	}
	for _, c := range cases {
		got := HumanizeEndIn(now, c.end)
		if got != c.want {
			t.Errorf("HumanizeEndIn(end=%v) = %q, want %q", c.end, got, c.want)
		}
	}
}

func TestFormatSignalPrompt_ShowsSignalOnly(t *testing.T) {
	s := FormatSignalPrompt(SignalPromptEvent{
		Nonce:   "n",
		Match:   "LoL: Shifters vs G2 Esports",
		Context: "Game 1 Winner",
		EndIn:   "2h 05m",
		Choices: []SignalChoice{
			{Slot: 0, Outcome: "Shifters", Mid: 0.235, IsSignal: true},
			{Slot: 1, Outcome: "G2 Esports", Mid: 0.765, IsSignal: false},
		},
		DeltaPP: 4.2, TailUps: 4, TailLen: 5, BuyRatio: 0.78,
		ExpiresIn: 10 * time.Minute,
	})
	for _, want := range []string{
		"Shifters ↑ @ 0.2350",
		"LoL: Shifters vs G2 Esports",
		"Game 1 Winner · 2h 05m",
		"Δ+4.20pp",
		"buy 78%",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("FormatSignalPrompt missing %q; got:\n%s", want, s)
		}
	}
	// Compact format drops the verbose "当前" / "tail x/y" / TTL line.
	for _, absent := range []string{"0.7650", "选 ", "← 信号", "当前 ", "tail ", "按钮 "} {
		if strings.Contains(s, absent) {
			t.Errorf("FormatSignalPrompt leaked %q; got:\n%s", absent, s)
		}
	}
}
