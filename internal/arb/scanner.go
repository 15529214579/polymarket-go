package arb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/odds"
)

// ScanConfig holds configurable thresholds for the arb scanner.
type ScanConfig struct {
	MinGapPP       float64 // minimum net EV in pp to flag (default 5)
	TradingCostPP  float64 // estimated round-trip cost in pp (default 2)
	MinBookCount   int     // consensus filter: min bookmakers (default 3)
	DeviationPP    float64 // consensus filter: max deviation from median (default 0.05)
	BaseBetUSDC    float64 // base bet size (default 2.5)
	MaxBetUSDC     float64 // max bet size (default 5.0)
	Tier2MinPP     float64 // net_ev threshold for tier 2 (default 10)
	Tier3MinPP     float64 // net_ev threshold for tier 3 (default 15)
	Tier2Mult      float64 // multiplier for tier 2 (default 1.5)
	Tier3Mult      float64 // multiplier for tier 3 (default 2.0)
}

func DefaultScanConfig() ScanConfig {
	return ScanConfig{
		MinGapPP:      5.0,
		TradingCostPP: 2.0,
		MinBookCount:  3,
		DeviationPP:   0.05,
		BaseBetUSDC:   2.5,
		MaxBetUSDC:    5.0,
		Tier2MinPP:    10.0,
		Tier3MinPP:    15.0,
		Tier2Mult:     1.5,
		Tier3Mult:     2.0,
	}
}

func (c ScanConfig) tieredBetSize(netEvPP float64) float64 {
	var size float64
	switch {
	case netEvPP >= c.Tier3MinPP:
		size = c.BaseBetUSDC * c.Tier3Mult
	case netEvPP >= c.Tier2MinPP:
		size = c.BaseBetUSDC * c.Tier2Mult
	default:
		size = c.BaseBetUSDC
	}
	if size > c.MaxBetUSDC {
		size = c.MaxBetUSDC
	}
	return size
}

// polyMarket is a minimal Gamma API market representation for arb matching.
type polyMarket struct {
	ConditionID string `json:"conditionId"`
	ID          string `json:"id"`
	Question    string `json:"question"`
	Slug        string `json:"slug"`
	EndDate     string `json:"endDate"`
	Active      bool   `json:"active"`
	Closed      bool   `json:"closed"`
	Tokens      json.RawMessage `json:"clobTokenIds"`
	Prices      json.RawMessage `json:"outcomePrices"`
	Outcomes    json.RawMessage `json:"outcomes"`
}

func (m polyMarket) tokenIDs() []string   { return parseJSONStringArray(m.Tokens) }
func (m polyMarket) outcomePrices() []string { return parseJSONStringArray(m.Prices) }
func (m polyMarket) outcomeList() []string { return parseJSONStringArray(m.Outcomes) }

func parseJSONStringArray(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	// Gamma returns these as JSON strings containing JSON arrays
	// (double-encoded), e.g. "[\"abc\",\"def\"]" as a string value.
	// Try direct array first, then unwrap string.
	var out []string
	if err := json.Unmarshal(raw, &out); err == nil {
		return out
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil && s != "" {
		var inner []string
		if err := json.Unmarshal([]byte(s), &inner); err == nil {
			return inner
		}
	}
	return nil
}

// extractTokenAndPrice returns (yesToken, noToken, yesPrice, noPrice, ok).
func (m polyMarket) extractTokenAndPrice() (string, string, float64, float64, bool) {
	outcomes := m.outcomeList()
	if len(outcomes) > 2 {
		return "", "", 0, 0, false
	}

	tokens := m.tokenIDs()
	if len(tokens) == 0 {
		return "", "", 0, 0, false
	}

	yesToken := ""
	noToken := ""
	if len(tokens) >= 1 {
		yesToken = tokens[0]
	}
	if len(tokens) >= 2 {
		noToken = tokens[1]
	}

	prices := m.outcomePrices()
	if len(prices) >= 2 {
		var yp, np float64
		fmt.Sscanf(prices[0], "%f", &yp)
		fmt.Sscanf(prices[1], "%f", &np)
		if yp > 0 || np > 0 {
			return yesToken, noToken, yp, np, true
		}
	}
	return "", "", 0, 0, false
}

// Scanner runs the arb scan pipeline.
type Scanner struct {
	oddsClient *odds.Client
	httpClient *http.Client
	gammaBase  string
	cfg        ScanConfig
	store      *Store
}

func NewScanner(oc *odds.Client, store *Store, cfg ScanConfig) *Scanner {
	return &Scanner{
		oddsClient: oc,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		gammaBase:  "https://gamma-api.polymarket.com",
		cfg:        cfg,
		store:      store,
	}
}

// Scan runs one full arb scan cycle: fetch PM markets, fetch bookmaker odds,
// match, compute gaps, filter, deduplicate, store snapshots.
func (s *Scanner) Scan(ctx context.Context) ([]odds.ArbOpportunity, error) {
	polyMarkets, err := s.fetchPolySportsMarkets(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch poly markets: %w", err)
	}
	if len(polyMarkets) == 0 {
		slog.Info("arb_scan_no_poly_markets")
		return nil, nil
	}

	// Fetch h2h odds + outrights in parallel.
	type fetchResult struct {
		items []odds.BookmakerOdds
		err   error
	}
	h2hCh := make(chan fetchResult, 1)
	outCh := make(chan fetchResult, 1)

	go func() {
		items, err := s.oddsClient.FetchH2H(ctx, nil)
		h2hCh <- fetchResult{items, err}
	}()
	go func() {
		items, err := s.oddsClient.FetchOutrights(ctx, nil)
		outCh <- fetchResult{items, err}
	}()

	h2hRes := <-h2hCh
	outRes := <-outCh
	if h2hRes.err != nil {
		slog.Warn("arb_h2h_fetch_fail", "err", h2hRes.err.Error())
	}
	if outRes.err != nil {
		slog.Warn("arb_outright_fetch_fail", "err", outRes.err.Error())
	}

	var allOdds []odds.BookmakerOdds
	allOdds = append(allOdds, h2hRes.items...)
	allOdds = append(allOdds, outRes.items...)
	if len(allOdds) == 0 {
		slog.Info("arb_scan_no_odds")
		return nil, nil
	}

	// Apply consensus filter.
	consensus := odds.ApplyConsensusFilter(allOdds, s.cfg.MinBookCount, s.cfg.DeviationPP)
	slog.Info("arb_consensus", "raw", len(allOdds), "consensus", len(consensus))
	if len(consensus) == 0 {
		return nil, nil
	}

	// Match and compute gaps.
	var opportunities []odds.ArbOpportunity
	totalMatched := 0

	for _, item := range consensus {
		market := s.matchToPolymarket(item, polyMarkets)
		if market == nil {
			continue
		}

		yesToken, noToken, yesPrice, noPrice, ok := market.extractTokenAndPrice()
		if !ok {
			continue
		}
		totalMatched++

		yesGap := (item.BookmakerProb - yesPrice) * 100
		noGap := ((1 - item.BookmakerProb) - noPrice) * 100

		var gap float64
		var direction, tokenID string
		var polyPrice float64

		if yesGap >= noGap {
			gap = yesGap
			direction = "BUY_YES"
			tokenID = yesToken
			polyPrice = yesPrice
		} else {
			gap = noGap
			direction = "BUY_NO"
			if noToken != "" {
				tokenID = noToken
			} else {
				tokenID = yesToken
			}
			polyPrice = noPrice
		}

		netEv := gap - s.cfg.TradingCostPP
		title := market.Question

		slog.Info("arb_match",
			"team", truncStr(item.TeamOrSide, 30),
			"market", truncStr(title, 50),
			"bk_prob", item.BookmakerProb,
			"poly", polyPrice,
			"gap_pp", math.Round(gap*10)/10,
			"net_ev_pp", math.Round(netEv*10)/10,
			"dir", direction,
		)

		if math.Abs(gap) > 45 {
			continue
		}
		if polyPrice < 0.02 || polyPrice > 0.95 {
			continue
		}

		// Store snapshot to DB regardless of threshold.
		if s.store != nil {
			s.store.Insert(odds.OddsSnapshot{
				MarketID:        tokenID,
				Sport:           item.Sport,
				EventName:       item.EventName + " — " + item.TeamOrSide,
				PolymarketPrice: polyPrice,
				Bookmaker:       item.Bookmaker,
				BookmakerProb:   item.BookmakerProb,
				GapPP:           math.Round(gap*100) / 100,
				SnapshotTimestamp: time.Now(),
			})
		}

		if netEv >= s.cfg.MinGapPP {
			opportunities = append(opportunities, odds.ArbOpportunity{
				TokenID:         tokenID,
				PolymarketPrice: polyPrice,
				BookmakerProb:   item.BookmakerProb,
				GapPP:           math.Round(gap*100) / 100,
				NetEvPP:         math.Round(netEv*100) / 100,
				Direction:       direction,
				MarketTitle:     truncStr(title, 100),
				Sport:           item.Sport,
				EventName:       item.EventName + " — " + item.TeamOrSide,
				Bookmaker:       item.Bookmaker,
				BetSizeUSDC:     s.cfg.tieredBetSize(netEv),
			})
		}
	}

	slog.Info("arb_scan_summary",
		"consensus", len(consensus),
		"matched", totalMatched,
		"opportunities", len(opportunities),
		"min_gap_pp", s.cfg.MinGapPP,
		"cost_pp", s.cfg.TradingCostPP,
	)

	// Deduplicate by tokenID, keep largest gap.
	seen := map[string]odds.ArbOpportunity{}
	for _, opp := range opportunities {
		if prev, ok := seen[opp.TokenID]; !ok || math.Abs(opp.GapPP) > math.Abs(prev.GapPP) {
			seen[opp.TokenID] = opp
		}
	}

	var final []odds.ArbOpportunity
	for _, opp := range seen {
		final = append(final, opp)
	}
	// Sort: priority tier desc, then gap desc.
	sortOpportunities(final)

	slog.Info("arb_scan_final", "opportunities", len(final))
	return final, nil
}

func sortOpportunities(opps []odds.ArbOpportunity) {
	for i := 0; i < len(opps); i++ {
		for j := i + 1; j < len(opps); j++ {
			pi := PriorityBoost(opps[i].Sport, opps[i].MarketTitle)
			pj := PriorityBoost(opps[j].Sport, opps[j].MarketTitle)
			if pj > pi || (pj == pi && math.Abs(opps[j].GapPP) > math.Abs(opps[i].GapPP)) {
				opps[i], opps[j] = opps[j], opps[i]
			}
		}
	}
}

// matchToPolymarket finds the best Polymarket market for a bookmaker odds item.
func (s *Scanner) matchToPolymarket(item odds.BookmakerOdds, markets []polyMarket) *polyMarket {
	team := normalizeTeamName(item.TeamOrSide)
	if team == "" || len(team) < 3 {
		return nil
	}

	sportKws := sportToPolyKeywords[item.Sport]
	isH2H := strings.HasPrefix(strings.ToLower(item.MarketName), "h2h")
	isDraw := isH2H && team == "draw"

	var drawTeamTokens []string
	if isDraw && item.EventName != "" {
		parts := regexp.MustCompile(`(?i)\s+vs\.?\s+`).Split(item.EventName, -1)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			norm := normalizeTeamName(p)
			words := strings.Fields(norm)
			if len(words) > 0 {
				last := words[len(words)-1]
				if len(last) >= 3 {
					drawTeamTokens = append(drawTeamTokens, last)
				}
			}
		}
	}

	var bestMatch *polyMarket
	bestScore := 0.0

	for i := range markets {
		m := &markets[i]
		question := strings.ToLower(m.Question)

		// Sport-context validation.
		if len(sportKws) > 0 {
			found := false
			for _, kw := range sportKws {
				if strings.Contains(question, kw) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Winner vs qualify mismatch.
		if strings.Contains(strings.ToLower(item.Sport), "winner") && strings.Contains(question, "qualify") {
			continue
		}

		// Seasonal filter for h2h.
		if isH2H {
			if seasonal, _ := isSeasonalMarket(m.Question); seasonal {
				continue
			}
		}

		var score float64
		if isDraw {
			if len(drawTeamTokens) > 0 {
				allFound := true
				for _, t := range drawTeamTokens {
					if !strings.Contains(question, t) {
						allFound = false
						break
					}
				}
				if allFound {
					score = 1.0
				} else {
					continue
				}
			} else if strings.Contains(question, "draw") {
				score = 0.5
			} else {
				continue
			}
		} else if strings.Contains(question, team) {
			score = 1.0
		} else {
			parts := strings.Fields(team)
			nickname := parts[len(parts)-1]
			if len(nickname) >= 4 && strings.Contains(question, nickname) {
				score = 0.8
			} else {
				continue
			}
		}

		// Cross-match guard for h2h soccer.
		if isH2H && strings.HasPrefix(item.Sport, "soccer_") {
			parts := regexp.MustCompile(`(?i)\s+vs\.?\s+`).Split(item.EventName, 2)
			home := ""
			away := ""
			if len(parts) >= 1 {
				home = strings.TrimSpace(parts[0])
			}
			if len(parts) >= 2 {
				away = strings.TrimSpace(parts[1])
			}
			ok, _ := crossMatchGuard(
				item.Sport,
				item.EventCommenceTime,
				home, away,
				m.Question,
				m.EndDate,
			)
			if !ok {
				continue
			}
		}

		if score > bestScore {
			bestScore = score
			bestMatch = m
		}
	}
	return bestMatch
}

// fetchPolySportsMarkets fetches active Polymarket markets and filters for sports.
func (s *Scanner) fetchPolySportsMarkets(ctx context.Context) ([]polyMarket, error) {
	var all []polyMarket

	// Pass 1: wide window (paginated, top by volume).
	for page := 0; page < 20; page++ {
		offset := page * 100
		params := url.Values{}
		params.Set("active", "true")
		params.Set("closed", "false")
		params.Set("limit", "100")
		params.Set("offset", fmt.Sprintf("%d", offset))

		items, err := s.fetchGammaPage(ctx, params)
		if err != nil {
			slog.Warn("gamma_batch_fail", "page", page, "err", err.Error())
			break
		}
		if len(items) == 0 {
			break
		}
		all = append(all, items...)
	}

	// Pass 2: narrow 96h window for game-level h2h markets.
	now := time.Now().UTC()
	narrowStart := now.Add(1 * time.Hour).Format("2006-01-02T15:04:05Z")
	narrowEnd := now.Add(96 * time.Hour).Format("2006-01-02T15:04:05Z")

	params := url.Values{}
	params.Set("active", "true")
	params.Set("closed", "false")
	params.Set("limit", "500")
	params.Set("end_date_min", narrowStart)
	params.Set("end_date_max", narrowEnd)

	narrow, err := s.fetchGammaPage(ctx, params)
	if err != nil {
		slog.Warn("gamma_narrow_fail", "err", err.Error())
	} else {
		all = append(all, narrow...)
	}

	// Deduplicate and filter.
	seen := map[string]bool{}
	var sports []polyMarket
	for _, m := range all {
		key := m.ConditionID
		if key == "" {
			key = m.ID
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		if IsPolySportsMarket(m.Question) {
			sports = append(sports, m)
		}
	}

	slog.Info("poly_sports_fetched", "total", len(all), "sports", len(sports))
	return sports, nil
}

func (s *Scanner) fetchGammaPage(ctx context.Context, params url.Values) ([]polyMarket, error) {
	u := s.gammaBase + "/markets?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gamma %d: %s", resp.StatusCode, truncStr(string(body), 200))
	}
	var markets []polyMarket
	if err := json.Unmarshal(body, &markets); err != nil {
		return nil, fmt.Errorf("gamma decode: %w", err)
	}
	return markets, nil
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
