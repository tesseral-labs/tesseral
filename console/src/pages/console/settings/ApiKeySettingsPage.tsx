import React from "react";

import { PageContent } from "@/components/page";
import { Title } from "@/components/page/Title";

import { ListBackendApiKeysCard } from "./api-keys/backend-api-keys/ListBackendApiKeysCard";
import { ListPublishableKeysCard } from "./api-keys/publishable-keys/ListPublishableKeysCard";

export function ApiKeySettingsPage() {
  return (
    <PageContent>
      <Title title="API Key Settings" />

      <div>
        <h1 className="text-xl font-bold">API Keys</h1>
        <p className="text-muted-foreground text-sm">
          Manage Publishable and Backend API Keys for your Project.
        </p>
      </div>
      <div className="space-y-8">
        <ListPublishableKeysCard />
        <ListBackendApiKeysCard />
      </div>
    </PageContent>
  );
}
