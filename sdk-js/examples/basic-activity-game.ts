/**
 * Example: Basic Activity Game Integration
 *
 * Demonstrates the golden path through the Game Designer SDK:
 * 1. Create/resume session
 * 2. Save game state
 * 3. Submit score
 * 4. Read leaderboard
 */
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({
  baseUrl: "http://localhost:8080",
});

async function main() {
  // 1. Login
  const session = await client.createOrResumeSession({
    playerId: "example-player",
    nickname: "DemoPlayer",
  });
  console.log(`Session: ${session.isNew ? "new" : "resumed"} player ${session.playerId}`);

  // 2. Save game state
  await client.saveGameState({
    data: { level: 3, coins: 150, items: ["sword"] },
    checkpoint: "level-3",
  });
  console.log("Game state saved");

  // 3. Load game state (simulate resume)
  const state = await client.getGameState();
  if (state) {
    console.log(`Resumed at checkpoint: ${state.checkpoint}`);
  }

  // 4. Submit score
  const scoreResult = await client.submitScore({ score: 2500 });
  console.log(`Score ${scoreResult.accepted ? "accepted" : "rejected"}, rank: ${scoreResult.rank}`);

  // 5. Read leaderboard
  const leaderboard = await client.getLeaderboard({ limit: 10 });
  console.log(`Leaderboard (${leaderboard.total} players):`);
  for (const entry of leaderboard.entries) {
    console.log(`  #${entry.rank} ${entry.nickname}: ${entry.score}`);
  }
}

main().catch(console.error);
