import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Users } from "lucide-react";
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
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function OrganizationScimCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Users />
          SCIM
        </CardTitle>
        <CardDescription>
          Configure SCIM user provisioning for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}
          </span>
          .
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex justify-between items-center gap-4">
            <div>
              <div className="font-semibold text-sm">SCIM Enable</div>
              <div className="text-xs text-muted-foreground">
                Allows automatic user management through SCIM-compatible
                identity providers.
              </div>
            </div>
            <Switch
              checked={getOrganizationResponse?.organization?.scimEnabled}
              disabled
            />
          </div>

          <ConfigureOrganizationScim />
        </div>
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  scimEnabled: z.boolean(),
});

function ConfigureOrganizationScim() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      scimEnabled: getOrganizationResponse?.organization?.scimEnabled ?? false,
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
        scimEnabled: data.scimEnabled,
      },
    });
    refetch();
    form.reset(data);
    setOpen(false);
    toast.success("SCIM configuration updated successfully.");
  }

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        scimEnabled: getOrganizationResponse.organization.scimEnabled || false,
      });
    }
  }, [getOrganizationResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="w-full" variant="outline">
          <Users />
          Configure SCIM
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Configure SCIM</DialogTitle>
          <DialogDescription>
            Configure SCIM user provisioning for{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            .
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="scimEnabled"
              render={({ field }) => (
                <FormItem className="flex gap-4">
                  <div>
                    <FormLabel>SCIM Enabled</FormLabel>
                    <FormDescription>
                      Allows automatic user management through SCIM-compatible
                      identity providers.
                    </FormDescription>
                    <FormMessage />
                  </div>
                  <FormControl className="flex items-center space-x-4">
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <DialogFooter className="justify-end mt-8 gap-2">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button disabled={!form.formState.isDirty} type="submit">
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
