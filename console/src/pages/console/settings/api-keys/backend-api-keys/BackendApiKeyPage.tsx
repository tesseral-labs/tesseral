import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, LoaderCircle, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
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
  deleteBackendAPIKey,
  getBackendAPIKey,
  updateBackendAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function BackendApiKeyPage() {
  const { backendApiKeyId } = useParams();

  const { data: getBackendApiKeyResponse, refetch } = useQuery(
    getBackendAPIKey,
    {
      id: backendApiKeyId,
    },
  );
  const updateBackendApiKeyMutation = useMutation(updateBackendAPIKey);

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: getBackendApiKeyResponse?.backendApiKey?.displayName || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateBackendApiKeyMutation.mutateAsync({
      id: backendApiKeyId,
      backendApiKey: {
        displayName: data.displayName,
      },
    });
    await refetch();
    toast.success("Backend API Key updated successfully");
  }

  useEffect(() => {
    form.reset({
      displayName: getBackendApiKeyResponse?.backendApiKey?.displayName || "",
    });
  }, [getBackendApiKeyResponse, form]);

  return (
    <PageContent>
      <Title title={`Backend API Key ${backendApiKeyId}`} />

      <div>
        <Link to="/settings/api-keys">
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to API Keys
          </Button>
        </Link>
      </div>
      <div className="">
        <h1 className="text-2xl font-semibold">
          {getBackendApiKeyResponse?.backendApiKey?.displayName}
        </h1>
        <ValueCopier
          value={getBackendApiKeyResponse?.backendApiKey?.id || ""}
          label="API Key ID"
        />
        <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
          <Badge className="border-0" variant="outline">
            Created{" "}
            {getBackendApiKeyResponse?.backendApiKey?.createTime &&
              DateTime.fromJSDate(
                timestampDate(
                  getBackendApiKeyResponse.backendApiKey.createTime,
                ),
              ).toRelative()}
          </Badge>
          <div>•</div>
          <Badge className="border-0" variant="outline">
            Updated{" "}
            {getBackendApiKeyResponse?.backendApiKey?.updateTime &&
              DateTime.fromJSDate(
                timestampDate(
                  getBackendApiKeyResponse.backendApiKey.updateTime,
                ),
              ).toRelative()}
          </Badge>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card className="mt-4">
            <CardHeader>
              <CardTitle>Backend API Key Details</CardTitle>
              <CardDescription>
                Update the details of this Backend API Key.
              </CardDescription>
              <CardAction>
                <Button
                  type="submit"
                  disabled={
                    !form.formState.isDirty ||
                    updateBackendApiKeyMutation.isPending
                  }
                >
                  {updateBackendApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {updateBackendApiKeyMutation.isPending
                    ? "Saving Changes"
                    : "Save Changes"}
                </Button>
              </CardAction>
            </CardHeader>
            <CardContent className="space-y-6">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display name</FormLabel>
                    <FormDescription>
                      The human-friendly name for this Backend API Key.
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

      {/* <ListBackendApiKeyAuditLogsCard /> */}

      <DangerZoneCard />
    </PageContent>
  );
}

function DangerZoneCard() {
  const { backendApiKeyId } = useParams();
  const navigate = useNavigate();

  const deleteBackendApiKeyMutation = useMutation(deleteBackendAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteBackendApiKeyMutation.mutateAsync({
      id: backendApiKeyId,
    });
    toast.success("Backend API Key deleted successfully");
    navigate("/settings/api-keys");
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
                <span>Delete Backend API Key</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Delete this Backend API Key. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete Backend API Key
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete this Backend API Key. This action
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete Backend API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
