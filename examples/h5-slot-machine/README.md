# H5 Slot Machine Example

Example H5 slot machine game demonstrating the Game Designer SDK golden path.

## Setup

```bash
cd examples/h5-slot-machine
npm install
```

## Test

```bash
npm test
```

## What It Does

This example walks through the complete slot machine loop:

1. Create or resume a player session
2. Fetch slot configuration (reels, paylines, wager limits)
3. Check virtual credit balance
4. Perform a spin with a wager
5. Display reels, payline wins, and payout
6. Read spin history
7. Show the slot leaderboard

The example uses only SDK methods — no hand-written fetch calls for slot endpoints.
