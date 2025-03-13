import React, { PropsWithChildren, useEffect, useState } from "react";
import { Helmet } from "react-helmet";

import { useIsMobile } from "@/hooks/use-mobile";
import {
  OrganizationContextProvider,
  ProjectContextProvider,
  UserContextProvider,
  useSession,
} from "@/lib/auth";
import { useDarkMode } from "@/lib/dark-mode";
import { useSettings } from "@/lib/settings";

import { SidebarProvider, SidebarTrigger } from "../ui/sidebar";
import { Toaster } from "../ui/sonner";
import { DashboardSidebar } from "./DashboardSidebar";

export function DashboardPage({ children }: PropsWithChildren) {
  const isDarkMode = useDarkMode();
  const isMobile = useIsMobile();
  const settings = useSettings();
  const session = useSession();

  const [favicon, setFavicon] = useState<string>("/apple-touch-icon.png");

  useEffect(() => {
    if (settings?.faviconUrl) {
      void (async () => {
        // Check if the favicon exists before setting it
        const faviconCheck = await fetch(settings.faviconUrl, {
          method: "HEAD",
        });

        setFavicon(
          faviconCheck.ok ? settings.faviconUrl : "/apple-touch-icon.png",
        );
      })();
    }
  }, [settings]);

  return (
    <div
      className={isDarkMode && settings?.detectDarkModeEnabled ? "dark" : ""}
    >
      <Helmet>
        <link rel="icon" href={favicon} />
        <link rel="apple-touch-icon" href={favicon} />
        <title>{session?.organization?.displayName || "Dashboard"}</title>
      </Helmet>

      <ProjectContextProvider value={session?.project}>
        <OrganizationContextProvider value={session?.organization}>
          <UserContextProvider value={session?.user}>
            <SidebarProvider>
              <DashboardSidebar />
              <main className="min-h-screen w-screen">
                {isMobile && <SidebarTrigger />}
                <div className="bg-background min-h-screen mx-auto items-center">
                  <div className="mx-auto px-6 lg:px-8">
                    <div className="py-8">{children}</div>
                  </div>
                </div>
              </main>
              <Toaster position="top-center" />
            </SidebarProvider>
          </UserContextProvider>
        </OrganizationContextProvider>
      </ProjectContextProvider>
    </div>
  );
}
