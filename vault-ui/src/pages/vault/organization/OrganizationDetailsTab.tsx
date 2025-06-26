import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
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
import { Input } from "@/components/ui/input";
import {
  getOrganization,
  updateOrganization,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function OrganizationDetailsTab() {
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization);
  const updateOrganizationMutation = useMutation(updateOrganization);

  const organization = getOrganizationResponse?.organization;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateOrganizationMutation.mutateAsync({
        organization: {
          displayName: data.displayName,
        },
      });
      await refetch();
      form.reset(data);
      toast.success("Organization details updated successfully.");
    } catch {
      toast.error("Failed to update organization details. Please try again.");
    }
  }

  useEffect(() => {
    if (organization) {
      form.reset({
        displayName: organization.displayName || "",
      });
    }
  }, [organization, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>Details</CardTitle>
            <CardDescription>
              Manage your organization settings.
            </CardDescription>
            <CardAction>
              <Button
                size="sm"
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  updateOrganizationMutation.isPending
                }
              >
                Save changes
              </Button>
            </CardAction>
          </CardHeader>
          <CardContent>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    This is the name that will be displayed for your
                    organization when users are logging in.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      className="max-w-lg"
                      {...field}
                      placeholder="ACME Corp"
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </CardContent>
        </Card>
      </form>
    </Form>
  );
}
