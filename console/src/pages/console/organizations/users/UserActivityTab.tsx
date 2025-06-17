import React from "react";
import { useParams } from "react-router";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ConsoleListAuditLogEventsRequest } from "@/gen/tesseral/backend/v1/backend_pb";

export function UserActivityTab() {
  const { userId } = useParams();

  return (
    <Card>
      <CardHeader>
        <CardTitle>User Activity</CardTitle>
        <CardDescription>
          View Logs of actions taken by this User.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ListAuditLogEventsTable
          listParams={
            {
              actorUserId: userId,
              pageToken: "",
            } as ConsoleListAuditLogEventsRequest
          }
        />
      </CardContent>
    </Card>
  );
}
