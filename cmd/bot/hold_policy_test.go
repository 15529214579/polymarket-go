package main

import (
	"testing"
	"time"
)

func TestSelectCopytradeHoldPolicy(t *testing.T) {
	defaultMax := 10 * time.Minute
	defaultEvent := 30 * time.Minute
	esports := 45 * time.Minute
	footballScore := 150 * time.Minute

	got := selectCopytradeHoldPolicy("Will AG.AL win the EWC League of Legends Tournament?", false, true, defaultMax, defaultEvent, esports, footballScore)
	if got.Name != "esports" || got.MaxHold != esports || got.EventHold != esports {
		t.Fatalf("paper esports policy=%+v", got)
	}
	got = selectCopytradeHoldPolicy("Will AG.AL win the EWC League of Legends Tournament?", false, false, defaultMax, defaultEvent, esports, footballScore)
	if got.Name != "default" || got.MaxHold != defaultMax || got.EventHold != defaultEvent {
		t.Fatalf("live policy changed=%+v", got)
	}
	got = selectCopytradeHoldPolicy("Exact Score: Spain 2-1 France", true, true, defaultMax, defaultEvent, esports, footballScore)
	if got.Name != "football_score" || got.MaxHold != footballScore || got.EventHold != footballScore {
		t.Fatalf("football score policy=%+v", got)
	}
}
