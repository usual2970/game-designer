---
name: debug-server-integration
description: Triage SDK, server, contract, and deployment failures
trigger: user reports integration failure, SDK errors, server errors, deploy issues, or verification failures
---

# debug-server-integration

Triage and diagnose SDK, server, contract, and deployment integration failures.

## When to Apply

- The user reports an integration error or unexpected behavior
- SDK calls fail with errors
- Server returns unexpected responses
- Contract or type mismatches
- Deploy or health check failures
- Verification script failures

## What This Skill Does

1. Identify the failure category
2. Run targeted diagnostics
3. Report findings and suggest fixes

## Failure Categories

### Category 1: Contract Mismatch

**Symptoms:** SDK types don't match server responses, 400 errors on valid requests
**Diagnosis:**
```bash
cd contracts && npm run validate
```
Compare the OpenAPI schemas with the Go server handler responses.

**Fix:** Update the contract or regenerate/align the SDK types.

### Category 2: SDK Integration Failure

**Symptoms:** `ApiError` thrown, type errors, missing methods
**Diagnosis:**
- Check error code: `INVALID_PARAMETERS`, `UNAUTHORIZED`, `NOT_FOUND`, `SESSION_EXPIRED`
- Verify session token is set after `createOrResumeSession`
- Check request body matches SDK types

**Fix:** Ensure SDK calls follow the golden path pattern from examples.

### Category 3: Server Error

**Symptoms:** 500 errors, unexpected JSON, missing endpoints
**Diagnosis:**
```bash
cd server-template && GOWORK=off go test ./... -v
```
Run the HTTP integration tests to verify server behavior.

**Fix:** Check handler implementation matches the contract schema.

### Category 4: Deploy Failure

**Symptoms:** CLI exits non-zero, health check fails
**Diagnosis:**
```bash
cd cli && go run ./cmd/game-designer preflight --server-path ../server-template
```
Check preflight output for specific failure points.

**Fix:** Use `prepare-deploy` to resolve issues before retrying.

### Category 5: Verification Failure

**Symptoms:** `verify-local.sh` or `verify-deployed.sh` fails
**Diagnosis:**
- Read the script output for which endpoint or check failed
- Check if the server is running locally on `:8080`
- Check if the deployed URL is accessible

**Fix:** Address the specific endpoint failure, then re-run verification.

## Checks

1. Contract validates: `cd contracts && npm run validate`
2. Server tests pass: `cd server-template && GOWORK=off go test ./...`
3. SDK tests pass: `cd sdk-js && npm test`
4. Server responds locally: `curl -X POST http://localhost:8080/api/v1/session -d '{"playerId":"test"}'`

## Success Output

```
Integration issue diagnosed.
- Category: [contract/sdk/server/deploy/verification]
- Root cause: [description]
- Suggested fix: [action]
```

## Failure Output

If diagnosis cannot determine root cause:
- Report all diagnostic results
- Suggest running all tests: `./scripts/verify-local.sh`
- Ask the user for specific error messages or logs
