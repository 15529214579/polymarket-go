# Strategy Lab Report

**Generated:** 2026-10-22 23:14 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (430)
- Valid strategies found: 103
- Candidate layers: 19 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 56 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 19
- Aggregate closed copy trades: 1675
- Aggregate copy ROI: 191.8%
- Aggregate copy PnL: $+38016.49
- Aggregate copy win rate: 94.9%
- Median wallet CopyROI: 146.9%
- Worst included wallet CopyROI: 101.3%
- Open copy cost / closed copy capital: 0.58x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x45e4...3c69` | A | 22.37 | 34 | 4% | 167.9% | $+5256.74 | 273 | 91.6% | 28.5% |
| `0x42c9...9623` | A | 12.84 | 0 | 0% | 139.1% | $+2879.75 | 155 | 100.0% | 49.3% |
| `0x5135...7217` | A | 17.89 | 1199 | 100% | 141.1% | $+1777.68 | 97 | 100.0% | 58.5% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x15dd...95ba` | A | 24.93 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x61b6...e14e` | A | 22.60 | 10 | 6% | 129.6% | $+453.50 | 34 | 67.7% | 16.4% |
| `0x2b17...d36d` | A | 9.60 | 329 | 72% | 120.4% | $+445.45 | 26 | 96.2% | 52.9% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 204.8% | $+348.11 | 13 | 69.2% | 41.8% |
| `0x162d...8944` | A | 24.96 | 3 | 1% | 101.3% | $+344.59 | 33 | 90.9% | 47.7% |
| `0x7d68...8239` | A | 10.74 | 20 | 22% | 112.4% | $+292.19 | 26 | 96.2% | 31.8% |
| `0x7bdc...9b4c` | A | 24.39 | 61 | 37% | 225.2% | $+247.74 | 11 | 90.9% | 32.0% |
| `0xddf3...8d34` | A | 2.98 | 31 | 74% | 220.0% | $+220.02 | 8 | 100.0% | 22.7% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 146.9% | $+132.20 | 9 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.80 | 323 | 31% | 253.3% | $+12133.61 | 450 | 98.2% | 56.9% |
| `0x0f00...73ce` | B | 26.63 | 177 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.08 | 135 | 52% | 246.9% | $+3135.87 | 127 | 100.0% | 61.9% |
| `0x1bd8...e5f8` | B | 26.16 | 22 | 10% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x092b...614e` | B | 27.37 | 115 | 99% | 102.5% | $+225.51 | 22 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.73 | 63 | 62% | 114.6% | $+217.75 | 18 | 100.0% | 47.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 966
- Aggregate copy ROI: 95.3%
- Aggregate copy PnL: $+11557.38
- Aggregate copy win rate: 84.7%
- Worst included CopyROI: 45.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 98.2 | 14.2 | 1091.0 | 48 | 44% | 418.6% | $+293.01 | 7 | 100.0% | $300 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8a43...30da` | B | 96.7 | 25.2 | 887.3 | 21 | 88% | 255.3% | $+383.01 | 12 | 91.7% | $1050 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.2 | 822.5 | 63 | 22% | 165.0% | $+742.52 | 44 | 97.7% | $258 | existing |
| `0x760f...326a` | B | 100.0 | 30.5 | 794.9 | 15 | 11% | 126.0% | $+1284.66 | 91 | 93.4% | $8930 | existing |
| `0xbe7e...6d00` | B | 100.0 | 30.5 | 730.6 | 99 | 12% | 121.4% | $+959.12 | 73 | 95.9% | $1760 | existing |
| `0x3520...1fa0` | B | 100.0 | 30.0 | 720.7 | 34 | 31% | 164.5% | $+378.40 | 22 | 95.5% | $835 | existing |
| `0x612f...05d8` | B | 100.0 | 32.7 | 642.6 | 80 | 46% | 105.0% | $+651.00 | 59 | 94.9% | $600 | existing |
| `0x9a0c...370e` | B | 100.0 | 30.5 | 641.2 | 184 | 79% | 82.6% | $+1494.78 | 152 | 91.5% | $327 | existing,holder |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2bf0...728f` | A | 100.0 | 19.9 | 529.5 | 4 | 2% | 82.8% | $+480.19 | 56 | 85.7% | $1876 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x5925...9984` | B | 100.0 | 31.7 | 487.4 | 48 | 96% | 97.0% | $+174.66 | 15 | 80.0% | $1547 | existing |
| `0x7972...594f` | B | 100.0 | 33.1 | 481.7 | 285 | 67% | 55.9% | $+1028.76 | 144 | 50.7% | $1146 | existing |
| `0x789f...20ec` | B | 99.4 | 26.9 | 444.5 | 52 | 96% | 71.6% | $+171.80 | 18 | 94.4% | $1083 | existing |
| `0x3724...69b5` | A | 100.0 | 16.9 | 440.2 | 1161 | 94% | 76.9% | $+284.56 | 18 | 83.3% | $439 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.6 | 421.6 | 143 | 100% | 48.0% | $+297.65 | 60 | 91.7% | $946 | existing |
| `0x0e24...7014` | B | 100.0 | 30.6 | 414.5 | 204 | 91% | 45.9% | $+404.13 | 73 | 78.1% | $1210 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa6de...256c` | B | 100.0 | 30.5 | 1029.5 | 123 | 63 | 11% | 277.9% | $+805.75 | 18 | 100.0% | 273.6% | existing | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.3 | 756.4 | 820 | 350 | 92% | 122.3% | $+1835.24 | 85 | 100.0% | 118.2% | existing | burst_trading,open_copy_exposure |
| `0x677d...4907` | A | 100.0 | 21.5 | 593.8 | 150 | 122 | 50% | 88.4% | $+583.34 | 65 | 96.9% | 83.8% | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 570.0 | 1150 | 1114 | 100% | 60.9% | $+4288.40 | 296 | 96.3% | 60.9% | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.3 | 489.7 | 184 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.5 | 425.1 | 178 | 67 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | opposite_side_same_market |
| `0x620a...2200` | B | 100.0 | 25.3 | 384.7 | 60 | 20 | 26% | 97.9% | $+68.51 | 7 | 28.6% | 82.5% | existing | opposite_side_same_market |
| `0x5ba7...8732` | B | 100.0 | 30.8 | 367.2 | 375 | 133 | 54% | 35.2% | $+341.46 | 93 | 82.8% | 49.0% | existing | opposite_side_same_market |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 357.1 | 25 | 19 | 96% | 92.2% | $+46.10 | 5 | 80.0% | 92.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x26eb...22d0` | B | 100.0 | 31.9 | 344.2 | 144 | 45 | 14% | 95.0% | $+66.46 | 5 | 100.0% | 104.7% | existing | burst_trading,open_copy_exposure |

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
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.2 | 1073 | 585 | 98% | 57.4% | 163 | 58.0% | 168 | $606 | existing | opposite_side_same_market |
| `0x89eb...cff9` | C | 100.0 | 37.0 | 237.5 | 441 | 410 | 90% | 28.9% | 62 | 39.7% | 75 | $2084 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 36.4 | 232.3 | 1218 | 406 | 88% | 96.9% | 33 | 97.0% | 40 | $690 | existing,holder | burst_trading |
| `0x8e74...5972` | C | 100.0 | 38.4 | 225.1 | 1433 | 374 | 100% | 118.5% | 40 | 118.8% | 42 | $428 | existing | burst_trading |
| `0x257a...b42a` | B | 100.0 | 29.2 | 218.2 | 248 | 198 | 99% | 23.3% | 44 | 23.3% | 44 | $335 | existing | open_copy_exposure,opposite_side_same_market |
| `0x234f...6320` | C | 100.0 | 43.4 | 216.1 | 434 | 335 | 100% | 27.2% | 149 | 27.2% | 149 | $540 | existing,holder | burst_trading,opposite_side_same_market |
| `0x4b59...3aa6` | C | 100.0 | 37.5 | 210.5 | 1387 | 278 | 98% | 108.3% | 41 | 106.8% | 43 | $167 | existing,holder | burst_trading |
| `0x1f3a...895d` | C | 100.0 | 42.9 | 209.3 | 949 | 452 | 68% | 134.6% | 18 | 133.4% | 41 | $480 | existing | burst_trading |
| `0xbdd2...6657` | C | 100.0 | 40.6 | 208.2 | 1047 | 344 | 78% | 111.9% | 14 | 111.0% | 21 | $402 | existing,holder | burst_trading,open_copy_exposure |
| `0x73c2...3a29` | B | 100.0 | 33.5 | 207.4 | 465 | 277 | 85% | 13.6% | 139 | 15.2% | 159 | $158 | existing | opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 60.9% | 296 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.6 | 1h edge -72.59pp over 1 samples | 45.9% | 73 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3ce...a57a` | blocked-edge | 1h edge -47.28pp over 1 samples | B | 100.0 | 34.3 | 623.0 | 1150 | 1114 | 100% | 60.9% | 296 | 60.9% | 296 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0xd106...c79c` | watch | - | C | 100.0 | 40.1 | 544.0 | 345 | 195 | 26% | 106.5% | 20 | 101.6% | 93 | existing,holder,sports_tape | burst_trading |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 19 | 1675 | 191.8% | $+38016.49 | 94.9% | 146.9% | 101.3% | 0.58x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 17 | 1602 | 193.9% | $+36987.17 | 95.9% | 146.9% | 101.3% | 0.53x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 17 | 1628 | 192.8% | $+37214.88 | 95.6% | 146.9% | 101.3% | 0.54x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 15 | 1520 | 197.4% | $+35863.15 | 96.9% | 146.9% | 101.3% | 0.54x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 16 | 1589 | 193.8% | $+36639.06 | 96.2% | 144.0% | 101.3% | 0.52x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 6 | 17 | 1658 | 191.9% | $+37664.27 | 94.8% | 141.1% | 101.3% | 0.56x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 15 | 1585 | 193.9% | $+36634.95 | 95.9% | 141.1% | 101.3% | 0.51x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 15 | 1611 | 192.9% | $+36862.66 | 95.6% | 141.1% | 101.3% | 0.52x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 15 | 1634 | 191.6% | $+37068.42 | 95.0% | 139.1% | 101.3% | 0.54x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 13 | 1503 | 197.5% | $+35510.93 | 96.9% | 141.1% | 101.3% | 0.51x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 11 | 14 | 1572 | 193.8% | $+36286.84 | 96.1% | 140.1% | 101.3% | 0.49x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 12 | 14 | 1600 | 192.7% | $+36614.92 | 95.6% | 140.1% | 101.3% | 0.51x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
