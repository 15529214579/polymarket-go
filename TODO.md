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

### Phase 0 — Bootstrap（1 天内）
- [x] `go mod init github.com/murphyismurphy/polymarket-go`
- [x] 目录骨架：`cmd/bot/`, `internal/{feed,strategy,order,risk,log,config}/`
- [x] Makefile + .gitignore + build 通过
- [x] git init + 首个 commit（3d072a7）
- [ ] 建 github public repo（gh auth 后 push）
- [ ] golangci-lint 配置

### Phase 1 — 数据层（下一步，2-3 天）
- [ ] Polymarket WSS 客户端（自动重连、心跳）
- [ ] gamma REST 客户端（LoL 市场筛选）
- [ ] orderbook 内存模型（bid/ask 深度、最近成交流）
- [ ] tick 采样器（1s 粒度，滑窗 60s）

## 💤 待启动

### Phase 2 — 策略层（2 天）
- [ ] 动量信号检测（N秒涨幅、tick 单调性、主动成交占比）
- [ ] 出场信号（反转、止损、超时）
- [ ] 仓位管理（单仓去重、总敞口）

### Phase 3 — 下单（1-2 天，方案 A：自签+broadcast）
- [ ] Bitwarden 取助记词 → 派生私钥（启动时只驻内存）
- [ ] EIP-712 typed data 签名（Polymarket CLOB order struct）
- [ ] CLOB REST `/order` POST 客户端
- [ ] 成交回执轮询 + status 机
- [ ] Paper mode（同一路径但不发真单，记模拟 fill）

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

## ❌ 不做

- 不接 1号 派的 python 活
- 不碰 python 项目任何文件
- 不依赖外部数据源
