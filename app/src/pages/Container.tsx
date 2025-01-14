import { Outlet } from 'react-router'
import React from 'react'

export function Container() {
  return (
    <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <Outlet />
    </div>
  )
}
