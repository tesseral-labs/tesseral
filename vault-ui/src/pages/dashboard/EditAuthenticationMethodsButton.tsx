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
  logInWithGithub: z.boolean(),
  logInWithGoogle: z.boolean(),
  logInWithMicrosoft: z.boolean(),
  logInWithEmail: z.boolean(),
  logInWithPassword: z.boolean(),
});

export function EditAuthenticationMethodsButton() {
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithEmail: false,
      logInWithGithub: false,
      logInWithGoogle: false,
      logInWithMicrosoft: false,
      logInWithPassword: false,
    },
  });

  const { data: getOrganizationResponse, refetch: refetchGetOrganization } =
    useQuery(getOrganization);
  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        logInWithEmail: getOrganizationResponse.organization.logInWithEmail,
        logInWithGithub: getOrganizationResponse.organization.logInWithGithub,
        logInWithGoogle: getOrganizationResponse.organization.logInWithGoogle,
        logInWithMicrosoft:
          getOrganizationResponse.organization.logInWithMicrosoft,
        logInWithPassword:
          getOrganizationResponse.organization.logInWithPassword,
      });
    }
  }, [form, getOrganizationResponse]);

  const { mutateAsync: updateOrganizationAsync } =
    useMutation(updateOrganization);

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateOrganizationAsync({
      organization: {
        logInWithEmail: values.logInWithEmail,
        logInWithGithub: values.logInWithGithub,
        logInWithGoogle: values.logInWithGoogle,
        logInWithMicrosoft: values.logInWithMicrosoft,
        logInWithPassword: values.logInWithPassword,
      },
    });
    await refetchGetOrganization();
    setOpen(false);
    toast.success("Authentication methods updated");
  }

  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit authentication settings</AlertDialogTitle>
          <AlertDialogDescription>
            Choose how users in this organization can authenticate.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            {getProjectResponse?.project?.logInWithGoogle && (
              <FormField
                control={form.control}
                name="logInWithGoogle"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Google</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Users authenticate by signing in with their Google
                      account.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {getProjectResponse?.project?.logInWithMicrosoft && (
              <FormField
                control={form.control}
                name="logInWithMicrosoft"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Microsoft</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Users authenticate by signing in with their Microsoft
                      account.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {getProjectResponse?.project?.logInWithGithub && (
              <FormField
                control={form.control}
                name="logInWithGithub"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with GitHub</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Users authenticate by signing in with their GitHub
                      account.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {getProjectResponse?.project?.logInWithEmail && (
              <FormField
                control={form.control}
                name="logInWithEmail"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Email (Magic Links)</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Users authenticate by visiting a magic link sent to their
                      email address.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {getProjectResponse?.project?.logInWithPassword && (
              <FormField
                control={form.control}
                name="logInWithPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log in with Password</FormLabel>
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormDescription>
                      Users authenticate by entering their email address and
                      password.
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
