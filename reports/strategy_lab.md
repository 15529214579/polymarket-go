# Strategy Lab Report

**Generated:** 2026-09-04 23:42 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (430)
- Valid strategies found: 84
- Candidate layers: 23 core + 20 watch + 10 sports + 5 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 61 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 704
- Aggregate copy ROI: 114.1%
- Aggregate copy PnL: $+11165.76
- Aggregate copy win rate: 93.0%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 44.4%
- Open copy cost / closed copy capital: 1.25x
- Params: tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x42c9...9623` | A | 19.05 | 0 | 0% | 119.4% | $+4119.72 | 224 | 100.0% | 49.4% |
| `0x78be...bde0` | A | 18.73 | 0 | 0% | 106.7% | $+2005.62 | 82 | 98.8% | 46.7% |
| `0x5be9...ef2b` | A | 22.91 | 88 | 60% | 181.8% | $+745.32 | 41 | 100.0% | 57.0% |
| `0x162d...8944` | A | 24.85 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x2c5b...1dcd` | A | 12.94 | 17 | 22% | 117.3% | $+199.41 | 17 | 100.0% | 47.7% |
| `0x2089...712e` | A | 24.02 | 13 | 16% | 86.1% | $+197.93 | 13 | 92.3% | 41.8% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0x1bd8...e5f8` | A | 24.46 | 16 | 12% | 159.8% | $+159.80 | 10 | 80.0% | 36.4% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x2fe2...9645` | A | 18.72 | 34 | 41% | 44.4% | $+39.94 | 9 | 88.9% | 48.2% |
| `0x4607...d542` | B | 25.34 | 98 | 60% | 144.9% | $+507.23 | 33 | 100.0% | 43.0% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x760f...326a` | B | 29.60 | 13 | 12% | 116.0% | $+440.92 | 33 | 90.9% | 47.9% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.78 | 161 | 97% | 96.6% | $+202.79 | 20 | 85.0% | 61.1% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.31 | 26 | 25% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |
| `0x0931...e78e` | B | 25.83 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | 26.9% |
| `0x4dee...8ad7` | B | 27.07 | 29 | 71% | 51.3% | $+82.10 | 9 | 66.7% | 24.6% |
| `0xaa29...9ce5` | B | 25.37 | 351 | 49% | 46.3% | $+60.17 | 13 | 76.9% | 60.8% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 483
- Aggregate copy ROI: 101.4%
- Aggregate copy PnL: $+6134.34
- Aggregate copy win rate: 77.8%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8cfd...95c4` | B | 100.0 | 33.3 | 816.6 | 72 | 33% | 134.9% | $+1308.56 | 90 | 54.4% | $690 | existing |
| `0x5461...92be` | B | 99.5 | 26.5 | 678.1 | 79 | 66% | 140.7% | $+605.05 | 22 | 100.0% | $997 | existing,holder |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.7 | 562.7 | 117 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $575 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x6b91...06c4` | B | 100.0 | 25.6 | 525.4 | 1 | 3% | 193.7% | $+116.22 | 5 | 100.0% | $4073 | existing |
| `0x075d...58c1` | A | 100.0 | 22.4 | 494.6 | 7 | 10% | 75.1% | $+315.27 | 41 | 87.8% | $3889 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x7fcf...80ac` | A | 100.0 | 10.8 | 456.5 | 18 | 34% | 109.6% | $+109.59 | 7 | 85.7% | $1771 | existing |
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
| `0x17e2...d472` | B | 100.0 | 31.3 | 1138.1 | 381 | 290 | 36% | 198.8% | $+2126.93 | 100 | 99.0% | 189.6% | existing | opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.6 | 737.4 | 1340 | 413 | 100% | 137.1% | $+1851.07 | 46 | 100.0% | 137.1% | existing,holder | burst_trading,open_copy_exposure |
| `0x7e33...73e7` | B | 100.0 | 34.7 | 722.8 | 216 | 75 | 49% | 198.1% | $+297.12 | 14 | 92.9% | 138.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4cb0...57e0` | B | 100.0 | 34.2 | 639.5 | 813 | 347 | 93% | 122.8% | $+700.19 | 38 | 100.0% | 121.0% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.5 | 534.0 | 1402 | 368 | 100% | 80.2% | $+1627.25 | 52 | 100.0% | 80.2% | existing,holder | burst_trading |
| `0x076d...8d4c` | B | 100.0 | 33.7 | 502.4 | 285 | 79 | 24% | 83.3% | $+449.91 | 33 | 100.0% | 142.5% | existing,holder,leaderboard_profit_all,leaderboard_volume_30d | burst_trading |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.0 | 384.7 | 116 | 64 | 32% | 73.6% | $+117.75 | 15 | 73.3% | 87.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 357.1 | 25 | 19 | 96% | 92.2% | $+46.10 | 5 | 80.0% | 92.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 5
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xd570...e4f8` | C | 100.0 | 43.6 | 253.3 | 1462 | 100% | 73.5% | 17 | 539 | $1459 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_30d | burst_trading |
| `0x983e...6748` | C | 100.0 | 44.0 | 242.2 | 301 | 66% | 142.3% | 1 | 213 | $16576 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0xc72d...fc21` | C | 100.0 | 38.6 | 210.5 | 523 | 37% | 0.0% | 0 | 351 | $559 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xb8e3...aaf2` | C | 100.0 | 42.0 | 202.5 | 12 | 1% | -2.2% | 1 | 702 | $606 | leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xd3a0...f136` | C | 99.3 | 44.5 | 193.5 | 106 | 43% | 100.2% | 4 | 82 | $6080 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.2 | 1084 | 585 | 98% | 57.4% | 163 | 58.0% | 168 | $599 | existing | opposite_side_same_market |
| `0x0467...fb2e` | C | 100.0 | 42.4 | 233.7 | 1284 | 517 | 100% | 14.3% | 7 | 14.3% | 7 | $349 | existing | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 36.6 | 226.8 | 1366 | 336 | 97% | 138.2% | 49 | 138.0% | 50 | $528 | existing,holder | burst_trading |
| `0x54f0...7725` | A | 100.0 | 24.2 | 220.0 | 331 | 186 | 82% | 10.8% | 52 | 9.7% | 61 | $558 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.3 | 219.4 | 643 | 406 | 89% | 28.9% | 91 | 27.4% | 98 | $241 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | B | 100.0 | 34.5 | 215.3 | 342 | 298 | 69% | 18.4% | 84 | 27.6% | 114 | $1028 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | C | 100.0 | 36.1 | 209.5 | 877 | 355 | 81% | 67.8% | 30 | 67.7% | 32 | $157 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.1 | 205.4 | 918 | 393 | 81% | 68.3% | 56 | 70.4% | 64 | $285 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 40.3 | 203.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing | burst_trading,open_copy_exposure |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |

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

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x076d...8d4c` | B | 100.0 | 33.7 | 15m edge -6.00pp over 4 samples | 142.5% | 101 | existing,holder,leaderboard_profit_all,leaderboard_volume_30d | burst_trading |
| `0x44c4...09cb` | B | 100.0 | 25.9 | 1h edge -15.87pp over 1 samples | 135.3% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.1 | 1h edge -52.95pp over 1 samples | 70.4% | 64 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x1944...8088` | prompt | - | A | 100.0 | 24.6 | 242.3 | 36 | 24 | 90% | 55.0% | 9 | 55.0% | 9 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 704 | 114.1% | $+11165.76 | 93.0% | 115.4% | 44.4% | 1.25x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 2 | 22 | 716 | 114.6% | $+11274.11 | 93.9% | 115.7% | 44.4% | 1.19x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 23 | 762 | 112.0% | $+11632.07 | 91.9% | 115.4% | 46.3% | 1.14x | tier>=B bot<30 copyT>=10 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 22 | 727 | 114.0% | $+11358.99 | 93.1% | 115.7% | 46.3% | 1.17x | tier>=B bot<30 copyT>=10 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 5 | 22 | 695 | 114.7% | $+11125.82 | 93.1% | 115.7% | 46.3% | 1.19x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 24 | 745 | 112.2% | $+11523.54 | 92.5% | 112.2% | 46.3% | 1.14x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 7 | 21 | 675 | 116.3% | $+10958.84 | 94.2% | 116.0% | 44.4% | 1.23x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 8 | 23 | 721 | 111.8% | $+11111.72 | 93.2% | 115.4% | 23.4% | 1.15x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 9 | 21 | 707 | 115.2% | $+11234.17 | 93.9% | 116.0% | 46.3% | 1.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 22 | 749 | 112.8% | $+11571.90 | 92.1% | 115.7% | 59.4% | 1.15x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 21 | 686 | 115.8% | $+11043.72 | 93.4% | 116.0% | 46.3% | 1.21x | tier>=B bot<30 copyT>=10 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 12 | 21 | 714 | 114.9% | $+11298.82 | 93.4% | 116.0% | 59.4% | 1.18x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
