# Strategy Lab Report

**Generated:** 2026-08-11 23:04 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (397)
- Valid strategies found: 86
- Candidate layers: 10 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 42 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 10
- Aggregate closed copy trades: 246
- Aggregate copy ROI: 120.9%
- Aggregate copy PnL: $+3494.53
- Aggregate copy win rate: 80.1%
- Median wallet CopyROI: 135.0%
- Worst included wallet CopyROI: 101.4%
- Open copy cost / closed copy capital: 2.22x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.34 | 325 | 100% | 101.4% | $+973.26 | 77 | 79.2% | 40.8% |
| `0xb2ed...4418` | A | 22.84 | 126 | 36% | 134.8% | $+606.44 | 37 | 67.6% | 57.7% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 27.23 | 135 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 55.3% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.08 | 42 | 64% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 580
- Aggregate copy ROI: 72.9%
- Aggregate copy PnL: $+4944.09
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
| `0x18c2...529a` | B | 100.0 | 27.9 | 496.7 | 857 | 74% | 72.4% | $+564.48 | 60 | 88.3% | $280 | existing |
| `0x819d...6e9c` | B | 100.0 | 30.2 | 496.7 | 26 | 33% | 125.6% | $+138.21 | 10 | 100.0% | $413 | existing |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.8 | 123 | 79% | 93.9% | $+178.48 | 16 | 75.0% | $447 | existing |
| `0x092b...614e` | B | 100.0 | 26.0 | 470.4 | 150 | 97% | 89.7% | $+170.42 | 18 | 83.3% | $1045 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.3 | 412.1 | 21 | 33% | 83.3% | $+99.92 | 9 | 66.7% | $1895 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x760f...326a` | B | 100.0 | 31.4 | 390.9 | 13 | 13% | 83.2% | $+108.16 | 11 | 72.7% | $4606 | existing |
| `0x0e24...7014` | B | 100.0 | 29.8 | 389.1 | 153 | 98% | 44.3% | $+309.94 | 45 | 77.8% | $1074 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.2 | 378.0 | 72 | 65% | 44.5% | $+173.61 | 32 | 62.5% | $572 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.0 | 626.0 | 381 | 278 | 38% | 105.7% | $+633.92 | 57 | 96.5% | 133.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 32.3 | 519.7 | 394 | 262 | 94% | 143.4% | $+215.15 | 10 | 80.0% | 134.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 412.8 | 461 | 238 | 45% | 85.3% | $+136.49 | 14 | 92.9% | 69.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 0
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No leaderboard-only wallets passed the scout filters.

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3ce...a57a` | B | 100.0 | 34.0 | 335.3 | 1176 | 1124 | 100% | 20.8% | 20 | 20.8% | 20 | $1593 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x2929...1dd0` | C | 100.0 | 43.9 | 260.7 | 818 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | $1132 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xd3d3...bb8e` | C | 100.0 | 39.7 | 236.8 | 816 | 581 | 70% | 13.3% | 168 | 13.3% | 248 | $638 | existing,holder | opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.1 | 220.8 | 266 | 170 | 84% | 12.5% | 50 | 11.7% | 57 | $676 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.9 | 213.4 | 275 | 123 | 100% | 25.8% | 40 | 25.8% | 40 | $181 | existing | opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 34.1 | 209.7 | 532 | 266 | 89% | 52.2% | 18 | 43.5% | 22 | $149 | existing,holder | open_copy_exposure,opposite_side_same_market |
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
| `0xa75b...772c` | A | 23.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 556.9 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.9 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 411.6 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.2 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 22.8 | 1h edge -19.00pp over 1 samples | 134.8% | 37 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.0 | 1h edge -47.28pp over 1 samples | 20.8% | 20 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 72.4% | 60 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.8 | 1h edge -72.59pp over 1 samples | 44.3% | 45 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.3 | 705.6 | 325 | 242 | 100% | 101.4% | 77 | 101.4% | 77 | existing,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.9 | 613.2 | 818 | 637 | 98% | 109.6% | 31 | 108.5% | 32 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 10 | 246 | 120.9% | $+3494.53 | 80.1% | 135.0% | 101.4% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 18 | 452 | 100.7% | $+5278.34 | 82.3% | 104.9% | 65.3% | 3.24x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 16 | 432 | 101.9% | $+5087.17 | 82.9% | 108.5% | 65.3% | 3.22x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 15 | 423 | 99.6% | $+4851.15 | 82.7% | 108.4% | 65.3% | 3.25x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 16 | 434 | 98.8% | $+4942.40 | 82.5% | 104.9% | 65.3% | 3.28x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 15 | 395 | 98.7% | $+4480.73 | 84.3% | 108.4% | 65.3% | 3.13x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 16 | 406 | 97.9% | $+4571.98 | 84.0% | 104.9% | 65.3% | 3.16x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 20 | 571 | 88.6% | $+5920.58 | 82.1% | 97.7% | 42.2% | 2.64x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 23 | 600 | 87.4% | $+6193.85 | 81.5% | 89.7% | 42.2% | 2.63x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 21 | 603 | 86.2% | $+6094.19 | 81.1% | 93.9% | 42.2% | 2.52x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 14 | 386 | 96.0% | $+4244.71 | 84.2% | 104.9% | 65.3% | 3.16x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 15 | 397 | 95.3% | $+4335.96 | 83.9% | 101.4% | 65.3% | 3.19x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
