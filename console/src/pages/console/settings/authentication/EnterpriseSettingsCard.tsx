import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Settings, Shield } from "lucide-react";
import React, { useEffect, useState } from "react";
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

export function EnterpriseSettingsCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Shield />
          Enterprise Auth
        </CardTitle>
        <CardDescription>
          Configure whether users can log in with SAML SSO and use SCIM
          provisioning.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">SAML SSO</div>
            <div>
              {getProjectResponse?.project?.logInWithSaml ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">OIDC SSO</div>
            <div>
              {getProjectResponse?.project?.logInWithOidc ? (
                <Badge>Enabled</Badge>
              ) : (
                <Badge variant="secondary">Disabled</Badge>
              )}
            </div>
          </div>
          <div className="flex items-center justify-between gap-4">
            <div className="font-semibold text-sm">SCIM Provisioning</div>
            <div>
              <Badge>Always Enabled</Badge>
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <div className="w-full">
          <ConfigureEnterpriseSettingsButton />
        </div>
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  logInWithSaml: z.boolean(),
  logInWithOidc: z.boolean(),
});

function ConfigureEnterpriseSettingsButton() {
  const { data: getProjectResponse, refetch } = useQuery(getProject);
  const updateProjectMutation = useMutation(updateProject);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      logInWithSaml: getProjectResponse?.project?.logInWithSaml ?? false,
      logInWithOidc: getProjectResponse?.project?.logInWithOidc ?? false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectMutation.mutateAsync({
      project: {
        logInWithSaml: data.logInWithSaml,
        logInWithOidc: data.logInWithOidc,
      },
    });
    await refetch();
    toast.success("Enterprise settings updated successfully");
  }

  useEffect(() => {
    if (getProjectResponse) {
      form.reset({
        logInWithSaml: getProjectResponse.project?.logInWithSaml ?? false,
        logInWithOidc: getProjectResponse.project?.logInWithOidc ?? false,
      });
    }
  }, [getProjectResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full">
          <Settings />
          Configure Enterprise Settings
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Enterprise Settings</DialogTitle>
          <DialogDescription>
            Configure whether users can log in with SAML/OIDC SSO.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="logInWithSaml"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with SAML</FormLabel>
                      <FormDescription>
                        Enable SAML SSO for your organization. This allows users
                        to log in using their SAML identity provider.
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
                name="logInWithOidc"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2">
                      <FormLabel>Log in with OIDC</FormLabel>
                      <FormDescription>
                        Enable OIDC SSO for your organization. This allows users
                        to log in using their OIDC identity provider.
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
            <DialogFooter className="mt-4">
              <Button variant="outline" onClick={() => setOpen(false)}>
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
                  : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
