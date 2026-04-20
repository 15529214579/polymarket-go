package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type stubHandler struct {
	calls  atomic.Int64
	last   struct{ nonce string; size float64 }
	ack    string
	err    error
	mu     sync.Mutex
}

func (s *stubHandler) OnBuy(ctx context.Context, nonce string, size float64) (string, error) {
	s.calls.Add(1)
	s.mu.Lock()
	s.last.nonce = nonce
	s.last.size = size
	s.mu.Unlock()
	return s.ack, s.err
}

func TestParseBuyCallback(t *testing.T) {
	cases := []struct {
		in      string
		nonce   string
		size    float64
		ok      bool
	}{
		{"buy:abc123:5", "abc123", 5, true},
		{"buy:abc123:0.5", "abc123", 0.5, true},
		{"buy:abc123:10", "abc123", 10, true},
		{"sell:abc:5", "", 0, false},
		{"buy::5", "", 0, false},
		{"buy:abc:-1", "", 0, false},
		{"buy:abc:bad", "", 0, false},
		{"garbage", "", 0, false},
		{"buy:abc", "", 0, false},
	}
	for _, c := range cases {
		n, s, ok := parseBuyCallback(c.in)
		if n != c.nonce || s != c.size || ok != c.ok {
			t.Errorf("parseBuyCallback(%q) = (%q,%v,%v), want (%q,%v,%v)",
				c.in, n, s, ok, c.nonce, c.size, c.ok)
		}
	}
}

// fakeTelegram is a minimal bot-API stub: serves one canned batch of updates
// then empties, and logs every answerCallbackQuery. Run exits via ctx cancel.
type fakeTelegram struct {
	mu       sync.Mutex
	batches  [][]tgUpdate
	served   int
	answered []answerCall
}

type answerCall struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

func (f *fakeTelegram) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/bot", http.NotFound) // just to anchor
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if contains(path, "/getUpdates") {
			f.mu.Lock()
			var batch []tgUpdate
			if f.served < len(f.batches) {
				batch = f.batches[f.served]
				f.served++
			}
			f.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": batch})
			return
		}
		if contains(path, "/answerCallbackQuery") {
			var ac answerCall
			_ = json.NewDecoder(r.Body).Decode(&ac)
			f.mu.Lock()
			f.answered = append(f.answered, ac)
			f.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			return
		}
		http.NotFound(w, r)
	})
	return mux
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle string) int {
	// minimal avoid-strings-dep helper
	if needle == "" {
		return 0
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

func newPoller(t *testing.T, srv *httptest.Server, h CallbackHandler, chatID int64) *LongPoll {
	t.Helper()
	return NewLongPoll(LongPollConfig{
		BotToken:       "test",
		ExpectedChatID: chatID,
		BaseURL:        srv.URL,
		PollTimeout:    1 * time.Second,
		HTTPTimeout:    2 * time.Second,
		BackoffOnErr:   50 * time.Millisecond,
	}, h)
}

func makeCallback(updateID, chatID int64, data string) tgUpdate {
	up := tgUpdate{UpdateID: updateID, CallbackQuery: &tgCallbackQuery{
		ID:   fmt.Sprintf("cq-%d", updateID),
		Data: data,
	}}
	up.CallbackQuery.From.ID = chatID
	up.CallbackQuery.From.Username = "boss"
	up.CallbackQuery.Message = &struct {
		MessageID int64 `json:"message_id"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	}{MessageID: updateID}
	up.CallbackQuery.Message.Chat.ID = chatID
	return up
}

func TestLongPoll_DispatchesValidBuy(t *testing.T) {
	ft := &fakeTelegram{batches: [][]tgUpdate{
		{makeCallback(1, 42, "buy:nonce1:5")},
	}}
	srv := httptest.NewServer(ft.handler())
	defer srv.Close()

	h := &stubHandler{ack: "已下 5U"}
	lp := newPoller(t, srv, h, 42)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	_ = lp.Run(ctx)

	if h.calls.Load() != 1 {
		t.Fatalf("handler called %d times, want 1", h.calls.Load())
	}
	if h.last.nonce != "nonce1" || h.last.size != 5 {
		t.Fatalf("got nonce=%q size=%v", h.last.nonce, h.last.size)
	}
	ft.mu.Lock()
	defer ft.mu.Unlock()
	if len(ft.answered) != 1 || ft.answered[0].Text != "已下 5U" {
		t.Fatalf("answerCallbackQuery wrong: %+v", ft.answered)
	}
}

func TestLongPoll_RejectsForeignChat(t *testing.T) {
	ft := &fakeTelegram{batches: [][]tgUpdate{
		{makeCallback(1, 999, "buy:nonce1:5")}, // foreign chat
	}}
	srv := httptest.NewServer(ft.handler())
	defer srv.Close()

	h := &stubHandler{}
	lp := newPoller(t, srv, h, 42) // expect 42

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	_ = lp.Run(ctx)

	if h.calls.Load() != 0 {
		t.Fatalf("handler should not be called for foreign chat, got %d", h.calls.Load())
	}
	ft.mu.Lock()
	defer ft.mu.Unlock()
	if len(ft.answered) != 1 || ft.answered[0].Text != "not authorized" {
		t.Fatalf("want not-authorized toast, got %+v", ft.answered)
	}
}

func TestLongPoll_RejectsMalformedData(t *testing.T) {
	ft := &fakeTelegram{batches: [][]tgUpdate{
		{makeCallback(1, 42, "buy::5")},     // empty nonce
		{makeCallback(2, 42, "buy:abc:-5")}, // negative size
		{makeCallback(3, 42, "hack:abc:5")}, // wrong action
	}}
	srv := httptest.NewServer(ft.handler())
	defer srv.Close()

	h := &stubHandler{}
	lp := newPoller(t, srv, h, 42)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = lp.Run(ctx)

	if h.calls.Load() != 0 {
		t.Fatalf("handler should not be called for malformed, got %d", h.calls.Load())
	}
	ft.mu.Lock()
	defer ft.mu.Unlock()
	if len(ft.answered) != 3 {
		t.Fatalf("want 3 bad-data answers, got %d", len(ft.answered))
	}
	for _, a := range ft.answered {
		if a.Text != "bad data" {
			t.Errorf("unexpected toast: %q", a.Text)
		}
	}
}

func TestLongPoll_SurfaceHandlerError(t *testing.T) {
	ft := &fakeTelegram{batches: [][]tgUpdate{
		{makeCallback(1, 42, "buy:nonce1:5")},
	}}
	srv := httptest.NewServer(ft.handler())
	defer srv.Close()

	h := &stubHandler{err: fmt.Errorf("已过期")}
	lp := newPoller(t, srv, h, 42)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	_ = lp.Run(ctx)

	ft.mu.Lock()
	defer ft.mu.Unlock()
	if len(ft.answered) != 1 {
		t.Fatalf("want 1 answer, got %d", len(ft.answered))
	}
	if ft.answered[0].Text != "❌ 已过期" || !ft.answered[0].ShowAlert {
		t.Fatalf("want alert with error prefix, got %+v", ft.answered[0])
	}
}
