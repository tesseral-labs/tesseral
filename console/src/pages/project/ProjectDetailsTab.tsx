import React, { useEffect, useState } from 'react';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import { Button } from '@/components/ui/button';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createStripeCheckoutLink,
  getProject,
  getProjectEntitlements,
  getProjectWebhookManagementURL,
  updateProject,
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
import { InputTags } from '@/components/input-tags';
import { EditAPIKeySettingsButton } from './project-ui-settings/EditAPIKeySettingsButton';
import { ConsoleCardDetails } from '@/components/ui/console-card';

export const ProjectDetailsTab = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getProjectWebhookManagementUrlResponse } = useQuery(
    getProjectWebhookManagementURL,
  );
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );

  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }

  return (
    <div className="space-y-8">
      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Redirect Settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Where Tesseral will redirect your users whenever they return from
              your project's Vault. You can special-case where they go after
              logging in or signing up.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditProjectRedirectURIsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
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
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Domains settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings related to domains that your project runs on.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <EditProjectDomainSettingsButton />
        </ConsoleCardHeader>
        <ConsoleCardContent>
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
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>API key settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings for API keys used by your customers with your product.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          {getProjectEntitlementsResponse?.entitledBackendApiKeys && (
            <EditAPIKeySettingsButton />
          )}
        </ConsoleCardHeader>
        <ConsoleCardContent>
          {!getProjectEntitlementsResponse?.entitledBackendApiKeys ? (
            <div className="text-sm my-8 w-full flex flex-col items-center justify-center space-y-6">
              <div className="font-medium">
                API Keys are available on the Growth Tier.
              </div>

              <div className="flex items-center gap-x-4">
                <Button onClick={handleUpgrade}>Upgrade to Growth Tier</Button>
                <span>
                  or{' '}
                  <a
                    href="https://cal.com/ned-o-leary-j8ydyi/30min"
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                  >
                    meet an expert
                  </a>
                  .
                </span>
              </div>
            </div>
          ) : (
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Status</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.apiKeysEnabled
                      ? 'Enabled'
                      : 'Disabled'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>API Key Prefix</DetailsGridKey>
                  <DetailsGridValue>
                    {getProjectResponse?.project?.apiKeySecretTokenPrefix ||
                      '-'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          )}
        </ConsoleCardContent>
      </ConsoleCard>

      <ConsoleCard>
        <ConsoleCardHeader>
          <ConsoleCardDetails>
            <ConsoleCardTitle>Webhook settings</ConsoleCardTitle>
            <ConsoleCardDescription>
              Settings for webhooks sent from Tesseral to your application about
              your project.
            </ConsoleCardDescription>
          </ConsoleCardDetails>
          <a href={getProjectWebhookManagementUrlResponse?.url} target="_blank">
            <Button variant="outline">Edit</Button>
          </a>
        </ConsoleCardHeader>
        <ConsoleCardContent></ConsoleCardContent>
      </ConsoleCard>
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
                  <FormDescription>
                    Where users will be redirected after visiting your project's
                    Vault.
                  </FormDescription>
                  <FormControl>
                    <Input placeholder="https://app.company.com/" {...field} />
                  </FormControl>

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
                  <FormDescription>
                    Where users will be redirected after logging in. If blank,
                    uses the default redirect URI.
                  </FormDescription>
                  <FormControl>
                    <Input
                      placeholder="https://app.company.com/after-login"
                      {...field}
                    />
                  </FormControl>

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
                  <FormDescription>
                    Where users will be redirected after signing up. If blank,
                    uses the default redirect URI.
                  </FormDescription>
                  <FormControl>
                    <Input
                      placeholder="https://app.company.com/after-signup"
                      {...field}
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
                  <FormDescription>
                    Client-side JavaScript on this domain and its subdomains
                    will have access to User access tokens. You cannot modify
                    this field until you have configured a custom Vault domain.
                  </FormDescription>
                  <FormControl>
                    <Input
                      placeholder="app.company.com"
                      {...field}
                      disabled={!getProjectResponse?.project?.vaultDomainCustom}
                    />
                  </FormControl>

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
