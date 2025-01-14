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

export function ListOrganizationsPage() {
  const { data: listOrganizationsResponse } = useQuery(listOrganizations, {})

  return (
    <div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Display Name</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listOrganizationsResponse?.organizations?.map((org) => (
            <TableRow key={org.id}>
              <TableCell>{org.id}</TableCell>
              <TableCell className="font-medium">{org.displayName}</TableCell>
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
