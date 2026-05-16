# Game Designer Server API Contract

This directory contains the OpenAPI 3.0 contract that defines the slot-machine Game Designer Server API surface.

## Contract File

- `game-server.openapi.yaml` — The single source of truth for all slot-machine API operations, request/response schemas, and error codes.

## Validation

Validate the contract using:

```bash
npx @apidevtools/swagger-cli validate contracts/game-server.openapi.yaml
```

Run contract structure tests:

```bash
node contracts/test/validate-contract.mjs
```

## Slot Machine Capabilities

The contract defines these capability groups:

| Tag | Operations | Purpose |
|-----|-----------|---------|
| Session | `POST /session` | Create or resume a player session |
| Profile | `GET /profile`, `PUT /profile` | Read and update player profile |
| Slot | `GET /slot/config`, `POST /spin` | Slot configuration and server-authoritative spin resolution |
| Balance | `GET /balance` | Read virtual credit balance |
| SpinHistory | `GET /spin/history` | Read past spin outcomes |
| Leaderboard | `GET /leaderboard` | Slot leaderboard ranked by highest balance |

## Golden Path

The slot machine loop follows this sequence:

1. `POST /session` — Identify or create the player
2. `GET /slot/config` — Fetch reel configuration, paylines, wager limits
3. `GET /balance` — Check current virtual credit balance
4. `POST /spin` — Wager credits, receive server-resolved reels and payout
5. `GET /spin/history` — Review past spin outcomes
6. `GET /leaderboard` — Show ranked players by highest balance

## Error Codes

All error responses use the `Error` schema with machine-readable codes:

| Code | Meaning | Agent Action |
|------|---------|-------------|
| `INVALID_PARAMETERS` | Request body or parameters invalid (e.g. wager out of range) | Fix request and retry |
| `UNAUTHORIZED` | Missing or invalid session token | Re-authenticate via `POST /session` |
| `NOT_FOUND` | Requested resource does not exist | Check resource identifiers |
| `SESSION_EXPIRED` | Session token has expired | Re-authenticate via `POST /session` |
| `INSUFFICIENT_BALANCE` | Player does not have enough virtual credits for the wager | Lower wager or check balance |
| `INTERNAL_ERROR` | Unexpected server error | Retry once, then stop and report |

## Alignment

- Go server handlers must conform to this contract.
- TypeScript SDK types are aligned with this contract.
- Verification scripts validate server responses against these schemas.
