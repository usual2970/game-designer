#!/usr/bin/env bash
# verify-deployed.sh — Run deployed verification against a live Game Designer Server
#
# Usage: ./scripts/verify-deployed.sh <deployed_url>
#
# Runs health plus slot-loop checks against a deployed backend.

set -euo pipefail

if [ -z "${1:-}" ]; then
  echo "Usage: $0 <deployed_url>"
  echo "Example: $0 https://my-game.fake.local"
  exit 1
fi

DEPLOYED_URL="$1"
TIMEOUT=30

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

pass=0
fail=0
errors=()

check() {
  local name="$1"
  shift
  if "$@" 2>/dev/null; then
    echo -e "  ${GREEN}PASS${NC} $name"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} $name"
    ((fail++))
    errors+=("$name")
  fi
}

echo "=== Game Designer Server — Deployed Verification ==="
echo "Target: $DEPLOYED_URL"
echo ""

# Step 1: Health check
echo "1. Health check"
HTTP_CODE=$(curl -sf -o /dev/null -w "%{http_code}" "$DEPLOYED_URL/api/v1/session" --max-time "$TIMEOUT" 2>/dev/null || echo "000")
if [ "$HTTP_CODE" != "000" ]; then
  echo -e "  ${GREEN}PASS${NC} endpoint reachable (HTTP $HTTP_CODE)"
  ((pass++))
else
  echo -e "  ${RED}FAIL${NC} endpoint unreachable"
  ((fail++))
  errors+=("endpoint-reachable")
fi

# Step 2: Session creation
echo ""
echo "2. Session"
SESSION_RESP=$(curl -sf -X POST "$DEPLOYED_URL/api/v1/session" \
  -H "Content-Type: application/json" \
  -d '{"playerId":"verify-deployed","nickname":"DeployBot"}' --max-time "$TIMEOUT" 2>/dev/null || echo "")

if [ -n "$SESSION_RESP" ] && echo "$SESSION_RESP" | grep -q '"token"'; then
  TOKEN=$(echo "$SESSION_RESP" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  echo -e "  ${GREEN}PASS${NC} session created"
  ((pass++))

  # Step 3: Slot config
  echo ""
  echo "3. Slot config"
  CONFIG_RESP=$(curl -sf "$DEPLOYED_URL/api/v1/slot/config" \
    -H "X-Session-Token: $TOKEN" --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$CONFIG_RESP" | grep -q '"reels"'; then
    echo -e "  ${GREEN}PASS${NC} slot config read"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} slot config"
    ((fail++))
    errors+=("slot-config")
  fi

  # Step 4: Balance
  echo ""
  echo "4. Balance"
  BAL_RESP=$(curl -sf "$DEPLOYED_URL/api/v1/balance" \
    -H "X-Session-Token: $TOKEN" --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$BAL_RESP" | grep -q '"balance"'; then
    echo -e "  ${GREEN}PASS${NC} balance read"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} balance read"
    ((fail++))
    errors+=("balance-read")
  fi

  # Step 5: Spin
  echo ""
  echo "5. Spin"
  SPIN_RESP=$(curl -sf -X POST "$DEPLOYED_URL/api/v1/spin" \
    -H "Content-Type: application/json" \
    -H "X-Session-Token: $TOKEN" \
    -d '{"wager":10}' --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$SPIN_RESP" | grep -q '"spinId"'; then
    echo -e "  ${GREEN}PASS${NC} spin resolved"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} spin"
    ((fail++))
    errors+=("spin")
  fi

  # Step 6: Spin history
  echo ""
  echo "6. Spin history"
  HIST_RESP=$(curl -sf "$DEPLOYED_URL/api/v1/spin/history" \
    -H "X-Session-Token: $TOKEN" --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$HIST_RESP" | grep -q '"entries"'; then
    echo -e "  ${GREEN}PASS${NC} spin history read"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} spin history"
    ((fail++))
    errors+=("spin-history")
  fi

  # Step 7: Leaderboard
  echo ""
  echo "7. Leaderboard"
  LB_RESP=$(curl -sf "$DEPLOYED_URL/api/v1/leaderboard" \
    -H "X-Session-Token: $TOKEN" --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$LB_RESP" | grep -q '"entries"'; then
    echo -e "  ${GREEN}PASS${NC} leaderboard read"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} leaderboard read"
    ((fail++))
    errors+=("leaderboard-read")
  fi
else
  echo -e "  ${RED}FAIL${NC} session creation failed"
  ((fail++))
  errors+=("session-create")
fi

# Summary
echo ""
echo "=== Results ==="
echo -e "  Passed: ${GREEN}$pass${NC}"
echo -e "  Failed: ${RED}$fail${NC}"

if [ "$fail" -gt 0 ]; then
  echo ""
  echo -e "${RED}FAILED:${NC}"
  for e in "${errors[@]}"; do
    echo "  - $e"
  done
  echo ""
  echo '{"success":false,"message":"Deployed verification failed","code":"VERIFICATION_FAILED","details":{"passed":'"$pass"',"failed":'"$fail"'}}'
  exit 1
fi

echo ""
echo '{"success":true,"message":"All deployed verification checks passed","code":"SUCCESS","details":{"passed":'"$pass"'}}'
