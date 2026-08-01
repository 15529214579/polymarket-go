package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	ClobBaseURL       = "https://clob.polymarket.com"
	usdcScale         = 1_000_000 // 6 decimals
	pollInterval      = 2 * time.Second
	pollTimeout       = 30 * time.Second
	maxCancelAttempts = 2
)

type V2Client struct {
	wallet   *Wallet
	creds    *APICredentials
	clobBase string
	http     *http.Client
	negRisk  bool
	feeMu    sync.Mutex
	feeRates map[string]float64
}

type OpenOrder struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Market       string `json:"market"`
	AssetID      string `json:"asset_id"`
	Side         string `json:"side"`
	OriginalSize string `json:"original_size"`
	SizeMatched  string `json:"size_matched"`
	Price        string `json:"price"`
}

type OpenOrdersResponse struct {
	Data       []OpenOrder `json:"data"`
	NextCursor string      `json:"next_cursor"`
	Limit      int         `json:"limit"`
	Count      int         `json:"count"`
}

type BalanceAllowanceResponse struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}

func NewV2Client(wallet *Wallet, creds *APICredentials, negRisk bool) *V2Client {
	return &V2Client{
		wallet:   wallet,
		creds:    creds,
		clobBase: ClobBaseURL,
		http:     &http.Client{Timeout: 30 * time.Second},
		negRisk:  negRisk,
		feeRates: make(map[string]float64),
	}
}

func (c *V2Client) Name() string { return "v2-live" }

func (c *V2Client) Submit(ctx context.Context, in Intent) (Result, error) {
	if err := validate(in); err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	exchange := common.HexToAddress(V2ExchangeAddress)
	if in.NegRisk {
		exchange = common.HexToAddress(V2NegRiskExchangeAddress)
	}

	now := time.Now()
	makerAmt, takerAmt, err := amountsForIntent(in)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}
	side := SigBuy
	if in.Side == Sell {
		side = SigSell
	}

	tokenID, ok := new(big.Int).SetString(in.AssetID, 10)
	if !ok || tokenID.Sign() < 0 {
		err := fmt.Errorf("order: invalid decimal asset ID %q", in.AssetID)
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	order := V2Order{
		Salt:          NewSalt(),
		Maker:         c.wallet.Address(),
		Signer:        c.wallet.Address(),
		TokenID:       tokenID,
		MakerAmount:   makerAmt,
		TakerAmount:   takerAmt,
		Side:          side,
		SignatureType: SigTypeEOA,
		Timestamp:     big.NewInt(now.UnixMilli()),
	}

	digest, err := EIP712HashV2Order(order, exchange)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, fmt.Errorf("order: eip712 hash: %w", err)
	}
	sig, err := c.wallet.SignDigest(digest)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, fmt.Errorf("order: sign: %w", err)
	}

	orderType := string(in.Type)
	if orderType == "" {
		orderType = string(GTC)
	}

	payload := sendOrderPayload{
		Order: orderJSON{
			Maker:         c.wallet.Address().Hex(),
			Signer:        c.wallet.Address().Hex(),
			TokenID:       in.AssetID,
			MakerAmount:   makerAmt.String(),
			TakerAmount:   takerAmt.String(),
			Side:          string(in.Side),
			Expiration:    "0",
			Timestamp:     fmt.Sprintf("%d", now.UnixMilli()),
			Metadata:      "0x" + fmt.Sprintf("%064x", 0),
			Builder:       "0x" + fmt.Sprintf("%064x", 0),
			Signature:     fmt.Sprintf("0x%x", sig),
			Salt:          order.Salt.Int64(),
			SignatureType: int(SigTypeEOA),
		},
		Owner:     c.creds.APIKey,
		OrderType: orderType,
		DeferExec: false,
		PostOnly:  false,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	headers := buildL2Headers(c.creds, c.wallet.Address(), "POST", "/order", string(bodyBytes))

	req, err := http.NewRequestWithContext(ctx, "POST", c.clobBase+"/order", nil)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}
	req.Body = io.NopCloser(bytesReader(bodyBytes))
	req.ContentLength = int64(len(bodyBytes))
	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, fmt.Errorf("order: POST /order: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("CLOB %d: %s", resp.StatusCode, respBody)
		return Result{Status: StatusRejected, Error: errMsg}, fmt.Errorf("order: %s", errMsg)
	}

	var clobResp clobOrderResponse
	if err := json.Unmarshal(respBody, &clobResp); err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, fmt.Errorf("order: parse response: %w", err)
	}

	slog.Info("v2_post_response",
		"order_id", clobResp.OrderID,
		"status", clobResp.Status,
		"success", clobResp.Success,
		"error_msg", clobResp.ErrorMsg,
		"trades", len(clobResp.TradeIDs))

	if !clobResp.Success && clobResp.ErrorMsg != "" {
		errMsg := fmt.Sprintf("CLOB rejected: %s", clobResp.ErrorMsg)
		return Result{Status: StatusRejected, Error: errMsg}, fmt.Errorf("order: %s", errMsg)
	}

	if strings.EqualFold(clobResp.Status, "matched") {
		filledSize, avgPrice := fillFromAmounts(in.Side, clobResp.MakingAmount, clobResp.TakingAmount, makerAmt, takerAmt)
		feeUSD := c.fillFeeUSD(ctx, in.Market, filledSize, avgPrice)
		slog.Info("v2_order_filled",
			"order_id", clobResp.OrderID,
			"trades", len(clobResp.TradeIDs),
			"filled_size", filledSize,
			"avg_price", avgPrice,
			"fee_usd", feeUSD)
		return Result{
			OrderID:    clobResp.OrderID,
			Status:     StatusFilled,
			FilledSize: filledSize,
			AvgPrice:   avgPrice,
			SubmitAt:   now,
			FilledAt:   time.Now(),
			FeeUSD:     feeUSD,
		}, nil
	}
	if strings.EqualFold(clobResp.Status, "delayed") {
		return c.pollUntilFilled(ctx, clobResp.OrderID, in, now)
	}

	// Anything other than "matched" means not immediately filled — cancel and fail.
	slog.Warn("v2_order_no_match", "order_id", clobResp.OrderID,
		"status", clobResp.Status, "limit_px", in.LimitPx,
		"hint", "not immediately filled, cancelling")
	c.tryCancelOrder(context.Background(), clobResp.OrderID)
	errMsg := fmt.Sprintf("order %s but not filled — limit %.4f below best ask, cancelled", clobResp.Status, in.LimitPx)
	return Result{
		OrderID:  clobResp.OrderID,
		Status:   StatusExpired,
		SubmitAt: now,
		Error:    errMsg,
	}, fmt.Errorf("order: %s", errMsg)
}

func (c *V2Client) pollUntilFilled(ctx context.Context, orderID string, in Intent, submitAt time.Time) (Result, error) {
	deadline := time.After(pollTimeout)

	// Sports orders can enter a one-second matching delay.
	select {
	case <-ctx.Done():
		c.tryCancelOrder(context.Background(), orderID)
		return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
			Error: "context cancelled"}, ctx.Err()
	case <-time.After(1 * time.Second):
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	consecutive404 := 0
	const max404 = 3

	for {
		select {
		case <-ctx.Done():
			c.tryCancelOrder(context.Background(), orderID)
			return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
				Error: "context cancelled"}, ctx.Err()
		case <-deadline:
			slog.Warn("v2_poll_timeout", "order_id", orderID, "timeout", pollTimeout)
			c.tryCancelOrder(context.Background(), orderID)
			return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
				Error: "poll timeout"}, nil
		case <-ticker.C:
			os, err := c.GetOrder(ctx, orderID)
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					consecutive404++
					if consecutive404 >= max404 {
						slog.Warn("v2_poll_404_bail", "order_id", orderID, "consecutive", consecutive404)
						return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
							Error: "order not found (404)"}, nil
					}
				}
				slog.Warn("v2_poll_err", "order_id", orderID, "err", err, "404_count", consecutive404)
				continue
			}
			consecutive404 = 0
			status := strings.ToUpper(os.Status)
			switch status {
			case "MATCHED", "ORDER_STATUS_MATCHED":
				slog.Info("v2_order_filled_after_poll", "order_id", orderID,
					"trades", len(os.AssociateTrades), "size_matched", os.SizeMatched)
				avgPx := in.LimitPx
				if os.AvgPrice > 0 {
					avgPx = os.AvgPrice
				}
				filledSize := in.SizeUSD / avgPx
				if os.SizeMatched > 0 {
					filledSize = os.SizeMatched
				}
				return Result{
					OrderID:    orderID,
					Status:     StatusFilled,
					FilledSize: filledSize,
					AvgPrice:   avgPx,
					SubmitAt:   submitAt,
					FilledAt:   time.Now(),
					FeeUSD:     c.fillFeeUSD(ctx, in.Market, filledSize, avgPx),
				}, nil
			case "CANCELLED", "CANCELED", "ORDER_STATUS_CANCELED", "ORDER_STATUS_CANCELED_MARKET_RESOLVED":
				if os.SizeMatched > 0 {
					return c.resultForPartialFill(ctx, orderID, in, submitAt, os), nil
				}
				slog.Info("v2_order_cancelled", "order_id", orderID)
				return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
					Error: "cancelled"}, nil
			default:
				if os.SizeMatched > 0 {
					c.tryCancelOrder(context.Background(), orderID)
					return c.resultForPartialFill(ctx, orderID, in, submitAt, os), nil
				}
				slog.Debug("v2_poll_still_pending", "order_id", orderID, "status", os.Status)
			}
		}
	}
}

type OrderStatusResponse struct {
	ID              string   `json:"id"`
	Status          string   `json:"status"`
	Side            string   `json:"side"`
	MakerAmount     string   `json:"maker_amount"`
	TakerAmount     string   `json:"taker_amount"`
	SizeMatchedRaw  string   `json:"size_matched"`
	SizeMatched     float64  `json:"-"`
	AvgPrice        float64  `json:"-"`
	AssociateTrades []string `json:"associate_trades"`
	OrigPrice       string   `json:"original_price"`
	Price           string   `json:"price"`
}

func (c *V2Client) GetOrder(ctx context.Context, orderID string) (*OrderStatusResponse, error) {
	path := "/data/order/" + orderID
	headers := buildL2Headers(c.creds, c.wallet.Address(), "GET", path, "")

	req, err := http.NewRequestWithContext(ctx, "GET", c.clobBase+path, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s %d: %s", path, resp.StatusCode, body)
	}

	var os OrderStatusResponse
	if err := json.Unmarshal(body, &os); err != nil {
		return nil, fmt.Errorf("parse order status: %w", err)
	}
	os.SizeMatched = parseFixedAmount(os.SizeMatchedRaw)
	os.AvgPrice, _ = strconv.ParseFloat(strings.TrimSpace(os.Price), 64)
	return &os, nil
}

func (c *V2Client) resultForPartialFill(ctx context.Context, orderID string, in Intent, submitAt time.Time, status *OrderStatusResponse) Result {
	avgPrice := status.AvgPrice
	if avgPrice <= 0 {
		avgPrice = in.LimitPx
	}
	return Result{
		OrderID:    orderID,
		Status:     StatusFilled,
		FilledSize: status.SizeMatched,
		AvgPrice:   avgPrice,
		SubmitAt:   submitAt,
		FilledAt:   time.Now(),
		FeeUSD:     c.fillFeeUSD(ctx, in.Market, status.SizeMatched, avgPrice),
	}
}

func (c *V2Client) CancelOrder(ctx context.Context, orderID string) error {
	path := "/order"
	payload, err := json.Marshal(struct {
		OrderID string `json:"orderID"`
	}{OrderID: orderID})
	if err != nil {
		return err
	}
	headers := buildL2Headers(c.creds, c.wallet.Address(), "DELETE", path, string(payload))

	req, err := http.NewRequestWithContext(ctx, "DELETE", c.clobBase+path, bytesReader(payload))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(payload))
	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("DELETE %s %d: %s", path, resp.StatusCode, body)
	}
	if resp.StatusCode == http.StatusOK && len(body) > 0 {
		var result struct {
			Canceled    []string          `json:"canceled"`
			NotCanceled map[string]string `json:"not_canceled"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("DELETE %s decode: %w", path, err)
		}
		for _, canceledID := range result.Canceled {
			if canceledID == orderID {
				slog.Info("v2_order_cancelled", "order_id", orderID)
				return nil
			}
		}
		if reason := result.NotCanceled[orderID]; reason != "" {
			return fmt.Errorf("cancel order %s: %s", orderID, reason)
		}
		return fmt.Errorf("cancel order %s: CLOB did not confirm cancellation", orderID)
	}
	slog.Info("v2_order_cancelled", "order_id", orderID)
	return nil
}

func (c *V2Client) tryCancelOrder(ctx context.Context, orderID string) {
	for i := 0; i < maxCancelAttempts; i++ {
		if err := c.CancelOrder(ctx, orderID); err != nil {
			slog.Warn("v2_cancel_attempt_failed", "order_id", orderID, "attempt", i+1, "err", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return
	}
}

func (c *V2Client) CancelAllOpen(ctx context.Context) error {
	path := "/cancel-all"
	headers := buildL2Headers(c.creds, c.wallet.Address(), "DELETE", path, "")
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.clobBase+path, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	slog.Info("v2_cancel_all_open", "status", resp.StatusCode, "body", string(body))
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("DELETE %s %d: %s", path, resp.StatusCode, body)
	}
	return nil
}

func (c *V2Client) GetOpenOrders(ctx context.Context) (*OpenOrdersResponse, error) {
	path := "/data/orders"
	headers := buildL2Headers(c.creds, c.wallet.Address(), "GET", path, "")
	req, err := http.NewRequestWithContext(ctx, "GET", c.clobBase+path, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf("GET %s response: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s %d: %s", path, resp.StatusCode, body)
	}
	var result OpenOrdersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GET %s decode: %w", path, err)
	}
	return &result, nil
}

func (c *V2Client) GetBalanceAllowance(ctx context.Context) (*BalanceAllowanceResponse, error) {
	path := "/balance-allowance"
	requestURL := c.clobBase + path + "?asset_type=COLLATERAL&signature_type=0"
	headers := buildL2Headers(c.creds, c.wallet.Address(), "GET", path, "")
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("GET %s response: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s %d: %s", path, resp.StatusCode, body)
	}
	var result BalanceAllowanceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GET %s decode: %w", path, err)
	}
	return &result, nil
}

// RefreshBalanceAllowance asks the CLOB to resync its collateral cache from
// chain state. It does not submit a Polygon transaction.
func (c *V2Client) RefreshBalanceAllowance(ctx context.Context) error {
	path := "/balance-allowance/update"
	requestURL := c.clobBase + path + "?asset_type=COLLATERAL&signature_type=0"
	headers := buildL2Headers(c.creds, c.wallet.Address(), "GET", path, "")
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("GET %s response: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("GET %s %d: %s", path, resp.StatusCode, body)
	}
	return nil
}

func GetCLOBVersion(ctx context.Context, clobBase string) (int, error) {
	if clobBase == "" {
		clobBase = ClobBaseURL
	}
	return getCLOBVersion(ctx, clobBase, &http.Client{Timeout: 15 * time.Second})
}

func getCLOBVersion(ctx context.Context, clobBase string, hc *http.Client) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", strings.TrimRight(clobBase, "/")+"/version", nil)
	if err != nil {
		return 0, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return 0, fmt.Errorf("GET /version: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return 0, fmt.Errorf("GET /version response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GET /version %d: %s", resp.StatusCode, body)
	}
	var result struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("GET /version decode: %w", err)
	}
	if result.Version == 0 {
		return 0, fmt.Errorf("GET /version returned empty version")
	}
	return result.Version, nil
}

func amountsForIntent(in Intent) (*big.Int, *big.Int, error) {
	price, ok := new(big.Rat).SetString(strconv.FormatFloat(in.LimitPx, 'f', 6, 64))
	if !ok || price.Sign() <= 0 {
		return nil, nil, fmt.Errorf("invalid limit price %.6f", in.LimitPx)
	}
	num := new(big.Int).Set(price.Num())
	den := new(big.Int).Set(price.Denom())

	if in.Side == Buy {
		target := toMicro(in.SizeUSD)
		k := new(big.Int).Quo(target, num)
		if k.Sign() <= 0 {
			return nil, nil, fmt.Errorf("buy size %.6f too small at %.6f", in.SizeUSD, in.LimitPx)
		}
		return new(big.Int).Mul(num, k), new(big.Int).Mul(den, k), nil
	}

	shares := in.SizeShares
	if shares <= 0 {
		shares = in.SizeUSD / in.LimitPx
	}
	target := toMicro(shares)
	k := new(big.Int).Quo(target, den)
	if k.Sign() <= 0 {
		return nil, nil, fmt.Errorf("sell size %.6f too small at %.6f", shares, in.LimitPx)
	}
	return new(big.Int).Mul(den, k), new(big.Int).Mul(num, k), nil
}

func toMicro(v float64) *big.Int {
	r, ok := new(big.Rat).SetString(strconv.FormatFloat(v, 'f', 6, 64))
	if !ok {
		return new(big.Int)
	}
	r.Mul(r, big.NewRat(usdcScale, 1))
	return new(big.Int).Quo(r.Num(), r.Denom())
}

func parseFixedAmount(raw string) float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if strings.Contains(raw, ".") {
		v, _ := strconv.ParseFloat(raw, 64)
		return v
	}
	v, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return 0
	}
	f, _ := new(big.Rat).SetFrac(v, big.NewInt(usdcScale)).Float64()
	return f
}

func fillFromAmounts(side Side, makingRaw, takingRaw string, fallbackMaker, fallbackTaker *big.Int) (filledSize, avgPrice float64) {
	making := parseFixedAmount(makingRaw)
	taking := parseFixedAmount(takingRaw)
	if making <= 0 {
		making = fixedIntFloat(fallbackMaker)
	}
	if taking <= 0 {
		taking = fixedIntFloat(fallbackTaker)
	}
	if side == Buy {
		filledSize = taking
		if taking > 0 {
			avgPrice = making / taking
		}
		return
	}
	filledSize = making
	if making > 0 {
		avgPrice = taking / making
	}
	return
}

func fixedIntFloat(v *big.Int) float64 {
	if v == nil {
		return 0
	}
	f, _ := new(big.Rat).SetFrac(v, big.NewInt(usdcScale)).Float64()
	return f
}

func (c *V2Client) fillFeeUSD(ctx context.Context, market string, shares, price float64) float64 {
	if shares <= 0 || price <= 0 || price >= 1 {
		return 0
	}
	rate := c.marketFeeRate(ctx, market)
	return shares * rate * price * (1 - price)
}

func (c *V2Client) marketFeeRate(ctx context.Context, market string) float64 {
	if market == "" {
		return 0
	}
	c.feeMu.Lock()
	rate, ok := c.feeRates[market]
	c.feeMu.Unlock()
	if ok {
		return rate
	}

	feeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(feeCtx, "GET", c.clobBase+"/clob-markets/"+market, nil)
	if err != nil {
		return 0
	}
	resp, err := c.http.Do(req)
	if err != nil {
		slog.Warn("v2_fee_rate_fetch_fail", "market", market, "err", err)
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Warn("v2_fee_rate_fetch_fail", "market", market, "status", resp.StatusCode)
		return 0
	}
	var info struct {
		FeeDetails struct {
			Rate float64 `json:"r"`
		} `json:"fd"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		slog.Warn("v2_fee_rate_decode_fail", "market", market, "err", err)
		return 0
	}
	rate = info.FeeDetails.Rate
	c.feeMu.Lock()
	c.feeRates[market] = rate
	c.feeMu.Unlock()
	return rate
}

type sendOrderPayload struct {
	Order     orderJSON `json:"order"`
	Owner     string    `json:"owner"`
	OrderType string    `json:"orderType"`
	DeferExec bool      `json:"deferExec"`
	PostOnly  bool      `json:"postOnly"`
}

type orderJSON struct {
	Maker         string `json:"maker"`
	Signer        string `json:"signer"`
	TokenID       string `json:"tokenId"`
	MakerAmount   string `json:"makerAmount"`
	TakerAmount   string `json:"takerAmount"`
	Side          string `json:"side"`
	Expiration    string `json:"expiration"`
	Timestamp     string `json:"timestamp"`
	Metadata      string `json:"metadata"`
	Builder       string `json:"builder"`
	Signature     string `json:"signature"`
	Salt          int64  `json:"salt"`
	SignatureType int    `json:"signatureType"`
}

type clobOrderResponse struct {
	Success      bool     `json:"success"`
	OrderID      string   `json:"orderID"`
	Status       string   `json:"status"`
	ErrorMsg     string   `json:"errorMsg"`
	MakingAmount string   `json:"makingAmount"`
	TakingAmount string   `json:"takingAmount"`
	TradeIDs     []string `json:"tradeIDs"`
}

type bytesReaderWrapper struct {
	data []byte
	pos  int
}

func bytesReader(data []byte) *bytesReaderWrapper {
	return &bytesReaderWrapper{data: data}
}

func (r *bytesReaderWrapper) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
