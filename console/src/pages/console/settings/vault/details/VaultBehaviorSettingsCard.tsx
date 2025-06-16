import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Workflow } from "lucide-react";
import React, { MouseEvent, useEffect, useState } from "react";
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
  getProjectUISettings,
  updateProjectUISettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultBehaviorSettingsCard() {
  const { data: getProjectUiSettingsResponse } = useQuery(getProjectUISettings);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Workflow />
          <span>Behavior Settings</span>
        </CardTitle>
        <CardDescription>Configure how the Vault behaves.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6 flex-grow">
        <div className="flex justify-between items-center gap-2">
          <div className="text-sm space-y-1">
            <p className="font-semibold">Auto-create Organizations</p>
            <p className="text-muted-foreground">
              Create Organizations automatically when a new user signs up.
            </p>
          </div>
          <Switch
            disabled
            checked={
              getProjectUiSettingsResponse?.projectUiSettings
                ?.autoCreateOrganizations
            }
          />
        </div>
      </CardContent>
      <CardFooter>
        <ConfigureVaultBehaviorSettingsButton />
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  autoCreateOrganizations: z.boolean(),
});

function ConfigureVaultBehaviorSettingsButton() {
  const { data: getProjectUiSettingsResponse, refetch } =
    useQuery(getProjectUISettings);
  const updateProjectUiSettingsMutation = useMutation(updateProjectUISettings);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      autoCreateOrganizations:
        getProjectUiSettingsResponse?.projectUiSettings
          ?.autoCreateOrganizations ?? false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateProjectUiSettingsMutation.mutateAsync({
      autoCreateOrganizations: data.autoCreateOrganizations,
    });
    await refetch();
    form.reset(data);
    setOpen(false);
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
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          className="w-full"
          size="lg"
          onClick={() => setOpen(true)}
        >
          Configure Vault Behavior
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure Vault Behavior</DialogTitle>
          <DialogDescription>
            Adjust the behavior settings for the Vault.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
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
            </div>

            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button type="submit" disabled={!form.formState.isDirty}>
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
