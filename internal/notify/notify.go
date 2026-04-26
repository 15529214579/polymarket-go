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
	// SignalPrompt pushes a signal DM with inline-keyboard rows (10U–100U gradient);
	// callback_data is "buy:<nonce>:<slot>:<sizeUSD>:<mode>" where slot indexes
	// into PendingIntent.Choices and mode is "l" (ladder) or "h" (hold).
	SignalPrompt(ev SignalPromptEvent)
	// EditSignalExpired rewrites the original prompt to "⌛ 已过期 · 未下单" and
	// strips the inline keyboard, called when the TTL-reaper evicts an unclicked
	// pending. No-op when messageID == 0 (prompt's send response hadn't landed).
	EditSignalExpired(messageID int64)
	// EditSignalFilled rewrites the original prompt to "✅ 已下单 …" and strips
	// the inline keyboard, called after a successful click.
	EditSignalFilled(ev FillReceiptEvent, messageID int64)
	// FillReceipt pushes a durable DM receipt for a successful manual open
	// (Phase 3.5 C — "成交凭据留档", complement to the callback toast which is
	// ephemeral).
	FillReceipt(ev FillReceiptEvent)
	// InjuryAlert pushes an NBA injury alert DM. Guarded by -injury_enabled flag;
	// to remove: delete this method + InjuryAlertEvent + FormatInjuryAlert.
	InjuryAlert(ev InjuryAlertEvent)
	// WhaleAlert pushes a smart-money whale trade notification. Guarded by
	// -whale_enabled flag; to remove: delete this method + WhaleAlertEvent +
	// FormatWhaleAlert.
	WhaleAlert(ev WhaleAlertEvent)
	// ClosePrompt pushes a DM with inline close/ignore buttons when a whale
	// sells an asset we hold. The boss clicks to confirm the close.
	ClosePrompt(ev ClosePromptEvent)
	// EditCloseDone rewrites the close prompt to show the result and strips
	// the inline keyboard.
	EditCloseDone(text string, messageID int64)
	Close(ctx context.Context) error
}

// RiskTripEvent carries the state snapshot at the moment the breaker flipped.
type RiskTripEvent struct {
	Reason        string  // daily_loss | drawdown | feed_silence | manual_pause
	DayPnLUSD     float64 // negative for losses
	DayLossCapUSD float64 // positive number — the cap magnitude
	DrawdownUSD   float64 // current drawdown from peak (0 if N/A)
	DrawdownCap   float64 // max allowed drawdown (0 if N/A)
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

	Slug       string // market slug for newshare link
	WhaleLabel string // if set, this is a whale-follow signal (shows 🐋 instead of ⚡)
	SizesUSD  []float64     // default {10, 20, ..., 100}
	ExpiresIn time.Duration // visual-only hint in the DM body

	// OnSent, if set, is called asynchronously by the Telegram backend after
	// the prompt has been delivered. messageID is the Telegram message_id of
	// the prompt DM (0 on error). Main uses this to stash the id in the
	// PendingStore so later TTL-expire / click-success edits can target it.
	OnSent func(messageID int64, err error)
}

// SignalChoice is one selectable outcome rendered as a button row.
type SignalChoice struct {
	Slot     int
	Outcome  string
	Mid      float64
	IsSignal bool
}

// FillReceiptEvent is a durable DM after a successful manual_open (Phase 3.5
// click-to-buy) so the boss has a record beyond the transient callback toast.
type FillReceiptEvent struct {
	Question string
	Match    string // optional, nicer header if we have it
	Outcome  string // YES / NO / team name
	SizeUSD  float64
	Units    float64
	FillPx   float64
	OrderID  string
	Source   string // "manual" for click-to-buy; "auto" if we ever use this for auto opens
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

// InjuryAlertEvent carries a star-player injury worth surfacing. Guarded by
// -injury_enabled flag; to remove: delete this type + FormatInjuryAlert +
// InjuryAlert method from Notifier.
type InjuryAlertEvent struct {
	Team       string
	StarPlayer string
	Status     string // "Out" / "Doubtful"
	Reason     string
	Impact     string // "franchise_player_out" / "co_star_out" / "rotation_star_out"
}

// WhaleAlertEvent carries a large trade from a tracked smart-money wallet.
// Guarded by -whale_enabled flag.
type WhaleAlertEvent struct {
	Wallet    string
	Label     string // human-readable whale name (e.g. "drpufferfish", "countryside")
	Side      string // BUY / SELL
	SizeUnits float64
	Price     float64
	Notional  float64
	Market    string
	Outcome   string
	TradeID    string
	LinkURL    string
	ProfileURL string // whale's Polymarket profile page (e.g. https://polymarket.com/@handle)
	Timestamp  time.Time
	// Position context: whale's total holding for this asset (0 if unknown).
	TotalShares float64
	AvgPrice    float64
	PctSold     float64 // percentage of position sold (SELL only, 0-100)
}

// ClosePromptEvent is the payload for a whale-sell close prompt. The boss
// sees matching open positions and decides whether to close or ignore.
type ClosePromptEvent struct {
	Nonce     string
	Market    string // question/title
	Outcome   string
	AssetID   string
	WhaleLabel string // human-readable whale name
	WhaleSize  float64 // whale's sell size in shares
	WhaleNotl  float64 // whale's sell notional USD
	WhalePrice float64
	LinkURL    string
	ProfileURL string // whale's Polymarket profile page
	// Positions lists our open positions matching this asset.
	Positions []ClosePosition
	OnSent    func(messageID int64, err error)
	// Whale's total position context (0 if unknown).
	WhaleTotalShares float64
	WhaleAvgPrice    float64
	WhalePctSold     float64 // percentage of total position being sold (0-100)
}

// ClosePosition is one open position shown in the close prompt.
type ClosePosition struct {
	PosID    string
	SizeUSD  float64
	Units    float64
	EntryMid float64
}

// ---- formatting helpers (exported so telegram_test can assert them) ----

// FormatRiskTrip renders a single-line title + body block. The body intentionally
// avoids markdown so we don't have to escape Telegram special chars.
func FormatRiskTrip(ev RiskTripEvent) string {
	var title string
	switch ev.Reason {
	case "daily_loss":
		title = "🔴 风控熔断：日亏损上限"
	case "drawdown":
		title = "🔴 风控熔断：组合回撤上限"
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
	if ev.DrawdownUSD > 0 {
		fmt.Fprintf(&b, "drawdown: %.2f / %.2f USDC\n", ev.DrawdownUSD, ev.DrawdownCap)
	}
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

// FormatSignalPrompt renders a compact DM body that accompanies the inline-keyboard.
// Only the signal side is surfaced — the boss picks amount, not direction.
//
// Layout (3 lines, keeps buttons visible on mobile without scrolling):
//
//	⚡ <signalOutcome> ↑ @ 0.xxxx
//	<Match>
//	<Context> · <EndIn> · Δ+x.xxpp buy 78%
func FormatSignalPrompt(ev SignalPromptEvent) string {
	var b strings.Builder

	sig, ok := signalChoice(ev.Choices)
	signalOutcome := "?"
	mid := 0.0
	if ok {
		signalOutcome = sig.Outcome
		mid = sig.Mid
	}

	if ev.WhaleLabel != "" {
		fmt.Fprintf(&b, "🐋 %s BUY %s @ %.4f\n", ev.WhaleLabel, truncateStr(signalOutcome, 40), mid)
	} else {
		fmt.Fprintf(&b, "⚡ %s ↑ @ %.4f\n", truncateStr(signalOutcome, 40), mid)
	}

	match := ev.Match
	if match == "" && len(ev.Choices) > 0 {
		match = ev.Choices[0].Outcome + " vs " + outcomeAt(ev.Choices, 1)
	}
	fmt.Fprintln(&b, truncateStr(match, 80))

	// Collapse context · endIn · Δ · buy into one line.
	parts := make([]string, 0, 4)
	if ev.Context != "" {
		parts = append(parts, truncateStr(ev.Context, 40))
	}
	if ev.EndIn != "" {
		parts = append(parts, ev.EndIn)
	}
	if ev.DeltaPP != 0 {
		parts = append(parts, fmt.Sprintf("Δ%+.2fpp", ev.DeltaPP))
	}
	if ev.TailLen > 0 || ev.BuyRatio > 0 {
		parts = append(parts, fmt.Sprintf("buy %.0f%%", ev.BuyRatio*100))
	}
	b.WriteString(strings.Join(parts, " · "))
	if ev.Slug != "" {
		fmt.Fprintf(&b, "\nhttps://newshare.bwb.online/zh/polymarket/event?slug=%s&_nobar=true&_needChain=matic", ev.Slug)
	}
	return b.String()
}

// signalChoice returns the first Choice flagged IsSignal, or the first entry
// if none are flagged (defensive; prompts are only emitted when a signal fires).
func signalChoice(cs []SignalChoice) (SignalChoice, bool) {
	for _, c := range cs {
		if c.IsSignal {
			return c, true
		}
	}
	if len(cs) > 0 {
		return cs[0], true
	}
	return SignalChoice{}, false
}

// DefaultSizesUSD is the Buy 10-100U inline-button gradient used when
// SignalPromptEvent.SizesUSD is empty.
var DefaultSizesUSD = []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

// FormatFillReceipt renders the durable DM body for a successful manual open.
// Kept minimal — the boss already saw the prompt + toast; this is the archive
// copy.
func FormatFillReceipt(ev FillReceiptEvent) string {
	tag := "🧾 成交凭据 · 手动"
	if ev.Source == "auto" {
		tag = "🧾 成交凭据 · 自动"
	}
	header := ev.Match
	if header == "" {
		header = ev.Question
	}
	if len(header) > 100 {
		header = header[:97] + "..."
	}
	return fmt.Sprintf(
		"%s\n%s · %s\n%gU @ %.4f · units %.2f\norder %s",
		tag, header, ev.Outcome, ev.SizeUSD, ev.FillPx, ev.Units, ev.OrderID,
	)
}

// FormatSignalExpired is the body used to rewrite an unclicked prompt once the
// TTL reaper evicts it.
func FormatSignalExpired() string { return "⌛ 已过期 · 未下单" }

// FormatSignalFilled rewrites the original prompt to reflect a successful fill
// after the boss clicks a size button.
func FormatSignalFilled(ev FillReceiptEvent) string {
	return fmt.Sprintf("✅ 已下单 · %s %gU @ %.4f\norder %s",
		ev.Outcome, ev.SizeUSD, ev.FillPx, ev.OrderID)
}

// FormatInjuryAlert renders a compact DM for a star-player injury report.
func FormatInjuryAlert(ev InjuryAlertEvent) string {
	icon := "🏥"
	if ev.Impact == "franchise_player_out" {
		icon = "🚨"
	}
	status := ev.Status
	if status == "" {
		status = "Out"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s · %s %s\n", icon, ev.Team, ev.StarPlayer, status)
	if ev.Reason != "" {
		fmt.Fprintf(&b, "原因: %s\n", ev.Reason)
	}
	fmt.Fprintf(&b, "影响: %s\n", ev.Impact)
	b.WriteString("对手 underdog 可能存在动量机会")
	return b.String()
}

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

// FormatClosePrompt renders the DM body for a whale-sell close decision.
func FormatClosePrompt(ev ClosePromptEvent) string {
	var b strings.Builder
	label := ev.WhaleLabel
	if label == "" {
		label = "鲸鱼"
	}
	fmt.Fprintf(&b, "🔻 %s 卖出 · 是否平仓?\n", label)
	mkt := ev.Market
	if len(mkt) > 80 {
		mkt = mkt[:77] + "..."
	}
	fmt.Fprintf(&b, "%s · %s\n", mkt, ev.Outcome)
	fmt.Fprintf(&b, "鲸鱼: SELL %.0f shares @ %.4f = $%.0f\n", ev.WhaleSize, ev.WhalePrice, ev.WhaleNotl)
	if ev.WhaleTotalShares > 0 {
		fmt.Fprintf(&b, "鲸鱼持仓: %.0f shares (均价 $%.4f)", ev.WhaleTotalShares, ev.WhaleAvgPrice)
		if ev.WhalePctSold > 0 {
			fmt.Fprintf(&b, " · 卖出 %.0f%%", ev.WhalePctSold)
			if ev.WhalePctSold >= 95 {
				b.WriteString(" 清仓")
			} else {
				b.WriteString(" 部分止盈")
			}
		}
		b.WriteByte('\n')
	}
	for _, p := range ev.Positions {
		fmt.Fprintf(&b, "我方持仓: %.2fU · entry %.4f · units %.2f\n", p.SizeUSD, p.EntryMid, p.Units)
	}
	if ev.LinkURL != "" {
		b.WriteString(ev.LinkURL)
		b.WriteByte('\n')
	}
	if ev.ProfileURL != "" {
		fmt.Fprintf(&b, "主页: %s", ev.ProfileURL)
	}
	return b.String()
}

// FormatWhaleAlert renders a compact DM for a smart-money large trade.
func FormatWhaleAlert(ev WhaleAlertEvent) string {
	icon := "🐋"
	if strings.ToUpper(ev.Side) == "SELL" {
		icon = "🔻"
	}
	addr := ev.Wallet
	if len(addr) > 10 {
		addr = addr[:6] + "…" + addr[len(addr)-4:]
	}
	mkt := ev.Market
	if len(mkt) > 80 {
		mkt = mkt[:77] + "..."
	}
	var b strings.Builder
	label := ev.Label
	if label == "" {
		label = addr
	}
	fmt.Fprintf(&b, "%s %s 大单\n", icon, label)
	fmt.Fprintf(&b, "%s · %s\n", mkt, ev.Outcome)
	fmt.Fprintf(&b, "%s %.0f shares @ %.4f = $%.0f\n", ev.Side, ev.SizeUnits, ev.Price, ev.Notional)
	if ev.TotalShares > 0 {
		fmt.Fprintf(&b, "持仓: %.0f shares (均价 $%.4f)", ev.TotalShares, ev.AvgPrice)
		if strings.ToUpper(ev.Side) == "SELL" && ev.PctSold > 0 {
			fmt.Fprintf(&b, " · 本次卖出 %.0f%%", ev.PctSold)
			if ev.PctSold >= 95 {
				b.WriteString(" 清仓")
			} else {
				b.WriteString(" 部分止盈")
			}
		}
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "钱包: %s\n", addr)
	if ev.LinkURL != "" {
		b.WriteString(ev.LinkURL)
		b.WriteByte('\n')
	}
	if ev.ProfileURL != "" {
		fmt.Fprintf(&b, "主页: %s", ev.ProfileURL)
	}
	return b.String()
}
