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
}

func NewV2Client(wallet *Wallet, creds *APICredentials, negRisk bool) *V2Client {
	return &V2Client{
		wallet:   wallet,
		creds:    creds,
		clobBase: ClobBaseURL,
		http:     &http.Client{Timeout: 30 * time.Second},
		negRisk:  negRisk,
	}
}

func (c *V2Client) Name() string { return "v2-live" }

func (c *V2Client) Submit(ctx context.Context, in Intent) (Result, error) {
	if err := validate(in); err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	exchange := common.HexToAddress(V2ExchangeAddress)
	if c.negRisk {
		exchange = common.HexToAddress(V2NegRiskExchangeAddress)
	}

	now := time.Now()
	priceCents := int64(in.LimitPx*100 + 0.5)
	var makerAmt, takerAmt *big.Int
	side := SigBuy
	if in.Side == Sell {
		side = SigSell
	}
	if side == SigBuy {
		rawTaker := int64(in.SizeUSD / in.LimitPx * usdcScale)
		rawTaker = rawTaker - (rawTaker % 10000)
		takerAmt = big.NewInt(rawTaker)
		makerAmt = big.NewInt(rawTaker * priceCents / 100)
	} else {
		rawMaker := int64(in.SizeUSD / in.LimitPx * usdcScale)
		rawMaker = rawMaker - (rawMaker % 10000)
		makerAmt = big.NewInt(rawMaker)
		takerAmt = big.NewInt(rawMaker * priceCents / 100)
	}

	tokenID := new(big.Int)
	tokenID.SetString(in.AssetID, 10)

	order := V2Order{
		Salt:          NewSalt(),
		Maker:         c.wallet.Address(),
		Signer:        c.wallet.Address(),
		TokenID:       tokenID,
		MakerAmount:   makerAmt,
		TakerAmount:   takerAmt,
		Side:          side,
		SignatureType:  SigTypeEOA,
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

	if clobResp.Status == "matched" {
		slog.Info("v2_order_filled",
			"order_id", clobResp.OrderID,
			"trades", len(clobResp.TradeIDs))
		return Result{
			OrderID:    clobResp.OrderID,
			Status:     StatusFilled,
			FilledSize: in.SizeUSD / in.LimitPx,
			AvgPrice:   in.LimitPx,
			SubmitAt:   now,
			FilledAt:   time.Now(),
		}, nil
	}

	slog.Info("v2_order_pending", "order_id", clobResp.OrderID, "polling_start", true)
	return c.pollUntilFilled(ctx, clobResp.OrderID, in, now)
}

func (c *V2Client) pollUntilFilled(ctx context.Context, orderID string, in Intent, submitAt time.Time) (Result, error) {
	deadline := time.After(pollTimeout)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

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
				slog.Warn("v2_poll_err", "order_id", orderID, "err", err)
				continue
			}
			switch os.Status {
			case "matched":
				slog.Info("v2_order_filled_after_poll", "order_id", orderID,
					"fills", len(os.Trades), "size_matched", os.SizeMatched)
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
				}, nil
			case "cancelled":
				slog.Info("v2_order_cancelled", "order_id", orderID)
				return Result{OrderID: orderID, Status: StatusExpired, SubmitAt: submitAt,
					Error: "cancelled"}, nil
			default:
				slog.Debug("v2_poll_still_pending", "order_id", orderID, "status", os.Status)
			}
		}
	}
}

type OrderStatusResponse struct {
	ID          string           `json:"id"`
	Status      string           `json:"status"`
	Side        string           `json:"side"`
	MakerAmount string           `json:"maker_amount"`
	TakerAmount string           `json:"taker_amount"`
	SizeMatched float64          `json:"-"`
	AvgPrice    float64          `json:"-"`
	Trades      []TradeResponse  `json:"associate_trades"`
	OrigPrice   string           `json:"original_price"`
	Price       string           `json:"price"`
}

type TradeResponse struct {
	ID    string `json:"id"`
	Price string `json:"price"`
	Size  string `json:"size"`
}

func (c *V2Client) GetOrder(ctx context.Context, orderID string) (*OrderStatusResponse, error) {
	path := "/order/" + orderID
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
	os.SizeMatched, os.AvgPrice = computeFillStats(os.Trades)
	return &os, nil
}

func (c *V2Client) CancelOrder(ctx context.Context, orderID string) error {
	path := "/order/" + orderID
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

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("DELETE %s %d: %s", path, resp.StatusCode, body)
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

func computeFillStats(trades []TradeResponse) (totalSize, avgPrice float64) {
	if len(trades) == 0 {
		return 0, 0
	}
	var totalNotional float64
	for _, t := range trades {
		sz, _ := strconv.ParseFloat(strings.TrimSpace(t.Size), 64)
		px, _ := strconv.ParseFloat(strings.TrimSpace(t.Price), 64)
		totalSize += sz
		totalNotional += sz * px
	}
	if totalSize > 0 {
		avgPrice = totalNotional / totalSize
	}
	return
}

type sendOrderPayload struct {
	Order     orderJSON `json:"order"`
	Owner     string    `json:"owner"`
	OrderType string    `json:"orderType"`
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
	Success   bool     `json:"success"`
	OrderID   string   `json:"orderID"`
	Status    string   `json:"status"`
	ErrorMsg  string   `json:"errorMsg"`
	TradeIDs  []string `json:"tradeIDs"`
}

func toFixedPoint(usd float64) *big.Int {
	scaled := int64(usd * usdcScale)
	return big.NewInt(scaled)
}

func roundTo(v *big.Int, granularity int64) *big.Int {
	g := big.NewInt(granularity)
	mod := new(big.Int).Mod(v, g)
	rounded := new(big.Int).Sub(v, mod)
	half := new(big.Int).Div(g, big.NewInt(2))
	if mod.Cmp(half) >= 0 {
		rounded.Add(rounded, g)
	}
	return rounded
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
