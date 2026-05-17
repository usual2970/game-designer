---
name: gd-theme-h5-game
description: Customize the Phaser H5 game theme with colors, assets, copy, sounds, and mobile layout
trigger: user asks to customize the game theme, change colors, update assets, or apply a campaign style
---

# gd-theme-h5-game

Customize a Phaser H5 game into a campaign-ready activity through theme tokens, assets, copy, sounds, and mobile layout settings.

## Prerequisites

- Phaser H5 frontend with slot gameplay via `gd-create-slot-game`
- Theme schema and default theme present in the frontend project

## When to Apply

- The user asks to customize the game appearance, colors, or branding
- The user wants to apply a campaign or activity theme
- The user asks to change game assets, sounds, or mobile layout

## What This Skill Does

1. Locate the theme configuration in `frontend/src/game/theme/`
2. Review the current `defaultTheme.ts` and `themeSchema.ts` for available customization points
3. Apply user-requested changes to theme tokens:
   - **Title/copy**: Update `title` and `subtitle`
   - **Colors**: Update `colors.*` values (background, primary, secondary, accent, text, textMuted)
   - **Symbol assets**: Update `symbols.*` keys to map to new asset texture keys
   - **Sounds**: Set `sounds.spin`, `sounds.win`, `sounds.buttonClick` paths relative to `assets/`
   - **Mobile layout**: Adjust `layout.safeAreaTop`, `layout.safeAreaBottom`, `layout.minTapTarget`
4. Validate the theme using `validateTheme()` — report errors and warnings
5. If the user provides new assets, place them in `frontend/public/assets/` and update symbol mappings
6. For missing optional sounds, ensure the game falls back to silent mode
7. For missing required symbol assets, report the error before build
8. Run `npm run build` to verify the themed game builds successfully

## Read Scope

- `frontend/src/game/theme/` — theme configuration and schema
- `frontend/public/assets/` — current game assets

## Write Scope

- `frontend/src/game/theme/defaultTheme.ts` — theme configuration values
- `frontend/public/assets/` — new asset files (images, audio)
- `frontend/src/game/scenes/PreloadScene.ts` — asset preload list if new assets are added
- Does not modify slot state logic, SDK integration, or server code

## Checks

1. `validateTheme()` reports no errors for the updated theme
2. TypeScript compiles without errors
3. Production build succeeds
4. Required symbol assets (cherry, lemon, orange, plum, bell, seven, bar) are mapped

## Success Output

```
Game theme updated.
- Theme: <name>
- Colors: background=<bg>, primary=<pri>, accent=<acc>
- Symbols: <count> mapped
- Sounds: <count> configured
- Validation: passed
- Build: OK
```

## Failure Output

- Validation errors: Report each invalid field with the expected format
- Missing required symbol: List the missing symbol names and suggest adding assets to `public/assets/`
- Build failure: Report the error and check for broken asset references
- Asset too large: Warn about bundle size and suggest compression or smaller dimensions
