import React, { FC, PropsWithChildren, useEffect, useState } from 'react'
import { Helmet } from 'react-helmet'
import { cn } from '@/lib/utils'

import useDarkMode from '@/lib/dark-mode'
import useSettings from '@/lib/settings'
import { Outlet } from 'react-router'
import {
  OrganizationContextProvider,
  ProjectContextProvider,
  UserContextProvider,
  useSession,
} from '@/lib/auth'
import Header from './Header'

const DashboardPage: FC<PropsWithChildren> = ({ children }) => {
  const isDarkMode = useDarkMode()
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
    <>
      <Helmet>
        <link rel="icon" href={favicon} />
        <link rel="apple-touch-icon" href={favicon} />
        <title>{session?.organization?.displayName || 'Dashboard'}</title>
      </Helmet>

      <ProjectContextProvider value={session?.project}>
        <OrganizationContextProvider value={session?.organization}>
          <UserContextProvider value={session?.user}>
            <div
              className={cn(
                'min-h-screen w-screen',
                isDarkMode && settings.detectDarkModeEnabled ? 'dark' : '',
              )}
            >
              <div className="bg-body w-screen min-h-screen mx-auto items-center">
                <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
                  <Header />
                  <div className="py-8">{children}</div>
                </div>
              </div>
            </div>
          </UserContextProvider>
        </OrganizationContextProvider>
      </ProjectContextProvider>
    </>
  )
}

export default DashboardPage
