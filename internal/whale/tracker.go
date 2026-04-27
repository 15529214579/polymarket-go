// Package whale polls target wallets' Polymarket trades via the public
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

// WalletEntry describes one tracked whale.
type WalletEntry struct {
	Address    string
	Label      string
	MinSizeUSD float64
	ProfileURL string
}

type Config struct {
	Enabled      bool
	Wallets      []WalletEntry
	PollInterval time.Duration

	// Legacy single-wallet fields (used when Wallets is empty).
	Wallet     string
	ProfileURL string
	MinSizeUSD float64
}

func DefaultConfig() Config {
	return Config{
		MinSizeUSD:   1000,
		PollInterval: 30 * time.Second,
	}
}

// ResolvedWallets returns the effective wallet list, falling back to the
// legacy single-wallet fields if Wallets is empty.
func (c Config) ResolvedWallets() []WalletEntry {
	if len(c.Wallets) > 0 {
		return c.Wallets
	}
	if c.Wallet == "" {
		return nil
	}
	return []WalletEntry{{
		Address:    c.Wallet,
		Label:      shortAddr(c.Wallet),
		MinSizeUSD: c.MinSizeUSD,
		ProfileURL: c.ProfileURL,
	}}
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
	Label       string // human-readable whale label (e.g. "drpufferfish")
	ProfileURL  string
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

// walletState holds per-wallet dedup state.
type walletState struct {
	mu       sync.Mutex
	lastTS   int64
	lastSeen map[string]struct{}
}

type Tracker struct {
	cfg    Config
	http   *http.Client
	alert  AlertFunc
	logger *slog.Logger

	states map[string]*walletState // keyed by wallet address
}

func NewTracker(cfg Config, alert AlertFunc) *Tracker {
	states := make(map[string]*walletState)
	for _, w := range cfg.ResolvedWallets() {
		states[strings.ToLower(w.Address)] = &walletState{
			lastSeen: make(map[string]struct{}),
		}
	}
	return &Tracker{
		cfg:    cfg,
		http:   &http.Client{Timeout: 15 * time.Second},
		alert:  alert,
		logger: slog.Default(),
		states: states,
	}
}

func (t *Tracker) Run(ctx context.Context) error {
	wallets := t.cfg.ResolvedWallets()
	if !t.cfg.Enabled || len(wallets) == 0 {
		return nil
	}

	for _, w := range wallets {
		t.logger.Info("whale_tracker.ready",
			"wallet", w.Address,
			"label", w.Label,
			"min_size_usd", w.MinSizeUSD,
			"poll_interval", t.cfg.PollInterval.String(),
		)
	}

	// Seed all wallets.
	for _, w := range wallets {
		if err := t.seed(ctx, w); err != nil {
			t.logger.Warn("whale_seed_fail", "wallet", w.Label, "err", err.Error())
		}
	}

	tk := time.NewTicker(t.cfg.PollInterval)
	defer tk.Stop()
	var pollCount int
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tk.C:
			pollCount++
			for _, w := range wallets {
				if err := t.poll(ctx, w); err != nil {
					t.logger.Warn("whale_poll_fail", "wallet", w.Label, "err", err.Error())
				}
			}
			if pollCount%10 == 0 {
				for _, w := range wallets {
					st := t.states[strings.ToLower(w.Address)]
					st.mu.Lock()
					ts := st.lastTS
					st.mu.Unlock()
					t.logger.Info("whale_poll_heartbeat",
						"wallet", w.Label,
						"polls", pollCount,
						"last_ts", ts,
						"last_age_sec", time.Now().Unix()-ts,
					)
				}
			}
		}
	}
}

func (t *Tracker) seed(ctx context.Context, w WalletEntry) error {
	trades, err := t.fetchTrades(ctx, w.Address)
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
	st := t.states[strings.ToLower(w.Address)]
	st.mu.Lock()
	st.lastTS = maxTS
	st.lastSeen = seen
	st.mu.Unlock()
	t.logger.Info("whale_seed_done", "wallet", w.Label, "trades_seen", len(trades), "last_ts", maxTS)
	return nil
}

func (t *Tracker) poll(ctx context.Context, w WalletEntry) error {
	trades, err := t.fetchTrades(ctx, w.Address)
	if err != nil {
		return err
	}
	if len(trades) == 0 {
		return nil
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Timestamp < trades[j].Timestamp
	})

	st := t.states[strings.ToLower(w.Address)]
	st.mu.Lock()
	lastTS := st.lastTS
	lastSeen := st.lastSeen
	st.mu.Unlock()

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
		if notional < w.MinSizeUSD {
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
			Wallet:      w.Address,
			Label:       w.Label,
			ProfileURL:  w.ProfileURL,
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

		if positions, err := t.FetchPositions(ctx, w.Address, tr.Asset); err == nil {
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
			t.logger.Warn("whale_position_fetch_fail", "wallet", w.Label, "err", err.Error())
		}

		t.alert(ev)

		t.logger.Info("whale_alert_fired",
			"wallet", w.Label,
			"tx", truncate(tr.TransactionHash, 16),
			"side", tr.Side,
			"notional_usd", notional,
			"market", tr.Title,
			"outcome", tr.Outcome,
		)
	}

	if newTS > 0 {
		st.mu.Lock()
		if newTS > st.lastTS {
			st.lastTS = newTS
			st.lastSeen = newSeen
		} else if newTS == st.lastTS {
			for k, v := range newSeen {
				st.lastSeen[k] = v
			}
		}
		st.mu.Unlock()
	}
	return nil
}

func (t *Tracker) fetchTrades(ctx context.Context, wallet string) ([]trade, error) {
	q := url.Values{}
	q.Set("user", wallet)
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

func shortAddr(addr string) string {
	if len(addr) > 10 {
		return addr[:6] + "…" + addr[len(addr)-4:]
	}
	return addr
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// ParseWallets parses a comma-separated wallet spec string.
// Format: "addr|label|minUSD|profileURL,addr|label|minUSD|profileURL,..."
func ParseWallets(s string) ([]WalletEntry, error) {
	if s == "" {
		return nil, nil
	}
	var wallets []WalletEntry
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.SplitN(part, "|", 4)
		if len(fields) < 3 {
			return nil, fmt.Errorf("invalid wallet spec %q (want addr|label|minUSD[|profileURL])", part)
		}
		var minUSD float64
		if _, err := fmt.Sscanf(fields[2], "%f", &minUSD); err != nil {
			return nil, fmt.Errorf("invalid minUSD %q in wallet spec %q", fields[2], part)
		}
		we := WalletEntry{
			Address:    fields[0],
			Label:      fields[1],
			MinSizeUSD: minUSD,
		}
		if len(fields) >= 4 {
			we.ProfileURL = fields[3]
		}
		wallets = append(wallets, we)
	}
	return wallets, nil
}
