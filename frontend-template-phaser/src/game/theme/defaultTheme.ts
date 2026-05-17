import type { ThemeConfig } from "./themeSchema";

export const defaultTheme: ThemeConfig = {
  name: "default",
  title: "Slot Machine",
  subtitle: "",
  colors: {
    background: "#1a1a2e",
    primary: "#00cc66",
    secondary: "#00ccff",
    accent: "#ffcc00",
    text: "#ffffff",
    textMuted: "#aaaaaa",
  },
  symbols: {
    cherry: "symbol_cherry",
    lemon: "symbol_lemon",
    orange: "symbol_orange",
    plum: "symbol_plum",
    bell: "symbol_bell",
    seven: "symbol_seven",
    bar: "symbol_bar",
  },
  sounds: {},
  layout: {
    safeAreaTop: 20,
    safeAreaBottom: 20,
    minTapTarget: 48,
  },
};
