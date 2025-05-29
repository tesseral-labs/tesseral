import React from 'react';

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
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
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
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { EditProjectGoogleSettingsButton } from './EditProjectGoogleSettingsButton';
import { EditProjectMicrosoftSettingsButton } from './EditProjectMicrosoftSettingsButton';
import { EditProjectGithubSettingsButton } from './EditProjectGithubSettingsButton';

export function LoginSettingsTab() {
  const { data: getProjectResponse } = useQuery(getProject, {});

  return (
    <div className="space-y-8">
      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Login Methods</ConsoleCardTitle>
            <ConsoleCardDescription>
              Primary and secondary authentication methods your users can use.
              Organizations can take this list and restrict it further, but they
              can't add to it.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditLoginMethodsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Email (Magic Links)</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithEmail
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Password</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithPassword
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>
                  Log in with Passkey (Secondary Factor)
                </DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithPasskey
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>
                  Log in with Authenticator App (Secondary Factor)
                </DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithAuthenticatorApp
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with SAML</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithSaml
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Google settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings for "Log in with Google" in your project.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditProjectGoogleSettingsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithGoogle
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google OAuth Client ID</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.googleOauthClientId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Google OAuth Client Secret</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.googleOauthClientId ? (
                    <div className="text-muted-foreground">Encrypted</div>
                  ) : (
                    '-'
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>
      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Microsoft settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings for "Log in with Microsoft" in your project.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditProjectMicrosoftSettingsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithMicrosoft
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft OAuth Client ID</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.microsoftOauthClientId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Microsoft OAuth Client Secret</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.microsoftOauthClientId ? (
                    <div className="text-muted-foreground">Encrypted</div>
                  ) : (
                    '-'
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>GitHub settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings for "Log in with GitHub" in your project.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditProjectGithubSettingsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithGithub
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>GitHub OAuth Client ID</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.githubOauthClientId || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>GitHub OAuth Client Secret</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.githubOauthClientId ? (
                    <div className="text-muted-foreground">Encrypted</div>
                  ) : (
                    '-'
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </ConsoleCardContent>
      </ConsoleCard>
    </div>
  );
}

const loginMethodsSchema = z.object({
  logInWithEmail: z.boolean(),
  logInWithPassword: z.boolean(),
  logInWithPasskey: z.boolean(),
  logInWithAuthenticatorApp: z.boolean(),
  logInWithSaml: z.boolean(),
});

const EditLoginMethodsButton = () => {
  const form = useForm<z.infer<typeof loginMethodsSchema>>({
    resolver: zodResolver(loginMethodsSchema),
    defaultValues: {
      logInWithEmail: false,
      logInWithPassword: false,
      logInWithPasskey: false,
      logInWithAuthenticatorApp: false,
      logInWithSaml: false,
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);

  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithEmail: getProjectResponse?.project?.logInWithEmail || false,
        logInWithPassword:
          getProjectResponse?.project?.logInWithPassword || false,
        logInWithPasskey:
          getProjectResponse?.project?.logInWithPasskey || false,
        logInWithAuthenticatorApp:
          getProjectResponse?.project?.logInWithAuthenticatorApp || false,
        logInWithSaml: getProjectResponse?.project?.logInWithSaml || false,
      });
    }
  }, [getProjectResponse]);

  const updateProjectMutation = useMutation(updateProject);
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof loginMethodsSchema>) => {
    if (!values.logInWithEmail) {
      if (
        !getProjectResponse?.project?.logInWithGoogle &&
        !getProjectResponse?.project?.logInWithMicrosoft
      ) {
        form.setError('logInWithEmail', {
          message:
            'At least one of Log in with Email, Log in with Google, or Log in with Microsoft must be enabled.',
        });
        return;
      }
    }

    await updateProjectMutation.mutateAsync({
      project: {
        ...values,
      },
    });
    await refetch();
    toast.success('Login methods updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Login Methods</AlertDialogTitle>
          <AlertDialogDescription>
            To enable Google or Microsoft, go to their respective settings
            section.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="logInWithEmail"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Email</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can log in with a Magic Link sent to their
                    email address.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="logInWithPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Password</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can log in using a password.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="logInWithPasskey"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Passkey</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can register a passkey as a secondary
                    authentication factor.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="logInWithAuthenticatorApp"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with Authenticator App</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Users can register an authenticator app as a
                    secondary authentication factor.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="logInWithSaml"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Log in with SAML</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether Organizations in this Project can enable SAML.
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
