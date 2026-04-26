package btc

import "math"

// KellyFraction computes the optimal bet fraction using the Kelly Criterion.
//
// For a binary bet at PM price p (probability of winning = p_fair, payout = 1/p - 1):
//   f* = (b*p_fair - q) / b
//   where b = net payout odds = (1/p - 1), q = 1 - p_fair
//
// We use half-Kelly (f*/2) for safety.
func KellyFraction(pmPrice, fairProb float64) float64 {
	if pmPrice <= 0 || pmPrice >= 1 || fairProb <= 0 || fairProb >= 1 {
		return 0
	}

	b := 1.0/pmPrice - 1.0 // net odds (e.g., buy at 0.40 → b = 1.5)
	q := 1.0 - fairProb

	f := (b*fairProb - q) / b
	if f <= 0 {
		return 0
	}

	// Half-Kelly for safety
	return f / 2.0
}

// KellySizeUSD returns the bet size in USD given a bankroll, PM price,
// and fair probability. Caps at maxBet.
func KellySizeUSD(bankroll, pmPrice, fairProb, maxBet float64) float64 {
	f := KellyFraction(pmPrice, fairProb)
	if f <= 0 {
		return 0
	}

	size := bankroll * f
	if size < 1.0 {
		return 0
	}
	if size > maxBet {
		size = maxBet
	}

	return math.Round(size*100) / 100
}

// ValueEdge computes the edge (in percentage points) of buying at pmPrice
// when fair value is fairProb. Positive = profitable.
func ValueEdge(pmPrice, fairProb float64) float64 {
	return (fairProb - pmPrice) * 100
}

// ExpectedValue computes the expected value per $1 bet.
// EV = fairProb * (1/pmPrice - 1) - (1 - fairProb)
func ExpectedValue(pmPrice, fairProb float64) float64 {
	if pmPrice <= 0 || pmPrice >= 1 {
		return 0
	}
	payout := 1.0/pmPrice - 1.0
	return fairProb*payout - (1 - fairProb)
}
