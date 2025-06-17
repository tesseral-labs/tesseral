import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, Trash, TriangleAlert } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
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
  deleteOrganization,
  getOrganization,
  updateOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function OrganizationDetailsTab() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse, refetch } = useQuery(getOrganization, {
    id: organizationId,
  });
  const updateOrganizationMutation = useMutation(updateOrganization);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: getOrganizationResponse?.organization?.displayName || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateOrganizationMutation.mutateAsync({
      id: organizationId,
      organization: {
        displayName: data.displayName,
      },
    });
    await refetch();
    form.reset(data);
    toast.success("Organization updated successfully");
  }

  useEffect(() => {
    if (getOrganizationResponse?.organization) {
      form.reset({
        displayName: getOrganizationResponse.organization.displayName,
      });
    }
  }, [getOrganizationResponse, form]);

  return (
    <div className="space-y-8">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>Organization details</CardTitle>
              <CardDescription>
                General settings for{" "}
                {getOrganizationResponse?.organization?.displayName}.
              </CardDescription>
              <CardAction>
                <Button
                  disabled={
                    !form.formState.isDirty ||
                    updateOrganizationMutation.isPending
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
              </CardAction>
            </CardHeader>
            <CardContent>
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display name</FormLabel>
                    <FormDescription>
                      The name shown to this Organization's users
                    </FormDescription>
                    <FormControl>
                      <Input
                        className="max-w-2xl"
                        placeholder="Display name"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
        </form>
      </Form>

      <DangerZoneCard />
    </div>
  );
}

function DangerZoneCard() {
  const { organizationId } = useParams();
  const navigate = useNavigate();

  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const deleteOrganizationMutation = useMutation(deleteOrganization);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteOrganizationMutation.mutateAsync({
      id: organizationId,
    });
    navigate("/organizations");
    toast.success("Organization deleted successfully");
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="text-destructive flex items-center gap-2">
            <TriangleAlert className="w-4 h-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-between gap-8">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Trash className="w-4 h-4" />
                <span>Delete Organization</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Completely delete the organization and all information
                associated with it. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete Organization
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Delete Confirmation Modal */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are your sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will completed delete the{" "}
              <span className="font-semibold">
                {getOrganizationResponse?.organization?.displayName}
              </span>{" "}
              Organization. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete Organization
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
