/**
 * Example: Basic Slot Machine Integration
 *
 * Demonstrates the golden path through the Game Designer SDK:
 * 1. Create/resume session
 * 2. Get slot config and balance
 * 3. Spin with a virtual credit wager
 * 4. Read spin history and slot leaderboard
 */
import { GameDesignerClient } from "@game-designer/sdk";

const client = new GameDesignerClient({
  baseUrl: "http://localhost:8080",
});

async function main() {
  // 1. Login
  const session = await client.createOrResumeSession({
    playerId: "example-player",
    nickname: "SlotPlayer",
  });
  console.log(`Session: ${session.isNew ? "new" : "resumed"} player ${session.playerId}`);

  // 2. Get slot configuration
  const config = await client.getSlotConfig();
  console.log(`Slot: ${config.reels} reels x ${config.rows} rows, ${config.paylines.length} paylines`);
  console.log(`Wager range: ${config.minWager}-${config.maxWager} credits`);

  // 3. Check balance
  const bal = await client.getBalance();
  console.log(`Balance: ${bal.balance} credits`);

  // 4. Spin
  const wager = 10;
  const result = await client.spin({ wager });
  console.log(`Wagered ${result.wager} credits`);
  console.log(`Reels:`);
  for (const row of result.reels) {
    console.log(`  ${row.join(" | ")}`);
  }
  if (result.paylineWins.length > 0) {
    for (const win of result.paylineWins) {
      console.log(`Payline ${win.paylineId}: ${win.count}x ${win.symbol} = ${win.payout} credits`);
    }
    console.log(`Total payout: ${result.totalPayout} credits`);
  } else {
    console.log("No winning paylines");
  }
  console.log(`Balance: ${result.balance} credits`);

  // 5. Spin history
  const history = await client.getSpinHistory({ limit: 5 });
  console.log(`\nSpin history (${history.total} total):`);
  for (const entry of history.entries) {
    console.log(`  ${entry.spinId}: wager=${entry.wager} payout=${entry.totalPayout} balance=${entry.balance}`);
  }

  // 6. Slot leaderboard
  const leaderboard = await client.getSlotLeaderboard({ limit: 10 });
  console.log(`\nSlot Leaderboard (${leaderboard.total} players):`);
  for (const entry of leaderboard.entries) {
    console.log(`  #${entry.rank} ${entry.nickname}: ${entry.balance} credits`);
  }
}

main().catch(console.error);
