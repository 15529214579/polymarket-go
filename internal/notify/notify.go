// Package notify pushes key events (risk breaker trips, large fills, daily
// reports) to the boss via Telegram. All sends are best-effort and async —
// a failed delivery must never block trading.
//
// The runtime picks Telegram when TELEGRAM_BOT_TOKEN + TELEGRAM_CHAT_ID are
// set in the environment; otherwise it uses the Nop notifier (tests + offline).
package notify

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Notifier is the narrow surface called from the trading loop. Implementations
// must be safe for concurrent use and must not block the caller for more than
// a few milliseconds (use goroutines + channels internally).
type Notifier interface {
	RiskTrip(ev RiskTripEvent)
	RiskResume(ev RiskResumeEvent)
	LargeFill(ev LargeFillEvent)
	// SignalPrompt pushes a signal DM with one inline-keyboard row per outcome
	// (Phase 3.5 UX). Each row is Buy 1U / 5U / 10U; callback_data is
	// "buy:<nonce>:<slot>:<sizeUSD>" where slot indexes into PendingIntent.Choices.
	SignalPrompt(ev SignalPromptEvent)
	Close(ctx context.Context) error
}

// RiskTripEvent carries the state snapshot at the moment the breaker flipped.
type RiskTripEvent struct {
	Reason        string  // daily_loss | feed_silence | manual_pause
	DayPnLUSD     float64 // negative for losses
	DayLossCapUSD float64 // positive number — the cap magnitude
	SilentSec     int     // feed silence duration (0 if N/A)
	OpenPositions int
}

// RiskResumeEvent reports when the breaker has been cleared. We surface this
// so the boss has a definite "trading is live again" signal — especially
// after a manual Pause/Resume cycle.
type RiskResumeEvent struct {
	PrevReason    string
	DayPnLUSD     float64
	DayLossCapUSD float64
}

// SignalPromptEvent is the payload for a new momentum signal needing a
// manual sizing + outcome decision. Nonce must match the PendingStore entry
// that holds the full Choices slice.
type SignalPromptEvent struct {
	Nonce   string
	Match   string // "LoL: Shifters vs G2 Esports"
	Context string // "Game 1 Winner" / "BO3 · LCK Rounds 1-2" / ""
	EndIn   string // human-readable "1h 23m" until market close; "" if unknown

	// Choices are rendered one row per entry (up to ~6 total). Slot = array idx
	// and encoded in callback_data. One Choice must have IsSignal=true — it
	// drives the DM title.
	Choices []SignalChoice

	// Signal context (describes the firing asset; shown in the subtitle).
	DeltaPP  float64
	TailUps  int
	TailLen  int
	BuyRatio float64

	SizesUSD  []float64     // default {1, 5, 10}
	ExpiresIn time.Duration // visual-only hint in the DM body
}

// SignalChoice is one selectable outcome rendered as a button row.
type SignalChoice struct {
	Slot     int
	Outcome  string
	Mid      float64
	IsSignal bool
}

// LargeFillEvent carries a closed-position summary worth surfacing in DM.
// Typically fired when |pnl| crosses a config threshold (see SPEC §6).
type LargeFillEvent struct {
	Question string
	AssetID  string
	Side     string // buy | sell (the close side)
	SizeUSD  float64
	PnLUSD   float64 // realized, negative for losses
	EntryPx  float64
	ExitPx   float64
	Reason   string // reversal_ticks | reversal_drawdown | stop_loss | timeout
	HeldSec  int
}

// ---- formatting helpers (exported so telegram_test can assert them) ----

// FormatRiskTrip renders a single-line title + body block. The body intentionally
// avoids markdown so we don't have to escape Telegram special chars.
func FormatRiskTrip(ev RiskTripEvent) string {
	var title string
	switch ev.Reason {
	case "daily_loss":
		title = "🔴 风控熔断：日亏损上限"
	case "feed_silence":
		title = "🟠 风控熔断：WSS 喂价静默"
	case "manual_pause":
		title = "⏸ 风控暂停：手动"
	default:
		title = "⚠️ 风控熔断：" + ev.Reason
	}
	var b strings.Builder
	fmt.Fprintln(&b, title)
	fmt.Fprintf(&b, "day pnl:  %+.2f USDC (cap %.2f)\n", ev.DayPnLUSD, ev.DayLossCapUSD)
	if ev.SilentSec > 0 {
		fmt.Fprintf(&b, "feed silent: %ds\n", ev.SilentSec)
	}
	if ev.OpenPositions > 0 {
		fmt.Fprintf(&b, "open positions: %d (exit 逻辑会各自跑平)\n", ev.OpenPositions)
	}
	fmt.Fprintln(&b, "\n老板手动恢复前不再开新仓。")
	return b.String()
}

func FormatRiskResume(ev RiskResumeEvent) string {
	return fmt.Sprintf(
		"🟢 风控恢复（之前：%s）\nday pnl: %+.2f USDC (cap %.2f)\n交易已继续。",
		ev.PrevReason, ev.DayPnLUSD, ev.DayLossCapUSD,
	)
}

// FormatSignalPrompt renders the DM body that accompanies the inline-keyboard
// rows. Formatting kept plain (no markdown escaping needed).
//
// Layout:
//
//	⚡ 动量信号 · <signalOutcome> ↑
//	<Match>
//	<Context> · 结算 <EndIn>
//	Δ +x.xxpp · tail 4/5 · buy 78%
//
//	选 <outcome> (当前 0.xxxx):
//	  [按钮行]
//	选 <outcome> (当前 0.xxxx):
//	  [按钮行]
//	按钮 60s 内有效
func FormatSignalPrompt(ev SignalPromptEvent) string {
	var b strings.Builder

	signalOutcome := "?"
	for _, c := range ev.Choices {
		if c.IsSignal {
			signalOutcome = c.Outcome
			break
		}
	}

	fmt.Fprintf(&b, "⚡ 动量信号 · %s ↑\n", signalOutcome)
	match := ev.Match
	if match == "" && len(ev.Choices) > 0 {
		match = ev.Choices[0].Outcome + " vs " + outcomeAt(ev.Choices, 1)
	}
	fmt.Fprintln(&b, truncateStr(match, 120))

	var ctxLine strings.Builder
	if ev.Context != "" {
		ctxLine.WriteString(ev.Context)
	}
	if ev.EndIn != "" {
		if ctxLine.Len() > 0 {
			ctxLine.WriteString(" · ")
		}
		ctxLine.WriteString("结算 ")
		ctxLine.WriteString(ev.EndIn)
	}
	if ctxLine.Len() > 0 {
		fmt.Fprintln(&b, ctxLine.String())
	}

	fmt.Fprintf(&b, "Δ %+.2fpp · tail %d/%d · buy %.0f%%\n",
		ev.DeltaPP, ev.TailUps, ev.TailLen, ev.BuyRatio*100)
	fmt.Fprintln(&b)

	for _, c := range ev.Choices {
		tag := ""
		if c.IsSignal {
			tag = "  ← 信号"
		}
		fmt.Fprintf(&b, "选 %s (当前 %.4f)%s:\n",
			truncateStr(c.Outcome, 40), c.Mid, tag)
	}

	ttl := ev.ExpiresIn
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	fmt.Fprintf(&b, "\n按钮 %ds 内有效", int(ttl.Seconds()))
	return b.String()
}

// DefaultSizesUSD is the Buy 1/5/10 inline-button row used when
// SignalPromptEvent.SizesUSD is empty.
var DefaultSizesUSD = []float64{1, 5, 10}

func FormatLargeFill(ev LargeFillEvent) string {
	tag := "💰 大单平仓"
	if ev.PnLUSD < 0 {
		tag = "📉 大单止损"
	}
	q := ev.Question
	if q == "" {
		q = ev.AssetID
	}
	if len(q) > 80 {
		q = q[:77] + "..."
	}
	return fmt.Sprintf(
		"%s %+.2f USDC\n%s\nentry %.4f → exit %.4f · %s · held %ds · size %.2fU",
		tag, ev.PnLUSD, q, ev.EntryPx, ev.ExitPx, ev.Reason, ev.HeldSec, ev.SizeUSD,
	)
}

// ---- question parsing ----

// ParseMarketTitle splits a Polymarket question into a match title and a
// context suffix. Best-effort: the gamma Question strings we see are of a few
// loose shapes:
//
//	"LoL: Shifters vs G2 Esports - Game 1 Winner"             → ("LoL: Shifters vs G2 Esports", "Game 1 Winner")
//	"LoL: Weibo Gaming vs Oh My God (BO3) - LPL Group Ascend" → ("LoL: Weibo Gaming vs Oh My God", "BO3 · LPL Group Ascend")
//	"Games Total: O/U 2.5"                                    → ("Games Total: O/U 2.5", "")
//	"Game Handicap: BLG (-1.5) vs Invictus Gaming (+1.5)"     → ("Game Handicap: BLG (-1.5) vs Invictus Gaming (+1.5)", "")
//
// The function never panics and always returns match=strings.TrimSpace(input)
// when no split is found.
func ParseMarketTitle(q string) (match, ctx string) {
	q = strings.TrimSpace(q)
	// Find " - " separator (use last to keep hyphens inside team names intact).
	idx := strings.LastIndex(q, " - ")
	if idx < 0 {
		return q, ""
	}
	head := strings.TrimSpace(q[:idx])
	tail := strings.TrimSpace(q[idx+3:])

	// Pull "(BO3)" / "(BO5)" out of head into ctx prefix so the match title
	// stays clean.
	if open := strings.LastIndex(head, "("); open > 0 {
		if close := strings.Index(head[open:], ")"); close > 0 {
			inside := head[open+1 : open+close]
			if isFormatTag(inside) {
				stripped := strings.TrimSpace(head[:open])
				tail = inside + " · " + tail
				head = stripped
			}
		}
	}
	return head, tail
}

// isFormatTag returns true for "BO3"/"BO5"/"BO7"-style inner text.
func isFormatTag(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 3 || len(s) > 4 {
		return false
	}
	if s[0] != 'B' || s[1] != 'O' {
		return false
	}
	for _, r := range s[2:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// HumanizeEndIn renders time-until-end as "2h 05m" / "45m" / "<1m".
// Empty string if end is zero or in the past.
func HumanizeEndIn(now, end time.Time) string {
	if end.IsZero() {
		return ""
	}
	d := end.Sub(now)
	if d <= 0 {
		return ""
	}
	if d < time.Minute {
		return "<1m"
	}
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	if h > 0 {
		return fmt.Sprintf("%dh %02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}

func outcomeAt(cs []SignalChoice, i int) string {
	if i < 0 || i >= len(cs) {
		return "?"
	}
	return cs[i].Outcome
}
