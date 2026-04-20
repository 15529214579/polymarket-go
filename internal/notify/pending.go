package notify

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Choice is one selectable outcome in a PendingIntent. Binary markets carry 2
// (YES / NO or team A / team B); multi-outcome markets could carry more. Slot
// is the index encoded in callback_data so the handler can resolve back
// without trusting the clicker with an asset_id.
type Choice struct {
	AssetID  string
	Outcome  string  // outcome label as Polymarket returns it
	Mid      float64 // latest mid at prompt time
	IsSignal bool    // true for the asset that fired the momentum signal
}

// PendingIntent is the snapshot captured when a signal prompt is sent. The
// callback consumer looks it up by nonce, picks Choices[slot], and replays it
// through the order path at the chosen size. Stored in memory only; process
// restarts wipe pending prompts (desired — 10min TTL anyway).
type PendingIntent struct {
	Nonce     string
	Market    string
	Question  string
	Choices   []Choice
	CreatedAt time.Time
	ExpiresAt time.Time
}

// PendingStore is a thread-safe TTL map keyed by nonce. Callers either Claim
// a nonce (one-shot, removes it) or let TTL evict it. A single reaper
// goroutine walks the map every few seconds.
type PendingStore struct {
	ttl time.Duration
	mu  sync.Mutex
	m   map[string]PendingIntent
}

// NewPendingStore returns a store with the given TTL and launches a reaper
// tied to the returned stop func. Call stop() on shutdown.
func NewPendingStore(ttl time.Duration) *PendingStore {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &PendingStore{ttl: ttl, m: make(map[string]PendingIntent)}
}

// Put stores an intent under a fresh nonce and returns it. Caller embeds the
// nonce in Telegram callback_data so the callback handler can Claim it back.
func (s *PendingStore) Put(in PendingIntent, now time.Time) PendingIntent {
	if in.Nonce == "" {
		in.Nonce = newNonce()
	}
	in.CreatedAt = now
	in.ExpiresAt = now.Add(s.ttl)
	s.mu.Lock()
	s.m[in.Nonce] = in
	s.mu.Unlock()
	return in
}

// Claim looks up and removes a nonce in one shot (prevents button replay).
// Returns (_, false) if missing or expired.
func (s *PendingStore) Claim(nonce string, now time.Time) (PendingIntent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.m[nonce]
	if !ok {
		return PendingIntent{}, false
	}
	delete(s.m, nonce)
	if now.After(p.ExpiresAt) {
		return PendingIntent{}, false
	}
	return p, true
}

// Reap drops expired entries. Safe to call from a background ticker.
func (s *PendingStore) Reap(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for k, v := range s.m {
		if now.After(v.ExpiresAt) {
			delete(s.m, k)
			n++
		}
	}
	return n
}

// Size returns the current count — useful for risk_status logs.
func (s *PendingStore) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.m)
}

// newNonce returns 8 bytes of hex (16 chars). Telegram caps callback_data at
// 64 bytes; we fit well under that with "buy:<nonce>:<slot>:<size>" ≈ 28 chars.
func newNonce() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
