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

## Current Provider

The MVP ships with a **fake provider** for testing:

```bash
game-designer deploy --provider fake
```

The fake provider simulates deploy, status, and health checks without touching a real PaaS.

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

## Configuration

Provider-specific configuration can be passed via:

- CLI flags (`--provider mypaas`)
- Environment variables (e.g., `PAAS_API_KEY`)
- Config file (future)

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
