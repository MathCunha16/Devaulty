import React, { useEffect, useState } from "react";
import { type Theme, ThemeContext } from "../hooks/useTheme";

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setTheme] = useState<Theme>(() => {
    try {
      const stored = localStorage.getItem("devaulty-theme");
      if (stored === "light" || stored === "dark") return stored;
    } catch {
      // Ignore storage errors in restricted/private browsing mode
    }
    return "dark"; // dark mode by default for developer tool aesthetic
  });

  useEffect(() => {
    const root = window.document.documentElement;
    root.classList.remove("light", "dark");
    root.classList.add(theme);
    try {
      localStorage.setItem("devaulty-theme", theme);
    } catch {
      // Ignore storage errors in restricted/private browsing mode
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === "light" ? "dark" : "light"));
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

