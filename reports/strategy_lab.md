# Strategy Lab Report

**Generated:** 2026-08-04 23:31 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (378)
- Valid strategies found: 167
- Candidate layers: 11 core + 20 watch + 10 sports + 2 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 45 total
- Live-edge blocked push wallets: 9
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 233
- Aggregate copy ROI: 127.3%
- Aggregate copy PnL: $+3500.11
- Aggregate copy win rate: 80.7%
- Median wallet CopyROI: 135.2%
- Worst included wallet CopyROI: 108.4%
- Open copy cost / closed copy capital: 2.25x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 23.87 | 229 | 100% | 116.2% | $+847.88 | 55 | 83.6% | 36.6% |
| `0xb2ed...4418` | A | 23.14 | 125 | 38% | 134.4% | $+591.26 | 36 | 66.7% | 58.6% |
| `0x2952...f50d` | A | 23.78 | 77 | 59% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 27.12 | 137 | 63% | 108.6% | $+521.13 | 44 | 77.3% | 47.6% |
| `0x5b1d...3721` | B | 27.04 | 47 | 96% | 196.7% | $+236.02 | 9 | 88.9% | 41.9% |
| `0x0ec9...1e0c` | B | 25.93 | 13 | 22% | 146.1% | $+146.14 | 10 | 70.0% | 11.1% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 660
- Aggregate copy ROI: 86.9%
- Aggregate copy PnL: $+7587.30
- Aggregate copy win rate: 86.7%
- Worst included CopyROI: 44.5%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xbb35...b62a` | B | 100.0 | 33.7 | 755.6 | 1340 | 100% | 118.9% | $+2723.78 | 104 | 97.1% | $1106 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 29 | 48% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0x17fe...b0ca` | B | 100.0 | 31.0 | 568.0 | 175 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $674 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.7 | 559.9 | 27 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $754 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.4 | 528.3 | 911 | 93% | 59.6% | $+1228.71 | 187 | 89.3% | $295 | existing,holder |
| `0x092b...614e` | B | 100.0 | 26.0 | 493.7 | 148 | 97% | 97.6% | $+185.38 | 18 | 83.3% | $1081 | existing |
| `0x819d...6e9c` | B | 100.0 | 31.1 | 488.9 | 23 | 33% | 134.1% | $+107.24 | 8 | 100.0% | $421 | existing |
| `0x18c2...529a` | B | 100.0 | 28.0 | 483.3 | 842 | 74% | 70.3% | $+520.04 | 57 | 86.0% | $267 | existing |
| `0x84cd...7565` | A | 100.0 | 23.4 | 475.3 | 124 | 79% | 93.9% | $+178.48 | 16 | 75.0% | $445 | existing |
| `0x2d44...4bae` | A | 100.0 | 14.9 | 462.8 | 15 | 35% | 117.9% | $+94.33 | 7 | 71.4% | $946 | existing |
| `0x578e...c3c0` | A | 100.0 | 24.1 | 451.7 | 57 | 34% | 59.8% | $+286.98 | 46 | 80.4% | $1281 | existing |
| `0x7124...f0b5` | A | 98.7 | 22.5 | 429.9 | 49 | 79% | 100.3% | $+160.57 | 8 | 62.5% | $7903 | existing |
| `0x7fcf...80ac` | A | 100.0 | 12.4 | 417.1 | 26 | 37% | 81.2% | $+105.62 | 10 | 70.0% | $2172 | existing |
| `0x7673...fa40` | B | 100.0 | 28.7 | 405.0 | 93 | 62% | 76.6% | $+122.51 | 14 | 78.6% | $545 | existing |
| `0xabb2...0bc2` | B | 100.0 | 31.3 | 402.2 | 108 | 53% | 55.9% | $+201.35 | 35 | 85.7% | $260 | existing |
| `0x760f...326a` | B | 100.0 | 32.2 | 389.7 | 13 | 15% | 83.2% | $+108.16 | 11 | 72.7% | $4094 | existing |
| `0x0e24...7014` | B | 100.0 | 29.0 | 388.7 | 144 | 98% | 44.5% | $+293.77 | 41 | 75.6% | $1049 | existing,holder |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 32.9 | 621.0 | 386 | 274 | 40% | 104.9% | $+618.86 | 56 | 96.4% | 127.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 33.4 | 517.3 | 338 | 206 | 94% | 149.4% | $+209.20 | 9 | 77.8% | 139.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 405.6 | 128 | 72 | 28% | 75.4% | $+150.82 | 18 | 72.2% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x0931...e78e` | B | 100.0 | 25.6 | 379.9 | 125 | 71 | 68% | 71.2% | $+106.79 | 15 | 73.3% | 55.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x68cb...6f57` | A | 100.0 | 21.4 | 375.0 | 72 | 58 | 67% | 54.0% | $+156.56 | 22 | 63.6% | 45.6% | existing | opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 371.9 | 61 | 55 | 38% | 62.3% | $+118.35 | 18 | 83.3% | 25.4% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x272d...bc2f` | B | 100.0 | 28.6 | 359.1 | 437 | 227 | 42% | 64.5% | $+116.01 | 16 | 81.2% | 61.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 346.9 | 31 | 26 | 91% | 84.4% | $+92.81 | 5 | 100.0% | 84.4% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 2
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3a6...71bd` | C | 100.0 | 37.0 | 218.9 | 405 | 31% | 0.0% | 0 | 384 | $764 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x2dc1...b33c` | C | 100.0 | 44.2 | 173.2 | 0 | 0% | 123.1% | 1 | 365 | $6827 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure,opposite_side_same_market |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x2929...1dd0` | C | 100.0 | 43.7 | 252.0 | 680 | 536 | 99% | 109.2% | 24 | 107.7% | 25 | $1070 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x2fbb...7e44` | A | 100.0 | 24.9 | 223.3 | 460 | 211 | 98% | 19.7% | 70 | 19.7% | 70 | $198 | existing | opposite_side_same_market |
| `0x4bba...cf14` | B | 100.0 | 26.4 | 217.4 | 768 | 213 | 100% | 20.0% | 6 | 20.0% | 6 | $103 | existing | open_copy_exposure,opposite_side_same_market |
| `0x9ffe...86b5` | B | 100.0 | 25.2 | 215.6 | 54 | 50 | 84% | 31.1% | 14 | 32.5% | 16 | $3259 | existing | opposite_side_same_market |
| `0xd052...563d` | C | 100.0 | 35.3 | 207.7 | 292 | 254 | 72% | 17.5% | 73 | 22.9% | 90 | $1129 | existing | open_copy_exposure,opposite_side_same_market |
| `0x682c...ef07` | B | 100.0 | 33.1 | 204.7 | 401 | 218 | 87% | 45.3% | 14 | 36.2% | 18 | $155 | existing,holder | open_copy_exposure,opposite_side_same_market |
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

- Wallets: 1
- Rule: recent basketball/soccer/esports large-order wallets; 5k+ direct whales or scored low-bot tape candidates; pushed through tape list with consensus gate

| Wallet | Tier | Smart | Bot | TapeHotScore | TargetT | TargetLarge | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x1b1f...d6d5` | C | 100.0 | 43.5 | 310.3 | 72 | 67 | 46.9% | 27 | 45.8% | 28 | $1279 | existing,holder,sports_tape | extreme_price_heavy,opposite_side_same_market |

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
| `0xa75b...772c` | A | 23.9 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 571.3 | opposite_side_same_market |
| `0x2929...1dd0` | C | 43.7 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 404.0 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 9
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x7124...f0b5` | A | 98.7 | 22.5 | 15m edge -1.22pp over 2 samples | 100.3% | 8 | existing | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.8 | 15m edge -2.63pp over 2 samples | 25.4% | 29 | existing | opposite_side_same_market |
| `0xbb35...b62a` | B | 100.0 | 33.7 | 15m edge -2.84pp over 3 samples | 118.9% | 104 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.1 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 27.0 | 1h edge -17.38pp over 1 samples | 196.7% | 9 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.1 | 1h edge -19.00pp over 1 samples | 134.4% | 36 | existing | open_copy_exposure,opposite_side_same_market |
| `0x119e...cb14` | B | 100.0 | 31.5 | 1h edge -36.88pp over 1 samples | 84.4% | 5 | existing | opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.0 | 1h edge -52.95pp over 1 samples | 70.3% | 57 | existing | open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 29.0 | 1h edge -72.59pp over 1 samples | 44.5% | 41 | existing,holder | opposite_side_same_market |

## Sports Tape Candidate Review

- Wallets: 3
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa75b...772c` | pushed | - | A | 100.0 | 23.9 | 735.8 | 229 | 179 | 100% | 116.2% | 55 | 116.2% | 55 | existing,holder,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | pushed | - | C | 100.0 | 43.7 | 565.9 | 680 | 536 | 99% | 109.2% | 24 | 107.7% | 25 | existing,holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x1b1f...d6d5` | pushed | - | C | 100.0 | 43.5 | 253.3 | 72 | 67 | 92% | 46.9% | 27 | 45.8% | 28 | existing,holder,sports_tape | extreme_price_heavy,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 233 | 127.3% | $+3500.11 | 80.7% | 135.2% | 108.4% | 2.25x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 20 | 454 | 101.9% | $+5389.85 | 81.9% | 108.5% | 61.0% | 3.21x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 10 | 197 | 125.9% | $+2908.85 | 83.2% | 135.4% | 108.4% | 1.95x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 10 | 224 | 124.1% | $+3264.09 | 80.4% | 134.8% | 108.4% | 2.26x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 19 | 443 | 103.0% | $+5282.42 | 82.2% | 108.6% | 61.0% | 3.26x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 6 | 19 | 445 | 99.7% | $+5153.83 | 81.8% | 108.4% | 61.0% | 3.24x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 19 | 418 | 98.9% | $+4798.59 | 83.3% | 108.4% | 61.0% | 3.15x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 8 | 18 | 434 | 100.7% | $+5046.40 | 82.0% | 108.5% | 61.0% | 3.29x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 9 | 22 | 551 | 92.2% | $+5985.46 | 80.8% | 103.0% | 44.5% | 2.81x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=10 smart>=70 |
| 10 | 18 | 407 | 100.0% | $+4691.16 | 83.5% | 108.5% | 61.0% | 3.20x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 11 | 17 | 411 | 101.7% | $+4809.49 | 82.2% | 108.4% | 61.0% | 3.43x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 12 | 23 | 560 | 91.2% | $+6067.56 | 80.5% | 97.6% | 44.5% | 2.75x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
