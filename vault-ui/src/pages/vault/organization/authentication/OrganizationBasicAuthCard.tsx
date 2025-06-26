import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Lock } from "lucide-react";
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
  getOrganization,
  getProject,
  updateOrganization,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  logInWithEmail: z.boolean(),
  logInWithPassword: z.boolean(),
});

export function OrganizationBasicAuthCard() {
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithEmail: false,
      logInWithPassword: false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateOrganizationMutation.mutateAsync({
        organization: {
          logInWithEmail: data.logInWithEmail,
          logInWithPassword: data.logInWithPassword,
        },
      });
      await refetch();
      form.reset(data);
      toast.success("Basic authentication settings updated successfully.");
    } catch {
      toast.error(
        "Failed to update basic authentication settings. Please try again.",
      );
    }
  }

  useEffect(() => {
    if (organization) {
      form.reset({
        logInWithEmail: organization.logInWithEmail || false,
        logInWithPassword: organization.logInWithPassword || false,
      });
    }
  }, [organization, form]);

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
              Configure basic authentication for your organization.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow space-y-6">
            {project?.logInWithEmail && (
              <FormField
                control={form.control}
                name="logInWithEmail"
                render={({ field }) => (
                  <FormItem className="flex items-start justify-between gap-x-4">
                    <div>
                      <FormLabel>Log in with Email Magic Link</FormLabel>
                      <FormDescription>
                        Allow users in this Organization to log in with an email
                        magic link.
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
            {project?.logInWithPassword && (
              <FormField
                control={form.control}
                name="logInWithPassword"
                render={({ field }) => (
                  <FormItem className="flex items-start justify-between gap-x-4">
                    <div>
                      <FormLabel>Log in with Password</FormLabel>
                      <FormDescription>
                        Allow users in this Organization to log in with a
                        password.
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
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
              disabled={
                !form.formState.isDirty || updateOrganizationMutation.isPending
              }
              type="submit"
            >
              Save changes
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
