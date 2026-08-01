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

- **历史 ladder TP + 止损**（`-exit_mode=ladder`，Phase 7.b 已实现，但 copytrade 不直接启用价格 TP/SL）：
  - **TP1** 价格较入场涨 15% → 清 50%（默认；`-ladder_tp1_pct=0.15` `-ladder_tp1_frac=0.50`）
  - **TP2** 价格较入场涨 30% → 清剩余 100%（默认；`-ladder_tp2_pct=0.30` `-ladder_tp2_frac=1.0`）
  - **Stop-loss** 价格较入场跌 5% → 清 100%（`-ladder_sl_pct=0.05`，04-20 22:42 收紧自 0.10）
  - **MaxHold** 4h 强平（避免锁死资金）
  - 相比老方案 (+30%/+60%/余量 hold) 更激进：拉早 TP1、补足 TP2 清仓、加硬止损、加超时，不保留 hold tail
  - SL 收紧依据：Phase 7.d sweep 显示 SL 是主导杠杆——python pool 中 baseline -30.41 → SL=5% 的 -0.23（top 10 全部 SL=5%）。TP 阈值在 5%~100% 之间只差 7 个点。
- **成本建模** — `-slippage_bp` 对买卖两边施加不利滑点；动态平台费按 `shares × rate × price × (1-price)` 计算，优先读取 CLOB `/clob-markets/{condition_id}` 的逐市场 `fd.r`，`-taker_fee_rate` 仅作接口失败时的回退；`-fee_bp` 保留给额外固定 builder fee。paper 双边计费写 journal，净 PnL = 毛 PnL − entry_fee − exit_fee。智能钱体育模拟盘默认每边 50bp 滑点、回退费率 0.05。
- **可成交退出价** — 未结算市场的 timeout 只允许使用订单簿最新 Best Bid，再由 paper client 施加卖出滑点和费率；Best Bid 缺失时延后退出并记录日志，禁止使用 Gamma outcome/mid 代替成交价。
- ~~**Phase 7.c 长尾市场**~~ — 老板 04-20 21:42 拍板**不做**：周期太长不适合 90 USDC 资金体量。
- **历史回放** — 把 python trades 的 entry_price × market_id 灌进 Go backtester 验 ladder_TP 期望曲线，Phase 7.d。
- **逐仓路径** — copytrade/重启恢复仓位会动态加入 WebSocket 标的订阅并记录 1Hz tick；每次 Start 创建独立文件，禁止跨重启向旧 `pN.jsonl` 追加。回测器拒绝混有多个 asset 的旧污染文件，并按真实时间戳计算 timeout。

### 2.5 出场模式

- `-exit_mode=hold`（当前默认，手动点单用）：**买了就等最终结果**——不看 SL/TP/timeout，开仓后**只等 market resolve**，按 gamma `OutcomePrices[SlotIdx]` 清算（赢家侧 1.0、输家侧 0.0）。settlement watcher 每 60s 轮询 gamma，`closed=true` 即清算；5 min 打一行 `hold_status` 便于 grep。
- `-exit_mode=auto`（legacy）：ExitTracker 按旧版（反转 3 tick / 回撤 2pp / 入场-3pp 止损 / 30min 超时）。
- `-exit_mode=ladder`（Phase 7.b）：TP1/TP2/SL/Timeout 分级，见 §2.4 参数。paper 期支持 tranche 级别的分批平仓，journal 每个 tranche 一行。
- 智能钱体育模拟盘的普通赛事仓位至少观察 30m：未来开赛的仓位截止到开赛后 30m，已开赛仓位截止到入场后 30m；Gamma/CLOB 时间缺失不会永久缓存，30s 后重试。到期检查每 5s（`-exit_poll_interval`），常规结算查询仍每 60s；timeout 后同市场默认冷却 30m（`-timeout_reentry_cooldown`）。
- copytrade 同步运行不下单的退出影子观测：10/20/30/45/60m、可成交 bid 连续 15s 命中 -20%/-25%、以及 +30%/+50%；所有候选退出均按 Best Bid 扣卖出滑点、逐市场平台费和固定费，并扣已支付入场费后输出净 PnL。实际仓位提前退出后，shadow 仍保留到 60m 样本完成；分批止盈的每个 tranche 独立计退出费。日志事件为 `copytrade_exit_shadow`，日报 `reports/smartmoney-exit-shadow.md` 按策略和品类去重汇总。
- 足球比分候选池每小时扫描最多 300 个 Exact Score 市场、每个 token 前 250 名持仓、最低 50 shares、最多保留 500 个地址，更新后自动重载智能钱模拟盘。比分地址保留通用 tier，同时按 Exact Score 的 Yes 侧交易样本、已结算 ROI 和扣滑点后的跟单 ROI 单独评 A/B/C/D；专项 A/B 仅能通过比分模拟盘的 B 级门槛。比分 BUY 必须在 2m 内，同场可跟多个比分但总敞口默认不超过 60U，单笔 20U，最长观察 150m；普通同场多盘口总敞口默认不超过 100U。
- paper journal 写入 `policy_version` 和完整跟单钱包地址。`reports/smartmoney-paper-pnl.md` 按版本、单笔金额、策略、成本口径和来源拆分毛 PnL、双边费用与净 PnL，同时列出开放敞口和按零价值计的保守 PnL。
- 本地地址策略至少收集 10 个已结仓仓位：净 PnL ≤ -5U 自动加入 `wallets.paper-demoted.txt`，净 PnL ≥ +5U 且 ROI ≥ 2% 才进入 B 级 `wallets.paper-promoted.txt`；每日和每小时任务刷新后重载模拟盘。
- 所有模式都保留日亏损熔断 + 单笔亏损 flag + feed-silence watchdog，且扣 fee 计净 PnL。
- 运行停用按组件隔离：`db/live-trading.disabled` 在启动时和运行中阻止实盘；实盘还必须有权限 `0600`、绑定钱包且最长 24h 有效的 `db/live-trading.enabled`。守卫在每次签名前复查，单笔默认上限 20U、单进程累计 BUY 默认上限 100U；任一检查失败会关闭当前实盘进程。`db/research.disabled` 停研究迭代，`db/monitoring.disabled` 停旧监控；模拟盘和地址迭代不受影响。

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
- 签名密钥只在本地内存持有，启动时从 Bitwarden 拉；禁止通过命令行参数传私钥，API 凭据和助记词不得进入日志

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

### 7.1 七天保守自治优化

- `com.polymarket-go.daily-p0-optimize` 每天本地时间 04:30 在日报稳定后运行一轮 7 天实验；Day 1 建基线，Day 2/4/6 允许最多一个小改动，Day 3/5 只观察，Day 7 只做周总结，然后停止自动改动。
- 全周最多接受 3 次改动且至少间隔 48h。单次只改一个可归因问题，避免多项同时变化后无法判断 PnL 或稳定性变化来源。
- 内层 Codex 在隔离的临时 Git worktree 中处理有运行证据的正确性、风险控制、数据完整性、可靠性或可观测性问题，也可在严格数据门槛下调整模拟盘/影子参数；不直接接触正在更新的主工作区，也不负责 Git 状态或网络部署。
- 外层控制器拒绝源码脏启动、远端分叉、控制文件改动、越界路径、超过 12 个文件、超过 600 行、二进制内容及疑似敏感值；通过 shell/plist 校验、`git diff --check` 和 `go test -race ./...` 后才 commit/push。
- 失败改动进入独立 Git stash，避免 23:55 的每日 Git 保存误推未验证代码。
- 实盘开关、运行数据库、生成报表、地址池、金额/晋升/总敞口继续由人工控制。模拟盘/影子参数只有在至少 30 个成熟独立样本、两个非重叠切片按双边手续费和滑点计算后都改善时，才允许单参数相对调整不超过 10%。

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
