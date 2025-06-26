import React, { PropsWithChildren } from "react";
import { Outlet } from "react-router";

import { ProjectSettingsProvider } from "@/lib/project-settings";

import { UISettingsInjector } from "../core/UISettingsInjector";
import { SidebarInset, SidebarProvider } from "../ui/sidebar";
import { VaultSidebar } from "./VaultSidebar";

export function Page() {
  return (
    <ProjectSettingsProvider>
      <UISettingsInjector>
        <SidebarProvider>
          <VaultSidebar />
          <SidebarInset>
            <PageInner />
          </SidebarInset>
        </SidebarProvider>
      </UISettingsInjector>
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

export function PageContent({ children }: PropsWithChildren) {
  return <div className="space-y-4 lg:space-y-8 w-full">{children}</div>;
}
