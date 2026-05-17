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

## Invalid Wager

**Symptoms:** Spin returns 400 with `INVALID_PARAMETERS`, wager rejected.

**Common causes:**
- Wager amount below `minWager` from slot config
- Wager amount above `maxWager` from slot config
- Wager is not a positive integer
- Wager amount exceeds player balance

**Fix:**
1. Call `getSlotConfig()` to retrieve valid wager range
2. Call `getBalance()` to confirm sufficient credits
3. Clamp the wager between `minWager` and `maxWager`
4. Disable spin button when balance is zero

## Insufficient Balance

**Symptoms:** Spin returns 400 with `INSUFFICIENT_BALANCE`, player cannot spin.

**Common causes:**
- Player has exhausted all virtual credits
- Previous spin deducted more credits than expected
- Balance not refreshed after a losing spin

**Fix:**
1. Call `getBalance()` before each spin to display current credits
2. Show a "no credits" state when balance is below `minWager`
3. If credits should reset, check the server-side balance initialization logic
4. Verify the spin endpoint correctly deducts the wager and adds payouts

## Payout Mismatch

**Symptoms:** Payout amount differs from expected value based on payline symbols, balance change does not match wager minus payout.

**Common causes:**
- Payline calculation does not match the slot config payline definitions
- Symbol multiplier table is out of sync between client and server
- Client displays wrong payout due to stale slot config cache

**Fix:**
1. Call `getSlotConfig()` to refresh the payline and symbol definitions
2. Cross-check the server payout calculation with the config multipliers
3. Verify the balance after spin equals `previousBalance - wager + payout`
4. Check server logs for the spin result details and payline evaluation

## Server Error

**Symptoms:** 500 responses, unexpected JSON, missing endpoints.

**Diagnosis:**
```bash
# In the plugin repository:
cd server-template && GOWORK=off go test ./... -v

# In a consuming game project:
cd server && GOWORK=off go test ./... -v
```

**Common causes:**
- Handler returns wrong status code: Check against OpenAPI contract
- Missing field in response: Verify all required fields are present
- Port conflict: Kill existing process on port 8080

## Missing PaaS Configuration

**Symptoms:** CLI exits with `PREFLIGHT_FAILED` or `CONFIG_ERROR`, deploy fails immediately.

**Fix:**
1. Run `game-designer preflight --server-path .` to see which checks fail (`../server` in consuming projects)
2. For production: set `GD_IDENTIFIER` and `GD_PASSWORD` environment variables
3. Check that required mode-specific flags are provided (e.g., `--package-path` for create)

## Deploy Failure

**Symptoms:** CLI exits non-zero with a structured error code.

**Diagnosis:**
- Read the structured JSON output for the specific failure code

**Recovery by code:**

| Code | Cause | Action |
|------|-------|--------|
| `CONFIG_ERROR` | Missing flags or invalid mode | Check required flags for your mode |
| `AUTH_FAILED` | Login rejected | Verify `GD_IDENTIFIER` / `GD_PASSWORD` |
| `LIST_FAILED` | Game list API error | Check auth token and network |
| `LOOKUP_FAILED` | Game lookup ambiguous | Use explicit `--game-uri` instead of name lookup |
| `UPLOAD_FAILED` | OSS upload rejected | Check package file exists and policy token is valid |
| `PUBLISH_FAILED` | Create/update API rejected | Check game payload, version format, and game state |
| `REVIEW_FAILED` | Review application rejected | Verify review URI and that game has a deployable version |
| `PARTIAL_SUCCESS` | Publish OK, review failed | Game is published. Retry review with `--mode apply-review` |
| `DEPLOY_FAILED` | Generic deploy error | Check full error message for details |
| `HEALTH_CHECK_FAILED` | Service deployed but unhealthy | Check server logs at the deployed URL |

## Production Auth Issues

**Symptoms:** `AUTH_FAILED` immediately after deploy starts.

**Common causes:**
- Wrong identifier or password
- Account not registered as developer
- Account suspended

**Fix:**
1. Verify credentials by calling the auth API directly:
   ```bash
   curl -X POST https://api.3sdk.yu3.co/common/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"identifier":"<email>","type":"password","data":"<password>"}'
   ```
2. Ensure the response contains `accessToken`
3. If `code` is non-zero, check the error message

## Upload Issues

**Symptoms:** `UPLOAD_FAILED` after successful auth.

**Common causes:**
- Package file not found at `--package-path`
- Policy token expired (tokens are short-lived)
- OSS rejected the upload (size limit, permission)

**Fix:**
1. Verify the file exists: `ls -la <package-path>`
2. For SQL files, note that `--sql-path` is optional
3. Re-run the deploy command (gets a fresh policy token)

## Partial Publish/Review Failure

**Symptoms:** Game was created or updated but review submission failed.

**Recovery:**
1. The game data is saved — no need to re-create
2. Run `--mode list` to find the game and its review URI
3. Use `--mode apply-review --review-uri <uri>` to retry review submission

## Verification Failure

**Symptoms:** `verify-local.sh` or `verify-deployed.sh` exits non-zero.

**Fix:**
1. Read the output for which specific check failed
2. Run the failing component's test suite directly
3. For local: ensure the server is running on the expected port
4. For deployed: ensure the URL is accessible and the server is running
