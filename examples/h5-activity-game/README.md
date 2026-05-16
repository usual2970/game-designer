# Example H5 Activity Game

Minimal example demonstrating the Game Designer SDK golden path.

## What It Does

1. Logs in (creates or resumes a session)
2. Resumes saved progress (if any)
3. Plays a round (simple level progression + random coins)
4. Saves progress
5. Submits a score
6. Shows the leaderboard

## Run Against Local Server

```bash
# Terminal 1: Start the Go server
cd ../../server-template && GOWORK=off go run ./cmd/server

# Terminal 2: Run the example
cd examples/h5-activity-game
npm install
npx tsx src/game.ts
```

## Test

```bash
cd examples/h5-activity-game
npm install
npm test
```

## Integration Points

The `ActivityGame` class in `src/game.ts` marks where a real H5 game would plug in:
- Replace `playRound()` with actual game logic
- The SDK calls remain the same regardless of game implementation
- Error handling uses `ApiError` for structured failure output
