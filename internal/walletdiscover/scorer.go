package walletdiscover

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func ScoreWallet(addr string, cand *Candidate, trades []Trade, closed []ClosedPosition, cfg Config) WalletScore {
	stats := buildStats(trades, closed, cfg)
	bot := botScore(stats)
	edge := edgeScore(stats)
	smart := smartMoneyScore(stats, cand, bot, edge)
	tier, reason := tierFor(stats, bot, edge, smart)
	action := followAction(tier, stats, bot, smart)
	flags := riskFlags(stats, bot)
	strengths := strengthFlags(stats, smart)
	score := smart - bot*0.35
	src := map[string]int{}
	if cand != nil {
		for k, v := range cand.Sources {
			src[k] = v
		}
	}
	return WalletScore{
		Address:         addr,
		Tier:            tier,
		Reason:          reason,
		FollowAction:    action,
		BotScore:        round2(bot),
		Edge:            round2(edge),
		SmartMoneyScore: round2(smart),
		Score:           round2(score),
		RiskFlags:       flags,
		Strengths:       strengths,
		Sources:         src,
		Stats:           stats,
	}
}

func buildStats(trades []Trade, closed []ClosedPosition, cfg Config) WalletStats {
	var st WalletStats
	st.RawTrades = len(trades)
	amountBuckets := map[string]int{}
	priceBuckets := map[string]int{}
	marketSides := map[string]map[string]struct{}{}
	marketCounts := map[string]int{}
	markets := map[string]struct{}{}
	categoryCounts := map[string]int{}
	minuteCounts := map[int64]int{}
	var totalNotional float64
	var validForRatios int

	for _, tr := range trades {
		if tr.Type != "" && !strings.EqualFold(tr.Type, "TRADE") {
			continue
		}
		side := strings.ToUpper(tr.Side)
		if side != "BUY" && side != "SELL" {
			continue
		}
		st.ValidTrades++
		if side == "BUY" {
			st.BuyTrades++
		} else {
			st.SellTrades++
		}
		notional := tr.NotionalUSD()
		totalNotional += notional
		if notional >= cfg.MinNotionalUSD {
			st.LargeTrades++
		}
		if tr.Price <= 0.05 || tr.Price >= 0.95 {
			st.ExtremePriceRatio++
		}
		amountBuckets[fmt.Sprintf("%.0f", notional)]++
		priceBuckets[fmt.Sprintf("%.3f", tr.Price)]++
		if tr.ConditionID != "" {
			markets[tr.ConditionID] = struct{}{}
			marketCounts[tr.ConditionID]++
			if _, ok := marketSides[tr.ConditionID]; !ok {
				marketSides[tr.ConditionID] = map[string]struct{}{}
			}
			marketSides[tr.ConditionID][side] = struct{}{}
		}
		targetCat := TradeTargetCategory(tr)
		if TargetCategoryAllowed(targetCat, cfg.TargetCategories) && targetCat != "other" {
			st.TargetTrades++
			if notional >= cfg.MinNotionalUSD {
				st.TargetLargeTrades++
			}
		}
		if cat := coarseCategory(tr.Title + " " + tr.Slug); cat != "" {
			categoryCounts[cat]++
		}
		if isPositiveFootballScoreTrade(tr) {
			st.FootballScoreTrades++
			if notional >= cfg.MinNotionalUSD {
				st.FootballScoreLargeTrades++
			}
		}
		if tr.Timestamp > 0 {
			minute := tr.Timestamp / 60
			minuteCounts[minute]++
			if minuteCounts[minute] > st.MaxTradesPerMinute {
				st.MaxTradesPerMinute = minuteCounts[minute]
			}
		}
		validForRatios++
	}
	if validForRatios > 0 {
		st.ExtremePriceRatio = st.ExtremePriceRatio / float64(validForRatios)
		st.FixedAmountRatio = maxBucketRatio(amountBuckets, validForRatios)
		st.FixedPriceRatio = maxBucketRatio(priceBuckets, validForRatios)
		st.AvgTradeNotional = totalNotional / float64(validForRatios)
		st.BuyRatio = float64(st.BuyTrades) / float64(validForRatios)
		st.TopMarketRatio = maxBucketRatio(marketCounts, validForRatios)
		st.TargetTradeRatio = float64(st.TargetTrades) / float64(validForRatios)
	}
	st.DistinctMarkets = len(markets)
	st.DistinctCategories = len(categoryCounts)
	st.TopCategory, st.TopCategoryRatio = topCategory(categoryCounts, validForRatios)
	for _, sides := range marketSides {
		if len(sides) > 1 {
			st.OppositeSideMarkets++
		}
	}
	for _, n := range minuteCounts {
		if n >= 4 {
			st.BurstTrades += n
		}
	}

	st.ClosedPositions = len(closed)
	for _, p := range closed {
		st.ClosedPnL += p.RealizedPnL
		st.ClosedCapital += p.TotalBought
		if p.RealizedPnL > 0 {
			st.PositiveClosed++
		}
		if isPositiveFootballScoreClosed(p) {
			st.FootballScoreClosed++
			st.FootballScoreClosedPnL += p.RealizedPnL
			st.FootballScoreClosedCapital += p.TotalBought
			if p.RealizedPnL > 0 {
				st.FootballScoreClosedWins++
			}
		}
	}
	if st.ClosedCapital > 0 {
		st.ClosedROI = st.ClosedPnL / st.ClosedCapital * 100
	}
	if st.ClosedPositions > 0 {
		st.ClosedWinRate = float64(st.PositiveClosed) / float64(st.ClosedPositions) * 100
	}
	if st.FootballScoreClosedCapital > 0 {
		st.FootballScoreClosedROI = st.FootballScoreClosedPnL / st.FootballScoreClosedCapital * 100
	}
	applyCopySimulation(&st, trades, cfg)
	applyTargetCopySimulation(&st, trades, cfg)
	applyFootballScoreCopySimulation(&st, trades, cfg)

	st.ExtremePriceRatio = round4(st.ExtremePriceRatio)
	st.FixedAmountRatio = round4(st.FixedAmountRatio)
	st.FixedPriceRatio = round4(st.FixedPriceRatio)
	st.BuyRatio = round4(st.BuyRatio)
	st.TopCategoryRatio = round4(st.TopCategoryRatio)
	st.TargetTradeRatio = round4(st.TargetTradeRatio)
	st.TopMarketRatio = round4(st.TopMarketRatio)
	st.ClosedPnL = round2(st.ClosedPnL)
	st.ClosedCapital = round2(st.ClosedCapital)
	st.ClosedROI = round2(st.ClosedROI)
	st.ClosedWinRate = round2(st.ClosedWinRate)
	st.AvgTradeNotional = round2(st.AvgTradeNotional)
	st.CopyPnL = round2(st.CopyPnL)
	st.CopyCapital = round2(st.CopyCapital)
	st.CopyROI = round2(st.CopyROI)
	st.CopyWinRate = round2(st.CopyWinRate)
	st.CopyOpenCost = round2(st.CopyOpenCost)
	st.TargetCopyPnL = round2(st.TargetCopyPnL)
	st.TargetCopyCapital = round2(st.TargetCopyCapital)
	st.TargetCopyROI = round2(st.TargetCopyROI)
	st.TargetCopyWinRate = round2(st.TargetCopyWinRate)
	st.TargetCopyOpenCost = round2(st.TargetCopyOpenCost)
	st.FootballScoreClosedPnL = round2(st.FootballScoreClosedPnL)
	st.FootballScoreClosedCapital = round2(st.FootballScoreClosedCapital)
	st.FootballScoreClosedROI = round2(st.FootballScoreClosedROI)
	st.FootballScoreCopyPnL = round2(st.FootballScoreCopyPnL)
	st.FootballScoreCopyCapital = round2(st.FootballScoreCopyCapital)
	st.FootballScoreCopyROI = round2(st.FootballScoreCopyROI)
	st.FootballScoreCopyWinRate = round2(st.FootballScoreCopyWinRate)
	st.FootballScoreCopyOpenCost = round2(st.FootballScoreCopyOpenCost)
	return st
}

type copyPosition struct {
	qty      float64
	cost     float64
	avgPrice float64
}

type copySimResult struct {
	Buys          int
	ClosedTrades  int
	Wins          int
	PnL           float64
	Capital       float64
	ROI           float64
	WinRate       float64
	OpenPositions int
	OpenCost      float64
}

func applyCopySimulation(st *WalletStats, trades []Trade, cfg Config) {
	res := simulateCopy(trades, cfg, nil)
	st.CopyBuys = res.Buys
	st.CopyClosedTrades = res.ClosedTrades
	st.CopyWins = res.Wins
	st.CopyPnL = res.PnL
	st.CopyCapital = res.Capital
	st.CopyROI = res.ROI
	st.CopyWinRate = res.WinRate
	st.CopyOpenPositions = res.OpenPositions
	st.CopyOpenCost = res.OpenCost
}

func applyTargetCopySimulation(st *WalletStats, trades []Trade, cfg Config) {
	res := simulateCopy(trades, cfg, func(tr Trade) bool {
		cat := TradeTargetCategory(tr)
		return cat != "other" && TargetCategoryAllowed(cat, cfg.TargetCategories)
	})
	st.TargetCopyBuys = res.Buys
	st.TargetCopyClosed = res.ClosedTrades
	st.TargetCopyWins = res.Wins
	st.TargetCopyPnL = res.PnL
	st.TargetCopyCapital = res.Capital
	st.TargetCopyROI = res.ROI
	st.TargetCopyWinRate = res.WinRate
	st.TargetCopyOpen = res.OpenPositions
	st.TargetCopyOpenCost = res.OpenCost
}

func applyFootballScoreCopySimulation(st *WalletStats, trades []Trade, cfg Config) {
	res := simulateCopy(trades, cfg, isPositiveFootballScoreTrade)
	st.FootballScoreCopyBuys = res.Buys
	st.FootballScoreCopyClosed = res.ClosedTrades
	st.FootballScoreCopyWins = res.Wins
	st.FootballScoreCopyPnL = res.PnL
	st.FootballScoreCopyCapital = res.Capital
	st.FootballScoreCopyROI = res.ROI
	st.FootballScoreCopyWinRate = res.WinRate
	st.FootballScoreCopyOpen = res.OpenPositions
	st.FootballScoreCopyOpenCost = res.OpenCost
}

func isPositiveFootballScoreTrade(tr Trade) bool {
	return IsFootballScoreTrade(tr) && strings.EqualFold(strings.TrimSpace(tr.Outcome), "yes")
}

func isPositiveFootballScoreClosed(p ClosedPosition) bool {
	text := p.Title + " " + p.Slug
	return feed.IsFootballScoreMarketText(text) && strings.EqualFold(strings.TrimSpace(p.Outcome), "yes")
}

func simulateCopy(trades []Trade, cfg Config, include func(Trade) bool) copySimResult {
	var res copySimResult
	if cfg.CopyStakeUSD <= 0 {
		return res
	}
	positions := map[string]*copyPosition{}
	slip := cfg.CopySlippageBP / 10000
	fee := cfg.CopyFeeBP / 10000
	sorted := append([]Trade(nil), trades...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Timestamp < sorted[j].Timestamp })

	closePosition := func(asset string, exitPrice float64) {
		pos := positions[asset]
		if pos == nil || pos.qty <= 0 || exitPrice <= 0 {
			return
		}
		netExit := exitPrice * (1 - slip - fee)
		if netExit < 0 {
			netExit = 0
		}
		proceeds := pos.qty * netExit
		pnl := proceeds - pos.cost
		res.PnL += pnl
		res.Capital += pos.cost
		res.ClosedTrades++
		if pnl > 0 {
			res.Wins++
		}
		delete(positions, asset)
	}

	for _, tr := range sorted {
		if include != nil && !include(tr) {
			continue
		}
		asset := tr.Asset
		if asset == "" {
			continue
		}
		typ := strings.ToUpper(tr.Type)
		side := strings.ToUpper(tr.Side)
		if typ == "REDEEM" {
			closePosition(asset, 1)
			continue
		}
		if typ != "" && typ != "TRADE" {
			continue
		}
		if side == "BUY" {
			notional := tr.NotionalUSD()
			if notional < cfg.MinNotionalUSD || tr.Price < cfg.MinPrice || tr.Price > cfg.MaxPrice {
				continue
			}
			entryPrice := tr.Price * (1 + slip + fee)
			if entryPrice <= 0 || entryPrice >= 1 {
				continue
			}
			qty := cfg.CopyStakeUSD / entryPrice
			pos := positions[asset]
			if pos == nil {
				pos = &copyPosition{}
				positions[asset] = pos
			}
			pos.qty += qty
			pos.cost += cfg.CopyStakeUSD
			pos.avgPrice = pos.cost / pos.qty
			res.Buys++
			continue
		}
		if side == "SELL" {
			closePosition(asset, tr.Price)
		}
	}

	res.OpenPositions = len(positions)
	for _, pos := range positions {
		res.OpenCost += pos.cost
	}
	if res.Capital > 0 {
		res.ROI = res.PnL / res.Capital * 100
	}
	if res.ClosedTrades > 0 {
		res.WinRate = float64(res.Wins) / float64(res.ClosedTrades) * 100
	}
	return res
}

func botScore(st WalletStats) float64 {
	var s float64
	if st.ValidTrades == 0 {
		return 50
	}
	s += clamp(st.ExtremePriceRatio*45, 0, 35)
	s += clamp(st.FixedAmountRatio*25, 0, 20)
	s += clamp(st.FixedPriceRatio*20, 0, 15)
	if st.MaxTradesPerMinute >= 10 {
		s += 25
	} else if st.MaxTradesPerMinute >= 6 {
		s += 15
	} else if st.MaxTradesPerMinute >= 4 {
		s += 8
	}
	if st.ValidTrades > 0 {
		s += clamp(float64(st.BurstTrades)/float64(st.ValidTrades)*20, 0, 15)
	}
	if st.DistinctCategories >= 6 && st.ValidTrades >= 50 {
		s += 15
	} else if st.DistinctCategories >= 4 && st.ValidTrades >= 50 {
		s += 8
	}
	if st.OppositeSideMarkets >= 5 {
		s += 15
	} else if st.OppositeSideMarkets >= 2 {
		s += 7
	}
	return clamp(s, 0, 100)
}

func edgeScore(st WalletStats) float64 {
	var s float64
	if st.LargeTrades >= 20 {
		s += 20
	} else {
		s += float64(st.LargeTrades)
	}
	if st.ClosedPositions >= 10 {
		s += 15
	} else {
		s += float64(st.ClosedPositions) * 1.2
	}
	if st.ClosedROI > 0 {
		s += clamp(st.ClosedROI, 0, 35)
	}
	if st.ClosedWinRate > 50 {
		s += clamp((st.ClosedWinRate-50)*0.8, 0, 20)
	}
	if st.ClosedPnL > 0 {
		s += clamp(math.Log10(st.ClosedPnL+1)*8, 0, 20)
	}
	if st.CopyClosedTrades >= 5 {
		if st.CopyROI > 0 {
			s += clamp(st.CopyROI, 0, 30)
		} else {
			s -= clamp(math.Abs(st.CopyROI)*0.6, 0, 25)
		}
		if st.CopyWinRate > 50 {
			s += clamp((st.CopyWinRate-50)*0.6, 0, 15)
		}
		if st.CopyPnL > 0 {
			s += clamp(math.Log10(st.CopyPnL+1)*6, 0, 12)
		}
	}
	if st.DistinctCategories > 0 && st.DistinctCategories <= 3 && st.ValidTrades >= 20 {
		s += 8
	}
	return clamp(s, 0, 100)
}

func smartMoneyScore(st WalletStats, cand *Candidate, bot, edge float64) float64 {
	var s float64
	s += edge * 0.75
	if st.ClosedPositions >= 25 {
		s += 14
	} else {
		s += float64(st.ClosedPositions) * 0.55
	}
	if st.ClosedPnL > 100 {
		s += clamp(math.Log10(st.ClosedPnL)*9, 0, 20)
	}
	if st.ClosedROI >= 15 {
		s += 15
	} else if st.ClosedROI >= 5 {
		s += 8
	}
	if st.ClosedWinRate >= 60 && st.ClosedPositions >= 10 {
		s += 10
	} else if st.ClosedWinRate >= 53 && st.ClosedPositions >= 10 {
		s += 5
	}
	if st.CopyClosedTrades >= 5 {
		if st.CopyROI >= 10 {
			s += 14
		} else if st.CopyROI > 0 {
			s += 6
		} else {
			s -= clamp(math.Abs(st.CopyROI), 0, 25)
		}
	}
	if st.TopCategoryRatio >= 0.45 && st.TopCategoryRatio <= 0.90 && st.DistinctMarkets >= 8 {
		s += 8
	}
	if st.BuyRatio >= 0.55 && st.BuyRatio <= 0.92 {
		s += 4
	}
	if cand != nil {
		if cand.Sources["holder"] > 0 && cand.Sources["recent_trade"] > 0 {
			s += 5
		}
		if cand.Sources["existing"] > 0 {
			s += 3
		}
	}
	s -= bot * 0.55
	if st.ValidTrades < 12 && st.ClosedPositions < 6 {
		s -= 20
	}
	if st.ClosedCapital < 100 && st.ClosedPositions > 0 {
		s -= 8
	}
	return clamp(s, 0, 100)
}

func tierFor(st WalletStats, bot, edge, smart float64) (string, string) {
	if st.ValidTrades < 8 && st.ClosedPositions < 3 {
		return "D", "insufficient sample"
	}
	if bot >= 60 {
		return "BOT", "bot-like flow"
	}
	if bot >= 45 && smart < 70 {
		return "BOT", "high bot score without enough edge"
	}
	if st.CopyClosedTrades >= 5 && st.CopyROI <= 0 {
		return "D", "negative copy-sim ROI"
	}
	if smart >= 70 && edge >= 55 && bot < 25 && st.LargeTrades >= 10 && st.CopyClosedTrades >= 5 && st.CopyROI >= 5 {
		return "A", "strong smart-money signal and low bot score"
	}
	if smart >= 50 && edge >= 35 && bot < 35 && st.LargeTrades >= 5 && (st.CopyClosedTrades < 5 || st.CopyROI > 0) {
		return "B", "positive edge, prompt before auto-follow"
	}
	if smart >= 28 && edge >= 20 && bot < 45 {
		return "C", "watchlist"
	}
	return "D", "weak edge"
}

func followAction(tier string, st WalletStats, bot, smart float64) string {
	switch tier {
	case "A":
		if bot < 18 && smart >= 80 && st.ClosedPositions >= 20 && st.CopyClosedTrades >= 8 && st.CopyROI >= 10 && st.CopyWinRate >= 50 {
			return "auto-small"
		}
		return "prompt"
	case "B":
		return "prompt"
	case "C":
		return "watch"
	case "BOT":
		return "reject-bot"
	default:
		return "reject"
	}
}

func riskFlags(st WalletStats, bot float64) []string {
	var flags []string
	if bot >= 45 {
		flags = append(flags, "bot_like_flow")
	}
	if st.FixedAmountRatio >= 0.45 {
		flags = append(flags, "fixed_amount")
	}
	if st.FixedPriceRatio >= 0.35 {
		flags = append(flags, "fixed_price")
	}
	if st.ExtremePriceRatio >= 0.35 {
		flags = append(flags, "extreme_price_heavy")
	}
	if st.MaxTradesPerMinute >= 6 {
		flags = append(flags, "burst_trading")
	}
	if st.OppositeSideMarkets >= 2 {
		flags = append(flags, "opposite_side_same_market")
	}
	if st.ClosedPositions > 0 && st.ClosedCapital < 100 {
		flags = append(flags, "small_closed_capital")
	}
	if st.CopyClosedTrades >= 5 && st.CopyROI <= 0 {
		flags = append(flags, "negative_copy_sim")
	}
	if st.CopyOpenCost > st.CopyCapital && st.CopyOpenCost >= 50 {
		flags = append(flags, "open_copy_exposure")
	}
	return flags
}

func strengthFlags(st WalletStats, smart float64) []string {
	var flags []string
	if smart >= 70 {
		flags = append(flags, "high_smart_score")
	}
	if st.ClosedROI >= 10 && st.ClosedPositions >= 10 {
		flags = append(flags, "profitable_closed_book")
	}
	if st.ClosedWinRate >= 55 && st.ClosedPositions >= 10 {
		flags = append(flags, "above_average_win_rate")
	}
	if st.CopyClosedTrades >= 5 && st.CopyROI >= 5 {
		flags = append(flags, "positive_copy_sim")
	}
	if st.TopCategoryRatio >= 0.45 && st.TopCategory != "" {
		flags = append(flags, "category_focus:"+st.TopCategory)
	}
	if st.TargetTrades >= 5 && st.TargetTradeRatio >= 0.20 {
		flags = append(flags, fmt.Sprintf("target_focus:%.0f%%", st.TargetTradeRatio*100))
	}
	if st.LargeTrades >= 10 {
		flags = append(flags, "large_trade_sample")
	}
	return flags
}

func coarseCategory(text string) string {
	text = strings.ToLower(text)
	switch {
	case strings.Contains(text, "nba"), strings.Contains(text, "wnba"), strings.Contains(text, "basketball"), strings.Contains(text, "spurs"), strings.Contains(text, "thunder"):
		return "basketball"
	case strings.Contains(text, "football"), strings.Contains(text, "soccer"), strings.Contains(text, " fc "), strings.Contains(text, "premier league"), strings.Contains(text, "fifwc"), strings.Contains(text, "fifa world cup"), strings.Contains(text, "arsenal"), strings.Contains(text, "burnley"):
		return "soccer"
	case strings.Contains(text, "tennis"), strings.Contains(text, "open:"):
		return "tennis"
	case strings.Contains(text, "dota"), strings.Contains(text, "lol"), strings.Contains(text, "cs2"), strings.Contains(text, "valorant"):
		return "esports"
	case strings.Contains(text, "trump"), strings.Contains(text, "election"), strings.Contains(text, "senate"), strings.Contains(text, "president"):
		return "politics"
	case strings.Contains(text, "bitcoin"), strings.Contains(text, "btc"), strings.Contains(text, "ethereum"), strings.Contains(text, "crypto"):
		return "crypto"
	default:
		return "other"
	}
}

func maxBucketRatio(buckets map[string]int, total int) float64 {
	if total <= 0 {
		return 0
	}
	var max int
	for _, n := range buckets {
		if n > max {
			max = n
		}
	}
	return float64(max) / float64(total)
}

func topCategory(buckets map[string]int, total int) (string, float64) {
	if total <= 0 {
		return "", 0
	}
	var top string
	var topN int
	for k, n := range buckets {
		if n > topN || (n == topN && k < top) {
			top = k
			topN = n
		}
	}
	return top, float64(topN) / float64(total)
}

func SortScores(scores []WalletScore) {
	tierRank := map[string]int{"A": 5, "B": 4, "C": 3, "D": 2, "BOT": 1}
	sort.Slice(scores, func(i, j int) bool {
		ri, rj := tierRank[scores[i].Tier], tierRank[scores[j].Tier]
		if ri != rj {
			return ri > rj
		}
		return scores[i].Score > scores[j].Score
	})
}

func TierAllowed(tier, minTier string) bool {
	order := map[string]int{"A": 4, "B": 3, "C": 2, "D": 1, "BOT": 0}
	if minTier == "" {
		minTier = "C"
	}
	return order[strings.ToUpper(tier)] >= order[strings.ToUpper(minTier)]
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
