#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

go build -o "$ROOT/bin/smartmoney-wallet-policy" ./cmd/smartmoney-wallet-policy
exec "$ROOT/bin/smartmoney-wallet-policy" \
  -journal "${SMARTMONEY_PAPER_JOURNAL:-$ROOT/db/smartmoney-paper/journal}" \
  -whale_log "${SMARTMONEY_PAPER_WHALE_LOG:-$ROOT/db/smartmoney-paper/journal/whale_trades.jsonl}" \
  -json_out "${SMARTMONEY_PAPER_WALLET_POLICY_JSON:-$ROOT/db/smartmoney-paper/wallet-performance.json}" \
  -report "${SMARTMONEY_PAPER_WALLET_POLICY_REPORT:-$ROOT/reports/smartmoney-paper-wallets.md}" \
  -promoted "${SMARTMONEY_PAPER_PROMOTED_WALLETS:-$ROOT/db/smartmoney-paper/wallets.paper-promoted.txt}" \
  -demoted "${SMARTMONEY_PAPER_DEMOTED_WALLETS:-$ROOT/db/smartmoney-paper/wallets.paper-demoted.txt}" \
  -min_positions "${SMARTMONEY_PAPER_WALLET_MIN_POSITIONS:-10}" \
  -promote_min_net "${SMARTMONEY_PAPER_WALLET_PROMOTE_MIN_NET:-5}" \
  -promote_min_roi "${SMARTMONEY_PAPER_WALLET_PROMOTE_MIN_ROI:-2}" \
  -demote_max_net "${SMARTMONEY_PAPER_WALLET_DEMOTE_MAX_NET:--5}" "$@"
