# Strategy Lab Report

**Generated:** 2026-09-01 23:25 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (429)
- Valid strategies found: 108
- Candidate layers: 18 core + 20 watch + 10 sports + 4 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 56 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 18
- Aggregate closed copy trades: 327
- Aggregate copy ROI: 122.7%
- Aggregate copy PnL: $+4443.09
- Aggregate copy win rate: 90.5%
- Median wallet CopyROI: 121.9%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 2.72x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x5be9...ef2b` | A | 23.37 | 84 | 62% | 192.3% | $+692.35 | 36 | 100.0% | 57.1% |
| `0x162d...8944` | A | 24.85 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x2089...712e` | A | 10.68 | 9 | 20% | 123.9% | $+185.88 | 10 | 100.0% | 43.9% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0xb853...0d33` | A | 7.50 | 4 | 17% | 129.9% | $+181.93 | 14 | 100.0% | 59.4% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0x1bd8...e5f8` | A | 24.62 | 15 | 13% | 192.2% | $+153.77 | 8 | 87.5% | 36.4% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.70 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x4607...d542` | B | 26.23 | 105 | 56% | 126.6% | $+556.99 | 41 | 100.0% | 43.0% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.65 | 163 | 97% | 96.6% | $+202.79 | 20 | 85.0% | 61.1% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.43 | 26 | 27% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |
| `0x60b1...8e7d` | B | 27.35 | 337 | 39% | 72.6% | $+79.88 | 10 | 100.0% | 44.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 481
- Aggregate copy ROI: 101.7%
- Aggregate copy PnL: $+6042.57
- Aggregate copy win rate: 77.8%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 942.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing,leaderboard_profit_7d |
| `0x8cfd...95c4` | B | 100.0 | 33.1 | 790.7 | 69 | 34% | 131.7% | $+1172.41 | 83 | 51.8% | $647 | existing |
| `0x5461...92be` | A | 98.5 | 21.7 | 715.5 | 37 | 55% | 163.7% | $+507.61 | 17 | 100.0% | $1181 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x760f...326a` | B | 100.0 | 30.1 | 578.4 | 13 | 13% | 120.8% | $+314.01 | 23 | 87.0% | $5609 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x6b91...06c4` | B | 100.0 | 26.3 | 524.7 | 1 | 3% | 193.7% | $+116.22 | 5 | 100.0% | $4157 | existing |
| `0x075d...58c1` | A | 100.0 | 23.6 | 472.2 | 7 | 12% | 72.2% | $+252.67 | 35 | 88.6% | $4022 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.3 | 450.6 | 72 | 55% | 63.5% | $+273.08 | 35 | 65.7% | $551 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.2 | 429.7 | 50 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1304 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.8 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x612f...05d8` | B | 100.0 | 12.5 | 396.5 | 10 | 62% | 125.7% | $+37.70 | 3 | 100.0% | $565 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1136.7 | 382 | 291 | 37% | 199.5% | $+2074.97 | 98 | 99.0% | 190.7% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.6 | 722.4 | 217 | 76 | 49% | 198.1% | $+297.12 | 14 | 92.9% | 129.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4cb0...57e0` | B | 100.0 | 33.2 | 624.4 | 691 | 314 | 92% | 125.8% | $+503.34 | 31 | 100.0% | 122.8% | existing,holder | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 31.8 | 593.7 | 1344 | 407 | 100% | 106.4% | $+1287.33 | 38 | 100.0% | 106.4% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.8 | 547.9 | 1408 | 348 | 100% | 86.2% | $+1706.24 | 46 | 100.0% | 86.2% | existing,holder | burst_trading |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 517.5 | 1155 | 1103 | 100% | 69.0% | $+1420.42 | 109 | 95.4% | 69.0% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 403.5 | 799 | 337 | 75% | 65.0% | $+201.38 | 29 | 86.2% | 56.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.5 | 402.7 | 122 | 69 | 32% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 4
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xd570...e4f8` | C | 100.0 | 43.1 | 259.4 | 1466 | 100% | 65.1% | 19 | 586 | $1743 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_30d | burst_trading |
| `0x21f3...35e3` | C | 100.0 | 41.2 | 219.1 | 95 | 8% | 0.0% | 0 | 587 | $583 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xf3a6...71bd` | C | 100.0 | 41.7 | 216.6 | 385 | 31% | 162.1% | 11 | 507 | $988 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xc72d...fc21` | C | 100.0 | 39.7 | 205.9 | 641 | 45% | 0.0% | 0 | 328 | $507 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x67ba...246c` | C | 100.0 | 38.8 | 249.8 | 1062 | 575 | 98% | 57.5% | 161 | 58.1% | 166 | $608 | existing | opposite_side_same_market |
| `0x0467...fb2e` | C | 100.0 | 42.1 | 233.2 | 1295 | 512 | 100% | 12.9% | 9 | 12.9% | 9 | $324 | existing | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.3 | 219.3 | 328 | 186 | 81% | 10.8% | 52 | 9.7% | 62 | $562 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.4 | 217.0 | 625 | 390 | 89% | 28.9% | 87 | 27.4% | 94 | $237 | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | C | 100.0 | 35.6 | 215.9 | 533 | 396 | 70% | 173.5% | 118 | 166.7% | 131 | $266 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 100.0 | 39.7 | 214.6 | 1373 | 318 | 97% | 137.4% | 42 | 137.2% | 43 | $437 | existing | burst_trading |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.2 | 908 | 393 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 40.3 | 203.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing | burst_trading,open_copy_exposure |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |

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

- Wallets: 6
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 25.9 | 1h edge -15.87pp over 1 samples | 135.3% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 69.0% | 109 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 18 | 327 | 122.7% | $+4443.09 | 90.5% | 121.9% | 65.3% | 2.72x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 2 | 17 | 317 | 124.3% | $+4363.21 | 90.2% | 123.9% | 65.3% | 2.72x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 3 | 19 | 362 | 118.3% | $+4695.76 | 90.3% | 119.9% | 65.3% | 2.50x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 18 | 352 | 119.6% | $+4615.88 | 90.1% | 121.9% | 65.3% | 2.50x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 17 | 318 | 121.6% | $+4169.76 | 91.5% | 123.9% | 44.4% | 2.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=80 closedROI>=5 smart>=70 |
| 6 | 17 | 344 | 118.8% | $+4382.49 | 91.3% | 123.9% | 65.3% | 1.79x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 16 | 309 | 123.6% | $+4129.82 | 91.6% | 124.0% | 65.3% | 1.95x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=5 smart>=70 |
| 8 | 14 | 261 | 132.9% | $+3867.92 | 90.8% | 125.3% | 104.0% | 2.73x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 19 | 362 | 116.4% | $+4716.17 | 88.1% | 119.9% | 63.5% | 2.47x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 10 | 18 | 352 | 117.7% | $+4636.29 | 87.8% | 121.9% | 63.5% | 2.47x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=5 smart>=70 |
| 11 | 13 | 251 | 133.7% | $+3664.16 | 91.6% | 126.6% | 104.0% | 1.79x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 20 | 397 | 112.9% | $+4968.84 | 88.2% | 117.6% | 63.5% | 2.30x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
