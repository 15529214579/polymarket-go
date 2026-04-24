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
	"strings"
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
	"Atlanta Hawks":          {"Trae Young", "Dejounte Murray", "Jalen Johnson"},
	"Boston Celtics":         {"Jayson Tatum", "Jaylen Brown", "Derrick White", "Kristaps Porzingis"},
	"Brooklyn Nets":          {"Mikal Bridges", "Cameron Johnson", "Ben Simmons"},
	"Charlotte Hornets":      {"LaMelo Ball", "Brandon Miller", "Miles Bridges"},
	"Chicago Bulls":          {"Zach LaVine", "DeMar DeRozan", "Coby White"},
	"Cleveland Cavaliers":    {"Donovan Mitchell", "Darius Garland", "Evan Mobley", "Jarrett Allen"},
	"Dallas Mavericks":       {"Luka Doncic", "Kyrie Irving", "Daniel Gafford"},
	"Denver Nuggets":         {"Nikola Jokic", "Jamal Murray", "Michael Porter Jr.", "Aaron Gordon"},
	"Detroit Pistons":        {"Cade Cunningham", "Jaden Ivey", "Ausar Thompson"},
	"Golden State Warriors":  {"Stephen Curry", "Klay Thompson", "Draymond Green", "Andrew Wiggins"},
	"Houston Rockets":        {"Jalen Green", "Alperen Sengun", "Jabari Smith Jr."},
	"Indiana Pacers":         {"Tyrese Haliburton", "Myles Turner", "Pascal Siakam"},
	"LA Clippers":            {"Kawhi Leonard", "Paul George", "James Harden"},
	"Los Angeles Lakers":     {"LeBron James", "Anthony Davis", "Austin Reaves"},
	"Memphis Grizzlies":      {"Ja Morant", "Desmond Bane", "Marcus Smart"},
	"Miami Heat":             {"Jimmy Butler", "Bam Adebayo", "Tyler Herro"},
	"Milwaukee Bucks":        {"Giannis Antetokounmpo", "Damian Lillard", "Khris Middleton"},
	"Minnesota Timberwolves": {"Anthony Edwards", "Karl-Anthony Towns", "Rudy Gobert"},
	"New Orleans Pelicans":   {"Zion Williamson", "Brandon Ingram", "CJ McCollum"},
	"New York Knicks":        {"Jalen Brunson", "Julius Randle", "OG Anunoby", "Donte DiVincenzo"},
	"Oklahoma City Thunder":  {"Shai Gilgeous-Alexander", "Jalen Williams", "Chet Holmgren"},
	"Orlando Magic":          {"Paolo Banchero", "Franz Wagner", "Jalen Suggs"},
	"Philadelphia 76ers":     {"Joel Embiid", "Tyrese Maxey"},
	"Phoenix Suns":           {"Kevin Durant", "Devin Booker", "Bradley Beal"},
	"Portland Trail Blazers": {"Anfernee Simons", "Scoot Henderson", "Jerami Grant"},
	"Sacramento Kings":       {"De'Aaron Fox", "Domantas Sabonis", "Keegan Murray"},
	"San Antonio Spurs":      {"Victor Wembanyama", "Keldon Johnson", "Tre Jones"},
	"Toronto Raptors":        {"Scottie Barnes", "RJ Barrett", "Immanuel Quickley"},
	"Utah Jazz":              {"Lauri Markkanen", "Collin Sexton", "Jordan Clarkson"},
	"Washington Wizards":     {"Kyle Kuzma", "Jordan Poole", "Deni Avdija"},
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

type Scanner struct {
	cfg    Config
	client *http.Client
	seen   map[string]time.Time // "team:player" → last alert time
}

func NewScanner(cfg Config) *Scanner {
	return &Scanner{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
		seen:   make(map[string]time.Time),
	}
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

	now := time.Now()
	var alerts []InjuryAlert

	byTeam := make(map[string][]InjuryEntry)
	for _, e := range entries {
		byTeam[e.Team] = append(byTeam[e.Team], e)
	}

	for team, teamEntries := range byTeam {
		for _, e := range teamEntries {
			if e.Status != StatusOut && e.Status != StatusDoubtful {
				continue
			}
			if s.cfg.StarOnly && !isStar(team, e.Player) {
				continue
			}

			key := team + ":" + e.Player
			if last, ok := s.seen[key]; ok && now.Sub(last) < 6*time.Hour {
				continue
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

	// prune old seen entries
	for k, t := range s.seen {
		if now.Sub(t) > 24*time.Hour {
			delete(s.seen, k)
		}
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

type espnResponse struct {
	Injuries []struct {
		Team struct {
			DisplayName string `json:"displayName"`
		} `json:"team"`
		Injuries []struct {
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
				Team:   team.Team.DisplayName,
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
