package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DailyMarket struct {
	Slug       string
	Question   string
	Strike     float64
	EndDate    time.Time
	YesPrice   float64
	NoPrice    float64
	CondID     string
	ClobTokens []string
}

type DailySignal struct {
	Market     DailyMarket
	Spot       float64
	ModelProb  float64
	PMPrice    float64
	Edge       float64 // ModelProb - PMPrice (positive = buy Yes)
	Side       string  // "YES" or "NO"
	HoursLeft  float64
	AnnualVol  float64
}

var reBTCDaily = regexp.MustCompile(`^bitcoin-above-(\d+)k?-on-`)

func FetchDailyBTCMarkets(ctx context.Context) ([]DailyMarket, error) {
	url := "https://gamma-api.polymarket.com/markets?limit=200&active=true&closed=false&order=volume24hr&ascending=false"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch daily btc: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("gamma HTTP %d: %s", resp.StatusCode, body)
	}

	var raw []struct {
		Slug         string `json:"slug"`
		Question     string `json:"question"`
		EndDate      string `json:"endDate"`
		OutcomePrices string `json:"outcomePrices"`
		ConditionID  string `json:"conditionId"`
		ClobTokenIDs string `json:"clobTokenIds"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var out []DailyMarket
	for _, m := range raw {
		slug := strings.ToLower(m.Slug)
		// Match bitcoin-above-XXk-on-<date> or will-bitcoin-reach/dip patterns
		if !strings.Contains(slug, "bitcoin") {
			continue
		}
		isDailyAbove := reBTCDaily.MatchString(slug)
		isMonthlyReach := strings.Contains(slug, "will-bitcoin-reach") || strings.Contains(slug, "will-bitcoin-dip")
		if !isDailyAbove && !isMonthlyReach {
			continue
		}

		strike := parseStrikeFromQuestion(m.Question)
		if strike <= 0 {
			continue
		}

		end, _ := time.Parse(time.RFC3339, m.EndDate)
		if end.IsZero() {
			end, _ = time.Parse("2006-01-02T15:04:05Z", m.EndDate)
		}
		if end.IsZero() {
			continue
		}

		yes, no := parsePrices(m.OutcomePrices)
		var tokens []string
		if m.ClobTokenIDs != "" {
			_ = json.Unmarshal([]byte(m.ClobTokenIDs), &tokens)
		}

		out = append(out, DailyMarket{
			Slug:       m.Slug,
			Question:   m.Question,
			Strike:     strike,
			EndDate:    end,
			YesPrice:   yes,
			NoPrice:    no,
			CondID:     m.ConditionID,
			ClobTokens: tokens,
		})
	}
	return out, nil
}

// AboveAtExpiry returns P(S_T > K) under GBM with zero drift.
func AboveAtExpiry(spot, strike, annualVol, yearsToExpiry float64) float64 {
	if annualVol <= 0 || yearsToExpiry <= 0 || spot <= 0 || strike <= 0 {
		if spot > strike {
			return 1.0
		}
		return 0.0
	}
	d2 := (math.Log(spot/strike) - 0.5*annualVol*annualVol*yearsToExpiry) / (annualVol * math.Sqrt(yearsToExpiry))
	return normCDF(d2)
}

// ScanDailyBTC fetches current BTC spot, volatility, and all daily threshold
// markets, then returns signals where edge exceeds minEdgePP percentage points.
func ScanDailyBTC(ctx context.Context, minEdgePP float64) ([]DailySignal, error) {
	candles, err := FetchCandles(ctx, "BTCUSDT", Interval1h, 168) // 7 days
	if err != nil {
		return nil, fmt.Errorf("fetch candles: %w", err)
	}
	if len(candles) < 10 {
		return nil, fmt.Errorf("insufficient candle data: %d", len(candles))
	}

	spot := candles[len(candles)-1].Close
	annualVol := EWMAVolatility(candles, 0.94)
	if annualVol <= 0 {
		annualVol = HistoricalVolatility(candles)
	}

	markets, err := FetchDailyBTCMarkets(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch markets: %w", err)
	}

	now := time.Now().UTC()
	var signals []DailySignal
	for _, m := range markets {
		hoursLeft := m.EndDate.Sub(now).Hours()
		if hoursLeft <= 0 {
			continue // already expired
		}
		yearsLeft := hoursLeft / 8760.0

		isDip := strings.Contains(strings.ToLower(m.Slug), "dip")
		var modelProb float64
		if isDip {
			// "will bitcoin dip to $K" = P(min < K) at any point
			modelProb = FirstPassageProb(spot, m.Strike, annualVol, yearsLeft)
		} else {
			// "above $K on date" = P(S_T > K) at expiry
			modelProb = AboveAtExpiry(spot, m.Strike, annualVol, yearsLeft)
		}

		// Check Yes side edge
		yesEdge := modelProb - m.YesPrice
		// Check No side edge
		noEdge := (1 - modelProb) - m.NoPrice

		var sig DailySignal
		sig.Market = m
		sig.Spot = spot
		sig.ModelProb = modelProb
		sig.HoursLeft = hoursLeft
		sig.AnnualVol = annualVol

		if yesEdge*100 >= minEdgePP {
			sig.Side = "YES"
			sig.PMPrice = m.YesPrice
			sig.Edge = yesEdge
			signals = append(signals, sig)
		} else if noEdge*100 >= minEdgePP {
			sig.Side = "NO"
			sig.PMPrice = m.NoPrice
			sig.Edge = noEdge
			signals = append(signals, sig)
		}
	}
	return signals, nil
}

// GetBTCSpot returns the current BTC/USDT price from the latest 1h candle.
func GetBTCSpot(ctx context.Context) (float64, error) {
	candles, err := FetchCandles(ctx, "BTCUSDT", Interval1h, 1)
	if err != nil {
		return 0, err
	}
	if len(candles) == 0 {
		return 0, fmt.Errorf("no candle data")
	}
	return candles[0].Close, nil
}

// FormatDailySignal formats a signal for Telegram notification.
func FormatDailySignal(s DailySignal) string {
	sgt := s.Market.EndDate.In(time.FixedZone("SGT", 8*3600))
	var b strings.Builder
	fmt.Fprintf(&b, "₿ BTC Daily Threshold\n")
	fmt.Fprintf(&b, "━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(&b, "📊 %s\n", s.Market.Question)
	fmt.Fprintf(&b, "💰 BTC Spot: $%.0f\n", s.Spot)
	fmt.Fprintf(&b, "🎯 Strike: $%.0f\n", s.Market.Strike)
	fmt.Fprintf(&b, "⏰ Expires: %s（SGT）\n", sgt.Format("2006-01-02 15:04"))
	fmt.Fprintf(&b, "⏳ %.1f hours left\n\n", s.HoursLeft)
	fmt.Fprintf(&b, "📈 Model: %.1f%% | PM: %.1f%% | Edge: +%.1fpp\n",
		s.ModelProb*100, s.PMPrice*100, s.Edge*100)
	fmt.Fprintf(&b, "📐 Vol: %.1f%% annualized\n\n", s.AnnualVol*100)

	if s.Side == "YES" {
		fmt.Fprintf(&b, "✅ BUY YES @ %.4f\n", s.PMPrice)
		var pct float64
		if s.PMPrice > 0 {
			pct = (1/s.PMPrice - 1) * 100
		}
		fmt.Fprintf(&b, "💵 Payout: $1.00 per share (%.0f%% return)\n", pct)
	} else {
		fmt.Fprintf(&b, "❌ BUY NO @ %.4f\n", s.PMPrice)
		var pct float64
		if s.PMPrice > 0 {
			pct = (1/s.PMPrice - 1) * 100
		}
		fmt.Fprintf(&b, "💵 Payout: $1.00 per share (%.0f%% return)\n", pct)
	}
	return b.String()
}

// ParseBTCDailyStrike extracts threshold from slug like "bitcoin-above-72k-on-april-27".
func ParseBTCDailyStrike(slug string) float64 {
	m := reBTCDaily.FindStringSubmatch(slug)
	if len(m) < 2 {
		return 0
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0
	}
	return v * 1000
}
