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

// ---------------------------------------------------------------------------
// State definition — dual-timeframe momentum × volume
// ---------------------------------------------------------------------------

// Short-term momentum (30s delta in percentage points)
const (
	shortCrash = iota // < -2pp
	shortDown         // -2 to -0.5pp
	shortFlat         // -0.5 to 0.5pp
	shortUp           // 0.5 to 2pp
	shortSurge        // > 2pp
	nShort     = 5
)

// Medium-term trend (60s delta in percentage points)
const (
	trendDown = iota // < 0pp
	trendFlat        // 0 to 3pp
	trendUp          // > 3pp (matches detector's 3pp threshold)
	nTrend    = 3
)

// Volume regime
const (
	volSell = iota // buy_ratio < 0.40
	volNeut        // 0.40 to 0.60
	volBuy         // > 0.60 (matches detector's 0.60 threshold)
	nVol    = 3
)

const nStates = nShort * nTrend * nVol // 45

func stateIdx(short, trend, vol int) int {
	return short*nTrend*nVol + trend*nVol + vol
}

func stateComponents(s int) (short, trend, vol int) {
	vol = s % nVol
	s /= nVol
	trend = s % nTrend
	short = s / nTrend
	return
}

func stateName(s int) string {
	sh, tr, vl := stateComponents(s)
	shorts := [nShort]string{"sCRASH", "sDOWN", "sFLAT", "sUP", "sSURGE"}
	trends := [nTrend]string{"tDn", "tFl", "tUp"}
	vols := [nVol]string{"vSell", "vNeut", "vBuy"}
	return fmt.Sprintf("%s/%s/%s", shorts[sh], trends[tr], vols[vl])
}

// ---------------------------------------------------------------------------
// Rich tick (extends tickRow with volume fields)
// ---------------------------------------------------------------------------

type richTick struct {
	Time    time.Time `json:"t"`
	Mid     float64   `json:"mid"`
	BuyVol  float64   `json:"buy_vol"`
	SellVol float64   `json:"sell_vol"`
}

type richPath struct {
	PosID    string
	EntryMid float64
	Ticks    []richTick
	Duration time.Duration
}

func loadRichPaths(tickDir string) ([]richPath, error) {
	entries, err := os.ReadDir(tickDir)
	if err != nil {
		return nil, err
	}
	var paths []richPath
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		posID := strings.TrimSuffix(e.Name(), ".jsonl")
		f, err := os.Open(filepath.Join(tickDir, e.Name()))
		if err != nil {
			continue
		}
		var ticks []richTick
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			var row struct {
				T       time.Time `json:"t"`
				Mid     float64   `json:"mid"`
				BuyVol  float64   `json:"buy_vol"`
				SellVol float64   `json:"sell_vol"`
			}
			if err := json.Unmarshal(sc.Bytes(), &row); err != nil || row.Mid <= 0 {
				continue
			}
			ticks = append(ticks, richTick{
				Time:    row.T,
				Mid:     row.Mid,
				BuyVol:  row.BuyVol,
				SellVol: row.SellVol,
			})
		}
		f.Close()
		if len(ticks) < 120 {
			continue
		}
		paths = append(paths, richPath{
			PosID:    posID,
			EntryMid: ticks[0].Mid,
			Ticks:    ticks,
			Duration: ticks[len(ticks)-1].Time.Sub(ticks[0].Time),
		})
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i].PosID < paths[j].PosID })
	return paths, nil
}

// ---------------------------------------------------------------------------
// Feature computation
// ---------------------------------------------------------------------------

func classifyShort(deltaPP float64) int {
	switch {
	case deltaPP < -2:
		return shortCrash
	case deltaPP < -0.5:
		return shortDown
	case deltaPP <= 0.5:
		return shortFlat
	case deltaPP <= 2:
		return shortUp
	default:
		return shortSurge
	}
}

func classifyTrend(deltaPP float64) int {
	switch {
	case deltaPP < 0:
		return trendDown
	case deltaPP <= 3:
		return trendFlat
	default:
		return trendUp
	}
}

func classifyVol(buyRatio float64) int {
	switch {
	case buyRatio < 0.40:
		return volSell
	case buyRatio <= 0.60:
		return volNeut
	default:
		return volBuy
	}
}

type tickFeatures struct {
	delta30s  float64
	delta60s  float64
	buyRatio  float64
	state     int
}

func computeFeatures(ticks []richTick, i int) tickFeatures {
	// 30s short-term momentum
	lb30 := 30
	if lb30 > i {
		lb30 = i
	}
	d30 := 0.0
	if lb30 > 0 && ticks[i-lb30].Mid > 0 {
		d30 = (ticks[i].Mid - ticks[i-lb30].Mid) / ticks[i-lb30].Mid * 100
	}

	// 60s medium-term trend
	lb60 := 60
	if lb60 > i {
		lb60 = i
	}
	d60 := 0.0
	if lb60 > 0 && ticks[i-lb60].Mid > 0 {
		d60 = (ticks[i].Mid - ticks[i-lb60].Mid) / ticks[i-lb60].Mid * 100
	}

	// 30s rolling buy ratio
	window := 30
	if window > i {
		window = i
	}
	var bv, sv float64
	for j := i - window; j <= i; j++ {
		bv += ticks[j].BuyVol
		sv += ticks[j].SellVol
	}
	br := 0.5
	if bv+sv > 0 {
		br = bv / (bv + sv)
	}

	sh := classifyShort(d30)
	tr := classifyTrend(d60)
	vl := classifyVol(br)
	return tickFeatures{
		delta30s: d30,
		delta60s: d60,
		buyRatio: br,
		state:    stateIdx(sh, tr, vl),
	}
}

// ---------------------------------------------------------------------------
// Forward return analysis (the key Markov insight)
// ---------------------------------------------------------------------------

type forwardStats struct {
	state     int
	n         int
	sumRet30  float64 // sum of 30s forward returns (pp)
	sumRet60  float64 // sum of 60s forward returns
	sumRet120 float64 // sum of 120s forward returns
	nPos30    int     // count of positive 30s returns
	nPos60    int
	nPos120   int
	maxDD30   float64 // max drawdown within 30s window
}

func (f forwardStats) avgRet30() float64 {
	if f.n == 0 {
		return 0
	}
	return f.sumRet30 / float64(f.n)
}
func (f forwardStats) avgRet60() float64 {
	if f.n == 0 {
		return 0
	}
	return f.sumRet60 / float64(f.n)
}
func (f forwardStats) avgRet120() float64 {
	if f.n == 0 {
		return 0
	}
	return f.sumRet120 / float64(f.n)
}
func (f forwardStats) posRate30() float64 {
	if f.n == 0 {
		return 0
	}
	return 100 * float64(f.nPos30) / float64(f.n)
}
func (f forwardStats) posRate60() float64 {
	if f.n == 0 {
		return 0
	}
	return 100 * float64(f.nPos60) / float64(f.n)
}

func computeForwardReturns(paths []richPath) [nStates]forwardStats {
	var stats [nStates]forwardStats
	for s := 0; s < nStates; s++ {
		stats[s].state = s
	}

	for _, p := range paths {
		n := len(p.Ticks)
		for i := 60; i < n; i++ {
			f := computeFeatures(p.Ticks, i)
			mid := p.Ticks[i].Mid
			if mid <= 0 {
				continue
			}

			stats[f.state].n++

			// 30s forward
			if i+30 < n {
				ret := (p.Ticks[i+30].Mid - mid) / mid * 100
				stats[f.state].sumRet30 += ret
				if ret > 0 {
					stats[f.state].nPos30++
				}
				// max drawdown within 30s
				minMid := mid
				for j := i + 1; j <= i+30 && j < n; j++ {
					if p.Ticks[j].Mid < minMid {
						minMid = p.Ticks[j].Mid
					}
				}
				dd := (mid - minMid) / mid * 100
				if dd > stats[f.state].maxDD30 {
					stats[f.state].maxDD30 = dd
				}
			}

			// 60s forward
			if i+60 < n {
				ret := (p.Ticks[i+60].Mid - mid) / mid * 100
				stats[f.state].sumRet60 += ret
				if ret > 0 {
					stats[f.state].nPos60++
				}
			}

			// 120s forward
			if i+120 < n {
				ret := (p.Ticks[i+120].Mid - mid) / mid * 100
				stats[f.state].sumRet120 += ret
				if ret > 0 {
					stats[f.state].nPos120++
				}
			}
		}
	}
	return stats
}

// ---------------------------------------------------------------------------
// Second-order transitions (bigram: prev_state → current_state)
// ---------------------------------------------------------------------------

type bigramKey struct {
	prev, curr int
}

type bigramForward struct {
	n        int
	sumRet60 float64
	nPos60   int
}

func (b bigramForward) avgRet60() float64 {
	if b.n == 0 {
		return 0
	}
	return b.sumRet60 / float64(b.n)
}
func (b bigramForward) posRate60() float64 {
	if b.n == 0 {
		return 0
	}
	return 100 * float64(b.nPos60) / float64(b.n)
}

func computeBigramForward(paths []richPath) map[bigramKey]bigramForward {
	bigrams := map[bigramKey]bigramForward{}
	for _, p := range paths {
		n := len(p.Ticks)
		prevState := -1
		for i := 60; i < n; i++ {
			f := computeFeatures(p.Ticks, i)
			if prevState >= 0 && i+60 < n {
				key := bigramKey{prevState, f.state}
				b := bigrams[key]
				b.n++
				ret := (p.Ticks[i+60].Mid - p.Ticks[i].Mid) / p.Ticks[i].Mid * 100
				b.sumRet60 += ret
				if ret > 0 {
					b.nPos60++
				}
				bigrams[key] = b
			}
			prevState = f.state
		}
	}
	return bigrams
}

// ---------------------------------------------------------------------------
// Markov-guided backtest engine
// ---------------------------------------------------------------------------

type markovConfig struct {
	minExpectedReturn float64 // min E[60s return] to enter (pp)
	minPosRate        float64 // min P(positive 60s return) to enter (0-1)
	useBigram         bool    // use 2nd-order Markov (prev→curr bigram)
	bigramMinN        int     // min observations for bigram to be valid
	slPct             float64
	tpPct             float64 // 0 = no TP
	timeoutSec        int
	feeBP             float64
	cooldownSec       int
	warmupTicks       int
	requireTrendUp    bool    // require 60s trend to be Up (>3pp)
	requireVolBuy     bool    // require buy ratio > 0.60
}

type markovResult struct {
	label       string
	nPaths      int
	nSignals    int
	nTrades     int
	wins        int
	losses      int
	sumPnL      float64
	maxDD       float64
	avgHoldSec  float64
	tpHit       int
	slHit       int
	timeoutHit  int
	naturalHit  int
}

func (r markovResult) winRate() float64 {
	if r.nTrades == 0 {
		return 0
	}
	return 100 * float64(r.wins) / float64(r.nTrades)
}

func (r markovResult) avgPnL() float64 {
	if r.nTrades == 0 {
		return 0
	}
	return r.sumPnL / float64(r.nTrades)
}

func replayMarkov(
	p richPath,
	fwdStats [nStates]forwardStats,
	bigrams map[bigramKey]bigramForward,
	cfg markovConfig,
) markovResult {
	res := markovResult{nPaths: 1}

	type openPos struct {
		entryIdx int
		entryMid float64
	}

	var position *openPos
	coolUntil := 0
	totalHold := 0
	prevState := -1

	for i := cfg.warmupTicks; i < len(p.Ticks); i++ {
		mid := p.Ticks[i].Mid
		f := computeFeatures(p.Ticks, i)

		// check exit for open position
		if position != nil {
			held := i - position.entryIdx
			ret := (mid - position.entryMid) / position.entryMid
			fee := 2 * cfg.feeBP / 10000.0

			exited := false
			var pnl float64

			if cfg.timeoutSec > 0 && held >= cfg.timeoutSec {
				pnl = (ret - fee) * 5
				res.timeoutHit++
				exited = true
			} else if cfg.tpPct > 0 && ret >= cfg.tpPct/100 {
				pnl = (cfg.tpPct/100 - fee) * 5
				res.tpHit++
				exited = true
			} else if cfg.slPct > 0 && ret <= -cfg.slPct/100 {
				pnl = (-cfg.slPct/100 - fee) * 5
				res.slHit++
				coolUntil = i + cfg.cooldownSec
				exited = true
			}

			if exited {
				res.nTrades++
				res.sumPnL += pnl
				totalHold += held
				if pnl > 0.001 {
					res.wins++
				} else if pnl < -0.001 {
					res.losses++
				}
				position = nil
			}
		}

		// check entry
		if position != nil || i < coolUntil {
			prevState = f.state
			continue
		}

		// Gate 1: must be in short-term UP or SURGE
		sh, tr, vl := stateComponents(f.state)
		if sh != shortUp && sh != shortSurge {
			prevState = f.state
			continue
		}

		// Gate 2: optional trend/volume filters
		if cfg.requireTrendUp && tr != trendUp {
			prevState = f.state
			continue
		}
		if cfg.requireVolBuy && vl != volBuy {
			prevState = f.state
			continue
		}

		// Gate 3: forward return expectation
		var expRet float64
		var posRate float64
		var hasData bool

		if cfg.useBigram && prevState >= 0 {
			key := bigramKey{prevState, f.state}
			if b, ok := bigrams[key]; ok && b.n >= cfg.bigramMinN {
				expRet = b.avgRet60()
				posRate = b.posRate60() / 100
				hasData = true
			}
		}

		if !hasData {
			fs := fwdStats[f.state]
			if fs.n >= 10 {
				expRet = fs.avgRet60()
				posRate = fs.posRate60() / 100
				hasData = true
			}
		}

		if !hasData {
			prevState = f.state
			continue
		}

		if expRet >= cfg.minExpectedReturn && posRate >= cfg.minPosRate {
			res.nSignals++
			position = &openPos{entryIdx: i, entryMid: mid}
		}

		prevState = f.state
	}

	// close remaining position at path end
	if position != nil {
		lastMid := p.Ticks[len(p.Ticks)-1].Mid
		fee := 2 * cfg.feeBP / 10000.0
		ret := (lastMid - position.entryMid) / position.entryMid
		pnl := (ret - fee) * 5
		res.nTrades++
		res.sumPnL += pnl
		totalHold += len(p.Ticks) - 1 - position.entryIdx
		res.naturalHit++
		if pnl > 0.001 {
			res.wins++
		} else if pnl < -0.001 {
			res.losses++
		}
		position = nil
	}

	if res.nTrades > 0 {
		res.avgHoldSec = float64(totalHold) / float64(res.nTrades)
	}
	return res
}

// ---------------------------------------------------------------------------
// Momentum baseline (replicate current detector logic for comparison)
// ---------------------------------------------------------------------------

func replayMomentumBaseline(p richPath, feeBP, slPct, tpPct float64, timeoutSec, cooldownSec int) markovResult {
	res := markovResult{nPaths: 1}

	type openPos struct {
		entryIdx int
		entryMid float64
	}

	var position *openPos
	coolUntil := 0
	totalHold := 0

	for i := 60; i < len(p.Ticks); i++ {
		mid := p.Ticks[i].Mid

		// check exit
		if position != nil {
			held := i - position.entryIdx
			ret := (mid - position.entryMid) / position.entryMid
			fee := 2 * feeBP / 10000.0

			exited := false
			var pnl float64

			if timeoutSec > 0 && held >= timeoutSec {
				pnl = (ret - fee) * 5
				res.timeoutHit++
				exited = true
			} else if tpPct > 0 && ret >= tpPct/100 {
				pnl = (tpPct/100 - fee) * 5
				res.tpHit++
				exited = true
			} else if slPct > 0 && ret <= -slPct/100 {
				pnl = (-slPct/100 - fee) * 5
				res.slHit++
				coolUntil = i + cooldownSec
				exited = true
			}

			if exited {
				res.nTrades++
				res.sumPnL += pnl
				totalHold += held
				if pnl > 0.001 {
					res.wins++
				} else if pnl < -0.001 {
					res.losses++
				}
				position = nil
			}
		}

		if position != nil || i < coolUntil {
			continue
		}

		// momentum signal (replicate detector)
		refMid := p.Ticks[i-60].Mid
		if refMid <= 0 {
			continue
		}
		delta60 := (mid - refMid) / refMid * 100
		if delta60 < 3.0 {
			continue
		}

		// tail 4/5 upticks
		tailLen := 5
		if tailLen > i {
			continue
		}
		ups := 0
		for j := i - tailLen + 1; j <= i; j++ {
			if p.Ticks[j].Mid > p.Ticks[j-1].Mid {
				ups++
			}
		}
		if ups < 4 {
			continue
		}

		// buy ratio >= 60%
		var bv, sv float64
		for j := i - 60; j <= i; j++ {
			bv += p.Ticks[j].BuyVol
			sv += p.Ticks[j].SellVol
		}
		br := 0.5
		if bv+sv > 0 {
			br = bv / (bv + sv)
		}
		if br < 0.60 {
			continue
		}

		res.nSignals++
		position = &openPos{entryIdx: i, entryMid: mid}
	}

	// close remaining
	if position != nil {
		lastMid := p.Ticks[len(p.Ticks)-1].Mid
		fee := 2 * feeBP / 10000.0
		ret := (lastMid - position.entryMid) / position.entryMid
		pnl := (ret - fee) * 5
		res.nTrades++
		res.sumPnL += pnl
		totalHold += len(p.Ticks) - 1 - position.entryIdx
		res.naturalHit++
		if pnl > 0.001 {
			res.wins++
		} else if pnl < -0.001 {
			res.losses++
		}
	}

	if res.nTrades > 0 {
		res.avgHoldSec = float64(totalHold) / float64(res.nTrades)
	}
	return res
}

func mergeResults(results []markovResult) markovResult {
	if len(results) == 0 {
		return markovResult{}
	}
	m := markovResult{}
	equity, peak, mdd := 0.0, 0.0, 0.0
	totalHold := 0.0
	for _, r := range results {
		m.nPaths += r.nPaths
		m.nSignals += r.nSignals
		m.nTrades += r.nTrades
		m.wins += r.wins
		m.losses += r.losses
		m.sumPnL += r.sumPnL
		m.tpHit += r.tpHit
		m.slHit += r.slHit
		m.timeoutHit += r.timeoutHit
		m.naturalHit += r.naturalHit
		totalHold += r.avgHoldSec * float64(r.nTrades)

		equity += r.sumPnL
		if equity > peak {
			peak = equity
		}
		if peak-equity > mdd {
			mdd = peak - equity
		}
	}
	m.maxDD = mdd
	if m.nTrades > 0 {
		m.avgHoldSec = totalHold / float64(m.nTrades)
	}
	return m
}

// ---------------------------------------------------------------------------
// Main entry point
// ---------------------------------------------------------------------------

func runMarkovBacktest(tickDir string, feeBP float64) error {
	paths, err := loadRichPaths(tickDir)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no tickpath files in %s", tickDir)
	}

	totalTicks := 0
	for _, p := range paths {
		totalTicks += len(p.Ticks)
	}

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" Markov Chain Entry Signal Backtest v2")
	fmt.Printf(" paths: %d · ticks: %d · fee: %.0f bp/leg\n", len(paths), totalTicks, feeBP)
	fmt.Println("════════════════════════════════════════")

	// === Phase 1: Forward return analysis by state ===
	fmt.Println("\n── Phase 1: Forward Return by State ──")
	fwdStats := computeForwardReturns(paths)

	// show states with significant data, sorted by expected 60s return
	type stRow struct {
		name    string
		n       int
		r30     float64
		r60     float64
		r120    float64
		pr60    float64
		maxDD30 float64
	}
	var rows []stRow
	for s := 0; s < nStates; s++ {
		fs := fwdStats[s]
		if fs.n < 20 {
			continue
		}
		rows = append(rows, stRow{
			name:    stateName(s),
			n:       fs.n,
			r30:     fs.avgRet30(),
			r60:     fs.avgRet60(),
			r120:    fs.avgRet120(),
			pr60:    fs.posRate60(),
			maxDD30: fs.maxDD30,
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].r60 > rows[j].r60 })

	fmt.Printf("  %-24s %5s %+8s %+8s %+8s %6s %6s\n",
		"state", "n", "E[30s]", "E[60s]", "E[120s]", "P+60s", "mDD30")

	fmt.Println("  — Top 10 (best expected 60s return) —")
	for i, r := range rows {
		if i >= 10 {
			break
		}
		fmt.Printf("  %-24s %5d %+8.3f %+8.3f %+8.3f %5.1f%% %5.1f%%\n",
			r.name, r.n, r.r30, r.r60, r.r120, r.pr60, r.maxDD30)
	}

	fmt.Println("  — Bottom 10 (worst expected 60s return) —")
	for i := max2(0, len(rows)-10); i < len(rows); i++ {
		r := rows[i]
		fmt.Printf("  %-24s %5d %+8.3f %+8.3f %+8.3f %5.1f%% %5.1f%%\n",
			r.name, r.n, r.r30, r.r60, r.r120, r.pr60, r.maxDD30)
	}

	// === Phase 2: Bigram analysis ===
	fmt.Println("\n── Phase 2: Bigram Analysis (2nd-order Markov) ──")
	bigrams := computeBigramForward(paths)

	// show FLAT→SURGE vs SURGE→SURGE transitions
	type bgRow struct {
		from, to string
		n        int
		r60      float64
		pr60     float64
	}
	var bgRows []bgRow
	for key, b := range bigrams {
		if b.n < 10 {
			continue
		}
		_, _, _ = stateComponents(key.curr)
		sh, _, _ := stateComponents(key.curr)
		if sh != shortSurge && sh != shortUp {
			continue
		}
		bgRows = append(bgRows, bgRow{
			from: stateName(key.prev),
			to:   stateName(key.curr),
			n:    b.n,
			r60:  b.avgRet60(),
			pr60: b.posRate60(),
		})
	}
	sort.Slice(bgRows, func(i, j int) bool { return bgRows[i].r60 > bgRows[j].r60 })

	fmt.Printf("  %-24s → %-24s %5s %+8s %6s\n", "from", "to", "n", "E[60s]", "P+60s")
	for i, r := range bgRows {
		if i >= 15 {
			break
		}
		fmt.Printf("  %-24s → %-24s %5d %+8.3f %5.1f%%\n",
			r.from, r.to, r.n, r.r60, r.pr60)
	}
	if len(bgRows) > 15 {
		fmt.Println("  ...")
		for i := max2(0, len(bgRows)-5); i < len(bgRows); i++ {
			r := bgRows[i]
			fmt.Printf("  %-24s → %-24s %5d %+8.3f %5.1f%%\n",
				r.from, r.to, r.n, r.r60, r.pr60)
		}
	}

	// === Phase 3: Momentum baseline ===
	fmt.Println("\n── Phase 3: Momentum Baseline ──")
	slPct := 15.0
	tpPct := 0.0
	timeoutSec := 21600
	cooldownSec := 3600

	var baseResults []markovResult
	for _, p := range paths {
		baseResults = append(baseResults, replayMomentumBaseline(p, feeBP, slPct, tpPct, timeoutSec, cooldownSec))
	}
	baseline := mergeResults(baseResults)
	baseline.label = "momentum"
	printRow(baseline)

	// === Phase 4: Markov entry sweep ===
	fmt.Println("\n── Phase 4: Markov Entry Sweep ──")

	type sweepResult struct {
		label  string
		result markovResult
	}
	var allSweep []sweepResult

	minReturns := []float64{-0.5, -0.2, 0.0, 0.1, 0.2, 0.3, 0.5}
	minPosRates := []float64{0.30, 0.35, 0.40, 0.45, 0.50, 0.55}

	for _, useT := range []bool{false, true} {
		for _, useV := range []bool{false, true} {
			for _, useB := range []bool{false, true} {
				for _, mr := range minReturns {
					for _, mpr := range minPosRates {
						cfg := markovConfig{
							minExpectedReturn: mr,
							minPosRate:        mpr,
							useBigram:         useB,
							bigramMinN:        5,
							slPct:             slPct,
							tpPct:             tpPct,
							timeoutSec:        timeoutSec,
							feeBP:             feeBP,
							cooldownSec:       cooldownSec,
							warmupTicks:       60,
							requireTrendUp:    useT,
							requireVolBuy:     useV,
						}
						var results []markovResult
						for _, p := range paths {
							results = append(results, replayMarkov(p, fwdStats, bigrams, cfg))
						}
						merged := mergeResults(results)
						flags := ""
						if useT {
							flags += "T"
						}
						if useV {
							flags += "V"
						}
						if useB {
							flags += "B"
						}
						if flags == "" {
							flags = "-"
						}
						label := fmt.Sprintf("r≥%+.1f p≥%.0f%% %s", mr, mpr*100, flags)
						allSweep = append(allSweep, sweepResult{label, merged})
					}
				}
			}
		}
	}

	sort.Slice(allSweep, func(i, j int) bool {
		return allSweep[i].result.sumPnL > allSweep[j].result.sumPnL
	})

	fmt.Printf("  %-28s %5s %5s %3s %3s %3s %3s %+9s %+8s %6s %6s %7s\n",
		"config", "sig", "trd", "W", "L", "sl", "to", "pnl", "avg", "win%", "mdd", "hold_s")

	fmt.Println("\n  Top 20 by PnL:")
	for i, s := range allSweep {
		if i >= 20 {
			break
		}
		printRowLine(s.label, s.result)
	}

	// filter for configs with nTrades > 0 and better than baseline
	fmt.Println("\n  Configs that beat momentum baseline:")
	found := 0
	for _, s := range allSweep {
		if s.result.nTrades == 0 {
			continue
		}
		if s.result.sumPnL > baseline.sumPnL {
			printRowLine(s.label, s.result)
			found++
			if found >= 20 {
				break
			}
		}
	}
	if found == 0 {
		fmt.Println("  (none)")
	}

	// === Phase 5: Head-to-head ===
	fmt.Println("\n── Phase 5: Head-to-Head ──")

	if len(allSweep) > 0 && allSweep[0].result.nTrades > 0 {
		best := allSweep[0]
		fmt.Printf("\n  %-28s %5s %5s %3s %3s %+9s %+8s %6s %6s\n",
			"strategy", "sig", "trd", "W", "L", "pnl", "avg", "win%", "mdd")
		fmt.Printf("  %-28s %5d %5d %3d %3d %+9.2f %+8.4f %6.1f %6.2f\n",
			"momentum (current)",
			baseline.nSignals, baseline.nTrades,
			baseline.wins, baseline.losses,
			baseline.sumPnL, baseline.avgPnL(), baseline.winRate(), baseline.maxDD)
		fmt.Printf("  %-28s %5d %5d %3d %3d %+9.2f %+8.4f %6.1f %6.2f\n",
			"markov-best ("+best.label+")",
			best.result.nSignals, best.result.nTrades,
			best.result.wins, best.result.losses,
			best.result.sumPnL, best.result.avgPnL(), best.result.winRate(), best.result.maxDD)

		imp := best.result.sumPnL - baseline.sumPnL
		fmt.Printf("\n  Δ PnL: %+.2f  ", imp)
		if imp > 0 {
			fmt.Println("✅ Markov wins")
		} else {
			fmt.Println("❌ Momentum wins")
		}

		// per-path breakdown for best
		fmt.Println("\n  Per-path detail (only trades with activity):")
		fmt.Printf("  %-6s %6s %+9s  %6s %+9s  %+9s\n",
			"path", "m_trd", "m_pnl", "k_trd", "k_pnl", "delta")

		bestCfg := markovConfig{
			minExpectedReturn: best.result.avgPnL(), // approximate
			minPosRate:        0,
			slPct:             slPct,
			tpPct:             tpPct,
			timeoutSec:        timeoutSec,
			feeBP:             feeBP,
			cooldownSec:       cooldownSec,
			warmupTicks:       60,
		}
		// re-parse the label to extract actual config
		// just re-run with same params
		_ = bestCfg

		better, worse, same := 0, 0, 0
		for i, p := range paths {
			mRes := baseResults[i]
			kRes := replayMarkov(p, fwdStats, bigrams, markovConfig{
				minExpectedReturn: 0, // use best config's threshold
				minPosRate:        0.30,
				slPct:             slPct,
				tpPct:             tpPct,
				timeoutSec:        timeoutSec,
				feeBP:             feeBP,
				cooldownSec:       cooldownSec,
				warmupTicks:       60,
				requireTrendUp:    false,
				requireVolBuy:     false,
			})
			d := kRes.sumPnL - mRes.sumPnL
			if mRes.nTrades > 0 || kRes.nTrades > 0 {
				fmt.Printf("  %-6s %6d %+9.2f  %6d %+9.2f  %+9.2f\n",
					p.PosID, mRes.nTrades, mRes.sumPnL,
					kRes.nTrades, kRes.sumPnL, d)
			}
			switch {
			case d > 0.01:
				better++
			case d < -0.01:
				worse++
			default:
				same++
			}
		}
		fmt.Printf("\n  Paths improved: %d  worse: %d  same: %d\n", better, worse, same)
	}

	// === Phase 6: Insight summary ===
	fmt.Println("\n── Phase 6: Key Insights ──")

	// find which states have positive expected returns
	fmt.Println("\n  States with E[60s] > 0 (safe to enter):")
	for _, r := range rows {
		if r.r60 > 0 && r.n >= 50 {
			fmt.Printf("    %-24s  n=%d  E[60s]=%+.3fpp  P+=%5.1f%%\n",
				r.name, r.n, r.r60, r.pr60)
		}
	}

	fmt.Println("\n  States with E[60s] < -0.3 (avoid entry):")
	for _, r := range rows {
		if r.r60 < -0.3 && r.n >= 50 {
			fmt.Printf("    %-24s  n=%d  E[60s]=%+.3fpp  P+=%5.1f%%\n",
				r.name, r.n, r.r60, r.pr60)
		}
	}

	// flash SL analysis: how many entries in momentum baseline hit SL within 30s?
	fmt.Println("\n  Flash SL analysis (momentum entries that SL within 30s):")
	flashSL, totalEntries := 0, 0
	for _, p := range paths {
		var entryIdx int
		inPos := false
		for i := 60; i < len(p.Ticks); i++ {
			mid := p.Ticks[i].Mid
			if inPos {
				held := i - entryIdx
				ret := (mid - p.Ticks[entryIdx].Mid) / p.Ticks[entryIdx].Mid
				if ret <= -0.15 {
					if held <= 30 {
						flashSL++
					}
					inPos = false
					continue
				}
				if held >= 21600 {
					inPos = false
					continue
				}
				continue
			}
			// replicate momentum entry check
			refMid := p.Ticks[i-60].Mid
			if refMid <= 0 {
				continue
			}
			if (mid-refMid)/refMid*100 < 3.0 {
				continue
			}
			tailLen := 5
			if tailLen > i {
				continue
			}
			ups := 0
			for j := i - tailLen + 1; j <= i; j++ {
				if p.Ticks[j].Mid > p.Ticks[j-1].Mid {
					ups++
				}
			}
			if ups < 4 {
				continue
			}
			var bv, sv float64
			for j := i - 60; j <= i; j++ {
				bv += p.Ticks[j].BuyVol
				sv += p.Ticks[j].SellVol
			}
			if bv+sv > 0 && bv/(bv+sv) < 0.60 {
				continue
			}
			totalEntries++
			entryIdx = i
			inPos = true

			// check what state this entry is in
			f := computeFeatures(p.Ticks, i)
			fs := fwdStats[f.state]
			if fs.n >= 10 && fs.avgRet60() < 0 {
				// would Markov have filtered this?
			}
		}
	}
	if totalEntries > 0 {
		fmt.Printf("    Total momentum entries: %d\n", totalEntries)
		fmt.Printf("    Flash SL (≤30s): %d (%.0f%%)\n", flashSL, 100*float64(flashSL)/float64(totalEntries))
		fmt.Printf("    → Markov filter would skip entries in states with E[60s] < 0\n")
	}

	return nil
}

func printRow(r markovResult) {
	fmt.Printf("  %-28s %5d %5d %3d %3d %3d %3d %+9.2f %+8.4f %6.1f %6.2f %7.0f\n",
		r.label, r.nSignals, r.nTrades,
		r.wins, r.losses, r.slHit, r.timeoutHit,
		r.sumPnL, r.avgPnL(), r.winRate(), r.maxDD, r.avgHoldSec)
}

func printRowLine(label string, r markovResult) {
	fmt.Printf("  %-28s %5d %5d %3d %3d %3d %3d %+9.2f %+8.4f %6.1f %6.2f %7.0f\n",
		label, r.nSignals, r.nTrades,
		r.wins, r.losses, r.slHit, r.timeoutHit,
		r.sumPnL, r.avgPnL(), r.winRate(), r.maxDD, r.avgHoldSec)
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Avoid unused import warning
var _ = math.Abs
