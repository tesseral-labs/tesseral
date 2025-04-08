import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { clsx } from "clsx";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Outlet } from "react-router";
import { Link, useLocation } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Card,
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

export function OrganizationSettingsPage() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);

  const tabs = [
    {
      root: true,
      name: "Users",
      url: `/organization-settings`,
    },
    {
      name: "Advanced Settings",
      url: `/organization-settings/advanced`,
    },
  ];

  const { pathname } = useLocation();
  const currentTab = tabs.find((tab) => tab.url === pathname);

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>Organization Settings</CardTitle>
          <CardDescription>Manage your organization settings.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-medium">Organization Name</div>
                <div className="text-sm">
                  {getOrganizationResponse?.organization?.displayName}
                </div>
              </div>

              <EditOrganizationNameButton />
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="flex">
        {tabs.map((tab) => (
          <Link
            key={tab.name}
            to={tab.url}
            className={clsx(
              tab.url === currentTab?.url
                ? "border-foreground text-foreground"
                : "border-transparent text-muted-foreground hover:text-foreground",
              "whitespace-nowrap border-b-2 px-4 py-2 pb-3 text-sm font-medium",
            )}
          >
            {tab.name}
          </Link>
        ))}
      </div>

      <Outlet />
    </div>
  );
}

const schema = z.object({
  displayName: z.string().nonempty(),
});

function EditOrganizationNameButton() {
  const { data: getOrganizationResponse, refetch: refetchOrganization } =
    useQuery(getOrganization);

  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        displayName: getOrganizationResponse.organization.displayName,
      });
    }
  }, [form, getOrganizationResponse]);

  const { mutateAsync: updateOrganizationAsync } =
    useMutation(updateOrganization);

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateOrganizationAsync({
      organization: {
        displayName: values.displayName,
      },
    });
    await refetchOrganization();
    setOpen(false);
    toast.success(`Organization name updated.`);
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Organization Name</AlertDialogTitle>
          <AlertDialogDescription>
            Update the name of your organization.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-8"
          >
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Organization Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    The name of your organization.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Update</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
