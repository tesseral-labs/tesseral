import React, { FC, useEffect, useState } from 'react'
import { Helmet } from 'react-helmet'
import { cn } from '@/lib/utils'

import useDarkMode from '@/lib/dark-mode'
import useProjectUiSettings from '@/lib/project-ui-settings'
import { Outlet } from 'react-router'
import {
  OrganizationContextProvider,
  ProjectContextProvider,
  UserContextProvider,
  useSession,
} from '@/lib/auth'
import Header from './Header'

const DashboardPage: FC = () => {
  const isDarkMode = useDarkMode()
  const projectUiSettings = useProjectUiSettings()
  const session = useSession()

  const [favicon, setFavicon] = useState<string>('/apple-touch-icon.png')

  useEffect(() => {
    if (projectUiSettings?.faviconUrl) {
      ;(async () => {
        // Check if the favicon exists before setting it
        const faviconCheck = await fetch(projectUiSettings.faviconUrl, {
          method: 'HEAD',
        })

        setFavicon(
          faviconCheck.ok
            ? projectUiSettings.faviconUrl
            : '/apple-touch-icon.png',
        )
      })()
    }
  }, [projectUiSettings])

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
                isDarkMode && projectUiSettings.detectDarkModeEnabled
                  ? 'dark bg-dark'
                  : 'light bg-muted',
              )}
            >
              <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
                <Header />
                <Outlet />
              </div>
            </div>
          </UserContextProvider>
        </OrganizationContextProvider>
      </ProjectContextProvider>
    </>
  )
}

export default DashboardPage
