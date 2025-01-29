import { z } from 'zod'
import { useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  getOrganizationGoogleHostedDomains,
  updateOrganizationGoogleHostedDomains,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import React, { useEffect, useState } from 'react'
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
import { InputTags } from '@/components/input-tags'

const schema = z.object({
  googleHostedDomains: z.array(z.string()),
})

export function EditOrganizationGoogleConfigurationButton() {
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
                    <InputTags
                      className="max-w-96"
                      placeholder="example.com"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    The set of Google workspaces associated with this
                    organization, identified by their primary domain.
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
