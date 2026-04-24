// Package odds fetches live bookmaker odds from The Odds API,
// converts decimal odds to juice-removed implied probabilities,
// and caches results to stay within the 500 req/month free tier.
package odds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// H2H sport keys (match-level moneyline).
var DefaultH2HSportKeys = []string{
	"soccer_epl",
	"soccer_uefa_champs_league",
	"soccer_spain_la_liga",
	"soccer_germany_bundesliga",
	"soccer_italy_serie_a",
	"basketball_nba",
	"americanfootball_nfl",
	"mma_mixed_martial_arts",
}

// Esports sport keys.
var EsportsSportKeys = []string{
	"esports_lol",
	"esports_csgo",
	"esports_valorant",
	"esports_dota2",
}

// Outright/futures sport keys.
var OutrightSportKeys = []string{
	"basketball_nba_championship_winner",
	"icehockey_nhl_championship_winner",
	"soccer_fifa_world_cup_winner",
	"golf_masters_tournament_winner",
}

const oddsAPIBase = "https://api.the-odds-api.com/v4/sports"

// APIUsage tracks quota from response headers.
type APIUsage struct {
	RequestsRemaining int `json:"requests_remaining"`
	RequestsUsed      int `json:"requests_used"`
}

// Client fetches odds from The Odds API with JSON file caching.
type Client struct {
	apiKey   string
	cacheDir string
	cacheTTL time.Duration
	http     *http.Client

	mu    sync.Mutex
	usage APIUsage
}

// NewClient creates an odds API client. apiKey defaults to ODDS_API_KEY env var.
// cacheDir defaults to db/.odds_cache under the working directory.
func NewClient(apiKey, cacheDir string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("ODDS_API_KEY")
	}
	if cacheDir == "" {
		cacheDir = filepath.Join("db", ".odds_cache")
	}
	return &Client{
		apiKey:   apiKey,
		cacheDir: cacheDir,
		cacheTTL: time.Hour,
		http:     &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Usage() APIUsage {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.usage
}

// FetchH2H fetches head-to-head odds for the given sport keys.
func (c *Client) FetchH2H(ctx context.Context, sportKeys []string) ([]BookmakerOdds, error) {
	if c.apiKey == "" {
		slog.Warn("odds_api_key_missing")
		return nil, nil
	}
	if len(sportKeys) == 0 {
		sportKeys = DefaultH2HSportKeys
	}
	return c.fetchOdds(ctx, sportKeys, "h2h", "us,eu,uk")
}

// FetchOutrights fetches outright/futures odds.
func (c *Client) FetchOutrights(ctx context.Context, sportKeys []string) ([]BookmakerOdds, error) {
	if c.apiKey == "" {
		slog.Warn("odds_api_key_missing")
		return nil, nil
	}
	if len(sportKeys) == 0 {
		sportKeys = OutrightSportKeys
	}
	return c.fetchOdds(ctx, sportKeys, "outrights", "us,eu,uk")
}

func (c *Client) fetchOdds(ctx context.Context, sportKeys []string, markets, regions string) ([]BookmakerOdds, error) {
	sorted := make([]string, len(sportKeys))
	copy(sorted, sportKeys)
	sort.Strings(sorted)
	cacheKey := fmt.Sprintf("%s_%s", markets, strings.Join(sorted, "_"))

	if cached := c.readCache(cacheKey); cached != nil {
		return cached, nil
	}

	var all []BookmakerOdds
	for _, sport := range sportKeys {
		items, err := c.fetchSport(ctx, sport, markets, regions)
		if err != nil {
			slog.Warn("odds_fetch_fail", "sport", sport, "err", err.Error())
			continue
		}
		all = append(all, items...)
	}
	slog.Info("odds_fetched", "markets", markets, "sports", len(sportKeys), "outcomes", len(all))

	if len(all) > 0 {
		c.writeCache(cacheKey, all)
	}
	return all, nil
}

func (c *Client) fetchSport(ctx context.Context, sport, markets, regions string) ([]BookmakerOdds, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/%s/odds", oddsAPIBase, sport))
	q := u.Query()
	q.Set("apiKey", c.apiKey)
	q.Set("regions", regions)
	q.Set("markets", markets)
	q.Set("oddsFormat", "decimal")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("odds api %s: %w", sport, err)
	}
	defer resp.Body.Close()

	c.mu.Lock()
	if v := resp.Header.Get("x-requests-remaining"); v != "" {
		fmt.Sscanf(v, "%d", &c.usage.RequestsRemaining)
	}
	if v := resp.Header.Get("x-requests-used"); v != "" {
		fmt.Sscanf(v, "%d", &c.usage.RequestsUsed)
	}
	c.mu.Unlock()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("odds api %s http %d: %s", sport, resp.StatusCode, truncBody(body))
	}

	var events []oddsEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("odds api %s decode: %w", sport, err)
	}

	var result []BookmakerOdds
	for _, ev := range events {
		eventName := ev.HomeTeam + " vs " + ev.AwayTeam
		if ev.HomeTeam == "" && ev.AwayTeam == "" {
			eventName = ev.SportTitle
			if eventName == "" {
				eventName = ev.ID
			}
		}
		for _, bm := range ev.Bookmakers {
			for _, mkt := range bm.Markets {
				// Only accept pure h2h (reject h2h_lay etc per TASK-067)
				if markets == "h2h" && mkt.Key != "h2h" {
					continue
				}
				if len(mkt.Outcomes) == 0 {
					continue
				}
				decimalOdds := make([]float64, len(mkt.Outcomes))
				names := make([]string, len(mkt.Outcomes))
				valid := true
				for i, o := range mkt.Outcomes {
					decimalOdds[i] = o.Price
					names[i] = o.Name
					if o.Price <= 0 {
						valid = false
					}
				}
				if !valid {
					continue
				}
				fairProbs := removeJuice(decimalOdds)
				for i, name := range names {
					result = append(result, BookmakerOdds{
						Sport:             sport,
						EventID:           ev.ID,
						EventName:         eventName,
						TeamOrSide:        name,
						BookmakerProb:     math.Round(fairProbs[i]*10000) / 10000,
						Bookmaker:         bm.Key,
						MarketName:        mkt.Key,
						EventCommenceTime: ev.CommenceTime,
					})
				}
			}
		}
	}
	slog.Info("odds_sport_fetched", "sport", sport, "outcomes", len(result))
	return result, nil
}

// removeJuice converts European decimal odds to fair probabilities.
// p_i = (1/odds_i) / sum(1/odds_j)
func removeJuice(decimalOdds []float64) []float64 {
	raw := make([]float64, len(decimalOdds))
	total := 0.0
	for i, o := range decimalOdds {
		if o > 0 {
			raw[i] = 1.0 / o
			total += raw[i]
		}
	}
	if total == 0 {
		return make([]float64, len(decimalOdds))
	}
	fair := make([]float64, len(decimalOdds))
	for i := range raw {
		fair[i] = raw[i] / total
	}
	return fair
}

// --- JSON file cache ---

type cacheEntry struct {
	Ts    float64         `json:"ts"`
	Items []BookmakerOdds `json:"items"`
}

func (c *Client) cachePath(key string) string {
	_ = os.MkdirAll(c.cacheDir, 0o755)
	return filepath.Join(c.cacheDir, key+".json")
}

func (c *Client) readCache(key string) []BookmakerOdds {
	p := c.cachePath(key)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil
	}
	age := time.Since(time.Unix(int64(entry.Ts), 0))
	if age > c.cacheTTL {
		slog.Info("odds_cache_expired", "key", key, "age_sec", int(age.Seconds()))
		return nil
	}
	slog.Info("odds_cache_hit", "key", key, "items", len(entry.Items), "age_sec", int(age.Seconds()))
	return entry.Items
}

func (c *Client) writeCache(key string, items []BookmakerOdds) {
	entry := cacheEntry{
		Ts:    float64(time.Now().Unix()),
		Items: items,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	p := c.cachePath(key)
	_ = os.WriteFile(p, data, 0o644)
	slog.Info("odds_cache_write", "key", key, "items", len(items))
}

// --- API response types ---

type oddsEvent struct {
	ID           string           `json:"id"`
	SportKey     string           `json:"sport_key"`
	SportTitle   string           `json:"sport_title"`
	CommenceTime string           `json:"commence_time"`
	HomeTeam     string           `json:"home_team"`
	AwayTeam     string           `json:"away_team"`
	Bookmakers   []oddsBookmaker  `json:"bookmakers"`
}

type oddsBookmaker struct {
	Key     string       `json:"key"`
	Title   string       `json:"title"`
	Markets []oddsMarket `json:"markets"`
}

type oddsMarket struct {
	Key      string        `json:"key"`
	Outcomes []oddsOutcome `json:"outcomes"`
}

type oddsOutcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func truncBody(b []byte) string {
	s := string(b)
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}
