import { useEffect, useState } from "react";
import { useProjectSettings } from "@/lib/project-settings";

export function useDarkMode ()  {
  const settings = useProjectSettings()

  const [isDarkMode, setIsDarkMode] = useState(() => {
    // Get the initial dark mode state
    const matcher =
      window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)");
    return matcher ? matcher.matches : false;
  });

  useEffect(() => {
    const matcher =
      window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)");

    if (matcher) {
      // eslint-disable-next-line func-style
      const handleDarkModeChange = (event: MediaQueryListEvent) => {
        setIsDarkMode(event.matches);
      };

      matcher.addEventListener("change", handleDarkModeChange);

      // Cleanup listener on unmount
      return () => {
        matcher.removeEventListener("change", handleDarkModeChange);
      };
    }
  }, []);

  return settings.detectDarkModeEnabled ? isDarkMode : false;
}
