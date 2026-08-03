package walletdiscover

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func testHTTPResponse(req *http.Request, status int, body string, headers http.Header) *http.Response {
	if headers == nil {
		headers = make(http.Header)
	}
	return &http.Response{
		StatusCode: status,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func TestIsMaxHistoricalOffsetError(t *testing.T) {
	err := errors.New(`GET https://data-api.polymarket.com/trades?limit=500&offset=3500: HTTP 400: {"error":"max historical activity offset of 3000 exceeded"}`)
	if !IsMaxHistoricalOffsetError(err) {
		t.Fatal("expected max historical offset error to be recognized")
	}
	if IsMaxHistoricalOffsetError(errors.New("HTTP 500")) {
		t.Fatal("unrelated errors should not be treated as offset limit")
	}
}

func TestClientRetriesRateLimitThenSucceeds(t *testing.T) {
	var attempts atomic.Int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if attempts.Add(1) < 3 {
			return testHTTPResponse(r, http.StatusTooManyRequests, "rate limited", http.Header{"Retry-After": []string{"0"}}), nil
		}
		return testHTTPResponse(r, http.StatusOK, `[{"proxyWallet":"0x0000000000000000000000000000000000000001"}]`, nil), nil
	})

	client := NewClient(Config{
		LeaderboardBase: "https://leaderboard.test",
		HTTPMaxAttempts: 3,
		HTTPRetryBase:   time.Millisecond,
		HTTPRetryMax:    2 * time.Millisecond,
	})
	client.http.Transport = transport
	rows, err := client.Leaderboard(context.Background(), "profit", "7d", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || attempts.Load() != 3 {
		t.Fatalf("rows=%d attempts=%d", len(rows), attempts.Load())
	}
	stats := client.Stats()
	if stats.Retries != 2 || stats.RateLimits != 2 || stats.Failures != 0 {
		t.Fatalf("stats=%+v", stats)
	}
}

func TestLeaderboardPageUsesOffset(t *testing.T) {
	var gotLimit, gotOffset string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/profit" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		gotLimit = r.URL.Query().Get("limit")
		gotOffset = r.URL.Query().Get("offset")
		return testHTTPResponse(r, http.StatusOK, `[{"proxyWallet":"0x0000000000000000000000000000000000000001"}]`, nil), nil
	})

	client := NewClient(Config{LeaderboardBase: "https://leaderboard.test"})
	client.http.Transport = transport
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
