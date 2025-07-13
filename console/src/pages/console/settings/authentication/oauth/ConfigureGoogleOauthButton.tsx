import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ExternalLink, LoaderCircle, Settings } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import { toast } from "sonner";
import { z } from "zod";
import clsx from "clsx";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
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
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import {
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  logInWithGoogle: z.boolean(),
  useDefaultGoogleClient: z.boolean(),
  googleOauthClientId: z.string(),
  googleOauthClientSecret: z.string(),
});

export function ConfigureGoogleOAuthButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGoogle: getProjectResponse?.project?.logInWithGoogle || false,
      useDefaultGoogleClient: !getProjectResponse?.project?.googleOauthClientId,
      googleOauthClientId:
        getProjectResponse?.project?.googleOauthClientId || "",
      googleOauthClientSecret:
        getProjectResponse?.project?.googleOauthClientSecret || "",
    },
  });

  function handleCancel(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    if (
      !data.logInWithGoogle &&
      !getProjectResponse?.project?.logInWithEmail &&
      !getProjectResponse?.project?.logInWithPassword &&
      !getProjectResponse?.project?.logInWithMicrosoft &&
      !getProjectResponse?.project?.logInWithGithub
    ) {
      form.setError("logInWithGoogle", {
        message:
          "At least one of Log in with Email, Log in with Password, Log in with Google, Log in with Microsoft, or Log in with GitHub must be enabled.",
      });
      return;
    }
    if (
      data.logInWithGoogle &&
      !data.useDefaultGoogleClient &&
      !data.googleOauthClientId
    ) {
      form.setError("googleOauthClientId", {
        message:
          "Google OAuth Client ID is required when using custom client.",
      });
      return;
    }
    if (
      data.logInWithGoogle &&
      !data.useDefaultGoogleClient &&
      !data.googleOauthClientSecret
    ) {
      form.setError("googleOauthClientSecret", {
        message:
          "Google OAuth Client Secret is required when using custom client.",
      });
      return;
    }
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithGoogle: data.logInWithGoogle,
        googleOauthClientId: data.useDefaultGoogleClient ? "" : data.googleOauthClientId,
        googleOauthClientSecret: data.useDefaultGoogleClient ? "" : data.googleOauthClientSecret,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("Google OAuth settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        logInWithGoogle: getProjectResponse.project?.logInWithGoogle,
        useDefaultGoogleClient: !getProjectResponse.project?.googleOauthClientId,
        googleOauthClientId:
          getProjectResponse.project?.googleOauthClientId || "",
        googleOauthClientSecret:
          getProjectResponse.project?.googleOauthClientSecret || "",
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Settings />
          <span className="hidden lg:block">Configure</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Google OAuth</DialogTitle>
          <DialogDescription>
            Configure Google OAuth settings for your project.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="logInWithGoogle"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between space-x-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with Google</FormLabel>
                      <FormDescription>
                        Whether Users can log in using their Google account.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="useDefaultGoogleClient"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>OAuth Client Configuration</FormLabel>
                    <FormDescription>
                      Choose whether to use Tesseral's default Google OAuth client or provide your own.
                    </FormDescription>
                    <FormControl>
                      <RadioGroup
                        value={field.value ? "default" : "custom"}
                        onValueChange={(value) => field.onChange(value === "default")}
                        className="mt-2 gap-0 -space-y-px"
                      >
                        <div className={clsx(
                          "flex items-start space-x-3 p-4 border rounded-t-lg",
                          field.value && "bg-muted/50 border-primary relative"
                        )}>
                          <RadioGroupItem value="default" className="mt-1" />
                          <div className="flex-1 space-y-1">
                            <Label className="text-sm font-medium">
                              Default Google OAuth Client
                            </Label>
                            <p className="text-xs text-muted-foreground">
                              Use Tesseral's preconfigured Google OAuth client. This is the easiest option and requires no additional setup.
                            </p>
                          </div>
                        </div>
                        <div className={clsx(
                          "border rounded-b-lg",
                          !field.value && "bg-muted/50 border-primary"
                        )}>
                          <div className="flex items-start space-x-3 p-4">
                            <RadioGroupItem value="custom" className="mt-1" />
                            <div className="flex-1 space-y-1">
                              <Label className="text-sm font-medium">
                                Custom Google OAuth Client
                              </Label>
                              <p className="text-xs text-muted-foreground">
                                Use your own Google OAuth application with custom branding and settings.
                              </p>
                            </div>
                          </div>
                          {!field.value && (
                            <div className="p-4 pt-0 space-y-4">
                              <div className="flex flex-col gap-2 text-sm">
                                <Label>Callback URL</Label>
                                <span>
                                  Use this as the Authorized redirect URI in your Google OAuth app settings.{" "}
                                  <Link
                                    to="https://tesseral.com/docs/login-methods/primary-factors/log-in-with-google"
                                    target="_blank"
                                    className="underline"
                                  >
                                    Docs <ExternalLink className="inline size-3" />
                                  </Link>
                                </span>
                                <ValueCopier
                                  value={`https://${getProjectResponse?.project?.vaultDomain}/google-oauth-callback`}
                                />
                              </div>
                              <FormField
                                control={form.control}
                                name="googleOauthClientId"
                                render={({ field: clientIdField }) => (
                                  <FormItem>
                                    <FormLabel>Client ID</FormLabel>
                                    <FormDescription>
                                      Your company's Google OAuth Client ID.
                                    </FormDescription>
                                    <FormControl>
                                      <Input
                                        placeholder="Google OAuth Client ID"
                                        {...clientIdField}
                                      />
                                    </FormControl>
                                    <FormMessage />
                                  </FormItem>
                                )}
                              />
                              <FormField
                                control={form.control}
                                name="googleOauthClientSecret"
                                render={({ field: clientSecretField }) => (
                                  <FormItem>
                                    <FormLabel>Client Secret</FormLabel>
                                    <FormDescription>
                                      Your company's Google OAuth Client Secret.
                                    </FormDescription>
                                    <FormControl>
                                      <Input
                                        type="password"
                                        placeholder="Google OAuth Client Secret"
                                        {...clientSecretField}
                                      />
                                    </FormControl>
                                    <FormMessage />
                                  </FormItem>
                                )}
                              />
                            </div>
                          )}
                        </div>
                      </RadioGroup>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
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
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
