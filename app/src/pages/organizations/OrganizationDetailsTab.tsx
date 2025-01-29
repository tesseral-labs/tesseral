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
            <EditGoogleConfigurationButton />
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
          <CardHeader>
            <CardTitle>Microsoft Configuration</CardTitle>
            <CardDescription>
              Settings related to logging into this organization with Microsoft.
            </CardDescription>
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

const schema = z.object({
  googleHostedDomains: z.array(z.string()),
})

function EditGoogleConfigurationButton() {
  const { organizationId } = useParams()
  const { data: getOrganizationGoogleHostedDomainsResponse, refetch } =
    useQuery(getOrganizationGoogleHostedDomains, {
      organizationId,
    })
  const updateOrganizationGoogleHostedDomainsMutation = useMutation(
    updateOrganizationGoogleHostedDomains,
  )
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      googleHostedDomains: [],
    },
  })
  useEffect(() => {
    if (
      getOrganizationGoogleHostedDomainsResponse?.organizationGoogleHostedDomains
    ) {
      form.reset({
        googleHostedDomains:
          getOrganizationGoogleHostedDomainsResponse
            .organizationGoogleHostedDomains.googleHostedDomains,
      })
    }
  }, [getOrganizationGoogleHostedDomainsResponse])

  const [open, setOpen] = useState(false)

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateOrganizationGoogleHostedDomainsMutation.mutateAsync({
      organizationId,
      organizationGoogleHostedDomains: {
        googleHostedDomains: values.googleHostedDomains,
      },
    })
    await refetch()
    setOpen(false)
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Google configuration</AlertDialogTitle>
          <AlertDialogDescription>
            Edit organization google configuration.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="googleHostedDomains"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Google Hosted Domains</FormLabel>
                  <FormControl>
                    <InputTags className="max-w-96" {...field} />
                  </FormControl>
                  <FormDescription>
                    The set of Google workspaces associated with this
                    organization. Google identifies workspaces by their "hosted
                    domains", e.g. "example.com".
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  )
}
