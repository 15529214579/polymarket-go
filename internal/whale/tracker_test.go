package whale

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestSeedReplaysRecentLargeTrades(t *testing.T) {
	now := time.Now().Unix()
	wallet := WalletEntry{Address: "0x1111111111111111111111111111111111111111", Label: "test", MinSizeUSD: 500}
	var got []AlertEvent
	tracker := NewTracker(Config{
		Enabled:      true,
		Wallets:      []WalletEntry{wallet},
		MinSizeUSD:   500,
		PollInterval: time.Hour,
		ReplayWindow: 2 * time.Minute,
	}, func(ev AlertEvent) {
		got = append(got, ev)
	})
	tracker.http = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/activity":
			return jsonResponse(`[
				{"proxyWallet":"` + wallet.Address + `","side":"BUY","asset":"asset-old","conditionId":"0xold","size":2000,"price":0.5,"timestamp":` + formatInt(now-300) + `,"title":"Will Spain win on 2026-07-10?","outcome":"Yes","transactionHash":"0xold","type":"TRADE"},
				{"proxyWallet":"` + wallet.Address + `","side":"BUY","asset":"asset-new","conditionId":"0xnew","size":2000,"price":0.5,"timestamp":` + formatInt(now-30) + `,"title":"Will Spain win on 2026-07-10?","outcome":"Yes","transactionHash":"0xnew","type":"TRADE"},
				{"proxyWallet":"` + wallet.Address + `","side":"BUY","asset":"asset-small","conditionId":"0xsmall","size":10,"price":0.5,"timestamp":` + formatInt(now-20) + `,"title":"Will Spain win on 2026-07-10?","outcome":"Yes","transactionHash":"0xsmall","type":"TRADE"}
			]`), nil
		case "/positions":
			return jsonResponse(`[]`), nil
		default:
			t.Fatalf("unexpected path %s", req.URL.Path)
			return nil, nil
		}
	})}

	if err := tracker.seed(context.Background(), wallet); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("replayed alerts=%d, want 1: %#v", len(got), got)
	}
	if got[0].TradeID != "0xnew" || got[0].Notional != 1000 {
		t.Fatalf("unexpected replay alert: %#v", got[0])
	}
}

func TestSeedReplayDisabledByDefault(t *testing.T) {
	now := time.Now().Unix()
	wallet := WalletEntry{Address: "0x2222222222222222222222222222222222222222", Label: "test", MinSizeUSD: 500}
	var got []AlertEvent
	tracker := NewTracker(Config{
		Enabled:      true,
		Wallets:      []WalletEntry{wallet},
		MinSizeUSD:   500,
		PollInterval: time.Hour,
	}, func(ev AlertEvent) {
		got = append(got, ev)
	})
	tracker.http = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/activity" {
			return jsonResponse(`[{"proxyWallet":"` + wallet.Address + `","side":"BUY","asset":"asset-new","conditionId":"0xnew","size":2000,"price":0.5,"timestamp":` + formatInt(now-30) + `,"title":"Will Spain win on 2026-07-10?","outcome":"Yes","transactionHash":"0xnew","type":"TRADE"}]`), nil
		}
		return jsonResponse(`[]`), nil
	})}

	if err := tracker.seed(context.Background(), wallet); err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("replayed alerts=%d, want 0", len(got))
	}
}

func TestFetchTradesPaginates(t *testing.T) {
	tracker := NewTracker(Config{MaxPages: maxActivityPages}, nil)
	requests := 0
	tracker.http = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		offset := req.URL.Query().Get("offset")
		if offset == "0" {
			rows := make([]trade, 500)
			for i := range rows {
				rows[i] = trade{Asset: "asset-" + strconv.Itoa(i), Timestamp: 100}
			}
			raw, _ := json.Marshal(rows)
			return jsonResponse(string(raw)), nil
		}
		if offset == "500" {
			return jsonResponse(`[{"asset":"asset-500","timestamp":100}]`), nil
		}
		t.Fatalf("unexpected offset %s", offset)
		return nil, nil
	})}
	trades, err := tracker.fetchTrades(context.Background(), "0xwallet", 99, maxActivityPages)
	if err != nil {
		t.Fatal(err)
	}
	if len(trades) != 501 || requests != 2 {
		t.Fatalf("trades=%d requests=%d, want 501/2", len(trades), requests)
	}
}

func TestTradeKeyDistinguishesFillsInSameTransaction(t *testing.T) {
	left := trade{TransactionHash: "0xsame", Asset: "asset-a", Side: "BUY", Size: 1, Price: 0.5, Timestamp: 100}
	right := left
	right.Asset = "asset-b"
	if tradeKey(&left) == tradeKey(&right) {
		t.Fatal("tradeKey collapsed distinct fills")
	}
}

func TestTrackerWatermarkPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watermarks.json")
	wallet := WalletEntry{Address: "0x1111111111111111111111111111111111111111"}
	tracker := NewTracker(Config{Wallets: []WalletEntry{wallet}, StatePath: path}, nil)
	tracker.advanceFloor(wallet.Address, 123)
	if err := tracker.persistState(); err != nil {
		t.Fatal(err)
	}
	reloaded := NewTracker(Config{Wallets: []WalletEntry{wallet}, StatePath: path}, nil)
	state := reloaded.states[strings.ToLower(wallet.Address)]
	if state.lastTS != 123 {
		t.Fatalf("lastTS=%d, want 123", state.lastTS)
	}
}

func TestLoadWalletsFileWithListMinsUsesListMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wallets.txt")
	body := strings.Join([]string{
		"0x1111111111111111111111111111111111111111 # list=sports tier=A",
		"0x2222222222222222222222222222222222222222 # list=watch tier=B",
		"0x3333333333333333333333333333333333333333 # tier=C",
		"0x4444444444444444444444444444444444444444 list=core tier=A",
	}, "\n") + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	wallets, err := LoadWalletsFileWithListMins(path, 500, map[string]float64{
		"core":   1000,
		"sports": 1500,
		"watch":  750,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wallets) != 4 {
		t.Fatalf("wallets len=%d, want 4", len(wallets))
	}
	got := map[string]float64{}
	for _, w := range wallets {
		got[strings.ToLower(w.Address)] = w.MinSizeUSD
	}
	cases := map[string]float64{
		"0x1111111111111111111111111111111111111111": 1500,
		"0x2222222222222222222222222222222222222222": 750,
		"0x3333333333333333333333333333333333333333": 500,
		"0x4444444444444444444444444444444444444444": 1000,
	}
	for addr, want := range cases {
		if got[addr] != want {
			t.Fatalf("%s min=%v, want %v (all=%#v)", addr, got[addr], want, got)
		}
	}
}

func formatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}
