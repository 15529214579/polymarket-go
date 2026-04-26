package btc

import (
	"testing"
)

func TestDefaultHMMPriors(t *testing.T) {
	m := DefaultHMMPriors()

	// Check Pi sums to ~1
	var piSum float64
	for _, p := range m.Pi {
		piSum += p
	}
	if piSum < 0.99 || piSum > 1.01 {
		t.Fatalf("Pi sums to %f, want ~1.0", piSum)
	}

	// Check each transition row sums to ~1
	for i := 0; i < NRegimes; i++ {
		var sum float64
		for j := 0; j < NRegimes; j++ {
			sum += m.Trans[i][j]
		}
		if sum < 0.99 || sum > 1.01 {
			t.Errorf("Trans[%d] sums to %f, want ~1.0", i, sum)
		}
	}

	// Check each emission row sums to ~1
	for i := 0; i < NRegimes; i++ {
		var sum float64
		for j := 0; j < NStates; j++ {
			sum += m.Emit[i][j]
		}
		if sum < 0.99 || sum > 1.01 {
			t.Errorf("Emit[%d] sums to %f, want ~1.0", i, sum)
		}
	}
}

func TestTrainHMMSmall(t *testing.T) {
	obs := make([]int, 100)
	for i := range obs {
		obs[i] = i % NStates
	}

	m := TrainHMM(obs, 5)

	// Check trained model still has valid probabilities
	for i := 0; i < NRegimes; i++ {
		var sum float64
		for j := 0; j < NRegimes; j++ {
			sum += m.Trans[i][j]
		}
		if sum < 0.98 || sum > 1.02 {
			t.Errorf("trained Trans[%d] sums to %f", i, sum)
		}
	}
}

func TestViterbiReturnsPath(t *testing.T) {
	m := DefaultHMMPriors()
	obs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	path := Viterbi(m, obs)
	if len(path) != len(obs) {
		t.Fatalf("path len %d, want %d", len(path), len(obs))
	}

	for i, r := range path {
		if r < 0 || r >= NRegimes {
			t.Errorf("path[%d] = %d, out of range", i, r)
		}
	}
}

func TestCandlesToObservations(t *testing.T) {
	candles := make([]Candle, 30)
	for i := range candles {
		candles[i] = Candle{
			Close:  float64(100 + i),
			Volume: float64(1000 + i*10),
		}
	}

	obs := CandlesToObservations(candles)
	if len(obs) == 0 {
		t.Fatal("expected non-empty observations")
	}

	for i, o := range obs {
		if o < 0 || o >= NStates {
			t.Errorf("obs[%d] = %d, out of range", i, o)
		}
	}
}

func TestHMMRegimeFilter(t *testing.T) {
	if !HMMRegimeFilter(RegimeTrend, 0.7) {
		t.Error("TREND with conf 0.7 should pass")
	}
	if HMMRegimeFilter(RegimeTrend, 0.3) {
		t.Error("TREND with conf 0.3 should not pass")
	}
	if HMMRegimeFilter(RegimeVolat, 0.9) {
		t.Error("VOLATILE should never pass")
	}
	if !HMMRegimeFilter(RegimeMR, 0.8) {
		t.Error("MEAN_REVERT with conf 0.8 should pass")
	}
	if HMMRegimeFilter(RegimeMR, 0.5) {
		t.Error("MEAN_REVERT with conf 0.5 should not pass")
	}
}

func TestEWMAVolatility(t *testing.T) {
	candles := make([]Candle, 100)
	candles[0] = Candle{Close: 100}
	for i := 1; i < 100; i++ {
		delta := 0.5
		if i%2 == 0 {
			delta = -0.5
		}
		candles[i] = Candle{Close: candles[i-1].Close + delta}
	}

	vol := EWMAVolatility(candles, 0.94)
	if vol <= 0 {
		t.Fatalf("EWMA vol should be positive, got %f", vol)
	}

	histVol := HistoricalVolatility(candles)
	if histVol <= 0 {
		t.Fatalf("hist vol should be positive, got %f", histVol)
	}
}

func TestBlendedVolatility(t *testing.T) {
	candles := make([]Candle, 200)
	candles[0] = Candle{Close: 100}
	for i := 1; i < 200; i++ {
		delta := 0.3
		if i%2 == 0 {
			delta = -0.3
		}
		if i > 150 {
			delta *= 0.1
		}
		candles[i] = Candle{Close: candles[i-1].Close + delta}
	}

	blended := BlendedVolatility(candles, 0.94, 0.6)
	if blended < 0.25 {
		t.Fatalf("blended vol should be >= floor 0.25, got %f", blended)
	}

	ewma := EWMAVolatility(candles, 0.94)
	hist := HistoricalVolatility(candles)
	if ewma >= hist && blended < ewma {
		t.Fatalf("blended should be >= EWMA when EWMA >= hist")
	}
	if ewma < hist && blended <= ewma {
		t.Fatalf("blended %f should be > pure EWMA %f when hist %f is higher", blended, ewma, hist)
	}
}

func TestBlendedVolFloor(t *testing.T) {
	candles := make([]Candle, 100)
	candles[0] = Candle{Close: 100}
	for i := 1; i < 100; i++ {
		candles[i] = Candle{Close: candles[i-1].Close + 0.001}
	}
	blended := BlendedVolatility(candles, 0.94, 0.6)
	if blended < 0.25 {
		t.Fatalf("blended vol floor not applied, got %f", blended)
	}
}

func TestVolSmileAdjust(t *testing.T) {
	base := 0.50

	// ATM should be close to base vol
	atm := VolSmileAdjust(base, 100, 100)
	if atm != base {
		t.Errorf("ATM vol should equal base: got %f", atm)
	}

	// OTM should be higher
	otm := VolSmileAdjust(base, 100, 50)
	if otm <= base {
		t.Errorf("OTM vol should be > base: got %f", otm)
	}

	// Far OTM should be even higher
	farOTM := VolSmileAdjust(base, 100, 10)
	if farOTM <= otm {
		t.Errorf("far OTM vol should be > OTM: got %f vs %f", farOTM, otm)
	}
}
