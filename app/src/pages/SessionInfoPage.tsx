import { Title } from '@/components/Title'
import { whoAmI } from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import { useQuery } from '@connectrpc/connect-query'
import React from 'react'

const SessionInfoPage = () => {
  const { data: whoamiRes } = useQuery(whoAmI)

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1>Hello, {whoamiRes?.email}</h1>
      </div>
    </>
  )
}

export default SessionInfoPage
