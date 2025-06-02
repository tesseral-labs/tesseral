import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganization,
  getOrganizationDomains,
  getProject,
  updateOrganization,
  updateOrganizationDomains,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React, { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
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
import { Link } from 'react-router-dom';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import { Input } from '@/components/ui/input';
import { PageContent, PageHeader, PageTitle } from '@/components/page';
import { toast } from 'sonner';
import { InputTags } from '@/components/input-tags';

const schema = z.object({
  apiKeysEnabled: z.boolean(),
  displayName: z.string(),
  logInWithGoogle: z.boolean(),
  logInWithMicrosoft: z.boolean(),
  logInWithGithub: z.boolean(),
  logInWithEmail: z.boolean(),
  logInWithPassword: z.boolean(),
  logInWithSaml: z.boolean(),
  logInWithAuthenticatorApp: z.boolean(),
  logInWithPasskey: z.boolean(),
  requireMfa: z.boolean(),
  scimEnabled: z.boolean(),
  domains: z.array(z.string()),
});

export const EditOrganizationPage = () => {
  const navigate = useNavigate();
  const { organizationId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getOrganizationDomainsResponse } = useQuery(
    getOrganizationDomains,
    {
      organizationId,
    },
  );

  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  });

  const updateOrganizationMutation = useMutation(updateOrganization);
  const updateOrganizationDomainsMutation = useMutation(
    updateOrganizationDomains,
  );

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      // Currently there's an issue with the types of react-hook-form and zod
      // preventing the compiler from inferring the correct types.
      form.reset({
        ...form.getValues(),
        apiKeysEnabled:
          getOrganizationResponse.organization.apiKeysEnabled || false,
        displayName: getOrganizationResponse.organization.displayName,
        logInWithGoogle: getOrganizationResponse.organization.logInWithGoogle,
        logInWithMicrosoft:
          getOrganizationResponse.organization.logInWithMicrosoft,
        logInWithGithub: getOrganizationResponse.organization.logInWithGithub,
        logInWithEmail: getOrganizationResponse.organization.logInWithEmail,
        logInWithPassword:
          getOrganizationResponse.organization.logInWithPassword,
        logInWithSaml: getOrganizationResponse.organization.logInWithSaml,
        logInWithAuthenticatorApp:
          getOrganizationResponse.organization.logInWithAuthenticatorApp,
        logInWithPasskey: getOrganizationResponse.organization.logInWithPasskey,
        requireMfa: getOrganizationResponse.organization.requireMfa,
        scimEnabled: getOrganizationResponse.organization.scimEnabled,
      });
    }

    if (getOrganizationDomainsResponse?.organizationDomains) {
      form.reset({
        ...form.getValues(),
        domains: getOrganizationDomainsResponse.organizationDomains.domains,
      });
    }
  }, [getOrganizationResponse, getOrganizationDomainsResponse]);

  const onSubmit = async (values: z.infer<typeof schema>) => {
    if (values.requireMfa) {
      if (!values.logInWithAuthenticatorApp && !values.logInWithPasskey) {
        form.setError('requireMfa', {
          message:
            'To require MFA, you must enable either Log in with Authenticator App or Log in with Passkey.',
        });
        return;
      }
    }

    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        apiKeysEnabled: values.apiKeysEnabled,
        displayName: values.displayName,
        logInWithGoogle: values.logInWithGoogle,
        logInWithMicrosoft: values.logInWithMicrosoft,
        logInWithGithub: values.logInWithGithub,
        logInWithEmail: values.logInWithEmail,
        logInWithPassword: values.logInWithPassword,
        logInWithSaml: values.logInWithSaml,
        logInWithAuthenticatorApp: values.logInWithAuthenticatorApp,
        logInWithPasskey: values.logInWithPasskey,
        requireMfa: values.requireMfa,
        scimEnabled: values.scimEnabled,
      },
    });

    await updateOrganizationDomainsMutation.mutateAsync({
      organizationId,
      organizationDomains: {
        domains: values.domains,
      },
    });

    toast.success('Organization updated successfully');
    navigate(`/organizations/${organizationId}`);
  };

  return (
    <>
      <PageHeader>
        <PageTitle>
          Edit {getOrganizationResponse?.organization?.displayName}
        </PageTitle>
      </PageHeader>

      <PageContent>
        <Form {...form}>
          {}
          {/** There's an issue with the types of react-hook-form and zod
        preventing the compiler from inferring the correct types.*/}
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="mt-8 space-y-8"
          >
            {/** eslint-enable @typescript-eslint/no-unsafe-call */}
            <ConsoleCard>
              <ConsoleCardHeader>
                <ConsoleCardDetails>
                  <ConsoleCardTitle>Organization settings</ConsoleCardTitle>
                  <ConsoleCardDescription>
                    Configure basic settings on this organization.
                  </ConsoleCardDescription>
                </ConsoleCardDetails>
              </ConsoleCardHeader>
              <ConsoleCardContent className="space-y-8">
                <FormField
                  control={form.control}
                  name="displayName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Display Name</FormLabel>
                      <FormDescription>
                        A human-friendly name for the organization.
                      </FormDescription>
                      <FormControl>
                        <Input
                          className="max-w-80"
                          placeholder="Acme Corporation"
                          {...field}
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
              </ConsoleCardContent>
            </ConsoleCard>

            <ConsoleCard>
              <ConsoleCardHeader>
                <ConsoleCardDetails>
                  <ConsoleCardTitle>Login settings</ConsoleCardTitle>
                  <ConsoleCardDescription>
                    Configure how users can log into this organization.
                  </ConsoleCardDescription>
                </ConsoleCardDetails>
              </ConsoleCardHeader>
              <ConsoleCardContent className="space-y-8">
                {getProjectResponse?.project?.logInWithGoogle && (
                  <FormField
                    control={form.control}
                    name="logInWithGoogle"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Google</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using
                          their Google account.
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
                )}

                {getProjectResponse?.project?.logInWithMicrosoft && (
                  <FormField
                    control={form.control}
                    name="logInWithMicrosoft"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Microsoft</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using
                          their Microsoft account.
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
                )}

                {getProjectResponse?.project?.logInWithGithub && (
                  <FormField
                    control={form.control}
                    name="logInWithGithub"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with GitHub</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using
                          their GitHub account.
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
                )}

                {getProjectResponse?.project?.logInWithEmail && (
                  <FormField
                    control={form.control}
                    name="logInWithEmail"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Email</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in an
                          email.
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
                )}

                {getProjectResponse?.project?.logInWithPassword && (
                  <FormField
                    control={form.control}
                    name="logInWithPassword"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Password</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using a
                          password.
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
                )}

                {getProjectResponse?.project?.logInWithAuthenticatorApp && (
                  <FormField
                    control={form.control}
                    name="logInWithAuthenticatorApp"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Authenticator App</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using an
                          Authenticator App as a secondary factor.
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
                )}

                {getProjectResponse?.project?.logInWithPasskey && (
                  <FormField
                    control={form.control}
                    name="logInWithPasskey"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with Passkey</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in using a
                          Passkey as a secondary factor.
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
                )}

                <FormField
                  control={form.control}
                  name="requireMfa"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>Require MFA</FormLabel>
                      <FormDescription>
                        Whether Users in this Organization must authenticate
                        with a secondary factor when logging in.
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
              </ConsoleCardContent>
            </ConsoleCard>
            {getProjectResponse?.project?.apiKeysEnabled && (
              <ConsoleCard>
                <ConsoleCardHeader>
                  <ConsoleCardDetails>
                    <ConsoleCardTitle>API Keys settings</ConsoleCardTitle>
                    <ConsoleCardDescription>
                      Configure whether this Organization can use API Keys.
                    </ConsoleCardDescription>
                  </ConsoleCardDetails>
                </ConsoleCardHeader>
                <ConsoleCardContent className="space-y-8">
                  <FormField
                    control={form.control}
                    name="apiKeysEnabled"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>API Keys Enabled</FormLabel>
                        <FormDescription>
                          Whether this Organization can authenticate to your
                          service using API Keys.
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
                </ConsoleCardContent>
              </ConsoleCard>
            )}
            <ConsoleCard>
              <ConsoleCardHeader>
                <ConsoleCardDetails>
                  <ConsoleCardTitle>Enterprise settings</ConsoleCardTitle>
                  <ConsoleCardDescription>
                    Configure whether this Organization can use SCIM.
                  </ConsoleCardDescription>
                </ConsoleCardDetails>
              </ConsoleCardHeader>
              <ConsoleCardContent className="space-y-8">
                {getProjectResponse?.project?.logInWithSaml && (
                  <FormField
                    control={form.control}
                    name="logInWithSaml"
                    render={({ field }: { field: any }) => (
                      <FormItem>
                        <FormLabel>Log in with SAML</FormLabel>
                        <FormDescription>
                          Whether this organization can configure SAML
                          Connections and use them to log in with SAML.
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
                )}

                <FormField
                  control={form.control}
                  name="scimEnabled"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>SCIM Enabled</FormLabel>
                      <FormDescription>
                        Whether this Organization can configure SCIM
                        ("Enterprise Directory Sync").
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
                  name="domains"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>SAML / SCIM Domains</FormLabel>
                      <FormDescription>
                        SAML and SCIM users must have emails from this list of
                        domains.
                      </FormDescription>
                      <FormControl>
                        <InputTags
                          className="max-w-96"
                          placeholder="example.com"
                          {...field}
                          value={field.value || []}
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
              </ConsoleCardContent>
            </ConsoleCard>

            <div className="flex justify-end gap-x-4 pb-8">
              <Button variant="outline" asChild>
                <Link to={`/organizations/${organizationId}`}>Cancel</Link>
              </Button>
              <Button type="submit">Save Changes</Button>
            </div>
          </form>
        </Form>
      </PageContent>
    </>
  );
};
