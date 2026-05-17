---
name: gd-debug-h5-game
description: Diagnose Phaser H5 game frontend failures including white screens, asset issues, and build problems
trigger: user reports a frontend issue, white screen, Phaser error, canvas problem, or asset loading failure
---

# gd-debug-h5-game

Diagnose and troubleshoot Phaser H5 game frontend failures separately from backend and deployment issues.

## Prerequisites

- Phaser H5 frontend project present
- Node.js 20.19+ or 22.12+ installed and on PATH
- Access to browser DevTools (for runtime diagnostics)

## When to Apply

- The user sees a white screen or blank canvas when opening the game
- Phaser boot or scene initialization fails
- Assets fail to load (404 errors in console)
- Canvas sizing or scaling is wrong
- Audio does not play (mobile audio unlock issue)
- Build output is missing or incomplete
- SDK API calls fail from the frontend (as opposed to backend failures)
- Deployed static assets have path issues

## What This Skill Does

1. Identify the failure category (see below)
2. Run targeted diagnostics
3. Suggest fixes
4. Route backend/server issues to `gd-debug-integration` when the root cause is outside the frontend

## Read Scope

- `frontend/` — full frontend project for diagnosis
- Browser DevTools console output

## Write Scope

- None — this skill diagnoses but does not modify project files. Suggest fixes for the user or agent to apply.

## Failure Categories

### Category 1: White Screen / Blank Canvas

**Symptoms:** Page loads but shows empty or solid-color background, no game content
**Diagnosis:**
- Open browser DevTools console — check for JavaScript errors
- Verify `#game-container` element exists in the DOM
- Check that Phaser `Game` config has valid `parent`, `width`, `height`
- Verify scene classes extend `Scene` and are listed in the config

**Fix:** Resolve any import errors, ensure scene classes are exported correctly, verify the `parent` element selector matches the HTML.

### Category 2: Phaser Boot Failure

**Symptoms:** Console shows "Phaser" errors, WebGL/Canvas context failures
**Diagnosis:**
- Check browser WebGL support: open `chrome://gpu` or equivalent
- Verify `type: AUTO` or `type: CANVAS` in game config
- Check for conflicting canvas elements

**Fix:** Use `CANVAS` renderer as fallback, or test in a browser with WebGL support.

### Category 3: Canvas Sizing Issues

**Symptoms:** Game renders but is too large, too small, or misaligned
**Diagnosis:**
- Check `scale.mode` (e.g., `FIT`) and `scale.autoCenter` (e.g., `CENTER_BOTH`)
- Verify the parent container has non-zero dimensions
- Check CSS for conflicting `width`/`height` on the canvas element

**Fix:** Set explicit dimensions on the game container, use `Scale.FIT` with `Scale.CENTER_BOTH`.

### Category 4: Asset 404 / Missing Assets

**Symptoms:** Console shows 404 errors for images or sounds, game shows missing textures
**Diagnosis:**
- Check `public/assets/` for expected files
- Verify asset paths in preload match the actual file names (case-sensitive)
- For production builds, check `vite.config.ts` `base` setting — should be `"./"` for relative paths
- Check that the build output `dist/assets/` contains the expected files

**Fix:** Add missing assets to `public/assets/`, fix path typos, set `base: "./"` in Vite config for relative asset paths.

### Category 5: Audio Unlock (Mobile)

**Symptoms:** Sounds do not play on mobile devices
**Diagnosis:**
- Mobile browsers require a user gesture to unlock audio context
- Check if audio is being played before any user interaction

**Fix:** Ensure audio playback starts only after a user tap/click (e.g., on the spin button).

### Category 6: CORS / Base URL / Session Issues

**Symptoms:** SDK calls fail with network errors, CORS blocked, or 401 Unauthorized
**Diagnosis:**
- Check the game server URL configuration (`meta[name="game-server-url"]` or default)
- Verify the server is running and accessible from the browser origin
- Check browser network tab for CORS headers on failed requests
- Verify session token is set after `createOrResumeSession`

**Fix:** Update the server URL, configure CORS on the server, or ensure the session flow completes before making authenticated calls.
**When to route to `gd-debug-integration`:** If the SDK error code is `SESSION_EXPIRED`, `UNAUTHORIZED`, or the server returns 5xx errors, the root cause is likely on the server side.

### Category 7: SDK API Errors

**Symptoms:** `ApiError` thrown from SDK calls, structured error responses
**Diagnosis:**
- Read the `code` field: `INVALID_PARAMETERS`, `INSUFFICIENT_BALANCE`, `INTERNAL_ERROR`
- Check request payload matches expected types
- Verify wager is within configured min/max range

**Fix:** Validate request data before calling SDK methods, handle `INSUFFICIENT_BALANCE` by blocking spin.
**When to route to `gd-debug-integration`:** If the error code is `INTERNAL_ERROR` or the server response is malformed, the issue is in the server-side logic.

### Category 8: Build Path Issues

**Symptoms:** Production build works locally but fails when deployed (missing assets, broken paths)
**Diagnosis:**
- Check `vite.config.ts` — `base: "./"` for relative paths
- Verify `dist/index.html` references assets with correct relative paths
- Check that all assets in `public/` are copied to `dist/` during build

**Fix:** Use relative base path, ensure all static assets are in `public/` directory.

### Category 9: Deployed Static Asset Issues

**Symptoms:** Game loads locally but not after deployment to PaaS
**Diagnosis:**
- Verify the deployed `index.html` is served correctly
- Check that `dist/` contents were packaged as the frontend surface
- Verify asset paths resolve correctly from the deployment URL
- Check PaaS static file serving configuration

**Fix:** Re-package with `gd-package-frontend` and redeploy. Check PaaS documentation for static file requirements.

## Checks

1. Browser console has no errors
2. Canvas element exists and has content
3. All assets load with 200 status
4. SDK session completes successfully
5. Spin action triggers server call and updates UI

## Success Output

```
Frontend issue diagnosed.
- Category: [category name]
- Root cause: [description]
- Suggested fix: [action]
```

## Failure Output

If diagnosis points to a backend issue:
```
Frontend is healthy. Issue appears to be on the server side.
- Route to: gd-debug-integration
- Evidence: [description of what points to the server]
```
