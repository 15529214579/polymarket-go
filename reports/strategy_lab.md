# Strategy Lab Report

**Generated:** 2026-08-05 23:00 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (381)
- Valid strategies found: 165
- Candidate layers: 10 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 42 total
- Live-edge blocked push wallets: 9
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 10
- Aggregate closed copy trades: 228
- Aggregate copy ROI: 125.4%
- Aggregate copy PnL: $+3398.86
- Aggregate copy win rate: 81.1%
- Median wallet CopyROI: 134.8%
- Worst included wallet CopyROI: 108.4%
- Open copy cost / closed copy capital: 2.24x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 24.05 | 247 | 100% | 113.0% | $+892.77 | 60 | 83.3% | 36.8% |
| `0xb2ed...4418` | A | 23.11 | 126 | 39% | 134.4% | $+591.26 | 36 | 66.7% | 58.3% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 27.12 | 137 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 47.6% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.9% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 680
- Aggregate copy ROI: 86.6%
- Aggregate copy PnL: $+7744.97
- Aggregate copy win rate: 86.5%
- Worst included CopyROI: 44.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xbb35...b62a` | B | 100.0 | 33.5 | 760.0 | 1344 | 100% | 118.0% | $+2797.61 | 109 | 96.3% | $1137 | existing,holder |
| `0xadba...a6e2` | B | 100.0 | 31.3 | 671.7 | 24 | 77% | 206.3% | $+165.05 | 8 | 100.0% | $991 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0x17fe...b0ca` | B | 100.0 | 31.0 | 568.0 | 175 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $674 | existing |
| `0xb624...de17` | B | 100.0 | 26.2 | 563.7 | 28 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $795 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.5 | 525.0 | 922 | 93% | 59.4% | $+1241.49 | 190 | 89.5% | $293 | existing |
| `0x18c2...529a` | B | 100.0 | 28.2 | 490.2 | 853 | 74% | 70.3% | $+548.64 | 60 | 86.7% | $270 | existing,holder |
| `0x819d...6e9c` | B | 100.0 | 31.1 | 488.9 | 23 | 33% | 134.1% | $+107.24 | 8 | 100.0% | $421 | existing |
| `0x0ec9...1e0c` | A | 100.0 | 23.7 | 481.1 | 23 | 31% | 105.4% | $+147.61 | 14 | 57.1% | $576 | existing,holder |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.8 | 123 | 79% | 93.9% | $+178.48 | 16 | 75.0% | $447 | existing |
| `0x092b...614e` | B | 100.0 | 26.2 | 474.8 | 145 | 97% | 92.0% | $+165.61 | 17 | 82.3% | $1095 | existing,holder |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 24.0 | 452.1 | 56 | 34% | 59.8% | $+286.98 | 46 | 80.4% | $1288 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.5 | 417.0 | 26 | 37% | 81.2% | $+105.62 | 10 | 70.0% | $2177 | existing |
| `0x7673...fa40` | B | 100.0 | 28.8 | 405.0 | 93 | 63% | 76.6% | $+122.51 | 14 | 78.6% | $554 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x0e24...7014` | B | 100.0 | 29.8 | 397.1 | 153 | 98% | 44.3% | $+309.94 | 45 | 77.8% | $1074 | existing,sports_tape |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.9 | 621.0 | 386 | 274 | 40% | 104.9% | $+618.86 | 56 | 96.4% | 127.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 33.2 | 517.6 | 344 | 212 | 94% | 149.4% | $+209.20 | 9 | 77.8% | 139.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.6 | 379.9 | 125 | 71 | 68% | 71.2% | $+106.79 | 15 | 73.3% | 55.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 379.4 | 438 | 235 | 43% | 68.2% | $+136.50 | 18 | 83.3% | 63.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 21.4 | 375.0 | 72 | 58 | 67% | 54.0% | $+156.56 | 22 | 63.6% | 45.6% | existing | opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 371.9 | 61 | 55 | 38% | 62.3% | $+118.35 | 18 | 83.3% | 25.4% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3a6...71bd` | C | 100.0 | 36.7 | 220.1 | 405 | 31% | 0.0% | 0 | 390 | $778 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2929...1dd0` | C | 100.0 | 43.8 | 261.8 | 781 | 614 | 99% | 110.8% | 29 | 109.6% | 30 | $1133 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.7 | 226.0 | 450 | 204 | 98% | 19.6% | 69 | 19.6% | 69 | $195 | existing,holder | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 33.3 | 205.3 | 417 | 224 | 87% | 45.3% | 14 | 36.2% | 18 | $154 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |
| `0x91b9...889f` | B | 100.0 | 31.6 | 188.8 | 105 | 50 | 81% | 42.1% | 18 | 55.8% | 24 | $461 | existing | opposite_side_same_market |
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
| `0xa75b...772c` | A | 24.1 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 566.5 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.8 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 413.7 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 9
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 15m edge -2.63pp over 2 samples | 25.4% | 29 | existing | opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 33.5 | 15m edge -2.84pp over 3 samples | 118.0% | 109 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.1 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.1 | 1h edge -19.00pp over 1 samples | 134.4% | 36 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.2 | 1h edge -52.95pp over 1 samples | 70.3% | 60 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.8 | 1h edge -72.59pp over 1 samples | 44.3% | 45 | existing,sports_tape | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 5
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 24.1 | 732.4 | 247 | 193 | 100% | 113.0% | 60 | 113.0% | 60 | existing,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.8 | 607.6 | 781 | 614 | 99% | 110.8% | 29 | 109.6% | 30 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | blocked-edge | 1h edge -72.59pp over 1 samples | B | 100.0 | 29.8 | 296.6 | 153 | 145 | 98% | 42.6% | 43 | 44.3% | 45 | existing,sports_tape | opposite_side_same_market |
| `0x0c99...0c88` | watch | - | C | 100.0 | 44.5 | 25.2 | 365 | 142 | 36% | 0.0% | 0 | 0.0% | 0 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0xe907...cff6` | blocked-edge | 15m edge -1.55pp over 26 samples | BOT | 100.0 | 82.9 | -10.7 | 783 | 124 | 59% | 6.0% | 28 | 6.0% | 35 | existing,holder,leaderboard_volume_30d,leaderboard_volume_7d | bot_like_flow,burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 10 | 228 | 125.4% | $+3398.86 | 81.1% | 134.8% | 108.4% | 2.24x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 19 | 453 | 100.3% | $+5317.92 | 82.3% | 108.4% | 63.3% | 3.18x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 18 | 442 | 101.4% | $+5210.49 | 82.6% | 108.5% | 63.3% | 3.23x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 4 | 17 | 420 | 100.9% | $+4883.70 | 82.4% | 108.4% | 63.3% | 3.37x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 17 | 433 | 99.1% | $+4974.47 | 82.4% | 108.4% | 63.3% | 3.26x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 6 | 17 | 406 | 98.3% | $+4619.23 | 84.0% | 108.4% | 63.3% | 3.17x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 7 | 10 | 201 | 120.4% | $+2925.81 | 80.6% | 123.7% | 81.2% | 2.26x | tier>=A bot<25 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 18 | 444 | 98.1% | $+5081.90 | 82.2% | 101.2% | 63.3% | 3.21x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 18 | 417 | 97.3% | $+4726.66 | 83.7% | 101.2% | 63.3% | 3.12x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 22 | 563 | 89.7% | $+6011.80 | 81.0% | 93.0% | 44.3% | 2.71x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 21 | 554 | 90.7% | $+5929.70 | 81.2% | 93.9% | 44.3% | 2.77x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 23 | 596 | 87.7% | $+6210.29 | 80.0% | 92.0% | 44.3% | 2.62x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
