# Game Designer Server Template

Go backend template for slot-machine H5 mini-games with virtual credits and server-authoritative spin resolution.

## Capabilities

| Capability | Endpoints | Description |
|-----------|----------|-------------|
| Session | `POST /api/v1/session` | Create or resume player sessions |
| Profile | `GET/PUT /api/v1/profile` | Player profile management |
| Slot Config | `GET /api/v1/slot/config` | Reel configuration, paylines, wager limits |
| Balance | `GET /api/v1/balance` | Virtual credit balance |
| Spin | `POST /api/v1/spin` | Server-authoritative spin resolution |
| Spin History | `GET /api/v1/spin/history` | Past spin outcomes |
| Leaderboard | `GET /api/v1/leaderboard` | Slot leaderboard ranked by highest balance |

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
  balance/        Virtual credit balance management
  slot/           Slot config, spin resolution, payline evaluation, leaderboard
  http/           HTTP handlers and routing
  store/          In-memory persistence (swappable)
```

The slot service resolves spins authoritatively — the client submits a wager and the server returns reels, payline wins, payout, and updated balance. Spin outcome generation uses an injectable RNG interface so tests can use deterministic reels.

The store layer is isolated behind a simple interface so PaaS-specific storage can replace the default in-memory implementation.
