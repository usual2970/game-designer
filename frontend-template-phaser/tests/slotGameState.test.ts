import { describe, it, expect } from "vitest";
import {
  createSlotGameState,
  setReady,
  setWager,
  startSpin,
  completeSpin,
  spinError,
  resetFromResult,
} from "../src/game/services/slotGameState";

describe("slotGameState", () => {
  const baseConfig = {
    reels: 3, rows: 3,
    paylines: [{ id: 1, positions: [0, 0, 0] }],
    symbols: [{ name: "Cherry", payoutMultiplier: 5 }],
    minWager: 1, maxWager: 100, defaultBalance: 1000,
  };

  it("transitions from loading to ready with balance and config", () => {
    const state = createSlotGameState();
    expect(state.phase).toBe("loading");

    const ready = setReady(state, 1000, baseConfig);
    expect(ready.phase).toBe("ready");
    expect(ready.balance).toBe(1000);
    expect(ready.minWager).toBe(1);
    expect(ready.maxWager).toBe(100);
  });

  it("transitions from ready to spinning to result on winning spin", () => {
    let state = setReady(createSlotGameState(), 1000, baseConfig);

    state = startSpin(state);
    expect(state.phase).toBe("spinning");

    const result = {
      spinId: "s1", wager: 10,
      reels: [["Cherry", "Cherry", "Cherry"], ["Lemon", "Bell", "Plum"], ["Orange", "Seven", "BAR"]],
      paylineWins: [{ paylineId: 1, symbol: "Cherry", count: 3, payout: 50 }],
      totalPayout: 50, balance: 1040,
    };

    state = completeSpin(state, result);
    expect(state.phase).toBe("result");
    expect(state.balance).toBe(1040);
    expect(state.lastResult?.totalPayout).toBe(50);
  });

  it("transitions to result with no-win spin", () => {
    let state = setReady(createSlotGameState(), 1000, baseConfig);
    state = startSpin(state);

    const result = {
      spinId: "s2", wager: 10,
      reels: [["Cherry", "Lemon", "Orange"], ["Bell", "Plum", "Seven"], ["BAR", "Cherry", "Lemon"]],
      paylineWins: [], totalPayout: 0, balance: 990,
    };

    state = completeSpin(state, result);
    expect(state.phase).toBe("result");
    expect(state.balance).toBe(990);
    expect(state.lastResult?.totalPayout).toBe(0);
  });

  it("blocks spin when balance is below wager", () => {
    let state = setReady(createSlotGameState(), 5, baseConfig);
    state = setWager(state, 10);

    state = startSpin(state);
    expect(state.phase).toBe("insufficient_balance");
    expect(state.errorMessage).toBe("Insufficient balance for this wager");
  });

  it("clamps wager to min/max range", () => {
    let state = setReady(createSlotGameState(), 1000, baseConfig);

    state = setWager(state, 0);
    expect(state.wager).toBe(1);

    state = setWager(state, 500);
    expect(state.wager).toBe(100);
  });

  it("enters error state on SDK failure and can recover", () => {
    let state = setReady(createSlotGameState(), 1000, baseConfig);
    state = startSpin(state);

    state = spinError(state, "Server unavailable");
    expect(state.phase).toBe("error");
    expect(state.errorMessage).toBe("Server unavailable");

    state = resetFromResult(state);
    expect(state.phase).toBe("ready");
  });

  it("prevents spin when not in ready phase", () => {
    const loading = createSlotGameState();
    expect(startSpin(loading).phase).toBe("loading");

    let state = setReady(createSlotGameState(), 1000, baseConfig);
    state = startSpin(state);
    expect(startSpin(state).phase).toBe("spinning");
    expect(startSpin(state)).toBe(state);
  });

  it("resets to insufficient_balance when balance is too low after spin", () => {
    let state = setReady(createSlotGameState(), 15, baseConfig);
    state = setWager(state, 10);
    state = startSpin(state);

    const result = {
      spinId: "s3", wager: 10,
      reels: [["A", "B", "C"], ["D", "E", "F"], ["G", "H", "I"]],
      paylineWins: [], totalPayout: 0, balance: 5,
    };

    state = completeSpin(state, result);
    state = resetFromResult(state);
    expect(state.phase).toBe("insufficient_balance");
  });
});
