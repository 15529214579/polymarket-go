package btc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type OnChainMetrics struct {
	MempoolTxs       int     // pending transactions
	AvgFee24h        float64 // average tx fee in USD (24h)
	HashRate         float64 // network hashrate (TH/s)
	Difficulty       float64 // current difficulty
	BlockHeight      int64   // latest block
	MempoolSize      float64 // mempool size in bytes
	ExchangeNetFlow  string  // "INFLOW" / "OUTFLOW" / "NEUTRAL" (from mempool heuristic)
	OnChainSignal    string  // "BULLISH" / "BEARISH" / "NEUTRAL"
	OnChainScore     float64 // -1.0 to +1.0
}

type blockchairStats struct {
	Data struct {
		MempoolTxs       int             `json:"mempool_transactions"`
		AvgFee24h        json.Number     `json:"average_transaction_fee_24h"`
		HashRate         json.Number     `json:"hashrate_24h"`
		Difficulty       json.Number     `json:"difficulty"`
		BestBlockHeight  int64           `json:"best_block_height"`
		MempoolSize      json.Number     `json:"mempool_size"`
	} `json:"data"`
}

func FetchOnChainMetrics(ctx context.Context) OnChainMetrics {
	var m OnChainMetrics
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.blockchair.com/bitcoin/stats", nil)
	if err != nil {
		slog.Warn("onchain.fetch_fail", "err", err.Error())
		return m
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("onchain.fetch_fail", "err", err.Error())
		return m
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32768))
	if err != nil {
		slog.Warn("onchain.read_fail", "err", err.Error())
		return m
	}

	var stats blockchairStats
	if err := json.Unmarshal(body, &stats); err != nil {
		slog.Warn("onchain.decode_fail", "err", err.Error())
		return m
	}

	m.MempoolTxs = stats.Data.MempoolTxs
	m.AvgFee24h, _ = stats.Data.AvgFee24h.Float64()
	m.HashRate, _ = stats.Data.HashRate.Float64()
	m.Difficulty, _ = stats.Data.Difficulty.Float64()
	m.BlockHeight = stats.Data.BestBlockHeight
	m.MempoolSize, _ = stats.Data.MempoolSize.Float64()

	score := 0.0

	// High mempool = congestion = high activity (slightly bullish — demand for block space)
	if m.MempoolTxs > 50000 {
		score += 0.15
	} else if m.MempoolTxs < 5000 {
		score -= 0.10 // low activity = less interest
	}

	// High fees = network demand (bullish signal)
	if m.AvgFee24h > 5.0 {
		score += 0.15
	} else if m.AvgFee24h < 0.5 {
		score -= 0.10
	}

	// Hashrate stability: above 800 EH/s is healthy (very rough heuristic)
	if m.HashRate > 800e18 {
		score += 0.10
	}

	m.OnChainScore = clampFloat(score, -1, 1)

	switch {
	case m.OnChainScore > 0.15:
		m.OnChainSignal = "BULLISH"
	case m.OnChainScore < -0.15:
		m.OnChainSignal = "BEARISH"
	default:
		m.OnChainSignal = "NEUTRAL"
	}

	// Mempool heuristic for exchange flow direction
	if m.MempoolTxs > 30000 && m.AvgFee24h > 3.0 {
		m.ExchangeNetFlow = "INFLOW" // high activity + high fees often means exchange deposits
	} else if m.MempoolTxs < 10000 {
		m.ExchangeNetFlow = "OUTFLOW" // low mempool = accumulation phase
	} else {
		m.ExchangeNetFlow = "NEUTRAL"
	}

	slog.Info("onchain.metrics",
		"mempool_txs", m.MempoolTxs,
		"avg_fee", fmt.Sprintf("$%.2f", m.AvgFee24h),
		"hashrate", fmt.Sprintf("%.2e", m.HashRate),
		"block", m.BlockHeight,
		"flow", m.ExchangeNetFlow,
		"signal", m.OnChainSignal,
		"score", fmt.Sprintf("%.2f", m.OnChainScore),
	)

	return m
}

func SaveOnChainMetrics(ctx context.Context, db *sql.DB, m OnChainMetrics) error {
	const ddl = `CREATE TABLE IF NOT EXISTS btc_onchain (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp     INTEGER NOT NULL,
		mempool_txs   INTEGER,
		avg_fee       REAL,
		hashrate      REAL,
		difficulty    REAL,
		block_height  INTEGER,
		mempool_size  REAL,
		exchange_flow TEXT,
		signal        TEXT,
		score         REAL
	);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx,
		`INSERT INTO btc_onchain(timestamp, mempool_txs, avg_fee, hashrate, difficulty, block_height, mempool_size, exchange_flow, signal, score)
		 VALUES(?,?,?,?,?,?,?,?,?,?)`,
		time.Now().Unix(), m.MempoolTxs, m.AvgFee24h, m.HashRate, m.Difficulty, m.BlockHeight, m.MempoolSize, m.ExchangeNetFlow, m.OnChainSignal, m.OnChainScore)
	return err
}

func (m OnChainMetrics) OnChainModifier(direction string, isReach bool) float64 {
	if m.OnChainSignal == "NEUTRAL" {
		return 1.0
	}
	mod := 1.0
	switch {
	case m.OnChainSignal == "BULLISH" && direction == "BUY_YES" && isReach:
		mod = 1.0 + m.OnChainScore*0.10
	case m.OnChainSignal == "BULLISH" && direction == "BUY_NO" && !isReach:
		mod = 1.0 - m.OnChainScore*0.08
	case m.OnChainSignal == "BEARISH" && direction == "BUY_NO" && !isReach:
		mod = 1.0 + (-m.OnChainScore)*0.10
	case m.OnChainSignal == "BEARISH" && direction == "BUY_YES" && isReach:
		mod = 1.0 + m.OnChainScore*0.08
	}
	return mod
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
