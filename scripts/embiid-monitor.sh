#!/bin/bash
# Embiid injury monitor - check every run, alert on status change
STATE_FILE="/Users/murphyma/work/polymarket-go/db/embiid_state.txt"
PUSH_BOT="8405812595:AAEwrI7MdS-miIQtbE_MClwdcqfsjjHMfLo"
CHAT_ID="6695538819"

# Get current status
RESULT=$(curl -s "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/injuries" | python3 -c "
import json,sys
data=json.load(sys.stdin)
for team_data in data.get('injuries',[]):
    if '76' in team_data.get('displayName',''):
        for inj in team_data.get('injuries',[]):
            if 'Embiid' in inj.get('athlete',{}).get('displayName',''):
                status=inj.get('status','?')
                short=inj.get('shortComment','')
                print(f'{status}|||{short}')
" 2>/dev/null)

if [ -z "$RESULT" ]; then
    # No injury entry = cleared from list = Available
    CURRENT_STATUS="Available"
    COMMENT="Embiid cleared from injury report"
else
    CURRENT_STATUS=$(echo "$RESULT" | cut -d'|' -f1)
    COMMENT=$(echo "$RESULT" | cut -d'|' -f4-)
fi

# Read previous status
PREV_STATUS=""
if [ -f "$STATE_FILE" ]; then
    PREV_STATUS=$(cat "$STATE_FILE")
fi

NOW=$(TZ=Asia/Singapore date "+%H:%M SGT")

# Always update state file
echo "$CURRENT_STATUS" > "$STATE_FILE"

# Alert if status changed OR first run
if [ "$PREV_STATUS" != "$CURRENT_STATUS" ]; then
    MSG="🚨 Embiid 伤病状态更新 ($NOW)

状态: $PREV_STATUS → $CURRENT_STATUS

$COMMENT

⏰ 76ers @ Celtics Game 5 · 今晚 07:00 SGT
📎 https://newshare.bwb.online/zh/polymarket/event?slug=nba-phi-bos-2026-04-28&_nobar=true&_needChain=matic"

    curl -s -X POST "https://api.telegram.org/bot${PUSH_BOT}/sendMessage" \
        -d chat_id="$CHAT_ID" \
        -d text="$MSG" \
        -d parse_mode="" > /dev/null
    echo "$(date): Status changed $PREV_STATUS -> $CURRENT_STATUS, alert sent"
else
    echo "$(date): Status unchanged ($CURRENT_STATUS)"
fi
