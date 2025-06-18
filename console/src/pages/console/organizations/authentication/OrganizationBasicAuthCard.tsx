import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Lock } from "lucide-react";
import React from "react";
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

const schema = z.object({
  logInWithEmail: z.boolean().optional(),
  logInWithPassword: z.boolean().optional(),
});

export function OrganizationBasicAuthCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithEmail:
        getOrganizationResponse?.organization?.logInWithEmail || false,
      logInWithPassword:
        getOrganizationResponse?.organization?.logInWithPassword || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        logInWithEmail: data.logInWithEmail,
        logInWithPassword: data.logInWithPassword,
      },
    });
    refetch();
    form.reset(data);
    toast.success("Basic authentication settings updated successfully");
  }

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Lock />
              Basic Authentication
            </CardTitle>
            <CardDescription>
              Configure basic authentication for your{" "}
              <span className="font-semibold">
                {getOrganizationResponse?.organization?.displayName}
              </span>
              .
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow">
            <div className="space-y-4">
              {getProjectResponse?.project?.logInWithEmail && (
                <FormField
                  control={form.control}
                  name="logInWithEmail"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div className="space-y-2">
                        <FormLabel>Log in with Email Magic Link</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in with an
                          email magic link.
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
              {getProjectResponse?.project?.logInWithPassword && (
                <FormField
                  control={form.control}
                  name="logInWithPassword"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div className="space-y-2">
                        <FormLabel>Log in with Password</FormLabel>
                        <FormDescription>
                          Whether Users in this Organization can log in with an
                          email and password.
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
          </CardContent>
          <CardFooter className="mt-4">
            <Button
              className="w-full"
              disabled={
                !form.formState.isDirty || updateOrganizationMutation.isPending
              }
              type="submit"
            >
              {updateOrganizationMutation.isPending && (
                <LoaderCircle className="animate-spin" />
              )}
              {updateOrganizationMutation.isPending
                ? "Saving changes"
                : "Save changes"}
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
