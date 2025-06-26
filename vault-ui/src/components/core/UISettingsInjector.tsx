import React, { useEffect } from "react";

import { useDarkMode } from "@/lib/dark-mode";
import { useProjectSettings } from "@/lib/project-settings";
import { isColorDark } from "@/lib/utils";

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

    const primaryColor = darkMode
      ? settings.darkModePrimaryColor || "#ffffff"
      : settings.primaryColor || "#0f172a";

    const primaryForeground = isColorDark(primaryColor) ? "#ffffff" : "#000000";

    document.body.style.setProperty("--primary", primaryColor);
    document.body.style.setProperty("--primary-foreground", primaryForeground);
  }, [darkMode, settings]);

  return <>{children}</>;
}
