import React from 'react'
import { useParams } from 'react-router'
import { useQuery } from '@connectrpc/connect-query'
import {
  listSAMLConnections,
  listUsers,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
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
import { Badge } from '@/components/ui/badge'

export function OrganizationSAMLConnectionsTab() {
  const { organizationId } = useParams()
  const { data: listSAMLConnectionsResponse } = useQuery(listSAMLConnections, {
    organizationId,
  })

  return (
    <Card>
      <CardHeader>
        <CardTitle>SAML Connections</CardTitle>
        <CardDescription>
          A SAML connection is a link between Tesseral and your customer's
          enterprise Identity Provider. Lorem ipsum dolor.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listSAMLConnectionsResponse?.samlConnections?.map(
              (samlConnection) => (
                <TableRow key={samlConnection.id}>
                  <TableCell>
                    <Link
                      className="font-mono font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organizations/${organizationId}/saml-connections/${samlConnection.id}`}
                    >
                      {samlConnection.id}
                    </Link>

                    {samlConnection.primary && (
                      <Badge variant="outline" className="ml-2">
                        Primary
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell>
                    {DateTime.fromJSDate(
                      timestampDate(samlConnection.createTime!),
                    ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {DateTime.fromJSDate(
                      timestampDate(samlConnection.updateTime!),
                    ).toRelative()}
                  </TableCell>
                </TableRow>
              ),
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
