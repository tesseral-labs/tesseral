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

export function OrganizationOAuthCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
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
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">Log in with Google</div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using their
                  Microsoft accounts.
                </div>
              </div>
              <Switch
                checked={getOrganizationResponse?.organization?.logInWithGoogle}
                disabled
              />
            </div>
          )}
          {getProjectResponse?.project?.logInWithMicrosoft && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">
                  Log in with Microsoft
                </div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using their
                  Microsoft accounts.
                </div>
              </div>
              <Switch
                checked={
                  getOrganizationResponse?.organization?.logInWithMicrosoft
                }
                disabled
              />
            </div>
          )}
          {getProjectResponse?.project?.logInWithGithub && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">Log in with GitHub</div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using their GitHub
                  accounts.
                </div>
              </div>
              <Switch
                checked={getOrganizationResponse?.organization?.logInWithGithub}
                disabled
              />
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <ConfigureOrganizationOauthButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithGoogle: z.boolean().optional(),
  logInWithGithub: z.boolean().optional(),
  logInWithMicrosoft: z.boolean().optional(),
});

function ConfigureOrganizationOauthButton() {
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
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full">
          Configure OAuth Providers
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure OAuth Providers</DialogTitle>
          <DialogDescription>
            Configure OAuth providers for{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            .
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
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
            <DialogFooter className="mt-8 justify-end gap-2">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                disabled={
                  !form.formState.isDirty ||
                  updateOrganizationMutation.isPending
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
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
