# Strategy Lab Report

**Generated:** 2026-08-14 23:03 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (402)
- Valid strategies found: 83
- Candidate layers: 11 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 45 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 259
- Aggregate copy ROI: 117.0%
- Aggregate copy PnL: $+3462.23
- Aggregate copy win rate: 84.2%
- Median wallet CopyROI: 125.6%
- Worst included wallet CopyROI: 102.5%
- Open copy cost / closed copy capital: 1.81x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.24 | 372 | 100% | 102.5% | $+1076.53 | 86 | 80.2% | 41.4% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 27.18 | 135 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 55.3% |
| `0x162d...8944` | B | 25.04 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.08 | 42 | 64% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x819d...6e9c` | B | 29.02 | 26 | 30% | 125.6% | $+138.21 | 10 | 100.0% | 47.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 453
- Aggregate copy ROI: 71.4%
- Aggregate copy PnL: $+3907.86
- Aggregate copy win rate: 82.3%
- Worst included CopyROI: 41.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 565.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x18c2...529a` | B | 100.0 | 27.9 | 499.4 | 865 | 75% | 73.0% | $+569.27 | 60 | 88.3% | $287 | existing |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.8 | 123 | 79% | 93.9% | $+178.48 | 16 | 75.0% | $447 | existing |
| `0x092b...614e` | B | 100.0 | 25.5 | 470.5 | 137 | 96% | 93.2% | $+158.52 | 16 | 81.2% | $1082 | existing |
| `0xa35c...a113` | B | 100.0 | 29.8 | 468.9 | 121 | 100% | 77.8% | $+225.64 | 26 | 96.2% | $688 | existing,sports_tape |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 418.0 | 178 | 98% | 49.8% | $+403.06 | 53 | 81.1% | $1101 | existing |
| `0x760f...326a` | B | 100.0 | 31.0 | 413.9 | 15 | 14% | 84.4% | $+135.00 | 13 | 76.9% | $5078 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.3 | 412.1 | 21 | 33% | 83.3% | $+99.92 | 9 | 66.7% | $1895 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0xa931...1387` | B | 73.8 | 27.0 | 392.7 | 7 | 70% | 145.1% | $+43.53 | 3 | 100.0% | $2600 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x0931...e78e` | B | 100.0 | 25.8 | 373.9 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | $298 | existing |
| `0x95c1...94d7` | B | 100.0 | 26.3 | 372.7 | 50 | 50% | 41.9% | $+134.17 | 30 | 93.3% | $893 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x5135...7217` | B | 100.0 | 32.3 | 814.4 | 1356 | 362 | 100% | 235.6% | $+753.80 | 13 | 100.0% | 235.6% | existing,holder | burst_trading,open_copy_exposure |
| `0x17e2...d472` | B | 100.0 | 33.1 | 727.2 | 380 | 288 | 38% | 124.4% | $+870.95 | 67 | 97.0% | 151.5% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 33.2 | 564.6 | 455 | 317 | 95% | 136.0% | $+285.52 | 16 | 87.5% | 130.0% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.8 | 406.2 | 499 | 259 | 49% | 80.9% | $+137.53 | 15 | 93.3% | 66.5% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 377.9 | 1184 | 1132 | 100% | 51.6% | $+567.66 | 43 | 86.0% | 51.6% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6ac5...4b6e` | C | 100.0 | 43.6 | 235.7 | 1141 | 83% | 96.4% | 11 | 366 | $691 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2929...1dd0` | C | 100.0 | 43.8 | 260.7 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | $1136 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.2 | 221.1 | 271 | 171 | 85% | 12.5% | 50 | 11.7% | 57 | $669 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0xb99c...7e96` | A | 100.0 | 20.9 | 213.8 | 67 | 53 | 84% | 77.3% | 6 | 77.3% | 6 | $544 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 213.5 | 565 | 344 | 88% | 29.4% | 74 | 27.7% | 81 | $227 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 34.6 | 213.3 | 623 | 300 | 91% | 47.5% | 21 | 40.5% | 25 | $145 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 207.7 | 292 | 254 | 72% | 17.5% | 73 | 22.9% | 90 | $1129 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2fbb...7e44` | B | 100.0 | 25.8 | 207.1 | 196 | 87 | 100% | 27.6% | 25 | 27.6% | 25 | $172 | existing | opposite_side_same_market |
| `0x6e32...ab65` | C | 100.0 | 35.2 | 204.3 | 888 | 329 | 71% | 153.6% | 17 | 143.2% | 29 | $130 | existing,holder | burst_trading,open_copy_exposure |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |

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

- Wallets: 3
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 23.2 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 564.6 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.8 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 411.7 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.2 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 1h edge -47.28pp over 1 samples | 51.6% | 43 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 73.0% | 60 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 49.8% | 53 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 3
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.2 | 731.0 | 372 | 273 | 100% | 102.5% | 86 | 102.5% | 86 | existing,holder,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.8 | 613.2 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa35c...a113` | pushed | - | B | 100.0 | 29.8 | 425.0 | 121 | 119 | 100% | 77.8% | 26 | 77.8% | 26 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 259 | 117.0% | $+3462.23 | 84.2% | 125.6% | 102.5% | 1.81x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 20 | 489 | 97.6% | $+5437.57 | 85.1% | 103.2% | 62.1% | 3.00x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 10 | 215 | 118.6% | $+2941.10 | 85.6% | 130.4% | 102.5% | 1.74x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 4 | 18 | 469 | 98.7% | $+5263.13 | 85.7% | 106.2% | 65.3% | 2.96x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 19 | 480 | 97.9% | $+5337.65 | 85.4% | 104.0% | 62.1% | 3.00x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 21 | 564 | 91.5% | $+5790.01 | 84.9% | 102.5% | 41.9% | 2.69x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 7 | 20 | 544 | 92.6% | $+5665.19 | 85.7% | 103.2% | 41.9% | 2.69x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 22 | 565 | 90.5% | $+5819.20 | 85.3% | 98.2% | 41.8% | 2.85x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=10 smart>=70 |
| 9 | 25 | 603 | 88.5% | $+6126.04 | 84.1% | 93.2% | 41.8% | 2.80x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 22 | 597 | 88.9% | $+5982.69 | 83.8% | 98.2% | 41.9% | 2.56x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 18 | 471 | 95.7% | $+5101.63 | 85.4% | 103.2% | 62.1% | 3.03x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 17 | 460 | 96.5% | $+5027.11 | 85.7% | 104.0% | 65.3% | 2.99x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
