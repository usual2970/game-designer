# Installation Guide

Install the Game Designer plugin for Claude Code or Codex to scaffold H5 game backends, connect the SDK, and deploy to PaaS.

## Quick Start

1. Clone this repository
2. Install the plugin for your agent (see below)
3. Run `setup-game-designer-cli` to build the deploy CLI
4. Follow the golden path: create server, connect SDK, prepare deploy, deploy

## Claude Code

### Option A: Quick Dev Testing (Recommended)

Load the plugin directly from disk for the current session. No marketplace setup needed.

```bash
# From the repository root
claude --plugin-dir .
```

Inside the session, reload after changes with `/reload-plugins`.

### Option B: Persistent Marketplace Install

Install the plugin so it persists across sessions.

```bash
# Add this repo as a local marketplace (note: ./ not .)
claude plugin marketplace add ./

# Install the plugin from the marketplace
claude plugin install game-designer@game-designer-marketplace

# Reload plugins in session
# /reload-plugins
```

If the marketplace install shows 0 skills (known bug), use a symlink workaround:

```bash
mkdir -p ~/.claude/plugins/marketplaces
ln -sfn "$(pwd)" ~/.claude/plugins/marketplaces/game-designer-marketplace
claude plugin marketplace add ./
claude plugin install game-designer@game-designer-marketplace
```

### Verify Installation

```bash
# In a Claude Code session, ask: "set up the game-designer deploy CLI"
# Or check skills are visible by invoking any skill

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

## Updating

When the plugin content (skills, manifests, docs) or the CLI source changes:

**1. Update the plugin**

If using `--plugin-dir`, just restart the session or run `/reload-plugins`.

If using marketplace install:

```bash
# Pull latest changes
git pull

# Reinstall to refresh the Claude Code cache
claude plugin remove game-designer@game-designer-marketplace
claude plugin install game-designer@game-designer-marketplace

# Then in session: /reload-plugins
```

**2. Rebuild the CLI**

Plugin update refreshes the source but does not recompile the binary:

```bash
# In a Claude Code session, ask: "set up the game-designer deploy CLI"
# Or manually:
cd cli && GOWORK=off go build -o game-designer ./cmd/game-designer
```

**3. Validate**

```bash
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
| Stale cached plugin | Claude Code cached an older version | `/reload-plugins` in session, or remove and reinstall |
| Marketplace install shows 0 skills | Known bug with local marketplace copy | Use the symlink workaround above |
| CLI build fails | Go not installed or wrong version | Install Go 1.24+ and verify with `go version` |
| SDK build fails | Node.js not installed or wrong version | Install Node.js 18+ and verify with `node --version` |
| Skills reference missing files | Installed from a subdirectory | Reinstall from the repository root |
