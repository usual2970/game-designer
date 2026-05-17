import { Scene } from "phaser";
import { createGameDesignerClient } from "../services/gameDesignerClient";
import {
  createSlotGameState,
  setReady,
  setWager,
  startSpin,
  completeSpin,
  spinError,
  resetFromResult,
} from "../services/slotGameState";
import type { SlotGameState } from "../services/slotGameState";

const SYMBOL_NAMES = ["cherry", "lemon", "orange", "plum", "bell", "seven", "bar"];
const REEL_COUNT = 3;
const ROW_COUNT = 3;
const CELL_SIZE = 64;
const CELL_GAP = 4;
const WAGER_DEFAULT = 10;

export class SlotScene extends Scene {
  private client: ReturnType<typeof createGameDesignerClient> | null = null;
  private state: SlotGameState = createSlotGameState();
  private reelContainers: Phaser.GameObjects.Container[] = [];
  private balanceText!: Phaser.GameObjects.Text;
  private wagerText!: Phaser.GameObjects.Text;
  private statusText!: Phaser.GameObjects.Text;
  private spinButton!: Phaser.GameObjects.Image;

  constructor() {
    super({ key: "SlotScene" });
  }

  async create(): Promise<void> {
    const cx = this.cameras.main.centerX;

    this.add.text(cx, 30, "Slot Machine", {
      fontSize: "28px",
      color: "#ffffff",
      fontStyle: "bold",
    }).setOrigin(0.5);

    const gridWidth = REEL_COUNT * (CELL_SIZE + CELL_GAP) - CELL_GAP;
    const gridX = cx - gridWidth / 2;
    const gridY = 80;

    for (let col = 0; col < REEL_COUNT; col++) {
      const container = this.add.container(gridX + col * (CELL_SIZE + CELL_GAP), gridY);
      for (let row = 0; row < ROW_COUNT; row++) {
        const symName = SYMBOL_NAMES[(col * ROW_COUNT + row) % SYMBOL_NAMES.length];
        const img = this.add.image(0, row * (CELL_SIZE + CELL_GAP), `symbol_${symName}`);
        img.setOrigin(0);
        container.add(img);
      }
      this.reelContainers.push(container);
    }

    const infoY = gridY + ROW_COUNT * (CELL_SIZE + CELL_GAP);

    this.balanceText = this.add.text(cx, infoY + 20, "Balance: --", {
      fontSize: "20px",
      color: "#00ccff",
    }).setOrigin(0.5);

    this.wagerText = this.add.text(cx, infoY + 50, `Wager: ${this.state.wager}`, {
      fontSize: "18px",
      color: "#ffcc00",
    }).setOrigin(0.5);

    this.spinButton = this.add.image(cx, infoY + 110, "btn_spin");
    this.spinButton.setInteractive({ useHandCursor: true });
    this.spinButton.on("pointerdown", () => this.handleSpin());

    this.add.text(cx, infoY + 110, "SPIN", {
      fontSize: "18px",
      color: "#ffffff",
      fontStyle: "bold",
    }).setOrigin(0.5).setDepth(1);

    this.statusText = this.add.text(cx, infoY + 160, "Connecting...", {
      fontSize: "14px",
      color: "#aaaaaa",
    }).setOrigin(0.5);

    await this.initClient();
  }

  private async initClient(): Promise<void> {
    const baseUrl = this.getBaseUrl();
    this.client = createGameDesignerClient(baseUrl);

    try {
      await this.client.createOrResumeSession({
        playerId: this.getPlayerId(),
        nickname: `Player-${this.getPlayerId().slice(0, 4)}`,
      });

      const config = await this.client.getSlotConfig();
      const balance = await this.client.getBalance();

      this.state = setReady(this.state, balance.balance, config);
      this.updateUI();
      this.statusText.setText("Ready");
    } catch {
      this.state = spinError(this.state, "Connection failed. Check server URL.");
      this.statusText.setText(this.state.errorMessage ?? "Connection failed");
    }
  }

  private async handleSpin(): Promise<void> {
    if (!this.client) return;

    this.state = startSpin(this.state);

    if (this.state.phase === "insufficient_balance" || this.state.phase === "error") {
      this.statusText.setText(this.state.errorMessage ?? "Cannot spin");
      return;
    }

    this.statusText.setText("Spinning...");

    try {
      const result = await this.client.spin({ wager: this.state.wager });
      this.state = completeSpin(this.state, result);
      this.animateReels(result.reels);

      if (result.totalPayout > 0) {
        this.statusText.setText(`Win! +${result.totalPayout} credits`);
      } else {
        this.statusText.setText("No win. Try again!");
      }
    } catch {
      this.state = spinError(this.state, "Spin failed. Check connection.");
      this.statusText.setText(this.state.errorMessage ?? "Spin failed");
    } finally {
      this.state = resetFromResult(this.state);
      this.updateUI();
    }
  }

  private animateReels(reels: string[][]): void {
    for (let col = 0; col < Math.min(reels.length, REEL_COUNT); col++) {
      const container = this.reelContainers[col];
      if (!container) continue;
      const symbols = reels[col];

      for (let row = 0; row < Math.min(symbols.length, ROW_COUNT); row++) {
        const symKey = `symbol_${symbols[row].toLowerCase()}`;
        const existing = container.getAt(row) as Phaser.GameObjects.Image;
        if (existing) {
          existing.setTexture(symKey);
        }
      }
    }
  }

  private updateUI(): void {
    this.balanceText.setText(`Balance: ${this.state.balance}`);
    this.wagerText.setText(`Wager: ${this.state.wager}`);
  }

  private getBaseUrl(): string {
    const meta = document.querySelector('meta[name="game-server-url"]');
    return meta?.getAttribute("content") || "http://localhost:8080";
  }

  private getPlayerId(): string {
    const params = new URLSearchParams(window.location.search);
    return params.get("playerId") || `player-${Date.now()}`;
  }
}
