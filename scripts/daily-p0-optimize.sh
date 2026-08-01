#!/bin/bash
set -euo pipefail

ROOT="${POLYMARKET_GO_ROOT:-/Users/murphyma/work/polymarket-go}"
CODEX_BIN="${CODEX_BIN:-/Applications/ChatGPT.app/Contents/Resources/codex}"
PROMPT_FILE="$ROOT/docs/daily-p0-optimizer-prompt.md"
STATE_DIR="$ROOT/db/self-optimization"
LOCK_DIR="$STATE_DIR/run.lock"
RUN_ID="$(date '+%Y%m%d-%H%M%S')"
DAY="$(date '+%Y-%m-%d')"
STARTED_AT="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
RUN_LOG="$ROOT/logs/daily-p0-optimize-${DAY}.log"
LAST_MESSAGE="$STATE_DIR/last-message.md"
STATE_FILE="$STATE_DIR/state.json"
LAST_RUN_FILE="$STATE_DIR/last-run-day"
CHANGED_FILE="$STATE_DIR/changed-${RUN_ID}.txt"
WORKTREE_DIR="${TMPDIR:-/private/tmp}/polymarket-go-p0-${RUN_ID}"
WORKTREE_BRANCH="automation/p0-${RUN_ID}"
MAX_FILES="${SELF_OPT_MAX_FILES:-12}"
MAX_DIFF_LINES="${SELF_OPT_MAX_DIFF_LINES:-600}"
CODEX_TIMEOUT_SECONDS="${SELF_OPT_CODEX_TIMEOUT_SECONDS:-2700}"
EXPECTED_REMOTE="${SELF_OPT_EXPECTED_REMOTE:-https://github.com/15529214579/polymarket-go.git}"
SUCCESS=0
AUTOMATION_COMMITTED=0
AUTOMATION_COMMIT=""
AUTOMATION_MERGED=0
WORKTREE_CREATED=0
LOCK_OWNED=0

SOURCE_GUARD_PATHS=(
  .gitignore .golangci.yml Makefile README.md SPEC.md TODO.md BTC_TODO.md PRINCIPLES.md
  go.mod go.sum archive docs cmd internal scripts
)
AUTONOMOUS_PATHS=(cmd internal scripts)

mkdir -p "$STATE_DIR" "$ROOT/logs"

log() {
  printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*" | tee -a "$RUN_LOG"
}

write_state() {
  local status="$1"
  local commit="${2:-}"
  local tmp="$STATE_FILE.tmp"
  printf '{\n  "status": "%s",\n  "run_id": "%s",\n  "started_at": "%s",\n  "finished_at": "%s",\n  "commit": "%s"\n}\n' \
    "$status" "$RUN_ID" "$STARTED_AT" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$commit" > "$tmp"
  mv "$tmp" "$STATE_FILE"
}

is_allowed_path() {
  case "$1" in
    cmd/*|internal/*|scripts/*)
      return 0
      ;;
  esac
  return 1
}

quarantine_worktree_changes() {
  [ "$WORKTREE_CREATED" -eq 1 ] || return 0
  if [ -n "$(git -C "$WORKTREE_DIR" status --porcelain)" ]; then
    git -C "$WORKTREE_DIR" stash push -u -m "daily-p0-failed-${RUN_ID}" >> "$RUN_LOG" 2>&1 || true
    log "Failed agent changes were quarantined in a Git stash."
  fi
}

remove_worktree() {
  [ "$WORKTREE_CREATED" -eq 1 ] || return 0
  cd "$ROOT"
  git -C "$ROOT" worktree remove "$WORKTREE_DIR" >> "$RUN_LOG" 2>&1
  WORKTREE_CREATED=0
  git -C "$ROOT" branch -d "$WORKTREE_BRANCH" >> "$RUN_LOG" 2>&1
}

cleanup() {
  local code=$?
  if [ "$code" -ne 0 ] && [ "$SUCCESS" -ne 1 ]; then
    quarantine_worktree_changes
    if [ "$WORKTREE_CREATED" -eq 1 ]; then
      git -C "$ROOT" worktree remove --force "$WORKTREE_DIR" >> "$RUN_LOG" 2>&1 || true
      WORKTREE_CREATED=0
    fi
    if [ "$AUTOMATION_COMMITTED" -eq 0 ] || [ "$AUTOMATION_MERGED" -eq 1 ]; then
      git -C "$ROOT" branch -D "$WORKTREE_BRANCH" >> "$RUN_LOG" 2>&1 || true
    fi
    write_state failed "${AUTOMATION_COMMIT:-$(git -C "$ROOT" rev-parse HEAD 2>/dev/null || true)}" || true
  fi
  rm -f "$CHANGED_FILE"
  if [ "$LOCK_OWNED" -eq 1 ]; then
    rmdir "$LOCK_DIR" 2>/dev/null || true
  fi
}
trap cleanup EXIT
trap 'exit 130' INT TERM

preflight() {
  cd "$ROOT"
  [ -x "$CODEX_BIN" ] || { log "Codex CLI is missing: $CODEX_BIN"; return 1; }
  [ -f "$PROMPT_FILE" ] || { log "Optimizer prompt is missing: $PROMPT_FILE"; return 1; }
  [ "$(git branch --show-current)" = "master" ] || { log "Refusing to run outside master."; return 1; }
  [ "$(git remote get-url origin)" = "$EXPECTED_REMOTE" ] || { log "Origin URL does not match the approved repository."; return 1; }
  git diff --cached --quiet || { log "Refusing to run with staged changes."; return 1; }
  if [ -n "$(git status --porcelain -- "${SOURCE_GUARD_PATHS[@]}")" ]; then
    log "Refusing to run while reviewable source paths contain user changes."
    return 1
  fi
}

if [ "${1:-run}" = "check" ]; then
  preflight
  log "Daily P0 optimizer preflight passed."
  SUCCESS=1
  exit 0
fi

if ! mkdir "$LOCK_DIR" 2>/dev/null; then
  log "Another daily P0 optimizer run is active; skipping."
  SUCCESS=1
  exit 0
fi
LOCK_OWNED=1

preflight

if [ "${SELF_OPT_FORCE:-0}" != "1" ] && [ -f "$LAST_RUN_FILE" ] && [ "$(cat "$LAST_RUN_FILE")" = "$DAY" ]; then
  log "Daily P0 optimizer already completed today; skipping."
  SUCCESS=1
  exit 0
fi

log "Fetching the approved master branch before analysis."
if ! git fetch origin master >> "$RUN_LOG" 2>&1; then
  log "Fetch failed; no code analysis was started."
  exit 1
fi
if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/master)" ]; then
  log "Local master and origin/master differ; refusing autonomous edits."
  exit 1
fi

log "Creating isolated worktree $WORKTREE_BRANCH."
git worktree add -b "$WORKTREE_BRANCH" "$WORKTREE_DIR" HEAD >> "$RUN_LOG" 2>&1
WORKTREE_CREATED=1
log "Starting unattended evidence review and one-P0 implementation."
if ! /usr/bin/perl -e '$timeout = shift @ARGV; alarm $timeout; exec @ARGV' \
  "$CODEX_TIMEOUT_SECONDS" "$CODEX_BIN" exec \
  --ephemeral \
  --sandbox workspace-write \
  --cd "$WORKTREE_DIR" \
  --config 'approval_policy="never"' \
  --output-last-message "$LAST_MESSAGE" \
  - < "$PROMPT_FILE" >> "$RUN_LOG" 2>&1; then
  log "Codex execution failed; generated changes will be quarantined."
  exit 1
fi

cd "$WORKTREE_DIR"
while IFS= read -r status_line; do
  [ -n "$status_line" ] || continue
  path="${status_line:3}"
  path="${path##* -> }"
  if ! is_allowed_path "$path"; then
    log "Rejected unexpected path from autonomous run: $path"
    exit 1
  fi
done < <(git status --porcelain)

if ! git diff --cached --quiet; then
  log "Rejected autonomous staging; the outer controller owns the index."
  exit 1
fi

{
  git diff --name-only -- "${AUTONOMOUS_PATHS[@]}"
  git ls-files --others --exclude-standard -- "${AUTONOMOUS_PATHS[@]}"
} | sort -u > "$CHANGED_FILE"

if [ ! -s "$CHANGED_FILE" ]; then
  remove_worktree
  printf '%s\n' "$DAY" > "$LAST_RUN_FILE"
  write_state no_action "$(git -C "$ROOT" rev-parse HEAD)"
  log "No evidence-backed P0 met the bar; no source change was made."
  SUCCESS=1
  exit 0
fi

while IFS= read -r path; do
  case "$path" in
    scripts/daily-p0-optimize.sh|scripts/install-daily-p0-launchd.sh|scripts/daily-git-save.sh)
      log "Rejected modification of autonomous control file: $path"
      exit 1
      ;;
  esac
done < "$CHANGED_FILE"

FILE_COUNT="$(wc -l < "$CHANGED_FILE" | tr -d ' ')"
if [ "$FILE_COUNT" -gt "$MAX_FILES" ]; then
  log "Rejected oversized change: $FILE_COUNT files exceeds $MAX_FILES."
  exit 1
fi

DIFF_LINES=0
while IFS= read -r path; do
  [ -n "$path" ] || continue
  if git ls-files --error-unmatch "$path" >/dev/null 2>&1; then
    if git diff --numstat -- "$path" | grep -qE '^-' ; then
      log "Rejected binary change: $path"
      exit 1
    fi
    lines="$(git diff --numstat -- "$path" | awk '{sum += $1 + $2} END {print sum + 0}')"
  else
    encoding="$(file -b --mime-encoding "$path")"
    case "$encoding" in
      us-ascii|utf-8) ;;
      *)
        log "Rejected non-text file: $path ($encoding)"
        exit 1
        ;;
    esac
    lines="$(wc -l < "$path" | tr -d ' ')"
  fi
  DIFF_LINES=$((DIFF_LINES + lines))
done < "$CHANGED_FILE"
if [ "$DIFF_LINES" -gt "$MAX_DIFF_LINES" ]; then
  log "Rejected oversized change: $DIFF_LINES lines exceeds $MAX_DIFF_LINES."
  exit 1
fi

GO_FILES=()
while IFS= read -r path; do
  [ -f "$path" ] || continue
  case "$path" in
    *.go) GO_FILES+=("$path") ;;
  esac
done < "$CHANGED_FILE"
if [ "${#GO_FILES[@]}" -gt 0 ]; then
  gofmt -w "${GO_FILES[@]}"
fi

log "Running shell, plist, diff, and full Go race-test gates."
for script in scripts/*.sh; do
  bash -n "$script"
done
for plist in launchd/*.plist; do
  plutil -lint "$plist" >> "$RUN_LOG"
done
git diff --check -- "${AUTONOMOUS_PATHS[@]}"
if ! go test -race ./... >> "$RUN_LOG" 2>&1; then
  log "Full race tests failed; generated changes will be quarantined."
  exit 1
fi

while IFS= read -r path; do
  [ -n "$path" ] || continue
  git add -- "$path"
done < "$CHANGED_FILE"

if git diff --cached -U0 | grep '^+' | grep -v '^+++' | \
  grep -Eqi '(BEGIN (RSA|OPENSSH|EC) PRIVATE KEY|PRIVATE_KEY[[:space:]]*=|MNEMONIC[[:space:]]*=|BOT_TOKEN[[:space:]]*=[[:space:]]*[^$]|API_KEY[[:space:]]*=[[:space:]]*[^$]|api\.telegram\.org/bot[0-9]+:)'; then
  log "Sensitive-value scan rejected the staged change."
  exit 1
fi

git commit -m "fix: autonomous P0 optimization ${DAY}" >> "$RUN_LOG"
AUTOMATION_COMMIT="$(git rev-parse HEAD)"
AUTOMATION_COMMITTED=1

log "Fast-forwarding validated commit into local master."
if ! git -C "$ROOT" merge --ff-only "$AUTOMATION_COMMIT" >> "$RUN_LOG" 2>&1; then
  log "Validated commit could not be fast-forwarded; retained branch $WORKTREE_BRANCH for review."
  exit 1
fi
AUTOMATION_MERGED=1

if [ "${SELF_OPT_PUSH:-1}" = "1" ]; then
  if ! git -C "$ROOT" push origin master >> "$RUN_LOG" 2>&1; then
    log "Commit $AUTOMATION_COMMIT passed all gates but push failed."
    exit 1
  fi
fi

remove_worktree
printf '%s\n' "$DAY" > "$LAST_RUN_FILE"
write_state success "$AUTOMATION_COMMIT"
log "Autonomous P0 passed all gates and was saved as $AUTOMATION_COMMIT."
SUCCESS=1
