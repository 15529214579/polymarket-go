# Strategy Lab Report

**Generated:** 2026-08-13 23:32 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (401)
- Valid strategies found: 86
- Candidate layers: 19 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 19
- Aggregate closed copy trades: 463
- Aggregate copy ROI: 97.6%
- Aggregate copy PnL: $+5162.19
- Aggregate copy win rate: 84.4%
- Median wallet CopyROI: 102.0%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 3.09x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.20 | 355 | 100% | 102.0% | $+1040.30 | 83 | 79.5% | 41.4% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x162d...8944` | A | 24.50 | 3 | 1% | 104.8% | $+314.31 | 29 | 89.7% | 46.0% |
| `0x84cd...7565` | A | 23.48 | 123 | 79% | 93.9% | $+178.48 | 16 | 75.0% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 12.31 | 21 | 33% | 83.3% | $+99.92 | 9 | 66.7% | 37.3% |
| `0x18c2...529a` | B | 27.94 | 857 | 74% | 73.2% | $+563.49 | 59 | 88.1% | 45.8% |
| `0x44c4...09cb` | B | 27.23 | 135 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 55.3% |
| `0x7ec5...9fe0` | B | 26.73 | 109 | 52% | 79.9% | $+279.81 | 35 | 88.6% | 29.1% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.08 | 42 | 64% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x272d...bc2f` | B | 28.60 | 474 | 46% | 69.9% | $+216.65 | 29 | 82.8% | 68.2% |
| `0x21cc...54bc` | B | 28.82 | 168 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x092b...614e` | B | 25.91 | 143 | 97% | 93.2% | $+158.52 | 16 | 81.2% | 61.6% |
| `0xa35c...a113` | B | 28.72 | 109 | 100% | 73.5% | $+139.68 | 17 | 100.0% | 79.9% |
| `0x7673...fa40` | B | 28.36 | 83 | 60% | 66.5% | $+73.21 | 10 | 70.0% | 53.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 359
- Aggregate copy ROI: 69.0%
- Aggregate copy PnL: $+2966.51
- Aggregate copy win rate: 79.7%
- Worst included CopyROI: 42.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.9 | 556.1 | 112 | 100% | 109.6% | $+263.03 | 23 | 82.6% | $582 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x819d...6e9c` | B | 100.0 | 30.2 | 496.7 | 26 | 33% | 125.6% | $+138.21 | 10 | 100.0% | $413 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x760f...326a` | B | 100.0 | 31.1 | 413.6 | 15 | 15% | 84.4% | $+135.00 | 13 | 76.9% | $4782 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x0e24...7014` | B | 100.0 | 30.0 | 393.8 | 158 | 98% | 44.5% | $+320.55 | 46 | 78.3% | $1087 | existing,holder |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x0931...e78e` | B | 100.0 | 25.8 | 373.9 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | $298 | existing |
| `0x95c1...94d7` | B | 100.0 | 25.8 | 372.5 | 49 | 51% | 42.2% | $+130.76 | 29 | 93.1% | $864 | existing |
| `0x8fd6...c2b8` | B | 100.0 | 6.7 | 366.8 | 16 | 32% | 94.6% | $+37.84 | 4 | 100.0% | $633 | existing |
| `0xb99c...7e96` | A | 100.0 | 21.8 | 362.6 | 45 | 78% | 77.3% | $+69.60 | 6 | 83.3% | $620 | existing,holder |
| `0x6d15...04c2` | B | 100.0 | 26.5 | 359.5 | 43 | 83% | 58.1% | $+87.15 | 12 | 58.3% | $520 | existing |
| `0x1abb...0901` | B | 100.0 | 26.7 | 350.1 | 30 | 71% | 71.3% | $+49.95 | 7 | 71.4% | $1293 | existing |
| `0x6ed4...8c3c` | A | 100.0 | 24.4 | 349.0 | 38 | 70% | 82.8% | $+49.65 | 6 | 50.0% | $255 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.9 | 725.8 | 379 | 284 | 38% | 125.7% | $+841.84 | 64 | 96.9% | 150.2% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 32.4 | 519.5 | 430 | 293 | 95% | 143.4% | $+215.15 | 10 | 80.0% | 134.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 358.2 | 1186 | 1133 | 100% | 50.6% | $+475.27 | 33 | 81.8% | 50.6% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 333.3 | 422 | 232 | 75% | 39.0% | $+225.89 | 50 | 68.0% | 34.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 317.5 | 109 | 105 | 60% | 55.6% | $+94.59 | 16 | 43.8% | 46.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.8 | 309.8 | 865 | 317 | 70% | 76.3% | $+99.20 | 5 | 100.0% | 103.2% | existing,holder | burst_trading,open_copy_exposure |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 36.5 | 218.6 | 1215 | 87% | 0.0% | 0 | 378 | $501 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 37.8 | 190.6 | 0 | 0% | 123.1% | 1 | 345 | $8161 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2929...1dd0` | C | 100.0 | 43.8 | 260.7 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | $1136 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.1 | 220.8 | 271 | 171 | 84% | 12.5% | 50 | 11.7% | 57 | $665 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 34.5 | 211.1 | 586 | 282 | 90% | 46.1% | 19 | 38.8% | 23 | $146 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.9 | 210.1 | 218 | 96 | 100% | 25.1% | 29 | 25.1% | 29 | $180 | existing | opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 207.7 | 292 | 254 | 72% | 17.5% | 73 | 22.9% | 90 | $1129 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | C | 100.0 | 41.2 | 192.2 | 716 | 185 | 87% | 129.6% | 12 | 141.4% | 16 | $354 | existing,holder | burst_trading,opposite_side_same_market |
| `0xb9ab...f282` | C | 100.0 | 41.7 | 189.7 | 1098 | 247 | 84% | 167.0% | 15 | 155.7% | 16 | $123 | existing,holder | burst_trading,open_copy_exposure |
| `0x51d3...1719` | B | 100.0 | 27.6 | 184.9 | 177 | 65 | 69% | 20.4% | 18 | 48.4% | 31 | $191 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 23.2 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 563.0 | opposite_side_same_market |
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
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 50.6% | 33 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 73.2% | 59 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.0 | 1h edge -72.59pp over 1 samples | 44.5% | 46 | existing,holder | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 3
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.2 | 721.9 | 355 | 263 | 100% | 102.0% | 83 | 102.0% | 83 | existing,holder,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.8 | 613.2 | 819 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa35c...a113` | pushed | - | B | 100.0 | 28.7 | 364.3 | 109 | 106 | 100% | 73.5% | 17 | 73.5% | 17 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 19 | 463 | 97.6% | $+5162.19 | 84.4% | 102.0% | 65.3% | 3.09x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 18 | 454 | 97.9% | $+5062.27 | 84.8% | 103.4% | 65.3% | 3.10x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 10 | 244 | 116.8% | $+3269.44 | 83.2% | 122.1% | 102.0% | 1.77x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 17 | 444 | 98.6% | $+4989.06 | 85.1% | 104.8% | 65.3% | 3.05x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 19 | 518 | 92.3% | $+5387.71 | 85.1% | 102.0% | 42.2% | 2.76x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 23 | 566 | 89.6% | $+5767.76 | 83.6% | 93.2% | 42.2% | 2.75x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 7 | 20 | 528 | 91.8% | $+5460.92 | 84.8% | 98.0% | 42.2% | 2.80x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=10 smart>=70 |
| 8 | 20 | 538 | 91.1% | $+5512.53 | 84.4% | 98.0% | 42.2% | 2.76x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 17 | 445 | 95.6% | $+4826.25 | 84.7% | 102.0% | 65.3% | 3.12x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 21 | 536 | 90.3% | $+5537.66 | 85.1% | 93.9% | 42.2% | 2.73x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 16 | 435 | 96.2% | $+4753.04 | 85.1% | 103.4% | 65.3% | 3.08x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 23 | 578 | 86.0% | $+5659.16 | 84.1% | 93.2% | 25.1% | 2.75x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=10 smart>=70 |
