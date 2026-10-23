# Strategy Lab Report

**Generated:** 2026-10-23 23:04 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (432)
- Valid strategies found: 162
- Candidate layers: 23 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 60 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 2013
- Aggregate copy ROI: 177.8%
- Aggregate copy PnL: $+44509.20
- Aggregate copy win rate: 95.5%
- Median wallet CopyROI: 134.8%
- Worst included wallet CopyROI: 100.5%
- Open copy cost / closed copy capital: 0.54x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x45e4...3c69` | A | 22.43 | 34 | 4% | 167.8% | $+5285.03 | 275 | 91.6% | 28.5% |
| `0x3724...69b5` | A | 16.93 | 1114 | 94% | 122.0% | $+4513.24 | 206 | 98.5% | 48.9% |
| `0x42c9...9623` | A | 12.69 | 0 | 0% | 134.8% | $+2845.25 | 161 | 100.0% | 49.3% |
| `0x2b17...d36d` | A | 9.76 | 303 | 70% | 124.8% | $+1797.38 | 100 | 99.0% | 52.9% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x15dd...95ba` | A | 24.93 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x61b6...e14e` | A | 23.16 | 10 | 5% | 131.4% | $+551.81 | 41 | 68.3% | 17.4% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 194.2% | $+446.58 | 16 | 75.0% | 41.8% |
| `0x162d...8944` | A | 24.69 | 3 | 1% | 105.5% | $+401.08 | 37 | 91.9% | 47.7% |
| `0xddf3...8d34` | A | 2.84 | 33 | 75% | 138.1% | $+386.53 | 22 | 100.0% | 14.6% |
| `0x768b...bf0f` | A | 2.83 | 14 | 18% | 103.3% | $+299.49 | 24 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.68 | 64 | 38% | 166.7% | $+266.69 | 16 | 93.8% | 32.0% |
| `0x1425...6e0c` | A | 23.89 | 174 | 64% | 108.8% | $+239.28 | 22 | 100.0% | 63.0% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.74 | 336 | 32% | 249.8% | $+12241.54 | 459 | 98.3% | 56.9% |
| `0x0f00...73ce` | B | 26.63 | 177 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.17 | 135 | 52% | 246.4% | $+3154.16 | 128 | 100.0% | 61.9% |
| `0x7967...f4a0` | B | 29.57 | 36 | 49% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.31 | 22 | 10% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x092b...614e` | B | 27.37 | 115 | 99% | 100.5% | $+391.98 | 36 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.86 | 63 | 62% | 101.1% | $+363.82 | 34 | 100.0% | 47.6% |
| `0xfeca...8c27` | B | 27.22 | 178 | 53% | 173.8% | $+312.83 | 18 | 100.0% | 69.8% |
| `0xbe12...728f` | B | 27.29 | 22 | 27% | 104.3% | $+198.14 | 19 | 100.0% | 52.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1171
- Aggregate copy ROI: 104.8%
- Aggregate copy PnL: $+15085.00
- Aggregate copy win rate: 88.8%
- Worst included CopyROI: 60.5%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 98.3 | 13.9 | 1285.2 | 47 | 44% | 401.8% | $+602.69 | 13 | 100.0% | $303 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8a43...30da` | B | 97.3 | 25.2 | 871.7 | 21 | 88% | 241.7% | $+386.75 | 13 | 92.3% | $1050 | existing |
| `0xbe7e...6d00` | B | 100.0 | 30.6 | 863.1 | 99 | 12% | 117.9% | $+2983.73 | 230 | 98.7% | $1782 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.2 | 804.2 | 63 | 23% | 146.7% | $+894.95 | 60 | 98.3% | $258 | existing |
| `0x3520...1fa0` | B | 100.0 | 30.0 | 801.4 | 34 | 31% | 175.9% | $+510.03 | 28 | 96.4% | $835 | existing |
| `0x760f...326a` | B | 100.0 | 30.8 | 798.6 | 15 | 11% | 127.0% | $+1283.09 | 90 | 93.3% | $9094 | existing |
| `0x9a0c...370e` | B | 100.0 | 30.7 | 642.5 | 185 | 78% | 83.0% | $+1535.84 | 156 | 91.7% | $333 | existing |
| `0x612f...05d8` | B | 100.0 | 32.9 | 641.1 | 81 | 42% | 101.9% | $+713.16 | 66 | 93.9% | $589 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x7d68...8239` | A | 100.0 | 10.7 | 595.6 | 20 | 22% | 94.5% | $+396.86 | 42 | 97.6% | $288 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x2bf0...728f` | A | 100.0 | 19.9 | 571.1 | 8 | 4% | 88.7% | $+558.59 | 61 | 85.2% | $1840 | existing,holder |
| `0x7972...594f` | B | 100.0 | 33.1 | 566.0 | 300 | 68% | 67.2% | $+1692.39 | 194 | 61.9% | $1153 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x5925...9984` | B | 100.0 | 31.7 | 487.4 | 48 | 96% | 97.0% | $+174.66 | 15 | 80.0% | $1547 | existing |
| `0x0d23...2d56` | B | 100.0 | 26.2 | 472.5 | 35 | 19% | 60.5% | $+344.55 | 56 | 98.2% | $5052 | existing |
| `0x6261...4adb` | A | 100.0 | 7.7 | 460.0 | 120 | 90% | 83.5% | $+292.35 | 12 | 91.7% | $1268 | existing |
| `0x578e...c3c0` | B | 100.0 | 25.5 | 456.9 | 6 | 19% | 151.5% | $+75.75 | 5 | 100.0% | $728 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa6de...256c` | B | 100.0 | 30.8 | 888.6 | 96 | 43 | 8% | 250.7% | $+551.54 | 14 | 100.0% | 216.1% | existing | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.6 | 771.6 | 792 | 343 | 92% | 103.3% | $+3255.21 | 196 | 100.0% | 101.1% | existing | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.1 | 744.3 | 1174 | 323 | 100% | 95.4% | $+2403.95 | 187 | 100.0% | 95.4% | existing,holder | burst_trading |
| `0x677d...4907` | A | 100.0 | 21.4 | 696.5 | 152 | 125 | 50% | 103.7% | $+923.13 | 88 | 97.7% | 98.5% | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 570.4 | 1146 | 1108 | 100% | 60.9% | $+4336.39 | 297 | 96.0% | 60.9% | existing,holder | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.3 | 489.7 | 184 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.6 | 477.0 | 152 | 75 | 80% | 71.6% | $+300.63 | 40 | 90.0% | 76.8% | existing,holder | opposite_side_same_market |
| `0x2b46...684d` | B | 100.0 | 34.7 | 442.9 | 84 | 76 | 8% | 123.8% | $+457.89 | 6 | 100.0% | 163.1% | existing | burst_trading |
| `0xf28e...cab2` | B | 100.0 | 33.4 | 441.8 | 93 | 46 | 10% | 138.6% | $+235.69 | 5 | 100.0% | 86.3% | existing | burst_trading |
| `0x21cc...54bc` | B | 100.0 | 26.6 | 424.5 | 182 | 69 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | open_copy_exposure,opposite_side_same_market |

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
| `0x67ba...246c` | C | 100.0 | 39.1 | 250.5 | 1034 | 582 | 98% | 57.4% | 163 | 58.0% | 168 | $627 | existing | opposite_side_same_market |
| `0x89eb...cff9` | C | 100.0 | 37.0 | 237.5 | 441 | 410 | 90% | 57.7% | 154 | 60.6% | 169 | $2084 | existing | opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 36.0 | 234.5 | 1196 | 420 | 87% | 91.1% | 51 | 89.7% | 64 | $721 | existing,holder | burst_trading |
| `0x1550...5481` | C | 100.0 | 43.2 | 228.5 | 1456 | 488 | 100% | 175.4% | 10 | 175.4% | 10 | $222 | existing,holder | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 38.5 | 221.8 | 1358 | 346 | 94% | 118.2% | 32 | 118.4% | 33 | $461 | existing,holder | burst_trading |
| `0x0e24...7014` | B | 100.0 | 30.8 | 219.8 | 206 | 197 | 91% | 47.6% | 85 | 50.4% | 91 | $1204 | existing | opposite_side_same_market |
| `0x234f...6320` | C | 100.0 | 43.1 | 216.3 | 461 | 356 | 100% | 28.2% | 156 | 28.2% | 156 | $543 | existing | burst_trading,opposite_side_same_market |
| `0x257a...b42a` | B | 100.0 | 29.4 | 215.8 | 221 | 179 | 99% | 47.1% | 76 | 47.1% | 76 | $345 | existing | opposite_side_same_market |
| `0x1396...62b5` | A | 100.0 | 24.3 | 214.2 | 283 | 76 | 97% | 52.5% | 11 | 50.8% | 13 | $325 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4da5...3baf` | B | 100.0 | 31.6 | 211.9 | 143 | 116 | 100% | 49.4% | 71 | 49.4% | 71 | $946 | existing | opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 60.9% | 297 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.8 | 1h edge -72.59pp over 1 samples | 50.4% | 91 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xde7b...5f4b` | reject-bot | - | D | 47.9 | $11733 | $11733 | 1 | 221.2% | 36 | bot_like_flow,burst_trading |
| `0xbcc3...ac9c` | reject-bot | - | D | 53.8 | $9666 | $9666 | 1 | 44.5% | 65 | bot_like_flow,burst_trading,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 7
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x689f...259c` | reject-bot | - | BOT | 100.0 | 61.5 | 1017.3 | 695 | 219 | 59% | 157.6% | 75 | 155.4% | 195 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xc0f5...3d04` | reject | - | D | 100.0 | 56.7 | 780.9 | 1315 | 569 | 100% | 122.7% | 68 | 122.7% | 68 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x8637...0c86` | reject | - | D | 100.0 | 48.3 | 516.3 | 133 | 76 | 62% | 87.2% | 55 | 83.0% | 79 | sports_tape | bot_like_flow,opposite_side_same_market |
| `0xefde...6505` | watch | - | C | 100.0 | 38.6 | 432.8 | 926 | 192 | 71% | 84.8% | 20 | 86.5% | 37 | holder,sports_tape | burst_trading,open_copy_exposure |
| `0xde7b...5f4b` | reject | - | D | 100.0 | 47.9 | 411.5 | 290 | 23 | 21% | 78.4% | 2 | 221.2% | 36 | existing,sports_tape | bot_like_flow,burst_trading |
| `0x1365...c82a` | reject | - | D | 100.0 | 58.5 | 364.1 | 1437 | 242 | 99% | 67.4% | 35 | 68.3% | 36 | holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbcc3...ac9c` | reject | - | D | 100.0 | 53.8 | 290.8 | 926 | 419 | 68% | 45.4% | 49 | 44.5% | 65 | existing,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 2013 | 177.8% | $+44509.20 | 95.5% | 134.8% | 100.5% | 0.54x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 22 | 1972 | 178.5% | $+43957.39 | 96.1% | 136.4% | 100.5% | 0.51x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 22 | 2002 | 178.0% | $+44366.59 | 95.5% | 136.4% | 100.5% | 0.54x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 21 | 1961 | 178.8% | $+43814.78 | 96.1% | 138.1% | 100.5% | 0.51x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 19 | 1848 | 181.3% | $+42159.08 | 97.3% | 134.8% | 100.5% | 0.50x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 20 | 1911 | 179.5% | $+42995.04 | 96.5% | 132.2% | 100.5% | 0.49x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 19 | 1900 | 179.8% | $+42852.43 | 96.5% | 134.8% | 100.5% | 0.49x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 18 | 1837 | 181.5% | $+42016.47 | 97.3% | 136.4% | 100.5% | 0.49x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 18 | 1826 | 181.8% | $+41772.55 | 97.3% | 132.2% | 100.5% | 0.50x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=20 smart>=70 |
| 10 | 17 | 1815 | 182.0% | $+41629.94 | 97.3% | 134.8% | 100.5% | 0.49x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=20 smart>=70 |
| 11 | 17 | 1918 | 176.4% | $+42376.37 | 95.6% | 131.4% | 100.5% | 0.49x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 40 | 2832 | 151.9% | $+53523.00 | 95.2% | 103.8% | 60.5% | 0.60x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
