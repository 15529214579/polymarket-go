# Strategy Lab Report

**Generated:** 2026-07-15 22:53 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 211
- Candidate layers: 18 core + 20 watch + 10 sports + 8 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 12
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 18
- Aggregate closed copy trades: 446
- Aggregate copy ROI: 99.2%
- Aggregate copy PnL: $+5100.99
- Aggregate copy win rate: 83.6%
- Median wallet CopyROI: 96.6%
- Worst included wallet CopyROI: 64.6%
- Open copy cost / closed copy capital: 2.69x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.88 | 113 | 44% | 135.9% | $+502.77 | 30 | 66.7% | 46.9% |
| `0x2952...f50d` | A | 24.09 | 71 | 50% | 119.4% | $+358.14 | 24 | 87.5% | 19.5% |
| `0xe745...5681` | A | 5.93 | 57 | 30% | 188.3% | $+338.99 | 18 | 100.0% | 30.5% |
| `0x7992...1fc1` | A | 19.10 | 21 | 24% | 81.1% | $+243.30 | 25 | 76.0% | 19.6% |
| `0x6f16...5fe7` | A | 23.45 | 33 | 57% | 212.9% | $+191.57 | 9 | 77.8% | 31.0% |
| `0x84cd...7565` | A | 23.75 | 114 | 73% | 93.7% | $+187.34 | 17 | 76.5% | 41.9% |
| `0xa75b...772c` | A | 22.23 | 47 | 100% | 76.1% | $+167.40 | 15 | 73.3% | 28.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x7fcf...80ac` | A | 11.62 | 32 | 41% | 87.2% | $+104.63 | 9 | 66.7% | 38.5% |
| `0xcc6e...fa6f` | B | 26.53 | 301 | 85% | 99.5% | $+706.45 | 62 | 91.9% | 45.6% |
| `0x44c4...09cb` | B | 26.79 | 161 | 59% | 109.5% | $+645.88 | 54 | 81.5% | 47.8% |
| `0xfbe8...bb28` | B | 28.38 | 777 | 93% | 64.6% | $+407.01 | 43 | 86.0% | 48.0% |
| `0x89cf...5f47` | B | 26.19 | 56 | 44% | 106.0% | $+381.50 | 35 | 85.7% | 45.3% |
| `0x21cc...54bc` | B | 29.24 | 151 | 100% | 71.6% | $+178.90 | 25 | 100.0% | 36.5% |
| `0x18c2...529a` | B | 28.12 | 281 | 57% | 73.0% | $+167.92 | 22 | 77.3% | 40.0% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x0931...e78e` | B | 26.10 | 114 | 66% | 67.0% | $+127.35 | 19 | 68.4% | 26.8% |
| `0xec56...1f87` | B | 27.33 | 24 | 28% | 66.5% | $+112.96 | 17 | 100.0% | 22.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 562
- Aggregate copy ROI: 88.5%
- Aggregate copy PnL: $+7011.46
- Aggregate copy win rate: 86.3%
- Worst included CopyROI: 58.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.8 | 1037.8 | 20 | 14% | 352.9% | $+388.23 | 10 | 80.0% | $1079 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.8 | 644.2 | 1336 | 100% | 95.3% | $+2153.60 | 102 | 96.1% | $964 | existing,holder |
| `0x0f00...73ce` | A | 100.0 | 24.3 | 604.2 | 93 | 79% | 214.6% | $+171.67 | 5 | 100.0% | $2596 | existing |
| `0xde24...4ded` | B | 100.0 | 30.8 | 578.9 | 62 | 100% | 150.8% | $+286.59 | 9 | 100.0% | $2556 | existing,sports_tape |
| `0x9caf...94dc` | B | 100.0 | 32.2 | 562.0 | 60 | 49% | 103.7% | $+300.79 | 28 | 96.4% | $2824 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.4 | 542.7 | 11 | 18% | 129.2% | $+142.10 | 11 | 90.9% | $2412 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 508.8 | 716 | 93% | 59.2% | $+964.71 | 148 | 88.5% | $302 | existing |
| `0xc117...7410` | B | 85.3 | 25.8 | 504.8 | 31 | 79% | 149.7% | $+149.68 | 7 | 57.1% | $2173 | existing,holder |
| `0x5b1d...3721` | B | 100.0 | 29.0 | 502.1 | 31 | 100% | 145.0% | $+101.51 | 6 | 100.0% | $7979 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.1 | 488.2 | 137 | 99% | 79.0% | $+260.85 | 32 | 93.8% | $478 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 485.1 | 134 | 41% | 89.3% | $+268.06 | 23 | 87.0% | $353 | existing,holder |
| `0x578e...c3c0` | A | 100.0 | 23.1 | 467.0 | 78 | 38% | 58.9% | $+353.66 | 58 | 81.0% | $1149 | existing,holder |
| `0xa8b9...775d` | B | 100.0 | 34.2 | 459.7 | 326 | 93% | 68.6% | $+576.05 | 45 | 53.3% | $1276 | existing |
| `0xeb8b...6d8a` | B | 100.0 | 29.1 | 457.8 | 20 | 49% | 117.4% | $+82.18 | 7 | 100.0% | $2708 | existing |
| `0x07a9...32c8` | B | 100.0 | 26.8 | 440.9 | 33 | 50% | 108.6% | $+76.05 | 6 | 100.0% | $19294 | existing,holder |
| `0x7c73...6ee3` | B | 100.0 | 26.1 | 437.2 | 44 | 96% | 120.8% | $+60.39 | 5 | 100.0% | $5138 | existing |
| `0x7124...f0b5` | A | 98.5 | 22.8 | 430.7 | 47 | 78% | 100.3% | $+160.57 | 8 | 62.5% | $7844 | existing |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 413.3 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | $3736 | existing |
| `0x119b...ac3e` | B | 100.0 | 33.5 | 410.7 | 68 | 55% | 60.4% | $+181.10 | 27 | 77.8% | $8094 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.5 | 429.5 | 81 | 63 | 62% | 90.5% | $+189.96 | 11 | 72.7% | 56.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb916...f248` | A | 100.0 | 21.5 | 404.7 | 141 | 111 | 75% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.5 | 360.1 | 1025 | 812 | 85% | 62.6% | $+181.64 | 19 | 79.0% | 60.7% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.2 | 359.7 | 56 | 52 | 36% | 59.5% | $+107.11 | 17 | 82.3% | 16.0% | existing | opposite_side_same_market |
| `0xc419...a8b0` | B | 100.0 | 28.3 | 348.2 | 36 | 35 | 86% | 98.8% | $+59.28 | 5 | 80.0% | 85.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 21.7 | 347.2 | 58 | 46 | 60% | 48.6% | $+111.87 | 18 | 66.7% | 40.1% | existing | opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 8
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x56ac...d77e` | B | 100.0 | 34.9 | 254.3 | 119 | 16% | 0.0% | 0 | 578 | $3255 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbc43...96d3` | C | 100.0 | 42.6 | 221.0 | 484 | 33% | 0.0% | 0 | 543 | $718 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xe558...dcd2` | C | 100.0 | 44.8 | 204.1 | 794 | 57% | -42.0% | 1 | 351 | $944 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 43.7 | 192.1 | 623 | 43% | 0.0% | 0 | 271 | $787 | existing,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 92.6 | 39.3 | 174.7 | 19 | 1% | 0.0% | 0 | 309 | $669 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 40.1 | 160.8 | 0 | 0% | 0.0% | 0 | 379 | $671 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 90.9 | 44.0 | 147.2 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 41.3 | 143.9 | 0 | 0% | 0.0% | 0 | 203 | $625 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 36.4 | 314.0 | 1296 | 1078 | 100% | 15.5% | 255 | 15.5% | 255 | $472 | existing | burst_trading,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.8 | 282.1 | 1214 | 906 | 100% | 76.2% | 230 | 76.2% | 230 | $321 | existing,holder | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.8 | 244.7 | 565 | 526 | 93% | 89.3% | 95 | 82.1% | 104 | $1742 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x2a8c...fd3d` | C | 100.0 | 37.8 | 243.5 | 1253 | 561 | 97% | 141.7% | 8 | 140.6% | 9 | $167 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.2 | 242.6 | 92 | 88 | 99% | 39.8% | 17 | 39.8% | 17 | $2734 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 23.1 | 238.1 | 314 | 199 | 86% | 13.0% | 6 | 13.0% | 10 | $2100 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.3 | 233.3 | 321 | 296 | 100% | 65.5% | 29 | 65.5% | 29 | $1571 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.3 | 230.4 | 235 | 210 | 100% | 40.9% | 27 | 40.9% | 27 | $410 | existing | open_copy_exposure,opposite_side_same_market |
| `0xa66c...cd67` | C | 100.0 | 39.8 | 229.0 | 529 | 451 | 63% | 10.8% | 106 | 13.9% | 149 | $5135 | existing | burst_trading,opposite_side_same_market |
| `0x42db...d512` | B | 100.0 | 30.3 | 222.6 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | $1194 | existing,sports_tape | opposite_side_same_market |

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
| `0xa75b...772c` | A | 22.2 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 496.6 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.5 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 338.5 | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 30.8 | $3000 | $3000 | 11 | 72.7% | +6.35 | +0.95 | +10.45 | -3.35 | 332.3 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 12
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.5 | 22.8 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.8 | 15m edge -1.77pp over 4 samples | 82.1% | 104 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0f00...73ce` | A | 100.0 | 24.3 | 15m edge -10.00pp over 2 samples | 214.6% | 5 | existing | open_copy_exposure,opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 34.8 | 15m edge -2.84pp over 3 samples | 95.3% | 102 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.3 | 15m edge -24.00pp over 2 samples | 40.9% | 27 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.8 | 1h edge -15.87pp over 1 samples | 109.5% | 54 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 29.0 | 1h edge -17.38pp over 1 samples | 145.0% | 6 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.9 | 1h edge -19.00pp over 1 samples | 135.9% | 30 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.3 | 1h edge -34.83pp over 1 samples | 65.5% | 29 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.1 | 1h edge -52.95pp over 1 samples | 73.0% | 22 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x7c1e...bab3` | C | 100.0 | 43.7 | 1h edge -8.50pp over 2 samples | 0.0% | 0 | existing,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 11.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | B | 25.2 | $15777 | $15777 | 1 | -18.4% | 2 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 67.9 | $7000 | $7000 | 1 | 120.7% | 114 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | BOT | 73.6 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,extreme_price_heavy,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.2 | 652.2 | 420 | 356 | 37% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 67.9 | 620.5 | 159 | 154 | 11% | 131.9% | 22 | 120.7% | 114 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.5 | 594.4 | 301 | 258 | 85% | 87.5% | 54 | 99.5% | 62 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | pushed | - | B | 100.0 | 30.8 | 545.0 | 62 | 53 | 100% | 150.8% | 9 | 150.8% | 9 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 22.2 | 369.0 | 47 | 46 | 100% | 76.1% | 15 | 76.1% | 15 | existing,holder,sports_tape | opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.8 | 276.3 | 910 | 819 | 73% | 29.3% | 280 | 28.9% | 362 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.9 | 238.2 | 149 | 138 | 19% | 33.0% | 55 | 37.3% | 240 | sports_tape | opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 75.7 | 174.1 | 151 | 148 | 12% | 28.4% | 47 | 44.2% | 287 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 83.7 | 21.8 | 169.5 | 22 | 17 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 40.3 | 166.1 | 156 | 139 | 70% | 23.8% | 50 | 18.4% | 72 | sports_tape | opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.4 | 132.9 | 54 | 39 | 31% | 117.0% | 1 | 61.2% | 3 | holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 17.4 | 112.1 | 7 | 5 | 64% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 71.0 | 11.0 | 91.5 | 10 | 6 | 100% | 25.8% | 1 | 25.8% | 1 | existing,holder,sports_tape | - |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 18 | 446 | 99.2% | $+5100.99 | 83.6% | 96.6% | 64.6% | 2.69x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 2 | 19 | 456 | 98.5% | $+5208.58 | 83.6% | 93.7% | 64.6% | 2.65x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 16 | 397 | 99.1% | $+4499.55 | 83.9% | 96.6% | 64.6% | 2.86x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 16 | 428 | 97.5% | $+4804.79 | 84.1% | 96.6% | 64.6% | 2.62x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 5 | 17 | 438 | 96.7% | $+4912.38 | 84.0% | 93.7% | 64.6% | 2.59x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 15 | 388 | 97.9% | $+4366.24 | 86.1% | 99.5% | 64.6% | 2.57x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 7 | 16 | 398 | 97.0% | $+4473.83 | 85.9% | 96.6% | 64.6% | 2.53x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 23 | 583 | 86.7% | $+6230.43 | 82.3% | 81.1% | 41.3% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 14 | 379 | 97.1% | $+4203.35 | 84.4% | 96.6% | 64.6% | 2.80x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 10 | 14 | 379 | 95.5% | $+4174.67 | 86.3% | 96.6% | 64.6% | 2.54x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 11 | 15 | 389 | 94.7% | $+4282.26 | 86.1% | 93.7% | 64.6% | 2.50x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 24 | 610 | 84.7% | $+6358.91 | 81.5% | 78.6% | 40.1% | 2.15x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=5 smart>=70 |
