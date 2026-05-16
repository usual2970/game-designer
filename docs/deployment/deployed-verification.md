# Deployed Verification

Run the deployed verification script after deploying to PaaS.

## Quick Run

```bash
./scripts/verify-deployed.sh https://my-game.fake.local
```

## What It Checks

1. **Endpoint reachability** — Deployed URL responds
2. **Session creation** — Player can create a session
3. **Slot config retrieval** — Slot machine configuration is returned
4. **Balance check** — Player virtual credit balance is returned
5. **Spin execution** — Spin with wager is accepted, result and payout returned
6. **Spin history** — Recent spin results are returned
7. **Slot leaderboard** — Top payouts ranking is returned

## Output

On success:
```json
{"success":true,"message":"All deployed verification checks passed","code":"SUCCESS","details":{"passed":7}}
```

On failure:
```json
{"success":false,"message":"Deployed verification failed","code":"VERIFICATION_FAILED","details":{"passed":4,"failed":3}}
```

## Failure Categories

| Failure | Cause | Action |
|---------|-------|--------|
| Endpoint unreachable | DNS, network, or server down | Check PaaS deployment status |
| Session creation failed | Server misconfiguration | Check server logs and env vars |
| Slot config failed | Config endpoint missing or malformed | Check slot config handler, verify reel/payline data |
| Balance check failed | Balance service error | Check virtual credits initialization |
| Spin failed | Validation error or INSUFFICIENT_BALANCE | Check wager validation, ensure initial balance is funded |
| Spin history failed | History storage issue | Check storage configuration |
| Slot leaderboard failed | Ranking engine error | Check leaderboard service |

## CI Integration

```bash
# After deploy
./scripts/verify-deployed.sh "$DEPLOYED_URL"
if [ $? -ne 0 ]; then
  echo "Deployed verification failed — see errors above"
  exit 1
fi
```
