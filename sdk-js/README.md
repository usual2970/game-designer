# @game-designer/sdk

TypeScript H5 SDK for the Game Designer Server.

## Install

```bash
npm install @game-designer/sdk
```

## Usage

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({
  baseUrl: "http://localhost:8080",
});

// Login
const session = await client.createOrResumeSession({
  playerId: "player-123",
  nickname: "Alice",
});

// Save game state
await client.saveGameState({
  data: { level: 5, coins: 200 },
  checkpoint: "level-5",
});

// Resume game state
const state = await client.getGameState();

// Submit score
await client.submitScore({ score: 1500 });

// Read leaderboard
const leaderboard = await client.getLeaderboard({ limit: 10 });
```

## Error Handling

```typescript
import { ApiError } from "@game-designer/sdk";

try {
  await client.submitScore({ score: 1500 });
} catch (error) {
  if (error instanceof ApiError) {
    console.log(error.code);    // "INVALID_PARAMETERS" | "UNAUTHORIZED" | ...
    console.log(error.message); // human-readable message
    console.log(error.details); // additional context
  }
}
```

## API Methods

| Method | Description |
|--------|-------------|
| `createOrResumeSession(request)` | Create or resume a player session |
| `getPlayerProfile()` | Get current player profile |
| `updatePlayerProfile(request)` | Update player profile |
| `saveGameState(request)` | Save game progress |
| `getGameState()` | Load saved game state (returns `null` if none) |
| `submitScore(request)` | Submit a player score |
| `getLeaderboard(options?)` | Read leaderboard with optional pagination |

## Build

```bash
npm run build
```

## Test

```bash
npm test
```
