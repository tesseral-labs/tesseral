import {
  AlertDialogContent,
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
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
import { Switch } from '@/components/ui/switch';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

const schema = z.object({
  apiKeysEnabled: z.boolean(),
  apiKeysPrefix: z
    .string()
    .max(64, {
      message: 'Prefix must be at most 64 characters long',
    })
    .regex(/^[a-z0-9_]+$/, {
      message:
        'Prefix can only contain lowercase letters, numbers, and underscores',
    }),
});

export function EditAPIKeySettingsButton() {
  const [open, setOpen] = useState(false);

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      apiKeysEnabled: false,
      apiKeysPrefix: '',
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    await updateProjectMutation.mutateAsync({
      project: {
        apiKeysEnabled: data.apiKeysEnabled,
        apiKeySecretTokenPrefix: data.apiKeysPrefix,
      },
    });

    await refetch();

    toast.success('API Key settings updated successfully');
    setOpen(false);
  };

  useEffect(() => {
    if (getProjectResponse && getProjectResponse.project) {
      form.setValue(
        'apiKeysEnabled',
        getProjectResponse.project.apiKeysEnabled || false,
      );
      form.setValue(
        'apiKeysPrefix',
        getProjectResponse.project.apiKeySecretTokenPrefix || '',
      );
    }
  }, [getProjectResponse]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit API Key Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Update the settings for your API keys.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="apiKeysEnabled"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>API Keys Enabled</FormLabel>
                    <FormControl className="block">
                      <Switch
                        id="apiKeysEnabled"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                    <FormDescription>
                      Whether or not Organizations are allowed to create API
                      Keys.
                    </FormDescription>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="apiKeysPrefix"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>API Keys Prefix</FormLabel>
                    <FormControl>
                      <Input type="text" id="apiKeysPrefix" {...field} />
                    </FormControl>
                    <FormMessage />
                    <FormDescription>
                      Set a prefix for your API keys. We recommend ending your
                      prefix with an underscore. For example:{' '}
                      <span className="bg-muted text-muted-foreground font-mono text-xs p-1">
                        myapp_sk_
                      </span>
                      .
                    </FormDescription>
                  </FormItem>
                )}
              />
            </div>

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
