# Strategy Lab Report

**Generated:** 2026-08-17 23:45 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (409)
- Valid strategies found: 75
- Candidate layers: 20 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 56 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 2 observation-only

## Selected Core Strategy

- Wallets: 20
- Aggregate closed copy trades: 459
- Aggregate copy ROI: 101.5%
- Aggregate copy PnL: $+5552.36
- Aggregate copy win rate: 84.3%
- Median wallet CopyROI: 108.7%
- Worst included wallet CopyROI: 62.9%
- Open copy cost / closed copy capital: 3.27x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.27 | 415 | 100% | 101.7% | $+1138.77 | 93 | 79.6% | 42.1% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0x5be9...ef2b` | A | 14.59 | 55 | 63% | 110.4% | $+132.44 | 12 | 100.0% | 30.1% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.88 | 19 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x18c2...529a` | B | 27.82 | 883 | 76% | 71.7% | $+559.45 | 59 | 89.8% | 45.2% |
| `0x44c4...09cb` | B | 26.98 | 133 | 63% | 107.2% | $+503.66 | 43 | 76.7% | 55.3% |
| `0x0d66...ff1d` | B | 25.26 | 77 | 99% | 109.4% | $+382.85 | 20 | 85.0% | 33.2% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x272d...bc2f` | B | 28.77 | 520 | 51% | 68.3% | $+204.82 | 28 | 82.1% | 67.8% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x092b...614e` | B | 25.68 | 134 | 96% | 93.2% | $+158.52 | 16 | 81.2% | 61.6% |
| `0x819d...6e9c` | B | 28.17 | 26 | 28% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |
| `0x7673...fa40` | B | 28.23 | 80 | 60% | 62.9% | $+62.87 | 9 | 66.7% | 53.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 436
- Aggregate copy ROI: 85.3%
- Aggregate copy PnL: $+4450.16
- Aggregate copy win rate: 83.3%
- Worst included CopyROI: 43.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3c2b...825e` | B | 100.0 | 33.7 | 840.5 | 504 | 96% | 156.6% | $+1080.70 | 61 | 95.1% | $265 | existing |
| `0x54fc...76eb` | B | 100.0 | 31.3 | 617.6 | 46 | 87% | 211.4% | $+232.59 | 6 | 66.7% | $639 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 524.0 | 15 | 14% | 112.2% | $+246.72 | 18 | 83.3% | $5032 | existing |
| `0xa35c...a113` | B | 100.0 | 33.2 | 487.5 | 155 | 100% | 75.7% | $+333.23 | 40 | 95.0% | $702 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.8 | 431.3 | 48 | 31% | 57.0% | $+267.89 | 45 | 80.0% | $1342 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0xa931...1387` | B | 100.0 | 29.3 | 427.1 | 10 | 77% | 128.8% | $+51.51 | 4 | 100.0% | $2630 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 415.8 | 192 | 98% | 48.1% | $+418.37 | 59 | 79.7% | $1108 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x95c1...94d7` | B | 100.0 | 27.6 | 382.1 | 55 | 50% | 43.8% | $+153.18 | 33 | 93.9% | $918 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x0931...e78e` | B | 100.0 | 25.8 | 373.9 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | $298 | existing |
| `0x8fd6...c2b8` | B | 100.0 | 6.7 | 366.8 | 16 | 32% | 94.6% | $+37.84 | 4 | 100.0% | $633 | existing |
| `0x6d15...04c2` | B | 100.0 | 26.5 | 359.5 | 43 | 83% | 58.1% | $+87.15 | 12 | 58.3% | $520 | existing |
| `0x1abb...0901` | B | 100.0 | 26.7 | 350.1 | 30 | 71% | 71.3% | $+49.95 | 7 | 71.4% | $1293 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.3 | 809.3 | 361 | 269 | 35% | 142.4% | $+1010.96 | 68 | 97.1% | 161.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.2 | 723.5 | 1361 | 386 | 100% | 141.1% | $+1467.78 | 39 | 100.0% | 141.1% | existing,holder | burst_trading,open_copy_exposure |
| `0x6e32...ab65` | B | 100.0 | 34.9 | 605.6 | 810 | 291 | 65% | 120.2% | $+733.23 | 31 | 100.0% | 122.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 333.3 | 422 | 232 | 75% | 39.0% | $+225.89 | 50 | 68.0% | 34.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.1 | 329.8 | 319 | 126 | 29% | 50.5% | $+242.66 | 21 | 38.1% | 67.8% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 317.5 | 109 | 105 | 60% | 55.6% | $+94.59 | 16 | 43.8% | 46.9% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 43.8 | 200.8 | 1260 | 90% | 0.0% | 0 | 420 | $569 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 36.5 | 192.6 | 0 | 0% | 123.1% | 1 | 327 | $9143 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6ac5...4b6e` | C | 100.0 | 38.4 | 231.2 | 1131 | 305 | 81% | 178.8% | 11 | 164.2% | 18 | $422 | existing,holder,leaderboard_profit_30d,leaderboard_volume_7d | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.1 | 216.9 | 696 | 333 | 91% | 42.7% | 22 | 36.7% | 26 | $143 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.4 | 212.7 | 591 | 360 | 89% | 28.8% | 82 | 27.2% | 89 | $225 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 210.8 | 293 | 255 | 72% | 17.5% | 73 | 23.5% | 91 | $1126 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | C | 100.0 | 41.6 | 199.2 | 1191 | 280 | 88% | 147.8% | 24 | 149.1% | 30 | $313 | existing | burst_trading,opposite_side_same_market |
| `0x07b1...6dfc` | C | 100.0 | 35.6 | 192.5 | 592 | 175 | 71% | 188.8% | 11 | 171.7% | 13 | $389 | existing,holder | burst_trading,open_copy_exposure |
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

- Wallets: 2
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 23.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 562.1 | opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 6
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.0 | 1h edge -15.87pp over 1 samples | 107.2% | 43 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.8 | 1h edge -52.95pp over 1 samples | 71.7% | 59 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 48.1% | 59 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x8002...8b6c` | reject-bot | - | BOT | 79.4 | $8003 | $8003 | 1 | -23.6% | 25 | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.3 | 738.4 | 415 | 295 | 100% | 101.7% | 93 | 101.7% | 93 | existing,sports_tape | opposite_side_same_market |
| `0x8002...8b6c` | reject-bot | - | BOT | 71.5 | 79.4 | -105.8 | 371 | 169 | 74% | -10.7% | 17 | -23.6% | 25 | existing,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 20 | 459 | 101.5% | $+5552.36 | 84.3% | 108.7% | 62.9% | 3.27x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 24 | 532 | 93.8% | $+6010.92 | 83.8% | 105.6% | 43.8% | 3.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 3 | 25 | 541 | 93.1% | $+6050.86 | 83.9% | 104.0% | 43.8% | 3.24x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 4 | 21 | 503 | 97.3% | $+5767.49 | 84.5% | 108.4% | 43.8% | 3.06x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 5 | 22 | 503 | 95.9% | $+5781.07 | 85.3% | 107.8% | 43.8% | 3.29x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 19 | 450 | 102.2% | $+5489.49 | 84.7% | 108.9% | 65.3% | 3.23x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 22 | 507 | 95.0% | $+5803.97 | 85.2% | 107.8% | 37.0% | 3.09x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 22 | 536 | 94.2% | $+5960.17 | 83.2% | 107.8% | 43.8% | 2.89x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 26 | 574 | 90.5% | $+6243.54 | 82.8% | 102.8% | 43.8% | 3.07x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 20 | 483 | 98.6% | $+5642.67 | 85.3% | 108.7% | 43.8% | 3.07x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 21 | 494 | 96.7% | $+5741.13 | 85.2% | 108.4% | 43.8% | 3.17x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 25 | 565 | 91.1% | $+6203.60 | 82.7% | 104.0% | 43.8% | 2.97x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
