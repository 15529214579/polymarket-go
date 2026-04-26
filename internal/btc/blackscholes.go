package btc

import "math"

// FirstPassageProb computes the probability that BTC will touch strike K at
// any point before time T (years), given current spot S and annualized
// volatility sigma. Uses GBM first-passage (barrier) probability with zero
// drift assumption (conservative).
//
// For "reach" markets (K > S): P(max(S_t) > K, t∈[0,T])
// For "dip" markets (K < S): P(min(S_t) < K, t∈[0,T])
func FirstPassageProb(spot, strike, sigma, yearsToExpiry float64) float64 {
	if sigma <= 0 || yearsToExpiry <= 0 || spot <= 0 || strike <= 0 {
		return 0
	}

	if strike > spot {
		return touchAbove(spot, strike, sigma, yearsToExpiry)
	}
	return touchBelow(spot, strike, sigma, yearsToExpiry)
}

// touchAbove: P(max(S_t) > K) for K > S using GBM first-passage.
// Log-price Y_t = (mu - sigma^2/2)*t + sigma*W_t, barrier b = ln(K/S) > 0.
// With mu=0: drift m = -sigma^2/2.
// P = N((mT - b)/(sigma*sqrt(T))) + exp(2mb/sigma^2) * N((-mT - b)/(sigma*sqrt(T)))
func touchAbove(s, k, sigma, t float64) float64 {
	b := math.Log(k / s)
	m := -0.5 * sigma * sigma
	sqrtT := math.Sqrt(t)
	sigSqrtT := sigma * sqrtT

	term1 := normCDF((m*t - b) / sigSqrtT)
	expFactor := math.Exp(2 * m * b / (sigma * sigma))
	term2 := expFactor * normCDF((-m*t-b)/sigSqrtT)

	p := term1 + term2
	if p > 1 {
		p = 1
	}
	if p < 0 {
		p = 0
	}
	return p
}

// touchBelow: P(min(S_t) < K) for K < S using GBM first-passage.
// barrier c = ln(K/S) < 0.
// P = N((c - mT)/(sigma*sqrt(T))) + exp(2mc/sigma^2) * N((c + mT)/(sigma*sqrt(T)))
func touchBelow(s, k, sigma, t float64) float64 {
	if k >= s {
		return 1.0
	}
	c := math.Log(k / s) // < 0
	m := -0.5 * sigma * sigma
	sqrtT := math.Sqrt(t)
	sigSqrtT := sigma * sqrtT

	term1 := normCDF((c - m*t) / sigSqrtT)
	expFactor := math.Exp(2 * m * c / (sigma * sigma))
	term2 := expFactor * normCDF((c+m*t)/sigSqrtT)

	p := term1 + term2
	if p > 1 {
		p = 1
	}
	if p < 0 {
		p = 0
	}
	return p
}

// HistoricalVolatility computes annualized volatility from hourly candles
// using close-to-close log returns.
func HistoricalVolatility(candles []Candle) float64 {
	if len(candles) < 2 {
		return 0
	}
	var sum, sumSq float64
	n := 0
	for i := 1; i < len(candles); i++ {
		if candles[i-1].Close <= 0 || candles[i].Close <= 0 {
			continue
		}
		r := math.Log(candles[i].Close / candles[i-1].Close)
		sum += r
		sumSq += r * r
		n++
	}
	if n < 2 {
		return 0
	}
	mean := sum / float64(n)
	variance := sumSq/float64(n) - mean*mean
	if variance < 0 {
		variance = 0
	}
	hourlyVol := math.Sqrt(variance)
	return hourlyVol * math.Sqrt(8760)
}

// EWMAVolatility computes annualized volatility using Exponentially Weighted
// Moving Average with the given lambda (0 < lambda < 1). Lambda=0.94 is the
// RiskMetrics standard. More responsive to recent volatility clusters than
// fixed-window historical vol.
func EWMAVolatility(candles []Candle, lambda float64) float64 {
	if len(candles) < 2 {
		return 0
	}
	if lambda <= 0 || lambda >= 1 {
		lambda = 0.94
	}

	r0 := math.Log(candles[1].Close / candles[0].Close)
	ewmaVar := r0 * r0

	for i := 2; i < len(candles); i++ {
		if candles[i-1].Close <= 0 || candles[i].Close <= 0 {
			continue
		}
		r := math.Log(candles[i].Close / candles[i-1].Close)
		ewmaVar = lambda*ewmaVar + (1-lambda)*r*r
	}

	hourlyVol := math.Sqrt(ewmaVar)
	return hourlyVol * math.Sqrt(8760)
}

// VolSmileAdjust applies a simple volatility smile: strikes far from spot
// get higher implied vol (fat tails). The adjustment is linear in log-moneyness.
func VolSmileAdjust(baseVol, spot, strike float64) float64 {
	if spot <= 0 || strike <= 0 || baseVol <= 0 {
		return baseVol
	}
	logMoneyness := math.Abs(math.Log(strike / spot))
	skewFactor := 1.0 + 0.5*logMoneyness
	return baseVol * skewFactor
}

// normCDF approximates the standard normal cumulative distribution function.
func normCDF(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

// BSGap represents the gap between first-passage fair value and PM price.
type BSGap struct {
	Strike    float64
	Question  string
	PMPrice   float64 // PM Yes price
	BSProb    float64 // First-passage probability
	GapPP     float64 // (BSProb - PMPrice) * 100
	Direction string  // BUY_YES (PM underpriced) or BUY_NO (PM overpriced)
	EdgeRatio float64 // |gap| / PMPrice — relative edge size
}

// FindBSGaps compares first-passage probabilities against PM prices for
// all BTC markets, returning opportunities where gap exceeds minGapPP.
// Uses vol smile adjustment: strikes far from spot get higher implied vol.
func FindBSGaps(markets []PMMarket, spot, sigma, yearsToExpiry, minGapPP float64) []BSGap {
	var gaps []BSGap
	for _, m := range markets {
		if m.Strike <= 0 || m.YesPrice <= 0 {
			continue
		}
		if m.YesPrice >= 0.99 || m.YesPrice <= 0.01 {
			continue
		}

		adjustedSigma := VolSmileAdjust(sigma, spot, m.Strike)
		bsProb := FirstPassageProb(spot, m.Strike, adjustedSigma, yearsToExpiry)
		gapPP := (bsProb - m.YesPrice) * 100

		if math.Abs(gapPP) < minGapPP {
			continue
		}

		dir := "BUY_YES"
		if gapPP < 0 {
			dir = "BUY_NO"
		}

		edgeRatio := 0.0
		if m.YesPrice > 0 {
			edgeRatio = math.Abs(gapPP) / (m.YesPrice * 100)
		}

		gaps = append(gaps, BSGap{
			Strike:    m.Strike,
			Question:  m.Question,
			PMPrice:   m.YesPrice,
			BSProb:    bsProb,
			GapPP:     gapPP,
			Direction: dir,
			EdgeRatio: edgeRatio,
		})
	}

	for i := 0; i < len(gaps); i++ {
		for j := i + 1; j < len(gaps); j++ {
			if math.Abs(gaps[j].GapPP) > math.Abs(gaps[i].GapPP) {
				gaps[i], gaps[j] = gaps[j], gaps[i]
			}
		}
	}
	return gaps
}
