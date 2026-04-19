#!/usr/bin/env bash
# heartbeat.sh — 单次项目自检，输出给 5号 session 用。
# 原则：只读，绝不改文件；失败也退 0，让 cron 继续下一跳。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 0

now=$(date '+%Y-%m-%d %H:%M:%S %Z')
echo "=== polymarket-go heartbeat @ ${now} ==="
echo

echo "[git]"
git -C "$ROOT" log -1 --format='  last: %h %cr — %s' 2>/dev/null || echo "  (no commits)"
dirty=$(git -C "$ROOT" status --porcelain 2>/dev/null | wc -l | tr -d ' ')
echo "  uncommitted: ${dirty} files"
branch=$(git -C "$ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null || echo "?")
remote=$(git -C "$ROOT" remote -v 2>/dev/null | awk 'NR==1{print $2}')
echo "  branch: ${branch}   remote: ${remote:-<none>}"

echo
echo "[build]"
export PATH="/usr/local/go/bin:$PATH"
if command -v go >/dev/null 2>&1; then
  if go build ./... 2>/tmp/polymarket-go-build.err; then
    echo "  go build: OK"
  else
    echo "  go build: FAIL"
    sed 's/^/    /' /tmp/polymarket-go-build.err | head -20
  fi
else
  echo "  go: not found"
fi

echo
echo "[todo]"
if [ -f "$ROOT/TODO.md" ]; then
  open_count=$(grep -c '^\- \[ \]' "$ROOT/TODO.md" 2>/dev/null || echo 0)
  done_count=$(grep -c '^\- \[x\]' "$ROOT/TODO.md" 2>/dev/null || echo 0)
  echo "  open: ${open_count}   done: ${done_count}"
  echo "  next 3 open:"
  grep -n '^\- \[ \]' "$ROOT/TODO.md" | head -3 | sed 's/^/    /'
else
  echo "  (no TODO.md)"
fi

echo
echo "[logs]"
if [ -d "$ROOT/logs" ] && [ "$(ls -A "$ROOT/logs" 2>/dev/null)" ]; then
  latest=$(ls -t "$ROOT/logs" | head -1)
  mtime=$(stat -f '%Sm' -t '%Y-%m-%d %H:%M:%S' "$ROOT/logs/$latest" 2>/dev/null)
  echo "  latest: ${latest} (${mtime})"
else
  echo "  (no logs yet)"
fi

echo
echo "[state]"
if [ -f "$ROOT/state.json" ]; then
  cat "$ROOT/state.json" | sed 's/^/  /'
else
  echo "  (no state.json)"
fi

echo
echo "=== end heartbeat ==="
