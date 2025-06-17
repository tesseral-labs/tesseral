import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Fingerprint, LoaderCircle } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function OrganizationMFACard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Fingerprint />
          Multi-Factor Authentication (MFA)
        </CardTitle>
        <CardDescription>
          Configure Multi-Factor Authentication for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}
          </span>
          .
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          {(getProjectResponse?.project?.logInWithAuthenticatorApp ||
            getProjectResponse?.project?.logInWithPasskey) && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">Require MFA</div>
                <div className="text-xs text-muted-foreground">
                  Require users to complete Multi-Factor Authentication when
                  logging into this Organization.
                </div>
              </div>
              <Switch
                checked={getOrganizationResponse?.organization?.requireMfa}
                disabled
              />
            </div>
          )}
          {getProjectResponse?.project?.logInWithAuthenticatorApp && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">
                  Log in with Authenticator App
                </div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using a TOTP-based
                  Authenticator App.
                </div>
              </div>
              <Switch
                checked={
                  getOrganizationResponse?.organization
                    ?.logInWithAuthenticatorApp
                }
                disabled
              />
            </div>
          )}
          {getProjectResponse?.project?.logInWithPasskey && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">Log in with Passkey</div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using Passkeys.
                </div>
              </div>
              <Switch
                checked={
                  getOrganizationResponse?.organization?.logInWithPasskey
                }
                disabled
              />
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <ConfigureOrganizationMfaButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithAuthenticatorApp: z.boolean(),
  logInWithPasskey: z.boolean(),
  requireMfa: z.boolean(),
});

function ConfigureOrganizationMfaButton() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithAuthenticatorApp:
        getOrganizationResponse?.organization?.logInWithAuthenticatorApp ||
        false,
      logInWithPasskey:
        getOrganizationResponse?.organization?.logInWithPasskey || false,
      requireMfa: getOrganizationResponse?.organization?.requireMfa || false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        logInWithAuthenticatorApp: data.logInWithAuthenticatorApp,
        logInWithPasskey: data.logInWithPasskey,
        requireMfa: data.requireMfa,
      },
    });
    form.reset(data);
    await refetch();
    setOpen(false);
    toast.success("MFA configuration updated successfully.");
  }

  useEffect(() => {
    if (getOrganizationResponse) {
      form.reset({
        logInWithAuthenticatorApp:
          getOrganizationResponse.organization?.logInWithAuthenticatorApp ||
          false,
        logInWithPasskey:
          getOrganizationResponse.organization?.logInWithPasskey || false,
        requireMfa: getOrganizationResponse.organization?.requireMfa || false,
      });
    }
  }, [getOrganizationResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          Configure MFA
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure MFA</DialogTitle>
          <DialogDescription>
            Configure Multi-Factor Authentication settings for{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            .
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              {(getProjectResponse?.project?.logInWithAuthenticatorApp ||
                getProjectResponse?.project?.logInWithPasskey) && (
                <FormField
                  control={form.control}
                  name="requireMfa"
                  render={({ field }) => (
                    <FormItem className="flex justify-between items-center gap-4">
                      <div>
                        <FormLabel>Require MFA for Login</FormLabel>
                        <FormDescription>
                          Require users to complete Multi-Factor Authentication
                          when logging into this organization.
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
              {getProjectResponse?.project?.logInWithAuthenticatorApp && (
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
              {getProjectResponse?.project?.logInWithPasskey && (
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
            </div>
            <DialogFooter className="mt-8 justify-end gap-2">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
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
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
