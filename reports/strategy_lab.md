# Strategy Lab Report

**Generated:** 2026-10-26 23:46 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (434)
- Valid strategies found: 142
- Candidate layers: 21 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 3 tape
- Push wallets after live-edge blocks: 61 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 21
- Aggregate closed copy trades: 2329
- Aggregate copy ROI: 175.7%
- Aggregate copy PnL: $+50154.57
- Aggregate copy win rate: 95.6%
- Median wallet CopyROI: 141.9%
- Worst included wallet CopyROI: 100.0%
- Open copy cost / closed copy capital: 0.34x
- Params: tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.30 | 32 | 4% | 173.1% | $+5765.16 | 293 | 92.5% | 28.5% |
| `0x3724...69b5` | A | 17.02 | 1112 | 94% | 122.0% | $+4478.86 | 206 | 99.0% | 48.9% |
| `0x2b17...d36d` | A | 9.76 | 303 | 70% | 171.8% | $+3865.25 | 153 | 99.3% | 53.2% |
| `0x677d...4907` | A | 21.59 | 148 | 50% | 111.1% | $+2633.50 | 234 | 97.9% | 54.3% |
| `0x42c9...9623` | A | 4.85 | 0 | 0% | 117.3% | $+2333.81 | 150 | 100.0% | 49.3% |
| `0x075d...58c1` | A | 19.62 | 18 | 17% | 169.4% | $+1422.97 | 76 | 92.1% | 40.5% |
| `0x61b6...e14e` | A | 22.93 | 15 | 6% | 141.9% | $+950.53 | 66 | 72.7% | 18.5% |
| `0x2d8a...e818` | A | 24.57 | 2 | 1% | 127.8% | $+447.15 | 29 | 86.2% | 52.1% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0x768b...bf0f` | A | 2.83 | 14 | 18% | 103.3% | $+299.49 | 24 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.48 | 67 | 39% | 166.7% | $+266.69 | 16 | 93.8% | 33.5% |
| `0x1425...6e0c` | A | 23.91 | 177 | 64% | 106.3% | $+244.51 | 23 | 100.0% | 63.4% |
| `0x17e2...d472` | B | 29.41 | 326 | 31% | 256.6% | $+12853.82 | 470 | 97.9% | 53.7% |
| `0x5be9...ef2b` | B | 27.04 | 138 | 52% | 241.4% | $+3283.07 | 136 | 100.0% | 62.5% |
| `0x15dd...95ba` | B | 25.00 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x7967...f4a0` | B | 29.37 | 38 | 50% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.58 | 22 | 9% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x162d...8944` | B | 25.52 | 3 | 1% | 125.3% | $+475.98 | 37 | 91.9% | 47.3% |
| `0x092b...614e` | B | 26.86 | 119 | 98% | 105.0% | $+409.52 | 36 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.77 | 65 | 61% | 100.0% | $+380.09 | 36 | 100.0% | 47.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1044
- Aggregate copy ROI: 98.9%
- Aggregate copy PnL: $+12906.19
- Aggregate copy win rate: 88.6%
- Worst included CopyROI: 61.2%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 100.0 | 13.6 | 1259.9 | 46 | 47% | 379.1% | $+606.57 | 14 | 100.0% | $258 | existing,holder |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 836.0 | 15 | 12% | 132.2% | $+1454.13 | 99 | 93.9% | $9198 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.2 | 803.5 | 64 | 23% | 147.3% | $+883.76 | 59 | 98.3% | $254 | existing |
| `0x3520...1fa0` | B | 100.0 | 30.1 | 800.9 | 34 | 30% | 175.9% | $+510.03 | 28 | 96.4% | $832 | existing |
| `0xddf3...8d34` | A | 100.0 | 2.6 | 727.7 | 37 | 77% | 134.5% | $+484.35 | 28 | 100.0% | $7172 | existing |
| `0x612f...05d8` | B | 100.0 | 32.3 | 658.0 | 95 | 42% | 101.1% | $+849.32 | 80 | 93.8% | $589 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x50d0...a6fc` | A | 100.0 | 18.3 | 600.2 | 35 | 21% | 73.0% | $+971.14 | 118 | 100.0% | $2425 | existing |
| `0x7d68...8239` | A | 100.0 | 10.9 | 594.6 | 20 | 21% | 94.5% | $+396.86 | 42 | 97.6% | $284 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x7972...594f` | B | 100.0 | 33.1 | 575.6 | 322 | 70% | 67.7% | $+1808.77 | 204 | 62.2% | $1190 | existing,holder |
| `0x2bf0...728f` | A | 100.0 | 19.8 | 559.9 | 10 | 5% | 84.0% | $+571.09 | 66 | 84.8% | $1817 | existing |
| `0x988d...9899` | B | 100.0 | 27.6 | 540.7 | 140 | 73% | 87.0% | $+452.63 | 37 | 94.6% | $826 | existing,holder |
| `0x6d3d...51b8` | A | 100.0 | 22.4 | 538.4 | 202 | 100% | 72.8% | $+487.85 | 59 | 98.3% | $343 | existing,sports_tape |
| `0x2c5b...1dcd` | A | 100.0 | 3.5 | 530.1 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | $295 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x0d23...2d56` | B | 100.0 | 34.8 | 466.7 | 34 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5811 | existing |
| `0x6261...4adb` | A | 100.0 | 7.7 | 460.0 | 120 | 90% | 83.5% | $+292.35 | 12 | 91.7% | $1268 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.7 | 784.6 | 780 | 338 | 92% | 104.3% | $+3472.36 | 210 | 100.0% | 101.9% | existing,holder | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.1 | 747.8 | 1169 | 312 | 100% | 96.5% | $+2345.70 | 183 | 100.0% | 96.5% | existing,holder | burst_trading |
| `0x3231...6e75` | A | 100.0 | 0.8 | 656.3 | 502 | 308 | 100% | 78.3% | $+1252.30 | 160 | 100.0% | 78.3% | existing | - |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 566.9 | 1150 | 1112 | 100% | 60.3% | $+4366.05 | 296 | 95.3% | 60.3% | existing,holder | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.1 | 490.0 | 191 | 40 | 54% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x3d40...5739` | A | 100.0 | 24.9 | 463.6 | 129 | 71 | 76% | 70.0% | $+266.09 | 36 | 91.7% | 74.5% | existing | opposite_side_same_market |
| `0xd5c1...af9e` | A | 100.0 | 19.6 | 444.5 | 22 | 12 | 48% | 124.3% | $+87.03 | 6 | 83.3% | 66.8% | existing,holder | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.9 | 443.9 | 190 | 74 | 98% | 80.4% | $+152.69 | 19 | 94.7% | 80.4% | existing,holder | opposite_side_same_market |
| `0x4da5...3baf` | B | 100.0 | 30.8 | 432.8 | 158 | 131 | 99% | 49.9% | $+423.98 | 82 | 92.7% | 49.6% | existing | opposite_side_same_market |
| `0x789f...20ec` | B | 99.4 | 26.9 | 422.0 | 52 | 32 | 96% | 71.6% | $+171.80 | 18 | 94.4% | 71.2% | existing | opposite_side_same_market |

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
| `0x67ba...246c` | C | 100.0 | 39.3 | 249.0 | 965 | 571 | 98% | 58.2% | 166 | 58.8% | 171 | $665 | existing | opposite_side_same_market |
| `0x3fe2...93b3` | A | 100.0 | 20.5 | 234.1 | 97 | 80 | 100% | 27.5% | 31 | 27.5% | 31 | $1528 | existing,holder | opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 35.9 | 232.9 | 1278 | 421 | 92% | 90.1% | 36 | 90.2% | 42 | $567 | existing | burst_trading |
| `0x249e...4fa6` | C | 100.0 | 36.5 | 232.1 | 1170 | 398 | 100% | 60.8% | 174 | 60.8% | 174 | $147 | existing,sports_tape | burst_trading,opposite_side_same_market |
| `0x5f61...519a` | B | 100.0 | 27.0 | 231.8 | 474 | 292 | 100% | 40.6% | 113 | 40.6% | 113 | $244 | existing | opposite_side_same_market |
| `0x1550...5481` | C | 100.0 | 44.4 | 228.8 | 1448 | 491 | 100% | 152.4% | 37 | 152.4% | 37 | $295 | existing,holder | burst_trading,open_copy_exposure |
| `0x9e07...6e49` | C | 100.0 | 43.8 | 228.7 | 1353 | 477 | 97% | 92.6% | 47 | 92.3% | 54 | $411 | existing,holder | burst_trading |
| `0x3a01...c89e` | C | 100.0 | 43.7 | 227.7 | 1328 | 498 | 99% | 84.0% | 88 | 84.0% | 88 | $206 | existing,holder | burst_trading |
| `0x0e24...7014` | B | 100.0 | 30.4 | 221.8 | 217 | 207 | 91% | 46.8% | 90 | 49.6% | 96 | $1200 | existing | opposite_side_same_market |
| `0xb208...c095` | C | 100.0 | 41.0 | 218.5 | 1403 | 331 | 100% | 91.9% | 10 | 91.9% | 10 | $433 | existing,holder | burst_trading,open_copy_exposure |

## Recent Trade Flow Scout

- Wallets: 0
- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No recent-trade wallets passed the flow filters.

## Sports Tape Hotlist

- Wallets: 3
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xe9b7...ce7b` | C | 100.0 | 41.6 | 405.0 | 221 | 120 | 74.2% | 12 | 66.5% | 32 | $27890 | existing,sports_tape | burst_trading |
| `0x64ce...6015` | C | 100.0 | 39.5 | 327.0 | 536 | 135 | 110.1% | 4 | 109.8% | 5 | $1125 | existing,sports_tape | burst_trading |
| `0x6ac5...4b6e` | C | 100.0 | 43.6 | 303.9 | 1003 | 288 | 65.4% | 10 | 96.1% | 18 | $948 | existing,holder,sports_tape | burst_trading,open_copy_exposure |

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
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 60.3% | 296 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.4 | 1h edge -72.59pp over 1 samples | 49.6% | 96 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe9b7...ce7b` | pushed | - | C | 41.6 | $244803 | $244803 | 1 | 66.5% | 32 | burst_trading |
| `0x6ac5...4b6e` | pushed | - | C | 43.6 | $13822 | $13822 | 1 | 96.1% | 18 | burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 10
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x249e...4fa6` | pushed | - | C | 100.0 | 36.5 | 515.5 | 1170 | 398 | 100% | 60.8% | 174 | 60.8% | 174 | existing,sports_tape | burst_trading,opposite_side_same_market |
| `0x23cb...80d0` | reject | - | D | 100.0 | 57.1 | 514.7 | 1194 | 706 | 91% | 64.4% | 151 | 63.5% | 165 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x6d3d...51b8` | pushed | - | A | 100.0 | 22.4 | 495.6 | 202 | 107 | 100% | 72.8% | 59 | 72.8% | 59 | existing,sports_tape | opposite_side_same_market |
| `0x6e55...68d0` | reject | - | D | 100.0 | 54.3 | 481.5 | 1348 | 998 | 100% | 56.9% | 132 | 56.9% | 132 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x5d92...da54` | prompt | - | B | 100.0 | 30.9 | 362.9 | 136 | 95 | 100% | 61.0% | 33 | 61.0% | 33 | existing,sports_tape | opposite_side_same_market |
| `0x6ac5...4b6e` | pushed | - | C | 100.0 | 43.6 | 316.8 | 1003 | 288 | 69% | 65.4% | 10 | 96.1% | 18 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0xe9b7...ce7b` | pushed | - | C | 100.0 | 41.6 | 315.9 | 221 | 120 | 52% | 74.2% | 12 | 66.5% | 32 | existing,sports_tape | burst_trading |
| `0x64ce...6015` | pushed | - | C | 100.0 | 39.5 | 305.4 | 536 | 135 | 100% | 110.1% | 4 | 109.8% | 5 | existing,sports_tape | burst_trading |
| `0x63bc...92c5` | reject-bot | - | BOT | 100.0 | 62.5 | 274.6 | 1191 | 603 | 82% | 39.6% | 75 | 35.4% | 78 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xb26f...9536` | reject-bot | - | BOT | 99.7 | 60.5 | 95.5 | 1477 | 802 | 100% | 4.5% | 129 | 4.5% | 129 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 21 | 2329 | 175.7% | $+50154.57 | 95.6% | 141.9% | 100.0% | 0.34x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 22 | 2340 | 175.5% | $+50297.18 | 95.6% | 135.8% | 100.0% | 0.35x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 2235 | 176.7% | $+48770.83 | 96.6% | 128.7% | 100.0% | 0.29x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 19 | 2224 | 176.9% | $+48628.22 | 96.6% | 127.8% | 100.0% | 0.28x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 18 | 2150 | 178.6% | $+47548.34 | 97.3% | 128.7% | 100.0% | 0.29x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 17 | 2139 | 178.8% | $+47405.73 | 97.3% | 127.8% | 100.0% | 0.28x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 18 | 2282 | 173.8% | $+48675.32 | 95.7% | 126.5% | 100.0% | 0.32x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 41 | 3294 | 148.8% | $+60306.22 | 95.5% | 103.3% | 64.3% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 40 | 3282 | 149.0% | $+60227.25 | 95.6% | 104.1% | 64.3% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 36 | 3236 | 150.5% | $+59420.61 | 95.6% | 104.9% | 66.8% | 0.44x | tier>=B bot<30 copyT>=15 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 15 | 2108 | 176.9% | $+46373.06 | 97.3% | 125.3% | 100.0% | 0.27x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 16 | 2177 | 174.9% | $+47148.97 | 96.8% | 123.7% | 100.0% | 0.26x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
