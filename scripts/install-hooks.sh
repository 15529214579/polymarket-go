#!/usr/bin/env bash
# install-hooks.sh — copy tracked git hooks into .git/hooks/ (idempotent).
set -eu
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
mkdir -p "$ROOT/.git/hooks"

cat > "$ROOT/.git/hooks/post-commit" <<'HOOK'
#!/usr/bin/env bash
# post-commit — refresh state.json immediately after a commit so heartbeat
# readers don't see up-to-20min-stale last_commit. Fire-and-forget, never
# blocks the commit.
set -u
ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
[ -x "$ROOT/scripts/cron-poke.sh" ] || exit 0
nohup "$ROOT/scripts/cron-poke.sh" >/dev/null 2>&1 &
disown 2>/dev/null || true
exit 0
HOOK
chmod +x "$ROOT/.git/hooks/post-commit"
echo "installed: .git/hooks/post-commit"
