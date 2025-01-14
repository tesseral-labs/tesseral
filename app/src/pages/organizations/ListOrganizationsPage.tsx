import React from 'react'
import { useQuery } from '@connectrpc/connect-query'
import {
  listOrganizations
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'

export function ListOrganizationsPage() {
  const { data: organizations } = useQuery(listOrganizations, {})

  return (
    <div>
      {organizations?.organizations?.map(org => (
        <div key={org.id}>{org.displayName}</div>
      ))}
    </div>
  )
}
