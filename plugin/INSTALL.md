# Installation Guide

Install the Game Designer plugin for Claude Code or Codex to scaffold H5 game backends, connect the SDK, and deploy to PaaS.

## Quick Start

1. Clone this repository
2. Install the plugin for your agent (see below)
3. Run `setup-game-designer-cli` to build the deploy CLI
4. Follow the golden path: create server, connect SDK, prepare deploy, deploy

## Claude Code

### Local Install

```bash
# From the repository root
claude plugin install .
```

Or add as a marketplace:

```bash
claude plugin marketplace add .
claude plugin install game-designer
```

### Verify Installation

```bash
# List installed skills
claude plugin skills

# Build the deploy CLI (first use)
# In a Claude Code session, ask: "set up the game-designer deploy CLI"

# Validate the plugin package
./scripts/verify-plugin-package.sh
```

## Codex

### Local Install

1. Open Codex Plugins UI
2. Add a local plugin source pointing to this repository root
3. Enable the `game-designer` plugin

### Verify Installation

```bash
# Build the deploy CLI (first use)
cd cli && GOWORK=off go build -o game-designer ./cmd/game-designer
./cli/game-designer version

# Validate the plugin package
./scripts/verify-plugin-package.sh
```

## Important Notes

- **Plugin installation does not compile the CLI.** Run `setup-game-designer-cli` after installing the plugin.
- **Plugin installation does not deploy the server.** Follow the golden path to create and deploy.
- **Install from the repository root**, not from `plugin/`. The root directory contains all bundled assets (`server-template/`, `cli/`, `sdk-js/`, `contracts/`, `examples/`, `scripts/`).

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.24+ | Build server and deploy CLI |
| Node.js | 18+ | Build SDK and run tests |
| npm | 9+ | SDK package management |

## Troubleshooting

| Problem | Cause | Fix |
|---------|-------|-----|
| Skills not visible after install | Installed `plugin/` instead of repo root | Reinstall from the repository root directory |
| Missing manifest error | `.claude-plugin/plugin.json` not found | Ensure you installed from the repo root |
| Stale cached plugin | Claude Code cached an older version | Remove and reinstall the plugin |
| CLI build fails | Go not installed or wrong version | Install Go 1.24+ and verify with `go version` |
| SDK build fails | Node.js not installed or wrong version | Install Node.js 18+ and verify with `node --version` |
| Skills reference missing files | Installed from a subdirectory | Reinstall from the repository root |
