import { GameDesignerClient, ApiError } from "@game-designer/sdk";

/**
 * SlotMachineGame demonstrates the full golden path through the Game Designer SDK.
 *
 * This is the integration point where a real H5 slot game would:
 * 1. Identify the player
 * 2. Load slot config and balance
 * 3. Perform spins
 * 4. Show spin history and leaderboard
 */
export class SlotMachineGame {
  private client: GameDesignerClient;
  private playerId: string;

  constructor(baseUrl: string, playerId: string) {
    this.client = new GameDesignerClient({ baseUrl });
    this.playerId = playerId;
  }

  async start(wager: number): Promise<void> {
    // Step 1: Login
    const session = await this.client.createOrResumeSession({
      playerId: this.playerId,
      nickname: `Player-${this.playerId.slice(0, 4)}`,
    });
    console.log(`Session: ${session.isNew ? "new" : "resumed"}`);

    // Step 2: Get slot configuration
    const config = await this.client.getSlotConfig();
    console.log(`Slot: ${config.reels}x${config.rows}, ${config.paylines.length} paylines, wager ${config.minWager}-${config.maxWager}`);

    // Step 3: Check balance
    const bal = await this.client.getBalance();
    console.log(`Balance: ${bal.balance} credits`);

    if (bal.balance < wager) {
      console.log("Insufficient credits for this wager");
      return;
    }

    // Step 4: Spin
    const result = await this.client.spin({ wager });
    console.log(`Wagered ${result.wager} credits`);
    console.log(`Reels:`);
    for (const row of result.reels) {
      console.log(`  ${row.join(" | ")}`);
    }
    if (result.paylineWins.length > 0) {
      for (const win of result.paylineWins) {
        console.log(`Payline ${win.paylineId}: ${win.count}x ${win.symbol} = ${win.payout}`);
      }
      console.log(`Total payout: ${result.totalPayout} credits`);
    } else {
      console.log("No winning paylines");
    }
    console.log(`Balance: ${result.balance} credits`);

    // Step 5: Spin history
    const history = await this.client.getSpinHistory({ limit: 5 });
    console.log(`\nSpin history (${history.total} total):`);
    for (const entry of history.entries) {
      console.log(`  ${entry.spinId}: wager=${entry.wager} payout=${entry.totalPayout} balance=${entry.balance}`);
    }

    // Step 6: Slot leaderboard
    const leaderboard = await this.client.getSlotLeaderboard({ limit: 5 });
    console.log(`\nSlot Leaderboard (top ${leaderboard.entries.length} of ${leaderboard.total}):`);
    for (const entry of leaderboard.entries) {
      console.log(`  #${entry.rank} ${entry.nickname}: ${entry.balance} credits`);
    }
  }
}
