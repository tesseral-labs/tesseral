import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganization,
  getSAMLConnection,
  updateSAMLConnection,
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
import { toast } from 'sonner';
import { PageContent, PageHeader, PageTitle } from '@/components/page';

const schema = z.object({
  primary: z.boolean(),
  idpEntityId: z.string().min(1, {
    message: 'IDP Entity ID must be non-empty.',
  }),
  idpRedirectUrl: z.string().url({
    message: 'IDP Redirect URL must be a valid URL.',
  }),
  idpX509Certificate: z.string().startsWith('-----BEGIN CERTIFICATE-----', {
    message: 'IDP Certificate must be a PEM-encoded X.509 certificate.',
  }),
});

export const EditSAMLConnectionPage = () => {
  const navigate = useNavigate();
  const { organizationId, samlConnectionId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getSAMLConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });
  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {},
  });
  const updateSAMLConnectionMutation = useMutation(updateSAMLConnection);

  useEffect(() => {
    if (getSAMLConnectionResponse?.samlConnection) {
      form.reset({
        primary: getSAMLConnectionResponse.samlConnection.primary,
        idpEntityId: getSAMLConnectionResponse.samlConnection.idpEntityId,
        idpRedirectUrl: getSAMLConnectionResponse.samlConnection.idpRedirectUrl,
        idpX509Certificate:
          getSAMLConnectionResponse.samlConnection.idpX509Certificate,
      });
    }
  }, [getSAMLConnectionResponse]);
  /* eslint-enable @typescript-eslint/no-unsafe-call */

  const onSubmit = async (values: z.infer<typeof schema>) => {
    await updateSAMLConnectionMutation.mutateAsync({
      id: samlConnectionId,
      samlConnection: {
        primary: values.primary,
        idpEntityId: values.idpEntityId,
        idpRedirectUrl: values.idpRedirectUrl,
        idpX509Certificate: values.idpX509Certificate,
      },
    });

    toast.success('SAML Connection updated');
    navigate(
      `/organizations/${organizationId}/saml-connections/${samlConnectionId}`,
    );
  };

  return (
    // TODO remove padding when app shell in place
    <>
      <PageHeader>
        <PageTitle>Edit SAML Connection</PageTitle>
      </PageHeader>

      <PageContent>
        <Form {...form}>
          {/* eslint-disable @typescript-eslint/no-unsafe-call */}
          {/** Currently there's an issue with the types of react-hook-form and zod 
        preventing the compiler from inferring the correct types.*/}
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            {/* eslint-enable @typescript-eslint/no-unsafe-call */}
            <Card>
              <CardHeader>
                <CardTitle>SAML connection settings</CardTitle>
                <CardDescription>
                  Configure basic settings on this SAML connection.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-8">
                <FormField
                  control={form.control}
                  name="primary"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Primary</FormLabel>
                      <FormControl>
                        <Switch
                          className="block"
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                      <FormDescription>
                        A primary SAML connection gets used by default within
                        its organization.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Service Provider settings</CardTitle>
                <CardDescription>
                  The configuration here is assigned automatically by Tesseral,
                  and needs to be inputted into your customer's Identity
                  Provider by their IT admin.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-8">
                <div>
                  <div className="text-sm font-medium leading-none">
                    Assertion Consumer Service (ACS) URL
                  </div>
                  <div className="mt-1">
                    {getSAMLConnectionResponse?.samlConnection?.spAcsUrl}
                  </div>
                </div>
                <div>
                  <div className="text-sm font-medium leading-none">
                    SP Entity ID
                  </div>
                  <div className="mt-1">
                    {getSAMLConnectionResponse?.samlConnection?.spEntityId}
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Identity Provider settings</CardTitle>
                <CardDescription>
                  The configuration here needs to be copied over from the
                  customer's Identity Provider ("IDP").
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-8">
                <FormField
                  control={form.control}
                  name="idpEntityId"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>IDP Entity ID</FormLabel>
                      <FormControl>
                        <Input className="max-w-96" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="idpRedirectUrl"
                  render={({ field }: { field: any }) => (
                    <FormItem>
                      <FormLabel>IDP Redirect URL</FormLabel>
                      <FormControl>
                        <Input className="max-w-96" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="idpX509Certificate"
                  render={({
                    field: { onChange },
                  }: {
                    field: { onChange: (value: string) => void };
                  }) => (
                    <FormItem>
                      <FormLabel>IDP Certificate</FormLabel>
                      <FormControl>
                        <Input
                          className="max-w-96"
                          type="file"
                          onChange={async (e) => {
                            // File inputs are special; they are necessarily "uncontrolled", and their value is a FileList.
                            // We just copy over the file's contents to the react-form-hook state manually on input change.
                            if (e.target.files) {
                              onChange(await e.target.files[0].text());
                            }
                          }}
                        />
                      </FormControl>
                      <FormDescription>
                        IDP Certificate, as a PEM-encoded X.509 certificate.
                        These start with '-----BEGIN CERTIFICATE-----' and end
                        with '-----END CERTIFICATE-----'.
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
      </PageContent>
    </>
  );
};
