package btc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	fngAPIURL     = "https://api.alternative.me/fng/?limit=1"
	fundingAPIURL = "https://fapi.binance.com/fapi/v1/fundingRate?symbol=BTCUSDT&limit=1"
)

type FearGreed struct {
	Value          int       // 0-100
	Classification string   // "Extreme Fear", "Fear", "Neutral", "Greed", "Extreme Greed"
	Timestamp      time.Time
}

type FundingRate struct {
	Rate      float64   // e.g. 0.0001 = 0.01%
	MarkPrice float64
	Timestamp time.Time
}

type Sentiment struct {
	FearGreed   *FearGreed
	FundingRate *FundingRate
}

// SentimentModifier returns a multiplier for BS gap sizing.
// >1 = amplify signal, <1 = dampen signal.
//
// Logic:
//   - Extreme Fear (<25) + BUY_YES on dip market → dampen (dip more likely, PM may be right)
//   - Extreme Fear (<25) + BUY_NO on dip market  → amplify (fear overestimates dip probability)
//   - Extreme Greed (>75) + BUY_YES on reach      → dampen (greed inflates reach expectations)
//   - Extreme Greed (>75) + BUY_NO on reach       → amplify
//   - Positive funding + BUY_YES on reach          → dampen (crowded long → pullback risk)
//   - Negative funding + BUY_YES on dip            → dampen (crowded short → bounce risk)
func (s *Sentiment) SentimentModifier(direction string, isReach bool) float64 {
	mod := 1.0

	if s.FearGreed != nil {
		fng := s.FearGreed.Value
		if fng < 25 {
			if isReach && direction == "BUY_YES" {
				mod *= 0.8
			} else if !isReach && direction == "BUY_NO" {
				mod *= 1.2
			}
		} else if fng > 75 {
			if isReach && direction == "BUY_YES" {
				mod *= 0.8
			} else if !isReach && direction == "BUY_NO" {
				mod *= 0.8
			}
		}
	}

	if s.FundingRate != nil {
		rate := s.FundingRate.Rate
		if rate > 0.0005 {
			if isReach && direction == "BUY_YES" {
				mod *= 0.85
			}
		} else if rate < -0.0005 {
			if !isReach && direction == "BUY_YES" {
				mod *= 0.85
			}
		}
	}

	return mod
}

func FetchFearGreed(ctx context.Context) (*FearGreed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fngAPIURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch fng: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("fng HTTP %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Data []struct {
			Value               string `json:"value"`
			ValueClassification string `json:"value_classification"`
			Timestamp           string `json:"timestamp"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode fng: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("fng: empty data")
	}

	d := result.Data[0]
	val, _ := strconv.Atoi(d.Value)
	ts, _ := strconv.ParseInt(d.Timestamp, 10, 64)

	return &FearGreed{
		Value:          val,
		Classification: d.ValueClassification,
		Timestamp:      time.Unix(ts, 0).UTC(),
	}, nil
}

func FetchFundingRate(ctx context.Context) (*FundingRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fundingAPIURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch funding: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("funding HTTP %d: %s", resp.StatusCode, body)
	}

	var result []struct {
		FundingRate string `json:"fundingRate"`
		MarkPrice   string `json:"markPrice"`
		FundingTime int64  `json:"fundingTime"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode funding: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("funding: empty data")
	}

	r := result[0]
	rate, _ := strconv.ParseFloat(r.FundingRate, 64)
	mark, _ := strconv.ParseFloat(r.MarkPrice, 64)

	return &FundingRate{
		Rate:      rate,
		MarkPrice: mark,
		Timestamp: time.UnixMilli(r.FundingTime).UTC(),
	}, nil
}

func FetchSentiment(ctx context.Context) *Sentiment {
	s := &Sentiment{}

	fng, err := FetchFearGreed(ctx)
	if err != nil {
		slog.Warn("sentiment.fng_fail", "err", err.Error())
	} else {
		s.FearGreed = fng
		slog.Info("sentiment.fng", "value", fng.Value, "class", fng.Classification)
	}

	fr, err := FetchFundingRate(ctx)
	if err != nil {
		slog.Warn("sentiment.funding_fail", "err", err.Error())
	} else {
		s.FundingRate = fr
		slog.Info("sentiment.funding", "rate", fmt.Sprintf("%.6f", fr.Rate), "mark", fmt.Sprintf("%.2f", fr.MarkPrice))
	}

	return s
}

func InitSentimentDB(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS btc_sentiment (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp  INTEGER NOT NULL,
    fng_value  INTEGER,
    fng_class  TEXT,
    funding_rate REAL,
    funding_mark REAL,
    UNIQUE(timestamp)
);`
	_, err := db.Exec(ddl)
	return err
}

func SaveSentiment(ctx context.Context, db *sql.DB, s *Sentiment) error {
	if err := InitSentimentDB(db); err != nil {
		return err
	}

	now := time.Now().Unix()
	// round to 5-minute boundary to avoid duplicate spam
	now = (now / 300) * 300

	var fngVal *int
	var fngClass *string
	var fundRate *float64
	var fundMark *float64

	if s.FearGreed != nil {
		v := s.FearGreed.Value
		c := s.FearGreed.Classification
		fngVal = &v
		fngClass = &c
	}
	if s.FundingRate != nil {
		r := s.FundingRate.Rate
		m := s.FundingRate.MarkPrice
		fundRate = &r
		fundMark = &m
	}

	_, err := db.ExecContext(ctx,
		`INSERT OR IGNORE INTO btc_sentiment(timestamp, fng_value, fng_class, funding_rate, funding_mark)
		 VALUES(?,?,?,?,?)`,
		now, fngVal, fngClass, fundRate, fundMark)
	return err
}
