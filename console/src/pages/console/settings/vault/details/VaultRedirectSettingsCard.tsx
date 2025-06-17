import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Split } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
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
import { Input } from "@/components/ui/input";
import {
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultRedirectSettingsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
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
        <div className="flex items-center justify-between gap-2">
          <div className="text-sm font-semibold">Default Redirect URI</div>
          {getProjectResponse?.project?.redirectUri && (
            <ValueCopier
              value={getProjectResponse?.project?.redirectUri}
              label="Default Redirect URI"
            />
          )}
        </div>
        <div className="flex items-center justify-between gap-2">
          <div className="text-sm font-semibold">After Login Redirect URI</div>
          {getProjectResponse?.project?.afterLoginRedirectUri ? (
            <ValueCopier
              value={getProjectResponse?.project?.afterLoginRedirectUri}
              label="Login Redirect URI"
            />
          ) : (
            <span className="text-muted-foreground">—</span>
          )}
        </div>
        <div className="flex flex-col lg:flex-row items-center justify-between gap-2 flex-wrap lg:flex-nowrap">
          <div className="w-full lg:w-auto lg:inline text-sm font-semibold">
            After Signup Redirect URI
          </div>
          <div className="w-full lg:w-auto lg:inline">
            {getProjectResponse?.project?.afterSignupRedirectUri ? (
              <ValueCopier
                value={getProjectResponse?.project?.afterSignupRedirectUri}
                label="Signup Redirect URI"
              />
            ) : (
              <span className="text-muted-foreground">—</span>
            )}
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <ConfigureVaultRedirectSettingsButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  redirectUri: z
    .string()
    .url("Must be a valid URL")
    .min(1, "Default Redirect URI is required"),
  afterLoginRedirectUri: z.string().url("Must be a valid URL").optional(),
  afterSignupRedirectUri: z.string().url("Must be a valid URL").optional(),
});

function ConfigureVaultRedirectSettingsButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

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

  function handleCancel(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

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
    setOpen(false);
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
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full" size="lg">
          Configure Vault Redirect Settings
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Vault Redirect Settings</DialogTitle>
          <DialogDescription>
            Set the URIs where users are redirected after authentication
            actions.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
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
            </div>

            <DialogFooter className="mt-8">
              <Button type="button" variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty || updateProjectMutation.isPending
                }
              >
                {updateProjectMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {updateProjectMutation.isPending
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
