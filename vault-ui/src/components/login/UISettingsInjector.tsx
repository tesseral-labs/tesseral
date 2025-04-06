import { useQuery } from "@connectrpc/connect-query";
import React, { useEffect, useRef } from "react";

import { getSettings } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useDarkMode } from "@/lib/dark-mode";
import { hexToHSL, isColorDark } from "@/lib/utils";
import { useProjectSettings } from "@/lib/project-settings";

export function UISettingsInjector({
  children,
}: {
  children?: React.ReactNode;
}) {
  const root = useRef<HTMLDivElement>(null);
  const settings = useProjectSettings()
  const darkMode = useDarkMode();

  useEffect(() => {
    if (!root.current) {
      return;
    }

    if (!darkMode && settings.primaryColor) {
      const foreground = isColorDark(settings.primaryColor)
        ? "0 0% 100%"
        : "0 0% 0%";

      root.current.style.setProperty(
        "--primary",
        hexToHSL(settings.primaryColor),
      );
      root.current.style.setProperty("--primary-foreground", foreground);
    }

    if (settings.darkModePrimaryColor && darkMode) {
      const darkForeground = isColorDark(
        settings.darkModePrimaryColor,
      )
        ? "0 0% 100%"
        : "0 0% 0%";

      root.current.style.setProperty(
        "--primary",
        hexToHSL(settings.darkModePrimaryColor),
      );
      root.current.style.setProperty("--primary-foreground", darkForeground);
    }
  }, [darkMode, settings]);

  return <div ref={root}>{children}</div>;
}
