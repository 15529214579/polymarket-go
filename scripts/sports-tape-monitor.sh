#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

LOG="${SPORTS_TAPE_LOG:-$ROOT/logs/sports-tape-monitor.log}"
INTERVAL="${SPORTS_TAPE_INTERVAL:-60}"
LOCKDIR="${SPORTS_TAPE_LOCKDIR:-$ROOT/db/sports-tape-monitor.lock}"
TAPE_PAGES="${SPORTS_TAPE_PAGES:-8}"
TAPE_TIMEOUT="${SPORTS_TAPE_TIMEOUT:-3m}"

mkdir -p "$ROOT/logs" "$ROOT/reports" "$ROOT/db/strategy_iteration"
if ! mkdir "$LOCKDIR" 2>/dev/null; then
  existing_pid="$(cat "$LOCKDIR/pid" 2>/dev/null || true)"
  if [ -n "$existing_pid" ] && kill -0 "$existing_pid" 2>/dev/null; then
    printf 'sports-tape-monitor.already_running pid=%s lock=%s\n' "$existing_pid" "$LOCKDIR" >> "$LOG"
    exit 0
  fi
  rm -f "$LOCKDIR/pid"
  rmdir "$LOCKDIR" 2>/dev/null || true
  mkdir "$LOCKDIR"
fi
printf '%s\n' "$$" > "$LOCKDIR/pid"
cleanup() {
  rm -f "$LOCKDIR/pid"
  rmdir "$LOCKDIR" 2>/dev/null || true
}
trap cleanup EXIT
trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM
trap 'printf "sports-tape-monitor.hup_exit pid=%s\n" "$$" >> "$LOG"; cleanup; exit 129' HUP

go build -o bin/sports-tape ./cmd/sports-tape
go build -o bin/sports-tape-alert ./cmd/sports-tape-alert
go build -o bin/sports-alert-report ./cmd/sports-alert-report

ALERT_MODES="${SPORTS_TAPE_ALERT_MODES:-FOLLOW-READY,CANDIDATE,PROBATION,EDGE-HOT,FLOW-SCOUT}"
case ",$ALERT_MODES," in
  *,FLOW-SCOUT,*) ;;
  *) ALERT_MODES="$ALERT_MODES,FLOW-SCOUT" ;;
esac
OBSERVE_MIN_NOTIONAL="${SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL:-3000}"
if [ "${SPORTS_TAPE_ALERT_ENABLE_OBSERVE:-1}" = "0" ]; then
  OBSERVE_MIN_NOTIONAL="0"
fi
OBSERVE_REQUIRE_KNOWN="${SPORTS_TAPE_ALERT_OBSERVE_REQUIRE_KNOWN:-true}"
OBSERVE_MIN_TIER="${SPORTS_TAPE_ALERT_OBSERVE_MIN_TIER:-B}"
TAPE_WALLET_STATUSES="${SPORTS_TAPE_WALLET_STATUSES:-$ROOT/wallets.football-score-push.txt,$ROOT/wallets.leaderboard-watch.txt,$ROOT/wallets.leaderboard-push.txt,$ROOT/wallets.leaderboard-sports-push.txt,$ROOT/wallets.sports-holders-push.txt,$ROOT/wallets.strategy-push.txt,$ROOT/wallets.strategy-tape.txt,$ROOT/wallets.strategy-tape-observe.txt,$ROOT/wallets.strategy-tape-probation.txt,$ROOT/wallets.strategy-tape-candidates.txt,$ROOT/wallets.strategy-tape-edgehot.txt,$ROOT/wallets.strategy-tape-follow.txt,$ROOT/wallets.strategy-consensus.txt,$ROOT/wallets.strategy-tape-reversal.txt,$ROOT/wallets.strategy-review-noise.txt}"
TAPE_EXCLUDE_WALLETS="${SPORTS_TAPE_EXCLUDE_WALLETS:-$ROOT/wallets.strategy-quarantine.txt,$ROOT/wallets.strategy-review-noise.txt}"
TAPE_ALERT_WALLET_STATUSES="${SPORTS_TAPE_ALERT_WALLET_STATUSES:-$ROOT/wallets.football-score-push.txt,$ROOT/wallets.leaderboard-watch.txt,$ROOT/wallets.leaderboard-push.txt,$ROOT/wallets.leaderboard-sports-push.txt,$ROOT/wallets.sports-holders-push.txt,$ROOT/wallets.strategy-push.txt,$ROOT/wallets.strategy-tape-follow.txt,$ROOT/wallets.strategy-tape-candidates.txt,$ROOT/wallets.strategy-tape-edgehot.txt,$ROOT/wallets.strategy-tape-probation.txt,$ROOT/wallets.strategy-tape-observe.txt,$ROOT/wallets.strategy-consensus.txt,$ROOT/wallets.strategy-tape-reversal.txt,$ROOT/wallets.strategy-review-noise.txt}"
POLICY_REFRESH_INTERVAL_SECONDS="${SPORTS_TAPE_POLICY_REFRESH_INTERVAL_SECONDS:-300}"
last_policy_refresh=0

exec >> "$LOG" 2>&1
echo "sports-tape-monitor.start pid=$$ interval=${INTERVAL}s"
echo "sports-tape-monitor.pages ${TAPE_PAGES}"
echo "sports-tape-monitor.timeout ${TAPE_TIMEOUT}"
echo "sports-tape-monitor.alert_modes ${ALERT_MODES}"
echo "sports-tape-monitor.observe_min_notional ${OBSERVE_MIN_NOTIONAL}"
echo "sports-tape-monitor.observe_burst_min_notional ${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}"
echo "sports-tape-monitor.shadow_log ${SPORTS_TAPE_SHADOW_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_shadow_alerts.jsonl}"
echo "sports-tape-monitor.observe_require_known ${OBSERVE_REQUIRE_KNOWN}"
echo "sports-tape-monitor.observe_min_tier ${OBSERVE_MIN_TIER}"
echo "sports-tape-monitor.insider_min_notional ${SPORTS_TAPE_ALERT_INSIDER_MIN_NOTIONAL:-25000}"
echo "sports-tape-monitor.insider_max_bot ${SPORTS_TAPE_ALERT_INSIDER_MAX_BOT:-35}"
echo "sports-tape-monitor.consensus_max_age ${SPORTS_TAPE_CONSENSUS_MAX_AGE:-15m}"
echo "sports-tape-monitor.unknown_flow_min_notional ${SPORTS_TAPE_UNKNOWN_FLOW_MIN_NOTIONAL:-4000}"
echo "sports-tape-monitor.unknown_flow_min_markets ${SPORTS_TAPE_UNKNOWN_FLOW_MIN_MARKETS:-2}"
echo "sports-tape-monitor.seed_flow_min_notional ${SPORTS_TAPE_SEED_FLOW_MIN_NOTIONAL:-3000}"
echo "sports-tape-monitor.seed_flow_min_markets ${SPORTS_TAPE_SEED_FLOW_MIN_MARKETS:-2}"
echo "sports-tape-monitor.scored_flow_min_notional ${SPORTS_TAPE_SCORED_FLOW_MIN_NOTIONAL:-4000}"
echo "sports-tape-monitor.scored_flow_min_markets ${SPORTS_TAPE_SCORED_FLOW_MIN_MARKETS:-2}"
echo "sports-tape-monitor.scored_flow_min_tier ${SPORTS_TAPE_SCORED_FLOW_MIN_TIER:-B}"
echo "sports-tape-monitor.scored_flow_max_bot ${SPORTS_TAPE_SCORED_FLOW_MAX_BOT:-35}"
echo "sports-tape-monitor.mode_policy ${SPORTS_TAPE_MODE_POLICY:-$ROOT/db/strategy_iteration/sports_policy_decision.json}"
echo "sports-tape-monitor.mode_policy_max_age ${SPORTS_TAPE_MODE_POLICY_MAX_AGE:-2h}"
echo "sports-tape-monitor.mode_policy_min_action ${SPORTS_TAPE_MODE_POLICY_MIN_ACTION:-COLLECT_POSITIVE}"
echo "sports-tape-monitor.policy_refresh_interval_seconds ${POLICY_REFRESH_INTERVAL_SECONDS}"
echo "sports-tape-monitor.policy_report ${SPORTS_TAPE_POLICY_REPORT:-$ROOT/reports/sports_alert_shadow_performance.md}"
echo "sports-tape-monitor.wallet_statuses ${TAPE_WALLET_STATUSES}"
echo "sports-tape-monitor.exclude_wallets ${TAPE_EXCLUDE_WALLETS}"
echo "sports-tape-monitor.alert_wallet_statuses ${TAPE_ALERT_WALLET_STATUSES}"

while :; do
  date '+===== %Y-%m-%d %H:%M:%S %Z ====='
  "$ROOT/bin/sports-tape" \
    -output "${SPORTS_TAPE_OUTPUT:-$ROOT/db/strategy_iteration/sports_tape.jsonl}" \
    -retain_input "${SPORTS_TAPE_RETAIN_INPUT:-$ROOT/db/strategy_iteration/sports_tape.jsonl}" \
    -report "${SPORTS_TAPE_REPORT:-$ROOT/reports/sports_tape.md}" \
    -wallets_out "${SPORTS_TAPE_WALLETS_OUT:-$ROOT/wallets.sports-tape.txt}" \
    -push_wallets "${SPORTS_TAPE_PUSH_WALLETS:-$ROOT/wallets.strategy-push.txt}" \
    -wallet_statuses "$TAPE_WALLET_STATUSES" \
    -exclude_wallets "$TAPE_EXCLUDE_WALLETS" \
    -scores "${SPORTS_TAPE_SCORES:-$ROOT/db/strategy_iteration/wallet_scores.json}" \
    -target_categories "${SPORTS_TAPE_CATEGORIES:-basketball,soccer,esports}" \
    -pages "$TAPE_PAGES" \
    -limit "${SPORTS_TAPE_LIMIT:-500}" \
    -min_notional "${SPORTS_TAPE_MIN_NOTIONAL:-500}" \
    -retain_window "${SPORTS_TAPE_RETAIN_WINDOW:-6h}" \
    -top "${SPORTS_TAPE_TOP:-25}" \
    -timeout "$TAPE_TIMEOUT" || true
  now_epoch="$(date +%s)"
  if [ "$POLICY_REFRESH_INTERVAL_SECONDS" -le 0 ] || [ $((now_epoch - last_policy_refresh)) -ge "$POLICY_REFRESH_INTERVAL_SECONDS" ]; then
    "$ROOT/bin/sports-alert-report" \
      -log "${SPORTS_TAPE_SHADOW_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_shadow_alerts.jsonl}" \
      -extra_log "${SPORTS_TAPE_POLICY_EXTRA_LOG:-${SPORTS_TAPE_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_alerts.jsonl}}" \
      -report "${SPORTS_TAPE_POLICY_REPORT:-$ROOT/reports/sports_alert_shadow_performance.md}" \
      -decision_json "${SPORTS_TAPE_MODE_POLICY:-$ROOT/db/strategy_iteration/sports_policy_decision.json}" \
      -current_policy_modes "${SPORTS_TAPE_POLICY_CURRENT_MODES:-CONSENSUS,OBSERVE-BURST}" \
      -current_exclude_wallets "${SPORTS_TAPE_POLICY_EXCLUDE_WALLETS:-$ROOT/wallets.strategy-quarantine.txt,$ROOT/wallets.strategy-review-noise.txt,$ROOT/wallets.strategy-tape-reversal.txt}" \
      -timeout "${SPORTS_TAPE_POLICY_TIMEOUT:-20s}" && last_policy_refresh="$now_epoch" || true
  fi
  "$ROOT/bin/sports-tape-alert" \
    -tape "${SPORTS_TAPE_OUTPUT:-$ROOT/db/strategy_iteration/sports_tape.jsonl}" \
    -state "${SPORTS_TAPE_ALERT_STATE:-$ROOT/db/strategy_iteration/sports_tape_alert_sent.json}" \
    -sent_log "${SPORTS_TAPE_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_alerts.jsonl}" \
    -shadow_log "${SPORTS_TAPE_SHADOW_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_shadow_alerts.jsonl}" \
    -diagnostic_report "${SPORTS_TAPE_ALERT_CANDIDATE_REPORT:-$ROOT/reports/sports_alert_candidates.md}" \
    -edge_snapshots "${SPORTS_TAPE_ALERT_EDGE_SNAPSHOTS:-$ROOT/db/strategy_iteration/whale_edge_snapshots.jsonl}" \
    -mode_policy "${SPORTS_TAPE_MODE_POLICY:-$ROOT/db/strategy_iteration/sports_policy_decision.json}" \
    -mode_policy_max_age "${SPORTS_TAPE_MODE_POLICY_MAX_AGE:-2h}" \
    -mode_policy_min_action "${SPORTS_TAPE_MODE_POLICY_MIN_ACTION:-COLLECT_POSITIVE}" \
    -wallet_statuses "$TAPE_ALERT_WALLET_STATUSES" \
    -env "${SPORTS_TAPE_ALERT_ENV:-$ROOT/.env.local}" \
    -min_notional "${SPORTS_TAPE_ALERT_MIN_NOTIONAL:-3000}" \
    -observe_min_notional "$OBSERVE_MIN_NOTIONAL" \
    -observe_burst_min_notional "${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}" \
    -observe_max_bot "${SPORTS_TAPE_ALERT_OBSERVE_MAX_BOT:-35}" \
    -observe_require_known="$OBSERVE_REQUIRE_KNOWN" \
    -observe_min_tier "$OBSERVE_MIN_TIER" \
    -insider_min_notional "${SPORTS_TAPE_ALERT_INSIDER_MIN_NOTIONAL:-25000}" \
    -insider_max_bot "${SPORTS_TAPE_ALERT_INSIDER_MAX_BOT:-35}" \
    -edge_hot_min_notional "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_NOTIONAL:-750}" \
    -edge_hot_max_bot "${SPORTS_TAPE_ALERT_EDGE_HOT_MAX_BOT:-45}" \
    -edge_hot_min_samples "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_SAMPLES:-2}" \
    -edge_hot_min_avg_pp "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_AVG_PP:-2}" \
    -edge_hot_min_win_rate "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_WIN_RATE:-60}" \
    -edge_hot_min_5m_avg_pp "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_5M_AVG_PP:-0.5}" \
    -edge_hot_min_15m_avg_pp "${SPORTS_TAPE_ALERT_EDGE_HOT_MIN_15M_AVG_PP:-0}" \
    -edge_hot_max_1h_neg_pp "${SPORTS_TAPE_ALERT_EDGE_HOT_MAX_1H_NEG_PP:--5}" \
    -edge_block_15m_samples "${SPORTS_TAPE_ALERT_EDGE_BLOCK_15M_SAMPLES:-2}" \
    -edge_block_15m_max_avg_pp "${SPORTS_TAPE_ALERT_EDGE_BLOCK_15M_MAX_AVG_PP:--1}" \
    -edge_block_1h_samples "${SPORTS_TAPE_ALERT_EDGE_BLOCK_1H_SAMPLES:-1}" \
    -edge_block_1h_max_avg_pp "${SPORTS_TAPE_ALERT_EDGE_BLOCK_1H_MAX_AVG_PP:--5}" \
    -require_positive_edge_modes "${SPORTS_TAPE_ALERT_REQUIRE_POSITIVE_EDGE_MODES:-CANDIDATE,PROBATION}" \
    -modes "$ALERT_MODES" \
    -max_age "${SPORTS_TAPE_ALERT_MAX_AGE:-10m}" \
    -diagnostic_age "${SPORTS_TAPE_ALERT_DIAGNOSTIC_AGE:-6h}" \
    -consensus_alerts="${SPORTS_TAPE_CONSENSUS_ALERTS:-true}" \
    -consensus_min_notional "${SPORTS_TAPE_CONSENSUS_MIN_NOTIONAL:-7500}" \
    -consensus_min_wallets "${SPORTS_TAPE_CONSENSUS_MIN_WALLETS:-2}" \
    -consensus_max_bot "${SPORTS_TAPE_CONSENSUS_MAX_BOT:-60}" \
    -consensus_max_age "${SPORTS_TAPE_CONSENSUS_MAX_AGE:-15m}" \
    -unknown_flow_min_notional "${SPORTS_TAPE_UNKNOWN_FLOW_MIN_NOTIONAL:-4000}" \
    -unknown_flow_min_markets "${SPORTS_TAPE_UNKNOWN_FLOW_MIN_MARKETS:-2}" \
    -unknown_flow_max_bot "${SPORTS_TAPE_UNKNOWN_FLOW_MAX_BOT:-45}" \
    -seed_flow_min_notional "${SPORTS_TAPE_SEED_FLOW_MIN_NOTIONAL:-3000}" \
    -seed_flow_min_markets "${SPORTS_TAPE_SEED_FLOW_MIN_MARKETS:-2}" \
    -scored_flow_min_notional "${SPORTS_TAPE_SCORED_FLOW_MIN_NOTIONAL:-4000}" \
    -scored_flow_min_markets "${SPORTS_TAPE_SCORED_FLOW_MIN_MARKETS:-2}" \
    -scored_flow_max_bot "${SPORTS_TAPE_SCORED_FLOW_MAX_BOT:-35}" \
    -scored_flow_min_tier "${SPORTS_TAPE_SCORED_FLOW_MIN_TIER:-B}" \
    -position_cooldown "${SPORTS_TAPE_ALERT_POSITION_COOLDOWN:-30m}" \
    -repeat_min_notional "${SPORTS_TAPE_ALERT_REPEAT_MIN_NOTIONAL:-15000}" \
    -max_alerts "${SPORTS_TAPE_ALERT_MAX_ALERTS:-5}" || true
  sleep "$INTERVAL" || true
done
