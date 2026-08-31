# Strategy Lab Report

**Generated:** 2026-08-31 23:45 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (429)
- Valid strategies found: 78
- Candidate layers: 13 core + 20 watch + 10 sports + 4 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 52 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 13
- Aggregate closed copy trades: 275
- Aggregate copy ROI: 121.3%
- Aggregate copy PnL: $+3676.46
- Aggregate copy win rate: 89.8%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 2.90x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x5be9...ef2b` | A | 23.49 | 83 | 62% | 192.0% | $+672.13 | 35 | 100.0% | 57.1% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.71 | 78 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.51 | 120 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.70 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x4607...d542` | B | 26.17 | 104 | 55% | 124.6% | $+560.79 | 42 | 100.0% | 43.6% |
| `0x44c4...09cb` | B | 25.94 | 114 | 61% | 135.3% | $+459.88 | 33 | 81.8% | 55.5% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.65 | 163 | 97% | 96.6% | $+202.79 | 20 | 85.0% | 61.1% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.43 | 26 | 27% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 671
- Aggregate copy ROI: 118.2%
- Aggregate copy PnL: $+9622.63
- Aggregate copy win rate: 88.5%
- Worst included CopyROI: 46.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3c2b...825e` | B | 100.0 | 34.9 | 1005.2 | 568 | 73% | 165.8% | $+2387.83 | 126 | 97.6% | $268 | existing,holder |
| `0x8ce7...3150` | B | 100.0 | 34.5 | 955.0 | 197 | 27% | 149.1% | $+2519.90 | 156 | 97.4% | $591 | existing,holder |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 942.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing,leaderboard_profit_7d |
| `0x5461...92be` | A | 98.5 | 21.7 | 715.5 | 37 | 55% | 163.7% | $+507.61 | 17 | 100.0% | $1181 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x760f...326a` | B | 100.0 | 30.0 | 531.3 | 13 | 13% | 109.7% | $+263.21 | 21 | 85.7% | $5470 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.3 | 450.6 | 72 | 55% | 63.5% | $+273.08 | 35 | 65.7% | $551 | existing |
| `0x075d...58c1` | A | 100.0 | 22.6 | 450.6 | 6 | 11% | 72.7% | $+196.26 | 27 | 85.2% | $3497 | existing |
| `0x789f...20ec` | B | 99.3 | 27.2 | 444.2 | 52 | 98% | 71.6% | $+171.80 | 18 | 94.4% | $1103 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.3 | 430.0 | 50 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1302 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 421.8 | 235 | 97% | 46.3% | $+509.54 | 79 | 77.2% | $1172 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x8bd3...f69a` | A | 100.0 | 24.7 | 385.7 | 20 | 31% | 67.1% | $+107.43 | 11 | 72.7% | $5404 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.4 | 1141.7 | 380 | 290 | 36% | 201.1% | $+2071.36 | 97 | 99.0% | 191.2% | existing | opposite_side_same_market |
| `0x7e33...73e7` | B | 100.0 | 34.4 | 722.7 | 217 | 76 | 49% | 198.1% | $+297.12 | 14 | 92.9% | 129.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x4cb0...57e0` | B | 100.0 | 33.3 | 624.3 | 686 | 310 | 93% | 125.8% | $+503.34 | 31 | 100.0% | 122.8% | existing,holder | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 32.4 | 568.6 | 1347 | 400 | 100% | 104.4% | $+1075.05 | 34 | 100.0% | 104.4% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.9 | 549.9 | 1402 | 346 | 100% | 88.0% | $+1653.45 | 44 | 100.0% | 88.0% | existing,holder | burst_trading |
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 515.2 | 1153 | 1102 | 100% | 68.9% | $+1391.99 | 106 | 95.3% | 68.9% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0571...c5cc` | A | 100.0 | 21.8 | 408.7 | 141 | 58 | 99% | 101.0% | $+121.19 | 10 | 20.0% | 101.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.7 | 403.4 | 787 | 334 | 74% | 65.0% | $+201.38 | 29 | 86.2% | 56.3% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 32.5 | 402.7 | 122 | 69 | 31% | 75.8% | $+144.06 | 17 | 70.6% | 87.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 4
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x565c...878d` | C | 100.0 | 39.7 | 245.3 | 475 | 62% | 183.3% | 178 | 718 | $1281 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |
| `0x983e...6748` | C | 100.0 | 44.2 | 223.4 | 301 | 68% | 142.3% | 1 | 206 | $16771 | existing,leaderboard_profit_30d,leaderboard_profit_all,leaderboard_volume_7d | burst_trading,open_copy_exposure |
| `0x168a...0d24` | C | 100.0 | 39.6 | 220.6 | 69 | 33% | 79.7% | 2 | 155 | $23012 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x684b...8409` | C | 100.0 | 37.8 | 220.5 | 109 | 8% | 0.0% | 0 | 315 | $543 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x67ba...246c` | C | 100.0 | 38.7 | 249.3 | 1049 | 569 | 98% | 58.0% | 158 | 58.7% | 163 | $614 | existing | opposite_side_same_market |
| `0x0467...fb2e` | C | 100.0 | 42.7 | 228.4 | 1290 | 486 | 100% | 12.4% | 8 | 12.4% | 8 | $297 | existing | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.1 | 327 | 185 | 81% | 10.8% | 52 | 9.7% | 62 | $561 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 216.3 | 619 | 385 | 89% | 29.1% | 86 | 27.6% | 93 | $237 | existing | open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 100.0 | 40.0 | 216.1 | 1378 | 312 | 97% | 139.2% | 39 | 139.0% | 40 | $421 | existing,holder | burst_trading |
| `0x18c2...529a` | C | 100.0 | 43.0 | 205.1 | 907 | 392 | 80% | 68.5% | 54 | 71.8% | 63 | $288 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
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

- Wallets: 1
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xae2b...90ea` | C | 100.0 | 44.5 | 623.4 | 397 | 362 | 119.8% | 92 | 139.3% | 143 | $1863 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |

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
| `0xf3ce...a57a` | B | 100.0 | 34.3 | 1h edge -47.28pp over 1 samples | 68.9% | 106 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | C | 100.0 | 43.0 | 1h edge -52.95pp over 1 samples | 71.8% | 63 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 46.3% | 79 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xae2b...90ea` | pushed | - | C | 100.0 | 44.5 | 864.0 | 397 | 362 | 40% | 119.8% | 92 | 139.3% | 143 | existing,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x19ae...395c` | reject-bot | - | BOT | 100.0 | 64.3 | 469.5 | 1279 | 398 | 90% | 140.5% | 8 | 138.2% | 11 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 13 | 275 | 121.3% | $+3676.46 | 89.8% | 115.4% | 65.3% | 2.90x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 2 | 14 | 285 | 119.6% | $+3756.34 | 90.2% | 112.2% | 65.3% | 2.89x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 13 | 277 | 120.4% | $+3646.83 | 90.6% | 115.4% | 65.3% | 2.89x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 10 | 219 | 130.9% | $+3181.17 | 90.4% | 121.9% | 104.0% | 2.96x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 14 | 302 | 117.4% | $+3872.72 | 89.4% | 112.2% | 65.3% | 2.71x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 6 | 12 | 267 | 122.2% | $+3566.95 | 90.3% | 117.6% | 65.3% | 2.91x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 15 | 312 | 115.9% | $+3952.60 | 89.7% | 108.9% | 65.3% | 2.70x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 13 | 294 | 118.0% | $+3763.21 | 89.8% | 115.4% | 65.3% | 2.71x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 14 | 304 | 116.5% | $+3843.09 | 90.1% | 112.2% | 65.3% | 2.70x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 16 | 347 | 110.0% | $+4225.68 | 87.3% | 108.7% | 63.5% | 2.44x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 17 | 357 | 108.3% | $+4212.00 | 88.5% | 108.4% | 44.4% | 3.47x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 16 | 348 | 109.8% | $+4172.06 | 88.5% | 108.7% | 56.3% | 3.34x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
