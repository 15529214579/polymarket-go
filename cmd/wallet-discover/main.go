package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func main() {
	def := walletdiscover.DefaultConfig()
	var cfg walletdiscover.Config
	flag.StringVar(&cfg.GammaBase, "gamma_base", def.GammaBase, "Gamma API base URL")
	flag.StringVar(&cfg.DataBase, "data_base", def.DataBase, "Polymarket data API base URL")
	flag.StringVar(&cfg.LeaderboardBase, "leaderboard_base", def.LeaderboardBase, "Polymarket leaderboard API base URL")
	flag.IntVar(&cfg.MarketsLimit, "markets", def.MarketsLimit, "max active markets to scan")
	flag.IntVar(&cfg.TradesPages, "trade_pages", def.TradesPages, "recent /trades pages to scan")
	flag.IntVar(&cfg.TradesLimit, "trade_limit", def.TradesLimit, "recent /trades page size")
	flag.IntVar(&cfg.HoldersLimit, "holders", def.HoldersLimit, "holders per market")
	flag.IntVar(&cfg.ActivityPages, "activity_pages", def.ActivityPages, "activity pages per candidate")
	flag.IntVar(&cfg.ActivityLimit, "activity_limit", def.ActivityLimit, "activity page size per candidate")
	flag.IntVar(&cfg.ClosedLimit, "closed_limit", def.ClosedLimit, "closed positions per candidate")
	flag.IntVar(&cfg.Concurrency, "concurrency", def.Concurrency, "wallet scoring concurrency")
	flag.IntVar(&cfg.MaxCandidates, "max_candidates", def.MaxCandidates, "max candidates to score")
	flag.IntVar(&cfg.Days, "days", def.Days, "lookback days for wallet activity")
	flag.Float64Var(&cfg.MinNotionalUSD, "min_notional", def.MinNotionalUSD, "minimum trade notional USD for discovery")
	flag.Float64Var(&cfg.MinHolderShares, "min_holder_shares", def.MinHolderShares, "minimum holder shares for discovery")
	flag.Float64Var(&cfg.MinPrice, "min_price", def.MinPrice, "minimum trade price for discovery")
	flag.Float64Var(&cfg.MaxPrice, "max_price", def.MaxPrice, "maximum trade price for discovery")
	flag.Float64Var(&cfg.CopyStakeUSD, "copy_stake", def.CopyStakeUSD, "fixed USD stake used for copy-simulation")
	flag.Float64Var(&cfg.CopySlippageBP, "copy_slippage_bp", def.CopySlippageBP, "round-trip copy-simulation slippage in basis points")
	flag.Float64Var(&cfg.CopyFeeBP, "copy_fee_bp", def.CopyFeeBP, "copy-simulation fee in basis points per entry/exit")
	flag.IntVar(&cfg.LeaderboardLimit, "leaderboard_limit", def.LeaderboardLimit, "top N leaderboard rows to crawl per kind/window; 0 disables")
	flag.StringVar(&cfg.LeaderboardWindows, "leaderboard_windows", def.LeaderboardWindows, "comma-separated leaderboard windows")
	flag.StringVar(&cfg.LeaderboardKinds, "leaderboard_kinds", def.LeaderboardKinds, "comma-separated leaderboard kinds, e.g. profit,volume")
	flag.StringVar(&cfg.TargetCategories, "target_categories", def.TargetCategories, "comma-separated discovery categories to scan, e.g. basketball,soccer,esports")
	flag.StringVar(&cfg.ExistingWallets, "existing_wallets", def.ExistingWallets, "existing wallets file to re-score")
	flag.StringVar(&cfg.SportsTapeWallets, "sports_tape_wallets", def.SportsTapeWallets, "recent target-category large-order wallets to score as sports_tape candidates")
	flag.StringVar(&cfg.RetainWallets, "retain_wallets", def.RetainWallets, "wallets to keep in the discovery candidate pool as retained strategy candidates")
	flag.StringVar(&cfg.OutputDir, "output_dir", def.OutputDir, "output directory for JSON/JSONL")
	flag.StringVar(&cfg.ReportPath, "report", def.ReportPath, "markdown report path")
	flag.StringVar(&cfg.GeneratedTier, "generated_tier", def.GeneratedTier, "minimum tier written to wallets.generated.txt")
	flag.StringVar(&cfg.GeneratedWalletsPath, "generated_wallets", def.GeneratedWalletsPath, "path for generated wallet list")
	flag.StringVar(&cfg.AutoWalletsPath, "auto_wallets", def.AutoWalletsPath, "path for auto-small smart-money wallet list")
	flag.StringVar(&cfg.PromptWalletsPath, "prompt_wallets", def.PromptWalletsPath, "path for prompt smart-money wallet list")
	flag.StringVar(&cfg.PositiveWalletsPath, "positive_wallets", def.PositiveWalletsPath, "path for positive copy-simulation production wallet list")
	timeout := flag.Duration("timeout", 20*time.Minute, "overall run timeout")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, err := walletdiscover.Run(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet-discover: %v\n", err)
		os.Exit(1)
	}
	counts := map[string]int{}
	for _, s := range result.Scores {
		counts[s.Tier]++
	}
	fmt.Printf("wallet-discover done: markets=%d candidates=%d A=%d B=%d C=%d D=%d BOT=%d\n",
		len(result.Markets), len(result.Scores),
		counts["A"], counts["B"], counts["C"], counts["D"], counts["BOT"])
	fmt.Printf("outputs: %s/wallet_candidates.jsonl %s/wallet_scores.json %s %s %s %s %s\n",
		cfg.OutputDir, cfg.OutputDir, cfg.ReportPath, cfg.GeneratedWalletsPath, cfg.AutoWalletsPath, cfg.PromptWalletsPath, cfg.PositiveWalletsPath)
}
