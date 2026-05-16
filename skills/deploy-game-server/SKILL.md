---
name: deploy-game-server
description: Deploy the game server to the team PaaS
trigger: user asks to deploy the game server, push to PaaS, release the backend
---

# deploy-game-server

Deploy the slot machine game server to the team PaaS using the CLI.

## Prerequisites

- The deploy CLI built and verified (run `setup-game-designer-cli` first)
- Local verification passed (run `prepare-deploy` first)
- PaaS provider credentials configured if not using the fake provider

## When to Apply

- The user asks to deploy the game server
- The user asks to push to PaaS or release the backend
- After local verification passes

## What This Skill Does

1. Verify the deploy CLI binary is available and reports its version
2. Run the deploy command:
   ```bash
   cd cli && ./game-designer deploy \
     --server-path ../server-template \
     --app-name <game-name> \
     --env production \
     --provider fake
   ```
3. Parse the structured JSON output
4. If successful, run deployed verification: `./scripts/verify-deployed.sh <url>`
5. Report the deployment result

## Read Scope

- `server-template/` — Go server source for deployment
- `cli/` — Deploy CLI for execution
- `scripts/verify-deployed.sh` — Deployed verification script

## Write Scope

- None — this skill triggers deployment; it does not modify local project files

## Checks

1. Deploy CLI is available (exits 0 with version output)
2. CLI exits with code 0
3. JSON output contains `"success": true`
4. Deployed URL is accessible
5. Health check passes
6. Deployed slot spin loop verification passes

## Success Output

```
Deployment successful.
- URL: https://<app-name>.fake.local
- Version: v0.1.0
- Health: OK
- Deployed verification: PASS (slot spin loop verified)
```

## Failure Output

- Deploy CLI missing: Run `setup-game-designer-cli` to build the CLI first
- `PREFLIGHT_FAILED`: Run `prepare-deploy` to fix pre-deploy issues
- `DEPLOY_FAILED`: Check PaaS provider logs and configuration
- `HEALTH_CHECK_FAILED`: The service deployed but is not responding correctly
- `INTERNAL_ERROR`: Retry once. If persistent, check CLI version and provider config

## Error Recovery

- If deploy CLI missing: use `setup-game-designer-cli` skill to build it
- If `PREFLIGHT_FAILED`: use `prepare-deploy` skill to resolve issues
- If `DEPLOY_FAILED`: check PaaS credentials and provider configuration
- If `HEALTH_CHECK_FAILED`: use `debug-server-integration` skill to diagnose
