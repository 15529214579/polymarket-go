# Strategy Lab Report

**Generated:** 2026-10-24 23:00 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (433)
- Valid strategies found: 195
- Candidate layers: 23 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 62 total
- Live-edge blocked push wallets: 2
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 2310
- Aggregate copy ROI: 174.3%
- Aggregate copy PnL: $+49498.45
- Aggregate copy win rate: 95.7%
- Median wallet CopyROI: 138.0%
- Worst included wallet CopyROI: 101.1%
- Open copy cost / closed copy capital: 0.36x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.42 | 32 | 4% | 171.1% | $+5818.25 | 298 | 92.3% | 28.5% |
| `0x3724...69b5` | A | 16.92 | 1118 | 94% | 122.1% | $+4530.25 | 207 | 98.5% | 48.9% |
| `0x2b17...d36d` | A | 9.76 | 303 | 70% | 171.8% | $+3865.25 | 153 | 99.3% | 53.2% |
| `0x42c9...9623` | A | 12.56 | 0 | 0% | 127.5% | $+2574.95 | 154 | 100.0% | 49.3% |
| `0x677d...4907` | A | 21.48 | 149 | 50% | 108.5% | $+2300.37 | 209 | 97.6% | 54.3% |
| `0x075d...58c1` | A | 19.49 | 18 | 17% | 166.5% | $+1432.08 | 78 | 91.0% | 40.5% |
| `0x61b6...e14e` | A | 23.18 | 10 | 4% | 138.0% | $+814.18 | 58 | 72.4% | 18.2% |
| `0x15dd...95ba` | A | 24.93 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0xddf3...8d34` | A | 2.84 | 33 | 75% | 138.1% | $+386.53 | 22 | 100.0% | 14.6% |
| `0x768b...bf0f` | A | 2.83 | 14 | 18% | 103.3% | $+299.49 | 24 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.68 | 64 | 38% | 166.7% | $+266.69 | 16 | 93.8% | 32.0% |
| `0x1425...6e0c` | A | 23.87 | 176 | 64% | 106.3% | $+244.51 | 23 | 100.0% | 63.4% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.59 | 332 | 31% | 249.2% | $+12159.97 | 457 | 98.0% | 56.9% |
| `0x5be9...ef2b` | B | 27.22 | 135 | 53% | 241.2% | $+3207.39 | 133 | 100.0% | 61.9% |
| `0x7967...f4a0` | B | 29.57 | 36 | 49% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.21 | 22 | 10% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x162d...8944` | B | 25.10 | 3 | 1% | 125.3% | $+475.98 | 37 | 91.9% | 47.3% |
| `0x092b...614e` | B | 26.55 | 111 | 98% | 103.8% | $+373.82 | 33 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.86 | 64 | 63% | 101.1% | $+363.82 | 34 | 100.0% | 47.6% |
| `0xbe12...728f` | B | 26.44 | 22 | 27% | 104.3% | $+198.14 | 19 | 100.0% | 52.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1312
- Aggregate copy ROI: 101.9%
- Aggregate copy PnL: $+16286.46
- Aggregate copy win rate: 89.8%
- Worst included CopyROI: 61.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 98.3 | 13.9 | 1253.4 | 47 | 44% | 379.1% | $+606.57 | 14 | 100.0% | $303 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8a43...30da` | A | 97.5 | 24.8 | 872.2 | 22 | 88% | 241.7% | $+386.75 | 13 | 92.3% | $1008 | existing |
| `0xbe7e...6d00` | B | 100.0 | 30.4 | 855.6 | 96 | 12% | 116.5% | $+2946.91 | 230 | 98.3% | $1824 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 835.7 | 15 | 11% | 131.3% | $+1483.87 | 102 | 94.1% | $9118 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.2 | 804.2 | 63 | 23% | 146.7% | $+894.95 | 60 | 98.3% | $258 | existing |
| `0x3520...1fa0` | B | 100.0 | 30.0 | 801.4 | 34 | 31% | 175.9% | $+510.03 | 28 | 96.4% | $835 | existing |
| `0x612f...05d8` | B | 100.0 | 32.6 | 656.2 | 86 | 43% | 103.2% | $+763.95 | 70 | 94.3% | $585 | existing,holder |
| `0x9a0c...370e` | B | 100.0 | 30.5 | 644.5 | 200 | 79% | 81.8% | $+1602.83 | 165 | 92.1% | $357 | existing,holder |
| `0x50d0...a6fc` | A | 100.0 | 18.1 | 603.0 | 30 | 18% | 72.2% | $+952.68 | 116 | 99.1% | $2440 | existing,sports_tape |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x7d68...8239` | A | 100.0 | 10.7 | 595.0 | 20 | 22% | 94.5% | $+396.86 | 42 | 97.6% | $287 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x2bf0...728f` | A | 100.0 | 20.0 | 569.6 | 10 | 5% | 87.7% | $+561.11 | 62 | 85.5% | $1826 | existing |
| `0x7972...594f` | B | 100.0 | 33.1 | 564.7 | 300 | 68% | 66.9% | $+1691.91 | 195 | 61.5% | $1159 | existing |
| `0x5925...9984` | B | 100.0 | 31.0 | 518.2 | 52 | 95% | 105.8% | $+200.96 | 16 | 81.2% | $1537 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x0d23...2d56` | B | 100.0 | 34.2 | 467.1 | 35 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5561 | existing |
| `0x6261...4adb` | A | 100.0 | 7.7 | 460.0 | 120 | 90% | 83.5% | $+292.35 | 12 | 91.7% | $1268 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.8 | 786.5 | 772 | 332 | 92% | 104.8% | $+3469.85 | 208 | 100.0% | 102.3% | existing | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.3 | 753.2 | 1179 | 320 | 100% | 97.0% | $+2473.17 | 188 | 100.0% | 97.0% | existing,holder | burst_trading |
| `0x8987...02f3` | B | 100.0 | 34.3 | 517.9 | 290 | 116 | 30% | 77.6% | $+814.62 | 53 | 79.2% | 89.4% | existing | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.1 | 490.0 | 187 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.4 | 477.3 | 146 | 75 | 79% | 71.6% | $+300.63 | 40 | 90.0% | 76.8% | existing | opposite_side_same_market |
| `0xeb2e...9b53` | B | 100.0 | 29.0 | 445.4 | 166 | 51 | 50% | 76.4% | $+359.20 | 20 | 90.0% | 62.7% | existing,holder | opposite_side_same_market |
| `0xf902...a267` | B | 100.0 | 26.8 | 427.9 | 21 | 18 | 84% | 89.2% | $+98.07 | 11 | 100.0% | 69.7% | existing,sports_tape | opposite_side_same_market |
| `0x4da5...3baf` | B | 100.0 | 31.1 | 424.7 | 148 | 121 | 100% | 49.7% | $+373.05 | 73 | 93.2% | 49.7% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.6 | 424.5 | 182 | 69 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.7 | 422.5 | 207 | 197 | 91% | 48.2% | $+472.73 | 86 | 81.4% | 50.9% | existing | opposite_side_same_market |

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
| `0x67ba...246c` | C | 100.0 | 39.4 | 249.5 | 986 | 577 | 98% | 58.4% | 168 | 59.0% | 173 | $654 | existing | opposite_side_same_market |
| `0x249e...4fa6` | C | 100.0 | 36.4 | 237.8 | 1168 | 417 | 100% | 59.9% | 180 | 59.9% | 180 | $151 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0x1550...5481` | C | 100.0 | 43.7 | 233.9 | 1451 | 493 | 100% | 174.4% | 18 | 174.4% | 18 | $244 | existing,sports_tape | burst_trading,open_copy_exposure |
| `0x4d6e...5cfd` | C | 100.0 | 35.8 | 232.1 | 1207 | 425 | 88% | 92.5% | 48 | 91.2% | 59 | $645 | existing | burst_trading |
| `0x3fe2...93b3` | A | 100.0 | 20.3 | 231.9 | 75 | 58 | 100% | 19.3% | 20 | 19.3% | 20 | $1579 | existing,holder | opposite_side_same_market |
| `0x3a01...c89e` | C | 100.0 | 43.1 | 231.8 | 1322 | 483 | 99% | 91.7% | 98 | 91.7% | 98 | $197 | existing,sports_tape | burst_trading |
| `0x5f61...519a` | B | 100.0 | 27.1 | 227.0 | 429 | 254 | 100% | 37.4% | 100 | 37.4% | 100 | $235 | existing | opposite_side_same_market |
| `0x89eb...cff9` | C | 100.0 | 44.2 | 224.9 | 456 | 424 | 88% | 59.1% | 162 | 61.9% | 177 | $2019 | existing | burst_trading,opposite_side_same_market |
| `0x684b...8409` | C | 100.0 | 35.9 | 216.0 | 1226 | 297 | 83% | 92.7% | 8 | 92.1% | 14 | $548 | existing,holder | burst_trading |
| `0x6ac5...4b6e` | C | 100.0 | 41.0 | 215.4 | 1469 | 297 | 100% | 162.7% | 12 | 162.7% | 12 | $500 | existing,holder | burst_trading |

## Recent Trade Flow Scout

- Wallets: 0
- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No recent-trade wallets passed the flow filters.

## Sports Tape Hotlist

- Wallets: 1
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x81f0...cbb1` | A | 100.0 | 2.9 | 329.4 | 27 | 26 | 37.4% | 13 | 41.2% | 17 | $1798 | existing,sports_tape | open_copy_exposure |

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

- Wallets: 2
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | A | 100.0 | 11.5 | 15m edge -10.00pp over 2 samples | 221.9% | 205 | existing | opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.7 | 1h edge -72.59pp over 1 samples | 50.9% | 92 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa5d9...88ad` | watch | - | B | 5.0 | $8158 | $13200 | 2 | 350.0% | 1 | open_copy_exposure |
| `0x81f0...cbb1` | pushed | - | A | 2.9 | $6045 | $6045 | 1 | 41.2% | 17 | open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 9
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x1550...5481` | pushed | - | C | 100.0 | 43.7 | 789.7 | 1451 | 493 | 100% | 174.4% | 18 | 174.4% | 18 | existing,sports_tape | burst_trading,open_copy_exposure |
| `0x3a01...c89e` | pushed | - | C | 100.0 | 43.1 | 660.9 | 1322 | 483 | 99% | 91.7% | 98 | 91.7% | 98 | existing,sports_tape | burst_trading |
| `0x249e...4fa6` | pushed | - | C | 100.0 | 36.4 | 513.6 | 1168 | 417 | 100% | 59.9% | 180 | 59.9% | 180 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0xa5d9...88ad` | prompt | - | B | 100.0 | 5.0 | 435.9 | 12 | 12 | 92% | 350.0% | 1 | 350.0% | 1 | existing,holder,sports_tape | open_copy_exposure |
| `0x50d0...a6fc` | pushed | - | A | 100.0 | 18.1 | 395.7 | 30 | 26 | 18% | 66.5% | 19 | 72.2% | 116 | existing,sports_tape | opposite_side_same_market |
| `0xf902...a267` | pushed | - | B | 100.0 | 26.8 | 354.3 | 21 | 18 | 84% | 89.2% | 11 | 69.7% | 13 | existing,sports_tape | opposite_side_same_market |
| `0x7e4d...b4ca` | prompt | - | B | 100.0 | 33.4 | 303.3 | 294 | 200 | 77% | 36.0% | 115 | 36.1% | 148 | existing,sports_tape | opposite_side_same_market |
| `0x81f0...cbb1` | pushed | - | A | 100.0 | 2.9 | 239.5 | 27 | 26 | 46% | 37.4% | 13 | 41.2% | 17 | existing,sports_tape | open_copy_exposure |
| `0x769c...2ee3` | prompt | - | B | 100.0 | 33.1 | 120.1 | 206 | 41 | 46% | 19.7% | 17 | 19.7% | 17 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 2310 | 174.3% | $+49498.45 | 95.7% | 138.0% | 101.1% | 0.36x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 22 | 2299 | 174.5% | $+49355.84 | 95.7% | 138.0% | 101.1% | 0.35x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 2191 | 175.8% | $+47721.92 | 96.7% | 128.6% | 101.1% | 0.30x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 19 | 2128 | 177.3% | $+46885.96 | 97.4% | 129.7% | 101.1% | 0.31x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 18 | 2117 | 177.5% | $+46743.35 | 97.4% | 132.8% | 101.1% | 0.30x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 19 | 2180 | 176.0% | $+47579.31 | 96.7% | 127.5% | 101.1% | 0.30x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 18 | 2106 | 177.7% | $+46499.43 | 97.4% | 128.6% | 101.1% | 0.30x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=20 smart>=70 |
| 8 | 18 | 2233 | 173.0% | $+47678.45 | 95.8% | 132.7% | 101.1% | 0.32x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 17 | 2095 | 177.9% | $+46356.82 | 97.4% | 127.5% | 101.1% | 0.30x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=20 smart>=70 |
| 10 | 41 | 3097 | 149.8% | $+58267.36 | 95.2% | 103.8% | 62.7% | 0.42x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 40 | 3084 | 150.1% | $+58176.80 | 95.2% | 104.1% | 62.7% | 0.42x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 40 | 3088 | 150.0% | $+58186.68 | 95.2% | 104.1% | 62.7% | 0.42x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
