import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { AlignLeft, Plus, Settings, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
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
  createOIDCConnection,
  deleteOIDCConnection,
  getOrganization,
  listOIDCConnections,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { OIDCConnection } from "@/gen/tesseral/backend/v1/models_pb";

export function ListOrganizationOidcConnectionsCard() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const {
    data: listOidcConnectionsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listOIDCConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const oidcConnections =
    listOidcConnectionsResponses?.pages?.flatMap(
      (page) => page.oidcConnections,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>OIDC Connections</CardTitle>
        <CardDescription>
          Configure OIDC identify providers for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}.
          </span>
        </CardDescription>
        <CardAction>
          <CreateOidcConnectionButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {oidcConnections.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-6">
                No OIDC connections found. Create a new connection to get
                started.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Connection</TableHead>
                    <TableHead>Issuer</TableHead>
                    <TableHead>Client ID</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {oidcConnections.map((oidcConnection) => (
                    <TableRow key={oidcConnection.id}>
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
                      <TableCell>{oidcConnection.clientId || "—"}</TableCell>
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

function CreateOidcConnectionButton() {
  const { organizationId } = useParams();
  const navigate = useNavigate();

  const { data: listOidcConnectionsResponse, refetch } = useInfiniteQuery(
    listOIDCConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createOidcConnectionMutation = useMutation(createOIDCConnection);

  async function handleCreateOidcConnection() {
    const isFirstOidcConnection =
      listOidcConnectionsResponse?.pages[0]?.oidcConnections.length === 0;
    const newOidcConnection = await createOidcConnectionMutation.mutateAsync({
      oidcConnection: {
        organizationId,
        primary: isFirstOidcConnection,
      },
    });

    if (!newOidcConnection.oidcConnection?.id) {
      toast.error("Failed to create OIDC Connection. Please try again.");
      return;
    }

    await refetch();
    toast.success("OIDC Connection created successfully");
    navigate(
      `/organizations/${organizationId}/oidc-connections/${newOidcConnection.oidcConnection.id}`,
    );
  }

  return (
    <Button onClick={handleCreateOidcConnection}>
      <Plus />
      Create OIDC Connection
    </Button>
  );
}

function ManageOidcConnectionButton({
  oidcConnection,
}: {
  oidcConnection: OIDCConnection;
}) {
  const { organizationId } = useParams();

  const { refetch } = useInfiniteQuery(
    listOIDCConnections,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const deleteOidcConnectionMutation = useMutation(deleteOIDCConnection);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteOidcConnectionMutation.mutateAsync({
      id: oidcConnection.id,
    });
    await refetch();
    toast.success("OIDC Connection deleted successfully");
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
              to={`/organizations/${organizationId}/oidc-connections/${oidcConnection.id}`}
            >
              <div className="w-full flex items-center gap-2">
                <AlignLeft />
                <span>OIDC Connection Details</span>
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
              Delete OIDC Connection
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
              <span className="font-semibold">{oidcConnection.id}</span> OIDC
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
