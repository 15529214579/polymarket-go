# Strategy Lab Report

**Generated:** 2026-08-01 23:07 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 199
- Candidate layers: 15 core + 20 watch + 10 sports + 10 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 58 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 15
- Aggregate closed copy trades: 285
- Aggregate copy ROI: 128.7%
- Aggregate copy PnL: $+4144.21
- Aggregate copy win rate: 84.2%
- Median wallet CopyROI: 135.5%
- Worst included wallet CopyROI: 103.7%
- Open copy cost / closed copy capital: 2.82x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xa75b...772c` | A | 24.89 | 180 | 99% | 112.3% | $+673.73 | 46 | 84.8% | 34.7% |
| `0xb2ed...4418` | A | 23.26 | 125 | 40% | 142.2% | $+583.20 | 34 | 67.7% | 58.2% |
| `0xe745...5681` | A | 5.86 | 64 | 32% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x2952...f50d` | A | 23.82 | 78 | 57% | 108.9% | $+326.79 | 22 | 86.4% | 19.5% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 58 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 62 | 93% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 27.09 | 145 | 64% | 108.6% | $+521.13 | 44 | 77.3% | 47.6% |
| `0x092b...614e` | B | 28.38 | 174 | 97% | 103.7% | $+279.93 | 25 | 88.0% | 55.6% |
| `0xe916...7e93` | B | 27.86 | 62 | 98% | 117.3% | $+152.51 | 12 | 83.3% | 44.7% |
| `0xffa1...6340` | B | 25.47 | 52 | 98% | 165.2% | $+148.73 | 9 | 88.9% | 18.0% |
| `0x0ec9...1e0c` | B | 26.70 | 7 | 13% | 183.6% | $+146.85 | 8 | 87.5% | 11.1% |
| `0xeb8b...6d8a` | B | 28.11 | 28 | 53% | 138.5% | $+124.64 | 9 | 100.0% | 42.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 787
- Aggregate copy ROI: 98.5%
- Aggregate copy PnL: $+10437.44
- Aggregate copy win rate: 89.2%
- Worst included CopyROI: 59.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.3 | 897.8 | 52 | 33% | 274.8% | $+412.21 | 12 | 75.0% | $1207 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 59 | 64% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 738.7 | 39 | 95% | 248.6% | $+223.76 | 7 | 100.0% | $7773 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.0 | 737.6 | 1337 | 100% | 115.0% | $+2667.48 | 102 | 97.1% | $1090 | existing,holder |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 29 | 48% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0x89cf...5f47` | B | 100.0 | 33.2 | 652.5 | 77 | 44% | 119.7% | $+526.61 | 43 | 88.4% | $669 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xf4e1...34a5` | B | 100.0 | 30.4 | 617.0 | 67 | 40% | 125.0% | $+374.96 | 25 | 84.0% | $5956 | existing |
| `0x2a35...9015` | B | 100.0 | 32.9 | 610.0 | 160 | 99% | 106.9% | $+545.10 | 38 | 92.1% | $1998 | existing,holder |
| `0xcc6e...fa6f` | B | 100.0 | 25.8 | 605.0 | 259 | 82% | 97.4% | $+574.59 | 51 | 88.2% | $2164 | existing,sports_tape |
| `0x17fe...b0ca` | B | 100.0 | 30.9 | 567.8 | 177 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $669 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 70 | 57% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0xb624...de17` | B | 100.0 | 26.7 | 559.9 | 27 | 100% | 195.0% | $+97.50 | 5 | 100.0% | $754 | existing |
| `0xde24...4ded` | B | 100.0 | 28.2 | 556.2 | 274 | 100% | 91.5% | $+668.26 | 37 | 100.0% | $1889 | existing,holder |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 554.3 | 16 | 26% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.3 | 525.7 | 899 | 93% | 59.8% | $+1220.55 | 185 | 89.2% | $296 | existing |
| `0x7992...1fc1` | A | 100.0 | 18.6 | 524.5 | 21 | 21% | 86.7% | $+294.71 | 29 | 79.3% | $2688 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |
| `0x7e01...b0b5` | B | 100.0 | 29.8 | 484.8 | 399 | 97% | 65.9% | $+566.73 | 60 | 80.0% | $333 | existing,holder |
| `0x18c2...529a` | B | 100.0 | 28.2 | 483.1 | 802 | 74% | 70.3% | $+520.04 | 57 | 86.0% | $272 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x17e2...d472` | B | 100.0 | 33.1 | 621.5 | 371 | 266 | 40% | 104.9% | $+618.86 | 56 | 96.4% | 128.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xcefd...d6aa` | B | 100.0 | 30.9 | 600.1 | 77 | 57 | 69% | 177.3% | $+248.17 | 9 | 100.0% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xeec5...1862` | B | 89.5 | 31.3 | 586.7 | 63 | 61 | 88% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7a26...1589` | B | 100.0 | 31.8 | 513.8 | 135 | 123 | 66% | 85.4% | $+392.90 | 37 | 91.9% | 69.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x84cd...7565` | A | 100.0 | 23.5 | 499.9 | 125 | 60 | 80% | 109.8% | $+186.61 | 14 | 85.7% | 93.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xc117...7410` | B | 84.8 | 25.4 | 464.8 | 33 | 27 | 77% | 143.0% | $+143.01 | 7 | 42.9% | 118.6% | existing | opposite_side_same_market |
| `0xa8b9...775d` | B | 100.0 | 33.1 | 446.4 | 397 | 262 | 94% | 68.0% | $+571.20 | 44 | 54.5% | 63.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x3c2b...825e` | B | 100.0 | 34.0 | 441.7 | 311 | 180 | 93% | 131.6% | $+157.95 | 7 | 71.4% | 121.8% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 28.8 | 424.4 | 168 | 82 | 99% | 65.3% | $+182.99 | 28 | 92.9% | 65.3% | existing | opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.1 | 423.1 | 87 | 86 | 64% | 73.8% | $+184.39 | 23 | 82.6% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 10
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xf3ce...a57a` | B | 100.0 | 25.9 | 301.0 | 1134 | 100% | 8.5% | 12 | 1094 | $1947 | existing,holder,leaderboard_profit_30d | open_copy_exposure,opposite_side_same_market |
| `0x16bb...8492` | B | 100.0 | 31.9 | 255.1 | 16 | 1% | 0.0% | 0 | 698 | $3766 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x27f7...44b0` | C | 100.0 | 43.5 | 230.3 | 365 | 28% | 0.0% | 0 | 678 | $723 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x1d1a...0df6` | C | 100.0 | 43.1 | 205.4 | 566 | 42% | -58.5% | 1 | 383 | $1144 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xbf00...1199` | C | 100.0 | 43.8 | 200.5 | 18 | 12% | 14.4% | 48 | 143 | $13927 | existing,leaderboard_profit_30d | opposite_side_same_market |
| `0x1465...7072` | B | 100.0 | 35.0 | 185.4 | 0 | 0% | 0.1% | 1 | 609 | $1018 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x656a...6077` | B | 100.0 | 10.2 | 180.2 | 12 | 14% | 0.0% | 0 | 58 | $753 | existing,leaderboard_profit_30d | open_copy_exposure |
| `0xc807...c12b` | C | 100.0 | 38.6 | 157.1 | 0 | 0% | 0.0% | 0 | 309 | $769 | leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 44.4 | 152.3 | 0 | 0% | 0.0% | 0 | 273 | $812 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x1b20...3ab0` | C | 86.5 | 42.4 | 134.4 | 9 | 1% | 0.0% | 0 | 229 | $1147 | leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xe872...819a` | B | 100.0 | 25.4 | 308.2 | 1060 | 840 | 90% | 74.0% | 16 | 74.0% | 16 | $1631 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2a99...51bb` | A | 100.0 | 24.1 | 265.9 | 706 | 500 | 99% | 42.1% | 127 | 42.1% | 127 | $296 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 42.5 | 250.8 | 625 | 572 | 100% | 71.6% | 62 | 71.6% | 62 | $1625 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x8dd1...5a9b` | C | 100.0 | 39.3 | 249.8 | 649 | 618 | 86% | 56.6% | 146 | 57.3% | 150 | $810 | existing | open_copy_exposure,opposite_side_same_market |
| `0xfbe8...bb28` | C | 100.0 | 35.2 | 245.7 | 1127 | 596 | 90% | 43.0% | 51 | 44.8% | 60 | $167 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd3d3...bb8e` | C | 100.0 | 39.3 | 239.9 | 858 | 618 | 71% | 13.8% | 180 | 12.7% | 255 | $637 | existing | opposite_side_same_market |
| `0x6e54...511e` | B | 100.0 | 32.0 | 237.2 | 977 | 474 | 88% | 13.3% | 20 | 12.0% | 21 | $142 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.0 | 236.8 | 93 | 89 | 93% | 40.5% | 18 | 41.6% | 19 | $2728 | existing | open_copy_exposure,opposite_side_same_market |
| `0x1e8e...f9b2` | C | 100.0 | 42.5 | 236.2 | 710 | 545 | 98% | 99.3% | 165 | 99.7% | 169 | $243 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x0ee5...e798` | C | 100.0 | 44.1 | 235.8 | 1168 | 612 | 81% | 127.2% | 107 | 128.1% | 125 | $410 | existing,holder | open_copy_exposure,opposite_side_same_market |

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
| `0x7af7...4c89` | B | 100.0 | 30.5 | 316.8 | 22 | 21 | 58.1% | 3 | 58.1% | 3 | $12821 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

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
| `0xa75b...772c` | A | 24.9 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 559.6 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 25.8 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 320.7 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | C | 41.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 218.1 | opposite_side_same_market |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0xbb35...b62a` | B | 100.0 | 34.0 | 15m edge -2.84pp over 3 samples | 115.0% | 102 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 28.2 | 1h edge -11.42pp over 4 samples | 91.5% | 37 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 27.1 | 1h edge -15.87pp over 1 samples | 108.6% | 44 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.3 | 1h edge -19.00pp over 1 samples | 142.2% | 34 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 42.5 | 1h edge -34.83pp over 1 samples | 71.6% | 62 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 25.9 | 1h edge -47.28pp over 1 samples | 8.5% | 12 | existing,holder,leaderboard_profit_30d | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.2 | 1h edge -52.95pp over 1 samples | 70.3% | 57 | existing | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 30.5 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 72.6 | $7000 | $7000 | 1 | 112.5% | 124 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 48.4 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 72.6 | 712.2 | 552 | 358 | 38% | 127.8% | 50 | 112.5% | 124 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 24.9 | 681.9 | 180 | 145 | 99% | 112.3% | 46 | 112.3% | 46 | existing,holder,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 25.8 | 535.6 | 259 | 217 | 82% | 81.7% | 41 | 97.4% | 51 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 63.2 | 510.9 | 339 | 303 | 31% | 126.5% | 14 | 112.7% | 38 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 70.0 | 293.8 | 1003 | 901 | 74% | 30.4% | 301 | 29.8% | 388 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 41.4 | 261.5 | 163 | 151 | 15% | 37.5% | 59 | 36.2% | 278 | sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject | - | D | 100.0 | 57.4 | 238.0 | 125 | 122 | 10% | 34.1% | 39 | 51.0% | 325 | sports_tape | bot_like_flow,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | D | 100.0 | 58.1 | 228.5 | 343 | 312 | 81% | 30.9% | 115 | 28.2% | 142 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x42db...d512` | prompt | - | B | 100.0 | 30.0 | 182.1 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 83.3 | 20.0 | 171.2 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 30.5 | 161.8 | 22 | 21 | 69% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | watch | - | C | 100.0 | 41.0 | 157.2 | 479 | 458 | 90% | 11.7% | 213 | 11.2% | 229 | existing,sports_tape | opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 74.9 | 10.0 | 96.6 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 24.4 | 94.2 | 20 | 5 | 31% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 15 | 285 | 128.7% | $+4144.21 | 84.2% | 135.5% | 103.7% | 2.82x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 14 | 251 | 126.7% | $+3561.01 | 86.5% | 135.4% | 103.7% | 2.65x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 3 | 27 | 655 | 99.0% | $+7950.63 | 85.2% | 108.4% | 65.3% | 3.40x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 4 | 28 | 666 | 98.4% | $+8058.06 | 85.0% | 106.1% | 65.3% | 3.37x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 13 | 207 | 130.5% | $+3039.88 | 88.4% | 135.5% | 103.7% | 2.72x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 12 | 246 | 128.1% | $+3521.84 | 83.7% | 135.4% | 103.7% | 3.04x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 26 | 621 | 96.7% | $+7367.43 | 86.2% | 106.1% | 65.3% | 3.38x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 8 | 27 | 632 | 96.1% | $+7474.86 | 85.9% | 103.7% | 65.3% | 3.34x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 9 | 11 | 237 | 127.7% | $+3397.20 | 83.1% | 135.2% | 103.7% | 3.06x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 10 | 12 | 259 | 125.8% | $+3723.99 | 83.4% | 126.3% | 103.7% | 2.86x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 11 | 212 | 125.6% | $+2938.64 | 86.3% | 135.2% | 103.7% | 2.88x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 12 | 23 | 587 | 97.4% | $+7033.55 | 85.3% | 103.7% | 65.3% | 3.62x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
