import React, { FC, SyntheticEvent } from 'react'

import useDarkMode from '@/lib/dark-mode'
import useSettings from '@/lib/settings'
import { Outlet } from 'react-router'

const SideBySideLayout: FC = () => {
  const isDarkMode = useDarkMode()
  const settings = useSettings()

  return (
    <div className="bg-body w-screen min-h-screen grid grid-cols-2 gap-0">
      <div className="bg-primary" />
      <div className="flex flex-col justify-center items-center p-4">
        <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
          <div className="mb-8">
            {/* TODO: Make this conditionally load an Organizations configured logo */}
            {isDarkMode && settings?.detectDarkModeEnabled ? (
              <img
                className="max-w-[180px]"
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
                className="max-w-[180px]"
                src={settings?.logoUrl || '/images/tesseral-logo-black.svg'}
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
    </div>
  )
}

export default SideBySideLayout
