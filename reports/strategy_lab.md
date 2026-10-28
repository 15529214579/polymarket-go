# Strategy Lab Report

**Generated:** 2026-10-28 23:27 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (434)
- Valid strategies found: 138
- Candidate layers: 26 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 63 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 26
- Aggregate closed copy trades: 2067
- Aggregate copy ROI: 154.8%
- Aggregate copy PnL: $+39715.12
- Aggregate copy win rate: 95.1%
- Median wallet CopyROI: 128.7%
- Worst included wallet CopyROI: 100.0%
- Open copy cost / closed copy capital: 0.49x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.35 | 32 | 4% | 176.5% | $+5666.11 | 281 | 92.5% | 28.5% |
| `0x3724...69b5` | A | 17.19 | 1123 | 95% | 121.3% | $+4525.43 | 209 | 98.6% | 49.4% |
| `0x2b17...d36d` | A | 9.87 | 280 | 70% | 154.7% | $+3139.71 | 139 | 99.3% | 53.2% |
| `0x677d...4907` | A | 21.59 | 148 | 50% | 111.1% | $+2633.50 | 234 | 97.9% | 54.3% |
| `0x42c9...9623` | A | 20.01 | 0 | 0% | 110.1% | $+2355.84 | 164 | 100.0% | 49.3% |
| `0x075d...58c1` | A | 19.71 | 18 | 17% | 170.7% | $+1399.43 | 74 | 91.9% | 40.5% |
| `0x61b6...e14e` | A | 23.23 | 18 | 6% | 122.9% | $+934.26 | 75 | 68.0% | 16.9% |
| `0x2d8a...e818` | A | 24.57 | 2 | 1% | 127.8% | $+447.15 | 29 | 86.2% | 52.1% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0x768b...bf0f` | A | 3.03 | 14 | 20% | 105.0% | $+272.91 | 22 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.74 | 68 | 38% | 166.7% | $+266.69 | 16 | 93.8% | 33.5% |
| `0x1425...6e0c` | A | 23.91 | 177 | 64% | 106.3% | $+244.51 | 23 | 100.0% | 63.4% |
| `0x3280...ea74` | A | 14.72 | 40 | 74% | 112.7% | $+214.20 | 16 | 100.0% | 46.1% |
| `0x1b4d...e246` | A | 21.94 | 29 | 55% | 105.2% | $+189.31 | 13 | 92.3% | 52.5% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | 47.7% |
| `0x5be9...ef2b` | B | 27.08 | 136 | 52% | 243.0% | $+3377.73 | 139 | 100.0% | 62.5% |
| `0x3ac2...f949` | B | 29.76 | 121 | 20% | 153.9% | $+1892.77 | 123 | 99.2% | 46.3% |
| `0x15dd...95ba` | B | 25.00 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x7967...f4a0` | B | 29.37 | 38 | 50% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.48 | 22 | 9% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x3520...1fa0` | B | 29.08 | 35 | 29% | 175.9% | $+510.03 | 28 | 96.4% | 58.9% |
| `0x162d...8944` | B | 25.91 | 3 | 2% | 133.0% | $+452.33 | 33 | 93.9% | 47.3% |
| `0x4607...d542` | B | 26.77 | 65 | 61% | 100.0% | $+380.09 | 36 | 100.0% | 47.6% |
| `0x092b...614e` | B | 27.37 | 114 | 98% | 103.6% | $+342.00 | 32 | 100.0% | 55.2% |
| `0xbe12...728f` | B | 25.81 | 25 | 25% | 105.3% | $+284.34 | 26 | 100.0% | 51.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1032
- Aggregate copy ROI: 95.9%
- Aggregate copy PnL: $+12072.25
- Aggregate copy win rate: 92.2%
- Worst included CopyROI: 52.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 100.0 | 13.6 | 1256.9 | 46 | 47% | 379.1% | $+606.57 | 14 | 100.0% | $258 | existing |
| `0x8a43...30da` | B | 100.0 | 25.8 | 1013.3 | 27 | 69% | 248.7% | $+621.69 | 22 | 90.9% | $919 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 836.0 | 15 | 12% | 132.2% | $+1454.13 | 99 | 93.9% | $9198 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.3 | 811.1 | 64 | 23% | 149.8% | $+883.69 | 58 | 98.3% | $253 | existing |
| `0x612f...05d8` | B | 100.0 | 32.1 | 660.5 | 96 | 42% | 101.0% | $+869.06 | 82 | 93.9% | $587 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x50d0...a6fc` | A | 100.0 | 18.2 | 599.3 | 45 | 26% | 72.6% | $+994.85 | 120 | 100.0% | $2467 | existing |
| `0x7d68...8239` | A | 100.0 | 10.9 | 594.6 | 20 | 21% | 94.5% | $+396.86 | 42 | 97.6% | $284 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x2bf0...728f` | A | 100.0 | 19.9 | 552.1 | 10 | 4% | 81.5% | $+578.32 | 69 | 82.6% | $1822 | existing |
| `0x988d...9899` | B | 100.0 | 28.0 | 533.8 | 134 | 72% | 87.7% | $+421.05 | 35 | 94.3% | $823 | existing |
| `0x6d3d...51b8` | A | 100.0 | 22.4 | 528.7 | 202 | 100% | 72.1% | $+482.93 | 60 | 98.3% | $342 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x5f61...519a` | B | 100.0 | 27.0 | 477.1 | 527 | 100% | 52.9% | $+794.00 | 127 | 78.0% | $280 | existing,holder |
| `0x0d23...2d56` | B | 100.0 | 34.6 | 467.1 | 33 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5942 | existing |
| `0x4da5...3baf` | B | 100.0 | 30.4 | 466.1 | 160 | 98% | 53.3% | $+469.25 | 85 | 92.9% | $1015 | existing |
| `0x6261...4adb` | A | 100.0 | 7.4 | 462.4 | 122 | 84% | 83.5% | $+292.35 | 12 | 91.7% | $1178 | existing,holder |
| `0xca01...20e4` | B | 100.0 | 30.8 | 458.3 | 31 | 54% | 70.9% | $+212.77 | 20 | 90.0% | $4889 | existing,sports_tape |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.1 | 782.0 | 745 | 337 | 92% | 103.7% | $+3442.29 | 210 | 99.5% | 101.7% | existing,holder | burst_trading |
| `0x1633...1e55` | B | 100.0 | 31.2 | 750.1 | 112 | 7 | 63% | 284.2% | $+170.54 | 6 | 100.0% | 267.6% | existing | opposite_side_same_market |
| `0x5135...7217` | A | 100.0 | 18.2 | 725.9 | 1159 | 295 | 100% | 92.7% | $+2151.63 | 180 | 100.0% | 92.7% | existing,holder | burst_trading |
| `0x3231...6e75` | A | 100.0 | 0.8 | 661.2 | 507 | 312 | 100% | 79.0% | $+1279.74 | 162 | 100.0% | 79.0% | existing | - |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 567.4 | 1156 | 1118 | 100% | 60.5% | $+4374.69 | 294 | 94.6% | 60.5% | existing | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 33.9 | 490.3 | 194 | 40 | 53% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0xc89e...0770` | B | 100.0 | 33.8 | 487.8 | 100 | 73 | 54% | 85.1% | $+374.52 | 26 | 100.0% | 61.2% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.9 | 470.3 | 191 | 75 | 99% | 90.2% | $+162.43 | 18 | 100.0% | 90.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb517...434b` | B | 100.0 | 32.2 | 469.9 | 173 | 129 | 68% | 69.3% | $+464.20 | 46 | 82.6% | 57.0% | existing | opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.6 | 464.6 | 121 | 65 | 76% | 73.2% | $+248.80 | 32 | 90.6% | 77.3% | existing | opposite_side_same_market |

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
| `0xc9bf...81c0` | C | 100.0 | 37.9 | 257.5 | 1303 | 665 | 95% | 104.8% | 77 | 105.2% | 79 | $258 | existing,holder | burst_trading |
| `0x67ba...246c` | C | 100.0 | 39.3 | 247.9 | 921 | 559 | 98% | 57.7% | 163 | 58.3% | 168 | $688 | existing | opposite_side_same_market |
| `0x57a8...ebc2` | C | 100.0 | 44.5 | 239.6 | 1466 | 508 | 99% | 64.2% | 9 | 64.2% | 9 | $1782 | existing | burst_trading |
| `0x1550...5481` | C | 100.0 | 44.4 | 238.1 | 1435 | 558 | 100% | 166.4% | 50 | 166.4% | 50 | $354 | existing,holder | burst_trading |
| `0x3a01...c89e` | C | 100.0 | 41.9 | 235.6 | 1321 | 556 | 99% | 82.7% | 93 | 82.7% | 93 | $226 | existing | burst_trading |
| `0x9e07...6e49` | C | 100.0 | 42.5 | 233.5 | 1363 | 480 | 98% | 94.2% | 71 | 94.0% | 74 | $503 | existing,holder | burst_trading |
| `0x4d6e...5cfd` | C | 100.0 | 37.4 | 230.0 | 1266 | 406 | 90% | 102.3% | 31 | 98.6% | 34 | $521 | existing,holder | burst_trading |
| `0x6ac5...4b6e` | C | 100.0 | 40.0 | 224.3 | 1426 | 331 | 96% | 52.4% | 13 | 53.9% | 14 | $1019 | existing,holder | burst_trading,open_copy_exposure |
| `0x0e24...7014` | B | 100.0 | 30.5 | 221.9 | 220 | 210 | 91% | 47.8% | 92 | 50.5% | 98 | $1228 | existing | opposite_side_same_market |
| `0x73c2...3a29` | B | 100.0 | 33.6 | 215.9 | 546 | 337 | 88% | 11.5% | 165 | 14.0% | 183 | $165 | existing | opposite_side_same_market |

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
| `0x0f00...73ce` | A | 100.0 | 11.5 | 15m edge -10.00pp over 2 samples | 221.9% | 205 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 60.5% | 294 | existing | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.5 | 1h edge -72.59pp over 1 samples | 50.5% | 98 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xca01...20e4` | pushed | - | B | 100.0 | 30.8 | 340.0 | 31 | 30 | 54% | 90.4% | 9 | 70.9% | 20 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 26 | 2067 | 154.8% | $+39715.12 | 95.1% | 128.7% | 100.0% | 0.49x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 25 | 1992 | 155.8% | $+38780.86 | 96.1% | 129.7% | 100.0% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 23 | 1968 | 156.3% | $+38448.94 | 96.1% | 133.0% | 100.0% | 0.45x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 24 | 2043 | 155.3% | $+39383.20 | 95.1% | 130.4% | 100.0% | 0.48x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 24 | 1953 | 155.9% | $+38205.04 | 96.6% | 128.7% | 100.0% | 0.44x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 22 | 1929 | 156.4% | $+37873.12 | 96.6% | 130.4% | 100.0% | 0.42x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 22 | 1868 | 157.2% | $+36982.55 | 97.4% | 128.7% | 100.0% | 0.44x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 8 | 20 | 1844 | 157.7% | $+36650.63 | 97.4% | 130.4% | 100.0% | 0.43x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 46 | 3113 | 131.2% | $+50687.51 | 95.3% | 104.9% | 63.7% | 0.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 45 | 3102 | 131.4% | $+50617.22 | 95.4% | 105.0% | 63.7% | 0.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 45 | 3038 | 131.3% | $+49753.25 | 96.0% | 104.8% | 63.7% | 0.55x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 44 | 3027 | 131.5% | $+49682.96 | 96.1% | 104.9% | 63.7% | 0.55x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
