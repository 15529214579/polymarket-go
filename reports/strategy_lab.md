# Strategy Lab Report

**Generated:** 2026-08-28 23:34 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (426)
- Valid strategies found: 40
- Candidate layers: 13 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 50 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 13
- Aggregate closed copy trades: 236
- Aggregate copy ROI: 119.4%
- Aggregate copy PnL: $+3115.13
- Aggregate copy win rate: 87.7%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 3.10x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x5be9...ef2b` | A | 23.64 | 80 | 64% | 184.8% | $+517.37 | 28 | 100.0% | 51.4% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.77 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x44c4...09cb` | B | 26.12 | 115 | 62% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.22 | 151 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.20 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 503
- Aggregate copy ROI: 110.9%
- Aggregate copy PnL: $+7100.52
- Aggregate copy win rate: 85.9%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x93b3...0bac` | A | 97.3 | 22.1 | 1284.4 | 14 | 67% | 513.9% | $+359.74 | 7 | 100.0% | $415 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.2 | 889.0 | 604 | 88% | 146.9% | $+1909.46 | 111 | 95.5% | $296 | existing |
| `0x4362...9e47` | B | 100.0 | 35.0 | 844.2 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1325 | existing |
| `0x5461...92be` | A | 98.5 | 21.7 | 715.5 | 37 | 55% | 163.7% | $+507.61 | 17 | 100.0% | $1181 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.8 | 587.4 | 170 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $730 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 550.9 | 15 | 14% | 119.0% | $+273.80 | 19 | 84.2% | $5209 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.4 | 430.4 | 49 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1312 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0xaabc...a190` | A | 100.0 | 4.4 | 417.8 | 22 | 21% | 101.1% | $+50.53 | 5 | 100.0% | $434 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.5 | 72 | 59% | 50.5% | $+207.07 | 34 | 64.7% | $543 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x2ba6...644c` | B | 100.0 | 28.7 | 389.4 | 25 | 74% | 62.8% | $+94.19 | 4 | 100.0% | $123509 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_volume_7d |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1054.3 | 379 | 290 | 36% | 183.6% | $+1835.46 | 94 | 98.9% | 183.1% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.9 | 740.0 | 218 | 79 | 46% | 194.5% | $+330.69 | 16 | 93.8% | 127.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.1 | 623.4 | 1353 | 386 | 100% | 128.5% | $+822.58 | 27 | 100.0% | 128.5% | existing,holder | burst_trading,open_copy_exposure |
| `0x2829...12a0` | B | 93.6 | 33.4 | 562.2 | 185 | 95 | 89% | 124.0% | $+558.14 | 18 | 100.0% | 124.0% | existing,holder | burst_trading,open_copy_exposure |
| `0x153a...9d6d` | B | 100.0 | 34.3 | 547.8 | 165 | 59 | 79% | 179.4% | $+412.60 | 6 | 100.0% | 179.4% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x4cb0...57e0` | B | 100.0 | 33.6 | 545.8 | 652 | 294 | 93% | 109.7% | $+329.00 | 26 | 100.0% | 108.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.3 | 542.9 | 1390 | 337 | 100% | 98.9% | $+1226.27 | 29 | 100.0% | 98.9% | existing | burst_trading,open_copy_exposure |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.8 | 401.8 | 125 | 70 | 30% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x168a...0d24` | B | 100.0 | 25.7 | 258.2 | 35 | 36% | 79.7% | 2 | 64 | $28644 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x65c1...2988` | C | 100.0 | 39.6 | 222.6 | 1313 | 92% | 0.0% | 0 | 345 | $521 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 25.0 | 298.9 | 1056 | 770 | 96% | 10.2% | 10 | 9.4% | 11 | $513 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.2 | 220.1 | 290 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $621 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.8 | 616 | 382 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | B | 100.0 | 34.8 | 212.9 | 317 | 275 | 70% | 18.1% | 78 | 25.6% | 100 | $1064 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 208.3 | 664 | 302 | 64% | 61.7% | 22 | 51.9% | 29 | $175 | existing | open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 100.0 | 42.1 | 206.2 | 1249 | 261 | 96% | 122.7% | 36 | 121.9% | 38 | $438 | existing,holder | burst_trading |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.1 | 907 | 392 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
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
| `0x44c4...09cb` | B | 100.0 | 26.1 | 1h edge -15.87pp over 1 samples | 135.3% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3b62...5cc8` | reject | - | D | 100.0 | 57.6 | 125.5 | 596 | 433 | 53% | 29.8% | 9 | 27.8% | 11 | holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 13 | 236 | 119.4% | $+3115.13 | 87.7% | 115.4% | 65.3% | 3.10x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 10 | 182 | 129.9% | $+2637.28 | 87.9% | 121.9% | 104.0% | 3.23x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 12 | 228 | 120.2% | $+3005.62 | 88.2% | 117.6% | 65.3% | 3.12x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 11 | 218 | 120.3% | $+2801.86 | 89.0% | 115.4% | 65.3% | 2.05x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=80 closedROI>=0 smart>=70 |
| 5 | 14 | 265 | 112.0% | $+3281.24 | 86.0% | 112.2% | 51.9% | 4.06x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 12 | 227 | 117.4% | $+2841.80 | 89.0% | 112.2% | 44.4% | 2.34x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 7 | 17 | 303 | 104.1% | $+3528.10 | 84.2% | 108.4% | 44.4% | 3.96x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 13 | 257 | 112.5% | $+3171.73 | 86.4% | 115.4% | 51.9% | 4.12x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 15 | 274 | 110.0% | $+3321.18 | 86.1% | 108.9% | 44.4% | 4.23x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 16 | 294 | 105.7% | $+3488.16 | 84.0% | 108.7% | 51.3% | 3.81x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 15 | 285 | 108.5% | $+3406.06 | 84.6% | 108.9% | 51.9% | 3.98x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 16 | 291 | 105.5% | $+3363.98 | 85.6% | 108.7% | 25.2% | 4.04x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
