# Strategy Lab Report

**Generated:** 2026-08-16 23:06 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (407)
- Valid strategies found: 77
- Candidate layers: 20 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 55 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 2 observation-only

## Selected Core Strategy

- Wallets: 20
- Aggregate closed copy trades: 457
- Aggregate copy ROI: 99.3%
- Aggregate copy PnL: $+5333.00
- Aggregate copy win rate: 84.7%
- Median wallet CopyROI: 107.8%
- Worst included wallet CopyROI: 62.3%
- Open copy cost / closed copy capital: 3.54x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.27 | 415 | 100% | 101.7% | $+1138.77 | 93 | 79.6% | 42.1% |
| `0x162d...8944` | A | 24.93 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0xb99c...7e96` | A | 21.70 | 131 | 91% | 66.7% | $+200.00 | 19 | 89.5% | 26.2% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 49.1% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x7fcf...80ac` | A | 10.88 | 19 | 33% | 99.5% | $+109.51 | 8 | 75.0% | 37.3% |
| `0x5be9...ef2b` | A | 14.85 | 53 | 63% | 108.9% | $+108.90 | 10 | 100.0% | 27.0% |
| `0x18c2...529a` | B | 27.89 | 869 | 75% | 73.0% | $+569.27 | 60 | 88.3% | 46.0% |
| `0x44c4...09cb` | B | 26.98 | 133 | 63% | 107.2% | $+503.66 | 43 | 76.7% | 55.3% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.1% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x272d...bc2f` | B | 28.76 | 521 | 51% | 68.3% | $+204.82 | 28 | 82.1% | 67.8% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x2b17...d36d` | B | 25.41 | 416 | 73% | 134.6% | $+174.98 | 8 | 62.5% | 52.3% |
| `0x092b...614e` | B | 25.68 | 134 | 96% | 93.2% | $+158.52 | 16 | 81.2% | 61.6% |
| `0x819d...6e9c` | B | 28.31 | 26 | 29% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |
| `0x7673...fa40` | B | 28.44 | 81 | 60% | 62.3% | $+68.57 | 10 | 70.0% | 53.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 415
- Aggregate copy ROI: 74.3%
- Aggregate copy PnL: $+3702.33
- Aggregate copy win rate: 81.9%
- Worst included CopyROI: 41.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x54fc...76eb` | B | 100.0 | 31.3 | 617.6 | 46 | 87% | 211.4% | $+232.59 | 6 | 66.7% | $639 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.6 | 580.4 | 487 | 95% | 126.5% | $+316.38 | 20 | 90.0% | $256 | existing,holder |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xa35c...a113` | B | 100.0 | 32.7 | 487.2 | 159 | 100% | 75.7% | $+333.23 | 40 | 95.0% | $705 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.7 | 442.5 | 48 | 32% | 57.0% | $+267.89 | 45 | 80.0% | $1407 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.9 | 415.8 | 192 | 98% | 48.1% | $+418.37 | 59 | 79.7% | $1108 | existing |
| `0x760f...326a` | B | 100.0 | 31.0 | 413.9 | 15 | 14% | 84.4% | $+135.00 | 13 | 76.9% | $5078 | existing |
| `0xa931...1387` | B | 100.0 | 27.4 | 413.5 | 9 | 75% | 145.1% | $+43.53 | 3 | 100.0% | $2517 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x68cb...6f57` | A | 100.0 | 21.1 | 392.9 | 72 | 64% | 48.2% | $+192.68 | 33 | 63.6% | $576 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x95c1...94d7` | B | 100.0 | 27.1 | 375.3 | 54 | 50% | 41.8% | $+137.97 | 31 | 93.5% | $885 | existing,holder |
| `0x119e...cb14` | B | 100.0 | 31.5 | 374.7 | 31 | 91% | 84.4% | $+92.81 | 5 | 100.0% | $1802 | existing |
| `0x0931...e78e` | B | 100.0 | 25.8 | 373.9 | 125 | 69% | 59.4% | $+124.82 | 20 | 65.0% | $298 | existing |
| `0x8fd6...c2b8` | B | 100.0 | 6.7 | 366.8 | 16 | 32% | 94.6% | $+37.84 | 4 | 100.0% | $633 | existing |
| `0x6d15...04c2` | B | 100.0 | 26.5 | 359.5 | 43 | 83% | 58.1% | $+87.15 | 12 | 58.3% | $520 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.5 | 790.7 | 374 | 279 | 37% | 137.2% | $+1001.56 | 70 | 97.1% | 157.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5135...7217` | B | 100.0 | 32.4 | 746.0 | 1359 | 366 | 100% | 154.3% | $+1357.94 | 34 | 100.0% | 154.3% | existing,holder | burst_trading,open_copy_exposure |
| `0x6e32...ab65` | B | 100.0 | 34.7 | 669.9 | 842 | 306 | 67% | 147.9% | $+635.95 | 25 | 100.0% | 140.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 341.4 | 97 | 73 | 86% | 46.4% | $+115.87 | 24 | 75.0% | 39.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x885b...a30d` | B | 100.0 | 32.9 | 338.6 | 89 | 36 | 74% | 97.1% | $+58.28 | 5 | 80.0% | 87.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x13c8...a814` | B | 100.0 | 32.8 | 333.3 | 422 | 232 | 75% | 39.0% | $+225.89 | 50 | 68.0% | 34.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7795...9a17` | B | 100.0 | 29.9 | 317.5 | 109 | 105 | 60% | 55.6% | $+94.59 | 16 | 43.8% | 46.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x1abb...0901` | B | 100.0 | 26.7 | 312.3 | 30 | 28 | 71% | 73.1% | $+43.84 | 6 | 66.7% | 71.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.0 | 300.8 | 318 | 126 | 29% | 43.3% | $+190.56 | 20 | 35.0% | 63.6% | existing | burst_trading,open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 43.8 | 207.0 | 1262 | 90% | 0.0% | 0 | 423 | $561 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 36.7 | 192.6 | 0 | 0% | 123.1% | 1 | 334 | $8934 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x6ac5...4b6e` | C | 100.0 | 38.2 | 243.8 | 1089 | 271 | 79% | 198.4% | 11 | 184.8% | 17 | $399 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x54f0...7725` | A | 100.0 | 24.2 | 220.5 | 273 | 171 | 84% | 12.5% | 49 | 11.7% | 56 | $658 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.0 | 215.9 | 682 | 325 | 91% | 42.7% | 22 | 36.7% | 26 | $144 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.4 | 212.7 | 591 | 360 | 89% | 28.8% | 82 | 27.2% | 89 | $225 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2248...bc0e` | B | 100.0 | 29.7 | 200.8 | 151 | 115 | 84% | 26.8% | 13 | 26.8% | 13 | $403 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6d36...f787` | C | 100.0 | 41.6 | 199.2 | 1191 | 280 | 88% | 147.8% | 24 | 149.1% | 30 | $313 | existing | burst_trading,opposite_side_same_market |
| `0x2fbb...7e44` | B | 100.0 | 28.8 | 197.2 | 127 | 55 | 100% | 34.0% | 14 | 34.0% | 14 | $158 | existing | opposite_side_same_market |
| `0x076d...8d4c` | C | 100.0 | 44.3 | 191.1 | 397 | 105 | 34% | 106.3% | 36 | 113.0% | 69 | $237 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xae2b...90ea` | C | 100.0 | 44.9 | 185.8 | 343 | 310 | 39% | 92.0% | 58 | 109.2% | 92 | $1951 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 23.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 565.1 | opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x076d...8d4c` | C | 100.0 | 44.3 | 15m edge -6.00pp over 4 samples | 113.0% | 69 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x44c4...09cb` | B | 100.0 | 27.0 | 1h edge -15.87pp over 1 samples | 107.2% | 43 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 73.0% | 60 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.9 | 1h edge -72.59pp over 1 samples | 48.1% | 59 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x8002...8b6c` | reject-bot | - | BOT | 77.4 | $8003 | $8003 | 1 | -11.0% | 22 | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 2
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.3 | 738.4 | 415 | 295 | 100% | 101.7% | 93 | 101.7% | 93 | existing,holder,sports_tape | opposite_side_same_market |
| `0x8002...8b6c` | reject-bot | - | BOT | 93.2 | 77.4 | -72.1 | 290 | 154 | 78% | -11.2% | 16 | -11.0% | 22 | existing,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 20 | 457 | 99.3% | $+5333.00 | 84.7% | 107.8% | 62.3% | 3.54x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 19 | 447 | 100.1% | $+5264.43 | 85.0% | 108.4% | 65.3% | 3.51x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 22 | 543 | 92.4% | $+5795.11 | 84.3% | 105.6% | 41.8% | 3.15x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 4 | 25 | 572 | 89.5% | $+6025.27 | 83.7% | 101.7% | 41.8% | 3.22x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 5 | 26 | 605 | 87.2% | $+6217.95 | 82.6% | 100.6% | 41.8% | 3.06x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 19 | 449 | 98.4% | $+5158.02 | 85.1% | 107.2% | 62.3% | 3.05x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 23 | 576 | 89.8% | $+5987.79 | 83.2% | 104.0% | 41.8% | 2.98x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 18 | 439 | 99.2% | $+5089.45 | 85.4% | 107.8% | 65.3% | 3.01x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 20 | 515 | 92.7% | $+5495.31 | 85.4% | 105.6% | 41.8% | 2.71x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 13 | 288 | 115.6% | $+3839.32 | 84.0% | 109.8% | 101.7% | 2.59x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 22 | 535 | 90.6% | $+5643.37 | 85.0% | 102.8% | 41.8% | 2.88x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 23 | 548 | 89.2% | $+5706.21 | 85.0% | 101.7% | 37.0% | 2.81x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
