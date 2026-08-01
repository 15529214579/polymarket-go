#!/bin/bash
set -euo pipefail

ROOT="${POLYMARKET_GO_ROOT:-/Users/murphyma/work/polymarket-go}"
LABEL="com.polymarket-go.hourly-live-redeem"
SOURCE="$ROOT/launchd/${LABEL}.plist"
TARGET="$HOME/Library/LaunchAgents/${LABEL}.plist"
DOMAIN="gui/$(id -u)"

plutil -lint "$SOURCE"
mkdir -p "$HOME/Library/LaunchAgents"
launchctl bootout "$DOMAIN/$LABEL" >/dev/null 2>&1 || true
cp "$SOURCE" "$TARGET"
plutil -lint "$TARGET"
launchctl bootstrap "$DOMAIN" "$TARGET"
launchctl enable "$DOMAIN/$LABEL"
launchctl print "$DOMAIN/$LABEL"
