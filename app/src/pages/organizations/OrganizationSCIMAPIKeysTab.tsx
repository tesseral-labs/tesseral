import React from 'react'
import { useParams } from 'react-router'
import { useQuery } from '@connectrpc/connect-query'
import {
  listSAMLConnections,
  listSCIMAPIKeys,
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

export function OrganizationSCIMAPIKeysTab() {
  const { organizationId } = useParams()
  const { data: listSCIMAPIKeysResponse } = useQuery(listSCIMAPIKeys, {
    organizationId,
  })

  return (
    <Card>
      <CardHeader>
        <CardTitle>SCIM API Keys</CardTitle>
        <CardDescription>
          A SCIM API key lets your customer do enterprise directory syncing.
          Lorem ipsum dolor.
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
            {listSCIMAPIKeysResponse?.scimApiKeys?.map((scimAPIKey) => (
              <TableRow key={scimAPIKey.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/saml-connections/${scimAPIKey.id}`}
                  >
                    {scimAPIKey.displayName}
                  </Link>
                </TableCell>
                <TableCell className="font-mono">{scimAPIKey.id}</TableCell>
                <TableCell>
                  {DateTime.fromJSDate(
                    timestampDate(scimAPIKey.createTime!),
                  ).toRelative()}
                </TableCell>
                <TableCell>
                  {DateTime.fromJSDate(
                    timestampDate(scimAPIKey.updateTime!),
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
