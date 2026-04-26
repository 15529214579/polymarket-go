package btc

import (
	"fmt"
	"log/slog"
	"math"
	"time"
)

type MacroEventType string

const (
	EventFOMC MacroEventType = "FOMC"
	EventCPI  MacroEventType = "CPI"
	EventNFP  MacroEventType = "NFP"
	EventPCE  MacroEventType = "PCE"
)

type MacroEvent struct {
	Type MacroEventType
	Date time.Time // UTC date+time of release
}

type MacroState struct {
	NextEvent     *MacroEvent
	HoursUntil    float64
	VolMultiplier float64 // >1 before events, <1 during cool-down
	Phase         string  // "pre_event", "post_event", "normal"
}

// 2026 macro calendar (UTC times)
// FOMC: 2pm ET = 18:00 UTC (statement release)
// CPI: 8:30am ET = 12:30 UTC
// NFP: 8:30am ET = 12:30 UTC
// PCE: 8:30am ET = 12:30 UTC
var macroCalendar2026 = []MacroEvent{
	// FOMC meetings (2-day, date = announcement day)
	{EventFOMC, time.Date(2026, 1, 29, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 3, 19, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 5, 7, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 6, 18, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 7, 30, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 9, 17, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 10, 29, 18, 0, 0, 0, time.UTC)},
	{EventFOMC, time.Date(2026, 12, 17, 18, 0, 0, 0, time.UTC)},
	// CPI releases (usually ~12th of month)
	{EventCPI, time.Date(2026, 1, 14, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 2, 12, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 3, 11, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 4, 14, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 5, 13, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 6, 10, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 7, 15, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 8, 12, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 9, 16, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 10, 14, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 11, 12, 12, 30, 0, 0, time.UTC)},
	{EventCPI, time.Date(2026, 12, 10, 12, 30, 0, 0, time.UTC)},
	// NFP (first Friday of month)
	{EventNFP, time.Date(2026, 1, 9, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 2, 6, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 3, 6, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 4, 3, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 5, 8, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 6, 5, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 7, 2, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 8, 7, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 9, 4, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 10, 2, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 11, 6, 12, 30, 0, 0, time.UTC)},
	{EventNFP, time.Date(2026, 12, 4, 12, 30, 0, 0, time.UTC)},
	// PCE (last Friday of month, ~10:00 UTC release)
	{EventPCE, time.Date(2026, 1, 30, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 2, 27, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 3, 27, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 4, 30, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 5, 29, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 6, 26, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 7, 31, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 8, 28, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 9, 25, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 10, 30, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 11, 25, 12, 30, 0, 0, time.UTC)},
	{EventPCE, time.Date(2026, 12, 23, 12, 30, 0, 0, time.UTC)},
}

func GetMacroState(now time.Time) MacroState {
	var next *MacroEvent
	for i := range macroCalendar2026 {
		ev := &macroCalendar2026[i]
		if ev.Date.After(now) {
			if next == nil || ev.Date.Before(next.Date) {
				next = ev
			}
		}
	}

	if next == nil {
		return MacroState{Phase: "normal", VolMultiplier: 1.0}
	}

	hours := next.Date.Sub(now).Hours()

	var recent *MacroEvent
	for i := range macroCalendar2026 {
		ev := &macroCalendar2026[i]
		if ev.Date.Before(now) && now.Sub(ev.Date).Hours() < 2 {
			recent = ev
			break
		}
	}

	if recent != nil {
		hoursAfter := now.Sub(recent.Date).Hours()
		cooldown := math.Max(0.85, 1.0-0.15*(1.0-hoursAfter/2.0))
		return MacroState{
			NextEvent:     next,
			HoursUntil:    hours,
			VolMultiplier: cooldown,
			Phase:         "post_event",
		}
	}

	if hours <= 24 {
		var bump float64
		switch next.Type {
		case EventFOMC:
			bump = 0.30
		case EventCPI:
			bump = 0.25
		case EventNFP:
			bump = 0.20
		case EventPCE:
			bump = 0.15
		}
		ramp := 1.0 - hours/24.0
		mult := 1.0 + bump*ramp
		return MacroState{
			NextEvent:     next,
			HoursUntil:    hours,
			VolMultiplier: mult,
			Phase:         "pre_event",
		}
	}

	return MacroState{
		NextEvent:  next,
		HoursUntil: hours,
		VolMultiplier: 1.0,
		Phase:      "normal",
	}
}

func MacroVolAdjust(baseVol float64, macro MacroState) float64 {
	return baseVol * macro.VolMultiplier
}

func LogMacroState(macro MacroState) {
	if macro.NextEvent == nil {
		slog.Info("macro.state", "phase", "normal", "next", "none")
		return
	}
	slog.Info("macro.state",
		"phase", macro.Phase,
		"next_event", string(macro.NextEvent.Type),
		"next_date", macro.NextEvent.Date.Format("2006-01-02 15:04 UTC"),
		"hours_until", fmt.Sprintf("%.1f", macro.HoursUntil),
		"vol_mult", fmt.Sprintf("%.3f", macro.VolMultiplier),
	)
}
