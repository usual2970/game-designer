# Game Designer Server API Contract

This directory contains the OpenAPI 3.0 contract that defines the MVP Game Designer Server API surface.

## Contract File

- `game-server.openapi.yaml` — The single source of truth for all MVP API operations, request/response schemas, and error codes.

## Validation

Validate the contract using:

```bash
npx @apidevtools/swagger-cli validate contracts/game-server.openapi.yaml
```

## MVP Capabilities

The contract defines these capability groups:

| Tag | Operations | Purpose |
|-----|-----------|---------|
| Session | `POST /session` | Create or resume a player session |
| Profile | `GET /profile`, `PUT /profile` | Read and update player profile |
| GameState | `GET /game-state`, `PUT /game-state` | Save and resume game progress |
| Score | `POST /scores` | Submit a player score |
| Leaderboard | `GET /leaderboard` | Read ranked player scores |

## Golden Path

The activity-game loop follows this sequence:

1. `POST /session` — Identify or create the player
2. `PUT /game-state` — Save progress during play
3. `GET /game-state` — Resume progress after interruption
4. `POST /scores` — Submit a final score
5. `GET /leaderboard` — Show rankings

## Error Codes

All error responses use the `Error` schema with machine-readable codes:

| Code | Meaning | Agent Action |
|------|---------|-------------|
| `INVALID_PARAMETERS` | Request body or parameters invalid | Fix request and retry |
| `UNAUTHORIZED` | Missing or invalid session token | Re-authenticate via `POST /session` |
| `NOT_FOUND` | Requested resource does not exist | Check resource identifiers |
| `SESSION_EXPIRED` | Session token has expired | Re-authenticate via `POST /session` |
| `INTERNAL_ERROR` | Unexpected server error | Retry once, then stop and report |

## Alignment

- Go server handlers must conform to this contract.
- TypeScript SDK types are generated from this contract.
- Verification scripts validate server responses against these schemas.
