import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Settings } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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
  logInWithGithub: z.boolean(),
  githubOauthClientId: z.string(),
  githubOauthClientSecret: z.string(),
});

export function ConfigureGithubOAuthButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGithub: getProjectResponse?.project?.logInWithGithub || false,
      githubOauthClientId:
        getProjectResponse?.project?.githubOauthClientId || "",
      githubOauthClientSecret:
        getProjectResponse?.project?.githubOauthClientSecret || "",
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
      !data.logInWithGithub &&
      !getProjectResponse?.project?.logInWithEmail &&
      !getProjectResponse?.project?.logInWithPassword &&
      !getProjectResponse?.project?.logInWithGoogle &&
      !getProjectResponse?.project?.logInWithMicrosoft
    ) {
      form.setError("logInWithGithub", {
        message:
          "At least one of Log in with Email, Log in with Password, Log in with Google, Log in with Microsoft, or Log in with GitHub must be enabled.",
      });
      return;
    }
    if (
      data.logInWithGithub &&
      data.githubOauthClientId === "" &&
      !getProjectResponse?.project?.githubOauthClientId
    ) {
      form.setError("githubOauthClientId", {
        message:
          "GitHub OAuth Client ID is required when enabling GitHub login.",
      });
      return;
    }
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithGithub: data.logInWithGithub,
        githubOauthClientId: data.githubOauthClientId,
        githubOauthClientSecret: data.githubOauthClientSecret,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("GitHub OAuth settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        logInWithGithub: getProjectResponse.project?.logInWithGithub,
        githubOauthClientId:
          getProjectResponse.project?.githubOauthClientId || "",
        githubOauthClientSecret:
          getProjectResponse.project?.githubOauthClientSecret || "",
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
          <DialogTitle>Configure GitHub OAuth</DialogTitle>
          <DialogDescription>
            Configure GitHub OAuth settings for your project. You will need to
            provide the Client ID and Client Secret obtained from your GitHub
            OAuth application.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="logInWithGithub"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between space-x-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with GitHub</FormLabel>
                      <FormDescription>
                        Whether Users can log in using their GitHub account.
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
                name="githubOauthClientId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel></FormLabel>
                    <FormDescription>
                      Your company's GitHub OAuth Client ID.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="GitHub OAuth Client ID" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="githubOauthClientSecret"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel></FormLabel>
                    <FormDescription>
                      Your company's GitHub OAuth Client Secret.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="GitHub OAuth Client Secret"
                        {...field}
                      />
                    </FormControl>
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
