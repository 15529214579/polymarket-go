# Strategy Lab Report

**Generated:** 2026-08-25 23:16 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (423)
- Valid strategies found: 55
- Candidate layers: 15 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 50 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 15
- Aggregate closed copy trades: 253
- Aggregate copy ROI: 117.5%
- Aggregate copy PnL: $+3338.28
- Aggregate copy win rate: 87.7%
- Median wallet CopyROI: 119.9%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 4.74x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 23.85 | 76 | 65% | 129.0% | $+322.55 | 25 | 100.0% | 48.2% |
| `0x84cd...7565` | A | 23.84 | 121 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.77 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x44c4...09cb` | B | 26.23 | 119 | 62% | 133.3% | $+479.87 | 35 | 82.9% | 55.5% |
| `0x0a7c...8964` | B | 27.13 | 953 | 82% | 191.7% | $+249.23 | 8 | 100.0% | 56.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.22 | 151 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 491
- Aggregate copy ROI: 109.3%
- Aggregate copy PnL: $+6829.13
- Aggregate copy win rate: 82.9%
- Worst included CopyROI: 45.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3b5a...87ef` | B | 100.0 | 33.0 | 876.1 | 60 | 27% | 167.3% | $+1170.90 | 53 | 73.6% | $538 | existing |
| `0x4362...9e47` | B | 100.0 | 35.0 | 844.2 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1325 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.9 | 800.6 | 578 | 96% | 143.3% | $+1117.96 | 68 | 92.7% | $317 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x5461...92be` | A | 94.4 | 23.2 | 645.2 | 34 | 61% | 158.9% | $+317.89 | 12 | 100.0% | $1113 | existing,holder |
| `0x0d66...ff1d` | B | 100.0 | 34.1 | 584.1 | 123 | 98% | 106.1% | $+519.73 | 32 | 87.5% | $816 | existing |
| `0x760f...326a` | B | 100.0 | 30.7 | 563.3 | 15 | 14% | 121.1% | $+290.65 | 20 | 85.0% | $5102 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.5 | 430.5 | 49 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1314 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 419.1 | 235 | 97% | 45.9% | $+491.54 | 78 | 76.9% | $1176 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.7 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 389.4 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_7d |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1054.3 | 379 | 290 | 36% | 183.6% | $+1835.46 | 94 | 98.9% | 183.4% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.9 | 740.1 | 216 | 79 | 46% | 194.5% | $+330.69 | 16 | 93.8% | 127.4% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 33.5 | 600.7 | 1366 | 361 | 100% | 106.2% | $+1253.27 | 41 | 100.0% | 106.2% | existing,holder | burst_trading,open_copy_exposure |
| `0x076d...8d4c` | B | 100.0 | 32.7 | 548.5 | 144 | 33 | 12% | 139.4% | $+250.84 | 11 | 100.0% | 149.9% | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading |
| `0x4cb0...57e0` | B | 100.0 | 33.8 | 545.6 | 643 | 285 | 93% | 109.7% | $+329.00 | 26 | 100.0% | 108.9% | existing | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.0 | 530.6 | 1227 | 302 | 100% | 102.4% | $+1095.53 | 23 | 100.0% | 102.4% | existing,holder | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 475.4 | 1156 | 1104 | 100% | 63.7% | $+1076.22 | 87 | 93.1% | 63.7% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.6 | 402.1 | 125 | 70 | 29% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x168a...0d24` | B | 99.1 | 27.3 | 246.1 | 35 | 46% | 79.7% | 2 | 53 | $29240 | existing,leaderboard_profit_30d,leaderboard_volume_7d | burst_trading |
| `0x65c1...2988` | C | 100.0 | 36.3 | 210.1 | 1360 | 96% | 395.0% | 1 | 380 | $508 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 287.8 | 969 | 702 | 96% | 10.2% | 10 | 9.4% | 11 | $501 | existing | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.8 | 616 | 382 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 208.7 | 916 | 396 | 80% | 68.5% | 54 | 71.8% | 63 | $287 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.7 | 205.3 | 586 | 306 | 57% | 64.7% | 19 | 57.7% | 29 | $190 | existing | open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 40.3 | 203.4 | 1051 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $388 | existing | burst_trading,open_copy_exposure |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |

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
| `0x076d...8d4c` | B | 100.0 | 32.7 | 15m edge -6.00pp over 4 samples | 149.9% | 104 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading |
| `0x44c4...09cb` | B | 100.0 | 26.2 | 1h edge -15.87pp over 1 samples | 133.3% | 35 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.4 | 1h edge -47.28pp over 1 samples | 63.7% | 87 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 45.9% | 78 | existing | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 15 | 253 | 117.5% | $+3338.28 | 87.7% | 119.9% | 65.3% | 4.74x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 12 | 199 | 126.6% | $+2860.43 | 87.9% | 126.5% | 104.0% | 5.27x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 16 | 282 | 111.5% | $+3522.85 | 86.2% | 117.6% | 57.7% | 5.59x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 13 | 235 | 118.2% | $+3025.01 | 88.9% | 124.0% | 65.3% | 3.96x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 11 | 189 | 127.1% | $+2656.67 | 88.9% | 129.0% | 104.0% | 4.25x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 17 | 291 | 109.6% | $+3562.79 | 86.3% | 115.4% | 44.4% | 5.70x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 17 | 302 | 108.2% | $+3647.67 | 84.8% | 115.4% | 57.7% | 5.42x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 14 | 244 | 115.7% | $+3064.95 | 88.9% | 119.7% | 44.4% | 4.16x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 9 | 18 | 311 | 105.7% | $+3729.77 | 84.2% | 112.2% | 51.3% | 5.19x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 18 | 308 | 105.4% | $+3605.59 | 85.7% | 112.2% | 25.2% | 5.45x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 19 | 320 | 104.1% | $+3769.71 | 84.4% | 108.9% | 44.4% | 5.30x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 18 | 328 | 103.1% | $+3753.98 | 83.5% | 112.2% | 39.4% | 5.11x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
