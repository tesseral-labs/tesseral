import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Key } from "lucide-react";
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
  logInWithGithub: z.boolean(),
  logInWithGoogle: z.boolean(),
  logInWithMicrosoft: z.boolean(),
});

export function OrganizationOAuthCard() {
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const organization = getOrganizationResponse?.organization;
  const project = getProjectResponse?.project;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGithub: false,
      logInWithGoogle: false,
      logInWithMicrosoft: false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateOrganizationMutation.mutateAsync({
        organization: {
          logInWithGithub: data.logInWithGithub,
          logInWithGoogle: data.logInWithGoogle,
          logInWithMicrosoft: data.logInWithMicrosoft,
        },
      });
      await refetch();
      form.reset(data);
      toast.success("Organization settings updated successfully.");
    } catch {
      toast.error("Failed to update organization settings. Please try again.");
    }
  }

  useEffect(() => {
    if (organization) {
      form.reset({
        logInWithGithub: organization.logInWithGithub || false,
        logInWithGoogle: organization.logInWithGoogle || false,
        logInWithMicrosoft: organization.logInWithMicrosoft || false,
      });
    }
  }, [organization, form]);

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Key />
              OAuth Providers
            </CardTitle>
            <CardDescription>
              Configure OAuth providers for your organization.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow space-y-6">
            {project?.logInWithGoogle && (
              <FormField
                control={form.control}
                name="logInWithGoogle"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Log in with Google</FormLabel>
                      <FormDescription>
                        Allows users to log into this organization using their
                        Google accounts.
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
            {project?.logInWithMicrosoft && (
              <FormField
                control={form.control}
                name="logInWithMicrosoft"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Log in with Microsoft</FormLabel>
                      <FormDescription>
                        Allows users to log into this organization using their
                        Microsoft accounts.
                      </FormDescription>
                    </div>
                    <FormMessage />
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
            {project?.logInWithGithub && (
              <FormField
                control={form.control}
                name="logInWithGithub"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Log in with GitHub</FormLabel>
                      <FormDescription>
                        Allows users to log into this organization using their
                        GitHub accounts.
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
