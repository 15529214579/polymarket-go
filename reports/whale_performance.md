# Whale Performance Report

**Generated:** 2026-08-02 00:14 +08

- Log: `/Users/murphyma/work/polymarket-go/db/journal/whale_trades.jsonl`
- Wallet filter: `/Users/murphyma/work/polymarket-go/wallets.strategy-push.txt,/Users/murphyma/work/polymarket-go/wallets.football-score-push.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-push.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-watch.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-sports-push.txt,/Users/murphyma/work/polymarket-go/wallets.sports-holders-push.txt`
- Fixed stake: $10.00 per BUY signal
- Minimum whale notional: $500
- List minimum notionals: core=$750, flow=$1000, leaderboard_push=$3000, leaderboard_watch=$3000, scout=$500, sports=$1000, tape=$1000, target=$500, watch=$500
- Policy since: 2026-08-01T16:11:25Z
- Repeat cooldown: 3m0s per wallet+asset; BUYs >= $5000 bypass cooldown

## Summary

- Raw matched BUY signals: 2
- Suppressed repeat BUYs: 0
- Logged asset-cooldown BUYs: 0
- Logged event-cooldown BUYs: 3
- Logged pending-consensus BUYs: 0
- Logged duplicate BUYs: 0
- Duplicate BUY alerts ignored: 0
- Evaluated signals: 2
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 2
- Still open/unmarked: 0
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $-0.09
- ROI incl. midpoint marks: -0.4%
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
| `0x3af1...d739` | sports_holders_push | D | 0 | 1 | 0 | 0 | 1 | $2747 | 08-02 00:12 |
| `0xf201...527e` | sports_holders_push | BOT | 0 | 1 | 0 | 0 | 1 | $2421 | 08-02 00:13 |
| `0xf68a...5b1b` | sports_holders_push | BOT | 0 | 1 | 0 | 0 | 1 | $1090 | 08-02 00:11 |

### By Wallet-Event

| Event | Wallet | List | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| lol: giantx vs sk gaming | `0x3af1...d739` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $2747 | 08-02 00:12 |
| lol: giantx vs sk gaming | `0xf201...527e` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $2421 | 08-02 00:13 |
| lol: vivo keyd stars vs furia esports | `0xf68a...5b1b` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $1090 | 08-02 00:11 |

## Event Cluster Summary

- Independent wallet-event clusters: 2
- Event-cluster win rate: 50.0%
- Event-cluster ROI: -0.4%
- Event-cluster PnL: $-0.09

## Event-Capped Strategy

- Rule: one fixed-stake entry per wallet-event cluster, using the first evaluated signal
- Entries: 2
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 2
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $-0.09
- ROI incl. midpoint marks: -0.4%
- Proven entries: 0
- Proven win rate: 0.0%
- Proven PnL: $+0.00
- Proven ROI: 0.0%

## Event-Capped By List

| List | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|
| sports_holders_push | 2 | 0 | 0 | 2 | 50.0% | -0.4% | $-0.09 |

## Event-Capped By Wallet

| Wallet | List | Tier | Label | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| `0x9e3e...f882` | sports_holders_push | BOT | w_9e3e…f882 | 2 | 0 | 0 | 2 | 50.0% | -0.4% | $-0.09 |
## By List

| List | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|
| sports_holders_push | 2 | 0 | 0 | 2 | 50.0% | -0.4% | $-0.09 |
## By Wallet

| Wallet | List | Tier | Label | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| `0x9e3e...f882` | sports_holders_push | BOT | w_9e3e…f882 | 2 | 0 | 0 | 2 | 50.0% | -0.4% | $-0.09 |

## By Event Cluster

| Event | Wallet | List | Tier | Signals | Markets | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| valorant: karmine corp vs team liquid | `0x9e3e...f882` | sports_holders_push | BOT | 1 | 1 | 0 | 0 | 1 | 0.0% | -0.9% | $-0.09 |
| lol: vivo keyd stars vs furia esports | `0x9e3e...f882` | sports_holders_push | BOT | 1 | 1 | 0 | 0 | 1 | 100.0% | 0.1% | $+0.01 |

## Event-Capped Entries

| Time | Event | Wallet | List | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---:|---:|---|---:|---:|---|
| 08-02 00:12 | valorant: karmine corp vs team ... | `0x9e3e...f882` | sports_holders_push | Team Liquid | 0.5300 | 0.5250 | mark | -0.9% | $-0.09 | Valorant: Karmine Corp vs Team Liquid (... |
| 08-02 00:12 | lol: vivo keyd stars vs furia e... | `0x9e3e...f882` | sports_holders_push | FURIA Esports | 0.4597 | 0.4600 | mark | +0.1% | $+0.01 | LoL: Vivo Keyd Stars vs FURIA Esports -... |

## Recent Signals

| Time | Wallet | List | Tier | Side | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---|---:|---:|---|---:|---:|---|
| 08-02 00:12 | `0x9e3e...f882` | sports_holders_push | BOT | BUY | Team Liquid | 0.5300 | 0.5250 | mark | -0.9% | $-0.09 | Valorant: Karmine Corp vs Team Liquid (BO3) -... |
| 08-02 00:12 | `0x9e3e...f882` | sports_holders_push | BOT | BUY | FURIA Esports | 0.4597 | 0.4600 | mark | +0.1% | $+0.01 | LoL: Vivo Keyd Stars vs FURIA Esports - Game ... |
