# Game Designer Plugin

Agent-facing plugin skills for the Game Designer Server MVP.

## Installation

This repository is an installable plugin for Claude Code and Codex. See [INSTALL.md](INSTALL.md) for platform-specific installation guides.

## Plugin Manifests

The repository root contains host-specific manifests that point to `plugin/skills/` as the shared skill source:

- `.claude-plugin/plugin.json` — Claude Code plugin manifest
- `.claude-plugin/marketplace.json` — Claude Code marketplace catalog
- `.codex-plugin/plugin.json` — Codex plugin manifest

## Skills

| Skill | Trigger | Purpose |
|-------|---------|---------|
| `setup-game-designer-cli` | Set up deploy CLI | Build and verify the Go deploy CLI |
| `create-game-server` | Create game backend | Scaffold and verify the Go server template |
| `connect-js-sdk` | Connect SDK to H5 game | Wire the TypeScript SDK into an H5 game project |
| `prepare-deploy` | Prepare deployment | Run preflight checks and local verification |
| `deploy-game-server` | Deploy to PaaS | Execute the full deploy lifecycle via CLI |
| `debug-server-integration` | Debug failures | Triage SDK, server, contract, and deploy issues |

## Golden Path Sequence

1. `setup-game-designer-cli` — Build and verify the deploy CLI (first use)
2. `create-game-server` — Set up the Go backend
3. `connect-js-sdk` — Connect the SDK to the H5 game
4. `prepare-deploy` — Verify everything locally
5. `deploy-game-server` — Deploy to PaaS
6. `debug-server-integration` — (if needed) Diagnose failures

## Skill Structure

Each skill defines:
- **When to apply** — What triggers the skill
- **What it does** — Step-by-step actions
- **Files changed** — What surfaces the skill may modify
- **Checks** — Verification steps
- **Success output** — What success looks like
- **Failure output** — What failures look like and how to recover

## Validation

Run `./scripts/verify-plugin-package.sh` to check that manifests, skills, and bundled assets are structurally sound.
