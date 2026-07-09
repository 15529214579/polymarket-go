# Sports Burst Performance

**Generated:** 2026-07-09 11:44 +08

- Tape: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/sports_tape.jsonl`
- Fixed paper stake: $10.00 per burst
- Mark source: live CLOB midpoint; closed markets fall back to Gamma settlement outcome prices and local cache
- Rule: same wallet + same asset, window=15m0s, min total=$5000, min trades=2, min leg=$1000

- Consensus rule: same asset across wallets, min total=$10000, min wallets=2, max bot=60.0; excludes review-noise, reversal-risk, negative-edge, and BOT tiers
- Negative-edge blocked wallets: 8

## Summary

- Bursts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+1.38
- ROI incl. midpoint marks: 13.8%
- Avg price delta: +11.68pp

## Mode Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| `CONSENSUS` | 1 | 1 | 100.0% | 13.8% | $+1.38 | COLLECT_POSITIVE |

## Scope Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| `consensus` | 1 | 1 | 100.0% | 13.8% | $+1.38 | COLLECT_POSITIVE |

## Wallet Gates

- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Collect-positive gate: ROI>0 with sample below promote gate

| Key | Bursts | Marked | Win | ROI | PnL | Action |
|---|---:|---:|---:|---:|---:|---|
| `multi:4` | 1 | 1 | 100.0% | 13.8% | $+1.38 | COLLECT_POSITIVE |

## Consensus Participants

- Rule: attribution table for wallets appearing in CONSENSUS bursts; PnL is attributed to each participant for ranking only

| Wallet | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |
|---|---:|---:|---:|---:|---:|---:|---:|
| `0x124e...1c58` | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |
| `0x7120...7c4e` | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |
| `0xc365...2e53` | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |
| `0xf202...9ea8` | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |

## Research-Only Consensus Watch

- Rule: lower-threshold cross-wallet same-asset BUY bursts for sample discovery only; not Telegram-pushed and not counted as official CONSENSUS
- Threshold: total>=$5000 wallets>=2
- Watch bursts: 0
- Marked to current midpoint/settlement: 0
- ROI incl. midpoint marks: 0.0%
- Durable watch events: 2
- Durable watch marked: 2
- Durable watch ROI: 4.5%

| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Participants | Market |
|---|---:|---:|---:|---:|---:|---:|---|---|
| n/a | 0 | 0 | $0 | 0.000 | 0.000 | 0.0% |  |  |

## Durable Consensus Event History

- Rule: persisted CONSENSUS events imported from live and shadow alert logs; this survives tape-window rollover
- Events: 3
- Marked to current midpoint/settlement: 3
- Marked samples still needed for promotion review: 2
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.84
- ROI incl. midpoint marks: 22.8%
- Gate action: COLLECT_POSITIVE

| Last | Category | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---:|---:|---:|---:|---:|---:|---|
| 07-09 11:40 | basketball | 4 | 6 | $11047 | 0.848 | 0.965 | 13.8% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 11:31 | basketball | 3 | 3 | $10694 | 0.790 | 0.964 | 22.0% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 04:29 | esports | 5 | 5 | $10336 | 0.754 | 1.000 | 32.6% | Dota 2: Virtus.pro vs 1win - Game 2 Winner |

## Durable Consensus Research Wallets

- Rule: persisted address-level attribution from positive CONSENSUS events; wallets stay research-only until repeated marked samples pass promotion gates

| Wallet | Tier | Bot | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x124e...1c58` | C | 40.2 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0xc365...2e53` | C | 42.1 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0xf202...9ea8` | D | 51.5 | 2 | 2 | 100.0% | 17.9% | $+3.58 | $21741 | +14.53 |
| `0x3833...33d1` | D | 48.0 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x67bb...b6ef` | D | 56.6 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x87da...53d1` | D | 57.3 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x93bd...bbc5` | D | 55.0 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x9679...9315` | D | 56.3 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x7120...7c4e` | D | 53.5 | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |

## Recent Bursts

| Last | Scope | Mode | Actor | Participants | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---|---|---|---:|---:|---:|---:|---:|---:|---|
| 4m | consensus | CONSENSUS | `multi:4` | `0x124e...1c58`, `0x7120...7c4e`, `0xc365...2e53`, `0xf202...9ea8` | 4 | 6 | $11047 | 0.848 | 0.965 | 13.8% | Indiana Fever vs. Los Angeles Sparks |

