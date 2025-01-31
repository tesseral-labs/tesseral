import React from 'react'
import { useNavigate, useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  createSAMLConnection,
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
import { Button } from '@/components/ui/button'
import { toast } from 'sonner'

export function OrganizationSAMLConnectionsTab() {
  const { organizationId } = useParams()
  const { data: listSAMLConnectionsResponse } = useQuery(listSAMLConnections, {
    organizationId,
  })
  const navigate = useNavigate()
  const createSAMLConnectionMutation = useMutation(createSAMLConnection)

  async function handleCreateSAMLConnection() {
    const { samlConnection } = await createSAMLConnectionMutation.mutateAsync({
      samlConnection: {
        organizationId,

        // if there are no saml connections on the org yet, default to making
        // the first one be primary
        primary: !!listSAMLConnectionsResponse?.samlConnections,
      },
    })

    toast.success('SAML Connection created')
    navigate(
      `/organizations/${organizationId}/saml-connections/${samlConnection!.id}`,
    )
  }

  return (
    <Card>
      <CardHeader className="flex-row justify-between items-center">
        <div className="flex flex-col space-y-1 5">
          <CardTitle>SAML Connections</CardTitle>
          <CardDescription>
            A SAML connection is a link between Tesseral and your customer's
            enterprise Identity Provider. Lorem ipsum dolor.
          </CardDescription>
        </div>
        <Button variant="outline" onClick={handleCreateSAMLConnection}>
          Create
        </Button>
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
