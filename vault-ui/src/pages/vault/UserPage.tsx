import { useQuery } from "@connectrpc/connect-query";
import { ChevronDown } from "lucide-react";
import React from "react";
import { Helmet } from "react-helmet";
import { Link, Outlet, useLocation } from "react-router";

import { PageContent } from "@/components/page";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getOrganization } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function UserPage() {
  return (
    <PageContent>
      <Helmet>
        <title>User Settings</title>
      </Helmet>
      <div>
        <h1 className="text-2xl font-semibold">Account settings</h1>
        <p className="text-muted-foreground">
          Manage your account settings, authentication methods, and view your
          history.
        </p>
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

  const { data: getOrganizationResponse } = useQuery(getOrganization);

  const organization = getOrganizationResponse?.organization;

  return (
    <>
      <Tabs className="hidden lg:inline-flex">
        <TabLink active={pathname === `/user`} to={`/user`}>
          Details
        </TabLink>
        {(organization?.logInWithAuthenticatorApp ||
          organization?.logInWithPasskey) && (
          <TabLink
            active={pathname === `/user/authentication`}
            to={`/user/authentication`}
          >
            Authentication
          </TabLink>
        )}
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
                {pathname === `user` && "Details"}
                {(organization?.logInWithAuthenticatorApp ||
                  organization?.logInWithPasskey) && (
                  <>{pathname === `user/authentication` && "Authentication"}</>
                )}
              </span>
              <ChevronDown className="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem asChild>
              <Link to={`/user`}>Details</Link>
            </DropdownMenuItem>
            {(organization?.logInWithAuthenticatorApp ||
              organization?.logInWithPasskey) && (
              <DropdownMenuItem asChild>
                <Link to={`/user/authentication`}>Authentication</Link>
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </>
  );
}
