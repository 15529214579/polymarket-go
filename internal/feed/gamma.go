// Package feed — gamma REST client for market discovery.
package feed

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

const gammaBase = "https://gamma-api.polymarket.com"

type GammaClient struct {
	http *http.Client
	base string
}

func NewGammaClient() *GammaClient {
	return &GammaClient{
		http: &http.Client{Timeout: 15 * time.Second},
		base: gammaBase,
	}
}

type Market struct {
	ID                string  `json:"id"`
	ConditionID       string  `json:"conditionId"`
	Slug              string  `json:"slug"`
	Question          string  `json:"question"`
	Category          string  `json:"category"`
	Active            bool    `json:"active"`
	Closed            bool    `json:"closed"`
	AcceptingOrders   bool    `json:"acceptingOrders"`
	EndDate           string  `json:"endDate"`
	Volume24hr        float64 `json:"volume24hr"`
	LiquidityClob     float64 `json:"liquidityClob"`
	ClobTokenIDsRaw   string  `json:"clobTokenIds"`
	OutcomePricesRaw  string  `json:"outcomePrices"`
	OutcomesRaw       string  `json:"outcomes"`
}

func (m Market) ClobTokenIDs() []string { return parseStringArray(m.ClobTokenIDsRaw) }
func (m Market) Outcomes() []string     { return parseStringArray(m.OutcomesRaw) }
func (m Market) OutcomePrices() []string { return parseStringArray(m.OutcomePricesRaw) }

func parseStringArray(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}

// ListActiveMarkets paginates through active+open markets.
func (c *GammaClient) ListActiveMarkets(ctx context.Context, pageLimit int) ([]Market, error) {
	if pageLimit <= 0 {
		pageLimit = 500
	}
	var all []Market
	offset := 0
	for {
		q := url.Values{}
		q.Set("active", "true")
		q.Set("closed", "false")
		q.Set("limit", fmt.Sprintf("%d", pageLimit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		q.Set("order", "volume24hr")
		q.Set("ascending", "false")

		req, err := http.NewRequestWithContext(ctx, "GET", c.base+"/markets?"+q.Encode(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("gamma GET: %w", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("gamma %d: %s", resp.StatusCode, truncate(string(body), 200))
		}
		var page []Market
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("gamma decode: %w", err)
		}
		all = append(all, page...)
		if len(page) < pageLimit {
			break
		}
		offset += pageLimit
		if offset >= 5000 {
			break // safety cap
		}
	}
	return all, nil
}

// IsLoLMarket — real LoL markets on Polymarket always have "LoL:" prefix
// in the question or "lol-" in the slug, so match on those to avoid matching
// substrings like "election" (contains "lec").
func IsLoLMarket(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if strings.HasPrefix(q, "lol:") || strings.HasPrefix(q, "lol ") {
		return true
	}
	if strings.Contains(q, "league of legends") {
		return true
	}
	if strings.HasPrefix(slug, "lol-") {
		return true
	}
	return false
}

// FilterLoL returns only LoL markets from a list.
func FilterLoL(ms []Market) []Market {
	out := make([]Market, 0, len(ms))
	for _, m := range ms {
		if IsLoLMarket(m) {
			out = append(out, m)
		}
	}
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
