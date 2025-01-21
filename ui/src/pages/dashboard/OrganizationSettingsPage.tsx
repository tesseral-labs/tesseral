import React, { FC, MouseEvent } from 'react'
import { DateTime } from 'luxon'
import { useOrganization, useProject, useUser } from '@/lib/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  listSAMLConnections,
  listUsers,
  updateUser,
} from '@/gen/openauth/frontend/v1/frontend-FrontendService_connectquery'
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

const OrganizationSettingsPage: FC = () => {
  const organization = useOrganization()
  const project = useProject()
  const user = useUser()

  const { data: usersData, refetch: refetchUsers } = useQuery(listUsers)
  const { data: samlConnectionsData, refetch: refetchSAMLConnections } =
    useQuery(listSAMLConnections)
  const updateUserMutation = useMutation(updateUser)

  const changeRole = async (userId: string, isOwner: boolean) => {
    await updateUserMutation.mutateAsync({
      id: userId,
      user: {
        owner: isOwner,
      },
    })

    refetchUsers()
  }

  return (
    <div className="dark:text-foreground">
      <div className="mb-4">
        <h1 className="text-2xl font-bold mb-2">{organization?.displayName}</h1>
        <span className="text-xs border px-2 py-1 rounded text-gray-400 dark:text-gray-700 bg-gray-200 dark:bg-gray-900 dark:border-gray-800">
          {organization?.id}
        </span>
      </div>

      <Card className="my-8">
        <CardHeader className="py-4">
          <CardTitle className="text-xl">General configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-x-2 text-sm">
            <div className="border-r border-gray-200 pr-8 dark:border-gray-700">
              <div className="font-semibold mb-2">Display Name</div>
              <div className="text-sm text-gray-500">
                {organization?.displayName}
              </div>
            </div>
            <div className="border-r border-gray-200 pl-8 pr-8 dark:border-gray-700">
              <div className="font-semibold mb-2">Created</div>
              <div className="text-sm text-gray-500">
                {organization?.createTime &&
                  DateTime.fromJSDate(
                    new Date(organization.updateTime),
                  ).toRelative()}
              </div>
            </div>
            <div className="px-8">
              <div className="font-semibold mb-2">Last updated</div>
              <div className="text-sm text-gray-500">
                {organization?.updateTime &&
                  DateTime.fromJSDate(
                    new Date(organization.updateTime),
                  ).toRelative()}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="my-8">
        <CardHeader className="py-4">
          <CardTitle className="text-xl">Users</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Role</TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {usersData?.users.map((u) => (
                <TableRow key={u.id}>
                  <TableCell className="flex items-center">{u.id}</TableCell>
                  <TableCell className="text-gray-500">{u.email}</TableCell>
                  <TableCell className="text-gray-500">
                    {u.owner ? 'Owner' : 'Member'}

                    {u.owner && u.id !== user?.id && (
                      <div
                        className="ml-2 rounded cursor-pointer text-primary border-border px-4 py-2 inline-block"
                        onClick={async (e: MouseEvent<HTMLSpanElement>) => {
                          e.stopPropagation()
                          e.preventDefault()

                          await changeRole(u.id, !u.owner)
                        }}
                      >
                        Make {u.owner ? 'Member' : 'Owner'}
                      </div>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card className="my-8">
        <CardHeader className="py-4">
          <CardTitle className="text-xl">SAML Connections</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>IDP Entity ID</TableCell>
                <TableCell>IDP Redirect URL</TableCell>
                <TableCell>IDP X509 Certificate</TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {samlConnectionsData?.samlConnections.map((c) => (
                <TableRow key={c.id}>
                  <TableCell className="flex items-center">{c.id}</TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpEntityId}
                  </TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpRedirectUrl}
                  </TableCell>
                  <TableCell className="text-gray-500">
                    {c.idpX509Certificate ? '✓' : 'x'}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default OrganizationSettingsPage
