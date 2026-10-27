# Strategy Lab Report

**Generated:** 2026-10-27 23:36 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (434)
- Valid strategies found: 118
- Candidate layers: 23 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 61 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 2474
- Aggregate copy ROI: 174.7%
- Aggregate copy PnL: $+52399.31
- Aggregate copy win rate: 95.8%
- Median wallet CopyROI: 147.7%
- Worst included wallet CopyROI: 100.0%
- Open copy cost / closed copy capital: 0.39x
- Params: tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.27 | 32 | 4% | 174.2% | $+5696.02 | 287 | 92.7% | 28.5% |
| `0x3724...69b5` | A | 17.16 | 1125 | 94% | 122.0% | $+4475.64 | 205 | 99.0% | 49.4% |
| `0x2b17...d36d` | A | 9.84 | 290 | 71% | 158.9% | $+3336.31 | 143 | 99.3% | 53.2% |
| `0x677d...4907` | A | 21.59 | 148 | 50% | 111.1% | $+2633.50 | 234 | 97.9% | 54.3% |
| `0x42c9...9623` | A | 5.11 | 0 | 0% | 106.9% | $+2202.74 | 157 | 100.0% | 49.3% |
| `0x075d...58c1` | A | 19.62 | 18 | 17% | 169.4% | $+1422.97 | 76 | 92.1% | 40.5% |
| `0x61b6...e14e` | A | 22.87 | 17 | 6% | 140.6% | $+956.30 | 67 | 73.1% | 18.7% |
| `0x2d8a...e818` | A | 24.57 | 2 | 1% | 127.8% | $+447.15 | 29 | 86.2% | 52.1% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0x768b...bf0f` | A | 2.79 | 14 | 18% | 103.3% | $+299.49 | 24 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.57 | 67 | 38% | 166.7% | $+266.69 | 16 | 93.8% | 33.4% |
| `0x1425...6e0c` | A | 23.91 | 177 | 64% | 106.3% | $+244.51 | 23 | 100.0% | 63.4% |
| `0x17e2...d472` | B | 29.31 | 327 | 31% | 264.0% | $+13436.77 | 474 | 97.9% | 52.9% |
| `0x5be9...ef2b` | B | 27.17 | 135 | 52% | 241.4% | $+3283.07 | 136 | 100.0% | 62.5% |
| `0x3ac2...f949` | B | 29.78 | 121 | 20% | 153.9% | $+1892.77 | 123 | 99.2% | 46.3% |
| `0x15dd...95ba` | B | 25.00 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x7967...f4a0` | B | 29.37 | 38 | 50% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.48 | 22 | 9% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x3520...1fa0` | B | 29.44 | 35 | 30% | 175.9% | $+510.03 | 28 | 96.4% | 58.9% |
| `0x162d...8944` | B | 25.85 | 3 | 1% | 125.3% | $+475.98 | 37 | 91.9% | 47.3% |
| `0x092b...614e` | B | 26.96 | 118 | 98% | 106.8% | $+395.11 | 35 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.77 | 65 | 61% | 100.0% | $+380.09 | 36 | 100.0% | 47.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 951
- Aggregate copy ROI: 97.9%
- Aggregate copy PnL: $+11512.22
- Aggregate copy win rate: 93.4%
- Worst included CopyROI: 53.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 100.0 | 13.6 | 1256.9 | 46 | 47% | 379.1% | $+606.57 | 14 | 100.0% | $258 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8a43...30da` | A | 99.6 | 24.1 | 848.1 | 27 | 84% | 222.2% | $+399.95 | 15 | 86.7% | $836 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 836.0 | 15 | 12% | 132.2% | $+1454.13 | 99 | 93.9% | $9198 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.3 | 811.1 | 64 | 23% | 149.8% | $+883.69 | 58 | 98.3% | $253 | existing |
| `0x612f...05d8` | B | 100.0 | 32.2 | 660.7 | 95 | 42% | 101.0% | $+869.06 | 82 | 93.9% | $587 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x50d0...a6fc` | A | 100.0 | 18.3 | 601.3 | 37 | 22% | 73.3% | $+974.51 | 118 | 100.0% | $2439 | existing |
| `0x7d68...8239` | A | 100.0 | 10.9 | 594.6 | 20 | 21% | 94.5% | $+396.86 | 42 | 97.6% | $284 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x2bf0...728f` | A | 100.0 | 19.8 | 559.9 | 10 | 5% | 84.0% | $+571.09 | 66 | 84.8% | $1817 | existing |
| `0x6d3d...51b8` | A | 100.0 | 22.3 | 533.2 | 205 | 100% | 72.4% | $+492.62 | 60 | 98.3% | $340 | existing,holder |
| `0x2c5b...1dcd` | A | 100.0 | 3.5 | 530.1 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | $295 | existing |
| `0x988d...9899` | B | 100.0 | 27.6 | 524.9 | 139 | 73% | 84.2% | $+429.62 | 36 | 94.4% | $828 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x0d23...2d56` | B | 100.0 | 34.6 | 467.1 | 33 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5942 | existing |
| `0x4da5...3baf` | B | 100.0 | 30.4 | 466.1 | 160 | 98% | 53.3% | $+469.25 | 85 | 92.9% | $1015 | existing |
| `0x6261...4adb` | A | 100.0 | 7.5 | 459.8 | 120 | 85% | 83.5% | $+292.35 | 12 | 91.7% | $1197 | existing |
| `0xb517...434b` | B | 100.0 | 31.6 | 459.7 | 167 | 67% | 55.8% | $+535.18 | 66 | 78.8% | $1244 | existing,sports_tape |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.7 | 781.6 | 787 | 341 | 92% | 103.3% | $+3493.08 | 214 | 99.5% | 101.2% | existing,holder | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.3 | 735.7 | 1173 | 299 | 100% | 95.1% | $+2111.27 | 177 | 100.0% | 95.1% | existing,holder | burst_trading |
| `0x1633...1e55` | B | 100.0 | 31.6 | 697.8 | 100 | 6 | 62% | 281.9% | $+140.94 | 5 | 100.0% | 262.8% | existing | opposite_side_same_market |
| `0x3231...6e75` | A | 100.0 | 0.8 | 656.4 | 505 | 310 | 100% | 78.2% | $+1258.80 | 161 | 100.0% | 78.2% | existing | - |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 565.0 | 1154 | 1116 | 100% | 60.1% | $+4337.13 | 295 | 94.6% | 60.1% | existing | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.1 | 490.0 | 193 | 40 | 54% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0xc89e...0770` | B | 100.0 | 33.6 | 488.0 | 100 | 73 | 53% | 85.1% | $+374.52 | 26 | 100.0% | 61.2% | existing | opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.3 | 462.8 | 130 | 70 | 77% | 70.6% | $+261.08 | 35 | 91.4% | 75.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.8 | 442.9 | 197 | 78 | 98% | 80.4% | $+152.69 | 19 | 94.7% | 80.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5f61...519a` | B | 100.0 | 26.8 | 422.5 | 510 | 325 | 100% | 43.6% | $+615.34 | 122 | 77.0% | 43.6% | existing,holder | opposite_side_same_market |

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
| `0x67ba...246c` | C | 100.0 | 39.2 | 248.4 | 942 | 564 | 98% | 57.9% | 165 | 58.5% | 170 | $677 | existing | opposite_side_same_market |
| `0x57a8...ebc2` | C | 100.0 | 44.2 | 245.9 | 1465 | 545 | 99% | 55.8% | 9 | 55.8% | 9 | $2146 | existing | burst_trading |
| `0x1550...5481` | C | 100.0 | 44.6 | 234.2 | 1444 | 532 | 100% | 153.2% | 42 | 153.2% | 42 | $328 | existing,holder | burst_trading |
| `0x9e07...6e49` | C | 100.0 | 42.8 | 233.7 | 1353 | 492 | 97% | 98.7% | 61 | 98.3% | 65 | $485 | existing,holder | burst_trading |
| `0x4d6e...5cfd` | C | 100.0 | 36.6 | 232.6 | 1270 | 412 | 91% | 105.6% | 33 | 101.2% | 38 | $520 | existing,holder | burst_trading |
| `0x3a01...c89e` | C | 100.0 | 42.9 | 231.7 | 1327 | 517 | 99% | 84.2% | 89 | 84.2% | 89 | $212 | existing,holder | burst_trading |
| `0x6ac5...4b6e` | C | 100.0 | 43.0 | 226.6 | 1329 | 418 | 91% | 62.0% | 11 | 67.4% | 20 | $1119 | existing,holder | burst_trading,open_copy_exposure |
| `0x1f3a...895d` | C | 100.0 | 42.3 | 224.2 | 1143 | 473 | 79% | 128.0% | 12 | 133.6% | 21 | $608 | existing,holder | burst_trading |
| `0x249e...4fa6` | C | 100.0 | 36.5 | 224.1 | 1169 | 397 | 100% | 60.1% | 175 | 60.1% | 175 | $149 | existing | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.5 | 221.9 | 218 | 208 | 91% | 47.0% | 91 | 49.7% | 97 | $1201 | existing | opposite_side_same_market |

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
| `0xce58...0a21` | B | 100.0 | 34.4 | 248.8 | 29 | 22 | 38.3% | 6 | 43.1% | 8 | $946 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 60.1% | 295 | existing | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.5 | 1h edge -72.59pp over 1 samples | 49.7% | 97 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 5
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6e55...68d0` | reject | - | D | 100.0 | 54.7 | 469.9 | 1349 | 973 | 100% | 55.7% | 131 | 55.7% | 131 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xb517...434b` | pushed | - | B | 100.0 | 31.6 | 410.1 | 167 | 123 | 67% | 68.3% | 45 | 55.8% | 66 | existing,sports_tape | opposite_side_same_market |
| `0x08c2...c023` | reject | - | D | 100.0 | 58.0 | 171.1 | 1404 | 634 | 98% | 17.4% | 92 | 19.9% | 101 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xce58...0a21` | pushed | - | B | 100.0 | 34.4 | 159.2 | 29 | 22 | 88% | 38.3% | 6 | 43.1% | 8 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x9155...32ef` | prompt | - | B | 100.0 | 31.6 | 111.4 | 106 | 55 | 94% | 12.8% | 21 | 12.8% | 21 | existing,holder,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 2474 | 174.7% | $+52399.31 | 95.8% | 147.7% | 100.0% | 0.39x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 21 | 2368 | 175.8% | $+50867.19 | 96.8% | 153.9% | 100.0% | 0.34x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 24 | 2485 | 174.5% | $+52541.92 | 95.8% | 144.1% | 100.0% | 0.39x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 19 | 2283 | 177.6% | $+49644.70 | 97.5% | 153.9% | 100.0% | 0.34x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 22 | 2379 | 175.7% | $+51009.80 | 96.8% | 141.8% | 100.0% | 0.34x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 20 | 2294 | 177.4% | $+49787.31 | 97.5% | 141.8% | 100.0% | 0.34x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 20 | 2427 | 172.8% | $+50920.06 | 95.9% | 134.2% | 100.0% | 0.37x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 42 | 3419 | 149.4% | $+62226.18 | 95.8% | 105.6% | 63.7% | 0.49x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 41 | 3407 | 149.6% | $+62147.21 | 95.9% | 106.3% | 63.7% | 0.49x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 18 | 2321 | 174.0% | $+49387.94 | 96.9% | 126.5% | 100.0% | 0.32x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 11 | 17 | 2252 | 175.8% | $+48612.03 | 97.5% | 127.8% | 100.0% | 0.32x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 37 | 3360 | 151.2% | $+61297.33 | 95.9% | 106.8% | 72.4% | 0.47x | tier>=B bot<30 copyT>=15 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
