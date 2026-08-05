# Sports Burst Performance

**Generated:** 2026-08-05 00:10 +08

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
- Watch bursts: 1
- Marked to current midpoint/settlement: 1
- ROI incl. midpoint marks: 20.3%
- Durable watch events: 4
- Durable watch marked: 4
- Durable watch ROI: 9.1%

| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Participants | Market |
|---|---:|---:|---:|---:|---:|---:|---|---|
| 08-05 00:06 | 3 | 3 | $5020 | 0.769 | 0.925 | 20.3% | `0x0c99...0c88`, `0x1b1f...d6d5`, `0x51be...12c2` | Dota 2: Team Falcons vs Team Liquid - Game 1 Winner |

### Durable Watch History

| Last | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---:|---:|---:|---:|---:|---:|---|
| 08-05 00:06 | 3 | 3 | $5020 | 0.769 | 0.925 | 20.3% | Dota 2: Team Falcons vs Team Liquid - Game 1 Winner |
| 08-02 00:07 | 3 | 6 | $8586 | 0.694 | 0.745 | 7.3% | Dota 2: GLYPH vs Zero Tenacity - Game 1 Winner |
| 07-09 11:28 | 3 | 4 | $8705 | 0.781 | 0.843 | 7.9% | Indiana Fever vs. Los Angeles Sparks |
| 07-09 10:29 | 2 | 3 | $8397 | 0.400 | 0.404 | 1.0% | Will Spain win on 2026-07-10? |

### Durable Watch Participants

| Wallet | Tier | Bot | Signals | Marked | Win | ROI | PnL | TotalNotional | AvgDeltaPP |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0c99...0c88` | - | 0.0 | 1 | 1 | 100.0% | 20.3% | $+2.03 | $5020 | +15.63 |
| `0x1b1f...d6d5` | C | 43.5 | 1 | 1 | 100.0% | 20.3% | $+2.03 | $5020 | +15.63 |
| `0x51be...12c2` | - | 0.0 | 1 | 1 | 100.0% | 20.3% | $+2.03 | $5020 | +15.63 |
| `0x124e...1c58` | C | 40.5 | 1 | 1 | 100.0% | 7.9% | $+0.79 | $8705 | +6.19 |
| `0xc365...2e53` | - | 0.0 | 1 | 1 | 100.0% | 7.9% | $+0.79 | $8705 | +6.19 |
| `0xf202...9ea8` | - | 0.0 | 1 | 1 | 100.0% | 7.9% | $+0.79 | $8705 | +6.19 |
| `0x4aec...4ce9` | D | 58.8 | 1 | 1 | 100.0% | 7.3% | $+0.73 | $8586 | +5.07 |
| `0xdf27...026f` | - | 0.0 | 1 | 1 | 100.0% | 7.3% | $+0.73 | $8586 | +5.07 |
| `0xe16d...5e30` | D | 45.7 | 1 | 1 | 100.0% | 7.3% | $+0.73 | $8586 | +5.07 |
| `0x204f...5e14` | D | 48.9 | 1 | 1 | 100.0% | 1.0% | $+0.10 | $8397 | +0.40 |
| `0x7ea5...de7b` | - | 0.0 | 1 | 1 | 100.0% | 1.0% | $+0.10 | $8397 | +0.40 |

## Durable Consensus Event History

- Rule: persisted CONSENSUS events imported from live and shadow alert logs; this survives tape-window rollover
- Events: 4
- Marked to current midpoint/settlement: 4
- Marked samples still needed for promotion review: 1
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $+1.38
- ROI incl. midpoint marks: 3.5%
- Gate action: COLLECT_POSITIVE

| Last | Category | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---:|---:|---:|---:|---:|---:|---|
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
| `0x7af7...4c89` | B | 30.5 | 1 | 1 | 100.0% | 59.0% | $+5.90 | $17277 | +37.09 |
| `0x93bd...bbc5` | D | 55.0 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x9679...9315` | D | 56.3 | 1 | 1 | 100.0% | 32.6% | $+3.26 | $10336 | +24.58 |
| `0x7120...7c4e` | D | 57.2 | 1 | 1 | 100.0% | 13.8% | $+1.38 | $11047 | +11.68 |

## Recent Bursts

| Last | Scope | Mode | Actor | Participants | Wallets | Trades | Total | VWAP | Mid | ROI | Market |
|---|---|---|---|---|---:|---:|---:|---:|---:|---:|---|
| n/a |  |  |  |  | 0 | 0 | $0 | 0.000 | 0.000 | 0.0% |  |
