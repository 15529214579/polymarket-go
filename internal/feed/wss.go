// Package feed — Polymarket WSS client (skeleton; TODO Phase 1.2).
package feed

import (
	"context"
	"log/slog"
	"time"
)

// WSSClient will maintain a persistent WebSocket connection to
// Polymarket's CLOB market channel (wss://ws-subscriptions-clob.polymarket.com/ws/market)
// and fan out decoded order-book diffs + trades to subscribers.
//
// Phase 1.2 scope (not yet implemented):
//   - dial with backoff (1s → 30s exp)
//   - subscribe to assets_ids (clobTokenIds) in batches
//   - decode "book"/"price_change"/"tick_size_change"/"last_trade_price"
//   - expose channels: Books() <-chan BookEvent, Trades() <-chan TradeEvent
//   - ping/pong + reconnect on 30s silence
type WSSClient struct {
	assetIDs []string
}

func NewWSSClient(assetIDs []string) *WSSClient {
	return &WSSClient{assetIDs: assetIDs}
}

// Run blocks until ctx is canceled. Phase 1.1 stub: logs intent only.
func (w *WSSClient) Run(ctx context.Context) error {
	slog.Info("wss.stub: would subscribe", "n_assets", len(w.assetIDs))
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			slog.Debug("wss.stub alive")
		}
	}
}
