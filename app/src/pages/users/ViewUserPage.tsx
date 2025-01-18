import { useParams } from 'react-router'
import { useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
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
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { PageCodeSubtitle, PageDescription, PageTitle } from '@/components/page'
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'

export function ViewUserPage() {
  const { organizationId, userId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  })
  const { data: listSessionsResponse } = useQuery(listSessions, {
    userId,
  })

  return (
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/organizations">Organizations</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}/users`}>Users</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{getUserResponse?.user?.email}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>{getUserResponse?.user?.email}</PageTitle>
      <PageCodeSubtitle>{userId}</PageCodeSubtitle>
      <PageDescription>
        A user is what people using your product log into. Lorem ipsum dolor.
      </PageDescription>

      <Card className="my-8">
        <CardHeader>
          <CardTitle>General settings</CardTitle>
          <CardDescription>Basic settings for this user.</CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Email</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.email}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Owner</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.owner ? 'Yes' : 'No'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google User ID</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.googleUserId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft User ID</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.microsoftUserId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Created</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(getUserResponse?.user?.createTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Updated</DetailsGridKey>
                <DetailsGridValue>
                  {getUserResponse?.user?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(getUserResponse?.user?.updateTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Sessions</CardTitle>
          <CardDescription>
            Every time your users log in or perform an action, that's associated
            with a session. Lorem ipsum dolor.
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
              {listSessionsResponse?.sessions?.map((session) => (
                <TableRow key={session.id}>
                  <TableCell className="font-medium font-mono">
                    <Link
                      className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
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
        </CardContent>
      </Card>
    </div>
  )
}
