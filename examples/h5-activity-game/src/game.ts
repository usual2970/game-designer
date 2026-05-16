import {
  GameDesignerClient,
  ApiError,
} from "@game-designer/sdk";

/**
 * ActivityGame demonstrates the full golden path through the Game Designer SDK.
 *
 * This is the integration point where a real H5 game would:
 * 1. Identify the player
 * 2. Save/resume game progress
 * 3. Submit a final score
 * 4. Show the leaderboard
 */
export class ActivityGame {
  private client: GameDesignerClient;
  private playerId: string;

  constructor(baseUrl: string, playerId: string) {
    this.client = new GameDesignerClient({ baseUrl });
    this.playerId = playerId;
  }

  async start(): Promise<void> {
    // Step 1: Login
    const session = await this.client.createOrResumeSession({
      playerId: this.playerId,
      nickname: `Player-${this.playerId.slice(0, 4)}`,
    });
    console.log(`Session: ${session.isNew ? "new" : "resumed"}`);

    // Step 2: Resume progress (if any)
    const savedState = await this.client.getGameState();
    let level = 1;
    let coins = 0;
    if (savedState) {
      level = (savedState.data.level as number) || 1;
      coins = (savedState.data.coins as number) || 0;
      console.log(`Resumed at level ${level} with ${coins} coins`);
    }

    // Step 3: Play (simulate)
    const result = this.playRound(level);
    level = result.level;
    coins += result.coinsEarned;
    console.log(`Round complete: level ${level}, coins ${coins}`);

    // Step 4: Save progress
    await this.client.saveGameState({
      data: { level, coins },
      checkpoint: `level-${level}`,
    });
    console.log("Progress saved");

    // Step 5: Submit score
    const scoreResult = await this.client.submitScore({ score: coins });
    console.log(
      `Score submitted: ${coins} points, rank #${scoreResult.rank}, best: ${scoreResult.bestScore}`
    );

    // Step 6: Show leaderboard
    const leaderboard = await this.client.getLeaderboard({ limit: 5 });
    console.log(`\nLeaderboard (top ${leaderboard.entries.length} of ${leaderboard.total}):`);
    for (const entry of leaderboard.entries) {
      console.log(`  #${entry.rank} ${entry.nickname}: ${entry.score}`);
    }
  }

  private playRound(currentLevel: number): {
    level: number;
    coinsEarned: number;
  } {
    // Simple progression: advance one level, earn random coins
    return {
      level: currentLevel + 1,
      coinsEarned: Math.floor(Math.random() * 100) + 10,
    };
  }
}
