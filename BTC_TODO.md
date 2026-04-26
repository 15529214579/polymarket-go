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

### 3. [P1] Fear & Greed Index 接入
- alternative.me API（免费）
- 极度恐慌 (<25) 时 dip 市场被高估概率更大
- 极度贪婪 (>75) 时 reach 市场被低估概率更大
- 作为 BS gap 的加权因子

### 4. [P1] Funding Rate 接入
- Binance perpetual BTCUSDT funding rate（8h 周期）
- 正 funding = 多头拥挤 → 回调风险
- 负 funding = 空头拥挤 → 反弹概率
- 用于修正马尔科夫转移概率

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

### 9. [P1] Hidden Markov Model (HMM)
- 隐状态: 趋势/震荡/反转 三个 regime
- 观测: 价格变动 + 成交量 + 波动率
- Baum-Welch 训练，Viterbi 解码当前 regime
- 只在趋势 regime 做方向性 bet

### 10. [P1] 动态波动率 (GARCH/EWMA)
- 当前用固定窗口历史波动率（90天）
- EWMA λ=0.94 对近期波动赋更高权重
- GARCH(1,1) 捕捉波动率聚集效应
- 替换 BS first-passage 里的 σ 输入

### 11. [P1] 波动率曲面建模
- 不同 strike 用不同隐含波动率
- 低 strike（dip $25K）→ 左尾肥尾 → σ 上调
- 高 strike（reach $200K）→ 右尾 → σ 上调
- 中间 strike（$80K-$100K）→ 基准 σ

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

### 15. [P1] 出场策略
- BS gap 从 >10pp 收窄到 <3pp → 平仓（alpha 已耗尽）
- BTC 价格向不利方向移动 >5% → 止损
- 持仓超过 7 天无 gap 变化 → timeout 退出

### 16. [P1] 仓位管理 (Kelly Criterion)
- f* = (bp - q) / b，其中 p=BS概率, q=1-p, b=PM赔率
- gap 越大 → Kelly 分配越大（但 cap 在 half-Kelly）
- 单市场最大敞口 $20（paper 模式）

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

### 19. [P1] PnL 归因分析
- 按模型拆分: BS gap / Markov / 联合
- 按市场拆分: reach vs dip / 不同 strike
- 按时段拆分: 亚洲时段 / 美洲时段 / 欧洲时段
- 每日回测自动输出归因报告

### 20. [P1] 历史 PM 价格补全
- 当前只有 1 个快照（Apr 26）
- 每小时采集一次 PM BTC 市场价格
- 积累 7 天后做 PM 价格变动 vs BTC 价格变动的回归分析

### 21. [P2] bonereaper 策略逆向
- 拿到钱包地址后拉全部交易
- 分析: 入场时机（BTC 处于什么状态时买入）
- 分析: 偏好 strike（总是买 reach 还是 dip）
- 分析: 持仓周期（短线翻转还是长期持有）
- 复制其信号加入 whale tracker

---

## 五、基建与回测

### 22. [P1] 回测引擎增强
- 加入 PM 滑点模拟（基于 orderbook 深度）
- 加入 PM 手续费 2% 扣除
- 加入持仓期间 BTC 波动对 PM 价格的动态影响
- 输出 Sharpe ratio / max drawdown / Calmar ratio

### 23. [P2] 多币种扩展 (ETH/SOL)
- PM 上有 ETH/SOL 价格预测市场
- 复用同一套 Markov + BS 框架
- 币种间相关性利用: BTC 涨 → ETH 通常跟涨 → 联合入场

### 24. [P2] Regime Detection 自动切换
- 牛市 regime: 只做 reach 市场（买 Yes）
- 熊市 regime: 只做 dip 市场（买 No）
- 震荡 regime: 双向扫 gap
- 用 HMM 隐状态或 200 日均线斜率判断

### 25. [P3] 信号质量评分系统
- 综合: BS gap 大小 + Markov 置信度 + 时间尺度一致性 + funding rate + F&G
- 0-100 分制
- >80 分自动入场 / 60-80 推按钮 / <60 静默记录
- 每日回测更新评分模型权重

---

## 进度追踪

| 日期 | 完成项 | 发现 |
|------|-------|------|
| Apr 26 | 基础 Markov + BS + PM tracker + 1h 回测 | Markov 49.9% 掷硬币; BS gap 10-23pp 有价值 |
| Apr 26 | BTC live strategy 上线 (1h scan, gap>7pp) | 首次扫描: $50K dip PM=42.5% BS=19% gap=-23pp |
| Apr 26 | 5m/15m K线 + multi-TF Markov + PM delta tracking | MIXED(bull=2.7%,bear=2.9%) 放行全部 BS gap 信号 |
