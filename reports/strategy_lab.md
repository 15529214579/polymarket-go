# Strategy Lab Report

**Generated:** 2026-08-22 23:20 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (416)
- Valid strategies found: 72
- Candidate layers: 11 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 47 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 185
- Aggregate copy ROI: 128.7%
- Aggregate copy PnL: $+2716.37
- Aggregate copy win rate: 87.6%
- Median wallet CopyROI: 124.0%
- Worst included wallet CopyROI: 104.0%
- Open copy cost / closed copy capital: 5.38x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 23.57 | 69 | 66% | 143.1% | $+271.82 | 19 | 100.0% | 43.5% |
| `0x84cd...7565` | A | 23.84 | 121 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 26.40 | 115 | 62% | 133.3% | $+479.87 | 35 | 82.9% | 55.3% |
| `0x0a7c...8964` | B | 27.22 | 991 | 85% | 180.8% | $+234.98 | 8 | 87.5% | 56.8% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 603
- Aggregate copy ROI: 119.7%
- Aggregate copy PnL: $+8788.86
- Aggregate copy win rate: 86.4%
- Worst included CopyROI: 45.0%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9c03...522d` | B | 100.0 | 33.2 | 1204.7 | 190 | 65% | 244.1% | $+1439.94 | 57 | 96.5% | $756 | existing |
| `0x4362...9e47` | B | 100.0 | 34.9 | 844.4 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1329 | existing |
| `0x3c2b...825e` | B | 100.0 | 34.1 | 829.7 | 554 | 96% | 150.4% | $+1127.82 | 66 | 95.5% | $299 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 32.3 | 777.7 | 56 | 29% | 150.8% | $+874.55 | 43 | 69.8% | $516 | existing |
| `0x888b...5d1a` | B | 100.0 | 34.9 | 730.6 | 37 | 31% | 146.5% | $+541.93 | 33 | 90.9% | $3106 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.2 | 588.3 | 117 | 98% | 106.1% | $+519.73 | 32 | 87.5% | $839 | existing,holder |
| `0x760f...326a` | B | 100.0 | 30.7 | 563.3 | 15 | 14% | 121.1% | $+290.65 | 20 | 85.0% | $5102 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x18c2...529a` | B | 100.0 | 27.8 | 495.1 | 917 | 80% | 70.9% | $+595.82 | 62 | 91.9% | $303 | existing |
| `0x092b...614e` | B | 100.0 | 25.4 | 494.3 | 149 | 97% | 97.5% | $+185.35 | 18 | 83.3% | $1033 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x7fcf...80ac` | A | 100.0 | 10.8 | 443.2 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | $1826 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.6 | 431.0 | 48 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1331 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.7 | 415.2 | 222 | 98% | 45.0% | $+454.87 | 72 | 77.8% | $1180 | existing,holder |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.9 | 1034.5 | 384 | 293 | 38% | 183.1% | $+1684.64 | 87 | 97.7% | 183.2% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.6 | 670.6 | 774 | 284 | 62% | 132.1% | $+1188.59 | 34 | 100.0% | 129.2% | existing,holder | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 33.8 | 583.1 | 1373 | 378 | 100% | 106.9% | $+1090.68 | 36 | 100.0% | 106.9% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.0 | 517.7 | 627 | 272 | 93% | 108.5% | $+271.29 | 21 | 100.0% | 107.7% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 32.8 | 490.6 | 1186 | 288 | 100% | 104.3% | $+792.39 | 16 | 100.0% | 104.3% | existing,holder | burst_trading,open_copy_exposure |
| `0x21cc...54bc` | B | 100.0 | 28.7 | 424.6 | 170 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 127 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 375.8 | 72 | 58 | 60% | 54.0% | $+156.56 | 22 | 63.6% | 50.5% | existing | opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.8 | 364.4 | 565 | 301 | 56% | 64.8% | $+123.12 | 17 | 82.3% | 57.3% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 36.2 | 231.9 | 1357 | 96% | 395.0% | 1 | 404 | $609 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 279.7 | 882 | 636 | 96% | 10.2% | 10 | 9.4% | 11 | $503 | existing | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.5 | 219.0 | 771 | 356 | 92% | 28.6% | 21 | 24.5% | 25 | $137 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 215.1 | 612 | 378 | 89% | 29.1% | 86 | 27.6% | 93 | $233 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd052...563d` | B | 100.0 | 35.0 | 211.6 | 305 | 265 | 70% | 17.4% | 76 | 22.4% | 95 | $1086 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 41.4 | 208.9 | 1111 | 347 | 82% | 121.7% | 23 | 126.8% | 28 | $392 | existing,holder | burst_trading,open_copy_exposure |
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

- Wallets: 5
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.4 | 1h edge -15.87pp over 1 samples | 133.3% | 35 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.8 | 1h edge -52.95pp over 1 samples | 70.9% | 62 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.7 | 1h edge -72.59pp over 1 samples | 45.0% | 72 | existing,holder | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xf1ef...8cd6` | reject-bot | - | D | 59.1 | $5600 | $5600 | 1 | 50.6% | 22 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf1ef...8cd6` | reject | - | D | 100.0 | 59.1 | 139.1 | 125 | 54 | 9% | 46.4% | 6 | 50.6% | 22 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 185 | 128.7% | $+2716.37 | 87.6% | 124.0% | 104.0% | 5.38x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 10 | 175 | 129.5% | $+2512.61 | 88.6% | 128.6% | 104.0% | 4.29x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 3 | 15 | 301 | 107.4% | $+3790.04 | 88.4% | 115.4% | 65.3% | 4.61x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 16 | 310 | 105.7% | $+3869.58 | 88.4% | 112.2% | 61.2% | 4.48x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 19 | 380 | 96.2% | $+4241.40 | 87.9% | 108.4% | 44.4% | 4.93x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=10 smart>=70 |
| 6 | 10 | 177 | 125.3% | $+2481.39 | 87.6% | 121.9% | 104.0% | 3.28x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 19 | 391 | 95.0% | $+4295.02 | 88.2% | 108.4% | 44.4% | 4.61x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 20 | 400 | 94.0% | $+4334.96 | 88.2% | 106.2% | 44.4% | 4.71x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 18 | 371 | 97.3% | $+4201.46 | 87.9% | 108.7% | 44.4% | 4.82x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=10 smart>=70 |
| 10 | 20 | 400 | 94.0% | $+4408.38 | 86.2% | 106.2% | 44.4% | 4.58x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 21 | 409 | 93.1% | $+4448.32 | 86.3% | 104.0% | 44.4% | 4.68x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 20 | 397 | 93.5% | $+4284.20 | 87.4% | 106.2% | 25.2% | 4.76x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=10 smart>=70 |
