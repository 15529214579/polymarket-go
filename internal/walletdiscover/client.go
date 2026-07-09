package walletdiscover

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	gammaBase string
	dataBase  string
	lbBase    string
	http      *http.Client
}

func NewClient(cfg Config) *Client {
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
	return &Client{
		gammaBase: gammaBase,
		dataBase:  dataBase,
		lbBase:    lbBase,
		http: &http.Client{
			Timeout: 30 * time.Second,
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "polymarket-go-wallet-discover/1.0")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("GET %s: HTTP %d: %s", reqURL, resp.StatusCode, trunc(string(body), 300))
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("decode %s: %w: %s", reqURL, err, trunc(string(body), 300))
	}
	return nil
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
