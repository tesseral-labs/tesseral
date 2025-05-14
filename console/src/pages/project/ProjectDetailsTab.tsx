import React, { useEffect, useMemo, useState } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  getProjectWebhookManagementURL,
  getVaultDomainSettings,
  updateProject,
  updateVaultDomainSettings,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
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
import { toast } from 'sonner';
import { Switch } from '@/components/ui/switch';
import { InputTags } from '@/components/input-tags';
import { EditProjectGoogleSettingsButton } from '@/pages/project/EditProjectGoogleSettingsButton';
import { EditProjectMicrosoftSettingsButton } from '@/pages/project/EditProjectMicrosoftSettingsButton';
import { EditProjectGithubSettingsButton } from './EditProjectGithubSettingsButton';

export const ProjectDetailsTab = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getProjectWebhookManagementUrlResponse } = useQuery(
    getProjectWebhookManagementURL,
  );

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Redirect Settings</CardTitle>
            <CardDescription>
              Where Tesseral will redirect your users whenever they return from
              your project's Vault. You can special-case where they go after
              logging in or signing up.
            </CardDescription>
          </div>
          <EditProjectRedirectURIsButton />
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Default Redirect URI</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.redirectUri || '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>After-Login Redirect URI</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.afterLoginRedirectUri || (
                    <span className="text-muted-foreground">
                      (Use Default Redirect URI)
                    </span>
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>After-Signup Redirect URI</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.afterSignupRedirectUri || (
                    <span className="text-muted-foreground">
                      (Use Default Redirect URI)
                    </span>
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Domains settings</CardTitle>
            <CardDescription>
              Settings related to domains that your project runs on.
            </CardDescription>
          </div>
          <EditProjectDomainSettingsButton />
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Cookie Domain</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.cookieDomain}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Trusted Domains</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.trustedDomains?.map(
                    (domain) => <div key={domain}>{domain}</div>,
                  )}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Login Methods</CardTitle>
            <CardDescription>
              Primary and secondary authentication methods your users can use.
              Organizations can take this list and restrict it further, but they
              can't add to it.
            </CardDescription>
          </div>
          <EditLoginMethodsButton />
        </CardHeader>
        <CardContent>
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Google settings</CardTitle>
            <CardDescription>
              Settings for "Log in with Google" in your project.
            </CardDescription>
          </div>
          <EditProjectGoogleSettingsButton />
        </CardHeader>
        <CardContent>
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
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Microsoft settings</CardTitle>
            <CardDescription>
              Settings for "Log in with Microsoft" in your project.
            </CardDescription>
          </div>
          <EditProjectMicrosoftSettingsButton />
        </CardHeader>
        <CardContent>
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>GitHub settings</CardTitle>
            <CardDescription>
              Settings for "Log in with GitHub" in your project.
            </CardDescription>
          </div>
          <EditProjectGithubSettingsButton />
        </CardHeader>
        <CardContent>
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Webhook settings</CardTitle>
            <CardDescription>
              Settings for webhooks sent from Tesseral to your application about
              your project.
            </CardDescription>
          </div>
          <a href={getProjectWebhookManagementUrlResponse?.url} target="_blank">
            <Button variant="outline">Edit</Button>
          </a>
        </CardHeader>
        <CardContent></CardContent>
      </Card>
    </div>
  );
};

const redirectURIsSchema = z.object({
  redirectUri: z.string().url(),
  afterLoginRedirectUri: z.string().url().or(z.literal('')).optional(),
  afterSignupRedirectUri: z.string().url().or(z.literal('')).optional(),
});

const EditProjectRedirectURIsButton = () => {
  const form = useForm<z.infer<typeof redirectURIsSchema>>({
    resolver: zodResolver(redirectURIsSchema),
    defaultValues: {
      redirectUri: '',
      afterLoginRedirectUri: '',
      afterSignupRedirectUri: '',
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        redirectUri: getProjectResponse?.project?.redirectUri,
        afterLoginRedirectUri:
          getProjectResponse?.project?.afterLoginRedirectUri,
        afterSignupRedirectUri:
          getProjectResponse?.project?.afterSignupRedirectUri,
      });
    }
  }, [getProjectResponse]);

  const updateProjectMutation = useMutation(updateProject);
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof redirectURIsSchema>) => {
    // if any of the redirect URIs have the scheme `http` or `https`,
    // automatically add those to the set of trusted domains.
    const trustedDomains = new Set(getProjectResponse!.project!.trustedDomains);
    if (
      values.redirectUri.startsWith('http://') ||
      values.redirectUri.startsWith('https://')
    ) {
      trustedDomains.add(new URL(values.redirectUri).host);
    }
    if (
      values.afterLoginRedirectUri?.startsWith('http://') ||
      values.afterLoginRedirectUri?.startsWith('https://')
    ) {
      trustedDomains.add(new URL(values.afterLoginRedirectUri).host);
    }
    if (
      values.afterSignupRedirectUri?.startsWith('http://') ||
      values.afterSignupRedirectUri?.startsWith('https://')
    ) {
      trustedDomains.add(new URL(values.afterSignupRedirectUri).host);
    }

    await updateProjectMutation.mutateAsync({
      project: {
        redirectUri: values.redirectUri,
        afterLoginRedirectUri: values.afterLoginRedirectUri,
        afterSignupRedirectUri: values.afterSignupRedirectUri,
        trustedDomains: Array.from(trustedDomains),
      },
    });
    await refetch();
    toast.success('Project Redirect URIs updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Redirect URIs</AlertDialogTitle>
          <AlertDialogDescription>
            Where Tesseral will redirect your users. You can special-case where
            they go after logging in or signing up.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="redirectUri"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Default Redirect URI</FormLabel>
                  <FormControl>
                    <Input placeholder="https://app.company.com/" {...field} />
                  </FormControl>
                  <FormDescription>
                    Where users will be redirected after visiting your project's
                    Vault.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="afterLoginRedirectUri"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>After-Login Redirect URI</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="https://app.company.com/after-login"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Where users will be redirected after logging in. If blank,
                    uses the default redirect URI.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="afterSignupRedirectUri"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>After-Signup Redirect URI</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="https://app.company.com/after-signup"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Where users will be redirected after signing up. If blank,
                    uses the default redirect URI.
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

const domainSettingsSchema = z.object({
  cookieDomain: z.string().regex(/^[^.]/, {
    message: "Cookie domain must not start with a dot ('.').",
  }),
  trustedDomains: z.array(
    z.string().regex(/^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+(:\d+)?$/),
  ),
});

const EditProjectDomainSettingsButton = () => {
  const form = useForm<z.infer<typeof domainSettingsSchema>>({
    resolver: zodResolver(domainSettingsSchema),
    defaultValues: {
      cookieDomain: '',
      trustedDomains: [],
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        cookieDomain: getProjectResponse?.project?.cookieDomain,
        trustedDomains: getProjectResponse?.project?.trustedDomains || [],
      });
    }
  }, [getProjectResponse]);

  const updateProjectMutation = useMutation(updateProject);
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof domainSettingsSchema>) => {
    if (
      !getProjectResponse?.project?.vaultDomain?.endsWith(values.cookieDomain)
    ) {
      form.setError('cookieDomain', {
        message: `Cookie Domain must be a parent domain of the Vault domain (${getProjectResponse?.project?.vaultDomain}).`,
      });
      return;
    }

    await updateProjectMutation.mutateAsync({
      project: {
        trustedDomains: values.trustedDomains,

        // only attempt to honor cookie domain if vault domain is custom
        ...(getProjectResponse?.project?.vaultDomainCustom
          ? { cookieDomain: values.cookieDomain }
          : {}),
      },
    });
    await refetch();
    toast.success('Trusted Domains updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Trusted Domains</AlertDialogTitle>
          <AlertDialogDescription>
            Edit the trusted domains for your project.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="cookieDomain"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Cookie Domain</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="app.company.com"
                      {...field}
                      disabled={!getProjectResponse?.project?.vaultDomainCustom}
                    />
                  </FormControl>
                  <FormDescription>
                    Client-side JavaScript on this domain and its subdomains
                    will have access to User access tokens. You cannot modify
                    this field until you have configured a custom Vault domain.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="trustedDomains"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Trusted Domains</FormLabel>
                  <FormControl>
                    <InputTags
                      placeholder="app.company.com, localhost:3000"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Add the domains that your app runs on, e.g.
                    "app.company.com". The Vault domain is always a trusted
                    domain.
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
