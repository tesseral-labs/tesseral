import { useQuery } from "@connectrpc/connect-query";
import React from "react";

import { TabContent } from "@/components/page/Tabs";
import {
  getOrganization,
  getProject,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { OrganizationBasicAuthCard } from "./authentication/OrganizationBasicAuthCard";
import { OrganizationMfaCard } from "./authentication/OrganizationMfaCard";
import { OrganizationOAuthCard } from "./authentication/OrganizationOAuthCard";
import { SamlConnectionsCard } from "./saml-connections/SamlConnectionsCard";
import { ScimApiKeysCard } from "./scim-api-keys/ScimApiKeysCard";

export function OrganizationAuthenticationTab() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;

  return (
    <TabContent>
      <div className="grid grid-cols-1 xl:grid-cols-3 gap-4 lg:gap-8 w-full">
        {(project?.logInWithEmail || project?.logInWithPassword) && (
          <OrganizationBasicAuthCard />
        )}
        {(project?.logInWithGoogle ||
          project?.logInWithMicrosoft ||
          project?.logInWithGithub) && <OrganizationOAuthCard />}
        {(project?.logInWithAuthenticatorApp || project?.logInWithPasskey) && (
          <OrganizationMfaCard />
        )}
      </div>

      {organization?.logInWithSaml && <SamlConnectionsCard />}
      {organization?.scimEnabled && <ScimApiKeysCard />}
    </TabContent>
  );
}
