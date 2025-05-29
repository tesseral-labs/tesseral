import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  deleteAPIKey,
  deleteAPIKeyRoleAssignment,
  getAPIKey,
  getRole,
  listAPIKeyRoleAssignments,
  revokeAPIKey,
  updateAPIKey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { APIKeyRoleAssignment } from "@/gen/tesseral/frontend/v1/models_pb";

import { AddAPIKeyRoleButton } from "./AddAPIKeyRoleButton";

export function ViewAPIKeyPage() {
  const { apiKeyId } = useParams();

  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const { data: listApiKeyRoleAssignmentsResponse } = useQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
    },
  );

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="py-4 flex flex-row items-center justify-between">
          <div>
            <CardTitle>API Key Details</CardTitle>
            <CardDescription></CardDescription>
          </div>

          <EditAPIKeyButton />
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-4 gap-4">
            <div>
              <div className="text-sm font-semibold">Display name</div>
              <div className="text-sm">
                {getAPIKeyResponse?.apiKey?.displayName}
              </div>
            </div>

            <div className="border-l pl-8">
              <div className="text-sm font-semibold">Status</div>
              <div className="text-sm">
                {getAPIKeyResponse?.apiKey?.revoked ? "Revoked" : "Active"}
              </div>
            </div>

            <div className="border-l pl-8">
              <div className="text-sm font-semibold">Expires</div>
              <div className="text-sm">
                {getAPIKeyResponse?.apiKey?.expireTime
                  ? DateTime.fromJSDate(
                      timestampDate(getAPIKeyResponse?.apiKey?.expireTime),
                    ).toRelative()
                  : "Never"}
              </div>
            </div>

            <div className="border-l pl-8">
              <div className="text-sm font-semibold">Created</div>
              <div className="text-sm">
                {getAPIKeyResponse?.apiKey?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(getAPIKeyResponse?.apiKey?.createTime),
                  ).toRelative()}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="py-4 flex flex-row items-center justify-between">
          <div>
            <CardTitle>API Key Roles</CardTitle>
            <CardDescription>
              Manage the roles assigned to this API key.
            </CardDescription>
          </div>

          <AddAPIKeyRoleButton />
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Role</TableHead>
                <TableHead>Actions</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {listApiKeyRoleAssignmentsResponse?.apiKeyRoleAssignments?.map(
                (roleAssignment) => (
                  <APIKeyRoleAssignmentRow
                    key={roleAssignment.id}
                    apiKeyRoleAssignment={roleAssignment}
                  />
                ),
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  );
}

function APIKeyRoleAssignmentRow({
  apiKeyRoleAssignment,
}: {
  apiKeyRoleAssignment: APIKeyRoleAssignment;
}) {
  const { data: getRoleResponse } = useQuery(getRole, {
    id: apiKeyRoleAssignment.roleId,
  });

  return (
    <TableRow key={getRoleResponse?.role?.id}>
      <TableCell>{getRoleResponse?.role?.displayName}</TableCell>
      <TableCell className="space-x-2">
        {getRoleResponse?.role?.actions.map((action) => (
          <span
            key={action}
            className="p-1 text-xs text-mono bg-muted text-muted-foreground rounded"
          >
            {action}
          </span>
        ))}
      </TableCell>
      <TableCell className="text-right">
        <RemoveRoleButton id={apiKeyRoleAssignment.id} />
      </TableCell>
    </TableRow>
  );
}

const schema = z.object({
  displayName: z.string().min(1, { message: "Display name is required" }),
});

function EditAPIKeyButton() {
  const [open, setOpen] = useState(false);

  const { apiKeyId } = useParams();
  const { data: getAPIKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });

  const updateAPIKeyMutation = useMutation(updateAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateAPIKeyMutation.mutateAsync({
      apiKey: {
        id: apiKeyId,
        displayName: data.displayName,
      },
    });

    await refetch();
    toast.success("API key updated successfully");
    setOpen(false);
  }

  useEffect(() => {
    if (getAPIKeyResponse?.apiKey) {
      form.reset({
        displayName: getAPIKeyResponse.apiKey.displayName,
      });
    }
  }, [getAPIKeyResponse, form]);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit API Key</AlertDialogTitle>
        </AlertDialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Display Name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function RemoveRoleButton({ id }: { id: string }) {
  const { apiKeyId } = useParams();
  const [open, setOpen] = useState(false);
  const { refetch } = useQuery(listAPIKeyRoleAssignments, {
    apiKeyId,
  });
  const deleteAPIKeyRoleAssignmentMutation = useMutation(
    deleteAPIKeyRoleAssignment,
  );

  async function handleDelete() {
    await deleteAPIKeyRoleAssignmentMutation.mutateAsync({
      id,
    });

    toast.success("Role removed successfully");
    await refetch();
    setOpen(false);
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="destructive">Unassign</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Unassign Role</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to unassign this role from the API key? This
            disable all actions associated with this role for this API key.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button variant="destructive" onClick={() => handleDelete()}>
            Unassign
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function DangerZoneCard() {
  const navigate = useNavigate();
  const { apiKeyId } = useParams();
  const { data: getAPIKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const deleteAPIKeyMutation = useMutation(deleteAPIKey);
  const revokeAPIKeyMutation = useMutation(revokeAPIKey);

  const [revokeOpen, setRevokeOpen] = useState(false);

  async function handleDelete() {
    await deleteAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success("API key deleted successfully");

    navigate(`/organization-settings/api-keys`);
  }

  async function handleRevoke() {
    await revokeAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success("API key revoked successfully");

    await refetch();
  }

  return (
    <Card className="border-destructive">
      <CardHeader>
        <CardTitle>Danger Zone</CardTitle>
        <CardDescription>
          Actions in this section cannot be undone. Please proceed with caution.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-8">
          {!getAPIKeyResponse?.apiKey?.revoked ? (
            <div className="flex justify-between items-center space-y-2">
              <div>
                <div className="font-semibold text-sm">Revoke this API Key</div>
                <div className="text-sm text-muted-foreground">
                  This action cannot be undone. The{" "}
                  <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will
                  no longer be usable.
                </div>
              </div>
              <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Revoke API Key</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <p>
                      This action cannot be undone. The{" "}
                      <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key
                      will no longer be usable.
                    </p>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <Button variant="destructive" onClick={handleRevoke}>
                      Revoke API Key
                    </Button>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          ) : (
            <div className="flex justify-between items-center space-y-2">
              <div>
                <div className="font-semibold text-sm">Delete this API Key</div>
                <div className="text-sm text-muted-foreground">
                  This action cannot be undone. The{" "}
                  <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will
                  be permanently deleted.
                </div>
              </div>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Delete API Key</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <p>
                      This action cannot be undone.{" "}
                      <b>{getAPIKeyResponse?.apiKey?.displayName}</b> will be
                      permanently deleted.
                    </p>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <Button variant="destructive" onClick={handleDelete}>
                      Delete API Key
                    </Button>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
