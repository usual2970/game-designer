import { Scene } from "phaser";

export class PreloadScene extends Scene {
  constructor() {
    super({ key: "PreloadScene" });
  }

  preload(): void {
    const cx = this.cameras.main.centerX;
    const cy = this.cameras.main.centerY;

    const barWidth = 200;
    const barHeight = 16;
    const barX = cx - barWidth / 2;

    const bg = this.add.graphics();
    bg.fillStyle(0x222244);
    bg.fillRect(barX, cy, barWidth, barHeight);

    const bar = this.add.graphics();
    this.load.on("progress", (value: number) => {
      bar.clear();
      bar.fillStyle(0x00ccff);
      bar.fillRect(barX, cy, barWidth * value, barHeight);
    });

    this.load.on("complete", () => {
      bar.destroy();
      bg.destroy();
    });

    // Load placeholder assets for the slot game.
    // In a real project, replace these with actual game assets.
    this.generatePlaceholderAssets();
  }

  private generatePlaceholderAssets(): void {
    // Generate colored rectangles as placeholder symbols.
    const colors = [0xff4444, 0x44ff44, 0x4444ff, 0xffff44, 0xff44ff, 0x44ffff, 0xff8844];
    const names = ["Cherry", "Lemon", "Orange", "Plum", "Bell", "Seven", "BAR"];

    for (let i = 0; i < names.length; i++) {
      const g = this.add.graphics();
      g.fillStyle(colors[i]);
      g.fillRoundedRect(0, 0, 64, 64, 8);

      const text = this.add.text(32, 32, names[i].slice(0, 2), {
        fontSize: "14px",
        color: "#ffffff",
      });
      text.setOrigin(0.5);

      g.generateTexture(`symbol_${names[i].toLowerCase()}`, 64, 64);
      g.destroy();
      text.destroy();
    }

    // Generate spin button placeholder.
    const btn = this.add.graphics();
    btn.fillStyle(0x00cc66);
    btn.fillRoundedRect(0, 0, 120, 48, 12);
    btn.generateTexture("btn_spin", 120, 48);
    btn.destroy();
  }

  create(): void {
    this.scene.start("SlotScene");
  }
}
