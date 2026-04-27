#!/usr/bin/env bash
# daily-iterate.sh — SGT 00:05 cron 触发，分析最近 7 天交易数据，生成迭代报告 + 推 Telegram。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 0
mkdir -p "$ROOT/logs"

now_local=$(date '+%Y-%m-%d %H:%M:%S %Z')
day=$(date '+%Y-%m-%d')
log="$ROOT/logs/daily-iterate-${day}.log"

export PATH="/usr/local/go/bin:$PATH"

{
  echo "===== ${now_local} ====="
  if [ ! -x "$ROOT/bin/bot" ]; then
    echo "bin/bot missing — running 'go build'"
    (cd "$ROOT" && go build -o bin/bot ./cmd/bot) || { echo "build failed"; exit 0; }
  fi
  "$ROOT/bin/bot" -mode=daily-iterate -iterate_window=7 -report_push
  echo
} >> "$log" 2>&1

exit 0
