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

// RiskTrip formats and enqueues a risk-breaker trip. Drops silently if the
// queue is saturated (a stuck Telegram pipe must not block the trading loop).
func (t *Telegram) RiskTrip(ev RiskTripEvent)     { t.enqueue(outgoing{text: FormatRiskTrip(ev), tag: "risk_trip"}) }
func (t *Telegram) RiskResume(ev RiskResumeEvent) { t.enqueue(outgoing{text: FormatRiskResume(ev), tag: "risk_resume"}) }
func (t *Telegram) LargeFill(ev LargeFillEvent)   { t.enqueue(outgoing{text: FormatLargeFill(ev), tag: "large_fill"}) }

// SignalPrompt enqueues a DM with a single inline-keyboard row for the signal
// side only (boss picks amount, not direction). Buttons are "Buy 1U / 5U / 10U";
// callback_data is "buy:<nonce>:<slot>:<sizeUSD>" where slot indexes into
// PendingIntent.Choices. The inbound callback handler resolves nonce via the
// PendingStore and executes Choices[slot].
func (t *Telegram) SignalPrompt(ev SignalPromptEvent) {
	sizes := ev.SizesUSD
	if len(sizes) == 0 {
		sizes = DefaultSizesUSD
	}
	sig, ok := signalChoice(ev.Choices)
	if !ok {
		// Defensive: caller forgot to populate Choices.
		sig = SignalChoice{Slot: 0, Outcome: "?", IsSignal: true}
	}
	row := make([]map[string]string, 0, len(sizes))
	for _, s := range sizes {
		row = append(row, map[string]string{
			"text":          buttonLabel(sig.Outcome, s, true),
			"callback_data": fmt.Sprintf("buy:%s:%d:%g", ev.Nonce, sig.Slot, s),
		})
	}
	kb := map[string]any{"inline_keyboard": [][]map[string]string{row}}
	tok := t.cfg.PromptBotToken
	if tok == "" {
		tok = t.cfg.BotToken
	}
	t.enqueue(outgoing{text: FormatSignalPrompt(ev), tag: "signal_prompt", replyMarkup: kb, sendToken: tok})
}

// buttonLabel trims the outcome to fit Telegram's ~40-char inline-button cap.
// Yes/No markets get ✅/❌ for explicit polarity; team-vs-team keeps ⚡ on the
// signal (upward-momentum) row so the buy side is obvious at a glance.
func buttonLabel(outcome string, sizeUSD float64, isSignal bool) string {
	name := outcome
	if name == "" {
		name = "?"
	}
	switch strings.ToLower(name) {
	case "yes":
		return fmt.Sprintf("✅ Yes %gU", sizeUSD)
	case "no":
		return fmt.Sprintf("❌ No %gU", sizeUSD)
	}
	if len(name) > 18 {
		name = name[:15] + "..."
	}
	if isSignal {
		return fmt.Sprintf("🟢 %s %gU", name, sizeUSD)
	}
	return fmt.Sprintf("🔴 %s %gU", name, sizeUSD)
}

type outgoing struct {
	text        string
	tag         string
	replyMarkup any
	// sendToken, when non-empty, overrides cfg.BotToken for this message. Used
	// for SignalPrompt so the message originates from the sidecar bot that the
	// LongPoll watches for callback_query.
	sendToken string
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
	url := fmt.Sprintf("%s/bot%s/sendMessage", t.cfg.BaseURL, tok)
	body := map[string]any{
		"chat_id":                  t.cfg.ChatID,
		"text":                     o.text,
		"disable_web_page_preview": true,
	}
	if o.replyMarkup != nil {
		body["reply_markup"] = o.replyMarkup
	}
	payload, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.SendTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		t.logger.Warn("notify_send_build_err", "err", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		t.logger.Warn("notify_send_http_err", "err", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		t.logger.Warn("notify_send_status",
			"status", resp.StatusCode,
			"body", string(body),
		)
	}
}
