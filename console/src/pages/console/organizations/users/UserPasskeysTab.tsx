import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import {
  AlignLeft,
  Settings,
  ShieldBan,
  ShieldPlus,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { Link, useParams } from "react-router";
import { toast } from "sonner";

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
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  deletePasskey,
  getUser,
  listPasskeys,
  updatePasskey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { Passkey } from "@/gen/tesseral/backend/v1/models_pb";
import { AAGUIDS } from "@/lib/passkeys";

export function UserPasskeysTab() {
  const { organizationId, userId } = useParams();

  const {
    data: listPasskeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listPasskeys,
    {
      userId: userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const passkeys =
    listPasskeysResponses?.pages?.flatMap((page) => page.passkeys) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Passkeys</CardTitle>
        <CardDescription>Passkeys associated with this User.</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={6} />
        ) : (
          <>
            {passkeys.length === 0 && (
              <div className="text-center text-muted-foreground py-6 text-sm">
                No passkeys found for this user
              </div>
            )}
            {passkeys.length > 0 && (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Passkey</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Public Key</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Updated</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {passkeys.map((passkey) => (
                    <TableRow key={passkey.id}>
                      <TableCell className="space-y-2">
                        <Link
                          to={`/organizations/${organizationId}/users/${userId}/passkeys/${passkey.id}`}
                        >
                          <span className="block font-semibold">
                            {passkey.aaguid && AAGUIDS[passkey.aaguid]
                              ? AAGUIDS[passkey.aaguid]
                              : "Unknown Vendor"}
                          </span>
                        </Link>
                        <ValueCopier value={passkey.id} label="Passkey ID" />
                      </TableCell>
                      <TableCell>
                        {passkey.disabled ? (
                          <Badge variant="secondary">Disabled</Badge>
                        ) : (
                          <Badge>Active</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {passkey.publicKeyPkix && (
                          <a
                            className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                            download={`Public Key ${passkey.id}.pem`}
                            href={`data:text/plain;base64,${btoa(passkey.publicKeyPkix)}`}
                          >
                            Download (.pem)
                          </a>
                        )}
                      </TableCell>
                      <TableCell>
                        {passkey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(passkey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {passkey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(passkey.updateTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManagePasskeyButton passkey={passkey} />
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
        <CardFooter className="justify-center">
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

function ManagePasskeyButton({ passkey }: { passkey: Passkey }) {
  const { organizationId, userId } = useParams();

  const { refetch } = useInfiniteQuery(
    listPasskeys,
    {
      userId: userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const deletePasskeyMutation = useMutation(deletePasskey);
  const updatePasskeyMutation = useMutation(updatePasskey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [enableOpen, setEnableOpen] = useState(false);

  async function handleDelete() {
    await deletePasskeyMutation.mutateAsync({
      id: passkey.id,
    });
    await refetch();
    setDeleteOpen(false);
    toast.success("Passkey deleted successfully.");
  }

  async function handleDisable() {
    await updatePasskeyMutation.mutateAsync({
      id: passkey.id,
      passkey: {
        disabled: true,
      },
    });
    await refetch();
    setDisableOpen(false);
    toast.success(`Passkey disabled successfully.`);
  }

  async function handleEnable() {
    await updatePasskeyMutation.mutateAsync({
      id: passkey.id,
      passkey: {
        disabled: false,
      },
    });
    await refetch();
    setEnableOpen(false);
    toast.success(`Passkey enabled successfully.`);
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
              to={`/organizations/${organizationId}/users/${userId}/passkeys/${passkey.id}`}
            >
              <div className="w-full flex items-center gap-2">
                <AlignLeft />
                Details
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() =>
              passkey.disabled ? setEnableOpen(true) : setDisableOpen(true)
            }
          >
            {passkey.disabled ? (
              <>
                <ShieldPlus className="text-destructive group-hover:text-destructive" />
                <span className="text-destructive group-hover:text-destructive">
                  Enable Passkey
                </span>
              </>
            ) : (
              <>
                <ShieldBan className="text-destructive group-hover:text-destructive" />
                <span className="text-destructive group-hover:text-destructive">
                  Disable Passkey
                </span>
              </>
            )}
          </DropdownMenuItem>
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete Passkey
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Disable Confirmation Dialog */}
      <AlertDialog open={disableOpen} onOpenChange={setDisableOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will disable this Passkey, disallowing{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>
              to authentication with this Passkey on future logins.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="mt-8">
            <Button variant="outline" onClick={() => setDisableOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDisable}>
              Disable
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Enable Confirmation Dialog */}
      <AlertDialog open={enableOpen} onOpenChange={setEnableOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will enable this Passkey.{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>{" "}
              will be required to authenticate with this passkey (or another
              active passkey) when logging in.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="mt-8">
            <Button variant="outline" onClick={() => setEnableOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleEnable}>
              Enable
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
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This cannot be undone. This will permanently delete the{" "}
              <span className="font-semibold">{passkey?.id}</span> Passkey for{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>
              .
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
