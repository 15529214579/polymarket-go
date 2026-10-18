# Strategy Lab Report

**Generated:** 2026-10-18 23:23 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (430)
- Valid strategies found: 70
- Candidate layers: 17 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 17
- Aggregate closed copy trades: 1486
- Aggregate copy ROI: 188.0%
- Aggregate copy PnL: $+41449.24
- Aggregate copy win rate: 97.0%
- Median wallet CopyROI: 146.9%
- Worst included wallet CopyROI: 102.5%
- Open copy cost / closed copy capital: 0.90x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xedc3...f923` | A | 23.82 | 844 | 77% | 172.2% | $+9503.91 | 211 | 99.0% | 49.1% |
| `0x42c9...9623` | A | 13.03 | 0 | 0% | 145.7% | $+3146.34 | 155 | 100.0% | 49.3% |
| `0xc9bf...81c0` | A | 24.07 | 1378 | 99% | 168.3% | $+2204.77 | 19 | 100.0% | 41.6% |
| `0x5135...7217` | A | 18.63 | 1193 | 100% | 151.5% | $+1817.68 | 89 | 100.0% | 58.5% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x612f...05d8` | A | 24.93 | 67 | 45% | 106.7% | $+586.62 | 52 | 94.2% | 53.6% |
| `0x1bd8...e5f8` | A | 23.63 | 22 | 12% | 141.5% | $+509.49 | 36 | 72.2% | 19.3% |
| `0x162d...8944` | A | 24.81 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 47.7% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 195.0% | $+312.07 | 12 | 66.7% | 41.8% |
| `0x2b17...d36d` | A | 17.50 | 383 | 72% | 119.9% | $+203.76 | 10 | 70.0% | 52.9% |
| `0x7d68...8239` | A | 10.25 | 20 | 24% | 108.6% | $+184.55 | 17 | 94.1% | 31.8% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 146.9% | $+132.20 | 9 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.88 | 331 | 31% | 235.5% | $+10551.08 | 422 | 98.1% | 56.9% |
| `0x0f00...73ce` | B | 26.64 | 176 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.21 | 135 | 53% | 250.5% | $+2980.81 | 119 | 100.0% | 61.3% |
| `0x092b...614e` | B | 27.48 | 114 | 99% | 102.5% | $+225.51 | 22 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.35 | 60 | 62% | 113.2% | $+203.74 | 17 | 100.0% | 47.4% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 785
- Aggregate copy ROI: 128.4%
- Aggregate copy PnL: $+13557.01
- Aggregate copy win rate: 84.5%
- Worst included CopyROI: 43.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xae2b...90ea` | B | 100.0 | 34.6 | 1448.9 | 187 | 35% | 232.9% | $+6335.38 | 192 | 93.8% | $1896 | existing |
| `0x3e35...4f69` | A | 98.5 | 13.7 | 1093.0 | 51 | 41% | 418.6% | $+293.01 | 7 | 100.0% | $367 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x7e33...73e7` | B | 100.0 | 33.9 | 821.1 | 61 | 22% | 164.8% | $+741.47 | 44 | 95.5% | $258 | existing |
| `0x760f...326a` | B | 100.0 | 30.5 | 793.7 | 15 | 10% | 126.0% | $+1284.66 | 91 | 93.4% | $8749 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 591.0 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6314 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x7972...594f` | B | 100.0 | 33.2 | 493.7 | 249 | 67% | 59.6% | $+947.95 | 124 | 50.0% | $1111 | existing,holder |
| `0x68cb...6f57` | B | 100.0 | 17.9 | 492.5 | 0 | 0% | 199.8% | $+99.91 | 4 | 100.0% | $407 | existing |
| `0x5925...9984` | B | 100.0 | 32.0 | 490.6 | 48 | 98% | 97.0% | $+174.66 | 15 | 80.0% | $1559 | existing,holder |
| `0xbb4d...de5a` | A | 100.0 | 15.2 | 451.1 | 8 | 36% | 129.8% | $+64.90 | 5 | 100.0% | $7419 | existing |
| `0x789f...20ec` | B | 99.4 | 26.9 | 444.5 | 52 | 96% | 71.6% | $+171.80 | 18 | 94.4% | $1083 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.7 | 422.4 | 139 | 100% | 48.7% | $+292.33 | 58 | 91.4% | $869 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 416.6 | 201 | 91% | 46.4% | $+408.49 | 73 | 78.1% | $1215 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |
| `0x0d23...2d56` | B | 100.0 | 26.2 | 353.6 | 31 | 18% | 43.9% | $+92.16 | 21 | 95.2% | $4388 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 346.7 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing |
| `0xaabc...a190` | B | 100.0 | 5.0 | 345.6 | 0 | 0% | 106.0% | $+31.81 | 3 | 100.0% | $298 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 33.9 | 747.2 | 904 | 374 | 93% | 125.5% | $+1656.19 | 73 | 100.0% | 122.7% | existing | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 588.2 | 1149 | 1112 | 100% | 64.3% | $+4372.80 | 292 | 96.6% | 64.3% | existing,holder | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.6 | 489.3 | 174 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.5 | 425.1 | 178 | 67 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.3 | 331.4 | 319 | 123 | 29% | 51.7% | $+242.86 | 20 | 40.0% | 68.4% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x3f3d...2e18` | B | 100.0 | 34.0 | 320.8 | 72 | 41 | 70% | 79.5% | $+111.27 | 5 | 80.0% | 87.5% | existing | burst_trading,open_copy_exposure |
| `0x1944...8088` | A | 100.0 | 23.4 | 315.4 | 39 | 26 | 85% | 48.9% | $+78.15 | 10 | 60.0% | 45.7% | existing | opposite_side_same_market |
| `0x73c2...3a29` | B | 100.0 | 33.6 | 261.2 | 446 | 268 | 85% | 13.4% | $+182.34 | 135 | 73.3% | 15.1% | existing | opposite_side_same_market |
| `0x257a...b42a` | B | 100.0 | 28.9 | 249.4 | 265 | 208 | 99% | 21.3% | $+97.89 | 42 | 54.8% | 21.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | B | 100.0 | 25.2 | 235.2 | 241 | 74 | 60% | 16.7% | $+38.43 | 21 | 61.9% | 16.1% | existing | opposite_side_same_market |

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
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.4 | 1085 | 586 | 98% | 57.4% | 163 | 58.0% | 168 | $600 | existing | opposite_side_same_market |
| `0xd570...e4f8` | C | 100.0 | 43.6 | 244.0 | 1451 | 537 | 100% | 73.5% | 17 | 73.5% | 17 | $1466 | existing | burst_trading |
| `0xfdc0...4b85` | C | 100.0 | 35.9 | 236.4 | 1264 | 458 | 100% | 132.2% | 38 | 132.2% | 38 | $172 | existing,holder | burst_trading,open_copy_exposure |
| `0x89eb...cff9` | C | 100.0 | 37.4 | 234.5 | 409 | 378 | 90% | 16.2% | 47 | 30.4% | 59 | $1662 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 36.8 | 228.6 | 1202 | 410 | 87% | 336.9% | 5 | 299.0% | 6 | $691 | existing | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 38.4 | 225.1 | 1433 | 374 | 100% | 118.5% | 40 | 118.8% | 42 | $428 | existing | burst_trading |
| `0x5674...8b14` | C | 100.0 | 43.9 | 225.0 | 872 | 573 | 78% | 154.0% | 217 | 156.1% | 286 | $216 | existing,holder | burst_trading,opposite_side_same_market |
| `0xc34c...2e29` | C | 100.0 | 41.9 | 217.3 | 605 | 510 | 66% | 69.8% | 419 | 71.1% | 593 | $478 | existing | opposite_side_same_market |
| `0x4b59...3aa6` | C | 100.0 | 36.5 | 213.1 | 1385 | 306 | 98% | 102.6% | 46 | 101.5% | 48 | $182 | existing | burst_trading |
| `0xbdd2...6657` | C | 100.0 | 40.3 | 209.7 | 1028 | 357 | 77% | 99.6% | 8 | 105.2% | 10 | $397 | existing,holder | burst_trading,open_copy_exposure |

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
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 1h edge -47.28pp over 1 samples | 64.3% | 292 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.4% | 73 | existing | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 17 | 1486 | 188.0% | $+41449.24 | 97.0% | 146.9% | 102.5% | 0.90x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 21 | 1585 | 180.5% | $+43042.34 | 96.6% | 141.5% | 76.1% | 0.93x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 1577 | 180.8% | $+42981.46 | 96.6% | 143.6% | 77.1% | 0.93x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 16 | 1450 | 188.7% | $+40939.75 | 97.7% | 149.2% | 102.5% | 0.88x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 16 | 1477 | 188.1% | $+41317.04 | 97.0% | 148.6% | 102.5% | 0.89x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 19 | 1541 | 181.4% | $+42471.97 | 97.1% | 145.7% | 77.1% | 0.92x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 15 | 1441 | 188.9% | $+40807.55 | 97.6% | 151.5% | 102.5% | 0.88x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 16 | 1474 | 187.9% | $+41137.17 | 97.3% | 146.3% | 102.5% | 0.89x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 19 | 1568 | 181.0% | $+42849.26 | 96.6% | 141.5% | 77.1% | 0.93x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 19 | 1565 | 180.7% | $+42669.39 | 96.8% | 141.5% | 77.1% | 0.92x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 15 | 1438 | 188.7% | $+40627.68 | 97.9% | 146.9% | 102.5% | 0.87x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 12 | 15 | 1465 | 188.1% | $+41004.97 | 97.3% | 145.7% | 102.5% | 0.88x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
