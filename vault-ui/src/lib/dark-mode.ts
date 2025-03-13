import { useEffect, useState } from "react";

const useDarkMode = () => {
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

  return isDarkMode;
};

export default useDarkMode;
