# Whale Performance Report

**Generated:** 2026-08-03 00:15 +08

- Log: `/Users/murphyma/work/polymarket-go/db/journal/whale_trades.jsonl`
- Wallet filter: `/Users/murphyma/work/polymarket-go/wallets.strategy-push.txt,/Users/murphyma/work/polymarket-go/wallets.football-score-push.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-push.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-watch.txt,/Users/murphyma/work/polymarket-go/wallets.leaderboard-sports-push.txt,/Users/murphyma/work/polymarket-go/wallets.sports-holders-push.txt`
- Fixed stake: $10.00 per BUY signal
- Minimum whale notional: $500
- List minimum notionals: core=$750, flow=$1000, leaderboard_push=$3000, leaderboard_watch=$3000, scout=$500, sports=$1000, tape=$1000, target=$500, watch=$500
- Policy since: 2026-08-02T15:36:45Z
- Repeat cooldown: 3m0s per wallet+asset; BUYs >= $5000 bypass cooldown

## Summary

- Raw matched BUY signals: 6
- Suppressed repeat BUYs: 0
- Logged asset-cooldown BUYs: 6
- Logged event-cooldown BUYs: 5
- Logged pending-consensus BUYs: 0
- Logged duplicate BUYs: 0
- Duplicate BUY alerts ignored: 0
- Evaluated signals: 6
- Realized via whale SELL: 3
- Settled by market resolution: 0
- Marked to current midpoint: 3
- Still open/unmarked: 0
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $+5.83
- ROI incl. midpoint marks: 9.7%
- Proven signals: 3
- Proven win rate: 66.7%
- Proven PnL: $+5.53
- Proven ROI: 18.4%

## Policy Violations

- Alerted BUYs outside current sports/esports/price policy: 1

| Time | Reason | Wallet | List | Tier | Notional | Price | Outcome | Market |
|---|---|---|---|---|---:|---:|---|---|
| 08-03 00:01 | category_filtered | `0x0993...525f` | sports_holders_push | D | $1854 | 0.4000 | ShindeN | Counter-Strike: ShindeN vs BESTIA (BO3) - StarLadder ... |

## Suppressed BUY Noise

These rows are detected large BUYs that were logged but not pushed/evaluated because duplicate, cooldown, or consensus gates suppressed them.

### By Wallet

| Wallet | List | Tier | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| `0x4aec...4ce9` | sports_holders_push | D | 5 | 1 | 0 | 0 | 6 | $19206 | 08-02 23:50 |
| `0x2929...1dd0` | sports_holders_push | C | 1 | 1 | 0 | 0 | 2 | $2000 | 08-03 00:06 |
| `0x9e3e...f882` | sports_holders_push | BOT | 0 | 2 | 0 | 0 | 2 | $1290 | 08-03 00:11 |
| `0x751b...73fd` | sports_holders_push | D | 0 | 1 | 0 | 0 | 1 | $9500 | 08-02 23:53 |

### By Wallet-Event

| Event | Wallet | List | AssetCD | EventCD | Pending | Dup | Total | Notional | Last |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| dota 2: team falcons vs 1win | `0x4aec...4ce9` | sports_holders_push | 5 | 1 | 0 | 0 | 6 | $19206 | 08-02 23:50 |
| dota 2: rekonix vs yakult brothers | `0x2929...1dd0` | sports_holders_push | 1 | 1 | 0 | 0 | 2 | $2000 | 08-03 00:06 |
| lol: sk gaming vs g2 esports | `0x751b...73fd` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $9500 | 08-02 23:53 |
| counter-strike: sinners vs sparta | `0x9e3e...f882` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $762 | 08-02 23:39 |
| lol: fluxo w7m vs loud | `0x9e3e...f882` | sports_holders_push | 0 | 1 | 0 | 0 | 1 | $529 | 08-03 00:11 |

## Event Cluster Summary

- Independent wallet-event clusters: 6
- Event-cluster win rate: 50.0%
- Event-cluster ROI: 9.7%
- Event-cluster PnL: $+5.83

## Event-Capped Strategy

- Rule: one fixed-stake entry per wallet-event cluster, using the first evaluated signal
- Entries: 6
- Realized via whale SELL: 3
- Settled by market resolution: 0
- Marked to current midpoint: 3
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $+5.83
- ROI incl. midpoint marks: 9.7%
- Proven entries: 3
- Proven win rate: 66.7%
- Proven PnL: $+5.53
- Proven ROI: 18.4%

## Event-Capped By List

| List | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|
| sports_holders_push | 5 | 2 | 0 | 3 | 60.0% | 15.5% | $+7.77 |
| football_score_push | 1 | 1 | 0 | 0 | 0.0% | -19.4% | $-1.94 |

## Event-Capped By Wallet

| Wallet | List | Tier | Label | Entries | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| `0x4aec...4ce9` | sports_holders_push | D | w_4aec…4ce9 | 2 | 2 | 0 | 0 | 100.0% | 37.3% | $+7.47 |
| `0x9e3e...f882` | sports_holders_push | BOT | w_9e3e…f882 | 1 | 0 | 0 | 1 | 100.0% | 4.8% | $+0.48 |
| `0x751b...73fd` | sports_holders_push | D | w_751b…73fd | 1 | 0 | 0 | 1 | 0.0% | -0.5% | $-0.05 |
| `0x0993...525f` | sports_holders_push | D | w_0993…525f | 1 | 0 | 0 | 1 | 0.0% | -1.3% | $-0.13 |
| `0x84d4...77a9` | football_score_push | D | w_84d4…77a9 | 1 | 1 | 0 | 0 | 0.0% | -19.4% | $-1.94 |
## By List

| List | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---:|---:|---:|---:|---:|---:|---:|
| sports_holders_push | 5 | 2 | 0 | 3 | 60.0% | 15.5% | $+7.77 |
| football_score_push | 1 | 1 | 0 | 0 | 0.0% | -19.4% | $-1.94 |
## By Wallet

| Wallet | List | Tier | Label | Signals | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| `0x4aec...4ce9` | sports_holders_push | D | w_4aec…4ce9 | 2 | 2 | 0 | 0 | 100.0% | 37.3% | $+7.47 |
| `0x9e3e...f882` | sports_holders_push | BOT | w_9e3e…f882 | 1 | 0 | 0 | 1 | 100.0% | 4.8% | $+0.48 |
| `0x751b...73fd` | sports_holders_push | D | w_751b…73fd | 1 | 0 | 0 | 1 | 0.0% | -0.5% | $-0.05 |
| `0x0993...525f` | sports_holders_push | D | w_0993…525f | 1 | 0 | 0 | 1 | 0.0% | -1.3% | $-0.13 |
| `0x84d4...77a9` | football_score_push | D | w_84d4…77a9 | 1 | 1 | 0 | 0 | 0.0% | -19.4% | $-1.94 |

## By Event Cluster

| Event | Wallet | List | Tier | Signals | Markets | Closed | Settled | Marked | Win | ROI | PnL |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| counter-strike: shinden vs bestia | `0x0993...525f` | sports_holders_push | D | 1 | 1 | 0 | 0 | 1 | 0.0% | -1.3% | $-0.13 |
| lol: sk gaming vs g2 esports | `0x751b...73fd` | sports_holders_push | D | 1 | 1 | 0 | 0 | 1 | 0.0% | -0.5% | $-0.05 |
| lol: fluxo w7m vs loud | `0x9e3e...f882` | sports_holders_push | BOT | 1 | 1 | 0 | 0 | 1 | 100.0% | 4.8% | $+0.48 |
| dota 2: rekonix vs yakult brothers | `0x4aec...4ce9` | sports_holders_push | D | 1 | 1 | 1 | 0 | 0 | 100.0% | 23.3% | $+2.33 |
| dota 2: team falcons vs 1win | `0x4aec...4ce9` | sports_holders_push | D | 1 | 1 | 1 | 0 | 0 | 100.0% | 51.4% | $+5.14 |
| dota 2: team falcons vs 1win | `0x84d4...77a9` | football_score_push | D | 1 | 1 | 1 | 0 | 0 | 0.0% | -19.4% | $-1.94 |

## Event-Capped Entries

| Time | Event | Wallet | List | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---:|---:|---|---:|---:|---|
| 08-03 00:01 | counter-strike: shinden vs bestia | `0x0993...525f` | sports_holders_push | ShindeN | 0.4000 | 0.3950 | mark | -1.3% | $-0.13 | Counter-Strike: ShindeN vs BESTIA (BO3)... |
| 08-02 23:50 | lol: sk gaming vs g2 esports | `0x751b...73fd` | sports_holders_push | G2 Esports | 0.9500 | 0.9450 | mark | -0.5% | $-0.05 | LoL: SK Gaming vs G2 Esports (BO3) - LE... |
| 08-02 23:48 | lol: fluxo w7m vs loud | `0x9e3e...f882` | sports_holders_push | LOUD | 0.5393 | 0.5650 | mark | +4.8% | $+0.48 | LoL: Fluxo W7M vs LOUD - Game 1 Winner |
| 08-02 23:41 | dota 2: rekonix vs yakult brothers | `0x4aec...4ce9` | sports_holders_push | REKONIX | 0.8100 | 0.9990 | sell | +23.3% | $+2.33 | Dota 2: REKONIX vs Yakult Brothers - Ga... |
| 08-02 23:38 | dota 2: team falcons vs 1win | `0x4aec...4ce9` | sports_holders_push | Team Falcons | 0.6600 | 0.9990 | sell | +51.4% | $+5.14 | Dota 2: Team Falcons vs 1win - Game 2 W... |
| 08-02 23:37 | dota 2: team falcons vs 1win | `0x84d4...77a9` | football_score_push | Team Falcons | 0.6100 | 0.4919 | sell | -19.4% | $-1.94 | Dota 2: Team Falcons vs 1win - Game 2 W... |

## Recent Signals

| Time | Wallet | List | Tier | Side | Outcome | Entry | Exit | Src | Ret | PnL | Market |
|---|---|---|---|---|---|---:|---:|---|---:|---:|---|
| 08-03 00:01 | `0x0993...525f` | sports_holders_push | D | BUY | ShindeN | 0.4000 | 0.3950 | mark | -1.3% | $-0.13 | Counter-Strike: ShindeN vs BESTIA (BO3) - Sta... |
| 08-02 23:50 | `0x751b...73fd` | sports_holders_push | D | BUY | G2 Esports | 0.9500 | 0.9450 | mark | -0.5% | $-0.05 | LoL: SK Gaming vs G2 Esports (BO3) - LEC Regu... |
| 08-02 23:48 | `0x9e3e...f882` | sports_holders_push | BOT | BUY | LOUD | 0.5393 | 0.5650 | mark | +4.8% | $+0.48 | LoL: Fluxo W7M vs LOUD - Game 1 Winner |
| 08-02 23:41 | `0x4aec...4ce9` | sports_holders_push | D | BUY | REKONIX | 0.8100 | 0.9990 | sell | +23.3% | $+2.33 | Dota 2: REKONIX vs Yakult Brothers - Game 1 W... |
| 08-02 23:38 | `0x4aec...4ce9` | sports_holders_push | D | BUY | Team Falcons | 0.6600 | 0.9990 | sell | +51.4% | $+5.14 | Dota 2: Team Falcons vs 1win - Game 2 Winner |
| 08-02 23:37 | `0x84d4...77a9` | football_score_push | D | BUY | Team Falcons | 0.6100 | 0.4919 | sell | -19.4% | $-1.94 | Dota 2: Team Falcons vs 1win - Game 2 Winner |
