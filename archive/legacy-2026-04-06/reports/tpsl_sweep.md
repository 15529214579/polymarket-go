# TP/SL Sweep — python closed trades endpoint replay

Run: `./bin/backtest -mode=tpsl-sweep -fee_bp=<0|20|100>`
Source: `~/.openclaw/workspace-dev3/polymarket-agent/db/polymarket_agent.db`
Sample: 183 closed trades (YES+NO merged via `pnl/capital` → natural return)
Enriched subset: 9 trades joinable to `odds_snapshot` for peak/trough proxy.

## Baseline (no TP/SL)

| fee/leg | avg_ret | hit% | sum_ret | mdd |
| ------- | ------- | ---- | ------- | --- |
| 0 bp    | -16.62% | 4.4  | -30.41  | 46.50 |
| 20 bp   | -17.02% | 4.4  | -31.14  | 47.19 |
| 100 bp  | -18.62% | 4.4  | -34.07  | 49.96 |

**Observation**: natural hit rate is 4.4% (8/183). Python's blended book was a losing pool by a wide margin.

## Top endpoint configs (fee=0 bp)

| TP%  | SL%  | tp_hit | sl_hit | natur | sum_ret | avg_ret | mdd |
| ---- | ---- | ------ | ------ | ----- | ------- | ------- | --- |
| 100  | 5    | 6      | 151    | 26    | -0.23   | -0.13%  | 7.5 |
| 75   | 5    | 6      | 151    | 26    | -1.73   | -0.95%  | 7.5 |
| 50   | 5    | 8      | 151    | 24    | -3.56   | -1.94%  | 7.5 |
| 40   | 5    | 8      | 151    | 24    | -4.36   | -2.38%  | 7.5 |
| 30   | 5    | 8      | 151    | 24    | -5.16   | -2.82%  | 7.5 |

## Bottom endpoint configs (fee=0 bp)

| TP%  | SL%  | sum_ret | mdd |
| ---- | ---- | ------- | --- |
| 5    | 0    | -48.96  | 49.25 |
| 10   | 0    | -48.56  | 49.15 |
| 15   | 0    | -48.16  | 49.05 |
| 20   | 0    | -47.76  | 48.95 |
| 30   | 0    | -46.96  | 48.75 |

No SL + any TP = worse than baseline. SL is the dominant loss-limiting lever.

## Findings

1. **Tight SL dominates**: Every top-10 config pins `SL=5%`. Moving to `SL=10%` roughly doubles the sum loss (-7.78 vs -0.23 at TP=100).
2. **TP impact is small**: Only 6-8/183 trades reach even 5% endpoint return. Setting TP lower doesn't boost wins, it just caps the few winners that exist.
3. **Fee sensitivity is linear**: 100 bp/leg (round-trip 2%) pushes the best config from -0.23 → -3.89 — still within striking distance of breakeven.
4. **Direction indifference**: natural return is computed as `pnl/capital`, so TP/SL thresholds are side-agnostic.

## Data limits (hard)

- **No tick path**: python DB stores entry + close + pnl only. Endpoint replay under-counts TP hits — price can touch TP intraday and retrace below TP by close; we'd miss that TP fire.
- **`odds_snapshot` join**: only 9/183 trades overlap with their `odds_snapshot` window. Path-aware sweep on n=9 is too small for a conclusion (all 9 show SL touched with the conservative "both touched → SL wins" rule).
- **Different strategy pool**: python trades are primarily `theodds_h2h` + `llm_scan` entries, not PM-internal momentum. Thresholds are directional guidance, not absolute calibration for our momentum strategy.

## Recommendations

Directional, not final:

1. **Tighten SL in SPEC**: current `StopLossPct=0.30` (30% adverse) is likely too loose. Try `SL=0.05-0.10` (5-10%) in SPEC, keep ladder TPs roomy (`TP1=0.30, TP2=0.60` or wider).
2. **Keep TP ladder loose**: aggressive TPs rarely trigger; clipping winners via tight TP is strictly worse in this data.
3. **Phase 7.e is the real fix**: persist 1 Hz sampler ticks to disk per open paper position. After 3-7 days of paper, we have dense path data for OUR entry logic — then re-run sweep with path-aware replay. The current 9-trade path sample is meaningless.
4. **Fee budget**: design for round-trip ≤ 100 bp (2%). Anything higher makes even best-case config (-3.89) a negative drift.

Not proposing auto-applying these — needs boss sign-off because (a) python sample is a different strategy and (b) endpoint replay under-counts TP. `SL=5%` in particular would feel tight on paper where we're deliberately collecting wide-range samples.
