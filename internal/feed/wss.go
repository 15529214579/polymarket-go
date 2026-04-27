// Package feed — Polymarket CLOB WSS client.
//
// Connects to wss://ws-subscriptions-clob.polymarket.com/ws/market,
// subscribes to a batch of assetIDs (clobTokenIds), and fans out
// decoded events to subscriber channels. On disconnect, reconnects
// with exponential backoff (1s → 30s) and re-subscribes.
package feed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wssURL         = "wss://ws-subscriptions-clob.polymarket.com/ws/market"
	pingInterval   = 10 * time.Second
	readIdleLimit  = 45 * time.Second
	backoffMin     = 1 * time.Second
	backoffMax     = 30 * time.Second
	maxAssetsBatch = 500
)

// Level is one side of the order book.
type Level struct {
	Price float64
	Size  float64
}

// BookEvent is a full book snapshot (type="book") or post-diff reconstructed
// state (we apply price_change locally and re-emit). Bids descending, asks ascending.
type BookEvent struct {
	AssetID   string
	Market    string
	Timestamp time.Time
	Bids      []Level
	Asks      []Level
	Raw       string // best-effort preserved for debugging
}

// TradeEvent is a last_trade_price message.
type TradeEvent struct {
	AssetID   string
	Market    string
	Timestamp time.Time
	Price     float64
	Size      float64
	Side      string // BUY | SELL (taker side)
}

// WSSClient maintains one persistent connection covering up to ~500 assets.
type WSSClient struct {
	assetIDs []string
	books    chan BookEvent
	trades   chan TradeEvent

	// lastEventNs is a unix-nano stamp of the most recent book/trade we
	// decoded. Read by the risk feed-silence watchdog (SPEC §6).
	lastEventNs atomic.Int64

	// connected is true while a WSS session is live. The risk watchdog uses
	// this to distinguish "the feed is quiet because the market is quiet"
	// (don't trip) from "the feed is quiet because the socket is gone"
	// (trip after grace period).
	connected atomic.Bool

	mu       sync.Mutex
	orderbks map[string]*bookState // per-assetID local reconstruction
}

type bookState struct {
	market string
	bids   map[string]float64 // price→size
	asks   map[string]float64
	hash   string
}

// NewWSSClient returns a client that will subscribe to assetIDs on Run.
func NewWSSClient(assetIDs []string) *WSSClient {
	if len(assetIDs) > maxAssetsBatch {
		assetIDs = assetIDs[:maxAssetsBatch]
	}
	return &WSSClient{
		assetIDs: assetIDs,
		books:    make(chan BookEvent, 256),
		trades:   make(chan TradeEvent, 256),
		orderbks: make(map[string]*bookState),
	}
}

// Books returns decoded book events (snapshots + reconstructed after diffs).
func (w *WSSClient) Books() <-chan BookEvent { return w.books }

// Trades returns last_trade_price events.
func (w *WSSClient) Trades() <-chan TradeEvent { return w.trades }

// LastEventAt returns the timestamp of the most recent decoded event. Zero
// time means no event has arrived yet.
func (w *WSSClient) LastEventAt() time.Time {
	ns := w.lastEventNs.Load()
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}

// Connected reports whether a WSS session is currently live. Flips to true
// right after a successful subscribe and back to false when runOnce returns
// (which is always followed by backoff + redial).
func (w *WSSClient) Connected() bool { return w.connected.Load() }

// Run blocks until ctx is canceled. Reconnects on error.
func (w *WSSClient) Run(ctx context.Context) error {
	if len(w.assetIDs) == 0 {
		return errors.New("wss: no assetIDs")
	}
	backoff := backoffMin
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		start := time.Now()
		err := w.runOnce(ctx)
		elapsed := time.Since(start)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		slog.Warn("wss: disconnected", "err", err, "uptime", elapsed.Round(time.Second).String())
		// reset backoff if the connection lived ≥1 min
		if elapsed > time.Minute {
			backoff = backoffMin
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		if backoff < backoffMax {
			backoff *= 2
			if backoff > backoffMax {
				backoff = backoffMax
			}
		}
	}
}

func (w *WSSClient) runOnce(ctx context.Context) error {
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	dctx, dcancel := context.WithTimeout(ctx, 15*time.Second)
	defer dcancel()
	conn, resp, err := dialer.DialContext(dctx, wssURL, http.Header{})
	if err != nil {
		if resp != nil {
			return fmt.Errorf("dial: %w (http %d)", err, resp.StatusCode)
		}
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()
	slog.Info("wss: connected", "n_assets", len(w.assetIDs))

	sub := map[string]any{
		"assets_ids": w.assetIDs,
		"type":       "market",
	}
	if err := conn.WriteJSON(sub); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	w.connected.Store(true)
	defer w.connected.Store(false)

	// reader deadline loop
	_ = conn.SetReadDeadline(time.Now().Add(readIdleLimit))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(readIdleLimit))
		return nil
	})

	// pinger
	pingCtx, pingCancel := context.WithCancel(ctx)
	defer pingCancel()
	go func() {
		t := time.NewTicker(pingInterval)
		defer t.Stop()
		for {
			select {
			case <-pingCtx.Done():
				return
			case <-t.C:
				_ = conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second))
			}
		}
	}()

	// read loop
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		_ = conn.SetReadDeadline(time.Now().Add(readIdleLimit))
		w.dispatch(data)
	}
}

// Polymarket sends both single-object and array payloads; handle both.
func (w *WSSClient) dispatch(data []byte) {
	trim := firstNonSpace(data)
	if trim == '[' {
		var arr []json.RawMessage
		if err := json.Unmarshal(data, &arr); err == nil {
			for _, m := range arr {
				w.dispatchOne(m)
			}
			return
		}
	}
	w.dispatchOne(data)
}

func (w *WSSClient) dispatchOne(data []byte) {
	var head struct {
		EventType string `json:"event_type"`
		AssetID   string `json:"asset_id"`
		Market    string `json:"market"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		slog.Debug("wss: non-json frame", "len", len(data))
		return
	}
	if os.Getenv("WSS_DUMP") != "" {
		n := len(data)
		if n > 400 {
			n = 400
		}
		slog.Info("wss.raw", "event", head.EventType, "len", len(data), "raw", string(data[:n]))
	}
	ts := parseTS(head.Timestamp)

	switch head.EventType {
	case "book":
		var p struct {
			AssetID string     `json:"asset_id"`
			Market  string     `json:"market"`
			Bids    []rawLevel `json:"bids"`
			Asks    []rawLevel `json:"asks"`
			Hash    string     `json:"hash"`
		}
		if err := json.Unmarshal(data, &p); err != nil {
			slog.Warn("wss_book_parse_err", "err", err)
			return
		}
		w.applyBookSnapshot(p.AssetID, p.Market, p.Bids, p.Asks, p.Hash)
		w.emitBook(p.AssetID, p.Market, ts, string(data))
	case "price_change":
		var p struct {
			Market       string           `json:"market"`
			PriceChanges []rawPriceChange `json:"price_changes"`
		}
		if err := json.Unmarshal(data, &p); err != nil {
			slog.Warn("wss_price_change_parse_err", "err", err)
			return
		}
		// group changes by asset_id so we emit one book event per asset
		byAsset := map[string][]rawPriceChange{}
		for _, c := range p.PriceChanges {
			if c.AssetID == "" {
				continue
			}
			byAsset[c.AssetID] = append(byAsset[c.AssetID], c)
		}
		for aid, changes := range byAsset {
			w.applyPriceChanges(aid, p.Market, changes)
			w.emitBook(aid, p.Market, ts, "")
		}
	case "last_trade_price":
		var p struct {
			AssetID string `json:"asset_id"`
			Market  string `json:"market"`
			Price   string `json:"price"`
			Size    string `json:"size"`
			Side    string `json:"side"`
		}
		if err := json.Unmarshal(data, &p); err != nil {
			slog.Warn("wss_trade_parse_err", "err", err)
			return
		}
		ev := TradeEvent{
			AssetID:   p.AssetID,
			Market:    p.Market,
			Timestamp: ts,
			Price:     parseFloat(p.Price),
			Size:      parseFloat(p.Size),
			Side:      p.Side,
		}
		w.lastEventNs.Store(time.Now().UnixNano())
		select {
		case w.trades <- ev:
		default:
			// drop if consumer is slow
		}
	case "tick_size_change":
		// silent; informational
	default:
		// Polymarket sometimes sends keep-alive / ack frames without event_type
	}
}

type rawLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// rawPriceChange mirrors one element of the `price_changes` array sent by the
// CLOB WSS. `best_bid`/`best_ask` are included on every change and let us skip
// full book reconciliation when we only care about top-of-book.
type rawPriceChange struct {
	AssetID string `json:"asset_id"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Hash    string `json:"hash"`
	BestBid string `json:"best_bid"`
	BestAsk string `json:"best_ask"`
}

func (w *WSSClient) applyBookSnapshot(assetID, market string, bids, asks []rawLevel, hash string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	st := &bookState{market: market, bids: map[string]float64{}, asks: map[string]float64{}}
	for _, l := range bids {
		if sz := parseFloat(l.Size); sz > 0 {
			st.bids[l.Price] = sz
		}
	}
	for _, l := range asks {
		if sz := parseFloat(l.Size); sz > 0 {
			st.asks[l.Price] = sz
		}
	}
	st.hash = hash
	w.orderbks[assetID] = st
}

func (w *WSSClient) applyPriceChanges(assetID, market string, changes []rawPriceChange) {
	w.mu.Lock()
	defer w.mu.Unlock()
	st, ok := w.orderbks[assetID]
	if !ok {
		st = &bookState{market: market, bids: map[string]float64{}, asks: map[string]float64{}}
		w.orderbks[assetID] = st
	}
	for _, c := range changes {
		sz := parseFloat(c.Size)
		m := st.asks
		if c.Side == "BUY" || c.Side == "buy" {
			m = st.bids
		}
		if sz <= 0 {
			delete(m, c.Price)
		} else {
			m[c.Price] = sz
		}
		if c.Hash != "" {
			st.hash = c.Hash
		}
	}
}

func (w *WSSClient) emitBook(assetID, market string, ts time.Time, raw string) {
	w.mu.Lock()
	st, ok := w.orderbks[assetID]
	if !ok {
		w.mu.Unlock()
		return
	}
	bids := levelsFromMap(st.bids, true)
	asks := levelsFromMap(st.asks, false)
	if market == "" {
		market = st.market
	}
	w.mu.Unlock()

	ev := BookEvent{
		AssetID:   assetID,
		Market:    market,
		Timestamp: ts,
		Bids:      bids,
		Asks:      asks,
		Raw:       raw,
	}
	w.lastEventNs.Store(time.Now().UnixNano())
	select {
	case w.books <- ev:
	default:
	}
}

func levelsFromMap(m map[string]float64, desc bool) []Level {
	out := make([]Level, 0, len(m))
	for p, s := range m {
		out = append(out, Level{Price: parseFloat(p), Size: s})
	}
	// simple insertion sort is fine for typical book sizes
	for i := 1; i < len(out); i++ {
		for j := i; j > 0; j-- {
			swap := false
			if desc {
				swap = out[j].Price > out[j-1].Price
			} else {
				swap = out[j].Price < out[j-1].Price
			}
			if !swap {
				break
			}
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseTS(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil || ms == 0 {
		return time.Now()
	}
	// polymarket timestamps are ms since epoch
	return time.UnixMilli(ms)
}

func firstNonSpace(b []byte) byte {
	for _, c := range b {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return c
		}
	}
	return 0
}
