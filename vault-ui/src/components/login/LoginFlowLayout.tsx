import React from "react";
import { Outlet } from "react-router";

import { UISettingsInjector } from "@/components/core/UISettingsInjector";
import { ProjectSettingsProvider } from "@/lib/project-settings";

export function LoginFlowLayout() {
  return (
    <ProjectSettingsProvider>
      <UISettingsInjector>
        <div className="bg-background w-full min-h-screen mx-auto flex flex-col justify-center items-center py-8">
          <div className="max-w-sm w-full mx-auto">
            <Outlet />
          </div>
        </div>
      </UISettingsInjector>
    </ProjectSettingsProvider>
  );
}
