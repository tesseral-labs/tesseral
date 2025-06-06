import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery } from "@connectrpc/connect-query";
import { DateTime } from "luxon";
import React from "react";
import { Link } from "react-router-dom";

import {
  Card,
  CardContent,
  CardDescription,
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
import { listSAMLConnections } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { CreateSAMLConnectionButton } from "./saml-connections/CreateSAMLConnectionButton";

export function OrganizationSAMLConnectionsTab() {
  const {
    data: listSAMLConnectionsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listSAMLConnections,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const samlConnections = listSAMLConnectionsResponses?.pages?.flatMap(
    (page) => page.samlConnections,
  );

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center space-x-4">
          <div className="space-y-2">
            <CardTitle>SAML Connections</CardTitle>
            <CardDescription>
              A SAML connection is a link your enterprise Identity Provider.
            </CardDescription>
          </div>
          <CreateSAMLConnectionButton />
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {samlConnections?.map((samlConnection) => (
                <TableRow key={samlConnection.id}>
                  <TableCell>
                    <Link
                      className="font-mono font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organization-settings/saml-connections/${samlConnection.id}`}
                    >
                      {samlConnection.id}
                    </Link>
                  </TableCell>
                  <TableCell>
                    {samlConnection.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(samlConnection.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {samlConnection.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(samlConnection.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
              {isFetchingNextPage && (
                <TableRow>
                  <TableCell colSpan={3}>Loading...</TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>

          {hasNextPage && (
            <div className="flex justify-center mt-4 mb-6">
              <button
                className="btn btn-primary"
                onClick={() => fetchNextPage()}
              >
                Load More
              </button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
