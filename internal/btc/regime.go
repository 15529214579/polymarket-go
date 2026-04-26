package btc

import "math"

// RegimeDirectionBias returns a multiplier [0.5, 1.5] that adjusts signal
// strength based on HMM regime and multi-TF direction alignment.
//
// TREND regime:
//   - bull trend + dip market BUY_NO → amplify (trend away from dip)
//   - bull trend + reach market BUY_YES → amplify (trend toward reach)
//   - bear trend + dip market BUY_YES → dampen (bear ≠ buy dip)
//   - bear trend + reach market BUY_NO → amplify
//
// MEAN_REVERT regime: contrarian — opposite of trend logic
// VOLATILE regime: reduce all signals (high uncertainty)
func RegimeDirectionBias(regime int, regimeConf float64, alignment string, direction string, isReach bool) float64 {
	if regimeConf < 0.4 {
		return 1.0
	}

	switch regime {
	case RegimeTrend:
		return trendBias(alignment, direction, isReach, regimeConf)
	case RegimeMR:
		return mrBias(alignment, direction, isReach, regimeConf)
	case RegimeVolat:
		return math.Max(0.5, 1.0-regimeConf*0.5)
	default:
		return 1.0
	}
}

func trendBias(alignment, direction string, isReach bool, conf float64) float64 {
	boost := 1.0 + conf*0.3

	switch alignment {
	case "ALIGNED_BULL":
		if isReach && direction == "BUY_YES" {
			return boost
		}
		if !isReach && direction == "BUY_NO" {
			return boost
		}
		return 1.0 / boost
	case "ALIGNED_BEAR":
		if !isReach && direction == "BUY_YES" {
			return boost
		}
		if isReach && direction == "BUY_NO" {
			return boost
		}
		return 1.0 / boost
	default:
		return 1.0
	}
}

func mrBias(alignment, direction string, isReach bool, conf float64) float64 {
	boost := 1.0 + conf*0.2
	switch alignment {
	case "ALIGNED_BULL":
		if !isReach && direction == "BUY_YES" {
			return boost
		}
		if isReach && direction == "BUY_NO" {
			return boost
		}
		return 1.0
	case "ALIGNED_BEAR":
		if isReach && direction == "BUY_YES" {
			return boost
		}
		if !isReach && direction == "BUY_NO" {
			return boost
		}
		return 1.0
	default:
		return 1.0
	}
}
