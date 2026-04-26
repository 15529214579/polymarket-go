package btc

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"
)

type MomentumAlert struct {
	Move1h   float64
	Move15m  float64
	Spot     float64
	OpenHour float64
	IsSharp  bool
}

func DetectMomentum(ctx context.Context) (*MomentumAlert, error) {
	candles5m, err := FetchCandles(ctx, "BTCUSDT", Interval5m, 12)
	if err != nil {
		return nil, fmt.Errorf("fetch 5m: %w", err)
	}
	if len(candles5m) < 12 {
		return nil, fmt.Errorf("too few 5m candles: %d", len(candles5m))
	}

	spot := candles5m[len(candles5m)-1].Close
	openHour := candles5m[0].Open
	open15m := candles5m[len(candles5m)-3].Open

	move1h := (spot - openHour) / openHour * 100
	move15m := (spot - open15m) / open15m * 100

	isSharp := math.Abs(move1h) >= 1.5 || math.Abs(move15m) >= 0.8

	alert := &MomentumAlert{
		Move1h:   move1h,
		Move15m:  move15m,
		Spot:     spot,
		OpenHour: openHour,
		IsSharp:  isSharp,
	}

	if isSharp {
		slog.Info("momentum.sharp_move",
			"move_1h_pct", fmt.Sprintf("%+.2f%%", move1h),
			"move_15m_pct", fmt.Sprintf("%+.2f%%", move15m),
			"spot", fmt.Sprintf("%.0f", spot),
		)
	}

	return alert, nil
}

type MomentumWatcher struct {
	checkInterval time.Duration
	onSharpMove   func()
	lastAlert     time.Time
	cooldown      time.Duration
}

func NewMomentumWatcher(onSharpMove func()) *MomentumWatcher {
	return &MomentumWatcher{
		checkInterval: 2 * time.Minute,
		onSharpMove:   onSharpMove,
		cooldown:      15 * time.Minute,
	}
}

func (w *MomentumWatcher) Run(ctx context.Context) {
	slog.Info("momentum_watcher.started", "check_interval", w.checkInterval.String())
	tk := time.NewTicker(w.checkInterval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			alert, err := DetectMomentum(ctx)
			if err != nil {
				slog.Warn("momentum_watcher.check_fail", "err", err.Error())
				continue
			}
			if alert.IsSharp && time.Since(w.lastAlert) > w.cooldown {
				w.lastAlert = time.Now()
				slog.Info("momentum_watcher.triggering_scan",
					"move_1h", fmt.Sprintf("%+.2f%%", alert.Move1h),
					"move_15m", fmt.Sprintf("%+.2f%%", alert.Move15m),
				)
				w.onSharpMove()
			}
		}
	}
}
