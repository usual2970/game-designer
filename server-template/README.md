# Game Designer Server Template

Go backend template for activity-style H5 mini-games.

## Capabilities

| Capability | Endpoints | Description |
|-----------|----------|-------------|
| Session | `POST /api/v1/session` | Create or resume player sessions |
| Profile | `GET/PUT /api/v1/profile` | Player profile management |
| Game State | `GET/PUT /api/v1/game-state` | Save and resume game progress |
| Score | `POST /api/v1/scores` | Submit player scores |
| Leaderboard | `GET /api/v1/leaderboard` | Ranked player scores |

## Run

```bash
cd server-template
go run ./cmd/server
```

Server starts on `:8080`.

## Test

```bash
cd server-template
GOWORK=off go test ./... -v
```

## Architecture

```
cmd/server/       Entry point
internal/
  session/        Session creation and token validation
  profile/        Player profile read/write
  gamestate/      Game state save/load
  leaderboard/    Score submission and ranking
  http/           HTTP handlers and routing
  store/          In-memory persistence (swappable)
```

The store layer is isolated behind a simple interface so PaaS-specific storage can replace the default in-memory implementation.
