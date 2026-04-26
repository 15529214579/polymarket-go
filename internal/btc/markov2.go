package btc

import (
	"fmt"
	"math"
)

// Second-order Markov: state pair (s[t-1], s[t]) → s[t+1].
// 15 base states → 225 pair states → predict next among 15.

const NPairStates = NStates * NStates // 225

func PairIdx(prev, curr int) int { return prev*NStates + curr }

func PairComponents(p int) (prev, curr int) {
	return p / NStates, p % NStates
}

func PairName(p int) string {
	prev, curr := PairComponents(p)
	return fmt.Sprintf("%s→%s", StateName(prev), StateName(curr))
}

// Trans2 is a second-order transition count matrix: [pair][next]count.
type Trans2 [NPairStates][NStates]int64

func Train2(candles []Candle) (Trans2, int) {
	var tm Trans2
	count := 0
	for i := 3; i < len(candles); i++ {
		prev, ok1 := CandleState(candles, i-2)
		curr, ok2 := CandleState(candles, i-1)
		next, ok3 := CandleState(candles, i)
		if !ok1 || !ok2 || !ok3 {
			continue
		}
		tm[PairIdx(prev, curr)][next]++
		count++
	}
	return tm, count
}

func (tm *Trans2) RowProbs(pair int) []float64 {
	var total int64
	for _, c := range tm[pair] {
		total += c
	}
	if total == 0 {
		return nil
	}
	probs := make([]float64, NStates)
	for to, c := range tm[pair] {
		probs[to] = float64(c) / float64(total)
	}
	return probs
}

// Predict2 returns a Prediction using second-order context.
func Predict2(prev, curr int, tm2 *Trans2, retStats [NStates]ReturnStats) Prediction {
	pair := PairIdx(prev, curr)
	probs := tm2.RowProbs(pair)

	pred := Prediction{
		CurrentState:     curr,
		CurrentStateName: fmt.Sprintf("%s (2nd: %s)", StateName(curr), PairName(pair)),
	}

	if probs == nil {
		pred.NextProbs = make([]float64, NStates)
		pred.BullProb = 0.5
		pred.BearProb = 0.5
		return pred
	}
	pred.NextProbs = probs

	for s, p := range probs {
		pred.ExpectedReturn += p * retStats[s].AvgRet()
		posRate := 0.5
		if retStats[s].N > 0 {
			posRate = float64(retStats[s].NPos) / float64(retStats[s].N)
		}
		pred.BullProb += p * posRate
		pred.BearProb += p * (1.0 - posRate)
	}
	return pred
}

// BlendedPrediction merges first-order and second-order predictions.
// When second-order has enough data (minObs), it gets higher weight.
func BlendedPrediction(candles []Candle) (Prediction, bool) {
	n := len(candles)
	if n < 4 {
		return Prediction{}, false
	}

	tm1, _ := Train(candles)
	retStats := BuildReturnStats(candles)

	currState, ok := CandleState(candles, n-1)
	if !ok {
		return Prediction{}, false
	}
	pred1 := Predict(currState, &tm1, retStats)

	prevState, ok := CandleState(candles, n-2)
	if !ok {
		return pred1, true
	}

	tm2, _ := Train2(candles)
	pair := PairIdx(prevState, currState)
	probs2 := tm2.RowProbs(pair)

	if probs2 == nil {
		return pred1, true
	}

	pred2 := Predict2(prevState, currState, &tm2, retStats)

	var totalObs int64
	for _, c := range tm2[pair] {
		totalObs += c
	}

	const minObs = 10
	w2 := 0.0
	if totalObs >= minObs {
		w2 = math.Min(float64(totalObs)/50.0, 0.6)
	}
	w1 := 1.0 - w2

	blended := Prediction{
		CurrentState:     currState,
		CurrentStateName: fmt.Sprintf("%s (blended w2=%.0f%%)", StateName(currState), w2*100),
		NextProbs:        make([]float64, NStates),
		ExpectedReturn:   w1*pred1.ExpectedReturn + w2*pred2.ExpectedReturn,
		BullProb:         w1*pred1.BullProb + w2*pred2.BullProb,
		BearProb:         w1*pred1.BearProb + w2*pred2.BearProb,
	}
	for s := range blended.NextProbs {
		p1, p2 := 0.0, 0.0
		if pred1.NextProbs != nil && s < len(pred1.NextProbs) {
			p1 = pred1.NextProbs[s]
		}
		if pred2.NextProbs != nil && s < len(pred2.NextProbs) {
			p2 = pred2.NextProbs[s]
		}
		blended.NextProbs[s] = w1*p1 + w2*p2
	}

	return blended, true
}
