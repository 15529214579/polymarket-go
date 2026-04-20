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
