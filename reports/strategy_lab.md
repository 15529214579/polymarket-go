# Strategy Lab Report

**Generated:** 2026-08-26 22:54 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (425)
- Valid strategies found: 40
- Candidate layers: 13 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 50 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 13
- Aggregate closed copy trades: 234
- Aggregate copy ROI: 113.1%
- Aggregate copy PnL: $+2928.80
- Aggregate copy win rate: 87.6%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 3.13x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 23.70 | 77 | 65% | 129.0% | $+322.55 | 25 | 100.0% | 48.2% |
| `0x84cd...7565` | A | 23.84 | 121 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.77 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x44c4...09cb` | B | 26.07 | 116 | 62% | 133.8% | $+468.37 | 34 | 82.3% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.22 | 151 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 530
- Aggregate copy ROI: 110.4%
- Aggregate copy PnL: $+7429.34
- Aggregate copy win rate: 83.4%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3c2b...825e` | B | 100.0 | 34.2 | 868.0 | 595 | 96% | 146.0% | $+1664.23 | 98 | 94.9% | $316 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 33.7 | 850.1 | 70 | 29% | 156.4% | $+1188.72 | 59 | 69.5% | $548 | existing,holder |
| `0x4362...9e47` | B | 100.0 | 35.0 | 844.2 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1325 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x5461...92be` | A | 94.7 | 22.6 | 640.9 | 37 | 59% | 158.9% | $+317.89 | 12 | 100.0% | $1171 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.6 | 580.8 | 129 | 98% | 103.4% | $+537.85 | 34 | 88.2% | $793 | existing |
| `0x760f...326a` | B | 100.0 | 30.8 | 563.3 | 15 | 14% | 121.1% | $+290.65 | 20 | 85.0% | $5140 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.5 | 430.5 | 49 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1314 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1176 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.8 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 389.4 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_7d |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1054.3 | 379 | 290 | 36% | 183.6% | $+1835.46 | 94 | 98.9% | 183.4% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.9 | 740.0 | 216 | 79 | 46% | 194.5% | $+330.69 | 16 | 93.8% | 127.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.5 | 617.8 | 1349 | 403 | 100% | 108.0% | $+1295.78 | 45 | 100.0% | 108.0% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 33.7 | 545.6 | 644 | 286 | 93% | 109.7% | $+329.00 | 26 | 100.0% | 108.9% | existing | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 32.8 | 532.2 | 1273 | 320 | 100% | 100.3% | $+1123.32 | 25 | 100.0% | 100.3% | existing,holder | burst_trading,open_copy_exposure |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.6 | 402.0 | 125 | 70 | 30% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.8 | 360.1 | 607 | 313 | 59% | 59.6% | $+137.16 | 21 | 76.2% | 54.7% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 357.1 | 25 | 19 | 96% | 92.2% | $+46.10 | 5 | 80.0% | 92.2% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x168a...0d24` | B | 99.1 | 27.3 | 246.1 | 35 | 46% | 79.7% | 2 | 53 | $29240 | existing,leaderboard_profit_30d,leaderboard_volume_7d | burst_trading |
| `0xb8e3...aaf2` | C | 100.0 | 42.6 | 163.4 | 2 | 0% | 0.0% | 0 | 368 | $597 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 292.1 | 982 | 712 | 96% | 10.2% | 10 | 9.4% | 11 | $506 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.8 | 616 | 382 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.1 | 908 | 392 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 40.3 | 203.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing | burst_trading,open_copy_exposure |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |
| `0x51d3...1719` | B | 100.0 | 27.6 | 184.9 | 177 | 65 | 69% | 20.4% | 18 | 48.4% | 31 | $191 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0x44c4...09cb` | B | 100.0 | 26.1 | 1h edge -15.87pp over 1 samples | 133.8% | 34 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 13 | 234 | 113.1% | $+2928.80 | 87.6% | 115.4% | 65.3% | 3.13x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 12 | 226 | 113.7% | $+2819.29 | 88.1% | 117.6% | 65.3% | 3.15x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 10 | 180 | 121.9% | $+2450.95 | 87.8% | 121.9% | 104.0% | 3.26x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 14 | 265 | 106.3% | $+3114.75 | 85.7% | 112.2% | 54.7% | 4.17x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 17 | 303 | 99.2% | $+3361.61 | 83.8% | 108.4% | 44.4% | 4.06x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 6 | 11 | 216 | 113.2% | $+2615.53 | 88.9% | 115.4% | 65.3% | 2.07x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 13 | 257 | 106.6% | $+3005.24 | 86.0% | 115.4% | 54.7% | 4.23x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 15 | 274 | 104.5% | $+3154.69 | 85.8% | 108.9% | 44.4% | 4.33x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 16 | 294 | 100.7% | $+3321.67 | 83.7% | 108.7% | 51.3% | 3.90x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 12 | 225 | 110.6% | $+2655.47 | 88.9% | 112.2% | 44.4% | 2.36x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 11 | 15 | 285 | 103.2% | $+3239.57 | 84.2% | 108.9% | 54.7% | 4.08x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 16 | 291 | 100.2% | $+3197.49 | 85.2% | 108.7% | 25.2% | 4.13x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
