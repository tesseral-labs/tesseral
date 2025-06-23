import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Shield } from "lucide-react";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { InputTags } from "@/components/core/InputTags";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
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
import { Switch } from "@/components/ui/switch";
import {
  getOrganization,
  getOrganizationDomains,
  getProject,
  updateOrganizationDomains,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { updateOrganization } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  domains: z.array(z.string()).optional(),
  logInWithSaml: z.boolean().default(false),
});

export function OrganizationSamlCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getOrganizationDomainsResponse, refetch: refetchDomains } =
    useQuery(getOrganizationDomains, {
      organizationId: organizationId,
    });
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);
  const updateOrganizationDomainsMutation = useMutation(
    updateOrganizationDomains,
  );

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      domains:
        getOrganizationDomainsResponse?.organizationDomains?.domains || [],
      logInWithSaml:
        getOrganizationResponse?.organization?.logInWithSaml || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        logInWithSaml: data.logInWithSaml,
      },
    });
    await updateOrganizationDomainsMutation.mutateAsync({
      organizationId: organizationId,
      organizationDomains: {
        domains: data.domains || [],
      },
    });
    await refetch();
    await refetchDomains();
    form.reset();
    toast.success("SAML configuration updated successfully.");
  }

  useEffect(() => {
    if (getOrganizationResponse && getOrganizationDomainsResponse) {
      form.reset({
        logInWithSaml:
          getOrganizationResponse.organization?.logInWithSaml || false,
        domains:
          getOrganizationDomainsResponse.organizationDomains?.domains || [],
      });
    }
  }, [getOrganizationResponse, getOrganizationDomainsResponse, form]);

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield />
              SAML SSO
            </CardTitle>
            <CardDescription>
              Configure SAML authentication for{" "}
              <span className="font-semibold">
                {getOrganizationResponse?.organization?.displayName}
              </span>
              .
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow">
            <div className="space-y-4">
              {getProjectResponse?.project?.logInWithSaml ? (
                <FormField
                  control={form.control}
                  name="logInWithSaml"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div>
                        <FormLabel>Log in with SAML</FormLabel>
                        <FormDescription>
                          Allows users to log into this organization with
                          SAML-based identity providers.
                        </FormDescription>
                        <FormMessage />
                      </div>
                      <FormControl>
                        <Switch
                          id="logInWithSaml"
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              ) : (
                <>
                  <div className="text-sm text-muted-foreground flex-grow">
                    SAML authentication is not enabled for this project. Please
                    enable SAML at the project level to configure it for this
                    organization.
                  </div>
                </>
              )}

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
            </div>
          </CardContent>
          <CardFooter className="mt-4">
            {getProjectResponse?.project?.logInWithSaml ? (
              <Button
                className="w-full"
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  updateOrganizationMutation.isPending
                }
              >
                {updateOrganizationMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {updateOrganizationMutation.isPending
                  ? "Saving changes"
                  : "Save changes"}
              </Button>
            ) : (
              <Link className="w-full" to="/settings/authentication">
                <Button className="w-full" variant="outline">
                  Manage Project SAML Settings
                </Button>
              </Link>
            )}
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
