import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery } from "@connectrpc/connect-query";
import React from "react";
import { useParams } from "react-router";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { consoleListAuditLogEvents } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleListAuditLogEventsRequest } from "@/gen/tesseral/backend/v1/backend_pb";
import { AuditLogEventResourceType } from "@/gen/tesseral/backend/v1/models_pb";

export function UserHistoryTab() {
  const { userId } = useParams();

  return (
    <Card>
      <CardHeader>
        <CardTitle>User History</CardTitle>
        <CardDescription>View Logs of changes to this User.</CardDescription>
      </CardHeader>
      <CardContent>
        <ListAuditLogEventsTable
          listParams={
            {
              resourceId: userId as string,
              resourceType: AuditLogEventResourceType.USER,
              pageToken: "",
            } as ConsoleListAuditLogEventsRequest
          }
        />
      </CardContent>
    </Card>
  );
}
