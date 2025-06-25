import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft, ChevronDown } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { Title } from "@/components/core/Title";
import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getUser } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function OrganizationUserPage() {
  const { userId } = useParams();

  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });

  const user = getUserResponse?.user;

  return (
    <PageContent>
      <Title title={user?.displayName || user?.email || "User"} />

      <div>
        <Link to={`/organization/users`}>
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to Users
          </Button>
        </Link>
      </div>

      <div className="">
        <h1 className="text-2xl font-semibold">
          {user?.displayName || user?.email}
        </h1>
        <ValueCopier value={user?.id || ""} label="User ID" />
        <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
          <Badge className="border-0" variant="outline">
            Created{" "}
            {user?.createTime &&
              DateTime.fromJSDate(timestampDate(user.createTime)).toRelative()}
          </Badge>
          <div>•</div>
          <Badge className="border-0" variant="outline">
            Updated{" "}
            {user?.updateTime &&
              DateTime.fromJSDate(timestampDate(user.updateTime)).toRelative()}
          </Badge>
        </div>
      </div>

      <UserTabs />

      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}

function UserTabs() {
  const { pathname } = useLocation();
  const { userId } = useParams();

  return (
    <>
      <Tabs className="hidden lg:inline-flex">
        <TabLink
          active={pathname === `/organization/users/${userId}`}
          to={`/organization/users/${userId}`}
        >
          Details
        </TabLink>
        <TabLink
          active={pathname === `/organization/users/${userId}/roles`}
          to={`/organization/users/${userId}/roles`}
        >
          Roles
        </TabLink>
      </Tabs>
      <div className="lg:hidden">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              className="flex items-center gap-2"
              variant="outline"
              size="sm"
            >
              <span>
                {pathname === `/organization/users/${userId}` && "Details"}
                {pathname === `/organization/users/${userId}/roles` && "Roles"}
              </span>
              <ChevronDown className="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem asChild>
              <Link to={`/organization/users/${userId}`}>Details</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to={`/organization/users/${userId}/roles`}>Roles</Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </>
  );
}
