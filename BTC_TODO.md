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

### 5. [P2] ETF 资金流数据
- SoSoValue / CoinGlass 公开 API
- BTC 现货 ETF 日净流入/流出
- 大资金流入 → 看多信号（修正 BS gap 权重）

### 6. [P2] 链上大额转账监控
- Whale Alert API 或 Blockchair
- 交易所净流入 > 阈值 → 抛压信号
- 交易所净流出 → 囤币信号

### 7. [P3] 宏观事件日历
- FOMC 利率决议 / CPI / 非农就业
- 事件前 24h 波动率上调 20-30%
- 事件后 2h 波动率回落 → 入场窗口

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

### 12. [P2] 二阶/三阶马尔科夫
- 当前: P(S_t+1 | S_t)
- 改为: P(S_t+1 | S_t, S_t-1) 或 P(S_t+1 | S_t, S_t-1, S_t-2)
- 捕捉路径依赖（连续上涨后反转概率递增）

### 13. [P2] 自动再训练
- 每日 cron 用最新 90 天数据重新训练转移矩阵
- 检测矩阵漂移（KL divergence > 阈值 → 告警）
- 滚动窗口 vs 全量数据 A/B 对比

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

### 17. [P2] PM 盘口深度分析
- gamma API orderbook 数据
- 深度不足的市场 → 滑点大 → 降低仓位或跳过
- 大单挂在某个价位 → 该价位有支撑/阻力

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

### 23. [P2] 多币种扩展 (ETH/SOL)
- PM 上有 ETH/SOL 价格预测市场
- 复用同一套 Markov + BS 框架
- 币种间相关性利用: BTC 涨 → ETH 通常跟涨 → 联合入场

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
