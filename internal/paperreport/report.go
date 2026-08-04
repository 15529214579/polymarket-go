// Package paperreport summarizes simulated positions and execution records.
package paperreport

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

type Cohort struct {
	Name       string  `json:"name"`
	Positions  int     `json:"positions"`
	Records    int     `json:"records"`
	CapitalUSD float64 `json:"capital_usd"`
	GrossPnL   float64 `json:"gross_pnl_usd"`
	FeesUSD    float64 `json:"fees_usd"`
	NetPnL     float64 `json:"net_pnl_usd"`
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	Flat       int     `json:"flat"`
}

type PnLScope struct {
	Closed                  Cohort  `json:"closed"`
	OpenPositions           int     `json:"open_positions"`
	OpenExposureUSD         float64 `json:"open_exposure_usd"`
	OpenEntryFeesUSD        float64 `json:"open_entry_fees_usd"`
	ConservativeOpenNetPnL  float64 `json:"conservative_open_net_pnl_usd"`
	ConservativeTotalNetPnL float64 `json:"conservative_total_net_pnl_usd"`
}

type Report struct {
	GeneratedAt                time.Time `json:"generated_at"`
	Records                    int       `json:"records"`
	ClosedPositions            int       `json:"closed_positions"`
	GrossPnLUSD                float64   `json:"gross_pnl_usd"`
	EntryFeesUSD               float64   `json:"entry_fees_usd"`
	ExitFeesUSD                float64   `json:"exit_fees_usd"`
	FeesUSD                    float64   `json:"fees_usd"`
	RealizedNetPnLUSD          float64   `json:"realized_net_pnl_usd"`
	OpenPositions              int       `json:"open_positions"`
	OpenExposureUSD            float64   `json:"open_exposure_usd"`
	OpenEntryFeesUSD           float64   `json:"open_entry_fees_usd"`
	ConservativeOpenNetPnLUSD  float64   `json:"conservative_open_net_pnl_usd"`
	ConservativeTotalNetPnLUSD float64   `json:"conservative_total_net_pnl_usd"`
	Tradable                   PnLScope  `json:"tradable"`
	BroadCollection            PnLScope  `json:"broad_collection"`
	ByPolicy                   []Cohort  `json:"by_policy"`
	ByStake                    []Cohort  `json:"by_stake"`
	ByStrategy                 []Cohort  `json:"by_strategy"`
	ByCostModel                []Cohort  `json:"by_cost_model"`
	BySignalSource             []Cohort  `json:"by_signal_source"`
}

type WalletPolicyConfig struct {
	MinPositions  int     `json:"min_independent_samples"`
	PromoteMinNet float64 `json:"promote_min_net_usd"`
	PromoteMinROI float64 `json:"promote_min_roi_pct"`
	DemoteMaxNet  float64 `json:"demote_max_net_usd"`
}

type WalletPerformance struct {
	Wallet             string  `json:"wallet,omitempty"`
	Source             string  `json:"source"`
	IndependentSamples int     `json:"independent_samples"`
	Positions          int     `json:"positions"`
	Records            int     `json:"records"`
	CapitalUSD         float64 `json:"capital_usd"`
	GrossPnL           float64 `json:"gross_pnl_usd"`
	FeesUSD            float64 `json:"fees_usd"`
	NetPnL             float64 `json:"net_pnl_usd"`
	ROI                float64 `json:"roi_pct"`
	Wins               int     `json:"wins"`
	Losses             int     `json:"losses"`
	Flat               int     `json:"flat"`
	Decision           string  `json:"decision"`
	Reason             string  `json:"reason"`
}

type WalletPolicyReport struct {
	GeneratedAt time.Time           `json:"generated_at"`
	Config      WalletPolicyConfig  `json:"config"`
	Wallets     []WalletPerformance `json:"wallets"`
	Promoted    int                 `json:"promoted"`
	Demoted     int                 `json:"demoted"`
	Unresolved  int                 `json:"unresolved"`
}

type positionResult struct {
	id        string
	records   int
	capital   float64
	gross     float64
	entryFees float64
	exitFees  float64
	net       float64
	policy    string
	strategy  string
	source    string
	market    string
	costed    bool
}

func Analyze(trades []journal.TradeRecord, open []strategy.Position) Report {
	report := Report{
		GeneratedAt:     time.Now(),
		Tradable:        PnLScope{Closed: Cohort{Name: "tradable"}},
		BroadCollection: PnLScope{Closed: Cohort{Name: "broad_collection"}},
	}
	positions := groupTrades(trades, &report)
	report.ClosedPositions = len(positions)
	report.FeesUSD = report.EntryFeesUSD + report.ExitFeesUSD

	policy := map[string]*Cohort{}
	stake := map[string]*Cohort{}
	strategyGroups := map[string]*Cohort{}
	cost := map[string]*Cohort{}
	sources := map[string]*Cohort{}
	for _, pos := range positions {
		policyName := pos.policy
		if policyName == "" {
			policyName = "legacy/unversioned"
		}
		strategyName := pos.strategy
		if strategyName == "" {
			strategyName = "other"
		}
		sourceName := pos.source
		if sourceName == "" {
			sourceName = "unknown"
		}
		costName := "legacy-no-cost"
		if pos.costed {
			costName = "fee-aware"
		}
		acc(policy, policyName, pos)
		acc(stake, stakeBucket(pos.capital), pos)
		acc(strategyGroups, strategyName, pos)
		acc(cost, costName, pos)
		acc(sources, sourceName, pos)
		if isCollectionStrategy(strategyName) {
			accCohort(&report.BroadCollection.Closed, pos)
		} else {
			accCohort(&report.Tradable.Closed, pos)
		}
	}
	report.ByPolicy = cohorts(policy, false)
	report.ByStake = cohorts(stake, false)
	report.ByStrategy = cohorts(strategyGroups, false)
	report.ByCostModel = cohorts(cost, false)
	report.BySignalSource = cohorts(sources, true)

	for _, pos := range open {
		report.OpenPositions++
		report.OpenExposureUSD += pos.SizeUSD
		remainingFee := pos.OpenFeeUSD - pos.EntryFeeChargedUSD
		if remainingFee < 0 {
			remainingFee = 0
		}
		report.OpenEntryFeesUSD += remainingFee
		scope := &report.Tradable
		if isCollectionPosition(pos) {
			scope = &report.BroadCollection
		}
		scope.OpenPositions++
		scope.OpenExposureUSD += pos.SizeUSD
		scope.OpenEntryFeesUSD += remainingFee
	}
	report.ConservativeOpenNetPnLUSD = -(report.OpenExposureUSD + report.OpenEntryFeesUSD)
	report.ConservativeTotalNetPnLUSD = report.RealizedNetPnLUSD + report.ConservativeOpenNetPnLUSD
	finalizeScope(&report.Tradable)
	finalizeScope(&report.BroadCollection)
	return report
}

func groupTrades(trades []journal.TradeRecord, report *Report) map[string]*positionResult {
	positions := make(map[string]*positionResult)
	for _, tr := range trades {
		if report != nil {
			report.Records++
			report.GrossPnLUSD += tr.PnLUSD
			report.EntryFeesUSD += tr.EntryFeeUSD
			report.ExitFeesUSD += tr.ExitFeeUSD
			report.RealizedNetPnLUSD += tradeNet(tr)
		}

		id := positionKey(tr)
		pos := positions[id]
		if pos == nil {
			pos = &positionResult{id: id}
			positions[id] = pos
		}
		pos.records++
		pos.capital += tr.SizeUSD
		pos.gross += tr.PnLUSD
		pos.entryFees += tr.EntryFeeUSD
		pos.exitFees += tr.ExitFeeUSD
		pos.net += tradeNet(tr)
		pos.costed = pos.costed || tr.EntryFeeUSD > 0 || tr.ExitFeeUSD > 0
		pos.policy = mergeLabel(pos.policy, strings.TrimSpace(tr.PolicyVersion))
		pos.source = mergeSource(pos.source, normalizedSource(tr.SignalSource))
		pos.strategy = mergeStrategy(pos.strategy, strategyLabel(tr))
		pos.market = mergeLabel(pos.market, strings.TrimSpace(tr.Market))
	}
	return positions
}

func AnalyzeWalletPolicy(trades []journal.TradeRecord, aliases map[string]string, cfg WalletPolicyConfig) WalletPolicyReport {
	if cfg.MinPositions <= 0 {
		cfg.MinPositions = 10
	}
	if cfg.PromoteMinNet == 0 {
		cfg.PromoteMinNet = 5
	}
	if cfg.PromoteMinROI == 0 {
		cfg.PromoteMinROI = 2
	}
	if cfg.DemoteMaxNet == 0 {
		cfg.DemoteMaxNet = -5
	}

	result := WalletPolicyReport{GeneratedAt: time.Now(), Config: cfg}
	byWallet := map[string]*WalletPerformance{}
	byWalletSample := map[string]map[string]float64{}
	for _, pos := range groupTrades(trades, nil) {
		wallet := walletFromSource(pos.source, aliases)
		if wallet == "" && !strings.HasPrefix(strings.ToLower(pos.source), "copytrade") {
			continue
		}
		key := wallet
		if key == "" {
			key = "unresolved:" + pos.source
		}
		row := byWallet[key]
		if row == nil {
			row = &WalletPerformance{Wallet: wallet, Source: pos.source}
			byWallet[key] = row
		}
		row.Positions++
		row.Records += pos.records
		row.CapitalUSD += pos.capital
		row.GrossPnL += pos.gross
		row.FeesUSD += pos.entryFees + pos.exitFees
		row.NetPnL += pos.net
		sample := strings.TrimSpace(pos.market)
		if sample == "" || sample == "mixed" {
			sample = "position:" + pos.id
		}
		if byWalletSample[key] == nil {
			byWalletSample[key] = map[string]float64{}
		}
		byWalletSample[key][sample] += pos.net
	}

	for key, row := range byWallet {
		for _, sampleNet := range byWalletSample[key] {
			row.IndependentSamples++
			switch {
			case sampleNet > 0:
				row.Wins++
			case sampleNet < 0:
				row.Losses++
			default:
				row.Flat++
			}
		}
		if row.CapitalUSD > 0 {
			row.ROI = row.NetPnL / row.CapitalUSD * 100
		}
		switch {
		case row.Wallet == "":
			row.Decision = "unresolved"
			row.Reason = "full wallet address unavailable"
			result.Unresolved++
		case row.IndependentSamples < cfg.MinPositions:
			row.Decision = "collect"
			row.Reason = fmt.Sprintf("%d/%d independent wallet-market samples", row.IndependentSamples, cfg.MinPositions)
		case row.NetPnL <= cfg.DemoteMaxNet:
			row.Decision = "demote"
			row.Reason = fmt.Sprintf("net %+.2fU <= %+.2fU", row.NetPnL, cfg.DemoteMaxNet)
			result.Demoted++
		case row.NetPnL >= cfg.PromoteMinNet && row.ROI >= cfg.PromoteMinROI:
			row.Decision = "promote"
			row.Reason = fmt.Sprintf("net %+.2fU, ROI %+.2f%%", row.NetPnL, row.ROI)
			result.Promoted++
		default:
			row.Decision = "keep"
			row.Reason = fmt.Sprintf("net %+.2fU, ROI %+.2f%%", row.NetPnL, row.ROI)
		}
		result.Wallets = append(result.Wallets, *row)
	}
	sort.Slice(result.Wallets, func(i, j int) bool {
		if result.Wallets[i].NetPnL != result.Wallets[j].NetPnL {
			return result.Wallets[i].NetPnL > result.Wallets[j].NetPnL
		}
		return result.Wallets[i].Source < result.Wallets[j].Source
	})
	return result
}

func walletFromSource(source string, aliases map[string]string) string {
	source = strings.ToLower(strings.TrimSpace(source))
	for _, prefix := range []string{"copytrade_wallet:", "copytrade_football_score_wallet:", "copytrade_collect_wallet:", "copytrade_collect_football_score_wallet:"} {
		if strings.HasPrefix(source, prefix) {
			wallet := strings.TrimPrefix(source, prefix)
			if len(wallet) == 42 && strings.HasPrefix(wallet, "0x") {
				return wallet
			}
		}
	}
	for _, key := range []string{source, strings.TrimPrefix(source, "copytrade_football_score_"), strings.TrimPrefix(source, "copytrade_")} {
		if wallet := strings.ToLower(strings.TrimSpace(aliases[key])); wallet != "" {
			return wallet
		}
	}
	return ""
}

func tradeNet(tr journal.TradeRecord) float64 {
	if tr.NetPnLUSD == 0 && tr.EntryFeeUSD == 0 && tr.ExitFeeUSD == 0 {
		return tr.PnLUSD
	}
	return tr.NetPnLUSD
}

func basePositionID(id string) string {
	if dot := strings.IndexByte(id, '.'); dot > 0 {
		return id[:dot]
	}
	return id
}

func positionKey(tr journal.TradeRecord) string {
	id := basePositionID(tr.ID)
	if tr.EntryTime.IsZero() {
		return id
	}
	// Position IDs can restart from p1 after a state reset. Entry time is
	// stable across tranches and separates otherwise unrelated historical IDs.
	return id + "|" + tr.EntryTime.UTC().Format(time.RFC3339Nano)
}

func normalizedSource(source string) string {
	source = strings.TrimSpace(source)
	if source == "" {
		return "unknown"
	}
	return source
}

func strategyLabel(tr journal.TradeRecord) string {
	return strategyLabelFrom(tr.SignalSource, tr.Question)
}

func strategyLabelFrom(signalSource, questionText string) string {
	source := strings.ToLower(signalSource)
	question := strings.ToLower(questionText)
	switch {
	case strings.HasPrefix(source, "copytrade_collect_football_score"):
		return "football_score_collect"
	case strings.HasPrefix(source, "copytrade_collect"):
		return "copytrade_collect"
	case strings.HasPrefix(source, "copytrade_football_score"), strings.Contains(question, "exact score"):
		return "football_score"
	case strings.HasPrefix(source, "copytrade"):
		return "copytrade"
	case source == "manual":
		return "manual"
	case source == "auto" || source == "":
		return "legacy_auto"
	default:
		return "other"
	}
}

func isCollectionStrategy(name string) bool {
	return strings.HasSuffix(strings.ToLower(strings.TrimSpace(name)), "_collect")
}

func isCollectionPosition(pos strategy.Position) bool {
	source := pos.SignalSource
	if strings.TrimSpace(source) == "" {
		source = pos.Source
	}
	return isCollectionStrategy(strategyLabelFrom(source, pos.Question))
}

func mergeLabel(current, next string) string {
	if next == "" {
		return current
	}
	if current == "" || current == next {
		return next
	}
	return "mixed"
}

func mergeSource(current, next string) string {
	if current == "" || current == "unknown" || current == "auto" {
		return next
	}
	if next == "" || next == "unknown" || next == "auto" || current == next {
		return current
	}
	return "mixed"
}

func mergeStrategy(current, next string) string {
	if current == "" || current == "legacy_auto" {
		return next
	}
	if next == "" || next == "legacy_auto" || current == next {
		return current
	}
	return "mixed"
}

func stakeBucket(capital float64) string {
	switch {
	case capital >= 19:
		return "20U+"
	case capital >= 9:
		return "10U"
	case capital >= 4.5:
		return "5U"
	default:
		return "<5U"
	}
}

func acc(groups map[string]*Cohort, name string, pos *positionResult) {
	group := groups[name]
	if group == nil {
		group = &Cohort{Name: name}
		groups[name] = group
	}
	accCohort(group, pos)
}

func accCohort(group *Cohort, pos *positionResult) {
	group.Positions++
	group.Records += pos.records
	group.CapitalUSD += pos.capital
	group.GrossPnL += pos.gross
	group.FeesUSD += pos.entryFees + pos.exitFees
	group.NetPnL += pos.net
	switch {
	case pos.net > 0:
		group.Wins++
	case pos.net < 0:
		group.Losses++
	default:
		group.Flat++
	}
}

func finalizeScope(scope *PnLScope) {
	scope.ConservativeOpenNetPnL = -(scope.OpenExposureUSD + scope.OpenEntryFeesUSD)
	scope.ConservativeTotalNetPnL = scope.Closed.NetPnL + scope.ConservativeOpenNetPnL
}

func cohorts(groups map[string]*Cohort, byNet bool) []Cohort {
	out := make([]Cohort, 0, len(groups))
	for _, group := range groups {
		out = append(out, *group)
	}
	if byNet {
		sort.Slice(out, func(i, j int) bool {
			if out[i].NetPnL != out[j].NetPnL {
				return out[i].NetPnL > out[j].NetPnL
			}
			return out[i].Name < out[j].Name
		})
	} else {
		sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	}
	return out
}

func FormatMarkdown(report Report) string {
	var b strings.Builder
	b.WriteString("# Smartmoney Paper PnL\n\n")
	fmt.Fprintf(&b, "Generated: %s\n\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- Tradable realized net PnL: **%+.2fU** (%d positions)\n", report.Tradable.Closed.NetPnL, report.Tradable.Closed.Positions)
	fmt.Fprintf(&b, "- Broad collection realized net PnL: %+.2fU (%d positions; research only)\n", report.BroadCollection.Closed.NetPnL, report.BroadCollection.Closed.Positions)
	fmt.Fprintf(&b, "- Combined observed realized net PnL: %+.2fU\n", report.RealizedNetPnLUSD)
	fmt.Fprintf(&b, "- Gross PnL: %+.2fU\n", report.GrossPnLUSD)
	fmt.Fprintf(&b, "- Fees: %.2fU (entry %.2fU / exit %.2fU)\n", report.FeesUSD, report.EntryFeesUSD, report.ExitFeesUSD)
	fmt.Fprintf(&b, "- Closed positions: %d (%d journal rows)\n", report.ClosedPositions, report.Records)
	fmt.Fprintf(&b, "- Open: %d positions / %.2fU exposure / %.2fU remaining entry fees\n", report.OpenPositions, report.OpenExposureUSD, report.OpenEntryFeesUSD)
	fmt.Fprintf(&b, "- Conservative open PnL: %+.2fU (open value marked to zero)\n", report.ConservativeOpenNetPnLUSD)
	fmt.Fprintf(&b, "- Conservative total PnL: **%+.2fU**\n", report.ConservativeTotalNetPnLUSD)
	fmt.Fprintf(&b, "- Tradable conservative total: **%+.2fU**; broad collection conservative total: %+.2fU\n", report.Tradable.ConservativeTotalNetPnL, report.BroadCollection.ConservativeTotalNetPnL)

	writeCohorts(&b, "By Policy", report.ByPolicy)
	writeCohorts(&b, "By Stake", report.ByStake)
	writeCohorts(&b, "By Strategy", report.ByStrategy)
	writeCohorts(&b, "By Cost Model", report.ByCostModel)
	writeCohorts(&b, "By Signal Source", report.BySignalSource)
	return b.String()
}

func writeCohorts(b *strings.Builder, title string, rows []Cohort) {
	fmt.Fprintf(b, "\n## %s\n\n", title)
	b.WriteString("| Cohort | Positions | Rows | Capital | Gross | Fees | Net | W/L/F |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, row := range rows {
		fmt.Fprintf(b, "| %s | %d | %d | %.2f | %+.2f | %.2f | %+.2f | %d/%d/%d |\n",
			row.Name, row.Positions, row.Records, row.CapitalUSD, row.GrossPnL, row.FeesUSD, row.NetPnL, row.Wins, row.Losses, row.Flat)
	}
}
