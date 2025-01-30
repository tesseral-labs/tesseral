import React, { FC, SyntheticEvent } from 'react'
import useDarkMode from '@/lib/dark-mode'
import useSettings from '@/lib/settings'
import { Link, useLocation } from 'react-router-dom'
import { cn } from '@/lib/utils'

const Header: FC = () => {
  const isDarkMode = useDarkMode()
  const location = useLocation()
  const settings = useSettings()

  return (
    <header className="flex flex-row justify-between text-foreground py-4 mb-4 border-b dark:border-gray-800">
      <div>
        {isDarkMode && settings.detectDarkModeEnabled ? (
          <img
            className="max-h-[30px] max-w-[150px]"
            src={settings?.darkModeLogoUrl || '/images/tesseral-logo-white.svg'}
            onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
              const target = e.target as HTMLImageElement
              target.onerror = null
              target.src = '/images/tesseral-logo-white.svg'
            }}
          />
        ) : (
          <img
            className="max-h-[30px] max-w-[150px]"
            src={settings?.logoUrl || '/images/tesseral-logo-black.svg'}
            onError={(e: SyntheticEvent<HTMLImageElement, Event>) => {
              const target = e.target as HTMLImageElement
              target.onerror = null
              target.src = '/images/tesseral-logo-black.svg'
            }}
          />
        )}
      </div>
      <div>
        <Link
          className={cn(
            'px-4',
            location.pathname === '/settings' ? 'font-bold' : '',
          )}
          to="/settings"
        >
          User Settings
        </Link>
        <Link
          className={cn(
            'px-4',
            location.pathname === '/organization' ? 'font-bold' : '',
          )}
          to="/organization"
        >
          Organization Settings
        </Link>
      </div>
    </header>
  )
}

export default Header
