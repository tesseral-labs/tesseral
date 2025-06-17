import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Users } from "lucide-react";
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
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  scimEnabled: z.boolean(),
});

export function OrganizationScimCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      scimEnabled: getOrganizationResponse?.organization?.scimEnabled ?? false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        scimEnabled: data.scimEnabled,
      },
    });
    refetch();
    form.reset(data);
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
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
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
