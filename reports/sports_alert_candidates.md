# Sports Alert Candidates

**Generated:** 2026-10-25 00:10 +08

- Recent BUY rows inspected: 26
- Diagnostic window: 6h0m0s
- Alert window: 10m0s
- Consensus alert window: 15m0s
- Currently alertable unsent rows: 1
- Eligible unsent rows in diagnostic window: 1
- Already sent rows in window: 0
- Positive-edge required modes: CANDIDATE,PROBATION
- Observe min notional: $3000
- Observe-burst min notional: $6000
- Unknown-flow min notional: $4000
- Unknown-flow min markets: 2
- Seed-flow min notional: $3000
- Seed-flow min markets: 2
- Scored-flow min notional: $4000
- Scored-flow min markets: 2
- Scored-flow min tier: B
- Scored-flow max bot: 35.0
- Observe min tier: B
- Insider-scout min notional: $25000
- Insider-scout max bot: 35.0
- Edge-hot wallets: 5
- Negative-edge blocked wallets: 30

## Edge-Hot Wallets

| Wallet | Reason |
|---|---|
| `0x2929...1dd0` | edge-hot 100.0% win avg +16.72pp over 4 samples |
| `0x4c9c...e006` | edge-hot 75.0% win avg +16.81pp over 12 samples |
| `0xa75b...772c` | edge-hot 100.0% win avg +30.74pp over 4 samples |
| `0xc6f1...5729` | edge-hot 65.5% win avg +7.80pp over 58 samples |
| `0xcc6e...fa6f` | edge-hot 66.7% win avg +6.70pp over 12 samples |

## Near Edge-Hot Wallets

- Rule: recent large BUYs from edge-hot candidate lists that are not yet eligible; this is diagnostics only and does not loosen Telegram pushes.

| Wallet | List | Tier | Bot | Notional | Samples | Win | AvgPP | 5m | 15m | 1h | Gap | Market |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x1550...5481` | target | C | 43.6 | $4735 | 0 | 0.0% | +0.00 | +0.00 | +0.00 | +0.00 | no edge snapshots | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |

## Negative-Edge Blocked Wallets

| Wallet | Reason |
|---|---|
| `0x076d...8d4c` | 15m edge -6.00pp over 4 samples |
| `0x07a9...32c8` | 1h edge -12.31pp over 2 samples |
| `0x0e24...7014` | 1h edge -72.59pp over 1 samples |
| `0x0f00...73ce` | 15m edge -10.00pp over 2 samples |
| `0x119e...cb14` | 1h edge -36.88pp over 1 samples |
| `0x124e...1c58` | 15m edge -25.28pp over 2 samples |
| `0x18c2...529a` | 1h edge -52.95pp over 1 samples |
| `0x2005...75ea` | 15m edge -2.49pp over 2 samples |
| `0x2039...2e5f` | 1h edge -18.95pp over 1 samples |
| `0x204f...5e14` | 15m edge -1.01pp over 13 samples |
| `0x2c33...0563` | 1h edge -7.50pp over 1 samples |
| `0x34a5...9450` | 15m edge -1.77pp over 4 samples |
| `0x44c4...09cb` | 1h edge -15.87pp over 1 samples |
| `0x4572...83ff` | 1h edge -34.83pp over 1 samples |
| `0x4bff...fc26` | 1h edge -15.00pp over 4 samples |
| `0x5b1d...3721` | 1h edge -17.38pp over 1 samples |
| `0x7124...f0b5` | 15m edge -1.22pp over 2 samples |
| `0x7c1e...bab3` | 1h edge -8.50pp over 2 samples |
| `0x85bb...9ada` | 15m edge -1.74pp over 7 samples |
| `0x9520...fa6e` | 15m edge -1.49pp over 2 samples |
| `0x96d7...3a17` | 1h edge -35.95pp over 1 samples |
| `0xa4fd...748b` | 1h edge -21.51pp over 2 samples |
| `0xb2ed...4418` | 1h edge -19.00pp over 1 samples |
| `0xba8a...b772` | 15m edge -2.63pp over 2 samples |
| `0xbb35...b62a` | 15m edge -2.84pp over 3 samples |
| `0xd8b5...54c4` | 15m edge -24.00pp over 2 samples |
| `0xde24...4ded` | 1h edge -11.42pp over 4 samples |
| `0xe907...cff6` | 15m edge -1.55pp over 26 samples |
| `0xf3ce...a57a` | 1h edge -47.28pp over 1 samples |
| `0xf4ec...c574` | 1h edge -27.79pp over 2 samples |

## Accumulation Bursts

- Rule: same wallet + same asset, strategy mode only, cumulative BUY notional within burst window

| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---|---|
| none |  |  | 0 | $0 | $0 |  |  |  |

## Observe Accumulation Bursts

- Rule: same wallet + same asset, OBSERVE mode only, cumulative BUY notional within burst window; requires scored low-bot wallet unless insider threshold applies separately; consensus_research wallets remain single-wallet blocked and require a CONSENSUS burst

| Status | Mode | Wallet | Trades | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---|---|
| none |  |  | 0 | $0 | $0 |  |  |  |

## Consensus Bursts

- Rule: same asset across wallets, cumulative BUY notional within burst window; excludes review-noise, reversal-risk, negative-edge, and BOT tiers

| Status | Wallets | Trades | Total | VWAP | LastAge | Participants | Reason | Market |
|---|---:|---:|---:|---:|---:|---|---|---|
| consensus | 9 | 9 | $62298 | 0.520 | 2m | `0x074a...f0e0`, `0x0df2...5d51`, `0x2d8a...e818`, `0x479e...91c3`, `0x6bc7...c4a9`, `0xb5db...5800`, `0xbc56...8f28`, `0xcf0c...f2bb`, +1 more | consensus burst in diagnostic window | Will Chelsea FC win on 2026-10-24? |
| consensus | 5 | 5 | $18325 | 0.666 | 3m | `0x00a9...36bb`, `0x1550...5481`, `0x4214...db99`, `0x5dfa...de9c`, `0xc5f0...40ec` | consensus burst in diagnostic window | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |

## Unknown Multi-Market Flow

- Rule: shadow-only unknown wallets buying multiple target markets within the burst window; used to collect edge before any Telegram promotion.

| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---:|---:|---:|---:|---:|---|---|
| none |  | 0 | 0 | $0 | $0 |  |  |  |

## Seed Multi-Market Flow

- Rule: lower-threshold shadow-only unknown wallets buying multiple target markets; used to catch early sports flow before UNKNOWN-FLOW size is reached.

| Status | Wallet | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---:|---:|---:|---:|---:|---|---|
| none |  | 0 | 0 | $0 | $0 |  |  |  |

## Scored Multi-Market Flow

- Rule: shadow-only scored low-bot wallets buying multiple target markets within the burst window; used to find leaderboard whale flow before any Telegram promotion.

| Status | Wallet | Tier | Bot | Trades | Markets | Total | LastBuy | LastAge | Reason | Market |
|---|---|---|---:|---:|---:|---:|---:|---:|---|---|
| none |  |  | 0.0 | 0 | 0 | $0 | $0 |  |  |  |

## Trade Rows

| Status | Mode | Wallet | List | Tier | Bot | Notional | Age | Reason | Market |
|---|---|---|---|---|---:|---:|---:|---|---|
| eligible | OBSERVE | `0x0df2...5d51` | - | - | 0.0 | $26101 | 7m | insider-scout huge whale | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0xcbc6...ec62` | - | D | 50.3 | $294000 | 6m | bot 50.3 > insider max 35.0 | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0xcf0c...f2bb` | - | - | 0.0 | $12804 | 7m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0x074a...f0e0` | - | - | 0.0 | $12804 | 7m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0x4214...db99` | - | - | 0.0 | $4792 | 5m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | EDGE-HOT | `0x1550...5481` | target | C | 43.6 | $4735 | 5m | edge-hot thresholds not met | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x5dfa...de9c` | - | - | 0.0 | $3810 | 3m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0xf9d4...f843` | - | - | 0.0 | $3505 | 7m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0xc5f0...40ec` | - | - | 0.0 | $3120 | 7m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x00a9...36bb` | - | - | 0.0 | $1868 | 8m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0xbc56...8f28` | - | - | 0.0 | $1800 | 6m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | REVERSAL-RISK | `0xe907...cff6` | tape_reversal | BOT | 80.1 | $1619 | 6m | reversal risk disabled | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0xb5db...5800` | - | - | 0.0 | $1464 | 5m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | EDGE-HOT | `0x1550...5481` | target | C | 43.6 | $1444 | 10m | edge-hot thresholds not met | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x479e...91c3` | - | - | 0.0 | $1410 | 8m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0x6bc7...c4a9` | - | - | 0.0 | $1408 | 7m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0x51be...12c2` | - | - | 0.0 | $1244 | 9m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x9697...18e9` | sports_holders_push | D | 54.7 | $1149 | 11m | bot 54.7 > 35.0 | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x2d8a...e818` | - | - | 0.0 | $1000 | 2m | observe wallet unscored | Will Chelsea FC win on 2026-10-24? |
| blocked | OBSERVE | `0x36d7...c7ed` | - | - | 0.0 | $1000 | 3m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x3231...6e75` | - | - | 0.0 | $692 | 6m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x309c...e6fb` | - | - | 0.0 | $690 | 8m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x77d9...06ed` | - | - | 0.0 | $600 | 5m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x77d9...06ed` | - | - | 0.0 | $600 | 5m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x77d9...06ed` | - | - | 0.0 | $600 | 5m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
| blocked | OBSERVE | `0x77d9...06ed` | - | - | 0.0 | $600 | 5m | observe wallet unscored | Dota 2: Team Spirit vs BetBoom Team - Game 2 Winner |
