# Whale Performance Report

**Generated:** 2026-07-08 23:00 +08

- Log: `/Users/murphyma/work/polymarket-go/db/journal/whale_trades.jsonl`
- Wallet filter: `/Users/murphyma/work/polymarket-go/wallets.strategy-push.txt`
- Fixed stake: $10.00 per BUY signal
- Minimum whale notional: $500
- Repeat cooldown: 3m0s per wallet+asset; BUYs >= $5000 bypass cooldown

## Summary

- Raw matched BUY signals: 17
- Suppressed repeat BUYs: 6
- Logged asset-cooldown BUYs: 9
- Logged event-cooldown BUYs: 17
- Logged duplicate BUYs: 0
- Duplicate BUY alerts ignored: 11
- Evaluated signals: 0
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 0
- Still open/unmarked: 11
- Win rate: 0.0%
- PnL: $+0.00
- ROI: 0.0%

## Suppressed BUY Noise

These rows are detected large BUYs that were logged but not pushed/evaluated because duplicate or cooldown gates suppressed them.

### By Wallet

| Wallet | List | Tier | AssetCD | EventCD | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---|
| `0xe16d...5e30` | target | C | 9 | 17 | 0 | 26 | $52000 | 07-08 22:58 |

### By Wallet-Event

| Event | Wallet | List | AssetCD | EventCD | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---|
| dota 2: rekonix vs mouz | `0xe16d...5e30` | target | 8 | 14 | 0 | 22 | $48000 | 07-08 22:58 |
| dota 2: nigma galaxy vs aurora | `0xe16d...5e30` | target | 1 | 3 | 0 | 4 | $4000 | 07-08 22:29 |

## Event Cluster Summary

- Independent wallet-event clusters: 0
- Event-cluster win rate: 0.0%
- Event-cluster ROI: 0.0%
- Event-cluster PnL: $+0.00

## Event-Capped Strategy

- Rule: one fixed-stake entry per wallet-event cluster, using the first evaluated signal
- Entries: 0
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 0
- Win rate: 0.0%
- PnL: $+0.00
- ROI: 0.0%

## Event-Capped By List

| List | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|

## Event-Capped By Wallet

| Wallet | List | Tier | Label | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
## By List

| List | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|
## By Wallet

| Wallet | List | Tier | Label | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|

## By Event Cluster

| Event | Wallet | List | Tier | Signals | Markets | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|

## Event-Capped Entries

| Time | Event | Wallet | List | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---:|---:|---|---:|---:|---|

## Recent Signals

| Time | Wallet | List | Tier | Side | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---|---:|---:|---|---:|---:|---|

No evaluated BUY signals matched the current filters yet. Keep the bot running; this report will become useful after the core list fires.
