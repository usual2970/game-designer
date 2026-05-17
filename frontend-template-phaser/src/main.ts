import { Game, AUTO, Scale } from "phaser";
import { BootScene } from "./game/scenes/BootScene";
import { PreloadScene } from "./game/scenes/PreloadScene";
import { SlotScene } from "./game/scenes/SlotScene";

const GAME_WIDTH = 375;
const GAME_HEIGHT = 667;

const config = {
  type: AUTO,
  width: GAME_WIDTH,
  height: GAME_HEIGHT,
  parent: "game-container",
  backgroundColor: "#1a1a2e",
  scale: {
    mode: Scale.FIT,
    autoCenter: Scale.CENTER_BOTH,
  },
  scene: [BootScene, PreloadScene, SlotScene],
};

const game = new Game(config);

export { game };
