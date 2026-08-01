#!/bin/bash
set -euo pipefail

umask 077

ROOT="${POLYMARKET_GO_ROOT:-/Users/murphyma/work/polymarket-go}"
EXPECTED_WALLET="0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e"
LIVE_DIR="$ROOT/db/live"
BIN_DIR="$ROOT/bin"
LOG_DIR="$ROOT/logs"
ENABLE_FILE="$LIVE_DIR/redeem.enabled"
DISABLE_FILE="$LIVE_DIR/redeem.disabled"
LOCK_DIR="$LIVE_DIR/hourly-redeem.lock"
STATE_FILE="$LIVE_DIR/hourly-redeem-state.json"
REDEEMED_FILE="$LIVE_DIR/redeemed.json"
MAINTENANCE_BIN="$BIN_DIR/trade-maintenance"
ACTION="${1:-run}"

if [ -L "$LIVE_DIR" ]; then
  printf 'hourly_redeem.failed reason=live_state_dir_is_symlink\n' >&2
  exit 1
fi
mkdir -p "$LIVE_DIR" "$BIN_DIR" "$LOG_DIR"
cd "$ROOT"

write_state() {
  local status="$1"
  local reason="$2"
  local exit_code="${3:-0}"
  local checked_at tmp
  checked_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  tmp="$STATE_FILE.tmp.$$"
  printf '{\n  "checked_at": "%s",\n  "status": "%s",\n  "reason": "%s",\n  "exit_code": %s,\n  "wallet": "%s",\n  "state_scope": "db/live"\n}\n' \
    "$checked_at" "$status" "$reason" "$exit_code" "$EXPECTED_WALLET" > "$tmp"
  chmod 600 "$tmp"
  mv "$tmp" "$STATE_FILE"
}

case "$ACTION" in
  status)
    if [ -f "$STATE_FILE" ]; then
      /bin/cat "$STATE_FILE"
    else
      printf 'hourly_redeem.status state=missing path=%s\n' "$STATE_FILE"
    fi
    exit 0
    ;;
  run) ;;
  *)
    printf 'usage: %s [run|status]\n' "$0" >&2
    exit 2
    ;;
esac

if [ -L "$LOCK_DIR" ]; then
  printf 'hourly_redeem.failed reason=lock_is_symlink\n' >&2
  exit 1
fi
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
  lock_pid="$(/bin/cat "$LOCK_DIR/pid" 2>/dev/null || true)"
  lock_active=0
  if [ -n "$lock_pid" ] && kill -0 "$lock_pid" 2>/dev/null; then
    lock_command="$(/bin/ps -o command= -p "$lock_pid" 2>/dev/null || true)"
    case "$lock_command" in
      *hourly-live-redeem.sh*) lock_active=1 ;;
      "") lock_active=1 ;;
    esac
  fi
  if [ "$lock_active" = "1" ]; then
    printf 'hourly_redeem.skip reason=already_running\n'
    exit 0
  fi
  rm -f "$LOCK_DIR/pid"
  if ! rmdir "$LOCK_DIR" 2>/dev/null || ! mkdir "$LOCK_DIR" 2>/dev/null; then
    printf 'hourly_redeem.skip reason=lock_unavailable\n'
    exit 0
  fi
fi
printf '%s\n' "$$" > "$LOCK_DIR/pid"

build_tmp=""
cleanup() {
  if [ -n "$build_tmp" ] && [ -f "$build_tmp" ]; then
    rm -f "$build_tmp"
  fi
  rm -f "$LOCK_DIR/pid"
  rmdir "$LOCK_DIR" 2>/dev/null || true
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

if [ -e "$DISABLE_FILE" ] || [ -L "$DISABLE_FILE" ]; then
  write_state "skipped" "disabled" 0
  printf 'hourly_redeem.skip reason=disabled\n'
  exit 0
fi

if [ ! -f "$ENABLE_FILE" ] || [ -L "$ENABLE_FILE" ]; then
  write_state "skipped" "not_armed" 0
  printf 'hourly_redeem.skip reason=not_armed\n'
  exit 0
fi

enable_mode="$(stat -f '%Lp' "$ENABLE_FILE" 2>/dev/null || stat -c '%a' "$ENABLE_FILE" 2>/dev/null || true)"
enable_owner="$(stat -f '%u' "$ENABLE_FILE" 2>/dev/null || stat -c '%u' "$ENABLE_FILE" 2>/dev/null || true)"
enable_wallet="$(/bin/cat "$ENABLE_FILE")"
if [ "$enable_mode" != "600" ]; then
  write_state "skipped" "invalid_arm_mode" 0
  printf 'hourly_redeem.skip reason=invalid_arm_mode\n'
  exit 0
fi
if [ "$enable_owner" != "$(id -u)" ]; then
  write_state "skipped" "invalid_arm_owner" 0
  printf 'hourly_redeem.skip reason=invalid_arm_owner\n'
  exit 0
fi
if [ "$enable_wallet" != "$EXPECTED_WALLET" ]; then
  write_state "skipped" "wallet_binding_mismatch" 0
  printf 'hourly_redeem.skip reason=wallet_binding_mismatch\n'
  exit 0
fi

if [ -z "${BW_SESSION:-}" ]; then
  write_state "skipped" "signer_unavailable" 0
  printf 'hourly_redeem.skip reason=signer_unavailable\n'
  exit 0
fi

GO_BIN="${POLYMARKET_GO_BIN:-/usr/local/go/bin/go}"
if [ ! -x "$GO_BIN" ]; then
  if [ -x "/opt/homebrew/bin/go" ]; then
    GO_BIN="/opt/homebrew/bin/go"
  else
    GO_BIN="$(command -v go 2>/dev/null || true)"
  fi
fi
BW_BIN="${POLYMARKET_BW_BIN:-/opt/homebrew/bin/bw}"
if [ ! -x "$BW_BIN" ]; then
  if [ -x "/usr/local/bin/bw" ]; then
    BW_BIN="/usr/local/bin/bw"
  else
    BW_BIN="$(command -v bw 2>/dev/null || true)"
  fi
fi
if [ -z "$GO_BIN" ] || [ ! -x "$GO_BIN" ] || [ -z "$BW_BIN" ] || [ ! -x "$BW_BIN" ]; then
  write_state "failed" "dependency_missing" 1
  printf 'hourly_redeem.failed reason=dependency_missing\n' >&2
  exit 1
fi
export POLYMARKET_BW_BIN="$BW_BIN"

if [ ! -f "$ROOT/go.mod" ] || [ ! -d "$ROOT/cmd/trade" ]; then
  write_state "failed" "invalid_project_root" 1
  printf 'hourly_redeem.failed reason=invalid_project_root\n' >&2
  exit 1
fi

write_state "running" "building_maintenance_binary" 0
build_tmp="$(mktemp "$BIN_DIR/.trade-maintenance.XXXXXX")"
if ! "$GO_BIN" build -o "$build_tmp" ./cmd/trade; then
  write_state "failed" "build_failed" 1
  printf 'hourly_redeem.failed reason=build_failed\n' >&2
  exit 1
fi
chmod 700 "$build_tmp"
mv "$build_tmp" "$MAINTENANCE_BIN"
build_tmp=""

write_state "running" "redeeming" 0
if "$MAINTENANCE_BIN" \
  -redeem-all \
  -redeemed-state "$REDEEMED_FILE" \
  -expected-address "$EXPECTED_WALLET"; then
  write_state "ok" "completed" 0
  printf 'hourly_redeem.result status=ok\n'
else
  rc=$?
  write_state "failed" "redeem_command_failed" "$rc"
  printf 'hourly_redeem.failed reason=redeem_command_failed exit_code=%s\n' "$rc" >&2
  exit "$rc"
fi
