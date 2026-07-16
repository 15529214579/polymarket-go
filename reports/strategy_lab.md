# Strategy Lab Report

**Generated:** 2026-07-16 23:16 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 197
- Candidate layers: 19 core + 20 watch + 10 sports + 7 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 13
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 19
- Aggregate closed copy trades: 466
- Aggregate copy ROI: 100.8%
- Aggregate copy PnL: $+5401.28
- Aggregate copy win rate: 83.9%
- Median wallet CopyROI: 99.5%
- Worst included wallet CopyROI: 63.6%
- Open copy cost / closed copy capital: 2.68x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.88 | 115 | 44% | 149.7% | $+553.94 | 31 | 67.7% | 45.5% |
| `0x2952...f50d` | A | 24.29 | 76 | 51% | 112.3% | $+381.87 | 26 | 88.5% | 19.5% |
| `0xe745...5681` | A | 5.93 | 57 | 30% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x7992...1fc1` | A | 19.05 | 21 | 23% | 83.5% | $+267.23 | 27 | 77.8% | 19.6% |
| `0x6f16...5fe7` | A | 23.30 | 34 | 58% | 212.9% | $+191.57 | 9 | 77.8% | 31.0% |
| `0x84cd...7565` | A | 23.92 | 115 | 73% | 93.7% | $+187.34 | 17 | 76.5% | 41.9% |
| `0xa75b...772c` | A | 22.23 | 47 | 100% | 76.1% | $+167.40 | 15 | 73.3% | 28.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x7fcf...80ac` | A | 11.62 | 32 | 41% | 87.2% | $+104.63 | 9 | 66.7% | 38.5% |
| `0xcc6e...fa6f` | B | 26.53 | 301 | 85% | 99.5% | $+706.45 | 62 | 91.9% | 45.6% |
| `0x44c4...09cb` | B | 26.80 | 163 | 59% | 109.5% | $+645.88 | 54 | 81.5% | 48.0% |
| `0xfbe8...bb28` | B | 28.30 | 801 | 93% | 68.5% | $+445.19 | 44 | 86.4% | 48.4% |
| `0x89cf...5f47` | B | 26.14 | 56 | 43% | 106.0% | $+381.50 | 35 | 85.7% | 45.8% |
| `0x18c2...529a` | B | 28.00 | 314 | 59% | 78.8% | $+196.91 | 24 | 79.2% | 40.1% |
| `0x21cc...54bc` | B | 28.98 | 162 | 100% | 71.4% | $+192.74 | 27 | 96.3% | 35.4% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x0931...e78e` | B | 26.06 | 117 | 67% | 63.6% | $+127.20 | 20 | 65.0% | 26.8% |
| `0xec56...1f87` | B | 27.56 | 26 | 30% | 66.5% | $+112.96 | 17 | 100.0% | 22.6% |
| `0xeb8b...6d8a` | B | 28.59 | 23 | 50% | 129.7% | $+103.78 | 8 | 100.0% | 40.0% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 548
- Aggregate copy ROI: 98.3%
- Aggregate copy PnL: $+7818.43
- Aggregate copy win rate: 86.5%
- Worst included CopyROI: 58.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 35.0 | 926.2 | 24 | 16% | 294.5% | $+412.31 | 11 | 81.8% | $1114 | existing |
| `0x94a8...5204` | B | 100.0 | 15.6 | 797.5 | 22 | 27% | 408.3% | $+244.99 | 3 | 100.0% | $2855 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 755.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing,holder |
| `0x5b1d...3721` | B | 100.0 | 29.3 | 740.9 | 34 | 97% | 248.6% | $+223.76 | 7 | 100.0% | $8130 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 676.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing,holder |
| `0xbb35...b62a` | B | 100.0 | 34.6 | 664.9 | 1339 | 100% | 100.1% | $+2292.18 | 100 | 96.0% | $992 | existing,holder |
| `0xde24...4ded` | B | 100.0 | 31.5 | 576.1 | 86 | 100% | 125.8% | $+352.35 | 15 | 100.0% | $2741 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 561.9 | 61 | 50% | 103.7% | $+300.79 | 28 | 96.4% | $2840 | existing |
| `0x2a35...9015` | B | 100.0 | 32.6 | 548.6 | 49 | 100% | 130.0% | $+181.95 | 12 | 91.7% | $1788 | existing,holder |
| `0x73e2...46ec` | A | 100.0 | 5.4 | 542.7 | 11 | 18% | 129.2% | $+142.10 | 11 | 90.9% | $2412 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 507.0 | 736 | 93% | 58.4% | $+981.85 | 153 | 88.2% | $297 | existing |
| `0x2929...1dd0` | B | 100.0 | 29.5 | 498.2 | 75 | 99% | 172.7% | $+155.42 | 4 | 100.0% | $1261 | existing,holder,sports_tape |
| `0xd5b1...1b71` | B | 100.0 | 31.1 | 488.2 | 137 | 99% | 79.0% | $+260.85 | 32 | 93.8% | $478 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 482.1 | 134 | 41% | 89.3% | $+268.06 | 23 | 87.0% | $353 | existing |
| `0xc117...7410` | B | 82.8 | 25.5 | 479.5 | 32 | 80% | 133.9% | $+147.33 | 8 | 50.0% | $2138 | existing |
| `0x5a56...10f5` | B | 100.0 | 33.8 | 469.7 | 46 | 100% | 124.4% | $+211.51 | 7 | 71.4% | $1129 | existing |
| `0x578e...c3c0` | A | 100.0 | 23.0 | 469.4 | 79 | 38% | 59.8% | $+364.97 | 59 | 81.4% | $1202 | existing |
| `0xa8b9...775d` | B | 100.0 | 34.1 | 459.5 | 332 | 94% | 68.6% | $+576.05 | 45 | 53.3% | $1264 | existing |
| `0x7c73...6ee3` | B | 100.0 | 26.6 | 453.3 | 46 | 96% | 119.6% | $+71.74 | 6 | 100.0% | $5556 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.5 | 429.2 | 83 | 64 | 62% | 90.5% | $+189.96 | 11 | 72.7% | 56.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.2 | 416.5 | 72 | 72 | 56% | 75.8% | $+159.16 | 19 | 89.5% | 60.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb916...f248` | A | 100.0 | 21.5 | 404.6 | 142 | 112 | 75% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 386.5 | 73 | 69 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 58.4% | existing | opposite_side_same_market |
| `0x07a9...32c8` | B | 100.0 | 26.6 | 374.3 | 34 | 32 | 49% | 107.8% | $+64.68 | 5 | 100.0% | 83.2% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 371.9 | 58 | 54 | 36% | 62.3% | $+118.35 | 18 | 83.3% | 25.4% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xf787...3c4e` | A | 92.6 | 21.3 | 366.3 | 23 | 23 | 55% | 85.0% | $+76.49 | 8 | 50.0% | 50.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 7
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x56ac...d77e` | C | 100.0 | 35.4 | 252.8 | 124 | 16% | 0.0% | 0 | 602 | $3453 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbc43...96d3` | C | 100.0 | 42.6 | 220.9 | 484 | 33% | 0.0% | 0 | 541 | $719 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 44.3 | 213.9 | 838 | 57% | 0.0% | 0 | 312 | $941 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 98.2 | 38.8 | 183.1 | 19 | 1% | 0.0% | 0 | 310 | $671 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 40.1 | 160.8 | 0 | 0% | 0.0% | 0 | 379 | $671 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 90.9 | 44.0 | 147.2 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 41.3 | 143.9 | 0 | 0% | 0.0% | 0 | 203 | $625 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 36.7 | 313.5 | 1290 | 1077 | 100% | 15.4% | 263 | 15.4% | 263 | $479 | existing | burst_trading,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.4 | 290.7 | 1028 | 810 | 86% | 62.6% | 19 | 60.7% | 20 | $1588 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.8 | 281.4 | 1236 | 926 | 100% | 76.0% | 234 | 76.0% | 234 | $323 | existing | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.6 | 245.1 | 566 | 527 | 93% | 89.0% | 94 | 81.7% | 103 | $1737 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.2 | 242.6 | 92 | 88 | 99% | 39.8% | 17 | 39.8% | 17 | $2734 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2a8c...fd3d` | C | 100.0 | 37.8 | 239.4 | 1251 | 530 | 97% | 135.9% | 7 | 134.9% | 8 | $158 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 22.9 | 236.2 | 325 | 206 | 86% | 13.0% | 6 | 13.0% | 10 | $2074 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.1 | 234.6 | 327 | 302 | 100% | 65.5% | 29 | 65.5% | 29 | $1606 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.2 | 233.7 | 237 | 212 | 100% | 40.9% | 27 | 40.9% | 27 | $408 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xa66c...cd67` | C | 100.0 | 39.8 | 228.9 | 531 | 451 | 63% | 10.8% | 106 | 14.1% | 150 | $5123 | existing | burst_trading,opposite_side_same_market |

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
| `0x7af7...4c89` | B | 100.0 | 26.9 | 331.1 | 20 | 19 | 58.1% | 3 | 58.1% | 3 | $19206 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 22.2 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 493.6 | opposite_side_same_market |
| `0x2929...1dd0` | B | 29.5 | $3000 | $3000 | 3 | 100.0% | +15.32 | +9.50 | +15.50 | +0.00 | 386.6 | open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.5 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 335.5 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 13
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.6 | 15m edge -1.77pp over 4 samples | 81.7% | 103 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 15m edge -2.63pp over 2 samples | 25.4% | 29 | existing | opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 34.6 | 15m edge -2.84pp over 3 samples | 100.1% | 100 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.2 | 15m edge -24.00pp over 2 samples | 40.9% | 27 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 31.5 | 1h edge -11.42pp over 4 samples | 125.8% | 15 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x07a9...32c8` | B | 100.0 | 26.6 | 1h edge -12.31pp over 2 samples | 83.2% | 7 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.8 | 1h edge -15.87pp over 1 samples | 109.5% | 54 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 29.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.9 | 1h edge -19.00pp over 1 samples | 149.7% | 31 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.1 | 1h edge -34.83pp over 1 samples | 65.5% | 29 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.0 | 1h edge -52.95pp over 1 samples | 78.8% | 24 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x7c1e...bab3` | C | 100.0 | 44.3 | 1h edge -8.50pp over 2 samples | 0.0% | 0 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_profit_all | burst_trading,open_copy_exposure |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 11.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 26.9 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 74.5 | $7000 | $7000 | 1 | 118.2% | 115 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | BOT | 61.5 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.2 | 652.4 | 420 | 357 | 37% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 74.5 | 616.9 | 186 | 181 | 13% | 130.0% | 25 | 118.2% | 115 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.5 | 594.4 | 301 | 258 | 85% | 87.5% | 54 | 99.5% | 62 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | B | 100.0 | 29.5 | 449.0 | 75 | 43 | 99% | 172.7% | 4 | 172.7% | 4 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 22.2 | 369.0 | 47 | 46 | 100% | 76.1% | 15 | 76.1% | 15 | existing,sports_tape | opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.9 | 272.7 | 939 | 847 | 73% | 28.3% | 290 | 28.3% | 373 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.9 | 238.2 | 149 | 138 | 19% | 33.0% | 55 | 37.3% | 240 | sports_tape | opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 41.4 | 177.6 | 184 | 164 | 74% | 24.6% | 60 | 19.6% | 82 | holder,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 94.8 | 21.1 | 177.4 | 23 | 18 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 26.9 | 174.7 | 20 | 19 | 95% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 75.6 | 174.2 | 151 | 148 | 12% | 28.4% | 47 | 44.1% | 292 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.3 | 133.3 | 56 | 41 | 32% | 117.0% | 1 | 61.2% | 3 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 17.4 | 112.1 | 7 | 5 | 64% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 19 | 466 | 100.8% | $+5401.28 | 83.9% | 99.5% | 63.6% | 2.68x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 2 | 20 | 476 | 100.0% | $+5508.87 | 83.8% | 96.6% | 63.6% | 2.65x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 17 | 413 | 101.1% | $+4752.18 | 84.0% | 99.5% | 63.6% | 2.89x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 16 | 406 | 98.8% | $+4615.51 | 86.5% | 102.7% | 66.5% | 2.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 5 | 17 | 416 | 98.0% | $+4723.10 | 86.3% | 99.5% | 66.5% | 2.54x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 24 | 620 | 88.5% | $+6578.77 | 82.3% | 85.3% | 51.2% | 2.24x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 7 | 16 | 440 | 98.6% | $+5001.30 | 84.1% | 96.6% | 63.6% | 2.62x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 8 | 17 | 450 | 97.9% | $+5108.89 | 84.0% | 93.7% | 63.6% | 2.59x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 25 | 649 | 86.5% | $+6729.22 | 81.4% | 83.5% | 43.0% | 2.16x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=5 smart>=70 |
| 10 | 26 | 659 | 86.2% | $+6836.81 | 81.3% | 81.1% | 43.0% | 2.15x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 28 | 671 | 85.9% | $+6842.70 | 81.7% | 77.4% | 40.0% | 2.36x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 14 | 353 | 98.9% | $+3966.41 | 87.0% | 102.7% | 66.5% | 2.80x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
