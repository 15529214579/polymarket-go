package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestBuyTimeStorePersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "buy-times.json")
	store, err := loadBuyTimeStore(path)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 8, 3, 1, 2, 3, 0, time.UTC)
	if err := store.Set("asset-1", want); err != nil {
		t.Fatal(err)
	}
	reloaded, err := loadBuyTimeStore(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := reloaded.Get("asset-1")
	if !ok || !got.Equal(want) {
		t.Fatalf("Get()=(%v,%v), want (%v,true)", got, ok, want)
	}
}
