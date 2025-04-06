import React from "react";
import { Outlet } from "react-router";

import { UISettingsInjector } from "@/components/login/UISettingsInjector";
import { ProjectSettingsProvider } from "@/lib/project-settings";

export function LoginFlowLayout() {
  return (
    <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="max-w-sm w-full mx-auto">
        <ProjectSettingsProvider>
          <UISettingsInjector>
            <Outlet />
          </UISettingsInjector>
        </ProjectSettingsProvider>
      </div>
    </div>
  );
}
