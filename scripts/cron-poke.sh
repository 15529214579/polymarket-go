#!/usr/bin/env bash
# cron-poke.sh — 被 macOS crontab 周期调用。
# 作用：运行 heartbeat、追加日志、更新 state.json、必要时写告警标记。
# 本脚本不直接推 telegram — 告警投递交给 5号 session 醒来后读 state.json 再决定。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 0
mkdir -p "$ROOT/logs"

now_iso=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
now_local=$(date '+%Y-%m-%d %H:%M:%S %Z')
day=$(date '+%Y-%m-%d')
hour=$(date '+%H')
log="$ROOT/logs/heartbeat-${day}.log"

# 夜间静默（SGT 00:00-07:59）：只记录，不触发任何告警通道
quiet=0
if [ "$hour" -lt 8 ]; then quiet=1; fi

{
  echo "----- ${now_local} (quiet=${quiet}) -----"
  "$ROOT/scripts/heartbeat.sh"
  echo
} >> "$log" 2>&1

# 提炼状态（从刚跑完的 heartbeat 重跑一次拿原始值 — 比解析日志稳）
out=$("$ROOT/scripts/heartbeat.sh" 2>/dev/null)
build_fail=$(echo "$out" | grep -c 'go build: FAIL' || true)
uncommitted=$(echo "$out" | sed -n 's/.*uncommitted: \([0-9]*\).*/\1/p' | head -1)
open_todo=$(echo "$out" | sed -n 's/.*open: \([0-9]*\).*/\1/p' | head -1)
last_commit=$(echo "$out" | grep '  last:' | head -1 | sed 's/^  last: //')

# alert 逻辑：
#  - go build 失败
#  - uncommitted > 20（说明工作停在半途）
#  - 连续多次无 commit（通过 state.json 跨跳追踪）
alert=""
[ "${build_fail:-0}" -gt 0 ] && alert="build-failing"

# 读上次 commit，对比是否长时间未推进
prev_commit=""
prev_ticks_no_progress=0
if [ -f "$ROOT/state.json" ]; then
  prev_commit=$(sed -n 's/.*"last_commit": *"\([^"]*\)".*/\1/p' "$ROOT/state.json" | head -1)
  prev_ticks_no_progress=$(sed -n 's/.*"ticks_no_progress": *\([0-9]*\).*/\1/p' "$ROOT/state.json" | head -1)
  [ -z "$prev_ticks_no_progress" ] && prev_ticks_no_progress=0
fi

if [ "$last_commit" = "$prev_commit" ] && [ -n "$last_commit" ]; then
  ticks_no_progress=$((prev_ticks_no_progress + 1))
else
  ticks_no_progress=0
fi

# 超过 24 个 tick（30min × 24 = 12h）没 commit 且 quiet 结束后，告警
if [ "$ticks_no_progress" -gt 24 ] && [ "$quiet" = "0" ] && [ -z "$alert" ]; then
  alert="stalled-12h"
fi

cat > "$ROOT/state.json" <<EOF
{
  "last_heartbeat": "${now_iso}",
  "last_commit": "${last_commit}",
  "uncommitted": ${uncommitted:-0},
  "open_todo": ${open_todo:-0},
  "build_fail": ${build_fail:-0},
  "ticks_no_progress": ${ticks_no_progress},
  "quiet_window": ${quiet},
  "alert": "${alert}"
}
EOF

# 告警投递（读 state.json，cooldown 2h，夜间静默）
"$ROOT/scripts/alert-dispatch.sh" >> "$log" 2>&1 || true

exit 0
