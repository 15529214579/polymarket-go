package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// TelegramConfig is the minimum to push to one chat.
type TelegramConfig struct {
	BotToken string
	ChatID   string
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

	queue chan string
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
		queue:  make(chan string, cfg.QueueSize),
		stop:   make(chan struct{}),
	}
	t.wg.Add(1)
	go t.drain()
	return t
}

// RiskTrip formats and enqueues a risk-breaker trip. Drops silently if the
// queue is saturated (a stuck Telegram pipe must not block the trading loop).
func (t *Telegram) RiskTrip(ev RiskTripEvent)     { t.enqueue(FormatRiskTrip(ev), "risk_trip") }
func (t *Telegram) RiskResume(ev RiskResumeEvent) { t.enqueue(FormatRiskResume(ev), "risk_resume") }
func (t *Telegram) LargeFill(ev LargeFillEvent)   { t.enqueue(FormatLargeFill(ev), "large_fill") }

func (t *Telegram) enqueue(msg, tag string) {
	select {
	case t.queue <- msg:
	default:
		t.logger.Warn("notify_drop", "reason", "queue_full", "tag", tag)
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
				case msg := <-t.queue:
					t.send(msg)
				default:
					return
				}
			}
		case msg := <-t.queue:
			t.send(msg)
		}
	}
}

// send is synchronous inside the drain goroutine. Errors are logged but not
// returned — the trading loop never waits on Telegram.
func (t *Telegram) send(msg string) {
	url := fmt.Sprintf("%s/bot%s/sendMessage", t.cfg.BaseURL, t.cfg.BotToken)
	payload, _ := json.Marshal(map[string]any{
		"chat_id":                  t.cfg.ChatID,
		"text":                     msg,
		"disable_web_page_preview": true,
	})
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
