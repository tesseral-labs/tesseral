import React, { PropsWithChildren } from "react";
import { Outlet } from "react-router";

import { SidebarProvider } from "../ui/sidebar";
import { VaultSidebar } from "./VaultSidebar";

export function Page() {
  return (
    <SidebarProvider>
      <VaultSidebar />
      <div className="container p-4 lg:p-8 w-full">
        <Outlet />
      </div>
    </SidebarProvider>
  );
}

export function PageContent({ children }: PropsWithChildren) {
  return <div className="space-y-4 lg:space-y-8 w-full">{children}</div>;
}
