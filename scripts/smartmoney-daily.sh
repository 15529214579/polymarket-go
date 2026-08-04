#!/usr/bin/env bash
# Daily smart-money pipeline: crawl leaderboard, rebuild the core wallet list,
# evaluate recent whale-push signals, and write an operator summary.
set -euo pipefail

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p "$ROOT/db" "$ROOT/logs" "$ROOT/reports"

LOCK_DIR="$ROOT/db/smartmoney-daily.lock"
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
  printf 'smartmoney daily already running lock=%s\n' "$LOCK_DIR"
  exit 0
fi
trap 'rmdir "$LOCK_DIR" 2>/dev/null || true' EXIT

DAY="$(date '+%Y-%m-%d')"
RUN_AT="$(date '+%Y-%m-%d %H:%M:%S %Z')"
LOG="${SMARTMONEY_DAILY_LOG:-$ROOT/logs/smartmoney-daily-$DAY.log}"
SUMMARY="${SMARTMONEY_DAILY_REPORT:-$ROOT/reports/smartmoney_daily.md}"
PAPER_PNL_REPORT="${SMARTMONEY_PAPER_PNL_REPORT:-$ROOT/reports/smartmoney-paper-pnl.md}"
PAPER_WALLET_POLICY_REPORT="${SMARTMONEY_PAPER_WALLET_POLICY_REPORT:-$ROOT/reports/smartmoney-paper-wallets.md}"
PAPER_SHADOW_REPORT="${SMARTMONEY_SHADOW_REPORT:-$ROOT/reports/smartmoney-exit-shadow.md}"
PAPER_PROMOTED_WALLETS="${SMARTMONEY_PAPER_PROMOTED_WALLETS:-$ROOT/db/smartmoney-paper/wallets.paper-promoted.txt}"
PAPER_DEMOTED_WALLETS="${SMARTMONEY_PAPER_DEMOTED_WALLETS:-$ROOT/db/smartmoney-paper/wallets.paper-demoted.txt}"
CORE="${SMARTMONEY_CORE_WALLETS:-$ROOT/wallets.strategy-core.txt}"
WATCH="${SMARTMONEY_WATCH_WALLETS:-$ROOT/wallets.strategy-watch.txt}"
SPORTS="${SMARTMONEY_SPORTS_WALLETS:-$ROOT/wallets.strategy-sports.txt}"
SCOUT="${SMARTMONEY_SCOUT_WALLETS:-$ROOT/wallets.strategy-scout.txt}"
TARGET="${SMARTMONEY_TARGET_WALLETS:-$ROOT/wallets.strategy-target.txt}"
FLOW="${SMARTMONEY_FLOW_WALLETS:-$ROOT/wallets.strategy-flow.txt}"
TAPE="${SMARTMONEY_TAPE_WALLETS:-$ROOT/wallets.strategy-tape.txt}"
TAPE_OBSERVE="${SMARTMONEY_TAPE_OBSERVE_WALLETS:-$ROOT/wallets.strategy-tape-observe.txt}"
TAPE_PROBATION="${SMARTMONEY_TAPE_PROBATION_WALLETS:-$ROOT/wallets.strategy-tape-probation.txt}"
TAPE_CANDIDATES="${SMARTMONEY_TAPE_CANDIDATE_WALLETS:-$ROOT/wallets.strategy-tape-candidates.txt}"
TAPE_FOLLOW="${SMARTMONEY_TAPE_FOLLOW_WALLETS:-$ROOT/wallets.strategy-tape-follow.txt}"
TAPE_REVERSAL="${SMARTMONEY_TAPE_REVERSAL_WALLETS:-$ROOT/wallets.strategy-tape-reversal.txt}"
TAPE_EDGEHOT="${SMARTMONEY_TAPE_EDGEHOT_WALLETS:-$ROOT/wallets.strategy-tape-edgehot.txt}"
CONSENSUS_RESEARCH="${SMARTMONEY_CONSENSUS_RESEARCH_WALLETS:-$ROOT/wallets.strategy-consensus.txt}"
PUSH="${SMARTMONEY_PUSH_WALLETS:-$ROOT/wallets.strategy-push.txt}"
DISCOVERY_REPORT="${SMARTMONEY_DISCOVERY_REPORT:-$ROOT/reports/strategy_iteration.md}"
LAB_REPORT="${SMARTMONEY_LAB_REPORT:-$ROOT/reports/strategy_lab.md}"
LEADERBOARD_WHALES_REPORT="${SMARTMONEY_LEADERBOARD_WHALES_REPORT:-$ROOT/reports/leaderboard_whales.md}"
LEADERBOARD_WATCH="${SMARTMONEY_LEADERBOARD_WATCH_WALLETS:-$ROOT/wallets.leaderboard-watch.txt}"
LEADERBOARD_PUSH="${SMARTMONEY_LEADERBOARD_PUSH_WALLETS:-$ROOT/wallets.leaderboard-push.txt}"
LEADERBOARD_SPORTS_PUSH="${SMARTMONEY_LEADERBOARD_SPORTS_PUSH_WALLETS:-$ROOT/wallets.leaderboard-sports-push.txt}"
SPORTS_HOLDERS_PUSH="${SMARTMONEY_SPORTS_HOLDERS_PUSH_WALLETS:-$ROOT/wallets.sports-holders-push.txt}"
FOOTBALL_SCORE_PUSH="${SMARTMONEY_FOOTBALL_SCORE_PUSH_WALLETS:-$ROOT/wallets.football-score-push.txt}"
SPORTS_TAPE_REPORT="${SMARTMONEY_SPORTS_TAPE_REPORT:-$ROOT/reports/sports_tape.md}"
SPORTS_ALERT_PERF_REPORT="${SPORTS_ALERT_PERF_REPORT:-$ROOT/reports/sports_alert_performance.md}"
SPORTS_ALERT_CANDIDATE_REPORT="${SPORTS_ALERT_CANDIDATE_REPORT:-$ROOT/reports/sports_alert_candidates.md}"
SPORTS_BURST_PERF_REPORT="${SPORTS_BURST_PERF_REPORT:-$ROOT/reports/sports_burst_performance.md}"
WHALE_EDGE_REPORT="${WHALE_EDGE_REPORT:-$ROOT/reports/whale_edge.md}"
WHALE_EDGE_SNAPSHOTS="${WHALE_EDGE_SNAPSHOTS:-$ROOT/db/strategy_iteration/whale_edge_snapshots.jsonl}"
SPORTS_TAPE_ALERT_STATE="${SPORTS_TAPE_ALERT_STATE:-$ROOT/db/strategy_iteration/sports_tape_alert_sent.json}"
SPORTS_TAPE_ALERT_LOG="${SPORTS_TAPE_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_alerts.jsonl}"
SPORTS_TAPE_SHADOW_ALERT_LOG="${SPORTS_TAPE_SHADOW_ALERT_LOG:-$ROOT/db/strategy_iteration/sports_tape_shadow_alerts.jsonl}"
SPORTS_TAPE_MODE_POLICY="${SPORTS_TAPE_MODE_POLICY:-$ROOT/db/strategy_iteration/sports_policy_decision.json}"
SPORTS_ALERT_MARK_CACHE="${SPORTS_ALERT_MARK_CACHE:-$ROOT/db/strategy_iteration/sports_alert_midpoints.json}"
SPORTS_ALERT_SHADOW_PERF_REPORT="${SPORTS_ALERT_SHADOW_PERF_REPORT:-$ROOT/reports/sports_alert_shadow_performance.md}"
SPORTS_CONSENSUS_EVENTS="${SPORTS_CONSENSUS_EVENTS:-$ROOT/db/strategy_iteration/sports_consensus_events.jsonl}"
SPORTS_CONSENSUS_WATCH_EVENTS="${SPORTS_CONSENSUS_WATCH_EVENTS:-$ROOT/db/strategy_iteration/sports_consensus_watch_events.jsonl}"
SPORTS_TAPE_ALERT_MODES_EFFECTIVE="${SPORTS_TAPE_ALERT_MODES:-FOLLOW-READY,CANDIDATE,PROBATION,EDGE-HOT,FLOW-SCOUT}"
case ",$SPORTS_TAPE_ALERT_MODES_EFFECTIVE," in
  *,FLOW-SCOUT,*) ;;
  *) SPORTS_TAPE_ALERT_MODES_EFFECTIVE="$SPORTS_TAPE_ALERT_MODES_EFFECTIVE,FLOW-SCOUT" ;;
esac
SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL_EFFECTIVE="${SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL:-3000}"
if [ "${SPORTS_TAPE_ALERT_ENABLE_OBSERVE:-1}" = "0" ]; then
  SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL_EFFECTIVE="0"
fi
SPORTS_TAPE_ALERT_OBSERVE_REQUIRE_KNOWN_EFFECTIVE="${SPORTS_TAPE_ALERT_OBSERVE_REQUIRE_KNOWN:-true}"
SPORTS_TAPE_ALERT_OBSERVE_MIN_TIER_EFFECTIVE="${SPORTS_TAPE_ALERT_OBSERVE_MIN_TIER:-B}"
SPORTS_ALERT_PROMOTE_MIN_MARKED="5"
SPORTS_ALERT_PROMOTE_MIN_ROI="5.0"
SPORTS_ALERT_PROMOTE_MIN_WIN="60.0"
CORE_PERF_REPORT="${WHALE_CORE_PERF_REPORT:-$ROOT/reports/whale_performance_core.md}"
SPORTS_PERF_REPORT="${WHALE_SPORTS_PERF_REPORT:-$ROOT/reports/whale_performance_sports.md}"
SCOUT_PERF_REPORT="${WHALE_SCOUT_PERF_REPORT:-$ROOT/reports/whale_performance_scout.md}"
TARGET_PERF_REPORT="${WHALE_TARGET_PERF_REPORT:-$ROOT/reports/whale_performance_target.md}"
FLOW_PERF_REPORT="${WHALE_FLOW_PERF_REPORT:-$ROOT/reports/whale_performance_flow.md}"
TAPE_PERF_REPORT="${WHALE_TAPE_PERF_REPORT:-$ROOT/reports/whale_performance_tape.md}"
PUSH_PERF_REPORT="${WHALE_PERF_REPORT:-$ROOT/reports/whale_performance.md}"
PUSH_PERF_JSON="${WHALE_PERF_SUMMARY_JSON:-${PUSH_PERF_REPORT%.md}.json}"
POLICY_ALERT_SENT="${SMARTMONEY_POLICY_ALERT_SENT:-$ROOT/logs/smartmoney-policy-violation-sent.json}"
MAINT_REPORT="${SMARTMONEY_MAINT_REPORT:-$ROOT/reports/wallet_maintenance.md}"
QUARANTINE="${SMARTMONEY_QUARANTINE_WALLETS:-$ROOT/wallets.strategy-quarantine.txt}"
REVIEW_NOISE="${SMARTMONEY_REVIEW_NOISE_WALLETS:-$ROOT/wallets.strategy-review-noise.txt}"
STRATEGY_EXCLUDE_FILE="${STRATEGY_EXCLUDE_WALLETS:-$ROOT/db/strategy_iteration/wallets.strategy-exclude.txt}"
MIN_SIGNALS="${WHALE_HEALTH_MIN_SIGNALS:-10}"
MIN_EVENT_CAPPED_ENTRIES="${WHALE_HEALTH_MIN_EVENT_CAPPED_ENTRIES:-5}"
MIN_PROVEN_SIGNALS="${WHALE_HEALTH_MIN_PROVEN_SIGNALS:-3}"
MIN_PROVEN_EVENT_CAPPED_ENTRIES="${WHALE_HEALTH_MIN_PROVEN_EVENT_CAPPED_ENTRIES:-2}"
MIN_ROI="${WHALE_HEALTH_MIN_ROI:-0}"
EVENT_MIN_SIGNALS="${MAINT_EVENT_MIN_SIGNALS:-1}"
EVENT_QUARANTINE_ROI="${MAINT_EVENT_QUARANTINE_ROI:--30}"
NOISE_MIN_SUPPRESSED="${MAINT_NOISE_MIN_SUPPRESSED:-3}"

hash_file() {
  if [ -s "$1" ]; then
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

wallet_count() {
  if [ -s "$1" ]; then
    grep -Eci '^0x[0-9a-fA-F]{40}([[:space:]]|$)' "$1" || true
  else
    printf '0'
  fi
}

json_sent_count() {
  local file="$1"
  if [ -s "$file" ]; then
    awk '
      /"sent"[[:space:]]*:/ {in_sent=1; next}
      in_sent && /^[[:space:]]*}/ {in_sent=0}
      in_sent && /^[[:space:]]*"[^"]+"[[:space:]]*:/ {n++}
      END {print n+0}
    ' "$file" 2>/dev/null || printf '0'
  else
    printf '0'
  fi
}

file_line_count() {
  if [ -s "$1" ]; then
    awk 'END {print NR+0}' "$1" 2>/dev/null || printf '0'
  else
    printf '0'
  fi
}

metric() {
  local file="$1"
  local label="$2"
  if [ -s "$file" ]; then
    awk -v prefix="- ${label}: " 'index($0, prefix) == 1 {print substr($0, length(prefix) + 1); exit}' "$file"
  fi
}

metric_any() {
  local file="$1"
  shift
  local value=''
  for label in "$@"; do
    value="$(metric "$file" "$label")"
    if [ -n "$value" ]; then
      printf '%s\n' "$value"
      return 0
    fi
  done
}

first_field() {
  awk '{print $1}'
}

nonnegative_int() {
  local value="${1:-0}"
  if [[ "$value" =~ ^[0-9]+$ ]]; then
    printf '%s\n' "$value"
  else
    printf '0\n'
  fi
}

restart_worker() {
  local screen_name="${SMARTMONEY_SCREEN:-polymarket-whale-push}"

  screen -S "$screen_name" -X quit >/dev/null 2>&1 || true
  sleep 1
  screen -dmS "$screen_name" "$ROOT/scripts/start-whale-push.sh"
  printf 'restarted screen %s' "$screen_name"
}

run_push_performance() {
  local push_file="$PUSH"
  if [ ! -s "$push_file" ]; then
    push_file="$CORE"
  fi
  WHALE_PERF_WALLETS="$push_file,$FOOTBALL_SCORE_PUSH,$LEADERBOARD_PUSH,$LEADERBOARD_WATCH,$LEADERBOARD_SPORTS_PUSH,$SPORTS_HOLDERS_PUSH" WHALE_PERF_REPORT="$PUSH_PERF_REPORT" WHALE_PERF_SUMMARY_JSON="$PUSH_PERF_JSON" ./scripts/whale-performance.sh
}

run_wallet_maintenance() {
  local push_file="$PUSH"
  if [ ! -s "$push_file" ]; then
    push_file="$CORE"
  fi
  local maint_wallets="$push_file,$FOOTBALL_SCORE_PUSH,$LEADERBOARD_PUSH,$LEADERBOARD_WATCH,$LEADERBOARD_SPORTS_PUSH,$SPORTS_HOLDERS_PUSH,$TAPE,$TAPE_OBSERVE,$TAPE_PROBATION,$TAPE_CANDIDATES,$TAPE_FOLLOW,$TAPE_REVERSAL"
  go build -o bin/wallet-maintain ./cmd/wallet-maintain
  "$ROOT/bin/wallet-maintain" \
    -log "$ROOT/db/journal/whale_trades.jsonl" \
    -push "$maint_wallets" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -performance_json "$PUSH_PERF_JSON" \
    -edge_snapshots "$WHALE_EDGE_SNAPSHOTS" \
    -report "$MAINT_REPORT" \
    -quarantine "$QUARANTINE" \
    -review_noise "$REVIEW_NOISE" \
    -tape_candidates "$TAPE_CANDIDATES" \
    -tape_follow "$TAPE_FOLLOW" \
    -tape_reversal "$TAPE_REVERSAL" \
    -stake "${WHALE_PERF_STAKE:-10}" \
    -min_signals "${MAINT_MIN_SIGNALS:-5}" \
    -min_roi "${MAINT_MIN_ROI:-0}" \
    -promote_roi "${MAINT_PROMOTE_ROI:-5}" \
    -edge_promote_min_samples "${MAINT_EDGE_PROMOTE_MIN_SAMPLES:-3}" \
    -edge_promote_min_avg_pp "${MAINT_EDGE_PROMOTE_MIN_AVG_PP:-1}" \
    -edge_promote_min_win_rate "${MAINT_EDGE_PROMOTE_MIN_WIN_RATE:-60}" \
    -edge_promote_max_bot "${MAINT_EDGE_PROMOTE_MAX_BOT:-45}" \
    -edge_promote_reversal_min_15m_samples "${MAINT_EDGE_PROMOTE_REVERSAL_MIN_15M_SAMPLES:-2}" \
    -edge_promote_reversal_max_15m_avg_pp "${MAINT_EDGE_PROMOTE_REVERSAL_MAX_15M_AVG_PP:--1}" \
    -edge_promote_severe_min_samples "${MAINT_EDGE_PROMOTE_SEVERE_MIN_SAMPLES:-1}" \
    -edge_promote_severe_max_avg_pp "${MAINT_EDGE_PROMOTE_SEVERE_MAX_AVG_PP:--20}" \
    -edge_promote_negative_min_samples "${MAINT_EDGE_PROMOTE_NEGATIVE_MIN_SAMPLES:-5}" \
    -edge_promote_negative_max_avg_pp "${MAINT_EDGE_PROMOTE_NEGATIVE_MAX_AVG_PP:--0.25}" \
    -edge_promote_negative_max_win_rate "${MAINT_EDGE_PROMOTE_NEGATIVE_MAX_WIN_RATE:-20}" \
    -tape_follow_min_samples "${MAINT_TAPE_FOLLOW_MIN_SAMPLES:-6}" \
    -tape_follow_min_avg_pp "${MAINT_TAPE_FOLLOW_MIN_AVG_PP:-1.5}" \
    -tape_follow_min_win_rate "${MAINT_TAPE_FOLLOW_MIN_WIN_RATE:-65}" \
    -tape_follow_min_5m_avg_pp "${MAINT_TAPE_FOLLOW_MIN_5M_AVG_PP:-0.5}" \
    -tape_follow_min_15m_avg_pp "${MAINT_TAPE_FOLLOW_MIN_15M_AVG_PP:-0}" \
    -tape_follow_max_bot "${MAINT_TAPE_FOLLOW_MAX_BOT:-45}" \
    -event_min_signals "$EVENT_MIN_SIGNALS" \
    -event_quarantine_roi "$EVENT_QUARANTINE_ROI" \
    -noise_min_suppressed "$NOISE_MIN_SUPPRESSED"
}

rebuild_strategy_exclude() {
  if [ -n "${STRATEGY_EXCLUDE_WALLETS:-}" ]; then
    return 0
  fi
  mkdir -p "$ROOT/db/strategy_iteration"
  awk '
    {
      line=$0
      sub(/#.*/, "", line)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
      if (line ~ /^0x[0-9a-fA-F]{40}$/) {
        print tolower(line)
      }
    }
  ' "$QUARANTINE" "$REVIEW_NOISE" 2>/dev/null | sort -u > "$STRATEGY_EXCLUDE_FILE"
}

filter_push_excludes() {
  local before_count='0'
  local after_count='0'
  local tmp=''

  rebuild_strategy_exclude
  if [ ! -s "$PUSH" ]; then
    return 0
  fi

  before_count="$(wallet_count "$PUSH")"
  tmp="$(mktemp "$ROOT/db/strategy_iteration/wallets.strategy-push.filtered.XXXXXX")"
  awk '
    FNR == NR {
      line=$0
      sub(/#.*/, "", line)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
      if (line ~ /^0x[0-9a-fA-F]{40}$/) {
        blocked[tolower(line)] = 1
      }
      next
    }
    {
      line=$0
      body=$0
      sub(/#.*/, "", body)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", body)
      if (body ~ /^0x[0-9a-fA-F]{40}$/ && blocked[tolower(body)]) {
        next
      }
      print line
    }
  ' "$STRATEGY_EXCLUDE_FILE" "$PUSH" > "$tmp"
  mv "$tmp" "$PUSH"

  after_count="$(wallet_count "$PUSH")"
  if [ "$before_count" != "$after_count" ]; then
    printf 'push exclude filter removed=%s before=%s after=%s exclude=%s\n' \
      "$((before_count - after_count))" "$before_count" "$after_count" "$STRATEGY_EXCLUDE_FILE"
  fi
}

section_metric() {
  local file="$1"
  local section="$2"
  local label="$3"
  if [ -s "$file" ]; then
    awk -v section="## ${section}" -v prefix="- ${label}: " '
      $0 == section {inside=1; next}
      /^## / {inside=0}
      inside && index($0, prefix) == 1 {print substr($0, length(prefix) + 1); exit}
    ' "$file"
  fi
}

section_metric_any() {
  local file="$1"
  local section="$2"
  shift 2
  local value=''
  for label in "$@"; do
    value="$(section_metric "$file" "$section" "$label")"
    if [ -n "$value" ]; then
      printf '%s\n' "$value"
      return 0
    fi
  done
}

send_policy_violation_alert() {
  local count="$1"
  local env_file="${SMARTMONEY_ENV_FILE:-$ROOT/.env.local}"
  local tag=''
  local last_tag=''
  local token=''
  local chat=''
  local body=''

  if [ "${SMARTMONEY_POLICY_ALERTS:-1}" = "0" ] || [ "${count:-0}" -le 0 ]; then
    return 0
  fi
  if [ ! -f "$env_file" ]; then
    return 0
  fi

  # shellcheck disable=SC1090
  set -a; . "$env_file"; set +a
  token="${PUSH_BOT_TOKEN:-${TELEGRAM_BOT_TOKEN:-}}"
  chat="${TELEGRAM_CHAT_ID:-}"
  if [ -z "$token" ] || [ -z "$chat" ]; then
    return 0
  fi

  mkdir -p "$ROOT/logs"
  tag="policy_violation:${DAY}:${count}"
  if [ -s "$POLICY_ALERT_SENT" ]; then
    last_tag="$(sed -n 's/.*"tag": *"\([^"]*\)".*/\1/p' "$POLICY_ALERT_SENT" | head -1)"
    if [ "$tag" = "$last_tag" ]; then
      return 0
    fi
  fi

  body="$(printf 'Polymarket smart-money policy violation\n\ncount: %s\nreport: %s\nlog: %s\n\naction: inspect violations before expanding or restarting push wallets' "$count" "$PUSH_PERF_REPORT" "$LOG")"
  curl -s --max-time 10 -X POST \
    "https://api.telegram.org/bot${token}/sendMessage" \
    --data-urlencode "chat_id=${chat}" \
    --data-urlencode "text=${body}" \
    --data-urlencode "disable_notification=false" >/dev/null 2>&1 || true

  printf '{\n  "tag": "%s",\n  "sent_at": "%s",\n  "count": %s\n}\n' \
    "$tag" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$count" > "$POLICY_ALERT_SENT"
}

run_pipeline() {
  printf '===== %s =====\n' "$RUN_AT"
  printf 'root=%s\n' "$ROOT"
  printf 'core_wallets=%s\n' "$CORE"
  printf 'push_wallets=%s\n' "$PUSH"

  run_push_performance
  run_wallet_maintenance
  filter_push_excludes

  go build -o bin/sports-tape ./cmd/sports-tape
  "$ROOT/bin/sports-tape" \
    -output "$ROOT/db/strategy_iteration/sports_tape.jsonl" \
    -report "$SPORTS_TAPE_REPORT" \
    -wallets_out "$ROOT/wallets.sports-tape.txt" \
    -push_wallets "$PUSH" \
    -wallet_statuses "$FOOTBALL_SCORE_PUSH,$LEADERBOARD_WATCH,$LEADERBOARD_PUSH,$LEADERBOARD_SPORTS_PUSH,$SPORTS_HOLDERS_PUSH,$PUSH,$TAPE,$TAPE_OBSERVE,$TAPE_PROBATION,$TAPE_CANDIDATES,$TAPE_EDGEHOT,$TAPE_FOLLOW,$CONSENSUS_RESEARCH,$TAPE_REVERSAL,$REVIEW_NOISE" \
    -exclude_wallets "$QUARANTINE,$REVIEW_NOISE" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -target_categories "${SPORTS_TAPE_TARGET_CATEGORIES:-basketball,soccer,esports}" \
    -pages "${SPORTS_TAPE_PAGES:-20}" \
    -limit "${SPORTS_TAPE_LIMIT:-500}" \
    -min_notional "${SPORTS_TAPE_MIN_NOTIONAL:-500}" \
    -retain_window "${SPORTS_TAPE_RETAIN_WINDOW:-6h}" \
    -top "${SPORTS_TAPE_TOP:-25}"

  go build -o bin/sports-alert-report ./cmd/sports-alert-report
  "$ROOT/bin/sports-alert-report" \
    -log "$SPORTS_TAPE_ALERT_LOG" \
    -report "$SPORTS_ALERT_PERF_REPORT" \
    -mark_cache "$SPORTS_ALERT_MARK_CACHE" \
    -stake "${SPORTS_ALERT_PERF_STAKE:-10}" \
    -current_policy_modes "${SPORTS_ALERT_CURRENT_POLICY_MODES:-FLOW-SCOUT,EDGE-HOT,FOLLOW-READY,CANDIDATE,PROBATION,CONSENSUS}" \
    -current_exclude_wallets "$QUARANTINE,$REVIEW_NOISE,$TAPE_REVERSAL,$TAPE_OBSERVE" \
    -timeout "${SPORTS_ALERT_PERF_TIMEOUT:-20s}"

  go build -o bin/sports-burst-report ./cmd/sports-burst-report
  "$ROOT/bin/sports-burst-report" \
    -tape "$ROOT/db/strategy_iteration/sports_tape.jsonl" \
    -report "$SPORTS_BURST_PERF_REPORT" \
    -wallet_statuses "$PUSH,$TAPE,$TAPE_OBSERVE,$TAPE_PROBATION,$TAPE_CANDIDATES,$TAPE_EDGEHOT,$TAPE_FOLLOW,$CONSENSUS_RESEARCH,$TAPE_REVERSAL,$REVIEW_NOISE" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -participants_out "$CONSENSUS_RESEARCH" \
    -participants_exclude_wallets "$QUARANTINE,$REVIEW_NOISE,$TAPE_REVERSAL,$TAPE_OBSERVE" \
    -consensus_events_out "$SPORTS_CONSENSUS_EVENTS" \
    -consensus_watch_events_out "$SPORTS_CONSENSUS_WATCH_EVENTS" \
    -consensus_alert_logs "$SPORTS_TAPE_ALERT_LOG,$SPORTS_TAPE_SHADOW_ALERT_LOG" \
    -mark_cache "$SPORTS_ALERT_MARK_CACHE" \
    -gamma_base "${SPORTS_ALERT_GAMMA_BASE:-https://gamma-api.polymarket.com}" \
    -stake "${SPORTS_BURST_PERF_STAKE:-10}" \
    -window "${SPORTS_BURST_WINDOW:-15m}" \
    -max_age "${SPORTS_BURST_MAX_AGE:-6h}" \
    -consensus_history_max_age "${SPORTS_CONSENSUS_HISTORY_MAX_AGE:-24h}" \
    -min_notional "${SPORTS_BURST_MIN_NOTIONAL:-5000}" \
    -min_trades "${SPORTS_BURST_MIN_TRADES:-2}" \
    -min_leg_notional "${SPORTS_BURST_MIN_LEG_NOTIONAL:-1000}" \
    -consensus_watch_min_notional "${SPORTS_BURST_CONSENSUS_WATCH_MIN_NOTIONAL:-5000}" \
    -consensus_watch_min_wallets "${SPORTS_BURST_CONSENSUS_WATCH_MIN_WALLETS:-2}" \
    -timeout "${SPORTS_BURST_TIMEOUT:-20s}"

  go build -o bin/sports-tape-alert ./cmd/sports-tape-alert
  "$ROOT/bin/sports-tape-alert" \
    -tape "$ROOT/db/strategy_iteration/sports_tape.jsonl" \
    -state "$SPORTS_TAPE_ALERT_STATE" \
    -sent_log "$SPORTS_TAPE_ALERT_LOG" \
    -shadow_log "$SPORTS_TAPE_SHADOW_ALERT_LOG" \
    -diagnostic_report "$SPORTS_ALERT_CANDIDATE_REPORT" \
    -edge_snapshots "$WHALE_EDGE_SNAPSHOTS" \
    -mode_policy "$SPORTS_TAPE_MODE_POLICY" \
    -mode_policy_max_age "${SPORTS_TAPE_MODE_POLICY_MAX_AGE:-2h}" \
    -mode_policy_min_action "${SPORTS_TAPE_MODE_POLICY_MIN_ACTION:-COLLECT_POSITIVE}" \
    -wallet_statuses "$FOOTBALL_SCORE_PUSH,$LEADERBOARD_WATCH,$LEADERBOARD_PUSH,$LEADERBOARD_SPORTS_PUSH,$SPORTS_HOLDERS_PUSH,$PUSH,$TAPE_FOLLOW,$TAPE_CANDIDATES,$TAPE_EDGEHOT,$TAPE_PROBATION,$TAPE_OBSERVE,$CONSENSUS_RESEARCH,$TAPE_REVERSAL,$REVIEW_NOISE" \
    -min_notional "${SPORTS_TAPE_ALERT_MIN_NOTIONAL:-3000}" \
    -observe_min_notional "$SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL_EFFECTIVE" \
    -observe_burst_min_notional "${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}" \
    -observe_max_bot "${SPORTS_TAPE_ALERT_OBSERVE_MAX_BOT:-35}" \
    -observe_require_known="$SPORTS_TAPE_ALERT_OBSERVE_REQUIRE_KNOWN_EFFECTIVE" \
    -observe_min_tier "$SPORTS_TAPE_ALERT_OBSERVE_MIN_TIER_EFFECTIVE" \
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
    -modes "$SPORTS_TAPE_ALERT_MODES_EFFECTIVE" \
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
    -max_alerts 0 \
    -dry_run

  "$ROOT/bin/sports-alert-report" \
    -log "$SPORTS_TAPE_SHADOW_ALERT_LOG" \
    -extra_log "$SPORTS_TAPE_ALERT_LOG" \
    -report "$SPORTS_ALERT_SHADOW_PERF_REPORT" \
    -decision_json "$SPORTS_TAPE_MODE_POLICY" \
    -mark_cache "$SPORTS_ALERT_MARK_CACHE" \
    -stake "${SPORTS_ALERT_PERF_STAKE:-10}" \
    -current_policy_modes "OBSERVE-BURST,CONSENSUS" \
    -current_exclude_wallets "$QUARANTINE,$REVIEW_NOISE,$TAPE_REVERSAL" \
    -timeout "${SPORTS_ALERT_PERF_TIMEOUT:-20s}"

  if [ "${SMARTMONEY_SKIP_DISCOVER:-0}" = "1" ]; then
    printf 'wallet discovery skipped by SMARTMONEY_SKIP_DISCOVER=1\n'
  else
    ./scripts/wallet-discover.sh
  fi

  go build -o bin/leaderboard-whales ./cmd/leaderboard-whales
  "$ROOT/bin/leaderboard-whales" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -push_wallets "$PUSH" \
    -exclude_wallets "$QUARANTINE,$REVIEW_NOISE" \
    -report "$LEADERBOARD_WHALES_REPORT" \
    -recommend_wallets "$LEADERBOARD_WATCH" \
    -push_wallets_out "$LEADERBOARD_PUSH" \
    -sports_push_wallets_out "$LEADERBOARD_SPORTS_PUSH" \
    -top "${LEADERBOARD_WHALES_TOP:-25}" \
    -min_smart "${LEADERBOARD_WHALES_MIN_SMART:-70}" \
    -max_bot "${LEADERBOARD_WHALES_MAX_BOT:-45}" \
    -min_large "${LEADERBOARD_WHALES_MIN_LARGE:-20}" \
    -min_avg_notional "${LEADERBOARD_WHALES_MIN_AVG_NOTIONAL:-500}" \
    -min_target_trades "${LEADERBOARD_WHALES_MIN_TARGET_TRADES:-1}" \
    -min_target_large "${LEADERBOARD_WHALES_MIN_TARGET_LARGE:-0}" \
    -whale_watch_min_smart "${LEADERBOARD_WHALES_WATCH_MIN_SMART:-80}" \
    -whale_watch_max_bot "${LEADERBOARD_WHALES_WATCH_MAX_BOT:-45}" \
    -whale_watch_min_large "${LEADERBOARD_WHALES_WATCH_MIN_LARGE:-100}" \
    -whale_watch_min_avg_notional "${LEADERBOARD_WHALES_WATCH_MIN_AVG_NOTIONAL:-300}" \
    -whale_watch_min_target_large "${LEADERBOARD_WHALES_WATCH_MIN_TARGET_LARGE:-20}" \
    -recommend_limit "${LEADERBOARD_WHALES_RECOMMEND_LIMIT:-50}" \
    -push_limit "${LEADERBOARD_WHALES_PUSH_LIMIT:-25}" \
    -push_min_tier "${LEADERBOARD_WHALES_PUSH_MIN_TIER:-B}" \
    -push_min_smart "${LEADERBOARD_WHALES_PUSH_MIN_SMART:-80}" \
    -push_max_bot "${LEADERBOARD_WHALES_PUSH_MAX_BOT:-35}" \
    -push_min_large "${LEADERBOARD_WHALES_PUSH_MIN_LARGE:-50}" \
    -push_min_avg_notional "${LEADERBOARD_WHALES_PUSH_MIN_AVG_NOTIONAL:-1000}" \
    -push_min_target_trades "${LEADERBOARD_WHALES_PUSH_MIN_TARGET_TRADES:-5}" \
    -push_min_target_large "${LEADERBOARD_WHALES_PUSH_MIN_TARGET_LARGE:-1}" \
    -sports_push_limit "${LEADERBOARD_SPORTS_PUSH_LIMIT:-200}" \
    -sports_push_min_tier "${LEADERBOARD_SPORTS_PUSH_MIN_TIER:-C}" \
    -sports_push_min_smart "${LEADERBOARD_SPORTS_PUSH_MIN_SMART:-55}" \
    -sports_push_max_bot "${LEADERBOARD_SPORTS_PUSH_MAX_BOT:-45}" \
    -sports_push_min_large "${LEADERBOARD_SPORTS_PUSH_MIN_LARGE:-5}" \
    -sports_push_min_avg_notional "${LEADERBOARD_SPORTS_PUSH_MIN_AVG_NOTIONAL:-300}" \
    -sports_push_min_target_trades "${LEADERBOARD_SPORTS_PUSH_MIN_TARGET_TRADES:-3}" \
    -sports_push_min_target_large "${LEADERBOARD_SPORTS_PUSH_MIN_TARGET_LARGE:-1}"

  go build -o bin/sports-holders-push ./cmd/sports-holders-push
  "$ROOT/bin/sports-holders-push" \
    -out "$SPORTS_HOLDERS_PUSH" \
    -exclude_wallets "$QUARANTINE,$REVIEW_NOISE" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -target_categories "${SPORTS_HOLDERS_PUSH_TARGET_CATEGORIES:-basketball,soccer,esports}" \
    -markets "${SPORTS_HOLDERS_PUSH_MARKETS:-300}" \
    -holders "${SPORTS_HOLDERS_PUSH_HOLDERS:-100}" \
    -max_wallets "${SPORTS_HOLDERS_PUSH_MAX_WALLETS:-200}" \
    -min_shares "${SPORTS_HOLDERS_PUSH_MIN_SHARES:-1000}" \
    -timeout "${SPORTS_HOLDERS_PUSH_TIMEOUT:-5m}"

  score_tmp="$(mktemp "$ROOT/db/football-score-push.daily.XXXXXX")"
  if "$ROOT/bin/sports-holders-push" \
    -out "$score_tmp" \
    -list football_score_push \
    -football_score_only \
    -exclude_wallets "$QUARANTINE,$REVIEW_NOISE" \
    -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
    -target_categories soccer \
    -markets "${FOOTBALL_SCORE_PUSH_MARKETS:-300}" \
    -holders "${FOOTBALL_SCORE_PUSH_HOLDERS:-250}" \
    -max_wallets "${FOOTBALL_SCORE_PUSH_MAX_WALLETS:-500}" \
    -min_shares "${FOOTBALL_SCORE_PUSH_MIN_SHARES:-50}" \
    -timeout "${FOOTBALL_SCORE_PUSH_TIMEOUT:-10m}" && \
    grep -Eq '^0x[0-9a-fA-F]{40}([[:space:]]|$)' "$score_tmp"; then
    mv "$score_tmp" "$FOOTBALL_SCORE_PUSH"
  else
    rm -f "$score_tmp"
    printf 'football score holder scan empty or failed; preserved %s\n' "$FOOTBALL_SCORE_PUSH"
  fi

  WHALE_PERF_WALLETS="$CORE" WHALE_PERF_REPORT="$CORE_PERF_REPORT" ./scripts/whale-performance.sh

  local sports_file="$SPORTS"
  if [ ! -s "$sports_file" ]; then
    sports_file="$ROOT/db/empty_wallets.txt"
    : > "$sports_file"
  fi
  WHALE_PERF_WALLETS="$sports_file" WHALE_PERF_REPORT="$SPORTS_PERF_REPORT" ./scripts/whale-performance.sh

  local scout_file="$SCOUT"
  if [ ! -s "$scout_file" ]; then
    scout_file="$ROOT/db/empty_wallets.txt"
    : > "$scout_file"
  fi
  WHALE_PERF_WALLETS="$scout_file" WHALE_PERF_REPORT="$SCOUT_PERF_REPORT" ./scripts/whale-performance.sh

  local target_file="$TARGET"
  if [ ! -s "$target_file" ]; then
    target_file="$ROOT/db/empty_wallets.txt"
    : > "$target_file"
  fi
  WHALE_PERF_WALLETS="$target_file" WHALE_PERF_REPORT="$TARGET_PERF_REPORT" ./scripts/whale-performance.sh

  local flow_file="$FLOW"
  if [ ! -s "$flow_file" ]; then
    flow_file="$ROOT/db/empty_wallets.txt"
    : > "$flow_file"
  fi
  WHALE_PERF_WALLETS="$flow_file" WHALE_PERF_REPORT="$FLOW_PERF_REPORT" ./scripts/whale-performance.sh

  local tape_file="$TAPE"
  if [ ! -s "$tape_file" ]; then
    tape_file="$ROOT/db/empty_wallets.txt"
    : > "$tape_file"
  fi
  WHALE_PERF_WALLETS="$tape_file" WHALE_PERF_REPORT="$TAPE_PERF_REPORT" ./scripts/whale-performance.sh

  run_push_performance
  run_wallet_maintenance
  filter_push_excludes
}

before_core_hash="$(hash_file "$CORE")"
before_push_hash="$(hash_file "$PUSH")"
before_paper_promoted_hash="$(hash_file "$PAPER_PROMOTED_WALLETS")"
before_paper_demoted_hash="$(hash_file "$PAPER_DEMOTED_WALLETS")"
before_core_count="$(wallet_count "$CORE")"
before_push_count="$(wallet_count "$PUSH")"

pipeline_status=0
set +e
run_pipeline >> "$LOG" 2>&1
pipeline_status=$?
set -e

paper_report_status=0
set +e
SMARTMONEY_PAPER_REPORT_OUT="$PAPER_PNL_REPORT" \
  "$ROOT/scripts/smartmoney-paper-report.sh" >> "$LOG" 2>&1
paper_report_status=$?
set -e
if [ "$pipeline_status" -eq 0 ] && [ "$paper_report_status" -ne 0 ]; then
  pipeline_status="$paper_report_status"
fi

paper_policy_status=0
set +e
"$ROOT/scripts/smartmoney-paper-wallet-policy.sh" >> "$LOG" 2>&1
paper_policy_status=$?
set -e
if [ "$pipeline_status" -eq 0 ] && [ "$paper_policy_status" -ne 0 ]; then
  pipeline_status="$paper_policy_status"
fi

paper_shadow_status=0
set +e
"$ROOT/scripts/smartmoney-shadow-report.sh" >> "$LOG" 2>&1
paper_shadow_status=$?
set -e
if [ "$pipeline_status" -eq 0 ] && [ "$paper_shadow_status" -ne 0 ]; then
  pipeline_status="$paper_shadow_status"
fi

after_core_hash="$(hash_file "$CORE")"
after_push_hash="$(hash_file "$PUSH")"
after_paper_promoted_hash="$(hash_file "$PAPER_PROMOTED_WALLETS")"
after_paper_demoted_hash="$(hash_file "$PAPER_DEMOTED_WALLETS")"
after_core_count="$(wallet_count "$CORE")"
after_watch_count="$(wallet_count "$WATCH")"
after_sports_count="$(wallet_count "$SPORTS")"
after_scout_count="$(wallet_count "$SCOUT")"
after_target_count="$(wallet_count "$TARGET")"
after_flow_count="$(wallet_count "$FLOW")"
after_tape_count="$(wallet_count "$TAPE")"
after_tape_observe_count="$(wallet_count "$TAPE_OBSERVE")"
after_tape_probation_count="$(wallet_count "$TAPE_PROBATION")"
after_tape_candidate_count="$(wallet_count "$TAPE_CANDIDATES")"
after_tape_follow_count="$(wallet_count "$TAPE_FOLLOW")"
after_tape_reversal_count="$(wallet_count "$TAPE_REVERSAL")"
after_tape_edgehot_count="$(wallet_count "$TAPE_EDGEHOT")"
after_consensus_research_count="$(wallet_count "$CONSENSUS_RESEARCH")"
after_push_count="$(wallet_count "$PUSH")"
core_changed='no'
push_changed='no'
paper_policy_changed='no'
if [ "$before_core_hash" != "$after_core_hash" ]; then
  core_changed='yes'
fi
if [ "$before_push_hash" != "$after_push_hash" ]; then
  push_changed='yes'
fi
if [ "$before_paper_promoted_hash" != "$after_paper_promoted_hash" ] || [ "$before_paper_demoted_hash" != "$after_paper_demoted_hash" ]; then
  paper_policy_changed='yes'
fi

paper_restart_status='not needed'
if [ "$pipeline_status" -eq 0 ] && [ "$paper_policy_changed" = "yes" ]; then
  if "$ROOT/scripts/start-smartmoney-paper.sh" restart >> "$LOG" 2>&1; then
    paper_restart_status='restarted after wallet policy change'
  else
    paper_restart_status='restart failed'
    pipeline_status=1
  fi
fi

restart_status='not requested'
restart_needed='no'
if [ "$push_changed" = "yes" ]; then
  restart_needed='yes'
fi
if [ "${SMARTMONEY_FORCE_RESTART:-0}" = "1" ]; then
  restart_needed='yes'
fi

if [ "${SMARTMONEY_RESTART:-0}" = "1" ]; then
  if [ "$pipeline_status" -ne 0 ]; then
    restart_status='skipped after pipeline error'
  elif [ "$push_changed" = "yes" ] || [ "${SMARTMONEY_FORCE_RESTART:-0}" = "1" ]; then
    restart_status="$(restart_worker)"
    restart_needed='no'
  else
    restart_status='skipped; push list unchanged'
  fi
fi

selected_wallets="$(metric "$LAB_REPORT" 'Wallets' | first_field)"
copy_trades="$(metric "$LAB_REPORT" 'Aggregate closed copy trades')"
copy_roi="$(metric "$LAB_REPORT" 'Aggregate copy ROI')"
copy_pnl="$(metric "$LAB_REPORT" 'Aggregate copy PnL')"
copy_win="$(metric "$LAB_REPORT" 'Aggregate copy win rate')"
worst_copy_roi="$(metric "$LAB_REPORT" 'Worst included wallet CopyROI')"
strategy_params="$(metric "$LAB_REPORT" 'Params')"
quarantine_count="$(wallet_count "$QUARANTINE")"
review_noise_count="$(wallet_count "$REVIEW_NOISE")"
edge_snapshots="$(metric "$WHALE_EDGE_REPORT" 'Snapshots' | first_field)"

core_evaluated="$(metric "$CORE_PERF_REPORT" 'Evaluated signals' | first_field)"
core_perf_pnl="$(metric_any "$CORE_PERF_REPORT" 'Proven PnL' 'PnL')"
core_perf_roi="$(metric_any "$CORE_PERF_REPORT" 'Proven ROI' 'ROI')"
sports_evaluated="$(metric "$SPORTS_PERF_REPORT" 'Evaluated signals' | first_field)"
sports_perf_pnl="$(metric_any "$SPORTS_PERF_REPORT" 'Proven PnL' 'PnL')"
sports_perf_roi="$(metric_any "$SPORTS_PERF_REPORT" 'Proven ROI' 'ROI')"
scout_evaluated="$(metric "$SCOUT_PERF_REPORT" 'Evaluated signals' | first_field)"
scout_perf_pnl="$(metric_any "$SCOUT_PERF_REPORT" 'Proven PnL' 'PnL')"
scout_perf_roi="$(metric_any "$SCOUT_PERF_REPORT" 'Proven ROI' 'ROI')"
target_evaluated="$(metric "$TARGET_PERF_REPORT" 'Evaluated signals' | first_field)"
target_perf_pnl="$(metric_any "$TARGET_PERF_REPORT" 'Proven PnL' 'PnL')"
target_perf_roi="$(metric_any "$TARGET_PERF_REPORT" 'Proven ROI' 'ROI')"
flow_evaluated="$(metric "$FLOW_PERF_REPORT" 'Evaluated signals' | first_field)"
flow_perf_pnl="$(metric_any "$FLOW_PERF_REPORT" 'Proven PnL' 'PnL')"
flow_perf_roi="$(metric_any "$FLOW_PERF_REPORT" 'Proven ROI' 'ROI')"
tape_evaluated="$(metric "$TAPE_PERF_REPORT" 'Evaluated signals' | first_field)"
tape_perf_pnl="$(metric_any "$TAPE_PERF_REPORT" 'Proven PnL' 'PnL')"
tape_perf_roi="$(metric_any "$TAPE_PERF_REPORT" 'Proven ROI' 'ROI')"
sports_alerts="$(metric "$SPORTS_ALERT_PERF_REPORT" 'Alerts' | first_field)"
sports_alert_marked="$(metric "$SPORTS_ALERT_PERF_REPORT" 'Marked to current midpoint' | first_field)"
sports_alert_unmarked="$(metric "$SPORTS_ALERT_PERF_REPORT" 'Unmarked' | first_field)"
sports_alert_win="$(metric "$SPORTS_ALERT_PERF_REPORT" 'Win rate incl. midpoint marks')"
sports_alert_pnl="$(metric_any "$SPORTS_ALERT_PERF_REPORT" 'PnL incl. midpoint marks' 'PnL')"
sports_alert_roi="$(metric_any "$SPORTS_ALERT_PERF_REPORT" 'ROI incl. midpoint marks' 'ROI')"
sports_alert_current_modes="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Modes')"
sports_alert_current_alerts="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Alerts' | first_field)"
sports_alert_current_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Marked to current midpoint' | first_field)"
sports_alert_current_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Win rate incl. midpoint marks')"
sports_alert_current_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_current_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_current_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Gate action')"
sports_alert_current_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Performance' 'Gate reason')"
sports_alert_position_modes="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Modes')"
sports_alert_position_positions="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Positions' | first_field)"
sports_alert_position_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Marked to current midpoint' | first_field)"
sports_alert_position_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Win rate incl. midpoint marks')"
sports_alert_position_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_position_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_position_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Gate action')"
sports_alert_position_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Current Policy Position-Capped Performance' 'Gate reason')"
sports_alert_observe_alerts="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'Alerts' | first_field)"
sports_alert_observe_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'Marked to current midpoint' | first_field)"
sports_alert_observe_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'Win rate incl. midpoint marks')"
sports_alert_observe_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_observe_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_observe_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'Gate action')"
sports_alert_observe_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Performance' 'Gate reason')"
sports_alert_observe_positions="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Position-Capped Performance' 'Positions' | first_field)"
sports_alert_observe_position_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_observe_burst_alerts="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Alerts' | first_field)"
sports_alert_observe_burst_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Marked to current midpoint' | first_field)"
sports_alert_observe_burst_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Win rate incl. midpoint marks')"
sports_alert_observe_burst_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_observe_burst_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_observe_burst_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Gate action')"
sports_alert_observe_burst_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Gate reason')"
sports_alert_observe_burst_positions="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Position-Capped Performance' 'Positions' | first_field)"
sports_alert_observe_burst_position_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental OBSERVE-BURST Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_burst_alerts="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Alerts' | first_field)"
sports_alert_shadow_burst_marked="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Marked to current midpoint' | first_field)"
sports_alert_shadow_burst_win="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Win rate incl. midpoint marks')"
sports_alert_shadow_burst_pnl="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_shadow_burst_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_burst_action="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Gate action')"
sports_alert_shadow_burst_reason="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Performance' 'Gate reason')"
sports_alert_shadow_burst_positions="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Position-Capped Performance' 'Positions' | first_field)"
sports_alert_shadow_burst_position_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental OBSERVE-BURST Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_consensus_alerts="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Alerts' | first_field)"
sports_alert_shadow_consensus_marked="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Marked to current midpoint' | first_field)"
sports_alert_shadow_consensus_win="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Win rate incl. midpoint marks')"
sports_alert_shadow_consensus_pnl="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_shadow_consensus_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_consensus_action="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Gate action')"
sports_alert_shadow_consensus_reason="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Gate reason')"
sports_alert_shadow_consensus_positions="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Position-Capped Performance' 'Positions' | first_field)"
sports_alert_shadow_consensus_position_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental CONSENSUS Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_scored_flow_alerts="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'Alerts' | first_field)"
sports_alert_shadow_scored_flow_marked="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'Marked to current midpoint' | first_field)"
sports_alert_shadow_scored_flow_win="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'Win rate incl. midpoint marks')"
sports_alert_shadow_scored_flow_pnl="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_shadow_scored_flow_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_scored_flow_action="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'Gate action')"
sports_alert_shadow_scored_flow_reason="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Performance' 'Gate reason')"
sports_alert_shadow_scored_flow_positions="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Position-Capped Performance' 'Positions' | first_field)"
sports_alert_shadow_scored_flow_position_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SCORED-FLOW Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_seed_flow_alerts="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'Alerts' | first_field)"
sports_alert_shadow_seed_flow_marked="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'Marked to current midpoint' | first_field)"
sports_alert_shadow_seed_flow_win="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'Win rate incl. midpoint marks')"
sports_alert_shadow_seed_flow_pnl="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_shadow_seed_flow_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_seed_flow_action="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'Gate action')"
sports_alert_shadow_seed_flow_reason="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Performance' 'Gate reason')"
sports_alert_shadow_seed_flow_positions="$(section_metric "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Position-Capped Performance' 'Positions' | first_field)"
sports_alert_shadow_seed_flow_position_roi="$(section_metric_any "$SPORTS_ALERT_SHADOW_PERF_REPORT" 'Experimental SEED-FLOW Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_shadow_consensus_marked_n="$(nonnegative_int "$sports_alert_shadow_consensus_marked")"
sports_alert_shadow_consensus_needed="$((SPORTS_ALERT_PROMOTE_MIN_MARKED - sports_alert_shadow_consensus_marked_n))"
if [ "$sports_alert_shadow_consensus_needed" -lt 0 ]; then
  sports_alert_shadow_consensus_needed="0"
fi
if [ "$sports_alert_shadow_consensus_needed" -gt 0 ]; then
  sports_alert_shadow_consensus_readiness="COLLECT (${sports_alert_shadow_consensus_needed} more marked samples before promotion review)"
else
  sports_alert_shadow_consensus_readiness="${sports_alert_shadow_consensus_action:-REVIEW}"
fi
sports_alert_insider_alerts="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'Alerts' | first_field)"
sports_alert_insider_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'Marked to current midpoint' | first_field)"
sports_alert_insider_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'Win rate incl. midpoint marks')"
sports_alert_insider_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_insider_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_insider_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'Gate action')"
sports_alert_insider_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Performance' 'Gate reason')"
sports_alert_insider_positions="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Position-Capped Performance' 'Positions' | first_field)"
sports_alert_insider_position_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental INSIDER-SCOUT Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_consensus_alerts="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Alerts' | first_field)"
sports_alert_consensus_marked="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Marked to current midpoint' | first_field)"
sports_alert_consensus_win="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Win rate incl. midpoint marks')"
sports_alert_consensus_pnl="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'PnL incl. midpoint marks' 'PnL')"
sports_alert_consensus_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_alert_consensus_action="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Gate action')"
sports_alert_consensus_reason="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Performance' 'Gate reason')"
sports_alert_consensus_positions="$(section_metric "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Position-Capped Performance' 'Positions' | first_field)"
sports_alert_consensus_position_roi="$(section_metric_any "$SPORTS_ALERT_PERF_REPORT" 'Experimental CONSENSUS Position-Capped Performance' 'ROI incl. midpoint marks' 'ROI')"
sports_burst_count="$(metric "$SPORTS_BURST_PERF_REPORT" 'Bursts' | first_field)"
sports_burst_marked="$(metric "$SPORTS_BURST_PERF_REPORT" 'Marked to current midpoint' | first_field)"
sports_burst_unmarked="$(metric "$SPORTS_BURST_PERF_REPORT" 'Unmarked' | first_field)"
sports_burst_win="$(metric "$SPORTS_BURST_PERF_REPORT" 'Win rate incl. midpoint marks')"
sports_burst_pnl="$(metric_any "$SPORTS_BURST_PERF_REPORT" 'PnL incl. midpoint marks' 'PnL')"
sports_burst_roi="$(metric_any "$SPORTS_BURST_PERF_REPORT" 'ROI incl. midpoint marks' 'ROI')"
sports_burst_avg_delta="$(metric "$SPORTS_BURST_PERF_REPORT" 'Avg price delta')"
sports_alert_candidates="$(metric "$SPORTS_ALERT_CANDIDATE_REPORT" 'Currently alertable unsent rows' | first_field)"
sports_alert_bursts="0"
sports_alert_consensus_bursts="0"
if [ -s "$SPORTS_ALERT_CANDIDATE_REPORT" ]; then
  sports_alert_bursts="$(awk '
    /^## Accumulation Bursts/ { in_section=1; next }
    /^## / { in_section=0 }
    in_section && /^\| (burst|covered|stale) \|/ { n++ }
    END { print n + 0 }
  ' "$SPORTS_ALERT_CANDIDATE_REPORT")"
  sports_alert_consensus_bursts="$(awk '
    /^## Consensus Bursts/ { in_section=1; next }
    /^## / { in_section=0 }
    in_section && /^\| (consensus|covered|stale) \|/ { n++ }
    END { print n + 0 }
  ' "$SPORTS_ALERT_CANDIDATE_REPORT")"
fi

evaluated="$(metric "$PUSH_PERF_REPORT" 'Evaluated signals' | first_field)"
logged_asset_cooldown="$(metric "$PUSH_PERF_REPORT" 'Logged asset-cooldown BUYs' | first_field)"
logged_event_cooldown="$(metric "$PUSH_PERF_REPORT" 'Logged event-cooldown BUYs' | first_field)"
logged_duplicates="$(metric "$PUSH_PERF_REPORT" 'Logged duplicate BUYs' | first_field)"
realized="$(metric "$PUSH_PERF_REPORT" 'Realized via whale SELL' | first_field)"
settled="$(metric "$PUSH_PERF_REPORT" 'Settled by market resolution' | first_field)"
marked="$(metric "$PUSH_PERF_REPORT" 'Marked to current midpoint' | first_field)"
open_unmarked="$(metric "$PUSH_PERF_REPORT" 'Still open/unmarked' | first_field)"
perf_pnl="$(metric_any "$PUSH_PERF_REPORT" 'PnL incl. midpoint marks' 'PnL')"
perf_roi="$(metric_any "$PUSH_PERF_REPORT" 'ROI incl. midpoint marks' 'ROI')"
perf_roi_num="${perf_roi%\%}"
proven_pnl="$(metric_any "$PUSH_PERF_REPORT" 'Proven PnL')"
proven_roi="$(metric_any "$PUSH_PERF_REPORT" 'Proven ROI')"
event_capped_entries="$(section_metric "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Entries')"
event_capped_realized="$(section_metric "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Realized via whale SELL')"
event_capped_settled="$(section_metric "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Settled by market resolution')"
event_capped_marked="$(section_metric "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Marked to current midpoint')"
event_capped_pnl="$(section_metric_any "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'PnL incl. midpoint marks' 'PnL')"
event_capped_roi="$(section_metric_any "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'ROI incl. midpoint marks' 'ROI')"
event_capped_roi_num="${event_capped_roi%\%}"
event_capped_proven_pnl="$(section_metric_any "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Proven PnL')"
event_capped_proven_roi="$(section_metric_any "$PUSH_PERF_REPORT" 'Event-Capped Strategy' 'Proven ROI')"
policy_violations="$(section_metric "$PUSH_PERF_REPORT" 'Policy Violations' 'Alerted BUYs outside current sports/esports/price policy' | first_field)"

selected_wallets="${selected_wallets:-0}"
copy_trades="${copy_trades:-n/a}"
copy_roi="${copy_roi:-n/a}"
copy_pnl="${copy_pnl:-n/a}"
copy_win="${copy_win:-n/a}"
worst_copy_roi="${worst_copy_roi:-n/a}"
strategy_params="${strategy_params:-n/a}"
quarantine_count="${quarantine_count:-0}"
review_noise_count="${review_noise_count:-0}"
after_consensus_research_count="${after_consensus_research_count:-0}"
after_tape_edgehot_count="${after_tape_edgehot_count:-0}"
edge_snapshots="${edge_snapshots:-0}"
sports_tape_alert_sent="$(json_sent_count "$SPORTS_TAPE_ALERT_STATE")"
sports_tape_alert_logged="$(file_line_count "$SPORTS_TAPE_ALERT_LOG")"
sports_tape_shadow_alert_logged="$(file_line_count "$SPORTS_TAPE_SHADOW_ALERT_LOG")"
sports_consensus_event_logged="$(file_line_count "$SPORTS_CONSENSUS_EVENTS")"
sports_consensus_watch_event_logged="$(file_line_count "$SPORTS_CONSENSUS_WATCH_EVENTS")"
core_evaluated="${core_evaluated:-0}"
core_perf_pnl="${core_perf_pnl:-n/a}"
core_perf_roi="${core_perf_roi:-0.0%}"
sports_evaluated="${sports_evaluated:-0}"
sports_perf_pnl="${sports_perf_pnl:-n/a}"
sports_perf_roi="${sports_perf_roi:-0.0%}"
scout_evaluated="${scout_evaluated:-0}"
scout_perf_pnl="${scout_perf_pnl:-n/a}"
scout_perf_roi="${scout_perf_roi:-0.0%}"
target_evaluated="${target_evaluated:-0}"
target_perf_pnl="${target_perf_pnl:-n/a}"
target_perf_roi="${target_perf_roi:-0.0%}"
flow_evaluated="${flow_evaluated:-0}"
flow_perf_pnl="${flow_perf_pnl:-n/a}"
flow_perf_roi="${flow_perf_roi:-0.0%}"
tape_evaluated="${tape_evaluated:-0}"
tape_perf_pnl="${tape_perf_pnl:-n/a}"
tape_perf_roi="${tape_perf_roi:-0.0%}"
sports_alerts="${sports_alerts:-0}"
sports_alert_marked="${sports_alert_marked:-0}"
sports_alert_unmarked="${sports_alert_unmarked:-0}"
sports_alert_win="${sports_alert_win:-0.0%}"
sports_alert_pnl="${sports_alert_pnl:-n/a}"
sports_alert_roi="${sports_alert_roi:-0.0%}"
sports_alert_current_modes="${sports_alert_current_modes:-n/a}"
sports_alert_current_alerts="${sports_alert_current_alerts:-0}"
sports_alert_current_marked="${sports_alert_current_marked:-0}"
sports_alert_current_win="${sports_alert_current_win:-0.0%}"
sports_alert_current_pnl="${sports_alert_current_pnl:-n/a}"
sports_alert_current_roi="${sports_alert_current_roi:-0.0%}"
sports_alert_current_action="${sports_alert_current_action:-COLLECT}"
sports_alert_current_reason="${sports_alert_current_reason:-n/a}"
sports_alert_position_modes="${sports_alert_position_modes:-n/a}"
sports_alert_position_positions="${sports_alert_position_positions:-0}"
sports_alert_position_marked="${sports_alert_position_marked:-0}"
sports_alert_position_win="${sports_alert_position_win:-0.0%}"
sports_alert_position_pnl="${sports_alert_position_pnl:-n/a}"
sports_alert_position_roi="${sports_alert_position_roi:-0.0%}"
sports_alert_position_action="${sports_alert_position_action:-COLLECT}"
sports_alert_position_reason="${sports_alert_position_reason:-n/a}"
sports_alert_observe_alerts="${sports_alert_observe_alerts:-0}"
sports_alert_observe_marked="${sports_alert_observe_marked:-0}"
sports_alert_observe_win="${sports_alert_observe_win:-0.0%}"
sports_alert_observe_pnl="${sports_alert_observe_pnl:-n/a}"
sports_alert_observe_roi="${sports_alert_observe_roi:-0.0%}"
sports_alert_observe_action="${sports_alert_observe_action:-COLLECT}"
sports_alert_observe_reason="${sports_alert_observe_reason:-n/a}"
sports_alert_observe_positions="${sports_alert_observe_positions:-0}"
sports_alert_observe_position_roi="${sports_alert_observe_position_roi:-0.0%}"
sports_alert_observe_burst_alerts="${sports_alert_observe_burst_alerts:-0}"
sports_alert_observe_burst_marked="${sports_alert_observe_burst_marked:-0}"
sports_alert_observe_burst_win="${sports_alert_observe_burst_win:-0.0%}"
sports_alert_observe_burst_pnl="${sports_alert_observe_burst_pnl:-n/a}"
sports_alert_observe_burst_roi="${sports_alert_observe_burst_roi:-0.0%}"
sports_alert_observe_burst_action="${sports_alert_observe_burst_action:-COLLECT}"
sports_alert_observe_burst_reason="${sports_alert_observe_burst_reason:-n/a}"
sports_alert_observe_burst_positions="${sports_alert_observe_burst_positions:-0}"
sports_alert_observe_burst_position_roi="${sports_alert_observe_burst_position_roi:-0.0%}"
sports_alert_insider_alerts="${sports_alert_insider_alerts:-0}"
sports_alert_insider_marked="${sports_alert_insider_marked:-0}"
sports_alert_insider_win="${sports_alert_insider_win:-0.0%}"
sports_alert_insider_pnl="${sports_alert_insider_pnl:-n/a}"
sports_alert_insider_roi="${sports_alert_insider_roi:-0.0%}"
sports_alert_insider_action="${sports_alert_insider_action:-COLLECT}"
sports_alert_insider_reason="${sports_alert_insider_reason:-n/a}"
sports_alert_insider_positions="${sports_alert_insider_positions:-0}"
sports_alert_insider_position_roi="${sports_alert_insider_position_roi:-0.0%}"
sports_alert_consensus_alerts="${sports_alert_consensus_alerts:-0}"
sports_alert_consensus_marked="${sports_alert_consensus_marked:-0}"
sports_alert_consensus_win="${sports_alert_consensus_win:-0.0%}"
sports_alert_consensus_pnl="${sports_alert_consensus_pnl:-n/a}"
sports_alert_consensus_roi="${sports_alert_consensus_roi:-0.0%}"
sports_alert_consensus_action="${sports_alert_consensus_action:-COLLECT}"
sports_alert_consensus_reason="${sports_alert_consensus_reason:-n/a}"
sports_alert_consensus_positions="${sports_alert_consensus_positions:-0}"
sports_alert_consensus_position_roi="${sports_alert_consensus_position_roi:-0.0%}"
sports_burst_count="${sports_burst_count:-0}"
sports_burst_marked="${sports_burst_marked:-0}"
sports_burst_unmarked="${sports_burst_unmarked:-0}"
sports_burst_win="${sports_burst_win:-0.0%}"
sports_burst_pnl="${sports_burst_pnl:-n/a}"
sports_burst_roi="${sports_burst_roi:-0.0%}"
sports_burst_avg_delta="${sports_burst_avg_delta:-n/a}"
sports_alert_candidates="${sports_alert_candidates:-0}"
sports_alert_bursts="${sports_alert_bursts:-0}"
sports_alert_consensus_bursts="${sports_alert_consensus_bursts:-0}"
sports_consensus_event_logged="${sports_consensus_event_logged:-0}"
evaluated="${evaluated:-0}"
logged_asset_cooldown="${logged_asset_cooldown:-0}"
logged_event_cooldown="${logged_event_cooldown:-0}"
logged_duplicates="${logged_duplicates:-0}"
realized="${realized:-0}"
settled="${settled:-0}"
marked="${marked:-0}"
open_unmarked="${open_unmarked:-0}"
perf_pnl="${perf_pnl:-n/a}"
perf_roi="${perf_roi:-0.0%}"
perf_roi_num="${perf_roi_num:-0}"
proven_pnl="${proven_pnl:-n/a}"
proven_roi="${proven_roi:-0.0%}"
event_capped_entries="${event_capped_entries:-0}"
event_capped_realized="${event_capped_realized:-0}"
event_capped_settled="${event_capped_settled:-0}"
event_capped_marked="${event_capped_marked:-0}"
event_capped_pnl="${event_capped_pnl:-n/a}"
event_capped_roi="${event_capped_roi:-0.0%}"
event_capped_roi_num="${event_capped_roi_num:-0}"
event_capped_proven_pnl="${event_capped_proven_pnl:-n/a}"
event_capped_proven_roi="${event_capped_proven_roi:-0.0%}"
policy_violations="${policy_violations:-0}"
proven_signals=$((realized + 0 + settled))
proven_event_capped=$((event_capped_realized + 0 + event_capped_settled))

health='learning'
health_reason="evaluated signals below $MIN_SIGNALS"
if [ "$pipeline_status" -ne 0 ]; then
  health='pipeline_error'
  health_reason="pipeline exited with status $pipeline_status"
elif [ "$policy_violations" -gt 0 ]; then
  health='policy_violation'
  health_reason="$policy_violations alerted BUYs are outside current sports/esports/price policy"
elif [ "$evaluated" -lt "$MIN_SIGNALS" ]; then
  health='learning'
  health_reason="evaluated signals $evaluated below $MIN_SIGNALS"
elif [ "$event_capped_entries" -lt "$MIN_EVENT_CAPPED_ENTRIES" ]; then
  health='learning'
  health_reason="event-capped entries $event_capped_entries below $MIN_EVENT_CAPPED_ENTRIES"
else
  if awk "BEGIN {exit !(($perf_roi_num >= $MIN_ROI) && ($event_capped_roi_num >= $MIN_ROI))}"; then
    if [ "$proven_signals" -ge "$MIN_PROVEN_SIGNALS" ] && [ "$proven_event_capped" -ge "$MIN_PROVEN_EVENT_CAPPED_ENTRIES" ]; then
      health='pass'
      health_reason="verified ROI $perf_roi and event-capped ROI $event_capped_roi with proven signals $proven_signals/${MIN_PROVEN_SIGNALS} and proven event-capped entries $proven_event_capped/${MIN_PROVEN_EVENT_CAPPED_ENTRIES}"
    else
      health='provisional_pass'
      health_reason="marked ROI $perf_roi and event-capped ROI $event_capped_roi are positive, but proven signals $proven_signals/${MIN_PROVEN_SIGNALS} and proven event-capped entries $proven_event_capped/${MIN_PROVEN_EVENT_CAPPED_ENTRIES} are not enough"
    fi
  else
    health='degrade'
    health_reason="ROI $perf_roi or event-capped ROI $event_capped_roi below threshold ${MIN_ROI}%"
  fi
fi

action='keep current push running and collect more signals'
if [ "$health" = "pipeline_error" ]; then
  action="inspect $LOG before changing the running worker"
elif [ "$health" = "policy_violation" ]; then
  action='pause wallet expansion; inspect policy violations and tighten runtime filters before restarting'
elif [ "$health" = "degrade" ]; then
  action='do not expand copied wallets; tighten thresholds or pause weak wallets'
elif [ "$health" = "provisional_pass" ]; then
  action='keep push-only mode; wait for closed or settled proof before promoting wallets'
elif [ "$restart_needed" = "yes" ]; then
  action='restart whale push so the worker loads the new push wallet list'
fi

{
  printf '# Smart Money Daily\n\n'
  printf '**Generated:** %s\n\n' "$(date '+%Y-%m-%d %H:%M:%S %Z')"
  printf '## Pipeline\n\n'
  printf -- '- Status: %s\n' "$health"
  printf -- '- Reason: %s\n' "$health_reason"
  printf -- '- Action: %s\n' "$action"
  printf -- '- Log: `%s`\n' "$LOG"
  printf -- '- Discovery report: `%s`\n' "$DISCOVERY_REPORT"
  printf -- '- Strategy report: `%s`\n' "$LAB_REPORT"
  printf -- '- Sports tape report: `%s`\n' "$SPORTS_TAPE_REPORT"
  printf -- '- Sports alert performance report: `%s`\n' "$SPORTS_ALERT_PERF_REPORT"
  printf -- '- Sports alert candidate report: `%s`\n' "$SPORTS_ALERT_CANDIDATE_REPORT"
  printf -- '- Sports burst performance report: `%s`\n' "$SPORTS_BURST_PERF_REPORT"
  printf -- '- Whale edge report: `%s`\n' "$WHALE_EDGE_REPORT"
  printf -- '- Core performance report: `%s`\n' "$CORE_PERF_REPORT"
  printf -- '- Sports performance report: `%s`\n' "$SPORTS_PERF_REPORT"
  printf -- '- Scout performance report: `%s`\n' "$SCOUT_PERF_REPORT"
  printf -- '- Target performance report: `%s`\n' "$TARGET_PERF_REPORT"
  printf -- '- Flow performance report: `%s`\n' "$FLOW_PERF_REPORT"
  printf -- '- Tape performance report: `%s`\n' "$TAPE_PERF_REPORT"
  printf -- '- Push performance report: `%s`\n' "$PUSH_PERF_REPORT"
  printf -- '- Paper PnL report: `%s` (status %s)\n' "$PAPER_PNL_REPORT" "$paper_report_status"
  printf -- '- Paper wallet policy: `%s` (status %s, changed %s, worker %s)\n' "$PAPER_WALLET_POLICY_REPORT" "$paper_policy_status" "$paper_policy_changed" "$paper_restart_status"
  printf -- '- Paper exit shadow: `%s` (status %s)\n' "$PAPER_SHADOW_REPORT" "$paper_shadow_status"
  printf -- '- Maintenance report: `%s`\n\n' "$MAINT_REPORT"

  printf '## Wallet Lists\n\n'
  printf -- '- Core before count: %s\n' "$before_core_count"
  printf -- '- Core after count: %s\n' "$after_core_count"
  printf -- '- Watch after count: %s\n' "$after_watch_count"
  printf -- '- Sports after count: %s\n' "$after_sports_count"
  printf -- '- Scout after count: %s\n' "$after_scout_count"
  printf -- '- Target after count: %s\n' "$after_target_count"
  printf -- '- Flow after count: %s\n' "$after_flow_count"
  printf -- '- Tape after count: %s\n' "$after_tape_count"
  printf -- '- Tape observe count: %s\n' "$after_tape_observe_count"
  printf -- '- Tape probation count: %s\n' "$after_tape_probation_count"
  printf -- '- Tape candidate count: %s\n' "$after_tape_candidate_count"
  printf -- '- Tape follow-ready count: %s\n' "$after_tape_follow_count"
  printf -- '- Tape edge-hot count: %s\n' "$after_tape_edgehot_count"
  printf -- '- Tape reversal count: %s\n' "$after_tape_reversal_count"
  printf -- '- Consensus research count: %s\n' "$after_consensus_research_count"
  printf -- '- Push before count: %s\n' "$before_push_count"
  printf -- '- Push after count: %s\n' "$after_push_count"
  printf -- '- Quarantine count: %s\n' "$quarantine_count"
  printf -- '- Review-noise exclude count: %s\n' "$review_noise_count"
  printf -- '- Edge snapshots: %s\n' "$edge_snapshots"
  printf -- '- Sports tape alert sent: %s\n' "$sports_tape_alert_sent"
  printf -- '- Sports tape alert logged: %s\n' "$sports_tape_alert_logged"
  printf -- '- Sports tape shadow alert logged: %s\n' "$sports_tape_shadow_alert_logged"
  printf -- '- Sports consensus event history: %s\n' "$sports_consensus_event_logged"
  printf -- '- Sports consensus watch event history: %s\n' "$sports_consensus_watch_event_logged"
  printf -- '- Sports alert eligible now: %s\n' "$sports_alert_candidates"
  printf -- '- Sports alert accumulation bursts: %s\n' "$sports_alert_bursts"
  printf -- '- Sports alert consensus bursts: %s\n' "$sports_alert_consensus_bursts"
  printf -- '- Selected core wallets: %s\n' "$selected_wallets"
  printf -- '- Core list changed: %s\n' "$core_changed"
  printf -- '- Push list changed: %s\n' "$push_changed"
  printf -- '- Restart needed: %s\n' "$restart_needed"
  printf -- '- Restart status: %s\n\n' "$restart_status"

  printf '## Backtest Filter\n\n'
  printf -- '- Closed copy trades: %s\n' "$copy_trades"
  printf -- '- Copy ROI: %s\n' "$copy_roi"
  printf -- '- Copy PnL: %s\n' "$copy_pnl"
  printf -- '- Copy win rate: %s\n' "$copy_win"
  printf -- '- Worst included CopyROI: %s\n' "$worst_copy_roi"
  printf -- '- Params: %s\n\n' "$strategy_params"

  printf '## Live Core Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$core_evaluated"
  printf -- '- PnL: %s\n' "$core_perf_pnl"
  printf -- '- ROI: %s\n\n' "$core_perf_roi"

  printf '## Live Sports Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$sports_evaluated"
  printf -- '- PnL: %s\n' "$sports_perf_pnl"
  printf -- '- ROI: %s\n\n' "$sports_perf_roi"

  printf '## Sports Alert Performance\n\n'
  printf -- '- Alerts: %s\n' "$sports_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_marked"
  printf -- '- Unmarked: %s\n' "$sports_alert_unmarked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n\n' "$sports_alert_roi"

  printf '## Sports Alert Current Policy\n\n'
  printf -- '- Modes: %s\n' "$sports_alert_current_modes"
  printf -- '- Alerts: %s\n' "$sports_alert_current_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_current_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_current_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_current_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_current_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_current_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_current_reason"

  printf '## Sports Alert Position-Capped Policy\n\n'
  printf -- '- Rule: first alert per wallet + asset\n'
  printf -- '- Modes: %s\n' "$sports_alert_position_modes"
  printf -- '- Positions: %s\n' "$sports_alert_position_positions"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_position_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_position_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_position_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_position_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_position_reason"

  printf '## Sports Alert OBSERVE Experiment\n\n'
  printf -- '- Rule: raw low-bot sports/esports whale BUY alerts, including same-wallet split-order bursts; not counted in current policy until promoted\n'
  printf -- '- Default min notional: $%s\n' "$SPORTS_TAPE_ALERT_OBSERVE_MIN_NOTIONAL_EFFECTIVE"
  printf -- '- Observe-burst min notional: $%s\n' "${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}"
  printf -- '- Require scored/listed wallet: %s\n' "$SPORTS_TAPE_ALERT_OBSERVE_REQUIRE_KNOWN_EFFECTIVE"
  printf -- '- Minimum tier: %s\n' "$SPORTS_TAPE_ALERT_OBSERVE_MIN_TIER_EFFECTIVE"
  printf -- '- Insider-scout min notional: $%s\n' "${SPORTS_TAPE_ALERT_INSIDER_MIN_NOTIONAL:-25000}"
  printf -- '- Insider-scout max bot: %s\n' "${SPORTS_TAPE_ALERT_INSIDER_MAX_BOT:-35}"
  printf -- '- Alerts: %s\n' "$sports_alert_observe_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_observe_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_observe_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_observe_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_observe_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_observe_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_observe_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_observe_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_observe_reason"

  printf '## Sports Alert OBSERVE-BURST Experiment\n\n'
  printf -- '- Rule: same-wallet split-order sports/esports BUY bursts; observation only until repeated positive ROI is proven\n'
  printf -- '- Min cumulative notional: $%s\n' "${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}"
  printf -- '- Alerts: %s\n' "$sports_alert_observe_burst_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_observe_burst_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_observe_burst_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_observe_burst_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_observe_burst_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_observe_burst_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_observe_burst_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_observe_burst_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_observe_burst_reason"

  printf '## Sports Alert SHADOW OBSERVE-BURST Experiment\n\n'
  printf -- '- Rule: stale same-wallet split-order bursts captured for research only; not counted as Telegram pushes or current policy\n'
  printf -- '- Shadow log: %s\n' "$SPORTS_TAPE_SHADOW_ALERT_LOG"
  printf -- '- Logged events: %s\n' "$sports_tape_shadow_alert_logged"
  printf -- '- Alerts: %s\n' "$sports_alert_shadow_burst_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_shadow_burst_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_shadow_burst_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_shadow_burst_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_shadow_burst_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_shadow_burst_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_shadow_burst_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_shadow_burst_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_shadow_burst_reason"

  printf '## Sports Alert SHADOW CONSENSUS Experiment\n\n'
  printf -- '- Rule: stale cross-wallet same-asset sports/esports BUY bursts captured for research only; not counted as Telegram pushes or current policy\n'
  printf -- '- Promote gate: marked>=%s, ROI>=%s%%, win>=%s%%\n' "$SPORTS_ALERT_PROMOTE_MIN_MARKED" "$SPORTS_ALERT_PROMOTE_MIN_ROI" "$SPORTS_ALERT_PROMOTE_MIN_WIN"
  printf -- '- Alerts: %s\n' "$sports_alert_shadow_consensus_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_shadow_consensus_marked"
  printf -- '- Marked samples still needed: %s\n' "$sports_alert_shadow_consensus_needed"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_shadow_consensus_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_shadow_consensus_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_shadow_consensus_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_shadow_consensus_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_shadow_consensus_position_roi"
  printf -- '- Readiness: %s\n' "$sports_alert_shadow_consensus_readiness"
  printf -- '- Gate action: %s\n' "$sports_alert_shadow_consensus_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_shadow_consensus_reason"

  printf '## Sports Alert SHADOW SEED-FLOW Experiment\n\n'
  printf -- '- Rule: lower-threshold unknown wallets buying multiple sports/esports markets; research-only before UNKNOWN-FLOW size is reached\n'
  printf -- '- Min cumulative notional: $%s\n' "${SPORTS_TAPE_SEED_FLOW_MIN_NOTIONAL:-3000}"
  printf -- '- Min markets: %s\n' "${SPORTS_TAPE_SEED_FLOW_MIN_MARKETS:-2}"
  printf -- '- Alerts: %s\n' "$sports_alert_shadow_seed_flow_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_shadow_seed_flow_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_shadow_seed_flow_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_shadow_seed_flow_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_shadow_seed_flow_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_shadow_seed_flow_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_shadow_seed_flow_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_shadow_seed_flow_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_shadow_seed_flow_reason"

  printf '## Sports Alert SHADOW SCORED-FLOW Experiment\n\n'
  printf -- '- Rule: scored low-bot leaderboard wallets buying multiple sports/esports markets; shadow-only until repeated positive ROI is proven\n'
  printf -- '- Min cumulative notional: $%s\n' "${SPORTS_TAPE_SCORED_FLOW_MIN_NOTIONAL:-4000}"
  printf -- '- Min markets: %s\n' "${SPORTS_TAPE_SCORED_FLOW_MIN_MARKETS:-2}"
  printf -- '- Min tier: %s\n' "${SPORTS_TAPE_SCORED_FLOW_MIN_TIER:-B}"
  printf -- '- Max bot: %s\n' "${SPORTS_TAPE_SCORED_FLOW_MAX_BOT:-35}"
  printf -- '- Alerts: %s\n' "$sports_alert_shadow_scored_flow_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_shadow_scored_flow_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_shadow_scored_flow_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_shadow_scored_flow_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_shadow_scored_flow_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_shadow_scored_flow_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_shadow_scored_flow_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_shadow_scored_flow_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_shadow_scored_flow_reason"

  printf '## Sports Alert INSIDER-SCOUT Experiment\n\n'
  printf -- '- Rule: 25k+ very large low-bot sports/esports whale BUY alerts; observation only until repeated positive ROI is proven\n'
  printf -- '- Alerts: %s\n' "$sports_alert_insider_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_insider_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_insider_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_insider_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_insider_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_insider_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_insider_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_insider_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_insider_reason"

  printf '## Sports Alert CONSENSUS Experiment\n\n'
  printf -- '- Rule: cross-wallet same-asset sports/esports BUY bursts; research/observation only until repeated positive ROI is proven\n'
  printf -- '- Alerts: %s\n' "$sports_alert_consensus_alerts"
  printf -- '- Marked to current midpoint: %s\n' "$sports_alert_consensus_marked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_alert_consensus_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_alert_consensus_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_alert_consensus_roi"
  printf -- '- Position-capped entries: %s\n' "$sports_alert_consensus_positions"
  printf -- '- Position-capped ROI: %s\n' "$sports_alert_consensus_position_roi"
  printf -- '- Gate action: %s\n' "$sports_alert_consensus_action"
  printf -- '- Gate reason: %s\n\n' "$sports_alert_consensus_reason"

  printf '## Sports Burst Performance\n\n'
  printf -- '- Bursts: %s\n' "$sports_burst_count"
  printf -- '- Marked to current midpoint: %s\n' "$sports_burst_marked"
  printf -- '- Unmarked: %s\n' "$sports_burst_unmarked"
  printf -- '- Win rate incl. midpoint marks: %s\n' "$sports_burst_win"
  printf -- '- PnL incl. midpoint marks: %s\n' "$sports_burst_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$sports_burst_roi"
  printf -- '- Avg price delta: %s\n\n' "$sports_burst_avg_delta"

  printf '## Live Scout Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$scout_evaluated"
  printf -- '- PnL: %s\n' "$scout_perf_pnl"
  printf -- '- ROI: %s\n\n' "$scout_perf_roi"

  printf '## Live Target Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$target_evaluated"
  printf -- '- PnL: %s\n' "$target_perf_pnl"
  printf -- '- ROI: %s\n\n' "$target_perf_roi"

  printf '## Live Flow Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$flow_evaluated"
  printf -- '- PnL: %s\n' "$flow_perf_pnl"
  printf -- '- ROI: %s\n\n' "$flow_perf_roi"

  printf '## Live Tape Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$tape_evaluated"
  printf -- '- PnL: %s\n' "$tape_perf_pnl"
  printf -- '- ROI: %s\n\n' "$tape_perf_roi"

  printf '## Live Push Performance\n\n'
  printf -- '- Evaluated signals: %s\n' "$evaluated"
  printf -- '- Event-capped entries: %s\n' "$event_capped_entries"
  printf -- '- Policy violations: %s\n' "$policy_violations"
  printf -- '- Proven signals: %s\n' "$proven_signals"
  printf -- '- Proven event-capped entries: %s\n' "$proven_event_capped"
  printf -- '- Logged asset-cooldown BUYs: %s\n' "$logged_asset_cooldown"
  printf -- '- Logged event-cooldown BUYs: %s\n' "$logged_event_cooldown"
  printf -- '- Logged duplicate BUYs: %s\n' "$logged_duplicates"
  printf -- '- Realized via whale SELL: %s\n' "$realized"
  printf -- '- Settled by market resolution: %s\n' "$settled"
  printf -- '- Marked to current midpoint: %s\n' "$marked"
  printf -- '- Still open/unmarked: %s\n' "$open_unmarked"
  printf -- '- PnL incl. midpoint marks: %s\n' "$perf_pnl"
  printf -- '- ROI incl. midpoint marks: %s\n' "$perf_roi"
  printf -- '- Proven PnL: %s\n' "$proven_pnl"
  printf -- '- Proven ROI: %s\n' "$proven_roi"
  printf -- '- Event-capped realized via whale SELL: %s\n' "$event_capped_realized"
  printf -- '- Event-capped settled by market resolution: %s\n' "$event_capped_settled"
  printf -- '- Event-capped marked to current midpoint: %s\n' "$event_capped_marked"
  printf -- '- Event-capped PnL incl. midpoint marks: %s\n' "$event_capped_pnl"
  printf -- '- Event-capped ROI incl. midpoint marks: %s\n' "$event_capped_roi"
  printf -- '- Event-capped proven PnL: %s\n' "$event_capped_proven_pnl"
  printf -- '- Event-capped proven ROI: %s\n' "$event_capped_proven_roi"
} > "$SUMMARY"

printf 'smartmoney daily status=%s report=%s log=%s\n' "$health" "$SUMMARY" "$LOG"

send_policy_violation_alert "$policy_violations"

if [ "$pipeline_status" -ne 0 ]; then
  exit "$pipeline_status"
fi
