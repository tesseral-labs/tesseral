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
import Loader from '@/components/ui/loader';
import { toast } from 'sonner';
import { Switch } from '@/components/ui/switch';

export const ProjectDetailsTab = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});

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
        <CardHeader>
          <CardTitle>Authentication settings</CardTitle>
          <CardDescription>
            Configure the login methods your customers can use to log in to your
            application.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Password</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithPassword
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Google</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithGoogle
                    ? 'Enabled'
                    : 'Disabled'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>Log in with Microsoft</DetailsGridKey>
                <DetailsGridValue>
                  {getProjectResponse?.project?.logInWithMicrosoft
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
        <CardHeader>
          <CardTitle>
            <div className="grid grid-cols-2 gap-8">
              <span>Google settings</span>
              <div className="text-right">
                <Link to="/project-settings/log-in-with-google/edit">
                  <Button variant="outline" size="sm">
                    Edit
                  </Button>
                </Link>
              </div>
            </div>
          </CardTitle>
          <CardDescription>
            Settings for "Log in with Google" in your project.
          </CardDescription>
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
        <CardHeader>
          <CardTitle>
            <div className="grid grid-cols-2 gap-8">
              <span>Microsoft settings</span>
              <div className="text-right">
                <Link to="/project-settings/log-in-with-microsoft/edit">
                  <Button variant="outline" size="sm">
                    Edit
                  </Button>
                </Link>
              </div>
            </div>
          </CardTitle>
          <CardDescription>
            Settings for "Log in with Microsoft" in your project.
          </CardDescription>
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
            <CardTitle>SAML Settings</CardTitle>
            <CardDescription>
              Settings for "Log in with SAML" in your Project.
            </CardDescription>
          </div>
          <EditProjectSAMLSettingsButton />
        </CardHeader>
        <CardContent>
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Status</DetailsGridKey>
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
    await updateProjectMutation.mutateAsync({
      project: {
        redirectUri: values.redirectUri,
        afterLoginRedirectUri: values.afterLoginRedirectUri,
        afterSignupRedirectUri: values.afterSignupRedirectUri,
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

const samlSchema = z.object({
  logInWithSaml: z.boolean(),
});

const EditProjectSAMLSettingsButton = () => {
  const form = useForm<z.infer<typeof samlSchema>>({
    resolver: zodResolver(samlSchema),
    defaultValues: {
      logInWithSaml: false,
    },
  });

  const { data: getProjectResponse, refetch } = useQuery(getProject);
  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithSaml: getProjectResponse?.project?.logInWithSaml || false,
      });
    }
  }, [getProjectResponse]);

  const updateProjectMutation = useMutation(updateProject);
  const [open, setOpen] = useState(false);
  const handleSubmit = async (values: z.infer<typeof samlSchema>) => {
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithSaml: values.logInWithSaml,
      },
    });
    await refetch();
    toast.success('SAML Settings updated');
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit SAML Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Enable or disable "Log in with SAML" for your project.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="logInWithSaml"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Status</FormLabel>
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
