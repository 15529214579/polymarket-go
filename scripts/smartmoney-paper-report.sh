#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

go build -o "$ROOT/bin/smartmoney-report" ./cmd/smartmoney-report
exec "$ROOT/bin/smartmoney-report" \
  -journal "${SMARTMONEY_PAPER_JOURNAL:-$ROOT/db/smartmoney-paper/journal}" \
  -positions "${SMARTMONEY_PAPER_POSITIONS:-$ROOT/db/smartmoney-paper/positions.json}" \
  -out "${SMARTMONEY_PAPER_REPORT_OUT:-$ROOT/reports/smartmoney-paper-pnl.md}" \
  -json_out "${SMARTMONEY_PAPER_REPORT_JSON:-$ROOT/db/smartmoney-paper/pnl-report.json}" "$@"
