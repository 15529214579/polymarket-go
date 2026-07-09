#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

DEFAULT_WALLETS="$ROOT/wallets.strategy-push.txt"
if [ ! -s "$DEFAULT_WALLETS" ] && [ -s "$ROOT/wallets.strategy-core.txt" ]; then
  DEFAULT_WALLETS="$ROOT/wallets.strategy-core.txt"
fi

REPORT_PATH="${WHALE_PERF_REPORT:-$ROOT/reports/whale_performance.md}"
SUMMARY_JSON="${WHALE_PERF_SUMMARY_JSON:-${REPORT_PATH%.md}.json}"
SNAPSHOT_JSONL="${WHALE_PERF_SNAPSHOT_JSONL:-$ROOT/db/strategy_iteration/whale_performance_snapshots.jsonl}"
WHALE_LIST_MIN_USD_EFFECTIVE="${WHALE_LIST_MIN_USD:-core=1000,sports=1500,watch=750,scout=750,target=750,tape=1500,flow=1500,leaderboard_push=5000,leaderboard_watch=5000}"
POLICY_START_FILE="${WHALE_POLICY_START_FILE:-$ROOT/db/whale-push.policy-start}"
WHALE_PERF_SINCE_EFFECTIVE="${WHALE_PERF_SINCE:-}"
if [ -z "$WHALE_PERF_SINCE_EFFECTIVE" ] && [ -s "$POLICY_START_FILE" ]; then
  WHALE_PERF_SINCE_EFFECTIVE="$(awk -F= '$1 == "started_at" {print $2; exit}' "$POLICY_START_FILE")"
fi

go build -o bin/whale-report ./cmd/whale-report

args=(
  -log "${WHALE_LOG:-$ROOT/db/journal/whale_trades.jsonl}" \
  -wallets "${WHALE_PERF_WALLETS:-$DEFAULT_WALLETS}" \
  -report "$REPORT_PATH" \
  -summary_json "$SUMMARY_JSON" \
  -snapshot_jsonl "$SNAPSHOT_JSONL" \
  -stake "${WHALE_PERF_STAKE:-10}" \
  -min_notional "${WHALE_MIN_USD:-750}" \
  -list_min_notional "$WHALE_LIST_MIN_USD_EFFECTIVE" \
  -actions "${WHALE_PERF_ACTIONS:-alert,followed}" \
  -repeat_cooldown "${WHALE_REPEAT_COOLDOWN:-3m}" \
  -repeat_min_notional "${WHALE_REPEAT_MIN_USD:-5000}" \
  -mark_current="${WHALE_PERF_MARK_CURRENT:-true}" \
  -max_open_age "${WHALE_PERF_MAX_OPEN_AGE:-0s}" \
  -timeout "${WHALE_PERF_TIMEOUT:-20s}"
)
if [ -n "$WHALE_PERF_SINCE_EFFECTIVE" ]; then
  args+=(-since "$WHALE_PERF_SINCE_EFFECTIVE")
fi

exec "$ROOT/bin/whale-report" "${args[@]}"
