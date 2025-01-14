import { useQuery } from '@connectrpc/connect-query'
import { getOrganization } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { useParams } from 'react-router'
import React from 'react'

export function ViewOrganizationPage() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })

  return (
    <div>
      {JSON.stringify(getOrganizationResponse?.organization?.displayName)}
    </div>
  )
}
