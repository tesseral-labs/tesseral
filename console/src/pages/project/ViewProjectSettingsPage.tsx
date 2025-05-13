import { useQuery } from '@connectrpc/connect-query';
import { getProject } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React from 'react';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { Outlet, useLocation } from 'react-router';
import { TabBar, TabBarLink } from '@/components/ui/tab-bar';
import { Settings2 } from 'lucide-react';

export const ViewProjectSettingsPage = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { pathname } = useLocation();

  const tabs = [
    {
      root: true,
      name: 'Details',
      url: `/project-settings`,
    },
    {
      name: 'Vault UI Settings',
      url: `/project-settings/vault-ui-settings`,
    },
    {
      name: 'Vault Domain Settings',
      url: `/project-settings/vault-domain-settings`,
    },
    {
      name: 'Role-Based Access Control Settings',
      url: `/project-settings/rbac-settings`,
    },
  ];

  const currentTab = tabs.find((tab) => tab.url === pathname);

  return (
    <>
      <TabBar>
        {tabs.map((tab) => (
          <TabBarLink
            key={tab.name}
            active={tab.url === currentTab?.url}
            url={tab.url}
            label={tab.name}
          />
        ))}
      </TabBar>
      <PageHeader>
        <PageTitle className="flex items-center">
          <Settings2 className="inline mr-2 w-6 h-6" />
          Project settings
        </PageTitle>
        <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
        <PageDescription>
          Everything you do in Tesseral happens inside a Project.
        </PageDescription>
      </PageHeader>
      <PageContent>
        <Card className="my-8">
          <CardHeader>
            <CardTitle>General configuration</CardTitle>
          </CardHeader>

          <CardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Display name</DetailsGridKey>
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

        <div className="mt-4">
          <Outlet />
        </div>
      </PageContent>
    </>
  );
};
