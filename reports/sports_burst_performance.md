# Sports Burst Performance

**Generated:** 2026-10-28 00:10 +08

- Tape: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/sports_tape.jsonl`
- Fixed paper stake: $10.00 per burst
- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices and local cache
- Rule: same wallet + same asset, window=15m0s, min total=$5000, min trades=2, min leg=$1000

- Consensus rule: same asset across wallets, min total=$10000, min wallets=2, max bot=60.0; excludes review-noise, reversal-risk, negative-edge, and BOT tiers
- Negative-edge blocked wallets: 30

## Summary

- Bursts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp

## Mode Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| n/a | 0 | 0 | 0.0% | 0.0% | $+0.00 | COLLECT |

## Scope Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| n/a | 0 | 0 | 0.0% | 0.0% | $+0.00 | COLLECT |

## Wallet Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| n/a | 0 | 0 | 0.0% | 0.0% | $+0.00 | COLLECT |

## Consensus Participants

- Rule: attribution table for wallets appearing in CONSENSUS bursts; PnL is attributed to each participant for ranking only

| Wallet | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|---:|
| n/a | 0 | 0 | 0.0% | 0.0% | $+0.00 | $0 | +0.00 |

## Research-Only Consensus Watch

- Rule: lower-threshold cross-wallet same-asset BUY bursts for sample discovery only; not Telegram-pushed and not counted as official CONSENSUS
- Threshold: total>=$5000 wallets>=2
- Watch bursts: 0
- Marked to current midpoint/settlement: 0
- ROI incl. midpoint marks: 0.0%
- Durable watch events: 8
- Durable watch marked: 8
- Durable watch ROI: 6.8%

| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Participants | Market |
|---|---:|---:|---:|---:|---:|---:|---|---|
| n/a | 0 | 0 | $0 | 0.000 | 0.000 | 0.0% |  |  |

## Durable Consensus Event History

- Rule: persisted CONSENSUS events imported from live and shadow alert logs; this survives tape-window rollover
- Events: 9
- Marked to current midpoint/settlement: 9
- Marked samples still needed for promotion review: 0
- Win rate incl. midpoint marks: 55.6%
- PnL incl. midpoint marks: $+3.22
- ROI incl. midpoint marks: 3.6%
- Gate action: COLLECT_POSITIVE

| Last | Category | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---:|---:|---:|---:|---:|---:|---|
| 10-27 00:01 | esports | 7 | 9 | $21915 | 0.909 | 1.000 | 10.0% | Dota 2: Team Yandex vs Team Liquid (BO3) - PARI Universe Group A |
| 10-26 00:02 | esports | 10 | 11 | $16566 | 0.826 | 1.000 | 21.0% | Dota 2: Team Spirit vs Astra Invicta - Game 1 Winner |
| 10-25 00:08 | soccer | 9 | 9 | $62298 | 0.520 | 0.515 | -1.0% | Will Chelsea FC win on 2026-10-24? |
| 10-25 00:07 | esports | 5 | 5 | $18325 | 0.666 | 0.565 | -15.1% | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| 10-23 00:03 | esports | 2 | 2 | $10994 | 0.430 | 0.445 | 3.5% | Dota 2: Team Liquid vs 1win - Game 2 Winner |
| 07-09 12:49 | soccer | 2 | 2 | $17277 | 0.629 | 1.000 | 59.0% | Will France win on 2026-07-09? |
| 07-09 11:40 | basketball | 4 | 4 | $11047 | 0.848 | 0.500 | -41.0% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 11:31 | basketball | 3 | 3 | $10694 | 0.790 | 0.500 | -36.7% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 04:29 | esports | 5 | 5 | $10336 | 0.754 | 1.000 | 32.6% | Dota 2: Virtus.pro vs 1win - Game 2 Winner |

## Durable Consensus Research Wallets

- Rule: persisted address-level attribution from positive CONSENSUS events; wallets stay research-only until repeated marked samples pass promotion gates

| Wallet | Tier | Bot | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x124e...1c58` | D | 58.1 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0xc365...2e53` | C | 38.4 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0xf202...9ea8` | BOT | 62.4 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0x20c6...45ae` | - | 0.0 | 1 | 1 | 100.0% | 59.0% | $+5.90 | $17277 | +37.09 |
| `0x7af7...4c89` | - | 0.0 | 1 | 1 | 100.0% | 59.0% | $+5.90 | $17277 | +37.09 |
| `0x93bd...bbc5` | D | 55.0 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x9679...9315` | D | 56.3 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x23cb...80d0` | D | 57.1 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0x36b3...a92a` | - | 0.0 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0x36d7...c7ed` | - | 0.0 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0x6d3d...51b8` | A | 22.3 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0x92dd...c5e5` | - | 0.0 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0xed09...e314` | - | 0.0 | 1 | 1 | 100.0% | 21.0% | $+2.10 | $16566 | +17.32 |
| `0x7120...7c4e` | D | 57.2 | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |
| `0x08c2...c023` | D | 58.0 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0x6e55...68d0` | D | 54.7 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0x8a36...4836` | - | 0.0 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0x9055...756a` | - | 0.0 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0xce58...0a21` | B | 34.4 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0xd5e4...2c05` | - | 0.0 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0xfa76...e588` | - | 0.0 | 1 | 1 | 100.0% | 10.0% | $+1.00 | $21915 | +9.10 |
| `0x689f...259c` | BOT | 61.5 | 1 | 1 | 100.0% | 3.5% | $+0.35 | $10994 | +1.50 |
| `0xbcc3...ac9c` | D | 53.8 | 1 | 1 | 100.0% | 3.5% | $+0.35 | $10994 | +1.50 |

## Recent Bursts

| Last | Scope | Mode | Actor | Participants | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---|---|---|---:|---:|---:|---:|---:|---:|---|
| n/a |  |  |  |  | 0 | 0 | $0 | 0.000 | 0.000 | 0.0% |  |
