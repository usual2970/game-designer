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
export interface SaveGameStateRequest {
    data: Record<string, unknown>;
    checkpoint?: string;
}
export interface GameStateResponse {
    data: Record<string, unknown>;
    checkpoint?: string;
    savedAt: string;
}
export interface SubmitScoreRequest {
    score: number;
    metadata?: Record<string, unknown>;
}
export interface SubmitScoreResponse {
    accepted: boolean;
    rank: number | null;
    bestScore: number;
    isNewBest: boolean;
}
export interface LeaderboardEntry {
    rank: number;
    playerId: string;
    nickname: string;
    score: number;
    achievedAt: string;
}
export interface LeaderboardResponse {
    entries: LeaderboardEntry[];
    total: number;
}
export type ErrorCode = "INVALID_PARAMETERS" | "UNAUTHORIZED" | "NOT_FOUND" | "SESSION_EXPIRED" | "INTERNAL_ERROR";
export interface ApiErrorResponse {
    error: string;
    code: ErrorCode;
    details?: Record<string, unknown>;
}
