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
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export function ViewOrganizationPage() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { pathname } = useLocation()

  const tabs = [
    {
      name: 'Details',
      url: `/organizations/${organizationId}`,
    },
    {
      name: 'Users',
      url: `/organizations/${organizationId}/users`,
    },
    {
      name: 'SAML Connections',
      url: `/organizations/${organizationId}/saml-connections`,
    },
    {
      name: 'SCIM API Keys',
      url: `/organizations/${organizationId}/scim-api-keys`,
    },
  ]

  return (
    // TODO remove padding when app shell in place
    <div className="pt-8">
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
            <BreadcrumbPage>
              {getOrganizationResponse?.organization?.displayName}
            </BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <h1 className="mt-4 font-semibold text-2xl">
        {getOrganizationResponse?.organization?.displayName}
      </h1>
      <span className="mt-1 inline-block border rounded bg-gray-100 py-1 px-2 font-mono text-xs text-muted-foreground">
        {organizationId}
      </span>
      <div className="mt-4">
        An organization represents one of your business customers. Lorem ipsum
        dolor.
      </div>

      <Card className="my-8">
        <CardHeader className="py-4">
          <CardTitle className="text-xl">General configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-x-2 text-sm">
            <div className="border-r border-gray-200 pr-8">
              <div className="font-semibold">Display Name</div>
              <div>{getOrganizationResponse?.organization?.displayName}</div>
            </div>
            <div className="border-r border-gray-200 pl-8 pr-8">
              <div className="font-semibold">Created</div>
              <div>
                {getOrganizationResponse?.organization?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(
                      getOrganizationResponse.organization.createTime,
                    ),
                  ).toRelative()}
              </div>
            </div>
            <div className="px-8">
              <div className="font-semibold">Last updated</div>
              <div>
                {getOrganizationResponse?.organization?.updateTime &&
                  DateTime.fromJSDate(
                    timestampDate(
                      getOrganizationResponse.organization.updateTime,
                    ),
                  ).toRelative()}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="border-b border-gray-200">
        <nav aria-label="Tabs" className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <Link
              key={tab.name}
              to={tab.url}
              className={clsx(
                pathname === tab.url
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
