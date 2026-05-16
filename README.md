# Game Designer Server Plugin MVP

A contract-first Go backend template, TypeScript H5 SDK, Go deploy CLI, and agent-facing plugin skills that let a code agent connect and deploy an activity-style H5 game backend.

## Quick Start

### 1. Start the server

```bash
cd server-template
GOWORK=off go run ./cmd/server
# Server starts on :8080
```

### 2. Use the SDK in an H5 game

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });

// Login
const session = await client.createOrResumeSession({
  playerId: "player-123",
  nickname: "Alice",
});

// Save progress
await client.saveGameState({ data: { level: 5 }, checkpoint: "level-5" });

// Submit score
await client.submitScore({ score: 1500 });

// Read leaderboard
const leaderboard = await client.getLeaderboard({ limit: 10 });
```

### 3. Verify locally

```bash
./scripts/verify-local.sh
```

### 4. Deploy

```bash
cd cli
GOWORK=off go run ./cmd/game-designer deploy \
  --server-path ../server-template \
  --app-name my-game \
  --provider fake
```

### 5. Verify deployed

```bash
./scripts/verify-deployed.sh https://my-game.fake.local
```

## Project Structure

```
contracts/          OpenAPI contract (single source of truth)
server-template/    Go backend template
sdk-js/             TypeScript H5 SDK
cli/                Go deploy CLI
plugin/skills/      Agent-facing plugin skills
examples/           Example H5 activity game
scripts/            Verification scripts
docs/               Documentation
```

## Test

```bash
# Go server
cd server-template && GOWORK=off go test ./... -v

# TypeScript SDK
cd sdk-js && npm test

# Go CLI
cd cli && GOWORK=off go test ./... -v

# Example game
cd examples/h5-activity-game && npm test

# All checks
./scripts/verify-local.sh
```

## Golden Path

1. **Create** — Use `create-game-server` skill to scaffold the Go backend
2. **Connect** — Use `connect-js-sdk` skill to wire the SDK into the H5 game
3. **Verify** — Run `./scripts/verify-local.sh`
4. **Deploy** — Use `deploy-game-server` skill via the CLI
5. **Verify deployed** — Run `./scripts/verify-deployed.sh <url>`

## Capabilities

| Capability | Description |
|-----------|-------------|
| Session | Create or resume player sessions |
| Profile | Player profile management |
| Game State | Save and resume game progress |
| Score | Submit player scores |
| Leaderboard | Ranked player scores |

## Documentation

- [Contract-first workflow](docs/integration/contract-first-workflow.md)
- [Local verification](docs/integration/local-verification.md)
- [Agent golden path](docs/integration/agent-golden-path.md)
- [SDK usage](docs/integration/sdk-usage.md)
- [PaaS provider](docs/deployment/paas-provider.md)
- [Deployed verification](docs/deployment/deployed-verification.md)
- [Troubleshooting](docs/deployment/troubleshooting.md)

## License

MIT
