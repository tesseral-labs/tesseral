import React from "react";

import { PageContent } from "@/components/page";

import { BasicSettingsCard } from "./BasicSettingsCard";
import { EnterpriseSettingsCard } from "./EnterpriseSettingsCard";
import { ManagedApiKeySettingsCard } from "./ManagedApiKeySettingsCard";
import { MfaSettingsCard } from "./MfaSettingsCard";
import { OAuthSettingsCard } from "./OauthSettingsCard";

export function AuthenticationSettingsPage() {
  return (
    <PageContent>
      <div>
        <h1 className="text-xl font-semibold">Authentication Settings</h1>
        <p className="text-muted-foreground text-sm">
          Configure global authentication methods and identity providers.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <BasicSettingsCard />
        <OAuthSettingsCard />
        <MfaSettingsCard />
        <EnterpriseSettingsCard />
        <ManagedApiKeySettingsCard />
      </div>
    </PageContent>
  );
}
