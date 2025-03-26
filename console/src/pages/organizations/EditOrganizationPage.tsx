import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganization, getOrganizationDomains,
  getProject,
  updateOrganization, updateOrganizationDomains,
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
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { PageTitle } from '@/components/page';
import { toast } from 'sonner';
import { InputTags } from '@/components/input-tags';

const schema = z.object({
  displayName: z.string(),
  logInWithGoogle: z.boolean(),
  logInWithMicrosoft: z.boolean(),
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
  const { data: getOrganizationDomainsResponse } = useQuery(getOrganizationDomains, {
    organizationId,
  })
   
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  });
   
  const updateOrganizationMutation = useMutation(updateOrganization);
  const updateOrganizationDomainsMutation = useMutation(updateOrganizationDomains);

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      // Currently there's an issue with the types of react-hook-form and zod
      // preventing the compiler from inferring the correct types.
      form.reset({
        ...form.getValues(),
        displayName: getOrganizationResponse.organization.displayName,
        logInWithGoogle: getOrganizationResponse.organization.logInWithGoogle,
        logInWithMicrosoft:
          getOrganizationResponse.organization.logInWithMicrosoft,
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
      })
    }
  }, [getOrganizationResponse, getOrganizationDomainsResponse]);

  const onSubmit = async (values: z.infer<typeof schema>) => {
    if (values.requireMfa) {
      if (!values.logInWithAuthenticatorApp && !values.logInWithPasskey) {
        form.setError('requireMfa', {
          message: 'To require MFA, you must enable either Log in with Authenticator App or Log in with Passkey.',
        })
        return;
      }
    }

    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        displayName: values.displayName,
        logInWithGoogle: values.logInWithGoogle,
        logInWithMicrosoft: values.logInWithMicrosoft,
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
    })

    toast.success('Organization updated successfully');
    navigate(`/organizations/${organizationId}`);
  };

  return (
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/organizations">Organizations</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Edit</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>
        Edit {getOrganizationResponse?.organization?.displayName}
      </PageTitle>

      <Form {...form}>
        { }
        {/** There's an issue with the types of react-hook-form and zod
        preventing the compiler from inferring the correct types.*/}
        <form onSubmit={form.handleSubmit(onSubmit)} className="mt-8 space-y-8">
          {/** eslint-enable @typescript-eslint/no-unsafe-call */}
          <Card>
            <CardHeader>
              <CardTitle>Organization settings</CardTitle>
              <CardDescription>
                Configure basic settings on this organization.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormControl>
                      <Input
                        className="max-w-80"
                        placeholder="Acme Corporation"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      A human-friendly name for the organization.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Login settings</CardTitle>
              <CardDescription>
                Configure how users can log into this organization.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              {getProjectResponse?.project?.logInWithGoogle && (
                <FormField
                  control={form.control}
                  name="logInWithGoogle"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>Log in with Google</FormLabel>
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
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Enterprise settings</CardTitle>
              <CardDescription>
                Configure whether this organization can use SCIM.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              {getProjectResponse?.project?.logInWithSaml && (
                <FormField
                  control={form.control}
                  name="logInWithSaml"
                  render={({ field }: { field: any }) => (
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
                        Whether this organization can configure SAML Connections
                        and use them to log in with SAML.
                      </FormDescription>
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
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Whether this organization can configure SCIM ("Enterprise
                      Directory Sync").
                    </FormDescription>
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
                    <FormControl>
                      <InputTags
                        className="max-w-96"
                        placeholder="example.com"
                        {...field}
                        value={field.value || []}
                      />
                    </FormControl>
                    <FormDescription>
                      SAML and SCIM users must have emails from this list of
                      domains.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          <div className="flex justify-end gap-x-4 pb-8">
            <Button variant="outline" asChild>
              <Link to={`/organizations/${organizationId}`}>Cancel</Link>
            </Button>
            <Button type="submit">Save Changes</Button>
          </div>
        </form>
      </Form>
    </div>
  );
};
