# Strategy Lab Report

**Generated:** 2026-07-11 23:04 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 3
- Candidate layers: 10 core + 20 watch + 8 sports + 4 scout + 10 target + 1 flow + 0 tape
- Push wallets after live-edge blocks: 43 total
- Live-edge blocked push wallets: 6
- Leaderboard scout push enabled: false
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 10
- Aggregate closed copy trades: 189
- Aggregate copy ROI: 61.7%
- Aggregate copy PnL: $+1598.62
- Aggregate copy win rate: 78.3%
- Median wallet CopyROI: 62.3%
- Worst included wallet CopyROI: 33.9%
- Open copy cost / closed copy capital: 2.51x
- Params: tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x7992...1fc1` | A | 19.87 | 21 | 28% | 70.6% | $+183.65 | 21 | 71.4% | 19.5% |
| `0xe6bf...c536` | A | 23.11 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 20.2% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xdbd5...cba7` | A | 18.30 | 90 | 99% | 39.8% | $+127.32 | 17 | 94.1% | 26.3% |
| `0xfbe8...bb28` | B | 28.60 | 753 | 94% | 64.6% | $+407.01 | 43 | 86.0% | 48.0% |
| `0xd8b5...54c4` | B | 25.05 | 210 | 100% | 60.0% | $+150.01 | 23 | 60.9% | 54.1% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0x0931...e78e` | B | 26.24 | 106 | 66% | 67.0% | $+127.35 | 19 | 68.4% | 26.6% |
| `0x0e24...7014` | B | 27.79 | 61 | 100% | 33.9% | $+101.70 | 15 | 73.3% | 29.1% |
| `0x5c7a...fdbe` | B | 27.60 | 31 | 29% | 43.2% | $+47.47 | 9 | 77.8% | 40.6% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 554
- Aggregate copy ROI: 45.3%
- Aggregate copy PnL: $+2903.79
- Aggregate copy win rate: 75.1%
- Worst included CopyROI: 12.9%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x0f00...73ce` | B | 96.5 | 25.1 | 573.8 | 79 | 77% | 233.0% | $+163.12 | 4 | 100.0% | $2316 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.6 | 502.2 | 650 | 93% | 58.8% | $+888.02 | 138 | 87.7% | $294 | existing |
| `0xd5b1...1b71` | B | 100.0 | 31.6 | 490.4 | 122 | 98% | 81.3% | $+252.01 | 30 | 93.3% | $523 | existing |
| `0xe872...819a` | B | 100.0 | 33.4 | 476.7 | 1419 | 87% | 82.6% | $+363.50 | 25 | 88.0% | $1278 | existing,sports_tape |
| `0x7c73...6ee3` | B | 100.0 | 25.8 | 437.7 | 43 | 96% | 120.8% | $+60.39 | 5 | 100.0% | $5230 | existing |
| `0x8bb0...ca79` | B | 100.0 | 32.4 | 426.9 | 66 | 32% | 87.0% | $+174.00 | 13 | 76.9% | $1518 | existing |
| `0xc419...a8b0` | B | 100.0 | 28.7 | 369.9 | 35 | 85% | 85.9% | $+60.12 | 6 | 83.3% | $1038 | existing |
| `0x9ecc...7850` | B | 100.0 | 30.4 | 362.9 | 4 | 4% | 53.5% | $+160.46 | 30 | 73.3% | $1563 | existing |
| `0x68cb...6f57` | A | 100.0 | 22.1 | 355.1 | 55 | 60% | 41.2% | $+123.72 | 25 | 64.0% | $550 | existing |
| `0x42db...d512` | B | 100.0 | 30.4 | 303.3 | 108 | 100% | 22.4% | $+103.10 | 42 | 54.8% | $1184 | existing,sports_tape |
| `0x8a35...904a` | A | 100.0 | 24.1 | 290.6 | 0 | 0% | 54.1% | $+75.74 | 6 | 100.0% | $896 | existing |
| `0xa8b8...36dd` | A | 100.0 | 20.2 | 290.4 | 75 | 54% | 37.8% | $+22.69 | 6 | 100.0% | $493 | existing |
| `0x2393...1a8d` | A | 100.0 | 21.8 | 284.6 | 34 | 61% | 29.6% | $+35.51 | 12 | 58.3% | $599 | existing,sports_tape |
| `0x8e17...9d3c` | B | 100.0 | 34.0 | 280.6 | 170 | 62% | 16.1% | $+159.88 | 89 | 57.3% | $1336 | existing |
| `0xbeea...0c83` | B | 100.0 | 16.0 | 277.8 | 7 | 70% | 46.4% | $+13.91 | 3 | 66.7% | $413 | existing,sports_tape |
| `0x1c62...5805` | B | 100.0 | 33.4 | 274.9 | 25 | 56% | 35.7% | $+28.53 | 8 | 75.0% | $344 | existing |
| `0xc6f1...5729` | B | 100.0 | 31.7 | 282.0 | 159 | 90% | 12.9% | $+100.30 | 75 | 69.3% | $1409 | existing,sports_tape |
| `0x9f03...766d` | A | 100.0 | 23.8 | 262.9 | 8 | 18% | 29.8% | $+17.85 | 6 | 50.0% | $25233 | leaderboard_profit_7d |
| `0x3c14...a20a` | B | 100.0 | 32.8 | 255.7 | 10 | 10% | 23.2% | $+64.92 | 26 | 92.3% | $423 | existing |
| `0xe4b1...a75e` | B | 100.0 | 30.0 | 252.5 | 31 | 41% | 22.5% | $+36.02 | 5 | 60.0% | $2712 | existing |

## Sports Positive Copy

- Wallets: 8
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 89.5 | 31.3 | 587.3 | 61 | 60 | 85% | 189.2% | $+245.96 | 9 | 44.4% | 144.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x59a8...eb4e` | B | 100.0 | 33.4 | 415.7 | 88 | 63 | 19% | 79.9% | $+151.82 | 17 | 76.5% | 86.2% | recent_trade | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 31.3 | 382.0 | 103 | 48 | 100% | 61.3% | $+110.27 | 18 | 100.0% | 61.3% | existing | opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x631a...a80f` | B | 100.0 | 33.6 | 360.9 | 145 | 67 | 76% | 60.7% | $+109.30 | 17 | 88.2% | 61.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xfe18...3980` | A | 84.0 | 24.4 | 251.8 | 94 | 70 | 65% | 27.5% | $+96.28 | 33 | 21.2% | 17.9% | existing,sports_tape | opposite_side_same_market |
| `0x7020...fd20` | B | 100.0 | 32.4 | 243.5 | 54 | 47 | 44% | 34.6% | $+41.48 | 7 | 85.7% | 11.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0x94de...e3ca` | B | 100.0 | 30.3 | 225.7 | 87 | 26 | 75% | 18.6% | $+14.90 | 8 | 100.0% | 9.6% | recent_trade | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 4
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xa710...23c4` | C | 100.0 | 41.6 | 178.3 | 0 | 0% | 0.0% | 0 | 430 | $551 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x4cb5...ef29` | C | 100.0 | 44.4 | 177.6 | 0 | 0% | 0.0% | 0 | 245 | $4587 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x3a8a...7699` | C | 100.0 | 41.7 | 172.9 | 1 | 0% | -30.7% | 2 | 459 | $813 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |
| `0x5e94...5ba1` | C | 100.0 | 43.2 | 155.1 | 0 | 0% | 0.0% | 0 | 335 | $506 | existing,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 35.3 | 352.9 | 1753 | 1389 | 100% | 22.2% | 302 | 22.2% | 302 | $460 | existing | burst_trading,opposite_side_same_market |
| `0xb5b5...205a` | C | 100.0 | 42.9 | 271.9 | 1146 | 849 | 100% | 74.3% | 217 | 74.3% | 217 | $317 | existing | opposite_side_same_market |
| `0x8f76...aea6` | C | 100.0 | 41.5 | 230.2 | 922 | 667 | 51% | 23.7% | 141 | 16.9% | 252 | $593 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.4 | 221.9 | 254 | 231 | 100% | 65.9% | 28 | 65.9% | 28 | $1477 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4c9c...e006` | C | 100.0 | 41.9 | 208.5 | 284 | 231 | 98% | 22.8% | 75 | 23.1% | 77 | $423 | existing,sports_tape | opposite_side_same_market |
| `0xaf6a...925f` | C | 100.0 | 37.8 | 196.6 | 93 | 85 | 100% | 86.8% | 21 | 86.8% | 21 | $886 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2391...554f` | C | 100.0 | 41.5 | 195.1 | 101 | 99 | 99% | 85.6% | 33 | 85.6% | 33 | $1551 | existing,retain | opposite_side_same_market |
| `0x96d7...3a17` | C | 100.0 | 37.7 | 189.6 | 91 | 52 | 100% | 53.9% | 5 | 53.9% | 5 | $553 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x124e...1c58` | C | 100.0 | 40.4 | 188.8 | 128 | 114 | 70% | 22.4% | 41 | 16.9% | 59 | $1135 | existing,sports_tape | opposite_side_same_market |
| `0x7020...d655` | C | 100.0 | 43.1 | 176.3 | 349 | 231 | 63% | 60.7% | 81 | 69.4% | 101 | $288 | existing | opposite_side_same_market |

## Recent Trade Flow Scout

- Wallets: 1
- Rule: recent qualifying trade source, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | FlowScore | RecentHits | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xd581...83b7` | C | 100.0 | 35.7 | 219.9 | 2 | 6 | 27% | 173.6% | 5 | 17 | $542 | recent_trade | open_copy_exposure,opposite_side_same_market |

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
| `0xe872...819a` | B | 33.4 | $2200 | $2200 | 51 | 70.6% | +8.72 | +1.61 | +6.46 | +16.96 | 309.6 | open_copy_exposure,opposite_side_same_market |
| `0x4c9c...e006` | C | 41.9 | $800 | $800 | 12 | 75.0% | +16.81 | +4.67 | +14.65 | +30.95 | 288.0 | opposite_side_same_market |
| `0xc6f1...5729` | B | 31.7 | $6870 | $6870 | 38 | 71.1% | +10.37 | +1.18 | +3.95 | +15.97 | 253.8 | opposite_side_same_market |

## Live Edge Blocked Push Wallets

- Wallets: 6
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | B | 96.5 | 25.1 | 15m edge -10.00pp over 2 samples | 233.0% | 4 | existing | open_copy_exposure,opposite_side_same_market |
| `0xd8b5...54c4` | B | 100.0 | 25.1 | 15m edge -24.00pp over 2 samples | 60.0% | 23 | existing | open_copy_exposure,opposite_side_same_market |
| `0x124e...1c58` | C | 100.0 | 40.4 | 15m edge -25.28pp over 2 samples | 16.9% | 59 | existing,sports_tape | opposite_side_same_market |
| `0x4572...83ff` | C | 100.0 | 35.4 | 1h edge -34.83pp over 1 samples | 65.9% | 28 | existing | open_copy_exposure,opposite_side_same_market |
| `0x96d7...3a17` | C | 100.0 | 37.7 | 1h edge -35.95pp over 1 samples | 53.9% | 5 | existing | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 27.8 | 1h edge -72.59pp over 1 samples | 33.9% | 15 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | watch | - | D | 18.3 | $18862 | $18862 | 1 | 25.8% | 1 | - |
| `0x7af7...4c89` | watch | - | B | 21.7 | $15777 | $15777 | 1 | -33.4% | 1 | open_copy_exposure |
| `0x9520...fa6e` | watch | - | B | 33.6 | $14778 | $14778 | 1 | 7.1% | 25 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 68.9 | $7000 | $7000 | 1 | 105.9% | 134 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | reject-bot | - | D | 54.9 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.1 | 653.4 | 439 | 368 | 37% | 155.5% | 18 | 126.2% | 40 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | pushed | - | B | 100.0 | 33.4 | 517.3 | 1419 | 1140 | 87% | 84.4% | 24 | 82.6% | 25 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 100.0 | 68.9 | 425.5 | 124 | 122 | 6% | 81.9% | 19 | 105.9% | 134 | existing,sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xc387...5a97` | reject-bot | - | BOT | 100.0 | 69.7 | 258.4 | 862 | 778 | 72% | 27.6% | 265 | 27.5% | 347 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 43.1 | 249.7 | 136 | 126 | 19% | 37.0% | 51 | 37.5% | 229 | sports_tape | opposite_side_same_market |
| `0x4c9c...e006` | pushed | - | C | 100.0 | 41.9 | 191.6 | 284 | 231 | 98% | 22.8% | 75 | 23.1% | 77 | existing,sports_tape | opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.4 | 181.5 | 108 | 103 | 100% | 22.4% | 42 | 22.4% | 42 | existing,sports_tape | opposite_side_same_market |
| `0xfe18...3980` | pushed | - | A | 84.0 | 24.4 | 169.1 | 94 | 70 | 65% | 27.5% | 33 | 17.9% | 44 | existing,sports_tape | opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 84.8 | 152.4 | 154 | 151 | 10% | 28.5% | 49 | 41.7% | 279 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,open_copy_exposure |
| `0x124e...1c58` | blocked-edge | 15m edge -25.28pp over 2 samples | C | 100.0 | 40.4 | 150.6 | 128 | 114 | 70% | 22.4% | 41 | 16.9% | 59 | existing,sports_tape | opposite_side_same_market |
| `0xc6f1...5729` | pushed | - | B | 100.0 | 31.7 | 148.5 | 159 | 152 | 90% | 15.0% | 69 | 12.9% | 75 | existing,sports_tape | opposite_side_same_market |
| `0x20c6...45ae` | watch | - | C | 100.0 | 43.6 | 128.8 | 53 | 38 | 33% | 117.0% | 1 | 68.1% | 2 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | pushed | - | B | 100.0 | 16.0 | 115.1 | 7 | 5 | 70% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0xbc25...aaba` | prompt | - | B | 100.0 | 3.9 | 89.2 | 17 | 7 | 94% | 0.0% | 0 | 0.0% | 0 | sports_tape | open_copy_exposure |
| `0xb3c1...4837` | prompt | - | B | 100.0 | 21.2 | 81.1 | 194 | 146 | 91% | -0.1% | 21 | 0.3% | 23 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 10 | 189 | 61.7% | $+1598.62 | 78.3% | 62.3% | 33.9% | 2.51x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 2 | 11 | 214 | 59.6% | $+1722.34 | 76.6% | 60.0% | 33.9% | 2.30x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 10 | 205 | 60.2% | $+1674.87 | 76.6% | 62.3% | 33.9% | 2.16x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
