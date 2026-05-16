# Game Designer Plugin

Agent-facing plugin skills for the Game Designer Server MVP.

## Skills

| Skill | Trigger | Purpose |
|-------|---------|---------|
| `create-game-server` | Create game backend | Scaffold and verify the Go server template |
| `connect-js-sdk` | Connect SDK to H5 game | Wire the TypeScript SDK into an H5 game project |
| `prepare-deploy` | Prepare deployment | Run preflight checks and local verification |
| `deploy-game-server` | Deploy to PaaS | Execute the full deploy lifecycle via CLI |
| `debug-server-integration` | Debug failures | Triage SDK, server, contract, and deploy issues |

## Golden Path Sequence

1. `create-game-server` — Set up the Go backend
2. `connect-js-sdk` — Connect the SDK to the H5 game
3. `prepare-deploy` — Verify everything locally
4. `deploy-game-server` — Deploy to PaaS
5. `debug-server-integration` — (if needed) Diagnose failures

## Skill Structure

Each skill defines:
- **When to apply** — What triggers the skill
- **What it does** — Step-by-step actions
- **Files changed** — What surfaces the skill may modify
- **Checks** — Verification steps
- **Success output** — What success looks like
- **Failure output** — What failures look like and how to recover
