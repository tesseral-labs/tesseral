import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  getProjectUISettings,
  updateProject,
  updateProjectUISettings,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
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
import { Switch } from '@/components/ui/switch';

const schema = z.object({
  autoCreateOrganizations: z.boolean(),
});

export function EditBehaviorSettingsButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      autoCreateOrganizations: false,
    },
  });

  const { data: getProjectUISettingsResponse, refetch } =
    useQuery(getProjectUISettings);
  useEffect(() => {
    if (getProjectUISettingsResponse?.projectUiSettings) {
      form.reset({
        autoCreateOrganizations:
          getProjectUISettingsResponse.projectUiSettings
            .autoCreateOrganizations,
      });
    }
  }, [getProjectUISettingsResponse]);

  const { mutateAsync: updateProjectUISettingsAsync } = useMutation(
    updateProjectUISettings,
  );
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    await updateProjectUISettingsAsync({
      autoCreateOrganizations: values.autoCreateOrganizations,
    });
    await refetch();
    toast.success('Behavior settings updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Behavior Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Modify the behavior of the login flow your Users will see.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="autoCreateOrganizations"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Auto-create organizations</FormLabel>
                  <FormDescription>
                    Instead of asking Users to create an Organization on signup,
                    auto-create one for them.
                  </FormDescription>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
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
}
