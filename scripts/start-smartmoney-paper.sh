#!/usr/bin/env bash
# Start/stop the 5000U smart-money paper copytrade runner.
set -euo pipefail

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:$PATH"
export CLOB_PROXY="${CLOB_PROXY:-direct}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCRIPT="$ROOT/scripts/start-smartmoney-paper.sh"
STATE_DIR="$ROOT/db/smartmoney-paper"
PIDFILE="$STATE_DIR/pid"
CHILD_PIDFILE="$STATE_DIR/child.pid"
LOCKDIR="$STATE_DIR/lock"
WALLETS_FILE="$STATE_DIR/wallets.txt"
POLICY_FILE="$STATE_DIR/policy-start"
LOG="$ROOT/logs/smartmoney-paper.log"
ERR="$ROOT/logs/smartmoney-paper.err"
SCREEN_NAME="${SMARTMONEY_PAPER_SCREEN_NAME:-polymarket-smartmoney-paper}"

mkdir -p "$STATE_DIR" "$ROOT/logs" "$ROOT/bin"

INITIAL_CAPITAL="${SMARTMONEY_PAPER_INITIAL_CAPITAL:-5000}"
MIN_TIER="${SMARTMONEY_PAPER_MIN_TIER:-B}"
COPYTRADE_SIZE="${SMARTMONEY_PAPER_COPYTRADE_SIZE:-10}"
MAX_OPEN_USD="${SMARTMONEY_PAPER_MAX_OPEN_USD:-1500}"
MAX_PER_MARKET_USD="${SMARTMONEY_PAPER_MAX_PER_MARKET_USD:-100}"
MAX_OPEN_POSITIONS="${SMARTMONEY_PAPER_MAX_OPEN_POSITIONS:-120}"
MARKETS="${SMARTMONEY_PAPER_MARKETS:-120}"
WHALE_MIN_USD="${SMARTMONEY_PAPER_WHALE_MIN_USD:-500}"
WALLET_TIERS="${SMARTMONEY_PAPER_WALLET_TIERS:-$ROOT/db/strategy_iteration/copytrade_backtest_results.generated.json}"
SOURCE_WALLETS="${SMARTMONEY_PAPER_SOURCE_WALLETS:-$ROOT/wallets.strategy-push.txt $ROOT/wallets.hourly-push.txt $ROOT/wallets.leaderboard-watch.txt $ROOT/wallets.leaderboard-sports-push.txt}"
EXCLUDE_WALLETS="${SMARTMONEY_PAPER_EXCLUDE_WALLETS:-$ROOT/wallets.strategy-quarantine.txt $ROOT/wallets.strategy-review-noise.txt $ROOT/db/strategy_iteration/wallets.strategy-exclude.txt}"

acquire_lock() {
  while ! mkdir "$LOCKDIR" 2>/dev/null; do
    local lock_pid
    lock_pid="$(cat "$LOCKDIR/pid" 2>/dev/null || true)"
    if [ -n "$lock_pid" ] && kill -0 "$lock_pid" 2>/dev/null; then
      echo "smartmoney-paper.lock_held lock_pid=$lock_pid self=$$ action=exit"
      exit 0
    fi
    echo "smartmoney-paper.lock_stale lock_pid=${lock_pid:-} action=remove"
    rm -rf "$LOCKDIR"
  done
  echo $$ > "$LOCKDIR/pid"
  echo "smartmoney-paper.lock_acquired pid=$$ lock=$LOCKDIR"
}

release_lock() {
  local lock_pid
  lock_pid="$(cat "$LOCKDIR/pid" 2>/dev/null || true)"
  if [ "$lock_pid" = "$$" ]; then
    rm -rf "$LOCKDIR"
  fi
}

is_running() {
  [ -s "$PIDFILE" ] || return 1
  local pid
  pid="$(cat "$PIDFILE" 2>/dev/null || true)"
  [ -n "$pid" ] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  local cmd
  cmd="$(ps -o command= -p "$pid" 2>/dev/null || true)"
  case "$cmd" in
    *"start-smartmoney-paper.sh run"*) return 0 ;;
    *) return 1 ;;
  esac
}

merge_wallet_files() {
  local out="$1"
  local tmp="$out.tmp"
  local blocked_tmp="$out.blocked.tmp"
  : > "$tmp"
  : > "$blocked_tmp"
  for file in $EXCLUDE_WALLETS; do
    [ -s "$file" ] || continue
    awk '
      {
        line=$0
        sub(/#.*/, "", line)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
        addr=tolower(line)
        if (addr ~ /^0x[0-9a-f]{40}$/) print addr
      }
    ' "$file" >> "$blocked_tmp"
  done
  for file in $SOURCE_WALLETS; do
    [ -s "$file" ] || continue
    printf '# source=%s\n' "$file" >> "$tmp"
    awk 'NF { print }' "$file" >> "$tmp"
  done
  awk '
    FNR == NR {
      blocked[$1] = 1
      next
    }
    /^[[:space:]]*#/ { print; next }
    {
      addr=tolower($1)
      if (addr ~ /^0x[0-9a-f]{40}$/ && !blocked[addr] && !seen[addr]++) print
    }
  ' "$blocked_tmp" "$tmp" > "$out"
  rm -f "$tmp" "$blocked_tmp"
}

prepare() {
  cd "$ROOT"
  go build -o bin/bot ./cmd/bot
  merge_wallet_files "$WALLETS_FILE"
  WALLET_COUNT="$(awk '/^0x[0-9a-fA-F]{40}/ { n++ } END { print n+0 }' "$WALLETS_FILE")"
  if [ "$WALLET_COUNT" -eq 0 ]; then
    echo "no wallets available for smartmoney paper runner" >&2
    exit 1
  fi
  {
    printf 'started_at=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    printf 'initial_capital=%s\n' "$INITIAL_CAPITAL"
    printf 'min_tier=%s\n' "$MIN_TIER"
    printf 'copytrade_size=%s\n' "$COPYTRADE_SIZE"
    printf 'max_open_usd=%s\n' "$MAX_OPEN_USD"
    printf 'max_per_market_usd=%s\n' "$MAX_PER_MARKET_USD"
    printf 'max_open_positions=%s\n' "$MAX_OPEN_POSITIONS"
    printf 'markets=%s\n' "$MARKETS"
    printf 'whale_min_usd=%s\n' "$WHALE_MIN_USD"
    printf 'wallet_count=%s\n' "$WALLET_COUNT"
    printf 'wallets_file=%s\n' "$WALLETS_FILE"
    printf 'wallet_tiers=%s\n' "$WALLET_TIERS"
  } > "$POLICY_FILE"
}

run_loop() {
  acquire_lock
  prepare
  echo $$ > "$PIDFILE"
  echo "smartmoney-paper.loop pid=$$ policy=$POLICY_FILE"
  trap 'child="$(cat "$CHILD_PIDFILE" 2>/dev/null || true)"; echo "smartmoney-paper.loop_signal signal=term wrapper_pid=$$ child_pid=${child:-}"; [ -n "$child" ] && kill "$child" 2>/dev/null || true; rm -f "$PIDFILE" "$CHILD_PIDFILE"; release_lock; exit 0' INT TERM HUP
  restart_count=0
  while :; do
    restart_count=$((restart_count + 1))
    spawn_ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    spawn_epoch="$(date +%s)"
    echo "smartmoney-paper.spawn restart=$restart_count ts=$spawn_ts capital=$INITIAL_CAPITAL min_tier=$MIN_TIER wallets=$WALLET_COUNT markets=$MARKETS max_open_usd=$MAX_OPEN_USD"
    "$ROOT/bin/bot" \
      -mode=detect \
      -signal_mode=copytrade \
      -exit_mode=ladder \
      -markets="$MARKETS" \
      -window=60 \
      -fee_bp=0 \
      -lottery_enabled=false \
      -ladder_sl_pct=0.20 \
      -ladder_max_hold=10m \
      -whale_enabled \
      -whale_min_usd="$WHALE_MIN_USD" \
      -whale_interval=10s \
      -whale_replay_window=15m \
      -wallets_file="$WALLETS_FILE" \
      -copytrade_size="$COPYTRADE_SIZE" \
      -wallet_tiers="$WALLET_TIERS" \
      -min_tier="$MIN_TIER" \
      -initial_capital="$INITIAL_CAPITAL" \
      -positions_state="$STATE_DIR/positions.json" \
      -risk_state="$STATE_DIR/risk_state.json" \
      -buy_times_state="$STATE_DIR/buy_times.json" \
      -journal_dir="$STATE_DIR/journal" \
      -tickpath_dir="$STATE_DIR/tickpath" \
      -pos_max_total_open_usd="$MAX_OPEN_USD" \
      -pos_max_per_market_usd="$MAX_PER_MARKET_USD" \
      -pos_max_open_positions="$MAX_OPEN_POSITIONS" &
    child_pid=$!
    echo "$child_pid" > "$CHILD_PIDFILE"
    echo "smartmoney-paper.child_start restart=$restart_count child_pid=$child_pid"
    status=0
    wait "$child_pid" || status=$?
    rm -f "$CHILD_PIDFILE"
    end_epoch="$(date +%s)"
    duration_sec=$((end_epoch - spawn_epoch))
    exit_kind="exit"
    signal_num=""
    if [ "$status" -ge 128 ]; then
      exit_kind="signal"
      signal_num=$((status - 128))
    elif [ "$status" -ne 0 ]; then
      exit_kind="error"
    fi
    restart_delay="${SMARTMONEY_PAPER_RESTART_DELAY:-10}"
    echo "smartmoney-paper.child_exit restart=$restart_count child_pid=$child_pid status=$status exit_kind=$exit_kind signal=${signal_num:-} duration_sec=$duration_sec restart_delay=${restart_delay}s"
    sleep "$restart_delay"
  done
}

start() {
  if is_running; then
    echo "smartmoney-paper already running pid=$(cat "$PIDFILE")"
    return 0
  fi
  prepare

  if command -v screen >/dev/null 2>&1; then
    screen -S "$SCREEN_NAME" -X quit >/dev/null 2>&1 || true
    screen -dmS "$SCREEN_NAME" bash -lc "exec '$SCRIPT' run >> '$LOG' 2>> '$ERR'"
  else
    nohup "$SCRIPT" run >> "$LOG" 2>> "$ERR" &
    echo $! > "$PIDFILE"
  fi
  sleep 2
  if is_running; then
    echo "smartmoney-paper started pid=$(cat "$PIDFILE") wallets=$WALLET_COUNT policy=$POLICY_FILE"
  else
    rm -f "$PIDFILE"
    echo "smartmoney-paper failed to start; see $ERR" >&2
    exit 1
  fi
}

stop() {
  if ! is_running; then
    echo "smartmoney-paper not running"
    rm -f "$PIDFILE"
    return 0
  fi
  local pid
  pid="$(cat "$PIDFILE")"
  local child
  child="$(cat "$CHILD_PIDFILE" 2>/dev/null || true)"
  [ -n "$child" ] && kill -TERM "$child" 2>/dev/null || true
  kill -TERM "$pid" 2>/dev/null || true
  for _ in 1 2 3 4 5; do
    sleep 1
    is_running || break
  done
  if is_running; then
    kill -KILL "$pid" 2>/dev/null || true
  fi
  rm -f "$PIDFILE" "$CHILD_PIDFILE"
  echo "smartmoney-paper stopped"
}

status() {
  if is_running; then
    local pid
    pid="$(cat "$PIDFILE")"
    echo "RUNNING pid=$pid"
    ps -o pid,etime,rss,command -p "$pid" 2>/dev/null | tail -1 || true
    [ -s "$POLICY_FILE" ] && cat "$POLICY_FILE"
  else
    echo "NOT RUNNING"
    [ -s "$ERR" ] && tail -10 "$ERR"
  fi
}

case "${1:-status}" in
  start) start ;;
  stop) stop ;;
  restart) stop; start ;;
  status) status ;;
  tail) tail -f "$LOG" ;;
  run) run_loop ;;
  *) echo "usage: $0 {start|stop|restart|status|tail}" >&2; exit 2 ;;
esac
