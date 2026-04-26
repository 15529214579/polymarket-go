package btc

import (
	"fmt"
	"log/slog"
	"math"
)

// KLDivergence computes KL(p || q) = Σ p(i) * log(p(i)/q(i)).
// Uses Laplace smoothing (add epsilon) to avoid log(0).
func KLDivergence(p, q []float64) float64 {
	if len(p) != len(q) || len(p) == 0 {
		return 0
	}
	const eps = 1e-10
	kl := 0.0
	for i := range p {
		pi := p[i] + eps
		qi := q[i] + eps
		kl += pi * math.Log(pi/qi)
	}
	return kl
}

// SymmetricKL returns (KL(p||q) + KL(q||p)) / 2 for a symmetric measure.
func SymmetricKL(p, q []float64) float64 {
	return (KLDivergence(p, q) + KLDivergence(q, p)) / 2.0
}

// DriftResult holds per-state and aggregate drift between two transition matrices.
type DriftResult struct {
	PerState   [NStates]float64 // symmetric KL divergence per state row
	MeanDrift  float64          // average across all states with data
	MaxDrift   float64          // worst-case single state drift
	MaxState   int              // which state has max drift
	StatesUsed int              // how many states had data in both matrices
	Drifted    bool             // true if MeanDrift > threshold
}

const DriftThreshold = 0.15 // symmetric KL > 0.15 = significant regime change

// MatrixDrift compares two first-order transition matrices.
func MatrixDrift(recent, full *TransitionMatrix) DriftResult {
	var dr DriftResult
	var sumKL float64

	for s := 0; s < NStates; s++ {
		pRecent := recent.RowProbs(s)
		pFull := full.RowProbs(s)
		if pRecent == nil || pFull == nil {
			continue
		}
		kl := SymmetricKL(pRecent, pFull)
		dr.PerState[s] = kl
		sumKL += kl
		dr.StatesUsed++
		if kl > dr.MaxDrift {
			dr.MaxDrift = kl
			dr.MaxState = s
		}
	}

	if dr.StatesUsed > 0 {
		dr.MeanDrift = sumKL / float64(dr.StatesUsed)
	}
	dr.Drifted = dr.MeanDrift > DriftThreshold
	return dr
}

// Matrix2Drift compares two second-order transition matrices.
func Matrix2Drift(recent, full *Trans2) float64 {
	var sumKL float64
	var n int
	for p := 0; p < NPairStates; p++ {
		pRecent := recent.RowProbs(p)
		pFull := full.RowProbs(p)
		if pRecent == nil || pFull == nil {
			continue
		}
		sumKL += SymmetricKL(pRecent, pFull)
		n++
	}
	if n == 0 {
		return 0
	}
	return sumKL / float64(n)
}

// CheckDrift trains on recent vs full candle history and reports drift.
// recentN is the number of most-recent candles for the "recent" window.
func CheckDrift(candles []Candle, recentN int) DriftResult {
	if len(candles) < recentN+30 {
		return DriftResult{}
	}

	tmFull, _ := Train(candles)
	recent := candles[len(candles)-recentN:]
	tmRecent, _ := Train(recent)

	dr := MatrixDrift(&tmRecent, &tmFull)

	// Also check second-order drift
	tm2Full, _ := Train2(candles)
	tm2Recent, _ := Train2(recent)
	drift2 := Matrix2Drift(&tm2Recent, &tm2Full)

	slog.Info("btc_strategy.drift_check",
		"candles_total", len(candles),
		"recent_window", recentN,
		"mean_drift_1st", fmt.Sprintf("%.4f", dr.MeanDrift),
		"max_drift_1st", fmt.Sprintf("%.4f (state=%s)", dr.MaxDrift, StateName(dr.MaxState)),
		"states_used", dr.StatesUsed,
		"mean_drift_2nd", fmt.Sprintf("%.4f", drift2),
		"drifted", dr.Drifted,
		"threshold", DriftThreshold,
	)

	if dr.Drifted {
		slog.Warn("btc_strategy.DRIFT_ALERT",
			"mean_drift", fmt.Sprintf("%.4f", dr.MeanDrift),
			"max_state", StateName(dr.MaxState),
			"max_drift", fmt.Sprintf("%.4f", dr.MaxDrift),
			"action", "recent window shows different transition dynamics vs full history",
		)
	}

	return dr
}
