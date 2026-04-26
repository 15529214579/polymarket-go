package btc

import (
	"context"
	"fmt"
	"log/slog"
)

// MultiTFPrediction combines Markov predictions from multiple timeframes.
type MultiTFPrediction struct {
	TF5m  *Prediction // nil if unavailable
	TF15m *Prediction
	TF1h  *Prediction

	CombinedBull   float64 // weighted bull probability [0,1]
	CombinedBear   float64 // weighted bear probability [0,1]
	CombinedReturn float64 // weighted expected return
	Alignment      string  // "ALIGNED_BULL" / "ALIGNED_BEAR" / "MIXED"
	Confidence     float64 // 0-1, higher = more timeframes agree
}

const (
	weight5m  = 0.20
	weight15m = 0.30
	weight1h  = 0.50
)

// PredictMultiTF fetches candles at 5m/15m/1h, trains per-TF Markov models,
// and returns a combined prediction. Any unavailable timeframe is skipped
// (its weight is redistributed).
func PredictMultiTF(ctx context.Context) (*MultiTFPrediction, error) {
	type tfResult struct {
		interval Interval
		candles  []Candle
		pred     Prediction
		ok       bool
	}

	intervals := []Interval{Interval5m, Interval15m, Interval1h}
	limits := []int{1000, 1000, 720}

	results := make([]tfResult, len(intervals))

	for i, iv := range intervals {
		candles, err := FetchCandles(ctx, "BTCUSDT", iv, limits[i])
		if err != nil {
			slog.Warn("multi_tf.fetch_fail", "interval", iv, "err", err.Error())
			continue
		}
		if len(candles) < 30 {
			slog.Warn("multi_tf.insufficient", "interval", iv, "count", len(candles))
			continue
		}

		tm, _ := Train(candles)
		retStats := BuildReturnStats(candles)
		pred, ok := PredictFromCandles(candles, &tm, retStats)
		results[i] = tfResult{interval: iv, candles: candles, pred: pred, ok: ok}
	}

	mtp := &MultiTFPrediction{}

	type weighted struct {
		pred   Prediction
		weight float64
	}

	var active []weighted
	weights := []float64{weight5m, weight15m, weight1h}

	for i, r := range results {
		if !r.ok {
			continue
		}
		p := r.pred
		switch intervals[i] {
		case Interval5m:
			mtp.TF5m = &p
		case Interval15m:
			mtp.TF15m = &p
		case Interval1h:
			mtp.TF1h = &p
		}
		active = append(active, weighted{pred: p, weight: weights[i]})
	}

	if len(active) == 0 {
		return nil, fmt.Errorf("no timeframe data available")
	}

	var totalW float64
	for _, a := range active {
		totalW += a.weight
	}
	for _, a := range active {
		w := a.weight / totalW
		mtp.CombinedBull += a.pred.BullProb * w
		mtp.CombinedBear += a.pred.BearProb * w
		mtp.CombinedReturn += a.pred.ExpectedReturn * w
	}

	bullCount, bearCount := 0, 0
	for _, a := range active {
		if a.pred.BullProb > a.pred.BearProb {
			bullCount++
		} else if a.pred.BearProb > a.pred.BullProb {
			bearCount++
		}
	}

	switch {
	case bullCount == len(active):
		mtp.Alignment = "ALIGNED_BULL"
		mtp.Confidence = mtp.CombinedBull
	case bearCount == len(active):
		mtp.Alignment = "ALIGNED_BEAR"
		mtp.Confidence = mtp.CombinedBear
	default:
		mtp.Alignment = "MIXED"
		majority := float64(max(bullCount, bearCount)) / float64(len(active))
		mtp.Confidence = majority * 0.5
	}

	slog.Info("multi_tf.prediction",
		"alignment", mtp.Alignment,
		"combined_bull", fmt.Sprintf("%.1f%%", mtp.CombinedBull*100),
		"combined_bear", fmt.Sprintf("%.1f%%", mtp.CombinedBear*100),
		"combined_return", fmt.Sprintf("%.3f%%", mtp.CombinedReturn),
		"confidence", fmt.Sprintf("%.2f", mtp.Confidence),
		"tf_count", len(active),
	)
	return mtp, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MultiTFEntryFilter returns true if the multi-timeframe signal supports
// entering a position in the given direction. Only blocks when there is
// strong directional conflict (all timeframes aligned against the trade).
// MIXED alignment always passes — BS gap is a long-term structural edge,
// short-term direction being unclear is not a reason to block.
func (m *MultiTFPrediction) MultiTFEntryFilter(direction string) bool {
	switch direction {
	case "BUY_YES":
		return m.Alignment != "ALIGNED_BEAR" || m.Confidence < 0.55
	case "BUY_NO":
		return m.Alignment != "ALIGNED_BULL" || m.Confidence < 0.55
	default:
		return true
	}
}
