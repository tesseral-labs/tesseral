import React from 'react'
import { Outlet } from 'react-router'
import useDarkMode from '@/lib/dark-mode'
import { cn } from '@/lib/utils'

const Page = () => {
  const isDarkMode = useDarkMode()

  return (
    <div
      className={cn(
        'mx-auto flex flex-col justify-center items-center min-h-screen w-screen py-8',
        isDarkMode ? 'dark bg-dark' : 'light bg-muted',
      )}
    >
      <div className="container flex justify-center">
        <div className="mb-8">
          {/* TODO: Make this conditionally load an Organizations configured logo */}
          {isDarkMode ? (
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
