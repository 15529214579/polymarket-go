#!/usr/bin/env bash
# alert-dispatch.sh — 读 state.json，如有 alert 则推送 Telegram 告警给老板。
# 设计要点：
#   - 夜间静默（state.json.quiet_window=1 时）不推送
#   - cooldown：同一 alert 标签 2h 内只推一次，写入 logs/alert-sent.json
#   - token/chat 从 .env.local 读取（gitignored）
#   - curl 失败也 exit 0，绝不让 cron 链断
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
STATE="$ROOT/state.json"
SENT="$ROOT/logs/alert-sent.json"
ENV_FILE="$ROOT/.env.local"

[ -f "$STATE" ] || exit 0
[ -f "$ENV_FILE" ] || { echo "alert-dispatch: .env.local missing" >&2; exit 0; }

# shellcheck disable=SC1090
set -a; . "$ENV_FILE"; set +a
[ -n "${TELEGRAM_BOT_TOKEN:-}" ] || exit 0
[ -n "${TELEGRAM_CHAT_ID:-}" ] || exit 0

alert=$(sed -n 's/.*"alert": *"\([^"]*\)".*/\1/p' "$STATE" | head -1)
quiet=$(sed -n 's/.*"quiet_window": *\([0-9]*\).*/\1/p' "$STATE" | head -1)

[ -z "$alert" ] && exit 0
[ "${quiet:-0}" = "1" ] && exit 0

# cooldown 2h
mkdir -p "$ROOT/logs"
now_ts=$(date +%s)
last_ts=0
last_tag=""
if [ -f "$SENT" ]; then
  last_tag=$(sed -n 's/.*"tag": *"\([^"]*\)".*/\1/p' "$SENT" | head -1)
  last_ts=$(sed -n 's/.*"ts": *\([0-9]*\).*/\1/p' "$SENT" | head -1)
  [ -z "$last_ts" ] && last_ts=0
fi
cooldown=7200
if [ "$alert" = "$last_tag" ] && [ $((now_ts - last_ts)) -lt "$cooldown" ]; then
  exit 0
fi

# 组装消息
case "$alert" in
  build-failing) title="🔴 polymarket-go build 失败" ;;
  stalled-12h)   title="🟡 polymarket-go 停滞 12h+ 未 commit" ;;
  *)             title="⚠️ polymarket-go: $alert" ;;
esac

last_commit=$(sed -n 's/.*"last_commit": *"\([^"]*\)".*/\1/p' "$STATE" | head -1)
open_todo=$(sed -n 's/.*"open_todo": *\([0-9]*\).*/\1/p' "$STATE" | head -1)
body=$(cat <<MSG
${title}

last commit: ${last_commit:-<none>}
open todo:   ${open_todo:-?}
alert tag:   ${alert}

—— 5号 自动告警（alert-dispatch.sh）
MSG
)

# 发送
resp=$(curl -s --max-time 10 -X POST \
  "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
  --data-urlencode "chat_id=${TELEGRAM_CHAT_ID}" \
  --data-urlencode "text=${body}" \
  --data-urlencode "disable_notification=false" 2>&1 || true)

ok=$(LC_ALL=C echo "$resp" | LC_ALL=C grep -oE '"ok":[[:space:]]*(true|false)' | head -1 | awk -F: '{gsub(/ /,"",$2); print $2}')
[ -z "$ok" ] && ok="unknown"

cat > "$SENT" <<EOF
{
  "tag": "${alert}",
  "ts": ${now_ts},
  "ok": "${ok:-unknown}",
  "sent_at": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
}
EOF

exit 0
