import React from 'react'

import { useUser } from '@/lib/auth'
import { Title } from '@/components/Title'

const SessionInfoPage = () => {
  const user = useUser()

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1 className="text-foreground">Hello, {user?.email}</h1>
      </div>
    </>
  )
}

export default SessionInfoPage
