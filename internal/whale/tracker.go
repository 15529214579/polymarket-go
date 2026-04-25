// Package whale polls a target wallet's Polymarket trades via the public
// data API and pushes Telegram alerts for large orders (>threshold USDC).
// Feature-flagged via -whale_enabled; to remove: delete this package +
// WhaleAlert from Notifier.
package whale

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const dataAPI = "https://data-api.polymarket.com"

type Config struct {
	Enabled      bool
	Wallet       string
	MinSizeUSD   float64
	PollInterval time.Duration
}

func DefaultConfig() Config {
	return Config{
		MinSizeUSD:   1000,
		PollInterval: 30 * time.Second,
	}
}

type trade struct {
	ProxyWallet     string  `json:"proxyWallet"`
	Side            string  `json:"side"`
	Asset           string  `json:"asset"`
	ConditionID     string  `json:"conditionId"`
	Size            float64 `json:"size"`
	Price           float64 `json:"price"`
	Timestamp       int64   `json:"timestamp"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	EventSlug       string  `json:"eventSlug"`
	Outcome         string  `json:"outcome"`
	OutcomeIndex    int     `json:"outcomeIndex"`
	TransactionHash string  `json:"transactionHash"`
	Type            string  `json:"type"`
}

func (t *trade) notionalUSD() float64 {
	return t.Size * t.Price
}

// AlertEvent is the payload passed to the AlertFunc callback for each
// qualifying trade. The caller (main.go) converts this to notify.WhaleAlertEvent.
type AlertEvent struct {
	Wallet      string
	Side        string
	SizeUnits   float64
	Price       float64
	Notional    float64
	Question    string
	Slug        string
	Outcome     string
	TradeID     string
	Timestamp   time.Time
	LinkURL     string
	AssetID     string // CLOB token ID
	ConditionID string // market condition ID
	// Position context (populated when position lookup succeeds).
	TotalShares float64
	AvgPrice    float64
	PctSold     float64 // SELL only: percentage of total position sold (0-100)
}

// AlertFunc is called for each trade that exceeds MinSizeUSD. The caller
// bridges this to the notify.Notifier.WhaleAlert method to avoid a circular
// import between whale → notify.
type AlertFunc func(ev AlertEvent)

type Tracker struct {
	cfg    Config
	http   *http.Client
	alert  AlertFunc
	logger *slog.Logger

	mu        sync.Mutex
	lastTS    int64
	lastSeen  map[string]struct{} // txHash set for dedup within same timestamp
}

func NewTracker(cfg Config, alert AlertFunc) *Tracker {
	return &Tracker{
		cfg:      cfg,
		http:     &http.Client{Timeout: 15 * time.Second},
		alert:    alert,
		logger:   slog.Default(),
		lastSeen: make(map[string]struct{}),
	}
}

func (t *Tracker) Run(ctx context.Context) error {
	if !t.cfg.Enabled || t.cfg.Wallet == "" {
		return nil
	}
	t.logger.Info("whale_tracker.ready",
		"wallet", t.cfg.Wallet,
		"min_size_usd", t.cfg.MinSizeUSD,
		"poll_interval", t.cfg.PollInterval.String(),
	)

	if err := t.seed(ctx); err != nil {
		t.logger.Warn("whale_seed_fail", "err", err.Error())
	}

	tk := time.NewTicker(t.cfg.PollInterval)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tk.C:
			if err := t.poll(ctx); err != nil {
				t.logger.Warn("whale_poll_fail", "err", err.Error())
			}
		}
	}
}

func (t *Tracker) seed(ctx context.Context) error {
	trades, err := t.fetchTrades(ctx)
	if err != nil {
		return err
	}
	if len(trades) == 0 {
		return nil
	}
	var maxTS int64
	seen := make(map[string]struct{})
	for _, tr := range trades {
		if tr.Timestamp > maxTS {
			maxTS = tr.Timestamp
		}
		if tr.TransactionHash != "" {
			seen[tr.TransactionHash] = struct{}{}
		}
	}
	t.mu.Lock()
	t.lastTS = maxTS
	t.lastSeen = seen
	t.mu.Unlock()
	t.logger.Info("whale_seed_done", "trades_seen", len(trades), "last_ts", maxTS)
	return nil
}

func (t *Tracker) poll(ctx context.Context) error {
	trades, err := t.fetchTrades(ctx)
	if err != nil {
		return err
	}
	if len(trades) == 0 {
		return nil
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Timestamp < trades[j].Timestamp
	})

	t.mu.Lock()
	lastTS := t.lastTS
	lastSeen := t.lastSeen
	t.mu.Unlock()

	var newTS int64
	newSeen := make(map[string]struct{})

	for i := range trades {
		tr := &trades[i]
		if tr.Timestamp < lastTS {
			continue
		}
		if tr.Timestamp == lastTS {
			if _, dup := lastSeen[tr.TransactionHash]; dup {
				continue
			}
		}

		if tr.Timestamp > newTS {
			newTS = tr.Timestamp
			newSeen = make(map[string]struct{})
		}
		if tr.Timestamp == newTS && tr.TransactionHash != "" {
			newSeen[tr.TransactionHash] = struct{}{}
		}

		if tr.Type != "" && strings.ToUpper(tr.Type) != "BUY" && strings.ToUpper(tr.Type) != "SELL" {
			continue
		}

		notional := tr.notionalUSD()
		if notional < t.cfg.MinSizeUSD {
			continue
		}

		slug := tr.EventSlug
		if slug == "" {
			slug = tr.Slug
		}
		linkURL := fmt.Sprintf(
			"https://newshare.bwb.online/zh/polymarket/event?slug=%s&_nobar=true&_needChain=matic",
			url.QueryEscape(slug),
		)

		ts := time.Unix(tr.Timestamp, 0)

		ev := AlertEvent{
			Wallet:      t.cfg.Wallet,
			Side:        strings.ToUpper(tr.Side),
			SizeUnits:   tr.Size,
			Price:       tr.Price,
			Notional:    notional,
			Question:    tr.Title,
			Slug:        slug,
			Outcome:     tr.Outcome,
			TradeID:     tr.TransactionHash,
			Timestamp:   ts,
			LinkURL:     linkURL,
			AssetID:     tr.Asset,
			ConditionID: tr.ConditionID,
		}

		if positions, err := t.FetchPositions(ctx, t.cfg.Wallet, tr.Asset); err == nil {
			for _, p := range positions {
				if p.Asset == tr.Asset {
					ev.TotalShares = p.Size
					ev.AvgPrice = p.AvgPrice
					if strings.ToUpper(tr.Side) == "SELL" && p.Size+tr.Size > 0 {
						ev.PctSold = (tr.Size / (p.Size + tr.Size)) * 100
					}
					break
				}
			}
		} else {
			t.logger.Warn("whale_position_fetch_fail", "err", err.Error())
		}

		t.alert(ev)

		t.logger.Info("whale_alert_fired",
			"tx", truncate(tr.TransactionHash, 16),
			"side", tr.Side,
			"notional_usd", notional,
			"market", tr.Title,
			"outcome", tr.Outcome,
		)
	}

	if newTS > 0 {
		t.mu.Lock()
		if newTS > t.lastTS {
			t.lastTS = newTS
			t.lastSeen = newSeen
		} else if newTS == t.lastTS {
			for k, v := range newSeen {
				t.lastSeen[k] = v
			}
		}
		t.mu.Unlock()
	}
	return nil
}

func (t *Tracker) fetchTrades(ctx context.Context) ([]trade, error) {
	q := url.Values{}
	q.Set("user", t.cfg.Wallet)
	q.Set("limit", "50")

	reqURL := dataAPI + "/activity?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("data-api GET /activity: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, truncate(string(body), 200))
	}

	var trades []trade
	if err := json.Unmarshal(body, &trades); err != nil {
		return nil, fmt.Errorf("data-api decode: %w (raw: %s)", err, truncate(string(body), 200))
	}
	return trades, nil
}

// Position represents a wallet's holding for a single outcome token,
// returned by the data-api /positions endpoint.
type Position struct {
	Size        float64 `json:"size"`
	AvgPrice    float64 `json:"avgPrice"`
	TotalBought float64 `json:"totalBought"`
	RealizedPnL float64 `json:"realizedPnl"`
	CurPrice    float64 `json:"curPrice"`
	Title       string  `json:"title"`
	Outcome     string  `json:"outcome"`
	Asset       string  `json:"asset"`
	ConditionID string  `json:"conditionId"`
}

// FetchPositions returns all positions for the given wallet, optionally
// filtered to a single asset (pass "" to get all).
func (t *Tracker) FetchPositions(ctx context.Context, wallet, assetID string) ([]Position, error) {
	q := url.Values{}
	q.Set("user", wallet)
	q.Set("sizeThreshold", "0.1")
	q.Set("limit", "100")
	if assetID != "" {
		q.Set("asset", assetID)
	}
	reqURL := dataAPI + "/positions?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("data-api GET /positions: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var positions []Position
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, fmt.Errorf("data-api decode positions: %w (raw: %s)", err, truncate(string(body), 200))
	}
	return positions, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
