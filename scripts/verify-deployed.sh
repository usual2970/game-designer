#!/usr/bin/env bash
# verify-deployed.sh — Run deployed verification against a live Game Designer Server
#
# Usage: ./scripts/verify-deployed.sh <deployed_url>
#
# Runs health plus activity-loop checks against a deployed backend.

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

  # Step 3: Save state
  echo ""
  echo "3. Game state"
  SAVE_RESP=$(curl -sf -X PUT "$DEPLOYED_URL/api/v1/game-state" \
    -H "Content-Type: application/json" \
    -H "X-Session-Token: $TOKEN" \
    -d '{"data":{"level":1},"checkpoint":"start"}' --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if [ -n "$SAVE_RESP" ]; then
    echo -e "  ${GREEN}PASS${NC} game state saved"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} game state save"
    ((fail++))
    errors+=("game-state-save")
  fi

  # Step 4: Load state
  LOAD_RESP=$(curl -sf "$DEPLOYED_URL/api/v1/game-state" \
    -H "X-Session-Token: $TOKEN" --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if [ -n "$LOAD_RESP" ]; then
    echo -e "  ${GREEN}PASS${NC} game state loaded"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} game state load"
    ((fail++))
    errors+=("game-state-load")
  fi

  # Step 5: Submit score
  echo ""
  echo "4. Score submission"
  SCORE_RESP=$(curl -sf -X POST "$DEPLOYED_URL/api/v1/scores" \
    -H "Content-Type: application/json" \
    -H "X-Session-Token: $TOKEN" \
    -d '{"score":100}' --max-time "$TIMEOUT" 2>/dev/null || echo "")
  if echo "$SCORE_RESP" | grep -q '"accepted":true'; then
    echo -e "  ${GREEN}PASS${NC} score submitted"
    ((pass++))
  else
    echo -e "  ${RED}FAIL${NC} score submission"
    ((fail++))
    errors+=("score-submit")
  fi

  # Step 6: Leaderboard
  echo ""
  echo "5. Leaderboard"
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
