package order

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAmountsForIntentPreservesTickPrice(t *testing.T) {
	maker, taker, err := amountsForIntent(Intent{
		AssetID: "asset", Side: Buy, SizeUSD: 2, LimitPx: 0.936,
	})
	if err != nil {
		t.Fatal(err)
	}
	if maker.String() != "1999998" || taker.String() != "2136750" {
		t.Fatalf("buy amounts: maker=%s taker=%s", maker, taker)
	}
	price, _ := new(big.Rat).SetFrac(maker, taker).Float64()
	if diff := price - 0.936; diff > 1e-12 || diff < -1e-12 {
		t.Fatalf("implied buy price=%v", price)
	}

	maker, taker, err = amountsForIntent(Intent{
		AssetID: "asset", Side: Sell, SizeUSD: 9.36, SizeShares: 10, LimitPx: 0.936,
	})
	if err != nil {
		t.Fatal(err)
	}
	if maker.String() != "10000000" || taker.String() != "9360000" {
		t.Fatalf("sell amounts: maker=%s taker=%s", maker, taker)
	}
}

func TestFillFromAmountsUsesActualMatchedAmounts(t *testing.T) {
	fallbackMaker := big.NewInt(20_000_000)
	fallbackTaker := big.NewInt(40_000_000)

	units, price := fillFromAmounts(Buy, "19500000", "39000000", fallbackMaker, fallbackTaker)
	if units != 39 || price != 0.5 {
		t.Fatalf("buy fill units=%v price=%v", units, price)
	}

	units, price = fillFromAmounts(Sell, "10000000", "4900000", fallbackTaker, fallbackMaker)
	if units != 10 || price < 0.49-1e-12 || price > 0.49+1e-12 {
		t.Fatalf("sell fill units=%v price=%v", units, price)
	}
}

func TestParseFixedAmount(t *testing.T) {
	if got := parseFixedAmount("100000000"); got != 100 {
		t.Fatalf("fixed amount=%v", got)
	}
	if got := parseFixedAmount("12.5"); got != 12.5 {
		t.Fatalf("decimal amount=%v", got)
	}
}

func TestGetOpenOrdersUsesV2DataEndpointAndSanitizesOwner(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.RequestURI() != "/data/orders" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		if r.Header.Get("POLY_API_KEY") == "" || r.Header.Get("POLY_SIGNATURE") == "" {
			t.Fatal("missing L2 authentication headers")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"order-1","owner":"secret-api-key","status":"LIVE","asset_id":"asset-1"}],"count":1}`))
	})

	result, err := client.GetOpenOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 1 || len(result.Data) != 1 || result.Data[0].ID != "order-1" {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestGetBalanceAllowanceUsesCollateralEndpoint(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.RequestURI() != "/balance-allowance?asset_type=COLLATERAL&signature_type=0" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"balance":"1181600","allowance":"999"}`))
	})

	result, err := client.GetBalanceAllowance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Balance != "1181600" {
		t.Fatalf("balance = %q", result.Balance)
	}
}

func TestRefreshBalanceAllowanceUsesV2UpdateEndpoint(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.RequestURI() != "/balance-allowance/update?asset_type=COLLATERAL&signature_type=0" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.RefreshBalanceAllowance(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrderUsesV2DataEndpoint(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.RequestURI() != "/data/order/order-1" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		_, _ = w.Write([]byte(`{"id":"order-1","status":"ORDER_STATUS_LIVE","size_matched":"1000000","price":"0.5"}`))
	})
	result, err := client.GetOrder(context.Background(), "order-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "order-1" || result.SizeMatched != 1 || result.AvgPrice != 0.5 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestCancelOrderUsesV2BodyEndpoint(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.RequestURI() != "/order" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		var body struct {
			OrderID string `json:"orderID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.OrderID != "order-1" {
			t.Fatalf("orderID = %q", body.OrderID)
		}
		_, _ = w.Write([]byte(`{"canceled":["order-1"],"not_canceled":{}}`))
	})
	if err := client.CancelOrder(context.Background(), "order-1"); err != nil {
		t.Fatal(err)
	}
}

func TestSubmitUsesCurrentV2Payload(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.RequestURI() != "/order" {
			t.Fatalf("request = %s %s", r.Method, r.URL.RequestURI())
		}
		var payload struct {
			Order struct {
				TokenID   string `json:"tokenId"`
				Timestamp string `json:"timestamp"`
				Metadata  string `json:"metadata"`
				Builder   string `json:"builder"`
			} `json:"order"`
			Owner     string `json:"owner"`
			OrderType string `json:"orderType"`
			DeferExec bool   `json:"deferExec"`
			PostOnly  bool   `json:"postOnly"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Order.TokenID != "123" || payload.Order.Timestamp == "" {
			t.Fatalf("unexpected order: %+v", payload.Order)
		}
		if len(payload.Order.Metadata) != 66 || len(payload.Order.Builder) != 66 {
			t.Fatalf("invalid bytes32 fields: metadata=%q builder=%q", payload.Order.Metadata, payload.Order.Builder)
		}
		if payload.Owner != "api-key" || payload.OrderType != "GTC" || payload.DeferExec || payload.PostOnly {
			t.Fatalf("unexpected envelope: %+v", payload)
		}
		_, _ = w.Write([]byte(`{"success":true,"orderID":"order-1","status":"matched","makingAmount":"5000000","takingAmount":"10000000"}`))
	})
	result, err := client.Submit(context.Background(), Intent{
		AssetID: "123",
		Side:    Buy,
		SizeUSD: 5,
		LimitPx: 0.5,
		Type:    GTC,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusFilled || result.FilledSize != 10 || result.AvgPrice != 0.5 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSubmitRejectsNonDecimalAssetID(t *testing.T) {
	client := testV2Client(t, func(http.ResponseWriter, *http.Request) {
		t.Fatal("HTTP must not be called for an invalid asset ID")
	})
	result, err := client.Submit(context.Background(), Intent{
		AssetID: "not-a-token",
		Side:    Buy,
		SizeUSD: 5,
		LimitPx: 0.5,
	})
	if err == nil || result.Status != StatusRejected {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestAuthenticatedReadRejectsHTTPError(t *testing.T) {
	client := testV2Client(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
	if _, err := client.GetOpenOrders(context.Background()); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 error, got %v", err)
	}
}

func TestGetCLOBVersion(t *testing.T) {
	hc := testHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/version" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"version":2}`))
	}))

	version, err := getCLOBVersion(context.Background(), "https://clob.test", hc)
	if err != nil {
		t.Fatal(err)
	}
	if version != 2 {
		t.Fatalf("version = %d", version)
	}
}

func testV2Client(t *testing.T, handler http.HandlerFunc) *V2Client {
	t.Helper()
	wallet, err := NewWalletFromHexKey(strings.Repeat("1", 64))
	if err != nil {
		t.Fatal(err)
	}
	client := NewV2Client(wallet, &APICredentials{
		APIKey:     "api-key",
		Secret:     "c2VjcmV0",
		Passphrase: "passphrase",
	}, false)
	client.clobBase = "https://clob.test"
	client.http = testHTTPClient(handler)
	return client
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testHTTPClient(handler http.Handler) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		return recorder.Result(), nil
	})}
}
