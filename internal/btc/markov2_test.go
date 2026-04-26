package btc

import (
	"math"
	"testing"
	"time"
)

func TestPairIdxRoundTrip(t *testing.T) {
	for p := 0; p < NStates; p++ {
		for c := 0; c < NStates; c++ {
			idx := PairIdx(p, c)
			gp, gc := PairComponents(idx)
			if gp != p || gc != c {
				t.Errorf("PairIdx(%d,%d)=%d → PairComponents=(%d,%d)", p, c, idx, gp, gc)
			}
		}
	}
}

func TestTrain2Empty(t *testing.T) {
	tm, n := Train2(nil)
	if n != 0 {
		t.Errorf("Train2(nil) = %d transitions, want 0", n)
	}
	_ = tm
}

func TestTrain2Flat(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, 60)
	for i := range candles {
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), 50000, 50010, 49990, 50000, 100)
	}
	tm, n := Train2(candles)
	if n == 0 {
		t.Fatal("expected transitions from flat candles")
	}

	for pair := 0; pair < NPairStates; pair++ {
		for next := 0; next < NStates; next++ {
			if tm[pair][next] > 0 {
				prev, curr := PairComponents(pair)
				prevRet, _ := StateComponents(prev)
				currRet, _ := StateComponents(curr)
				nextRet, _ := StateComponents(next)
				if prevRet != RetFlat || currRet != RetFlat || nextRet != RetFlat {
					t.Errorf("non-flat transition in flat candles: %s→%s (count=%d)",
						PairName(pair), StateName(next), tm[pair][next])
				}
			}
		}
	}
}

func TestTrans2RowProbsSumToOne(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, 200)
	price := 50000.0
	for i := range candles {
		if i%3 == 0 {
			price *= 1.01
		} else if i%5 == 0 {
			price *= 0.99
		}
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), price, price*1.005, price*0.995, price, 100+float64(i%20)*10)
	}
	tm, _ := Train2(candles)

	for pair := 0; pair < NPairStates; pair++ {
		probs := tm.RowProbs(pair)
		if probs == nil {
			continue
		}
		sum := 0.0
		for _, p := range probs {
			sum += p
		}
		if math.Abs(sum-1.0) > 1e-9 {
			t.Errorf("RowProbs(%s) sums to %f", PairName(pair), sum)
		}
	}
}

func TestPredict2NoData(t *testing.T) {
	var tm2 Trans2
	var retStats [NStates]ReturnStats

	pred := Predict2(0, 0, &tm2, retStats)
	if pred.ExpectedReturn != 0 {
		t.Errorf("expected return=0 with no data, got %f", pred.ExpectedReturn)
	}
	if pred.BullProb != 0.5 || pred.BearProb != 0.5 {
		t.Errorf("expected 50/50 with no data, got bull=%.2f bear=%.2f", pred.BullProb, pred.BearProb)
	}
}

func TestBlendedPrediction(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, 100)
	price := 50000.0
	for i := range candles {
		if i%2 == 0 {
			price *= 1.005
		} else {
			price *= 0.995
		}
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), price, price*1.001, price*0.999, price, 100)
	}

	pred, ok := BlendedPrediction(candles)
	if !ok {
		t.Fatal("BlendedPrediction returned ok=false")
	}
	if pred.BullProb+pred.BearProb > 1.0+1e-9 {
		t.Errorf("bull+bear = %f > 1", pred.BullProb+pred.BearProb)
	}

	_, ok = BlendedPrediction(candles[:3])
	if ok {
		t.Error("BlendedPrediction should fail with <4 candles")
	}
}
