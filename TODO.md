# polymarket-go TODO

**维护人：** 5号 (monitor) — 我自己接活自己处理，不接 1号 派单。

## 🛠 进行中

### P00 — 自我迭代工具（最高优先级）✅
> 没这个我醒不过来，先做。cron + heartbeat 是后续所有 phase 能自我推进的前提。

- [x] `scripts/heartbeat.sh`：单次自检（git/build/uncommitted/开放 TODO/最近日志）
- [x] `scripts/cron-poke.sh`：周期入口，写日志 + 更新 state.json + 计算 alert
- [x] macOS crontab 注册：~~每 30 min~~ → **每 20 min**（`*/20 * * * *`，04-19 23:50 改）
- [x] 夜间静默（00:00-07:59 SGT）：quiet_window 标记，不触发 alert 升级
- [x] `state.json`：last_heartbeat / last_commit / uncommitted / ticks_no_progress / alert
- [x] OpenClaw 唤醒 cron：`17,47 8-22 * * *`（session-only，需每次开 session 时重装）
- [x] alert 升级通道：`scripts/alert-dispatch.sh` — 读 state.json 若 alert 非空 → Telegram Bot API 直推老板，2h cooldown、夜间静默、.env.local 存 token（gitignored）。04-19 23:49 端到端验证 ok=true。

### Phase 0 — Bootstrap ✅
- [x] `go mod init github.com/15529214579/polymarket-go`
- [x] 目录骨架：`cmd/bot/`, `internal/{feed,strategy,order,risk,log,config}/`
- [x] Makefile + .gitignore + build 通过
- [x] git init + 首个 commit（3d072a7）
- [x] github public repo（`github.com/15529214579/polymarket-go`）
- [ ] golangci-lint 配置（非阻塞）

### Phase 1 — 数据层（进行中）
- [x] gamma REST 客户端（LoL 市场筛选）— 04-20 00:02 跑通，`./bin/bot -mode=discover` 拉到 59 个活跃 LoL 市场（LPL/LCK/LEC/LCS）
- [x] Polymarket WSS 客户端（自动重连、心跳、book/price_change/last_trade_price 解码）— 04-20 00:09 跑通，`./bin/bot -mode=feed -markets=8` 20s 采到 44 book + 2 trade 事件，VIT/GIANTX 活局 mid 0.83/0.84 稳定
- [x] orderbook 内存模型（bid/ask 深度、最近成交流）— 在 WSS 客户端内，price_change 增量合并到本地 bookState
- [x] tick 采样器（1s 粒度，滑窗 60s）— 04-20 00:30 跑通，`./bin/bot -mode=sample -markets=8 -window=20` 25s 内采 352 tick + 32 window summary，VIT vs GIANTX LEC Game 2 实时 mid 0.175 跟盘。单测 3 个全过。**Phase 1 完成。**

## 💤 待启动

### Phase 2 — 策略层（2 天）
- [x] 动量信号检测（N秒涨幅、tick 单调性、主动成交占比）— 04-20 00:42 `internal/strategy/detector.go` 上线，`./bin/bot -mode=detect -markets=10 -window=30` 75s 实测首发信号：LOUD vs Leviatan LCS Game 1 Winner, Δ+4.00pp, tail 4/5 ups, buy_ratio 1.00. 5min per-asset cooldown 生效。
- [x] 出场信号（反转、止损、超时）— 04-20 08:45 `internal/strategy/exit.go` 上线，ExitTracker：反转(3下行tick/2pp回撤) + 3pp硬止损 + 30min超时。5 个单测全过；主 detect 模式 paper-open→exit 链路已接；75s 实盘烟测此时段 LoL 盘全平静未触发（预期）。
- [x] 仓位管理（单仓去重、总敞口）— 04-20 10:41 `internal/strategy/position.go` 上线，`PositionManager`：per-asset/per-market 双重去重（拒绝同市场 YES+NO 同开）、MaxOpenPositions=6、MaxTotalOpen=45 USDC（50% of 90 本金）。实时 PnL 用 units×(exit-entry)。6 个单测全过（open/close/dedupe-by-market/dedupe-by-asset/max caps/reopen）。detect 模式已接 `pm.Open → exit.Open → pm.Close`；45s 实盘烟测此时段无信号（market quiet，预期）。**Phase 2 完成。**

### Phase 3 — 下单（1-2 天，方案 A：自签+broadcast，**直接按 V2 出生**）
> 04-20 10:36 老板拍板：Paper 顺延到 Apr 29，cutover 后 wrap + 实盘。
- [x] 2026-04-20 11:4x — Phase 3.a 骨架：`internal/order/` Intent/Result/Client 接口 + PaperClient（slippage 模型）+ 7 单测 ✅
- [x] 2026-04-20 11:4x — Phase 3.b：PaperClient 接进 detect 循环。signal → `Buy Intent` Submit → 用 fill 价开仓；exit → `Sell Intent` Submit → 用 fill 价平仓 → 实现 slippage-priced PnL。新增 `-slippage_bp` flag，order_id 写入日志。35s 实盘烟测：paper_client.ready 正常，无信号（market quiet）。build+test 绿。
- [ ] Phase 3.0 前置（Apr 28 19:00 SGT cutover 后执行）：Bitwarden 取私钥 → Collateral Onramp `wrap(90.41 USDC.e)` → 拿等量 pUSD
- [ ] Bitwarden 取助记词 → 派生私钥（启动时只驻内存）
- [ ] EIP-712 typed data 签名（**V2** order struct：去 `taker/expiration/nonce/feeRateBps`，加 `timestamp/metadata/builder`；domain version `"2"`，V2 Exchange 地址）
- [ ] CLOB **V2** REST `/order` POST 客户端
- [ ] 成交回执轮询 + status 机
- [ ] **Apr 28 cutover 当天 WSS 烟测**：18:45 SGT 待机 → 20:15 SGT 跑 `-mode=detect` 20min，验 3 种消息类型帧结构

### Phase 3.5 — 半自动点选下单（Hybrid UX）
> 信号推 DM → 按钮点选 1U / 5U / 10U → 回调触发同一签名路径下单。依赖 Phase 2 信号 + Phase 3 下单。
- [x] 2026-04-20 12:0x — Phase 3.5.a（outbound）：`notify.SignalPrompt` 把信号打成 DM，附 inline keyboard "Buy 1U / 5U / 10U"；callback_data = `buy:<nonce>:<sizeUSD>`。`internal/notify/pending.go` 上线（`PendingStore` TTL=60s、one-shot Claim、hex8 nonce）。3 pending 单测 + 1 telegram signal 单测全过。**outbound 独立于 callback 方案 A/B，已锁。**
- [ ] **⚠️ 等老板拍板：callback 路径 A/B**（见下）
- [ ] Phase 3.5.b（inbound）：callback 消费 → Claim nonce → `paper.Submit`（paper）或 Phase 3 CLOB 签名（Day 9 起实盘）
- [ ] 超时作废（>60s 自动在 Telegram 上编辑原 DM 为 "已过期"，callback 返回 alert=expired）
- [ ] 安全：sender_id 过滤只认老板；PendingStore 已做 one-shot，callback_data 带 nonce 防重放
- [ ] Paper 期间：点了按钮走 paper 路径；Day 9 起自动走真下单
- [ ] 成交回执：下单成功/失败都回 DM 单条小回执（`notify.LargeFill` 已有路径，小单复用一个 `notify.OrderResult` 事件即可）

**callback 路径选项（开工前拍板）：**
- **A. OpenClaw 转发**：OpenClaw telegram provider 接到 callback_query 后把"老板点了 buy:xxx:5"路由进当前 5号 session，5号 session 处理 → 调 bot 的 CLI 或本地 RPC 完成下单。
  - 优：零新服务、凭据全落 OpenClaw、安全模型复用
  - 劣：耦合 OpenClaw 的 callback 转发能力，目前尚未验证是否支持
- **B. bot 自起 long-poll**：`cmd/bot` 里开一个 goroutine 调 Telegram `getUpdates` 消费 callback_query（同一个 bot token），命中 nonce → 内部直接下单。
  - 优：完全 in-process、不依赖外部
  - 劣：同一 token 如果 OpenClaw 也在 poll，会抢 update；需要确认 `.env.local` 的 token 是否和 OpenClaw 用的同一个，否则就要开个 sidecar bot

### Phase 4 — 风控 + 可观测（1 天）
- [x] 2026-04-20 11:4x — Phase 4.a：`internal/risk/risk.go` 上线。日亏损熔断（15% × 90.41 = -13.56 USDC 上限）+ per-trade 单笔损失 ≥ 3 USDC 计数旗标 + WSS feed-silence watchdog（30s 无 book/trade → trip）+ 手动 Pause/Resume。SGT 日切滚动但不自动解除 breaker（SPEC §6 "等老板手动恢复"）。8 个单测全过。接进 detect 循环：开仓前 `AllowOpen` 门控，close 后 `OnClose` 累计，5s 心跳 `CheckFeed`，60s `risk_status` 日志。35s 实盘烟测：`risk.ready` + 无 trip（预期，LoL 市场平静）。
- [x] 2026-04-20 11:55 — Phase 4.b：`internal/notify/` 上线（Notifier interface + TelegramNotifier + Nop）。事件：`RiskTrip`（daily_loss / feed_silence / manual_pause，每次 trip 一次）+ `RiskResume`（breaker 恢复）+ `LargeFill`（|realized_pnl| ≥ `-large_fill_usd`，默认 3 USDC）。异步 goroutine queue（32 buf，满则 drop 不阻塞交易）。detect 循环两个 trip site + close site 都接好；resume 在 watchdog 跟踪 prevBlocked→!blocked 转换触发。`internal/config/dotenv.go` 辅助（启动加载 .env.local，不覆盖已有 env）。5 notify 单测 + 2 dotenv 单测全过，实盘 smoke ok (`notify.ready mode=telegram`)。
- [ ] 日结报表 cron（SGT 00:00 汇总 realized PnL / trade 分布 / block 时长）

### Phase 5 — Paper 跑 8 天（Apr 20 → Apr 28 cutover）
> 04-20 10:36 老板拍板 A：顺延 Paper 跨过 V2 cutover，Apr 29 起实盘（原 7 天 → 8 天）
- [ ] Day 1-3（Apr 20-22）：数据/信号合理性检查
- [ ] Day 4-7（Apr 23-26）：策略参数微调
- [ ] Day 8（Apr 27-28）：出 paper 报表；28 日晚 cutover 烟测 + wrap USDC→pUSD
- [ ] **Apr 29**：老板 review paper + V2 验证 → 实盘启用

## ✅ 已完成

- [x] 2026-04-19 23:33 — 策略方向对齐（动量跟进 / LoL 练手 / 7 天 paper）
- [x] 2026-04-19 23:31 — 独立钱包入 Bitwarden（`Polymarket-Go Wallet`）
- [x] 2026-04-19 23:35 — SPEC.md / TODO.md 初稿
- [x] 2026-04-19 23:34 — 下单通道敲定 A（自签+broadcast）
- [x] 2026-04-19 23:58 — PRINCIPLES.md 上线（7 条老板拍板原则持久化到 repo）
- [x] 2026-04-20 00:02 — Phase 1.1 完成：gamma LoL 市场发现 +  WSS 骨架（commit d5c67b9）
- [x] 2026-04-20 00:09 — Phase 1.2/1.3 完成：真 WSS dial + book/price_change/last_trade_price 解码 + 本地 orderbook 重建；活 LEC 盘 VIT/GIANTX 实时 bid/ask 跑通
- [x] 2026-04-20 00:13 — 老板预存启动资金到 Go 钱包：90.405327 USDC.e + 111.030024 POL，Day 0 PnL 基准已锁定（SPEC §5.1）
- [x] 2026-04-20 00:42 — Phase 2.1 完成：momentum detector 上线并在实盘 LCS Game 1 Winner 触发首个信号（Δ+4pp、tail 4/5、buy_ratio 1.0）
- [x] 2026-04-20 08:45 — Phase 2.2 完成：ExitTracker（反转/止损/超时 4 规则，5 单测）
- [x] 2026-04-20 10:36 — SPEC/TODO 对齐 Polymarket V2 迁移：Paper 顺延 Apr 29，Phase 3 直接按 V2 签名，cutover 当天 wrap USDC→pUSD + WSS 烟测
- [x] 2026-04-20 10:41 — Phase 2.3 完成：PositionManager（双重去重 + 敞口/仓位数上限，6 单测）；detect 链路 pm.Open → exit.Open → pm.Close 闭环。**Phase 2 整层通关。**
- [x] 2026-04-20 11:4x — Phase 3.a/3.b 完成：`order.Client`/`PaperClient` + slippage + detect 循环改走 `Submit→Fill→Open/Close`（Apr 28 cutover 后换真 V2 client 零改动）
- [x] 2026-04-20 11:4x — Phase 4.a 完成：RiskManager（日亏损 -13.56 USDC 熔断 + 30s feed-silence watchdog + 手动 Pause/Resume + 8 单测），detect 闭环已装风控门控

## ❌ 不做

- 不接 1号 派的 python 活
- 不碰 python 项目任何文件
- 不依赖外部数据源
