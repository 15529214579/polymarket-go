# Strategy Lab Report

**Generated:** 2026-07-14 23:34 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 99
- Candidate layers: 17 core + 20 watch + 10 sports + 8 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 62 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 5 observation-only

## Selected Core Strategy

- Wallets: 17
- Aggregate closed copy trades: 422
- Aggregate copy ROI: 96.6%
- Aggregate copy PnL: $+4772.16
- Aggregate copy win rate: 83.4%
- Median wallet CopyROI: 98.7%
- Worst included wallet CopyROI: 64.6%
- Open copy cost / closed copy capital: 3.05x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.64 | 110 | 43% | 142.2% | $+511.76 | 29 | 69.0% | 46.9% |
| `0x2952...f50d` | A | 24.56 | 66 | 49% | 119.4% | $+358.14 | 24 | 87.5% | 19.5% |
| `0x7992...1fc1` | A | 19.40 | 21 | 25% | 81.1% | $+243.30 | 25 | 76.0% | 19.6% |
| `0x6f16...5fe7` | A | 23.45 | 33 | 57% | 212.9% | $+191.57 | 9 | 77.8% | 31.0% |
| `0x84cd...7565` | A | 23.74 | 113 | 75% | 98.7% | $+187.44 | 16 | 81.2% | 41.9% |
| `0xa75b...772c` | A | 22.39 | 46 | 100% | 76.1% | $+167.40 | 15 | 73.3% | 28.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xcc6e...fa6f` | B | 27.06 | 307 | 85% | 100.8% | $+745.53 | 64 | 92.2% | 45.6% |
| `0x44c4...09cb` | B | 26.92 | 160 | 59% | 112.5% | $+663.56 | 54 | 83.3% | 47.8% |
| `0xfbe8...bb28` | B | 28.45 | 766 | 93% | 64.6% | $+407.01 | 43 | 86.0% | 48.0% |
| `0x89cf...5f47` | B | 26.24 | 55 | 44% | 106.0% | $+381.50 | 35 | 85.7% | 45.3% |
| `0x21cc...54bc` | B | 29.57 | 140 | 100% | 69.1% | $+158.90 | 23 | 100.0% | 36.0% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x18c2...529a` | B | 27.98 | 258 | 56% | 68.0% | $+135.95 | 19 | 73.7% | 38.7% |
| `0x0931...e78e` | B | 26.17 | 113 | 66% | 67.0% | $+127.35 | 19 | 68.4% | 26.8% |
| `0xec56...1f87` | B | 27.33 | 24 | 28% | 66.5% | $+112.96 | 17 | 100.0% | 22.6% |
| `0x8e4e...1dfc` | B | 27.70 | 97 | 21% | 72.1% | $+100.91 | 8 | 62.5% | 48.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 520
- Aggregate copy ROI: 87.5%
- Aggregate copy PnL: $+6780.25
- Aggregate copy win rate: 86.5%
- Worst included CopyROI: 56.3%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.2 | 1038.8 | 20 | 14% | 352.9% | $+388.23 | 10 | 80.0% | $1082 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.9 | 641.1 | 1337 | 100% | 94.2% | $+2176.97 | 104 | 96.2% | $960 | existing,holder |
| `0x0f00...73ce` | A | 100.0 | 24.3 | 607.2 | 93 | 79% | 214.6% | $+171.67 | 5 | 100.0% | $2596 | existing,holder |
| `0x73e2...46ec` | A | 100.0 | 5.4 | 545.7 | 11 | 18% | 129.2% | $+142.10 | 11 | 90.9% | $2412 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 31.9 | 542.0 | 60 | 49% | 99.3% | $+278.15 | 27 | 96.3% | $2666 | existing |
| `0xc117...7410` | B | 100.0 | 25.8 | 521.0 | 31 | 79% | 149.7% | $+149.68 | 7 | 57.1% | $2173 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 508.5 | 700 | 93% | 59.2% | $+958.50 | 147 | 88.4% | $306 | existing |
| `0x5b1d...3721` | B | 100.0 | 29.5 | 503.0 | 30 | 100% | 145.0% | $+101.51 | 6 | 100.0% | $8079 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.3 | 487.5 | 133 | 99% | 79.7% | $+255.01 | 31 | 93.5% | $488 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 482.6 | 131 | 41% | 89.3% | $+268.06 | 23 | 87.0% | $340 | existing |
| `0xde24...4ded` | B | 100.0 | 28.8 | 459.9 | 38 | 100% | 127.8% | $+153.36 | 5 | 100.0% | $2155 | existing,holder |
| `0xa8b9...775d` | B | 100.0 | 34.3 | 458.3 | 325 | 94% | 68.0% | $+564.27 | 44 | 52.3% | $1267 | existing,holder |
| `0x7c73...6ee3` | B | 100.0 | 25.8 | 437.7 | 43 | 96% | 120.8% | $+60.39 | 5 | 100.0% | $5230 | existing |
| `0x07a9...32c8` | B | 100.0 | 27.1 | 437.4 | 34 | 51% | 108.6% | $+76.05 | 6 | 100.0% | $19036 | existing |
| `0xeb8b...6d8a` | B | 100.0 | 28.4 | 435.3 | 19 | 48% | 113.0% | $+67.80 | 6 | 100.0% | $1729 | existing,holder |
| `0x7124...f0b5` | A | 98.4 | 23.0 | 431.4 | 45 | 78% | 100.3% | $+160.57 | 8 | 62.5% | $7672 | existing |
| `0x119b...ac3e` | B | 100.0 | 32.9 | 433.2 | 62 | 55% | 64.8% | $+181.43 | 26 | 80.8% | $7162 | existing,sports_tape |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 413.3 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | $3736 | existing |
| `0xb916...f248` | A | 100.0 | 21.7 | 406.7 | 137 | 74% | 56.3% | $+292.83 | 24 | 79.2% | $618 | existing,holder |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x552b...f8bc` | B | 100.0 | 34.3 | 675.8 | 66 | 28 | 63% | 202.0% | $+242.42 | 10 | 70.0% | 215.1% | existing | extreme_price_heavy,opposite_side_same_market |
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.6 | 430.1 | 77 | 61 | 61% | 90.5% | $+189.96 | 11 | 72.7% | 58.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.3 | 415.4 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.4 | 365.2 | 1021 | 809 | 85% | 64.9% | $+181.74 | 18 | 83.3% | 62.9% | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.2 | 359.7 | 56 | 52 | 36% | 59.5% | $+107.11 | 17 | 82.3% | 16.0% | existing | opposite_side_same_market |
| `0xc419...a8b0` | B | 100.0 | 28.7 | 347.7 | 35 | 34 | 85% | 98.8% | $+59.28 | 5 | 80.0% | 85.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 33.1 | 346.1 | 28 | 24 | 90% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 21.9 | 335.6 | 56 | 44 | 60% | 46.1% | $+101.51 | 17 | 64.7% | 38.1% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 8
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbc43...96d3` | C | 100.0 | 40.4 | 221.2 | 191 | 13% | 0.0% | 0 | 505 | $663 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x361b...74fe` | C | 100.0 | 42.5 | 220.1 | 32 | 2% | 0.0% | 0 | 699 | $2946 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xb61b...8a06` | C | 100.0 | 40.7 | 194.9 | 147 | 10% | 0.0% | 0 | 227 | $630 | existing,holder,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 44.0 | 192.2 | 622 | 42% | 0.0% | 0 | 259 | $679 | existing,holder,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0x1465...7072` | C | 100.0 | 38.8 | 181.7 | 0 | 0% | 0.0% | 0 | 570 | $1076 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 40.1 | 160.8 | 0 | 0% | 0.0% | 0 | 379 | $671 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 90.9 | 44.0 | 147.2 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 41.3 | 143.9 | 0 | 0% | 0.0% | 0 | 203 | $625 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 36.1 | 325.5 | 1297 | 1078 | 100% | 14.8% | 246 | 14.8% | 246 | $478 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.8 | 279.6 | 1189 | 885 | 100% | 75.4% | 225 | 75.4% | 225 | $320 | existing,holder | opposite_side_same_market |
| `0x2a8c...fd3d` | C | 100.0 | 37.8 | 245.1 | 1257 | 573 | 97% | 141.7% | 8 | 140.6% | 9 | $170 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.6 | 239.7 | 547 | 508 | 93% | 97.0% | 89 | 88.7% | 98 | $1732 | existing | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.2 | 239.6 | 92 | 88 | 99% | 39.8% | 17 | 39.8% | 17 | $2734 | existing | open_copy_exposure,opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 23.8 | 234.0 | 295 | 189 | 88% | 13.0% | 6 | 13.0% | 10 | $2163 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.6 | 232.5 | 228 | 206 | 100% | 40.9% | 27 | 40.9% | 27 | $418 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.2 | 232.0 | 308 | 283 | 100% | 65.5% | 29 | 65.5% | 29 | $1583 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0562...18a4` | C | 100.0 | 35.9 | 230.3 | 292 | 265 | 99% | 84.0% | 72 | 84.0% | 72 | $3299 | existing | open_copy_exposure,opposite_side_same_market |
| `0xa66c...cd67` | C | 100.0 | 39.9 | 229.1 | 527 | 451 | 64% | 10.8% | 106 | 14.2% | 147 | $5188 | existing | burst_trading,opposite_side_same_market |

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

- Wallets: 5
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 22.4 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 493.7 | opposite_side_same_market |
| `0x8d85...c436` | C | 36.1 | $1822 | $1822 | 4 | 100.0% | +27.09 | +24.50 | +27.95 | +27.95 | 388.7 | burst_trading,opposite_side_same_market |
| `0x119b...ac3e` | B | 32.9 | $10000 | $10000 | 8 | 75.0% | +7.50 | +1.00 | +1.37 | +16.62 | 320.5 | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 33.4 | $9973 | $9973 | 63 | 60.3% | +6.67 | +1.28 | +5.14 | +12.53 | 266.2 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 3
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | A | 100.0 | 24.3 | 15m edge -10.00pp over 2 samples | 214.6% | 5 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | A | 100.0 | 24.6 | 15m edge -24.00pp over 2 samples | 40.9% | 27 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.2 | 1h edge -34.83pp over 1 samples | 65.5% | 29 | existing,holder | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 13.8 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | B | 25.8 | $15777 | $15777 | 1 | -18.4% | 2 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.3 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 68.0 | $7000 | $7000 | 1 | 120.1% | 114 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 59.6 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.2 | 652.0 | 419 | 355 | 37% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 68.0 | 606.7 | 155 | 150 | 11% | 129.9% | 21 | 120.1% | 114 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | pushed | - | B | 100.0 | 33.4 | 377.3 | 1021 | 809 | 85% | 64.9% | 18 | 62.9% | 19 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | pushed | - | B | 100.0 | 32.9 | 374.1 | 62 | 62 | 55% | 84.5% | 17 | 64.8% | 26 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 22.4 | 368.7 | 46 | 46 | 100% | 76.1% | 15 | 76.1% | 15 | existing,sports_tape | opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.6 | 262.5 | 880 | 789 | 72% | 28.0% | 268 | 27.8% | 350 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x8d85...c436` | pushed | - | C | 100.0 | 36.1 | 245.4 | 1297 | 1078 | 100% | 14.8% | 246 | 14.8% | 246 | existing,holder,sports_tape | burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 42.9 | 238.2 | 149 | 138 | 19% | 33.0% | 55 | 37.3% | 240 | sports_tape | opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 93.6 | 23.2 | 173.4 | 20 | 15 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 75.9 | 171.6 | 151 | 148 | 12% | 28.4% | 47 | 43.4% | 286 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 39.8 | 154.3 | 144 | 128 | 69% | 22.1% | 46 | 16.9% | 68 | holder,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.4 | 132.9 | 53 | 38 | 32% | 117.0% | 1 | 61.2% | 3 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 17.4 | 112.1 | 7 | 5 | 64% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 17 | 422 | 96.6% | $+4772.16 | 83.4% | 98.7% | 64.6% | 3.05x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 15 | 405 | 95.1% | $+4479.68 | 84.0% | 98.7% | 64.6% | 2.39x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 15 | 373 | 96.1% | $+4170.72 | 83.6% | 98.7% | 64.6% | 3.31x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 4 | 14 | 366 | 94.9% | $+4032.14 | 85.8% | 99.7% | 64.6% | 2.30x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 13 | 356 | 94.4% | $+3878.24 | 84.3% | 98.7% | 64.6% | 2.56x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 13 | 357 | 92.3% | $+3840.57 | 86.0% | 98.7% | 64.6% | 2.26x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 7 | 20 | 483 | 89.4% | $+5416.97 | 82.6% | 78.6% | 56.3% | 2.70x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 16 | 410 | 88.8% | $+4500.20 | 85.4% | 89.9% | 56.3% | 2.07x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 12 | 317 | 94.0% | $+3430.70 | 86.4% | 99.7% | 64.6% | 2.48x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 10 | 22 | 516 | 87.0% | $+5577.77 | 81.6% | 74.1% | 40.0% | 2.61x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 23 | 524 | 86.6% | $+5621.86 | 81.5% | 72.1% | 40.0% | 2.68x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 10 | 298 | 95.8% | $+3293.92 | 88.3% | 103.4% | 64.6% | 2.26x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
