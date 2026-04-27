# polymarket-go — 项目原则

本文件持久化老板拍过板的原则性决策与可复用经验。任何后续 agent / 未来的我，clone 下来第一件事读这份文档 + SPEC.md + TODO.md。

原则比代码稳定，先 commit 原则再动代码。

---

## P1. 项目完整隔离（2026-04-19）

polymarket-go 与 python 的 polymarket-agent 在以下所有维度必须物理隔离：

- **Repo**：独立 GitHub repo，不是子目录也不是 branch。
- **钱包**：独立助记词（存 Bitwarden `Polymarket-Go Wallet`），独立 Polygon 地址。**绝不复用** python 项目的钱包/私钥。
- **TODO**：独立 `TODO.md`，不碰 python 项目的 `PROJECT_TODO`。
- **订单数据**：下单走自签，不经过 python 的 order path；不污染 python 项目的订单簿/成交记录。
- **告警通道**：独立 bot token / chat 配置（虽然目的地都是同一个老板）。

**Why:** 隔离是为了出事能定位到人，以及不让 python 项目的订单/PnL 数据被这边的实验污染。老板 04-19 23:27 + 23:31 反复强调。

## P2. 自我迭代层（P00）优先于业务层（2026-04-19）

在写任何策略/数据/下单代码之前，先搭：

- **心跳**（heartbeat.sh）：定期自检 git/build/TODO/日志
- **周期驱动**（crontab + OpenClaw cron）：保证无人值守也能往前推
- **状态持久化**（state.json）：下一次醒来知道上一次在哪
- **告警升级**（alert-dispatch.sh）：出事能主动喊人，不是等被发现

**Why:** 5号 是 session-based，没有常驻循环。没有 P00 =一次 session 结束项目就停。业务层再好也白搭。老板 04-19 23:35 拍板"P00 必须先做完"。

**How:** 任何新项目复制这套模板：`scripts/{heartbeat,cron-poke,alert-dispatch}.sh` + `state.json` + crontab `*/20 * * * *` + 夜间静默。

## P3. 告警通道端到端自测（2026-04-19）

新接一个告警出口（Telegram / Slack / webhook），必须真发一次测试消息验证到达。不能假设"token 配了就能通"。

**Why:** 04-19 23:49 搭 alert-dispatch 时实测发现即使 bot token 和 chat_id 都配对，默认参数下消息不一定到。只有真发一次才确认 pipeline 通了。

**How:** `heartbeat.sh` 或等价脚本里留一个 `--self-test` 入口，新接出口时跑一次。正常运行态不触发。

## P4. 敏感凭据走 Bitwarden，不落 git（2026-04-19，复用 TOOLS.md）

- 助记词 / 私钥 / bot token / API key → Bitwarden `Agent` 文件夹
- 只允许 `.env.local`（chmod 600 + gitignored）缓存运行时需要的 token
- **绝不** commit `.env` / `*.key` / 助记词原文
- 接收到的助记词图片等，Bitwarden 存档后立刻从 inbound 目录删除

## P5. 群聊打断 ≠ 交付（跨项目，auto-memory 也有记录）

收到群聊/DM 消息时，回完必须在**同一 turn 内**继续推进主线工作（代码/commit/spec），不能只回话就 end turn 等下一轮。

**Why:** 5号 是 session-based 的，end turn = 失去推进动力。老板 04-18 15:26 拍板。

## P6. 大 TODO 完成要在当前老板私聊 DM 公示（2026-04-19 立规，2026-04-20 改正渠道）

Phase / P0x 级别的 TODO 完成时，在**当前老板 telegram 私聊 DM**（chat_id 6695538819）同步一句里程碑消息，不贴长报告。小条目不用。

**Why:** polymarket-go 是 5号 独立接的活，没有所谓"团队群"——草台班子团队群是 python 项目的，5号 04-19 23:40 已被踢出。汇报渠道就是当前 DM。老板 04-20 00:03 改正（原 04-19 23:58 版本误写"团队群"）。

**How:** commit 消息带 `milestone:` 前缀时，heartbeat 把它放到 state.json 的 milestones 字段，下一次 5号 醒来推送到当前 DM。（机制待实现，TODO.md P00 改进项。）

## P7. 原则持久化到 md，不只留对话里（2026-04-19，元规则）

老板反复强调的做法、拍板的决策、"这是原则"类表述，立刻写入：
- 跨项目行为偏好 → auto-memory `feedback_*`
- 本项目规则 → 本文件 `PRINCIPLES.md`
- 业务规格 → `SPEC.md`

**Why:** 对话会被 compaction 冲掉，md 是跨 session/跨 agent 稳定的承载。老板 04-19 23:58 拍板。

---

## P8. Paper 期初出场策略 = hold-to-settlement（2026-04-20 17:21）

Paper 阶段**不挂 SL/TP/timeout**。开仓后只等 market resolve，按 gamma `OutcomePrices[SlotIdx]` 清算（赢家 1.0 / 输家 0.0）。settlement watcher goroutine 每 60s 查 gamma。

**Why:** 老板 04-20 17:21 看 17:17 Day-1 样本后拍板——4/4 平仓全是 `reversal_drawdown` 几秒内止损，detector 门槛（Δ≥3pp）在抓短期顶点。在还没有足够样本校准阈值前，用 **"持到结果"** 给策略一个真正的 ground-truth PnL 分布，避免 exit 规则把好信号洗掉。

**How to apply:** 默认 `-exit_mode=hold`（bot-daemon.sh 已内置）。legacy `-exit_mode=auto` 走旧 ExitTracker（反转/回撤/3pp 止损/30min 超时），做阈值对照或手动压测时用。两种模式都保留日亏损熔断 + feed-silence watchdog——那是风控，不是策略。

## 变更日志

- 2026-04-19 — P1~P7 初版（5号 开工首日）
- 2026-04-20 17:29 — P8 hold-to-settlement 出场策略
- 2026-04-20 21:15 — P8 新增"infra 前先跑历史回测"
- 2026-04-20 21:30 — P9 新增"从邻居 DB 具体扬长避短"

---

## 8. 在建 infra 上烧钱前，先用历史数据做离线验证

**2026-04-20 21:1x 拍板。**

策略想法（追涨 / 宽 gap arb）在没有样本验证前容易自信。这个项目第一天就全 infra 堆栈铺开，但 Phase 6 方向讨论时才想起姐妹项目（python polymarket-agent）有 1325 条 theodds_h2h 快照 + 361 笔真实盘。一次 `cmd/backtest`（modernc.org/sqlite 只读打开）3 分钟产出硬数据：13/13 全败、-47.54% ROI、最大回撤 -101.76 USDC——直接推翻了 "PM vs bookmaker >5pp" 的主 edge 假设。

**规则**：

1. 每条新策略 SPEC 之前，先检查邻近项目 / 现有 DB / 现有日志里能不能凑出历史样本做离线验证
2. 历史数据即使只有 1-2 周，胜率 0/13 这种极端值就已经足以下决定不走这条路
3. 大 gap（>15pp）/ 大 delta（>10pp）这些极端信号大概率是 data pipeline 的 mismatch，不是真机会——**任何 gap 阈值上都必须先看 gap 的分布**

**不要**：不做这一步就直接 "先把 infra 搭好再边跑边看"。先验证想法，再烧代码。

---

## 9. 扬长避短 — 从邻居 DB 学具体而非抽象

**2026-04-20 21:30 拍板。**

第一次用 python DB 只问了 "PM vs bookmaker >5pp 能不能赚钱"，答案是不能就直接跳到 "不抄 python"。这是**抽象否定**——漏掉了具体模式里还有 8 笔赢家（ladder_TP 0.13→0.25, 0.36→0.60, 0.6455→0.895 …）值得学。

**规则**：

1. 否定一个策略后，**必须下钻**到 winners 清单和 losers 清单，按 entry_price / exit_reason / market_type 看具体长相
2. 拒绝 "整体 ROI 是负的就全部作废"——把**赢家共性**抽出来作为新 filter（本项目的 `reports/python_autopsy.md` §5 / §8 就是这一步的产物）
3. 同时把**输家共性**列成黑名单（本项目：theodds_h2h / football favorite 高价 / TIME_STOP 长守护）作为 prompt 过滤条件
4. 这条规则和 §8 互补：§8 要求用数据验证，§9 要求不只看聚合，必须看子集

**不要**：看到"总 PnL 为负"就把整个 DB 丢一边。

---

## P10. 主动日志审计 — 异常不等老板指出（2026-04-27）

**背景**：04-27 一天内出了 injury scanner 从未生效、BTC 长期盘反复推送、LCK Challengers 穿透过滤器、lottery 注额 $5→$1 不一致、仓位重启丢失等 5+ 个 bug，全部是老板先发现/指出才修。日志里其实都有明显异常信号（`injury_alert: 0`、`btc_strategy.signal_pushed` 不该出现、`lottery_open size_usd:5` 与配置不符），但没人看。

**规则**：

1. **每次心跳（cron-poke 20min）自动扫最近 20min 日志**，检查以下异常模式：
   - 不该出现的 event（如 `btc_strategy` 在 `-btc_enabled=false` 时出现）
   - 期望出现但缺失的 event（如 `injury_alert` 连续 1h 为 0）
   - 数值异常（如 `size_usd` 与配置不符、`pnl` 突变）
   - 错误率突增（`_err` / `_failed` event 密度）
2. **发现异常立刻写 TODO**（`BTC_TODO.md` 或 `TODO.md`），标 P0/P1，不等下次对话
3. **能自修的立刻修**（编译通过 + 测试通过 → commit + restart daemon），不能自修的推老板 DM
4. **每次 daemon 重启后 5min 内做一次 smoke check**：关键 event 是否正常出现

**频率**：
- cron-poke 每 20min 做基础异常扫描（grep 错误模式）
- 5号 session 每次醒来做完整日志审计（最近 1-2h）
- daemon 重启后 5min smoke check

**Why**：日志是 daemon 唯一的"嘴"——它不会主动喊人，只会默默记录。不看日志 = 聋了。老板 04-27 16:08 拍板"不应该等我指出去修复"。

---

## P11 · 核心功能必须落日志，不许静默 (04-27 19:45 老板拍板)

**铁律**：所有核心功能（开仓/平仓/信号/下单/风控/推送/数据拉取/状态持久化）**必须有日志输出**。错误不许吞（`_ = someFunc()` 禁止用于有副作用的操作）。

**具体要求**：

1. **每个 error 路径必须 `slog.Warn`**：函数返回 error 后，调用方必须 log，不允许 `if err != nil { return }` 不带日志
2. **状态持久化必须 log 失败**：`SaveState` / `WriteFile` / `Marshal` 失败 → `slog.Warn`，不允许 `_ =` 丢弃
3. **外部 API 调用必须 log 失败**：HTTP 请求错误、非 200 状态码、JSON 解析失败 → 全部 `slog.Warn` 带上下文
4. **数据丢弃必须可见**：sampler tick drop、WSS 帧解析失败、dedup 跳过 → 用 counter 或 warn 让运维看到
5. **回调错误不许吞**：`OnSent` / `OnClose` 等回调中 err != nil → `slog.Warn`
6. **新写代码必须过审计**：每个新函数写完后，自查"如果这个操作失败了，日志里看得到吗？"

**检查清单**（代码审计用）：
- `_ = ` 赋值 → 目标函数有无副作用？有 → 必须改为 `if err := ...; err != nil { slog.Warn(...) }`
- `if err != nil { return }` → 上面有没有 `slog.Warn`？没有 → 加
- 函数超过 20 行 → 有没有至少一个 `slog` 调用？没有 → 考虑加
- `continue` 在循环里 → 跳过了什么？值得 log 吗？

**Why**：问题静默 = 问题不存在（对运维而言）。04-27 一天发现 5 个静默 bug（injury scanner 从未工作、lottery 注额错误、仓位重启丢失），全部因为没有日志才拖到老板问了才查。

---
