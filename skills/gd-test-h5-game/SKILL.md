---
name: gd-test-h5-game
description: Verify a Phaser H5 game through static checks, unit tests, build, and browser validation
trigger: user asks to test the H5 game, run frontend verification, or validate the browser build
---

# gd-test-h5-game

Run frontend-specific verification for a Phaser H5 game: static checks, unit tests, production build, local preview, browser smoke testing, mobile viewport checks, and console error checks.

## Prerequisites

- Phaser H5 frontend created via `gd-create-h5-game`
- Node.js 20.19+ or 22.12+ installed and on PATH
- (Optional) Game Designer server running for live SDK integration tests

## When to Apply

- The user asks to test, verify, or validate the H5 frontend
- The user wants to check the game works in a browser
- Before packaging the frontend via `gd-package-frontend`
- After making changes to gameplay, theme, or assets

## What This Skill Does

1. **Static checks**: Run TypeScript compiler in check mode
   ```bash
   cd frontend && npx tsc --noEmit
   ```

2. **Unit tests**: Run the Vitest test suite
   ```bash
   cd frontend && npm test
   ```

3. **Production build**: Verify the Vite build produces static output
   ```bash
   cd frontend && npm run build
   ```
   Confirm `dist/index.html` exists and assets are bundled.

4. **Local preview**: Start the preview server
   ```bash
   cd frontend && npm run preview
   ```

5. **Browser smoke test** (manual or automated):
   - Canvas is non-blank — the Phaser game area renders visible content
   - Controls are clickable/tappable — spin button responds to input
   - Spin animation or result rendering occurs after tapping spin
   - No visible text overlaps at common mobile sizes (375x667, 390x844)
   - No JavaScript console errors (check DevTools console)

6. **Asset load check**: Verify all referenced assets in the theme configuration are present in the build output
   - Symbol texture keys referenced in `theme.symbols` exist
   - Sound paths (if configured) resolve without 404s
   - Background images (if configured) load correctly

7. **Mobile viewport check**:
   - Game fits within viewport without horizontal scroll
   - Touch targets meet the minimum 44px accessibility guideline
   - Content respects safe area insets on notched devices

## Read Scope

- `frontend/` — full frontend project for verification

## Write Scope

- None — this skill is read-only; it verifies but does not modify project files

## Checks

1. TypeScript compiles with no errors
2. All unit tests pass
3. Production build produces `dist/index.html`
4. Browser renders non-blank canvas
5. No console errors on page load
6. All assets load without 404 errors
7. Game fits mobile viewport without overflow

## Success Output

```
H5 frontend verification passed.
- TypeScript: OK
- Unit tests: <N> passed
- Build: OK (dist/index.html + assets)
- Browser smoke: canvas non-blank, controls responsive, no console errors
- Mobile viewport: fits 375x667 and 390x844
- Assets: all loaded
```

## Failure Output

- TypeScript errors: Report the compiler errors with file and line numbers
- Test failures: Report failing test names and assertions
- Build failure: Report the Vite/Rollup error message
- Blank canvas: Check Phaser boot configuration, scene registration, and asset preload
- Console errors: Report each error with the source file and line
- Asset 404s: List missing asset paths and check `public/assets/` directory
- Mobile overflow: Report which elements exceed the viewport width
