import { useQuery } from '@connectrpc/connect-query'
import { getProject } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import React from 'react'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'

export function ViewProjectPage() {
  const { data: getProjectResponse } = useQuery(getProject, {})

  return (
    <div>
      <h1 className="font-semibold text-2xl">
        {getProjectResponse?.project?.displayName}
      </h1>

      {getProjectResponse?.project && (
        <div>
          <div>ID: {getProjectResponse?.project?.id}</div>

          <div>
            Created At:
            {DateTime.fromJSDate(
              timestampDate(getProjectResponse.project.createTime!),
            ).toRelative()}
          </div>
          <div>
            Updated At:
            {DateTime.fromJSDate(
              timestampDate(getProjectResponse.project.updateTime!),
            ).toRelative()}
          </div>
          <div>
            Log in with Google Enabled?
            {getProjectResponse.project.logInWithPasswordEnabled ? 'Yes' : 'No'}
          </div>
          <div>
            Log in with Microsoft Enabled?
            {getProjectResponse.project.logInWithMicrosoftEnabled
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Log in with Password Enabled?
            {getProjectResponse.project.logInWithPasswordEnabled ? 'Yes' : 'No'}
          </div>
          <div>
            Google OAuth Client ID:
            {getProjectResponse.project.googleOauthClientId}
          </div>
          <div>
            Microsoft OAuth Client ID:
            {getProjectResponse.project.microsoftOauthClientId}
          </div>
          <div>
            Projects have SAML enabled by default:
            {getProjectResponse.project.organizationsSamlEnabledDefault
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Projects have SCIM enabled by default:
            {getProjectResponse.project.organizationsScimEnabledDefault
              ? 'Yes'
              : 'No'}
          </div>
          <div>
            Auth domain:
            {getProjectResponse.project.authDomain}
          </div>
        </div>
      )}
    </div>
  )
}
