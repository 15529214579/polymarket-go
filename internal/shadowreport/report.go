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

type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	Samples     int       `json:"samples"`
	ByPolicy    []Group   `json:"by_policy"`
	ByCategory  []Group   `json:"by_category_policy"`
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
		GeneratedAt: time.Now(),
		Samples:     len(deduped),
		ByPolicy:    finish(policy),
		ByCategory:  finish(category),
	}
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
	case strings.Contains(text, "dota"), strings.Contains(text, "lol:"), strings.Contains(text, "league of legends"):
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
	return b.String()
}

func writeGroups(b *strings.Builder, title string, groups []Group) {
	fmt.Fprintf(b, "\n## %s\n\n", title)
	b.WriteString("| Cohort | Samples | Positive | Post-actual | Net | Avg net | Positive rate |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|\n")
	for _, group := range groups {
		fmt.Fprintf(b, "| %s | %d | %d | %d | %+.2f | %+.2f | %.1f%% |\n", group.Name, group.Samples, group.Positive, group.PostActual, group.NetPnLUSD, group.AverageNet, group.PositivePct)
	}
}
