package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/15529214579/polymarket-go/internal/btc"
)

func runBTCUpDownBacktest(days *int, feeBP *float64) error {
	ctx := context.Background()
	nDays := *days
	fee := *feeBP / 10000.0

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" BTC Up/Down 1h Backtest")
	fmt.Printf(" days=%d  fee=%.0f bp\n", nDays, *feeBP)
	fmt.Println("════════════════════════════════════════")

	end := time.Now().UTC().Truncate(time.Hour)
	start := end.Add(-time.Duration(nDays) * 24 * time.Hour)

	fmt.Printf("\n── Fetching 1h candles (%d days) ──\n", nDays)
	candles, err := btc.FetchCandlesRange(ctx, "BTCUSDT", btc.Interval1h, start, end)
	if err != nil {
		return fmt.Errorf("fetch candles: %w", err)
	}
	if len(candles) < 48 {
		return fmt.Errorf("too few candles (%d)", len(candles))
	}
	fmt.Printf("  fetched %d hourly candles\n", len(candles))

	sigmaHist := btc.HistoricalVolatility(candles)
	sigmaEWMA := btc.EWMAVolatility(candles, 0.94)
	fmt.Printf("  vol: hist=%.1f%% ewma=%.1f%%\n", sigmaHist*100, sigmaEWMA*100)

	// HMM regime analysis
	regime, rName, rConf := btc.DetectCurrentRegime(candles)
	_ = regime
	fmt.Printf("  HMM regime: %s (conf=%.0f%%)\n", rName, rConf*100)

	// Split: first 67% train, last 33% test
	splitIdx := int(float64(len(candles)) * 0.67)
	if splitIdx < 100 {
		splitIdx = 100
	}
	trainCandles := candles[:splitIdx]
	testCandles := candles[splitIdx:]
	fmt.Printf("  train: %d bars  test: %d bars\n", len(trainCandles), len(testCandles))

	// Train Markov on training set
	tm, _ := btc.Train(trainCandles)
	retStats := btc.BuildReturnStats(trainCandles)

	// Also train HMM on training set
	trainObs := btc.CandlesToObservations(trainCandles)
	hmmModel := btc.TrainHMM(trainObs, 20)

	// Backtest: simulate 1h Up/Down bets on test set
	type betResult struct {
		Hour      int
		Predicted string
		Actual    string
		Correct   bool
		PnL       float64
		PMPrice   float64
		Regime    int
	}

	var results []betResult
	var equity, peak, maxDD float64
	var totalPnL float64
	wins, losses := 0, 0
	byHour := make(map[int]*struct{ n, w int })
	byRegime := make(map[int]*struct{ n, w int })

	for i := 1; i+1 < len(testCandles); i++ {
		// Predict direction using Markov
		state, ok := btc.CandleState(testCandles, i)
		if !ok {
			continue
		}
		pred := btc.Predict(state, &tm, retStats)
		predDir := "Up"
		if pred.BearProb > pred.BullProb {
			predDir = "Down"
		}

		// Confidence filter
		conf := pred.BullProb
		if predDir == "Down" {
			conf = pred.BearProb
		}
		if conf < 0.40 {
			continue
		}

		// HMM regime filter
		testObs := btc.CandlesToObservations(testCandles[:i+1])
		var curRegime int
		if len(testObs) > 0 {
			path := btc.Viterbi(hmmModel, testObs)
			if len(path) > 0 {
				curRegime = path[len(path)-1]
			}
		}

		// Check actual result
		curr := testCandles[i]
		next := testCandles[i+1]
		if curr.Close <= 0 || next.Close <= 0 {
			continue
		}
		actualDir := "Down"
		if next.Close >= next.Open {
			actualDir = "Up"
		}

		correct := predDir == actualDir
		pmPrice := 0.50 // assume 50/50 market (no PM data in backtest)
		var pnl float64
		if correct {
			if pmPrice > 0 {
				pnl = (1.0/pmPrice-1.0)*5.0 - fee*5.0
			}
		} else {
			pnl = -5.0 - fee*5.0
		}

		if correct {
			wins++
		} else {
			losses++
		}
		totalPnL += pnl
		equity += pnl
		if equity > peak {
			peak = equity
		}
		dd := peak - equity
		if dd > maxDD {
			maxDD = dd
		}

		hr := curr.Timestamp.Hour()
		if _, ok := byHour[hr]; !ok {
			byHour[hr] = &struct{ n, w int }{}
		}
		byHour[hr].n++
		if correct {
			byHour[hr].w++
		}

		if _, ok := byRegime[curRegime]; !ok {
			byRegime[curRegime] = &struct{ n, w int }{}
		}
		byRegime[curRegime].n++
		if correct {
			byRegime[curRegime].w++
		}

		results = append(results, betResult{
			Hour:      hr,
			Predicted: predDir,
			Actual:    actualDir,
			Correct:   correct,
			PnL:       pnl,
			PMPrice:   pmPrice,
			Regime:    curRegime,
		})
	}

	// Calculate metrics
	total := wins + losses
	winRate := 0.0
	if total > 0 {
		winRate = float64(wins) / float64(total) * 100
	}

	// Sharpe ratio (annualized, assume 16 bets/day)
	var sumRet, sumRetSq float64
	for _, r := range results {
		sumRet += r.PnL
		sumRetSq += r.PnL * r.PnL
	}
	nBets := float64(len(results))
	meanRet := 0.0
	stdRet := 0.0
	sharpe := 0.0
	if nBets > 1 {
		meanRet = sumRet / nBets
		variance := sumRetSq/nBets - meanRet*meanRet
		if variance > 0 {
			stdRet = math.Sqrt(variance)
			sharpe = (meanRet / stdRet) * math.Sqrt(16*365) // annualized
		}
	}

	// Calmar ratio
	calmar := 0.0
	if maxDD > 0 {
		annualPnL := totalPnL / float64(nDays) * 365
		calmar = annualPnL / maxDD
	}

	// Print results
	fmt.Println("\n── Results ──")
	fmt.Printf("  Bets:           %d\n", total)
	fmt.Printf("  Wins:           %d\n", wins)
	fmt.Printf("  Losses:         %d\n", losses)
	fmt.Printf("  Win Rate:       %.1f%%\n", winRate)
	fmt.Printf("  Total PnL:      $%+.2f\n", totalPnL)
	fmt.Printf("  Max Drawdown:   $%.2f\n", maxDD)
	fmt.Printf("  Sharpe (ann):   %.2f\n", sharpe)
	fmt.Printf("  Calmar:         %.2f\n", calmar)
	fmt.Printf("  Avg PnL/bet:    $%+.3f\n", meanRet)
	fmt.Printf("  Std PnL/bet:    $%.3f\n", stdRet)

	// By hour
	fmt.Println("\n── By Hour (UTC) ──")
	fmt.Printf("  %-6s %5s %5s %8s\n", "hour", "bets", "wins", "winrate")
	for hr := 0; hr < 24; hr++ {
		h, ok := byHour[hr]
		if !ok || h.n == 0 {
			continue
		}
		wr := float64(h.w) / float64(h.n) * 100
		marker := ""
		if wr >= 55 {
			marker = " ✅"
		} else if wr < 45 {
			marker = " ❌"
		}
		fmt.Printf("  %-6d %5d %5d %7.1f%%%s\n", hr, h.n, h.w, wr, marker)
	}

	// By regime
	fmt.Println("\n── By HMM Regime ──")
	fmt.Printf("  %-14s %5s %5s %8s\n", "regime", "bets", "wins", "winrate")
	for r := 0; r < btc.NRegimes; r++ {
		rg, ok := byRegime[r]
		if !ok || rg.n == 0 {
			continue
		}
		wr := float64(rg.w) / float64(rg.n) * 100
		fmt.Printf("  %-14s %5d %5d %7.1f%%\n", btc.RegimeName(r), rg.n, rg.w, wr)
	}

	// Regime-filtered strategy
	fmt.Println("\n── Regime-Filtered Strategy (TREND only) ──")
	var filtWins, filtLosses int
	var filtPnL float64
	for _, r := range results {
		if r.Regime != btc.RegimeTrend {
			continue
		}
		if r.Correct {
			filtWins++
		} else {
			filtLosses++
		}
		filtPnL += r.PnL
	}
	filtTotal := filtWins + filtLosses
	filtWR := 0.0
	if filtTotal > 0 {
		filtWR = float64(filtWins) / float64(filtTotal) * 100
	}
	fmt.Printf("  Bets:     %d (vs %d unfiltered)\n", filtTotal, total)
	fmt.Printf("  Win Rate: %.1f%% (vs %.1f%%)\n", filtWR, winRate)
	fmt.Printf("  PnL:      $%+.2f (vs $%+.2f)\n", filtPnL, totalPnL)

	return nil
}
