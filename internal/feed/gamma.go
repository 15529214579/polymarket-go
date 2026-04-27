// Package feed — gamma REST client for market discovery.
package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
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
	ID               string  `json:"id"`
	ConditionID      string  `json:"conditionId"`
	Slug             string  `json:"slug"`
	Question         string  `json:"question"`
	Category         string  `json:"category"`
	Active           bool    `json:"active"`
	Closed           bool    `json:"closed"`
	AcceptingOrders  bool    `json:"acceptingOrders"`
	EndDate          string  `json:"endDate"`
	Volume24hr       float64 `json:"volume24hr"`
	LiquidityClob    float64 `json:"liquidityClob"`
	ClobTokenIDsRaw  string  `json:"clobTokenIds"`
	OutcomePricesRaw string  `json:"outcomePrices"`
	OutcomesRaw      string  `json:"outcomes"`
}

func (m Market) ClobTokenIDs() []string  { return parseStringArray(m.ClobTokenIDsRaw) }
func (m Market) Outcomes() []string      { return parseStringArray(m.OutcomesRaw) }
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

// allowedLoLLeagues — only these leagues pass the LoL filter. LJL (Japan),
// LEC (Europe), and other minor leagues had negative EV in paper trading.
var allowedLoLLeagues = []string{"lck", "lpl"}

// IsLoLMarket — real LoL markets on Polymarket always have "LoL:" prefix
// in the question or "lol-" in the slug, so match on those to avoid matching
// substrings like "election" (contains "lec"). Excludes non-moneyline
// derivatives (handicap / totals / O-U) that ride the same slug prefix.
// Additionally filters to only allowed leagues (LCK, LPL).
func IsLoLMarket(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if !isMoneylineQuestion(q) || !isMoneylineSlug(slug) {
		return false
	}
	isLoL := false
	if strings.HasPrefix(q, "lol:") || strings.HasPrefix(q, "lol ") {
		isLoL = true
	}
	if strings.Contains(q, "league of legends") {
		isLoL = true
	}
	if strings.HasPrefix(slug, "lol-") {
		isLoL = true
	}
	if !isLoL {
		return false
	}
	return isAllowedLoLLeague(q, slug)
}

// excludedLoLKeywords — academy / challengers leagues are minor-tier even
// if they contain "lck" or "lpl" in the name. Filter them out.
var excludedLoLKeywords = []string{"challengers", "academy", "amateur", "developing"}

func isAllowedLoLLeague(q, slug string) bool {
	for _, kw := range excludedLoLKeywords {
		if strings.Contains(q, kw) || strings.Contains(slug, kw) {
			return false
		}
	}
	for _, league := range allowedLoLLeagues {
		if strings.Contains(q, league) || strings.Contains(slug, league) {
			return true
		}
	}
	return false
}

// isMoneylineQuestion — reject questions that are derivatives (handicap,
// totals, over/under, prop, parlay). Polymarket surfaces these under the
// same event slug as the moneyline so the slug-only filter isn't enough.
func isMoneylineQuestion(q string) bool {
	for _, bad := range []string{
		"game handicap", "games total", "total:",
		"over/under", "o/u", "spread", "prop", "parlay",
		"exact score", "end in a draw", "leading at halftime",
		"halftime result", "halftime winner", "first to score",
		"both teams to score", "correct score",
	} {
		if strings.Contains(q, bad) {
			return false
		}
	}
	return true
}

// In-play daily sport matchups: slug shape `<league>-<teamA>-<teamB>-YYYY-MM-DD...`.
// Seasonal futures (e.g. "will-the-lakers-win-the-2026-nba-finals") do not match
// and are intentionally excluded — they don't move on our 60s momentum horizon.
var (
	reNBADaily    = regexp.MustCompile(`^nba-[a-z]{2,4}-[a-z]{2,4}-\d{4}-\d{2}-\d{2}`)
	reNBAPlayoffs = regexp.MustCompile(`^nba-playoffs-`) // series-winner in-play
	reEPLDaily    = regexp.MustCompile(`^epl-[a-z]{2,4}-[a-z]{2,4}-\d{4}-\d{2}-\d{2}`)
	reDota2Daily  = regexp.MustCompile(`^dota2-[a-z0-9]+-[a-z0-9]+-\d{4}-\d{2}-\d{2}`)
	reWTADaily    = regexp.MustCompile(`^wta-[a-z]+-[a-z]+-\d{4}-\d{2}-\d{2}`)
	reATPDaily    = regexp.MustCompile(`^atp-[a-z]+-[a-z]+-\d{4}-\d{2}-\d{2}`)
)

// isMoneylineSlug — exclude derivatives (spread / total / over-under / prop)
// so we only take clean win-probability markets where momentum semantics hold.
func isMoneylineSlug(slug string) bool {
	for _, bad := range []string{"-spread-", "-total-", "-ou-", "-over-", "-under-", "-prop-", "-parlay-"} {
		if strings.Contains(slug, bad) {
			return false
		}
	}
	return true
}

// IsBasketballMarket — NBA daily matchups + NBA playoff series winners, moneyline only.
func IsBasketballMarket(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if !isMoneylineSlug(slug) || !isMoneylineQuestion(q) {
		return false
	}
	return reNBADaily.MatchString(slug) || reNBAPlayoffs.MatchString(slug)
}

// IsFootballMarket — soccer daily matchups (EPL only for now), moneyline only.
func IsFootballMarket(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if !isMoneylineSlug(slug) || !isMoneylineQuestion(q) {
		return false
	}
	return reEPLDaily.MatchString(slug)
}

// IsDota2Market — Dota 2 daily matchups, moneyline only.
// Slug pattern: dota2-<team1>-<team2>-YYYY-MM-DD[-game1|-game2].
func IsDota2Market(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if !isMoneylineSlug(slug) || !isMoneylineQuestion(q) {
		return false
	}
	return reDota2Daily.MatchString(slug)
}

// IsTennisMarket — WTA/ATP daily matchups, moneyline only.
func IsTennisMarket(m Market) bool {
	slug := strings.ToLower(m.Slug)
	if !isMoneylineSlug(slug) {
		return false
	}
	return reWTADaily.MatchString(slug) || reATPDaily.MatchString(slug)
}

// IsSportsMarket — union of LoL + basketball + football (soccer) + Dota 2 + tennis.
// Used for subscription targeting. Keep narrow: only in-play daily / series markets.
func IsSportsMarket(m Market) bool {
	return IsLoLMarket(m) || IsBasketballMarket(m) || IsFootballMarket(m) || IsDota2Market(m) || IsTennisMarket(m)
}

// FilterLoL returns only LoL markets from a list.
func FilterLoL(ms []Market) []Market {
	return filterBy(ms, IsLoLMarket)
}

// FilterSports — LoL + NBA (daily+playoffs) + EPL daily.
func FilterSports(ms []Market) []Market {
	return filterBy(ms, IsSportsMarket)
}

// GetByConditionIDs fetches a batch of markets by their conditionId. The gamma
// /markets endpoint accepts repeated `condition_ids=<hex>` query params and
// returns only matching rows (ignoring active/closed state), which is exactly
// what we want for settlement polling: we need to see closed=true markets too.
func (c *GammaClient) GetByConditionIDs(ctx context.Context, ids []string) ([]Market, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := url.Values{}
	for _, id := range ids {
		if id == "" {
			continue
		}
		q.Add("condition_ids", id)
	}
	q.Set("limit", fmt.Sprintf("%d", len(ids)+5))
	req, err := http.NewRequestWithContext(ctx, "GET", c.base+"/markets?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gamma GET: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gamma %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var out []Market
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("gamma decode: %w", err)
	}
	return out, nil
}

func filterBy(ms []Market, pred func(Market) bool) []Market {
	out := make([]Market, 0, len(ms))
	for _, m := range ms {
		if pred(m) {
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
