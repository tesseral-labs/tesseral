import { z } from 'zod';
import { useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganizationMicrosoftTenantIDs,
  updateOrganizationMicrosoftTenantIDs,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery';
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
  microsoftTenantIds: z.array(z.string()),
});

export const EditOrganizationMicrosoftConfigurationButton = () => {
  const { organizationId } = useParams();
  const { data: getOrganizationMicrosoftTenantIdsResponse, refetch } = useQuery(
    getOrganizationMicrosoftTenantIDs,
    {
      organizationId,
    },
  );
  const updateOrganizationMicrosoftTenantIdsMutation = useMutation(
    updateOrganizationMicrosoftTenantIDs,
  );
  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      microsoftTenantIds: [],
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */
  useEffect(() => {
    if (
      getOrganizationMicrosoftTenantIdsResponse?.organizationMicrosoftTenantIds
    ) {
      /* eslint-disable @typescript-eslint/no-unsafe-call */
      // Currently there's an issue with the types of react-hook-form and zod
      // preventing the compiler from inferring the correct types.
      form.reset({
        microsoftTenantIds:
          getOrganizationMicrosoftTenantIdsResponse
            .organizationMicrosoftTenantIds.microsoftTenantIds,
      });
      /* eslint-enable @typescript-eslint/no-unsafe-call */
    }
  }, [getOrganizationMicrosoftTenantIdsResponse]);

  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateOrganizationMicrosoftTenantIdsMutation.mutateAsync({
      organizationId,
      organizationMicrosoftTenantIds: {
        microsoftTenantIds: values.microsoftTenantIds,
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
          <AlertDialogTitle>Edit Microsoft configuration</AlertDialogTitle>
          <AlertDialogDescription>
            Edit organization Microsoft configuration.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          {/* eslint-disable @typescript-eslint/no-unsafe-call */}
          {/** Currently there's an issue with the types of react-hook-form and zod
          preventing the compiler from inferring the correct types.*/}
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            {/* eslint-enable @typescript-eslint/no-unsafe-call */}
            <FormField
              control={form.control}
              name="microsoftTenantIds"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Microsoft Tenant IDs</FormLabel>
                  <FormControl>
                    <InputTags
                      className="max-w-96"
                      placeholder="3b465a84-801e-..."
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    The set of Microsoft Tenant IDs associated with this
                    organization.
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
  );
};
