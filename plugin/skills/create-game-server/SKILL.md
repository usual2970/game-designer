---
name: create-game-server
description: Scaffold or attach the Go server template for a game backend
trigger: user asks to create a game backend, add server support, or set up game-designer backend
---

# create-game-server

Create or attach the Go server template for an activity-style H5 game backend.

## When to Apply

- The user asks to create a game backend or add server-side support to an H5 project
- The user references "game-designer" server or backend setup
- No existing `server-template/` directory is present in the project

## What This Skill Does

1. Copy or reference the Go server template from `server-template/`
2. Verify the template builds: `cd server-template && GOWORK=off go build ./...`
3. Run the server locally: `cd server-template && GOWORK=off go run ./cmd/server`
4. Verify the server responds on `:8080`
5. Report the result

## Files Changed

- May create: `server-template/` (if not present)
- Reads: `contracts/game-server.openapi.yaml`

## Checks

1. `go build` succeeds in the server directory
2. Server starts on port 8080
3. `POST /api/v1/session` returns 200 with a valid session response

## Success Output

```
Game server created and verified.
- Server path: server-template/
- Build: OK
- Local endpoint: http://localhost:8080
- Session endpoint: POST /api/v1/session
```

## Failure Output

- Build failure: Report the Go compiler error and suggest checking the Go version (1.24+)
- Port conflict: Suggest using a different port or stopping the existing process
- Missing go.mod: Ensure the template was copied correctly
