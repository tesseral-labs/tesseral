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
  getProjectUISettings,
  updateProjectUISettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  autoCreateOrganizations: z.boolean(),
});

export function VaultBehaviorSettingsCard() {
  const { data: getProjectUiSettingsResponse, refetch } =
    useQuery(getProjectUISettings);
  const updateProjectUiSettingsMutation = useMutation(updateProjectUISettings);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      autoCreateOrganizations:
        getProjectUiSettingsResponse?.projectUiSettings
          ?.autoCreateOrganizations ?? false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectUiSettingsMutation.mutateAsync({
      autoCreateOrganizations: data.autoCreateOrganizations,
    });
    await refetch();
    form.reset(data);
    toast.success("Vault behavior settings updated successfully");
  }

  useEffect(() => {
    if (getProjectUiSettingsResponse && form) {
      form.reset({
        autoCreateOrganizations:
          getProjectUiSettingsResponse.projectUiSettings
            ?.autoCreateOrganizations ?? false,
      });
    }
  }, [getProjectUiSettingsResponse, form]);

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
