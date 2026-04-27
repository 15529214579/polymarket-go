// Package eurovision provides Eurovision PM market scanning and betting odds comparison.
package eurovision

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Market struct {
	Slug     string
	Question string
	Country  string
	YesPrice float64
	NoPrice  float64
	CondID   string
	Volume24 float64
}

type OddsEntry struct {
	Country     string
	ImpliedProb float64 // from betting odds
	BookOdds    float64 // decimal odds
	Source      string
}

type Signal struct {
	Market    Market
	BookProb  float64 // bookmaker implied probability
	PMPrice   float64 // PM Yes price
	Edge      float64 // BookProb - PMPrice
	Side      string  // "YES" or "NO"
	Source    string
}

// FetchEurovisionMarkets gets active Eurovision markets from PM gamma.
func FetchEurovisionMarkets(ctx context.Context) ([]Market, error) {
	url := "https://gamma-api.polymarket.com/markets?limit=200&active=true&closed=false&order=volume24hr&ascending=false"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch eurovision: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("gamma HTTP %d: %s", resp.StatusCode, body)
	}

	var raw []struct {
		Slug          string  `json:"slug"`
		Question      string  `json:"question"`
		OutcomePrices string  `json:"outcomePrices"`
		ConditionID   string  `json:"conditionId"`
		Volume24hr    float64 `json:"volume24hr"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var out []Market
	for _, m := range raw {
		q := strings.ToLower(m.Question)
		if !strings.Contains(q, "eurovision") {
			continue
		}
		country := parseCountry(m.Question)
		yes, no := parsePrices(m.OutcomePrices)
		out = append(out, Market{
			Slug:     m.Slug,
			Question: m.Question,
			Country:  country,
			YesPrice: yes,
			NoPrice:  no,
			CondID:   m.ConditionID,
			Volume24: m.Volume24hr,
		})
	}
	return out, nil
}

func parseCountry(q string) string {
	// "Will Finland win Eurovision 2026?" → "Finland"
	q = strings.TrimPrefix(q, "Will ")
	q = strings.TrimPrefix(q, "will ")
	parts := strings.SplitN(q, " win ", 2)
	if len(parts) >= 1 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

func parsePrices(raw string) (yes, no float64) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		raw = strings.ReplaceAll(raw, `\"`, `"`)
	}
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return
	}
	if len(arr) >= 2 {
		yes, _ = strconv.ParseFloat(arr[0], 64)
		no, _ = strconv.ParseFloat(arr[1], 64)
	}
	return
}

// FetchEurovisionOdds fetches Eurovision winner odds from the-odds-api.
func FetchEurovisionOdds(ctx context.Context, apiKey string) ([]OddsEntry, error) {
	// The Odds API Eurovision sport key
	url := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/entertainment_eurovision/odds?apiKey=%s&regions=eu&markets=outrights&oddsFormat=decimal",
		apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("odds api: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("odds api HTTP %d: %s", resp.StatusCode, body)
	}

	var events []struct {
		Bookmakers []struct {
			Key     string `json:"key"`
			Title   string `json:"title"`
			Markets []struct {
				Key      string `json:"key"`
				Outcomes []struct {
					Name  string  `json:"name"`
					Price float64 `json:"price"`
				} `json:"outcomes"`
			} `json:"markets"`
		} `json:"bookmakers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode odds: %w", err)
	}

	var entries []OddsEntry
	for _, ev := range events {
		for _, bk := range ev.Bookmakers {
			for _, mkt := range bk.Markets {
				if mkt.Key != "outrights" {
					continue
				}
				for _, o := range mkt.Outcomes {
					implied := 0.0
					if o.Price > 0 {
						implied = 1.0 / o.Price
					}
					entries = append(entries, OddsEntry{
						Country:     o.Name,
						ImpliedProb: implied,
						BookOdds:    o.Price,
						Source:      bk.Title,
					})
				}
			}
		}
	}
	return entries, nil
}

// ConsensusOdds averages implied probabilities across bookmakers for each country.
func ConsensusOdds(entries []OddsEntry) map[string]float64 {
	sums := make(map[string]float64)
	counts := make(map[string]int)
	for _, e := range entries {
		sums[e.Country] += e.ImpliedProb
		counts[e.Country]++
	}
	out := make(map[string]float64)
	for k, s := range sums {
		out[k] = s / float64(counts[k])
	}
	return out
}

// EvalSignals compares PM prices vs bookmaker consensus to find edge.
func EvalSignals(markets []Market, consensus map[string]float64, minEdgePP float64) []Signal {
	var signals []Signal
	for _, m := range markets {
		bookProb, ok := consensus[m.Country]
		if !ok {
			// Try fuzzy match
			for k, v := range consensus {
				if strings.EqualFold(k, m.Country) || strings.Contains(strings.ToLower(k), strings.ToLower(m.Country)) {
					bookProb = v
					ok = true
					break
				}
			}
		}
		if !ok {
			continue
		}

		yesEdge := bookProb - m.YesPrice
		noEdge := (1 - bookProb) - m.NoPrice

		if yesEdge*100 >= minEdgePP {
			signals = append(signals, Signal{
				Market:   m,
				BookProb: bookProb,
				PMPrice:  m.YesPrice,
				Edge:     yesEdge,
				Side:     "YES",
				Source:   "consensus",
			})
		} else if noEdge*100 >= minEdgePP {
			signals = append(signals, Signal{
				Market:   m,
				BookProb: bookProb,
				PMPrice:  m.NoPrice,
				Edge:     noEdge,
				Side:     "NO",
				Source:   "consensus",
			})
		}
	}
	return signals
}

// FormatSignal formats a Eurovision signal for Telegram notification.
func FormatSignal(s Signal) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🎤 Eurovision\n")
	fmt.Fprintf(&b, "━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(&b, "🏳️ %s\n", s.Market.Country)
	fmt.Fprintf(&b, "📊 %s\n\n", s.Market.Question)
	fmt.Fprintf(&b, "📈 Bookmaker: %.1f%% | PM: %.1f%% | Edge: +%.1fpp\n",
		s.BookProb*100, s.PMPrice*100, s.Edge*100)
	fmt.Fprintf(&b, "💰 Volume: $%.0f/24h\n\n", s.Market.Volume24)
	if s.Side == "YES" {
		fmt.Fprintf(&b, "✅ BUY YES @ %.4f\n", s.PMPrice)
	} else {
		fmt.Fprintf(&b, "❌ BUY NO @ %.4f\n", s.PMPrice)
	}
	return b.String()
}

// ScannerConfig for the periodic Eurovision scanner.
type ScannerConfig struct {
	OddsAPIKey string
	MinEdgePP  float64
	Interval   time.Duration
}
