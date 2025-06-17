import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Key, LoaderCircle } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
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
  getOrganization,
  getProject,
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  logInWithGoogle: z.boolean().optional(),
  logInWithGithub: z.boolean().optional(),
  logInWithMicrosoft: z.boolean().optional(),
});

export function OrganizationOAuthCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithGoogle:
        getOrganizationResponse?.organization?.logInWithGoogle || false,
      logInWithGithub:
        getOrganizationResponse?.organization?.logInWithGithub || false,
      logInWithMicrosoft:
        getOrganizationResponse?.organization?.logInWithMicrosoft || false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        logInWithGoogle: data.logInWithGoogle,
        logInWithGithub: data.logInWithGithub,
        logInWithMicrosoft: data.logInWithMicrosoft,
      },
    });
    form.reset(data);
    await refetch();
    setOpen(false);
    toast.success("OAuth settings updated successfully");
  }

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        logInWithGoogle: getOrganizationResponse.organization.logInWithGoogle,
        logInWithGithub: getOrganizationResponse.organization.logInWithGithub,
        logInWithMicrosoft:
          getOrganizationResponse.organization.logInWithMicrosoft,
      });
    }
  }, [getOrganizationResponse, form]);

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
              Configure OAuth providers for{" "}
              <span className="font-semibold">
                {getOrganizationResponse?.organization?.displayName}
              </span>
              .
            </CardDescription>
          </CardHeader>
          <CardContent className="flex-grow">
            <div className="space-y-4">
              {getProjectResponse?.project?.logInWithGoogle && (
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
              {getProjectResponse?.project?.logInWithMicrosoft && (
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
              {getProjectResponse?.project?.logInWithGithub && (
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
