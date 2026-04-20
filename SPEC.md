# polymarket-go — SPEC

**Owner:** 5号 (monitor)  
**Repo:** (待建) github public  
**Wallet:** `0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e` (Bitwarden → `Polymarket-Go Wallet`)  
**Status:** draft — 2026-04-19

---

## 1. 项目目标

用 Go 重写 Polymarket 交易代理，**不污染** python 项目。独立钱包、独立 repo、独立 TODO。

**为什么 Go：** 低延迟、goroutine 并发、持久 WSS 长连接稳定。

## 2. 策略（2026-04-20 21:30 R3+R4 改版）

> **戳破 04-19 版本的追涨假设 + Phase 6 backtest 把 "PM vs bookmaker >5pp" 也证伪之后，第二次重写。详细推翻依据见 [`reports/python_autopsy.md`](reports/python_autopsy.md)。**

### 2.0 决策框架（R3 + R4）

- **R3** — auto-open 关闭，daemon 默认 `-signal_mode=prompt`。所有信号走 Telegram DM prompt，老板**手动点按钮**才下单。
- **R4** — Apr 29 不切实盘。Phase 3 V2 签名暂缓，等 Phase 7 的 ladder_TP 策略回测通过再讨论实盘。
- **策略定位** — "信号推荐 + 人工过滤"，不是"全自动 momentum"。模型和触发条件向 python DB 里真正赚过钱的模式靠拢。

### 2.1 扬（python DB 证明过能赚的三件事）

1. **ladder TP 出场**：+30% 出 1/3、+60% 出 1/2 剩余、剩余 hold 到结算（6/6 赢家的唯一共性）
2. **中间价带入场**（0.15–0.70）：所有赢家都落在这段，偏弱侧
3. **长尾高 payoff 市场**（政治事件、明确 asymmetric 赔率）：可作战赛道，但筛选靠人工/LLM，不是纯 gap

### 2.2 避（python DB 证明过会亏的三件事）

1. **theodds_h2h >5pp gap**：13/13 全败，gap 来自 league mismatch 而非真 arb
2. **高价 favorite 追入**（>0.70）：一次翻车归零抹平多笔盈利
3. **裸追涨 Δ>3pp in-play momentum**：Day-1 自己 0W4L，python DB 里没对应 scan_type 佐证能赚

### 2.3 当前触发条件（Phase 7.a — 价带过滤版）

- 旧：60s Δ≥3pp + tail 4/5 + buy ≥60%（momentum 原生信号，保留作为候选池）
- 新增过滤器：
  - `entry_price` 必须在 `[min_entry_price, max_entry_price]` 闭区间（默认 0.15–0.70）
  - 命中过滤的信号才发 `SignalPrompt` DM；没命中的只留 `signal` 日志不发 DM
- 老板点按钮后走原 manual_open 路径（paper 期间），hold 到结算（`-exit_mode=hold`）

### 2.4 待验证（Phase 7.b+）

- **激进 ladder TP + 止损 + 手续费**（`-exit_mode=ladder`，Phase 7.b 已实现）：
  - **TP1** 价格较入场涨 15% → 清 50%（默认；`-ladder_tp1_pct=0.15` `-ladder_tp1_frac=0.50`）
  - **TP2** 价格较入场涨 30% → 清剩余 100%（默认；`-ladder_tp2_pct=0.30` `-ladder_tp2_frac=1.0`）
  - **Stop-loss** 价格较入场跌 10% → 清 100%（`-ladder_sl_pct=0.10`）
  - **MaxHold** 4h 强平（避免锁死资金）
  - 相比老方案 (+30%/+60%/余量 hold) 更激进：拉早 TP1、补足 TP2 清仓、加硬止损、加超时，不保留 hold tail
- **手续费建模** — `-fee_bp`（per-side 基点，默认 0 匹配 CLOB V1 实测；V2 官方数字发布后更新）。paper 双边计费写 journal，净 PnL = 毛 PnL − entry_fee − exit_fee。
- ~~**Phase 7.c 长尾市场**~~ — 老板 04-20 21:42 拍板**不做**：周期太长不适合 90 USDC 资金体量。
- **历史回放** — 把 python trades 的 entry_price × market_id 灌进 Go backtester 验 ladder_TP 期望曲线，Phase 7.d。

### 2.5 出场模式

- `-exit_mode=hold`（当前默认，手动点单用）：**买了就等最终结果**——不看 SL/TP/timeout，开仓后**只等 market resolve**，按 gamma `OutcomePrices[SlotIdx]` 清算（赢家侧 1.0、输家侧 0.0）。settlement watcher 每 60s 轮询 gamma，`closed=true` 即清算；5 min 打一行 `hold_status` 便于 grep。
- `-exit_mode=auto`（legacy）：ExitTracker 按旧版（反转 3 tick / 回撤 2pp / 入场-3pp 止损 / 30min 超时）。
- `-exit_mode=ladder`（Phase 7.b）：TP1/TP2/SL/Timeout 分级，见 §2.4 参数。paper 期支持 tranche 级别的分批平仓，journal 每个 tranche 一行。
- 所有模式都保留日亏损熔断 + 单笔亏损 flag + feed-silence watchdog，且扣 fee 计净 PnL。

### 2.6 仓位（prompt 模式）

- Paper 阶段：按钮档 1U / 5U / 10U 由老板手选，`PositionManager.OpenSized` 已支持可变 size。
- 实盘阶段（Phase 7 过后才考虑）：总资金 × 5%/笔，由回测结果定。

## 3. 数据源

**仅 Polymarket 官方 API，不依赖外部：**
- `wss://ws-subscriptions-clob.polymarket.com/ws/` — orderbook 实时订阅
- `https://clob.polymarket.com/` — 市场元数据、下单
- `https://gamma-api.polymarket.com/` — 市场列表、LoL 赛事筛选

**LoL 赛事筛选：**
- gamma events 按 `tag=League of Legends` 或 title 正则匹配
- 只订阅 `live=true` 的 markets

## 4. 下单通道（**老板 04-19 23:34 拍板：A**）

**A：独立钱包自己 sign+broadcast** ✅
- 新钱包已独立，助记词/私钥已入 Bitwarden（`Polymarket-Go Wallet`）
- Go 侧用 `go-ethereum` 本地 EIP-712 签名 → Polymarket CLOB REST API 下单
- 零 python 耦合、零订单污染
- 签名密钥只在本地内存持有，启动时从 Bitwarden 拉

## 5. 生命周期（2026-04-20 10:36 调整：Polymarket V2 cutover 对齐）

- **Day 0-8：** Paper trade，从 Apr 20 到 **Apr 28 cutover 结束**（原 7 天 → 8 天，跨过 V2 切换窗口）
- **Apr 28 19:00 SGT：** Polymarket CLOB V2 cutover（~1h downtime，open order 清空，collateral 换 pUSD）
  - Cutover 后立即执行：USDC.e `wrap()` → pUSD（Phase 3.0）
  - 执行 WSS 帧烟测（3 种消息类型验证）
- **Apr 29：** 老板 review paper + V2 验证结果 → 实盘启用
- **实盘上限：** 启动资金 `90.41 USDC.e`（老板 2026-04-20 00:13 预存）

### 5.2 Polymarket V2 迁移要点（2026-04-20 10:36 入档）

| 接口 | 项目用途 | V2 影响 |
|---|---|---|
| `gamma-api.polymarket.com/markets` | Phase 1.1 市场发现 | 🟢 基本不变 |
| `wss://ws-subscriptions-clob.polymarket.com/ws/market` | Phase 1.2/1.3/1.4 | 🟢 URL 与 book/price_change 结构基本不变，cutover 当天仍需烟测 |
| CLOB REST `/order` POST + EIP-712 签名 | Phase 3（未写） | 🔴 schema 完全改，直接按 V2 出生 |
| Collateral | USDC.e → **pUSD** | 🔴 cutover 后必须 `wrap()`，否则无法下单 |

**Phase 3 签名代码直接按 V2 写：**
- EIP-712 domain version `"2"`，使用新 Exchange 合约地址
- Order struct 去掉 `taker/expiration/nonce/feeRateBps`，加 `timestamp/metadata/builder`
- 不实现 V1 兼容分支

### 5.1 启动资金（2026-04-20 00:13 SGT 快照）

| 资产 | 余额 | 用途 |
|---|---|---|
| USDC.e (`0x2791…a174`) | **90.405327** | 交易本金 |
| POL (native) | **111.030024** | gas 储备 |

- 来源链：Polygon mainnet
- 钱包：`0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e`
- 以此快照为 Day 0 基准，PnL 以此为起点计算

## 6. 风控硬限

- 单笔最大亏损 ≤ 3 USDC（paper）/ 打款额的 5%（实盘）
- 日亏损达 15% → 自动暂停，等老板手动恢复
- WSS 断线 > 30s → 关闭所有开仓（市价或挂接近市价）
- 钱包余额 < 预留 gas → 暂停下单，只平仓

## 7. 可观测性

- stdout JSON log → `~/work/polymarket-go/logs/bot.log`
- 关键事件（进场、出场、错误、风控触发）→ **telegram 私聊推送**
- 日结报表 → 每天 00:00 SGT push 一次

## 8. 不做什么（边界）

- ❌ 不碰 python 项目任何文件
- ❌ 不共享钱包、不共享数据库
- ❌ 不接 1号 派的 python 相关活
- ❌ 不依赖任何外部数据源（no bookmaker, no sports API）

## 9. 未决项（等老板拍板）

1. ~~下单通道 A/B~~ → **A 已定（04-19 23:34）**
2. "上升利好"具体参数（N秒、M tick、阈值）— 先按默认跑 paper，1-2 天后调
3. ~~Paper → 实盘切换日~~ → **Apr 29（V2 cutover 后，04-20 10:36 定）**
4. 是否需要 Discord/其他告警冗余（目前只推 telegram 私聊）
