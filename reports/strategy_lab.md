# Strategy Lab Report

**Generated:** 2026-07-11 23:14 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 19
- Candidate layers: 11 core + 20 watch + 10 sports + 6 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 52 total
- Live-edge blocked push wallets: 5
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 11
- Aggregate closed copy trades: 258
- Aggregate copy ROI: 88.0%
- Aggregate copy PnL: $+2921.63
- Aggregate copy win rate: 81.0%
- Median wallet CopyROI: 70.6%
- Worst included wallet CopyROI: 58.4%
- Open copy cost / closed copy capital: 3.56x
- Params: tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0xda19...c13b` | A | 15.38 | 18 | 15% | 449.7% | $+404.73 | 8 | 100.0% | 66.1% |
| `0x7992...1fc1` | A | 19.87 | 21 | 28% | 70.6% | $+183.65 | 21 | 71.4% | 19.5% |
| `0xbe18...3be4` | A | 18.71 | 76 | 61% | 58.9% | $+176.75 | 17 | 64.7% | 20.0% |
| `0xe6bf...c536` | A | 23.11 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 20.2% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xcc6e...fa6f` | B | 28.24 | 344 | 85% | 98.6% | $+917.11 | 77 | 92.2% | 45.6% |
| `0xfbe8...bb28` | B | 28.60 | 753 | 94% | 64.6% | $+407.01 | 43 | 86.0% | 48.0% |
| `0xd8b5...54c4` | B | 25.05 | 210 | 100% | 60.0% | $+150.01 | 23 | 60.9% | 54.1% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x0931...e78e` | B | 26.24 | 106 | 66% | 67.0% | $+127.35 | 19 | 68.4% | 26.6% |
| `0x8e4e...1dfc` | B | 27.88 | 91 | 20% | 72.1% | $+100.91 | 8 | 62.5% | 48.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 556
- Aggregate copy ROI: 65.1%
- Aggregate copy PnL: $+4793.96
- Aggregate copy win rate: 82.6%
- Worst included CopyROI: 22.4%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x0f00...73ce` | B | 96.5 | 25.1 | 573.8 | 79 | 77% | 233.0% | $+163.12 | 4 | 100.0% | $2316 | existing |
| `0x3abd...855d` | B | 100.0 | 23.5 | 568.5 | 24 | 92% | 214.7% | $+85.90 | 4 | 100.0% | $3297 | holder |
| `0xf4ec...c574` | B | 100.0 | 33.7 | 560.0 | 950 | 89% | 71.7% | $+1584.17 | 143 | 90.9% | $2072 | holder |
| `0xb36f...53d0` | B | 100.0 | 30.7 | 504.8 | 651 | 93% | 58.6% | $+891.03 | 139 | 87.8% | $294 | existing,holder |
| `0xc117...7410` | B | 86.4 | 17.9 | 498.9 | 25 | 89% | 194.3% | $+136.05 | 4 | 50.0% | $1691 | holder |
| `0xd5b1...1b71` | B | 100.0 | 31.6 | 490.4 | 122 | 98% | 81.3% | $+252.01 | 30 | 93.3% | $523 | existing |
| `0x2cf0...df85` | B | 100.0 | 33.1 | 483.9 | 130 | 41% | 91.2% | $+264.52 | 22 | 86.4% | $316 | holder |
| `0x7c73...6ee3` | B | 100.0 | 25.8 | 437.7 | 43 | 96% | 120.8% | $+60.39 | 5 | 100.0% | $5230 | existing |
| `0x07a9...32c8` | B | 100.0 | 27.8 | 437.0 | 34 | 59% | 108.6% | $+76.05 | 6 | 100.0% | $21206 | existing |
| `0x8bb0...ca79` | B | 100.0 | 32.4 | 426.9 | 66 | 32% | 87.0% | $+174.00 | 13 | 76.9% | $1518 | existing |
| `0xe872...819a` | B | 100.0 | 33.4 | 418.9 | 1026 | 85% | 68.5% | $+219.25 | 20 | 85.0% | $1548 | existing,holder,sports_tape |
| `0xc419...a8b0` | B | 100.0 | 28.7 | 369.9 | 35 | 85% | 85.9% | $+60.12 | 6 | 83.3% | $1038 | existing |
| `0x9ecc...7850` | B | 100.0 | 30.4 | 362.9 | 4 | 4% | 53.5% | $+160.46 | 30 | 73.3% | $1563 | existing |
| `0x68cb...6f57` | A | 100.0 | 22.1 | 355.1 | 55 | 60% | 41.2% | $+123.72 | 25 | 64.0% | $550 | existing |
| `0xdbd5...cba7` | A | 100.0 | 18.3 | 354.0 | 90 | 99% | 39.8% | $+127.32 | 17 | 94.1% | $2620 | existing |
| `0x6b51...ad13` | B | 100.0 | 23.6 | 347.6 | 10 | 53% | 79.5% | $+39.73 | 4 | 100.0% | $6564 | holder |
| `0xb53e...8617` | B | 100.0 | 34.5 | 323.2 | 45 | 43% | 40.2% | $+120.58 | 19 | 57.9% | $1464 | holder |
| `0x0e24...7014` | B | 100.0 | 27.8 | 317.4 | 61 | 100% | 33.9% | $+101.70 | 15 | 73.3% | $1004 | existing,holder |
| `0x6d15...04c2` | B | 100.0 | 26.2 | 308.9 | 31 | 79% | 46.1% | $+50.74 | 8 | 50.0% | $570 | holder |
| `0x42db...d512` | B | 100.0 | 30.4 | 303.3 | 108 | 100% | 22.4% | $+103.10 | 42 | 54.8% | $1184 | existing,sports_tape |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.4 | 415.7 | 88 | 63 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 31.3 | 382.0 | 103 | 48 | 100% | 61.3% | $+110.27 | 18 | 100.0% | 61.3% | existing,holder | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x631a...a80f` | B | 100.0 | 33.6 | 360.9 | 145 | 67 | 76% | 60.7% | $+109.30 | 17 | 88.2% | 61.9% | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0xabb2...0bc2` | B | 100.0 | 32.0 | 324.6 | 82 | 59 | 43% | 53.2% | $+74.53 | 14 | 85.7% | 58.4% | holder | open_copy_exposure,opposite_side_same_market |
| `0xba8a...b772` | B | 100.0 | 32.1 | 319.0 | 54 | 50 | 35% | 47.4% | $+80.54 | 16 | 81.2% | 9.8% | holder | opposite_side_same_market |
| `0xd8f6...15fb` | A | 100.0 | 21.6 | 308.2 | 113 | 68 | 44% | 35.1% | $+112.27 | 21 | 57.1% | 28.3% | holder | opposite_side_same_market |
| `0xc6f1...5729` | B | 100.0 | 31.6 | 255.5 | 160 | 153 | 90% | 14.8% | $+108.14 | 70 | 71.4% | 12.7% | existing,sports_tape | opposite_side_same_market |
| `0xfe18...3980` | A | 84.0 | 24.4 | 251.8 | 94 | 70 | 65% | 27.5% | $+96.28 | 33 | 21.2% | 17.9% | existing,sports_tape | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 6
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x361b...74fe` | C | 100.0 | 44.6 | 213.3 | 23 | 2% | 0.0% | 0 | 583 | $2477 | leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xe549...2d54` | C | 100.0 | 42.3 | 200.7 | 588 | 40% | 0.0% | 0 | 264 | $612 | existing,holder,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0xa710...23c4` | C | 100.0 | 41.7 | 174.0 | 0 | 0% | 0.0% | 0 | 364 | $621 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 100.0 | 44.4 | 165.6 | 0 | 0% | 0.0% | 0 | 245 | $4587 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x3a8a...7699` | C | 100.0 | 42.2 | 164.0 | 0 | 0% | -33.4% | 2 | 380 | $852 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 42.9 | 147.6 | 0 | 0% | 0.0% | 0 | 232 | $521 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 35.8 | 316.4 | 1290 | 1063 | 100% | 14.4% | 234 | 14.4% | 234 | $482 | existing,holder | burst_trading,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.9 | 274.9 | 1146 | 849 | 100% | 74.3% | 217 | 74.3% | 217 | $317 | existing,holder | opposite_side_same_market |
| `0x34a5...9450` | C | 100.0 | 42.5 | 240.9 | 538 | 499 | 95% | 101.2% | 85 | 98.1% | 90 | $1705 | holder | open_copy_exposure,opposite_side_same_market |
| `0x0562...18a4` | C | 100.0 | 35.8 | 231.5 | 293 | 266 | 99% | 83.7% | 72 | 83.7% | 72 | $3272 | holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.4 | 224.9 | 254 | 231 | 100% | 65.9% | 28 | 65.9% | 28 | $1477 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x8f76...aea6` | C | 100.0 | 41.5 | 211.4 | 689 | 488 | 51% | 28.6% | 95 | 21.9% | 165 | $568 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4c9c...e006` | C | 100.0 | 41.9 | 208.5 | 284 | 231 | 98% | 22.8% | 75 | 23.1% | 77 | $423 | existing,sports_tape | opposite_side_same_market |
| `0x2391...554f` | C | 100.0 | 41.5 | 198.1 | 101 | 99 | 99% | 85.6% | 33 | 85.6% | 33 | $1551 | existing,holder,retain | opposite_side_same_market |
| `0xaf6a...925f` | C | 100.0 | 37.8 | 196.6 | 93 | 85 | 100% | 86.8% | 21 | 86.8% | 21 | $886 | existing | open_copy_exposure,opposite_side_same_market |
| `0x96d7...3a17` | C | 100.0 | 37.7 | 192.6 | 91 | 52 | 100% | 53.9% | 5 | 53.9% | 5 | $553 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |

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
| `0xe872...819a` | B | 33.4 | $2200 | $2200 | 51 | 70.6% | +8.72 | +1.61 | +6.46 | +16.96 | 294.5 | open_copy_exposure,opposite_side_same_market |
| `0x4c9c...e006` | C | 41.9 | $800 | $800 | 12 | 75.0% | +16.81 | +4.67 | +14.65 | +30.95 | 288.0 | opposite_side_same_market |
| `0xc6f1...5729` | B | 31.6 | $6870 | $6870 | 38 | 71.1% | +10.37 | +1.18 | +3.95 | +15.97 | 253.7 | opposite_side_same_market |

## Live Edge Blocked Push Wallets

- Wallets: 5
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | B | 96.5 | 25.1 | 15m edge -10.00pp over 2 samples | 233.0% | 4 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | B | 100.0 | 25.1 | 15m edge -24.00pp over 2 samples | 60.0% | 23 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.4 | 1h edge -34.83pp over 1 samples | 65.9% | 28 | existing,holder | open_copy_exposure,opposite_side_same_market |
| `0x96d7...3a17` | C | 100.0 | 37.7 | 1h edge -35.95pp over 1 samples | 53.9% | 5 | existing,holder | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 27.8 | 1h edge -72.59pp over 1 samples | 33.9% | 15 | existing,holder | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | D | 18.3 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | B | 21.7 | $15777 | $15777 | 1 | -33.4% | 1 | open_copy_exposure |
| `0x9520...fa6e` | watch | - | B | 33.6 | $14778 | $14778 | 1 | 7.1% | 25 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 68.4 | $7000 | $7000 | 1 | 115.1% | 114 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | reject-bot | - | D | 55.3 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.2 | 651.7 | 417 | 351 | 36% | 158.3% | 17 | 127.0% | 39 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 100.0 | 68.4 | 437.5 | 124 | 122 | 8% | 81.9% | 19 | 115.1% | 114 | existing,holder,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | pushed | - | B | 100.0 | 33.4 | 406.2 | 1026 | 815 | 85% | 70.6% | 19 | 68.5% | 20 | existing,holder,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.7 | 258.4 | 862 | 778 | 72% | 27.6% | 265 | 27.5% | 347 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 43.1 | 249.7 | 136 | 126 | 19% | 37.0% | 51 | 37.5% | 229 | holder,sports_tape | opposite_side_same_market |
| `0x4c9c...e006` | pushed | - | C | 100.0 | 41.9 | 191.6 | 284 | 231 | 98% | 22.8% | 75 | 23.1% | 77 | existing,sports_tape | opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.4 | 181.5 | 108 | 103 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0xfe18...3980` | pushed | - | A | 84.0 | 24.4 | 169.1 | 94 | 70 | 65% | 27.5% | 33 | 17.9% | 44 | existing,sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 76.5 | 167.6 | 150 | 147 | 12% | 28.4% | 47 | 42.5% | 270 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,opposite_side_same_market |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 40.2 | 150.7 | 128 | 114 | 70% | 22.4% | 41 | 16.9% | 59 | existing,sports_tape | opposite_side_same_market |
| `0xc6f1...5729` | pushed | - | B | 100.0 | 31.6 | 148.1 | 160 | 153 | 90% | 14.8% | 70 | 12.7% | 76 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.6 | 128.8 | 53 | 38 | 33% | 117.0% | 1 | 68.1% | 2 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 16.0 | 115.1 | 7 | 5 | 70% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xbc25...aaba` | prompt | - | B | 100.0 | 3.9 | 89.2 | 17 | 7 | 94% | 0.0% | 0 | 0.0% | 0 | sports_tape | open_copy_exposure |
| `0xb3c1...4837` | prompt | - | B | 100.0 | 21.2 | 81.1 | 194 | 146 | 91% | -0.1% | 21 | 0.3% | 23 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 11 | 258 | 88.0% | $+2921.63 | 81.0% | 70.6% | 58.4% | 3.56x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 2 | 12 | 267 | 86.6% | $+2969.10 | 80.9% | 68.8% | 43.2% | 3.63x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 3 | 12 | 283 | 84.1% | $+3045.35 | 79.5% | 68.8% | 41.2% | 3.30x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 13 | 292 | 82.9% | $+3092.82 | 79.5% | 67.0% | 41.2% | 3.37x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 10 | 237 | 89.5% | $+2737.98 | 81.9% | 69.6% | 58.4% | 3.76x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=50 copyWin>=60 closedROI>=20 smart>=70 |
| 6 | 11 | 246 | 87.9% | $+2785.45 | 81.7% | 67.0% | 43.2% | 3.83x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 13 | 290 | 80.0% | $+3150.65 | 81.4% | 67.0% | 33.9% | 3.12x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 16 | 315 | 76.9% | $+3270.32 | 81.0% | 62.3% | 33.1% | 3.09x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 9 | 14 | 315 | 77.2% | $+3274.37 | 80.0% | 65.8% | 33.9% | 2.93x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 17 | 340 | 74.6% | $+3394.04 | 79.7% | 60.0% | 33.1% | 2.92x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 12 | 269 | 80.6% | $+2967.00 | 82.2% | 65.8% | 33.9% | 3.26x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=20 smart>=70 |
| 12 | 11 | 240 | 81.7% | $+2688.84 | 86.2% | 64.6% | 33.9% | 2.22x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=70 closedROI>=0 smart>=70 |
