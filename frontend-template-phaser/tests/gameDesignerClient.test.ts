import { describe, it, expect, vi, beforeEach } from "vitest";
import { GameDesignerClient } from "@game-designer/sdk";

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

describe("GameDesignerClient integration", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("completes session, config, balance, and spin flow", async () => {
    globalThis.fetch = mockFetch([
      {
        status: 200,
        body: {
          sessionId: "s1", token: "tok-1", playerId: "p1",
          isNew: true, expiresAt: "2026-12-31T00:00:00Z",
        },
      },
      {
        status: 200,
        body: {
          reels: 3, rows: 3,
          paylines: [{ id: 1, positions: [0, 0, 0] }],
          symbols: [{ name: "Cherry", payoutMultiplier: 5 }],
          minWager: 1, maxWager: 100, defaultBalance: 1000,
        },
      },
      { status: 200, body: { balance: 1000 } },
      {
        status: 200,
        body: {
          spinId: "spin-1", wager: 10,
          reels: [["Cherry", "Cherry", "Cherry"], ["Lemon", "Bell", "Plum"], ["Orange", "Seven", "BAR"]],
          paylineWins: [{ paylineId: 1, symbol: "Cherry", count: 3, payout: 50 }],
          totalPayout: 50, balance: 1040,
        },
      },
    ]);

    const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });

    const session = await client.createOrResumeSession({ playerId: "p1" });
    expect(session.token).toBe("tok-1");

    const config = await client.getSlotConfig();
    expect(config.reels).toBe(3);

    const balance = await client.getBalance();
    expect(balance.balance).toBe(1000);

    const result = await client.spin({ wager: 10 });
    expect(result.totalPayout).toBe(50);
    expect(result.balance).toBe(1040);
  });

  it("handles insufficient balance error from server", async () => {
    globalThis.fetch = mockFetch([
      {
        status: 200,
        body: {
          sessionId: "s1", token: "tok-1", playerId: "p1",
          isNew: true, expiresAt: "2026-12-31T00:00:00Z",
        },
      },
      { status: 200, body: { balance: 5 } },
      {
        status: 400,
        body: { error: "Insufficient balance", code: "INSUFFICIENT_BALANCE" },
      },
    ]);

    const client = new GameDesignerClient({ baseUrl: "http://localhost:8080" });
    await client.createOrResumeSession({ playerId: "p1" });
    await client.getBalance();

    await expect(client.spin({ wager: 10 })).rejects.toThrow("Insufficient balance");
  });

  it("handles server connection failure gracefully", async () => {
    globalThis.fetch = vi.fn(async () => {
      throw new TypeError("Failed to fetch");
    });

    const client = new GameDesignerClient({ baseUrl: "http://unreachable:9999" });
    await expect(client.getBalance()).rejects.toThrow("Failed to fetch");
  });
});
