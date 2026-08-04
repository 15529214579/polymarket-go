#!/usr/bin/env bash
# daily-iterate.sh — SGT 00:05 cron 触发，分析最近 7 天交易数据，生成迭代报告 + 推 Telegram。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 1
. "$ROOT/scripts/component-flags.sh"
component_disabled "$ROOT" research && exit 0
mkdir -p "$ROOT/logs"

LOCK_DIR="$ROOT/db/daily-iterate.lock"
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
	lock_pid="$(cat "$LOCK_DIR/pid" 2>/dev/null || true)"
	if [ -n "$lock_pid" ] && kill -0 "$lock_pid" 2>/dev/null; then
		lock_command="$(ps -o command= -p "$lock_pid" 2>/dev/null || true)"
		case "$lock_command" in
			*daily-iterate.sh*|"") exit 0 ;;
		esac
	fi
	rm -f "$LOCK_DIR/pid"
	rmdir "$LOCK_DIR" 2>/dev/null || exit 0
	mkdir "$LOCK_DIR" 2>/dev/null || exit 0
fi
printf '%s\n' "$$" > "$LOCK_DIR/pid"
trap 'rm -f "$LOCK_DIR/pid"; rmdir "$LOCK_DIR" 2>/dev/null || true' EXIT

now_local=$(date '+%Y-%m-%d %H:%M:%S %Z')
day=$(date '+%Y-%m-%d')
log="$ROOT/logs/daily-iterate-${day}.log"
JOURNAL_DIR="${DAILY_ITERATE_JOURNAL_DIR:-${SMARTMONEY_PAPER_JOURNAL:-$ROOT/db/smartmoney-paper/journal}}"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

{
  echo "===== ${now_local} ====="
  if [ ! -d "$JOURNAL_DIR" ]; then
    echo "active paper journal missing: $JOURNAL_DIR"
    exit 1
  fi
  if [ ! -x "$ROOT/bin/bot" ]; then
    echo "bin/bot missing — running 'go build'"
    (cd "$ROOT" && go build -o bin/bot ./cmd/bot) || { echo "build failed"; exit 1; }
  fi
  if ! "$ROOT/bin/bot" -mode=daily-iterate -journal_dir="$JOURNAL_DIR" -iterate_window=7 -report_push; then
    echo "daily-iterate failed"
    exit 1
  fi
  echo
} >> "$log" 2>&1

exit 0
