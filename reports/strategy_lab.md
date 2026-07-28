# Strategy Lab Report

**Generated:** 2026-07-28 23:29 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 164
- Candidate layers: 14 core + 20 watch + 10 sports + 7 scout + 10 target + 0 flow + 1 tape
- Push wallets after live-edge blocks: 54 total
- Live-edge blocked push wallets: 8
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 14
- Aggregate closed copy trades: 248
- Aggregate copy ROI: 129.5%
- Aggregate copy PnL: $+3793.38
- Aggregate copy win rate: 82.7%
- Median wallet CopyROI: 135.4%
- Worst included wallet CopyROI: 103.8%
- Open copy cost / closed copy capital: 3.21x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xb2ed...4418` | A | 23.41 | 117 | 39% | 144.0% | $+575.85 | 33 | 66.7% | 53.5% |
| `0xa75b...772c` | A | 24.26 | 127 | 99% | 103.8% | $+467.06 | 33 | 78.8% | 32.0% |
| `0xe745...5681` | A | 5.86 | 57 | 29% | 187.3% | $+355.81 | 19 | 100.0% | 29.9% |
| `0x2952...f50d` | A | 24.10 | 79 | 57% | 109.8% | $+340.44 | 23 | 87.0% | 19.5% |
| `0xfc92...5667` | A | 23.57 | 127 | 91% | 157.1% | $+251.29 | 9 | 88.9% | 56.2% |
| `0x6f16...5fe7` | A | 24.93 | 42 | 63% | 192.8% | $+231.31 | 12 | 83.3% | 33.9% |
| `0x90e0...21a2` | A | 22.22 | 53 | 98% | 143.1% | $+171.66 | 12 | 91.7% | 36.5% |
| `0xabff...9e8f` | A | 24.76 | 57 | 46% | 135.5% | $+149.04 | 11 | 100.0% | 48.7% |
| `0xb715...c3bb` | A | 24.63 | 60 | 90% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0x44c4...09cb` | B | 26.68 | 144 | 60% | 109.5% | $+536.49 | 45 | 77.8% | 47.6% |
| `0x3c2b...825e` | B | 25.88 | 279 | 92% | 121.8% | $+158.40 | 8 | 75.0% | 56.1% |
| `0xe916...7e93` | B | 28.09 | 60 | 100% | 117.3% | $+152.51 | 12 | 83.3% | 44.7% |
| `0xeb8b...6d8a` | B | 28.11 | 27 | 51% | 138.5% | $+124.64 | 9 | 100.0% | 42.5% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 732
- Aggregate copy ROI: 100.5%
- Aggregate copy PnL: $+9684.30
- Aggregate copy win rate: 89.1%
- Worst included CopyROI: 60.1%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xda19...c13b` | B | 100.0 | 34.3 | 897.8 | 25 | 16% | 274.8% | $+412.21 | 12 | 75.0% | $1207 | existing |
| `0xeca2...44fc` | B | 100.0 | 32.5 | 752.1 | 58 | 63% | 222.3% | $+311.21 | 11 | 90.9% | $354 | existing |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 738.7 | 39 | 95% | 248.6% | $+223.76 | 7 | 100.0% | $7773 | existing |
| `0xbb35...b62a` | B | 100.0 | 34.3 | 709.9 | 1340 | 100% | 109.9% | $+2494.88 | 102 | 96.1% | $1032 | existing |
| `0xc367...3066` | B | 100.0 | 31.2 | 673.6 | 28 | 47% | 264.1% | $+158.44 | 5 | 100.0% | $326 | existing |
| `0x89cf...5f47` | B | 100.0 | 33.3 | 653.2 | 70 | 41% | 119.7% | $+526.61 | 43 | 88.4% | $672 | existing |
| `0xb05d...dcdc` | B | 99.1 | 34.4 | 642.8 | 85 | 83% | 215.4% | $+150.76 | 7 | 100.0% | $256 | existing |
| `0xcc6e...fa6f` | B | 100.0 | 25.9 | 614.3 | 276 | 83% | 97.4% | $+623.47 | 56 | 89.3% | $2113 | existing,sports_tape |
| `0x2a35...9015` | B | 100.0 | 33.1 | 602.5 | 117 | 99% | 115.3% | $+415.06 | 28 | 92.9% | $1858 | existing |
| `0xf4e1...34a5` | B | 100.0 | 30.6 | 574.8 | 68 | 41% | 115.2% | $+322.47 | 23 | 82.6% | $5973 | existing |
| `0x17fe...b0ca` | B | 100.0 | 30.2 | 568.6 | 189 | 100% | 95.6% | $+420.79 | 42 | 76.2% | $629 | existing |
| `0xde24...4ded` | B | 100.0 | 30.8 | 563.6 | 145 | 100% | 103.8% | $+508.64 | 25 | 100.0% | $2689 | existing |
| `0x9caf...94dc` | B | 100.0 | 32.1 | 562.0 | 63 | 51% | 103.7% | $+300.79 | 28 | 96.4% | $2918 | existing |
| `0x73e2...46ec` | A | 100.0 | 5.2 | 546.1 | 12 | 19% | 129.2% | $+142.10 | 11 | 90.9% | $2834 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.4 | 526.7 | 889 | 93% | 60.1% | $+1219.76 | 184 | 89.1% | $298 | existing |
| `0x7992...1fc1` | A | 100.0 | 18.6 | 524.5 | 21 | 21% | 86.7% | $+294.71 | 29 | 79.3% | $2688 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.0 | 488.4 | 137 | 98% | 79.0% | $+260.85 | 32 | 93.8% | $475 | existing |
| `0x2cf0...df85` | B | 100.0 | 32.9 | 482.5 | 144 | 43% | 88.5% | $+274.36 | 24 | 87.5% | $359 | existing |
| `0x18c2...529a` | B | 100.0 | 28.4 | 478.7 | 661 | 71% | 72.9% | $+444.95 | 47 | 83.0% | $273 | existing |
| `0x84cd...7565` | A | 100.0 | 23.9 | 474.7 | 113 | 71% | 93.9% | $+178.48 | 16 | 75.0% | $439 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7a26...1589` | B | 100.0 | 31.7 | 499.9 | 83 | 78 | 56% | 94.1% | $+263.37 | 23 | 95.7% | 71.4% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xc117...7410` | B | 84.8 | 25.4 | 464.8 | 33 | 27 | 77% | 143.0% | $+143.01 | 7 | 42.9% | 118.6% | existing | opposite_side_same_market |
| `0xa8b9...775d` | B | 100.0 | 33.2 | 446.5 | 381 | 258 | 93% | 68.0% | $+571.20 | 44 | 54.5% | 63.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 29.0 | 444.0 | 162 | 79 | 100% | 71.4% | $+192.74 | 27 | 96.3% | 71.4% | existing | opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.6 | 434.6 | 423 | 397 | 100% | 67.6% | $+554.10 | 39 | 94.9% | 67.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x119b...ac3e` | B | 100.0 | 33.1 | 431.6 | 80 | 80 | 59% | 76.9% | $+184.49 | 22 | 86.4% | 62.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x5a56...10f5` | B | 100.0 | 31.6 | 420.5 | 90 | 75 | 85% | 94.5% | $+245.67 | 11 | 63.6% | 88.6% | existing | open_copy_exposure,opposite_side_same_market |
| `0x578e...c3c0` | A | 100.0 | 23.9 | 415.6 | 57 | 53 | 33% | 71.3% | $+135.52 | 18 | 88.9% | 59.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.2 | 415.5 | 89 | 64 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |

## Leaderboard Scout

- Wallets: 7
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x5f65...2036` | C | 100.0 | 37.7 | 320.9 | 40 | 3% | 82.5% | 265 | 1385 | $10177 | existing,leaderboard_profit_7d | open_copy_exposure,opposite_side_same_market |
| `0x16bb...8492` | B | 100.0 | 31.5 | 256.8 | 16 | 1% | 0.0% | 0 | 715 | $3622 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x56ac...d77e` | C | 100.0 | 40.3 | 228.5 | 46 | 6% | 0.0% | 0 | 520 | $2489 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0x0d18...24a6` | C | 100.0 | 43.2 | 199.9 | 552 | 38% | 0.0% | 0 | 333 | $632 | leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x7c1e...bab3` | C | 100.0 | 43.9 | 194.5 | 645 | 44% | 0.0% | 0 | 271 | $1070 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 42.5 | 171.9 | 0 | 0% | 0.0% | 0 | 326 | $813 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x476e...396d` | C | 100.0 | 38.9 | 169.9 | 19 | 1% | 0.0% | 0 | 309 | $674 | existing,leaderboard_profit_30d,leaderboard_profit_all | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xe872...819a` | B | 100.0 | 25.4 | 300.7 | 1016 | 797 | 86% | 76.6% | 15 | 74.0% | 16 | $1624 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2a99...51bb` | A | 100.0 | 24.5 | 263.0 | 717 | 484 | 99% | 40.5% | 121 | 40.5% | 121 | $284 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x8dd1...5a9b` | C | 100.0 | 39.5 | 255.6 | 701 | 660 | 89% | 56.4% | 153 | 57.0% | 157 | $774 | existing | open_copy_exposure,opposite_side_same_market |
| `0xfbe8...bb28` | C | 100.0 | 35.5 | 244.3 | 1070 | 554 | 93% | 57.6% | 48 | 57.2% | 57 | $163 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xdbd5...cba7` | A | 100.0 | 18.0 | 236.8 | 93 | 89 | 93% | 40.5% | 18 | 41.6% | 19 | $2728 | existing | open_copy_exposure,opposite_side_same_market |
| `0x6e54...511e` | B | 100.0 | 31.9 | 233.4 | 973 | 447 | 87% | 12.7% | 22 | 11.5% | 23 | $135 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x1e8e...f9b2` | C | 100.0 | 42.5 | 230.5 | 665 | 501 | 98% | 97.1% | 152 | 97.5% | 156 | $231 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xe6bf...c536` | A | 100.0 | 23.1 | 228.3 | 73 | 69 | 95% | 58.4% | 20 | 58.4% | 20 | $3677 | existing | opposite_side_same_market |
| `0x42db...d512` | B | 100.0 | 30.0 | 222.5 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | $1174 | existing,sports_tape | opposite_side_same_market |
| `0x0378...7634` | B | 100.0 | 25.4 | 217.7 | 311 | 130 | 100% | 28.3% | 45 | 28.3% | 45 | $338 | existing | opposite_side_same_market |

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

- Wallets: 3
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xa75b...772c` | A | 24.3 | $3000 | $3000 | 4 | 100.0% | +30.74 | +5.50 | +20.50 | +50.95 | 536.8 | opposite_side_same_market |
| `0xcc6e...fa6f` | B | 25.9 | $14442 | $14442 | 12 | 66.7% | +6.70 | +6.49 | +6.91 | +7.24 | 325.4 | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 8
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0xbb35...b62a` | B | 100.0 | 34.3 | 15m edge -2.84pp over 3 samples | 109.9% | 102 | existing | open_copy_exposure,opposite_side_same_market |
| `0xde24...4ded` | B | 100.0 | 30.8 | 1h edge -11.42pp over 4 samples | 103.8% | 25 | existing | open_copy_exposure,opposite_side_same_market |
| `0x44c4...09cb` | B | 100.0 | 26.7 | 1h edge -15.87pp over 1 samples | 109.5% | 45 | existing | open_copy_exposure,opposite_side_same_market |
| `0x5b1d...3721` | B | 100.0 | 28.3 | 1h edge -17.38pp over 1 samples | 248.6% | 7 | existing | open_copy_exposure,opposite_side_same_market |
| `0xb2ed...4418` | A | 100.0 | 23.4 | 1h edge -19.00pp over 1 samples | 144.0% | 33 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 34.6 | 1h edge -34.83pp over 1 samples | 67.6% | 39 | existing | open_copy_exposure,opposite_side_same_market |
| `0x18c2...529a` | B | 100.0 | 28.4 | 1h edge -52.95pp over 1 samples | 72.9% | 47 | existing | open_copy_exposure,opposite_side_same_market |
| `0x7c1e...bab3` | C | 100.0 | 43.9 | 1h edge -8.50pp over 2 samples | 0.0% | 0 | existing,leaderboard_profit_30d | burst_trading,open_copy_exposure |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | B | 10.0 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | pushed | - | B | 30.5 | $15777 | $15777 | 1 | 58.1% | 3 | open_copy_exposure,opposite_side_same_market |
| `0x9520...fa6e` | blocked-edge | 15m edge -1.49pp over 2 samples | B | 33.2 | $14778 | $14778 | 1 | 6.9% | 28 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 73.4 | $7000 | $7000 | 1 | 115.4% | 131 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | blocked-edge | 15m edge -1.01pp over 13 samples | D | 55.4 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa4fd...748b` | blocked-edge | 1h edge -21.51pp over 2 samples | BOT | 100.0 | 73.4 | 641.2 | 267 | 254 | 18% | 119.5% | 38 | 115.4% | 131 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.3 | 621.6 | 389 | 340 | 34% | 149.0% | 17 | 124.4% | 41 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa75b...772c` | pushed | - | A | 100.0 | 24.3 | 586.9 | 127 | 105 | 99% | 103.8% | 33 | 103.8% | 33 | existing,sports_tape | opposite_side_same_market |
| `0xcc6e...fa6f` | pushed | - | B | 100.0 | 25.9 | 555.3 | 276 | 231 | 83% | 83.2% | 46 | 97.4% | 56 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 70.1 | 297.2 | 1008 | 903 | 74% | 30.7% | 310 | 29.9% | 397 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 41.4 | 261.5 | 163 | 151 | 15% | 37.5% | 59 | 36.2% | 278 | sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject | - | D | 100.0 | 58.3 | 228.1 | 137 | 134 | 11% | 32.4% | 43 | 48.3% | 325 | sports_tape | bot_like_flow,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | D | 100.0 | 57.8 | 200.8 | 288 | 262 | 79% | 28.6% | 93 | 25.7% | 119 | holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.0 | 182.1 | 110 | 104 | 99% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0x8da2...7fea` | prompt | - | B | 83.3 | 20.0 | 171.2 | 25 | 20 | 100% | 59.7% | 3 | 59.7% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | prompt | - | B | 83.1 | 25.3 | 164.2 | 99 | 73 | 65% | 26.7% | 34 | 17.4% | 46 | existing,sports_tape | opposite_side_same_market |
| `0x7af7...4c89` | pushed | - | B | 100.0 | 30.5 | 160.9 | 21 | 20 | 66% | 58.1% | 3 | 58.1% | 3 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xe0cc...aafc` | prompt | - | B | 74.9 | 10.0 | 96.6 | 11 | 7 | 100% | 25.8% | 1 | 25.8% | 1 | existing,sports_tape | - |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 24.2 | 94.7 | 20 | 5 | 32% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xb3c1...4837` | prompt | - | B | 100.0 | 21.1 | 81.7 | 198 | 150 | 91% | -0.1% | 21 | 0.3% | 23 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 14 | 248 | 129.5% | $+3793.38 | 82.7% | 135.4% | 103.8% | 3.21x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 13 | 225 | 131.8% | $+3452.94 | 82.2% | 135.5% | 103.8% | 3.46x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 3 | 25 | 527 | 103.8% | $+6466.39 | 83.9% | 108.4% | 66.5% | 4.13x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 24 | 516 | 104.8% | $+6358.96 | 84.1% | 109.0% | 66.5% | 4.19x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 5 | 13 | 215 | 127.2% | $+3217.53 | 85.1% | 135.2% | 103.8% | 3.12x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 6 | 13 | 229 | 123.6% | $+3400.15 | 81.7% | 135.2% | 81.2% | 2.67x | tier>=A bot<25 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 7 | 22 | 464 | 105.6% | $+5723.81 | 84.3% | 109.0% | 66.5% | 4.54x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 8 | 12 | 192 | 129.6% | $+2877.09 | 84.9% | 135.4% | 103.8% | 3.40x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=70 closedROI>=20 smart>=70 |
| 9 | 23 | 483 | 102.0% | $+5783.11 | 85.3% | 108.4% | 66.5% | 4.22x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=5 smart>=70 |
| 10 | 11 | 222 | 127.8% | $+3259.05 | 82.0% | 135.2% | 103.8% | 2.80x | tier>=B bot<30 copyT>=10 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 24 | 494 | 101.0% | $+5890.54 | 85.0% | 106.1% | 66.5% | 4.15x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
| 12 | 10 | 174 | 135.0% | $+2821.34 | 83.3% | 139.3% | 103.8% | 2.89x | tier>=A bot<25 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
