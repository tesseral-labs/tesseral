import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  ArrowLeft,
  Ban,
  ChevronDown,
  LoaderCircle,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useLocation, useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { PageLoading } from "@/components/page/PageLoading";
import { Tab, Tabs } from "@/components/page/Tabs";
import { Title } from "@/components/page/Title";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
  deleteSCIMAPIKey,
  getSCIMAPIKey,
  revokeSCIMAPIKey,
  updateSCIMAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function OrganizationScimApiKeyDetailsTab() {
  const { scimApiKeyId } = useParams();

  const { data: getScimApiKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });

  const updateScimApiKeyMutation = useMutation(updateSCIMAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: getScimApiKeyResponse?.scimApiKey?.displayName || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
      scimApiKey: {
        displayName: data.displayName,
      },
    });
    await refetch();
    form.reset(data);
    toast.success("SCIM API Key updated successfully");
  }

  useEffect(() => {
    if (getScimApiKeyResponse?.scimApiKey) {
      form.reset({
        displayName: getScimApiKeyResponse.scimApiKey.displayName,
      });
    }
  }, [getScimApiKeyResponse, form]);

  return (
    <div className="space-y-8">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>SCIM API Key Details</CardTitle>
              <CardDescription>
                Update basic information about this SCIM API Key.
              </CardDescription>
              <CardAction>
                <Button
                  type="submit"
                  disabled={
                    !form.formState.isDirty ||
                    updateScimApiKeyMutation.isPending
                  }
                >
                  {updateScimApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {updateScimApiKeyMutation.isPending
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
                      The human-friendly name for this SCIM API Key.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        className="max-w-2xl"
                        placeholder="Display name"
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
  const { organizationId, scimApiKeyId } = useParams();
  const navigate = useNavigate();

  const { data: getScimApiKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  const revokeScimApiKeyMutation = useMutation(revokeSCIMAPIKey);
  const deleteScimApiKeyMutation = useMutation(deleteSCIMAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [revokeOpen, setRevokeOpen] = useState(false);

  async function handleRevoke() {
    await revokeScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });
    toast.success("SCIM API Key revoked successfully");
    setRevokeOpen(false);
    await refetch();
  }

  async function handleDelete() {
    await deleteScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });
    toast.success("SCIM API Key deleted successfully");
    navigate(`/organizations/${organizationId}/authentication`);
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <TriangleAlert className="h-4 w-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Ban className="w-4 h-4" />
                <span>Revoke SCIM API Key</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Revoke the SCIM API Key. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setRevokeOpen(true)}
              disabled={getScimApiKeyResponse?.scimApiKey?.revoked}
            >
              Revoke SCIM API Key
            </Button>
          </div>
          <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Trash className="w-4 h-4" />
                <span>Delete SCIM API Key</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Completely delete the SCIM API Key. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
              disabled={!getScimApiKeyResponse?.scimApiKey?.revoked}
            >
              Delete SCIM API Key
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Revoke Confirmation Dialog */}
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently revoke the{" "}
              <span className="font-semibold">
                {getScimApiKeyResponse?.scimApiKey?.displayName ||
                  getScimApiKeyResponse?.scimApiKey?.id}
              </span>{" "}
              SCIM API Key. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setRevokeOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleRevoke}>
              Revoke SCIM API Key
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
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete the{" "}
              <span className="font-semibold">
                {getScimApiKeyResponse?.scimApiKey?.displayName ||
                  getScimApiKeyResponse?.scimApiKey?.id}
              </span>{" "}
              SCIM API Key. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete SCIM API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
