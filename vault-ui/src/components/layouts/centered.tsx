import React, { FC, SyntheticEvent } from "react";
import { Outlet } from "react-router";

import useDarkMode from "@/lib/dark-mode";
import useSettings from "@/lib/settings";

const CenteredLayout: FC = () => {
  const isDarkMode = useDarkMode();
  const settings = useSettings();

  return (
    <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
      <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
        <div className="mb-8">
          {/* TODO: Make this conditionally load an Organizations configured logo */}
          {isDarkMode && settings?.detectDarkModeEnabled ? (
            <img
              className="max-w-[180px]"
              src={
                settings?.darkModeLogoUrl || "/images/tesseral-logo-white.svg"
              }
              onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                const target = e.target as HTMLImageElement;
                target.onerror = null;
                target.src = "/images/tesseral-logo-white.svg";
              }}
            />
          ) : (
            <img
              className="max-w-[180px]"
              src={settings?.logoUrl || "/images/tesseral-logo-black.svg"}
              onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                const target = e.target as HTMLImageElement;
                target.onerror = null;
                target.src = "/images/tesseral-logo-black.svg";
              }}
            />
          )}
        </div>
      </div>
      <Outlet />
    </div>
  );
};

export default CenteredLayout;
