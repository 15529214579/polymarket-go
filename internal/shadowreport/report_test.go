package shadowreport

import (
	"testing"
	"time"
)

func TestAnalyzeDedupesPositionPolicyAndTracksPostActual(t *testing.T) {
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	actual := entry.Add(12 * time.Minute)
	report := Analyze([]Observation{
		{PosID: "p1", Policy: "timeout_10m", EntryTime: entry, ObservedAt: entry.Add(10 * time.Minute), Question: "Dota 2: A vs B", NetPnLUSD: -1},
		{PosID: "p1", Policy: "timeout_10m", EntryTime: entry, ObservedAt: entry.Add(11 * time.Minute), Question: "Dota 2: A vs B", NetPnLUSD: -2},
		{PosID: "p1", Policy: "timeout_45m", EntryTime: entry, ObservedAt: entry.Add(45 * time.Minute), ActualCloseAt: actual, Question: "Dota 2: A vs B", NetPnLUSD: 3},
	})
	if report.Samples != 2 {
		t.Fatalf("samples=%d", report.Samples)
	}
	if report.ByPolicy[1].Name != "timeout_45m" || report.ByPolicy[1].PostActual != 1 || report.ByPolicy[1].NetPnLUSD != 3 {
		t.Fatalf("groups=%+v", report.ByPolicy)
	}
}

func TestAnalyzeBuildsMatchedTimeoutComparisons(t *testing.T) {
	entry := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	report := Analyze([]Observation{
		{PosID: "p1", Policy: "timeout_30m", EntryTime: entry, ObservedAt: entry.Add(30 * time.Minute), Question: "Dota 2: A vs B", NetPnLUSD: 1},
		{PosID: "p1", Policy: "timeout_45m", EntryTime: entry, ObservedAt: entry.Add(45 * time.Minute), Question: "Dota 2: A vs B", NetPnLUSD: 3},
		{PosID: "p2", Policy: "timeout_30m", EntryTime: entry.Add(time.Hour), ObservedAt: entry.Add(90 * time.Minute), Question: "CS2: C vs D", NetPnLUSD: -1},
		{PosID: "p2", Policy: "timeout_45m", EntryTime: entry.Add(time.Hour), ObservedAt: entry.Add(105 * time.Minute), ActualCloseAt: entry.Add(100 * time.Minute), Question: "CS2: C vs D", NetPnLUSD: 1},
		{PosID: "unmatched", Policy: "timeout_45m", EntryTime: entry, ObservedAt: entry.Add(45 * time.Minute), Question: "Dota 2: E vs F", NetPnLUSD: 100},
	})

	var got PairedComparison
	for _, row := range report.PairedTimeouts {
		if row.Category == "esports" && row.EarlierPolicy == "timeout_30m" && row.LaterPolicy == "timeout_45m" {
			got = row
			break
		}
	}
	if got.Samples != 2 || got.LaterBetter != 2 || got.PostActual != 1 {
		t.Fatalf("pair=%+v", got)
	}
	if got.EarlierNetPnLUSD != 0 || got.LaterNetPnLUSD != 4 || got.NetUpliftUSD != 4 || got.AverageUpliftUSD != 2 || got.LaterBetterPct != 100 {
		t.Fatalf("pair economics=%+v", got)
	}
}
