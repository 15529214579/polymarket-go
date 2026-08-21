# Strategy Lab Report

**Generated:** 2026-08-21 23:04 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (413)
- Valid strategies found: 78
- Candidate layers: 11 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 47 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 190
- Aggregate copy ROI: 121.5%
- Aggregate copy PnL: $+2600.01
- Aggregate copy win rate: 85.8%
- Median wallet CopyROI: 119.9%
- Worst included wallet CopyROI: 104.0%
- Open copy cost / closed copy capital: 3.21x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 15.30 | 65 | 65% | 149.5% | $+254.13 | 17 | 100.0% | 42.6% |
| `0x84cd...7565` | A | 23.66 | 123 | 83% | 109.8% | $+186.61 | 14 | 85.7% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 46.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 28.4% |
| `0x44c4...09cb` | B | 26.25 | 122 | 62% | 110.8% | $+465.42 | 39 | 76.9% | 55.3% |
| `0x6f16...5fe7` | B | 25.16 | 41 | 65% | 192.8% | $+231.31 | 12 | 83.3% | 35.6% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 558
- Aggregate copy ROI: 111.7%
- Aggregate copy PnL: $+7717.70
- Aggregate copy win rate: 84.2%
- Worst included CopyROI: 44.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x9c03...522d` | B | 100.0 | 33.3 | 1207.6 | 190 | 66% | 244.1% | $+1439.94 | 57 | 96.5% | $758 | existing,holder |
| `0x3c2b...825e` | B | 100.0 | 34.3 | 829.1 | 541 | 96% | 150.4% | $+1127.82 | 66 | 95.5% | $284 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 32.3 | 777.7 | 56 | 29% | 150.8% | $+874.55 | 43 | 69.8% | $516 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.4 | 609.3 | 111 | 98% | 114.7% | $+527.80 | 30 | 90.0% | $861 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0x760f...326a` | B | 100.0 | 30.8 | 524.0 | 15 | 14% | 112.2% | $+246.72 | 18 | 83.3% | $4990 | existing |
| `0x092b...614e` | B | 100.0 | 25.4 | 497.3 | 149 | 97% | 97.5% | $+185.35 | 18 | 83.3% | $1033 | existing,holder |
| `0x18c2...529a` | B | 100.0 | 27.8 | 495.2 | 918 | 80% | 70.4% | $+591.26 | 62 | 90.3% | $301 | existing,holder |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x7fcf...80ac` | A | 100.0 | 10.8 | 443.2 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | $1826 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.7 | 431.2 | 48 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1337 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 30.5 | 408.1 | 215 | 98% | 44.8% | $+434.45 | 68 | 76.5% | $1154 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.9 | 952.8 | 367 | 276 | 36% | 172.0% | $+1324.23 | 74 | 97.3% | 175.9% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.1 | 750.5 | 786 | 293 | 63% | 146.4% | $+1376.61 | 41 | 100.0% | 138.2% | existing,holder | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 33.4 | 589.6 | 1384 | 365 | 100% | 106.4% | $+1127.63 | 38 | 100.0% | 106.4% | existing,holder | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 34.1 | 504.3 | 618 | 263 | 93% | 109.8% | $+241.48 | 18 | 100.0% | 108.8% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 32.3 | 450.8 | 1185 | 297 | 100% | 102.0% | $+479.53 | 12 | 100.0% | 102.0% | existing,holder | burst_trading,open_copy_exposure |
| `0x0a7c...8964` | B | 100.0 | 34.3 | 449.3 | 994 | 404 | 85% | 133.4% | $+146.70 | 7 | 85.7% | 138.6% | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.7 | 424.6 | 170 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 127 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | B | 100.0 | 34.5 | 238.9 | 1290 | 91% | 395.0% | 1 | 412 | $623 | existing,holder,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.4 | 218.7 | 755 | 351 | 92% | 33.2% | 22 | 28.6% | 26 | $138 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 217.5 | 605 | 373 | 89% | 28.8% | 85 | 27.3% | 92 | $230 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0xa35c...a113` | C | 100.0 | 35.1 | 206.4 | 139 | 139 | 100% | 75.7% | 40 | 75.7% | 40 | $699 | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.7 | 206.2 | 559 | 299 | 55% | 64.8% | 17 | 55.1% | 26 | $189 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x8e74...5972` | C | 89.9 | 43.8 | 205.4 | 955 | 190 | 97% | 115.1% | 20 | 114.4% | 22 | $407 | existing,holder,leaderboard_profit_7d | burst_trading |
| `0x42f1...f53f` | B | 100.0 | 27.9 | 203.3 | 97 | 73 | 86% | 46.4% | 24 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 43.6 | 203.0 | 1116 | 350 | 83% | 119.6% | 26 | 124.3% | 31 | $387 | existing | burst_trading,open_copy_exposure |
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
| `0x44c4...09cb` | B | 100.0 | 26.2 | 1h edge -15.87pp over 1 samples | 110.8% | 39 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.8 | 1h edge -52.95pp over 1 samples | 70.4% | 62 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.5 | 1h edge -72.59pp over 1 samples | 44.8% | 68 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x88d1...18da` | reject-bot | - | BOT | 64.7 | $5000 | $5000 | 1 | 67.8% | 37 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x88d1...18da` | reject-bot | - | BOT | 100.0 | 64.7 | 211.2 | 585 | 212 | 46% | 43.8% | 13 | 67.8% | 37 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 190 | 121.5% | $+2600.01 | 85.8% | 119.9% | 104.0% | 3.21x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 19 | 392 | 92.1% | $+4108.15 | 87.2% | 108.4% | 44.4% | 3.85x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 19 | 395 | 91.7% | $+4190.24 | 85.1% | 108.4% | 44.4% | 3.70x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 4 | 15 | 306 | 103.1% | $+3669.12 | 86.9% | 109.8% | 65.3% | 3.31x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 20 | 404 | 90.8% | $+4230.18 | 85.1% | 106.2% | 44.4% | 3.82x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 6 | 19 | 392 | 91.2% | $+4066.06 | 86.2% | 108.4% | 25.2% | 3.87x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 7 | 18 | 383 | 93.1% | $+4068.21 | 87.2% | 108.7% | 44.4% | 3.73x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 18 | 375 | 93.8% | $+4023.26 | 86.7% | 108.7% | 44.4% | 4.00x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 9 | 20 | 409 | 89.7% | $+4150.95 | 86.8% | 106.2% | 25.2% | 3.73x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 10 | 18 | 386 | 93.2% | $+4108.14 | 85.5% | 108.7% | 44.4% | 3.82x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 11 | 19 | 420 | 89.5% | $+4315.21 | 83.8% | 108.4% | 44.4% | 3.54x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 20 | 429 | 88.3% | $+4397.31 | 83.4% | 106.2% | 44.4% | 3.43x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=5 smart>=70 |
