import React, { SyntheticEvent, useEffect, useRef, useState } from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn, hexToHSL, isColorDark } from '@/lib/utils'
import useSettings from '@/lib/settings'
import { Helmet } from 'react-helmet'
import { Settings } from '@/gen/openauth/intermediate/v1/intermediate_pb'

const Page = () => {
  const isDarkMode = useDarkMode()
  const settings = useSettings()

  const [favicon, setFavicon] = useState<string>('/apple-touch-icon.png')

  const applyTheme = () => {
    console.log('Applying theme:', settings)

    const root = document.documentElement
    const primary = isDarkMode
      ? settings?.darkModePrimaryColor
      : settings?.primaryColor

    if (primary) {
      const foreground = isColorDark(primary) ? '0 0% 100%' : '0 0% 0%'

      root.style.setProperty('--primary', hexToHSL(primary))
      root.style.setProperty('--primary-foreground', foreground)

      console.log(
        'Primary:',
        getComputedStyle(root).getPropertyValue('--primary').trim(),
      )

      console.log(
        'Primary foreground:',
        getComputedStyle(root).getPropertyValue('--primary-foreground').trim(),
      )
    }
  }

  useEffect(() => {
    if (settings) {
      applyTheme()
    }

    if (settings?.faviconUrl) {
      ;(async () => {
        try {
          // Check if the favicon exists before setting it
          const faviconCheck = await fetch(settings?.faviconUrl, {
            method: 'HEAD',
          })

          setFavicon(
            faviconCheck.ok ? settings?.faviconUrl : '/apple-touch-icon.png',
          )
        } catch (error) {
          setFavicon('/apple-touch-icon.png')
        }
      })()
    }
  }, [settings])

  useEffect(() => {
    applyTheme()
  }, [isDarkMode])

  return (
    <>
      <Helmet>
        <link rel="icon" href={favicon} />
        <link rel="apple-touch-icon" href={favicon} />
      </Helmet>

      <div
        className={cn(
          'min-h-screen w-screen',
          isDarkMode && settings?.detectDarkModeEnabled
            ? 'dark bg-dark'
            : 'light bg-body',
        )}
      >
        <div className="bg-body w-screen min-h-screen mx-auto flex flex-col justify-center items-center py-8">
          <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 flex justify-center">
            <div className="mb-8">
              {/* TODO: Make this conditionally load an Organizations configured logo */}
              {isDarkMode && settings?.detectDarkModeEnabled ? (
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
