# Strategy Lab Report

**Generated:** 2026-10-19 23:01 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (429)
- Valid strategies found: 85
- Candidate layers: 18 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 55 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 18
- Aggregate closed copy trades: 1352
- Aggregate copy ROI: 194.9%
- Aggregate copy PnL: $+31648.92
- Aggregate copy win rate: 96.2%
- Median wallet CopyROI: 145.3%
- Worst included wallet CopyROI: 102.5%
- Open copy cost / closed copy capital: 0.76x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x42c9...9623` | A | 12.68 | 0 | 0% | 143.8% | $+3119.68 | 155 | 100.0% | 49.3% |
| `0x5135...7217` | A | 18.68 | 1221 | 100% | 148.9% | $+1697.96 | 85 | 100.0% | 58.5% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x15dd...95ba` | A | 24.99 | 192 | 48% | 117.1% | $+749.44 | 61 | 83.6% | 59.0% |
| `0x1bd8...e5f8` | A | 23.58 | 22 | 11% | 147.9% | $+547.10 | 37 | 73.0% | 19.3% |
| `0xaca0...9b86` | A | 19.13 | 14 | 16% | 330.6% | $+429.71 | 13 | 100.0% | 55.2% |
| `0x162d...8944` | A | 24.81 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 47.7% |
| `0xe4b1...a75e` | A | 24.09 | 0 | 0% | 195.0% | $+312.07 | 12 | 66.7% | 41.8% |
| `0x7d68...8239` | A | 10.71 | 20 | 24% | 114.7% | $+229.39 | 20 | 95.0% | 31.8% |
| `0x2b17...d36d` | A | 17.54 | 367 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 52.9% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 146.9% | $+132.20 | 9 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.91 | 333 | 32% | 244.8% | $+11115.99 | 427 | 98.1% | 56.9% |
| `0x0f00...73ce` | B | 26.64 | 176 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.12 | 136 | 53% | 245.9% | $+3024.91 | 123 | 100.0% | 61.6% |
| `0x612f...05d8` | B | 25.03 | 68 | 45% | 105.9% | $+592.79 | 53 | 94.3% | 53.6% |
| `0x092b...614e` | B | 27.37 | 115 | 99% | 102.5% | $+225.51 | 22 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.50 | 60 | 62% | 113.2% | $+203.74 | 17 | 100.0% | 47.1% |
| `0xc402...c7f1` | B | 29.39 | 40 | 82% | 118.6% | $+177.96 | 12 | 100.0% | 49.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 680
- Aggregate copy ROI: 144.6%
- Aggregate copy PnL: $+13290.74
- Aggregate copy win rate: 90.7%
- Worst included CopyROI: 46.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xae2b...90ea` | B | 100.0 | 34.6 | 1448.9 | 187 | 35% | 232.9% | $+6335.38 | 192 | 93.8% | $1896 | existing |
| `0x3e35...4f69` | A | 98.4 | 13.8 | 1092.9 | 51 | 42% | 418.6% | $+293.01 | 7 | 100.0% | $372 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.0 | 834.3 | 61 | 22% | 169.0% | $+743.75 | 43 | 97.7% | $259 | existing |
| `0x760f...326a` | B | 100.0 | 30.5 | 793.7 | 15 | 10% | 126.0% | $+1284.66 | 91 | 93.4% | $8749 | existing |
| `0xddf3...8d34` | A | 100.0 | 2.9 | 756.9 | 19 | 66% | 240.5% | $+216.46 | 7 | 100.0% | $5221 | existing,holder |
| `0x3520...1fa0` | B | 100.0 | 33.0 | 726.5 | 30 | 33% | 172.1% | $+361.41 | 20 | 95.0% | $856 | existing |
| `0xfa19...bb8f` | B | 100.0 | 30.9 | 600.7 | 42 | 33% | 140.5% | $+224.74 | 16 | 75.0% | $291 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 591.0 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6314 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x68cb...6f57` | B | 100.0 | 17.5 | 492.8 | 0 | 0% | 199.8% | $+99.91 | 4 | 100.0% | $394 | existing |
| `0x5925...9984` | B | 100.0 | 31.7 | 487.4 | 48 | 96% | 97.0% | $+174.66 | 15 | 80.0% | $1547 | existing |
| `0xbb4d...de5a` | A | 100.0 | 15.2 | 451.1 | 8 | 36% | 129.8% | $+64.90 | 5 | 100.0% | $7419 | existing |
| `0x789f...20ec` | B | 99.4 | 26.9 | 444.5 | 52 | 96% | 71.6% | $+171.80 | 18 | 94.4% | $1083 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.6 | 420.7 | 143 | 100% | 48.0% | $+293.09 | 59 | 91.5% | $946 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 416.6 | 201 | 91% | 46.4% | $+408.49 | 73 | 78.1% | $1215 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 346.7 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.1 | 747.2 | 887 | 369 | 92% | 125.5% | $+1656.19 | 73 | 100.0% | 122.7% | existing | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 590.7 | 1148 | 1110 | 100% | 64.4% | $+4461.34 | 299 | 96.7% | 64.4% | existing,holder | burst_trading,opposite_side_same_market |
| `0x677d...4907` | A | 100.0 | 21.8 | 590.2 | 139 | 109 | 50% | 90.5% | $+525.09 | 57 | 98.2% | 85.6% | existing,holder | opposite_side_same_market |
| `0x7972...594f` | B | 100.0 | 33.3 | 518.2 | 256 | 242 | 66% | 72.0% | $+849.07 | 86 | 55.8% | 57.7% | existing | opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.8 | 489.1 | 176 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.5 | 425.1 | 178 | 67 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | opposite_side_same_market |
| `0xc16e...8f62` | B | 100.0 | 30.4 | 360.8 | 39 | 19 | 20% | 99.9% | $+69.90 | 5 | 80.0% | 104.7% | existing | burst_trading,open_copy_exposure |
| `0x8987...02f3` | B | 100.0 | 34.3 | 331.4 | 319 | 123 | 29% | 51.7% | $+242.86 | 20 | 40.0% | 68.4% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x3f3d...2e18` | B | 100.0 | 34.0 | 320.8 | 72 | 41 | 70% | 79.5% | $+111.27 | 5 | 80.0% | 87.5% | existing | burst_trading,open_copy_exposure |
| `0x1944...8088` | A | 100.0 | 23.4 | 315.4 | 39 | 26 | 85% | 48.9% | $+78.15 | 10 | 60.0% | 45.7% | existing | opposite_side_same_market |

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
| `0xc9bf...81c0` | C | 100.0 | 36.4 | 265.7 | 1364 | 714 | 97% | 180.8% | 19 | 180.8% | 19 | $202 | existing,holder | burst_trading,open_copy_exposure |
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.4 | 1085 | 586 | 98% | 57.4% | 163 | 58.0% | 168 | $600 | existing | opposite_side_same_market |
| `0xd570...e4f8` | C | 100.0 | 43.6 | 244.0 | 1451 | 537 | 100% | 73.5% | 17 | 73.5% | 17 | $1466 | existing | burst_trading |
| `0xfdc0...4b85` | C | 100.0 | 35.3 | 242.5 | 1262 | 492 | 100% | 134.2% | 35 | 134.2% | 35 | $197 | existing,holder | burst_trading,open_copy_exposure |
| `0xe21c...b754` | B | 100.0 | 29.2 | 241.1 | 795 | 420 | 98% | 32.7% | 62 | 30.9% | 65 | $209 | existing | open_copy_exposure,opposite_side_same_market |
| `0x89eb...cff9` | C | 100.0 | 37.3 | 233.0 | 413 | 382 | 90% | 16.2% | 47 | 31.7% | 60 | $1944 | existing | open_copy_exposure,opposite_side_same_market |
| `0x9e07...6e49` | C | 100.0 | 45.0 | 230.6 | 1423 | 511 | 99% | 89.6% | 16 | 89.3% | 22 | $348 | existing,holder | burst_trading,open_copy_exposure |
| `0x4d6e...5cfd` | C | 100.0 | 36.8 | 228.6 | 1201 | 410 | 87% | 336.9% | 5 | 299.0% | 6 | $692 | existing | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 38.4 | 225.1 | 1433 | 374 | 100% | 118.5% | 40 | 118.8% | 42 | $428 | existing | burst_trading |
| `0x257a...b42a` | B | 100.0 | 28.9 | 219.8 | 265 | 208 | 99% | 21.3% | 42 | 21.3% | 42 | $325 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 64.4% | 299 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.4% | 73 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9052...3c2d` | reject-bot | - | BOT | 64.7 | $16720 | $26500 | 2 | 134.6% | 69 | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x927c...6633` | reject-bot | - | BOT | 64.8 | $8360 | $8360 | 1 | -5.6% | 58 | bot_like_flow,burst_trading,negative_copy_sim,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 5
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x9052...3c2d` | reject-bot | - | BOT | 100.0 | 64.7 | 595.6 | 691 | 242 | 51% | 109.9% | 27 | 134.6% | 69 | existing,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x3f45...180b` | reject-bot | - | BOT | 100.0 | 62.6 | 242.3 | 34 | 23 | 83% | 109.8% | 4 | 109.8% | 4 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x821d...cf16` | reject | - | D | 100.0 | 52.8 | -4.3 | 217 | 25 | 15% | 0.0% | 0 | 0.0% | 0 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure |
| `0x5cb5...7ce9` | reject | - | D | 100.0 | 59.0 | -10.3 | 49 | 13 | 16% | 0.6% | 1 | 5.2% | 14 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x927c...6633` | reject-bot | - | BOT | 100.0 | 64.8 | -17.0 | 630 | 252 | 45% | -2.5% | 21 | -5.6% | 58 | existing,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 18 | 1352 | 194.9% | $+31648.92 | 96.2% | 145.3% | 102.5% | 0.76x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 17 | 1315 | 196.0% | $+31101.82 | 96.8% | 143.8% | 102.5% | 0.74x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 17 | 1343 | 195.1% | $+31516.72 | 96.1% | 143.8% | 102.5% | 0.75x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 17 | 1340 | 194.9% | $+31336.85 | 96.4% | 143.8% | 102.5% | 0.75x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 15 | 1293 | 196.8% | $+30585.99 | 97.3% | 143.8% | 102.5% | 0.56x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 16 | 1306 | 196.3% | $+30969.62 | 96.8% | 131.8% | 102.5% | 0.73x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 16 | 1303 | 196.0% | $+30789.75 | 97.1% | 131.8% | 102.5% | 0.72x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 8 | 16 | 1331 | 195.2% | $+31204.65 | 96.4% | 131.8% | 102.5% | 0.74x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 13 | 1296 | 195.6% | $+30393.22 | 96.5% | 143.8% | 102.5% | 0.53x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 12 | 1279 | 196.5% | $+30189.48 | 96.5% | 145.8% | 102.5% | 0.50x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 14 | 1284 | 197.1% | $+30453.79 | 97.3% | 131.2% | 102.5% | 0.54x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 11 | 1242 | 197.7% | $+29642.38 | 97.2% | 143.8% | 102.5% | 0.46x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
