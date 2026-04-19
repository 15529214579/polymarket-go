package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func main() {
	mode := flag.String("mode", "run", "run | discover")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch *mode {
	case "discover":
		if err := runDiscover(ctx); err != nil {
			slog.Error("discover failed", "err", err)
			os.Exit(1)
		}
	case "run":
		slog.Info("polymarket-go starting", "mode", "paper")
		// Phase 1.2: wire real WSSClient here once discover gives us assetIDs.
		<-ctx.Done()
		slog.Info("shutdown")
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n", *mode)
		os.Exit(2)
	}
}

func runDiscover(ctx context.Context) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	lol := feed.FilterLoL(all)
	slog.Info("gamma.discover", "total_active", len(all), "lol", len(lol))
	for _, m := range lol {
		tokens := m.ClobTokenIDs()
		slog.Info("lol_market",
			"q", m.Question,
			"slug", m.Slug,
			"vol24h", m.Volume24hr,
			"liq_clob", m.LiquidityClob,
			"accepting", m.AcceptingOrders,
			"end", m.EndDate,
			"tokens", len(tokens),
		)
	}
	return nil
}
