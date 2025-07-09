import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlignLeft,
  Ban,
  LoaderCircle,
  Plus,
  Settings,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { SecretCopier } from "@/components/core/SecretCopier";
import { ValueCopier } from "@/components/core/ValueCopier";
import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
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
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createSCIMAPIKey,
  deleteSCIMAPIKey,
  listSCIMAPIKeys,
  revokeSCIMAPIKey,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { SCIMAPIKey } from "@/gen/tesseral/frontend/v1/models_pb";

export function ScimApiKeysCard() {
  const {
    data: listScimApiKeysResponses,
    isLoading,
    hasNextPage,
    fetchNextPage,
    isFetchingNextPage,
    refetch,
  } = useInfiniteQuery(
    listSCIMAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const scimApiKeys =
    listScimApiKeysResponses?.pages.flatMap((page) => page.scimApiKeys) || [];

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>SCIM API Keys</CardTitle>
        <CardDescription>
          A SCIM API key lets you do enterprise directory syncing.
        </CardDescription>
        <CardAction>
          <CreateScimApiKeyButton refetch={refetch} />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {scimApiKeys?.length === 0 ? (
              <div className="text-muted-foreground text-sm text-center pt-8">
                No SCIM API keys found for this organization.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>API Key</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Updated</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {scimApiKeys?.map((scimApiKey) => (
                    <TableRow key={scimApiKey.id}>
                      <TableCell className="space-y-1">
                        <div className="font-medium">
                          <Link
                            to={`/organization/scim-api-keys/${scimApiKey.id}`}
                          >
                            {scimApiKey.displayName || "No display name"}
                          </Link>
                        </div>

                        <ValueCopier
                          value={scimApiKey.id}
                          label="SCIM API Key ID"
                        />
                      </TableCell>
                      <TableCell>
                        {scimApiKey.revoked ? (
                          <Badge variant="secondary">Revoked</Badge>
                        ) : (
                          <Badge>Active</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {scimApiKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(scimApiKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {scimApiKey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(scimApiKey.updateTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManageScimApiKeyButton
                          scimApiKey={scimApiKey}
                          refetch={refetch}
                        />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </>
        )}
      </CardContent>
      {hasNextPage && (
        <CardFooter className="mt-4 justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            Load More
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

function CreateScimApiKeyButton({ refetch }: { refetch: () => Promise<any> }) {
  const createScimApiKeyMutation = useMutation(createSCIMAPIKey);

  const [createOpen, setCreateOpen] = useState(false);
  const [secretOpen, setSecretOpen] = useState(false);
  const [scimApiKey, setScimApiKey] = useState<SCIMAPIKey>();

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setCreateOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    const { scimApiKey } = await createScimApiKeyMutation.mutateAsync({
      scimApiKey: {
        displayName: data.displayName,
      },
    });

    if (!scimApiKey) {
      toast.error("Failed to create SCIM API key. Please try again.");
      setCreateOpen(false);
      return;
    }

    setScimApiKey(scimApiKey);
    setSecretOpen(true);
    setCreateOpen(false);

    toast.success("SCIM API key created successfully");
    form.reset();

    await refetch();
  }

  return (
    <>
      <Dialog
        open={!!scimApiKey?.secretToken && secretOpen}
        onOpenChange={setSecretOpen}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>API Key Created</DialogTitle>
            <DialogDescription>
              API Key was created successfully.
            </DialogDescription>
          </DialogHeader>

          <div className="text-sm font-medium leading-none">
            API Key Secret Token
          </div>

          {scimApiKey?.secretToken && (
            <SecretCopier
              placeholder={`tesseral_secret_scim_api_key_•••••••••••••••••••••••••••••••••••••••••••••••••••••••`}
              secret={scimApiKey.secretToken}
            />
          )}

          <div className="text-sm text-muted-foreground">
            Store this secret in your secrets manager. You will not be able to
            see this secret token again later.
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setSecretOpen(false)}>
              Close
            </Button>
            {!!scimApiKey?.id && (
              <Link to={`/organization/scim-api-keys/${scimApiKey.id}`}>
                <Button>View API Key</Button>
              </Link>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogTrigger asChild>
          <Button size="sm">
            <Plus />
            Create SCIM API Key
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create SCIM API Key</DialogTitle>
            <DialogDescription>
              Create a new SCIM API key for this organization.
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <div className="space-y-6">
                <FormField
                  control={form.control}
                  name="displayName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Display Name</FormLabel>
                      <FormDescription>
                        The human-friendly name for this SCIM API key.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input placeholder="Display name" {...field} />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
              <DialogFooter className="mt-8">
                <Button variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button
                  disabled={
                    !form.formState.isDirty ||
                    createScimApiKeyMutation.isPending
                  }
                  type="submit"
                >
                  {createScimApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {createScimApiKeyMutation.isPending
                    ? "Creating SCIM API Key"
                    : "Create SCIM API Key"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </>
  );
}

function ManageScimApiKeyButton({
  scimApiKey,
  refetch,
}: {
  scimApiKey: SCIMAPIKey;
  refetch: () => Promise<any>;
}) {
  const revokeScimApiKeyMutation = useMutation(revokeSCIMAPIKey);
  const deleteScimApiKeyMutation = useMutation(deleteSCIMAPIKey);

  const [revokeOpen, setRevokeOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleRevoke() {
    try {
      await revokeScimApiKeyMutation.mutateAsync({
        id: scimApiKey.id,
      });
      await refetch();
      setRevokeOpen(false);
      toast.success("SCIM API key revoked successfully");
    } catch {
      toast.error("Failed to revoke SCIM API key. Please try again.");
      setRevokeOpen(false);
    }
  }

  async function handleDelete() {
    try {
      await deleteScimApiKeyMutation.mutateAsync({
        id: scimApiKey.id,
      });
      await refetch();
      setDeleteOpen(false);
      toast.success("SCIM API key deleted successfully");
    } catch {
      toast.error("Failed to delete SCIM API key. Please try again.");
      setDeleteOpen(false);
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm">
            <Settings />
            Manage
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organization/scim-api-keys/${scimApiKey.id}`}
            >
              <div className="w-full flex items-center gap-2">
                <AlignLeft />
                <span>Details</span>
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => setRevokeOpen(true)}
            disabled={scimApiKey.revoked}
          >
            <Ban className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Revoke SCIM API Key
            </span>
          </DropdownMenuItem>
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
            disabled={!scimApiKey.revoked}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete SCIM API Key
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Revoke Confirmation Dialog */}
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert className="w-4 h-4" />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently revoke the{" "}
              <span className="font-semibold">
                {scimApiKey.displayName || scimApiKey.id}
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
              <TriangleAlert className="w-4 h-4" />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete the{" "}
              <span className="font-semibold">
                {scimApiKey.displayName || scimApiKey.id}
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
