# Local Verification

Run the local verification script before deploying to PaaS. For plugin package checks, use `./scripts/verify-plugin-package.sh` instead.

## Related Verification Scripts

| Script | Purpose |
|--------|---------|
| `./scripts/verify-plugin-package.sh` | Validate plugin manifests, skills, and bundled assets |
| `./scripts/verify-local.sh` | Build, test, and activity loop against local server |
| `./scripts/verify-deployed.sh <url>` | Activity loop against deployed server |

## Quick Run

```bash
# Start the server in a separate terminal
cd server-template && GOWORK=off go run ./cmd/server

# Run verification
./scripts/verify-local.sh
```

With a custom server URL:

```bash
./scripts/verify-local.sh http://localhost:9090
```

## What It Checks

1. **Contract validation** — OpenAPI schema is valid
2. **Server build** — Go server compiles without errors
3. **Server tests** — All Go unit and integration tests pass
4. **SDK build** — TypeScript SDK compiles
5. **SDK tests** — SDK test suite passes
6. **CLI preflight** — Deploy CLI preflight checks pass
7. **Server connectivity** — Server is reachable (if running)
8. **Activity loop** — Full golden path through live server (if running)

## Output

On success:
```json
{"success":true,"message":"All local verification checks passed","code":"SUCCESS","details":{"passed":12}}
```

On failure:
```json
{"success":false,"message":"Local verification failed","code":"VERIFICATION_FAILED","details":{"passed":8,"failed":2}}
```

## Skipping Live Checks

The activity loop checks (steps 7-8) are skipped if the server is not running. This is normal when running CI without a server — the build, test, and contract checks still provide meaningful coverage.

## Troubleshooting

- **Contract validation fails**: Run `cd contracts && npm run validate` to see details
- **Server build fails**: Check Go version (1.24+) and run `GOWORK=off go build ./...` in `server-template/`
- **SDK build fails**: Run `cd sdk-js && npm run build` for TypeScript errors
