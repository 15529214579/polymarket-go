// Package whale polls target wallets' Polymarket trades via the public
// data API and pushes Telegram alerts for large orders (>threshold USDC).
// Feature-flagged via -whale_enabled; to remove: delete this package +
// WhaleAlert from Notifier.
package whale

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const dataAPI = "https://data-api.polymarket.com"

const maxActivityPages = 11 // offsets 0..5000 at 500 rows per page

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
	ReplayWindow time.Duration
	StatePath    string
	MaxPages     int

	// Legacy single-wallet fields (used when Wallets is empty).
	Wallet     string
	ProfileURL string
	MinSizeUSD float64
}

func DefaultConfig() Config {
	return Config{
		MinSizeUSD:   1000,
		PollInterval: 30 * time.Second,
		ReplayWindow: 0,
		MaxPages:     maxActivityPages,
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
	USDCSize        float64 `json:"usdcSize"`
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
	if t.USDCSize > 0 {
		return t.USDCSize
	}
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

	states    map[string]*walletState // keyed by wallet address
	persistMu sync.Mutex
}

func NewTracker(cfg Config, alert AlertFunc) *Tracker {
	if cfg.MaxPages <= 0 || cfg.MaxPages > maxActivityPages {
		cfg.MaxPages = maxActivityPages
	}
	states := make(map[string]*walletState)
	for _, w := range cfg.ResolvedWallets() {
		states[strings.ToLower(w.Address)] = &walletState{
			lastSeen: make(map[string]struct{}),
		}
	}
	tracker := &Tracker{
		cfg:    cfg,
		http:   &http.Client{Timeout: 15 * time.Second},
		alert:  alert,
		logger: slog.Default(),
		states: states,
	}
	if err := tracker.loadState(); err != nil {
		tracker.logger.Warn("whale_state_load_fail", "path", cfg.StatePath, "err", err)
	}
	return tracker
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
		seedFloor := time.Now()
		if t.cfg.ReplayWindow > 0 {
			seedFloor = seedFloor.Add(-t.cfg.ReplayWindow)
		}
		if err := t.seed(ctx, w); err != nil {
			t.logger.Warn("whale_seed_fail", "wallet", w.Label, "err", err.Error())
			t.advanceFloor(w.Address, seedFloor.Unix())
			if persistErr := t.persistState(); persistErr != nil {
				t.logger.Warn("whale_state_save_fail", "wallet", w.Label, "err", persistErr)
			}
		}
	}

	tk := time.NewTicker(t.cfg.PollInterval)
	defer tk.Stop()
	// Bound concurrent /activity HTTP fans-out so a 70-wallet copytrade
	// list doesn't burst-hammer data-api.polymarket.com on every tick.
	sem := make(chan struct{}, 10)
	var pollCount int
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tk.C:
			pollCount++
			var wg sync.WaitGroup
			for _, w := range wallets {
				wg.Add(1)
				sem <- struct{}{}
				go func(w WalletEntry) {
					defer wg.Done()
					defer func() { <-sem }()
					if err := t.poll(ctx, w); err != nil {
						t.logger.Warn("whale_poll_fail", "wallet", w.Label, "err", err.Error())
					}
				}(w)
			}
			wg.Wait()
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
	st := t.states[strings.ToLower(w.Address)]
	st.mu.Lock()
	persistedTS := st.lastTS
	st.mu.Unlock()
	if persistedTS > 0 {
		return t.poll(ctx, w)
	}

	seedFloor := time.Now()
	since := int64(0)
	maxPages := -1 // latest page is enough when startup replay is disabled
	if t.cfg.ReplayWindow > 0 {
		seedFloor = seedFloor.Add(-t.cfg.ReplayWindow)
		since = seedFloor.Unix()
		maxPages = t.cfg.MaxPages
	}
	trades, err := t.fetchTrades(ctx, w.Address, since, maxPages)
	if err != nil {
		return err
	}
	if len(trades) == 0 {
		t.advanceFloor(w.Address, seedFloor.Unix())
		return t.persistState()
	}
	var maxTS int64
	seen := make(map[string]struct{})
	for _, tr := range trades {
		if tr.Timestamp > maxTS {
			maxTS = tr.Timestamp
			seen = make(map[string]struct{})
		}
		if tr.Timestamp == maxTS {
			seen[tradeKey(&tr)] = struct{}{}
		}
	}
	replayed := 0
	if t.cfg.ReplayWindow > 0 {
		cutoff := time.Now().Add(-t.cfg.ReplayWindow).Unix()
		sort.Slice(trades, func(i, j int) bool {
			return trades[i].Timestamp < trades[j].Timestamp
		})
		for i := range trades {
			tr := &trades[i]
			if tr.Timestamp < cutoff {
				continue
			}
			if t.emitTradeAlert(ctx, w, tr) {
				replayed++
			}
		}
	}
	st.mu.Lock()
	st.lastTS = maxTS
	st.lastSeen = seen
	st.mu.Unlock()
	if err := t.persistState(); err != nil {
		return err
	}
	t.logger.Info("whale_seed_done", "wallet", w.Label, "trades_seen", len(trades), "last_ts", maxTS, "replayed", replayed, "replay_window", t.cfg.ReplayWindow.String())
	return nil
}

type persistedWalletState struct {
	LastTS   int64    `json:"last_ts"`
	LastSeen []string `json:"last_seen,omitempty"`
}

func (t *Tracker) loadState() error {
	if strings.TrimSpace(t.cfg.StatePath) == "" {
		return nil
	}
	raw, err := os.ReadFile(t.cfg.StatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var saved map[string]persistedWalletState
	if err := json.Unmarshal(raw, &saved); err != nil {
		return err
	}
	for address, state := range saved {
		st := t.states[strings.ToLower(address)]
		if st == nil {
			continue
		}
		st.lastTS = state.LastTS
		st.lastSeen = make(map[string]struct{}, len(state.LastSeen))
		for _, key := range state.LastSeen {
			st.lastSeen[key] = struct{}{}
		}
	}
	return nil
}

func (t *Tracker) persistState() error {
	if strings.TrimSpace(t.cfg.StatePath) == "" {
		return nil
	}
	t.persistMu.Lock()
	defer t.persistMu.Unlock()
	saved := make(map[string]persistedWalletState, len(t.states))
	for address, st := range t.states {
		st.mu.Lock()
		state := persistedWalletState{LastTS: st.lastTS, LastSeen: make([]string, 0, len(st.lastSeen))}
		for key := range st.lastSeen {
			state.LastSeen = append(state.LastSeen, key)
		}
		st.mu.Unlock()
		sort.Strings(state.LastSeen)
		saved[address] = state
	}
	raw, err := json.MarshalIndent(saved, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(t.cfg.StatePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".whale-watermarks-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, t.cfg.StatePath); err != nil {
		return err
	}
	dirHandle, err := os.Open(dir)
	if err != nil {
		return err
	}
	syncErr := dirHandle.Sync()
	closeErr := dirHandle.Close()
	return errors.Join(syncErr, closeErr)
}

func (t *Tracker) advanceFloor(wallet string, timestamp int64) {
	st := t.states[strings.ToLower(wallet)]
	if st == nil {
		return
	}
	st.mu.Lock()
	if timestamp > st.lastTS {
		st.lastTS = timestamp
		st.lastSeen = make(map[string]struct{})
	}
	st.mu.Unlock()
}

func cloneSeen(source map[string]struct{}) map[string]struct{} {
	clone := make(map[string]struct{}, len(source))
	for key := range source {
		clone[key] = struct{}{}
	}
	return clone
}

func tradeKey(tr *trade) string {
	payload := fmt.Sprintf("%s\x00%s\x00%s\x00%.12g\x00%.12g\x00%d",
		strings.ToLower(tr.TransactionHash), tr.Asset, strings.ToUpper(tr.Side), tr.Size, tr.Price, tr.Timestamp)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
}

func (t *Tracker) poll(ctx context.Context, w WalletEntry) error {
	st := t.states[strings.ToLower(w.Address)]
	st.mu.Lock()
	lastTS := st.lastTS
	lastSeen := cloneSeen(st.lastSeen)
	st.mu.Unlock()

	trades, err := t.fetchTrades(ctx, w.Address, lastTS, t.cfg.MaxPages)
	if err != nil {
		return err
	}
	if len(trades) == 0 {
		return nil
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].Timestamp < trades[j].Timestamp
	})

	var newTS int64
	newSeen := make(map[string]struct{})

	for i := range trades {
		tr := &trades[i]
		if tr.Timestamp < lastTS {
			continue
		}
		if tr.Timestamp == lastTS {
			if _, dup := lastSeen[tradeKey(tr)]; dup {
				continue
			}
		}

		if tr.Timestamp > newTS {
			newTS = tr.Timestamp
			newSeen = make(map[string]struct{})
		}
		if tr.Timestamp == newTS {
			newSeen[tradeKey(tr)] = struct{}{}
		}

		t.emitTradeAlert(ctx, w, tr)
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
		if err := t.persistState(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tracker) emitTradeAlert(ctx context.Context, w WalletEntry, tr *trade) bool {
	side := strings.ToUpper(tr.Side)
	if side != "BUY" && side != "SELL" {
		return false
	}

	notional := tr.notionalUSD()
	if notional < w.MinSizeUSD {
		return false
	}

	slug := tr.EventSlug
	if slug == "" {
		slug = tr.Slug
	}
	linkURL := fmt.Sprintf(
		"https://newshare.bwb.online/zh/polymarket/event?slug=%s&_nobar=true&_needChain=matic",
		url.QueryEscape(slug),
	)

	ev := AlertEvent{
		Wallet:      w.Address,
		Label:       w.Label,
		ProfileURL:  w.ProfileURL,
		Side:        side,
		SizeUnits:   tr.Size,
		Price:       tr.Price,
		Notional:    notional,
		Question:    tr.Title,
		Slug:        slug,
		Outcome:     tr.Outcome,
		TradeID:     tr.TransactionHash,
		Timestamp:   time.Unix(tr.Timestamp, 0),
		LinkURL:     linkURL,
		AssetID:     tr.Asset,
		ConditionID: tr.ConditionID,
	}

	if positions, err := t.FetchPositions(ctx, w.Address, tr.Asset); err == nil {
		for _, p := range positions {
			if p.Asset == tr.Asset {
				ev.TotalShares = p.Size
				ev.AvgPrice = p.AvgPrice
				if side == "SELL" && p.Size+tr.Size > 0 {
					ev.PctSold = (tr.Size / (p.Size + tr.Size)) * 100
				}
				break
			}
		}
	} else {
		t.logger.Warn("whale_position_fetch_fail", "wallet", w.Label, "err", err.Error())
	}

	if t.alert != nil {
		t.alert(ev)
	}

	t.logger.Info("whale_trade_detected",
		"wallet", w.Label,
		"tx", truncate(tr.TransactionHash, 16),
		"side", tr.Side,
		"notional_usd", notional,
		"market", tr.Title,
		"outcome", tr.Outcome,
	)
	return true
}

func (t *Tracker) fetchTrades(ctx context.Context, wallet string, since int64, maxPages int) ([]trade, error) {
	const pageLimit = 500
	allowTruncated := maxPages < 0
	sortDirection := "ASC"
	if allowTruncated {
		maxPages = -maxPages
		sortDirection = "DESC"
	}
	if maxPages <= 0 || maxPages > maxActivityPages {
		maxPages = maxActivityPages
	}
	var trades []trade
	for page := 0; page < maxPages; page++ {
		offset := page * pageLimit
		q := url.Values{}
		q.Set("user", wallet)
		q.Set("type", "TRADE")
		q.Set("sortBy", "TIMESTAMP")
		q.Set("sortDirection", sortDirection)
		q.Set("limit", fmt.Sprint(pageLimit))
		q.Set("offset", fmt.Sprint(offset))
		if since > 0 {
			q.Set("start", fmt.Sprint(since))
		}

		reqURL := dataAPI + "/activity?" + q.Encode()
		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := t.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("data-api GET /activity: %w", err)
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		closeErr := resp.Body.Close()
		if readErr != nil || closeErr != nil {
			return nil, errors.Join(readErr, closeErr)
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, truncate(string(body), 200))
		}

		var pageTrades []trade
		if err := json.Unmarshal(body, &pageTrades); err != nil {
			return nil, fmt.Errorf("data-api decode: %w (raw: %s)", err, truncate(string(body), 200))
		}
		trades = append(trades, pageTrades...)
		if len(pageTrades) < pageLimit {
			return trades, nil
		}
	}
	if len(trades) == maxPages*pageLimit && !allowTruncated {
		t.logger.Warn("data_api_activity_window_partial",
			"wallet", wallet,
			"since", since,
			"rows", len(trades),
			"next_poll", "resume from newest processed timestamp")
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
	const pageLimit = 500
	var positions []Position
	for offset := 0; offset <= 10000; offset += pageLimit {
		q := url.Values{}
		q.Set("user", wallet)
		q.Set("sizeThreshold", "0.1")
		q.Set("limit", fmt.Sprint(pageLimit))
		q.Set("offset", fmt.Sprint(offset))
		reqURL := dataAPI + "/positions?" + q.Encode()
		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := t.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("data-api GET /positions: %w", err)
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		closeErr := resp.Body.Close()
		if readErr != nil || closeErr != nil {
			return nil, errors.Join(readErr, closeErr)
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("data-api %d: %s", resp.StatusCode, truncate(string(body), 200))
		}
		var page []Position
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("data-api decode positions: %w (raw: %s)", err, truncate(string(body), 200))
		}
		for _, position := range page {
			if assetID == "" || position.Asset == assetID {
				positions = append(positions, position)
			}
		}
		if assetID != "" && len(positions) > 0 {
			return positions, nil
		}
		if len(page) < pageLimit {
			return positions, nil
		}
		if offset == 10000 {
			return nil, errors.New("data-api positions exceeded pagination limit")
		}
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

// LoadWalletsFile reads addresses from a newline-delimited file and returns a
// WalletEntry list ready to drop into Config.Wallets. Each non-empty,
// non-comment line is one 0x… address; the label is auto-derived
// "w_<first6>…<last4>". MinSizeUSD is applied uniformly to every wallet.
func LoadWalletsFile(path string, minSizeUSD float64) ([]WalletEntry, error) {
	return LoadWalletsFileWithListMins(path, minSizeUSD, nil)
}

// LoadWalletsFileWithListMins is LoadWalletsFile plus optional thresholds by
// wallet-list metadata. Lines generated by strategy-lab look like:
//
//	0xabc... # list=sports tier=A ...
//
// When listMinSizeUSD has a matching "sports" entry, that value becomes the
// wallet's MinSizeUSD. Wallets without list metadata keep minSizeUSD.
func LoadWalletsFileWithListMins(path string, minSizeUSD float64, listMinSizeUSD map[string]float64) ([]WalletEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open wallets file %s: %w", path, err)
	}
	defer f.Close()

	var wallets []WalletEntry
	seen := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		comment := ""
		if i := strings.Index(line, "#"); i >= 0 {
			comment = strings.TrimSpace(line[i+1:])
			line = strings.TrimSpace(line[:i])
		}
		// Take the first whitespace-delimited token in case the line has
		// trailing tags ("0xabc... alice").
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		line = fields[0]
		if comment == "" && len(fields) > 1 {
			comment = strings.Join(fields[1:], " ")
		}
		if line == "" {
			continue
		}
		addr := strings.ToLower(line)
		if _, dup := seen[addr]; dup {
			continue
		}
		seen[addr] = struct{}{}
		effectiveMin := minSizeUSD
		if list := walletListTag(comment); list != "" {
			if v, ok := listMinSizeUSD[list]; ok {
				effectiveMin = v
			}
		}
		wallets = append(wallets, WalletEntry{
			Address:    line,
			Label:      walletLabel(line),
			MinSizeUSD: effectiveMin,
		})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read wallets file %s: %w", path, err)
	}
	return wallets, nil
}

func walletListTag(comment string) string {
	for _, field := range strings.Fields(comment) {
		k, v, ok := strings.Cut(field, "=")
		if !ok || strings.ToLower(strings.TrimSpace(k)) != "list" {
			continue
		}
		return strings.ToLower(strings.TrimSpace(v))
	}
	return ""
}

// walletLabel returns the auto-generated short label used for file-loaded
// wallets — e.g. 0xd1c7…31d2b → w_d1c7…1d2b.
func walletLabel(addr string) string {
	if len(addr) < 10 {
		return "w_" + addr
	}
	return "w_" + addr[2:6] + "…" + addr[len(addr)-4:]
}
