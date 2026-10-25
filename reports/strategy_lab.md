# Strategy Lab Report

**Generated:** 2026-10-25 22:54 +08

- Scores: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallet_scores.json`
- Exclude wallets: `/Users/murphyma/work/polymarket-go/db/strategy_iteration/wallets.strategy-exclude.txt` (434)
- Valid strategies found: 156
- Candidate layers: 23 core + 20 watch + 10 sports + 0 scout + 10 target + 0 flow + 0 tape
- Push wallets after live-edge blocks: 60 total
- Live-edge blocked push wallets: 3
- Leaderboard scout push enabled: true
- Tape probation wallets: 0 observation-only
- Tape edge-hot wallets: 1 observation-only

## Selected Core Strategy

- Wallets: 23
- Aggregate closed copy trades: 2355
- Aggregate copy ROI: 172.6%
- Aggregate copy PnL: $+49655.86
- Aggregate copy win rate: 95.8%
- Median wallet CopyROI: 127.8%
- Worst included wallet CopyROI: 101.1%
- Open copy cost / closed copy capital: 0.35x
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Core Wallets

| Wallet | Tier | Bot | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | ClosedROI |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `0x0f00...73ce` | A | 11.47 | 177 | 46% | 221.9% | $+7479.88 | 205 | 98.5% | 53.1% |
| `0x45e4...3c69` | A | 22.34 | 32 | 4% | 171.5% | $+5796.31 | 296 | 92.2% | 28.5% |
| `0x3724...69b5` | A | 17.01 | 1116 | 94% | 123.0% | $+4550.20 | 208 | 99.0% | 48.9% |
| `0x2b17...d36d` | A | 9.76 | 303 | 70% | 171.8% | $+3865.25 | 153 | 99.3% | 53.2% |
| `0x677d...4907` | A | 21.59 | 148 | 50% | 110.5% | $+2487.01 | 222 | 97.8% | 54.3% |
| `0x42c9...9623` | A | 4.62 | 0 | 0% | 116.5% | $+2340.71 | 152 | 100.0% | 49.3% |
| `0x075d...58c1` | A | 19.62 | 18 | 17% | 169.4% | $+1422.97 | 76 | 92.1% | 40.5% |
| `0x61b6...e14e` | A | 22.93 | 15 | 6% | 141.9% | $+950.53 | 66 | 72.7% | 18.5% |
| `0x15dd...95ba` | A | 24.93 | 216 | 51% | 104.8% | $+775.91 | 69 | 79.7% | 59.5% |
| `0x2d8a...e818` | A | 24.57 | 2 | 1% | 127.8% | $+447.15 | 29 | 86.2% | 52.1% |
| `0xbfd8...eeca` | A | 12.34 | 17 | 7% | 107.2% | $+310.87 | 29 | 100.0% | 65.7% |
| `0x768b...bf0f` | A | 2.83 | 14 | 18% | 103.3% | $+299.49 | 24 | 100.0% | 47.3% |
| `0x7bdc...9b4c` | A | 24.56 | 65 | 38% | 166.7% | $+266.69 | 16 | 93.8% | 32.0% |
| `0x1425...6e0c` | A | 23.91 | 177 | 64% | 106.3% | $+244.51 | 23 | 100.0% | 63.4% |
| `0x2c5b...1dcd` | A | 3.53 | 5 | 15% | 129.7% | $+142.61 | 11 | 100.0% | 47.7% |
| `0x17e2...d472` | B | 29.48 | 323 | 30% | 249.1% | $+12281.37 | 463 | 98.1% | 55.6% |
| `0x5be9...ef2b` | B | 27.12 | 135 | 52% | 241.2% | $+3207.39 | 133 | 100.0% | 61.9% |
| `0x7967...f4a0` | B | 29.37 | 38 | 50% | 510.6% | $+765.98 | 15 | 100.0% | 45.6% |
| `0x1bd8...e5f8` | B | 26.50 | 22 | 9% | 147.7% | $+575.82 | 39 | 74.4% | 19.3% |
| `0x162d...8944` | B | 25.15 | 3 | 1% | 125.3% | $+475.98 | 37 | 91.9% | 47.3% |
| `0x092b...614e` | B | 26.71 | 115 | 98% | 101.1% | $+384.20 | 35 | 100.0% | 55.2% |
| `0x4607...d542` | B | 26.77 | 65 | 62% | 101.1% | $+363.82 | 34 | 100.0% | 47.6% |
| `0xbe12...728f` | B | 26.41 | 22 | 26% | 105.3% | $+221.21 | 20 | 100.0% | 51.7% |

## Active Watchlist

- Wallets: 20
- Aggregate closed copy trades: 1385
- Aggregate copy ROI: 99.1%
- Aggregate copy PnL: $+16570.27
- Aggregate copy win rate: 90.1%
- Worst included CopyROI: 50.6%

| Wallet | Tier | Smart | Bot | WatchScore | TargetT | Target% | CopyROI | CopyPnL | CopyT | CopyWin | AvgNotional | Sources |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0x3e35...4f69` | A | 100.0 | 13.4 | 1257.0 | 46 | 46% | 379.1% | $+606.57 | 14 | 100.0% | $255 | existing |
| `0x43e9...a4fb` | B | 100.0 | 34.6 | 924.2 | 0 | 0% | 281.2% | $+393.62 | 13 | 100.0% | $1129 | existing |
| `0x8a43...30da` | A | 97.5 | 24.8 | 872.2 | 22 | 88% | 241.7% | $+386.75 | 13 | 92.3% | $1008 | existing |
| `0xbe7e...6d00` | B | 100.0 | 30.4 | 855.7 | 96 | 12% | 116.5% | $+2946.91 | 230 | 98.3% | $1828 | existing |
| `0x760f...326a` | B | 100.0 | 31.0 | 835.9 | 15 | 12% | 131.9% | $+1464.30 | 100 | 94.0% | $9204 | existing |
| `0x7e33...73e7` | B | 100.0 | 34.2 | 803.9 | 64 | 22% | 146.7% | $+894.95 | 60 | 98.3% | $254 | existing |
| `0x3520...1fa0` | B | 100.0 | 30.1 | 800.9 | 34 | 30% | 175.9% | $+510.03 | 28 | 96.4% | $832 | existing |
| `0x612f...05d8` | B | 100.0 | 32.1 | 661.2 | 93 | 43% | 102.6% | $+810.88 | 75 | 94.7% | $580 | existing,holder |
| `0x9a0c...370e` | B | 100.0 | 30.5 | 643.6 | 200 | 79% | 82.0% | $+1631.27 | 167 | 92.2% | $357 | existing |
| `0x586a...be99` | B | 100.0 | 35.0 | 601.6 | 49 | 84% | 127.1% | $+304.98 | 20 | 90.0% | $1689 | existing |
| `0x7d68...8239` | A | 100.0 | 10.7 | 595.0 | 20 | 22% | 94.5% | $+396.86 | 42 | 97.6% | $287 | existing |
| `0x50d0...a6fc` | A | 100.0 | 18.0 | 594.2 | 32 | 19% | 71.9% | $+955.71 | 117 | 99.2% | $2495 | existing |
| `0x2089...712e` | B | 100.0 | 28.8 | 590.9 | 141 | 40% | 92.5% | $+1221.19 | 51 | 88.2% | $6315 | existing |
| `0x0d66...ff1d` | B | 100.0 | 34.5 | 587.2 | 174 | 99% | 100.0% | $+619.91 | 42 | 88.1% | $717 | existing |
| `0x7972...594f` | B | 100.0 | 33.1 | 568.1 | 302 | 68% | 66.9% | $+1698.14 | 196 | 61.7% | $1169 | existing,holder |
| `0x2bf0...728f` | A | 100.0 | 19.8 | 559.9 | 10 | 5% | 84.0% | $+571.09 | 66 | 84.8% | $1817 | existing |
| `0x68cb...6f57` | B | 100.0 | 16.8 | 493.3 | 1 | 3% | 199.8% | $+99.91 | 4 | 100.0% | $375 | existing |
| `0x0d23...2d56` | B | 100.0 | 34.5 | 466.9 | 34 | 19% | 61.2% | $+355.12 | 57 | 98.2% | $5715 | existing |
| `0x6261...4adb` | A | 100.0 | 7.7 | 460.0 | 120 | 90% | 83.5% | $+292.35 | 12 | 91.7% | $1268 | existing |
| `0x4da5...3baf` | B | 100.0 | 31.4 | 448.4 | 153 | 100% | 50.6% | $+409.73 | 78 | 93.6% | $980 | existing |

## Sports Positive Copy

- Wallets: 10
- Rule: target-category closed copy sample, positive TargetCopyROI/PnL, low bot score, not already core/watch

| Wallet | Tier | Smart | Bot | SportsScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyPnL | TargetCopyT | TargetCopyWin | CopyROI | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x4cb0...57e0` | B | 100.0 | 34.8 | 786.5 | 772 | 332 | 92% | 104.8% | $+3469.85 | 208 | 100.0% | 102.3% | existing | burst_trading |
| `0x5135...7217` | A | 100.0 | 18.4 | 751.6 | 1178 | 314 | 100% | 97.2% | $+2381.20 | 184 | 100.0% | 97.2% | existing,holder | burst_trading |
| `0x3231...6e75` | A | 100.0 | 0.8 | 654.2 | 498 | 305 | 100% | 78.3% | $+1222.15 | 156 | 100.0% | 78.3% | existing,holder,sports_tape | - |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 571.3 | 1150 | 1112 | 100% | 61.1% | $+4435.01 | 295 | 95.6% | 61.1% | existing,recent_trade | burst_trading,opposite_side_same_market |
| `0x8987...02f3` | B | 100.0 | 34.1 | 515.6 | 278 | 115 | 30% | 77.4% | $+804.82 | 52 | 78.8% | 87.9% | existing | burst_trading,opposite_side_same_market |
| `0x3830...d2d5` | B | 100.0 | 34.1 | 490.0 | 190 | 40 | 55% | 88.6% | $+221.55 | 25 | 88.0% | 70.0% | existing | opposite_side_same_market |
| `0x3d40...5739` | B | 100.0 | 25.6 | 488.6 | 131 | 68 | 79% | 76.4% | $+290.15 | 37 | 91.9% | 79.4% | existing | opposite_side_same_market |
| `0x21cc...54bc` | B | 100.0 | 26.6 | 424.5 | 182 | 69 | 98% | 77.1% | $+131.00 | 17 | 94.1% | 77.1% | existing | open_copy_exposure,opposite_side_same_market |
| `0x789f...20ec` | B | 99.4 | 26.9 | 422.0 | 52 | 32 | 96% | 71.6% | $+171.80 | 18 | 94.4% | 71.2% | existing | opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.3 | 417.7 | 213 | 203 | 91% | 47.0% | $+469.67 | 88 | 79.5% | 49.7% | existing | opposite_side_same_market |

## Leaderboard Scout

- Wallets: 0
- Rule: profit leaderboard source, not already selected in core/watch/sports, bot below threshold, no fixed-flow/negative-copy flags; research-only by default

| Wallet | Tier | Smart | Bot | ScoutScore | TargetT | Target% | CopyROI | CopyT | Large | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|---|

No leaderboard-only wallets passed the scout filters.

## Target Category Scout

- Wallets: 10
- Rule: basketball/soccer/esports target activity, not already selected in higher-priority layers, bot below threshold, no fixed-flow/negative-copy flags

| Wallet | Tier | Smart | Bot | TargetScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | AvgNotional | Sources | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x67ba...246c` | C | 100.0 | 39.3 | 249.0 | 965 | 571 | 98% | 58.2% | 166 | 58.8% | 171 | $665 | existing | opposite_side_same_market |
| `0x4d6e...5cfd` | C | 100.0 | 36.1 | 237.9 | 1255 | 443 | 91% | 90.7% | 44 | 90.7% | 51 | $603 | existing,holder | burst_trading |
| `0x1550...5481` | C | 100.0 | 44.0 | 235.7 | 1448 | 482 | 100% | 160.8% | 28 | 160.8% | 28 | $265 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0x3fe2...93b3` | A | 100.0 | 19.8 | 233.6 | 83 | 66 | 100% | 20.9% | 24 | 20.9% | 24 | $1544 | existing,holder | opposite_side_same_market |
| `0x5f61...519a` | B | 100.0 | 27.1 | 230.5 | 437 | 259 | 100% | 37.9% | 101 | 37.9% | 101 | $234 | existing,holder | opposite_side_same_market |
| `0x249e...4fa6` | C | 100.0 | 36.4 | 225.3 | 1165 | 406 | 100% | 62.9% | 178 | 62.9% | 178 | $148 | existing | burst_trading,opposite_side_same_market |
| `0x3a01...c89e` | C | 100.0 | 43.2 | 224.4 | 1318 | 491 | 99% | 83.3% | 93 | 83.3% | 93 | $195 | existing | burst_trading |
| `0x73c2...3a29` | B | 100.0 | 33.6 | 216.0 | 504 | 314 | 87% | 12.6% | 156 | 14.4% | 174 | $167 | existing,holder | opposite_side_same_market |
| `0x257a...b42a` | B | 100.0 | 29.6 | 214.4 | 212 | 171 | 99% | 45.4% | 72 | 45.4% | 72 | $351 | existing | opposite_side_same_market |
| `0xb208...c095` | C | 100.0 | 42.0 | 210.9 | 1222 | 285 | 100% | 107.6% | 5 | 107.6% | 5 | $423 | existing,holder | burst_trading,open_copy_exposure |

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

- Wallets: 1
- Rule: recent sports-tape wallets with positive measured 5m/15m edge; observation-only and promoted only after alert ROI proves out

| Wallet | Tier | Bot | MaxBuy | BuyNotional | EdgeN | EdgeWin | AvgPP | 5mPP | 15mPP | 1hPP | Score | Risks |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| `0xc6f1...5729` | TAPE | 0.0 | $6870 | $6870 | 58 | 65.5% | +7.80 | +1.04 | +3.86 | +12.54 | 244.6 | - |

## Live Edge Blocked Push Wallets

- Wallets: 3
- Rule: wallets remain in research layers, but are removed from Telegram push when measured 15m/1h edge breaches the live block gate.

| Wallet | Tier | Smart | Bot | Reason | CopyROI | CopyT | Sources | Risks |
|---|---|---:|---:|---|---:|---:|---|---|
| `0x0f00...73ce` | A | 100.0 | 11.5 | 15m edge -10.00pp over 2 samples | 221.9% | 205 | existing | opposite_side_same_market |
| `0xf3ce...a57a` | B | 100.0 | 34.2 | 1h edge -47.28pp over 1 samples | 61.1% | 295 | existing,recent_trade | burst_trading,opposite_side_same_market |
| `0x0e24...7014` | B | 100.0 | 30.3 | 1h edge -72.59pp over 1 samples | 49.7% | 94 | existing | opposite_side_same_market |

## Direct Sports Tape Gate Review

- Rule: 5000+ single-buy sports tape wallets are reviewed for direct whale push; scored BOT/bot-like flow stays blocked.

| Wallet | Status | EdgeBlock | Tier | Bot | MaxBuy | BuyNotional | Buys | CopyROI | CopyT | Risks |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `0xcbc6...ec62` | reject-bot | - | D | 50.3 | $294000 | $294000 | 1 | 112.5% | 4 | bot_like_flow,burst_trading |
| `0x0df2...5d51` | reject-bot | - | BOT | 62.6 | $26101 | $26101 | 1 | 101.6% | 26 | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x074a...f0e0` | reject-bot | - | BOT | 77.6 | $12804 | $12804 | 1 | 74.7% | 54 | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xcf0c...f2bb` | reject-bot | - | D | 51.9 | $12804 | $12804 | 1 | 95.0% | 28 | bot_like_flow,burst_trading |

## Sports Tape Candidate Review

- Wallets: 9
- Rule: latest target-category large-order wallets, reviewed after full scoring; pushed only if they pass strategy layer filters

| Wallet | Status | EdgeBlock | Tier | Smart | Bot | TapeScore | TargetT | TargetLarge | Target% | TargetCopyROI | TargetCopyT | CopyROI | CopyT | Sources | Risks |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| `0x1550...5481` | pushed | - | C | 100.0 | 44.0 | 827.6 | 1448 | 482 | 100% | 160.8% | 28 | 160.8% | 28 | existing,holder,sports_tape | burst_trading,open_copy_exposure |
| `0x3231...6e75` | pushed | - | A | 100.0 | 0.8 | 687.5 | 498 | 305 | 100% | 78.3% | 156 | 78.3% | 156 | existing,holder,sports_tape | - |
| `0xcbc6...ec62` | reject | - | D | 100.0 | 50.3 | 291.8 | 1155 | 304 | 78% | 112.5% | 4 | 112.5% | 4 | existing,sports_tape | bot_like_flow,burst_trading |
| `0x0df2...5d51` | reject-bot | - | BOT | 100.0 | 62.6 | 291.5 | 122 | 60 | 44% | 70.7% | 10 | 101.6% | 26 | existing,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xcf0c...f2bb` | reject | - | D | 100.0 | 51.9 | 270.6 | 73 | 27 | 5% | 75.8% | 6 | 95.0% | 28 | existing,holder,sports_tape | bot_like_flow,burst_trading |
| `0x074a...f0e0` | reject-bot | - | BOT | 100.0 | 77.6 | 265.8 | 1176 | 107 | 84% | 53.2% | 24 | 74.7% | 54 | existing,holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0x2d8a...e818` | pushed | - | A | 100.0 | 24.6 | 215.0 | 2 | 1 | 1% | 0.0% | 0 | 127.8% | 29 | existing,sports_tape | opposite_side_same_market |
| `0x9697...18e9` | reject-bot | - | BOT | 100.0 | 72.1 | 171.9 | 1459 | 417 | 100% | 28.8% | 50 | 28.6% | 51 | existing,holder,sports_tape | bot_like_flow,burst_trading,opposite_side_same_market |
| `0xe907...cff6` | blocked-edge | 15m edge -1.55pp over 26 samples | BOT | 97.4 | 78.0 | -69.3 | 475 | 55 | 37% | -9.0% | 16 | -0.7% | 34 | existing,holder,sports_tape | bot_like_flow,burst_trading,negative_copy_sim,opposite_side_same_market |

## Top Strategies

| Rank | Wallets | CopyT | CopyROI | CopyPnL | CopyWin | MedianROI | WorstROI | Open/Cap | Params |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 23 | 2355 | 172.6% | $+49655.86 | 95.8% | 127.8% | 101.1% | 0.35x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 2 | 22 | 2344 | 172.8% | $+49513.25 | 95.8% | 126.5% | 101.1% | 0.35x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 3 | 20 | 2181 | 175.6% | $+47353.60 | 97.4% | 126.5% | 101.1% | 0.30x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 4 | 21 | 2250 | 173.7% | $+48129.51 | 96.8% | 125.3% | 101.1% | 0.29x | tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 5 | 19 | 2170 | 175.8% | $+47210.99 | 97.4% | 125.3% | 101.1% | 0.30x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=80 closedROI>=0 smart>=70 |
| 6 | 20 | 2239 | 173.9% | $+47986.90 | 96.8% | 124.1% | 101.1% | 0.29x | tier>=B bot<30 copyT>=15 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=20 smart>=70 |
| 7 | 20 | 2313 | 171.0% | $+48480.58 | 95.8% | 124.1% | 101.1% | 0.33x | tier>=B bot<30 copyT>=20 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 8 | 40 | 3203 | 149.1% | $+58662.65 | 95.7% | 104.1% | 63.1% | 0.45x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 9 | 39 | 3194 | 149.3% | $+58599.53 | 95.7% | 104.8% | 64.3% | 0.45x | tier>=B bot<30 copyT>=10 copyROI>=60 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 10 | 38 | 3182 | 149.6% | $+58520.56 | 95.8% | 105.1% | 64.3% | 0.45x | tier>=B bot<30 copyT>=8 copyROI>=60 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
| 11 | 44 | 3325 | 145.7% | $+59283.00 | 95.4% | 101.1% | 41.2% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70 |
| 12 | 39 | 3254 | 147.5% | $+58888.46 | 95.4% | 104.8% | 45.4% | 0.46x | tier>=B bot<30 copyT>=8 copyROI>=40 copyPnL>=100 copyWin>=60 closedROI>=0 smart>=70 |
