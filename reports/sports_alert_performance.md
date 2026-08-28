# Sports Alert Performance

**Generated:** 2026-08-28 00:10 +08

- Alert log: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/sports_tape_alerts.jsonl`
- Fixed paper stake: $10.00 per Telegram alert
- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices

- Current policy exclude wallets: `/Users/murphyma/work/polymarket-go/wallets.strategy-quarantine.txt,/Users/murphyma/work/polymarket-go/wallets.strategy-review-noise.txt,/Users/murphyma/work/polymarket-go/wallets.strategy-tape-reversal.txt,/Users/murphyma/work/polymarket-go/wallets.strategy-tape-observe.txt`
- Historical alerts excluded by wallet lists: 3
- Historical alerts excluded by current market filter: 1

## Summary

- Alerts: 5
- Marked to current midpoint: 5
- Unmarked: 0
- Win rate incl. midpoint marks: 20.0%
- PnL incl. midpoint marks: $-34.10
- ROI incl. midpoint marks: -68.2%
- Avg price delta: -16.57pp

## Current Policy Performance

- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Current Policy Position-Capped Performance

- Rule: first alert per wallet + asset
- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Effective Tradable Policy Performance

- Rule: current policy modes after removing CUT/PROBATION modes
- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Effective Tradable Position-Capped Performance

- Rule: first alert per wallet + asset after removing CUT/PROBATION modes
- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental OBSERVE Performance

- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted
- Modes: OBSERVE,OBSERVE-BURST
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental OBSERVE Position-Capped Performance

- Rule: first OBSERVE/OBSERVE-BURST alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental OBSERVE-BURST Performance

- Rule: same-wallet split-order sports/esports BUY bursts; not counted in current policy until repeated positive ROI is proven
- Modes: OBSERVE-BURST
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental OBSERVE-BURST Position-Capped Performance

- Rule: first OBSERVE-BURST alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental UNKNOWN-FLOW Performance

- Rule: shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: UNKNOWN-FLOW
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental UNKNOWN-FLOW Position-Capped Performance

- Rule: first UNKNOWN-FLOW alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental SEED-FLOW Performance

- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: SEED-FLOW
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental SEED-FLOW Position-Capped Performance

- Rule: first SEED-FLOW alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental SCORED-FLOW Performance

- Rule: shadow-only scored low-bot leaderboard wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: SCORED-FLOW
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental SCORED-FLOW Position-Capped Performance

- Rule: first SCORED-FLOW alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental INSIDER-SCOUT Performance

- Rule: very large low-bot sports/esports whale BUY alerts; observation only until repeated positive ROI is proven
- Modes: INSIDER-SCOUT
- Alerts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental INSIDER-SCOUT Position-Capped Performance

- Rule: first INSIDER-SCOUT alert per wallet + asset
- Positions: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Experimental CONSENSUS Performance

- Rule: cross-wallet same-asset sports/esports BUY bursts; research/observation only until repeated positive ROI is proven
- Modes: CONSENSUS
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental CONSENSUS Position-Capped Performance

- Rule: first CONSENSUS alert per asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Avg price delta: +37.09pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Mode Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Cut gate: marked>=3 and ROI<=-20.0% or win<=35.0%; severe single loss is probation

| Key | Alerts | Marked | Win | ROI | PnL | Action | Reason |
|---|---:|---:|---:|---:|---:|---|---|
| `CONSENSUS` | 1 | 1 | 100.0% | 59.0% | $+5.90 | COLLECT_POSITIVE | positive but sample below promote gate |

## Wallet Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Cut gate: marked>=3 and ROI<=-20.0% or win<=35.0%; severe single loss is probation

| Key | Alerts | Marked | Win | ROI | PnL | Action | Reason |
|---|---:|---:|---:|---:|---:|---|---|
| `multi:2` | 1 | 1 | 100.0% | 59.0% | $+5.90 | COLLECT_POSITIVE | positive but sample below promote gate |

## By Mode

| Key | Alerts | Marked | Win | ROI | PnL | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|
| `CONSENSUS` | 1 | 1 | 100.0% | 59.0% | $+5.90 | +37.09 |

## By Wallet

| Key | Alerts | Marked | Win | ROI | PnL | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|
| `multi:2` | 1 | 1 | 100.0% | 59.0% | $+5.90 | +37.09 |

## Recent Alerts

| Sent | Mode | Wallet | Notional | Entry | Mid | ROI | Market |
|---|---|---|---:|---:|---:|---:|---|
| 07-09 12:51 | CONSENSUS | `multi:2` | $17277 | 0.629 | 1.000 | 59.0% | Will France win on 2026-07-09? |
| 07-09 05:52 | OBSERVE | `0xd728...9e9b` | $55000 | 0.326 | 0.000 | -100.0% | Will France win the 2026 FIFA World Cup? |
| 07-09 03:42 | FLOW-SCOUT | `0x620d...c67f` | $6000 | 0.195 | 0.000 | -100.0% | Will Argentina win the 2026 FIFA World Cup? |
| 07-09 03:42 | FLOW-SCOUT | `0x620d...c67f` | $14000 | 0.199 | 0.000 | -100.0% | Will Argentina win the 2026 FIFA World Cup? |
| 07-09 02:49 | OBSERVE | `0xf3ce...a57a` | $7675 | 0.479 | 0.000 | -100.0% | Dota 2: Virtus.pro vs 1win - Game 1 Winner |

