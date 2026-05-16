import { describe, it, expect, vi, beforeEach } from "vitest";
import { SlotMachineGame } from "./game";

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

describe("SlotMachineGame", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("completes the full slot loop", async () => {
    globalThis.fetch = mockFetch([
      // Session
      {
        status: 200,
        body: {
          sessionId: "s1",
          token: "tok-1",
          playerId: "p1",
          isNew: true,
          expiresAt: "2026-12-31T00:00:00Z",
        },
      },
      // Slot config
      {
        status: 200,
        body: {
          reels: 3,
          rows: 3,
          paylines: [{ id: 1, positions: [0, 0, 0] }],
          symbols: [{ name: "Cherry", payoutMultiplier: 5 }],
          minWager: 1,
          maxWager: 100,
          defaultBalance: 1000,
        },
      },
      // Balance
      {
        status: 200,
        body: { balance: 1000 },
      },
      // Spin (win)
      {
        status: 200,
        body: {
          spinId: "spin-1",
          wager: 10,
          reels: [["Cherry", "Cherry", "Cherry"], ["Lemon", "Bell", "Plum"], ["Orange", "Seven", "BAR"]],
          paylineWins: [{ paylineId: 1, symbol: "Cherry", count: 3, payout: 50 }],
          totalPayout: 50,
          balance: 1040,
        },
      },
      // Spin history
      {
        status: 200,
        body: {
          entries: [
            { spinId: "spin-1", wager: 10, totalPayout: 50, balance: 1040, reels: [["A"]], paylineWins: [], spunAt: "2026-05-16T12:00:00Z" },
          ],
          total: 1,
        },
      },
      // Leaderboard
      {
        status: 200,
        body: {
          entries: [
            { rank: 1, playerId: "p1", nickname: "Play", balance: 1040, updatedAt: "2026-05-16T12:00:00Z" },
          ],
          total: 1,
        },
      },
    ]);

    const game = new SlotMachineGame("http://localhost:8080", "p1");
    await game.start(10);
  });

  it("handles no-win spin", async () => {
    globalThis.fetch = mockFetch([
      // Session
      {
        status: 200,
        body: {
          sessionId: "s1",
          token: "tok-1",
          playerId: "p1",
          isNew: true,
          expiresAt: "2026-12-31T00:00:00Z",
        },
      },
      // Slot config
      {
        status: 200,
        body: {
          reels: 3, rows: 3, paylines: [], symbols: [],
          minWager: 1, maxWager: 100, defaultBalance: 1000,
        },
      },
      // Balance
      {
        status: 200,
        body: { balance: 1000 },
      },
      // Spin (no win)
      {
        status: 200,
        body: {
          spinId: "spin-2",
          wager: 10,
          reels: [["Cherry", "Lemon", "Orange"], ["Bell", "Plum", "Seven"], ["BAR", "Cherry", "Lemon"]],
          paylineWins: [],
          totalPayout: 0,
          balance: 990,
        },
      },
      // Spin history
      {
        status: 200,
        body: { entries: [], total: 0 },
      },
      // Leaderboard
      {
        status: 200,
        body: { entries: [], total: 0 },
      },
    ]);

    const game = new SlotMachineGame("http://localhost:8080", "p1");
    await game.start(10);
  });

  it("skips spin when balance is insufficient", async () => {
    globalThis.fetch = mockFetch([
      // Session
      {
        status: 200,
        body: {
          sessionId: "s1",
          token: "tok-1",
          playerId: "p1",
          isNew: true,
          expiresAt: "2026-12-31T00:00:00Z",
        },
      },
      // Slot config
      {
        status: 200,
        body: {
          reels: 3, rows: 3, paylines: [], symbols: [],
          minWager: 1, maxWager: 100, defaultBalance: 1000,
        },
      },
      // Balance (too low)
      {
        status: 200,
        body: { balance: 5 },
      },
    ]);

    const game = new SlotMachineGame("http://localhost:8080", "p1");
    await game.start(10);
    // No spin call made — only 3 fetch calls (session, config, balance)
  });
});
