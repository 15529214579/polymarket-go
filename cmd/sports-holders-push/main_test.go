package main

import (
	"testing"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func TestPositiveOutcomeTokens(t *testing.T) {
	m := walletdiscover.Market{
		ClobTokenIDsRaw: `["yes-token","no-token"]`,
		OutcomesRaw:     `["Yes","No"]`,
	}
	got := positiveOutcomeTokens(m)
	if !got["yes-token"] || got["no-token"] || len(got) != 1 {
		t.Fatalf("positiveOutcomeTokens() = %#v, want Yes token only", got)
	}
}

func TestPositiveOutcomeTokensKeepsConcreteScores(t *testing.T) {
	m := walletdiscover.Market{
		ClobTokenIDsRaw: `["score-10","score-11","other"]`,
		OutcomesRaw:     `["1-0","1-1","Other"]`,
	}
	got := positiveOutcomeTokens(m)
	if !got["score-10"] || !got["score-11"] || got["other"] || len(got) != 2 {
		t.Fatalf("positiveOutcomeTokens() = %#v, want concrete scores only", got)
	}
}
