#!/usr/bin/env bash
# Start the smart-money whale order push worker in the current process.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LOG="$ROOT/db/agent.log"
ERR="$ROOT/db/agent.err"
PIDFILE="$ROOT/db/bot.pid"
WRAPPER_PIDFILE="$ROOT/db/whale-push.pid"
CHILD_PIDFILE="$ROOT/db/whale-push.child.pid"
LOCKDIR="$ROOT/db/whale-push.lock"
POLICY_START_FILE="${WHALE_POLICY_START_FILE:-$ROOT/db/whale-push.policy-start}"
WHALE_BASE_WALLETS_FILE="${WHALE_WALLETS_FILE:-$ROOT/wallets.strategy-push.txt}"
WHALE_EXTRA_WALLETS_FILES="${WHALE_EXTRA_WALLETS_FILES:-$ROOT/wallets.football-score-push.txt $ROOT/wallets.leaderboard-push.txt $ROOT/wallets.leaderboard-watch.txt $ROOT/wallets.leaderboard-sports-push.txt $ROOT/wallets.sports-holders-push.txt $ROOT/wallets.hourly-push.txt}"
WHALE_EXCLUDE_WALLETS_FILES="${WHALE_EXCLUDE_WALLETS_FILES:-$ROOT/wallets.strategy-quarantine.txt $ROOT/wallets.strategy-review-noise.txt $ROOT/db/strategy_iteration/wallets.strategy-exclude.txt}"
WHALE_RUNTIME_WALLETS_FILE="${WHALE_RUNTIME_WALLETS_FILE:-$ROOT/db/whale-push.wallets.txt}"
WHALE_WALLETS_FILE="$WHALE_BASE_WALLETS_FILE"

mkdir -p "$ROOT/db"
cd "$ROOT"
exec >> "$LOG" 2>> "$ERR"

acquire_lock() {
  while ! mkdir "$LOCKDIR" 2>/dev/null; do
    local lock_pid
    lock_pid="$(cat "$LOCKDIR/pid" 2>/dev/null || true)"
    if [ -n "$lock_pid" ] && kill -0 "$lock_pid" 2>/dev/null; then
      echo "whale-push.lock_held lock_pid=$lock_pid self=$$ action=exit"
      exit 0
    fi
    echo "whale-push.lock_stale lock_pid=${lock_pid:-} action=remove"
    rm -rf "$LOCKDIR"
  done
  echo $$ > "$LOCKDIR/pid"
  echo "whale-push.lock_acquired pid=$$ lock=$LOCKDIR"
}

release_lock() {
  local lock_pid
  lock_pid="$(cat "$LOCKDIR/pid" 2>/dev/null || true)"
  if [ "$lock_pid" = "$$" ]; then
    rm -rf "$LOCKDIR"
  fi
}

merge_wallet_files() {
  local out="$1"
  local tmp="$out.tmp"
  local blocked_tmp="$out.blocked.tmp"
  : > "$tmp"
  : > "$blocked_tmp"
  for file in $WHALE_EXCLUDE_WALLETS_FILES; do
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
  for file in "$WHALE_BASE_WALLETS_FILE" $WHALE_EXTRA_WALLETS_FILES; do
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

merge_wallet_files "$WHALE_RUNTIME_WALLETS_FILE"
if [ -s "$WHALE_RUNTIME_WALLETS_FILE" ]; then
  WHALE_WALLETS_FILE="$WHALE_RUNTIME_WALLETS_FILE"
elif [ -s "$ROOT/wallets.strategy-core.txt" ]; then
  WHALE_WALLETS_FILE="$ROOT/wallets.strategy-core.txt"
fi

pid_command() {
  ps -o command= -p "$1" 2>/dev/null || true
}

is_whale_push_pid() {
  local pid="$1"
  local cmd
  cmd="$(pid_command "$pid")"
  case "$cmd" in
    *start-whale-push.sh*|*"/bin/bot -mode=detect -signal_mode=whale"*|*"bin/bot -mode=detect -signal_mode=whale"*)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

stop_pidfile() {
  local file="$1"
  local pid
  pid="$(cat "$file" 2>/dev/null || true)"
  [ -n "$pid" ] || return 0
  [ "$pid" != "$$" ] || return 0
  kill -0 "$pid" 2>/dev/null || {
    rm -f "$file"
    return 0
  }
  is_whale_push_pid "$pid" || return 0
  kill "$pid" 2>/dev/null || true
  sleep 1
  kill -0 "$pid" 2>/dev/null && kill -KILL "$pid" 2>/dev/null || true
  rm -f "$file"
}

stop_pidfile "$CHILD_PIDFILE"
stop_pidfile "$WRAPPER_PIDFILE"
stop_pidfile "$PIDFILE"
acquire_lock

export RESTART_REASON="${RESTART_REASON:-whale-push}"
export CLOB_PROXY="${CLOB_PROXY:-direct}"
export WHALE_CONFIRM_LISTS="${WHALE_CONFIRM_LISTS:-watch,scout,target,tape,sports,leaderboard_watch}"
export WHALE_CONFIRM_WINDOW="${WHALE_CONFIRM_WINDOW:-30m}"
export WHALE_CONFIRM_MIN_WALLETS="${WHALE_CONFIRM_MIN_WALLETS:-2}"
export WHALE_CONFIRM_BYPASS_USD="${WHALE_CONFIRM_BYPASS_USD:-3000}"
export WHALE_CONFIRM_MAX_WORSE_PRICE="${WHALE_CONFIRM_MAX_WORSE_PRICE:-0.02}"
export WHALE_LIST_MIN_USD="${WHALE_LIST_MIN_USD:-core=750,sports=1000,watch=500,scout=500,target=500,tape=1000,flow=1000,football_score_push=300,leaderboard_push=3000,leaderboard_watch=3000,leaderboard_sports_push=1000,sports_holders_push=1000}"

policy_started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
{
  printf 'started_at=%s\n' "$policy_started_at"
  printf 'wallets_file=%s\n' "$WHALE_WALLETS_FILE"
  printf 'base_wallets_file=%s\n' "$WHALE_BASE_WALLETS_FILE"
  printf 'extra_wallets_files=%s\n' "$WHALE_EXTRA_WALLETS_FILES"
  printf 'exclude_wallets_files=%s\n' "$WHALE_EXCLUDE_WALLETS_FILES"
  printf 'whale_min_usd=%s\n' "${WHALE_MIN_USD:-500}"
  printf 'whale_markets=%s\n' "${WHALE_MARKETS:-120}"
  printf 'whale_list_min_usd=%s\n' "$WHALE_LIST_MIN_USD"
  printf 'whale_confirm_lists=%s\n' "$WHALE_CONFIRM_LISTS"
  printf 'whale_confirm_min_wallets=%s\n' "$WHALE_CONFIRM_MIN_WALLETS"
  printf 'whale_confirm_bypass_usd=%s\n' "$WHALE_CONFIRM_BYPASS_USD"
} > "$POLICY_START_FILE"

echo $$ > "$PIDFILE"
echo $$ > "$WRAPPER_PIDFILE"

echo "whale-push.start pid=$$ wallets=$WHALE_WALLETS_FILE policy_started_at=$policy_started_at policy_file=$POLICY_START_FILE"

child_pid=""
terminate() {
  echo "whale-push.loop_signal signal=term wrapper_pid=$$ child_pid=${child_pid:-}"
  if [ -n "$child_pid" ]; then
    kill "$child_pid" 2>/dev/null || true
    wait "$child_pid" 2>/dev/null || true
  fi
  rm -f "$CHILD_PIDFILE" "$WRAPPER_PIDFILE"
  release_lock
  exit 0
}
trap terminate INT TERM HUP

restart_count=0
while :; do
  restart_count=$((restart_count + 1))
  spawn_ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  spawn_epoch="$(date +%s)"
  echo "whale-push.spawn restart=$restart_count ts=$spawn_ts wallets=$WHALE_WALLETS_FILE min_usd=${WHALE_MIN_USD:-500} markets=${WHALE_MARKETS:-120} min_tier=A"
  "$ROOT/bin/bot" \
    -mode=detect \
    -signal_mode=whale \
    -exit_mode=ladder \
    -markets="${WHALE_MARKETS:-120}" \
    -window=60 \
    -fee_bp=0 \
    -lottery_enabled=false \
    -ladder_sl_pct=0.20 \
    -ladder_max_hold=10m \
    -whale_enabled \
    -whale_min_usd="${WHALE_MIN_USD:-500}" \
    -whale_interval="${WHALE_INTERVAL:-10s}" \
    -whale_replay_window="${WHALE_REPLAY_WINDOW:-15m}" \
    -wallets_file="$WHALE_WALLETS_FILE" \
    -copytrade_size=0 \
    -wallet_tiers="$ROOT/db/user_wallet_review/copytrade_backtest_results.generated.json" \
    -min_tier=A &
  child_pid=$!
  echo "whale-push.child_start restart=$restart_count child_pid=$child_pid"
  echo "$child_pid" > "$CHILD_PIDFILE"
  status=0
  wait "$child_pid" || status=$?
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
  child_pid=""
  rm -f "$CHILD_PIDFILE"
  restart_delay="${WHALE_RESTART_DELAY:-10}"
  echo "whale-push.child_exit restart=$restart_count status=$status exit_kind=$exit_kind signal=${signal_num:-} duration_sec=$duration_sec restart_delay=${restart_delay}s"
  sleep "$restart_delay"
done
