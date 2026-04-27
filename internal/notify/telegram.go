package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TelegramConfig is the minimum to push to one chat.
type TelegramConfig struct {
	BotToken string
	ChatID   string
	// PromptBotToken, when set, is used *only* for SignalPrompt messages so the
	// inline-keyboard message originates from the same bot the callback long-poll
	// watches. Leave empty to send prompts via BotToken (only correct if that bot
	// is also the one being long-polled for callback_query).
	PromptBotToken string
	// PushBotToken, when set, is used for informational push notifications
	// (injury alerts, text alerts, whale alerts, risk events, large fills)
	// so they don't clutter the order bot's chat. Falls back to BotToken.
	PushBotToken string
	// BaseURL defaults to Telegram's production Bot API. Override for tests.
	BaseURL string
	// SendTimeout bounds each HTTP attempt. Default 5s.
	SendTimeout time.Duration
	// QueueSize caps how many pending messages may buffer before we start
	// dropping (prevents a stuck Telegram pipe from blocking trading). Default 32.
	QueueSize int
}

// Telegram implements Notifier against the Bot API. All event calls enqueue
// a formatted message into a buffered channel; a single background goroutine
// drains it and POSTs. Close flushes and stops the goroutine.
type Telegram struct {
	cfg    TelegramConfig
	client *http.Client
	logger *slog.Logger

	queue chan outgoing
	wg    sync.WaitGroup
	stop  chan struct{}
	once  sync.Once
}

// NewTelegram wires a ready notifier. Caller must Close() before exit to
// drain in-flight sends (up to queue size).
func NewTelegram(cfg TelegramConfig) *Telegram {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.telegram.org"
	}
	if cfg.SendTimeout <= 0 {
		cfg.SendTimeout = 5 * time.Second
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 32
	}
	t := &Telegram{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.SendTimeout},
		logger: slog.Default(),
		queue:  make(chan outgoing, cfg.QueueSize),
		stop:   make(chan struct{}),
	}
	t.wg.Add(1)
	go t.drain()
	return t
}

func (t *Telegram) pushToken() string {
	if t.cfg.PushBotToken != "" {
		return t.cfg.PushBotToken
	}
	return t.cfg.BotToken
}

func (t *Telegram) RiskTrip(ev RiskTripEvent) {
	t.enqueue(outgoing{text: FormatRiskTrip(ev), tag: "risk_trip", sendToken: t.pushToken()})
}
func (t *Telegram) RiskResume(ev RiskResumeEvent) {
	t.enqueue(outgoing{text: FormatRiskResume(ev), tag: "risk_resume", sendToken: t.pushToken()})
}
func (t *Telegram) LargeFill(ev LargeFillEvent) {
	t.enqueue(outgoing{text: FormatLargeFill(ev), tag: "large_fill", sendToken: t.pushToken()})
}

// SignalPrompt enqueues a DM with inline-keyboard rows for the signal side:
//   - 🟢 rows: 10U–100U ladder mode (SL + 4h timeout), 5 buttons per row
//   - 🔒 rows: 10U–100U hold to settlement (no SL, no timeout), 5 per row
//
// callback_data is "buy:<nonce>:<slot>:<sizeUSD>:<mode>" where mode is "l"
// (ladder) or "h" (hold). The inbound callback handler resolves nonce via
// the PendingStore and routes to the appropriate exit strategy.
func (t *Telegram) SignalPrompt(ev SignalPromptEvent) {
	sizes := ev.SizesUSD
	if len(sizes) == 0 {
		sizes = DefaultSizesUSD
	}
	sig, ok := signalChoice(ev.Choices)
	if !ok {
		sig = SignalChoice{Slot: 0, Outcome: "?", IsSignal: true}
	}

	const rowCap = 5
	var rows [][]map[string]string
	// Ladder rows
	var cur []map[string]string
	for _, s := range sizes {
		cur = append(cur, map[string]string{
			"text":          buttonLabel(sig.Outcome, s, true),
			"callback_data": fmt.Sprintf("buy:%s:%d:%g:l", ev.Nonce, sig.Slot, s),
		})
		if len(cur) == rowCap {
			rows = append(rows, cur)
			cur = nil
		}
	}
	if len(cur) > 0 {
		rows = append(rows, cur)
	}
	// Hold rows
	cur = nil
	for _, s := range sizes {
		cur = append(cur, map[string]string{
			"text":          holdButtonLabel(s),
			"callback_data": fmt.Sprintf("buy:%s:%d:%g:h", ev.Nonce, sig.Slot, s),
		})
		if len(cur) == rowCap {
			rows = append(rows, cur)
			cur = nil
		}
	}
	if len(cur) > 0 {
		rows = append(rows, cur)
	}
	kb := map[string]any{"inline_keyboard": rows}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{
		text: FormatSignalPrompt(ev), tag: "signal_prompt",
		replyMarkup: kb, sendToken: tok, onSent: ev.OnSent,
	})
}

// EditSignalExpired rewrites the original prompt to "已过期" + strips buttons.
// Called from the pending reaper. Uses the prompt bot (the one the boss sees
// the DM from); falls back to the alert bot when PromptBotToken is unset.
func (t *Telegram) EditSignalExpired(messageID int64) {
	if messageID == 0 {
		return
	}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{
		tag: "edit_expired", text: FormatSignalExpired(),
		editMessageID: messageID, stripKeyboard: true, sendToken: tok,
	})
}

// EditSignalFilled rewrites the original prompt to "已下单 …" and strips the
// keyboard. Called after a successful callback click (Phase 3.5 C).
func (t *Telegram) EditSignalFilled(ev FillReceiptEvent, messageID int64) {
	if messageID == 0 {
		return
	}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{
		tag: "edit_filled", text: FormatSignalFilled(ev),
		editMessageID: messageID, stripKeyboard: true, sendToken: tok,
	})
}

// FillReceipt enqueues a durable archive DM for a manual open. Goes via the
// alert bot (not the prompt bot) so the boss's prompt-bot chat stays lean.
func (t *Telegram) FillReceipt(ev FillReceiptEvent) {
	t.enqueue(outgoing{text: FormatFillReceipt(ev), tag: "fill_receipt"})
}

// InjuryAlert enqueues a star-player injury DM. Guarded by -injury_enabled flag
// at the call site; to remove: delete this method + InjuryAlertEvent + FormatInjuryAlert.
func (t *Telegram) InjuryAlert(ev InjuryAlertEvent) {
	t.enqueue(outgoing{text: FormatInjuryAlert(ev), tag: "injury_alert", sendToken: t.pushToken()})
}

func (t *Telegram) TextAlert(text string) {
	t.enqueue(outgoing{text: text, tag: "text_alert", sendToken: t.pushToken()})
}

// ClosePrompt enqueues a DM with a close button when a whale sells an asset
// we hold. The boss clicks "✅ 平仓" to confirm or ignores.
func (t *Telegram) ClosePrompt(ev ClosePromptEvent) {
	closeRow := []map[string]string{
		{"text": "✅ 平仓", "callback_data": fmt.Sprintf("close:%s", ev.Nonce)},
		{"text": "❌ 忽略", "callback_data": fmt.Sprintf("skip:%s", ev.Nonce)},
	}
	kb := map[string]any{"inline_keyboard": [][]map[string]string{closeRow}}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{
		text: FormatClosePrompt(ev), tag: "close_prompt",
		replyMarkup: kb, sendToken: tok, onSent: ev.OnSent,
	})
}

// EditCloseDone rewrites a close prompt to show the result and strips buttons.
func (t *Telegram) EditCloseDone(text string, messageID int64) {
	if messageID == 0 {
		return
	}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{
		tag: "edit_close_done", text: text,
		editMessageID: messageID, stripKeyboard: true, sendToken: tok,
	})
}

// WhaleAlert enqueues a smart-money whale trade DM via the push bot.
func (t *Telegram) WhaleAlert(ev WhaleAlertEvent) {
	t.enqueue(outgoing{text: FormatWhaleAlert(ev), tag: "whale_alert", sendToken: t.pushToken()})
}

// buttonLabel builds a short inline-button caption for ladder-mode buttons.
func buttonLabel(outcome string, sizeUSD float64, isSignal bool) string {
	switch strings.ToLower(strings.TrimSpace(outcome)) {
	case "yes":
		return fmt.Sprintf("✅ Yes %gU", sizeUSD)
	case "no":
		return fmt.Sprintf("❌ No %gU", sizeUSD)
	}
	if isSignal {
		return fmt.Sprintf("🟢 %gU", sizeUSD)
	}
	return fmt.Sprintf("🔴 %gU", sizeUSD)
}

// holdButtonLabel builds a caption for the hold-to-settlement row.
func holdButtonLabel(sizeUSD float64) string {
	return fmt.Sprintf("🔒 %gU", sizeUSD)
}

type outgoing struct {
	text        string
	tag         string
	replyMarkup any
	// sendToken, when non-empty, overrides cfg.BotToken for this message. Used
	// for SignalPrompt so the message originates from the sidecar bot that the
	// LongPoll watches for callback_query.
	sendToken string
	// editMessageID, when non-zero, switches the dispatch from sendMessage to
	// editMessageText on that existing message. Used to rewrite a prompt to
	// "已过期" / "已下单" post-hoc.
	editMessageID int64
	// stripKeyboard, when true on an edit, replaces reply_markup with an empty
	// inline_keyboard so the buttons disappear.
	stripKeyboard bool
	// onSent, if set, is called after the HTTP send completes with the returned
	// message_id (0 on error). Used by SignalPrompt to hand the message id to
	// the PendingStore for later edits.
	onSent func(messageID int64, err error)
}

func (t *Telegram) enqueue(o outgoing) {
	select {
	case t.queue <- o:
	default:
		t.logger.Warn("notify_drop", "reason", "queue_full", "tag", o.tag)
	}
}

// Close waits for the drain goroutine to finish sending queued messages, or
// for ctx to cancel.
func (t *Telegram) Close(ctx context.Context) error {
	t.once.Do(func() { close(t.stop) })
	done := make(chan struct{})
	go func() { t.wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Telegram) drain() {
	defer t.wg.Done()
	for {
		select {
		case <-t.stop:
			// Flush whatever's already queued, then exit.
			for {
				select {
				case o := <-t.queue:
					t.send(o)
				default:
					return
				}
			}
		case o := <-t.queue:
			t.send(o)
		}
	}
}

// send is synchronous inside the drain goroutine. Errors are logged but not
// returned — the trading loop never waits on Telegram.
func (t *Telegram) send(o outgoing) {
	tok := o.sendToken
	if tok == "" {
		tok = t.cfg.BotToken
	}

	endpoint := "sendMessage"
	body := map[string]any{
		"chat_id":                  t.cfg.ChatID,
		"text":                     o.text,
		"disable_web_page_preview": true,
	}
	if o.editMessageID != 0 {
		endpoint = "editMessageText"
		body["message_id"] = o.editMessageID
		if o.stripKeyboard {
			body["reply_markup"] = map[string]any{"inline_keyboard": [][]map[string]string{}}
		}
	} else if o.replyMarkup != nil {
		body["reply_markup"] = o.replyMarkup
	}

	url := fmt.Sprintf("%s/bot%s/%s", t.cfg.BaseURL, tok, endpoint)
	payload, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.SendTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		t.logger.Warn("notify_send_build_err", "err", err.Error())
		if o.onSent != nil {
			o.onSent(0, err)
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.Warn("notify_send_http_err", "err", err.Error())
		if o.onSent != nil {
			o.onSent(0, err)
		}
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		t.logger.Warn("notify_send_status",
			"status", resp.StatusCode,
			"tag", o.tag,
			"body", string(raw),
		)
		if o.onSent != nil {
			o.onSent(0, fmt.Errorf("telegram status=%d", resp.StatusCode))
		}
		return
	}
	if o.onSent != nil {
		msgID := parseMessageID(raw)
		o.onSent(msgID, nil)
	}
}

// parseMessageID pulls result.message_id out of a Bot API success envelope.
// Returns 0 on any parse failure.
func parseMessageID(body []byte) int64 {
	var env struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return 0
	}
	if !env.OK {
		return 0
	}
	return env.Result.MessageID
}
