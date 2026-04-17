export type ThemeName = "violet" | "forest" | "sunset" | "ocean" | "paper";

type ThemeTokens = {
  bg: string;
  bgStrong: string;
  panel: string;
  panel2: string;
  border: string;
  borderStrong: string;
  primary: string;
  accent: string;
  muted: string;
  fg: string;
  ok: string;
  error: string;
  shadow: string;
  colorScheme: "dark" | "light";
};

const themes: Record<ThemeName, ThemeTokens> = {
  violet: {
    bg: "#0f172a",
    bgStrong: "#0b1120",
    panel: "rgba(17, 24, 39, 0.92)",
    panel2: "rgba(31, 41, 55, 0.82)",
    border: "#374151",
    borderStrong: "rgba(167, 139, 250, 0.42)",
    primary: "#7c3aed",
    accent: "#a78bfa",
    muted: "#94a3b8",
    fg: "#f9fafb",
    ok: "#10b981",
    error: "#ef4444",
    shadow: "0 24px 80px rgba(2, 6, 23, 0.38)",
    colorScheme: "dark",
  },
  forest: {
    bg: "#071a17",
    bgStrong: "#041311",
    panel: "rgba(9, 24, 21, 0.92)",
    panel2: "rgba(15, 41, 35, 0.82)",
    border: "#29443f",
    borderStrong: "rgba(45, 212, 191, 0.42)",
    primary: "#0f766e",
    accent: "#2dd4bf",
    muted: "#9ca3af",
    fg: "#ecfdf5",
    ok: "#22c55e",
    error: "#ef4444",
    shadow: "0 24px 80px rgba(1, 10, 9, 0.38)",
    colorScheme: "dark",
  },
  sunset: {
    bg: "#1c0f0a",
    bgStrong: "#130905",
    panel: "rgba(33, 17, 12, 0.94)",
    panel2: "rgba(52, 28, 18, 0.82)",
    border: "#5c4033",
    borderStrong: "rgba(251, 146, 60, 0.44)",
    primary: "#c2410c",
    accent: "#fb923c",
    muted: "#d6b8a6",
    fg: "#fffbeb",
    ok: "#16a34a",
    error: "#dc2626",
    shadow: "0 24px 80px rgba(20, 8, 2, 0.4)",
    colorScheme: "dark",
  },
  ocean: {
    bg: "#081824",
    bgStrong: "#06111a",
    panel: "rgba(11, 25, 38, 0.92)",
    panel2: "rgba(18, 40, 56, 0.82)",
    border: "#334155",
    borderStrong: "rgba(56, 189, 248, 0.42)",
    primary: "#0369a1",
    accent: "#38bdf8",
    muted: "#94a3b8",
    fg: "#f8fafc",
    ok: "#10b981",
    error: "#ef4444",
    shadow: "0 24px 80px rgba(3, 9, 16, 0.38)",
    colorScheme: "dark",
  },
  paper: {
    bg: "#f8fafc",
    bgStrong: "#e2e8f0",
    panel: "rgba(255, 255, 255, 0.96)",
    panel2: "rgba(248, 250, 252, 0.94)",
    border: "#cbd5e1",
    borderStrong: "rgba(29, 78, 216, 0.28)",
    primary: "#1d4ed8",
    accent: "#0f766e",
    muted: "#64748b",
    fg: "#111827",
    ok: "#15803d",
    error: "#b91c1c",
    shadow: "0 24px 80px rgba(148, 163, 184, 0.3)",
    colorScheme: "light",
  },
};

export function applyTheme(themeName: string): void {
  const theme = themes[(themeName as ThemeName) in themes ? (themeName as ThemeName) : "violet"];
  const root = document.documentElement;
  root.style.setProperty("--bg", theme.bg);
  root.style.setProperty("--bg-strong", theme.bgStrong);
  root.style.setProperty("--panel", theme.panel);
  root.style.setProperty("--panel-2", theme.panel2);
  root.style.setProperty("--border", theme.border);
  root.style.setProperty("--border-strong", theme.borderStrong);
  root.style.setProperty("--primary", theme.primary);
  root.style.setProperty("--accent", theme.accent);
  root.style.setProperty("--muted", theme.muted);
  root.style.setProperty("--fg", theme.fg);
  root.style.setProperty("--ok", theme.ok);
  root.style.setProperty("--error", theme.error);
  root.style.setProperty("--shadow", theme.shadow);
  root.style.setProperty("color-scheme", theme.colorScheme);
  document.body.dataset.theme = themeName;
}

export function themeNames(): ThemeName[] {
  return Object.keys(themes) as ThemeName[];
}
