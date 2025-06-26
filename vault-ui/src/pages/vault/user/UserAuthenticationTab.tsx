import { useQuery } from "@connectrpc/connect-query";
import React from "react";

import { TabContent } from "@/components/page/Tabs";
import { getOrganization } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { UserAuthenticatorAppCard } from "./authentication/UserAuthenticatorAppCard";
import { UserPasskeysCard } from "./authentication/UserPasskeysCard";

export function UserAuthenticationTab() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);

  const organization = getOrganizationResponse?.organization;

  return (
    <TabContent>
      {organization?.logInWithAuthenticatorApp && <UserAuthenticatorAppCard />}
      {organization?.logInWithPasskey && <UserPasskeysCard />}
    </TabContent>
  );
}
