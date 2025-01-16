import React from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn } from '@/lib/utils'
import { useQuery } from '@connectrpc/connect-query'
import { getProjectUISettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'

const Page = () => {
  const isDarkMode = useDarkMode()
  const { data: uiSettingsRes } = useQuery(getProjectUISettings)

  const shouldDetectDarkMode = () => {
    return uiSettingsRes?.projectUiSettings
      ? uiSettingsRes.projectUiSettings.detectDarkModeEnabled
      : true
  }

  return (
    <div
      className={cn(
        'mx-auto flex flex-col justify-center items-center min-h-screen w-screen py-8',
        isDarkMode && shouldDetectDarkMode()
          ? 'dark bg-dark'
          : 'light bg-muted',
      )}
    >
      <div className="container flex justify-center">
        <div className="mb-8">
          {/* TODO: Make this conditionally load an Organizations configured logo */}
          {isDarkMode && shouldDetectDarkMode() ? (
            <img
              className="max-w-[240px]"
              src="/images/tesseral-logo-white.svg"
            />
          ) : (
            <img
              className="max-w-[240px]"
              src="/images/tesseral-logo-black.svg"
            />
          )}
        </div>
      </div>
      <Outlet />
    </div>
  )
}

export default Page
