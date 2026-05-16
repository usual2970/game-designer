import { describe, it, vi, beforeEach } from "vitest";
import { ActivityGame } from "./game";
function mockFetch(responses) {
    const queue = [...responses];
    return vi.fn(async () => {
        const next = queue.shift();
        if (!next)
            throw new Error("no more mock responses");
        if (next.status === 204) {
            return new Response(null, { status: 204 });
        }
        return new Response(JSON.stringify(next.body), {
            status: next.status,
            headers: { "Content-Type": "application/json" },
        });
    });
}
describe("ActivityGame", () => {
    beforeEach(() => {
        vi.restoreAllMocks();
    });
    it("completes the full activity loop", async () => {
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
            // Get game state (no save exists)
            { status: 204, body: null },
            // Save game state
            {
                status: 200,
                body: {
                    data: { level: 2, coins: 50 },
                    checkpoint: "level-2",
                    savedAt: "2026-05-16T12:00:00Z",
                },
            },
            // Submit score
            {
                status: 200,
                body: {
                    accepted: true,
                    rank: 1,
                    bestScore: 50,
                    isNewBest: true,
                },
            },
            // Get leaderboard
            {
                status: 200,
                body: {
                    entries: [
                        {
                            rank: 1,
                            playerId: "p1",
                            nickname: "Play",
                            score: 50,
                            achievedAt: "2026-05-16T12:00:00Z",
                        },
                    ],
                    total: 1,
                },
            },
        ]);
        const game = new ActivityGame("http://localhost:8080", "p1");
        await game.start();
        // If we get here without throwing, the loop completed
    });
    it("resumes from saved state", async () => {
        globalThis.fetch = mockFetch([
            // Session (returning player)
            {
                status: 200,
                body: {
                    sessionId: "s2",
                    token: "tok-2",
                    playerId: "p1",
                    isNew: false,
                    expiresAt: "2026-12-31T00:00:00Z",
                },
            },
            // Get game state (exists)
            {
                status: 200,
                body: {
                    data: { level: 5, coins: 300 },
                    checkpoint: "level-5",
                    savedAt: "2026-05-16T11:00:00Z",
                },
            },
            // Save game state (updated)
            {
                status: 200,
                body: {
                    data: { level: 6, coins: 370 },
                    checkpoint: "level-6",
                    savedAt: "2026-05-16T12:00:00Z",
                },
            },
            // Submit score
            {
                status: 200,
                body: {
                    accepted: true,
                    rank: 1,
                    bestScore: 370,
                    isNewBest: true,
                },
            },
            // Get leaderboard
            {
                status: 200,
                body: {
                    entries: [
                        {
                            rank: 1,
                            playerId: "p1",
                            nickname: "Play",
                            score: 370,
                            achievedAt: "2026-05-16T12:00:00Z",
                        },
                    ],
                    total: 1,
                },
            },
        ]);
        const game = new ActivityGame("http://localhost:8080", "p1");
        await game.start();
    });
});
