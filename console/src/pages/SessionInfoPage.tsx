import { Title } from '@/components/Title'
import { whoAmI } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { useQuery } from '@connectrpc/connect-query'
import React from 'react'
import { useAccessToken, useUser } from '@/lib/use-access-token'

const SessionInfoPage = () => {
  const user = useUser()

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1>Hello, {user?.id} {user?.email}</h1>
      </div>
    </>
  )
}

export default SessionInfoPage
