package btc

import (
	"testing"
)

func TestComputeDepth_DeepBook(t *testing.T) {
	book := clobBook{
		Bids: []clobOrder{
			{Price: "0.50", Size: "200"},
			{Price: "0.49", Size: "300"},
			{Price: "0.48", Size: "400"},
		},
		Asks: []clobOrder{
			{Price: "0.52", Size: "200"},
			{Price: "0.53", Size: "300"},
		},
	}
	d := computeDepth("tok1", 55000, book)
	if d.DepthScore < 0.6 {
		t.Errorf("deep book should have high score, got %.2f", d.DepthScore)
	}
	if d.Spread < 0.01 {
		t.Errorf("spread should be ~0.02, got %.4f", d.Spread)
	}
}

func TestComputeDepth_EmptyBook(t *testing.T) {
	d := computeDepth("tok2", 50000, clobBook{})
	if d.DepthScore != 0.0 {
		t.Errorf("empty book should have score 0, got %.2f", d.DepthScore)
	}
}

func TestComputeDepth_ThinBook(t *testing.T) {
	book := clobBook{
		Bids: []clobOrder{{Price: "0.40", Size: "5"}},
		Asks: []clobOrder{{Price: "0.55", Size: "5"}},
	}
	d := computeDepth("tok3", 45000, book)
	if d.DepthScore >= 0.5 {
		t.Errorf("thin book with wide spread should score low, got %.2f", d.DepthScore)
	}
}

func TestDepthModifier_IlliquidMarket(t *testing.T) {
	md := MarketDepth{
		Depths: []OrderbookDepth{
			{Strike: 55000, DepthScore: 0.2},
			{Strike: 50000, DepthScore: 0.8},
		},
	}
	if mod := DepthModifier(md, 55000); mod > 0.6 {
		t.Errorf("illiquid market should get low modifier, got %.2f", mod)
	}
	if mod := DepthModifier(md, 50000); mod != 1.0 {
		t.Errorf("liquid market should get 1.0, got %.2f", mod)
	}
}

func TestDepthModifier_NoData(t *testing.T) {
	md := MarketDepth{}
	mod := DepthModifier(md, 60000)
	if mod != 0.8 {
		t.Errorf("no data should return 0.8, got %.2f", mod)
	}
}
