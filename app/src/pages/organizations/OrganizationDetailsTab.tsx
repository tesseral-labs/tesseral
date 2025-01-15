import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { useParams } from 'react-router'
import { useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  getProject,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'

export function OrganizationDetailsTab() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: getProjectResponse } = useQuery(getProject, {})

  return (
    <Card className="my-8">
      <CardHeader className="py-4">
        <CardTitle className="text-xl">Details</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-3 gap-x-2">
          <div className="border-r border-gray-200 pr-8 flex flex-col gap-4">
            <div>
              <div className="font-semibold">Override Login Methods</div>
              <div>
                {getOrganizationResponse?.organization?.overrideLogInMethods
                  ? 'Yes'
                  : 'No'}
              </div>
            </div>

            {getProjectResponse?.project?.logInWithGoogleEnabled && (
              <div>
                <div className="font-semibold">Log in with Google</div>
                <div>
                  {getOrganizationResponse?.organization?.logInWithGoogleEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </div>
              </div>
            )}

            {getProjectResponse?.project?.logInWithMicrosoftEnabled && (
              <div>
                <div className="font-semibold">Log in with Microsoft</div>
                <div>
                  {getOrganizationResponse?.organization
                    ?.logInWithMicrosoftEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </div>
              </div>
            )}

            {getProjectResponse?.project?.logInWithPasswordEnabled && (
              <div>
                <div className="font-semibold">Log in with Password</div>
                <div>
                  {getOrganizationResponse?.organization
                    ?.logInWithPasswordEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </div>
              </div>
            )}
          </div>
          <div className="border-r border-gray-200 pl-8 pr-8 flex flex-col gap-4">
            {getProjectResponse?.project?.logInWithGoogleEnabled && (
              <div>
                <div className="font-semibold">Google Hosted Domain</div>
                <div>
                  {getOrganizationResponse?.organization?.googleHostedDomain ||
                    '-'}
                </div>
              </div>
            )}

            {getProjectResponse?.project?.logInWithMicrosoftEnabled && (
              <div>
                <div className="font-semibold">Microsoft Tenant ID</div>
                <div>
                  {getOrganizationResponse?.organization?.microsoftTenantId ||
                    '-'}
                </div>
              </div>
            )}
          </div>
          <div className="border-gray-200 pl-8 flex flex-col gap-4">
            <div>
              <div className="font-semibold">Configuring SAML</div>
              <div>
                {getOrganizationResponse?.organization?.samlEnabled
                  ? 'Enabled'
                  : 'Disabled'}
              </div>
            </div>
            <div>
              <div className="font-semibold">Configuring SCIM</div>
              <div>
                {getOrganizationResponse?.organization?.scimEnabled
                  ? 'Enabled'
                  : 'Disabled'}
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
