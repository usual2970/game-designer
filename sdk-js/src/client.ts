import type {
  SessionRequest,
  SessionResponse,
  ProfileResponse,
  UpdateProfileRequest,
  SlotConfigResponse,
  BalanceResponse,
  SpinRequest,
  SpinResult,
  SpinHistoryResponse,
  SlotLeaderboardResponse,
  ApiErrorResponse,
} from "./types";
import { ApiError } from "./error";

export interface GameDesignerConfig {
  baseUrl: string;
}

export class GameDesignerClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(config: GameDesignerConfig) {
    this.baseUrl = config.baseUrl.replace(/\/+$/, "");
  }

  async createOrResumeSession(
    request: SessionRequest
  ): Promise<SessionResponse> {
    const response = await this.request("POST", "/api/v1/session", request);
    const data: SessionResponse = await response.json();
    this.token = data.token;
    return data;
  }

  async getPlayerProfile(): Promise<ProfileResponse> {
    const response = await this.request("GET", "/api/v1/profile");
    return response.json();
  }

  async updatePlayerProfile(
    request: UpdateProfileRequest
  ): Promise<ProfileResponse> {
    const response = await this.request("PUT", "/api/v1/profile", request);
    return response.json();
  }

  async getSlotConfig(): Promise<SlotConfigResponse> {
    const response = await this.request("GET", "/api/v1/slot/config");
    return response.json();
  }

  async getBalance(): Promise<BalanceResponse> {
    const response = await this.request("GET", "/api/v1/balance");
    return response.json();
  }

  async spin(request: SpinRequest): Promise<SpinResult> {
    const response = await this.request("POST", "/api/v1/spin", request);
    return response.json();
  }

  async getSpinHistory(
    options?: { limit?: number; offset?: number }
  ): Promise<SpinHistoryResponse> {
    const params = new URLSearchParams();
    if (options?.limit !== undefined) params.set("limit", String(options.limit));
    if (options?.offset !== undefined) params.set("offset", String(options.offset));
    const qs = params.toString() ? `?${params.toString()}` : "";
    const response = await this.request("GET", `/api/v1/spin/history${qs}`);
    return response.json();
  }

  async getSlotLeaderboard(
    options?: { limit?: number; offset?: number }
  ): Promise<SlotLeaderboardResponse> {
    const params = new URLSearchParams();
    if (options?.limit !== undefined) params.set("limit", String(options.limit));
    if (options?.offset !== undefined) params.set("offset", String(options.offset));
    const qs = params.toString() ? `?${params.toString()}` : "";
    const response = await this.request("GET", `/api/v1/leaderboard${qs}`);
    return response.json();
  }

  setToken(token: string): void {
    this.token = token;
  }

  getToken(): string | null {
    return this.token;
  }

  private async request(
    method: string,
    path: string,
    body?: unknown
  ): Promise<Response> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this.token) {
      headers["X-Session-Token"] = this.token;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok && response.status !== 204) {
      const errBody: ApiErrorResponse = await response.json();
      throw new ApiError(errBody);
    }

    return response;
  }
}
