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
