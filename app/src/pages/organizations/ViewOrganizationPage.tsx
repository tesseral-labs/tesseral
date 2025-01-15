import { useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  listSAMLConnections,
  listSCIMAPIKeys,
  listUsers,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { Outlet, useLocation, useParams } from 'react-router'
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
import { ChevronDownIcon } from 'lucide-react'
import { clsx } from 'clsx'

export function ViewOrganizationPage() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { pathname } = useLocation()

  const tabs = [
    {
      name: 'Users',
      url: `/organizations/${organizationId}/users`,
      alternativeUrl: `/organizations/${organizationId}`,
    },
    {
      name: 'SAML Connections',
      url: `/organizations/${organizationId}/saml-connections`,
    },
    {
      name: 'SCIM API Keys',
      url: `/organizations/${organizationId}/scim-api-keys`,
    },
    {
      name: 'Organization Settings',
      url: `/organizations/${organizationId}/settings`,
    },
  ]

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

      <div className="border-b border-gray-200">
        <nav aria-label="Tabs" className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <Link
              key={tab.name}
              to={tab.url}
              className={clsx(
                pathname === tab.url || pathname === tab.alternativeUrl
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                'whitespace-nowrap border-b-2 px-1 py-4 text-sm font-medium',
              )}
            >
              {tab.name}
            </Link>
          ))}
        </nav>
      </div>

      <div className="mt-4">
        <Outlet />
      </div>
    </div>
  )
}
