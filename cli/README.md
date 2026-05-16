# Game Designer Deploy CLI

Go CLI for deploying the Game Designer Server to PaaS.

## Build

```bash
cd cli
GOWORK=off go build -o game-designer ./cmd/game-designer
```

Agents using the Game Designer plugin should run the `setup-game-designer-cli` skill instead of building manually. That skill handles Go version checks, build, and binary verification in one step.

## Commands

### Version

```bash
game-designer version
```

### Preflight Checks

```bash
game-designer preflight --server-path ../server-template
```

Validates:
- Server path exists
- Server builds successfully

### Deploy

```bash
game-designer deploy \
  --server-path ../server-template \
  --app-name my-game \
  --env production \
  --provider fake
```

Runs the full lifecycle:
1. Preflight checks
2. Deploy through provider
3. Health check
4. Status verification

Output is structured JSON for agent parsing:

```json
{
  "success": true,
  "message": "Deployed to https://my-game.fake.local",
  "code": "SUCCESS",
  "details": { ... }
}
```

## Error Codes

| Code | Meaning | Agent Action |
|------|---------|-------------|
| `SUCCESS` | Deploy succeeded | Continue |
| `PREFLIGHT_FAILED` | Pre-deploy checks failed | Fix issues and retry |
| `DEPLOY_FAILED` | Provider deploy error | Check provider details |
| `HEALTH_CHECK_FAILED` | Post-deploy health check failed | Investigate endpoint |
| `INTERNAL_ERROR` | Unexpected error | Retry once, then report |

## Test

```bash
cd cli
GOWORK=off go test ./... -v
```
