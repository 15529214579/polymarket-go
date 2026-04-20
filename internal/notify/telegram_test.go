package notify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestTelegram_SendsFormattedRiskTrip(t *testing.T) {
	var got atomic.Value // json payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/bot-TOK/sendMessage") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		got.Store(m)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{
		BotToken: "-TOK", ChatID: "42", BaseURL: srv.URL,
		SendTimeout: 2 * time.Second, QueueSize: 4,
	})
	tg.RiskTrip(RiskTripEvent{
		Reason: "daily_loss", DayPnLUSD: -14.23, DayLossCapUSD: 13.56, OpenPositions: 2,
	})
	if err := tg.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	m, ok := got.Load().(map[string]any)
	if !ok {
		t.Fatalf("no payload captured")
	}
	if m["chat_id"] != "42" {
		t.Errorf("chat_id: %v", m["chat_id"])
	}
	text, _ := m["text"].(string)
	if !strings.Contains(text, "日亏损上限") || !strings.Contains(text, "-14.23") {
		t.Errorf("text missing expected bits: %q", text)
	}
}

func TestTelegram_DropsWhenQueueFull(t *testing.T) {
	// Server that blocks forever so the drain goroutine can't make progress.
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
	}))
	defer srv.Close()
	defer close(block)

	tg := NewTelegram(TelegramConfig{
		BotToken: "x", ChatID: "1", BaseURL: srv.URL,
		SendTimeout: 100 * time.Millisecond, QueueSize: 1,
	})

	// First send enters the queue and is picked up (now blocks at HTTP).
	// Second fills the 1-slot queue. Third+ should drop without blocking.
	for i := 0; i < 10; i++ {
		tg.LargeFill(LargeFillEvent{AssetID: "a", PnLUSD: 1})
	}
	// If we got here without deadlocking, drop path worked.
}

func TestTelegram_DrainsPendingOnClose(t *testing.T) {
	var count atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{BotToken: "x", ChatID: "1", BaseURL: srv.URL, QueueSize: 8})
	for i := 0; i < 3; i++ {
		tg.LargeFill(LargeFillEvent{PnLUSD: float64(i) + 1})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := tg.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}
	if got := count.Load(); got != 3 {
		t.Errorf("expected 3 sends, got %d", got)
	}
}

func TestFormatLargeFill(t *testing.T) {
	ev := LargeFillEvent{
		Question: "LEC VIT vs GIANTX Game 2 Winner",
		AssetID:  "xyz", EntryPx: 0.1754, ExitPx: 0.0912,
		PnLUSD: -3.57, Reason: "stop_loss", HeldSec: 428, SizeUSD: 5,
	}
	s := FormatLargeFill(ev)
	if !strings.Contains(s, "📉") || !strings.Contains(s, "-3.57") || !strings.Contains(s, "stop_loss") {
		t.Errorf("bad format: %q", s)
	}
}

func TestTelegram_SignalPromptAttachesInlineKeyboard(t *testing.T) {
	var got atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		got.Store(m)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{BotToken: "t", ChatID: "1", BaseURL: srv.URL, QueueSize: 4})
	tg.SignalPrompt(SignalPromptEvent{
		Nonce: "abcd1234", Question: "LEC VIT vs GIANTX Game 2",
		AssetID: "a", Mid: 0.42, DeltaPP: 3.5, TailUps: 4, TailLen: 5, BuyRatio: 0.8,
		ExpiresIn: 60 * time.Second,
	})
	if err := tg.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	m, _ := got.Load().(map[string]any)
	if m == nil {
		t.Fatal("no payload captured")
	}
	text, _ := m["text"].(string)
	if !strings.Contains(text, "动量信号") || !strings.Contains(text, "VIT") {
		t.Errorf("text: %q", text)
	}
	rm, ok := m["reply_markup"].(map[string]any)
	if !ok {
		t.Fatalf("reply_markup missing: %v", m["reply_markup"])
	}
	kb, ok := rm["inline_keyboard"].([]any)
	if !ok || len(kb) != 1 {
		t.Fatalf("inline_keyboard shape: %v", kb)
	}
	row, _ := kb[0].([]any)
	if len(row) != 3 {
		t.Fatalf("expected 3 buttons, got %d", len(row))
	}
	b0, _ := row[0].(map[string]any)
	if b0["text"] != "Buy 1U" || b0["callback_data"] != "buy:abcd1234:1" {
		t.Errorf("button[0]: %+v", b0)
	}
	b2, _ := row[2].(map[string]any)
	if b2["text"] != "Buy 10U" || b2["callback_data"] != "buy:abcd1234:10" {
		t.Errorf("button[2]: %+v", b2)
	}
}

func TestFormatRiskTrip_FeedSilence(t *testing.T) {
	s := FormatRiskTrip(RiskTripEvent{Reason: "feed_silence", SilentSec: 47})
	if !strings.Contains(s, "WSS") || !strings.Contains(s, "47s") {
		t.Errorf("bad format: %q", s)
	}
}
