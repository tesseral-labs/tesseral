import { Title } from '@/components/Title'
import { whoami } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useQuery } from '@connectrpc/connect-query'
import React from 'react'

const SessionInfoPage = () => {
  const { data: whoamiRes } = useQuery(whoami)

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1 className="text-foreground">Hello, {whoamiRes?.email}</h1>
      </div>
    </>
  )
}

export default SessionInfoPage
