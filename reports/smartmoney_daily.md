# Smart Money Daily

**Generated:** 2026-07-09 11:44:41 +08

## Pipeline

- Status: learning
- Reason: evaluated signals 0 below 10
- Action: keep current push running and collect more signals
- Log: `/Users/murphyma/work/polymarket-go/logs/smartmoney-daily-2026-07-09.log`
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
- Maintenance report: `/Users/murphyma/work/polymarket-go/reports/wallet_maintenance.md`

## Wallet Lists

- Core before count: 10
- Core after count: 10
- Watch after count: 20
- Sports after count: 6
- Scout after count: 1
- Target after count: 6
- Flow after count: 0
- Tape after count: 0
- Tape observe count: 5
- Tape probation count: 0
- Tape candidate count: 0
- Tape follow-ready count: 0
- Tape edge-hot count: 3
- Tape reversal count: 8
- Consensus research count: 9
- Push before count: 39
- Push after count: 39
- Quarantine count: 3
- Review-noise exclude count: 4
- Edge snapshots: 254
- Sports tape alert sent: 4
- Sports tape alert logged: 4
- Sports tape shadow alert logged: 3
- Sports consensus event history: 3
- Sports consensus watch event history: 2
- Sports alert eligible now: 0
- Sports alert accumulation bursts: 0
- Sports alert consensus bursts: 1
- Selected core wallets: 10
- Core list changed: no
- Push list changed: no
- Restart needed: no
- Restart status: not requested

## Backtest Filter

- Closed copy trades: 182
- Copy ROI: 65.9%
- Copy PnL: $+1442.89
- Copy win rate: 75.3%
- Worst included CopyROI: 42.0%
- Params: tier>=B bot<30 copyT>=8 copyROI>=10 copyPnL>=50 copyWin>=60 closedROI>=0 smart>=70

## Live Core Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Live Sports Performance

- Evaluated signals: 0
- PnL: $+0.00
- ROI: 0.0%

## Sports Alert Performance

- Alerts: 4
- Marked to current midpoint: 4
- Unmarked: 0
- Win rate incl. midpoint marks: 50.0%
- PnL incl. midpoint marks: $-9.75
- ROI incl. midpoint marks: -24.4%

## Sports Alert Current Policy

- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Alerts: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Alert Position-Capped Policy

- Rule: first alert per wallet + asset
- Modes: CANDIDATE,CONSENSUS,EDGE-HOT,FLOW-SCOUT,FOLLOW-READY,PROBATION
- Positions: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Alert OBSERVE Experiment

- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted
- Default min notional: $5000
- Observe-burst min notional: $8000
- Require scored/listed wallet: true
- Minimum tier: B
- Insider-scout min notional: $25000
- Insider-scout max bot: 35
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $-0.11
- ROI incl. midpoint marks: -1.1%
- Position-capped entries: 1
- Position-capped ROI: -1.1%
- Gate action: COLLECT
- Gate reason: sample below promote gate

## Sports Alert OBSERVE-BURST Experiment

- Rule: same-wallet split-order sports/esports BUY bursts; observation only until repeated positive ROI is proven
- Min cumulative notional: $8000
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
- Logged events: 3
- Alerts: 1
- Marked to current midpoint: 1
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+21.10
- ROI incl. midpoint marks: 211.0%
- Position-capped entries: 1
- Position-capped ROI: 211.0%
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

## Sports Alert SHADOW CONSENSUS Experiment

- Rule: stale cross-wallet same-asset sports/esports BUY bursts captured for research only; not counted as Telegram pushes or current policy
- Promote gate: marked>=5, ROI>=5.0%, win>=60.0%
- Alerts: 2
- Marked to current midpoint: 2
- Marked samples still needed: 3
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+5.46
- ROI incl. midpoint marks: 27.3%
- Position-capped entries: 2
- Position-capped ROI: 27.3%
- Readiness: COLLECT (3 more marked samples before promotion review)
- Gate action: COLLECT_POSITIVE
- Gate reason: positive but sample below promote gate

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
- Alerts: 0
- Marked to current midpoint: 0
- Win rate incl. midpoint marks: 0.0%
- PnL incl. midpoint marks: $+0.00
- ROI incl. midpoint marks: 0.0%
- Position-capped entries: 0
- Position-capped ROI: 0.0%
- Gate action: COLLECT
- Gate reason: no marked alerts yet

## Sports Burst Performance

- Bursts: 1
- Marked to current midpoint: 1
- Unmarked: 0
- Win rate incl. midpoint marks: 100.0%
- PnL incl. midpoint marks: $+1.38
- ROI incl. midpoint marks: 13.8%
- Avg price delta: +11.68pp

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
