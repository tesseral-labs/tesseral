import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Switch } from "@/components/ui/switch";
import {
  getOrganization,
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  customRolesEnabled: z.boolean(),
});

export function OrganizationRoleSettingsCard() {
  const { organizationId } = useParams();

  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      customRolesEnabled:
        getOrganizationResponse?.organization?.customRolesEnabled || false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        customRolesEnabled: data.customRolesEnabled,
      },
    });
    form.reset(data);
    await refetch();
    toast.success("Role settings updated successfully");
  }

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        customRolesEnabled:
          getOrganizationResponse.organization.customRolesEnabled || false,
      });
    }
  }, [form, getOrganizationResponse]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>Role Settings</CardTitle>
          </CardHeader>
          <CardContent>
            <FormField
              control={form.control}
              name="customRolesEnabled"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Allow Custom Roles</FormLabel>
                  <FormDescription>
                    Enables custom, Organization-level roles for this
                    Organization
                  </FormDescription>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormItem>
              )}
            />
          </CardContent>
          <CardFooter>
            <Button
              type="submit"
              disabled={
                !form.formState.isDirty || updateOrganizationMutation.isPending
              }
            >
              Save Changes
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
