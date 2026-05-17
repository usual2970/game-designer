---
name: gd-create-h5-game
description: Create or attach a Phaser H5 frontend project from the bundled template
trigger: user asks to create an H5 game frontend, add Phaser support, or set up a browser-playable game
---

# gd-create-h5-game

Create or attach a Phaser + TypeScript + Vite H5 frontend project from the plugin's bundled template.

## Prerequisites

- Node.js 20.19+ or 22.12+ installed and on PATH
- The plugin installed with `frontend-template-phaser/` and `sdk-js/` accessible
- (Optional) Game Designer server running if the frontend needs live SDK calls

## When to Apply

- The user asks to create an H5 game frontend or add Phaser browser-game support
- The user references "game-designer" frontend or Phaser setup
- No existing `frontend/` directory is present in the target project, or the user wants to initialize a new H5 project

## What This Skill Does

1. Check whether a frontend directory already exists in the target project
   - If `frontend/` exists and contains unrelated files (e.g., not a Phaser project), stop and ask the user to confirm the target directory
   - If `frontend/` exists and is already a Phaser H5 project, report that the frontend is already set up
2. Locate the frontend template at `${CLAUDE_PLUGIN_ROOT}/frontend-template-phaser/` (the plugin's bundled asset)
3. Copy the template into the current target project as `frontend/`
4. Wire the SDK dependency:
   - The template already imports `@game-designer/sdk` — no separate SDK wiring step is needed
   - If the consuming project has the SDK available locally, update `frontend/package.json` to point `"@game-designer/sdk"` at the correct relative path
   - Ensure `node_modules/@game-designer/sdk` resolves correctly after `npm install`
5. Adjust `tsconfig.json` paths if needed — the template's `paths` and `include` entries reference `../../sdk-js/src`, which must be updated to match the consuming project's SDK location
6. Install dependencies: `cd frontend && npm install`
7. Verify TypeScript compiles: `cd frontend && npx tsc --noEmit`
8. Verify the production build produces static output: `cd frontend && npm run build`
9. Run unit tests: `cd frontend && npm test`
10. Report the result

## Read Scope

- `${CLAUDE_PLUGIN_ROOT}/frontend-template-phaser/` — Phaser H5 frontend template (plugin-bundled asset)
- `${CLAUDE_PLUGIN_ROOT}/sdk-js/` — TypeScript SDK (referenced as a file dependency by the template)

## Write Scope

- Target project directory — creates `frontend/` by copying from the plugin's `frontend-template-phaser/`
- Does not modify `server-template/`, `sdk-js/`, `cli/`, `contracts/`, or any skill files

## Checks

1. `npm install` succeeds without errors
2. `npx tsc --noEmit` passes with no type errors
3. `npm run build` produces `dist/index.html` and bundled assets
4. `npm test` passes all unit tests

## Success Output

```
Phaser H5 frontend created and verified.
- Frontend path: frontend/
- Template: Phaser + TypeScript + Vite
- Build: OK
- Tests: OK
- Build output: frontend/dist/ (static H5 files)
- Dev server: cd frontend && npm run dev
```

## Failure Output

- `npm install` failure: Check Node.js version (requires 20.19+ or 22.12+). Report the npm error output
- TypeScript errors: Report the compiler errors and suggest checking that the SDK symlink resolves correctly (`node_modules/@game-designer/sdk -> ../../../sdk-js`)
- Build failure: Report the Vite/Rollup error. Common causes include missing Phaser dependency or incompatible Node.js version
- Existing frontend directory: Stop and ask the user to confirm the target. Do not overwrite unrelated files
