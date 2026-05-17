# PaaS Provider

The deploy CLI uses a provider interface to abstract deployment targets.

## Architecture

```
Provider Interface
├── Name()       — Provider identifier
├── Deploy()     — Deploy the server
├── Status()     — Check deployment status
└── HealthCheck() — Verify endpoint health
```

## Available Providers

### Fake Provider (Offline/Dry Run)

```bash
game-designer deploy --provider fake
```

The fake provider simulates deploy, status, and health checks without touching a real PaaS. Use for:
- Local testing and smoke checks
- Verifying CLI flag parsing
- Agent dry-run validation

### 3os Provider (Production)

```bash
game-designer deploy --provider 3os --mode create \
  --identifier $GD_IDENTIFIER --password $GD_PASSWORD \
  --game-name "My Game" --package-path ./game.zip --version 1.0.0
```

The 3os provider publishes games through the production PaaS API. It supports five modes:

| Mode | Description |
|------|-------------|
| `create` | Create a new game with an initial version |
| `update-info` | Update game base information |
| `update-version` | Publish a new version of an existing game |
| `list` | List developer's games (paginated) |
| `apply-review` | Submit a game for review |

#### Production Workflow

1. **Auth** — Login with identifier/password, receive `accessToken`
2. **Upload** — Get OSS policy token, upload package (and optional SQL) via POST V4
3. **Publish** — Call game create/update API with uploaded URLs
4. **Review** — Optionally submit the game for review

#### Configuration

| Setting | Flag | Env Var | Default |
|---------|------|---------|---------|
| API base URL | `--base-url` | `GD_BASE_URL` | `https://api.3sdk.yu3.co` |
| Identifier | `--identifier` | `GD_IDENTIFIER` | — |
| Password | `--password` | `GD_PASSWORD` | — |

Credentials should be set via environment variables, not command-line flags, to avoid shell history exposure.

#### API Dependencies

The 3os provider calls these production endpoints:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/common/v1/auth/login` | POST | Password login |
| `/developer/v1/file/policy-token` | GET | OSS upload credentials |
| `/developer/v1/game` | GET | List developer games |
| `/developer/v1/game/create-with-version` | POST | Create game + version |
| `/developer/v1/game/:uri` | PUT | Update game info |
| `/developer/v1/game/update-with-version` | POST | Add version to game |
| `/developer/v1/game/apply-review` | POST | Submit for review |

All authenticated endpoints require `Authorization: Bearer <token>`. The backend returns `code: 0` for success and non-zero for application errors.

## Adding a New Provider

1. Create a new package under `cli/internal/provider/`:

```go
package mypaas

type MyPaaSProvider struct {
    // config fields
}

func (p *MyPaaSProvider) Name() string { return "mypaas" }
func (p *MyPaaSProvider) Deploy(ctx context.Context, config provider.DeployConfig) (*provider.DeployResult, error) {
    // Call your PaaS deploy API
}
func (p *MyPaaSProvider) Status(ctx context.Context, config provider.DeployConfig) (*provider.StatusResult, error) {
    // Call your PaaS status API
}
func (p *MyPaaSProvider) HealthCheck(ctx context.Context, url string) (*provider.HealthResult, error) {
    // HTTP health check against deployed URL
}
```

2. Register in `cli/internal/commands/commands.go` `resolveProvider()`

3. Add tests with a mock PaaS client

## Frontend Packaging

Before deploying the frontend surface, the Phaser H5 build output must be packaged:

```bash
cd frontend && npm run build
```

The resulting `dist/` directory contains `index.html` and bundled static assets. Use `gd-package-frontend` to verify the build is ready for deployment, then pass the directory via `--frontend-dir` to the deploy command.

For details on the frontend build and packaging workflow, see the `gd-package-frontend` and `gd-create-h5-game` skills.

## Structured Output

All providers must return results through the `reporting.Result` type for consistent agent-readable output:

```json
{
  "success": true,
  "message": "Deployed to https://...",
  "code": "SUCCESS",
  "details": { ... }
}
```
