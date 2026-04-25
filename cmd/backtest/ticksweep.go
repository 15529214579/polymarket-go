package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type tickRow struct {
	PosID string    `json:"pos_id"`
	Time  time.Time `json:"t"`
	Mid   float64   `json:"mid"`
}

type posPath struct {
	PosID    string
	EntryMid float64
	Ticks    []float64 // mid prices at 1Hz
	Duration time.Duration
	Journal  *journalTrade // nil if not matched
}

type journalTrade struct {
	ID         string  `json:"id"`
	Question   string  `json:"question"`
	Outcome    string  `json:"outcome"`
	EntryMid   float64 `json:"entry_mid"`
	ExitMid    float64 `json:"exit_mid"`
	ExitReason string  `json:"exit_reason"`
	HeldSec    int     `json:"held_sec"`
	PnlUSD     float64 `json:"pnl_usd"`
	SizeUSD    float64 `json:"size_usd"`
	Source     string  `json:"signal_source"`
	Tranche    string  `json:"tranche"`
}

type tsResult struct {
	TPPct    float64
	SLPct    float64
	Timeout  int // seconds, 0 = no timeout
	N        int
	TPHit    int
	SLHit    int
	Timeout_ int // timeout exits
	Natural  int // path ended before TP/SL/timeout
	SumPnL   float64
	Wins     int
	Losses   int
	MaxDD    float64
	AvgHold  float64 // seconds
}

func (r tsResult) avgPnL() float64 {
	if r.N == 0 {
		return 0
	}
	return r.SumPnL / float64(r.N)
}

func (r tsResult) winRate() float64 {
	if r.N == 0 {
		return 0
	}
	return 100 * float64(r.Wins) / float64(r.N)
}

func loadTickPaths(tickDir string) ([]posPath, error) {
	entries, err := os.ReadDir(tickDir)
	if err != nil {
		return nil, fmt.Errorf("read tickpath dir: %w", err)
	}
	var paths []posPath
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		posID := strings.TrimSuffix(e.Name(), ".jsonl")
		fpath := filepath.Join(tickDir, e.Name())

		f, err := os.Open(fpath)
		if err != nil {
			continue
		}
		var ticks []tickRow
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			var row tickRow
			if err := json.Unmarshal(sc.Bytes(), &row); err != nil {
				continue
			}
			if row.Mid > 0 {
				ticks = append(ticks, row)
			}
		}
		f.Close()

		if len(ticks) < 2 {
			continue
		}

		mids := make([]float64, len(ticks))
		for i, t := range ticks {
			mids[i] = t.Mid
		}
		dur := ticks[len(ticks)-1].Time.Sub(ticks[0].Time)
		paths = append(paths, posPath{
			PosID:    posID,
			EntryMid: mids[0],
			Ticks:    mids,
			Duration: dur,
		})
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i].PosID < paths[j].PosID })
	return paths, nil
}

func loadJournalTrades(journalDir string) ([]journalTrade, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, fmt.Errorf("read journal dir: %w", err)
	}
	var trades []journalTrade
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		f, err := os.Open(filepath.Join(journalDir, e.Name()))
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			var t journalTrade
			if err := json.Unmarshal(sc.Bytes(), &t); err != nil {
				continue
			}
			trades = append(trades, t)
		}
		f.Close()
	}
	return trades, nil
}

func replayPath(p posPath, tpPct, slPct float64, timeoutSec int, feeBP float64) (pnl float64, exitType string, holdSec int) {
	entry := p.EntryMid
	if entry <= 0 {
		return 0, "skip", 0
	}
	tp := tpPct / 100.0
	sl := slPct / 100.0
	fee := 2 * feeBP / 10000.0

	for i := 1; i < len(p.Ticks); i++ {
		mid := p.Ticks[i]
		ret := (mid - entry) / entry

		if timeoutSec > 0 && i >= timeoutSec {
			pnlVal := ret - fee
			return pnlVal * entry * 5 / entry, "timeout", i // normalize to $5 position
		}
		if tp > 0 && ret >= tp {
			pnlVal := tp - fee
			return pnlVal * 5, "tp", i
		}
		if sl > 0 && ret <= -sl {
			pnlVal := -sl - fee
			return pnlVal * 5, "sl", i
		}
	}
	// path ended (position closed by other means or data ends)
	lastRet := (p.Ticks[len(p.Ticks)-1] - entry) / entry
	pnlVal := lastRet - fee
	return pnlVal * 5, "natural", len(p.Ticks) - 1
}

func sweepTickPaths(paths []posPath, tpPct, slPct float64, timeoutSec int, feeBP float64) tsResult {
	res := tsResult{TPPct: tpPct, SLPct: slPct, Timeout: timeoutSec}
	equity := 0.0
	peak := 0.0
	totalHold := 0

	for _, p := range paths {
		pnl, exit, holdSec := replayPath(p, tpPct, slPct, timeoutSec, feeBP)
		res.N++
		res.SumPnL += pnl
		totalHold += holdSec

		switch exit {
		case "tp":
			res.TPHit++
		case "sl":
			res.SLHit++
		case "timeout":
			res.Timeout_++
		default:
			res.Natural++
		}

		if pnl > 0.001 {
			res.Wins++
		} else if pnl < -0.001 {
			res.Losses++
		}

		equity += pnl
		if equity > peak {
			peak = equity
		}
		if peak-equity > res.MaxDD {
			res.MaxDD = peak - equity
		}
	}

	if res.N > 0 {
		res.AvgHold = float64(totalHold) / float64(res.N)
	}
	return res
}

func runTickPathSweep(tickDir, journalDir string, feeBP float64) error {
	paths, err := loadTickPaths(tickDir)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no tickpath files found in %s", tickDir)
	}

	journal, _ := loadJournalTrades(journalDir)

	// basic stats
	totalTicks := 0
	var minDur, maxDur time.Duration
	entryBands := map[string]int{}
	for i, p := range paths {
		totalTicks += len(p.Ticks)
		if i == 0 || p.Duration < minDur {
			minDur = p.Duration
		}
		if p.Duration > maxDur {
			maxDur = p.Duration
		}
		switch {
		case p.EntryMid < 0.15:
			entryBands["<0.15"]++
		case p.EntryMid < 0.30:
			entryBands["0.15-0.30"]++
		case p.EntryMid < 0.50:
			entryBands["0.30-0.50"]++
		case p.EntryMid < 0.70:
			entryBands["0.50-0.70"]++
		default:
			entryBands["0.70+"]++
		}
	}

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" TP/SL tickpath sweep · Go daemon 真实 1Hz 路径")
	fmt.Printf(" paths: %d · ticks: %d · fee: %.0f bp/leg\n", len(paths), totalTicks, feeBP)
	fmt.Printf(" duration: %s → %s\n", minDur.Round(time.Second), maxDur.Round(time.Second))
	fmt.Println("════════════════════════════════════════")

	fmt.Println("\n── 入场价分布 ──")
	for _, band := range []string{"<0.15", "0.15-0.30", "0.30-0.50", "0.50-0.70", "0.70+"} {
		fmt.Printf("  %-12s  %d paths\n", band, entryBands[band])
	}

	// journal actual performance
	if len(journal) > 0 {
		fmt.Printf("\n── journal 真实出场（%d trades）──\n", len(journal))
		reasons := map[string]struct{ n int; pnl float64 }{}
		for _, t := range journal {
			r := reasons[t.ExitReason]
			r.n++
			r.pnl += t.PnlUSD
			reasons[t.ExitReason] = r
		}
		type kv struct {
			k string
			n int
			p float64
		}
		var sorted_ []kv
		for k, v := range reasons {
			sorted_ = append(sorted_, kv{k, v.n, v.pnl})
		}
		sort.Slice(sorted_, func(i, j int) bool { return sorted_[i].n > sorted_[j].n })
		fmt.Printf("  %-22s %5s %10s\n", "exit_reason", "n", "pnl_usd")
		for _, s := range sorted_ {
			fmt.Printf("  %-22s %5d %+10.2f\n", s.k, s.n, s.p)
		}
	}

	// current params baseline
	fmt.Println("\n── 当前参数 baseline (TP1=15% TP2=30% SL=5% timeout=4h) ──")
	baseline := sweepTickPaths(paths, 15, 5, 14400, feeBP)
	printTSRow(baseline)

	// sweep grid
	tps := []float64{5, 10, 15, 20, 25, 30, 40, 50, 75, 100}
	sls := []float64{2, 3, 5, 7, 10, 15, 20}
	timeouts := []int{0, 300, 600, 1800, 3600, 7200, 14400}

	// Phase 1: TP × SL sweep (fixed timeout 4h)
	fmt.Println("\n── TP × SL sweep (timeout=4h, top 15 by sum_pnl) ──")
	fmt.Printf("  %-5s %-5s %5s %5s %5s %5s %5s %9s %8s %8s %8s %8s\n",
		"TP%", "SL%", "n", "tp", "sl", "tout", "nat", "sum_pnl", "avg_pnl", "win%", "mdd", "avg_s")

	var allResults []tsResult
	for _, tp := range tps {
		for _, sl := range sls {
			r := sweepTickPaths(paths, tp, sl, 14400, feeBP)
			allResults = append(allResults, r)
		}
	}
	sort.Slice(allResults, func(i, j int) bool { return allResults[i].SumPnL > allResults[j].SumPnL })
	for i, r := range allResults {
		if i >= 15 {
			break
		}
		printTSRowLine(r)
	}

	fmt.Println("\n  ... bottom 5 ...")
	for i := max(0, len(allResults)-5); i < len(allResults); i++ {
		printTSRowLine(allResults[i])
	}

	// Phase 2: timeout sweep (best TP/SL from above)
	if len(allResults) > 0 {
		bestTP := allResults[0].TPPct
		bestSL := allResults[0].SLPct
		fmt.Printf("\n── timeout sweep (TP=%.0f%% SL=%.0f%%, best from above) ──\n", bestTP, bestSL)
		fmt.Printf("  %-8s %5s %5s %5s %5s %5s %9s %8s %8s %8s %8s\n",
			"timeout", "n", "tp", "sl", "tout", "nat", "sum_pnl", "avg_pnl", "win%", "mdd", "avg_s")
		for _, tout := range timeouts {
			r := sweepTickPaths(paths, bestTP, bestSL, tout, feeBP)
			label := "none"
			if tout > 0 {
				label = fmt.Sprintf("%ds", tout)
			}
			fmt.Printf("  %-8s %5d %5d %5d %5d %5d %+9.2f %+8.4f %8.1f %8.2f %8.0f\n",
				label, r.N, r.TPHit, r.SLHit, r.Timeout_, r.Natural,
				r.SumPnL, r.avgPnL(), r.winRate(), r.MaxDD, r.AvgHold)
		}
	}

	// Phase 3: entry price band analysis
	fmt.Println("\n── 按入场价带分拆 (当前参数 TP=15% SL=5%) ──")
	bands := []struct {
		name   string
		lo, hi float64
	}{
		{"<0.15", 0, 0.15},
		{"0.15-0.30", 0.15, 0.30},
		{"0.30-0.50", 0.30, 0.50},
		{"0.50-0.70", 0.50, 0.70},
		{"0.70+", 0.70, 1.01},
	}
	for _, b := range bands {
		var subset []posPath
		for _, p := range paths {
			if p.EntryMid >= b.lo && p.EntryMid < b.hi {
				subset = append(subset, p)
			}
		}
		if len(subset) == 0 {
			fmt.Printf("  %-12s  0 paths\n", b.name)
			continue
		}
		r := sweepTickPaths(subset, 15, 5, 14400, feeBP)
		fmt.Printf("  %-12s  %3d paths  %2dW/%2dL  pnl %+7.2f  win%% %5.1f  mdd %5.2f  avg_hold %4.0fs\n",
			b.name, r.N, r.Wins, r.Losses, r.SumPnL, r.winRate(), r.MaxDD, r.AvgHold)
	}

	// Phase 4: per-band optimal params
	fmt.Println("\n── 按入场价带最优参数 ──")
	for _, b := range bands {
		var subset []posPath
		for _, p := range paths {
			if p.EntryMid >= b.lo && p.EntryMid < b.hi {
				subset = append(subset, p)
			}
		}
		if len(subset) < 3 {
			continue
		}
		best := tsResult{SumPnL: -math.MaxFloat64}
		for _, tp := range tps {
			for _, sl := range sls {
				r := sweepTickPaths(subset, tp, sl, 14400, feeBP)
				if r.SumPnL > best.SumPnL {
					best = r
				}
			}
		}
		fmt.Printf("  %-12s  best TP=%.0f%% SL=%.0f%%  %dW/%dL  pnl %+.2f  win%% %.1f\n",
			b.name, best.TPPct, best.SLPct, best.Wins, best.Losses, best.SumPnL, best.winRate())
	}

	// Peak/trough analysis
	fmt.Println("\n── 路径峰值/谷值分析 ──")
	var peakRets, troughRets []float64
	for _, p := range paths {
		entry := p.EntryMid
		if entry <= 0 {
			continue
		}
		peak, trough := entry, entry
		for _, mid := range p.Ticks {
			if mid > peak {
				peak = mid
			}
			if mid < trough {
				trough = mid
			}
		}
		peakRets = append(peakRets, (peak-entry)/entry*100)
		troughRets = append(troughRets, (trough-entry)/entry*100)
	}
	if len(peakRets) > 0 {
		sort.Float64s(peakRets)
		sort.Float64s(troughRets)
		n := len(peakRets)
		fmt.Printf("  peak return:   p25=%.1f%%  p50=%.1f%%  p75=%.1f%%  max=%.1f%%\n",
			peakRets[n/4], peakRets[n/2], peakRets[3*n/4], peakRets[n-1])
		fmt.Printf("  trough return: p25=%.1f%%  p50=%.1f%%  p75=%.1f%%  min=%.1f%%\n",
			troughRets[n/4], troughRets[n/2], troughRets[3*n/4], troughRets[0])
	}

	return nil
}

func printTSRow(r tsResult) {
	fmt.Printf("  n=%d  tp=%d sl=%d tout=%d nat=%d  pnl=%+.2f  avg=%+.4f  win%%=%.1f  mdd=%.2f  avg_hold=%.0fs\n",
		r.N, r.TPHit, r.SLHit, r.Timeout_, r.Natural, r.SumPnL, r.avgPnL(), r.winRate(), r.MaxDD, r.AvgHold)
}

func printTSRowLine(r tsResult) {
	fmt.Printf("  %-5.0f %-5.0f %5d %5d %5d %5d %5d %+9.2f %+8.4f %8.1f %8.2f %8.0f\n",
		r.TPPct, r.SLPct, r.N, r.TPHit, r.SLHit, r.Timeout_, r.Natural,
		r.SumPnL, r.avgPnL(), r.winRate(), r.MaxDD, r.AvgHold)
}

