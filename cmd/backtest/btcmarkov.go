package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/15529214579/polymarket-go/internal/btc"
	_ "modernc.org/sqlite" //nolint:revive // register sqlite driver
)

// ---------------------------------------------------------------------------
// Entry point called from main.go dispatch
// ---------------------------------------------------------------------------

func runBTCMarkovBacktest(days int, trainPct float64, dbDir string) error {
	ctx := context.Background()

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" BTC Markov Chain Backtest")
	fmt.Printf(" days=%d  train=%.0f%%  test=%.0f%%\n",
		days, trainPct*100, (1-trainPct)*100)
	fmt.Println("════════════════════════════════════════")

	// ── 1. Open / initialise SQLite ────────────────────────────────────────
	dbPath := filepath.Join(dbDir, "btc_markov.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open btc_markov.db: %w", err)
	}
	defer db.Close()

	if err := btc.InitDB(db); err != nil {
		return fmt.Errorf("init db: %w", err)
	}

	// ── 2. Fetch hourly candles (90 days → 2160 hourly bars) ───────────────
	fmt.Printf("\n── Step 1: Fetch %d days of 1h BTC candles from Binance ──\n", days)
	end := time.Now().UTC().Truncate(time.Hour)
	start := end.Add(-time.Duration(days) * 24 * time.Hour)

	candles, err := btc.FetchCandlesRange(ctx, "BTCUSDT", btc.Interval1h, start, end)
	if err != nil {
		return fmt.Errorf("fetch candles: %w", err)
	}
	if len(candles) < 48 {
		return fmt.Errorf("too few candles (%d); check connectivity", len(candles))
	}
	fmt.Printf("  fetched %d hourly candles  [%s → %s]\n",
		len(candles),
		candles[0].Timestamp.Format("2006-01-02"),
		candles[len(candles)-1].Timestamp.Format("2006-01-02"))

	if err := btc.SaveCandles(ctx, db, candles, btc.Interval1h); err != nil {
		fmt.Fprintf(os.Stderr, "warn: save candles: %v\n", err)
	}

	// ── 3. Train / test split ───────────────────────────────────────────────
	splitIdx := int(float64(len(candles)) * trainPct)
	if splitIdx < 24 || splitIdx >= len(candles)-24 {
		return fmt.Errorf("train/test split too tight: splitIdx=%d len=%d", splitIdx, len(candles))
	}
	trainCandles := candles[:splitIdx]
	testCandles := candles[splitIdx-1:] // overlap by 1 so first test candle has a previous

	fmt.Printf("  train: %d bars  test: %d bars\n", len(trainCandles), len(testCandles)-1)

	// ── 4. Train model ─────────────────────────────────────────────────────
	fmt.Println("\n── Step 2: Train Markov transition matrix ──")
	tm, nTransitions := btc.Train(trainCandles)
	retStats := btc.BuildReturnStats(trainCandles)

	fmt.Printf("  transitions recorded: %d\n", nTransitions)
	printTransitionMatrix(&tm)
	printStateReturnTable(retStats)

	// ── 5. Backtest on test window ─────────────────────────────────────────
	fmt.Println("\n── Step 3: Backtest on test window ──")
	result := backtestBTCMarkov(testCandles, &tm, retStats, db, ctx)
	printBTCBacktestResult(result)

	// ── 6. Fetch PM markets & compute gap opportunities ────────────────────
	fmt.Println("\n── Step 4: Fetch PM BTC markets & compute gaps ──")
	markets, err := btc.FetchBTCMarkets(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  warn: fetch PM markets: %v\n", err)
		fmt.Println("  (skipping PM gap analysis — check connectivity)")
	} else {
		fmt.Printf("  fetched %d PM BTC markets\n", len(markets))
		if err := btc.SavePMPrices(ctx, db, markets); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: save PM prices: %v\n", err)
		}
		printPMMarkets(markets)
		analyzePMGaps(candles, &tm, retStats, markets)
	}

	fmt.Println("\n════════════════════════════════════════")
	fmt.Println(" Backtest complete.  Results in:", dbPath)
	fmt.Println("════════════════════════════════════════")
	return nil
}

// ---------------------------------------------------------------------------
// Backtest engine
// ---------------------------------------------------------------------------

type btcDayResult struct {
	Date          time.Time
	PredictedBull bool   // model predicted bullish
	ActualBull    bool   // next-day close > open
	PredState     int
	ActualState   int
	PredReturn    float64 // model's expected return (pp)
	ActualReturn  float64 // actual 1h return
	SimPnL        float64 // simulated PnL for this bar
}

type btcBacktestResult struct {
	Days         int
	Correct      int
	Wrong        int
	SumSimPnL    float64
	MaxDD        float64
	BestState    int
	WorstState   int
	StateResults [btc.NStates]struct {
		N       int
		Correct int
		SumPnL  float64
	}
}

func (r *btcBacktestResult) accuracy() float64 {
	total := r.Correct + r.Wrong
	if total == 0 {
		return 0
	}
	return 100 * float64(r.Correct) / float64(total)
}

func backtestBTCMarkov(
	candles []btc.Candle,
	tm *btc.TransitionMatrix,
	retStats [btc.NStates]btc.ReturnStats,
	db *sql.DB,
	ctx context.Context,
) btcBacktestResult {
	var result btcBacktestResult

	// Simulate with a unit bet per bar
	equity, peak := 0.0, 0.0

	for i := 1; i+1 < len(candles); i++ {
		state, ok := btc.CandleState(candles, i)
		if !ok {
			continue
		}

		pred := btc.Predict(state, tm, retStats)

		// Determine actual next-bar direction
		curr := candles[i]
		next := candles[i+1]
		if curr.Close <= 0 || next.Close <= 0 {
			continue
		}
		actualRet := (next.Close - curr.Close) / curr.Close * 100
		actualBull := actualRet > 0

		predBull := pred.BullProb > pred.BearProb

		correct := predBull == actualBull
		if correct {
			result.Correct++
		} else {
			result.Wrong++
		}
		result.Days++

		// Simple simulation: bet +1 if bull pred, -1 if bear pred.
		// Gain/loss = |actualRet| if correct, else -|actualRet|. Fee: 2bp.
		fee := 0.0002
		var simPnL float64
		if correct {
			simPnL = math.Abs(actualRet) - fee*100
		} else {
			simPnL = -(math.Abs(actualRet) + fee*100)
		}
		result.SumSimPnL += simPnL

		// Track by state
		sr := &result.StateResults[state]
		sr.N++
		if correct {
			sr.Correct++
		}
		sr.SumPnL += simPnL

		// Drawdown tracking
		equity += simPnL
		if equity > peak {
			peak = equity
		}
		dd := peak - equity
		if dd > result.MaxDD {
			result.MaxDD = dd
		}

		// Persist prediction
		if db != nil {
			nextState, ok2 := btc.CandleState(candles, i+1)
			_ = ok2
			_, _ = db.ExecContext(ctx, `
INSERT INTO btc_predictions(timestamp, predicted_state, actual_state, predicted_return, actual_return, bull_prob, bear_prob)
VALUES(?,?,?,?,?,?,?)`,
				curr.Timestamp.Unix(), state, nextState, pred.ExpectedReturn, actualRet,
				pred.BullProb, pred.BearProb,
			)
		}
	}

	// Find best/worst states by accuracy
	bestAcc, worstAcc := -1.0, 101.0
	for s := 0; s < btc.NStates; s++ {
		sr := result.StateResults[s]
		if sr.N < 5 {
			continue
		}
		acc := 100 * float64(sr.Correct) / float64(sr.N)
		if acc > bestAcc {
			bestAcc = acc
			result.BestState = s
		}
		if acc < worstAcc {
			worstAcc = acc
			result.WorstState = s
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Display helpers
// ---------------------------------------------------------------------------

func printTransitionMatrix(tm *btc.TransitionMatrix) {
	fmt.Println("\n  Transition matrix (row=from, col=to) — probabilities:")
	header := fmt.Sprintf("  %-14s", "from \\ to")
	for s := 0; s < btc.NStates; s++ {
		header += fmt.Sprintf(" %8s", btc.StateName(s))
	}
	fmt.Println(header)

	for from := 0; from < btc.NStates; from++ {
		probs := tm.RowProbs(from)
		row := fmt.Sprintf("  %-14s", btc.StateName(from))
		for to := 0; to < btc.NStates; to++ {
			if probs == nil {
				row += fmt.Sprintf(" %8s", "-")
			} else {
				row += fmt.Sprintf(" %8.3f", probs[to])
			}
		}
		fmt.Println(row)
	}
}

func printStateReturnTable(retStats [btc.NStates]btc.ReturnStats) {
	fmt.Println("\n  Per-state forward-return statistics (training set):")
	fmt.Printf("  %-14s %6s %+9s %7s\n", "state", "n", "E[ret%]", "P(>0)")

	type row struct {
		name   string
		n      int
		avgRet float64
		posR   float64
	}
	var rows []row
	for s := 0; s < btc.NStates; s++ {
		rs := retStats[s]
		if rs.N == 0 {
			continue
		}
		rows = append(rows, row{btc.StateName(s), rs.N, rs.AvgRet(), rs.PosRate()})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].avgRet > rows[j].avgRet })
	for _, r := range rows {
		fmt.Printf("  %-14s %6d %+9.4f %6.1f%%\n", r.name, r.n, r.avgRet, r.posR*100)
	}
}

func printBTCBacktestResult(r btcBacktestResult) {
	fmt.Printf("\n  Test bars:      %d\n", r.Days)
	fmt.Printf("  Correct:        %d\n", r.Correct)
	fmt.Printf("  Wrong:          %d\n", r.Wrong)
	fmt.Printf("  Accuracy:       %.1f%%\n", r.accuracy())
	fmt.Printf("  Simulated PnL:  %+.2f pp\n", r.SumSimPnL)
	fmt.Printf("  Max Drawdown:   %.2f pp\n", r.MaxDD)
	fmt.Printf("  Best state:     %s\n", btc.StateName(r.BestState))
	fmt.Printf("  Worst state:    %s\n", btc.StateName(r.WorstState))

	fmt.Println("\n  Per-state backtest accuracy (test set, min 5 observations):")
	fmt.Printf("  %-14s %5s %5s %8s %+9s\n", "state", "n", "corr", "acc%", "simPnL")

	type row struct {
		name   string
		n      int
		corr   int
		acc    float64
		simPnL float64
	}
	var rows []row
	for s := 0; s < btc.NStates; s++ {
		sr := r.StateResults[s]
		if sr.N < 5 {
			continue
		}
		acc := 100 * float64(sr.Correct) / float64(sr.N)
		rows = append(rows, row{btc.StateName(s), sr.N, sr.Correct, acc, sr.SumPnL})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].acc > rows[j].acc })
	for _, r := range rows {
		fmt.Printf("  %-14s %5d %5d %7.1f%% %+9.2f\n", r.name, r.n, r.corr, r.acc, r.simPnL)
	}
}

func printPMMarkets(markets []btc.PMMarket) {
	fmt.Printf("\n  %-12s %-8s %-8s  %s\n", "strike", "yes_px", "no_px", "question (truncated)")
	for _, m := range markets {
		q := m.Question
		if len(q) > 55 {
			q = q[:52] + "..."
		}
		fmt.Printf("  $%-11.0f %-8.3f %-8.3f  %s\n",
			m.Strike, m.YesPrice, m.NoPrice, q)
	}
}

func analyzePMGaps(
	candles []btc.Candle,
	tm *btc.TransitionMatrix,
	retStats [btc.NStates]btc.ReturnStats,
	markets []btc.PMMarket,
) {
	if len(candles) < 2 {
		return
	}

	pred, ok := btc.PredictFromCandles(candles, tm, retStats)
	if !ok {
		fmt.Println("  (could not derive current state from candles)")
		return
	}

	fmt.Printf("\n  Current BTC state:    %s\n", pred.CurrentStateName)
	fmt.Printf("  Expected return:      %+.4f%%\n", pred.ExpectedReturn)
	fmt.Printf("  Bull prob (SURGE+UP): %.1f%%\n", pred.BullProb*100)
	fmt.Printf("  Bear prob (CRASH+DN): %.1f%%\n", pred.BearProb*100)

	// Current BTC price
	spot := candles[len(candles)-1].Close
	fmt.Printf("  Current BTC price:    $%.0f\n", spot)

	// Historical volatility (annualized)
	sigma := btc.HistoricalVolatility(candles)
	fmt.Printf("  90d annualized vol:   %.1f%%\n", sigma*100)

	// Time to expiry: Dec 31 2026
	expiry := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
	yearsToExpiry := expiry.Sub(time.Now().UTC()).Hours() / 8760.0
	fmt.Printf("  Time to expiry:       %.2f years\n", yearsToExpiry)

	// Black-Scholes gap analysis
	gaps := btc.FindBSGaps(markets, spot, sigma, yearsToExpiry, 5.0)

	if len(gaps) == 0 {
		fmt.Println("\n  No PM vs Black-Scholes gaps >5pp found.")
		return
	}

	fmt.Printf("\n  PM vs Black-Scholes gaps (minGap=5pp) — top opportunities:\n")
	fmt.Printf("  %-12s %-8s %-8s %-8s %-10s %-10s\n",
		"strike", "pm_yes", "bs_prob", "gap_pp", "edge_ratio", "direction")
	for i, g := range gaps {
		if i >= 15 {
			fmt.Printf("  ... (%d more)\n", len(gaps)-15)
			break
		}
		fmt.Printf("  $%-11.0f %-8.3f %-8.3f %+8.1f %-10.1f%% %-10s\n",
			g.Strike, g.PMPrice, g.BSProb, g.GapPP, g.EdgeRatio*100, g.Direction)
	}

	// Summary recommendation
	fmt.Println("\n  📊 BS模型建议:")
	buyYes := 0
	buyNo := 0
	for _, g := range gaps {
		if g.Direction == "BUY_YES" {
			buyYes++
		} else {
			buyNo++
		}
	}
	fmt.Printf("  BUY_YES（PM低估）: %d 个市场\n", buyYes)
	fmt.Printf("  BUY_NO（PM高估）:  %d 个市场\n", buyNo)
	if len(gaps) > 0 {
		best := gaps[0]
		fmt.Printf("  最大价差: $%.0f strike | PM=%.3f vs BS=%.3f | gap=%+.1fpp | %s\n",
			best.Strike, best.PMPrice, best.BSProb, best.GapPP, best.Direction)
	}
}
