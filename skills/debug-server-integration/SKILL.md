---
name: debug-server-integration
description: Triage SDK, server, contract, and deployment failures
trigger: user reports integration failure, SDK errors, server errors, deploy issues, or verification failures
---

# debug-server-integration

Triage and diagnose SDK, server, contract, and deployment integration failures for the slot machine backend.

## Prerequisites

- Go 1.24+ installed and on PATH
- Node.js 18+ installed and on PATH
- Access to the project directory with `server-template/`, `sdk-js/`, `contracts/`, and `cli/`

## When to Apply

- The user reports an integration error or unexpected behavior
- SDK calls fail with errors
- Server returns unexpected responses
- Contract or type mismatches
- Deploy or health check failures
- Verification script failures
- Spin failures, balance issues, or payout mismatches

## What This Skill Does

1. Identify the failure category
2. Run targeted diagnostics
3. Report findings and suggest fixes

## Read Scope

- `contracts/` — OpenAPI contract for validation
- `server-template/` — Go server source for test diagnostics
- `sdk-js/` — TypeScript SDK source for test diagnostics
- `cli/` — Deploy CLI for preflight diagnostics
- `scripts/` — Verification scripts

## Write Scope

- None — this skill is read-only; it diagnoses but does not modify project files

## Failure Categories

### Category 1: Contract Mismatch

**Symptoms:** SDK types don't match server responses, 400 errors on valid requests
**Diagnosis:**
```bash
cd contracts && npm run validate
```
Compare the OpenAPI schemas with the Go server handler responses.

**Fix:** Update the contract or align the SDK types with the slot machine API surface.

### Category 2: SDK Integration Failure

**Symptoms:** `ApiError` thrown, type errors, missing methods
**Diagnosis:**
- Check error code: `INVALID_PARAMETERS`, `UNAUTHORIZED`, `NOT_FOUND`, `SESSION_EXPIRED`, `INSUFFICIENT_BALANCE`
- Verify session token is set after `createOrResumeSession`
- Check request body matches SDK types for spin, balance, and config endpoints

**Fix:** Ensure SDK calls follow the golden path pattern from `sdk-js/examples/basic-slot-machine.ts`.

### Category 3: Invalid Wager

**Symptoms:** 400 with `INVALID_PARAMETERS`, wager rejected
**Diagnosis:**
- Check wager is within min/max range from `getSlotConfig()`
- Verify wager is a positive integer
- Check that the wager field is present in the request body

**Fix:** Use a wager value between the configured min and max (default: 1-100).

### Category 4: Insufficient Balance

**Symptoms:** 400 with `INSUFFICIENT_BALANCE`, spin rejected
**Diagnosis:**
- Check current balance via `getBalance()`
- Verify the player has enough virtual credits for the wager
- Check if previous spins drained the balance

**Fix:** Lower the wager or check that the player received the default starting balance.

### Category 5: Payout Mismatch

**Symptoms:** Balance or payout values don't match expected results
**Diagnosis:**
- Spin resolution is server-authoritative — the server's result is the source of truth
- Check payline win evaluation matches the slot configuration
- Verify balance = previous balance - wager + total payout

**Fix:** The client should trust server-returned spin results, not calculate payouts locally.

### Category 6: Server Error

**Symptoms:** 500 errors, unexpected JSON, missing endpoints
**Diagnosis:**
```bash
cd server-template && GOWORK=off go test ./... -v
```
Run the HTTP integration tests to verify server behavior.

**Fix:** Check handler implementation matches the contract schema for slot endpoints.

### Category 7: Deploy Failure

**Symptoms:** CLI exits non-zero, health check fails
**Diagnosis:**
```bash
cd cli && go run ./cmd/game-designer preflight --server-path ../server-template
```
Check preflight output for specific failure points.

**Fix:** Use `prepare-deploy` to resolve issues before retrying.

### Category 8: Verification Failure

**Symptoms:** `verify-local.sh` or `verify-deployed.sh` fails
**Diagnosis:**
- Read the script output for which slot endpoint or check failed
- Check if the server is running locally on `:8080`
- Check if the deployed URL is accessible
- Verify the slot spin loop completes: session → config → balance → spin → history → leaderboard

**Fix:** Address the specific endpoint failure, then re-run verification.

## Checks

1. Contract validates: `cd contracts && npm run validate`
2. Server tests pass: `cd server-template && GOWORK=off go test ./...`
3. SDK tests pass: `cd sdk-js && npm test`
4. Server responds locally: `curl -X POST http://localhost:8080/api/v1/session -d '{"playerId":"test"}'`
5. Slot config available: `curl http://localhost:8080/api/v1/slot/config -H 'X-Session-Token: <token>'`

## Success Output

```
Integration issue diagnosed.
- Category: [contract/sdk/wager/balance/payout/server/deploy/verification]
- Root cause: [description]
- Suggested fix: [action]
```

## Failure Output

If diagnosis cannot determine root cause:
- Report all diagnostic results
- Suggest running all tests: `./scripts/verify-local.sh`
- Suggest running plugin package validation: `./scripts/verify-plugin-package.sh`
- Ask the user for specific error messages or logs
