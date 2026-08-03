# Whale Performance Report

**Generated:** 2026-08-03 00:15 +08

- Log: `/Users/murphyma/work/polymarket-go/db/journal/whale_trades.jsonl`
- Wallet filter: `/Users/murphyma/work/polymarket-go/wallets.strategy-target.txt`
- Fixed stake: $10.00 per BUY signal
- Minimum whale notional: $500
- List minimum notionals: core=$750, flow=$1000, leaderboard_push=$3000, leaderboard_watch=$3000, scout=$500, sports=$1000, tape=$1000, target=$500, watch=$500
- Policy since: 2026-08-02T15:36:45Z
- Repeat cooldown: 3m0s per wallet+asset; BUYs >= $5000 bypass cooldown

## Summary

- Raw matched BUY signals: 0
- Suppressed repeat BUYs: 0
- Logged asset-cooldown BUYs: 1
- Logged event-cooldown BUYs: 1
- Logged pending-consensus BUYs: 0
- Logged duplicate BUYs: 0
- Duplicate BUY alerts ignored: 0
- Evaluated signals: 0
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 0
- Still open/unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Proven signals: 0
- Proven win rate: 0.0%
- Proven PnL: $+0.00
- Proven ROI: 0.0%

## Policy Violations

- Alerted BUYs outside current sports/esports/price policy: 0

## Suppressed BUY Noise

These rows are detected large BUYs that were logged but not pushed/evaluated because duplicate, cooldown, or consensus gates suppressed them.

### By Wallet

| Wallet | List | Tier | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| `0x2929...1dd0` | target | C | 1 | 1 | 0 | 0 | 2 | $2000 | 08-03 00:06 |

### By Wallet-Event

| Event | Wallet | List | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| dota 2: rekonix vs yakult brothers | `0x2929...1dd0` | target | 1 | 1 | 0 | 0 | 2 | $2000 | 08-03 00:06 |

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
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Proven entries: 0
- Proven win rate: 0.0%
- Proven PnL: $+0.00
- Proven ROI: 0.0%

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
