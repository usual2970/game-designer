import type { SlotConfigResponse, SpinResult } from "@game-designer/sdk";

export type SlotPhase =
  | "loading"
  | "ready"
  | "spinning"
  | "result"
  | "insufficient_balance"
  | "error";

export interface SlotGameState {
  phase: SlotPhase;
  balance: number;
  wager: number;
  minWager: number;
  maxWager: number;
  lastResult: SpinResult | null;
  errorMessage: string | null;
}

export function createSlotGameState(config?: {
  balance?: number;
  minWager?: number;
  maxWager?: number;
}): SlotGameState {
  return {
    phase: "loading",
    balance: config?.balance ?? 0,
    wager: config?.minWager ?? 10,
    minWager: config?.minWager ?? 1,
    maxWager: config?.maxWager ?? 100,
    lastResult: null,
    errorMessage: null,
  };
}

export function setReady(state: SlotGameState, balance: number, config: SlotConfigResponse): SlotGameState {
  return {
    ...state,
    phase: "ready",
    balance,
    wager: Math.max(config.minWager, Math.min(state.wager, config.maxWager)),
    minWager: config.minWager,
    maxWager: config.maxWager,
    lastResult: null,
    errorMessage: null,
  };
}

export function setWager(state: SlotGameState, wager: number): SlotGameState {
  const clamped = Math.max(state.minWager, Math.min(wager, state.maxWager));
  return { ...state, wager: clamped };
}

export function startSpin(state: SlotGameState): SlotGameState {
  if (state.phase !== "ready") return state;
  if (state.balance < state.wager) {
    return { ...state, phase: "insufficient_balance", errorMessage: "Insufficient balance for this wager" };
  }
  return { ...state, phase: "spinning", errorMessage: null };
}

export function completeSpin(state: SlotGameState, result: SpinResult): SlotGameState {
  return {
    ...state,
    phase: "result",
    balance: result.balance,
    lastResult: result,
    errorMessage: null,
  };
}

export function spinError(state: SlotGameState, message: string): SlotGameState {
  return { ...state, phase: "error", errorMessage: message };
}

export function resetFromResult(state: SlotGameState): SlotGameState {
  if (state.phase === "result" || state.phase === "error" || state.phase === "insufficient_balance") {
    const nextPhase = state.balance >= state.wager ? "ready" : "insufficient_balance";
    return { ...state, phase: nextPhase, errorMessage: nextPhase === "insufficient_balance" ? "Insufficient balance" : null };
  }
  return state;
}
