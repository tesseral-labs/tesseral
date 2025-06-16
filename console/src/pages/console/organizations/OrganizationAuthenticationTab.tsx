import React from "react";

import { ListOrganizationSamlConnectionsCard } from "./authentication/ListOrganizationSamlConnectionsCard";
import { ListOrganizationScimApiKeysCard } from "./authentication/ListOrganizationScimApiKeysCard";
import { OrganizationBasicAuthCard } from "./authentication/OrganizationBasicAuthCard";
import { OrganizationMFACard } from "./authentication/OrganizationMfaCard";
import { OrganizationOAuthCard } from "./authentication/OrganizationOauthCard";
import { OrganizationSamlCard } from "./authentication/OrganizationSamlCard";
import { OrganizationScimCard } from "./authentication/OrganizationScimCard";

export function OrganizationAuthentication() {
  return (
    <div className="space-y-8">
      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <OrganizationBasicAuthCard />
        <OrganizationOAuthCard />
        <OrganizationMFACard />
        <OrganizationSamlCard />
        <OrganizationScimCard />
      </div>

      <ListOrganizationSamlConnectionsCard />
      <ListOrganizationScimApiKeysCard />
    </div>
  );
}
