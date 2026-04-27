// Package elon provides Elon Musk tweet counting and PM market matching
// for "how many tweets will Elon post from date X to date Y" markets.
package elon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TweetCountMarket represents a PM market on Elon's tweet count.
type TweetCountMarket struct {
	Slug     string
	Question string
	RangeLo  int // lower bound of tweet range (inclusive)
	RangeHi  int // upper bound (inclusive), -1 for "less than" type
	Start    time.Time
	End      time.Time
	YesPrice float64
	NoPrice  float64
	CondID   string
}

// TweetSignal is emitted when the current tweet count implies edge vs PM price.
type TweetSignal struct {
	Market      TweetCountMarket
	CurrentCount int
	ModelProb   float64
	PMPrice     float64
	Edge        float64
	Side        string // "YES" or "NO"
	HoursLeft   float64
}

var (
	reElonSlug   = regexp.MustCompile(`elon-musk-of-tweets-`)
	reRange      = regexp.MustCompile(`(\d+)-(\d+)\s+tweets`)
	reLessThan   = regexp.MustCompile(`<\s*(\d+)\s+tweets`)
	reMoreThan   = regexp.MustCompile(`(\d+)\+\s+tweets`)
	reDateRange  = regexp.MustCompile(`from\s+(\w+\s+\d+)\s+to\s+(\w+\s+\d+),?\s+(\d{4})`)
)

// FetchElonMarkets retrieves active Elon tweet count markets from gamma.
func FetchElonMarkets(ctx context.Context) ([]TweetCountMarket, error) {
	url := "https://gamma-api.polymarket.com/markets?limit=200&active=true&closed=false&order=volume24hr&ascending=false"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch elon markets: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("gamma HTTP %d: %s", resp.StatusCode, body)
	}

	var raw []struct {
		Slug          string `json:"slug"`
		Question      string `json:"question"`
		EndDate       string `json:"endDate"`
		OutcomePrices string `json:"outcomePrices"`
		ConditionID   string `json:"conditionId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var out []TweetCountMarket
	for _, m := range raw {
		slug := strings.ToLower(m.Slug)
		if !reElonSlug.MatchString(slug) {
			continue
		}

		q := strings.ToLower(m.Question)
		lo, hi := parseRange(q)
		start, end := parseDateRange(q)
		if end.IsZero() {
			t, _ := time.Parse(time.RFC3339, m.EndDate)
			end = t
		}

		yes, no := parsePricesJSON(m.OutcomePrices)

		out = append(out, TweetCountMarket{
			Slug:     m.Slug,
			Question: m.Question,
			RangeLo:  lo,
			RangeHi:  hi,
			Start:    start,
			End:      end,
			YesPrice: yes,
			NoPrice:  no,
			CondID:   m.ConditionID,
		})
	}
	return out, nil
}

func parseRange(q string) (lo, hi int) {
	if m := reRange.FindStringSubmatch(q); len(m) >= 3 {
		lo, _ = strconv.Atoi(m[1])
		hi, _ = strconv.Atoi(m[2])
		return
	}
	if m := reLessThan.FindStringSubmatch(q); len(m) >= 2 {
		hi, _ = strconv.Atoi(m[1])
		hi-- // "less than 40" = 0-39
		return 0, hi
	}
	if m := reMoreThan.FindStringSubmatch(q); len(m) >= 2 {
		lo, _ = strconv.Atoi(m[1])
		return lo, 9999
	}
	return 0, 0
}

func parseDateRange(q string) (start, end time.Time) {
	m := reDateRange.FindStringSubmatch(q)
	if len(m) < 4 {
		return
	}
	year := m[3]
	start, _ = time.Parse("January 2 2006", m[1]+" "+year)
	end, _ = time.Parse("January 2 2006", m[2]+" "+year)
	if end.IsZero() {
		start, _ = time.Parse("January 2, 2006", m[1]+", "+year)
		end, _ = time.Parse("January 2, 2006", m[2]+", "+year)
	}
	return
}

func parsePricesJSON(raw string) (yes, no float64) {
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

// CountTweets fetches Elon's tweet count using the X API v2.
// Requires a Bearer token. Returns the count of tweets in the given time range.
func CountTweets(ctx context.Context, bearerToken string, start, end time.Time) (int, error) {
	// X API v2 tweet counts endpoint (academic/basic access)
	url := fmt.Sprintf("https://api.twitter.com/2/tweets/counts/all?query=from:elonmusk&start_time=%s&end_time=%s&granularity=day",
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("x api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("x api %d: %s (may need elevated access)", resp.StatusCode, body)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("x api %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Data []struct {
			TweetCount int `json:"tweet_count"`
		} `json:"data"`
		Meta struct {
			TotalTweetCount int `json:"total_tweet_count"`
		} `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode: %w", err)
	}
	return result.Meta.TotalTweetCount, nil
}

// CountTweetsViaSearch uses the user timeline search endpoint (v2 basic access).
// Paginated — counts all tweets from @elonmusk in the time range.
func CountTweetsViaSearch(ctx context.Context, bearerToken string, start, end time.Time) (int, error) {
	count := 0
	paginationToken := ""
	for {
		url := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=from:elonmusk&start_time=%s&end_time=%s&max_results=100",
			start.UTC().Format(time.RFC3339),
			end.UTC().Format(time.RFC3339))
		if paginationToken != "" {
			url += "&next_token=" + paginationToken
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return count, err
		}
		req.Header.Set("Authorization", "Bearer "+bearerToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return count, fmt.Errorf("x api search: %w", err)
		}

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			return count, fmt.Errorf("x api search %d: %s", resp.StatusCode, body)
		}

		var result struct {
			Data []json.RawMessage `json:"data"`
			Meta struct {
				NextToken   string `json:"next_token"`
				ResultCount int    `json:"result_count"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return count, fmt.Errorf("decode search: %w", err)
		}
		resp.Body.Close()

		count += result.Meta.ResultCount
		if result.Meta.NextToken == "" {
			break
		}
		paginationToken = result.Meta.NextToken
	}
	return count, nil
}

// EvalSignals evaluates all Elon tweet markets against current tweet count.
// For markets near expiry, if the count is already outside a range's bounds,
// that range is definitely not going to win → edge is large.
func EvalSignals(markets []TweetCountMarket, currentCount int, hoursLeft float64, minEdgePP float64) []TweetSignal {
	// Simple heuristic: if the current count already exceeds the range's upper bound,
	// that range is impossible → buy No. If count is already within range and
	// hours left are few, probability increases → buy Yes if PM underprices.
	var signals []TweetSignal
	for _, m := range markets {
		if m.RangeHi <= 0 && m.RangeLo <= 0 {
			continue
		}

		var modelProb float64
		if currentCount > m.RangeHi {
			modelProb = 0.0 // already exceeded upper bound
		} else if currentCount >= m.RangeLo && currentCount <= m.RangeHi {
			if hoursLeft < 4 {
				// Near expiry and count is in range — high probability
				remaining := m.RangeHi - currentCount
				if remaining > 20 {
					modelProb = 0.90 // lots of room left
				} else if remaining > 10 {
					modelProb = 0.70
				} else if remaining > 5 {
					modelProb = 0.50
				} else {
					modelProb = 0.30 // very little room, might overflow
				}
			} else {
				continue // too far from expiry, model uncertain
			}
		} else if currentCount < m.RangeLo {
			if hoursLeft < 2 {
				modelProb = 0.0 // can't catch up in 2h
			} else {
				continue // still possible, skip
			}
		}

		// Check Yes edge
		yesEdge := modelProb - m.YesPrice
		noEdge := (1 - modelProb) - m.NoPrice

		if yesEdge*100 >= minEdgePP {
			signals = append(signals, TweetSignal{
				Market:       m,
				CurrentCount: currentCount,
				ModelProb:    modelProb,
				PMPrice:      m.YesPrice,
				Edge:         yesEdge,
				Side:         "YES",
				HoursLeft:    hoursLeft,
			})
		} else if noEdge*100 >= minEdgePP {
			signals = append(signals, TweetSignal{
				Market:       m,
				CurrentCount: currentCount,
				ModelProb:    modelProb,
				PMPrice:      m.NoPrice,
				Edge:         noEdge,
				Side:         "NO",
				HoursLeft:    hoursLeft,
			})
		}
	}
	return signals
}

// FormatTweetSignal formats a signal for Telegram notification.
func FormatTweetSignal(s TweetSignal) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🐦 Elon Tweet Count\n")
	fmt.Fprintf(&b, "━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(&b, "📊 %s\n", s.Market.Question)
	fmt.Fprintf(&b, "🔢 Current Count: %d tweets\n", s.CurrentCount)
	fmt.Fprintf(&b, "🎯 Range: %d-%d\n", s.Market.RangeLo, s.Market.RangeHi)
	fmt.Fprintf(&b, "⏳ %.1f hours left\n\n", s.HoursLeft)
	fmt.Fprintf(&b, "📈 Model: %.1f%% | PM: %.1f%% | Edge: +%.1fpp\n",
		s.ModelProb*100, s.PMPrice*100, s.Edge*100)
	if s.Side == "YES" {
		fmt.Fprintf(&b, "✅ BUY YES @ %.4f\n", s.PMPrice)
	} else {
		fmt.Fprintf(&b, "❌ BUY NO @ %.4f\n", s.PMPrice)
	}
	return b.String()
}
