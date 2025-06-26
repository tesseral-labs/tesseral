import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Fingerprint } from "lucide-react";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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
  getProject,
  updateOrganization,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  logInWithAuthenticatorApp: z.boolean(),
  logInWithPasskey: z.boolean(),
  requireMfa: z.boolean(),
});

export function OrganizationMfaCard() {
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithAuthenticatorApp: false,
      logInWithPasskey: false,
      requireMfa: false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      if (
        data.requireMfa &&
        !data.logInWithAuthenticatorApp &&
        !data.logInWithPasskey
      ) {
        form.setError("requireMfa", {
          message:
            "Either Authenticator App or Passkey log in must be enabled when MFA is required.",
        });
        return;
      }

      await updateOrganizationMutation.mutateAsync({
        organization: {
          logInWithAuthenticatorApp: data.logInWithAuthenticatorApp,
          logInWithPasskey: data.logInWithPasskey,
          requireMfa: data.requireMfa,
        },
      });
      await refetch();
      form.reset(data);
      toast.success("MFA settings updated successfully.");
    } catch {
      toast.error("Failed to update MFA settings. Please try again later.");
    }
  }

  useEffect(() => {
    if (organization) {
      form.reset({
        logInWithAuthenticatorApp:
          organization.logInWithAuthenticatorApp || false,
        logInWithPasskey: organization.logInWithPasskey || false,
        requireMfa: organization.requireMfa || false,
      });
    }
  }, [organization, form]);

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Fingerprint />
              Multi-Factor Authentication (MFA)
            </CardTitle>
            <CardDescription>
              Configure Multi-Factor Authentication for your organization.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow space-y-6">
            <FormField
              control={form.control}
              name="requireMfa"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div>
                    <FormLabel>Require MFA for Login</FormLabel>
                    <FormDescription>
                      Require users to complete Multi-Factor Authentication when
                      logging into this organization.
                    </FormDescription>
                    <FormMessage />
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            {project?.logInWithAuthenticatorApp && (
              <FormField
                control={form.control}
                name="logInWithAuthenticatorApp"
                render={({ field }) => (
                  <FormItem className="flex justify-between items-center gap-4">
                    <div>
                      <FormLabel>Log in with Authenticator App</FormLabel>
                      <FormDescription>
                        Allows users to log into this organization using a
                        TOTP-based Authenticator App.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            )}
            {project?.logInWithPasskey && (
              <FormField
                control={form.control}
                name="logInWithPasskey"
                render={({ field }) => (
                  <FormItem className="flex justify-between items-center gap-4">
                    <div>
                      <FormLabel>Log in with Passkey</FormLabel>
                      <FormDescription>
                        Allows users to log into this organization using
                        Passkeys.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            )}
          </CardContent>
          <CardFooter>
            <Button
              disabled={
                !form.formState.isDirty || updateOrganizationMutation.isPending
              }
              type="submit"
              className="w-full"
            >
              Save changes
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
