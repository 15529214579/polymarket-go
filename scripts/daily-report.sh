#!/usr/bin/env bash
# daily-report.sh — 被 macOS crontab 在 SGT 00:00 触发，汇总「昨日」trade 并推 Telegram。
# 设计：不依赖 bot 进程在跑，直接读取活跃 smartmoney 模拟盘的平仓日志。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 1
mkdir -p "$ROOT/logs"

now_local=$(date '+%Y-%m-%d %H:%M:%S %Z')
day=$(date '+%Y-%m-%d')
log="$ROOT/logs/daily-report-${day}.log"
JOURNAL_DIR="${DAILY_REPORT_JOURNAL_DIR:-${SMARTMONEY_PAPER_JOURNAL:-$ROOT/db/smartmoney-paper/journal}}"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

{
  echo "===== ${now_local} =====
"
  if [ ! -d "$JOURNAL_DIR" ]; then
    echo "active paper journal missing: $JOURNAL_DIR"
    exit 1
  fi
  if [ ! -x "$ROOT/bin/bot" ]; then
    echo "bin/bot missing — running 'go build'"
    (cd "$ROOT" && go build -o bin/bot ./cmd/bot) || { echo "build failed"; exit 1; }
  fi
  if ! "$ROOT/bin/bot" -mode=daily-report -journal_dir="$JOURNAL_DIR" -report_push; then
    echo "daily-report failed"
    exit 1
  fi
  echo
} >> "$log" 2>&1

exit 0
