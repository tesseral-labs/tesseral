import React from "react";
import { Outlet } from "react-router";

import { ProjectSettingsProvider } from "../../lib/project-settings";

interface LoginFlowLayoutProps {
  background?: string;
}

export function LoginFlowLayout({ background }: LoginFlowLayoutProps) {
  return (
    <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="max-w-sm w-sm mx-auto">
        <ProjectSettingsProvider>
          <Outlet />
        </ProjectSettingsProvider>
      </div>
    </div>
  );
}
