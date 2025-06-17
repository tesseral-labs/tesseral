import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Lock } from "lucide-react";
import React, { MouseEvent, useState } from "react";
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

export function OrganizationBasicAuthCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
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
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">
                  Log in with Magic Link
                </div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using an email
                  magic link.
                </div>
              </div>
              <Switch
                checked={getOrganizationResponse?.organization?.logInWithEmail}
                disabled
              />
            </div>
          )}
        </div>
        <div className="space-y-4">
          {getProjectResponse?.project?.logInWithPassword && (
            <div className="flex justify-between items-center gap-4">
              <div>
                <div className="font-semibold text-sm">
                  Log in with Password
                </div>
                <div className="text-xs text-muted-foreground">
                  Allows users to log into this organization using an email and
                  password.
                </div>
              </div>
              <Switch
                checked={
                  getOrganizationResponse?.organization?.logInWithPassword
                }
                disabled
              />
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <ConfigureOrganizationBasicAuthButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithEmail: z.boolean().optional(),
  logInWithPassword: z.boolean().optional(),
});

function ConfigureOrganizationBasicAuthButton() {
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
      logInWithEmail:
        getOrganizationResponse?.organization?.logInWithEmail || false,
      logInWithPassword:
        getOrganizationResponse?.organization?.logInWithPassword || false,
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
        logInWithEmail: data.logInWithEmail,
        logInWithPassword: data.logInWithPassword,
      },
    });
    refetch();
    form.reset(data);
    setOpen(false);
    toast.success("Basic authentication settings updated successfully");
  }

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
          <DialogDescription></DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
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

            <DialogFooter className="mt-8 justify-end gap-2">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  updateOrganizationMutation.isPending
                }
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
