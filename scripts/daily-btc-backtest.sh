#!/bin/bash
# Daily BTC Markov + Black-Scholes backtest
# Runs at 08:00 SGT (00:00 UTC) via crontab
# Schedule: 0 0 * * * /Users/murphyma/work/polymarket-go/scripts/daily-btc-backtest.sh

set -euo pipefail

WORKDIR="$HOME/work/polymarket-go"
LOG="$WORKDIR/db/btc_backtest.log"
DB_DIR="$WORKDIR/db"
DATE=$(date -u +%Y-%m-%d)
TS=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

cd "$WORKDIR"

echo "" >> "$LOG"
echo "═══════════════════════════════════════════════════" >> "$LOG"
echo "[$TS] Daily BTC backtest starting" >> "$LOG"

go run ./cmd/backtest \
    -mode=btc-markov \
    -days=90 \
    -train_pct=0.67 \
    -db_dir="$DB_DIR" \
    2>&1 | tee -a "$LOG"

echo "[$DATE] BTC backtest completed" >> "$LOG"

# Extract key metrics for quick review
tail -20 "$LOG" | grep -E "(Accuracy|Simulated PnL|最大价差|BUY_YES|BUY_NO|Current BTC)" >> /dev/null 2>&1 || true
