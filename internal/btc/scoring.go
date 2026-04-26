package btc

import "math"

// SignalScore is a 0-100 composite quality score for a BTC signal.
type SignalScore struct {
	Total          int     // 0-100
	GapScore       int     // 0-35: BS gap magnitude
	SentimentScore int     // 0-15: sentiment alignment
	RegimeScore    int     // 0-20: HMM regime support
	TFScore        int     // 0-15: multi-TF alignment
	EdgeScore      int     // 0-15: edge ratio (gap relative to price)
	Tier           string  // "AUTO" (>80) / "SIGNAL" (60-80) / "LOG" (<60)
}

// ScoreSignal computes a composite quality score.
func ScoreSignal(
	gapPP float64,
	sentimentMod float64,
	regimeBias float64,
	alignment string,
	confidence float64,
	edgeRatio float64,
) SignalScore {
	s := SignalScore{}

	absGap := math.Abs(gapPP)
	switch {
	case absGap >= 30:
		s.GapScore = 35
	case absGap >= 20:
		s.GapScore = 25 + int((absGap-20)/10*10)
	case absGap >= 10:
		s.GapScore = 15 + int((absGap-10)/10*10)
	case absGap >= 7:
		s.GapScore = int(absGap / 7 * 15)
	default:
		s.GapScore = int(absGap * 2)
	}

	switch {
	case sentimentMod >= 1.15:
		s.SentimentScore = 15
	case sentimentMod >= 1.05:
		s.SentimentScore = 10
	case sentimentMod >= 0.95:
		s.SentimentScore = 7
	case sentimentMod >= 0.85:
		s.SentimentScore = 3
	default:
		s.SentimentScore = 0
	}

	switch {
	case regimeBias >= 1.2:
		s.RegimeScore = 20
	case regimeBias >= 1.1:
		s.RegimeScore = 15
	case regimeBias >= 0.95:
		s.RegimeScore = 10
	case regimeBias >= 0.8:
		s.RegimeScore = 5
	default:
		s.RegimeScore = 0
	}

	switch alignment {
	case "ALIGNED_BULL", "ALIGNED_BEAR":
		s.TFScore = int(confidence * 15)
	case "MIXED":
		s.TFScore = 7
	default:
		s.TFScore = 5
	}

	switch {
	case edgeRatio >= 1.0:
		s.EdgeScore = 15
	case edgeRatio >= 0.5:
		s.EdgeScore = 10 + int((edgeRatio-0.5)*10)
	case edgeRatio >= 0.2:
		s.EdgeScore = 5 + int((edgeRatio-0.2)/0.3*5)
	default:
		s.EdgeScore = int(edgeRatio / 0.2 * 5)
	}

	s.Total = s.GapScore + s.SentimentScore + s.RegimeScore + s.TFScore + s.EdgeScore
	if s.Total > 100 {
		s.Total = 100
	}

	switch {
	case s.Total >= 80:
		s.Tier = "AUTO"
	case s.Total >= 60:
		s.Tier = "SIGNAL"
	default:
		s.Tier = "LOG"
	}

	return s
}
