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
  ( cd "$ROOT" && go build -o bin/bot ./cmd/bot ) || { echo "build failed"; exit 1; }
  cd "$ROOT" || exit 1
  shift_args=("${@:2}")
  # Default mode (2026-04-25 SGT): whale-follow mode — momentum auto-open
  # and DM buttons disabled; whale BUY → SignalPrompt with buy buttons,
  # whale SELL → close prompt with buttons. Lottery scanner still runs.
  # OddsPapi: Pinnacle sharp-line scanner for football (EPL/UCL/La Liga) at 3h.
  # fee_bp=0 matches CLOB V1; update after V2 cutover.
  args=(-mode=detect -signal_mode=whale -exit_mode=ladder -markets=20 -window=60 -fee_bp=0 -injury_enabled -injury_interval=15m -whale_enabled "-whale_wallets=0xdb27bf2ac5d428a9c63dbc914611036855a6c56e|drpufferfish|1000|https://polymarket.com/@drpufferfish,0xbddf61af533ff524d27154e589d2d7a81510c684|countryside|1500|https://polymarket.com/@countryside,0xc2e7800b5af46e6093872b177b7a5e7f0563be51|beachboy4|5000|https://polymarket.com/@beachboy4" -oddspapi_enabled -oddspapi_interval=3h -oddspapi_bookmaker=pinnacle -oddspapi_sports=soccer_epl,soccer_spain_la_liga,soccer_uefa_champs_league -updown_enabled -updown_interval=10m -updown_confidence=0.40 -updown_size=5 -updown_max_daily=40 -updown_db=db/btc.db)
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
