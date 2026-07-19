# Strategy Lab Report

**Generated:** 2026-07-19 23:11 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 161
- Candidate layers: 11 core + 20 watch + 10 sports + 5 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 48 total
- Live-edge blocked push wallets: 9
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 246
- Aggregate copy ROI: 125.4%
- Aggregate copy PnL: $+3499.49
- Aggregate copy win rate: 82.5%
- Median wallet CopyROI: 129.7%
- Worst included wallet CopyROI: 103.7%
- Open copy cost / closed copy capital: 2.37x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.79 | 117 | 43% | 143.7% | $+545.89 | 32 | 65.6% | 47.1% |
| `0xa75b...772c` | A | 23.93 | 106 | 99% | 108.2% | $+432.72 | 28 | 75.0% | 31.2% |
| `0x2952...f50d` | A | 24.23 | 77 | 52% | 112.3% | $+381.87 | 26 | 88.5% | 19.5% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x90e0...21a2` | A | 22.50 | 51 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 26.66 | 155 | 58% | 109.6% | $+613.96 | 51 | 80.4% | 48.0% |
| `0x89cf...5f47` | B | 25.96 | 57 | 42% | 103.7% | $+383.61 | 36 | 86.1% | 45.8% |
| `0x6f16...5fe7` | B | 25.08 | 41 | 62% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0xb715...c3bb` | B | 25.08 | 57 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xeb8b...6d8a` | B | 27.25 | 25 | 49% | 129.7% | $+103.78 | 8 | 100.0% | 41.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 601
- Aggregate copy ROI: 101.0%
- Aggregate copy PnL: $+8325.62
- Aggregate copy win rate: 89.5%
- Worst included CopyROI: 59.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.6 | 897.4 | 24 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1214 | existing |
| `0x94a8...5204` | B | 100.0 | 15.8 | 797.4 | 24 | 27% | 408.3% | $+244.99 | 3 | 100.0% | $2895 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 29.0 | 740.2 | 38 | 97% | 248.6% | $+223.76 | 7 | 100.0% | $8095 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.4 | 678.8 | 1341 | 100% | 104.2% | $+2301.73 | 98 | 95.9% | $1007 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xcc6e...fa6f` | B | 100.0 | 26.2 | 616.2 | 304 | 84% | 94.8% | $+672.70 | 62 | 90.3% | $2164 | existing,holder,sports_tape |
| `0xde24...4ded` | B | 100.0 | 30.5 | 568.1 | 133 | 100% | 107.0% | $+470.58 | 23 | 100.0% | $2547 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 31.8 | 561.5 | 62 | 50% | 103.7% | $+300.79 | 28 | 96.4% | $2833 | existing |
| `0x4ca8...ecd4` | B | 100.0 | 15.1 | 558.1 | 23 | 88% | 238.4% | $+71.52 | 3 | 100.0% | $620 | existing,holder |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 546.1 | 12 | 19% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x2a35...9015` | B | 100.0 | 33.3 | 532.2 | 57 | 98% | 114.4% | $+205.85 | 15 | 93.3% | $1760 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.5 | 524.8 | 813 | 93% | 59.9% | $+1131.60 | 172 | 89.5% | $297 | existing,holder |
| `0x7992...1fc1` | A | 100.0 | 18.6 | 505.8 | 21 | 21% | 83.5% | $+267.23 | 27 | 77.8% | $2688 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.1 | 488.2 | 137 | 99% | 79.0% | $+260.85 | 32 | 93.8% | $478 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.0 | 485.7 | 142 | 42% | 88.5% | $+274.36 | 24 | 87.5% | $358 | existing,holder |
| `0x84cd...7565` | A | 100.0 | 23.9 | 482.6 | 115 | 73% | 93.7% | $+187.34 | 17 | 76.5% | $432 | existing |
| `0xc117...7410` | A | 84.2 | 25.0 | 475.2 | 32 | 76% | 118.6% | $+154.23 | 10 | 50.0% | $2225 | existing |
| `0x18c2...529a` | B | 100.0 | 28.3 | 470.1 | 374 | 61% | 78.0% | $+272.97 | 31 | 83.9% | $271 | existing,holder |
| `0x5a56...10f5` | B | 100.0 | 32.2 | 466.1 | 84 | 95% | 108.8% | $+261.16 | 10 | 70.0% | $1070 | existing,holder |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xa8b9...775d` | B | 100.0 | 33.7 | 460.8 | 353 | 246 | 93% | 72.4% | $+586.68 | 42 | 57.1% | 67.7% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 29.0 | 444.0 | 162 | 79 | 100% | 71.4% | $+192.74 | 27 | 96.3% | 71.4% | existing | opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.3 | 437.8 | 76 | 76 | 58% | 80.1% | $+184.30 | 21 | 85.7% | 64.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.3 | 421.5 | 87 | 68 | 63% | 85.5% | $+188.08 | 12 | 66.7% | 54.3% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7c73...6ee3` | B | 100.0 | 26.6 | 414.8 | 46 | 45 | 96% | 119.6% | $+71.74 | 6 | 100.0% | 119.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.4 | 404.8 | 366 | 341 | 100% | 64.6% | $+420.17 | 30 | 93.3% | 64.6% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb53e...8617` | B | 100.0 | 34.5 | 404.7 | 58 | 51 | 46% | 88.9% | $+177.87 | 11 | 63.6% | 69.8% | existing | opposite_side_same_market |
| `0xb916...f248` | A | 100.0 | 21.5 | 404.5 | 143 | 113 | 75% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing,holder | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 5
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x56ac...d77e` | B | 100.0 | 27.6 | 267.6 | 104 | 14% | 0.0% | 0 | 568 | $3629 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xe015...6ffd` | C | 100.0 | 39.9 | 214.9 | 1046 | 72% | 1202.7% | 1 | 392 | $585 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 98.2 | 38.8 | 183.1 | 19 | 1% | 0.0% | 0 | 310 | $673 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 90.9 | 44.0 | 147.2 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 42.9 | 147.0 | 0 | 0% | 0.0% | 0 | 226 | $989 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 37.5 | 312.1 | 1287 | 1078 | 100% | 16.1% | 264 | 16.1% | 264 | $476 | existing | burst_trading,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.4 | 289.1 | 1023 | 797 | 85% | 76.6% | 15 | 74.0% | 16 | $1577 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.9 | 285.9 | 1231 | 939 | 100% | 76.6% | 235 | 76.6% | 235 | $328 | existing,holder | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.0 | 243.9 | 566 | 525 | 90% | 88.4% | 96 | 84.4% | 107 | $1707 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.2 | 240.1 | 93 | 89 | 99% | 40.5% | 18 | 40.5% | 18 | $2875 | existing | open_copy_exposure,opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 21.9 | 239.1 | 507 | 227 | 88% | 13.0% | 6 | 13.0% | 10 | $1420 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 23.9 | 236.0 | 253 | 227 | 100% | 37.0% | 30 | 37.0% | 30 | $404 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 228.3 | 73 | 69 | 95% | 58.4% | 20 | 58.4% | 20 | $3677 | existing | opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.6 | 225.3 | 296 | 234 | 89% | 19.5% | 23 | 19.5% | 23 | $766 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x42db...d512` | B | 100.0 | 30.3 | 222.6 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | $1194 | existing,sports_tape | opposite_side_same_market |

## Recent Trade Flow Scout

- Wallets: 0
- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No recent-trade wallets passed the flow filters.

## Sports Tape Hotlist

- Wallets: 1
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x7af7...4c89` | B | 100.0 | 26.4 | 334.2 | 21 | 20 | 58.1% | 3 | 58.1% | 3 | $18593 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |

## Sports Tape Probation

- Wallets: 0
- Rule: scored sports-tape wallets with positive target-copy and sub-45 bot score, but soft flow risks; observation-only until edge windows prove out

| Wallet | Tier | Smart | Bot | ProbationScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|

No sports-tape wallets qualified for probation edge observation.

## Sports Tape Edge-Hot

- Wallets: 3
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 23.9 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 541.4 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.2 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 331.6 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 9
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x34a5...9450` | C | 100.0 | 42.0 | 15m edge -1.77pp over 4 samples | 84.4% | 107 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 34.4 | 15m edge -2.84pp over 3 samples | 104.2% | 98 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 23.9 | 15m edge -24.00pp over 2 samples | 37.0% | 30 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 30.5 | 1h edge -11.42pp over 4 samples | 107.0% | 23 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.7 | 1h edge -15.87pp over 1 samples | 109.6% | 51 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 29.0 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.8 | 1h edge -19.00pp over 1 samples | 143.7% | 32 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.4 | 1h edge -34.83pp over 1 samples | 64.6% | 30 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.3 | 1h edge -52.95pp over 1 samples | 78.0% | 31 | existing,holder | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 26.4 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 74.0 | $7000 | $7000 | 1 | 117.9% | 116 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 52.6 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.6 | 653.4 | 408 | 357 | 35% | 158.3% | 17 | 126.9% | 41 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 74.0 | 623.1 | 209 | 203 | 14% | 128.1% | 27 | 117.9% | 116 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.9 | 583.8 | 106 | 89 | 99% | 108.2% | 28 | 108.2% | 28 | existing,holder,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.2 | 565.7 | 304 | 258 | 84% | 82.6% | 53 | 94.8% | 62 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.9 | 292.4 | 979 | 886 | 74% | 30.2% | 302 | 29.9% | 385 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.7 | 262.2 | 159 | 148 | 19% | 38.0% | 58 | 37.1% | 258 | sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 41.5 | 189.9 | 212 | 190 | 76% | 25.5% | 68 | 20.8% | 90 | holder,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 100.0 | 20.0 | 184.2 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 74.9 | 177.9 | 151 | 148 | 12% | 28.4% | 47 | 44.9% | 305 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 26.4 | 175.8 | 21 | 20 | 95% | 58.1% | 3 | 58.1% | 3 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 15.8 | 112.5 | 7 | 5 | 54% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 74.9 | 10.0 | 96.6 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xb3c1...4837` | prompt | - | B | 100.0 | 21.1 | 81.7 | 198 | 150 | 91% | -0.1% | 21 | 0.3% | 23 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 246 | 125.4% | $+3499.49 | 82.5% | 129.7% | 103.7% | 2.37x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 19 | 457 | 105.8% | $+5438.25 | 83.6% | 108.2% | 63.6% | 2.41x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 3 | 20 | 468 | 104.6% | $+5545.68 | 83.3% | 105.9% | 63.6% | 2.38x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 10 | 220 | 127.2% | $+3117.62 | 81.8% | 132.5% | 103.7% | 2.53x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 18 | 449 | 105.4% | $+5334.47 | 83.3% | 105.9% | 63.6% | 2.40x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 6 | 19 | 460 | 104.3% | $+5441.90 | 83.0% | 103.7% | 63.6% | 2.38x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 17 | 404 | 106.9% | $+4789.15 | 83.7% | 108.2% | 63.6% | 2.56x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 17 | 405 | 104.5% | $+4765.16 | 85.9% | 108.2% | 66.5% | 2.30x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 9 | 10 | 238 | 125.3% | $+3395.71 | 81.9% | 123.8% | 103.7% | 2.35x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 18 | 416 | 103.2% | $+4872.59 | 85.6% | 105.9% | 66.5% | 2.28x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 16 | 396 | 106.5% | $+4685.37 | 83.3% | 105.9% | 63.6% | 2.55x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 12 | 16 | 397 | 104.0% | $+4661.38 | 85.6% | 105.9% | 66.5% | 2.29x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
