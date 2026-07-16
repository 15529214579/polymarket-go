#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

mkdir -p "$ROOT/db/strategy_iteration"

go build -o bin/wallet-discover ./cmd/wallet-discover

if [ "$#" -gt 0 ]; then
  exec "$ROOT/bin/wallet-discover" "$@"
fi

EXISTING_WALLETS_FILE="${DISCOVER_EXISTING_WALLETS:-$ROOT/db/strategy_iteration/wallets.existing-seed.txt}"
if [ -z "${DISCOVER_EXISTING_WALLETS:-}" ]; then
  awk '
    {
      sub(/#.*/, "")
      for (i = 1; i <= NF; i++) {
        if ($i ~ /^0x[0-9a-fA-F]{40}$/) {
          print tolower($i)
          break
        }
      }
    }
  ' \
    "$ROOT/wallets.whale-push.txt" \
    "$ROOT/wallets.strategy-push.txt" \
    "$ROOT/wallets.strategy-core.txt" \
    "$ROOT/wallets.strategy-watch.txt" \
    "$ROOT/wallets.strategy-sports.txt" \
    "$ROOT/wallets.strategy-target.txt" \
    "$ROOT/wallets.strategy-scout.txt" \
    "$ROOT/wallets.strategy-flow.txt" \
    "$ROOT/wallets.leaderboard-sports-push.txt" \
    "$ROOT/wallets.sports-holders-push.txt" \
    "$ROOT/wallets.strategy-tape.txt" \
    "$ROOT/wallets.strategy-tape-observe.txt" \
    "$ROOT/wallets.strategy-tape-probation.txt" \
    "$ROOT/wallets.strategy-tape-candidates.txt" \
    "$ROOT/wallets.strategy-tape-follow.txt" \
    "$ROOT/wallets.strategy-tape-reversal.txt" \
    "$ROOT/wallets.strategy-review-noise.txt" \
    "$ROOT/wallets.strategy-positive.txt" 2>/dev/null | sort -u > "$EXISTING_WALLETS_FILE"
fi

EXCLUDE_WALLETS_FILE="${STRATEGY_EXCLUDE_WALLETS:-$ROOT/db/strategy_iteration/wallets.strategy-exclude.txt}"
if [ -z "${STRATEGY_EXCLUDE_WALLETS:-}" ]; then
  awk '
    {
      sub(/#.*/, "")
      for (i = 1; i <= NF; i++) {
        if ($i ~ /^0x[0-9a-fA-F]{40}$/) {
          print tolower($i)
          break
        }
      }
    }
  ' \
    "$ROOT/wallets.strategy-quarantine.txt" \
    "$ROOT/wallets.strategy-review-noise.txt" 2>/dev/null | sort -u > "$EXCLUDE_WALLETS_FILE"
fi

set -- \
  -existing_wallets "$EXISTING_WALLETS_FILE" \
  -sports_tape_wallets "${SPORTS_TAPE_WALLETS:-$ROOT/wallets.sports-tape.txt}" \
  -retain_wallets "${DISCOVER_RETAIN_WALLETS:-$ROOT/wallets.strategy-retain.txt}" \
  -output_dir "$ROOT/db/strategy_iteration" \
  -report "$ROOT/reports/strategy_iteration.md" \
  -generated_tier C \
  -generated_wallets "$ROOT/db/strategy_iteration/wallets.generated.txt" \
  -auto_wallets "$ROOT/db/strategy_iteration/wallets.smartmoney-auto.txt" \
  -prompt_wallets "$ROOT/db/strategy_iteration/wallets.smartmoney-prompt.txt" \
  -positive_wallets "$ROOT/wallets.strategy-positive.txt" \
  -leaderboard_limit "${LEADERBOARD_LIMIT:-500}" \
  -leaderboard_windows "${LEADERBOARD_WINDOWS:-7d,30d,all}" \
  -leaderboard_kinds "${LEADERBOARD_KINDS:-profit,volume}" \
  -target_categories "${DISCOVER_TARGET_CATEGORIES:-basketball,soccer,esports}" \
  -markets "${DISCOVER_MARKETS:-25}" \
  -trade_pages "${DISCOVER_TRADE_PAGES:-7}" \
  -holders "${DISCOVER_HOLDERS:-1}" \
  -max_candidates "${DISCOVER_MAX_CANDIDATES:-1500}" \
  -activity_pages "${DISCOVER_ACTIVITY_PAGES:-4}" \
  -copy_stake "${COPY_STAKE_USD:-10}" \
  -copy_slippage_bp "${COPY_SLIPPAGE_BP:-50}" \
  -copy_fee_bp "${COPY_FEE_BP:-0}" \
  -timeout "${DISCOVER_TIMEOUT:-45m}"

"$ROOT/bin/wallet-discover" "$@"

go build -o bin/strategy-lab ./cmd/strategy-lab
exec "$ROOT/bin/strategy-lab" \
  -scores "$ROOT/db/strategy_iteration/wallet_scores.json" \
  -report "$ROOT/reports/strategy_lab.md" \
  -exclude_wallets "$EXCLUDE_WALLETS_FILE" \
  -core_wallets "$ROOT/wallets.strategy-core.txt" \
  -watch_wallets "$ROOT/wallets.strategy-watch.txt" \
  -sports_wallets "$ROOT/wallets.strategy-sports.txt" \
  -scout_wallets "$ROOT/wallets.strategy-scout.txt" \
  -target_wallets "$ROOT/wallets.strategy-target.txt" \
  -flow_wallets "$ROOT/wallets.strategy-flow.txt" \
  -tape_wallets "$ROOT/wallets.strategy-tape.txt" \
  -tape_probation_wallets "$ROOT/wallets.strategy-tape-probation.txt" \
  -tape_observe_wallets "$ROOT/wallets.strategy-tape-observe.txt" \
  -tape_edge_hot_wallets "${SMARTMONEY_TAPE_EDGEHOT_WALLETS:-$ROOT/wallets.strategy-tape-edgehot.txt}" \
  -sports_tape_wallets "$ROOT/wallets.sports-tape.txt" \
  -push_wallets "$ROOT/wallets.strategy-push.txt" \
  -min_wallets "${CORE_MIN_WALLETS:-10}" \
  -min_total_copy_trades "${CORE_MIN_TOTAL_COPY_TRADES:-120}" \
  -floor_min_copy_trades "${CORE_MIN_COPY_TRADES:-8}" \
  -floor_min_copy_roi "${CORE_MIN_COPY_ROI:-10}" \
  -floor_min_copy_pnl "${CORE_MIN_COPY_PNL:-25}" \
  -floor_min_copy_win_rate "${CORE_MIN_COPY_WIN_RATE:-60}" \
  -floor_min_closed_roi "${CORE_MIN_CLOSED_ROI:-0}" \
  -floor_min_smart "${CORE_MIN_SMART:-70}" \
  -ceiling_max_bot "${CORE_MAX_BOT:-30}" \
  -watch_max_wallets "${WATCH_MAX_WALLETS:-20}" \
  -watch_max_bot "${WATCH_MAX_BOT:-35}" \
  -watch_min_smart "${WATCH_MIN_SMART:-60}" \
  -watch_min_copy_trades "${WATCH_MIN_COPY_TRADES:-3}" \
  -watch_min_copy_roi "${WATCH_MIN_COPY_ROI:-5}" \
  -watch_min_copy_pnl "${WATCH_MIN_COPY_PNL:-10}" \
  -watch_min_copy_win_rate "${WATCH_MIN_COPY_WIN_RATE:-50}" \
  -watch_min_large_trades "${WATCH_MIN_LARGE_TRADES:-5}" \
  -watch_min_avg_notional "${WATCH_MIN_AVG_NOTIONAL:-250}" \
  -sports_max_wallets "${SPORTS_MAX_WALLETS:-10}" \
  -sports_max_bot "${SPORTS_MAX_BOT:-35}" \
  -sports_min_smart "${SPORTS_MIN_SMART:-70}" \
  -sports_min_edge "${SPORTS_MIN_EDGE:-50}" \
  -sports_min_target_trades "${SPORTS_MIN_TARGET_TRADES:-10}" \
  -sports_min_target_ratio "${SPORTS_MIN_TARGET_RATIO:-0.05}" \
  -sports_min_target_copy_trades "${SPORTS_MIN_TARGET_COPY_TRADES:-5}" \
  -sports_min_target_copy_roi "${SPORTS_MIN_TARGET_COPY_ROI:-5}" \
  -sports_min_target_copy_pnl "${SPORTS_MIN_TARGET_COPY_PNL:-5}" \
  -sports_min_target_large "${SPORTS_MIN_TARGET_LARGE:-5}" \
  -scout_max_wallets "${SCOUT_MAX_WALLETS:-10}" \
  -scout_max_bot "${SCOUT_MAX_BOT:-45}" \
  -scout_min_smart "${SCOUT_MIN_SMART:-80}" \
  -scout_min_edge "${SCOUT_MIN_EDGE:-60}" \
  -scout_min_large_trades "${SCOUT_MIN_LARGE_TRADES:-50}" \
  -scout_min_avg_notional "${SCOUT_MIN_AVG_NOTIONAL:-500}" \
  -scout_push_enabled="${SCOUT_PUSH_ENABLED:-false}" \
  -target_max_wallets "${TARGET_MAX_WALLETS:-10}" \
  -target_max_bot "${TARGET_MAX_BOT:-45}" \
  -target_min_smart "${TARGET_MIN_SMART:-70}" \
  -target_min_edge "${TARGET_MIN_EDGE:-40}" \
  -target_min_trades "${TARGET_MIN_TRADES:-50}" \
  -target_min_ratio "${TARGET_MIN_RATIO:-0.20}" \
  -target_min_copy_trades "${TARGET_MIN_COPY_TRADES:-5}" \
  -target_min_copy_roi "${TARGET_MIN_COPY_ROI:-10}" \
  -target_min_large_trades "${TARGET_MIN_LARGE_TRADES:-50}" \
  -target_min_avg_notional "${TARGET_MIN_AVG_NOTIONAL:-100}" \
  -flow_max_wallets "${FLOW_MAX_WALLETS:-15}" \
  -flow_max_bot "${FLOW_MAX_BOT:-45}" \
  -flow_min_smart "${FLOW_MIN_SMART:-70}" \
  -flow_min_edge "${FLOW_MIN_EDGE:-35}" \
  -flow_min_recent_trades "${FLOW_MIN_RECENT_TRADES:-2}" \
  -flow_min_large_trades "${FLOW_MIN_LARGE_TRADES:-5}" \
  -flow_min_avg_notional "${FLOW_MIN_AVG_NOTIONAL:-500}" \
  -tape_max_wallets "${TAPE_MAX_WALLETS:-8}" \
  -tape_max_bot "${TAPE_MAX_BOT:-45}" \
  -tape_min_direct_max_buy "${TAPE_MIN_DIRECT_MAX_BUY:-5000}" \
  -tape_observe_min_buy_notional "${TAPE_OBSERVE_MIN_BUY_NOTIONAL:-${SPORTS_TAPE_ALERT_OBSERVE_BURST_MIN_NOTIONAL:-6000}}" \
  -tape_min_scored_max_buy "${TAPE_MIN_SCORED_MAX_BUY:-2500}" \
  -tape_min_smart "${TAPE_MIN_SMART:-70}" \
  -tape_min_target_copy_trades "${TAPE_MIN_TARGET_COPY_TRADES:-2}" \
  -tape_min_target_copy_roi "${TAPE_MIN_TARGET_COPY_ROI:-25}" \
  -edge_snapshots "${WHALE_EDGE_SNAPSHOTS:-$ROOT/db/strategy_iteration/whale_edge_snapshots.jsonl}" \
  -tape_edge_block_15m_samples "${TAPE_EDGE_BLOCK_15M_SAMPLES:-2}" \
  -tape_edge_block_15m_max_avg_pp "${TAPE_EDGE_BLOCK_15M_MAX_AVG_PP:--1}" \
  -tape_edge_block_1h_samples "${TAPE_EDGE_BLOCK_1H_SAMPLES:-1}" \
  -tape_edge_block_1h_max_avg_pp "${TAPE_EDGE_BLOCK_1H_MAX_AVG_PP:--5}" \
  -top_n "${CORE_TOP_STRATEGIES:-12}"
