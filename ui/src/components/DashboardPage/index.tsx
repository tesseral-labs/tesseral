import React, { FC, PropsWithChildren, useEffect, useState } from 'react'
import { Helmet } from 'react-helmet'
import { cn } from '@/lib/utils'

import useDarkMode from '@/lib/dark-mode'
import useSettings from '@/lib/settings'
import {
  OrganizationContextProvider,
  ProjectContextProvider,
  UserContextProvider,
  useSession,
} from '@/lib/auth'
import Header from './Header'
import { SidebarProvider, SidebarTrigger } from '../ui/sidebar'
import DashboardSidebar from './DashboardSidebar'
import { Toaster } from '../ui/sonner'
import { useIsMobile } from '@/hooks/use-mobile'

const DashboardPage: FC<PropsWithChildren> = ({ children }) => {
  const isDarkMode = useDarkMode()
  const isMobile = useIsMobile()
  const settings = useSettings()
  const session = useSession()

  const [favicon, setFavicon] = useState<string>('/apple-touch-icon.png')

  useEffect(() => {
    if (settings?.faviconUrl) {
      ;(async () => {
        // Check if the favicon exists before setting it
        const faviconCheck = await fetch(settings.faviconUrl, {
          method: 'HEAD',
        })

        setFavicon(
          faviconCheck.ok ? settings.faviconUrl : '/apple-touch-icon.png',
        )
      })()
    }
  }, [settings])

  return (
    <div
      className={isDarkMode && settings?.detectDarkModeEnabled ? 'dark' : ''}
    >
      <Helmet>
        <link rel="icon" href={favicon} />
        <link rel="apple-touch-icon" href={favicon} />
        <title>{session?.organization?.displayName || 'Dashboard'}</title>
      </Helmet>

      <ProjectContextProvider value={session?.project}>
        <OrganizationContextProvider value={session?.organization}>
          <UserContextProvider value={session?.user}>
            <SidebarProvider>
              <DashboardSidebar />
              <main className="min-h-screen w-screen">
                {isMobile && <SidebarTrigger />}
                <div className="bg-background min-h-screen mx-auto items-center">
                  <div className="mx-auto px-6 lg:px-8">
                    <div className="py-8">{children}</div>
                  </div>
                </div>
              </main>
              <Toaster position="top-center" />
            </SidebarProvider>
          </UserContextProvider>
        </OrganizationContextProvider>
      </ProjectContextProvider>
    </div>
  )
}

export default DashboardPage
