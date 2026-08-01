#!/usr/bin/env bash
# Promote high-quality scored wallets into an hourly push overlay.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

LOCK_DIR="$ROOT/db/hourly-wallet-promotion.lock"
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
  printf 'hourly-wallet-promotion.already_running lock=%s\n' "$LOCK_DIR"
  exit 0
fi
trap 'rmdir "$LOCK_DIR" 2>/dev/null || true' EXIT

mkdir -p "$ROOT/db" "$ROOT/logs" "$ROOT/reports"

SCORES="${HOURLY_WALLET_SCORES:-$ROOT/db/strategy_iteration/wallet_scores.json}"
OUT="${HOURLY_PUSH_WALLETS:-$ROOT/wallets.hourly-push.txt}"
TMP="$(mktemp "$ROOT/db/hourly-wallet-promotion.XXXXXX")"
FOOTBALL_SCORE_OUT="${HOURLY_FOOTBALL_SCORE_WALLETS:-$ROOT/wallets.football-score-push.txt}"

printf 'hourly-wallet-promotion.start ts=%s scores=%s out=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$SCORES" "$OUT"

if [ "${HOURLY_WALLET_DISCOVER:-1}" = "1" ]; then
  printf 'hourly-wallet-promotion.discover start\n'
  SCORES_BAK="$ROOT/db/hourly-wallet-promotion.wallet_scores.backup.json"
  if [ -s "$SCORES" ]; then
    cp "$SCORES" "$SCORES_BAK"
  fi
  if LEADERBOARD_LIMIT="${HOURLY_LEADERBOARD_LIMIT:-300}" \
    DISCOVER_MARKETS="${HOURLY_DISCOVER_MARKETS:-15}" \
    DISCOVER_TRADE_PAGES="${HOURLY_DISCOVER_TRADE_PAGES:-4}" \
    DISCOVER_HOLDERS="${HOURLY_DISCOVER_HOLDERS:-0}" \
    DISCOVER_ACTIVITY_PAGES="${HOURLY_DISCOVER_ACTIVITY_PAGES:-3}" \
    DISCOVER_MAX_CANDIDATES="${HOURLY_DISCOVER_MAX_CANDIDATES:-800}" \
    DISCOVER_TIMEOUT="${HOURLY_DISCOVER_TIMEOUT:-35m}" \
    SCOUT_PUSH_ENABLED="${HOURLY_SCOUT_PUSH_ENABLED:-true}" \
    "$ROOT/scripts/wallet-discover.sh"; then
    printf 'hourly-wallet-promotion.discover ok\n'
  else
    printf 'hourly-wallet-promotion.discover failed; promoting from existing scores\n'
    if [ -s "$SCORES_BAK" ]; then
      cp "$SCORES_BAK" "$SCORES"
      printf 'hourly-wallet-promotion.discover restored scores=%s\n' "$SCORES"
    fi
  fi
fi

if [ "${HOURLY_FOOTBALL_SCORE_DISCOVER:-1}" = "1" ]; then
  score_tmp="$(mktemp "$ROOT/db/football-score-push.XXXXXX")"
  printf 'hourly-wallet-promotion.football_score start out=%s\n' "$FOOTBALL_SCORE_OUT"
  go build -o "$ROOT/bin/sports-holders-push" ./cmd/sports-holders-push
  if "$ROOT/bin/sports-holders-push" \
    -out "$score_tmp" \
    -list football_score_push \
    -football_score_only \
    -exclude_wallets "$ROOT/wallets.strategy-quarantine.txt,$ROOT/wallets.strategy-review-noise.txt" \
    -scores "$SCORES" \
    -target_categories soccer \
    -markets "${HOURLY_FOOTBALL_SCORE_MARKETS:-40}" \
    -holders "${HOURLY_FOOTBALL_SCORE_HOLDERS:-50}" \
    -max_wallets "${HOURLY_FOOTBALL_SCORE_MAX_WALLETS:-100}" \
    -min_shares "${HOURLY_FOOTBALL_SCORE_MIN_SHARES:-250}" \
    -timeout "${HOURLY_FOOTBALL_SCORE_TIMEOUT:-4m}" && \
    grep -Eq '^0x[0-9a-fA-F]{40}([[:space:]]|$)' "$score_tmp"; then
    mv "$score_tmp" "$FOOTBALL_SCORE_OUT"
    printf 'hourly-wallet-promotion.football_score updated count=%s\n' "$(awk '/^0x[0-9a-fA-F]{40}/ {n++} END {print n+0}' "$FOOTBALL_SCORE_OUT")"
  else
    rm -f "$score_tmp"
    printf 'hourly-wallet-promotion.football_score empty_or_failed preserved=%s\n' "$FOOTBALL_SCORE_OUT"
  fi
fi

python3 - "$SCORES" "$OUT" "$TMP" "$ROOT" <<'PY'
import json
import os
import re
import sys
from datetime import datetime, timezone

scores_path, out_path, tmp_path, root = sys.argv[1:5]

limit = int(os.environ.get("HOURLY_PUSH_LIMIT", "25"))
min_smart = float(os.environ.get("HOURLY_PUSH_MIN_SMART", "70"))
max_bot = float(os.environ.get("HOURLY_PUSH_MAX_BOT", "45"))
min_edge = float(os.environ.get("HOURLY_PUSH_MIN_EDGE", "35"))
min_large = int(os.environ.get("HOURLY_PUSH_MIN_LARGE", "10"))
min_avg_notional = float(os.environ.get("HOURLY_PUSH_MIN_AVG_NOTIONAL", "300"))
min_copy_trades = int(os.environ.get("HOURLY_PUSH_MIN_COPY_TRADES", "3"))
min_copy_roi = float(os.environ.get("HOURLY_PUSH_MIN_COPY_ROI", "5"))
min_target_large = int(os.environ.get("HOURLY_PUSH_MIN_TARGET_LARGE", "3"))
min_target_copy_trades = int(os.environ.get("HOURLY_PUSH_MIN_TARGET_COPY_TRADES", "2"))
min_target_copy_roi = float(os.environ.get("HOURLY_PUSH_MIN_TARGET_COPY_ROI", "5"))

address_re = re.compile(r"0x[0-9a-fA-F]{40}")
bad_flags = {
    "bot_like_flow",
    "fixed_price",
    "fixed_amount",
    "negative_copy_sim",
    "extreme_price_heavy",
}

exclude_files = [
    "wallets.strategy-quarantine.txt",
    "wallets.strategy-review-noise.txt",
    "db/strategy_iteration/wallets.strategy-exclude.txt",
]
already_pushed_files = [
    "wallets.strategy-push.txt",
    "wallets.leaderboard-push.txt",
    "wallets.leaderboard-watch.txt",
    "wallets.leaderboard-sports-push.txt",
    "wallets.sports-holders-push.txt",
    "wallets.football-score-push.txt",
]

def wallet_set(paths):
    out = set()
    for rel in paths:
        path = rel if os.path.isabs(rel) else os.path.join(root, rel)
        try:
            with open(path, encoding="utf-8") as fh:
                for line in fh:
                    match = address_re.search(line)
                    if match:
                        out.add(match.group(0).lower())
        except FileNotFoundError:
            pass
    return out

excluded = wallet_set(exclude_files)
already = wallet_set(already_pushed_files)

with open(scores_path, encoding="utf-8") as fh:
    rows = json.load(fh)

def sources_summary(sources):
    keys = sorted(k for k, v in (sources or {}).items() if v)
    return ",".join(keys[:5]) if keys else "-"

def source_bonus(sources):
    sources = sources or {}
    score = 0.0
    for key, value in sources.items():
        if not value:
            continue
        if key == "leaderboard_profit_7d":
            score += 20
        elif key == "leaderboard_profit_30d":
            score += 14
        elif key == "sports_tape":
            score += 10
        elif key == "holder":
            score += 4
    return score

def qualifies(row):
    addr = str(row.get("address", "")).lower()
    if not address_re.fullmatch(addr):
        return False
    if addr in excluded or addr in already:
        return False
    if row.get("tier") not in {"A", "B", "C"}:
        return False
    if float(row.get("smart_money_score") or 0) < min_smart:
        return False
    if float(row.get("bot_score") or 999) > max_bot:
        return False
    if float(row.get("edge_score") or 0) < min_edge:
        return False
    if bad_flags & set(row.get("risk_flags") or []):
        return False
    stats = row.get("stats") or {}
    if int(stats.get("large_trades") or 0) < min_large:
        return False
    if float(stats.get("avg_trade_notional") or 0) < min_avg_notional:
        return False
    copy_ok = (
        int(stats.get("copy_closed_trades") or 0) >= min_copy_trades
        and float(stats.get("copy_roi") or 0) >= min_copy_roi
        and float(stats.get("copy_pnl") or 0) > 0
    )
    target_ok = (
        int(stats.get("target_large_trades") or 0) >= min_target_large
        and int(stats.get("target_copy_closed_trades") or 0) >= min_target_copy_trades
        and float(stats.get("target_copy_roi") or 0) >= min_target_copy_roi
        and float(stats.get("target_copy_pnl") or 0) > 0
    )
    leaderboard_profit = any(k.startswith("leaderboard_profit") and v for k, v in (row.get("sources") or {}).items())
    return copy_ok or target_ok or leaderboard_profit

def rank(row):
    stats = row.get("stats") or {}
    return (
        float(row.get("smart_money_score") or 0) * 1.2
        + float(row.get("edge_score") or 0) * 0.6
        - float(row.get("bot_score") or 0) * 1.5
        + min(float(stats.get("large_trades") or 0), 250) * 0.08
        + min(float(stats.get("target_large_trades") or 0), 120) * 0.10
        + max(float(stats.get("copy_roi") or 0), 0) * 0.20
        + max(float(stats.get("target_copy_roi") or 0), 0) * 0.25
        + source_bonus(row.get("sources"))
    )

selected = sorted((row for row in rows if qualifies(row)), key=lambda row: (-rank(row), str(row.get("address", ""))))[:limit]

lines = [
    f"# generated by hourly-wallet-promotion at {datetime.now(timezone.utc).isoformat()}",
    "# list=hourly_push; overlay merged by scripts/start-whale-push.sh",
]
for row in selected:
    stats = row.get("stats") or {}
    lines.append(
        "{addr} # list=hourly_push tier={tier} smart={smart:.1f} bot={bot:.1f} "
        "hourlyScore={score:.1f} copyROI={copy_roi:.1f}% copyT={copy_t} "
        "targetCopyROI={target_copy_roi:.1f}% targetCopyT={target_copy_t} "
        "large={large} targetLarge={target_large} avgNotional=${avg:.0f} sources={sources}".format(
            addr=str(row.get("address", "")).lower(),
            tier=row.get("tier") or "?",
            smart=float(row.get("smart_money_score") or 0),
            bot=float(row.get("bot_score") or 0),
            score=rank(row),
            copy_roi=float(stats.get("copy_roi") or 0),
            copy_t=int(stats.get("copy_closed_trades") or 0),
            target_copy_roi=float(stats.get("target_copy_roi") or 0),
            target_copy_t=int(stats.get("target_copy_closed_trades") or 0),
            large=int(stats.get("large_trades") or 0),
            target_large=int(stats.get("target_large_trades") or 0),
            avg=float(stats.get("avg_trade_notional") or 0),
            sources=sources_summary(row.get("sources")),
        )
    )

with open(tmp_path, "w", encoding="utf-8") as fh:
    fh.write("\n".join(lines) + "\n")

print(f"hourly-wallet-promotion.selected count={len(selected)} limit={limit}")
for row in selected[:10]:
    stats = row.get("stats") or {}
    print(
        "hourly-wallet-promotion.candidate "
        f"wallet={str(row.get('address', '')).lower()} tier={row.get('tier')} "
        f"smart={float(row.get('smart_money_score') or 0):.1f} "
        f"bot={float(row.get('bot_score') or 0):.1f} "
        f"large={int(stats.get('large_trades') or 0)} "
        f"sources={sources_summary(row.get('sources'))}"
    )
PY

if [ -f "$OUT" ] && cmp -s "$TMP" "$OUT"; then
  rm -f "$TMP"
  printf 'hourly-wallet-promotion.unchanged out=%s\n' "$OUT"
  changed=0
else
  mv "$TMP" "$OUT"
  printf 'hourly-wallet-promotion.updated out=%s count=%s\n' "$OUT" "$(awk '/^0x[0-9a-fA-F]{40}/ {n++} END {print n+0}' "$OUT")"
  changed=1
fi

if [ "${HOURLY_PROMOTION_RESTART:-1}" = "1" ]; then
  screen_name="${HOURLY_PROMOTION_SCREEN:-polymarket-whale-push}"
  printf 'hourly-wallet-promotion.restart screen=%s changed=%s\n' "$screen_name" "$changed"
  screen -S "$screen_name" -X quit >/dev/null 2>&1 || true
  old_pid="$(cat "$ROOT/db/bot.pid" 2>/dev/null || true)"
  if [ -n "$old_pid" ]; then
    kill "$old_pid" >/dev/null 2>&1 || true
  fi
  sleep 1
  screen -dmS "$screen_name" "$ROOT/scripts/start-whale-push.sh"
fi

printf 'hourly-wallet-promotion.done ts=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
