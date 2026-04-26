package btc

import "fmt"

// ---------------------------------------------------------------------------
// State discretisation
// ---------------------------------------------------------------------------
//
// BTC hourly return bucket (5 states):
//   SURGE  >+2%
//   UP     +0.5% to +2%
//   FLAT   -0.5% to +0.5%
//   DOWN   -2% to -0.5%
//   CRASH  < -2%
//
// Volume regime relative to 24-hour rolling average (3 states):
//   HIGH   > 1.5× avg
//   MED    0.75× to 1.5× avg
//   LOW    < 0.75× avg
//
// Combined state = returnBucket*3 + volumeRegime  → 15 states total

const (
	RetSurge = iota // >+2%
	RetUp           // +0.5% to +2%
	RetFlat         // -0.5% to +0.5%
	RetDown         // -2% to -0.5%
	RetCrash        // <-2%
	nRetBuckets = 5
)

const (
	VolHigh = iota // > 1.5× 24h avg
	VolMed         // 0.75× to 1.5×
	VolLow         // < 0.75×
	nVolRegimes = 3
)

// NStates is the total number of combined (return × volume) Markov states.
const NStates = nRetBuckets * nVolRegimes // 15

// StateIdx encodes a (returnBucket, volumeRegime) pair into a flat index.
func StateIdx(ret, vol int) int { return ret*nVolRegimes + vol }

// StateComponents decodes a flat state index.
func StateComponents(s int) (ret, vol int) {
	return s / nVolRegimes, s % nVolRegimes
}

var retNames = [nRetBuckets]string{"SURGE", "UP", "FLAT", "DOWN", "CRASH"}
var volNames = [nVolRegimes]string{"HIGH", "MED", "LOW"}

// StateName returns a human-readable label for state s.
func StateName(s int) string {
	ret, vol := StateComponents(s)
	if ret < 0 || ret >= nRetBuckets || vol < 0 || vol >= nVolRegimes {
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
	return fmt.Sprintf("%s/%s", retNames[ret], volNames[vol])
}

// ClassifyReturn maps a percentage return into a return bucket.
func ClassifyReturn(pct float64) int {
	switch {
	case pct > 2.0:
		return RetSurge
	case pct > 0.5:
		return RetUp
	case pct >= -0.5:
		return RetFlat
	case pct >= -2.0:
		return RetDown
	default:
		return RetCrash
	}
}

// ClassifyVolume maps a candle's volume relative to the rolling 24h average
// into a volume regime.
func ClassifyVolume(vol, avg24h float64) int {
	if avg24h <= 0 {
		return VolMed
	}
	ratio := vol / avg24h
	switch {
	case ratio > 1.5:
		return VolHigh
	case ratio >= 0.75:
		return VolMed
	default:
		return VolLow
	}
}

// CandleState derives the Markov state for candle i given the preceding candles.
// It requires at least 2 candles (i >= 1) and uses a 24-bar rolling window for
// volume normalisation.
func CandleState(candles []Candle, i int) (int, bool) {
	if i < 1 || i >= len(candles) {
		return 0, false
	}
	prev := candles[i-1]
	curr := candles[i]
	if prev.Close <= 0 || curr.Close <= 0 {
		return 0, false
	}

	pctReturn := (curr.Close - prev.Close) / prev.Close * 100

	// rolling 24-bar volume average ending at i-1
	start := i - 24
	if start < 0 {
		start = 0
	}
	var sumVol float64
	n := 0
	for j := start; j < i; j++ {
		sumVol += candles[j].Volume
		n++
	}
	avg24h := 0.0
	if n > 0 {
		avg24h = sumVol / float64(n)
	}

	ret := ClassifyReturn(pctReturn)
	vol := ClassifyVolume(curr.Volume, avg24h)
	return StateIdx(ret, vol), true
}

// ---------------------------------------------------------------------------
// Transition matrix
// ---------------------------------------------------------------------------

// TransitionMatrix is an NStates×NStates count matrix.
// transMatrix[from][to] = number of observed (from→to) transitions.
type TransitionMatrix [NStates][NStates]int64

// Train builds a transition count matrix from a slice of candles.
// Returns the number of transitions recorded.
func Train(candles []Candle) (TransitionMatrix, int) {
	var tm TransitionMatrix
	count := 0
	for i := 2; i < len(candles); i++ {
		from, ok1 := CandleState(candles, i-1)
		to, ok2 := CandleState(candles, i)
		if !ok1 || !ok2 {
			continue
		}
		tm[from][to]++
		count++
	}
	return tm, count
}

// RowProbs converts the raw counts for state s into a probability distribution.
// Returns nil if the row has no observations.
func (tm *TransitionMatrix) RowProbs(s int) []float64 {
	var total int64
	for _, c := range tm[s] {
		total += c
	}
	if total == 0 {
		return nil
	}
	probs := make([]float64, NStates)
	for to, c := range tm[s] {
		probs[to] = float64(c) / float64(total)
	}
	return probs
}

// ---------------------------------------------------------------------------
// Expected-return table
// ---------------------------------------------------------------------------

// ReturnStats accumulates actual returns per state for evaluating the model.
type ReturnStats struct {
	N        int
	SumRet   float64 // sum of 1-period % returns following this state
	NPos     int     // count where next-period return > 0
}

func (r *ReturnStats) Add(ret float64) {
	r.N++
	r.SumRet += ret
	if ret > 0 {
		r.NPos++
	}
}

// AvgRet returns the mean return or 0 if no observations.
func (r ReturnStats) AvgRet() float64 {
	if r.N == 0 {
		return 0
	}
	return r.SumRet / float64(r.N)
}

// PosRate returns the fraction of observations with positive next-period return.
func (r ReturnStats) PosRate() float64 {
	if r.N == 0 {
		return 0
	}
	return float64(r.NPos) / float64(r.N)
}

// BuildReturnStats computes per-state forward-return statistics from candles.
func BuildReturnStats(candles []Candle) [NStates]ReturnStats {
	var stats [NStates]ReturnStats
	for i := 1; i+1 < len(candles); i++ {
		state, ok := CandleState(candles, i)
		if !ok {
			continue
		}
		curr := candles[i]
		next := candles[i+1]
		if curr.Close <= 0 || next.Close <= 0 {
			continue
		}
		fwdRet := (next.Close - curr.Close) / curr.Close * 100
		stats[state].Add(fwdRet)
	}
	return stats
}

// ---------------------------------------------------------------------------
// Prediction
// ---------------------------------------------------------------------------

// Prediction holds the model's output for a given current state.
type Prediction struct {
	CurrentState      int
	CurrentStateName  string
	NextProbs         []float64 // probability of transitioning to each next state
	ExpectedReturn    float64   // E[next-period return] = Σ p(s') × avgRet(s')
	BullProb          float64   // P(next state is SURGE or UP)
	BearProb          float64   // P(next state is CRASH or DOWN)
}

// Predict returns a Prediction from the trained model for the given current state.
func Predict(currentState int, tm *TransitionMatrix, retStats [NStates]ReturnStats) Prediction {
	probs := tm.RowProbs(currentState)

	pred := Prediction{
		CurrentState:     currentState,
		CurrentStateName: StateName(currentState),
	}

	if probs == nil {
		pred.NextProbs = make([]float64, NStates)
		return pred
	}
	pred.NextProbs = probs

	for s, p := range probs {
		pred.ExpectedReturn += p * retStats[s].AvgRet()
		ret, _ := StateComponents(s)
		if ret == RetSurge || ret == RetUp {
			pred.BullProb += p
		}
		if ret == RetCrash || ret == RetDown {
			pred.BearProb += p
		}
	}
	return pred
}

// PredictFromCandles is a convenience wrapper: it derives the current state
// from the last candle in the slice and returns a Prediction.
func PredictFromCandles(candles []Candle, tm *TransitionMatrix, retStats [NStates]ReturnStats) (Prediction, bool) {
	last := len(candles) - 1
	state, ok := CandleState(candles, last)
	if !ok {
		return Prediction{}, false
	}
	return Predict(state, tm, retStats), true
}
