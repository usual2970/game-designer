# SDK Usage Guide

How to use the `@game-designer/sdk` TypeScript SDK in an H5 game.

## Installation

```bash
npm install @game-designer/sdk
```

For local development, use the file reference:

```json
{
  "dependencies": {
    "@game-designer/sdk": "file:../sdk-js"
  }
}
```

## Initialization

```typescript
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({
  baseUrl: "http://localhost:8080", // or deployed URL
});
```

## Session Management

```typescript
// Create or resume a session
const session = await client.createOrResumeSession({
  playerId: "user-123",
  nickname: "Alice",
  avatarUrl: "https://example.com/avatar.png",
});

// The SDK automatically stores the session token for subsequent requests
console.log(session.isNew); // true for first login, false for returning
```

## Player Profile

```typescript
// Get profile
const profile = await client.getPlayerProfile();
console.log(profile.nickname);

// Update profile
const updated = await client.updatePlayerProfile({
  nickname: "NewName",
  avatarUrl: "https://example.com/new-avatar.png",
});
```

## Game State

```typescript
// Save progress
await client.saveGameState({
  data: {
    level: 5,
    coins: 200,
    items: ["sword", "shield"],
  },
  checkpoint: "level-5",
});

// Resume progress (returns null if no saved state)
const state = await client.getGameState();
if (state) {
  console.log(`Resuming from ${state.checkpoint}`);
  console.log(`Level: ${state.data.level}`);
}
```

## Score and Leaderboard

```typescript
// Submit a score
const result = await client.submitScore({
  score: 1500,
  metadata: { level: 5, duration: 120 },
});
console.log(`Rank: #${result.rank}, Best: ${result.bestScore}`);
console.log(`New personal best: ${result.isNewBest}`);

// Read leaderboard
const leaderboard = await client.getLeaderboard({ limit: 10, offset: 0 });
for (const entry of leaderboard.entries) {
  console.log(`#${entry.rank} ${entry.nickname}: ${entry.score}`);
}
```

## Error Handling

```typescript
import { ApiError } from "@game-designer/sdk";

try {
  await client.submitScore({ score: 1500 });
} catch (error) {
  if (error instanceof ApiError) {
    switch (error.code) {
      case "INVALID_PARAMETERS":
        // Fix request and retry
        break;
      case "UNAUTHORIZED":
        // Re-authenticate
        break;
      case "NOT_FOUND":
        // Check resource identifiers
        break;
      case "SESSION_EXPIRED":
        // Re-authenticate via createOrResumeSession
        break;
      case "INTERNAL_ERROR":
        // Retry once, then report
        break;
    }
  }
}
```

## Full Example

See `sdk-js/examples/basic-activity-game.ts` for a complete working example.

See `examples/h5-activity-game/` for a game integration example.
