# Plugin Installation

How to install the Game Designer plugin for Claude Code and Codex, verify the installation, and prepare for first use.

## Overview

The Game Designer repository is an installable code-agent plugin. The repository root is the plugin root — it contains host manifests, shared skills, and all bundled assets needed to create browser-playable H5 games, scaffold backends, and deploy.

```
game-designer-backend/           <- Install from here (the plugin root)
├── .claude-plugin/plugin.json   <- Claude Code manifest
├── .codex-plugin/plugin.json    <- Codex manifest
├── skills/                      <- Twelve shared skills
├── server-template/             <- Go server template (slot machine with virtual credits)
├── frontend-template-phaser/    <- Phaser + TypeScript + Vite H5 frontend template
├── cli/                         <- Go deploy CLI source
├── sdk-js/                      <- TypeScript SDK
├── contracts/                   <- OpenAPI contract
├── examples/                    <- Example H5 slot machine games
└── scripts/                     <- Verification scripts
```

## Claude Code Installation

### Quick Dev Testing (Recommended)

Load the plugin directly from disk for the current session. No marketplace setup needed:

```bash
git clone <repo-url> game-designer-backend
cd game-designer-backend
claude --plugin-dir .
```

Inside the session, reload after changes with `/reload-plugins`.

### Persistent Marketplace Install

Install the plugin so it persists across sessions:

```bash
cd game-designer-backend

# Add this repo as a local marketplace from the repository root
claude plugin marketplace add ./

# Install the plugin
claude plugin install game-designer@game-designer-marketplace
```

If the marketplace install shows 0 skills, first confirm `.claude-plugin/plugin.json` points to `./skills/` and that the command ran from the repository root. For local development, `claude --plugin-dir .` is the most direct path.

If a local marketplace cache is stale, remove and reinstall the marketplace:

```bash
claude plugin marketplace remove game-designer-marketplace
claude plugin marketplace add ./
claude plugin install game-designer@game-designer-marketplace
```

### GitHub Repository Install

If the repository is hosted on GitHub:

```bash
claude plugin marketplace add <owner>/<repo>
claude plugin install game-designer@game-designer-marketplace
```

### Post-Install Verification

After installing, verify the plugin is working:

1. List available skills — the agent should see twelve Game Designer skills
2. Build the deploy CLI — ask the agent to run `gd-setup-cli`
3. Run package validation:
   ```bash
   ./scripts/verify-plugin-package.sh
   ```

## Codex Installation

### Local Import

1. Open the Codex Plugins UI
2. Add a local plugin source pointing to the repository root directory
3. Enable the `game-designer` plugin
4. Invoke skills using `$skill-name` in the thread (for example, `$gd-setup-cli`)

### Post-Install Verification

Same as Claude Code — build the CLI and run package validation.

## First-Use Setup

Plugin installation does **not** compile the deploy CLI. After installation, the first operational step is to build and verify the CLI:

```bash
# Option A: Ask the agent to run the setup skill
# "set up the game-designer deploy CLI"

# Option B: Build manually
cd cli
GOWORK=off go build -o game-designer ./cmd/game-designer
./game-designer version
```

After the CLI is built, follow the [Agent Golden Path](agent-golden-path.md) to create a slot machine game backend, connect the SDK, and deploy.

## Installation vs Build vs Deploy

| Phase | What happens | When |
|-------|-------------|------|
| **Plugin install** | Agent discovers skills and bundled assets | Once, when adding the plugin |
| **CLI build** | Go binary compiled from `cli/` source | First use, or when source changes |
| **Server create** | Go server template scaffolded in target project | Per game project |
| **Deploy** | Server pushed to PaaS | Per deployment |

## Common Issues

### Skills not visible after install

The most common cause is installing from a subdirectory instead of the repository root. The plugin root must be the repository root so that `skills/`, `server-template/`, `cli/`, `sdk-js/`, and other bundled assets are available.

Fix: reinstall from the repository root.

### Missing manifest error

Claude Code expects `.claude-plugin/plugin.json` at the plugin root. If installed from the wrong directory, this file is missing.

Fix: reinstall from the repository root.

### Stale cached plugin

Claude Code caches installed plugins. After updating the repository, the cached copy may be stale.

Fix: run `/reload-plugins` in session, or remove and reinstall.

### Marketplace install shows 0 skills

Local marketplace caches can become stale after moving skills or changing manifests.

Fix: remove and re-add the marketplace from the repository root, or use `claude --plugin-dir .` for local development.

### CLI build fails

The deploy CLI requires Go 1.24+. Verify with `go version`.

Fix: install or update Go, then rebuild with `gd-setup-cli`.

### Codex plugin not recognized

Verify `.codex-plugin/plugin.json` exists at the repository root and contains valid JSON.

Fix: check the manifest with `./scripts/verify-plugin-package.sh`.
