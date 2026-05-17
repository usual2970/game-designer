---
name: gd-create-server
description: Scaffold or attach the Go server template for a game backend
trigger: user asks to create a game backend, add server support, or set up game-designer backend
---

# gd-create-server

Create or attach the Go server template for a slot-machine H5 game backend.

## Prerequisites

- Go 1.24+ installed and on PATH
- The plugin installed with `server-template/` and `contracts/` accessible

## When to Apply

- The user asks to create a game backend or add server-side support to an H5 project
- The user references "game-designer" server or backend setup
- No existing `server/` directory is present in the target project

## What This Skill Does

1. Locate the server template at `${CLAUDE_PLUGIN_ROOT}/server-template/` (the plugin's bundled source asset)
2. Copy the template into the current target project as `server/`
3. Verify the copy builds: `cd server && GOWORK=off go build ./...`
4. Run the server locally: `cd server && GOWORK=off go run ./cmd/server`
5. Verify the server responds on `:8080` by sending `POST /api/v1/session`
6. Verify slot endpoints: `GET /api/v1/slot/config` and `GET /api/v1/balance`
7. Report the result

## Read Scope

- `${CLAUDE_PLUGIN_ROOT}/server-template/` — Go server template source (plugin-bundled asset)
- `contracts/game-server.openapi.yaml` — OpenAPI contract for endpoint verification

## Write Scope

- Target project directory — creates `server/` by copying from the plugin's `server-template/`
- Does not modify `contracts/`, `sdk-js/`, `cli/`, or any skill files

## Checks

1. `go build` succeeds in the server directory
2. Server starts on port 8080
3. `POST /api/v1/session` returns 200 with a valid session response
4. `GET /api/v1/slot/config` returns slot configuration with reels, paylines, and wager limits
5. `GET /api/v1/balance` returns virtual credit balance

## Success Output

```
Slot machine game server created and verified.
- Server path: server/
- Build: OK
- Local endpoint: http://localhost:8080
- Session endpoint: POST /api/v1/session
- Slot config: GET /api/v1/slot/config
- Balance: GET /api/v1/balance
- Spin: POST /api/v1/spin
```

## Failure Output

- Build failure: Report the Go compiler error and suggest checking the Go version (1.24+)
- Port conflict: Suggest using a different port or stopping the existing process
- Missing go.mod: Ensure the template was copied correctly from `${CLAUDE_PLUGIN_ROOT}/server-template/` into the target project's `server/` directory
