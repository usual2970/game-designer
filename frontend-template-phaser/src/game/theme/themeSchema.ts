export interface ThemeConfig {
  name: string;
  title: string;
  subtitle?: string;
  colors: {
    background: string;
    primary: string;
    secondary: string;
    accent: string;
    text: string;
    textMuted: string;
  };
  symbols: Record<string, string>;
  sounds: {
    spin?: string;
    win?: string;
    buttonClick?: string;
  };
  layout: {
    safeAreaTop: number;
    safeAreaBottom: number;
    minTapTarget: number;
  };
}

export interface ThemeValidationResult {
  valid: boolean;
  errors: string[];
  warnings: string[];
}

const REQUIRED_SYMBOLS = ["cherry", "lemon", "orange", "plum", "bell", "seven", "bar"];

export function validateTheme(theme: Partial<ThemeConfig>): ThemeValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!theme.name || theme.name.trim() === "") {
    errors.push("Theme name is required");
  }

  if (!theme.title || theme.title.trim() === "") {
    errors.push("Theme title is required");
  }

  if (!theme.colors) {
    errors.push("Colors section is required");
  } else {
    for (const key of ["background", "primary", "text"] as const) {
      if (!theme.colors[key]) {
        errors.push(`colors.${key} is required`);
      }
    }
    for (const [key, value] of Object.entries(theme.colors)) {
      if (value && !isValidColor(value)) {
        errors.push(`colors.${key} "${value}" is not a valid CSS color`);
      }
    }
  }

  if (theme.symbols) {
    for (const required of REQUIRED_SYMBOLS) {
      if (!theme.symbols[required]) {
        errors.push(`Missing required symbol asset: ${required}`);
      }
    }
  }

  if (theme.sounds) {
    for (const [key, path] of Object.entries(theme.sounds)) {
      if (path && !path.startsWith("assets/")) {
        warnings.push(`Sound ${key} path should be relative to the assets directory`);
      }
    }
  }

  if (theme.layout) {
    if (theme.layout.minTapTarget < 44) {
      warnings.push("minTapTarget should be at least 44px for mobile accessibility");
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

function isValidColor(value: string): boolean {
  return /^(#[0-9a-fA-F]{3,8}|rgb|hsl)/.test(value) || /^[a-zA-Z]+$/.test(value);
}
