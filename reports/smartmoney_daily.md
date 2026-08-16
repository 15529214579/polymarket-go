# Smart Money Daily

**Generated:** 2026-08-16 00:17:31 +08

## Pipeline

- Status: learning
- Reason: evaluated signals 0 below 10
- Action: keep current push running and collect more signals
- Log: `/Users/murphyma/work/polymarket-go/logs/smartmoney-daily-2026-08-16.log`
- Discovery report: `/Users/murphyma/work/polymarket-go/reports/strategy_iteration.md`
- Strategy report: `/Users/murphyma/work/polymarket-go/reports/strategy_lab.md`
- Sports tape report: `/Users/murphyma/work/polymarket-go/reports/sports_tape.md`
- Sports alert performance report: `/Users/murphyma/work/polymarket-go/reports/sports_alert_performance.md`
- Sports alert candidate report: `/Users/murphyma/work/polymarket-go/reports/sports_alert_candidates.md`
- Sports burst performance report: `/Users/murphyma/work/polymarket-go/reports/sports_burst_performance.md`
- Whale edge report: `/Users/murphyma/work/polymarket-go/reports/whale_edge.md`
- Core performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_core.md`
- Sports performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_sports.md`
- Scout performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_scout.md`
- Target performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_target.md`
- Flow performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_flow.md`
- Tape performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance_tape.md`
- Push performance report: `/Users/murphyma/work/polymarket-go/reports/whale_performance.md`
- Paper PnL report: `/Users/murphyma/work/polymarket-go/reports/smartmoney-paper-pnl.md` (status 0)
- Paper wallet policy: `/Users/murphyma/work/polymarket-go/reports/smartmoney-paper-wallets.md` (status 0, changed no, worker not needed)
- Paper exit shadow: `/Users/murphyma/work/polymarket-go/reports/smartmoney-exit-shadow.md` (status 0)
- Maintenance report: `/Users/murphyma/work/polymarket-go/reports/wallet_maintenance.md`

## Wallet Lists

- Core before count: 14
- Core after count: 14
- Watch after count: 20
- Sports after count: 10
- Scout after count: 0
- Target after count: 10
- Flow after count: 0
- Tape after count: 0
- Tape observe count: 1
- Tape probation count: 0
- Tape candidate count: 0
- Tape follow-ready count: 0
- Tape edge-hot count: 2
- Tape reversal count: 8
- Consensus research count: 8
- Push before count: 48
- Push after count: 46
- Quarantine count: 3
- Review-noise exclude count: 404
- Edge snapshots: 1282
- Sports tape alert sent: 5
- Sports tape alert logged: 5
- Sports tape shadow alert logged: 7
- Sports consensus event history: 4
- Sports consensus watch event history: 4
- Sports alert eligible now: 0
- Sports alert accumulation bursts: 0
- Sports alert consensus bursts: 0
- Selected core wallets: 14
- Core list changed: yes
- Push list changed: yes
- Restart needed: no
- Restart status: restarted screen polymarket-whale-push

## Backtest Filter

- Closed copy trades: 295
- Copy ROI: 116.1%
- Copy PnL: $+3947.33
- Copy win rate: 83.7%
- Worst included CopyROI: 101.3%
- Params: tier>=B bot<30 copyT>=8 copyROI>=100 copyPnL>=25 copyWin>=60 closedROI>=0 smart>=70

## Live Core Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Sports Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Sports Alert Performance

- Alerts: 5
- Marked to current midpoint: 5
- Unmarked: 0
- Win rate incl. midpoint marks: 20.0%
- PnL incl. midpoint marks: $-34.10
- ROI incl. midpoint marks: -68.2%

## Sports Alert Current Policy

- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert Position-Capped Policy

- Rule: first alert per wallet + asset
- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Positions: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert OBSERVE Experiment

- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted
- Default min notional: $3000
- Observe-burst min notional: $6000
- Require scored/listed wallet: true
- Minimum tier: B
- Insider-scout min notional: $25000
- Insider-scout max bot: 35
- Alerts: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Position-capped entries: 0
- Position-capped ROI: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Alert OBSERVE-BURST Experiment

- Rule: same-wallet split-order sports/esports BUY bursts; observation only until repeated positive ROI is proven
- Min cumulative notional: $6000
- Alerts: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Position-capped entries: 0
- Position-capped ROI: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Alert SHADOW OBSERVE-BURST Experiment

- Rule: stale same-wallet split-order bursts captured for research only; not counted as Telegram pushes or current policy
- Shadow log: /Users/murphyma/work/polymarket-go/db/strategy_iteration/sports_tape_shadow_alerts.jsonl
- Logged events: 7
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Position-capped entries: 1
- Position-capped ROI: 61.3%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert SHADOW CONSENSUS Experiment

- Rule: stale cross-wallet same-asset sports/esports BUY bursts captured for research only; not counted as Telegram pushes or current policy
- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Alerts: 4
- Marked to current midpoint: 4
- Marked samples still needed: 1
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $+1.38
- ROI incl. midpoint marks: 3.5%
- Position-capped entries: 3
- Position-capped ROI: 18.3%
- Readiness: COLLECT (1 more marked samples before promotion review)
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert SHADOW SEED-FLOW Experiment

- Rule: lower-threshold unknown wallets buying multiple sports/esports markets; research-only before UNKNOWN-FLOW size is reached
- Min cumulative notional: $3000
- Min markets: 2
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+6.13
- ROI incl. midpoint marks: 61.3%
- Position-capped entries: 1
- Position-capped ROI: 61.3%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert SHADOW SCORED-FLOW Experiment

- Rule: scored low-bot leaderboard wallets buying multiple sports/esports markets; shadow-only until repeated positive ROI is proven
- Min cumulative notional: $4000
- Min markets: 2
- Min tier: B
- Max bot: 35
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-10.00
- ROI incl. midpoint marks: -100.0%
- Position-capped entries: 1
- Position-capped ROI: -100.0%
- Gate action: PROBATION
- Gate reason: severe drawdown on limited sample

## Sports Alert INSIDER-SCOUT Experiment

- Rule: 25k+ very large low-bot sports/esports whale BUY alerts; observation only until repeated positive ROI is proven
- Alerts: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Position-capped entries: 0
- Position-capped ROI: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Alert CONSENSUS Experiment

- Rule: cross-wallet same-asset sports/esports BUY bursts; research/observation only until repeated positive ROI is proven
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.90
- ROI incl. midpoint marks: 59.0%
- Position-capped entries: 1
- Position-capped ROI: 59.0%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Burst Performance

- Bursts: 0
- Marked to current midpoint: 0
- Unmarked: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Avg price delta: +0.00pp

## Live Scout Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Target Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Flow Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Tape Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Push Performance

- Evaluated signals: 0
- Event-capped entries: 0
- Policy violations: 0
- Proven signals: 0
- Proven event-capped entries: 0
- Logged asset-cooldown BUYs: 0
- Logged event-cooldown BUYs: 0
- Logged duplicate BUYs: 0
- Realized via whale SELL: 0
- Settled by market resolution: 0
- Marked to current midpoint: 0
- Still open/unmarked: 0
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Proven PnL: $+0.00
- Proven ROI: 0.0%
- Event-capped realized via whale SELL: 0
- Event-capped settled by market resolution: 0
- Event-capped marked to current midpoint: 0
- Event-capped PnL incl. midpoint marks: $+0.00
- Event-capped ROI incl. midpoint marks: 0.0%
- Event-capped proven PnL: $+0.00
- Event-capped proven ROI: 0.0%
