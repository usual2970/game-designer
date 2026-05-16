---
name: prepare-deploy
description: Run pre-deploy checks and prepare the game server for deployment
trigger: user asks to prepare deploy, run preflight, check deployment readiness
---

# prepare-deploy

Run pre-deploy checks and prepare the game server for deployment.

## When to Apply

- The user asks to prepare for deployment
- The user asks to run preflight or pre-deploy checks
- Before deploying to PaaS

## What This Skill Does

1. Run preflight checks using the CLI:
   ```bash
   cd cli && go run ./cmd/game-designer preflight --server-path ../server-template
   ```
2. Verify the server builds without errors
3. Check that the SDK is built: `cd sdk-js && npm run build`
4. Run local verification: `./scripts/verify-local.sh`
5. Report readiness

## Files Changed

- None (read-only checks)

## Checks

1. Server builds: `cd server-template && GOWORK=off go build ./...`
2. Server tests pass: `cd server-template && GOWORK=off go test ./... -v`
3. SDK builds: `cd sdk-js && npm run build`
4. SDK tests pass: `cd sdk-js && npm test`
5. CLI preflight passes
6. Local verification passes

## Success Output

```
Deployment preparation complete.
- Server build: OK
- Server tests: PASS (21 tests)
- SDK build: OK
- SDK tests: PASS (8 tests)
- CLI preflight: PASS
- Local verification: PASS
Ready to deploy.
```

## Failure Output

- Server build failure: Report Go compiler errors
- Test failure: Report which tests failed and what they expect
- SDK build failure: Report TypeScript errors
- Local verification failure: Report which endpoint checks failed
- Not ready: List specific issues that must be resolved before deploy
