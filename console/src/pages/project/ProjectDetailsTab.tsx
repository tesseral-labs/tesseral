import React from 'react'
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
import { useQuery } from '@connectrpc/connect-query'
import { getProject } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { Button } from '@/components/ui/button'
import { Link } from 'react-router-dom'

export function ProjectDetailsTab() {
  const { data: getProjectResponse } = useQuery(getProject, {})

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>Authentication settings</CardTitle>
          <CardDescription>
            Configure the login methods your customers can use to log in to your
            application.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Password</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithPassword
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Google</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithGoogle
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithMicrosoft
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>
            <div className="grid grid-cols-2 gap-8">
              <span>Google settings</span>
              <div className="text-right">
                <Link to="/project-settings/log-in-with-google/edit">
                  <Button variant="outline" size="sm">
                    Edit
                  </Button>
                </Link>
              </div>
            </div>
          </CardTitle>
          <CardDescription>
            Settings for "Log in with Google" in your project.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithGoogle
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google OAuth Client ID</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.googleOauthClientId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google OAuth Client Secret</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.googleOauthClientId ? (
                    <div className="text-muted-foreground">Encrypted</div>
                  ) : (
                    '-'
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>
            <div className="grid grid-cols-2 gap-8">
              <span>Microsoft settings</span>
              <div className="text-right">
                <Link to="/project-settings/log-in-with-microsoft/edit">
                  <Button variant="outline" size="sm">
                    Edit
                  </Button>
                </Link>
              </div>
            </div>
          </CardTitle>
          <CardDescription>
            Settings for "Log in with Microsoft" in your project.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithMicrosoft
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft OAuth Client ID</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.microsoftOauthClientId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft OAuth Client Secret</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.microsoftOauthClientId ? (
                    <div className="text-muted-foreground">Encrypted</div>
                  ) : (
                    '-'
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>
    </div>
  )
}
