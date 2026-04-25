package odds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const oddsPapiBase = "https://api.oddspapi.io/v4"

// OddsPapiConfig holds flags for the high-frequency OddsPapi scanner.
type OddsPapiConfig struct {
	Enabled   bool
	Interval  time.Duration
	Bookmaker string
	SportKeys []string // e.g. ["soccer_epl", "soccer_spain_la_liga", "soccer_uefa_champs_league"]
}

// DefaultFootballTournaments maps sport keys to OddsPapi tournament IDs.
var DefaultFootballTournaments = map[string]int{
	"soccer_epl":                17, // English Premier League
	"soccer_spain_la_liga":      8,  // La Liga
	"soccer_uefa_champs_league": 7,  // UEFA Champions League
	"soccer_germany_bundesliga": 35, // Bundesliga
	"soccer_italy_serie_a":      23, // Serie A
	"soccer_france_ligue_1":     34, // Ligue 1
}

// leagueFilter maps tournament IDs to required (tournamentSlug, categorySlug) pairs
// to avoid false positives (e.g. "premier-league" in Mongolia).
var leagueFilter = map[int]struct{ slug, category string }{
	17: {"premier-league", "england"},
	8:  {"la-liga", "spain"},
	7:  {"", "international-clubs"}, // UCL
	35: {"bundesliga", "germany"},
	23: {"serie-a", "italy"},
	34: {"ligue-1", "france"},
}

// OddsPapiClient fetches sharp-line odds from OddsPapi (Pinnacle, bet365, etc).
type OddsPapiClient struct {
	apiKey     string
	bookmaker  string
	cacheDir   string
	cacheTTL   time.Duration
	httpClient *http.Client

	mu    sync.Mutex
	usage OddsPapiUsage

	tournamentsMu sync.RWMutex
	tournaments   map[string]int
}

type OddsPapiUsage struct {
	RequestsRemaining int `json:"requests_remaining"`
	RequestsUsed      int `json:"requests_used"`
}

func NewOddsPapiClient(apiKey, bookmaker, cacheDir string) *OddsPapiClient {
	if apiKey == "" {
		apiKey = os.Getenv("ODDSPAPI_API_KEY")
	}
	if bookmaker == "" {
		bookmaker = "pinnacle"
	}
	if cacheDir == "" {
		cacheDir = filepath.Join("db", ".oddspapi_cache")
	}
	return &OddsPapiClient{
		apiKey:      apiKey,
		bookmaker:   bookmaker,
		cacheDir:    cacheDir,
		cacheTTL:    30 * time.Minute,
		httpClient:  &http.Client{Timeout: 20 * time.Second},
		tournaments: make(map[string]int),
	}
}

func (c *OddsPapiClient) HasKey() bool { return c.apiKey != "" }

func (c *OddsPapiClient) Usage() OddsPapiUsage {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.usage
}

// FetchFootballOdds fetches H2H odds for the specified football leagues.
// Two-step: /fixtures (team names) + /odds-by-tournaments (bookmaker odds), merged by fixtureId.
func (c *OddsPapiClient) FetchFootballOdds(ctx context.Context, sportKeys []string) ([]BookmakerOdds, error) {
	if c.apiKey == "" {
		slog.Warn("oddspapi_key_missing")
		return nil, nil
	}
	if len(sportKeys) == 0 {
		sportKeys = []string{
			"soccer_epl",
			"soccer_spain_la_liga",
			"soccer_uefa_champs_league",
		}
	}

	ids, err := c.resolveTournamentIDs(ctx, sportKeys)
	if err != nil {
		return nil, fmt.Errorf("resolve tournament IDs: %w", err)
	}
	if len(ids) == 0 {
		slog.Warn("oddspapi_no_tournament_ids", "keys", sportKeys)
		return nil, nil
	}

	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrs = append(idStrs, strconv.Itoa(id))
	}
	tournamentList := strings.Join(idStrs, ",")

	cacheKey := fmt.Sprintf("oddspapi_%s_%s", c.bookmaker, tournamentList)
	if cached := c.readCache(cacheKey); cached != nil {
		return cached, nil
	}

	// Step 1: Fetch fixtures to get team names.
	nameMap, err := c.fetchFixtureNames(ctx, ids)
	if err != nil {
		slog.Warn("oddspapi_fixtures_failed", "err", err)
		// Non-fatal: proceed with empty names.
	}

	// Step 2: Fetch odds via /odds-by-tournaments (single API call for all tournaments).
	u, _ := url.Parse(oddsPapiBase + "/odds-by-tournaments")
	q := u.Query()
	q.Set("apiKey", c.apiKey)
	q.Set("bookmaker", c.bookmaker)
	q.Set("tournamentIds", tournamentList)
	q.Set("oddsFormat", "decimal")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oddspapi fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("oddspapi http %d: %s", resp.StatusCode, truncBody(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("oddspapi read body: %w", err)
	}

	var fixtures []oddsPapiFixture
	if err := json.Unmarshal(body, &fixtures); err != nil {
		return nil, fmt.Errorf("oddspapi decode: %w", err)
	}

	tidToSport := make(map[int]string)
	for sport, tid := range ids {
		tidToSport[tid] = sport
	}

	result := c.convertFixtures(fixtures, tidToSport, nameMap)
	slog.Info("oddspapi_fetched",
		"bookmaker", c.bookmaker,
		"fixtures", len(fixtures),
		"outcomes", len(result),
		"tournaments", tournamentList,
	)

	if len(result) > 0 {
		c.writeCache(cacheKey, result)
	}
	return result, nil
}

// fetchFixtureNames calls /fixtures to get participant names keyed by fixtureId.
type fixtureMeta struct {
	Home       string
	Away       string
	Tournament string
	Category   string
}

func (c *OddsPapiClient) fetchFixtureNames(ctx context.Context, ids map[string]int) (map[string]fixtureMeta, error) {
	now := time.Now().UTC()
	from := now.Format("2006-01-02T15:04:05Z")
	to := now.Add(7 * 24 * time.Hour).Format("2006-01-02T15:04:05Z")

	u, _ := url.Parse(oddsPapiBase + "/fixtures")
	q := u.Query()
	q.Set("apiKey", c.apiKey)
	q.Set("sportId", "10") // soccer
	q.Set("from", from)
	q.Set("to", to)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fixtures http %d: %s", resp.StatusCode, truncBody(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		FixtureID        string `json:"fixtureId"`
		Participant1Name string `json:"participant1Name"`
		Participant2Name string `json:"participant2Name"`
		TournamentID     int    `json:"tournamentId"`
		TournamentSlug   string `json:"tournamentSlug"`
		CategorySlug     string `json:"categorySlug"`
		TournamentName   string `json:"tournamentName"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// Build set of target tournament IDs.
	targetTIDs := make(map[int]bool)
	for _, tid := range ids {
		targetTIDs[tid] = true
	}

	result := make(map[string]fixtureMeta, len(raw)/10)
	for _, f := range raw {
		if !targetTIDs[f.TournamentID] {
			continue
		}
		// Verify league filter to avoid false positives.
		if filter, ok := leagueFilter[f.TournamentID]; ok {
			if filter.slug != "" && f.TournamentSlug != filter.slug {
				continue
			}
			if filter.category != "" && f.CategorySlug != filter.category {
				continue
			}
		}
		result[f.FixtureID] = fixtureMeta{
			Home:       f.Participant1Name,
			Away:       f.Participant2Name,
			Tournament: f.TournamentName,
			Category:   f.CategorySlug,
		}
	}

	slog.Info("oddspapi_fixtures_loaded", "total", len(raw), "matched", len(result))
	return result, nil
}

func (c *OddsPapiClient) resolveTournamentIDs(_ context.Context, sportKeys []string) (map[string]int, error) {
	result := make(map[string]int)
	var missing []string

	for _, sk := range sportKeys {
		if id, ok := DefaultFootballTournaments[sk]; ok {
			result[sk] = id
		} else {
			missing = append(missing, sk)
		}
	}

	if len(missing) > 0 {
		slog.Warn("oddspapi_unknown_tournaments", "missing", missing)
	}
	return result, nil
}

// convertFixtures transforms OddsPapi response into our BookmakerOdds format.
func (c *OddsPapiClient) convertFixtures(fixtures []oddsPapiFixture, tidToSport map[int]string, nameMap map[string]fixtureMeta) []BookmakerOdds {
	var result []BookmakerOdds

	for _, fix := range fixtures {
		sport := tidToSport[fix.TournamentID]
		if sport == "" {
			sport = fmt.Sprintf("tournament_%d", fix.TournamentID)
		}

		// Get participant names from the name map (fetched via /fixtures).
		meta, hasMeta := nameMap[fix.FixtureID]
		home := fix.Participant1Name
		away := fix.Participant2Name
		if home == "" && hasMeta {
			home = meta.Home
		}
		if away == "" && hasMeta {
			away = meta.Away
		}
		if home == "" {
			home = fmt.Sprintf("Team_%d", fix.Participant1ID)
		}
		if away == "" {
			away = fmt.Sprintf("Team_%d", fix.Participant2ID)
		}
		eventName := home + " vs " + away

		bmOdds, ok := fix.BookmakerOdds[c.bookmaker]
		if !ok {
			continue
		}

		for marketID, market := range bmOdds.Markets {
			if !isH2HMarketID(marketID) {
				continue
			}

			for _, outcome := range market.Outcomes {
				for _, player := range outcome.Players {
					if !player.Active || player.Price <= 1.0 {
						continue
					}

					teamOrSide := resolveOutcomeName(player.BookmakerOutcomeID, home, away)

					impliedProb := 1.0 / player.Price

					result = append(result, BookmakerOdds{
						Sport:             sport,
						EventID:           fix.FixtureID,
						EventName:         eventName,
						TeamOrSide:        teamOrSide,
						BookmakerProb:     impliedProb,
						Bookmaker:         c.bookmaker,
						MarketName:        "h2h",
						EventCommenceTime: fix.StartTime,
					})
				}
			}
		}
	}

	result = normalizeOddsByEvent(result)
	return result
}

func normalizeOddsByEvent(items []BookmakerOdds) []BookmakerOdds {
	type groupKey struct {
		eventID   string
		bookmaker string
	}

	groups := make(map[groupKey][]int)
	for i, item := range items {
		k := groupKey{item.EventID, item.Bookmaker}
		groups[k] = append(groups[k], i)
	}

	for _, indices := range groups {
		total := 0.0
		for _, idx := range indices {
			total += items[idx].BookmakerProb
		}
		if total <= 0 {
			continue
		}
		for _, idx := range indices {
			items[idx].BookmakerProb = items[idx].BookmakerProb / total
			items[idx].BookmakerProb = float64(int(items[idx].BookmakerProb*10000+0.5)) / 10000
		}
	}
	return items
}

// isH2HMarketID checks if the market ID represents a moneyline/1x2 market.
// OddsPapi uses "101" for the standard 1X2 (home/draw/away) moneyline market.
func isH2HMarketID(id string) bool {
	switch id {
	case "101", "1", "1x2", "moneyline", "h2h", "match_winner", "full_time_result":
		return true
	}
	return false
}

// resolveOutcomeName maps OddsPapi bookmakerOutcomeId to team names or "Draw".
func resolveOutcomeName(bookmakerOutcomeID string, home, away string) string {
	switch strings.ToLower(bookmakerOutcomeID) {
	case "home", "1":
		return home
	case "away", "2":
		return away
	case "draw", "x":
		return "Draw"
	}
	return bookmakerOutcomeID
}

// --- OddsPapi API response types ---

type oddsPapiFixture struct {
	FixtureID        string                           `json:"fixtureId"`
	Participant1ID   int                              `json:"participant1Id"`
	Participant2ID   int                              `json:"participant2Id"`
	Participant1Name string                           `json:"participant1Name"`
	Participant2Name string                           `json:"participant2Name"`
	SportID          int                              `json:"sportId"`
	TournamentID     int                              `json:"tournamentId"`
	SeasonID         int                              `json:"seasonId"`
	StatusID         int                              `json:"statusId"`
	StartTime        string                           `json:"startTime"`
	BookmakerOdds    map[string]oddsPapiBookmakerOdds `json:"bookmakerOdds"`
}

type oddsPapiBookmakerOdds struct {
	BookmakerIsActive bool                      `json:"bookmakerIsActive"`
	Markets           map[string]oddsPapiMarket `json:"markets"`
}

type oddsPapiMarket struct {
	Outcomes map[string]oddsPapiOutcome `json:"outcomes"`
}

type oddsPapiOutcome struct {
	Players map[string]oddsPapiPlayer `json:"players"`
}

type oddsPapiPlayer struct {
	Active             bool    `json:"active"`
	Price              float64 `json:"price"`
	BookmakerOutcomeID string  `json:"bookmakerOutcomeId"`
	PlayerName         string  `json:"playerName"`
}

// --- JSON file cache ---

func (c *OddsPapiClient) cachePath(key string) string {
	_ = os.MkdirAll(c.cacheDir, 0o755)
	return filepath.Join(c.cacheDir, key+".json")
}

func (c *OddsPapiClient) readCache(key string) []BookmakerOdds {
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
		return nil
	}
	slog.Info("oddspapi_cache_hit", "key", key, "items", len(entry.Items), "age_sec", int(age.Seconds()))
	return entry.Items
}

func (c *OddsPapiClient) writeCache(key string, items []BookmakerOdds) {
	entry := cacheEntry{
		Ts:    float64(time.Now().Unix()),
		Items: items,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	_ = os.WriteFile(c.cachePath(key), data, 0o644)
}
