// Package arb matches bookmaker odds against Polymarket markets and
// identifies cross-venue price gaps.
package arb

import (
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"
)

// --- Non-sports blacklist ---

var nonSportsBlacklist = []string{
	"president", "election", "ceasefire", "war", "peace", "treaty",
	"congress", "senate", "parliament", "democrat", "republican",
	"tariff", "inflation", "gdp", "recession", "fed ", "interest rate",
	"jesus", "christ", "pope", "god", "bible", "religion", "church",
	"gta", "movie", "album", "oscar", "grammy", "emmy", "netflix",
	"spacex", "mars", "moon", "nasa", "launch", "asteroid",
	"bitcoin", "ethereum", "crypto", "token", "blockchain",
	"ai ", "artificial intelligence", "openai", "chatgpt",
	"covid", "pandemic", "vaccine", "virus",
	"earthquake", "hurricane", "flood", "climate",
	"trump", "biden", "putin", "zelensky", "xi jinping",
	"tiktok", "twitter", "meta ", "apple ", "google ",
	"nuclear", "missile", "sanctions",
	"assassination", "impeach", "resign",
}

// --- Sports keywords ---

var sportsKeywords = []string{
	"vs", "nba", "nfl", "nhl", "mlb", "epl", "premier league", "la liga",
	"bundesliga", "serie a", "ligue 1", "champions league", "laliga",
	"world cup", "super bowl", "stanley cup", "march madness",
	"basketball", "football", "soccer", "baseball", "hockey",
	"tennis", "golf", "boxing", "ufc", "mma",
	"lol", "dota", "valorant", "esport", "worlds",
	"lck", "lpl", "lec", "lcs", "msi", "vct",
	"cs2", "csgo", "bellator", "ncaa",
}

var polySportsKeywords = []string{
	"nba", "nfl", "nhl", "mlb", "world cup", "fifa", "champion", "finals",
	"qualify", "super bowl", "stanley cup", "premier league", "epl",
	"la liga", "bundesliga", "serie a", "ligue 1", "champions league",
	"soccer", "football", "basketball", "baseball", "hockey",
	"tennis", "golf", "boxing", "ufc", "mma", "cricket", "rugby",
	"masters", "wimbledon", "pga", "ncaa", "march madness",
	"lol", "valorant", "cs2", "csgo", "dota", "esport", "worlds",
	"lck", "lpl", "lec", "vct",
	" vs. ", " vs ",
}

// --- Seasonal market filter ---

var seasonalKeywords = []string{
	"finish in top", "be relegated", "win the league", "win the championship",
	"win the cup", "win the champions league", "win the europa league",
	"top 4", "top scorer", "top goal scorer", "golden boot",
	"make the semi", "make the final", "make the quarter",
	"relegated", "promoted",
	"nba finals", "stanley cup", "super bowl", "world series",
	"conference final", "division champion",
	"mvp", "rookie of the year", "defensive player of the year",
	"coach of the year", "fighter of the year", "ballon d",
}

// Go regexp2 doesn't support negative lookahead (?!), so we match the
// base pattern and reject ISO-date continuations (e.g. "2026-04-19") in code.
var seasonRegex = regexp.MustCompile(`\b20\d{2}[-\x{2013}/]\d{2}\b`)
var seasonYearPrefix = regexp.MustCompile(`(?i)\bwin\s+the\s+20\d{2}\b`)

func isSeasonalMarket(title string) (bool, string) {
	lower := strings.ToLower(title)
	for _, kw := range seasonalKeywords {
		if strings.Contains(lower, kw) {
			return true, kw
		}
	}
	if locs := seasonRegex.FindStringIndex(title); locs != nil {
		end := locs[1]
		// Reject if followed by "-\d" (ISO date like 2026-04-19).
		if end < len(title)-1 && title[end] == '-' && end+1 < len(title) && title[end+1] >= '0' && title[end+1] <= '9' {
			// looks like an ISO date continuation, skip
		} else {
			return true, title[locs[0]:locs[1]]
		}
	}
	if m := seasonYearPrefix.FindString(title); m != "" {
		return true, m
	}
	return false, ""
}

// --- League identification ---

type leagueRule struct {
	pattern *regexp.Regexp
	key     string
}

var polyLeagueMap = []leagueRule{
	{regexp.MustCompile(`(?i)\bpremier league\b|\bepl\b`), "soccer_epl"},
	{regexp.MustCompile(`(?i)\befl championship\b|\bchampionship\b`), "soccer_efl_champ"},
	{regexp.MustCompile(`(?i)\bla liga\b|\blaliga\b`), "soccer_spain_la_liga"},
	{regexp.MustCompile(`(?i)\bserie a\b`), "soccer_italy_serie_a"},
	{regexp.MustCompile(`(?i)\bbundesliga\b`), "soccer_germany_bundesliga"},
	{regexp.MustCompile(`(?i)\bligue 1\b`), "soccer_france_ligue_1"},
	{regexp.MustCompile(`(?i)\bk league\b|\bkleague\b`), "soccer_korea_kleague1"},
	{regexp.MustCompile(`(?i)\bj league\b|\bjleague\b|\bj1 league\b`), "soccer_japan_j_league"},
	{regexp.MustCompile(`(?i)\bmls\b|\bmajor league soccer\b`), "soccer_usa_mls"},
}

func identifyPolyLeague(title string) string {
	for _, r := range polyLeagueMap {
		if r.pattern.MatchString(title) {
			return r.key
		}
	}
	return ""
}

// --- Team name Jaccard ---

var teamStripRe = regexp.MustCompile(`(?i)\b(fc|cf|sc|afc|club|united|city|albion|rovers|wanderers|athletic|town|county|hotspur|palace|united|villa|west brom|the)\b`)
var punctRe = regexp.MustCompile(`[^a-z0-9\s]`)

func jaccardTokens(text string) map[string]struct{} {
	text = strings.ToLower(text)
	text = teamStripRe.ReplaceAllString(text, " ")
	text = punctRe.ReplaceAllString(text, " ")
	tokens := map[string]struct{}{}
	for _, t := range strings.Fields(text) {
		if len(t) >= 3 {
			tokens[t] = struct{}{}
		}
	}
	return tokens
}

// crossMatchGuard applies the 3-layer guard from TASK-088.
// Returns (true, "ok") if all pass, or (false, reason).
func crossMatchGuard(
	sportKey string,
	commenceTimeISO string,
	homeTeam string,
	awayTeam string,
	polyTitle string,
	polyEndDateISO string,
) (bool, string) {
	// Guard 1: league strict equality
	polyLeague := identifyPolyLeague(polyTitle)
	if polyLeague == "" {
		return false, "league_unrecognised"
	}
	if sportKey != polyLeague {
		return false, "league_mismatch"
	}

	// Guard 2: date proximity ≤ 72h
	if commenceTimeISO != "" && polyEndDateISO != "" {
		tBk, err1 := parseISO(commenceTimeISO)
		tPoly, err2 := parseISO(polyEndDateISO)
		if err1 != nil || err2 != nil {
			return false, "date_far"
		}
		diffHours := math.Abs(tBk.Sub(tPoly).Hours())
		if diffHours > 72 {
			slog.Debug("cross_match_date_far", "diff_h", diffHours, "poly", polyTitle[:min(80, len(polyTitle))])
			return false, "date_far"
		}
	}

	// Guard 3: team-name recall-based Jaccard ≥ 0.6
	bkTokens := jaccardTokens(homeTeam + " " + awayTeam)
	polyTokens := jaccardTokens(polyTitle)
	if len(bkTokens) == 0 {
		return false, "name_jaccard_low"
	}
	matched := 0
	for t := range bkTokens {
		if _, ok := polyTokens[t]; ok {
			matched++
		}
	}
	recall := float64(matched) / float64(len(bkTokens))
	if recall < 0.6 || matched < 2 {
		return false, "name_jaccard_low"
	}
	return true, "ok"
}

func parseISO(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "Z") {
		s = s[:len(s)-1] + "+00:00"
	}
	for _, layout := range []string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable: %s", s)
}

// normalizeTeamName lowercases and strips common soccer-only suffixes.
// Does NOT strip "city" or "united" when they're part of the team
// identity (e.g. "Oklahoma City Thunder", "Manchester United").
func normalizeTeamName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, suffix := range []string{" fc", " cf", " sc"} {
		if strings.HasSuffix(name, suffix) {
			name = strings.TrimSpace(name[:len(name)-len(suffix)])
		}
	}
	return name
}

// IsPolySportsMarket checks if a Polymarket question matches any sports keyword.
func IsPolySportsMarket(question string) bool {
	q := strings.ToLower(question)
	for _, bl := range nonSportsBlacklist {
		if strings.Contains(q, bl) {
			return false
		}
	}
	for _, kw := range polySportsKeywords {
		if strings.Contains(q, kw) {
			return true
		}
	}
	return false
}

// sportToPolyKeywords maps bookmaker sport keys to required Polymarket title keywords.
var sportToPolyKeywords = map[string][]string{
	"basketball_nba_championship_winner": {"nba", "nba finals"},
	"americanfootball_nfl_super_bowl_winner": {"nfl", "super bowl"},
	"icehockey_nhl_championship_winner":     {"nhl", "stanley cup"},
	"soccer_fifa_world_cup_winner":          {"world cup", "fifa"},
	"baseball_mlb_world_series_winner":      {"mlb", "world series"},
	"golf_masters_tournament_winner":        {"masters"},
	"golf_pga_championship_winner":          {"pga"},
	"esports_lol":                           {"lol", "league of legends", "lck", "lpl", "lec", "worlds"},
	"esports_csgo":                          {"csgo", "cs2", "counter-strike", "esl", "blast"},
	"esports_valorant":                      {"valorant", "vct"},
}

// highPriorityKeywords for short-cycle sports.
var highPriorityKeywords = []string{
	"lol", "league of legends", "lck", "lpl", "lec", "lcs", "msi", "worlds",
	"valorant", "vct", "cs2", "csgo", "dota", "esport", "esports",
	"basketball", "nba",
}

// PriorityBoost returns 2 for short-cycle sports, 1 otherwise.
func PriorityBoost(sport, title string) int {
	combined := strings.ToLower(sport + " " + title)
	for _, kw := range highPriorityKeywords {
		if strings.Contains(combined, kw) {
			return 2
		}
	}
	return 1
}

