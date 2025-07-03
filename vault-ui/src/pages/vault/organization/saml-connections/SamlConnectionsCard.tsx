import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import {
  AlignLeft,
  Plus,
  Settings,
  Trash,
  TriangleAlert,
  WandSparkles,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";

import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
import {
  AlertDialog,
  AlertDialogCancel,
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
  listSAMLConnections,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { SAMLConnection } from "@/gen/tesseral/frontend/v1/models_pb";
import { cn } from "@/lib/utils";

export function SamlConnectionsCard() {
  const {
    data: listSamlConnectionsResponses,
    hasNextPage,
    fetchNextPage,
    isFetchingNextPage,
    isLoading,
    refetch,
  } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const createSamlConnectionMutation = useMutation(createSAMLConnection);

  const samlConnections =
    listSamlConnectionsResponses?.pages.flatMap(
      (page) => page.samlConnections,
    ) || [];

  const [createdSamlConnection, setCreatedSamlConnection] =
    useState<SAMLConnection | null>(null);
  const navigate = useNavigate();

  async function handleCreate() {
    try {
      const { samlConnection } = await createSamlConnectionMutation.mutateAsync(
        {
          samlConnection: {},
        },
      );
      if (!samlConnection) {
        toast.error(
          "Failed to create SAML connection. Please try again later.",
        );
        return;
      }

      await refetch();
      navigate(`/organization/saml-connections/${samlConnection.id}/setup`);
      toast.success("SAML connection created successfully.");
    } catch {
      toast.error("Failed to create SAML connection. Please try again later.");
    }
  }

  useEffect(() => {
    if (createdSamlConnection) {
      setTimeout(() => {
        setCreatedSamlConnection(null);
      }, 2000);
    }
  }, [createdSamlConnection, setCreatedSamlConnection]);

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-col lg:flex-row items-start justify-between space-y-2 lg:space-y-0">
        <div className="space-y-2">
          <CardTitle>SAML Connections</CardTitle>
          <CardDescription>
            Configure SAML identify providers for your Organization.
          </CardDescription>
        </div>
        <CardAction className="w-full lg:w-auto">
          <Button
            onClick={handleCreate}
            disabled={createSamlConnectionMutation.isPending}
            size="sm"
          >
            <Plus className="mr-2" />
            Create SAML Connection
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {samlConnections.length > 0 ? (
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
                    <TableRow
                      key={samlConnection.id}
                      className={cn(
                        createdSamlConnection?.id === samlConnection.id
                          ? "bg-muted animate-pulse duration-500"
                          : "",
                      )}
                    >
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
            ) : (
              <div className="text-center text-muted-foreground text-sm pt-8">
                No SAML Connections found. Create one to get started.
              </div>
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
            Load more
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function ManageSamlConnectionButton({
  samlConnection,
}: {
  samlConnection: SAMLConnection;
}) {
  const { refetch } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deleteSamlConnectionMutation = useMutation(deleteSAMLConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete(samlConnection: SAMLConnection) {
    try {
      await deleteSamlConnectionMutation.mutateAsync({
        id: samlConnection.id,
      });

      await refetch();
      setDeleteOpen(false);
      toast.success("SAML connection deleted successfully.");
    } catch {
      toast.error("Failed to delete SAML connection. Please try again later.");
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
          <DropdownMenuItem asChild>
            <Link to={`/organization/saml-connections/${samlConnection.id}`}>
              <AlignLeft />
              Details
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <Link
              to={`/organization/saml-connections/${samlConnection.id}/setup`}
            >
              <WandSparkles />
              Setup Wizard
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Deleting this SAML Connection is permanent and connot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel asChild>
              <Button variant="outline">Cancel</Button>
            </AlertDialogCancel>
            <Button
              onClick={() => handleDelete(samlConnection)}
              variant="destructive"
            >
              Delete SAML Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
