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
import { Input } from '@/components/ui/input';
import { InputTags } from '@/components/input-tags';
import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { toast } from 'sonner';
import { Switch } from '@/components/ui/switch';

const schema = z.object({
  logInWithMicrosoft: z.boolean(),
  microsoftOAuthClientId: z.string(),
  microsoftOAuthClientSecret: z.string(),
});

export function EditProjectMicrosoftSettingsButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithMicrosoft: true,
      microsoftOAuthClientId: '',
      microsoftOAuthClientSecret: '',
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithMicrosoft: getProjectResponse.project.logInWithMicrosoft,
        microsoftOAuthClientId:
          getProjectResponse.project.microsoftOauthClientId,
        microsoftOAuthClientSecret: '',
      });
    }
  }, [getProjectResponse]);

  const { mutateAsync: updateProjectMutationAsync } =
    useMutation(updateProject);
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    if (!values.logInWithMicrosoft) {
      if (
        !getProjectResponse?.project?.logInWithEmail &&
        !getProjectResponse?.project?.logInWithGoogle &&
        !getProjectResponse?.project?.logInWithGithub
      ) {
        form.setError('logInWithMicrosoft', {
          message:
            'At least one of Log in with Email, Log in with Microsoft, Log in with Google, or Log in with GitHub must be enabled.',
        });
        return;
      }
    }

    await updateProjectMutationAsync({
      project: {
        logInWithMicrosoft: values.logInWithMicrosoft,
        microsoftOauthClientId: values.microsoftOAuthClientId,
        microsoftOauthClientSecret: values.microsoftOAuthClientSecret,
      },
    });
    await refetch();
    toast.success('Microsoft settings updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Microsoft Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Edit the settings for "Log in with Microsoft" in your project.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="logInWithMicrosoft"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Microsoft</FormLabel>
                  <FormDescription>
                    Whether Users can log in using their Microsoft account.
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

            <FormField
              control={form.control}
              name="microsoftOAuthClientId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Microsoft OAuth Client ID</FormLabel>
                  <FormDescription>
                    Your company's Microsoft OAuth Client ID.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="microsoftOAuthClientSecret"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Microsoft OAuth Client Secret</FormLabel>
                  <FormDescription>
                    Your company's Microsoft OAuth Client Secret.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} />
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
