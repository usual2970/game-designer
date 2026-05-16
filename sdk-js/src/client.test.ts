import { describe, it, expect, vi, beforeEach } from "vitest";
import { GameDesignerClient } from "./client";
import { ApiError } from "./error";

function mockFetch(responses: { status: number; body: unknown }[]) {
  const queue = [...responses];
  return vi.fn(async () => {
    const next = queue.shift();
    if (!next) throw new Error("no more mock responses");
    return new Response(JSON.stringify(next.body), {
      status: next.status,
      headers: { "Content-Type": "application/json" },
    });
  });
}

describe("GameDesignerClient", () => {
  let client: GameDesignerClient;

  beforeEach(() => {
    client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });
    vi.restoreAllMocks();
  });

  it("creates a session and stores token", async () => {
    const sessionResp = {
      sessionId: "s1",
      token: "tok-123",
      playerId: "p1",
      isNew: true,
      expiresAt: "2026-12-31T00:00:00Z",
    };
    globalThis.fetch = mockFetch([{ status: 200, body: sessionResp }]);

    const result = await client.createOrResumeSession({
      playerId: "p1",
      nickname: "Alice",
    });

    expect(result.token).toBe("tok-123");
    expect(client.getToken()).toBe("tok-123");
  });

  it("gets player profile", async () => {
    client.setToken("tok-123");
    const profileResp = {
      playerId: "p1",
      nickname: "Alice",
      avatarUrl: null,
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
    };
    globalThis.fetch = mockFetch([{ status: 200, body: profileResp }]);

    const result = await client.getPlayerProfile();
    expect(result.nickname).toBe("Alice");
  });

  it("saves and loads game state", async () => {
    client.setToken("tok-123");
    const savedAt = "2026-05-16T12:00:00Z";
    const saveResp = { data: { level: 5 }, checkpoint: "level-5", savedAt };
    const loadResp = { data: { level: 5 }, checkpoint: "level-5", savedAt };

    globalThis.fetch = mockFetch([
      { status: 200, body: saveResp },
      { status: 200, body: loadResp },
    ]);

    await client.saveGameState({ data: { level: 5 }, checkpoint: "level-5" });
    const state = await client.getGameState();

    expect(state?.data.level).toBe(5);
    expect(state?.checkpoint).toBe("level-5");
  });

  it("returns null for no game state (204)", async () => {
    client.setToken("tok-123");
    globalThis.fetch = vi.fn(async () => new Response(null, { status: 204 }));

    const state = await client.getGameState();
    expect(state).toBeNull();
  });

  it("submits a score", async () => {
    client.setToken("tok-123");
    const scoreResp = {
      accepted: true,
      rank: 1,
      bestScore: 1500,
      isNewBest: true,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: scoreResp }]);

    const result = await client.submitScore({ score: 1500 });
    expect(result.accepted).toBe(true);
    expect(result.rank).toBe(1);
  });

  it("reads leaderboard with pagination", async () => {
    client.setToken("tok-123");
    const lbResp = {
      entries: [
        { rank: 1, playerId: "p1", nickname: "Alice", score: 2000, achievedAt: "2026-05-16T12:00:00Z" },
      ],
      total: 1,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: lbResp }]);

    const result = await client.getLeaderboard({ limit: 10, offset: 0 });
    expect(result.entries).toHaveLength(1);
    expect(result.total).toBe(1);
  });

  it("throws ApiError on server error", async () => {
    client.setToken("tok-123");
    globalThis.fetch = mockFetch([
      { status: 401, body: { error: "Invalid token", code: "UNAUTHORIZED" } },
      { status: 401, body: { error: "Invalid token", code: "UNAUTHORIZED" } },
    ]);

    await expect(client.getPlayerProfile()).rejects.toThrow(ApiError);

    try {
      await client.getPlayerProfile();
    } catch (e) {
      expect((e as ApiError).code).toBe("UNAUTHORIZED");
    }
  });

  it("throws ApiError with details", async () => {
    client.setToken("tok-123");
    globalThis.fetch = mockFetch([
      {
        status: 400,
        body: {
          error: "Invalid request parameters",
          code: "INVALID_PARAMETERS",
          details: { fields: ["score"] },
        },
      },
      {
        status: 400,
        body: {
          error: "Invalid request parameters",
          code: "INVALID_PARAMETERS",
          details: { fields: ["score"] },
        },
      },
    ]);

    await expect(client.submitScore({ score: -1 })).rejects.toThrow(ApiError);

    try {
      await client.submitScore({ score: -1 });
    } catch (e) {
      const err = e as ApiError;
      expect(err.code).toBe("INVALID_PARAMETERS");
      expect(err.details?.fields).toEqual(["score"]);
    }
  });
});
