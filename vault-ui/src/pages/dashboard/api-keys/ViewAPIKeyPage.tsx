import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
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
  getOrganization,
  listAPIKeyRoleAssignments,
  revokeAPIKey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

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
                {getAPIKeyResponse?.apiKey?.expireTime &&
                  DateTime.fromJSDate(
                    timestampDate(getAPIKeyResponse?.apiKey?.expireTime),
                  ).toRelative()}
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
          {/* <DetailsGrid>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Created at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.createTime &&
                        DateTime.fromJSDate(
                          timestampDate(getAPIKeyResponse?.apiKey?.createTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Updated at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.updateTime &&
                        DateTime.fromJSDate(
                          timestampDate(getAPIKeyResponse?.apiKey?.updateTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Status</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.revoked
                        ? 'Revoked'
                        : 'Active'}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Expires at</DetailsGridKey>
                    <DetailsGridValue>
                      {getAPIKeyResponse?.apiKey?.expireTime ? (
                        <>
                          {DateTime.fromJSDate(
                            timestampDate(
                              getAPIKeyResponse?.apiKey?.expireTime,
                            ),
                          ).toRelative()}
                        </>
                      ) : (
                        <>{'never'}</>
                      )}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
              </DetailsGrid> */}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="py-4 flex flex-row items-center justify-between">
          <div>
            <CardTitle>API Key Roles</CardTitle>
            <CardDescription>
              Manage the roles associated with this API key.
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
                  <TableRow key={roleAssignment.id}>
                    <TableCell>{roleAssignment.role?.displayName}</TableCell>
                    <TableCell className="space-x-2">
                      {roleAssignment.role?.actions.map((action) => (
                        <span
                          key={action}
                          className="p-1 text-xs text-mono bg-muted text-muted-foreground rounded"
                        >
                          {action}
                        </span>
                      ))}
                    </TableCell>
                    <TableCell className="text-right">
                      <RemoveRoleButton id={roleAssignment.id} />
                    </TableCell>
                  </TableRow>
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

const schema = z.object({
  displayName: z.string().min(1, { message: "Display name is required" }),
});

const EditAPIKeyButton = () => {
  const { organizationId, apiKeyId } = useParams();
  const { data: getAPIKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
    organizationId,
  });

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      displayName: "",
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {};

  useEffect(() => {
    if (getAPIKeyResponse?.apiKey) {
      form.reset({
        displayName: getAPIKeyResponse.apiKey.displayName,
      });
    }
  }, [getAPIKeyResponse]);

  return (
    <AlertDialog>
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
};

const RemoveRoleButton = ({ id }: { id: string }) => {
  const { apiKeyId } = useParams();
  const [open, setOpen] = useState(false);
  const { refetch } = useQuery(listAPIKeyRoleAssignments, {
    apiKeyId,
  });
  const deleteAPIKeyRoleAssignmentMutation = useMutation(
    deleteAPIKeyRoleAssignment,
  );

  const handleDelete = async () => {
    await deleteAPIKeyRoleAssignmentMutation.mutateAsync({
      id,
    });

    toast.success("Role removed successfully");
    await refetch();
    setOpen(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="destructive">Remove</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remove Role</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to remove this role from the API key? This
            disable all actions associated with this role for this API key.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button variant="destructive" onClick={() => handleDelete()}>
            Remove
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

const DangerZoneCard = () => {
  const navigate = useNavigate();
  const { organizationId, apiKeyId } = useParams();
  const { data: getAPIKeyResponse, refetch } = useQuery(getAPIKey, {
    id: apiKeyId,
    organizationId,
  });
  const deleteAPIKeyMutation = useMutation(deleteAPIKey);
  const revokeAPIKeyMutation = useMutation(revokeAPIKey);

  const [revokeOpen, setRevokeOpen] = useState(false);

  const handleDelete = async () => {
    await deleteAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success("API key deleted successfully");

    navigate(`/organizations/${organizationId}/api-keys`);
  };

  const handleRevoke = async () => {
    await revokeAPIKeyMutation.mutateAsync({
      id: apiKeyId,
    });

    toast.success("API key revoked successfully");

    await refetch();
  };

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
          {getAPIKeyResponse?.apiKey?.secretTokenSuffix ? (
            <div className="space-y-2">
              <div className="font-semibold text-sm">Revoke this API Key</div>
              <div className="text-sm text-muted-foreground">
                This action cannot be undone. The{" "}
                <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will no
                longer be usable, but all database entries will be retained.
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
                      will no longer be usable, but all database entries will be
                      retained.
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
            <div className="space-y-2">
              <div className="font-semibold text-sm">Delete this API Key</div>
              <div className="text-sm text-muted-foreground">
                This action cannot be undone. The{" "}
                <b>{getAPIKeyResponse?.apiKey?.displayName}</b> API key will no
                longer be usable and all database entries will be permanently
                deleted.
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
                      <b>{getAPIKeyResponse?.apiKey?.displayName}</b> will no
                      longer be usable and all database entries will be
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
};
