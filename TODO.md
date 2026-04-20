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
- [ ] 仓位管理（单仓去重、总敞口）

### Phase 3 — 下单（1-2 天，方案 A：自签+broadcast）
- [ ] Bitwarden 取助记词 → 派生私钥（启动时只驻内存）
- [ ] EIP-712 typed data 签名（Polymarket CLOB order struct）
- [ ] CLOB REST `/order` POST 客户端
- [ ] 成交回执轮询 + status 机
- [ ] Paper mode（同一路径但不发真单，记模拟 fill）

### Phase 3.5 — 半自动点选下单（Hybrid UX，不急）
> 信号推 DM → 按钮点选 1U / 5U / 10U → 回调触发同一签名路径下单。依赖 Phase 2 信号 + Phase 3 下单。
- [ ] 信号推送格式：市场 / 方向 / 当前 mid / 触发指标 + 三颗 inline 按钮（1U / 5U / 10U）
- [ ] callback 接收通道：OpenClaw 把 Telegram callback_query 转发到 5号 session（或走独立 bot webhook → 写入 pending 队列，5号 轮询消费）—— 开工前确认走哪条
- [ ] 回调 → 下单：复用 Phase 3 签名路径，单仓、去重、超时作废（信号推出 60s 内未点击则按钮过期）
- [ ] 安全：按钮只对老板 chat_id 生效，callback_data 带 nonce 防重放
- [ ] Paper 期间：点了按钮走 paper 路径（记模拟 fill）；Day 7 转实盘后按钮直接真下单
- [ ] 告警：下单成功/失败都回 DM 单条小回执

### Phase 4 — 风控 + 可观测（1 天）
- [ ] 日亏损熔断
- [ ] WSS 断线保护
- [ ] telegram webhook 推送关键事件
- [ ] 日结报表 cron

### Phase 5 — Paper 跑 7 天
- [ ] Day 1-3：数据/信号合理性检查
- [ ] Day 4-7：策略参数微调
- [ ] Day 7：出 paper 报表 → 老板审 → 打款 → 实盘

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

## ❌ 不做

- 不接 1号 派的 python 活
- 不碰 python 项目任何文件
- 不依赖外部数据源
