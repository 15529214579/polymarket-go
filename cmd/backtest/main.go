// cmd/backtest 离线分析 python polymarket-agent 的交易/快照数据。
// 只读打开，不碰 python 进程。核心问题："PM vs bookmaker >5pp 这条路历史上有 edge 吗？"
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite" //nolint:revive // register sqlite driver
)

type tradeRow struct {
	ID         int64
	MarketID   string
	Status     string
	EntryPx    float64
	Qty        float64
	ClosePx    sql.NullFloat64
	PnL        sql.NullFloat64
	Timestamp  string
	ScanType   string
	GapAtEntry sql.NullFloat64
}

type bucketStats struct {
	N             int
	NClosed       int
	Wins          int
	Losses        int
	BreakEven     int
	OpenOrPending int
	SumPnL        float64
	SumCapital    float64
	PnLs          []float64
}

func (b *bucketStats) add(t tradeRow) {
	b.N++
	capital := t.Qty * t.EntryPx
	switch t.Status {
	case "CLOSED":
		b.NClosed++
		b.SumCapital += capital
		if t.PnL.Valid {
			b.SumPnL += t.PnL.Float64
			b.PnLs = append(b.PnLs, t.PnL.Float64)
			switch {
			case t.PnL.Float64 > 0.01:
				b.Wins++
			case t.PnL.Float64 < -0.01:
				b.Losses++
			default:
				b.BreakEven++
			}
		}
	default:
		b.OpenOrPending++
	}
}

func (b bucketStats) roiPct() float64 {
	if b.SumCapital == 0 {
		return 0
	}
	return 100 * b.SumPnL / b.SumCapital
}

func (b bucketStats) hitRatePct() float64 {
	if b.NClosed == 0 {
		return 0
	}
	return 100 * float64(b.Wins) / float64(b.NClosed)
}

func (b bucketStats) maxDrawdownUSD() float64 {
	if len(b.PnLs) == 0 {
		return 0
	}
	peak, equity, mdd := 0.0, 0.0, 0.0
	for _, p := range b.PnLs {
		equity += p
		if equity > peak {
			peak = equity
		}
		if peak-equity > mdd {
			mdd = peak - equity
		}
	}
	return mdd
}

func main() {
	defaultDB := filepath.Join(os.Getenv("HOME"), ".openclaw", "workspace-dev3", "polymarket-agent", "db", "polymarket_agent.db")
	dbPath := flag.String("db", defaultDB, "python polymarket-agent sqlite db path")
	minGap := flag.Float64("min_gap", 5.0, "minimum abs(gap_pp) to include in analysis")
	mode := flag.String("mode", "summary", "summary | tpsl-sweep | tickpath-sweep | markov | btc-markov")
	feeBP := flag.Float64("fee_bp", 0, "per-leg fee (bp); round-trip = 2x")
	defaultTickDir := filepath.Join(os.Getenv("HOME"), "work", "polymarket-go", "db", "tickpath")
	tickDir := flag.String("tick_dir", defaultTickDir, "directory of per-position .jsonl tick recordings")
	defaultJournalDir := filepath.Join(os.Getenv("HOME"), "work", "polymarket-go", "db", "journal")
	journalDir := flag.String("journal_dir", defaultJournalDir, "directory of journal .jsonl files")
	btcDays := flag.Int("days", 90, "btc-markov: number of historical days to fetch")
	btcTrainPct := flag.Float64("train_pct", 0.67, "btc-markov: fraction of data used for training (0..1)")
	defaultDBDir := filepath.Join(os.Getenv("HOME"), "work", "polymarket-go", "db")
	btcDBDir := flag.String("db_dir", defaultDBDir, "btc-markov: directory for btc_markov.db")
	flag.Parse()

	if err := dispatch(*mode, *dbPath, *minGap, *feeBP, *tickDir, *journalDir, *btcDays, *btcTrainPct, *btcDBDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dispatch(mode, dbPath string, minGap, feeBP float64, tickDir, journalDir string, btcDays int, btcTrainPct float64, btcDBDir string) error {
	switch mode {
	case "summary":
		return run(dbPath, minGap)
	case "tpsl-sweep":
		uri := fmt.Sprintf("file:%s?mode=ro&immutable=1", dbPath)
		db, err := sql.Open("sqlite", uri)
		if err != nil {
			return fmt.Errorf("open db: %w", err)
		}
		defer db.Close()
		return runTpslSweep(context.Background(), db, feeBP)
	case "tickpath-sweep":
		return runTickPathSweep(tickDir, journalDir, feeBP)
	case "markov":
		return runMarkovBacktest(tickDir, feeBP)
	case "btc-markov":
		return runBTCMarkovBacktest(btcDays, btcTrainPct, btcDBDir)
	default:
		return fmt.Errorf("unknown mode: %s", mode)
	}
}

func run(dbPath string, minGap float64) error {
	// 只读 URI：file:...?mode=ro&immutable=1
	uri := fmt.Sprintf("file:%s?mode=ro&immutable=1", dbPath)
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	trades, err := loadTrades(ctx, db)
	if err != nil {
		return fmt.Errorf("load trades: %w", err)
	}
	snaps, err := loadScanDist(ctx, db)
	if err != nil {
		return fmt.Errorf("load scan_dist: %w", err)
	}
	pnlCurve, err := loadPnLCurve(ctx, db)
	if err != nil {
		return fmt.Errorf("load pnl curve: %w", err)
	}

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" polymarket-go OFFLINE backtest")
	fmt.Println(" source:", dbPath)
	fmt.Println(" trades:", len(trades))
	fmt.Println("════════════════════════════════════════")

	printScanDist(snaps)
	printByScanType(trades)
	printByGapBucket(trades, minGap)
	printPnLCurve(pnlCurve)

	fmt.Println()
	fmt.Println("── 结论提示 ──")
	fmt.Println("  theodds_h2h 的已平仓胜率 / ROI 是 gap-bucket 策略最直接的实盘答卷。")
	fmt.Println("  gap 越大（>15pp）多是 league mismatch / outright 错配，不是真 arb。")
	fmt.Println("  窄 gap 段（5-10pp）样本稀少，需要更多数据才能判断是否可交易。")
	return nil
}

func loadTrades(ctx context.Context, db *sql.DB) ([]tradeRow, error) {
	const q = `
SELECT t.id, t.market_id, t.status, t.entry_price, t.quantity, t.close_price, t.pnl_usdc,
       t.timestamp,
       COALESCE((SELECT ol.scan_type FROM opportunity_log ol
                  WHERE ol.market_id = t.market_id
                    AND ol.scan_timestamp <= t.timestamp
                  ORDER BY ol.scan_timestamp DESC LIMIT 1), 'none') AS scan_type,
       (SELECT ol.gap_pp FROM opportunity_log ol
         WHERE ol.market_id = t.market_id
           AND ol.scan_timestamp <= t.timestamp
           AND ol.scan_type='theodds_h2h'
         ORDER BY ol.scan_timestamp DESC LIMIT 1) AS gap_at_entry
FROM trades t
ORDER BY t.timestamp`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []tradeRow
	for rows.Next() {
		var t tradeRow
		if err := rows.Scan(&t.ID, &t.MarketID, &t.Status, &t.EntryPx, &t.Qty,
			&t.ClosePx, &t.PnL, &t.Timestamp, &t.ScanType, &t.GapAtEntry); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return out, nil
}

type scanDistRow struct {
	ScanType string
	N        int
	AvgGap   sql.NullFloat64
	AvgPM    sql.NullFloat64
	AvgBook  sql.NullFloat64
}

func loadScanDist(ctx context.Context, db *sql.DB) ([]scanDistRow, error) {
	const q = `
SELECT scan_type, COUNT(*), AVG(ABS(gap_pp)), AVG(polymarket_price), AVG(bookmaker_prob)
FROM opportunity_log
GROUP BY scan_type
ORDER BY COUNT(*) DESC`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []scanDistRow
	for rows.Next() {
		var r scanDistRow
		if err := rows.Scan(&r.ScanType, &r.N, &r.AvgGap, &r.AvgPM, &r.AvgBook); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

type pnlPoint struct {
	Timestamp string
	Equity    float64
}

func loadPnLCurve(ctx context.Context, db *sql.DB) ([]pnlPoint, error) {
	const q = `
SELECT close_time, pnl_usdc FROM trades
WHERE status='CLOSED' AND close_time IS NOT NULL AND pnl_usdc IS NOT NULL
ORDER BY close_time`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []pnlPoint
	equity := 0.0
	for rows.Next() {
		var ts string
		var pnl float64
		if err := rows.Scan(&ts, &pnl); err != nil {
			return nil, err
		}
		equity += pnl
		out = append(out, pnlPoint{Timestamp: ts, Equity: equity})
	}
	return out, nil
}

func printScanDist(rows []scanDistRow) {
	fmt.Println()
	fmt.Println("── opportunity_log × scan_type ──")
	fmt.Printf("  %-15s %6s %10s %10s %10s\n", "scan_type", "n", "avg_gap", "avg_pm", "avg_book")
	for _, r := range rows {
		gap, pm, book := "-", "-", "-"
		if r.AvgGap.Valid {
			gap = fmt.Sprintf("%.2f", r.AvgGap.Float64)
		}
		if r.AvgPM.Valid {
			pm = fmt.Sprintf("%.3f", r.AvgPM.Float64)
		}
		if r.AvgBook.Valid {
			book = fmt.Sprintf("%.3f", r.AvgBook.Float64)
		}
		fmt.Printf("  %-15s %6d %10s %10s %10s\n", r.ScanType, r.N, gap, pm, book)
	}
}

func printByScanType(trades []tradeRow) {
	groups := map[string]*bucketStats{}
	for _, t := range trades {
		b, ok := groups[t.ScanType]
		if !ok {
			b = &bucketStats{}
			groups[t.ScanType] = b
		}
		b.add(t)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return groups[keys[i]].SumPnL > groups[keys[j]].SumPnL })

	fmt.Println()
	fmt.Println("── trades × scan_type ──")
	fmt.Printf("  %-15s %5s %7s %5s %5s %10s %10s %10s %8s %8s\n",
		"strategy", "n", "closed", "wins", "loss", "pnl_usd", "capital", "roi%", "hit%", "mdd_usd")
	for _, k := range keys {
		b := groups[k]
		fmt.Printf("  %-15s %5d %7d %5d %5d %10.2f %10.2f %+10.2f %8.1f %8.2f\n",
			k, b.N, b.NClosed, b.Wins, b.Losses,
			b.SumPnL, b.SumCapital, b.roiPct(), b.hitRatePct(), b.maxDrawdownUSD())
	}
}

func printByGapBucket(trades []tradeRow, minGap float64) {
	fmt.Println()
	fmt.Println("── theodds_h2h × gap_bucket （实盘已平仓样本）──")
	buckets := []struct {
		name    string
		lo, hi  float64
		stats   bucketStats
		samples int
	}{
		{name: "5-7pp", lo: 5, hi: 7},
		{name: "7-10pp", lo: 7, hi: 10},
		{name: "10-15pp", lo: 10, hi: 15},
		{name: "15-25pp", lo: 15, hi: 25},
		{name: "25-50pp", lo: 25, hi: 50},
		{name: "50+pp", lo: 50, hi: math.MaxFloat64},
	}
	for _, t := range trades {
		if t.ScanType != "theodds_h2h" || !t.GapAtEntry.Valid {
			continue
		}
		absGap := math.Abs(t.GapAtEntry.Float64)
		if absGap < minGap {
			continue
		}
		for i := range buckets {
			if absGap >= buckets[i].lo && absGap < buckets[i].hi {
				buckets[i].stats.add(t)
				buckets[i].samples++
				break
			}
		}
	}
	fmt.Printf("  %-10s %5s %7s %5s %5s %10s %10s %10s %8s\n",
		"gap", "n", "closed", "wins", "loss", "pnl_usd", "capital", "roi%", "hit%")
	for _, b := range buckets {
		s := b.stats
		if s.N == 0 {
			fmt.Printf("  %-10s %5s %7s %5s %5s %10s %10s %10s %8s\n",
				b.name, "0", "-", "-", "-", "-", "-", "-", "-")
			continue
		}
		fmt.Printf("  %-10s %5d %7d %5d %5d %10.2f %10.2f %+10.2f %8.1f\n",
			b.name, s.N, s.NClosed, s.Wins, s.Losses,
			s.SumPnL, s.SumCapital, s.roiPct(), s.hitRatePct())
	}
}

func printPnLCurve(pts []pnlPoint) {
	if len(pts) == 0 {
		return
	}
	fmt.Println()
	fmt.Println("── 累计 PnL 曲线（全策略合并，按 close_time 排序）──")
	fmt.Printf("  第一笔 close:  %s  equity=%+.2f\n", pts[0].Timestamp, pts[0].Equity)
	n := len(pts)
	mid := n / 2
	fmt.Printf("  中段 (#%d):     %s  equity=%+.2f\n", mid, pts[mid].Timestamp, pts[mid].Equity)
	fmt.Printf("  最新 (#%d):     %s  equity=%+.2f\n", n-1, pts[n-1].Timestamp, pts[n-1].Equity)
	// 峰值 & 最大回撤
	peak, mdd := pts[0].Equity, 0.0
	peakAt, troughAt := pts[0].Timestamp, pts[0].Timestamp
	for _, p := range pts {
		if p.Equity > peak {
			peak = p.Equity
			peakAt = p.Timestamp
		}
		if peak-p.Equity > mdd {
			mdd = peak - p.Equity
			troughAt = p.Timestamp
		}
	}
	fmt.Printf("  峰值 equity=%+.2f @ %s\n", peak, peakAt)
	fmt.Printf("  最大回撤  -%.2f USDC  谷底 @ %s\n", mdd, troughAt)
}
