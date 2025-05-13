import React from 'react';
import { Outlet } from 'react-router';

import { ProjectSettingsProvider } from '@/lib/project-settings';

export function LoginFlowLayout() {
  return (
    <div className="bg-zinc-950 w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8 relative">
      <div className="absolute flex justify-center items-center blur-3xl w-full z-5">
        <div className="relative rounded-full w-[750px] h-[750px] bg-indigo-600/30 blur-3xl m-auto" />
      </div>
      <div className="max-w-sm w-full mx-auto z-10">
        <ProjectSettingsProvider>
          <Outlet />
        </ProjectSettingsProvider>
      </div>
    </div>
  );
}
