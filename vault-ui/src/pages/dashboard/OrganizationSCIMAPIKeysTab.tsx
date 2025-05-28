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
import { listSCIMAPIKeys } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { CreateSCIMAPIKeyButton } from "./scim-api-keys/CreateSCIMAPIKeyButton";

export function OrganizationSCIMAPIKeysTab() {
  const {
    data: listSCIMAPIKeysResponse,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listSCIMAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const scimAPIKeys = listSCIMAPIKeysResponse?.pages?.flatMap(
    (page) => page.scimApiKeys,
  );

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex-row justify-between items-center space-x-4">
          <div className="space-y-2">
            <CardTitle>SCIM API Keys</CardTitle>
            <CardDescription>
              A SCIM API key allows for enterprise directory syncing.
            </CardDescription>
          </div>

          <CreateSCIMAPIKeyButton />
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Display Name</TableHead>
                <TableHead>ID</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created At</TableHead>
                <TableHead>Updated At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {scimAPIKeys?.map((scimAPIKey) => (
                <TableRow key={scimAPIKey.id}>
                  <TableCell>
                    <Link
                      className="font-mono font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organization-settings/scim-api-keys/${scimAPIKey.id}`}
                    >
                      {scimAPIKey.displayName}
                    </Link>
                  </TableCell>
                  <TableCell>
                    <Link
                      className="font-mono font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organization-settings/scim-api-keys/${scimAPIKey.id}`}
                    >
                      {scimAPIKey.id}
                    </Link>
                  </TableCell>
                  <TableCell>
                    {scimAPIKey.revoked ? "Revoked" : "Active"}
                  </TableCell>
                  <TableCell>
                    {scimAPIKey.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(scimAPIKey.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {scimAPIKey.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(scimAPIKey.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {hasNextPage && (
            <div className="flex justify-center mt-4">
              <button
                className="btn btn-primary"
                onClick={() => fetchNextPage()}
                disabled={isFetchingNextPage}
              >
                {isFetchingNextPage ? "Loading..." : "Load More"}
              </button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
