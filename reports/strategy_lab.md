# Strategy Lab Report

**Generated:** 2026-08-30 22:55 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (427)
- Valid strategies found: 51
- Candidate layers: 15 core + 20 watch + 10 sports + 3 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 52 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 15
- Aggregate closed copy trades: 294
- Aggregate copy ROI: 121.6%
- Aggregate copy PnL: $+3926.86
- Aggregate copy win rate: 90.8%
- Median wallet CopyROI: 119.9%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 2.81x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x5be9...ef2b` | A | 23.63 | 83 | 63% | 191.3% | $+612.17 | 32 | 100.0% | 56.1% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0xb853...0d33` | A | 7.05 | 3 | 14% | 146.1% | $+219.08 | 15 | 100.0% | 59.4% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0x7fcf...80ac` | A | 10.70 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x4607...d542` | B | 26.29 | 104 | 55% | 123.4% | $+567.69 | 43 | 100.0% | 43.6% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.29 | 150 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |
| `0x60b1...8e7d` | B | 26.52 | 336 | 39% | 67.0% | $+60.29 | 8 | 100.0% | 43.3% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 477
- Aggregate copy ROI: 107.6%
- Aggregate copy PnL: $+6638.37
- Aggregate copy win rate: 84.9%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x93b3...0bac` | B | 100.0 | 29.4 | 1340.5 | 16 | 59% | 468.5% | $+468.51 | 10 | 100.0% | $545 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.2 | 916.4 | 619 | 80% | 152.6% | $+1998.47 | 112 | 95.5% | $278 | existing |
| `0x5461...92be` | A | 98.5 | 21.7 | 715.5 | 37 | 55% | 163.7% | $+507.61 | 17 | 100.0% | $1181 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x760f...326a` | B | 100.0 | 30.2 | 511.5 | 13 | 13% | 107.0% | $+235.30 | 19 | 84.2% | $5375 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.3 | 430.0 | 50 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1302 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.5 | 72 | 59% | 50.5% | $+207.07 | 34 | 64.7% | $543 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 379.4 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing,leaderboard_profit_30d,leaderboard_profit_7d |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1124.3 | 382 | 291 | 36% | 197.8% | $+2017.56 | 96 | 99.0% | 186.9% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.4 | 722.6 | 217 | 76 | 48% | 198.1% | $+297.12 | 14 | 92.9% | 128.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 32.3 | 616.8 | 1357 | 392 | 100% | 126.8% | $+925.39 | 27 | 100.0% | 126.8% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 33.3 | 547.6 | 684 | 308 | 93% | 108.8% | $+337.34 | 27 | 100.0% | 108.1% | existing | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.6 | 543.0 | 1394 | 355 | 100% | 89.4% | $+1546.15 | 39 | 100.0% | 89.4% | existing | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 514.9 | 1156 | 1105 | 100% | 69.3% | $+1372.19 | 103 | 95.2% | 69.3% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.4 | 402.8 | 122 | 69 | 31% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 381.4 | 747 | 319 | 71% | 62.6% | $+162.72 | 24 | 83.3% | 53.2% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 3
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x983e...6748` | C | 100.0 | 44.5 | 224.5 | 301 | 74% | 142.3% | 1 | 188 | $16858 | existing,holder,leaderboard_profit_30d,leaderboard_volume_7d | burst_trading,open_copy_exposure |
| `0x6e9f...0eae` | C | 100.0 | 44.5 | 201.5 | 267 | 47% | 100.1% | 49 | 293 | $1127 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 37.0 | 194.4 | 0 | 0% | 123.1% | 1 | 353 | $9804 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 305.7 | 1123 | 825 | 96% | 10.5% | 11 | 9.8% | 12 | $508 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x67ba...246c` | C | 100.0 | 38.7 | 248.6 | 1043 | 563 | 98% | 58.4% | 156 | 59.1% | 161 | $615 | existing | opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.2 | 222.4 | 319 | 182 | 81% | 10.8% | 51 | 11.0% | 60 | $572 | existing,holder | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 216.3 | 619 | 385 | 89% | 29.1% | 86 | 27.6% | 93 | $237 | existing | open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 100.0 | 41.7 | 210.4 | 1363 | 290 | 96% | 122.8% | 39 | 122.0% | 41 | $437 | existing,holder | burst_trading |
| `0xadfa...324e` | C | 100.0 | 40.3 | 206.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing,holder | burst_trading,open_copy_exposure |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.1 | 907 | 392 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 69.3% | 103 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xe907...cff6` | blocked-edge | 15m edge -1.55pp over 26 samples | BOT | 91.2 | 75.5 | -110.1 | 518 | 52 | 41% | -27.5% | 8 | -8.2% | 17 | existing,holder,leaderboard_volume_30d,leaderboard_volume_7d | bot_like_flow,burst_trading,negative_copy_sim,open_copy_exposure |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 15 | 294 | 121.6% | $+3926.86 | 90.8% | 119.9% | 65.3% | 2.81x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 2 | 16 | 302 | 120.1% | $+3975.67 | 90.1% | 117.6% | 61.0% | 3.07x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 14 | 286 | 123.1% | $+3866.57 | 90.6% | 121.6% | 65.3% | 2.82x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 13 | 278 | 124.0% | $+3757.06 | 91.0% | 123.4% | 65.3% | 2.82x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 13 | 276 | 122.5% | $+3613.59 | 92.0% | 123.4% | 65.3% | 1.95x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 10 | 222 | 133.3% | $+3184.96 | 92.3% | 129.6% | 104.0% | 1.79x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 11 | 232 | 132.4% | $+3388.72 | 91.4% | 124.0% | 104.0% | 2.85x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 14 | 285 | 120.2% | $+3653.53 | 91.9% | 119.4% | 44.4% | 2.16x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 12 | 268 | 124.2% | $+3553.30 | 91.8% | 123.7% | 65.3% | 1.93x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=80 closedROI>=0 smart>=70 |
| 10 | 16 | 325 | 115.1% | $+4107.66 | 89.2% | 117.6% | 53.2% | 3.53x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 15 | 317 | 116.3% | $+4047.37 | 89.0% | 119.9% | 53.2% | 3.56x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 17 | 334 | 113.3% | $+4147.60 | 89.2% | 115.4% | 44.4% | 3.67x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
