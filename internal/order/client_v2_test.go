package order

import (
	"math/big"
	"testing"
)

func TestAmountsForIntentPreservesTickPrice(t *testing.T) {
	maker, taker, err := amountsForIntent(Intent{
		AssetID: "asset", Side: Buy, SizeUSD: 2, LimitPx: 0.936,
	})
	if err != nil {
		t.Fatal(err)
	}
	if maker.String() != "1999998" || taker.String() != "2136750" {
		t.Fatalf("buy amounts: maker=%s taker=%s", maker, taker)
	}
	price, _ := new(big.Rat).SetFrac(maker, taker).Float64()
	if diff := price - 0.936; diff > 1e-12 || diff < -1e-12 {
		t.Fatalf("implied buy price=%v", price)
	}

	maker, taker, err = amountsForIntent(Intent{
		AssetID: "asset", Side: Sell, SizeUSD: 9.36, SizeShares: 10, LimitPx: 0.936,
	})
	if err != nil {
		t.Fatal(err)
	}
	if maker.String() != "10000000" || taker.String() != "9360000" {
		t.Fatalf("sell amounts: maker=%s taker=%s", maker, taker)
	}
}

func TestFillFromAmountsUsesActualMatchedAmounts(t *testing.T) {
	fallbackMaker := big.NewInt(20_000_000)
	fallbackTaker := big.NewInt(40_000_000)

	units, price := fillFromAmounts(Buy, "19500000", "39000000", fallbackMaker, fallbackTaker)
	if units != 39 || price != 0.5 {
		t.Fatalf("buy fill units=%v price=%v", units, price)
	}

	units, price = fillFromAmounts(Sell, "10000000", "4900000", fallbackTaker, fallbackMaker)
	if units != 10 || price < 0.49-1e-12 || price > 0.49+1e-12 {
		t.Fatalf("sell fill units=%v price=%v", units, price)
	}
}

func TestParseFixedAmount(t *testing.T) {
	if got := parseFixedAmount("100000000"); got != 100 {
		t.Fatalf("fixed amount=%v", got)
	}
	if got := parseFixedAmount("12.5"); got != 12.5 {
		t.Fatalf("decimal amount=%v", got)
	}
}
