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

  it("gets slot config", async () => {
    client.setToken("tok-123");
    const configResp = {
      reels: 3,
      rows: 3,
      paylines: [{ id: 1, positions: [0, 0, 0] }],
      symbols: [{ name: "Cherry", payoutMultiplier: 5 }],
      minWager: 1,
      maxWager: 100,
      defaultBalance: 1000,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: configResp }]);

    const result = await client.getSlotConfig();
    expect(result.reels).toBe(3);
    expect(result.symbols[0].name).toBe("Cherry");
  });

  it("gets balance", async () => {
    client.setToken("tok-123");
    globalThis.fetch = mockFetch([{ status: 200, body: { balance: 1000 } }]);

    const result = await client.getBalance();
    expect(result.balance).toBe(1000);
  });

  it("sends a spin request and returns typed result", async () => {
    client.setToken("tok-123");
    const spinResp = {
      spinId: "spin-1",
      wager: 10,
      reels: [["Cherry", "Lemon", "Orange"], ["Cherry", "Bell", "Plum"], ["Cherry", "Seven", "BAR"]],
      paylineWins: [{ paylineId: 1, symbol: "Cherry", count: 3, payout: 50 }],
      totalPayout: 50,
      balance: 1040,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: spinResp }]);

    const result = await client.spin({ wager: 10 });
    expect(result.spinId).toBe("spin-1");
    expect(result.totalPayout).toBe(50);
    expect(result.balance).toBe(1040);
    expect(result.paylineWins).toHaveLength(1);
    expect(result.paylineWins[0].symbol).toBe("Cherry");
  });

  it("handles a no-win spin with zero payout", async () => {
    client.setToken("tok-123");
    const spinResp = {
      spinId: "spin-2",
      wager: 10,
      reels: [["Cherry", "Lemon", "Orange"], ["Bell", "Plum", "Seven"], ["BAR", "Cherry", "Lemon"]],
      paylineWins: [],
      totalPayout: 0,
      balance: 990,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: spinResp }]);

    const result = await client.spin({ wager: 10 });
    expect(result.totalPayout).toBe(0);
    expect(result.paylineWins).toHaveLength(0);
    expect(result.balance).toBe(990);
  });

  it("gets spin history with pagination", async () => {
    client.setToken("tok-123");
    const histResp = {
      entries: [
        { spinId: "s1", wager: 10, totalPayout: 50, balance: 1040, reels: [["A"]], paylineWins: [], spunAt: "2026-05-16T12:00:00Z" },
      ],
      total: 5,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: histResp }]);

    const result = await client.getSpinHistory({ limit: 1, offset: 0 });
    expect(result.entries).toHaveLength(1);
    expect(result.total).toBe(5);
  });

  it("reads slot leaderboard with pagination", async () => {
    client.setToken("tok-123");
    const lbResp = {
      entries: [
        { rank: 1, playerId: "p1", nickname: "Alice", balance: 2000, updatedAt: "2026-05-16T12:00:00Z" },
      ],
      total: 1,
    };
    globalThis.fetch = mockFetch([{ status: 200, body: lbResp }]);

    const result = await client.getSlotLeaderboard({ limit: 10, offset: 0 });
    expect(result.entries).toHaveLength(1);
    expect(result.entries[0].balance).toBe(2000);
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

  it("throws ApiError for invalid wager", async () => {
    client.setToken("tok-123");
    globalThis.fetch = mockFetch([
      {
        status: 400,
        body: {
          error: "Invalid request parameters",
          code: "INVALID_PARAMETERS",
          details: { fields: ["wager"] },
        },
      },
    ]);

    try {
      await client.spin({ wager: 0 });
      expect.unreachable("expected ApiError");
    } catch (e) {
      expect(e).toBeInstanceOf(ApiError);
      const err = e as ApiError;
      expect(err.code).toBe("INVALID_PARAMETERS");
      expect(err.details?.fields).toEqual(["wager"]);
    }
  });

  it("throws ApiError for insufficient balance", async () => {
    client.setToken("tok-123");
    globalThis.fetch = mockFetch([
      {
        status: 400,
        body: {
          error: "insufficient virtual credits",
          code: "INSUFFICIENT_BALANCE",
        },
      },
    ]);

    try {
      await client.spin({ wager: 99999 });
      expect.unreachable("expected ApiError");
    } catch (e) {
      expect(e).toBeInstanceOf(ApiError);
      expect((e as ApiError).code).toBe("INSUFFICIENT_BALANCE");
    }
  });
});
