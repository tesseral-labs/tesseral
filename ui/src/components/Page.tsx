import React, { SyntheticEvent, useRef } from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn } from '@/lib/utils'
import { useQuery } from '@connectrpc/connect-query'
import { getProjectUISettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import useProjectUiSettings from '@/lib/project-ui-settings'

const Page = () => {
  const isDarkMode = useDarkMode()
  const projectUiSettings = useProjectUiSettings()

  return (
    <div
      className={cn(
        'mx-auto flex flex-col justify-center items-center min-h-screen w-screen py-8',
        isDarkMode && projectUiSettings.detectDarkModeEnabled
          ? 'dark bg-dark'
          : 'light bg-muted',
      )}
    >
      <div className="container flex justify-center">
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
                projectUiSettings?.logoUrl || '/images/tesseral-logo-black.svg'
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
  )
}

export default Page
