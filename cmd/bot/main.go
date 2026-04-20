package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/15529214579/polymarket-go/internal/config"
	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/notify"
	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/15529214579/polymarket-go/internal/risk"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func main() {
	mode := flag.String("mode", "run", "run | discover | feed | sample | detect")
	maxMarkets := flag.Int("markets", 20, "top-N LoL markets by vol24h to subscribe")
	windowSec := flag.Int("window", 60, "sampler window in seconds")
	slippageBp := flag.Float64("slippage_bp", 0, "paper fill slippage in bp applied against you")
	largeFillUSD := flag.Float64("large_fill_usd", 3.0, "DM notifier threshold on |realized pnl|")
	envFile := flag.String("env_file", ".env.local", "dotenv file to load before reading env")
	signalMode := flag.String("signal_mode", "auto", "auto (paper-submit on signal) | prompt (DM + Buy 1/5/10 inline buttons, boss picks size)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := config.LoadDotEnv(*envFile); err != nil {
		slog.Warn("dotenv_load_warn", "path", *envFile, "err", err.Error())
	}

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
		if err := runDetect(ctx, *maxMarkets, *windowSec, *slippageBp, *largeFillUSD, *signalMode); err != nil && ctx.Err() == nil {
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

func runDetect(ctx context.Context, topN, windowSec int, slippageBp, largeFillUSD float64, signalMode string) error {
	if signalMode != "auto" && signalMode != "prompt" {
		return fmt.Errorf("invalid signal_mode %q (want auto|prompt)", signalMode)
	}
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
	exitCfg := strategy.DefaultExitConfig()
	exit := strategy.NewExitTracker(exitCfg)
	posCfg := strategy.DefaultPositionConfig()
	pm := strategy.NewPositionManager(posCfg)
	paper := order.NewPaperClient(slippageBp)
	riskCfg := risk.DefaultConfig()
	rm := risk.New(riskCfg, time.Now())
	notifier := buildNotifier()
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = notifier.Close(sctx)
	}()
	pending := notify.NewPendingStore(60 * time.Second)
	slog.Info("paper_client.ready", "slippage_bp", slippageBp, "per_pos_usd", posCfg.PerPositionUSD)
	slog.Info("risk.ready",
		"bankroll_usd", riskCfg.StartingBankrollUSD,
		"daily_loss_cap_usd", rm.State().DayLossCapUSD,
		"max_single_loss_usd", riskCfg.MaxSingleLossUSD,
		"feed_silence_sec", riskCfg.FeedSilenceSec,
		"large_fill_usd", largeFillUSD,
	)
	slog.Info("signal_mode.ready", "mode", signalMode)

	// Inbound callback consumer (Phase 3.5.b). Only runs if a DEDICATED sidecar
	// bot token is configured — we never long-poll the alert bot's token because
	// OpenClaw may also be polling it, and Telegram delivers updates competitively.
	sidecarToken := os.Getenv("SIDECAR_BOT_TOKEN")
	sidecarChat := os.Getenv("SIDECAR_CHAT_ID")
	if sidecarChat == "" {
		sidecarChat = os.Getenv("TELEGRAM_CHAT_ID")
	}
	if sidecarToken != "" && sidecarChat != "" {
		chatID, err := strconv.ParseInt(sidecarChat, 10, 64)
		if err != nil {
			slog.Warn("sidecar_chat_id_parse_fail", "err", err.Error())
		} else {
			h := &buyHandler{
				pm:       pm,
				exit:     exit,
				paper:    paper,
				rm:       rm,
				pending:  pending,
				notifier: notifier,
				assetToQ: assetToQ,
				largeFillUSD: largeFillUSD,
			}
			lp := notify.NewLongPoll(notify.LongPollConfig{
				BotToken:       sidecarToken,
				ExpectedChatID: chatID,
			}, h)
			go func() {
				slog.Info("sidecar_longpoll.ready", "chat_id", chatID)
				if err := lp.Run(ctx); err != nil && ctx.Err() == nil {
					slog.Warn("sidecar_longpoll_exit", "err", err.Error())
				}
			}()
		}
	} else if signalMode == "prompt" {
		slog.Warn("signal_mode_prompt_without_sidecar",
			"hint", "prompt mode needs SIDECAR_BOT_TOKEN + chat_id — buttons will arrive but clicks won't be consumed")
	}

	// Pending-store reaper so expired button prompts don't accumulate.
	go func() {
		tk := time.NewTicker(15 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				if n := pending.Reap(now); n > 0 {
					slog.Info("pending_reap", "expired", n, "remaining", pending.Size())
				}
			}
		}
	}()

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

	// Fan-out ticks to the exit tracker (only tracks opened positions).
	// Uses a fresh Sampler subscription via a side goroutine: we tap the detector's
	// upstream by reading ticks *through* the sampler's Ticks() channel which the
	// detector already consumes. To avoid a fight for one channel, run a dedicated
	// tap via TickTail polling on each detected open asset instead.
	// Simpler: have detect subscribe to ticks directly alongside the detector.
	//
	// Here we piggyback on the fact that only the detector consumes sampler.Ticks().
	// Instead of stealing them, expose a separate TickTail-based poller below.
	// Update: cleanest is to have the sampler fan-out — but we only have one
	// consumer right now. Workaround: poll sampler.Window for open positions.
	go func() {
		tk := time.NewTicker(1 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				for _, w := range sampler.Snapshot() {
					if !exit.Has(w.AssetID) {
						continue
					}
					tail, ok := sampler.TickTail(w.AssetID, 1)
					if !ok || len(tail) == 0 {
						continue
					}
					if sig, fired := exit.OnTick(tail[0]); fired {
						sellIntent := order.Intent{
							AssetID: sig.AssetID,
							Market:  sig.Market,
							Side:    order.Sell,
							SizeUSD: posCfg.PerPositionUSD,
							LimitPx: sig.ExitMid,
							Type:    order.GTC,
						}
						res, err := paper.Submit(ctx, sellIntent)
						if err != nil {
							slog.Warn("paper_sell_reject",
								"asset", short(sig.AssetID),
								"limit", sig.ExitMid,
								"err", err.Error())
							continue
						}
						// Override exit mid with the actual fill price so realized PnL
						// reflects paper slippage.
						sig.ExitMid = res.AvgPrice
						sig.ChangePP = (res.AvgPrice - sig.EntryMid) * 100
						closed, err := pm.Close(sig.AssetID, sig)
						if err != nil {
							slog.Warn("paper_close_miss", "asset", short(sig.AssetID), "err", err.Error())
							continue
						}
						stats := pm.Stats()
						if tripped := rm.OnClose(closed.PnLUSD, sig.Time); tripped {
							rst := rm.State()
							slog.Error("risk_trip",
								"reason", string(rst.BlockReason),
								"day_pnl_usd", rst.DayRealizedPnL,
								"cap_usd", rst.DayLossCapUSD,
							)
							notifier.RiskTrip(notify.RiskTripEvent{
								Reason:        string(rst.BlockReason),
								DayPnLUSD:     rst.DayRealizedPnL,
								DayLossCapUSD: rst.DayLossCapUSD,
								OpenPositions: stats.Open,
							})
						}
						if closed.PnLUSD <= -largeFillUSD || closed.PnLUSD >= largeFillUSD {
							notifier.LargeFill(notify.LargeFillEvent{
								Question: assetToQ[sig.AssetID],
								AssetID:  sig.AssetID,
								Side:     "sell",
								SizeUSD:  posCfg.PerPositionUSD,
								PnLUSD:   closed.PnLUSD,
								EntryPx:  sig.EntryMid,
								ExitPx:   res.AvgPrice,
								Reason:   string(sig.Reason),
								HeldSec:  int(sig.HeldFor.Seconds()),
							})
						}
						slog.Info("exit",
							"asset", short(sig.AssetID),
							"q", assetToQ[sig.AssetID],
							"reason", string(sig.Reason),
							"order_id", res.OrderID,
							"entry", sig.EntryMid,
							"peak", sig.PeakMid,
							"exit_fill", res.AvgPrice,
							"delta_pp", sig.ChangePP,
							"drawdown_pp", sig.DrawdownPP,
							"held_sec", int(sig.HeldFor.Seconds()),
							"pnl_usd", closed.PnLUSD,
							"open_positions", stats.Open,
							"realized_pnl", stats.RealizedPnLUSD,
						)
					}
				}
			}
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
				// Risk gate first — daily-loss breaker / feed-silence / manual pause.
				if err := rm.AllowOpen(time.Now()); err != nil {
					st := rm.State()
					slog.Warn("risk_block_open",
						"asset", short(sig.AssetID),
						"q", assetToQ[sig.AssetID],
						"reason", string(st.BlockReason),
						"day_pnl_usd", st.DayRealizedPnL,
						"cap_usd", st.DayLossCapUSD,
					)
					continue
				}
				// Dedupe checks before we bother submitting to the paper client;
				// avoids polluting history with rejected paper orders.
				if pm.Has(sig.AssetID) || pm.HasMarket(sig.Market) {
					slog.Info("paper_open_skip",
						"asset", short(sig.AssetID),
						"q", assetToQ[sig.AssetID],
						"reason", "already_open",
					)
					continue
				}

				// Prompt mode: publish the signal as a DM with Buy 1/5/10 buttons
				// and stash the intent in the pending store. The callback longpoll
				// (above) claims the nonce and executes via buyHandler.
				if signalMode == "prompt" {
					p := pending.Put(notify.PendingIntent{
						AssetID:  sig.AssetID,
						Market:   sig.Market,
						Question: assetToQ[sig.AssetID],
						Mid:      sig.Mid,
					}, time.Now())
					notifier.SignalPrompt(notify.SignalPromptEvent{
						Nonce:     p.Nonce,
						Question:  assetToQ[sig.AssetID],
						AssetID:   sig.AssetID,
						Mid:       sig.Mid,
						DeltaPP:   sig.DeltaPP,
						TailUps:   sig.TailUps,
						TailLen:   sig.TailLen,
						BuyRatio:  sig.BuyRatio,
						ExpiresIn: 60 * time.Second,
					})
					slog.Info("signal_prompt_sent",
						"asset", short(sig.AssetID),
						"nonce", p.Nonce,
						"mid", sig.Mid,
					)
					continue
				}

				buyIntent := order.Intent{
					AssetID: sig.AssetID,
					Market:  sig.Market,
					Side:    order.Buy,
					SizeUSD: posCfg.PerPositionUSD,
					LimitPx: sig.Mid,
					Type:    order.GTC,
				}
				res, err := paper.Submit(ctx, buyIntent)
				if err != nil {
					slog.Warn("paper_buy_reject",
						"asset", short(sig.AssetID),
						"limit", sig.Mid,
						"err", err.Error())
					continue
				}
				// Book the position at the actual fill price so slippage is priced in.
				entryTick := feed.Tick{
					AssetID: sig.AssetID, Market: sig.Market,
					Time: sig.Time, Mid: res.AvgPrice,
				}
				pos, err := pm.Open(sig.AssetID, sig.Market, entryTick)
				if err != nil {
					slog.Info("paper_open_skip",
						"asset", short(sig.AssetID),
						"q", assetToQ[sig.AssetID],
						"order_id", res.OrderID,
						"reason", err.Error(),
					)
					continue
				}
				exit.Open(sig.AssetID, sig.Market, entryTick)
				stats := pm.Stats()
				slog.Info("paper_open",
					"id", pos.ID,
					"order_id", res.OrderID,
					"asset", short(sig.AssetID),
					"q", assetToQ[sig.AssetID],
					"signal_mid", sig.Mid,
					"entry_fill", res.AvgPrice,
					"size_usd", pos.SizeUSD,
					"units", pos.Units,
					"open_positions", stats.Open,
					"total_exposure_usd", stats.TotalExposure,
				)
			}
		}
	}()

	// Feed-silence watchdog + periodic risk snapshot. SPEC §6: >30s WSS
	// silence trips breaker. We also push a risk summary every 60s so the
	// heartbeat log has a recent snapshot.
	go func() {
		tk := time.NewTicker(5 * time.Second)
		defer tk.Stop()
		lastSummary := time.Now()
		prevBlocked := false
		prevReason := ""
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				if at := ws.LastEventAt(); !at.IsZero() {
					rm.OnFeedHeartbeat(at)
				}
				silent, tripped := rm.CheckFeed(now)
				st := rm.State()
				if tripped {
					slog.Error("risk_trip",
						"reason", string(risk.BlockFeedSilence),
						"silent_sec", int(silent.Seconds()),
					)
					notifier.RiskTrip(notify.RiskTripEvent{
						Reason:        string(risk.BlockFeedSilence),
						DayPnLUSD:     st.DayRealizedPnL,
						DayLossCapUSD: st.DayLossCapUSD,
						SilentSec:     int(silent.Seconds()),
						OpenPositions: pm.Stats().Open,
					})
				}
				// Detect resume transition (boss manually resumed, or day rolled
				// over while the breaker was a pure daily-loss one — in practice
				// SPEC says we don't auto-clear, but keep this wired for when
				// we add a /resume command).
				if prevBlocked && !st.Blocked {
					notifier.RiskResume(notify.RiskResumeEvent{
						PrevReason:    prevReason,
						DayPnLUSD:     st.DayRealizedPnL,
						DayLossCapUSD: st.DayLossCapUSD,
					})
				}
				prevBlocked = st.Blocked
				if st.Blocked {
					prevReason = string(st.BlockReason)
				}
				if now.Sub(lastSummary) >= 60*time.Second {
					lastSummary = now
					st := rm.State()
					slog.Info("risk_status",
						"day", st.Day,
						"day_pnl_usd", st.DayRealizedPnL,
						"cap_usd", st.DayLossCapUSD,
						"blocked", st.Blocked,
						"reason", string(st.BlockReason),
						"single_loss_flags", st.SingleLossFlags,
						"feed_silent_sec", int(silent.Seconds()),
					)
				}
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

// buyHandler wires a click on Buy 1/5/10 → same paper-submit → pm.Open path
// the auto-mode signal loop uses, but honors the size the boss picked and the
// frozen mid captured at signal time. Executes synchronously on the longpoll
// goroutine; Telegram dispatch of the resulting DM is async via notifier.
type buyHandler struct {
	pm           *strategy.PositionManager
	exit         *strategy.ExitTracker
	paper        order.Client
	rm           *risk.Manager
	pending      *notify.PendingStore
	notifier     notify.Notifier
	assetToQ     map[string]string
	largeFillUSD float64
}

func (h *buyHandler) OnBuy(ctx context.Context, nonce string, sizeUSD float64) (string, error) {
	now := time.Now()
	p, ok := h.pending.Claim(nonce, now)
	if !ok {
		return "", fmt.Errorf("已过期或已点过")
	}
	if err := h.rm.AllowOpen(now); err != nil {
		st := h.rm.State()
		return "", fmt.Errorf("风控阻止: %s (day_pnl=%.2f)", st.BlockReason, st.DayRealizedPnL)
	}
	if h.pm.Has(p.AssetID) || h.pm.HasMarket(p.Market) {
		return "", fmt.Errorf("已有同市场仓位")
	}
	intent := order.Intent{
		AssetID: p.AssetID,
		Market:  p.Market,
		Side:    order.Buy,
		SizeUSD: sizeUSD,
		LimitPx: p.Mid,
		Type:    order.GTC,
	}
	res, err := h.paper.Submit(ctx, intent)
	if err != nil {
		return "", fmt.Errorf("下单失败: %s", err.Error())
	}
	entryTick := feed.Tick{
		AssetID: p.AssetID, Market: p.Market,
		Time: now, Mid: res.AvgPrice,
	}
	pos, err := h.pm.OpenSized(p.AssetID, p.Market, entryTick, sizeUSD)
	if err != nil {
		return "", fmt.Errorf("开仓失败: %s", err.Error())
	}
	h.exit.Open(p.AssetID, p.Market, entryTick)
	stats := h.pm.Stats()
	slog.Info("manual_open",
		"id", pos.ID,
		"order_id", res.OrderID,
		"asset", short(p.AssetID),
		"q", h.assetToQ[p.AssetID],
		"size_usd", sizeUSD,
		"signal_mid", p.Mid,
		"entry_fill", res.AvgPrice,
		"units", pos.Units,
		"open_positions", stats.Open,
		"total_exposure_usd", stats.TotalExposure,
	)
	return fmt.Sprintf("✅ %gU @ %.4f · order %s", sizeUSD, res.AvgPrice, short(res.OrderID)), nil
}

// buildNotifier returns a Telegram notifier when TELEGRAM_BOT_TOKEN + _CHAT_ID
// are present, otherwise a Nop so the trading loop is unconditional.
func buildNotifier() notify.Notifier {
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		slog.Info("notify.ready", "mode", "nop", "reason", "telegram_env_missing")
		return notify.Nop{}
	}
	slog.Info("notify.ready", "mode", "telegram", "chat_id", chat)
	return notify.NewTelegram(notify.TelegramConfig{BotToken: tok, ChatID: chat})
}
