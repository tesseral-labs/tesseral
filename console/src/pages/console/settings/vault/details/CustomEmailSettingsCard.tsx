import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Mail } from "lucide-react";
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
  getProject,
  updateProject,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  customEmailVerifyEmail: z.boolean().optional(),
  customEmailPasswordReset: z.boolean().optional(),
  customEmailUserInvite: z.boolean().optional(),
});

export function CustomEmailSettingsCard() {
  const { data: getProjectResponse, refetch: refetchProject } =
    useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      customEmailVerifyEmail:
        getProjectResponse?.project?.customEmailVerifyEmail ?? false,
      customEmailPasswordReset:
        getProjectResponse?.project?.customEmailPasswordReset ?? false,
      customEmailUserInvite:
        getProjectResponse?.project?.customEmailUserInvite ?? false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        customEmailVerifyEmail: data.customEmailVerifyEmail,
        customEmailPasswordReset: data.customEmailPasswordReset,
        customEmailUserInvite: data.customEmailUserInvite,
      },
    });
    await refetchProject();
    form.reset(data);
    toast.success("Custom email settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse && form) {
      form.reset({
        customEmailVerifyEmail:
          getProjectResponse?.project?.customEmailVerifyEmail ?? false,
        customEmailPasswordReset:
          getProjectResponse?.project?.customEmailPasswordReset ?? false,
        customEmailUserInvite:
          getProjectResponse?.project?.customEmailUserInvite ?? false,
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Mail />
              <span>Custom Email Settings</span>
            </CardTitle>
            <CardDescription>
              By default, Tesseral will send emails on your behalf. You can have
              Tesseral send you a webhook instead.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 flex-grow">
            <FormField
              control={form.control}
              name="customEmailVerifyEmail"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div className="space-y-2">
                    <FormLabel>Custom Email Verification</FormLabel>
                    <FormDescription>
                      Send a{" "}
                      <code className="text-xs">custom_email.verify_email</code>{" "}
                      webhook event when Users verify their email address.
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
            <FormField
              control={form.control}
              name="customEmailPasswordReset"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div className="space-y-2">
                    <FormLabel>Custom Password Reset</FormLabel>
                    <FormDescription>
                      Send a{" "}
                      <code className="text-xs">
                        custom_email.password_reset
                      </code>{" "}
                      webhook event when Users forget their password.
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
            <FormField
              control={form.control}
              name="customEmailUserInvite"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div className="space-y-2">
                    <FormLabel>Custom User Invite</FormLabel>
                    <FormDescription>
                      Send a{" "}
                      <code className="text-xs">custom_email.user_invite</code>{" "}
                      webhook event when Users are invited.
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
          </CardContent>
          <CardFooter className="mt-4">
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
