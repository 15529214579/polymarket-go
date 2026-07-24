# Strategy Lab Report

**Generated:** 2026-07-24 22:55 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 46
- Candidate layers: 13 core + 20 watch + 10 sports + 0 scout + 3 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 38 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 13
- Aggregate closed copy trades: 243
- Aggregate copy ROI: 126.2%
- Aggregate copy PnL: $+3571.19
- Aggregate copy win rate: 82.7%
- Median wallet CopyROI: 135.2%
- Worst included wallet CopyROI: 100.3%
- Open copy cost / closed copy capital: 2.73x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.49 | 117 | 42% | 143.7% | $+545.89 | 32 | 65.6% | 0.0% |
| `0xa75b...772c` | A | 24.26 | 127 | 99% | 103.8% | $+467.06 | 33 | 78.8% | 0.0% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 0.0% |
| `0x2952...f50d` | A | 23.79 | 79 | 55% | 107.3% | $+343.22 | 24 | 87.5% | 0.0% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 0.0% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 0.0% |
| `0x7124...f0b5` | A | 22.50 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | 0.0% |
| `0xabff...9e8f` | A | 24.76 | 57 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 0.0% |
| `0xb715...c3bb` | A | 24.63 | 60 | 90% | 135.2% | $+148.75 | 10 | 80.0% | 0.0% |
| `0x73e2...46ec` | A | 5.24 | 12 | 19% | 129.2% | $+142.10 | 11 | 90.9% | 0.0% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 0.0% |
| `0x44c4...09cb` | B | 26.41 | 158 | 60% | 109.3% | $+601.01 | 50 | 80.0% | 0.0% |
| `0xeb8b...6d8a` | B | 28.11 | 27 | 51% | 138.5% | $+124.64 | 9 | 100.0% | 0.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 717
- Aggregate copy ROI: 100.4%
- Aggregate copy PnL: $+9953.30
- Aggregate copy win rate: 90.0%
- Worst included CopyROI: 60.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 69.9 | 34.3 | 863.2 | 25 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1207 | existing |
| `0xfdff...4adc` | B | 74.2 | 33.7 | 721.7 | 17 | 24% | 251.8% | $+251.77 | 8 | 75.0% | $2953 | existing |
| `0xeca2...44fc` | B | 74.9 | 32.5 | 721.5 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 69.2 | 28.3 | 703.5 | 39 | 95% | 248.6% | $+223.76 | 7 | 100.0% | $7773 | existing |
| `0xbb35...b62a` | B | 61.9 | 34.3 | 668.9 | 1341 | 100% | 109.9% | $+2494.88 | 102 | 96.1% | $1032 | existing |
| `0xc367...3066` | B | 67.6 | 31.2 | 637.2 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0x89cf...5f47` | B | 68.4 | 33.4 | 612.8 | 68 | 41% | 119.7% | $+526.61 | 43 | 88.4% | $649 | existing |
| `0xb05d...dcdc` | B | 73.8 | 34.4 | 612.1 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xcc6e...fa6f` | B | 78.5 | 26.0 | 587.9 | 289 | 83% | 96.9% | $+639.67 | 58 | 89.7% | $2147 | existing,sports_tape |
| `0x2a35...9015` | B | 66.5 | 33.1 | 569.9 | 107 | 99% | 118.1% | $+401.60 | 26 | 92.3% | $1826 | existing,holder |
| `0xf4e1...34a5` | B | 70.0 | 30.6 | 534.7 | 68 | 41% | 115.2% | $+322.47 | 23 | 82.6% | $5973 | existing |
| `0x9caf...94dc` | B | 75.4 | 31.5 | 531.1 | 64 | 50% | 103.7% | $+300.79 | 28 | 96.4% | $2865 | existing |
| `0xde24...4ded` | B | 67.8 | 30.8 | 527.3 | 145 | 100% | 103.8% | $+508.64 | 25 | 100.0% | $2689 | existing |
| `0xb36f...53d0` | B | 76.0 | 30.4 | 500.1 | 882 | 93% | 60.1% | $+1219.76 | 184 | 89.1% | $298 | existing,holder |
| `0x7992...1fc1` | A | 82.5 | 18.6 | 500.0 | 21 | 21% | 86.7% | $+294.71 | 29 | 79.3% | $2688 | existing |
| `0x141a...d05a` | B | 68.3 | 29.9 | 463.5 | 169 | 100% | 74.3% | $+616.98 | 37 | 94.6% | $1675 | existing |
| `0x84cd...7565` | A | 79.7 | 23.7 | 454.2 | 116 | 72% | 93.7% | $+187.34 | 17 | 76.5% | $433 | existing |
| `0x18c2...529a` | B | 67.2 | 28.3 | 452.2 | 563 | 70% | 80.9% | $+396.49 | 39 | 82.0% | $271 | existing |
| `0xd5b1...1b71` | B | 67.7 | 31.0 | 452.0 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |
| `0x2cf0...df85` | B | 68.6 | 32.9 | 441.3 | 144 | 43% | 88.5% | $+274.36 | 24 | 87.5% | $359 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x59a8...eb4e` | B | 74.5 | 33.2 | 390.3 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5a56...10f5` | B | 71.5 | 31.8 | 389.5 | 90 | 75 | 87% | 94.5% | $+245.67 | 11 | 63.6% | 88.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb916...f248` | A | 72.9 | 21.5 | 378.3 | 143 | 113 | 75% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0x578e...c3c0` | A | 72.0 | 23.2 | 374.4 | 70 | 66 | 36% | 60.5% | $+151.16 | 24 | 83.3% | 55.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 74.8 | 25.4 | 369.8 | 1019 | 800 | 86% | 76.6% | $+199.20 | 15 | 73.3% | 74.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 72.0 | 23.1 | 359.9 | 73 | 69 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 58.4% | existing | opposite_side_same_market |
| `0x0931...e78e` | B | 72.8 | 25.5 | 347.9 | 124 | 71 | 67% | 71.2% | $+106.79 | 15 | 73.3% | 55.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 70.2 | 33.7 | 341.4 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x07a9...32c8` | B | 71.1 | 27.5 | 339.5 | 35 | 33 | 46% | 107.8% | $+64.68 | 5 | 100.0% | 83.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x53ad...e49e` | B | 78.1 | 26.1 | 337.0 | 59 | 52 | 82% | 94.2% | $+75.36 | 6 | 83.3% | 84.5% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 0
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No leaderboard-only wallets passed the scout filters.

## Target Category Scout

- Wallets: 3
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xdbd5...cba7` | A | 74.8 | 18.0 | 207.2 | 93 | 89 | 93% | 40.5% | 18 | 41.6% | 19 | $2728 | existing | open_copy_exposure,opposite_side_same_market |
| `0xc3ed...9f79` | C | 72.1 | 37.5 | 177.3 | 144 | 142 | 46% | 50.0% | 24 | 55.1% | 56 | $8991 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |
| `0x42f1...f53f` | B | 74.8 | 27.9 | 158.1 | 73 | 64 | 65% | 39.3% | 22 | 39.4% | 26 | $671 | existing | open_copy_exposure,opposite_side_same_market |

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

- Wallets: 3
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 24.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 529.7 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.0 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 321.2 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 74.8 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 61.9 | 34.3 | 15m edge -2.84pp over 3 samples | 109.9% | 102 | existing | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 67.8 | 30.8 | 1h edge -11.42pp over 4 samples | 103.8% | 25 | existing | open_copy_exposure,opposite_side_same_market |
| `0x07a9...32c8` | B | 71.1 | 27.5 | 1h edge -12.31pp over 2 samples | 83.2% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 78.2 | 26.4 | 1h edge -15.87pp over 1 samples | 109.3% | 50 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 69.2 | 28.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 75.6 | 23.5 | 1h edge -19.00pp over 1 samples | 143.7% | 32 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 67.2 | 28.3 | 1h edge -52.95pp over 1 samples | 80.9% | 39 | existing | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | D | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | D | 30.5 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | C | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 73.7 | $7000 | $7000 | 1 | 116.7% | 121 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | BOT | 48.3 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 46.2 | 73.7 | 590.8 | 234 | 226 | 16% | 121.7% | 32 | 116.7% | 121 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbe7e...6d00` | reject-bot | - | BOT | 36.4 | 64.3 | 582.1 | 392 | 344 | 34% | 149.0% | 17 | 122.3% | 42 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 71.4 | 24.3 | 568.1 | 127 | 105 | 99% | 103.8% | 33 | 103.8% | 33 | existing,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 78.5 | 26.0 | 544.2 | 289 | 244 | 83% | 83.1% | 48 | 96.9% | 58 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 41.4 | 69.8 | 253.9 | 1012 | 908 | 74% | 29.8% | 310 | 29.2% | 397 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 51.2 | 41.5 | 229.2 | 163 | 151 | 16% | 37.5% | 59 | 36.3% | 270 | sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 50.9 | 59.7 | 190.1 | 139 | 136 | 11% | 32.0% | 44 | 47.1% | 315 | sports_tape | bot_like_flow,extreme_price_heavy,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | BOT | 48.3 | 58.0 | 163.6 | 263 | 238 | 78% | 29.5% | 85 | 25.2% | 110 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 53.5 | 30.0 | 147.3 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 57.8 | 25.3 | 146.9 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | reject | - | D | 14.0 | 20.0 | 124.0 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x7af7...4c89` | reject | - | D | 19.2 | 30.5 | 95.8 | 21 | 20 | 66% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xbea4...fac9` | prompt | - | B | 59.4 | 34.1 | 44.7 | 4 | 4 | 7% | 44.0% | 1 | 24.7% | 7 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xe0cc...aafc` | reject | - | D | 0.0 | 10.0 | 35.4 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xaa9e...e588` | reject | - | D | 0.0 | 16.2 | 30.3 | 4 | 4 | 100% | 33.6% | 1 | 33.6% | 1 | sports_tape | fixed_price |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 13 | 243 | 126.2% | $+3571.19 | 82.7% | 135.2% | 100.3% | 2.73x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 11 | 226 | 127.4% | $+3285.98 | 82.7% | 135.2% | 103.8% | 2.77x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 21 | 436 | 107.6% | $+5498.45 | 83.5% | 107.3% | 66.5% | 4.05x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 13 | 232 | 121.2% | $+3272.64 | 81.9% | 129.2% | 81.2% | 2.66x | tier>=A bot<25 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 11 | 184 | 129.9% | $+2845.54 | 82.6% | 135.2% | 100.3% | 2.83x | tier>=A bot<25 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 11 | 203 | 125.1% | $+2864.73 | 86.2% | 135.2% | 103.8% | 2.60x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 13 | 280 | 113.6% | $+3430.11 | 88.9% | 129.2% | 66.5% | 2.71x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 8 | 10 | 176 | 132.3% | $+2684.97 | 83.5% | 135.4% | 103.8% | 2.88x | tier>=A bot<25 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 19 | 419 | 107.3% | $+5213.24 | 83.5% | 107.3% | 66.5% | 4.14x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 14 | 240 | 120.0% | $+3433.21 | 81.2% | 118.8% | 81.2% | 2.64x | tier>=A bot<25 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 10 | 170 | 130.3% | $+2397.67 | 87.6% | 135.4% | 107.3% | 3.11x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 12 | 26 | 574 | 92.8% | $+6539.97 | 82.6% | 98.6% | 41.6% | 3.28x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
