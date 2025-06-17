import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Fingerprint, LoaderCircle } from "lucide-react";
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

export function MfaSettingsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Fingerprint />
          Multi-Factor Authentication (MFA)
        </CardTitle>
        <CardDescription>
          Configure which Multi-factor Authentication methods are available to
          users when authenticating.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">Authenticator Apps</div>
            <div>
              {getProjectResponse?.project?.logInWithAuthenticatorApp ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">Passkeys</div>
            <div>
              {getProjectResponse?.project?.logInWithPasskey ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <ConfigureMfaButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithAuthenticatorApp: z.boolean(),
  logInWithPasskey: z.boolean(),
});

export function ConfigureMfaButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithAuthenticatorApp:
        getProjectResponse?.project?.logInWithAuthenticatorApp ?? false,
      logInWithPasskey: getProjectResponse?.project?.logInWithPasskey ?? false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithAuthenticatorApp: data.logInWithAuthenticatorApp,
        logInWithPasskey: data.logInWithPasskey,
      },
    });
    await refetch();
    form.reset(data);
    setOpen(false);
    toast.success("MFA settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse?.project) {
      form.reset({
        logInWithAuthenticatorApp:
          getProjectResponse.project.logInWithAuthenticatorApp,
        logInWithPasskey: getProjectResponse.project.logInWithPasskey,
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full">
          Configure MFA Settings
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Mutli-factor Authentication Settings</DialogTitle>
          <DialogDescription>
            Configure multi-factor authentication (MFA) options for your
            organization. You can enable login with an authenticator app or
            passkey.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="logInWithAuthenticatorApp"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Authenticator Apps</FormLabel>
                      <FormDescription>
                        Enable MFA using authenticator apps like Google
                        Authenticator or Authy.
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
                name="logInWithPasskey"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div>
                      <FormLabel>Passkeys</FormLabel>
                      <FormDescription>
                        Enable MFA using passkeys for enhanced security.
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
            </div>
            <DialogFooter className="mt-8">
              <Button type="button" variant="outline" onClick={handleCancel}>
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
