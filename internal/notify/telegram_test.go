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
		Nonce:   "abcd1234",
		Match:   "LoL: VIT vs GIANTX",
		Context: "Game 2 Winner",
		EndIn:   "1h 23m",
		Choices: []SignalChoice{
			{Slot: 0, Outcome: "VIT", Mid: 0.42, IsSignal: true},
			{Slot: 1, Outcome: "GIANTX", Mid: 0.58, IsSignal: false},
		},
		DeltaPP: 3.5, TailUps: 4, TailLen: 5, BuyRatio: 0.8,
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
	for _, want := range []string{"动量信号", "VIT", "Game 2 Winner", "结算 1h 23m"} {
		if !strings.Contains(text, want) {
			t.Errorf("text missing %q; got: %q", want, text)
		}
	}
	// Contrarian mid price must not be rendered in the body (match title may contain the team name legitimately).
	if strings.Contains(text, "0.5800") || strings.Contains(text, "选 GIANTX") {
		t.Errorf("contrarian side leaked into prompt body: %q", text)
	}
	rm, ok := m["reply_markup"].(map[string]any)
	if !ok {
		t.Fatalf("reply_markup missing: %v", m["reply_markup"])
	}
	kb, ok := rm["inline_keyboard"].([]any)
	if !ok || len(kb) != 1 {
		t.Fatalf("inline_keyboard shape: want 1 signal row, got %v", kb)
	}
	row0, _ := kb[0].([]any)
	if len(row0) != 3 {
		t.Fatalf("row0 buttons: got %d", len(row0))
	}
	b0, _ := row0[0].(map[string]any)
	if b0["text"] != "🟢 VIT 1U" || b0["callback_data"] != "buy:abcd1234:0:1" {
		t.Errorf("row0.button0: %+v", b0)
	}
	b2, _ := row0[2].(map[string]any)
	if b2["text"] != "🟢 VIT 10U" || b2["callback_data"] != "buy:abcd1234:0:10" {
		t.Errorf("row0.button2: %+v", b2)
	}
}

func TestButtonLabel_YesNoPolarity(t *testing.T) {
	cases := []struct {
		outcome  string
		isSignal bool
		size     float64
		want     string
	}{
		{"Yes", true, 1, "✅ Yes 1U"},
		{"yes", false, 5, "✅ Yes 5U"},
		{"No", false, 10, "❌ No 10U"},
		{"NO", true, 1, "❌ No 1U"},
		{"Gen.G", true, 5, "🟢 Gen.G 5U"},
		{"Nongshim", false, 1, "🔴 Nongshim 1U"},
	}
	for _, c := range cases {
		got := buttonLabel(c.outcome, c.size, c.isSignal)
		if got != c.want {
			t.Errorf("buttonLabel(%q, %v, %v) = %q, want %q", c.outcome, c.size, c.isSignal, got, c.want)
		}
	}
}

func TestTelegram_SignalPromptRoutesThroughPromptBot(t *testing.T) {
	// Verifies the fix for the Phase 3.5.b wiring bug: SignalPrompt must be sent
	// from PromptBotToken when set, so clicks reach the bot the LongPoll watches.
	var alertHits, promptHits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/botALERT/"):
			alertHits.Add(1)
		case strings.Contains(r.URL.Path, "/botPROMPT/"):
			promptHits.Add(1)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{
		BotToken: "ALERT", PromptBotToken: "PROMPT", ChatID: "1",
		BaseURL: srv.URL, QueueSize: 4,
	})
	tg.RiskTrip(RiskTripEvent{Reason: "daily_loss", DayPnLUSD: -14, DayLossCapUSD: 13.56})
	tg.LargeFill(LargeFillEvent{PnLUSD: 5})
	tg.SignalPrompt(SignalPromptEvent{
		Nonce: "n", Match: "m",
		Choices:   []SignalChoice{{Slot: 0, Outcome: "Yes", Mid: 0.5, IsSignal: true}},
		ExpiresIn: time.Minute,
	})
	if err := tg.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	if alertHits.Load() != 2 {
		t.Errorf("alert bot sends: got %d, want 2", alertHits.Load())
	}
	if promptHits.Load() != 1 {
		t.Errorf("prompt bot sends: got %d, want 1", promptHits.Load())
	}
}

func TestFormatRiskTrip_FeedSilence(t *testing.T) {
	s := FormatRiskTrip(RiskTripEvent{Reason: "feed_silence", SilentSec: 47})
	if !strings.Contains(s, "WSS") || !strings.Contains(s, "47s") {
		t.Errorf("bad format: %q", s)
	}
}

type capturedHit struct {
	path string
	body map[string]any
}

func newCaptureServer(hits chan<- capturedHit) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		hits <- capturedHit{path: r.URL.Path, body: m}
		w.WriteHeader(200)
		switch {
		case strings.Contains(r.URL.Path, "/sendMessage"):
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":4242}}`))
		case strings.Contains(r.URL.Path, "/editMessageText"):
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":4242}}`))
		default:
			_, _ = w.Write([]byte(`{"ok":true}`))
		}
	}))
}

func TestTelegram_SignalPrompt_OnSent_ReportsMessageID(t *testing.T) {
	hits := make(chan capturedHit, 4)
	srv := newCaptureServer(hits)
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{
		BotToken: "ALERT", PromptBotToken: "PROMPT", ChatID: "1",
		BaseURL: srv.URL, QueueSize: 4,
	})

	gotID := make(chan int64, 1)
	tg.SignalPrompt(SignalPromptEvent{
		Nonce:     "n",
		Choices:   []SignalChoice{{Slot: 0, Outcome: "Yes", Mid: 0.5, IsSignal: true}},
		ExpiresIn: time.Minute,
		OnSent: func(id int64, err error) {
			if err != nil {
				t.Errorf("OnSent err: %v", err)
			}
			gotID <- id
		},
	})
	if err := tg.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	select {
	case id := <-gotID:
		if id != 4242 {
			t.Errorf("OnSent msgID = %d, want 4242", id)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnSent never fired")
	}
	// Sanity: prompt path was used.
	h := <-hits
	if !strings.Contains(h.path, "/botPROMPT/sendMessage") {
		t.Errorf("prompt sent on wrong path: %s", h.path)
	}
}

func TestTelegram_EditSignalExpired_StripsKeyboard(t *testing.T) {
	hits := make(chan capturedHit, 2)
	srv := newCaptureServer(hits)
	defer srv.Close()

	tg := NewTelegram(TelegramConfig{
		BotToken: "ALERT", PromptBotToken: "PROMPT", ChatID: "9",
		BaseURL: srv.URL, QueueSize: 4,
	})
	tg.EditSignalExpired(77)
	if err := tg.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	h := <-hits
	if !strings.Contains(h.path, "/botPROMPT/editMessageText") {
		t.Errorf("expected editMessageText via PROMPT, got %s", h.path)
	}
	if h.body["message_id"].(float64) != 77 {
		t.Errorf("message_id = %v, want 77", h.body["message_id"])
	}
	text, _ := h.body["text"].(string)
	if !strings.Contains(text, "已过期") {
		t.Errorf("text missing 已过期: %q", text)
	}
	rm, ok := h.body["reply_markup"].(map[string]any)
	if !ok {
		t.Fatal("reply_markup missing")
	}
	kb, _ := rm["inline_keyboard"].([]any)
	if len(kb) != 0 {
		t.Errorf("inline_keyboard should be empty, got %d rows", len(kb))
	}
}

func TestTelegram_EditSignalExpired_ZeroIsNoOp(t *testing.T) {
	hits := make(chan capturedHit, 1)
	srv := newCaptureServer(hits)
	defer srv.Close()
	tg := NewTelegram(TelegramConfig{BotToken: "A", ChatID: "1", BaseURL: srv.URL, QueueSize: 2})
	tg.EditSignalExpired(0)
	_ = tg.Close(context.Background())
	select {
	case h := <-hits:
		t.Fatalf("unexpected hit: %s", h.path)
	default:
	}
}

func TestTelegram_EditSignalFilled_Renders(t *testing.T) {
	hits := make(chan capturedHit, 2)
	srv := newCaptureServer(hits)
	defer srv.Close()
	tg := NewTelegram(TelegramConfig{
		BotToken: "ALERT", PromptBotToken: "PROMPT", ChatID: "1",
		BaseURL: srv.URL, QueueSize: 4,
	})
	tg.EditSignalFilled(FillReceiptEvent{
		Outcome: "Gen.G", SizeUSD: 5, FillPx: 0.4321, OrderID: "paper-abc",
	}, 123)
	_ = tg.Close(context.Background())
	h := <-hits
	if !strings.Contains(h.path, "/botPROMPT/editMessageText") {
		t.Errorf("path: %s", h.path)
	}
	text, _ := h.body["text"].(string)
	if !strings.Contains(text, "已下单") || !strings.Contains(text, "Gen.G") || !strings.Contains(text, "0.4321") {
		t.Errorf("text missing bits: %q", text)
	}
}

func TestTelegram_FillReceipt_GoesToAlertBot(t *testing.T) {
	hits := make(chan capturedHit, 2)
	srv := newCaptureServer(hits)
	defer srv.Close()
	tg := NewTelegram(TelegramConfig{
		BotToken: "ALERT", PromptBotToken: "PROMPT", ChatID: "1",
		BaseURL: srv.URL, QueueSize: 4,
	})
	tg.FillReceipt(FillReceiptEvent{
		Question: "LoL: Gen.G vs Nongshim - Game 1 Winner",
		Match:    "LoL: Gen.G vs Nongshim",
		Outcome:  "Gen.G", SizeUSD: 1, Units: 2, FillPx: 0.5, OrderID: "paper-xyz",
		Source: "manual",
	})
	_ = tg.Close(context.Background())
	h := <-hits
	if !strings.Contains(h.path, "/botALERT/sendMessage") {
		t.Errorf("receipt should go via alert bot, got path %s", h.path)
	}
	text, _ := h.body["text"].(string)
	if !strings.Contains(text, "成交凭据") || !strings.Contains(text, "paper-xyz") {
		t.Errorf("text missing bits: %q", text)
	}
}
