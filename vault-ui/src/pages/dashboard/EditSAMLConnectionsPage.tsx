import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import {
  getSAMLConnection,
  updateSAMLConnection,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { parseErrorMessage } from "@/lib/errors";

const schema = z.object({
  primary: z.boolean(),
  idpEntityId: z.string().min(1, {
    message: "IDP Entity ID must be non-empty.",
  }),
  idpRedirectUrl: z.string().url({
    message: "IDP Redirect URL must be a valid URL.",
  }),
  idpX509Certificate: z.string().startsWith("-----BEGIN CERTIFICATE-----", {
    message: "IDP Certificate must be a PEM-encoded X.509 certificate.",
  }),
});

export function EditSAMLConnectionsPage () {
  const navigate = useNavigate();
  const params = useParams();

  const { data } = useQuery(getSAMLConnection, {
    id: params.samlConnectionId,
  });

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {},
  });
  const updateSAMLConnectionMutation = useMutation(updateSAMLConnection);

  async function handleSubmit (values: z.infer<typeof schema>) {
    try {
      await updateSAMLConnectionMutation.mutateAsync({
        id: params.samlConnectionId,
        samlConnection: {
          primary: values.primary,
          idpEntityId: values.idpEntityId,
          idpRedirectUrl: values.idpRedirectUrl,
          idpX509Certificate: values.idpX509Certificate,
        },
      });

      navigate(`/organization`);
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error("Could not update SAML connection", {
        description: message,
      });
    }
  };

  useEffect(() => {
    if (data?.samlConnection) {
      form.reset({
        primary: data.samlConnection.primary,
        idpEntityId: data.samlConnection.idpEntityId,
        idpRedirectUrl: data.samlConnection.idpRedirectUrl,
        idpX509Certificate: data.samlConnection.idpX509Certificate,
      });
    }
  }, [data]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="dark:text-foreground">
      <div className="mb-4">
        <h1 className="text-2xl font-bold mb-2">Edit SAML Connection</h1>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>SAML connection settings</CardTitle>
              <p>Configure basic settings on this SAML connection</p>
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
                      A primary SAML connection gets used by default within its
                      organization.
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
                and needs to be inputted into your customer's Identity Provider
                by their IT admin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-8">
              <div>
                <div className="text-sm font-medium leading-none">
                  Assertion Consumer Service (ACS) URL
                </div>
                <div className="mt-1">{data?.samlConnection?.spAcsUrl}</div>
              </div>
              <div>
                <div className="text-sm font-medium leading-none">
                  SP Entity ID
                </div>
                <div className="mt-1">{data?.samlConnection?.spEntityId}</div>
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
                render={({ field }) => (
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
                render={({ field }) => (
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
                render={({ field: { onChange } }) => (
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
                      IDP Certificate, as a PEM-encoded X.509 certificate. These
                      start with '-----BEGIN CERTIFICATE-----' and end with
                      '-----END CERTIFICATE-----'.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          <div className="flex justify-end gap-x-4 pb-8">
            <Button variant="outline" asChild>
              <Link to={`/organization`}>Cancel</Link>
            </Button>
            <Button type="submit">Save Changes</Button>
          </div>
        </form>
      </Form>
    </div>
  );
};
