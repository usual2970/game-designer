# Agent Golden Path

Step-by-step instructions for a code agent to create, connect, and deploy a slot machine H5 game with frontend and backend.

## Overview

The golden path covers both backend creation and frontend development. Agents can follow the backend-only path if they already have an H5 frontend, or the full path to create a Phaser browser-playable game.

### All Skills

| Skill | Purpose |
|-------|---------|
| `gd-setup-cli` | Build the deploy CLI (first use only) |
| `gd-create-server` | Scaffold the Go backend into `server/` |
| `gd-connect-sdk` | Wire the TypeScript SDK into an H5 project |
| `gd-create-h5-game` | Create a Phaser H5 frontend from the bundled template |
| `gd-create-slot-game` | Add slot-machine gameplay to the frontend |
| `gd-theme-h5-game` | Customize game appearance and assets |
| `gd-test-h5-game` | Verify the frontend through browser validation |
| `gd-debug-h5-game` | Diagnose frontend-specific failures |
| `gd-package-frontend` | Package the frontend build for deployment |
| `gd-deploy-game` | Deploy the game package via CLI |
| `gd-prepare-deploy` | Run preflight checks |
| `gd-debug-integration` | Triage SDK/server/deploy failures |

Before starting, the agent must have the Game Designer plugin installed. See [Plugin Installation](plugin-installation.md) for setup instructions.

## Backend-Only Path

Use this path when you already have an H5 frontend and only need the backend:

0. **gd-setup-cli** — Build the deploy CLI (first use only)
1. **gd-create-server** — Scaffold the Go server into `server/`
2. **gd-connect-sdk** — Wire the TypeScript SDK into the H5 game
3. **gd-prepare-deploy** — Run preflight checks
4. **gd-deploy-game** — Deploy the game package via CLI
5. **gd-debug-integration** — (if needed) Triage failures

## Full Path: Backend + Phaser Frontend

Use this path to create a complete browser-playable slot machine from scratch:

### Phase 1: Backend

0. **gd-setup-cli** — Build the deploy CLI (first use only)
1. **gd-create-server** — Scaffold the Go server into `server/`
2. **gd-connect-sdk** — Wire the TypeScript SDK into the H5 project

### Phase 2: Frontend

3. **gd-create-h5-game** — Create a Phaser H5 frontend project
4. **gd-create-slot-game** — Add slot-machine gameplay with SDK integration
5. **gd-theme-h5-game** — (optional) Customize appearance, colors, assets, sounds
6. **gd-test-h5-game** — Run frontend verification (TypeScript, tests, build, browser)
7. **gd-debug-h5-game** — (if needed) Diagnose frontend failures

### Phase 3: Deploy

8. **gd-package-frontend** — Package the frontend build output
9. **gd-prepare-deploy** — Run preflight checks
10. **gd-deploy-game** — Deploy the game package with all three surfaces
11. **Verify deployed** — Run `./scripts/verify-deployed.sh <url>`

## Step Details

### Step 0: Set Up the Deploy CLI

```
Skill: gd-setup-cli
```

Build the Go deploy CLI from source. This step is required on first use and when the CLI source changes.

```bash
cd cli
GOWORK=off go build -o game-designer ./cmd/game-designer
./game-designer version
```

Expected: binary builds and reports its version without errors.

### Step 1: Create the Server

```
Skill: gd-create-server
```

The Go server template provides session, profile, slot config, balance, spin, spin history, and slot leaderboard capabilities behind the OpenAPI contract. The server manages virtual credits for slot machine gameplay.

- Copy the plugin's `server-template/` into the target project as `server/`
- Build: `cd server && GOWORK=off go build ./...`
- Start: `cd server && GOWORK=off go run ./cmd/server`
- Verify: `POST /api/v1/session` with `{"playerId":"test"}` returns 200

### Step 2: Connect the SDK

```
Skill: gd-connect-sdk
```

The TypeScript SDK wraps all API calls with typed methods.

- Reference `sdk-js/` as a dependency
- Import: `import { GameDesignerClient } from "@game-designer/sdk"`
- Initialize: `new GameDesignerClient({ baseUrl: "http://localhost:8080" })`
- Follow the pattern in `sdk-js/examples/basic-slot-machine.ts`

Key integration points:
- `createOrResumeSession({ playerId })` — call on game start
- `getSlotConfig()` — call to retrieve reel configuration, symbols, and paylines
- `getBalance()` — call to check player's virtual credit balance
- `spin({ wager })` — call when player spins (deducts wager, returns result and payout)
- `getSpinHistory({ limit })` — call to show recent spin results
- `getSlotLeaderboard({ limit })` — call to show top payouts

### Step 3: Create the H5 Frontend

```
Skill: gd-create-h5-game
```

Create a browser-playable Phaser + TypeScript + Vite frontend from the bundled template.

- Copy the plugin's `frontend-template-phaser/` into the target project as `frontend/`
- Install: `cd frontend && npm install`
- Build: `cd frontend && npm run build`
- Test: `cd frontend && npm test`

### Step 4: Add Slot Gameplay

```
Skill: gd-create-slot-game
```

Add slot-machine gameplay that uses the SDK golden path. All spin outcomes come from the server.

### Step 5: Theme the Game (Optional)

```
Skill: gd-theme-h5-game
```

Customize colors, assets, sounds, and mobile layout through data-driven theme tokens.

### Step 6: Test the Frontend

```
Skill: gd-test-h5-game
```

Verify the frontend through TypeScript checks, unit tests, production build, and browser smoke testing.

### Step 7: Debug Frontend Issues (If Needed)

```
Skill: gd-debug-h5-game
```

Diagnose frontend-specific failures: white screen, asset 404, canvas sizing, audio issues, SDK connectivity. Routes backend issues to `gd-debug-integration`.

### Step 8: Package the Frontend

```
Skill: gd-package-frontend
```

Verify the build output is ready for deployment: `index.html` present, assets bundled, relative paths configured.

### Step 9: Prepare for Deploy

```
Skill: gd-prepare-deploy
```

Run pre-deploy checks to verify readiness.

```bash
./scripts/verify-local.sh
```

Expected output:
```json
{"success":true,"message":"All local verification checks passed","code":"SUCCESS"}
```

If any check fails, fix the issue before proceeding.

### Step 10: Deploy the Game

```
Skill: gd-deploy-game
```

Deploy using the CLI. Production deployment publishes three surfaces — backend, socket, and frontend — through `buildConfig`.

Dry run:

```bash
cd cli && GOWORK=off go run ./cmd/game-designer deploy \
  --server-path ../server-template \
  --app-name <game-name> \
  --provider fake
```

Production:

```bash
cd cli && GOWORK=off go run ./cmd/game-designer deploy \
  --provider 3os \
  --mode create \
  --identifier "$GD_IDENTIFIER" \
  --password "$GD_PASSWORD" \
  --game-name "<game-name>" \
  --package-path <path-to-zip> \
  --version <version> \
  --change-log "<description>" \
  --backend-dir "<dir>" \
  --backend-cmd "<cmd>" \
  --frontend-dir "<dir>" \
  --socket-dir "<dir>" \
  --socket-cmd "<cmd>"
```

Expected output:
```json
{"success":true,"message":"Deployed to https://<app>.3os.local","code":"SUCCESS"}
```

### Step 11: Verify Deployed

```bash
./scripts/verify-deployed.sh https://<deployed-url>
```

## If Something Fails

### Backend Issues

Use `gd-debug-integration` to triage:

| Symptom | Category | Fix |
|---------|----------|-----|
| SDK type errors | Contract mismatch | Validate contract, update SDK |
| 400/401 responses | SDK integration | Check session token, request format |
| INSUFFICIENT_BALANCE | Wager/balance | Check balance before spin, validate wager amount |
| Payout mismatch | Balance/payout | Verify payout table matches slot config, check multiplier logic |
| 500 responses | Server error | Check server logs, run tests |
| CLI exits non-zero | Deploy failure | Check preflight, provider config |
| Verification fails | Integration | Run debug skill, check endpoint |
| CLI binary missing | Setup incomplete | Run `gd-setup-cli` first |

### Frontend Issues

Use `gd-debug-h5-game` to triage:

| Symptom | Category | Fix |
|---------|----------|-----|
| White screen / blank canvas | Boot failure | Check Phaser config, scene registration |
| Missing textures / 404s | Asset loading | Check asset paths, public/assets/ directory |
| Game too large/small | Canvas sizing | Check scale mode, parent container |
| No sound on mobile | Audio unlock | Play audio only after user gesture |
| CORS / network errors | SDK connectivity | Check server URL, CORS config |
| API errors in frontend | SDK API | Check request data, handle error codes |
| Build works locally but not deployed | Build paths | Use `base: "./"` in Vite config |
