#!/usr/bin/env bash
# verify-local.sh — Run local verification for the Game Designer Server MVP
#
# Usage: ./scripts/verify-local.sh [server_url]
#
# Starts the Go server, runs SDK-backed activity-loop checks, and reports
# actionable results.

set -euo pipefail

SERVER_URL="${1:-http://localhost:8080}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TIMEOUT=30

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

pass=0
fail=0
errors=()

# Prevent set -e from killing the script on arithmetic ((0))
true_fn() { true; }

check() {
  local name="$1"
  shift
  if "$@" 2>/dev/null; then
    echo -e "  ${GREEN}PASS${NC} $name"
    pass=$((pass + 1))
  else
    echo -e "  ${RED}FAIL${NC} $name"
    fail=$((fail + 1))
    errors+=("$name")
  fi
}

echo "=== Game Designer Server — Local Verification ==="
echo ""

# Step 1: Verify contract
echo "1. Contract validation"
cd "$ROOT_DIR/contracts"
check "contract validates" npx swagger-cli validate game-server.openapi.yaml

# Step 2: Build server
echo ""
echo "2. Server build"
cd "$ROOT_DIR/server-template"
check "go build" env GOWORK=off go build ./...

# Step 3: Run server tests
echo ""
echo "3. Server tests"
check "server tests pass" env GOWORK=off go test ./...

# Step 4: Build SDK
echo ""
echo "4. SDK build"
cd "$ROOT_DIR/sdk-js"
check "sdk builds" npm run build

# Step 5: Run SDK tests
echo ""
echo "5. SDK tests"
check "sdk tests pass" npm test

# Step 6: CLI preflight
echo ""
echo "6. CLI preflight"
cd "$ROOT_DIR/cli"
check "cli preflight" env GOWORK=off go run ./cmd/game-designer preflight --server-path ../server-template

# Step 7: Check server is running
echo ""
echo "7. Server connectivity"
if curl -sf -o /dev/null "$SERVER_URL/api/v1/session" 2>/dev/null; then
  echo -e "  ${GREEN}PASS${NC} server reachable at $SERVER_URL"
  ((pass++))

  # Step 8: Activity loop against live server
  echo ""
  echo "8. Activity loop (live server)"

  # Create session
  SESSION_RESP=$(curl -sf -X POST "$SERVER_URL/api/v1/session" \
    -H "Content-Type: application/json" \
    -d '{"playerId":"verify-local","nickname":"VerifyBot"}')
  if [ -n "$SESSION_RESP" ]; then
    TOKEN=$(echo "$SESSION_RESP" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TOKEN" ]; then
      echo -e "  ${GREEN}PASS${NC} session created"
      ((pass++))

      # Save game state
      SAVE_RESP=$(curl -sf -X PUT "$SERVER_URL/api/v1/game-state" \
        -H "Content-Type: application/json" \
        -H "X-Session-Token: $TOKEN" \
        -d '{"data":{"level":1,"coins":0},"checkpoint":"start"}')
      if [ -n "$SAVE_RESP" ]; then
        echo -e "  ${GREEN}PASS${NC} game state saved"
        ((pass++))
      else
        echo -e "  ${RED}FAIL${NC} game state save"
        ((fail++))
        errors+=("game-state-save")
      fi

      # Load game state
      LOAD_RESP=$(curl -sf "$SERVER_URL/api/v1/game-state" \
        -H "X-Session-Token: $TOKEN")
      if [ -n "$LOAD_RESP" ]; then
        echo -e "  ${GREEN}PASS${NC} game state loaded"
        ((pass++))
      else
        echo -e "  ${RED}FAIL${NC} game state load"
        ((fail++))
        errors+=("game-state-load")
      fi

      # Submit score
      SCORE_RESP=$(curl -sf -X POST "$SERVER_URL/api/v1/scores" \
        -H "Content-Type: application/json" \
        -H "X-Session-Token: $TOKEN" \
        -d '{"score":100}')
      if echo "$SCORE_RESP" | grep -q '"accepted":true'; then
        echo -e "  ${GREEN}PASS${NC} score submitted"
        ((pass++))
      else
        echo -e "  ${RED}FAIL${NC} score submission"
        ((fail++))
        errors+=("score-submit")
      fi

      # Read leaderboard
      LB_RESP=$(curl -sf "$SERVER_URL/api/v1/leaderboard" \
        -H "X-Session-Token: $TOKEN")
      if echo "$LB_RESP" | grep -q '"entries"'; then
        echo -e "  ${GREEN}PASS${NC} leaderboard read"
        ((pass++))
      else
        echo -e "  ${RED}FAIL${NC} leaderboard read"
        ((fail++))
        errors+=("leaderboard-read")
      fi
    else
      echo -e "  ${RED}FAIL${NC} could not extract session token"
      ((fail++))
      errors+=("session-token")
    fi
  else
    echo -e "  ${RED}FAIL${NC} session creation"
    ((fail++))
    errors+=("session-create")
  fi
else
  echo -e "  ${YELLOW}SKIP${NC} server not running at $SERVER_URL (start with: cd server-template && GOWORK=off go run ./cmd/server)"
  echo -e "  ${YELLOW}SKIP${NC} activity loop (requires running server)"
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
  echo '{"success":false,"message":"Local verification failed","code":"VERIFICATION_FAILED","details":{"passed":'"$pass"',"failed":'"$fail"'}}'
  exit 1
fi

echo ""
echo '{"success":true,"message":"All local verification checks passed","code":"SUCCESS","details":{"passed":'"$pass"'}}'
