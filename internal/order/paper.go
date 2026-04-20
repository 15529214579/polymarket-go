package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// PaperClient is the no-network fill simulator used during Paper Day 0..7.
// Fills immediately at LimitPx (strategy is expected to pass current mid),
// with optional bp slippage and a configurable per-side fee (bp of notional)
// so net-PnL accounting can model the V2 fee reality ahead of cutover.
type PaperClient struct {
	slippageBp float64
	feeBp      float64

	mu     sync.Mutex
	orders []Result
}

// NewPaperClient — slippageBp ≥ 0 pulls fill price against you (BUY fills
// higher, SELL fills lower). Pass 0 for clean paper. Default feeBp is 0.
func NewPaperClient(slippageBp float64) *PaperClient {
	return &PaperClient{slippageBp: slippageBp}
}

// NewPaperClientWithFee is the ladder-era constructor that takes both
// slippage and a per-side fee in basis points of notional. The strategy
// layer charges this on each buy + each sell so tranche PnL is net of fees.
func NewPaperClientWithFee(slippageBp, feeBp float64) *PaperClient {
	return &PaperClient{slippageBp: slippageBp, feeBp: feeBp}
}

func (p *PaperClient) Name() string { return "paper" }

func (p *PaperClient) Submit(ctx context.Context, in Intent) (Result, error) {
	if err := validate(in); err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	now := time.Now().UTC()
	px := applySlippage(in.LimitPx, in.Side, p.slippageBp)
	if px <= 0 || px >= 1 {
		return Result{Status: StatusRejected, Error: "slipped out of (0,1)"},
			fmt.Errorf("paper: slipped price %.4f out of (0,1)", px)
	}
	units := in.SizeUSD / px
	fee := in.SizeUSD * p.feeBp / 10_000

	r := Result{
		OrderID:    "paper-" + randHex(6),
		Status:     StatusFilled,
		FilledSize: units,
		AvgPrice:   px,
		SubmitAt:   now,
		FilledAt:   now,
		FeeUSD:     fee,
	}

	p.mu.Lock()
	p.orders = append(p.orders, r)
	p.mu.Unlock()
	return r, nil
}

// History returns a copy of all paper fills so far (safe to call from any
// goroutine; cheap enough for paper volumes).
func (p *PaperClient) History() []Result {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]Result, len(p.orders))
	copy(out, p.orders)
	return out
}

func validate(in Intent) error {
	if in.AssetID == "" {
		return errors.New("empty AssetID")
	}
	if in.Side != Buy && in.Side != Sell {
		return fmt.Errorf("bad side %q", in.Side)
	}
	if in.SizeUSD <= 0 {
		return fmt.Errorf("non-positive SizeUSD %v", in.SizeUSD)
	}
	if in.LimitPx <= 0 || in.LimitPx >= 1 {
		return fmt.Errorf("LimitPx %v out of (0,1)", in.LimitPx)
	}
	return nil
}

func applySlippage(px float64, side Side, bp float64) float64 {
	if bp == 0 {
		return px
	}
	adj := px * bp / 10_000
	if side == Buy {
		return px + adj
	}
	return px - adj
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
