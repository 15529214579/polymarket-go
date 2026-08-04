package main

import (
	"testing"

	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func TestSettlementTokenFallbackMapsRecoveredMarketToPersistedCondition(t *testing.T) {
	open := []strategy.Position{
		{ID: "p1", Market: "old-condition-a", AssetID: "token-a"},
		{ID: "p2", Market: "old-condition-b", AssetID: "token-b"},
	}
	missing := []string{"old-condition-a", "old-condition-b"}
	tokens := settlementFallbackTokenIDs(open, missing)
	if len(tokens) != 2 || tokens[0] != "token-a" || tokens[1] != "token-b" {
		t.Fatalf("fallback tokens=%v", tokens)
	}

	byCondition := map[string]feed.Market{}
	recovered := indexSettlementMarketsByToken(open, missing, []feed.Market{{
		ConditionID:     "real-condition-a",
		ClobTokenIDsRaw: `["token-a"]`,
		Closed:          true,
	}}, byCondition)
	if recovered != 1 {
		t.Fatalf("recovered=%d", recovered)
	}
	market, ok := byCondition["old-condition-a"]
	if !ok || market.ConditionID != "real-condition-a" || !market.Closed {
		t.Fatalf("mapped market=%+v ok=%v", market, ok)
	}
	remaining := missingSettlementConditionIDs(missing, byCondition)
	if len(remaining) != 1 || remaining[0] != "old-condition-b" {
		t.Fatalf("remaining=%v", remaining)
	}
}
