# Strategy Lab Report

**Generated:** 2026-08-02 23:36 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (215)
- Valid strategies found: 200
- Candidate layers: 16 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 50 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 16
- Aggregate closed copy trades: 297
- Aggregate copy ROI: 130.5%
- Aggregate copy PnL: $+4437.13
- Aggregate copy win rate: 83.5%
- Median wallet CopyROI: 135.4%
- Worst included wallet CopyROI: 102.7%
- Open copy cost / closed copy capital: 2.75x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.89 | 211 | 100% | 113.9% | $+797.15 | 53 | 83.0% | 36.0% |
| `0xb2ed...4418` | A | 23.17 | 125 | 39% | 134.4% | $+591.26 | 36 | 66.7% | 58.8% |
| `0xe745...5681` | A | 5.86 | 64 | 32% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 26.79 | 143 | 64% | 108.6% | $+521.13 | 44 | 77.3% | 47.6% |
| `0x5b1d...3721` | B | 27.55 | 45 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.9% |
| `0x092b...614e` | B | 26.93 | 158 | 97% | 102.7% | $+205.35 | 19 | 84.2% | 55.4% |
| `0xe916...7e93` | B | 27.86 | 62 | 98% | 117.3% | $+152.51 | 12 | 83.3% | 44.7% |
| `0xffa1...6340` | B | 25.47 | 52 | 98% | 165.2% | $+148.73 | 9 | 88.9% | 18.0% |
| `0x0ec9...1e0c` | B | 26.70 | 7 | 13% | 183.6% | $+146.85 | 8 | 87.5% | 11.1% |
| `0xeb8b...6d8a` | B | 28.11 | 28 | 53% | 138.5% | $+124.64 | 9 | 100.0% | 42.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 679
- Aggregate copy ROI: 91.5%
- Aggregate copy PnL: $+8437.83
- Aggregate copy win rate: 89.1%
- Worst included CopyROI: 59.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xbb35...b62a` | B | 100.0 | 33.7 | 754.3 | 1333 | 100% | 118.8% | $+2613.22 | 101 | 97.0% | $1093 | existing,holder |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 59 | 64% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 29 | 48% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xcc6e...fa6f` | B | 100.0 | 25.5 | 605.4 | 257 | 82% | 97.4% | $+574.59 | 51 | 88.2% | $2172 | existing,sports_tape |
| `0x17fe...b0ca` | B | 100.0 | 31.0 | 568.0 | 175 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $674 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.7 | 559.9 | 27 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $754 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 525.7 | 899 | 93% | 59.8% | $+1220.55 | 185 | 89.2% | $296 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |
| `0x18c2...529a` | B | 100.0 | 28.1 | 486.2 | 820 | 74% | 70.3% | $+520.04 | 57 | 86.0% | $271 | existing,holder |
| `0x819d...6e9c` | B | 100.0 | 31.8 | 485.5 | 21 | 32% | 134.1% | $+107.24 | 8 | 100.0% | $419 | existing |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.2 | 125 | 80% | 93.9% | $+178.48 | 16 | 75.0% | $443 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0xe410...d4d1` | B | 100.0 | 30.7 | 461.6 | 64 | 74% | 85.8% | $+180.14 | 19 | 94.7% | $969 | existing |
| `0x2931...7b81` | B | 100.0 | 30.3 | 456.8 | 251 | 99% | 66.0% | $+587.85 | 33 | 93.9% | $634 | existing,holder |
| `0x578e...c3c0` | A | 100.0 | 24.1 | 451.5 | 57 | 34% | 59.8% | $+286.98 | 46 | 80.4% | $1277 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7c73...6ee3` | B | 100.0 | 25.9 | 427.9 | 49 | 96% | 102.0% | $+71.40 | 7 | 85.7% | $5610 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.0 | 621.1 | 382 | 272 | 40% | 104.9% | $+618.86 | 56 | 96.4% | 127.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xcefd...d6aa` | B | 100.0 | 30.0 | 598.6 | 86 | 62 | 68% | 177.3% | $+248.17 | 9 | 100.0% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xeec5...1862` | B | 89.5 | 31.3 | 586.7 | 63 | 61 | 88% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 33.7 | 442.1 | 320 | 189 | 93% | 131.6% | $+157.95 | 7 | 71.4% | 121.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x2a99...51bb` | A | 100.0 | 24.1 | 403.4 | 688 | 498 | 99% | 41.3% | $+520.57 | 123 | 76.4% | 41.3% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xf56d...c3d9` | B | 100.0 | 26.5 | 399.3 | 281 | 171 | 94% | 54.5% | $+310.75 | 38 | 71.0% | 55.5% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 25.4 | 392.7 | 1058 | 839 | 90% | 74.0% | $+199.72 | 16 | 75.0% | 74.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 386.5 | 73 | 69 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 58.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3f3a...e8fd` | C | 100.0 | 44.3 | 212.5 | 334 | 24% | 9.3% | 2 | 561 | $662 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x2dc1...b33c` | C | 100.0 | 44.0 | 173.9 | 0 | 0% | 123.1% | 1 | 372 | $6712 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure,opposite_side_same_market |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xd3d3...bb8e` | C | 100.0 | 39.4 | 239.4 | 850 | 613 | 72% | 14.1% | 179 | 12.3% | 254 | $645 | existing | opposite_side_same_market |
| `0xde24...4ded` | C | 100.0 | 37.6 | 232.1 | 344 | 311 | 100% | 91.7% | 45 | 91.7% | 45 | $1812 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2929...1dd0` | C | 100.0 | 44.0 | 236.2 | 517 | 405 | 99% | 109.2% | 24 | 109.2% | 24 | $1105 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0378...7634` | B | 100.0 | 26.3 | 226.1 | 378 | 181 | 100% | 34.7% | 63 | 34.7% | 63 | $385 | existing,holder | opposite_side_same_market |
| `0x4bba...cf14` | B | 100.0 | 26.3 | 217.5 | 772 | 213 | 100% | 20.0% | 6 | 20.0% | 6 | $102 | existing | open_copy_exposure,opposite_side_same_market |
| `0xe896...a4c7` | C | 100.0 | 38.2 | 216.0 | 571 | 328 | 94% | 37.7% | 72 | 38.4% | 80 | $249 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.2 | 215.9 | 125 | 118 | 98% | 42.5% | 34 | 44.6% | 36 | $1024 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x7e01...b0b5` | C | 100.0 | 37.8 | 211.2 | 459 | 274 | 97% | 66.6% | 72 | 66.6% | 72 | $342 | existing | opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 207.7 | 292 | 254 | 72% | 17.5% | 73 | 22.9% | 90 | $1129 | existing | open_copy_exposure,opposite_side_same_market |

## Recent Trade Flow Scout

- Wallets: 0
- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No recent-trade wallets passed the flow filters.

## Sports Tape Hotlist

- Wallets: 0
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No recent sports-tape wallets passed the hotlist filters.

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
| `0xa75b...772c` | A | 23.9 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 567.4 | opposite_side_same_market |
| `0x2929...1dd0` | C | 44.0 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 404.2 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | B | 25.5 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 321.1 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 33.7 | 15m edge -2.84pp over 3 samples | 118.8% | 101 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | C | 100.0 | 37.6 | 1h edge -11.42pp over 4 samples | 91.7% | 45 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.8 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.2 | 1h edge -19.00pp over 1 samples | 134.4% | 36 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.1 | 1h edge -52.95pp over 1 samples | 70.3% | 57 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.2 | 1h edge -72.59pp over 1 samples | 44.6% | 36 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 6
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.9 | 716.2 | 211 | 167 | 100% | 113.9% | 53 | 113.9% | 53 | existing,holder,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 44.0 | 555.3 | 517 | 405 | 99% | 109.2% | 24 | 109.2% | 24 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 25.5 | 536.0 | 257 | 215 | 82% | 81.7% | 41 | 97.4% | 51 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x4aec...4ce9` | reject | - | D | 100.0 | 58.6 | 314.2 | 1491 | 600 | 100% | 54.5% | 30 | 54.5% | 30 | existing,holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xdf27...026f` | reject | - | D | 100.0 | 49.5 | 175.5 | 1124 | 334 | 98% | 41.7% | 8 | 41.7% | 8 | holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xf2bb...3250` | prompt | - | B | 100.0 | 17.7 | 57.8 | 109 | 95 | 48% | 0.0% | 0 | -53.4% | 1 | holder,sports_tape | open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 16 | 297 | 130.5% | $+4437.13 | 83.5% | 135.4% | 102.7% | 2.75x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 15 | 261 | 129.9% | $+3845.87 | 85.8% | 135.5% | 102.7% | 2.59x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 14 | 217 | 134.1% | $+3324.74 | 87.6% | 137.0% | 102.7% | 2.64x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 4 | 27 | 565 | 103.9% | $+7022.39 | 84.2% | 108.6% | 62.9% | 3.83x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 13 | 258 | 130.2% | $+3814.76 | 82.9% | 135.2% | 102.7% | 2.95x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 26 | 554 | 104.8% | $+6914.96 | 84.5% | 108.8% | 62.9% | 3.88x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 7 | 26 | 529 | 101.8% | $+6431.13 | 85.4% | 108.5% | 62.9% | 3.83x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 25 | 518 | 102.7% | $+6323.70 | 85.7% | 108.6% | 62.9% | 3.89x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 9 | 12 | 222 | 129.5% | $+3223.50 | 85.6% | 135.4% | 102.7% | 2.80x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 10 | 19 | 405 | 108.4% | $+4995.13 | 88.4% | 117.3% | 65.3% | 2.62x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 11 | 21 | 471 | 101.2% | $+5457.34 | 87.5% | 113.9% | 58.4% | 2.39x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 23 | 515 | 102.7% | $+6292.59 | 84.3% | 108.4% | 62.9% | 4.07x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
