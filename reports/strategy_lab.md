# Strategy Lab Report

**Generated:** 2026-07-09 13:46 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (7)
- Valid strategies found: 8
- Candidate layers: 10 core + 20 watch + 7 sports + 1 scout + 6 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 40 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: false
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 3 observation-only

## Selected Core Strategy

- Wallets: 10
- Aggregate closed copy trades: 182
- Aggregate copy ROI: 65.9%
- Aggregate copy PnL: $+1442.89
- Aggregate copy win rate: 75.3%
- Median wallet CopyROI: 64.3%
- Worst included wallet CopyROI: 42.0%
- Open copy cost / closed copy capital: 1.65x
- Params: tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x7992...1fc1` | A | 20.07 | 21 | 29% | 73.7% | $+184.35 | 20 | 75.0% | 23.4% |
| `0xe6bf...c536` | A | 23.11 | 70 | 95% | 58.4% | $+175.23 | 20 | 85.0% | 20.2% |
| `0x68cb...6f57` | A | 21.60 | 50 | 60% | 56.1% | $+151.58 | 24 | 66.7% | 8.4% |
| `0x9ecc...7850` | A | 22.33 | 4 | 5% | 50.8% | $+132.13 | 26 | 73.1% | 35.9% |
| `0xa451...449e` | A | 24.48 | 32 | 55% | 108.4% | $+130.13 | 12 | 83.3% | 29.3% |
| `0xdbd5...cba7` | A | 18.49 | 85 | 99% | 42.0% | $+117.75 | 16 | 93.8% | 27.3% |
| `0xb715...c3bb` | B | 25.24 | 56 | 89% | 135.2% | $+148.75 | 10 | 80.0% | 47.4% |
| `0xd8b5...54c4` | B | 25.60 | 181 | 99% | 73.1% | $+146.13 | 18 | 66.7% | 54.2% |
| `0x0cd6...5c6c` | B | 28.26 | 49 | 45% | 61.7% | $+129.49 | 17 | 70.6% | 4.2% |
| `0x0931...e78e` | B | 26.67 | 101 | 66% | 67.0% | $+127.35 | 19 | 68.4% | 24.2% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 447
- Aggregate copy ROI: 58.7%
- Aggregate copy PnL: $+3253.97
- Aggregate copy win rate: 77.6%
- Worst included CopyROI: 15.8%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x552b...f8bc` | B | 100.0 | 32.3 | 906.7 | 45 | 63% | 258.6% | $+413.78 | 14 | 64.3% | $400 | existing |
| `0x0f00...73ce` | B | 100.0 | 25.2 | 586.6 | 78 | 76% | 233.0% | $+163.12 | 4 | 100.0% | $2290 | existing |
| `0xe872...819a` | B | 100.0 | 33.5 | 507.3 | 1426 | 88% | 87.2% | $+427.04 | 30 | 90.0% | $1227 | existing,sports_tape |
| `0xd5b1...1b71` | B | 100.0 | 32.6 | 495.1 | 109 | 98% | 84.5% | $+244.89 | 28 | 92.9% | $566 | existing |
| `0xb36f...53d0` | B | 100.0 | 30.6 | 496.4 | 581 | 92% | 57.6% | $+754.29 | 121 | 86.8% | $280 | existing,sports_tape |
| `0x8bb0...ca79` | B | 100.0 | 33.6 | 474.2 | 65 | 28% | 98.6% | $+246.44 | 16 | 75.0% | $1498 | existing |
| `0x7c73...6ee3` | B | 100.0 | 26.0 | 437.5 | 42 | 95% | 120.8% | $+60.39 | 5 | 100.0% | $5348 | existing |
| `0x4572...83ff` | B | 100.0 | 33.9 | 398.0 | 143 | 99% | 61.7% | $+240.50 | 21 | 90.5% | $1135 | existing |
| `0xc419...a8b0` | B | 100.0 | 29.0 | 369.4 | 34 | 85% | 85.9% | $+60.12 | 6 | 83.3% | $1026 | existing |
| `0x42db...d512` | B | 100.0 | 30.3 | 318.9 | 104 | 100% | 26.6% | $+117.02 | 40 | 57.5% | $1190 | existing,sports_tape |
| `0x96d7...3a17` | B | 100.0 | 15.2 | 307.2 | 45 | 100% | 60.2% | $+18.06 | 3 | 100.0% | $543 | existing |
| `0x2393...1a8d` | A | 100.0 | 22.4 | 305.6 | 32 | 63% | 39.0% | $+42.94 | 11 | 63.6% | $609 | sports_tape |
| `0x5c7a...fdbe` | B | 100.0 | 27.8 | 294.7 | 31 | 30% | 43.2% | $+47.47 | 9 | 77.8% | $292 | existing |
| `0x8a35...904a` | A | 100.0 | 21.8 | 291.4 | 0 | 0% | 57.6% | $+74.95 | 5 | 100.0% | $1143 | existing |
| `0x0e24...7014` | B | 100.0 | 26.6 | 289.7 | 25 | 100% | 35.4% | $+46.04 | 6 | 66.7% | $894 | existing |
| `0xa8b8...36dd` | A | 100.0 | 20.5 | 284.7 | 65 | 55% | 38.1% | $+19.06 | 5 | 100.0% | $487 | existing |
| `0x8e17...9d3c` | B | 100.0 | 34.0 | 277.9 | 160 | 61% | 15.8% | $+148.39 | 84 | 57.1% | $1342 | existing |
| `0x1c62...5805` | B | 100.0 | 33.9 | 274.4 | 25 | 57% | 35.7% | $+28.53 | 8 | 75.0% | $351 | existing |
| `0x3c14...a20a` | B | 100.0 | 32.8 | 255.7 | 10 | 10% | 23.2% | $+64.92 | 26 | 92.3% | $423 | existing |
| `0xe4b1...a75e` | B | 100.0 | 29.8 | 253.5 | 30 | 41% | 22.5% | $+36.02 | 5 | 60.0% | $2777 | existing |

## Sports Positive Copy

- Wallets: 7
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xeec5...1862` | B | 100.0 | 32.0 | 620.7 | 55 | 55 | 86% | 205.2% | $+246.27 | 8 | 50.0% | 154.8% | existing | open_copy_exposure,opposite_side_same_market |
| `0xb907...2c9e` | B | 100.0 | 33.7 | 369.0 | 67 | 43 | 33% | 84.9% | $+93.37 | 8 | 100.0% | 95.7% | existing | open_copy_exposure,opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 30.8 | 341.6 | 38 | 17 | 100% | 75.4% | $+45.23 | 6 | 100.0% | 75.4% | existing | opposite_side_same_market |
| `0xfe18...3980` | A | 100.0 | 24.9 | 303.5 | 77 | 55 | 63% | 36.9% | $+99.72 | 25 | 24.0% | 22.5% | existing,sports_tape | opposite_side_same_market |
| `0x631a...a80f` | B | 100.0 | 32.9 | 253.8 | 85 | 40 | 72% | 31.6% | $+34.74 | 11 | 81.8% | 38.4% | existing | open_copy_exposure,opposite_side_same_market |
| `0x7020...fd20` | B | 100.0 | 32.5 | 244.0 | 52 | 46 | 43% | 34.6% | $+41.48 | 7 | 85.7% | 11.9% | existing | open_copy_exposure,opposite_side_same_market |
| `0xc6f1...5729` | B | 100.0 | 29.8 | 231.3 | 103 | 98 | 85% | 11.8% | $+53.23 | 43 | 65.1% | 9.0% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 1
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb5...ef29` | C | 100.0 | 44.4 | 177.6 | 0 | 0% | 0.0% | 0 | 245 | $4587 | existing,leaderboard_profit_30d,leaderboard_profit_7d | burst_trading,open_copy_exposure |

## Target Category Scout

- Wallets: 6
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x8d85...c436` | C | 100.0 | 35.0 | 351.0 | 1741 | 1369 | 100% | 22.5% | 291 | 22.5% | 291 | $456 | existing | burst_trading,opposite_side_same_market |
| `0xaf6a...925f` | C | 100.0 | 38.3 | 195.4 | 90 | 82 | 100% | 93.4% | 20 | 93.4% | 20 | $894 | existing | open_copy_exposure,opposite_side_same_market |
| `0x2391...554f` | C | 100.0 | 40.6 | 192.3 | 70 | 68 | 99% | 88.7% | 23 | 88.7% | 23 | $1479 | existing,retain | opposite_side_same_market |
| `0x124e...1c58` | C | 100.0 | 40.5 | 185.1 | 108 | 96 | 70% | 14.0% | 35 | 10.6% | 49 | $928 | existing,sports_tape | opposite_side_same_market |
| `0x4c9c...e006` | C | 100.0 | 43.4 | 176.9 | 111 | 85 | 97% | 23.2% | 30 | 24.9% | 31 | $292 | existing | opposite_side_same_market |
| `0x7020...d655` | C | 100.0 | 42.7 | 172.5 | 308 | 202 | 61% | 62.0% | 69 | 72.0% | 87 | $275 | existing | opposite_side_same_market |

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
| `0xb36f...53d0` | B | 30.6 | $551 | $551 | 4 | 100.0% | +26.13 | +8.50 | +26.50 | +34.95 | 442.5 | opposite_side_same_market |
| `0xe872...819a` | B | 33.5 | $650 | $650 | 4 | 100.0% | +21.99 | +8.50 | +11.00 | +47.95 | 440.3 | open_copy_exposure,opposite_side_same_market |
| `0x124e...1c58` | C | 40.5 | $4366 | $5523 | 4 | 100.0% | +10.55 | +3.74 | +4.94 | +17.44 | 247.7 | opposite_side_same_market |

## Live Edge Blocked Push Wallets

- Wallets: 3
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0xd8b5...54c4` | B | 100.0 | 25.6 | 15m edge -24.00pp over 2 samples | 73.1% | 18 | existing | open_copy_exposure,opposite_side_same_market |
| `0x4572...83ff` | B | 100.0 | 33.9 | 1h edge -34.83pp over 1 samples | 61.7% | 21 | existing | open_copy_exposure,opposite_side_same_market |
| `0x96d7...3a17` | B | 100.0 | 15.2 | 1h edge -35.95pp over 1 samples | 60.2% | 3 | existing | open_copy_exposure,opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xe0cc...aafc` | reject-flow | - | D | 21.7 | $18862 | $18862 | 1 | 0.0% | 0 | fixed_price |
| `0x7af7...4c89` | reject-flow | - | D | 23.8 | $15777 | $15777 | 1 | -33.4% | 1 | fixed_price,open_copy_exposure |
| `0x9520...fa6e` | watch | - | B | 33.8 | $14778 | $14778 | 1 | 7.1% | 25 | open_copy_exposure,opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 62.0 | $7000 | $7000 | 1 | 116.3% | 110 | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x204f...5e14` | reject-bot | - | D | 52.9 | $5209 | $14402 | 15 | 0.0% | 0 | bot_like_flow,burst_trading,open_copy_exposure |

## Sports Tape Candidate Review

- Wallets: 15
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0xbe7e...6d00` | reject-bot | - | BOT | 100.0 | 64.1 | 652.0 | 433 | 369 | 36% | 152.4% | 19 | 125.3% | 41 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xe872...819a` | pushed | - | B | 100.0 | 33.5 | 561.1 | 1426 | 1156 | 88% | 88.9% | 29 | 87.2% | 30 | existing,sports_tape | open_copy_exposure,opposite_side_same_market |
| `0xb36f...53d0` | pushed | - | B | 100.0 | 30.6 | 462.6 | 581 | 348 | 92% | 57.7% | 109 | 57.6% | 121 | existing,sports_tape | opposite_side_same_market |
| `0xc387...5a97` | reject | - | D | 100.0 | 52.5 | 281.6 | 840 | 757 | 75% | 27.1% | 259 | 26.3% | 337 | sports_tape | bot_like_flow,opposite_side_same_market |
| `0xe790...698c` | watch | - | C | 100.0 | 43.2 | 234.4 | 122 | 112 | 18% | 35.6% | 47 | 35.1% | 219 | sports_tape | opposite_side_same_market |
| `0xfe18...3980` | pushed | - | A | 100.0 | 24.9 | 212.4 | 77 | 55 | 63% | 36.9% | 25 | 22.5% | 36 | existing,sports_tape | opposite_side_same_market |
| `0x42db...d512` | pushed | - | B | 100.0 | 30.3 | 201.6 | 104 | 100 | 100% | 26.6% | 40 | 26.6% | 40 | existing,sports_tape | opposite_side_same_market |
| `0xa4fd...748b` | reject-bot | - | BOT | 100.0 | 62.0 | 162.9 | 28 | 28 | 1% | -29.3% | 2 | 116.3% | 110 | sports_tape | bot_like_flow,burst_trading,open_copy_exposure,opposite_side_same_market |
| `0x6982...a165` | reject-bot | - | BOT | 100.0 | 84.6 | 151.8 | 148 | 145 | 9% | 28.4% | 47 | 42.7% | 265 | sports_tape | bot_like_flow,burst_trading,extreme_price_heavy,open_copy_exposure |
| `0xcc87...521a` | prompt | - | B | 100.0 | 25.0 | 147.5 | 11 | 7 | 55% | 97.8% | 1 | 97.8% | 1 | sports_tape | extreme_price_heavy,open_copy_exposure |
| `0x20c6...45ae` | watch | - | C | 100.0 | 44.3 | 127.2 | 49 | 34 | 33% | 117.0% | 1 | 68.1% | 2 | sports_tape | burst_trading,open_copy_exposure,opposite_side_same_market |
| `0xbeea...0c83` | prompt | - | B | 100.0 | 18.2 | 109.8 | 5 | 3 | 62% | 27.8% | 1 | 46.4% | 3 | existing,sports_tape | opposite_side_same_market |
| `0x124e...1c58` | pushed | - | C | 100.0 | 40.5 | 104.0 | 108 | 96 | 70% | 14.0% | 35 | 10.6% | 49 | existing,sports_tape | opposite_side_same_market |
| `0x9438...ad81` | reject-bot | - | BOT | 100.0 | 64.9 | 85.2 | 1874 | 480 | 99% | 9.0% | 73 | 8.8% | 75 | sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xbc25...aaba` | prompt | - | B | 97.7 | 10.9 | 82.9 | 16 | 6 | 89% | 0.0% | 0 | 37.5% | 1 | sports_tape | open_copy_exposure,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 10 | 182 | 65.9% | $+1442.89 | 75.3% | 64.3% | 42.0% | 1.65x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 11 | 193 | 64.6% | $+1485.83 | 74.6% | 61.7% | 39.0% | 1.68x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 12 | 202 | 63.6% | $+1533.30 | 74.8% | 60.0% | 39.0% | 1.86x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 4 | 11 | 191 | 64.8% | $+1490.36 | 75.4% | 61.7% | 42.0% | 1.84x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 5 | 10 | 176 | 64.9% | $+1356.34 | 75.0% | 62.7% | 39.0% | 1.75x | tier>=B bot<30 copyT>=10 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 6 | 10 | 174 | 65.1% | $+1360.87 | 75.9% | 62.7% | 42.0% | 1.92x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
| 7 | 10 | 161 | 64.9% | $+1252.23 | 76.4% | 62.7% | 39.0% | 2.14x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=10 smart>=70 |
| 8 | 11 | 185 | 63.8% | $+1403.81 | 75.1% | 58.4% | 39.0% | 1.94x | tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=25 copyWin>=60 closedROI>=5 smart>=70 |
