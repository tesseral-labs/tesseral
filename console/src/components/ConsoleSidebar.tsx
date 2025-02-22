import React, { FC } from 'react';
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
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
  useSidebar,
} from './ui/sidebar';
import { Building2, Settings2Icon, UserIcon } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useQuery } from '@connectrpc/connect-query';
import { whoami } from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';

const ConsoleSidebar: FC = () => {
  const { data: whoamiResponse } = useQuery(whoami, {});
  const { state } = useSidebar();

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <Link className="flex flex-col" to="/">
          {state === 'expanded' ? (
            <img
              className="justify-self-start py-2 h-10"
              src="/images/tesseral-logo-black.svg"
            />
          ) : (
            <img
              className="justify-self-start py-2 h-10"
              src="/images/tesseral-icon-black.svg"
            />
          )}
        </Link>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/organizations">
                    <Building2 />
                    Organizations
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>

            <SidebarMenu>
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
                      <Link to="/project-settings">General</Link>
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
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton>
              <UserIcon />
              {whoamiResponse?.user?.email}
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
};

export default ConsoleSidebar;
