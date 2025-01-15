import { useParams } from 'react-router'
import { useQuery } from '@connectrpc/connect-query'
import {
  getUser,
  listSessions,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
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
import { timestampDate } from '@bufbuild/protobuf/dist/esm/wkt'

export function ViewUserPage() {
  const { organizationId, userId } = useParams()
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  })
  const { data: listSessionsResponse } = useQuery(listSessions, {
    userId,
  })

  return (
    <div>
      <h1 className="text-2xl">{getUserResponse?.user?.email}</h1>

      <div>ID: {getUserResponse?.user?.email}</div>
      <div>Google user ID: {getUserResponse?.user?.googleUserId}</div>
      <div>Microsoft user ID: {getUserResponse?.user?.microsoftUserId}</div>
      <div>Owner? {getUserResponse?.user?.owner ? 'yes' : 'no'}</div>

      <h2 className="font-semibold text-xl">Sessions</h2>
      <Table className="mt-4">
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listSessionsResponse?.sessions?.map((session) => (
            <TableRow key={session.id}>
              <TableCell className="font-medium font-mono">
                <Link
                  to={`/organizations/${organizationId}/users/${userId}/sessions/${session.id}`}
                >
                  {session.id}
                </Link>
              </TableCell>
              <TableCell>TODO</TableCell>
              <TableCell>TODO</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
