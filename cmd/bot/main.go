package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func main() {
	mode := flag.String("mode", "run", "run | discover | feed | sample | detect")
	maxMarkets := flag.Int("markets", 20, "top-N LoL markets by vol24h to subscribe")
	windowSec := flag.Int("window", 60, "sampler window in seconds")
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
	case "feed":
		if err := runFeed(ctx, *maxMarkets); err != nil {
			slog.Error("feed failed", "err", err)
			os.Exit(1)
		}
	case "sample":
		if err := runSample(ctx, *maxMarkets, *windowSec); err != nil && ctx.Err() == nil {
			slog.Error("sample failed", "err", err)
			os.Exit(1)
		}
	case "detect":
		if err := runDetect(ctx, *maxMarkets, *windowSec); err != nil && ctx.Err() == nil {
			slog.Error("detect failed", "err", err)
			os.Exit(1)
		}
	case "run":
		slog.Info("polymarket-go starting", "mode", "paper")
		// Phase 2+: strategy loop. For now, bot -mode=feed exercises the data layer.
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

func runFeed(ctx context.Context, topN int) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	lol := feed.FilterLoL(all)
	if len(lol) == 0 {
		return fmt.Errorf("no active LoL markets")
	}
	if topN > len(lol) {
		topN = len(lol)
	}
	lol = lol[:topN]

	// Collect clobTokenIDs; map back to question for log context.
	assetToQ := map[string]string{}
	var assetIDs []string
	for _, m := range lol {
		for _, id := range m.ClobTokenIDs() {
			if id == "" {
				continue
			}
			assetToQ[id] = m.Question
			assetIDs = append(assetIDs, id)
		}
	}
	slog.Info("feed.start", "markets", len(lol), "assets", len(assetIDs))

	ws := feed.NewWSSClient(assetIDs)

	// consumer: log tick summaries
	go func() {
		throttle := map[string]time.Time{}
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-ws.Books():
				if !ok {
					return
				}
				// throttle per-asset book logs to 1/s to keep log volume sane
				if t, seen := throttle[ev.AssetID]; seen && time.Since(t) < time.Second {
					continue
				}
				throttle[ev.AssetID] = time.Now()
				bestBid, bestAsk := 0.0, 0.0
				if len(ev.Bids) > 0 {
					bestBid = ev.Bids[0].Price
				}
				if len(ev.Asks) > 0 {
					bestAsk = ev.Asks[0].Price
				}
				slog.Info("book",
					"asset", short(ev.AssetID),
					"q", assetToQ[ev.AssetID],
					"bid", bestBid,
					"ask", bestAsk,
					"n_bids", len(ev.Bids),
					"n_asks", len(ev.Asks),
				)
			case tr, ok := <-ws.Trades():
				if !ok {
					return
				}
				slog.Info("trade",
					"asset", short(tr.AssetID),
					"q", assetToQ[tr.AssetID],
					"price", tr.Price,
					"size", tr.Size,
					"side", tr.Side,
				)
			}
		}
	}()

	return ws.Run(ctx)
}

func runSample(ctx context.Context, topN, windowSec int) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	lol := feed.FilterLoL(all)
	if len(lol) == 0 {
		return fmt.Errorf("no active LoL markets")
	}
	if topN > len(lol) {
		topN = len(lol)
	}
	lol = lol[:topN]

	assetToQ := map[string]string{}
	var assetIDs []string
	for _, m := range lol {
		for _, id := range m.ClobTokenIDs() {
			if id == "" {
				continue
			}
			assetToQ[id] = m.Question
			assetIDs = append(assetIDs, id)
		}
	}
	slog.Info("sample.start", "markets", len(lol), "assets", len(assetIDs), "window_sec", windowSec)

	ws := feed.NewWSSClient(assetIDs)
	sampler := feed.NewSampler(windowSec)

	go func() {
		if err := sampler.Run(ctx, ws.Books(), ws.Trades()); err != nil && ctx.Err() == nil {
			slog.Error("sampler exited", "err", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-sampler.Ticks():
				if !ok {
					return
				}
				slog.Info("tick",
					"asset", short(t.AssetID),
					"q", assetToQ[t.AssetID],
					"bid", t.BestBid,
					"ask", t.BestAsk,
					"mid", t.Mid,
					"trades", t.Trades,
					"buy_vol", t.BuyVol,
					"sell_vol", t.SellVol,
				)
			}
		}
	}()

	// periodic window summary, every 10s
	go func() {
		tk := time.NewTicker(10 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				for _, w := range sampler.Snapshot() {
					slog.Info("window",
						"asset", short(w.AssetID),
						"q", assetToQ[w.AssetID],
						"samples", w.Samples,
						"start_mid", w.StartMid,
						"end_mid", w.EndMid,
						"delta_pp", w.DeltaPP,
						"up", w.Upticks,
						"down", w.Downticks,
						"flat", w.Flats,
						"buy_ratio", w.BuyRatio,
					)
				}
			}
		}
	}()

	return ws.Run(ctx)
}

func runDetect(ctx context.Context, topN, windowSec int) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	lol := feed.FilterLoL(all)
	if len(lol) == 0 {
		return fmt.Errorf("no active LoL markets")
	}
	if topN > len(lol) {
		topN = len(lol)
	}
	lol = lol[:topN]

	assetToQ := map[string]string{}
	var assetIDs []string
	for _, m := range lol {
		for _, id := range m.ClobTokenIDs() {
			if id == "" {
				continue
			}
			assetToQ[id] = m.Question
			assetIDs = append(assetIDs, id)
		}
	}
	slog.Info("detect.start", "markets", len(lol), "assets", len(assetIDs), "window_sec", windowSec)

	ws := feed.NewWSSClient(assetIDs)
	sampler := feed.NewSampler(windowSec)

	cfg := strategy.DefaultConfig()
	cfg.WindowSec = windowSec
	if windowSec < cfg.MinSamplesWarm {
		cfg.MinSamplesWarm = windowSec / 2
	}
	det := strategy.NewDetector(cfg, sampler)

	go func() {
		if err := sampler.Run(ctx, ws.Books(), ws.Trades()); err != nil && ctx.Err() == nil {
			slog.Error("sampler exited", "err", err)
		}
	}()
	go func() {
		if err := det.Run(ctx); err != nil && ctx.Err() == nil {
			slog.Error("detector exited", "err", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case sig, ok := <-det.Signals():
				if !ok {
					return
				}
				slog.Info("signal",
					"asset", short(sig.AssetID),
					"q", assetToQ[sig.AssetID],
					"mid", sig.Mid,
					"delta_pp", sig.DeltaPP,
					"tail_ups", sig.TailUps,
					"tail_len", sig.TailLen,
					"buy_ratio", sig.BuyRatio,
					"reason", sig.Reason,
				)
			}
		}
	}()

	go func() {
		tk := time.NewTicker(30 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				top := topWindow(sampler.Snapshot(), 5)
				for _, w := range top {
					slog.Info("top_window",
						"asset", short(w.AssetID),
						"q", assetToQ[w.AssetID],
						"samples", w.Samples,
						"mid", w.EndMid,
						"delta_pp", w.DeltaPP,
						"up", w.Upticks,
						"buy_ratio", w.BuyRatio,
					)
				}
			}
		}
	}()

	return ws.Run(ctx)
}

func topWindow(ws []feed.WindowStats, n int) []feed.WindowStats {
	// sort by DeltaPP desc; naive O(n^2) good enough for small sets
	out := append([]feed.WindowStats(nil), ws...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].DeltaPP > out[i].DeltaPP {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if n > len(out) {
		n = len(out)
	}
	return out[:n]
}

func short(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:6] + ".." + id[len(id)-4:]
}
