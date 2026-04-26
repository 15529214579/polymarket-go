package btc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"
)

type InstitutionalFlow struct {
	OpenInterest   float64 // BTC open interest in contracts
	LongShortRatio float64 // long/short account ratio (>1 = crowd long)
	FuturesPremium float64 // (mark - index) / index as percentage
	FlowSignal     string  // "BULLISH" / "BEARISH" / "NEUTRAL"
	FlowScore      float64 // -1.0 to +1.0 (positive = bullish flow)
}

func FetchInstitutionalFlow(ctx context.Context) InstitutionalFlow {
	var flow InstitutionalFlow

	type oiResp struct {
		OpenInterest string `json:"openInterest"`
	}
	type lsResp struct {
		LongShortRatio string `json:"longShortRatio"`
	}
	type premResp struct {
		MarkPrice  string `json:"markPrice"`
		IndexPrice string `json:"indexPrice"`
	}

	client := &http.Client{Timeout: 10 * time.Second}

	if resp, err := fetchJSON[oiResp](ctx, client, "https://fapi.binance.com/fapi/v1/openInterest?symbol=BTCUSDT"); err == nil {
		flow.OpenInterest, _ = strconv.ParseFloat(resp.OpenInterest, 64)
	}

	if resp, err := fetchJSONArray[lsResp](ctx, client, "https://fapi.binance.com/futures/data/globalLongShortAccountRatio?symbol=BTCUSDT&period=1h&limit=1"); err == nil && len(resp) > 0 {
		flow.LongShortRatio, _ = strconv.ParseFloat(resp[0].LongShortRatio, 64)
	}

	if resp, err := fetchJSON[premResp](ctx, client, "https://fapi.binance.com/fapi/v1/premiumIndex?symbol=BTCUSDT"); err == nil {
		mark, _ := strconv.ParseFloat(resp.MarkPrice, 64)
		index, _ := strconv.ParseFloat(resp.IndexPrice, 64)
		if index > 0 {
			flow.FuturesPremium = (mark - index) / index * 100
		}
	}

	score := 0.0

	// Long/short ratio signal: <0.85 = crowd short (contrarian bullish), >1.15 = crowd long (contrarian bearish)
	if flow.LongShortRatio > 0 {
		if flow.LongShortRatio < 0.85 {
			score += 0.4 // crowd short → contrarian bullish
		} else if flow.LongShortRatio > 1.15 {
			score -= 0.4 // crowd long → contrarian bearish
		}
	}

	// Futures premium: positive = futures at premium (bullish), negative = backwardation (bearish)
	if flow.FuturesPremium > 0.05 {
		score += 0.3
	} else if flow.FuturesPremium < -0.05 {
		score -= 0.3
	}

	// OI direction would need history comparison — for now just log
	flow.FlowScore = math.Max(-1, math.Min(1, score))

	switch {
	case flow.FlowScore > 0.2:
		flow.FlowSignal = "BULLISH"
	case flow.FlowScore < -0.2:
		flow.FlowSignal = "BEARISH"
	default:
		flow.FlowSignal = "NEUTRAL"
	}

	slog.Info("institutional.flow",
		"oi_btc", fmt.Sprintf("%.0f", flow.OpenInterest),
		"long_short", fmt.Sprintf("%.3f", flow.LongShortRatio),
		"futures_prem", fmt.Sprintf("%.4f%%", flow.FuturesPremium),
		"signal", flow.FlowSignal,
		"score", fmt.Sprintf("%.2f", flow.FlowScore),
	)

	return flow
}

func SaveInstitutionalFlow(ctx context.Context, db *sql.DB, flow InstitutionalFlow) error {
	const ddl = `CREATE TABLE IF NOT EXISTS btc_institutional (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp       INTEGER NOT NULL,
		open_interest   REAL,
		long_short      REAL,
		futures_premium REAL,
		flow_signal     TEXT,
		flow_score      REAL
	);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx,
		`INSERT INTO btc_institutional(timestamp, open_interest, long_short, futures_premium, flow_signal, flow_score)
		 VALUES(?,?,?,?,?,?)`,
		time.Now().Unix(), flow.OpenInterest, flow.LongShortRatio, flow.FuturesPremium, flow.FlowSignal, flow.FlowScore)
	return err
}

// InstitutionalModifier returns a multiplier for BS gap signals based on institutional flow.
// Bullish flow amplifies BUY_YES on reach markets and dampens BUY_YES on dip markets.
func (f InstitutionalFlow) InstitutionalModifier(direction string, isReach bool) float64 {
	if f.FlowSignal == "NEUTRAL" {
		return 1.0
	}
	mod := 1.0
	switch {
	case f.FlowSignal == "BULLISH" && direction == "BUY_YES" && isReach:
		mod = 1.0 + f.FlowScore*0.15 // bullish flow supports reach BUY_YES
	case f.FlowSignal == "BULLISH" && direction == "BUY_NO" && !isReach:
		mod = 1.0 - f.FlowScore*0.10 // bullish flow weakens dip BUY_NO (less likely to dip)
	case f.FlowSignal == "BEARISH" && direction == "BUY_NO" && !isReach:
		mod = 1.0 + math.Abs(f.FlowScore)*0.15 // bearish flow supports dip BUY_NO
	case f.FlowSignal == "BEARISH" && direction == "BUY_YES" && isReach:
		mod = 1.0 - math.Abs(f.FlowScore)*0.10 // bearish flow weakens reach BUY_YES
	}
	return mod
}

func fetchJSON[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var zero T
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return zero, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, err
	}
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return zero, fmt.Errorf("json decode: %w (body: %s)", err, string(body[:min(len(body), 200)]))
	}
	return result, nil
}

func fetchJSONArray[T any](ctx context.Context, client *http.Client, url string) ([]T, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
