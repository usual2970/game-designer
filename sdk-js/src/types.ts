export interface SessionRequest {
  playerId: string;
  nickname?: string;
  avatarUrl?: string;
}

export interface SessionResponse {
  sessionId: string;
  token: string;
  playerId: string;
  isNew: boolean;
  expiresAt: string;
}

export interface ProfileResponse {
  playerId: string;
  nickname: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateProfileRequest {
  nickname?: string;
  avatarUrl?: string;
}

export interface SymbolDefinition {
  name: string;
  payoutMultiplier: number;
}

export interface PaylineDefinition {
  id: number;
  positions: number[];
}

export interface SlotConfigResponse {
  reels: number;
  rows: number;
  paylines: PaylineDefinition[];
  symbols: SymbolDefinition[];
  minWager: number;
  maxWager: number;
  defaultBalance: number;
}

export interface BalanceResponse {
  balance: number;
}

export interface SpinRequest {
  wager: number;
}

export interface PaylineWin {
  paylineId: number;
  symbol: string;
  count: number;
  payout: number;
}

export interface SpinResult {
  spinId: string;
  wager: number;
  reels: string[][];
  paylineWins: PaylineWin[];
  totalPayout: number;
  balance: number;
}

export interface SpinHistoryEntry {
  spinId: string;
  wager: number;
  totalPayout: number;
  balance: number;
  reels: string[][];
  paylineWins: PaylineWin[];
  spunAt: string;
}

export interface SpinHistoryResponse {
  entries: SpinHistoryEntry[];
  total: number;
}

export interface SlotLeaderboardEntry {
  rank: number;
  playerId: string;
  nickname: string;
  balance: number;
  updatedAt: string;
}

export interface SlotLeaderboardResponse {
  entries: SlotLeaderboardEntry[];
  total: number;
}

export type ErrorCode =
  | "INVALID_PARAMETERS"
  | "UNAUTHORIZED"
  | "NOT_FOUND"
  | "SESSION_EXPIRED"
  | "INSUFFICIENT_BALANCE"
  | "INTERNAL_ERROR";

export interface ApiErrorResponse {
  error: string;
  code: ErrorCode;
  details?: Record<string, unknown>;
}
