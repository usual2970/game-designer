# Troubleshooting

Common issues and their resolution paths.

## Contract Mismatch

**Symptoms:** SDK types don't match server responses, 400 errors on valid requests, type errors after contract update.

**Diagnosis:**
```bash
cd contracts && npm run validate
```

**Fix:**
1. Validate the contract
2. Compare contract schemas with Go server handler responses
3. Update the contract or server implementation
4. Re-run `./scripts/verify-local.sh`

## SDK Integration Failure

**Symptoms:** `ApiError` thrown, methods not found, type errors.

**Common causes:**
- Missing session token: Call `createOrResumeSession` before other methods
- Invalid request body: Ensure data matches the TypeScript types
- Import error: Check the SDK package is properly referenced

**Diagnosis:**
```typescript
// Check the error code
catch (error) {
  if (error instanceof ApiError) {
    console.log(error.code, error.message, error.details);
  }
}
```

## Server Error

**Symptoms:** 500 responses, unexpected JSON, missing endpoints.

**Diagnosis:**
```bash
cd server-template && GOWORK=off go test ./... -v
```

**Common causes:**
- Handler returns wrong status code: Check against OpenAPI contract
- Missing field in response: Verify all required fields are present
- Port conflict: Kill existing process on port 8080

## Missing PaaS Configuration

**Symptoms:** CLI exits with `PREFLIGHT_FAILED`, deploy fails immediately.

**Fix:**
1. Run `game-designer preflight --server-path .` to see which checks fail
2. Ensure all required flags are provided (`--app-name`, `--provider`)
3. Check environment variables for provider-specific config

## Deploy Failure

**Symptoms:** CLI exits non-zero, `DEPLOY_FAILED` or `HEALTH_CHECK_FAILED` code.

**Diagnosis:**
- Read the structured JSON output for the specific failure
- Check provider logs for the error details

**Fix:**
- `DEPLOY_FAILED`: Check PaaS credentials and provider configuration
- `HEALTH_CHECK_FAILED`: The service deployed but isn't healthy — check server logs

## Verification Failure

**Symptoms:** `verify-local.sh` or `verify-deployed.sh` exits non-zero.

**Fix:**
1. Read the output for which specific check failed
2. Run the failing component's test suite directly
3. For local: ensure the server is running on the expected port
4. For deployed: ensure the URL is accessible and the server is running
