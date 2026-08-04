#!/usr/bin/env bash
# bot-daemon.sh — start | stop | status | tail | restart for the paper detect bot.
# This is the legacy 5U paper runner. Keep its process state and logs separate
# from the whale-push worker, which has its own lifecycle.
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PIDFILE="$ROOT/db/legacy-paper.pid"
LOG="$ROOT/logs/legacy-paper.log"
ERR="$ROOT/logs/legacy-paper.err"
mkdir -p "$ROOT/db" "$ROOT/logs"

action="${1:-status}"

is_running() {
  [ -f "$PIDFILE" ] || return 1
  local pid
  pid=$(cat "$PIDFILE" 2>/dev/null || echo "")
  [ -n "$pid" ] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  # Verify PID is actually our bot binary (not a reused PID)
  local cmd
  cmd=$(ps -o comm= -p "$pid" 2>/dev/null || echo "")
  case "$cmd" in
    */bot|bot) return 0 ;;
    *) return 1 ;;
  esac
}

start() {
  if is_running; then
    echo "already running pid=$(cat "$PIDFILE")"
    return 0
  fi
  export PATH="/usr/local/go/bin:$PATH"
  export RESTART_REASON="${RESTART_REASON:-manual}"
  ( cd "$ROOT" && go build -o bin/bot ./cmd/bot ) || { echo "build failed"; exit 1; }
  cd "$ROOT" || exit 1
  shift_args=("${@:2}")
  # Default mode (2026-05-30 SGT): copytrade mode, push/follow only the
  # manually selected wallet. Tiered sizing: A=$20, B=$10, C/D=$5.
  # Min trade $100. Poll 60s. Ladder exit SL 20% / timeout 10m.
  # Live mode is never inherited from the environment. It requires an explicit
  # custom -live argument plus the short-lived arm file checked by the bot.
  args=(-mode=detect -signal_mode=copytrade -exit_mode=ladder -markets=20 -window=60 -fee_bp=0 -ladder_sl_pct=0.20 -ladder_max_hold=10m -injury_enabled -injury_interval=1m -whale_enabled -whale_min_usd=100 -whale_interval=10s -wallets_file="$ROOT/wallets.push-only.txt" -copytrade_size=5 -wallet_tiers="$ROOT/db/copytrade_tiers_push_only.json" -min_tier=A)
  if [ "${#shift_args[@]}" -gt 0 ]; then
    args=("${shift_args[@]}")
  fi
  nohup "$ROOT/bin/bot" "${args[@]}" >> "$LOG" 2>> "$ERR" &
  echo $! > "$PIDFILE"
  sleep 1
  if is_running; then
    echo "started pid=$(cat "$PIDFILE")  args=${args[*]}"
  else
    echo "FAILED to start — see $ERR"
    rm -f "$PIDFILE"
    exit 1
  fi
}

stop() {
  if ! is_running; then
    echo "not running"
    rm -f "$PIDFILE"
    return 0
  fi
  local pid; pid=$(cat "$PIDFILE")
  echo "stopping pid=$pid"
  kill -TERM "$pid" 2>/dev/null || true
  for _ in 1 2 3 4 5; do
    sleep 1
    is_running || break
  done
  if is_running; then
    echo "still alive after TERM, sending KILL"
    kill -KILL "$pid" 2>/dev/null || true
  fi
  rm -f "$PIDFILE"
  echo "stopped"
}

status() {
  if is_running; then
    local pid; pid=$(cat "$PIDFILE")
    echo "RUNNING pid=$pid"
    ps -o pid,etime,rss,command -p "$pid" 2>/dev/null | tail -1
    echo "log: $LOG  ($(wc -l < "$LOG" 2>/dev/null || echo 0) lines)"
  else
    echo "NOT RUNNING"
    [ -f "$ERR" ] && echo "last stderr:" && tail -5 "$ERR"
  fi
}

case "$action" in
  start) start "$@" ;;
  stop) stop ;;
  restart) stop; start "$@" ;;
  status) status ;;
  tail) tail -f "$LOG" ;;
  *) echo "usage: $0 {start|stop|restart|status|tail}"; exit 2 ;;
esac
