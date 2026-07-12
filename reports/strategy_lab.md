# Strategy Lab Report

**Generated:** 2026-07-12 23:39 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 89
- Candidate layers: 12 core + 20 watch + 10 sports + 5 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 52 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 12
- Aggregate closed copy trades: 337
- Aggregate copy ROI: 99.5%
- Aggregate copy PnL: $+4031.12
- Aggregate copy win rate: 84.3%
- Median wallet CopyROI: 103.5%
- Worst included wallet CopyROI: 64.6%
- Open copy cost / closed copy capital: 2.41x
- Params: tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 24.06 | 109 | 45% | 148.0% | $+517.89 | 28 | 71.4% | 42.2% |
| `0x2952...f50d` | A | 24.70 | 63 | 47% | 119.4% | $+358.14 | 24 | 87.5% | 19.5% |
| `0x7992...1fc1` | A | 19.40 | 21 | 25% | 81.1% | $+243.30 | 25 | 76.0% | 19.6% |
| `0x6f16...5fe7` | A | 23.45 | 33 | 57% | 212.9% | $+191.57 | 9 | 77.8% | 31.0% |
| `0xa75b...772c` | A | 22.39 | 46 | 100% | 76.1% | $+167.40 | 15 | 73.3% | 28.7% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xcc6e...fa6f` | B | 27.98 | 347 | 85% | 98.6% | $+917.11 | 77 | 92.2% | 45.6% |
| `0x44c4...09cb` | B | 26.93 | 162 | 60% | 112.5% | $+663.56 | 54 | 83.3% | 47.8% |
| `0xfbe8...bb28` | B | 28.54 | 762 | 94% | 64.6% | $+407.01 | 43 | 86.0% | 48.0% |
| `0x21cc...54bc` | B | 29.88 | 125 | 100% | 68.9% | $+151.55 | 22 | 100.0% | 40.6% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x18c2...529a` | B | 27.92 | 247 | 56% | 70.9% | $+134.71 | 18 | 72.2% | 39.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 668
- Aggregate copy ROI: 83.5%
- Aggregate copy PnL: $+8348.68
- Aggregate copy win rate: 87.6%
- Worst included CopyROI: 58.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.4 | 1038.6 | 20 | 14% | 352.9% | $+388.23 | 10 | 80.0% | $1083 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 676.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing,holder |
| `0xbb35...b62a` | B | 100.0 | 34.9 | 636.0 | 1334 | 100% | 93.2% | $+2142.39 | 104 | 96.2% | $964 | existing,holder |
| `0x0f00...73ce` | A | 100.0 | 24.7 | 603.2 | 88 | 79% | 214.6% | $+171.67 | 5 | 100.0% | $2386 | existing |
| `0xf4ec...c574` | B | 100.0 | 33.7 | 559.9 | 944 | 89% | 71.2% | $+1567.58 | 143 | 90.9% | $2100 | existing,holder |
| `0x9caf...94dc` | B | 100.0 | 32.0 | 554.0 | 63 | 50% | 101.7% | $+294.88 | 28 | 96.4% | $2660 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.5 | 539.4 | 10 | 17% | 129.2% | $+142.10 | 11 | 90.9% | $2119 | existing |
| `0xc117...7410` | B | 100.0 | 26.1 | 524.3 | 31 | 82% | 149.7% | $+149.68 | 7 | 57.1% | $2009 | existing,holder |
| `0xb36f...53d0` | B | 100.0 | 30.2 | 511.0 | 686 | 93% | 59.1% | $+952.20 | 146 | 88.4% | $309 | existing,holder |
| `0x5b1d...3721` | B | 100.0 | 29.5 | 503.0 | 30 | 100% | 145.0% | $+101.51 | 6 | 100.0% | $8079 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.3 | 127 | 98% | 79.7% | $+255.01 | 31 | 93.5% | $507 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 482.6 | 131 | 41% | 89.3% | $+268.06 | 23 | 87.0% | $340 | existing |
| `0xa8b9...775d` | B | 100.0 | 34.3 | 455.8 | 320 | 94% | 68.0% | $+564.27 | 44 | 52.3% | $1273 | existing |
| `0x07a9...32c8` | B | 100.0 | 27.3 | 440.3 | 35 | 54% | 108.6% | $+76.05 | 6 | 100.0% | $19476 | existing,holder |
| `0x7c73...6ee3` | B | 100.0 | 25.8 | 437.7 | 43 | 96% | 120.8% | $+60.39 | 5 | 100.0% | $5230 | existing |
| `0x7124...f0b5` | A | 98.4 | 23.0 | 431.4 | 45 | 78% | 100.3% | $+160.57 | 8 | 62.5% | $7672 | existing |
| `0x4572...83ff` | B | 100.0 | 34.8 | 425.6 | 280 | 100% | 65.9% | $+382.31 | 28 | 92.9% | $1543 | existing,holder |
| `0x119b...ac3e` | B | 100.0 | 32.9 | 433.2 | 62 | 55% | 64.8% | $+181.43 | 26 | 80.8% | $7162 | existing,sports_tape |
| `0x8bb0...ca79` | B | 100.0 | 32.1 | 420.4 | 66 | 34% | 87.0% | $+156.68 | 12 | 75.0% | $1557 | existing |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 413.3 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | $3736 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xbe18...3be4` | A | 100.0 | 18.7 | 430.1 | 77 | 61 | 62% | 90.5% | $+189.96 | 11 | 72.7% | 58.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.3 | 415.4 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb916...f248` | A | 100.0 | 21.8 | 405.4 | 135 | 105 | 74% | 61.5% | $+301.54 | 21 | 85.7% | 56.3% | existing,holder | opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 26.3 | 392.9 | 109 | 62 | 65% | 76.4% | $+106.94 | 14 | 78.6% | 67.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | B | 100.0 | 25.1 | 367.1 | 210 | 188 | 100% | 60.0% | $+150.01 | 23 | 60.9% | 60.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 100.0 | 33.4 | 365.3 | 1023 | 806 | 85% | 64.9% | $+181.74 | 18 | 83.3% | 62.9% | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.2 | 359.7 | 56 | 52 | 36% | 59.5% | $+107.11 | 17 | 82.3% | 16.0% | existing | opposite_side_same_market |
| `0xc419...a8b0` | B | 100.0 | 28.7 | 347.7 | 35 | 34 | 85% | 98.8% | $+59.28 | 5 | 80.0% | 85.9% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 5
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x361b...74fe` | C | 100.0 | 42.9 | 246.9 | 42 | 3% | 0.0% | 0 | 690 | $2928 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 44.5 | 206.6 | 589 | 40% | 0.0% | 0 | 269 | $641 | existing,leaderboard_profit_30d,leaderboard_profit_7d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 100.0 | 44.0 | 180.3 | 0 | 0% | 0.0% | 0 | 256 | $6525 | existing,leaderboard_profit_7d,leaderboard_volume_7d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 40.5 | 176.6 | 0 | 0% | 0.0% | 0 | 368 | $639 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 42.9 | 150.9 | 0 | 0% | 0.0% | 0 | 243 | $673 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 35.9 | 314.7 | 1298 | 1076 | 100% | 14.9% | 241 | 14.9% | 241 | $481 | existing | burst_trading,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.9 | 277.2 | 1167 | 868 | 100% | 74.4% | 222 | 74.4% | 222 | $318 | existing,holder | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.8 | 240.1 | 552 | 513 | 93% | 96.7% | 92 | 88.6% | 101 | $1732 | existing | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.3 | 238.9 | 90 | 86 | 99% | 39.8% | 17 | 39.8% | 17 | $2620 | existing | open_copy_exposure,opposite_side_same_market |
| `0xaf48...36c2` | A | 100.0 | 21.2 | 237.4 | 220 | 178 | 85% | 13.0% | 6 | 13.0% | 10 | $2769 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0562...18a4` | C | 100.0 | 35.9 | 230.8 | 297 | 269 | 99% | 83.8% | 73 | 83.8% | 73 | $3264 | existing | open_copy_exposure,opposite_side_same_market |
| `0xa66c...cd67` | C | 100.0 | 39.9 | 229.2 | 527 | 451 | 64% | 10.8% | 106 | 14.2% | 147 | $5194 | existing | burst_trading,opposite_side_same_market |
| `0x42db...d512` | B | 100.0 | 30.3 | 222.6 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | $1194 | existing,sports_tape | opposite_side_same_market |
| `0xc6f1...5729` | B | 100.0 | 31.4 | 224.3 | 187 | 180 | 89% | 14.7% | 82 | 13.1% | 90 | $1327 | existing,sports_tape | opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 27.6 | 215.5 | 66 | 62 | 99% | 36.3% | 16 | 36.3% | 16 | $1031 | existing,holder | opposite_side_same_market |

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
| `0xa75b...772c` | A | 22.4 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 493.7 | opposite_side_same_market |
| `0x119b...ac3e` | B | 32.9 | $10000 | $10000 | 8 | 75.0% | +7.50 | +1.00 | +1.37 | +16.62 | 320.5 | open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | B | 33.4 | $9973 | $9973 | 60 | 63.3% | +7.02 | +1.40 | +5.52 | +12.53 | 270.6 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | B | 31.4 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 230.4 | opposite_side_same_market |

## Live Edge Blocked Push Wallets

- Wallets: 5
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | A | 100.0 | 24.7 | 15m edge -10.00pp over 2 samples | 214.6% | 5 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | B | 100.0 | 25.1 | 15m edge -24.00pp over 2 samples | 60.0% | 23 | existing | open_copy_exposure,opposite_side_same_market |
| `0xf4ec...c574` | B | 100.0 | 33.7 | 1h edge -27.79pp over 2 samples | 71.2% | 143 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.8 | 1h edge -34.83pp over 1 samples | 65.9% | 28 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 27.6 | 1h edge -72.59pp over 1 samples | 36.3% | 16 | existing,holder | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | D | 18.3 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | B | 26.4 | $15777 | $15777 | 1 | -18.4% | 2 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.3 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 67.9 | $7000 | $7000 | 1 | 119.9% | 114 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | reject-bot | - | D | 50.7 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.2 | 652.2 | 422 | 356 | 37% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 67.9 | 606.5 | 155 | 150 | 11% | 129.9% | 21 | 119.9% | 114 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | pushed | - | B | 100.0 | 33.4 | 377.1 | 1023 | 806 | 85% | 64.9% | 18 | 62.9% | 19 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | pushed | - | B | 100.0 | 32.9 | 374.1 | 62 | 62 | 55% | 84.5% | 17 | 64.8% | 26 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 22.4 | 368.7 | 46 | 46 | 100% | 76.1% | 15 | 76.1% | 15 | existing,sports_tape | opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.5 | 261.1 | 870 | 783 | 72% | 27.9% | 266 | 27.7% | 348 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 43.1 | 237.7 | 149 | 138 | 20% | 33.0% | 55 | 37.2% | 237 | holder,sports_tape | opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.3 | 181.8 | 109 | 104 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 93.6 | 23.2 | 173.4 | 20 | 15 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 76.3 | 167.9 | 150 | 147 | 12% | 28.4% | 47 | 42.4% | 275 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | A | 83.3 | 24.9 | 165.2 | 98 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 40.2 | 160.2 | 137 | 122 | 68% | 23.8% | 44 | 17.9% | 66 | sports_tape | opposite_side_same_market |
| `0xc6f1...5729` | pushed | - | B | 100.0 | 31.4 | 153.7 | 187 | 180 | 89% | 14.7% | 82 | 13.1% | 90 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.4 | 133.1 | 53 | 38 | 32% | 117.0% | 1 | 61.2% | 3 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 16.0 | 116.3 | 7 | 5 | 70% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 12 | 337 | 99.5% | $+4031.12 | 84.3% | 103.5% | 64.6% | 2.41x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 2 | 11 | 328 | 97.0% | $+3839.55 | 84.5% | 98.6% | 64.6% | 2.38x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 10 | 288 | 99.4% | $+3429.68 | 84.7% | 103.5% | 64.6% | 2.62x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 4 | 14 | 381 | 92.4% | $+4499.18 | 84.0% | 89.9% | 56.3% | 2.15x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=70 closedROI>=0 smart>=70 |
| 5 | 15 | 387 | 95.2% | $+4409.39 | 81.7% | 81.1% | 60.0% | 3.25x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 6 | 16 | 396 | 94.8% | $+4473.44 | 81.3% | 78.6% | 60.0% | 3.30x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 15 | 391 | 91.6% | $+4563.89 | 83.6% | 81.1% | 56.3% | 2.13x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 16 | 399 | 91.1% | $+4607.98 | 83.5% | 78.6% | 55.1% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 21 | 475 | 86.7% | $+5227.05 | 80.4% | 71.2% | 55.1% | 2.91x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 22 | 500 | 84.5% | $+5350.77 | 79.6% | 71.0% | 41.2% | 2.80x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 20 | 467 | 87.1% | $+5182.96 | 80.5% | 71.6% | 56.3% | 2.84x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 12 | 21 | 492 | 84.9% | $+5306.68 | 79.7% | 71.2% | 41.2% | 2.73x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
