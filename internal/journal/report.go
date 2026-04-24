package journal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// DailySummary is the aggregate over one SGT day's TradeRecords. Built by
// Summarize; rendered by FormatTelegram for the cron-pushed message.
//
// RealizedPnLUSD is net of fees when any record carries NetPnLUSD / fee
// fields (Phase 7.b ladder); older records with zero fee fields contribute
// gross PnL (net == gross in that case).
// SourceStats holds stats for one source (auto or manual).
type SourceStats struct {
	Count        int
	Wins         int
	Losses       int
	Breakevens   int
	WinRate      float64
	PnLUSD       float64
	AvgPnLUSD    float64
	BiggestWin   float64
	BiggestLoss  float64
	AvgHeldSec   int
	ExitReasons  map[string]int
}

type DailySummary struct {
	Day string

	// Auto-only headline stats (manual excluded from these).
	Trades          int
	Wins            int
	Losses          int
	Breakevens      int
	WinRate         float64
	RealizedPnLUSD  float64
	GrossPnLUSD     float64
	FeesUSD         float64
	GrossWinUSD     float64
	GrossLossUSD    float64
	AvgPnLUSD       float64
	BiggestWinUSD   float64
	BiggestLossUSD  float64
	AvgHeldSec      int
	ExitReasonCount map[string]int

	// Separate accounting.
	Auto   SourceStats
	Manual SourceStats
}

// Summarize buckets a slice of trades into a DailySummary. Wins are pnl>0,
// losses pnl<0, breakevens exactly 0 (slippage 0 paper can produce these).
func tradeNet(t TradeRecord) float64 {
	net := t.NetPnLUSD
	if net == 0 && t.EntryFeeUSD == 0 && t.ExitFeeUSD == 0 {
		net = t.PnLUSD
	}
	return net
}

func accSource(ss *SourceStats, t TradeRecord, net float64) {
	ss.Count++
	ss.PnLUSD += net
	switch {
	case net > 0:
		ss.Wins++
		if net > ss.BiggestWin {
			ss.BiggestWin = net
		}
	case net < 0:
		ss.Losses++
		if net < ss.BiggestLoss {
			ss.BiggestLoss = net
		}
	default:
		ss.Breakevens++
	}
	if t.ExitReason != "" {
		if ss.ExitReasons == nil {
			ss.ExitReasons = map[string]int{}
		}
		ss.ExitReasons[t.ExitReason]++
	}
}

func finalizeSource(ss *SourceStats, heldTotal int) {
	if ss.Count > 0 {
		ss.AvgPnLUSD = ss.PnLUSD / float64(ss.Count)
		ss.AvgHeldSec = heldTotal / ss.Count
		decided := ss.Wins + ss.Losses
		if decided > 0 {
			ss.WinRate = float64(ss.Wins) / float64(decided)
		}
	}
}

func Summarize(day string, trades []TradeRecord) DailySummary {
	s := DailySummary{Day: day, ExitReasonCount: map[string]int{}}
	if len(trades) == 0 {
		return s
	}
	var autoHeld, manualHeld int
	for _, t := range trades {
		net := tradeNet(t)
		isManual := t.SignalSource == "manual"
		if isManual {
			accSource(&s.Manual, t, net)
			manualHeld += t.HeldSec
		} else {
			accSource(&s.Auto, t, net)
			autoHeld += t.HeldSec
			// Headline stats = auto only.
			s.Trades++
			s.GrossPnLUSD += t.PnLUSD
			s.FeesUSD += t.EntryFeeUSD + t.ExitFeeUSD
			s.RealizedPnLUSD += net
			switch {
			case net > 0:
				s.Wins++
				s.GrossWinUSD += net
				if net > s.BiggestWinUSD {
					s.BiggestWinUSD = net
				}
			case net < 0:
				s.Losses++
				s.GrossLossUSD += net
				if net < s.BiggestLossUSD {
					s.BiggestLossUSD = net
				}
			default:
				s.Breakevens++
			}
			if t.ExitReason != "" {
				s.ExitReasonCount[t.ExitReason]++
			}
		}
	}
	finalizeSource(&s.Auto, autoHeld)
	finalizeSource(&s.Manual, manualHeld)
	if s.Trades > 0 {
		s.AvgPnLUSD = s.RealizedPnLUSD / float64(s.Trades)
		s.AvgHeldSec = autoHeld / s.Trades
		decided := s.Wins + s.Losses
		if decided > 0 {
			s.WinRate = float64(s.Wins) / float64(decided)
		}
	}
	return s
}

// FormatTelegram renders a Markdown-light summary suitable for sendMessage's
// default parse mode (plain text — no escaping needed). When both auto and
// manual trades exist, a per-source breakdown is appended.
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
	fmt.Fprintf(&b, "• 实现 PnL(净): %s%.4f USDC\n", pnlSign, s.RealizedPnLUSD)
	if s.FeesUSD > 0 {
		fmt.Fprintf(&b, "• 毛 PnL: %+.4f · 手续费 %.4f\n", s.GrossPnLUSD, s.FeesUSD)
	}
	fmt.Fprintf(&b, "• 成交 %d 笔  胜 %d / 负 %d / 平 %d  (胜率 %.0f%%)\n",
		s.Trades, s.Wins, s.Losses, s.Breakevens, s.WinRate*100)
	fmt.Fprintf(&b, "• 平均 PnL/笔 %.4f USDC\n", s.AvgPnLUSD)
	if s.Wins > 0 || s.Losses > 0 {
		fmt.Fprintf(&b, "• 最大胜 +%.4f / 最大负 %.4f\n",
			s.BiggestWinUSD, s.BiggestLossUSD)
	}
	fmt.Fprintf(&b, "• 平均持仓 %s\n", humanizeSec(s.AvgHeldSec))
	if len(s.ExitReasonCount) > 0 {
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
	if s.Manual.Count > 0 {
		fmt.Fprintf(&b, "\n👤 手动单独结算: %d笔 %+.4f USDC", s.Manual.Count, s.Manual.PnLUSD)
		if s.Manual.Wins > 0 || s.Manual.Losses > 0 {
			fmt.Fprintf(&b, "  胜%d/负%d", s.Manual.Wins, s.Manual.Losses)
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
