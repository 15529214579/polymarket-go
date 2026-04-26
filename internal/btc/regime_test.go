package btc

import "testing"

func TestRegimeDirectionBias_TrendBull(t *testing.T) {
	bias := RegimeDirectionBias(RegimeTrend, 0.8, "ALIGNED_BULL", "BUY_NO", false)
	if bias <= 1.0 {
		t.Errorf("trend+bull+dip BUY_NO should amplify, got %.2f", bias)
	}
}

func TestRegimeDirectionBias_TrendBullContra(t *testing.T) {
	bias := RegimeDirectionBias(RegimeTrend, 0.8, "ALIGNED_BULL", "BUY_YES", false)
	if bias >= 1.0 {
		t.Errorf("trend+bull+dip BUY_YES should dampen, got %.2f", bias)
	}
}

func TestRegimeDirectionBias_Volatile(t *testing.T) {
	bias := RegimeDirectionBias(RegimeVolat, 0.9, "ALIGNED_BULL", "BUY_YES", true)
	if bias >= 1.0 {
		t.Errorf("volatile regime should dampen all, got %.2f", bias)
	}
}

func TestRegimeDirectionBias_LowConf(t *testing.T) {
	bias := RegimeDirectionBias(RegimeTrend, 0.3, "ALIGNED_BULL", "BUY_YES", true)
	if bias != 1.0 {
		t.Errorf("low confidence should return 1.0, got %.2f", bias)
	}
}

func TestRegimeDirectionBias_MRContrarian(t *testing.T) {
	bias := RegimeDirectionBias(RegimeMR, 0.7, "ALIGNED_BULL", "BUY_NO", true)
	if bias <= 1.0 {
		t.Errorf("mean-revert+bull should amplify contrarian reach BUY_NO, got %.2f", bias)
	}
}
