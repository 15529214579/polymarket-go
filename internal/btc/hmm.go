package btc

import (
	"fmt"
	"log/slog"
	"math"
)

// HMM implements a simple Hidden Markov Model with 3 hidden states (regimes)
// and NStates (15) observation symbols derived from (return_bucket, volume_regime).
//
// Regimes:
//   0 = Trending   (directional moves persist)
//   1 = MeanRevert (moves reverse quickly)
//   2 = Volatile   (big moves, no clear direction)

const (
	RegimeTrend  = 0
	RegimeMR     = 1
	RegimeVolat  = 2
	NRegimes     = 3
)

var regimeNames = [NRegimes]string{"TREND", "MEAN_REVERT", "VOLATILE"}

// RegimeName returns a human-readable label for a regime.
func RegimeName(r int) string {
	if r < 0 || r >= NRegimes {
		return fmt.Sprintf("UNKNOWN(%d)", r)
	}
	return regimeNames[r]
}

// HMMModel holds the trained parameters.
type HMMModel struct {
	Pi    [NRegimes]float64             // initial state distribution
	Trans [NRegimes][NRegimes]float64   // transition matrix
	Emit  [NRegimes][NStates]float64    // emission matrix
}

// DefaultHMMPriors returns hand-tuned priors based on BTC market structure.
// These are starting points that Baum-Welch will refine.
func DefaultHMMPriors() HMMModel {
	var m HMMModel

	m.Pi = [NRegimes]float64{0.4, 0.3, 0.3}

	m.Trans = [NRegimes][NRegimes]float64{
		{0.85, 0.10, 0.05}, // trend: sticky, sometimes switches to MR
		{0.10, 0.80, 0.10}, // mean-revert: sticky
		{0.10, 0.10, 0.80}, // volatile: sticky
	}

	// Emission priors: trending regime favors SURGE/UP + HIGH vol;
	// mean-revert favors FLAT + any vol; volatile favors extreme moves.
	for obs := 0; obs < NStates; obs++ {
		ret, vol := StateComponents(obs)
		switch {
		case (ret == RetSurge || ret == RetUp || ret == RetDown || ret == RetCrash) && vol == VolHigh:
			m.Emit[RegimeTrend][obs] = 0.12
			m.Emit[RegimeMR][obs] = 0.03
			m.Emit[RegimeVolat][obs] = 0.10
		case ret == RetFlat:
			m.Emit[RegimeTrend][obs] = 0.02
			m.Emit[RegimeMR][obs] = 0.15
			m.Emit[RegimeVolat][obs] = 0.03
		case (ret == RetSurge || ret == RetCrash):
			m.Emit[RegimeTrend][obs] = 0.08
			m.Emit[RegimeMR][obs] = 0.02
			m.Emit[RegimeVolat][obs] = 0.12
		default:
			m.Emit[RegimeTrend][obs] = 0.05
			m.Emit[RegimeMR][obs] = 0.06
			m.Emit[RegimeVolat][obs] = 0.05
		}
	}

	for r := 0; r < NRegimes; r++ {
		normalizeRow(m.Emit[r][:])
	}

	return m
}

func normalizeRow(row []float64) {
	var sum float64
	for _, v := range row {
		sum += v
	}
	if sum > 0 {
		for i := range row {
			row[i] /= sum
		}
	}
}

// TrainHMM runs Baum-Welch (EM) to fit the HMM to observed state sequences.
func TrainHMM(observations []int, maxIter int) HMMModel {
	m := DefaultHMMPriors()
	T := len(observations)
	if T < 10 || maxIter <= 0 {
		return m
	}

	for iter := 0; iter < maxIter; iter++ {
		alpha := forward(m, observations)
		beta := backward(m, observations)

		ll := logLikelihood(alpha, T)

		gamma := make([][NRegimes]float64, T)
		xi := make([][NRegimes][NRegimes]float64, T-1)

		for t := 0; t < T; t++ {
			var denom float64
			for i := 0; i < NRegimes; i++ {
				denom += alpha[t][i] * beta[t][i]
			}
			if denom == 0 {
				denom = 1e-300
			}
			for i := 0; i < NRegimes; i++ {
				gamma[t][i] = alpha[t][i] * beta[t][i] / denom
			}
		}

		for t := 0; t < T-1; t++ {
			var denom float64
			for i := 0; i < NRegimes; i++ {
				for j := 0; j < NRegimes; j++ {
					denom += alpha[t][i] * m.Trans[i][j] * m.Emit[j][observations[t+1]] * beta[t+1][j]
				}
			}
			if denom == 0 {
				denom = 1e-300
			}
			for i := 0; i < NRegimes; i++ {
				for j := 0; j < NRegimes; j++ {
					xi[t][i][j] = alpha[t][i] * m.Trans[i][j] * m.Emit[j][observations[t+1]] * beta[t+1][j] / denom
				}
			}
		}

		// M-step: update parameters
		for i := 0; i < NRegimes; i++ {
			m.Pi[i] = gamma[0][i]
		}
		normalizeRow(m.Pi[:])

		for i := 0; i < NRegimes; i++ {
			var gammaSum float64
			for t := 0; t < T-1; t++ {
				gammaSum += gamma[t][i]
			}
			if gammaSum == 0 {
				gammaSum = 1e-300
			}
			for j := 0; j < NRegimes; j++ {
				var xiSum float64
				for t := 0; t < T-1; t++ {
					xiSum += xi[t][i][j]
				}
				m.Trans[i][j] = xiSum / gammaSum
			}
			normalizeRow(m.Trans[i][:])
		}

		for i := 0; i < NRegimes; i++ {
			var gammaSum float64
			for t := 0; t < T; t++ {
				gammaSum += gamma[t][i]
			}
			if gammaSum == 0 {
				gammaSum = 1e-300
			}
			for obs := 0; obs < NStates; obs++ {
				var sum float64
				for t := 0; t < T; t++ {
					if observations[t] == obs {
						sum += gamma[t][i]
					}
				}
				m.Emit[i][obs] = sum / gammaSum
			}
			normalizeRow(m.Emit[i][:])
		}

		if iter > 0 {
			_ = ll // convergence check would go here
		}
	}

	return m
}

func forward(m HMMModel, obs []int) [][NRegimes]float64 {
	T := len(obs)
	alpha := make([][NRegimes]float64, T)

	for i := 0; i < NRegimes; i++ {
		alpha[0][i] = m.Pi[i] * m.Emit[i][obs[0]]
	}
	scaleRow(alpha[0][:])

	for t := 1; t < T; t++ {
		for j := 0; j < NRegimes; j++ {
			var sum float64
			for i := 0; i < NRegimes; i++ {
				sum += alpha[t-1][i] * m.Trans[i][j]
			}
			alpha[t][j] = sum * m.Emit[j][obs[t]]
		}
		scaleRow(alpha[t][:])
	}
	return alpha
}

func backward(m HMMModel, obs []int) [][NRegimes]float64 {
	T := len(obs)
	beta := make([][NRegimes]float64, T)

	for i := 0; i < NRegimes; i++ {
		beta[T-1][i] = 1.0
	}
	scaleRow(beta[T-1][:])

	for t := T - 2; t >= 0; t-- {
		for i := 0; i < NRegimes; i++ {
			var sum float64
			for j := 0; j < NRegimes; j++ {
				sum += m.Trans[i][j] * m.Emit[j][obs[t+1]] * beta[t+1][j]
			}
			beta[t][i] = sum
		}
		scaleRow(beta[t][:])
	}
	return beta
}

func scaleRow(row []float64) {
	var sum float64
	for _, v := range row {
		sum += v
	}
	if sum > 0 {
		for i := range row {
			row[i] /= sum
		}
	}
}

func logLikelihood(alpha [][NRegimes]float64, T int) float64 {
	var sum float64
	for i := 0; i < NRegimes; i++ {
		sum += alpha[T-1][i]
	}
	if sum <= 0 {
		return math.Inf(-1)
	}
	return math.Log(sum)
}

// Viterbi returns the most likely regime sequence for the observations.
func Viterbi(m HMMModel, obs []int) []int {
	T := len(obs)
	if T == 0 {
		return nil
	}

	delta := make([][NRegimes]float64, T)
	psi := make([][NRegimes]int, T)

	for i := 0; i < NRegimes; i++ {
		delta[0][i] = math.Log(m.Pi[i]+1e-300) + math.Log(m.Emit[i][obs[0]]+1e-300)
	}

	for t := 1; t < T; t++ {
		for j := 0; j < NRegimes; j++ {
			best := math.Inf(-1)
			bestI := 0
			for i := 0; i < NRegimes; i++ {
				v := delta[t-1][i] + math.Log(m.Trans[i][j]+1e-300)
				if v > best {
					best = v
					bestI = i
				}
			}
			delta[t][j] = best + math.Log(m.Emit[j][obs[t]]+1e-300)
			psi[t][j] = bestI
		}
	}

	path := make([]int, T)
	best := math.Inf(-1)
	for i := 0; i < NRegimes; i++ {
		if delta[T-1][i] > best {
			best = delta[T-1][i]
			path[T-1] = i
		}
	}
	for t := T - 2; t >= 0; t-- {
		path[t] = psi[t+1][path[t+1]]
	}
	return path
}

// CandlesToObservations converts candles to Markov state indices for HMM.
func CandlesToObservations(candles []Candle) []int {
	if len(candles) < 25 {
		return nil
	}

	var avg24hVol float64
	for i := 0; i < 24 && i < len(candles); i++ {
		avg24hVol += candles[i].Volume
	}
	avg24hVol /= float64(min24(24, len(candles)))

	obs := make([]int, 0, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		if candles[i-1].Close <= 0 {
			continue
		}
		ret := (candles[i].Close/candles[i-1].Close - 1) * 100
		retBucket := ClassifyReturn(ret)
		volBucket := ClassifyVolume(candles[i].Volume, avg24hVol)

		// rolling avg volume update
		if i >= 24 {
			avg24hVol = 0
			for j := i - 23; j <= i; j++ {
				avg24hVol += candles[j].Volume
			}
			avg24hVol /= 24
		}

		obs = append(obs, StateIdx(retBucket, volBucket))
	}
	return obs
}

func min24(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DetectCurrentRegime trains an HMM on the given candles and returns the
// most likely current regime plus the full regime sequence.
func DetectCurrentRegime(candles []Candle) (currentRegime int, regimeName_ string, confidence float64) {
	obs := CandlesToObservations(candles)
	if len(obs) < 30 {
		return RegimeVolat, "VOLATILE", 0.0
	}

	model := TrainHMM(obs, 20)
	path := Viterbi(model, obs)
	if len(path) == 0 {
		return RegimeVolat, "VOLATILE", 0.0
	}

	current := path[len(path)-1]

	// Confidence: fraction of last 10 ticks in the same regime
	lookback := 10
	if lookback > len(path) {
		lookback = len(path)
	}
	same := 0
	for i := len(path) - lookback; i < len(path); i++ {
		if path[i] == current {
			same++
		}
	}
	conf := float64(same) / float64(lookback)

	slog.Info("hmm.regime",
		"current", RegimeName(current),
		"confidence", fmt.Sprintf("%.2f", conf),
		"obs_count", len(obs),
		"last_10_regime_pct", fmt.Sprintf("%.0f%%", conf*100),
	)

	return current, RegimeName(current), conf
}

// HMMRegimeFilter returns true if the current regime supports trading in
// the Up/Down direction. Only trade directionally in TREND regime.
// MEAN_REVERT can be used for contrarian bets. VOLATILE → skip.
func HMMRegimeFilter(regime int, confidence float64) bool {
	switch regime {
	case RegimeTrend:
		return confidence >= 0.5
	case RegimeMR:
		return confidence >= 0.7
	case RegimeVolat:
		return false
	default:
		return false
	}
}
