import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { getOrganization } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function OrganizationPage() {
  const { pathname } = useLocation();
  const { organizationId } = useParams();

  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });

  return (
    <PageContent>
      <div>
        <Link to="/organizations">
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to Organizations
          </Button>
        </Link>
      </div>

      <div className="">
        <h1 className="text-2xl font-semibold">
          {getOrganizationResponse?.organization?.displayName}
        </h1>
        <ValueCopier
          value={getOrganizationResponse?.organization?.id || ""}
          label="Organization ID"
        />
        <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
          <Badge className="border-0" variant="outline">
            Created{" "}
            {getOrganizationResponse?.organization?.createTime &&
              DateTime.fromJSDate(
                timestampDate(getOrganizationResponse.organization.createTime),
              ).toRelative()}
          </Badge>
          <div>•</div>
          <Badge className="border-0" variant="outline">
            Updated{" "}
            {getOrganizationResponse?.organization?.updateTime &&
              DateTime.fromJSDate(
                timestampDate(getOrganizationResponse.organization.updateTime),
              ).toRelative()}
          </Badge>
        </div>
      </div>
      <Tabs>
        <TabLink
          active={pathname === `/organizations/${organizationId}`}
          to={`/organizations/${organizationId}`}
        >
          Details
        </TabLink>
        <TabLink
          active={
            pathname === `/organizations/${organizationId}/authentication`
          }
          to={`/organizations/${organizationId}/authentication`}
        >
          Authentication
        </TabLink>
        <TabLink
          active={pathname.startsWith(
            `/organizations/${organizationId}/api-keys`,
          )}
          to={`/organizations/${organizationId}/api-keys`}
        >
          API Keys
        </TabLink>
        <TabLink
          active={pathname.startsWith(`/organizations/${organizationId}/users`)}
          to={`/organizations/${organizationId}/users`}
        >
          Users
        </TabLink>
        <TabLink
          active={pathname === `/organizations/${organizationId}/logs`}
          to={`/organizations/${organizationId}/logs`}
        >
          Audit Logs
        </TabLink>
      </Tabs>
      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}
