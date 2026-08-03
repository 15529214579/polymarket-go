# Strategy Lab Report

**Generated:** 2026-08-03 23:01 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (315)
- Valid strategies found: 159
- Candidate layers: 11 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 45 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 229
- Aggregate copy ROI: 127.8%
- Aggregate copy PnL: $+3450.09
- Aggregate copy win rate: 81.2%
- Median wallet CopyROI: 135.2%
- Worst included wallet CopyROI: 108.4%
- Open copy cost / closed copy capital: 2.28x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.89 | 211 | 100% | 113.9% | $+797.15 | 53 | 83.0% | 36.0% |
| `0xb2ed...4418` | A | 23.14 | 125 | 38% | 134.4% | $+591.26 | 36 | 66.7% | 58.6% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 26.84 | 142 | 64% | 108.6% | $+521.13 | 44 | 77.3% | 47.6% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.9% |
| `0x0ec9...1e0c` | B | 26.70 | 7 | 13% | 183.6% | $+146.85 | 8 | 87.5% | 11.1% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 694
- Aggregate copy ROI: 74.7%
- Aggregate copy PnL: $+5788.28
- Aggregate copy win rate: 84.6%
- Worst included CopyROI: 41.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 29 | 48% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xcc6e...fa6f` | B | 100.0 | 25.6 | 613.9 | 254 | 82% | 99.6% | $+577.47 | 51 | 88.2% | $2191 | existing,sports_tape |
| `0x17fe...b0ca` | B | 100.0 | 31.0 | 568.0 | 175 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $674 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.7 | 559.9 | 27 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $754 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 525.7 | 899 | 93% | 59.8% | $+1220.55 | 185 | 89.2% | $296 | existing |
| `0x092b...614e` | B | 100.0 | 25.9 | 493.9 | 147 | 97% | 97.6% | $+185.38 | 18 | 83.3% | $1081 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |
| `0x18c2...529a` | B | 100.0 | 28.1 | 486.3 | 837 | 74% | 70.3% | $+520.04 | 57 | 86.0% | $268 | existing,holder |
| `0x819d...6e9c` | B | 100.0 | 31.6 | 485.6 | 21 | 31% | 134.1% | $+107.24 | 8 | 100.0% | $420 | existing |
| `0x84cd...7565` | A | 100.0 | 23.5 | 475.2 | 125 | 80% | 93.9% | $+178.48 | 16 | 75.0% | $443 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 24.1 | 451.5 | 57 | 34% | 59.8% | $+286.98 | 46 | 80.4% | $1277 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.3 | 417.3 | 28 | 38% | 81.2% | $+105.62 | 10 | 70.0% | $2180 | existing |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 412.5 | 73 | 95% | 58.4% | $+175.23 | 20 | 85.0% | $3677 | existing |
| `0x2a99...51bb` | A | 100.0 | 24.0 | 408.2 | 698 | 99% | 41.1% | $+522.65 | 124 | 76.6% | $311 | existing,holder |
| `0x7673...fa40` | B | 100.0 | 28.7 | 405.0 | 93 | 62% | 76.6% | $+122.51 | 14 | 78.6% | $545 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.0 | 620.9 | 385 | 274 | 40% | 104.9% | $+618.86 | 56 | 96.4% | 127.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xcefd...d6aa` | B | 100.0 | 29.5 | 597.6 | 95 | 65 | 70% | 177.3% | $+248.17 | 9 | 100.0% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 33.6 | 457.2 | 327 | 196 | 93% | 130.6% | $+169.75 | 8 | 75.0% | 121.6% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.6 | 379.9 | 125 | 71 | 68% | 71.2% | $+106.79 | 15 | 73.3% | 55.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 21.4 | 375.0 | 72 | 58 | 67% | 54.0% | $+156.56 | 22 | 63.6% | 45.6% | existing | opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 371.9 | 61 | 55 | 38% | 62.3% | $+118.35 | 18 | 83.3% | 25.4% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 359.0 | 433 | 225 | 42% | 64.5% | $+116.01 | 16 | 81.2% | 59.7% | existing,holder | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4b40...2f6a` | B | 100.0 | 32.0 | 209.8 | 26 | 26% | 0.0% | 0 | 74 | $1615 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 44.0 | 173.8 | 0 | 0% | 123.1% | 1 | 370 | $6799 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure,opposite_side_same_market |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbfab...f067` | C | 100.0 | 35.8 | 276.2 | 1328 | 778 | 100% | 11.8% | 217 | 11.8% | 217 | $321 | existing | opposite_side_same_market |
| `0x2929...1dd0` | C | 100.0 | 43.8 | 241.3 | 569 | 446 | 99% | 109.2% | 24 | 107.7% | 25 | $1107 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xd3d3...bb8e` | C | 100.0 | 39.3 | 238.9 | 832 | 606 | 72% | 15.3% | 176 | 13.6% | 249 | $655 | existing | opposite_side_same_market |
| `0xde24...4ded` | C | 100.0 | 37.4 | 237.0 | 386 | 351 | 100% | 91.7% | 45 | 91.7% | 45 | $1720 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd3bc...5bd2` | C | 100.0 | 42.2 | 232.2 | 437 | 403 | 100% | 47.4% | 121 | 47.4% | 121 | $1278 | existing,holder | opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.7 | 224.6 | 471 | 218 | 98% | 20.2% | 74 | 20.2% | 74 | $199 | existing | opposite_side_same_market |
| `0x4bba...cf14` | B | 100.0 | 26.4 | 220.4 | 768 | 213 | 100% | 20.0% | 6 | 20.0% | 6 | $103 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 28.8 | 217.1 | 129 | 122 | 98% | 41.4% | 35 | 43.4% | 37 | $1020 | existing | opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0xcbb5...5ed9` | B | 100.0 | 26.7 | 208.9 | 248 | 197 | 62% | 11.5% | 12 | 15.6% | 20 | $542 | existing,holder | open_copy_exposure,opposite_side_same_market |

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

- Wallets: 4
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 23.9 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 564.4 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.8 | $1000 | $2000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 404.4 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | B | 25.6 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 324.1 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 15m edge -2.63pp over 2 samples | 25.4% | 29 | existing | opposite_side_same_market |
| `0xde24...4ded` | C | 100.0 | 37.4 | 1h edge -11.42pp over 4 samples | 91.7% | 45 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.8 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.1 | 1h edge -19.00pp over 1 samples | 134.4% | 36 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.1 | 1h edge -52.95pp over 1 samples | 70.3% | 57 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 28.8 | 1h edge -72.59pp over 1 samples | 43.4% | 37 | existing | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 4
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.9 | 716.2 | 211 | 167 | 100% | 113.9% | 53 | 113.9% | 53 | existing,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.8 | 558.5 | 569 | 446 | 99% | 109.2% | 24 | 107.7% | 25 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 25.6 | 547.6 | 254 | 212 | 82% | 84.0% | 41 | 99.6% | 51 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xe59f...0d08` | reject | - | D | 100.0 | 59.0 | 119.3 | 1473 | 725 | 100% | 8.4% | 132 | 8.4% | 132 | existing,holder,leaderboard_profit_30d,leaderboard_volume_30d | bot_like_flow,burst_trading,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 229 | 127.8% | $+3450.09 | 81.2% | 135.2% | 108.4% | 2.28x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 21 | 482 | 101.9% | $+5910.93 | 83.0% | 108.4% | 62.9% | 2.55x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 471 | 102.9% | $+5803.50 | 83.2% | 108.5% | 62.9% | 2.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 4 | 10 | 193 | 126.5% | $+2858.83 | 83.9% | 135.4% | 108.4% | 1.98x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 19 | 435 | 100.2% | $+5212.24 | 84.6% | 108.4% | 62.9% | 2.47x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 6 | 20 | 446 | 99.2% | $+5319.67 | 84.3% | 104.0% | 62.9% | 2.44x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 18 | 441 | 101.3% | $+5329.86 | 83.0% | 104.0% | 62.9% | 2.67x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 18 | 454 | 99.6% | $+5420.63 | 83.0% | 104.0% | 62.9% | 2.60x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 9 | 19 | 465 | 98.7% | $+5528.06 | 82.8% | 99.6% | 62.9% | 2.57x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 26 | 752 | 83.3% | $+7375.82 | 80.9% | 95.8% | 41.1% | 2.58x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 27 | 761 | 82.8% | $+7457.92 | 80.7% | 93.9% | 41.1% | 2.54x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 27 | 783 | 81.8% | $+7548.98 | 80.1% | 93.9% | 41.1% | 2.49x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=5 smart>=70 |
