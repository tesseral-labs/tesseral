import { useNavigate, useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  getProject,
  getSAMLConnection,
  updateOrganization,
  updateSAMLConnection,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import React, { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Switch } from '@/components/ui/switch'
import { Link } from 'react-router-dom'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { toast } from 'sonner'

const schema = z.object({
  idpEntityId: z.string().min(1, {
    message: 'IDP Entity ID must be non-empty.',
  }),
  idpRedirectUrl: z.string().url({
    message: 'IDP Redirect URL must be a valid URL.',
  }),
  idpX509Certificate: z.string(),
  // idpX509Certificate: z.string().startsWith('-----BEGIN CERTIFICATE-----', {
  //   message: 'IDP Certificate must be a PEM-encoded X.509 certificate.',
  // }),
})

export function EditSAMLConnectionPage() {
  const navigate = useNavigate()
  const { organizationId, samlConnectionId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: getSAMLConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  })
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {},
  })
  const updateSAMLConnectionMutation = useMutation(updateSAMLConnection)

  useEffect(() => {
    if (getSAMLConnectionResponse?.samlConnection) {
      form.reset({
        idpEntityId: getSAMLConnectionResponse.samlConnection.idpEntityId,
        idpRedirectUrl: getSAMLConnectionResponse.samlConnection.idpRedirectUrl,
        idpX509Certificate:
          getSAMLConnectionResponse.samlConnection.idpX509Certificate,
      })
    }
  }, [getSAMLConnectionResponse])

  async function onSubmit(values: z.infer<typeof schema>) {
    await updateSAMLConnectionMutation.mutateAsync({
      id: samlConnectionId,
      samlConnection: {
        idpEntityId: values.idpEntityId,
        idpRedirectUrl: values.idpRedirectUrl,
        idpX509Certificate: values.idpX509Certificate,
      },
    })

    toast.success('SAML Connection updated')
    navigate(
      `/organizations/${organizationId}/saml-connections/${samlConnectionId}`,
    )
  }

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
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}/saml-connections`}>
                SAML Connections
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link
                to={`/organizations/${organizationId}/saml-connections/${samlConnectionId}`}
              >
                {samlConnectionId}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Edit</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <h1 className="mt-4 mb-8 font-semibold text-2xl">Edit SAML Connection</h1>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>Service Provider Configuration</CardTitle>
              <CardDescription>
                The configuration here is assigned automatically by Tesseral,
                and needs to be inputted into your customer's Identity Provider
                by their IT admin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <div>
                <div className="text-sm font-medium leading-none">
                  Assertion Consumer Service (ACS) URL
                </div>
                <div>{getSAMLConnectionResponse?.samlConnection?.spAcsUrl}</div>
              </div>
              <div>
                <div className="text-sm font-medium leading-none">
                  SP Entity ID
                </div>
                <div>
                  {getSAMLConnectionResponse?.samlConnection?.spEntityId}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Identity Provider Configuration</CardTitle>
              <CardDescription>
                The configuration here needs to be copied over from the
                customer's Identity Provider ("IDP").
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="idpEntityId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>IDP Entity ID</FormLabel>
                    <FormControl>
                      <Input className="max-w-96" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="idpRedirectUrl"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>IDP Redirect URL</FormLabel>
                    <FormControl>
                      <Input className="max-w-96" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="idpX509Certificate"
                render={({ field: { onChange } }) => (
                  <FormItem>
                    <FormLabel>IDP Certificate</FormLabel>
                    <FormControl>
                      <Input
                        className="max-w-96"
                        type="file"
                        onChange={async (e) => {
                          // File inputs are special; they are necessarily "uncontrolled", and their value is a FileList.
                          // We just copy over the file's contents to the react-form-hook state manually on input change.
                          if (e.target.files) {
                            onChange(await e.target.files[0].text())
                          }
                        }}
                      />
                    </FormControl>
                    <FormDescription>
                      IDP Certificate, as a PEM-encoded X.509 certificate. These
                      start with '-----BEGIN CERTIFICATE-----' and end with
                      '-----END CERTIFICATE-----'.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          <div className="flex justify-end gap-x-4 pb-8">
            <Button variant="outline" asChild>
              <Link to={`/organizations/${organizationId}`}>Cancel</Link>
            </Button>
            <Button type="submit">Save Changes</Button>
          </div>
        </form>
      </Form>
    </div>
  )
}
