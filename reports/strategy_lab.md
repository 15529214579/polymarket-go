# Strategy Lab Report

**Generated:** 2026-07-31 23:44 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 191
- Candidate layers: 15 core + 20 watch + 10 sports + 7 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 56 total
- Live-edge blocked push wallets: 7
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 4 observation-only

## Selected Core Strategy

- Wallets: 15
- Aggregate closed copy trades: 291
- Aggregate copy ROI: 126.6%
- Aggregate copy PnL: $+4544.22
- Aggregate copy win rate: 85.6%
- Median wallet CopyROI: 135.5%
- Worst included wallet CopyROI: 102.8%
- Open copy cost / closed copy capital: 2.55x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.31 | 120 | 39% | 142.2% | $+583.20 | 34 | 67.7% | 58.2% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x2952...f50d` | A | 24.23 | 77 | 56% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0xfc92...5667` | A | 23.53 | 143 | 92% | 147.5% | $+265.52 | 11 | 90.9% | 54.8% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 57 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 60 | 90% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xa75b...772c` | B | 25.09 | 156 | 99% | 113.4% | $+623.71 | 42 | 83.3% | 34.0% |
| `0xde24...4ded` | B | 28.46 | 208 | 100% | 102.8% | $+595.93 | 31 | 100.0% | 31.3% |
| `0x44c4...09cb` | B | 27.08 | 139 | 60% | 109.5% | $+536.49 | 45 | 77.8% | 47.6% |
| `0xe916...7e93` | B | 27.86 | 62 | 98% | 117.3% | $+152.51 | 12 | 83.3% | 44.7% |
| `0xffa1...6340` | B | 25.47 | 52 | 98% | 165.2% | $+148.73 | 9 | 88.9% | 18.0% |
| `0xeb8b...6d8a` | B | 28.11 | 27 | 51% | 138.5% | $+124.64 | 9 | 100.0% | 42.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 726
- Aggregate copy ROI: 104.3%
- Aggregate copy PnL: $+9772.63
- Aggregate copy win rate: 89.5%
- Worst included CopyROI: 59.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.3 | 897.8 | 25 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1207 | existing |
| `0xc4f7...a5e1` | B | 100.0 | 34.2 | 762.7 | 82 | 96% | 167.2% | $+434.81 | 26 | 92.3% | $899 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 738.7 | 39 | 95% | 248.6% | $+223.76 | 7 | 100.0% | $7773 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.2 | 719.3 | 1340 | 100% | 111.0% | $+2552.86 | 103 | 96.1% | $1063 | existing,holder |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0x89cf...5f47` | B | 100.0 | 33.3 | 653.0 | 70 | 40% | 119.7% | $+526.61 | 43 | 88.4% | $674 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xf4e1...34a5` | B | 100.0 | 30.4 | 617.0 | 67 | 40% | 125.0% | $+374.96 | 25 | 84.0% | $5956 | existing |
| `0x2a35...9015` | B | 100.0 | 33.4 | 610.2 | 148 | 99% | 108.6% | $+521.40 | 36 | 91.7% | $1961 | existing,holder |
| `0xcefd...d6aa` | B | 100.0 | 32.9 | 606.8 | 75 | 87% | 169.7% | $+254.53 | 10 | 100.0% | $255 | existing,holder |
| `0xcc6e...fa6f` | B | 100.0 | 26.0 | 611.6 | 265 | 82% | 97.7% | $+605.53 | 54 | 88.9% | $2153 | existing,sports_tape |
| `0x17fe...b0ca` | B | 100.0 | 30.7 | 568.0 | 179 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $663 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 63 | 51% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.7 | 559.9 | 27 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $754 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 546.1 | 12 | 19% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 525.7 | 899 | 93% | 59.8% | $+1220.55 | 185 | 89.2% | $296 | existing |
| `0x7992...1fc1` | A | 100.0 | 18.6 | 524.5 | 21 | 21% | 86.7% | $+294.71 | 29 | 79.3% | $2688 | existing |
| `0x18c2...529a` | B | 100.0 | 28.3 | 488.7 | 764 | 73% | 71.6% | $+508.26 | 55 | 85.5% | $275 | existing,holder |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.1 | 621.0 | 354 | 256 | 39% | 105.2% | $+609.96 | 55 | 96.4% | 130.0% | existing | open_copy_exposure,opposite_side_same_market |
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7a26...1589` | B | 100.0 | 31.9 | 530.1 | 124 | 114 | 64% | 91.7% | $+394.37 | 34 | 94.1% | 73.9% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x84cd...7565` | A | 100.0 | 23.5 | 512.7 | 111 | 56 | 71% | 116.7% | $+186.71 | 13 | 92.3% | 93.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7e01...b0b5` | B | 100.0 | 30.1 | 470.2 | 368 | 206 | 97% | 66.0% | $+527.78 | 54 | 81.5% | 66.0% | holder | opposite_side_same_market |
| `0xc117...7410` | B | 84.8 | 25.4 | 464.8 | 33 | 27 | 77% | 143.0% | $+143.01 | 7 | 42.9% | 118.6% | existing | opposite_side_same_market |
| `0xa8b9...775d` | B | 100.0 | 33.1 | 446.5 | 385 | 259 | 93% | 68.0% | $+571.20 | 44 | 54.5% | 63.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 34.1 | 441.6 | 308 | 177 | 93% | 131.6% | $+157.95 | 7 | 71.4% | 121.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.1 | 431.6 | 80 | 80 | 59% | 76.9% | $+184.49 | 22 | 86.4% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 7
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x5f65...2036` | C | 100.0 | 37.3 | 330.9 | 34 | 2% | 78.9% | 276 | 1375 | $10123 | existing,leaderboard_profit_7d,leaderboard_volume_7d | opposite_side_same_market |
| `0x16bb...8492` | B | 100.0 | 31.9 | 254.1 | 16 | 1% | 0.0% | 0 | 690 | $3645 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x27f7...44b0` | C | 100.0 | 43.7 | 230.4 | 378 | 29% | 0.0% | 0 | 687 | $715 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x56ac...d77e` | C | 100.0 | 40.3 | 224.4 | 39 | 5% | 0.0% | 0 | 485 | $2161 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 100.0 | 38.9 | 169.9 | 19 | 1% | 0.0% | 0 | 309 | $674 | existing,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 44.9 | 163.4 | 0 | 0% | 0.0% | 0 | 282 | $746 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x1b20...3ab0` | C | 89.8 | 43.2 | 135.1 | 9 | 1% | 0.0% | 0 | 193 | $1205 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xe872...819a` | B | 100.0 | 25.4 | 300.4 | 1013 | 795 | 85% | 76.6% | 15 | 74.0% | 16 | $1627 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2a99...51bb` | A | 100.0 | 24.2 | 261.7 | 698 | 491 | 99% | 42.1% | 125 | 42.1% | 125 | $295 | existing | open_copy_exposure,opposite_side_same_market |
| `0x8dd1...5a9b` | C | 100.0 | 39.3 | 250.9 | 660 | 629 | 86% | 56.5% | 147 | 57.1% | 151 | $803 | existing | open_copy_exposure,opposite_side_same_market |
| `0xfbe8...bb28` | C | 100.0 | 35.3 | 248.2 | 1119 | 588 | 92% | 45.1% | 50 | 46.6% | 59 | $165 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 41.9 | 247.8 | 564 | 513 | 100% | 66.5% | 52 | 66.5% | 52 | $1623 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xd3d3...bb8e` | C | 100.0 | 39.4 | 243.1 | 864 | 622 | 71% | 14.0% | 182 | 12.9% | 257 | $632 | existing,holder | opposite_side_same_market |
| `0x6e54...511e` | B | 100.0 | 31.9 | 237.1 | 979 | 472 | 88% | 12.2% | 22 | 11.1% | 23 | $140 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.0 | 236.8 | 93 | 89 | 93% | 40.5% | 18 | 41.6% | 19 | $2728 | existing | open_copy_exposure,opposite_side_same_market |
| `0x1e8e...f9b2` | C | 100.0 | 42.6 | 232.5 | 705 | 540 | 98% | 99.4% | 164 | 99.7% | 168 | $243 | existing | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 228.3 | 73 | 69 | 95% | 58.4% | 20 | 58.4% | 20 | $3677 | existing | opposite_side_same_market |

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
| `0x7af7...4c89` | B | 100.0 | 30.5 | 316.8 | 21 | 20 | 58.1% | 3 | 58.1% | 3 | $12821 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | B | 25.1 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 555.6 | opposite_side_same_market |
| `0x2929...1dd0` | C | 42.7 | $3000 | $3000 | 4 | 100.0% | +16.72 | +9.50 | +15.50 | +20.95 | 405.7 | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | B | 26.0 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 324.0 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 7
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0xbb35...b62a` | B | 100.0 | 34.2 | 15m edge -2.84pp over 3 samples | 111.0% | 103 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 28.5 | 1h edge -11.42pp over 4 samples | 102.8% | 31 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.1 | 1h edge -15.87pp over 1 samples | 109.5% | 45 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.3 | 1h edge -19.00pp over 1 samples | 142.2% | 34 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 41.9 | 1h edge -34.83pp over 1 samples | 66.5% | 52 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.3 | 1h edge -52.95pp over 1 samples | 71.6% | 55 | existing,holder | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 30.5 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 72.9 | $7000 | $7000 | 1 | 114.4% | 122 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 47.7 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 72.9 | 673.7 | 518 | 326 | 36% | 121.3% | 45 | 114.4% | 122 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | B | 100.0 | 25.1 | 672.3 | 156 | 131 | 99% | 113.4% | 42 | 113.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x2929...1dd0` | watch | - | C | 100.0 | 42.7 | 556.3 | 344 | 254 | 99% | 122.0% | 18 | 122.0% | 18 | holder,sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 26.0 | 548.7 | 265 | 222 | 82% | 82.9% | 44 | 97.7% | 54 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 63.4 | 531.9 | 348 | 311 | 32% | 129.7% | 15 | 114.3% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 70.1 | 294.9 | 1007 | 905 | 74% | 30.4% | 304 | 29.8% | 390 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 41.4 | 261.5 | 163 | 151 | 15% | 37.5% | 59 | 36.2% | 278 | sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject | - | D | 100.0 | 57.5 | 236.4 | 130 | 127 | 10% | 33.6% | 41 | 50.3% | 326 | sports_tape | bot_like_flow,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | D | 100.0 | 57.9 | 235.4 | 331 | 302 | 81% | 32.4% | 110 | 29.3% | 137 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 100.0 | 30.0 | 182.1 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 83.3 | 20.0 | 171.2 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 30.5 | 160.9 | 21 | 20 | 66% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 74.9 | 10.0 | 96.6 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 24.2 | 94.7 | 20 | 5 | 32% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 15 | 291 | 126.6% | $+4544.22 | 85.6% | 135.5% | 102.8% | 2.55x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 26 | 577 | 103.9% | $+7252.85 | 85.3% | 108.7% | 65.3% | 3.67x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 14 | 257 | 124.6% | $+3961.02 | 87.9% | 135.4% | 102.8% | 2.37x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 4 | 25 | 566 | 104.8% | $+7145.42 | 85.5% | 108.9% | 65.3% | 3.71x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 5 | 13 | 260 | 127.1% | $+4068.70 | 85.4% | 135.5% | 102.8% | 2.69x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 13 | 273 | 125.2% | $+4270.85 | 85.0% | 135.2% | 102.8% | 2.56x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 13 | 212 | 127.3% | $+3424.53 | 90.1% | 135.5% | 102.8% | 2.38x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 8 | 25 | 543 | 101.5% | $+6669.65 | 86.4% | 108.4% | 65.3% | 3.65x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 12 | 251 | 126.8% | $+3944.06 | 84.9% | 135.4% | 102.8% | 2.69x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 10 | 24 | 532 | 102.4% | $+6562.22 | 86.7% | 108.7% | 65.3% | 3.70x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 11 | 23 | 548 | 103.5% | $+6872.05 | 85.2% | 108.4% | 65.3% | 3.75x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 12 | 24 | 559 | 102.6% | $+6979.48 | 85.0% | 105.6% | 65.3% | 3.70x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
