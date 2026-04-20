package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
)

// tpslTrade 只保留回放必须的字段，return/direction 用原始 pnl/capital 还原，避开方向语义歧义。
type tpslTrade struct {
	ID        int64
	EntryPx   float64
	ClosePx   float64
	Capital   float64
	NaturalR  float64 // pnl / capital —— python 实际关仓时的单笔回报率
	FeeBP     float64 // 每笔总 round-trip 手续费（BP，供 sweep 扣减）
	Timestamp string
	CloseTime string
	// 可选：该笔交易持仓期间 PM 快照最高 / 最低价（仅 join 上 odds_snapshot 的子集有）
	HasPath bool
	PeakR   float64 // max( (px-entry)/entry ) across snaps between open/close
	TroughR float64 // min( (px-entry)/entry )
}

func loadTpslTrades(ctx context.Context, db *sql.DB, feeBP float64) ([]tpslTrade, error) {
	const q = `
SELECT t.id, t.entry_price, t.close_price, t.quantity, t.pnl_usdc,
       t.timestamp, t.close_time, t.token_id
FROM trades t
WHERE t.status='CLOSED' AND t.close_price IS NOT NULL AND t.pnl_usdc IS NOT NULL
  AND t.entry_price > 0 AND t.quantity > 0
ORDER BY t.close_time`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []tpslTrade
	for rows.Next() {
		var (
			id               int64
			entry, closePx   float64
			qty, pnl         float64
			ts, closeTime    sql.NullString
			tokenID          sql.NullString
			timestampStr, ct string
		)
		if err := rows.Scan(&id, &entry, &closePx, &qty, &pnl, &ts, &closeTime, &tokenID); err != nil {
			return nil, err
		}
		if ts.Valid {
			timestampStr = ts.String
		}
		if closeTime.Valid {
			ct = closeTime.String
		}
		capital := qty * entry
		if capital <= 0 {
			continue
		}
		out = append(out, tpslTrade{
			ID:        id,
			EntryPx:   entry,
			ClosePx:   closePx,
			Capital:   capital,
			NaturalR:  pnl / capital,
			FeeBP:     feeBP,
			Timestamp: timestampStr,
			CloseTime: ct,
		})
		_ = tokenID // reserved for future per-trade path enrichment via join (done below in batch)
	}
	return out, rows.Err()
}

// enrichWithSnapshots: 对能 join 到 odds_snapshot 的 trade 补 peak/trough。
// PM 快照的 price 指 YES token price；python `entry_price` 也是 token 入场价。
// 用 (px - entry)/entry 作为 raw return 代理，**方向符号交给后续 sweep 决定**
// （对 YES 卖方向做正向，对 NO 做反向在 natural return 已经封装好；
// 这里我们只关心价带峰谷，用来捕捉 TP 真正可能触达的最大涨幅/跌幅）。
func enrichWithSnapshots(ctx context.Context, db *sql.DB, trades []tpslTrade) ([]tpslTrade, error) {
	if len(trades) == 0 {
		return trades, nil
	}
	// 批量一次性拉 所有 可能 join 的 snapshot (范围限于 trades 的最早 open 到最晚 close)
	const q = `
SELECT o.market_id, o.polymarket_price, o.snapshot_timestamp
FROM odds_snapshot o
WHERE o.polymarket_price IS NOT NULL`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type snap struct {
		ts string
		px float64
	}
	byMarket := map[string][]snap{}
	for rows.Next() {
		var mid, ts string
		var px float64
		if err := rows.Scan(&mid, &px, &ts); err != nil {
			return nil, err
		}
		byMarket[mid] = append(byMarket[mid], snap{ts: ts, px: px})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 再查一次 trades 拿 token_id 做 lookup（上面 load 里 tokenID 没 return 出来，简单起见重查）
	const q2 = `SELECT id, token_id FROM trades WHERE status='CLOSED' AND token_id IS NOT NULL`
	tokenByID := map[int64]string{}
	r2, err := db.QueryContext(ctx, q2)
	if err != nil {
		return nil, err
	}
	defer r2.Close()
	for r2.Next() {
		var id int64
		var tok sql.NullString
		if err := r2.Scan(&id, &tok); err != nil {
			return nil, err
		}
		if tok.Valid {
			tokenByID[id] = tok.String
		}
	}

	for i := range trades {
		t := &trades[i]
		tok, ok := tokenByID[t.ID]
		if !ok {
			continue
		}
		snaps := byMarket[tok]
		if len(snaps) == 0 {
			continue
		}
		peak, trough := math.Inf(-1), math.Inf(1)
		seen := 0
		for _, s := range snaps {
			if s.ts < t.Timestamp || (t.CloseTime != "" && s.ts > t.CloseTime) {
				continue
			}
			if s.px > peak {
				peak = s.px
			}
			if s.px < trough {
				trough = s.px
			}
			seen++
		}
		if seen == 0 {
			continue
		}
		t.HasPath = true
		t.PeakR = (peak - t.EntryPx) / t.EntryPx
		t.TroughR = (trough - t.EntryPx) / t.EntryPx
	}
	return trades, nil
}

// sweepResult 单个 (TP, SL) 组合在样本池上的表现汇总。
type sweepResult struct {
	TP, SL      float64 // 正数（百分比），SL=0 代表不设止损
	N           int
	TPHit       int
	SLHit       int
	Natural     int
	SumRet      float64 // 样本回报率之和（每笔 -fee 后的几何累加近似：简单相加便于比较）
	SumRetNoFee float64
	MaxDD       float64
	Wins        int
	Losses      int
}

func (r sweepResult) avgRet() float64 {
	if r.N == 0 {
		return 0
	}
	return r.SumRet / float64(r.N)
}
func (r sweepResult) hit() float64 {
	if r.N == 0 {
		return 0
	}
	return 100 * float64(r.Wins) / float64(r.N)
}

// sweepEndpoint 端点近似：只看 close 时刻回报率是否跨过 TP/SL。
// 对 natural return 正向为赢（YES 策略下 pnl>0 → 价格朝我动），反向则亏；
// SL/TP 阈值都基于 natural return 的绝对值。
// feeBP 从 sum return 里扣 round-trip 两腿。
func sweepEndpoint(trades []tpslTrade, tpPct, slPct float64) sweepResult {
	res := sweepResult{TP: tpPct, SL: slPct}
	equity := 0.0
	peak := 0.0
	feeAdj := 2 * trades[0].FeeBP / 10000.0 // 每笔两腿手续费
	if trades[0].FeeBP == 0 {
		feeAdj = 0
	}
	_ = feeAdj // 单独算
	for _, t := range trades {
		r := t.NaturalR
		tp := tpPct / 100.0
		sl := slPct / 100.0
		switch {
		case tp > 0 && r >= tp:
			r = tp
			res.TPHit++
		case sl > 0 && r <= -sl:
			r = -sl
			res.SLHit++
		default:
			res.Natural++
		}
		// 扣 round-trip 手续费
		fee := 2 * t.FeeBP / 10000.0
		rNet := r - fee
		res.SumRet += rNet
		res.SumRetNoFee += r
		res.N++
		if rNet > 0.0005 {
			res.Wins++
		} else if rNet < -0.0005 {
			res.Losses++
		}
		equity += rNet
		if equity > peak {
			peak = equity
		}
		if peak-equity > res.MaxDD {
			res.MaxDD = peak - equity
		}
	}
	return res
}

// sweepPath 基于 peak/trough 的"真 TP 触达"近似：若 TP 已被 peak 超过则直接按 TP 结算。
// 仅对 HasPath=true 子集生效，样本小但更接近真实回放。
func sweepPath(trades []tpslTrade, tpPct, slPct float64) sweepResult {
	res := sweepResult{TP: tpPct, SL: slPct}
	for _, t := range trades {
		if !t.HasPath {
			continue
		}
		r := t.NaturalR
		tp := tpPct / 100.0
		sl := slPct / 100.0
		tpTouched := tp > 0 && t.PeakR >= tp
		slTouched := sl > 0 && t.TroughR <= -sl
		switch {
		case tpTouched && slTouched:
			// 都触过 —— 无法从稀疏快照判断先后；保守假设 SL 先触（劣势路径）
			r = -sl
			res.SLHit++
		case tpTouched:
			r = tp
			res.TPHit++
		case slTouched:
			r = -sl
			res.SLHit++
		default:
			res.Natural++
		}
		fee := 2 * t.FeeBP / 10000.0
		rNet := r - fee
		res.SumRet += rNet
		res.SumRetNoFee += r
		res.N++
		if rNet > 0.0005 {
			res.Wins++
		} else if rNet < -0.0005 {
			res.Losses++
		}
	}
	return res
}

func runTpslSweep(ctx context.Context, db *sql.DB, feeBP float64) error {
	trades, err := loadTpslTrades(ctx, db, feeBP)
	if err != nil {
		return fmt.Errorf("load trades: %w", err)
	}
	trades, err = enrichWithSnapshots(ctx, db, trades)
	if err != nil {
		return fmt.Errorf("enrich: %w", err)
	}
	pathN := 0
	for _, t := range trades {
		if t.HasPath {
			pathN++
		}
	}
	fmt.Println("════════════════════════════════════════")
	fmt.Println(" TP/SL sweep · python closed trades")
	fmt.Printf(" pool: %d total (%d with intra-hold snaps)\n", len(trades), pathN)
	fmt.Printf(" fee_bp: %.1f each leg (round-trip %.1f bp)\n", feeBP, feeBP*2)
	fmt.Println("════════════════════════════════════════")

	// 用自然回报率做 baseline（没有 TP/SL 覆盖）
	base := sweepEndpoint(trades, 0, 0)
	fmt.Printf("\n-- baseline (no TP/SL, fee=%.1fbp) --\n", feeBP)
	fmt.Printf("  avg_ret %+.4f · hit %.1f%% · sum_ret %+.2f · mdd %.2f · n=%d\n",
		base.avgRet(), base.hit(), base.SumRet, base.MaxDD, base.N)

	tps := []float64{5, 10, 15, 20, 30, 40, 50, 75, 100}
	sls := []float64{0, 5, 10, 15, 20, 30}

	var endpointResults []sweepResult
	for _, tp := range tps {
		for _, sl := range sls {
			endpointResults = append(endpointResults, sweepEndpoint(trades, tp, sl))
		}
	}
	sort.Slice(endpointResults, func(i, j int) bool {
		return endpointResults[i].SumRet > endpointResults[j].SumRet
	})

	fmt.Println("\n-- endpoint approximation (top 10 by sum_ret) --")
	fmt.Printf("  %-5s %-5s %6s %7s %7s %7s %9s %8s %8s %8s\n",
		"TP%", "SL%", "n", "tp_hit", "sl_hit", "natur", "sum_ret", "avg_ret", "hit%", "mdd")
	for i, r := range endpointResults {
		if i >= 10 {
			break
		}
		fmt.Printf("  %-5.0f %-5.0f %6d %7d %7d %7d %+9.3f %+8.4f %8.1f %8.3f\n",
			r.TP, r.SL, r.N, r.TPHit, r.SLHit, r.Natural,
			r.SumRet, r.avgRet(), r.hit(), r.MaxDD)
	}

	fmt.Println("\n-- endpoint approximation (bottom 5) --")
	for i := len(endpointResults) - 5; i < len(endpointResults); i++ {
		if i < 0 {
			continue
		}
		r := endpointResults[i]
		fmt.Printf("  %-5.0f %-5.0f %6d %7d %7d %7d %+9.3f %+8.4f %8.1f %8.3f\n",
			r.TP, r.SL, r.N, r.TPHit, r.SLHit, r.Natural,
			r.SumRet, r.avgRet(), r.hit(), r.MaxDD)
	}

	if pathN > 0 {
		var pathResults []sweepResult
		for _, tp := range tps {
			for _, sl := range sls {
				pathResults = append(pathResults, sweepPath(trades, tp, sl))
			}
		}
		sort.Slice(pathResults, func(i, j int) bool {
			return pathResults[i].SumRet > pathResults[j].SumRet
		})
		fmt.Printf("\n-- path-aware (n=%d, 保守：同时触及 TP/SL 判 SL) --\n", pathN)
		fmt.Printf("  %-5s %-5s %6s %7s %7s %7s %9s %8s %8s\n",
			"TP%", "SL%", "n", "tp_hit", "sl_hit", "natur", "sum_ret", "avg_ret", "hit%")
		for i, r := range pathResults {
			if i >= 10 {
				break
			}
			fmt.Printf("  %-5.0f %-5.0f %6d %7d %7d %7d %+9.3f %+8.4f %8.1f\n",
				r.TP, r.SL, r.N, r.TPHit, r.SLHit, r.Natural,
				r.SumRet, r.avgRet(), r.hit())
		}
	}

	fmt.Println("\n── 数据限制提示 ──")
	fmt.Println("  1) python trades 没有持仓期 tick 数据，端点回放只能看 close 时刻是否跨阈值——")
	fmt.Println("     这会【低估】TP 命中率（价格中途触达 TP 后又跌回的没被计入）。")
	fmt.Println("  2) 仅 ~9 笔能 join odds_snapshot 拿到 peak/trough 代理，样本太小结论不稳。")
	fmt.Println("  3) 正确做法：本项目 daemon 已采 1Hz tick，下步把 tick 写盘累积，")
	fmt.Println("     用自家 momentum entry 的真实 path 做回放（Phase 7.e 计划中）。")
	fmt.Println("  4) python trades 里有 184 笔是 -NO entry 逻辑，方向归一到 natural_return")
	fmt.Println("     (pnl/capital) 后再 sweep；TP/SL 阈值本身是方向无关的 % 回报率。")

	return nil
}
