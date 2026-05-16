# Agent Golden Path

Step-by-step instructions for a code agent to connect and deploy a slot machine H5 game backend.

## Overview

The golden path has six steps, each mapped to a plugin skill:

0. **gd-setup-cli** — Build and verify the deploy CLI (first use only)
1. **gd-create-server** — Scaffold the Go server
2. **gd-connect-sdk** — Wire the TypeScript SDK into the H5 game
3. **gd-prepare-deploy** — Run preflight checks
4. **gd-deploy-server** — Deploy via CLI
5. **gd-debug-integration** — (if needed) Triage failures

Before starting, the agent must have the Game Designer plugin installed. See [Plugin Installation](plugin-installation.md) for setup instructions.

## Step 0: Set Up the Deploy CLI

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

## Step 1: Create the Server

```
Skill: gd-create-server
```

The Go server template provides session, profile, slot config, balance, spin, spin history, and slot leaderboard capabilities behind the OpenAPI contract. The server manages virtual credits for slot machine gameplay.

- Copy `server-template/` into the project
- Build: `cd server-template && GOWORK=off go build ./...`
- Start: `cd server-template && GOWORK=off go run ./cmd/server`
- Verify: `POST /api/v1/session` with `{"playerId":"test"}` returns 200

## Step 2: Connect the SDK

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

## Step 3: Prepare for Deploy

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

## Step 4: Deploy

```
Skill: gd-deploy-server
```

Deploy using the CLI:

```bash
cd cli && GOWORK=off go run ./cmd/game-designer deploy \
  --server-path ../server-template \
  --app-name <game-name> \
  --provider fake
```

Expected output:
```json
{"success":true,"message":"Deployed to https://<app>.fake.local","code":"SUCCESS"}
```

## Step 5: Verify Deployed

```bash
./scripts/verify-deployed.sh https://<deployed-url>
```

## If Something Fails

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
