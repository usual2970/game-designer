# Contract-First Integration Workflow

This document explains how the Game Designer Server uses a contract-first approach to keep the Go server, TypeScript SDK, and integration tests aligned.

## Overview

The contract-first workflow ensures that:

1. **Single Source of Truth**: The OpenAPI contract defines all MVP capabilities
2. **Type Safety**: TypeScript SDK types are generated from the contract
3. **Alignment**: Server and client stay synchronized through validation
4. **Agent-Friendly**: Code agents can understand the API surface without reading prose

## Workflow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     OpenAPI Contract                             │
│              (contracts/game-server.openapi.yaml)                │
└───────────────────────────┬─────────────────────────────────────┘
                            │
            ┌───────────────┴───────────────┐
            │                               │
            ▼                               ▼
┌───────────────────────┐   ┌───────────────────────────┐
│   Go Server Template  │   │   TypeScript SDK           │
│   (server-template/)  │   │   (sdk-js/)                │
│                       │   │                           │
│ • Implement handlers  │   │ • Generate client types   │
│ • Validate requests   │   │ • Add ergonomic wrappers  │
│ • Return responses    │   │ • Include examples        │
└───────────┬───────────┘   └─────────────┬─────────────┘
            │                               │
            └───────────────┬───────────────┘
                            │
                            ▼
                ┌───────────────────────┐
                │   Integration Tests   │
                │   (scripts/verify-*)  │
                │                       │
                │ • Contract validation │
                │ • Local verification  │
                │ • Deployed verification│
                └───────────────────────┘
```

## Step-by-Step Process

### 1. Define the Contract

**File**: `contracts/game-server.openapi.yaml`

Define your API operations, request/response schemas, and error handling using OpenAPI 3.0.

```yaml
paths:
  /session:
    post:
      operationId: createOrResumeSession
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SessionRequest'
      responses:
        '200':
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SessionResponse'
```

### 2. Validate the Contract

Before implementation, validate the contract syntax and structure:

```bash
# Using npm script
npm run validate:contract

# Or direct validation
npx @apidevtools/swagger-cli validate contracts/game-server.openapi.yaml
```

### 3. Generate TypeScript SDK

Generate client types and methods from the contract:

```bash
# Generate SDK
./scripts/generate-sdk.sh

# Or using OpenAPI Generator
openapi-generator-cli generate \
  -i contracts/game-server.openapi.yaml \
  -g typescript-axios \
  -o sdk-js/generated/
```

This creates:
- Type definitions for all request/response schemas
- Client methods for each operation
- API documentation

### 4. Implement Go Server

Implement handlers that conform to the contract:

```go
// server-template/internal/http/handlers.go
func (h *Handler) CreateOrResumeSession(w http.ResponseWriter, r *http.Request) {
    var req SessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // Return structured error matching contract
        writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", err)
        return
    }

    // Business logic...
    resp := h.sessionService.CreateOrResume(req)

    // Return response matching contract schema
    json.NewEncoder(w).Encode(resp)
}
```

### 5. Run Integration Tests

Verify that server and SDK stay aligned:

```bash
# Local verification
./scripts/verify-local.sh

# Deployed verification
./scripts/verify-deployed.sh
```

These tests verify:
- Contract schema validation
- Server responses match contract
- SDK can call server endpoints
- Error handling is consistent

## Updating the Contract

When you need to change the API:

1. **Update the contract**: Edit `contracts/game-server.openapi.yaml`
2. **Validate**: Run `npm run validate:contract`
3. **Regenerate SDK**: Run `./scripts/generate-sdk.sh`
4. **Update server**: Modify Go handlers to match
5. **Test**: Run verification scripts

## Contract Drift Detection

The integration workflow detects drift between server and contract:

**Local Verification** (`scripts/verify-local.sh`):
- Starts local Go server
- Runs SDK against all endpoints
- Validates responses against contract schemas
- Reports mismatches before deployment

**Deployed Verification** (`scripts/verify-deployed.sh`):
- Calls deployed service endpoints
- Validates responses match contract
- Catches runtime divergences

## Error Handling Patterns

The contract defines structured error responses:

```yaml
components:
  schemas:
    Error:
      type: object
      required: [error]
      properties:
        error:
          type: string
          example: "Invalid request parameters"
        code:
          type: string
          example: "INVALID_PARAMETERS"
        details:
          type: object
          additionalProperties: true
```

Both server and SDK use this pattern:

**Go Server**:
```go
return writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", map[string]interface{}{
    "fields": []string{"score"},
})
```

**TypeScript SDK**:
```typescript
try {
    await sdk.spin({ wager: 0 });
} catch (error: ApiError) {
    console.error(error.code); // "INVALID_PARAMETERS"
    console.error(error.details?.fields); // ["wager"]
}
```

## Benefits for Code Agents

This contract-first approach is designed for agent autonomy:

1. **Discoverable**: All operations are in one OpenAPI file
2. **Validated**: Agents can verify contract syntax before generating code
3. **Typed**: Generated TypeScript types prevent API misuse
4. **Testable**: Verification scripts provide clear success/failure signals
5. **Structured**: Error codes help agents decide retry/fix/stop behavior

## Example: Agent Golden Path

A code agent can follow this workflow:

1. **Read Contract**: Parse OpenAPI to understand available operations
2. **Generate SDK**: Run `./scripts/generate-sdk.sh` for TypeScript client
3. **Integrate**: Add SDK calls to H5 game using examples from `sdk-js/examples/`
4. **Verify Locally**: Run `./scripts/verify-local.sh` before deploy
5. **Deploy**: Use CLI to deploy validated backend
6. **Verify Deployed**: Run `./scripts/verify-deployed.sh` after release

Each step provides clear output that the agent can interpret without human intervention.

## Tools Used

- **Contract Validation**: `@apidevtools/swagger-cli`
- **TypeScript Generation**: `openapi-generator-cli` (typescript-axios)
- **Go Server**: Standard library `net/http`
- **Testing**: Custom verification scripts using curl/sdk-js

## Further Reading

- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Contracts README](../../contracts/README.md)
- [TypeScript SDK Guide](../../sdk-js/README.md)
- [Local Verification](../../scripts/verify-local.sh)
- [Deployed Verification](../../scripts/verify-deployed.sh)
