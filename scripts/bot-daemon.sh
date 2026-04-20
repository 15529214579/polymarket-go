#!/usr/bin/env bash
# bot-daemon.sh — start | stop | status | tail | restart for the paper detect bot.
# Logs go to db/agent.log (stdout) + db/agent.err (stderr).
# Single instance; PID file at db/bot.pid.
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PIDFILE="$ROOT/db/bot.pid"
LOG="$ROOT/db/agent.log"
ERR="$ROOT/db/agent.err"
mkdir -p "$ROOT/db"

action="${1:-status}"

is_running() {
  [ -f "$PIDFILE" ] || return 1
  local pid
  pid=$(cat "$PIDFILE" 2>/dev/null || echo "")
  [ -n "$pid" ] || return 1
  kill -0 "$pid" 2>/dev/null
}

start() {
  if is_running; then
    echo "already running pid=$(cat "$PIDFILE")"
    return 0
  fi
  export PATH="/usr/local/go/bin:$PATH"
  ( cd "$ROOT" && go build -o bin/bot ./cmd/bot ) || { echo "build failed"; exit 1; }
  cd "$ROOT" || exit 1
  shift_args=("${@:2}")
  # Default mode: prompt-only (R3: no auto-open; boss hand-picks via DM buttons),
  # 20 markets, 60s window, hold-to-settlement. Override by passing extra args.
  args=(-mode=detect -signal_mode=prompt -exit_mode=hold -markets=20 -window=60)
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
