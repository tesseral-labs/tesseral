import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoaderCircle, ShieldBan, Trash, TriangleAlert } from "lucide-react";
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
  deleteAPIKey,
  getAPIKey,
  revokeAPIKey,
  updateAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function OrganizationApiKeyDetailsTab() {
  const { apiKeyId } = useParams();

  const { data: getApiKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const updateApiKeyMutation = useMutation(updateAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateApiKeyMutation.mutateAsync({
      id: apiKeyId,
      apiKey: {
        displayName: data.displayName,
      },
    });
    await refetch();
    form.reset(data);
    toast.success("API Key updated successfully.");
  }

  useEffect(() => {
    if (getApiKeyResponse?.apiKey) {
      form.reset({
        displayName: getApiKeyResponse.apiKey.displayName || "",
      });
    }
  }, [getApiKeyResponse, form]);

  return (
    <div className="space-y-8">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>API Key Details</CardTitle>
              <CardDescription>
                View and manage the details of your API key.
              </CardDescription>
              <CardAction>
                <Button
                  disabled={
                    !form.formState.isDirty || updateApiKeyMutation.isPending
                  }
                  type="submit"
                >
                  {updateApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {updateApiKeyMutation.isPending
                    ? "Saving changes"
                    : "Save changes"}
                </Button>
              </CardAction>
            </CardHeader>
            <CardContent className="space-y-6">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      The human-friendly name for this API Key.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        className="max-w-2xl"
                        placeholder="Enter display name"
                        {...field}
                      />
                    </FormControl>
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
  const { apiKeyId, organizationId } = useParams();
  const navigate = useNavigate();

  const { data: getApiKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });

  const deleteApiKeyMutation = useMutation(deleteAPIKey);
  const revokeApiKeyMutation = useMutation(revokeAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [revokeOpen, setRevokeOpen] = useState(false);

  async function handleDelete() {
    await deleteApiKeyMutation.mutateAsync({ id: apiKeyId });
    toast.success("API Key deleted successfully.");
    navigate(`/organizations/${organizationId}/api-keys`);
  }

  async function handleRevoke() {
    await revokeApiKeyMutation.mutateAsync({ id: apiKeyId });
    await refetch();
    setRevokeOpen(false);
    toast.success("API Key revoked successfully.");
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <TriangleAlert className="w-4 h-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {getApiKeyResponse?.apiKey?.revoked ? (
            <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
              <div className="space-y-1">
                <div className="text-sm font-semibold flex items-center gap-2">
                  <Trash className="w-4 h-4" />
                  <span>Delete API Key</span>
                </div>
                <div className="text-sm text-muted-foreground">
                  Permanently delete the API Key. This cannot be undone.
                </div>
              </div>
              <Button
                variant="destructive"
                size="sm"
                onClick={() => setDeleteOpen(true)}
              >
                Delete API Key
              </Button>
            </div>
          ) : (
            <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
              <div className="space-y-1">
                <div className="text-sm font-semibold flex items-center gap-2">
                  <ShieldBan className="w-4 h-4" />
                  <span>Revoke API Key</span>
                </div>
                <div className="text-sm text-muted-foreground">
                  Revoke the API Key. This cannot be undone.
                </div>
              </div>
              <Button
                className="border-destructive text-destructive hover:bg-destructive hover:text-white"
                variant="outline"
                size="sm"
                onClick={() => setRevokeOpen(true)}
              >
                Revoke API Key
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Revoke Confirmation Dialog */}
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are your sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will revoke the{" "}
              <span className="font-semibold">
                {getApiKeyResponse?.apiKey?.displayName ||
                  getApiKeyResponse?.apiKey?.id}
              </span>{" "}
              API Key. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setRevokeOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleRevoke}>
              Revoke API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are your sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete the{" "}
              <span className="font-semibold">
                {getApiKeyResponse?.apiKey?.displayName ||
                  getApiKeyResponse?.apiKey?.id}
              </span>{" "}
              API Key. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
