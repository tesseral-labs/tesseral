import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Settings } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

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
import { Switch } from "@/components/ui/switch";
import {
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  logInWithMicrosoft: z.boolean(),
  microsoftOauthClientId: z.string(),
  microsoftOauthClientSecret: z.string(),
});

export function ConfigureMicrosoftOAuthButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithMicrosoft:
        getProjectResponse?.project?.logInWithMicrosoft || false,
      microsoftOauthClientId:
        getProjectResponse?.project?.microsoftOauthClientId || "",
      microsoftOauthClientSecret:
        getProjectResponse?.project?.microsoftOauthClientSecret || "",
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
      !data.logInWithMicrosoft &&
      !getProjectResponse?.project?.logInWithEmail &&
      !getProjectResponse?.project?.logInWithPassword &&
      !getProjectResponse?.project?.logInWithGoogle &&
      !getProjectResponse?.project?.logInWithGithub
    ) {
      form.setError("logInWithMicrosoft", {
        message:
          "At least one of Log in with Email, Log in with Password, Log in with Google, Log in with Microsoft, or Log in with GitHub must be enabled.",
      });
      return;
    }
    if (
      data.logInWithMicrosoft &&
      data.microsoftOauthClientId === "" &&
      !getProjectResponse?.project?.microsoftOauthClientId
    ) {
      form.setError("microsoftOauthClientId", {
        message:
          "GitHub OAuth Client ID is required when enabling GitHub login.",
      });
      return;
    }
    if (
      data.logInWithMicrosoft &&
      data.microsoftOauthClientSecret === "" &&
      !getProjectResponse?.project?.microsoftOauthClientSecret
    ) {
      form.setError("microsoftOauthClientSecret", {
        message:
          "GitHub OAuth Client Secret is required when enabling GitHub login.",
      });
      return;
    }
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithMicrosoft: data.logInWithMicrosoft,
        microsoftOauthClientId: data.microsoftOauthClientId,
        microsoftOauthClientSecret: data.microsoftOauthClientSecret,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("Microsoft OAuth settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        logInWithMicrosoft: getProjectResponse.project?.logInWithMicrosoft,
        microsoftOauthClientId:
          getProjectResponse.project?.microsoftOauthClientId || "",
        microsoftOauthClientSecret:
          getProjectResponse.project?.microsoftOauthClientSecret || "",
      });
    }
  }, [getProjectResponse, form]);

  const callbackUrl =
    getProjectResponse &&
    `https://${getProjectResponse.project?.vaultDomain}/microsoft-oauth-callback`;

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
          <DialogTitle>Configure Microsoft OAuth</DialogTitle>
          <DialogDescription>
            Configure Microsoft OAuth settings for your project. You will need
            to provide the Client ID and Client Secret obtained from your
            Microsoft OAuth application.
            <br />
            <br />
            For instructions on how to set up a Microsoft OAuth application,
            please refer to the{" "}
            <Link
              to="https://tesseral.com/docs/login-methods/primary-factors/log-in-with-microsoft"
              target="_blank"
              className="underline"
            >
              documentation
            </Link>
            .
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="logInWithMicrosoft"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between space-x-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with Microsoft</FormLabel>
                      <FormDescription>
                        Whether Users can log in using their Microsoft account.
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
                name="microsoftOauthClientId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel></FormLabel>
                    <FormDescription>
                      Your company's Microsoft OAuth Client ID.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        placeholder="Microsoft OAuth Client ID"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="microsoftOauthClientSecret"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel></FormLabel>
                    <FormDescription>
                      Your company's Microsoft OAuth Client Secret.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Microsoft OAuth Client Secret"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <div className="mt-4 space-y-2">
                <FormLabel>Callback URL</FormLabel>
                <FormDescription>
                  Use this as the Redirect URI in your Microsoft app
                  registration.
                </FormDescription>
                <Input
                  value={callbackUrl || ""}
                  readOnly
                  tabIndex={0}
                  className="bg-muted text-muted-foreground cursor-default"
                  onFocus={(e) => e.target.select()}
                />
              </div>
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
// Compare
