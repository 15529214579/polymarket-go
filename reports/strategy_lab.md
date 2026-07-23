# Strategy Lab Report

**Generated:** 2026-07-23 23:37 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 135
- Candidate layers: 11 core + 20 watch + 10 sports + 6 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 48 total
- Live-edge blocked push wallets: 10
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 224
- Aggregate copy ROI: 127.7%
- Aggregate copy PnL: $+3268.52
- Aggregate copy win rate: 83.0%
- Median wallet CopyROI: 135.5%
- Worst included wallet CopyROI: 103.8%
- Open copy cost / closed copy capital: 2.76x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.54 | 117 | 42% | 143.7% | $+545.89 | 32 | 65.6% | 48.7% |
| `0xa75b...772c` | A | 24.26 | 127 | 99% | 103.8% | $+467.06 | 33 | 78.8% | 32.0% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x2952...f50d` | A | 23.79 | 79 | 55% | 107.3% | $+343.22 | 24 | 87.5% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 57 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 60 | 90% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 26.45 | 159 | 60% | 109.3% | $+601.01 | 50 | 80.0% | 47.6% |
| `0xeb8b...6d8a` | B | 28.11 | 27 | 51% | 138.5% | $+124.64 | 9 | 100.0% | 42.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 652
- Aggregate copy ROI: 106.0%
- Aggregate copy PnL: $+9218.93
- Aggregate copy win rate: 90.3%
- Worst included CopyROI: 60.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.3 | 897.8 | 25 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1207 | existing |
| `0x94a8...5204` | B | 100.0 | 25.1 | 828.1 | 28 | 29% | 378.2% | $+264.76 | 4 | 100.0% | $2966 | existing |
| `0xfdff...4adc` | B | 100.0 | 33.7 | 752.8 | 17 | 24% | 251.8% | $+251.77 | 8 | 75.0% | $2953 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 738.7 | 39 | 95% | 248.6% | $+223.76 | 7 | 100.0% | $7773 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.2 | 698.1 | 1338 | 100% | 108.3% | $+2404.26 | 98 | 95.9% | $1021 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0x2a35...9015` | B | 100.0 | 34.3 | 633.7 | 88 | 99% | 132.6% | $+371.23 | 22 | 95.5% | $1854 | existing,holder |
| `0x89cf...5f47` | B | 100.0 | 33.5 | 621.5 | 67 | 42% | 111.7% | $+480.36 | 42 | 88.1% | $655 | existing |
| `0xcc6e...fa6f` | B | 100.0 | 26.2 | 619.4 | 294 | 83% | 96.6% | $+666.69 | 61 | 90.2% | $2141 | existing,sports_tape |
| `0xf4e1...34a5` | B | 100.0 | 30.6 | 574.8 | 68 | 41% | 115.2% | $+322.47 | 23 | 82.6% | $5973 | existing |
| `0xde24...4ded` | B | 100.0 | 30.8 | 563.6 | 145 | 100% | 103.8% | $+508.64 | 25 | 100.0% | $2689 | existing |
| `0x9caf...94dc` | B | 100.0 | 31.5 | 561.2 | 64 | 50% | 103.7% | $+300.79 | 28 | 96.4% | $2865 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 546.1 | 12 | 19% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x4ca8...ecd4` | B | 100.0 | 14.0 | 538.8 | 27 | 90% | 194.2% | $+77.67 | 4 | 100.0% | $661 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.4 | 526.8 | 877 | 93% | 60.1% | $+1219.76 | 184 | 89.1% | $300 | existing |
| `0x7992...1fc1` | A | 100.0 | 18.6 | 524.5 | 21 | 21% | 86.7% | $+294.71 | 29 | 79.3% | $2688 | existing |
| `0x18c2...529a` | B | 100.0 | 28.3 | 494.6 | 563 | 70% | 80.9% | $+396.49 | 39 | 82.0% | $271 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x84cd...7565` | A | 100.0 | 23.7 | 517.3 | 116 | 57 | 72% | 115.0% | $+195.57 | 14 | 92.9% | 93.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xc117...7410` | B | 84.8 | 25.4 | 464.8 | 33 | 27 | 77% | 143.0% | $+143.01 | 7 | 42.9% | 118.6% | existing | opposite_side_same_market |
| `0x141a...d05a` | B | 100.0 | 29.1 | 462.9 | 145 | 136 | 100% | 72.8% | $+502.13 | 30 | 93.3% | 72.8% | existing | opposite_side_same_market |
| `0xa8b9...775d` | B | 100.0 | 33.6 | 448.5 | 364 | 252 | 93% | 68.9% | $+571.98 | 43 | 55.8% | 64.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 29.0 | 444.0 | 162 | 79 | 100% | 71.4% | $+192.74 | 27 | 96.3% | 71.4% | existing | opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.9 | 432.9 | 403 | 377 | 100% | 67.7% | $+548.58 | 38 | 94.7% | 67.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.1 | 431.6 | 80 | 80 | 59% | 76.9% | $+184.49 | 22 | 86.4% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5a56...10f5` | B | 100.0 | 31.8 | 420.2 | 90 | 75 | 87% | 94.5% | $+245.67 | 11 | 63.6% | 88.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 6
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xc3ed...9f79` | C | 100.0 | 36.5 | 231.5 | 134 | 54% | 57.4% | 44 | 248 | $9404 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |
| `0x9b12...5ee4` | B | 100.0 | 33.8 | 221.4 | 96 | 7% | 0.0% | 0 | 369 | $613 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x162c...6c47` | C | 100.0 | 39.9 | 210.6 | 314 | 29% | 0.0% | 0 | 364 | $655 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 43.9 | 194.5 | 645 | 44% | 0.0% | 0 | 271 | $1070 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x1465...7072` | C | 100.0 | 35.7 | 188.2 | 0 | 0% | 0.0% | 0 | 594 | $988 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 100.0 | 38.9 | 169.9 | 19 | 1% | 0.0% | 0 | 309 | $674 | existing,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 37.6 | 312.5 | 1286 | 1083 | 100% | 16.5% | 267 | 16.5% | 267 | $480 | existing | burst_trading,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 25.4 | 301.2 | 1020 | 801 | 86% | 76.6% | 15 | 74.0% | 16 | $1620 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 43.0 | 278.9 | 1162 | 904 | 100% | 77.1% | 224 | 77.1% | 224 | $336 | existing | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 41.5 | 243.2 | 550 | 511 | 90% | 85.2% | 95 | 81.4% | 106 | $1762 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.0 | 236.8 | 93 | 89 | 93% | 40.5% | 18 | 41.6% | 19 | $2728 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.0 | 232.1 | 251 | 222 | 100% | 31.4% | 29 | 31.4% | 29 | $397 | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 228.3 | 73 | 69 | 95% | 58.4% | 20 | 58.4% | 20 | $3677 | existing | opposite_side_same_market |
| `0x42db...d512` | B | 100.0 | 30.0 | 222.5 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | $1174 | existing,sports_tape | opposite_side_same_market |
| `0xb57f...0c96` | C | 100.0 | 42.0 | 219.4 | 872 | 552 | 64% | 70.2% | 110 | 71.1% | 181 | $379 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.2 | 216.5 | 90 | 71 | 62% | 74.2% | 14 | 48.9% | 21 | $2360 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0x7af7...4c89` | B | 100.0 | 30.5 | 316.8 | 21 | 20 | 58.1% | 3 | 58.1% | 3 | $12821 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 24.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 536.8 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.2 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 328.6 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 10
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x34a5...9450` | C | 100.0 | 41.5 | 15m edge -1.77pp over 4 samples | 81.4% | 106 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 34.2 | 15m edge -2.84pp over 3 samples | 108.3% | 98 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.0 | 15m edge -24.00pp over 2 samples | 31.4% | 29 | existing | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 30.8 | 1h edge -11.42pp over 4 samples | 103.8% | 25 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.4 | 1h edge -15.87pp over 1 samples | 109.3% | 50 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.5 | 1h edge -19.00pp over 1 samples | 143.7% | 32 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.9 | 1h edge -34.83pp over 1 samples | 67.7% | 38 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.3 | 1h edge -52.95pp over 1 samples | 80.9% | 39 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7c1e...bab3` | C | 100.0 | 43.9 | 1h edge -8.50pp over 2 samples | 0.0% | 0 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 30.5 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 73.9 | $7000 | $7000 | 1 | 117.3% | 119 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 55.4 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 73.9 | 627.3 | 224 | 217 | 15% | 123.9% | 31 | 117.3% | 119 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.5 | 619.5 | 399 | 351 | 35% | 149.0% | 17 | 122.5% | 41 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 24.3 | 586.9 | 127 | 105 | 99% | 103.8% | 33 | 103.8% | 33 | existing,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.2 | 567.8 | 294 | 249 | 83% | 83.4% | 51 | 96.6% | 61 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.9 | 295.8 | 993 | 895 | 74% | 30.6% | 304 | 30.1% | 388 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.0 | 261.9 | 159 | 148 | 17% | 38.0% | 58 | 36.4% | 266 | sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 66.9 | 205.5 | 142 | 139 | 12% | 31.7% | 45 | 46.0% | 311 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.0 | 182.1 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | D | 100.0 | 57.9 | 181.7 | 257 | 232 | 78% | 27.3% | 83 | 22.3% | 107 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 83.3 | 20.0 | 171.2 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 30.5 | 160.9 | 21 | 20 | 66% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 21.8 | 102.3 | 20 | 5 | 44% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 74.9 | 10.0 | 96.6 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xb3c1...4837` | prompt | - | B | 100.0 | 21.1 | 81.7 | 198 | 150 | 91% | -0.1% | 21 | 0.3% | 23 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 224 | 127.7% | $+3268.52 | 83.0% | 135.5% | 103.8% | 2.76x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 10 | 200 | 130.6% | $+2925.30 | 82.5% | 137.0% | 103.8% | 2.98x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 10 | 215 | 127.3% | $+3143.88 | 82.3% | 135.4% | 103.8% | 2.77x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 21 | 505 | 100.8% | $+6206.73 | 85.1% | 103.8% | 66.5% | 3.68x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 5 | 10 | 192 | 124.9% | $+2722.63 | 85.9% | 135.4% | 103.8% | 2.60x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 22 | 516 | 99.9% | $+6314.16 | 84.9% | 100.2% | 66.5% | 3.63x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 10 | 168 | 129.2% | $+2492.61 | 81.0% | 135.4% | 81.2% | 3.10x | tier>=A bot<25 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 12 | 221 | 120.9% | $+3130.54 | 81.4% | 121.8% | 81.2% | 2.64x | tier>=A bot<25 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 20 | 496 | 100.2% | $+6082.09 | 84.9% | 100.2% | 66.5% | 3.70x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 10 | 19 | 452 | 101.3% | $+5568.80 | 85.4% | 103.8% | 66.5% | 3.96x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 11 | 21 | 507 | 99.4% | $+6189.52 | 84.6% | 96.6% | 66.5% | 3.64x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 20 | 473 | 97.9% | $+5660.84 | 86.5% | 100.2% | 66.5% | 3.68x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
