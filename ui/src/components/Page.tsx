import React from 'react'
import { Outlet } from 'react-router'

const Page = () => {
  return (
    <div className="container mx-auto flex flex-col justify-center items-center h-screen py-8">
      <div className="flex justify-center">
        <div className="mb-8">
          {/* TODO: Make this conditionally load an Organizations configured logo */}
          <img
            className="max-w-[240px]"
            src="/images/tesseral-logo-black.svg"
          />
        </div>
      </div>
      <Outlet />
    </div>
  )
}

export default Page
