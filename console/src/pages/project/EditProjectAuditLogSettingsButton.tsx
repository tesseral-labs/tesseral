import {
  AlertDialog,
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
import Loader from '@/components/ui/loader';
import { Switch } from '@/components/ui/switch';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { AlertDialogCancel } from '@radix-ui/react-alert-dialog';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const schema = z.object({
  auditLogsEnabled: z.boolean(),
});

export function EditProjectAuditLogSettingsButton() {
  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      auditLogsEnabled: true,
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        auditLogsEnabled: data.auditLogsEnabled,
      },
    });

    await refetch();
    setOpen(false);
    toast.success('Audit Log settings updated successfully.');
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Audit Log Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Edit the settings for how Audit Logs are displayed for your
            customers.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="auditLogsEnabled"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Audit logs enabled</FormLabel>
                  <FormDescription>
                    When enabled, your customers will be able to view audit log
                    events in the Vault.
                  </FormDescription>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={(checked) => field.onChange(checked)}
                      disabled={updateProjectMutation.isPending}
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit" disabled={updateProjectMutation.isPending}>
                {updateProjectMutation.isPending && <Loader />}
                {updateProjectMutation.isPending ? 'Saving...' : 'Save'}
              </Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
