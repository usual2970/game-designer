# Game Designer Phaser Slot Template

A browser-playable H5 slot machine built with Phaser, TypeScript, and Vite. Uses the Game Designer SDK for server-authoritative gameplay.

## Quick Start

```bash
npm install
npm run dev
```

Open `http://localhost:3000` in a browser. Make sure the Game Designer server is running on `http://localhost:8080`.

## Project Structure

```
src/
  main.ts              Entry point, Phaser game config
  game/
    scenes/
      BootScene.ts     Initial boot scene
      PreloadScene.ts  Asset preloading and placeholder generation
      SlotScene.ts     Slot machine gameplay and SDK integration
    services/
      gameDesignerClient.ts  SDK client factory
    ui/                UI components (extend here)
    styles/
  main.css             Global styles
public/
  assets/              Game assets (images, audio)
tests/                 Vitest unit tests
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start dev server with hot reload |
| `npm run build` | TypeScript check + production build |
| `npm run preview` | Preview production build locally |
| `npm test` | Run unit tests |

## SDK Integration

The slot game uses the Game Designer SDK for:

- **Session** — player login via `createOrResumeSession`
- **Balance** — credit balance via `getBalance`
- **Spin** — server-authoritative spin via `spin`
- **Config** — slot machine config via `getSlotConfig`

All gameplay outcomes come from the server. The client never calculates payouts.

## Configuration

- **Server URL**: Set via `<meta name="game-server-url">` in `index.html`, or defaults to `http://localhost:8080`.
- **Player ID**: Set via `?playerId=xxx` query parameter, or auto-generated.

## Production Build

```bash
npm run build
```

Output goes to `dist/` as static files suitable for deployment as the `frontend` surface in the Game Designer deploy CLI.
