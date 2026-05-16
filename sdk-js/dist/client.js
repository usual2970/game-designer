import { ApiError } from "./error";
export class GameDesignerClient {
    constructor(config) {
        this.token = null;
        this.baseUrl = config.baseUrl.replace(/\/+$/, "");
    }
    async createOrResumeSession(request) {
        const response = await this.request("POST", "/api/v1/session", request);
        const data = await response.json();
        this.token = data.token;
        return data;
    }
    async getPlayerProfile() {
        const response = await this.request("GET", "/api/v1/profile");
        return response.json();
    }
    async updatePlayerProfile(request) {
        const response = await this.request("PUT", "/api/v1/profile", request);
        return response.json();
    }
    async saveGameState(request) {
        const response = await this.request("PUT", "/api/v1/game-state", request);
        return response.json();
    }
    async getGameState() {
        const response = await this.request("GET", "/api/v1/game-state");
        if (response.status === 204) {
            return null;
        }
        return response.json();
    }
    async submitScore(request) {
        const response = await this.request("POST", "/api/v1/scores", request);
        return response.json();
    }
    async getLeaderboard(options) {
        const params = new URLSearchParams();
        if (options?.limit !== undefined)
            params.set("limit", String(options.limit));
        if (options?.offset !== undefined)
            params.set("offset", String(options.offset));
        const qs = params.toString() ? `?${params.toString()}` : "";
        const response = await this.request("GET", `/api/v1/leaderboard${qs}`);
        return response.json();
    }
    setToken(token) {
        this.token = token;
    }
    getToken() {
        return this.token;
    }
    async request(method, path, body) {
        const headers = {
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
            const errBody = await response.json();
            throw new ApiError(errBody);
        }
        return response;
    }
}
