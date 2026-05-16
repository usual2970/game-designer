---
name: deploy-game-server
description: Deploy the game server to the team PaaS
trigger: user asks to deploy the game server, push to PaaS, release the backend
---

# deploy-game-server

Deploy the game server to the team PaaS using the CLI.

## When to Apply

- The user asks to deploy the game server
- The user asks to push to PaaS or release the backend
- After local verification passes

## What This Skill Does

1. Run the deploy command:
   ```bash
   cd cli && go run ./cmd/game-designer deploy \
     --server-path ../server-template \
     --app-name <game-name> \
     --env production \
     --provider fake
   ```
2. Parse the structured JSON output
3. If successful, run deployed verification: `./scripts/verify-deployed.sh <url>`
4. Report the deployment result

## Files Changed

- None (triggers deployment)

## Checks

1. CLI exits with code 0
2. JSON output contains `"success": true`
3. Deployed URL is accessible
4. Health check passes

## Success Output

```
Deployment successful.
- URL: https://<app-name>.fake.local
- Version: v0.1.0
- Health: OK
- Deployed verification: PASS
```

## Failure Output

- `PREFLIGHT_FAILED`: Run `prepare-deploy` to fix pre-deploy issues
- `DEPLOY_FAILED`: Check PaaS provider logs and configuration
- `HEALTH_CHECK_FAILED`: The service deployed but is not responding correctly
- `INTERNAL_ERROR`: Retry once. If persistent, check CLI version and provider config

## Error Recovery

- If `PREFLIGHT_FAILED`: use `prepare-deploy` skill to resolve issues
- If `DEPLOY_FAILED`: check PaaS credentials and provider configuration
- If `HEALTH_CHECK_FAILED`: use `debug-server-integration` skill to diagnose
