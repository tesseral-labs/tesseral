import React, { PropsWithChildren, useEffect, useRef } from "react";
import { Outlet } from "react-router";

import { useDarkMode } from "@/lib/dark-mode";
import {
  ProjectSettingsProvider,
  useProjectSettings,
} from "@/lib/project-settings";

import { SidebarInset, SidebarProvider } from "../ui/sidebar";
import { VaultSidebar } from "./VaultSidebar";

export function Page() {
  return (
    <ProjectSettingsProvider>
      <PageOuter>
        <SidebarProvider>
          <VaultSidebar />
          <SidebarInset>
            <PageInner />
          </SidebarInset>
        </SidebarProvider>
      </PageOuter>
    </ProjectSettingsProvider>
  );
}

function PageInner() {
  return (
    <div className="container p-4 lg:p-8 w-full bg-background">
      <Outlet />
    </div>
  );
}

function PageOuter({ children }: PropsWithChildren) {
  const rootRef = useRef<HTMLDivElement>(
    document.getElementById("react-root") as HTMLDivElement,
  );
  const darkMode = useDarkMode();
  const projectSettings = useProjectSettings();

  useEffect(() => {
    if (rootRef.current) {
      if (darkMode) {
        rootRef.current.classList.add("dark");
        if (projectSettings?.darkModePrimaryColor) {
          rootRef.current.style.setProperty(
            "--primary",
            projectSettings.darkModePrimaryColor,
          );
        }
      } else {
        rootRef.current.classList.remove("dark");
        if (projectSettings?.primaryColor) {
          rootRef.current.style.setProperty(
            "--primary",
            projectSettings.primaryColor,
          );
        }
      }
    }
  }, [darkMode, projectSettings, rootRef]);

  return <>{children}</>;
}

export function PageContent({ children }: PropsWithChildren) {
  return <div className="space-y-4 lg:space-y-8 w-full">{children}</div>;
}
