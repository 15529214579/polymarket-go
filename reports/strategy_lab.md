# Strategy Lab Report

**Generated:** 2026-08-19 23:15 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (412)
- Valid strategies found: 69
- Candidate layers: 15 core + 20 watch + 10 sports + 3 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 53 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 2 observation-only

## Selected Core Strategy

- Wallets: 15
- Aggregate closed copy trades: 320
- Aggregate copy ROI: 117.9%
- Aggregate copy PnL: $+4517.35
- Aggregate copy win rate: 83.8%
- Median wallet CopyROI: 119.9%
- Worst included wallet CopyROI: 101.7%
- Open copy cost / closed copy capital: 2.21x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.25 | 419 | 100% | 101.7% | $+1138.77 | 93 | 79.6% | 42.1% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 14.61 | 57 | 63% | 165.8% | $+232.07 | 14 | 100.0% | 39.7% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0x586a...be99` | A | 24.55 | 35 | 80% | 162.6% | $+146.36 | 8 | 75.0% | 33.9% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 26.19 | 123 | 62% | 110.8% | $+465.42 | 39 | 76.9% | 55.3% |
| `0x0d66...ff1d` | B | 26.31 | 87 | 99% | 107.2% | $+418.25 | 23 | 87.0% | 34.1% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x819d...6e9c` | B | 27.67 | 26 | 27% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 543
- Aggregate copy ROI: 108.1%
- Aggregate copy PnL: $+7110.79
- Aggregate copy win rate: 84.7%
- Worst included CopyROI: 48.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9c03...522d` | B | 100.0 | 33.2 | 1205.0 | 186 | 65% | 244.1% | $+1439.94 | 57 | 96.5% | $767 | existing |
| `0x3c2b...825e` | B | 100.0 | 34.2 | 840.8 | 527 | 96% | 154.3% | $+1110.96 | 64 | 95.3% | $273 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 32.2 | 787.0 | 54 | 28% | 155.7% | $+841.00 | 41 | 68.3% | $470 | existing |
| `0x54fc...76eb` | B | 100.0 | 31.3 | 617.6 | 46 | 87% | 211.4% | $+232.59 | 6 | 66.7% | $639 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x760f...326a` | B | 100.0 | 30.8 | 524.0 | 15 | 14% | 112.2% | $+246.72 | 18 | 83.3% | $4990 | existing |
| `0x18c2...529a` | B | 100.0 | 27.6 | 492.0 | 895 | 78% | 70.5% | $+563.79 | 59 | 89.8% | $295 | existing,holder |
| `0xa35c...a113` | B | 100.0 | 33.9 | 487.7 | 149 | 100% | 75.7% | $+333.23 | 40 | 95.0% | $709 | existing |
| `0x092b...614e` | B | 100.0 | 25.6 | 470.3 | 135 | 96% | 93.2% | $+158.52 | 16 | 81.2% | $1079 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x7fcf...80ac` | A | 100.0 | 10.8 | 443.2 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | $1826 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.8 | 431.3 | 48 | 31% | 57.0% | $+267.89 | 45 | 80.0% | $1342 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0xa931...1387` | B | 100.0 | 29.3 | 427.1 | 10 | 77% | 128.8% | $+51.51 | 4 | 100.0% | $2630 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 415.8 | 192 | 98% | 48.1% | $+418.37 | 59 | 79.7% | $1108 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.0 | 403.2 | 72 | 63% | 50.5% | $+207.07 | 34 | 64.7% | $569 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.9 | 947.2 | 364 | 270 | 36% | 172.5% | $+1276.29 | 71 | 97.2% | 175.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.9 | 710.8 | 809 | 304 | 64% | 145.9% | $+1065.39 | 34 | 100.0% | 142.6% | existing,holder | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 33.5 | 609.5 | 1380 | 364 | 100% | 114.4% | $+1087.02 | 36 | 100.0% | 114.4% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.4 | 459.7 | 603 | 249 | 93% | 129.2% | $+129.19 | 8 | 100.0% | 129.2% | existing,holder | burst_trading,open_copy_exposure |
| `0x21cc...54bc` | B | 100.0 | 28.7 | 424.6 | 170 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.7 | 388.7 | 529 | 279 | 52% | 73.8% | $+132.83 | 16 | 87.5% | 59.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4b59...3aa6` | B | 100.0 | 32.1 | 377.1 | 1177 | 295 | 100% | 99.3% | $+238.41 | 6 | 100.0% | 99.3% | existing,holder | burst_trading,open_copy_exposure |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 3
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x772f...8145` | B | 100.0 | 19.7 | 216.0 | 5 | 4% | 2.3% | 28 | 119 | $7628 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |
| `0x65c1...2988` | C | 100.0 | 43.7 | 198.8 | 1275 | 91% | 0.0% | 0 | 400 | $527 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xd0fd...a7da` | B | 100.0 | 30.7 | 169.1 | 3 | 4% | 60.5% | 5 | 55 | $1293 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.2 | 217.8 | 722 | 342 | 92% | 36.2% | 24 | 31.5% | 28 | $143 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.4 | 216.2 | 594 | 363 | 89% | 28.8% | 82 | 27.2% | 89 | $226 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.1 | 210.9 | 297 | 258 | 71% | 17.5% | 73 | 22.6% | 92 | $1110 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0a7c...8964` | C | 100.0 | 42.7 | 206.6 | 1001 | 401 | 86% | 112.7% | 9 | 123.8% | 12 | $113 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb9ab...f282` | C | 100.0 | 40.2 | 202.6 | 1159 | 242 | 90% | 169.2% | 30 | 155.0% | 32 | $122 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
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

- Wallets: 2
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 23.2 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 565.0 | opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 5
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.2 | 1h edge -15.87pp over 1 samples | 110.8% | 39 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.6 | 1h edge -52.95pp over 1 samples | 70.5% | 59 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 48.1% | 59 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 4
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xb9ab...f282` | pushed | - | C | 100.0 | 40.2 | 852.0 | 1159 | 242 | 90% | 169.2% | 30 | 155.0% | 32 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.2 | 738.7 | 419 | 298 | 100% | 101.7% | 93 | 101.7% | 93 | existing,holder,sports_tape | opposite_side_same_market |
| `0x4c9c...e006` | reject | - | D | 100.0 | 54.1 | 182.1 | 1368 | 776 | 97% | 14.5% | 182 | 15.7% | 188 | existing,holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x183f...85e3` | watch | - | C | 100.0 | 41.9 | 95.3 | 1311 | 762 | 99% | 0.0% | 0 | 0.0% | 0 | existing,sports_tape | burst_trading,open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 15 | 320 | 117.9% | $+4517.35 | 83.8% | 119.9% | 101.7% | 2.21x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 24 | 529 | 97.1% | $+6106.31 | 84.1% | 107.8% | 44.4% | 3.09x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 3 | 25 | 542 | 94.8% | $+6050.47 | 84.7% | 107.2% | 25.2% | 2.97x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 25 | 563 | 94.2% | $+6313.38 | 82.9% | 107.2% | 44.4% | 2.93x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 23 | 520 | 97.8% | $+6066.37 | 84.0% | 108.4% | 44.4% | 2.99x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 6 | 19 | 431 | 106.6% | $+5532.16 | 84.9% | 109.8% | 65.3% | 2.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 22 | 500 | 99.7% | $+5899.39 | 85.2% | 108.7% | 44.4% | 3.17x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 24 | 554 | 94.9% | $+6273.44 | 82.9% | 107.8% | 44.4% | 2.82x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 22 | 511 | 99.1% | $+5984.27 | 84.3% | 108.7% | 44.4% | 3.06x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 22 | 504 | 98.7% | $+5922.29 | 85.1% | 108.7% | 37.0% | 2.98x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 23 | 545 | 96.0% | $+6191.34 | 83.1% | 108.4% | 44.4% | 2.89x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 21 | 491 | 100.5% | $+5859.45 | 85.1% | 108.9% | 44.4% | 3.06x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
