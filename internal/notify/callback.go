package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// CallbackHandler is the inbound half of the Phase 3.5 click-to-buy flow. It
// is called once per valid "buy:<nonce>:<slot>:<sizeUSD>:<mode>" callback_query.
// Slot indexes into PendingIntent.Choices (YES/NO etc). Mode is "ladder" or
// "hold" (hold-to-settlement, no SL/timeout). The ack string (if any) is
// surfaced as the Telegram toast on the clicker's screen.
type CallbackHandler interface {
	OnBuy(ctx context.Context, nonce string, slot int, sizeUSD float64, mode string, messageID int64) (ack string, err error)
}

// LongPollConfig bounds a single long-poll consumer. Use a DEDICATED bot token
// here — never reuse a token another process is also polling (updates are
// consumed competitively on the Bot API).
type LongPollConfig struct {
	BotToken string
	// ExpectedChatID is the Telegram chat ID we trust to click buttons.
	// Anything else is ignored with a toast.
	ExpectedChatID int64
	BaseURL        string        // default "https://api.telegram.org"
	PollTimeout    time.Duration // getUpdates long-poll seconds, default 25
	HTTPTimeout    time.Duration // per-request HTTP timeout, default 35s (must > PollTimeout)
	BackoffOnErr   time.Duration // sleep after transport failure, default 3s
}

// LongPoll drives a blocking getUpdates loop on a dedicated bot token.
// Safe for concurrent use only via Run (single-flight).
type LongPoll struct {
	cfg     LongPollConfig
	client  *http.Client
	handler CallbackHandler
	logger  *slog.Logger
}

// NewLongPoll returns a poller ready to Run. Defaults applied.
func NewLongPoll(cfg LongPollConfig, h CallbackHandler) *LongPoll {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.telegram.org"
	}
	if cfg.PollTimeout <= 0 {
		cfg.PollTimeout = 25 * time.Second
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = cfg.PollTimeout + 10*time.Second
	}
	if cfg.BackoffOnErr <= 0 {
		cfg.BackoffOnErr = 3 * time.Second
	}
	return &LongPoll{
		cfg:     cfg,
		client:  &http.Client{Timeout: cfg.HTTPTimeout},
		handler: h,
		logger:  slog.Default(),
	}
}

// Run blocks until ctx is cancelled. Transport errors are logged + retried
// after BackoffOnErr. Individual malformed callbacks toast the user but never
// stop the loop.
func (l *LongPoll) Run(ctx context.Context) error {
	var offset int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		updates, next, err := l.getUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			l.logger.Warn("longpoll_get_err", "err", err.Error())
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(l.cfg.BackoffOnErr):
			}
			continue
		}
		for _, u := range updates {
			l.dispatch(ctx, u)
		}
		offset = next
	}
}

type tgUpdate struct {
	UpdateID      int64            `json:"update_id"`
	CallbackQuery *tgCallbackQuery `json:"callback_query,omitempty"`
}

type tgCallbackQuery struct {
	ID   string `json:"id"`
	From struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"from"`
	Message *struct {
		MessageID int64 `json:"message_id"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message,omitempty"`
	Data string `json:"data"`
}

func (l *LongPoll) getUpdates(ctx context.Context, offset int64) ([]tgUpdate, int64, error) {
	v := url.Values{}
	v.Set("timeout", strconv.Itoa(int(l.cfg.PollTimeout.Seconds())))
	if offset > 0 {
		v.Set("offset", strconv.FormatInt(offset, 10))
	}
	v.Set("allowed_updates", `["callback_query"]`)
	u := fmt.Sprintf("%s/bot%s/getUpdates?%s", l.cfg.BaseURL, l.cfg.BotToken, v.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, offset, err
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, offset, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, offset, fmt.Errorf("getUpdates status=%d body=%s", resp.StatusCode, string(body))
	}
	var out struct {
		OK     bool       `json:"ok"`
		Result []tgUpdate `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, offset, err
	}
	if !out.OK {
		return nil, offset, fmt.Errorf("getUpdates returned ok=false")
	}
	next := offset
	for _, up := range out.Result {
		if up.UpdateID >= next {
			next = up.UpdateID + 1
		}
	}
	return out.Result, next, nil
}

// dispatch is best-effort: every branch must call answerCallbackQuery so the
// Telegram spinner goes away on the boss's screen.
func (l *LongPoll) dispatch(ctx context.Context, u tgUpdate) {
	if u.CallbackQuery == nil {
		return
	}
	cq := u.CallbackQuery

	// Only honor clicks from the expected chat. We also trust From.ID because
	// for private chats chat.ID == from.ID, and for groups we'd only ever DM
	// the boss anyway.
	if cq.Message == nil {
		l.answerCallback(ctx, cq.ID, "ignored", true)
		return
	}
	if cq.Message.Chat.ID != l.cfg.ExpectedChatID {
		l.logger.Warn("longpoll_reject_chat",
			"chat_id", cq.Message.Chat.ID,
			"from_id", cq.From.ID,
			"expected", l.cfg.ExpectedChatID,
		)
		l.answerCallback(ctx, cq.ID, "not authorized", true)
		return
	}

	nonce, slot, size, mode, ok := parseBuyCallback(cq.Data)
	if !ok {
		l.answerCallback(ctx, cq.ID, "bad data", true)
		return
	}

	l.logger.Info("callback_click",
		"nonce", nonce,
		"slot", slot,
		"size_usd", size,
		"mode", mode,
		"from", cq.From.Username,
	)

	ack, err := l.handler.OnBuy(ctx, nonce, slot, size, mode, cq.Message.MessageID)
	if err != nil {
		l.answerCallback(ctx, cq.ID, "❌ "+truncate(err.Error(), 180), true)
		return
	}
	if ack == "" {
		ack = "✅ 已下单"
	}
	l.answerCallback(ctx, cq.ID, ack, false)
}

// parseBuyCallback validates and splits "buy:<nonce>:<slot>:<sizeUSD>:<mode>".
// Mode is "l" (ladder) or "h" (hold-to-settlement); 4-part legacy format
// defaults to "ladder". Returns the expanded mode string ("ladder"/"hold").
func parseBuyCallback(s string) (nonce string, slot int, sizeUSD float64, mode string, ok bool) {
	parts := strings.SplitN(s, ":", 5)
	if len(parts) < 4 || parts[0] != "buy" || parts[1] == "" {
		return "", 0, 0, "", false
	}
	sl, err := strconv.Atoi(parts[2])
	if err != nil || sl < 0 {
		return "", 0, 0, "", false
	}
	sz, err := strconv.ParseFloat(parts[3], 64)
	if err != nil || sz <= 0 {
		return "", 0, 0, "", false
	}
	m := "ladder"
	if len(parts) == 5 {
		switch parts[4] {
		case "h":
			m = "hold"
		case "l":
			m = "ladder"
		default:
			return "", 0, 0, "", false
		}
	}
	return parts[1], sl, sz, m, true
}

func (l *LongPoll) answerCallback(ctx context.Context, cqID string, text string, showAlert bool) {
	if cqID == "" {
		return
	}
	body := map[string]any{
		"callback_query_id": cqID,
		"text":              truncate(text, 200), // Telegram caps at 200 chars
		"show_alert":        showAlert,
	}
	payload, _ := json.Marshal(body)
	u := fmt.Sprintf("%s/bot%s/answerCallbackQuery", l.cfg.BaseURL, l.cfg.BotToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(string(payload)))
	if err != nil {
		l.logger.Warn("answer_cb_build_err", "err", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := l.client.Do(req)
	if err != nil {
		l.logger.Warn("answer_cb_http_err", "err", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		l.logger.Warn("answer_cb_status", "status", resp.StatusCode, "body", string(b))
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
