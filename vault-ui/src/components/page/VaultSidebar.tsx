import { useQuery } from "@connectrpc/connect-query";
import {
  AlignLeft,
  ChevronsUpDown,
  ChevronsUpDownIcon,
  Key,
  LayoutGridIcon,
  Lock,
  LogOut,
  Logs,
  UserPlus,
  Users,
} from "lucide-react";
import React from "react";
import { Link, useLocation, useNavigate } from "react-router";

import {
  getOrganization,
  getProject,
  listSwitchableOrganizations,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { cn } from "@/lib/utils";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarSeparator,
} from "../ui/sidebar";

export function VaultSidebar() {
  const { pathname } = useLocation();

  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: whoamiResponse } = useQuery(whoami);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;
  const user = whoamiResponse?.user;

  return (
    <Sidebar variant="inset">
      <SidebarContent className="overflow-x-hidden">
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <OrganizationSwitcher />
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {user?.owner && (
          <SidebarGroup>
            <SidebarGroupLabel>Organization</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      className={cn(
                        pathname === "/organization"
                          ? "bg-muted text-primary font-medium"
                          : "",
                      )}
                      to="/organization"
                    >
                      <AlignLeft />
                      Details
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      className={cn(
                        pathname === "/organization/authentication"
                          ? "bg-muted text-primary font-medium"
                          : "",
                      )}
                      to="/organization/authentication"
                    >
                      <Lock />
                      Authentication
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      className={cn(
                        pathname.startsWith("/organization/users")
                          ? "bg-muted text-primary font-medium"
                          : "",
                      )}
                      to="/organization/users"
                    >
                      <Users />
                      Users
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      className={cn(
                        pathname.startsWith("/organization/user-invites")
                          ? "bg-muted text-primary font-medium"
                          : "",
                      )}
                      to="/organization/user-invites"
                    >
                      <UserPlus />
                      User Invites
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                {project?.apiKeysEnabled && organization?.apiKeysEnabled && (
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link
                        className={cn(
                          pathname.startsWith("/organization/api-keys")
                            ? "bg-muted text-primary font-medium"
                            : "",
                        )}
                        to="/organization/api-keys"
                      >
                        <Key />
                        API Keys
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                )}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}
        <SidebarGroup>
          <SidebarGroupLabel>User</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link
                    className={cn(
                      pathname === "/user"
                        ? "bg-muted text-primary font-medium"
                        : "",
                    )}
                    to="/user"
                  >
                    <AlignLeft />
                    Details
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              {(organization?.logInWithAuthenticatorApp ||
                organization?.logInWithPasskey) && (
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      className={cn(
                        pathname === "/user/authentication"
                          ? "bg-muted text-primary font-medium"
                          : "",
                      )}
                      to="/user/authentication"
                    >
                      <Lock />
                      Authentication
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {user?.owner && project?.auditLogsEnabled && (
          <>
            <SidebarSeparator />
            <SidebarGroup>
              <SidebarGroupLabel>System</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link
                        className={cn(
                          pathname === "/logs"
                            ? "bg-muted text-primary font-medium"
                            : "",
                        )}
                        to="/logs"
                      >
                        <Logs />
                        Audit Logs
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </>
        )}
      </SidebarContent>
      <SidebarFooter>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <UserMenu />
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarFooter>
    </Sidebar>
  );
}

function OrganizationSwitcher() {
  const navigate = useNavigate();

  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
  );

  const organization = getOrganizationResponse?.organization;
  const organizations =
    listSwitchableOrganizationsResponse?.switchableOrganizations || [];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <SidebarMenuButton className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
          <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
            <LayoutGridIcon className="size-4" />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold">
              {organization?.displayName}
            </span>
            <span className="truncate text-xs">{organization?.id}</span>
          </div>
          <ChevronsUpDownIcon className="ml-auto" />
        </SidebarMenuButton>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
        align="start"
        side="right"
        sideOffset={4}
      >
        <DropdownMenuLabel className="text-xs text-muted-foreground">
          Organizations
        </DropdownMenuLabel>
        {organizations?.map((org) => (
          <DropdownMenuItem
            key={org.id}
            className="gap-2 p-2"
            onClick={() => {
              if (org.id !== getOrganizationResponse?.organization?.id) {
                navigate(`/switch-organizations/${org.id}`);
              }
            }}
          >
            <div className="flex size-6 items-center justify-center rounded-sm border">
              <LayoutGridIcon className="size-4 shrink-0" />
            </div>
            {org.displayName}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function UserMenu() {
  const { data: whoamiResponse } = useQuery(whoami);

  const user = whoamiResponse?.user;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className="flex items-center justify-start gap-2 w-full">
        {user?.profilePictureUrl ? (
          <img
            src={user.profilePictureUrl}
            alt="User Avatar"
            className="w-8 h-8 rounded-full"
          />
        ) : (
          <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
            <span className="text-gray-600">
              {(user?.displayName || user?.email)?.charAt(0).toUpperCase() ||
                "U"}
            </span>
          </div>
        )}
        <div className="flex-shrink overflow-x-hidden">
          {user?.displayName && (
            <div className="text-xs text-muted-foreground text-left">
              {user.displayName}
            </div>
          )}
          <div className="text-sm font-medium">{user?.email}</div>
        </div>
        <ChevronsUpDown className="ml-auto h-4" />
      </DropdownMenuTrigger>
      <DropdownMenuContent side="right">
        <DropdownMenuItem>
          {user?.profilePictureUrl ? (
            <img
              src={user.profilePictureUrl}
              alt="User Avatar"
              className="w-8 h-8 rounded-full"
            />
          ) : (
            <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
              <span className="text-gray-600">
                {(user?.displayName || user?.email)?.charAt(0).toUpperCase() ||
                  "U"}
              </span>
            </div>
          )}
          <div className="flex-shrink overflow-x-hidden">
            {user?.displayName && (
              <div className="text-xs text-muted-foreground">
                {user.displayName}
              </div>
            )}
            <div className="text-sm font-medium">{user?.email}</div>
          </div>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Link to="/logout">
            <LogOut />
            Logout
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
