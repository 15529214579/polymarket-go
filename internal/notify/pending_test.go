package notify

import (
	"testing"
	"time"
)

func TestPendingStore_PutClaim_OneShot(t *testing.T) {
	s := NewPendingStore(60 * time.Second)
	now := time.Now()
	p := s.Put(PendingIntent{AssetID: "a", Market: "m", Mid: 0.42}, now)
	if p.Nonce == "" {
		t.Fatal("nonce should be auto-generated")
	}
	if p.ExpiresAt.Sub(p.CreatedAt) != 60*time.Second {
		t.Errorf("TTL not applied: %v", p.ExpiresAt.Sub(p.CreatedAt))
	}
	got, ok := s.Claim(p.Nonce, now.Add(1*time.Second))
	if !ok || got.AssetID != "a" {
		t.Errorf("Claim(inside TTL): ok=%v got=%+v", ok, got)
	}
	// second Claim returns false (one-shot)
	if _, ok := s.Claim(p.Nonce, now.Add(2*time.Second)); ok {
		t.Error("double Claim should fail")
	}
}

func TestPendingStore_Claim_ExpiredRemoves(t *testing.T) {
	s := NewPendingStore(30 * time.Second)
	now := time.Now()
	p := s.Put(PendingIntent{AssetID: "x"}, now)
	if _, ok := s.Claim(p.Nonce, now.Add(61*time.Second)); ok {
		t.Error("Claim after TTL should fail")
	}
	// even a re-try at t=0 should fail — expired claims remove the entry
	if _, ok := s.Claim(p.Nonce, now); ok {
		t.Error("entry should have been deleted on expired Claim")
	}
}

func TestPendingStore_Reap(t *testing.T) {
	s := NewPendingStore(10 * time.Second)
	now := time.Now()
	s.Put(PendingIntent{AssetID: "1"}, now.Add(-20*time.Second))
	s.Put(PendingIntent{AssetID: "2"}, now.Add(-5*time.Second))
	s.Put(PendingIntent{AssetID: "3"}, now)
	if got := s.Size(); got != 3 {
		t.Fatalf("size=%d", got)
	}
	n := s.Reap(now)
	if n != 1 {
		t.Errorf("reap removed %d, want 1 (only the -20s entry is past its 10s TTL)", n)
	}
	if got := s.Size(); got != 2 {
		t.Errorf("size after reap=%d, want 2", got)
	}
}
