package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/paperreport"
)

func main() {
	journalDir := flag.String("journal", "db/smartmoney-paper/journal", "paper journal directory")
	whaleLog := flag.String("whale_log", "db/smartmoney-paper/journal/whale_trades.jsonl", "whale event log used to resolve legacy labels")
	jsonOut := flag.String("json_out", "db/smartmoney-paper/wallet-performance.json", "wallet policy JSON output")
	reportOut := flag.String("report", "reports/smartmoney-paper-wallets.md", "wallet policy Markdown output")
	promotedOut := flag.String("promoted", "db/smartmoney-paper/wallets.paper-promoted.txt", "promoted wallet overlay")
	demotedOut := flag.String("demoted", "db/smartmoney-paper/wallets.paper-demoted.txt", "demoted wallet block list")
	minPositions := flag.Int("min_positions", 10, "minimum independent wallet-market samples before a decision")
	promoteMinNet := flag.Float64("promote_min_net", 5, "minimum net PnL in U for promotion")
	promoteMinROI := flag.Float64("promote_min_roi", 2, "minimum net ROI percentage for promotion")
	promoteMinWinRate := flag.Float64("promote_min_win_rate", 45, "minimum independent-sample win rate percentage for promotion")
	promoteMinTrimmedNet := flag.Float64("promote_min_trimmed_net", 1, "minimum net PnL after removing the best sample")
	maxBestSampleShare := flag.Float64("max_best_sample_share", 60, "maximum percentage of total net PnL contributed by the best sample")
	maxTwoSidedMarkets := flag.Int("max_two_sided_markets", 0, "maximum markets where a wallet traded multiple outcomes")
	demoteMaxNet := flag.Float64("demote_max_net", -5, "maximum net PnL in U before demotion")
	flag.Parse()

	trades, err := journal.ReadAll(*journalDir)
	if err != nil {
		fatalf("read journal: %v", err)
	}
	aliases, err := readAliases(*whaleLog)
	if err != nil {
		fatalf("read whale aliases: %v", err)
	}
	report := paperreport.AnalyzeWalletPolicy(trades, aliases, paperreport.WalletPolicyConfig{
		MinPositions:         *minPositions,
		PromoteMinNet:        *promoteMinNet,
		PromoteMinROI:        *promoteMinROI,
		PromoteMinWinRate:    *promoteMinWinRate,
		PromoteMinTrimmedNet: *promoteMinTrimmedNet,
		MaxBestSampleShare:   *maxBestSampleShare,
		MaxTwoSidedMarkets:   *maxTwoSidedMarkets,
		DemoteMaxNet:         *demoteMaxNet,
	})

	body, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fatalf("marshal report: %v", err)
	}
	if err := atomicWrite(*jsonOut, append(body, '\n')); err != nil {
		fatalf("write JSON: %v", err)
	}
	markdown := formatMarkdown(report)
	if err := atomicWrite(*reportOut, []byte(markdown)); err != nil {
		fatalf("write report: %v", err)
	}
	if err := atomicWrite(*promotedOut, []byte(formatWalletList(report, "promote"))); err != nil {
		fatalf("write promoted wallets: %v", err)
	}
	if err := atomicWrite(*demotedOut, []byte(formatWalletList(report, "demote"))); err != nil {
		fatalf("write demoted wallets: %v", err)
	}
	fmt.Printf("smartmoney-wallet-policy promoted=%d demoted=%d unresolved=%d report=%s\n", report.Promoted, report.Demoted, report.Unresolved, *reportOut)
}

func readAliases(path string) (map[string]string, error) {
	aliases := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return aliases, nil
		}
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	buf := make([]byte, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		var event struct {
			Wallet string `json:"wallet"`
			Label  string `json:"label"`
		}
		if json.Unmarshal(sc.Bytes(), &event) != nil {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(event.Wallet))
		label := strings.ToLower(strings.TrimSpace(event.Label))
		if len(wallet) != 42 || !strings.HasPrefix(wallet, "0x") || label == "" {
			continue
		}
		aliases[label] = wallet
		aliases["copytrade_"+label] = wallet
		aliases["copytrade_football_score_"+label] = wallet
		aliases["copytrade_collect_"+label] = wallet
		aliases["copytrade_collect_football_score_"+label] = wallet
	}
	return aliases, sc.Err()
}

func formatMarkdown(report paperreport.WalletPolicyReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Smartmoney Paper Wallet Policy\n\nGenerated: %s\n\n", report.GeneratedAt.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Fprintf(&b, "- Rule: at least %d independent wallet-market samples; promote at net >= %+.2fU, ROI >= %+.2f%%, win rate >= %.1f%%, and trimmed net >= %+.2fU\n", report.Config.MinPositions, report.Config.PromoteMinNet, report.Config.PromoteMinROI, report.Config.PromoteMinWinRate, report.Config.PromoteMinTrimmedNet)
	fmt.Fprintf(&b, "- Robustness: best sample <= %.1f%% of net; two-sided markets <= %d; demote at net <= %+.2fU\n", report.Config.MaxBestSampleShare, report.Config.MaxTwoSidedMarkets, report.Config.DemoteMaxNet)
	fmt.Fprintf(&b, "- Decisions: %d promoted / %d demoted / %d unresolved\n\n", report.Promoted, report.Demoted, report.Unresolved)
	b.WriteString("| Decision | Wallet | Samples | Positions | Capital | Fees | Net | Trimmed | ROI | Win | Two-sided | Max/sample | Reason |\n")
	b.WriteString("|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")
	for _, row := range report.Wallets {
		wallet := row.Wallet
		if wallet == "" {
			wallet = row.Source
		}
		fmt.Fprintf(&b, "| %s | %s | %d | %d | %.2f | %.2f | %+.2f | %+.2f | %+.2f%% | %.1f%% | %d | %d | %s |\n",
			row.Decision, wallet, row.IndependentSamples, row.Positions, row.CapitalUSD, row.FeesUSD, row.NetPnL, row.TrimmedNetPnL, row.ROI, row.WinRate, row.TwoSidedMarkets, row.MaxPositionsPerSample, row.Reason)
	}
	return b.String()
}

func formatWalletList(report paperreport.WalletPolicyReport, decision string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# generated by smartmoney-wallet-policy; decision=%s\n", decision)
	for _, row := range report.Wallets {
		if row.Decision != decision || row.Wallet == "" {
			continue
		}
		list := "paper_" + decision + "d"
		tier := ""
		if decision == "promote" {
			tier = " tier=B"
		}
		fmt.Fprintf(&b, "%s # list=%s%s localPaperSamples=%d localPaperPositions=%d localPaperROI=%.2f%% localPaperPnL=%+.2fU\n",
			row.Wallet, list, tier, row.IndependentSamples, row.Positions, row.ROI, row.NetPnL)
	}
	return b.String()
}

func atomicWrite(path string, body []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".wallet-policy-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(body); err != nil {
		f.Close()
		return err
	}
	if err := f.Chmod(0o644); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
