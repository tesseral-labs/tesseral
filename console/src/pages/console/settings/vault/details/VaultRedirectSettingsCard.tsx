import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Split } from "lucide-react";
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
import { Input } from "@/components/ui/input";
import {
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  redirectUri: z
    .string()
    .url("Must be a valid URL")
    .min(1, "Default Redirect URI is required"),
  afterLoginRedirectUri: z
    .string()
    .url("Must be a valid URL")
    .optional()
    .or(z.literal("")),
  afterSignupRedirectUri: z
    .string()
    .url("Must be a valid URL")
    .optional()
    .or(z.literal("")),
});

export function VaultRedirectSettingsCard() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      redirectUri: getProjectResponse?.project?.redirectUri || "",
      afterLoginRedirectUri:
        getProjectResponse?.project?.afterLoginRedirectUri || "",
      afterSignupRedirectUri:
        getProjectResponse?.project?.afterSignupRedirectUri || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    // Ensure trusted domains are updated with domains from the provided URIs
    const trustedDomains = new Set(getProjectResponse!.project!.trustedDomains);
    if (
      data.redirectUri.startsWith("http://") ||
      data.redirectUri.startsWith("https://")
    ) {
      trustedDomains.add(new URL(data.redirectUri).host);
    }
    if (
      data.afterLoginRedirectUri?.startsWith("http://") ||
      data.afterLoginRedirectUri?.startsWith("https://")
    ) {
      trustedDomains.add(new URL(data.afterLoginRedirectUri).host);
    }
    if (
      data.afterSignupRedirectUri?.startsWith("http://") ||
      data.afterSignupRedirectUri?.startsWith("https://")
    ) {
      trustedDomains.add(new URL(data.afterSignupRedirectUri).host);
    }

    await updateProjectMutation.mutateAsync({
      project: {
        afterLoginRedirectUri: data.afterLoginRedirectUri,
        afterSignupRedirectUri: data.afterSignupRedirectUri,
        redirectUri: data.redirectUri,
        trustedDomains: Array.from(trustedDomains),
      },
    });
    await refetch();
    form.reset(data);
    toast.success("Vault redirect settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse?.project && form) {
      form.reset({
        redirectUri: getProjectResponse.project.redirectUri || "",
        afterLoginRedirectUri:
          getProjectResponse.project.afterLoginRedirectUri || "",
        afterSignupRedirectUri:
          getProjectResponse.project.afterSignupRedirectUri || "",
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Split />
              <span>Redirect Settings</span>
            </CardTitle>
            <CardDescription>
              Configure where your users are redirected after they authenticate.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 flex-grow">
            <FormField
              control={form.control}
              name="redirectUri"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Default Redirect URI</FormLabel>
                  <FormDescription>
                    The URI users are redirected to after authentication.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      {...field}
                      type="url"
                      placeholder="https://example.com/redirect"
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="afterLoginRedirectUri"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>After Login Redirect URI</FormLabel>
                  <FormDescription>
                    The URI users are redirected to after logging in.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      {...field}
                      type="url"
                      placeholder="https://example.com/login-redirect"
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="afterSignupRedirectUri"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>After Signup Redirect URI</FormLabel>
                  <FormDescription>
                    The URI users are redirected to after signing up.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      {...field}
                      type="url"
                      placeholder="https://example.com/signup-redirect"
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
              disabled={
                !form.formState.isDirty || updateProjectMutation.isPending
              }
              type="submit"
            >
              {updateProjectMutation.isPending && (
                <LoaderCircle className="animate-spin" />
              )}
              {updateProjectMutation.isPending
                ? "Saving changes"
                : "Save changes"}
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
