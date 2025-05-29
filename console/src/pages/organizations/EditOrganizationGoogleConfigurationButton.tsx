import { z } from 'zod';
import { useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganizationGoogleHostedDomains,
  updateOrganizationGoogleHostedDomains,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect, useState } from 'react';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { InputTags } from '@/components/input-tags';

const schema = z.object({
  googleHostedDomains: z.array(z.string()),
});

export const EditOrganizationGoogleConfigurationButton = () => {
  const { organizationId } = useParams();
  const { data: getOrganizationGoogleHostedDomainsResponse, refetch } =
    useQuery(getOrganizationGoogleHostedDomains, {
      organizationId,
    });
  const updateOrganizationGoogleHostedDomainsMutation = useMutation(
    updateOrganizationGoogleHostedDomains,
  );
  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      googleHostedDomains: [],
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */
  useEffect(() => {
    if (
      getOrganizationGoogleHostedDomainsResponse?.organizationGoogleHostedDomains
    ) {
      /* eslint-disable @typescript-eslint/no-unsafe-call */
      // Currently there's an issue with the types of react-hook-form and zod
      // preventing the compiler from inferring the correct types.
      form.reset({
        googleHostedDomains:
          getOrganizationGoogleHostedDomainsResponse
            .organizationGoogleHostedDomains.googleHostedDomains,
      });
      /* eslint-enable @typescript-eslint/no-unsafe-call */
    }
  }, [getOrganizationGoogleHostedDomainsResponse]);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateOrganizationGoogleHostedDomainsMutation.mutateAsync({
      organizationId,
      organizationGoogleHostedDomains: {
        googleHostedDomains: values.googleHostedDomains,
      },
    });
    await refetch();
    setOpen(false);
  };

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
          {/* eslint-disable @typescript-eslint/no-unsafe-call */}
          {/* Currently there's an issue with the types of react-hook-form and zod
          preventing the compiler from inferring the correct types.*/}
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            {/* eslint-enable @typescript-eslint/no-unsafe-call */}
            <FormField
              control={form.control}
              name="googleHostedDomains"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Google Hosted Domains</FormLabel>
                  <FormDescription>
                    The set of Google workspaces associated with this
                    organization, identified by their primary domain.
                  </FormDescription>
                  <FormControl>
                    <InputTags
                      className="max-w-96"
                      placeholder="example.com"
                      {...field}
                    />
                  </FormControl>
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
  );
};
