# Strategy Lab Report

**Generated:** 2026-09-03 22:49 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (430)
- Valid strategies found: 91
- Candidate layers: 23 core + 20 watch + 10 sports + 6 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 64 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 660
- Aggregate copy ROI: 109.9%
- Aggregate copy PnL: $+8791.86
- Aggregate copy win rate: 93.0%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 44.4%
- Open copy cost / closed copy capital: 1.58x
- Params: tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x42c9...9623` | A | 19.24 | 0 | 0% | 103.5% | $+3291.72 | 221 | 100.0% | 49.4% |
| `0x5be9...ef2b` | A | 23.08 | 88 | 62% | 185.5% | $+723.59 | 39 | 100.0% | 57.2% |
| `0xbbbc...5921` | A | 4.33 | 6 | 7% | 121.9% | $+353.45 | 28 | 96.4% | 43.5% |
| `0x162d...8944` | A | 24.85 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x075d...58c1` | A | 22.73 | 7 | 11% | 74.4% | $+304.90 | 40 | 87.5% | 13.8% |
| `0x2c5b...1dcd` | A | 13.00 | 17 | 22% | 117.3% | $+199.41 | 17 | 100.0% | 47.7% |
| `0x2089...712e` | A | 24.02 | 13 | 16% | 86.1% | $+197.93 | 13 | 92.3% | 41.8% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0x1bd8...e5f8` | A | 24.56 | 16 | 13% | 159.8% | $+143.79 | 9 | 77.8% | 36.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.70 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x2fe2...9645` | A | 18.72 | 34 | 41% | 44.4% | $+39.94 | 9 | 88.9% | 48.2% |
| `0x4607...d542` | B | 25.47 | 103 | 60% | 138.8% | $+527.63 | 36 | 100.0% | 43.0% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x760f...326a` | B | 29.67 | 13 | 12% | 119.0% | $+357.04 | 27 | 88.9% | 48.8% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.60 | 164 | 97% | 96.6% | $+202.79 | 20 | 85.0% | 61.1% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.53 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |
| `0xaa29...9ce5` | B | 25.43 | 350 | 48% | 46.3% | $+60.17 | 13 | 76.9% | 60.8% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 436
- Aggregate copy ROI: 103.4%
- Aggregate copy PnL: $+5698.29
- Aggregate copy win rate: 76.6%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8cfd...95c4` | B | 100.0 | 33.2 | 781.1 | 72 | 34% | 128.6% | $+1182.65 | 86 | 52.3% | $694 | existing |
| `0x5461...92be` | B | 97.3 | 26.9 | 704.7 | 44 | 54% | 160.6% | $+513.82 | 18 | 100.0% | $1077 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x6b91...06c4` | B | 100.0 | 25.6 | 525.4 | 1 | 3% | 193.7% | $+116.22 | 5 | 100.0% | $4073 | existing |
| `0xb939...da70` | A | 100.0 | 3.5 | 504.3 | 0 | 0% | 157.0% | $+156.98 | 6 | 83.3% | $5523 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.3 | 450.6 | 72 | 55% | 63.5% | $+273.08 | 35 | 65.7% | $551 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.2 | 429.7 | 50 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1304 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.8 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x612f...05d8` | B | 100.0 | 12.5 | 396.5 | 10 | 62% | 125.7% | $+37.70 | 3 | 100.0% | $565 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1136.8 | 382 | 290 | 37% | 199.5% | $+2074.97 | 98 | 99.0% | 190.8% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.8 | 722.1 | 217 | 76 | 49% | 198.1% | $+297.12 | 14 | 92.9% | 131.5% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.9 | 698.8 | 1344 | 423 | 100% | 137.9% | $+1503.36 | 36 | 100.0% | 137.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.5 | 622.5 | 782 | 339 | 93% | 122.5% | $+539.00 | 34 | 100.0% | 119.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.6 | 530.0 | 1407 | 358 | 100% | 80.8% | $+1592.84 | 48 | 100.0% | 80.8% | existing,holder | burst_trading |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.4 | 369.5 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.8 | 359.0 | 118 | 65 | 32% | 63.4% | $+107.79 | 16 | 68.8% | 85.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 357.1 | 25 | 19 | 96% | 92.2% | $+46.10 | 5 | 80.0% | 92.2% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 6
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xd570...e4f8` | C | 100.0 | 43.6 | 253.2 | 1463 | 100% | 73.5% | 17 | 539 | $1458 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_30d | burst_trading |
| `0x983e...6748` | C | 100.0 | 44.0 | 242.2 | 301 | 66% | 142.3% | 1 | 213 | $16576 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0xc72d...fc21` | C | 100.0 | 39.7 | 207.4 | 600 | 43% | 0.0% | 0 | 341 | $531 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x684b...8409` | C | 100.0 | 37.0 | 194.3 | 86 | 6% | 0.0% | 0 | 285 | $504 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xd3a0...f136` | C | 99.3 | 44.5 | 193.5 | 106 | 43% | 100.2% | 4 | 82 | $6080 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xb8e3...aaf2` | C | 100.0 | 42.1 | 182.7 | 2 | 0% | -2.2% | 1 | 659 | $607 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x67ba...246c` | C | 100.0 | 38.7 | 251.2 | 1084 | 585 | 98% | 57.4% | 163 | 58.0% | 168 | $599 | existing | opposite_side_same_market |
| `0x0467...fb2e` | C | 100.0 | 42.2 | 237.4 | 1291 | 543 | 100% | 14.3% | 7 | 14.3% | 7 | $359 | existing | burst_trading,open_copy_exposure |
| `0x8e74...5972` | C | 100.0 | 36.0 | 229.1 | 1369 | 349 | 97% | 135.6% | 48 | 135.5% | 49 | $507 | existing,holder | burst_trading |
| `0x54f0...7725` | A | 100.0 | 24.1 | 220.0 | 328 | 186 | 81% | 10.8% | 52 | 9.7% | 61 | $561 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.3 | 219.1 | 641 | 404 | 89% | 28.9% | 91 | 27.4% | 98 | $239 | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | C | 100.0 | 36.3 | 206.3 | 840 | 342 | 78% | 65.0% | 29 | 63.5% | 32 | $160 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.3 | 910 | 393 | 80% | 68.5% | 54 | 71.8% | 63 | $287 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
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

- Wallets: 5
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 25.9 | 1h edge -15.87pp over 1 samples | 135.3% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 660 | 109.9% | $+8791.86 | 93.0% | 115.4% | 44.4% | 1.58x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 2 | 22 | 620 | 111.8% | $+8486.96 | 93.4% | 116.3% | 44.4% | 1.65x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 3 | 22 | 651 | 110.6% | $+8751.92 | 93.1% | 116.3% | 46.3% | 1.51x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 23 | 640 | 110.0% | $+8653.94 | 92.2% | 115.4% | 46.3% | 1.59x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 22 | 658 | 110.3% | $+8816.57 | 92.6% | 116.3% | 59.4% | 1.56x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 6 | 21 | 611 | 112.6% | $+8447.02 | 93.5% | 117.3% | 46.3% | 1.58x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=20 smart>=70 |
| 7 | 21 | 638 | 111.7% | $+8691.75 | 93.4% | 117.3% | 65.3% | 1.53x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 21 | 618 | 112.3% | $+8511.67 | 92.9% | 117.3% | 59.4% | 1.64x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 9 | 23 | 693 | 108.0% | $+9089.65 | 91.2% | 115.4% | 59.4% | 1.51x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 23 | 684 | 108.0% | $+8922.88 | 91.7% | 115.4% | 39.4% | 1.56x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 20 | 598 | 113.8% | $+8386.85 | 93.8% | 118.2% | 65.3% | 1.60x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=20 smart>=70 |
| 12 | 22 | 673 | 109.2% | $+8964.83 | 92.0% | 116.3% | 63.5% | 1.47x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
