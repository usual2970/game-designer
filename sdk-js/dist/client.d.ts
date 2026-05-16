import type { SessionRequest, SessionResponse, ProfileResponse, UpdateProfileRequest, SaveGameStateRequest, GameStateResponse, SubmitScoreRequest, SubmitScoreResponse, LeaderboardResponse } from "./types";
export interface GameDesignerConfig {
    baseUrl: string;
}
export declare class GameDesignerClient {
    private baseUrl;
    private token;
    constructor(config: GameDesignerConfig);
    createOrResumeSession(request: SessionRequest): Promise<SessionResponse>;
    getPlayerProfile(): Promise<ProfileResponse>;
    updatePlayerProfile(request: UpdateProfileRequest): Promise<ProfileResponse>;
    saveGameState(request: SaveGameStateRequest): Promise<GameStateResponse>;
    getGameState(): Promise<GameStateResponse | null>;
    submitScore(request: SubmitScoreRequest): Promise<SubmitScoreResponse>;
    getLeaderboard(options?: {
        limit?: number;
        offset?: number;
    }): Promise<LeaderboardResponse>;
    setToken(token: string): void;
    getToken(): string | null;
    private request;
}
