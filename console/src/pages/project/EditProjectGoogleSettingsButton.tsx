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
  logInWithGoogle: z.boolean(),
  googleOAuthClientId: z.string(),
  googleOAuthClientSecret: z.string(),
});

export function EditProjectGoogleSettingsButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGoogle: true,
      googleOAuthClientId: '',
      googleOAuthClientSecret: '',
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithGoogle: getProjectResponse.project.logInWithGoogle,
        googleOAuthClientId: getProjectResponse.project.googleOauthClientId,
        googleOAuthClientSecret: '',
      });
    }
  }, [getProjectResponse]);

  const { mutateAsync: updateProjectMutationAsync } =
    useMutation(updateProject);
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    if (!values.logInWithGoogle) {
      if (
        !getProjectResponse?.project?.logInWithEmail &&
        !getProjectResponse?.project?.logInWithMicrosoft &&
        !getProjectResponse?.project?.logInWithGithub
      ) {
        form.setError('logInWithGoogle', {
          message:
            'At least one of Log in with Email, Log in with Google, Log in with Microsoft, or Log in with GitHub must be enabled.',
        });
        return;
      }
    }

    await updateProjectMutationAsync({
      project: {
        logInWithGoogle: values.logInWithGoogle,
        googleOauthClientId: values.googleOAuthClientId,
        googleOauthClientSecret: values.googleOAuthClientSecret,
      },
    });
    await refetch();
    toast.success('Google settings updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Google Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Edit the settings for "Log in with Google" in your project.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="logInWithGoogle"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Google</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can log in using their Google account.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="googleOAuthClientId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Google OAuth Client ID</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    Your company's Google OAuth Client ID.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="googleOAuthClientSecret"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Google OAuth Client Secret</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    Your company's Google OAuth Client Secret.
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
}
