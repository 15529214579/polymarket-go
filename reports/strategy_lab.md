# Strategy Lab Report

**Generated:** 2026-08-24 23:41 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (421)
- Valid strategies found: 65
- Candidate layers: 10 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 47 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 10
- Aggregate closed copy trades: 180
- Aggregate copy ROI: 122.1%
- Aggregate copy PnL: $+2455.08
- Aggregate copy win rate: 87.8%
- Median wallet CopyROI: 121.9%
- Worst included wallet CopyROI: 104.0%
- Open copy cost / closed copy capital: 3.25x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 23.78 | 74 | 64% | 131.3% | $+315.18 | 24 | 100.0% | 47.4% |
| `0x84cd...7565` | A | 23.84 | 121 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 26.29 | 118 | 62% | 133.3% | $+479.87 | 35 | 82.9% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 554
- Aggregate copy ROI: 102.5%
- Aggregate copy PnL: $+7153.44
- Aggregate copy win rate: 83.2%
- Worst included CopyROI: 46.5%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x4362...9e47` | B | 100.0 | 35.0 | 844.2 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1325 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.9 | 803.6 | 571 | 96% | 143.3% | $+1117.96 | 68 | 92.7% | $315 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 33.2 | 781.6 | 60 | 29% | 150.7% | $+904.09 | 45 | 68.9% | $503 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.1 | 584.1 | 123 | 98% | 106.1% | $+519.73 | 32 | 87.5% | $816 | existing |
| `0x760f...326a` | B | 100.0 | 30.7 | 563.3 | 15 | 14% | 121.1% | $+290.65 | 20 | 85.0% | $5102 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x5461...92be` | A | 100.0 | 19.1 | 559.2 | 15 | 56% | 155.6% | $+186.77 | 7 | 100.0% | $1202 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x18c2...529a` | B | 100.0 | 27.9 | 504.0 | 912 | 80% | 72.1% | $+619.78 | 63 | 92.1% | $307 | existing,holder |
| `0x092b...614e` | B | 100.0 | 25.2 | 494.4 | 151 | 97% | 97.5% | $+185.35 | 18 | 83.3% | $1020 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x7fcf...80ac` | A | 100.0 | 10.8 | 443.2 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | $1826 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.5 | 430.5 | 49 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1314 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 31.3 | 421.1 | 233 | 98% | 46.5% | $+492.83 | 77 | 77.9% | $1186 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.6 | 1049.7 | 389 | 299 | 38% | 181.9% | $+1855.03 | 96 | 97.9% | 182.5% | existing | opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.9 | 734.5 | 644 | 229 | 52% | 149.0% | $+1266.57 | 34 | 100.0% | 132.2% | existing | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 34.2 | 616.6 | 1370 | 358 | 100% | 111.6% | $+1272.13 | 40 | 100.0% | 111.6% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 33.8 | 545.5 | 638 | 281 | 93% | 109.7% | $+329.00 | 26 | 100.0% | 108.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.0 | 529.6 | 1222 | 300 | 100% | 105.7% | $+1036.10 | 21 | 100.0% | 105.7% | existing,holder | burst_trading,open_copy_exposure |
| `0x21cc...54bc` | B | 100.0 | 28.7 | 424.6 | 170 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.9 | 406.0 | 127 | 72 | 29% | 75.4% | $+150.82 | 18 | 72.2% | 86.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.7 | 370.4 | 576 | 304 | 56% | 64.7% | $+135.77 | 19 | 79.0% | 57.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 357.1 | 25 | 19 | 96% | 92.2% | $+46.10 | 5 | 80.0% | 92.2% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x168a...0d24` | B | 99.1 | 27.3 | 246.1 | 35 | 46% | 79.7% | 2 | 53 | $29240 | existing,leaderboard_profit_30d,leaderboard_volume_7d | burst_trading |
| `0x65c1...2988` | C | 100.0 | 36.3 | 211.1 | 1360 | 96% | 395.0% | 1 | 389 | $518 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 290.3 | 961 | 697 | 96% | 10.2% | 10 | 9.4% | 11 | $504 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x6db5...e279` | B | 100.0 | 29.5 | 220.2 | 667 | 318 | 61% | 23.7% | 6 | 37.7% | 16 | $1171 | existing | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.8 | 616 | 382 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 201.2 | 422 | 232 | 75% | 39.0% | 50 | 34.8% | 62 | $270 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |
| `0x51d3...1719` | B | 100.0 | 27.6 | 184.9 | 177 | 65 | 69% | 20.4% | 18 | 48.4% | 31 | $191 | existing | open_copy_exposure,opposite_side_same_market |
| `0x06ee...66b5` | B | 100.0 | 26.6 | 176.0 | 86 | 66 | 38% | 11.2% | 24 | 27.7% | 47 | $464 | existing | opposite_side_same_market |

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
| `0x44c4...09cb` | B | 100.0 | 26.3 | 1h edge -15.87pp over 1 samples | 133.3% | 35 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 72.1% | 63 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 31.3 | 1h edge -72.59pp over 1 samples | 46.5% | 77 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3b62...5cc8` | reject | - | D | 100.0 | 57.8 | 129.6 | 590 | 420 | 53% | 28.0% | 10 | 32.8% | 13 | existing,leaderboard_profit_7d,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 10 | 180 | 122.1% | $+2455.08 | 87.8% | 121.9% | 104.0% | 3.25x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 14 | 297 | 103.0% | $+3552.71 | 88.6% | 112.2% | 65.3% | 3.32x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 13 | 289 | 103.1% | $+3443.20 | 88.9% | 115.4% | 65.3% | 3.34x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 17 | 355 | 95.3% | $+3944.20 | 85.4% | 108.4% | 51.3% | 3.94x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 17 | 345 | 95.8% | $+3851.75 | 86.7% | 108.4% | 44.4% | 4.16x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 18 | 364 | 94.2% | $+3984.14 | 85.4% | 106.2% | 44.4% | 4.07x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 16 | 335 | 97.9% | $+3777.22 | 87.2% | 108.7% | 44.4% | 4.29x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 8 | 18 | 365 | 93.5% | $+4018.73 | 84.9% | 106.2% | 46.6% | 3.84x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 16 | 346 | 97.0% | $+3862.10 | 85.8% | 108.7% | 57.7% | 4.08x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 19 | 374 | 92.5% | $+4058.67 | 85.0% | 104.0% | 44.4% | 3.96x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 16 | 336 | 97.0% | $+3811.81 | 86.6% | 108.7% | 46.6% | 4.03x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 17 | 380 | 92.7% | $+4069.17 | 83.9% | 108.4% | 50.5% | 3.74x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
