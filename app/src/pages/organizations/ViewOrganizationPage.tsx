import { useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  listUsers,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { useParams } from 'react-router'
import React from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Link } from 'react-router-dom'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'

export function ViewOrganizationPage() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: listUsersResponse } = useQuery(listUsers, {
    organizationId,
  })

  return (
    <div>
      <h1 className="font-semibold text-2xl">
        {getOrganizationResponse?.organization?.displayName}
      </h1>

      <h2 className="font-display text-xl">Users</h2>
      <Table className="mt-4">
        <TableHeader>
          <TableRow>
            <TableHead>Email</TableHead>
            <TableHead>ID</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listUsersResponse?.users?.map((user) => (
            <TableRow key={user.id}>
              <TableCell className="font-medium">
                <Link to={`/organizations/${organizationId}/users/${user.id}`}>
                  {user.email}
                </Link>
              </TableCell>
              <TableCell className="font-mono">{user.id}</TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(user.createTime!),
                ).toRelative()}
              </TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(user.updateTime!),
                ).toRelative()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
