import React, { FC, useEffect, useState } from "react";
import { Helmet } from "react-helmet";

import { useDarkMode } from "@/lib/dark-mode";
import { useLayout, useSettings } from "@/lib/settings";
import { cn, hexToHSL, isColorDark } from "@/lib/utils";
import { LoginLayouts } from "@/lib/views";

import { CenteredLayout } from "./layouts/centered";
import { SideBySideLayout } from "./layouts/side-by-side";
import { Toaster } from "./ui/sonner";

const layoutMap: Record<string, FC> = {
  [`${LoginLayouts.Centered}`]: CenteredLayout,
  [`${LoginLayouts.SideBySide}`]: SideBySideLayout,
};

export function Page() {
  const isDarkMode = useDarkMode();
  const layout = useLayout();
  const settings = useSettings();

  const [favicon, setFavicon] = useState<string>("/apple-touch-icon.png");
  const Layout =
    layout && layoutMap[layout] ? layoutMap[layout] : CenteredLayout;

  // eslint-disable-next-line func-style
  const applyTheme = () => {
    const root = document.documentElement;
    const darkRoot = document.querySelector(".dark");

    const primary = settings?.primaryColor;
    const darkPrimary = settings?.darkModePrimaryColor;

    if (primary) {
      const foreground = isColorDark(primary) ? "0 0% 100%" : "0 0% 0%";

      root.style.setProperty("--primary", hexToHSL(primary));
      root.style.setProperty("--primary-foreground", foreground);
    }

    if (darkPrimary && darkRoot) {
      const root = darkRoot as HTMLElement;
      const darkForeground = isColorDark(darkPrimary) ? "0 0% 100%" : "0 0% 0%";

      root.style.setProperty("--primary", hexToHSL(darkPrimary));
      root.style.setProperty("--primary-foreground", darkForeground);
    }
  };

  useEffect(() => {
    if (settings) {
      applyTheme();
    }

    if (settings?.faviconUrl) {
      void (async () => {
        try {
          // Check if the favicon exists before setting it
          const faviconCheck = await fetch(settings?.faviconUrl, {
            method: "HEAD",
          });

          setFavicon(
            faviconCheck.ok ? settings?.faviconUrl : "/apple-touch-icon.png",
          );
        } catch {
          setFavicon("/apple-touch-icon.png");
        }
      })();
    }
  }, [settings]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    applyTheme();
  }, [isDarkMode]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div
      className={cn(
        "min-h-screen w-screen",
        isDarkMode && settings?.detectDarkModeEnabled
          ? "dark"
          : "light bg-body",
      )}
    >
      <div className="bg-background min-h-screen w-full">
        <Helmet>
          <link rel="icon" href={favicon} />
          <link rel="apple-touch-icon" href={favicon} />
        </Helmet>

        <Layout />

        <Toaster
          position={
            layout === LoginLayouts.SideBySide ? "top-right" : "top-center"
          }
        />
      </div>
    </div>
  );
}
