---
title: "refactor: Clarify Phaser SDK skill boundary"
type: refactor
status: completed
date: 2026-05-17
---

# refactor: Clarify Phaser SDK skill boundary

## Summary

Update the H5 frontend skill guidance so the bundled Phaser template is treated as SDK-ready by default. `gd-create-h5-game` should create a working template frontend without routing through `gd-connect-sdk`, while `gd-connect-sdk` remains available for users who already have a custom H5 project and need to add Game Designer SDK integration.

## Requirements

- R1. New projects created from the Phaser template must not require a separate `gd-connect-sdk` step.
- R2. `gd-create-h5-game` must still verify that `@game-designer/sdk` resolves correctly after copying the template.
- R3. `gd-connect-sdk` must remain documented as the path for existing or custom H5 frontends.
- R4. Skill docs and package validation should prevent future guidance from reintroducing the incorrect dependency between the template path and `gd-connect-sdk`.

## Scope Boundaries

- This plan does not remove `gd-connect-sdk`.
- This plan does not copy SDK source into `frontend-template-phaser`.
- This plan does not change server API behavior, SDK client behavior, or Phaser gameplay behavior.
- This plan does not change deployment packaging beyond clarifying that the template build already includes the SDK through normal frontend bundling.

## Context & Research

### Relevant Code and Patterns

- `frontend-template-phaser/package.json` already declares `@game-designer/sdk` as a dependency.
- `frontend-template-phaser/src/game/services/gameDesignerClient.ts` already imports and wraps `GameDesignerClient` from `@game-designer/sdk`.
- `frontend-template-phaser/tsconfig.json` already maps `@game-designer/sdk` to `sdk-js/src/index.ts` for local development.
- `frontend-template-phaser/tests/gameDesignerClient.test.ts` already verifies the SDK golden path through the template dependency.
- `skills/gd-create-h5-game/SKILL.md` currently says to run `gd-connect-sdk` if the SDK is not available locally, which conflicts with the template-as-ready-path model.
- `skills/gd-connect-sdk/SKILL.md` already describes wiring the SDK into an existing H5 game project.
- `scripts/verify-plugin-package.sh` already validates skill docs and bundled assets, and is the right place for a small regression check on skill-boundary wording.

### Institutional Learnings

- No `docs/solutions/` directory or institutional learning files were present during planning.

### External References

- External research was not needed. This is a repo-local documentation and skill-boundary correction using existing template and SDK patterns.

## Key Technical Decisions

- Keep one SDK implementation source in `sdk-js/`: copying SDK code into `frontend-template-phaser` would create drift risk.
- Treat `frontend-template-phaser` as batteries-included for SDK usage: the template owns the import pattern, thin client wrapper, tests, and Phaser scene usage.
- Narrow `gd-connect-sdk` to custom or existing frontend projects: it remains useful, but should not be a prerequisite in the official template creation path.
- Add package-level wording validation: a lightweight script check can catch accidental reintroduction of "run `gd-connect-sdk` first" guidance in the template skill.

## Open Questions

### Resolved During Planning

- Should `gd-connect-sdk` be deleted? Resolved: no, keep it for non-template/custom frontend projects.
- Should SDK source be copied into the Phaser template? Resolved: no, keep the SDK as a dependency and bundle it into the built frontend through Vite.

### Deferred to Implementation

- Exact wording of the revised skill docs: choose concise language during editing while preserving the responsibility boundary described here.

## Implementation Units

### U1. Update Phaser template creation guidance

**Goal:** Make `gd-create-h5-game` describe the Phaser template path as SDK-ready and remove the instruction to route through `gd-connect-sdk`.

**Requirements:** R1, R2

**Dependencies:** None

**Files:**
- Modify: `skills/gd-create-h5-game/SKILL.md`

**Approach:**
- Replace the current SDK wiring step with guidance that the template already imports `@game-designer/sdk`.
- Keep the local dependency/path adjustment requirement, because copied templates may need a different relative `sdk-js` location in consuming projects.
- Keep install, TypeScript, build, and unit-test verification as the proof that the SDK dependency resolves correctly.
- Update failure guidance so SDK-not-found errors point to dependency/path resolution, not to running `gd-connect-sdk`.

**Patterns to follow:**
- Existing `gd-create-h5-game` structure: prerequisites, when to apply, what it does, scopes, checks, success output, failure output.
- Existing template dependency pattern in `frontend-template-phaser/package.json` and `frontend-template-phaser/tsconfig.json`.

**Test scenarios:**
- Happy path: reading the skill shows a new-template flow that copies the Phaser template, adjusts SDK dependency resolution, installs, builds, and tests without invoking `gd-connect-sdk`.
- Error path: SDK resolution failure guidance tells the implementer to fix `@game-designer/sdk` dependency/path resolution rather than route to `gd-connect-sdk`.

**Verification:**
- `skills/gd-create-h5-game/SKILL.md` no longer says to run `gd-connect-sdk` for the bundled Phaser template path.
- The skill still includes explicit SDK dependency resolution checks.

### U2. Clarify custom-project role for SDK connection

**Goal:** Make `gd-connect-sdk` explicitly read as the skill for existing or custom H5 projects, not the default path after creating the bundled Phaser template.

**Requirements:** R3

**Dependencies:** None

**Files:**
- Modify: `skills/gd-connect-sdk/SKILL.md`

**Approach:**
- Tighten "When to Apply" so it names existing/custom H5 projects.
- Add a short boundary note that new frontends created from the bundled Phaser template already include the SDK integration pattern.
- Preserve the current golden-path SDK usage guidance for non-template projects.

**Patterns to follow:**
- Current `gd-connect-sdk` golden-path section.
- Current `frontend-template-phaser` service wrapper pattern for `GameDesignerClient`.

**Test scenarios:**
- Happy path: a user with an existing H5 project can still identify `gd-connect-sdk` as the right integration skill.
- Boundary path: a user creating a new Phaser template frontend is directed to `gd-create-h5-game`, not asked to run both skills.

**Verification:**
- The two skill docs describe complementary responsibilities with no circular or duplicate setup instructions.

### U3. Add package validation for the skill boundary

**Goal:** Add a lightweight regression check so package verification fails if the Phaser template creation skill again instructs agents to run `gd-connect-sdk` as part of the official template path.

**Requirements:** R4

**Dependencies:** U1, U2

**Files:**
- Modify: `scripts/verify-plugin-package.sh`

**Approach:**
- Add a documentation accuracy check that rejects guidance in `gd-create-h5-game` which routes missing SDK availability through `gd-connect-sdk`.
- Add or adjust a positive check that `gd-connect-sdk` still references existing/custom H5 project integration.
- Keep the script style consistent with existing Python-backed text assertions.

**Patterns to follow:**
- Existing documentation accuracy checks in `scripts/verify-plugin-package.sh`.
- Existing skill consistency checks for slot machine concepts.

**Test scenarios:**
- Happy path: package verification passes when `gd-create-h5-game` treats the template as SDK-ready and `gd-connect-sdk` is scoped to existing/custom projects.
- Regression path: package verification fails if `gd-create-h5-game` reintroduces "run `gd-connect-sdk` first" for the template path.

**Verification:**
- Plugin package verification includes the new boundary check and passes after the doc updates.

### U4. Verify template SDK behavior remains intact

**Goal:** Ensure the documentation change does not mask an actual template dependency problem.

**Requirements:** R1, R2

**Dependencies:** U1

**Files:**
- Test: `frontend-template-phaser/tests/gameDesignerClient.test.ts`
- Test: `frontend-template-phaser/tests/slotGameState.test.ts`
- Test: `frontend-template-phaser/tests/themeSchema.test.ts`

**Approach:**
- Use the existing template tests and build checks as the verification surface.
- Treat failures as dependency-resolution or template correctness issues, not as reasons to restore the `gd-connect-sdk` prerequisite.

**Patterns to follow:**
- Existing Vitest coverage in `frontend-template-phaser/tests/`.
- Existing `gd-create-h5-game` checks for install, TypeScript, production build, and unit tests.

**Test scenarios:**
- Happy path: the template can import `@game-designer/sdk` and complete the session/config/balance/spin test flow.
- Error path: server/API failures remain surfaced through the SDK client and do not produce a blank or unhandled template state.
- Integration: the template production build bundles the SDK dependency into static frontend output.

**Verification:**
- Template tests and production build pass.
- Package verification passes.

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| Agents may still treat `gd-connect-sdk` as mandatory because older wording remains elsewhere. | Search skill docs and integration docs for `gd-connect-sdk` references and keep the distinction explicit where relevant. |
| Consumers without a local `sdk-js/` path may hit dependency resolution errors. | Keep `gd-create-h5-game` focused on fixing the `@game-designer/sdk` dependency path for the copied template. |
| Verification script becomes too brittle by matching exact prose. | Check for the incorrect relationship at a phrase level, not the exact final wording of every sentence. |

## System-Wide Impact

- **Skill routing:** New Phaser template creation becomes a single-skill path; custom frontend SDK integration remains a separate skill path.
- **Template behavior:** Runtime behavior is unchanged; the template already imports the SDK.
- **Deployment:** Static frontend output remains the deployable surface, with SDK code bundled by the frontend build.
- **Unchanged invariants:** `sdk-js/` remains the single source of SDK implementation, and `gd-connect-sdk` remains available for non-template projects.
