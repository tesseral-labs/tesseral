import { useNavigate, useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  getOrganization,
  updateOrganization,
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

const schema = z.object({
  displayName: z.string(),
  overrideLogInMethods: z.boolean(),
  logInWithPasswordEnabled: z.boolean(),
  logInWithGoogleEnabled: z.boolean(),
  logInWithMicrosoftEnabled: z.boolean(),
  samlEnabled: z.boolean(),
  scimEnabled: z.boolean(),
})

export function EditOrganizationPage() {
  const navigate = useNavigate()
  const { organizationId } = useParams()
  const { data: getOrganizationResponse, refetch: refetchOrganization } =
    useQuery(getOrganization, {
      id: organizationId,
    })
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      // form needs to read this to change appearance; just default its value initially
      overrideLogInMethods: false,
    },
  })
  const updateOrganizationMutation = useMutation(updateOrganization)

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        displayName: getOrganizationResponse.organization.displayName,
        overrideLogInMethods:
          getOrganizationResponse.organization.overrideLogInMethods,
        logInWithGoogleEnabled:
          getOrganizationResponse.organization.logInWithGoogleEnabled,
        logInWithPasswordEnabled:
          getOrganizationResponse.organization.logInWithPasswordEnabled,
        logInWithMicrosoftEnabled:
          getOrganizationResponse.organization.logInWithMicrosoftEnabled,
        samlEnabled: getOrganizationResponse.organization.samlEnabled,
        scimEnabled: getOrganizationResponse.organization.scimEnabled,
      })
    }
  }, [getOrganizationResponse])

  async function onSubmit(values: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        displayName: values.displayName,
        overrideLogInMethods: values.overrideLogInMethods,
        logInWithGoogleEnabled: values.logInWithGoogleEnabled,
        logInWithPasswordEnabled: values.logInWithPasswordEnabled,
        logInWithMicrosoftEnabled: values.logInWithMicrosoftEnabled,
        samlEnabled: values.samlEnabled,
        scimEnabled: values.scimEnabled,
      },
    })

    navigate(`/organizations/${organizationId}`)
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
            <BreadcrumbLink>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Edit</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <h1 className="mt-4 mb-8 font-semibold text-2xl">
        Edit {getOrganizationResponse?.organization?.displayName}
      </h1>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>Organization settings</CardTitle>
              <CardDescription>
                Configure basic settings on this organization.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormControl>
                      <Input
                        className="max-w-80"
                        placeholder="Acme Corporation"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      A human-friendly name for the organization.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Login settings</CardTitle>
              <CardDescription>
                Configure what login methods users can use to log into this
                organization.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="overrideLogInMethods"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Override Login Methods</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      If your customer wants to restrict their login methods to
                      a subset of the project-level ones, you can have their
                      organization override the supported login methods.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="logInWithGoogleEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Google</FormLabel>
                    <FormControl>
                      <Switch
                        disabled={
                          !form.getValues('overrideLogInMethods').valueOf()
                        }
                        className="block"
                        checked={
                          form.getValues('overrideLogInMethods').valueOf()
                            ? field.value
                            : true
                        }
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      You can only modify this if "Override Login Methods" is
                      enabled for this organization.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="logInWithMicrosoftEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Microsoft</FormLabel>
                    <FormControl>
                      <Switch
                        disabled={
                          !form.getValues('overrideLogInMethods').valueOf()
                        }
                        className="block"
                        checked={
                          form.getValues('overrideLogInMethods').valueOf()
                            ? field.value
                            : true
                        }
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      You can only modify this if "Override Login Methods" is
                      enabled for this organization.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="logInWithPasswordEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Password</FormLabel>
                    <FormControl>
                      <Switch
                        disabled={
                          !form.getValues('overrideLogInMethods').valueOf()
                        }
                        className="block"
                        checked={
                          form.getValues('overrideLogInMethods').valueOf()
                            ? field.value
                            : true
                        }
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      You can only modify this if "Override Login Methods" is
                      enabled for this organization.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Enterprise settings</CardTitle>
              <CardDescription>
                Configure whether this organization can use SAML or SCIM.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="samlEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SAML Enabled</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Whether this organization can configure SAML ("Enterprise
                      Single Sign-On").
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="scimEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SCIM Enabled</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Whether this organization can configure SCIM ("Enterprise
                      Directory Sync").
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
