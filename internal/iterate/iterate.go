package iterate

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/journal"
)

type SportBreakdown struct {
	Sport    string
	Trades   int
	Wins     int
	Losses   int
	WinRate  float64
	PnLUSD   float64
	AvgEntry float64
	AvgHeld  int
}

type ExitBreakdown struct {
	Reason string
	Count  int
	PnLUSD float64
	AvgPnL float64
}

type PriceBand struct {
	Label   string
	Min     float64
	Max     float64
	Trades  int
	Wins    int
	WinRate float64
	PnLUSD  float64
}

type DayStats struct {
	Day    string
	Trades int
	WR     float64
	PnL    float64
}

type IterationReport struct {
	Day           string
	WindowDays    int
	TotalTrades   int
	TotalWins     int
	TotalLosses   int
	WinRate       float64
	CumulativePnL float64
	AvgPnLPerDay  float64
	AvgHeldSec    int

	DailyBreakdown []DayStats
	SportBreakdown []SportBreakdown
	ExitBreakdown  []ExitBreakdown
	PriceBands     []PriceBand

	SLHitRate      float64
	TPHitRate      float64
	TimeoutRate    float64
	SettleRate     float64

	Suggestions []string
}

func classifySport(question string) string {
	q := strings.ToLower(question)
	switch {
	case strings.Contains(q, "lck") || strings.Contains(q, "lpl") || strings.Contains(q, "lol") ||
		strings.Contains(q, " vs ") && (strings.Contains(q, "t1") || strings.Contains(q, "gen.g") || strings.Contains(q, "hanwha")):
		return "LoL"
	case strings.Contains(q, "dota"):
		return "Dota2"
	case strings.Contains(q, "lakers") || strings.Contains(q, "celtics") || strings.Contains(q, "nba") ||
		strings.Contains(q, "rockets") || strings.Contains(q, "warriors") || strings.Contains(q, "thunder") ||
		strings.Contains(q, "nuggets") || strings.Contains(q, "wolves") || strings.Contains(q, "spurs") ||
		strings.Contains(q, "blazers") || strings.Contains(q, "cavaliers") || strings.Contains(q, "pistons") ||
		strings.Contains(q, "magic") || strings.Contains(q, "hawks") || strings.Contains(q, "knicks") ||
		strings.Contains(q, "raptors") || strings.Contains(q, "suns") || strings.Contains(q, "76ers") ||
		strings.Contains(q, "bucks") || strings.Contains(q, "heat") || strings.Contains(q, "nets") ||
		strings.Contains(q, "bulls") || strings.Contains(q, "moneyline"):
		return "NBA"
	case strings.Contains(q, "wta") || strings.Contains(q, "atp") || strings.Contains(q, "tennis"):
		return "Tennis"
	case strings.Contains(q, "epl") || strings.Contains(q, "premier league") || strings.Contains(q, "la liga") ||
		strings.Contains(q, "serie a") || strings.Contains(q, "bundesliga") || strings.Contains(q, "champions league"):
		return "Football"
	default:
		return "Other"
	}
}

func Analyze(journalDir string, windowDays int) (*IterationReport, error) {
	now := time.Now().In(journal.SGT)

	var allTrades []journal.TradeRecord
	var dailyStats []DayStats

	for i := 0; i < windowDays; i++ {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		trades, err := journal.Read(journalDir, day)
		if err != nil {
			continue
		}
		var autoTrades []journal.TradeRecord
		for _, t := range trades {
			if t.SignalSource == "manual" {
				continue
			}
			autoTrades = append(autoTrades, t)
		}
		allTrades = append(allTrades, autoTrades...)

		ds := DayStats{Day: day, Trades: len(autoTrades)}
		var w, l int
		for _, t := range autoTrades {
			net := t.NetPnLUSD
			if net == 0 && t.EntryFeeUSD == 0 {
				net = t.PnLUSD
			}
			ds.PnL += net
			if net > 0 {
				w++
			} else if net < 0 {
				l++
			}
		}
		if w+l > 0 {
			ds.WR = float64(w) / float64(w+l)
		}
		dailyStats = append(dailyStats, ds)
	}

	sort.Slice(dailyStats, func(i, j int) bool { return dailyStats[i].Day < dailyStats[j].Day })

	r := &IterationReport{
		Day:            now.AddDate(0, 0, -1).Format("2006-01-02"),
		WindowDays:     windowDays,
		TotalTrades:    len(allTrades),
		DailyBreakdown: dailyStats,
	}

	if len(allTrades) == 0 {
		r.Suggestions = append(r.Suggestions, "无成交数据，检查 daemon 是否正常运行")
		return r, nil
	}

	sportMap := map[string]*SportBreakdown{}
	exitMap := map[string]*ExitBreakdown{}
	bands := []PriceBand{
		{Label: "0.05-0.15", Min: 0.05, Max: 0.15},
		{Label: "0.15-0.25", Min: 0.15, Max: 0.25},
		{Label: "0.25-0.35", Min: 0.25, Max: 0.35},
		{Label: "0.35-0.50", Min: 0.35, Max: 0.50},
		{Label: "0.50-0.70", Min: 0.50, Max: 0.70},
		{Label: "0.70-1.00", Min: 0.70, Max: 1.00},
	}

	var totalHeld int
	var slCount, tpCount, timeoutCount, settleCount int

	for _, t := range allTrades {
		net := t.NetPnLUSD
		if net == 0 && t.EntryFeeUSD == 0 {
			net = t.PnLUSD
		}

		r.CumulativePnL += net
		totalHeld += t.HeldSec

		if net > 0 {
			r.TotalWins++
		} else if net < 0 {
			r.TotalLosses++
		}

		sport := classifySport(t.Question)
		sb, ok := sportMap[sport]
		if !ok {
			sb = &SportBreakdown{Sport: sport}
			sportMap[sport] = sb
		}
		sb.Trades++
		sb.PnLUSD += net
		sb.AvgEntry += t.EntryMid
		sb.AvgHeld += t.HeldSec
		if net > 0 {
			sb.Wins++
		} else if net < 0 {
			sb.Losses++
		}

		reason := t.ExitReason
		if reason == "" {
			reason = "unknown"
		}
		eb, ok := exitMap[reason]
		if !ok {
			eb = &ExitBreakdown{Reason: reason}
			exitMap[reason] = eb
		}
		eb.Count++
		eb.PnLUSD += net

		switch {
		case strings.Contains(reason, "sl"):
			slCount++
		case strings.Contains(reason, "tp"):
			tpCount++
		case strings.Contains(reason, "timeout"):
			timeoutCount++
		case strings.Contains(reason, "settle"):
			settleCount++
		}

		for i := range bands {
			if t.EntryMid >= bands[i].Min && t.EntryMid < bands[i].Max {
				bands[i].Trades++
				bands[i].PnLUSD += net
				if net > 0 {
					bands[i].Wins++
				}
				break
			}
		}
	}

	decided := r.TotalWins + r.TotalLosses
	if decided > 0 {
		r.WinRate = float64(r.TotalWins) / float64(decided)
	}
	r.AvgPnLPerDay = r.CumulativePnL / float64(windowDays)
	r.AvgHeldSec = totalHeld / len(allTrades)

	total := slCount + tpCount + timeoutCount + settleCount
	if total > 0 {
		r.SLHitRate = float64(slCount) / float64(total)
		r.TPHitRate = float64(tpCount) / float64(total)
		r.TimeoutRate = float64(timeoutCount) / float64(total)
		r.SettleRate = float64(settleCount) / float64(total)
	}

	for _, sb := range sportMap {
		if sb.Trades > 0 {
			sb.AvgEntry /= float64(sb.Trades)
			sb.AvgHeld /= sb.Trades
			d := sb.Wins + sb.Losses
			if d > 0 {
				sb.WinRate = float64(sb.Wins) / float64(d)
			}
		}
		r.SportBreakdown = append(r.SportBreakdown, *sb)
	}
	sort.Slice(r.SportBreakdown, func(i, j int) bool { return r.SportBreakdown[i].Trades > r.SportBreakdown[j].Trades })

	for _, eb := range exitMap {
		if eb.Count > 0 {
			eb.AvgPnL = eb.PnLUSD / float64(eb.Count)
		}
		r.ExitBreakdown = append(r.ExitBreakdown, *eb)
	}
	sort.Slice(r.ExitBreakdown, func(i, j int) bool { return r.ExitBreakdown[i].Count > r.ExitBreakdown[j].Count })

	for i := range bands {
		d := bands[i].Wins + (bands[i].Trades - bands[i].Wins)
		if d > 0 {
			bands[i].WinRate = float64(bands[i].Wins) / float64(bands[i].Trades)
		}
	}
	r.PriceBands = bands

	r.Suggestions = generateSuggestions(r)

	return r, nil
}

func generateSuggestions(r *IterationReport) []string {
	var s []string

	if r.WinRate < 0.40 && r.TotalTrades >= 10 {
		s = append(s, fmt.Sprintf("⚠️ 胜率 %.0f%% 偏低（<40%%），建议收紧入场条件或提高信号阈值", r.WinRate*100))
	}
	if r.WinRate > 0.65 && r.TotalTrades >= 10 {
		s = append(s, fmt.Sprintf("✅ 胜率 %.0f%% 偏高（>65%%），考虑提高单笔注额或放宽入场条件", r.WinRate*100))
	}

	if r.SLHitRate > 0.60 && r.TotalTrades >= 10 {
		s = append(s, fmt.Sprintf("⚠️ SL 触发率 %.0f%%（>60%%），SL 可能过紧，建议从当前值适度放宽", r.SLHitRate*100))
	}
	if r.TimeoutRate > 0.40 && r.TotalTrades >= 10 {
		s = append(s, fmt.Sprintf("⚠️ Timeout 占比 %.0f%%（>40%%），持仓时间可能过长或 TP 目标过高", r.TimeoutRate*100))
	}

	for _, sb := range r.SportBreakdown {
		if sb.Trades >= 5 && sb.WinRate < 0.30 {
			s = append(s, fmt.Sprintf("⚠️ %s 胜率 %.0f%%（%d笔），考虑暂停该品类或调整参数", sb.Sport, sb.WinRate*100, sb.Trades))
		}
		if sb.Trades >= 5 && sb.PnLUSD < -5 {
			s = append(s, fmt.Sprintf("⚠️ %s 累亏 $%.2f（%d笔），拖累整体表现", sb.Sport, sb.PnLUSD, sb.Trades))
		}
	}

	var bestBand, worstBand string
	var bestPnL, worstPnL float64
	for _, b := range r.PriceBands {
		if b.Trades >= 3 {
			if b.PnLUSD > bestPnL {
				bestPnL = b.PnLUSD
				bestBand = b.Label
			}
			if b.PnLUSD < worstPnL {
				worstPnL = b.PnLUSD
				worstBand = b.Label
			}
		}
	}
	if bestBand != "" && bestPnL > 0 {
		s = append(s, fmt.Sprintf("📈 最佳入场价带 %s（PnL +$%.2f），可考虑加权", bestBand, bestPnL))
	}
	if worstBand != "" && worstPnL < -3 {
		s = append(s, fmt.Sprintf("📉 最差入场价带 %s（PnL $%.2f），考虑收窄或排除", worstBand, worstPnL))
	}

	if r.AvgHeldSec > 0 {
		avgMin := r.AvgHeldSec / 60
		if avgMin < 3 {
			s = append(s, fmt.Sprintf("⚡ 平均持仓 %dm，出场极快，可能在抖动中被洗出", avgMin))
		}
	}

	var trend []float64
	for _, ds := range r.DailyBreakdown {
		if ds.Trades > 0 {
			trend = append(trend, ds.PnL)
		}
	}
	if len(trend) >= 3 {
		recent := trend[len(trend)-1]
		earlier := 0.0
		for _, v := range trend[:len(trend)-1] {
			earlier += v
		}
		earlier /= float64(len(trend) - 1)
		if recent < earlier-3 {
			s = append(s, fmt.Sprintf("📉 最近一天 PnL $%.2f 低于前几天均值 $%.2f，注意策略衰退", recent, earlier))
		}
		if recent > earlier+3 {
			s = append(s, fmt.Sprintf("📈 最近一天 PnL $%.2f 高于前几天均值 $%.2f，策略表现提升", recent, earlier))
		}
	}

	if len(s) == 0 {
		s = append(s, "✅ 各项指标正常，继续积累数据")
	}

	return s
}

func FormatMarkdown(r *IterationReport) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# 每日迭代报告 %s SGT\n\n", r.Day)
	fmt.Fprintf(&b, "**窗口**: %d 天 · **总成交**: %d 笔 · **累计 PnL**: $%.2f\n\n", r.WindowDays, r.TotalTrades, r.CumulativePnL)

	b.WriteString("## 总览\n\n")
	fmt.Fprintf(&b, "| 指标 | 值 |\n|------|----|\n")
	fmt.Fprintf(&b, "| 胜/负/平 | %d / %d / %d |\n", r.TotalWins, r.TotalLosses, r.TotalTrades-r.TotalWins-r.TotalLosses)
	fmt.Fprintf(&b, "| 胜率 | %.1f%% |\n", r.WinRate*100)
	fmt.Fprintf(&b, "| 日均 PnL | $%.2f |\n", r.AvgPnLPerDay)
	fmt.Fprintf(&b, "| 平均持仓 | %s |\n", humanizeSec(r.AvgHeldSec))
	fmt.Fprintf(&b, "| SL 触发率 | %.0f%% |\n", r.SLHitRate*100)
	fmt.Fprintf(&b, "| TP 触发率 | %.0f%% |\n", r.TPHitRate*100)
	fmt.Fprintf(&b, "| Timeout 率 | %.0f%% |\n", r.TimeoutRate*100)
	fmt.Fprintf(&b, "| Settlement 率 | %.0f%% |\n\n", r.SettleRate*100)

	b.WriteString("## 每日明细\n\n")
	b.WriteString("| 日期 | 笔数 | 胜率 | PnL |\n|------|------|------|-----|\n")
	for _, ds := range r.DailyBreakdown {
		pnl := fmt.Sprintf("$%.2f", ds.PnL)
		if ds.PnL > 0 {
			pnl = "+$" + fmt.Sprintf("%.2f", ds.PnL)
		}
		wr := "-"
		if ds.Trades > 0 {
			wr = fmt.Sprintf("%.0f%%", ds.WR*100)
		}
		fmt.Fprintf(&b, "| %s | %d | %s | %s |\n", ds.Day, ds.Trades, wr, pnl)
	}
	b.WriteString("\n")

	if len(r.SportBreakdown) > 0 {
		b.WriteString("## 品类分析\n\n")
		b.WriteString("| 品类 | 笔数 | 胜率 | PnL | 均价 | 均持仓 |\n|------|------|------|-----|------|--------|\n")
		for _, sb := range r.SportBreakdown {
			fmt.Fprintf(&b, "| %s | %d | %.0f%% | $%.2f | %.3f | %s |\n",
				sb.Sport, sb.Trades, sb.WinRate*100, sb.PnLUSD, sb.AvgEntry, humanizeSec(sb.AvgHeld))
		}
		b.WriteString("\n")
	}

	if len(r.ExitBreakdown) > 0 {
		b.WriteString("## 出场分析\n\n")
		b.WriteString("| 出场原因 | 次数 | PnL | 均 PnL |\n|----------|------|-----|--------|\n")
		for _, eb := range r.ExitBreakdown {
			fmt.Fprintf(&b, "| %s | %d | $%.2f | $%.4f |\n", eb.Reason, eb.Count, eb.PnLUSD, eb.AvgPnL)
		}
		b.WriteString("\n")
	}

	b.WriteString("## 入场价带分析\n\n")
	b.WriteString("| 价带 | 笔数 | 胜率 | PnL |\n|------|------|------|-----|\n")
	for _, pb := range r.PriceBands {
		if pb.Trades == 0 {
			continue
		}
		fmt.Fprintf(&b, "| %s | %d | %.0f%% | $%.2f |\n", pb.Label, pb.Trades, pb.WinRate*100, pb.PnLUSD)
	}
	b.WriteString("\n")

	b.WriteString("## 算法迭代建议\n\n")
	for _, s := range r.Suggestions {
		fmt.Fprintf(&b, "- %s\n", s)
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "---\n*生成时间: %s SGT*\n", time.Now().In(journal.SGT).Format("2006-01-02 15:04:05"))
	return b.String()
}

func FormatTelegram(r *IterationReport) string {
	var b strings.Builder

	sign := ""
	if r.CumulativePnL > 0 {
		sign = "+"
	}
	fmt.Fprintf(&b, "🔄 每日迭代 %s SGT\n", r.Day)
	fmt.Fprintf(&b, "窗口 %d天 · %d笔 · PnL %s$%.2f\n", r.WindowDays, r.TotalTrades, sign, r.CumulativePnL)
	fmt.Fprintf(&b, "胜率 %.0f%% · SL %.0f%% / TP %.0f%% / TO %.0f%%\n\n", r.WinRate*100, r.SLHitRate*100, r.TPHitRate*100, r.TimeoutRate*100)

	if len(r.SportBreakdown) > 0 {
		for _, sb := range r.SportBreakdown {
			s := ""
			if sb.PnLUSD > 0 {
				s = "+"
			}
			fmt.Fprintf(&b, "  %s: %d笔 %.0f%% %s$%.2f\n", sb.Sport, sb.Trades, sb.WinRate*100, s, sb.PnLUSD)
		}
		b.WriteString("\n")
	}

	if len(r.Suggestions) > 0 {
		b.WriteString("💡 建议:\n")
		for _, s := range r.Suggestions {
			fmt.Fprintf(&b, "  %s\n", s)
		}
	}
	return b.String()
}

func humanizeSec(sec int) string {
	if sec <= 0 {
		return "0s"
	}
	d := time.Duration(sec) * time.Second
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", sec)
	case d < time.Hour:
		return fmt.Sprintf("%dm%02ds", sec/60, sec%60)
	default:
		h := sec / 3600
		m := (sec % 3600) / 60
		return fmt.Sprintf("%dh%02dm", h, m)
	}
}

func Abs(f float64) float64 { return math.Abs(f) }
