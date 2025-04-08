import React from "react";
import { Outlet } from "react-router";

import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ProjectSettingsProvider } from "@/lib/project-settings";
import { DashboardSidebar } from "@/pages/dashboard/DashboardSidebar";

export function DashboardLayout() {
  return (
    <ProjectSettingsProvider>
      <SidebarProvider>
        <DashboardSidebar />
        <SidebarInset>
          <main className="bg-background w-full">
            <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 py-8">
              <Outlet />
            </div>
          </main>
        </SidebarInset>
      </SidebarProvider>
    </ProjectSettingsProvider>
  );
}
