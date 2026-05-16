# Plugin Installation

How to install the Game Designer plugin for Claude Code and Codex, verify the installation, and prepare for first use.

## Overview

The Game Designer repository is an installable code-agent plugin. The repository root is the plugin root — it contains host manifests, shared skills, and all bundled assets needed to scaffold and deploy H5 game backends.

```
game-designer-backend/           <- Install from here (the plugin root)
├── .claude-plugin/plugin.json   <- Claude Code manifest
├── .codex-plugin/plugin.json    <- Codex manifest
├── plugin/skills/               <- Six shared skills
├── server-template/             <- Go server template
├── cli/                         <- Go deploy CLI source
├── sdk-js/                      <- TypeScript SDK
├── contracts/                   <- OpenAPI contract
├── examples/                    <- Example H5 games
└── scripts/                     <- Verification scripts
```

## Claude Code Installation

### Local Development

Install directly from a cloned repository:

```bash
git clone <repo-url> game-designer-backend
cd game-designer-backend
claude plugin install .
```

### Marketplace Install

Add the repository as a marketplace source and install:

```bash
claude plugin marketplace add .
claude plugin install game-designer
```

### GitHub Repository Install

If the repository is hosted on GitHub:

```bash
claude plugin marketplace add <owner>/<repo>
claude plugin install game-designer
```

### Post-Install Verification

After installing, verify the plugin is working:

1. List available skills — the agent should see six Game Designer skills
2. Build the deploy CLI — ask the agent to run `setup-game-designer-cli`
3. Run package validation:
   ```bash
   ./scripts/verify-plugin-package.sh
   ```

## Codex Installation

### Local Import

1. Open the Codex Plugins UI
2. Add a local plugin source pointing to the repository root directory
3. Enable the `game-designer` plugin
4. Invoke skills using `$skill-name` in the thread (for example, `$setup-game-designer-cli`)

### Post-Install Verification

Same as Claude Code — build the CLI and run package validation.

## First-Use Setup

Plugin installation does **not** compile the deploy CLI. After installation, the first operational step is to build and verify the CLI:

```bash
# Option A: Ask the agent to run the setup skill
# "set up the game-designer deploy CLI"

# Option B: Build manually
cd cli && GOWORK=off go build -o game-designer ./cmd/game-designer
./cli/game-designer version
```

After the CLI is built, follow the [Agent Golden Path](agent-golden-path.md) to create a game backend, connect the SDK, and deploy.

## Installation vs Build vs Deploy

| Phase | What happens | When |
|-------|-------------|------|
| **Plugin install** | Agent discovers skills and bundled assets | Once, when adding the plugin |
| **CLI build** | Go binary compiled from `cli/` source | First use, or when source changes |
| **Server create** | Go server template scaffolded in target project | Per game project |
| **Deploy** | Server pushed to PaaS | Per deployment |

## Common Issues

### Skills not visible after install

The most common cause is installing from the `plugin/` subdirectory instead of the repository root. The plugin root must be the repository root so that `server-template/`, `cli/`, `sdk-js/`, and other bundled assets are available.

Fix: reinstall from the repository root.

### Missing manifest error

Claude Code expects `.claude-plugin/plugin.json` at the plugin root. If installed from the wrong directory, this file is missing.

Fix: reinstall from the repository root.

### Stale cached plugin

Claude Code caches installed plugins. After updating the repository, the cached copy may be stale.

Fix: remove and reinstall the plugin.

### CLI build fails

The deploy CLI requires Go 1.24+. Verify with `go version`.

Fix: install or update Go, then rebuild with `setup-game-designer-cli`.

### Codex plugin not recognized

Verify `.codex-plugin/plugin.json` exists at the repository root and contains valid JSON.

Fix: check the manifest with `./scripts/verify-plugin-package.sh`.
