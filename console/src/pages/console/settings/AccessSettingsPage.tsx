import React from "react";

import { PageContent } from "@/components/page";

import { ListProjectActionsCard } from "./access/ListProjectActionsCard";
import { ListProjectRolesCard } from "./access/ListProjectRolesCard";

export function AccessSettingsPage() {
  return (
    <PageContent>
      <div>
        <h1 className="text-xl font-bold">Access Control</h1>
        <p className="text-muted-foreground text-sm">
          Manage Role-based Access Control (RBAC) policies and roles.
        </p>
      </div>
      <ListProjectActionsCard />
      <ListProjectRolesCard />
    </PageContent>
  );
}
