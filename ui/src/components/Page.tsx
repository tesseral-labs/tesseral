import React, { SyntheticEvent, useEffect, useRef, useState } from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn } from '@/lib/utils'
import useSettings from '@/lib/settings'
import { Helmet } from 'react-helmet'

const Page = () => {
  const isDarkMode = useDarkMode()
  const settings = useSettings()

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
      </Helmet>

      <div
        className={cn(
          'min-h-screen w-screen',
          isDarkMode && settings.detectDarkModeEnabled
            ? 'dark bg-dark'
            : 'light bg-body',
        )}
      >
        <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
          <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
            <div className="mb-8">
              {/* TODO: Make this conditionally load an Organizations configured logo */}
              {isDarkMode && settings.detectDarkModeEnabled ? (
                <img
                  className="max-w-[180px]"
                  src={
                    settings?.darkModeLogoUrl ||
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
    </>
  )
}

export default Page
