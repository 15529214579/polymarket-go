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
//
// MessageID is filled asynchronously after Telegram confirms the send — so
// TTL-expiry / click-success can edit the original prompt to "已过期" / "已下单".
// A zero MessageID means the Telegram response hasn't landed yet (extremely
// rare for normal flows; edits are skipped silently if still zero).
type PendingIntent struct {
	Nonce     string
	Market    string
	Question  string
	Choices   []Choice
	MessageID int64
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
// Returns (_, false) if missing or expired. Prefer Peek for paper stacking.
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

// Peek looks up a nonce without deleting it, so the same prompt can be
// clicked repeatedly inside its TTL. Used by paper-mode stacking. Returns
// (_, false) if missing or expired; an expired entry is left to the reaper.
func (s *PendingStore) Peek(nonce string, now time.Time) (PendingIntent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.m[nonce]
	if !ok {
		return PendingIntent{}, false
	}
	if now.After(p.ExpiresAt) {
		return PendingIntent{}, false
	}
	return p, true
}

// Reap drops expired entries and returns them so the caller can edit the
// original prompt DM to "已过期" (or whatever Phase 3.5 TTL UX is). Safe to
// call from a background ticker.
func (s *PendingStore) Reap(now time.Time) []PendingIntent {
	s.mu.Lock()
	defer s.mu.Unlock()
	var evicted []PendingIntent
	for k, v := range s.m {
		if now.After(v.ExpiresAt) {
			delete(s.m, k)
			evicted = append(evicted, v)
		}
	}
	return evicted
}

// SetMessageID records the Telegram message_id that a prompt was sent as.
// Called from the Telegram drain goroutine once the send succeeds. Returns
// false if the nonce has already been claimed or reaped.
func (s *PendingStore) SetMessageID(nonce string, id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.m[nonce]
	if !ok {
		return false
	}
	p.MessageID = id
	s.m[nonce] = p
	return true
}

// Size returns the current count — useful for risk_status logs.
func (s *PendingStore) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.m)
}

// CloseIntent stores context for a whale-sell close prompt so the callback
// handler can look it up by nonce and execute the close.
type CloseIntent struct {
	Nonce     string
	AssetID   string
	Market    string
	Question  string
	Outcome   string
	WhalePrice float64
	MessageID int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

// CloseStore is a thread-safe TTL map for close prompts, keyed by nonce.
type CloseStore struct {
	ttl time.Duration
	mu  sync.Mutex
	m   map[string]CloseIntent
}

func NewCloseStore(ttl time.Duration) *CloseStore {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &CloseStore{ttl: ttl, m: make(map[string]CloseIntent)}
}

func (s *CloseStore) Put(in CloseIntent, now time.Time) CloseIntent {
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

func (s *CloseStore) Claim(nonce string, now time.Time) (CloseIntent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.m[nonce]
	if !ok {
		return CloseIntent{}, false
	}
	delete(s.m, nonce)
	if now.After(p.ExpiresAt) {
		return CloseIntent{}, false
	}
	return p, true
}

func (s *CloseStore) SetMessageID(nonce string, id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.m[nonce]
	if !ok {
		return false
	}
	p.MessageID = id
	s.m[nonce] = p
	return true
}

func (s *CloseStore) Reap(now time.Time) []CloseIntent {
	s.mu.Lock()
	defer s.mu.Unlock()
	var evicted []CloseIntent
	for k, v := range s.m {
		if now.After(v.ExpiresAt) {
			delete(s.m, k)
			evicted = append(evicted, v)
		}
	}
	return evicted
}

// newNonce returns 8 bytes of hex (16 chars). Telegram caps callback_data at
// 64 bytes; we fit well under that with "buy:<nonce>:<slot>:<size>" ≈ 28 chars.
func newNonce() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
