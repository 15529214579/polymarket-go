# Sports Alert Performance

**Generated:** 2026-07-09 13:57 +08

- Alert log: `db/strategy_iteration/sports_tape_shadow_alerts.jsonl`
- Extra alert log: `db/strategy_iteration/sports_tape_alerts.jsonl`
- Fixed paper stake: $10.00 per Telegram alert
- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices

- Current policy exclude wallets: `wallets.strategy-quarantine.txt,wallets.strategy-review-noise.txt,wallets.strategy-tape-reversal.txt`
- Historical alerts excluded by wallet lists: 3
- Historical alerts excluded by current market filter: 1

## Summary

- Alerts: 12
- Marked to current midpoint: 12
- Unmarked: 0
- Win rate incl. midpoint marks: 33.3%
- PnL incl. midpoint marks: $-8.68
- ROI incl. midpoint marks: -7.2%
- Avg price delta: -5.90pp

## Current Policy Performance

- Modes: CONSENSUS,OBSERVE-BURST
- Alerts: 5
- Marked to current midpoint: 5
- Unmarked: 0
- Win rate incl. midpoint marks: 40.0%
- PnL incl. midpoint marks: $+1.39
- ROI incl. midpoint marks: 2.8%
- Avg price delta: -4.32pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Current Policy Position-Capped Performance

- Rule: first alert per wallet + asset
- Modes: CONSENSUS,OBSERVE-BURST
- Positions: 4
- Marked to current midpoint: 4
- Unmarked: 0
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $+5.49
- ROI incl. midpoint marks: 13.7%
- Avg price delta: +3.29pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Effective Tradable Policy Performance

- Rule: current policy modes after removing CUT/PROBATION modes
- Modes: OBSERVE-BURST
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Effective Tradable Position-Capped Performance

- Rule: first alert per wallet + asset after removing CUT/PROBATION modes
- Modes: OBSERVE-BURST
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental OBSERVE Performance

- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted
- Modes: OBSERVE,OBSERVE-BURST
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental OBSERVE Position-Capped Performance

- Rule: first OBSERVE/OBSERVE-BURST alert per wallet + asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental OBSERVE-BURST Performance

- Rule: same-wallet split-order sports/esports BUY bursts; not counted in current policy until repeated positive ROI is proven
- Modes: OBSERVE-BURST
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental OBSERVE-BURST Position-Capped Performance

- Rule: first OBSERVE-BURST alert per wallet + asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Avg price delta: +19.00pp
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Experimental UNKNOWN-FLOW Performance

- Rule: shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: UNKNOWN-FLOW
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.07
- ROI incl. midpoint marks: -0.7%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Experimental UNKNOWN-FLOW Position-Capped Performance

- Rule: first UNKNOWN-FLOW alert per wallet + asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.07
- ROI incl. midpoint marks: -0.7%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Experimental SEED-FLOW Performance

- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: SEED-FLOW
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.08
- ROI incl. midpoint marks: -0.8%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Experimental SEED-FLOW Position-Capped Performance

- Rule: first SEED-FLOW alert per wallet + asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.08
- ROI incl. midpoint marks: -0.8%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Experimental SCORED-FLOW Performance

- Rule: shadow-only scored low-bot leaderboard wallets buying multiple target markets; not counted in current policy until repeated positive ROI is proven
- Modes: SCORED-FLOW
- Alerts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.14
- ROI incl. midpoint marks: -1.4%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Experimental SCORED-FLOW Position-Capped Performance

- Rule: first SCORED-FLOW alert per wallet + asset
- Positions: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.14
- ROI incl. midpoint marks: -1.4%
- Avg price delta: -0.50pp
- Gate action: COLLECT
- Gate reason: sample below promote gate

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
- Alerts: 4
- Marked to current midpoint: 4
- Unmarked: 0
- Win rate incl. midpoint marks: 25.0%
- PnL incl. midpoint marks: $-4.74
- ROI incl. midpoint marks: -11.8%
- Avg price delta: -10.15pp
- Gate action: CUT
- Gate reason: negative marked sample

## Experimental CONSENSUS Position-Capped Performance

- Rule: first CONSENSUS alert per asset
- Positions: 3
- Marked to current midpoint: 3
- Unmarked: 0
- Win rate incl. midpoint marks: 33.3%
- PnL incl. midpoint marks: $-0.64
- ROI incl. midpoint marks: -2.1%
- Avg price delta: -1.95pp
- Gate action: CUT
- Gate reason: negative marked sample

## Mode Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Cut gate: marked>=3 and ROI<=-20.0% or win<=35.0%; severe single loss is probation

| Key | Alerts | Marked | Win | ROI | PnL | Action | Reason |
|---|---:|---:|---:|---:|---:|---|---|
| `CONSENSUS` | 4 | 4 | 25.0% | -11.8% | $-4.74 | CUT | negative marked sample |
| `OBSERVE-BURST` | 1 | 1 | 100.0% | 61.3% | $+6.13 | COLLECT_POSITIVE | positive but sample below promote gate |
| `UNKNOWN-FLOW` | 1 | 1 | 0.0% | -0.7% | $-0.07 | COLLECT | sample below promote gate |
| `SEED-FLOW` | 1 | 1 | 0.0% | -0.8% | $-0.08 | COLLECT | sample below promote gate |
| `SCORED-FLOW` | 1 | 1 | 0.0% | -1.4% | $-0.14 | COLLECT | sample below promote gate |

## Wallet Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Cut gate: marked>=3 and ROI<=-20.0% or win<=35.0%; severe single loss is probation

| Key | Alerts | Marked | Win | ROI | PnL | Action | Reason |
|---|---:|---:|---:|---:|---:|---|---|
| `0x6027...4320` | 1 | 1 | 100.0% | 61.3% | $+6.13 | COLLECT_POSITIVE | positive but sample below promote gate |
| `multi:5` | 1 | 1 | 100.0% | 32.6% | $+3.26 | COLLECT_POSITIVE | positive but sample below promote gate |
| `0xb748...dbdd` | 1 | 1 | 0.0% | -0.7% | $-0.07 | COLLECT | sample below promote gate |
| `0x2393...1a8d` | 1 | 1 | 0.0% | -0.8% | $-0.08 | COLLECT | sample below promote gate |
| `0xb3c1...4837` | 1 | 1 | 0.0% | -1.4% | $-0.14 | COLLECT | sample below promote gate |
| `multi:2` | 1 | 1 | 0.0% | -2.2% | $-0.22 | COLLECT | sample below promote gate |
| `multi:3` | 1 | 1 | 0.0% | -36.7% | $-3.67 | COLLECT | sample below promote gate |
| `multi:4` | 1 | 1 | 0.0% | -41.0% | $-4.10 | COLLECT | sample below promote gate |

## By Mode

| Key | Alerts | Marked | Win | ROI | PnL | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|
| `OBSERVE-BURST` | 1 | 1 | 100.0% | 61.3% | $+6.13 | +19.00 |
| `UNKNOWN-FLOW` | 1 | 1 | 0.0% | -0.7% | $-0.07 | -0.50 |
| `SEED-FLOW` | 1 | 1 | 0.0% | -0.8% | $-0.08 | -0.50 |
| `SCORED-FLOW` | 1 | 1 | 0.0% | -1.4% | $-0.14 | -0.50 |
| `CONSENSUS` | 4 | 4 | 25.0% | -11.8% | $-4.74 | -10.15 |

## By Wallet

| Key | Alerts | Marked | Win | ROI | PnL | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|
| `0x6027...4320` | 1 | 1 | 100.0% | 61.3% | $+6.13 | +19.00 |
| `multi:5` | 1 | 1 | 100.0% | 32.6% | $+3.26 | +24.58 |
| `0xb748...dbdd` | 1 | 1 | 0.0% | -0.7% | $-0.07 | -0.50 |
| `0x2393...1a8d` | 1 | 1 | 0.0% | -0.8% | $-0.08 | -0.50 |
| `0xb3c1...4837` | 1 | 1 | 0.0% | -1.4% | $-0.14 | -0.50 |
| `multi:2` | 1 | 1 | 0.0% | -2.2% | $-0.22 | -1.41 |
| `multi:3` | 1 | 1 | 0.0% | -36.7% | $-3.67 | -29.02 |
| `multi:4` | 1 | 1 | 0.0% | -41.0% | $-4.10 | -34.77 |

## Recent Alerts

| Sent | Mode | Wallet | Notional | Entry | Mid | ROI | Market |
|---|---|---|---:|---:|---:|---:|---|
| 07-09 13:56 | SEED-FLOW | `0x2393...1a8d` | $3000 | 0.620 | 0.615 | -0.8% | Will France win on 2026-07-09? |
| 07-09 13:21 | SCORED-FLOW | `0xb3c1...4837` | $7781 | 0.350 | 0.345 | -1.4% | Norway vs. England: Team to Advance |
| 07-09 12:51 | CONSENSUS | `multi:2` | $17277 | 0.629 | 0.615 | -2.2% | Will France win on 2026-07-09? |
| 07-09 12:38 | UNKNOWN-FLOW | `0xb748...dbdd` | $6000 | 0.740 | 0.735 | -0.7% | Argentina vs. Switzerland: Team to Advance |
| 07-09 11:50 | CONSENSUS | `multi:4` | $11047 | 0.848 | 0.500 | -41.0% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 11:41 | CONSENSUS | `multi:3` | $10694 | 0.790 | 0.500 | -36.7% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 09:50 | CONSENSUS | `multi:5` | $10336 | 0.754 | 1.000 | 32.6% | Dota 2: Virtus.pro vs 1win - Game 2 Winner |
| 07-09 09:38 | OBSERVE-BURST | `0x6027...4320` | $9084 | 0.310 | 0.500 | 61.3% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 05:52 | OBSERVE | `0xd728...9e9b` | $55000 | 0.326 | 0.322 | -1.4% | Will France win the 2026 FIFA World Cup? |
| 07-09 03:42 | FLOW-SCOUT | `0x620d...c67f` | $6000 | 0.195 | 0.201 | 2.8% | Will Argentina win the 2026 FIFA World Cup? |
| 07-09 03:42 | FLOW-SCOUT | `0x620d...c67f` | $14000 | 0.199 | 0.201 | 0.8% | Will Argentina win the 2026 FIFA World Cup? |
| 07-09 02:49 | OBSERVE | `0xf3ce...a57a` | $7675 | 0.479 | 0.000 | -100.0% | Dota 2: Virtus.pro vs 1win - Game 1 Winner |

