import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { Tab, Tabs } from "@/components/page/Tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { getUser } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function UserPage() {
  const { pathname } = useLocation();
  const { organizationId, userId } = useParams();
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });

  return (
    <PageContent>
      <div>
        <Link to={`/organizations/${organizationId}/users`}>
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to Users
          </Button>
        </Link>
      </div>

      <div className="">
        <h1 className="text-2xl font-semibold">
          {getUserResponse?.user?.displayName || getUserResponse?.user?.email}
        </h1>
        <ValueCopier value={getUserResponse?.user?.id || ""} label="User ID" />
        <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
          <Badge className="border-0" variant="outline">
            Created{" "}
            {getUserResponse?.user?.createTime &&
              DateTime.fromJSDate(
                timestampDate(getUserResponse.user.createTime),
              ).toRelative()}
          </Badge>
          <div>•</div>
          <Badge className="border-0" variant="outline">
            Updated{" "}
            {getUserResponse?.user?.updateTime &&
              DateTime.fromJSDate(
                timestampDate(getUserResponse.user.updateTime),
              ).toRelative()}
          </Badge>
        </div>
      </div>

      <Tabs>
        <Tab
          active={
            pathname === `/organizations/${organizationId}/users/${userId}`
          }
        >
          <Link to={`/organizations/${organizationId}/users/${userId}`}>
            Details
          </Link>
        </Tab>
        <Tab
          active={
            pathname ===
            `/organizations/${organizationId}/users/${userId}/sessions`
          }
        >
          <Link
            to={`/organizations/${organizationId}/users/${userId}/sessions`}
          >
            Sessions
          </Link>
        </Tab>
        <Tab
          active={
            pathname ===
            `/organizations/${organizationId}/users/${userId}/roles`
          }
        >
          <Link to={`/organizations/${organizationId}/users/${userId}/roles`}>
            Roles
          </Link>
        </Tab>
        <Tab
          active={
            pathname ===
            `/organizations/${organizationId}/users/${userId}/passkeys`
          }
        >
          <Link
            to={`/organizations/${organizationId}/users/${userId}/passkeys`}
          >
            Passkeys
          </Link>
        </Tab>
        <Tab
          active={
            pathname === `/organizations/${organizationId}/users/${userId}/logs`
          }
        >
          <Link to={`/organizations/${organizationId}/users/${userId}/logs`}>
            Audit Logs
          </Link>
        </Tab>
      </Tabs>

      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}
