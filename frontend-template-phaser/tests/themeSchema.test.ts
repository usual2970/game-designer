import { describe, it, expect } from "vitest";
import { validateTheme } from "../src/game/theme/themeSchema";
import { defaultTheme } from "../src/game/theme/defaultTheme";

describe("themeSchema", () => {
  it("validates the default theme as valid", () => {
    const result = validateTheme(defaultTheme);
    expect(result.valid).toBe(true);
    expect(result.errors).toHaveLength(0);
  });

  it("reports missing required fields", () => {
    const result = validateTheme({});
    expect(result.valid).toBe(false);
    expect(result.errors).toContain("Theme name is required");
    expect(result.errors).toContain("Theme title is required");
    expect(result.errors).toContain("Colors section is required");
  });

  it("reports missing required symbol assets", () => {
    const result = validateTheme({
      name: "test",
      title: "Test",
      colors: { background: "#000", primary: "#fff", text: "#fff", secondary: "", accent: "", textMuted: "" },
      symbols: { cherry: "sym" },
    });
    expect(result.valid).toBe(false);
    expect(result.errors).toContain("Missing required symbol asset: lemon");
  });

  it("warns about small tap targets", () => {
    const result = validateTheme({
      name: "test",
      title: "Test",
      colors: { background: "#000", primary: "#fff", text: "#fff", secondary: "", accent: "", textMuted: "" },
      layout: { safeAreaTop: 0, safeAreaBottom: 0, minTapTarget: 30 },
    });
    expect(result.warnings.length).toBeGreaterThan(0);
    expect(result.warnings.some((w) => w.includes("minTapTarget"))).toBe(true);
  });

  it("warns about non-assets sound paths", () => {
    const result = validateTheme({
      name: "test",
      title: "Test",
      colors: { background: "#000", primary: "#fff", text: "#fff", secondary: "", accent: "", textMuted: "" },
      sounds: { spin: "sounds/spin.mp3" },
    });
    expect(result.warnings.some((w) => w.includes("Sound"))).toBe(true);
    expect(result.warnings.some((w) => w.includes("spin"))).toBe(true);
  });

  it("accepts valid hex and named colors", () => {
    const result = validateTheme({
      name: "test",
      title: "Test",
      colors: { background: "#1a1a2e", primary: "green", text: "rgb(255,255,255)", secondary: "hsl(0,0%,0%)", accent: "#fff", textMuted: "#ccc" },
    });
    expect(result.errors).toHaveLength(0);
  });
});
