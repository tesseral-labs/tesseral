import React from "react";

import { PageContent } from "@/components/page";
import { Title } from "@/components/page/Title";

import { ApiKeysCard } from "./overview/APIKeysCard";
import { AuthenticationCard } from "./overview/AuthenticationCard";
import { VaultCustomizationCard } from "./overview/VaultCustomizationCard";

export function SettingsOverviewPage() {
  return (
    <PageContent>
      <Title title="Settings" />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <AuthenticationCard />
        <ApiKeysCard />
        <VaultCustomizationCard />
      </div>
    </PageContent>
  );
}
