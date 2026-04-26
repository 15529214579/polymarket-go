package btc

import (
	"testing"
)

func TestMomentumAlert(t *testing.T) {
	a := MomentumAlert{
		Move1h:   2.0,
		Move15m:  0.5,
		Spot:     80000,
		OpenHour: 78431,
		IsSharp:  true,
	}
	if !a.IsSharp {
		t.Error("2% 1h move should be sharp")
	}

	b := MomentumAlert{
		Move1h:  0.3,
		Move15m: 0.1,
		IsSharp: false,
	}
	if b.IsSharp {
		t.Error("0.3% move should not be sharp")
	}
}

func TestNewMomentumWatcher(t *testing.T) {
	called := false
	w := NewMomentumWatcher(func() { called = true })
	if w.checkInterval == 0 {
		t.Error("checkInterval should be set")
	}
	if w.onSharpMove == nil {
		t.Error("callback should be set")
	}
	w.onSharpMove()
	if !called {
		t.Error("callback should have been called")
	}
}
