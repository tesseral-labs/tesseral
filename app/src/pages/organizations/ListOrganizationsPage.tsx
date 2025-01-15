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
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

export function ListOrganizationsPage() {
  const { data: listOrganizationsResponse } = useQuery(listOrganizations, {})

  return (
    <Card>
      <CardHeader>
        <CardTitle>Organizations</CardTitle>
        <CardDescription>
          An organization represents one of your business customers. Lorem ipsum
          dolor.
        </CardDescription>
      </CardHeader>
      <CardContent>
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
                  <Link
                    className="underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${org.id}`}
                  >
                    {org.displayName}
                  </Link>
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
      </CardContent>
    </Card>
  )
}
