import { useQuery } from '@connectrpc/connect-query'
import { getProject } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import React from 'react'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Link } from 'react-router-dom'
import { PageCodeSubtitle, PageDescription, PageTitle } from '@/components/page'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid'
import { clsx } from 'clsx'
import { Outlet, useLocation } from 'react-router'

export function ViewProjectSettingsPage() {
  const { data: getProjectResponse } = useQuery(getProject, {})
  const { pathname } = useLocation()

  const tabs = [
    {
      root: true,
      name: 'Details',
      url: `/project-settings`,
    },
    {
      name: 'Hosted Portal Settings',
      url: `/project-settings/hosted-portal`,
    },
  ]

  const currentTab = tabs.find((tab) => tab.url === pathname)!

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
            <BreadcrumbPage>Project Settings</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Project Settings</PageTitle>
      <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
      <PageDescription>
        Everything you do in Tesseral happens inside a project.
      </PageDescription>

      <Card className="my-8">
        <CardHeader>
          <CardTitle className="text-xl">General configuration</CardTitle>
        </CardHeader>

        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Display Name</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.displayName}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Created</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(getProjectResponse.project.createTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Updated</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(getProjectResponse.project.updateTime),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <div className="border-b border-gray-200">
        <nav aria-label="Tabs" className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <Link
              key={tab.name}
              to={tab.url}
              className={clsx(
                tab.url === currentTab.url
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
