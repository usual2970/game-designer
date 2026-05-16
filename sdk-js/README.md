# @game-designer/sdk

TypeScript H5 SDK for the Game Designer Slot Machine Server.

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

// Get slot configuration
const config = await client.getSlotConfig();

// Check balance
const balance = await client.getBalance();

// Spin with virtual credit wager
const result = await client.spin({ wager: 10 });
console.log(result.reels);       // symbol grid
console.log(result.paylineWins); // winning paylines
console.log(result.totalPayout); // total payout in credits
console.log(result.balance);     // updated balance

// Spin history
const history = await client.getSpinHistory({ limit: 20 });

// Slot leaderboard
const leaderboard = await client.getSlotLeaderboard({ limit: 10 });
```

## Error Handling

```typescript
import { ApiError } from "@game-designer/sdk";

try {
  await client.spin({ wager: 99999 });
} catch (error) {
  if (error instanceof ApiError) {
    console.log(error.code);    // "INSUFFICIENT_BALANCE" | "INVALID_PARAMETERS" | ...
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
| `getSlotConfig()` | Get slot machine configuration (reels, paylines, wager limits) |
| `getBalance()` | Get current virtual credit balance |
| `spin(request)` | Perform a server-authoritative spin with a wager |
| `getSpinHistory(options?)` | Read spin history with optional pagination |
| `getSlotLeaderboard(options?)` | Read slot leaderboard ranked by highest balance |

## Build

```bash
npm run build
```

## Test

```bash
npm test
```
