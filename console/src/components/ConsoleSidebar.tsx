import React, { FC, useCallback } from 'react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent, SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
} from './ui/sidebar';
import {
  BadgeCheckIcon, BookIcon, BookOpenIcon, BugIcon,
  Building2Icon,
  ChevronsUpDownIcon, HomeIcon,
  LayoutGridIcon, LifeBuoyIcon,
  LogOutIcon,
  PlusIcon,
  Settings2Icon,
  UserIcon,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  listSwitchableOrganizations,
  logout,
  whoami,
} from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { getProject } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useNavigate } from 'react-router';
import { useIsMobile } from '@/hooks/use-mobile';
import { API_URL } from '@/config';
import { toast } from 'sonner';

const ConsoleSidebar: FC = () => {
  const { data: whoamiResponse } = useQuery(whoami, {});
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
    {},
  );

  const isMobile = useIsMobile();
  const navigate = useNavigate();

  const { mutateAsync: logoutAsync } = useMutation(logout);
  const handleLogout = async () => {
    await logoutAsync({});
    toast.success('You have been logged out.');
    navigate('/login');
  };

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
                      {getProjectResponse?.project?.displayName}
                    </span>
                    <span className="truncate text-xs">
                      {getProjectResponse?.project?.vaultDomain}
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
                  Projects
                </DropdownMenuLabel>
                {listSwitchableOrganizationsResponse?.switchableOrganizations?.map(
                  (org) => (
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
                  ),
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem className="gap-2 p-2">
                  <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                    <PlusIcon className="size-4" />
                  </div>
                  <div className="font-medium text-muted-foreground">
                    Create new project
                  </div>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Project</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/">
                    <HomeIcon />
                    Project Home
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/organizations">
                    <Building2Icon />
                    Organizations
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/project-settings">
                    <Settings2Icon />
                    Project Settings
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuSub>
                  <SidebarMenuSubItem>
                    <SidebarMenuSubButton asChild>
                      <Link to="/project-settings">General Settings</Link>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                  <SidebarMenuSubItem>
                    <SidebarMenuSubButton asChild>
                      <Link to="/project-settings/api-keys">
                        Project API Keys
                      </Link>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                </SidebarMenuSub>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarGroup>
          <SidebarGroupLabel>Resources</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="https://tesseral.com/docs/quickstart" target="_blank">
                    <BookOpenIcon />
                    Tesseral Documentation
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuButton asChild>
                  <Link to="https://github.com/tesseral-labs/tesseral/issues/new" target="_blank">
                    <BugIcon />
                    Report an Issue
                  </Link>
                </SidebarMenuButton>
                <SidebarMenuButton asChild>
                  <Link to="mailto:support@tesseral.com" target="_blank">
                    <LifeBuoyIcon />
                    Contact Support
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
          <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <Avatar className="h-8 w-8 rounded-lg">
                    <AvatarFallback className="rounded-lg">
                      {whoamiResponse?.user?.email
                        ?.substring(0, 1)
                        ?.toUpperCase()}
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
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                side={isMobile ? 'bottom' : 'right'}
                align="end"
                sideOffset={4}
              >
                <DropdownMenuLabel className="p-0 font-normal">
                  <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <Avatar className="h-8 w-8 rounded-lg">
                      <AvatarFallback className="rounded-lg">
                        {whoamiResponse?.user?.email
                          ?.substring(0, 1)
                          ?.toUpperCase()}
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
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                  <DropdownMenuItem asChild>
                    <Link to={`${API_URL}/user-settings`}>
                      <UserIcon />
                      User Settings
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link to={`${API_URL}/organization-settings`}>
                      <Building2Icon />
                      Collaboration Settings
                    </Link>
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOutIcon />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
};

export default ConsoleSidebar;
