---
name: gd-package-frontend
description: Package the Phaser build output as the frontend surface for deployment
trigger: user asks to package the frontend, prepare frontend for deploy, or bundle the H5 build
---

# gd-package-frontend

Package the Phaser build output into the static frontend surface expected by the deploy CLI and `gd-deploy-game` skill.

## Prerequisites

- Phaser H5 frontend with a successful production build (`npm run build` completed)
- Frontend verification passed via `gd-test-h5-game`
- Deploy CLI built and available (run `gd-setup-cli` first)

## When to Apply

- The user asks to package the frontend for deployment
- Before running `gd-deploy-game` with a frontend surface
- After frontend verification passes and the game is ready to ship

## What This Skill Does

1. Verify the frontend build output exists at `frontend/dist/`
2. Check for required files:
   - `dist/index.html` — entry point for the H5 client
   - `dist/assets/` — bundled JavaScript, CSS, and other assets
3. Validate the build configuration:
   - `vite.config.ts` uses `base: "./"` for relative asset paths (required for PaaS deployment)
   - No absolute paths in asset references
4. Report the package details:
   - Total size of `dist/` contents
   - Number of files
   - Presence of `index.html`
5. Confirm the frontend directory is ready for the deploy CLI's `--frontend-dir` flag

## Packaging Workflow

The frontend packaging step fits between frontend testing and deployment:

```
gd-test-h5-game → gd-package-frontend → gd-deploy-game
```

The deploy CLI expects a directory path containing `index.html` and static assets. The `--frontend-dir` flag in `gd-deploy-game` points to this directory within the uploaded package.

## Read Scope

- `frontend/dist/` — production build output
- `frontend/vite.config.ts` — build configuration

## Write Scope

- None — this skill verifies and reports; it does not modify build output. Build is done by `npm run build`.

## Checks

1. `frontend/dist/` directory exists
2. `frontend/dist/index.html` exists and is non-empty
3. `frontend/dist/assets/` directory exists with at least one bundled file
4. Vite config uses `base: "./"` (relative paths)
5. No absolute URLs in `index.html` asset references

## Success Output

```
Frontend package ready for deployment.
- Build output: frontend/dist/
- index.html: present
- Assets: <count> files
- Total size: <size> KB
- Base path: relative ("./")
- Ready for: gd-deploy-game --frontend-dir
```

## Failure Output

- Missing build output: Run `cd frontend && npm run build` first
- Missing `index.html`: The build may have failed. Check Vite config and TypeScript errors
- Empty assets directory: Check that Phaser and other dependencies are bundled correctly
- Absolute base path: Update `vite.config.ts` to use `base: "./"` for PaaS deployment
- Package too large: Consider code-splitting or reducing asset sizes
