import React from "react";
import { Outlet } from "react-router";

import {
  AccessTokenProvider,
  useAccessToken,
} from "@/lib/access-token-provider";
import { cn } from "@/lib/utils";

import { Navigation } from "./Navigation";

export function PageShell() {
  return (
    <AccessTokenProvider>
      <PageShellInner />
    </AccessTokenProvider>
  );
}

function PageShellInner() {
  const accessToken = useAccessToken();
  if (!accessToken) {
    return null;
  }

  return (
    <>
      <main className="w-full min-h-screen">
        <Navigation />

        <div>
          <Outlet />
        </div>
      </main>
    </>
  );
}

export function PageContent({
  className,
  children,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("container p-4 pb-16 m-auto space-y-8", className)}
      {...props}
    >
      {children}
    </div>
  );
}
PageContent.displayName = "PageContent";
