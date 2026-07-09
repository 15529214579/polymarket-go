package walletdiscover

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadWallets(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []string
	seen := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(strings.Split(sc.Text(), "#")[0])
		addr := normalizeAddress(line)
		if addr == "" {
			continue
		}
		if _, ok := seen[addr]; ok {
			continue
		}
		seen[addr] = struct{}{}
		out = append(out, addr)
	}
	return out, sc.Err()
}

func LoadCachedActivity(dir, addr string) ([]Trade, error) {
	path := filepath.Join(dir, "wallet_activity", addr+".jsonl")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []Trade
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var tr Trade
		if err := json.Unmarshal(sc.Bytes(), &tr); err == nil {
			out = append(out, tr)
		}
	}
	return out, sc.Err()
}

func SaveActivity(dir, addr string, trades []Trade) error {
	if err := os.MkdirAll(filepath.Join(dir, "wallet_activity"), 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "wallet_activity", addr+".jsonl")
	seen := map[string]struct{}{}
	sort.Slice(trades, func(i, j int) bool { return trades[i].Timestamp < trades[j].Timestamp })
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, tr := range trades {
		key := tr.TransactionHash + "|" + tr.Asset + "|" + fmt.Sprint(tr.Timestamp)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		if err := enc.Encode(tr); err != nil {
			return err
		}
	}
	return nil
}

func SaveResult(result *Result, cfg Config) error {
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return err
	}
	if err := writeCandidates(filepath.Join(cfg.OutputDir, "wallet_candidates.jsonl"), result.Candidates); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(cfg.OutputDir, "wallet_scores.json"), result.Scores); err != nil {
		return err
	}
	if err := writeTierJSON(filepath.Join(cfg.OutputDir, "copytrade_backtest_results.generated.json"), result.Scores); err != nil {
		return err
	}
	if err := writeGeneratedWallets(cfg.GeneratedWalletsPath, result.Scores, cfg.GeneratedTier); err != nil {
		return err
	}
	if err := writeFollowActionWallets(cfg.AutoWalletsPath, result.Scores, "auto-small"); err != nil {
		return err
	}
	if err := writeFollowActionWallets(cfg.PromptWalletsPath, result.Scores, "prompt"); err != nil {
		return err
	}
	if err := writePositiveCopyWallets(cfg.PositiveWalletsPath, result.Scores); err != nil {
		return err
	}
	if cfg.ReportPath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.ReportPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(cfg.ReportPath, []byte(RenderReport(result, cfg)), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeCandidates(path string, candidates []*Candidate) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, c := range candidates {
		if err := enc.Encode(c); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func writeTierJSON(path string, scores []WalletScore) error {
	out := map[string]map[string]any{}
	for _, s := range scores {
		out[s.Address] = map[string]any{
			"tier":              s.Tier,
			"reason":            s.Reason,
			"follow_action":     s.FollowAction,
			"bot_score":         s.BotScore,
			"edge":              s.Edge,
			"smart_money_score": s.SmartMoneyScore,
			"risk_flags":        s.RiskFlags,
			"strengths":         s.Strengths,
			"stats":             s.Stats,
		}
	}
	return writeJSON(path, out)
}

func writeGeneratedWallets(path string, scores []WalletScore, minTier string) error {
	var lines []string
	for _, s := range scores {
		if TierAllowed(s.Tier, minTier) {
			lines = append(lines, fmt.Sprintf("%s # tier=%s action=%s smart=%.2f score=%.2f bot=%.2f edge=%.2f %s",
				s.Address, s.Tier, s.FollowAction, s.SmartMoneyScore, s.Score, s.BotScore, s.Edge, s.Reason))
		}
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func writeFollowActionWallets(path string, scores []WalletScore, action string) error {
	var lines []string
	for _, s := range scores {
		if s.FollowAction != action {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s # tier=%s action=%s smart=%.2f bot=%.2f edge=%.2f %s",
			s.Address, s.Tier, s.FollowAction, s.SmartMoneyScore, s.BotScore, s.Edge, s.Reason))
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func writePositiveCopyWallets(path string, scores []WalletScore) error {
	var lines []string
	for _, s := range scores {
		if !TierAllowed(s.Tier, "B") {
			continue
		}
		if s.FollowAction != "auto-small" && s.FollowAction != "prompt" {
			continue
		}
		if s.BotScore >= 35 {
			continue
		}
		if s.Stats.CopyClosedTrades < 3 || s.Stats.CopyROI <= 0 || s.Stats.CopyPnL <= 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s # tier=%s action=%s smart=%.2f bot=%.2f edge=%.2f copyROI=%.1f%% copyPnL=$%+.2f copyT=%d %s",
			s.Address, s.Tier, s.FollowAction, s.SmartMoneyScore, s.BotScore, s.Edge,
			s.Stats.CopyROI, s.Stats.CopyPnL, s.Stats.CopyClosedTrades, s.Reason))
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
