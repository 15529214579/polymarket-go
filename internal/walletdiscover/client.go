package walletdiscover

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Client struct {
	gammaBase   string
	dataBase    string
	lbBase      string
	http        *http.Client
	maxAttempts int
	retryBase   time.Duration
	retryMax    time.Duration
	retries     atomic.Int64
	rateLimits  atomic.Int64
	failures    atomic.Int64
}

func NewClient(cfg Config) *Client {
	def := DefaultConfig()
	gammaBase := cfg.GammaBase
	if gammaBase == "" {
		gammaBase = defaultGammaBase
	}
	dataBase := cfg.DataBase
	if dataBase == "" {
		dataBase = defaultDataBase
	}
	lbBase := cfg.LeaderboardBase
	if lbBase == "" {
		lbBase = defaultLBBase
	}
	if cfg.HTTPMaxAttempts <= 0 {
		cfg.HTTPMaxAttempts = def.HTTPMaxAttempts
	}
	if cfg.HTTPRetryBase <= 0 {
		cfg.HTTPRetryBase = def.HTTPRetryBase
	}
	if cfg.HTTPRetryMax <= 0 {
		cfg.HTTPRetryMax = def.HTTPRetryMax
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = def.HTTPTimeout
	}
	return &Client{
		gammaBase:   gammaBase,
		dataBase:    dataBase,
		lbBase:      lbBase,
		maxAttempts: cfg.HTTPMaxAttempts,
		retryBase:   cfg.HTTPRetryBase,
		retryMax:    cfg.HTTPRetryMax,
		http: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

func (c *Client) ListMarkets(ctx context.Context, limit, offset int) ([]Market, error) {
	q := url.Values{}
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("active", "true")
	q.Set("closed", "false")
	q.Set("order", "volume24hr")
	q.Set("ascending", "false")
	var out []Market
	if err := c.getJSON(ctx, c.gammaBase+"/markets?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Leaderboard(ctx context.Context, kind, window string, limit int) ([]LeaderboardEntry, error) {
	return c.LeaderboardPage(ctx, kind, window, limit, 0)
}

func (c *Client) LeaderboardPage(ctx context.Context, kind, window string, limit, offset int) ([]LeaderboardEntry, error) {
	if kind == "" {
		kind = "profit"
	}
	if window == "" {
		window = "all"
	}
	if limit <= 0 {
		limit = 100
	}
	q := url.Values{}
	q.Set("window", window)
	q.Set("limit", fmt.Sprintf("%d", limit))
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	var out []LeaderboardEntry
	if err := c.getJSON(ctx, c.lbBase+"/"+kind+"?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RecentTrades(ctx context.Context, limit, offset int) ([]Trade, error) {
	q := url.Values{}
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	var out []Trade
	if err := c.getJSON(ctx, c.dataBase+"/trades?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Holders(ctx context.Context, conditionID string, limit int) ([]HolderResponse, error) {
	q := url.Values{}
	q.Set("market", conditionID)
	q.Set("limit", fmt.Sprintf("%d", limit))
	var out []HolderResponse
	if err := c.getJSON(ctx, c.dataBase+"/holders?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Activity(ctx context.Context, wallet string, limit, offset int) ([]Trade, error) {
	q := url.Values{}
	q.Set("user", wallet)
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	var out []Trade
	if err := c.getJSON(ctx, c.dataBase+"/activity?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ClosedPositions(ctx context.Context, wallet string, limit int) ([]ClosedPosition, error) {
	q := url.Values{}
	q.Set("user", wallet)
	q.Set("limit", fmt.Sprintf("%d", limit))
	var out []ClosedPosition
	if err := c.getJSON(ctx, c.dataBase+"/closed-positions?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) getJSON(ctx context.Context, reqURL string, dst any) error {
	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		retryable, retryAfter, status, err := c.getJSONAttempt(ctx, reqURL, dst)
		if err == nil {
			return nil
		}
		if status == http.StatusTooManyRequests {
			c.rateLimits.Add(1)
		}
		if !retryable || attempt == c.maxAttempts || ctx.Err() != nil {
			c.failures.Add(1)
			return err
		}
		delay := c.retryDelay(attempt)
		if retryAfter > delay {
			delay = retryAfter
		}
		if delay > c.retryMax {
			delay = c.retryMax
		}
		c.retries.Add(1)
		slog.Warn("wallet_discover.http_retry", "attempt", attempt, "max_attempts", c.maxAttempts, "status", status, "delay", delay, "url", reqURL, "err", err)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			c.failures.Add(1)
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("GET %s: retry loop exhausted", reqURL)
}

func (c *Client) getJSONAttempt(ctx context.Context, reqURL string, dst any) (bool, time.Duration, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return false, 0, 0, err
	}
	req.Header.Set("User-Agent", "polymarket-go-wallet-discover/1.0")
	resp, err := c.http.Do(req)
	if err != nil {
		return ctx.Err() == nil, 0, 0, err
	}
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	_ = resp.Body.Close()
	if readErr != nil {
		return true, 0, resp.StatusCode, readErr
	}
	if resp.StatusCode >= 400 {
		err = fmt.Errorf("GET %s: HTTP %d: %s", reqURL, resp.StatusCode, trunc(string(body), 300))
		return retryableStatus(resp.StatusCode), parseRetryAfter(resp.Header.Get("Retry-After")), resp.StatusCode, err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return false, 0, resp.StatusCode, fmt.Errorf("decode %s: %w: %s", reqURL, err, trunc(string(body), 300))
	}
	return false, 0, resp.StatusCode, nil
}

func (c *Client) retryDelay(attempt int) time.Duration {
	delay := c.retryBase
	for i := 1; i < attempt && delay < c.retryMax; i++ {
		delay *= 2
	}
	if delay > c.retryMax {
		return c.retryMax
	}
	return delay
}

func (c *Client) Stats() HTTPStats {
	return HTTPStats{
		Retries:    c.retries.Load(),
		RateLimits: c.rateLimits.Load(),
		Failures:   c.failures.Load(),
	}
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusTooEarly, http.StatusTooManyRequests:
		return true
	default:
		return status >= 500 && status <= 599
	}
}

func parseRetryAfter(raw string) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(raw); err == nil {
		if seconds <= 0 {
			return 0
		}
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(raw); err == nil {
		if delay := time.Until(when); delay > 0 {
			return delay
		}
	}
	return 0
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func IsMaxHistoricalOffsetError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "max historical activity offset")
}
