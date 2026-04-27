// Package injury — NBA injury report scanner.
// Guarded by -injury_enabled flag; all exports are safe to call when disabled.
// To remove: delete this package + grep "injury" in cmd/bot/main.go.
package injury

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type PlayerStatus string

const (
	StatusOut      PlayerStatus = "Out"
	StatusDoubtful PlayerStatus = "Doubtful"
	StatusQuest    PlayerStatus = "Questionable"
	StatusProb     PlayerStatus = "Probable"
	StatusAvail    PlayerStatus = "Available"
)

type InjuryEntry struct {
	Player string       `json:"player"`
	Team   string       `json:"team"`
	Status PlayerStatus `json:"status"`
	Reason string       `json:"reason"`
}

type InjuryAlert struct {
	Team       string        `json:"team"`
	Opponent   string        `json:"opponent"`
	GameDate   string        `json:"game_date"`
	StarPlayer string        `json:"star_player"`
	Status     PlayerStatus  `json:"status"`
	Reason     string        `json:"reason"`
	Impact     string        `json:"impact"`
	Entries    []InjuryEntry `json:"entries"`
	FetchedAt  time.Time     `json:"fetched_at"`
}

type Config struct {
	Enabled      bool
	ScanInterval time.Duration
	StarOnly     bool // only alert on star players (top ~15 per team)
}

func DefaultConfig() Config {
	return Config{
		Enabled:      false,
		ScanInterval: 30 * time.Minute,
		StarOnly:     true,
	}
}

var nbaStars = map[string][]string{
	"Atlanta Hawks":          {"Jalen Johnson", "De'Andre Hunter", "Clint Capela"},
	"Boston Celtics":         {"Jayson Tatum", "Jaylen Brown", "Derrick White", "Kristaps Porzingis"},
	"Brooklyn Nets":          {"Mikal Bridges", "Cameron Johnson", "Michael Porter Jr."},
	"Charlotte Hornets":      {"LaMelo Ball", "Brandon Miller", "Miles Bridges"},
	"Chicago Bulls":          {"Anfernee Simons", "Josh Giddey", "Coby White", "Rob Dillingham"},
	"Cleveland Cavaliers":    {"Donovan Mitchell", "Darius Garland", "Evan Mobley", "Jarrett Allen", "James Harden"},
	"Dallas Mavericks":       {"Kyrie Irving", "P.J. Washington", "Caleb Martin", "Cooper Flagg"},
	"Denver Nuggets":         {"Nikola Jokic", "Jamal Murray", "Aaron Gordon"},
	"Detroit Pistons":        {"Cade Cunningham", "Jaden Ivey", "Ausar Thompson"},
	"Golden State Warriors":  {"Stephen Curry", "Jimmy Butler III", "Draymond Green", "Andrew Wiggins"},
	"Houston Rockets":        {"Kevin Durant", "Fred VanVleet", "Alperen Sengun", "Jabari Smith Jr.", "Amen Thompson", "Reed Sheppard"},
	"Indiana Pacers":         {"Tyrese Haliburton", "Pascal Siakam", "Myles Turner", "Ivica Zubac"},
	"LA Clippers":            {"Kawhi Leonard", "Bradley Beal", "Norman Powell"},
	"Los Angeles Lakers":     {"LeBron James", "Luka Doncic", "Austin Reaves", "Deandre Ayton"},
	"Memphis Grizzlies":      {"Ja Morant", "Desmond Bane", "Zach Edey"},
	"Miami Heat":             {"Bam Adebayo", "Tyler Herro"},
	"Milwaukee Bucks":        {"Giannis Antetokounmpo", "Kyle Kuzma", "Bobby Portis"},
	"Minnesota Timberwolves": {"Anthony Edwards", "Julius Randle", "Rudy Gobert", "Donte DiVincenzo"},
	"New Orleans Pelicans":   {"Zion Williamson", "Dejounte Murray", "CJ McCollum"},
	"New York Knicks":        {"Jalen Brunson", "OG Anunoby", "Karl-Anthony Towns"},
	"Oklahoma City Thunder":  {"Shai Gilgeous-Alexander", "Jalen Williams", "Chet Holmgren"},
	"Orlando Magic":          {"Paolo Banchero", "Franz Wagner", "Jalen Suggs"},
	"Philadelphia 76ers":     {"Joel Embiid", "Tyrese Maxey", "Paul George"},
	"Phoenix Suns":           {"Devin Booker", "Mark Williams"},
	"Portland Trail Blazers": {"Damian Lillard", "Scoot Henderson", "Jerami Grant"},
	"Sacramento Kings":       {"DeMar DeRozan", "Domantas Sabonis", "Zach LaVine", "Keegan Murray"},
	"San Antonio Spurs":      {"Victor Wembanyama", "Keldon Johnson", "Tre Jones"},
	"Toronto Raptors":        {"Scottie Barnes", "RJ Barrett", "Immanuel Quickley"},
	"Utah Jazz":              {"Lauri Markkanen", "Collin Sexton", "Jaren Jackson Jr."},
	"Washington Wizards":     {"Anthony Davis", "Trae Young", "Alex Sarr", "Bilal Coulibaly"},
}

func isStar(team, player string) bool {
	stars, ok := nbaStars[team]
	if !ok {
		return false
	}
	lp := strings.ToLower(player)
	for _, s := range stars {
		if strings.ToLower(s) == lp {
			return true
		}
	}
	return false
}

type GameInfo struct {
	HomeTeam   string
	AwayTeam   string
	Tipoff     time.Time
	Status     string // "Scheduled" / "In Progress" / "Final"
	SeriesNote string // "Game 5 · Series tied 2-2" etc.
}

type Scanner struct {
	cfg      Config
	client   *http.Client
	seen     map[string]time.Time // "2006-01-02:team:player:status" → alert time
	seenPath string               // persistence file path

	mu        sync.RWMutex
	cache     map[string][]InjuryEntry // team → current star injuries (refreshed each Scan)
	allCache  map[string][]InjuryEntry // team → ALL injuries (stars + non-stars, OUT/Doubtful/Questionable)
	games     map[string]GameInfo      // team name → today's game (both teams point to same GameInfo)
}

func NewScanner(cfg Config, dbDir string) *Scanner {
	s := &Scanner{
		cfg:      cfg,
		client:   &http.Client{Timeout: 15 * time.Second},
		seen:     make(map[string]time.Time),
		seenPath: filepath.Join(dbDir, "injury_seen.json"),
		cache:    make(map[string][]InjuryEntry),
		allCache: make(map[string][]InjuryEntry),
		games:    make(map[string]GameInfo),
	}
	s.loadSeen()
	return s
}

func (s *Scanner) seenKey(team, player string, status PlayerStatus) string {
	date := time.Now().Format("2006-01-02")
	return date + ":" + team + ":" + player + ":" + string(status)
}

func (s *Scanner) loadSeen() {
	data, err := os.ReadFile(s.seenPath)
	if err != nil {
		return
	}
	var m map[string]time.Time
	if json.Unmarshal(data, &m) == nil {
		s.seen = m
	}
}

func (s *Scanner) saveSeen() {
	data, err := json.Marshal(s.seen)
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(s.seenPath), 0755)
	_ = os.WriteFile(s.seenPath, data, 0644)
}

// InjuredStars returns OUT/Doubtful star players for the given team.
// Safe to call from any goroutine; returns nil if no injuries or team unknown.
func (s *Scanner) InjuredStars(team string) []InjuryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.cache[team]; ok {
		return v
	}
	lt := strings.ToLower(team)
	for k, v := range s.cache {
		if strings.Contains(strings.ToLower(k), lt) || strings.Contains(lt, strings.ToLower(k)) {
			return v
		}
	}
	return nil
}

// HasInjuredStar reports whether the team has at least one star OUT or Doubtful.
func (s *Scanner) HasInjuredStar(team string) bool {
	return len(s.InjuredStars(team)) > 0
}

// AllInjuries returns all OUT/Doubtful/Questionable players (stars and non-stars) for the given team.
func (s *Scanner) AllInjuries(team string) []InjuryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.allCache[team]; ok {
		return v
	}
	lt := strings.ToLower(team)
	for k, v := range s.allCache {
		if strings.Contains(strings.ToLower(k), lt) || strings.Contains(lt, strings.ToLower(k)) {
			return v
		}
	}
	return nil
}

func (s *Scanner) Enabled() bool { return s.cfg.Enabled }

// Scan fetches the latest NBA injury data and returns alerts for
// significant status changes (star player OUT/Doubtful).
func (s *Scanner) Scan(ctx context.Context) ([]InjuryAlert, error) {
	if !s.cfg.Enabled {
		return nil, nil
	}

	entries, err := s.fetchInjuries(ctx)
	if err != nil {
		return nil, fmt.Errorf("injury fetch: %w", err)
	}

	if games, err := s.fetchScoreboard(ctx); err == nil {
		s.mu.Lock()
		s.games = games
		s.mu.Unlock()
	}

	now := time.Now()
	var alerts []InjuryAlert

	byTeam := make(map[string][]InjuryEntry)
	for _, e := range entries {
		byTeam[e.Team] = append(byTeam[e.Team], e)
	}

	// Rebuild the shared injury cache (read by momentum/lottery filters).
	starOut := make(map[string][]InjuryEntry)
	allInj := make(map[string][]InjuryEntry)
	for team, teamEntries := range byTeam {
		for _, e := range teamEntries {
			if e.Status == StatusOut || e.Status == StatusDoubtful || e.Status == StatusQuest {
				allInj[team] = append(allInj[team], e)
			}
			if (e.Status == StatusOut || e.Status == StatusDoubtful) && isStar(team, e.Player) {
				starOut[team] = append(starOut[team], e)
			}
		}
	}
	s.mu.Lock()
	s.cache = starOut
	s.allCache = allInj
	s.mu.Unlock()

	for team, teamEntries := range byTeam {
		for _, e := range teamEntries {
			if e.Status != StatusOut && e.Status != StatusDoubtful {
				continue
			}
			if s.cfg.StarOnly && !isStar(team, e.Player) {
				continue
			}

			key := s.seenKey(team, e.Player, e.Status)
			if _, ok := s.seen[key]; ok {
				continue // already pushed today for this player+status
			}

			impact := assessImpact(team, e)
			alert := InjuryAlert{
				Team:       team,
				StarPlayer: e.Player,
				Status:     e.Status,
				Reason:     e.Reason,
				Impact:     impact,
				Entries:    teamEntries,
				FetchedAt:  now,
			}
			alerts = append(alerts, alert)
			s.seen[key] = now
		}
	}

	// prune entries older than 48h
	for k, t := range s.seen {
		if now.Sub(t) > 48*time.Hour {
			delete(s.seen, k)
		}
	}

	if len(alerts) > 0 {
		s.saveSeen()
	}

	return alerts, nil
}

func assessImpact(team string, e InjuryEntry) string {
	stars := nbaStars[team]
	if len(stars) == 0 {
		return "unknown_team"
	}
	lp := strings.ToLower(e.Player)
	for i, s := range stars {
		if strings.ToLower(s) == lp {
			if i == 0 {
				return "franchise_player_out"
			}
			if i <= 1 {
				return "co_star_out"
			}
			return "rotation_star_out"
		}
	}
	return "role_player_out"
}

// PlayerRole returns a human-readable role label for a player on a team.
func PlayerRole(team, player string) string {
	stars := nbaStars[team]
	lp := strings.ToLower(player)
	for i, s := range stars {
		if strings.ToLower(s) == lp {
			switch i {
			case 0:
				return "核心/当家球星"
			case 1:
				return "二当家"
			case 2:
				return "第三核心"
			default:
				return "主力轮换"
			}
		}
	}
	return "角色球员"
}

// PlayerImpactPct estimates a player's contribution to team strength (0-100).
// Based on position in the nbaStars list: franchise ~35%, co-star ~25%, 3rd ~15%, rotation ~10%, role ~5%.
func PlayerImpactPct(team, player string) int {
	stars := nbaStars[team]
	lp := strings.ToLower(player)
	for i, s := range stars {
		if strings.ToLower(s) == lp {
			switch i {
			case 0:
				return 35
			case 1:
				return 25
			case 2:
				return 15
			default:
				return 10
			}
		}
	}
	return 5
}

type espnResponse struct {
	Injuries []struct {
		DisplayName string `json:"displayName"`
		Injuries    []struct {
			Athlete struct {
				DisplayName string `json:"displayName"`
			} `json:"athlete"`
			Status  string `json:"status"`
			Details struct {
				Detail string `json:"detail"`
			} `json:"details"`
		} `json:"injuries"`
	} `json:"injuries"`
}

func (s *Scanner) fetchInjuries(ctx context.Context) ([]InjuryEntry, error) {
	url := "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/injuries"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "polymarket-go/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("espn status %d: %s", resp.StatusCode, body)
	}

	var data espnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("espn decode: %w", err)
	}

	var entries []InjuryEntry
	for _, team := range data.Injuries {
		for _, inj := range team.Injuries {
			entries = append(entries, InjuryEntry{
				Player: inj.Athlete.DisplayName,
				Team:   team.DisplayName,
				Status: parseStatus(inj.Status),
				Reason: inj.Details.Detail,
			})
		}
	}

	slog.Info("injury_fetch", "entries", len(entries), "teams", len(data.Injuries))
	return entries, nil
}

func parseStatus(s string) PlayerStatus {
	ls := strings.ToLower(s)
	switch {
	case strings.Contains(ls, "out"):
		return StatusOut
	case strings.Contains(ls, "doubtful"):
		return StatusDoubtful
	case strings.Contains(ls, "questionable"):
		return StatusQuest
	case strings.Contains(ls, "probable"):
		return StatusProb
	default:
		return StatusAvail
	}
}

func (s *Scanner) GameFor(team string) (GameInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if g, ok := s.games[team]; ok {
		return g, true
	}
	lt := strings.ToLower(team)
	for k, g := range s.games {
		if strings.Contains(strings.ToLower(k), lt) || strings.Contains(lt, strings.ToLower(k)) {
			return g, true
		}
	}
	return GameInfo{}, false
}

type espnScoreboard struct {
	Events []struct {
		Date         string `json:"date"`
		Competitions []struct {
			Status struct {
				Type struct {
					Description string `json:"description"`
				} `json:"type"`
			} `json:"status"`
			Competitors []struct {
				HomeAway string `json:"homeAway"`
				Team     struct {
					DisplayName string `json:"displayName"`
				} `json:"team"`
			} `json:"competitors"`
			Series struct {
				Summary string `json:"summary"`
				Title   string `json:"title"`
			} `json:"series"`
		} `json:"competitions"`
	} `json:"events"`
}

func (s *Scanner) fetchScoreboard(ctx context.Context) (map[string]GameInfo, error) {
	url := "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "polymarket-go/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("espn scoreboard %d: %s", resp.StatusCode, body)
	}

	var data espnScoreboard
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("espn scoreboard decode: %w", err)
	}

	games := make(map[string]GameInfo)
	for _, ev := range data.Events {
		if len(ev.Competitions) == 0 {
			continue
		}
		comp := ev.Competitions[0]
		var home, away string
		for _, c := range comp.Competitors {
			if c.HomeAway == "home" {
				home = c.Team.DisplayName
			} else {
				away = c.Team.DisplayName
			}
		}
		tipoff, _ := time.Parse(time.RFC3339, ev.Date)

		seriesNote := ""
		if comp.Series.Title != "" {
			seriesNote = comp.Series.Title
			if comp.Series.Summary != "" {
				seriesNote += " · " + comp.Series.Summary
			}
		}

		gi := GameInfo{
			HomeTeam:   home,
			AwayTeam:   away,
			Tipoff:     tipoff,
			Status:     comp.Status.Type.Description,
			SeriesNote: seriesNote,
		}
		games[home] = gi
		games[away] = gi
	}

	slog.Info("scoreboard_fetch", "games", len(data.Events))
	return games, nil
}
