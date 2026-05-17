---
name: gd-create-slot-game
description: Generate playable slot-machine gameplay using the Game Designer SDK golden path
trigger: user asks to create a slot game, add slot gameplay, or build a slot-machine frontend
---

# gd-create-slot-game

Turn a Phaser H5 shell into a playable slot-machine frontend that consumes the Game Designer SDK golden path: session, slot config, balance, spin, history, and leaderboard.

## Prerequisites

- Phaser H5 frontend created via `gd-create-h5-game` (or an existing Phaser project)
- Game Designer SDK connected via `gd-connect-sdk`
- Game Designer server running on a reachable URL
- Game Designer SDK installed as a dependency (`@game-designer/sdk`)

## When to Apply

- The user asks to create a slot-machine game or add slot gameplay to an H5 project
- The user wants a playable slot frontend that integrates with the server-authoritative spin API
- The user references "slot game", "slot machine", or "spin" in the context of an H5 game

## What This Skill Does

1. Verify the frontend project has the Game Designer SDK installed and importable
2. Verify the slot state machine module exists or create it:
   - Define phases: `loading`, `ready`, `spinning`, `result`, `insufficient_balance`, `error`
   - State includes balance, wager, last result, and error message
   - All state transitions are pure functions for testability
3. Implement or verify the slot scene:
   - Boot sequence: connect to server, create session, fetch slot config and balance
   - Wager controls: adjust wager within server-configured min/max range
   - Spin action: validate balance, call SDK `spin()` with current wager, update state with result
   - Result display: show reel symbols, payline wins, payout, and updated balance
   - Error handling: display recoverable errors, route to `insufficient_balance` or `error` states
4. Ensure the client never calculates payouts locally — all outcomes come from the server
5. Verify unit tests cover:
   - Full win spin state transition
   - No-win spin state transition
   - Insufficient balance blocking
   - Wager clamping to min/max
   - Duplicate spin prevention
   - Error state and recovery
6. Run the test suite and report results

## Read Scope

- `${CLAUDE_PLUGIN_ROOT}/frontend-template-phaser/` — template reference for slot scene structure
- `${CLAUDE_PLUGIN_ROOT}/sdk-js/src/` — SDK types and client for API integration
- `${CLAUDE_PLUGIN_ROOT}/examples/h5-slot-machine/src/game.ts` — reference SDK golden-path usage
- Target project's `frontend/src/` — existing frontend code

## Write Scope

- Target project's `frontend/src/game/scenes/SlotScene.ts` — slot gameplay scene
- Target project's `frontend/src/game/services/slotGameState.ts` — slot state machine
- Target project's `frontend/tests/` — test files for slot state and gameplay
- Does not modify SDK, server, or deployment files

## Checks

1. Slot state machine tests pass (state transitions for all phases)
2. No frontend code calculates payouts independently — all results come from `client.spin()`
3. Wager is clamped between server-returned `minWager` and `maxWager`
4. Spin action is blocked when `spinning` phase is active (no duplicate spins)
5. TypeScript compiles without errors

## Success Output

```
Slot gameplay implemented and verified.
- State machine: frontend/src/game/services/slotGameState.ts
- Slot scene: frontend/src/game/scenes/SlotScene.ts
- Tests: frontend/tests/slotGameState.test.ts
- State phases: loading -> ready -> spinning -> result -> ready
- Server-authoritative: all spin outcomes from API
```

## Failure Output

- SDK import failure: Ensure `gd-connect-sdk` ran successfully and `@game-designer/sdk` is installed
- Missing slot config: Verify the server is running and `GET /api/v1/slot/config` returns valid data
- State machine test failure: Report the failing transition and expected vs. actual state
- Payout calculation detected in frontend code: Remove local payout logic and rely on server response
