#!/usr/bin/env bash
# cron-poke.sh — 被 macOS crontab 周期调用。
# 作用：运行 heartbeat、追加日志、更新 state.json、必要时写告警标记。
# 本脚本不直接推 telegram — 告警投递交给 5号 session 醒来后读 state.json 再决定。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 0
. "$ROOT/scripts/component-flags.sh"
component_disabled "$ROOT" monitoring && exit 0
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

# Keep tracking ticks_no_progress for state visibility, but do not page on a
# stale commit alone. This repo can run unattended while paper-trading.

cat > "$ROOT/state.json" <<EOF
{
  "last_heartbeat": "${now_iso}",
  "last_commit": "${last_commit}",
  "uncommitted": ${uncommitted:-0},
  "open_todo": ${open_todo:-0},
  "build_fail": ${build_fail:-0},
  "ticks_no_progress": ${ticks_no_progress},
  "quiet_window": ${quiet},
  "alert": "${alert}",
  "legacy_paper_archived": true
}
EOF

# ── P10 日志异常自动扫描（每 20min cron-poke 触发） ──
anomaly_log="$ROOT/logs/anomaly-${day}.log"
daemon_log="$ROOT/logs/smartmoney-paper.log"
if [ -f "$daemon_log" ]; then
  cutoff=$(date -v-20M '+%Y-%m-%dT%H:%M' 2>/dev/null || date -d '20 minutes ago' '+%Y-%m-%dT%H:%M' 2>/dev/null || echo "")
  if [ -n "$cutoff" ]; then
    recent=$(awk -v t="$cutoff" '$0 ~ /"ts":"[0-9T:-]+"/ { if (index($0, t) > 0 || $0 > t) print }' "$daemon_log" | tail -200)

    err_count=$(echo "$recent" | grep -c '"_err\|_failed\|_error\|panic\|FATAL' || true)
    btc_leak=$(echo "$recent" | grep -c 'btc_strategy' || true)
    injury_alert=$(echo "$recent" | grep -c 'injury_alert' || true)
    injury_fetch=$(echo "$recent" | grep -c 'injury_fetch' || true)

    anomalies=""
    [ "$err_count" -gt 5 ] && anomalies="${anomalies}error_spike(${err_count}) "
    [ "$btc_leak" -gt 0 ] && anomalies="${anomalies}btc_strategy_leak(${btc_leak}) "
    if [ "$injury_fetch" -gt 0 ] && [ "$injury_alert" -eq 0 ]; then
      anomalies="${anomalies}injury_fetch_no_alert "
    fi

    if [ -n "$anomalies" ]; then
      echo "[${now_local}] ANOMALY: ${anomalies}" >> "$anomaly_log"
      if [ "${quiet:-0}" = "0" ] && [ -z "$alert" ]; then
        alert="log-anomaly: ${anomalies}"
      fi
    fi
  fi
fi

# 告警投递（读 state.json，cooldown 2h，夜间静默）
"$ROOT/scripts/alert-dispatch.sh" >> "$log" 2>&1 || true

exit 0
