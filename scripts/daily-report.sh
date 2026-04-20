#!/usr/bin/env bash
# daily-report.sh — 被 macOS crontab 在 SGT 00:00 触发，汇总「昨日」trade 并推 Telegram。
# 设计：
#   - 不依赖 bot 进程在跑（直接读 db/journal/trades-YYYY-MM-DD.jsonl）
#   - 失败也 exit 0，不让 cron 链断
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 0
mkdir -p "$ROOT/logs"

now_local=$(date '+%Y-%m-%d %H:%M:%S %Z')
day=$(date '+%Y-%m-%d')
log="$ROOT/logs/daily-report-${day}.log"

export PATH="/usr/local/go/bin:$PATH"

{
  echo "===== ${now_local} =====
"
  if [ ! -x "$ROOT/bin/bot" ]; then
    echo "bin/bot missing — running 'go build'"
    (cd "$ROOT" && go build -o bin/bot ./cmd/bot) || { echo "build failed"; exit 0; }
  fi
  "$ROOT/bin/bot" -mode=daily-report -report_push
  echo
} >> "$log" 2>&1

exit 0
