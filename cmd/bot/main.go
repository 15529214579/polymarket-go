package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	nethttp "net/http"
	neturl "net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/15529214579/polymarket-go/internal/config"
	"github.com/15529214579/polymarket-go/internal/feed"
	"github.com/15529214579/polymarket-go/internal/journal"
	"github.com/15529214579/polymarket-go/internal/notify"
	"github.com/15529214579/polymarket-go/internal/order"
	"github.com/15529214579/polymarket-go/internal/risk"
	"github.com/15529214579/polymarket-go/internal/strategy"
)

func main() {
	mode := flag.String("mode", "run", "run | discover | feed | sample | detect | prompt-test | daily-report")
	maxMarkets := flag.Int("markets", 20, "top-N sports markets (LoL + NBA daily/playoffs + EPL daily) by vol24h to subscribe")
	windowSec := flag.Int("window", 60, "sampler window in seconds")
	slippageBp := flag.Float64("slippage_bp", 0, "paper fill slippage in bp applied against you")
	largeFillUSD := flag.Float64("large_fill_usd", 3.0, "DM notifier threshold on |realized pnl|")
	envFile := flag.String("env_file", ".env.local", "dotenv file to load before reading env")
	signalMode := flag.String("signal_mode", "auto", "auto (paper-submit on signal) | prompt (DM + Buy 1/5/10 inline buttons, boss picks size)")
	journalDir := flag.String("journal_dir", "db/journal", "trade-journal directory (one JSONL per SGT day)")
	reportDay := flag.String("report_day", "", "daily-report mode: SGT day YYYY-MM-DD (default: yesterday SGT)")
	reportPush := flag.Bool("report_push", false, "daily-report mode: also push summary via Telegram alert bot")
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
		if err := runDetect(ctx, *maxMarkets, *windowSec, *slippageBp, *largeFillUSD, *signalMode, *journalDir); err != nil && ctx.Err() == nil {
			slog.Error("detect failed", "err", err)
			os.Exit(1)
		}
	case "daily-report":
		if err := runDailyReport(ctx, *journalDir, *reportDay, *reportPush); err != nil {
			slog.Error("daily-report failed", "err", err)
			os.Exit(1)
		}
	case "prompt-test":
		if err := runPromptTest(ctx, *slippageBp); err != nil && ctx.Err() == nil {
			slog.Error("prompt-test failed", "err", err)
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
	mkts := feed.FilterSports(all)
	slog.Info("gamma.discover",
		"total_active", len(all),
		"sports", len(mkts),
		"lol", len(feed.FilterLoL(all)),
		"nba_epl_playoffs", len(mkts)-len(feed.FilterLoL(all)),
	)
	for _, m := range mkts {
		tokens := m.ClobTokenIDs()
		slog.Info("sports_market",
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
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	slog.Info("feed.start", "markets", len(mkts), "assets", len(assetIDs))

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
					"q", metaQ(meta, ev.AssetID),
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
					"q", metaQ(meta, tr.AssetID),
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
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	slog.Info("sample.start", "markets", len(mkts), "assets", len(assetIDs), "window_sec", windowSec)

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
					"q", metaQ(meta, t.AssetID),
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
						"q", metaQ(meta, w.AssetID),
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

func runDetect(ctx context.Context, topN, windowSec int, slippageBp, largeFillUSD float64, signalMode, journalDir string) error {
	if signalMode != "auto" && signalMode != "prompt" {
		return fmt.Errorf("invalid signal_mode %q (want auto|prompt)", signalMode)
	}
	jrn, err := journal.New(journalDir)
	if err != nil {
		return fmt.Errorf("journal init: %w", err)
	}
	defer jrn.Close()
	src := newSourceTracker()
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 500)
	if err != nil {
		return err
	}
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
	}
	if topN > len(mkts) {
		topN = len(mkts)
	}
	mkts = mkts[:topN]

	meta := buildAssetMeta(mkts)
	assetIDs := make([]string, 0, len(meta))
	for id := range meta {
		assetIDs = append(assetIDs, id)
	}
	slog.Info("detect.start",
		"markets", len(mkts),
		"lol", countBy(mkts, feed.IsLoLMarket),
		"basketball", countBy(mkts, feed.IsBasketballMarket),
		"football", countBy(mkts, feed.IsFootballMarket),
		"assets", len(assetIDs),
		"window_sec", windowSec,
	)

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
	pending := notify.NewPendingStore(10 * time.Minute)
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
				pm:           pm,
				exit:         exit,
				paper:        paper,
				rm:           rm,
				pending:      pending,
				notifier:     notifier,
				meta:         meta,
				src:          src,
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
	// For each evicted entry we rewrite the original DM to "已过期" and strip
	// its keyboard — so the boss's chat history shows the outcome of every
	// prompt (Phase 3.5 B).
	go func() {
		tk := time.NewTicker(15 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-tk.C:
				evicted := pending.Reap(now)
				if len(evicted) == 0 {
					continue
				}
				edited := 0
				for _, p := range evicted {
					if p.MessageID != 0 {
						notifier.EditSignalExpired(p.MessageID)
						edited++
					}
				}
				slog.Info("pending_reap",
					"expired", len(evicted),
					"edited_expired_dm", edited,
					"remaining", pending.Size(),
				)
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
								Question: metaQ(meta, sig.AssetID),
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
						source, openOID := src.Take(sig.AssetID)
						if err := jrn.Append(journal.TradeRecord{
							ID: closed.ID, AssetID: closed.AssetID, Market: closed.Market,
							Question:     metaQ(meta, closed.AssetID),
							Outcome:      metaOutcome(meta, closed.AssetID),
							Side:         "buy",
							SizeUSD:      closed.SizeUSD,
							Units:        closed.Units,
							EntryMid:     closed.EntryMid,
							EntryTime:    closed.EntryTime,
							ExitMid:      closed.ExitMid,
							ExitTime:     closed.ExitTime,
							ExitReason:   string(closed.ExitReason),
							HeldSec:      int(sig.HeldFor.Seconds()),
							PnLUSD:       closed.PnLUSD,
							OpenOrderID:  openOID,
							CloseOrderID: res.OrderID,
							Mode:         "paper",
							SignalSource: source,
						}); err != nil {
							slog.Warn("journal_append_fail", "asset", short(sig.AssetID), "err", err.Error())
						}
						slog.Info("exit",
							"asset", short(sig.AssetID),
							"q", metaQ(meta, sig.AssetID),
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
					"q", metaQ(meta, sig.AssetID),
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
						"q", metaQ(meta, sig.AssetID),
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
						"q", metaQ(meta, sig.AssetID),
						"reason", "already_open",
					)
					continue
				}

				// Prompt mode: publish the signal as a DM with one button row per
				// outcome (YES/NO or team-A/team-B) and stash the full Choices slice
				// in the pending store. The callback longpoll (above) claims the
				// nonce, picks Choices[slot], and executes via buyHandler.
				if signalMode == "prompt" {
					me := meta[sig.AssetID]
					choices := []notify.Choice{{
						AssetID: sig.AssetID, Outcome: outcomeOrDefault(me, "Yes"),
						Mid: sig.Mid, IsSignal: true,
					}}
					sigChoices := []notify.SignalChoice{{
						Slot: 0, Outcome: choices[0].Outcome, Mid: sig.Mid, IsSignal: true,
					}}
					if me != nil && me.Sibling != "" {
						sibMid := 1.0 - sig.Mid // fallback: binary complement
						if w, ok := sampler.Window(me.Sibling); ok && w.Samples > 0 {
							sibMid = w.EndMid
						}
						sibOutcome := me.SiblingOutcome
						if sibOutcome == "" {
							sibOutcome = "No"
						}
						choices = append(choices, notify.Choice{
							AssetID: me.Sibling, Outcome: sibOutcome, Mid: sibMid,
						})
						sigChoices = append(sigChoices, notify.SignalChoice{
							Slot: 1, Outcome: sibOutcome, Mid: sibMid,
						})
					}
					p := pending.Put(notify.PendingIntent{
						Market:   sig.Market,
						Question: metaQ(meta, sig.AssetID),
						Choices:  choices,
					}, time.Now())
					var match, ctxLine, endIn string
					if me != nil {
						match = me.Match
						ctxLine = me.Context
						endIn = notify.HumanizeEndIn(time.Now(), me.EndTime)
					}
					nonceSnap := p.Nonce
					notifier.SignalPrompt(notify.SignalPromptEvent{
						Nonce:     p.Nonce,
						Match:     match,
						Context:   ctxLine,
						EndIn:     endIn,
						Choices:   sigChoices,
						DeltaPP:   sig.DeltaPP,
						TailUps:   sig.TailUps,
						TailLen:   sig.TailLen,
						BuyRatio:  sig.BuyRatio,
						ExpiresIn: 10 * time.Minute,
						OnSent: func(msgID int64, err error) {
							if err != nil || msgID == 0 {
								return
							}
							pending.SetMessageID(nonceSnap, msgID)
						},
					})
					slog.Info("signal_prompt_sent",
						"asset", short(sig.AssetID),
						"nonce", p.Nonce,
						"mid", sig.Mid,
						"choices", len(choices),
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
						"q", metaQ(meta, sig.AssetID),
						"order_id", res.OrderID,
						"reason", err.Error(),
					)
					continue
				}
				exit.Open(sig.AssetID, sig.Market, entryTick)
				src.Mark(sig.AssetID, "auto", res.OrderID)
				stats := pm.Stats()
				slog.Info("paper_open",
					"id", pos.ID,
					"order_id", res.OrderID,
					"asset", short(sig.AssetID),
					"q", metaQ(meta, sig.AssetID),
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
						"q", metaQ(meta, w.AssetID),
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

func countBy(ms []feed.Market, pred func(feed.Market) bool) int {
	n := 0
	for _, m := range ms {
		if pred(m) {
			n++
		}
	}
	return n
}

func short(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:6] + ".." + id[len(id)-4:]
}

// assetMeta carries per-asset context used by signal prompts and log lines.
// Built once at startup from the gamma market list so the hot path never
// touches gamma again.
type assetMeta struct {
	Question       string
	Match          string // parsed title, e.g. "LoL: Shifters vs G2 Esports"
	Context        string // parsed context, e.g. "Game 1 Winner" or "BO3 · LCK ..."
	Outcome        string // this asset's outcome label ("Shifters", "Yes", ...)
	Sibling        string // sibling asset_id (the other outcome) — empty if market is non-binary
	SiblingOutcome string
	EndTime        time.Time // parsed from market.EndDate; zero if unparseable
}

// buildAssetMeta walks a market list and produces an asset_id-keyed view that
// pairs each asset with its sibling outcome. Only binary markets get sibling
// info; multi-outcome markets (rare in LoL) degrade to "no sibling" which
// renders as a single-row prompt.
func buildAssetMeta(ms []feed.Market) map[string]*assetMeta {
	out := make(map[string]*assetMeta, len(ms)*2)
	for _, m := range ms {
		tokens := m.ClobTokenIDs()
		outcomes := m.Outcomes()
		match, ctx := notify.ParseMarketTitle(m.Question)
		var endTime time.Time
		if m.EndDate != "" {
			if t, err := time.Parse(time.RFC3339, m.EndDate); err == nil {
				endTime = t
			}
		}
		for i, id := range tokens {
			if id == "" {
				continue
			}
			me := &assetMeta{
				Question: m.Question,
				Match:    match,
				Context:  ctx,
				EndTime:  endTime,
			}
			if i < len(outcomes) {
				me.Outcome = outcomes[i]
			}
			// Sibling: for a 2-outcome market, the "other" token.
			if len(tokens) == 2 {
				sibIdx := 1 - i
				me.Sibling = tokens[sibIdx]
				if sibIdx < len(outcomes) {
					me.SiblingOutcome = outcomes[sibIdx]
				}
			}
			out[id] = me
		}
	}
	return out
}

// metaOutcome returns the Outcome label for an asset, or "" if unknown.
func metaOutcome(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Outcome
	}
	return ""
}

// metaQ returns the Question string for an asset, or "" if unknown. Used by log
// lines that previously indexed a plain map[string]string.
func metaQ(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Question
	}
	return ""
}

// metaMatch returns the parsed match title for an asset, or "" if unknown.
// Falls back to empty so the FillReceipt formatter can use Question instead.
func metaMatch(m map[string]*assetMeta, id string) string {
	if me := m[id]; me != nil {
		return me.Match
	}
	return ""
}

// outcomeOrDefault pulls the outcome label for a meta entry, falling back to
// def when the market has no outcome list (rare but defensive).
func outcomeOrDefault(me *assetMeta, def string) string {
	if me == nil || me.Outcome == "" {
		return def
	}
	return me.Outcome
}

// runPromptTest drives the Phase 3.5.b button loop end-to-end against real
// Telegram APIs, without needing a live momentum signal. Picks the top-vol LoL
// market, seeds a PendingStore entry, sends a SignalPrompt DM, runs the sidecar
// LongPoll, and waits up to 2 minutes for one click — then logs manual_open
// or the error toast and exits. Paper only.
func runPromptTest(ctx context.Context, slippageBp float64) error {
	gc := feed.NewGammaClient()
	all, err := gc.ListActiveMarkets(ctx, 200)
	if err != nil {
		return err
	}
	mkts := feed.FilterSports(all)
	if len(mkts) == 0 {
		return fmt.Errorf("no active sports markets")
	}
	top := mkts[0]
	tokens := top.ClobTokenIDs()
	if len(tokens) == 0 {
		return fmt.Errorf("top market has no clob tokens: %s", top.Slug)
	}
	assetID := tokens[0]
	meta := buildAssetMeta(mkts)

	posCfg := strategy.DefaultPositionConfig()
	pm := strategy.NewPositionManager(posCfg)
	exit := strategy.NewExitTracker(strategy.DefaultExitConfig())
	paper := order.NewPaperClient(slippageBp)
	rm := risk.New(risk.DefaultConfig(), time.Now())
	notifier := buildNotifier()
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = notifier.Close(sctx)
	}()
	pending := notify.NewPendingStore(2 * time.Minute)

	sidecarToken := os.Getenv("SIDECAR_BOT_TOKEN")
	sidecarChat := os.Getenv("SIDECAR_CHAT_ID")
	if sidecarChat == "" {
		sidecarChat = os.Getenv("TELEGRAM_CHAT_ID")
	}
	if sidecarToken == "" || sidecarChat == "" {
		return fmt.Errorf("prompt-test needs SIDECAR_BOT_TOKEN and chat_id in env")
	}
	chatID, err := strconv.ParseInt(sidecarChat, 10, 64)
	if err != nil {
		return fmt.Errorf("parse chat_id: %w", err)
	}

	h := &buyHandler{
		pm: pm, exit: exit, paper: paper, rm: rm,
		pending: pending, notifier: notifier,
		meta: meta, largeFillUSD: 3.0,
	}
	lp := notify.NewLongPoll(notify.LongPollConfig{
		BotToken: sidecarToken, ExpectedChatID: chatID,
	}, h)

	lpCtx, lpCancel := context.WithCancel(ctx)
	defer lpCancel()
	lpDone := make(chan error, 1)
	go func() {
		slog.Info("prompt_test.longpoll_start", "chat_id", chatID)
		lpDone <- lp.Run(lpCtx)
	}()

	// Seed one pending intent at the current mid (use 0.50 as a placeholder —
	// paper fill math is the same since slippage is bp-relative). Build full
	// Choices so the prompt shows both YES/NO (or team-A/team-B) rows.
	mid := 0.50
	me := meta[assetID]
	choices := []notify.Choice{{
		AssetID: assetID, Outcome: outcomeOrDefault(me, "Yes"),
		Mid: mid, IsSignal: true,
	}}
	sigChoices := []notify.SignalChoice{{
		Slot: 0, Outcome: choices[0].Outcome, Mid: mid, IsSignal: true,
	}}
	if me != nil && me.Sibling != "" {
		sibOutcome := me.SiblingOutcome
		if sibOutcome == "" {
			sibOutcome = "No"
		}
		sibMid := 1.0 - mid
		choices = append(choices, notify.Choice{
			AssetID: me.Sibling, Outcome: sibOutcome, Mid: sibMid,
		})
		sigChoices = append(sigChoices, notify.SignalChoice{
			Slot: 1, Outcome: sibOutcome, Mid: sibMid,
		})
	}
	p := pending.Put(notify.PendingIntent{
		Market:   top.Slug,
		Question: top.Question,
		Choices:  choices,
	}, time.Now())

	var match, ctxLine, endIn string
	if me != nil {
		match = me.Match
		ctxLine = me.Context
		endIn = notify.HumanizeEndIn(time.Now(), me.EndTime)
	}
	if ctxLine != "" {
		ctxLine += " · [PROMPT-TEST]"
	} else {
		ctxLine = "[PROMPT-TEST]"
	}
	nonceSnap := p.Nonce
	notifier.SignalPrompt(notify.SignalPromptEvent{
		Nonce:     p.Nonce,
		Match:     match,
		Context:   ctxLine,
		EndIn:     endIn,
		Choices:   sigChoices,
		DeltaPP:   0,
		TailUps:   0,
		TailLen:   0,
		BuyRatio:  0,
		ExpiresIn: 2 * time.Minute,
		OnSent: func(msgID int64, err error) {
			if err != nil || msgID == 0 {
				return
			}
			pending.SetMessageID(nonceSnap, msgID)
		},
	})
	slog.Info("prompt_test.sent",
		"nonce", p.Nonce,
		"q", top.Question,
		"mid", mid,
		"asset", short(assetID),
		"choices", len(choices),
	)

	// Wait for either: first fill (pm.Stats().Open > 0), the user's click
	// arriving but failing (history shows a Result), context cancel, or timeout.
	deadline := time.After(2 * time.Minute)
	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-deadline:
			slog.Warn("prompt_test.timeout", "hint", "no click within 2 min")
			lpCancel()
			<-lpDone
			return nil
		case <-tk.C:
			if pm.Stats().Open > 0 {
				// Give the toast 1s to flush, then exit.
				time.Sleep(1 * time.Second)
				lpCancel()
				<-lpDone
				slog.Info("prompt_test.done", "open_positions", pm.Stats().Open)
				return nil
			}
		}
	}
}

// buyHandler wires a click on one outcome's Buy 1/5/10 → same paper-submit →
// pm.Open path the auto-mode signal loop uses, but honors the size the boss
// picked and the Choice (YES/NO) resolved from PendingIntent.Choices[slot].
// Executes synchronously on the longpoll goroutine; Telegram dispatch of the
// resulting DM is async via notifier.
type buyHandler struct {
	pm           *strategy.PositionManager
	exit         *strategy.ExitTracker
	paper        order.Client
	rm           *risk.Manager
	pending      *notify.PendingStore
	notifier     notify.Notifier
	meta         map[string]*assetMeta
	src          *sourceTracker
	largeFillUSD float64
}

func (h *buyHandler) OnBuy(ctx context.Context, nonce string, slot int, sizeUSD float64) (string, error) {
	now := time.Now()
	p, ok := h.pending.Claim(nonce, now)
	if !ok {
		return "", fmt.Errorf("已过期或已点过")
	}
	if slot < 0 || slot >= len(p.Choices) {
		return "", fmt.Errorf("选项越界 slot=%d", slot)
	}
	choice := p.Choices[slot]
	if err := h.rm.AllowOpen(now); err != nil {
		st := h.rm.State()
		return "", fmt.Errorf("风控阻止: %s (day_pnl=%.2f)", st.BlockReason, st.DayRealizedPnL)
	}
	if h.pm.Has(choice.AssetID) || h.pm.HasMarket(p.Market) {
		return "", fmt.Errorf("已有同市场仓位")
	}
	intent := order.Intent{
		AssetID: choice.AssetID,
		Market:  p.Market,
		Side:    order.Buy,
		SizeUSD: sizeUSD,
		LimitPx: choice.Mid,
		Type:    order.GTC,
	}
	res, err := h.paper.Submit(ctx, intent)
	if err != nil {
		return "", fmt.Errorf("下单失败: %s", err.Error())
	}
	entryTick := feed.Tick{
		AssetID: choice.AssetID, Market: p.Market,
		Time: now, Mid: res.AvgPrice,
	}
	pos, err := h.pm.OpenSized(choice.AssetID, p.Market, entryTick, sizeUSD)
	if err != nil {
		return "", fmt.Errorf("开仓失败: %s", err.Error())
	}
	h.exit.Open(choice.AssetID, p.Market, entryTick)
	if h.src != nil {
		h.src.Mark(choice.AssetID, "manual", res.OrderID)
	}
	stats := h.pm.Stats()
	slog.Info("manual_open",
		"id", pos.ID,
		"order_id", res.OrderID,
		"asset", short(choice.AssetID),
		"q", metaQ(h.meta, choice.AssetID),
		"outcome", choice.Outcome,
		"slot", slot,
		"size_usd", sizeUSD,
		"signal_mid", choice.Mid,
		"entry_fill", res.AvgPrice,
		"units", pos.Units,
		"open_positions", stats.Open,
		"total_exposure_usd", stats.TotalExposure,
	)

	// Phase 3.5 B + C: rewrite the original prompt to "已下单 …" + strip buttons,
	// and push a durable archive DM so the boss has a permanent record.
	receipt := notify.FillReceiptEvent{
		Question: metaQ(h.meta, choice.AssetID),
		Match:    metaMatch(h.meta, choice.AssetID),
		Outcome:  choice.Outcome,
		SizeUSD:  sizeUSD,
		Units:    pos.Units,
		FillPx:   res.AvgPrice,
		OrderID:  res.OrderID,
		Source:   "manual",
	}
	if h.notifier != nil {
		h.notifier.EditSignalFilled(receipt, p.MessageID)
		h.notifier.FillReceipt(receipt)
	}

	return fmt.Sprintf("✅ %s %gU @ %.4f · order %s",
		choice.Outcome, sizeUSD, res.AvgPrice, short(res.OrderID)), nil
}

// buildNotifier returns a Telegram notifier when TELEGRAM_BOT_TOKEN + _CHAT_ID
// are present, otherwise a Nop so the trading loop is unconditional.
//
// SIDECAR_BOT_TOKEN, when set, routes SignalPrompt messages through that bot so
// the inline-keyboard message originates from the same bot the LongPoll watches
// for callback_query (Phase 3.5.b). Other events stay on the alert bot.
func buildNotifier() notify.Notifier {
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		slog.Info("notify.ready", "mode", "nop", "reason", "telegram_env_missing")
		return notify.Nop{}
	}
	cfg := notify.TelegramConfig{BotToken: tok, ChatID: chat, PromptBotToken: os.Getenv("SIDECAR_BOT_TOKEN")}
	slog.Info("notify.ready",
		"mode", "telegram",
		"chat_id", chat,
		"prompt_via_sidecar", cfg.PromptBotToken != "",
	)
	return notify.NewTelegram(cfg)
}

// sourceTracker remembers which path opened a position (auto detector vs manual
// click) so the journal can attribute closed trades correctly. Position state
// itself is in PositionManager; this is a small sidecar keyed by AssetID.
type sourceTracker struct {
	mu sync.Mutex
	m  map[string]sourceEntry
}

type sourceEntry struct {
	source      string // "auto" or "manual"
	openOrderID string
}

func newSourceTracker() *sourceTracker {
	return &sourceTracker{m: map[string]sourceEntry{}}
}

func (s *sourceTracker) Mark(assetID, source, openOID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[assetID] = sourceEntry{source: source, openOrderID: openOID}
}

// Take returns the recorded source + open order id and removes the entry.
// Missing entries default to "auto" with empty order id (safe for legacy).
func (s *sourceTracker) Take(assetID string) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[assetID]
	if !ok {
		return "auto", ""
	}
	delete(s.m, assetID)
	return e.source, e.openOrderID
}

// runDailyReport reads one SGT day's trade journal, prints the summary to
// stdout, and (when -report_push is set) DMs it via the Telegram alert bot.
// Default day = yesterday SGT — this matches the cron firing at 00:00:30 SGT
// to summarize the day that just ended.
func runDailyReport(ctx context.Context, dir, day string, push bool) error {
	if day == "" {
		yesterday := time.Now().In(journal.SGT).AddDate(0, 0, -1)
		day = yesterday.Format("2006-01-02")
	}
	trades, err := journal.Read(dir, day)
	if err != nil {
		return fmt.Errorf("read journal: %w", err)
	}
	summary := journal.Summarize(day, trades)
	out := journal.FormatTelegram(summary)
	fmt.Print(out)
	if !push {
		return nil
	}
	tok := os.Getenv("TELEGRAM_BOT_TOKEN")
	chat := os.Getenv("TELEGRAM_CHAT_ID")
	if tok == "" || chat == "" {
		return fmt.Errorf("report_push: TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID missing")
	}
	if err := sendTelegram(ctx, tok, chat, out); err != nil {
		return fmt.Errorf("telegram push: %w", err)
	}
	slog.Info("daily_report.pushed", "day", day, "trades", summary.Trades, "pnl_usd", summary.RealizedPnLUSD)
	return nil
}

// sendTelegram fires a single sendMessage to the alert bot. We don't use the
// notify package because that's wired for fire-and-forget queues; the cron
// wants a synchronous send-and-exit.
func sendTelegram(ctx context.Context, token, chat, body string) error {
	api := "https://api.telegram.org/bot" + token + "/sendMessage"
	form := neturl.Values{}
	form.Set("chat_id", chat)
	form.Set("text", body)
	form.Set("disable_web_page_preview", "true")
	req, err := nethttp.NewRequestWithContext(ctx, "POST", api, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cl := &nethttp.Client{Timeout: 10 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("telegram http %d", resp.StatusCode)
	}
	return nil
}
