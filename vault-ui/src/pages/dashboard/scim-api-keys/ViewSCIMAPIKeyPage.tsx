import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { get, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import { Link } from "react-router-dom";
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
import Loader from "@/components/ui/loader";
import {
  deleteSCIMAPIKey,
  getProject,
  getSCIMAPIKey,
  revokeSCIMAPIKey,
  updateSCIMAPIKey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ViewSCIMAPIKeyPage() {
  const { scimApiKeyId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getSCIMAPIKeyResponse } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });

  return (
    <div className="space-y-8">
      <Link to="/organization-settings/scim-api-keys">
        <Button variant="ghost" size="sm">
          <ChevronLeft className="h-4 w-4" />
          Back
        </Button>
      </Link>

      <Card>
        <CardHeader className="flex-row justify-between items-center space-x-4">
          <div className="space-y-2">
            <CardTitle>SCIM API Key</CardTitle>
            <CardDescription>View details of a SCIM API key.</CardDescription>
          </div>
          <EditSCIMAPIKeyButton />
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-8">
            <div className="py-8">
              <div className="text-sm font-semibold">Display Name</div>
              <p className="text-sm">
                {getSCIMAPIKeyResponse?.scimApiKey?.displayName}
              </p>

              <div className="text-sm font-semibold mt-4">SCIM Base URL</div>
              <p className="text-sm">
                {getProjectResponse?.project?.vaultDomain}/api/scim/v1
              </p>

              <div className="text-sm font-semibold mt-4">Status</div>
              <p className="text-sm">
                {getSCIMAPIKeyResponse?.scimApiKey?.revoked
                  ? "Revoked"
                  : "Active"}
              </p>
            </div>
            <div className="border-l pl-8 py-8">
              <div className="text-sm font-semibold mt-4">Created At</div>
              <p className="text-sm">
                {getSCIMAPIKeyResponse?.scimApiKey?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(
                      getSCIMAPIKeyResponse?.scimApiKey?.createTime,
                    ),
                  ).toRelative()}
              </p>

              <div className="text-sm font-semibold mt-4">Updated At</div>
              <p className="text-sm">
                {getSCIMAPIKeyResponse?.scimApiKey?.updateTime &&
                  DateTime.fromJSDate(
                    timestampDate(
                      getSCIMAPIKeyResponse?.scimApiKey?.updateTime,
                    ),
                  ).toRelative()}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

function EditSCIMAPIKeyButton() {
  const { scimApiKeyId } = useParams();
  const [open, setOpen] = useState(false);

  const { data: getSCIMAPIKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  const updateSCIMAPIKeyMutation = useMutation(updateSCIMAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(values: z.infer<typeof schema>) {
    await updateSCIMAPIKeyMutation.mutateAsync({
      id: scimApiKeyId,
      scimApiKey: {
        displayName: values.displayName,
      },
    });

    toast.success("SCIM API Key updated successfully");
    await refetch();
    form.reset();
    setOpen(false);
  }

  useEffect(() => {
    if (getSCIMAPIKeyResponse?.scimApiKey) {
      form.reset({
        displayName: getSCIMAPIKeyResponse.scimApiKey.displayName,
      });
    }
  }, [form, getSCIMAPIKeyResponse]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit SCIM API Key</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit SCIM API Key</AlertDialogTitle>
          <AlertDialogDescription>
            Manage this SCIM API Key.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    The human-friendly name for this SCIM API Key.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button
                type="submit"
                disabled={updateSCIMAPIKeyMutation.isPending}
              >
                {updateSCIMAPIKeyMutation.isPending && <Loader />}
                Update SCIM API Key
              </Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function DangerZoneCard() {
  const { scimApiKeyId } = useParams();
  const navigate = useNavigate();
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [revokeOpen, setRevokeOpen] = useState(false);

  const { data: getSCIMAPIKeyResponse, refetch } = useQuery(getSCIMAPIKey, {
    id: scimApiKeyId,
  });
  const deleteSCIMAPIKeyMutation = useMutation(deleteSCIMAPIKey);
  const recokeSCIMAPIKeyMutation = useMutation(revokeSCIMAPIKey);

  function handleDelete() {
    setDeleteOpen(true);
  }

  function handleRevoke() {
    setRevokeOpen(true);
  }

  async function handleConfirmDelete() {
    await deleteSCIMAPIKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });
    toast.success("SCIM API Key deleted successfully");
    setDeleteOpen(false);
    navigate(`/organization-settings/scim-api-keys`);
  }

  async function handleConfirmRevoke() {
    await recokeSCIMAPIKeyMutation.mutateAsync({
      id: scimApiKeyId,
    });
    toast.success("SCIM API Key revoked successfully");
    setRevokeOpen(false);
    await refetch();
  }

  return (
    <>
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Revoke SCIM API Key?</AlertDialogTitle>
            <AlertDialogDescription>
              Revoking a SCIM API Key cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmRevoke}>
              Revoke SCIM API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete SCIM API Key?</AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a SCIM API Key cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Permanently Delete SCIM API Key
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Card className="border-destructive">
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
          <CardDescription>
            Actions in this section are irreversible and can affect your
            organization.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-8">
            {!getSCIMAPIKeyResponse?.scimApiKey?.revoked && (
              <div className="flex justify-between items-center">
                <div>
                  <div className="text-sm font-semibold">
                    Revoke SCIM API Key
                  </div>
                  <p className="text-sm">
                    Revoke this SCIM API Key. This cannot be undone.
                  </p>
                </div>

                <Button variant="destructive" onClick={handleRevoke}>
                  Revoke SCIM API Key
                </Button>
              </div>
            )}

            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-semibold">Delete SCIM API Key</div>
                <p className="text-sm">
                  Delete this SCIM API Key. This cannot be undone.
                </p>
              </div>

              <Button
                disabled={!getSCIMAPIKeyResponse?.scimApiKey?.revoked}
                variant="destructive"
                onClick={handleDelete}
              >
                Delete SCIM API Key
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
}
