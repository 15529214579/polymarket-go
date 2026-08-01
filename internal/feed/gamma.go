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
	"strconv"
	"strings"
	"time"
)

const (
	gammaBase = "https://gamma-api.polymarket.com"
	clobBase  = "https://clob.polymarket.com"
)

type GammaClient struct {
	http     *http.Client
	base     string
	clobBase string
}

func NewGammaClient() *GammaClient {
	return &GammaClient{
		http:     &http.Client{Timeout: 15 * time.Second},
		base:     gammaBase,
		clobBase: clobBase,
	}
}

type MarketEvent struct {
	StartTime string `json:"startTime"`
}

type Market struct {
	ID               string        `json:"id"`
	ConditionID      string        `json:"conditionId"`
	Slug             string        `json:"slug"`
	Question         string        `json:"question"`
	Category         string        `json:"category"`
	Active           bool          `json:"active"`
	Closed           bool          `json:"closed"`
	AcceptingOrders  bool          `json:"acceptingOrders"`
	EndDate          string        `json:"endDate"`
	GameStartTime    string        `json:"gameStartTime"`
	Volume24hr       float64       `json:"volume24hr"`
	LiquidityClob    float64       `json:"liquidityClob"`
	ClobTokenIDsRaw  string        `json:"clobTokenIds"`
	OutcomePricesRaw string        `json:"outcomePrices"`
	OutcomesRaw      string        `json:"outcomes"`
	NegRisk          bool          `json:"negRisk"`
	Events           []MarketEvent `json:"events"`
}

func (m Market) ClobTokenIDs() []string  { return parseStringArray(m.ClobTokenIDsRaw) }
func (m Market) Outcomes() []string      { return parseStringArray(m.OutcomesRaw) }
func (m Market) OutcomePrices() []string { return parseStringArray(m.OutcomePricesRaw) }

// EventStartTime returns the scheduled event start exposed by Gamma. Sports
// markets usually carry it on the event; some feeds also expose it directly.
func (m Market) EventStartTime() time.Time {
	if t := parseMarketTime(m.GameStartTime); !t.IsZero() {
		return t
	}
	for _, event := range m.Events {
		if t := parseMarketTime(event.StartTime); !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

type CLOBMarketInfo struct {
	GameStart    time.Time
	TakerFeeRate float64
	FeeRateKnown bool
}

// GetCLOBMarketInfo reads the authoritative game start and fee curve used by
// CLOB. Both values are market-specific and may legitimately be zero.
func (c *GammaClient) GetCLOBMarketInfo(ctx context.Context, conditionID string) (CLOBMarketInfo, error) {
	if strings.TrimSpace(conditionID) == "" {
		return CLOBMarketInfo{}, fmt.Errorf("empty condition id")
	}
	base := c.clobBase
	if base == "" {
		base = clobBase
	}
	req, err := http.NewRequestWithContext(ctx, "GET", base+"/clob-markets/"+url.PathEscape(conditionID), nil)
	if err != nil {
		return CLOBMarketInfo{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return CLOBMarketInfo{}, fmt.Errorf("clob market info GET: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return CLOBMarketInfo{}, err
	}
	if resp.StatusCode >= 400 {
		return CLOBMarketInfo{}, fmt.Errorf("clob market info %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var info struct {
		GameStart  string `json:"gst"`
		FeeDetails *struct {
			Rate float64 `json:"r"`
		} `json:"fd"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return CLOBMarketInfo{}, fmt.Errorf("clob market info decode: %w", err)
	}
	out := CLOBMarketInfo{GameStart: parseMarketTime(info.GameStart)}
	if info.FeeDetails != nil {
		out.TakerFeeRate = info.FeeDetails.Rate
		out.FeeRateKnown = true
	}
	return out, nil
}

// GetCLOBEventStart is kept for callers that only need sports timing.
func (c *GammaClient) GetCLOBEventStart(ctx context.Context, conditionID string) (time.Time, error) {
	info, err := c.GetCLOBMarketInfo(ctx, conditionID)
	if err != nil {
		return time.Time{}, err
	}
	if info.GameStart.IsZero() {
		return time.Time{}, fmt.Errorf("clob market info has no game start")
	}
	return info.GameStart, nil
}

func parseMarketTime(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05Z"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return time.Time{}
}

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
		pageLimit = 100
	}
	if pageLimit > 100 {
		pageLimit = 100
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
			if resp.StatusCode == http.StatusUnprocessableEntity && len(all) > 0 && strings.Contains(strings.ToLower(string(body)), "offset too large") {
				break
			}
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

// IsDerivativeFollowMarketText rejects derivatives (handicap, totals,
// over/under, prop, parlay). Polymarket surfaces these under the same event
// slug as the moneyline so the slug-only filter isn't enough.
func IsDerivativeFollowMarketText(text string) bool {
	text = strings.ToLower(text)
	if strings.Contains(text, " and ") || strings.Contains(text, " & ") {
		return true
	}
	for _, bad := range []string{
		"game handicap", "games total", "total:",
		"over/under", "o/u", "spread", "prop", "parlay",
		"exact score", "end in a draw", "leading at halftime",
		"halftime result", "halftime winner", "first to score",
		"both teams to score", "correct score", "to win 2-0",
		"draw (1-1)", "draw at halftime", "match result", "penalty shootout",
		"neither team to score", "to score first", "score first",
		"extra time", "second half", "first half",
		"shots", "goals during", "+ goals", "goalscorer", "goal scorer",
		"cards", "corners",
		" score ",
		"first blood", "roshan", "rampage", "ultra kill",
		"barracks", "daytime", "total kills", "both teams beat",
	} {
		if strings.Contains(text, bad) {
			return true
		}
	}
	return false
}

var (
	footballScoreResultRE = regexp.MustCompile(`\b(?:to win|draw)\s*\(?\d+\s*[-:]\s*\d+\)?`)
	footballScoreLabelRE  = regexp.MustCompile(`\bscore\s*[:=-]?\s*\d+\s*[-:]\s*\d+\b`)
)

// IsFootballScoreMarketText identifies football correct-score markets. These
// stay outside automatic following, but are useful to the push-only monitor.
func IsFootballScoreMarketText(text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return false
	}
	for _, otherSport := range []string{
		"basketball", "nba", "wnba", "tennis", " atp", " wta",
		"league of legends", " lol", "dota", "cs2", "csgo", "valorant",
		"series score", "map score", "game score",
	} {
		if strings.Contains(text, otherSport) {
			return false
		}
	}
	if strings.Contains(text, "exact score") || strings.Contains(text, "correct score") {
		return true
	}
	if footballScoreResultRE.MatchString(text) {
		return true
	}
	return footballScoreLabelRE.MatchString(text)
}

// IsOutrightFollowMarketText rejects long-horizon championship/futures markets.
// They can attract large orders but do not fit the short-window whale-follow
// signals used by the sports/ esports copy-trading system.
func IsOutrightFollowMarketText(text string) bool {
	text = strings.ToLower(text)
	for _, bad := range []string{
		"world cup winner",
		"win the world cup",
		"win the 2026 fifa world cup",
		"win the fifa world cup",
		"win the nba finals",
		"win the stanley cup",
		"win the super bowl",
		"win the champions league",
		"win the premier league",
		"win dota 2 the international",
		"win the international",
		"win lol worlds",
		"win worlds",
	} {
		if strings.Contains(text, bad) {
			return true
		}
	}
	for _, badSlug := range []string{
		"world-cup-winner",
		"win-the-2026-fifa-world-cup",
		"win-the-fifa-world-cup",
		"win-the-nba-finals",
		"win-the-stanley-cup",
		"win-the-super-bowl",
		"win-the-champions-league",
		"win-the-premier-league",
		"win-dota-2-the-international",
		"win-the-international",
		"win-lol-worlds",
		"win-worlds",
	} {
		if strings.Contains(text, badSlug) {
			return true
		}
	}
	return false
}

// isMoneylineQuestion keeps only moneyline-style questions.
func isMoneylineQuestion(q string) bool {
	if IsDerivativeFollowMarketText(q) || IsOutrightFollowMarketText(q) {
		return false
	}
	return true
}

// In-play daily sport matchups: slug shape `<league>-<teamA>-<teamB>-YYYY-MM-DD...`.
// Seasonal futures (e.g. "will-the-lakers-win-the-2026-nba-finals") do not match
// and are intentionally excluded — they don't move on our 60s momentum horizon.
var (
	reNBADaily                = regexp.MustCompile(`^nba-[a-z]{2,4}-[a-z]{2,4}-\d{4}-\d{2}-\d{2}`)
	reNBAPlayoffs             = regexp.MustCompile(`^nba-playoffs-`) // series-winner in-play
	reEPLDaily                = regexp.MustCompile(`^epl-[a-z]{2,4}-[a-z]{2,4}-\d{4}-\d{2}-\d{2}`)
	reFIFWCDaily              = regexp.MustCompile(`^fifwc-[a-z0-9]+-[a-z0-9]+-\d{4}-\d{2}-\d{2}`)
	reNationalFootballDaily   = regexp.MustCompile(`^will ([a-z][a-z .'-]+) win on \d{4}-\d{2}-\d{2}\??$`)
	reNationalFootballAdvance = regexp.MustCompile(`^([a-z][a-z .'-]+) vs\.? ([a-z][a-z .'-]+): team to advance\??$`)
	reDota2Daily              = regexp.MustCompile(`^dota2-[a-z0-9]+-[a-z0-9]+-\d{4}-\d{2}-\d{2}`)
	reWTADaily                = regexp.MustCompile(`^wta-[a-z]+-[a-z]+-\d{4}-\d{2}-\d{2}`)
	reATPDaily                = regexp.MustCompile(`^atp-[a-z]+-[a-z]+-\d{4}-\d{2}-\d{2}`)
)

var footballNations = map[string]struct{}{
	"algeria": {}, "argentina": {}, "australia": {}, "austria": {},
	"belgium": {}, "brazil": {}, "canada": {}, "chile": {}, "china": {},
	"colombia": {}, "croatia": {}, "denmark": {}, "ecuador": {},
	"egypt": {}, "england": {}, "france": {}, "germany": {}, "ghana": {},
	"greece": {}, "hungary": {}, "iran": {}, "ireland": {}, "italy": {},
	"japan": {}, "mexico": {}, "morocco": {}, "netherlands": {},
	"new zealand": {}, "nigeria": {}, "norway": {}, "paraguay": {},
	"peru": {}, "poland": {}, "portugal": {}, "qatar": {}, "romania": {},
	"saudi arabia": {}, "scotland": {}, "senegal": {}, "serbia": {},
	"south africa": {}, "south korea": {}, "spain": {}, "sweden": {},
	"switzerland": {}, "tunisia": {}, "turkey": {}, "ukraine": {},
	"united states": {}, "uruguay": {}, "usa": {}, "wales": {},
}

var basketballTeamAliases = [][]string{
	{"atlanta hawks", "hawks"},
	{"boston celtics", "celtics"},
	{"brooklyn nets", "nets"},
	{"charlotte hornets", "hornets"},
	{"chicago bulls", "bulls"},
	{"cleveland cavaliers", "cavaliers", "cavs"},
	{"dallas mavericks", "mavericks", "mavs"},
	{"denver nuggets", "nuggets"},
	{"detroit pistons", "pistons"},
	{"golden state warriors", "warriors"},
	{"houston rockets", "rockets"},
	{"indiana pacers", "pacers"},
	{"los angeles clippers", "clippers"},
	{"los angeles lakers", "lakers"},
	{"memphis grizzlies", "grizzlies"},
	{"miami heat"},
	{"milwaukee bucks", "bucks"},
	{"minnesota timberwolves", "timberwolves", "wolves"},
	{"new orleans pelicans", "pelicans"},
	{"new york knicks", "knicks"},
	{"oklahoma city thunder", "thunder"},
	{"orlando magic", "magic"},
	{"philadelphia 76ers", "76ers", "sixers"},
	{"phoenix suns", "suns"},
	{"portland trail blazers", "trail blazers"},
	{"sacramento kings", "kings"},
	{"san antonio spurs", "spurs"},
	{"toronto raptors", "raptors"},
	{"utah jazz"},
	{"washington wizards", "wizards"},
	{"atlanta dream", "dream"},
	{"chicago sky"},
	{"connecticut sun"},
	{"dallas wings", "wings"},
	{"golden state valkyries", "valkyries"},
	{"indiana fever", "fever"},
	{"las vegas aces", "aces"},
	{"los angeles sparks", "sparks"},
	{"minnesota lynx", "lynx"},
	{"new york liberty", "liberty"},
	{"phoenix mercury", "mercury"},
	{"seattle storm", "storm"},
	{"toronto tempo", "tempo"},
	{"washington mystics", "mystics"},
}

// isMoneylineSlug — exclude derivatives (spread / total / over-under / prop)
// so we only take clean win-probability markets where momentum semantics hold.
func isMoneylineSlug(slug string) bool {
	if IsOutrightFollowMarketText(slug) {
		return false
	}
	for _, bad := range []string{"-spread-", "-total-", "-ou-", "-over-", "-under-", "-prop-", "-parlay-", "-match-result-"} {
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
	return reNBADaily.MatchString(slug) || reNBAPlayoffs.MatchString(slug) || isBasketballTeamMatchQuestion(q)
}

func isBasketballTeamMatchQuestion(q string) bool {
	if !(strings.Contains(q, " vs ") || strings.Contains(q, " vs. ") || strings.Contains(q, " at ")) {
		return false
	}
	teams := 0
	for _, aliases := range basketballTeamAliases {
		for _, alias := range aliases {
			if containsPhrase(q, alias) {
				teams++
				break
			}
		}
		if teams >= 2 {
			return true
		}
	}
	return false
}

func containsPhrase(text, phrase string) bool {
	if phrase == "" {
		return false
	}
	return strings.Contains(" "+text+" ", " "+phrase+" ")
}

// IsFootballMarket — soccer daily matchups, moneyline only.
func IsFootballMarket(m Market) bool {
	q := strings.ToLower(m.Question)
	slug := strings.ToLower(m.Slug)
	if !isMoneylineSlug(slug) || !isMoneylineQuestion(q) {
		return false
	}
	return reEPLDaily.MatchString(slug) || reFIFWCDaily.MatchString(slug) || isNationalFootballWinQuestion(q) || isNationalFootballAdvanceQuestion(q)
}

func isNationalFootballWinQuestion(q string) bool {
	match := reNationalFootballDaily.FindStringSubmatch(strings.TrimSpace(q))
	if len(match) != 2 {
		return false
	}
	team := strings.TrimSpace(match[1])
	_, ok := footballNations[team]
	return ok
}

func isNationalFootballAdvanceQuestion(q string) bool {
	match := reNationalFootballAdvance.FindStringSubmatch(strings.TrimSpace(q))
	if len(match) != 3 {
		return false
	}
	_, okA := footballNations[strings.TrimSpace(match[1])]
	_, okB := footballNations[strings.TrimSpace(match[2])]
	return okA && okB
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
// Used by older scanners. Keep narrow: only in-play daily / series markets.
func IsSportsMarket(m Market) bool {
	return IsLoLMarket(m) || IsBasketballMarket(m) || IsFootballMarket(m) || IsDota2Market(m) || IsTennisMarket(m)
}

// IsFollowTargetMarket is the stricter universe for smart-money whale following:
// basketball, soccer/football, and esports. Tennis and other sports stay out
// because they are not part of the current copy-trading mandate.
func IsFollowTargetMarket(m Market) bool {
	return IsLoLMarket(m) || IsBasketballMarket(m) || IsFootballMarket(m) || IsDota2Market(m)
}

// FilterLoL returns only LoL markets from a list.
func FilterLoL(ms []Market) []Market {
	return filterBy(ms, IsLoLMarket)
}

// FilterSports — LoL + NBA (daily+playoffs) + EPL daily.
func FilterSports(ms []Market) []Market {
	return filterBy(ms, IsSportsMarket)
}

func FilterFollowTargets(ms []Market) []Market {
	return filterBy(ms, IsFollowTargetMarket)
}

// FilterTradablePriceBand removes markets whose displayed outcome prices are
// all outside the configured entry band. Markets without parseable prices are
// kept so discovery does not drop valid targets when Gamma omits price fields.
func FilterTradablePriceBand(ms []Market, minPrice, maxPrice float64) []Market {
	if minPrice <= 0 && maxPrice >= 1 {
		return ms
	}
	out := make([]Market, 0, len(ms))
	for _, m := range ms {
		if marketHasTradableOutcomePrice(m, minPrice, maxPrice) {
			out = append(out, m)
		}
	}
	return out
}

func marketHasTradableOutcomePrice(m Market, minPrice, maxPrice float64) bool {
	prices := m.OutcomePrices()
	if len(prices) == 0 {
		return true
	}
	parseable := 0
	for _, raw := range prices {
		p, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			continue
		}
		parseable++
		if p >= minPrice && p <= maxPrice {
			return true
		}
	}
	return parseable == 0
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
	return c.getMarkets(ctx, q)
}

// GetByClobTokenIDs fetches markets containing the given CLOB token ids. Gamma
// supports repeated `clob_token_ids=<token>` params, which lets report tooling
// backfill conditionId/outcomePrices for older trade logs that only stored the
// token id.
func (c *GammaClient) GetByClobTokenIDs(ctx context.Context, ids []string) ([]Market, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := url.Values{}
	for _, id := range ids {
		if id == "" {
			continue
		}
		q.Add("clob_token_ids", id)
	}
	q.Set("limit", fmt.Sprintf("%d", len(ids)+5))
	return c.getMarkets(ctx, q)
}

func (c *GammaClient) getMarkets(ctx context.Context, q url.Values) ([]Market, error) {
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
