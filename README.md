# Game Designer Server Plugin

A contract-first Go backend template, Phaser H5 frontend template, TypeScript H5 SDK, Go deploy CLI, and agent-facing plugin skills that let a code agent create, connect, and deploy a slot-machine H5 game with virtual credits and server-authoritative spin resolution.

## Install as a Code Agent Plugin

This repository is an installable plugin for Claude Code and Codex. Install from the repository root to get all twelve skills plus bundled assets.

```bash
# Claude Code (quick dev testing)
claude --plugin-dir .

# Validate the package
./scripts/verify-plugin-package.sh
```

See [Plugin Installation](docs/integration/plugin-installation.md) for Claude Code and Codex guides, prerequisites, and troubleshooting.

## Quick Start

### 1. Start the server

```bash
cd server-template
GOWORK=off go run ./cmd/server
# Server starts on :8080
```

### 2. Create a Phaser H5 frontend

```bash
cd frontend-template-phaser
npm install
npm run dev
# Frontend dev server starts on :3000
```

### 3. Use the SDK in an H5 slot game

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });

// Login
const session = await client.createOrResumeSession({
  playerId: "player-123",
  nickname: "Alice",
});

// Get slot config
const config = await client.getSlotConfig();

// Check balance
const balance = await client.getBalance();

// Spin with virtual credits
const result = await client.spin({ wager: 10 });
// result.reels, result.paylineWins, result.totalPayout, result.balance

// Spin history
const history = await client.getSpinHistory({ limit: 20 });

// Slot leaderboard
const leaderboard = await client.getSlotLeaderboard({ limit: 10 });
```

### 4. Verify locally

```bash
./scripts/verify-local.sh
```

### 5. Deploy

```bash
cd cli
GOWORK=off go run ./cmd/game-designer deploy \
  --server-path ../server-template \
  --app-name my-game \
  --provider fake
```

### 6. Verify deployed

```bash
./scripts/verify-deployed.sh https://my-game.fake.local
```

## Project Structure

```
contracts/                OpenAPI contract (single source of truth)
server-template/          Go slot machine backend template
frontend-template-phaser/ Phaser + TypeScript + Vite H5 frontend template
sdk-js/                   TypeScript H5 SDK
cli/                      Go deploy CLI
skills/                   Agent-facing plugin skills
examples/                 Example H5 slot machine game
scripts/                  Verification scripts
docs/                     Documentation
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
cd examples/h5-slot-machine && npm test

# Phaser frontend template
cd frontend-template-phaser && npm test

# All checks
./scripts/verify-local.sh
```

## Golden Path

### Backend + SDK (existing)

0. **Setup CLI** — Use `gd-setup-cli` skill to build the deploy CLI (first use only)
1. **Create** — Use `gd-create-server` skill to scaffold the Go backend into `server/`
2. **Connect** — Use `gd-connect-sdk` skill to wire the SDK into the H5 slot game
3. **Verify** — Run `./scripts/verify-local.sh`
4. **Deploy** — Use `gd-deploy-game` skill to publish the game package
5. **Verify deployed** — Run `./scripts/verify-deployed.sh <url>`

### Frontend + Gameplay (new)

1. **Create frontend** — Use `gd-create-h5-game` skill to scaffold a Phaser H5 frontend
2. **Add gameplay** — Use `gd-create-slot-game` skill to add slot-machine gameplay
3. **Theme** — Use `gd-theme-h5-game` skill to customize the game appearance
4. **Test** — Use `gd-test-h5-game` skill to verify the frontend
5. **Debug** — Use `gd-debug-h5-game` skill to diagnose frontend issues
6. **Package** — Use `gd-package-frontend` skill to prepare for deployment
7. **Deploy** — Use `gd-deploy-game` skill with `--frontend-dir` to publish

## Capabilities

| Capability | Description |
|-----------|-------------|
| Session | Create or resume player sessions |
| Profile | Player profile management |
| Slot Config | Reel configuration, paylines, wager limits |
| Balance | Virtual credit balance |
| Spin | Server-authoritative spin resolution |
| Spin History | Past spin outcomes |
| Leaderboard | Slot leaderboard ranked by highest balance |
| Phaser Frontend | Browser-playable H5 slot machine game |
| Theme | Data-driven game customization |

## Documentation

- [Plugin installation](docs/integration/plugin-installation.md)
- [Agent golden path](docs/integration/agent-golden-path.md)
- [Contract-first workflow](docs/integration/contract-first-workflow.md)
- [Local verification](docs/integration/local-verification.md)
- [SDK usage](docs/integration/sdk-usage.md)
- [PaaS provider](docs/deployment/paas-provider.md)
- [Deployed verification](docs/deployment/deployed-verification.md)
- [Troubleshooting](docs/deployment/troubleshooting.md)

## License

MIT
