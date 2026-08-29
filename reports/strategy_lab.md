# Strategy Lab Report

**Generated:** 2026-08-29 23:17 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (426)
- Valid strategies found: 49
- Candidate layers: 16 core + 20 watch + 10 sports + 3 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 53 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 16
- Aggregate closed copy trades: 309
- Aggregate copy ROI: 120.1%
- Aggregate copy PnL: $+4093.78
- Aggregate copy win rate: 89.0%
- Median wallet CopyROI: 117.6%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 2.92x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x5be9...ef2b` | A | 23.50 | 81 | 64% | 183.9% | $+551.81 | 30 | 100.0% | 52.6% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.70 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x4607...d542` | B | 26.47 | 102 | 55% | 124.5% | $+560.16 | 42 | 100.0% | 44.0% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x760f...326a` | B | 29.95 | 12 | 12% | 107.0% | $+235.30 | 19 | 84.2% | 45.0% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.29 | 150 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.43 | 26 | 27% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 498
- Aggregate copy ROI: 112.1%
- Aggregate copy PnL: $+7129.57
- Aggregate copy win rate: 85.7%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x93b3...0bac` | B | 98.4 | 30.3 | 1401.6 | 14 | 61% | 518.0% | $+466.15 | 9 | 100.0% | $466 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.0 | 920.3 | 608 | 85% | 152.6% | $+1998.47 | 112 | 95.5% | $290 | existing,holder |
| `0x4362...9e47` | B | 100.0 | 35.0 | 844.2 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1325 | existing |
| `0x5461...92be` | A | 98.5 | 21.7 | 715.5 | 37 | 55% | 163.7% | $+507.61 | 17 | 100.0% | $1181 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 590.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing,holder |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.3 | 430.0 | 50 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1302 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.8 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.5 | 72 | 59% | 50.5% | $+207.07 | 34 | 64.7% | $543 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 389.4 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_7d |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1124.4 | 381 | 291 | 36% | 197.8% | $+2017.56 | 96 | 99.0% | 187.6% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.5 | 722.6 | 214 | 75 | 47% | 198.1% | $+297.12 | 14 | 92.9% | 127.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.3 | 678.5 | 1350 | 395 | 100% | 136.5% | $+1078.43 | 33 | 100.0% | 136.5% | existing,holder | burst_trading,open_copy_exposure |
| `0x2829...12a0` | B | 93.6 | 33.4 | 559.6 | 211 | 110 | 88% | 124.0% | $+558.14 | 18 | 100.0% | 124.0% | existing | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 33.5 | 547.4 | 673 | 302 | 93% | 108.8% | $+337.34 | 27 | 100.0% | 108.1% | existing | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.7 | 526.0 | 1393 | 356 | 100% | 88.8% | $+1376.43 | 34 | 100.0% | 88.8% | existing,holder | burst_trading,open_copy_exposure |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.4 | 402.8 | 122 | 69 | 31% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.5 | 372.2 | 695 | 314 | 67% | 61.7% | $+148.03 | 22 | 81.8% | 51.9% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 3
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x168a...0d24` | C | 99.9 | 41.4 | 242.5 | 66 | 35% | 79.7% | 2 | 133 | $20538 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_7d | burst_trading,open_copy_exposure |
| `0x65c1...2988` | C | 100.0 | 40.7 | 217.6 | 1306 | 91% | 0.0% | 0 | 347 | $509 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 36.8 | 193.7 | 0 | 0% | 123.1% | 1 | 340 | $9920 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 302.9 | 1094 | 802 | 96% | 10.2% | 10 | 9.4% | 11 | $510 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x67ba...246c` | C | 100.0 | 38.6 | 247.9 | 1031 | 555 | 98% | 58.1% | 153 | 58.8% | 158 | $620 | existing | opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.0 | 220.2 | 303 | 177 | 83% | 10.8% | 51 | 11.0% | 60 | $599 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.9 | 617 | 383 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing | open_copy_exposure,opposite_side_same_market |
| `0xaa4d...0bd2` | C | 100.0 | 44.9 | 207.4 | 832 | 220 | 86% | 110.1% | 11 | 106.4% | 13 | $382 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.1 | 907 | 392 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 100.0 | 42.0 | 203.5 | 1250 | 262 | 96% | 122.6% | 37 | 121.8% | 39 | $448 | existing | burst_trading |
| `0xadfa...324e` | C | 100.0 | 40.3 | 203.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing | burst_trading,open_copy_exposure |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x076d...8d4c` | C | 100.0 | 42.2 | 201.7 | 313 | 117 | 27% | 116.4% | 34 | 149.4% | 108 | $296 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading |

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
| `0x076d...8d4c` | C | 100.0 | 42.2 | 15m edge -6.00pp over 4 samples | 149.4% | 108 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading |
| `0x44c4...09cb` | B | 100.0 | 25.9 | 1h edge -15.87pp over 1 samples | 135.3% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x076d...8d4c` | blocked-edge | 15m edge -6.00pp over 4 samples | C | 100.0 | 42.2 | 719.3 | 313 | 117 | 27% | 116.4% | 34 | 149.4% | 108 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 16 | 309 | 120.1% | $+4093.78 | 89.0% | 117.6% | 65.3% | 2.92x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 13 | 255 | 127.8% | $+3615.93 | 89.4% | 124.0% | 104.0% | 2.98x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 15 | 301 | 120.7% | $+3984.27 | 89.4% | 119.9% | 65.3% | 2.93x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 12 | 245 | 128.3% | $+3412.17 | 90.2% | 124.2% | 104.0% | 2.03x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 14 | 291 | 120.8% | $+3780.51 | 90.0% | 119.7% | 65.3% | 2.12x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 17 | 338 | 114.2% | $+4259.89 | 87.6% | 115.4% | 51.9% | 3.69x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 15 | 300 | 118.6% | $+3820.45 | 90.0% | 115.4% | 44.4% | 2.33x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 8 | 16 | 330 | 114.7% | $+4150.38 | 87.9% | 117.6% | 51.9% | 3.73x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 18 | 347 | 112.6% | $+4299.83 | 87.6% | 112.2% | 44.4% | 3.83x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 18 | 358 | 111.3% | $+4384.71 | 86.3% | 112.2% | 51.9% | 3.65x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 17 | 350 | 111.6% | $+4275.20 | 86.6% | 115.4% | 51.9% | 3.68x | tier>=B bot<30 copyT>=10 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 17 | 347 | 110.6% | $+4193.18 | 87.3% | 115.4% | 25.2% | 3.59x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
