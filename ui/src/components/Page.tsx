import React, { SyntheticEvent, useEffect, useRef, useState } from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn } from '@/lib/utils'
import { useQuery } from '@connectrpc/connect-query'
import { getProjectUISettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import useProjectUiSettings from '@/lib/project-ui-settings'
import { Helmet } from 'react-helmet'

const Page = () => {
  const isDarkMode = useDarkMode()
  const projectUiSettings = useProjectUiSettings()

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
      </Helmet>

      <div
        className={cn(
          'mx-auto flex flex-col justify-center items-center min-h-screen w-screen py-8',
          isDarkMode && projectUiSettings.detectDarkModeEnabled
            ? 'dark bg-dark'
            : 'light bg-muted',
        )}
      >
        <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
          <div className="mb-8">
            {/* TODO: Make this conditionally load an Organizations configured logo */}
            {isDarkMode && projectUiSettings.detectDarkModeEnabled ? (
              <img
                className="max-w-[240px]"
                src={
                  projectUiSettings?.darkModeLogoUrl ||
                  '/images/tesseral-logo-white.svg'
                }
                onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                  const target = e.target as HTMLImageElement
                  target.onerror = null
                  target.src = '/images/tesseral-logo-white.svg'
                }}
              />
            ) : (
              <img
                className="max-w-[240px]"
                src={
                  projectUiSettings?.logoUrl ||
                  '/images/tesseral-logo-black.svg'
                }
                onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
                  const target = e.target as HTMLImageElement
                  target.onerror = null
                  target.src = '/images/tesseral-logo-black.svg'
                }}
              />
            )}
          </div>
        </div>
        <Outlet />
      </div>
    </>
  )
}

export default Page
