import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useQuery } from "@connectrpc/connect-query";
import { AlignLeft, Copy } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, useParams } from "react-router";
import { toast } from "sonner";

import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  getUser,
  listSessions,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { PrimaryAuthFactor } from "@/gen/tesseral/backend/v1/models_pb";

export function UserSessionsTab() {
  const { organizationId, userId } = useParams();
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const {
    data: listSessionsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listSessions,
    {
      userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const sessions =
    listSessionsResponses?.pages.flatMap((page) => page.sessions) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>User Sessions</CardTitle>
        <CardDescription>
          Sessions created by {getUserResponse?.user?.email}.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {sessions.length === 0 ? (
              <div className="text-center text-muted-foreground py-6">
                No sessions found for this User.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>Auth Method</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Last Active</TableHead>
                    <TableHead>Expires</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {sessions.map((session) => (
                    <TableRow key={session.id}>
                      <TableCell>
                        <span
                          className="bg-muted text-muted-foreground px-2 py-1 rounded text-xs font-mono cursor-pointer"
                          onClick={() => {
                            navigator.clipboard.writeText(session.id);
                            toast.success("Session ID copied to clipboard");
                          }}
                        >
                          {session.id}
                          <Copy className="inline w-3 h-3 ml-1" />
                        </span>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {primaryAuthFactorLabel(session.primaryAuthFactor)}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {session.revoked ? (
                          <Badge variant="secondary">Revoked</Badge>
                        ) : (
                          <Badge>Active</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {session.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(session.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {session.lastActiveTime &&
                          DateTime.fromJSDate(
                            timestampDate(session.lastActiveTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {session.expireTime &&
                          DateTime.fromJSDate(
                            timestampDate(session.expireTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <Link
                          to={`/organizations/${organizationId}/users/${userId}/sessions/${session.id}`}
                        >
                          <Button variant="outline" size="sm">
                            <AlignLeft />
                            Session Details
                          </Button>
                        </Link>
                      </TableCell>
                    </TableRow>
                  ))}
                  {isFetchingNextPage && (
                    <TableRow>
                      <TableCell colSpan={6}>
                        Loading more sessions...
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            )}
          </>
        )}

        {hasNextPage && (
          <CardFooter>
            <Button onClick={() => fetchNextPage()} variant="outline" size="sm">
              Load More
            </Button>
          </CardFooter>
        )}
      </CardContent>
    </Card>
  );
}

// Handles proper display of primary auth factor labels, e.g. `Oidc` -> `OIDC`.
function primaryAuthFactorLabel(primaryAuthFactor: PrimaryAuthFactor) {
  switch (primaryAuthFactor) {
    case PrimaryAuthFactor.EMAIL:
      return "Email";
    case PrimaryAuthFactor.GOOGLE:
      return "Google";
    case PrimaryAuthFactor.GITHUB:
      return "GitHub";
    case PrimaryAuthFactor.MICROSOFT:
      return "Microsoft";
    case PrimaryAuthFactor.SAML:
      return "SAML";
    case PrimaryAuthFactor.OIDC:
      return "OIDC";
    case PrimaryAuthFactor.IMPERSONATION:
      return "Impersonation";
    default:
      return "Unknown";
  }
}
