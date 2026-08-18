# Strategy Lab Report

**Generated:** 2026-08-18 23:24 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (409)
- Valid strategies found: 69
- Candidate layers: 24 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 59 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 2 observation-only

## Selected Core Strategy

- Wallets: 24
- Aggregate closed copy trades: 527
- Aggregate copy ROI: 95.9%
- Aggregate copy PnL: $+5992.21
- Aggregate copy win rate: 84.1%
- Median wallet CopyROI: 108.7%
- Worst included wallet CopyROI: 44.4%
- Open copy cost / closed copy capital: 3.13x
- Params: tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.27 | 415 | 100% | 101.7% | $+1138.77 | 93 | 79.6% | 42.1% |
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0x586a...be99` | A | 24.19 | 34 | 79% | 162.6% | $+146.36 | 8 | 75.0% | 33.9% |
| `0x5be9...ef2b` | A | 14.61 | 57 | 63% | 110.4% | $+132.44 | 12 | 100.0% | 30.8% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.64 | 18 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x2fe2...9645` | A | 19.06 | 40 | 37% | 44.4% | $+39.94 | 9 | 88.9% | 48.2% |
| `0x18c2...529a` | B | 27.48 | 902 | 78% | 69.5% | $+528.48 | 57 | 89.5% | 45.6% |
| `0x44c4...09cb` | B | 26.80 | 126 | 61% | 109.7% | $+482.80 | 41 | 78.0% | 55.3% |
| `0x0d66...ff1d` | B | 25.42 | 82 | 99% | 110.0% | $+407.18 | 21 | 85.7% | 34.0% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x272d...bc2f` | B | 28.68 | 521 | 51% | 60.6% | $+181.86 | 28 | 78.6% | 67.8% |
| `0x95c1...94d7` | B | 27.46 | 56 | 50% | 44.4% | $+159.96 | 34 | 94.1% | 20.0% |
| `0x092b...614e` | B | 25.68 | 134 | 96% | 93.2% | $+158.52 | 16 | 81.2% | 61.6% |
| `0x819d...6e9c` | B | 28.04 | 26 | 28% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |
| `0x0931...e78e` | B | 25.83 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | 26.9% |
| `0x4dee...8ad7` | B | 27.07 | 29 | 71% | 51.3% | $+82.10 | 9 | 66.7% | 24.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 470
- Aggregate copy ROI: 110.7%
- Aggregate copy PnL: $+6267.38
- Aggregate copy win rate: 84.5%
- Worst included CopyROI: 48.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9c03...522d` | B | 100.0 | 32.8 | 1208.8 | 180 | 65% | 244.1% | $+1439.94 | 57 | 96.5% | $789 | existing,holder |
| `0x3c2b...825e` | B | 100.0 | 33.8 | 843.5 | 508 | 96% | 156.6% | $+1080.70 | 61 | 95.1% | $268 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 31.9 | 736.6 | 37 | 23% | 149.2% | $+656.26 | 34 | 73.5% | $465 | existing |
| `0x54fc...76eb` | B | 100.0 | 31.3 | 617.6 | 46 | 87% | 211.4% | $+232.59 | 6 | 66.7% | $639 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x760f...326a` | B | 100.0 | 30.9 | 524.0 | 15 | 14% | 112.2% | $+246.72 | 18 | 83.3% | $5032 | existing |
| `0xa35c...a113` | B | 100.0 | 33.4 | 487.6 | 153 | 100% | 75.7% | $+333.23 | 40 | 95.0% | $703 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.8 | 431.3 | 48 | 31% | 57.0% | $+267.89 | 45 | 80.0% | $1342 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0xa931...1387` | B | 100.0 | 29.3 | 427.1 | 10 | 77% | 128.8% | $+51.51 | 4 | 100.0% | $2630 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 415.8 | 192 | 98% | 48.1% | $+418.37 | 59 | 79.7% | $1108 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 395.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing,holder |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x8fd6...c2b8` | B | 100.0 | 6.7 | 366.8 | 16 | 32% | 94.6% | $+37.84 | 4 | 100.0% | $633 | existing |
| `0x6d15...04c2` | B | 100.0 | 26.5 | 359.5 | 43 | 83% | 58.1% | $+87.15 | 12 | 58.3% | $520 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.1 | 919.9 | 361 | 268 | 35% | 168.1% | $+1193.66 | 68 | 97.1% | 171.0% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 31.1 | 709.8 | 1370 | 399 | 100% | 135.2% | $+1513.92 | 41 | 100.0% | 135.2% | existing,holder | burst_trading,open_copy_exposure |
| `0x07b1...6dfc` | B | 100.0 | 35.0 | 666.7 | 651 | 185 | 72% | 188.8% | $+1038.10 | 11 | 100.0% | 171.7% | existing | burst_trading,open_copy_exposure |
| `0x6e32...ab65` | B | 100.0 | 34.7 | 605.7 | 822 | 297 | 66% | 120.2% | $+733.23 | 31 | 100.0% | 129.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 333.3 | 422 | 232 | 75% | 39.0% | $+225.89 | 50 | 68.0% | 34.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.1 | 329.8 | 319 | 126 | 29% | 50.5% | $+242.66 | 21 | 38.1% | 67.8% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 43.7 | 200.6 | 1257 | 90% | 0.0% | 0 | 414 | $574 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6ac5...4b6e` | C | 100.0 | 39.7 | 250.2 | 1128 | 327 | 81% | 178.8% | 11 | 164.2% | 18 | $448 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.1 | 217.7 | 708 | 339 | 91% | 36.2% | 24 | 31.5% | 28 | $143 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.4 | 212.7 | 591 | 360 | 89% | 28.8% | 82 | 27.2% | 89 | $225 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.2 | 210.9 | 296 | 257 | 72% | 17.5% | 73 | 23.5% | 91 | $1114 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 199.9 | 109 | 105 | 60% | 55.6% | 16 | 46.9% | 25 | $3028 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | C | 100.0 | 41.6 | 199.2 | 1191 | 280 | 88% | 147.8% | 24 | 149.1% | 30 | $313 | existing | burst_trading,opposite_side_same_market |
| `0x565c...878d` | C | 100.0 | 39.5 | 194.1 | 256 | 243 | 60% | 147.4% | 56 | 132.7% | 102 | $969 | existing,holder | open_copy_exposure,opposite_side_same_market |

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
| `0x44c4...09cb` | B | 100.0 | 26.8 | 1h edge -15.87pp over 1 samples | 109.7% | 41 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.5 | 1h edge -52.95pp over 1 samples | 69.5% | 57 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 48.1% | 59 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x8002...8b6c` | reject-bot | - | BOT | 79.1 | $8003 | $8003 | 1 | -23.6% | 25 | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.3 | 738.4 | 415 | 295 | 100% | 101.7% | 93 | 101.7% | 93 | existing,sports_tape | opposite_side_same_market |
| `0x8002...8b6c` | reject-bot | - | BOT | 71.9 | 79.1 | -104.5 | 377 | 174 | 74% | -10.7% | 17 | -23.6% | 25 | existing,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 24 | 527 | 95.9% | $+5992.21 | 84.1% | 108.7% | 44.4% | 3.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 2 | 25 | 540 | 93.6% | $+5936.37 | 84.6% | 108.4% | 25.2% | 3.00x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 25 | 560 | 93.0% | $+6184.89 | 82.9% | 108.4% | 44.4% | 2.96x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 23 | 518 | 96.6% | $+5952.27 | 84.0% | 108.9% | 44.4% | 3.01x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 5 | 24 | 551 | 93.7% | $+6144.95 | 82.8% | 108.7% | 44.4% | 2.85x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 22 | 498 | 98.4% | $+5785.29 | 85.1% | 109.3% | 44.4% | 3.21x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 22 | 509 | 97.8% | $+5870.17 | 84.3% | 109.3% | 44.4% | 3.08x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 22 | 502 | 97.5% | $+5808.19 | 85.1% | 109.3% | 37.0% | 3.01x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 20 | 455 | 102.9% | $+5585.39 | 84.4% | 109.8% | 60.6% | 3.26x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 23 | 542 | 94.7% | $+6062.85 | 83.0% | 108.9% | 44.4% | 2.91x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 21 | 489 | 99.2% | $+5745.35 | 85.1% | 109.7% | 44.4% | 3.09x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 15 | 318 | 116.1% | $+4424.03 | 83.6% | 110.4% | 101.7% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
