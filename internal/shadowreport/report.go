// Package shadowreport compares actual exits with unexecuted policy alternatives.
package shadowreport

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Observation struct {
	ObservedAt    time.Time
	EntryTime     time.Time
	ActualCloseAt time.Time
	PosID         string
	Policy        string
	Question      string
	Source        string
	HoldProfile   string
	NetPnLUSD     float64
}

type Group struct {
	Name        string  `json:"name"`
	Samples     int     `json:"samples"`
	Positive    int     `json:"positive"`
	PostActual  int     `json:"post_actual_close"`
	NetPnLUSD   float64 `json:"net_pnl_usd"`
	AverageNet  float64 `json:"average_net_pnl_usd"`
	PositivePct float64 `json:"positive_pct"`
}

type PairedComparison struct {
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	EarlierPolicy    string  `json:"earlier_policy"`
	LaterPolicy      string  `json:"later_policy"`
	Samples          int     `json:"samples"`
	LaterBetter      int     `json:"later_better"`
	Equal            int     `json:"equal"`
	PostActual       int     `json:"later_post_actual_close"`
	EarlierNetPnLUSD float64 `json:"earlier_net_pnl_usd"`
	LaterNetPnLUSD   float64 `json:"later_net_pnl_usd"`
	NetUpliftUSD     float64 `json:"net_uplift_usd"`
	AverageUpliftUSD float64 `json:"average_uplift_usd"`
	LaterBetterPct   float64 `json:"later_better_pct"`
}

type Report struct {
	GeneratedAt    time.Time          `json:"generated_at"`
	Samples        int                `json:"samples"`
	ByPolicy       []Group            `json:"by_policy"`
	ByCategory     []Group            `json:"by_category_policy"`
	PairedTimeouts []PairedComparison `json:"paired_timeouts"`
}

func Analyze(observations []Observation) Report {
	deduped := map[string]Observation{}
	for _, obs := range observations {
		if obs.PosID == "" || obs.Policy == "" {
			continue
		}
		entry := obs.EntryTime
		if entry.IsZero() && !obs.ObservedAt.IsZero() {
			entry = obs.ObservedAt
		}
		key := fmt.Sprintf("%s|%s|%d", obs.PosID, obs.Policy, entry.Unix())
		previous, exists := deduped[key]
		if !exists || (!obs.ObservedAt.IsZero() && obs.ObservedAt.Before(previous.ObservedAt)) {
			deduped[key] = obs
		}
	}

	policy := map[string]*Group{}
	category := map[string]*Group{}
	for _, obs := range deduped {
		acc(policy, obs.Policy, obs)
		acc(category, classify(obs.Question)+" / "+obs.Policy, obs)
	}
	return Report{
		GeneratedAt:    time.Now(),
		Samples:        len(deduped),
		ByPolicy:       finish(policy),
		ByCategory:     finish(category),
		PairedTimeouts: analyzeTimeoutPairs(deduped),
	}
}

func analyzeTimeoutPairs(observations map[string]Observation) []PairedComparison {
	byPosition := map[string]map[string]Observation{}
	for _, obs := range observations {
		entry := obs.EntryTime
		if entry.IsZero() {
			entry = obs.ObservedAt
		}
		key := fmt.Sprintf("%s|%d", obs.PosID, entry.Unix())
		if byPosition[key] == nil {
			byPosition[key] = map[string]Observation{}
		}
		byPosition[key][obs.Policy] = obs
	}
	pairs := [][2]string{
		{"timeout_10m", "timeout_20m"},
		{"timeout_20m", "timeout_30m"},
		{"timeout_30m", "timeout_45m"},
		{"timeout_45m", "timeout_60m"},
	}
	groups := map[string]*PairedComparison{}
	for _, policies := range byPosition {
		for _, pair := range pairs {
			earlier, earlierOK := policies[pair[0]]
			later, laterOK := policies[pair[1]]
			if !earlierOK || !laterOK {
				continue
			}
			category := classify(later.Question)
			key := category + "|" + pair[0] + "|" + pair[1]
			group := groups[key]
			if group == nil {
				group = &PairedComparison{
					Name:     category + " / " + pair[0] + " -> " + pair[1],
					Category: category, EarlierPolicy: pair[0], LaterPolicy: pair[1],
				}
				groups[key] = group
			}
			group.Samples++
			group.EarlierNetPnLUSD += earlier.NetPnLUSD
			group.LaterNetPnLUSD += later.NetPnLUSD
			uplift := later.NetPnLUSD - earlier.NetPnLUSD
			switch {
			case uplift > 1e-9:
				group.LaterBetter++
			case uplift >= -1e-9:
				group.Equal++
			}
			if !later.ActualCloseAt.IsZero() && !later.ObservedAt.Before(later.ActualCloseAt) {
				group.PostActual++
			}
		}
	}
	out := make([]PairedComparison, 0, len(groups))
	for _, group := range groups {
		group.NetUpliftUSD = group.LaterNetPnLUSD - group.EarlierNetPnLUSD
		if group.Samples > 0 {
			group.AverageUpliftUSD = group.NetUpliftUSD / float64(group.Samples)
			group.LaterBetterPct = float64(group.LaterBetter) / float64(group.Samples) * 100
		}
		out = append(out, *group)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func acc(groups map[string]*Group, name string, obs Observation) {
	group := groups[name]
	if group == nil {
		group = &Group{Name: name}
		groups[name] = group
	}
	group.Samples++
	group.NetPnLUSD += obs.NetPnLUSD
	if obs.NetPnLUSD > 0 {
		group.Positive++
	}
	if !obs.ActualCloseAt.IsZero() && !obs.ObservedAt.Before(obs.ActualCloseAt) {
		group.PostActual++
	}
}

func finish(groups map[string]*Group) []Group {
	out := make([]Group, 0, len(groups))
	for _, group := range groups {
		if group.Samples > 0 {
			group.AverageNet = group.NetPnLUSD / float64(group.Samples)
			group.PositivePct = float64(group.Positive) / float64(group.Samples) * 100
		}
		out = append(out, *group)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func classify(question string) string {
	text := strings.ToLower(question)
	switch {
	case strings.Contains(text, "exact score"), strings.Contains(text, "correct score"):
		return "football_score"
	case strings.Contains(text, "dota"), strings.Contains(text, "lol:"), strings.Contains(text, "league of legends"), strings.Contains(text, "cs2"), strings.Contains(text, "counter-strike"), strings.Contains(text, "valorant"):
		return "esports"
	case strings.Contains(text, "nba"), strings.Contains(text, "wnba"), strings.Contains(text, "basketball"):
		return "basketball"
	case strings.Contains(text, " fc"), strings.Contains(text, "soccer"), strings.Contains(text, "football"):
		return "football"
	case strings.TrimSpace(text) == "":
		return "legacy_unknown"
	default:
		return "other"
	}
}

func FormatMarkdown(report Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Smartmoney Exit Shadow\n\nGenerated: %s\n\n- Independent position-policy samples: %d\n", report.GeneratedAt.Format(time.RFC3339), report.Samples)
	writeGroups(&b, "By Policy", report.ByPolicy)
	writeGroups(&b, "By Category / Policy", report.ByCategory)
	writePairs(&b, report.PairedTimeouts)
	return b.String()
}

func writePairs(b *strings.Builder, groups []PairedComparison) {
	b.WriteString("\n## Matched Timeout Comparisons\n\n")
	b.WriteString("Only positions observed at both horizons are compared. Uplift is later net PnL minus earlier net PnL for the same position.\n\n")
	b.WriteString("| Category / Pair | Matched | Later better | Later post-actual | Earlier net | Later net | Uplift | Avg uplift | Better rate |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, group := range groups {
		fmt.Fprintf(b, "| %s | %d | %d | %d | %+.2f | %+.2f | %+.2f | %+.2f | %.1f%% |\n",
			group.Name, group.Samples, group.LaterBetter, group.PostActual, group.EarlierNetPnLUSD, group.LaterNetPnLUSD,
			group.NetUpliftUSD, group.AverageUpliftUSD, group.LaterBetterPct)
	}
}

func writeGroups(b *strings.Builder, title string, groups []Group) {
	fmt.Fprintf(b, "\n## %s\n\n", title)
	b.WriteString("| Cohort | Samples | Positive | Post-actual | Net | Avg net | Positive rate |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|\n")
	for _, group := range groups {
		fmt.Fprintf(b, "| %s | %d | %d | %d | %+.2f | %+.2f | %.1f%% |\n", group.Name, group.Samples, group.Positive, group.PostActual, group.NetPnLUSD, group.AverageNet, group.PositivePct)
	}
}
