#!/usr/bin/env bash
# Commit and push the project's reviewable state once per day.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

REMOTE_URL="$(git remote get-url origin)"
EXPECTED_REMOTE="https://github.com/15529214579/polymarket-go.git"
if [ "$REMOTE_URL" != "$EXPECTED_REMOTE" ]; then
  echo "daily-git-save: unexpected origin remote: $REMOTE_URL" >&2
  exit 1
fi

DAY="$(date '+%Y-%m-%d')"

# Keep the daily save scoped to files that are useful in Git. Runtime state,
# logs, SQLite DBs, pid files, and compiled binaries are intentionally excluded.
git add \
  .gitignore .golangci.yml Makefile README.md SPEC.md TODO.md BTC_TODO.md PRINCIPLES.md \
  go.mod go.sum \
  archive docs cmd internal scripts reports \
  wallets*.txt

if git diff --cached --quiet; then
  echo "daily-git-save: no reviewable changes to commit"
else
  git commit -m "chore: daily git save ${DAY}"
fi

git push origin master
