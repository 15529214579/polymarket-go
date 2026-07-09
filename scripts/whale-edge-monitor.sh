#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

LOG="${WHALE_EDGE_LOG:-$ROOT/logs/whale-edge-monitor.log}"
INTERVAL="${WHALE_EDGE_INTERVAL:-60}"

mkdir -p "$ROOT/logs" "$ROOT/reports" "$ROOT/db/strategy_iteration"

go build -o bin/whale-edge ./cmd/whale-edge

exec >> "$LOG" 2>&1
echo "whale-edge-monitor.start pid=$$ interval=${INTERVAL}s"

while :; do
  date '+===== %Y-%m-%d %H:%M:%S %Z ====='
  "$ROOT/bin/whale-edge" \
    -log "${WHALE_LOG:-$ROOT/db/journal/whale_trades.jsonl}" \
    -tape_log "${WHALE_EDGE_TAPE_LOG:-$ROOT/db/strategy_iteration/sports_tape.jsonl}" \
    -wallets "${WHALE_EDGE_WALLETS:-$ROOT/wallets.strategy-push.txt,$ROOT/wallets.strategy-tape.txt,$ROOT/wallets.strategy-tape-observe.txt,$ROOT/wallets.strategy-tape-probation.txt,$ROOT/wallets.strategy-tape-candidates.txt,$ROOT/wallets.strategy-tape-follow.txt,$ROOT/wallets.strategy-tape-reversal.txt}" \
    -snapshots "${WHALE_EDGE_SNAPSHOTS:-$ROOT/db/strategy_iteration/whale_edge_snapshots.jsonl}" \
    -report "${WHALE_EDGE_REPORT:-$ROOT/reports/whale_edge.md}" \
    -actions "${WHALE_EDGE_ACTIONS:-alert,followed}" \
    -horizons "${WHALE_EDGE_HORIZONS:-0m,5m,15m,30m,60m}" \
    -tolerance "${WHALE_EDGE_TOLERANCE:-2m}" \
    -min_notional "${WHALE_EDGE_MIN_NOTIONAL:-500}" \
    -timeout "${WHALE_EDGE_TIMEOUT:-20s}" || true
  sleep "$INTERVAL"
done
