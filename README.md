# polymarket-go

Go-based Polymarket research, alerting, and paper-trading workspace.

## Current Track: July 2026

The active work is the July smart-money whale pipeline:

- `cmd/bot` runs the whale signal worker.
- `cmd/strategy-lab` rebuilds wallet strategy layers from scored wallets.
- `cmd/whale-report` and `cmd/whale-edge` evaluate live signal quality.
- `cmd/sports-tape` and `cmd/sports-tape-alert` collect and gate sports/esports whale flow.
- `cmd/wallet-discover` and `cmd/wallet-maintain` refresh wallet candidates and quarantine/review lists.

The current production-like worker is started through:

```sh
scripts/start-whale-push.sh
```

The July daily pipeline is:

```sh
scripts/smartmoney-daily.sh
```

The project is still in learning/collection mode. The whale worker is configured for alerts and research with `copytrade_size=0`; do not treat it as an enabled real-money copytrader unless that flag and wallet policy are intentionally changed.

## Useful Reports

- `reports/smartmoney_daily.md` - operator summary for the latest daily run.
- `reports/strategy_lab.md` - selected core/watch/sports/scout/target layers.
- `reports/wallet_maintenance.md` - wallet promotion, review, and quarantine status.
- `reports/whale_edge.md` - measured edge after whale buy signals.
- `reports/sports_tape.md` - recent basketball/soccer/esports whale flow.
- `reports/sports_alert_performance.md` - sports alert performance and promotion gates.

## Legacy Archive

April-June paper-trading and backtest reports were moved under:

```text
archive/legacy-2026-04-06/
```

`SPEC.md`, `TODO.md`, and `BTC_TODO.md` are retained as historical planning material. Prefer this README and `docs/JULY_2026_CONTINUATION.md` for the current operating direction.

## Development

```sh
make test
make build
```

Large raw runtime datasets, local logs, pid files, and wallet activity crawls are intentionally ignored. Commit code, scripts, wallet policy lists, and summarized reports; keep raw market data local unless there is a deliberate reason to publish it.
