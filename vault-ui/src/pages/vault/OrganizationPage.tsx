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
import {
  getOrganization,
  getProject,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function OrganizationPage() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const organization = getOrganizationResponse?.organization;

  return (
    <PageContent>
      <Helmet>
        <title>{organization?.displayName || "Organization"} Settings</title>
      </Helmet>

      <div>
        <h1 className="text-2xl font-semibold">Organization</h1>
        <p className="text-muted-foreground">
          Manage your organization settings, users, and authentication methods.
        </p>
      </div>

      <OrganizationTabs />

      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}

function OrganizationTabs() {
  const { pathname } = useLocation();

  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;

  return (
    <>
      <Tabs className="hidden lg:inline-flex">
        <TabLink active={pathname === "/organization"} to="/organization">
          Details
        </TabLink>
        <TabLink
          active={pathname === "/organization/authentication"}
          to="/organization/authentication"
        >
          Authentication
        </TabLink>
        <TabLink
          active={pathname.startsWith("/organization/users")}
          to="/organization/users"
        >
          Users
        </TabLink>
        <TabLink
          active={pathname.startsWith("/organization/user-invites")}
          to="/organization/user-invites"
        >
          User Invites
        </TabLink>
        {organization?.apiKeysEnabled && (
          <TabLink
            active={pathname === "/organization/api-keys"}
            to="/organization/api-keys"
          >
            API Keys
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
                {pathname === `/organization` && "Details"}
                {pathname.startsWith(`/organization/users`) && "Users"}
                {pathname.startsWith(`/organization/user-invites`) &&
                  "User Invites"}
                {pathname === `/organization/authentication` &&
                  "Authentication"}
                {organization?.apiKeysEnabled &&
                  pathname.startsWith(`/organization/api-keys`) &&
                  "API Keys"}
              </span>
              <ChevronDown className="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem asChild>
              <Link to="/organization">Details</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to="/organization/authentication">Authentication</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to="/organization/users">Users</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to="/organization/user-invites">User Invites</Link>
            </DropdownMenuItem>
            {project?.apiKeysEnabled && organization?.apiKeysEnabled && (
              <DropdownMenuItem asChild>
                <Link to="/organization/api-keys">API Keys</Link>
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </>
  );
}
