# Strategy Lab Report

**Generated:** 2026-08-23 23:02 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (418)
- Valid strategies found: 54
- Candidate layers: 13 core + 20 watch + 10 sports + 1 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 47 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 13
- Aggregate closed copy trades: 282
- Aggregate copy ROI: 102.9%
- Aggregate copy PnL: $+3395.61
- Aggregate copy win rate: 88.7%
- Median wallet CopyROI: 115.4%
- Worst included wallet CopyROI: 65.3%
- Open copy cost / closed copy capital: 3.39x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x162d...8944` | A | 24.89 | 3 | 1% | 104.0% | $+332.66 | 31 | 90.3% | 46.5% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 24.0% |
| `0x5be9...ef2b` | A | 24.02 | 73 | 65% | 137.2% | $+288.21 | 21 | 100.0% | 45.6% |
| `0x84cd...7565` | A | 23.84 | 121 | 83% | 115.4% | $+184.60 | 13 | 84.6% | 41.8% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 35.3% |
| `0x6f16...5fe7` | A | 23.16 | 38 | 67% | 161.6% | $+161.64 | 10 | 80.0% | 35.6% |
| `0x7fcf...80ac` | A | 10.77 | 18 | 34% | 99.5% | $+109.51 | 8 | 75.0% | 36.5% |
| `0x18c2...529a` | B | 27.87 | 913 | 80% | 72.1% | $+619.78 | 63 | 92.1% | 46.2% |
| `0x44c4...09cb` | B | 26.34 | 116 | 62% | 133.3% | $+479.87 | 35 | 82.9% | 55.3% |
| `0x2b17...d36d` | B | 25.40 | 417 | 73% | 119.9% | $+203.76 | 10 | 70.0% | 53.7% |
| `0x092b...614e` | B | 25.29 | 150 | 97% | 97.5% | $+185.35 | 18 | 83.3% | 61.6% |
| `0x21cc...54bc` | B | 28.66 | 170 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 32.1% |
| `0x819d...6e9c` | B | 27.31 | 26 | 26% | 124.0% | $+148.79 | 11 | 100.0% | 46.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 533
- Aggregate copy ROI: 99.9%
- Aggregate copy PnL: $+6643.49
- Aggregate copy win rate: 82.6%
- Worst included CopyROI: 44.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x4362...9e47` | B | 100.0 | 34.9 | 844.4 | 137 | 47% | 172.5% | $+724.43 | 41 | 92.7% | $1329 | existing |
| `0x3c2b...825e` | B | 100.0 | 33.9 | 803.6 | 563 | 96% | 143.3% | $+1117.96 | 68 | 92.7% | $313 | existing,holder |
| `0x3b5a...87ef` | B | 100.0 | 33.1 | 782.3 | 59 | 29% | 150.7% | $+904.09 | 45 | 68.9% | $502 | existing |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 662.9 | 49 | 96% | 196.7% | $+236.02 | 9 | 88.9% | $6758 | existing |
| `0x54fc...76eb` | B | 100.0 | 30.7 | 646.9 | 51 | 88% | 210.5% | $+252.56 | 7 | 71.4% | $828 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.2 | 585.3 | 117 | 98% | 106.1% | $+519.73 | 32 | 87.5% | $839 | existing |
| `0x760f...326a` | B | 100.0 | 30.7 | 563.3 | 15 | 14% | 121.1% | $+290.65 | 20 | 85.0% | $5102 | existing |
| `0x17fe...b0ca` | B | 100.0 | 33.5 | 562.4 | 118 | 100% | 110.2% | $+275.50 | 24 | 83.3% | $571 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 532.6 | 49 | 84% | 129.7% | $+194.50 | 12 | 83.3% | $1689 | existing |
| `0xe760...f2f9` | B | 100.0 | 30.6 | 503.2 | 137 | 100% | 69.1% | $+400.80 | 40 | 87.5% | $5431 | existing,leaderboard_profit_7d |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | B | 100.0 | 31.6 | 430.8 | 48 | 30% | 57.0% | $+267.89 | 45 | 80.0% | $1326 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x0e24...7014` | B | 100.0 | 31.3 | 421.1 | 233 | 98% | 46.5% | $+492.83 | 77 | 77.9% | $1186 | existing |
| `0x68cb...6f57` | A | 100.0 | 20.9 | 402.4 | 72 | 60% | 50.5% | $+207.07 | 34 | 64.7% | $547 | existing |
| `0x66ef...a836` | B | 100.0 | 18.0 | 392.6 | 32 | 97% | 122.4% | $+48.97 | 3 | 100.0% | $943 | existing |
| `0x95c1...94d7` | B | 100.0 | 27.3 | 386.9 | 56 | 50% | 44.4% | $+159.96 | 34 | 94.1% | $1068 | existing |
| `0x8bd3...f69a` | B | 100.0 | 25.2 | 386.9 | 20 | 32% | 67.1% | $+107.43 | 11 | 72.7% | $5314 | existing |
| `0x2c63...0ca5` | A | 100.0 | 17.4 | 383.6 | 25 | 96% | 92.2% | $+46.10 | 5 | 80.0% | $1330 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 31.7 | 1043.5 | 388 | 299 | 38% | 181.5% | $+1814.84 | 94 | 97.9% | 180.8% | existing | opposite_side_same_market |
| `0x6e32...ab65` | B | 100.0 | 34.9 | 668.6 | 689 | 244 | 54% | 139.3% | $+1072.32 | 28 | 100.0% | 133.0% | existing | burst_trading,open_copy_exposure |
| `0x5135...7217` | B | 100.0 | 34.4 | 608.7 | 1365 | 362 | 100% | 109.2% | $+1310.34 | 40 | 100.0% | 109.2% | existing | burst_trading,open_copy_exposure |
| `0x4cb0...57e0` | B | 100.0 | 33.9 | 545.1 | 632 | 277 | 93% | 110.9% | $+321.63 | 25 | 100.0% | 110.1% | existing,holder | burst_trading,open_copy_exposure |
| `0x4b59...3aa6` | B | 100.0 | 33.0 | 512.0 | 1190 | 292 | 100% | 108.7% | $+891.57 | 17 | 100.0% | 108.7% | existing,holder | burst_trading,open_copy_exposure |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 467.8 | 1161 | 1109 | 100% | 62.6% | $+1051.00 | 84 | 92.9% | 62.6% | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.0 | 406.0 | 127 | 72 | 29% | 75.4% | $+150.82 | 18 | 72.2% | 86.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.8 | 379.6 | 125 | 71 | 69% | 71.2% | $+106.79 | 15 | 73.3% | 59.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.8 | 370.4 | 577 | 306 | 56% | 64.7% | $+135.77 | 19 | 79.0% | 57.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x65c1...2988` | C | 100.0 | 36.2 | 229.1 | 1367 | 97% | 395.0% | 1 | 383 | $539 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x3724...69b5` | A | 100.0 | 24.9 | 288.1 | 938 | 679 | 96% | 10.2% | 10 | 9.4% | 11 | $505 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x682c...ef07` | C | 100.0 | 35.6 | 220.4 | 799 | 367 | 92% | 24.7% | 22 | 21.3% | 26 | $137 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x54f0...7725` | A | 100.0 | 24.4 | 219.7 | 287 | 175 | 83% | 10.9% | 50 | 11.0% | 59 | $628 | existing | opposite_side_same_market |
| `0x257a...b42a` | C | 100.0 | 38.5 | 218.8 | 616 | 382 | 89% | 29.1% | 86 | 27.6% | 93 | $236 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x07b1...6dfc` | C | 100.0 | 41.1 | 209.8 | 1351 | 319 | 94% | 108.0% | 28 | 106.9% | 31 | $400 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xadfa...324e` | C | 100.0 | 40.4 | 206.3 | 1052 | 331 | 77% | 110.0% | 20 | 117.6% | 25 | $387 | existing,holder | burst_trading,open_copy_exposure |
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
| `0x44c4...09cb` | B | 100.0 | 26.3 | 1h edge -15.87pp over 1 samples | 133.3% | 35 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 34.6 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.5 | 1h edge -47.28pp over 1 samples | 62.6% | 84 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 27.9 | 1h edge -52.95pp over 1 samples | 72.1% | 63 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 31.3 | 1h edge -72.59pp over 1 samples | 46.5% | 77 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0x34dd...fd83` | reject-flow | - | C | 37.5 | $375861 | $465978 | 3 | 116.0% | 1 | burst_trading,fixed_price |

## Sports Tape Candidate Review

- Wallets: 1
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x34dd...fd83` | watch | - | C | 92.9 | 37.5 | 154.0 | 57 | 27 | 100% | 116.0% | 1 | 116.0% | 1 | existing,leaderboard_profit_30d,leaderboard_profit_7d,sports_tape | burst_trading,fixed_price |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 13 | 282 | 102.9% | $+3395.61 | 88.7% | 115.4% | 65.3% | 3.39x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 12 | 274 | 103.0% | $+3286.10 | 89.1% | 117.6% | 65.3% | 3.42x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 16 | 356 | 92.7% | $+3827.58 | 87.9% | 106.4% | 44.4% | 3.80x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 15 | 345 | 94.0% | $+3740.14 | 87.8% | 108.9% | 44.4% | 3.92x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 17 | 365 | 91.6% | $+3867.52 | 87.9% | 104.0% | 44.4% | 3.93x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 16 | 365 | 92.2% | $+3864.96 | 86.6% | 106.4% | 44.4% | 3.87x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 7 | 18 | 385 | 89.7% | $+4034.50 | 86.2% | 101.8% | 44.4% | 3.64x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 14 | 337 | 93.8% | $+3630.63 | 88.1% | 112.2% | 44.4% | 3.96x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 11 | 264 | 102.1% | $+3082.34 | 89.8% | 115.4% | 65.3% | 2.61x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 10 | 15 | 348 | 92.5% | $+3718.07 | 88.2% | 108.9% | 44.4% | 3.83x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 11 | 19 | 394 | 88.8% | $+4074.44 | 86.3% | 99.5% | 44.4% | 3.76x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 16 | 368 | 90.8% | $+3842.89 | 87.0% | 106.4% | 44.4% | 3.78x | tier>=B bot<30 copyT>=10 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
