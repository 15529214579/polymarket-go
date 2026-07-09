package walletdiscover

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsMaxHistoricalOffsetError(t *testing.T) {
	err := errors.New(`GET https://data-api.polymarket.com/trades?limit=500&offset=3500: HTTP 400: {"error":"max historical activity offset of 3000 exceeded"}`)
	if !IsMaxHistoricalOffsetError(err) {
		t.Fatal("expected max historical offset error to be recognized")
	}
	if IsMaxHistoricalOffsetError(errors.New("HTTP 500")) {
		t.Fatal("unrelated errors should not be treated as offset limit")
	}
}

func TestLeaderboardPageUsesOffset(t *testing.T) {
	var gotLimit, gotOffset string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/profit" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		gotLimit = r.URL.Query().Get("limit")
		gotOffset = r.URL.Query().Get("offset")
		_ = json.NewEncoder(w).Encode([]LeaderboardEntry{{ProxyWallet: "0x0000000000000000000000000000000000000001"}})
	}))
	defer srv.Close()

	client := NewClient(Config{LeaderboardBase: srv.URL})
	entries, err := client.LeaderboardPage(context.Background(), "profit", "7d", 50, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if gotLimit != "50" || gotOffset != "100" {
		t.Fatalf("unexpected query limit=%q offset=%q", gotLimit, gotOffset)
	}
}
