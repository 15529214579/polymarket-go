# Strategy Lab Report

**Generated:** 2026-08-12 23:49 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (399)
- Valid strategies found: 90
- Candidate layers: 11 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 43 total
- Live-edge blocked push wallets: 9
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 354
- Aggregate copy ROI: 123.6%
- Aggregate copy PnL: $+6502.64
- Aggregate copy win rate: 84.7%
- Median wallet CopyROI: 134.8%
- Worst included wallet CopyROI: 100.7%
- Open copy cost / closed copy capital: 2.62x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 22.97 | 349 | 100% | 100.7% | $+1017.04 | 82 | 79.3% | 41.2% |
| `0xb2ed...4418` | A | 22.82 | 127 | 36% | 134.8% | $+606.44 | 37 | 67.6% | 57.7% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0xbb35...b62a` | B | 25.02 | 1324 | 100% | 127.8% | $+2964.33 | 103 | 96.1% | 70.4% |
| `0x44c4...09cb` | B | 27.23 | 135 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 55.3% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.08 | 42 | 64% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 581
- Aggregate copy ROI: 72.9%
- Aggregate copy PnL: $+4951.55
- Aggregate copy win rate: 84.5%
- Worst included CopyROI: 44.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xadba...a6e2` | B | 100.0 | 30.8 | 671.0 | 24 | 75% | 206.3% | $+165.05 | 8 | 100.0% | $977 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.2 | 560.7 | 28 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $795 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.9 | 556.1 | 112 | 100% | 109.6% | $+263.03 | 23 | 82.6% | $582 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.4 | 526.0 | 929 | 93% | 59.5% | $+1250.07 | 191 | 89.5% | $293 | existing |
| `0x18c2...529a` | B | 100.0 | 27.9 | 496.7 | 855 | 74% | 72.4% | $+564.48 | 60 | 88.3% | $280 | existing |
| `0x819d...6e9c` | B | 100.0 | 30.2 | 496.7 | 26 | 33% | 125.6% | $+138.21 | 10 | 100.0% | $413 | existing |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.8 | 123 | 79% | 93.9% | $+178.48 | 16 | 75.0% | $447 | existing |
| `0x092b...614e` | B | 100.0 | 26.1 | 470.3 | 149 | 97% | 89.7% | $+170.42 | 18 | 83.3% | $1049 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.3 | 412.1 | 21 | 33% | 83.3% | $+99.92 | 9 | 66.7% | $1895 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x760f...326a` | B | 100.0 | 31.4 | 397.0 | 13 | 13% | 82.6% | $+115.62 | 12 | 75.0% | $4606 | existing |
| `0x0e24...7014` | B | 100.0 | 29.8 | 389.1 | 153 | 98% | 44.3% | $+309.94 | 45 | 77.8% | $1074 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.2 | 378.0 | 72 | 65% | 44.5% | $+173.61 | 32 | 62.5% | $572 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.0 | 643.6 | 380 | 282 | 38% | 107.8% | $+689.59 | 61 | 96.7% | 136.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 32.0 | 520.1 | 409 | 277 | 95% | 143.4% | $+215.15 | 10 | 80.0% | 134.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | B | 100.0 | 34.8 | 442.9 | 506 | 111 | 94% | 113.7% | $+670.82 | 7 | 100.0% | 111.3% | existing,holder | burst_trading |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.5 | 413.0 | 463 | 240 | 45% | 85.3% | $+136.49 | 14 | 92.9% | 69.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2dc1...b33c` | C | 100.0 | 37.8 | 190.9 | 0 | 0% | 123.1% | 1 | 348 | $8132 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 335.0 | 1181 | 1127 | 100% | 39.6% | 23 | 39.6% | 23 | $1542 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x2929...1dd0` | C | 100.0 | 43.8 | 260.7 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | $1136 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xd3d3...bb8e` | C | 100.0 | 39.8 | 237.1 | 823 | 586 | 70% | 17.5% | 170 | 16.4% | 251 | $634 | existing,holder | opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.1 | 220.8 | 266 | 170 | 84% | 12.5% | 50 | 11.7% | 57 | $676 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.7 | 212.0 | 241 | 108 | 100% | 27.4% | 33 | 27.4% | 33 | $185 | existing | opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 34.3 | 210.3 | 550 | 273 | 89% | 52.2% | 18 | 43.5% | 22 | $148 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 207.7 | 292 | 254 | 72% | 17.5% | 73 | 22.9% | 90 | $1129 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 23.0 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 561.0 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.8 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 411.7 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 9
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 25.0 | 15m edge -2.84pp over 3 samples | 127.8% | 103 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.2 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 22.8 | 1h edge -19.00pp over 1 samples | 134.8% | 37 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 39.6% | 23 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 72.4% | 60 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.8 | 1h edge -72.59pp over 1 samples | 44.3% | 45 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.0 | 712.3 | 349 | 259 | 100% | 100.7% | 82 | 100.7% | 82 | existing,holder,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.8 | 613.2 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 354 | 123.6% | $+6502.64 | 84.7% | 134.8% | 100.7% | 2.62x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 19 | 559 | 108.9% | $+8268.60 | 84.8% | 108.4% | 65.3% | 3.21x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 17 | 540 | 110.0% | $+8095.47 | 85.4% | 108.6% | 65.3% | 3.19x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 10 | 345 | 121.9% | $+6266.62 | 84.6% | 131.3% | 100.7% | 2.64x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 10 | 317 | 122.6% | $+5896.20 | 86.8% | 131.5% | 100.7% | 2.47x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 17 | 541 | 107.9% | $+7932.66 | 85.0% | 108.4% | 65.3% | 3.24x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 17 | 513 | 107.7% | $+7562.24 | 86.4% | 108.4% | 65.3% | 3.15x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 16 | 531 | 108.6% | $+7859.45 | 85.3% | 108.5% | 65.3% | 3.21x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 16 | 503 | 108.4% | $+7489.03 | 86.7% | 108.5% | 65.3% | 3.12x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 21 | 679 | 98.7% | $+8928.88 | 84.2% | 100.7% | 42.2% | 2.77x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 22 | 711 | 96.4% | $+9102.49 | 83.3% | 97.3% | 42.2% | 2.67x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 24 | 707 | 97.3% | $+9184.11 | 83.6% | 91.8% | 42.2% | 2.76x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
