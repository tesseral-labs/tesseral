import React, { useEffect } from "react";

import { useDarkMode } from "@/lib/dark-mode";
import { useProjectSettings } from "@/lib/project-settings";
import { hexToHSL, isColorDark } from "@/lib/utils";

export function UISettingsInjector({
  children,
}: {
  children?: React.ReactNode;
}) {
  const settings = useProjectSettings();
  const darkMode = useDarkMode();

  useEffect(() => {
    if (darkMode) {
      document.body.classList.add("dark");
    } else {
      document.body.classList.remove("dark");
    }

    if (!darkMode && settings.primaryColor) {
      const foreground = isColorDark(settings.primaryColor)
        ? "0 0% 100%"
        : "0 0% 0%";

      document.body.style.setProperty(
        "--primary",
        hexToHSL(settings.primaryColor),
      );
      document.body.style.setProperty("--primary-foreground", foreground);
    }

    if (settings.darkModePrimaryColor && darkMode) {
      const darkForeground = isColorDark(settings.darkModePrimaryColor)
        ? "0 0% 100%"
        : "0 0% 0%";

      document.body.style.setProperty(
        "--primary",
        hexToHSL(settings.darkModePrimaryColor),
      );
      document.body.style.setProperty("--primary-foreground", darkForeground);
    }
  }, [darkMode, settings]);

  return <>{children}</>;
}
