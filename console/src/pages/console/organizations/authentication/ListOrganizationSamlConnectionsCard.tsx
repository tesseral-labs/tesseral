import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { AlignLeft, Plus, Settings, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";

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
  createSAMLConnection,
  deleteSAMLConnection,
  getOrganization,
  listSAMLConnections,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { SAMLConnection } from "@/gen/tesseral/backend/v1/models_pb";

export function ListOrganizationSamlConnectionsCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const {
    data: listSamlConnectionsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listSAMLConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const samlConnections =
    listSamlConnectionsResponses?.pages?.flatMap(
      (page) => page.samlConnections,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>SAML Connections</CardTitle>
        <CardDescription>
          Configure SAML identify providers for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}.
          </span>
        </CardDescription>
        <CardAction>
          <CreateSamlConnectionButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {samlConnections.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-6">
                No SAML connections found. Create a new connection to get
                started.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Connection</TableHead>
                    <TableHead>IDP Entity ID</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {samlConnections.map((samlConnection) => (
                    <TableRow key={samlConnection.id}>
                      <TableCell>
                        <span className="px-2 py-1 bg-muted font-mono text-xs rounded mr-2">
                          {samlConnection.id}
                        </span>
                        {samlConnection.primary && (
                          <Badge variant="outline">Primary</Badge>
                        )}
                      </TableCell>
                      <TableCell>{samlConnection.idpEntityId || "—"}</TableCell>
                      <TableCell>
                        {samlConnection.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(samlConnection.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManageSamlConnectionButton
                          samlConnection={samlConnection}
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

function CreateSamlConnectionButton() {
  const { organizationId } = useParams();
  const navigate = useNavigate();

  const { data: listSamlConnectionsResponse, refetch } = useInfiniteQuery(
    listSAMLConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createSamlConnectionMutation = useMutation(createSAMLConnection);

  async function handleCreateSamlConnection() {
    const isFirstSamlConnection =
      listSamlConnectionsResponse?.pages[0]?.samlConnections.length === 0;
    const newSamlConnection = await createSamlConnectionMutation.mutateAsync({
      samlConnection: {
        organizationId,
        primary: isFirstSamlConnection,
      },
    });

    if (!newSamlConnection.samlConnection?.id) {
      toast.error("Failed to create SAML Connection. Please try again.");
      return;
    }

    await refetch();
    toast.success("SAML Connection created successfully");
    navigate(
      `/organizations/${organizationId}/saml-connections/${newSamlConnection.samlConnection.id}`,
    );
  }

  return (
    <Button onClick={handleCreateSamlConnection}>
      <Plus />
      Create SAML Connection
    </Button>
  );
}

function ManageSamlConnectionButton({
  samlConnection,
}: {
  samlConnection: SAMLConnection;
}) {
  const { organizationId } = useParams();

  const { refetch } = useInfiniteQuery(
    listSAMLConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const deleteSamlConnectionMutation = useMutation(deleteSAMLConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteSamlConnectionMutation.mutateAsync({
      id: samlConnection.id,
    });
    await refetch();
    toast.success("SAML Connection deleted successfully");
    setDeleteOpen(false);
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
              to={`/organizations/${organizationId}/saml-connections/${samlConnection.id}`}
            >
              <div className="w-full flex items-center gap-2">
                <AlignLeft />
                <span>SAML Connection Details</span>
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete SAML Connection
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

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
              <span className="font-semibold">{samlConnection.id}</span> SAML
              Connection.
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
