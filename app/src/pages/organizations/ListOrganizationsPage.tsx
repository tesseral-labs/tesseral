import React from 'react'
import { useQuery } from '@connectrpc/connect-query'
import { listOrganizations } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { Link } from 'react-router-dom'

export function ListOrganizationsPage() {
  const { data: listOrganizationsResponse } = useQuery(listOrganizations, {})

  return (
    <div>
      <h1 className="font-semibold text-2xl">Organizations</h1>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Display Name</TableHead>
            <TableHead>ID</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listOrganizationsResponse?.organizations?.map((org) => (
            <TableRow key={org.id}>
              <TableCell className="font-medium">
                <Link to={`/organizations/${org.id}`}>{org.displayName}</Link>
              </TableCell>
              <TableCell className="font-mono">{org.id}</TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(org.createTime!),
                ).toRelative()}
              </TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(org.updateTime!),
                ).toRelative()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
