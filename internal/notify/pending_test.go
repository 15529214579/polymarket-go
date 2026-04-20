package notify

import (
	"testing"
	"time"
)

func TestPendingStore_PutClaim_OneShot(t *testing.T) {
	s := NewPendingStore(60 * time.Second)
	now := time.Now()
	p := s.Put(PendingIntent{
		Market:   "m",
		Question: "q",
		Choices: []Choice{
			{AssetID: "a", Outcome: "Yes", Mid: 0.42, IsSignal: true},
			{AssetID: "b", Outcome: "No", Mid: 0.58},
		},
	}, now)
	if p.Nonce == "" {
		t.Fatal("nonce should be auto-generated")
	}
	if p.ExpiresAt.Sub(p.CreatedAt) != 60*time.Second {
		t.Errorf("TTL not applied: %v", p.ExpiresAt.Sub(p.CreatedAt))
	}
	got, ok := s.Claim(p.Nonce, now.Add(1*time.Second))
	if !ok || len(got.Choices) != 2 || got.Choices[0].AssetID != "a" {
		t.Errorf("Claim(inside TTL): ok=%v got=%+v", ok, got)
	}
	if _, ok := s.Claim(p.Nonce, now.Add(2*time.Second)); ok {
		t.Error("double Claim should fail")
	}
}

func TestPendingStore_Claim_ExpiredRemoves(t *testing.T) {
	s := NewPendingStore(30 * time.Second)
	now := time.Now()
	p := s.Put(PendingIntent{Choices: []Choice{{AssetID: "x"}}}, now)
	if _, ok := s.Claim(p.Nonce, now.Add(61*time.Second)); ok {
		t.Error("Claim after TTL should fail")
	}
	if _, ok := s.Claim(p.Nonce, now); ok {
		t.Error("entry should have been deleted on expired Claim")
	}
}

func TestPendingStore_Reap(t *testing.T) {
	s := NewPendingStore(10 * time.Second)
	now := time.Now()
	s.Put(PendingIntent{Choices: []Choice{{AssetID: "1"}}}, now.Add(-20*time.Second))
	s.Put(PendingIntent{Choices: []Choice{{AssetID: "2"}}}, now.Add(-5*time.Second))
	s.Put(PendingIntent{Choices: []Choice{{AssetID: "3"}}}, now)
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
