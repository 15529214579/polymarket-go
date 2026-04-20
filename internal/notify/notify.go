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
	// SignalPrompt pushes a signal DM with inline buttons (Phase 3.5 UX).
	// Telegram impl renders a Buy 1U / 5U / 10U row whose callback_data is
	// "buy:<nonce>:<sizeUSD>". Callers embed the nonce via PendingStore so
	// the inbound callback handler can resolve intent → order.
	SignalPrompt(ev SignalPromptEvent)
	Close(ctx context.Context) error
}

// RiskTripEvent carries the state snapshot at the moment the breaker flipped.
type RiskTripEvent struct {
	Reason         string  // daily_loss | feed_silence | manual_pause
	DayPnLUSD      float64 // negative for losses
	DayLossCapUSD  float64 // positive number — the cap magnitude
	SilentSec      int     // feed silence duration (0 if N/A)
	OpenPositions  int
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
// manual sizing decision. Nonce must match the PendingStore entry.
type SignalPromptEvent struct {
	Nonce    string
	Question string
	AssetID  string
	Mid      float64
	DeltaPP  float64 // 60s delta in pp
	TailUps  int
	TailLen  int
	BuyRatio float64 // 0..1
	ExpiresIn time.Duration // visual-only hint in the DM body
	// SizesUSD defines the inline-button row, left→right. Default {1, 5, 10}.
	SizesUSD []float64
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

// FormatSignalPrompt renders the DM body that accompanies the Buy 1/5/10
// inline-keyboard row. Formatting kept plain (no markdown escaping needed).
func FormatSignalPrompt(ev SignalPromptEvent) string {
	q := ev.Question
	if q == "" {
		q = ev.AssetID
	}
	if len(q) > 80 {
		q = q[:77] + "..."
	}
	ttl := ev.ExpiresIn
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	return fmt.Sprintf(
		"⚡ 动量信号\n%s\nmid %.4f · Δ %+.2fpp · tail %d/%d · buy %.0f%%\n按钮 %ds 内有效",
		q, ev.Mid, ev.DeltaPP, ev.TailUps, ev.TailLen, ev.BuyRatio*100, int(ttl.Seconds()),
	)
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
