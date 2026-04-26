package btc

import "testing"

func TestScoreSignal_HighGap(t *testing.T) {
	s := ScoreSignal(-30.0, 1.2, 1.3, "ALIGNED_BEAR", 0.8, 0.7)
	if s.Total < 70 {
		t.Errorf("high gap + aligned + good regime should score >70, got %d", s.Total)
	}
	if s.Tier != "AUTO" && s.Tier != "SIGNAL" {
		t.Errorf("should be AUTO or SIGNAL tier, got %s", s.Tier)
	}
}

func TestScoreSignal_LowGap(t *testing.T) {
	s := ScoreSignal(-7.0, 1.0, 1.0, "MIXED", 0.5, 0.1)
	if s.Total >= 60 {
		t.Errorf("low gap + neutral should score <60, got %d", s.Total)
	}
	if s.Tier != "LOG" {
		t.Errorf("should be LOG tier, got %s", s.Tier)
	}
}

func TestScoreSignal_Tiers(t *testing.T) {
	auto := ScoreSignal(-35.0, 1.2, 1.3, "ALIGNED_BEAR", 0.9, 1.0)
	if auto.Tier != "AUTO" {
		t.Errorf("max inputs should be AUTO, got %s (score=%d)", auto.Tier, auto.Total)
	}

	log := ScoreSignal(-7.0, 0.8, 0.7, "MIXED", 0.3, 0.05)
	if log.Tier != "LOG" {
		t.Errorf("min inputs should be LOG, got %s (score=%d)", log.Tier, log.Total)
	}
}

func TestScoreSignal_Components(t *testing.T) {
	s := ScoreSignal(-20.0, 1.1, 1.15, "ALIGNED_BULL", 0.7, 0.5)
	if s.GapScore <= 0 {
		t.Error("gap score should be > 0")
	}
	if s.RegimeScore <= 0 {
		t.Error("regime score should be > 0")
	}
	if s.TFScore <= 0 {
		t.Error("TF score should be > 0")
	}
	if s.EdgeScore <= 0 {
		t.Error("edge score should be > 0")
	}
}
