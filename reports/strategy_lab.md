# Strategy Lab Report

**Generated:** 2026-07-18 23:54 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 199
- Candidate layers: 18 core + 20 watch + 10 sports + 3 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 50 total
- Live-edge blocked push wallets: 12
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 18
- Aggregate closed copy trades: 436
- Aggregate copy ROI: 106.7%
- Aggregate copy PnL: $+5205.55
- Aggregate copy win rate: 83.7%
- Median wallet CopyROI: 106.1%
- Worst included wallet CopyROI: 63.6%
- Open copy cost / closed copy capital: 2.40x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.82 | 116 | 43% | 149.7% | $+553.94 | 31 | 67.7% | 45.5% |
| `0xa75b...772c` | A | 24.48 | 77 | 100% | 116.7% | $+385.06 | 25 | 76.0% | 29.9% |
| `0x2952...f50d` | A | 24.29 | 77 | 52% | 112.3% | $+381.87 | 26 | 88.5% | 19.5% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x7992...1fc1` | A | 18.97 | 21 | 23% | 83.5% | $+267.23 | 27 | 77.8% | 19.6% |
| `0x6f16...5fe7` | A | 23.95 | 37 | 60% | 201.7% | $+201.68 | 10 | 80.0% | 31.7% |
| `0x84cd...7565` | A | 23.92 | 115 | 73% | 93.7% | $+187.34 | 17 | 76.5% | 41.9% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x7fcf...80ac` | A | 11.62 | 32 | 41% | 87.2% | $+104.63 | 9 | 66.7% | 38.5% |
| `0xcc6e...fa6f` | B | 26.48 | 298 | 84% | 98.1% | $+676.72 | 60 | 91.7% | 45.6% |
| `0x44c4...09cb` | B | 26.73 | 158 | 59% | 109.3% | $+622.92 | 52 | 80.8% | 48.0% |
| `0x89cf...5f47` | B | 26.02 | 57 | 43% | 103.7% | $+383.61 | 36 | 86.1% | 45.8% |
| `0x18c2...529a` | B | 28.25 | 352 | 59% | 81.6% | $+269.18 | 30 | 83.3% | 46.3% |
| `0x21cc...54bc` | B | 28.98 | 162 | 100% | 71.4% | $+192.74 | 27 | 96.3% | 35.4% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x0931...e78e` | B | 25.99 | 117 | 66% | 63.6% | $+127.20 | 20 | 65.0% | 26.8% |
| `0xec56...1f87` | B | 27.56 | 26 | 30% | 66.5% | $+112.96 | 17 | 100.0% | 22.6% |
| `0xeb8b...6d8a` | B | 27.25 | 25 | 49% | 129.7% | $+103.78 | 8 | 100.0% | 41.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 563
- Aggregate copy ROI: 98.2%
- Aggregate copy PnL: $+7956.78
- Aggregate copy win rate: 86.3%
- Worst included CopyROI: 56.7%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.8 | 896.7 | 24 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1119 | existing |
| `0x94a8...5204` | B | 100.0 | 15.6 | 797.5 | 22 | 27% | 408.3% | $+244.99 | 3 | 100.0% | $2855 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 29.3 | 740.9 | 34 | 97% | 248.6% | $+223.76 | 7 | 100.0% | $8130 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.4 | 679.9 | 1342 | 100% | 103.5% | $+2339.25 | 99 | 96.0% | $989 | existing,holder |
| `0xc367...3066` | B | 100.0 | 31.2 | 676.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing,holder |
| `0xde24...4ded` | B | 100.0 | 30.6 | 576.5 | 118 | 100% | 114.6% | $+435.61 | 20 | 100.0% | $2508 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 561.9 | 61 | 50% | 103.7% | $+300.79 | 28 | 96.4% | $2840 | existing |
| `0x4ca8...ecd4` | B | 100.0 | 14.4 | 555.8 | 22 | 88% | 238.4% | $+71.52 | 3 | 100.0% | $604 | existing |
| `0x2a35...9015` | B | 100.0 | 32.3 | 545.9 | 49 | 98% | 130.0% | $+181.95 | 12 | 91.7% | $1754 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.4 | 542.7 | 11 | 18% | 129.2% | $+142.10 | 11 | 90.9% | $2412 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.5 | 525.7 | 789 | 92% | 60.7% | $+1092.35 | 164 | 89.0% | $294 | existing,holder |
| `0xd5b1...1b71` | B | 100.0 | 31.1 | 488.2 | 137 | 99% | 79.0% | $+260.85 | 32 | 93.8% | $478 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 482.1 | 134 | 41% | 89.3% | $+268.06 | 23 | 87.0% | $353 | existing |
| `0xc117...7410` | B | 82.8 | 25.5 | 479.5 | 32 | 80% | 133.9% | $+147.33 | 8 | 50.0% | $2138 | existing |
| `0x5a56...10f5` | B | 100.0 | 31.5 | 471.0 | 72 | 99% | 115.0% | $+230.02 | 9 | 66.7% | $1085 | existing,holder |
| `0xa8b9...775d` | B | 100.0 | 33.7 | 456.4 | 352 | 93% | 67.7% | $+575.12 | 46 | 52.2% | $1208 | existing |
| `0x7c73...6ee3` | B | 100.0 | 26.6 | 453.3 | 46 | 96% | 119.6% | $+71.74 | 6 | 100.0% | $5556 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.0 | 452.9 | 72 | 36% | 56.7% | $+328.91 | 56 | 80.4% | $1233 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.5 | 429.2 | 83 | 64 | 62% | 90.5% | $+189.96 | 11 | 72.7% | 56.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.2 | 416.5 | 72 | 72 | 56% | 75.8% | $+159.16 | 19 | 89.5% | 60.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.6 | 405.2 | 353 | 328 | 100% | 65.5% | $+412.89 | 29 | 93.1% | 65.5% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb916...f248` | A | 100.0 | 21.5 | 404.6 | 142 | 112 | 75% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 386.5 | 73 | 69 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 58.4% | existing | opposite_side_same_market |
| `0x07a9...32c8` | B | 100.0 | 26.6 | 374.3 | 34 | 32 | 49% | 107.8% | $+64.68 | 5 | 100.0% | 83.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 371.9 | 58 | 54 | 36% | 62.3% | $+118.35 | 18 | 83.3% | 25.4% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 3
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x56ac...d77e` | C | 100.0 | 35.0 | 254.4 | 111 | 15% | 0.0% | 0 | 574 | $3536 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x476e...396d` | C | 98.2 | 38.8 | 183.1 | 19 | 1% | 0.0% | 0 | 310 | $672 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 90.9 | 44.0 | 147.2 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 37.3 | 315.5 | 1288 | 1079 | 100% | 16.1% | 265 | 16.1% | 265 | $479 | existing,holder | burst_trading,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.4 | 290.2 | 1033 | 806 | 86% | 49.4% | 14 | 47.6% | 15 | $1547 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.8 | 286.1 | 1253 | 939 | 100% | 76.0% | 236 | 76.0% | 236 | $322 | existing,holder | opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 20.9 | 243.3 | 462 | 222 | 87% | 13.0% | 6 | 13.0% | 10 | $1534 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.0 | 241.0 | 567 | 526 | 90% | 88.4% | 96 | 84.4% | 107 | $1705 | existing | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.2 | 239.6 | 92 | 88 | 99% | 39.8% | 17 | 39.8% | 17 | $2734 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.0 | 235.0 | 246 | 220 | 100% | 38.0% | 29 | 38.0% | 29 | $405 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2a8c...fd3d` | C | 100.0 | 38.1 | 233.3 | 1250 | 490 | 97% | 120.1% | 5 | 119.8% | 6 | $140 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xa66c...cd67` | C | 100.0 | 39.8 | 228.9 | 531 | 451 | 63% | 10.8% | 106 | 14.1% | 150 | $5123 | existing | burst_trading,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.5 | 227.8 | 309 | 248 | 90% | 18.2% | 22 | 18.2% | 22 | $795 | existing,holder | open_copy_exposure,opposite_side_same_market |

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
| `0x7af7...4c89` | B | 100.0 | 26.9 | 331.1 | 20 | 19 | 58.1% | 3 | 58.1% | 3 | $19206 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Sports Tape Probation

- Wallets: 0
- Rule: scored sports-tape wallets with positive target-copy and sub-45 bot score, but soft flow risks; observation-only until edge windows prove out

| Wallet | Tier | Smart | Bot | ProbationScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|

No sports-tape wallets qualified for probation edge observation.

## Sports Tape Edge-Hot

- Wallets: 4
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 24.5 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 548.5 | opposite_side_same_market |
| `0x2929...1dd0` | C | 36.0 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 383.0 | open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.5 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 331.8 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 12
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.0 | 15m edge -1.77pp over 4 samples | 84.4% | 107 | existing | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 15m edge -2.63pp over 2 samples | 25.4% | 29 | existing | opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 34.4 | 15m edge -2.84pp over 3 samples | 103.5% | 99 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.0 | 15m edge -24.00pp over 2 samples | 38.0% | 29 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 30.6 | 1h edge -11.42pp over 4 samples | 114.6% | 20 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x07a9...32c8` | B | 100.0 | 26.6 | 1h edge -12.31pp over 2 samples | 83.2% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.7 | 1h edge -15.87pp over 1 samples | 109.3% | 52 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 29.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.8 | 1h edge -19.00pp over 1 samples | 149.7% | 31 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.6 | 1h edge -34.83pp over 1 samples | 65.5% | 29 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.2 | 1h edge -52.95pp over 1 samples | 81.6% | 30 | existing,holder | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 26.9 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 74.0 | $7000 | $7000 | 1 | 117.9% | 116 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | BOT | 61.3 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.0 | 652.6 | 414 | 360 | 36% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 74.0 | 623.1 | 209 | 203 | 14% | 128.1% | 27 | 117.9% | 116 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 24.5 | 603.7 | 77 | 69 | 100% | 116.7% | 25 | 116.7% | 25 | existing,holder,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.5 | 579.3 | 298 | 255 | 84% | 85.4% | 52 | 98.1% | 60 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x2929...1dd0` | watch | - | C | 100.0 | 36.0 | 412.7 | 108 | 65 | 99% | 144.5% | 5 | 144.5% | 5 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.9 | 292.3 | 978 | 885 | 74% | 30.2% | 302 | 29.9% | 385 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.8 | 252.7 | 155 | 144 | 19% | 35.9% | 56 | 37.6% | 251 | holder,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 41.6 | 185.0 | 206 | 184 | 76% | 24.9% | 67 | 20.2% | 89 | holder,sports_tape | opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 94.8 | 21.1 | 177.4 | 23 | 18 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 75.0 | 176.3 | 151 | 148 | 12% | 28.4% | 47 | 44.4% | 302 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 26.9 | 174.7 | 20 | 19 | 95% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.2 | 133.4 | 56 | 41 | 32% | 117.0% | 1 | 61.2% | 3 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 15.8 | 112.5 | 7 | 5 | 54% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 18 | 436 | 106.7% | $+5205.55 | 83.7% | 106.1% | 63.6% | 2.40x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 2 | 19 | 446 | 105.6% | $+5313.14 | 83.6% | 103.7% | 63.6% | 2.37x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 10 | 229 | 126.6% | $+3267.55 | 82.5% | 123.2% | 103.7% | 2.41x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 16 | 419 | 106.8% | $+4997.14 | 83.8% | 106.1% | 63.6% | 2.34x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 5 | 17 | 429 | 105.7% | $+5104.73 | 83.7% | 103.7% | 63.6% | 2.32x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 16 | 383 | 108.0% | $+4556.45 | 83.8% | 106.1% | 63.6% | 2.59x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 16 | 386 | 104.3% | $+4527.37 | 86.3% | 106.1% | 66.5% | 2.20x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 15 | 376 | 105.5% | $+4419.78 | 86.4% | 108.4% | 66.5% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 9 | 14 | 366 | 108.2% | $+4348.04 | 83.9% | 106.1% | 63.6% | 2.53x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 10 | 15 | 378 | 103.8% | $+4423.59 | 86.0% | 103.7% | 66.5% | 2.18x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 14 | 368 | 105.0% | $+4316.00 | 86.1% | 106.1% | 66.5% | 2.21x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 12 | 13 | 323 | 106.8% | $+3770.68 | 87.0% | 108.4% | 66.5% | 2.41x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
