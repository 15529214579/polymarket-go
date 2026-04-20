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
- [x] 2026-04-20 15:36 — Phase 0 golangci-lint 配置（非阻塞）：`.golangci.yml` v2 格式，启用 errcheck/govet/ineffassign/staticcheck/unused/gocritic/gosec/misspell/unconvert + gofmt/goimports 格式化。gosec 排除 G104/G301/G302/G304/G306/G704/G706（配置驱动路径 + 自家 API HTTP）；`_test.go` 排除 errcheck/gosec；`cmd/bot/main.go` 排除 gocritic（wiring-heavy）。13 个现有文件 gofmt 过一遍。commit `05a90b4` 已入。

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
- [x] 2026-04-20 12:09 — 老板拍板选 **B（sidecar bot long-poll）**。
- [x] 2026-04-20 12:1x — Phase 3.5.b 代码层完成：
  - `internal/notify/callback.go` 上线：`LongPoll` 消费 `getUpdates`（allowed_updates=callback_query），chat_id 白名单，每个点击都走 `answerCallbackQuery` 回 toast；`CallbackHandler` 接口
  - `cmd/bot/main.go` 加 `-signal_mode auto|prompt` flag（prompt 模式下信号→DM+pending，callback 执行）；`buyHandler` 负责 Claim/risk 门控/dedupe/`paper.Submit`/`pm.OpenSized`/`exit.Open`/toast
  - `strategy.PositionManager` 新增 `OpenSized(sizeUSD)`，保留所有 dedupe+敞口上限
  - 只有 `SIDECAR_BOT_TOKEN` 存在时 longpoll 才启动（防止和 OpenClaw 抢 `TELEGRAM_BOT_TOKEN` 的 update）
  - 新 callback 单测 4 个（parse/分派成功/跨 chat 拒/坏 data/handler err → toast）全过；`./bin/bot -mode=detect -signal_mode=prompt` 启动 smoke OK
- [x] 2026-04-20 12:15 — Sidecar bot token 入档：老板 BotFather 建 `@Murphyoderbot` (id 8760736438)，token 存 Bitwarden `Polymarket-Go Sidecar Bot` + 写 `.env.local` `SIDECAR_BOT_TOKEN=...`（`SIDECAR_CHAT_ID` 复用 `TELEGRAM_CHAT_ID=6695538819`）
- [x] 2026-04-20 12:17 — 修 Phase 3.5.b 接线 bug：`SignalPrompt` 原先走 alert bot `@Murphy005bot` 发消息，点击 callback 进 alert bot 队列，但 LongPoll 盯的是 sidecar → 永远收不到点击。修法：`TelegramConfig.PromptBotToken` 字段 + `outgoing.sendToken` per-message 覆盖，prompt 走 sidecar，其它事件继续走 alert bot。`TestTelegram_SignalPromptRoutesThroughPromptBot` 锁死。
- [x] 2026-04-20 12:20 — 端到端实盘验收：`./bin/bot -mode=prompt-test` 从 @Murphyoderbot 发 Buy 1U/5U/10U 给 `LCK CL Gen.G vs Nongshim BO3` (mid 0.50)；老板点 Buy 1U，14s 内 `callback_click` + `manual_open` id=p1 order=paper-eeff0905 2 units @ 0.50，toast 回 `✅ 1U @ 0.5000 · order paper-e..0905`。sidecar → pending.Claim → paper.Submit → pm.OpenSized → exit.Open 整条回路通。
- [x] 2026-04-20 12:35 — Phase 3.5.c 按钮 UX 改版（老板反馈 "应该给我买 yes 和 no 的选项；赛事信息不全"）：`SignalPromptEvent` 拆出 `Match/Context/EndIn/Choices[]`，`PendingIntent.Choices []Choice{AssetID,Outcome,Mid,IsSignal}`，callback_data 扩成 `buy:<nonce>:<slot>:<sizeUSD>`；`buildAssetMeta` 从 gamma 一次性做 asset→{Question,Match,Context,Outcome,Sibling,EndTime} 索引，prompt 渲染时 signal 行 ⚡ 标记 + 对手盘行并列。`notify` 包 4 个测试（ParseMarketTitle/HumanizeEndIn/FormatSignalPrompt/InlineKeyboard 2 行 3 列）+ callback parse 4-part 格式全绿。实盘验收：Gen.G vs Nongshim BO3 prompt 两行（⚡ Gen.G / Nongshim），老板点 slot 0 × 1U，`callback_click nonce=5b4fa5be slot=0 size=1` → `manual_open p1 outcome=Gen.G Global Academy order=paper-8de4f5598b3a` 18 秒闭环。
- [x] 2026-04-20 12:45 — Phase 3.5.e：按钮有效期 60s→10m（老板反馈 "可能没那么及时看群消息"）。`NewPendingStore` + `SignalPromptEvent.ExpiresIn` 改成 10 分钟，`FormatSignalPrompt` 末行渲染走 `humanizeTTL`（≥1m 用 "10m"，<1m 走 "30s"）。`notify_test.go` 对齐。commit `7422206` 已 push。
- [x] 2026-04-20 12:45 — Phase 3.5.f：只推信号侧按钮（老板反馈 "你觉得对的让我选金额就行"）。`FormatSignalPrompt` 去掉对手盘 "选 X (当前 ...)" 段，正文尾行改成 `当前 0.xxxx`；`telegram.SignalPrompt` 只渲染 `IsSignal=true` 的 choice 行（1U/5U/10U）。`signalChoice()` 辅助在 notify.go。`notify_test.go` / `telegram_test.go` 锁死"only signal shown"。
- [x] 2026-04-20 15:36 — Phase 3.5.g：超时作废视觉升级。`PendingStore.Reap` 返回 `[]PendingIntent`（原 int）→ reaper 循环对每条有 `MessageID` 的条目调 `notifier.EditSignalExpired`，正文改写成 "⌛ 已过期 · 未下单" + 清空 `reply_markup.inline_keyboard`。`PendingIntent.MessageID int64` 新字段，`PendingStore.SetMessageID(nonce, id)` 由 Telegram drain 回调异步填入。`SignalPromptEvent.OnSent func(messageID int64, err error)` 暴露 message_id 给主程序（runDetect + prompt-test 两处设置 OnSent → `pending.SetMessageID`）。Telegram.send 现在解析 `result.message_id`，支持 `editMessageText` 路径（`outgoing.editMessageID` + `stripKeyboard`）；edits 经 prompt bot 发（和原 DM 同一只）。Notifier 接口扩了 `EditSignalExpired(int64)`，Nop noop。新单测：`TestPendingStore_SetMessageID` + `TestTelegram_SignalPrompt_OnSent_ReportsMessageID` + `TestTelegram_EditSignalExpired_StripsKeyboard` + `TestTelegram_EditSignalExpired_ZeroIsNoOp` 全过。
- [x] 2026-04-20 15:36 — Phase 3.5.h：成交凭据留档（C）。Notifier 接口扩了 `EditSignalFilled(FillReceiptEvent, int64)` + `FillReceipt(FillReceiptEvent)`（Nop noop）。`buyHandler.OnBuy` 在 `manual_open` 成功后调用两者：edit 原 prompt → "✅ 已下单 · `<outcome>` `<size>U` @ `<px>`" + 清空键盘（走 prompt bot）；并额外发一条 "🧾 成交凭据 · 手动" DM（走 alert bot，留档）。Telegram.FillReceipt / EditSignalFilled 实现 + `TestTelegram_EditSignalFilled_Renders` + `TestTelegram_FillReceipt_GoesToAlertBot` 全过。
- [ ] Paper 期间：点了按钮走 paper 路径；Day 9 起自动走真下单（Phase 3 V2 签名 client 接同一 `order.Client` 接口）

### Phase 4 — 风控 + 可观测（1 天）
- [x] 2026-04-20 11:4x — Phase 4.a：`internal/risk/risk.go` 上线。日亏损熔断（15% × 90.41 = -13.56 USDC 上限）+ per-trade 单笔损失 ≥ 3 USDC 计数旗标 + WSS feed-silence watchdog（30s 无 book/trade → trip）+ 手动 Pause/Resume。SGT 日切滚动但不自动解除 breaker（SPEC §6 "等老板手动恢复"）。8 个单测全过。接进 detect 循环：开仓前 `AllowOpen` 门控，close 后 `OnClose` 累计，5s 心跳 `CheckFeed`，60s `risk_status` 日志。35s 实盘烟测：`risk.ready` + 无 trip（预期，LoL 市场平静）。
- [x] 2026-04-20 11:55 — Phase 4.b：`internal/notify/` 上线（Notifier interface + TelegramNotifier + Nop）。事件：`RiskTrip`（daily_loss / feed_silence / manual_pause，每次 trip 一次）+ `RiskResume`（breaker 恢复）+ `LargeFill`（|realized_pnl| ≥ `-large_fill_usd`，默认 3 USDC）。异步 goroutine queue（32 buf，满则 drop 不阻塞交易）。detect 循环两个 trip site + close site 都接好；resume 在 watchdog 跟踪 prevBlocked→!blocked 转换触发。`internal/config/dotenv.go` 辅助（启动加载 .env.local，不覆盖已有 env）。5 notify 单测 + 2 dotenv 单测全过，实盘 smoke ok (`notify.ready mode=telegram`)。
- [x] 2026-04-20 13:53 — Phase 4.c：日结报表 cron。`internal/journal/` 包（JSONL 持久化 + Summarize + FormatTelegram）；detect 关仓后 append（auto/manual 都标 source）；`-mode=daily-report [-report_day=YYYY-MM-DD] [-report_push]` 读取 SGT 当日 JSONL → 渲染 → 可选 Telegram alert bot 推送；`scripts/daily-report.sh` 包装 + crontab `0 0 * * *` SGT；6 单测全过；wrapper 烟测推 04-19 "无成交" + 04-20 seed 数据 PnL/胜率/持仓/出场原因渲染都对。

### Phase 5 — Paper 跑 8 天（Apr 20 → Apr 28 cutover）
> 04-20 10:36 老板拍板 A：顺延 Paper 跨过 V2 cutover，Apr 29 起实盘（原 7 天 → 8 天）
- [x] 2026-04-20 13:54 — Day-1 paper detect daemon 起飞：`scripts/bot-daemon.sh start` → `-mode=detect -signal_mode=auto -markets=20 -window=60`，logs 进 `db/agent.{log,err}`，pidfile `db/bot.pid`。`cron-poke.sh` 每 20 min 自动重启如果挂掉。startup ok：detect.start markets=20 assets=40，wss connected，risk 0/13.56 cap，notify telegram，sidecar long-poll ready。
- [x] 2026-04-20 16:03 — **多体育扩容**（老板 15:57 问"lol 数据不够，足球篮球也加上"）：`feed.FilterSports` 扩到 LoL + NBA daily/playoffs + EPL daily（moneyline 独占，`-spread-/-total-/-ou-/-over-/-under-/-prop-/-parlay-` 后缀一律排除）。gamma.discover 从 42 条 LoL → 111 条 sports；当前 top-20 订阅分布 `lol=5 / basketball=11 / football=4`。gamma_test.go 新增 LoL/NBA/EPL/union 4 套测试覆盖 seasonal futures 排除、derivative 排除、union 顺序。daemon 已用新 bin 重启（pid=69791，detect.start 绿）。
- [x] 2026-04-20 16:03 — **heartbeat state.json 即时刷新**（老板 A：修 last_commit 20min 滞后 race）：`scripts/install-hooks.sh` + `.git/hooks/post-commit` 上线，每次 commit 后异步 disown `cron-poke.sh` 立刻重写 state.json，0s 滞后；`make hooks` 一键重装。不阻塞 commit UX（nohup + disown + 吞 error）。
- [x] 2026-04-20 17:29 — **出场策略改 hold-to-settlement**（老板 17:21：不加 SL/TP，买了等最终结果）：新 `-exit_mode=hold|auto`（默认 hold）；hold 模式跳过 `exit.Open`，新增 settlement watcher goroutine 每 60s 查 gamma `GetByConditionIDs`，`closed=true` 即按 `OutcomePrices[SlotIdx]` 清算平仓（1.0 赢家 / 0.0 输家）；同样走 pm.Close → risk.OnClose → large-fill/risk-trip notify → journal。`assetMeta` 加 ConditionID + SlotIdx；新 `ExitSettlement` reason；两条新单测（httptest gamma mock + empty-input）。daemon 已切到 hold，`scripts/bot-daemon.sh` 默认 `-exit_mode=hold`。
- [ ] Day 1-3（Apr 20-22）：信号密度 + 假阳性观察（每 4-6h 看一次 db/agent.log，凑足 20+ 信号样本）
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
