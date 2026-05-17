import { GameDesignerClient } from "@game-designer/sdk";

export function createGameDesignerClient(baseUrl: string): GameDesignerClient {
  return new GameDesignerClient({ baseUrl });
}

export type { GameDesignerClient };
