import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { useParams } from "react-router";

import { getOrganization } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { ListOrganizationRolesCard } from "./roles/ListOrganizationRolesCard";
import { OrganizationRoleSettingsCard } from "./roles/OrganizationRoleSettingsCard";

export function OrganizationRolesTab() {
  const { organizationId } = useParams();

  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });

  const organization = getOrganizationResponse?.organization;

  return (
    <div className="space-y-8">
      <div className="grid grid-cols-3">
        <OrganizationRoleSettingsCard />
      </div>

      {organization?.customRolesEnabled && <ListOrganizationRolesCard />}
    </div>
  );
}
