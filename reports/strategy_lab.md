# Strategy Lab Report

**Generated:** 2026-08-20 23:05 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (413)
- Valid strategies found: 60
- Candidate layers: 21 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 56 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 21
- Aggregate closed copy trades: 425
- Aggregate copy ROI: 92.8%
- Aggregate copy PnL: $+4659.09
- Aggregate copy win rate: 85.4%
- Median wallet CopyROI: 108.3%
- Worst included wallet CopyROI: 44.4%
- Open copy cost / closed copy capital: 3.59x
- Params: tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 14.82 | 60 | 64% | 160.5% | $+240.75 | 15 | 100.0% | 41.3% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.77 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x2fe2...9645` | A | 19.10 | 40 | 38% | 44.4% | $+39.94 | 9 | 88.9% | 48.2% |
| `0x18c2...529a` | B | 27.62 | 898 | 78% | 70.5% | $+563.79 | 59 | 89.8% | 45.6% |
| `0x0d66...ff1d` | B | 26.67 | 101 | 99% | 108.3% | $+476.46 | 28 | 89.3% | 37.9% |
| `0x44c4...09cb` | B | 26.19 | 123 | 62% | 110.8% | $+465.42 | 39 | 76.9% | 55.3% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x092b...614e` | B | 25.77 | 137 | 96% | 93.9% | $+168.94 | 17 | 82.3% | 61.6% |
| `0x272d...bc2f` | B | 28.59 | 534 | 52% | 60.7% | $+163.95 | 25 | 76.0% | 67.8% |
| `0x95c1...94d7` | B | 27.46 | 56 | 50% | 44.4% | $+159.96 | 34 | 94.1% | 20.0% |
| `0x819d...6e9c` | B | 27.67 | 26 | 27% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |
| `0x0931...e78e` | B | 25.83 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | 26.9% |
| `0x4dee...8ad7` | B | 27.07 | 29 | 71% | 51.3% | $+82.10 | 9 | 66.7% | 24.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 495
- Aggregate copy ROI: 114.0%
- Aggregate copy PnL: $+6850.81
- Aggregate copy win rate: 84.2%
- Worst included CopyROI: 46.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9c03...522d` | B | 100.0 | 33.2 | 1205.0 | 186 | 65% | 244.1% | $+1439.94 | 57 | 96.5% | $767 | existing |
| `0x3c2b...825e` | B | 100.0 | 34.4 | 829.0 | 534 | 96% | 150.4% | $+1127.82 | 66 | 95.5% | $282 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 32.1 | 778.0 | 56 | 29% | 150.8% | $+874.55 | 43 | 69.8% | $519 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.2 | 619.5 | 50 | 88% | 211.4% | $+232.59 | 6 | 66.7% | $660 | existing,holder |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 535.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing,holder |
| `0x760f...326a` | B | 100.0 | 30.8 | 524.0 | 15 | 14% | 112.2% | $+246.72 | 18 | 83.3% | $4990 | existing |
| `0xa35c...a113` | B | 100.0 | 34.5 | 488.2 | 143 | 100% | 75.7% | $+333.23 | 40 | 95.0% | $704 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.8 | 431.3 | 48 | 31% | 57.0% | $+267.89 | 45 | 80.0% | $1342 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.7 | 414.2 | 206 | 98% | 46.2% | $+430.14 | 64 | 76.6% | $1159 | existing,holder |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x8fd6...c2b8` | B | 100.0 | 6.7 | 366.8 | 16 | 32% | 94.6% | $+37.84 | 4 | 100.0% | $633 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.9 | 947.2 | 364 | 270 | 36% | 172.5% | $+1276.29 | 71 | 97.2% | 175.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.9 | 817.3 | 793 | 293 | 63% | 165.0% | $+1451.74 | 41 | 100.0% | 150.7% | existing | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 33.8 | 568.2 | 1380 | 350 | 100% | 102.7% | $+1006.18 | 36 | 100.0% | 102.7% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.3 | 499.4 | 607 | 253 | 93% | 112.6% | $+225.18 | 16 | 100.0% | 111.4% | existing,holder | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 435.1 | 1172 | 1120 | 100% | 59.3% | $+866.13 | 64 | 90.6% | 59.3% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x90a1...fb15` | B | 100.0 | 34.5 | 388.4 | 291 | 39 | 47% | 73.2% | $+109.73 | 14 | 78.6% | 67.8% | existing | opposite_side_same_market |
| `0x4b59...3aa6` | B | 100.0 | 32.2 | 381.6 | 1179 | 300 | 100% | 88.7% | $+345.93 | 8 | 100.0% | 88.7% | existing | burst_trading,open_copy_exposure |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6ac5...4b6e` | C | 100.0 | 40.0 | 244.2 | 1149 | 80% | 352.0% | 2 | 424 | $520 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x65c1...2988` | C | 100.0 | 35.5 | 217.2 | 1282 | 91% | 395.0% | 1 | 407 | $531 | existing,holder,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.3 | 218.2 | 741 | 346 | 92% | 33.2% | 22 | 28.6% | 26 | $140 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 214.0 | 602 | 370 | 89% | 28.8% | 85 | 27.3% | 92 | $230 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0a7c...8964` | C | 100.0 | 42.3 | 207.0 | 994 | 402 | 85% | 122.2% | 8 | 130.8% | 11 | $114 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xb9ab...f282` | C | 100.0 | 40.2 | 201.9 | 1163 | 240 | 89% | 169.2% | 30 | 155.0% | 32 | $121 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | C | 100.0 | 41.6 | 199.2 | 1191 | 280 | 88% | 147.8% | 24 | 149.1% | 30 | $313 | existing | burst_trading,opposite_side_same_market |

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

- Wallets: 1
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.2 | 1h edge -15.87pp over 1 samples | 110.8% | 39 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 1h edge -47.28pp over 1 samples | 59.3% | 64 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.6 | 1h edge -52.95pp over 1 samples | 70.5% | 59 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.7 | 1h edge -72.59pp over 1 samples | 46.2% | 64 | existing,holder | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xb9ab...f282` | pushed | - | C | 100.0 | 40.2 | 851.6 | 1163 | 240 | 89% | 169.2% | 30 | 155.0% | 32 | existing,holder,sports_tape | burst_trading,open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 21 | 425 | 92.8% | $+4659.09 | 85.4% | 108.3% | 44.4% | 3.59x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 2 | 21 | 426 | 91.3% | $+4557.81 | 86.4% | 108.3% | 25.2% | 3.50x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 20 | 416 | 93.7% | $+4619.15 | 85.3% | 108.4% | 44.4% | 3.46x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 4 | 12 | 216 | 119.7% | $+3063.09 | 86.1% | 115.3% | 104.0% | 2.77x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 21 | 450 | 90.4% | $+4826.22 | 83.8% | 108.3% | 44.4% | 3.23x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 22 | 459 | 89.6% | $+4866.16 | 83.9% | 106.1% | 44.4% | 3.35x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 19 | 396 | 95.7% | $+4452.17 | 86.9% | 108.4% | 44.4% | 3.73x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 19 | 407 | 95.1% | $+4537.05 | 85.7% | 108.4% | 44.4% | 3.57x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 20 | 441 | 91.6% | $+4744.12 | 84.1% | 108.4% | 44.4% | 3.32x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 19 | 400 | 94.6% | $+4475.07 | 86.8% | 108.4% | 37.0% | 3.48x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 17 | 353 | 101.2% | $+4252.27 | 86.1% | 108.9% | 60.7% | 3.86x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 18 | 387 | 96.8% | $+4412.23 | 86.8% | 108.7% | 44.4% | 3.60x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
