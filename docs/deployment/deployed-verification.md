# Deployed Verification

Run the deployed verification script after deploying to PaaS.

## Quick Run

```bash
./scripts/verify-deployed.sh https://my-game.fake.local
```

## What It Checks

1. **Endpoint reachability** — Deployed URL responds
2. **Session creation** — Player can create a session
3. **Game state save** — Progress can be saved
4. **Game state load** — Progress can be resumed
5. **Score submission** — Scores are accepted
6. **Leaderboard read** — Rankings are returned

## Output

On success:
```json
{"success":true,"message":"All deployed verification checks passed","code":"SUCCESS","details":{"passed":6}}
```

On failure:
```json
{"success":false,"message":"Deployed verification failed","code":"VERIFICATION_FAILED","details":{"passed":3,"failed":3}}
```

## Failure Categories

| Failure | Cause | Action |
|---------|-------|--------|
| Endpoint unreachable | DNS, network, or server down | Check PaaS deployment status |
| Session creation failed | Server misconfiguration | Check server logs and env vars |
| Game state failed | Persistence layer issue | Check storage configuration |
| Score submission failed | Validation or state error | Check server logs |
| Leaderboard failed | Ranking engine error | Check leaderboard service |

## CI Integration

```bash
# After deploy
./scripts/verify-deployed.sh "$DEPLOYED_URL"
if [ $? -ne 0 ]; then
  echo "Deployed verification failed — see errors above"
  exit 1
fi
```
