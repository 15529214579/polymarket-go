package btc

import (
	"math"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// State classification tests
// ---------------------------------------------------------------------------

func TestClassifyReturn(t *testing.T) {
	tests := []struct {
		pct  float64
		want int
	}{
		{3.0, RetSurge},
		{2.01, RetSurge},
		{2.0, RetUp},
		{1.0, RetUp},
		{0.51, RetUp},
		{0.5, RetFlat},
		{0.0, RetFlat},
		{-0.5, RetFlat},
		{-0.51, RetDown},
		{-1.5, RetDown},
		{-2.0, RetDown},
		{-2.01, RetCrash},
		{-5.0, RetCrash},
	}
	for _, tt := range tests {
		got := ClassifyReturn(tt.pct)
		if got != tt.want {
			t.Errorf("ClassifyReturn(%.2f) = %d (%s), want %d (%s)",
				tt.pct, got, retNames[got], tt.want, retNames[tt.want])
		}
	}
}

func TestClassifyVolume(t *testing.T) {
	tests := []struct {
		vol, avg float64
		want     int
	}{
		{150, 100, VolHigh}, // ratio=1.5 exactly → boundary → VolMed per >=0.75
		{151, 100, VolHigh},
		{100, 100, VolMed},
		{75, 100, VolMed},
		{74, 100, VolLow},
		{0, 100, VolLow},
		{100, 0, VolMed}, // avg=0 guard
	}
	// boundary: ratio=1.5 → code: ratio>1.5 → VolHigh, else ratio>=0.75 → VolMed
	// fix expected for exactly 1.5
	tests[0].want = VolMed // ratio == 1.5, not > 1.5

	for _, tt := range tests {
		got := ClassifyVolume(tt.vol, tt.avg)
		if got != tt.want {
			t.Errorf("ClassifyVolume(%.0f, %.0f) = %d (%s), want %d (%s)",
				tt.vol, tt.avg, got, volNames[got], tt.want, volNames[tt.want])
		}
	}
}

func TestStateIdxRoundTrip(t *testing.T) {
	for r := 0; r < nRetBuckets; r++ {
		for v := 0; v < nVolRegimes; v++ {
			idx := StateIdx(r, v)
			gotR, gotV := StateComponents(idx)
			if gotR != r || gotV != v {
				t.Errorf("StateIdx(%d,%d)=%d → StateComponents gives (%d,%d)", r, v, idx, gotR, gotV)
			}
		}
	}
}

func TestStateNameDoesNotPanic(t *testing.T) {
	for s := 0; s < NStates; s++ {
		name := StateName(s)
		if name == "" {
			t.Errorf("StateName(%d) returned empty string", s)
		}
	}
	// out-of-range should not panic
	_ = StateName(-1)
	_ = StateName(NStates + 5)
}

// ---------------------------------------------------------------------------
// CandleState tests
// ---------------------------------------------------------------------------

func makeCandle(ts time.Time, open, high, low, close_, volume float64) Candle {
	return Candle{Timestamp: ts, Open: open, High: high, Low: low, Close: close_, Volume: volume}
}

func TestCandleState(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Build 26 candles: constant volume=100, then a big up candle at index 25
	candles := make([]Candle, 26)
	for i := 0; i < 25; i++ {
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), 50000, 50100, 49900, 50000, 100)
	}
	// index 25: +5% surge, double volume
	candles[25] = makeCandle(base.Add(25*time.Hour), 50000, 52600, 49900, 52500, 200)

	state, ok := CandleState(candles, 25)
	if !ok {
		t.Fatal("CandleState returned ok=false for valid candle")
	}
	ret, vol := StateComponents(state)
	if ret != RetSurge {
		t.Errorf("expected RetSurge (pct≈+5%%), got %s", retNames[ret])
	}
	if vol != VolHigh {
		t.Errorf("expected VolHigh (ratio=2.0), got %s", volNames[vol])
	}

	// Test out-of-bounds
	_, ok = CandleState(candles, 0)
	if ok {
		t.Error("CandleState(candles, 0) should return ok=false (need i>=1)")
	}
}

// ---------------------------------------------------------------------------
// Train & transition matrix tests
// ---------------------------------------------------------------------------

func TestTrainEmpty(t *testing.T) {
	tm, n := Train(nil)
	if n != 0 {
		t.Errorf("Train(nil) transitions = %d, want 0", n)
	}
	_ = tm
}

func TestTrainSingleState(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	// All flat candles: each bar closes at exactly the previous bar's close.
	candles := make([]Candle, 50)
	for i := 0; i < 50; i++ {
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), 50000, 50010, 49990, 50000, 100)
	}
	tm, n := Train(candles)
	if n == 0 {
		t.Fatal("expected at least some transitions from flat candles")
	}

	// All transitions should stay in the FLAT/MED or FLAT/LOW state
	for from := 0; from < NStates; from++ {
		for to := 0; to < NStates; to++ {
			count := tm[from][to]
			if count > 0 {
				fromRet, _ := StateComponents(from)
				toRet, _ := StateComponents(to)
				if fromRet != RetFlat || toRet != RetFlat {
					t.Errorf("unexpected transition from %s to %s (count=%d)",
						StateName(from), StateName(to), count)
				}
			}
		}
	}
}

func TestRowProbsSumsToOne(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	// Mix of up/down candles
	candles := make([]Candle, 100)
	price := 50000.0
	for i := 0; i < 100; i++ {
		if i%3 == 0 {
			price *= 1.01
		} else if i%5 == 0 {
			price *= 0.99
		}
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), price, price*1.005, price*0.995, price, 100)
	}
	tm, _ := Train(candles)

	for s := 0; s < NStates; s++ {
		probs := tm.RowProbs(s)
		if probs == nil {
			continue // row with no observations; OK
		}
		sum := 0.0
		for _, p := range probs {
			sum += p
		}
		if math.Abs(sum-1.0) > 1e-9 {
			t.Errorf("RowProbs(%s) sums to %f, want 1.0", StateName(s), sum)
		}
	}
}

// ---------------------------------------------------------------------------
// ReturnStats tests
// ---------------------------------------------------------------------------

func TestBuildReturnStats(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := []Candle{
		makeCandle(base, 100, 101, 99, 100, 100),
		makeCandle(base.Add(time.Hour), 100, 102, 100, 102, 100), // +2% → RetSurge
		makeCandle(base.Add(2*time.Hour), 102, 103, 101, 103, 100), // +0.98% → RetUp
	}
	stats := BuildReturnStats(candles)

	// At index 1: state from candle[0→1] is surge (return=+2%), then
	// actual forward return to candle[2] = (103-102)/102*100 ≈ 0.98%
	// At least one state should have N>0
	total := 0
	for _, rs := range stats {
		total += rs.N
	}
	if total == 0 {
		t.Error("BuildReturnStats returned all-zero stats for non-trivial candles")
	}
}

// ---------------------------------------------------------------------------
// Predict tests
// ---------------------------------------------------------------------------

func TestPredictNoData(t *testing.T) {
	var tm TransitionMatrix
	var retStats [NStates]ReturnStats

	pred := Predict(0, &tm, retStats)
	if pred.CurrentState != 0 {
		t.Errorf("unexpected current state: %d", pred.CurrentState)
	}
	// No transitions → no probs → expected return should be 0
	if pred.ExpectedReturn != 0 {
		t.Errorf("expected return = %f, want 0 when no training data", pred.ExpectedReturn)
	}
}

func TestPredictBullBearSumLE1(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, 50)
	price := 50000.0
	for i := 0; i < 50; i++ {
		if i%2 == 0 {
			price *= 1.005
		} else {
			price *= 0.995
		}
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), price, price*1.001, price*0.999, price, 100+float64(i))
	}
	tm, _ := Train(candles)
	retStats := BuildReturnStats(candles)

	for s := 0; s < NStates; s++ {
		if tm.RowProbs(s) == nil {
			continue
		}
		pred := Predict(s, &tm, retStats)
		if pred.BullProb+pred.BearProb > 1.0+1e-9 {
			t.Errorf("state %s: BullProb+BearProb = %f > 1.0", StateName(s), pred.BullProb+pred.BearProb)
		}
	}
}

func TestPredictFromCandles(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, 30)
	price := 50000.0
	for i := 0; i < 30; i++ {
		price *= 1.001
		candles[i] = makeCandle(base.Add(time.Duration(i)*time.Hour), price, price*1.002, price*0.998, price, 100)
	}
	tm, _ := Train(candles)
	retStats := BuildReturnStats(candles)

	_, ok := PredictFromCandles(candles, &tm, retStats)
	if !ok {
		t.Error("PredictFromCandles returned ok=false for valid candles")
	}

	// Edge case: too few candles
	_, ok = PredictFromCandles(candles[:1], &tm, retStats)
	if ok {
		t.Error("PredictFromCandles should return ok=false for 1-candle slice")
	}
}
