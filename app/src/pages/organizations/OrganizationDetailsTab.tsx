import React, { useEffect, useState } from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  getOrganizationGoogleHostedDomains,
  getOrganizationMicrosoftTenantIDs,
  getProject,
  getProjectAPIKey,
  updateOrganizationGoogleHostedDomains,
  updateProjectAPIKey,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { Button } from '@/components/ui/button'
import { Link } from 'react-router-dom'
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { InputTags } from '@/components/input-tags'
import { EditOrganizationGoogleConfigurationButton } from '@/pages/organizations/EditOrganizationGoogleConfigurationButton'
import { EditOrganizationMicrosoftConfigurationButton } from '@/pages/organizations/EditOrganizationMicrosoftConfigurationButton'

export function OrganizationDetailsTab() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: getProjectResponse } = useQuery(getProject, {})
  const { data: getOrganizationGoogleHostedDomainsResponse } = useQuery(
    getOrganizationGoogleHostedDomains,
    {
      organizationId,
    },
  )
  const { data: getOrganizationMicrosoftTenantIdsResponse } = useQuery(
    getOrganizationMicrosoftTenantIDs,
    {
      organizationId,
    },
  )

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Details</CardTitle>
            <CardDescription>
              Additional details about your organization. Lorem ipsum dolor.
            </CardDescription>
          </div>
          <Button variant="outline" asChild>
            <Link to={`/organizations/${organizationId}/edit`}>Edit</Link>
          </Button>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Override Login Methods</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.overrideLogInMethods
                    ? 'Yes'
                    : 'No'}
                </DetailsGridValue>
              </DetailsGridEntry>

              {getProjectResponse?.project?.logInWithGoogleEnabled && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Google</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization
                      ?.logInWithGoogleEnabled
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithMicrosoftEnabled && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization
                      ?.logInWithMicrosoftEnabled
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}

              {getProjectResponse?.project?.logInWithPasswordEnabled && (
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Password</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationResponse?.organization
                      ?.logInWithPasswordEnabled
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              )}
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Configuring SAML</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.samlEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Configuring SCIM</DetailsGridKey>
                <DetailsGridValue>
                  {getOrganizationResponse?.organization?.scimEnabled
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      {getOrganizationResponse?.organization?.logInWithGoogleEnabled && (
        <Card>
          <CardHeader className="flex-row justify-between items-center">
            <div className="flex flex-col space-y-1 5">
              <CardTitle>Google configuration</CardTitle>
              <CardDescription>
                Settings related to logging into this organization with Google.
              </CardDescription>
            </div>
            <EditOrganizationGoogleConfigurationButton />
          </CardHeader>
          <CardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Google</DetailsGridKey>
                  <DetailsGridValue>Enabled</DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Google Hosted Domains</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationGoogleHostedDomainsResponse
                      ?.organizationGoogleHostedDomains?.googleHostedDomains
                      ? getOrganizationGoogleHostedDomainsResponse.organizationGoogleHostedDomains.googleHostedDomains.map(
                          (s) => <div key={s}>{s}</div>,
                        )
                      : '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </CardContent>
        </Card>
      )}

      {getOrganizationResponse?.organization?.logInWithMicrosoftEnabled && (
        <Card>
          <CardHeader className="flex-row justify-between items-center">
            <div className="flex flex-col space-y-1 5">
              <CardTitle>Microsoft configuration</CardTitle>
              <CardDescription>
                Settings related to logging into this organization with
                Microsoft.
              </CardDescription>
            </div>
            <EditOrganizationMicrosoftConfigurationButton />
          </CardHeader>
          <CardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                  <DetailsGridValue>Enabled</DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Microsoft Tenant IDs</DetailsGridKey>
                  <DetailsGridValue>
                    {getOrganizationMicrosoftTenantIdsResponse
                      ?.organizationMicrosoftTenantIds?.microsoftTenantIds
                      ? getOrganizationMicrosoftTenantIdsResponse.organizationMicrosoftTenantIds.microsoftTenantIds.map(
                          (s) => <div key={s}>{s}</div>,
                        )
                      : '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
