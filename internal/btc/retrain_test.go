package btc

import (
	"math"
	"testing"
	"time"
)

func TestKLDivergence_Identical(t *testing.T) {
	p := []float64{0.5, 0.3, 0.2}
	kl := KLDivergence(p, p)
	if kl > 1e-8 {
		t.Fatalf("KL of identical distributions should be ~0, got %f", kl)
	}
}

func TestKLDivergence_Different(t *testing.T) {
	p := []float64{0.9, 0.05, 0.05}
	q := []float64{0.1, 0.45, 0.45}
	kl := KLDivergence(p, q)
	if kl < 0.5 {
		t.Fatalf("KL of very different distributions should be large, got %f", kl)
	}
}

func TestSymmetricKL(t *testing.T) {
	p := []float64{0.7, 0.2, 0.1}
	q := []float64{0.3, 0.4, 0.3}
	skl := SymmetricKL(p, q)
	if skl <= 0 {
		t.Fatal("symmetric KL should be positive")
	}
	// Symmetric: should equal (KL(p||q) + KL(q||p)) / 2
	expected := (KLDivergence(p, q) + KLDivergence(q, p)) / 2.0
	if math.Abs(skl-expected) > 1e-10 {
		t.Fatalf("expected %f, got %f", expected, skl)
	}
}

func TestMatrixDrift_Same(t *testing.T) {
	candles := makeSyntheticCandles(200)
	tm, _ := Train(candles)
	dr := MatrixDrift(&tm, &tm)
	if dr.MeanDrift > 1e-8 {
		t.Fatalf("drift of same matrix should be ~0, got %f", dr.MeanDrift)
	}
	if dr.Drifted {
		t.Fatal("should not flag drift for same matrix")
	}
}

func TestMatrixDrift_DifferentWindows(t *testing.T) {
	candles := makeSyntheticCandles(500)
	tmFull, _ := Train(candles)
	recent := candles[len(candles)-100:]
	tmRecent, _ := Train(recent)
	dr := MatrixDrift(&tmRecent, &tmFull)
	if dr.StatesUsed == 0 {
		t.Fatal("should have states with data in both")
	}
	// Some drift is expected between windows
	t.Logf("drift: mean=%.4f max=%.4f states=%d drifted=%v",
		dr.MeanDrift, dr.MaxDrift, dr.StatesUsed, dr.Drifted)
}

func TestCheckDrift(t *testing.T) {
	candles := makeSyntheticCandles(300)
	dr := CheckDrift(candles, 100)
	if dr.StatesUsed == 0 {
		t.Fatal("should have states")
	}
	t.Logf("drift: mean=%.4f max=%.4f drifted=%v", dr.MeanDrift, dr.MaxDrift, dr.Drifted)
}

func makeSyntheticCandles(n int) []Candle {
	candles := make([]Candle, n)
	price := 50000.0
	for i := range candles {
		delta := float64(i%7-3) * 0.3
		candles[i] = Candle{
			Timestamp: time.Unix(int64(i)*3600, 0),
			Close:     price * (1 + delta/100),
			Volume:    1000 + float64(i%5)*200,
		}
	}
	return candles
}
