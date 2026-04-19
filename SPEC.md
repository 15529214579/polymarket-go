# polymarket-go — SPEC

**Owner:** 5号 (monitor)  
**Repo:** (待建) github public  
**Wallet:** `0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e` (Bitwarden → `Polymarket-Go Wallet`)  
**Status:** draft — 2026-04-19

---

## 1. 项目目标

用 Go 重写 Polymarket 交易代理，**不污染** python 项目。独立钱包、独立 repo、独立 TODO。

**为什么 Go：** 低延迟、goroutine 并发、持久 WSS 长连接稳定。

## 2. 策略（老板 2026-04-19 23:33 定调）

**类型：** 赛中动量跟进（in-play momentum follow）

**触发条件：**
- 必须是**已开赛**（game_state = live）的 LoL 比赛市场
- **任一侧** yes/no 价格出现"持续上升利好"信号 → 跟进该侧
- **不限**初始价格区间（0.3 也可以进，0.8 也可以进）

**"上升利好"量化（默认，待老板 review）：**
- 最近 N 秒（默认 **60s**）价格净涨幅 ≥ **3pp**（百分点）
- 且最近 M 个 tick（默认 **5**）收盘价单调或准单调上升（至少 4/5 上行）
- 且 orderbook 买一方向主动成交占比 ≥ 60%（避免挂单堆叠误导）

**出场：**
- 止盈：持有至反转信号（3 个 tick 连续下行 或 净回撤 2pp）
- 止损：入场价 -3pp 硬止损
- 结算：市场 resolve 时按最终结果清算
- 最大持仓：单笔 **30 分钟**上限（比赛慢节奏时强制出）

**仓位：**
- Paper 阶段：单笔 **5 USDC**，不叠仓，同一市场同时只持 1 仓
- 实盘阶段：总资金 × 5%/笔，由老板打钱后定

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

## 5. 生命周期

- **Day 0-7：** Paper trade（只记录不发单，对照实盘价格模拟 PnL）
- **Day 7：** 老板 review paper 结果 → 打款到新钱包 → 实盘
- **实盘上限：** 老板首轮打款额度定（建议 50-100 USDC 起步）

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
3. 实盘首轮打款额度
4. 是否需要 Discord/其他告警冗余（目前只推 telegram 私聊）
