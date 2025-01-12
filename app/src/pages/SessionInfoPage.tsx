import { Title } from '@/components/Title'
import { whoAmI } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { useQuery } from '@connectrpc/connect-query'
import React from 'react'
import { useAccessToken } from '@/lib/use-access-token'

const SessionInfoPage = () => {
  const { data: whoamiRes } = useQuery(whoAmI)
  const accessToken = useAccessToken()

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1>Hello, {whoamiRes?.email} {accessToken}</h1>
      </div>
    </>
  )
}

export default SessionInfoPage
