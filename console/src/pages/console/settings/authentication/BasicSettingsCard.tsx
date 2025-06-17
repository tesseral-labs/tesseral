import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Lock } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Badge } from "@/components/ui/badge";
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
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function BasicSettingsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Lock />
          Basic Authentication
        </CardTitle>
        <CardDescription>
          Configure basic authentication settings for how users can
          authentication to Organizations in your Project.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">Email Magic Links</div>
            <div>
              {getProjectResponse?.project?.logInWithEmail ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">Password Authentication</div>
            <div>
              {getProjectResponse?.project?.logInWithPassword ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <ConfigureBasicAuthenticationButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithEmail: z.boolean(),
  logInWithPassword: z.boolean(),
});

function ConfigureBasicAuthenticationButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithEmail: getProjectResponse?.project?.logInWithEmail ?? false,
      logInWithPassword:
        getProjectResponse?.project?.logInWithPassword ?? false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    if (
      !data.logInWithEmail &&
      !data.logInWithPassword &&
      !getProjectResponse?.project?.logInWithGoogle &&
      !getProjectResponse?.project?.logInWithMicrosoft &&
      !getProjectResponse?.project?.logInWithGithub
    ) {
      form.setError("logInWithEmail", {
        message:
          "At least one of Log in with Email, Log in with Password, Log in with Google, Log in with Microsoft, or Log in with GitHub must be enabled.",
      });
      return;
    }
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithEmail: data.logInWithEmail,
        logInWithPassword: data.logInWithPassword,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("Basic Auth settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithEmail: getProjectResponse.project.logInWithEmail,
        logInWithPassword: getProjectResponse.project.logInWithPassword,
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          Configure Basic Authentication
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Basic Authentication</DialogTitle>
          <DialogDescription>
            Configure what basic authentication methods are available to
            Organization users in this Project.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="logInWithEmail"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with Email Magic Link</FormLabel>
                      <FormDescription>
                        Allow users to log in using their email address. A magic
                        link will be sent to their email for authentication.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={(checked) => field.onChange(checked)}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="logInWithPassword"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with Password</FormLabel>
                      <FormDescription>
                        Allow users to log in using an email and password.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={(checked) => field.onChange(checked)}
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
                type="submit"
                disabled={
                  !form.formState.isDirty || updateProjectMutation.isPending
                }
              >
                {updateProjectMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {updateProjectMutation.isPending
                  ? "Saving Changes"
                  : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
