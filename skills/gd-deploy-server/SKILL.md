---
name: gd-deploy-server
description: Deploy the game server to the team PaaS
trigger: user asks to deploy the game server, push to PaaS, release the backend
---

# gd-deploy-server

Deploy the slot machine game server to the team PaaS using the CLI.

## Prerequisites

- The deploy CLI built and verified (run `gd-setup-cli` first)
- Local verification passed (run `gd-prepare-deploy` first)
- For production: `GD_IDENTIFIER` and `GD_PASSWORD` environment variables set

## When to Apply

- The user asks to deploy the game server
- The user asks to push to PaaS or release the backend
- After local verification passes

## What This Skill Does

### Dry Run (Fake Provider)

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
4. Report the deployment result

### Production Deploy (3os Provider)

1. Verify the deploy CLI binary is available
2. Verify `GD_IDENTIFIER` and `GD_PASSWORD` environment variables are set
3. List existing games (optional):
   ```bash
   cd cli && ./game-designer deploy \
     --provider 3os \
     --mode list \
     --identifier "$GD_IDENTIFIER" \
     --password "$GD_PASSWORD"
   ```
4. Deploy with the appropriate mode:
   ```bash
   # Create a new game
   cd cli && ./game-designer deploy \
     --provider 3os \
     --mode create \
     --identifier "$GD_IDENTIFIER" \
     --password "$GD_PASSWORD" \
     --game-name "<game-name>" \
     --package-path <path-to-zip> \
     --version <version> \
     --change-log "<description>" \
     --screen-type 1 \
     --half-support 2 \
     --half-ratio "0.75" \
     --backend-dir "<dir>" \
     --backend-cmd "<cmd>" \
     --frontend-dir "<dir>" \
     --socket-dir "<dir>" \
     --socket-cmd "<cmd>"

   # Update version for existing game
   cd cli && ./game-designer deploy \
     --provider 3os \
     --mode update-version \
     --identifier "$GD_IDENTIFIER" \
     --password "$GD_PASSWORD" \
     --game-uri <game-uri> \
     --package-path <path-to-zip> \
     --version <version> \
     --change-log "<description>"
   ```
5. Parse the structured JSON output
6. If successful and review is needed:
   ```bash
   cd cli && ./game-designer deploy \
     --provider 3os \
     --mode apply-review \
     --identifier "$GD_IDENTIFIER" \
     --password "$GD_PASSWORD" \
     --review-uri <review-uri>
   ```

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
- Provider: 3os (or fake)
- Mode: create
- Game URI: <game-uri>
- URL: https://...
- Version: 1.0.0
- Health: OK
- Review: applied (or skipped)
```

## Failure Output

- Deploy CLI missing: Run `gd-setup-cli` to build the CLI first
- `CONFIG_ERROR`: Check required flags for the selected mode
- `PREFLIGHT_FAILED`: Run `gd-prepare-deploy` to fix pre-deploy issues
- `AUTH_FAILED`: Verify `GD_IDENTIFIER` / `GD_PASSWORD` credentials
- `UPLOAD_FAILED`: Check package file exists at `--package-path`
- `PUBLISH_FAILED`: Check game payload and API response details
- `REVIEW_FAILED`: Check review URI and game state
- `PARTIAL_SUCCESS`: Game published but review failed — retry review separately
- `DEPLOY_FAILED`: Check PaaS provider logs and configuration
- `HEALTH_CHECK_FAILED`: The service deployed but is not responding correctly
- `INTERNAL_ERROR`: Retry once. If persistent, check CLI version and provider config

## Error Recovery

- If deploy CLI missing: use `gd-setup-cli` skill to build it
- If `PREFLIGHT_FAILED`: use `gd-prepare-deploy` skill to resolve issues
- If `AUTH_FAILED`: verify credentials with the user
- If `UPLOAD_FAILED`: ensure package file exists and is readable
- If `PUBLISH_FAILED`/`REVIEW_FAILED`: use `gd-debug-integration` skill to diagnose
- If `PARTIAL_SUCCESS`: game is published, retry `--mode apply-review` separately
