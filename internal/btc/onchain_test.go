package btc

import (
	"testing"
)

func TestOnChainModifier_Bullish(t *testing.T) {
	m := OnChainMetrics{OnChainSignal: "BULLISH", OnChainScore: 0.3}
	mod := m.OnChainModifier("BUY_YES", true) // reach BUY_YES + bullish
	if mod <= 1.0 {
		t.Errorf("bullish should amplify reach BUY_YES, got %.2f", mod)
	}
}

func TestOnChainModifier_Bearish(t *testing.T) {
	m := OnChainMetrics{OnChainSignal: "BEARISH", OnChainScore: -0.3}
	mod := m.OnChainModifier("BUY_NO", false) // dip BUY_NO + bearish
	if mod <= 1.0 {
		t.Errorf("bearish should amplify dip BUY_NO, got %.2f", mod)
	}
}

func TestOnChainModifier_Neutral(t *testing.T) {
	m := OnChainMetrics{OnChainSignal: "NEUTRAL", OnChainScore: 0.0}
	mod := m.OnChainModifier("BUY_YES", true)
	if mod != 1.0 {
		t.Errorf("neutral should return 1.0, got %.2f", mod)
	}
}

func TestClampFloat(t *testing.T) {
	if v := clampFloat(1.5, -1, 1); v != 1.0 {
		t.Errorf("expected 1.0, got %.2f", v)
	}
	if v := clampFloat(-2.0, -1, 1); v != -1.0 {
		t.Errorf("expected -1.0, got %.2f", v)
	}
	if v := clampFloat(0.5, -1, 1); v != 0.5 {
		t.Errorf("expected 0.5, got %.2f", v)
	}
}
