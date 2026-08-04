#!/usr/bin/env bash
# Health check for the active smartmoney-paper and whale-push screen runners.
set -euo pipefail

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p "$ROOT/db/launchd-health" "$ROOT/logs"

LOCK_DIR="$ROOT/db/launchd-health.lock"
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
	lock_pid="$(cat "$LOCK_DIR/pid" 2>/dev/null || true)"
	if [ -n "$lock_pid" ] && kill -0 "$lock_pid" 2>/dev/null; then
		lock_command="$(ps -o command= -p "$lock_pid" 2>/dev/null || true)"
		case "$lock_command" in
			*smartmoney-health-check.sh*|"")
				printf 'health.already_running lock=%s pid=%s\n' "$LOCK_DIR" "$lock_pid"
				exit 0
				;;
		esac
	fi
	rm -f "$LOCK_DIR/pid"
	rmdir "$LOCK_DIR" 2>/dev/null || { printf 'health.lock_unavailable lock=%s\n' "$LOCK_DIR"; exit 0; }
	mkdir "$LOCK_DIR" 2>/dev/null || { printf 'health.already_running lock=%s\n' "$LOCK_DIR"; exit 0; }
fi
printf '%s\n' "$$" > "$LOCK_DIR/pid"
trap 'rm -f "$LOCK_DIR/pid"; rmdir "$LOCK_DIR" 2>/dev/null || true' EXIT

STATE="$ROOT/db/launchd-health/state.json"
MAX_LOG_AGE_SEC="${POLYMARKET_HEALTH_MAX_LOG_AGE_SEC:-300}"
RESTART_ENABLED="${POLYMARKET_HEALTH_RESTART:-1}"
SMART_SCREEN="${SMARTMONEY_PAPER_SCREEN_NAME:-polymarket-smartmoney-paper}"
WHALE_SCREEN="${WHALE_PUSH_SCREEN_NAME:-polymarket-whale-push}"

now_epoch="$(date +%s)"
checked_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

bool() {
  if [ "${1:-0}" = "1" ]; then
    printf 'true'
  else
    printf 'false'
  fi
}

file_age_sec() {
  local path="$1"
  local mtime
  if [ ! -f "$path" ]; then
    printf '999999999'
    return 0
  fi
  mtime="$(stat -f '%m' "$path" 2>/dev/null || printf '0')"
  if [ -z "$mtime" ] || [ "$mtime" = "0" ]; then
    printf '999999999'
    return 0
  fi
  printf '%s' "$((now_epoch - mtime))"
}

screen_running() {
  local name="$1"
  local out
  out="$(screen -ls 2>/dev/null || true)"
  printf '%s\n' "$out" | grep -Eq "[0-9]+[.]${name}[[:space:]]"
}

pid_matches() {
  local file="$1"
  local pattern="$2"
  local pid cmd
  pid="$(cat "$file" 2>/dev/null || true)"
  [ -n "$pid" ] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  cmd="$(ps -o command= -p "$pid" 2>/dev/null || true)"
  case "$cmd" in
    *"$pattern"*) return 0 ;;
    *) return 1 ;;
  esac
}

restart_smartmoney() {
  printf 'health.restart component=smartmoney-paper reason=%s\n' "$smart_reason"
  "$ROOT/scripts/start-smartmoney-paper.sh" restart >/dev/null 2>&1 || true
}

restart_whale_push() {
  local old_pid child_pid
  printf 'health.restart component=whale-push reason=%s\n' "$whale_reason"
  screen -S "$WHALE_SCREEN" -X quit >/dev/null 2>&1 || true
  child_pid="$(cat "$ROOT/db/whale-push.child.pid" 2>/dev/null || true)"
  old_pid="$(cat "$ROOT/db/whale-push.pid" 2>/dev/null || true)"
  if [ -n "$child_pid" ] && pid_matches "$ROOT/db/whale-push.child.pid" "-signal_mode=whale"; then
    kill "$child_pid" >/dev/null 2>&1 || true
  fi
  if [ -n "$old_pid" ] && pid_matches "$ROOT/db/whale-push.pid" "start-whale-push.sh"; then
    kill "$old_pid" >/dev/null 2>&1 || true
  fi
  sleep 1
  screen -dmS "$WHALE_SCREEN" "$ROOT/scripts/start-whale-push.sh"
}

smart_status="$("$ROOT/scripts/start-smartmoney-paper.sh" status 2>&1 | head -1 || true)"
smart_running=0
case "$smart_status" in
  RUNNING*) smart_running=1 ;;
esac
smart_screen=0
screen_running "$SMART_SCREEN" && smart_screen=1
smart_log_age="$(file_age_sec "$ROOT/logs/smartmoney-paper.log")"
smart_blocked=0
if grep -q '"blocked"[[:space:]]*:[[:space:]]*true' "$ROOT/db/smartmoney-paper/risk_state.json" 2>/dev/null; then
  smart_blocked=1
fi

smart_reason="ok"
if [ "$smart_running" -ne 1 ]; then
  smart_reason="not_running"
elif [ "$smart_log_age" -gt "$MAX_LOG_AGE_SEC" ]; then
  smart_reason="stale_log"
elif [ "$smart_blocked" -eq 1 ]; then
  smart_reason="risk_blocked"
fi

if [ "$RESTART_ENABLED" = "1" ] && { [ "$smart_reason" = "not_running" ] || [ "$smart_reason" = "stale_log" ]; }; then
  restart_smartmoney
  smart_status="$("$ROOT/scripts/start-smartmoney-paper.sh" status 2>&1 | head -1 || true)"
  smart_running=0
  case "$smart_status" in
    RUNNING*) smart_running=1 ;;
  esac
  smart_screen=0
  screen_running "$SMART_SCREEN" && smart_screen=1
  smart_log_age="$(file_age_sec "$ROOT/logs/smartmoney-paper.log")"
  if [ "$smart_running" -eq 1 ]; then
    smart_reason="restarted"
  else
    smart_reason="restart_failed"
  fi
fi

whale_wrapper=0
whale_child=0
whale_screen=0
pid_matches "$ROOT/db/whale-push.pid" "start-whale-push.sh" && whale_wrapper=1
pid_matches "$ROOT/db/whale-push.child.pid" "-signal_mode=whale" && whale_child=1
screen_running "$WHALE_SCREEN" && whale_screen=1
whale_log_age="$(file_age_sec "$ROOT/db/agent.log")"

whale_reason="ok"
if [ "$whale_wrapper" -ne 1 ] || [ "$whale_child" -ne 1 ]; then
  whale_reason="not_running"
elif [ "$whale_log_age" -gt "$MAX_LOG_AGE_SEC" ]; then
  whale_reason="stale_log"
fi

if [ "$RESTART_ENABLED" = "1" ] && { [ "$whale_reason" = "not_running" ] || [ "$whale_reason" = "stale_log" ]; }; then
  restart_whale_push
  whale_wrapper=0
  whale_child=0
  whale_screen=0
  pid_matches "$ROOT/db/whale-push.pid" "start-whale-push.sh" && whale_wrapper=1
  pid_matches "$ROOT/db/whale-push.child.pid" "-signal_mode=whale" && whale_child=1
  screen_running "$WHALE_SCREEN" && whale_screen=1
  whale_log_age="$(file_age_sec "$ROOT/db/agent.log")"
  if [ "$whale_wrapper" -eq 1 ] && [ "$whale_child" -eq 1 ]; then
    whale_reason="restarted"
  else
    whale_reason="restart_failed"
  fi
fi

live_trading_disabled=0
live_trading_arm_present=0
live_redeem_disabled=0
live_redeem_arm_present=0
research_disabled=0
monitoring_disabled=0
legacy_project_disabled=0
[ -f "$ROOT/db/live-trading.disabled" ] && live_trading_disabled=1
[ -f "$ROOT/db/live-trading.enabled" ] && live_trading_arm_present=1
[ -f "$ROOT/db/live/redeem.disabled" ] && live_redeem_disabled=1
[ -f "$ROOT/db/live/redeem.enabled" ] && live_redeem_arm_present=1
[ -f "$ROOT/db/research.disabled" ] && research_disabled=1
[ -f "$ROOT/db/monitoring.disabled" ] && monitoring_disabled=1
[ -f "$ROOT/db/project.disabled" ] && legacy_project_disabled=1

overall="ok"
case "$smart_reason:$whale_reason" in
  *restart_failed*|*risk_blocked*) overall="warn" ;;
  *not_running*|*stale_log*) overall="warn" ;;
  *restarted*) overall="recovered" ;;
esac

cat > "$STATE" <<EOF
{
  "checked_at": "$checked_at",
  "overall": "$overall",
  "component_flags": {
    "live_trading_disabled": $(bool "$live_trading_disabled"),
    "live_trading_arm_present": $(bool "$live_trading_arm_present"),
    "live_redeem_disabled": $(bool "$live_redeem_disabled"),
    "live_redeem_arm_present": $(bool "$live_redeem_arm_present"),
    "research_disabled": $(bool "$research_disabled"),
    "monitoring_disabled": $(bool "$monitoring_disabled"),
    "legacy_project_disabled_present": $(bool "$legacy_project_disabled")
  },
  "restart_enabled": $(bool "$([ "$RESTART_ENABLED" = "1" ] && printf 1 || printf 0)"),
  "max_log_age_sec": $MAX_LOG_AGE_SEC,
  "smartmoney_paper": {
    "reason": "$smart_reason",
    "running": $(bool "$smart_running"),
    "screen": $(bool "$smart_screen"),
    "log_age_sec": $smart_log_age,
    "risk_blocked": $(bool "$smart_blocked")
  },
  "whale_push": {
    "reason": "$whale_reason",
    "wrapper": $(bool "$whale_wrapper"),
    "child": $(bool "$whale_child"),
    "screen": $(bool "$whale_screen"),
    "log_age_sec": $whale_log_age
  }
}
EOF

printf 'health.result ts=%s overall=%s smart=%s smart_log_age=%s whale=%s whale_log_age=%s live_disabled=%s live_redeem_disabled=%s live_redeem_armed=%s research_disabled=%s monitoring_disabled=%s legacy_project_disabled=%s\n' \
  "$checked_at" "$overall" "$smart_reason" "$smart_log_age" "$whale_reason" "$whale_log_age" \
  "$live_trading_disabled" "$live_redeem_disabled" "$live_redeem_arm_present" "$research_disabled" "$monitoring_disabled" "$legacy_project_disabled"
