---
name: prepare-deploy
description: Run pre-deploy checks and prepare the game server for deployment
trigger: user asks to prepare deploy, run preflight, check deployment readiness
---

# prepare-deploy

Run pre-deploy checks and prepare the slot machine game server for deployment.

## Prerequisites

- Go 1.24+ installed and on PATH
- Node.js 18+ installed and on PATH
- The deploy CLI built and verified (run `setup-game-designer-cli` first)
- The game server created (run `create-game-server` first)
- The SDK connected (run `connect-js-sdk` first)

## When to Apply

- The user asks to prepare for deployment
- The user asks to run preflight or pre-deploy checks
- Before deploying to PaaS

## What This Skill Does

1. Verify the deploy CLI binary is available and reports its version
2. Run preflight checks using the CLI:
   ```bash
   cd cli && ./game-designer preflight --server-path ../server-template
   ```
3. Verify the server builds without errors: `cd server-template && GOWORK=off go build ./...`
4. Verify the server tests pass: `cd server-template && GOWORK=off go test ./... -v`
5. Check that the SDK is built: `cd sdk-js && npm run build`
6. Check that the SDK tests pass: `cd sdk-js && npm test`
7. Run local verification: `./scripts/verify-local.sh`
8. Report readiness

## Read Scope

- `server-template/` — Go server source for build verification
- `sdk-js/` — TypeScript SDK source for build verification
- `cli/` — Deploy CLI for preflight execution
- `scripts/verify-local.sh` — Local verification script

## Write Scope

- None — this skill is read-only; it checks but does not modify project files

## Checks

1. Deploy CLI is available (exits 0 with version output)
2. Server builds: `cd server-template && GOWORK=off go build ./...`
3. Server tests pass: `cd server-template && GOWORK=off go test ./... -v`
4. SDK builds: `cd sdk-js && npm run build`
5. SDK tests pass: `cd sdk-js && npm test`
6. CLI preflight passes
7. Local verification passes (includes slot spin loop check)

## Success Output

```
Deployment preparation complete.
- Deploy CLI: available
- Server build: OK
- Server tests: PASS
- SDK build: OK
- SDK tests: PASS
- CLI preflight: PASS
- Local verification: PASS
Ready to deploy.
```

## Failure Output

- Deploy CLI missing: Run `setup-game-designer-cli` to build the CLI first
- Server build failure: Report Go compiler errors
- Test failure: Report which tests failed and what they expect
- SDK build failure: Report TypeScript errors
- Local verification failure: Report which slot endpoint checks failed
- Not ready: List specific issues that must be resolved before deploy
