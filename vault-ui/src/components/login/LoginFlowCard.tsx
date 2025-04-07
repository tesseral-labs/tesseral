import React from "react";

import { Card } from "@/components/ui/card";
import { useProjectSettings } from "@/lib/project-settings";
import { useDarkMode } from "@/lib/dark-mode";

export function LoginFlowCard({ children }: { children?: React.ReactNode }) {
  const settings = useProjectSettings();
  const isDarkMode = useDarkMode();

  const logo = isDarkMode ? settings?.darkModeLogoUrl : settings?.logoUrl

  return (
    <>
      <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
        {logo && (
          <img alt="logo" src={logo} className="mb-8 max-w-[180px] max-h-[80px]" />
        )}
      </div>
      <Card className="w-full">{children}</Card>
    </>
  );
}
