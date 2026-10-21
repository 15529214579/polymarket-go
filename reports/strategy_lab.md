# Strategy Lab Report

**Generated:** 2026-10-21 23:23 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (430)
- Valid strategies found: 60
- Candidate layers: 17 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 17
- Aggregate closed copy trades: 1608
- Aggregate copy ROI: 193.2%
- Aggregate copy PnL: $+36841.36
- Aggregate copy win rate: 95.0%
- Median wallet CopyROI: 145.6%
- Worst included wallet CopyROI: 102.5%
- Open copy cost / closed copy capital: 0.60x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x45e4...3c69` | A | 22.48 | 34 | 4% | 167.2% | $+5268.53 | 275 | 91.6% | 28.5% |
| `0x42c9...9623` | A | 12.73 | 0 | 0% | 141.3% | $+2896.52 | 153 | 100.0% | 49.3% |
| `0x5135...7217` | A | 17.89 | 1206 | 100% | 145.6% | $+1761.43 | 92 | 100.0% | 58.5% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x61b6...e14e` | A | 23.42 | 9 | 7% | 131.9% | $+395.70 | 29 | 72.4% | 18.9% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 47.7% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 195.0% | $+312.07 | 12 | 66.7% | 41.8% |
| `0x7d68...8239` | A | 10.58 | 20 | 23% | 113.6% | $+249.84 | 22 | 95.5% | 31.8% |
| `0x2b17...d36d` | A | 9.60 | 338 | 72% | 128.0% | $+204.75 | 9 | 77.8% | 52.9% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 146.9% | $+132.20 | 9 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.79 | 322 | 30% | 250.1% | $+11855.11 | 445 | 98.0% | 56.9% |
| `0x0f00...73ce` | B | 26.63 | 177 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.12 | 136 | 53% | 247.1% | $+3113.93 | 126 | 100.0% | 61.6% |
| `0x15dd...95ba` | B | 25.00 | 196 | 48% | 117.3% | $+774.21 | 63 | 84.1% | 59.0% |
| `0x1bd8...e5f8` | B | 25.74 | 22 | 10% | 147.9% | $+547.10 | 37 | 73.0% | 19.3% |
| `0x092b...614e` | B | 27.37 | 115 | 99% | 102.5% | $+225.51 | 22 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.72 | 62 | 63% | 114.6% | $+217.75 | 18 | 100.0% | 47.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 787
- Aggregate copy ROI: 99.8%
- Aggregate copy PnL: $+9704.08
- Aggregate copy win rate: 90.6%
- Worst included CopyROI: 43.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 98.3 | 14.1 | 1091.7 | 50 | 44% | 418.6% | $+293.01 | 7 | 100.0% | $330 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.0 | 834.3 | 61 | 22% | 169.0% | $+743.75 | 43 | 97.7% | $258 | existing |
| `0x760f...326a` | B | 100.0 | 30.4 | 794.7 | 15 | 11% | 126.0% | $+1284.66 | 91 | 93.4% | $8873 | existing |
| `0x3520...1fa0` | B | 100.0 | 32.0 | 726.0 | 32 | 33% | 172.1% | $+361.41 | 20 | 95.0% | $836 | existing |
| `0xbe7e...6d00` | B | 100.0 | 30.4 | 709.4 | 100 | 12% | 126.5% | $+720.75 | 52 | 94.2% | $1708 | existing |
| `0x612f...05d8` | B | 100.0 | 33.0 | 636.0 | 75 | 46% | 105.3% | $+610.90 | 55 | 94.5% | $602 | existing |
| `0x9a0c...370e` | B | 100.0 | 30.6 | 635.9 | 181 | 78% | 82.1% | $+1444.96 | 147 | 91.8% | $322 | existing,holder |
| `0x2089...712e` | B | 100.0 | 28.8 | 591.0 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6314 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2bf0...728f` | A | 100.0 | 19.9 | 529.5 | 4 | 2% | 82.8% | $+480.19 | 56 | 85.7% | $1876 | existing |
| `0x68cb...6f57` | B | 100.0 | 17.5 | 492.8 | 0 | 0% | 199.8% | $+99.91 | 4 | 100.0% | $394 | existing |
| `0x5925...9984` | B | 100.0 | 31.7 | 487.4 | 48 | 96% | 97.0% | $+174.66 | 15 | 80.0% | $1547 | existing |
| `0x789f...20ec` | B | 99.4 | 26.9 | 444.5 | 52 | 96% | 71.6% | $+171.80 | 18 | 94.4% | $1083 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.6 | 420.7 | 143 | 100% | 48.0% | $+293.09 | 59 | 91.5% | $946 | existing |
| `0x0e24...7014` | B | 100.0 | 30.8 | 416.6 | 201 | 90% | 46.4% | $+408.49 | 73 | 78.1% | $1223 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |
| `0x0d23...2d56` | B | 100.0 | 26.2 | 353.9 | 33 | 18% | 43.9% | $+92.21 | 21 | 95.2% | $4648 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.4 | 747.5 | 841 | 358 | 92% | 125.5% | $+1656.19 | 73 | 100.0% | 122.7% | existing | burst_trading,open_copy_exposure |
| `0x677d...4907` | A | 100.0 | 21.4 | 584.9 | 144 | 116 | 49% | 87.7% | $+543.88 | 61 | 96.7% | 83.9% | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 577.3 | 1150 | 1112 | 100% | 62.1% | $+4382.20 | 298 | 96.3% | 62.1% | existing,holder | burst_trading,opposite_side_same_market |
| `0x7972...594f` | B | 100.0 | 33.2 | 495.8 | 276 | 262 | 67% | 66.6% | $+825.43 | 90 | 53.3% | 53.8% | existing | opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.4 | 489.6 | 182 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.5 | 425.1 | 178 | 67 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | opposite_side_same_market |
| `0x620a...2200` | B | 100.0 | 25.3 | 384.7 | 60 | 20 | 26% | 97.9% | $+68.51 | 7 | 28.6% | 82.5% | existing | opposite_side_same_market |
| `0x5ba7...8732` | B | 100.0 | 30.8 | 362.8 | 373 | 132 | 54% | 34.4% | $+327.14 | 91 | 83.5% | 48.4% | existing | opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.3 | 331.4 | 319 | 123 | 29% | 51.7% | $+242.86 | 20 | 40.0% | 68.4% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x3f3d...2e18` | B | 100.0 | 34.0 | 320.8 | 72 | 41 | 70% | 79.5% | $+111.27 | 5 | 80.0% | 87.5% | existing | burst_trading,open_copy_exposure |

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
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.3 | 1081 | 585 | 98% | 57.4% | 163 | 58.0% | 168 | $602 | existing | opposite_side_same_market |
| `0xc9bf...81c0` | C | 100.0 | 42.0 | 241.2 | 1385 | 603 | 97% | 101.2% | 12 | 101.2% | 12 | $175 | existing,holder | burst_trading,open_copy_exposure |
| `0x89eb...cff9` | C | 100.0 | 37.0 | 239.7 | 436 | 405 | 90% | 18.9% | 49 | 33.1% | 62 | $2062 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 36.1 | 233.4 | 1189 | 413 | 86% | 336.9% | 5 | 299.0% | 6 | $744 | existing,holder | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 38.4 | 225.1 | 1433 | 374 | 100% | 118.5% | 40 | 118.8% | 42 | $428 | existing | burst_trading |
| `0x257a...b42a` | B | 100.0 | 29.0 | 219.3 | 261 | 205 | 99% | 21.3% | 42 | 21.3% | 42 | $326 | existing | open_copy_exposure,opposite_side_same_market |
| `0x234f...6320` | C | 100.0 | 43.5 | 218.3 | 410 | 314 | 100% | 24.9% | 140 | 24.9% | 140 | $529 | existing,sports_tape | burst_trading,opposite_side_same_market |
| `0xbdd2...6657` | C | 100.0 | 40.2 | 212.2 | 1044 | 369 | 78% | 100.2% | 6 | 108.2% | 8 | $423 | existing,holder | burst_trading,open_copy_exposure |
| `0x73c2...3a29` | B | 100.0 | 33.4 | 207.2 | 463 | 275 | 85% | 13.4% | 138 | 15.0% | 158 | $157 | existing | opposite_side_same_market |
| `0x168a...0d24` | C | 100.0 | 39.9 | 198.5 | 349 | 200 | 47% | 156.1% | 9 | 136.6% | 18 | $18041 | existing | burst_trading,open_copy_exposure |

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

- Wallets: 3
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | B | 100.0 | 26.6 | 15m edge -10.00pp over 2 samples | 225.7% | 195 | existing | burst_trading,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 62.1% | 298 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.8 | 1h edge -72.59pp over 1 samples | 46.4% | 73 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x234f...6320` | pushed | - | C | 100.0 | 43.5 | 229.3 | 410 | 314 | 100% | 24.9% | 140 | 24.9% | 140 | existing,sports_tape | burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 17 | 1608 | 193.2% | $+36841.36 | 95.0% | 145.6% | 102.5% | 0.60x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 16 | 1596 | 193.2% | $+36529.29 | 95.2% | 143.4% | 102.5% | 0.58x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 15 | 1542 | 195.1% | $+35898.56 | 96.0% | 145.6% | 102.5% | 0.55x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 15 | 1590 | 194.0% | $+36504.41 | 95.1% | 145.6% | 102.5% | 0.47x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 14 | 1530 | 195.1% | $+35586.49 | 96.2% | 143.4% | 102.5% | 0.53x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 6 | 14 | 1578 | 194.0% | $+36192.34 | 95.3% | 143.4% | 102.5% | 0.45x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 13 | 1524 | 195.9% | $+35561.61 | 96.1% | 145.6% | 102.5% | 0.41x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 13 | 1521 | 195.7% | $+35381.74 | 96.3% | 145.6% | 102.5% | 0.40x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 13 | 1560 | 194.8% | $+35974.59 | 95.3% | 145.6% | 102.5% | 0.42x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 23 | 1905 | 175.3% | $+40117.10 | 94.2% | 128.0% | 77.1% | 0.68x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 12 | 1512 | 195.9% | $+35249.54 | 96.3% | 143.4% | 102.5% | 0.39x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 12 | 11 | 1494 | 196.8% | $+35031.79 | 96.3% | 145.6% | 102.5% | 0.36x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
