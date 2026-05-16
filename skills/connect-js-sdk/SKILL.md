---
name: connect-js-sdk
description: Connect the TypeScript H5 SDK to an H5 game project
trigger: user asks to connect SDK, add game-designer SDK, integrate backend SDK into H5 game
---

# connect-js-sdk

Connect the Game Designer TypeScript SDK to an H5 game project.

## Prerequisites

- Node.js 18+ installed and on PATH
- The plugin installed with `sdk-js/` accessible
- A running game server (from `create-game-server`) or a known server URL

## When to Apply

- The user asks to add backend/SDK integration to an H5 game
- The user references connecting the game to the server
- The SDK needs to be wired into an existing H5 project

## What This Skill Does

1. Locate the SDK source at `${CLAUDE_PLUGIN_ROOT}/sdk-js/` (or `sdk-js/` relative to the plugin root)
2. Install the SDK: reference `sdk-js/` or install from the package
3. Import the SDK client: `import { GameDesignerClient } from "@game-designer/sdk"`
4. Add SDK initialization code to the H5 game entry point
5. Wire up the golden path calls using patterns from `sdk-js/examples/basic-activity-game.ts`
6. Run the SDK tests to verify: `cd sdk-js && npm test`

## Read Scope

- `sdk-js/` — TypeScript SDK source and examples
- `sdk-js/examples/basic-activity-game.ts` — golden path integration pattern
- `contracts/game-server.openapi.yaml` — API contract for type reference

## Write Scope

- Target H5 game project only — entry point or integration module
- Does not modify `sdk-js/`, `server-template/`, `cli/`, or any skill files

## SDK Integration Pattern

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });

// Login
const session = await client.createOrResumeSession({
  playerId: currentUser.id,
  nickname: currentUser.name,
});

// Save progress
await client.saveGameState({ data: currentProgress, checkpoint: "level-3" });

// Submit score
const result = await client.submitScore({ score: finalScore });

// Read leaderboard
const leaderboard = await client.getLeaderboard({ limit: 10 });
```

## Checks

1. SDK builds without errors: `cd sdk-js && npm run build`
2. SDK tests pass: `cd sdk-js && npm test`
3. No hand-written HTTP calls — all calls go through the SDK client
4. Error handling uses `ApiError` class with structured codes

## Success Output

```
SDK connected to H5 game.
- SDK import: OK
- Session flow: wired
- Game state: wired
- Score submission: wired
- Leaderboard: wired
- Error handling: using ApiError
```

## Failure Output

- Build errors: Check TypeScript version compatibility (5.4+) and import paths
- Missing fetch: The SDK requires a browser environment with native fetch
- Type mismatch: Ensure SDK types align with game data shapes
- SDK not found: Verify `sdk-js/` exists at the plugin root or install the package
