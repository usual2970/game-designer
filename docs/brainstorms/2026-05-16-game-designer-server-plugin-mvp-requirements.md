---
date: 2026-05-16
topic: game-designer-server-plugin-mvp
---

# Game Designer Server Plugin MVP Requirements

## Summary

Build a Game Designer server plugin MVP that lets a code agent autonomously connect an activity-style H5 mini-game to a provided backend and deploy it through a CLI to the team's own PaaS. The MVP is a golden-path loop across server template, JS SDK, agent skills, and deployment, centered on login/player profile, game-state save, score submission, and leaderboard.

---

## Problem Frame

The first target user is a product or operations-oriented H5 mini-game creator who can use a code agent but does not want every activity game to become a custom backend project. Today, these users usually ask engineers to write a temporary backend or modify an existing backend for each campaign.

That workaround is slow and hard to repeat. It couples lightweight game launches to backend engineering availability, and every one-off integration increases the chance that SDK calls, server behavior, deployment steps, and runtime checks drift from each other.

The product bet for the MVP is that the code agent, not the human developer, should be able to perform the backend connection and deployment workflow from plugin guidance. Human-readable documentation remains useful, but it is not the primary success condition.

---

## Actors

- A1. Product/operations H5 mini-game creator: Wants to launch an activity-style mini-game with backend-backed participation, progress, score, and leaderboard behavior.
- A2. Code agent: Uses the Game Designer plugin skills, template, SDK, and CLI to connect the game to backend capabilities and deploy it.
- A3. Backend/platform maintainer: Provides the server template, SDK contract, deploy provider, and PaaS integration that the agent relies on.
- A4. Game player: Plays the H5 mini-game and expects identity, saved progress, score submission, and leaderboard behavior to work.
- A5. Team PaaS: The first real deployment target for the generated backend service.

---

## Key Flows

- F1. Agent connects an H5 game to the backend
  - **Trigger:** The developer asks the code agent to add backend support to an activity-style H5 game.
  - **Actors:** A1, A2, A3
  - **Steps:** The agent identifies the game integration surface, uses plugin guidance to create or attach the server template, connects the JS SDK, wires login/player profile, game-state save, score submission, and leaderboard usage, then runs local verification.
  - **Outcome:** The H5 game can exercise the MVP backend capabilities through the JS SDK in a local or test environment.
  - **Covered by:** R1, R2, R3, R4, R5, R6, R7

- F2. Agent deploys the connected backend to the team PaaS
  - **Trigger:** The developer asks the code agent to deploy the game backend after integration passes local checks.
  - **Actors:** A1, A2, A5
  - **Steps:** The agent runs pre-deploy checks, selects the default PaaS provider, supplies required configuration, deploys through the CLI, waits for health/status feedback, and reports the deployed service result.
  - **Outcome:** The backend service is deployed to the team PaaS with enough verification output for the agent and developer to know whether the release succeeded.
  - **Covered by:** R8, R9, R10, R11, R12

- F3. Player completes the activity-game backend loop
  - **Trigger:** A player opens the deployed H5 mini-game.
  - **Actors:** A4
  - **Steps:** The player gets or resumes a session, the game stores player/profile state, progress is saved, the player finishes a run, score is submitted, and leaderboard position can be retrieved.
  - **Outcome:** The deployed game demonstrates the complete backend-backed activity loop.
  - **Covered by:** R2, R3, R4, R5, R13

---

## Requirements

**Agent-first workflow**
- R1. The MVP must provide plugin skills that guide a code agent through the end-to-end workflow: create or attach server backend, connect JS SDK, verify integration, prepare deploy, and deploy to PaaS.
- R2. The workflow must optimize for code-agent autonomy: the agent should be able to discover required steps, run checks, and interpret success or failure without relying on a human reading long setup documentation.
- R3. The plugin must expose a clear golden path for one activity-style H5 game rather than a broad menu of unrelated backend features.

**Backend capabilities**
- R4. The server template must include login/session and player profile capability sufficient for an activity game to identify a returning player.
- R5. The server template must include game-state save capability sufficient for a game to persist and resume simple player progress.
- R6. The server template must include score submission and leaderboard capability sufficient for an activity game to rank players.
- R7. These capabilities must work together as one flow: identify player, save or resume state, submit score, and read leaderboard.

**JS SDK**
- R8. The JS SDK must provide H5-facing integration primitives for the MVP capabilities: login/session, player profile, game-state save, score submission, and leaderboard read.
- R9. The SDK and server behavior must stay aligned through an explicit contract or generation workflow so agent-made changes do not silently drift.
- R10. The SDK integration must include examples or guidance that the code agent can apply inside a game project without inventing call patterns.

**Deploy CLI and PaaS provider**
- R11. The deploy CLI must support a default provider for the team's own PaaS.
- R12. The deploy CLI must run a deploy lifecycle rather than a blind upload: preflight checks, deploy execution, status or health verification, and a final result that the agent can report.
- R13. The deploy system must keep a provider boundary so future cloud providers can be added later without changing the MVP product shape.

**Verification and feedback**
- R14. The MVP must include a local or test verification path that proves the SDK can call the backend capabilities before deployment.
- R15. The MVP must include a deployed verification path that proves the activity-game loop works after release to PaaS.
- R16. Failure output from verification and deployment must be structured enough for a code agent to decide whether to retry, fix integration, ask for missing configuration, or stop with a clear error.

---

## Acceptance Examples

- AE1. **Covers R1, R2, R4, R5, R6, R7, R8.** Given an activity-style H5 game without backend support, when the developer asks the code agent to connect the Game Designer backend, the agent uses plugin skills and SDK guidance to add login/player profile, state save, score submission, and leaderboard behavior without requiring the developer to manually design backend APIs.
- AE2. **Covers R11, R12, R15, R16.** Given a locally verified backend integration and valid PaaS configuration, when the developer asks the code agent to deploy, the agent runs the deploy CLI, receives status or health feedback, and reports either a successful deployed service or a clear actionable failure.
- AE3. **Covers R3, R7, R15.** Given the deployed game, when a player opens it, plays a round, resumes progress, submits a score, and views the leaderboard, the flow demonstrates one coherent activity-game backend loop rather than disconnected service calls.
- AE4. **Covers R9, R10, R14.** Given a server capability change during development, when the agent updates the integration, the SDK contract or generation workflow keeps client usage aligned and local verification catches mismatches before deploy.

---

## Success Criteria

- A code agent can autonomously connect an H5 activity game to the MVP backend capabilities using the plugin.
- A code agent can deploy the connected backend to the team's own PaaS using the plugin and CLI.
- The MVP proves one complete backend-backed game loop: player identity/profile, game-state save, score submission, and leaderboard read.
- The developer does not need to ask a backend engineer to write temporary service code or manually modify an existing backend for the golden-path case.
- A downstream planner can identify the template, SDK, CLI, skill, provider, and verification workstreams without inventing product scope.

---

## Scope Boundaries

- In scope: server-side MVP only.
- In scope: activity-style H5 mini-games, especially campaign or operations-driven games with lightweight backend needs.
- In scope: one default PaaS provider for the team's own platform.
- In scope: a provider boundary for future expansion, but not multiple production providers.
- Deferred: rewards, redemption codes, entitlement issuance, payment, and monetization capabilities.
- Deferred: public cloud providers such as Vercel, Cloudflare Workers, Tencent Cloud, Alibaba Cloud, or others.
- Deferred: visual H5 editor integration and game-play generation.
- Deferred: template marketplace or multi-game template catalog.
- Out of scope: optimizing for medium/heavy online games that need complex realtime multiplayer, matchmaking, or long-running economy systems.

---

## Key Decisions

- Agent autonomy is the core MVP success condition: The primary proof is not that a human can follow docs, but that a code agent can connect and deploy using plugin-provided workflows.
- Use one complete game loop as the MVP spine: Login/profile, state save, score submission, and leaderboard are all included because together they demonstrate real backend value for an activity game.
- Prefer a golden path over a broad platform surface: The MVP should be easy for a code agent to execute reliably before it becomes a generic game backend platform.
- Use provider abstraction with one real provider: The first deployment target is the team's own PaaS, while the product shape leaves room for later providers.

---

## Dependencies / Assumptions

- The team has an existing PaaS that can act as the default deployment target.
- The PaaS can expose enough deploy, status, health, and failure information for an agent-facing CLI workflow.
- The first target game type is an activity-style H5 mini-game, not a complex online game.
- The plugin system can package multiple skills and make them discoverable to the code agent.
- The server template, JS SDK, and CLI will be versioned together or otherwise checked for compatibility.

---

## Outstanding Questions

### Deferred to Planning

- [Affects R9][Technical] What exact contract mechanism should align server capabilities and JS SDK generation or validation?
- [Affects R11, R12][Technical] What deploy operations and status APIs does the team PaaS already expose to the CLI?
- [Affects R14, R15][Technical] What form should local and deployed verification take so it is fast enough for agent iteration but meaningful enough to catch integration errors?
- [Affects R1][Technical] How many skills should ship in the first plugin version versus being folded into fewer MVP skills?
