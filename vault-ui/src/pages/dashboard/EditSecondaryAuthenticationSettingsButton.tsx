import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
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

export function EditSecondaryAuthenticationSettingsButton() {
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithAuthenticatorApp: false,
      logInWithPasskey: false,
      requireMfa: false,
    },
  });

  const { data: getOrganizationResponse, refetch: refetchGetOrganization } =
    useQuery(getOrganization);
  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        logInWithAuthenticatorApp:
          getOrganizationResponse.organization.logInWithAuthenticatorApp,
        logInWithPasskey: getOrganizationResponse.organization.logInWithPasskey,
        requireMfa: getOrganizationResponse.organization.requireMfa,
      });
    }
  }, [form, getOrganizationResponse]);

  const { mutateAsync: updateOrganizationAsync } =
    useMutation(updateOrganization);

  async function handleSubmit(values: z.infer<typeof schema>) {
    if (values.requireMfa) {
      if (!values.logInWithAuthenticatorApp && !values.logInWithPasskey) {
        form.setError("requireMfa", {
          message: 'To require MFA, you must enable either Log in with Authenticator App or Log in with Passkey.',
        })
        return;
      }
    }

    await updateOrganizationAsync({
      organization: {
        logInWithAuthenticatorApp: values.logInWithAuthenticatorApp,
        logInWithPasskey: values.logInWithPasskey,
        requireMfa: values.requireMfa,
      },
    });
    await refetchGetOrganization();
    setOpen(false);
    toast.success("Secondary authentication settings updated");
  }

  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Edit secondary authentication settings
          </AlertDialogTitle>
          <AlertDialogDescription>
            Configure multi-factor authentication.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="requireMfa"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Require MFA</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Require that users configure multi-factor authentication.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            {getProjectResponse?.project?.logInWithAuthenticatorApp && (
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
                      Users can configure an authenticator app.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {getProjectResponse?.project?.logInWithPasskey && (
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
                      Users can configure passkeys.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Update</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
