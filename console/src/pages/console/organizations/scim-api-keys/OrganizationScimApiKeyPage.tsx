import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
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
  deleteSCIMAPIKey,
  getSCIMAPIKey,
  updateSCIMAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

export function OrganizationScimApiKeyPage() {
  const { organizationId, scimApiKeyId } = useParams();

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
    <PageContent>
      <div>
        <Link to={`/organizations/${organizationId}/authentication`}>
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to Authentication
          </Button>
        </Link>
      </div>

      <div>
        <div>
          <h1 className="text-2xl font-semibold">
            {getScimApiKeyResponse?.scimApiKey?.displayName || "SCIM API Key"}
          </h1>
          <ValueCopier
            value={getScimApiKeyResponse?.scimApiKey?.id || ""}
            label="SCIM API Key ID"
          />
          <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
            {getScimApiKeyResponse?.scimApiKey?.revoked ? (
              <Badge variant="secondary">Revoked</Badge>
            ) : (
              <Badge>Active</Badge>
            )}
            <Badge className="border-0" variant="outline">
              Created{" "}
              {getScimApiKeyResponse?.scimApiKey?.createTime &&
                DateTime.fromJSDate(
                  timestampDate(getScimApiKeyResponse.scimApiKey.createTime),
                ).toRelative()}
            </Badge>
            <div>•</div>
            <Badge className="border-0" variant="outline">
              Updated{" "}
              {getScimApiKeyResponse?.scimApiKey?.updateTime &&
                DateTime.fromJSDate(
                  timestampDate(getScimApiKeyResponse.scimApiKey.updateTime),
                ).toRelative()}
            </Badge>
          </div>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>SCIM API Key Details</CardTitle>
              <CardDescription>
                Update basic information about this SCIM API Key.
              </CardDescription>
              <CardAction>
                <Button type="submit" disabled={!form.formState.isDirty}>
                  Save changes
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
    </PageContent>
  );
}

function DangerZoneCard() {
  const { organizationId, scimApiKeyId } = useParams();
  const navigate = useNavigate();

  const { data: getScimApiKeyResponse } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  const deleteScimApiKeyMutation = useMutation(deleteSCIMAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteScimApiKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });
    toast.success("SCIM API Key deleted successfully");
    navigate(`/organizations/${organizationId}/authentication`);
  }

  return (
    <>
      <Card className="bg-red-50 border-red-20">
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
          <div className="flex items-center justify-between gap-8">
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
            >
              Delete SCIM API Key
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
