# BTC 策略优化 TODO

> 优先级: P0 紧急 / P1 本周 / P2 下周 / P3 backlog
> 状态: ✅ 完成 / 🛠 进行中 / 💤 未开始

---

## 一、数据源扩展

### 1. [P0] ✅ 5min K线接入
- Binance BTCUSDT 5min candles（免费无 key）
- 每次 scan 自动 fetch 1000 根 + 存 btc.db（interval='5m'）
- done: commit `2bcbe3f`

### 2. [P0] ✅ 15min K线接入
- Binance BTCUSDT 15min candles
- 每次 scan 自动 fetch 1000 根 + 存 btc.db（interval='15m'）
- done: commit `2bcbe3f`

### 3. [P1] ✅ Fear & Greed Index 接入
- alternative.me API（免费）
- 极度恐慌 (<25) 时 dip 市场被高估概率更大
- 极度贪婪 (>75) 时 reach 市场被低估概率更大
- SentimentModifier: fear dampens reach BUY_YES, amplifies dip BUY_NO
- 持久化到 btc_sentiment 表
- done: commit `141ef3d`

### 4. [P1] ✅ Funding Rate 接入
- Binance perpetual BTCUSDT funding rate（8h 周期）
- 正 funding >0.05% = 多头拥挤 → dampen reach BUY_YES
- 负 funding <-0.05% = 空头拥挤 → dampen dip BUY_YES
- 集成到 SentimentModifier 统一输出乘数
- done: commit `141ef3d`

### 5. [P2] ✅ 机构资金流（Institutional Flow Proxy）
- ETF API 需付费 → 改用 Binance 免费 API 做 proxy
- `institutional.go`: Open Interest + Long/Short ratio + Futures Premium
- FlowScore (-1~+1): L/S<0.85 contrarian bullish(+0.4), premium>0.05% bullish(+0.3)
- InstitutionalModifier: 调整 BS gap 信号权重（bullish→amplify reach, dampen dip）
- 持久化到 btc_institutional 表
- 首次结果: OI=96K L/S=0.76(crowd short) → BULLISH(0.40) inst_mod=0.96
- done: commit `77f7ff1`

### 6. [P2] ✅ 链上指标监控（On-Chain Metrics）
- `onchain.go`: Blockchair API → mempool/fees/hashrate/difficulty/block height
- OnChainScore (-1~+1): 高mempool+高fee=活跃(bullish), 低活动=bearish
- ExchangeNetFlow 启发式: 高mempool+高fee=INFLOW, 低mempool=OUTFLOW
- OnChainModifier: 调整 BS gap 信号权重
- 持久化到 btc_onchain 表
- 首次结果: mempool=860 fee=$238 hashrate=1.03e21 OUTFLOW NEUTRAL(0.15)
- done: commit `559dcee` + `f35b665`

### 7. [P3] ✅ 宏观事件日历
- `macro.go`: 2026 全年 FOMC/CPI/NFP/PCE 日期硬编码（免费 API 已下线）
- 事件前 24h: vol 线性上调 — FOMC 30% / CPI 25% / NFP 20% / PCE 15%
- 事件后 2h: vol cooldown 15%（入场窗口）
- MacroVolAdjust 集成到 scanOnceWithState，sigma = macro × blended
- scan_done 日志新增 sigma_macro + macro label
- 首次结果: phase=normal, next=PCE 2026-04-30 (87h), vol_mult=1.00
- 6 单测覆盖
- done: commit `3b9df2b`

---

## 二、模型优化

### 8. [P0] ✅ 多时间尺度马尔科夫
- `multi_tf.go`: PredictMultiTF 加权共识（5m=20%, 15m=30%, 1h=50%）
- Alignment 检测: ALIGNED_BULL / ALIGNED_BEAR / MIXED
- Confidence 评分 + entry filter（只在强反向时阻止）
- done: commit `2bcbe3f` + `f85a577`

### 9. [P1] ✅ Hidden Markov Model (HMM)
- 3 隐状态: TREND / MEAN_REVERT / VOLATILE
- Baum-Welch (EM) 训练，Viterbi 解码
- `hmm.go`: forward/backward/Viterbi + CandlesToObservations
- VOLATILE regime 过滤（一致性最差）
- 7 个单测覆盖
- done: commit `2e0052c`

### 10. [P1] ✅ 动态波动率 (EWMA → Blended)
- EWMA λ=0.94 替换固定窗口历史波动率
- 90d hist=37.7% vs ewma=16.7%（近期 BTC 波动低）
- **问题**: 纯 EWMA 低估尾部风险，$55K dip 产生 -48.8pp 虚假信号
- **修复**: BlendedVolatility (60% EWMA + 40% hist) + 25% vol floor
- 效果: 信号从 17→14 个，gap 缩小 ~40%（$55K: -48.8→-30.6pp）
- done: commit `2e0052c` (EWMA) + `1a08f92` (blended)

### 11. [P1] ✅ 波动率微笑 (Vol Smile)
- `VolSmileAdjust`: 按 log-moneyness 线性调整
- 远离 ATM 的 strike 用更高 implied vol
- FindBSGaps 已应用 per-strike 调整
- done: commit `2e0052c`

### 12. [P2] ✅ 二阶/三阶马尔科夫
- ✅ markov2.go: 225 pair states (s[t-1],s[t])→s[t+1], adaptive w2 weight
- ✅ BlendedPrediction: 1st+2nd order merged (w2 up to 60% when ≥50 obs)
- ✅ multi_tf.go 已切换到 BlendedPrediction
- ✅ 8 新测试全部通过
- 三阶暂不实现（数据量不足以支撑 3375 个三元组状态）

### 13. [P2] ✅ 自动再训练 (KL Drift Detection)
- `retrain.go`: KLDivergence, SymmetricKL, MatrixDrift, Matrix2Drift, CheckDrift
- 7天滚动窗口 vs 全量历史比较（1st + 2nd order 同时检测）
- Symmetric KL > 0.15 触发 DRIFT_ALERT 日志告警
- 集成到 scan 循环：每 6 次 scan (~6h) 自动运行
- 6 单测覆盖 (identical/different distributions, same/different windows)
- done: commit `95983a5`

---

## 三、入场/出场/仓位

### 14. [P0] ✅ 入场择时优化
- MultiTFEntryFilter 联合 BS gap + 多时间尺度方向
- ALIGNED_BEAR 且 confidence>0.55 时阻止 BUY_YES，反之亦然
- MIXED 全部放行（BS gap 是结构性 edge，短线无方向时不阻止）
- done: commit `f85a577`

### 15. [P1] ✅ 出场策略
- `exit.go`: btc_positions 表 + CheckExits 每轮扫描检查
- 三出场条件: gap<3pp(alpha 耗尽) / BTC±5% 止损 / 7 天 timeout
- `scanOnceWithState` 返回市场状态供 exit checker 使用
- `RunStrategyWithExit` 封装 entry+exit 完整循环
- 信号触发即自动记录仓位, 每轮 scan 自动检查退出
- 6 单测覆盖全部退出路径 + hold 场景
- done: commit `e0e4c5c`

### 16. [P1] ✅ 仓位管理 (Kelly Criterion)
- `kelly.go`: KellyFraction (half-Kelly), KellySizeUSD, ValueEdge, ExpectedValue
- `kelly_test.go`: 4 个单测覆盖
- updown.go 集成: Kelly 动态 sizing 替换固定 SizeUSD
- bankroll=$90, fair=0.50, maxBet=SizeUSD*3
- EV per dollar 写入日志
- done: 本次 commit

### 17. [P2] ✅ PM 盘口深度分析 (CLOB Orderbook)
- `orderbook.go`: CLOB API → per-market bid/ask depth + spread + mid price
- DepthScore (0-1): $500+ depth=1.0, <$10=0.2, wide spread penalty
- DepthModifier: illiquid(<0.3) → 0.5x, thin(<0.5) → 0.75x, deep → 1.0x
- PMMarket 新增 ClobTokenIDs 字段（从 gamma API 传播）
- 持久化到 btc_orderbook 表
- 首次结果: 34 markets, avg_depth=0.56, min_depth=0.00
- done: commit `559dcee`

---

## 四、数据追踪与分析

### 18. [P0] ✅ PM 价格 delta 追踪
- `pm_btc_deltas` 表每次 scan 记录 PM 价格 + BTC spot
- 积累数据后做 PM 调价速度回归分析
- done: commit `2bcbe3f`
- TODO: 积累 7 天后写分析查询

### 19. [P1] 🛠 PnL 归因分析
- `cmd/backtest -mode=btc-pnl`: 读 btc.db 输出实盘 PnL 报告
- 含 PM 定价效率分析（deviation 分布图）
- 每小时自动 logPnLSummary（累计胜率/PnL/快照数）
- 解盘 log 增强: EV/Kelly/candle range/body 全记录
- **关键发现: 1h Up/Down PM 定价高效（84% 偏差<0.5pp），几乎无 edge**
- 下一步: 按市场/时段/regime 分解，需更多数据
- done: btcpnl.go + resolution 增强 commit `60cde2b`

### 20. [P1] 🛠 历史 PM 价格补全
- `updown_prices` 表: 每次 scan 记录 Up/Down 价格 + spread + deviation
- 每 30min scan 自动采集（随 daemon 运行）
- 积累 7 天后做 PM 价格变动 vs BTC 价格变动的回归分析
- done: updown_prices 表 + logUpDownPrices 本次 commit

### 21. [P2] bonereaper 策略逆向
- 拿到钱包地址后拉全部交易
- 分析: 入场时机（BTC 处于什么状态时买入）
- 分析: 偏好 strike（总是买 reach 还是 dip）
- 分析: 持仓周期（短线翻转还是长期持有）
- 复制其信号加入 whale tracker

---

## 五、基建与回测

### 22. [P1] ✅ 回测引擎增强
- btc-updown 回测模式：Sharpe/Calmar/hour/regime 分析
- PM 手续费 2% (200bp) 扣除
- 按 HMM regime 拆分胜率
- 按 UTC 小时拆分胜率
- 关键发现：30d vs 90d 小时模式完全不一致，不可靠
- 关键发现：BTC 1h 涨跌 = 掷硬币 (49.7%)，无小时偏差
- 关键发现：简单 Markov 无法稳定 >54% 胜率
- **策略转向：价值投注（买 PM 定价偏低的一侧）**
- done: commit `2e0052c`

### 23. [P2] ✅ 多币种扩展 (ETH/SOL)
- `multicoin.go`: CoinConfig + ScanCoinOnce + coin_candles/coin_pm_prices 表
- 重构 FetchBTCMarkets 共享 Gamma API 通用 fetcher
- 修复 parseStrikeFromQuestion 支持小额 strike ($80 SOL)
- ETH: 16 markets, 9 gaps, top: $1500 dip -18.1pp
- SOL: 18 markets, 10 gaps, top: $40 dip -27.9pp (66/SIGNAL!)
- 每轮 BTC scan 后自动扫 ETH+SOL
- 4 新测试覆盖
- done: commit `fd79ade`

### 24. [P2] ✅ Regime Detection 自动切换
- `regime.go`: RegimeDirectionBias 按 HMM regime + multi-TF alignment 调信号权重
- TREND: 顺势放大(1.3x) + 逆势压缩(0.77x)
- MEAN_REVERT: 逆势放大(1.2x) + 顺势中性
- VOLATILE: 全局压缩(0.5-1.0x)
- 集成到 scanOnce, 每个 signal 带 regime_bias 字段
- 5 单测覆盖
- done: commit `220deaa`

### 25. [P2] ✅ 信号质量评分系统
- `scoring.go`: ScoreSignal 综合 5 维度 → 0-100 分
  - GapScore (0-35): BS gap 大小
  - SentimentScore (0-15): F&G + funding rate 对齐度
  - RegimeScore (0-20): HMM regime 支持度
  - TFScore (0-15): 多时间尺度一致性
  - EdgeScore (0-15): 相对边际(gap/price)
- 三档: AUTO(>80) / SIGNAL(60-80) / LOG(<60)
- Signal struct 新增 Score 字段, 日志输出 score/tier
- 4 单测覆盖
- done: commit `220deaa`

---

## 进度追踪

| 日期 | 完成项 | 发现 |
|------|-------|------|
| Apr 26 | 基础 Markov + BS + PM tracker + 1h 回测 | Markov 49.9% 掷硬币; BS gap 10-23pp 有价值 |
| Apr 26 | BTC live strategy 上线 (1h scan, gap>7pp) | 首次扫描: $50K dip PM=42.5% BS=19% gap=-23pp |
| Apr 26 | 5m/15m K线 + multi-TF Markov + PM delta tracking | MIXED(bull=2.7%,bear=2.9%) 放行全部 BS gap 信号 |
| Apr 27 | F&G Index + Funding Rate 接入 + SentimentModifier | F&G=33(Fear), FR≈0, sent_mod=1.0（中性区间无调整） |
| Apr 27 | **砍长期盘 → 专注 1h Up/Down**; 修 BullProb (3%→51%); confidence 0.52→0.40; 候选窗 1-4h; max 40单/天; PM edge filter ≤0.52 | 首发 3 笔 Up 信号 · Markov 终于能出手了 |
| Apr 27 | EWMA vol (λ=0.94) + Vol Smile + HMM 3态 regime + btc-updown 回测 | hist=52%/ewma=17%·BTC 1h=掷硬币·策略转向价值投注 |
| Apr 27 | **策略转向**: 从方向预测→PM 定价偏差套利; 买 PM 定价<0.49 的一侧 | Markov 做 tiebreaker 而非主驱动 |
| Apr 27 | Kelly Criterion sizing + updown_prices 采集 + 4 单测 | half-Kelly 动态仓位; PM 价格分布数据积累中 |
| Apr 27 | PnL 报告 + 定价效率分析 + 动量触发 + 解盘增强 | **PM 1h Up/Down 定价高效 (84%<0.5pp)**; 需找低效市场 |
| Apr 27 | 启用 btc_enabled 价格级策略; 128 snapshots 分析 Up/Down 死局 | **Up/Down 最大偏差 0.5pp，无法覆盖 2% fee**; 价格级策略首扫: $55K/-48.8pp $50K/-41.4pp $45K/-33.3pp 全 BUY_NO; ⚠️ EWMA=16.7% vs hist=37.7% 尾部风险低估 |
| Apr 27 | **BlendedVol** (60%EWMA+40%hist+25%floor) + bot-daemon.sh 修复 btc_enabled | gap 缩小 40%: $55K -48.8→-30.6pp; 修复 cron-poke 重启不带 btc_enabled 导致双 daemon 409 冲突 |
| Apr 27 | **#15 Exit Strategy**: exit.go + 3 退出条件 + 仓位追踪 + 6 单测 | gap<3pp 平仓 / BTC±5% 止损 / 7d timeout; scanOnceWithState 返回完整市场状态 |

| Apr 27 | **#24 Regime Detection** + **#25 Signal Scoring** | HMM=MEAN_REVERT(100%); 信号评分: K=69/SIGNAL K=69/SIGNAL K=67/SIGNAL; regime_bias=1.0 (MR+bull+dip=中性) |
| Apr 27 | **#12 二阶马尔科夫** — markov2.go + BlendedPrediction + multi_tf 集成 | 225 pair states; w2 adaptive (≥50obs→60%); multi_tf 已切 ALIGNED_BULL(0.51); 三阶暂缓(数据不够3375态) |
| Apr 27 | **#13 自动再训练** — retrain.go + KL divergence drift detection | 7d滚动vs全量; SymmetricKL>0.15告警; 每6h检查1st+2nd order drift; 6单测 |
| Apr 27 | **#5 机构资金流** — institutional.go + OI/L-S/FuturesPrem proxy | OI=96K L/S=0.76(crowd short)→BULLISH(0.40); inst_mod=0.96 for dip BUY_NO |
| Apr 27 | **#6 链上指标** + **#17 盘口深度** — onchain.go + orderbook.go | mempool=860 fee=$238 NEUTRAL(0.15); 34 markets avg_depth=0.56; 信号升至 73/SIGNAL (regime_bias=1.12) |
| Apr 27 | **#23 多币种 ETH/SOL** — multicoin.go + parseStrike 修复 + 通用 Gamma fetcher | ETH 16mkt 9gaps top $1500/-18.1pp; SOL 18mkt 10gaps top $40/-27.9pp(66/SIGNAL); 信号宇宙 3x 扩展 |
| Apr 27 | **#7 宏观日历** — macro.go + FOMC/CPI/NFP/PCE 2026 全年 + MacroVolAdjust | phase=normal, next=PCE 04-30(87h), vol_mult=1.00; 事件前24h vol↑15-30%, 事件后2h cooldown 15%; 6单测 |
