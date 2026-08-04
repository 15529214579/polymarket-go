package main

import (
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func TestPaperTimeoutExitPriceRequiresBestBid(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		tick feed.Tick
		want float64
		ok   bool
	}{
		{name: "best bid", tick: feed.Tick{BestBid: 0.28, BestBidSize: 100, QuoteTime: now, Mid: 0.40}, want: 0.28, ok: true},
		{name: "no mid fallback", tick: feed.Tick{Mid: 0.40}, ok: false},
		{name: "invalid bid", tick: feed.Tick{BestBid: 1, BestBidSize: 100, QuoteTime: now, Mid: 0.40}, ok: false},
		{name: "stale bid", tick: feed.Tick{BestBid: 0.28, BestBidSize: 100, QuoteTime: now.Add(-time.Minute)}, ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := paperTimeoutExitPrice(tt.tick, now)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("paperTimeoutExitPrice() = %.4f/%v, want %.4f/%v", got, ok, tt.want, tt.ok)
			}
		})
	}
}
