package btc

import "testing"

func TestParseStrikeMultiCoin(t *testing.T) {
	tests := []struct {
		question string
		want     float64
	}{
		{"Will Bitcoin reach $120,000 before 2027?", 120000},
		{"Will the price of Bitcoin be above $72,000 on April 29?", 72000},
		{"Will the price of Ethereum be above $2,000 on May 1?", 2000},
		{"Will the price of Solana be above $80 on April 27?", 80},
		{"Will Ethereum reach $2,900 April 20-26?", 2900},
		{"Will the price of Solana be above $200 on May 3?", 200},
		{"No dollar sign here 2027", 0},
		{"Will Bitcoin hit $50K before June?", 50000},
	}
	for _, tt := range tests {
		got := parseStrikeFromQuestion(tt.question)
		if got != tt.want {
			t.Errorf("parseStrikeFromQuestion(%q) = %.0f, want %.0f", tt.question, got, tt.want)
		}
	}
}

func TestCoinConfigs(t *testing.T) {
	if CoinETH.BinancePair != "ETHUSDT" {
		t.Errorf("ETH pair = %s", CoinETH.BinancePair)
	}
	if CoinSOL.BinancePair != "SOLUSDT" {
		t.Errorf("SOL pair = %s", CoinSOL.BinancePair)
	}
	if CoinETH.GammaSlug == "" || CoinSOL.GammaSlug == "" {
		t.Error("missing gamma slug")
	}
}
