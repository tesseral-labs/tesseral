import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { PageLoading } from "@/components/page/PageLoading";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Title } from "@/components/page/Title";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { getUser } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

export function UserPage() {
  const { pathname } = useLocation();
  const { organizationId, userId } = useParams();
  const {
    data: getUserResponse,
    isError,
    isLoading,
  } = useQuery(
    getUser,
    {
      id: userId,
    },
    {
      retry: false,
    },
  );

  return (
    <>
      {isLoading ? (
        <PageLoading />
      ) : isError ? (
        <NotFound />
      ) : (
        <PageContent>
          <Title title={getUserResponse?.user?.email || "User"} />

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
              {getUserResponse?.user?.displayName ||
                getUserResponse?.user?.email}
            </h1>
            <ValueCopier
              value={getUserResponse?.user?.id || ""}
              label="User ID"
            />
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
            <TabLink
              active={
                pathname === `/organizations/${organizationId}/users/${userId}`
              }
              to={`/organizations/${organizationId}/users/${userId}`}
            >
              Details
            </TabLink>
            <TabLink
              active={
                pathname ===
                `/organizations/${organizationId}/users/${userId}/sessions`
              }
              to={`/organizations/${organizationId}/users/${userId}/sessions`}
            >
              Sessions
            </TabLink>
            <TabLink
              active={
                pathname ===
                `/organizations/${organizationId}/users/${userId}/roles`
              }
              to={`/organizations/${organizationId}/users/${userId}/roles`}
            >
              Roles
            </TabLink>
            <TabLink
              active={
                pathname ===
                `/organizations/${organizationId}/users/${userId}/passkeys`
              }
              to={`/organizations/${organizationId}/users/${userId}/passkeys`}
            >
              Passkeys
            </TabLink>
            <TabLink
              active={
                pathname ===
                `/organizations/${organizationId}/users/${userId}/history`
              }
              to={`/organizations/${organizationId}/users/${userId}/history`}
            >
              User History
            </TabLink>
            <TabLink
              active={
                pathname ===
                `/organizations/${organizationId}/users/${userId}/activity`
              }
              to={`/organizations/${organizationId}/users/${userId}/activity`}
            >
              User Activity
            </TabLink>
          </Tabs>

          <div>
            <Outlet />
          </div>
        </PageContent>
      )}
    </>
  );
}
