import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import React, { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery, useMutation } from '@connectrpc/connect-query'
import {
  getOrganization,
  updateOrganization,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
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
import { Checkbox } from '@/components/ui/checkbox'
import { Switch } from '@/components/ui/switch'
import { useParams } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'

const schema = z.object({
  overrideLogInMethods: z.boolean(),
  logInWithPasswordEnabled: z.boolean(),
  logInWithGoogleEnabled: z.boolean(),
  logInWithMicrosoftEnabled: z.boolean(),
  samlEnabled: z.boolean(),
  scimEnabled: z.boolean(),
})

export function EditOrganizationDetailsButton() {
  const { organizationId } = useParams()
  const { data: getOrganizationResponse, refetch: refetchOrganization } =
    useQuery(getOrganization, {
      id: organizationId,
    })
  const [open, setOpen] = useState(false)
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  })
  const updateOrganizationMutation = useMutation(updateOrganization)

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
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
        overrideLogInMethods: values.overrideLogInMethods,
        logInWithGoogleEnabled: values.logInWithGoogleEnabled,
        logInWithPasswordEnabled: values.logInWithPasswordEnabled,
        logInWithMicrosoftEnabled: values.logInWithMicrosoftEnabled,
        samlEnabled: values.samlEnabled,
        scimEnabled: values.scimEnabled,
      },
    })

    await refetchOrganization()
    setOpen(false)
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Organization Details</AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogDescription>
          Please update the organization details below.
        </AlertDialogDescription>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
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
                    If your customer wants to restrict their login methods to a
                    subset of the project-level ones, you can have their
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
                    Whether this organization can configure SAML.
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
                    Whether this organization can configure SCIM.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter>
              <Button type="submit">Save Changes</Button>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  )
}
