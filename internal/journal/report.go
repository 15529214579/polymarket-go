package journal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// DailySummary is the aggregate over one SGT day's TradeRecords. Built by
// Summarize; rendered by FormatTelegram for the cron-pushed message.
type DailySummary struct {
	Day             string
	Trades          int
	Wins            int
	Losses          int
	Breakevens      int
	WinRate         float64 // 0..1, undefined if Trades==0
	RealizedPnLUSD  float64
	GrossWinUSD     float64
	GrossLossUSD    float64
	AvgPnLUSD       float64
	BiggestWinUSD   float64
	BiggestLossUSD  float64
	AvgHeldSec      int
	ExitReasonCount map[string]int
	ManualCount     int
	AutoCount       int
}

// Summarize buckets a slice of trades into a DailySummary. Wins are pnl>0,
// losses pnl<0, breakevens exactly 0 (slippage 0 paper can produce these).
func Summarize(day string, trades []TradeRecord) DailySummary {
	s := DailySummary{Day: day, ExitReasonCount: map[string]int{}}
	if len(trades) == 0 {
		return s
	}
	var heldTotal int
	for _, t := range trades {
		s.Trades++
		s.RealizedPnLUSD += t.PnLUSD
		switch {
		case t.PnLUSD > 0:
			s.Wins++
			s.GrossWinUSD += t.PnLUSD
			if t.PnLUSD > s.BiggestWinUSD {
				s.BiggestWinUSD = t.PnLUSD
			}
		case t.PnLUSD < 0:
			s.Losses++
			s.GrossLossUSD += t.PnLUSD
			if t.PnLUSD < s.BiggestLossUSD {
				s.BiggestLossUSD = t.PnLUSD
			}
		default:
			s.Breakevens++
		}
		heldTotal += t.HeldSec
		if t.ExitReason != "" {
			s.ExitReasonCount[t.ExitReason]++
		}
		switch t.SignalSource {
		case "manual":
			s.ManualCount++
		case "auto", "":
			s.AutoCount++
		}
	}
	if s.Trades > 0 {
		s.AvgPnLUSD = s.RealizedPnLUSD / float64(s.Trades)
		s.AvgHeldSec = heldTotal / s.Trades
		decided := s.Wins + s.Losses
		if decided > 0 {
			s.WinRate = float64(s.Wins) / float64(decided)
		}
	}
	return s
}

// FormatTelegram renders a Markdown-light summary suitable for sendMessage's
// default parse mode (plain text — no escaping needed).
func FormatTelegram(s DailySummary) string {
	var b strings.Builder
	pnlSign := ""
	if s.RealizedPnLUSD > 0 {
		pnlSign = "+"
	}
	fmt.Fprintf(&b, "📊 polymarket-go 日结 %s SGT\n", s.Day)
	if s.Trades == 0 {
		b.WriteString("无成交。\n")
		return b.String()
	}
	fmt.Fprintf(&b, "• 实现 PnL: %s%.4f USDC\n", pnlSign, s.RealizedPnLUSD)
	fmt.Fprintf(&b, "• 成交 %d 笔  胜 %d / 负 %d / 平 %d  (胜率 %.0f%%)\n",
		s.Trades, s.Wins, s.Losses, s.Breakevens, s.WinRate*100)
	fmt.Fprintf(&b, "• 平均 PnL/笔 %.4f USDC\n", s.AvgPnLUSD)
	if s.Wins > 0 || s.Losses > 0 {
		fmt.Fprintf(&b, "• 最大胜 +%.4f / 最大负 %.4f\n",
			s.BiggestWinUSD, s.BiggestLossUSD)
	}
	fmt.Fprintf(&b, "• 平均持仓 %s\n", humanizeSec(s.AvgHeldSec))
	if s.ManualCount > 0 || s.AutoCount > 0 {
		fmt.Fprintf(&b, "• 来源: auto %d / manual %d\n", s.AutoCount, s.ManualCount)
	}
	if len(s.ExitReasonCount) > 0 {
		// Stable ordering for snapshot tests + readability.
		keys := make([]string, 0, len(s.ExitReasonCount))
		for k := range s.ExitReasonCount {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString("• 出场: ")
		for i, k := range keys {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%s×%d", k, s.ExitReasonCount[k])
		}
		b.WriteString("\n")
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
		return fmt.Sprintf("%dh%02dm", sec/3600, (sec%3600)/60)
	}
}
