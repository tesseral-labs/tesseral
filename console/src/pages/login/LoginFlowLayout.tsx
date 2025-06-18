import React from "react";
import { Outlet } from "react-router";

import { ProjectSettingsProvider } from "@/lib/project-settings";

export function LoginFlowLayout() {
  return (
    <div className="w-full min-h-screen mx-auto flex flex-col justify-center items-center py-8 relative">
      <div className="max-w-sm w-full mx-auto z-10">
        <ProjectSettingsProvider>
          <Outlet />
        </ProjectSettingsProvider>
      </div>
    </div>
  );
}
