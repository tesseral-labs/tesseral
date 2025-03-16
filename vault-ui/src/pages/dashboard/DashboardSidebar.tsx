import { useQuery } from "@connectrpc/connect-query";
import {
  Building2Icon,
  ChevronsUpDownIcon,
  LayoutGridIcon,
  UserIcon,
} from "lucide-react";
import React from "react";
import { Link, useNavigate } from "react-router-dom";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import {
  getOrganization,
  listSwitchableOrganizations,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function DashboardSidebar() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
  );
  const { data: whoamiResponse } = useQuery(whoami);

  const navigate = useNavigate();

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                    <LayoutGridIcon className="size-4" />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      {getOrganizationResponse?.organization?.displayName}
                    </span>
                    <span className="truncate text-xs">
                      {getOrganizationResponse?.organization?.id}
                    </span>
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
                {listSwitchableOrganizationsResponse?.switchableOrganizations
                  ?.filter(
                    (org) =>
                      org.id !== getOrganizationResponse?.organization?.id,
                  )
                  .map((org) => (
                    <DropdownMenuItem
                      key={org.id}
                      className="gap-2 p-2"
                      onClick={() => {
                        navigate(`/switch-organizations/${org.id}`);
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
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/organization-settings">
                    <Building2Icon />
                    Organization Settings
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>

            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/user-settings">
                    <UserIcon />
                    User Settings
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarFallback className="rounded-lg">
                  {whoamiResponse?.user?.email?.substring(0, 1)?.toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">
                  {whoamiResponse?.user?.email}
                </span>
                <span className="truncate text-xs">
                  {whoamiResponse?.user?.email}
                </span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
