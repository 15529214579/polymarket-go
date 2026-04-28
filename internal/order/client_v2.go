package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	ClobBaseURL = "https://clob.polymarket.com"
	usdcScale   = 1_000_000 // 6 decimals
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
	makerAmt := toFixedPoint(in.SizeUSD)
	takerAmt := toFixedPoint(in.SizeUSD / in.LimitPx)
	side := SigBuy
	if in.Side == Sell {
		side = SigSell
		makerAmt = toFixedPoint(in.SizeUSD / in.LimitPx)
		takerAmt = toFixedPoint(in.SizeUSD)
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

	status := StatusPending
	if clobResp.Status == "matched" {
		status = StatusFilled
	}

	return Result{
		OrderID:    clobResp.OrderID,
		Status:     status,
		FilledSize: in.SizeUSD / in.LimitPx,
		AvgPrice:   in.LimitPx,
		SubmitAt:   now,
		FilledAt:   time.Now(),
	}, nil
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
