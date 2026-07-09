# Sports Alert Candidates

**Generated:** 2026-07-09 13:59 +08

- Recent BUY rows inspected: 50
- Diagnostic window: 6h0m0s
- Alert window: 10m0s
- Consensus alert window: 15m0s
- Currently alertable unsent rows: 0
- Eligible unsent rows in diagnostic window: 3
- Already sent rows in window: 0
- Positive-edge required modes: CANDIDATE,PROBATION
- Observe min notional: $5000
- Observe-burst min notional: $8000
- Unknown-flow min notional: $6000
- Unknown-flow min markets: 2
- Seed-flow min notional: $3000
- Seed-flow min markets: 2
- Scored-flow min notional: $6000
- Scored-flow min markets: 2
- Scored-flow min tier: B
- Scored-flow max bot: 35.0
- Observe min tier: B
- Insider-scout min notional: $25000
- Insider-scout max bot: 35.0
- Edge-hot wallets: 3
- Negative-edge blocked wallets: 8

## Edge-Hot Wallets

| Wallet | Reason |
|---|---|
| `0x124e...1c58` | edge-hot 100.0% win avg +10.55pp over 4 samples |
| `0xb36f...53d0` | edge-hot 100.0% win avg +26.13pp over 4 samples |
| `0xe872...819a` | edge-hot 100.0% win avg +21.99pp over 4 samples |

## Near Edge-Hot Wallets

- Rule: recent large BUYs from edge-hot candidate lists that are not yet eligible; this is diagnostics only and does not loosen Telegram pushes.

| Wallet | List | Tier | Bot | Notional | Samples | Win | AvgPP | 5m | 15m | 1h | Gap | Market |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x42db...d512` | watch | B | 30.3 | $1000 | 4 | 0.0% | -0.50 | -0.50 | -0.50 | -0.50 | avg -0.50pp < +2.00pp | Will Spain win on 2026-07-10? |
| `0xfe18...3980` | sports | A | 24.9 | $952 | 2 | 0.0% | -0.40 | +0.00 | +0.00 | -0.40 | avg -0.40pp < +2.00pp | Will Spain win on 2026-07-10? |
| `0x2393...1a8d` | watch | A | 22.4 | $2000 | 2 | 0.0% | -0.50 | +0.00 | -0.50 | +0.00 | avg -0.50pp < +2.00pp | Will France win on 2026-07-09? |

## Negative-Edge Blocked Wallets

| Wallet | Reason |
|---|---|
| `0x2005...75ea` | 15m edge -2.49pp over 2 samples |
| `0x2039...2e5f` | 1h edge -18.95pp over 1 samples |
| `0x2c33...0563` | 1h edge -7.50pp over 1 samples |
| `0x4572...83ff` | 1h edge -34.83pp over 1 samples |
| `0x4bff...fc26` | 1h edge -15.00pp over 4 samples |
| `0x96d7...3a17` | 1h edge -35.95pp over 1 samples |
| `0xd8b5...54c4` | 15m edge -24.00pp over 2 samples |
| `0xf3ce...a57a` | 1h edge -47.28pp over 1 samples |

## Accumulation Bursts

- Rule: same wallet + same asset, strategy mode only, cumulative BUY notional within burst window

| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---|---|
| none |  |  | 0 | $0 | $0 |  |  |  |

## Observe Accumulation Bursts

- Rule: same wallet + same asset, OBSERVE mode only, cumulative BUY notional within burst window; requires scored low-bot wallet unless insider threshold applies separately; consensus_research wallets remain single-wallet blocked and require a CONSENSUS burst

| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---|---|
| stale | OBSERVE-BURST | `0x6027...4320` | 2 | $9084 | $4466 | 4.9h | stale observe burst: last age 4.9h > alert window 10m0s | Indiana Fever vs. Los Angeles Sparks |

## Consensus Bursts

- Rule: same asset across wallets, cumulative BUY notional within burst window; excludes review-noise, reversal-risk, negative-edge, and BOT tiers

| Status | Wallets | Trades | Total | VWAP | LastAge | Participants | Reason | Market |
|---|---:|---:|---:|---:|---:|---|---|---|
| stale | 4 | 6 | $11047 | 0.848 | 2.3h | `0x124e...1c58`, `0x7120...7c4e`, `0xc365...2e53`, `0xf202...9ea8` | stale consensus: last age 2.3h > alert window 15m0s | Indiana Fever vs. Los Angeles Sparks |

## Unknown Multi-Market Flow

- Rule: shadow-only unknown wallets buying multiple target markets within the burst window; used to collect edge before any Telegram promotion.

| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---:|---:|---:|---:|---:|---|---|
| none |  | 0 | 0 | $0 | $0 |  |  |  |

## Seed Multi-Market Flow

- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; used to catch early sports flow before UNKNOWN-FLOW size is reached.

| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---:|---:|---:|---:|---:|---|---|
| stale | `0x2393...1a8d` | 2 | 2 | $3000 | $2000 | 25m | stale seed flow: last age 25m > alert window 10m0s | Will France win on 2026-07-09? |

## Scored Multi-Market Flow

- Rule: shadow-only scored low-bot wallets buying multiple target markets within the burst window; used to find leaderboard whale flow before any Telegram promotion.

| Status | Wallet | Tier | Bot | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---:|---:|---|---|
| stale | `0xb3c1...4837` | B | 21.4 | 3 | 3 | $7781 | $1681 | 1.1h | stale scored flow: last age 1.1h > alert window 10m0s | Norway vs. England: Team to Advance |

## Trade Rows

| Status | Mode | Wallet | List | Tier | Bot | Notional | Age | Reason | Market |
|---|---|---|---|---|---:|---:|---:|---|---|
| stale | OBSERVE | `0x9520...fa6e` | tape_observe | B | 33.8 | $14778 | 2.7h | stale: age 2.7h > alert window 10m0s | Will Argentina win on 2026-07-11? |
| stale | EDGE-HOT | `0x124e...1c58` | tape_edgehot | C | 40.5 | $4366 | 2.5h | stale: age 2.5h > alert window 10m0s | Indiana Fever vs. Los Angeles Sparks |
| stale | EDGE-HOT | `0x124e...1c58` | tape_edgehot | C | 40.5 | $1157 | 3.7h | stale: age 3.7h > alert window 10m0s | Indiana Fever vs. Los Angeles Sparks |
| blocked | REVIEW-NOISE | `0xe0cc...aafc` | review_noise | D | 21.7 | $18862 | 43m | review-noise excluded | France vs. Morocco: Team to Advance |
| blocked | REVIEW-NOISE | `0x7af7...4c89` | review_noise | D | 23.8 | $15777 | 1.3h | review-noise excluded | Will France win on 2026-07-09? |
| blocked | REVIEW-NOISE | `0xa4fd...748b` | review_noise | BOT | 62.0 | $7000 | 40m | review-noise excluded | France vs. Morocco: Team to Advance |
| blocked | REVIEW-NOISE | `0x204f...5e14` | review_noise | D | 52.9 | $5209 | 3.7h | review-noise excluded | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0x6027...4320` | tape_observe | B | 34.3 | $4618 | 4.9h | notional $4618 < observe min $5000 | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x6027...4320` | tape_observe | B | 34.3 | $4466 | 4.9h | notional $4466 < observe min $5000 | Indiana Fever vs. Los Angeles Sparks |
| blocked | REVERSAL-RISK | `0xe907...cff6` | tape_reversal | BOT | 80.6 | $3883 | 4.1h | reversal risk disabled | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0x3eec...6b3d` | - | D | 57.7 | $3500 | 2.3h | bot 57.7 > 35.0 | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0xb3c1...4837` | - | B | 21.4 | $3100 | 1.1h | notional $3100 < observe min $5000 | Spain vs. Belgium: Team to Advance |
| blocked | OBSERVE | `0xb3c1...4837` | - | B | 21.4 | $3000 | 1.1h | notional $3000 < observe min $5000 | France vs. Morocco: Team to Advance |
| blocked | OBSERVE | `0xb748...dbdd` | - | C | 39.4 | $3000 | 1.5h | bot 39.4 > 35.0 | Argentina vs. Switzerland: Team to Advance |
| blocked | OBSERVE | `0xb748...dbdd` | - | C | 39.4 | $3000 | 1.6h | bot 39.4 > 35.0 | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0xbe7e...6d00` | - | BOT | 64.1 | $3000 | 2.9h | BOT tier | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x6982...a165` | - | BOT | 84.6 | $2500 | 2.7h | BOT tier | Indiana Fever vs. Los Angeles Sparks |
| blocked | EDGE-HOT | `0x2393...1a8d` | watch | A | 22.4 | $2000 | 25m | edge-hot thresholds not met | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0xc365...2e53` | consensus_research | C | 42.1 | $1989 | 2.5h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x7ea5...de7b` | - | C | 36.1 | $1987 | 3.5h | bot 36.1 > 35.0 | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0xc365...2e53` | consensus_research | C | 42.1 | $1978 | 2.7h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x5268...135d` | leaderboard_watch | C | 43.3 | $1752 | 5.1h | bot 43.3 > 35.0 | Minnesota Lynx vs. Connecticut Sun |
| blocked | OBSERVE | `0xb3c1...4837` | - | B | 21.4 | $1681 | 1.1h | notional $1681 < observe min $5000 | Norway vs. England: Team to Advance |
| blocked | OBSERVE | `0x6a34...c216` | - | - | 0.0 | $1640 | 11m | observe wallet unscored | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0x20c6...45ae` | - | C | 44.3 | $1500 | 1.2h | bot 44.3 > 35.0 | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0xbd0c...c1f2` | - | D | 29.8 | $1300 | 2.9h | tier D below observe min B | Indiana Fever vs. Phoenix Mercury |
| blocked | OBSERVE | `0xf202...9ea8` | consensus_research | D | 51.5 | $1284 | 2.3h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x43cf...88fa` | - | BOT | 63.0 | $1260 | 2.0h | BOT tier | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0xf202...9ea8` | consensus_research | D | 51.5 | $1250 | 2.6h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | REVIEW-NOISE | `0x204f...5e14` | review_noise | D | 52.9 | $1201 | 3.7h | review-noise excluded | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0x7120...7c4e` | consensus_research | D | 53.5 | $1151 | 2.4h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0x7120...7c4e` | consensus_research | D | 53.5 | $1146 | 2.4h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0xf202...9ea8` | consensus_research | D | 51.5 | $1111 | 2.5h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0xa9a6...0009` | - | D | 57.9 | $1081 | 5.1h | bot 57.9 > 35.0 | Minnesota Lynx vs. Connecticut Sun |
| blocked | OBSERVE | `0xc079...7bd5` | - | B | 9.2 | $1034 | 2.4h | notional $1034 < observe min $5000 | France vs. Morocco: Team to Advance |
| blocked | OBSERVE | `0xaa9e...e588` | - | D | 21.7 | $1008 | 4.7h | tier D below observe min B | Will Argentina win on 2026-07-11? |
| blocked | EDGE-HOT | `0x2393...1a8d` | watch | A | 22.4 | $1000 | 25m | edge-hot thresholds not met | France vs. Morocco: Team to Advance |
| blocked | OBSERVE | `0x473e...cfe2` | - | D | 47.8 | $1000 | 40m | bot 47.8 > 35.0 | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0x473e...cfe2` | - | D | 47.8 | $1000 | 40m | bot 47.8 > 35.0 | Will Argentina win on 2026-07-11? |
| blocked | OBSERVE | `0x9c98...f4f4` | - | D | 50.3 | $1000 | 1.7h | bot 50.3 > 35.0 | France vs. Morocco: Team to Advance |
| blocked | OBSERVE | `0xaeaa...ed6d` | - | BOT | 65.4 | $1000 | 3.6h | BOT tier | France vs. Morocco: Team to Advance |
| blocked | EDGE-HOT | `0x42db...d512` | watch | B | 30.3 | $1000 | 4.7h | edge-hot thresholds not met | Will Spain win on 2026-07-10? |
| blocked | REVERSAL-RISK | `0xe907...cff6` | tape_reversal | BOT | 80.6 | $971 | 4.3h | reversal risk disabled | Will France win on 2026-07-09? |
| blocked | OBSERVE | `0x974b...cc9a` | - | D | 22.9 | $961 | 3.2h | tier D below observe min B | Will France win on 2026-07-09? |
| blocked | EDGE-HOT | `0xfe18...3980` | sports | A | 24.9 | $952 | 2.1h | edge-hot thresholds not met | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0x6d3c...d294` | - | BOT | 66.8 | $920 | 5.1h | BOT tier | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0x6d3c...d294` | - | BOT | 66.8 | $920 | 5.1h | BOT tier | Will Spain win on 2026-07-10? |
| blocked | OBSERVE | `0x7120...7c4e` | consensus_research | D | 53.5 | $858 | 2.3h | consensus research wallet requires consensus burst | Indiana Fever vs. Los Angeles Sparks |
| blocked | OBSERVE | `0xbeea...0c83` | - | B | 18.2 | $846 | 3.0h | notional $846 < observe min $5000 | France vs. Morocco: Team to Advance |
| blocked | OBSERVE | `0xec3c...5dc4` | - | BOT | 60.4 | $819 | 3.6h | BOT tier | France vs. Morocco: Team to Advance |
