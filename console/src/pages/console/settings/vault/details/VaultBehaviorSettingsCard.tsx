import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Workflow } from "lucide-react";
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
  getProjectUISettings,
  updateProject,
  updateProjectUISettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  auditLogsEnabled: z.boolean().optional(),
  autoCreateOrganizations: z.boolean(),
});

export function VaultBehaviorSettingsCard() {
  const { data: getProjectResponse, refetch: refetchProject } =
    useQuery(getProject);
  const {
    data: getProjectUiSettingsResponse,
    refetch: refetchProjectUiSettings,
  } = useQuery(getProjectUISettings);
  const updateProjectMutation = useMutation(updateProject);
  const updateProjectUiSettingsMutation = useMutation(updateProjectUISettings);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      auditLogsEnabled: getProjectResponse?.project?.auditLogsEnabled ?? false,
      autoCreateOrganizations:
        getProjectUiSettingsResponse?.projectUiSettings
          ?.autoCreateOrganizations ?? false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectUiSettingsMutation.mutateAsync({
      autoCreateOrganizations: data.autoCreateOrganizations,
    });
    await updateProjectMutation.mutateAsync({
      project: {
        auditLogsEnabled: data.auditLogsEnabled,
      },
    });
    await refetchProject();
    await refetchProjectUiSettings();
    form.reset(data);
    toast.success("Vault behavior settings updated successfully");
  }

  useEffect(() => {
    if (getProjectUiSettingsResponse && getProjectResponse && form) {
      form.reset({
        auditLogsEnabled:
          getProjectResponse?.project?.auditLogsEnabled ?? false,
        autoCreateOrganizations:
          getProjectUiSettingsResponse.projectUiSettings
            ?.autoCreateOrganizations ?? false,
      });
    }
  }, [getProjectUiSettingsResponse, getProjectResponse, form]);

  return (
    <Form {...form}>
      <form className="flex-grow" onSubmit={form.handleSubmit(handleSubmit)}>
        <Card className="h-full">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Workflow />
              <span>Behavior Settings</span>
            </CardTitle>
            <CardDescription>Configure how the Vault behaves.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 flex-grow">
            <FormField
              control={form.control}
              name="autoCreateOrganizations"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div className="space-y-2">
                    <FormLabel>Auto-create Organizations</FormLabel>
                    <FormDescription>
                      Automatically create Organizations when a new user signs
                      up. When disabled, new users will be prompted to name
                      their Organization upon first login.
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
              name="auditLogsEnabled"
              render={({ field }) => (
                <FormItem className="flex justify-between items-center gap-4">
                  <div className="space-y-2">
                    <FormLabel>Audit Logs</FormLabel>
                    <FormDescription>
                      Whether to show audit logs to your users in the Vault UI.
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
                !form.formState.isDirty ||
                updateProjectUiSettingsMutation.isPending
              }
              type="submit"
            >
              {updateProjectUiSettingsMutation.isPending && (
                <LoaderCircle className="animate-spin" />
              )}
              {updateProjectUiSettingsMutation.isPending
                ? "Saving changes"
                : "Save changes"}
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
