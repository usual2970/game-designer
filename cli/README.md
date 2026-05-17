# Game Designer Deploy CLI

Go CLI for deploying the Game Designer Server to PaaS.

## Build

```bash
cd cli
GOWORK=off go build -o game-designer ./cmd/game-designer
```

Agents using the Game Designer plugin should run the `gd-setup-cli` skill instead of building manually. That skill handles Go version checks, build, and binary verification in one step.

Plugin installation does not compile this binary. Re-run `gd-setup-cli` after updating the plugin or changing files under `cli/`.

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

#### Dry Run (Offline Verification)

```bash
game-designer deploy \
  --server-path ../server-template \
  --app-name my-game \
  --env production \
  --provider fake
```

#### Production Deploy

```bash
# List your games
game-designer deploy \
  --provider 3os \
  --mode list \
  --identifier $GD_IDENTIFIER \
  --password $GD_PASSWORD

# Create a new game with initial version
game-designer deploy \
  --provider 3os \
  --mode create \
  --identifier $GD_IDENTIFIER \
  --password $GD_PASSWORD \
  --game-name "My Game" \
  --game-desc "Game description" \
  --package-path ./game.zip \
  --sql-path ./init.sql \
  --version 1.0.0 \
  --change-log "Initial release" \
  --screen-type 1 \
  --half-support 2 \
  --half-ratio "0.75" \
  --backend-dir "admin" \
  --backend-cmd "./server -type admin" \
  --frontend-dir "h5" \
  --socket-dir "logic" \
  --socket-cmd "./server -type logic"

# Update game base info
game-designer deploy \
  --provider 3os \
  --mode update-info \
  --identifier $GD_IDENTIFIER \
  --password $GD_PASSWORD \
  --game-uri <game-uri> \
  --game-desc "Updated description"

# Publish a new version
game-designer deploy \
  --provider 3os \
  --mode update-version \
  --identifier $GD_IDENTIFIER \
  --password $GD_PASSWORD \
  --game-uri <game-uri> \
  --package-path ./game-v2.zip \
  --version 1.1.0 \
  --change-log "Bug fixes"

# Submit for review
game-designer deploy \
  --provider 3os \
  --mode apply-review \
  --identifier $GD_IDENTIFIER \
  --password $GD_PASSWORD \
  --review-uri <review-uri>
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
  "message": "Deployed to https://...",
  "code": "SUCCESS",
  "details": { ... }
}
```

## Environment Variables

Use environment variables for credentials instead of flags to avoid shell history leaks:

| Variable | Flag | Description |
|----------|------|-------------|
| `GD_IDENTIFIER` | `--identifier` | Login identifier (email) |
| `GD_PASSWORD` | `--password` | Login password |
| `GD_BASE_URL` | `--base-url` | API base URL (default: `https://api.3sdk.yu3.co`) |
| `GD_GAME_URI` | `--game-uri` | Game URI for update operations |

## Publish Modes

| Mode | Required Flags | Description |
|------|---------------|-------------|
| `create` | `--game-name`, `--package-path`, `--version` | Create a new game with initial version |
| `update-info` | `--game-uri` | Update game base info (name, description, logo) |
| `update-version` | `--game-uri`, `--package-path`, `--version` | Publish a new version of an existing game |
| `list` | (auth only) | List your developer games with pagination |
| `apply-review` | `--review-uri` | Submit a game review application |

## Error Codes

| Code | Meaning | Agent Action |
|------|---------|-------------|
| `SUCCESS` | Deploy succeeded | Continue |
| `PREFLIGHT_FAILED` | Pre-deploy checks failed | Fix issues and retry |
| `DEPLOY_FAILED` | Provider deploy error | Check provider details |
| `HEALTH_CHECK_FAILED` | Post-deploy health check failed | Investigate endpoint |
| `INTERNAL_ERROR` | Unexpected error | Retry once, then report |
| `CONFIG_ERROR` | Missing or invalid CLI configuration | Check required flags/env vars |
| `AUTH_FAILED` | Login credentials rejected | Verify identifier/password |
| `LIST_FAILED` | Game list API error | Check auth and network |
| `LOOKUP_FAILED` | Game lookup returned zero or multiple matches | Use explicit `--game-uri` |
| `UPLOAD_FAILED` | OSS upload failed | Check package path and policy token |
| `PUBLISH_FAILED` | Game create/update API rejected | Check payload and game state |
| `REVIEW_FAILED` | Review application failed | Check review URI and game state |
| `PARTIAL_SUCCESS` | Publish succeeded but review failed | Check review state, retry review |

## Test

```bash
cd cli
GOWORK=off go test ./... -v
```
