import { useQuery } from '@connectrpc/connect-query';
import { getOrganization } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery';
import { Outlet, useLocation, useParams } from 'react-router';
import React from 'react';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { clsx } from 'clsx';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  PageCodeSubtitle,
  PageDescription,
  PageTitle,
} from '@/components/page';

export const ViewOrganizationPage = () => {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { pathname } = useLocation();

  const tabs = [
    {
      root: true,
      name: 'Details',
      url: `/organizations/${organizationId}`,
    },
    {
      name: 'Users',
      url: `/organizations/${organizationId}/users`,
    },
    {
      name: 'User Invites',
      url: `/organizations/${organizationId}/user-invites`,
    },
    {
      name: 'SAML Connections',
      url: `/organizations/${organizationId}/saml-connections`,
    },
    {
      name: 'SCIM API Keys',
      url: `/organizations/${organizationId}/scim-api-keys`,
    },
  ];

  const currentTab = tabs.find((tab) => tab.url === pathname);

  return (
    // TODO remove padding when app shell in place
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
          {currentTab?.root ? (
            <BreadcrumbItem>
              <BreadcrumbPage>
                {getOrganizationResponse?.organization?.displayName}
              </BreadcrumbPage>
            </BreadcrumbItem>
          ) : (
            <>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to={`/organizations/${organizationId}`}>
                    {getOrganizationResponse?.organization?.displayName}
                  </Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>{currentTab?.name}</BreadcrumbPage>
              </BreadcrumbItem>
            </>
          )}
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>
        {getOrganizationResponse?.organization?.displayName}
      </PageTitle>
      <PageCodeSubtitle>{organizationId}</PageCodeSubtitle>
      <PageDescription>
        An organization represents one of your business customers. Lorem ipsum
        dolor.
      </PageDescription>

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
                tab.url === currentTab?.url
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
  );
};
