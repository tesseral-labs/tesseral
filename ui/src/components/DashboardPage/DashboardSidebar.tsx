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
} from '../ui/sidebar'
import useSettings from '@/lib/settings'
import useDarkMode from '@/lib/dark-mode'
import { Building2, UserCog } from 'lucide-react'

const DashboardSidebar: FC = () => {
  const settings = useSettings()
  const isDarkMode = useDarkMode()

  const items = [
    {
      icon: UserCog,
      title: 'User Settings',
      url: '/settings',
    },
    {
      icon: Building2,
      title: 'Organization Settings',
      url: '/organization',
    },
  ]

  return (
    <Sidebar className={isDarkMode ? 'dark' : ''} collapsible="icon">
      <SidebarHeader>
        <div>
          {isDarkMode && settings?.detectDarkModeEnabled ? (
            <img
              className="max-h-[20px] max-w-[150px]"
              src={
                settings?.darkModeLogoUrl || '/images/tesseral-logo-white.svg'
              }
              onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                const target = e.target as HTMLImageElement
                target.onerror = null
                target.src = '/images/tesseral-logo-white.svg'
              }}
            />
          ) : (
            <img
              className="max-h-[20px] max-w-[150px]"
              src={settings?.logoUrl || '/images/tesseral-logo-black.svg'}
              onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                const target = e.target as HTMLImageElement
                target.onerror = null
                target.src = '/images/tesseral-logo-black.svg'
              }}
            />
          )}
        </div>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup title="Settings">
          <SidebarGroupLabel>Settings</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
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

export default DashboardSidebar
