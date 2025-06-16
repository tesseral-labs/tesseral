import React from "react";

import { PageContent } from "@/components/page";

import { ListBackendApiKeysCard } from "./backend-api-keys/ListBackendApiKeysCard";

export function BackendApiKeysPage() {
  return (
    <PageContent>
      <div className="">
        <h1 className="text-2xl font-semibold">Backend API Keys</h1>
        <p className="text-muted-foreground text-sm">
          Manage your backend API keys for accessing Tesseral services.
        </p>
      </div>

      <ListBackendApiKeysCard />
    </PageContent>
  );
}
