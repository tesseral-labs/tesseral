import React, { FC, SyntheticEvent } from 'react'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from './ui/sidebar'
import { Building2, FolderGit, KeyRound } from 'lucide-react'

const ConsoleSidebar: FC = () => {
  const overviewItems = [
    {
      icon: Building2,
      title: 'Organizations',
      url: '/organizations',
    },
  ]
  const projectItems = [
    {
      icon: FolderGit,
      title: 'Project Settings',
      url: '/project-settings',
    },
    {
      icon: KeyRound,
      title: 'Project API Keys',
      url: '/project-api-keys',
    },
  ]

  return (
    <Sidebar className="min-h-screen" collapsible="none" variant="inset">
      <SidebarHeader>
        <div className="px-2 pt-4">
          <img className="max-h-[24px]" src="/images/tesseral-logo-white.svg" />
        </div>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup title="Overview">
          <SidebarGroupLabel>Overview</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {overviewItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup title="Project">
          <SidebarGroupLabel>Project</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {projectItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}

export default ConsoleSidebar
