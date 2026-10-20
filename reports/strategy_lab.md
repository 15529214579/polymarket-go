# Strategy Lab Report

**Generated:** 2026-10-20 23:45 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (429)
- Valid strategies found: 66
- Candidate layers: 16 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 53 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 16
- Aggregate closed copy trades: 1337
- Aggregate copy ROI: 195.5%
- Aggregate copy PnL: $+31242.59
- Aggregate copy win rate: 96.3%
- Median wallet CopyROI: 142.2%
- Worst included wallet CopyROI: 102.5%
- Open copy cost / closed copy capital: 0.72x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x42c9...9623` | A | 12.92 | 0 | 0% | 137.6% | $+2902.28 | 151 | 100.0% | 49.3% |
| `0x5135...7217` | A | 17.99 | 1219 | 100% | 149.4% | $+1703.40 | 88 | 100.0% | 58.5% |
| `0x075d...58c1` | A | 20.13 | 18 | 16% | 173.7% | $+1354.97 | 70 | 90.0% | 38.2% |
| `0x15dd...95ba` | A | 24.99 | 192 | 48% | 117.1% | $+749.44 | 61 | 83.6% | 59.0% |
| `0x612f...05d8` | A | 24.97 | 69 | 45% | 105.9% | $+592.79 | 53 | 94.3% | 53.6% |
| `0x1bd8...e5f8` | A | 24.16 | 22 | 11% | 147.9% | $+547.10 | 37 | 73.0% | 19.3% |
| `0xaca0...9b86` | A | 19.18 | 14 | 16% | 330.6% | $+429.71 | 13 | 100.0% | 55.2% |
| `0x162d...8944` | A | 24.81 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 47.7% |
| `0x7d68...8239` | A | 10.71 | 20 | 24% | 114.7% | $+229.39 | 20 | 95.0% | 31.8% |
| `0x2b17...d36d` | A | 17.63 | 348 | 72% | 119.9% | $+203.76 | 10 | 70.0% | 52.9% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 146.9% | $+132.20 | 9 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.90 | 332 | 31% | 244.8% | $+11336.77 | 435 | 97.9% | 56.9% |
| `0x0f00...73ce` | B | 26.64 | 176 | 29% | 225.7% | $+7199.08 | 195 | 98.5% | 53.1% |
| `0x5be9...ef2b` | B | 27.12 | 137 | 53% | 248.0% | $+3099.79 | 125 | 100.0% | 61.6% |
| `0x092b...614e` | B | 27.37 | 115 | 99% | 102.5% | $+225.51 | 22 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.50 | 60 | 62% | 113.2% | $+203.74 | 17 | 100.0% | 47.1% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 804
- Aggregate copy ROI: 134.8%
- Aggregate copy PnL: $+14438.23
- Aggregate copy win rate: 91.4%
- Worst included CopyROI: 46.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xae2b...90ea` | B | 100.0 | 34.6 | 1448.9 | 187 | 35% | 232.9% | $+6335.38 | 192 | 93.8% | $1896 | existing |
| `0x3e35...4f69` | A | 98.3 | 14.0 | 1091.9 | 50 | 43% | 418.6% | $+293.01 | 7 | 100.0% | $331 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.1 | 834.2 | 60 | 22% | 169.0% | $+743.75 | 43 | 97.7% | $259 | existing |
| `0x760f...326a` | B | 100.0 | 30.6 | 793.8 | 15 | 10% | 126.0% | $+1284.66 | 91 | 93.4% | $8774 | existing |
| `0xddf3...8d34` | A | 100.0 | 2.6 | 751.6 | 23 | 70% | 240.5% | $+216.46 | 7 | 100.0% | $5250 | existing |
| `0x3520...1fa0` | B | 100.0 | 32.5 | 725.4 | 30 | 32% | 172.1% | $+361.41 | 20 | 95.0% | $846 | existing |
| `0x9a0c...370e` | B | 100.0 | 30.6 | 627.3 | 172 | 80% | 81.7% | $+1372.23 | 140 | 92.9% | $320 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 591.0 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6314 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x68cb...6f57` | B | 100.0 | 17.5 | 492.8 | 0 | 0% | 199.8% | $+99.91 | 4 | 100.0% | $394 | existing |
| `0x5925...9984` | B | 100.0 | 31.7 | 487.4 | 48 | 96% | 97.0% | $+174.66 | 15 | 80.0% | $1547 | existing |
| `0xbb4d...de5a` | A | 100.0 | 15.2 | 451.1 | 8 | 36% | 129.8% | $+64.90 | 5 | 100.0% | $7419 | existing |
| `0x789f...20ec` | B | 99.4 | 26.9 | 444.5 | 52 | 96% | 71.6% | $+171.80 | 18 | 94.4% | $1083 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.6 | 420.7 | 143 | 100% | 48.0% | $+293.09 | 59 | 91.5% | $946 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 416.6 | 201 | 91% | 46.4% | $+408.49 | 73 | 78.1% | $1212 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 346.7 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.4 | 747.2 | 853 | 362 | 92% | 125.5% | $+1656.19 | 73 | 100.0% | 122.7% | existing | burst_trading,open_copy_exposure |
| `0x677d...4907` | A | 100.0 | 21.6 | 594.7 | 138 | 112 | 48% | 91.2% | $+538.11 | 58 | 98.3% | 86.2% | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 591.3 | 1148 | 1109 | 100% | 64.6% | $+4458.26 | 297 | 96.6% | 64.6% | existing,holder | burst_trading,opposite_side_same_market |
| `0x7972...594f` | B | 100.0 | 33.3 | 514.3 | 266 | 252 | 67% | 71.0% | $+845.46 | 87 | 55.2% | 56.7% | existing,holder | opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.7 | 489.2 | 178 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.5 | 425.1 | 178 | 67 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | opposite_side_same_market |
| `0x5ba7...8732` | B | 100.0 | 30.9 | 365.6 | 368 | 128 | 54% | 35.0% | $+329.44 | 90 | 84.4% | 48.8% | existing | opposite_side_same_market |
| `0xc16e...8f62` | B | 100.0 | 29.8 | 360.7 | 44 | 20 | 20% | 99.9% | $+69.90 | 5 | 80.0% | 104.7% | existing | burst_trading,open_copy_exposure |
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
| `0xc9bf...81c0` | C | 100.0 | 40.5 | 257.5 | 1376 | 710 | 96% | 115.3% | 14 | 115.3% | 14 | $203 | existing,holder | burst_trading,open_copy_exposure |
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.3 | 1083 | 586 | 98% | 57.4% | 163 | 58.0% | 168 | $601 | existing | opposite_side_same_market |
| `0xd570...e4f8` | C | 100.0 | 43.6 | 244.0 | 1451 | 537 | 100% | 73.5% | 17 | 73.5% | 17 | $1466 | existing | burst_trading |
| `0xe21c...b754` | B | 100.0 | 29.2 | 241.1 | 796 | 420 | 98% | 32.7% | 62 | 30.9% | 65 | $209 | existing | open_copy_exposure,opposite_side_same_market |
| `0x89eb...cff9` | C | 100.0 | 37.1 | 237.1 | 420 | 389 | 89% | 16.2% | 47 | 31.7% | 60 | $1989 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9e07...6e49` | C | 100.0 | 44.8 | 229.1 | 1395 | 502 | 97% | 89.5% | 15 | 89.2% | 21 | $355 | existing,holder | burst_trading,open_copy_exposure |
| `0x4d6e...5cfd` | C | 100.0 | 36.6 | 228.1 | 1190 | 405 | 86% | 336.9% | 5 | 299.0% | 6 | $702 | existing | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 38.4 | 225.1 | 1433 | 374 | 100% | 118.5% | 40 | 118.8% | 42 | $428 | existing | burst_trading |
| `0x257a...b42a` | B | 100.0 | 29.0 | 219.4 | 262 | 206 | 99% | 21.3% | 42 | 21.3% | 42 | $326 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbdd2...6657` | C | 100.0 | 40.3 | 211.6 | 1045 | 367 | 78% | 104.7% | 7 | 110.8% | 9 | $411 | existing,holder | burst_trading,open_copy_exposure |

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
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 1h edge -47.28pp over 1 samples | 64.6% | 297 | existing,holder | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.4% | 73 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9052...3c2d` | reject-bot | - | BOT | 64.6 | $16720 | $26500 | 2 | 132.7% | 68 | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x927c...6633` | reject-bot | - | BOT | 64.4 | $8360 | $8360 | 1 | -3.9% | 59 | bot_like_flow,burst_trading,negative_copy_sim,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 3
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x9052...3c2d` | reject-bot | - | BOT | 100.0 | 64.6 | 577.0 | 705 | 243 | 52% | 106.5% | 26 | 132.7% | 68 | existing,holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x821d...cf16` | reject | - | D | 100.0 | 55.7 | -9.3 | 226 | 27 | 15% | 0.0% | 0 | 0.0% | 0 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure |
| `0x927c...6633` | reject-bot | - | BOT | 100.0 | 64.4 | -11.6 | 651 | 268 | 47% | -2.4% | 22 | -3.9% | 59 | existing,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 16 | 1337 | 195.5% | $+31242.59 | 96.3% | 142.2% | 102.5% | 0.72x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 14 | 1290 | 197.5% | $+30491.73 | 97.2% | 142.2% | 102.5% | 0.54x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 3 | 15 | 1300 | 196.6% | $+30695.49 | 97.0% | 137.6% | 102.5% | 0.70x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 15 | 1328 | 195.8% | $+31110.39 | 96.3% | 137.6% | 102.5% | 0.71x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 13 | 1281 | 197.8% | $+30359.53 | 97.2% | 137.6% | 102.5% | 0.53x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 12 | 1288 | 196.5% | $+30273.18 | 96.4% | 142.7% | 102.5% | 0.49x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 13 | 1305 | 195.5% | $+30476.92 | 96.5% | 137.6% | 102.5% | 0.53x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 14 | 1291 | 196.9% | $+30563.29 | 97.0% | 128.7% | 102.5% | 0.68x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 9 | 11 | 1251 | 197.6% | $+29726.08 | 97.1% | 137.6% | 102.5% | 0.46x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 10 | 12 | 1268 | 196.6% | $+29929.82 | 97.2% | 127.3% | 102.5% | 0.49x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 11 | 21 | 1556 | 178.4% | $+33879.09 | 96.0% | 117.1% | 76.1% | 0.78x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 22 | 1576 | 177.0% | $+33963.19 | 95.9% | 115.9% | 42.0% | 0.84x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
