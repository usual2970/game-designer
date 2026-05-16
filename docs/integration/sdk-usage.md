# SDK Usage Guide

How to use the `@game-designer/sdk` TypeScript SDK in an H5 slot machine game.

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

## Slot Config

```typescript
// Get slot machine configuration
const config = await client.getSlotConfig();
console.log(`Reels: ${config.reels.length}`);
console.log(`Paylines: ${config.paylines.length}`);
console.log(`Min wager: ${config.minWager}, Max wager: ${config.maxWager}`);

// Use config to render reels, symbols, and payline indicators
for (const reel of config.reels) {
  console.log(`Reel symbols: ${reel.symbols.join(", ")}`);
}
```

## Balance

```typescript
// Check virtual credit balance
const balance = await client.getBalance();
console.log(`Current balance: ${balance.credits} credits`);
```

## Spin

```typescript
// Place a spin with a wager amount
const spinResult = await client.spin({
  wager: 100,
});

console.log(`Reel stops: ${spinResult.reelStops.join(", ")}`);
console.log(`Winning paylines: ${spinResult.winningPaylines.length}`);
console.log(`Payout: ${spinResult.payout} credits`);
console.log(`New balance: ${spinResult.balance} credits`);

// Check if the spin was a win
if (spinResult.payout > 0) {
  for (const payline of spinResult.winningPaylines) {
    console.log(`Payline ${payline.paylineId}: ${payline.symbol} x${payline.count} = ${payline.payout}`);
  }
}
```

## Spin History

```typescript
// Get recent spin results
const history = await client.getSpinHistory({ limit: 10, offset: 0 });
for (const entry of history.spins) {
  console.log(
    `Spin #${entry.id}: wager=${entry.wager}, payout=${entry.payout}, balance=${entry.balanceAfter}`
  );
}
```

## Slot Leaderboard

```typescript
// Read leaderboard
const leaderboard = await client.getSlotLeaderboard({ limit: 10, offset: 0 });
for (const entry of leaderboard.entries) {
  console.log(`#${entry.rank} ${entry.nickname}: ${entry.bestPayout} credits (biggest win)`);
}
```

## Error Handling

```typescript
import { ApiError } from "@game-designer/sdk";

try {
  await client.spin({ wager: 100 });
} catch (error) {
  if (error instanceof ApiError) {
    switch (error.code) {
      case "INVALID_PARAMETERS":
        // Fix request and retry (e.g. invalid wager amount)
        break;
      case "INSUFFICIENT_BALANCE":
        // Not enough credits for the wager — show balance to player
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

See `sdk-js/examples/basic-slot-machine.ts` for a complete working example.

See `examples/h5-slot-machine/` for a game integration example.
