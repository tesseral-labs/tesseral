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
  logInWithGithub: z.boolean(),
  githubOAuthClientId: z.string(),
  githubOAuthClientSecret: z.string(),
});

export function EditProjectGithubSettingsButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGithub: true,
      githubOAuthClientId: '',
      githubOAuthClientSecret: '',
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithGithub: getProjectResponse.project.logInWithGithub,
        githubOAuthClientId: getProjectResponse.project.githubOauthClientId,
        githubOAuthClientSecret: '',
      });
    }
  }, [getProjectResponse]);

  const { mutateAsync: updateProjectMutationAsync } =
    useMutation(updateProject);
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    if (!values.logInWithGithub) {
      if (
        !getProjectResponse?.project?.logInWithEmail &&
        !getProjectResponse?.project?.logInWithGoogle &&
        !getProjectResponse?.project?.logInWithMicrosoft
      ) {
        form.setError('logInWithGithub', {
          message:
            'At least one of Log in with Email, Log in with Microsoft, Log in with Google, or Log in with GitHub must be enabled.',
        });
        return;
      }
    }

    await updateProjectMutationAsync({
      project: {
        logInWithGithub: values.logInWithGithub,
        githubOauthClientId: values.githubOAuthClientId,
        githubOauthClientSecret: values.githubOAuthClientSecret,
      },
    });
    await refetch();
    toast.success('GitHub settings updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit GitHub Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Edit the settings for "Log in with GitHub" in your project.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="logInWithGithub"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with GitHub</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can log in using their GitHub account.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="githubOAuthClientId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>GitHub OAuth Client ID</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    Your company's GitHub OAuth Client ID.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="githubOAuthClientSecret"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>GitHub OAuth Client Secret</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    Your company's GitHub OAuth Client Secret.
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
