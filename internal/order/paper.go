package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

// PaperClient is the no-network fill simulator used during Paper Day 0..7.
// Fills immediately at LimitPx (strategy is expected to pass current mid),
// with optional bp slippage and a configurable per-side fee (bp of notional)
// so net-PnL accounting can model the V2 fee reality ahead of cutover.
type PaperClient struct {
	slippageBp   float64
	feeBp        float64
	takerFeeRate float64
	maxQuoteAge  time.Duration

	mu     sync.Mutex
	orders []Result
}

// NewPaperClient — slippageBp ≥ 0 pulls fill price against you (BUY fills
// higher, SELL fills lower). Pass 0 for clean paper. Default feeBp is 0.
func NewPaperClient(slippageBp float64) *PaperClient {
	return &PaperClient{slippageBp: slippageBp, maxQuoteAge: 5 * time.Second}
}

// NewPaperClientWithFee is the ladder-era constructor that takes both
// slippage and a per-side fee in basis points of notional. The strategy
// layer charges this on each buy + each sell so tranche PnL is net of fees.
func NewPaperClientWithFee(slippageBp, feeBp float64) *PaperClient {
	return &PaperClient{slippageBp: slippageBp, feeBp: feeBp, maxQuoteAge: 5 * time.Second}
}

// NewPaperClientWithFeeModel adds Polymarket's dynamic taker fee curve:
// shares x feeRate x price x (1-price). feeBp remains available for an
// additional flat builder fee on filled notional.
func NewPaperClientWithFeeModel(slippageBp, feeBp, takerFeeRate float64) *PaperClient {
	return &PaperClient{
		slippageBp:   slippageBp,
		feeBp:        feeBp,
		takerFeeRate: takerFeeRate,
		maxQuoteAge:  5 * time.Second,
	}
}

func (p *PaperClient) Name() string { return "paper" }

func (p *PaperClient) Submit(ctx context.Context, in Intent) (Result, error) {
	if err := validate(in); err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}

	now := time.Now().UTC()
	referencePx, model, quoteAge, fillable, reason := p.executionReference(in, now)
	if !fillable {
		r := Result{
			OrderID: "paper-" + randHex(6), Status: StatusExpired,
			SubmitAt: now, Error: reason, ReferencePrice: referencePx,
			ExecutionModel: model, QuoteAge: quoteAge,
		}
		p.record(r)
		return r, nil
	}
	px := applySlippage(referencePx, in.Side, p.slippageBp)
	if px <= 0 || px >= 1 {
		return Result{Status: StatusRejected, Error: "slipped out of (0,1)"},
			fmt.Errorf("paper: slipped price %.4f out of (0,1)", px)
	}
	units := in.SizeUSD / px
	if in.SizeShares > 0 {
		units = in.SizeShares
	}
	notional := units * px
	flatFee := notional * p.feeBp / 10_000
	takerFeeRate := p.takerFeeRate
	if in.TakerFeeRateOverride != nil {
		takerFeeRate = *in.TakerFeeRateOverride
	}
	platformFee := units * takerFeeRate * px * (1 - px)
	if platformFee > 0 {
		platformFee = math.Round(platformFee*100_000) / 100_000
	}
	fee := flatFee + platformFee

	r := Result{
		OrderID:        "paper-" + randHex(6),
		Status:         StatusFilled,
		FilledSize:     units,
		AvgPrice:       px,
		SubmitAt:       now,
		FilledAt:       now,
		FeeUSD:         fee,
		ReferencePrice: referencePx,
		ExecutionModel: model,
		QuoteAge:       quoteAge,
	}
	p.record(r)
	return r, nil
}

func (p *PaperClient) executionReference(in Intent, now time.Time) (float64, string, time.Duration, bool, string) {
	fallback := in.PaperReferencePx
	if fallback <= 0 || fallback >= 1 {
		fallback = in.LimitPx
	}
	if in.PaperQuoteAt.IsZero() {
		if in.PaperRequireQuote {
			return fallback, "quote_required", 0, false, "no executable order-book quote"
		}
		model := "limit_fallback"
		if in.PaperReferencePx > 0 && in.PaperReferencePx < 1 {
			model = "signal_fallback"
		}
		return fallback, model, 0, true, ""
	}
	age := now.Sub(in.PaperQuoteAt)
	if age < 0 {
		age = 0
	}
	if age > p.maxQuoteAge {
		if !in.PaperRequireQuote {
			return fallback, "stale_quote_fallback", age, true, ""
		}
		return fallback, "stale_quote", age, false, "order-book quote is stale"
	}
	quote := in.PaperBestAsk
	depth := in.PaperBestAskSize
	if in.Side == Sell {
		quote = in.PaperBestBid
		depth = in.PaperBestBidSize
	}
	if quote <= 0 || quote >= 1 {
		return quote, "top_of_book", age, false, "no executable top-of-book quote"
	}
	if in.Side == Buy && quote > in.LimitPx+1e-9 {
		return quote, "top_of_book", age, false, "best ask above buy limit"
	}
	if in.Side == Sell && quote < in.LimitPx-1e-9 {
		return quote, "top_of_book", age, false, "best bid below sell limit"
	}
	if in.PaperRequireQuote {
		units := in.SizeShares
		if units <= 0 {
			units = in.SizeUSD / quote
		}
		if depth+1e-9 < units {
			return quote, "top_of_book", age, false, fmt.Sprintf("insufficient top-of-book depth: need %.6f shares, have %.6f", units, depth)
		}
	}
	return quote, "top_of_book", age, true, ""
}

func (p *PaperClient) record(result Result) {
	p.mu.Lock()
	p.orders = append(p.orders, result)
	p.mu.Unlock()
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
	if in.SizeUSD <= 0 || math.IsNaN(in.SizeUSD) || math.IsInf(in.SizeUSD, 0) {
		return fmt.Errorf("non-positive SizeUSD %v", in.SizeUSD)
	}
	if in.LimitPx <= 0 || in.LimitPx >= 1 || math.IsNaN(in.LimitPx) || math.IsInf(in.LimitPx, 0) {
		return fmt.Errorf("LimitPx %v out of (0,1)", in.LimitPx)
	}
	if in.TakerFeeRateOverride != nil && (*in.TakerFeeRateOverride < 0 || *in.TakerFeeRateOverride > 1 || math.IsNaN(*in.TakerFeeRateOverride) || math.IsInf(*in.TakerFeeRateOverride, 0)) {
		return fmt.Errorf("TakerFeeRateOverride %v out of [0,1]", *in.TakerFeeRateOverride)
	}
	if in.PaperReferencePx < 0 || in.PaperReferencePx >= 1 || math.IsNaN(in.PaperReferencePx) || math.IsInf(in.PaperReferencePx, 0) {
		return fmt.Errorf("PaperReferencePx %v out of [0,1)", in.PaperReferencePx)
	}
	for name, quote := range map[string]float64{"PaperBestBid": in.PaperBestBid, "PaperBestAsk": in.PaperBestAsk} {
		if quote < 0 || quote >= 1 || math.IsNaN(quote) || math.IsInf(quote, 0) {
			return fmt.Errorf("%s %v out of [0,1)", name, quote)
		}
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
