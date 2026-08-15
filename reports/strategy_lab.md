# Strategy Lab Report

**Generated:** 2026-08-15 23:36 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (403)
- Valid strategies found: 125
- Candidate layers: 14 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 48 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 14
- Aggregate closed copy trades: 295
- Aggregate copy ROI: 116.1%
- Aggregate copy PnL: $+3947.33
- Aggregate copy win rate: 83.7%
- Median wallet CopyROI: 120.8%
- Worst included wallet CopyROI: 101.3%
- Open copy cost / closed copy capital: 2.55x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.27 | 401 | 100% | 101.3% | $+1114.74 | 91 | 80.2% | 41.7% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.60 | 123 | 82% | 102.8% | $+185.06 | 15 | 80.0% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x5be9...ef2b` | A | 14.84 | 51 | 63% | 115.9% | $+104.32 | 9 | 100.0% | 24.0% |
| `0x44c4...09cb` | B | 26.98 | 133 | 63% | 107.2% | $+503.66 | 43 | 76.7% | 55.3% |
| `0x162d...8944` | B | 25.00 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.47 | 385 | 74% | 134.6% | $+174.98 | 8 | 62.5% | 53.4% |
| `0x819d...6e9c` | B | 28.58 | 26 | 29% | 125.6% | $+138.21 | 10 | 100.0% | 47.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 461
- Aggregate copy ROI: 75.2%
- Aggregate copy PnL: $+4227.91
- Aggregate copy win rate: 82.9%
- Worst included CopyROI: 48.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3c2b...825e` | B | 100.0 | 33.6 | 582.6 | 474 | 95% | 129.4% | $+310.68 | 19 | 89.5% | $250 | existing,holder |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x18c2...529a` | B | 100.0 | 27.9 | 499.4 | 865 | 75% | 73.0% | $+569.27 | 60 | 88.3% | $287 | existing |
| `0xa35c...a113` | B | 100.0 | 32.6 | 499.3 | 156 | 100% | 76.6% | $+329.36 | 39 | 94.9% | $696 | existing,holder,sports_tape |
| `0x092b...614e` | B | 100.0 | 25.6 | 470.3 | 135 | 96% | 93.2% | $+158.52 | 16 | 81.2% | $1082 | existing |
| `0xb99c...7e96` | A | 100.0 | 20.1 | 466.4 | 104 | 89% | 83.7% | $+200.81 | 18 | 94.4% | $526 | existing,holder |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 423.3 | 187 | 98% | 49.5% | $+420.55 | 57 | 80.7% | $1103 | existing,holder |
| `0x760f...326a` | B | 100.0 | 31.0 | 413.9 | 15 | 14% | 84.4% | $+135.00 | 13 | 76.9% | $5078 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.4 | 412.0 | 21 | 33% | 83.3% | $+99.92 | 9 | 66.7% | $1884 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0xa931...1387` | B | 73.3 | 29.3 | 392.6 | 8 | 73% | 145.1% | $+43.53 | 3 | 100.0% | $2545 | existing,holder |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x0931...e78e` | B | 100.0 | 25.8 | 373.9 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | $298 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.8 | 781.6 | 378 | 285 | 37% | 135.9% | $+978.32 | 69 | 97.1% | 157.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 32.5 | 650.5 | 1362 | 364 | 100% | 141.7% | $+963.36 | 24 | 100.0% | 141.7% | existing,holder | burst_trading,open_copy_exposure |
| `0x21cc...54bc` | B | 100.0 | 28.7 | 424.6 | 170 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.9 | 406.0 | 513 | 270 | 50% | 80.9% | $+137.53 | 15 | 93.3% | 68.3% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 400.0 | 1179 | 1127 | 100% | 54.5% | $+664.68 | 51 | 88.2% | 54.5% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 333.3 | 422 | 232 | 75% | 39.0% | $+225.89 | 50 | 68.0% | 34.8% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 44.5 | 206.5 | 1181 | 85% | 0.0% | 0 | 402 | $533 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2929...1dd0` | C | 100.0 | 43.7 | 283.0 | 929 | 723 | 98% | 107.9% | 38 | 107.1% | 39 | $1228 | existing,holder,recent_trade,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x6ac5...4b6e` | C | 100.0 | 37.4 | 245.1 | 1071 | 281 | 77% | 122.5% | 16 | 121.1% | 20 | $355 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.2 | 221.1 | 271 | 171 | 85% | 12.5% | 50 | 11.7% | 57 | $669 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 34.8 | 214.9 | 653 | 315 | 91% | 47.5% | 21 | 40.5% | 25 | $145 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 214.3 | 577 | 351 | 88% | 28.0% | 78 | 26.4% | 85 | $226 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | C | 100.0 | 35.1 | 206.2 | 900 | 335 | 72% | 147.9% | 25 | 140.9% | 38 | $137 | existing,holder | burst_trading,open_copy_exposure |
| `0x2fbb...7e44` | B | 100.0 | 27.8 | 201.1 | 157 | 70 | 100% | 37.9% | 19 | 37.9% | 19 | $165 | existing | opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 23.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 561.0 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.7 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 429.6 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.0 | 1h edge -15.87pp over 1 samples | 107.2% | 43 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 1h edge -47.28pp over 1 samples | 54.5% | 51 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 73.0% | 60 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 49.5% | 57 | existing,holder | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 3
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.3 | 732.5 | 401 | 287 | 100% | 101.3% | 91 | 101.3% | 91 | existing,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.7 | 642.3 | 929 | 723 | 98% | 107.9% | 38 | 107.1% | 39 | existing,holder,recent_trade,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa35c...a113` | pushed | - | B | 100.0 | 32.6 | 459.2 | 156 | 155 | 100% | 76.6% | 39 | 76.6% | 39 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 14 | 295 | 116.1% | $+3947.33 | 83.7% | 120.8% | 101.3% | 2.55x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 21 | 465 | 100.5% | $+5438.18 | 84.5% | 107.2% | 62.1% | 3.51x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 447 | 101.3% | $+5237.37 | 84.1% | 107.8% | 62.1% | 3.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 26 | 579 | 90.6% | $+6126.65 | 83.6% | 102.1% | 41.8% | 3.19x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 5 | 27 | 612 | 88.3% | $+6319.33 | 82.5% | 101.3% | 41.8% | 3.04x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 19 | 445 | 101.8% | $+5263.74 | 85.2% | 108.4% | 65.3% | 3.50x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 25 | 561 | 90.9% | $+5925.84 | 83.2% | 102.8% | 41.8% | 3.23x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 22 | 540 | 93.9% | $+5790.62 | 84.4% | 105.6% | 41.9% | 3.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 13 | 287 | 115.4% | $+3772.35 | 84.3% | 115.9% | 101.3% | 1.79x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 21 | 522 | 94.3% | $+5589.81 | 84.1% | 107.2% | 41.9% | 3.17x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 11 | 18 | 427 | 102.7% | $+5062.93 | 84.8% | 108.7% | 65.3% | 3.55x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 12 | 19 | 448 | 100.1% | $+5163.28 | 85.3% | 107.2% | 62.1% | 3.07x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
