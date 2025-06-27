import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { AlignLeft, Plus, Settings, Trash, TriangleAlert } from "lucide-react";
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
  createOIDCConnection,
  deleteOIDCConnection,
  listOIDCConnections,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { OIDCConnection } from "@/gen/tesseral/frontend/v1/models_pb";
import { cn } from "@/lib/utils";

export function OidcConnectionsCard() {
  const {
    data: listOidcConnectionsResponses,
    hasNextPage,
    fetchNextPage,
    isFetchingNextPage,
    isLoading,
    refetch,
  } = useInfiniteQuery(
    listOIDCConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const createOidcConnectionMutation = useMutation(createOIDCConnection);

  const oidcConnections =
    listOidcConnectionsResponses?.pages.flatMap(
      (page) => page.oidcConnections,
    ) || [];

  const [createdOidcConnection, setCreatedOidcConnection] =
    useState<OIDCConnection | null>(null);
  const navigate = useNavigate();

  async function handleCreate() {
    try {
      const { oidcConnection } = await createOidcConnectionMutation.mutateAsync(
        {
          oidcConnection: {},
        },
      );
      if (!oidcConnection) {
        toast.error(
          "Failed to create OIDC connection. Please try again later.",
        );
        return;
      }

      await refetch();
      navigate(`/organization/oidc-connections/${oidcConnection.id}`);
      toast.success("OIDC connection created successfully.");
    } catch {
      toast.error("Failed to create OIDC connection. Please try again later.");
    }
  }

  useEffect(() => {
    if (createdOidcConnection) {
      setTimeout(() => {
        setCreatedOidcConnection(null);
      }, 2000);
    }
  }, [createdOidcConnection, setCreatedOidcConnection]);

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-col lg:flex-row items-start justify-between space-y-2 lg:space-y-0">
        <div className="space-y-2">
          <CardTitle>OIDC Connections</CardTitle>
          <CardDescription>
            Configure OIDC identify providers for your Organization.
          </CardDescription>
        </div>
        <CardAction className="w-full lg:w-auto">
          <Button
            onClick={handleCreate}
            disabled={createOidcConnectionMutation.isPending}
            size="sm"
          >
            <Plus className="mr-2" />
            Create OIDC Connection
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {oidcConnections.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Connection</TableHead>
                    <TableHead>Issuer</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {oidcConnections.map((oidcConnection) => (
                    <TableRow
                      key={oidcConnection.id}
                      className={cn(
                        createdOidcConnection?.id === oidcConnection.id
                          ? "bg-muted animate-pulse duration-500"
                          : "",
                      )}
                    >
                      <TableCell>
                        <span className="px-2 py-1 bg-muted font-mono text-xs rounded mr-2">
                          {oidcConnection.id}
                        </span>
                        {oidcConnection.primary && (
                          <Badge variant="outline">Primary</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <Issuer
                          configurationUrl={oidcConnection.configurationUrl}
                        />
                      </TableCell>
                      <TableCell>
                        {oidcConnection.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(oidcConnection.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManageOidcConnectionButton
                          oidcConnection={oidcConnection}
                        />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <div className="text-center text-muted-foreground text-sm pt-8">
                No OIDC Connections found. Create one to get started.
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

function Issuer({ configurationUrl }: { configurationUrl: string }) {
  const [issuer, setIssuer] = useState<string | null>(null);
  useEffect(() => {
    if (!configurationUrl) {
      setIssuer(null);
      return;
    }
    (async () => {
      try {
        const response = await fetch(configurationUrl);
        if (!response.ok) {
          setIssuer("—");
          return;
        }
        const data = await response.json();
        setIssuer(data.issuer);
      } catch (error) {
        console.error("Failed to fetch OIDC issuer:", error);
        setIssuer("—");
      }
    })();
  }, [configurationUrl]);

  return <>{issuer || "—"}</>;
}

function ManageOidcConnectionButton({
  oidcConnection,
}: {
  oidcConnection: OIDCConnection;
}) {
  const { refetch } = useInfiniteQuery(
    listOIDCConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deleteOidcConnectionMutation = useMutation(deleteOIDCConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete(oidcConnection: OIDCConnection) {
    try {
      await deleteOidcConnectionMutation.mutateAsync({
        id: oidcConnection.id,
      });

      await refetch();
      setDeleteOpen(false);
      toast.success("OIDC connection deleted successfully.");
    } catch {
      toast.error("Failed to delete OIDC connection. Please try again later.");
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
            <Link to={`/organization/oidc-connections/${oidcConnection.id}`}>
              <AlignLeft />
              Details
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
              Deleting this OIDC Connection is permanent and connot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel asChild>
              <Button variant="outline">Cancel</Button>
            </AlertDialogCancel>
            <Button
              onClick={() => handleDelete(oidcConnection)}
              variant="destructive"
            >
              Delete OIDC Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
