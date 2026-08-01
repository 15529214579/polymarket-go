#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export PATH="/usr/local/go/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

go build -o "$ROOT/bin/smartmoney-shadow-report" ./cmd/smartmoney-shadow-report
exec "$ROOT/bin/smartmoney-shadow-report" \
  -log "${SMARTMONEY_PAPER_LOG:-$ROOT/logs/smartmoney-paper.log}" \
  -out "${SMARTMONEY_SHADOW_REPORT:-$ROOT/reports/smartmoney-exit-shadow.md}" \
  -json_out "${SMARTMONEY_SHADOW_REPORT_JSON:-$ROOT/db/smartmoney-paper/exit-shadow-report.json}" "$@"
