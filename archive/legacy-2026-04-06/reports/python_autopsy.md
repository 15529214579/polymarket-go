# Python 项目战绩解剖 · 2026-04-20

**目的：** 用 `~/.openclaw/workspace-dev3/polymarket-agent/db/polymarket_agent.db`（python 项目真实盘历史）回答一个问题——
python 真的赚到过钱吗？赚在哪？亏在哪？我们该抄什么、躲什么？

**数据范围：** `trades` 表 361 行（182 已结 / 113 TIME_STOP 复写伪样本 / 61 OPEN 未平 / 5 CANCELLED）；`opportunity_log` 表 4085 行扫描流；时间窗 2026-04-16 → 2026-04-20。

---

## 1. 全盘累计 — 严重亏损

| 指标 | 值 |
|---|---|
| closed 交易 | 182 笔 |
| 累计 PnL | **−97.50 USDC** |
| 峰值 equity | +1.81 USDC（04-12） |
| 最大回撤 | −101.76 USDC（04-19） |

## 2. 按 scan_type 看（只看有对应 opportunity_log 的子集）

| scan_type | closed | wins | PnL | avg |
|---|---|---|---|---|
| llm_scan | 4 | 2 | +0.50 | +0.125 |
| **theodds_h2h** | **13** | **0** | **−36.51** | **−2.81** |

## 3. 按出场类别（去重 TIME_STOP 伪样本）

| cat | n | wins | PnL | avg | 备注 |
|---|---|---|---|---|---|
| **ladder_TP（阶梯止盈 +30/+60）** | **6** | **6** | **+3.77** | **+0.63** | ⭐ 唯一通杀 |
| force_close（临近截止 <2h 强平） | 3 | 2 | +0.50 | +0.17 | ✅ 小样本正 EV |
| write-off:market_resolved | 23 | 0 | 0.00 | 0 | 到期自然 |
| write-off:low_price (<0.01) | 20 | 0 | −2.66 | -0.13 | 底价归零 |
| TIME_STOP (18-24h 复写) | 114 | 0 | −62.13 | -0.55 | **实为 2 笔 × 57 行日志**，是同一单每次心跳写一行 |
| [THEODDS] × 市场名 | 9 | 0 | −31.27 | -3.47 | ⚠️ 全员阵亡 |
| manual | 3 | 0 | −4.47 | -1.49 | 人工介入也翻车 |

## 4. 按入场价分档（去 TIME_STOP 伪样本）

| 价带 | n | wins | PnL |
|---|---|---|---|
| 0.00–0.10 | 13 | 0 | −0.15 |
| 0.10–0.25 | 10 | 1 | −4.92 |
| 0.25–0.50 | 14 | 1 | −8.63 |
| 0.50–0.75 | 8 | 1 | −14.54 |
| 0.75–0.90 | 2 | 0 | −10.44 |

全价带都是负 EV —— 但 ladder_TP 的 6 个赢家全部落在 **0.13 ~ 0.65** 这一段。

## 5. 赢家清单（8 笔 winning trades）

| entry | close | PnL | reason |
|---|---|---|---|
| 0.13 | 0.245 | +1.834 | 阶梯止盈 2 档（Iran x Israel conflict） |
| 0.36 | 0.745 | +0.891 | 临近截止强平 |
| 0.13 | 0.185 | +0.611 | 阶梯止盈 1 档（Iran x Israel conflict） |
| 0.36 | 0.600 | +0.556 | 阶梯止盈 2 档 |
| 0.6455 | 0.948 | +0.469 | 临近截止强平 |
| 0.36 | 0.500 | +0.324 | 阶梯止盈 1 档 |
| 0.2995 | 0.4485 | +0.249 | 阶梯止盈 1 档 |
| 0.6455 | 0.895 | +0.193 | 阶梯止盈 1 档 |

**特征：**
- 入场 **0.13 ~ 0.65**（中间价带、偏弱侧）
- 出场 **+30% / +60%** 阶梯 TP，**或** 临近截止 <2h 强平
- 2 个最大的赢家来自一个**政治预测**市场（Iran x Israel conflict by April 30），不是体育
- 体育侧赢家都来自 scan_type 空白（非 theodds，也非 llm）—— 可能是 `internal_arb` 或早期策略

## 6. 输家清单（头部）

| entry | close | PnL | reason |
|---|---|---|---|
| 0.77 | 0.0 | −10.43 | [THEODDS] Eintracht × Leipzig |
| 0.80 | 0.0 | −8.99 | [THEODDS] Hellas Verona × Lecce |
| 0.68 | 0.0 | −3.98 | [THEODDS] Bayern × VfB（×2 单） |
| 0.65 | 0.0 | −2.85 | [THEODDS] Leeds × Wolves |
| 0.60 | 0.0 | −2.17 | [THEODDS] Genoa × Como |
| TIME_STOP 重复写 | 0 / 0.003 | −62 合计 | 同一坏单日志刷 57 行 |

**特征：**
- theodds_h2h 策略在足球 favorite 上 **100% 阵亡**（赛果本就不可预测，gap 来自 league mismatch / bookmaker map 错 / 低流动性，不是真 arb）
- 高价入场（>0.60）+ 结算归零 = 本金完全损失
- TIME_STOP 两单在 18-24h 前就被挂出去，长时间无 TP/SL 守护

## 7. bookmaker / gap 分布（opportunity_log 总览）

- 总 scan 1325 条，但 **1036 条 gap ≥ 15pp**，avg gap 33.72pp → 大概率是 **league 错配或 outright vs season 计价错位**，不是真 arb
- 5–10pp 窄 gap 只有 16 条 scan，**全数没被执行也没结算** → 真正可做 arb 的样本极少

## 8. 战略结论（2026-04-20 21:3x）

### 扬（python 证明过能赚的三件事）
1. **ladder_TP 出场**：+30% / +60% 两档止盈 —— 6/6 赢家的唯一共性
2. **中间价带入场**（0.13 ~ 0.65）—— 所有赢家落在这一带
3. **长尾高 payoff 市场**（比如政治事件解 0.13 → 0.94 结算）存在 —— 但筛选靠 LLM 或人工，不靠 gap

### 避（python 证明过会亏的三件事）
1. **theodds_h2h >5pp gap 策略**：13/13 全败，基本都是伪 arb
2. **高价 favorite 追入**（>0.65 且非近端）：一次翻车归零就抹平多笔盈利
3. **TIME_STOP 长守护**（18h+ 挂着不 TP/SL）：靠运气吃回撤，不是策略

### 不采用
- 继续做 PM 内部动量追涨（Day-1 自己 0W4L 证实，python DB 也没显著赚钱迹象）
- 盲目上 cross-venue gap：真正窄 gap 样本太稀薄，python 在这条路上的历史就是反例

---

## 9. 落地到 polymarket-go

**R3（prompt-only）已生效**：`scripts/bot-daemon.sh` 默认 `-signal_mode=prompt`，daemon 60395（04-20 21:32 起）不再 auto-open。

**R4（Apr 29 不切实盘）确认**：Phase 3 V2 签名暂缓。Paper 继续顺延到 ladder_TP 思路验证为止。

**Phase 7 计划（whitelist 过滤 + ladder TP）：**
- 7.a — prompt 过滤：`-min_entry=0.15 -max_entry=0.70` + 跳过足球 moneyline favorite，信号只在可作战价带出现
- 7.b — ladder TP exit：`-exit_mode=ladder` 新模式，+30% 出 1/3、+60% 出 1/2 剩余、余下 hold 到结算
- 7.c — 政治/长尾市场探索：沿 gamma tag 扩出 `politics/news` 类，人工筛选、ladder TP 默认开
- 7.d — backtest 回放：把 python trades 的 market_id × entry_price 灌 Go backtester，跑 ladder_TP 策略得期望曲线，决定是否 Apr 29 实盘

回测先行，代码先过 7.a（最简单、风险最低），再看 7.b。
