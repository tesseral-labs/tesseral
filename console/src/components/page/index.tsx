import React from "react";
import { Outlet } from "react-router";

import {
  AccessTokenProvider,
  useAccessToken,
} from "@/lib/access-token-provider";
import { GlobalSearchProvider } from "@/lib/search";
import { cn } from "@/lib/utils";

import { Search } from "../core/Search";
import { Navigation } from "./Navigation";

export function PageShell() {
  return (
    <AccessTokenProvider>
      <GlobalSearchProvider>
        <PageShellInner />
      </GlobalSearchProvider>
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

      <Search />
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
