---
name: setup-game-designer-cli
description: Build and verify the Go deploy CLI for first use
trigger: user asks to set up deploy CLI, build game-designer CLI, or prepare CLI for deployment
---

# setup-game-designer-cli

Build and verify the Go deploy CLI after plugin installation.

## Prerequisites

- Go 1.24+ installed and on PATH
- The plugin installed with `cli/` accessible

## When to Apply

- First use after plugin installation
- The user asks to set up or build the deploy CLI
- Before running `prepare-deploy` or `deploy-game-server`
- When the CLI binary is missing or may be stale

## What This Skill Does

1. Check for Go: `go version` (requires 1.24+)
2. Locate the CLI source at `${CLAUDE_PLUGIN_ROOT}/cli/` (or `cli/` relative to the plugin root)
3. Build the CLI binary:
   ```bash
   cd cli && GOWORK=off go build -o game-designer ./cmd/game-designer
   ```
4. Verify the binary works:
   ```bash
   ./cli/game-designer version
   ```
5. Report the result

## Read Scope

- `cli/` — Go CLI source code

## Write Scope

- `cli/game-designer` — the built binary (should be gitignored)
- Does not modify source files in `cli/`, `server-template/`, `sdk-js/`, or any skill files

## Checks

1. `go version` reports 1.24 or later
2. `go build` succeeds in `cli/`
3. The resulting binary exits 0 when run with `version`

## Success Output

```
Deploy CLI built and verified.
- Source: cli/
- Binary: cli/game-designer
- Build: OK
- Version: reports successfully
```

## Failure Output

- Go not found: Install Go 1.24+ and ensure it is on PATH
- Build failure: Report the Go compiler error. Check that `cli/go.mod` and `cli/cmd/game-designer/main.go` exist
- Binary not executable: On Unix, run `chmod +x cli/game-designer`. On Windows, the `.exe` extension is required
- Stale binary: Rebuild by running this skill again. The binary is built from source each time
