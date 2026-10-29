# Strategy Lab Report

**Generated:** 2026-10-29 22:25 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (434)
- Valid strategies found: 209
- Candidate layers: 28 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 66 total
- Live-edge blocked push wallets: 2
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 28
- Aggregate closed copy trades: 2199
- Aggregate copy ROI: 153.2%
- Aggregate copy PnL: $+41612.61
- Aggregate copy win rate: 95.1%
- Median wallet CopyROI: 131.6%
- Worst included wallet CopyROI: 101.7%
- Open copy cost / closed copy capital: 0.50x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.19 | 32 | 4% | 178.6% | $+5715.68 | 278 | 92.5% | 28.5% |
| `0x3724...69b5` | A | 17.18 | 1118 | 94% | 120.4% | $+4491.07 | 211 | 98.6% | 49.4% |
| `0x2b17...d36d` | A | 9.80 | 269 | 70% | 157.8% | $+2967.43 | 129 | 99.2% | 53.2% |
| `0x677d...4907` | A | 21.59 | 148 | 50% | 111.1% | $+2633.50 | 234 | 97.9% | 54.3% |
| `0x42c9...9623` | A | 20.17 | 0 | 0% | 114.2% | $+2511.55 | 170 | 100.0% | 49.3% |
| `0x075d...58c1` | A | 19.76 | 18 | 17% | 171.8% | $+1391.78 | 73 | 91.8% | 40.5% |
| `0x61b6...e14e` | A | 23.52 | 19 | 5% | 133.5% | $+1214.72 | 89 | 68.5% | 21.5% |
| `0x2d8a...e818` | A | 24.57 | 2 | 1% | 127.8% | $+447.15 | 29 | 86.2% | 52.1% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0x7bdc...9b4c` | A | 24.86 | 73 | 40% | 166.7% | $+266.69 | 16 | 93.8% | 33.6% |
| `0x1425...6e0c` | A | 24.19 | 181 | 64% | 104.8% | $+251.58 | 24 | 100.0% | 62.8% |
| `0x768b...bf0f` | A | 2.96 | 14 | 21% | 106.0% | $+243.70 | 19 | 100.0% | 47.3% |
| `0x3280...ea74` | A | 13.95 | 42 | 74% | 113.0% | $+237.32 | 18 | 100.0% | 47.2% |
| `0x1b4d...e246` | A | 11.89 | 22 | 41% | 101.7% | $+233.95 | 16 | 93.8% | 52.5% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | 47.7% |
| `0x5be9...ef2b` | B | 27.08 | 136 | 52% | 243.0% | $+3377.73 | 139 | 100.0% | 62.5% |
| `0x3ac2...f949` | B | 29.76 | 121 | 20% | 153.9% | $+1892.77 | 123 | 99.2% | 46.3% |
| `0x36e1...28eb` | B | 26.59 | 88 | 14% | 103.3% | $+1312.29 | 121 | 100.0% | 47.9% |
| `0x15dd...95ba` | B | 25.00 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x7967...f4a0` | B | 29.37 | 38 | 50% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.48 | 22 | 9% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x3520...1fa0` | B | 28.85 | 35 | 29% | 164.2% | $+525.56 | 30 | 96.7% | 58.9% |
| `0x162d...8944` | B | 25.66 | 3 | 2% | 135.8% | $+434.70 | 31 | 93.5% | 47.3% |
| `0xca01...20e4` | B | 29.85 | 39 | 60% | 116.1% | $+429.48 | 23 | 91.3% | 12.4% |
| `0xbe12...728f` | B | 26.56 | 26 | 24% | 110.5% | $+342.47 | 29 | 100.0% | 52.9% |
| `0x092b...614e` | B | 27.37 | 114 | 98% | 104.2% | $+333.32 | 31 | 100.0% | 55.2% |
| `0xb27f...4a88` | B | 25.62 | 108 | 35% | 142.8% | $+171.39 | 11 | 100.0% | 44.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1067
- Aggregate copy ROI: 99.0%
- Aggregate copy PnL: $+12549.07
- Aggregate copy win rate: 92.6%
- Worst included CopyROI: 53.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 100.0 | 13.6 | 1256.9 | 46 | 47% | 379.1% | $+606.57 | 14 | 100.0% | $258 | existing |
| `0x8a43...30da` | B | 93.7 | 33.3 | 1202.8 | 34 | 59% | 280.3% | $+1037.00 | 33 | 93.9% | $898 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 836.0 | 15 | 12% | 132.2% | $+1454.13 | 99 | 93.9% | $9198 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.4 | 811.3 | 64 | 23% | 149.8% | $+883.69 | 58 | 98.3% | $253 | existing |
| `0x612f...05d8` | B | 100.0 | 32.3 | 667.7 | 105 | 43% | 101.7% | $+915.09 | 86 | 94.2% | $595 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x50d0...a6fc` | A | 100.0 | 18.3 | 600.7 | 49 | 28% | 72.4% | $+1020.34 | 124 | 100.0% | $2543 | existing |
| `0x7d68...8239` | A | 100.0 | 10.9 | 594.6 | 20 | 21% | 94.5% | $+396.86 | 42 | 97.6% | $284 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x4607...d542` | B | 100.0 | 26.4 | 589.4 | 68 | 62% | 98.2% | $+392.82 | 37 | 100.0% | $4994 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x2bf0...728f` | A | 100.0 | 20.2 | 540.4 | 11 | 5% | 78.1% | $+570.15 | 70 | 81.4% | $1799 | existing |
| `0x988d...9899` | B | 100.0 | 28.1 | 533.9 | 133 | 72% | 87.7% | $+421.05 | 35 | 94.3% | $825 | existing |
| `0x6d3d...51b8` | A | 100.0 | 22.5 | 526.4 | 206 | 100% | 71.3% | $+484.87 | 61 | 98.4% | $338 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0xa5d9...88ad` | B | 100.0 | 4.1 | 491.2 | 18 | 64% | 174.0% | $+69.61 | 3 | 100.0% | $5679 | existing |
| `0x5f61...519a` | B | 100.0 | 27.1 | 480.4 | 548 | 100% | 53.2% | $+824.99 | 132 | 78.8% | $293 | existing,holder |
| `0x4da5...3baf` | B | 100.0 | 30.4 | 468.2 | 160 | 98% | 53.6% | $+477.17 | 86 | 93.0% | $1015 | existing |
| `0x0d23...2d56` | B | 100.0 | 34.6 | 467.1 | 33 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5942 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.2 | 781.9 | 731 | 335 | 92% | 103.6% | $+3439.38 | 211 | 99.5% | 101.6% | existing | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.2 | 731.1 | 1171 | 292 | 100% | 94.1% | $+2117.33 | 177 | 100.0% | 94.1% | existing,holder | burst_trading |
| `0x3231...6e75` | A | 100.0 | 0.8 | 660.8 | 512 | 317 | 100% | 78.6% | $+1297.30 | 165 | 100.0% | 78.6% | existing | - |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 562.0 | 1167 | 1124 | 100% | 59.9% | $+4319.47 | 288 | 93.1% | 59.9% | existing,holder | burst_trading,opposite_side_same_market |
| `0xbdd2...6657` | B | 100.0 | 32.0 | 535.4 | 1142 | 330 | 83% | 82.3% | $+1193.90 | 50 | 100.0% | 82.6% | existing,holder | burst_trading,open_copy_exposure |
| `0x3830...d2d5` | B | 100.0 | 33.7 | 490.6 | 195 | 40 | 53% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0xc89e...0770` | B | 100.0 | 33.7 | 487.9 | 102 | 73 | 55% | 85.1% | $+374.52 | 26 | 100.0% | 61.2% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 27.2 | 474.6 | 198 | 80 | 99% | 90.2% | $+171.45 | 19 | 100.0% | 90.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.6 | 449.3 | 119 | 63 | 75% | 69.1% | $+235.08 | 32 | 87.5% | 74.0% | existing | opposite_side_same_market |
| `0x6261...4adb` | A | 100.0 | 7.4 | 426.9 | 122 | 77 | 84% | 81.9% | $+270.13 | 10 | 90.0% | 83.5% | existing | open_copy_exposure |

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
| `0xc9bf...81c0` | C | 100.0 | 40.3 | 251.5 | 1318 | 637 | 95% | 115.0% | 71 | 116.2% | 74 | $322 | existing,holder | burst_trading |
| `0x67ba...246c` | C | 100.0 | 39.4 | 246.5 | 886 | 548 | 98% | 58.3% | 160 | 58.9% | 165 | $709 | existing | opposite_side_same_market |
| `0x1550...5481` | C | 100.0 | 44.1 | 241.4 | 1429 | 581 | 100% | 179.4% | 54 | 179.4% | 54 | $357 | existing,holder | burst_trading |
| `0x3a01...c89e` | C | 100.0 | 42.3 | 236.7 | 1320 | 548 | 99% | 73.8% | 92 | 73.8% | 92 | $221 | existing,holder | burst_trading |
| `0x9e07...6e49` | C | 100.0 | 42.4 | 235.7 | 1408 | 486 | 100% | 92.5% | 73 | 92.5% | 73 | $510 | existing,holder | burst_trading |
| `0x4d6e...5cfd` | C | 100.0 | 36.9 | 233.9 | 1305 | 424 | 91% | 101.7% | 30 | 98.0% | 33 | $571 | existing,holder | burst_trading |
| `0x1f3a...895d` | C | 100.0 | 42.5 | 225.8 | 1110 | 497 | 78% | 143.3% | 16 | 144.7% | 36 | $882 | existing | burst_trading |
| `0x6ac5...4b6e` | C | 100.0 | 43.3 | 223.5 | 1322 | 407 | 90% | 79.0% | 9 | 78.6% | 11 | $1004 | existing,holder | burst_trading,open_copy_exposure |
| `0x0286...f61d` | C | 100.0 | 37.0 | 222.5 | 1009 | 399 | 81% | 92.5% | 50 | 90.7% | 66 | $384 | existing,holder | burst_trading,open_copy_exposure |
| `0x7e4d...b4ca` | B | 100.0 | 25.0 | 221.0 | 290 | 201 | 83% | 38.8% | 111 | 41.0% | 134 | $363 | existing,holder | opposite_side_same_market |

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

- Wallets: 2
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | A | 100.0 | 11.5 | 15m edge -10.00pp over 2 samples | 221.9% | 205 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 59.9% | 288 | existing,holder | burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 28 | 2199 | 153.2% | $+41612.61 | 95.1% | 131.6% | 101.7% | 0.50x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 27 | 2110 | 153.9% | $+40397.89 | 96.3% | 129.7% | 101.7% | 0.47x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 26 | 2137 | 153.8% | $+40607.31 | 95.6% | 131.6% | 101.7% | 0.49x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 26 | 2177 | 153.4% | $+41298.61 | 95.1% | 130.6% | 101.7% | 0.50x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 25 | 2048 | 154.5% | $+39392.59 | 96.7% | 129.7% | 101.7% | 0.45x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 6 | 24 | 2115 | 154.0% | $+40293.31 | 95.5% | 130.6% | 101.7% | 0.48x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 25 | 2088 | 154.1% | $+40083.89 | 96.2% | 127.8% | 101.7% | 0.46x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 24 | 1986 | 155.1% | $+38599.58 | 97.4% | 128.7% | 101.7% | 0.45x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 23 | 1963 | 155.7% | $+38170.10 | 97.5% | 129.7% | 101.7% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=20 smart>=70 |
| 10 | 23 | 2026 | 154.7% | $+39078.59 | 96.7% | 127.8% | 101.7% | 0.45x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 11 | 50 | 3304 | 130.1% | $+53153.76 | 95.3% | 104.5% | 63.9% | 0.58x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 49 | 3293 | 130.3% | $+53083.47 | 95.4% | 104.8% | 64.0% | 0.58x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
