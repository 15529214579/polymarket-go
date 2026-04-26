// Package btc provides BTC price data fetching from Binance and Polymarket
// BTC market price tracking for Markov-based prediction.
package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const binanceKlinesURL = "https://api.binance.com/api/v3/klines"

// Candle holds one OHLCV bar from Binance.
type Candle struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// Interval controls the kline granularity.
type Interval string

const (
	Interval1h Interval = "1h"
	Interval1d Interval = "1d"
)

// FetchCandles fetches up to limit klines of the given interval from Binance.
// No API key is required for public market data.
func FetchCandles(ctx context.Context, symbol string, interval Interval, limit int) ([]Candle, error) {
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	url := fmt.Sprintf("%s?symbol=%s&interval=%s&limit=%d",
		binanceKlinesURL, symbol, interval, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch candles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("binance HTTP %d: %s", resp.StatusCode, body)
	}

	// Binance returns a JSON array of arrays:
	// [openTime, open, high, low, close, volume, closeTime, ...]
	var raw [][]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode klines: %w", err)
	}

	candles := make([]Candle, 0, len(raw))
	for _, row := range raw {
		if len(row) < 6 {
			continue
		}
		c, err := parseKlineRow(row)
		if err != nil {
			continue
		}
		candles = append(candles, c)
	}
	return candles, nil
}

// FetchCandlesRange fetches candles between startTime and endTime by making
// multiple paginated requests (Binance caps each at 1000 bars).
func FetchCandlesRange(ctx context.Context, symbol string, interval Interval, start, end time.Time) ([]Candle, error) {
	var all []Candle
	cursor := start

	for cursor.Before(end) {
		url := fmt.Sprintf("%s?symbol=%s&interval=%s&limit=1000&startTime=%d&endTime=%d",
			binanceKlinesURL, symbol, interval,
			cursor.UnixMilli(), end.UnixMilli())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch candles: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			resp.Body.Close()
			return nil, fmt.Errorf("binance HTTP %d: %s", resp.StatusCode, body)
		}

		var raw [][]json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode klines: %w", err)
		}
		resp.Body.Close()

		if len(raw) == 0 {
			break
		}

		for _, row := range raw {
			if len(row) < 6 {
				continue
			}
			c, err := parseKlineRow(row)
			if err != nil {
				continue
			}
			if c.Timestamp.After(end) {
				continue
			}
			all = append(all, c)
		}

		// advance cursor past the last candle we received
		last := raw[len(raw)-1]
		if len(last) < 1 {
			break
		}
		var openTimeMs int64
		if err := json.Unmarshal(last[0], &openTimeMs); err != nil {
			break
		}
		nextCursor := time.UnixMilli(openTimeMs).Add(1 * time.Millisecond)
		if !nextCursor.After(cursor) {
			break // guard against infinite loop
		}
		cursor = nextCursor

		// be a good citizen — short pause between pages
		select {
		case <-ctx.Done():
			return all, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return all, nil
}

func parseKlineRow(row []json.RawMessage) (Candle, error) {
	// field 0: openTime (int ms)
	var openTimeMs int64
	if err := json.Unmarshal(row[0], &openTimeMs); err != nil {
		return Candle{}, fmt.Errorf("openTime: %w", err)
	}

	parseFloat := func(r json.RawMessage) (float64, error) {
		var s string
		if err := json.Unmarshal(r, &s); err != nil {
			// might be a bare number
			var f float64
			if err2 := json.Unmarshal(r, &f); err2 != nil {
				return 0, fmt.Errorf("parse float: %w", err)
			}
			return f, nil
		}
		return strconv.ParseFloat(s, 64)
	}

	open, err := parseFloat(row[1])
	if err != nil {
		return Candle{}, fmt.Errorf("open: %w", err)
	}
	high, err := parseFloat(row[2])
	if err != nil {
		return Candle{}, fmt.Errorf("high: %w", err)
	}
	low, err := parseFloat(row[3])
	if err != nil {
		return Candle{}, fmt.Errorf("low: %w", err)
	}
	close_, err := parseFloat(row[4])
	if err != nil {
		return Candle{}, fmt.Errorf("close: %w", err)
	}
	volume, err := parseFloat(row[5])
	if err != nil {
		return Candle{}, fmt.Errorf("volume: %w", err)
	}

	return Candle{
		Timestamp: time.UnixMilli(openTimeMs).UTC(),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close_,
		Volume:    volume,
	}, nil
}
