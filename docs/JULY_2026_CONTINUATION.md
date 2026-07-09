# July 2026 Continuation Plan

## Goal

Continue the smart-money whale pipeline from July 2026 and keep older April-June paper-trading material archived.

## Current Runtime Shape

- Whale push worker: `scripts/start-whale-push.sh`
- Edge monitor: `scripts/whale-edge-monitor.sh`
- Daily strategy refresh: `scripts/smartmoney-daily.sh`
- Sports flow collection: `scripts/sports-tape-monitor.sh`

The worker currently uses wallet tiers and list gates, with `copytrade_size=0` for research/alerting. Promotion to active copytrading requires a separate review of signal count, ROI, win rate, and live edge.

## Active Strategy Layers

- Core: proven low-bot wallets from copytrade backtests.
- Watch: higher-sample wallets that are promising but not core.
- Sports: basketball/soccer/esports positive-copy wallets.
- Scout/target: research layers for new wallets and category specialists.
- Tape/edge-hot: recent large sports-flow candidates, observation only until promotion gates pass.

## Continue

1. Keep `whale-push` and `whale-edge` running while live samples accumulate.
2. Run `scripts/smartmoney-daily.sh` after meaningful new data.
3. Review `reports/smartmoney_daily.md` first, then drill into `strategy_lab`, `wallet_maintenance`, `whale_edge`, and sports alert reports.
4. Promote only after minimum sample gates pass; current sample counts are still below promotion thresholds.
5. Keep raw wallet activity and large generated datasets out of Git.

## Archive Boundary

The old April-June paper/backtest materials live in `archive/legacy-2026-04-06/`. Do not use those reports as the current operating source of truth unless explicitly comparing strategy generations.
