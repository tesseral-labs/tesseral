import { useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  listSAMLConnections,
  listSCIMAPIKeys,
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
  const { data: listSAMLConnectionsResponse } = useQuery(listSAMLConnections, {
    organizationId,
  })
  const { data: listSCIMAPIKeysResponse } = useQuery(listSCIMAPIKeys, {
    organizationId,
  })

  return (
    <div>
      <h1 className="font-semibold text-2xl">
        {getOrganizationResponse?.organization?.displayName}
      </h1>

      <div>ID: {organizationId}</div>

      {getOrganizationResponse?.organization && (
        <div>
          <div>
            Created At:
            {DateTime.fromJSDate(
              timestampDate(getOrganizationResponse.organization.createTime!),
            ).toRelative()}
          </div>
          <div>
            Updated At:
            {DateTime.fromJSDate(
              timestampDate(getOrganizationResponse.organization.updateTime!),
            ).toRelative()}
          </div>
          <div>
            Override Log in Methods?
            {getOrganizationResponse.organization.overrideLogInMethods
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Log in with Google?
            {getOrganizationResponse.organization.logInWithGoogleEnabled
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Override Log in Microsoft?
            {getOrganizationResponse.organization.logInWithMicrosoftEnabled
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Override Log in Password?
            {getOrganizationResponse.organization.logInWithPasswordEnabled
              ? 'Yes'
              : 'No'}
          </div>

          <div>
            Google Hosted Domain:
            {getOrganizationResponse.organization.googleHostedDomain}
          </div>
          <div>
            Microsoft Tenant ID:
            {getOrganizationResponse.organization.microsoftTenantId}
          </div>

          <div>
            SAML Enabled?
            {getOrganizationResponse.organization.samlEnabled}
          </div>

          <div>
            SCIM Enabled?
            {getOrganizationResponse.organization.scimEnabled}
          </div>
        </div>
      )}

      <h2 className="font-semibold text-xl">Users</h2>
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

      <h2 className="font-semibold text-xl">SAML Connections</h2>
      <Table className="mt-4">
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
                <TableCell className="font-medium">
                  <Link
                    to={`/organizations/${organizationId}/saml-connections/${samlConnection.id}`}
                  >
                    {samlConnection.id}
                  </Link>
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
      <h2 className="font-semibold text-xl">SCIM API Keys</h2>
      <Table className="mt-4">
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listSCIMAPIKeysResponse?.scimApiKeys?.map((scimAPIKey) => (
            <TableRow key={scimAPIKey.id}>
              <TableCell className="font-medium">
                <Link
                  to={`/organizations/${organizationId}/scim-api-keys/${scimAPIKey.id}`}
                >
                  {scimAPIKey.id}
                </Link>
              </TableCell>
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
    </div>
  )
}
