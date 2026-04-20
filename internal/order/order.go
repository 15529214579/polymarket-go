// Package order models Polymarket CLOB V2 orders and provides a Client
// abstraction used by strategy + paper layers. Real signing / HTTP lives in
// sibling files behind the same Client interface.
package order

import (
	"context"
	"time"
)

// Side is BUY / SELL. Paper mode only issues BUY for momentum entries, but
// the type is here for completeness.
type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

// OrderType — V2 supports GTC, FOK, GTD, FAK. Strategy only needs GTC-IOC-ish
// for now; keeping the enum open.
type OrderType string

const (
	GTC OrderType = "GTC"
	FOK OrderType = "FOK"
	FAK OrderType = "FAK"
	GTD OrderType = "GTD"
)

// Status is the fill state reported back by the CLOB (or simulated by paper).
type Status string

const (
	StatusPending  Status = "pending"
	StatusFilled   Status = "filled"
	StatusPartial  Status = "partial"
	StatusRejected Status = "rejected"
	StatusExpired  Status = "expired"
)

// Intent is the minimum the strategy layer owns. The signer/client turns this
// into a V2-shaped Order before submission.
type Intent struct {
	AssetID  string // ERC1155 token id
	Market   string // conditionID (for dedupe + logs)
	Side     Side
	SizeUSD  float64 // notional; units = SizeUSD / Price
	LimitPx  float64 // 0..1 probability
	Type     OrderType
	Deadline time.Time // for GTD; zero for GTC
}

// Result is what Submit returns — unified for paper + real.
type Result struct {
	OrderID    string // paper: local uuid; real: CLOB id
	Status     Status
	FilledSize float64 // units filled
	AvgPrice   float64
	SubmitAt   time.Time
	FilledAt   time.Time
	// FeeUSD is the platform fee paid on this fill, in USDC. Paper mode
	// computes it from the configured feeBp; real V2 will populate it from
	// the CLOB fill receipt. Net PnL accounting subtracts this per-leg.
	FeeUSD float64
	Error  string // non-empty only on rejected
}

// Client is the submission surface. Paper + V2-real both satisfy it.
// All impls must be goroutine-safe; the strategy loop may fire concurrently.
type Client interface {
	Submit(ctx context.Context, in Intent) (Result, error)
	Name() string
}
